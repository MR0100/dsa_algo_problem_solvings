package main

import (
	"fmt"
	"strconv"
)

// ── Approach 1: Single Pass, Track Range Start ───────────────────────────────
//
// summaryRanges scans the sorted array once, and whenever consecutiveness
// breaks it emits the range [start..prev] that just ended.
//
// Intuition:
//
//	The array is sorted and has no duplicates, so a "range" is a maximal run of
//	consecutive integers. Walk left to right remembering where the current run
//	started; the run continues as long as nums[i] == nums[i-1]+1. The moment
//	that fails (or we reach the end), the run [start..nums[i-1]] is complete —
//	format it and begin a new run at nums[i].
//
// Algorithm:
//  1. If empty, return an empty list.
//  2. Let start = nums[0].
//  3. For i from 1 to n: if i == n OR nums[i] != nums[i-1]+1, the run ending at
//     nums[i-1] is done — append its formatted string, then start = nums[i].
//  4. Return the collected strings.
//
// Time:  O(n) — one pass; formatting each range is proportional to output size.
// Space: O(1) extra — ignoring the output list.
func summaryRanges(nums []int) []string {
	res := []string{}
	n := len(nums)
	if n == 0 {
		return res // no numbers → no ranges
	}

	start := nums[0] // first element of the run currently being built
	for i := 1; i <= n; i++ {
		// A run ends when we run off the end OR the next value is not consecutive.
		if i == n || nums[i] != nums[i-1]+1 {
			if start == nums[i-1] { // single-element range "a"
				res = append(res, strconv.Itoa(start))
			} else { // multi-element range "a->b"
				res = append(res, strconv.Itoa(start)+"->"+strconv.Itoa(nums[i-1]))
			}
			if i < n { // begin the next run at the current (non-consecutive) value
				start = nums[i]
			}
		}
	}
	return res
}

// ── Approach 2: Two Pointers (Explicit Window Extension) ─────────────────────
//
// summaryRangesTwoPointers uses an outer pointer i to anchor each range's start
// and an inner pointer j to extend it as far as the run of consecutive numbers
// allows, then jumps i past the finished run.
//
// Intuition:
//
//	Same runs, framed as a window. Fix the left edge i at a range's start, then
//	push a right edge j forward while nums[j+1] == nums[j]+1. When j can go no
//	further, [i..j] is a maximal range; format it and restart with i = j+1. This
//	makes the "extend the window" structure explicit.
//
// Algorithm:
//  1. i = 0.
//  2. While i < n: set j = i; advance j while j+1 < n and nums[j+1] == nums[j]+1.
//  3. Emit [nums[i]..nums[j]] (single or "a->b"); set i = j + 1.
//  4. Return the list.
//
// Time:  O(n) — i and j together advance at most n steps.
// Space: O(1) extra — ignoring the output list.
func summaryRangesTwoPointers(nums []int) []string {
	res := []string{}
	n := len(nums)
	i := 0
	for i < n {
		j := i // extend the right edge as far as consecutiveness holds
		for j+1 < n && nums[j+1] == nums[j]+1 {
			j++
		}
		if i == j { // window is a single element
			res = append(res, strconv.Itoa(nums[i]))
		} else { // window spans nums[i]..nums[j]
			res = append(res, strconv.Itoa(nums[i])+"->"+strconv.Itoa(nums[j]))
		}
		i = j + 1 // jump past the finished range
	}
	return res
}

func main() {
	fmt.Println("=== Approach 1: Single Pass, Track Range Start ===")
	fmt.Println(summaryRanges([]int{0, 1, 2, 4, 5, 7}))    // [0->2 4->5 7]
	fmt.Println(summaryRanges([]int{0, 2, 3, 4, 6, 8, 9})) // [0 2->4 6 8->9]
	fmt.Println(summaryRanges([]int{}))                    // []

	fmt.Println("=== Approach 2: Two Pointers (Explicit Window Extension) ===")
	fmt.Println(summaryRangesTwoPointers([]int{0, 1, 2, 4, 5, 7}))    // [0->2 4->5 7]
	fmt.Println(summaryRangesTwoPointers([]int{0, 2, 3, 4, 6, 8, 9})) // [0 2->4 6 8->9]
	fmt.Println(summaryRangesTwoPointers([]int{}))                    // []
}
