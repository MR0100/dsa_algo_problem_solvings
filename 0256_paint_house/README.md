# 0256 — Paint House

> LeetCode #256 · Difficulty: Medium
> **Categories:** Dynamic Programming, Array

---

## Problem Statement

There is a row of `n` houses, where each house can be painted one of three colors: red, blue, or green. The cost of painting each house with a certain color is different. You have to paint all the houses such that no two adjacent houses have the same color.

The cost of painting each house with a certain color is represented by an `n x 3` cost matrix `costs`.

- For example, `costs[0][0]` is the cost of painting house `0` with the color red; `costs[1][2]` is the cost of painting house `1` with color green, and so on...

Return *the minimum cost to paint all houses*.

**Example 1:**

```
Input: costs = [[17,2,17],[16,16,5],[14,3,19]]
Output: 10
Explanation: Paint house 0 into blue, paint house 1 into green, paint house 2 into blue.
Minimum cost: 2 + 5 + 3 = 10.
```

**Example 2:**

```
Input: costs = [[7,6,2]]
Output: 2
```

**Constraints:**

- `costs.length == n`
- `costs[i].length == 3`
- `1 <= n <= 100`
- `1 <= costs[i][j] <= 20`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★☆☆☆ Low        | 2023          |
| Apple      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Dynamic Programming (1D)** — the answer for each house depends only on the previous house's three color totals; a linear scan carrying three rolling values solves it → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | DP Top-Down (Memoized) | O(n) | O(n) | Most intuitive translation of "min cost from here given last color" |
| 2 | DP Bottom-Up (Full Table) | O(n) | O(n) | Clearest table formulation; easy to read off dp[i][c] |
| 3 | DP Rolling O(1) Space (Optimal) | O(n) | O(1) | Production answer; only three scalars carried |

---

## Approach 1 — DP Top-Down (Memoized)

### Intuition

Paint houses left to right. The only constraint linking houses is that adjacent ones differ, so the state you must remember when moving to house `i` is just the color used on house `i-1`. Define `solve(i, prev)` = the minimum cost to finish houses `i..n-1` given house `i-1` was `prev`. At house `i` try each color `c != prev`, pay `costs[i][c]`, and recurse. Memoize on `(i, prev)` because the same state is reached many ways.

### Algorithm

1. `solve(i, prev)`: if `i == n`, return `0` (nothing left to paint).
2. For each color `c` in `{0,1,2}` with `c != prev`: `candidate = costs[i][c] + solve(i+1, c)`.
3. Return the minimum candidate and cache it in `memo[i][prev+1]`.
4. Answer = `solve(0, -1)` (house 0 has no previous color).

### Complexity

- **Time:** O(n) — there are `n * 4` states (`prev` shifted into 0..3) and each does O(3) work.
- **Space:** O(n) — memo table plus recursion depth up to `n`.

### Code

```go
func dpTopDown(costs [][]int) int {
	n := len(costs)
	if n == 0 {
		return 0
	}
	// memo[i][prev+1]: prev ranges -1..2, shift by +1 to index 0..3.
	memo := make([][4]int, n)
	seen := make([][4]bool, n)

	var solve func(i, prev int) int
	solve = func(i, prev int) int {
		if i == n { // painted every house
			return 0
		}
		if seen[i][prev+1] { // already computed this state
			return memo[i][prev+1]
		}
		best := 1 << 30 // large sentinel (no valid cost can reach this)
		for c := 0; c < 3; c++ {
			if c == prev { // adjacent houses cannot share a color
				continue
			}
			cost := costs[i][c] + solve(i+1, c) // paint i with c, recurse
			if cost < best {
				best = cost
			}
		}
		seen[i][prev+1] = true // memoize the result for this state
		memo[i][prev+1] = best
		return best
	}
	return solve(0, -1) // house 0 has no left neighbor
}
```

### Dry Run

Example 1: `costs = [[17,2,17],[16,16,5],[14,3,19]]`, colors 0=red, 1=blue, 2=green.

| Call | prev | tries (c, cost) | returns |
|------|------|-----------------|---------|
| solve(2, ·) blue-parent | — | c=0:14, c=2:19 | 14 (red) |
| solve(2, ·) green-parent | — | c=0:14, c=1:3 | 3 (blue) |
| solve(1, blue=1) | 1 | c=0:16+solve(2,red-path)…, c=2:5+solve(2,2)=5+3=8 | 8 (green) |
| solve(0, -1) | -1 | c=1:2+solve(1,1)=2+8=10 (best) | 10 |

Result: `10` ✔ (house 0 blue → house 1 green → house 2 blue: 2 + 5 + 3).

---

## Approach 2 — DP Bottom-Up (Full Table)

### Intuition

