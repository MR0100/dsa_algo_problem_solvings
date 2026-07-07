# 0231 — Power of Two

> LeetCode #231 · Difficulty: Easy
> **Categories:** Math, Bit Manipulation, Recursion

---

## Problem Statement

Given an integer `n`, return `true` if it is a power of two. Otherwise, return `false`.

An integer `n` is a power of two, if there exists an integer `x` such that `n == 2ˣ`.

**Example 1:**
```
Input: n = 1
Output: true
Explanation: 2^0 = 1
```

**Example 2:**
```
Input: n = 16
Output: true
Explanation: 2^4 = 16
```

**Example 3:**
```
Input: n = 3
Output: false
```

**Constraints:**
- `-2³¹ <= n <= 2³¹ - 1`

**Follow up:** Could you solve it without loops/recursion?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used
- **Bit Manipulation** — a power of two has exactly one set bit; `n & (n-1)` clears it → see [`/dsa/bit_manipulation.md`](/dsa/bit_manipulation.md)
- **Number Theory** — powers of two are the numbers whose only prime factor is 2 → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)

---

## Approaches Overview
| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Iterative Division | O(log n) | O(1) | Intuitive baseline; strip factors of 2 |
| 2 | Bit Count | O(log n) | O(1) | Emphasize the single-set-bit property |
| 3 | `n & (n-1)` (Optimal) | O(1) | O(1) | Loop-free follow-up answer |
| 4 | Divisor of Max Power | O(1) | O(1) | Loop-free via range assumption |

---

## Approach 1 — Iterative Division

### Intuition
A power of two is `2^k`, whose only prime factor is 2. Repeatedly dividing a positive number by 2 while it stays even reduces a true power of two exactly to 1; anything with another prime factor gets stuck on an odd number greater than 1.

### Algorithm
1. Powers of two are positive, so reject `n <= 0`.
2. While `n` is even, divide `n` by 2.
3. Return `true` iff the leftover equals 1.

### Complexity
- **Time:** O(log n) — at most `log₂(n)` halvings.
- **Space:** O(1) — one mutable integer.

### Code
```go
func iterativeDivision(n int) bool {
	if n <= 0 { // powers of two (2^k, k>=0) are all >= 1, so non-positive fails
		return false
	}
	for n%2 == 0 { // strip a factor of 2 as long as one is present
		n /= 2
	}
	return n == 1 // only a pure power of two reduces exactly to 1
}
```

### Dry Run
Trace `n = 16`:

| Step | n before | n % 2 == 0? | n after |
|------|----------|-------------|---------|
| init | 16       | —           | 16      |
| 1    | 16       | yes         | 8       |
| 2    | 8        | yes         | 4       |
| 3    | 4        | yes         | 2       |
| 4    | 2        | yes         | 1       |
| 5    | 1        | no (loop ends) | 1    |

Final `n == 1` → return `true`.

---

## Approach 2 — Bit Count

### Intuition
In binary, `2^k` is a single 1 followed by k zeros (`1`, `10`, `100`, …). So a positive number is a power of two exactly when its binary representation has one set bit. Count the set bits and compare with 1.

### Algorithm
1. Reject `n <= 0`.
2. Walk the bits, counting how many are 1.
3. Return `true` iff the count is exactly 1.

### Complexity
- **Time:** O(log n) — one pass over `~log₂(n)` bits.
- **Space:** O(1) — a counter.

### Code
```go
func bitCount(n int) bool {
	if n <= 0 { // negatives and zero are never powers of two
		return false
	}
	count := 0
	for n > 0 { // examine each bit from least significant upward
		count += n & 1 // add 1 if the current lowest bit is set
		n >>= 1        // shift to inspect the next bit
	}
	return count == 1 // a lone set bit means n == 2^k
}
```

### Dry Run
Trace `n = 16` (binary `10000`):

| Step | n (binary) | n & 1 | count | n >> 1 |
|------|------------|-------|-------|--------|
| 1    | 10000      | 0     | 0     | 1000   |
| 2    | 1000       | 0     | 0     | 100    |
| 3    | 100        | 0     | 0     | 10     |
| 4    | 10         | 0     | 0     | 1      |
| 5    | 1          | 1     | 1     | 0      |

Loop ends, `count == 1` → return `true`.

---

## Approach 3 — Brian Kernighan / `n & (n-1)` (Optimal)

### Intuition
Subtracting 1 from a power of two flips its single set bit to 0 and turns every lower bit to 1 (`1000 - 1 = 0111`). ANDing the two yields 0. Any number with two or more set bits keeps its higher bits after the AND, giving a non-zero result. So `n > 0 && n & (n-1) == 0` characterises powers of two in O(1) — the loop-free follow-up answer.

### Algorithm
1. Reject `n <= 0`.
2. Return `true` iff `n & (n-1) == 0`.

### Complexity
- **Time:** O(1) — one subtraction and one AND.
- **Space:** O(1).

### Code
```go
func bitTrick(n int) bool {
	// n>0 rules out zero/negatives; n&(n-1)==0 means n has a single set bit.
	return n > 0 && n&(n-1) == 0
}
```

### Dry Run
Trace `n = 16`:

| Expression | Value |
|------------|-------|
| n          | `10000` (16) |
| n - 1      | `01111` (15) |
| n & (n-1)  | `00000` (0)  |
| n > 0      | true |
| result     | `true` (both conditions hold) |

For `n = 3`: `n-1 = 2` = `10`, `n & (n-1) = 11 & 10 = 10 = 2 ≠ 0` → `false`.

---

## Approach 4 — Divisor of a Max Power of Two

### Intuition
Within a fixed integer width, every power of two divides the largest power of two that fits. For 32-bit signed ints that is `2^30`. Any power of two `n` in `[1, 2^30]` divides `2^30` evenly; a non-power never does. That reduces the test to one modulo — also loop-free — at the cost of relying on the value range.

### Algorithm
1. Reject `n <= 0`.
2. Let `maxPow = 2^30`.
3. Return `true` iff `maxPow % n == 0`.

### Complexity
- **Time:** O(1) — one modulo.
- **Space:** O(1).

### Code
```go
func maxPowerDivisor(n int) bool {
	if n <= 0 { // guard against non-positive and avoid % by zero
		return false
	}
	const maxPow = 1 << 30 // 2^30, the largest power of two < 2^31 (int32 range)
	return maxPow%n == 0    // every power of two in range divides 2^30 exactly
}
```

### Dry Run
Trace `n = 16`, `maxPow = 2^30 = 1073741824`:

| Step | Expression | Value |
|------|------------|-------|
| 1    | n <= 0     | false → continue |
| 2    | maxPow % n | `1073741824 % 16 = 0` |
| 3    | result     | `0 == 0` → `true` |

For `n = 3`: `1073741824 % 3 = 1 ≠ 0` → `false`.

---

## Key Takeaways
- `n & (n-1)` clears the lowest set bit — the single most useful bit trick for "power of two", "count set bits", and "is only one bit set" questions.
- A power of two ⇔ exactly one set bit ⇔ only prime factor is 2. Three equivalent lenses, three different solutions.
- Always guard `n <= 0` first: negatives and zero are never powers of two, and modulo needs a non-zero divisor.
- "Without loops or recursion" almost always hints at a constant-time bit or arithmetic identity.

## Related Problems
- LeetCode #191 — Number of 1 Bits (same `n & (n-1)` trick)
- LeetCode #326 — Power of Three (analogous, but 3 has no bit shortcut)
- LeetCode #342 — Power of Four (power of two plus a bit-mask check)
- LeetCode #338 — Counting Bits (builds on single-set-bit reasoning)
