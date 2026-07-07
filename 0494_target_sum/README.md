# 0494 — Target Sum

> LeetCode #494 · Difficulty: Medium
> **Categories:** Array, Dynamic Programming, Backtracking

---

## Problem Statement

You are given an integer array `nums` and an integer `target`.

You want to build an **expression** out of nums by adding one of the symbols `'+'` and `'-'` before each integer in nums and then concatenate all the integers.

- For example, if `nums = [2, 1]`, you can add a `'+'` before `2` and a `'-'` before `1` and concatenate them to build the expression `"+2-1"`.

Return the number of different **expressions** that you can build, which evaluates to `target`.

**Example 1:**

```
Input: nums = [1,1,1,1,1], target = 3
Output: 5
Explanation: There are 5 ways to assign symbols to make the sum of nums be target 3.
-1 + 1 + 1 + 1 + 1 = 3
+1 - 1 + 1 + 1 + 1 = 3
+1 + 1 - 1 + 1 + 1 = 3
+1 + 1 + 1 - 1 + 1 = 3
+1 + 1 + 1 + 1 - 1 = 3
```

**Example 2:**

```
Input: nums = [1], target = 1
Output: 1
```

**Constraints:**

- `1 <= nums.length <= 20`
- `0 <= nums[i] <= 1000`
- `0 <= sum(nums[i]) <= 1000`
- `-1000 <= target <= 1000`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Facebook   | ★★★★☆ High       | 2024          |
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Backtracking (binary sign choice)** — the brute force explores a `2ⁿ` decision tree, choosing `+` or `−` at each index → see [`/dsa/backtracking.md`](/dsa/backtracking.md)
- **1-D Dynamic Programming (subset-sum / 0-1 knapsack count)** — the algebraic reduction `P = (total+target)/2` turns the task into "count subsets that sum to `P`", solved with a downward-iterated 1-D DP array → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force DFS | O(2ⁿ) | O(n) | Fine for n ≤ 20; shows the raw decision tree |
| 2 | Top-Down DP (memo on `(i, sum)`) | O(n · total) | O(n · total) | Keeps the recursion but caches overlapping states |
| 3 | Subset-Sum 1-D DP (Optimal) | O(n · s) | O(s) | The intended answer; `s = (total+target)/2` |

---

## Approach 1 — Brute Force DFS (Try Both Signs)

### Intuition

Each of the `n` positions independently gets a `+` or `−`, so there are `2ⁿ` expressions. Walk that binary tree depth-first, carrying the running signed sum. At a leaf (all numbers consumed) count it iff the sum equals `target`. No memoization, so identical `(index, sum)` states are recomputed across branches.

### Algorithm

1. `dfs(i, sum)`: if `i == n`, return `1` when `sum == target`, else `0`.
2. Otherwise return `dfs(i+1, sum + nums[i]) + dfs(i+1, sum − nums[i])`.
3. Answer is `dfs(0, 0)`.

### Complexity

- **Time:** O(2ⁿ) — two branches per index; `n ≤ 20` → ~`10⁶` leaves, acceptable.
- **Space:** O(n) — recursion depth.

### Code

```go
func bruteForceDFS(nums []int, target int) int {
	var dfs func(i, sum int) int
	dfs = func(i, sum int) int {
		if i == len(nums) {
			if sum == target {
				return 1 // one valid signed expression reached target
			}
			return 0
		}
		// Branch on the sign of nums[i]: '+' then '-'.
		return dfs(i+1, sum+nums[i]) + dfs(i+1, sum-nums[i])
	}
	return dfs(0, 0)
}
```

### Dry Run

Example 2 (small enough to enumerate the whole tree): `nums = [1]`, `target = 1`.

| Call | Action | Return |
|------|--------|--------|
| dfs(0, 0) | branch `+1` and `−1` | sum of children |
| dfs(1, 0+1=1) | i==n, sum 1 == target → 1 | 1 |
| dfs(1, 0−1=−1) | i==n, sum −1 ≠ target → 0 | 0 |
| dfs(0, 0) total | 1 + 0 | **1** ✔ |

For Example 1 (`[1,1,1,1,1], target=3`) the same tree has `2⁵ = 32` leaves, exactly `5` of which sum to `3` → output `5`.

---

## Approach 2 — Top-Down DP with Memoization

### Intuition

The only state that determines the remaining count is `(i, running sum)` — how we arrived there is irrelevant. Many sign prefixes converge on the same `(i, sum)`, so cache each state's result. Running sums lie in `[−total, +total]`, giving at most `n · (2·total + 1)` distinct states; storing them per index (as a map to stay sparse and dodge negative indexing) turns the exponential tree polynomial.

### Algorithm

1. Compute `total`; if `|target| > total`, return `0` (unreachable magnitude).
2. `memo[i]` maps `sum → ways`.
3. `dfs(i, sum)`: base case as in Approach 1; else look up `memo[i][sum]`, and if absent compute `dfs(i+1, sum+nums[i]) + dfs(i+1, sum−nums[i])`, store it, return it.

### Complexity

- **Time:** O(n · total) — number of distinct `(i, sum)` states, each solved once with O(1) work.
- **Space:** O(n · total) — the memo tables, plus O(n) recursion stack.

### Code

```go
func dpTopDown(nums []int, target int) int {
	total := 0
	for _, v := range nums {
		total += v // maximum absolute reachable sum
	}
	// If target is unreachable in magnitude, no assignment can work.
	if target > total || target < -total {
		return 0
	}
	// memo[i] maps a running sum -> ways; use a map to keep it sparse and to
	// sidestep negative indexing cleanly.
	memo := make([]map[int]int, len(nums))
	for i := range memo {
		memo[i] = map[int]int{}
	}
	var dfs func(i, sum int) int
	dfs = func(i, sum int) int {
		if i == len(nums) {
			if sum == target {
				return 1
			}
			return 0
		}
		if v, ok := memo[i][sum]; ok {
			return v // state already solved
		}
		ways := dfs(i+1, sum+nums[i]) + dfs(i+1, sum-nums[i])
		memo[i][sum] = ways // cache before returning
		return ways
	}
	return dfs(0, 0)
}
```

