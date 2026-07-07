package main

import "fmt"

// ── Approach 1: Top-Down DP (Memoized Minimax) ───────────────────────────────
//
// dpTopDown computes the minimum amount of money that GUARANTEES a win over the
// worst-case hidden number in [1, n]. Guessing x costs x if wrong; we then pay
// for the worse of the two remaining subranges. We minimise this maximum.
//
// Intuition:
//
//	Define cost(lo, hi) = min money to guarantee finding any number in [lo, hi].
//	If we guess x, the adversary picks the branch that costs us more, so this
//	guess costs x + max(cost(lo, x-1), cost(x+1, hi)). We try every x and take
//	the cheapest — a classic minimax / interval DP. Base case: a range of size
//	≤ 1 costs 0 (we already know the answer).
//
// Algorithm:
//  1. memo[lo][hi] caches solved ranges.
//  2. solve(lo, hi): if lo >= hi return 0.
//  3. For x = lo..hi: candidate = x + max(solve(lo, x-1), solve(x+1, hi)).
//     Keep the minimum candidate. Store and return it.
//
// Time:  O(n^3) — O(n^2) distinct (lo, hi) ranges, each scanning O(n) guesses.
// Space: O(n^2) — the memo table plus recursion depth O(n).
func dpTopDown(n int) int {
	// memo[lo][hi]; -1 marks "not computed yet". Size n+2 to allow x+1 up to n+1.
	memo := make([][]int, n+2)
	for i := range memo {
		memo[i] = make([]int, n+2)
		for j := range memo[i] {
			memo[i][j] = -1
		}
	}

	var solve func(lo, hi int) int
	solve = func(lo, hi int) int {
		if lo >= hi { // 0 or 1 candidate → no cost to be certain
			return 0
		}
		if memo[lo][hi] != -1 { // reuse a solved range
			return memo[lo][hi]
		}
		best := 1 << 30 // +infinity sentinel
		for x := lo; x <= hi; x++ {
			left := solve(lo, x-1)  // worst cost if the pick is below x
			right := solve(x+1, hi) // worst cost if the pick is above x
			worse := left           // adversary forces the more expensive side
			if right > worse {
				worse = right
			}
			cost := x + worse // pay x for the wrong guess, then the worse branch
			if cost < best {  // keep the cheapest guarantee
				best = cost
			}
		}
		memo[lo][hi] = best
		return best
	}
	return solve(1, n)
}

// ── Approach 2: Bottom-Up Interval DP (Optimal) ──────────────────────────────
//
// dpBottomUp fills the same table iteratively by increasing interval length, so
// every sub-interval a guess depends on is already solved.
//
// Intuition:
//
//	dp[lo][hi] = min over x in [lo,hi] of x + max(dp[lo][x-1], dp[x+1][hi]).
//	Ranges of length L only depend on ranges of length < L, so processing by
//	ascending length (or descending lo / ascending hi) removes recursion.
//
// Algorithm:
//  1. dp is (n+2) x (n+2), all zero (empty/singleton ranges cost 0).
//  2. For lo from n-1 down to 1, for hi from lo+1 to n:
//     dp[lo][hi] = min_{x in [lo,hi]} ( x + max(dp[lo][x-1], dp[x+1][hi]) ).
//  3. Answer is dp[1][n].
//
// Time:  O(n^3). Space: O(n^2).
func dpBottomUp(n int) int {
	// dp[lo][hi]; indices 0..n+1 so dp[x+1][hi] and dp[lo][x-1] stay in bounds.
	dp := make([][]int, n+2)
	for i := range dp {
		dp[i] = make([]int, n+2)
	}

	// Shorter ranges first: lo decreasing guarantees dp[lo][x-1] (smaller hi via
	// same lo won't help) — we rely on dp[x+1][hi] having larger lo (already done)
	// and dp[lo][x-1] having smaller hi (already done in this lo's inner loop).
	for lo := n - 1; lo >= 1; lo-- {
		for hi := lo + 1; hi <= n; hi++ {
			best := 1 << 30
			for x := lo; x <= hi; x++ {
				left := dp[lo][x-1]  // already computed (smaller hi)
				right := dp[x+1][hi] // already computed (larger lo)
				worse := left
				if right > worse {
					worse = right
				}
				cost := x + worse
				if cost < best {
					best = cost
				}
			}
			dp[lo][hi] = best
		}
	}
	return dp[1][n]
}

func main() {
	fmt.Println("=== Approach 1: Top-Down DP (Memoized Minimax) ===")
	fmt.Println(dpTopDown(10)) // expected 16
	fmt.Println(dpTopDown(1))  // expected 0
	fmt.Println(dpTopDown(2))  // expected 1

	fmt.Println("=== Approach 2: Bottom-Up Interval DP ===")
	fmt.Println(dpBottomUp(10)) // expected 16
	fmt.Println(dpBottomUp(1))  // expected 0
	fmt.Println(dpBottomUp(2))  // expected 1
}
