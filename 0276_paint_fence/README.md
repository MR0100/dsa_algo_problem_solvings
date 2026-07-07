# 0276 — Paint Fence

> LeetCode #276 · Difficulty: Medium
> **Categories:** Dynamic Programming

---

## Problem Statement

You are painting a fence of `n` posts with `k` different colors. You must paint
all the posts such that **no more than two adjacent** fence posts have the same
color.

Return _the number of ways_ you can paint the fence.

**Example 1:**

```
Input: n = 3, k = 2
Output: 6
Explanation: All the possibilities are shown.
Note that painting all posts red or all posts green is invalid because there
cannot be three posts in a row with the same color.
   post1  post2  post3
1   R      R      G
2   R      G      R
3   R      G      G
4   G      R      R
5   G      R      G
6   G      G      R
```

**Example 2:**

```
Input: n = 1, k = 1
Output: 1
```

**Example 3:**

```
Input: n = 7, k = 2
Output: 42
```

**Constraints:**

- `1 <= n <= 50`
- `1 <= k <= 10^5`
- The testcases are generated such that the answer is in the range `[0, 2^31 - 1]`
  for the given `n` and `k`.

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Google    | ★★★☆☆ Medium     | 2023          |
| Amazon    | ★★★☆☆ Medium     | 2023          |
| Microsoft | ★★☆☆☆ Low        | 2022          |
| Adobe     | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **1D Dynamic Programming** — the answer for `i` posts is a linear recurrence on
  the answers for `i-1` and `i-2` posts → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Recursion | O(2^n) | O(n) | Understand the recurrence; too slow to submit |
| 2 | DP Bottom-Up (table) | O(n) | O(n) | Clear, memoized version of the recurrence |
| 3 | DP O(1) Space (Optimal) | O(n) | O(1) | Only the last two states are needed |

---

## Approach 1 — Recursion

### Intuition
When we paint post `i`, we either **reuse** the color of post `i-1` or pick a
**new** color different from it. Reusing is legal only if posts `i-1` and `i-2`
already differ (else three in a row match). Counting:
- "different from previous": `(k-1)` new-color choices applied to `total(i-1)`.
- "same as previous": legal only when the pair before differed, which is exactly
  `total(i-2)`, again scaled by the `(k-1)` choices that created that difference.

Both branches share the `(k-1)` factor, giving
`total(i) = (k-1) * (total(i-1) + total(i-2))`.

### Algorithm
1. If `n == 0` return `0`.
2. If `n == 1` return `k`.
3. If `n == 2` return `k*k`.
4. Otherwise return `(k-1) * (recursion(n-1) + recursion(n-2))`.

### Complexity
- **Time:** O(2^n) — each call spawns two more with no memoization.
- **Space:** O(n) — maximum recursion depth.

### Code
```go
func recursion(n, k int) int {
	if n == 0 { // no posts → no way to paint
		return 0
	}
	if n == 1 { // a single post can be any of the k colors
		return k
	}
	if n == 2 { // two posts: k choices for first, k for second (may match)
		return k * k
	}
	// (k-1) new-color choices multiply BOTH the "reuse" (n-2) and "new" (n-1) cases
	return (k - 1) * (recursion(n-1, k) + recursion(n-2, k))
}
```

### Dry Run
Trace `recursion(3, 2)`:

| Call | n | Returns |
|------|---|---------|
| `recursion(3,2)` | 3 | `(2-1)*(recursion(2)+recursion(1))` |
| `recursion(2,2)` | 2 | `k*k = 4` |
| `recursion(1,2)` | 1 | `k = 2` |
| back to `recursion(3,2)` | 3 | `1*(4+2) = 6` |

Result: **6**.

---

## Approach 2 — DP Bottom-Up (Table)

### Intuition
The recursion recomputes overlapping subproblems. Store each `total(i)` once in
an array and iterate upward from the base cases.

### Algorithm
1. Handle `n == 0` (→0) and `n == 1` (→k) directly.
2. Set `dp[1] = k`, `dp[2] = k*k`.
3. For `i = 3..n`: `dp[i] = (k-1) * (dp[i-1] + dp[i-2])`.
4. Return `dp[n]`.

### Complexity
- **Time:** O(n) — a single pass filling the table.
- **Space:** O(n) — the `dp` array.

### Code
```go
func dpBottomUp(n, k int) int {
	if n == 0 {
		return 0
	}
	if n == 1 {
		return k
	}
	dp := make([]int, n+1) // dp[i] = ways to paint the first i posts
	dp[1] = k              // 1 post → k ways
	dp[2] = k * k          // 2 posts → k*k ways
	for i := 3; i <= n; i++ {
		dp[i] = (k - 1) * (dp[i-1] + dp[i-2]) // fold the recurrence
	}
	return dp[n]
}
```

### Dry Run
Trace `dpBottomUp(3, 2)`:

| i | dp[i] | Computation |
|---|-------|-------------|
| 1 | 2 | base: `k` |
| 2 | 4 | base: `k*k` |
| 3 | 6 | `(2-1)*(dp[2]+dp[1]) = 1*(4+2)` |

Return `dp[3]` = **6**.

---

## Approach 3 — DP O(1) Space (Optimal)

### Intuition
`total(i)` depends only on the previous two values, so two rolling scalars
replace the whole table.

### Algorithm
1. `prev2 = k` (value at i=1), `prev1 = k*k` (value at i=2).
2. For `i = 3..n`: `cur = (k-1)*(prev1+prev2)`; then slide `prev2 = prev1`,
   `prev1 = cur`.
3. Return `prev1`.

### Complexity
- **Time:** O(n) — single loop.
- **Space:** O(1) — two scalars.

### Code
```go
func dpConstantSpace(n, k int) int {
	if n == 0 {
		return 0
	}
	if n == 1 {
		return k
	}
	prev2 := k     // ways for i-2 (starts at i=1)
	prev1 := k * k // ways for i-1 (starts at i=2)
	for i := 3; i <= n; i++ {
		cur := (k - 1) * (prev1 + prev2) // recurrence
		prev2 = prev1                    // slide the window forward
		prev1 = cur
	}
	return prev1 // holds ways for i=n
}
```

### Dry Run
Trace `dpConstantSpace(3, 2)`:

| i | prev2 | prev1 | cur | Note |
|---|-------|-------|-----|------|
| start | 2 | 4 | — | prev2=k, prev1=k*k |
| 3 | 4 | 6 | 6 | cur=1*(4+2)=6; slide |

Return `prev1` = **6**.

---

## Key Takeaways

- The classic "no more than two adjacent same" constraint collapses to the
  recurrence `f(i) = (k-1) * (f(i-1) + f(i-2))` — a Fibonacci-shaped relation.
- Splitting counts by "reuse previous color" vs "pick a new color" is a reusable
  DP framing for coloring/arrangement problems.
- Any linear recurrence that looks back a fixed number of steps can drop from
  O(n) space to O(1) with rolling variables.

---

## Related Problems

- LeetCode #70 — Climbing Stairs (same two-term recurrence shape)
- LeetCode #198 — House Robber (adjacency-constrained DP)
- LeetCode #256 — Paint House (coloring DP with per-color states)
- LeetCode #265 — Paint House II (coloring DP, k colors)
