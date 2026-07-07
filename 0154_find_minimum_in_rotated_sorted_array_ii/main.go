package main

import "fmt"

// ── Approach 1: Brute Force (Linear Scan) ────────────────────────────────────
//
// bruteForce solves Find Minimum in Rotated Sorted Array II by scanning every
// element and keeping the smallest.
//
// Intuition:
//
//	Duplicates or not, the minimum of any array is found by one full pass.
//	This is also, importantly, the true WORST-CASE cost of any comparison
//	algorithm on this input (see README): with duplicates an adversary can
//	hide the minimum anywhere, e.g. [2,2,2,...,2,0,2,...,2].
//
// Algorithm:
//  1. Start with minVal = nums[0].
//  2. Compare every remaining element, keeping the smaller.
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

// ── Approach 2: Binary Search with Duplicate Shrinking (Optimal) ─────────────
//
// binarySearch solves Find Minimum in Rotated Sorted Array II by comparing
// nums[mid] with nums[hi] and shrinking hi by one when they are equal.
//
// Intuition:
//
//	Same skeleton as LC 153: the minimum is the pivot where values drop.
//	  - nums[mid] > nums[hi]: pivot strictly right of mid → lo = mid+1.
//	  - nums[mid] < nums[hi]: mid..hi sorted → pivot at mid or left → hi = mid.
//	  - nums[mid] == nums[hi]: AMBIGUOUS. The pivot could be on either side
//	    (e.g. [2,2,2,0,2] vs [2,0,2,2,2] both give mid==hi==2). But nums[hi]
//	    has a twin at mid, so discarding JUST nums[hi] can never lose the
//	    minimum: even if nums[hi] IS the minimum, an equal value survives at
//	    mid. Shrink hi by one and retry.
//	Each equality step removes one element, so the worst case (all equal)
//	degrades to O(n) — provably unavoidable with duplicates.
//
// Algorithm:
//  1. lo, hi = 0, n-1.
//  2. While lo < hi:
//     a. mid = lo + (hi-lo)/2.
//     b. If nums[mid] > nums[hi]: lo = mid + 1.
//     c. Else if nums[mid] < nums[hi]: hi = mid.
//     d. Else (equal): hi--.
//  3. Return nums[lo].
//
// Time:  O(log n) average / O(n) worst case (all duplicates).
// Space: O(1) — two indices.
func binarySearch(nums []int) int {
	lo, hi := 0, len(nums)-1 // inclusive window guaranteed to hold the min
	for lo < hi {
		mid := lo + (hi-lo)/2 // overflow-safe midpoint
		switch {
		case nums[mid] > nums[hi]:
			// mid is in the first (larger) run → pivot strictly right of mid
			lo = mid + 1
		case nums[mid] < nums[hi]:
			// mid..hi is sorted → pivot at mid or to its left; keep mid
			hi = mid
		default:
			// nums[mid] == nums[hi]: cannot tell which side holds the pivot.
			// Safe to drop nums[hi] — its value also exists at mid, so the
			// minimum value is still inside the window.
			hi--
		}
	}
	return nums[lo] // window collapsed onto (one copy of) the minimum
}

// ── Approach 3: Divide and Conquer (Recursive) ───────────────────────────────
//
// divideAndConquer solves Find Minimum in Rotated Sorted Array II by
// recursing into both halves whenever sortedness cannot be proven.
//
// Intuition:
//
//	If nums[lo] < nums[hi], the range is strictly sorted → min is nums[lo].
//	With duplicates, nums[lo] == nums[hi] proves nothing (the pivot may hide
//	between equal walls, e.g. [2,2,0,2,2]), so both halves must be explored.
//	Sorted sub-ranges still cut off entire branches early, giving O(log n)
//	behaviour on inputs with few duplicates and O(n) only in the worst case.
//
// Algorithm:
//  1. rec(lo, hi):
//     a. If lo == hi → single element, return it.
//     b. If nums[lo] < nums[hi] → strictly sorted range, return nums[lo].
//     c. Split at mid; return min(rec(lo, mid), rec(mid+1, hi)).
//  2. Answer is rec(0, n-1).
//
// Time:  O(log n) average / O(n) worst case (both branches recurse on every level when the array is one big plateau).
// Space: O(log n) — ranges still halve, so the recursion depth stays logarithmic.
func divideAndConquer(nums []int) int {
	var rec func(lo, hi int) int
	rec = func(lo, hi int) int {
		if lo == hi {
			return nums[lo] // single element is its own minimum
		}
		if nums[lo] < nums[hi] {
			return nums[lo] // strictly ascending endpoints → sorted range
		}
		mid := lo + (hi-lo)/2 // split point
		// pivot could be in either half when endpoints are equal → check both
		return min(rec(lo, mid), rec(mid+1, hi))
	}
	return rec(0, len(nums)-1)
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Linear Scan) ===")
	fmt.Println(bruteForce([]int{1, 3, 5}))       // 1
	fmt.Println(bruteForce([]int{2, 2, 2, 0, 1})) // 0

	fmt.Println("=== Approach 2: Binary Search with Duplicate Shrinking (Optimal) ===")
	fmt.Println(binarySearch([]int{1, 3, 5}))       // 1
	fmt.Println(binarySearch([]int{2, 2, 2, 0, 1})) // 0

	fmt.Println("=== Approach 3: Divide and Conquer ===")
	fmt.Println(divideAndConquer([]int{1, 3, 5}))       // 1
	fmt.Println(divideAndConquer([]int{2, 2, 2, 0, 1})) // 0
}
