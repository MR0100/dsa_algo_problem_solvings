# 0153 — Find Minimum in Rotated Sorted Array

> LeetCode #153 · Difficulty: Medium
> **Categories:** Array, Binary Search, Divide and Conquer

---

## Problem Statement

Suppose an array of length `n` sorted in ascending order is **rotated** between `1` and `n` times. For example, the array `nums = [0,1,2,4,5,6,7]` might become:

- `[4,5,6,7,0,1,2]` if it was rotated `4` times.
- `[0,1,2,4,5,6,7]` if it was rotated `7` times.

Notice that **rotating** an array `[a[0], a[1], a[2], ..., a[n-1]]` 1 time results in the array `[a[n-1], a[0], a[1], a[2], ..., a[n-2]]`.

Given the sorted rotated array `nums` of **unique** elements, return *the minimum element of this array*.

You must write an algorithm that runs in `O(log n) time`.

**Example 1:**
```
Input: nums = [3,4,5,1,2]
Output: 1
Explanation: The original array was [1,2,3,4,5] rotated 3 times.
```

**Example 2:**
```
Input: nums = [4,5,6,7,0,1,2]
Output: 0
Explanation: The original array was [0,1,2,4,5,6,7] and it was rotated 4 times.
```

**Example 3:**
```
Input: nums = [11,13,15,17]
Output: 11
Explanation: The original array was [11,13,15,17] and it was rotated 4 times.
```

**Constraints:**
- `n == nums.length`
- `1 <= n <= 5000`
- `-5000 <= nums[i] <= 5000`
- All the integers of `nums` are **unique**.
- `nums` is sorted and rotated between `1` and `n` times.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Meta       | ★★★★★ Very High  | 2024          |
| Amazon     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★★☆ High       | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |
| Goldman Sachs | ★★★☆☆ Medium  | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Binary Search on answer position** — halve the window by deciding which side of `mid` the rotation pivot lies on → see [`/dsa/binary_search.md`](/dsa/binary_search.md)
- **Divide and Conquer** — a fully sorted sub-range answers in O(1); only the half containing the pivot recurses → see [`/dsa/divide_and_conquer.md`](/dsa/divide_and_conquer.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Linear Scan) | O(n) | O(1) | Baseline; ignores structure, fails the O(log n) requirement |
| 2 | Binary Search | O(log n) | O(1) | Required answer; iterative, no stack |
| 3 | Divide and Conquer | O(log n) | O(log n) stack | Same idea recursively; nice for explaining *why* log n works |

---

## Approach 1 — Brute Force (Linear Scan)

### Intuition
The minimum of *any* array is found by inspecting every element. This throws away the rotated-sorted structure, so it can't meet the required O(log n) — but it's the trivially correct baseline to verify the clever versions against.

### Algorithm
1. Set `minVal = nums[0]`.
2. For each remaining element `v`: if `v < minVal`, set `minVal = v`.
3. Return `minVal`.

### Complexity
- **Time:** O(n) — one look at each of the n elements.
- **Space:** O(1) — a single scalar.

### Code
```go
func bruteForce(nums []int) int {
	minVal := nums[0] // candidate minimum
	for _, v := range nums[1:] {
		if v < minVal {
			minVal = v // found a smaller element
		}
	}
	return minVal
}
```

### Dry Run
Example 1: `nums = [3,4,5,1,2]`

| v | v < minVal? | minVal |
|---|-------------|--------|
| (init) | — | 3 |
| 4 | no | 3 |
| 5 | no | 3 |
| 1 | yes | **1** |
| 2 | no | 1 |

Return `1` ✓

---

## Approach 2 — Binary Search (Optimal)

### Intuition
A rotated sorted array is two ascending runs glued together; the minimum is the head of the second run — the **pivot**, the only place where the values "drop". Compare `nums[mid]` against `nums[hi]`:

- `nums[mid] > nums[hi]` → `mid` is still in the *first* (larger-valued) run, so the drop is strictly to the **right** of `mid`: search `(mid, hi]`.
- `nums[mid] < nums[hi]` → `mid..hi` is sorted, so no drop exists there; the minimum is `mid` itself or to its **left**: search `[lo, mid]`.

Comparing against `nums[hi]` (never `nums[lo]`) is the key detail — it makes the fully-sorted case (rotated exactly `n` times, e.g. Example 3) fall out naturally: `nums[mid] < nums[hi]` always holds, the window slides left, and `lo` finishes at index 0. Uniqueness guarantees `nums[mid] != nums[hi]` while `lo < hi`, so no third case exists.

### Algorithm
1. `lo, hi = 0, n-1` (inclusive window always containing the minimum).
2. While `lo < hi`:
   1. `mid = lo + (hi-lo)/2` (floor; overflow-safe form).
   2. If `nums[mid] > nums[hi]`: `lo = mid + 1` (min strictly right of mid).
   3. Else: `hi = mid` (min at mid or left — keep `mid`, it may be the answer).
3. Loop exits when `lo == hi`; return `nums[lo]`.

