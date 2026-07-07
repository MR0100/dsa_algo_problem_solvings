package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Nested Scan with Used Marks (Brute Force) ─────────────────────
//
// bruteForce pairs each element of nums1 with an as-yet-unused equal element of
// nums2, so each value is matched by multiplicity.
//
// Intuition:
//
//	Unlike #349, duplicates DO count here: if a value appears twice in both
//	arrays, it must appear twice in the answer. So each element of nums1 should
//	consume one matching, not-yet-consumed element of nums2. A boolean "used"
//	array on nums2 prevents matching the same slot twice.
//
// Algorithm:
//  1. used := bool slice over nums2.
//  2. For each x in nums1: scan nums2 for the first j with nums2[j]==x and
//     !used[j]; if found, append x and mark used[j].
//
// Time:  O(n·m) — each nums1 element may scan all of nums2.
// Space: O(m) — the used marks (plus the result).
func bruteForce(nums1, nums2 []int) []int {
	used := make([]bool, len(nums2)) // which nums2 slots are already matched
	res := []int{}
	for _, x := range nums1 {
		for j := 0; j < len(nums2); j++ {
			if !used[j] && nums2[j] == x { // an unclaimed equal element
				res = append(res, x) // pair them up
				used[j] = true       // this slot is now consumed
				break
			}
		}
	}
	return res
}

// ── Approach 2: Frequency Map (Optimal, unsorted) ────────────────────────────
//
// hashMap counts occurrences in the smaller array, then decrements while
// scanning the other, emitting a value each time its remaining count is > 0.
//
// Intuition:
//
//	Multiplicity of a shared value in the result is min(count_in_nums1,
//	count_in_nums2). Build a frequency map of one array; walk the other, and for
//	each value with a positive remaining count, emit it and decrement. When the
//	count hits zero, that value is exhausted — exactly the min behaviour.
//
// Algorithm:
//  1. Build count[v] over nums1.
//  2. For each y in nums2: if count[y] > 0, append y and count[y]--.
//
// Time:  O(n + m) — one pass to count, one to consume.
// Space: O(min(n,m)) — map over the array we choose to count (the smaller one).
func hashMap(nums1, nums2 []int) []int {
	// Count the smaller array to minimise map memory.
	if len(nums1) > len(nums2) {
		nums1, nums2 = nums2, nums1
	}
	count := make(map[int]int) // value → remaining available occurrences
	for _, x := range nums1 {
		count[x]++
	}
	res := []int{}
	for _, y := range nums2 {
		if count[y] > 0 { // still have an unmatched occurrence of y
			res = append(res, y)
			count[y]-- // consume it
		}
	}
	return res
}

// ── Approach 3: Sort + Two Pointers (Optimal, sorted) ────────────────────────
//
// twoPointers sorts both arrays and merges them, emitting one value for each
// matched pair.
//
// Intuition:
//
//	After sorting, equal values line up. A merge walk advances the smaller side;
//	on a match it emits the value ONCE and advances BOTH pointers, so a value
//	shared k times produces exactly k outputs (min of the two run lengths). This
//	is also the answer to the follow-ups: sorted inputs need no hash map, and if
//	nums2 is on disk you can stream it past a sorted nums1.
//
// Algorithm:
//  1. Sort both arrays.
//  2. i=j=0. While both in range:
//     - a[i] < b[j] → i++.
//     - a[i] > b[j] → j++.
//     - equal → append value, i++, j++ (consume one from each).
//
// Time:  O(n log n + m log m) — the sorts dominate.
// Space: O(1) extra beyond the sorted copies and output.
func twoPointers(nums1, nums2 []int) []int {
	a := append([]int(nil), nums1...) // copy to avoid mutating inputs
	b := append([]int(nil), nums2...)
	sort.Ints(a)
	sort.Ints(b)

	res := []int{}
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		switch {
		case a[i] < b[j]:
			i++ // advance the smaller side
		case a[i] > b[j]:
			j++
		default: // matched pair: emit once, consume from both
			res = append(res, a[i])
			i++
			j++
		}
	}
	return res
}

// sortInts returns a sorted copy for deterministic printing of unordered output.
func sortInts(s []int) []int {
	out := append([]int(nil), s...)
	sort.Ints(out)
	return out
}

func main() {
	// Example 1: nums1=[1,2,2,1], nums2=[2,2] → [2,2]
	// Example 2: nums1=[4,9,5], nums2=[9,4,9,8,4] → [4,9] (any order; [9,4] also valid)
	// Result order is unspecified, so we sort before printing for a stable check.

	fmt.Println("=== Approach 1: Nested Scan with Used Marks (Brute Force) ===")
	fmt.Println(sortInts(bruteForce([]int{1, 2, 2, 1}, []int{2, 2})))       // expected [2 2]
	fmt.Println(sortInts(bruteForce([]int{4, 9, 5}, []int{9, 4, 9, 8, 4}))) // expected [4 9]

	fmt.Println("=== Approach 2: Frequency Map (Optimal, unsorted) ===")
	fmt.Println(sortInts(hashMap([]int{1, 2, 2, 1}, []int{2, 2})))       // expected [2 2]
	fmt.Println(sortInts(hashMap([]int{4, 9, 5}, []int{9, 4, 9, 8, 4}))) // expected [4 9]

	fmt.Println("=== Approach 3: Sort + Two Pointers ===")
	fmt.Println(twoPointers([]int{1, 2, 2, 1}, []int{2, 2}))                 // expected [2 2]
	fmt.Println(sortInts(twoPointers([]int{4, 9, 5}, []int{9, 4, 9, 8, 4}))) // expected [4 9]
}
