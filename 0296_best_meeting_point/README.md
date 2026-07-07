# 0296 — Best Meeting Point

> LeetCode #296 · Difficulty: Hard
> **Categories:** Array, Math, Sorting, Matrix

---

## Problem Statement

Given an `m x n` binary grid `grid` where each `1` marks the home of one friend, return *the minimal total travel distance*.

The total travel distance is the sum of the distances between the houses of the friends and the meeting point.

The distance is calculated using **Manhattan Distance**, where `distance(p1, p2) = |p2.x - p1.x| + |p2.y - p1.y|`.

**Example 1:**

```
Input: grid = [[1,0,0,0,1],[0,0,0,0,0],[0,0,1,0,0]]
Output: 6
Explanation: Given three friends living at (0,0), (0,4), and (2,2).
The point (0,2) is an ideal meeting point, as the total travel distance
of 2 + 2 + 2 = 6 is minimal.
So return 6.
```

**Example 2:**

```
Input: grid = [[1,1]]
Output: 1
```

**Constraints:**

- `m == grid.length`
- `n == grid[i].length`
- `1 <= m, n <= 200`
- `grid[i][j]` is either `0` or `1`.
- There will be **at least two** friends in the `grid`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2023          |
| Facebook   | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★★☆☆ Medium     | 2022          |
| Twitter    | ★★☆☆☆ Low        | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Median minimises sum of absolute distances** — the key math fact that makes this optimal → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)
- **Sorting** — to locate the median of the coordinates → see [`/dsa/sorting.md`](/dsa/sorting.md)
- **Two Pointers** — pairing outer extremes to sum median distances without an explicit median lookup → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **Matrix Traversal** — row-major and column-major scans to gather coordinates in sorted order → see [`/dsa/matrix_traversal.md`](/dsa/matrix_traversal.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (try every cell) | O(m·n·k) | O(k) | Tiny grids; clarity over speed |
| 2 | Median via Sorting | O(k log k) | O(k) | Clean, idiomatic optimal |
| 3 | Two-Pointer (row/col scan) | O(m·n) | O(k) | Avoid sorting entirely |

*(k = number of homes.)*

---

## Approach 1 — Brute Force

### Intuition
The meeting point can be any of the `m·n` cells. Try them all: for each candidate, add up the Manhattan distance to every home and keep the smallest total.

### Algorithm
1. Collect every home coordinate into a list.
2. For every cell `(r, c)`, sum `|r-hr| + |c-hc|` over all homes.
3. Track the minimum total and return it.

### Complexity
- **Time:** O(m·n·k) — each of `m·n` candidate cells is scored against all `k` homes.
- **Space:** O(k) — the list of home coordinates.

### Code
```go
func bruteForce(grid [][]int) int {
	m := len(grid)
	if m == 0 {
		return 0
	}
	n := len(grid[0])

	type point struct{ r, c int }
	homes := []point{}
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if grid[r][c] == 1 {
				homes = append(homes, point{r, c})
			}
		}
	}

	best := -1
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			total := 0
			for _, h := range homes {
				total += abs(r-h.r) + abs(c-h.c)
			}
			if best == -1 || total < best {
				best = total
			}
		}
	}
	if best == -1 {
		return 0
	}
	return best
}
```

### Dry Run
Example 1, homes at `(0,0)`, `(0,4)`, `(2,2)`. A few candidate cells:

| Candidate (r,c) | dist to (0,0) | dist to (0,4) | dist to (2,2) | total | best so far |
|-----------------|---------------|---------------|---------------|-------|-------------|
| (0,0)           | 0             | 4             | 4             | 8     | 8           |
| (0,2)           | 2             | 2             | 2             | 6     | 6           |
| (1,2)           | 3             | 3             | 1             | 7     | 6           |
| (2,2)           | 4             | 4             | 0             | 8     | 6           |

Minimum found = **6**.

---

## Approach 2 — Median via Sorting

### Intuition
Manhattan distance splits into independent x and y parts:
`total = Σ|r - hr| + Σ|c - hc|`. Minimise each 1-D sum separately. In one dimension, the point minimising the sum of absolute distances is the **median** of the points.

### Algorithm
1. Gather `rows[]` and `cols[]` of every home.
2. Sort both arrays.
3. Return `minDistance1D(rows) + minDistance1D(cols)`, where the helper walks two pointers inward summing outer gaps (equal to the sum of distances to the median).

### Complexity
- **Time:** O(k log k) — dominated by sorting the two coordinate lists.
- **Space:** O(k) — the coordinate lists.

### Code
```go
func medianSort(grid [][]int) int {
	rows := []int{}
	cols := []int{}
	for r := 0; r < len(grid); r++ {
		for c := 0; c < len(grid[0]); c++ {
			if grid[r][c] == 1 {
				rows = append(rows, r)
				cols = append(cols, c)
			}
		}
	}
	sort.Ints(rows)
	sort.Ints(cols)
	return minDistance1D(rows) + minDistance1D(cols)
}

func minDistance1D(sorted []int) int {
	total := 0
	i, j := 0, len(sorted)-1
	for i < j {
		total += sorted[j] - sorted[i]
		i++
		j--
	}
	return total
}
```

### Dry Run
Example 1: homes `(0,0)`, `(0,4)`, `(2,2)`.

| step | rows (sorted) | cols (sorted) |
|------|---------------|---------------|
| collect | [0, 0, 2] | [0, 4, 2] |
| sort    | [0, 0, 2] | [0, 2, 4] |

`minDistance1D(rows)`: pair `(0,2)` → gap `2`; pointers meet → total `2`.
`minDistance1D(cols)`: pair `(0,4)` → gap `4`; pointers meet → total `4`.
Answer = `2 + 4` = **6**.

---

## Approach 3 — Two-Pointer (Row/Col Scan) (Optimal)

### Intuition
Sorting isn't even necessary: a **row-major** scan naturally yields row indices in ascending order, and a **column-major** scan yields column indices in ascending order. Then the same two-pointer "sum of outer gaps" trick gives the median distance directly.

### Algorithm
1. Row scan (r outer, c inner): append `r` for each home → `rows` already sorted.
2. Column scan (c outer, r inner): append `c` for each home → `cols` already sorted.
3. Return `twoPointerGap(rows) + twoPointerGap(cols)`.

### Complexity
- **Time:** O(m·n) — two full grid scans, no sorting step.
- **Space:** O(k) — the two coordinate lists.

### Code
```go
func twoPointerOptimal(grid [][]int) int {
	m := len(grid)
	if m == 0 {
		return 0
	}
	n := len(grid[0])

	rows := []int{}
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if grid[r][c] == 1 {
				rows = append(rows, r)
			}
		}
	}

	cols := []int{}
	for c := 0; c < n; c++ {
		for r := 0; r < m; r++ {
			if grid[r][c] == 1 {
				cols = append(cols, c)
			}
		}
	}
	return twoPointerGap(rows) + twoPointerGap(cols)
}

func twoPointerGap(sorted []int) int {
	total := 0
	i, j := 0, len(sorted)-1
	for i < j {
		total += sorted[j] - sorted[i]
		i++
		j--
	}
	return total
}
```

### Dry Run
Example 1:

| scan | order visited | collected |
|------|---------------|-----------|
| row-major  | (0,0),(0,4),(2,2) | rows = [0, 0, 2] (ascending) |
| col-major  | (0,0),(2,2),(0,4) | cols = [0, 2, 4] (ascending) |

`twoPointerGap(rows)` = `2 - 0` = `2`.
`twoPointerGap(cols)` = `4 - 0` = `4`.
Answer = `2 + 4` = **6**.

---

## Key Takeaways
- **Manhattan distance decomposes** into independent 1-D problems on x and y.
- The **median** minimises the sum of absolute deviations in 1-D (not the mean — that minimises squared distance).
- The sum of distances to the median equals the **sum of gaps between symmetric outer pairs** — no need to locate the median value itself; a two-pointer sweep of the sorted list gives it.
- You can obtain sorted coordinates "for free" by choosing the traversal order (row-major for rows, column-major for columns).

---

## Related Problems
- LeetCode #462 — Minimum Moves to Equal Array Elements II (median in 1-D)
- LeetCode #317 — Shortest Distance from All Buildings (BFS variant of a meeting point)
- LeetCode #64 — Minimum Path Sum (grid distance, but constrained paths)
