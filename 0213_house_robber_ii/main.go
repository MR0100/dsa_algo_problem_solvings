package main

import "fmt"

// ── Approach 1: Two DP Passes over a 1D array (Optimal) ───────────────────────
//
// dpArray solves House Robber II by reducing the circular street to two linear
// House Robber I problems and taking the better result.
//
// Intuition:
//
//	Houses are in a CIRCLE, so the first and last houses are adjacent and can
//	never both be robbed. Split into two independent linear cases:
//	  (a) rob among houses [0 .. n-2]  (exclude the last house), and
//	  (b) rob among houses [1 .. n-1]  (exclude the first house).
//	Each case forbids the pair that touches, so neither can pick both ends of
//	the circle. Each subproblem is plain House Robber I (max sum of a linear
//	array with no two adjacent picks). The answer is max of the two. Handle
//	n == 1 specially (single house, no wrap-around).
//
// Algorithm:
//
//  1. If n == 1, return nums[0].
//  2. robLinear(lo, hi) computes the best non-adjacent sum over nums[lo..hi]
//     with a rolling DP: prev2, prev1 → cur = max(prev1, prev2 + nums[i]).
//  3. Return max(robLinear(0, n-2), robLinear(1, n-1)).
//
// Time:  O(n) — two linear scans.
// Space: O(1) — two rolling variables, no DP table.
func dpArray(nums []int) int {
	n := len(nums)
	if n == 1 {
		return nums[0] // single house: no neighbour, just take it
	}
	// robLinear = best non-adjacent sum over nums[lo..hi] inclusive.
	robLinear := func(lo, hi int) int {
		prev2, prev1 := 0, 0 // best up to i-2 and i-1 respectively
		for i := lo; i <= hi; i++ {
			cur := max(prev1, prev2+nums[i]) // skip house i, or rob it + best two back
			prev2, prev1 = prev1, cur        // slide the window forward
		}
		return prev1 // best over the whole range
	}
	// Case A excludes the last house; Case B excludes the first house.
	return max(robLinear(0, n-2), robLinear(1, n-1))
}

// ── Approach 2: Explicit DP table (two passes) ────────────────────────────────
//
// dpTable solves House Robber II with the same split but a full dp[] table per
// case, making the recurrence explicit and easy to trace.
//
// Intuition:
//
//	Identical circular-split idea, but instead of two rolling variables we fill
//	a dp array where dp[i] = best loot considering the sub-range up to local
//	index i. Clearer for a dry run at the cost of O(n) extra space.
//
// Algorithm:
//
//  1. If n == 1 return nums[0].
//  2. robRange(lo, hi): let m = hi-lo+1 houses.
//     dp[0] = nums[lo]; dp[1] = max(nums[lo], nums[lo+1]);
//     dp[k] = max(dp[k-1], dp[k-2] + nums[lo+k]).
//  3. Return max(robRange(0, n-2), robRange(1, n-1)).
//
// Time:  O(n).
// Space: O(n) for the dp table.
func dpTable(nums []int) int {
	n := len(nums)
	if n == 1 {
		return nums[0]
	}
	robRange := func(lo, hi int) int {
		m := hi - lo + 1
		if m == 1 {
			return nums[lo] // range of one house
		}
		dp := make([]int, m)
		dp[0] = nums[lo]                  // rob the first house of the range
		dp[1] = max(nums[lo], nums[lo+1]) // best of the first two
		for k := 2; k < m; k++ {
			// either skip house k (dp[k-1]) or rob it plus best two back (dp[k-2])
			dp[k] = max(dp[k-1], dp[k-2]+nums[lo+k])
		}
		return dp[m-1]
	}
	return max(robRange(0, n-2), robRange(1, n-1))
}

func main() {
	fmt.Println("=== Approach 1: Two DP Passes, O(1) space (Optimal) ===")
	fmt.Println(dpArray([]int{2, 3, 2}))    // 3
	fmt.Println(dpArray([]int{1, 2, 3, 1})) // 4
	fmt.Println(dpArray([]int{1, 2, 3}))    // 3

	fmt.Println("=== Approach 2: Explicit DP Table ===")
	fmt.Println(dpTable([]int{2, 3, 2}))    // 3
	fmt.Println(dpTable([]int{1, 2, 3, 1})) // 4
	fmt.Println(dpTable([]int{1, 2, 3}))    // 3
}
