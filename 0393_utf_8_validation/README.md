# 0393 — UTF-8 Validation

> LeetCode #393 · Difficulty: Medium
> **Categories:** Bit Manipulation, Array

---

## Problem Statement

Given an integer array `data` representing the data, return whether it is a valid **UTF-8**
encoding (i.e. it translates to a sequence of valid UTF-8 encoded characters).

A character in UTF-8 can be from **1 to 4 bytes** long, subjected to the following rules:

1. For a **1-byte** character, the first bit is a `0`, followed by its Unicode code.
2. For an **n-bytes** character, the first `n` bits are all one's, the `n + 1` bit is `0`,
   followed by `n - 1` bytes with the most significant `2` bits being `10`.

This is how the UTF-8 encoding would work:

```
     Number of Bytes   |        UTF-8 Octet Sequence
                       |              (binary)
   --------------------+-----------------------------------------
            1          |   0xxxxxxx
            2          |   110xxxxx 10xxxxxx
            3          |   1110xxxx 10xxxxxx 10xxxxxx
            4          |   11110xxx 10xxxxxx 10xxxxxx 10xxxxxx
```

`x` denotes a bit in the binary form of a byte that may be either `0` or `1`.

**Note:** The input is an array of integers. Only the **least significant 8 bits** of each
integer is used to store the data. This means each integer represents only 1 byte of data.

**Example 1:**

```
Input: data = [197,130,1]
Output: true
Explanation: data represents the octet sequence: 11000101 10000010 00000001.
It is a valid utf-8 encoding for a 2-bytes character followed by a 1-byte character.
```

**Example 2:**

```
Input: data = [235,140,4]
Output: false
Explanation: data represented the octet sequence: 11101011 10001100 00000100.
The first 3 bits are all one's and the 4th bit is 0 means it is a 3-bytes character.
The next byte is a continuation byte which starts with 10 and that's correct.
But the second continuation byte does not start with 10, so it is invalid.
```

**Constraints:**

- `1 <= data.length <= 2 * 10^4`
- `0 <= data[i] <= 255`

---

## Company Frequency

