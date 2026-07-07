# 0059 — Spiral Matrix II

> LeetCode #59 · Difficulty: Medium
> **Categories:** Array, Matrix, Simulation

---

## Problem Statement

Given a positive integer `n`, generate an `n × n` matrix filled with elements from `1` to `n²` in spiral order.

**Example 1**
```
Input:  n = 3
Output: [[1,2,3],[8,9,4],[7,6,5]]
```

**Example 2**
```
Input:  n = 1
Output: [[1]]
```

**Constraints**
- `1 <= n <= 20`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★☆ High      | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Meta      | ★★★☆☆ Medium    | 2023          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |
| Apple     | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Boundary Shrinking** — same as #54 Spiral Matrix; maintain top/bottom/left/right boundaries.
- **Direction Vector Simulation** — walk and turn; detect walls by bounds or by non-zero fill.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Layer-by-Layer Fill ✅ | O(n²) | O(n²) | Mirror of #54; boundary pointer approach |
| 2 | Direction Vector Simulation | O(n²) | O(n²) | Slightly simpler code; uses matrix values as visited markers |

---

## Approach 1 — Layer-by-Layer Fill (Recommended ✅)

### Intuition
Maintain four boundaries (top, bottom, left, right). Fill right across top, down right side, left across bottom, up left side. Shrink the corresponding boundary after each direction. Repeat until all n² cells are filled.

### Algorithm
```
num = 1
while num <= n²:
  right across top row; top++
  down right col; right--
  left across bottom row; bottom--
  up left col; left++
```

### Complexity
- **Time:** O(n²) — each of the n² cells written once.
- **Space:** O(n²) — the output matrix.

### Code
```go
// layerFill solves Spiral Matrix II by filling an n×n matrix in spiral order
// using the same boundary-shrinking technique as #54 (Spiral Matrix).
//
// Time:  O(n²) — each cell filled once.
// Space: O(n²) — the output matrix.
func layerFill(n int) [][]int {
	matrix := make([][]int, n)
	for i := range matrix {
		matrix[i] = make([]int, n)
	}

	num := 1
	top, bottom, left, right := 0, n-1, 0, n-1

	for num <= n*n {
		// fill right across top row
		for c := left; c <= right && num <= n*n; c++ {
			matrix[top][c] = num
			num++
		}
		top++

		// fill down right column
		for r := top; r <= bottom && num <= n*n; r++ {
			matrix[r][right] = num
			num++
		}
		right--

		// fill left across bottom row
		for c := right; c >= left && num <= n*n; c-- {
			matrix[bottom][c] = num
			num++
		}
		bottom--

		// fill up left column
		for r := bottom; r >= top && num <= n*n; r-- {
			matrix[r][left] = num
			num++
		}
		left++
	}

	return matrix
}
```

### Dry Run — `n = 3`
```
Initial: top=0, bottom=2, left=0, right=2, num=1

Iteration 1:
  Right (row 0): matrix[0][0..2] = 1,2,3. num=4. top=1.
  Down (col 2, rows 1–2): matrix[1][2]=4, matrix[2][2]=5. num=6. right=1.
  Left (row 2, cols 1–0): matrix[2][1]=6, matrix[2][0]=7. num=8. bottom=1.
  Up (col 0, rows 1–1): matrix[1][0]=8. num=9. left=1.

Iteration 2: top=1, bottom=1, left=1, right=1
  Right (row 1, col 1): matrix[1][1]=9. num=10. top=2.
  10 > 9 → stop.

Result: [[1,2,3],[8,9,4],[7,6,5]] ✓
```

---

## Approach 2 — Direction Vector Simulation

### Intuition
Walk one step at a time. If the next cell is out of bounds or already non-zero (already filled), rotate 90° clockwise. Fill each cell with the current number.

Zero serves as the "unvisited" marker since valid fill values start at 1.

### Complexity
- **Time:** O(n²).
- **Space:** O(n²).

### Code
```go
// simulation solves Spiral Matrix II by walking with a direction vector and
// turning right when blocked or out of bounds.
//
// Time:  O(n²)
// Space: O(n²)
func simulation(n int) [][]int {
	matrix := make([][]int, n)
	for i := range matrix {
		matrix[i] = make([]int, n)
	}

	dr := []int{0, 1, 0, -1} // right, down, left, up
	dc := []int{1, 0, -1, 0}
	dir := 0
	r, c := 0, 0

	for num := 1; num <= n*n; num++ {
		matrix[r][c] = num
		nr, nc := r+dr[dir], c+dc[dir]
		if nr < 0 || nr >= n || nc < 0 || nc >= n || matrix[nr][nc] != 0 {
			dir = (dir + 1) % 4 // turn right
			nr, nc = r+dr[dir], c+dc[dir]
		}
		r, c = nr, nc
	}

	return matrix
}
```

### Dry Run — `n = 3`

Directions cycle right(0)→down(1)→left(2)→up(3). Turn when the next cell is out of bounds or already non-zero.

| num | write (r,c) | value | next (nr,nc) via dir | blocked? | new dir | move to (r,c) |
|-----|-------------|-------|----------------------|----------|---------|---------------|
| 1 | (0,0) | 1 | dir0 →(0,1) | no | 0 | (0,1) |
| 2 | (0,1) | 2 | dir0 →(0,2) | no | 0 | (0,2) |
| 3 | (0,2) | 3 | dir0 →(0,3) oob | yes | 1 (down) | (1,2) |
| 4 | (1,2) | 4 | dir1 →(2,2) | no | 1 | (2,2) |
| 5 | (2,2) | 5 | dir1 →(3,2) oob | yes | 2 (left) | (2,1) |
| 6 | (2,1) | 6 | dir2 →(2,0) | no | 2 | (2,0) |
| 7 | (2,0) | 7 | dir2 →(2,-1) oob | yes | 3 (up) | (1,0) |
| 8 | (1,0) | 8 | dir3 →(0,0) filled | yes | 0 (right) | (1,1) |
| 9 | (1,1) | 9 | loop ends (num=9=n²) | — | — | — |

Result: `[[1,2,3],[8,9,4],[7,6,5]]` ✓

---

## Key Takeaways

- **This is #54 in reverse** — #54 reads from a matrix in spiral order; #59 writes to a matrix in spiral order. The same layer-peel / direction-simulation frameworks apply identically.
- **Direction simulation is slightly cleaner here** — because zeros serve as a natural "unvisited" marker, no separate visited array is needed (unlike #54 where elements could be zero).
- **Layer fill avoids boundary guard complexity** — the `while num <= n²` condition eliminates the need for `if top<=bottom` guards that appear in #54's layer peel.

---

## Related Problems

- LeetCode #54 — Spiral Matrix (read matrix in spiral order)
- LeetCode #48 — Rotate Image (in-place matrix transformation)
