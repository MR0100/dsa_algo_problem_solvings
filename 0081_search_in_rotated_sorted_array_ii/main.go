package main

import "fmt"

// ── Approach 1: Linear Scan ───────────────────────────────────────────────────
//
// linearScan solves Search in Rotated Sorted Array II by scanning all elements.
//
// Intuition:
//   With duplicates we can't always determine which half is sorted.
//   Fall back to O(n) linear search.
//
// Time:  O(n)
// Space: O(1)
func linearScan(nums []int, target int) bool {
	for _, v := range nums {
		if v == target {
			return true
		}
	}
	return false
}

// ── Approach 2: Binary Search (with duplicate skip) ───────────────────────────
//
// binarySearch solves Search in Rotated Sorted Array II using modified binary
// search that handles duplicates by skipping equal boundary values.
//
// Intuition:
//   Standard binary search on a rotated array (#33) fails when nums[lo]==nums[mid]
//   because we can't tell which half is sorted. Fix: when nums[lo]==nums[mid],
//   increment lo (skip one duplicate). This degrades to O(n) in the worst case
//   (all same values) but stays O(log n) on average.
//
// Algorithm:
//   lo=0, hi=n-1
//   while lo <= hi:
//     mid = lo + (hi-lo)/2
//     if nums[mid] == target: return true
//     if nums[lo] == nums[mid]: lo++  // can't determine which side; skip
//     else if nums[lo] <= nums[mid]:  // left half is sorted
//       if target in [nums[lo], nums[mid]): hi=mid-1
//       else: lo=mid+1
//     else:  // right half is sorted
//       if target in (nums[mid], nums[hi]]: lo=mid+1
//       else: hi=mid-1
//   return false
//
// Time:  O(log n) average, O(n) worst case (all duplicates).
// Space: O(1)
func binarySearch(nums []int, target int) bool {
	lo, hi := 0, len(nums)-1
	for lo <= hi {
		mid := lo + (hi-lo)/2
		if nums[mid] == target {
			return true
		}
		if nums[lo] == nums[mid] {
			lo++ // can't determine sorted half; skip one duplicate
			continue
		}
		if nums[lo] <= nums[mid] {
			// left half [lo..mid] is sorted
			if nums[lo] <= target && target < nums[mid] {
				hi = mid - 1
			} else {
				lo = mid + 1
			}
		} else {
			// right half [mid..hi] is sorted
			if nums[mid] < target && target <= nums[hi] {
				lo = mid + 1
			} else {
				hi = mid - 1
			}
		}
	}
	return false
}

func main() {
	fmt.Println("=== Approach 1: Linear Scan ===")
	fmt.Printf("nums=[2,5,6,0,0,1,2] target=0  got=%v  expected true\n", linearScan([]int{2, 5, 6, 0, 0, 1, 2}, 0))
	fmt.Printf("nums=[2,5,6,0,0,1,2] target=3  got=%v  expected false\n", linearScan([]int{2, 5, 6, 0, 0, 1, 2}, 3))

	fmt.Println("=== Approach 2: Binary Search ===")
	fmt.Printf("nums=[2,5,6,0,0,1,2] target=0  got=%v  expected true\n", binarySearch([]int{2, 5, 6, 0, 0, 1, 2}, 0))
	fmt.Printf("nums=[2,5,6,0,0,1,2] target=3  got=%v  expected false\n", binarySearch([]int{2, 5, 6, 0, 0, 1, 2}, 3))
	fmt.Printf("nums=[1,0,1,1,1] target=0  got=%v  expected true\n", binarySearch([]int{1, 0, 1, 1, 1}, 0))
	fmt.Printf("nums=[1,1,1,1,1] target=2  got=%v  expected false\n", binarySearch([]int{1, 1, 1, 1, 1}, 2))
}
