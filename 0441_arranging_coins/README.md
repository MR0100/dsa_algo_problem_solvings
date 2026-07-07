# 0441 — Arranging Coins

> LeetCode #441 · Difficulty: Easy
> **Categories:** Math, Binary Search

---

## Problem Statement

You have `n` coins and you want to build a staircase with these coins. The staircase consists of `k` rows where the `ith` row has exactly `i` coins. The last row of the staircase **may be** incomplete.

Given the integer `n`, return *the number of **complete rows** of the staircase you will build*.

**Example 1:**

```
Input: n = 5
Output: 2
Explanation: Because the 3rd row is incomplete, we return 2.
```

(The staircase looks like:
```
¤
¤ ¤
¤ ¤     ← incomplete: needs 3, only 2 left
```
Rows 1 and 2 are complete = 2 coins + ... = uses 1+2 = 3 coins, leaving 2 coins for row 3 which needs 3 → incomplete.)

**Example 2:**

```
Input: n = 8
Output: 3
Explanation: Because the 4th row is incomplete, we return 3.
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
| Apple      | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★☆☆☆☆ Rare       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Triangular Numbers / Quadratic Equation** — `k` complete rows cost `T(k) = k(k+1)/2` coins; solving `k(k+1)/2 ≤ n` closes the problem in O(1) via the quadratic formula → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)
- **Binary Search on a Monotone Predicate** — `T(k)` increases with `k`, so "does `k` rows fit in `n` coins?" is a monotone true→false test; binary-search the largest `k` that fits → see [`/dsa/binary_search.md`](/dsa/binary_search.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Simulation) | O(√n) | O(1) | Clearest correctness; fine since k ≈ √(2n) ≤ ~65 535 |
| 2 | Binary Search | O(log n) | O(1) | The canonical "search a monotone predicate" answer |
| 3 | Quadratic Formula (Optimal) | O(1) | O(1) | One-liner; watch float precision, but exact within 2⁵³ here |

---

## Approach 1 — Brute Force (Simulation)

### Intuition

Just build the staircase. Row `k` costs `k` coins, so peel coins off the pile one row at a time. The last row you can fully pay for is the number of complete rows. Because the coins spent grow as `1 + 2 + … + k`, you exhaust the pile after only about `√(2n)` rows — so even the "naive" loop is cheap.

### Algorithm

1. Initialise `row = 0`, `remaining = n`.
2. While `remaining >= row + 1` (you can afford the next row): increment `row`, then subtract its cost with `remaining -= row`.
3. When the next row is unaffordable, stop and return `row`.

### Complexity

- **Time:** O(√n) — the number of complete rows `k` satisfies `k(k+1)/2 ≤ n`, so `k ≈ √(2n)`; the loop runs that many iterations (≤ ~65 535 for max `n`).
- **Space:** O(1) — two integer counters.

### Code

```go
func bruteForce(n int) int {
	row := 0       // number of complete rows built so far
	remaining := n // coins left in the pile
	// Try to build row (row+1); its cost equals its index.
	for remaining >= row+1 {
		row++             // this row is complete
		remaining -= row  // pay for the row we just completed
	}
	return row // last fully-built row count
}
```

### Dry Run

Example 1: `n = 5`.

| Step | row (before) | need = row+1 | remaining (before) | remaining >= need? | Action | row (after) | remaining (after) |
|------|--------------|--------------|--------------------|--------------------|--------|-------------|-------------------|
| 1 | 0 | 1 | 5 | yes | build row 1 | 1 | 4 |
| 2 | 1 | 2 | 4 | yes | build row 2 | 2 | 2 |
| 3 | 2 | 3 | 2 | no (2 < 3) | stop | 2 | 2 |

Return `row = 2` ✔ — row 3 needs 3 coins but only 2 remain.

---

## Approach 2 — Binary Search

### Intuition

The total coins to complete `k` rows is the triangular number `T(k) = k(k+1)/2`, which strictly increases with `k`. We want the **largest** `k` with `T(k) ≤ n`. "Does `k` rows fit?" is a monotone predicate (true for small `k`, false once `k` is too big), so binary search finds the boundary in `O(log n)`.

### Algorithm

1. Set `lo = 1`, `hi = n` (since `T(k)` is quadratic, `k` can never exceed `n`).
2. While `lo <= hi`: compute `mid = lo + (hi-lo)/2` and `curr = mid*(mid+1)/2` (in 64-bit to avoid overflow).
3. If `curr == n`, return `mid` (exact fit). If `curr < n`, `mid` rows fit — search right (`lo = mid+1`). Else search left (`hi = mid-1`).
4. On exit, `hi` is the largest `k` with `T(k) ≤ n`; return `hi`.

### Complexity

- **Time:** O(log n) — the search interval halves each iteration.
- **Space:** O(1) — a handful of 64-bit scalars.

### Code

```go
func binarySearch(n int) int {
	lo, hi := 1, n // k lies in [1, n]
	for lo <= hi {
		mid := lo + (hi-lo)/2 // candidate row count (avoid lo+hi overflow)
		// Triangular number T(mid); use int64 to dodge overflow near n=2^31-1.
		curr := int64(mid) * int64(mid+1) / 2
		target := int64(n)
		switch {
		case curr == target:
			return mid // exact fit: mid complete rows, none left over
		case curr < target:
			lo = mid + 1 // mid rows fit with coins to spare — try more
		default:
			hi = mid - 1 // mid rows cost too much — try fewer
		}
	}
	// Loop exits with hi = largest k where T(k) < n < T(k+1); hi is the answer.
	return hi
}
```

### Dry Run

Example 1: `n = 5`.

| Step | lo | hi | mid | curr = mid(mid+1)/2 | Compare to 5 | Action |
|------|----|----|-----|---------------------|--------------|--------|
| 1 | 1 | 5 | 3 | 6 | 6 > 5 | hi = 2 |
| 2 | 1 | 2 | 1 | 1 | 1 < 5 | lo = 2 |
| 3 | 2 | 2 | 2 | 3 | 3 < 5 | lo = 3 |
| 4 | 3 | 2 | — | — | lo > hi → exit | return hi = 2 |

Return `2` ✔.

---

## Approach 3 — Quadratic Formula (Optimal)

### Intuition

Solve the inequality directly. We need the largest `k` with `k(k+1)/2 ≤ n`, i.e. `k² + k − 2n ≤ 0`. The positive root of `k² + k − 2n = 0` is `k = (−1 + √(1 + 8n)) / 2`. The floor of that root is exactly the number of complete rows.

### Algorithm

1. Compute `k = (−1 + √(1 + 8n)) / 2` in floating point.
2. Return `⌊k⌋` (integer truncation).

### Complexity

- **Time:** O(1) — a single square root.
- **Space:** O(1).

### Code

```go
func mathFormula(n int) int {
	// Positive root of k² + k − 2n = 0, then floored.
	return int((-1 + math.Sqrt(1+8*float64(n))) / 2)
}
```

> Precision note: `8n` peaks around `1.7×10¹⁰`, far below float64's exact-integer limit `2⁵³ ≈ 9×10¹⁵`, so `1 + 8n` is represented exactly and the `sqrt`/floor land on the right integer for every valid `n`.

### Dry Run

Example 1: `n = 5`.

| Step | Expression | Value |
|------|------------|-------|
| 1 | `8n` | 40 |
| 2 | `1 + 8n` | 41 |
| 3 | `√41` | ≈ 6.403 |
| 4 | `−1 + √41` | ≈ 5.403 |
| 5 | `(…)/2` | ≈ 2.701 |
| 6 | `⌊2.701⌋` | 2 |

Return `2` ✔.

---

## Key Takeaways

- **Recognise triangular numbers.** "Row `i` has `i` items" ⇒ total after `k` rows is `k(k+1)/2`. Any problem phrased as a growing staircase / cumulative sum of `1..k` reduces to this closed form.
- **Monotone predicate ⇒ binary search.** "Largest `k` such that `f(k) ≤ n`" with `f` increasing is the textbook last-true search; set `hi` to a safe over-estimate (`n` here).
- **Guard against overflow, not the search.** With `n` up to `2^31 − 1`, `mid*(mid+1)` overflows 32-bit and even signed 32-bit `int`; compute the triangular number in `int64`.
- **Closed form beats iteration when it exists** — but sanity-check floating-point precision against the `2⁵³` exact-integer boundary before trusting `sqrt`/`floor`.

---

## Related Problems

- LeetCode #69 — Sqrt(x) (integer square root via binary search / Newton)
- LeetCode #367 — Valid Perfect Square (same monotone-search / formula toolkit)
- LeetCode #278 — First Bad Version (last-true / first-false boundary search)
- LeetCode #1250 — Check If It Is a Good Array (number-theory closed form)
