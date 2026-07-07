package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Backtracking (Fill Four Sides) ───────────────────────────────
//
// backtracking solves Matchsticks to Square by trying to place every stick into
// one of the four sides of the square, backtracking whenever a placement can no
// longer lead to a valid square.
//
// Intuition:
//
//	A square has four equal sides, each of length total/4. Model the four sides
//	as four buckets. Recursively assign matchstick i to any bucket whose running
//	length would not exceed the target side. When all sticks are placed and all
//	four buckets equal the target, a square exists.
//
// Algorithm:
//  1. If total sum is 0 or not divisible by 4, return false immediately.
//  2. Let side = total/4. Keep an array sides[4] of running lengths.
//  3. dfs(i): if i == n, success (every bucket must already equal side).
//  4. For each bucket b: if sides[b]+sticks[i] <= side, add, recurse, undo.
//  5. Return true if any branch succeeds.
//
// Time:  O(4^n) — each of n sticks may branch into 4 buckets.
// Space: O(n) — recursion depth plus the 4-length buckets array.
func backtracking(matchsticks []int) bool {
	total := 0
	for _, m := range matchsticks {
		total += m // accumulate the perimeter
	}
	// A square needs a positive perimeter divisible by 4.
	if total == 0 || total%4 != 0 {
		return false
	}
	side := total / 4 // target length of every side
	sides := [4]int{} // running length of each of the four sides

	var dfs func(i int) bool
	dfs = func(i int) bool {
		if i == len(matchsticks) {
			// All sticks used; sides are valid iff each already reached `side`.
			// Because we never let a bucket exceed `side`, and the totals add up
			// to 4*side, reaching the end already guarantees all four equal side.
			return true
		}
		for b := 0; b < 4; b++ {
			// Only place stick i into bucket b if it still fits.
			if sides[b]+matchsticks[i] <= side {
				sides[b] += matchsticks[i] // place stick i on side b
				if dfs(i + 1) {
					return true // this placement led to a full square
				}
				sides[b] -= matchsticks[i] // undo — try the next bucket
			}
		}
		return false // stick i fit nowhere on a valid square
	}
	return dfs(0)
}

// ── Approach 2: Backtracking + Pruning (Sort Desc, Skip Duplicates) ──────────
//
// backtrackingPruned solves Matchsticks to Square with the same four-bucket
// search, sped up by two classic pruning rules that make it pass comfortably.
//
// Intuition:
//
//	Two prunes collapse the search tree massively:
//	  (1) Sort sticks in DESCENDING order. Placing big sticks first fails fast:
//	      a stick longer than `side` is rejected at depth 0, and large pieces
//	      constrain buckets early instead of late.
//	  (2) Skip symmetric buckets. If two buckets currently hold the same length,
//	      putting the stick in either is the same decision — try only the first.
//	      Also, if adding the stick makes a bucket exactly full or empty→value,
//	      identical empty buckets are interchangeable, so only try one empty one.
//
// Algorithm:
//  1. Same divisibility check and target side.
//  2. Sort matchsticks descending; if the largest > side, return false.
//  3. dfs(i, sides): for each bucket b, skip if sides[b] == sides[b'] already
//     tried (b' < b with equal value); place if it fits; recurse; undo.
//  4. Break out of the bucket loop after the first empty bucket is tried
//     (all empty buckets are equivalent).
//
// Time:  O(4^n) worst case, but pruning makes it near-instant for n ≤ 15.
// Space: O(n) — recursion stack.
func backtrackingPruned(matchsticks []int) bool {
	total := 0
	for _, m := range matchsticks {
		total += m
	}
	if total == 0 || total%4 != 0 {
		return false
	}
	side := total / 4
	// Sort descending so large sticks are placed first (fail fast).
	sort.Sort(sort.Reverse(sort.IntSlice(matchsticks)))
	// Largest stick cannot exceed one side.
	if matchsticks[0] > side {
		return false
	}
	sides := [4]int{}

	var dfs func(i int) bool
	dfs = func(i int) bool {
		if i == len(matchsticks) {
			return true // all placed, all buckets full
		}
		for b := 0; b < 4; b++ {
			// Symmetry prune: buckets with the same running length are
			// interchangeable — only try the first such bucket.
			dup := false
			for k := 0; k < b; k++ {
				if sides[k] == sides[b] {
					dup = true
					break
				}
			}
			if dup {
				continue
			}
			if sides[b]+matchsticks[i] <= side {
				sides[b] += matchsticks[i]
				if dfs(i + 1) {
					return true
				}
				sides[b] -= matchsticks[i]
			}
		}
		return false
	}
	return dfs(0)
}

