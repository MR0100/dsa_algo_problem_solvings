# 0334 — Increasing Triplet Subsequence

> LeetCode #334 · Difficulty: Medium
> **Categories:** Array, Greedy

---

## Problem Statement

Given an integer array `nums`, return `true` if there exists a triple of indices `(i, j, k)` such that `i < j < k` and `nums[i] < nums[j] < nums[k]`. If no such indices exists, return `false`.

**Example 1:**

```
Input: nums = [1,2,3,4,5]
Output: true
Explanation: Any triplet where i < j < k is valid.
```

**Example 2:**

```
Input: nums = [5,4,3,2,1]
Output: false
Explanation: No triplet exists.
```

**Example 3:**

```
Input: nums = [2,1,5,0,4,6]
Output: true
Explanation: The triplet (3, 4, 5) is valid because nums[3] == 0 < nums[4] == 4 < nums[5] == 6.
```

**Constraints:**

- `1 <= nums.length <= 5 * 10^5`
- `-2^31 <= nums[i] <= 2^31 - 1`

**Follow-up:** Could you implement a solution that runs in `O(n)` time complexity and `O(1)` space complexity?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2024          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |
| Apple      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Greedy** — the optimal solution greedily keeps the smallest possible "start" and "middle" tails, which is enough to detect a triple without tracking indices → see [`/dsa/greedy.md`](/dsa/greedy.md)
- **Prefix / Suffix Aggregates** — Approach 2 precomputes running min-from-left and max-from-right so each middle index is testable in O(1) → see [`/dsa/prefix_sum.md`](/dsa/prefix_sum.md)
- **Longest Increasing Subsequence intuition** — this is the LIS problem specialized to length 3; the two-smallest trick is the patience-sorting idea reduced to a two-element tails array → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Three Nested Loops) | O(n³) | O(1) | Tiny inputs / sanity baseline only |
| 2 | Precomputed Min-Left / Max-Right | O(n) | O(n) | Clear O(n) that generalizes to "find a valid middle" |
| 3 | Two Smallest (Greedy, Optimal) | O(n) | O(1) | The follow-up answer: single pass, constant space |

---

## Approach 1 — Brute Force (Three Nested Loops)

### Intuition

The problem asks whether *any* increasing triple of indices exists, so enumerate all triples `i < j < k` and test `nums[i] < nums[j] < nums[k]`. Correct by construction; the two early-`continue`s prune impossible middles/starts but the worst case is still cubic.

### Algorithm

1. For each `i`, for each `j > i` with `nums[j] > nums[i]`, for each `k > j`: if `nums[k] > nums[j]` return `true`.
2. If nothing qualifies, return `false`.

### Complexity

- **Time:** O(n³) — three nested loops over the array.
- **Space:** O(1) — only loop indices.

### Code

```go
func bruteForce(nums []int) bool {
	n := len(nums)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if nums[j] <= nums[i] {
				continue // need nums[i] < nums[j]
			}
			for k := j + 1; k < n; k++ {
				if nums[k] > nums[j] {
					return true // found nums[i] < nums[j] < nums[k]
				}
			}
		}
	}
	return false
}
```

### Dry Run

Example 1: `nums = [1,2,3,4,5]`.

| i | nums[i] | j | nums[j] > nums[i]? | k | nums[k] > nums[j]? | result |
|---|---------|---|--------------------|---|--------------------|--------|
| 0 | 1 | 1 | 2 > 1 yes | 2 | 3 > 2 yes | **return true** |

Found `(0,1,2)` → `1 < 2 < 3`. Result: `true` ✔

---

## Approach 2 — Precomputed Min-Left / Max-Right

### Intuition

A triple exists exactly when some index `j` can act as the **middle**: there is a strictly smaller value somewhere to its left and a strictly larger value somewhere to its right. Precompute, for every position, the minimum of everything to the left (`prefixMin`) and the maximum of everything to the right (`suffixMax`). Then `j` is a valid middle iff `prefixMin[j-1] < nums[j] < suffixMax[j+1]`.

### Algorithm

1. Build `prefixMin[i] = min(nums[0..i])` left-to-right.
2. Build `suffixMax[i] = max(nums[i..n-1])` right-to-left.
3. For each middle `j` in `[1, n-2]`: if `prefixMin[j-1] < nums[j] < suffixMax[j+1]`, return `true`.
4. Otherwise `false`.

### Complexity

- **Time:** O(n) — three linear passes.
- **Space:** O(n) — the two auxiliary arrays (does not meet the O(1) follow-up).

### Code

