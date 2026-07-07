# 0273 — Integer to English Words

> LeetCode #273 · Difficulty: Hard
> **Categories:** Math, String, Recursion

---

## Problem Statement

Convert a non-negative integer `num` to its English words representation.

**Example 1:**

```
Input: num = 123
Output: "One Hundred Twenty Three"
```

**Example 2:**

```
Input: num = 12345
Output: "Twelve Thousand Three Hundred Forty Five"
```

**Example 3:**

```
Input: num = 1234567
Output: "One Million Two Hundred Thirty Four Thousand Five Hundred Sixty Seven"
```

**Constraints:**

- `0 <= num <= 2^31 - 1`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★★ Very High  | 2024          |
| Microsoft  | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Facebook   | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Math / Number Decomposition** — the number is peeled apart by powers of ten (thousands groups, hundreds, tens, ones) → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)
- **String Building** — the answer is assembled from lookup tables and scale words, with careful spacing/trimming → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)
- **Recursion** — Approach 1 recurses on each three-digit group and scale → see [`/dsa/divide_and_conquer.md`](/dsa/divide_and_conquer.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Recursive Divide by Scale | O(1) | O(1) | Cleanest to write; recursion mirrors the language structure |
| 2 | Iterative Grouping by Thousands | O(1) | O(1) | Loop-based, easy to reason about group-by-group |

> Both are O(1): `num < 2^31` has at most 10 digits (4 three-digit groups), a fixed bound.

---

## Approach 1 — Recursive Divide by Scale

### Intuition

English number words are structured in groups of three digits. Each group is read the same way ("one hundred twenty three") and then tagged with a scale word — Thousand, Million, Billion. That structure is naturally recursive: the words for `N` are `words(N / scale) + scaleWord + words(N % scale)` for the largest scale that fits, and below 1000 we spell hundreds, tens, and ones directly. Helpers return `""` for a zero group so parents concatenate cleanly, and the top level trims stray spaces and substitutes `"Zero"`.

### Algorithm

1. Top level: if `num == 0` return `"Zero"`; else return `TrimSpace(helper(num))`.
2. In `helper`: `num == 0` → `""`.
3. `num < 20` → `belowTwenty[num]`.
4. `num < 100` → `tens[num/10] + " " + helper(num%10)`.
5. `num < 1000` → `helper(num/100) + " Hundred " + helper(num%100)`.
6. Otherwise pick the largest of Thousand / Million / Billion that fits: `helper(num/scale) + " " + word + " " + helper(num%scale)`.

### Complexity

- **Time:** O(1) — recursion depth and work are bounded by the fixed number of scales.
- **Space:** O(1) — bounded recursion stack; output length is bounded by a constant.

### Code

```go
func recursive(num int) string {
	if num == 0 {
		return "Zero" // only reachable at the top level; helpers return "" for 0
	}
	return strings.TrimSpace(recursiveHelper(num))
}

func recursiveHelper(num int) string {
	switch {
	case num == 0:
		return "" // contributes nothing to its parent group
	case num < 20:
		return belowTwenty[num] // direct table hit for 1..19
	case num < 100:
		// tens digit word, then recurse on the ones digit (which may be 0 → "").
		return tens[num/10] + " " + recursiveHelper(num%10)
	case num < 1000:
		// hundreds digit, the literal "Hundred", then the remaining two digits.
		return recursiveHelper(num/100) + " Hundred " + recursiveHelper(num%100)
	case num < 1000000:
		return recursiveHelper(num/1000) + " Thousand " + recursiveHelper(num%1000)
	case num < 1000000000:
		return recursiveHelper(num/1000000) + " Million " + recursiveHelper(num%1000000)
	default:
		return recursiveHelper(num/1000000000) + " Billion " + recursiveHelper(num%1000000000)
	}
}
```

### Dry Run

Example 1: `num = 123`.

| Call | branch | expands to |
|------|--------|-----------|
| `helper(123)` | `< 1000` | `helper(1)` + `" Hundred "` + `helper(23)` |
| `helper(1)` | `< 20` | `"One"` |
| `helper(23)` | `< 100` | `tens[2]` + `" "` + `helper(3)` = `"Twenty "` + `"Three"` |
| combine | — | `"One" + " Hundred " + "Twenty Three"` = `"One Hundred Twenty Three"` |

Top level `TrimSpace` → `"One Hundred Twenty Three"` ✔

---

## Approach 2 — Iterative Grouping by Thousands

### Intuition

Rather than recurse, peel the number 1000 at a time. The lowest three digits get no scale word; the next three get "Thousand", then "Million", then "Billion". Convert each non-zero three-digit group with a helper that spells `[1, 999]`, prepend the scale word, and stitch the groups together most-significant-first.

### Algorithm

1. If `num == 0`, return `"Zero"`.
2. `scales = ["", "Thousand", "Million", "Billion"]`, `i = 0`, `parts = []`.
3. While `num > 0`: `g = num % 1000`; if `g != 0`, prepend `three(g) + " " + scales[i]` (trimmed) to `parts`; then `num /= 1000`, `i++`.
4. Join `parts` with spaces.

### Complexity

- **Time:** O(1) — at most 4 groups, each converted in O(1).
- **Space:** O(1) — at most 4 parts.

### Code

```go
func iterativeGroups(num int) string {
	if num == 0 {
		return "Zero"
	}
	scales := []string{"", "Thousand", "Million", "Billion"}
	var parts []string // most-significant group ends up first
	i := 0
	for num > 0 {
		if g := num % 1000; g != 0 {
			// Convert this three-digit group and tag with its scale word.
			chunk := strings.TrimSpace(three(g) + " " + scales[i])
			// Prepend so higher scales precede lower ones in the final order.
			parts = append([]string{chunk}, parts...)
		}
		num /= 1000 // move to the next higher three-digit group
		i++
	}
	return strings.Join(parts, " ")
}

func three(n int) string {
	var b []string
	if n >= 100 {
		// hundreds digit + literal "Hundred"
		b = append(b, belowTwenty[n/100], "Hundred")
		n %= 100
	}
	if n >= 20 {
		// tens word (e.g. "Forty"), drop to the ones digit
		b = append(b, tens[n/10])
		n %= 10
	}
	if n > 0 {
		// remaining 1..19 as a single word
		b = append(b, belowTwenty[n])
	}
	return strings.Join(b, " ")
}
```

### Dry Run

Example 2: `num = 12345`.

| iter | num (start) | g = num%1000 | three(g) | scale[i] | chunk prepended | parts | num after /1000 |
|------|-------------|--------------|----------|----------|-----------------|-------|-----------------|
| i=0 | 12345 | 345 | `Three Hundred Forty Five` | `""` | `Three Hundred Forty Five` | `[Three Hundred Forty Five]` | 12 |
| i=1 | 12 | 12 | `Twelve` | `Thousand` | `Twelve Thousand` | `[Twelve Thousand, Three Hundred Forty Five]` | 0 |

Join with spaces → `"Twelve Thousand Three Hundred Forty Five"` ✔

---

## Key Takeaways

- **Group by thousands.** Every scale word (Thousand/Million/Billion) governs a three-digit block; solve the block once (`three`) and reuse it.
- **Return `""` for empty groups** so concatenation stays clean, then `TrimSpace`/`Join` at the boundaries — this avoids double spaces and trailing spaces without special-casing every combination.
- **Handle `0` explicitly** at the top level: it is the one input that produces a non-empty word from an "empty" group.
- Lookup tables for `belowTwenty` and `tens` turn the messy irregular part of English numerals (eleven, twelve, thirteen, twenty, forty) into a constant-time array index.
- No number below `2^31` needs "Trillion", so four scale words suffice.

---

## Related Problems

- LeetCode #12 — Integer to Roman (digit-place decomposition into symbols)
- LeetCode #13 — Roman to Integer (inverse parse)
- LeetCode #65 — Valid Number (string ↔ number formatting rules)
- LeetCode #8 — String to Integer (atoi) (careful string/number conversion)
