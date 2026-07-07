# 0213 — House Robber II

> LeetCode #213 · Difficulty: Medium
> **Categories:** Dynamic Programming, Array

---

## Problem Statement

You are a professional robber planning to rob houses along a street. Each house has a certain amount of money stashed. All houses at this place are **arranged in a circle.** That means the first house is the neighbor of the last one. Meanwhile, adjacent houses have a security system connected, and **it will automatically contact the police if two adjacent houses were broken into on the same night**.

Given an integer array `nums` representing the amount of money of each house, return *the maximum amount of money you can rob tonight **without alerting the police***.

**Example 1:**
```
Input: nums = [2,3,2]
Output: 3
Explanation: You cannot rob house 1 (money = 2) and then rob house 3 (money = 2),
because they are adjacent houses.
```

**Example 2:**
```
Input: nums = [1,2,3,1]
Output: 4
Explanation: Rob house 1 (money = 1) and then rob house 3 (money = 3).
Total amount you can rob = 1 + 3 = 4.
```

**Example 3:**
```
Input: nums = [1,2,3]
Output: 3
```

**Constraints:**
- `1 <= nums.length <= 100`
- `0 <= nums[i] <= 1000`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2024          |
| Google     | ★★★☆☆ Medium     | 2024          |
| Adobe      | ★★☆☆☆ Low        | 2023          |
| Apple      | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **1D Dynamic Programming** — the House Robber recurrence `dp[i] = max(dp[i-1], dp[i-2] + nums[i])`, applied twice → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)
- **Array Handling** — reducing a circular array to two linear sub-ranges via index windows → see [`/dsa/arrays.md`](/dsa/arrays.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Two DP passes, rolling O(1) space (Optimal) | O(n) | O(1) | The go-to answer |
| 2 | Two DP passes, explicit table | O(n) | O(n) | When you want the recurrence visible for debugging/teaching |

---

## Approach 1 — Two DP Passes, Rolling O(1) Space (Optimal)

### Intuition
The houses form a **circle**, so house `0` and house `n-1` are neighbours and can never both be robbed. That single coupling is the only thing separating this from linear House Robber I. Break the circle by considering two independent linear problems:
- **Case A:** rob among houses `[0 .. n-2]` (exclude the last house).
- **Case B:** rob among houses `[1 .. n-1]` (exclude the first house).

Each case removes one endpoint, so in neither can you pick both ends of the circle simultaneously. Every case is plain House Robber I. The answer is the larger of the two. Handle `n == 1` specially (a lone house has no wrap-around neighbour).

### Algorithm
1. If `n == 1`, return `nums[0]`.
2. Define `robLinear(lo, hi)` = best non-adjacent sum over `nums[lo..hi]` using two rolling variables: `cur = max(prev1, prev2 + nums[i])`.
3. Return `max(robLinear(0, n-2), robLinear(1, n-1))`.

### Complexity
- **Time:** O(n) — two linear scans of at most n elements.
- **Space:** O(1) — two scalars (`prev2`, `prev1`); no DP array.

### Code
```go
func dpArray(nums []int) int {
	n := len(nums)
	if n == 1 {
		return nums[0]
	}
	robLinear := func(lo, hi int) int {
		prev2, prev1 := 0, 0
		for i := lo; i <= hi; i++ {
			cur := max(prev1, prev2+nums[i]) // skip i, or rob i + best two back
			prev2, prev1 = prev1, cur
		}
		return prev1
	}
	return max(robLinear(0, n-2), robLinear(1, n-1))
}
```

### Dry Run (Example 1: `nums = [2,3,2]`, n = 3)

**Case A — houses [0..1] = [2, 3]:**

| i | nums[i] | prev2 | prev1 | cur = max(prev1, prev2+nums[i]) | new prev2 | new prev1 |
|---|---------|-------|-------|---------------------------------|-----------|-----------|
| 0 | 2 | 0 | 0 | max(0, 0+2)=2 | 0 | 2 |
| 1 | 3 | 0 | 2 | max(2, 0+3)=3 | 2 | 3 |

Case A result = **3**.

**Case B — houses [1..2] = [3, 2]:**

| i | nums[i] | prev2 | prev1 | cur | new prev2 | new prev1 |
|---|---------|-------|-------|-----|-----------|-----------|
| 1 | 3 | 0 | 0 | max(0, 0+3)=3 | 0 | 3 |
| 2 | 2 | 0 | 3 | max(3, 0+2)=3 | 3 | 3 |

Case B result = **3**.

Answer = `max(3, 3)` = **3** ✓ (matches expected `3`).

---

## Approach 2 — Two DP Passes, Explicit Table

### Intuition
Identical circular-split idea, but instead of two rolling scalars, fill a `dp` array per case where `dp[k]` = best loot considering the sub-range up to local index `k`. Slightly more memory, but the recurrence is spelled out, which makes the trace easy to follow.

### Algorithm
1. If `n == 1`, return `nums[0]`.
2. `robRange(lo, hi)` over `m = hi-lo+1` houses:
   - `dp[0] = nums[lo]`, `dp[1] = max(nums[lo], nums[lo+1])`.
   - `dp[k] = max(dp[k-1], dp[k-2] + nums[lo+k])` for `k ≥ 2`.
   - return `dp[m-1]`.
3. Return `max(robRange(0, n-2), robRange(1, n-1))`.

### Complexity
- **Time:** O(n) — two linear fills.
- **Space:** O(n) — the `dp` table.

### Code
```go
func dpTable(nums []int) int {
	n := len(nums)
	if n == 1 {
		return nums[0]
	}
	robRange := func(lo, hi int) int {
		m := hi - lo + 1
		if m == 1 {
			return nums[lo]
		}
		dp := make([]int, m)
		dp[0] = nums[lo]
		dp[1] = max(nums[lo], nums[lo+1])
		for k := 2; k < m; k++ {
			dp[k] = max(dp[k-1], dp[k-2]+nums[lo+k])
		}
		return dp[m-1]
	}
	return max(robRange(0, n-2), robRange(1, n-1))
}
```

### Dry Run (Example 1: `nums = [2,3,2]`, n = 3)

**Case A — range [0..1] = [2,3], m = 2:**

| k | local house | dp value |
|---|-------------|----------|
| 0 | nums[0]=2 | dp[0] = 2 |
| 1 | nums[1]=3 | dp[1] = max(2, 3) = 3 |

Case A = dp[1] = **3**.

**Case B — range [1..2] = [3,2], m = 2:**

| k | local house | dp value |
|---|-------------|----------|
| 0 | nums[1]=3 | dp[0] = 3 |
| 1 | nums[2]=2 | dp[1] = max(3, 2) = 3 |

Case B = dp[1] = **3**.

Answer = `max(3, 3)` = **3** ✓

---

## Key Takeaways

- **Break a circular constraint into two linear ones.** "First and last are adjacent" ⇒ solve twice, once excluding each endpoint, and take the max. This "remove one endpoint" trick recurs whenever a cycle couples two boundary elements.
- **Reuse House Robber I unchanged.** The core recurrence `dp[i] = max(dp[i-1], dp[i-2] + nums[i])` is applied to each sub-range; no new DP is invented.
- **Guard `n == 1` explicitly** — the split into `[0..n-2]` and `[1..n-1]` degenerates when there is a single house (the range `[0..-1]` is empty).
- **Rolling variables give O(1) space** because each state depends only on the previous two — a standard 1D-DP space optimisation.

---

## Related Problems

- LeetCode #198 — House Robber (the linear base case)
- LeetCode #337 — House Robber III (tree-shaped houses, DP on trees)
- LeetCode #740 — Delete and Earn (reduces to House Robber)
- LeetCode #256 — Paint House (per-index DP with choices)
- LeetCode #91 — Decode Ways (1D DP with a two-step recurrence)
