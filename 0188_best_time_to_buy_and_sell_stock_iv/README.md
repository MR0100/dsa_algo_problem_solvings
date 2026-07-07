# 0188 — Best Time to Buy and Sell Stock IV

> LeetCode #188 · Difficulty: Hard
> **Categories:** Array, Dynamic Programming, Greedy

---

## Problem Statement

You are given an integer array `prices` where `prices[i]` is the price of a given stock on the `i`-th day, and an integer `k`.

Find the maximum profit you can achieve. You may complete at most `k` transactions: i.e. you may buy at most `k` times and sell at most `k` times.

**Note:** You may not engage in multiple transactions simultaneously (i.e., you must sell the stock before you buy again).

**Example 1:**
```
Input: k = 2, prices = [2,4,1]
Output: 2
Explanation: Buy on day 1 (price = 2) and sell on day 2 (price = 4), profit = 4-2 = 2.
```

**Example 2:**
```
Input: k = 2, prices = [3,2,6,5,0,3]
Output: 7
Explanation: Buy on day 2 (price = 2) and sell on day 3 (price = 6), profit = 6-2 = 4.
Then buy on day 5 (price = 0) and sell on day 6 (price = 3), profit = 3-0 = 3.
```

**Constraints:**
- `1 <= k <= 100`
- `1 <= prices.length <= 1000`
- `0 <= prices[i] <= 1000`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Meta       | ★★★☆☆ Medium     | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |
| Citadel    | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **2D Dynamic Programming** — the natural state is (transactions used, day); dp[t][d] tables with the maxDiff sweep give O(n·k) → see [`/dsa/dynamic_programming_2d.md`](/dsa/dynamic_programming_2d.md)
- **1D DP / Rolling State (state machine)** — only the previous relaxation of each of the 2k trading states matters, so the day dimension compresses into two O(k) arrays → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)
- **Greedy** — when `k >= n/2` the transaction cap can never bind, and summing every positive day-to-day rise is provably optimal → see [`/dsa/greedy.md`](/dsa/greedy.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Exhaustive Recursion) | O(2ⁿ) | O(n) stack | Never in production — only to state the decision space |
| 2 | DP Top-Down (Memoization) | O(n·k) | O(n·k) | Fastest path from recursion to an accepted solution |
| 3 | DP Bottom-Up (dp[t][d] + maxDiff) | O(n·k) | O(n·k) | Shows the classic table + the O(n²·k) → O(n·k) optimization |
| 4 | State Machine, Space Optimized (Optimal) | O(n·k) | O(k) | The interview finisher: minimal memory + greedy shortcut for huge k |

---

## Approach 1 — Brute Force (Exhaustive Recursion)

### Intuition
Each day offers at most two sensible moves: **do nothing**, or **trade** — buy if we hold nothing, sell if we hold a share. Recursing over both choices for every day enumerates every legal schedule with at most `k` completed transactions. Counting a transaction on the **sell** (not the buy) makes the "at most k" bookkeeping trivial: `remaining` only decreases when a buy–sell pair actually completes.

### Algorithm
1. Define `solve(day, remaining, holding)` = best extra profit from `day` onward.
2. Base case: `day == n` or `remaining == 0` → return 0.
3. Option A: skip today → `solve(day+1, remaining, holding)`.
4. Option B: if holding, sell → `prices[day] + solve(day+1, remaining-1, false)`; if not holding, buy → `-prices[day] + solve(day+1, remaining, true)`.
5. Return the max of the options; the answer is `solve(0, k, false)`.

### Complexity
- **Time:** O(2ⁿ) — two independent branches per day and no state reuse, so the decision tree doubles each level.
- **Space:** O(n) — recursion depth is one stack frame per day.

### Code
```go
func bruteForce(k int, prices []int) int {
	var solve func(day, remaining int, holding bool) int
	solve = func(day, remaining int, holding bool) int {
		// no days left, or no transactions left to earn with → nothing more to gain
		if day == len(prices) || remaining == 0 {
			return 0
		}
		best := solve(day+1, remaining, holding) // choice 1: do nothing today
		if holding {
			// choice 2a: sell today — realise the price, one transaction completes
			if p := prices[day] + solve(day+1, remaining-1, false); p > best {
				best = p
			}
		} else {
			// choice 2b: buy today — pay the price, transaction counted on the future sell
			if p := -prices[day] + solve(day+1, remaining, true); p > best {
				best = p
			}
		}
		return best
	}
	return solve(0, k, false) // start on day 0, k transactions available, holding nothing
}
```

