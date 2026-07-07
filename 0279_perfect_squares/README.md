# 0279 — Perfect Squares

> LeetCode #279 · Difficulty: Medium
> **Categories:** Dynamic Programming, BFS, Math

---

## Problem Statement

Given an integer `n`, return _the least number of perfect square numbers that
sum to_ `n`.

A **perfect square** is an integer that is the square of an integer; in other
words, it is the product of some integer with itself. For example, `1`, `4`, `9`,
and `16` are perfect squares while `3` and `11` are not.

**Example 1:**

```
Input: n = 12
Output: 3
Explanation: 12 = 4 + 4 + 4.
```

**Example 2:**

```
Input: n = 13
Output: 2
Explanation: 13 = 4 + 9.
```

**Constraints:**

- `1 <= n <= 10^4`

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Amazon    | ★★★☆☆ Medium     | 2023          |
| Google    | ★★★☆☆ Medium     | 2023          |
| Microsoft | ★★★☆☆ Medium     | 2022          |
| Adobe     | ★★☆☆☆ Low        | 2022          |
| Bloomberg | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **1D Dynamic Programming** — `dp[i]` = fewest squares summing to `i`, a
  coin-change recurrence → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)
- **BFS (shortest path)** — model integers as nodes and squares as unit-cost
  edges; the answer is the shortest path `0 → n` → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Number theory** — Lagrange's four-square and Legendre's three-square
  theorems bound the answer to `{1,2,3,4}` → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | DP Bottom-Up | O(n√n) | O(n) | Clear coin-change framing; always works |
| 2 | BFS | O(n√n) | O(n) | Shortest-path view; early exit on first reach |
| 3 | Lagrange's Four-Square (Optimal) | O(√n) | O(1) | Fastest; uses number-theory classification |

---

## Approach 1 — DP Bottom-Up

### Intuition
Let `dp[i]` be the minimum number of perfect squares summing to `i`. To build
`i`, pick some square `j*j ≤ i` as the last term; the rest, `i - j*j`, is formed
optimally. So `dp[i] = min over j of dp[i - j*j] + 1`, with `dp[0] = 0`. This is
exactly the coin-change "fewest coins" recurrence with square-valued coins.

### Algorithm
1. `dp[0] = 0`; every other `dp[i]` starts at +∞.
2. For `i = 1..n`, for each `j` with `j*j ≤ i`: `dp[i] = min(dp[i], dp[i-j*j]+1)`.
3. Return `dp[n]`.

### Complexity
- **Time:** O(n√n) — each `i` tries up to `√i` squares.
- **Space:** O(n) — the `dp` table.

### Code
```go
func dpBottomUp(n int) int {
	dp := make([]int, n+1) // dp[i] = fewest squares that sum to i
	for i := 1; i <= n; i++ {
		dp[i] = math.MaxInt32 // start each as "unreachable / infinity"
		for j := 1; j*j <= i; j++ {
			// take square j*j as the last term; add 1 to the best for the remainder
			if dp[i-j*j]+1 < dp[i] {
				dp[i] = dp[i-j*j] + 1
			}
		}
	}
	return dp[n]
}
```

### Dry Run
`n = 12` (showing the relevant cells):

| i | squares tried (j*j) | dp[i] |
|---|---------------------|-------|
| 1 | 1 | dp[0]+1 = 1 |
| 4 | 1,4 | min(dp[3]+1, dp[0]+1) = 1 |
| 8 | 1,4 | min(…, dp[4]+1) = 2 |
| 12 | 1,4,9 | min(dp[11]+1, dp[8]+1, dp[3]+1) = dp[8]+1 = **3** |

`dp[12]` = **3** (12 = 4 + 4 + 4).

---

## Approach 2 — BFS

### Intuition
Treat each integer `0..n` as a graph node. From value `v` you may jump to
`v + j*j` for any square `j*j`, each jump costing one square. BFS from `0`
explores nodes in increasing order of jumps used, so the first time it reaches
`n`, the number of jumps is the minimum.

### Algorithm
1. Precompute all squares `≤ n`.
2. BFS level-by-level from `0`; the level number = squares used so far.
3. From each frontier value, jump by every square; the first path to reach `n`
   returns its level.

### Complexity
- **Time:** O(n√n) — each of ≤ `n` nodes expands ≤ `√n` edges.
- **Space:** O(n) — visited set + queue.

