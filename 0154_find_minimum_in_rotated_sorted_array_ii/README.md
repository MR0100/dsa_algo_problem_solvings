# 0154 — Find Minimum in Rotated Sorted Array II

> LeetCode #154 · Difficulty: Hard
> **Categories:** Array, Binary Search, Divide and Conquer

---

## Problem Statement

Suppose an array of length `n` sorted in ascending order is **rotated** between `1` and `n` times. For example, the array `nums = [0,1,4,4,5,6,7]` might become:

- `[4,5,6,7,0,1,4]` if it was rotated `4` times.
- `[0,1,4,4,5,6,7]` if it was rotated `7` times.

Notice that **rotating** an array `[a[0], a[1], a[2], ..., a[n-1]]` 1 time results in the array `[a[n-1], a[0], a[1], a[2], ..., a[n-2]]`.

Given the sorted rotated array `nums` that may contain **duplicates**, return *the minimum element of this array*.

You must decrease the overall operation steps as much as possible.

**Example 1:**
```
Input: nums = [1,3,5]
Output: 1
```

**Example 2:**
```
Input: nums = [2,2,2,0,1]
Output: 0
```

**Constraints:**
- `n == nums.length`
- `1 <= n <= 5000`
- `-5000 <= nums[i] <= 5000`
- `nums` is sorted and rotated between `1` and `n` times.

**Follow-up:** This problem is similar to [Find Minimum in Rotated Sorted Array](../0153_find_minimum_in_rotated_sorted_array/README.md), but `nums` may contain **duplicates**. Would this affect the runtime complexity? How and why?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2024          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Binary Search with degraded worst case** — the LC 153 pivot search plus a "shrink by one on equality" escape hatch for ambiguous comparisons → see [`/dsa/binary_search.md`](/dsa/binary_search.md)
- **Divide and Conquer** — recurse into both halves only when sortedness cannot be proven from the endpoints → see [`/dsa/divide_and_conquer.md`](/dsa/divide_and_conquer.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Linear Scan) | O(n) | O(1) | Matches the worst-case lower bound; simplest correct answer |
| 2 | Binary Search + Duplicate Shrinking | O(log n) avg, O(n) worst | O(1) | The expected interview answer; optimal given the lower bound |
| 3 | Divide and Conquer | O(log n) avg, O(n) worst | O(log n) | Recursive framing; makes the "why O(n) is unavoidable" argument vivid |

---

## Approach 1 — Brute Force (Linear Scan)

### Intuition
Duplicates or not, one pass over the array finds the minimum. What makes this more than a throwaway baseline here: **O(n) is the information-theoretic worst case for this problem.** Consider `[2,2,2,...,2]` versus the same array with a single `0` planted at an arbitrary position — every comparison-based algorithm must, in the worst case, examine every element to tell those inputs apart. So the linear scan is not merely acceptable; no algorithm can beat it on worst-case inputs.

### Algorithm
1. Set `minVal = nums[0]`.
2. For each remaining element `v`: if `v < minVal`, update `minVal = v`.
3. Return `minVal`.

### Complexity
- **Time:** O(n) — each element inspected exactly once.
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
Example 1: `nums = [1,3,5]`

| v | v < minVal? | minVal |
|---|-------------|--------|
| (init) | — | 1 |
| 3 | no | 1 |
| 5 | no | 1 |

Return `1` ✓

---

## Approach 2 — Binary Search with Duplicate Shrinking (Optimal)

### Intuition
Start from the LC 153 skeleton — compare `nums[mid]` with `nums[hi]` to decide which side of the rotation pivot (= the minimum) you're on:

- `nums[mid] > nums[hi]` → the drop is strictly right of `mid` → `lo = mid + 1`.
- `nums[mid] < nums[hi]` → `mid..hi` is sorted → the minimum is at `mid` or left → `hi = mid`.
- `nums[mid] == nums[hi]` → **ambiguous**. Both `[2,2,2,0,2]` (pivot right of mid) and `[2,0,2,2,2]` (pivot left of mid) produce this comparison. No half can be safely discarded.

The escape hatch: when `nums[mid] == nums[hi]`, discard **just `nums[hi]`** (`hi--`). This is always safe: even if `nums[hi]` happened to be the minimum, an element with the *identical value* still sits at `mid` inside the window — the minimum **value** (which is all we must return) survives. Each equality step eliminates one element, so a pathological all-equal array degrades gracefully to O(n), which the lower-bound argument above shows is unavoidable.

**Follow-up answer:** yes, duplicates change the complexity — from a guaranteed O(log n) to *average* O(log n) with an O(n) worst case, because equality between probes destroys the information needed to discard half the window, and an adversary can force equality at every probe.

### Algorithm
1. `lo, hi = 0, n-1`.
2. While `lo < hi`:
   1. `mid = lo + (hi-lo)/2`.
   2. If `nums[mid] > nums[hi]`: `lo = mid + 1`.
   3. Else if `nums[mid] < nums[hi]`: `hi = mid`.
   4. Else (`nums[mid] == nums[hi]`): `hi--`.
3. Return `nums[lo]`.

### Complexity
- **Time:** O(log n) on inputs where probes rarely tie; O(n) worst case — each tie only shrinks the window by 1, and all-duplicate arrays tie on every probe. This meets the proven lower bound.
- **Space:** O(1) — two indices.

### Code
```go
func binarySearch(nums []int) int {
	lo, hi := 0, len(nums)-1 // inclusive window guaranteed to hold the min
	for lo < hi {
		mid := lo + (hi-lo)/2 // overflow-safe midpoint
		switch {
		case nums[mid] > nums[hi]:
			// mid is in the first (larger) run → pivot strictly right of mid
			lo = mid + 1
		case nums[mid] < nums[hi]:
			// mid..hi is sorted → pivot at mid or to its left; keep mid
			hi = mid
		default:
			// nums[mid] == nums[hi]: cannot tell which side holds the pivot.
			// Safe to drop nums[hi] — its value also exists at mid, so the
			// minimum value is still inside the window.
			hi--
		}
	}
	return nums[lo] // window collapsed onto (one copy of) the minimum
}
```

