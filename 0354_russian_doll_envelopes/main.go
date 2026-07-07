package main

import (
	"fmt"
	"sort"
)

// Russian Doll Envelopes: given envelopes [w, h], one envelope fits inside
// another only if BOTH its width and height are strictly smaller. Find the
// longest chain of envelopes that can be nested (Russian-doll style).
//
// This is a 2-D generalization of Longest Increasing Subsequence (LIS): sort to
// fix one dimension, then run LIS on the other.
//
// Two approaches:
//   1. dpQuadratic — sort by width, O(n^2) LIS on height.
//   2. lisBinarySearch (optimal) — sort by width asc, height DESC on ties, then
//      patience-sorting LIS on heights in O(n log n).

// ── Approach 1: Sort + O(n^2) DP (LIS on height) ─────────────────────────────
//
// dpQuadratic sorts envelopes by width (then height) and runs the classic
// O(n^2) longest-increasing-subsequence DP on the sorted array, requiring BOTH
// dimensions to strictly increase along the chain.
//
// Intuition:
//
//	If we sort by width ascending, then any valid nesting chain appears as a
//	left-to-right subsequence. But equal widths cannot nest, so a pair only
//	extends a chain when width[j] < width[i] AND height[j] < height[i]. That is
//	exactly LIS with a 2-D strict-increase test.
//
// Algorithm:
//  1. Sort envelopes by width asc, and by height asc on width ties.
//  2. dp[i] = longest chain ending with envelope i (>=1).
//  3. For each i, for each j < i: if both dims of j are strictly < those of i,
//     dp[i] = max(dp[i], dp[j]+1).
//  4. Answer = max over dp.
//
// Time:  O(n^2) — the double loop.
// Space: O(n) for dp.
func dpQuadratic(envelopes [][]int) int {
	n := len(envelopes)
	if n == 0 {
		return 0
	}
	// Sort by width asc; on equal width, by height asc (order irrelevant here
	// because the strict-both test handles equal widths).
	sort.Slice(envelopes, func(i, j int) bool {
		if envelopes[i][0] != envelopes[j][0] {
			return envelopes[i][0] < envelopes[j][0]
		}
		return envelopes[i][1] < envelopes[j][1]
	})

	dp := make([]int, n) // dp[i] = longest nesting chain ending at i
	best := 0
	for i := 0; i < n; i++ {
		dp[i] = 1 // an envelope alone is a chain of length 1
		for j := 0; j < i; j++ {
			// j can sit strictly inside i only if BOTH dims are smaller.
			if envelopes[j][0] < envelopes[i][0] && envelopes[j][1] < envelopes[i][1] {
				if dp[j]+1 > dp[i] {
					dp[i] = dp[j] + 1
				}
			}
		}
		if dp[i] > best {
			best = dp[i]
		}
	}
	return best
}

// ── Approach 2: Sort (w asc, h desc) + O(n log n) LIS on height (Optimal) ─────
//
// lisBinarySearch fixes the width dimension by sorting width ascending and
// height DESCENDING within equal widths, so a strictly-increasing subsequence of
// the height sequence corresponds exactly to a valid nesting chain. It then runs
// patience-sorting LIS (binary search) on the heights.
//
// Intuition:
//
//	After sorting by width asc, if we ran LIS on height we would wrongly allow two
//	envelopes of the SAME width (heights increasing) to "nest". Sorting height
//	DESC on equal widths makes those heights non-increasing, so LIS (strictly
//	increasing) can never pick two of them — enforcing the strict-width rule for
//	free. Then the answer is just the longest strictly-increasing subsequence of
//	the height array.
//
// Algorithm:
//  1. Sort by width asc; on tie, height DESC.
//  2. tails[] = smallest possible tail height of an increasing subsequence of
//     each length (patience piles).
//  3. For each height h: binary-search the first tail >= h; replace it (or append
//     if h is larger than all tails).
//  4. Answer = len(tails).
//
// Time:  O(n log n) — sort plus a binary search per element.
// Space: O(n) for tails.
func lisBinarySearch(envelopes [][]int) int {
	n := len(envelopes)
	if n == 0 {
		return 0
	}
	// width asc; on equal width, height DESCENDING (key trick).
	sort.Slice(envelopes, func(i, j int) bool {
		if envelopes[i][0] != envelopes[j][0] {
			return envelopes[i][0] < envelopes[j][0]
		}
		return envelopes[i][1] > envelopes[j][1]
	})

	tails := []int{} // tails[k] = min tail height of an incr. subseq of length k+1
	for _, e := range envelopes {
		h := e[1]
		// First index in tails with value >= h (strictly-increasing LIS).
		lo := sort.Search(len(tails), func(i int) bool { return tails[i] >= h })
		if lo == len(tails) {
			tails = append(tails, h) // h extends the longest chain
		} else {
			tails[lo] = h // h can start/continue a shorter chain more tightly
		}
	}
	return len(tails)
}

func main() {
	fmt.Println("=== Approach 1: Sort + O(n^2) DP ===")
	fmt.Println(dpQuadratic([][]int{{5, 4}, {6, 4}, {6, 7}, {2, 3}}))         // expected 3
	fmt.Println(dpQuadratic([][]int{{1, 1}, {1, 1}, {1, 1}}))                 // expected 1
	fmt.Println(dpQuadratic([][]int{{4, 5}, {4, 6}, {6, 7}, {2, 3}, {1, 1}})) // expected 4

	fmt.Println("=== Approach 2: Sort (w asc, h desc) + O(n log n) LIS (Optimal) ===")
	fmt.Println(lisBinarySearch([][]int{{5, 4}, {6, 4}, {6, 7}, {2, 3}}))         // expected 3
	fmt.Println(lisBinarySearch([][]int{{1, 1}, {1, 1}, {1, 1}}))                 // expected 1
	fmt.Println(lisBinarySearch([][]int{{4, 5}, {4, 6}, {6, 7}, {2, 3}, {1, 1}})) // expected 4
}
