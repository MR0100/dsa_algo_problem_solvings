# 0238 — Product of Array Except Self

> LeetCode #238 · Difficulty: Medium
> **Categories:** Array, Prefix Sum

---

## Problem Statement

Given an integer array `nums`, return *an array* `answer` *such that* `answer[i]` *is equal to the product of all the elements of* `nums` *except* `nums[i]`.

The product of any prefix or suffix of `nums` is **guaranteed** to fit in a **32-bit** integer.

You must write an algorithm that runs in `O(n)` time and without using the division operation.

**Example 1:**

```
Input: nums = [1,2,3,4]
Output: [24,12,8,6]
```

**Example 2:**

```
Input: nums = [-1,1,0,-3,3]
Output: [0,0,9,0,0]
```

**Constraints:**

- `2 <= nums.length <= 10^5`
- `-30 <= nums[i] <= 30`
- The input is generated such that `answer[i]` is **guaranteed** to fit in a **32-bit** integer.

**Follow up:** Can you solve the problem in `O(1)` extra space complexity? (The output array **does not** count as extra space for space complexity analysis.)

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★★ Very High  | 2024          |
| Facebook   | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★★☆ High       | 2024          |
| Apple      | ★★★☆☆ Medium     | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Prefix / Suffix Products** — the answer at each index is the product of a left-prefix and a right-suffix; precomputing running products in both directions replaces division → see [`/dsa/prefix_sum.md`](/dsa/prefix_sum.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(n²) | O(1) extra | Baseline only; TLEs at n = 10⁵ |
| 2 | Prefix × Suffix Arrays | O(n) | O(n) | Clear, easy-to-explain O(n) version |
| 3 | Two-Pass O(1) Extra Space (Optimal) | O(n) | O(1) extra | The follow-up answer — no division, no extra array |

---

## Approach 1 — Brute Force

### Intuition

Follow the definition literally: for each index `i`, multiply together every element whose index isn't `i`. Two nested loops. This ignores the O(n) / no-division constraints but establishes correctness.

### Algorithm

1. For each `i`, start `product = 1`.
2. For each `j != i`, multiply `product *= nums[j]`.
3. Store `answer[i] = product`.

### Complexity

- **Time:** O(n²) — for every index we scan the whole array.
- **Space:** O(1) extra (output array excluded).

### Code

```go
func bruteForce(nums []int) []int {
	n := len(nums)
	answer := make([]int, n)
	for i := 0; i < n; i++ {
		product := 1 // running product of every element except nums[i]
		for j := 0; j < n; j++ {
			if j != i {
				product *= nums[j] // multiply in every other element
			}
		}
		answer[i] = product
	}
	return answer
}
```

### Dry Run

Example 1: `nums = [1,2,3,4]`.

| i | elements multiplied (j ≠ i) | answer[i] |
|---|-----------------------------|-----------|
| 0 | 2·3·4 | 24 |
| 1 | 1·3·4 | 12 |
| 2 | 1·2·4 | 8 |
| 3 | 1·2·3 | 6 |

Result: `[24, 12, 8, 6]` ✔

---

## Approach 2 — Prefix × Suffix Arrays

### Intuition

The product of "everything except `i`" factors cleanly:

```
answer[i] = (product of nums left of i) × (product of nums right of i)
```

Precompute the left products in a forward sweep and the right products in a backward sweep, then multiply the two per index. No division required.

### Algorithm

1. `prefix[0] = 1`; for `i ≥ 1`, `prefix[i] = prefix[i-1] * nums[i-1]`.
2. `suffix[n-1] = 1`; for `i ≤ n-2`, `suffix[i] = suffix[i+1] * nums[i+1]`.
3. `answer[i] = prefix[i] * suffix[i]`.

### Complexity

- **Time:** O(n) — three linear passes.
- **Space:** O(n) — two auxiliary arrays.

### Code

```go
func prefixSuffixArrays(nums []int) []int {
	n := len(nums)
	prefix := make([]int, n) // prefix[i] = product of everything left of i
	suffix := make([]int, n) // suffix[i] = product of everything right of i
	answer := make([]int, n)

	prefix[0] = 1
	for i := 1; i < n; i++ {
		prefix[i] = prefix[i-1] * nums[i-1] // accumulate leftward products
	}
	suffix[n-1] = 1
	for i := n - 2; i >= 0; i-- {
		suffix[i] = suffix[i+1] * nums[i+1] // accumulate rightward products
	}
	for i := 0; i < n; i++ {
		answer[i] = prefix[i] * suffix[i] // left product × right product
	}
	return answer
}
```

### Dry Run

Example 1: `nums = [1,2,3,4]`.

| i | prefix[i] (left) | suffix[i] (right) | answer[i] = prefix·suffix |
|---|------------------|-------------------|---------------------------|
| 0 | 1 | 2·3·4 = 24 | 24 |
| 1 | 1 | 3·4 = 12 | 12 |
| 2 | 1·2 = 2 | 4 | 8 |
| 3 | 1·2·3 = 6 | 1 | 6 |

Result: `[24, 12, 8, 6]` ✔

---

## Approach 3 — Two-Pass O(1) Extra Space (Optimal)

### Intuition

We don't need two separate arrays. Store the **prefix** products directly in the output array (the output slot isn't counted as extra space). Then make a single right-to-left pass carrying a scalar `right` that accumulates the suffix product, multiplying it into each slot on the way. That merges left and right using only one extra variable.

