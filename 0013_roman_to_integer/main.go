package main

import "fmt"

// ── Approach 1: Left-to-Right with Lookahead ─────────────────────────────────
//
// leftToRightLookahead scans the Roman numeral left to right. When a symbol is
// followed by a larger symbol, it is subtractive (subtract); otherwise additive.
//
// Intuition:
//   In any valid Roman numeral, a smaller value before a larger value means
//   subtraction (e.g. IV = 5-4 = 4, IX = 10-1 = 9). Build a value map and
//   check s[i] vs s[i+1] at each position.
//
// Time:  O(n) — single pass.
// Space: O(1) — fixed map with 7 entries.
func leftToRightLookahead(s string) int {
	val := map[byte]int{
		'I': 1, 'V': 5, 'X': 10, 'L': 50,
		'C': 100, 'D': 500, 'M': 1000,
	}

	result := 0
	for i := 0; i < len(s); i++ {
		cur := val[s[i]]
		// If there's a next symbol and it's larger, subtract current.
		if i+1 < len(s) && val[s[i+1]] > cur {
			result -= cur
		} else {
			result += cur
		}
	}
	return result
}

// ── Approach 2: Right-to-Left Running Total ───────────────────────────────────
//
// rightToLeft scans from right to left, maintaining a running "previous" value.
// If the current value is less than prev, subtract; otherwise add.
//
// Intuition:
//   Scanning right to left means we've already seen the "larger" symbol before
//   the "smaller" one. If current < prev, the current symbol is part of a
//   subtractive pair, so subtract it.
//
// Time:  O(n).
// Space: O(1).
func rightToLeft(s string) int {
	val := map[byte]int{
		'I': 1, 'V': 5, 'X': 10, 'L': 50,
		'C': 100, 'D': 500, 'M': 1000,
	}

	result := 0
	prev := 0
	for i := len(s) - 1; i >= 0; i-- {
		cur := val[s[i]]
		if cur < prev {
			result -= cur // subtractive position
		} else {
			result += cur
		}
		prev = cur
	}
	return result
}

func main() {
	examples := []struct {
		s      string
		expect int
	}{
		{"III", 3},
		{"LVIII", 58},
		{"MCMXCIV", 1994},
		{"IV", 4},
		{"IX", 9},
		{"XL", 40},
		{"XC", 90},
		{"CD", 400},
		{"CM", 900},
		{"MMMCMXCIX", 3999},
	}

	approaches := []struct {
		name string
		fn   func(string) int
	}{
		{"Approach 1: Left-to-Right Lookahead  O(n) T | O(1) S", leftToRightLookahead},
		{"Approach 2: Right-to-Left          ✅ O(n) T | O(1) S", rightToLeft},
	}

	for _, ex := range examples {
		fmt.Printf("s=%-12q  expect=%d\n", ex.s, ex.expect)
		for _, ap := range approaches {
			fmt.Printf("  %-54s → %d\n", ap.name, ap.fn(ex.s))
		}
		fmt.Println()
	}
}
