package main

import "fmt"

// A frog crosses a river landing only on stones at the given (strictly
// increasing) positions. The first jump must be exactly 1 unit; if the last
// jump was k units, the next must be k-1, k, or k+1 (and forward, ≥ 1). Return
// whether the frog can reach the LAST stone.

// ── Approach 1: Plain Recursion (exponential) ────────────────────────────────
//
// bruteForceRecursion tries, from the current stone and last jump size, all
// three next jump sizes and recurses, using a position->index map to test
// whether a landing spot is actually a stone.
//
// Intuition:
//
//	State is (current position, last jump k). From here the frog may jump
//	k-1, k, or k+1. A jump succeeds only if it is ≥ 1 AND lands exactly on a
//	stone. Recurse on each successful jump; success is reaching the final
//	stone's position. No memo yet, so overlapping (position, k) states are
//	recomputed — exponential.
//
// Algorithm:
//  1. Build pos -> true set of stone positions and note the target position.
//  2. dfs(position, k): if position == target return true. For step in {k-1,k,k+1}
//     with step ≥ 1: if position+step is a stone, recurse; return true if any does.
//  3. Start dfs(stones[0], 0); the first real jump becomes step = 1 (0+1).
//
// Time:  O(3^n) worst case — three branches per stone, no reuse.
// Space: O(n) recursion depth + O(n) position set.
func bruteForceRecursion(stones []int) bool {
	target := stones[len(stones)-1]
	stoneSet := make(map[int]bool, len(stones)) // fast "is there a stone here?"
	for _, s := range stones {
		stoneSet[s] = true
	}

	var dfs func(position, k int) bool
	dfs = func(position, k int) bool {
		if position == target {
			return true // reached the far bank
		}
		// Try the three permitted next jump sizes.
		for _, step := range []int{k - 1, k, k + 1} {
			if step <= 0 {
				continue // must jump forward by at least 1 unit
			}
			next := position + step
			if stoneSet[next] && dfs(next, step) {
				return true // some continuation from `next` succeeds
			}
		}
		return false
	}
	// The problem fixes the first jump at 1: starting with k=0 makes k+1 == 1.
	return dfs(stones[0], 0)
}

// ── Approach 2: Top-Down DP (memoized recursion) ─────────────────────────────
//
// dpTopDown is the same search but caches results per (stoneIndex, k) so each
// state is solved once.
//
// Intuition:
//
//	The exponential blow-up is pure repetition: many jump paths arrive at the
//	same stone with the same last-jump k. Since the answer from (stone, k) is
//	fixed, memoize it. There are at most n stones and jump sizes bounded by n,
//	so the state space is O(n²).
//
// Algorithm:
//  1. Map each position to its stone index; memo keyed by (index, k).
//  2. dfs(index, k): if index is the last stone return true. For step in
//     {k-1,k,k+1}, step ≥ 1: if stones[index]+step is a stone at index j,
//     recurse dfs(j, step). Cache and return whether any branch succeeds.
//  3. Start dfs(0, 0).
//
// Time:  O(n²) — O(n²) states, each doing O(1) work over 3 transitions.
// Space: O(n²) memo + O(n) recursion depth.
func dpTopDown(stones []int) bool {
	n := len(stones)
	indexOf := make(map[int]int, n) // stone position -> its index
	for i, s := range stones {
		indexOf[s] = i
	}

	// memo[index][k] caches whether (stone index, last jump k) can finish.
	// k ranges 0..n (a jump can't exceed n stones), so width n+1 is safe.
	memo := make([]map[int]bool, n)
	for i := range memo {
		memo[i] = make(map[int]bool)
	}

	var dfs func(index, k int) bool
	dfs = func(index, k int) bool {
		if index == n-1 {
			return true // standing on the last stone
		}
		if v, seen := memo[index][k]; seen {
			return v // already solved this exact state
		}
		res := false
		for _, step := range []int{k - 1, k, k + 1} {
			if step <= 0 {
				continue
			}
			if j, ok := indexOf[stones[index]+step]; ok {
				if dfs(j, step) {
					res = true
					break // one successful continuation is enough
				}
			}
		}
		memo[index][k] = res // remember for next time
		return res
	}
	return dfs(0, 0)
}

