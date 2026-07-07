package main

import "fmt"

// Given a string s of upper/lowercase letters (case sensitive), return the LENGTH
// of the longest palindrome that can be assembled from its letters. A palindrome
// uses each character an even number of times, plus at most one character an odd
// number of times (the single center).

// ── Approach 1: Hash Map Frequency ───────────────────────────────────────────
//
// hashMap solves Longest Palindrome by counting each character with a map and
// summing how many of each can face itself across the palindrome's center.
//
// Intuition:
//
//	A palindrome mirrors around its center, so every character except possibly
//	the center one must appear an even number of times. For a character seen c
//	times we can place c rounded down to even = c - (c&1) of them (matched pairs).
//	If ANY character has an odd count, exactly one leftover single may sit at the
//	center, adding 1. So answer = (sum of even-floored counts) + (1 if any odd
//	count exists else 0).
//
// Algorithm:
//  1. Count occurrences of each byte in a map.
//  2. length = 0; hasOdd = false.
//  3. For each count c: add c - (c%2) to length; if c is odd set hasOdd = true.
//  4. If hasOdd, add 1 (one character goes in the middle).
//
// Time:  O(n) — one pass to count, one over the ≤ distinct-char buckets.
// Space: O(k) — k distinct characters (≤ 52 for letters).
func hashMap(s string) int {
	counts := make(map[byte]int) // character → how many times it appears
	for i := 0; i < len(s); i++ {
		counts[s[i]]++ // tally this character
	}

	length := 0     // running palindrome length
	hasOdd := false // did we see any character with an odd count?
	for _, c := range counts {
		length += c - c%2 // use the largest even number of this char (full pairs)
		if c%2 == 1 {
			hasOdd = true // at least one odd-count char exists → a center is available
		}
	}
	if hasOdd {
		length++ // one leftover single character can occupy the exact middle
	}
	return length
}

// ── Approach 2: Fixed Count Array + Odd Tally (Optimal) ───────────────────────
//
// countArray solves Longest Palindrome using a 128-slot ASCII count array and a
// single accumulator of how many characters have an odd count.
//
// Intuition:
//
//	Same parity insight, but replace the map with a flat array indexed by ASCII
//	code (letters are a small fixed alphabet) for constant-factor speed and O(1)
//	extra space. Track oddCount = number of characters currently appearing an odd
//	number of times, updated on the fly. Every odd count wastes exactly one
//	character (the unpaired one), so the answer is n minus the wasted singles,
//	except we get ONE of them back for free as the palindrome's center. Hence:
//	if oddCount > 0, answer = n - oddCount + 1; else answer = n.
//
// Algorithm:
//  1. count[128]; for each char increment its slot and flip oddCount: +1 when the
//     slot becomes odd, -1 when it becomes even.
//  2. If oddCount > 0, return n - oddCount + 1 (keep one odd char as center,
//     discard the other oddCount-1 unpaired singles).
//  3. Else return n (everything pairs perfectly).
//
// Time:  O(n) — a single pass over the string.
// Space: O(1) — a fixed 128-entry array regardless of input size.
func countArray(s string) int {
	var count [128]int // ASCII code → occurrence count (letters fit in 0..127)
	oddCount := 0      // how many characters currently have an odd count
	for i := 0; i < len(s); i++ {
		count[s[i]]++ // record this character
		if count[s[i]]%2 == 1 {
			oddCount++ // count for this char just turned odd
		} else {
			oddCount-- // it turned even again — this char is fully paired now
		}
	}

	n := len(s)
	if oddCount > 0 {
		// Each odd-count char contributes one unpaired single; we may keep exactly
		// one of them in the center, so we discard (oddCount - 1) characters.
		return n - oddCount + 1
	}
	return n // all counts even → the whole string is usable
}

func main() {
	fmt.Println("=== Approach 1: Hash Map Frequency ===")
	fmt.Printf("s=\"abccccdd\" got=%d  expected 7\n", hashMap("abccccdd"))
	fmt.Printf("s=\"a\"        got=%d  expected 1\n", hashMap("a"))
	fmt.Printf("s=\"bb\"       got=%d  expected 2\n", hashMap("bb"))
	fmt.Printf("s=\"Aa\"       got=%d  expected 1\n", hashMap("Aa")) // case sensitive

	fmt.Println("=== Approach 2: Fixed Count Array + Odd Tally (Optimal) ===")
	fmt.Printf("s=\"abccccdd\" got=%d  expected 7\n", countArray("abccccdd"))
	fmt.Printf("s=\"a\"        got=%d  expected 1\n", countArray("a"))
	fmt.Printf("s=\"bb\"       got=%d  expected 2\n", countArray("bb"))
	fmt.Printf("s=\"Aa\"       got=%d  expected 1\n", countArray("Aa"))
	fmt.Printf("s=\"ccc\"      got=%d  expected 3\n", countArray("ccc")) // one center + a pair
}
