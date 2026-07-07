# 0309 — Best Time to Buy and Sell Stock with Cooldown

> LeetCode #309 · Difficulty: Medium
> **Categories:** Array, Dynamic Programming, State Machine

---

## Problem Statement

You are given an array `prices` where `prices[i]` is the price of a given stock on the `iᵗʰ` day.

Find the maximum profit you can achieve. You may complete as many transactions as you like (i.e., buy one and sell one share of the stock multiple times) with the following restrictions:

- After you sell your stock, you **cannot** buy stock on the next day (i.e., cooldown one day).

**Note:** You may not engage in multiple transactions simultaneously (i.e., you must sell the stock before you buy again).

**Example 1:**

```
Input: prices = [1,2,3,0,2]
Output: 3
Explanation: transactions = [buy, sell, cooldown, buy, sell]
```

**Example 2:**

```
Input: prices = [1]
Output: 0
```

**Constraints:**

- `1 <= prices.length <= 5000`
- `0 <= prices[i] <= 1000`

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Amazon    | ★★★★☆ High       | 2024          |
| Google    | ★★★★☆ High       | 2024          |
| Bloomberg | ★★★☆☆ Medium     | 2023          |
| Microsoft | ★★★☆☆ Medium     | 2023          |
| Meta      | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Dynamic Programming (state machine)** — model each day as being in one of {hold, sold, rest} and define transitions between them → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)
- **Top-Down Memoization** — recurse on (day, holding) and cache overlapping subproblems → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)
- **Rolling-array optimization** — transitions depend only on the previous day, so O(1) scalars suffice → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Top-Down DP (memoized) | O(n) | O(n) | Most intuitive; mirrors the choices |
| 2 | Bottom-Up DP (3 state arrays) | O(n) | O(n) | Clear state-machine formulation |
| 3 | O(1) Space State Machine (Optimal) | O(n) | O(1) | Final answer; least memory |

---

## Approach 1 — Top-Down DP (Memoized)

### Intuition

At each day, in each state (holding a share or not), pick the better of acting or waiting. If not holding: skip or buy (→ holding). If holding: skip or sell (→ cooldown, so the next buy can only happen on `day+2`). Subproblems on `(day, holding)` overlap, so memoize.

### Algorithm

1. `solve(day, holding)`: if `day >= n` return 0.
2. `skip = solve(day+1, holding)`.
3. If holding: `act = prices[day] + solve(day+2, 0)` (sell, then cooldown).
4. Else: `act = -prices[day] + solve(day+1, 1)` (buy).
5. Return `max(skip, act)`; answer is `solve(0, 0)`.

### Complexity

- **Time:** O(n) — 2n states, each resolved in O(1).
- **Space:** O(n) — memo table plus recursion stack.

### Code

```go
func dpTopDown(prices []int) int {
	n := len(prices)
	if n == 0 {
		return 0
	}
	memo := make([][2]int, n)
	seen := make([][2]bool, n)
	var solve func(day, holding int) int
	solve = func(day, holding int) int {
		if day >= n {
			return 0
		}
		if seen[day][holding] {
			return memo[day][holding]
		}
		skip := solve(day+1, holding)
		var act int
		if holding == 1 {
			act = prices[day] + solve(day+2, 0)
		} else {
			act = -prices[day] + solve(day+1, 1)
		}
		best := max(skip, act)
		memo[day][holding] = best
		seen[day][holding] = true
		return best
	}
	return solve(0, 0)
}
```

### Dry Run

`prices = [1,2,3,0,2]`, entry `solve(0, 0)` (indices are days):

| Call | Choice taken | Value |
|------|--------------|-------|
| solve(3,0) | buy at 0 → -0 + solve(4,1) | 2 |
| solve(4,1) | sell at 2 → 2 + solve(6,0)=2 | 2 |
| solve(1,0) | buy at 2 → best downstream | 2 |
| solve(0,0) | max(skip=solve(1,0), buy at 1 → -1 + solve(1,1)) | **3** |

The optimal path buys day 0 (−1), sells day 2 (+3 → net 2), cooldown day 3, buys day 3 (0), sells day 4 (+2). Total profit **3**.

---

## Approach 2 — Bottom-Up DP (State Arrays)

### Intuition

Model days as a state machine with three states per day:

- `hold[i]` — max profit on day `i` while **holding** a share.
- `sold[i]` — max profit on day `i` having **just sold** (cooldown begins).
- `rest[i]` — max profit on day `i` **idle** (free to buy).

Transitions bake in the cooldown: you can only buy from `rest`, and you reach `rest` only one day after `sold`.

### Algorithm

