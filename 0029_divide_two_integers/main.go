package main

import (
	"fmt"
	"math"
)

// ── Approach 1: Brute Force (Repeated Subtraction) ───────────────────────────
//
// bruteForce solves Divide Two Integers by repeatedly subtracting divisor
// from dividend and counting the subtractions.
//
// Intuition: Division is just "how many times can divisor fit into dividend?"
// Subtract divisor repeatedly until dividend < divisor.
//
// WARNING: O(|dividend/divisor|) — TLE for cases like INT_MIN / 1.
//
// Algorithm:
//  1. Handle sign and work with absolute values (use int64 to avoid overflow).
//  2. Subtract |divisor| from |dividend| until |dividend| < |divisor|; count.
//  3. Apply sign; clamp to [INT32_MIN, INT32_MAX].
//
// Time:  O(dividend/divisor) — up to 2^31 iterations in the worst case
// Space: O(1)
func bruteForce(dividend, divisor int) int {
	// special case: only overflow possible
	if dividend == math.MinInt32 && divisor == -1 {
		return math.MaxInt32
	}
	// determine sign of result
	negative := (dividend < 0) != (divisor < 0)

	// work in int64 to avoid overflow when negating MinInt32
	a, b := int64(dividend), int64(divisor)
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}

	count := int64(0)
	for a >= b {
		a -= b
		count++
	}

	if negative {
		return int(-count)
	}
	return int(count)
}

// ── Approach 2: Bit Shifting / Exponential Speedup (Optimal) ─────────────────
//
// bitShift solves Divide Two Integers without using multiplication, division,
// or mod operators, using bit shifting to achieve O(log² n) time.
//
// Intuition: Instead of subtracting divisor one at a time, try to subtract
// the largest multiple of divisor that fits: 2^k * divisor.
// Double that multiple (<<1) each time, doubling the count subtracted.
// When the multiple exceeds remaining dividend, restart from divisor×1.
//
// Algorithm:
//  1. Handle overflow edge case.
//  2. Determine sign; work with absolute int64 values.
//  3. While a >= b:
//     find largest shift k such that (b << k) <= a.
//     subtract (b << k) from a; add (1 << k) to quotient.
//  4. Apply sign; return quotient (clamped by overflow check).
//
// Time:  O(log² n) — outer loop runs O(log n) times; inner doubling runs O(log n)
// Space: O(1)
func bitShift(dividend, divisor int) int {
	// only overflow: MinInt32 / -1 would give MaxInt32+1
	if dividend == math.MinInt32 && divisor == -1 {
		return math.MaxInt32
	}

	negative := (dividend < 0) != (divisor < 0)

	a := int64(dividend)
	b := int64(divisor)
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}

	quotient := int64(0)
	for a >= b {
		temp := b
		multiple := int64(1)
		// double temp until it would exceed a
		for a >= (temp << 1) {
			temp <<= 1    // temp = b * 2^k
			multiple <<= 1 // multiple = 2^k
		}
		a -= temp
		quotient += multiple
	}

	if negative {
		return int(-quotient)
	}
	return int(quotient)
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Repeated Subtraction) ===")
	fmt.Printf("10 / 3   = %d  expected 3\n", bruteForce(10, 3))
	fmt.Printf("7  / -2  = %d  expected -3\n", bruteForce(7, -2))
	fmt.Printf("0  / 1   = %d  expected 0\n", bruteForce(0, 1))
	fmt.Printf("-1 / 1   = %d  expected -1\n", bruteForce(-1, 1))
	fmt.Printf("MIN/(-1) = %d  expected 2147483647\n", bruteForce(math.MinInt32, -1))

	fmt.Println("\n=== Approach 2: Bit Shift (Optimal) ===")
	fmt.Printf("10 / 3   = %d  expected 3\n", bitShift(10, 3))
	fmt.Printf("7  / -2  = %d  expected -3\n", bitShift(7, -2))
	fmt.Printf("0  / 1   = %d  expected 0\n", bitShift(0, 1))
	fmt.Printf("-1 / 1   = %d  expected -1\n", bitShift(-1, 1))
	fmt.Printf("MIN/(-1) = %d  expected 2147483647\n", bitShift(math.MinInt32, -1))
	fmt.Printf("MIN/ 1   = %d  expected -2147483648\n", bitShift(math.MinInt32, 1))
}
