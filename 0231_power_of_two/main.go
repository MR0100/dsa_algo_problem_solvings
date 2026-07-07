package main

import (
	"fmt"
	"math"
)

// ── Approach 1: Iterative Division ───────────────────────────────────────────
//
// iterativeDivision solves Power of Two by repeatedly dividing out factors of 2.
//
// Intuition:
//
//	A power of two is 2^k for some k >= 0, i.e. its only prime factor is 2.
//	If we keep dividing a positive number by 2 while it stays even, a genuine
//	power of two collapses all the way down to 1; anything with another prime
//	factor gets stuck on an odd number > 1.
//
// Algorithm:
//
//  1. Powers of two are positive, so reject n <= 0 immediately.
//  2. While n is even (n % 2 == 0), divide n by 2.
//  3. n is a power of two iff the leftover equals 1.
//
// Time:  O(log n) — at most log2(n) halvings.
// Space: O(1) — a single mutable integer.
func iterativeDivision(n int) bool {
	if n <= 0 { // powers of two (2^k, k>=0) are all >= 1, so non-positive fails
		return false
	}
	for n%2 == 0 { // strip a factor of 2 as long as one is present
		n /= 2
	}
	return n == 1 // only a pure power of two reduces exactly to 1
}

// ── Approach 2: Bit Count ────────────────────────────────────────────────────
//
// bitCount solves Power of Two by checking that exactly one bit is set.
//
// Intuition:
//
//	In binary, 2^k is a 1 followed by k zeros (1, 10, 100, 1000, ...). So a
//	positive number is a power of two exactly when its binary representation
//	contains a single set bit. Count the set bits and compare with 1.
//
// Algorithm:
//
//  1. Reject n <= 0.
//  2. Walk the bits, counting how many are 1.
//  3. Return true iff the count is exactly 1.
//
// Time:  O(log n) — one pass over the ~log2(n) bits.
// Space: O(1) — a counter.
func bitCount(n int) bool {
	if n <= 0 { // negatives and zero are never powers of two
		return false
	}
	count := 0
	for n > 0 { // examine each bit from least significant upward
		count += n & 1 // add 1 if the current lowest bit is set
		n >>= 1        // shift to inspect the next bit
	}
	return count == 1 // a lone set bit means n == 2^k
}

// ── Approach 3: Brian Kernighan / n & (n-1) (Optimal) ────────────────────────
//
// bitTrick solves Power of Two with the classic n & (n-1) == 0 test.
//
// Intuition:
//
//	Subtracting 1 from a power of two flips its single set bit to 0 and turns
//	every lower bit to 1 (e.g. 1000 - 1 = 0111). ANDing the two therefore
//	yields 0. For any number with two or more set bits, the highest bits
//	survive the AND, giving a non-zero result. So "n > 0 and n & (n-1) == 0"
//	characterises powers of two in O(1).
//
// Algorithm:
//
//  1. Reject n <= 0.
//  2. Return true iff n & (n-1) == 0.
//
// Time:  O(1) — one subtraction and one AND.
// Space: O(1).
func bitTrick(n int) bool {
	// n>0 rules out zero/negatives; n&(n-1)==0 means n has a single set bit.
	return n > 0 && n&(n-1) == 0
}

// ── Approach 4: Divisor of a Max Power of Two ────────────────────────────────
//
// maxPowerDivisor solves Power of Two by checking divisibility of the largest
// power of two representable, avoiding any per-bit loop.
//
// Intuition:
//
//	Within a fixed integer width, every power of two divides the largest
//	power of two that fits. For 32-bit signed ints the biggest is 2^30. Any
//	power of two n in [1, 2^30] divides 2^30 evenly; a non-power never does.
//	This turns the test into a single modulo, at the cost of assuming the
//	value range (n <= 2^31 - 1 per the constraints).
//
// Algorithm:
//
//  1. Reject n <= 0.
//  2. Let maxPow = 2^30 (the greatest power of two below 2^31).
//  3. Return true iff maxPow % n == 0.
//
// Time:  O(1) — one modulo.
// Space: O(1).
func maxPowerDivisor(n int) bool {
	if n <= 0 { // guard against non-positive and avoid % by zero
		return false
	}
	const maxPow = 1 << 30 // 2^30, the largest power of two < 2^31 (int32 range)
	return maxPow%n == 0   // every power of two in range divides 2^30 exactly
}

func main() {
	// Sanity: 2^30 must not overflow the values we test against.
	_ = math.MaxInt32

	fmt.Println("=== Approach 1: Iterative Division ===")
	fmt.Println(iterativeDivision(1))  // expected true   (2^0)
	fmt.Println(iterativeDivision(16)) // expected true   (2^4)
	fmt.Println(iterativeDivision(3))  // expected false

	fmt.Println("=== Approach 2: Bit Count ===")
	fmt.Println(bitCount(1))  // expected true
	fmt.Println(bitCount(16)) // expected true
	fmt.Println(bitCount(3))  // expected false

	fmt.Println("=== Approach 3: Brian Kernighan / n & (n-1) (Optimal) ===")
	fmt.Println(bitTrick(1))  // expected true
	fmt.Println(bitTrick(16)) // expected true
	fmt.Println(bitTrick(3))  // expected false

	fmt.Println("=== Approach 4: Divisor of a Max Power of Two ===")
	fmt.Println(maxPowerDivisor(1))  // expected true
	fmt.Println(maxPowerDivisor(16)) // expected true
	fmt.Println(maxPowerDivisor(3))  // expected false
}
