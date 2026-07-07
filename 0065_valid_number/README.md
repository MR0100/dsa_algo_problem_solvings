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
