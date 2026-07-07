package main

import "fmt"

// ── Approach 1: Brute Force ───────────────────────────────────────────────────
//
// bruteForce generates every string of length 2n using '(' and ')', then
// filters out invalid ones.
//
// Intuition:
//   2n positions, each either '(' or ')' → 2^(2n) candidate strings.
//   Check each for validity with a balance counter.
//
// Time:  O(2^(2n) · n) — generate 2^(2n) strings, each validated in O(n).
// Space: O(2^(2n) · n) — all generated strings.
func bruteForce(n int) []string {
	var result []string
	generate(make([]byte, 2*n), 0, n, &result)
	return result
}

func generate(current []byte, pos, n int, result *[]string) {
	if pos == 2*n {
		if isValid(current) {
			*result = append(*result, string(current))
		}
		return
	}
	current[pos] = '('
	generate(current, pos+1, n, result)
	current[pos] = ')'
	generate(current, pos+1, n, result)
}

func isValid(s []byte) bool {
	balance := 0
	for _, ch := range s {
		if ch == '(' {
			balance++
		} else {
			balance--
		}
		if balance < 0 {
			return false
		}
	}
	return balance == 0
}

// ── Approach 2: Backtracking ──────────────────────────────────────────────────
//
// backtracking builds only valid strings by tracking open and close counts.
//
// Intuition:
//   At each position we can:
//     - Add '(' if we still have opens left (open < n).
//     - Add ')' if there are more opens than closes (close < open).
//   This prunes all invalid branches — we never over-close or over-open.
//   When both counts reach n, record the complete string.
//
// Time:  O(4^n / sqrt(n)) — the n-th Catalan number; tight bound on valid strings.
// Space: O(n) — recursion stack depth (2n levels) + current path buffer.
func backtracking(n int) []string {
	var result []string
	btHelper(n, 0, 0, []byte{}, &result)
	return result
}

func btHelper(n, open, close int, path []byte, result *[]string) {
	if len(path) == 2*n {
		*result = append(*result, string(path))
		return
	}
	if open < n {
		btHelper(n, open+1, close, append(path, '('), result)
	}
	if close < open {
		btHelper(n, open, close+1, append(path, ')'), result)
	}
}

// ── Approach 3: Dynamic Programming ──────────────────────────────────────────
//
// dpApproach builds dp[k] (all valid strings of k pairs) from dp[0..k-1].
//
// Intuition:
//   Every valid parenthesisation of n pairs can be written as:
//     ( dp[i] ) dp[j]  where i + j = n-1, i,j >= 0.
//   The first '(' matches some ')' at position 2i+1, splitting the string
//   into an "inner" part (dp[i]) and an "outer" part (dp[j]).
//   Build dp[0..n] from dp[0] = [""].
//
// Time:  O(4^n / sqrt(n)) — same count of valid strings to produce.
// Space: O(4^n / sqrt(n)) — all generated strings stored in dp table.
func dpApproach(n int) []string {
	// dp[k] = all valid parenthesisations of k pairs.
	dp := make([][]string, n+1)
	dp[0] = []string{""}

	for k := 1; k <= n; k++ {
		var combinations []string
		for i := 0; i < k; i++ {
			j := k - 1 - i
			for _, inner := range dp[i] {
				for _, outer := range dp[j] {
					combinations = append(combinations, "("+inner+")"+outer)
				}
			}
		}
		dp[k] = combinations
	}
	return dp[n]
}

func main() {
	examples := []struct {
		n      int
		expect int // number of valid combinations (Catalan number)
	}{
		{1, 1},
		{2, 2},
		{3, 5},
		{4, 14},
	}

	approaches := []struct {
		name string
		fn   func(int) []string
	}{
		{"Approach 1: Brute Force      O(2^2n·n) T", bruteForce},
		{"Approach 2: Backtracking   ✅ O(4^n/√n) T", backtracking},
		{"Approach 3: DP (Catalan)     O(4^n/√n) T", dpApproach},
	}

	for _, ex := range examples {
		fmt.Printf("n=%d  expect=%d valid combinations\n", ex.n, ex.expect)
		for _, ap := range approaches {
			res := ap.fn(ex.n)
			fmt.Printf("  %-40s → count=%d  %v\n", ap.name, len(res), res)
		}
		fmt.Println()
	}
}