```go
func minLeftMaxRight(nums []int) bool {
	n := len(nums)
	if n < 3 {
		return false // need at least three elements
	}
	prefixMin := make([]int, n) // prefixMin[i] = min(nums[0..i])
	prefixMin[0] = nums[0]
	for i := 1; i < n; i++ {
		prefixMin[i] = min(prefixMin[i-1], nums[i]) // extend the running min
	}
	suffixMax := make([]int, n) // suffixMax[i] = max(nums[i..n-1])
	suffixMax[n-1] = nums[n-1]
	for i := n - 2; i >= 0; i-- {
		suffixMax[i] = max(suffixMax[i+1], nums[i]) // extend the running max
	}
	for j := 1; j < n-1; j++ {
		if prefixMin[j-1] < nums[j] && nums[j] < suffixMax[j+1] {
			return true
		}
	}
	return false
}
```

### Dry Run

Example 3: `nums = [2,1,5,0,4,6]`.

`prefixMin = [2,1,1,0,0,0]`, `suffixMax = [6,6,6,6,6,6]`.

| j | nums[j] | prefixMin[j-1] | suffixMax[j+1] | prefixMin < nums[j] < suffixMax? |
|---|---------|----------------|----------------|----------------------------------|
| 1 | 1 | 2 | 6 | 2 < 1? no |
| 2 | 5 | 1 | 6 | 1 < 5 < 6? **yes → true** |

Middle `j=2` (value 5) has 1 before it and 6 after it. Result: `true` ✔

---

## Approach 3 — Two Smallest (Greedy, Optimal)

### Intuition

Keep two running values: `first` = the smallest number seen so far, and `second` = the smallest number that has *some strictly smaller number before it* (a valid "middle" tail). Scan left to right:

- If `x <= first`, `x` is a new, smaller potential start → update `first`.
- Else if `x <= second`, `x` is a better middle → update `second`. (Crucially, `second` being set means a smaller `first` occurred *before* it.)
- Else `x > second`, and since `second` was preceded by some smaller `first`, we have `first < second < x` in index order → a triple exists.

Reassigning `first` to a later, smaller value does **not** corrupt the invariant: it only remembers a smaller start that appeared before the *current* `second`; any future `x > second` still had a smaller-than-`second` element before `second`. This is exactly the length-3 case of LIS via patience sorting.

### Algorithm

1. `first = second = +∞`.
2. For each `x`: if `x <= first` set `first = x`; else if `x <= second` set `second = x`; else return `true`.
3. Return `false`.

### Complexity

- **Time:** O(n) — a single pass.
- **Space:** O(1) — two scalars (meets the follow-up).

### Code

```go
func twoSmallest(nums []int) bool {
	first, second := 1<<62, 1<<62 // smallest and second-smallest valid tails
	for _, x := range nums {
		switch {
		case x <= first:
			first = x // new smallest candidate for the triple's start
		case x <= second:
			second = x // x can serve as a middle (some smaller `first` precedes it)
		default:
			return true // x beats both → increasing triple exists
		}
	}
	return false
}
```

### Dry Run

Example 3: `nums = [2,1,5,0,4,6]`, start `first = second = +∞`.

| x | x ≤ first? | x ≤ second? | Action | first | second |
|---|-----------|-------------|--------|-------|--------|
| 2 | yes | — | first = 2 | 2 | +∞ |
| 1 | yes | — | first = 1 | 1 | +∞ |
| 5 | no | yes | second = 5 | 1 | 5 |
| 0 | yes | — | first = 0 | 0 | 5 |
| 4 | no | yes | second = 4 | 0 | 4 |
| 6 | no | no | **return true** | 0 | 4 |

At `x=6`: `6 > second=4`, and `second=4` was set after some smaller `first` → triple `(1<4<6` from values, indices `1,4,5`). Result: `true` ✔

> Note the `first = 0` update at index 3 happens *after* `second = 5`; it looks like it breaks ordering, but `second` still points to a middle preceded by an earlier `first=1`. The check `x > second` remains a valid witness. That is the classic subtlety of this trick.

---

## Key Takeaways

- **This is LIS with target length 3.** The two-variable greedy is patience sorting's tails array truncated to size 2 — generalize to `k` smallest for "increasing k-tuple".
- **`first` can be overwritten by a later smaller value without breaking correctness**: `second` already encodes "a smaller element came before me", so a witness `x > second` is always valid regardless of what `first` currently holds. Understanding *why* this is safe is the whole interview.
- **When you need only existence, not the actual subsequence,** you can drop index bookkeeping and track just the minimal tails — O(1) space.
- The prefix-min / suffix-max approach is the more *explainable* O(n) and generalizes to "is there a valid middle" style problems, at the cost of O(n) space.

---

## Related Problems

- LeetCode #300 — Longest Increasing Subsequence (the general version; patience sorting)
- LeetCode #128 — Longest Consecutive Sequence (existence over an array via clever bookkeeping)
- LeetCode #674 — Longest Continuous Increasing Subsequence (contiguous variant)
- LeetCode #673 — Number of Longest Increasing Subsequence (counting LIS)
- LeetCode #42 — Trapping Rain Water (same prefix-max / suffix-max precompute idea)
