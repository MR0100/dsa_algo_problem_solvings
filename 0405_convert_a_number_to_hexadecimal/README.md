# 0405 — Convert a Number to Hexadecimal

> LeetCode #405 · Difficulty: Easy
> **Categories:** Bit Manipulation, Math

---

## Problem Statement

Given a 32-bit integer `num`, return *a string representing its hexadecimal representation*. For negative integers, [two's complement](https://en.wikipedia.org/wiki/Two%27s_complement) method is used.

All the letters in the answer string should be lowercase characters, and there should not be any leading zeros in the answer except for the zero itself.

**Note:** You are not allowed to use any built-in library method to directly solve this problem.

**Example 1:**

```
Input: num = 26
Output: "1a"
```

**Example 2:**

```
Input: num = -1
Output: "ffffffff"
```

**Constraints:**

- `-2^31 <= num <= 2^31 - 1`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Facebook   | ★★☆☆☆ Low        | 2022          |
| Apple      | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Bit Manipulation** — hex maps 4 bits ↔ 1 digit; reinterpreting the signed int as `uint32` gives the two's-complement bit pattern for free, and `& 0xf` / `>> 4` peel off one nibble at a time → see [`/dsa/bit_manipulation.md`](/dsa/bit_manipulation.md)
- **Math / Number Theory (base conversion)** — the same result comes from repeated `% 16` and `/ 16`, the classic manual base-conversion algorithm → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Bit Masking (nibbles, low→high) | O(8) = O(1) | O(1) | The idiomatic answer; `& 0xf` + `>> 4` per hex digit |
| 2 | Fixed 8-Nibble Scan (high→low) | O(8) = O(1) | O(1) | Emits digits in order (no reversal); skip leading zeros inline |
| 3 | Repeated Division / Modulo | O(8) = O(1) | O(1) | Arithmetic view; shows `%16 ≡ &0xf`, `/16 ≡ >>4` |

---

## Approach 1 — Bit Masking

### Intuition

Hexadecimal is base 16, so each hex digit encodes exactly 4 bits (a "nibble"). Reinterpreting the signed `num` as an unsigned `uint32` *is* the two's-complement representation for negatives — for example `-1` becomes `0xFFFFFFFF`. Then repeatedly grab the lowest 4 bits with `& 0xf` to get a value in `0..15`, translate it to a hex character, and shift right by 4 to expose the next nibble. Because digits come out least-significant first, reverse them at the end.

### Algorithm

1. If `num == 0`, return `"0"`.
2. `n = uint32(num)` — negatives become their two's-complement bit pattern.
3. While `n != 0`: `nibble = n & 0xf`; record `hexDigits[nibble]`; `n >>= 4`.
4. Reverse the recorded digits and return them as a string.

### Complexity

- **Time:** O(8) = O(1) — a 32-bit number has at most 8 nibbles.
- **Space:** O(8) = O(1) — the small output buffer.

### Code

```go
func bitMasking(num int) string {
	if num == 0 {
		return "0" // the only case that legitimately prints a single zero
	}
	// Reinterpreting the signed int as uint32 yields the two's-complement bits,
	// e.g. -1 -> 0xFFFFFFFF, which is exactly what the problem wants.
	n := uint32(num)

	var sb strings.Builder
	// Collect nibbles from least-significant to most-significant.
	var digits []byte
	for n != 0 {
		nibble := n & 0xf                        // isolate the low 4 bits (0..15)
		digits = append(digits, hexDigits[nibble]) // its hex character
		n >>= 4                                   // drop those 4 bits
	}
	// digits are reversed (low first); write them back high-to-low.
	for i := len(digits) - 1; i >= 0; i-- {
		sb.WriteByte(digits[i])
	}
	return sb.String()
}
```

`hexDigits` is the shared lookup table `"0123456789abcdef"`.

### Dry Run

Example 1: `num = 26`. As a bit pattern `26 = 0001 1010₂`, so `n = 0x1A`.

| iteration | n (hex) | nibble = n & 0xf | hex char | n >>= 4 |
|-----------|---------|------------------|----------|---------|
| 1 | `0x1A` | `0xA` = 10 | `a` | `0x1` |
| 2 | `0x1`  | `0x1` = 1  | `1` | `0x0` |

Collected low→high: `['a','1']`. Reversed: `"1a"` ✔

---

## Approach 2 — Fixed 8-Nibble Scan

### Intuition

A 32-bit value is precisely 8 nibbles, so walk them from most-significant to least by shifting right 28, 24, …, 0 bits and masking `0xf`. Skip zero nibbles at the top (leading zeros) until the first significant nibble appears; from that point emit every nibble, including interior zeros. Because we go high→low, the digits are produced in final order with no reversal needed.

### Algorithm

1. If `num == 0`, return `"0"`.
2. `n = uint32(num)`; set `leading = true`.
3. For `shift` from `28` down to `0` in steps of `4`: `nibble = (n >> shift) & 0xf`.
   - If `nibble != 0`, set `leading = false`.
   - If `!leading`, append `hexDigits[nibble]`.
4. Return the assembled string (non-empty since `num != 0`).

### Complexity

- **Time:** O(8) = O(1) — exactly 8 iterations.
- **Space:** O(8) = O(1) — the output buffer.

### Code

```go
func topDownFixed(num int) string {
	if num == 0 {
		return "0"
	}
	n := uint32(num)
	var sb strings.Builder
	leading := true // still skipping high-order zero nibbles?
	for shift := 28; shift >= 0; shift -= 4 {
		nibble := (n >> uint(shift)) & 0xf // the nibble at this position
		if nibble != 0 {
			leading = false // first significant nibble reached
		}
		if !leading {
			sb.WriteByte(hexDigits[nibble]) // emit real digits and interior zeros
		}
	}
	return sb.String()
}
```

### Dry Run

Example 1: `num = 26`, `n = 0x0000001A`. Nibbles high→low:

| shift | nibble = (n >> shift) & 0xf | leading before | action | output so far |
|-------|-----------------------------|----------------|--------|---------------|
| 28 | 0 | true | skip | `` |
| 24 | 0 | true | skip | `` |
| 20 | 0 | true | skip | `` |
| 16 | 0 | true | skip | `` |
| 12 | 0 | true | skip | `` |
| 8  | 0 | true | skip | `` |
| 4  | 1 | true → false | emit `1` | `1` |
| 0  | A | false | emit `a` | `1a` |

Result: `"1a"` ✔

---

## Approach 3 — Repeated Division / Modulo

### Intuition

Base conversion done by hand: the last hex digit is `value % 16`, then divide the value by 16 and repeat. Performing this on `uint32(num)` handles negatives automatically through the two's-complement bit pattern. Digits emerge least-significant first, so reverse at the end. This is the arithmetic mirror of Approach 1 — `% 16` is exactly `& 0xf` and `/ 16` is exactly `>> 4` — and it uses no bitwise operators at all.

### Algorithm

1. If `num == 0`, return `"0"`.
2. `n = uint32(num)`. While `n > 0`: record `hexDigits[n % 16]`; `n /= 16`.
3. Reverse the recorded digits and return the string.

### Complexity

- **Time:** O(8) = O(1) — at most 8 divisions.
- **Space:** O(8) = O(1) — the output buffer.

### Code

```go
func divisionMod(num int) string {
	if num == 0 {
		return "0"
	}
	n := uint32(num) // two's-complement bit pattern viewed as an unsigned value

	var digits []byte
	for n > 0 {
		digits = append(digits, hexDigits[n%16]) // low-order hex digit
		n /= 16                                  // shift right one hex place
	}
	// Reverse: digits were produced least-significant first.
	var sb strings.Builder
	for i := len(digits) - 1; i >= 0; i-- {
		sb.WriteByte(digits[i])
	}
	return sb.String()
}
```

### Dry Run

Example 1: `num = 26`, `n = 26`.

| iteration | n | n % 16 | hex char | n /= 16 |
|-----------|---|--------|----------|---------|
| 1 | 26 | 10 | `a` | 1 |
| 2 | 1  | 1  | `1` | 0 |

Collected low→high: `['a','1']`. Reversed: `"1a"` ✔

---

## Key Takeaways

- **Signed → unsigned reinterpretation gives two's complement for free.** `uint32(num)` turns `-1` into `0xFFFFFFFF`, so negative handling needs no special arithmetic — the reason Example 2 prints `ffffffff`.
- **Hex ↔ 4 bits.** Any base that is a power of two lets you convert digit-by-digit with a mask (`& (base-1)`) and shift (`>> log2(base)`), avoiding general division.
- **`% 16 ≡ & 0xf` and `/ 16 ≡ >> 4`.** The division/modulo and bitwise formulations are literally the same computation; know both so the "no built-in conversion" constraint is never a blocker.
- **Order and the zero case are the two gotchas.** Low-first extraction needs a reversal (or a high→low scan), and `num == 0` must short-circuit to `"0"` so the loop doesn't emit an empty string.

---

## Related Problems

- LeetCode #190 — Reverse Bits (nibble/bit-level manipulation of a 32-bit word)
- LeetCode #191 — Number of 1 Bits (mask-and-shift over 32 bits)
- LeetCode #504 — Base 7 (repeated division base conversion)
- LeetCode #168 — Excel Sheet Column Title (base-26 conversion with a twist)
- LeetCode #401 — Binary Watch (bit-count interpretation of numbers)
