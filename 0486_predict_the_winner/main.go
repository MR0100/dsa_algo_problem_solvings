package main

import "fmt"

// ── Approach 1: Brute Force Recursion (Minimax on score difference) ───────────
//
// bruteForce solves Predict the Winner by simulating the whole game with plain
// recursion, returning the best score DIFFERENCE (current player − opponent)
// the mover can force from the sub-array nums[lo..hi].
//
// Intuition:
//
//	Both players play optimally, so this is a two-player zero-sum game. Instead
//	of tracking two separate totals, track a single number: the net advantage
//	of whoever is about to move. If the mover takes nums[lo], they gain
//	nums[lo] now and then FACE the opponent playing optimally on [lo+1..hi];
//	from the mover's perspective the opponent's best difference is subtracted.
//	The mover picks whichever end maximises (their pick − opponent's best diff).
//	Player 1 wins iff the difference over the whole array is ≥ 0 (ties go to P1).
//
// Algorithm:
//  1. score(lo, hi): if lo == hi, only one number is left → return nums[lo].
//  2. Otherwise return max(
//     nums[lo] − score(lo+1, hi),   // take the left end
//     nums[hi] − score(lo, hi-1)).  // take the right end
//  3. PredictTheWinner = score(0, n-1) >= 0.
//
// Time:  O(2^n) — each call branches into two, no memoisation.
// Space: O(n) — recursion depth equals the number of picks.
func bruteForce(nums []int) bool {
	// score returns the best achievable (mover − opponent) difference on nums[lo..hi].
	var score func(lo, hi int) int
	score = func(lo, hi int) int {
		if lo == hi {
			return nums[lo] // last remaining number goes straight to the mover
		}
		// Take the left end: gain nums[lo], then subtract the opponent's best
		// difference on the smaller range (roles flip → subtract).
		takeLeft := nums[lo] - score(lo+1, hi)
		// Take the right end symmetrically.
		takeRight := nums[hi] - score(lo, hi-1)
		if takeLeft > takeRight {
			return takeLeft // prefer the end that yields the larger net advantage
		}
		return takeRight
	}
	return score(0, len(nums)-1) >= 0 // ≥ 0 ⇒ Player 1 at least ties ⇒ wins
}

// ── Approach 2: Top-Down DP (Memoised Minimax) ───────────────────────────────
//
// dpTopDown solves Predict the Winner with the same minimax recursion as brute
// force, but caches score(lo, hi) so each (lo, hi) pair is computed once.
//
// Intuition:
//
//	The recursion only ever depends on the pair (lo, hi), and there are just
//	n² such pairs. The exponential blow-up came from recomputing the same
//	sub-array difference over and over. Memoising collapses the 2^n tree into
//	an n² table lookup.
//
// Algorithm:
//  1. memo[lo][hi] holds score(lo, hi); seen[lo][hi] marks it as computed.
//  2. score(lo, hi): return the cache if present; else compute
//     max(nums[lo] − score(lo+1,hi), nums[hi] − score(lo,hi-1)), store, return.
//  3. Answer is score(0, n-1) >= 0.
//
// Time:  O(n²) — n² distinct states, O(1) work each.
// Space: O(n²) memo table + O(n) recursion stack.
func dpTopDown(nums []int) bool {
	n := len(nums)
	memo := make([][]int, n)  // memo[lo][hi] = best difference on nums[lo..hi]
	seen := make([][]bool, n) // seen[lo][hi] = has this state been computed?
	for i := range memo {
		memo[i] = make([]int, n)
		seen[i] = make([]bool, n)
	}
	var score func(lo, hi int) int
	score = func(lo, hi int) int {
		if lo == hi {
			return nums[lo] // base case: single number
		}
		if seen[lo][hi] {
			return memo[lo][hi] // already solved this sub-array
		}
		takeLeft := nums[lo] - score(lo+1, hi)  // take left end
		takeRight := nums[hi] - score(lo, hi-1) // take right end
		best := takeLeft
		if takeRight > best {
			best = takeRight
		}
		seen[lo][hi] = true // memoise before returning
		memo[lo][hi] = best
		return best
	}
	return score(0, n-1) >= 0
}

