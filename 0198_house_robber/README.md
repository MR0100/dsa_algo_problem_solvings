# 0198 — House Robber

> LeetCode #198 · Difficulty: Medium
> **Categories:** Dynamic Programming, Array

---

## Problem Statement

You are a professional robber planning to rob houses along a street. Each house has a certain amount of money stashed, the only constraint stopping you from robbing each of them is that adjacent houses have security systems connected and **it will automatically contact the police if two adjacent houses were broken into on the same night**.

Given an integer array `nums` representing the amount of money of each house, return *the maximum amount of money you can rob tonight **without alerting the police***.

**Example 1:**

```
Input: nums = [1,2,3,1]
Output: 4
Explanation: Rob house 1 (money = 1) and then rob house 3 (money = 3).
Total amount you can rob = 1 + 3 = 4.
```

**Example 2:**

```
Input: nums = [2,7,9,3,1]
Output: 12
Explanation: Rob house 1 (money = 2), rob house 3 (money = 9) and rob house 5 (money = 1).
Total amount you can rob = 2 + 9 + 1 = 12.
```

**Constraints:**

- `1 <= nums.length <= 100`
- `0 <= nums[i] <= 400`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★★ Very High  | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★★☆ High       | 2024          |
| Apple      | ★★★☆☆ Medium     | 2023          |
| Adobe      | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **1D Dynamic Programming** — the answer for the whole street is built from a single-index state `dp[i]` = "best loot considering only the first `i` houses", where each state depends on just the two states before it (`dp[i-1]` and `dp[i-2]`). This is the canonical "linear DP with a small fixed lookback" pattern → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)
- **Optimal Substructure & Overlapping Subproblems** — the best plan from house `i` onward reuses the best plan from `i+1` and `i+2`; those subproblems recur, which is exactly what memoization / tabulation exploits.
- **Space Optimization via Rolling Variables** — because the recurrence never reaches further back than `dp[i-2]`, the whole table collapses into two integers, the same trick that reduces Fibonacci to `O(1)` space.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Recursion) | O(2ⁿ) | O(n) | Baseline only; shows the rob/skip decision tree — TLE for `n` beyond ~30 |
| 2 | DP Top-Down (Memoization) | O(n) | O(n) | Natural first optimization; keeps the recursive shape but caches subproblems |
| 3 | DP Bottom-Up (Tabulation) | O(n) | O(n) | Iterative, no recursion stack; clearest way to reason about the table |
| 4 | Space-Optimized DP (Optimal) | O(n) | O(1) | The interview answer; same time, constant memory via two rolling variables |

---

## Approach 1 — Brute Force

### Intuition

Standing in front of house `i` there are exactly two legal futures: **rob** house `i` and jump to house `i+2` (the alarm forbids the neighbour `i+1`), or **skip** house `i` and move on to house `i+1`. The best loot from house `i` onward is simply the better of those two futures. Recursing on both branches enumerates every valid subset of non-adjacent houses without ever having to construct them explicitly.

### Algorithm

1. Define `robFrom(i)` = best loot obtainable from houses `i..n-1`.
2. Base case: `i >= n` → return `0` (no houses left to rob).
3. Recurrence: `robFrom(i) = max(nums[i] + robFrom(i+2), robFrom(i+1))`.
4. The answer is `robFrom(0)`.

### Complexity

- **Time:** O(2ⁿ) — each call spawns two recursive calls and identical subproblems are recomputed exponentially often (the recursion tree is Fibonacci-shaped).
- **Space:** O(n) — the maximum recursion stack depth is the number of houses.

### Code

```go
func bruteForce(nums []int) int {
	var robFrom func(i int) int
	robFrom = func(i int) int {
		if i >= len(nums) {
			return 0 // ran past the last house: nothing more to steal
		}
		take := nums[i] + robFrom(i+2) // rob house i → house i+1 is off-limits
		skip := robFrom(i + 1)         // leave house i → free to consider i+1
		return max(take, skip)
	}
	return robFrom(0)
}
```

### Dry Run

Example 1: `nums = [1, 2, 3, 1]`, `n = 4`. Tracing `robFrom(0)`; each node returns `max(take, skip)`.

| Call | take = nums[i] + robFrom(i+2) | skip = robFrom(i+1) | returns |
|------|-------------------------------|---------------------|---------|
| robFrom(4) | — | — | 0 (base case) |
| robFrom(3) | nums[3] + robFrom(5) = 1 + 0 = 1 | robFrom(4) = 0 | max(1, 0) = **1** |
| robFrom(2) | nums[2] + robFrom(4) = 3 + 0 = 3 | robFrom(3) = 1 | max(3, 1) = **3** |
| robFrom(1) | nums[1] + robFrom(3) = 2 + 1 = 3 | robFrom(2) = 3 | max(3, 3) = **3** |
| robFrom(0) | nums[0] + robFrom(2) = 1 + 3 = 4 | robFrom(1) = 3 | max(4, 3) = **4** |

