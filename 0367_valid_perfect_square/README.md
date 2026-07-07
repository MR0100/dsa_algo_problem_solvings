# 0367 — Valid Perfect Square

> LeetCode #367 · Difficulty: Easy
> **Categories:** Math, Binary Search

---

## Problem Statement

Given a positive integer `num`, return `true` if `num` is a perfect square or `false` otherwise.

A **perfect square** is an integer that is the square of an integer. In other words, it is the product of some integer with itself.

You must not use any built-in library function, such as `sqrt`.

**Example 1:**

```
Input: num = 16
Output: true
Explanation: We return true because 4 * 4 = 16 and 4 is an integer.
```

**Example 2:**

```
Input: num = 14
Output: false
Explanation: We return false because 3.742 * 3.742 = 14 and 3.742 is not an integer.
```

**Constraints:**

- `1 <= num <= 2^31 - 1`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Google     | ★★☆☆☆ Low        | 2023          |
| Apple      | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Binary Search on the answer** — the candidate roots `1..num` are sorted and `mid*mid` is monotonically increasing, so we can binary-search for an exact root → see [`/dsa/binary_search.md`](/dsa/binary_search.md)
- **Number theory / Newton's method** — an integer square-root computation using the fixed-point iteration `x ← (x + num/x)/2` → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Linear Scan (Brute Force) | O(√num) | O(1) | Trivial baseline; fine for small num |
| 2 | Binary Search (Optimal) | O(log num) | O(1) | The standard interview answer |
| 3 | Newton's Method (Optimal) | O(log num) | O(1) | Fewest iterations; the "clever" answer |

---

## Approach 1 — Linear Scan (Brute Force)

### Intuition

A perfect square is `num == i*i` for some integer `i`. Try `i = 1, 2, 3, …`. Because squares increase monotonically, the first time `i*i` reaches or passes `num` we can stop — nothing larger can equal `num` either.

### Algorithm

1. For `i = 1, 2, 3, …` while `i*i <= num`:
   - If `i*i == num`, return `true`.
2. If the loop ends, return `false`.

### Complexity

- **Time:** O(√num) — the loop runs until `i` reaches `√num`.
- **Space:** O(1).

### Code

```go
func linearScan(num int) bool {
	for i := 1; i*i <= num; i++ { // stop once the square meets/exceeds num
		if i*i == num { // exact hit ⇒ perfect square
			return true
		}
	}
	return false // never landed exactly on num
}
```

### Dry Run

Example 1: `num = 16`.

| i | i*i | i*i <= 16? | i*i == 16? |
|---|-----|-----------|-----------|
| 1 | 1 | yes | no |
| 2 | 4 | yes | no |
| 3 | 9 | yes | no |
| 4 | 16 | yes | **yes → return true** |

Result: `true` ✔

---

## Approach 2 — Binary Search (Optimal)

### Intuition

The roots `1..num` are sorted and `mid*mid` is monotonically increasing, so binary search applies. Pick the midpoint, compare its square to `num`, and discard half the range each step: too small → search right, too big → search left, equal → found.

### Algorithm

1. `lo = 1`, `hi = num`.
2. While `lo <= hi`:
   - `mid = lo + (hi-lo)/2`, `sq = mid*mid`.
   - `sq == num` → return `true`.
   - `sq < num` → `lo = mid + 1`.
   - `sq > num` → `hi = mid - 1`.
3. Return `false`.

### Complexity

- **Time:** O(log num) — the interval halves each iteration.
- **Space:** O(1).

### Code

```go
func binarySearch(num int) bool {
	lo, hi := 1, num
	for lo <= hi {
		mid := lo + (hi-lo)/2 // overflow-safe midpoint
		sq := mid * mid       // candidate square
		switch {
		case sq == num:
			return true // found the exact root
		case sq < num:
			lo = mid + 1 // too small; go right
		default:
			hi = mid - 1 // too big; go left
		}
	}
	return false
}
```

### Dry Run

Example 1: `num = 16`.

| Step | lo | hi | mid | sq = mid² | Compare to 16 | Action |
|------|----|----|-----|-----------|----------------|--------|
| 1 | 1 | 16 | 8 | 64 | 64 > 16 | hi = 7 |
| 2 | 1 | 7 | 4 | 16 | 16 == 16 | **return true** |

Result: `true` ✔

---

## Approach 3 — Newton's Method (Optimal)

### Intuition

To solve `x² = num`, apply Newton's method to `f(x) = x² − num`, whose update is `x ← (x + num/x)/2`. Starting from `x = num`, the sequence converges quadratically to `⌊√num⌋`. Because integer division floors, iterating until `x` stops decreasing (i.e. `x*x <= num`) lands on the floor of the true root; then `x*x == num` decides the answer.

### Algorithm

1. `x = num`.
2. While `x*x > num`: `x = (x + num/x) / 2`.
3. Return `x*x == num`.

### Complexity

- **Time:** O(log num) — quadratic convergence needs very few iterations.
- **Space:** O(1).

### Code

```go
func newtonsMethod(num int) bool {
	x := num                 // initial guess (an over-estimate)
	for x*x > num {          // shrink until x is the floor of the root
		x = (x + num/x) / 2 // Newton step toward √num
	}
	return x*x == num // exact only if num is a perfect square
}
```

### Dry Run

Example 1: `num = 16`.

| Step | x before | x*x > 16? | x = (x + 16/x)/2 |
|------|----------|-----------|-------------------|
| 1 | 16 | 256 > 16 yes | (16 + 1)/2 = 8 |
| 2 | 8 | 64 > 16 yes | (8 + 2)/2 = 5 |
| 3 | 5 | 25 > 16 yes | (5 + 3)/2 = 4 |
| 4 | 4 | 16 > 16 no → exit | — |

Final check: `4*4 == 16` → `true` ✔

---

## Key Takeaways

- **Monotone predicate ⇒ binary search on the answer.** `mid*mid` increasing is all you need to binary-search the root without computing `sqrt`.
- **Newton's method for integer roots:** `x ← (x + num/x)/2` from `x = num` converges to `⌊√num⌋`; verify with a final `x*x == num`.
- Use `lo + (hi-lo)/2` to avoid `lo + hi` overflow when `num` approaches `2^31 − 1`.
- The same "guess a root, square it, compare" template solves LeetCode #69 (Sqrt(x)).

---

## Related Problems

- LeetCode #69 — Sqrt(x) (integer square root, same three techniques)
- LeetCode #633 — Sum of Square Numbers (two-pointer over squares)
- LeetCode #279 — Perfect Squares (DP built on perfect-square checks)
- LeetCode #50 — Pow(x, n) (fast exponentiation, another math/search hybrid)
