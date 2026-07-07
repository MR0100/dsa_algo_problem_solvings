# 0324 — Wiggle Sort II

> LeetCode #324 · Difficulty: Medium
> **Categories:** Array, Sorting, Quickselect, Divide and Conquer, Two Pointers

---

## Problem Statement

Given an integer array `nums`, reorder it such that
`nums[0] < nums[1] > nums[2] < nums[3]...`.

You may assume the input array always has a valid answer.

**Example 1:**

```
Input:  nums = [1,5,1,1,6,4]
Output: [1,6,1,5,1,4]
Explanation: [1,4,1,5,1,6] is also accepted.
```

**Example 2:**

```
Input:  nums = [1,3,2,2,3,1]
Output: [2,3,1,3,1,2]
```

**Constraints:**

- `1 <= nums.length <= 5 * 10^4`
- `0 <= nums[i] <= 5000`
- It is guaranteed that there will be an answer for the given input `nums`.

**Follow Up:** Can you do it in `O(n)` time and/or in-place with `O(1)` extra
space?

> Note: there can be many valid answers. This repo's `main()` verifies the strict
> wiggle property `nums[0] < nums[1] > nums[2] < ...` rather than matching one
> exact permutation.

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Google    | ★★★☆☆ Medium     | 2023          |
| Amazon    | ★★★☆☆ Medium     | 2023          |
| Microsoft | ★★☆☆☆ Low        | 2022          |
| Apple     | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Sorting** — the O(n log n) approach sorts then interleaves the halves → see
  [`/dsa/sorting.md`](/dsa/sorting.md)
- **Quickselect** — the optimal approach finds the median in O(n) average → see
  [`/dsa/quickselect.md`](/dsa/quickselect.md)
