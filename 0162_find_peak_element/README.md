# 0162 — Find Peak Element

> LeetCode #162 · Difficulty: Medium
> **Categories:** Array, Binary Search

---

## Problem Statement

A peak element is an element that is **strictly greater** than its neighbors.

Given a **0-indexed** integer array `nums`, find a peak element, and return its index. If the array contains multiple peaks, return the index to **any of the peaks**.

You may imagine that `nums[-1] = nums[n] = -∞`. In other words, an element is always considered to be strictly greater than a neighbor that is outside the array.

You must write an algorithm that runs in `O(log n)` time.

**Example 1:**
```
Input: nums = [1,2,3,1]
Output: 2
Explanation: 3 is a peak element and your function should return the index number 2.
```

**Example 2:**
```
Input: nums = [1,2,1,3,5,6,4]
Output: 5
Explanation: Your function can return either index number 1 where the peak element is 2, or index number 5 where the peak element is 6.
```

**Constraints:**
- `1 <= nums.length <= 1000`
- `-2^31 <= nums[i] <= 2^31 - 1`
- `nums[i] != nums[i + 1]` for all valid `i`.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Meta      | ★★★★★ Very High | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Amazon    | ★★★★☆ High      | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |
| Uber      | ★★★☆☆ Medium    | 2023          |
| Apple     | ★★☆☆☆ Low       | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Binary Search** — comparing `nums[mid]` with `nums[mid+1]` tells you which half *must* contain a peak, so the range halves each step even though the array is unsorted → see [`/dsa/binary_search.md`](/dsa/binary_search.md)
- **Divide and Conquer** — the recursive formulation splits the problem into one half-sized subproblem whose answer is the full answer → see [`/dsa/divide_and_conquer.md`](/dsa/divide_and_conquer.md)

---

## Approaches Overview

| # | Approach                        | Time     | Space    | When to use                                                |
|---|---------------------------------|----------|----------|------------------------------------------------------------|
| 1 | Brute Force (Linear Scan)       | O(n)     | O(1)     | Baseline; fine when `n` is tiny, violates the O(log n) ask |
| 2 | Recursive Binary Search         | O(log n) | O(log n) | Same idea as #3; shows the divide-and-conquer structure    |
| 3 | Iterative Binary Search (Optimal) | O(log n) | O(1)   | Always — meets the required bound with constant space      |

---

## Approach 1 — Brute Force (Linear Scan)

### Intuition
Because the virtual sentinels `nums[-1]` and `nums[n]` are −∞ and adjacent elements are never equal, the array conceptually *rises out of* −∞. Walk left to right: the first index `i` where `nums[i] > nums[i+1]` is a peak — it beats its right neighbour by the test, and it beats its left neighbour because we only arrived at `i` by strictly climbing (every earlier step had `nums[k] < nums[k+1]`). If no such step exists the array is strictly increasing, so the last element (backed by the −∞ sentinel on its right) is a peak.

### Algorithm
1. For `i` from `0` to `n−2`:
   1. If `nums[i] > nums[i+1]`, return `i`.
2. Loop finished without a downhill step → return `n−1`.

### Complexity
- **Time:** O(n) — a strictly increasing array forces a full pass.
- **Space:** O(1) — a single loop index.

### Code
```go
func bruteForce(nums []int) int {
	for i := 0; i < len(nums)-1; i++ {
		if nums[i] > nums[i+1] { // first downhill step → i is a peak
			return i
		}
	}
	return len(nums) - 1 // strictly increasing → last element is the peak
}
```

### Dry Run
`nums = [1,2,3,1]` (Example 1):

| i | nums[i] | nums[i+1] | nums[i] > nums[i+1]? | action        |
|---|---------|-----------|----------------------|---------------|
| 0 | 1       | 2         | no (still climbing)  | continue      |
| 1 | 2       | 3         | no (still climbing)  | continue      |
| 2 | 3       | 1         | **yes**              | return **2** ✅ |

---

## Approach 2 — Recursive Binary Search

### Intuition
Binary search normally needs a sorted array — here it only needs the *slope* at the midpoint. Compare `nums[mid]` with `nums[mid+1]`:

- `nums[mid] > nums[mid+1]` — we stand on a descending slope. The values rose from −∞ somewhere on the left and are falling at `mid`, so a peak must exist in `[lo, mid]` (possibly `mid` itself).
- `nums[mid] < nums[mid+1]` — we stand on an ascending slope. The values must eventually fall back to −∞ on the right, so a peak must exist in `[mid+1, hi]`.

Either way half the range is discarded *without ever losing every peak*. (Equality is impossible: neighbours are guaranteed distinct.)

### Algorithm
1. `peakHelper(nums, lo, hi)` with the invariant "`[lo, hi]` contains a peak":
2. If `lo == hi`, return `lo` — the last candidate standing is a peak.
3. `mid = lo + (hi−lo)/2` (note `mid < hi`, so `mid+1` is always in bounds).
4. If `nums[mid] > nums[mid+1]`, recurse on `[lo, mid]`; else recurse on `[mid+1, hi]`.