### Algorithm

1. First pass (left→right): `answer[i]` = product of all elements **left** of `i` (`answer[0] = 1`).
2. Second pass (right→left) with running `right = 1`: `answer[i] *= right`, then `right *= nums[i]`.

### Complexity

- **Time:** O(n) — two passes.
- **Space:** O(1) extra — the output array plus one scalar.

### Code

```go
func prefixSuffixInPlace(nums []int) []int {
	n := len(nums)
	answer := make([]int, n)

	answer[0] = 1
	for i := 1; i < n; i++ {
		answer[i] = answer[i-1] * nums[i-1] // answer[i] = product of left part
	}
	right := 1 // running product of everything to the right of i
	for i := n - 1; i >= 0; i-- {
		answer[i] *= right // fold in the right-side product
		right *= nums[i]   // extend the suffix product to include nums[i]
	}
	return answer
}
```

### Dry Run

Example 1: `nums = [1,2,3,4]`.

**Pass 1 (prefix into answer):**

| i | answer[i] = left product |
|---|--------------------------|
| 0 | 1 |
| 1 | 1 |
| 2 | 1·2 = 2 |
| 3 | 1·2·3 = 6 |

`answer = [1, 1, 2, 6]`

**Pass 2 (right→left, running `right`):**

| i | answer[i] *= right | right updated to |
|---|--------------------|------------------|
| 3 | 6 × 1 = 6 | 1·4 = 4 |
| 2 | 2 × 4 = 8 | 4·3 = 12 |
| 1 | 1 × 12 = 12 | 12·2 = 24 |
| 0 | 1 × 24 = 24 | 24·1 = 24 |

Result: `[24, 12, 8, 6]` ✔

---

## Key Takeaways

- **Split "all except i" into left × right.** Any per-index aggregate over "everything but this element" decomposes into a prefix contribution and a suffix contribution — the core prefix/suffix pattern.
- **The output array can double as scratch space.** Storing the prefix pass in the result, then folding the suffix in with one scalar, achieves the O(1)-extra-space follow-up.
- **Avoiding division** matters when zeros are present (Example 2): a single zero makes the total product zero and division breaks. The prefix/suffix product handles zeros naturally — a slot excludes exactly one element.
- The two-pass idea reappears in trapping-rain-water and candy-distribution style problems: sweep once each direction, combine.

---

## Related Problems

- LeetCode #42 — Trapping Rain Water (left-max / right-max two-sweep pattern)
- LeetCode #135 — Candy (two directional sweeps combined)
- LeetCode #152 — Maximum Product Subarray
- LeetCode #303 — Range Sum Query - Immutable (prefix sums)
- LeetCode #724 — Find Pivot Index (prefix vs suffix balance)