### Code
```go
func bfs(n int) int {
	if n == 0 {
		return 0
	}
	// build the list of usable square values
	squares := []int{}
	for j := 1; j*j <= n; j++ {
		squares = append(squares, j*j)
	}
	visited := make([]bool, n+1) // avoid revisiting a value
	queue := []int{0}            // start BFS at sum 0
	visited[0] = true
	level := 0 // number of squares used to reach the current frontier
	for len(queue) > 0 {
		level++                 // one more square added at this BFS layer
		next := []int{}         // frontier for the next layer
		for _, v := range queue {
			for _, s := range squares {
				nv := v + s
				if nv == n {
					return level // reached target with `level` squares
				}
				if nv < n && !visited[nv] {
					visited[nv] = true      // mark to prevent duplicates
					next = append(next, nv) // explore later
				}
				if nv > n {
					break // squares are ascending; further ones overshoot too
				}
			}
		}
		queue = next
	}
	return level
}
```

### Dry Run
`n = 12`, squares `[1,4,9]`:

| level | frontier | jumps reaching new nodes | reached 12? |
|-------|----------|--------------------------|-------------|
| 1 | {0} | 1, 4, 9 | no |
| 2 | {1,4,9} | 2,5,10 / 5,8,13✗ / 10,13✗,18✗ | no (8,10,… added) |
| 3 | {2,5,8,10,…} | 8+4=12 → hit! | **yes** |

First reach of 12 is at level **3**.

---

## Approach 3 — Lagrange's Four-Square Theorem (Optimal)

### Intuition
Lagrange's theorem guarantees every natural number is a sum of at most four
squares, so the answer is `1`, `2`, `3`, or `4`. Classify directly:
- **1** if `n` is itself a perfect square.
- **4** iff `n = 4^a · (8b + 7)` (Legendre's three-square theorem: exactly the
  numbers not expressible with three squares).
- **2** if `n = a² + b²` for some `a` (scan `a` with `a² ≤ n`).
- **3** otherwise.

### Algorithm
1. If `n` is a perfect square → `1`.
2. Strip factors of `4` from `n`; if the remainder ≡ `7 (mod 8)` → `4`.
3. For `a` with `a² ≤ n`: if `n - a²` is a perfect square → `2`.
4. Otherwise → `3`.

### Complexity
- **Time:** O(√n) — the two-square scan dominates.
- **Space:** O(1).

### Code
```go
func mathFourSquare(n int) int {
	isSquare := func(x int) bool {
		r := int(math.Sqrt(float64(x))) // candidate integer root
		return r*r == x                 // exact when r*r reproduces x
	}
	if isSquare(n) {
		return 1 // n is itself a perfect square
	}
	// Legendre's three-square theorem: answer is 4 iff n = 4^a*(8b+7).
	m := n
	for m%4 == 0 { // strip out factors of 4
		m /= 4
	}
	if m%8 == 7 {
		return 4
	}
	// Try to write n as a sum of two squares.
	for a := 1; a*a <= n; a++ {
		if isSquare(n - a*a) {
			return 2
		}
	}
	return 3 // not 1, 2, or 4 → must be 3 by Lagrange
}
```

### Dry Run
`n = 12`:

| Step | Check | Result |
|------|-------|--------|
| 1 | `isSquare(12)`? √12≈3, 3²=9≠12 | not 1 |
| 2 | strip 4s: 12/4=3; `3 % 8 = 3` ≠ 7 | not 4 |
| 3 | a=1: 12-1=11 not square; a=2: 12-4=8 not square; a=3: 12-9=3 not square (a²≤12 exhausted) | not 2 |
| 4 | fall through | **3** |

Return **3**.

---

## Key Takeaways

- Perfect Squares is coin-change in disguise: the "coins" are `1, 4, 9, 16, …`
  and we want the fewest coins summing to `n`.
- "Fewest steps / min count" over unit-cost transitions is naturally BFS — the
  first time the target node is dequeued gives the minimum.
- Number theory can crush a DP: Lagrange (≤4 squares) + Legendre (the `4^a(8b+7)`
  form) reduce the whole problem to an O(√n) classification.

---

## Related Problems

- LeetCode #322 — Coin Change (identical min-coins DP)
- LeetCode #1137 — N-th Tribonacci Number (1D DP recurrence)
- LeetCode #746 — Min Cost Climbing Stairs (min-cost DP)
- LeetCode #127 — Word Ladder (BFS shortest transformation)
