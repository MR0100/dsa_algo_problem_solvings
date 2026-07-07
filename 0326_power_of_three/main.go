package main

import (
	"fmt"
	"math"
)

// ── Approach 1: Loop Division ────────────────────────────────────────────────
//
// loopDivision solves Power of Three by repeatedly dividing n by 3.
//
// Intuition:
//
//	A power of three, 3^x, is built up by multiplying 3 into 1 exactly x times.
//	Run that in reverse: keep dividing by 3 while the number is still evenly
//	divisible. If n is a genuine power of three we strip it all the way down to
//	1; if any leftover factor other than 3 exists, a division will leave a
//	non-zero remainder and we can reject early. Non-positive numbers can never
//	be a power of three (3^x is always ≥ 1), so they are rejected up front.
//
// Algorithm:
//  1. If n <= 0, return false (3^x ≥ 1 for all integer x ≥ 0).
//  2. While n is divisible by 3 (n % 3 == 0), divide n by 3.
//  3. After the loop, n is a power of three iff it has been reduced to 1.
//
// Time:  O(log₃ n) — each division shrinks n by a factor of 3.
// Space: O(1) — only mutates n in place.
func loopDivision(n int) bool {
	if n <= 0 {
		return false // powers of three are always positive (3^0 = 1)
	}
	for n%3 == 0 { // peel off one factor of 3 per iteration while it divides cleanly
		n /= 3
	}
	return n == 1 // fully reduced to 1 ⇒ it was purely 3^x
}

// ── Approach 2: Logarithm ────────────────────────────────────────────────────
//
// logarithm solves Power of Three using logarithms and a tolerance check.
//
// Intuition:
//
//	If n == 3^x then x = log₃(n) = ln(n) / ln(3) must be a whole number. We
//	compute that ratio and test whether it is (near) an integer. Because
//	floating-point logs carry rounding error, we round to the nearest integer
//	and compare against a small epsilon rather than testing for exact equality.
//
// Algorithm:
//  1. If n <= 0, return false (log of a non-positive number is undefined here).
//  2. Compute x = log₁₀(n) / log₁₀(3)  (any common base cancels out).
//  3. Return true iff x is within a tiny epsilon of its nearest integer.
//
// Time:  O(1) — a couple of constant-time math library calls.
// Space: O(1) — a few scalars.
func logarithm(n int) bool {
	if n <= 0 {
		return false // logarithm is undefined / meaningless for n ≤ 0
	}
	x := math.Log10(float64(n)) / math.Log10(3) // change-of-base: log₃(n)
	// A true power of three gives an integer x; FP error means we test nearness
	// to the nearest integer instead of exact equality.
	return math.Abs(x-math.Round(x)) < 1e-10
}

// ── Approach 3: Integer Limit (No Loops — Optimal) ───────────────────────────
//
// integerLimit solves Power of Three without any loop or recursion by using the
// largest power of three that fits in a signed 32-bit integer.
//
// Intuition:
//
//	3 is prime, so the divisors of 3^19 are exactly 1, 3, 9, 27, …, 3^19 — the
//	powers of three and nothing else. Within int32 the largest power of three is
//	3^19 = 1162261467 (3^20 overflows int32). Therefore n is a power of three
//	iff n is positive AND n divides 1162261467 evenly. No looping required.
//
// Algorithm:
//  1. Let MAX = 1162261467 = 3^19, the biggest power of three ≤ 2³¹ − 1.
//  2. Return true iff n > 0 and MAX % n == 0.
//
// Time:  O(1) — one comparison and one modulo.
// Space: O(1) — no extra storage.
func integerLimit(n int) bool {
	const maxPow3 = 1162261467     // 3^19, the largest power of three within int32
	return n > 0 && maxPow3%n == 0 // divisors of 3^19 (a prime power) are exactly the powers of three
}

func main() {
	fmt.Println("=== Approach 1: Loop Division ===")
	fmt.Printf("n=27  got=%v  expected true\n", loopDivision(27))  // expected true
	fmt.Printf("n=0   got=%v  expected false\n", loopDivision(0))  // expected false
	fmt.Printf("n=-1  got=%v  expected false\n", loopDivision(-1)) // expected false
	fmt.Printf("n=9   got=%v  expected true\n", loopDivision(9))   // expected true
	fmt.Printf("n=45  got=%v  expected false\n", loopDivision(45)) // expected false
	fmt.Printf("n=1   got=%v  expected true\n", loopDivision(1))   // expected true

	fmt.Println("=== Approach 2: Logarithm ===")
	fmt.Printf("n=27  got=%v  expected true\n", logarithm(27))  // expected true
	fmt.Printf("n=0   got=%v  expected false\n", logarithm(0))  // expected false
	fmt.Printf("n=-1  got=%v  expected false\n", logarithm(-1)) // expected false
	fmt.Printf("n=9   got=%v  expected true\n", logarithm(9))   // expected true
	fmt.Printf("n=45  got=%v  expected false\n", logarithm(45)) // expected false
	fmt.Printf("n=1   got=%v  expected true\n", logarithm(1))   // expected true

	fmt.Println("=== Approach 3: Integer Limit (No Loops — Optimal) ===")
	fmt.Printf("n=27  got=%v  expected true\n", integerLimit(27))  // expected true
	fmt.Printf("n=0   got=%v  expected false\n", integerLimit(0))  // expected false
	fmt.Printf("n=-1  got=%v  expected false\n", integerLimit(-1)) // expected false
	fmt.Printf("n=9   got=%v  expected true\n", integerLimit(9))   // expected true
	fmt.Printf("n=45  got=%v  expected false\n", integerLimit(45)) // expected false
	fmt.Printf("n=1   got=%v  expected true\n", integerLimit(1))   // expected true
}
