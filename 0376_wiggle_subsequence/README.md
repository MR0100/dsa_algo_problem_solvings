# 0376 — Wiggle Subsequence

> LeetCode #376 · Difficulty: Medium
> **Categories:** Dynamic Programming, Greedy, Array

---

## Problem Statement

A **wiggle sequence** is a sequence where the differences between successive numbers strictly alternate between positive and negative. The first difference (if one exists) may be either positive or negative. A sequence with one element and a sequence with two non-equal elements are trivially wiggle sequences.

- For example, `[1, 7, 4, 9, 2, 5]` is a **wiggle sequence** because the differences `(6, -3, 5, -7, 3)` alternate between positive and negative.
- In contrast, `[1, 4, 7, 2, 5]` and `[1, 7, 4, 5, 5]` are not wiggle sequences. The first is not because its first two differences are positive, and the second is not because its last difference is zero.

A **subsequence** is obtained by deleting some elements (possibly zero) from the original sequence, leaving the remaining elements in their original order.

Given an integer array `nums`, return *the length of the longest **wiggle subsequence** of* `nums`.

**Example 1:**

```
Input: nums = [1,7,4,9,2,5]
Output: 6
Explanation: The entire sequence is a wiggle sequence with differences (6, -3, 5, -7, 3).
```

**Example 2:**

```
Input: nums = [1,17,5,10,13,15,10,5,16,8]
Output: 7
Explanation: There are several subsequences that achieve this length.
One is [1, 17, 10, 13, 10, 16, 8] with differences (16, -7, 3, -3, 6, -8).
```

**Example 3:**

```
Input: nums = [1,2,3,4,5,6,7,8,9]
Output: 2
```

**Constraints:**

- `1 <= nums.length <= 1000`
- `0 <= nums[i] <= 1000`

