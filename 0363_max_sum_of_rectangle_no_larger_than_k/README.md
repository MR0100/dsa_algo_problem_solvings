# 0363 — Max Sum of Rectangle No Larger Than K

> LeetCode #363 · Difficulty: Hard
> **Categories:** Array, Matrix, Prefix Sum, Binary Search, Ordered Set

---

## Problem Statement

Given an `m x n` matrix `matrix` and an integer `k`, return *the max sum of a rectangle in the matrix such that its sum is no larger than* `k`.

It is **guaranteed** that there will be a rectangle with a sum no larger than `k`.

**Example 1:**

```
Input: matrix = [[1,0,1],[0,-2,3]], k = 2
Output: 2
Explanation: Because the sum of the blue rectangle [[0, 1], [-2, 3]] is 2,
and 2 is the max number no larger than k (k = 2).
```

**Example 2:**

```
Input: matrix = [[2,2,-1]], k = 3
Output: 3
```

**Constraints:**

- `m == matrix.length`
- `n == matrix[i].length`
- `1 <= m, n <= 100`
- `-100 <= matrix[i][j] <= 100`
- `-10^5 <= k <= 10^5`

**Follow up:** What if the number of rows is much larger than the number of columns?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Facebook   | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Citadel    | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Prefix Sum** — compressing a band of rows into per-column sums, then prefix sums of that 1-D array, is the core reduction → see [`/dsa/prefix_sum.md`](/dsa/prefix_sum.md)
- **Binary Search on an ordered set** — for each running prefix `P` we binary-search the smallest earlier prefix `≥ P-k` to maximise a subarray sum capped at `k` → see [`/dsa/binary_search.md`](/dsa/binary_search.md)
- **Matrix Traversal** — the row-pair compression sweeps the grid to reduce 2-D to 1-D → see [`/dsa/matrix_traversal.md`](/dsa/matrix_traversal.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (all corner pairs) | O(m³·n³) | O(1) | Tiny grids / correctness reference only |
| 2 | Column Compression + Prefix Sums | O(m²·n²) | O(n) | Clear intermediate; passes with n,m ≤ 100 |
| 3 | Compression + Sorted-Prefix Binary Search (Optimal) | O(m²·n log n) | O(n) | Intended answer; transpose so the smaller side is squared |

---

## Approach 1 — Brute Force

### Intuition

A rectangle is fixed by its top-left and bottom-right corners. Enumerate all corner pairs, sum every cell inside, and keep the largest sum that does not exceed `k`. Obviously correct, hopelessly slow.

### Algorithm

1. For every top-left `(r1,c1)` and bottom-right `(r2,c2)` with `r2≥r1, c2≥c1`:
   1. Sum all cells in `[r1..r2] × [c1..c2]`.
   2. If `sum ≤ k`, update `best`.
2. Return `best`.

### Complexity

- **Time:** O(m³·n³) — O(m²n²) rectangles, each re-summing up to O(mn) cells.
- **Space:** O(1) — a couple of accumulators.

### Code

```go
func bruteForce(matrix [][]int, k int) int {
	m, n := len(matrix), len(matrix[0])
	const negInf = -1 << 60
	best := negInf
	for r1 := 0; r1 < m; r1++ {
		for c1 := 0; c1 < n; c1++ {
			for r2 := r1; r2 < m; r2++ {
				for c2 := c1; c2 < n; c2++ {
					// Sum every cell of rectangle [r1..r2] x [c1..c2].
					sum := 0
					for i := r1; i <= r2; i++ {
						for j := c1; j <= c2; j++ {
							sum += matrix[i][j]
						}
					}
					// Keep the largest sum that stays within the cap k.
					if sum <= k && sum > best {
						best = sum
					}
				}
			}
		}
	}
	return best
}
```

### Dry Run

Example 1: `matrix = [[1,0,1],[0,-2,3]], k = 2`. A few rectangles (top-left → bottom-right):

| Rectangle (r1,c1)-(r2,c2) | Cells | Sum | ≤ k=2? | best |
|---------------------------|-------|-----|--------|------|
| (0,0)-(0,0) | [1] | 1 | yes | 1 |
| (0,0)-(1,2) | whole matrix | 1+0+1+0-2+3 = 3 | no | 1 |
| (0,1)-(1,2) | [0,1],[-2,3] | 0+1-2+3 = 2 | yes | **2** |
| (1,2)-(1,2) | [3] | 3 | no | 2 |

Max sum ≤ 2 is `2` ✔

---

## Approach 2 — Column Compression + Prefix Sums

### Intuition

Fixing a top row `r1` and a bottom row `r2` collapses the 2-D problem to 1-D: `colSum[c]` = sum of column `c` between those rows. The best rectangle spanning `r1..r2` is now the best contiguous **subarray** of `colSum` with sum ≤ k. Extend the band downward one row at a time so `colSum` is built incrementally instead of recomputed.

### Algorithm

1. For each top row `r1`, zero `colSum`.
2. For each bottom row `r2 ≥ r1`, add row `r2` into `colSum`.
3. Scan every subarray `colSum[i..j]` (O(n²)); if its sum ≤ k, update `best`.
4. Return `best`.

### Complexity

- **Time:** O(m²·n²) — O(m²) row pairs × O(n²) subarray scan.
- **Space:** O(n) — the compressed column array.

### Code

```go
func prefixSumRowSearch(matrix [][]int, k int) int {
	m, n := len(matrix), len(matrix[0])
	const negInf = -1 << 60
	best := negInf
	for r1 := 0; r1 < m; r1++ {
		colSum := make([]int, n) // band sum per column for rows r1..r2
		for r2 := r1; r2 < m; r2++ {
			for c := 0; c < n; c++ {
				colSum[c] += matrix[r2][c] // extend the band down by one row
			}
			// Best contiguous subarray of colSum with sum ≤ k, O(n²) scan.
			for i := 0; i < n; i++ {
				sum := 0
				for j := i; j < n; j++ {
					sum += colSum[j] // subarray colSum[i..j]
					if sum <= k && sum > best {
						best = sum
					}
				}
			}
		}
	}
	return best
}
```

### Dry Run

Example 1, `k=2`. Row band `r1=0, r2=1` (both rows), so `colSum = [1+0, 0-2, 1+3] = [1, -2, 4]`.

| Subarray of colSum | Sum | ≤ 2? | best |
|--------------------|-----|------|------|
| [1]        | 1  | yes | 1 |
| [1,-2]     | -1 | yes | 1 |
| [1,-2,4]   | 3  | no  | 1 |
| [-2]       | -2 | yes | 1 |
| [-2,4]     | 2  | yes | **2** |
| [4]        | 4  | no  | 2 |

Band `[-2,4]` corresponds to columns 1–2 across both rows = rectangle `[[0,1],[-2,3]]`, sum 2 ✔ (single-row bands `r1=r2` yield nothing larger and ≤ 2).

---

## Approach 3 — Compression + Sorted-Prefix Binary Search (Optimal)

### Intuition

Keep the row-band compression, but solve the 1-D "max subarray sum ≤ k" faster. A subarray sum equals `prefix[j] - prefix[i]`. To maximise this while staying ≤ k, for the current prefix `P` we need the **smallest earlier prefix ≥ P - k** (that gives the largest `P - prefix` not exceeding k). Keep earlier prefixes in a sorted set and binary-search that lower bound. Seed the set with `0` so subarrays starting at column 0 are considered.

### Algorithm

1. For each top row `r1`, zero `colSum`.
2. For each bottom row `r2`, extend `colSum` by row `r2`.
3. Walk a running `prefix` over `colSum`, keeping a sorted slice `seen` (seeded `{0}`):
   1. Binary-search `seen` for the smallest value `≥ prefix - k`.
   2. If found (`lo`), candidate `= prefix - lo ≤ k`; update `best`.
   3. Insert `prefix` into `seen`, keeping it sorted.
4. Return `best`. (If `m > n`, transpose first so the smaller dimension is squared — the follow-up.)

### Complexity

- **Time:** O(m²·n log n) — O(m²) row pairs, each an O(n log n) prefix scan (binary search + sorted insert). Transposing makes it O(min(m,n)²·max(m,n)·log(min(m,n))).
- **Space:** O(n) — `colSum` plus the sorted prefix set.

### Code

```go
func sortedPrefixBinarySearch(matrix [][]int, k int) int {
	m, n := len(matrix), len(matrix[0])
	const negInf = -1 << 60
	best := negInf
	for r1 := 0; r1 < m; r1++ {
		colSum := make([]int, n) // band sum per column for rows r1..r2
		for r2 := r1; r2 < m; r2++ {
			for c := 0; c < n; c++ {
				colSum[c] += matrix[r2][c] // extend band down by one row
			}
			// Find max subarray sum ≤ k in colSum using sorted prefixes.
			seen := []int{0} // prefixes seen so far; 0 = empty prefix
			prefix := 0
			for c := 0; c < n; c++ {
				prefix += colSum[c] // running prefix P = sum of colSum[0..c]
				// We want the smallest earlier prefix ≥ prefix-k so that
				// prefix - thatPrefix ≤ k and is as large as possible.
				target := prefix - k
				idx := sort.SearchInts(seen, target) // first seen[idx] ≥ target
				if idx < len(seen) {
					if cand := prefix - seen[idx]; cand > best {
						best = cand // best sum ≤ k for a subarray ending at c
					}
				}
				// Insert prefix into `seen` keeping it sorted (insertion sort).
				pos := sort.SearchInts(seen, prefix)
				seen = append(seen, 0)
				copy(seen[pos+1:], seen[pos:]) // shift right to open a gap
				seen[pos] = prefix
			}
		}
	}
	return best
}
```

### Dry Run

Example 1, `k=2`, row band `r1=0, r2=1`: `colSum = [1, -2, 4]`. `seen = {0}`, `best = -∞`.

| c | colSum[c] | prefix P | target = P-k | smallest seen ≥ target | candidate P-lo | best | seen after insert |
|---|-----------|----------|--------------|------------------------|----------------|------|-------------------|
| 0 | 1  | 1  | 1-2 = -1 | 0 (≥ -1)   | 1-0 = 1 | 1 | {0,1} |
| 1 | -2 | -1 | -1-2 = -3 | 0 (≥ -3) — wait smallest ≥ -3 in {0,1} is 0 | -1-0 = -1 | 1 | {-1,0,1} |
| 2 | 4  | 3  | 3-2 = 1  | 1 (≥ 1)    | 3-1 = **2** | **2** | {-1,0,1,3} |

At `c=2`, the smallest prefix ≥ 1 is `1` (the prefix after column 0), giving subarray `colSum[1..2] = -2+4 = 2 ≤ k`. `best = 2` ✔

---

## Key Takeaways

- **Reduce 2-D to 1-D by fixing two rows** and compressing the band into per-column sums — the workhorse trick for "best rectangle" problems.
- **"Max subarray sum ≤ k" in 1-D** is solved with prefix sums + an ordered set: for prefix `P`, find the smallest earlier prefix `≥ P-k`. This is the same lower-bound-search idea as counting subarrays with bounded sums.
- **Seed the prefix set with 0** so subarrays anchored at the start are eligible; forgetting this misses whole-prefix rectangles.
- **Transpose to square the smaller dimension.** The cost is O(min²·max·log min); when rows ≫ columns, iterating over column pairs instead answers the follow-up.
- A Go `[]int` with `sort.SearchInts` + insertion emulates an ordered multiset (a TreeSet/`SortedList` in other languages).

---

## Related Problems

- LeetCode #53 — Maximum Subarray (1-D, no cap)
- LeetCode #325 — Maximum Size Subarray Sum Equals k (prefix + hash map)
- LeetCode #560 — Subarray Sum Equals K (prefix-sum counting)
- LeetCode #85 — Maximal Rectangle (row compression, different objective)
- LeetCode #304 — Range Sum Query 2D Immutable (2D prefix sums)
