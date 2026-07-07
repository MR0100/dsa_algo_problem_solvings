package main

import "fmt"

// ── Approach 1: Linear Scan (Brute Force) ────────────────────────────────────
//
// linearScan solves Valid Perfect Square by trying every candidate root i and
// checking whether i*i equals num.
//
// Intuition:
//
//	A perfect square is num == i*i for some non-negative integer i. Try i = 1,
//	2, 3, ... The moment i*i meets or exceeds num we can stop: squares grow
//	monotonically, so nothing larger can equal num either.
//
// Algorithm:
//  1. For i = 1, 2, 3, ...: compute sq = i*i.
//  2. If sq == num, return true.
//  3. If sq > num, return false (we overshot).
//
// Time:  O(√num) — the loop runs until i reaches √num.
// Space: O(1).
func linearScan(num int) bool {
	for i := 1; i*i <= num; i++ { // stop once the square meets/exceeds num
		if i*i == num { // exact hit ⇒ perfect square
			return true
		}
	}
	return false // never landed exactly on num
}

// ── Approach 2: Binary Search (Optimal) ──────────────────────────────────────
//
// binarySearch solves Valid Perfect Square by binary-searching the integer root
// in the range [1, num].
//
// Intuition:
//
//	The candidate roots 1..num are sorted, and i*i is monotonically increasing,
//	so we can binary-search for a mid whose square equals num. Each step halves
//	the search space instead of stepping one at a time.
//
// Algorithm:
//  1. lo = 1, hi = num.
//  2. While lo <= hi: mid = lo + (hi-lo)/2, sq = mid*mid.
//     - sq == num → return true.
//     - sq <  num → search the upper half (lo = mid+1).
//     - sq >  num → search the lower half (hi = mid-1).
//  3. Return false.
//
// Time:  O(log num) — halving the interval each iteration.
// Space: O(1).
func binarySearch(num int) bool {
	lo, hi := 1, num
	for lo <= hi {
		mid := lo + (hi-lo)/2 // overflow-safe midpoint
		sq := mid * mid       // candidate square
		switch {
		case sq == num:
			return true // found the exact root
		case sq < num:
			lo = mid + 1 // too small; go right
		default:
			hi = mid - 1 // too big; go left
		}
	}
	return false
}

// ── Approach 3: Newton's Method (Optimal) ────────────────────────────────────
//
// newtonsMethod solves Valid Perfect Square by iterating Newton's root-finding
// update until it converges, then checking the result squared.
//
// Intuition:
//
//	To solve x² = num, Newton's method for f(x) = x² − num gives the update
//	x ← (x + num/x) / 2. Starting from x = num, the sequence converges
//	quadratically to ⌊√num⌋. Because integer division floors, iterating until x
//	stops decreasing lands on the floor of the true root; then x*x == num tells
//	us whether num is a perfect square.
//
// Algorithm:
//  1. x = num.
//  2. While x*x > num: x = (x + num/x) / 2.
//  3. Return x*x == num.
//
// Time:  O(log num) — quadratic convergence needs very few iterations.
// Space: O(1).
func newtonsMethod(num int) bool {
	x := num        // initial guess (an over-estimate)
	for x*x > num { // shrink until x is the floor of the root
		x = (x + num/x) / 2 // Newton step toward √num
	}
	return x*x == num // exact only if num is a perfect square
}

func main() {
	// Example 1: num = 16 → true   (4*4)
	// Example 2: num = 14 → false

	fmt.Println("=== Approach 1: Linear Scan (Brute Force) ===")
	fmt.Println(linearScan(16)) // expected true
	fmt.Println(linearScan(14)) // expected false
	fmt.Println(linearScan(1))  // expected true

	fmt.Println("=== Approach 2: Binary Search (Optimal) ===")
	fmt.Println(binarySearch(16)) // expected true
	fmt.Println(binarySearch(14)) // expected false
	fmt.Println(binarySearch(1))  // expected true

	fmt.Println("=== Approach 3: Newton's Method (Optimal) ===")
	fmt.Println(newtonsMethod(16)) // expected true
	fmt.Println(newtonsMethod(14)) // expected false
	fmt.Println(newtonsMethod(1))  // expected true
}
