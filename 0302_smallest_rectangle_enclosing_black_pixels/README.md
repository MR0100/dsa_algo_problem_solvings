# 0302 — Smallest Rectangle Enclosing Black Pixels

> LeetCode #302 · Difficulty: Hard (Premium / Locked)
> **Categories:** Array, Binary Search, Matrix

---

## Problem Statement

You are given an `image` that is represented by an `m x n` binary matrix, where `image[i][j]` is `'0'` if it represents a white pixel and `'1'` if it represents a black pixel.

The black pixels are connected (i.e., there is only one black region). Pixels are connected horizontally and vertically.

Given two integers `x` and `y` that represent the location of one of the black pixels, return _the area of the smallest (axis-aligned) rectangle that encloses all black pixels_.

You must write an algorithm with less than `O(m * n)` runtime complexity.

**Example 1:**

```
Input:  image = ["0010","0110","0100"], x = 0, y = 2
Output: 6
```

Explanation: the black pixels occupy rows 0..2 and columns 1..2, so the enclosing rectangle is 3 rows × 2 columns = 6.

**Example 2:**

```
Input:  image = ["1"], x = 0, y = 0
Output: 1
```

**Constraints:**

- `m == image.length`
- `n == image[i].length`
- `1 <= m, n <= 100`
- `image[i][j]` is either `'0'` or `'1'`.
- `0 <= x < m`
- `0 <= y < n`
- `image[x][y] == '1'`.
- The black pixels in the `image` only form one component.

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Google    | ★★★★☆ High       | 2023          |
| Facebook  | ★★★☆☆ Medium     | 2022          |
| Amazon    | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Binary Search on a monotone predicate** — the "row/column contains a black pixel" predicate is monotone around the connected region, enabling boundary search → see [`/dsa/binary_search.md`](/dsa/binary_search.md)
- **Matrix traversal** — projecting the 2D region onto row and column axes → see [`/dsa/matrix_traversal.md`](/dsa/matrix_traversal.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute-force scan | O(m·n) | O(1) | Baseline; simplest to verify but fails the sub-O(m·n) follow-up |
| 2 | Binary search on boundaries (Optimal) | O(m·log n + n·log m) | O(1) | Meets the follow-up; exploits connectivity for a monotone predicate |

---

## Approach 1 — Brute-Force Scan

### Intuition
The smallest enclosing rectangle is fully determined by the extreme black pixels: the minimum and maximum black rows and columns. Scanning every cell and relaxing four running extremes finds them directly.

### Algorithm
1. Initialise `minRow/minCol` high and `maxRow/maxCol` low (inverted).
2. For each `'1'` cell, update the four extremes.
3. If no black pixel was seen, return 0.
4. Return `(maxRow − minRow + 1) · (maxCol − minCol + 1)`.

### Complexity
- **Time:** O(m·n) — every pixel is visited once.
- **Space:** O(1) — four scalar extremes.

### Code
```go
func bruteForce(image []string, x, y int) int {
	if len(image) == 0 || len(image[0]) == 0 {
		return 0
	}
	minRow, maxRow := len(image), -1     // row bounds (init inverted)
	minCol, maxCol := len(image[0]), -1  // col bounds (init inverted)
	for r := 0; r < len(image); r++ {
		for c := 0; c < len(image[0]); c++ {
			if image[r][c] == '1' { // a black pixel widens the bounds
				if r < minRow {
					minRow = r
				}
				if r > maxRow {
					maxRow = r
				}
				if c < minCol {
					minCol = c
				}
				if c > maxCol {
					maxCol = c
				}
			}
		}
	}
	if maxRow == -1 { // no black pixel at all
		return 0
	}
	return (maxRow - minRow + 1) * (maxCol - minCol + 1)
}
```

### Dry Run
Input `image = ["0010","0110","0100"]`, `x=0, y=2`.

| Cell visited | Value | minRow | maxRow | minCol | maxCol |
|---|---|---|---|---|---|
| (0,2) | 1 | 0 | 0 | 2 | 2 |
| (1,1) | 1 | 0 | 1 | 1 | 2 |
| (1,2) | 1 | 0 | 1 | 1 | 2 |
| (2,1) | 1 | 0 | 2 | 1 | 2 |

Area = `(2−0+1)·(2−1+1)` = `3·2` = **6**.

---

## Approach 2 — Binary Search on Boundaries (Optimal)

### Intuition
Because all black pixels form one connected region, the set of rows containing a black pixel is a **contiguous interval** — there are no empty rows between the top and bottom black rows. So "row `r` has a black pixel" is a monotone predicate: false above the region, true inside, false below. Binary search finds each of the four edges. Starting from the known black pixel `(x, y)`, search the top edge in `[0, x]`, the bottom edge in `[x+1, m]`, and the left/right edges in the column axis similarly. Each probe scans one row (O(n)) or column (O(m)).

### Algorithm
1. Write helpers `hasBlackInRow(r)` and `hasBlackInCol(c)`.
2. Binary-search the smallest black row `top` in `[0, x)` (predicate = has black).
3. Binary-search the first non-black row `bottom` in `[x+1, m)` (predicate = has no black).
4. Binary-search `left` in `[0, y)` and `right` in `[y+1, n)` on columns the same way.
5. `top`/`left` are inclusive, `bottom`/`right` are exclusive; area = `(bottom − top)·(right − left)`.

### Complexity
- **Time:** O(m·log n + n·log m) — each search does O(log dim) probes, each probe scanning a full row/column.
- **Space:** O(1) — only index variables.

### Code
```go
func binarySearch(image []string, x, y int) int {
	if len(image) == 0 || len(image[0]) == 0 {
		return 0
	}
	m, n := len(image), len(image[0])

	hasBlackInRow := func(r int) bool { // any '1' in row r?
		for c := 0; c < n; c++ {
			if image[r][c] == '1' {
				return true
			}
		}
		return false
	}
	hasBlackInCol := func(c int) bool { // any '1' in column c?
		for r := 0; r < m; r++ {
			if image[r][c] == '1' {
				return true
			}
		}
		return false
	}

	searchRows := func(lo, hi int, findFirstBlack bool) int {
		for lo < hi {
			mid := (lo + hi) / 2
			if hasBlackInRow(mid) == findFirstBlack {
				hi = mid // condition met → boundary is at or above mid
			} else {
				lo = mid + 1 // condition not met → search below
			}
		}
		return lo
	}
	searchCols := func(lo, hi int, findFirstBlack bool) int {
		for lo < hi {
			mid := (lo + hi) / 2
			if hasBlackInCol(mid) == findFirstBlack {
				hi = mid
			} else {
				lo = mid + 1
			}
		}
		return lo
	}

	top := searchRows(0, x, true)       // first black row in [0, x)
	bottom := searchRows(x+1, m, false) // first NON-black row in [x+1, m)
	left := searchCols(0, y, true)      // first black col in [0, y)
	right := searchCols(y+1, n, false)  // first NON-black col in [y+1, n)

	return (bottom - top) * (right - left)
}
```

### Dry Run
Input `image = ["0010","0110","0100"]`, `x=0, y=2` (m=3, n=4).

| Search | Range | Result | Meaning |
|---|---|---|---|
| `top` = searchRows(0,0,true) | empty | 0 | row 0 already black (loop skipped) |
| `bottom` = searchRows(1,3,false) | rows 1,2 both black; row 3 out | 3 | first empty row index past region |
| `left` = searchCols(0,2,true) | col 0 empty, col 1 black | 1 | leftmost black column |
| `right` = searchCols(3,4,false) | col 3 empty | 3 | first empty col past region |

Area = `(bottom − top)·(right − left)` = `(3−0)·(3−1)` = `3·2` = **6**.

---

## Key Takeaways

- **Connectivity ⇒ contiguous projection.** A single connected region projects to a gap-free interval on each axis, which makes "does this line contain the region" a monotone predicate — the precondition for binary search.
- Search **four boundaries independently**; the known interior point `(x, y)` splits each axis into a "find first true" half and a "find first false" half.
- Keep boundaries **half-open** (`top`/`left` inclusive, `bottom`/`right` exclusive) so the area is a clean subtraction with no off-by-one.
- This is a classic case where an interview's O(m·n) baseline is upgraded to sub-linear-in-area by leveraging a structural guarantee.

---

## Related Problems

- LeetCode #74 — Search a 2D Matrix (binary search over a matrix)
- LeetCode #240 — Search a 2D Matrix II (monotone 2D search)
- LeetCode #200 — Number of Islands (connected black region, flood fill)
- LeetCode #35 — Search Insert Position (find-first-true binary search template)
