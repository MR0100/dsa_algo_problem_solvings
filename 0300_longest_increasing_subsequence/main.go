package main

import (
	"fmt"
	"sort"
)

// Longest Increasing Subsequence
//
// Given an integer array nums, return the length of the longest STRICTLY
// increasing subsequence (elements chosen in order but not necessarily
// contiguous, each strictly greater than the previous).

// ── Approach 1: Dynamic Programming O(n^2) ───────────────────────────────────
//
// dpQuadratic computes, for each index i, the length of the longest increasing
// subsequence that ENDS at i.
//
// Intuition:
//
//	dp[i] = 1 + max(dp[j]) over all j < i with nums[j] < nums[i]; if no such j,
//	dp[i] = 1 (the element alone). The answer is the maximum entry in dp.
//
// Algorithm:
//  1. dp[i] = 1 for all i.
//  2. For i from 1..n-1, for each j < i: if nums[j] < nums[i],
//     dp[i] = max(dp[i], dp[j]+1).
//  3. Return max(dp).
//
// Time:  O(n^2) — nested loops over all pairs.
// Space: O(n) — the dp array.
func dpQuadratic(nums []int) int {
	n := len(nums)
	if n == 0 {
		return 0
	}
	dp := make([]int, n)
	best := 1
	for i := 0; i < n; i++ {
		dp[i] = 1 // every element is a subsequence of length 1 by itself
		for j := 0; j < i; j++ {
			// extend any shorter run ending at an earlier, smaller value
			if nums[j] < nums[i] && dp[j]+1 > dp[i] {
				dp[i] = dp[j] + 1
			}
		}
		if dp[i] > best {
			best = dp[i] // track the global longest
		}
	}
	return best
}

// ── Approach 2: Patience Sorting + Binary Search O(n log n) (Optimal) ────────
//
// patienceBinarySearch maintains `tails`, where tails[k] is the smallest
// possible tail value of an increasing subsequence of length k+1.
//
// Intuition:
//
//	Greedily keeping each length's tail as small as possible leaves the most
//	room to extend later. For each number, binary-search the first tail >= it:
//	replace that tail (keeps it small) or append (grows the LIS). The length of
//	`tails` at the end is the LIS length. `tails` itself is NOT a valid LIS, but
//	its LENGTH is correct.
//
// Algorithm:
//  1. tails = [].
//  2. For each x: find the leftmost index in tails with value >= x.
//     - if none (x larger than all), append x (LIS grew).
//     - else overwrite tails[index] = x (smaller tail for that length).
//  3. Return len(tails).
//
// Time:  O(n log n) — binary search per element.
// Space: O(n) — the tails array.
func patienceBinarySearch(nums []int) int {
	tails := []int{} // tails[k] = smallest tail of an increasing subseq of length k+1
	for _, x := range nums {
		// first position whose tail is >= x (strictly increasing -> >=)
		i := sort.SearchInts(tails, x)
		if i == len(tails) {
			tails = append(tails, x) // x extends the longest run seen so far
		} else {
			tails[i] = x // replace to keep this length's tail minimal
		}
	}
	return len(tails)
}

func main() {
	// Example 1: nums = [10,9,2,5,3,7,101,18] -> 4  (e.g. [2,3,7,101])
	ex1 := []int{10, 9, 2, 5, 3, 7, 101, 18}
	// Example 2: nums = [0,1,0,3,2,3] -> 4
	ex2 := []int{0, 1, 0, 3, 2, 3}
	// Example 3: nums = [7,7,7,7,7,7,7] -> 1 (strictly increasing)
	ex3 := []int{7, 7, 7, 7, 7, 7, 7}

	fmt.Println("=== Approach 1: DP O(n^2) ===")
	fmt.Println(dpQuadratic(ex1)) // expected 4
	fmt.Println(dpQuadratic(ex2)) // expected 4
	fmt.Println(dpQuadratic(ex3)) // expected 1

	fmt.Println("=== Approach 2: Patience + Binary Search (Optimal) ===")
	fmt.Println(patienceBinarySearch(ex1)) // expected 4
	fmt.Println(patienceBinarySearch(ex2)) // expected 4
	fmt.Println(patienceBinarySearch(ex3)) // expected 1
}