### Dry Run
Example 2: `nums = [2,2,2,0,1]` (indices 0–4) — Example 1 finishes in two sorted-side steps, so the duplicate-handling one is traced instead:

| Iter | lo | hi | mid | nums[mid] | nums[hi] | Case | Action |
|------|----|----|-----|-----------|----------|------|--------|
| 1 | 0 | 4 | 2 | 2 | 1 | 2 > 1 | `lo = 3` |
| 2 | 3 | 4 | 3 | 0 | 1 | 0 < 1 | `hi = 3` |
| 3 | 3 | 3 | — | — | — | `lo == hi` | exit |

Return `nums[3] = 0` ✓

Duplicate-equality case, `nums = [2,2,2,0,2]`:

| Iter | lo | hi | mid | nums[mid] | nums[hi] | Case | Action |
|------|----|----|-----|-----------|----------|------|--------|
| 1 | 0 | 4 | 2 | 2 | 2 | equal | `hi-- → 3` |
| 2 | 0 | 3 | 1 | 2 | 0 | 2 > 0 | `lo = 2` |
| 3 | 2 | 3 | 2 | 2 | 0 | 2 > 0 | `lo = 3` |
| 4 | 3 | 3 | — | — | — | `lo == hi` | exit |

Return `nums[3] = 0` ✓

---

## Approach 3 — Divide and Conquer (Recursive)

### Intuition
In LC 153, `nums[lo] <= nums[hi]` proved a range sorted. With duplicates only the **strict** inequality `nums[lo] < nums[hi]` still proves it — equal endpoints can hide a pivot between them (`[2,2,0,2,2]`). So: strictly ascending endpoints → return `nums[lo]`; otherwise split and take the min of **both** halves. Sorted sub-ranges still prune whole branches early (that's where the average O(log n) comes from), but a plateau array forces both branches at every level → O(n) worst case, mirroring Approach 2's degradation.

### Algorithm
1. Define `rec(lo, hi)`:
   1. If `lo == hi`: single element → return `nums[lo]`.
   2. If `nums[lo] < nums[hi]`: strictly sorted range → return `nums[lo]`.
   3. Else `mid = lo + (hi-lo)/2`; return `min(rec(lo, mid), rec(mid+1, hi))`.
2. Answer is `rec(0, n-1)`.

### Complexity
- **Time:** O(log n) average (sorted halves cut off instantly), O(n) worst case — with all-equal elements the recursion visits every element, matching the lower bound.
- **Space:** O(log n) — ranges halve at every level, so the recursion stack depth is logarithmic even in the worst case.

### Code
```go
func divideAndConquer(nums []int) int {
	var rec func(lo, hi int) int
	rec = func(lo, hi int) int {
		if lo == hi {
			return nums[lo] // single element is its own minimum
		}
		if nums[lo] < nums[hi] {
			return nums[lo] // strictly ascending endpoints → sorted range
		}
		mid := lo + (hi-lo)/2 // split point
		// pivot could be in either half when endpoints are equal → check both
		return min(rec(lo, mid), rec(mid+1, hi))
	}
	return rec(0, len(nums)-1)
}
```

### Dry Run
Example 2: `nums = [2,2,2,0,1]`

| Call | lo == hi? | nums[lo] < nums[hi]? | Action | Returns |
|------|-----------|----------------------|--------|---------|
| `rec(0,4)` | no | 2 < 1? no | split mid=2 → `min(rec(0,2), rec(3,4))` | min(2, 0) = **0** |
| `rec(0,2)` | no | 2 < 2? no | split mid=1 → `min(rec(0,1), rec(2,2))` | min(2, 2) = 2 |
| `rec(0,1)` | no | 2 < 2? no | split mid=0 → `min(rec(0,0), rec(1,1))` | min(2, 2) = 2 |
| `rec(0,0)`, `rec(1,1)`, `rec(2,2)` | yes | — | base case | 2 each |
| `rec(3,4)` | no | 0 < 1? yes | sorted range | 0 |

Return `0` ✓ — the duplicate plateau on the left forced full exploration of that half, while the sorted right half answered in O(1).

---

## Key Takeaways

- **Duplicates break binary search's half-elimination guarantee**: when the probe equals the boundary, neither half can be discarded — the fix is to shed one element (`hi--`), trading worst-case speed for correctness.
- `hi--` is safe *because the minimum is a value, not an index*: a twin of `nums[hi]` remains inside the window at `mid`. (If the problem asked for the *index* of the first minimum, this trick would need more care.)
- **Know the lower-bound argument**: `[2,2,...,2]` vs the same with one hidden `0` is indistinguishable without inspecting every slot — that's why O(n) worst case is *optimal*, and it is the crisp answer to the follow-up.
- The strict vs non-strict comparison (`nums[lo] < nums[hi]` here vs `<=` in LC 153) is exactly the price of duplicates — plateaus make equality uninformative.
- Pattern pair: LC 33 → LC 81 (search) degrades the same way as LC 153 → LC 154 (minimum) for the same reason.

---

## Related Problems

- LeetCode #153 — Find Minimum in Rotated Sorted Array (no duplicates; strict O(log n))
- LeetCode #33 — Search in Rotated Sorted Array (search variant, unique elements)
- LeetCode #81 — Search in Rotated Sorted Array II (search variant with duplicates — same degradation)
- LeetCode #852 — Peak Index in a Mountain Array (boundary binary search)
