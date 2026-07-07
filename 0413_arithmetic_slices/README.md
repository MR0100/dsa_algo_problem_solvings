# 0413 — Arithmetic Slices

> LeetCode #413 · Difficulty: Medium
> **Categories:** Array, Dynamic Programming

---

## Problem Statement

An integer array is called **arithmetic** if it consists of **at least three elements** and if the difference between any two consecutive elements is the same.

- For example, `[1,3,5,7,9]`, `[7,7,7,7]`, and `[3,-1,-5,-9]` are arithmetic sequences.

Given an integer array `nums`, return *the number of arithmetic **subarrays** of `nums`*.

A **subarray** is a contiguous subsequence of the array.

**Example 1:**

```
Input: nums = [1,2,3,4]
Output: 3
Explanation: We have 3 arithmetic slices in nums: [1, 2, 3], [2, 3, 4] and [1,2,3,4] itself.
```

**Example 2:**

```
Input: nums = [1]
Output: 0
```

**Constraints:**

- `1 <= nums.length <= 5000`
- `-1000 <= nums[i] <= 1000`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |
| Facebook   | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **1-D Dynamic Programming** — define `dp[i]` = number of arithmetic slices ending at index `i`; the recurrence `dp[i] = dp[i-1] + 1` (when the difference holds) reuses the previous state, and the answer is `Σ dp[i]` → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)
- **Array / consecutive-difference scan** — the whole problem is about maximal runs of equal adjacent differences in the array; each run of length `k` contributes a triangular count → see [`/dsa/arrays.md`](/dsa/arrays.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Check Every Subarray) | O(n²) | O(1) | Baseline; fine for n ≤ 5000 but wasteful |
| 2 | DP by Ending Index (Bottom-Up) | O(n) | O(n) | Clean statement of the recurrence; easy to explain |
| 3 | O(1) Space Counting (Optimal) | O(n) | O(1) | Best: fold the DP array into one running counter |

---

## Approach 1 — Brute Force (Check Every Subarray)

### Intuition

A subarray `[i..j]` of length ≥ 3 is arithmetic iff **every** consecutive difference inside it equals the first one, `nums[i+1]-nums[i]`. Fix a start `i`, record that reference difference, then extend `j` rightward: as long as the newest difference `nums[j]-nums[j-1]` matches the reference, `[i..j]` is arithmetic and we count it. As soon as a difference breaks, every *longer* slice starting at `i` is also broken, so we stop and advance the start.

### Algorithm

1. `count = 0`.
2. For each start `i` from `0` to `n-3`:
   - `diff = nums[i+1] - nums[i]`.
   - For `j` from `i+2` to `n-1`: if `nums[j]-nums[j-1] == diff`, `count++`; else `break`.
3. Return `count`.

### Complexity

- **Time:** O(n²) — each of the `n` starts may extend across up to `n` elements.
- **Space:** O(1) — only counters.

### Code

```go
func bruteForce(nums []int) int {
	n := len(nums)
	count := 0
	// A start needs at least two more elements to form a length-3 slice.
	for i := 0; i+2 < n; i++ {
		diff := nums[i+1] - nums[i] // the constant difference this run must keep
		for j := i + 2; j < n; j++ {
			if nums[j]-nums[j-1] == diff {
				count++ // [i..j] is arithmetic (length >= 3)
			} else {
				break // difference changed → longer slices from i are impossible
			}
		}
	}
	return count
}
```

### Dry Run

Example 1: `nums = [1,2,3,4]` (all differences are `1`).

| start i | diff | j | nums[j]-nums[j-1] | match? | count after |
|---------|------|---|-------------------|--------|-------------|
| 0 | 1 | 2 | 1 | yes → `[1,2,3]` | 1 |
| 0 | 1 | 3 | 1 | yes → `[1,2,3,4]` | 2 |
| 1 | 1 | 3 | 1 | yes → `[2,3,4]` | 3 |
| 2 | — | — | (i+2 = 4 ≮ 4) start loop ends | — | 3 |

Result: `3` ✔

---

## Approach 2 — DP by Ending Index (Bottom-Up)

### Intuition

Let `dp[i]` be the number of arithmetic slices that **end exactly at index `i`**. Consider the last three elements `nums[i-2], nums[i-1], nums[i]`. If they share one common difference, then two things happen: every arithmetic slice that ended at `i-1` can be stretched one step to end at `i`, **and** a fresh length-3 slice `[i-2..i]` is born. That is precisely `dp[i] = dp[i-1] + 1`. If the difference breaks at `i`, nothing arithmetic ends there, so `dp[i] = 0`. Summing `dp` over all `i` counts every slice exactly once (charged to its right endpoint).

### Algorithm

