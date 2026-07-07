# 0167 — Two Sum II - Input Array Is Sorted

> LeetCode #167 · Difficulty: Medium
> **Categories:** Array, Two Pointers, Binary Search

---

## Problem Statement

Given a **1-indexed** array of integers `numbers` that is already **sorted in non-decreasing order**, find two numbers such that they add up to a specific `target` number. Let these two numbers be `numbers[index1]` and `numbers[index2]` where `1 <= index1 < index2 <= numbers.length`.

Return *the indices of the two numbers,* `index1` *and* `index2`*, **added by one** as an integer array* `[index1, index2]` *of length 2*.

The tests are generated such that there is **exactly one solution**. You **may not** use the same element twice.

Your solution must use only constant extra space.

**Example 1:**

```
Input: numbers = [2,7,11,15], target = 9
Output: [1,2]
Explanation: The sum of 2 and 7 is 9. Therefore, index1 = 1, index2 = 2. We return [1, 2].
```

**Example 2:**

```
Input: numbers = [2,3,4], target = 6
Output: [1,3]
Explanation: The sum of 2 and 4 is 6. Therefore index1 = 1, index2 = 3. We return [1, 3].
```

**Example 3:**

```
Input: numbers = [-1,0], target = -1
Output: [1,2]
Explanation: The sum of -1 and 0 is -1. Therefore index1 = 1, index2 = 2. We return [1, 2].
```

**Constraints:**

