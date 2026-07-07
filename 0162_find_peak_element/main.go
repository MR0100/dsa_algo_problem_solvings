package main

import "fmt"

// ── Approach 1: Brute Force (Linear Scan) ────────────────────────────────────
//
// bruteForce solves Find Peak Element by scanning left to right for the
// first "downhill" step.
//
// Intuition:
//
//	Because nums[-1] and nums[n] count as −∞ and neighbours are never equal,
//	the array starts by (conceptually) rising from −∞. The first index i
//	where nums[i] > nums[i+1] is therefore a peak: it is greater than its
//	right neighbour by the test, and greater than its left neighbour because
//	we only reached i by climbing. If no such i exists the array is strictly
//	increasing, so the last element is the peak.
//
// Algorithm:
//  1. For i from 0 to n−2: if nums[i] > nums[i+1], return i.
//  2. If the loop finishes, return n−1.
//
// Time:  O(n) — one pass in the worst case (strictly increasing array).
// Space: O(1) — only the loop index.
func bruteForce(nums []int) int {
	for i := 0; i < len(nums)-1; i++ {
		if nums[i] > nums[i+1] { // first downhill step → i is a peak
			return i
		}
	}
	return len(nums) - 1 // strictly increasing → last element is the peak
}

// ── Approach 2: Recursive Binary Search ──────────────────────────────────────
//
// recursiveBinarySearch solves Find Peak Element by recursively halving the
// search space toward a rising slope.
//
// Intuition:
//
//	Look at the middle element and its right neighbour. If nums[mid] >
//	nums[mid+1] we are on a descending slope, so some peak exists at mid or
//	to its left (the values must have risen from −∞ before falling here).
//	Otherwise nums[mid] < nums[mid+1]: we are ascending, so a peak must lie
//	strictly to the right (values eventually fall back to −∞ at the end).
//	Either way, half of the range is discarded — a peak is never lost.
//
// Algorithm:
//  1. If lo == hi, that single index is a peak — return it.
//  2. mid = lo + (hi−lo)/2.
//  3. If nums[mid] > nums[mid+1], recurse on [lo, mid].
//  4. Else recurse on [mid+1, hi].
//
// Time:  O(log n) — the range halves at every level of recursion.
// Space: O(log n) — recursion stack depth.
func recursiveBinarySearch(nums []int) int {
	return peakHelper(nums, 0, len(nums)-1)
}

// peakHelper narrows [lo, hi] (which always contains a peak) to one index.
func peakHelper(nums []int, lo, hi int) int {
	if lo == hi { // range shrunk to one candidate → it is a peak
		return lo
	}
	mid := lo + (hi-lo)/2        // overflow-safe midpoint; mid < hi so mid+1 is valid
	if nums[mid] > nums[mid+1] { // descending slope → peak is at mid or left of it
		return peakHelper(nums, lo, mid)
	}
	return peakHelper(nums, mid+1, hi) // ascending slope → peak is right of mid
}

// ── Approach 3: Iterative Binary Search (Optimal) ────────────────────────────
//
// iterativeBinarySearch solves Find Peak Element with the same slope-chasing
// idea as Approach 2, but in a loop with O(1) space.
//
// Intuition:
//
//	Maintain the invariant "the range [lo, hi] contains at least one peak".
//	Comparing nums[mid] with nums[mid+1] tells us which side of mid must
//	contain a peak, so we can shrink the range by half each iteration until
//	a single index remains. This meets the required O(log n) bound.
//
// Algorithm:
//  1. lo = 0, hi = n−1.
//  2. While lo < hi:
//     a. mid = lo + (hi−lo)/2.
//     b. If nums[mid] > nums[mid+1], a peak is in [lo, mid] → hi = mid.
//     c. Else a peak is in [mid+1, hi] → lo = mid+1.
//  3. Return lo (== hi).
//
// Time:  O(log n) — the search range halves every iteration.
// Space: O(1) — two pointers and a midpoint.
func iterativeBinarySearch(nums []int) int {
	lo, hi := 0, len(nums)-1 // invariant: [lo, hi] always contains a peak
	for lo < hi {
		mid := lo + (hi-lo)/2 // mid < hi, so nums[mid+1] is always in bounds
		if nums[mid] > nums[mid+1] {
			hi = mid // descending → mid itself may be the peak; keep it
		} else {
			lo = mid + 1 // ascending → mid cannot be a peak; drop it
		}
	}
	return lo // range collapsed to the peak index
}

func main() {
	examples := [][]int{
		{1, 2, 3, 1},          // expected 2 (value 3 is the only peak)
		{1, 2, 1, 3, 5, 6, 4}, // expected 1 or 5 (values 2 and 6 are both peaks)
	}

	fmt.Println("=== Approach 1: Brute Force (Linear Scan) ===")
	for _, ex := range examples {
		fmt.Printf("nums=%v  got=%d\n", ex, bruteForce(ex)) // expected 2, then 1 (index 5 also valid)
	}

	fmt.Println("=== Approach 2: Recursive Binary Search ===")
	for _, ex := range examples {
		fmt.Printf("nums=%v  got=%d\n", ex, recursiveBinarySearch(ex)) // expected 2, then 5 (index 1 also valid)
	}

	fmt.Println("=== Approach 3: Iterative Binary Search (Optimal) ===")
	for _, ex := range examples {
		fmt.Printf("nums=%v  got=%d\n", ex, iterativeBinarySearch(ex)) // expected 2, then 5 (index 1 also valid)
	}
}
