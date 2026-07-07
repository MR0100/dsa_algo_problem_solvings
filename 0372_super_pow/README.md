# 0372 — Super Pow

> LeetCode #372 · Difficulty: Medium
> **Categories:** Math, Divide and Conquer

---

## Problem Statement

Your task is to calculate `a^b mod 1337` where `a` is a positive integer and `b` is an extremely large positive integer given in the form of an array.

**Example 1:**

```
Input: a = 2, b = [3]
Output: 8
```

**Example 2:**

```
Input: a = 2, b = [1,0]
Output: 1024
```

**Example 3:**

```
Input: a = 1, b = [4,3,3,8,5,2]
Output: 1
```

**Constraints:**

- `1 <= a <= 2^31 - 1`
- `1 <= b.length <= 2000`
- `0 <= b[i] <= 9`
- `b` does not contain leading zeros.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Uber       | ★★☆☆☆ Low        | 2022          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Modular Exponentiation** — `(x*y) mod m = ((x mod m)*(y mod m)) mod m` keeps intermediate products bounded, and fast (binary) exponentiation raises to a power in O(log exp) → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)
- **Divide and Conquer on the Exponent** — peel one decimal digit at a time using `a^(10x+d) = (a^x)^10 · a^d` → see [`/dsa/divide_and_conquer.md`](/dsa/divide_and_conquer.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Digit-by-digit Horner expansion | O(n) | O(1) | Iterative, no stack growth |
| 2 | Recursive split on last digit | O(n) | O(n) | Cleanest expression of the recurrence |

*(n = len(b); each digit costs O(1) modular exponentiations with tiny fixed exponents.)*

---

## Approach 1 — Digit-by-Digit Horner Expansion

### Intuition
`b` is written out digit by digit, so use **Horner's method on the exponent**. If we already have `r = a^(prefix) mod 1337`, then appending digit `d` makes the exponent `prefix*10 + d`, and the new result is `r^10 · a^d (mod 1337)`. Start from `r = 1` (the empty-prefix `a^0`) and fold in each digit left to right. We never form the astronomically large `b`.

### Algorithm
1. `result = 1`; reduce `a` modulo 1337 once.
2. For each digit `d` in `b` (most significant first):
   - `result = powMod(result, 10) * powMod(a, d) mod 1337`.
3. Return `result`.

Here `powMod(base, exp)` is binary exponentiation under mod 1337.

### Complexity
- **Time:** O(n) — one iteration per digit; each does O(1) `powMod`s whose exponents (10 and ≤9) are constant.
- **Space:** O(1) — a handful of scalars.

### Code
```go
func superPowHorner(a int, b []int) int {
	a %= mod    // shrink the base first so products stay small
	result := 1 // a^0 for the empty prefix
	for _, d := range b {
		// Raise the running result to the 10th power (shift exponent one decimal
		// place) then multiply by a^d for the new least-significant digit.
		result = powMod(result, 10) * powMod(a, d) % mod
	}
	return result
}

func powMod(base, exp int) int {
	base %= mod   // keep base in range
	res := 1      // multiplicative identity
	for exp > 0 { // square-and-multiply over the bits of exp
		if exp&1 == 1 { // this bit is set → include current base power
			res = res * base % mod
		}
		base = base * base % mod // square the base for the next bit
		exp >>= 1                // move to the next higher bit
	}
	return res
}
```

### Dry Run
`a = 2`, `b = [3]` (so we want `2^3 mod 1337 = 8`):

| Step | digit d | result (before) | powMod(result,10) | powMod(a=2,d) | result = product mod 1337 |
|------|---------|-----------------|-------------------|---------------|---------------------------|
| init | — | 1 | — | — | 1 |
| 1 | 3 | 1 | 1^10 = 1 | 2^3 = 8 | 1 * 8 = **8** |

Return `8`. ✓

For `b = [1,0]` (`2^10 = 1024`): after digit `1`, result = `1^10 * 2^1 = 2`; after digit `0`, result = `2^10 * 2^0 = 1024 mod 1337 = 1024`. ✓

---

## Approach 2 — Recursive Split on Last Digit

### Intuition
Same math, expressed by peeling the **last** digit each call. With `b = [b0,…,bk]`, we have `b = 10 · [b0…b_{k-1}] + bk`, so `a^b = (a^[b0…b_{k-1}])^10 · a^bk`. Recurse on the shorter prefix, then combine. Base case: an empty exponent array is `a^0 = 1`.

### Algorithm
1. If `b` is empty, return 1.
2. `last = b[len-1]`, `prefix = b[:len-1]`.
3. `part1 = powMod(a, last)`.
4. `part2 = powMod(superPowRecursive(a, prefix), 10)`.
5. Return `part1 * part2 mod 1337`.

### Complexity
- **Time:** O(n) recursive calls, each O(log) for the fixed small exponents.
- **Space:** O(n) — recursion stack of depth n.

### Code
```go
func superPowRecursive(a int, b []int) int {
	if len(b) == 0 { // empty exponent means a^0 = 1
		return 1
	}
	last := b[len(b)-1]                               // least-significant digit
	prefix := b[:len(b)-1]                            // everything above it
	part1 := powMod(a, last)                          // a^(last digit)
	part2 := powMod(superPowRecursive(a, prefix), 10) // (a^prefix)^10
	return part1 * part2 % mod                        // combine under the modulus
}
```

### Dry Run
`a = 2`, `b = [3]`:

| Call | b | last | prefix | recurse(prefix) | part1=2^last | part2=recurse^10 | return |
|------|---|------|--------|-----------------|--------------|------------------|--------|
| 1 | [3] | 3 | [] | `f(2,[]) = 1` | 2^3 = 8 | 1^10 = 1 | 8 * 1 = **8** |

Return `8`. ✓

---

## Key Takeaways
- **Never materialize a giant exponent.** Decompose it digit by digit via `a^(10x+d) = (a^x)^10 · a^d`.
- **Modular reduction after every multiply** keeps numbers small and the answer correct (`(x·y) mod m` folds through products).
- **Binary exponentiation** (`square-and-multiply`) computes `base^exp mod m` in O(log exp); here the exponents are the tiny constants 10 and a single digit.
- The modulus 1337 is not prime (`1337 = 7 · 191`), so Fermat/Euler shortcuts don't cleanly apply — the digit decomposition is the robust route.

---

## Related Problems
- LeetCode #50 — Pow(x, n) (fast exponentiation core)
- LeetCode #69 — Sqrt(x) (numeric divide-and-conquer)
- LeetCode #372 relatives: any modular-arithmetic / big-exponent problem
