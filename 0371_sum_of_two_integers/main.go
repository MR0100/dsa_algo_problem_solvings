package main

import "fmt"

// ── Approach 1: Bitwise Iterative (Carry Loop) ───────────────────────────────
//
// bitwiseIterative adds two integers without using + or - by simulating the way
// a hardware adder works: XOR gives the partial sum (add without carry) and
// AND<<1 gives the carry. Repeat until there is no carry left.
//
// Intuition:
//
//	When you add two bits, XOR is the sum bit (1^1=0, 1^0=1, 0^0=0) and AND is
//	the carry (1&1=1). Shifting the carry left by one lines it up with the next
//	column, exactly like grade-school addition. Feeding (sum, carry) back into
//	the same rule until the carry becomes 0 yields the full sum.
//
// Algorithm:
//  1. While b != 0:
//     a. carry = (a & b) << 1  — bits where both are 1, moved to the next column.
//     b. a = a ^ b             — sum of a and b ignoring carry.
//     c. b = carry             — carry becomes the new addend.
//  2. Return a.
//
// Go note: Go forbids shifting negative signed values in some contexts and the
// carry chain relies on two's-complement wraparound, so we compute in uint and
// convert back. int is at least 32 bits; using uint (64-bit here) is safe for
// the −1000..1000 range and for full 32-bit inputs alike.
//
// Time:  O(1) — at most 32 (here 64) iterations, one per bit position.
// Space: O(1) — a couple of scalars.
func bitwiseIterative(a int, b int) int {
	ua, ub := uint(a), uint(b) // work in unsigned to get clean two's-complement wrap
	for ub != 0 {              // loop until nothing left to carry
		carry := (ua & ub) << 1 // columns where both bits are 1 carry into the next
		ua = ua ^ ub            // add the two numbers ignoring carry
		ub = carry              // the carry becomes the next thing to add
	}
	return int(ua) // reinterpret the bit pattern as a signed int
}

// ── Approach 2: Bitwise Recursive ────────────────────────────────────────────
//
// bitwiseRecursive is the same XOR/carry idea expressed recursively: the sum of
// a and b equals the sum of (a^b) and the carry (a&b)<<1, with base case b==0.
//
// Intuition:
//
//	getSum(a, b) = getSum(a^b, (a&b)<<1). Each recursive step pushes the
//	remaining carry one column left; because the carry can only move upward and
//	eventually falls off the top of the word, the recursion terminates when the
//	carry is 0.
//
// Algorithm:
//  1. If b == 0, return a (no carry left).
//  2. Otherwise return bitwiseRecursive(a^b, (a&b)<<1).
//
// Time:  O(1) — bounded by the word size (≤ 64 calls).
// Space: O(1) — recursion depth ≤ word size, effectively constant.
func bitwiseRecursive(a int, b int) int {
	if b == 0 { // no carry remains → a already holds the full sum
		return a
	}
	ua, ub := uint(a), uint(b)                           // unsigned for clean wraparound
	return bitwiseRecursive(int(ua^ub), int((ua&ub)<<1)) // sum-without-carry, carry
}

func main() {
	fmt.Println("=== Approach 1: Bitwise Iterative ===")
	fmt.Println(bitwiseIterative(1, 2))        // expected 3
	fmt.Println(bitwiseIterative(2, 3))        // expected 5
	fmt.Println(bitwiseIterative(-2, 3))       // expected 1
	fmt.Println(bitwiseIterative(-1, -1))      // expected -2
	fmt.Println(bitwiseIterative(1000, -1000)) // expected 0

	fmt.Println("=== Approach 2: Bitwise Recursive ===")
	fmt.Println(bitwiseRecursive(1, 2))        // expected 3
	fmt.Println(bitwiseRecursive(2, 3))        // expected 5
	fmt.Println(bitwiseRecursive(-2, 3))       // expected 1
	fmt.Println(bitwiseRecursive(-1, -1))      // expected -2
	fmt.Println(bitwiseRecursive(1000, -1000)) // expected 0
}
