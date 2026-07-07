package main

import (
	"fmt"
	"math/bits"
)

// ── Approach 1: Brute Force (Compare Bit by Bit) ─────────────────────────────
//
// bruteForce solves Hamming Distance by walking all 32 bit positions of x and
// y and counting the positions where the two bits disagree.
//
// Intuition:
//
//	The Hamming distance is, by definition, "how many bit positions differ".
//	So take each position i, extract bit i of x and bit i of y, and if they
//	are not equal, tally one. No cleverness — just the literal definition.
//
// Algorithm:
//  1. distance = 0.
//  2. For i = 0 .. 31: shift x and y right by i, mask the low bit of each.
//  3. If the two extracted bits differ, distance++.
//  4. Return distance.
//
// Time:  O(32) = O(1) — a fixed 32 iterations regardless of input value.
// Space: O(1) — one counter.
func bruteForce(x int, y int) int {
	distance := 0 // running count of differing bit positions
	// 32 positions cover every value in [0, 2^31 - 1].
	for i := 0; i < 32; i++ {
		bitX := (x >> i) & 1 // isolate bit i of x (0 or 1)
		bitY := (y >> i) & 1 // isolate bit i of y (0 or 1)
		if bitX != bitY {    // positions disagree → contributes to the distance
			distance++
		}
	}
	return distance
}

// ── Approach 2: XOR then Count Set Bits (Loop) ───────────────────────────────
//
// xorCountLoop solves Hamming Distance by XORing x and y — which lights up
// exactly the differing positions — then counting the 1 bits with a shift loop.
//
// Intuition:
//
//	XOR has the property a^b has a 1 exactly where a and b differ. So the
//	Hamming distance is simply popcount(x ^ y): the number of set bits in the
//	XOR. Counting set bits by repeatedly testing the lowest bit and shifting
//	right is the straightforward way to do it.
//
// Algorithm:
//  1. xor = x ^ y (differing positions become 1).
//  2. While xor != 0: add (xor & 1) to distance, then xor >>= 1.
//  3. Return distance.
//
// Time:  O(log(max(x,y))) — loops until the highest set bit is consumed, ≤ 31.
// Space: O(1) — one accumulator.
func xorCountLoop(x int, y int) int {
	xor := x ^ y  // 1s mark exactly the positions where x and y differ
	distance := 0 // number of set bits so far
	// Peel off the lowest bit each iteration until nothing is left.
	for xor != 0 {
		distance += xor & 1 // add 1 if the lowest bit is set
		xor >>= 1           // drop the lowest bit and continue
	}
	return distance
}

// ── Approach 3: XOR then Brian Kernighan's Trick (Optimal) ───────────────────
//
// xorKernighan solves Hamming Distance by XORing x and y and counting the set
// bits of the result with Kernighan's `n & (n-1)` lowest-set-bit clearing.
//
// Intuition:
//
//	Same XOR insight, but count the 1 bits faster: `n & (n-1)` clears the
//	single lowest set bit of n in one operation. Each loop iteration removes
//	exactly one 1, so the loop runs popcount(xor) times — not once per bit
//	position, but once per differing bit. This is the fewest-iterations way
//	to count bits by hand.
//
// Algorithm:
//  1. xor = x ^ y.
//  2. While xor != 0: xor &= xor - 1 (clear lowest set bit), distance++.
//  3. Return distance.
//
// Time:  O(popcount(x^y)) ≤ O(32) — one step per differing bit only.
// Space: O(1).
func xorKernighan(x int, y int) int {
	xor := x ^ y  // differing positions become 1
	distance := 0 // count of set bits removed
	// Each pass deletes the lowest remaining 1 bit.
	for xor != 0 {
		xor &= xor - 1 // Kernighan: clears exactly the lowest set bit
		distance++     // we just accounted for one differing position
	}
	return distance
}

// ── Approach 4: Built-in Population Count ────────────────────────────────────
//
// builtinPopcount solves Hamming Distance by delegating the set-bit count of
// the XOR to the standard library's hardware-backed popcount.
//
// Intuition:
//
//	The whole answer is popcount(x ^ y). Go's math/bits.OnesCount exposes a
//	single POPCNT-class instruction on modern CPUs, so this is the shortest,
//	fastest real-world implementation once you accept the XOR insight.
//
// Algorithm:
//  1. Return bits.OnesCount(uint(x ^ y)).
//
// Time:  O(1) — one CPU instruction on hardware with POPCNT.
// Space: O(1).
func builtinPopcount(x int, y int) int {
	// x^y marks differing bits; OnesCount reports how many there are.
	return bits.OnesCount(uint(x ^ y))
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Compare Bit by Bit) ===")
	fmt.Printf("x=1, y=4  got=%d  expected 2\n", bruteForce(1, 4)) // 0001 vs 0100 → 2
	fmt.Printf("x=3, y=1  got=%d  expected 1\n", bruteForce(3, 1)) // 0011 vs 0001 → 1

	fmt.Println("=== Approach 2: XOR then Count Set Bits (Loop) ===")
	fmt.Printf("x=1, y=4  got=%d  expected 2\n", xorCountLoop(1, 4))
	fmt.Printf("x=3, y=1  got=%d  expected 1\n", xorCountLoop(3, 1))

	fmt.Println("=== Approach 3: XOR then Brian Kernighan's Trick (Optimal) ===")
	fmt.Printf("x=1, y=4  got=%d  expected 2\n", xorKernighan(1, 4))
	fmt.Printf("x=3, y=1  got=%d  expected 1\n", xorKernighan(3, 1))

	fmt.Println("=== Approach 4: Built-in Population Count ===")
	fmt.Printf("x=1, y=4  got=%d  expected 2\n", builtinPopcount(1, 4))
	fmt.Printf("x=3, y=1  got=%d  expected 1\n", builtinPopcount(3, 1))
}
