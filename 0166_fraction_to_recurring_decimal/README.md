# 0166 — Fraction to Recurring Decimal

> LeetCode #166 · Difficulty: Medium
> **Categories:** Hash Table, Math, String

---

## Problem Statement

Given two integers representing the `numerator` and `denominator` of a fraction, return *the fraction in string format*.

If the fractional part is repeating, enclose the repeating part in parentheses.

If multiple answers are possible, return **any of them**.

It is **guaranteed** that the length of the answer string is less than `10^4` for all the given inputs.

**Example 1:**

```
Input: numerator = 1, denominator = 2
Output: "0.5"
```

**Example 2:**

```
Input: numerator = 2, denominator = 1
Output: "2"
```

**Example 3:**

```
Input: numerator = 4, denominator = 333
Output: "0.(012)"
```

**Constraints:**

- `-2^31 <= numerator, denominator <= 2^31 - 1`
- `denominator != 0`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Meta       | ★★★★★ Very High  | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |
| Goldman Sachs | ★★☆☆☆ Low     | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Map** — map each long-division remainder to the string position of the digit it produced, so a repeated remainder (the start of the cycle) is detected in O(1) → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Math / Number Theory** — the whole solution is schoolbook long division; a decimal repeats iff a remainder repeats, and remainders live in `[1, den-1]`, so a cycle is guaranteed within `den` steps → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)
- **String Algorithms** — careful incremental string building (sign, integer part, decimal point, parentheses splice) with `strings.Builder` → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (linear remainder scan) | O(L²) | O(L) | To understand *why* remainder repetition = digit cycle; fine for tiny outputs |
| 2 | Hash Map of Remainders (Optimal) | O(L) | O(L) | Always — the canonical interview answer |

*(L = length of the output string, guaranteed < 10^4.)*

---

## Approach 1 — Brute Force (Long Division + Linear Remainder Scan)

### Intuition

Do the division exactly the way you learned in school: divide, take the remainder, multiply it by 10, divide again, and so on. The key mathematical fact is that **the digit stream is fully determined by the current remainder** — if the same remainder ever shows up twice, the digits between the two occurrences will repeat forever. The brute-force way to detect that repeat is to store every remainder in a slice and linearly scan it before producing each new digit.

### Algorithm

1. If `numerator == 0`, return `"0"` immediately.
2. Write a leading `-` iff exactly one of the two operands is negative (XOR of signs).
3. Convert both operands to `int64` and take absolute values (`|-2^31|` overflows 32-bit).
4. Append the integer part `num / den`. If `num % den == 0`, the division is exact — return.
5. Append `"."` and start long division with `rem = num % den`:
   - Before emitting a digit, scan the slice of past remainders. If `rem` was seen at index `i`, then fractional digits `digits[i:]` form the cycle → emit `digits[:i] + "(" + digits[i:] + ")"` and return.
   - Otherwise append `rem` to the slice, emit digit `rem*10 / den`, and set `rem = rem*10 % den`.
6. If `rem` reaches `0`, the decimal terminates — emit all digits with no parentheses.

### Complexity

- **Time:** O(L²) — each of the up-to-`L` fractional digits does a linear scan over up to `L` stored remainders.
- **Space:** O(L) — the remainder slice and the digit buffer each hold at most one entry per output digit.

### Code

```go
func bruteForce(numerator int, denominator int) string {
	// Zero numerator short-circuits everything: 0 / anything = "0".
	if numerator == 0 {
		return "0"
	}

	var sb strings.Builder
	// The result is negative iff exactly one operand is negative.
	if (numerator < 0) != (denominator < 0) {
		sb.WriteByte('-')
	}

	// Promote to int64 before negating: -(-2^31) overflows a 32-bit int.
	num, den := int64(numerator), int64(denominator)
	if num < 0 {
		num = -num
	}
	if den < 0 {
		den = -den
	}

	// Integer part of the quotient.
	sb.WriteString(strconv.FormatInt(num/den, 10))

	rem := num % den
	// Division is exact — no fractional part needed.
	if rem == 0 {
		return sb.String()
	}
	sb.WriteByte('.')

	// remainders[i] is the remainder that produced fractional digit i.
	remainders := []int64{}
	// digits accumulates the fractional digits produced so far.
	digits := []byte{}

	for rem != 0 {
		// Linear scan: has this remainder produced a digit before?
		for i, r := range remainders {
			if r == rem {
				// Cycle found: digits[i:] repeat forever → wrap in parens.
				sb.Write(digits[:i])
				sb.WriteByte('(')
				sb.Write(digits[i:])
				sb.WriteByte(')')
				return sb.String()
			}
		}
		// Record the remainder responsible for the digit we emit next.
		remainders = append(remainders, rem)
		rem *= 10                                  // bring down a zero, as in long division
		digits = append(digits, byte('0'+rem/den)) // next fractional digit
		rem %= den                                 // remainder carried into the next step
	}

	// Remainder hit 0 → terminating decimal, no parentheses.
	sb.Write(digits)
	return sb.String()
}
```

### Dry Run

Example 1: `numerator = 1, denominator = 2`.

| Step | Action | num | den | rem | remainders | digits | Output so far |
|------|--------|-----|-----|-----|------------|--------|---------------|
| 1 | numerator ≠ 0, signs equal → no `-` | 1 | 2 | — | [] | [] | `""` |
| 2 | integer part `1/2 = 0` appended | 1 | 2 | — | [] | [] | `"0"` |
| 3 | `rem = 1 % 2 = 1` ≠ 0 → append `"."` | 1 | 2 | 1 | [] | [] | `"0."` |
| 4 | scan `[]` — no repeat; store rem 1 | 1 | 2 | 1 | [1] | [] | `"0."` |
| 5 | `rem = 1*10 = 10`; digit `10/2 = 5`; `rem = 10%2 = 0` | 1 | 2 | 0 | [1] | ['5'] | `"0."` |
| 6 | loop exits (`rem == 0`) → write digits | 1 | 2 | 0 | [1] | ['5'] | `"0.5"` |

