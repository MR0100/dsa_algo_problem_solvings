# 0474 — Ones and Zeroes

> LeetCode #474 · Difficulty: Medium
> **Categories:** Array, String, Dynamic Programming, 0/1 Knapsack

---

## Problem Statement

You are given an array of binary strings `strs` and two integers `m` and `n`.

Return *the size of the largest subset of `strs` such that there are **at most** `m` `0`'s and `n` `1`'s in the subset*.

A set `x` is a **subset** of a set `y` if all elements of `x` are also elements of `y`.

**Example 1:**

```
Input: strs = ["10","0001","111001","1","0"], m = 5, n = 3
Output: 4
Explanation: The largest subset with at most 5 0's and 3 1's is {"10", "0001", "1", "0"}, so the answer is 4.
Other valid but smaller subsets include {"0001", "1"} and {"10", "1", "0"}.
{"111001"} is an invalid subset because it contains 4 1's, greater than the maximum of 3.
```

**Example 2:**

```
Input: strs = ["10","0","1"], m = 1, n = 1
Output: 2
Explanation: The largest subset is {"0", "1"}, so the answer is 2.
```

**Constraints:**

- `1 <= strs.length <= 600`
- `1 <= strs[i].length <= 100`
- `strs[i]` consists only of digits `'0'` and `'1'`.
- `1 <= m, n <= 100`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Adobe      | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **0/1 Knapsack (two capacities)** — each string is an item taken at most once; its "weight" is the pair (zeros, ones) and its "value" is 1. Maximise count under two independent capacity limits → see [`/dsa/knapsack.md`](/dsa/knapsack.md)
- **Rolling-array DP** — the item dimension collapses to a single 2D table swept downward, the hallmark trick for keeping 0/1 knapsack at `O(capacity)` space → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (subset enum) | O(2^k) | O(k) | Understand the take/skip structure; TLE for k up to 600 |
| 2 | 3D DP (item × m × n) | O(k·m·n) | O(k·m·n) | Explicit, easy to reason about; memory-heavy |
| 3 | 2D Rolling DP (Optimal) | O(k·m·n) | O(m·n) | The accepted answer; standard 0/1-knapsack space trick |

---

## Approach 1 — Brute Force (Enumerate Every Subset)

### Intuition

Every string is an independent include/exclude choice. Recurse over the strings: for string `i`, either skip it, or — if it still fits the leftover 0- and 1-budgets — take it, spend its cost, and add 1. The answer is the best over all `2^k` combinations. This is the literal definition the DP will accelerate.

### Algorithm

1. Pre-count `(zeros[i], ones[i])` for every string.
2. `rec(i, remM, remN)`: if `i == k`, return `0`.
3. `best = rec(i+1, remM, remN)` (skip string `i`).
4. If `zeros[i] <= remM` and `ones[i] <= remN`, `best = max(best, 1 + rec(i+1, remM-zeros[i], remN-ones[i]))`.
5. Return `best`.

### Complexity

- **Time:** O(2^k) — two branches per string; exponential, so only viable for tiny inputs.
- **Space:** O(k) — recursion depth (plus the O(k) precomputed counts).

### Code

```go
func bruteForce(strs []string, m int, n int) int {
	// Pre-count zeros/ones once so recursion is cheap.
	zeros := make([]int, len(strs))
	ones := make([]int, len(strs))
	for i, s := range strs {
		zeros[i], ones[i] = count(s)
	}

	var rec func(i, remM, remN int) int
	rec = func(i, remM, remN int) int {
		if i == len(strs) {
			return 0 // no strings left to consider
		}
		best := rec(i+1, remM, remN) // Option A: skip string i
		// Option B: take string i, only if its cost fits the remaining budget.
		if zeros[i] <= remM && ones[i] <= remN {
			take := 1 + rec(i+1, remM-zeros[i], remN-ones[i])
			if take > best {
				best = take
			}
		}
		return best
	}
	return rec(0, m, n)
}
```

### Dry Run

