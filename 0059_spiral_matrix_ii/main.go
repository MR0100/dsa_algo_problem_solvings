package main

import "fmt"

// ── Approach 1: Layer-by-Layer Fill ───────────────────────────────────────────
//
// layerFill solves Spiral Matrix II by filling an n×n matrix in spiral order
// using the same boundary-shrinking technique as #54 (Spiral Matrix).
//
// Intuition:
//   Maintain four boundaries (top, bottom, left, right). Fill right across top,
//   down right col, left across bottom, up left col; shrink boundaries inward.
//   Repeat until num > n².
//
// Algorithm:
//   matrix = n×n zeros; num = 1
//   while num <= n²:
//     fill right across top; top++
//     fill down right; right--
//     fill left across bottom; bottom--
//     fill up left; left++
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

// ── Approach 2: Direction Vector Simulation ───────────────────────────────────
//
// simulation solves Spiral Matrix II by walking with a direction vector and
// turning right when blocked or out of bounds.
//
// Intuition:
//   Walk in the current direction; if the next cell is out of bounds or already
//   filled, rotate 90° clockwise. Fill each cell with the current number.
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

func printMatrix(m [][]int) {
	for _, row := range m {
		fmt.Println(row)
	}
}

func main() {
	fmt.Println("=== Approach 1: Layer Fill ===")
	fmt.Println("n=3:")
	printMatrix(layerFill(3))
	fmt.Println("n=1:")
	printMatrix(layerFill(1))

	fmt.Println("=== Approach 2: Simulation ===")
	fmt.Println("n=3:")
	printMatrix(simulation(3))
	fmt.Println("n=1:")
	printMatrix(simulation(1))
}
