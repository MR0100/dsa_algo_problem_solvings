package main

import "fmt"

// ── Approach 1: In-Place Row Update (O(k) Space) ─────────────────────────────
//
// getRow solves Pascal's Triangle II returning only the rowIndex-th row.
//
// Intuition:
//   Maintain a single row and update it in-place. To compute row i from row i-1,
//   walk from right to left: row[j] += row[j-1].
//   Walking right-to-left ensures row[j-1] used is still from the previous row.
//
// Time:  O(rowIndex^2)
// Space: O(rowIndex) — single row.
func getRow(rowIndex int) []int {
	row := make([]int, rowIndex+1)
	row[0] = 1

	for i := 1; i <= rowIndex; i++ {
		// walk right to left to avoid overwriting values we still need
		for j := i; j >= 1; j-- {
			row[j] += row[j-1]
		}
	}
	return row
}

// ── Approach 2: Combinatorial Formula ────────────────────────────────────────
//
// getRowCombinatorial solves Pascal's Triangle II using C(n,k) formula.
//
// Intuition:
//   The k-th element of row n is C(n,k) = n!/(k!*(n-k)!).
//   Use the recurrence C(n,k) = C(n,k-1) * (n-k+1) / k to avoid large
//   intermediate factorials.
//
// Time:  O(rowIndex)
// Space: O(rowIndex)
func getRowCombinatorial(rowIndex int) []int {
	row := make([]int, rowIndex+1)
	row[0] = 1
	for k := 1; k <= rowIndex; k++ {
		// C(n,k) = C(n,k-1) * (n-k+1) / k
		row[k] = row[k-1] * (rowIndex - k + 1) / k
	}
	return row
}

func main() {
	fmt.Println("=== Approach 1: In-Place Update ===")
	fmt.Printf("rowIndex=3  got=%v  expected [1 3 3 1]\n", getRow(3))
	fmt.Printf("rowIndex=0  got=%v  expected [1]\n", getRow(0))
	fmt.Printf("rowIndex=1  got=%v  expected [1 1]\n", getRow(1))

	fmt.Println("=== Approach 2: Combinatorial Formula ===")
	fmt.Printf("rowIndex=3  got=%v  expected [1 3 3 1]\n", getRowCombinatorial(3))
	fmt.Printf("rowIndex=4  got=%v  expected [1 4 6 4 1]\n", getRowCombinatorial(4))
}
