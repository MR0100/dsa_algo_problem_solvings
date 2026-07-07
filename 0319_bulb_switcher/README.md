# 0319 — Bulb Switcher

> LeetCode #319 · Difficulty: Medium
> **Categories:** Math, Number Theory, Brainteaser

---

## Problem Statement

There are `n` bulbs that are initially off. You first turn on all the bulbs,
then you turn off every second bulb.

On the third round, you toggle every third bulb (turning on if it's off or
turning off if it's on). For the `i`-th round, you toggle every `i` bulb. For
the `n`-th round, you only toggle the last bulb.

Return the number of bulbs that are on after `n` rounds.

**Example 1:**

```
Input:  n = 3
Output: 1
```

Explanation:
```
At first, the three bulbs are [off, off, off].
After the first round, the three bulbs are [on, on, on].
After the second round, the three bulbs are [on, off, on].
After the third round, the three bulbs are [on, off, off].
So you should return 1 because there is only one bulb is on.
```

**Example 2:**

```
Input:  n = 0
Output: 0
```

**Example 3:**

```
Input:  n = 1
Output: 1
```

**Constraints:**

- `0 <= n <= 10^9`

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Amazon    | ★★★☆☆ Medium     | 2023          |
| Google    | ★★★☆☆ Medium     | 2023          |
| Microsoft | ★★☆☆☆ Low        | 2022          |
| Adobe     | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Number Theory (divisor counting)** — a bulb ends on iff its index has an odd
  number of divisors, which happens only for perfect squares →
  see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)
- **Math / floor(sqrt(n))** — counting perfect squares ≤ n reduces to one square
  root → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force simulation | O(n log n) | O(n) | Verifying the pattern for small n |
| 2 | Count divisors / perfect squares | O(n) | O(1) | Explaining the "why perfect squares" insight |
| 3 | Integer square root (Optimal) | O(1) | O(1) | The real answer — floor(sqrt(n)) |

---

## Approach 1 — Brute Force Simulation

### Intuition
Follow the statement literally: keep a boolean per bulb and simulate every
round's toggles, then count what's left on.

### Algorithm
1. `bulbs` = `n` booleans, all `false` (off).
2. For `r = 1..n`: for `pos = r, 2r, 3r, ...` up to `n`: flip `bulbs[pos-1]`.
3. Count `true` entries.

### Complexity
- **Time:** O(n log n) — total toggles = `n/1 + n/2 + ... + n/n ≈ n·Hₙ`.
- **Space:** O(n) — the bulb array.

### Code
```go
func bruteForce(n int) int {
	bulbs := make([]bool, n)
	for r := 1; r <= n; r++ {
		for pos := r; pos <= n; pos += r {
			bulbs[pos-1] = !bulbs[pos-1]
		}
	}
	count := 0
	for _, on := range bulbs {
		if on {
			count++
		}
	}
	return count
}
```

### Dry Run
Input `n = 3`. Bulbs indexed 1..3.

| Round r | toggles positions | bulbs after (1,2,3) |
|---------|-------------------|---------------------|
| start   | —                 | off, off, off       |
| 1       | 1,2,3             | on, on, on          |
| 2       | 2                 | on, off, on         |
| 3       | 3                 | on, off, off        |

On count = 1. ✓

---

## Approach 2 — Count Divisors (Perfect Squares)

### Intuition
Bulb `i` is toggled once per divisor of `i`: round `r` toggles bulb `i` exactly
when `r | i`. It ends ON iff toggled an odd number of times, i.e. iff `i` has an
odd number of divisors. Divisors pair up as `(d, i/d)`; the only time a divisor
is unpaired (`d == i/d`) is when `i` is a **perfect square**. So exactly the
perfect-square-indexed bulbs stay on.

### Algorithm
1. For `i = 1..n`: check whether `i` is a perfect square; if so, count it.

### Complexity
- **Time:** O(n) — a constant-time perfect-square test per bulb.
- **Space:** O(1).

### Code
```go
func countDivisors(n int) int {
	count := 0
	for i := 1; i <= n; i++ {
		root := int(math.Sqrt(float64(i)))
		if root*root == i {
			count++
		}
	}
	return count
}
```

### Dry Run
Input `n = 3`.

| i | sqrt(i) rounded | root² == i? | perfect square? |
|---|-----------------|-------------|-----------------|
| 1 | 1               | 1 == 1 ✓    | yes → count=1   |
| 2 | 1               | 1 != 2      | no              |
| 3 | 1               | 1 != 3      | no              |

Count = 1. ✓ (For `n=10`, the perfect squares are 1, 4, 9 → count 3.)

---

## Approach 3 — Integer Square Root (Optimal)

### Intuition
Only perfect-square-indexed bulbs stay on, so the answer is the count of perfect
squares in `[1, n]`, which is exactly `⌊√n⌋` (the squares `1², 2², ..., ⌊√n⌋²`).
This is a single square root — O(1).

### Algorithm
1. Compute `root = floor(sqrt(n))`, correcting for floating-point error so that
   `root² ≤ n < (root+1)²`.
2. Return `root`.

### Complexity
- **Time:** O(1).
- **Space:** O(1).

### Code
```go
func integerSqrt(n int) int {
	root := int(math.Sqrt(float64(n)))
	for root*root > n {
		root--
	}
	for (root+1)*(root+1) <= n {
		root++
	}
	return root
}
```

### Dry Run
Input `n = 3`.

| Step | Value |
|------|-------|
| `math.Sqrt(3)` | ≈ 1.732 → `root = 1` |
| `root² > n`? | `1 > 3`? no |
| `(root+1)² ≤ n`? | `4 ≤ 3`? no |
| return | `1` |

Answer = 1. ✓ (For `n = 99999999`, sqrt ≈ 9999.9999 → root 9999, guard confirms
`9999² = 99980001 ≤ n` and `10000² > n`, so 9999.)

---

## Key Takeaways

- **Odd divisor count ⇔ perfect square.** The classic number-theory fact: every
  divisor `d < √i` pairs with `i/d > √i`; only `√i` is self-paired.
- **"Toggle every k-th" problems** map to divisor-parity questions — reframe
  "how many times is index i touched?" as "how many divisors does i have?".
- **Collapse simulation to a formula.** Recognizing the invariant turns an
  O(n log n) simulation into an O(1) `floor(sqrt(n))`.
- **Guard `int(math.Sqrt(...))` for large n** — floating point can be off by
  one near perfect squares; verify with integer multiplication.

---

## Related Problems

- LeetCode #672 — Bulb Switcher II (parity brainteaser)
- LeetCode #1375 — Bulb Switcher III
- LeetCode #367 — Valid Perfect Square (integer sqrt)
- LeetCode #69 — Sqrt(x) (integer square root)
