# 0315 — Count of Smaller Numbers After Self

> LeetCode #315 · Difficulty: Hard
> **Categories:** Array, Binary Indexed Tree, Segment Tree, Merge Sort, Divide and Conquer, Ordered Set

---

## Problem Statement

Given an integer array `nums`, return an integer array `counts` where `counts[i]` is the number of smaller elements to the right of `nums[i]`.

**Example 1:**

```
Input: nums = [5,2,6,1]
Output: [2,1,1,0]
Explanation:
To the right of 5 there are 2 smaller elements (2 and 1).
To the right of 2 there is only 1 smaller element (1).
To the right of 6 there is 1 smaller element (1).
To the right of 1 there is 0 smaller element.
```

**Example 2:**

```
Input: nums = [-1]
Output: [0]
```

**Example 3:**

```
Input: nums = [-1,-1]
Output: [0,0]
```

**Constraints:**

- `1 <= nums.length <= 10^5`
- `-10^4 <= nums[i] <= 10^4`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★★☆☆ Medium     | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Binary Indexed Tree (Fenwick) / Segment Tree** — O(log n) prefix-count and point-update over value ranks → see [`/dsa/segment_tree_fenwick.md`](/dsa/segment_tree_fenwick.md)
- **Merge sort / inversion counting** — the answer is a per-element inversion count computed during merging → see [`/dsa/divide_and_conquer.md`](/dsa/divide_and_conquer.md)
- **Coordinate compression / sorting** — mapping values to compact ranks → see [`/dsa/sorting.md`](/dsa/sorting.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (count pairs) | O(n²) | O(1) | Baseline; tiny inputs |
| 2 | Binary Indexed Tree (Fenwick) | O(n log n) | O(n) | Clean, general "count smaller seen" pattern |
| 3 | Merge Sort (inversion count, Optimal) | O(n log n) | O(n) | No coordinate compression needed |

---

## Approach 1 — Brute Force

### Intuition
The definition is literal: `counts[i]` = number of `j > i` with `nums[j] < nums[i]`. Just check every such pair.

### Algorithm
1. For each `i`, set `c = 0`.
2. For each `j > i`, if `nums[j] < nums[i]`, increment `c`.
3. `counts[i] = c`.

### Complexity
- **Time:** O(n²) — every pair `(i, j)` with `j > i`.
- **Space:** O(1) extra (excluding output).

### Code
```go
func bruteForce(nums []int) []int {
	n := len(nums)
	counts := make([]int, n)
	for i := 0; i < n; i++ {
		c := 0
		for j := i + 1; j < n; j++ {
			if nums[j] < nums[i] {
				c++
			}
		}
		counts[i] = c
	}
	return counts
}
```

### Dry Run
`nums = [5,2,6,1]`.

| i | nums[i] | elements to the right | smaller ones | counts[i] |
|---|---------|-----------------------|--------------|-----------|
| 0 | 5 | 2,6,1 | 2,1 | 2 |
| 1 | 2 | 6,1 | 1 | 1 |
| 2 | 6 | 1 | 1 | 1 |
| 3 | 1 | (none) | — | 0 |

Result: `[2,1,1,0]`.

---

## Approach 2 — Binary Indexed Tree (Fenwick)

### Intuition
Process **right-to-left** so "already inserted" means "to my right". A Fenwick tree over compressed value-ranks answers, in O(log n): "how many inserted values have rank strictly less than mine?" Then insert the current value. Coordinate-compress values to ranks `1..m` to bound the tree.

### Algorithm
1. Compress values to sorted distinct ranks (1-based).
2. Iterate `i` from `n-1` down to `0`: `counts[i] = query(rank[i]-1)` (inserted values with smaller rank), then `update(rank[i])`.

### Complexity
- **Time:** O(n log n) — `n` updates and queries, each O(log n).
- **Space:** O(n) — the Fenwick tree and rank map.

### Code
```go
func fenwick(nums []int) []int {
	n := len(nums)
	counts := make([]int, n)

	sorted := make([]int, n)
	copy(sorted, nums)
	sort.Ints(sorted)
	rank := map[int]int{}
	r := 0
	for _, v := range sorted {
		if _, ok := rank[v]; !ok {
			r++
			rank[v] = r
		}
	}

	tree := make([]int, r+1)
	update := func(i int) {
		for ; i <= r; i += i & (-i) {
			tree[i]++
		}
	}
	query := func(i int) int {
		s := 0
		for ; i > 0; i -= i & (-i) {
			s += tree[i]
		}
		return s
	}

	for i := n - 1; i >= 0; i-- {
		ri := rank[nums[i]]
		counts[i] = query(ri - 1)
		update(ri)
	}
	return counts
}
```

### Dry Run
`nums = [5,2,6,1]`. Sorted distinct = `[1,2,5,6]` → ranks: `1→1, 2→2, 5→3, 6→4`. Process right to left.

| i | nums[i] | rank | query(rank-1) = counts[i] | after update, inserted ranks |
|---|---------|------|---------------------------|------------------------------|
| 3 | 1 | 1 | query(0) = 0 | {1} |
| 2 | 6 | 4 | query(3) = 1 (rank 1 present) | {1,4} |
| 1 | 2 | 2 | query(1) = 1 (rank 1 present) | {1,2,4} |
| 0 | 5 | 3 | query(2) = 2 (ranks 1,2 present) | {1,2,3,4} |

Result: `[2,1,1,0]`.

---

## Approach 3 — Merge Sort (Inversion Count, Optimal)

### Intuition
`counts[i]` is the number of inversions `(i, j)`, `j > i`, `nums[j] < nums[i]`. Merge sort counts inversions naturally. Sort **indices** (not values) so each original index keeps its own tally. While merging two sorted halves, track how many right-half elements have already been placed; when we place a **left** index, that many right elements are both smaller and to its right, so add them to its count.

### Algorithm
1. `idx = [0..n-1]`, `counts = zeros`.
2. Recursively sort `idx` by `nums` value. During merge, keep `rightMerged` = right-half elements already placed; when taking a left index, `counts[idx[i]] += rightMerged`.
3. Return `counts`.

### Complexity
- **Time:** O(n log n) — merge sort.
- **Space:** O(n) — index buffers and recursion.

### Code
```go
func mergeSortCount(nums []int) []int {
	n := len(nums)
	counts := make([]int, n)
	idx := make([]int, n)
	for i := range idx {
		idx[i] = i
	}

	var sortRange func(lo, hi int)
	sortRange = func(lo, hi int) {
		if hi-lo <= 1 {
			return
		}
		mid := (lo + hi) / 2
		sortRange(lo, mid)
		sortRange(mid, hi)

		merged := make([]int, 0, hi-lo)
		i, j := lo, mid
		rightMerged := 0
		for i < mid && j < hi {
			if nums[idx[j]] < nums[idx[i]] {
				rightMerged++
				merged = append(merged, idx[j])
				j++
			} else {
				counts[idx[i]] += rightMerged
				merged = append(merged, idx[i])
				i++
			}
		}
		for i < mid {
			counts[idx[i]] += rightMerged
			merged = append(merged, idx[i])
			i++
		}
		for j < hi {
			merged = append(merged, idx[j])
			j++
		}
		copy(idx[lo:hi], merged)
	}
	sortRange(0, n)
	return counts
}
```

### Dry Run
`nums = [5,2,6,1]`, `idx = [0,1,2,3]`.

Split into `[0,1]` (values 5,2) and `[2,3]` (values 6,1).

Merge left `[0,1]`: values 5,2 → right cand 2<5 so place idx 1 first (`rightMerged=1`), then place idx 0 with `counts[0] += 1`. Left sorted = `[1,0]` (values 2,5), `counts[0]=1`.

Merge right `[2,3]`: values 6,1 → place idx 3 (`rightMerged=1`), then idx 2 with `counts[2] += 1`. Right sorted = `[3,2]` (values 1,6), `counts[2]=1`.

Final merge of `[1,0]` (2,5) and `[3,2]` (1,6):

| compare | action | rightMerged | counts update |
|---------|--------|-------------|---------------|
| 1 (val 1) < 2 | place idx 3 | 1 | — |
| 2 vs 6 | place idx 1 (left) | 1 | counts[1] += 1 → 1 |
| 5 vs 6 | place idx 0 (left) | 1 | counts[0] += 1 → 2 |
| left done | place idx 2 | — | — |

Final `counts = [2,1,1,0]`.

---

## Key Takeaways
- **"Count smaller to the right" = inversion count per element.** Two canonical O(n log n) tools: a Fenwick/segment tree over value ranks, or merge sort that tallies inversions during the merge.
- **Process right-to-left with a Fenwick tree** so "already inserted" equals "to the right"; query the prefix strictly below the current rank.
- **Coordinate compression** keeps the Fenwick tree small and handles negative values.
- **Sort indices, not values,** in the merge-sort variant so each element retains its own running count.

---

## Related Problems
- LeetCode #493 — Reverse Pairs (merge sort inversion variant)
- LeetCode #327 — Count of Range Sum (merge sort / BIT on prefix sums)
- LeetCode #308 — Range Sum Query 2D Mutable (BIT/segment tree)
- LeetCode #1649 — Create Sorted Array through Instructions (Fenwick counting)