- **Two Pointers** — both approaches walk two cursors (halves, or a
  Dutch-National-Flag partition) → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **Divide and Conquer** — quickselect recursively narrows to the median →
  see [`/dsa/divide_and_conquer.md`](/dsa/divide_and_conquer.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Sort + Reverse Interleave | O(n log n) | O(n) | Clear, easy to reason about; the go-to |
| 2 | Median + 3-Way Partition + Index Map (Optimal) | O(n) avg | O(1) extra | Meets the O(n)/O(1) follow-up |

---

## Approach 1 — Sort + Reverse Interleave

### Intuition
Sort, then split into a smaller half `S` and larger half `L`. Valleys (even
indices) need small values, peaks (odd indices) need large values. The subtlety
is duplicates near the median: filling each half **from its high end downward**
pushes equal values as far apart as possible, so equal medians can never end up
adjacent when a valid answer exists.

### Algorithm
1. `sorted = sort(nums)`.
2. `mid = (n+1)/2`; small half `= sorted[:mid]`, large half `= sorted[mid:]`.
3. Fill even indices `0,2,4,...` from the top of the small half downward.
   Fill odd indices `1,3,5,...` from the top of the large half downward.
4. Copy the result back into `nums`.

### Complexity
- **Time:** O(n log n) — dominated by the sort.
- **Space:** O(n) — sorted copy plus output buffer.

### Code
```go
func sortInterleave(nums []int) {
	n := len(nums)
	sorted := make([]int, n)
	copy(sorted, nums)
	sort.Ints(sorted)

	mid := (n + 1) / 2
	res := make([]int, n)
	j, k := mid-1, n-1
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			res[i] = sorted[j]
			j--
		} else {
			res[i] = sorted[k]
			k--
		}
	}
	copy(nums, res)
}
```

### Dry Run
Example 1: `nums = [1,5,1,1,6,4]`, `n = 6`.
`sorted = [1,1,1,4,5,6]`, `mid = 3`. Small `= [1,1,1]` (idx 0..2),
Large `= [4,5,6]` (idx 3..5). `j = 2`, `k = 5`.

| i | i%2 | source | value taken | res so far        |
|---|-----|--------|-------------|-------------------|
| 0 | 0   | small[j=2] | 1       | [1,_,_,_,_,_]     |
| 1 | 1   | large[k=5] | 6       | [1,6,_,_,_,_]     |
| 2 | 0   | small[j=1] | 1       | [1,6,1,_,_,_]     |
| 3 | 1   | large[k=4] | 5       | [1,6,1,5,_,_]     |
| 4 | 0   | small[j=0] | 1       | [1,6,1,5,1,_]     |
| 5 | 1   | large[k=3] | 4       | [1,6,1,5,1,4]     |

Result `[1,6,1,5,1,4]` — matches the expected output; `isWiggle` is true.

---

## Approach 2 — Median + 3-Way Partition + Index Mapping (Optimal)

### Intuition
We can achieve the sort+interleave layout without a full sort. The **median**
splits values into greater / equal / less. A **virtual index map** that visits
odd slots first (peaks) then even slots (valleys) lets a **Dutch-National-Flag**
3-way partition drop greater-than-median values onto peaks, less-than-median onto
valleys, and equal-to-median values into the middle band — automatically spread
apart. This is O(n) average and O(1) extra space.

### Algorithm
1. Find the median `m` with quickselect (the `n/2`-th smallest).
2. Define `mapped(i) = (2*i + 1) % (n | 1)` — logical order visits odd indices,
   then even indices.
3. DNF partition over `mapped` indices: `> m` swap to the left region, `< m` swap
   to the right region, `== m` advance. Pointers `left`, `i`, `right`.
4. `nums` is now wiggle-sorted in place.

### Complexity
- **Time:** O(n) average — quickselect O(n) + partition O(n).
- **Space:** O(1) extra beyond the input (quickselect uses a scratch copy here
  for clarity; can be done fully in place).

### Code
```go
func medianPartition(nums []int) {
	n := len(nums)
	if n < 2 {
		return
	}
	m := quickselectMedian(nums)

	mapped := func(i int) int { return (2*i + 1) % (n | 1) }

	i, left, right := 0, 0, n-1
	for i <= right {
		if nums[mapped(i)] > m {
			nums[mapped(left)], nums[mapped(i)] = nums[mapped(i)], nums[mapped(left)]
			left++
			i++
		} else if nums[mapped(i)] < m {
			nums[mapped(right)], nums[mapped(i)] = nums[mapped(i)], nums[mapped(right)]
			right--
		} else {
			i++
		}
	}
}
```

### Dry Run
Example 1: `nums = [1,5,1,1,6,4]`, `n = 6`, `n|1 = 7`.
Quickselect median (`n/2 = 3`-rd smallest of `[1,1,1,4,5,6]`) → `m = 4`.
`mapped(i) = (2i+1) % 7`: `mapped = [1,3,5,0,2,4]`.

| i | mapped(i) | nums[mapped(i)] vs m=4 | action                   |
|---|-----------|------------------------|--------------------------|
| 0 | 1         | 5 > 4                  | swap into peak, left→1, i→1 |
| 1 | 3         | 1 < 4                  | swap into valley, right→4 |
| 2 | 5         | 4 == 4                 | i→3 |
| 3 | 0 (right=4→ mapped 2) | continues per pointers | peaks get {5,6}, valleys get {1,1} |

Equal medians land in the centre band. Final array is a valid wiggle
(`isWiggle` true; e.g. `[1,5,1,6,1,4]`).

---

## Key Takeaways
- **Halves + reverse interleave** is the robust trick: filling both halves from
  their high ends prevents equal medians from touching.
- The follow-up O(n)/O(1) demands **quickselect for the median** plus a
  **3-way (Dutch National Flag) partition** over a **virtual index map**
  `(2*i+1) % (n|1)` that interleaves peaks and valleys.
- Contrast with Wiggle Sort I (#280), where a single greedy pass suffices because
  only `<=`/`>=` (non-strict) is required.

---

## Related Problems
- LeetCode #280 — Wiggle Sort (non-strict; one greedy pass)
- LeetCode #215 — Kth Largest Element (quickselect building block)
- LeetCode #75 — Sort Colors (Dutch National Flag 3-way partition)
- LeetCode #973 — K Closest Points to Origin (quickselect partition)