### Dry Run

Example 1: `nums = [1,1,1,1,1]`, `target = 3`, `total = 5`. States are `(i, sum)`. A slice of the memo after the recursion (values are "ways to finish from here and hit 3"):

| State (i, sum) | Meaning | Ways |
|----------------|---------|------|
| (5, 3) | all used, sum 3 == target | 1 |
| (5, x≠3) | all used, wrong sum | 0 |
| (4, 2) | one `1` left; `+1`→(5,3)=1, `−1`→(5,1)=0 | 1 |
| (4, 4) | `+1`→(5,5)=0, `−1`→(5,3)=1 | 1 |
| (3, 1) | two left; reaches 3 via net `+2` | 1 |
| … | reuse of cached (4,·) states avoids recompute | … |
| (0, 0) | root | **5** |

The `5` distinct sign assignments summing to `3` are counted with each `(i, sum)` computed once. Output `5` ✔

---

## Approach 3 — Subset-Sum 0/1 Knapsack, 1-D DP (Optimal)

### Intuition

Partition `nums` into the `+` group (sum `P`) and the `−` group (sum `N`). Then `P − N = target` and `P + N = total`, so adding the equations: `P = (total + target) / 2`. The answer is simply **the number of subsets of `nums` that sum to `P`** — the counting variant of subset-sum. Build `dp[j] = number of ways to reach sum j`; process each number and iterate `j` **downward** so each item is used at most once (the hallmark of 0/1 knapsack in one array). If `total + target` is odd or `|target| > total`, `P` is not a valid integer subset sum → answer `0`.

### Algorithm

1. Compute `total`. If `|target| > total` **or** `(total + target)` is odd, return `0`.
2. `s = (total + target) / 2`; allocate `dp[0..s]` with `dp[0] = 1` (empty subset).
3. For each `num`: for `j` from `s` down to `num`: `dp[j] += dp[j − num]`.
4. Return `dp[s]`.

### Complexity

- **Time:** O(n · s) — `n` numbers × capacity `s ≤ total ≤ 1000`.
- **Space:** O(s) — one 1-D DP array.

### Code

```go
func subsetSumDP(nums []int, target int) int {
	total := 0
	for _, v := range nums {
		total += v
	}
	// Need P = (total+target)/2 to be a non-negative integer.
	if target > total || target < -total || (total+target)%2 != 0 {
		return 0
	}
	s := (total + target) / 2 // required '+'-group subset sum
	dp := make([]int, s+1)    // dp[j] = # of subsets summing to j
	dp[0] = 1                 // empty subset makes sum 0 in exactly one way
	for _, num := range nums {
		// Iterate downward so num contributes to each j at most once (0/1).
		for j := s; j >= num; j-- {
			dp[j] += dp[j-num] // ways to reach j using this num
		}
	}
	return dp[s]
}
```

### Dry Run

Example 1: `nums = [1,1,1,1,1]`, `target = 3`, `total = 5`. `s = (5+3)/2 = 4`. Start `dp = [1,0,0,0,0]` (indices 0..4). Each row processes one `1`, iterating `j` from 4 down to 1 (`dp[j] += dp[j-1]`):

| After processing | dp[0] | dp[1] | dp[2] | dp[3] | dp[4] |
|------------------|-------|-------|-------|-------|-------|
| init             | 1 | 0 | 0 | 0 | 0 |
| 1st `1`          | 1 | 1 | 0 | 0 | 0 |
| 2nd `1`          | 1 | 2 | 1 | 0 | 0 |
| 3rd `1`          | 1 | 3 | 3 | 1 | 0 |
| 4th `1`          | 1 | 4 | 6 | 4 | 1 |
| 5th `1`          | 1 | 5 | 10 | 10 | 5 |

`dp[s] = dp[4] = 5` ✔ — exactly the `C(5,4) = 5` ways to choose which four `1`s are positive (the fifth negative gives net `4 − 1 = 3`).

---

## Key Takeaways

- **Sign-assignment → subset-sum via algebra:** `P − N = target`, `P + N = total` ⇒ `P = (total+target)/2`. Recognizing this collapses a `2ⁿ` search into an O(n·sum) DP. The same split appears in "partition into two equal halves" (#416).
- **Feasibility gates first:** if `(total+target)` is odd or `|target| > total`, the answer is `0` — no subset can hit a non-integer or out-of-range `P`.
- **0/1 knapsack in 1-D iterates the capacity downward.** Iterating upward would let a single item be reused (that is the *unbounded* knapsack). This down-vs-up direction is the crux to memorize.
- **Zeros matter for counting:** a `0` in `nums` can be signed `+` or `−` without changing the sum, so it doubles the number of expressions. The subset-sum count handles this automatically — `num = 0` makes `dp[j] += dp[j]`, doubling every reachable count.

---

## Related Problems

- LeetCode #416 — Partition Equal Subset Sum (subset-sum feasibility, same 1-D DP)
- LeetCode #698 — Partition to K Equal Sum Subsets (subset partition, backtracking)
- LeetCode #1049 — Last Stone Weight II (min `|P−N|`, same P/N split)
- LeetCode #322 — Coin Change (knapsack-style DP, unbounded variant)
- LeetCode #518 — Coin Change II (count subsets/combinations DP)
