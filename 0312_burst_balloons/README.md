# 0312 — Burst Balloons

> LeetCode #312 · Difficulty: Hard
> **Categories:** Array, Dynamic Programming, Interval DP

---

## Problem Statement

You are given `n` balloons, indexed from `0` to `n - 1`. Each balloon is painted with a number on it represented by an array `nums`. You are asked to burst all the balloons.

If you burst the `i`th balloon, you will get `nums[i - 1] * nums[i] * nums[i + 1]` coins. If `i - 1` or `i + 1` goes out of bounds of the array, then treat it as if there is a balloon with a `1` painted on it.

Return the maximum coins you can collect by bursting the balloons wisely.

**Example 1:**

```
Input: nums = [3,1,5,8]
Output: 167
Explanation:
nums = [3,1,5,8] --> [3,5,8] --> [3,8] --> [8] --> []
coins =  3*1*5    +   3*5*8   +  1*3*8  + 1*8*1 = 167
```

**Example 2:**

```
Input: nums = [1,5]
Output: 10
```

**Constraints:**

- `n == nums.length`
- `1 <= n <= 300`
- `0 <= nums[i] <= 100`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★★☆☆ Medium     | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Interval Dynamic Programming** — solving over ranges `(l, r)` where the "last" event fixes the boundaries → see [`/dsa/dynamic_programming_2d.md`](/dsa/dynamic_programming_2d.md)
- **Divide and conquer / memoisation** — decomposing an interval into independent subintervals around a split point → see [`/dsa/divide_and_conquer.md`](/dsa/divide_and_conquer.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (all orders) | O(n!) | O(n) | Tiny `n`; illustrates the naive decision |
| 2 | Interval DP Top-Down (memo) | O(n³) | O(n²) | Clean recursion; natural "last balloon" framing |
| 3 | Interval DP Bottom-Up (Optimal) | O(n³) | O(n²) | Production; no recursion overhead |

---

## Approach 1 — Brute Force

### Intuition
The naive framing is a sequence of choices: which balloon to burst **next**. Popping balloon `i` earns `left*nums[i]*right` for its **current** surviving neighbours. Try all choices, recurse on the shrunken array, take the best. This explores every ordering — `n!` of them.

### Algorithm
1. Pad the array with `1` on both ends so boundary balloons have neighbours.
2. Keep `remaining` = padded indices of alive balloons.
3. For each alive balloon, its neighbours are the alive balloons immediately before/after it in `remaining`; burst it, add its coins, recurse on the rest.

### Complexity
- **Time:** O(n!) — every ordering of bursts is tried.
- **Space:** O(n) recursion depth plus per-call copies.

### Code
```go
func bruteForce(nums []int) int {
	padded := make([]int, len(nums)+2)
	padded[0], padded[len(padded)-1] = 1, 1
	copy(padded[1:], nums)

	remaining := make([]int, len(nums))
	for i := range remaining {
		remaining[i] = i + 1
	}

	var solve func(rem []int) int
	solve = func(rem []int) int {
		if len(rem) == 0 {
			return 0
		}
		best := 0
		for idx := 0; idx < len(rem); idx++ {
			left := 1
			if idx > 0 {
				left = padded[rem[idx-1]]
			} else {
				left = padded[0]
			}
			right := 1
			if idx < len(rem)-1 {
				right = padded[rem[idx+1]]
			} else {
				right = padded[len(padded)-1]
			}
			coins := left * padded[rem[idx]] * right

			next := make([]int, 0, len(rem)-1)
			next = append(next, rem[:idx]...)
			next = append(next, rem[idx+1:]...)

			if got := coins + solve(next); got > best {
				best = got
			}
		}
		return best
	}
	return solve(remaining)
}
```

### Dry Run
`nums = [3,1,5,8]`, padded = `[1,3,1,5,8,1]`, remaining = `[1,2,3,4]` (padded indices).

Only the optimal branch is shown; the recursion actually explores all 4! = 24 orders.

| step | burst (value) | neighbours | coins | remaining values |
|------|---------------|-----------|-------|-------------------|
| 1 | index 2 (=1) | 3, 5 | 3·1·5 = 15 | [3,5,8] |
| 2 | index 1 (=5) | 3, 8 | 3·5·8 = 120 | [3,8] |
| 3 | index 0 (=3) | 1, 8 | 1·3·8 = 24 | [8] |
| 4 | index 0 (=8) | 1, 1 | 1·8·1 = 8 | [] |

Total = 15 + 120 + 24 + 8 = **167** (the max over all orders).

---

## Approach 2 — Interval DP Top-Down (Memoised)

### Intuition
Reframe from "burst first" to "burst **last**". If balloon `k` is the last to pop inside open interval `(l, r)`, then when it bursts both boundaries `l` and `r` are still present, earning `nums[l]*nums[k]*nums[r]`. Everything in `(l,k)` and `(k,r)` was already burst independently, so:

```
dp(l,r) = max over k in (l,r) of  dp(l,k) + nums[l]*nums[k]*nums[r] + dp(k,r)
```

Choosing "last" makes the boundaries **fixed**, decoupling the two subproblems — "first" would leave them changing.

### Algorithm
1. Pad with `1` on both ends → valid balloons at indices `1..n`.
2. `dp(l,r)` = best coins burstable strictly between `l` and `r`.
3. Base case `r - l <= 1` → 0. Memoise in a 2-D table.

### Complexity
- **Time:** O(n³) — O(n²) intervals, each scanning O(n) split points.
- **Space:** O(n²) memo + O(n) recursion depth.

### Code
```go
func dpTopDown(nums []int) int {
	n := len(nums)
	padded := make([]int, n+2)
	padded[0], padded[n+1] = 1, 1
	copy(padded[1:], nums)

	memo := make([][]int, n+2)
	for i := range memo {
		memo[i] = make([]int, n+2)
		for j := range memo[i] {
			memo[i][j] = -1
		}
	}

	var dp func(l, r int) int
	dp = func(l, r int) int {
		if r-l <= 1 {
			return 0
		}
		if memo[l][r] != -1 {
			return memo[l][r]
		}
		best := 0
		for k := l + 1; k < r; k++ {
			coins := padded[l]*padded[k]*padded[r] + dp(l, k) + dp(k, r)
			if coins > best {
				best = coins
			}
		}
		memo[l][r] = best
		return best
	}
	return dp(0, n+1)
}
```

### Dry Run
`nums = [3,1,5,8]`, padded = `[1,3,1,5,8,1]` (indices 0..5). Compute `dp(0,5)`.

| call | k choices tried | best combination | value |
|------|-----------------|------------------|-------|
| dp(1,3) between 3 and 5 | k=2 (val 1) | 3·1·5 | 15 |
| dp(3,5) between 5 and 1 | k=4 (val 8) | 5·8·1 | 40 |
| dp(1,5) between 3 and 1 | k=4 last: dp(1,4)+3·8·1+dp(4,5) | ... | 159 |
| dp(0,5) between 1 and 1 | k=4 last: dp(0,4)+1·8·1+dp(4,5) | 1·8·1 + dp(0,4) | **167** |

The winning split at `dp(0,5)` is `k=4` (value 8 last): `dp(0,4) = 159`, plus `1·8·1 = 8`, total 167.

---

## Approach 3 — Interval DP Bottom-Up (Optimal)

### Intuition
`dp[l][r]` depends only on strictly smaller intervals `dp[l][k]` and `dp[k][r]`. Process intervals from shortest to longest so every dependency is already computed, removing recursion overhead.

### Algorithm
1. Pad with `1` on both ends.
2. For `length` from 2 to `n+1` (distance `r-l`), for each left `l` with `r = l+length`, set `dp[l][r] = max over k in (l,r) of padded[l]*padded[k]*padded[r] + dp[l][k] + dp[k][r]`.
3. Answer is `dp[0][n+1]`.

### Complexity
- **Time:** O(n³) — intervals × split points.
- **Space:** O(n²) — the DP table.

### Code
```go
func dpBottomUp(nums []int) int {
	n := len(nums)
	padded := make([]int, n+2)
	padded[0], padded[n+1] = 1, 1
	copy(padded[1:], nums)

	dp := make([][]int, n+2)
	for i := range dp {
		dp[i] = make([]int, n+2)
	}

	for length := 2; length <= n+1; length++ {
		for l := 0; l+length <= n+1; l++ {
			r := l + length
			for k := l + 1; k < r; k++ {
				coins := padded[l]*padded[k]*padded[r] + dp[l][k] + dp[k][r]
				if coins > dp[l][r] {
					dp[l][r] = coins
				}
			}
		}
	}
	return dp[0][n+1]
}
```

### Dry Run
`nums = [3,1,5,8]`, padded = `[1,3,1,5,8,1]`. Fill by increasing `length`.

| length | interval (l,r) | boundaries | best k (value) | dp[l][r] |
|--------|----------------|-----------|----------------|----------|
| 2 | (0,2) | 1,1 | k=1 (3): 1·3·1 | 3 |
| 2 | (1,3) | 3,5 | k=2 (1): 3·1·5 | 15 |
| 2 | (2,4) | 1,8 | k=3 (5): 1·5·8 | 40 |
| 2 | (3,5) | 5,1 | k=4 (8): 5·8·1 | 40 |
| 3 | (1,4) | 3,8 | k=3: 3·5·8 + dp(1,3) | 135 |
| ... | ... | ... | ... | ... |
| 5 | (0,5) | 1,1 | k=4 (8) last: dp(0,4)+1·8·1 | **167** |

`dp[0][5] = 167`.

---

## Key Takeaways
- **"Think last, not first."** When neighbours change as you remove elements, fix the LAST removed element in an interval — its neighbours are the interval boundaries, which decouples the subproblems. This is the signature move of interval DP.
- **Padding with sentinels** (`1` on each side) removes special-casing for out-of-bounds neighbours.
- Interval DP is filled by **increasing interval length** (bottom-up) or via memoised recursion (top-down); both are O(n³) here.

---

## Related Problems
- LeetCode #1000 — Minimum Cost to Merge Stones (interval DP)
- LeetCode #1039 — Minimum Score Triangulation of Polygon (interval DP, "last" triangle)
- LeetCode #516 — Longest Palindromic Subsequence (interval DP)
- LeetCode #486 — Predict the Winner (interval DP, game)
