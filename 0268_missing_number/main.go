package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Sorting ──────────────────────────────────────────────────────
//
// sortScan solves Missing Number by sorting and finding the first index whose
// value does not equal the index.
//
// Intuition:
//
//	If no number were missing, the sorted array would be [0,1,2,...,n], with
//	nums[i] == i everywhere. The first place that breaks is the missing value.
//	If every index matches, the missing number is n itself.
//
// Algorithm:
//  1. Sort nums.
//  2. Scan; return the first i where nums[i] != i.
//  3. If none, return n.
//
// Time:  O(n log n) — dominated by the sort.
// Space: O(1) — in-place sort (ignoring sort's internal use).
func sortScan(nums []int) int {
	sort.Ints(nums) // arrange 0..n (with one gap) in order
	for i, v := range nums {
		if v != i { // first index whose value slipped
			return i
		}
	}
	return len(nums) // gap is at the very end -> n is missing
}

// ── Approach 2: Hash Set ─────────────────────────────────────────────────────
//
// hashSet solves Missing Number by putting every value in a set and probing
// 0..n for the absentee.
//
// Intuition:
//
//	Store all present numbers, then ask which of 0..n is not in the set.
//
// Algorithm:
//  1. Insert every value into a set.
//  2. For i in 0..n: if i is not in the set, return i.
//
// Time:  O(n) — build set + probe.
// Space: O(n) — the set.
func hashSet(nums []int) int {
	present := make(map[int]struct{}, len(nums)) // set of present values
	for _, v := range nums {
		present[v] = struct{}{}
	}
	for i := 0; i <= len(nums); i++ { // 0..n inclusive
		if _, ok := present[i]; !ok { // i never appeared
			return i
		}
	}
	return -1 // unreachable for valid input
}

// ── Approach 3: Gauss Sum ────────────────────────────────────────────────────
//
// gaussSum solves Missing Number by subtracting the actual sum from the
// expected sum of 0..n.
//
// Intuition:
//
//	The complete set 0..n sums to n(n+1)/2. The array is that set minus one
//	number, so expected - actual = the missing number.
//
// Algorithm:
//  1. expected = n(n+1)/2 where n = len(nums).
//  2. actual = sum of nums.
//  3. Return expected - actual.
//
// Time:  O(n) — one pass to sum.
// Space: O(1).
func gaussSum(nums []int) int {
	n := len(nums)
	expected := n * (n + 1) / 2 // sum of 0..n
	actual := 0
	for _, v := range nums {
		actual += v // sum of present values
	}
	return expected - actual // the gap
}

// ── Approach 4: XOR (Optimal) ────────────────────────────────────────────────
//
// xorBits solves Missing Number by XOR-ing all indices, n, and all values.
//
// Intuition:
//
//	XOR-ing a number with itself yields 0. If we XOR together every index
//	0..n and every value in nums, every present number cancels with its own
//	index, leaving only the missing number. This avoids the overflow risk of
//	summation.
//
// Algorithm:
//  1. result = n (accounts for the top index that has no matching element).
//  2. For each i: result ^= i ^ nums[i].
//  3. Return result.
//
// Time:  O(n) — single pass.
// Space: O(1).
func xorBits(nums []int) int {
	result := len(nums) // seed with n (index with no element)
	for i, v := range nums {
		result ^= i ^ v // cancel index against value
	}
	return result // only the missing number survives
}

func main() {
	fmt.Println("=== Approach 1: Sorting ===")
	fmt.Println(sortScan([]int{3, 0, 1}))                   // expected 2
	fmt.Println(sortScan([]int{0, 1}))                      // expected 2
	fmt.Println(sortScan([]int{9, 6, 4, 2, 3, 5, 7, 0, 1})) // expected 8

	fmt.Println("=== Approach 2: Hash Set ===")
	fmt.Println(hashSet([]int{3, 0, 1}))                   // expected 2
	fmt.Println(hashSet([]int{0, 1}))                      // expected 2
	fmt.Println(hashSet([]int{9, 6, 4, 2, 3, 5, 7, 0, 1})) // expected 8

	fmt.Println("=== Approach 3: Gauss Sum ===")
	fmt.Println(gaussSum([]int{3, 0, 1}))                   // expected 2
	fmt.Println(gaussSum([]int{0, 1}))                      // expected 2
	fmt.Println(gaussSum([]int{9, 6, 4, 2, 3, 5, 7, 0, 1})) // expected 8

	fmt.Println("=== Approach 4: XOR (Optimal) ===")
	fmt.Println(xorBits([]int{3, 0, 1}))                   // expected 2
	fmt.Println(xorBits([]int{0, 1}))                      // expected 2
	fmt.Println(xorBits([]int{9, 6, 4, 2, 3, 5, 7, 0, 1})) // expected 8
}
