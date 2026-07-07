# 0172 — Factorial Trailing Zeroes

> LeetCode #172 · Difficulty: Medium
> **Categories:** Math

---

## Problem Statement

Given an integer `n`, return *the number of trailing zeroes in `n!`*.

Note that `n! = n * (n - 1) * (n - 2) * ... * 3 * 2 * 1`.

**Example 1:**

```
Input: n = 3
Output: 0
Explanation: 3! = 6, no trailing zero.
```

**Example 2:**

```
Input: n = 5
Output: 1
Explanation: 5! = 120, one trailing zero.
```

**Example 3:**

```
Input: n = 0
Output: 0
```

**Constraints:**

- `0 <= n <= 10^4`

**Follow-up:** Could you write a solution that works in logarithmic time complexity?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Bloomberg  | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Meta       | ★★☆☆☆ Low        | 2023          |
| Adobe      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Math / Number Theory (prime factorization, Legendre's formula)** — trailing zeros = factors of 10 = paired factors of 2 and 5; since 2s always outnumber 5s in n!, the answer is the exponent of 5 in n!, given by ⌊n/5⌋ + ⌊n/25⌋ + ⌊n/125⌋ + … → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (compute n! with `math/big`) | O(n² log n) bit-work | O(n log n) | Only to verify the math; n! overflows `int64` at n = 21 |
| 2 | Count Factors of 5 per Multiple | O(n) | O(1) | When you've spotted "count the 5s" but not yet the closed form |
| 3 | Logarithmic Division — Legendre's Formula (Optimal) | O(log₅ n) | O(1) | Always — answers the follow-up in ~5 lines |

---

## Approach 1 — Brute Force (Compute n! with Big Integers)

### Intuition

The most direct reading: build the factorial, print it, count the `'0'` characters at the end. The catch is size — n! overflows `int64` already at n = 21 (21! ≈ 5.1×10¹⁹), and 10000! has **35,660 digits** — so Go's `math/big` arbitrary-precision integers are mandatory. This works, but it computes tens of thousands of digits only to look at the last few.

### Algorithm

1. Initialize `f = 1` as a `big.Int`.
2. For `i` from 2 to n: `f = f × i`.
3. Convert `f` to its decimal string.
4. Scan the string from the right, counting characters while they equal `'0'`.
5. Return the count.

### Complexity

- **Time:** O(n² log n) bit-work — n big-integer multiplications; the k-th multiplication costs O(digits of k!) and the digit count grows to Θ(n log n).
- **Space:** O(n log n) — the digits of n! itself (35,660 digits for n = 10⁴).

### Code

```go
func bruteForce(n int) int {
	f := big.NewInt(1) // running factorial, arbitrary precision
	for i := 2; i <= n; i++ {
		f.Mul(f, big.NewInt(int64(i))) // f *= i, never overflows
	}
	s := f.String() // full decimal expansion of n!
	zeros := 0
	// Walk from the last digit backwards while we keep seeing '0'.
	for i := len(s) - 1; i >= 0 && s[i] == '0'; i-- {
		zeros++
	}
	return zeros
}
```

### Dry Run

Example 1: `n = 3`.

| Step | i | f (big.Int) after `f *= i` |
|------|---|----------------------------|
| 1 | init | 1 |
| 2 | 2 | 2 |
| 3 | 3 | 6 |

| Step | string s | scan position | char | zeros |
|------|----------|---------------|------|-------|
| 4 | `"6"` | index 0 (last) | `'6'` ≠ `'0'` → stop | 0 |

Result: `0` ✔

Bonus trace of Example 2 (`n = 5`): f grows 1→2→6→24→120; s = `"120"`; scan: index 2 = `'0'` → zeros = 1; index 1 = `'2'` → stop. Result: `1` ✔

---

## Approach 2 — Count Factors of 5 per Multiple

### Intuition

A trailing zero is a factor of 10 = 2×5. In the product 1·2·…·n, factors of 2 vastly outnumber factors of 5 (every 2nd number is even; only every 5th number contributes a 5), so **every 5 finds a 2 to pair with** — the number of trailing zeros equals the exponent of 5 in n!. So skip building n! entirely: visit each multiple of 5 up to n and count how many times 5 divides it (25 = 5² contributes two, 125 = 5³ three, …).

### Algorithm

1. Set `zeros = 0`.
2. For `i = 5, 10, 15, …, n` (only multiples of 5 carry any factor of 5):
   1. Copy `x = i`; while `x % 5 == 0`: increment `zeros`, set `x /= 5`.
3. Return `zeros`.

### Complexity

- **Time:** O(n) — n/5 multiples visited; the inner loop adds only n/25 + n/125 + … extra iterations, a geometric series that keeps the total linear.
- **Space:** O(1) — two counters.

### Code

```go
func countFactorsOfFive(n int) int {
	zeros := 0
	for i := 5; i <= n; i += 5 { // only multiples of 5 carry any factor of 5
		for x := i; x%5 == 0; x /= 5 {
			zeros++ // one zero per factor of 5 inside this multiple (25→2, 125→3, …)
		}
	}
	return zeros
}
```

### Dry Run

Example 1: `n = 3`.

| Step | i | condition `i <= n` | inner loop | zeros |
|------|---|--------------------|------------|-------|
| 1 | init | — | — | 0 |
| 2 | 5 | 5 ≤ 3 is false → outer loop never runs | — | 0 |

Result: `0` ✔

Bonus trace of `n = 25` (shows the double-count at 25):

| Step | i | x values in inner loop | 5s found | zeros |
|------|---|------------------------|----------|-------|
| 1 | 5 | 5 → 1 | 1 | 1 |
| 2 | 10 | 10 → 2 | 1 | 2 |
| 3 | 15 | 15 → 3 | 1 | 3 |
| 4 | 20 | 20 → 4 | 1 | 4 |
| 5 | 25 | 25 → 5 → 1 | **2** | 6 |

Result: `6` ✔ (25! = …4000000 with six trailing zeros)

---

## Approach 3 — Logarithmic Division — Legendre's Formula (Optimal)

### Intuition

Flip the counting direction. Instead of asking each number "how many 5s do you contribute?", count in layers: ⌊n/5⌋ numbers contribute *at least one* 5; of those, ⌊n/25⌋ contribute *a second* 5; ⌊n/125⌋ a *third*; and so on. Summing the layers counts every factor of 5 exactly once — this is Legendre's formula for the exponent of a prime p in n!: Σ ⌊n/pᵏ⌋. Each term is one integer division and terms shrink 5× per step, so the loop runs log₅(n) times — exactly the logarithmic solution the follow-up requests.

### Algorithm

1. Set `zeros = 0`.
2. While `n > 0`:
   1. `n /= 5` — n now holds ⌊original/5ᵏ⌋ after the k-th division.
   2. `zeros += n` — add the count of numbers carrying a k-th factor of 5.
3. Return `zeros`.

### Complexity

- **Time:** O(log₅ n) — one division per power of 5 that fits in n (≤ 6 iterations for n ≤ 10⁴, ≤ 14 for `int32`).
- **Space:** O(1) — a single counter.

### Code

```go
func logarithmicDivision(n int) int {
	zeros := 0
	for n > 0 {
		n /= 5     // n is now ⌊n/5⌋, ⌊n/25⌋, ⌊n/125⌋, ... on successive turns
		zeros += n // add how many numbers contribute yet another factor of 5
	}
	return zeros
}
```

### Dry Run

Example 1: `n = 3`.

| Step | n (loop entry) | n after `/= 5` | zeros |
|------|----------------|----------------|-------|
| 1 | 3 | 0 | 0 + 0 = 0 |
| 2 | 0 → loop exits | — | 0 |

Result: `0` ✔

Bonus trace of `n = 10000` (constraint upper bound):

| Step | n (loop entry) | n after `/= 5` | meaning | zeros |
|------|----------------|----------------|---------|-------|
| 1 | 10000 | 2000 | multiples of 5 | 2000 |
| 2 | 2000 | 400 | multiples of 25 | 2400 |
| 3 | 400 | 80 | multiples of 125 | 2480 |
| 4 | 80 | 16 | multiples of 625 | 2496 |
| 5 | 16 | 3 | multiples of 3125 | 2499 |
| 6 | 3 | 0 | multiples of 15625 (none) | 2499 |

Result: `2499` ✔

---

## Key Takeaways

- **Trailing zeros = min(exponent of 2, exponent of 5) in the factorization** — and in n! the 5s are always the bottleneck, so count only 5s. Stating *why* (2s outnumber 5s) is the interview checkpoint.
- **Legendre's formula** — exponent of prime p in n! is Σₖ ⌊n/pᵏ⌋ — is a reusable tool: it answers "how many times does p divide n!" for any p, e.g. binomial-coefficient divisibility problems.
- The `n /= 5; zeros += n` loop is the tightest way to code the series — no explicit powers, so no overflow from computing 5ᵏ.
- Don't forget squares and cubes: 25 contributes **two** 5s, 125 three. A plain `n/5` (single term) is the classic off-by-a-layer wrong answer.
- Know the brute-force pitfall cold: n! exceeds `int64` at n = 21, so any "just compute it" answer must mention big integers.

---

## Related Problems

- LeetCode #793 — Preimage Size of Factorial Zeroes Function (binary search over this exact function)
- LeetCode #204 — Count Primes (prime-factor counting mindset)
- LeetCode #50 — Pow(x, n) (another "make it logarithmic" follow-up)
- LeetCode #43 — Multiply Strings (what the brute force is really doing under the hood)
