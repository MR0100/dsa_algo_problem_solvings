# 0485 — Max Consecutive Ones

> LeetCode #485 · Difficulty: Easy
> **Categories:** Array, Sliding Window

---

## Problem Statement

Given a binary array `nums`, return *the maximum number of consecutive* `1`*'s in the array*.

**Example 1:**

```
Input: nums = [1,1,0,1,1,1]
Output: 3
Explanation: The first two digits or the last three digits are consecutive 1s. The maximum number of consecutive 1s is 3.
```

**Example 2:**

```
Input: nums = [1,0,1,1,0,1]
Output: 2
```

**Constraints:**

- `1 <= nums.length <= 10^5`
- `nums[i]` is either `0` or `1`.

**Follow up:** If the input is a stream, i.e. you can only read the array one element at a time, how would you solve it? *(The running-count and sliding-window approaches below both work unchanged on a stream — they only ever look at the current element.)*

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Facebook   | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Array single-pass scan** — the answer is a property of contiguous elements, computable by a single left-to-right sweep with two counters → see [`/dsa/arrays.md`](/dsa/arrays.md)
- **Sliding Window** — framing the streak as a window that never contains a `0` gives the exact template that generalises to the "allow up to k zeros" follow-ups (#487 / #1004) → see [`/dsa/sliding_window.md`](/dsa/sliding_window.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Re-scan From Each Start) | O(n²) | O(1) | Baseline; illustrates the redundant re-scanning the optimal avoids |
| 2 | Single-Pass Running Count (Optimal) | O(n) | O(1) | The go-to answer; simplest correct one-pass solution |
| 3 | Sliding Window (Framing for the Follow-Up) | O(n) | O(1) | Same cost, but the reusable template for "at most k zeros" variants |

---

## Approach 1 — Brute Force (Re-scan From Each Start)

### Intuition

The answer is literally "the longest stretch of consecutive 1's", so the most direct method is: from each index `i`, if `nums[i]` is a 1, walk forward counting 1's until you hit a 0 or the end, and keep the largest count. It re-reads elements many times but requires no insight — a useful baseline that shows exactly what work the optimal pass eliminates.

### Algorithm

1. `best = 0`.
2. For each start `i`: set `run = 0`; extend `j` from `i` while `nums[j] == 1`, incrementing `run`. Update `best = max(best, run)`.
3. Return `best`.

### Complexity

- **Time:** O(n²) — an all-1's array makes every start scan to the end.
- **Space:** O(1) — a couple of counters.

### Code

```go
func bruteForce(nums []int) int {
	best := 0
	for i := 0; i < len(nums); i++ {
		run := 0
		for j := i; j < len(nums) && nums[j] == 1; j++ { // extend the run from i
			run++
		}
		if run > best {
			best = run // remember the longest run seen so far
		}
	}
	return best
}
```

### Dry Run

Example 1: `nums = [1,1,0,1,1,1]`.

| start i | nums[i] | run extends over | run | best after |
|---------|---------|-------------------|-----|------------|
| 0 | 1 | idx 0,1 (then 0 at 2) | 2 | 2 |
| 1 | 1 | idx 1 (then 0 at 2) | 1 | 2 |
| 2 | 0 | none | 0 | 2 |
| 3 | 1 | idx 3,4,5 (end) | 3 | 3 |
| 4 | 1 | idx 4,5 | 2 | 3 |
| 5 | 1 | idx 5 | 1 | 3 |

Result: `3` ✔

---

## Approach 2 — Single-Pass Running Count (Optimal)

### Intuition

Consecutive 1's are contiguous, so one left-to-right pass is enough. Keep `cur`, the length of the streak of 1's ending at the current index: each 1 extends it (`cur++`), each 0 breaks it (`cur = 0`). The answer is the maximum value `cur` ever attains. Every element is examined exactly once, and it works verbatim on a stream (only the current element is inspected).

### Algorithm

1. `best = 0`, `cur = 0`.
2. For each `x` in `nums`: if `x == 1`, `cur++` and `best = max(best, cur)`; else `cur = 0`.
3. Return `best`.

### Complexity

- **Time:** O(n) — a single pass.
- **Space:** O(1) — two counters.

### Code

```go
func runningCount(nums []int) int {
	best, cur := 0, 0
	for _, x := range nums {
		if x == 1 {
			cur++ // extend the current streak of 1's
			if cur > best {
				best = cur // new longest streak
			}
		} else {
			cur = 0 // a 0 breaks the streak; start counting fresh
		}
	}
	return best
}
```

### Dry Run

Example 1: `nums = [1,1,0,1,1,1]`.

| idx | x | cur (after) | best (after) |
|-----|---|-------------|--------------|
| 0 | 1 | 1 | 1 |
| 1 | 1 | 2 | 2 |
| 2 | 0 | 0 | 2 |
| 3 | 1 | 1 | 2 |
| 4 | 1 | 2 | 2 |
| 5 | 1 | 3 | 3 |

Result: `3` ✔

---

## Approach 3 — Sliding Window (Framing for the Follow-Up)

### Intuition

Frame the streak as a window `[left, right]` that never contains a 0. `right` advances through the array; when `nums[right]` is 0, the window can hold no 1's ending here, so snap `left` to `right+1` (empty it). Otherwise `[left..right]` is a valid all-1's block whose length `right-left+1` is a candidate answer. This is deliberately more machinery than the base problem needs — but it is the exact template that generalises to **Max Consecutive Ones II/III** (allow up to `k` zeros): there you keep the window valid by advancing `left` only while it contains more than `k` zeros, instead of resetting on every 0.

### Algorithm

1. `left = 0`, `best = 0`.
2. For `right = 0 … n-1`: if `nums[right] == 0`, set `left = right+1` (reset the window); else `best = max(best, right-left+1)`.
3. Return `best`.

### Complexity

- **Time:** O(n) — `left` and `right` each move forward monotonically.
- **Space:** O(1).

### Code

```go
func slidingWindow(nums []int) int {
	left, best := 0, 0
	for right := 0; right < len(nums); right++ {
		if nums[right] == 0 {
			left = right + 1 // a 0 can't be in an all-1's window → restart after it
		} else if right-left+1 > best {
			best = right - left + 1 // current all-1's window is the longest yet
		}
	}
	return best
}
```

### Dry Run

Example 1: `nums = [1,1,0,1,1,1]`.

| right | nums[right] | left (after) | window len (right-left+1) | best (after) |
|-------|-------------|--------------|----------------------------|--------------|
| 0 | 1 | 0 | 1 | 1 |
| 1 | 1 | 0 | 2 | 2 |
| 2 | 0 | 3 | — (reset) | 2 |
| 3 | 1 | 3 | 1 | 2 |
| 4 | 1 | 3 | 2 | 2 |
| 5 | 1 | 3 | 3 | 3 |

Result: `3` ✔

---

## Key Takeaways

- **Contiguous-property problems yield to a single pass.** When the answer depends only on runs of adjacent elements, one sweep with a running counter beats any re-scanning approach.
- **Reset-on-break** is the pattern: extend a counter on the "good" element, zero it on the "bad" one, track the max.
- **Sliding window is the generalisation.** The base problem resets the window on every 0; the follow-ups (#487 allow one flip, #1004 allow `k` flips) only change the *invalidation rule* — shrink `left` while the window holds more than `k` zeros. Learn the window frame here and the harder versions are a one-line change.
- All three work on a **stream** except the brute force — they inspect only the current element.

---

## Related Problems

- LeetCode #487 — Max Consecutive Ones II (flip at most one 0)
- LeetCode #1004 — Max Consecutive Ones III (flip at most k 0's)
- LeetCode #424 — Longest Repeating Character Replacement (window with a ≤k-changes budget)
- LeetCode #3 — Longest Substring Without Repeating Characters (canonical sliding window)
