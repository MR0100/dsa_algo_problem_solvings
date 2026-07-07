package main

import (
	"fmt"
	"strings"
)

// Integer to English Words (LeetCode #273)
//
// Convert a non-negative integer num to its English words representation.
// 0 <= num <= 2^31 - 1 = 2,147,483,647, so the answer never exceeds "Two
// Billion ...".

// Shared lookup tables used by both approaches.
var belowTwenty = []string{
	"", "One", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine",
	"Ten", "Eleven", "Twelve", "Thirteen", "Fourteen", "Fifteen", "Sixteen",
	"Seventeen", "Eighteen", "Nineteen",
}
var tens = []string{
	"", "", "Twenty", "Thirty", "Forty", "Fifty", "Sixty", "Seventy",
	"Eighty", "Ninety",
}

// ── Approach 1: Recursive Divide by Scale (Optimal) ──────────────────────────
//
// recursive builds the words by recursively decomposing num against descending
// scales: Billion (1e9), Million (1e6), Thousand (1e3), and the sub-1000 base.
//
// Intuition:
//
//	English number words are structured in groups of three digits, each group
//	read the same way ("one hundred twenty three") and then tagged with a
//	scale word (Thousand / Million / Billion). Recurse: the words for N are
//	"words(N / scale) + scaleWord + words(N % scale)" for the largest scale
//	that fits; below 1000 we spell out hundreds, tens, and ones directly.
//
// Algorithm:
//  1. Base case: num == 0 → "" (the caller trims / substitutes "Zero").
//  2. num < 20: return belowTwenty[num].
//  3. num < 100: tens[num/10] + " " + words(num%10).
//  4. num < 1000: words(num/100) + " Hundred " + words(num%100).
//  5. Else, for scale in {1e9→"Billion", 1e6→"Million", 1e3→"Thousand"}: if
//     num >= scale, return words(num/scale) + " " + scaleWord + " " + words(num%scale).
//  6. Trim stray spaces at the top level; special-case 0 → "Zero".
//
// Time:  O(1) — the recursion depth and work are bounded by the fixed number of
//
//	scales (num < 2^31 has at most 10 digits / 4 groups).
//
// Space: O(1) — bounded recursion; output length is bounded by a constant.
func recursive(num int) string {
	if num == 0 {
		return "Zero" // only reachable at the top level; helpers return "" for 0
	}
	return strings.TrimSpace(recursiveHelper(num))
}

// recursiveHelper returns the (possibly space-padded) words for num, using ""
// for zero so parents can concatenate cleanly.
func recursiveHelper(num int) string {
	switch {
	case num == 0:
		return "" // contributes nothing to its parent group
	case num < 20:
		return belowTwenty[num] // direct table hit for 1..19
	case num < 100:
		// tens digit word, then recurse on the ones digit (which may be 0 → "").
		return tens[num/10] + " " + recursiveHelper(num%10)
	case num < 1000:
		// hundreds digit, the literal "Hundred", then the remaining two digits.
		return recursiveHelper(num/100) + " Hundred " + recursiveHelper(num%100)
	case num < 1000000:
		return recursiveHelper(num/1000) + " Thousand " + recursiveHelper(num%1000)
	case num < 1000000000:
		return recursiveHelper(num/1000000) + " Million " + recursiveHelper(num%1000000)
	default:
		return recursiveHelper(num/1000000000) + " Billion " + recursiveHelper(num%1000000000)
	}
}

// ── Approach 2: Iterative Grouping by Thousands ──────────────────────────────
//
// iterativeGroups splits num into groups of three digits from the least
// significant end, converts each group, and tags it with its scale word.
//
// Intuition:
//
//	Rather than recurse, peel the number 1000 at a time: the lowest three
//	digits get no scale word, the next three get "Thousand", then "Million",
//	then "Billion". Convert each non-zero group with a helper that spells a
//	value in [1, 999], prepend the scale word, and stitch groups together in
//	most-significant-first order.
//
// Algorithm:
//  1. If num == 0, return "Zero".
//  2. scales = ["", "Thousand", "Million", "Billion"]; i = 0; parts = [].
//  3. While num > 0: g = num % 1000; if g != 0, prepend
//     "three(g) + scales[i]" to parts; num /= 1000; i++.
//  4. Join parts with single spaces and trim.
//
// Time:  O(1) — at most 4 groups, each O(1). Space: O(1).
func iterativeGroups(num int) string {
	if num == 0 {
		return "Zero"
	}
	scales := []string{"", "Thousand", "Million", "Billion"}
	var parts []string // most-significant group ends up first
	i := 0
	for num > 0 {
		if g := num % 1000; g != 0 {
			// Convert this three-digit group and tag with its scale word.
			chunk := strings.TrimSpace(three(g) + " " + scales[i])
			// Prepend so higher scales precede lower ones in the final order.
			parts = append([]string{chunk}, parts...)
		}
		num /= 1000 // move to the next higher three-digit group
		i++
	}
	return strings.Join(parts, " ")
}

// three spells a value in [1, 999] as English words (no scale word, no leading
// or trailing spaces).
func three(n int) string {
	var b []string
	if n >= 100 {
		// hundreds digit + literal "Hundred"
		b = append(b, belowTwenty[n/100], "Hundred")
		n %= 100
	}
	if n >= 20 {
		// tens word (e.g. "Forty"), drop to the ones digit
		b = append(b, tens[n/10])
		n %= 10
	}
	if n > 0 {
		// remaining 1..19 as a single word
		b = append(b, belowTwenty[n])
	}
	return strings.Join(b, " ")
}

func main() {
	fmt.Println("=== Approach 1: Recursive Divide by Scale ===")
	fmt.Println(recursive(123))        // expected One Hundred Twenty Three
	fmt.Println(recursive(12345))      // expected Twelve Thousand Three Hundred Forty Five
	fmt.Println(recursive(1234567))    // expected One Million Two Hundred Thirty Four Thousand Five Hundred Sixty Seven
	fmt.Println(recursive(0))          // expected Zero
	fmt.Println(recursive(2147483647)) // expected Two Billion One Hundred Forty Seven Million Four Hundred Eighty Three Thousand Six Hundred Forty Seven

	fmt.Println("=== Approach 2: Iterative Grouping by Thousands ===")
	fmt.Println(iterativeGroups(123))        // expected One Hundred Twenty Three
	fmt.Println(iterativeGroups(12345))      // expected Twelve Thousand Three Hundred Forty Five
	fmt.Println(iterativeGroups(1234567))    // expected One Million Two Hundred Thirty Four Thousand Five Hundred Sixty Seven
	fmt.Println(iterativeGroups(0))          // expected Zero
	fmt.Println(iterativeGroups(2147483647)) // expected Two Billion One Hundred Forty Seven Million Four Hundred Eighty Three Thousand Six Hundred Forty Seven
}
