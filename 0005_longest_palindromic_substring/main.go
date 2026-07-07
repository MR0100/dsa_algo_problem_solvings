package main

import "fmt"

// ── Approach 1: Brute Force ───────────────────────────────────────────────────
//
// bruteForce checks every substring of s and returns the longest palindrome.
//
// Intuition:
//   Generate every (i,j) pair, extract s[i..j], test if it is a palindrome
//   by comparing characters from both ends inward, keep the longest one found.
//
// Time:  O(n³) — O(n²) substrings × O(n) palindrome check each.
// Space: O(1)  — only indices; no extra allocation.
func bruteForce(s string) string {
	n := len(s)
	best := s[0:1] // single character is always a palindrome

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if j-i+1 > len(best) && isPalindrome(s, i, j) {
				best = s[i : j+1]
			}
		}
	}
	return best
}

// isPalindrome reports whether s[l..r] is a palindrome.
func isPalindrome(s string, l, r int) bool {
	for l < r {
		if s[l] != s[r] {
			return false
		}
		l++
		r--
	}
	return true
}

// ── Approach 2: Dynamic Programming ──────────────────────────────────────────
//
// dpApproach builds a 2-D boolean table where dp[i][j] = true means s[i..j]
// is a palindrome, filling from shorter to longer substrings.
//
// Intuition:
//   s[i..j] is a palindrome iff s[i]==s[j] AND s[i+1..j-1] is a palindrome.
//   Base cases: every single character (length 1) and every matching adjacent
//   pair (length 2) are palindromes. Build up to length n.
//
// Algorithm:
//   1. dp[i][i] = true for all i.
//   2. For each length L from 2 to n:
//        For each start i where end j = i+L-1 < n:
//          dp[i][j] = (s[i]==s[j]) && (L==2 || dp[i+1][j-1])
//          Track longest.
//
// Time:  O(n²) — fill n² cells.
// Space: O(n²) — the DP table.
func dpApproach(s string) string {
	n := len(s)
	// dp[i][j] = true means s[i..j] is a palindrome.
	dp := make([][]bool, n)
	for i := range dp {
		dp[i] = make([]bool, n)
		dp[i][i] = true // single characters are palindromes
	}

	start, maxLen := 0, 1

	// Fill by increasing substring length.
	for length := 2; length <= n; length++ {
		for i := 0; i <= n-length; i++ {
			j := i + length - 1 // end index

			if length == 2 {
				dp[i][j] = s[i] == s[j]
			} else {
				// s[i..j] is palindrome iff outer chars match AND inner is palindrome.
				dp[i][j] = (s[i] == s[j]) && dp[i+1][j-1]
			}

			if dp[i][j] && length > maxLen {
				start, maxLen = i, length
			}
		}
	}
	return s[start : start+maxLen]
}

// ── Approach 3: Expand Around Center ─────────────────────────────────────────
//
// expandAroundCenter treats each of the 2n-1 possible centers and expands
// outward while characters match, recording the widest palindrome seen.
//
// Intuition:
//   Every palindrome mirrors around a center. The center is either:
//     - A single character (odd-length palindromes like "aba").
//     - The gap between two characters (even-length like "abba").
//   There are n odd-center + (n-1) even-center = 2n-1 centers total.
//   For each center, expand outward while s[l]==s[r], tracking the best span.
//
// Algorithm:
//   For each center (represented as an index in a virtual doubled string):
//     expand(l, r) outward while s[l]==s[r] and in bounds.
//     Update [start, end] if r-l+1 > maxLen.
//
// Time:  O(n²) — 2n-1 centers × up to O(n) expansion each.
// Space: O(1)  — only index variables; no table.
func expandAroundCenter(s string) string {
	start, end := 0, 0

	for i := 0; i < len(s); i++ {
		// Odd-length: center is a single character at i.
		l1, r1 := expand(s, i, i)
		// Even-length: center is the gap between i and i+1.
		l2, r2 := expand(s, i, i+1)

		if r1-l1 > end-start {
			start, end = l1, r1
		}
		if r2-l2 > end-start {
			start, end = l2, r2
		}
	}
	return s[start : end+1]
}

// expand grows outward from center (l, r) while characters match,
// returning the inclusive [l, r] of the largest palindrome found.
func expand(s string, l, r int) (int, int) {
	for l >= 0 && r < len(s) && s[l] == s[r] {
		l--
		r++
	}
	// Step back: l and r overshot by one on each side.
	return l + 1, r - 1
}

// ── Approach 4: Manacher's Algorithm ─────────────────────────────────────────
//
// manacher computes the longest palindromic substring in O(n) using the
// Manacher's algorithm, which reuses previously computed palindrome radii.
//
// Intuition:
//   Transform s into a new string T with separators (#) between every character
//   (e.g. "abc" → "#a#b#c#") so all palindromes become odd-length and centers
//   are always at character positions.
//   Maintain p[i] = radius of the palindrome centered at T[i].
//   Use a "current rightmost palindrome" [c, r] to avoid re-expanding:
//     If T[i] is within [c, r], its mirror i' = 2c-i already has a computed
//     radius p[i'], giving us a free starting radius for T[i].
//   Expand beyond the guaranteed radius, then update [c, r] if needed.
//
// Time:  O(n) — amortised; each character is visited at most twice.
// Space: O(n) — the transformed string T and the radius array p.
func manacher(s string) string {
	// Transform: "abc" → "#a#b#c#"
	t := "#"
	for _, ch := range s {
		t += string(ch) + "#"
	}
	n := len(t)
	p := make([]int, n) // p[i] = palindrome radius centered at t[i]

	center, right := 0, 0 // center and right boundary of rightmost palindrome

	for i := 0; i < n; i++ {
		mirror := 2*center - i // mirror of i around center

		if i < right {
			// Use the mirror's radius, but don't go beyond right boundary.
			p[i] = min(right-i, p[mirror])
		}

		// Attempt to expand beyond current radius.
		l, r := i-p[i]-1, i+p[i]+1
		for l >= 0 && r < n && t[l] == t[r] {
			p[i]++
			l--
			r++
		}

		// Update the rightmost palindrome if this one extends further right.
		if i+p[i] > right {
			center = i
			right = i + p[i]
		}
	}

	// Find the maximum radius and its center in the original string.
	maxRadius, maxCenter := 0, 0
	for i, radius := range p {
		if radius > maxRadius {
			maxRadius, maxCenter = radius, i
		}
	}

	// Map back to the original string:
	// center in T at maxCenter, radius maxRadius → original start = (maxCenter-maxRadius)/2
	start := (maxCenter - maxRadius) / 2
	return s[start : start+maxRadius]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	examples := []struct {
		s      string
		expect string // one valid answer (multiple may exist)
	}{
		{"babad", "bab"},   // "aba" also valid
		{"cbbd", "bb"},
		{"a", "a"},
		{"ac", "a"},
		{"racecar", "racecar"},
	}

	approaches := []struct {
		name string
		fn   func(string) string
	}{
		{"Approach 1: Brute Force             O(n³) T | O(1)   S", bruteForce},
		{"Approach 2: Dynamic Programming     O(n²) T | O(n²)  S", dpApproach},
		{"Approach 3: Expand Around Center ✅ O(n²) T | O(1)   S", expandAroundCenter},
		{"Approach 4: Manacher's Algorithm    O(n)  T | O(n)   S", manacher},
	}

	for _, ex := range examples {
		fmt.Printf("s=%q\n", ex.s)
		for _, ap := range approaches {
			fmt.Printf("  %-60s → %q\n", ap.name, ap.fn(ex.s))
		}
		fmt.Println()
	}
}
