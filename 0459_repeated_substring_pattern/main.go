package main

import (
	"fmt"
	"strings"
)

// ── Approach 1: Try Every Divisor Period (Brute Force) ───────────────────────
//
// divisorBruteForce solves Repeated Substring Pattern by testing every candidate
// period length that divides n and checking whether repeating that prefix rebuilds s.
//
// Intuition:
//
//	If s is k copies of a block of length L, then L divides n and every
//	character equals the one L positions earlier: s[i] == s[i-L]. A valid block
//	length must be a proper divisor of n (L < n, n % L == 0). Try each such L
//	from smallest up; for each, verify the tiling condition across the whole
//	string. The first L that passes proves s is periodic.
//
// Algorithm:
//  1. n = len(s). For L = 1 .. n/2:
//     a. Skip L if n % L != 0 (block must tile the string exactly).
//     b. Check s[i] == s[i-L] for all i in [L, n); if all match, return true.
//  2. If no divisor works, return false.
//
// Time:  O(n * d(n)) ⊆ O(n^2) — for each of the ≤ n/2 divisors, an O(n) scan.
// Space: O(1) — index arithmetic only.
func divisorBruteForce(s string) bool {
	n := len(s)
	for l := 1; l <= n/2; l++ { // candidate block length; a real period is ≤ n/2
		if n%l != 0 {
			continue // block of length l can't tile n characters evenly
		}
		ok := true
		// Verify periodicity: each char must equal the one one block earlier.
		for i := l; i < n; i++ {
			if s[i] != s[i-l] {
				ok = false // mismatch → l is not a valid period
				break
			}
		}
		if ok {
			return true // s is the length-l block repeated n/l times
		}
	}
	return false // no proper divisor period rebuilds s
}

// ── Approach 2: Concatenation Trick (s+s Doubling) ───────────────────────────
//
// concatTrick solves Repeated Substring Pattern with the identity: s is periodic
// iff s appears inside (s+s) with the first and last characters removed.
//
// Intuition:
//
//	Build t = s + s. If s is k ≥ 2 copies of a block, then s starts again at
//	offset L (the block length) inside t, i.e. somewhere in position 1..n-1 of
//	t. Stripping the first and last character of t deletes the two trivial
//	occurrences (at offset 0 and offset n), so s is found in t[1:2n-1] IFF a
//	non-trivial period exists. One substring search settles it.
//
// Algorithm:
//  1. Form doubled = s + s.
//  2. Remove the first and last characters: middle = doubled[1 : len(doubled)-1].
//  3. Return whether middle contains s.
//
// Time:  O(n^2) with a naive substring search; O(n) if the search uses KMP/Z
//
//	(Go's strings.Contains is a tuned mix and is effectively linear here).
//
// Space: O(n) — the doubled string.
func concatTrick(s string) bool {
	doubled := s + s                      // two back-to-back copies
	middle := doubled[1 : len(doubled)-1] // drop the trivial offset-0 and offset-n matches
	return strings.Contains(middle, s)    // a surviving match ⇒ a real internal period
}

// ── Approach 3: KMP Failure Function Period (Optimal) ────────────────────────
//
// kmpFailure solves Repeated Substring Pattern by computing the longest
// proper-prefix-that-is-also-suffix (LPS) and checking the induced period.
//
// Intuition:
//
//	Let lps[n-1] be the length of the longest proper prefix of s that is also a
//	suffix. Then p = n - lps[n-1] is the smallest period of s. The string is a
//	repetition of a smaller block IFF that period p divides n AND p < n (so
//	lps[n-1] > 0). This computes the answer with a single O(n) prefix-function
//	pass and no string doubling.
//
// Algorithm:
//  1. Build lps[] (KMP prefix function) for s in O(n).
//  2. Let last = lps[n-1] and p = n - last.
//  3. Return last > 0 AND n % p == 0.
//
// Time:  O(n) — one prefix-function computation.
// Space: O(n) — the lps array.
func kmpFailure(s string) bool {
	n := len(s)
	if n < 2 {
		return false // a single character can't be a repeat of a shorter block
	}
	lps := make([]int, n) // lps[i] = length of longest proper prefix==suffix of s[0..i]
	length := 0           // length of the current matching prefix
	// Standard KMP prefix-function build.
	for i := 1; i < n; i++ {
		// Fall back through shorter borders while characters disagree.
		for length > 0 && s[i] != s[length] {
			length = lps[length-1] // reuse the next-longest border
		}
		if s[i] == s[length] {
			length++ // extend the current border by one character
		}
		lps[i] = length // record the border length ending at i
	}
	last := lps[n-1] // longest border of the whole string
	period := n - last
	// Periodic iff a non-empty border exists and its induced period tiles n.
	return last > 0 && n%period == 0
}

func main() {
	fmt.Println("=== Approach 1: Try Every Divisor Period (Brute Force) ===")
	fmt.Println(divisorBruteForce("abab"))         // expected true
	fmt.Println(divisorBruteForce("aba"))          // expected false
	fmt.Println(divisorBruteForce("abcabcabcabc")) // expected true

	fmt.Println("=== Approach 2: Concatenation Trick (s+s Doubling) ===")
	fmt.Println(concatTrick("abab"))         // expected true
	fmt.Println(concatTrick("aba"))          // expected false
	fmt.Println(concatTrick("abcabcabcabc")) // expected true

	fmt.Println("=== Approach 3: KMP Failure Function Period (Optimal) ===")
	fmt.Println(kmpFailure("abab"))         // expected true
	fmt.Println(kmpFailure("aba"))          // expected false
	fmt.Println(kmpFailure("abcabcabcabc")) // expected true
}
