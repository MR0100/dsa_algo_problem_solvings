package main

import "fmt"

// ── Approach 1: Linear Scan ───────────────────────────────────────────────────
//
// linearScan solves Search Insert Position with a forward scan.
//
// Intuition: Walk the sorted array; the answer is the first index where
// nums[i] >= target. If no such index exists, the target belongs at the end.
//
// Time:  O(n)
// Space: O(1)
func linearScan(nums []int, target int) int {
	for i, v := range nums {
		if v >= target {
			return i
		}
	}
	return len(nums) // target is larger than all elements
}

// ── Approach 2: Binary Search (Optimal) ──────────────────────────────────────
//
// binarySearch solves Search Insert Position in O(log n).
//
// Intuition: This is exactly the "lower_bound" query: find the leftmost index
// where nums[i] >= target. Binary search naturally converges to this:
// - When nums[mid] < target: the insertion point is to the right → lo=mid+1.
// - When nums[mid] >= target: could be the answer, but there might be an
//   earlier occurrence → hi=mid-1, but record mid as a candidate via lo.
// At the end of the loop, lo == the insertion position.
//
// Key insight: lo converges to the smallest index where nums[lo] >= target,
// which is exactly the insertion position whether or not target exists.
//
// Algorithm:
//  lo=0, hi=n-1
//  while lo<=hi:
//    mid=(lo+hi)/2
//    if nums[mid] < target: lo=mid+1
//    else: hi=mid-1
//  return lo
//
// Time:  O(log n)
// Space: O(1)
func binarySearch(nums []int, target int) int {
	lo, hi := 0, len(nums)-1
	for lo <= hi {
		mid := (lo + hi) / 2
		if nums[mid] < target {
			lo = mid + 1 // target is in the right half
		} else {
			hi = mid - 1 // this index could be the answer; search left for earlier
		}
	}
	return lo // lo is the insertion point
}

func main() {
	cases := []struct {
		nums   []int
		target int
		want   int
	}{
		{[]int{1, 3, 5, 6}, 5, 2},
		{[]int{1, 3, 5, 6}, 2, 1},
		{[]int{1, 3, 5, 6}, 7, 4},
		{[]int{1, 3, 5, 6}, 0, 0},
		{[]int{1}, 0, 0},
		{[]int{1}, 2, 1},
	}

	fmt.Println("=== Approach 1: Linear Scan ===")
	for _, c := range cases {
		got := linearScan(c.nums, c.target)
		fmt.Printf("nums=%v target=%d  got=%d  expected=%d\n", c.nums, c.target, got, c.want)
	}

	fmt.Println("\n=== Approach 2: Binary Search (Optimal) ===")
	for _, c := range cases {
		got := binarySearch(c.nums, c.target)
		fmt.Printf("nums=%v target=%d  got=%d  expected=%d\n", c.nums, c.target, got, c.want)
	}
}
