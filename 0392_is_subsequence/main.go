package main

import "fmt"

// ── Approach 1: Two Pointers (Optimal for a single query) ────────────────────
//
// twoPointers checks whether s is a subsequence of t by walking both strings
// once with two indices.
//
// Intuition:
//
//	A subsequence keeps relative order but may skip characters. Scan t left to
//	right; every time the current t-character equals the next needed
//	s-character, consume it (advance the s pointer). If we consume all of s
//	before t runs out, s is a subsequence. Greedily matching the earliest
//	possible t-character is always safe: matching later can never help.
//
// Algorithm:
//  1. i = 0 (index into s), j = 0 (index into t).
//  2. While i < len(s) and j < len(t):
//     if s[i] == t[j], advance i (matched one char); always advance j.
//  3. s is a subsequence iff i reached len(s).
//
// Time:  O(n) where n = len(t) — single pass over t.
// Space: O(1) — two integer indices.
func twoPointers(s string, t string) bool {
	i := 0 // how many chars of s are matched so far
	// Walk through t once, consuming s characters greedily in order.
	for j := 0; i < len(s) && j < len(t); j++ {
		if s[i] == t[j] { // current t char is the next one s needs
			i++ // consume it
		}
	}
	// Matched every character of s ⇒ subsequence.
	return i == len(s)
}

// ── Approach 2: DP Table (subsequence-matching lattice) ──────────────────────
//
// dpBottomUp computes, via a classic 2D boolean table, whether s is a
// subsequence of t. It is overkill for one query but shows the structure that
// generalizes to edit-distance-style problems.
//
// Intuition:
//
//	Let dp[i][j] = "is s[i:] a subsequence of t[j:]?". An empty s is always a
//	subsequence, so the last row is all true. Working backwards, s[i:] matches
//	t[j:] if either we skip t[j] (dp[i][j+1]) or, when s[i]==t[j], we match it
//	and recurse on dp[i+1][j+1].
//
// Algorithm:
//  1. dp has dimensions (len(s)+1) x (len(t)+1), all false; set dp[len(s)][*] = true.
//  2. Fill from i = len(s)-1 down to 0, j = len(t)-1 down to 0:
//     dp[i][j] = dp[i][j+1] OR (s[i]==t[j] AND dp[i+1][j+1]).
//  3. Answer is dp[0][0].
//
// Time:  O(m·n) where m = len(s), n = len(t).
// Space: O(m·n) for the table.
func dpBottomUp(s string, t string) bool {
	m, n := len(s), len(t)
	// dp[i][j] answers: is s[i:] a subsequence of t[j:]?
	dp := make([][]bool, m+1)
	for i := range dp {
		dp[i] = make([]bool, n+1)
	}
	// Base case: empty suffix of s (i == m) is a subsequence of anything.
	for j := 0; j <= n; j++ {
		dp[m][j] = true
	}
	// Fill backwards so each cell reads already-computed neighbours.
	for i := m - 1; i >= 0; i-- {
		for j := n - 1; j >= 0; j-- {
			skip := dp[i][j+1]                   // ignore t[j], keep looking for s[i]
			take := s[i] == t[j] && dp[i+1][j+1] // match s[i] with t[j]
			dp[i][j] = skip || take
		}
	}
	return dp[0][0]
}

// ── Approach 3: Preprocessed Next-Position Jump Table (Follow-up) ─────────────
//
// followUpManyQueries handles the follow-up: many s queries against the same
// t. It precomputes, for each position in t and each letter, the index of the
// NEXT occurrence of that letter, so each query runs in O(len(s)) with only
// jumps — no rescanning of t.
//
// Intuition:
//
//	If we must answer thousands of s's against one t, re-walking t each time is
//	wasteful. Build nxt[j][c] = smallest index >= j in t where letter c
//	appears (or len(t) if none). Then for a query s, jump letter by letter:
//	stand at position pos, look up nxt[pos][s[i]]; if it's len(t) the letter is
//	missing → fail, else move pos just past it.
//
// Algorithm:
//  1. Build nxt with rows len(t)+1 and 26 columns, filled bottom-up:
//     nxt[j][c] = j if t[j]==c else nxt[j+1][c]; last row all len(t).
//  2. For each query s: pos = 0; for each char, pos = nxt[pos][char]; if that
//     equals len(t) return false, else pos++ (move past the matched char).
//  3. Survive all chars ⇒ true.
//
// Time:  O(n·26) preprocessing, then O(len(s)) per query.
// Space: O(n·26) for the jump table.
func followUpManyQueries(s string, t string) bool {
	n := len(t)
	// nxt[j][c] = index of next occurrence of letter c at or after position j.
	nxt := make([][26]int, n+1)
	// Sentinel row: no letter occurs at or after position n.
	for c := 0; c < 26; c++ {
		nxt[n][c] = n
	}
	// Fill upward: each row copies the row below, then overrides its own letter.
	for j := n - 1; j >= 0; j-- {
		nxt[j] = nxt[j+1]    // default: next occurrence is wherever it was after j+1
		nxt[j][t[j]-'a'] = j // t[j] itself is available exactly at position j
	}
	// Answer this query by jumping through the table.
	pos := 0 // current search position in t
	for i := 0; i < len(s); i++ {
		c := s[i] - 'a'
		j := nxt[pos][c] // next place letter c appears from pos onward
		if j == n {
			return false // letter c never appears again ⇒ not a subsequence
		}
		pos = j + 1 // consume it, continue searching after it
	}
	return true
}

func main() {
	fmt.Println("=== Approach 1: Two Pointers ===")
	fmt.Printf("s=\"abc\", t=\"ahbgdc\": got=%v  expected true\n", twoPointers("abc", "ahbgdc"))  // expected true
	fmt.Printf("s=\"axc\", t=\"ahbgdc\": got=%v  expected false\n", twoPointers("axc", "ahbgdc")) // expected false
	fmt.Printf("s=\"\",    t=\"ahbgdc\": got=%v  expected true\n", twoPointers("", "ahbgdc"))     // expected true

	fmt.Println("=== Approach 2: DP Table ===")
	fmt.Printf("s=\"abc\", t=\"ahbgdc\": got=%v  expected true\n", dpBottomUp("abc", "ahbgdc"))  // expected true
	fmt.Printf("s=\"axc\", t=\"ahbgdc\": got=%v  expected false\n", dpBottomUp("axc", "ahbgdc")) // expected false
	fmt.Printf("s=\"\",    t=\"ahbgdc\": got=%v  expected true\n", dpBottomUp("", "ahbgdc"))     // expected true

	fmt.Println("=== Approach 3: Preprocessed Jump Table (Follow-up) ===")
	fmt.Printf("s=\"abc\", t=\"ahbgdc\": got=%v  expected true\n", followUpManyQueries("abc", "ahbgdc"))  // expected true
	fmt.Printf("s=\"axc\", t=\"ahbgdc\": got=%v  expected false\n", followUpManyQueries("axc", "ahbgdc")) // expected false
	fmt.Printf("s=\"\",    t=\"ahbgdc\": got=%v  expected true\n", followUpManyQueries("", "ahbgdc"))     // expected true
}
