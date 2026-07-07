package main

import (
	"fmt"
	"strconv"
)

// ── Approach 1: String Conversion ────────────────────────────────────────────
//
// stringConversion converts x to a string and checks if the string reads the
// same forwards and backwards.
//
// Intuition:
//   The simplest definition of a palindrome applies directly to strings.
//   Convert and compare characters from both ends inward.
//
// Time:  O(log x) — the number of digits is O(log₁₀ x).
// Space: O(log x) — the string.
func stringConversion(x int) bool {
	if x < 0 {
		return false // negative numbers are never palindromes
	}
	s := strconv.Itoa(x)
	for l, r := 0, len(s)-1; l < r; l, r = l+1, r-1 {
		if s[l] != s[r] {
			return false
		}
	}
	return true
}

// ── Approach 2: Reverse Full Number ──────────────────────────────────────────
//
// reverseFullNumber reverses x entirely using integer arithmetic and compares
// the result to the original.
//
// Intuition:
//   Pop each digit from x and push onto reversed. If x == reversed at the end
//   the number is a palindrome. Negative numbers and numbers ending in 0
//   (unless the number is 0 itself) are immediately false.
//
// Time:  O(log x) — one iteration per digit.
// Space: O(1) — no extra allocation.
func reverseFullNumber(x int) bool {
	if x < 0 || (x%10 == 0 && x != 0) {
		return false
	}
	original := x
	reversed := 0
	for x > 0 {
		reversed = reversed*10 + x%10
		x /= 10
	}
	return original == reversed
}

// ── Approach 3: Reverse Only Half (Optimal) ──────────────────────────────────
//
// reverseHalf reverses only the second half of the number and compares it to
// the first half, avoiding a potential overflow when reversing the full number.
//
// Intuition:
//   We don't need to reverse the whole number. Stop when the reversed half
//   becomes >= remaining (we've reversed exactly half the digits).
//   For even-length palindromes: remaining == reversed.
//   For odd-length palindromes: remaining == reversed/10 (middle digit ignored).
//
// Time:  O(log x) — half the digits.
// Space: O(1).
func reverseHalf(x int) bool {
	// Negative, or positive ending in 0 (would need leading 0 to be palindrome).
	if x < 0 || (x%10 == 0 && x != 0) {
		return false
	}

	reversed := 0
	for x > reversed {
		reversed = reversed*10 + x%10
		x /= 10
	}

	// Even length: x == reversed.
	// Odd length:  x == reversed/10 (skip middle digit).
	return x == reversed || x == reversed/10
}

func main() {
	examples := []struct {
		x      int
		expect bool
	}{
		{121, true},
		{-121, false},
		{10, false},
		{0, true},
		{1221, true},
		{12321, true},
		{123, false},
	}

	approaches := []struct {
		name string
		fn   func(int) bool
	}{
		{"Approach 1: String Conversion       O(log x) T | O(log x) S", stringConversion},
		{"Approach 2: Reverse Full Number     O(log x) T | O(1)     S", reverseFullNumber},
		{"Approach 3: Reverse Half          ✅ O(log x) T | O(1)     S", reverseHalf},
	}

	for _, ex := range examples {
		fmt.Printf("x=%-8d  expect=%v\n", ex.x, ex.expect)
		for _, ap := range approaches {
			fmt.Printf("  %-62s → %v\n", ap.name, ap.fn(ex.x))
		}
		fmt.Println()
	}
}
