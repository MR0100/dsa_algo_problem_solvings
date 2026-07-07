package main

import "fmt"

// ── Approach 1: Memoized Recursion ────────────────────────────────────────────
//
// isInterleave solves Interleaving String using top-down DP.
//
// Intuition:
//   s3 is an interleaving of s1 and s2 iff we can assign each character of
//   s3 to either s1 or s2 in order, consuming them from left to right.
//   dp(i, j) = true iff s3[i+j:] can be formed by interleaving s1[i:] and s2[j:].
//
//   At each step:
//   - If s1[i] == s3[i+j]: try consuming from s1 → dp(i+1, j).
//   - If s2[j] == s3[i+j]: try consuming from s2 → dp(i, j+1).
//
// Time:  O(m × n) — m*n unique states, O(1) per state.
// Space: O(m × n) — memo table.
func isInterleave(s1 string, s2 string, s3 string) bool {
	m, n := len(s1), len(s2)
	if m+n != len(s3) {
		return false
	}
	// memo[i][j]: -1=not visited, 0=false, 1=true
	memo := make([][]int, m+1)
	for i := range memo {
		memo[i] = make([]int, n+1)
		for j := range memo[i] {
			memo[i][j] = -1
		}
	}

	var dp func(i, j int) bool
	dp = func(i, j int) bool {
		if i == m && j == n {
			return true // consumed both strings
		}
		if memo[i][j] != -1 {
			return memo[i][j] == 1
		}
		result := false
		k := i + j // current index in s3
		if i < m && s1[i] == s3[k] {
			result = dp(i+1, j)
		}
		if !result && j < n && s2[j] == s3[k] {
			result = dp(i, j+1)
		}
		if result {
			memo[i][j] = 1
		} else {
			memo[i][j] = 0
		}
		return result
	}
	return dp(0, 0)
}

// ── Approach 2: 2D Bottom-Up DP ───────────────────────────────────────────────
//
// isInterleaveDP solves Interleaving String using bottom-up DP.
//
// Intuition:
//   dp[i][j] = true iff s3[0:i+j] can be formed by interleaving s1[0:i] and s2[0:j].
//   dp[0][0] = true.
//   dp[i][j] = (dp[i-1][j] && s1[i-1]==s3[i+j-1]) || (dp[i][j-1] && s2[j-1]==s3[i+j-1])
//
// Time:  O(m × n)
// Space: O(m × n), reducible to O(n) with rolling row.
func isInterleaveDP(s1 string, s2 string, s3 string) bool {
	m, n := len(s1), len(s2)
	if m+n != len(s3) {
		return false
	}
	dp := make([][]bool, m+1)
	for i := range dp {
		dp[i] = make([]bool, n+1)
	}
	dp[0][0] = true

	// fill first row (using only s2)
	for j := 1; j <= n; j++ {
		dp[0][j] = dp[0][j-1] && s2[j-1] == s3[j-1]
	}
	// fill first column (using only s1)
	for i := 1; i <= m; i++ {
		dp[i][0] = dp[i-1][0] && s1[i-1] == s3[i-1]
	}
	// fill rest
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			dp[i][j] = (dp[i-1][j] && s1[i-1] == s3[i+j-1]) ||
				(dp[i][j-1] && s2[j-1] == s3[i+j-1])
		}
	}
	return dp[m][n]
}

// ── Approach 3: O(n) Space DP ─────────────────────────────────────────────────
//
// isInterleaveO1 solves Interleaving String using a single 1D DP array.
//
// Time:  O(m × n)
// Space: O(n)
func isInterleaveO1(s1 string, s2 string, s3 string) bool {
	m, n := len(s1), len(s2)
	if m+n != len(s3) {
		return false
	}
	dp := make([]bool, n+1)
	dp[0] = true
	for j := 1; j <= n; j++ {
		dp[j] = dp[j-1] && s2[j-1] == s3[j-1]
	}
	for i := 1; i <= m; i++ {
		dp[0] = dp[0] && s1[i-1] == s3[i-1]
		for j := 1; j <= n; j++ {
			dp[j] = (dp[j] && s1[i-1] == s3[i+j-1]) || (dp[j-1] && s2[j-1] == s3[i+j-1])
		}
	}
	return dp[n]
}

func main() {
	fmt.Println("=== Approach 1: Memoized Recursion ===")
	fmt.Printf("s1=%q s2=%q s3=%q  got=%v  expected true\n", "aabcc", "dbbca", "aadbbcbcac", isInterleave("aabcc", "dbbca", "aadbbcbcac"))
	fmt.Printf("s1=%q s2=%q s3=%q  got=%v  expected false\n", "aabcc", "dbbca", "aadbbbaccc", isInterleave("aabcc", "dbbca", "aadbbbaccc"))
	fmt.Printf("s1=%q s2=%q s3=%q  got=%v  expected true\n", "", "", "", isInterleave("", "", ""))

	fmt.Println("=== Approach 2: 2D Bottom-Up DP ===")
	fmt.Printf("s1=%q s2=%q s3=%q  got=%v  expected true\n", "aabcc", "dbbca", "aadbbcbcac", isInterleaveDP("aabcc", "dbbca", "aadbbcbcac"))
	fmt.Printf("s1=%q s2=%q s3=%q  got=%v  expected false\n", "aabcc", "dbbca", "aadbbbaccc", isInterleaveDP("aabcc", "dbbca", "aadbbbaccc"))

	fmt.Println("=== Approach 3: O(n) Space ===")
	fmt.Printf("s1=%q s2=%q s3=%q  got=%v  expected true\n", "aabcc", "dbbca", "aadbbcbcac", isInterleaveO1("aabcc", "dbbca", "aadbbcbcac"))
	fmt.Printf("s1=%q s2=%q s3=%q  got=%v  expected false\n", "aabcc", "dbbca", "aadbbbaccc", isInterleaveO1("aabcc", "dbbca", "aadbbbaccc"))
}
