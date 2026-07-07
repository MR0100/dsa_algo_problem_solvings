package main

import "fmt"

// ── Approach 1: Breadth-First Search (Shortest Path) ─────────────────────────
//
// bfs finds the minimum number of single-character mutations from startGene to
// endGene, where every intermediate gene must appear in bank.
//
// Intuition:
//
//	Model genes as graph nodes; connect two genes with an edge when they differ
//	in exactly one character. A "mutation" is one edge. The fewest mutations is
//	then the shortest path from startGene to endGene in this graph. Unweighted
//	shortest path = BFS: explore all genes 1 mutation away, then 2 away, etc.;
//	the first time we pop endGene, the current level is the answer.
//
// Algorithm:
//
//  1. Put bank into a set; if endGene isn't in it, return -1 (unreachable).
//  2. BFS from startGene, tracking a step count per level.
//  3. To expand a gene, try every position × every letter in {A,C,G,T}; a
//     neighbour is any resulting string that is still in the (unused) bank set.
//  4. Mark neighbours used (delete from the set) so we never revisit them.
//  5. Return the level at which endGene is first reached; -1 if BFS drains.
//
// Time:  O(N · L · 4) = O(N · L) — N genes, each expanded by trying L positions
//
//	× 4 letters; L = 8 here. Building each candidate string costs O(L).
//
// Space: O(N · L) — the bank set and BFS queue hold up to N genes of length L.
func bfs(startGene, endGene string, bank []string) int {
	valid := make(map[string]bool, len(bank)) // set of usable target genes
	for _, g := range bank {
		valid[g] = true
	}
	if !valid[endGene] {
		return -1 // target not in bank → impossible
	}
	choices := []byte{'A', 'C', 'G', 'T'} // the only allowed characters
	queue := []string{startGene}          // BFS frontier, starts at the source
	steps := 0                            // mutations used to reach the current level
	for len(queue) > 0 {
		next := []string{} // genes reachable in steps+1 mutations
		for _, gene := range queue {
			if gene == endGene {
				return steps // popped the target: shortest distance found
			}
			b := []byte(gene) // mutable copy to try single-char edits
			for i := 0; i < len(b); i++ {
				original := b[i] // remember to restore after trying each letter
				for _, c := range choices {
					if c == original {
						continue // no mutation if the char is unchanged
					}
					b[i] = c
					cand := string(b)
					if valid[cand] { // neighbour exists and is unused
						valid[cand] = false // mark visited to avoid cycles/repeats
						next = append(next, cand)
					}
				}
				b[i] = original // restore this position before moving on
			}
		}
		queue = next // advance to the next BFS level
		steps++      // one more mutation applied
	}
	return -1 // exhausted all reachable genes without hitting endGene
}

// ── Approach 2: Bidirectional BFS (Optimal) ──────────────────────────────────
//
// bidirectionalBFS runs two BFS frontiers — one growing out from startGene, one
// from endGene — and stops as soon as they intersect.
//
// Intuition:
//
//	A single BFS explores a ball of radius d around the start, whose size grows
//	exponentially with d. Searching from both ends and meeting in the middle
//	means each side only needs radius ~d/2, cutting the explored volume roughly
//	from b^d to 2·b^(d/2). At every step we expand the SMALLER frontier (fewer
//	nodes to branch from), which keeps the search balanced and fast. The moment
//	a newly generated neighbour lies in the opposite frontier, the two paths join
//	and the total mutation count is known.
//
// Algorithm:
//
//  1. Bank → set; if endGene not present, return -1.
//  2. Maintain two frontier sets, begin={start} and end={end}, plus a visited
//     set. steps starts at 0.
//  3. Loop while both frontiers are non-empty: always expand the smaller one.
//     For each gene, generate all valid one-char neighbours in the bank:
//     • if a neighbour is in the OTHER frontier → return steps+1 (they meet).
//     • else if unvisited → add to the next frontier and mark visited.
//     Replace the expanded frontier with the next one; steps++.
//  4. If the loop ends without meeting, return -1.
//
// Time:  O(N · L) worst case, but explores far fewer nodes than one-directional
//
//	BFS in practice (√ of the search volume).
//
// Space: O(N · L) — two frontier sets plus the visited/bank sets.
func bidirectionalBFS(startGene, endGene string, bank []string) int {
	valid := make(map[string]bool, len(bank))
	for _, g := range bank {
		valid[g] = true
	}
	if !valid[endGene] {
		return -1
	}
	if startGene == endGene {
		return 0 // already there (defensive; problem says they differ)
	}
	choices := []byte{'A', 'C', 'G', 'T'}
	begin := map[string]bool{startGene: true} // frontier growing from the start
	end := map[string]bool{endGene: true}     // frontier growing from the target
	visited := map[string]bool{startGene: true, endGene: true}
	steps := 0
	for len(begin) > 0 && len(end) > 0 {
		// Always expand the smaller frontier to keep branching minimal.
		if len(begin) > len(end) {
			begin, end = end, begin // swap so `begin` is the smaller side
		}
		next := map[string]bool{} // the smaller frontier's next layer
		for gene := range begin {
			b := []byte(gene)
			for i := 0; i < len(b); i++ {
				original := b[i]
				for _, c := range choices {
					if c == original {
						continue
					}
					b[i] = c
					cand := string(b)
					if end[cand] {
						return steps + 1 // frontiers meet: paths join here
					}
					if valid[cand] && !visited[cand] {
						visited[cand] = true // claim it so the other side won't re-expand
						next[cand] = true
					}
				}
				b[i] = original
			}
		}
		begin = next // the smaller frontier advances one layer
		steps++      // that advance cost one mutation
	}
	return -1 // no meeting point → unreachable
}

