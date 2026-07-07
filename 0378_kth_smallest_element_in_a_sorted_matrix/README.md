# 0378 — Kth Smallest Element in a Sorted Matrix

> LeetCode #378 · Difficulty: Medium
> **Categories:** Binary Search, Heap (Priority Queue), Matrix, Sorting

---

## Problem Statement

Given an `n x n` `matrix` where each of the rows and columns is sorted in ascending order, return *the* `kth` *smallest element in the matrix*.

Note that it is the `kth` smallest element **in the sorted order**, not the `kth` **distinct** element.

You must find a solution with a memory complexity better than `O(n²)`.

**Example 1:**

```
Input: matrix = [[1,5,9],[10,11,13],[12,13,15]], k = 8
Output: 13
Explanation: The elements in the matrix are [1,5,9,10,11,12,13,13,15], and the 8th smallest number is 13.
```

**Example 2:**

```
Input: matrix = [[-5]], k = 1
Output: -5
```

**Constraints:**

- `n == matrix.length == matrix[i].length`
- `1 <= n <= 300`
- `-10^9 <= matrix[i][j] <= 10^9`
- All the rows and columns of `matrix` are **guaranteed** to be sorted in **non-decreasing order**.
- `1 <= k <= n²`

**Follow-up:** Could you solve the problem with a constant memory (i.e., `O(1)` memory complexity)?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Facebook   | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Uber       | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Binary Search** — search the *answer value*, not an index; count entries ≤ mid to decide which half → see [`/dsa/binary_search.md`](/dsa/binary_search.md)
- **Heap / Priority Queue** — merge n sorted rows with a min-heap of row fronts → see [`/dsa/heap_priority_queue.md`](/dsa/heap_priority_queue.md)
- **Matrix Traversal** — the O(n) "staircase" count walks from a corner exploiting both sortings → see [`/dsa/matrix_traversal.md`](/dsa/matrix_traversal.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (flatten + sort) | O(n² log n) | O(n²) | Baseline; violates the sub-n² memory hint |
| 2 | Min-Heap (k-way merge) | O(k log n) | O(n) | Great when k is small; uses the row sortedness |
| 3 | Binary Search on Value (Optimal) | O(n log(max−min)) | O(1) | Meets the O(1) follow-up; best for large k |

---

## Approach 1 — Brute Force (Flatten and Sort)

### Intuition

Forget the structure: dump all `n²` values into one slice, sort it, and read index `k-1`. Trivially correct and a good sanity check against the clever methods, but it uses O(n²) memory — which the problem explicitly asks us to beat.

### Algorithm

1. Flatten every row into one slice.
2. Sort ascending.
3. Return element at index `k-1` (k is 1-based).

### Complexity

- **Time:** O(n² log n) — sorting all `n²` values.
- **Space:** O(n²) — the flattened array.

### Code

```go
func bruteForce(matrix [][]int, k int) int {
	var flat []int
	for _, row := range matrix {
		flat = append(flat, row...)
	}
	sort.Ints(flat)
	return flat[k-1]
}
```

### Dry Run

Example 1: `matrix = [[1,5,9],[10,11,13],[12,13,15]], k = 8`.

| Step | Value |
|------|-------|
| Flatten | `[1,5,9,10,11,13,12,13,15]` |
| Sort | `[1,5,9,10,11,12,13,13,15]` |
| Index k-1 = 7 | `13` |

Result: **13** ✔

---

## Approach 2 — Min-Heap (K-Way Merge)

### Intuition

Each **row** is sorted, so the matrix is n sorted lists to be merged. The next-smallest global value is always the minimum across the current front of each row — a min-heap of at most n items. Pop k times; each pop advances one step to the right in that row.

### Algorithm

1. Push `(matrix[r][0], r, 0)` for every row `r`.
2. Pop the smallest. If it has a right neighbour `col+1` in its row, push `(matrix[r][col+1], r, col+1)`.
3. After `k` pops, the last popped value is the answer.

### Complexity

- **Time:** O(k log n) — k pops, each heap operation O(log n); heap holds ≤ n items.
- **Space:** O(n) — one heap entry per row.

### Code

```go
func minHeapMerge(matrix [][]int, k int) int {
	n := len(matrix)
	h := &itemHeap{}
	for r := 0; r < n; r++ {
		heap.Push(h, minHeapItem{val: matrix[r][0], row: r, col: 0})
	}
	var popped minHeapItem
	for i := 0; i < k; i++ {
		popped = heap.Pop(h).(minHeapItem)
		if popped.col+1 < n {
			heap.Push(h, minHeapItem{
				val: matrix[popped.row][popped.col+1],
				row: popped.row,
				col: popped.col + 1,
			})
		}
	}
	return popped.val
}
```

### Dry Run

Example 1: `matrix = [[1,5,9],[10,11,13],[12,13,15]], k = 8`. Heap seeded with row fronts `1,10,12`.

| Pop # | popped val | push next in row | heap fronts after |
|-------|-----------|------------------|-------------------|
| 1 | 1 | 5 | 5,10,12 |
| 2 | 5 | 9 | 9,10,12 |
| 3 | 9 | — (row end) | 10,12 |
| 4 | 10 | 11 | 11,12 |
| 5 | 11 | 13 | 12,13 |
| 6 | 12 | 13 | 13,13 |
| 7 | 13 | 15 | 13,15 |
| 8 | 13 | — | … |

8th pop = **13** ✔

---

## Approach 3 — Binary Search on Value (Optimal)

### Intuition

The answer is some value in `[matrix[0][0], matrix[n-1][n-1]]`. For a candidate `mid`, the count of entries `≤ mid` is **monotonically non-decreasing** in `mid`. We binary-search for the smallest value whose count is `≥ k`; that value must itself be a matrix element. Counting is O(n) with a **staircase walk** from the bottom-left corner: at `(row, col)`, if the value `≤ mid` the whole column up to `row` qualifies — add `row+1` and step right; otherwise step up.

### Algorithm

1. `lo = matrix[0][0]`, `hi = matrix[n-1][n-1]`.
2. While `lo < hi`: `mid = lo + (hi-lo)/2`. If `countLessEqual(mid) < k`, set `lo = mid+1`; else `hi = mid`.
3. Return `lo`.

### Complexity

- **Time:** O(n · log(max − min)) — each of ~log(value range) iterations counts in O(n).
- **Space:** O(1) — meets the follow-up.

### Code

```go
func binarySearchValue(matrix [][]int, k int) int {
	n := len(matrix)
	lo, hi := matrix[0][0], matrix[n-1][n-1]
	countLessEqual := func(target int) int {
		count := 0
		row, col := n-1, 0
		for row >= 0 && col < n {
			if matrix[row][col] <= target {
				count += row + 1
				col++
			} else {
				row--
			}
		}
		return count
	}
	for lo < hi {
		mid := lo + (hi-lo)/2
		if countLessEqual(mid) < k {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return lo
}
```

### Dry Run

Example 1: `matrix = [[1,5,9],[10,11,13],[12,13,15]], k = 8`. `lo=1, hi=15`.

| lo | hi | mid | countLessEqual(mid) | ≥ k=8? | new range |
|----|----|-----|---------------------|--------|-----------|
| 1 | 15 | 8 | 3 (1,5,... ) | no (<8) | lo=9 |
| 9 | 15 | 12 | 6 | no (<8) | lo=13 |
| 13 | 15 | 14 | 8 | yes | hi=14 |
| 13 | 14 | 13 | 8 | yes | hi=13 |
| 13 | 13 | — | loop ends | | |

Return `lo = 13`. Result: **13** ✔

---

## Key Takeaways

- **Binary search on the answer** applies whenever a predicate ("how many ≤ x?") is monotone in x — even when the answer isn't an array index.
- The **staircase count** from a corner is the standard O(n) trick for a row-and-column-sorted matrix (also used in Search a 2D Matrix II).
- **Min-heap k-way merge** is the go-to for "merge n sorted lists / find kth across sorted streams"; it wins when k is small relative to n².
- The binary-search answer meets the **O(1) memory** follow-up that both other methods miss.

---

## Related Problems

- LeetCode #240 — Search a 2D Matrix II (same staircase walk)
- LeetCode #373 — Find K Pairs with Smallest Sums (heap k-way merge)
- LeetCode #668 — Kth Smallest Number in Multiplication Table (binary search on value)
- LeetCode #719 — Find K-th Smallest Pair Distance (binary search on value + count)
- LeetCode #23 — Merge k Sorted Lists (heap k-way merge)
