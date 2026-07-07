package main

import "fmt"

// ── Approach 1: Recursion (Brute Force) ──────────────────────────────────────
//
// recursion matches string s against pattern p using plain recursion.
//
// Intuition:
//   Two pointers i (into s) and j (into p). At each step:
//   - If p[j+1] == '*': either skip the '*' pair (zero occurrences) or,
//     if p[j] matches s[i], consume one char from s and recurse with same j.
//   - Otherwise: match the current characters and advance both.
//   Base case: j == len(p) means the pattern is exhausted; return i == len(s).
//
// Time:  O(2^(m+n)) worst case — exponential due to repeated sub-problems.
// Space: O(m+n) — recursion stack depth.
func recursion(s, p string) bool {
	return recurse(s, p, 0, 0)
}

func recurse(s, p string, i, j int) bool {
	// Pattern exhausted: true only if string is also exhausted.
	if j == len(p) {
		return i == len(s)
	}

	// Does the current pattern character match the current string character?
	firstMatch := i < len(s) && (p[j] == s[i] || p[j] == '.')

	if j+1 < len(p) && p[j+1] == '*' {
		// Two choices for '*':
		//   1. Zero occurrences: skip the "x*" pair entirely.
		//   2. One+ occurrences: if first char matches, consume one from s.
		return recurse(s, p, i, j+2) || (firstMatch && recurse(s, p, i+1, j))
	}

	// No '*' following: must match current chars and advance both pointers.
	return firstMatch && recurse(s, p, i+1, j+1)
}

// ── Approach 2: Top-Down DP (Memoisation) ────────────────────────────────────
//
// topDownDP adds a memo cache to the recursion to avoid recomputing (i,j) pairs.
//
// Intuition:
//   The recursive solution recomputes the same (i,j) sub-problems many times.
//   A 2-D cache of size (m+1) × (n+1) stores known results.
//
// Time:  O(m × n) — at most (m+1)*(n+1) unique states.
// Space: O(m × n) — the memo table plus O(m+n) stack.
func topDownDP(s, p string) bool {
	memo := make([][]int, len(s)+1)
	for i := range memo {
		memo[i] = make([]int, len(p)+1)
		// 0 = unvisited, 1 = true, -1 = false
	}
	return memoRecurse(s, p, 0, 0, memo)
}

func memoRecurse(s, p string, i, j int, memo [][]int) bool {
	if memo[i][j] != 0 {
		return memo[i][j] == 1
	}
	var result bool
	if j == len(p) {
		result = i == len(s)
	} else {
		firstMatch := i < len(s) && (p[j] == s[i] || p[j] == '.')
		if j+1 < len(p) && p[j+1] == '*' {
			result = memoRecurse(s, p, i, j+2, memo) || (firstMatch && memoRecurse(s, p, i+1, j, memo))
		} else {
			result = firstMatch && memoRecurse(s, p, i+1, j+1, memo)
		}
	}
	if result {
		memo[i][j] = 1
	} else {
		memo[i][j] = -1
	}
	return result
}

// ── Approach 3: Bottom-Up DP (Optimal) ───────────────────────────────────────
//
// bottomUpDP fills a 2-D DP table iteratively.
//
// Intuition:
//   dp[i][j] = true means s[i:] matches p[j:].
//   Fill from the bottom-right corner upward.
//   Base cases:
//     dp[len(s)][len(p)] = true  (both exhausted)
//     For j < len(p): dp[len(s)][j] = true only if the remaining pattern
//       consists entirely of "x*" pairs, e.g. "a*b*" matches empty string.
//   Transition:
//     If p[j+1] == '*':
//       dp[i][j] = dp[i][j+2]                          (zero occurrences)
//                  || (firstMatch && dp[i+1][j])        (one or more)
//     Else:
//       dp[i][j] = firstMatch && dp[i+1][j+1]
//
// Time:  O(m × n).
// Space: O(m × n) — the DP table; can be reduced to O(n) with rolling rows.
func bottomUpDP(s, p string) bool {
	m, n := len(s), len(p)
	// dp[i][j] = true means s[i:] matches p[j:]
	dp := make([][]bool, m+1)
	for i := range dp {
		dp[i] = make([]bool, n+1)
	}

	dp[m][n] = true // both exhausted

	// Fill the last row: pattern-only matching against empty string.
	// "a*b*c*" can match "": each x* pair can mean zero occurrences.
	for j := n - 2; j >= 0; j -= 2 {
		if p[j+1] == '*' {
			dp[m][j] = dp[m][j+2]
		}
	}

	// Fill table from bottom-right to top-left.
	for i := m - 1; i >= 0; i-- {
		for j := n - 1; j >= 0; j-- {
			firstMatch := p[j] == s[i] || p[j] == '.'

			if j+1 < n && p[j+1] == '*' {
				// Zero occurrences of x* OR one+ if current chars match.
				dp[i][j] = dp[i][j+2] || (firstMatch && dp[i+1][j])
			} else {
				dp[i][j] = firstMatch && dp[i+1][j+1]
			}
		}
	}

	return dp[0][0]
}

func main() {
	examples := []struct {
		s, p   string
		expect bool
	}{
		{"aa", "a", false},
		{"aa", "a*", true},
		{"ab", ".*", true},
		{"aab", "c*a*b", true},
		{"mississippi", "mis*is*p*.", false},
		{"", "a*", true},
		{"a", ".", true},
		{"aaa", "a*a", true},
	}

	approaches := []struct {
		name string
		fn   func(string, string) bool
	}{
		{"Approach 1: Recursion (brute)        O(2^(m+n)) T | O(m+n) S", recursion},
		{"Approach 2: Top-Down DP (memo)       O(m×n)     T | O(m×n) S", topDownDP},
		{"Approach 3: Bottom-Up DP           ✅ O(m×n)     T | O(m×n) S", bottomUpDP},
	}

	for _, ex := range examples {
		fmt.Printf("s=%-14q  p=%-10q  expect=%v\n", ex.s, ex.p, ex.expect)
		for _, ap := range approaches {
			fmt.Printf("  %-62s → %v\n", ap.name, ap.fn(ex.s, ex.p))
		}
		fmt.Println()
	}
}
