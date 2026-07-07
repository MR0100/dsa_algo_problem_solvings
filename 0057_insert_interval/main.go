package main

import "fmt"

// ── Approach 1: Linear Scan with Three Phases ─────────────────────────────────
//
// insert solves Insert Interval in a single pass over the sorted non-overlapping
// interval list.
//
// Intuition:
//   The input is already sorted and non-overlapping. Walk through three phases:
//   1. Add all intervals that end before newInterval starts (no overlap).
//   2. Merge all intervals that overlap with newInterval (expand newInterval).
//   3. Add all remaining intervals after newInterval ends.
//
// Algorithm:
//   i = 0
//   Phase 1: while i<n and intervals[i][1] < newInterval[0]: append intervals[i]; i++
//   Phase 2: while i<n and intervals[i][0] <= newInterval[1]:
//              newInterval[0] = min(...); newInterval[1] = max(...)
//              i++
//            append newInterval
//   Phase 3: append intervals[i:]
//
// Time:  O(n) — single pass; each interval examined once.
// Space: O(n) — output list.
func insert(intervals [][]int, newInterval []int) [][]int {
	result := [][]int{}
	i, n := 0, len(intervals)

	// Phase 1: intervals entirely before newInterval
	for i < n && intervals[i][1] < newInterval[0] {
		result = append(result, intervals[i])
		i++
	}

	// Phase 2: merge all overlapping intervals into newInterval
	for i < n && intervals[i][0] <= newInterval[1] {
		if intervals[i][0] < newInterval[0] {
			newInterval[0] = intervals[i][0] // extend start
		}
		if intervals[i][1] > newInterval[1] {
			newInterval[1] = intervals[i][1] // extend end
		}
		i++
	}
	result = append(result, newInterval)

	// Phase 3: intervals entirely after newInterval
	result = append(result, intervals[i:]...)
	return result
}

func main() {
	fmt.Println("=== Approach 1: Linear Scan (Three Phases) ===")

	i1 := [][]int{{1, 3}, {6, 9}}
	fmt.Printf("intervals=%v newInterval=[2,5]  got=%v  expected [[1 5] [6 9]]\n",
		i1, insert(i1, []int{2, 5}))

	i2 := [][]int{{1, 2}, {3, 5}, {6, 7}, {8, 10}, {12, 16}}
	fmt.Printf("intervals=%v newInterval=[4,8]  got=%v  expected [[1 2] [3 10] [12 16]]\n",
		i2, insert(i2, []int{4, 8}))

	i3 := [][]int{}
	fmt.Printf("intervals=[] newInterval=[5,7]  got=%v  expected [[5 7]]\n",
		insert(i3, []int{5, 7}))

	i4 := [][]int{{1, 5}}
	fmt.Printf("intervals=%v newInterval=[2,3]  got=%v  expected [[1 5]]\n",
		i4, insert(i4, []int{2, 3}))

	i5 := [][]int{{1, 5}}
	fmt.Printf("intervals=%v newInterval=[6,8]  got=%v  expected [[1 5] [6 8]]\n",
		i5, insert(i5, []int{6, 8}))
}