1. If `n < 3`, return `0`. Allocate `dp` of length `n`; `total = 0`.
2. For `i` from `2` to `n-1`: if `nums[i]-nums[i-1] == nums[i-1]-nums[i-2]`, set `dp[i] = dp[i-1] + 1` and `total += dp[i]`.
3. Return `total`.

### Complexity

- **Time:** O(n) — one pass with O(1) work per index.
- **Space:** O(n) — the `dp` array (removable — see Approach 3).

### Code

```go
func dpBottomUp(nums []int) int {
	n := len(nums)
	if n < 3 {
		return 0 // impossible to have a length-3 slice
	}
	dp := make([]int, n) // dp[i] = number of arithmetic slices ending at i
	total := 0
	for i := 2; i < n; i++ {
		// Same difference across the last three elements?
		if nums[i]-nums[i-1] == nums[i-1]-nums[i-2] {
			dp[i] = dp[i-1] + 1 // extend all slices ending at i-1, plus the new triple
			total += dp[i]      // accumulate into the running answer
		}
		// else dp[i] stays 0 (the zero value): no slice ends here
	}
	return total
}
```

### Dry Run

Example 1: `nums = [1,2,3,4]`.

| i | nums[i]-nums[i-1] | nums[i-1]-nums[i-2] | equal? | dp[i] = dp[i-1]+1 | total |
|---|-------------------|---------------------|--------|-------------------|-------|
| 2 | 3-2 = 1 | 2-1 = 1 | yes | dp[2] = 0+1 = 1 | 1 |
| 3 | 4-3 = 1 | 3-2 = 1 | yes | dp[3] = 1+1 = 2 | 3 |

Result: `total = 3` ✔ (dp = `[0,0,1,2]`).

---

## Approach 3 — O(1) Space Counting (Optimal)

### Intuition

`dp[i]` reads only `dp[i-1]`, so a single scalar `cur` replaces the whole array. `cur` is "how many arithmetic slices end at the current index": increment it while the adjacent difference stays constant, reset it to `0` when the difference breaks, and add it into `total` at every step. Viewed globally, a maximal **run** of `k` equal consecutive differences (i.e. `k+1` numbers in arithmetic progression) contributes `1 + 2 + … + (k-1) = (k-1)k/2` slices — and the incrementing counter computes exactly that triangular number on the fly.

### Algorithm

1. `total = 0`, `cur = 0`.
2. For `i` from `2` to `n-1`: if the difference holds, `cur++` then `total += cur`; else `cur = 0`.
3. Return `total`.

### Complexity

- **Time:** O(n) — one pass.
- **Space:** O(1) — two integer counters, no DP array.

### Code

```go
func countingOptimal(nums []int) int {
	total, cur := 0, 0 // cur = arithmetic slices ending at the current index
	for i := 2; i < len(nums); i++ {
		if nums[i]-nums[i-1] == nums[i-1]-nums[i-2] {
			cur++        // one more slice ends here than ended at i-1 (+ the new triple)
			total += cur // triangular accumulation: 1,2,3,... across a run
		} else {
			cur = 0 // difference broke → no slice ends at i; restart the run
		}
	}
	return total
}
```

### Dry Run

Example 1: `nums = [1,2,3,4]`.

| i | difference holds? | cur after | total after |
|---|-------------------|-----------|-------------|
| 2 | 1 == 1 → yes | 1 | 1 |
| 3 | 1 == 1 → yes | 2 | 3 |

Result: `3` ✔. (For `[1,3,5,7,9]` the run has `k=4` equal diffs, `cur` climbs `1,2,3` → `total = 6 = (4-1)·4/2`.)

---

## Key Takeaways

- **Charge each slice to its right endpoint.** Defining `dp[i]` = slices *ending at* `i` gives the tidy recurrence `dp[i] = dp[i-1] + 1` and avoids double counting — a reusable DP framing for "count all subarrays with property P".
- **Extending an arithmetic run adds `cur+1` new slices,** not one: the new element lengthens every slice that ended at the previous index and creates one new triple. This is why the count grows triangularly, not linearly.
- **A run of `k` equal adjacent differences yields `(k-1)k/2` slices.** Recognising the closed form is a fast sanity check for any answer.
- **DP → O(1) space** whenever the recurrence looks back only one step; keep just the previous value in a scalar.

---

## Related Problems

- LeetCode #446 — Arithmetic Slices II — Subsequence (non-contiguous, hashmap DP)
- LeetCode #1027 — Longest Arithmetic Subsequence (DP over differences)
- LeetCode #53 — Maximum Subarray (run-based 1-D DP, Kadane)
- LeetCode #1218 — Longest Arithmetic Subsequence of Given Difference
