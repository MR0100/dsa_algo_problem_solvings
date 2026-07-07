# 0370 — Range Addition

> LeetCode #370 · Difficulty: Medium
> **Categories:** Array, Prefix Sum

---

## Problem Statement

You are given an integer `length` and an array `updates` where `updates[i] = [startIdxi, endIdxi, inci]`.

You have an array `arr` of length `length` with all zeros, and you have some operation to apply on `arr`. In the `ith` operation, you should increment all the elements `arr[startIdxi], arr[startIdxi + 1], ..., arr[endIdxi]` by `inci`.

Return `arr` *after applying all the* `updates`.

**Example 1:**

```
Input: length = 5, updates = [[1,3,2],[2,4,3],[0,2,-2]]
Output: [-2,0,3,5,3]
```

**Example 2:**

```
Input: length = 10, updates = [[2,4,6],[5,6,8],[1,9,-4]]
Output: [0,-4,2,2,2,4,4,-4,-4,-4]
```

**Constraints:**

- `1 <= length <= 10^5`
- `0 <= updates.length <= 10^4`
- `0 <= startIdxi <= endIdxi < length`
- `-1000 <= inci <= 1000`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Facebook   | ★★★☆☆ Medium     | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Difference Array / Prefix Sum** — record each range update as two O(1) endpoint edits, then reconstruct the array with a single prefix-sum sweep → see [`/dsa/prefix_sum.md`](/dsa/prefix_sum.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Naive Per-Element Updates (Brute Force) | O(length·k) | O(1) aux | Tiny inputs; clarifies the spec |
| 2 | Difference Array (Optimal) | O(length + k) | O(length) | The interview answer; scales to the constraints |

*(k = number of updates.)*

---

## Approach 1 — Naive Per-Element Updates (Brute Force)

### Intuition

The statement says "increment `arr[start..end]` by `inc`". Do exactly that: for each update, loop over its range and add `inc` to each cell. Correct but slow — a single wide update can touch the entire array, and with up to `10^4` updates over `10^5` elements this is `10^9` operations.

### Algorithm

1. `arr` = zeros of size `length`.
2. For each update `[start, end, inc]`: for `i` in `start..end`, `arr[i] += inc`.
3. Return `arr`.

### Complexity

- **Time:** O(length·k) — k updates each potentially spanning the whole array.
- **Space:** O(length) for the output, O(1) auxiliary.

### Code

```go
func bruteForce(length int, updates [][]int) []int {
	arr := make([]int, length) // starts all zeros
	for _, u := range updates {
		start, end, inc := u[0], u[1], u[2]
		for i := start; i <= end; i++ { // touch every index in the range
			arr[i] += inc
		}
	}
	return arr
}
```

### Dry Run

Example 1: `length = 5`, updates `[[1,3,2],[2,4,3],[0,2,-2]]`.

| Update | Applied to | arr after |
|--------|-----------|-----------|
| start `[0,0,0,0,0]` | — | `[0,0,0,0,0]` |
| `[1,3,2]` | idx 1,2,3 += 2 | `[0,2,2,2,0]` |
| `[2,4,3]` | idx 2,3,4 += 3 | `[0,2,5,5,3]` |
| `[0,2,-2]` | idx 0,1,2 += -2 | `[-2,0,3,5,3]` |

Result: `[-2,0,3,5,3]` ✔

---

## Approach 2 — Difference Array (Optimal)

### Intuition

Adding `inc` to the whole range `[start, end]` shows up in the **difference array** as exactly two changes: `diff[start] += inc` (values step up here) and `diff[end+1] -= inc` (they step back down just past the range). After all updates are recorded this way — each in O(1) — a running prefix sum of `diff` reconstructs the real array: every index inside `[start, end]` inherits the `+inc`, and nothing outside does. The `-inc` at `end+1` is omitted when `end+1` is out of bounds (nothing beyond the array to cancel).

### Algorithm

1. `diff` = zeros of size `length` (reused as the output).
2. For each update `[start, end, inc]`: `diff[start] += inc`; if `end+1 < length`, `diff[end+1] -= inc`.
3. Prefix-sum in place: `diff[i] += diff[i-1]` for `i = 1..length-1`.
4. Return `diff`.

### Complexity

- **Time:** O(length + k) — O(1) per update plus one prefix-sum pass.
- **Space:** O(length) — the difference/answer array.

### Code

```go
func differenceArray(length int, updates [][]int) []int {
	diff := make([]int, length) // difference array, reused as the output
	for _, u := range updates {
		start, end, inc := u[0], u[1], u[2]
		diff[start] += inc // value steps up at the range start
		if end+1 < length {
			diff[end+1] -= inc // and steps back down just past the range end
		}
	}
	// Prefix sum turns the difference array back into actual values.
	for i := 1; i < length; i++ {
		diff[i] += diff[i-1]
	}
	return diff
}
```

### Dry Run

Example 1: `length = 5`, updates `[[1,3,2],[2,4,3],[0,2,-2]]`.

Record endpoint edits into `diff` (start size `[0,0,0,0,0]`):

| Update | diff[start] += inc | diff[end+1] -= inc | diff after |
|--------|--------------------|--------------------|------------|
| `[1,3,2]` | diff[1] += 2 | diff[4] -= 2 | `[0,2,0,0,-2]` |
| `[2,4,3]` | diff[2] += 3 | end+1=5 out of bounds | `[0,2,3,0,-2]` |
| `[0,2,-2]` | diff[0] += -2 | diff[3] -= -2 | `[-2,2,3,2,-2]` |

Prefix sum `[-2,2,3,2,-2]`:

| i | diff[i] += diff[i-1] | running |
|---|----------------------|---------|
| 0 | — | -2 |
| 1 | 2 + (-2) | 0 |
| 2 | 3 + 0 | 3 |
| 3 | 2 + 3 | 5 |
| 4 | -2 + 5 | 3 |

Result: `[-2,0,3,5,3]` ✔

---

## Key Takeaways

- **Range-update, query-once ⇒ difference array.** Turn each O(range) update into two O(1) endpoint edits, then one prefix-sum sweep materialises the whole array — O(length + k) total.
- **The two edits:** `+inc` at `start`, `-inc` at `end+1`. The prefix sum "activates" the increment from `start` onward and "cancels" it after `end`.
- **Guard `end+1`:** skip the cancel edit when the range reaches the last index (nothing beyond to undo).
- If updates and queries were *interleaved*, you'd reach for a Fenwick/segment tree instead — the difference array works precisely because all updates finish before the single read.

---

## Related Problems

- LeetCode #1109 — Corporate Flight Bookings (difference array, identical pattern)
- LeetCode #1094 — Car Pooling (difference array over stops)
- LeetCode #598 — Range Addition II (range updates on a matrix, min-overlap trick)
- LeetCode #303 — Range Sum Query Immutable (prefix sum for range queries)
