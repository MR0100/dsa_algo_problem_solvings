package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Nested Scan (Brute Force) ────────────────────────────────────
//
// bruteForce checks, for every element of nums1, whether it appears in nums2,
// using a seen-set to keep the result free of duplicates.
//
// Intuition:
//
//	The intersection is "values present in both arrays." The most literal check
//	is: for each value in nums1, linearly search nums2. Because the result must
//	contain each common value only once, track which values we have already
//	emitted in a set.
//
// Algorithm:
//  1. For each x in nums1: linear-search nums2 for x.
//  2. On a hit not already emitted, add x to the result and mark it emitted.
//
// Time:  O(n·m) — every element of nums1 may scan all of nums2.
// Space: O(min(n,m)) — the emitted-set and result.
func bruteForce(nums1, nums2 []int) []int {
	emitted := make(map[int]bool) // values already placed in the result
	res := []int{}
	for _, x := range nums1 {
		if emitted[x] { // already output this common value
			continue
		}
		for _, y := range nums2 { // linear search nums2 for x
			if x == y {
				res = append(res, x) // x is in both arrays
				emitted[x] = true    // never emit it again
				break
			}
		}
	}
	return res
}

// ── Approach 2: Hash Set (Optimal, unsorted) ─────────────────────────────────
//
// hashSet stores nums1 in a set, then keeps each element of nums2 that is in
// that set, removing it so duplicates in nums2 are not emitted twice.
//
// Intuition:
//
//	Membership testing is what a hash set is for. Dump nums1 into a set; then a
//	single pass over nums2 asks "is this value in nums1?" in O(1). Deleting the
//	value on the first hit guarantees the result is duplicate-free.
//
// Algorithm:
//  1. Build set from nums1.
//  2. For each y in nums2: if y ∈ set, append y and delete it from the set.
//
// Time:  O(n + m) — one pass to build, one to probe.
// Space: O(n) — the set of nums1's distinct values.
func hashSet(nums1, nums2 []int) []int {
	set := make(map[int]bool) // distinct values of nums1
	for _, x := range nums1 {
		set[x] = true
	}
	res := []int{}
	for _, y := range nums2 {
		if set[y] { // y appears in nums1 → it is in the intersection
			res = append(res, y)
			delete(set, y) // remove so a repeated y in nums2 is not re-added
		}
	}
	return res
}

// ── Approach 3: Sort + Two Pointers (Optimal, sorted output) ─────────────────
//
// twoPointers sorts both arrays and walks them together, emitting each shared
// value once.
//
// Intuition:
//
//	Once both arrays are sorted, a merge-style walk finds equal values: advance
//	the pointer at the smaller value; when they match, record it and skip past
//	all copies of that value in both arrays so it is emitted only once. Uses no
//	extra hash structure and yields sorted output.
//
// Algorithm:
//  1. Sort nums1 and nums2.
//  2. i=j=0. While both in range:
//     - nums1[i] < nums2[j] → i++.
//     - nums1[i] > nums2[j] → j++.
//     - equal → append value, then skip duplicates in both, advance both.
//
// Time:  O(n log n + m log m) — dominated by the two sorts.
// Space: O(1) extra (ignoring the sort and output) — in-place pointers.
func twoPointers(nums1, nums2 []int) []int {
	a := append([]int(nil), nums1...) // copy so we do not mutate the inputs
	b := append([]int(nil), nums2...)
	sort.Ints(a)
	sort.Ints(b)

	res := []int{}
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		switch {
		case a[i] < b[j]:
			i++ // a's value too small; advance a
		case a[i] > b[j]:
			j++ // b's value too small; advance b
		default: // a[i] == b[j]: a common value
			res = append(res, a[i])
			val := a[i]
			for i < len(a) && a[i] == val { // skip all copies in a
				i++
			}
			for j < len(b) && b[j] == val { // skip all copies in b
				j++
			}
		}
	}
	return res
}

// sortInts returns a sorted copy so unordered results print deterministically.
func sortInts(s []int) []int {
	out := append([]int(nil), s...)
	sort.Ints(out)
	return out
}

func main() {
	// Example 1: nums1=[1,2,2,1], nums2=[2,2] → [2]
	// Example 2: nums1=[4,9,5], nums2=[9,4,9,8,4] → [9,4] (any order; [4,9] also valid)
	// The result order is unspecified, so we sort before printing for a stable check.

	fmt.Println("=== Approach 1: Nested Scan (Brute Force) ===")
	fmt.Println(sortInts(bruteForce([]int{1, 2, 2, 1}, []int{2, 2})))       // expected [2]
	fmt.Println(sortInts(bruteForce([]int{4, 9, 5}, []int{9, 4, 9, 8, 4}))) // expected [4 9]

	fmt.Println("=== Approach 2: Hash Set (Optimal, unsorted) ===")
	fmt.Println(sortInts(hashSet([]int{1, 2, 2, 1}, []int{2, 2})))       // expected [2]
	fmt.Println(sortInts(hashSet([]int{4, 9, 5}, []int{9, 4, 9, 8, 4}))) // expected [4 9]

	fmt.Println("=== Approach 3: Sort + Two Pointers ===")
	fmt.Println(twoPointers([]int{1, 2, 2, 1}, []int{2, 2}))       // expected [2]
	fmt.Println(twoPointers([]int{4, 9, 5}, []int{9, 4, 9, 8, 4})) // expected [4 9]
}
