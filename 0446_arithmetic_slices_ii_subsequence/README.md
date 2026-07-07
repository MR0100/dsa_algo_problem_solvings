# 0446 — Arithmetic Slices II - Subsequence

> LeetCode #446 · Difficulty: Hard
> **Categories:** Array, Dynamic Programming, Hash Table

---

## Problem Statement

Given an integer array `nums`, return *the number of all the **arithmetic subsequences** of* `nums`.

A sequence of numbers is called arithmetic if it consists of **at least three elements** and if the difference between any two consecutive elements is the same.

- For example, `[1, 3, 5, 7, 9]`, `[7, 7, 7, 7]`, and `[3, -1, -5, -9]` are arithmetic sequences.
- For example, `[1, 1, 2, 5, 7]` is not an arithmetic sequence.

A **subsequence** of an array is a sequence that can be formed by removing some elements (possibly none) of the array.

- For example, `[2,5,10]` is a subsequence of `[1,2,1,2,4,1,5,10]`.

The test cases are generated so that the answer fits in **32-bit** integer.

**Example 1:**

```
Input: nums = [2,4,6,8,10]
Output: 7
Explanation: All arithmetic subsequence slices are:
[2,4,6]
[4,6,8]
[6,8,10]
[2,4,6,8]
[4,6,8,10]
[2,4,6,8,10]
[2,6,10]
```

**Example 2:**

```
Input: nums = [7,7,7,7,7]
Output: 16
```

**Constraints:**

- `1 <= nums.length <= 1000`
- `-2^31 <= nums[i] <= 2^31 - 1`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |
| Adobe      | ★★☆☆☆ Low        | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Dynamic Programming (2-D state)** — the answer is built from `dp[i][d]` = number of arithmetic subsequences ending at index `i` with common difference `d`; each ordered pair `(j, i)` promotes weak (length-2) runs at `j` into strong (length ≥ 3) subsequences at `i` → see [`/dsa/dynamic_programming_2d.md`](/dsa/dynamic_programming_2d.md)
- **Hash Map** — because differences span the full 32-bit range they cannot index an array, so each `dp[i]` is a `map[difference]count`, giving O(1) lookup/update per pair → see [`/dsa/hash_map.md`](/dsa/hash_map.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Recursive Enumeration) | O(2^n) | O(n) | Understanding only; explodes on all-equal / long arithmetic arrays |
| 2 | DP with Hash Maps (Optimal) | O(n²) | O(n²) | The intended solution; n ≤ 1000 makes n² fine |

---

## Approach 1 — Brute Force (Recursive Enumeration)

### Intuition

An arithmetic subsequence is pinned down by *its last index* and *its common difference `d`*. Fix the first two chosen elements as an ordered pair `(i, j)`; that locks `d = nums[j] - nums[i]`. Now recursively try to append any later index `k` whose value equals `nums[j] + d`. The moment we append a **third** element we have a valid slice (length ≥ 3), so we count it — and we keep recursing to grow it into length 4, 5, … Doing this for every seed pair enumerates every arithmetic subsequence exactly once. It is correct but exponential: an all-equal array like `[7,7,7,7,7]` lets *every* subset of size ≥ 3 qualify.

### Algorithm

1. For every ordered starting pair `(i, j)` with `i < j`, compute `d = nums[j] - nums[i]`.
2. Call `extend(j, d)`: scan `k` from `j+1`; whenever `nums[k] - nums[j] == d`, that append is the ≥ 3rd element → `count++`, then recurse `extend(k, d)`.
3. Return the accumulated `count`.

### Complexity

- **Time:** O(2^n) — in the worst case (all equal, or one long arithmetic run) every subset of size ≥ 3 is generated; the recursion tree branches into essentially every subsequence.
- **Space:** O(n) — recursion depth is at most the array length; no extra tables.

### Code

```go
func bruteForce(nums []int) int {
	n := len(nums)
	count := 0 // total valid arithmetic subsequences (length >= 3)

	// extend continues an arithmetic subsequence whose last element sits at
	// index last and whose common difference is d. Any successful append here
	// is the 3rd (or later) element, so it forms a valid slice.
	var extend func(last int, d int64)
	extend = func(last int, d int64) {
		for k := last + 1; k < n; k++ {
			// Does appending nums[k] preserve the common difference d?
			if int64(nums[k])-int64(nums[last]) == d {
				count++      // length >= 3 reached → a valid subsequence
				extend(k, d) // try to grow it further from k
			}
		}
	}

	// Every pair (i, j) seeds a difference; recursion supplies the >=3rd term.
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			d := int64(nums[j]) - int64(nums[i]) // fixed common difference
			extend(j, d)                         // count all extensions
		}
	}
	return count
}
```

### Dry Run

Example 1: `nums = [2,4,6,8,10]` (indices 0..4). Only the diff `d = 2` seeds (from adjacent evens) produce chains; the seed pair `(0,4)` gives `d = 8` and pair `(0,2)` gives `d = 4` → the `[2,6,10]` chain. Tracing the seed pair `(i=0, j=1)`, `d = 2`:

| Call | last | scanning k | match? | count after | recurse |
|------|------|------------|--------|-------------|---------|
| extend(1, 2) | 1 (val 4) | k=2 (val 6) | 6-4=2 ✓ | 1 → `[2,4,6]` | extend(2,2) |
| extend(2, 2) | 2 (val 6) | k=3 (val 8) | 8-6=2 ✓ | 2 → `[2,4,6,8]` | extend(3,2) |
| extend(3, 2) | 3 (val 8) | k=4 (val 10) | 10-8=2 ✓ | 3 → `[2,4,6,8,10]` | extend(4,2) |
| extend(4, 2) | 4 | none | — | 3 | — |

