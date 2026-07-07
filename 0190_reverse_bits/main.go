package main

import (
	"fmt"
	"strconv"
)

// ── Approach 1: Brute Force (Binary String Round-Trip) ───────────────────────
//
// bruteForce solves Reverse Bits by formatting the number as a 32-character
// binary string, reversing the characters, and parsing the result back.
//
// Intuition:
//
//	Treat the task literally: the problem statement shows the input as a
//	32-bit binary string, so materialise exactly that string, mirror it,
//	and convert back. Zero bit-twiddling insight required — a correctness
//	baseline that makes the bit-level approaches easy to verify against.
//
// Algorithm:
//  1. Format num as a zero-padded 32-character binary string ("%032b").
//  2. Build the mirrored string: reversed[i] = bits[31-i].
//  3. Parse the mirrored string back with base-2 ParseUint.
//
// Time:  O(32) = O(1) — fixed-size format, mirror, and parse passes.
// Space: O(32) = O(1) — two fixed 32-byte buffers (heap-allocating, unlike Approaches 2-4).
func bruteForce(num uint32) uint32 {
	bits := fmt.Sprintf("%032b", num) // zero-padded so all 32 positions are explicit
	reversed := make([]byte, 32)
	for i := 0; i < 32; i++ {
		reversed[i] = bits[31-i] // mirror character positions around the middle
	}
	// error can't occur: the string is exactly 32 chars of '0'/'1'
	v, _ := strconv.ParseUint(string(reversed), 2, 32)
	return uint32(v)
}

// ── Approach 2: Bit by Bit ───────────────────────────────────────────────────
//
// bitByBit solves Reverse Bits by peeling the lowest bit off the input and
// pushing it onto the result, 32 times.
//
// Intuition:
//
//	Reversing is "last in, first out" — exactly what repeated shifting gives:
//	the bit popped from num's little end is pushed onto result's little end,
//	so the first-popped (lowest) bit ends up highest. After 32 pops the whole
//	number is mirrored, like reversing a list by pushing onto a stack.
//
// Algorithm:
//  1. result = 0.
//  2. Repeat 32 times: shift result left one, OR in num's lowest bit,
//     shift num right one.
//
// Time:  O(32) = O(1) — one constant-work iteration per bit.
// Space: O(1) — a single accumulator, no allocations.
func bitByBit(num uint32) uint32 {
	var result uint32
	for i := 0; i < 32; i++ {
		result <<= 1      // make room at the bottom for the next bit
		result |= num & 1 // take num's current lowest bit
		num >>= 1         // consume it
	}
	return result
}

// ── Approach 3: Byte Lookup Table (Follow-Up: Called Many Times) ─────────────
//
// byteTable solves Reverse Bits by reversing each of the 4 bytes via a
// precomputed 256-entry table and swapping the bytes into mirrored positions.
//
// Intuition:
//
//	If the function is called many times (the follow-up), stop re-deriving
//	per-bit work: precompute the reversal of every possible byte once
//	(256 entries), then answer each query with 4 table lookups. Reversing a
//	32-bit word = reversing each byte AND reversing the order of the bytes,
//	so the lowest input byte lands (bit-reversed) in the highest output slot.
//
// Algorithm:
//  1. Precompute revByte[b] = the 8-bit reversal of b for b in 0..255.
//  2. Split num into 4 bytes; look each up in the table.
//  3. Reassemble in mirrored byte order: low byte → bits 24-31, ..., high
//     byte → bits 0-7.
//
// Time:  O(1) per call — 4 lookups and 3 ORs (after a one-time 256-entry build).
// Space: O(256) = O(1) — the shared lookup table, amortised across all calls.
func byteTable(num uint32) uint32 {
	return revByte[num&0xff]<<24 | // lowest byte, reversed, becomes highest
		revByte[(num>>8)&0xff]<<16 | // 2nd byte → 3rd slot
		revByte[(num>>16)&0xff]<<8 | // 3rd byte → 2nd slot
		revByte[num>>24] // highest byte, reversed, becomes lowest
}

