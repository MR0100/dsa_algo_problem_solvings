package main

import (
	"fmt"
	"math"
	"sort"
)

// ── Approach 1: Merge then Find Median ───────────────────────────────────────
//
// mergeAndFind merges both sorted arrays into a single sorted array, then
// reads the median directly from the middle position(s).
//
// Intuition:
//   The median of a combined set of m+n numbers is the middle element (odd
//   total) or the average of the two middle elements (even total). If we had
//   the merged sorted array, finding the median is O(1). The cost is the merge.
//
// Algorithm:
//   1. Merge nums1 and nums2 using the merge step of merge sort.
//   2. If (m+n) is odd  → return merged[(m+n)/2].
//      If (m+n) is even → return (merged[(m+n)/2-1] + merged[(m+n)/2]) / 2.0.
//
// Time:  O(m+n) — one pass to merge.
// Space: O(m+n) — the merged array.
func mergeAndFind(nums1 []int, nums2 []int) float64 {
	// Standard two-pointer merge of two sorted arrays.
	merged := make([]int, 0, len(nums1)+len(nums2))
	i, j := 0, 0
	for i < len(nums1) && j < len(nums2) {
		if nums1[i] <= nums2[j] {
			merged = append(merged, nums1[i])
			i++
		} else {
			merged = append(merged, nums2[j])
			j++
		}
	}
	merged = append(merged, nums1[i:]...)
	merged = append(merged, nums2[j:]...)

	total := len(merged)
	mid := total / 2
	if total%2 == 1 {
		return float64(merged[mid])
	}
	return float64(merged[mid-1]+merged[mid]) / 2.0
}

// ── Approach 2: Concatenate and Sort ─────────────────────────────────────────
//
// concatAndSort naively concatenates both arrays, sorts, then reads the median.
//
// Intuition:
//   The simplest imaginable approach: ignore the "already sorted" property,
//   combine everything, re-sort from scratch.
//
// Algorithm:
//   1. combined = append(nums1, nums2...).
//   2. sort.Ints(combined).
//   3. Read median.
//
// Time:  O((m+n) log(m+n)) — dominated by the sort.
// Space: O(m+n) — the combined array.
//
// Note: strictly worse than mergeAndFind; included to show the naive baseline.
func concatAndSort(nums1 []int, nums2 []int) float64 {
	combined := make([]int, len(nums1)+len(nums2))
	copy(combined, nums1)
	copy(combined[len(nums1):], nums2)
	sort.Ints(combined)

	total := len(combined)
	mid := total / 2
	if total%2 == 1 {
		return float64(combined[mid])
	}
	return float64(combined[mid-1]+combined[mid]) / 2.0
}

// ── Approach 3: Two-Pointer Walk to Median ───────────────────────────────────
//
// twoPointerWalk uses two pointers (one per array) and counts steps up to the
// median position without building the full merged array.
//
// Intuition:
//   Instead of storing the merged array we can advance the pointers together,
//   counting until we reach position (m+n)/2. We only need to remember the
//   last two values seen (for even-total median).
//
// Algorithm:
//   1. Advance whichever pointer has the smaller current value, counting steps.
//   2. Stop at index (m+n-1)/2 and (m+n)/2; those are the two "middle" values.
//   3. Return their average (or the single middle for odd total).
//
// Time:  O(m+n) — same as mergeAndFind but O(1) space.
// Space: O(1) — only pointers and two stored values.
func twoPointerWalk(nums1 []int, nums2 []int) float64 {
	m, n := len(nums1), len(nums2)
	total := m + n
	// We need positions (total-1)/2 and total/2 (they're the same for odd total).
	targetHigh := total / 2
	targetLow := targetHigh - 1 + (total % 2) // (total-1)/2, kept for clarity

	i, j := 0, 0
	prev, cur := 0, 0 // values at positions targetLow and targetHigh

	_ = targetLow // used conceptually; targetHigh drives the loop
	for step := 0; step <= targetHigh; step++ {
		prev = cur // shift: prev becomes what was at step-1
		// Pick the smaller of the two current front elements.
		if i < m && (j >= n || nums1[i] <= nums2[j]) {
			cur = nums1[i]
			i++
		} else {
			cur = nums2[j]
			j++
		}
	}

	if total%2 == 1 {
		return float64(cur) // odd total: single middle element
	}
	return float64(prev+cur) / 2.0 // even total: average of two middles
}

