# 0326 — Power of Three

> LeetCode #326 · Difficulty: Easy
> **Categories:** Math, Recursion

---

## Problem Statement

Given an integer `n`, return `true` if it is a power of three. Otherwise, return `false`.

An integer `n` is a power of three, if there exists an integer `x` such that `n == 3^x`.

**Example 1:**

```
Input: n = 27
Output: true
Explanation: 27 = 3^3
```

**Example 2:**

```
Input: n = 0
Output: false
Explanation: There is no x where 3^x = 0.
```

**Example 3:**

```
Input: n = -1
Output: false
Explanation: There is no x where 3^x = (-1).
```

**Constraints:**

- `-2^31 <= n <= 2^31 - 1`

**Follow-up:** Could you solve it without loops/recursion?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |
| Adobe      | ★☆☆☆☆ Rare       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Number Theory (prime factorisation & divisibility)** — a power of three has 3 as its *only* prime factor; the no-loop trick leans on the fact that the divisors of `3^19` (a prime power) are exactly the powers of three, so a simple modulo test decides membership → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)
- **Logarithms / change of base** — `n == 3^x` iff `log₃(n)` is an integer; combined with floating-point tolerance handling → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Loop Division | O(log₃ n) | O(1) | The obvious, robust baseline; no tricks, no FP worries |
| 2 | Logarithm | O(1) | O(1) | Constant time, but needs careful epsilon handling |
| 3 | Integer Limit (No Loops — Optimal) | O(1) | O(1) | Answers the follow-up: no loop/recursion, one modulo |

---

## Approach 1 — Loop Division

### Intuition

`3^x` is nothing but `3` multiplied into `1` exactly `x` times. Undo that: keep dividing `n` by `3` as long as it divides evenly. A genuine power of three collapses all the way down to `1`; anything carrying a foreign factor leaves a non-zero remainder at some step, which means the leftover value stalls above `1`. Non-positive inputs are rejected up front because `3^x` is always at least `1`.

### Algorithm

1. If `n <= 0`, return `false` (every power of three is `≥ 1`).
2. While `n % 3 == 0`, do `n /= 3` (peel off one factor of three).
3. After the loop, return `true` iff `n == 1`.

### Complexity

- **Time:** O(log₃ n) — each iteration divides `n` by 3, so the number of steps is the exponent, at most ~19 for 32-bit inputs.
- **Space:** O(1) — mutates `n` in place, no extra storage.

### Code

```go
func loopDivision(n int) bool {
	if n <= 0 {
		return false // powers of three are always positive (3^0 = 1)
	}
	for n%3 == 0 { // peel off one factor of 3 per iteration while it divides cleanly
		n /= 3
	}
	return n == 1 // fully reduced to 1 ⇒ it was purely 3^x
}
```

### Dry Run

Example 1: `n = 27`.

| Step | n before | n <= 0? | n % 3 == 0? | Action | n after |
|------|----------|---------|-------------|--------|---------|
| 0 | 27 | no | — | initial guard passes | 27 |
| 1 | 27 | — | yes (27 % 3 = 0) | `n /= 3` | 9 |
| 2 | 9  | — | yes (9 % 3 = 0)  | `n /= 3` | 3 |
| 3 | 3  | — | yes (3 % 3 = 0)  | `n /= 3` | 1 |
| 4 | 1  | — | no (1 % 3 = 1)   | exit loop | 1 |

Final check: `n == 1` → `true` ✔

---

## Approach 2 — Logarithm

### Intuition

If `n == 3^x`, then taking a log gives `x = log₃(n)`, and `x` must be a whole number. Compute `log₃(n)` via change of base (`log(n) / log(3)`), then test whether the result is an integer. Floating-point logs are inexact, so instead of demanding exact equality we round to the nearest integer and confirm the gap is smaller than a tiny epsilon. Non-positive `n` is rejected before touching the log at all.

### Algorithm

1. If `n <= 0`, return `false`.
2. Compute `x = log₁₀(n) / log₁₀(3)` (the base cancels, giving `log₃(n)`).
3. Return `true` iff `|x − round(x)| < 1e-10`.

### Complexity

- **Time:** O(1) — two logarithm calls plus a round and a compare, all constant time.
- **Space:** O(1) — a couple of scalar floats.

