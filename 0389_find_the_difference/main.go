package main

import "fmt"

// ── Approach 1: Frequency Count (Hash / 26-array) ────────────────────────────
//
// countArray solves Find the Difference by counting each letter in t and
// subtracting the count of each letter in s; the letter left with a positive
// count is the extra one.
//
// Intuition:
//
//	t is s shuffled with exactly one extra letter inserted. If we add 1 for
//	every letter of t and subtract 1 for every letter of s, all the shared
//	letters cancel to 0 and only the added letter ends at count 1.
//
// Algorithm:
//  1. count[t[i]-'a']++ for all i.
//  2. count[s[i]-'a']-- for all i.
//  3. Return the single letter whose count is > 0.
//
// Time:  O(n) — two linear passes.
// Space: O(1) — 26 counters.
func countArray(s string, t string) byte {
	var count [26]int
	for i := 0; i < len(t); i++ {
		count[t[i]-'a']++ // t is one char longer; tally it up
	}
	for i := 0; i < len(s); i++ {
		count[s[i]-'a']-- // cancel out every letter that also appears in s
	}
	for i := 0; i < 26; i++ {
		if count[i] > 0 { // the leftover positive count is the added letter
			return byte('a' + i)
		}
	}
	return 0 // unreachable given the problem guarantees
}

// ── Approach 2: XOR Bit Manipulation (Optimal) ───────────────────────────────
//
// xorBits solves Find the Difference by XOR-ing every character of both
// strings together; identical characters cancel and the odd one out remains.
//
// Intuition:
//
//	XOR is its own inverse: x ^ x == 0, and it is commutative/associative. Each
//	letter present in both s and t appears an even number of times across the
//	combined stream, so it cancels to 0. The single extra letter appears an odd
//	number of times and survives as the XOR of everything.
//
// Algorithm:
//  1. acc = 0.
//  2. XOR acc with every byte of s and every byte of t.
//  3. Return acc — the surviving character.
//
// Time:  O(n) — one pass over each string.
// Space: O(1) — a single accumulator byte.
func xorBits(s string, t string) byte {
	var acc byte // XOR accumulator
	for i := 0; i < len(s); i++ {
		acc ^= s[i] // fold in each char of s
	}
	for i := 0; i < len(t); i++ {
		acc ^= t[i] // fold in each char of t; matched pairs cancel
	}
	return acc // odd one out remains
}

// ── Approach 3: ASCII Sum Difference ─────────────────────────────────────────
//
// sumDiff solves Find the Difference by summing the character codes of t and
// subtracting the sum of s; the difference is exactly the added character's code.
//
// Intuition:
//
//	Every shared letter contributes its ASCII value to both sums, so those
//	contributions cancel in (sum of t) − (sum of s). What remains is precisely
//	the code of the one extra character. Simple integer arithmetic, no tables.
//
// Algorithm:
//  1. diff = 0. Add each t[i] to diff; subtract each s[i] from diff.
//  2. Return byte(diff).
//
// Time:  O(n) — one pass over each string.
// Space: O(1) — a single integer accumulator.
func sumDiff(s string, t string) byte {
	diff := 0 // (sum of t codes) − (sum of s codes)
	for i := 0; i < len(t); i++ {
		diff += int(t[i]) // add every code of t
	}
	for i := 0; i < len(s); i++ {
		diff -= int(s[i]) // remove every code of s
	}
	return byte(diff) // leftover = code of the extra letter
}

func main() {
	fmt.Println("=== Approach 1: Frequency Count ===")
	fmt.Printf("%c\n", countArray("abcd", "abcde")) // expected e
	fmt.Printf("%c\n", countArray("", "y"))         // expected y

	fmt.Println("=== Approach 2: XOR Bit Manipulation (Optimal) ===")
	fmt.Printf("%c\n", xorBits("abcd", "abcde")) // expected e
	fmt.Printf("%c\n", xorBits("", "y"))         // expected y

	fmt.Println("=== Approach 3: ASCII Sum Difference ===")
	fmt.Printf("%c\n", sumDiff("abcd", "abcde")) // expected e
	fmt.Printf("%c\n", sumDiff("", "y"))         // expected y
}
