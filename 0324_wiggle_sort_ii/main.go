package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Sort + Reverse Interleave ────────────────────────────────────
//
// sortInterleave solves Wiggle Sort II by sorting the array, then filling the
// odd positions (the "peaks") with the largest halves and the even positions
// (the "valleys") with the smallest halves, each read from the middle outward so
// equal duplicates never land next to each other.
//
// Intuition:
//
//	After sorting, split into a smaller half S and a larger half L. Valleys (even
//	indices 0,2,4,...) must be small and peaks (odd 1,3,5,...) must be large. If
//	we just filled left-to-right, duplicates straddling the median could touch.
//	Filling BOTH halves from their high end downward pushes equal values as far
//	apart as possible, guaranteeing strict wiggle when a valid answer exists.
//
// Algorithm:
//  1. sorted = sort(nums).
//  2. Let n = len; mid = (n+1)/2. S = sorted[:mid] (small), L = sorted[mid:] (large).
//  3. Place S from its top down into even indices 0,2,4,...
//     Place L from its top down into odd indices 1,3,5,...
//  4. Copy result back into nums.
//
// Time:  O(n log n) — dominated by the sort.
// Space: O(n) — the sorted copy plus output buffer.
func sortInterleave(nums []int) {
	n := len(nums)
	sorted := make([]int, n)
	copy(sorted, nums)
	sort.Ints(sorted) // ascending

	mid := (n + 1) / 2 // small half gets the extra element when n is odd
	res := make([]int, n)
	// j walks the small half from its largest element downward → even slots.
	// k walks the large half from its largest element downward → odd slots.
	j, k := mid-1, n-1
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			res[i] = sorted[j] // valley: take next-largest of small half
			j--
		} else {
			res[i] = sorted[k] // peak: take next-largest of large half
			k--
		}
	}
	copy(nums, res) // write answer back in place (by reference)
}

// ── Approach 2: Median + Three-Way Partition + Index Mapping (Optimal) ────────
//
// medianPartition solves Wiggle Sort II in O(n) time and O(1) extra space
// (beyond the input) by finding the median with quickselect, then doing a
// Dutch-National-Flag partition using a virtual index mapping that interleaves
// peaks and valleys directly.
//
// Intuition:
//
//	The sort+interleave idea can be done without a full sort. The median splits
//	values into "greater", "equal", "less". Using a bijective index map that
//	visits odd indices first (peaks) then even indices (valleys), a 3-way
//	partition places greater-than-median at the front (peaks), less-than-median
//	at the back (valleys), and equals in the middle — exactly the wiggle layout,
//	with equal values spread apart.
//
// Algorithm:
//  1. Find the median value m via quickselect (kth = n/2 smallest).
//  2. Define mapped(i) = (2*i + 1) % (n | 1): odd slots then even slots.
//  3. Three-way partition over mapped indices: >m to the left region, <m to the
//     right region, ==m stays; standard DNF pointers i, j (current), k.
//  4. nums is now wiggle-sorted in place.
//
// Time:  O(n) average — quickselect is O(n), partition is O(n).
// Space: O(1) extra — in place beyond the recursion of quickselect.
func medianPartition(nums []int) {
	n := len(nums)
	if n < 2 {
		return // 0 or 1 element is trivially wiggle-sorted
	}
	m := quickselectMedian(nums) // the (n/2)-th smallest = median

	// mapped remaps a logical index into the physical wiggle position:
	// logical 0,1,2,... → physical odd slots 1,3,5,... then even 0,2,4,...
	mapped := func(i int) int { return (2*i + 1) % (n | 1) }

	i, left, right := 0, 0, n-1 // DNF pointers over the logical order
	for i <= right {
		if nums[mapped(i)] > m {
			// greater-than-median belongs among the peaks (front region)
			nums[mapped(left)], nums[mapped(i)] = nums[mapped(i)], nums[mapped(left)]
			left++
			i++
		} else if nums[mapped(i)] < m {
			// less-than-median belongs among the valleys (back region)
			nums[mapped(right)], nums[mapped(i)] = nums[mapped(i)], nums[mapped(right)]
			right--
		} else {
			i++ // equals the median: leave it in the middle band
		}
	}
}

// quickselectMedian returns the (n/2)-th smallest element (the upper median) of a
// COPY of nums, so the original slice order is not disturbed before partitioning.
//
// Time:  O(n) average. Space: O(n) for the copy.
func quickselectMedian(nums []int) int {
	arr := make([]int, len(nums))
	copy(arr, nums)
	target := len(arr) / 2 // 0-based index of the median in sorted order
	lo, hi := 0, len(arr)-1
	for lo < hi {
		p := partition(arr, lo, hi) // Lomuto partition; p is pivot's final index
		if p == target {
			break // pivot sits exactly at the median position
		} else if p < target {
			lo = p + 1 // median is to the right
		} else {
			hi = p - 1 // median is to the left
		}
	}
	return arr[target]
}

// partition is Lomuto partition around arr[hi]; returns the pivot's final index.
func partition(arr []int, lo, hi int) int {
	pivot := arr[hi] // choose last element as pivot
	i := lo          // boundary of elements < pivot
	for j := lo; j < hi; j++ {
		if arr[j] < pivot {
			arr[i], arr[j] = arr[j], arr[i] // push smaller element left
			i++
		}
	}
	arr[i], arr[hi] = arr[hi], arr[i] // move pivot to its sorted spot
	return i
}

// isWiggle verifies the strict property nums[0] < nums[1] > nums[2] < nums[3]...
// Used in main() to prove correctness regardless of which valid layout appears.
func isWiggle(nums []int) bool {
	for i := 1; i < len(nums); i++ {
		if i%2 == 1 { // odd index must be a strict peak
			if !(nums[i] > nums[i-1]) {
				return false
			}
		} else { // even index must be a strict valley
			if !(nums[i] < nums[i-1]) {
				return false
			}
		}
	}
	return true
}

func main() {
	fmt.Println("=== Approach 1: Sort + Reverse Interleave ===")
	a1 := []int{1, 5, 1, 1, 6, 4}
	sortInterleave(a1)
	fmt.Println(a1, "wiggle?", isWiggle(a1)) // expected e.g. [1 6 1 5 1 4] wiggle? true
	a2 := []int{1, 3, 2, 2, 3, 1}
	sortInterleave(a2)
	fmt.Println(a2, "wiggle?", isWiggle(a2)) // expected e.g. [2 3 1 3 1 2] wiggle? true

	fmt.Println("=== Approach 2: Median + 3-Way Partition (Optimal) ===")
	b1 := []int{1, 5, 1, 1, 6, 4}
	medianPartition(b1)
	fmt.Println(b1, "wiggle?", isWiggle(b1)) // expected some valid layout, wiggle? true
	b2 := []int{1, 3, 2, 2, 3, 1}
	medianPartition(b2)
	fmt.Println(b2, "wiggle?", isWiggle(b2)) // expected some valid layout, wiggle? true
}