### Complexity
- **Time:** O(log n) — each call halves the range.
- **Space:** O(log n) — one stack frame per halving (Go does not guarantee tail-call elimination).

### Code
```go
func recursiveBinarySearch(nums []int) int {
	return peakHelper(nums, 0, len(nums)-1)
}

// peakHelper narrows [lo, hi] (which always contains a peak) to one index.
func peakHelper(nums []int, lo, hi int) int {
	if lo == hi { // range shrunk to one candidate → it is a peak
		return lo
	}
	mid := lo + (hi-lo)/2        // overflow-safe midpoint; mid < hi so mid+1 is valid
	if nums[mid] > nums[mid+1] { // descending slope → peak is at mid or left of it
		return peakHelper(nums, lo, mid)
	}
	return peakHelper(nums, mid+1, hi) // ascending slope → peak is right of mid
}
```

### Dry Run
`nums = [1,2,3,1]` (Example 1):

| call | lo | hi | mid | nums[mid] vs nums[mid+1] | decision                  |
|------|----|----|-----|--------------------------|---------------------------|
| 1    | 0  | 3  | 1   | 2 < 3 (ascending)        | recurse on `[2, 3]`       |
| 2    | 2  | 3  | 2   | 3 > 1 (descending)       | recurse on `[2, 2]`       |
| 3    | 2  | 2  | —   | `lo == hi`               | return **2** ✅            |

---

## Approach 3 — Iterative Binary Search (Optimal)

### Intuition
Identical slope-chasing logic as Approach 2, expressed as a loop: keep the invariant "`[lo, hi]` contains at least one peak", shrink the range by half each iteration based on the slope at `mid`, and stop when one index remains. Same O(log n) time, but the recursion stack is gone — O(1) space and no function-call overhead. This is the canonical interview answer.

### Algorithm
1. `lo = 0`, `hi = n−1`.
2. While `lo < hi`:
   1. `mid = lo + (hi−lo)/2`.
   2. If `nums[mid] > nums[mid+1]`, the peak is in `[lo, mid]` → `hi = mid` (keep `mid`: it may itself be the peak).
   3. Else the peak is in `[mid+1, hi]` → `lo = mid + 1` (drop `mid`: an ascending `mid` cannot be a peak).
3. Return `lo` (== `hi`), the surviving peak index.

### Complexity
- **Time:** O(log n) — the window `[lo, hi]` halves every iteration; ~10 iterations for n = 1000.
- **Space:** O(1) — two boundary pointers and a midpoint.

### Code
```go
func iterativeBinarySearch(nums []int) int {
	lo, hi := 0, len(nums)-1 // invariant: [lo, hi] always contains a peak
	for lo < hi {
		mid := lo + (hi-lo)/2 // mid < hi, so nums[mid+1] is always in bounds
		if nums[mid] > nums[mid+1] {
			hi = mid // descending → mid itself may be the peak; keep it
		} else {
			lo = mid + 1 // ascending → mid cannot be a peak; drop it
		}
	}
	return lo // range collapsed to the peak index
}
```

### Dry Run
`nums = [1,2,3,1]` (Example 1):

| iteration | lo | hi | mid | nums[mid] vs nums[mid+1] | update      |
|-----------|----|----|-----|--------------------------|-------------|
| 1         | 0  | 3  | 1   | 2 < 3 (ascending)        | `lo = 2`    |
| 2         | 2  | 3  | 2   | 3 > 1 (descending)       | `hi = 2`    |
| exit      | 2  | 2  | —   | `lo == hi` → loop ends   | return **2** ✅ |

(On Example 2 `[1,2,1,3,5,6,4]` this approach returns index `5`, while the linear scan returns index `1` — both are valid peaks per the statement.)

---

## Key Takeaways

- **Binary search does not require a sorted array** — it requires a predicate that reliably discards half the space. Here the predicate is the local slope `nums[mid] > nums[mid+1]`.
- The −∞ sentinels are what *guarantee* a peak exists: an array rising out of −∞ and falling back into −∞ must turn around somewhere. Always articulate why the invariant "this half contains an answer" holds.
- `hi = mid` vs `lo = mid + 1` asymmetry: keep the index that might still be the answer, drop the one that provably is not. Pairing this with `mid = lo + (hi−lo)/2` (which rounds down, so `mid < hi`) prevents both infinite loops and out-of-bounds access at `mid+1`.
- "Return **any** peak" is the licence that makes O(log n) possible — you never need the global maximum, just *a* local one. Read the output spec carefully before assuming you must scan everything.
- The distinct-neighbours constraint (`nums[i] != nums[i+1]`) removes plateaus; with equal neighbours allowed, this slope argument breaks and the worst case degrades to O(n).

---

## Related Problems

- LeetCode #852 — Peak Index in a Mountain Array (same slope binary search, single guaranteed peak)
- LeetCode #1901 — Find a Peak Element II (the 2-D generalisation)
- LeetCode #33 — Search in Rotated Sorted Array (binary search on a not-fully-sorted array)
- LeetCode #153 — Find Minimum in Rotated Sorted Array (slope/boundary-chasing binary search)
- LeetCode #278 — First Bad Version (binary search over a monotone predicate)