The invariant "the window `[lo, hi]` contains the minimum" holds at every step, and the window shrinks every iteration, guaranteeing termination at the pivot.

### Complexity
- **Time:** O(log n) — the window halves each iteration.
- **Space:** O(1) — two indices.

### Code
```go
func binarySearch(nums []int) int {
	lo, hi := 0, len(nums)-1 // inclusive search window that contains the min
	for lo < hi {
		mid := lo + (hi-lo)/2 // overflow-safe midpoint (floor)
		if nums[mid] > nums[hi] {
			// mid sits in the FIRST (larger) run → pivot is strictly right
			lo = mid + 1
		} else {
			// nums[mid] < nums[hi] → mid..hi sorted → min at mid or left
			hi = mid // keep mid: it may itself be the minimum
		}
	}
	return nums[lo] // window collapsed onto the pivot
}
```

### Dry Run
Example 1: `nums = [3,4,5,1,2]` (indices 0–4)

| Iter | lo | hi | mid | nums[mid] | nums[hi] | Comparison | Action |
|------|----|----|-----|-----------|----------|------------|--------|
| 1 | 0 | 4 | 2 | 5 | 2 | 5 > 2 | pivot right → `lo = 3` |
| 2 | 3 | 4 | 3 | 1 | 2 | 1 < 2 | sorted → `hi = 3` |
| 3 | 3 | 3 | — | — | — | `lo == hi` | exit loop |

Return `nums[3] = 1` ✓

---

## Approach 3 — Divide and Conquer (Recursive)

### Intuition
If a range `lo..hi` satisfies `nums[lo] <= nums[hi]`, it is fully sorted (unique elements — no plateaus), and its minimum is trivially `nums[lo]`. If not, the range wraps around the pivot. Split it in half: the pivot lies in **exactly one** half, so the other half is sorted and returns in O(1). Only one branch of the recursion ever goes deep — that is precisely why the cost stays logarithmic instead of linear.

### Algorithm
1. Define `rec(lo, hi)`:
   1. If `nums[lo] <= nums[hi]` (sorted or single element), return `nums[lo]`.
   2. Otherwise `mid = lo + (hi-lo)/2`; return `min(rec(lo, mid), rec(mid+1, hi))`.
2. Answer is `rec(0, n-1)`.

### Complexity
- **Time:** O(log n) — at every level, one of the two calls hits the sorted-range base case immediately; only the pivot-containing half recurses.
- **Space:** O(log n) — recursion stack proportional to the depth of halving.

### Code
```go
func divideAndConquer(nums []int) int {
	var rec func(lo, hi int) int
	rec = func(lo, hi int) int {
		if nums[lo] <= nums[hi] {
			// sorted (or single-element) range → smallest is the first entry
			return nums[lo]
		}
		mid := lo + (hi-lo)/2 // split point
		// pivot is in exactly one half; the other half returns immediately
		return min(rec(lo, mid), rec(mid+1, hi))
	}
	return rec(0, len(nums)-1)
}
```

### Dry Run
Example 1: `nums = [3,4,5,1,2]`

| Call | nums[lo], nums[hi] | Sorted? | Action | Returns |
|------|--------------------|---------|--------|---------|
| `rec(0,4)` | 3, 2 | no (3 > 2) | split at mid=2 → `min(rec(0,2), rec(3,4))` | min(3, 1) = **1** |
| `rec(0,2)` | 3, 5 | yes (3 ≤ 5) | base case | 3 |
| `rec(3,4)` | 1, 2 | yes (1 ≤ 2) | base case | 1 |

Return `1` ✓ — note the wrapped half `rec(3,4)` bottomed out instantly; no branch went deeper than one level here.

---

## Key Takeaways

- **Compare `mid` against `hi`, not `lo`** — `nums[lo] <= nums[mid]` is ambiguous when the array isn't rotated at all; `nums[mid]` vs `nums[hi]` cleanly decides which side of the drop you're on in every case.
- **`hi = mid` vs `lo = mid + 1` asymmetry:** keep `mid` when it might be the answer (`nums[mid] < nums[hi]`), discard it when it provably isn't (`nums[mid] > nums[hi]`). This is the standard "binary search for a boundary" shape.
- The minimum of a rotated sorted array **is the rotation pivot** — finding it also tells you the rotation count (`index of min`), which is step 1 of searching rotated arrays (#33).
- Divide and conquer stays O(log n) here only because **one branch always terminates immediately**; if both halves could recurse (as with duplicates in #154), the bound degrades.
- Uniqueness is load-bearing: it eliminates the `nums[mid] == nums[hi]` case. See #154 for what breaks — and how to fix it — when duplicates appear.

---

## Related Problems

- LeetCode #33 — Search in Rotated Sorted Array (same pivot reasoning, then search)
- LeetCode #81 — Search in Rotated Sorted Array II (duplicates allowed)
- LeetCode #154 — Find Minimum in Rotated Sorted Array II (this problem + duplicates)
- LeetCode #702 — Search in a Sorted Array of Unknown Size (binary search variants)
- LeetCode #852 — Peak Index in a Mountain Array (boundary-finding binary search)
