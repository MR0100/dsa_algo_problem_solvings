package main

import "fmt"

// ── Approach 1: Row Binary Search + Column Binary Search ─────────────────────
//
// twoStepBinarySearch solves Search a 2D Matrix with two independent binary
// searches: find the correct row, then search within that row.
//
// Intuition:
//   Since each row is sorted and the first element of each row is greater than
//   the last element of the previous row, we can binary search to find the row
//   where target could be, then binary search within that row.
//
// Time:  O(log m + log n)
// Space: O(1)
func twoStepBinarySearch(matrix [][]int, target int) bool {
	m, n := len(matrix), len(matrix[0])

	// find the row: last row where matrix[row][0] <= target
	lo, hi := 0, m-1
	row := -1
	for lo <= hi {
		mid := lo + (hi-lo)/2
		if matrix[mid][0] <= target {
			row = mid
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	if row == -1 {
		return false // target is less than the smallest element
	}

	// binary search within the found row
	lo, hi = 0, n-1
	for lo <= hi {
		mid := lo + (hi-lo)/2
		if matrix[row][mid] == target {
			return true
		} else if matrix[row][mid] < target {
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	return false
}

// ── Approach 2: Treat as Flat Sorted Array ────────────────────────────────────
//
// flatBinarySearch solves Search a 2D Matrix by treating the m×n matrix as a
// single sorted array of m*n elements and performing one binary search.
//
// Intuition:
//   The matrix has the property that if we flatten it row by row, the result is
//   a sorted array. Index `mid` in this flat array maps to:
//   row = mid / n; col = mid % n.
//
// Algorithm:
//   lo=0, hi=m*n-1
//   while lo<=hi:
//     mid=(lo+hi)/2; val=matrix[mid/n][mid%n]
//     if val==target: return true
//     elif val<target: lo=mid+1
//     else: hi=mid-1
//
// Time:  O(log(m × n)) = O(log m + log n)
// Space: O(1)
func flatBinarySearch(matrix [][]int, target int) bool {
	m, n := len(matrix), len(matrix[0])
	lo, hi := 0, m*n-1

	for lo <= hi {
		mid := lo + (hi-lo)/2
		val := matrix[mid/n][mid%n] // convert flat index to 2D
		if val == target {
			return true
		} else if val < target {
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	return false
}

func main() {
	m1 := [][]int{{1, 3, 5, 7}, {10, 11, 16, 20}, {23, 30, 34, 60}}

	fmt.Println("=== Approach 1: Two-Step Binary Search ===")
	fmt.Printf("matrix, target=3   got=%v  expected true\n", twoStepBinarySearch(m1, 3))
	fmt.Printf("matrix, target=13  got=%v  expected false\n", twoStepBinarySearch(m1, 13))
	fmt.Printf("matrix, target=1   got=%v  expected true\n", twoStepBinarySearch(m1, 1))
	fmt.Printf("matrix, target=60  got=%v  expected true\n", twoStepBinarySearch(m1, 60))

	fmt.Println("=== Approach 2: Flat Binary Search ===")
	fmt.Printf("matrix, target=3   got=%v  expected true\n", flatBinarySearch(m1, 3))
	fmt.Printf("matrix, target=13  got=%v  expected false\n", flatBinarySearch(m1, 13))
	fmt.Printf("matrix, target=1   got=%v  expected true\n", flatBinarySearch(m1, 1))
	fmt.Printf("matrix, target=60  got=%v  expected true\n", flatBinarySearch(m1, 60))
}
