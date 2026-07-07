# 0461 — Hamming Distance

> LeetCode #461 · Difficulty: Easy
> **Categories:** Bit Manipulation

---

## Problem Statement

The [Hamming distance](https://en.wikipedia.org/wiki/Hamming_distance) between two integers is the number of positions at which the corresponding bits are different.

Given two integers `x` and `y`, return *the Hamming distance between them*.

**Example 1:**

```
Input: x = 1, y = 4
Output: 2
Explanation:
1   (0 0 0 1)
4   (0 1 0 0)
       ↑   ↑
The above arrows point to positions where the corresponding bits are different.
```

**Example 2:**

```
Input: x = 3, y = 1
Output: 1
```

**Constraints:**

- `0 <= x, y <= 2^31 - 1`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Facebook   | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Google     | ★★☆☆☆ Low        | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Adobe      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Bit Manipulation** — the entire problem is a property of the individual bits. `x ^ y` produces a 1 in exactly the positions where the two numbers differ, so the answer is `popcount(x ^ y)`; Kernighan's `n & (n-1)` clears the lowest set bit to count them in the fewest steps → see [`/dsa/bit_manipulation.md`](/dsa/bit_manipulation.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (compare each bit) | O(32) = O(1) | O(1) | Baseline; makes the definition explicit |
| 2 | XOR + count set bits (shift loop) | O(log max) ≤ 31 | O(1) | Clean, no library, once you know the XOR trick |
| 3 | XOR + Brian Kernighan's trick (Optimal) | O(popcount) ≤ 32 | O(1) | Fewest iterations; classic bit idiom |
| 4 | XOR + built-in popcount | O(1) | O(1) | Shortest real-world code; hardware POPCNT |

---

## Approach 1 — Brute Force (Compare Bit by Bit)

### Intuition

The Hamming distance is defined as "the number of positions where the bits differ", so the most literal solution walks all 32 bit positions, extracts bit `i` from both numbers, and tallies a position whenever the two bits disagree. Because both inputs fit in 31 bits, checking positions `0..31` covers every possible value.

### Algorithm

1. Initialise `distance = 0`.
2. For each position `i` from `0` to `31`:
   - `bitX = (x >> i) & 1` and `bitY = (y >> i) & 1`.
   - If `bitX != bitY`, increment `distance`.
3. Return `distance`.

### Complexity

- **Time:** O(32) = O(1) — a fixed 32 iterations regardless of the input magnitude.
- **Space:** O(1) — a single counter.

### Code

```go
func bruteForce(x int, y int) int {
	distance := 0 // running count of differing bit positions
	// 32 positions cover every value in [0, 2^31 - 1].
	for i := 0; i < 32; i++ {
		bitX := (x >> i) & 1 // isolate bit i of x (0 or 1)
		bitY := (y >> i) & 1 // isolate bit i of y (0 or 1)
		if bitX != bitY {    // positions disagree → contributes to the distance
			distance++
		}
	}
	return distance
}
```

### Dry Run

Example 1: `x = 1 (…0001), y = 4 (…0100)`. Only the low 4 positions matter; positions `4..31` are `0` in both, so they never differ.

| Position i | bit i of x | bit i of y | differ? | distance after |
|------------|------------|------------|---------|----------------|
| 0 | 1 | 0 | yes | 1 |
| 1 | 0 | 0 | no | 1 |
| 2 | 0 | 1 | yes | 2 |
| 3 | 0 | 0 | no | 2 |
| 4..31 | 0 | 0 | no | 2 |

Result: `2` ✔

---

## Approach 2 — XOR then Count Set Bits (Loop)

### Intuition

XOR flips exactly the bits that differ: `x ^ y` has a `1` in every position where `x` and `y` disagree, and `0` everywhere they agree. So the Hamming distance equals the number of set bits (the population count) of `x ^ y`. Counting those bits by testing the lowest bit and shifting right is the plain way to do it.

### Algorithm

1. Compute `xor = x ^ y`.
2. While `xor != 0`:
   - Add `xor & 1` to `distance` (adds 1 if the lowest bit is set).
   - Shift `xor >>= 1`.
3. Return `distance`.

### Complexity

- **Time:** O(log max(x, y)) — the loop runs until the highest set bit is shifted out, at most 31 times.
- **Space:** O(1) — one accumulator.

### Code

```go
func xorCountLoop(x int, y int) int {
	xor := x ^ y  // 1s mark exactly the positions where x and y differ
	distance := 0 // number of set bits so far
	// Peel off the lowest bit each iteration until nothing is left.
	for xor != 0 {
		distance += xor & 1 // add 1 if the lowest bit is set
		xor >>= 1           // drop the lowest bit and continue
	}
	return distance
}
```

### Dry Run

Example 1: `x = 1, y = 4`, so `xor = 1 ^ 4 = 5 = 101₂`.

| Step | xor (bin) | xor & 1 | distance after | xor after >>1 |
|------|-----------|---------|----------------|---------------|
| 1 | `101` (5) | 1 | 1 | `10` (2) |
| 2 | `10` (2)  | 0 | 1 | `1` (1) |
| 3 | `1` (1)   | 1 | 2 | `0` (0) |

Loop ends (`xor == 0`). Result: `2` ✔

---

## Approach 3 — XOR then Brian Kernighan's Trick (Optimal)

### Intuition

Same XOR insight, faster counting. `n & (n - 1)` clears the **lowest set bit** of `n` in a single operation. If we keep applying it to `xor` until it becomes `0`, the number of iterations equals the number of `1` bits — i.e. the Hamming distance. Unlike the shift loop, this skips straight from one set bit to the next, so it never runs more than `popcount(xor)` times.

### Algorithm

1. Compute `xor = x ^ y`.
2. While `xor != 0`:
   - `xor &= xor - 1` (erase the lowest set bit).
   - Increment `distance`.
3. Return `distance`.

### Complexity

- **Time:** O(popcount(x ^ y)) ≤ O(32) — exactly one iteration per differing bit.
- **Space:** O(1).

### Code

```go
func xorKernighan(x int, y int) int {
	xor := x ^ y  // differing positions become 1
	distance := 0 // count of set bits removed
	// Each pass deletes the lowest remaining 1 bit.
	for xor != 0 {
		xor &= xor - 1 // Kernighan: clears exactly the lowest set bit
		distance++     // we just accounted for one differing position
	}
	return distance
}
```

### Dry Run

Example 1: `x = 1, y = 4`, so `xor = 5 = 101₂` (two set bits → two iterations).

| Step | xor before | xor - 1 | xor & (xor-1) | distance after |
|------|------------|---------|---------------|----------------|
| 1 | `101` (5) | `100` (4) | `100` (4) | 1 |
| 2 | `100` (4) | `011` (3) | `000` (0) | 2 |

Loop ends (`xor == 0`). Result: `2` ✔ — two bit-clears, so two differing positions.

---

## Approach 4 — Built-in Population Count

### Intuition

Once you accept the answer is `popcount(x ^ y)`, there is nothing left to hand-roll: Go's `math/bits.OnesCount` maps to a single `POPCNT`-class CPU instruction on modern hardware. This is the shortest and fastest production version.

### Algorithm

1. Return `bits.OnesCount(uint(x ^ y))`.

### Complexity

- **Time:** O(1) — a single hardware instruction where available.
- **Space:** O(1).

### Code

```go
func builtinPopcount(x int, y int) int {
	// x^y marks differing bits; OnesCount reports how many there are.
	return bits.OnesCount(uint(x ^ y))
}
```

### Dry Run

Example 1: `x = 1, y = 4`.

| Step | Expression | Value |
|------|------------|-------|
| 1 | `x ^ y` | `5` = `101₂` |
| 2 | `bits.OnesCount(5)` | `2` |

Result: `2` ✔

---

## Key Takeaways

- **`x ^ y` isolates differences.** Any "how many positions differ" question over bit patterns reduces to a population count of the XOR.
- **`n & (n - 1)` clears the lowest set bit** — the go-to idiom for counting set bits in `popcount` iterations rather than one-per-bit-position (see #191 Number of 1 Bits, #201 Bitwise AND of Numbers Range).
- **Reach for `math/bits`** in real code: `OnesCount`, `LeadingZeros`, `TrailingZeros` compile to single instructions and are clearer than hand loops.
- The extension #477 (Total Hamming Distance) sums this over all pairs — solved by counting, per bit position, how many numbers have that bit set, instead of `O(n²)` pairwise XORs.

---

## Related Problems

- LeetCode #477 — Total Hamming Distance (sum of pairwise Hamming distances)
- LeetCode #191 — Number of 1 Bits (popcount of a single number)
- LeetCode #201 — Bitwise AND of Numbers Range (Kernighan's lowest-set-bit loop)
- LeetCode #190 — Reverse Bits (per-position bit manipulation)
- LeetCode #338 — Counting Bits (`dp[n] = dp[n & (n-1)] + 1`)
