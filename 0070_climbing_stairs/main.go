package main

import "fmt"

// ── Approach 1: Recursion with Memoization ────────────────────────────────────
//
// memoization solves Climbing Stairs using top-down DP.
//
// Intuition:
//   To reach step n, the last move was either 1 step (from n-1) or 2 steps
//   (from n-2). So ways(n) = ways(n-1) + ways(n-2). This is Fibonacci.
//   Memoize to avoid recomputation.
//
// Time:  O(n)
// Space: O(n) — memo array + recursion stack.
func memoization(n int) int {
	memo := make([]int, n+1)
	var dp func(i int) int
	dp = func(i int) int {
		if i <= 1 {
			return 1 // base: 0 steps → 1 way; 1 step → 1 way
		}
		if memo[i] != 0 {
			return memo[i]
		}
		memo[i] = dp(i-1) + dp(i-2)
		return memo[i]
	}
	return dp(n)
}

// ── Approach 2: DP Bottom-Up ──────────────────────────────────────────────────
//
// dpBottomUp solves Climbing Stairs using an iterative DP table.
//
// Intuition:
//   dp[i] = ways(i). Build from dp[0]=1, dp[1]=1 upward.
//
// Time:  O(n)
// Space: O(n)
func dpBottomUp(n int) int {
	if n <= 1 {
		return 1
	}
	dp := make([]int, n+1)
	dp[0], dp[1] = 1, 1
	for i := 2; i <= n; i++ {
		dp[i] = dp[i-1] + dp[i-2]
	}
	return dp[n]
}

// ── Approach 3: Two Variables (Optimal) ──────────────────────────────────────
//
// twoVars solves Climbing Stairs in O(1) space by keeping only prev and curr.
//
// Intuition:
//   Since we only need the last two Fibonacci values, we can replace the
//   array with two variables and update in-place.
//
// Time:  O(n)
// Space: O(1)
func twoVars(n int) int {
	if n <= 1 {
		return 1
	}
	prev, curr := 1, 1
	for i := 2; i <= n; i++ {
		prev, curr = curr, prev+curr
	}
	return curr
}

func main() {
	cases := []struct {
		n        int
		expected int
	}{
		{1, 1},
		{2, 2},
		{3, 3},
		{4, 5},
		{5, 8},
		{10, 89},
		{45, 1836311903},
	}

	fmt.Println("=== Approach 1: Memoization ===")
	for _, c := range cases {
		fmt.Printf("n=%-3d  got=%-12d  expected=%d\n", c.n, memoization(c.n), c.expected)
	}

	fmt.Println("=== Approach 2: DP Bottom-Up ===")
	for _, c := range cases {
		fmt.Printf("n=%-3d  got=%-12d  expected=%d\n", c.n, dpBottomUp(c.n), c.expected)
	}

	fmt.Println("=== Approach 3: Two Variables ===")
	for _, c := range cases {
		fmt.Printf("n=%-3d  got=%-12d  expected=%d\n", c.n, twoVars(c.n), c.expected)
	}
}
