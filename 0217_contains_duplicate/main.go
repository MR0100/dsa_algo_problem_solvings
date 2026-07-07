package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Brute Force ──────────────────────────────────────────────────
//
// bruteForce solves Contains Duplicate by comparing every pair of elements.
//
// Intuition:
//
//	A duplicate exists iff some two DIFFERENT positions hold the same value.
//	The most literal way to check this is to compare each element against
//	every element after it. As soon as we find a match, we are done.
//
// Algorithm:
//  1. For each index i, for each index j > i:
//  2. If nums[i] == nums[j], a duplicate exists → return true.
//  3. If no pair matched, return false.
//
// Time:  O(n²) — every pair (i, j) is examined in the worst case.
// Space: O(1) — no auxiliary storage.
func bruteForce(nums []int) bool {
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ { // only pairs with j > i, no self-compare
			if nums[i] == nums[j] { // found two equal values at different indices
				return true
			}
		}
	}
	return false // no matching pair anywhere
}

// ── Approach 2: Sorting ──────────────────────────────────────────────────────
//
// sorting solves Contains Duplicate by sorting a copy and scanning for equal
// adjacent elements.
//
// Intuition:
//
//	If a value appears twice, after sorting the two copies land next to each
//	other. So duplicates in the whole array reduce to "any two ADJACENT
//	elements equal" in the sorted array — a single linear scan.
//
// Algorithm:
//  1. Copy nums (so we don't mutate the caller's slice) and sort the copy.
//  2. Scan i from 1..n-1: if arr[i] == arr[i-1], return true.
//  3. Otherwise return false.
//
// Time:  O(n log n) — dominated by the sort; the scan is O(n).
// Space: O(n) — the sorted copy (O(1) extra if sorting in place is allowed).
func sorting(nums []int) bool {
	arr := make([]int, len(nums)) // copy so the caller's slice is untouched
	copy(arr, nums)
	sort.Ints(arr) // equal values become adjacent after sorting
	for i := 1; i < len(arr); i++ {
		if arr[i] == arr[i-1] { // adjacent equal → duplicate
			return true
		}
	}
	return false
}

// ── Approach 3: Hash Set (Optimal) ───────────────────────────────────────────
//
// hashSet solves Contains Duplicate by remembering every value seen so far in
// a set and reporting the first collision.
//
// Intuition:
//
//	Walk the array once. Keep a set of values already encountered. The moment
//	the current value is already in the set, we have seen it before — that is
//	a duplicate. Otherwise add it and continue.
//
// Algorithm:
//  1. Create an empty set.
//  2. For each value v: if v is in the set, return true; else insert v.
//  3. If the loop finishes, return false.
//
// Time:  O(n) — one pass; each set lookup/insert is O(1) average.
// Space: O(n) — the set may hold up to n distinct values.
func hashSet(nums []int) bool {
	seen := make(map[int]struct{}, len(nums)) // struct{} is a zero-byte "present" marker
	for _, v := range nums {
		if _, ok := seen[v]; ok { // v already recorded → second occurrence
			return true
		}
		seen[v] = struct{}{} // record v as seen
	}
	return false
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(bruteForce([]int{1, 2, 3, 1}))                   // expected true
	fmt.Println(bruteForce([]int{1, 2, 3, 4}))                   // expected false
	fmt.Println(bruteForce([]int{1, 1, 1, 3, 3, 4, 3, 2, 4, 2})) // expected true

	fmt.Println("=== Approach 2: Sorting ===")
	fmt.Println(sorting([]int{1, 2, 3, 1}))                   // expected true
	fmt.Println(sorting([]int{1, 2, 3, 4}))                   // expected false
	fmt.Println(sorting([]int{1, 1, 1, 3, 3, 4, 3, 2, 4, 2})) // expected true

	fmt.Println("=== Approach 3: Hash Set (Optimal) ===")
	fmt.Println(hashSet([]int{1, 2, 3, 1}))                   // expected true
	fmt.Println(hashSet([]int{1, 2, 3, 4}))                   // expected false
	fmt.Println(hashSet([]int{1, 1, 1, 3, 3, 4, 3, 2, 4, 2})) // expected true
}
