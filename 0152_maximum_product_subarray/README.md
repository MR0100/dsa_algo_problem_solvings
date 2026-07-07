# 0152 — Maximum Product Subarray

> LeetCode #152 · Difficulty: Medium
> **Categories:** Array, Dynamic Programming

---

## Problem Statement

Given an integer array `nums`, find a subarray that has the largest product, and return *the product*.

The test cases are generated so that the answer will fit in a **32-bit** integer.

**Example 1:**
```
Input: nums = [2,3,-2,4]
Output: 6
Explanation: [2,3] has the largest product 6.
```

**Example 2:**
```
Input: nums = [-2,0,-1]
Output: 0
Explanation: The result cannot be 2, because [-2,-1] is not a subarray.
```

**Constraints:**
- `1 <= nums.length <= 2 * 10^4`
- `-10 <= nums[i] <= 10`
- The product of any subarray of `nums` is **guaranteed** to fit in a **32-bit** integer.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| LinkedIn   | ★★★★★ Very High  | 2024          |
| Amazon     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2024          |
| Apple      | ★★★☆☆ Medium     | 2023          |
| Adobe      | ★★★☆☆ Medium     | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **1-D Dynamic Programming (Kadane variant)** — "best subarray ending at i" state, extended to a (max, min) pair because multiplication by negatives flips extremes → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)
- **Prefix/Suffix products** — running products from both ends, resetting at zeros, cover every candidate cut around negatives → see [`/dsa/prefix_sum.md`](/dsa/prefix_sum.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(n²) | O(1) | Baseline; sanity-check the fast versions |
| 2 | DP Max/Min Tracking (Kadane variant) | O(n) | O(1) | The canonical interview answer |
| 3 | Prefix/Suffix Sweep | O(n) | O(1) | Elegant alternative; great for explaining *why* the answer is a prefix/suffix of a zero-free block |

---

## Approach 1 — Brute Force

### Intuition
Every subarray is fixed by a start `i` and an end `j`. Instead of recomputing each product from scratch (O(n³)), fix `i` and grow `j` rightward while maintaining a running product: `product(nums[i..j]) = product(nums[i..j-1]) * nums[j]`. Compare every running product against the best seen.

### Algorithm
1. Seed `best = nums[0]` (subarrays are non-empty, so the answer is at least one element).
2. For each start index `i` from `0` to `n-1`:
   1. Reset `prod = 1`.
   2. For each end index `j` from `i` to `n-1`: multiply `prod *= nums[j]`, then update `best = max(best, prod)`.
3. Return `best`.

### Complexity
- **Time:** O(n²) — n start positions, each extended up to n times with O(1) work.
- **Space:** O(1) — two scalar accumulators.

### Code
```go
func bruteForce(nums []int) int {
	best := nums[0] // best product seen so far; seeded with a valid subarray
	for i := 0; i < len(nums); i++ {
		prod := 1 // running product of nums[i..j]
		for j := i; j < len(nums); j++ {
			prod *= nums[j] // extend the subarray by one element
			if prod > best {
				best = prod // record a new maximum
			}
		}
	}
	return best
}
```

### Dry Run
Example 1: `nums = [2,3,-2,4]`, seed `best = 2`

| i | j | prod | best |
|---|---|------|------|
| 0 | 0 | 2 | 2 |
| 0 | 1 | 6 | **6** |
| 0 | 2 | -12 | 6 |
| 0 | 3 | -48 | 6 |
| 1 | 1 | 3 | 6 |
| 1 | 2 | -6 | 6 |
| 1 | 3 | -24 | 6 |
| 2 | 2 | -2 | 6 |
| 2 | 3 | -8 | 6 |
| 3 | 3 | 4 | 6 |

Return `6` ✓

---

## Approach 2 — DP with Max/Min Tracking (Kadane Variant, Optimal)

### Intuition
Kadane's maximum-*sum* algorithm keeps only the best sum ending at `i`, because adding can never turn a bad prefix into a good one. Multiplication can: a hugely **negative** running product becomes hugely **positive** the moment it meets another negative number. So the state must carry *two* values per position — the maximum product `maxEnd` **and** the minimum product `minEnd` of a subarray ending exactly at `i`. When `nums[i]` is negative, the two swap roles before extending.

At each element there are only two choices: extend the previous subarray, or abandon it and start fresh at `nums[i]` (crucial after zeros, which annihilate any product through them).

### Algorithm
1. Initialize `maxEnd = minEnd = best = nums[0]`.
2. For each `i` from `1` to `n-1` with `n = nums[i]`:
   1. If `n < 0`, swap `maxEnd` and `minEnd` — multiplying by a negative maps the largest value to the smallest and vice versa.
   2. `maxEnd = max(n, maxEnd*n)` — extend or restart.
   3. `minEnd = min(n, minEnd*n)` — extend or restart the *worst* product too.
   4. `best = max(best, maxEnd)`.
3. Return `best`.

### Complexity
- **Time:** O(n) — one pass with constant work per element.
- **Space:** O(1) — only the previous DP state is needed, so three scalars replace the table.

### Code
```go
func dpMinMax(nums []int) int {
	maxEnd := nums[0] // max product of a subarray ENDING exactly at i
	minEnd := nums[0] // min product of a subarray ENDING exactly at i
	best := nums[0]   // global answer over all end positions
	for i := 1; i < len(nums); i++ {
		n := nums[i]
		if n < 0 {
			// a negative factor turns the biggest product into the smallest
			// and the smallest into the biggest — swap before extending
			maxEnd, minEnd = minEnd, maxEnd
		}
		// either extend the previous subarray or start a new one at n
		maxEnd = max(n, maxEnd*n)
		minEnd = min(n, minEnd*n)
		if maxEnd > best {
			best = maxEnd // new global maximum
		}
	}
	return best
}
```

### Dry Run
Example 1: `nums = [2,3,-2,4]`, init `maxEnd=2, minEnd=2, best=2`

| i | nums[i] | swap? | maxEnd = max(n, maxEnd·n) | minEnd = min(n, minEnd·n) | best |
|---|---------|-------|---------------------------|---------------------------|------|
| 1 | 3 | no | max(3, 2·3) = **6** | min(3, 2·3) = 3 | **6** |
| 2 | -2 | yes → maxEnd=3, minEnd=6 | max(-2, 3·-2) = **-2** | min(-2, 6·-2) = **-12** | 6 |
| 3 | 4 | no | max(4, -2·4) = **4** | min(4, -12·4) = **-48** | 6 |

Return `6` ✓ — note how `minEnd = -12` at i=2 stood ready to become `+` if another negative had appeared.

---

## Approach 3 — Prefix/Suffix Sweep (Optimal, No DP State)

### Intuition
Split the array at zeros (a zero annihilates every product through it). Inside a zero-free block:
- If the block has an **even** number of negatives, the whole block's product is the maximum.
- If it has an **odd** number, the best subarray drops everything up to and including either the *first* negative or the *last* negative — i.e. the answer is a **suffix** or a **prefix** of the block.

Either way, the answer is always some prefix or suffix product of a zero-free block. Two sweeps — left→right for prefixes, right→left for suffixes — with the running product reset to 1 whenever it hits zero, examine all these candidates. Zeros themselves are still compared (before the reset), so an all-negative-and-zero array like `[-2,0,-1]` correctly answers 0.

### Algorithm
1. `best = nums[0]`, `prod = 1`.
2. Sweep `i = 0 .. n-1`: `prod *= nums[i]`; `best = max(best, prod)`; if `prod == 0`, reset `prod = 1`.
3. Reset `prod = 1`; sweep `i = n-1 .. 0` with the identical update.
4. Return `best`.

### Complexity
- **Time:** O(n) — exactly two linear passes.
- **Space:** O(1) — one running product plus the answer.

### Code
```go
func prefixSuffix(nums []int) int {
	best := nums[0] // must hold at least one element
	prod := 1       // running prefix product
	for i := 0; i < len(nums); i++ {
		prod *= nums[i] // extend the prefix product
		if prod > best {
			best = prod
		}
		if prod == 0 {
			prod = 1 // a zero kills every product through it — restart after it
		}
	}
	prod = 1 // reset for the suffix sweep
	for i := len(nums) - 1; i >= 0; i-- {
		prod *= nums[i] // extend the suffix product
		if prod > best {
			best = prod
		}
		if prod == 0 {
			prod = 1 // restart past the zero
		}
	}
	return best
}
```

### Dry Run
Example 1: `nums = [2,3,-2,4]`, seed `best = 2`

Left→right (prefix) sweep:

| i | nums[i] | prod | best |
|---|---------|------|------|
| 0 | 2 | 2 | 2 |
| 1 | 3 | 6 | **6** |
| 2 | -2 | -12 | 6 |
| 3 | 4 | -48 | 6 |

Right→left (suffix) sweep:

| i | nums[i] | prod | best |
|---|---------|------|------|
| 3 | 4 | 4 | 6 |
| 2 | -2 | -8 | 6 |
| 1 | 3 | -24 | 6 |
| 0 | 2 | -48 | 6 |

Return `6` ✓ — the odd negative at index 2 means the best candidate is the prefix `[2,3]`, found by the first sweep.

---

## Key Takeaways

- **Kadane's for products needs a (max, min) pair** — negatives flip extremes, so the most negative running product is as valuable as the most positive one. This "track both extremes" idea generalizes to any DP where an operation can invert ordering.
- **`max(n, state·n)` encodes "extend or restart"** — the same shape as maximum-sum Kadane; restarting is what makes zeros harmless.
- **Zeros are hard resets** for products: nothing useful crosses a zero, but the zero itself may still be the answer (all other candidates negative).
- The prefix/suffix argument — *the optimal subarray of a zero-free block is a prefix or suffix of it* — is a neat exchange argument worth remembering; it converts a DP problem into two brainless sweeps.
- Small constraint (`|nums[i]| <= 10`, product fits int32) is what makes running products safe; with unbounded values you'd worry about overflow and consider logarithms or big integers.

---

## Related Problems

- LeetCode #53 — Maximum Subarray (the sum version; plain Kadane)
- LeetCode #628 — Maximum Product of Three Numbers (same negative-flip insight, no contiguity)
- LeetCode #713 — Subarray Product Less Than K (sliding window over products)
- LeetCode #238 — Product of Array Except Self (prefix/suffix product machinery)
- LeetCode #1567 — Maximum Length of Subarray With Positive Product (sign-tracking DP)
