# 0479 — Largest Palindrome Product

> LeetCode #479 · Difficulty: Hard
> **Categories:** Math, Enumeration, Number Theory

---

## Problem Statement

Given an integer `n`, return *the **largest palindromic integer** that can be represented as the product of two `n`-digits integers*. Since the answer can be very large, return it **modulo** `1337`.

**Example 1:**

```
Input: n = 2
Output: 987
Explanation: 99 x 91 = 9009, 9009 % 1337 = 987
```

**Example 2:**

```
Input: n = 1
Output: 9
```

**Constraints:**

- `1 <= n <= 8`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★☆☆☆ Low        | 2023          |
| Yahoo      | ★★☆☆☆ Low        | 2022          |
| Google     | ★☆☆☆☆ Rare       | 2022          |
| Microsoft  | ★☆☆☆☆ Rare       | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Math / Number Theory** — the answer for `n > 1` is a `2n`-digit palindrome fully determined by its first half; solving reduces to enumerating halves and testing factorization by trial division up to `√P` → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)
- **Enumeration with Pruning** — both approaches generate candidates in *descending* order and stop at the first success, so the largest valid candidate is found immediately → see [`/dsa/greedy.md`](/dsa/greedy.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Descending Product Search) | O((10ⁿ)²) | O(n) | Clear and correct up to `n ≤ 4`; TLE beyond |
| 2 | Construct Palindrome, Test Factorization (Optimal) | ~O(10ⁿ · short) | O(1) | The intended answer; handles all `n ≤ 8` fast |

---

## Approach 1 — Brute Force (Descending Product Search)

### Intuition

`n`-digit numbers lie in `[10^(n-1), 10^n − 1]`. To find the **largest** palindrome product, examine large products first. Loop `i` from the top down and `j` from `i` down (each unordered pair once); keep the biggest product that happens to be a palindrome. Two prunes make it survivable for small `n`: if `i·i ≤ best`, no larger product can appear, so stop entirely; within a row, once `i·j ≤ best`, break to the next `i`.

### Algorithm

1. Special-case `n == 1` → return `9`.
2. `hi = 10^n − 1`, `lo = 10^(n-1)`, `best = 0`.
3. For `i` from `hi` down to `lo`:
   - if `i·i ≤ best`, break.
   - For `j` from `i` down to `lo`:
     - if `i·j ≤ best`, break.
     - if `i·j` is a palindrome, set `best = i·j` and break the inner loop.
4. Return `best % 1337`.

### Complexity

- **Time:** O((10ⁿ)²) worst case — the full pair space; acceptable for `n ≤ 4` thanks to pruning, hopeless at `n = 8`.
- **Space:** O(n) for the string used in the palindrome check.

### Code

```go
func bruteForce(n int) int {
	if n == 1 {
		return 9 // 3*3 = 9 is the largest single-digit-product palindrome
	}
	hi := pow10(n) - 1 // largest n-digit number, e.g. n=2 → 99
	lo := pow10(n - 1) // smallest n-digit number, e.g. n=2 → 10
	best := 0
	for i := hi; i >= lo; i-- {
		if i*i <= best { // even i*i can't exceed best → no larger product remains
			break
		}
		for j := i; j >= lo; j-- {
			prod := i * j
			if prod <= best { // products only shrink as j drops → give up on this i
				break
			}
			if isPalindrome(prod) {
				best = prod // record a new champion; inner products now all smaller
				break
			}
		}
	}
	return best % 1337
}
```

### Dry Run

Example 1: `n = 2`, so `hi = 99`, `lo = 10`, `best = 0`. (Only the decisive iterations shown.)

| i | j | prod = i·j | palindrome? | best after |
|---|---|-----------|-------------|------------|
| 99 | 99 | 9801 | no (`9801` vs `1089`) | 0 |
| 99 | 98 | 9702 | no | 0 |
| 99 | 97 | 9603 | no | 0 |
| … | … | …    | …  | 0 |
| 99 | 91 | 9009 | **yes** | 9009 (break inner) |
| 98 | 98 | 9604 | `98·98 = 9604 > 9009`, but not palindrome; scanning down, all `98·j ≤ 9604` soon drop `≤ 9009` → break | 9009 |
| 97 | 97 | 9409 | `> 9009` yet no palindrome above 9009 in this row → break | 9009 |
| 95 | 95 | 9025 | `i·i = 9025 > 9009`… next `i=94`: `94² = 8836 ≤ 9009` → outer break | 9009 |

