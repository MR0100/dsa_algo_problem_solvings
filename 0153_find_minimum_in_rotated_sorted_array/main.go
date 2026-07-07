package main

import "fmt"

// ── Approach 1: Brute Force (Linear Scan) ────────────────────────────────────
//
// bruteForce solves Find Minimum in Rotated Sorted Array by scanning every
// element and keeping the smallest.
//
// Intuition:
//
//	The minimum of ANY array can be found by looking at every element once.
//	This ignores the rotated-sorted structure entirely, so it cannot meet the
//	required O(log n), but it is the correctness baseline.
//
// Algorithm:
//  1. Start with minVal = nums[0].
//  2. Compare every remaining element against minVal, keeping the smaller.
//  3. Return minVal.
//
// Time:  O(n) — every element inspected once.
// Space: O(1) — a single scalar.
func bruteForce(nums []int) int {
	minVal := nums[0] // candidate minimum
	for _, v := range nums[1:] {
		if v < minVal {
			minVal = v // found a smaller element
		}
	}
	return minVal
}

// ── Approach 2: Binary Search (Optimal) ──────────────────────────────────────
//
// binarySearch solves Find Minimum in Rotated Sorted Array in O(log n) by
// comparing the midpoint against the rightmost element.
//
// Intuition:
//
//	A rotated sorted array is two sorted runs; the minimum is the head of the
//	second run (the "pivot"). Compare nums[mid] with nums[hi]:
//	  - nums[mid] > nums[hi]: the drop (pivot) is strictly RIGHT of mid, so
//	    the minimum lives in (mid, hi].
//	  - nums[mid] < nums[hi]: mid..hi is sorted, so the minimum is at mid or
//	    LEFT of it — the pivot cannot be inside a sorted run.
//	Elements are unique, so equality never happens while lo < hi. Comparing
//	against nums[hi] (not nums[lo]) is what handles the "not rotated" case
//	(rotated n times = original array) without special-casing.
//
// Algorithm:
//  1. lo, hi = 0, n-1.
//  2. While lo < hi:
//     a. mid = lo + (hi-lo)/2.
//     b. If nums[mid] > nums[hi], the minimum is right of mid: lo = mid+1.
//     c. Else, the minimum is at mid or left of it: hi = mid.
//  3. When lo == hi, that index holds the minimum.
//
// Time:  O(log n) — the search range halves every iteration.
// Space: O(1) — two indices.
func binarySearch(nums []int) int {
	lo, hi := 0, len(nums)-1 // inclusive search window that contains the min
	for lo < hi {
		mid := lo + (hi-lo)/2 // overflow-safe midpoint (floor)
		if nums[mid] > nums[hi] {
			// mid sits in the FIRST (larger) run → pivot is strictly right
			lo = mid + 1
		} else {
			// nums[mid] < nums[hi] → mid..hi sorted → min at mid or left
			hi = mid // keep mid: it may itself be the minimum
		}
	}
	return nums[lo] // window collapsed onto the pivot
}

// ── Approach 3: Divide and Conquer (Recursive) ───────────────────────────────
//
// divideAndConquer solves Find Minimum in Rotated Sorted Array recursively:
// a sorted range answers immediately, otherwise split and recurse.
//
// Intuition:
//
//	If nums[lo] <= nums[hi], the range lo..hi is fully sorted (uniqueness
//	guarantees no flat plateaus) and its minimum is simply nums[lo].
//	Otherwise the range wraps around the pivot; split it in half — the pivot
//	lies in exactly one half, and the OTHER half is sorted and returns in
//	O(1). So only one branch ever goes deep, keeping the cost logarithmic.
//
// Algorithm:
//  1. rec(lo, hi):
//     a. If nums[lo] <= nums[hi], the range is sorted → return nums[lo].
//     b. Split at mid; return min(rec(lo, mid), rec(mid+1, hi)).
//  2. Answer is rec(0, n-1).
//
// Time:  O(log n) — the sorted half terminates instantly, so only the pivot half recurses further.
// Space: O(log n) — recursion stack depth.
func divideAndConquer(nums []int) int {
	var rec func(lo, hi int) int
	rec = func(lo, hi int) int {
		if nums[lo] <= nums[hi] {
			// sorted (or single-element) range → smallest is the first entry
			return nums[lo]
		}
		mid := lo + (hi-lo)/2 // split point
		// pivot is in exactly one half; the other half returns immediately
		return min(rec(lo, mid), rec(mid+1, hi))
	}
	return rec(0, len(nums)-1)
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Linear Scan) ===")
	fmt.Println(bruteForce([]int{3, 4, 5, 1, 2}))       // 1
	fmt.Println(bruteForce([]int{4, 5, 6, 7, 0, 1, 2})) // 0
	fmt.Println(bruteForce([]int{11, 13, 15, 17}))      // 11

	fmt.Println("=== Approach 2: Binary Search (Optimal) ===")
	fmt.Println(binarySearch([]int{3, 4, 5, 1, 2}))       // 1
	fmt.Println(binarySearch([]int{4, 5, 6, 7, 0, 1, 2})) // 0
	fmt.Println(binarySearch([]int{11, 13, 15, 17}))      // 11

	fmt.Println("=== Approach 3: Divide and Conquer ===")
	fmt.Println(divideAndConquer([]int{3, 4, 5, 1, 2}))       // 1
	fmt.Println(divideAndConquer([]int{4, 5, 6, 7, 0, 1, 2})) // 0
	fmt.Println(divideAndConquer([]int{11, 13, 15, 17}))      // 11
}
