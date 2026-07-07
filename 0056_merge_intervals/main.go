package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Sort + Linear Merge ──────────────────────────────────────────
//
// merge solves Merge Intervals by sorting on start time then merging overlapping
// intervals in a single pass.
//
// Intuition:
//   After sorting by start, any overlapping intervals are adjacent. Walk through
//   and merge each interval into the last element of the result if they overlap
//   (current start <= last end). Otherwise append as a new interval.
//
// Algorithm:
//   sort intervals by start
//   result = [intervals[0]]
//   for each interval in intervals[1:]:
//     last = result[len-1]
//     if interval[0] <= last[1]: last[1] = max(last[1], interval[1])
//     else: result = append(result, interval)
//
// Time:  O(n log n) — dominated by sort; merge pass is O(n).
// Space: O(n)       — output list.
func merge(intervals [][]int) [][]int {
	if len(intervals) == 0 {
		return nil
	}
	// sort by start time
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	result := [][]int{intervals[0]}
	for _, curr := range intervals[1:] {
		last := result[len(result)-1]
		if curr[0] <= last[1] {
			// overlap: extend the end if needed
			if curr[1] > last[1] {
				last[1] = curr[1]
			}
		} else {
			result = append(result, curr)
		}
	}
	return result
}

func main() {
	fmt.Println("=== Approach 1: Sort + Linear Merge ===")

	i1 := [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}
	fmt.Printf("intervals=%v  got=%v  expected [[1 6] [8 10] [15 18]]\n", i1, merge(i1))

	i2 := [][]int{{1, 4}, {4, 5}}
	fmt.Printf("intervals=%v  got=%v  expected [[1 5]]\n", i2, merge(i2))

	i3 := [][]int{{1, 4}, {0, 4}}
	fmt.Printf("intervals=%v  got=%v  expected [[0 4]]\n", i3, merge(i3))

	i4 := [][]int{{1, 4}, {2, 3}}
	fmt.Printf("intervals=%v  got=%v  expected [[1 4]]\n", i4, merge(i4))
}
