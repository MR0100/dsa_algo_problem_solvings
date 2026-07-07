package main

import (
	"fmt"
	"strings"
)

// ── Approach 1: Brute Force (Sliding Window Character Comparison) ─────────────
//
// bruteForce solves Find the Index of the First Occurrence by checking every
// possible start position in haystack.
//
// Intuition: At each position i in haystack (0 to len(haystack)-len(needle)),
// compare haystack[i..i+len(needle)-1] with needle character by character.
// Return i on the first match.
//
// Algorithm:
//  1. For i = 0 to len(haystack)-len(needle):
//     compare haystack[i:i+m] == needle.
//     if all chars match: return i.
//  2. Return -1.
//
// Time:  O((n-m+1) * m) ≈ O(n*m) where n=len(haystack), m=len(needle)
// Space: O(1)
func bruteForce(haystack, needle string) int {
	n, m := len(haystack), len(needle)
	if m == 0 {
		return 0
	}
	for i := 0; i <= n-m; i++ {
		j := 0
		for j < m && haystack[i+j] == needle[j] { // compare character by character
			j++
		}
		if j == m { // all m characters matched
			return i
		}
	}
	return -1
}

// ── Approach 2: Go Standard Library ──────────────────────────────────────────
//
// useStdlib solves Find the Index of the First Occurrence using strings.Index.
//
// Intuition: Go's standard library implements an optimised string search
// (Rabin-Karp or similar). In an interview this shows awareness of stdlib,
// but you should also know how to implement it manually.
//
// Time:  O(n) average (Go's strings.Index uses an optimised algorithm)
// Space: O(1)
func useStdlib(haystack, needle string) int {
	return strings.Index(haystack, needle)
}

// ── Approach 3: KMP (Knuth-Morris-Pratt) — Optimal ───────────────────────────
//
// kmp solves Find the Index of the First Occurrence using the KMP algorithm.
//
// Intuition: Build a "failure function" (also called the LPS — Longest Proper
// Prefix which is also Suffix — array) for the needle. During matching, when
// a mismatch occurs, the failure function tells us how far we can shift the
// needle without missing a match, instead of starting over from scratch.
//
// Algorithm:
//  1. Build lps[0..m-1]:
//     lps[i] = length of the longest proper prefix of needle[0..i] that is
//     also a suffix. E.g. needle="AAACAAAA" → lps=[0,1,2,0,1,2,3,3].
//  2. Match with two pointers i (haystack) and j (needle):
//     if haystack[i]==needle[j]: advance both.
//     if j==m: found at i-m; reset j = lps[j-1].
//     else if mismatch and j>0: j = lps[j-1] (don't move i).
//     else: advance i.
//
// Time:  O(n+m) — O(m) to build lps, O(n) to search
// Space: O(m) — lps array
func kmp(haystack, needle string) int {
	n, m := len(haystack), len(needle)
	if m == 0 {
		return 0
	}

	// build LPS (Longest Proper Prefix-Suffix) table
	lps := make([]int, m)
	length := 0 // length of previous longest prefix-suffix
	for i := 1; i < m; {
		if needle[i] == needle[length] {
			length++
			lps[i] = length
			i++
		} else if length != 0 {
			// don't increment i; try the previous longest prefix-suffix
			length = lps[length-1]
		} else {
			lps[i] = 0
			i++
		}
	}

	// search
	i, j := 0, 0 // i scans haystack, j scans needle
	for i < n {
		if haystack[i] == needle[j] {
			i++
			j++
		}
		if j == m { // full needle matched
			return i - m
		} else if i < n && haystack[i] != needle[j] {
			if j != 0 {
				j = lps[j-1] // shift needle using failure function
			} else {
				i++ // no prefix to fall back to; advance haystack
			}
		}
	}
	return -1
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Printf("sadbutsad / sad  = %d  expected 0\n", bruteForce("sadbutsad", "sad"))
	fmt.Printf("leetcode  / leeto = %d  expected -1\n", bruteForce("leetcode", "leeto"))
	fmt.Printf("\"\" / \"\"            = %d  expected 0\n", bruteForce("", ""))
	fmt.Printf("a / a             = %d  expected 0\n", bruteForce("a", "a"))

	fmt.Println("\n=== Approach 2: Standard Library ===")
	fmt.Printf("sadbutsad / sad  = %d  expected 0\n", useStdlib("sadbutsad", "sad"))
	fmt.Printf("leetcode  / leeto = %d  expected -1\n", useStdlib("leetcode", "leeto"))

	fmt.Println("\n=== Approach 3: KMP (Optimal) ===")
	fmt.Printf("sadbutsad / sad  = %d  expected 0\n", kmp("sadbutsad", "sad"))
	fmt.Printf("leetcode  / leeto = %d  expected -1\n", kmp("leetcode", "leeto"))
	fmt.Printf("aabaabaaf / aabaaf = %d  expected 3\n", kmp("aabaabaaf", "aabaaf"))
	fmt.Printf("mississippi / issip = %d  expected 4\n", kmp("mississippi", "issip"))
}