### Dry Run
Example 1: `k = 2, prices = [2,4,1]` — call tree of `solve(0, 2, false)` (profitable spine shown):

| Call | Options evaluated | Value |
|------|-------------------|-------|
| `solve(2, 1, true)` | skip → 0; sell at 1 → `1 + 0 = 1` | 1 |
| `solve(1, 2, true)` | skip → `solve(2,2,true)` = 1; sell at 4 → `4 + solve(2,1,false)` = `4 + max(0, -1+0)` = 4 | **4** |
| `solve(1, 2, false)` | skip → `solve(2,2,false)` = 0; buy at 4 → `-4 + solve(2,2,true)` = `-4+1` = -3 | 0 |
| `solve(0, 2, false)` | skip → `solve(1,2,false)` = 0; **buy at 2** → `-2 + solve(1,2,true)` = `-2+4` = **2** | **2** ✓ |

Best schedule found: buy day 0 (price 2), sell day 1 (price 4) → profit 2.

---

## Approach 2 — DP Top-Down (Memoized Recursion)

### Intuition
The brute force keeps re-solving identical futures: once we stand at day `d` with `t` transactions left and a given holding flag, the best continuation is fixed — *how* we got there is irrelevant. There are only `n · (k+1) · 2` distinct states, so caching each state's answer collapses the exponential tree into one visit per state.

### Algorithm
1. Allocate `memo[day][remaining][holding]`, initialised to `-1` ("unset" — safe because every real answer is ≥ 0: doing nothing forever earns 0).
2. Run the identical recursion as Approach 1, but return the cached value when present and store `best` before returning.
3. Answer: `solve(0, k, 0)`.

### Complexity
- **Time:** O(n·k) — `n·(k+1)·2` states, each resolved once with O(1) transition work.
- **Space:** O(n·k) — the memo table, plus an O(n) recursion stack.

### Code
```go
func dpTopDown(k int, prices []int) int {
	n := len(prices)
	if n == 0 || k == 0 {
		return 0 // no days or no allowed transactions → no profit possible
	}
	// memo[day][remaining][holding]; -1 = unset (all real answers are >= 0)
	memo := make([][][2]int, n)
	for d := range memo {
		memo[d] = make([][2]int, k+1)
		for t := range memo[d] {
			memo[d][t] = [2]int{-1, -1}
		}
	}
	var solve func(day, remaining, holding int) int
	solve = func(day, remaining, holding int) int {
		if day == n || remaining == 0 {
			return 0 // out of days or out of transactions
		}
		if memo[day][remaining][holding] != -1 {
			return memo[day][remaining][holding] // state already solved once
		}
		best := solve(day+1, remaining, holding) // skip today
		if holding == 1 {
			// sell: pocket today's price, the transaction completes
			if p := prices[day] + solve(day+1, remaining-1, 0); p > best {
				best = p
			}
		} else {
			// buy: spend today's price, stay on the same transaction budget
			if p := -prices[day] + solve(day+1, remaining, 1); p > best {
				best = p
			}
		}
		memo[day][remaining][holding] = best // cache before returning
		return best
	}
	return solve(0, k, 0)
}
```

### Dry Run
Example 1: `k = 2, prices = [2,4,1]` — states resolved (deepest first):

| State (day, remaining, holding) | skip | trade | memo value |
|---------------------------------|------|-------|------------|
| (2, 1, 1) | 0 | sell: `1+0 = 1` | 1 |
| (2, 2, 1) | 0 | sell: `1+0 = 1` | 1 |
| (2, 1, 0) | 0 | buy: `-1+0 = -1` | 0 |
| (2, 2, 0) | 0 | buy: `-1+0 = -1` | 0 |
| (1, 2, 1) | memo(2,2,1) = 1 | sell: `4 + memo(2,1,0)` = 4 | **4** |
| (1, 2, 0) | memo(2,2,0) = 0 | buy: `-4 + memo(2,2,1)` = -3 | 0 |
| (0, 2, 0) | memo(1,2,0) = 0 | buy: `-2 + memo(1,2,1)` = **2** | **2** ✓ |

Every state is computed exactly once; the top state returns 2.

---

## Approach 3 — DP Bottom-Up (Transactions × Days + maxDiff Trick)

### Intuition
Define `dp[t][d]` = best profit using **at most** `t` transactions within days `0..d`. On day `d` we either don't trade (`dp[t][d-1]`) or we **sell** on day `d` having bought on some earlier day `m`, worth `prices[d] − prices[m] + dp[t-1][m]`. Naively scanning every `m` gives O(n²·k). The rescue: the sell-on-`d` term is `prices[d] + (dp[t-1][m] − prices[m])`, and the parenthesised part is independent of `d` — so carry its running maximum, `maxDiff`, forward as `d` advances. Each row becomes a single O(n) sweep.

