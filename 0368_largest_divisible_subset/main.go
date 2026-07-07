package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: DP with Parent Pointers (Optimal) ────────────────────────────
//
// dpWithParents solves Largest Divisible Subset by sorting the numbers and
// running a longest-chain DP where each element can extend any earlier element
// it is divisible by.
//
// Intuition:
//
//	Divisibility on a sorted array behaves like ≤: if a | b and b | c then
//	a | c. So after sorting ascending, ANY subset in which each element divides
//	the next larger chosen element is fully pairwise-divisible — we only need
//	each consecutive pair to divide. That turns the problem into "longest chain
//	where nums[j] divides nums[i] for j < i", a Longest-Increasing-Subsequence
//	shaped DP. dp[i] = length of the best chain ending at i; parent[i] remembers
//	which earlier index we extended, so we can reconstruct the actual subset.
//
// Algorithm:
//  1. Sort nums ascending.
//  2. dp[i] = 1, parent[i] = -1 for all i.
//  3. For i from 0..n−1, for j from 0..i−1:
//     if nums[i] % nums[j] == 0 and dp[j]+1 > dp[i], set dp[i]=dp[j]+1,
//     parent[i]=j.
//  4. Track the index with the maximum dp value.
//  5. Walk parent pointers back from that index to rebuild the subset, then
//     reverse it.
//
// Time:  O(n²) — the double loop dominates the O(n log n) sort.
// Space: O(n) — dp and parent arrays.
func dpWithParents(nums []int) []int {
	n := len(nums)
	if n == 0 {
		return []int{}
	}
	sort.Ints(nums) // divisibility acts like ≤ once sorted ascending

	dp := make([]int, n)     // dp[i] = longest divisible chain ending at i
	parent := make([]int, n) // parent[i] = previous index in that chain
	best := 0                // index where the overall-longest chain ends

	for i := 0; i < n; i++ {
		dp[i] = 1      // a single element is always a valid chain
		parent[i] = -1 // no predecessor yet
		for j := 0; j < i; j++ {
			// nums[j] < nums[i] (sorted); if it divides nums[i] we may extend
			// the chain that ended at j.
			if nums[i]%nums[j] == 0 && dp[j]+1 > dp[i] {
				dp[i] = dp[j] + 1
				parent[i] = j
			}
		}
		if dp[i] > dp[best] { // remember the global best end-point
			best = i
		}
	}

	// Reconstruct by following parent pointers from best back to the start.
	var result []int
	for i := best; i != -1; i = parent[i] {
		result = append(result, nums[i])
	}
	// We collected largest→smallest; reverse to ascending for a tidy answer.
	for l, r := 0, len(result)-1; l < r; l, r = l+1, r-1 {
		result[l], result[r] = result[r], result[l]
	}
	return result
}

// ── Approach 2: DP Storing Full Subsets (Brute Force–ish) ────────────────────
//
// dpFullSubsets solves Largest Divisible Subset with the same DP recurrence but
// stores the ENTIRE best subset at each index instead of parent pointers.
//
// Intuition:
//
//	Conceptually identical to Approach 1, but easier to read at the cost of
//	memory: subsets[i] holds the actual largest divisible subset that ends at
//	nums[i]. To extend, copy the best predecessor's subset and append nums[i].
//	Simpler to reason about; O(n²) extra space from copying slices.
//
// Algorithm:
//  1. Sort ascending.
//  2. subsets[i] starts as [nums[i]].
//  3. For each i, for each j < i with nums[i] % nums[j] == 0: if
//     len(subsets[j])+1 > len(subsets[i]), replace subsets[i] with a copy of
//     subsets[j] plus nums[i].
//  4. Return the longest subsets[i].
//
// Time:  O(n²) comparisons, plus O(n) per copy ⇒ O(n²) to O(n³) worst copying.
// Space: O(n²) — each subsets[i] may hold up to n elements.
func dpFullSubsets(nums []int) []int {
	n := len(nums)
	if n == 0 {
		return []int{}
	}
	sort.Ints(nums)

	subsets := make([][]int, n) // subsets[i] = best divisible subset ending at i
	best := 0
	for i := 0; i < n; i++ {
		subsets[i] = []int{nums[i]} // at minimum, the element by itself
		for j := 0; j < i; j++ {
			if nums[i]%nums[j] == 0 && len(subsets[j])+1 > len(subsets[i]) {
				// Copy predecessor's subset so we don't alias/mutate it.
				cp := make([]int, len(subsets[j]))
				copy(cp, subsets[j])
				subsets[i] = append(cp, nums[i])
			}
		}
		if len(subsets[i]) > len(subsets[best]) {
			best = i
		}
	}
	return subsets[best]
}

func main() {
	// Example 1: nums = [1,2,3]   → [1,2] (or [1,3])
	// Example 2: nums = [1,2,4,8] → [1,2,4,8]

	fmt.Println("=== Approach 1: DP with Parent Pointers (Optimal) ===")
	fmt.Println(dpWithParents([]int{1, 2, 3}))    // expected [1 2]
	fmt.Println(dpWithParents([]int{1, 2, 4, 8})) // expected [1 2 4 8]

	fmt.Println("=== Approach 2: DP Storing Full Subsets ===")
	fmt.Println(dpFullSubsets([]int{1, 2, 3}))    // expected [1 2]
	fmt.Println(dpFullSubsets([]int{1, 2, 4, 8})) // expected [1 2 4 8]
}