// ── Approach 3: Bottom-Up DP with Reachable Jump Sets (Optimal) ───────────────
//
// dpBottomUp stores, for each stone, the set of jump sizes k by which the frog
// could have ARRIVED there, then propagates k-1, k, k+1 forward to later stones.
//
// Intuition:
//
//	Turn the recursion into forward propagation. reach[i] = set of last-jump
//	sizes with which stone i is reachable. Seed reach[0] = {0}. For each stone i
//	and each arrival jump k, the frog can leave by k-1, k, k+1 (≥1); if that
//	lands on a stone j, add the new jump size to reach[j]. The far bank is
//	reachable iff reach[last] is non-empty.
//
// Algorithm:
//  1. indexOf: position -> stone index. reach[i] = map of arrival jump sizes.
//  2. reach[0] = {0}. For i = 0..n-1, for each k in reach[i], for step in
//     {k-1,k,k+1} with step ≥ 1: if stones[i]+step is stone j (j>i), add step
//     to reach[j].
//  3. Return len(reach[n-1]) > 0.
//
// Time:  O(n²) — each stone holds O(n) distinct arrival jump sizes; 3 transitions each.
// Space: O(n²) — the reachable-jump sets across all stones.
func dpBottomUp(stones []int) bool {
	n := len(stones)
	indexOf := make(map[int]int, n)
	for i, s := range stones {
		indexOf[s] = i
	}

	// reach[i] holds every jump size by which stone i can be reached.
	reach := make([]map[int]bool, n)
	for i := range reach {
		reach[i] = make(map[int]bool)
	}
	reach[0][0] = true // start on stone 0, "arrived" with a phantom jump of 0

	for i := 0; i < n; i++ {
		for k := range reach[i] { // every way we could be standing on stone i
			for _, step := range []int{k - 1, k, k + 1} {
				if step <= 0 {
					continue // forward jumps only
				}
				if j, ok := indexOf[stones[i]+step]; ok && j > i {
					reach[j][step] = true // record how stone j was reached
				}
			}
		}
	}
	// The far bank is crossable iff we recorded any arrival on the last stone.
	return len(reach[n-1]) > 0
}

func main() {
	fmt.Println("=== Approach 1: Brute Force Recursion ===")
	fmt.Printf("[0,1,3,5,6,8,12,17] -> %v\n", bruteForceRecursion([]int{0, 1, 3, 5, 6, 8, 12, 17})) // expected true
	fmt.Printf("[0,1,2,3,4,8,9,11]  -> %v\n", bruteForceRecursion([]int{0, 1, 2, 3, 4, 8, 9, 11}))  // expected false

	fmt.Println("=== Approach 2: Top-Down DP ===")
	fmt.Printf("[0,1,3,5,6,8,12,17] -> %v\n", dpTopDown([]int{0, 1, 3, 5, 6, 8, 12, 17})) // expected true
	fmt.Printf("[0,1,2,3,4,8,9,11]  -> %v\n", dpTopDown([]int{0, 1, 2, 3, 4, 8, 9, 11}))  // expected false

	fmt.Println("=== Approach 3: Bottom-Up DP (Reachable Jump Sets) ===")
	fmt.Printf("[0,1,3,5,6,8,12,17] -> %v\n", dpBottomUp([]int{0, 1, 3, 5, 6, 8, 12, 17})) // expected true
	fmt.Printf("[0,1,2,3,4,8,9,11]  -> %v\n", dpBottomUp([]int{0, 1, 2, 3, 4, 8, 9, 11}))  // expected false
}
