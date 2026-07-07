package main

import "fmt"

// ── Approach 1: Brute Force (Iterative Multiplication) ───────────────────────
//
// bruteForce solves Pow(x, n) by multiplying x by itself |n| times.
//
// Intuition: x^n = x * x * ... * x (n times). Handle negatives by computing
// x^|n| and taking the reciprocal.
//
// WARNING: O(|n|) — TLE for large n such as n = 2^31-1.
//
// Time:  O(|n|)
// Space: O(1)
func bruteForce(x float64, n int) float64 {
	if n < 0 {
		x = 1 / x
		n = -n
	}
	result := 1.0
	for i := 0; i < n; i++ {
		result *= x
	}
	return result
}

// ── Approach 2: Fast Power / Exponentiation by Squaring (Optimal) ─────────────
//
// fastPow solves Pow(x, n) using the "exponentiation by squaring" technique.
//
// Intuition:
//   x^n = (x²)^(n/2)        if n is even
//   x^n = x * (x²)^((n-1)/2) if n is odd
//
// At each step, we halve the exponent (O(log n) steps total) and square x.
// This works for negative n by computing x^(-n) = (1/x)^n.
//
// Algorithm (iterative):
//  if n < 0: x = 1/x, n = -n
//  result = 1.0
//  while n > 0:
//    if n is odd: result *= x
//    x = x * x
//    n >>= 1   (n /= 2)
//  return result
//
// Time:  O(log |n|)
// Space: O(1)
func fastPow(x float64, n int) float64 {
	if n < 0 {
		x = 1 / x
		n = -n
	}
	result := 1.0
	for n > 0 {
		if n%2 == 1 { // n is odd: absorb one factor of x into result
			result *= x
		}
		x *= x // square x: move to x^(2k) for next iteration
		n >>= 1 // n = n/2
	}
	return result
}

// ── Approach 3: Recursive Fast Power ─────────────────────────────────────────
//
// fastPowRecursive solves Pow(x, n) with recursion.
//
// Time:  O(log |n|)
// Space: O(log |n|) — recursion stack
func fastPowRecursive(x float64, n int) float64 {
	if n == 0 {
		return 1
	}
	if n < 0 {
		return fastPowRecursive(1/x, -n)
	}
	if n%2 == 0 {
		half := fastPowRecursive(x, n/2)
		return half * half // avoid double recursion
	}
	return x * fastPowRecursive(x, n-1)
}

func main() {
	cases := []struct {
		x    float64
		n    int
		want float64
	}{
		{2.0, 10, 1024.0},
		{2.1, 3, 9.261},
		{2.0, -2, 0.25},
		{1.0, 2147483647, 1.0},
		{2.0, 0, 1.0},
		{0.0, 0, 1.0},
	}

	fmt.Println("=== Approach 1: Brute Force ===")
	for _, c := range cases[:3] { // skip large n for brute force
		fmt.Printf("x=%.1f n=%d  got=%.5f  expected=%.5f\n", c.x, c.n, bruteForce(c.x, c.n), c.want)
	}

	fmt.Println("\n=== Approach 2: Fast Power Iterative (Optimal) ===")
	for _, c := range cases {
		fmt.Printf("x=%.1f n=%d  got=%.5f  expected=%.5f\n", c.x, c.n, fastPow(c.x, c.n), c.want)
	}

	fmt.Println("\n=== Approach 3: Fast Power Recursive ===")
	for _, c := range cases {
		fmt.Printf("x=%.1f n=%d  got=%.5f  expected=%.5f\n", c.x, c.n, fastPowRecursive(c.x, c.n), c.want)
	}
}