### Algorithm
1. Allocate `dp` with `k+1` rows and `n` columns; row `t=0` (no transactions) and column `d=0` (single day) stay 0.
2. For each `t = 1..k`:
   1. `maxDiff = dp[t-1][0] − prices[0]` (day 0 as the initial candidate buy day).
   2. For each `d = 1..n-1`: `dp[t][d] = max(dp[t][d-1], prices[d] + maxDiff)`, then `maxDiff = max(maxDiff, dp[t-1][d] − prices[d])`.
3. Return `dp[k][n-1]`.

### Complexity
- **Time:** O(n·k) — k rows, each swept once in O(n) thanks to the incremental `maxDiff` (down from O(n²·k) without it).
- **Space:** O(n·k) — the full table is kept for clarity; only the previous row is actually read, so O(n) is possible (Approach 4 goes further to O(k)).

### Code
```go
func dpBottomUp(k int, prices []int) int {
	n := len(prices)
	if n < 2 || k == 0 {
		return 0 // need at least two days (a buy and a sell) to profit
	}
	// dp[t][d] = best profit with at most t transactions over days 0..d
	dp := make([][]int, k+1)
	for t := range dp {
		dp[t] = make([]int, n) // row t=0 and column d=0 remain 0 (no trades possible)
	}
	for t := 1; t <= k; t++ {
		// maxDiff = best (dp[t-1][m] - prices[m]) over buy days m seen so far
		maxDiff := dp[t-1][0] - prices[0]
		for d := 1; d < n; d++ {
			// don't trade today, or sell today on top of the best earlier buy
			dp[t][d] = max(dp[t][d-1], prices[d]+maxDiff)
			// let day d compete as the buy day for future sell days
			maxDiff = max(maxDiff, dp[t-1][d]-prices[d])
		}
	}
	return dp[k][n-1] // at most k transactions over all n days
}
```

### Dry Run
Example 1: `k = 2, prices = [2,4,1]` (n = 3)

Row `t = 1` (start `maxDiff = dp[0][0] − 2 = −2`):

| d | prices[d] | dp[1][d] = max(dp[1][d−1], prices[d]+maxDiff) | maxDiff after update |
|---|-----------|------------------------------------------------|----------------------|
| 1 | 4 | max(0, 4 + (−2)) = **2** | max(−2, 0−4) = −2 |
| 2 | 1 | max(2, 1 + (−2)) = 2 | max(−2, 0−1) = −1 |

Row `t = 2` (start `maxDiff = dp[1][0] − 2 = −2`):

| d | prices[d] | dp[2][d] = max(dp[2][d−1], prices[d]+maxDiff) | maxDiff after update |
|---|-----------|------------------------------------------------|----------------------|
| 1 | 4 | max(0, 4 + (−2)) = **2** | max(−2, dp[1][1]−4 = −2) = −2 |
| 2 | 1 | max(2, 1 + (−2)) = 2 | max(−2, dp[1][2]−1 = 1) = 1 |

`dp[2][2] = 2` ✓ — a second transaction cannot beat the single 2→4 trade on three days.

---

## Approach 4 — State Machine, Space Optimized (Optimal)

### Intuition
A schedule with at most `k` transactions walks through `2k+1` states:

```
start ──buy#1──▶ bought₁ ──sell#1──▶ sold₁ ──buy#2──▶ bought₂ ──sell#2──▶ sold₂ … soldₖ
```

Track, per state, the **best cash balance** achievable so far: `buy[j]` (right after the j-th buy) and `sell[j]` (right after the j-th sell). Every incoming price relaxes each state in O(k): buying transaction `j` today turns `sell[j-1]` into `sell[j-1] − p`; selling it today turns `buy[j]` into `buy[j] + p`. The day dimension disappears entirely.

