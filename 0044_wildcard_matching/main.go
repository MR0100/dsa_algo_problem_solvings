package main

import "fmt"

// ── Approach 1: Recursion (Brute Force) ──────────────────────────────────────
//
// recursion solves Wildcard Matching with pure recursion.
//
// Intuition:
//   - '?' matches any single character: recurse with (s[1:], p[1:]).
//   - '*' matches zero characters: try (s, p[1:]).
//          OR one or more characters: try (s[1:], p) [keep '*' to match more].
//   - Literal match: recurse with (s[1:], p[1:]).
//   - No match: return false.
//
// Time:  O(2^(len(s)+len(p))) — exponential due to '*' branching
// Space: O(len(s)+len(p)) — recursion stack
func recursion(s, p string) bool {
	if len(p) == 0 {
		return len(s) == 0
	}
	if p[0] == '*' {
		return recursion(s, p[1:]) || // '*' matches zero chars
			(len(s) > 0 && recursion(s[1:], p)) // '*' matches one or more
	}
	if len(s) > 0 && (p[0] == '?' || p[0] == s[0]) {
		return recursion(s[1:], p[1:])
	}
	return false
}

// ── Approach 2: DP Bottom-Up (Optimal) ───────────────────────────────────────
//
// dpBottomUp solves Wildcard Matching using a 2D DP table.
//
// Intuition: dp[i][j] = true if s[0..i-1] matches p[0..j-1].
//
// Base cases:
//   dp[0][0] = true (empty string matches empty pattern).
//   dp[0][j] = true if p[0..j-1] is all '*'s.
//
// Transitions (1-indexed):
//   if p[j-1] == '*':
//     dp[i][j] = dp[i][j-1]  ('*' matches zero chars, advance pattern)
//             || dp[i-1][j]   ('*' matches one more char of s, keep pattern)
//   elif p[j-1] == '?' or p[j-1] == s[i-1]:
//     dp[i][j] = dp[i-1][j-1]
//   else:
//     dp[i][j] = false
//
// Time:  O(len(s) * len(p))
// Space: O(len(s) * len(p)) — can be reduced to O(len(s)) with rolling array
func dpBottomUp(s, p string) bool {
	m, n := len(s), len(p)
	dp := make([][]bool, m+1)
	for i := range dp {
		dp[i] = make([]bool, n+1)
	}
	dp[0][0] = true // empty matches empty

	// empty string can only match all-'*' pattern
	for j := 1; j <= n; j++ {
		if p[j-1] == '*' {
			dp[0][j] = dp[0][j-1]
		}
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if p[j-1] == '*' {
				// '*' matches zero chars from s  → dp[i][j-1]
				// '*' matches one more char of s → dp[i-1][j]
				dp[i][j] = dp[i][j-1] || dp[i-1][j]
			} else if p[j-1] == '?' || p[j-1] == s[i-1] {
				dp[i][j] = dp[i-1][j-1]
			}
		}
	}
	return dp[m][n]
}

// ── Approach 3: Two Pointers with '*' Bookmark (Optimal O(1) Space) ──────────
//
// twoPointers solves Wildcard Matching in O(n) space (O(1) if s/p are arrays)
// using a greedy two-pointer approach with bookmarking.
//
// Intuition: Walk s and p simultaneously. When we see '*':
//   - Record bookmark: starIdx = j, match = i (position in s when '*' was seen).
//   - Assume '*' matches zero chars: advance only j.
// On mismatch:
//   - If we have a bookmark, '*' can match one more char of s:
//     advance match; reset i=match, j=starIdx+1.
//
// Time:  O(len(s) * len(p)) worst case (many '*'); O(len(s)) for simple cases
// Space: O(1)
func twoPointers(s, p string) bool {
	i, j := 0, 0           // pointers into s and p
	starIdx, match := -1, 0 // bookmark: position of last '*' in p, and where in s

	for i < len(s) {
		if j < len(p) && (p[j] == '?' || p[j] == s[i]) {
			i++; j++ // direct match or '?' match
		} else if j < len(p) && p[j] == '*' {
			starIdx = j  // record '*' position in p
			match = i    // '*' currently matches zero chars of s
			j++          // advance pattern; i stays
		} else if starIdx != -1 {
			// mismatch but we have a '*' to fall back on: extend it by one char
			match++
			i = match    // restart s from the next position after '*'
			j = starIdx + 1 // restart p from after the '*'
		} else {
			return false // no '*' to fall back on
		}
	}

	// consume any trailing '*'s in p
	for j < len(p) && p[j] == '*' {
		j++
	}
	return j == len(p)
}

func main() {
	cases := []struct {
		s, p string
		want bool
	}{
		{"aa", "a", false},
		{"aa", "*", true},
		{"cb", "?a", false},
		{"adceb", "*a*b", true},
		{"acdcb", "a*c?b", false},
		{"", "", true},
		{"", "*", true},
		{"", "**", true},
	}

	fmt.Println("=== Approach 1: Recursion ===")
	for _, c := range cases {
		fmt.Printf("s=%-10q p=%-8q => %v  expected %v\n", c.s, c.p, recursion(c.s, c.p), c.want)
	}

	fmt.Println("\n=== Approach 2: DP Bottom-Up (Optimal) ===")
	for _, c := range cases {
		fmt.Printf("s=%-10q p=%-8q => %v  expected %v\n", c.s, c.p, dpBottomUp(c.s, c.p), c.want)
	}

	fmt.Println("\n=== Approach 3: Two Pointers with Bookmark ===")
	for _, c := range cases {
		fmt.Printf("s=%-10q p=%-8q => %v  expected %v\n", c.s, c.p, twoPointers(c.s, c.p), c.want)
	}
}