Example 2: `strs = ["10","0","1"]`, `m = 1`, `n = 1`. Costs: `"10"`→(1,1), `"0"`→(1,0), `"1"`→(0,1).

| Call | budget (remM,remN) | decision | contributes |
|------|--------------------|----------|-------------|
| rec(0) | (1,1) | skip "10" → rec(1); also try take "10" (fits) → 1+rec(1 with (0,0)) | max of two branches |
| ↳ take "10" | (0,0) | rec(1,(0,0)): "0" needs 1 zero ✗, "1" needs 1 one ✗ → 0 | 1 + 0 = **1** |
| ↳ skip "10" | (1,1) | rec(1,(1,1)) | see below |
| rec(1) | (1,1) | take "0"(1,0) → 1 + rec(2,(0,1)); take "1" path via skip | |
| ↳ take "0" | (0,1) | rec(2,(0,1)): take "1"(0,1) → 1 + rec(3) = 1 | 1 + 1 = **2** |

Best over branches = `2` ✔ (the subset `{"0","1"}`).

---

## Approach 2 — 3D DP (Item × m × n)

### Intuition

Turn the search into a table. This is a 0/1 knapsack where every item has **two** weights — a zero-cost and a one-cost — and unit value. Define `dp[i][j][k]` = the largest subset drawn from the first `i` strings that uses at most `j` zeros and `k` ones. Each string is either skipped (inherit `dp[i-1][j][k]`) or taken (add 1 to `dp[i-1][j-z][k-o]`).

### Algorithm

1. Initialise `dp[0][*][*] = 0` — no strings means an empty subset.
2. For each string `i` (1-indexed) with cost `(z, o)`, for all `j` in `0..m` and `k` in `0..n`:
   - `dp[i][j][k] = dp[i-1][j][k]` (skip).
   - If `j >= z` and `k >= o`: `dp[i][j][k] = max(dp[i][j][k], dp[i-1][j-z][k-o] + 1)` (take).
3. Answer: `dp[k][m][n]`.

### Complexity

- **Time:** O(k·m·n) — every one of the `(k+1)(m+1)(n+1)` cells is filled in O(1).
- **Space:** O(k·m·n) — the full 3D table; for `k=600, m=n=100` that is ~6M ints (large but shows the structure).

### Code

```go
func dp3D(strs []string, m int, n int) int {
	k := len(strs)
	// dp[i][j][k]: use first i strings, budget j zeros and k ones.
	dp := make([][][]int, k+1)
	for i := range dp {
		dp[i] = make([][]int, m+1)
		for j := range dp[i] {
			dp[i][j] = make([]int, n+1)
		}
	}
	for i := 1; i <= k; i++ {
		z, o := count(strs[i-1]) // cost of the i-th string (1-indexed)
		for j := 0; j <= m; j++ {
			for l := 0; l <= n; l++ {
				dp[i][j][l] = dp[i-1][j][l] // skip string i
				// Take string i if both budgets allow.
				if j >= z && l >= o {
					cand := dp[i-1][j-z][l-o] + 1
					if cand > dp[i][j][l] {
						dp[i][j][l] = cand
					}
				}
			}
		}
	}
	return dp[k][m][n]
}
```

### Dry Run

Example 2: `strs = ["10","0","1"]`, `m=1`, `n=1`. Costs: (1,1), (1,0), (0,1). Tracking only the answer cell `dp[i][1][1]`:

| i | string (cost) | dp[i][1][1] = max(skip, take) |
|---|---------------|-------------------------------|
| 0 | — | `0` (base) |
| 1 | "10" (1,1) | max(dp[0][1][1]=0, dp[0][0][0]+1=1) = **1** |
| 2 | "0" (1,0) | max(dp[1][1][1]=1, dp[1][0][1]+1). dp[1][0][1]=0 → 1) = **1** |
| 3 | "1" (0,1) | max(dp[2][1][1]=1, dp[2][1][0]+1). dp[2][1][0]=1 (took "0") → 2) = **2** |

