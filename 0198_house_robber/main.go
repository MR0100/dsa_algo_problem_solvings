package main

import "fmt"

// ── Approach 1: Brute Force ──────────────────────────────────────────────────
//
// bruteForce solves House Robber by trying, at every house, both choices —
// rob it or skip it — with plain recursion and no caching.
//
// Intuition:
//   Standing in front of house i there are exactly two legal futures: rob
//   house i and jump to house i+2 (the alarm forbids i+1), or skip house i
//   and move to house i+1. The best loot from house i onward is the better
//   of those two futures. Recursing on both branches enumerates every valid
//   subset of non-adjacent houses.
//
// Algorithm:
//   1. Define robFrom(i) = best loot obtainable from houses i..n-1.
//   2. Base case: i >= n → 0 (no houses left).
//   3. Recurrence: robFrom(i) = max(nums[i] + robFrom(i+2), robFrom(i+1)).
//   4. The answer is robFrom(0).
//
// Time:  O(2ⁿ) — each call branches twice and identical subproblems are
//         recomputed exponentially often (Fibonacci-shaped recursion tree).
// Space: O(n) — maximum recursion stack depth.
func bruteForce(nums []int) int {
	var robFrom func(i int) int
	robFrom = func(i int) int {
		if i >= len(nums) {
			return 0 // ran past the last house: nothing more to steal
		}
		take := nums[i] + robFrom(i+2) // rob house i → house i+1 is off-limits
		skip := robFrom(i + 1)         // leave house i → free to consider i+1
		return max(take, skip)
	}
	return robFrom(0)
}

// ── Approach 2: DP Top-Down (Memoization) ────────────────────────────────────
//
// dpTopDown adds a memo table to the brute-force recursion so every
// subproblem robFrom(i) is solved exactly once.
//
// Intuition:
//   The brute force recomputes robFrom(i) for the same i over and over, yet
//   there are only n distinct subproblems. Caching each answer the first time
//   it is computed collapses the exponential recursion tree into n calls.
//
// Algorithm:
//   1. memo[i] caches the answer for houses i..n-1; -1 marks "not computed"
//      (safe sentinel because loot is never negative).
//   2. On each call return the cached value when present.
//   3. Otherwise compute max(nums[i] + robFrom(i+2), robFrom(i+1)), store it
//      in memo[i], and return it.
//
// Time:  O(n) — n distinct subproblems, each computed once in O(1).
// Space: O(n) — the memo table plus the recursion stack.
func dpTopDown(nums []int) int {
	memo := make([]int, len(nums))
	for i := range memo {
		memo[i] = -1 // -1 = not yet computed (valid answers are always >= 0)
	}
	var robFrom func(i int) int
	robFrom = func(i int) int {
		if i >= len(nums) {
			return 0
		}
		if memo[i] != -1 {
			return memo[i] // already solved: reuse instead of recomputing
		}
		memo[i] = max(nums[i]+robFrom(i+2), robFrom(i+1))
		return memo[i]
	}
	return robFrom(0)
}

// ── Approach 3: DP Bottom-Up (Tabulation) ────────────────────────────────────
//
// dpBottomUp fills a table dp[i] = best loot from the first i houses,
// iterating from the smallest prefix upward.
//
// Intuition:
//   Flip the recursion around: instead of "best from house i to the end",
//   compute "best among the first i houses". Each new house offers the same
//   two choices — rob it (its cash plus the best of the first i-2 houses) or
//   skip it (carry over the best of the first i-1 houses).
//
// Algorithm:
//   1. dp[0] = 0 (no houses), dp[1] = nums[0] (one house: rob it).
//   2. For i = 2..n: dp[i] = max(dp[i-1], dp[i-2] + nums[i-1]).
//   3. The answer is dp[n].
//
// Time:  O(n) — one pass over the houses.
// Space: O(n) — the dp table of n+1 entries.
func dpBottomUp(nums []int) int {
	n := len(nums)
	dp := make([]int, n+1) // dp[i] = max loot using only the first i houses
	dp[0] = 0              // zero houses → zero loot
	dp[1] = nums[0]        // one house → rob it (constraints guarantee n >= 1)
	for i := 2; i <= n; i++ {
		skip := dp[i-1]             // leave house i-1 (0-indexed) alone
		take := dp[i-2] + nums[i-1] // rob it: its neighbour must be skipped
		dp[i] = max(skip, take)
	}
	return dp[n]
}

// ── Approach 4: Space-Optimized DP (Optimal) ─────────────────────────────────
//
// spaceOptimized keeps only the last two dp values in two rolling variables,
// since the recurrence never looks further back than dp[i-2].
//
// Intuition:
//   dp[i] = max(dp[i-1], dp[i-2] + nums[i-1]) touches only the two previous
//   entries, so the whole table collapses into a rolling pair (prev2, prev1)
//   — the same trick that reduces Fibonacci to constant space.
//
// Algorithm:
//   1. Start prev2 = 0 and prev1 = 0 (best loot two/one houses back).
//   2. For each house value v: curr = max(prev1, prev2 + v), then slide the
//      window: prev2 = prev1, prev1 = curr.
//   3. prev1 holds the answer after the loop.
//
// Time:  O(n) — one pass over the houses.
// Space: O(1) — two integers regardless of input size.
func spaceOptimized(nums []int) int {
	prev2, prev1 := 0, 0 // dp[i-2] and dp[i-1] of the tabulation
	for _, v := range nums {
		curr := max(prev1, prev2+v) // skip this house (prev1) or rob it (prev2+v)
		prev2, prev1 = prev1, curr  // slide the two-value window forward
	}
	return prev1
}

func main() {
	nums1 := []int{1, 2, 3, 1}    // Example 1
	nums2 := []int{2, 7, 9, 3, 1} // Example 2

	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(bruteForce(nums1)) // 4
	fmt.Println(bruteForce(nums2)) // 12

	fmt.Println("=== Approach 2: DP Top-Down (Memoization) ===")
	fmt.Println(dpTopDown(nums1)) // 4
	fmt.Println(dpTopDown(nums2)) // 12

	fmt.Println("=== Approach 3: DP Bottom-Up (Tabulation) ===")
	fmt.Println(dpBottomUp(nums1)) // 4
	fmt.Println(dpBottomUp(nums2)) // 12

	fmt.Println("=== Approach 4: Space-Optimized DP (Optimal) ===")
	fmt.Println(spaceOptimized(nums1)) // 4
	fmt.Println(spaceOptimized(nums2)) // 12
}
