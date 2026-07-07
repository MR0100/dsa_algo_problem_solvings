package main

import "fmt"

// ── Approach 1: Brute Force (Check Every Subarray) ───────────────────────────
//
// bruteForce solves Arithmetic Slices by, for each start index, extending the
// subarray one element at a time and checking whether it stays arithmetic.
//
// Intuition:
//
//	A slice [i..j] (length >= 3) is arithmetic iff every consecutive difference in
//	it equals the first difference nums[i+1]-nums[i]. Fix the start i and record
//	that reference difference; then walk j forward. As long as the newest
//	difference nums[j]-nums[j-1] still equals the reference, the subarray [i..j]
//	is arithmetic, so count it. The moment a difference breaks, no longer slice
//	starting at i can be arithmetic, so stop extending and move the start.
//
// Algorithm:
//  1. count = 0.
//  2. For each start i from 0..n-3:
//     diff = nums[i+1]-nums[i]; for j = i+2..n-1:
//     if nums[j]-nums[j-1] == diff → count++ (slice [i..j] is arithmetic)
//     else break (extension broke).
//  3. Return count.
//
// Time:  O(n^2) — each start can extend up to n elements.
// Space: O(1) — a couple of scalars.
func bruteForce(nums []int) int {
	n := len(nums)
	count := 0
	// A start needs at least two more elements to form a length-3 slice.
	for i := 0; i+2 < n; i++ {
		diff := nums[i+1] - nums[i] // the constant difference this run must keep
		for j := i + 2; j < n; j++ {
			if nums[j]-nums[j-1] == diff {
				count++ // [i..j] is arithmetic (length >= 3)
			} else {
				break // difference changed → longer slices from i are impossible
			}
		}
	}
	return count
}

// ── Approach 2: DP by Ending Index (Bottom-Up) ───────────────────────────────
//
// dpBottomUp solves Arithmetic Slices with a 1-D DP where dp[i] is the number of
// arithmetic slices that END exactly at index i.
//
// Intuition:
//
//	If the last three elements ...nums[i-2], nums[i-1], nums[i] share a common
//	difference, then every arithmetic slice ending at i-1 can be extended by one
//	to end at i, and additionally the brand-new length-3 slice [i-2..i] appears.
//	So dp[i] = dp[i-1] + 1. If the difference breaks at i, no arithmetic slice ends
//	at i, so dp[i] = 0. The total answer is the sum of all dp[i].
//
// Algorithm:
//  1. dp = make([]int, n); total = 0.
//  2. For i = 2..n-1: if nums[i]-nums[i-1] == nums[i-1]-nums[i-2],
//     dp[i] = dp[i-1] + 1; total += dp[i].
//  3. Return total.
//
// Time:  O(n) — a single pass.
// Space: O(n) — the dp array (reducible to O(1), see Approach 3).
func dpBottomUp(nums []int) int {
	n := len(nums)
	if n < 3 {
		return 0 // impossible to have a length-3 slice
	}
	dp := make([]int, n) // dp[i] = number of arithmetic slices ending at i
	total := 0
	for i := 2; i < n; i++ {
		// Same difference across the last three elements?
		if nums[i]-nums[i-1] == nums[i-1]-nums[i-2] {
			dp[i] = dp[i-1] + 1 // extend all slices ending at i-1, plus the new triple
			total += dp[i]      // accumulate into the running answer
		}
		// else dp[i] stays 0 (the zero value): no slice ends here
	}
	return total
}

// ── Approach 3: O(1) Space Counting (Optimal) ────────────────────────────────
//
// countingOptimal solves Arithmetic Slices by collapsing the DP array into a
// single running counter, and by summing consecutive-difference "runs".
//
// Intuition:
//
//	dp[i] only ever reads dp[i-1], so one variable `cur` suffices. `cur` counts how
//	many arithmetic slices end at the current index: bump it while the difference
//	holds, reset it to 0 when the difference breaks, and add it to the total each
//	step. Equivalently: a maximal run of k consecutive equal differences (i.e. k+1
//	numbers in arithmetic progression) contributes 1+2+...+(k-1) = (k-1)k/2 slices;
//	the running counter computes exactly that triangular sum incrementally.
//
// Algorithm:
//  1. total = 0, cur = 0.
//  2. For i = 2..n-1: if difference holds, cur++ and total += cur; else cur = 0.
//  3. Return total.
//
// Time:  O(n) — one pass.
// Space: O(1) — two integer counters.
func countingOptimal(nums []int) int {
	total, cur := 0, 0 // cur = arithmetic slices ending at the current index
	for i := 2; i < len(nums); i++ {
		if nums[i]-nums[i-1] == nums[i-1]-nums[i-2] {
			cur++        // one more slice ends here than ended at i-1 (+ the new triple)
			total += cur // triangular accumulation: 1,2,3,... across a run
		} else {
			cur = 0 // difference broke → no slice ends at i; restart the run
		}
	}
	return total
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Check Every Subarray) ===")
	fmt.Println(bruteForce([]int{1, 2, 3, 4}))        // expected 3
	fmt.Println(bruteForce([]int{1}))                 // expected 0
	fmt.Println(bruteForce([]int{1, 3, 5, 7, 9}))     // expected 6
	fmt.Println(bruteForce([]int{7, 7, 7, 7}))        // expected 3
	fmt.Println(bruteForce([]int{1, 2, 3, 8, 9, 10})) // expected 2

	fmt.Println("=== Approach 2: DP by Ending Index (Bottom-Up) ===")
	fmt.Println(dpBottomUp([]int{1, 2, 3, 4}))        // expected 3
	fmt.Println(dpBottomUp([]int{1}))                 // expected 0
	fmt.Println(dpBottomUp([]int{1, 3, 5, 7, 9}))     // expected 6
	fmt.Println(dpBottomUp([]int{7, 7, 7, 7}))        // expected 3
	fmt.Println(dpBottomUp([]int{1, 2, 3, 8, 9, 10})) // expected 2

	fmt.Println("=== Approach 3: O(1) Space Counting (Optimal) ===")
	fmt.Println(countingOptimal([]int{1, 2, 3, 4}))        // expected 3
	fmt.Println(countingOptimal([]int{1}))                 // expected 0
	fmt.Println(countingOptimal([]int{1, 3, 5, 7, 9}))     // expected 6
	fmt.Println(countingOptimal([]int{7, 7, 7, 7}))        // expected 3
	fmt.Println(countingOptimal([]int{1, 2, 3, 8, 9, 10})) // expected 2
}
