# 0283 — Move Zeroes

> LeetCode #283 · Difficulty: Easy
> **Categories:** Array, Two Pointers

---

## Problem Statement

Given an integer array `nums`, move all `0`'s to the end of it while maintaining the relative order of the non-zero elements.

**Note** that you must do this in-place without making a copy of the array.

**Example 1:**
```
Input: nums = [0,1,0,3,12]
Output: [1,3,12,0,0]
```

**Example 2:**
```
Input: nums = [0]
Output: [0]
```

**Constraints:**
- `1 <= nums.length <= 10⁴`
- `-2³¹ <= nums[i] <= 2³¹ - 1`

**Follow-up:** Could you minimize the total number of operations done?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Meta       | ★★★★★ Very High  | 2024          |
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used
- **Two pointers** — a read pointer and a write/boundary pointer compact non-zeros in one pass → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **In-place array partitioning** — rearrange within the same buffer with O(1) extra space → see [`/dsa/sorting.md`](/dsa/sorting.md)

---

## Approaches Overview
| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (extra array) | O(n) | O(n) | Clearest baseline; violates in-place rule |
| 2 | Two Passes (overwrite then fill) | O(n) | O(1) | In-place, easy to reason about |
| 3 | Two Pointers Swap (Optimal) | O(n) | O(1) | Fewest writes; answers the follow-up |

---

## Approach 1 — Brute Force

### Intuition
The result is "non-zeros in order, then zeros". Most literally: collect the non-zeros into a buffer, pad the rest with zeros, copy back. It ignores the in-place requirement but is the obvious correctness baseline.

### Algorithm
1. Walk `nums`, appending each non-zero to a buffer.
2. Extend the buffer with `0`s until it has `len(nums)` entries.
3. Copy the buffer back over `nums`.

### Complexity
- **Time:** O(n) — two linear passes.
- **Space:** O(n) — the auxiliary buffer (violates the in-place follow-up).

### Code
```go
func bruteForce(nums []int) {
	buf := make([]int, 0, len(nums))
	for _, v := range nums {
		if v != 0 {
			buf = append(buf, v)
		}
	}
	for len(buf) < len(nums) {
		buf = append(buf, 0)
	}
	copy(nums, buf)
}
```

### Dry Run
`nums = [0,1,0,3,12]`.

| step | action | buf |
|------|--------|-----|
| scan 0 | skip | [] |
| scan 1 | append | [1] |
| scan 0 | skip | [1] |
| scan 3 | append | [1,3] |
| scan 12 | append | [1,3,12] |
| pad | +0,+0 | [1,3,12,0,0] |

Copy back ⇒ `nums = [1,3,12,0,0]`.

---

## Approach 2 — Two Passes, Overwrite Then Fill

### Intuition
Keep a write cursor `insert`. Scan left to right; every non-zero is written at `nums[insert]` and `insert` advances. After the scan, everything from `insert` to the end is leftover and must be zero, so overwrite it. O(1) space.

### Algorithm
1. `insert = 0`. For each `v` in `nums`: if `v != 0`, `nums[insert] = v`; `insert++`.
2. For `i` from `insert` to `n-1`: set `nums[i] = 0`.

### Complexity
- **Time:** O(n) — a compaction pass plus a fill pass.
- **Space:** O(1) — in place.

### Code
```go
func twoPass(nums []int) {
	insert := 0
	for _, v := range nums {
		if v != 0 {
			nums[insert] = v
			insert++
		}
	}
	for i := insert; i < len(nums); i++ {
		nums[i] = 0
	}
}
```

### Dry Run
`nums = [0,1,0,3,12]`.

| i | nums[i] | insert (in) | write? | array | insert (out) |
|---|---------|-------------|--------|-------|--------------|
| 0 | 0 | 0 | no | [0,1,0,3,12] | 0 |
| 1 | 1 | 0 | nums[0]=1 | [1,1,0,3,12] | 1 |
| 2 | 0 | 1 | no | [1,1,0,3,12] | 1 |
| 3 | 3 | 1 | nums[1]=3 | [1,3,0,3,12] | 2 |
| 4 | 12 | 2 | nums[2]=12 | [1,3,12,3,12] | 3 |

Fill pass: `nums[3]=0`, `nums[4]=0` ⇒ `[1,3,12,0,0]`.

---

## Approach 3 — Two Pointers, Swap (Optimal)

### Intuition
Maintain `last` = the index where the next non-zero belongs (everything before it is non-zero). When scanning meets a non-zero at `i`, swap `nums[i]` with `nums[last]` and advance `last`. When `i == last` nothing moves; a real swap only happens when a zero sits at `nums[last]`. So each element is touched at most once — the fewest writes, order preserved. This directly answers the follow-up.

### Algorithm
1. `last = 0`. For `i` from `0` to `n-1`:
2. If `nums[i] != 0`: swap `nums[i]` and `nums[last]`; `last++`.

### Complexity
- **Time:** O(n) — single pass.
- **Space:** O(1) — in place, order-preserving.

### Code
```go
func twoPointers(nums []int) {
	last := 0
	for i := 0; i < len(nums); i++ {
		if nums[i] != 0 {
			nums[i], nums[last] = nums[last], nums[i]
			last++
		}
	}
}
```

### Dry Run
`nums = [0,1,0,3,12]`.

| i | nums[i] | last (in) | swap | array | last (out) |
|---|---------|-----------|------|-------|------------|
| 0 | 0 | 0 | — | [0,1,0,3,12] | 0 |
| 1 | 1 | 0 | swap(1,0) | [1,0,0,3,12] | 1 |
| 2 | 0 | 1 | — | [1,0,0,3,12] | 1 |
| 3 | 3 | 1 | swap(3,1) | [1,3,0,0,12] | 2 |
| 4 | 12 | 2 | swap(4,2) | [1,3,12,0,0] | 3 |

Output: `[1,3,12,0,0]`.

---

## Key Takeaways
- **Slow/fast (write/read) pointers** are the canonical pattern for stable in-place partitioning — reuse it for "remove element", "remove duplicates", etc.
- The **swap variant minimises writes**: it only moves data when a non-zero must jump over a zero, which is what the follow-up rewards.
- Filling the tail with a sentinel (here `0`) after compaction is a clean two-pass alternative when swapping feels error-prone.

---

## Related Problems
- LeetCode #27 — Remove Element (in-place removal, write pointer)
- LeetCode #26 — Remove Duplicates from Sorted Array (write pointer)
- LeetCode #75 — Sort Colors (Dutch national flag, three-way partition)
- LeetCode #905 — Sort Array By Parity (partition by predicate)
