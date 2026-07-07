package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: DP (Longest Non-Overlapping Chain) ───────────────────────────
//
// dpLongestChain removes the fewest intervals by first finding the LARGEST set
// of mutually non-overlapping intervals to keep; the answer is total − kept.
//
// Intuition:
//
//	"Minimum removals to make the rest non-overlapping" is the complement of
//	"maximum intervals we can keep that don't overlap". If we can keep k of the n
//	intervals, we must delete n − k. Finding the max keep-set is a weighted-
//	interval / longest-chain problem: sort by start, then dp[i] = the largest
//	non-overlapping chain ending at interval i, extending any earlier interval j
//	whose end ≤ start[i]. This mirrors Longest Increasing Subsequence but with
//	the "increasing" test replaced by "non-overlapping".
//
// Algorithm:
//
//  1. Sort intervals by start.
//  2. dp[i] = 1 (interval i alone). For each i, for each j < i with
//     intervals[j].end ≤ intervals[i].start: dp[i] = max(dp[i], dp[j] + 1).
//  3. keep = max(dp); return n − keep.
//
// Time:  O(n²) — the nested chain relaxation (n intervals × n predecessors).
// Space: O(n) — the dp table.
func dpLongestChain(intervals [][]int) int {
	n := len(intervals)
	if n == 0 {
		return 0
	}
	// Sort by start so a valid predecessor for i is any earlier-ending j.
	sort.Slice(intervals, func(a, b int) bool {
		return intervals[a][0] < intervals[b][0]
	})
	dp := make([]int, n) // dp[i] = longest non-overlapping chain ending at i
	keep := 0            // best chain length seen overall
	for i := 0; i < n; i++ {
		dp[i] = 1 // interval i can always stand alone
		for j := 0; j < i; j++ {
			// j can precede i iff it finishes at or before i starts
			// (touching endpoints like [1,2] and [2,3] don't overlap).
			if intervals[j][1] <= intervals[i][0] && dp[j]+1 > dp[i] {
				dp[i] = dp[j] + 1 // extend j's chain by interval i
			}
		}
		if dp[i] > keep {
			keep = dp[i] // track the maximum keep-set size
		}
	}
	return n - keep // everything not in the largest kept chain must be removed
}

// ── Approach 2: Greedy by Earliest End Time (Optimal) ─────────────────────────
//
// greedyEarliestEnd removes the fewest intervals by scanning left to right and,
// whenever two intervals clash, discarding the one that ends later.
//
// Intuition:
//
//	This is the classic "activity selection" argument. Sort by END time. Greedily
//	keep an interval iff it starts at or after the end of the last interval we
//	kept — because among clashing intervals, the one that finishes earliest
//	leaves the most room for everything after it, so it is always safe to keep
//	the earlier-ending one and drop the later-ending one. Every time we must drop
//	one, that's a forced removal. Counting the removals directly avoids the O(n²)
//	chain DP.
//
// Algorithm:
//
//  1. Sort intervals by end.
//  2. Track prevEnd = end of the last kept interval (start at −∞).
//  3. For each interval in order: if its start ≥ prevEnd, keep it (update
//     prevEnd = its end); otherwise it overlaps → remove it (removals++).
//  4. Return removals.
//
// Time:  O(n log n) — dominated by the sort; the scan is O(n).
// Space: O(1) — a couple of scalars (in-place sort aside).
func greedyEarliestEnd(intervals [][]int) int {
	n := len(intervals)
	if n == 0 {
		return 0
	}
	// Sort by END time: the earliest finisher is the safest to keep.
	sort.Slice(intervals, func(a, b int) bool {
		return intervals[a][1] < intervals[b][1]
	})
	removals := 0
	prevEnd := intervals[0][1] // keep the very first (earliest-ending) interval
	for i := 1; i < n; i++ {
		if intervals[i][0] >= prevEnd {
			// No overlap with the last kept interval → keep this one.
			prevEnd = intervals[i][1] // advance the "kept until" boundary
		} else {
			// Overlaps: it ends no earlier than prevEnd (we sorted by end),
			// so dropping THIS one is optimal.
			removals++
		}
	}
	return removals
}

