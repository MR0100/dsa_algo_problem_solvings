# 0400 — Nth Digit

> LeetCode #400 · Difficulty: Medium
> **Categories:** Math, Binary Search (conceptual), Digit Counting

---

## Problem Statement

Given an integer `n`, return the `nth` digit of the infinite integer sequence `[1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, ...]`.

**Example 1:**

```
Input: n = 3
Output: 3
```

**Example 2:**

```
Input: n = 11
Output: 0
Explanation: The 11th digit of the sequence 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, ... is a 0, which is part of the number 10.
```

**Constraints:**

- `1 <= n <= 2^31 - 1`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |
| Apple      | ★☆☆☆☆ Rare       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Math / Counting by Digit Blocks** — numbers group by length: `9` one-digit, `90` two-digit, `900` three-digit … Each length-`L` block contributes `L·9·10^(L-1)` digits; skipping whole blocks pinpoints the target in O(log n) → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(n) | O(1) | Correctness reference for small `n`; TLE for `n` near `2^31` |
| 2 | Math / Digit Blocks (Optimal) | O(log n) | O(1) | Always — jumps over blocks, handles `n = 2^31 − 1` instantly |

---

## Approach 1 — Brute Force

### Intuition

The sequence is just the integers written out and glued together. Walk `1, 2, 3, …`, and for each number subtract its digit-count from `n`. The moment `n` no longer exceeds the current number's length, the target digit lives inside that number.

### Algorithm

1. `num = 1`.
2. Loop: let `s` be the decimal string of `num`. If `n <= len(s)`, return the `(n-1)`-th character of `s` as an integer. Otherwise `n -= len(s)`, `num++`.

### Complexity

- **Time:** O(n) — up to ~`n` numbers are visited in the worst case.
- **Space:** O(1) — one short string reused each iteration.

### Code

```go
func bruteForce(n int) int {
	num := 1
	for {
		s := strconv.Itoa(num) // decimal digits of the current number
		if n <= len(s) {
			// The n-th remaining digit is the (n-1)-th char of this number.
			return int(s[n-1] - '0')
		}
		n -= len(s) // skip all of this number's digits
		num++
	}
}
```

### Dry Run

Example 2: `n = 11`.

| num | s | len | n <= len? | action | n after |
|-----|---|-----|-----------|--------|---------|
| 1 | "1" | 1 | 11≤1? no | n -= 1 | 10 |
| 2 | "2" | 1 | no | n -= 1 | 9 |
| 3..9 | … | 1 each | no | n -= 1 seven times | 2 |
| 10 | "10" | 2 | 2≤2? yes | return s[2-1]=s[1]='0' | — |

After consuming digits of 1..9 (nine digits), `n = 2`; the number `10` supplies its 2nd digit `0`. Result: `0` ✔

---

## Approach 2 — Math / Digit Blocks (Optimal)

### Intuition

Instead of one number at a time, skip whole **blocks of equal length**:

| length L | count of numbers | digits contributed | example range |
|----------|------------------|--------------------|---------------|
| 1 | 9 | 1·9 = 9 | 1..9 |
| 2 | 90 | 2·90 = 180 | 10..99 |
| 3 | 900 | 3·900 = 2700 | 100..999 |
| L | 9·10^(L-1) | L·9·10^(L-1) | 10^(L-1)..10^L−1 |

Subtract each block's total digits from `n` until `n` lands inside a block. Then arithmetic gives the exact number and the exact digit within it — no iteration over individual numbers.

### Algorithm

1. `length = 1`, `count = 9`, `start = 1`.
2. While `n > length*count`: `n -= length*count`; `length++`; `count *= 10`; `start *= 10`.
3. The target number is `number = start + (n-1)/length`.
4. The digit index inside it is `(n-1)%length`; return that digit of `number`.

### Complexity

- **Time:** O(log n) — one iteration per digit-length block (≤ 10 for 32-bit `n`).
- **Space:** O(1) — a handful of integer accumulators.

### Code

```go
func mathBlocks(n int) int {
	length := 1   // current block's number length (digits per number)
	count := 9    // how many numbers have this length: 9,90,900,...
	start := 1    // first number of this length: 1,10,100,...
	for n > length*count {
		n -= length * count // consume this block's digits
		length++            // move to longer numbers
		count *= 10         // 10x as many of them
		start *= 10         // block now starts at the next power of ten
	}
	number := start + (n-1)/length     // which number holds the n-th digit
	digitIndex := (n - 1) % length     // which digit inside that number
	s := strconv.Itoa(number)
	return int(s[digitIndex] - '0')
}
```

### Dry Run

Example 2: `n = 11`.

| Step | length | count | start | n > length*count? | action | n after |
|------|--------|-------|-------|-------------------|--------|---------|
| init | 1 | 9 | 1 | 11 > 1·9=9? yes | n -= 9; length=2; count=90; start=10 | 2 |
| loop | 2 | 90 | 10 | 2 > 2·90=180? no | exit loop | 2 |

Now `n = 2`, `length = 2`, `start = 10`:
- `number = 10 + (2-1)/2 = 10 + 0 = 10`.
- `digitIndex = (2-1)%2 = 1`.
- `strconv.Itoa(10)[1] = '0'` → `0`.

Result: `0` ✔ — matches Example 2's explanation (the 11th digit is the `0` in `10`).

---

## Key Takeaways

- **Group by digit length** to turn an O(n) walk into O(log n): block `L` holds `9·10^(L-1)` numbers and `L·9·10^(L-1)` digits.
- After locating the block, two integer operations finish the job: `(n-1)/length` selects the number offset from `start`, `(n-1)%length` selects the digit within it. Watch the **1-based → 0-based** conversion (`n-1`).
- Beware overflow: for very large `n`, `length*count` can grow large — Go's `int` is 64-bit on modern platforms, so `n = 2^31 − 1` is safe here, but in strict 32-bit environments use `int64` for the block-size products.
- This "skip fixed-size blocks, then index within one" pattern recurs in problems about positions in structured infinite sequences.

---

## Related Problems

- LeetCode #233 — Number of Digit One (digit-position counting)
- LeetCode #172 — Factorial Trailing Zeroes (block/period counting math)
- LeetCode #60 — Permutation Sequence (locate an item by dividing into fixed-size blocks)
- LeetCode #1015 — Smallest Integer Divisible by K (digit-sequence math)
