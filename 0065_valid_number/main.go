package main

import "fmt"

// ── Approach 1: Deterministic Finite Automaton (DFA) ─────────────────────────
//
// isNumberDFA solves Valid Number using a state machine that processes
// each character and transitions between states.
//
// Intuition:
//   A valid number matches the pattern:
//     [sign] (integer | decimal) [exponent]
//   where:
//     integer = digit+
//     decimal = digit* '.' digit+  OR  digit+ '.'  OR  digit+ '.' digit+
//     exponent = ('e'|'E') [sign] digit+
//
//   We model this as a DFA with states:
//     0: initial
//     1: sign before digits
//     2: digits before decimal point
//     3: decimal point (no digits yet before it)
//     4: digits after decimal point
//     5: decimal point after digits (e.g. "1.")
//     6: 'e'/'E' seen
//     7: sign after 'e'
//     8: digits after 'e'
//   Accept states: 2, 4, 5, 8 (states where a valid number can end)
//
// Time:  O(n) — one pass.
// Space: O(1)
func isNumberDFA(s string) bool {
	// transitions[state][charType] = nextState, -1 = invalid
	// charType: 0=digit, 1=sign(+/-), 2=dot, 3=e/E
	transitions := [9][4]int{
		/*0*/ {2, 1, 3, -1},
		/*1*/ {2, -1, 3, -1},
		/*2*/ {2, -1, 5, 6},
		/*3*/ {4, -1, -1, -1},
		/*4*/ {4, -1, -1, 6},
		/*5*/ {4, -1, -1, 6},
		/*6*/ {8, 7, -1, -1},
		/*7*/ {8, -1, -1, -1},
		/*8*/ {8, -1, -1, -1},
	}
	acceptStates := map[int]bool{2: true, 4: true, 5: true, 8: true}

	state := 0
	for _, ch := range s {
		var ct int
		switch {
		case ch >= '0' && ch <= '9':
			ct = 0
		case ch == '+' || ch == '-':
			ct = 1
		case ch == '.':
			ct = 2
		case ch == 'e' || ch == 'E':
			ct = 3
		default:
			return false // invalid character
		}
		state = transitions[state][ct]
		if state == -1 {
			return false
		}
	}
	return acceptStates[state]
}

// ── Approach 2: Manual Parsing ────────────────────────────────────────────────
//
// isNumberParse solves Valid Number by parsing the string manually according to
// the grammar:
//   valid number = (integer | decimal) optional-exponent
//
// Intuition:
//   Scan past optional leading sign, then check:
//   - Did we see any digit?
//   - Did we see a dot? If so, did we see any digit before or after?
//   - If 'e'/'E' found, require at least one digit before and after.
//
// Time:  O(n)
// Space: O(1)
func isNumberParse(s string) bool {
	i, n := 0, len(s)
	if n == 0 {
		return false
	}
	// optional leading sign
	if s[i] == '+' || s[i] == '-' {
		i++
	}

	seenDigit := false
	seenDot := false

	for i < n && s[i] != 'e' && s[i] != 'E' {
		if s[i] >= '0' && s[i] <= '9' {
			seenDigit = true
		} else if s[i] == '.' && !seenDot {
			seenDot = true
		} else {
			return false
		}
		i++
	}

	if !seenDigit {
		return false // must have at least one digit in mantissa
	}

	// optional exponent
	if i < n && (s[i] == 'e' || s[i] == 'E') {
		i++
		if i < n && (s[i] == '+' || s[i] == '-') {
			i++ // optional sign after e
		}
		seenExpDigit := false
		for i < n {
			if s[i] >= '0' && s[i] <= '9' {
				seenExpDigit = true
			} else {
				return false
			}
			i++
		}
		if !seenExpDigit {
			return false // exponent must have digits
		}
	}

	return i == n
}

func main() {
	cases := []struct {
		s        string
		expected bool
	}{
		{"2", true},
		{"0089", true},
		{"-0.1", true},
		{"+3.14", true},
		{"4.", true},
		{"-.9", true},
		{"2e10", true},
		{"-90E3", true},
		{"3e+7", true},
		{"+6e-1", true},
		{"53.5e93", true},
		{"-123.456e789", true},
		{"abc", false},
		{"1a", false},
		{"1e", false},
		{"e3", false},
		{"99e2.5", false},
		{"--6", false},
		{"-+3", false},
		{"95a54e53", false},
		{".", false},
		{".e1", false},
	}

	fmt.Println("=== Approach 1: DFA ===")
	for _, c := range cases {
		got := isNumberDFA(c.s)
		status := "✓"
		if got != c.expected {
			status = "✗ FAIL"
		}
		fmt.Printf("  %q  got=%v expected=%v %s\n", c.s, got, c.expected, status)
	}

	fmt.Println("=== Approach 2: Manual Parse ===")
	for _, c := range cases {
		got := isNumberParse(c.s)
		status := "✓"
		if got != c.expected {
			status = "✗ FAIL"
		}
		fmt.Printf("  %q  got=%v expected=%v %s\n", c.s, got, c.expected, status)
	}
}
