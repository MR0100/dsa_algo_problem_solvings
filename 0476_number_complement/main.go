package main

import (
	"fmt"
	"math/bits"
)

// ── Approach 1: Bit by Bit ───────────────────────────────────────────────────
//
// bitByBit solves Number Complement by walking every bit position from the
// least-significant bit up to (and including) the highest set bit, flipping
// each one and reassembling the result.
//
// Intuition:
//
//	The "complement" only spans the bits up to the most-significant 1 — there
//	are no leading zeros to flip. So we scan bit positions 0,1,2,... while any
//	bit of num remains. For each position we read num's bit, invert it, and
//	place the inverted bit into the answer at the same position.
//
// Algorithm:
//  1. result = 0, position = 0.
//  2. While a copy of num is non-zero:
//     - take its lowest bit (copy & 1), flip it (1 - bit), and OR the flipped
//     bit shifted to `position` into result.
//     - shift the copy right by 1, position++.
//  3. Return result.
//
// Time:  O(number of bits) = O(log num) — one pass over ~31 bits.
// Space: O(1) — a handful of integer accumulators.
func bitByBit(num int) int {
	result := 0      // the complement we are assembling
	position := 0    // which bit position we are currently at
	remaining := num // consume a copy so num stays intact
	for remaining > 0 {
		bit := remaining & 1          // read the lowest bit of the remaining number
		flipped := 1 - bit            // invert it: 0→1, 1→0
		result |= flipped << position // drop the flipped bit at its position
		remaining >>= 1               // advance to the next-higher bit
		position++                    // and remember the new position
	}
	return result
}

// ── Approach 2: XOR with an All-Ones Mask (Optimal) ──────────────────────────
//
// xorMask solves Number Complement by XOR-ing num with a mask that is all 1s
// exactly as wide as num's binary representation.
//
// Intuition:
//
//	Flipping bit b is exactly b XOR 1. To flip *every* meaningful bit at once,
//	XOR num with a mask of all 1s that is the same bit-length as num. If num
//	uses L bits, that mask is (1 << L) - 1. Then num ^ mask flips precisely the
//	L relevant bits and leaves nothing above them (mask is 0 there, so those
//	high bits stay 0).
//
// Algorithm:
//  1. Find L = number of bits in num (position of highest set bit + 1).
//  2. mask = (1 << L) - 1  (L ones).
//  3. Return num ^ mask.
//
// Edge case: num == 0 has L = 0, mask = 0, and 0 ^ 0 = 0, which is the correct
// complement of the empty/zero representation.
//
// Time:  O(1) — bits.Len plus a couple of arithmetic ops.
// Space: O(1).
func xorMask(num int) int {
	length := bits.Len(uint(num)) // number of significant bits (0 for num==0)
	mask := (1 << length) - 1     // 'length' consecutive 1s: 5(101)→111
	return num ^ mask             // XOR flips exactly those bits
}

// ── Approach 3: Smear Bits then XOR ──────────────────────────────────────────
//
// smearMask solves Number Complement by "smearing" num's highest bit downward
// to build the all-ones mask without computing a length, then XOR-ing.
//
// Intuition:
//
//	OR-ing num with itself shifted right by 1,2,4,8,16 turns the leading 1 and
//	every bit below it into 1 — producing exactly the same all-ones mask as
//	Approach 2, but purely with bit tricks (this is the classic "fill bits
//	below the MSB" idiom). XOR num with that mask to flip the relevant bits.
//
// Algorithm:
//  1. mask = num.
//  2. mask |= mask>>1; |= mask>>2; |= mask>>4; |= mask>>8; |= mask>>16.
//  3. Return num ^ mask.
//
// Time:  O(1) — a fixed five shifts/ORs (covers 32 bits).
// Space: O(1).
func smearMask(num int) int {
	mask := num        // start from num
	mask |= mask >> 1  // fill 1 bit below each set bit
	mask |= mask >> 2  // fill 2 more
	mask |= mask >> 4  // 4 more
	mask |= mask >> 8  // 8 more
	mask |= mask >> 16 // 16 more → every bit from MSB down is now 1
	return num ^ mask  // flip exactly the smeared region
}

func main() {
	fmt.Println("=== Approach 1: Bit by Bit ===")
	fmt.Printf("num=5  got=%d  expected 2\n", bitByBit(5)) // 101 -> 010
	fmt.Printf("num=1  got=%d  expected 0\n", bitByBit(1)) // 1 -> 0
	fmt.Printf("num=0  got=%d  expected 0\n", bitByBit(0)) // edge: 0 -> 0
	fmt.Printf("num=7  got=%d  expected 0\n", bitByBit(7)) // 111 -> 000

	fmt.Println("=== Approach 2: XOR with All-Ones Mask (Optimal) ===")
	fmt.Printf("num=5  got=%d  expected 2\n", xorMask(5))
	fmt.Printf("num=1  got=%d  expected 0\n", xorMask(1))
	fmt.Printf("num=0  got=%d  expected 0\n", xorMask(0))
	fmt.Printf("num=7  got=%d  expected 0\n", xorMask(7))

	fmt.Println("=== Approach 3: Smear Bits then XOR ===")
	fmt.Printf("num=5  got=%d  expected 2\n", smearMask(5))
	fmt.Printf("num=1  got=%d  expected 0\n", smearMask(1))
	fmt.Printf("num=0  got=%d  expected 0\n", smearMask(0))
	fmt.Printf("num=7  got=%d  expected 0\n", smearMask(7))
}
