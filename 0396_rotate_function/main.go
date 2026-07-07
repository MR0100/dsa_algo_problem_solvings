package main

import "fmt"

// ── Approach 1: Brute Force ──────────────────────────────────────────────────
//
// bruteForce solves Rotate Function by literally computing F(k) for every
// rotation k and keeping the maximum.
//
// Intuition:
//
//	The definition F(k) = Σ i·arrk[i] is directly computable. Rotating the
//	array clockwise by k moves each element; instead of physically rotating,
//	we index into the original array with (i - k + n) % n, which is where the
//	element sitting at position i of the rotated array came from.
//
// Algorithm:
//  1. n = len(nums).
//  2. For each rotation k in 0..n-1:
//     a. sum = 0.
//     b. For each position i in 0..n-1: add i * nums[(i-k+n)%n].
//     c. Track the maximum sum seen.
//  3. Return the maximum.
//
// Time:  O(n²) — n rotations, each an O(n) summation.
// Space: O(1) — only running scalars.
func bruteForce(nums []int) int {
	n := len(nums)
	best := 0                // will hold the maximum F(k)
	for k := 0; k < n; k++ { // try every rotation amount
		sum := 0
		for i := 0; i < n; i++ {
			// In arrk, the element at index i came from original index
			// (i-k) modulo n; +n keeps the result non-negative before %.
			sum += i * nums[(i-k+n)%n]
		}
		if k == 0 || sum > best { // seed best on first k, then maximise
			best = sum
		}
	}
	return best
}

// ── Approach 2: Rolling Recurrence (Optimal) ─────────────────────────────────
//
// rollingRecurrence solves Rotate Function with the O(n) increment identity
// F(k) = F(k-1) + sum - n*nums[n-k].
//
// Intuition:
//
//	Compare F(k) with F(k-1). Rotating one more step increases every element's
//	coefficient by 1 (adding one whole `sum`), except the element that wraps
//	from coefficient (n-1) back to coefficient 0, which loses n·value. That
//	element is nums[n-k]. So each F(k) is derived from the previous one in O(1).
//
// Algorithm:
//  1. Compute sum = Σ nums and f = F(0) = Σ i·nums[i].
//  2. best = f.
//  3. For k in 1..n-1: f = f + sum - n*nums[n-k]; best = max(best, f).
//  4. Return best.
//
// Time:  O(n) — one pass to seed, one pass for the recurrence.
// Space: O(1) — a handful of accumulators.
func rollingRecurrence(nums []int) int {
	n := len(nums)
	sum := 0 // Σ nums[i]
	f := 0   // running F(k); starts as F(0)
	for i, v := range nums {
		sum += v   // total of all elements
		f += i * v // F(0) = Σ i·nums[i]
	}
	best := f // F(0) is our first candidate
	for k := 1; k < n; k++ {
		// F(k) = F(k-1) + sum - n*nums[n-k]: every coefficient +1 (adds sum),
		// but the wrap element nums[n-k] drops from coeff (n-1) to 0.
		f = f + sum - n*nums[n-k]
		if f > best {
			best = f
		}
	}
	return best
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Printf("nums=[4,3,2,6]  got=%d  expected 26\n", bruteForce([]int{4, 3, 2, 6}))
	fmt.Printf("nums=[100]      got=%d  expected 0\n", bruteForce([]int{100}))

	fmt.Println("=== Approach 2: Rolling Recurrence (Optimal) ===")
	fmt.Printf("nums=[4,3,2,6]  got=%d  expected 26\n", rollingRecurrence([]int{4, 3, 2, 6}))
	fmt.Printf("nums=[100]      got=%d  expected 0\n", rollingRecurrence([]int{100}))
}
