package main

import (
	"fmt"
	"math"
)

// ── Approach 1: Brute Force Simulation ───────────────────────────────────────
//
// bruteForce solves Bulb Switcher by literally simulating every round: on round
// r it toggles every r-th bulb, then counts how many remain on.
//
// Intuition:
//
//	Follow the problem statement exactly. Keep a boolean per bulb (off). Round
//	r toggles bulbs at positions r, 2r, 3r, ... Count the trues at the end.
//
// Algorithm:
//  1. bulbs = n booleans, all false (off).
//  2. For r = 1..n: for pos = r, 2r, ... ≤ n: flip bulbs[pos-1].
//  3. Count true entries.
//
// Time:  O(n log n) — the toggle loop runs n/1 + n/2 + ... + n/n ≈ n·Hₙ times.
// Space: O(n) — the bulb array.
func bruteForce(n int) int {
	bulbs := make([]bool, n) // all off initially
	for r := 1; r <= n; r++ {
		for pos := r; pos <= n; pos += r {
			bulbs[pos-1] = !bulbs[pos-1] // toggle every r-th bulb
		}
	}
	count := 0
	for _, on := range bulbs {
		if on {
			count++ // tally bulbs left on
		}
	}
	return count
}

// ── Approach 2: Count Divisors (Perfect Squares) ─────────────────────────────
//
// countDivisors solves Bulb Switcher using the insight that bulb i ends ON iff
// i has an odd number of divisors — checking divisor parity per bulb.
//
// Intuition:
//
//	Bulb i is toggled once for each divisor of i (round r toggles i exactly when
//	r divides i). It ends ON iff it was toggled an odd number of times, i.e. iff
//	i has an odd number of divisors. Divisors pair up (d, i/d), so the count is
//	odd only when i is a perfect square (the pair d == i/d is unpaired).
//
// Algorithm:
//  1. For i = 1..n: count divisors; if odd, increment the answer.
//     (Equivalently: check whether i is a perfect square.)
//
// Time:  O(n) using the perfect-square test per bulb.
// Space: O(1).
func countDivisors(n int) int {
	count := 0
	for i := 1; i <= n; i++ {
		root := int(math.Sqrt(float64(i)))
		if root*root == i { // i is a perfect square → odd divisor count → ON
			count++
		}
	}
	return count
}

// ── Approach 3: Integer Square Root (Optimal) ────────────────────────────────
//
// integerSqrt solves Bulb Switcher in O(1): the answer is simply the number of
// perfect squares in [1, n], which is floor(sqrt(n)).
//
// Intuition:
//
//	From Approach 2, exactly the perfect-square-indexed bulbs stay on. The count
//	of perfect squares ≤ n is floor(sqrt(n)) (namely 1², 2², ..., ⌊√n⌋²). So the
//	whole problem collapses to one square root.
//
// Algorithm:
//  1. Return floor(sqrt(n)), guarding float rounding by adjusting root so that
//     root² ≤ n < (root+1)².
//
// Time:  O(1).
// Space: O(1).
func integerSqrt(n int) int {
	root := int(math.Sqrt(float64(n))) // may be off by one due to float error
	// Guard against floating-point over/undershoot at large n.
	for root*root > n {
		root-- // shrink if we overshot
	}
	for (root+1)*(root+1) <= n {
		root++ // grow if we undershot
	}
	return root
}

func main() {
	fmt.Println("=== Approach 1: Brute Force Simulation ===")
	fmt.Printf("n=3       -> %d  expected 1\n", bruteForce(3))
	fmt.Printf("n=0       -> %d  expected 0\n", bruteForce(0))
	fmt.Printf("n=1       -> %d  expected 1\n", bruteForce(1))
	fmt.Printf("n=10      -> %d  expected 3\n", bruteForce(10)) // 1,4,9

	fmt.Println("=== Approach 2: Count Divisors (Perfect Squares) ===")
	fmt.Printf("n=3       -> %d  expected 1\n", countDivisors(3))
	fmt.Printf("n=0       -> %d  expected 0\n", countDivisors(0))
	fmt.Printf("n=1       -> %d  expected 1\n", countDivisors(1))
	fmt.Printf("n=10      -> %d  expected 3\n", countDivisors(10))

	fmt.Println("=== Approach 3: Integer Square Root (Optimal) ===")
	fmt.Printf("n=3       -> %d  expected 1\n", integerSqrt(3))
	fmt.Printf("n=0       -> %d  expected 0\n", integerSqrt(0))
	fmt.Printf("n=1       -> %d  expected 1\n", integerSqrt(1))
	fmt.Printf("n=10      -> %d  expected 3\n", integerSqrt(10))
	fmt.Printf("n=99999999-> %d  expected 9999\n", integerSqrt(99999999)) // large-n float guard
}
