# 0265 — Paint House II

> LeetCode #265 · Difficulty: Hard
> **Categories:** Array, Dynamic Programming

---

## Problem Statement

There are a row of `n` houses, each house can be painted with one of the `k` colors. The cost of painting each house with a certain color is different. You have to paint all the houses such that no two adjacent houses have the same color.

The cost of painting each house with a certain color is represented by an `n x k` cost matrix `costs`.

- For example, `costs[0][0]` is the cost of painting house `0` with color `0`; `costs[1][2]` is the cost of painting house `1` with color `2`, and so on...

Return the **minimum cost** to paint all houses.

**Example 1:**

```
Input: costs = [[1,5,3],[2,9,4]]
Output: 5
Explanation:
Paint house 0 into color 0, paint house 1 into color 2. Minimum cost: 1 + 4 = 5;
Or paint house 0 into color 2, paint house 1 into color 0. Minimum cost: 3 + 2 = 5.
```

**Example 2:**

```
Input: costs = [[1,3],[2,4]]
Output: 5
```

**Constraints:**

- `costs.length == n`
- `costs[i].length == k`
- `1 <= n <= 100`
- `2 <= k <= 20`
- `1 <= costs[i][j] <= 20`

**Follow up:** Could you solve it in `O(nk)` runtime?

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Amazon    | ★★★☆☆ Medium     | 2023          |
| Google    | ★★★☆☆ Medium     | 2023          |
| Microsoft | ★★☆☆☆ Low        | 2022          |
| Meta      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Dynamic Programming (1-D rolling state)** — `dp[j]` = best cost to finish the current house with color `j` → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)
- **Min/Second-Min trick** — track the two smallest previous costs to drop the inner O(k) scan → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | DP Full Table | O(n·k²) | O(n·k) | Clear baseline; small k |
| 2 | DP Min1/Min2 (Optimal) | O(n·k) | O(k) | Meets the follow-up bound |

---

## Approach 1 — DP Full Table

### Intuition
Let `dp[i][j]` be the minimum cost to paint houses `0..i` with house `i` colored `j`. To color house `i` with `j`, the previous house can be any color `p != j`, so `dp[i][j] = costs[i][j] + min over p!=j of dp[i-1][p]`. The answer is the minimum of the final row.

### Algorithm
1. Initialise `dp[0] = costs[0]`.
2. For each house `i ≥ 1` and color `j`: scan all `p != j` for the min of `dp[i-1][p]`, add `costs[i][j]`.
3. Return the min over `dp[n-1]`.

### Complexity
- **Time:** O(n·k²) — for each of `n` houses and `k` colors we scan `k` predecessors.
- **Space:** O(n·k) — the full table (reducible to two rows).

### Code
```go
func dpFullTable(costs [][]int) int {
	if len(costs) == 0 {
		return 0
	}
	n, k := len(costs), len(costs[0])
	dp := make([][]int, n) // dp[i][j] = best cost for houses 0..i ending color j
	for i := range dp {
		dp[i] = make([]int, k)
	}
	copy(dp[0], costs[0]) // base row: cost of painting only the first house
	for i := 1; i < n; i++ {
		for j := 0; j < k; j++ {
			best := -1 // min dp[i-1][p] over p != j
			for p := 0; p < k; p++ {
				if p == j { // adjacent houses cannot share a color
					continue
				}
				if best == -1 || dp[i-1][p] < best {
					best = dp[i-1][p]
				}
			}
			dp[i][j] = costs[i][j] + best // extend the cheapest allowed previous color
		}
	}
	ans := dp[n-1][0] // minimum over the last house's colors
	for j := 1; j < k; j++ {
		if dp[n-1][j] < ans {
			ans = dp[n-1][j]
		}
	}
	return ans
}
```

### Dry Run
Example 1: `costs = [[1,5,3],[2,9,4]]`, `n=2`, `k=3`. Base `dp[0] = [1,5,3]`.

House `i=1`, min over `p != j` of `dp[0]`:

| j | min of dp[0][p], p≠j | costs[1][j] | dp[1][j] |
|---|----------------------|-------------|----------|
| 0 | min(5,3)=3           | 2           | 5        |
| 1 | min(1,3)=1           | 9           | 10       |
| 2 | min(1,5)=1           | 4           | 5        |

Answer = `min(5,10,5) = 5`.

---

## Approach 2 — DP Min1/Min2 (Optimal)

### Intuition
For color `j`, the cheapest allowed previous cost is the overall smallest of the previous row — **unless** that smallest sits at color `j` itself, in which case we use the second smallest. Precompute `min1` (value + index `idx1`) and `min2` of the previous row, and the per-color choice becomes O(1), removing the inner scan.

### Algorithm
1. `prev = costs[0]`; compute `(min1, idx1, min2)` of `prev`.
2. For each house `i` and color `j`: `cur[j] = costs[i][j] + (min1 if j != idx1 else min2)`; recompute `(min1, idx1, min2)` from `cur`.
3. Return `min1` after the last house.

### Complexity
- **Time:** O(n·k) — one linear pass per house plus O(k) min-tracking.
- **Space:** O(k) — two rolling rows.

### Code
```go
func dpMinTwo(costs [][]int) int {
	if len(costs) == 0 {
		return 0
	}
	n, k := len(costs), len(costs[0])
	if k == 1 { // single color: houses can't be adjacent-distinct unless n==1
		if n == 1 {
			return costs[0][0]
		}
		return 0 // per constraints k>=2 when n>1; guard anyway
	}
	prev := make([]int, k)
	copy(prev, costs[0]) // base row = first house costs
	// min1 = smallest in prev, idx1 its index, min2 = second smallest.
	min1, idx1, min2 := minTwo(prev)
	for i := 1; i < n; i++ {
		cur := make([]int, k)
		for j := 0; j < k; j++ {
			best := min1 // cheapest previous cost avoiding color j
			if j == idx1 {
				best = min2 // smallest is at j itself ⇒ take second smallest
			}
			cur[j] = costs[i][j] + best
		}
		min1, idx1, min2 = minTwo(cur) // refresh trackers for next house
		prev = cur
	}
	return min1 // min over the last row is the answer
}
```

### Dry Run
Example 1: `costs = [[1,5,3],[2,9,4]]`. Base `prev = [1,5,3]` → `min1=1 (idx1=0)`, `min2=3`.

House `i=1`:

| j | j == idx1(0)? | best used | costs[1][j] | cur[j] |
|---|---------------|-----------|-------------|--------|
| 0 | yes           | min2 = 3  | 2           | 5      |
| 1 | no            | min1 = 1  | 9           | 10     |
| 2 | no            | min1 = 1  | 4           | 5      |

`cur = [5,10,5]` → new `min1 = 5`. Last house → answer `5`.

---

## Key Takeaways

- **"Different from the neighbour" DPs don't need the full O(k) predecessor scan.** The best previous choice is the global min, or the runner-up when the min collides with the current color — the classic *min1/min2* trick.
- Rolling two rows (`prev`/`cur`) drops space from O(n·k) to O(k).
- This is the k-color generalisation of Paint House (#256), which fixes `k = 3` and can hardcode the two "other" colors.

---

## Related Problems

- LeetCode #256 — Paint House (fixed 3 colors)
- LeetCode #276 — Paint Fence (adjacency with a run-length twist)
- LeetCode #198 — House Robber (adjacent-constraint 1-D DP)
- LeetCode #1289 — Minimum Falling Path Sum II (same min1/min2 optimisation on a grid)
