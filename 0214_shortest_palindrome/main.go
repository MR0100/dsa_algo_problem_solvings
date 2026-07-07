package main

import "fmt"

// ── Approach 1: Brute Force (Longest Palindromic Prefix) ──────────────────────
//
// bruteForce solves Shortest Palindrome by finding the longest prefix of s that
// is itself a palindrome, then prepending the reverse of the remaining suffix.
//
// Intuition:
//
//	We may only add characters to the FRONT. Whatever palindrome we build must
//	have s as its suffix. The characters we prepend mirror the tail of s that
//	is NOT already symmetric. Concretely: find the longest prefix s[0..k) that
//	is a palindrome — it needs nothing added in front of it. The rest, s[k:],
//	must be mirrored and placed at the very front: reverse(s[k:]) + s. Minimise
//	added characters ⇔ maximise the palindromic prefix length k.
//
// Algorithm:
//
//  1. For j = len(s) down to 0, test whether s[0..j) is a palindrome.
//  2. The first (largest) such j is the longest palindromic prefix.
//  3. Return reverse(s[j:]) + s.
//
// Time:  O(n²) — up to n prefix checks, each O(n).
// Space: O(n) for the reversed suffix and result.
func bruteForce(s string) string {
	n := len(s)
	for j := n; j >= 0; j-- { // longest prefix first
		if isPalindrome(s[:j]) { // s[0..j) already symmetric?
			return reverse(s[j:]) + s // mirror the non-palindromic tail to the front
		}
	}
	return s // unreachable (empty prefix is always a palindrome)
}

// isPalindrome reports whether t reads the same forwards and backwards.
func isPalindrome(t string) bool {
	for i, j := 0, len(t)-1; i < j; i, j = i+1, j-1 {
		if t[i] != t[j] {
			return false
		}
	}
	return true
}

// reverse returns t with its bytes in reverse order.
func reverse(t string) string {
	b := []byte(t)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i] // swap ends inward
	}
	return string(b)
}

// ── Approach 2: KMP Failure Function (Optimal) ───────────────────────────────
//
// kmp solves Shortest Palindrome in O(n) by building the KMP failure array of
// the string  s + '#' + reverse(s)  to find the longest palindromic prefix.
//
// Intuition:
//
//	The longest palindromic prefix of s equals the longest prefix of s that is
//	also a suffix of reverse(s). Concatenate  combined = s + sep + reverse(s)
//	with a separator that never appears in s (so matches can't cross it). The
//	final value of KMP's failure function over `combined` is exactly the length
//	of the longest prefix-of-s that is a suffix-of-reverse(s) — i.e. the length
//	k of the longest palindromic prefix. Then prepend reverse(s[k:]).
//
// Algorithm:
//
//  1. combined = s + "#" + reverse(s).
//  2. Build lps[] (longest proper prefix that is also suffix) for combined.
//  3. k = lps[last] = longest palindromic prefix length.
//  4. Return reverse(s[k:]) + s.
//
// Time:  O(n) — a single KMP failure-array build over a length-~2n string.
// Space: O(n) for the lps array and result.
func kmp(s string) string {
	if len(s) == 0 {
		return ""
	}
	rev := reverse(s)
	combined := s + "#" + rev // '#' cannot appear in lowercase s
	lps := buildLPS(combined)
	k := lps[len(lps)-1]      // length of longest palindromic prefix of s
	return rev[:len(s)-k] + s // reverse(s[k:]) == rev[:n-k]; prepend it
}

// buildLPS computes the KMP failure function: lps[i] = length of the longest
// proper prefix of t[:i+1] that is also a suffix of it.
func buildLPS(t string) []int {
	lps := make([]int, len(t))
	length := 0 // length of the current longest prefix-suffix
	for i := 1; i < len(t); i++ {
		for length > 0 && t[i] != t[length] {
			length = lps[length-1] // fall back along the failure links
		}
		if t[i] == t[length] {
			length++ // extend the matched prefix-suffix by one
		}
		lps[i] = length
	}
	return lps
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Longest Palindromic Prefix) ===")
	fmt.Printf("%q\n", bruteForce("aacecaaa")) // "aaacecaaa"
	fmt.Printf("%q\n", bruteForce("abcd"))     // "dcbabcd"
	fmt.Printf("%q\n", bruteForce(""))         // ""

	fmt.Println("=== Approach 2: KMP Failure Function (Optimal) ===")
	fmt.Printf("%q\n", kmp("aacecaaa")) // "aaacecaaa"
	fmt.Printf("%q\n", kmp("abcd"))     // "dcbabcd"
	fmt.Printf("%q\n", kmp(""))         // ""
}