| Company   | Frequency         | Last Reported |
|-----------|-------------------|---------------|
| Google    | ★★★☆☆ Medium      | 2023          |
| Facebook  | ★★★☆☆ Medium      | 2023          |
| Amazon    | ★★☆☆☆ Low         | 2022          |
| Cisco     | ★★☆☆☆ Low         | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Bit Manipulation** — inspect leading bits with shifts/masks to classify each byte → see [`/dsa/bit_manipulation.md`](/dsa/bit_manipulation.md)
- **Array single-pass state machine** — track "continuation bytes remaining" while scanning → see [`/dsa/arrays.md`](/dsa/arrays.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Count Continuation Bytes (Bit Masks) | O(n) | O(1) | Optimal; the interview answer |
| 2 | Bit-Pattern String Simulation | O(n) | O(1) | Most readable, mirrors the spec directly |

---

## Approach 1 — Count Continuation Bytes (Bit Masks)

### Intuition

Scan bytes left to right holding a counter `remaining` of continuation bytes still owed by
the current character. When `remaining == 0` we are at a character boundary and classify
the byte by its leading bits (`0`, `110`, `1110`, `11110`) to set how many continuation
bytes follow. Otherwise the byte must be a `10xxxxxx` continuation byte and we decrement.
Any bad leader, missing continuation, or leftover expectation is invalid.

### Algorithm

1. `remaining = 0`.
2. For each byte `b` (low 8 bits):
   - If `remaining == 0`: `0xxxxxxx`→0, `110`→1, `1110`→2, `11110`→3, else invalid.
   - Else: `b` must be `10xxxxxx`; decrement `remaining`.
3. Valid iff `remaining == 0` at the end.

### Complexity

- **Time:** O(n) — one pass, O(1) per byte.
- **Space:** O(1) — single counter.

### Code

```go
func validUtf8(data []int) bool {
	remaining := 0
	for _, num := range data {
		b := num & 0xFF
		if remaining == 0 {
			switch {
			case b>>7 == 0b0:
				remaining = 0
			case b>>5 == 0b110:
				remaining = 1
			case b>>4 == 0b1110:
				remaining = 2
			case b>>3 == 0b11110:
				remaining = 3
			default:
				return false
			}
		} else {
			if b>>6 != 0b10 {
				return false
			}
			remaining--
		}
	}
	return remaining == 0
}
```

### Dry Run

`data = [197, 130, 1]` → bytes `11000101`, `10000010`, `00000001`:

| byte | binary | remaining before | classification | remaining after |
|------|--------|-------------------|----------------|-----------------|
| 197 | 11000101 | 0 | `110xxxxx` → 2-byte leader | 1 |
| 130 | 10000010 | 1 | `10xxxxxx` continuation ✓ | 0 |
| 1 | 00000001 | 0 | `0xxxxxxx` → 1-byte char | 0 |

End with `remaining == 0` ⇒ **`true`**.

---

## Approach 2 — Bit-Pattern String Simulation

### Intuition

The count of leading `1` bits of a leader byte tells its length: `0`→1-byte, `k` (2..4)→
`k`-byte char needing `k-1` continuation bytes, and a lone leading `1` (`10xxxxxx`) can
never start a character. Rendering each byte as an 8-char binary string lets us count
leading ones and check continuation prefixes almost verbatim to the spec.

### Algorithm

1. Format each byte's low 8 bits as an 8-char binary string.
2. `ones` = leading `1` count. `0`→1-byte, `1` or `>4`→invalid.
3. For `ones` in 2..4: check the next `ones-1` bytes all start with `10`; bail if too few
   bytes remain.
4. Advance index past the whole character.

### Complexity

- **Time:** O(n) — each byte formatted/scanned a constant number of times.
- **Space:** O(1) — fixed 8-char scratch strings.

### Code

```go
func validUtf8Strings(data []int) bool {
	n := len(data)
	i := 0
	for i < n {
		bits := byteToBits(data[i])
		ones := leadingOnes(bits)
		switch {
		case ones == 0:
			i++
		case ones == 1 || ones > 4:
			return false
		default:
			need := ones - 1
			if i+need >= n {
				return false
			}
			for k := 1; k <= need; k++ {
				cont := byteToBits(data[i+k])
				if cont[0] != '1' || cont[1] != '0' {
					return false
				}
			}
			i += 1 + need
		}
	}
	return true
}
```

### Dry Run

`data = [197, 130, 1]`:

| i | byte | bits | leading ones | action | i after |
|---|------|------|--------------|--------|---------|
| 0 | 197 | 11000101 | 2 | 2-byte char; check data[1] starts "10" → "10000010" ✓ | 2 |
| 2 | 1 | 00000001 | 0 | 1-byte char | 3 |

Loop ends cleanly ⇒ **`true`**.

---

## Key Takeaways

- **Leader-byte leading-ones count = character length.** `0`→1, `110`→2, `1110`→3,
  `11110`→4; `10…` is only ever a continuation byte.
- Validate with a tiny **state machine**: one integer "continuation bytes remaining".
- `b >> k == pattern` is a clean way to test a fixed-length bit prefix; mask with `& 0xFF`
  first since only the low 8 bits count.
- Reject at three failure points: bad leader, non-`10` continuation, truncated tail
  (`remaining != 0` at end).

---

## Related Problems

- LeetCode #191 — Number of 1 Bits (bit counting)
- LeetCode #201 — Bitwise AND of Numbers Range (bit-prefix reasoning)
- LeetCode #338 — Counting Bits (per-value bit analysis)
- LeetCode #405 — Convert a Number to Hexadecimal (byte/nibble bit slicing)