`best = 9009`; return `9009 % 1337 = 987`. ✔

---

## Approach 2 — Construct Palindrome, Test Factorization (Optimal)

### Intuition

For `n > 1`, the largest palindrome product has exactly `2n` digits and is even-length, so it is determined entirely by its **first half**. Enumerate halves from the largest (`999…`) downward, mirror each into a full palindrome `P = half · 10ⁿ + reverse(half)`, and ask: does `P` factor into two `n`-digit numbers? Because we scan `P` from largest down, the *first* `P` that factors is the answer. To test factorization, trial-divide `P` by candidate divisors `d` from `hi` downward; stop once `d·d < P`, since beyond `√P` the cofactor `P/d` would exceed `hi` (all valid pairs are already checked by symmetry).

### Algorithm

1. Special-case `n == 1` → return `9`.
2. `hi = 10^n − 1`, `lo = 10^(n-1)`.
3. For `half` from `hi` down to `lo`:
   - Build `P` = `half` followed by its reversed digits (a `2n`-digit palindrome).
   - For `d` from `hi` down while `d·d ≥ P`:
     - if `P % d == 0` and `P/d ≤ hi`, return `P % 1337`.
4. (Unreachable for valid `n`.)

### Complexity

- **Time:** roughly O(10ⁿ) candidate palindromes, each with a short divisor scan; in practice the first (or very early) palindrome factors, so it is effectively fast for all `n ≤ 8`. Products reach ~`10^16`, handled in 64-bit.
- **Space:** O(1).

### Code

```go
func buildAndFactor(n int) int {
	if n == 1 {
		return 9
	}
	hi := pow10(n) - 1 // largest n-digit factor
	lo := pow10(n - 1) // smallest n-digit factor
	for half := hi; half >= lo; half-- {
		p := makePalindrome(half) // mirror half → full 2n-digit palindrome (int64-safe)
		// Trial divide P by n-digit numbers from the top; stop past sqrt(P).
		for d := hi; d*d >= p; d-- {
			if p%d == 0 { // d divides P
				cofactor := p / d
				if cofactor <= hi { // cofactor is n-digit (it is ≥ lo automatically here)
					return int(p % 1337)
				}
			}
		}
	}
	return -1 // not reachable for 1 <= n <= 8
}

func makePalindrome(half int) int {
	pal := half
	rev := half
	for rev > 0 {
		pal = pal*10 + rev%10 // shift pal left and append reversed digit
		rev /= 10
	}
	return pal
}
```

### Dry Run

Example 1: `n = 2`, `hi = 99`, `lo = 10`.

| half | P = makePalindrome(half) | divisor scan (d from 99, while d²≥P) | result |
|------|--------------------------|--------------------------------------|--------|
| 99 | `9009` | d=99: 9009%99≠0; d=98…; **d=91: 9009%91=0**, cofactor `9009/91 = 99 ≤ 99` ✔ | return `9009 % 1337` |

`9009 % 1337 = 987`. ✔ — the very first (largest) half `99` already yields a factorable palindrome.

---

## Key Takeaways

- **Generate answers in the order you want them.** Enumerating candidates from largest to smallest turns "find the maximum" into "return the first hit" — no full search needed.
- **A palindrome is defined by half its digits.** For an even-length palindrome, pick the top half and mirror it; this shrinks the search space from `10^{2n}` products to `~10^n` halves.
- **Trial division only needs to reach `√P`.** Factor pairs are symmetric, so scanning divisors from `hi` down to `√P` covers every `n`-digit factorization.
- Watch the width: `2n`-digit products for `n = 8` are ~`10^16`, so use 64-bit integers; apply `% 1337` only at the end.

---

## Related Problems

- LeetCode #9 — Palindrome Number (palindrome test / half-reversal)
- LeetCode #866 — Prime Palindrome (construct palindromes, test a property)
- LeetCode #906 — Super Palindromes (enumerate palindromes by their half)
- LeetCode #564 — Find the Closest Palindrome (palindrome construction from a half)
- LeetCode #2417 — Closest Fair Integer (digit-construction enumeration)