// revByte[b] holds b with its 8 bits reversed; built once at program start.
var revByte = buildRevByteTable()

// buildRevByteTable computes the 8-bit reversal of every byte value 0..255
// using the bit-by-bit method (only 256 × 8 steps, done a single time).
func buildRevByteTable() [256]uint32 {
	var table [256]uint32
	for b := 0; b < 256; b++ {
		r := 0
		for i := 0; i < 8; i++ {
			r = r<<1 | b>>i&1 // push b's i-th bit onto r (same LIFO idea as bitByBit)
		}
		table[b] = uint32(r)
	}
	return table
}

// ── Approach 4: Divide and Conquer Mask-Swap (Optimal) ───────────────────────
//
// divideAndConquer solves Reverse Bits with five mask-and-shift swaps:
// halves, bytes, nibbles, pairs, then single bits.
//
// Intuition:
//
//	A reversed word is: the two 16-bit halves swapped, with each half itself
//	reversed. Apply that definition recursively — swap halves, then swap the
//	bytes inside each half, then nibbles inside each byte, then bit-pairs,
//	then adjacent bits. Each level swaps ALL blocks of one size in parallel
//	using constant masks (0xff00ff00 picks every high byte, 0xf0f0f0f0 every
//	high nibble, 0xcccccccc every high pair, 0xaaaaaaaa every odd bit), so
//	log2(32) = 5 operations reverse the whole word — no loop at all.
//
// Algorithm:
//  1. Swap the 16-bit halves.
//  2. Swap adjacent bytes within each half.
//  3. Swap adjacent nibbles within each byte.
//  4. Swap adjacent 2-bit pairs within each nibble.
//  5. Swap adjacent single bits within each pair.
//
// Time:  O(1) — exactly 5 shift/mask/OR lines, branch-free and loop-free.
// Space: O(1) — the value is transformed in a register.
func divideAndConquer(num uint32) uint32 {
	num = (num >> 16) | (num << 16)                             // 1) swap the two 16-bit halves
	num = ((num & 0xff00ff00) >> 8) | ((num & 0x00ff00ff) << 8) // 2) swap bytes inside each half
	num = ((num & 0xf0f0f0f0) >> 4) | ((num & 0x0f0f0f0f) << 4) // 3) swap nibbles inside each byte
	num = ((num & 0xcccccccc) >> 2) | ((num & 0x33333333) << 2) // 4) swap 2-bit pairs inside each nibble
	num = ((num & 0xaaaaaaaa) >> 1) | ((num & 0x55555555) << 1) // 5) swap adjacent bits inside each pair
	return num
}

func main() {
	// Example 1: n = 00000010100101000001111010011100 (43261596)
	//            → 00111001011110000101001010000000 (964176192)
	// Example 2: n = 11111111111111111111111111111101 (4294967293)
	//            → 10111111111111111111111111111101 (3221225471)
	ex1 := uint32(0b00000010100101000001111010011100)
	ex2 := uint32(0b11111111111111111111111111111101)

	fmt.Println("=== Approach 1: Brute Force (Binary String Round-Trip) ===")
	fmt.Println(bruteForce(ex1)) // expected: 964176192
	fmt.Println(bruteForce(ex2)) // expected: 3221225471

	fmt.Println("=== Approach 2: Bit by Bit ===")
	fmt.Println(bitByBit(ex1)) // expected: 964176192
	fmt.Println(bitByBit(ex2)) // expected: 3221225471

	fmt.Println("=== Approach 3: Byte Lookup Table ===")
	fmt.Println(byteTable(ex1)) // expected: 964176192
	fmt.Println(byteTable(ex2)) // expected: 3221225471

	fmt.Println("=== Approach 4: Divide and Conquer Mask-Swap (Optimal) ===")
	fmt.Println(divideAndConquer(ex1)) // expected: 964176192
	fmt.Println(divideAndConquer(ex2)) // expected: 3221225471
}
