package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Brute Force Pairwise Check ───────────────────────────────────
//
// bruteForce solves Meeting Rooms by comparing every pair of intervals to see
// whether any two overlap.
//
// Intuition:
//
//	A person can attend all meetings iff no two meetings overlap. The most direct
//	test is to check all O(n²) pairs. Two intervals [s1,e1] and [s2,e2] overlap
//	when s1 < e2 AND s2 < e1 (strict, because touching at an endpoint is allowed:
//	a meeting ending at 10 and another starting at 10 are fine).
//
// Algorithm:
//  1. For every pair i < j, test overlap: intervals[i][0] < intervals[j][1]
//     && intervals[j][0] < intervals[i][1].
//  2. If any pair overlaps, return false.
//  3. Otherwise return true.
//
// Time:  O(n²) — all pairs.
// Space: O(1).
func bruteForce(intervals [][]int) bool {
	n := len(intervals)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			// Overlap test with strict inequalities so shared endpoints are OK.
			if intervals[i][0] < intervals[j][1] && intervals[j][0] < intervals[i][1] {
				return false // found two meetings that clash
			}
		}
	}
	return true // no clashing pair ⇒ all meetings attendable
}

// ── Approach 2: Sort by Start Time (Optimal) ─────────────────────────────────
//
// sortByStart solves Meeting Rooms by sorting on start time and checking only
// adjacent intervals for overlap.
//
// Intuition:
//
//	Once meetings are sorted by start time, a conflict can only occur between
//	consecutive meetings: if meeting i does not overlap meeting i+1, and starts
//	are non-decreasing, it cannot overlap any later meeting either. So one linear
//	sweep after sorting suffices — a meeting conflicts exactly when it starts
//	before the previous one ended.
//
// Algorithm:
//  1. Sort intervals by start time.
//  2. Walk i = 1..n-1; if intervals[i][0] < intervals[i-1][1], return false.
//  3. Return true.
//
// Time:  O(n log n) — dominated by the sort.
// Space: O(1) extra (in-place sort; O(n) if you count sort scratch).
func sortByStart(intervals [][]int) bool {
	// Sort ascending by start time.
	sort.Slice(intervals, func(a, b int) bool {
		return intervals[a][0] < intervals[b][0]
	})
	for i := 1; i < len(intervals); i++ {
		// Current meeting starts before the previous one ended ⇒ overlap.
		if intervals[i][0] < intervals[i-1][1] {
			return false
		}
	}
	return true // swept all adjacent pairs with no conflict
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(bruteForce([][]int{{0, 30}, {5, 10}, {15, 20}})) // expected false
	fmt.Println(bruteForce([][]int{{7, 10}, {2, 4}}))            // expected true

	fmt.Println("=== Approach 2: Sort by Start (Optimal) ===")
	fmt.Println(sortByStart([][]int{{0, 30}, {5, 10}, {15, 20}})) // expected false
	fmt.Println(sortByStart([][]int{{7, 10}, {2, 4}}))            // expected true
}
