package main

import (
	"fmt"
	"strings"
)

// count returns how many zeros and ones a binary string contains.
func count(s string) (zeros, ones int) {
	zeros = strings.Count(s, "0") // number of '0' characters
	ones = len(s) - zeros         // the rest are '1' characters
	return
}

// ── Approach 1: Brute Force (Enumerate Every Subset) ─────────────────────────
//
// bruteForce solves Ones and Zeroes by trying, via recursion, to take or skip
// every string and keeping the largest feasible subset.
//
// Intuition:
//
//	Each string is a binary include/exclude decision. Recurse over the strings;
//	if the current string still fits the remaining 0- and 1-budget, branch into
//	taking it (spend its zeros/ones, +1 to the count) or skipping it, and keep
//	the maximum. This is the definitional search behind the DP.
//
// Algorithm:
//  1. rec(i, m, n): if i == len(strs), return 0.
//  2. best = rec(i+1, m, n)                          // skip string i
//  3. if it fits (z0 <= m, o1 <= n):
//     best = max(best, 1 + rec(i+1, m-z0, n-o1))     // take string i
//  4. Return best.
//
// Time:  O(2^k) — k strings, take/skip each (exponential; baseline only).
// Space: O(k) — recursion depth.
func bruteForce(strs []string, m int, n int) int {
	// Pre-count zeros/ones once so recursion is cheap.
	zeros := make([]int, len(strs))
	ones := make([]int, len(strs))
	for i, s := range strs {
		zeros[i], ones[i] = count(s)
	}

	var rec func(i, remM, remN int) int
	rec = func(i, remM, remN int) int {
		if i == len(strs) {
			return 0 // no strings left to consider
		}
		best := rec(i+1, remM, remN) // Option A: skip string i
		// Option B: take string i, only if its cost fits the remaining budget.
		if zeros[i] <= remM && ones[i] <= remN {
			take := 1 + rec(i+1, remM-zeros[i], remN-ones[i])
			if take > best {
				best = take
			}
		}
		return best
	}
	return rec(0, m, n)
}

// ── Approach 2: 3D DP (Item × m × n) ─────────────────────────────────────────
//
// dp3D solves Ones and Zeroes as a 0/1 knapsack with TWO capacities (zeros and
// ones), materialising the full item dimension for clarity.
//
// Intuition:
//
//	Classic knapsack: dp[i][j][k] = largest subset using the first i strings with
//	at most j zeros and k ones. Each string is one item whose "weight" is a pair
//	(zeros, ones) and whose "value" is 1. Take-or-skip gives the transition.
//
// Algorithm:
//  1. dp[0][*][*] = 0 (no strings ⇒ empty subset).
//  2. For each string i (1-indexed) with cost (z, o):
//     for j in 0..m, k in 0..n:
//     dp[i][j][k] = dp[i-1][j][k]                        // skip
//     if j >= z and k >= o:
//     dp[i][j][k] = max(dp[i][j][k], dp[i-1][j-z][k-o]+1) // take
//  3. Answer: dp[len][m][n].
//
// Time:  O(k · m · n) — fill each of the (k · m · n) cells in O(1).
// Space: O(k · m · n) — the full 3D table (before the rolling optimisation).
func dp3D(strs []string, m int, n int) int {
	k := len(strs)
	// dp[i][j][k]: use first i strings, budget j zeros and k ones.
	dp := make([][][]int, k+1)
	for i := range dp {
		dp[i] = make([][]int, m+1)
		for j := range dp[i] {
			dp[i][j] = make([]int, n+1)
		}
	}
	for i := 1; i <= k; i++ {
		z, o := count(strs[i-1]) // cost of the i-th string (1-indexed)
		for j := 0; j <= m; j++ {
			for l := 0; l <= n; l++ {
				dp[i][j][l] = dp[i-1][j][l] // skip string i
				// Take string i if both budgets allow.
				if j >= z && l >= o {
					cand := dp[i-1][j-z][l-o] + 1
					if cand > dp[i][j][l] {
						dp[i][j][l] = cand
					}
				}
			}
		}
	}
	return dp[k][m][n]
}

// ── Approach 3: 2D Rolling DP (Optimal) ──────────────────────────────────────
//
// dp2DRolling solves Ones and Zeroes by collapsing the item dimension: a single
// (m+1)×(n+1) table updated in place, iterating the two capacities DOWNWARD so
// each string is used at most once (0/1 knapsack, not unbounded).
//
// Intuition:
//
//	dp[j][k] only ever depends on the same-or-smaller (j,k) from the PREVIOUS
//	item. If we sweep j and k from high to low while folding in string i, the
//	cells dp[j-z][k-o] we read still hold the "without string i" value — exactly
//	what 0/1 knapsack needs. Going low→high would let one string be counted
//	multiple times (that would be the unbounded variant).
//
// Algorithm:
//  1. dp is (m+1)×(n+1) of zeros.
//  2. For each string with cost (z, o):
//     for j from m down to z, for l from n down to o:
//     dp[j][l] = max(dp[j][l], dp[j-z][l-o] + 1)
//  3. Answer: dp[m][n].
//
// Time:  O(k · m · n) — same work, one string at a time.
// Space: O(m · n) — a single 2D table reused across all strings.
func dp2DRolling(strs []string, m int, n int) int {
	// dp[j][k]: best subset size achievable with budget j zeros and k ones,
	// considering the strings processed so far.
	dp := make([][]int, m+1)
	for j := range dp {
		dp[j] = make([]int, n+1)
	}
	for _, s := range strs {
		z, o := count(s) // cost of this string
		// Iterate DOWNWARD so dp[j-z][l-o] is still the "before this string"
		// value — this is what makes each string usable at most once.
		for j := m; j >= z; j-- {
			for l := n; l >= o; l-- {
				cand := dp[j-z][l-o] + 1 // take this string
				if cand > dp[j][l] {
					dp[j][l] = cand
				}
			}
		}
	}
	return dp[m][n]
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Printf("[10 0001 111001 1 0], m=5,n=3 => %d  (expected 4)\n",
		bruteForce([]string{"10", "0001", "111001", "1", "0"}, 5, 3))
	fmt.Printf("[10 0 1], m=1,n=1             => %d  (expected 2)\n",
		bruteForce([]string{"10", "0", "1"}, 1, 1))

	fmt.Println("=== Approach 2: 3D DP (Item × m × n) ===")
	fmt.Printf("[10 0001 111001 1 0], m=5,n=3 => %d  (expected 4)\n",
		dp3D([]string{"10", "0001", "111001", "1", "0"}, 5, 3))
	fmt.Printf("[10 0 1], m=1,n=1             => %d  (expected 2)\n",
		dp3D([]string{"10", "0", "1"}, 1, 1))

	fmt.Println("=== Approach 3: 2D Rolling DP (Optimal) ===")
	fmt.Printf("[10 0001 111001 1 0], m=5,n=3 => %d  (expected 4)\n",
		dp2DRolling([]string{"10", "0001", "111001", "1", "0"}, 5, 3))
	fmt.Printf("[10 0 1], m=1,n=1             => %d  (expected 2)\n",
		dp2DRolling([]string{"10", "0", "1"}, 1, 1))
}
