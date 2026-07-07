package main

import "fmt"

// ── Approach 1: Brute Force (Linear Scan) ────────────────────────────────────
//
// bruteForce solves Search in Rotated Sorted Array with a simple linear scan.
//
// Intuition: Ignore the sorted+rotated structure and just scan for the target.
//
// Time:  O(n)
// Space: O(1)
func bruteForce(nums []int, target int) int {
	for i, v := range nums {
		if v == target {
			return i
		}
	}
	return -1
}

// ── Approach 2: Binary Search (Optimal) ──────────────────────────────────────
//
// binarySearch solves Search in Rotated Sorted Array in O(log n) using
// a modified binary search that identifies the sorted half at each step.
//
// Intuition: After rotation, at least one half of [lo..hi] is always sorted.
// We can tell which half by comparing nums[mid] with nums[lo]:
//   - If nums[lo] <= nums[mid]: the LEFT half is sorted.
//     → if target ∈ [nums[lo], nums[mid]): search left; else search right.
//   - Else: the RIGHT half is sorted.
//     → if target ∈ (nums[mid], nums[hi]]: search right; else search left.
//
// Algorithm:
//  1. lo=0, hi=n-1.
//  2. While lo <= hi:
//     mid = (lo+hi)/2.
//     if nums[mid] == target: return mid.
//     if nums[lo] <= nums[mid] (left is sorted):
//       if nums[lo] <= target < nums[mid]: hi = mid-1.
//       else: lo = mid+1.
//     else (right is sorted):
//       if nums[mid] < target <= nums[hi]: lo = mid+1.
//       else: hi = mid-1.
//  3. Return -1.
//
// Time:  O(log n)
// Space: O(1)
func binarySearch(nums []int, target int) int {
	lo, hi := 0, len(nums)-1
	for lo <= hi {
		mid := (lo + hi) / 2
		if nums[mid] == target {
			return mid
		}
		// determine which half is sorted
		if nums[lo] <= nums[mid] { // left half is sorted
			if nums[lo] <= target && target < nums[mid] {
				hi = mid - 1 // target in sorted left half
			} else {
				lo = mid + 1 // target in right half
			}
		} else { // right half is sorted
			if nums[mid] < target && target <= nums[hi] {
				lo = mid + 1 // target in sorted right half
			} else {
				hi = mid - 1 // target in left half
			}
		}
	}
	return -1
}

func main() {
	cases := []struct {
		nums   []int
		target int
		want   int
	}{
		{[]int{4, 5, 6, 7, 0, 1, 2}, 0, 4},
		{[]int{4, 5, 6, 7, 0, 1, 2}, 3, -1},
		{[]int{1}, 0, -1},
		{[]int{1}, 1, 0},
		{[]int{3, 1}, 1, 1},
		{[]int{5, 1, 3}, 5, 0},
	}

	fmt.Println("=== Approach 1: Brute Force ===")
	for _, c := range cases {
		fmt.Printf("nums=%v target=%d  got=%d  expected=%d\n", c.nums, c.target, bruteForce(c.nums, c.target), c.want)
	}

	fmt.Println("\n=== Approach 2: Binary Search (Optimal) ===")
	for _, c := range cases {
		fmt.Printf("nums=%v target=%d  got=%d  expected=%d\n", c.nums, c.target, binarySearch(c.nums, c.target), c.want)
	}
}
