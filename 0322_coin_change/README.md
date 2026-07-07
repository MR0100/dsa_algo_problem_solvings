# 0322 — Coin Change

> LeetCode #322 · Difficulty: Medium
> **Categories:** Dynamic Programming, BFS, Greedy, Unbounded Knapsack

---

## Problem Statement

You are given an integer array `coins` representing coins of different
denominations and an integer `amount` representing a total amount of money.

Return the fewest number of coins that you need to make up that amount. If that
amount of money cannot be made up by any combination of the coins, return `-1`.

You may assume that you have an infinite number of each kind of coin.

**Example 1:**

```
Input:  coins = [1,2,5], amount = 11
Output: 3
Explanation: 11 = 5 + 5 + 1
```

**Example 2:**

```
Input:  coins = [2], amount = 3
Output: -1
```

**Example 3:**

```
Input:  coins = [1], amount = 0
Output: 0
```

**Constraints:**

- `1 <= coins.length <= 12`
- `1 <= coins[i] <= 2^31 - 1`
- `0 <= amount <= 10^4`

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Amazon    | ★★★★★ Very High  | 2024          |
| Google    | ★★★★☆ High       | 2024          |
| Microsoft | ★★★★☆ High       | 2024          |
| Uber      | ★★★☆☆ Medium     | 2023          |
| Bloomberg | ★★★☆☆ Medium     | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Dynamic Programming (1D)** — `dp[a]` = fewest coins for amount `a`, an
  unbounded-knapsack recurrence → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)