Turn the recursion around: `dp[i][c]` = minimum total cost to paint houses `0..i` with house `i` painted color `c`. To paint house `i` color `c`, house `i-1` must be one of the other two colors, so add the cheaper of those. The answer is the minimum of the last row.

### Algorithm

1. Initialise `dp[0] = costs[0]`.
2. For `i = 1..n-1` and each color `c`: `dp[i][c] = costs[i][c] + min(dp[i-1][c'] for c' != c)`.
3. Return `min(dp[n-1][0], dp[n-1][1], dp[n-1][2])`.

### Complexity

- **Time:** O(n) — `n` rows, constant work per cell.
- **Space:** O(n) — the full `n x 3` table.

### Code

```go
func dpBottomUp(costs [][]int) int {
	n := len(costs)
	if n == 0 {
		return 0
	}
	dp := make([][3]int, n)
	dp[0] = [3]int{costs[0][0], costs[0][1], costs[0][2]} // base row = first house's costs
	for i := 1; i < n; i++ {
		// Each color depends on the best of the OTHER two colors above it.
		dp[i][0] = costs[i][0] + min2(dp[i-1][1], dp[i-1][2])
		dp[i][1] = costs[i][1] + min2(dp[i-1][0], dp[i-1][2])
		dp[i][2] = costs[i][2] + min2(dp[i-1][0], dp[i-1][1])
	}
	last := dp[n-1]
	return min2(last[0], min2(last[1], last[2])) // cheapest way to finish
}
```

### Dry Run

Example 1: `costs = [[17,2,17],[16,16,5],[14,3,19]]`.

| i | dp[i][0] (red) | dp[i][1] (blue) | dp[i][2] (green) |
|---|----------------|-----------------|------------------|
| 0 | 17 | 2 | 17 |
| 1 | 16 + min(2,17) = 18 | 16 + min(17,17) = 33 | 5 + min(17,2) = 7 |
| 2 | 14 + min(33,7) = 21 | 3 + min(18,7) = 10 | 19 + min(18,33) = 37 |

`min(21, 10, 37) = 10`. Result: `10` ✔

---

## Approach 3 — DP Rolling O(1) Space (Optimal)

### Intuition

`dp[i]` reads only `dp[i-1]`, so keep just the previous house's three totals in variables `r, b, g` and roll them forward. No table needed.

### Algorithm

1. `r, b, g = costs[0][0], costs[0][1], costs[0][2]`.
2. For each next house compute `nr = costs[i][0]+min(b,g)`, `nb = costs[i][1]+min(r,g)`, `ng = costs[i][2]+min(r,b)`, then set `r,b,g = nr,nb,ng`.
3. Return `min(r, b, g)`.

### Complexity

- **Time:** O(n) — single pass.
- **Space:** O(1) — three rolling scalars.

### Code

```go
func dpRolling(costs [][]int) int {
	n := len(costs)
	if n == 0 {
		return 0
	}
	r, b, g := costs[0][0], costs[0][1], costs[0][2] // prev house's 3 totals
	for i := 1; i < n; i++ {
		// Compute this house's totals from the previous three, then roll over.
		nr := costs[i][0] + min2(b, g) // red now: cheaper of prev blue/green
		nb := costs[i][1] + min2(r, g) // blue now: cheaper of prev red/green
		ng := costs[i][2] + min2(r, b) // green now: cheaper of prev red/blue
		r, b, g = nr, nb, ng
	}
	return min2(r, min2(b, g)) // best final choice
}
```

### Dry Run

Example 1: `costs = [[17,2,17],[16,16,5],[14,3,19]]`.

| i | r (red) | b (blue) | g (green) |
|---|---------|----------|-----------|
| 0 (init) | 17 | 2 | 17 |
| 1 | 16 + min(2,17) = 18 | 16 + min(17,17) = 33 | 5 + min(17,2) = 7 |
| 2 | 14 + min(33,7) = 21 | 3 + min(18,7) = 10 | 19 + min(18,33) = 37 |

`min(21, 10, 37) = 10`. Result: `10` ✔

---

## Key Takeaways

- **State = the minimal history that constrains the next choice.** Here that is just the previous house's color, giving a tiny 3-way DP rather than an exponential search.
- **"Min over the other two" is the recurring move** in fixed-color-count DP (Paint House, Paint House II uses the same idea generalized to k colors with a min/second-min trick).
- **Table → rolling variables** whenever `dp[i]` reads only `dp[i-1]`: drops space from O(n) to O(1) for free.

---

## Related Problems

- LeetCode #265 — Paint House II (k colors; min + second-min optimization)
- LeetCode #198 — House Robber (adjacent-constraint 1D DP)
- LeetCode #213 — House Robber II (circular variant)
- LeetCode #746 — Min Cost Climbing Stairs (rolling-variable 1D DP)
