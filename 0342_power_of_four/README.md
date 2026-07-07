# 0342 — Power of Four

> LeetCode #342 · Difficulty: Easy
> **Categories:** Math, Bit Manipulation, Recursion

---

## Problem Statement

Given an integer `n`, return `true` if it is a power of four. Otherwise, return `false`.

An integer `n` is a power of four, if there exists an integer `x` such that `n == 4^x`.

**Example 1:**

```
Input: n = 16
Output: true
```

**Example 2:**

```
Input: n = 5
Output: false
```

**Example 3:**

```
Input: n = 1
Output: true
```

**Constraints:**

- `-2^31 <= n <= 2^31 - 1`

**Follow up:** Could you solve it without loops/recursion?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |
| TikTok     | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Bit Manipulation** — a power of four is a single set bit at an *even* position; `n & (n-1) == 0` checks single-bit, and a mask `0x55555555` checks even position → see [`/dsa/bit_manipulation.md`](/dsa/bit_manipulation.md)
- **Number Theory** — the identity `4^k ≡ 1 (mod 3)` and logarithm tests come from basic number theory → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Iterative Division | O(log₄ n) | O(1) | Simple, obviously correct baseline |
| 2 | Bit Manipulation (Optimal) | O(1) | O(1) | Answers the "no loops" follow-up cleanly |
| 3 | Modulo-3 Property | O(1) | O(1) | Elegant number-theory one-liner |
| 4 | Logarithm Check | O(1) | O(1) | Shows the log idea; needs float-precision care |

---

## Approach 1 — Iterative Division

### Intuition
`n` is a power of four iff you can divide it by 4 repeatedly, with zero remainder each time, until you reach exactly 1.

### Algorithm
1. Reject `n <= 0`.
2. While `n % 4 == 0`, do `n /= 4`.
3. Return `n == 1`.

### Complexity
- **Time:** O(log₄ n) — one division per factor of four.
- **Space:** O(1).

### Code
```go
func iterativeDivision(n int) bool {
	if n <= 0 {
		return false
	}
	for n%4 == 0 {
		n /= 4
	}
	return n == 1
}
```

### Dry Run
Input `n = 16`:

| Step | n | `n % 4 == 0`? | Action |
|------|---|---------------|--------|
| 0 | 16 | yes | n = 16/4 = 4 |
| 1 | 4 | yes | n = 4/4 = 1 |
| 2 | 1 | no (1%4=1) | exit loop |

`n == 1` → return **true**.

---

## Approach 2 — Bit Manipulation (Optimal)

### Intuition
Powers of four in binary are `1, 100, 10000, …` — a single `1` bit sitting at an **even** index (0, 2, 4, …). Two checks capture this: `n & (n-1) == 0` means exactly one bit is set (power of two); `n & 0x55555555 != 0` means that bit is at an even position (the mask `0101…0101` has ones only at even indices).

### Algorithm
1. `n > 0`.
2. `n & (n-1) == 0` — single set bit.
3. `n & 0x55555555 != 0` — that bit is at an even index.

### Complexity
- **Time:** O(1) — constant number of bit ops.
- **Space:** O(1).

### Code
```go
func bitTrick(n int) bool {
	return n > 0 && n&(n-1) == 0 && n&0x55555555 != 0
}
```

### Dry Run
Input `n = 16` = `0b10000`:

| Check | Computation | Result |
|-------|-------------|--------|
| `n > 0` | 16 > 0 | true |
| `n & (n-1)` | `10000 & 01111 = 0` | single bit ✓ |
| `n & 0x55555555` | bit 4 is even; mask has 1 at bit 4 → non-zero | even position ✓ |

All pass → **true**. (For `n = 8` = `1000`, bit 3 is odd, so `8 & 0x55555555 == 0` → false.)

---

## Approach 3 — Modulo-3 Property

### Intuition
`4 ≡ 1 (mod 3)`, so `4^k ≡ 1 (mod 3)` for all k. A power of two that is *not* a power of four is `2^odd`, and `2^odd ≡ 2 (mod 3)`. So among powers of two, exactly the powers of four leave remainder 1 mod 3.

### Algorithm
1. `n > 0` and `n & (n-1) == 0` (power of two).
2. `n % 3 == 1`.

### Complexity
- **Time:** O(1).
- **Space:** O(1).

### Code
```go
func moduloThree(n int) bool {
	return n > 0 && n&(n-1) == 0 && n%3 == 1
}
```

### Dry Run
Input `n = 16`:

| Check | Computation | Result |
|-------|-------------|--------|
| `n > 0` | 16 > 0 | true |
| `n & (n-1)` | 0 | power of two ✓ |
| `n % 3` | 16 % 3 = 1 | ✓ |

All pass → **true**. (`n = 8`: `8 % 3 = 2` → false; `n = 5`: not power of two → false.)

---

## Approach 4 — Logarithm Check

### Intuition
`n = 4^k` iff `log(n)/log(4)` is a whole number. Floating-point rounding forces us to compare against the rounded value and re-verify, rather than test equality directly.

### Algorithm
1. `n > 0`.
2. `x = log(n)/log(4)`; round to `r`.
3. Accept iff `x` is within tolerance of `r` and `4^r == n`.

### Complexity
- **Time:** O(1).
- **Space:** O(1).

### Code
```go
func logCheck(n int) bool {
	if n <= 0 {
		return false
	}
	x := math.Log(float64(n)) / math.Log(4)
	r := math.Round(x)
	return math.Abs(x-r) < 1e-10 && int(math.Pow(4, r)) == n
}
```

### Dry Run
Input `n = 16`:

| Step | Value |
|------|-------|
| `x = log16/log4` | 2.0 |
| `r = round(x)` | 2 |
| `|x - r| < 1e-10` | true |
| `4^r == n` | `16 == 16` ✓ |

Return **true**. (`n = 5`: `x ≈ 1.16`, not near an integer → false.)

---

## Key Takeaways

- Power-of-two check is `n > 0 && n & (n-1) == 0`; power-of-four adds a constraint that the single bit is at an **even** position.
- The even-position mask `0x55555555` (`0101…`) is a reusable trick; its complement `0xAAAAAAAA` catches odd positions.
- `4^k ≡ 1 (mod 3)` is a neat number-theory shortcut; the general form `a^k ≡ (a mod m)^k (mod m)` recurs in these "power of X" problems.
- Prefer the O(1) bit/modulo checks for the "no loops/recursion" follow-up; avoid floats unless you re-verify, since `log` is imprecise.

---

## Related Problems

- LeetCode #231 — Power of Two (single-bit check, no even-position constraint)
- LeetCode #326 — Power of Three (mod / log variant with `3^k`)
- LeetCode #191 — Number of 1 Bits (popcount, same bit toolkit)
- LeetCode #29 — Divide Two Integers (bit-shift arithmetic)