### Code

```go
func logarithm(n int) bool {
	if n <= 0 {
		return false // logarithm is undefined / meaningless for n ≤ 0
	}
	x := math.Log10(float64(n)) / math.Log10(3) // change-of-base: log₃(n)
	// A true power of three gives an integer x; FP error means we test nearness
	// to the nearest integer instead of exact equality.
	return math.Abs(x-math.Round(x)) < 1e-10
}
```

### Dry Run

Example 1: `n = 27`.

| Step | Expression | Value |
|------|------------|-------|
| 1 | `n <= 0?` | no (27 > 0) |
| 2 | `log₁₀(27)` | ≈ 1.4313637642 |
| 3 | `log₁₀(3)` | ≈ 0.4771212547 |
| 4 | `x = log₁₀(27) / log₁₀(3)` | ≈ 3.0000000000 |
| 5 | `round(x)` | 3 |
| 6 | `|x − round(x)|` | ≈ 2e-16 < 1e-10 → `true` |

Result: `true` ✔ (`log₃(27) = 3`, an integer).

---

## Approach 3 — Integer Limit (No Loops — Optimal)

### Intuition

This answers the follow-up: no loop, no recursion. The key fact is that `3` is prime, so the only divisors of `3^k` are `1, 3, 9, …, 3^k` — nothing else can divide a prime power. Within a signed 32-bit integer, the biggest power of three is `3^19 = 1162261467` (`3^20` overflows). So for any valid `n`, `n` is a power of three **iff** `n` is positive and `n` divides `1162261467` exactly. One comparison and one modulo settle it.

### Algorithm

1. Let `MAX = 1162261467 = 3^19`, the largest power of three that fits in int32.
2. Return `true` iff `n > 0` **and** `MAX % n == 0`.

### Complexity

- **Time:** O(1) — a single positivity check and a single modulo, no iteration.
- **Space:** O(1) — one constant.

### Code

```go
func integerLimit(n int) bool {
	const maxPow3 = 1162261467      // 3^19, the largest power of three within int32
	return n > 0 && maxPow3%n == 0 // divisors of 3^19 (a prime power) are exactly the powers of three
}
```

### Dry Run

Example 1: `n = 27`.

| Step | Expression | Value |
|------|------------|-------|
| 1 | `maxPow3` | 1162261467 (= 3^19) |
| 2 | `n > 0?` | yes (27 > 0) |
| 3 | `1162261467 % 27` | 0 (since 1162261467 = 27 × 43046721 = 3^19) |
| 4 | `n > 0 && 1162261467 % 27 == 0` | `true` |

Result: `true` ✔ — `27` divides `3^19`, so it must itself be a power of three.

Contrast with `n = 45`: `1162261467 % 45 = 27 ≠ 0` → `false` (45 = 9 × 5 carries the foreign factor 5).

---

## Key Takeaways

- **A power of a prime `p` is characterised purely by divisibility:** `n` is a power of `p` iff `n > 0` and `n` divides the largest power of `p` in the integer range. This collapses the whole problem to one modulo — the same trick works for #231 (Power of Two, via `2^30`) and #342 conceptually.
- **Guard the sign first.** `3^x ≥ 1` always, so `n <= 0` is an instant `false`; forgetting this is the classic bug (e.g. `-1`, `0`, or negatives that still satisfy a naive modulo).
- **Logs give O(1) but demand epsilon discipline.** `log₃(n)` being "an integer" must be tested with a tolerance, never `==`, because floating-point rounding can nudge a true integer off by ~1e-16.
- **Loop division is the safe default** — no overflow, no FP, easy to reason about — but the integer-limit trick is what an interviewer means by "without loops or recursion".

---

## Related Problems

- LeetCode #231 — Power of Two (same idea; `n & (n-1) == 0` bit trick or divide-out-2s)
- LeetCode #342 — Power of Four (power-of-two that also passes a mod-3 / bit-mask test)
- LeetCode #263 — Ugly Number (divide out prime factors 2, 3, 5)
- LeetCode #1780 — Check if Number is a Sum of Powers of Three
- LeetCode #50 — Pow(x, n) (fast exponentiation, the inverse operation)