// ── Approach 3: Bottom-Up Interval DP (2D table) ─────────────────────────────
//
// dpBottomUp solves Predict the Winner by filling a 2D interval table by
// increasing sub-array length, so no recursion is needed.
//
// Intuition:
//
//	dp[lo][hi] is the best difference (mover − opponent) on nums[lo..hi]. A
//	length-1 interval has dp[i][i] = nums[i]. Every longer interval depends
//	only on strictly shorter ones (lo+1..hi and lo..hi-1), so if we process
//	intervals from shortest to longest, the dependencies are always ready.
//	This is the classic interval-DP fill order.
//
// Algorithm:
//  1. Initialise dp[i][i] = nums[i] for all i.
//  2. For length = 2..n, for each lo with hi = lo+length-1:
//     dp[lo][hi] = max(nums[lo] − dp[lo+1][hi], nums[hi] − dp[lo][hi-1]).
//  3. Return dp[0][n-1] >= 0.
//
// Time:  O(n²) — two nested loops over interval endpoints.
// Space: O(n²) — the full triangular table.
func dpBottomUp(nums []int) bool {
	n := len(nums)
	dp := make([][]int, n) // dp[lo][hi] = best difference on nums[lo..hi]
	for i := range dp {
		dp[i] = make([]int, n)
		dp[i][i] = nums[i] // base case: a single number is a pure gain for the mover
	}
	// Grow the interval length; every state below depends only on shorter ones.
	for length := 2; length <= n; length++ {
		for lo := 0; lo+length-1 < n; lo++ {
			hi := lo + length - 1
			takeLeft := nums[lo] - dp[lo+1][hi]  // take the left end
			takeRight := nums[hi] - dp[lo][hi-1] // take the right end
			if takeLeft > takeRight {
				dp[lo][hi] = takeLeft
			} else {
				dp[lo][hi] = takeRight
			}
		}
	}
	return dp[0][n-1] >= 0
}

// ── Approach 4: Space-Optimised DP (1D rolling array) (Optimal) ───────────────
//
// dpOneDim solves Predict the Winner by collapsing the 2D interval table into a
// single 1D array, since each cell only needs its left neighbour and its value
// from the previous (shorter) pass.
//
// Intuition:
//
//	In the bottom-up fill, dp[lo][hi] reads dp[lo+1][hi] (from the previous
//	shorter length, same column hi) and dp[lo][hi-1] (already updated this
//	pass, one column left). If we reuse a 1D array `dp` indexed by hi and
//	iterate lo downward for each hi, then before we overwrite dp[hi] it still
//	holds the "lo+1" value, and dp[hi-1] holds the freshly computed "same lo"
//	value. That is exactly the two inputs we need — so one row suffices.
//
// Algorithm:
//  1. dp[i] = nums[i] initially (all length-1 intervals: dp[lo][lo]).
//  2. For lo from n-2 down to 0, for hi from lo+1 to n-1:
//     dp[hi] = max(nums[lo] − dp[hi], nums[hi] − dp[hi-1]).
//     (dp[hi] before write = old dp[lo+1][hi]; dp[hi-1] = new dp[lo][hi-1].)
//  3. Return dp[n-1] >= 0.
//
// Time:  O(n²) — same double loop.
// Space: O(n) — a single rolling row instead of the full table.
func dpOneDim(nums []int) bool {
	n := len(nums)
	dp := make([]int, n) // dp[hi] doubles as dp[lo][hi] for the current lo
	copy(dp, nums)       // length-1 intervals: dp[i] = nums[i]
	// Sweep lo from the second-to-last start downward so shorter intervals are ready.
	for lo := n - 2; lo >= 0; lo-- {
		for hi := lo + 1; hi < n; hi++ {
			// dp[hi]   still holds the value for start lo+1 (previous outer pass).
			// dp[hi-1] already holds the value for start lo   (this pass).
			takeLeft := nums[lo] - dp[hi]    // take left end: opponent solves [lo+1..hi]
			takeRight := nums[hi] - dp[hi-1] // take right end: opponent solves [lo..hi-1]
			if takeLeft > takeRight {
				dp[hi] = takeLeft
			} else {
				dp[hi] = takeRight
			}
		}
	}
	return dp[n-1] >= 0
}

func main() {
	fmt.Println("=== Approach 1: Brute Force Recursion ===")
	fmt.Println(bruteForce([]int{1, 5, 2}))      // expected false
	fmt.Println(bruteForce([]int{1, 5, 233, 7})) // expected true

	fmt.Println("=== Approach 2: Top-Down DP (Memoised) ===")
	fmt.Println(dpTopDown([]int{1, 5, 2}))      // expected false
	fmt.Println(dpTopDown([]int{1, 5, 233, 7})) // expected true

	fmt.Println("=== Approach 3: Bottom-Up Interval DP ===")
	fmt.Println(dpBottomUp([]int{1, 5, 2}))      // expected false
	fmt.Println(dpBottomUp([]int{1, 5, 233, 7})) // expected true

	fmt.Println("=== Approach 4: Space-Optimised 1D DP (Optimal) ===")
	fmt.Println(dpOneDim([]int{1, 5, 2}))      // expected false
	fmt.Println(dpOneDim([]int{1, 5, 233, 7})) // expected true
}
