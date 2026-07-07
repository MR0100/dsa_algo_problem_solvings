package main

import "fmt"

// ── Approach 1: Stack ─────────────────────────────────────────────────────────
//
// stackApproach uses a stack to match closing brackets against the most
// recently opened bracket.
//
// Intuition:
//   Push every opening bracket onto the stack.
//   On a closing bracket, check if the stack's top is the matching opener.
//   If not (or the stack is empty), the string is invalid.
//   At the end, the stack must be empty (all openers were matched).
//
// Time:  O(n) — one pass over the string.
// Space: O(n) — the stack holds at most n/2 brackets.
func stackApproach(s string) bool {
	stack := make([]byte, 0, len(s)/2)
	// Maps each closer to its matching opener.
	match := map[byte]byte{')': '(', ']': '[', '}': '{'}

	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch == '(' || ch == '[' || ch == '{' {
			stack = append(stack, ch) // push opener
		} else {
			// Closing bracket: stack must have a matching opener on top.
			if len(stack) == 0 || stack[len(stack)-1] != match[ch] {
				return false
			}
			stack = stack[:len(stack)-1] // pop
		}
	}
	return len(stack) == 0
}

// ── Approach 2: Counter (single bracket type only — included for contrast) ───
//
// counterApproach is valid ONLY when the input contains a single bracket type
// (e.g. only parentheses). It uses a simple count rather than a stack.
//
// Included to show WHY a stack is necessary for mixed brackets: a counter
// cannot detect mismatches like "([)]".
//
// Time:  O(n).
// Space: O(1).
//
// Note: This approach gives WRONG answers for mixed brackets. It is included
//       as an educational contrast, not as a correct solution.
func counterApproach(s string) bool {
	count := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(', '[', '{':
			count++
		case ')', ']', '}':
			count--
			if count < 0 {
				return false
			}
		}
	}
	return count == 0
}

func main() {
	examples := []struct {
		s      string
		expect bool
	}{
		{"()", true},
		{"()[]{}", true},
		{"(]", false},
		{"([)]", false}, // counter gives wrong answer here
		{"{[]}", true},
		{"", true},
		{"((", false},
		{"]", false},
	}

	approaches := []struct {
		name string
		fn   func(string) bool
	}{
		{"Approach 1: Stack         ✅ O(n) T | O(n) S (correct for all)", stackApproach},
		{"Approach 2: Counter       ⚠️  O(n) T | O(1) S (single type only)", counterApproach},
	}

	for _, ex := range examples {
		fmt.Printf("s=%-10q  expect=%v\n", ex.s, ex.expect)
		for _, ap := range approaches {
			fmt.Printf("  %-60s → %v\n", ap.name, ap.fn(ex.s))
		}
		fmt.Println()
	}
}