Bonus insight that also guards against wastefully large `k`: **one transaction consumes at least two days** (buy strictly before sell), so at most `⌊n/2⌋` transactions can ever complete. If `k >= n/2`, the cap never binds and the problem degenerates to *unlimited* transactions (LeetCode #122), solved greedily by pocketing every day-to-day rise — any single trade decomposes into consecutive daily rises, so harvesting all rises dominates every schedule.

### Algorithm
1. If `n < 2` or `k == 0`, return 0.
2. **Greedy shortcut:** if `k >= n/2`, return `Σ max(0, prices[d] − prices[d−1])`.
3. Initialise `sell[0..k] = 0` and `buy[1..k] = −∞` ("not reachable yet").
4. For every price `p`, for every slot `j = 1..k`:
   `buy[j] = max(buy[j], sell[j-1] − p)`, then `sell[j] = max(sell[j], buy[j] + p)` (same-day update allows a same-day buy+sell, which is a harmless zero-profit no-op).
5. Return `sell[k]` (monotone in `j`, so it dominates using fewer transactions).

### Complexity
- **Time:** O(n·k) — k constant-time relaxations per price; O(n) when the greedy shortcut fires.
- **Space:** O(k) — two arrays of `k+1` running states; no per-day storage at all.

### Code
```go
func stateMachine(k int, prices []int) int {
	n := len(prices)
	if n < 2 || k == 0 {
		return 0 // can't complete any transaction
	}
	// A completed transaction needs a buy day and a later sell day (>= 2 days),
	// so at most n/2 transactions ever fit: for k >= n/2 the cap is irrelevant.
	if k >= n/2 {
		profit := 0
		for d := 1; d < n; d++ {
			if prices[d] > prices[d-1] {
				profit += prices[d] - prices[d-1] // harvest every ascent (unlimited trades)
			}
		}
		return profit
	}
	const negInf = -1 << 60 // effectively -infinity: "state not reachable yet"
	buy := make([]int, k+1)  // buy[j]  = best cash right after the j-th buy
	sell := make([]int, k+1) // sell[j] = best cash right after the j-th sell (sell[0]=0: no trades)
	for j := 1; j <= k; j++ {
		buy[j] = negInf // no buy has happened yet
	}
	for _, p := range prices {
		for j := 1; j <= k; j++ {
			// start transaction j today: cash after (j-1)-th sell minus today's price
			buy[j] = max(buy[j], sell[j-1]-p)
			// close transaction j today: cash after its buy plus today's price
			sell[j] = max(sell[j], buy[j]+p)
		}
	}
	return sell[k] // best cash after at most k completed transactions
}
```

### Dry Run
Example 1: `k = 2, prices = [2,4,1]` — here `n = 3`, so `k = 2 >= n/2 = 1` and the **greedy shortcut** fires:

| d | prices[d−1] → prices[d] | rise? | profit |
|---|--------------------------|-------|--------|
| 1 | 2 → 4 | yes, +2 | 2 |
| 2 | 4 → 1 | no (−3) | 2 |

Return **2** ✓

Example 2 (`k = 2 < n/2 = 3`) exercises the relaxation loop — state after each price:

| p | buy[1] | sell[1] | buy[2] | sell[2] |
|---|--------|---------|--------|---------|
| init | −∞ | 0 | −∞ | 0 |
| 3 | −3 | 0 | −3 | 0 |
| 2 | −2 | 0 | −2 | 0 |
| 6 | −2 | **4** | −2 | 4 |
| 5 | −2 | 4 | −1 | 4 |
| 0 | 0 | 4 | **4** | 4 |
| 3 | 0 | 4 | 4 | **7** |

`sell[2] = 7` ✓ (buy 2 → sell 6, then buy 0 → sell 3).

---

## Key Takeaways

- **The stock-series master recipe:** state = (day, transactions used, holding?). It solves #121/#122/#123/#188/#309/#714 — only the constraint on transitions changes.
- **Count transactions on the sell** — it keeps "at most k" bookkeeping to a single decrement and makes buy legs free to abandon.
- **maxDiff trick:** when a DP transition is `max over m < d of (f(m)) + g(d)`, carry `max f(m)` forward incrementally to drop a factor of n.
- **`k >= n/2` ⇒ unlimited transactions:** a transaction needs two days, so an over-generous cap degenerates the Hard problem into greedy #122 — both an optimization and a defence against `k` blowing up the DP size.
- Initialise unreachable states to −∞, not 0 — otherwise the DP can "sell" shares it never bought.

---

## Related Problems

- LeetCode #121 — Best Time to Buy and Sell Stock (k = 1 special case)
- LeetCode #122 — Best Time to Buy and Sell Stock II (k = ∞; the greedy shortcut)
- LeetCode #123 — Best Time to Buy and Sell Stock III (k = 2 hard-wired; same state machine)
- LeetCode #309 — Best Time to Buy and Sell Stock with Cooldown (extra state edge)
- LeetCode #714 — Best Time to Buy and Sell Stock with Transaction Fee (transition cost)
