package main

import "fmt"

// ── Approach 1: Iterative Row-by-Row ─────────────────────────────────────────
//
// generate solves Pascal's Triangle iteratively.
//
// Intuition:
//   Each row starts and ends with 1. Interior elements are the sum of the two
//   elements directly above: triangle[i][j] = triangle[i-1][j-1] + triangle[i-1][j].
//   Build each row from the previous row.
//
// Time:  O(numRows^2) — n(n+1)/2 elements total.
// Space: O(numRows^2) — all elements stored.
func generate(numRows int) [][]int {
	result := make([][]int, numRows)

	for i := 0; i < numRows; i++ {
		row := make([]int, i+1)
		row[0] = 1      // first element always 1
		row[i] = 1      // last element always 1
		for j := 1; j < i; j++ {
			row[j] = result[i-1][j-1] + result[i-1][j]
		}
		result[i] = row
	}
	return result
}

func main() {
	fmt.Println("=== Approach 1: Iterative ===")
	fmt.Printf("numRows=5  got=%v\n  expected [[1] [1 1] [1 2 1] [1 3 3 1] [1 4 6 4 1]]\n", generate(5))
	fmt.Printf("numRows=1  got=%v  expected [[1]]\n", generate(1))
}
