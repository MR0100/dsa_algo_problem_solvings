# 0065 — Valid Number

> LeetCode #65 · Difficulty: Hard
> **Categories:** String

---

## Problem Statement

Given a string `s`, return whether `s` is a **valid number**.

For example, all the following are valid numbers: `"2"`, `"0089"`, `"-0.1"`, `"+3.14"`, `"4."`, `"-.9"`, `"2e10"`, `"-90E3"`, `"3e+7"`, `"+6e-1"`, `"53.5e93"`, `"-123.456e789"`.

And the following are **not** valid numbers: `"abc"`, `"1a"`, `"1e"`, `"e3"`, `"99e2.5"`, `"--6"`, `"-+3"`, `"95a54e53"`.

Formally, a valid number is defined using one of the following definitions:
1. A **decimal number** can be defined as `[sign] integer ['.' [integer]]` or `[sign] '.' integer`, where `sign` is `+/-`, and `integer` is a sequence of digits.
2. An **integer** can be defined as `[sign] integer`.
3. A decimal or integer number can include an optional exponent `('e'|'E') integer-with-optional-sign`.

**Constraints**
- `1 <= s.length <= 20`
- `s` consists of only English letters, digits, `'+'`, `'-'`, or `'.'`.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★☆ High      | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Meta      | ★★★☆☆ Medium    | 2023          |
| Microsoft | ★★★☆☆ Medium    | 2023          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |
| Apple     | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Deterministic Finite Automaton (DFA)** — model valid number as a state machine; each character causes a state transition.
- **Manual Parsing** — parse grammar rules directly for readability.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | DFA ✅ | O(n) | O(1) | Handles all cases systematically; common in interviews |
| 2 | Manual Parsing | O(n) | O(1) | More readable; grammar-driven |

---

## Approach 1 — DFA (Recommended ✅)

### Intuition
Model the problem as a finite automaton where each character type (digit, sign, dot, e/E) causes a state transition. Reject if any transition leads to an invalid state or if the final state is non-accepting.

**States:**
```
0: initial
1: after leading sign
2: integer digits (before decimal point)
3: decimal point with no digits before it (e.g. ".5")
4: digits after decimal point with NO prior integer (e.g. ".5")
5: decimal point after integer digits (e.g. "1.")
6: after 'e'/'E'
7: after sign following 'e'
8: digits after 'e'
```

**Accept states:** 2, 4, 5, 8 (a valid number can end at any of these).

**Transition table (charType: 0=digit, 1=sign, 2=dot, 3=e/E):**
```
State 0: digit→2, sign→1, dot→3,  e→-1
State 1: digit→2, sign→-1, dot→3, e→-1
State 2: digit→2, sign→-1, dot→5, e→6
State 3: digit→4, sign→-1, dot→-1, e→-1
State 4: digit→4, sign→-1, dot→-1, e→6
State 5: digit→4, sign→-1, dot→-1, e→6
State 6: digit→8, sign→7,  dot→-1, e→-1
State 7: digit→8, sign→-1, dot→-1, e→-1
State 8: digit→8, sign→-1, dot→-1, e→-1
```

### Complexity
- **Time:** O(n).
- **Space:** O(1) — fixed-size transition table.

### Code
```go
// isNumberDFA solves Valid Number using a state machine that processes
// each character and transitions between states.
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
```

### Dry Run — `s = "-0.1"`
```
State 0: '-' (sign) → state 1
State 1: '0' (digit) → state 2
State 2: '.' (dot) → state 5
State 5: '1' (digit) → state 4

Final state: 4 (accept state) → true ✓
```

### Dry Run — `s = "1e"`
```
State 0: '1' → state 2
State 2: 'e' → state 6
Final state: 6 (NOT accept) → false ✓
```

---

## Approach 2 — Manual Parsing

### Intuition
Parse the string according to the grammar directly:
1. Skip optional leading sign.
2. Scan for digits/dot to form the mantissa. Track whether we saw a digit and whether we saw a dot.
3. If we see `e`/`E`, skip it, skip optional sign, then require at least one digit.

### Complexity
- **Time:** O(n).
- **Space:** O(1).

### Code
```go
// isNumberParse solves Valid Number by parsing the string manually according to
// the grammar:
//   valid number = (integer | decimal) optional-exponent
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
```

### Dry Run — `s = "-0.1"`

`n = 4`. `s[0] = '-'` → skip sign, `i = 1`. `seenDigit = false`, `seenDot = false`:

| `i` | `s[i]` | Branch | State after |
|-----|--------|--------|-------------|
| 1 | `0` | digit | `seenDigit = true`, `i = 2` |
| 2 | `.` | dot & `!seenDot` | `seenDot = true`, `i = 3` |
| 3 | `1` | digit | `seenDigit = true`, `i = 4` |

Mantissa loop ends (`i == n`). `seenDigit = true`. No `e`/`E`, so exponent skipped. `return i == n` → `4 == 4` → **true**.

---

## Key Takeaways

- **`"4."` is valid** — a dot with digits before it and no digits after is OK.
- **`"."` is NOT valid** — a dot with no digits on either side is invalid.
- **`".e1"` is NOT valid** — needs at least one digit before the `e`.
- **`"1e"` is NOT valid** — exponent part requires at least one digit.
- **DFA is the systematic approach** — in an interview, drawing the state diagram first and then translating it to code shows deep understanding. For ad-hoc parsing, flag-based parsing is more readable.
- **Accept states encode the grammar** — states 2, 4, 5, 8 correspond to: integer, "." digit, "integer.", "integer e integer" — all valid endpoints.

---

## Related Problems

- LeetCode #8 — String to Integer (atoi) (simpler parsing; only integers)
- LeetCode #150 — Evaluate Reverse Polish Notation (parsing-like expression eval)