Result: `4` ✔ — corresponds to robbing house 0 (money 1) then house 2 (money 3).

---

## Approach 2 — DP Top-Down (Memoization)

### Intuition

The brute force recomputes `robFrom(i)` for the same `i` over and over, yet there are only `n` distinct subproblems — one per starting house. Caching each answer the first time it is computed collapses the exponential recursion tree into a linear number of real calls; every repeat visit is now an `O(1)` table lookup.

### Algorithm

1. Create `memo[i]` to cache the answer for houses `i..n-1`; initialise every entry to `-1` as a "not computed" sentinel (safe because loot is never negative).
2. On each call, if `memo[i] != -1` return it immediately.
3. Otherwise compute `max(nums[i] + robFrom(i+2), robFrom(i+1))`, store it in `memo[i]`, and return it.
4. The answer is `robFrom(0)`.

### Complexity

- **Time:** O(n) — there are `n` distinct subproblems and each is computed exactly once in `O(1)`.
- **Space:** O(n) — the memo table plus the recursion stack.

### Code

```go
func dpTopDown(nums []int) int {
	memo := make([]int, len(nums))
	for i := range memo {
		memo[i] = -1 // -1 = not yet computed (valid answers are always >= 0)
	}
	var robFrom func(i int) int
	robFrom = func(i int) int {
		if i >= len(nums) {
			return 0
		}
		if memo[i] != -1 {
			return memo[i] // already solved: reuse instead of recomputing
		}
		memo[i] = max(nums[i]+robFrom(i+2), robFrom(i+1))
		return memo[i]
	}
	return robFrom(0)
}
```

### Dry Run

Example 1: `nums = [1, 2, 3, 1]`. Calls resolve deepest-first; each result is written to `memo` once and reused.

| Call | computes | memo after |
|------|----------|------------|
| robFrom(3) | max(nums[3]+robFrom(5), robFrom(4)) = max(1+0, 0) = 1 | `[-1, -1, -1, 1]` |
| robFrom(2) | max(nums[2]+robFrom(4), robFrom(3)) = max(3+0, 1) = 3 | `[-1, -1, 3, 1]` |
| robFrom(1) | max(nums[1]+robFrom(3), robFrom(2)) = max(2+1, 3) = 3 (robFrom(3)=memo hit) | `[-1, 3, 3, 1]` |
| robFrom(0) | max(nums[0]+robFrom(2), robFrom(1)) = max(1+3, 3) = 4 (robFrom(2)=memo hit) | `[4, 3, 3, 1]` |

Result: `robFrom(0) = 4` ✔ — the two memo hits are the subproblems the brute force would have recomputed.

---

## Approach 3 — DP Bottom-Up (Tabulation)

### Intuition

Flip the recursion around. Instead of "best loot from house `i` to the end", compute "best loot among the first `i` houses" and build it up from the smallest prefix. Each new house offers the same two choices — **rob it** (its cash plus the best of the first `i-2` houses) or **skip it** (carry over the best of the first `i-1` houses). No recursion, no stack: just fill a table left to right.

### Algorithm

1. Let `dp[i]` = max loot using only the first `i` houses. Size the table `n+1`.
2. Base cases: `dp[0] = 0` (no houses) and `dp[1] = nums[0]` (one house → rob it; constraints guarantee `n >= 1`).
3. For `i = 2..n`: `dp[i] = max(dp[i-1], dp[i-2] + nums[i-1])` (note the `i-1` index into `nums` because `dp` is 1-based over houses).
4. The answer is `dp[n]`.

### Complexity

- **Time:** O(n) — a single pass over the houses.
- **Space:** O(n) — the `dp` table of `n+1` entries.

### Code

```go
func dpBottomUp(nums []int) int {
	n := len(nums)
	dp := make([]int, n+1) // dp[i] = max loot using only the first i houses
	dp[0] = 0              // zero houses → zero loot
	dp[1] = nums[0]        // one house → rob it (constraints guarantee n >= 1)
	for i := 2; i <= n; i++ {
		skip := dp[i-1]             // leave house i-1 (0-indexed) alone
		take := dp[i-2] + nums[i-1] // rob it: its neighbour must be skipped
		dp[i] = max(skip, take)
	}
	return dp[n]
}
```

