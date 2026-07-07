package main

import (
	"fmt"
	"math"
	"strconv"
)

// ── Approach 1: String Conversion ────────────────────────────────────────────
//
// stringConversion converts x to a string, reverses it, then parses it back.
//
// Intuition:
//   The simplest way: convert to string, handle sign separately, reverse the
//   digit characters, re-parse. Use strconv.ParseInt for overflow detection.
//
// Time:  O(log x) — the number of digits is O(log₁₀ x).
// Space: O(log x) — the string copy.
func stringConversion(x int) int {
	if x == 0 {
		return 0
	}

	sign := 1
	if x < 0 {
		sign = -1
		x = -x
	}

	s := strconv.Itoa(x)
	// Reverse the string.
	runes := []byte(s)
	for l, r := 0, len(runes)-1; l < r; l, r = l+1, r-1 {
		runes[l], runes[r] = runes[r], runes[l]
	}

	// Parse back; catch overflow.
	val, err := strconv.ParseInt(string(runes), 10, 64)
	if err != nil || val > math.MaxInt32 {
		return 0
	}
	return sign * int(val)
}

// ── Approach 2: Math — Pop and Push Digits ───────────────────────────────────
//
// mathPopPush reverses x by repeatedly popping the last digit and pushing it
// onto a result, checking for 32-bit overflow before each push.
//
// Intuition:
//   Pop: digit = x % 10, x /= 10.
//   Push: result = result * 10 + digit.
//   Before pushing, check that result won't exceed INT32_MAX / INT32_MIN.
//   The overflow check uses INT32_MAX/10 and INT32_MIN/10 as thresholds.
//
// Time:  O(log x) — one iteration per digit.
// Space: O(1) — no extra allocation.
func mathPopPush(x int) int {
	result := 0

	for x != 0 {
		digit := x % 10
		x /= 10

		// Check overflow before multiplying result by 10 and adding digit.
		// INT32_MAX = 2147483647, INT32_MIN = -2147483648
		if result > math.MaxInt32/10 || (result == math.MaxInt32/10 && digit > 7) {
			return 0
		}
		if result < math.MinInt32/10 || (result == math.MinInt32/10 && digit < -8) {
			return 0
		}

		result = result*10 + digit
	}
	return result
}

func main() {
	examples := []struct {
		x      int
		expect int
	}{
		{123, 321},
		{-123, -321},
		{120, 21},
		{0, 0},
		{1534236469, 0},  // overflows INT32
		{-2147483648, 0}, // -INT32_MIN overflows
	}

	approaches := []struct {
		name string
		fn   func(int) int
	}{
		{"Approach 1: String Conversion     O(log x) T | O(log x) S", stringConversion},
		{"Approach 2: Math Pop/Push       ✅ O(log x) T | O(1)     S", mathPopPush},
	}

	for _, ex := range examples {
		fmt.Printf("x=%-14d  expect=%d\n", ex.x, ex.expect)
		for _, ap := range approaches {
			fmt.Printf("  %-55s → %d\n", ap.name, ap.fn(ex.x))
		}
		fmt.Println()
	}
}