1. `hold[0] = -prices[0]`, `sold[0] = 0`, `rest[0] = 0`.
2. For `i` from 1: 
   - `hold[i] = max(hold[i-1], rest[i-1] - prices[i])`
   - `sold[i] = hold[i-1] + prices[i]`
   - `rest[i] = max(rest[i-1], sold[i-1])`
3. Answer = `max(sold[n-1], rest[n-1])` (holding at the end is never optimal).

### Complexity

- **Time:** O(n) — one pass, constant work per day.
- **Space:** O(n) — three arrays (reducible to O(1)).

### Code

```go
func dpBottomUp(prices []int) int {
	n := len(prices)
	if n == 0 {
		return 0
	}
	hold := make([]int, n)
	sold := make([]int, n)
	rest := make([]int, n)
	hold[0] = -prices[0]
	sold[0] = 0
	rest[0] = 0
	for i := 1; i < n; i++ {
		hold[i] = max(hold[i-1], rest[i-1]-prices[i])
		sold[i] = hold[i-1] + prices[i]
		rest[i] = max(rest[i-1], sold[i-1])
	}
	return max(sold[n-1], rest[n-1])
}
```

### Dry Run

`prices = [1,2,3,0,2]`:

| i | price | hold | sold | rest |
|---|-------|------|------|------|
| 0 | 1 | -1 | 0 | 0 |
| 1 | 2 | max(-1, 0-2)=-1 | -1+2=1 | max(0,0)=0 |
| 2 | 3 | max(-1, 0-3)=-1 | -1+3=2 | max(0,1)=1 |
| 3 | 0 | max(-1, 1-0)=1 | -1+0=-1 | max(1,2)=2 |
| 4 | 2 | max(1, 2-2)=1 | 1+2=3 | max(2,-1)=2 |

Answer = `max(sold[4]=3, rest[4]=2) = ` **3**.

---

## Approach 3 — O(1) Space State Machine (Optimal)

### Intuition

Each transition depends only on day `i-1`, so three scalars replace the arrays. This is the tightest form of the cooldown state machine.

### Algorithm

1. `hold = -prices[0]`, `sold = 0`, `rest = 0`.
2. For each subsequent price: 
   - `prevSold = sold`
   - `sold = hold + price`
   - `hold = max(hold, rest - price)`
   - `rest = max(rest, prevSold)`
3. Answer = `max(sold, rest)`.

> Order matters: capture `prevSold` before overwriting `sold`, and compute `sold` from the old `hold` before updating `hold`.

### Complexity

- **Time:** O(n) — single pass.
- **Space:** O(1) — three scalars.

### Code

```go
func dpConstantSpace(prices []int) int {
	n := len(prices)
	if n == 0 {
		return 0
	}
	hold := -prices[0]
	sold := 0
	rest := 0
	for i := 1; i < n; i++ {
		prevSold := sold
		sold = hold + prices[i]
		hold = max(hold, rest-prices[i])
		rest = max(rest, prevSold)
	}
	return max(sold, rest)
}
```

### Dry Run

`prices = [1,2,3,0,2]`:

| i | price | prevSold | sold | hold | rest |
|---|-------|----------|------|------|------|
| start | — | — | 0 | -1 | 0 |
| 1 | 2 | 0 | -1+2=1 | max(-1,0-2)=-1 | max(0,0)=0 |
| 2 | 3 | 1 | -1+3=2 | max(-1,0-3)=-1 | max(0,1)=1 |
| 3 | 0 | 2 | -1+0=-1 | max(-1,1-0)=1 | max(1,2)=2 |
| 4 | 2 | -1 | 1+2=3 | max(1,2-2)=1 | max(2,-1)=2 |

Answer = `max(sold=3, rest=2) = ` **3**.

---

## Key Takeaways

- **State-machine DP:** when transactions have a lifecycle (buy → hold → sell → cooldown), name the states and write transitions between them — the cooldown is enforced structurally, not with special-case code.
- **Only `rest` can buy**, and `rest` is reachable only one day after `sold`; that single edge encodes the one-day cooldown.
- **Rolling optimization:** any DP whose transitions look back only one step collapses to O(1) space — just be careful about update order.
- **Holding at the end is never optimal**, so the answer maxes over the two "cash" states.

---

## Related Problems

- LeetCode #121 — Best Time to Buy and Sell Stock (single transaction)
- LeetCode #122 — Best Time to Buy and Sell Stock II (unlimited, no cooldown)
- LeetCode #123 — Best Time to Buy and Sell Stock III (at most two transactions)
- LeetCode #188 — Best Time to Buy and Sell Stock IV (at most k transactions)
- LeetCode #714 — Best Time to Buy and Sell Stock with Transaction Fee (fee instead of cooldown)
