# 0171 — Excel Sheet Column Number

> LeetCode #171 · Difficulty: Easy
> **Categories:** Math, String

---

## Problem Statement

Given a string `columnTitle` that represents the column title as appears in an Excel sheet, return *its corresponding column number*.

For example:

```
A -> 1
B -> 2
C -> 3
...
Z -> 26
AA -> 27
AB -> 28
...
```

**Example 1:**

```
Input: columnTitle = "A"
Output: 1
```

**Example 2:**

```
Input: columnTitle = "AB"
Output: 28
```

**Example 3:**

```
Input: columnTitle = "ZY"
Output: 701
```

**Constraints:**

- `1 <= columnTitle.length <= 7`
- `columnTitle` consists only of uppercase English letters.
- `columnTitle` is in the range `["A", "FXSHRXW"]`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Microsoft  | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2023          |
| Adobe      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Math / Number Theory (base conversion)** — the title is a number written in *bijective* base-26: digits run A=1..Z=26 with no zero, so `value = Σ digit·26^position` → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)
- **String Algorithms** — character-by-character parsing of a string into a number (the base-26 twin of `atoi`'s Horner loop) → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (recompute power per character) | O(k²), k = len ≤ 7 | O(1) | To derive the place-value formula from first principles |
| 2 | Right-to-Left with Running Power | O(k) | O(1) | When you naturally think "ones place, 26s place, 676s place…" |
| 3 | Left-to-Right Horner's Method (Optimal) | O(k) | O(1) | Always — the canonical `result = result*base + digit` parse loop |

---

## Approach 1 — Brute Force (Recompute Power per Character)

### Intuition

`"AB"` means A·26¹ + B·26⁰ with digit values A=1 … Z=26 — the same way the decimal string `"28"` means 2·10¹ + 8·10⁰. The most literal translation of that formula: for each letter, count how many positions sit to its right, compute 26 raised to that count with a fresh inner loop, and add `digit × power` to a running total. No cleverness, just the definition of positional notation.

### Algorithm

1. Set `total = 0`, `k = len(columnTitle)`.
2. For every index `i` from `0` to `k-1`:
   1. `digit = columnTitle[i] - 'A' + 1` (map `'A'..'Z'` → `1..26`).
   2. Compute `power = 26^(k-1-i)` with an inner loop of `k-1-i` multiplications.
   3. `total += digit * power`.
3. Return `total`.

### Complexity

- **Time:** O(k²) — each of the k characters recomputes its power in up to k−1 multiplications; with k ≤ 7 this is at most 21 multiplies, effectively constant.
- **Space:** O(1) — only scalar accumulators.

### Code

```go
func bruteForce(columnTitle string) int {
	total := 0
	k := len(columnTitle)
	for i := 0; i < k; i++ {
		digit := int(columnTitle[i]-'A') + 1 // 'A'→1 ... 'Z'→26 (no zero digit!)
		// Recompute 26^(k-1-i) naively — the "brute" part of this approach.
		power := 1
		for p := 0; p < k-1-i; p++ {
			power *= 26 // one factor of 26 per position right of index i
		}
		total += digit * power // place value contribution of this letter
	}
	return total
}
```

### Dry Run

Example 1: `columnTitle = "A"` (k = 1).

| Step | i | char | digit | inner loop runs | power | contribution | total |
|------|---|------|-------|-----------------|-------|--------------|-------|
| 1 | init | — | — | — | — | — | 0 |
| 2 | 0 | `'A'` | 1 | k−1−i = 0 times | 1 | 1 × 1 = 1 | 1 |

Result: `1` ✔

Bonus trace of Example 3 (`columnTitle = "ZY"`, k = 2):

| Step | i | char | digit | inner loop runs | power | contribution | total |
|------|---|------|-------|-----------------|-------|--------------|-------|
| 1 | 0 | `'Z'` | 26 | 1 time (26¹) | 26 | 26 × 26 = 676 | 676 |
| 2 | 1 | `'Y'` | 25 | 0 times (26⁰) | 1 | 25 × 1 = 25 | 701 |

Result: `701` ✔

---

## Approach 2 — Right-to-Left with Running Power

### Intuition

Approach 1 wastes work: 26^(i+1) is just 26^i × 26, so there is no reason to rebuild each power from scratch. Walk the string from the rightmost (least-significant) letter, keep the current place value in one variable, and grow it by ×26 every time we step one position left. Same formula, single pass, no inner loop.

### Algorithm

1. Set `total = 0`, `power = 1`.
2. For `i` from `len-1` down to `0`:
   1. `digit = columnTitle[i] - 'A' + 1`.
   2. `total += digit * power`.
   3. `power *= 26` (the next position left is 26× more significant).
3. Return `total`.

### Complexity

- **Time:** O(k) — exactly one multiply-add per character.
- **Space:** O(1) — two scalar accumulators (`total`, `power`).

### Code

```go
func rightToLeftPower(columnTitle string) int {
	total := 0
	power := 1 // place value of the current position: 26^0, 26^1, ...
	for i := len(columnTitle) - 1; i >= 0; i-- {
		digit := int(columnTitle[i]-'A') + 1 // bijective digit 1..26
		total += digit * power               // add this letter's contribution
		power *= 26                          // next position left is 26× more significant
	}
	return total
}
```

### Dry Run

Example 1: `columnTitle = "A"`.

| Step | i | char | digit | power (before) | total after `+= digit*power` | power after `*= 26` |
|------|---|------|-------|----------------|------------------------------|---------------------|
| 1 | init | — | — | 1 | 0 | — |
| 2 | 0 | `'A'` | 1 | 1 | 0 + 1·1 = 1 | 26 |
| 3 | loop ends (i = −1) | — | — | — | 1 | — |

Result: `1` ✔

Bonus trace of Example 2 (`columnTitle = "AB"`):

| Step | i | char | digit | power (before) | total | power after |
|------|---|------|-------|----------------|-------|-------------|
| 1 | 1 | `'B'` | 2 | 1 | 0 + 2·1 = 2 | 26 |
| 2 | 0 | `'A'` | 1 | 26 | 2 + 1·26 = 28 | 676 |

Result: `28` ✔

---

## Approach 3 — Left-to-Right Horner's Method (Optimal)

### Intuition

Horner's rule: A·26¹ + B·26⁰ = (A)·26 + B. Reading left to right, every time one more letter appears on the right, *everything read so far* becomes 26× more significant. So keep a single accumulator: multiply it by 26, add the new digit. No powers tracked at all — this is exactly how `atoi` parses `"28"` (`result = result*10 + digit`), just in base 26 with digits starting at 1 instead of 0.

### Algorithm

1. Set `result = 0`.
2. For each character `c` from left to right: `result = result*26 + (c - 'A' + 1)`.
3. Return `result`.

### Complexity

- **Time:** O(k) — one multiply-add per character, single forward pass.
- **Space:** O(1) — a single accumulator.

### Code

```go
func hornersMethod(columnTitle string) int {
	result := 0
	for i := 0; i < len(columnTitle); i++ {
		// Shift everything seen so far one position left (×26),
		// then drop the new least-significant digit (1..26) in.
		result = result*26 + int(columnTitle[i]-'A') + 1
	}
	return result
}
```

### Dry Run

Example 1: `columnTitle = "A"`.

| Step | i | char | digit | result = result·26 + digit |
|------|---|------|-------|----------------------------|
| 1 | init | — | — | 0 |
| 2 | 0 | `'A'` | 1 | 0·26 + 1 = **1** |

Result: `1` ✔

Bonus trace of Example 3 (`columnTitle = "ZY"`):

| Step | i | char | digit | result = result·26 + digit |
|------|---|------|-------|----------------------------|
| 1 | 0 | `'Z'` | 26 | 0·26 + 26 = 26 |
| 2 | 1 | `'Y'` | 25 | 26·26 + 25 = **701** |

Result: `701` ✔

---

## Key Takeaways

- **String → number is always Horner's loop:** `result = result*base + digit`. Decimal parsing, binary parsing, and Excel titles are the same three lines with a different base and digit mapping.
- **Bijective base-26:** digits run 1..26 with *no zero* — encoded here by the `+1` in `c - 'A' + 1`. Title→number is the easy direction; the inverse (number→title, LeetCode #168) is where the missing zero bites and forces a `-1` before every mod/div.
- The constraint upper bound `"FXSHRXW"` = 2³¹ − 1 is a hint: the answer fits an `int32`, and no overflow handling is needed in Go's `int`.
- Doing #171 and #168 back-to-back cements the bijective-numeration mapping in both directions.

---

## Related Problems

- LeetCode #168 — Excel Sheet Column Title (the exact inverse mapping, number → title)
- LeetCode #8 — String to Integer (atoi) (same Horner parse loop in base 10)
- LeetCode #13 — Roman to Integer (symbol-by-symbol string → number accumulation)
- LeetCode #7 — Reverse Integer (digit-by-digit mod/div processing)