### Dry Run

Example 1: `nums = [1, 2, 3, 1]`, `n = 4`. Seed `dp[0] = 0`, `dp[1] = nums[0] = 1`.

| i | nums[i-1] | skip = dp[i-1] | take = dp[i-2] + nums[i-1] | dp[i] = max |
|---|-----------|----------------|----------------------------|-------------|
| 2 | 2 | dp[1] = 1 | dp[0] + 2 = 0 + 2 = 2 | max(1, 2) = **2** |
| 3 | 3 | dp[2] = 2 | dp[1] + 3 = 1 + 3 = 4 | max(2, 4) = **4** |
| 4 | 1 | dp[3] = 4 | dp[2] + 1 = 2 + 1 = 3 | max(4, 3) = **4** |

Final table: `dp = [0, 1, 2, 4, 4]`. Answer `dp[4] = 4` ✔

---

## Approach 4 — Space-Optimized DP (Optimal)

### Intuition

The tabulation recurrence `dp[i] = max(dp[i-1], dp[i-2] + nums[i-1])` only ever touches the two previous entries. Everything older is dead weight. So the entire table collapses into a rolling pair `(prev2, prev1)` = `(dp[i-2], dp[i-1])` that slides forward one house at a time — the same window trick that reduces Fibonacci to constant space.

### Algorithm

1. Start `prev2 = 0` and `prev1 = 0` (best loot two houses back / one house back, before any house is seen).
2. For each house value `v`: compute `curr = max(prev1, prev2 + v)` — skip this house (`prev1`) or rob it (`prev2 + v`).
3. Slide the window: `prev2 = prev1`, `prev1 = curr`.
4. After the loop, `prev1` holds the answer.

### Complexity

- **Time:** O(n) — a single pass over the houses.
- **Space:** O(1) — exactly two integers regardless of input size.

### Code

```go
func spaceOptimized(nums []int) int {
	prev2, prev1 := 0, 0 // dp[i-2] and dp[i-1] of the tabulation
	for _, v := range nums {
		curr := max(prev1, prev2+v) // skip this house (prev1) or rob it (prev2+v)
		prev2, prev1 = prev1, curr  // slide the two-value window forward
	}
	return prev1
}
```

### Dry Run

Example 1: `nums = [1, 2, 3, 1]`. Start `prev2 = 0`, `prev1 = 0`.

| v | curr = max(prev1, prev2 + v) | prev2 after | prev1 after |
|---|------------------------------|-------------|-------------|
| 1 | max(0, 0 + 1) = 1 | 0 | 1 |
| 2 | max(1, 0 + 2) = 2 | 1 | 2 |
| 3 | max(2, 1 + 3) = 4 | 2 | 4 |
| 1 | max(4, 2 + 1) = 4 | 4 | 4 |

Loop ends. Answer `prev1 = 4` ✔ — matches the tabulation's final `dp[n]`, using only two variables.

---

## Key Takeaways

- **"Pick non-adjacent elements for max sum" is the House Robber signature.** The moment a problem forbids using two neighbours together, reach for `dp[i] = max(dp[i-1], dp[i-2] + value[i])`.
- **The rob/skip decision generalizes:** at each element you either take it (and forfeit the previous one) or leave it (and inherit the previous best). That single recurrence powers all four approaches — only the *bookkeeping* changes.
- **Memoization and tabulation are the same DP** viewed top-down vs. bottom-up. Top-down keeps the recursive intuition; bottom-up removes the stack and makes the fixed lookback obvious.
- **When a recurrence only looks back a constant distance, drop the table.** Here `dp[i]` needs only `dp[i-1]` and `dp[i-2]`, so two rolling variables give `O(1)` space with no loss of clarity — the standard interview-grade answer.
- **Base cases matter:** seeding `dp[0] = 0` and `dp[1] = nums[0]` (or `prev2 = prev1 = 0`) is what makes single-house and two-house inputs correct without special-casing.

---

## Related Problems

- LeetCode #213 — House Robber II (houses arranged in a **circle**; run the linear DP twice, excluding either the first or the last house)
- LeetCode #337 — House Robber III (houses form a **binary tree**; DP on the tree returning rob/not-rob pairs)
- LeetCode #740 — Delete and Earn (reduces to House Robber after bucketing values)
- LeetCode #70 — Climbing Stairs (same two-step Fibonacci-shaped recurrence)
- LeetCode #746 — Min Cost Climbing Stairs (linear DP with a min instead of max)
- LeetCode #198 pattern also underlies #2320 — Count Number of Ways to Place Houses
