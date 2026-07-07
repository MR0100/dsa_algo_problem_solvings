# 0343 — Integer Break

> LeetCode #343 · Difficulty: Medium
> **Categories:** Math, Dynamic Programming

---

## Problem Statement

Given an integer `n`, break it into the sum of `k` **positive integers**, where `k >= 2`, and maximize the product of those integers.

Return *the maximum product you can get*.

**Example 1:**

```
Input: n = 2
Output: 1
Explanation: 2 = 1 + 1, 1 × 1 = 1.
```

**Example 2:**

```
Input: n = 10
Output: 36
Explanation: 10 = 3 + 3 + 4, 3 × 3 × 4 = 36.
```

**Constraints:**

- `2 <= n <= 58`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Google     | ★★☆☆☆ Low        | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Apple      | ★★☆☆☆ Low        | 2022          |
| ByteDance  | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Dynamic Programming (1D)** — `dp[i]` = max product of breaking `i`; each value builds on smaller ones over all first cuts → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)
- **Greedy / Number Theory** — the optimal factorization is all 3s (with a 2 or two 2s to mop up the remainder), a provable greedy fact → see [`/dsa/greedy.md`](/dsa/greedy.md)
- **Math** — the "prefer 3s over 2s" insight comes from comparing product-per-unit → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | DP Bottom-Up | O(n²) | O(n) | Safe default; no proof needed, easy to reason about |
| 2 | Math — Break into 3s (Optimal) | O(log n) | O(1) | When you know/can prove the 3s rule; fastest |
| 3 | DP Top-Down (Memoised) | O(n²) | O(n) | Same recurrence, recursive phrasing |

---

## Approach 1 — DP Bottom-Up

### Intuition
To break `i`, choose a first piece `j` (from 1 to i-1). The remainder `i-j` can either be left whole (product `j*(i-j)`) or broken further (product `j*dp[i-j]`). `dp[i]` is the best over all choices of `j`.

### Algorithm
1. `dp[1] = 1` (a leftover 1 acts as a factor).
2. For `i = 2..n`, for `j = 1..i-1`: `dp[i] = max(dp[i], j*(i-j), j*dp[i-j])`.
3. Return `dp[n]`.

### Complexity
- **Time:** O(n²) — nested loops over `i` and `j`.
- **Space:** O(n) — the dp table.

### Code
```go
func dpBottomUp(n int) int {
	dp := make([]int, n+1)
	dp[1] = 1
	for i := 2; i <= n; i++ {
		for j := 1; j < i; j++ {
			best := max(j*(i-j), j*dp[i-j])
			dp[i] = max(dp[i], best)
		}
	}
	return dp[n]
}
```

### Dry Run
Input `n = 2`:

| i | j | `j*(i-j)` | `j*dp[i-j]` | dp[i] |
|---|---|-----------|-------------|-------|
| 2 | 1 | 1*(1)=1 | 1*dp[1]=1*1=1 | max(0,1,1)=**1** |

`dp[2] = 1` → return **1**. (For reference, `n=10` yields `dp[10]=36`.)

---

## Approach 2 — Math — Break into 3s (Optimal)

### Intuition
The factor `3` maximises product-per-unit: e.g. `6 → 3*3 = 9` beats `2*2*2 = 8`. Never keep a factor of 1. So pull out as many 3s as possible. A remainder of 1 is wasteful (a 3+1 → `3*1=3`), so trade one 3 for two 2s (`2*2 = 4 > 3`). A remainder of 2 just multiplies by 2.

### Algorithm
1. If `n <= 3`, return `n-1` (forced to split: 2→1, 3→2).
2. Pull out 3s while `n > 4`: `product *= 3; n -= 3`.
3. Multiply the remaining `n` (which is 2, 3, or 4) into the product. (Remainder 4 → `2*2`, already optimal.)

### Complexity
- **Time:** O(log n) if exponentiation is used, or O(n/3) with the subtraction loop.
- **Space:** O(1).

### Code
```go
func mathThrees(n int) int {
	if n <= 3 {
		return n - 1
	}
	product := 1
	for n > 4 {
		product *= 3
		n -= 3
	}
	return product * n
}
```

### Dry Run
Input `n = 2`: `n <= 3` → return `n-1 = 1` → **1**.

Input `n = 10` (illustrative): `product=1`; loop: n=10>4 → product=3,n=7; n=7>4 → product=9,n=4; n=4 not >4 → stop; return `9*4 = 36`.

---

## Approach 3 — DP Top-Down (Memoised)

### Intuition
Same recurrence as bottom-up, phrased recursively: `solve(i)` tries every first cut `j` and takes `max(j*(i-j), j*solve(i-j))`. Memoise to avoid recomputing overlapping subproblems.

### Algorithm
1. `memo[i]` caches `solve(i)`; `solve(1) = 1`.
2. For `j = 1..i-1`, `best = max(best, j*(i-j), j*solve(i-j))`; store and return.

### Complexity
- **Time:** O(n²) — each of `n` states does O(n) work once.
- **Space:** O(n) — memo table + recursion depth.

### Code
```go
func dpTopDown(n int) int {
	memo := make([]int, n+1)
	var solve func(i int) int
	solve = func(i int) int {
		if i == 1 {
			return 1
		}
		if memo[i] != 0 {
			return memo[i]
		}
		best := 0
		for j := 1; j < i; j++ {
			best = max(best, max(j*(i-j), j*solve(i-j)))
		}
		memo[i] = best
		return best
	}
	return solve(n)
}
```

### Dry Run
Input `n = 2`:

| Call | Loop j | `j*(i-j)` | `j*solve(i-j)` | best |
|------|--------|-----------|----------------|------|
| solve(2) | j=1 | 1*(1)=1 | 1*solve(1)=1*1=1 | 1 |

`memo[2]=1`, return **1**.

---

## Key Takeaways

- Classic DP recurrence: `dp[i] = max over j of max(j*(i-j), j*dp[i-j])` — the two terms mean "stop breaking the remainder" vs "keep breaking it."
- The mathematical optimum is "use as many 3s as possible; if the remainder is 1, convert a 3 into 2+2." This is worth memorising — it appears in several product-maximisation problems.
- Watch the base case: `n = 2` and `n = 3` are forced to split into at least two parts, so they return `n-1`, not `n`.
- The 3s rule is why the DP answer for `n=10` is `3*3*4`, not `3*3*3*1`.

---

## Related Problems

- LeetCode #279 — Perfect Squares (DP over decompositions)
- LeetCode #96 — Unique Binary Search Trees (Catalan-style DP over splits)
- LeetCode #1808 — Maximize Number of Nice Divisors (same "break into 3s" math at scale)
- LeetCode #650 — 2 Keys Keyboard (prime-factor / break-into-factors flavour)
