# 0371 — Sum of Two Integers

> LeetCode #371 · Difficulty: Medium
> **Categories:** Bit Manipulation, Math

---

## Problem Statement

Given two integers `a` and `b`, return *the sum of the two integers without using the operators* `+` *and* `-`.

**Example 1:**

```
Input: a = 1, b = 2
Output: 3
```

**Example 2:**

```
Input: a = 2, b = 3
Output: 5
```

**Constraints:**

- `-1000 <= a, b <= 1000`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Bit Manipulation** — XOR is addition-without-carry, AND finds carry bits, and `<<1` shifts the carry into the next column → see [`/dsa/bit_manipulation.md`](/dsa/bit_manipulation.md)
- **Two's Complement Arithmetic** — negative numbers and wraparound handled by working in unsigned space → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Bitwise Iterative (carry loop) | O(1) | O(1) | Standard, no recursion overhead |
| 2 | Bitwise Recursive | O(1) | O(1) | Cleaner expression of the same identity |

---

## Approach 1 — Bitwise Iterative

### Intuition
Adding two bits produces a **sum bit** (XOR: `1^1=0`, `1^0=1`, `0^0=0`) and a **carry bit** (AND: only `1&1=1`). Shifting the carry left by one lines it up with the next column, exactly like grade-school addition. Feeding `(sum, carry)` back into the same rule until the carry is 0 gives the full sum. Working in `uint` gives clean two's-complement wraparound so negatives just work.

### Algorithm
1. While `b != 0`:
   1. `carry = (a & b) << 1` — bits where both are 1, moved to the next column.
   2. `a = a ^ b` — sum ignoring carry.
   3. `b = carry` — carry becomes the new addend.
2. Return `a`.

### Complexity
- **Time:** O(1) — at most 64 iterations (one per bit position), independent of value magnitude.
- **Space:** O(1) — two scalars.

### Code
```go
func bitwiseIterative(a int, b int) int {
	ua, ub := uint(a), uint(b) // work in unsigned to get clean two's-complement wrap
	for ub != 0 {              // loop until nothing left to carry
		carry := (ua & ub) << 1 // columns where both bits are 1 carry into the next
		ua = ua ^ ub            // add the two numbers ignoring carry
		ub = carry              // the carry becomes the next thing to add
	}
	return int(ua) // reinterpret the bit pattern as a signed int
}
```

### Dry Run
`a = 2 (010)`, `b = 3 (011)`:

| Step | ua (before) | ub (before) | carry=(ua&ub)<<1 | ua=ua^ub | ub=carry |
|------|-------------|-------------|------------------|----------|----------|
| 1 | 010 (2) | 011 (3) | (010)<<1 = 100 (4) | 001 (1) | 100 (4) |
| 2 | 001 (1) | 100 (4) | (000)<<1 = 000 (0) | 101 (5) | 000 (0) |
| 3 | 101 (5) | 000 (0) | loop exits (ub==0) | — | — |

Return `5`. ✓

---

## Approach 2 — Bitwise Recursive

### Intuition
The same identity, recursive: `getSum(a, b) = getSum(a^b, (a&b)<<1)`. Each call pushes the remaining carry one column left; since the carry only moves upward and eventually falls off the top of the word, the base case `b == 0` is always reached.

### Algorithm
1. If `b == 0`, return `a` (no carry left).
2. Otherwise return `bitwiseRecursive(a^b, (a&b)<<1)`.

### Complexity
- **Time:** O(1) — recursion depth bounded by word size (≤ 64 calls).
- **Space:** O(1) — call stack bounded by word size.

### Code
```go
func bitwiseRecursive(a int, b int) int {
	if b == 0 { // no carry remains → a already holds the full sum
		return a
	}
	ua, ub := uint(a), uint(b)                           // unsigned for clean wraparound
	return bitwiseRecursive(int(ua^ub), int((ua&ub)<<1)) // sum-without-carry, carry
}
```

### Dry Run
`a = 2 (010)`, `b = 3 (011)`:

| Call | a | b | a^b | (a&b)<<1 | Next call |
|------|---|---|-----|----------|-----------|
| 1 | 010 (2) | 011 (3) | 001 (1) | 100 (4) | `f(1, 4)` |
| 2 | 001 (1) | 100 (4) | 101 (5) | 000 (0) | `f(5, 0)` |
| 3 | 101 (5) | 000 (0) | — | — | `b==0 → return 5` |

Return `5`. ✓

---

## Key Takeaways
- **XOR = add-without-carry, AND<<1 = carry.** This is the fundamental "full adder" identity for building arithmetic from bit ops.
- **Loop/recurse until the carry is zero** — the carry chain always terminates because carries only move to higher bits.
- **Negative numbers need no special case** if you compute in unsigned and rely on two's-complement wraparound (in Go, cast to `uint` then back to `int`).

---

## Related Problems
- LeetCode #67 — Add Binary (add without native + on bit strings)
- LeetCode #2 — Add Two Numbers (carry propagation on linked lists)
- LeetCode #29 — Divide Two Integers (arithmetic without the forbidden operator)
- LeetCode #190 — Reverse Bits (bit-level manipulation)
