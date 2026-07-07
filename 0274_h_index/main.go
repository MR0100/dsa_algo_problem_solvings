package main

import (
	"fmt"
	"sort"
)

// H-Index (LeetCode #274)
//
// Given an array citations where citations[i] is the number of citations a
// researcher received for their i-th paper, return the researcher's h-index.
//
// The h-index is the maximum value h such that the researcher has published at
// least h papers that have each been cited at least h times.

// ── Approach 1: Sort Descending (Optimal-ish) ────────────────────────────────
//
// sortDescending sorts citations high-to-low and finds the largest position i
// (1-indexed) where the i-th paper still has at least i citations.
//
// Intuition:
//
//	Sort papers by citations, most-cited first. Walk down the list: after
//	seeing i papers, all i of them have at least citations[i-1] citations
//	(sorted desc). So h can be at least i as long as the i-th paper has >= i
//	citations. The h-index is the largest such i. The moment a paper's
//	citation count drops below its 1-based rank, no larger h is possible.
//
// Algorithm:
//  1. Sort citations in descending order.
//  2. For i from 0..n-1: if citations[i] >= i+1, h = i+1; else break.
//  3. Return h.
//
// Time:  O(n log n) — the sort.
// Space: O(1) — in-place sort (ignoring sort's internal use).
func sortDescending(citations []int) int {
	// Sort a copy would be cleaner, but sorting in place is standard here.
	sort.Sort(sort.Reverse(sort.IntSlice(citations)))
	h := 0
	for i := 0; i < len(citations); i++ {
		// citations[i] is the (i+1)-th largest; if it still has >= i+1 cites,
		// we have i+1 papers each cited >= i+1 times.
		if citations[i] >= i+1 {
			h = i + 1
		} else {
			break // counts only decrease from here — no larger h possible
		}
	}
	return h
}

// ── Approach 2: Counting Sort / Buckets (Optimal O(n)) ───────────────────────
//
// countingBuckets avoids the comparison sort by bucketing papers by citation
// count, capping counts above n at n (since h can never exceed n).
//
// Intuition:
//
//	The h-index is at most n (you can't have more papers than you published).
//	So bucket papers by citation count, lumping everything >= n into bucket n.
//	Then scan buckets from high to low, accumulating how many papers have AT
//	LEAST that many citations. The first citation level h where that running
//	count >= h is the h-index.
//
// Algorithm:
//  1. buckets[0..n], buckets[min(c, n)]++ for each citation c.
//  2. total = 0; for h from n down to 0: total += buckets[h];
//     if total >= h, return h.
//  3. (Loop always returns; h = 0 is the floor.)
//
// Time:  O(n) — one pass to bucket, one pass over n+1 buckets.
// Space: O(n) — the bucket array.
func countingBuckets(citations []int) int {
	n := len(citations)
	buckets := make([]int, n+1) // index = citation count, capped at n
	for _, c := range citations {
		if c >= n {
			buckets[n]++ // everything >= n lands in the top bucket
		} else {
			buckets[c]++
		}
	}
	total := 0 // papers with AT LEAST h citations, accumulated high→low
	for h := n; h >= 0; h-- {
		total += buckets[h]
		if total >= h {
			// h papers each have >= h citations → h is achievable, and since we
			// scan from the top, this is the maximum such h.
			return h
		}
	}
	return 0
}

// ── Approach 3: Binary Search on the Answer ──────────────────────────────────
//
// binarySearchAnswer binary-searches the candidate h in [0, n], using a
// feasibility check "are there at least h papers with >= h citations?".
//
// Intuition:
//
//	The predicate "at least h papers have >= h citations" is monotone in h:
//	if it holds for h, it holds for all smaller h. That monotonicity lets us
//	binary-search the largest feasible h instead of scanning linearly.
//
// Algorithm:
//  1. lo = 0, hi = n. While lo < hi (bias high): mid = (lo+hi+1)/2.
//  2. Count papers with >= mid citations; if count >= mid, lo = mid, else hi = mid-1.
//  3. Return lo.
//
// Time:  O(n log n) — log n iterations, each an O(n) count.
// Space: O(1).
func binarySearchAnswer(citations []int) int {
	n := len(citations)
	lo, hi := 0, n
	for lo < hi {
		mid := (lo + hi + 1) / 2 // bias upward to make progress toward larger h
		count := 0
		for _, c := range citations {
			if c >= mid {
				count++ // this paper qualifies for a candidate h of `mid`
			}
		}
		if count >= mid {
			lo = mid // feasible: at least mid papers with >= mid cites → try larger
		} else {
			hi = mid - 1 // infeasible → h must be smaller
		}
	}
	return lo
}

func main() {
	ex1 := []int{3, 0, 6, 1, 5} // expected 3
	ex2 := []int{1, 3, 1}       // expected 1

	fmt.Println("=== Approach 1: Sort Descending ===")
	fmt.Println(sortDescending(append([]int(nil), ex1...))) // expected 3
	fmt.Println(sortDescending(append([]int(nil), ex2...))) // expected 1

	fmt.Println("=== Approach 2: Counting Buckets ===")
	fmt.Println(countingBuckets(ex1)) // expected 3
	fmt.Println(countingBuckets(ex2)) // expected 1

	fmt.Println("=== Approach 3: Binary Search on the Answer ===")
	fmt.Println(binarySearchAnswer(ex1)) // expected 3
	fmt.Println(binarySearchAnswer(ex2)) // expected 1
}