// ── Approach 3: Greedy by Earliest Start (Keep Shorter on Clash) ──────────────
//
// greedyEarliestStart sorts by start and, on every overlap, keeps whichever of
// the two clashing intervals ends earlier (discarding the later-ending one).
//
// Intuition:
//
//	Same greedy optimum reached from a start-sorted view. Walking in start order,
//	when the current interval overlaps the last kept one, one of them must go —
//	keep the one with the smaller end (it blocks the least future space) and drop
//	the other. Because we sorted by start, "keep the smaller end" is a simple
//	min() on prevEnd. This shows the greedy choice is really about END times no
//	matter which key you sort on.
//
// Algorithm:
//
//  1. Sort intervals by start.
//  2. prevEnd = end of the first interval.
//  3. For each next interval: if its start ≥ prevEnd → no clash, keep it
//     (prevEnd = its end). Else clash → removals++ and set
//     prevEnd = min(prevEnd, its end) (retain the earlier finisher).
//  4. Return removals.
//
// Time:  O(n log n) — the sort dominates.
// Space: O(1) — scalar state.
func greedyEarliestStart(intervals [][]int) int {
	n := len(intervals)
	if n == 0 {
		return 0
	}
	// Sort by START time this time.
	sort.Slice(intervals, func(a, b int) bool {
		return intervals[a][0] < intervals[b][0]
	})
	removals := 0
	prevEnd := intervals[0][1]
	for i := 1; i < n; i++ {
		if intervals[i][0] >= prevEnd {
			// Disjoint from the kept interval → keep it.
			prevEnd = intervals[i][1]
		} else {
			// Overlap: drop one, and keep whichever ends earlier to leave
			// maximum room for later intervals.
			removals++
			if intervals[i][1] < prevEnd {
				prevEnd = intervals[i][1] // the new interval finishes sooner
			}
		}
	}
	return removals
}

func main() {
	ex1 := [][]int{{1, 2}, {2, 3}, {3, 4}, {1, 3}}
	ex2 := [][]int{{1, 2}, {1, 2}, {1, 2}}
	ex3 := [][]int{{1, 2}, {2, 3}}

	// Note: each approach sorts its input in place, so pass a fresh copy per call.
	fmt.Println("=== Approach 1: DP (Longest Non-Overlapping Chain) ===")
	fmt.Println(dpLongestChain(copyIntervals(ex1))) // expected 1
	fmt.Println(dpLongestChain(copyIntervals(ex2))) // expected 2
	fmt.Println(dpLongestChain(copyIntervals(ex3))) // expected 0

	fmt.Println("=== Approach 2: Greedy by Earliest End Time (Optimal) ===")
	fmt.Println(greedyEarliestEnd(copyIntervals(ex1))) // expected 1
	fmt.Println(greedyEarliestEnd(copyIntervals(ex2))) // expected 2
	fmt.Println(greedyEarliestEnd(copyIntervals(ex3))) // expected 0

	fmt.Println("=== Approach 3: Greedy by Earliest Start (Keep Shorter) ===")
	fmt.Println(greedyEarliestStart(copyIntervals(ex1))) // expected 1
	fmt.Println(greedyEarliestStart(copyIntervals(ex2))) // expected 2
	fmt.Println(greedyEarliestStart(copyIntervals(ex3))) // expected 0
}

// copyIntervals returns a deep copy so an in-place sort in one approach doesn't
// disturb the shared example data used by the others.
func copyIntervals(src [][]int) [][]int {
	out := make([][]int, len(src))
	for i, iv := range src {
		out[i] = []int{iv[0], iv[1]} // copy the two endpoints
	}
	return out
}
