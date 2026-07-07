# 0396 — Rotate Function

> LeetCode #396 · Difficulty: Medium
> **Categories:** Array, Math, Dynamic Programming

---

## Problem Statement

You are given an integer array `nums` of length `n`.

Assume `arrk` to be an array obtained by rotating `nums` by `k` positions clockwise. We define the **rotation function** `F` on `nums` as follow:

- `F(k) = 0 * arrk[0] + 1 * arrk[1] + ... + (n - 1) * arrk[n - 1]`.

Return *the maximum value of* `F(0), F(1), ..., F(n-1)`.

The test cases are generated so that the answer fits in a **32-bit** integer.

**Example 1:**

```
Input: nums = [4,3,2,6]
Output: 26
Explanation:
F(0) = (0 * 4) + (1 * 3) + (2 * 2) + (3 * 6) = 0 + 3 + 4 + 18 = 25
F(1) = (0 * 6) + (1 * 4) + (2 * 3) + (3 * 2) = 0 + 4 + 6 + 6 = 16
F(2) = (0 * 2) + (1 * 6) + (2 * 4) + (3 * 3) = 0 + 6 + 8 + 9 = 23
F(3) = (0 * 3) + (1 * 2) + (2 * 6) + (3 * 4) = 0 + 2 + 12 + 12 = 26
So the maximum value of F(0), F(1), F(2), F(3) is F(3) = 26.
```

**Example 2:**

```
Input: nums = [100]
Output: 0
```

**Constraints:**

- `n == nums.length`
- `1 <= n <= 10^5`
- `-100 <= nums[i] <= 100`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |
| Adobe      | ★☆☆☆☆ Rare       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Math / Algebraic Recurrence** — the whole trick is expanding `F(k)` and `F(k-1)` symbolically to discover the O(1) increment `F(k) = F(k-1) + sum − n·nums[n-k]` → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)
- **Prefix/Running Aggregate** — we precompute `sum = Σ nums` and `F(0) = Σ i·nums[i]` in a single pass, then reuse them → see [`/dsa/prefix_sum.md`](/dsa/prefix_sum.md)
- **Dynamic Programming (rolling state)** — each `F(k)` depends only on `F(k-1)`, a classic 1-D rolling recurrence → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(n²) | O(1) | Small `n`; verifying the recurrence; too slow for `n = 10⁵` |
| 2 | Rolling Recurrence (Optimal) | O(n) | O(1) | Always — derives each `F(k)` from `F(k-1)` in O(1) |

---

## Approach 1 — Brute Force

### Intuition

`F(k)` is defined explicitly, so just compute all `n` of them and take the max. Rather than physically rotating the array, note that in `arrk` the element sitting at index `i` came from original index `(i − k)` modulo `n`. Indexing that way lets us evaluate any `F(k)` directly against the untouched `nums`.

### Algorithm

1. Let `n = len(nums)`.
2. For each rotation `k` in `0..n-1`:
   1. `sum = 0`.
   2. For each index `i` in `0..n-1`: add `i * nums[(i-k+n)%n]`.
   3. Update the running maximum with `sum`.
3. Return the maximum.

### Complexity

- **Time:** O(n²) — `n` rotations, each requiring an O(n) summation.
- **Space:** O(1) — only scalar accumulators; no array copies.

### Code

```go
func bruteForce(nums []int) int {
	n := len(nums)
	best := 0            // will hold the maximum F(k)
	for k := 0; k < n; k++ { // try every rotation amount
		sum := 0
		for i := 0; i < n; i++ {
			// In arrk, the element at index i came from original index
			// (i-k) modulo n; +n keeps the result non-negative before %.
			sum += i * nums[(i-k+n)%n]
		}
		if k == 0 || sum > best { // seed best on first k, then maximise
			best = sum
		}
	}
	return best
}
```

### Dry Run

Example 1: `nums = [4,3,2,6]`, `n = 4`.

