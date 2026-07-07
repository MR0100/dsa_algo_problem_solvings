package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: BFS Topological Sort (Kahn's Algorithm) ──────────────────────
//
// bfsTopoSort solves Alien Dictionary by building a precedence graph from
// adjacent word pairs and ordering letters by Kahn's algorithm.
//
// Intuition:
//
//	The dictionary is sorted by an unknown alphabet. Comparing two adjacent
//	words, the first position where they differ reveals one ordering rule:
//	the earlier word's char comes before the later word's char. Collect all
//	such rules into a directed graph over letters, then any topological order
//	of that graph is a valid alien alphabet. If the graph has a cycle the
//	ordering is contradictory -> return "".
//
// Algorithm:
//  1. Seed every letter that appears with in-degree 0 (so lone letters are
//     included in the output).
//  2. For each adjacent word pair, find the first differing char c1,c2 and add
//     edge c1 -> c2 (once). Handle the invalid prefix case: if word1 is longer
//     than word2 but is a prefix of it (e.g. "abc" before "ab") return "".
//  3. Kahn: push all in-degree-0 letters onto a queue; repeatedly pop a letter,
//     append it, and decrement neighbours, enqueuing those that hit 0.
//  4. If the result includes every letter, return it; else a cycle exists -> "".
//
// Time:  O(C) — C = total characters across all words (graph build) plus
//
//	O(V + E) for the sort, with V ≤ 26 letters and E ≤ 26² edges.
//
// Space: O(V + E) — adjacency + in-degree over the ≤ 26 letters.
func bfsTopoSort(words []string) string {
	adj := make(map[byte]map[byte]struct{}) // c -> set of letters after c
	indeg := make(map[byte]int)             // c -> number of prerequisites

	// Seed all appearing letters with in-degree 0.
	for _, w := range words {
		for i := 0; i < len(w); i++ {
			if _, ok := adj[w[i]]; !ok {
				adj[w[i]] = make(map[byte]struct{})
				indeg[w[i]] = 0
			}
		}
	}

	// Derive ordering rules from adjacent word pairs.
	for i := 0; i+1 < len(words); i++ {
		w1, w2 := words[i], words[i+1]
		minLen := len(w1)
		if len(w2) < minLen {
			minLen = len(w2)
		}
		j := 0
		for j < minLen && w1[j] == w2[j] {
			j++ // skip the common prefix
		}
		if j == minLen {
			// No differing char within the shared length. Invalid only if the
			// first word is the longer one ("abc" cannot precede "ab").
			if len(w1) > len(w2) {
				return ""
			}
			continue
		}
		c1, c2 := w1[j], w2[j] // first difference gives rule c1 -> c2
		if _, exists := adj[c1][c2]; !exists {
			adj[c1][c2] = struct{}{}
			indeg[c2]++ // c2 gains a prerequisite
		}
	}

	// Kahn's algorithm. We keep the ready set sorted so the output is
	// deterministic (lexicographically smallest valid order).
	ready := make([]byte, 0)
	for c, d := range indeg {
		if d == 0 {
			ready = append(ready, c) // no prerequisites -> ready
		}
	}
	sort.Slice(ready, func(i, j int) bool { return ready[i] < ready[j] })
	result := make([]byte, 0, len(indeg))
	for len(ready) > 0 {
		c := ready[0] // smallest available letter
		ready = ready[1:]
		result = append(result, c) // commit this letter
		freed := make([]byte, 0)
		for next := range adj[c] {
			indeg[next]--
			if indeg[next] == 0 {
				freed = append(freed, next) // all its prereqs placed
			}
		}
		sort.Slice(freed, func(i, j int) bool { return freed[i] < freed[j] })
		ready = append(ready, freed...)
		sort.Slice(ready, func(i, j int) bool { return ready[i] < ready[j] })
	}
	if len(result) < len(indeg) {
		return "" // some letters never freed -> cycle
	}
	return string(result)
}

// ── Approach 2: DFS Topological Sort (Optimal) ───────────────────────────────
//
// dfsTopoSort solves Alien Dictionary with the same precedence graph but orders
// letters via post-order DFS, detecting cycles with a colour marking.
//
// Intuition:
//
//	A topological order is the reverse of DFS finish times. Colour nodes
//	white (unseen) / grey (in current path) / black (done). Encountering a
//	grey node again means a back edge -> cycle -> "". Otherwise append each
//	node after its descendants and reverse at the end.
//
// Algorithm:
//  1. Build the same adjacency graph and letter set.
//  2. DFS each letter; on a grey re-visit signal a cycle. Post-append blacks.
//  3. Reverse the post-order to get the alphabet.
//
// Time:  O(C + V + E) — build + traversal, V ≤ 26, E ≤ 26².
// Space: O(V + E) — graph + recursion stack.
func dfsTopoSort(words []string) string {
	adj := make(map[byte]map[byte]struct{})
	letters := make(map[byte]struct{})
	for _, w := range words {
		for i := 0; i < len(w); i++ {
			letters[w[i]] = struct{}{}
			if _, ok := adj[w[i]]; !ok {
				adj[w[i]] = make(map[byte]struct{})
			}
		}
	}
	for i := 0; i+1 < len(words); i++ {
		w1, w2 := words[i], words[i+1]
		minLen := len(w1)
		if len(w2) < minLen {
			minLen = len(w2)
		}
		j := 0
		for j < minLen && w1[j] == w2[j] {
			j++
		}
		if j == minLen {
			if len(w1) > len(w2) {
				return "" // invalid prefix ordering
			}
			continue
		}
		adj[w1[j]][w2[j]] = struct{}{}
	}

	const (
		white = 0 // unvisited
		grey  = 1 // on current DFS path
		black = 2 // fully processed
	)
	color := make(map[byte]int)
	order := make([]byte, 0, len(letters))
	// sortedKeys returns a map's byte keys in ascending order for determinism.
	sortedKeys := func(m map[byte]struct{}) []byte {
		ks := make([]byte, 0, len(m))
		for k := range m {
			ks = append(ks, k)
		}
		sort.Slice(ks, func(i, j int) bool { return ks[i] < ks[j] })
		return ks
	}
	var dfs func(c byte) bool
	dfs = func(c byte) bool {
		color[c] = grey // enter the recursion stack
		for _, next := range sortedKeys(adj[c]) {
			if color[next] == grey {
				return false // back edge -> cycle
			}
			if color[next] == white && !dfs(next) {
				return false // cycle found deeper
			}
		}
		color[c] = black         // finished
		order = append(order, c) // post-order append
		return true
	}
	for _, c := range sortedKeys(letters) {
		if color[c] == white {
			if !dfs(c) {
				return "" // cyclic constraints
			}
		}
	}
	// order is reverse topological; flip it.
	for i, j := 0, len(order)-1; i < j; i, j = i+1, j-1 {
		order[i], order[j] = order[j], order[i]
	}
	return string(order)
}

func main() {
	fmt.Println("=== Approach 1: BFS Topological Sort ===")
	fmt.Println(bfsTopoSort([]string{"wrt", "wrf", "er", "ett", "rftt"})) // expected wertf
	fmt.Println(bfsTopoSort([]string{"z", "x"}))                          // expected zx
	fmt.Println(bfsTopoSort([]string{"z", "x", "z"}))                     // expected (empty)

	fmt.Println("=== Approach 2: DFS Topological Sort (Optimal) ===")
	fmt.Println(dfsTopoSort([]string{"wrt", "wrf", "er", "ett", "rftt"})) // expected wertf
	fmt.Println(dfsTopoSort([]string{"z", "x"}))                          // expected zx
	fmt.Println(dfsTopoSort([]string{"z", "x", "z"}))                     // expected (empty)
}
