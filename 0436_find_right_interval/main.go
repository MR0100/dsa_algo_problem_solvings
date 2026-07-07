package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Brute Force ──────────────────────────────────────────────────
//
// bruteForce solves Find Right Interval by, for each interval i, scanning every
// other interval j and keeping the one with the smallest start that is still
// >= end_i.
//
// Intuition:
//
//	The right interval of i is defined directly: among all j with start_j >=
//	end_i, pick the smallest start_j. Just try every j for every i and track
//	the best (smallest qualifying start, remembering its index).
//
// Algorithm:
//  1. For each i, set best start = +infinity and answer = -1.
//  2. For each j, if start_j >= end_i and start_j < best, update best and
//     record j as the answer.
//  3. Store the recorded index (or -1) in res[i].
//
// Time:  O(n^2) — every pair (i, j) is examined once.
// Space: O(1) extra beyond the output (the result slice itself is O(n)).
func bruteForce(intervals [][]int) []int {
	n := len(intervals)
	res := make([]int, n) // res[i] = index of the right interval of i, or -1
	for i := 0; i < n; i++ {
		endI := intervals[i][1] // the value every candidate start must reach
		best := -1              // index of the current best right interval
		bestStart := 1 << 62    // smallest qualifying start seen so far (huge sentinel)
		for j := 0; j < n; j++ {
			startJ := intervals[j][0]
			// j qualifies if its start is at least end_i; among qualifiers we
			// want the minimal start, so compare against bestStart.
			if startJ >= endI && startJ < bestStart {
				bestStart = startJ // tighten the minimum
				best = j           // remember which interval achieved it
			}
		}
		res[i] = best
	}
	return res
}

// ── Approach 2: Sort Starts + Binary Search ──────────────────────────────────
//
// binarySearchSorted solves Find Right Interval by building a sorted list of
// (start, originalIndex) pairs, then for each interval binary-searching the
// first start that is >= its end.
//
// Intuition:
//
//	"Smallest start that is >= end_i" is exactly a lower-bound query on the
//	multiset of starts. Sort the starts once; each query is then a binary
//	search. We must carry each start's ORIGINAL index because the answer is an
//	index into the input, not into the sorted array.
//
// Algorithm:
//  1. Build starts = [(start_i, i) for all i] and sort by start.
//  2. For each interval i, binary-search starts for the first pair whose start
//     >= end_i (a lower bound).
//  3. If found, res[i] = that pair's original index; otherwise -1.
//
// Time:  O(n log n) — one sort plus n binary searches of O(log n) each.
// Space: O(n) — the array of (start, index) pairs.
func binarySearchSorted(intervals [][]int) []int {
	n := len(intervals)
	// starts[k] = {start value, original interval index}; sorting reorders these.
	type pair struct{ start, idx int }
	starts := make([]pair, n)
	for i := 0; i < n; i++ {
		starts[i] = pair{intervals[i][0], i}
	}
	// Sort ascending by start so binary search can find the lower bound.
	sort.Slice(starts, func(a, b int) bool { return starts[a].start < starts[b].start })

	res := make([]int, n)
	for i := 0; i < n; i++ {
		endI := intervals[i][1]
		// sort.Search returns the smallest index pos in [0, n] such that the
		// predicate is true; here: first start >= endI (the lower bound).
		pos := sort.Search(n, func(k int) bool { return starts[k].start >= endI })
		if pos < n {
			res[i] = starts[pos].idx // map the sorted hit back to its original index
		} else {
			res[i] = -1 // no start reaches endI → no right interval
		}
	}
	return res
}

// ── Approach 3: Two Sorted Arrays + Two Pointers ─────────────────────────────
//
// twoPointers solves Find Right Interval by sorting intervals by start AND by
// end, then sweeping the two orders with a single moving pointer.
//
// Intuition:
//
//	Process intervals in increasing order of their END. As end grows, the first
//	start that catches up to it only moves forward, never back — so a single
//	pointer into the start-sorted list, advanced monotonically, answers every
//	query. This trades the per-query log factor for a linear sweep after the
//	sort.
//
// Algorithm:
//  1. byStart = indices sorted by interval start; byEnd = indices sorted by
//     interval end.
//  2. Walk byEnd in order. Maintain pointer p into byStart.
//  3. For the current interval (smallest remaining end), advance p while the
//     start it points to is < that end. The interval now at p is the answer.
//  4. Record it (or -1 if p ran off the end) at the current interval's index.
//
// Time:  O(n log n) — two sorts dominate; the sweep is O(n) amortised.
// Space: O(n) — two index arrays.
func twoPointers(intervals [][]int) []int {
	n := len(intervals)
	byStart := make([]int, n) // interval indices ordered by start ascending
	byEnd := make([]int, n)   // interval indices ordered by end ascending
	for i := range intervals {
		byStart[i] = i
		byEnd[i] = i
	}
	sort.Slice(byStart, func(a, b int) bool { return intervals[byStart[a]][0] < intervals[byStart[b]][0] })
	sort.Slice(byEnd, func(a, b int) bool { return intervals[byEnd[a]][1] < intervals[byEnd[b]][1] })

	res := make([]int, n)
	p := 0 // pointer into byStart; only ever moves forward across the whole sweep
	// Consider intervals from the smallest end to the largest.
	for _, i := range byEnd {
		endI := intervals[i][1]
		// Skip every start strictly smaller than endI; they can never be the
		// right interval for this end (or any larger end still to come).
		for p < n && intervals[byStart[p]][0] < endI {
			p++
		}
		if p < n {
			res[i] = byStart[p] // first start >= endI, in original-index terms
		} else {
			res[i] = -1 // all starts exhausted → no right interval
		}
	}
	return res
}

// equal reports whether two int slices match — used to label expected output.
func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func main() {
	ex1 := [][]int{{1, 2}}
	ex2 := [][]int{{3, 4}, {2, 3}, {1, 2}}
	ex3 := [][]int{{1, 4}, {2, 3}, {3, 4}} // extra: overlaps + a self/next match

	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(bruteForce(ex1), equal(bruteForce(ex1), []int{-1}))        // [-1] true
	fmt.Println(bruteForce(ex2), equal(bruteForce(ex2), []int{-1, 0, 1}))  // [-1 0 1] true
	fmt.Println(bruteForce(ex3), equal(bruteForce(ex3), []int{-1, 2, -1})) // [-1 2 -1] true

	fmt.Println("=== Approach 2: Sort Starts + Binary Search ===")
	fmt.Println(binarySearchSorted(ex1), equal(binarySearchSorted(ex1), []int{-1}))        // [-1] true
	fmt.Println(binarySearchSorted(ex2), equal(binarySearchSorted(ex2), []int{-1, 0, 1}))  // [-1 0 1] true
	fmt.Println(binarySearchSorted(ex3), equal(binarySearchSorted(ex3), []int{-1, 2, -1})) // [-1 2 -1] true

	fmt.Println("=== Approach 3: Two Sorted Arrays + Two Pointers ===")
	fmt.Println(twoPointers(ex1), equal(twoPointers(ex1), []int{-1}))        // [-1] true
	fmt.Println(twoPointers(ex2), equal(twoPointers(ex2), []int{-1, 0, 1}))  // [-1 0 1] true
	fmt.Println(twoPointers(ex3), equal(twoPointers(ex3), []int{-1, 2, -1})) // [-1 2 -1] true
}
