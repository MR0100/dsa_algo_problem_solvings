package main

import (
	"fmt"
	"strings"
)

// ── Approach 1: Simulate with Row Buffers ────────────────────────────────────
//
// simulate distributes each character into one of numRows string builders,
// using a "bouncing" direction flag to zig-zag down then back up.
//
// Intuition:
//   Write the string diagonally down (direction +1) until hitting the last
//   row, then diagonally up (direction -1) until hitting row 0, then repeat.
//   Each character goes into the buffer for the current row; concatenate all
//   row buffers at the end.
//
// Time:  O(n) — one pass over s plus O(n) concatenation.
// Space: O(n) — the row buffers hold all n characters.
func simulate(s string, numRows int) string {
	if numRows == 1 || numRows >= len(s) {
		return s
	}

	rows := make([]strings.Builder, numRows)
	curRow := 0
	goingDown := false

	for _, ch := range s {
		rows[curRow].WriteRune(ch)
		// Bounce direction at the top and bottom rows.
		if curRow == 0 || curRow == numRows-1 {
			goingDown = !goingDown
		}
		if goingDown {
			curRow++
		} else {
			curRow--
		}
	}

	var result strings.Builder
	for _, row := range rows {
		result.WriteString(row.String())
	}
	return result.String()
}

// ── Approach 2: Math — Direct Index Formula ──────────────────────────────────
//
// mathFormula reads characters in output order using the cycle-length formula,
// without physically distributing them into rows first.
//
// Intuition:
//   The pattern repeats every cycleLen = 2*(numRows-1) characters.
//   For each row r and each cycle start j (0, cycleLen, 2*cycleLen, ...):
//     Going-down character:  s[j + r]
//     Going-up  character:   s[j + cycleLen - r]  (only for 0 < r < numRows-1)
//   Reading them in this order directly produces the zigzag output.
//
// Time:  O(n) — every character visited once.
// Space: O(n) — the result builder.
func mathFormula(s string, numRows int) string {
	if numRows == 1 || numRows >= len(s) {
		return s
	}

	n := len(s)
	cycleLen := 2 * (numRows - 1)
	var result strings.Builder

	for row := 0; row < numRows; row++ {
		for j := 0; j+row < n; j += cycleLen {
			result.WriteByte(s[j+row]) // down-stroke character

			// Up-stroke character only for interior rows (not top/bottom).
			upIdx := j + cycleLen - row
			if row != 0 && row != numRows-1 && upIdx < n {
				result.WriteByte(s[upIdx])
			}
		}
	}
	return result.String()
}

func main() {
	examples := []struct {
		s       string
		numRows int
		expect  string
	}{
		{"PAYPALISHIRING", 3, "PAHNAPLSIIGYIR"},
		{"PAYPALISHIRING", 4, "PINALSIGYAHRPI"},
		{"A", 1, "A"},
		{"AB", 1, "AB"},
	}

	approaches := []struct {
		name string
		fn   func(string, int) string
	}{
		{"Approach 1: Simulate (row buffers) O(n) T | O(n) S", simulate},
		{"Approach 2: Math formula         ✅ O(n) T | O(n) S", mathFormula},
	}

	for _, ex := range examples {
		fmt.Printf("s=%q  numRows=%d  expect=%q\n", ex.s, ex.numRows, ex.expect)
		for _, ap := range approaches {
			fmt.Printf("  %-55s → %q\n", ap.name, ap.fn(ex.s, ex.numRows))
		}
		fmt.Println()
	}
}