**Follow-up:** Can you solve this in `O(n)` time?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Facebook   | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Dynamic Programming (1D)** — track two running lengths (last step up / last step down) and extend them element by element → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)
- **Greedy** — the optimal wiggle subsequence is exactly the set of peaks and valleys; counting direction flips is a local, greedy decision that provably yields the global optimum → see [`/dsa/greedy.md`](/dsa/greedy.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force Recursion | O(2ⁿ) | O(n) | Understanding the up/down state only; TLE for n>~20 |
| 2 | DP up/down Tables | O(n²) | O(n) | Clear formulation; fine for n ≤ 1000 |
| 3 | Linear DP (rolling) | O(n) | O(1) | Collapses the tables into two scalars |
| 4 | Greedy Slope Counting (Optimal) | O(n) | O(1) | The follow-up answer; count direction flips |

---

## Approach 1 — Brute Force Recursion

### Intuition

From any committed element, the remaining wiggle length depends on a single piece of state: whether the next difference must go **up** or **down**. Recurse over that state, trying to extend to every later element that points the required way, and take the best. The first move is free to be up or down.

### Algorithm

1. Define `calc(i, isUp)` = best number of *additional* elements we can chain after committing index `i`, where the next kept difference must be positive iff `isUp`.
2. For each `j > i`: if `isUp` and `nums[j] > nums[i]`, recurse `1 + calc(j, false)`; if `!isUp` and `nums[j] < nums[i]`, recurse `1 + calc(j, true)`.
3. Answer = `1 + max(calc(0, true), calc(0, false))` (the `1` counts element 0 itself).

### Complexity

- **Time:** O(2ⁿ) — each element independently branches into keep/skip for both directions.
- **Space:** O(n) — recursion stack depth.

### Code

```go
func bruteForce(nums []int) int {
	n := len(nums)
	if n == 0 {
		return 0
	}
	var calc func(i int, isUp bool) int
	calc = func(i int, isUp bool) int {
		best := 0
		for j := i + 1; j < len(nums); j++ {
			if isUp && nums[j] > nums[i] {
				best = max(best, 1+calc(j, false))
			} else if !isUp && nums[j] < nums[i] {
				best = max(best, 1+calc(j, true))
			}
		}
		return best
	}
	return 1 + max(calc(0, true), calc(0, false))
}
```

### Dry Run

Example 1: `nums = [1,7,4,9,2,5]`, starting at index 0 (value 1).

| Call | State | Extends to | Contribution |
|------|-------|-----------|--------------|
| `calc(0, true)` | need up from 1 | j=1 (7), j=2 (4), j=3 (9), j=5 (5) | best path 1→7→4→9→2→5 gives +5 |
| inside: `calc(1,false)` | need down from 7 | j=2 (4), j=4 (2) | 7→4 then up... |
| … deepest chain | 1<7>4<9>2<5 | full alternation | 5 extra elements |

`1 + calc(0, true) = 1 + 5 = 6`. Result: **6** ✔

---

## Approach 2 — Dynamic Programming (up/down tables)

### Intuition

Define two tables. `up[i]` is the longest wiggle subsequence ending at `i` whose **last** difference is positive; `down[i]` the same for a negative last difference. To land at `i` with an up step you must have arrived at some earlier `j` (with `nums[j] < nums[i]`) whose last step was **down** — the directions must alternate.

### Algorithm

1. Initialise `up[i] = down[i] = 1` for all `i` (a single element).
2. For each `i`, for each `j < i`:
   - if `nums[i] > nums[j]`: `up[i] = max(up[i], down[j] + 1)`.
   - if `nums[i] < nums[j]`: `down[i] = max(down[i], up[j] + 1)`.
3. Answer = `max(up[n-1], down[n-1])`.

### Complexity

- **Time:** O(n²) — nested loop over earlier indices.
- **Space:** O(n) — two length-n tables.

### Code

```go
func dpTables(nums []int) int {
	n := len(nums)
	if n == 0 {
		return 0
	}
	up := make([]int, n)
	down := make([]int, n)
	for i := range nums {
		up[i], down[i] = 1, 1
	}
	for i := 1; i < n; i++ {
		for j := 0; j < i; j++ {
			if nums[i] > nums[j] {
				up[i] = max(up[i], down[j]+1)
			} else if nums[i] < nums[j] {
				down[i] = max(down[i], up[j]+1)
			}
		}
	}
	return max(up[n-1], down[n-1])
}
```

### Dry Run

Example 1: `nums = [1,7,4,9,2,5]`.

| i | nums[i] | up[i] | down[i] |
|---|---------|-------|---------|
| 0 | 1 | 1 | 1 |
| 1 | 7 | 2 (7>1, down[0]+1) | 1 |
| 2 | 4 | 2 | 3 (4<7, up[1]+1) |
| 3 | 9 | 4 (9>4, down[2]+1) | 1 |
| 4 | 2 | 2 | 5 (2<9, up[3]+1) |
| 5 | 5 | 6 (5>2, down[4]+1) | 3 |

`max(up[5], down[5]) = max(6, 3) = 6`. Result: **6** ✔

---

## Approach 3 — Linear DP (rolling up/down)

### Intuition

Scanning the tables, notice that when `nums[i] > nums[i-1]` the only useful earlier state is the best down-ending length so far, and it only improves as we move right. So we never need the full arrays — two rolling scalars `up` and `down` suffice, updated from the *previous element* alone.

### Algorithm

1. `up, down = 1, 1`.
2. For `i` from 1: if `nums[i] > nums[i-1]` set `up = down + 1`; else if `nums[i] < nums[i-1]` set `down = up + 1`; equal → unchanged.
3. Return `max(up, down)`.

### Complexity

- **Time:** O(n) — one pass.
- **Space:** O(1) — two scalars.

### Code

```go
func dpLinear(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	up, down := 1, 1
	for i := 1; i < len(nums); i++ {
		if nums[i] > nums[i-1] {
			up = down + 1
		} else if nums[i] < nums[i-1] {
			down = up + 1
		}
	}
	return max(up, down)
}
```

### Dry Run

Example 1: `nums = [1,7,4,9,2,5]`.

| i | nums[i-1]→nums[i] | direction | up | down |
|---|-------------------|-----------|----|------|
| — | — | init | 1 | 1 |
| 1 | 1→7 | up | 2 | 1 |
| 2 | 7→4 | down | 2 | 3 |
| 3 | 4→9 | up | 4 | 3 |
| 4 | 9→2 | down | 4 | 5 |
| 5 | 2→5 | up | 6 | 5 |

`max(6, 5) = 6`. Result: **6** ✔

---

## Approach 4 — Greedy Slope Counting (Optimal)

### Intuition

The longest wiggle subsequence is precisely the sequence of **turning points** (local peaks and valleys) of the array. Walk through consecutive differences and count each time the sign flips from `+` to `−` or `−` to `+`. Flat steps (`diff == 0`) never contribute. This is the direct answer to the O(n) follow-up.

### Algorithm

1. If fewer than 2 elements, return the length.
2. `count = 1`, `prevDiff = 0`.
3. For `i` from 1: `diff = nums[i] - nums[i-1]`. If (`diff > 0 && prevDiff <= 0`) or (`diff < 0 && prevDiff >= 0`): `count++`, `prevDiff = diff`.
4. Return `count`.

### Complexity

- **Time:** O(n) — single pass.
- **Space:** O(1).

### Code

```go
func greedy(nums []int) int {
	if len(nums) < 2 {
		return len(nums)
	}
	count := 1
	prevDiff := 0
	for i := 1; i < len(nums); i++ {
		diff := nums[i] - nums[i-1]
		if (diff > 0 && prevDiff <= 0) || (diff < 0 && prevDiff >= 0) {
			count++
			prevDiff = diff
		}
	}
	return count
}
```

### Dry Run

Example 1: `nums = [1,7,4,9,2,5]`.

| i | diff | prevDiff | flip? | count | prevDiff after |
|---|------|----------|-------|-------|----------------|
| — | — | 0 | — | 1 | 0 |
| 1 | +6 | 0 | yes (+ vs ≤0) | 2 | +6 |
| 2 | −3 | +6 | yes (− vs ≥0) | 3 | −3 |
| 3 | +5 | −3 | yes | 4 | +5 |
| 4 | −7 | +5 | yes | 5 | −7 |
| 5 | +3 | −7 | yes | 6 | +3 |

Result: **6** ✔

---

## Key Takeaways

- **Wiggle length = number of peaks and valleys.** The optimal subsequence keeps exactly the turning points; interior monotone runs are collapsed to their endpoints.
- Carrying a **direction bit** (last step up vs down) is the key state that turns an exponential search into linear DP.
- The `up`/`down` rolling recurrence is a reusable trick for "alternating" sequence problems (compare with LIS variants).
- Equal adjacent elements are neutral — handle `diff == 0` explicitly so it neither counts nor resets the direction.

---

## Related Problems

- LeetCode #300 — Longest Increasing Subsequence (monotone version of the DP)
- LeetCode #53 — Maximum Subarray (running-scalar DP)
- LeetCode #135 — Candy (two-direction greedy sweep)
- LeetCode #128 — Longest Consecutive Sequence (sequence-structure greedy)
