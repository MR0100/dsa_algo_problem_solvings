# 0377 — Combination Sum IV

> LeetCode #377 · Difficulty: Medium
> **Categories:** Dynamic Programming, Array

---

## Problem Statement

Given an array of **distinct** integers `nums` and a target integer `target`, return *the number of possible combinations that add up to* `target`.

The test cases are generated so that the answer can fit in a **32-bit** integer.

**Example 1:**

```
Input: nums = [1,2,3], target = 4
Output: 7
Explanation:
The possible combination ways are:
(1, 1, 1, 1)
(1, 1, 2)
(1, 2, 1)
(2, 1, 1)
(1, 3)
(3, 1)
(2, 2)
Note that different sequences are counted as different combinations.
```

**Example 2:**

```
Input: nums = [9], target = 3
Output: 0
```

**Constraints:**

- `1 <= nums.length <= 200`
- `1 <= nums[i] <= 1000`
- All the elements of `nums` are **unique**.
- `1 <= target <= 1000`

**Follow-up:** What if negative numbers are allowed in the given array? How does it change the problem? What limitation do we need to add to the question to allow negative numbers?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Facebook   | ★★☆☆☆ Low        | 2023          |
| Snapchat   | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Dynamic Programming (1D)** — `dp[t]` counts ordered sequences summing to `t`, built from smaller subtotals → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)
- **Backtracking** — the brute-force enumeration explores every ordered pick of a next term → see [`/dsa/backtracking.md`](/dsa/backtracking.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force Backtracking | O(targetᵗᵃʳᵍᵉᵗ) | O(target) | Understanding "ordered" counting only; TLE |
| 2 | DP Top-Down (memoized) | O(target × n) | O(target) | Natural recursive derivation with a cache |
| 3 | DP Bottom-Up (Optimal) | O(target × n) | O(target) | Cleanest iterative solution; amount loop outer |

---

## Approach 1 — Brute Force Backtracking

### Intuition

The problem counts **ordered** sequences: `(1,3)` and `(3,1)` are different. So from any partial sum, the next term may be *any* value in `nums`. Recurse on the remaining amount; a remaining of exactly 0 is one valid ordering, and a negative remaining is pruned.

### Algorithm

1. `count(remaining)`: if `remaining == 0` return 1.
2. For each `num` in `nums` with `num <= remaining`, add `count(remaining - num)`.
3. Return the sum over all choices.

### Complexity

- **Time:** O(target^target) — exponential branching with no reuse.
- **Space:** O(target) — recursion depth.

### Code

```go
func bruteForce(nums []int, target int) int {
	var count func(remaining int) int
	count = func(remaining int) int {
		if remaining == 0 {
			return 1
		}
		total := 0
		for _, num := range nums {
			if num <= remaining {
				total += count(remaining - num)
			}
		}
		return total
	}
	return count(target)
}
```

### Dry Run

Example 1: `nums = [1,2,3], target = 4`.

| Call | Expands into | Returns |
|------|--------------|---------|
| `count(4)` | count(3)+count(2)+count(1) | 4+2+1 = 7 |
| `count(3)` | count(2)+count(1)+count(0) | 2+1+1 = 4 |
| `count(2)` | count(1)+count(0) | 1+1 = 2 |
| `count(1)` | count(0) | 1 |
| `count(0)` | base case | 1 |

`count(4) = 7`. Result: **7** ✔

---

## Approach 2 — DP Top-Down (Memoized Recursion)

### Intuition

The brute force recomputes `count(remaining)` for the same `remaining` repeatedly. That value is fixed, so cache it in an array indexed by the remaining amount.

### Algorithm

1. `memo[r]` caches the answer for subtotal `r` (`-1` = unknown).
2. `count(0) = 1`. Otherwise if cached, return it; else sum `count(remaining - num)` over valid `num`, store into `memo`, return.

### Complexity

- **Time:** O(target × n) — `target` distinct states, each scanning `n` numbers once.
- **Space:** O(target) — memo array plus recursion depth.

### Code

```go
func dpTopDown(nums []int, target int) int {
	memo := make([]int, target+1)
	for i := range memo {
		memo[i] = -1
	}
	var count func(remaining int) int
	count = func(remaining int) int {
		if remaining == 0 {
			return 1
		}
		if memo[remaining] != -1 {
			return memo[remaining]
		}
		total := 0
		for _, num := range nums {
			if num <= remaining {
				total += count(remaining - num)
			}
		}
		memo[remaining] = total
		return total
	}
	return count(target)
}
```

### Dry Run

Example 1: `nums = [1,2,3], target = 4`. First-computation order (deepest first):

| remaining | computed as | memo stored |
|-----------|-------------|-------------|
| 0 | base | 1 |
| 1 | count(0) | 1 |
| 2 | count(1)+count(0) | 2 |
| 3 | count(2)+count(1)+count(0) | 4 |
| 4 | count(3)+count(2)+count(1) | 7 |

Result: **7** ✔ — each subtotal computed exactly once.

---

## Approach 3 — DP Bottom-Up (Optimal)

### Intuition

Let `dp[t]` be the number of ordered sequences summing to `t`. Every such sequence has a **last** term `num`; deleting it leaves a sequence summing to `t - num`. Summing over all possible last terms gives `dp[t] = Σ dp[t-num]`. Crucially the **amount loop is outer** and the **number loop is inner** — that ordering counts permutations (order matters). The reverse loop order would count unordered combinations (the classic Coin Change II).

### Algorithm

1. `dp[0] = 1` (the empty sequence).
2. For `t` from 1 to `target`, for each `num` with `num <= t`: `dp[t] += dp[t-num]`.
3. Return `dp[target]`.

### Complexity

- **Time:** O(target × n).
- **Space:** O(target) — one array.

### Code

```go
func dpBottomUp(nums []int, target int) int {
	dp := make([]int, target+1)
	dp[0] = 1
	for t := 1; t <= target; t++ {
		for _, num := range nums {
			if num <= t {
				dp[t] += dp[t-num]
			}
		}
	}
	return dp[target]
}
```

### Dry Run

Example 1: `nums = [1,2,3], target = 4`.

| t | contributions | dp[t] |
|---|---------------|-------|
| 0 | init | 1 |
| 1 | dp[0] | 1 |
| 2 | dp[1]+dp[0] | 2 |
| 3 | dp[2]+dp[1]+dp[0] | 4 |
| 4 | dp[3]+dp[2]+dp[1] | 7 |

`dp[4] = 7`. Result: **7** ✔

---

## Key Takeaways

- **Order matters here** — this is really a permutation count. The DP loop order encodes that: **amount outer, coin inner** counts ordered sequences; the reverse counts unordered combinations.
- Contrast with Coin Change II (#518): identical values but `dp` loops swapped.
- **Follow-up (negatives):** with negative numbers, infinitely many sequences can sum to the target (e.g. `+1, -1` repeated), so you must cap the sequence length (or forbid cycles) to keep the count finite.
- Recognising "last element removed leaves a smaller subproblem" is the reusable DP framing for counting sums.

---

## Related Problems

- LeetCode #518 — Coin Change II (unordered combinations — swapped loops)
- LeetCode #322 — Coin Change (min coins, same subproblem shape)
- LeetCode #39 — Combination Sum (return the actual combinations)
- LeetCode #70 — Climbing Stairs (Combination Sum IV with nums = [1,2])