| k | sum computation | F(k) | best after |
|---|-----------------|------|------------|
| 0 | 0·4 + 1·3 + 2·2 + 3·6 | 25 | 25 |
| 1 | 0·6 + 1·4 + 2·3 + 3·2 | 16 | 25 |
| 2 | 0·2 + 1·6 + 2·4 + 3·3 | 23 | 25 |
| 3 | 0·3 + 1·2 + 2·6 + 3·4 | 26 | **26** |

Result: `26` ✔

---

## Approach 2 — Rolling Recurrence (Optimal)

### Intuition

Expand the difference `F(k) − F(k-1)`. Rotating clockwise by one more position bumps every element's coefficient up by 1 — that adds one full `sum = Σ nums` to the total. The single exception is the element whose coefficient wraps from `(n-1)` back to `0`: it loses `n · value`. That wrapping element is `nums[n-k]`. Hence:

```
F(k) = F(k-1) + sum − n · nums[n-k]
```

Compute `F(0)` and `sum` once, then slide `k` from `1` to `n-1` applying this O(1) update, tracking the max.

### Algorithm

1. In one pass compute `sum = Σ nums[i]` and `f = F(0) = Σ i·nums[i]`.
2. Initialise `best = f`.
3. For `k` from `1` to `n-1`: set `f = f + sum − n·nums[n-k]`, then `best = max(best, f)`.
4. Return `best`.

### Complexity

- **Time:** O(n) — one pass to seed `sum`/`F(0)`, one pass for the recurrence.
- **Space:** O(1) — three integer accumulators, no extra arrays.

### Code

```go
func rollingRecurrence(nums []int) int {
	n := len(nums)
	sum := 0 // Σ nums[i]
	f := 0   // running F(k); starts as F(0)
	for i, v := range nums {
		sum += v     // total of all elements
		f += i * v   // F(0) = Σ i·nums[i]
	}
	best := f // F(0) is our first candidate
	for k := 1; k < n; k++ {
		// F(k) = F(k-1) + sum - n*nums[n-k]: every coefficient +1 (adds sum),
		// but the wrap element nums[n-k] drops from coeff (n-1) to 0.
		f = f + sum - n*nums[n-k]
		if f > best {
			best = f
		}
	}
	return best
}
```

### Dry Run

Example 1: `nums = [4,3,2,6]`, `n = 4`, `sum = 15`, `F(0) = 25`.

| k | nums[n-k] | update f = f + sum − n·nums[n-k] | f | best |
|---|-----------|----------------------------------|---|------|
| — | — | seed F(0) | 25 | 25 |
| 1 | nums[3]=6 | 25 + 15 − 4·6 = 25 + 15 − 24 | 16 | 25 |
| 2 | nums[2]=2 | 16 + 15 − 4·2 = 16 + 15 − 8 | 23 | 25 |
| 3 | nums[1]=3 | 23 + 15 − 4·3 = 23 + 15 − 12 | 26 | **26** |

Result: `26` ✔ — identical values to the brute force, produced in linear time.

---

## Key Takeaways

- **Differencing adjacent terms** (`F(k) − F(k-1)`) is the go-to move whenever you must evaluate a family of related sums — it usually collapses O(n) recomputation into an O(1) update.
- The wrap term is `nums[n-k]`: the element that just fell from the highest coefficient `(n-1)` back to `0`. Getting this index right is the whole difficulty.
- Precomputing `sum` and `F(0)` once turns the problem into a **1-D rolling DP** where the state is a single integer `f`.
- Range/rotation aggregates often have a closed-form step. When you see "compute this for every shift", look for the increment identity before iterating naively.

---

## Related Problems

- LeetCode #189 — Rotate Array (the rotation mechanic itself)
- LeetCode #918 — Maximum Sum Circular Subarray (circular-array aggregates)
- LeetCode #53 — Maximum Subarray (rolling-state maximisation)
- LeetCode #1888 — Minimum Number of Flips (sliding aggregate over rotations)