// ── Approach 3: DFS with Memo of Best Depth ──────────────────────────────────
//
// dfs explores mutation paths depth-first, remembering the best (smallest) step
// count with which each gene was reached to prune worse revisits.
//
// Intuition:
//
//	The same graph, walked depth-first. Plain DFS can loop forever and does not
//	naturally yield the shortest path, so we (a) track which genes are on the
//	current path to avoid cycles, and (b) memoise the minimum depth at which each
//	gene has been reached: if we arrive again at an equal-or-greater depth, that
//	branch cannot improve any answer, so we cut it. DFS is included to contrast
//	with BFS — it works, but needs this bookkeeping precisely because it lacks
//	BFS's built-in shortest-path guarantee.
//
// Algorithm:
//
//  1. Bank → set; if endGene absent, return -1.
//  2. Recurse from startGene at depth 0. At each gene:
//     • if gene == endGene, record depth as a candidate answer.
//     • else for every valid unused neighbour reached at a strictly better
//     depth, recurse at depth+1.
//  3. Use bestDepth[gene] to skip neighbours already reached as cheaply.
//  4. Return the minimum recorded depth, or -1 if none.
//
// Time:  O(N · L) amortised via the depth memo (each gene is meaningfully
//
//	expanded only when reached more cheaply than before).
//
// Space: O(N · L) for the sets/memo plus O(N) recursion depth.
func dfs(startGene, endGene string, bank []string) int {
	valid := make(map[string]bool, len(bank))
	for _, g := range bank {
		valid[g] = true
	}
	if !valid[endGene] {
		return -1
	}
	choices := []byte{'A', 'C', 'G', 'T'}
	bestDepth := map[string]int{} // gene → smallest depth it was reached at
	best := -1                    // best answer so far; -1 means "none yet"

	var explore func(gene string, depth int)
	explore = func(gene string, depth int) {
		if gene == endGene {
			// Reached the target: keep the smaller mutation count.
			if best == -1 || depth < best {
				best = depth
			}
			return
		}
		// Prune: if we've been at this gene at an equal/shallower depth, stop.
		if d, seen := bestDepth[gene]; seen && d <= depth {
			return
		}
		bestDepth[gene] = depth // record the best depth for this gene
		b := []byte(gene)
		for i := 0; i < len(b); i++ {
			original := b[i]
			for _, c := range choices {
				if c == original {
					continue
				}
				b[i] = c
				cand := string(b)
				if valid[cand] { // only step onto genes that exist in the bank
					explore(cand, depth+1) // go one mutation deeper
				}
			}
			b[i] = original
		}
	}
	explore(startGene, 0)
	return best
}

func main() {
	// Example 1: one mutation away.
	fmt.Println("=== Approach 1: Breadth-First Search ===")
	fmt.Println(bfs("AACCGGTT", "AACCGGTA", []string{"AACCGGTA"}))                         // expected 1
	fmt.Println(bfs("AACCGGTT", "AAACGGTA", []string{"AACCGGTA", "AACCGCTA", "AAACGGTA"})) // expected 2
	fmt.Println(bfs("AAAAACCC", "AACCCCCC", []string{"AAAACCCC", "AAACCCCC", "AACCCCCC"})) // expected 3
	fmt.Println(bfs("AACCGGTT", "AACCGGTA", []string{}))                                   // expected -1 (empty bank)

	fmt.Println("=== Approach 2: Bidirectional BFS (Optimal) ===")
	fmt.Println(bidirectionalBFS("AACCGGTT", "AACCGGTA", []string{"AACCGGTA"}))                         // expected 1
	fmt.Println(bidirectionalBFS("AACCGGTT", "AAACGGTA", []string{"AACCGGTA", "AACCGCTA", "AAACGGTA"})) // expected 2
	fmt.Println(bidirectionalBFS("AAAAACCC", "AACCCCCC", []string{"AAAACCCC", "AAACCCCC", "AACCCCCC"})) // expected 3
	fmt.Println(bidirectionalBFS("AACCGGTT", "AACCGGTA", []string{}))                                   // expected -1

	fmt.Println("=== Approach 3: DFS with Memo of Best Depth ===")
	fmt.Println(dfs("AACCGGTT", "AACCGGTA", []string{"AACCGGTA"}))                         // expected 1
	fmt.Println(dfs("AACCGGTT", "AAACGGTA", []string{"AACCGGTA", "AACCGCTA", "AAACGGTA"})) // expected 2
	fmt.Println(dfs("AAAAACCC", "AACCCCCC", []string{"AAAACCCC", "AAACCCCC", "AACCCCCC"})) // expected 3
	fmt.Println(dfs("AACCGGTT", "AACCGGTA", []string{}))                                   // expected -1
}
