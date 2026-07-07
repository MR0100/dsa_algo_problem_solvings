package main

import (
	"fmt"
	"math"
)

// ── Approach 1: Brute Force (Simulation) ─────────────────────────────────────
//
// bruteForce solves Arranging Coins by literally laying down rows one at a time.
//
// Intuition:
//
//	Building the staircase is the definition of the answer. Row k costs k coins.
//	Keep subtracting the next row's cost from the remaining coins; the last row
//	we can fully afford is the number of complete rows.
//
// Algorithm:
//  1. row = 0, remaining = n.
//  2. Repeatedly try the next row: its cost is row+1. If remaining >= row+1,
//     spend it (remaining -= row+1) and increment row.
//  3. Stop when we can no longer afford the next row; return row.
//
// Time:  O(√n) — the number of complete rows k satisfies k(k+1)/2 ≤ n, so
//
//	k ≈ √(2n); the loop runs that many times.
//
// Space: O(1) — two counters.
func bruteForce(n int) int {
	row := 0       // number of complete rows built so far
	remaining := n // coins left in the pile
	// Try to build row (row+1); its cost equals its index.
	for remaining >= row+1 {
		row++            // this row is complete
		remaining -= row // pay for the row we just completed
	}
	return row // last fully-built row count
}

// ── Approach 2: Binary Search ────────────────────────────────────────────────
//
// binarySearch solves Arranging Coins by searching for the largest k with
// k(k+1)/2 ≤ n.
//
// Intuition:
//
//	The total coins needed for k complete rows is the triangular number
//	T(k) = k(k+1)/2, which is monotonically increasing in k. We want the
//	largest k such that T(k) ≤ n — a classic "last true" binary search on a
//	monotone predicate.
//
// Algorithm:
//  1. lo = 1, hi = n (k can never exceed n since T(k) grows quadratically).
//  2. While lo <= hi: mid = (lo+hi)/2, curr = mid*(mid+1)/2.
//  3. If curr == n return mid. If curr < n, mid rows fit — record and search
//     right (lo = mid+1). Else search left (hi = mid-1).
//  4. Return hi — the largest k whose triangular number did not exceed n.
//
// Time:  O(log n) — halve the search space each step.
// Space: O(1) — a few 64-bit scalars.
func binarySearch(n int) int {
	lo, hi := 1, n // k lies in [1, n]
	for lo <= hi {
		mid := lo + (hi-lo)/2 // candidate row count (avoid lo+hi overflow)
		// Triangular number T(mid); use int64 to dodge overflow near n=2^31-1.
		curr := int64(mid) * int64(mid+1) / 2
		target := int64(n)
		switch {
		case curr == target:
			return mid // exact fit: mid complete rows, none left over
		case curr < target:
			lo = mid + 1 // mid rows fit with coins to spare — try more
		default:
			hi = mid - 1 // mid rows cost too much — try fewer
		}
	}
	// Loop exits with hi = largest k where T(k) < n < T(k+1); hi is the answer.
	return hi
}

// ── Approach 3: Quadratic Formula (Optimal) ──────────────────────────────────
//
// mathFormula solves Arranging Coins by solving the triangular-number equation
// directly.
//
// Intuition:
//
//	We need the largest k with k(k+1)/2 ≤ n, i.e. k² + k − 2n ≤ 0. Solving the
//	quadratic k² + k − 2n = 0 gives the positive root k = (−1 + √(1 + 8n)) / 2.
//	The floor of that root is exactly the number of complete rows.
//
// Algorithm:
//  1. Compute k = (−1 + √(1 + 8n)) / 2 in floating point.
//  2. Return ⌊k⌋.
//
// Time:  O(1) — one square root.
// Space: O(1).
//
// Note: 8n can reach ~1.7×10¹⁰, well within float64's exact-integer range
// (2⁵³), so the sqrt is precise enough here; floor gives the correct answer.
func mathFormula(n int) int {
	// Positive root of k² + k − 2n = 0, then floored.
	return int((-1 + math.Sqrt(1+8*float64(n))) / 2)
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Simulation) ===")
	fmt.Printf("n=5           got=%d  expected 2\n", bruteForce(5))
	fmt.Printf("n=8           got=%d  expected 3\n", bruteForce(8))
	fmt.Printf("n=1           got=%d  expected 1\n", bruteForce(1))

	fmt.Println("=== Approach 2: Binary Search ===")
	fmt.Printf("n=5           got=%d  expected 2\n", binarySearch(5))
	fmt.Printf("n=8           got=%d  expected 3\n", binarySearch(8))
	fmt.Printf("n=1           got=%d  expected 1\n", binarySearch(1))
	fmt.Printf("n=2147483647  got=%d  expected 65535\n", binarySearch(2147483647)) // max-input edge

	fmt.Println("=== Approach 3: Quadratic Formula (Optimal) ===")
	fmt.Printf("n=5           got=%d  expected 2\n", mathFormula(5))
	fmt.Printf("n=8           got=%d  expected 3\n", mathFormula(8))
	fmt.Printf("n=1           got=%d  expected 1\n", mathFormula(1))
	fmt.Printf("n=2147483647  got=%d  expected 65535\n", mathFormula(2147483647))
}
