package main

import "fmt"

// Split nums into k non-empty CONTIGUOUS subarrays so that the LARGEST subarray
// sum is as SMALL as possible; return that minimized largest sum.

// ── Approach 1: Binary Search on the Answer (Optimal) ─────────────────────────
//
// binarySearch solves Split Array Largest Sum by binary-searching the value of
// the largest allowed subarray sum and greedily checking feasibility.
//
// Intuition:
//
//	The answer (the largest subarray sum) is a number in [max(nums), sum(nums)]:
//	it can't be smaller than the biggest single element (that element sits in some
//	part alone at minimum), and it can't exceed the total (k=1 puts everything in
//	one part). "Can we split so no part exceeds cap?" is MONOTONE in cap: a larger
//	cap is always at least as easy. So binary-search the smallest feasible cap.
//	Feasibility check: greedily grow the current part; whenever adding the next
//	number would exceed cap, cut a new part. Count the parts; feasible iff parts
//	<= k.
//
// Algorithm:
//  1. lo = max(nums), hi = sum(nums).
//  2. While lo < hi: mid = (lo+hi)/2; if canSplit(mid) shrink hi=mid else lo=mid+1.
//  3. Return lo — the least cap that needs at most k parts.
//
// Time:  O(n * log(sum - max)) — each feasibility scan is O(n), search is log(range).
// Space: O(1).
func binarySearch(nums []int, k int) int {
	lo, hi := 0, 0
	for _, v := range nums {
		if v > lo {
			lo = v // lower bound: no part can be smaller than the largest element
		}
		hi += v // upper bound: everything in one part
	}

	// canSplit reports whether nums can be cut into <= k parts each summing <= cap.
	canSplit := func(cap int) bool {
		parts := 1   // we always have at least one part
		current := 0 // running sum of the part we are currently filling
		for _, v := range nums {
			if current+v > cap {
				// v doesn't fit in the current part → start a new part with v.
				parts++
				current = v
				if parts > k {
					return false // needed more than k parts ⇒ cap too small
				}
			} else {
				current += v // v fits, keep filling this part
			}
		}
		return true
	}

	// Standard "find smallest feasible value" binary search.
	for lo < hi {
		mid := lo + (hi-lo)/2 // candidate cap (avoids overflow)
		if canSplit(mid) {
			hi = mid // mid works; try to do even smaller
		} else {
			lo = mid + 1 // mid too small; need a bigger cap
		}
	}
	return lo // smallest cap that is feasible = minimized largest sum
}

// ── Approach 2: Dynamic Programming (Partition DP) ────────────────────────────
//
// dpBottomUp solves Split Array Largest Sum by computing the best "minimized
// largest sum" for prefixes split into a given number of parts.
//
// Intuition:
//
//	Let dp[i][j] = the minimized largest-subarray-sum when splitting the first i
//	elements into exactly j parts. To fill dp[i][j], let the LAST part cover
//	elements (p..i-1] for some p; that last part's sum is prefix[i]-prefix[p], and
//	the first p elements are split into j-1 parts optimally = dp[p][j-1]. The cost
//	of this choice is the MAX of those two (the largest part overall). Minimize
//	over all valid split points p.
//	    dp[i][j] = min over p of max(dp[p][j-1], prefix[i]-prefix[p]).
//
// Algorithm:
//  1. prefix[i] = sum of first i elements.
//  2. dp[0][0] = 0; everything else +inf.
//  3. For i=1..n, j=1..min(i,k): dp[i][j] = min over p in [j-1, i-1] of
//     max(dp[p][j-1], prefix[i]-prefix[p]).
//  4. Answer = dp[n][k].
//
// Time:  O(n^2 * k) — for each (i, j) we scan up to n split points.
// Space: O(n * k) — the DP table.
func dpBottomUp(nums []int, k int) int {
	n := len(nums)
	const inf = int(1e18) // sentinel for "not reachable"

	// prefix[i] = nums[0] + ... + nums[i-1].
	prefix := make([]int, n+1)
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i] + nums[i]
	}

	// dp[i][j] = min largest-part sum splitting first i elems into j parts.
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, k+1)
		for j := range dp[i] {
			dp[i][j] = inf // start unreachable
		}
	}
	dp[0][0] = 0 // zero elements, zero parts, zero cost

	for i := 1; i <= n; i++ {
		// Can't have more parts than elements, nor more than k.
		for j := 1; j <= k && j <= i; j++ {
			// Try every split point p: last part is (p..i-1], value prefix[i]-prefix[p].
			for p := j - 1; p < i; p++ {
				if dp[p][j-1] == inf {
					continue // first p elements can't be split into j-1 parts
				}
				lastPart := prefix[i] - prefix[p] // sum of the last (j-th) part
				candidate := dp[p][j-1]           // largest part among the first j-1
				if lastPart > candidate {
					candidate = lastPart // the overall largest part is the max of the two
				}
				if candidate < dp[i][j] {
					dp[i][j] = candidate // keep the split that minimizes the largest part
				}
			}
		}
	}
	return dp[n][k]
}

func main() {
	fmt.Println("=== Approach 1: Binary Search on the Answer (Optimal) ===")
	fmt.Printf("nums=[7,2,5,10,8], k=2 got=%d  expected 18\n", binarySearch([]int{7, 2, 5, 10, 8}, 2))
	fmt.Printf("nums=[1,2,3,4,5], k=2 got=%d  expected 9\n", binarySearch([]int{1, 2, 3, 4, 5}, 2))
	fmt.Printf("nums=[1,4,4],     k=3 got=%d  expected 4\n", binarySearch([]int{1, 4, 4}, 3))

	fmt.Println("=== Approach 2: Dynamic Programming (Partition DP) ===")
	fmt.Printf("nums=[7,2,5,10,8], k=2 got=%d  expected 18\n", dpBottomUp([]int{7, 2, 5, 10, 8}, 2))
	fmt.Printf("nums=[1,2,3,4,5], k=2 got=%d  expected 9\n", dpBottomUp([]int{1, 2, 3, 4, 5}, 2))
	fmt.Printf("nums=[1,4,4],     k=3 got=%d  expected 4\n", dpBottomUp([]int{1, 4, 4}, 3))
}
