# 0476 — Number Complement

> LeetCode #476 · Difficulty: Easy
> **Categories:** Bit Manipulation

---

## Problem Statement

The **complement** of an integer is the integer you get when you flip all the `0`'s to `1`'s and all the `1`'s to `0`'s in its binary representation.

- For example, The integer `5` is `"101"` in binary and its **complement** is `"010"` which is the integer `2`.

Given an integer `num`, return *its complement*.

**Example 1:**

```
Input: num = 5
Output: 2
Explanation: The binary representation of 5 is 101 (no leading zero bits), and its complement is 010. So you need to output 2.
```

**Example 2:**

```
Input: num = 1
Output: 0
Explanation: The binary representation of 1 is 1 (no leading zero bits), and its complement is 0. So you need to output 0.
```

**Constraints:**

- `1 <= num < 2^31`

**Note:** This question is the same as [1009: Complement of Base 10 Integer](https://leetcode.com/problems/complement-of-base-10-integer/).

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Adobe      | ★★☆☆☆ Low        | 2022          |
| Cloudera   | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Bit Manipulation** — the entire task is bit flipping restricted to the bits at or below the most-significant set bit; the clean solution builds an all-ones mask of exactly the number's bit-width and XORs (`b ^ 1` flips a bit) → see [`/dsa/bit_manipulation.md`](/dsa/bit_manipulation.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Bit by Bit | O(log num) | O(1) | Most explicit; good to explain the idea from first principles |
| 2 | XOR with All-Ones Mask (Optimal) | O(1) | O(1) | The clean interview answer — one `bits.Len` + XOR |
| 3 | Smear Bits then XOR | O(1) | O(1) | Same idea with no length call; shows the "fill below MSB" idiom |

---

## Approach 1 — Bit by Bit

### Intuition

The complement is defined only over the meaningful bits — the ones from position 0 up to the highest set bit. There are **no leading zeros to flip**. So walk the bits of `num` from least to most significant while any remain: read each bit, invert it (`1 - bit`), and drop the inverted bit into the answer at the same position.

### Algorithm

1. Initialise `result = 0`, `position = 0`, and `remaining = num`.
2. While `remaining > 0`:
   - `bit = remaining & 1` — the lowest bit still to process.
   - `flipped = 1 - bit` — invert it.
   - `result |= flipped << position` — place it at the correct position.
   - `remaining >>= 1`, `position++`.
3. Return `result`.

### Complexity

- **Time:** O(log num) — one iteration per significant bit, at most ~31.
- **Space:** O(1) — three integer accumulators.

### Code

```go
func bitByBit(num int) int {
	result := 0      // the complement we are assembling
	position := 0    // which bit position we are currently at
	remaining := num // consume a copy so num stays intact
	for remaining > 0 {
		bit := remaining & 1          // read the lowest bit of the remaining number
		flipped := 1 - bit            // invert it: 0→1, 1→0
		result |= flipped << position // drop the flipped bit at its position
		remaining >>= 1               // advance to the next-higher bit
		position++                    // and remember the new position
	}
	return result
}
```

### Dry Run

Example 1: `num = 5` (binary `101`).

| Step | remaining (bin) | bit | flipped | flipped << position | result (bin) | position after |
|------|-----------------|-----|---------|---------------------|--------------|----------------|
| 1 | `101` (5) | 1 | 0 | `0 << 0` = 0 | `000` (0) | 1 |
| 2 | `10` (2)  | 0 | 1 | `1 << 1` = 2 | `010` (2) | 2 |
| 3 | `1` (1)   | 1 | 0 | `0 << 2` = 0 | `010` (2) | 3 |
| 4 | `0` (0)   | — | — | loop ends       | `010` (2) | — |

Result: `2` ✔

---

## Approach 2 — XOR with All-Ones Mask (Optimal)

### Intuition

Flipping one bit `b` is exactly `b XOR 1`. To flip *every* meaningful bit at once, XOR `num` with a mask that is all `1`s and **exactly as wide** as `num`. If `num` occupies `L` bits, that mask is `(1 << L) - 1`. Then `num ^ mask` inverts precisely those `L` bits and, because the mask is `0` above bit `L-1`, no spurious high bits appear.

For `num = 5 = 101`, `L = 3`, mask `= 111`, and `101 ^ 111 = 010 = 2`.

### Algorithm

1. `L = bits.Len(uint(num))` — the number of significant bits.
2. `mask = (1 << L) - 1` — `L` consecutive ones.
3. Return `num ^ mask`.

### Complexity

- **Time:** O(1) — `bits.Len` plus a shift, a subtract, and a XOR.
- **Space:** O(1).

### Code

```go
func xorMask(num int) int {
	length := bits.Len(uint(num)) // number of significant bits (0 for num==0)
	mask := (1 << length) - 1     // 'length' consecutive 1s: 5(101)→111
	return num ^ mask             // XOR flips exactly those bits
}
```

### Dry Run

Example 1: `num = 5` (`101`).

| Step | Expression | Value |
|------|------------|-------|
| 1 | `length = bits.Len(101₂)` | 3 |
| 2 | `mask = (1 << 3) - 1` | `111₂` = 7 |
| 3 | `num ^ mask = 101 ^ 111` | `010₂` = 2 |

Result: `2` ✔

---

## Approach 3 — Smear Bits then XOR

### Intuition

We need the same all-ones mask as Approach 2 but can build it with pure bit tricks. OR-ing `num` with itself shifted right by `1, 2, 4, 8, 16` "smears" the highest set bit downward, turning it and every lower bit into `1` — the classic *fill-all-bits-below-the-MSB* idiom. XOR `num` with that mask to flip the relevant region.

### Algorithm

1. `mask = num`.
2. `mask |= mask >> 1`, then `>> 2`, `>> 4`, `>> 8`, `>> 16` — after these five steps every bit from the MSB down is `1`.
3. Return `num ^ mask`.

### Complexity

- **Time:** O(1) — a fixed five shift/OR pairs covers all 32 bits.
- **Space:** O(1).

### Code

```go
func smearMask(num int) int {
	mask := num        // start from num
	mask |= mask >> 1  // fill 1 bit below each set bit
	mask |= mask >> 2  // fill 2 more
	mask |= mask >> 4  // 4 more
	mask |= mask >> 8  // 8 more
	mask |= mask >> 16 // 16 more → every bit from MSB down is now 1
	return num ^ mask  // flip exactly the smeared region
}
```

### Dry Run

Example 1: `num = 5` (`101`).

| Step | Operation | mask (bin) |
|------|-----------|------------|
| 0 | `mask = num` | `101` |
| 1 | `mask |= mask >> 1` (`101|010`) | `111` |
| 2 | `mask |= mask >> 2` (`111|001`) | `111` |
| 3 | `mask |= mask >> 4/8/16` | `111` (no change) |
| 4 | `num ^ mask = 101 ^ 111` | `010` = 2 |

Result: `2` ✔

---

## Key Takeaways

- **Complement = XOR with an all-ones mask of the number's own bit-width.** The subtlety is *width*: you must not flip leading zeros, so the mask has to stop at the highest set bit.
- Two ways to size the mask: compute the bit length (`(1 << L) - 1`) or smear the MSB downward with shift/OR — both give the same value.
- `b ^ 1` flips a single bit; XOR against a mask flips a whole selection of bits in one operation.
- Edge care: LeetCode #476 constrains `num >= 1`, but the mask logic also handles `num = 0` (length 0 → mask 0 → result 0), which is what the twin problem #1009 requires.

---

## Related Problems

- LeetCode #1009 — Complement of Base 10 Integer (identical problem)
- LeetCode #191 — Number of 1 Bits (bit scanning)
- LeetCode #190 — Reverse Bits (per-bit reconstruction)
- LeetCode #201 — Bitwise AND of Numbers Range (masks / common prefix)
- LeetCode #338 — Counting Bits (bit DP)
