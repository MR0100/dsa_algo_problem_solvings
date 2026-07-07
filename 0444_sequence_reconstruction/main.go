package main

import "fmt"

// ── Approach 1: Topological Sort — Uniqueness Check (Kahn's BFS, Optimal) ─────
//
// topoUniqueness solves Sequence Reconstruction by asking whether the DAG built
// from all pairs in seqs has EXACTLY ONE topological order, and whether that
// order equals nums.
//
// Intuition:
//
//	Each seq imposes ordering constraints on consecutive elements: a before b.
//	Collecting every such edge gives a DAG. nums is the UNIQUE reconstruction iff
//	Kahn's algorithm always has exactly one node with in-degree 0 to choose at
//	each step (no ambiguity in ordering), the produced order equals nums, and
//	every value 1..n actually appears. If ever two nodes have in-degree 0 at the
//	same time, at least two valid orders exist → not unique.
//
// Algorithm:
//  1. Build adjacency list and in-degree over all values seen in seqs; validate
//     each value is within [1, n].
//  2. Seed a queue with all in-degree-0 nodes.
//  3. Repeatedly: if queue size != 1, reconstruction is ambiguous or stuck →
//     false. Pop the single node; it must equal nums[index]. Decrement its
//     neighbours' in-degrees, enqueue any that hit 0.
//  4. After processing, success iff we emitted exactly n nodes in nums's order.
//
// Time:  O(V + E) — V = n distinct values, E = total pairs across seqs.
// Space: O(V + E) — adjacency list, in-degree map, and queue.
func topoUniqueness(nums []int, sequences [][]int) bool {
	n := len(nums)

	adj := make(map[int]map[int]struct{}) // a -> set of b (dedup parallel edges)
	indeg := make(map[int]int)            // value -> number of prerequisites
	seen := make(map[int]bool)            // which values appear in seqs at all

	// Register a value node the first time we encounter it.
	touch := func(v int) {
		if !seen[v] {
			seen[v] = true
			indeg[v] = 0
			adj[v] = make(map[int]struct{})
		}
	}

	total := 0 // count of distinct values seen across seqs (must equal n)
	for _, seq := range sequences {
		for _, v := range seq {
			// Any value outside [1, n] means seqs can't reconstruct a
			// permutation of 1..n → impossible.
			if v < 1 || v > n {
				return false
			}
			if !seen[v] {
				touch(v)
				total++
			}
		}
		// Add an edge for every consecutive pair (dedup to keep in-degree exact).
		for i := 0; i+1 < len(seq); i++ {
			a, b := seq[i], seq[i+1]
			if _, exists := adj[a][b]; !exists {
				adj[a][b] = struct{}{}
				indeg[b]++ // b now has one more prerequisite
			}
		}
	}

	// If not every value 1..n showed up, nums cannot be reconstructed at all.
	if total != n {
		return false
	}

	// Seed queue with all sources (in-degree 0).
	queue := make([]int, 0, n)
	for v := range seen {
		if indeg[v] == 0 {
			queue = append(queue, v)
		}
	}

	index := 0 // position in nums we expect to emit next
	for len(queue) > 0 {
		// Unique next element required: more than one source ⇒ ambiguous order.
		if len(queue) > 1 {
			return false
		}
		node := queue[0]
		queue = queue[1:]

		// The forced next node must match nums at this position.
		if index >= n || node != nums[index] {
			return false
		}
		index++

		// Relax outgoing edges; enqueue neighbours that become sources.
		for nb := range adj[node] {
			indeg[nb]--
			if indeg[nb] == 0 {
				queue = append(queue, nb)
			}
		}
	}

	// Unique reconstruction iff we emitted all n values (order already checked).
	return index == n
}

