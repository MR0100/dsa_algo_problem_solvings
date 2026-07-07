package main

import "fmt"

// ── Approach 1: DP Bottom-Up ─────────────────────────────────────────────────
//
// dpBottomUp solves Integer Break by building the best product for each value
// from 2 up to n.
//
// Intuition:
//
//	To break n, pick a first piece j (1..n-1). The rest, n-j, can either be
//	left whole (product j*(n-j)) or itself be broken further (product
//	j*dp[n-j]). dp[i] is the max product achievable by breaking i into AT
//	LEAST two positive integers; take the best over all first cuts j.
//
// Algorithm:
//  1. dp[1] = 1 (used only as a factor).
//  2. For i = 2..n: for j = 1..i-1: dp[i] = max(dp[i], j*(i-j), j*dp[i-j]).
//  3. Return dp[n].
//
// Time:  O(n^2) — nested loops over i and j.
// Space: O(n) — the dp table.
func dpBottomUp(n int) int {
	dp := make([]int, n+1) // dp[i] = max product of breaking i into >=2 parts
	dp[1] = 1              // a leftover 1 contributes as a multiplicative factor
	for i := 2; i <= n; i++ {
		for j := 1; j < i; j++ {
			// j*(i-j): leave the remainder (i-j) whole.
			// j*dp[i-j]: break the remainder further.
			best := max(j*(i-j), j*dp[i-j])
			dp[i] = max(dp[i], best)
		}
	}
	return dp[n]
}

// ── Approach 2: Math — Break into 3s (Optimal) ───────────────────────────────
//
// mathThrees solves Integer Break using the fact that the optimal factors are
// 3s (with 2s to mop up the remainder).
//
// Intuition:
//
//	Splitting a factor into 3s maximises the product because 3 gives the best
//	product-per-unit among integers (3 > 2*... etc.; e.g. 6 → 3*3=9 beats
//	2*2*2=8). Never use a factor of 1. If the remainder mod 3 is 1, it is
//	better to turn one 3 into 2*2 (since 2*2=4 > 3*1=3).
//
// Algorithm:
//  1. n <= 3: the answer is n-1 (must break into >=2 parts: 2→1, 3→2).
//  2. Take out as many 3s as possible.
//  3. If remainder == 0: product is 3^(n/3).
//     If remainder == 1: use one fewer 3 and two 2s → 3^(n/3 - 1) * 4.
//     If remainder == 2: product is 3^(n/3) * 2.
//
// Time:  O(log n) for the exponentiation (or O(n/3) if multiplied in a loop).
// Space: O(1).
func mathThrees(n int) int {
	if n <= 3 {
		return n - 1 // 2→1*1=1, 3→1*2=2 (forced to break)
	}
	product := 1
	for n > 4 { // keep pulling out 3s while >4 remains
		product *= 3
		n -= 3
	}
	// Now n is 2, 3, or 4 — multiply the small remainder in directly.
	// (n==4 → 2*2=4, which is optimal for a remainder of 4.)
	return product * n
}

// ── Approach 3: DP Top-Down (Memoised Recursion) ─────────────────────────────
//
// dpTopDown solves Integer Break with recursion + memoisation.
//
// Intuition:
//
//	Same recurrence as bottom-up, expressed recursively. break(i) tries every
//	first cut j and takes max(j*(i-j), j*break(i-j)). Memoise to avoid
//	recomputing overlapping subproblems.
//
// Algorithm:
//  1. memo[i] caches the answer for i.
//  2. break(i): for j=1..i-1, best = max(j*(i-j), j*break(i-j)); cache & return.
//
// Time:  O(n^2). Space: O(n) for memo + recursion.
func dpTopDown(n int) int {
	memo := make([]int, n+1)
	var solve func(i int) int
	solve = func(i int) int {
		if i == 1 {
			return 1 // factor of 1
		}
		if memo[i] != 0 {
			return memo[i]
		}
		best := 0
		for j := 1; j < i; j++ {
			best = max(best, max(j*(i-j), j*solve(i-j)))
		}
		memo[i] = best
		return best
	}
	return solve(n)
}

func main() {
	fmt.Println("=== Approach 1: DP Bottom-Up ===")
	fmt.Println(dpBottomUp(2))  // expected 1
	fmt.Println(dpBottomUp(10)) // expected 36

	fmt.Println("=== Approach 2: Math — Break into 3s (Optimal) ===")
	fmt.Println(mathThrees(2))  // expected 1
	fmt.Println(mathThrees(10)) // expected 36

	fmt.Println("=== Approach 3: DP Top-Down (Memoised) ===")
	fmt.Println(dpTopDown(2))  // expected 1
	fmt.Println(dpTopDown(10)) // expected 36
}
