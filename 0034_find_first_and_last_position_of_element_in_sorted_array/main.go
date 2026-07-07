package main

import "fmt"

// ── Approach 1: Linear Scan ───────────────────────────────────────────────────
//
// linearScan solves Find First and Last Position of Element in Sorted Array
// with a single pass to find the first and last occurrence.
//
// Intuition: Walk the array; record the first index where target appears and
// keep updating the last index as long as target appears.
//
// Time:  O(n)
// Space: O(1)
func linearScan(nums []int, target int) []int {
	first, last := -1, -1
	for i, v := range nums {
		if v == target {
			if first == -1 {
				first = i
			}
			last = i
		}
	}
	return []int{first, last}
}

// ── Approach 2: Two Binary Searches (Optimal) ────────────────────────────────
//
// binarySearch solves Find First and Last Position in O(log n) by running two
// separate binary searches — one for the leftmost and one for the rightmost
// occurrence.
//
// Intuition: Standard binary search finds *a* position. By biasing mid toward
// the left (when nums[mid]==target, record mid and continue searching left)
// we find the first occurrence. The mirror bias finds the last.
//
// Algorithm — findFirst:
//  lo=0, hi=n-1, result=-1
//  while lo<=hi:
//    mid=(lo+hi)/2
//    if nums[mid]==target: result=mid; hi=mid-1 (keep searching left)
//    elif nums[mid]<target: lo=mid+1
//    else: hi=mid-1
//
// Algorithm — findLast: same but on match: result=mid; lo=mid+1
//
// Time:  O(log n) — two binary searches
// Space: O(1)
func binarySearch(nums []int, target int) []int {
	return []int{findFirst(nums, target), findLast(nums, target)}
}

func findFirst(nums []int, target int) int {
	lo, hi, result := 0, len(nums)-1, -1
	for lo <= hi {
		mid := (lo + hi) / 2
		if nums[mid] == target {
			result = mid   // record candidate
			hi = mid - 1  // keep searching left for an earlier occurrence
		} else if nums[mid] < target {
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	return result
}

func findLast(nums []int, target int) int {
	lo, hi, result := 0, len(nums)-1, -1
	for lo <= hi {
		mid := (lo + hi) / 2
		if nums[mid] == target {
			result = mid  // record candidate
			lo = mid + 1 // keep searching right for a later occurrence
		} else if nums[mid] < target {
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	return result
}

func main() {
	cases := []struct {
		nums   []int
		target int
		want   []int
	}{
		{[]int{5, 7, 7, 8, 8, 10}, 8, []int{3, 4}},
		{[]int{5, 7, 7, 8, 8, 10}, 6, []int{-1, -1}},
		{[]int{}, 0, []int{-1, -1}},
		{[]int{1}, 1, []int{0, 0}},
		{[]int{1, 1, 1, 1}, 1, []int{0, 3}},
	}

	fmt.Println("=== Approach 1: Linear Scan ===")
	for _, c := range cases {
		got := linearScan(c.nums, c.target)
		fmt.Printf("nums=%v target=%d  got=%v  expected=%v\n", c.nums, c.target, got, c.want)
	}

	fmt.Println("\n=== Approach 2: Binary Search (Optimal) ===")
	for _, c := range cases {
		got := binarySearch(c.nums, c.target)
		fmt.Printf("nums=%v target=%d  got=%v  expected=%v\n", c.nums, c.target, got, c.want)
	}
}
