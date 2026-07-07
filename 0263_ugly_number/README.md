# 0263 — Ugly Number

> LeetCode #263 · Difficulty: Easy
> **Categories:** Math, Number Theory

---

## Problem Statement

An **ugly number** is a positive integer which does not have a prime factor other than `2`, `3`, and `5`.

Given an integer `n`, return `true` if `n` is an ugly number.

**Example 1:**

```
Input: n = 6
Output: true
Explanation: 6 = 2 × 3
```

**Example 2:**

```
Input: n = 1
Output: true
Explanation: 1 has no prime factors, therefore all of its prime factors are limited to 2, 3, and 5.
```

**Example 3:**

```
Input: n = 14
Output: false
Explanation: 14 is not ugly since it includes the prime factor 7.
```

**Constraints:**

- `-2³¹ <= n <= 2³¹ - 1`

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Amazon    | ★★★☆☆ Medium     | 2023          |
| Microsoft | ★★★☆☆ Medium     | 2023          |
| Google    | ★★☆☆☆ Low        | 2022          |
| Adobe     | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Prime factorisation / trial division** — strip the allowed primes and check the remainder → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)
- **Recursion** — express the divide-out as a self-call on the quotient → see [`/dsa/divide_and_conquer.md`](/dsa/divide_and_conquer.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Trial Division (Optimal) | O(log n) | O(1) | Default — iterative, constant space |
| 2 | Recursion | O(log n) | O(log n) | Cleaner base-case framing; uses stack |

---

## Approach 1 — Trial Division (Optimal)

### Intuition
If `n` is built only from the primes 2, 3 and 5, then dividing out every one of those factors leaves exactly `1`. If `n` contains any other prime factor (like 7), that factor cannot be removed and the leftover is greater than 1.

### Algorithm
1. If `n <= 0`, return `false` (ugly numbers are positive).
2. For each `p` in `{2, 3, 5}`: while `n % p == 0`, set `n /= p`.
3. Return `true` iff the remaining `n == 1`.

### Complexity
- **Time:** O(log n) — every division at least halves `n` (the factor-2 loop dominates).
- **Space:** O(1) — only the running integer `n`.

### Code
```go
func trialDivision(n int) bool {
	if n <= 0 { // 0 and negatives are never ugly
		return false
	}
	for _, p := range []int{2, 3, 5} { // divide out each allowed prime fully
		for n%p == 0 { // while p divides n evenly
			n /= p // remove one factor of p
		}
	}
	return n == 1 // only 2/3/5 factors ⇒ nothing but 1 is left
}
```

### Dry Run
Example 1: `n = 6`.

| Prime p | Loop action        | n after |
|---------|--------------------|---------|
| 2       | 6 % 2 == 0 → 6/2   | 3       |
| 2       | 3 % 2 != 0 → stop  | 3       |
| 3       | 3 % 3 == 0 → 3/3   | 1       |
| 3       | 1 % 3 != 0 → stop  | 1       |
| 5       | 1 % 5 != 0 → stop  | 1       |

Final `n == 1` → return `true`.

---

## Approach 2 — Recursion

### Intuition
`n` is ugly iff it is divisible by one of 2/3/5 **and** the quotient is also ugly, bottoming out at `1`. Same divide-out logic, framed recursively.

### Algorithm
1. Base cases: `n <= 0` → `false`; `n == 1` → `true`.
2. If `n` is divisible by 2, 3, or 5, recurse on `n / p`.
3. If divisible by none of them, `n` has a forbidden factor → `false`.

### Complexity
- **Time:** O(log n) — one prime factor removed per call.
- **Space:** O(log n) — recursion stack equals the number of prime factors.

### Code
```go
func recursive(n int) bool {
	if n <= 0 { // non-positive ⇒ not ugly
		return false
	}
	if n == 1 { // fully reduced ⇒ ugly
		return true
	}
	for _, p := range []int{2, 3, 5} {
		if n%p == 0 { // divisible by an allowed prime
			return recursive(n / p) // recurse on the quotient
		}
	}
	return false // divisible by none of 2/3/5 ⇒ has a forbidden factor
}
```

### Dry Run
Example 1: `n = 6`.

| Call         | n  | divisible by | recurse on |
|--------------|----|--------------|------------|
| recursive(6) | 6  | 2            | 3          |
| recursive(3) | 3  | 3            | 1          |
| recursive(1) | 1  | base case    | → true     |

Returns `true`.

---

## Key Takeaways

- **Divide out the allowed primes; the residue tells the story.** Residue `1` → ugly; residue `> 1` → a forbidden prime factor remains.
- Guard `n <= 0` first: negatives and zero are never ugly, and 1 is ugly by the "no prime factors" convention.
- This factor-stripping template generalises to "is `n` a product of only primes in set S" problems.

---

## Related Problems

- LeetCode #264 — Ugly Number II (find the nth ugly number via DP / heap)
- LeetCode #313 — Super Ugly Number (arbitrary prime set)
- LeetCode #204 — Count Primes (sieve-based factorisation)
- LeetCode #1201 — Ugly Number III (counting with inclusion–exclusion)