// ── Approach 4: Binary Search on Partition (Optimal) ─────────────────────────
//
// binarySearchPartition finds the correct partition of the smaller array using
// binary search, achieving O(log(min(m,n))) time.
//
// Intuition:
//   The median divides the combined sorted sequence into two equal halves.
//   We want to partition nums1 at position i and nums2 at position j such that:
//     i + j = (m + n + 1) / 2          (left half has ⌈(m+n)/2⌉ elements)
//     nums1[i-1] ≤ nums2[j]            (max of left ≤ min of right)
//     nums2[j-1] ≤ nums1[i]            (max of left ≤ min of right)
//
//   Binary-search i over [0, m]. j is derived: j = half - i.
//   If nums1[i-1] > nums2[j]: i is too large → search left half.
//   If nums2[j-1] > nums1[i]: i is too small → search right half.
//   When correct:
//     median = max(left halves)             (odd total)
//     median = (max(left) + min(right)) / 2 (even total)
//
//   Always binary-search on the shorter array to keep the range O(log(min(m,n))).
//
// Time:  O(log(min(m,n))) — binary search on the smaller array.
// Space: O(1) — only index variables and boundary values.
func binarySearchPartition(nums1 []int, nums2 []int) float64 {
	// Ensure nums1 is the shorter array so we binary-search on it.
	if len(nums1) > len(nums2) {
		nums1, nums2 = nums2, nums1
	}
	m, n := len(nums1), len(nums2)
	half := (m + n + 1) / 2 // size of the left partition

	lo, hi := 0, m

	for lo <= hi {
		i := (lo + hi) / 2 // cut in nums1: i elements go to the left half
		j := half - i      // cut in nums2: j elements go to the left half

		// Boundary values; use ±Inf for edge cuts to avoid bounds-checking.
		nums1LeftMax := math.MinInt64
		if i > 0 {
			nums1LeftMax = nums1[i-1]
		}
		nums1RightMin := math.MaxInt64
		if i < m {
			nums1RightMin = nums1[i]
		}
		nums2LeftMax := math.MinInt64
		if j > 0 {
			nums2LeftMax = nums2[j-1]
		}
		nums2RightMin := math.MaxInt64
		if j < n {
			nums2RightMin = nums2[j]
		}

		if nums1LeftMax <= nums2RightMin && nums2LeftMax <= nums1RightMin {
			// Correct partition found.
			leftMax := nums1LeftMax
			if nums2LeftMax > leftMax {
				leftMax = nums2LeftMax
			}
			if (m+n)%2 == 1 {
				return float64(leftMax) // odd total: median is max of left half
			}
			rightMin := nums1RightMin
			if nums2RightMin < rightMin {
				rightMin = nums2RightMin
			}
			return float64(leftMax+rightMin) / 2.0
		} else if nums1LeftMax > nums2RightMin {
			hi = i - 1 // i too large, move partition left in nums1
		} else {
			lo = i + 1 // i too small, move partition right in nums1
		}
	}
	return 0.0 // unreachable for valid input
}

func main() {
	examples := []struct {
		nums1, nums2 []int
		expect       float64
	}{
		{[]int{1, 3}, []int{2}, 2.0},
		{[]int{1, 2}, []int{3, 4}, 2.5},
		{[]int{0, 0}, []int{0, 0}, 0.0},
		{[]int{}, []int{1}, 1.0},
		{[]int{2}, []int{}, 2.0},
	}

	approaches := []struct {
		name string
		fn   func([]int, []int) float64
	}{
		{"Approach 1: Merge & Find          O(m+n) T | O(m+n) S", mergeAndFind},
		{"Approach 2: Concat & Sort         O((m+n)log(m+n)) T | O(m+n) S", concatAndSort},
		{"Approach 3: Two-Pointer Walk      O(m+n) T | O(1)   S", twoPointerWalk},
		{"Approach 4: Binary Search (Opt) ✅ O(log(min(m,n))) T | O(1) S", binarySearchPartition},
	}

	for _, ex := range examples {
		fmt.Printf("nums1=%v  nums2=%v  expect=%.5f\n", ex.nums1, ex.nums2, ex.expect)
		for _, ap := range approaches {
			result := ap.fn(ex.nums1, ex.nums2)
			fmt.Printf("  %-65s → %.5f\n", ap.name, result)
		}
		fmt.Println()
	}
}