- **BFS** — amounts as graph nodes; shortest path from `amount` to `0` is the
  fewest coins → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Greedy** — sort coins descending to find a tight bound fast for DFS pruning
  → see [`/dsa/greedy.md`](/dsa/greedy.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Top-Down DP (Memoized) | O(amount·n) | O(amount) | Natural recursive framing |
| 2 | Bottom-Up DP (Optimal) | O(amount·n) | O(amount) | Cleanest, no recursion overhead |
| 3 | BFS on Amounts | O(amount·n) | O(amount) | "Fewest steps" ↔ shortest path intuition |
| 4 | Greedy + DFS with Pruning | Exponential (fast in practice) | O(n) | When coin count is tiny; pruning shines |

(`n = len(coins)`)

---

## Approach 1 — Top-Down DP (Memoized Recursion)

### Intuition
`minCoins(r) = 1 + min over coins c of minCoins(r-c)`, with `minCoins(0) = 0`.
Recursion revisits the same remaining amount many times, so cache each answer.

### Algorithm
1. `memo[r]` caches the answer for amount `r`; `-2` means "not computed".
2. `solve(r)`: return 0 if `r == 0`, return -1 if `r < 0`.
3. Try each coin, recurse on `r-c`, keep the smallest `sub+1` where `sub >= 0`.
4. Cache and return; `-1` if no coin leads anywhere.

### Complexity
- **Time:** O(amount·n) — each of `amount+1` states tries every coin once.
- **Space:** O(amount) — memo table plus recursion depth.

### Code
```go
func dpTopDown(coins []int, amount int) int {
	memo := make([]int, amount+1)
	for i := range memo {
		memo[i] = -2
	}
	var solve func(r int) int
	solve = func(r int) int {
		if r == 0 {
			return 0
		}
		if r < 0 {
			return -1
		}
		if memo[r] != -2 {
			return memo[r]
		}
		best := -1
		for _, c := range coins {
			sub := solve(r - c)
			if sub >= 0 && (best == -1 || sub+1 < best) {
				best = sub + 1
			}
		}
		memo[r] = best
		return best
	}
	return solve(amount)
}
```

### Dry Run
Example 1: `coins = [1,2,5]`, `amount = 11`.

| Call | tries c=1,2,5 → sub-calls | result |
|------|---------------------------|--------|
| solve(11) | 1+min(solve10, solve9, solve6) | 3 |
| solve(6)  | 1+min(solve5, solve4, solve1) | 2 (5+1) |
| solve(5)  | 1+min(solve4, solve3, solve0=0) | 1 |
| solve(1)  | 1+solve(0)=0 | 1 |
| solve(0)  | base | 0 |

`solve(11) = 1 + solve(6) = 1 + (1 + solve(5)) = 1 + 1 + 1 = 3`. Output `3`.

---

## Approach 2 — Bottom-Up DP (Unbounded Knapsack) (Optimal)

### Intuition
To make amount `a`, pick a last coin `c <= a`, leaving the already-solved
sub-amount `a-c`. So `dp[a] = 1 + min over c<=a of dp[a-c]`. Build from small `a`
upward. Any entry left at the "infinity" sentinel is unreachable.

### Algorithm
1. `dp[0] = 0`; `dp[a] = amount+1` (infinity) for `a >= 1`.
2. For `a` from 1..amount, for each coin `c <= a`: `dp[a] = min(dp[a], dp[a-c]+1)`.
3. Return `dp[amount]` if finite, else `-1`.

### Complexity
- **Time:** O(amount·n) — nested loops over amounts and coins.
- **Space:** O(amount) — the dp table.

### Code
```go
func dpBottomUp(coins []int, amount int) int {
	inf := amount + 1
	dp := make([]int, amount+1)
	for a := 1; a <= amount; a++ {
		dp[a] = inf
	}
	for a := 1; a <= amount; a++ {
		for _, c := range coins {
			if c <= a && dp[a-c]+1 < dp[a] {
				dp[a] = dp[a-c] + 1
			}
		}
	}
	if dp[amount] == inf {
		return -1
	}
	return dp[amount]
}
```

### Dry Run
Example 1: `coins = [1,2,5]`, `amount = 11`, `inf = 12`.

| a  | dp[a-1]+1 (c=1) | dp[a-2]+1 (c=2) | dp[a-5]+1 (c=5) | dp[a] |
|----|-----------------|-----------------|-----------------|-------|
| 0  | —               | —               | —               | 0 |
| 1  | dp[0]+1=1       | —               | —               | 1 |
| 2  | dp[1]+1=2       | dp[0]+1=1       | —               | 1 |
| 3  | dp[2]+1=2       | dp[1]+1=2       | —               | 2 |
| 4  | dp[3]+1=3       | dp[2]+1=2       | —               | 2 |
| 5  | dp[4]+1=3       | dp[3]+1=3       | dp[0]+1=1       | 1 |
| 6  | dp[5]+1=2       | dp[4]+1=3       | dp[1]+1=2       | 2 |
| ...| ...             | ...             | ...             | ... |
| 10 | dp[9]+1         | dp[8]+1         | dp[5]+1=2       | 2 |
| 11 | dp[10]+1=3      | dp[9]+1         | dp[6]+1=3       | 3 |

`dp[11] = 3`. Output `3`.

---

## Approach 3 — BFS on Amounts

### Intuition
Model amounts as nodes; from `a` an edge goes to `a-c` for each coin `c`. The
fewest coins is the shortest path (fewest edges) from `amount` to `0`. BFS visits
by distance, so the first arrival at `0` uses the minimum coins.

### Algorithm
1. Queue starts with `amount` at level 0; mark visited.
2. Each level, subtract every coin; on hitting `0` return `level+1`; push new
   positive amounts not yet visited.
3. If the queue drains without reaching `0`, return `-1`.

### Complexity
- **Time:** O(amount·n) — each amount enqueued at most once.
- **Space:** O(amount) — visited set and queue.

### Code
```go
func bfs(coins []int, amount int) int {
	if amount == 0 {
		return 0
	}
	visited := make([]bool, amount+1)
	visited[amount] = true
	queue := []int{amount}
	level := 0
	for len(queue) > 0 {
		level++
		next := []int{}
		for _, cur := range queue {
			for _, c := range coins {
				rem := cur - c
				if rem == 0 {
					return level
				}
				if rem > 0 && !visited[rem] {
					visited[rem] = true
					next = append(next, rem)
				}
			}
		}
		queue = next
	}
	return -1
}
```

### Dry Run
Example 1: `coins = [1,2,5]`, `amount = 11`.

| level | frontier (queue)              | new amounts (via -1,-2,-5) |
|-------|-------------------------------|----------------------------|
| 1     | {11}                          | 10, 9, 6                   |
| 2     | {10, 9, 6}                    | 8, 7, 5, 4, 1              |
| 3     | {8, 7, 5, 4, 1}               | 6-1=5..., **1-1=0 → return 3** |

At level 3, `1 - 1 = 0`, so BFS returns `3`. Output `3`.

---

## Approach 4 — Greedy + DFS with Pruning

### Intuition
Plain greedy (always take the biggest coin) is wrong for arbitrary coin sets,
but sorting coins **descending** makes a near-optimal answer appear early, giving
a tight `best` bound. DFS then explores fewer copies of each big coin, pruning
any branch whose optimistic coin count can't beat `best`.

### Algorithm
1. Sort coins descending.
2. `dfs(idx, remaining, count)`: if `remaining == 0`, update `best`.
3. Take `k = remaining/coins[idx]` copies down to 0; skip a branch when
   `best != -1 && count+k >= best` (can't improve).

### Complexity
- **Time:** Exponential worst case, but pruning is fast on the small `coins.length <= 12`.
- **Space:** O(n) — recursion depth.

### Code
```go
func greedyDFS(coins []int, amount int) int {
	sorted := make([]int, len(coins))
	copy(sorted, coins)
	sort.Sort(sort.Reverse(sort.IntSlice(sorted)))
	best := -1
	var dfs func(idx, remaining, count int)
	dfs = func(idx, remaining, count int) {
		if remaining == 0 {
			if best == -1 || count < best {
				best = count
			}
			return
		}
		if idx == len(sorted) {
			return
		}
		coin := sorted[idx]
		maxK := remaining / coin
		for k := maxK; k >= 0; k-- {
			if best != -1 && count+k >= best {
				continue
			}
			dfs(idx+1, remaining-k*coin, count+k)
		}
	}
	dfs(0, amount, 0)
	return best
}
```

### Dry Run
Example 1: `coins = [1,2,5]` → sorted `[5,2,1]`, `amount = 11`.

| step | idx (coin) | remaining | count | note |
|------|-----------|-----------|-------|------|
| 1 | 0 (5) | 11 | 0 | maxK=2, try k=2 |
| 2 | 1 (2) | 1  | 2 | maxK=0, k=0 |
| 3 | 2 (1) | 1  | 2 | k=1 → remaining 0, count 3 → best=3 |
| 4 | back to idx 0, k=1 (one 5) | 6 | 1 | prune paths where count+k ≥ 3 |
| … | further branches pruned by `best=3` |   |   | |

First full path 5+5+1 gives `best = 3`; pruning discards worse branches. Output `3`.

---

## Key Takeaways
- Coin Change is the canonical **unbounded knapsack / minimization DP**:
  `dp[a] = 1 + min(dp[a-c])`.
- "Fewest steps to reach a target" maps naturally onto **BFS shortest path**.
- Greedy is not correct alone here; it earns its keep only as a **bound for DFS pruning**.
- Use a sentinel larger than any possible count (`amount+1`) as "infinity" so a
  single `min` comparison also detects unreachable amounts.

---

## Related Problems
- LeetCode #518 — Coin Change II (count combinations, not minimum)
- LeetCode #279 — Perfect Squares (same min-count DP with square "coins")
- LeetCode #983 — Minimum Cost For Tickets (interval DP variant)
- LeetCode #377 — Combination Sum IV (ordered count DP)
