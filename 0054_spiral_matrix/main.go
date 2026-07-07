package main

import "fmt"

// ── Approach 1: Layer-by-Layer Peeling ───────────────────────────────────────
//
// layerPeel solves Spiral Matrix by peeling the outermost ring at each step.
//
// Intuition:
//   Maintain four boundaries: top, bottom, left, right.
//   Each iteration: traverse right across top row, down right col, left across
//   bottom row, up left col. After each traversal, shrink the corresponding
//   boundary inward.
//
// Algorithm:
//   top=0, bottom=m-1, left=0, right=n-1
//   while top<=bottom and left<=right:
//     right across top row; top++
//     down right col; right--
//     if top<=bottom: left across bottom row; bottom--
//     if left<=right: up left col; left++
//
// Time:  O(m × n) — every element visited once.
// Space: O(1)     — no extra data structures (output slice not counted).
func layerPeel(matrix [][]int) []int {
	if len(matrix) == 0 {
		return nil
	}
	m, n := len(matrix), len(matrix[0])
	result := make([]int, 0, m*n)
	top, bottom, left, right := 0, m-1, 0, n-1

	for top <= bottom && left <= right {
		// traverse right across top row
		for c := left; c <= right; c++ {
			result = append(result, matrix[top][c])
		}
		top++

		// traverse down right column
		for r := top; r <= bottom; r++ {
			result = append(result, matrix[r][right])
		}
		right--

		// traverse left across bottom row (guard: may have collapsed)
		if top <= bottom {
			for c := right; c >= left; c-- {
				result = append(result, matrix[bottom][c])
			}
			bottom--
		}

		// traverse up left column (guard: may have collapsed)
		if left <= right {
			for r := bottom; r >= top; r-- {
				result = append(result, matrix[r][left])
			}
			left++
		}
	}

	return result
}

// ── Approach 2: Direction-Vector Simulation ───────────────────────────────────
//
// simulation solves Spiral Matrix by walking in the current direction and
// turning right when the next step would go out of bounds or revisit a cell.
//
// Intuition:
//   Follow the spiral path step by step. At each cell, mark it visited.
//   Try to continue in the current direction; if blocked, rotate 90° clockwise.
//   Direction order: right → down → left → up → right → ...
//
// Algorithm:
//   dirs = [(0,1),(1,0),(0,-1),(-1,0)]; dir=0
//   r,c = 0,0; visited[r][c] = true
//   for i=0 to m*n-1:
//     record matrix[r][c]
//     nr,nc = r+dirs[dir][0], c+dirs[dir][1]
//     if (nr,nc) out of bounds or visited: dir = (dir+1)%4; nr,nc = r+dr,c+dc
//     r,c = nr,nc
//
// Time:  O(m × n)
// Space: O(m × n) — visited array.
func simulation(matrix [][]int) []int {
	if len(matrix) == 0 {
		return nil
	}
	m, n := len(matrix), len(matrix[0])
	result := make([]int, 0, m*n)
	visited := make([][]bool, m)
	for i := range visited {
		visited[i] = make([]bool, n)
	}

	// right, down, left, up
	dr := []int{0, 1, 0, -1}
	dc := []int{1, 0, -1, 0}
	dir := 0
	r, c := 0, 0

	for i := 0; i < m*n; i++ {
		result = append(result, matrix[r][c])
		visited[r][c] = true

		nr, nc := r+dr[dir], c+dc[dir]
		if nr < 0 || nr >= m || nc < 0 || nc >= n || visited[nr][nc] {
			dir = (dir + 1) % 4 // turn right
			nr, nc = r+dr[dir], c+dc[dir]
		}
		r, c = nr, nc
	}

	return result
}

func main() {
	m1 := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	m2 := [][]int{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}}

	fmt.Println("=== Approach 1: Layer Peel ===")
	fmt.Printf("matrix=[[1,2,3],[4,5,6],[7,8,9]]  got=%v  expected [1 2 3 6 9 8 7 4 5]\n", layerPeel(m1))
	fmt.Printf("matrix=4x3  got=%v  expected [1 2 3 4 8 12 11 10 9 5 6 7]\n", layerPeel(m2))

	fmt.Println("=== Approach 2: Direction Simulation ===")
	fmt.Printf("matrix=[[1,2,3],[4,5,6],[7,8,9]]  got=%v  expected [1 2 3 6 9 8 7 4 5]\n", simulation(m1))
	fmt.Printf("matrix=4x3  got=%v  expected [1 2 3 4 8 12 11 10 9 5 6 7]\n", simulation(m2))
}
