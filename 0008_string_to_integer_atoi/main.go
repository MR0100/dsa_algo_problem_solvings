package main

import (
	"fmt"
	"math"
	"strconv"
)

// ── Approach 1: Standard Library ─────────────────────────────────────────────
//
// useStdlib delegates to strconv.Atoi / strconv.ParseInt after manually
// stripping the leading whitespace and clipping overflow.
//
// Intuition:
//   Trim leading whitespace, call strconv.ParseInt with bitSize=32, then
//   map the range error to INT32_MIN / INT32_MAX. Handles the sign and digit
//   parsing for us; does NOT handle the "stop at first non-digit" rule exactly
//   the same way, so it is an approximation.
//
// Note: strconv.ParseInt stops at the first non-digit, which is correct,
//       but it returns an error for empty strings after trimming sign — we
//       guard that. This approach is useful as a reference only; the manual
//       approach below is the interview answer.
//
// Time:  O(n) — linear scan of the string.
// Space: O(1).
func useStdlib(s string) int {
	// Trim leading whitespace.
	i := 0
	for i < len(s) && s[i] == ' ' {
		i++
	}
	s = s[i:]

	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		// Range error: clamp.
		if numErr, ok := err.(*strconv.NumError); ok && numErr.Err == strconv.ErrRange {
			if len(s) > 0 && s[0] == '-' {
				return math.MinInt32
			}
			return math.MaxInt32
		}
		// Syntax error: partial parse — strconv doesn't do partial parse like atoi.
		// Fall back to 0 for this approximation.
		return 0
	}
	if val > math.MaxInt32 {
		return math.MaxInt32
	}
	if val < math.MinInt32 {
		return math.MinInt32
	}
	return int(val)
}

// ── Approach 2: Manual Linear Scan (Correct atoi) ─────────────────────────
//
// manualScan implements the exact myAtoi specification:
//   1. Skip leading whitespace.
//   2. Read optional '+' or '-'.
//   3. Read decimal digits until non-digit or end; ignore leading zeros.
//   4. Clamp to [INT32_MIN, INT32_MAX] on overflow.
//
// Intuition:
//   A straightforward state machine: whitespace → sign → digits → done.
//   Accumulate result as int64 so we can detect 32-bit overflow cleanly.
//
// Time:  O(n) — one pass.
// Space: O(1) — only counters.
func manualScan(s string) int {
	i, n := 0, len(s)

	// Step 1: skip leading whitespace.
	for i < n && s[i] == ' ' {
		i++
	}
	if i == n {
		return 0
	}

	// Step 2: read sign.
	sign := 1
	if s[i] == '+' {
		i++
	} else if s[i] == '-' {
		sign = -1
		i++
	}

	// Step 3: read digits, accumulate in int64 for overflow detection.
	var result int64
	for i < n && s[i] >= '0' && s[i] <= '9' {
		digit := int64(s[i] - '0')
		result = result*10 + digit

		// Step 4: clamp early if we've already exceeded INT32.
		if sign == 1 && result > math.MaxInt32 {
			return math.MaxInt32
		}
		if sign == -1 && -result < math.MinInt32 {
			return math.MinInt32
		}
		i++
	}

	return sign * int(result)
}

func main() {
	examples := []struct {
		s      string
		expect int
	}{
		{"42", 42},
		{"   -042", -42},
		{"1337c0d3", 1337},
		{"0-1", 0},
		{"words and 987", 0},
		{"-91283472332", math.MinInt32}, // overflow
		{"2147483648", math.MaxInt32},   // overflow
		{"  +  413", 0},                 // sign then space → stop at space
	}

	approaches := []struct {
		name string
		fn   func(string) int
	}{
		{"Approach 1: stdlib ParseInt  (approx) O(n) T | O(1) S", useStdlib},
		{"Approach 2: Manual scan      ✅ (exact) O(n) T | O(1) S", manualScan},
	}

	for _, ex := range examples {
		fmt.Printf("s=%-22q  expect=%d\n", ex.s, ex.expect)
		for _, ap := range approaches {
			fmt.Printf("  %-57s → %d\n", ap.name, ap.fn(ex.s))
		}
		fmt.Println()
	}
}