`dp[3][1][1] = 2` ✔

---

## Approach 3 — 2D Rolling DP (Optimal)

### Intuition

`dp[i][…]` only reads `dp[i-1][…]` at the same or smaller `(j,k)`. So drop the item axis and keep one `(m+1)×(n+1)` grid, folding in strings one at a time. The catch is direction: sweep `j` and `k` **downward** (high → low). Then when we read `dp[j-z][k-o]`, it still holds the value from *before* this string was folded in — precisely the 0/1-knapsack requirement that each string be used at most once. Sweeping upward would let a single string be re-taken (the *unbounded* knapsack, wrong here).

### Algorithm

1. `dp` is `(m+1)×(n+1)`, all zeros.
2. For each string with cost `(z, o)`:
   - For `j` from `m` down to `z`, for `k` from `n` down to `o`:
     - `dp[j][k] = max(dp[j][k], dp[j-z][k-o] + 1)`.
3. Answer: `dp[m][n]`.

### Complexity

- **Time:** O(k·m·n) — same total work, processed string by string.
- **Space:** O(m·n) — one table reused across all strings (from ~6M ints down to ~10K).

### Code

```go
func dp2DRolling(strs []string, m int, n int) int {
	// dp[j][k]: best subset size achievable with budget j zeros and k ones,
	// considering the strings processed so far.
	dp := make([][]int, m+1)
	for j := range dp {
		dp[j] = make([]int, n+1)
	}
	for _, s := range strs {
		z, o := count(s) // cost of this string
		// Iterate DOWNWARD so dp[j-z][l-o] is still the "before this string"
		// value — this is what makes each string usable at most once.
		for j := m; j >= z; j-- {
			for l := n; l >= o; l-- {
				cand := dp[j-z][l-o] + 1 // take this string
				if cand > dp[j][l] {
					dp[j][l] = cand
				}
			}
		}
	}
	return dp[m][n]
}
```

### Dry Run

Example 2: `strs = ["10","0","1"]`, `m=1`, `n=1`. Grid cells shown as `dp[j][k]` for `(j,k)` in `{(0,0),(0,1),(1,0),(1,1)}`, all starting at 0.

| After folding | Downward updates | dp[1][1] | dp[1][0] | dp[0][1] |
|---------------|------------------|----------|----------|----------|
| — (init) | — | 0 | 0 | 0 |
| "10" (1,1) | `dp[1][1]=max(0, dp[0][0]+1)=1` | 1 | 0 | 0 |
| "0" (1,0) | `dp[1][1]=max(1, dp[0][1]+1=1)=1`; `dp[1][0]=max(0, dp[0][0]+1)=1` | 1 | 1 | 0 |
| "1" (0,1) | `dp[1][1]=max(1, dp[1][0]+1=2)=2`; `dp[0][1]=max(0, dp[0][0]+1)=1` | **2** | 1 | 1 |

`dp[1][1] = 2` ✔ — folding "1" last lets it stack on top of the already-taken "0".

---

## Key Takeaways

- **Two capacities ⇒ two DP dimensions.** Any "at most X of this AND at most Y of that" constraint is just a knapsack with a 2D weight. The item axis is what you compress away.
- **Direction encodes take-count.** 0/1 knapsack sweeps capacities **downhill** (each item once); unbounded knapsack sweeps **uphill** (item reusable). Same loop body, opposite order.
- **Value = 1 means "maximise cardinality".** The DP maximises how *many* items fit, not their summed value — a common special case.
- Pre-count costs once; recomputing `zeros/ones` inside the innermost loop would multiply the runtime by the string length.

---

## Related Problems

- LeetCode #416 — Partition Equal Subset Sum (single-capacity 0/1 knapsack)
- LeetCode #494 — Target Sum (0/1 knapsack counting variant)
- LeetCode #322 — Coin Change (unbounded knapsack, sweep uphill)
- LeetCode #879 — Profitable Schemes (two-capacity knapsack counting)