- `2 <= numbers.length <= 3 * 10^4`
- `-1000 <= numbers[i] <= 1000`
- `numbers` is sorted in **non-decreasing order**.
- `-1000 <= target <= 1000`
- The tests are generated such that there is **exactly one solution**.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★★ Very High  | 2024          |
| Apple      | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2024          |
| Adobe      | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Yahoo      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Two Pointers** — converging pointers on a sorted array: the sum is monotone in each pointer, so one comparison per step safely discards an element → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **Binary Search** — sortedness turns "find the complement" into an O(log n) search → see [`/dsa/binary_search.md`](/dsa/binary_search.md)
- **Hash Map** — the classic Two Sum (#1) complement lookup, included for contrast even though it wastes the sortedness and violates the O(1)-space requirement → see [`/dsa/hash_map.md`](/dsa/hash_map.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(n²) | O(1) | Baseline only; too slow for n = 3·10⁴ in theory-heavy interviews |
| 2 | Hash Map | O(n) | O(n) | When the array is *not* sorted (i.e., LeetCode #1); violates the O(1)-space rule here |
| 3 | Binary Search | O(n log n) | O(1) | When you only remember binary search; meets the space rule but not optimal |
| 4 | Two Pointers (Optimal) | O(n) | O(1) | Always — the intended solution exploiting sortedness |

---

## Approach 1 — Brute Force

### Intuition

Test every pair. With exactly one solution guaranteed, the first pair that sums to `target` is the answer. This ignores the sorted property entirely — it is the reference point the other approaches improve on.

### Algorithm

1. Loop `i` from `0` to `n-2`.
2. Loop `j` from `i+1` to `n-1`.
3. If `numbers[i] + numbers[j] == target`, return `[i+1, j+1]` (1-indexed).

### Complexity

- **Time:** O(n²) — all ~n²/2 pairs may be inspected before the match.
- **Space:** O(1) — just two loop counters.

### Code

```go
func bruteForce(numbers []int, target int) []int {
	for i := 0; i < len(numbers)-1; i++ {
		for j := i + 1; j < len(numbers); j++ {
			// Exactly one solution exists, so return on the first hit.
			if numbers[i]+numbers[j] == target {
				return []int{i + 1, j + 1} // problem wants 1-indexed positions
			}
		}
	}
	return nil // unreachable: a solution is guaranteed
}
```

### Dry Run

Example 1: `numbers = [2,7,11,15], target = 9`.

| Step | i | j | numbers[i] | numbers[j] | Sum | Sum == 9? | Action |
|------|---|---|-----------|-----------|-----|-----------|--------|
| 1 | 0 | 1 | 2 | 7 | 9 | yes | return `[0+1, 1+1]` = `[1, 2]` |

Result: `[1, 2]` ✔ (lucky early exit; worst case scans all pairs).

---

## Approach 2 — Hash Map

### Intuition

Treat it as classic Two Sum: while scanning, ask "have I already seen `target - x`?" A value → index map answers in O(1). This is optimal *time* without needing sortedness — but it costs O(n) memory, so it breaks this problem's explicit constant-space requirement. Know it as the contrast case.

### Algorithm

1. Create an empty map `seen` (value → 0-based index).
2. For each index `i` with value `x`:
   - If `seen[target-x]` exists as `j`, return `[j+1, i+1]`.
   - Otherwise set `seen[x] = i`.

### Complexity

- **Time:** O(n) — one pass with O(1) average map operations.
- **Space:** O(n) — the map can hold nearly every element before the pair is found.

### Code

```go
func hashMap(numbers []int, target int) []int {
	seen := map[int]int{} // value → 0-based index where it was seen
	for i, x := range numbers {
		// Did an earlier element complete the pair?
		if j, ok := seen[target-x]; ok {
			return []int{j + 1, i + 1} // earlier index first, 1-indexed
		}
		seen[x] = i // remember this value for future complements
	}
	return nil // unreachable: a solution is guaranteed
}
```

### Dry Run

Example 1: `numbers = [2,7,11,15], target = 9`.

| Step | i | x | target−x | In seen? | Action | seen after |
|------|---|---|----------|----------|--------|------------|
| 1 | 0 | 2 | 7 | no | store 2→0 | {2:0} |
| 2 | 1 | 7 | 2 | yes, j = 0 | return `[0+1, 1+1]` = `[1, 2]` | — |

Result: `[1, 2]` ✔

---

## Approach 3 — Binary Search

### Intuition

The array is sorted, so once an element `numbers[i]` is fixed, its required complement `target - numbers[i]` can be binary searched instead of linearly scanned. Searching only the suffix `i+1 .. n-1` guarantees we never reuse the same element and never report a pair twice.

### Algorithm

1. For each `i` from `0` to `n-2`:
2. Compute `need = target - numbers[i]`.
3. Binary search `need` in `numbers[i+1 .. n-1]`:
   - `mid = lo + (hi-lo)/2`;
   - equal → return `[i+1, mid+1]`; smaller → `lo = mid+1`; larger → `hi = mid-1`.

### Complexity

- **Time:** O(n log n) — up to n outer iterations, each with an O(log n) search.
- **Space:** O(1) — only the search bounds.

### Code

```go
func binarySearch(numbers []int, target int) []int {
	for i := 0; i < len(numbers)-1; i++ {
		need := target - numbers[i] // complement we must find
		lo, hi := i+1, len(numbers)-1
		for lo <= hi {
			mid := lo + (hi-lo)/2 // overflow-safe midpoint
			switch {
			case numbers[mid] == need:
				return []int{i + 1, mid + 1} // found the pair, 1-indexed
			case numbers[mid] < need:
				lo = mid + 1 // complement lies to the right
			default:
				hi = mid - 1 // complement lies to the left
			}
		}
	}
	return nil // unreachable: a solution is guaranteed
}
```

### Dry Run

Example 1: `numbers = [2,7,11,15], target = 9`.

| Step | i | need | lo | hi | mid | numbers[mid] | Comparison | Action |
|------|---|------|----|----|-----|--------------|------------|--------|
| 1 | 0 | 9−2 = 7 | 1 | 3 | 2 | 11 | 11 > 7 | hi = 1 |
| 2 | 0 | 7 | 1 | 1 | 1 | 7 | 7 == 7 | return `[0+1, 1+1]` = `[1, 2]` |

Result: `[1, 2]` ✔

---

## Approach 4 — Two Pointers (Optimal)

### Intuition

Place one pointer at each end. Their sum is *steerable*: advancing `left` can only raise the sum (values grow to the right), retreating `right` can only lower it. So compare the sum with `target` and move exactly the pointer that pushes the sum the right way. The discard is safe: if `sum < target`, then `numbers[left]` paired with *anything* still inside the window is even smaller than the current sum — `numbers[left]` can never be part of the answer, so drop it (and symmetrically for `right` when `sum > target`). Exactly one element is eliminated per step, giving linear time with zero extra memory.

### Algorithm

1. `left = 0`, `right = n-1`.
2. While `left < right`:
   - `sum = numbers[left] + numbers[right]`.
   - `sum == target` → return `[left+1, right+1]`.
   - `sum < target` → `left++`.
   - `sum > target` → `right--`.

### Complexity

- **Time:** O(n) — each iteration moves one pointer inward; at most n−1 moves total.
- **Space:** O(1) — two integer indices (satisfies the problem's constant-space rule).

### Code

```go
func twoPointers(numbers []int, target int) []int {
	left, right := 0, len(numbers)-1
	for left < right {
		sum := numbers[left] + numbers[right]
		switch {
		case sum == target:
			return []int{left + 1, right + 1} // 1-indexed answer
		case sum < target:
			left++ // sum too small → advance the small end
		default:
			right-- // sum too big → retreat the large end
		}
	}
	return nil // unreachable: a solution is guaranteed
}
```

### Dry Run

Example 1: `numbers = [2,7,11,15], target = 9`.

| Step | left | right | numbers[left] | numbers[right] | sum | sum vs target | Action |
|------|------|-------|---------------|----------------|-----|---------------|--------|
| 1 | 0 | 3 | 2 | 15 | 17 | 17 > 9 | right-- → 2 |
| 2 | 0 | 2 | 2 | 11 | 13 | 13 > 9 | right-- → 1 |
| 3 | 0 | 1 | 2 | 7 | 9 | 9 == 9 | return `[0+1, 1+1]` = `[1, 2]` |

Result: `[1, 2]` ✔

---

## Key Takeaways

- **Sorted + pair-with-condition ⇒ two pointers.** The converging-pointer pattern works whenever moving each pointer changes the objective monotonically; it is the backbone of 3Sum (#15), 3Sum Closest (#16), and Container With Most Water (#11).
- **The discard argument is the interview answer.** Be ready to justify *why* `left++` is safe when `sum < target`: everything `numbers[left]` could pair with is already ≤ the current sum, so it can never reach `target`.
- **Match the tool to the input's structure.** Unsorted → hash map (#1). Sorted → two pointers (this problem). Sorted but you only need membership → binary search. The three variants of Two Sum exist precisely to teach this decision.
- Watch the **1-indexed** output — returning 0-based indices is the most common wrong-answer here.

---

## Related Problems

- LeetCode #1 — Two Sum (unsorted version: hash map)
- LeetCode #15 — 3Sum (sort + fix one element + this two-pointer scan)
- LeetCode #16 — 3Sum Closest (same converging-pointer engine, tracking nearest sum)
- LeetCode #11 — Container With Most Water (two pointers with a monotone discard argument)
- LeetCode #170 — Two Sum III - Data structure design (streaming variant)
- LeetCode #653 — Two Sum IV - Input is a BST (same idea over a tree)
