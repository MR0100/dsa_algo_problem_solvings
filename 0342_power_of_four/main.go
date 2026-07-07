package main

import (
	"fmt"
	"math"
)

// ── Approach 1: Iterative Division ───────────────────────────────────────────
//
// iterativeDivision solves Power of Four by repeatedly dividing by 4.
//
// Intuition:
//
//	n is a power of four iff you can divide it by 4 down to exactly 1 with no
//	remainder along the way. Any leftover remainder means it is not 4^k.
//
// Algorithm:
//  1. Reject n <= 0 (powers of four are positive).
//  2. While n % 4 == 0, divide n by 4.
//  3. n is a power of four iff it has shrunk to 1.
//
// Time:  O(log4 n) — one division per power of four.
// Space: O(1).
func iterativeDivision(n int) bool {
	if n <= 0 { // powers of four are 1,4,16,... all positive
		return false
	}
	for n%4 == 0 { // strip factors of 4
		n /= 4
	}
	return n == 1 // only survivor of a pure power of four is 1
}

// ── Approach 2: Bit Manipulation (Optimal) ───────────────────────────────────
//
// bitTrick solves Power of Four with a constant-time bit check.
//
// Intuition:
//
//	A power of four is a power of two (single set bit) whose bit sits at an
//	EVEN position (bits 0,2,4,...): 4^0=1 (bit0), 4^1=4 (bit2), 4^2=16 (bit4).
//	Two conditions: (a) exactly one bit set → n & (n-1) == 0, and (b) that bit
//	is at an even index → it overlaps the mask 0x55555555 (…0101), which has
//	1s only at even positions.
//
// Algorithm:
//  1. n > 0.
//  2. n & (n-1) == 0  → n is a power of two (single set bit).
//  3. n & 0x55555555 != 0 → that bit is at an even position.
//
// Time:  O(1).
// Space: O(1).
func bitTrick(n int) bool {
	// 0x55555555 = 0101...0101, ones at even bit positions only.
	return n > 0 && n&(n-1) == 0 && n&0x55555555 != 0
}

// ── Approach 3: Modulo-3 Property ────────────────────────────────────────────
//
// moduloThree solves Power of Four using the identity 4^k ≡ 1 (mod 3).
//
// Intuition:
//
//	4 ≡ 1 (mod 3), so 4^k ≡ 1^k ≡ 1 (mod 3) for every k. A power of TWO that
//	is NOT a power of four is 2^odd; 2 ≡ 2 (mod 3), so 2^odd ≡ 2 (mod 3). So
//	among powers of two, the ones with remainder 1 mod 3 are exactly the
//	powers of four.
//
// Algorithm:
//  1. n > 0 and n is a power of two (n & (n-1) == 0).
//  2. n % 3 == 1.
//
// Time:  O(1).
// Space: O(1).
func moduloThree(n int) bool {
	return n > 0 && n&(n-1) == 0 && n%3 == 1
}

// ── Approach 4: Logarithm Check ──────────────────────────────────────────────
//
// logCheck solves Power of Four by asking whether log4(n) is an integer.
//
// Intuition:
//
//	n = 4^k  ⇔  log(n)/log(4) = k is a whole number. Because floating point is
//	imprecise, compare against the rounded value rather than testing equality.
//
// Algorithm:
//  1. n > 0.
//  2. Compute x = log(n)/log(4); n is a power of four iff x is (near) integer.
//
// Time:  O(1).
// Space: O(1).
func logCheck(n int) bool {
	if n <= 0 {
		return false
	}
	x := math.Log(float64(n)) / math.Log(4) // log base 4 of n
	// Round and check the rounded exponent reproduces n exactly.
	r := math.Round(x)
	return math.Abs(x-r) < 1e-10 && int(math.Pow(4, r)) == n
}

func main() {
	fmt.Println("=== Approach 1: Iterative Division ===")
	fmt.Println(iterativeDivision(16)) // expected true
	fmt.Println(iterativeDivision(5))  // expected false
	fmt.Println(iterativeDivision(1))  // expected true

	fmt.Println("=== Approach 2: Bit Manipulation (Optimal) ===")
	fmt.Println(bitTrick(16)) // expected true
	fmt.Println(bitTrick(5))  // expected false
	fmt.Println(bitTrick(1))  // expected true

	fmt.Println("=== Approach 3: Modulo-3 Property ===")
	fmt.Println(moduloThree(16)) // expected true
	fmt.Println(moduloThree(5))  // expected false
	fmt.Println(moduloThree(1))  // expected true

	fmt.Println("=== Approach 4: Logarithm Check ===")
	fmt.Println(logCheck(16)) // expected true
	fmt.Println(logCheck(5))  // expected false
	fmt.Println(logCheck(1))  // expected true
}
