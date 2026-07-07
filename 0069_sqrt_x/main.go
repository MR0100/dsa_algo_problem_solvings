package main

import "fmt"

// ── Approach 1: Linear Scan ───────────────────────────────────────────────────
//
// linearScan solves Sqrt(x) by trying each integer from 0 upward.
//
// Intuition:
//   Find the largest k such that k² ≤ x. Increment until k² > x.
//
// Time:  O(√x) — loop runs √x iterations.
// Space: O(1)
func linearScan(x int) int {
	if x == 0 {
		return 0
	}
	k := 1
	for k*k <= x {
		k++
	}
	return k - 1 // last k where k²<=x
}

// ── Approach 2: Binary Search ─────────────────────────────────────────────────
//
// binarySearch solves Sqrt(x) using binary search on the answer space [0, x].
//
// Intuition:
//   The answer k satisfies k² ≤ x < (k+1)². Binary search: if mid*mid <= x,
//   the answer is ≥ mid (store it, move lo right). Otherwise move hi left.
//
// Algorithm:
//   lo=0, hi=x; ans=0
//   while lo<=hi:
//     mid = (lo+hi)/2
//     if mid*mid <= x: ans=mid; lo=mid+1
//     else: hi=mid-1
//   return ans
//
// Time:  O(log x)
// Space: O(1)
func binarySearch(x int) int {
	if x == 0 {
		return 0
	}
	lo, hi, ans := 0, x, 0
	for lo <= hi {
		mid := lo + (hi-lo)/2
		if mid*mid <= x {
			ans = mid    // mid is a candidate
			lo = mid + 1 // try to find a larger one
		} else {
			hi = mid - 1
		}
	}
	return ans
}

// ── Approach 3: Newton's Method ───────────────────────────────────────────────
//
// newtonMethod solves Sqrt(x) using Newton-Raphson iteration.
//
// Intuition:
//   We want to find root of f(r) = r² - x = 0.
//   Newton update: r' = r - f(r)/f'(r) = r - (r²-x)/(2r) = (r + x/r) / 2
//   Starting from r = x, converge quadratically to √x.
//
// Time:  O(log x) — quadratic convergence means few iterations.
// Space: O(1)
func newtonMethod(x int) int {
	if x == 0 {
		return 0
	}
	r := x
	for r*r > x {
		r = (r + x/r) / 2 // Newton update (integer division)
	}
	return r
}

func main() {
	cases := []struct {
		x        int
		expected int
	}{
		{4, 2},
		{8, 2},
		{0, 0},
		{1, 1},
		{9, 3},
		{100, 10},
		{2147395600, 46340},
	}

	fmt.Println("=== Approach 1: Linear Scan ===")
	for _, c := range cases {
		fmt.Printf("x=%-12d  got=%d  expected=%d\n", c.x, linearScan(c.x), c.expected)
	}

	fmt.Println("=== Approach 2: Binary Search ===")
	for _, c := range cases {
		fmt.Printf("x=%-12d  got=%d  expected=%d\n", c.x, binarySearch(c.x), c.expected)
	}

	fmt.Println("=== Approach 3: Newton's Method ===")
	for _, c := range cases {
		fmt.Printf("x=%-12d  got=%d  expected=%d\n", c.x, newtonMethod(c.x), c.expected)
	}
}
