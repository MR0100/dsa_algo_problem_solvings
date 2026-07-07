package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Brute Force (Enumerate Every Rectangle) ──────────────────────
//
// bruteForce solves Max Sum of Rectangle No Larger Than K by trying every pair
// of corners, summing the enclosed rectangle, and keeping the best sum ≤ k.
//
// Intuition:
//
//	A rectangle is defined by a top-left (r1,c1) and bottom-right (r2,c2).
//	Enumerate all such pairs, sum the cells inside, and track the largest sum
//	that does not exceed k. Correct but very slow — four nested corner loops
//	plus the interior sum.
//
// Algorithm:
//  1. For every (r1,c1) as top-left and (r2,c2) as bottom-right:
//     a. Sum all cells in that rectangle.
//     b. If the sum ≤ k, update best.
//  2. Return best.
//
// Time:  O(m² · n² · m · n) = O(m³ · n³) naively (interior re-summed each time).
// Space: O(1).
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

// ── Approach 2: Column-Pair Compression + Prefix Sums ────────────────────────
//
// prefixSumRowSearch fixes a pair of top/bottom rows, compresses the band into a
// 1-D array of column sums, then finds the best subarray with sum ≤ k by
// checking every subarray via running prefix sums.
//
// Intuition:
//
//	Fixing the top row r1 and bottom row r2 collapses the 2-D problem to 1-D:
//	rowSum[c] = sum of column c between rows r1..r2. Now we want the maximum
//	contiguous subarray of rowSum whose sum is ≤ k — a classic 1-D task.
//	Accumulating column sums incrementally as r2 grows avoids recomputing.
//
// Algorithm:
//  1. For each top row r1:
//     a. Reset colSum[] to zeros.
//     b. For each bottom row r2 ≥ r1: add row r2 into colSum (band sums).
//     c. Find the best subarray of colSum with sum ≤ k (O(n²) here) and
//     update the global best.
//  2. Return best.
//
// Time:  O(m² · n²) — O(m²) row pairs, each doing an O(n²) subarray scan.
// Space: O(n) — the compressed column-sum array.
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

// ── Approach 3: Column Compression + Sorted-Prefix Binary Search (Optimal) ───
//
// sortedPrefixBinarySearch keeps the row-band compression but replaces the inner
// O(n²) subarray scan with an O(n log n) search: for each running prefix, find
// the smallest earlier prefix that is ≥ prefix-k using an ordered set.
//
// Intuition:
//
//	After compressing rows r1..r2 into colSum, a subarray sum equals
//	prefix[j] - prefix[i]. We want the largest such value that is ≤ k, i.e.
//	for the current prefix P we want the SMALLEST earlier prefix ≥ P - k.
//	Maintain earlier prefixes in a sorted structure and binary-search for
//	that lower bound; each candidate is P - (found prefix). Always insert the
//	initial prefix 0 so single-cell-from-start subarrays are considered.
//
// Algorithm:
//  1. For each top row r1: reset colSum.
//  2. For each bottom row r2: extend colSum by row r2.
//  3. Run a prefix sum over colSum; keep a sorted slice `seen` of prior
//     prefixes (seeded with 0). For each new prefix P:
//     a. Binary-search `seen` for the smallest value ≥ P - k.
//     b. If found (call it lo), candidate = P - lo ≤ k; update best.
//     c. Insert P into `seen` keeping it sorted.
//  4. Return best.
//
// Time:  O(m² · n log n) — O(m²) row pairs, each an O(n log n) prefix search.
//
//	When m > n, transpose first so the smaller dimension is squared.
//
// Space: O(n) — colSum plus the sorted prefix set.
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

func main() {
	// Example 1: matrix=[[1,0,1],[0,-2,3]], k=2 → 2 (rectangle [[0,1],[-2,3]] sums to 2).
	m1 := [][]int{{1, 0, 1}, {0, -2, 3}}
	// Example 2: matrix=[[2,2,-1]], k=3 → 3 (the subarray [2,-1]... actually 2+2-1=3).
	m2 := [][]int{{2, 2, -1}}

	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(bruteForce(m1, 2)) // expected 2
	fmt.Println(bruteForce(m2, 3)) // expected 3

	fmt.Println("=== Approach 2: Column Compression + Prefix Sums ===")
	fmt.Println(prefixSumRowSearch(m1, 2)) // expected 2
	fmt.Println(prefixSumRowSearch(m2, 3)) // expected 3

	fmt.Println("=== Approach 3: Compression + Sorted-Prefix Binary Search (Optimal) ===")
	fmt.Println(sortedPrefixBinarySearch(m1, 2)) // expected 2
	fmt.Println(sortedPrefixBinarySearch(m2, 3)) // expected 3
}
