package main

import "fmt"

// ── Approach 1: Brute Force Backtracking ─────────────────────────────────────
//
// bruteForce solves Combination Sum IV by enumerating every ordered sequence of
// numbers (with repetition) whose values sum to target and counting them.
//
// Intuition:
//
//	Despite the name, this problem counts *ordered* combinations (permutations):
//	[1,2] and [2,1] are distinct. So the search tree, at every partial sum, may
//	pick ANY number in nums as the next term. When the remaining amount hits 0 we
//	found one valid ordering; when it goes negative we prune.
//
// Algorithm:
//  1. count(remaining): if remaining == 0 return 1.
//  2. For each num in nums with num <= remaining, add count(remaining-num).
//  3. Return the total.
//
// Time:  O(target^target) worst case — exponential branching, no memo.
// Space: O(target) — recursion depth.
func bruteForce(nums []int, target int) int {
	// count returns the number of ordered sequences summing to `remaining`.
	var count func(remaining int) int
	count = func(remaining int) int {
		if remaining == 0 {
			return 1 // an empty tail completes a valid ordering
		}
		total := 0
		for _, num := range nums {
			if num <= remaining {
				total += count(remaining - num) // pick num, recurse on the rest
			}
		}
		return total
	}
	return count(target)
}

// ── Approach 2: DP Top-Down (Memoized Recursion) ─────────────────────────────
//
// dpTopDown solves Combination Sum IV by caching count(remaining) so each
// subtotal is computed once.
//
// Intuition:
//
//	The brute-force recursion revisits the same `remaining` value over and over.
//	Since the number of orderings summing to a given remaining is fixed, memoize
//	it in an array indexed by remaining.
//
// Algorithm:
//  1. memo[remaining] caches the answer for that subtotal (-1 = unknown).
//  2. count(0) = 1; otherwise sum count(remaining-num) over valid nums, store, return.
//
// Time:  O(target × n) — target distinct states, each scans n numbers.
// Space: O(target) — memo array + recursion depth.
func dpTopDown(nums []int, target int) int {
	memo := make([]int, target+1) // memo[r] = orderings summing to r
	for i := range memo {
		memo[i] = -1 // -1 marks "not yet computed"
	}
	var count func(remaining int) int
	count = func(remaining int) int {
		if remaining == 0 {
			return 1 // base case: exactly hit the target
		}
		if memo[remaining] != -1 {
			return memo[remaining] // reuse the cached subtotal
		}
		total := 0
		for _, num := range nums {
			if num <= remaining {
				total += count(remaining - num)
			}
		}
		memo[remaining] = total // cache before returning
		return total
	}
	return count(target)
}

// ── Approach 3: DP Bottom-Up (Optimal) ───────────────────────────────────────
//
// dpBottomUp solves Combination Sum IV by filling dp[t] = number of ordered
// sequences summing to t, for t from 1 up to target.
//
// Intuition:
//
//	dp[t] counts orderings that reach exactly t. Any such ordering has a LAST
//	term num; removing it leaves an ordering summing to t-num. Summing over all
//	choices of the last term:  dp[t] = Σ dp[t-num] for num ≤ t.
//	Because the last term is chosen freely, order is respected — this is why the
//	amount loop is OUTER and the coin loop is INNER (the opposite of unordered
//	coin-change counting).
//
// Algorithm:
//  1. dp[0] = 1 (one empty ordering).
//  2. For t = 1..target, for each num ≤ t: dp[t] += dp[t-num].
//  3. Return dp[target].
//
// Time:  O(target × n).
// Space: O(target).
func dpBottomUp(nums []int, target int) int {
	dp := make([]int, target+1) // dp[t] = ordered sequences summing to t
	dp[0] = 1                   // the empty sequence sums to 0, one way
	for t := 1; t <= target; t++ {
		for _, num := range nums {
			if num <= t {
				dp[t] += dp[t-num] // append num as the final term
			}
		}
	}
	return dp[target]
}

func main() {
	fmt.Println("=== Approach 1: Brute Force Backtracking ===")
	fmt.Println(bruteForce([]int{1, 2, 3}, 4)) // expected 7
	fmt.Println(bruteForce([]int{9}, 3))       // expected 0

	fmt.Println("=== Approach 2: DP Top-Down (Memoized) ===")
	fmt.Println(dpTopDown([]int{1, 2, 3}, 4)) // expected 7
	fmt.Println(dpTopDown([]int{9}, 3))       // expected 0

	fmt.Println("=== Approach 3: DP Bottom-Up (Optimal) ===")
	fmt.Println(dpBottomUp([]int{1, 2, 3}, 4)) // expected 7
	fmt.Println(dpBottomUp([]int{9}, 3))       // expected 0
}
