package main

import "fmt"

// ── Approach 1: Brute Force (Full Edit-Distance DP) ──────────────────────────
//
// bruteForce solves One Edit Distance by computing the full Levenshtein
// edit distance and checking whether it equals exactly 1.
//
// Intuition:
//
//	"One edit distance apart" literally means the classic edit distance
//	(LeetCode #72) between s and t is exactly 1. So the most direct — if
//	heavyweight — solution is to fill the whole (m+1)×(n+1) DP table and
//	test dp[m][n] == 1. Note that equal strings have distance 0, which
//	correctly returns false.
//
// Algorithm:
//  1. Build table dp where dp[i][j] = edit distance between s[:i] and t[:j].
//  2. Base cases: dp[i][0] = i (delete i chars), dp[0][j] = j (insert j chars).
//  3. Transition: if s[i-1] == t[j-1], dp[i][j] = dp[i-1][j-1];
//     otherwise dp[i][j] = 1 + min(delete, insert, replace).
//  4. Return dp[m][n] == 1.
//
// Time:  O(m·n) — every cell of the table is filled once.
// Space: O(m·n) — the full DP table is stored.
func bruteForce(s, t string) bool {
	m, n := len(s), len(t)
	dp := make([][]int, m+1) // dp[i][j] = edit distance of s[:i] vs t[:j]
	for i := range dp {
		dp[i] = make([]int, n+1)
		dp[i][0] = i // turning s[:i] into "" needs i deletions
	}
	for j := 0; j <= n; j++ {
		dp[0][j] = j // turning "" into t[:j] needs j insertions
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if s[i-1] == t[j-1] {
				dp[i][j] = dp[i-1][j-1] // chars match → no edit needed here
			} else {
				del := dp[i-1][j]   // delete s[i-1]
				ins := dp[i][j-1]   // insert t[j-1]
				rep := dp[i-1][j-1] // replace s[i-1] with t[j-1]
				dp[i][j] = 1 + min(min(del, ins), rep)
			}
		}
	}
	return dp[m][n] == 1 // exactly one edit — 0 (equal) must return false
}

// ── Approach 2: First Mismatch + Suffix Comparison ───────────────────────────
//
// suffixCompare solves One Edit Distance by locating the first mismatching
// index and comparing the remaining suffixes directly.
//
// Intuition:
//
//	If s and t are one edit apart, they must agree on some (possibly empty)
//	prefix, then differ at exactly one position, then agree on the rest.
//	At the first mismatch there are only three possible repairs:
//	replace (equal lengths), insert into the shorter, or delete from the
//	longer — and each repair is verified by one suffix equality check.
//
// Algorithm:
//  1. If |len(s) − len(t)| > 1, more than one insert/delete is needed → false.
//  2. Walk both strings until the first index i where s[i] != t[i].
//  3. At that mismatch:
//     - equal lengths  → true iff s[i+1:] == t[i+1:] (replace s[i]).
//     - s shorter      → true iff s[i:]   == t[i+1:] (insert t[i] into s).
//     - s longer       → true iff s[i+1:] == t[i:]   (delete s[i]).
//  4. No mismatch found in the overlap → true iff lengths differ by exactly 1
//     (the extra trailing character is the single edit).
//
// Time:  O(min(m, n)) scan + O(n) suffix comparison → O(m + n) overall.
// Space: O(1) — Go string slicing creates views, not copies.
func suffixCompare(s, t string) bool {
	m, n := len(s), len(t)
	diff := m - n
	if diff < 0 {
		diff = -diff // absolute length difference
	}
	if diff > 1 {
		return false // would need at least two inserts/deletes
	}
	shorter := m
	if n < m {
		shorter = n // only the overlapping prefix can be compared index-wise
	}
	for i := 0; i < shorter; i++ {
		if s[i] != t[i] { // first mismatch — decide which single edit fixes it
			switch {
			case m == n:
				return s[i+1:] == t[i+1:] // replace s[i] with t[i]
			case m < n:
				return s[i:] == t[i+1:] // insert t[i] into s at position i
			default:
				return s[i+1:] == t[i:] // delete s[i] from s
			}
		}
	}
	// The overlap matched entirely: strings are equal (diff == 0 → false,
	// because zero edits is not one edit) or one has a single extra tail char.
	return diff == 1
}

// ── Approach 3: One-Pass Two Pointers (Optimal) ──────────────────────────────
//
// twoPointers solves One Edit Distance with a single simultaneous walk over
// both strings, allowing at most one divergence.
//
// Intuition:
//
//	Walk pointers i (over the shorter string s) and j (over the longer t)
//	together. The first time the characters differ, "spend" the one allowed
//	edit: on equal lengths skip both characters (replace); on unequal
//	lengths skip only the longer string's character (insert/delete). Any
//	second difference proves the distance exceeds 1.
//
// Algorithm:
//  1. Swap so that s is the shorter (or equal) string; if the length gap
//     exceeds 1, return false immediately.
//  2. Advance i and j while characters match.
//  3. On the first mismatch: mark the edit as used; advance j always, and
//     advance i too only when lengths are equal (replacement).
//  4. On a second mismatch return false.
//  5. After the loop: true if one edit was used, or if the untouched tail of
//     t is exactly one character (lengths differ by 1 with no mismatch).
//
// Time:  O(min(m, n)) — each pointer moves forward only, one pass.
// Space: O(1) — two indices and one flag.
func twoPointers(s, t string) bool {
	m, n := len(s), len(t)
	if m > n {
		s, t = t, s // ensure s is the shorter string
		m, n = n, m
	}
	if n-m > 1 {
		return false // one insert/delete can bridge a gap of at most 1
	}
	i, j := 0, 0
	usedEdit := false // whether the single allowed edit has been consumed
	for i < m && j < n {
		if s[i] == t[j] { // characters agree → advance both pointers
			i++
			j++
			continue
		}
		if usedEdit {
			return false // second mismatch → at least two edits required
		}
		usedEdit = true
		if m == n {
			i++ // equal lengths → this mismatch must be a replacement
		}
		j++ // unequal lengths → skip t[j] (delete from t / insert into s)
	}
	// Either the edit was already spent (tails matched afterwards), or the
	// strings matched fully and t has exactly one leftover character.
	return usedEdit || n-m == 1
}

// min returns the smaller of two ints (helper for the DP approach).
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	type example struct {
		s, t string
	}
	examples := []example{
		{"ab", "acb"}, // expected true  (insert 'c' into s)
		{"", ""},      // expected false (0 edits, not 1)
	}

	fmt.Println("=== Approach 1: Brute Force (Full Edit-Distance DP) ===")
	for _, ex := range examples {
		fmt.Printf("s=%q t=%q  got=%v\n", ex.s, ex.t, bruteForce(ex.s, ex.t)) // expected true, false
	}

	fmt.Println("=== Approach 2: First Mismatch + Suffix Comparison ===")
	for _, ex := range examples {
		fmt.Printf("s=%q t=%q  got=%v\n", ex.s, ex.t, suffixCompare(ex.s, ex.t)) // expected true, false
	}

	fmt.Println("=== Approach 3: One-Pass Two Pointers (Optimal) ===")
	for _, ex := range examples {
		fmt.Printf("s=%q t=%q  got=%v\n", ex.s, ex.t, twoPointers(ex.s, ex.t)) // expected true, false
	}
}