Result: `"0.5"` ✔ (terminating decimal, no parentheses).

---

## Approach 2 — Hash Map of Remainders (Optimal)

### Intuition

Approach 1's only inefficiency is the linear scan asking "have I seen this remainder before?" A hash map answers that in O(1). Even better, if the map stores the **string index where each remainder's digit begins**, then on the first repeat we know exactly where to splice the `(` — everything from that index onward is the repeating block.

### Algorithm

1. Handle `numerator == 0` → `"0"`.
2. Write `-` iff the operand signs differ; convert to `int64` absolute values.
3. Append the integer part; if `num % den == 0`, return (exact division).
4. Append `"."`. Create `seen := map[int64]int{}` (remainder → index in the built string).
5. While `rem != 0`:
   - If `seen[rem]` exists at position `pos`: return `s[:pos] + "(" + s[pos:] + ")"`.
   - Else record `seen[rem] = sb.Len()`, then `rem *= 10`, append digit `rem/den`, set `rem %= den`.
6. If the loop exits (`rem == 0`), the decimal terminates — return the built string.

### Complexity

- **Time:** O(L) — each output digit costs one O(1) map lookup + insert; L < 10^4 by the problem guarantee.
- **Space:** O(L) — at most one map entry per fractional digit.

### Code

```go
func hashMap(numerator int, denominator int) string {
	// Zero numerator short-circuits everything: 0 / anything = "0".
	if numerator == 0 {
		return "0"
	}

	var sb strings.Builder
	// Negative result iff signs differ.
	if (numerator < 0) != (denominator < 0) {
		sb.WriteByte('-')
	}

	// int64 avoids overflow when taking absolute values of -2^31.
	num, den := int64(numerator), int64(denominator)
	if num < 0 {
		num = -num
	}
	if den < 0 {
		den = -den
	}

	// Integer part.
	sb.WriteString(strconv.FormatInt(num/den, 10))

	rem := num % den
	if rem == 0 {
		return sb.String() // terminating with no fractional part
	}
	sb.WriteByte('.')

	// seen maps a remainder to the string index where its digit begins.
	seen := map[int64]int{}

	for rem != 0 {
		if pos, ok := seen[rem]; ok {
			// Repeat detected: everything from pos onward is the cycle.
			s := sb.String()
			return s[:pos] + "(" + s[pos:] + ")"
		}
		// The digit generated by this remainder starts at the current length.
		seen[rem] = sb.Len()
		rem *= 10                                      // long-division step
		sb.WriteString(strconv.FormatInt(rem/den, 10)) // emit next digit
		rem %= den                                     // carry the new remainder
	}

	// rem became 0 → the decimal terminates.
	return sb.String()
}
```

### Dry Run

Example 1: `numerator = 1, denominator = 2`.

| Step | Action | rem | seen | Built string |
|------|--------|-----|------|--------------|
| 1 | numerator ≠ 0; signs equal → no `-` | — | {} | `""` |
| 2 | integer part `1/2 = 0` | — | {} | `"0"` |
| 3 | `rem = 1 % 2 = 1` ≠ 0 → append `"."` | 1 | {} | `"0."` |
| 4 | `1 ∉ seen` → record `seen[1] = 2` (current length) | 1 | {1: 2} | `"0."` |
| 5 | `rem = 10`; append `10/2 = 5`; `rem = 10 % 2 = 0` | 0 | {1: 2} | `"0.5"` |
| 6 | `rem == 0` → loop exits, return | 0 | {1: 2} | `"0.5"` |

Result: `"0.5"` ✔

Bonus trace of Example 3 (`4/333`) to show the cycle splice:

| Step | rem before | seen lookup | seen after | digit emitted | Built string |
|------|-----------|-------------|------------|---------------|--------------|
| 1 | 4 | miss | {4: 2} | `40/333 = 0` | `"0.0"` |
| 2 | 40 | miss | {4:2, 40:3} | `400/333 = 1` | `"0.01"` |
| 3 | 67 | miss | {4:2, 40:3, 67:4} | `670/333 = 2` | `"0.012"` |
| 4 | 4 | **hit at pos 2** | — | — | `"0." + "(" + "012" + ")"` |

Result: `"0.(012)"` ✔

---

## Key Takeaways

- **Repeating decimal ⇔ repeating remainder.** The digit stream of long division is a pure function of the current remainder, so cycle detection reduces to "first repeated remainder" — a classic hash-map fingerprint pattern (same trick as LeetCode #202 Happy Number).
- **Store positions, not just presence.** Mapping remainder → string index makes inserting the `(` trivial; a plain set would tell you *that* a cycle exists but not *where* it starts.
- **Sign and overflow first.** `-2^31 / -1` and `abs(-2^31)` are the classic traps — promote to `int64` before negating, and derive the sign with XOR-of-signs before taking absolute values (`-1/2` must print `-0.5`, which naive integer division would render as `0.5`).
- The remainder is bounded by `den - 1`, so the fractional part must terminate or cycle within `den` steps — this is why the answer length guarantee exists.

---

## Related Problems

- LeetCode #202 — Happy Number (cycle detection via repeated state)
- LeetCode #29 — Divide Two Integers (manual division with sign/overflow edge cases)
- LeetCode #43 — Multiply Strings (digit-by-digit arithmetic as string building)
- LeetCode #972 — Equal Rational Numbers (reasoning about repeating decimal representations)