// ── Approach 2: Adjacent-Coverage Check (Constraint Counting) ────────────────
//
// adjacentCoverage solves Sequence Reconstruction with a clever observation that
// avoids an explicit topological sort.
//
// Intuition:
//
//	nums is the UNIQUE supersequence iff two conditions hold:
//	  (a) seqs pins down the order of EVERY adjacent pair (nums[i], nums[i+1]) —
//	      i.e. some seq lists nums[i] immediately before nums[i+1]. If even one
//	      adjacent pair is unconstrained, its two elements could swap, giving a
//	      second valid order.
//	  (b) Every value in seqs is a valid index into nums (in [1, n]) and every
//	      value 1..n is covered.
//	Build pos[v] = index of v in nums. A pair (a, b) that appears consecutively in
//	some seq must satisfy pos[a] + 1 == pos[b] to be consistent with nums; if it
//	is exactly one apart, it "covers" that adjacency. Count how many of the n-1
//	adjacencies get covered; all covered ⇒ unique.
//
// Algorithm:
//  1. Build pos map from nums; require 1..n present.
//  2. For each consecutive pair (a, b) in seqs: validate values in range and
//     present. If pos[a] > pos[b], seqs contradict nums → false. If
//     pos[a] + 1 == pos[b], mark that adjacency as covered.
//  3. All n-1 adjacencies covered ⇒ true.
//
// Time:  O(V + E) — build pos in O(n), scan all pairs once.
// Space: O(n) — position map and coverage flags.
func adjacentCoverage(nums []int, sequences [][]int) bool {
	n := len(nums)
	pos := make(map[int]int, n) // value -> its index in nums
	for i, v := range nums {
		pos[v] = i
	}

	covered := make([]bool, n)    // covered[i] = adjacency (nums[i], nums[i+1]) pinned
	toCover := n - 1              // number of adjacencies still needing a constraint
	sawValue := make([]bool, n+1) // sawValue[v] = value v appeared in seqs

	for _, seq := range sequences {
		for i := 0; i < len(seq); i++ {
			v := seq[i]
			// Value must index into nums (1..n); otherwise reconstruction fails.
			if v < 1 || v > n {
				return false
			}
			sawValue[v] = true

			if i+1 < len(seq) {
				a, b := seq[i], seq[i+1]
				if b < 1 || b > n {
					return false
				}
				// If a is supposed to come after b in nums, seqs contradict it.
				if pos[a] > pos[b] {
					return false
				}
				// Exactly adjacent in nums and not yet counted → cover it.
				if pos[a]+1 == pos[b] && !covered[pos[a]] {
					covered[pos[a]] = true
					toCover--
				}
			}
		}
	}

	// Every value 1..n must have appeared somewhere in seqs.
	for v := 1; v <= n; v++ {
		if !sawValue[v] {
			return false
		}
	}

	// Unique iff all n-1 adjacencies were pinned by some seq.
	return toCover == 0
}

func main() {
	fmt.Println("=== Approach 1: Topological Sort — Uniqueness (Kahn BFS) ===")
	fmt.Printf("nums=[1 2 3] seqs=[[1 2] [1 3]]           got=%v  expected false\n",
		topoUniqueness([]int{1, 2, 3}, [][]int{{1, 2}, {1, 3}}))
	fmt.Printf("nums=[1 2 3] seqs=[[1 2]]                 got=%v  expected false\n",
		topoUniqueness([]int{1, 2, 3}, [][]int{{1, 2}}))
	fmt.Printf("nums=[1 2 3] seqs=[[1 2] [1 3] [2 3]]     got=%v  expected true\n",
		topoUniqueness([]int{1, 2, 3}, [][]int{{1, 2}, {1, 3}, {2, 3}}))
	fmt.Printf("nums=[4 1 5 2 6 3] seqs=[[5 2 6 3][4 1 5 2]] got=%v  expected true\n",
		topoUniqueness([]int{4, 1, 5, 2, 6, 3}, [][]int{{5, 2, 6, 3}, {4, 1, 5, 2}}))

	fmt.Println("=== Approach 2: Adjacent-Coverage Check ===")
	fmt.Printf("nums=[1 2 3] seqs=[[1 2] [1 3]]           got=%v  expected false\n",
		adjacentCoverage([]int{1, 2, 3}, [][]int{{1, 2}, {1, 3}}))
	fmt.Printf("nums=[1 2 3] seqs=[[1 2]]                 got=%v  expected false\n",
		adjacentCoverage([]int{1, 2, 3}, [][]int{{1, 2}}))
	fmt.Printf("nums=[1 2 3] seqs=[[1 2] [1 3] [2 3]]     got=%v  expected true\n",
		adjacentCoverage([]int{1, 2, 3}, [][]int{{1, 2}, {1, 3}, {2, 3}}))
	fmt.Printf("nums=[4 1 5 2 6 3] seqs=[[5 2 6 3][4 1 5 2]] got=%v  expected true\n",
		adjacentCoverage([]int{4, 1, 5, 2, 6, 3}, [][]int{{5, 2, 6, 3}, {4, 1, 5, 2}}))
}
