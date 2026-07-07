# 0240 — Search a 2D Matrix II

> LeetCode #240 · Difficulty: Medium
> **Categories:** Array, Binary Search, Divide and Conquer, Matrix

---

## Problem Statement

Write an efficient algorithm that searches for a value `target` in an `m x n` integer matrix `matrix`. This matrix has the following properties:

- Integers in each row are sorted in ascending from left to right.
- Integers in each column are sorted in ascending from top to bottom.

**Example 1:**

```
Input: matrix = [[1,4,7,11,15],[2,5,8,12,19],[3,6,9,16,22],[10,13,14,17,24],[18,21,23,26,30]], target = 5
Output: true
```

**Example 2:**

```
Input: matrix = [[1,4,7,11,15],[2,5,8,12,19],[3,6,9,16,22],[10,13,14,17,24],[18,21,23,26,30]], target = 20
Output: false
```

**Constraints:**

- `m == matrix.length`
- `n == matrix[i].length`
- `1 <= n, m <= 300`
- `-10^9 <= matrix[i][j] <= 10^9`
- All the integers in each row are **sorted** in ascending order.
- All the integers in each column are **sorted** in ascending order.
- `-10^9 <= target <= 10^9`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2024          |
| Facebook   | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Matrix Traversal (staircase / elimination walk)** — starting at a corner and eliminating a full row or column each step is the O(m+n) optimal → see [`/dsa/matrix_traversal.md`](/dsa/matrix_traversal.md)
- **Binary Search** — each row (and column) is sorted, so per-row binary search gives an O(m log n) middle-ground solution → see [`/dsa/binary_search.md`](/dsa/binary_search.md)
- **Divide and Conquer** — the sorted 2D structure also admits a quadrant-partitioning recursion (the classic D&C framing) → see [`/dsa/divide_and_conquer.md`](/dsa/divide_and_conquer.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(m·n) | O(1) | Baseline; ignores the sorting |
| 2 | Binary Search Each Row | O(m log n) | O(1) | Uses row order only; simple and fast enough |
| 3 | Staircase Search from Top-Right (Optimal) | O(m + n) | O(1) | The intended answer — exploits both orderings |

---

## Approach 1 — Brute Force

### Intuition

Ignore the sorted structure and check every cell. Correct, and useful as the correctness baseline.

### Algorithm

1. Loop over every row and column.
2. Return `true` on the first cell equal to `target`.
3. Return `false` if the scan finishes with no match.

### Complexity

- **Time:** O(m·n) — visits every cell.
- **Space:** O(1).

### Code

```go
func bruteForce(matrix [][]int, target int) bool {
	for _, row := range matrix {
		for _, v := range row {
			if v == target {
				return true // found it
			}
		}
	}
	return false // no cell matched
}
```

### Dry Run

Example 1: `target = 5`. Scan row by row.

| Cell visited | value | == 5? |
|--------------|-------|-------|
| (0,0) | 1 | no |
| (0,1) | 4 | no |
| (0,2) | 7 | no |
| ... | ... | ... |
| (1,0) | 2 | no |
| (1,1) | 5 | **yes → true** |

Result: `true` ✔

---

## Approach 2 — Binary Search Each Row

### Intuition

Each row is sorted ascending, so a binary search finds `target` within a row in O(log n). Run it on all `m` rows. It exploits the row ordering but not the column ordering — still a solid improvement over brute force.

### Algorithm

1. For each row, binary-search for `target`.
2. Return `true` on the first hit; `false` if none of the rows contain it.

### Complexity

- **Time:** O(m log n) — `m` binary searches, each over `n` columns.
- **Space:** O(1).

### Code

```go
func binarySearchRows(matrix [][]int, target int) bool {
	for _, row := range matrix {
		lo, hi := 0, len(row)-1
		for lo <= hi {
			mid := lo + (hi-lo)/2 // avoid overflow
			switch {
			case row[mid] == target:
				return true
			case row[mid] < target:
				lo = mid + 1 // target is to the right
			default:
				hi = mid - 1 // target is to the left
			}
		}
	}
	return false
}
```

### Dry Run

Example 1: `target = 5`.

| Row | binary search | result |
|-----|---------------|--------|
| `[1,4,7,11,15]` | mid=7 → left; mid=4 → right; mid=... miss | not found |
| `[2,5,8,12,19]` | lo=0,hi=4 → mid=2 (8) → left; lo=0,hi=1 → mid=0 (2) → right; lo=1,hi=1 → mid=1 (**5**) | **true** |

Result: `true` ✔

---

## Approach 3 — Staircase Search from Top-Right (Optimal)

### Intuition

Stand at the **top-right** corner. That cell is the **largest in its row** and the **smallest in its column** — a perfect decision pivot:

- If it's **greater** than `target`, then nothing below it in this column can be `target` either (columns increase downward), so drop the whole column → move **left**.
- If it's **less** than `target`, then nothing to its left in this row can be `target` (rows increase rightward), so drop the whole row → move **down**.

Each step eliminates an entire row or column, so we finish in at most `m + n` steps.

### Algorithm

1. Start at `row = 0`, `col = n-1`.
2. While in bounds, compare `matrix[row][col]` to `target`:
   - equal → return `true`.
   - greater → `col--` (move left).
   - less → `row++` (move down).
3. If we walk off the grid, return `false`.

### Complexity

- **Time:** O(m + n) — every step decrements `col` or increments `row`; at most `m + n` steps.
- **Space:** O(1).

### Code

```go
func staircaseSearch(matrix [][]int, target int) bool {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return false
	}
	row := 0                  // start at the top row
	col := len(matrix[0]) - 1 // ...and the rightmost column
	for row < len(matrix) && col >= 0 {
		switch {
		case matrix[row][col] == target:
			return true // exact hit
		case matrix[row][col] > target:
			col-- // current is too big; drop this column, move left
		default:
			row++ // current is too small; drop this row, move down
		}
	}
	return false // walked off the grid without finding target
}
```

### Dry Run

Example 1: `target = 5`. Start at `(0, 4)` = `15`.

| Step | (row, col) | value | vs 5 | move |
|------|-----------|-------|------|------|
| 1 | (0,4) | 15 | > | left → col=3 |
| 2 | (0,3) | 11 | > | left → col=2 |
| 3 | (0,2) | 7 | > | left → col=1 |
| 4 | (0,1) | 4 | < | down → row=1 |
| 5 | (1,1) | 5 | = | **return true** |

Result: `true` ✔

---

## Key Takeaways

- **Pick a corner that is a min in one direction and a max in the other.** Top-right (or bottom-left) gives an unambiguous "eliminate a whole row or column" decision each step — the staircase / elimination walk. Top-left and bottom-right do **not** work, because there both moves increase or both decrease.
- The staircase search is O(m+n) with O(1) space and no recursion — cleaner than the divide-and-conquer quadrant recursion, which is O(n^1.58) for a square matrix.
- This is a fundamentally different problem from **#74 Search a 2D Matrix**, where the whole grid is one globally sorted sequence and a single binary search over `m·n` works. Here rows and columns are sorted independently, so global binary search does not apply.
- When a structure is monotonic along two axes, look for a pivot cell whose comparison rules out an entire line at once.

---

## Related Problems

- LeetCode #74 — Search a 2D Matrix (globally sorted → single binary search)
- LeetCode #378 — Kth Smallest Element in a Sorted Matrix (heap / binary-search-on-value)
- LeetCode #4 — Median of Two Sorted Arrays (binary search on a 2D-flavored structure)
- LeetCode #1428 — Leftmost Column with at Least a One (staircase walk)
