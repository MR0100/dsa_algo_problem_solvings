# 0360 — Sort Transformed Array

> LeetCode #360 · Difficulty: Medium
> **Categories:** Array, Math, Two Pointers, Sorting

---

## Problem Statement

Given a **sorted** integer array `nums` and three integers `a`, `b` and `c`, apply a quadratic function of the form `f(x) = ax² + bx + c` to each element `nums[i]` in the array, and return *the array in a sorted order*.

**Example 1:**

```
Input: nums = [-4,-2,2,4], a = 1, b = 3, c = 5
Output: [3,9,15,33]
```

**Example 2:**

```
Input: nums = [-4,-2,2,4], a = -1, b = 3, c = 5
Output: [-23,-5,1,7]
```

**Constraints:**

- `1 <= nums.length <= 200`
- `-100 <= nums[i], a, b, c <= 100`
- `nums` is sorted in **ascending** order.

**Follow-up:** Could you solve it in `O(n)` time?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Two Pointers** — a parabola over a sorted array has its extremes at the two ends (a > 0) or the middle (a < 0); converging pointers merge them in O(n) → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **Math (Quadratic / Parabola shape)** — the sign of `a` dictates where the largest/smallest transformed values live → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)
- **Sorting** — the brute-force baseline transforms then sorts → see [`/dsa/sorting.md`](/dsa/sorting.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (transform + sort) | O(n log n) | O(n) | Simple; ignores the sorted-input gift |
| 2 | Two Pointers (Optimal) | O(n) | O(n) | Meets the O(n) follow-up |

---

## Approach 1 — Brute Force

### Intuition

Forget that `nums` is sorted. Apply `f(x) = ax² + bx + c` to every element and sort the results with a comparison sort. Always correct, and a perfect oracle to validate the linear solution.

### Algorithm

1. For each `x` in `nums`, compute `y = a·x² + b·x + c`.
2. Sort the `y`-values ascending.
3. Return them.

### Complexity

- **Time:** O(n log n) — the sort dominates.
- **Space:** O(n) — the output slice.

### Code

```go
func bruteForce(nums []int, a, b, c int) []int {
	result := make([]int, len(nums))
	for i, x := range nums {
		result[i] = apply(x, a, b, c) // transform each element
	}
	sort.Ints(result) // then sort the transformed values
	return result
}
```

### Dry Run

Example 1: `nums = [-4,-2,2,4], a=1, b=3, c=5`.

| x | f(x) = x² + 3x + 5 |
|---|--------------------|
| -4 | 16 - 12 + 5 = 9 |
| -2 | 4 - 6 + 5 = 3 |
| 2 | 4 + 6 + 5 = 15 |
| 4 | 16 + 12 + 5 = 33 |

Transformed: `[9, 3, 15, 33]` → sort → **`[3, 9, 15, 33]`** ✔

---

## Approach 2 — Two Pointers (Optimal)

### Intuition

`f(x) = ax² + bx + c` is a parabola. Over a **sorted** input:

- If `a > 0` (upward parabola), the transformed values are **largest at the two ends** and smallest near the vertex. So compare the two ends and fill the result **from the back** with the larger one.
- If `a < 0` (downward parabola), the values are **smallest at the two ends** and largest near the vertex. Fill the result **from the front** with the smaller end.
- If `a == 0`, the function is linear/monotone; treat it like `a ≥ 0` (largest at the back — correct because a monotone function keeps the sorted order or its reverse, and the end-comparison handles both).

Two pointers from both ends converge in a single pass, placing one element per step.

### Algorithm

1. Evaluate `fl = f(nums[left])` and `fr = f(nums[right])` lazily each step.
2. If `a ≥ 0`: fill from the back (`idx = n-1` downward) with `max(fl, fr)`, advancing that pointer.
3. If `a < 0`: fill from the front (`idx = 0` upward) with `min(fl, fr)`, advancing that pointer.
4. Stop when `left > right`.

### Complexity

- **Time:** O(n) — each element is placed exactly once.
- **Space:** O(n) — the output slice; O(1) auxiliary beyond it.

### Code

```go
func twoPointers(nums []int, a, b, c int) []int {
	n := len(nums)
	result := make([]int, n)
	left, right := 0, n-1 // scan sorted nums from both ends

	if a >= 0 {
		// Upward parabola (or linear): extremes are at the ends → largest first.
		idx := n - 1 // fill position, from the back
		for left <= right {
			fl := apply(nums[left], a, b, c)  // value at the left end
			fr := apply(nums[right], a, b, c) // value at the right end
			if fl >= fr {
				result[idx] = fl // left end is the larger extreme
				left++
			} else {
				result[idx] = fr // right end is the larger extreme
				right--
			}
			idx-- // next slot toward the front
		}
	} else {
		// Downward parabola: extremes (smallest values) are at the ends → fill
		// smallest first from the front.
		idx := 0 // fill position, from the front
		for left <= right {
			fl := apply(nums[left], a, b, c)
			fr := apply(nums[right], a, b, c)
			if fl <= fr {
				result[idx] = fl // left end is the smaller value
				left++
			} else {
				result[idx] = fr // right end is the smaller value
				right--
			}
			idx++ // next slot toward the back
		}
	}
	return result
}
```

### Dry Run

Example 1: `nums = [-4,-2,2,4], a=1` (so `a ≥ 0`, fill from the back). Values: `f(-4)=9, f(-2)=3, f(2)=15, f(4)=33`.

| Step | left | right | fl = f(nums[left]) | fr = f(nums[right]) | pick (larger) | idx | result so far |
|------|------|-------|--------------------|--------------------|--------------|-----|---------------|
| 1 | 0 (-4) | 3 (4) | 9 | 33 | fr=33 → right-- | 3 | `[_,_,_,33]` |
| 2 | 0 (-4) | 2 (2) | 9 | 15 | fr=15 → right-- | 2 | `[_,_,15,33]` |
| 3 | 0 (-4) | 1 (-2) | 9 | 3 | fl=9 → left++ | 1 | `[_,9,15,33]` |
| 4 | 1 (-2) | 1 (-2) | 3 | 3 | fl=3 → left++ | 0 | `[3,9,15,33]` |

`left (2) > right (1)` ⇒ stop. Result: **`[3, 9, 15, 33]`** ✔

---

## Key Takeaways

- **The sign of `a` tells you where the extremes are.** For an upward parabola the biggest transformed values sit at the sorted array's ends; for a downward parabola the ends hold the smallest. This is the whole trick behind the O(n) merge.
- **Fill direction matches the extreme:** `a ≥ 0` → place largest first from the back; `a < 0` → place smallest first from the front. Two pointers converge, one placement per step.
- **`a == 0` is just the monotone case** and folds into the `a ≥ 0` branch — the end-comparison naturally handles both increasing and decreasing linear functions.
- This is a cousin of **Squares of a Sorted Array (#977)**, which is exactly `f(x) = x²` (a = 1, b = c = 0) and uses the same two-pointer merge.

---

## Related Problems

- LeetCode #977 — Squares of a Sorted Array (special case `x²`, two-pointer merge)
- LeetCode #88 — Merge Sorted Array (two-pointer fill from the back)
- LeetCode #4 — Median of Two Sorted Arrays (merge intuition)
- LeetCode #167 — Two Sum II (converging pointers on a sorted array)