// ── Approach 3: Bitmask DP (Optimal for n ≤ 15) ──────────────────────────────
//
// bitmaskDP solves Matchsticks to Square by iterating over subsets of sticks and
// tracking, for each subset, the leftover length modulo `side` on the current
// partially-built side.
//
// Intuition:
//
//	Encode which sticks are used as a bitmask (n ≤ 15 ⇒ ≤ 32768 masks). dp[mask]
//	= the used length on the CURRENT side, i.e. (sum of sticks in mask) mod side,
//	but only reachable if we can partition `mask` into completed sides plus this
//	partial one. Transition: from a reachable mask, add an unused stick j if it
//	fits in the remaining room of the current side (side - dp[mask]). Because
//	total = 4*side, filling all n sticks (mask = full) with dp == 0 means the
//	four sides all closed exactly.
//
// Algorithm:
//  1. Divisibility check; side = total/4.
//  2. dp array of size 2^n, dp[0] = 0 reachable, others = -1 (unreachable).
//  3. For each mask with dp[mask] >= 0, for each stick j not in mask:
//     if dp[mask]+sticks[j] <= side, set newUsed = (dp[mask]+sticks[j]) % side
//     and mark dp[mask | 1<<j] = newUsed (reachable).
//  4. Answer: dp[fullMask] == 0.
//
// Time:  O(2^n · n) — every mask examines n sticks.
// Space: O(2^n) — the dp table.
func bitmaskDP(matchsticks []int) bool {
	n := len(matchsticks)
	total := 0
	for _, m := range matchsticks {
		total += m
	}
	if total == 0 || total%4 != 0 {
		return false
	}
	side := total / 4
	// Any single stick longer than a side is an immediate no.
	for _, m := range matchsticks {
		if m > side {
			return false
		}
	}
	full := (1 << n) - 1 // mask with all sticks used
	// dp[mask] = length occupied on the current (in-progress) side, or -1 if the
	// subset `mask` cannot be arranged into complete-sides + this partial side.
	dp := make([]int, 1<<n)
	for i := range dp {
		dp[i] = -1
	}
	dp[0] = 0 // empty set: current side is empty, reachable
	for mask := 0; mask <= full; mask++ {
		if dp[mask] < 0 {
			continue // this configuration is not reachable
		}
		for j := 0; j < n; j++ {
			if mask&(1<<j) != 0 {
				continue // stick j already used
			}
			// Stick j must fit in the room left on the current side.
			if dp[mask]+matchsticks[j] <= side {
				// Adding it; wrap to 0 when the side closes exactly.
				used := (dp[mask] + matchsticks[j]) % side
				next := mask | (1 << j)
				// Prefer any reachable state; all lead to the same closure logic.
				if dp[next] == -1 {
					dp[next] = used
				}
			}
		}
	}
	// All sticks used AND the current side ended exactly closed ⇒ square.
	return dp[full] == 0
}

func main() {
	fmt.Println("=== Approach 1: Backtracking (Fill Four Sides) ===")
	fmt.Printf("[1 1 2 2 2]  => %v  (expected true)\n", backtracking([]int{1, 1, 2, 2, 2}))
	fmt.Printf("[3 3 3 3 4]  => %v  (expected false)\n", backtracking([]int{3, 3, 3, 3, 4}))
	fmt.Printf("[]           => %v  (expected false)\n", backtracking([]int{}))

	fmt.Println("=== Approach 2: Backtracking + Pruning ===")
	fmt.Printf("[1 1 2 2 2]  => %v  (expected true)\n", backtrackingPruned([]int{1, 1, 2, 2, 2}))
	fmt.Printf("[3 3 3 3 4]  => %v  (expected false)\n", backtrackingPruned([]int{3, 3, 3, 3, 4}))
	fmt.Printf("[5 5 5 5]    => %v  (expected true)\n", backtrackingPruned([]int{5, 5, 5, 5}))

	fmt.Println("=== Approach 3: Bitmask DP (Optimal) ===")
	fmt.Printf("[1 1 2 2 2]  => %v  (expected true)\n", bitmaskDP([]int{1, 1, 2, 2, 2}))
	fmt.Printf("[3 3 3 3 4]  => %v  (expected false)\n", bitmaskDP([]int{3, 3, 3, 3, 4}))
	fmt.Printf("[5 5 5 5]    => %v  (expected true)\n", bitmaskDP([]int{5, 5, 5, 5}))
}