Seed `(0,1)` contributed 3 slices. Summing over all seed pairs — `(1,2)→[4,6,8],[4,6,8,10]`, `(2,3)→[6,8,10]`, plus the `d=4` chain `[2,6,10]` — yields **7** total. ✔

---

## Approach 2 — DP with Hash Maps (Optimal)

### Intuition

Track **weak** arithmetic subsequences — those of length ≥ 2 (we relax the ≥ 3 rule to length 2 so pairs can seed longer runs). Let

```
dp[i][d] = number of weak arithmetic subsequences ending at index i with common difference d.
```

Consider an ending pair `(j, i)` with `j < i` and `d = nums[i] - nums[j]`. Every weak subsequence ending at `j` with that **same** `d` already has length ≥ 2; appending `nums[i]` makes it length ≥ 3 — a real answer. So we add `dp[j][d]` to the global count. We also grow `dp[i][d]` by `dp[j][d] + 1`: the `dp[j][d]` extended runs, plus **1** for the brand-new length-2 pair `(j, i)` itself. Because a subsequence is only counted when its *final* element is appended (using the state stored at its second-to-last element), every arithmetic subsequence is counted exactly once.

### Algorithm

1. Give each index `i` a map `dp[i]: difference → count of weak subsequences ending at i`.
2. For each `i`, for each `j < i`:
   - `d = nums[i] - nums[j]`.
   - `cnt = dp[j][d]` — weak runs ending at `j` with difference `d`.
   - `ans += cnt` — each is promoted to a valid length ≥ 3 subsequence by `nums[i]`.
   - `dp[i][d] += cnt + 1` — extend those runs to end at `i`, plus the fresh pair `(j, i)`.
3. Return `ans`.

### Complexity

- **Time:** O(n²) — each of the `~n²/2` ordered pairs `(j, i)` does O(1) amortised hash-map work.
- **Space:** O(n²) — each index may accumulate up to O(n) distinct differences across all the `j < i`.

### Code

```go
func dpHashMap(nums []int) int {
	n := len(nums)
	// dp[i][d] = # of weak (len>=2) arithmetic subsequences ending at i, diff d.
	dp := make([]map[int64]int, n)
	for i := range dp {
		dp[i] = make(map[int64]int)
	}

	ans := 0 // count of STRONG (len>=3) arithmetic subsequences
	for i := 0; i < n; i++ {
		for j := 0; j < i; j++ {
			// Common difference contributed by ending pair (j, i).
			d := int64(nums[i]) - int64(nums[j])
			// Weak subsequences ending at j with this exact difference:
			// each already has length >= 2, so nums[i] promotes it to >= 3.
			cnt := dp[j][d]
			ans += cnt // these are genuine arithmetic subsequences now
			// Extend those weak ones to end at i, and add 1 for the fresh
			// length-2 pair (j, i) that starts a new arithmetic run.
			dp[i][d] += cnt + 1
		}
	}
	return ans
}
```

### Dry Run

Example 1: `nums = [2,4,6,8,10]`. Showing the key `(j, i)` transitions with `d = nums[i]-nums[j]` (only pairs that promote a run are annotated):

| i | j | d | cnt = dp[j][d] | ans += cnt | dp[i][d] += cnt+1 |
|---|---|---|----------------|-----------|--------------------|
| 1 | 0 | 2 | 0 | 0 | dp[1][2]=1 |
| 2 | 0 | 4 | 0 | 0 | dp[2][4]=1 |
| 2 | 1 | 2 | dp[1][2]=1 | **+1** (=1) → `[2,4,6]` | dp[2][2]=2 |
| 3 | 1 | 4 | 0 | 0 | dp[3][4]=1 |
| 3 | 2 | 2 | dp[2][2]=2 | **+2** (=3) → `[2,4,6,8],[4,6,8]` | dp[3][2]=3 |
| 4 | 2 | 4 | dp[2][4]=1 | **+1** (=4) → `[2,6,10]` | dp[4][4]=2 |
| 4 | 3 | 2 | dp[3][2]=3 | **+3** (=7) → `[…6,8,10]` chains | dp[4][2]=4 |

(pairs producing `cnt = 0` still create length-2 seeds but add nothing to `ans`). Final `ans = 7`. ✔

---

## Key Takeaways

- **Encode DP state as (index, difference).** For "arithmetic subsequence" problems, the common difference is part of the state — you cannot fold it away.
- **Weak → strong counting.** Store length-≥ 2 runs so pairs can seed. Count a subsequence at the step where its *last* element is appended, reading the count from its second-to-last element's state. This avoids double counting completely.
- **Hash map replaces an impossible array dimension.** Differences here span the full 32-bit range, so `dp[i]` must be a `map[diff]count`, not a fixed-width table.
- **Guard against overflow.** `nums[i]` covers the full `int32` range, so a difference of two elements can exceed `int32`; compute differences in `int64`.
- The `+1` term is the single most error-prone line: it represents the *new* length-2 pair, which is exactly why longer subsequences are eventually counted.

---

## Related Problems

- LeetCode #413 — Arithmetic Slices (contiguous version, 1-D DP)
- LeetCode #1027 — Longest Arithmetic Subsequence (same (index, diff) DP, tracks length not count)
- LeetCode #1218 — Longest Arithmetic Subsequence of Given Difference (fixed `d`, 1-D hash DP)
- LeetCode #300 — Longest Increasing Subsequence (subsequence DP cousin)
