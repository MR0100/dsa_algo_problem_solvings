package main

import (
	"fmt"
	"math/bits"
)

// ── Approach 1: Brute Force (Bit-by-Bit Mask Scan) ───────────────────────────
//
// bruteForce solves Number of 1 Bits by testing every one of the 32 bit
// positions with a shifting mask.
//
// Intuition:
//
//	A uint32 has exactly 32 bit slots. If we peek into each slot and ask
//	"is this bit a 1?", the number of yes-answers is the Hamming weight.
//	No cleverness — inspect everything.
//
// Algorithm:
//  1. Initialise count = 0.
//  2. For every position i in 0..31, build the mask 1<<i (only bit i set).
//  3. AND the mask with n; a non-zero result means bit i of n is 1 → count++.
//  4. Return count.
//
// Time:  O(32) = O(1) — always exactly 32 iterations, independent of the value.
// Space: O(1) — only a counter and a mask.
func bruteForce(n uint32) int {
	count := 0                // running total of set bits found so far
	for i := 0; i < 32; i++ { // visit every bit position 0..31
		mask := uint32(1) << i // mask with only bit i set, e.g. i=3 → 0b1000
		if n&mask != 0 {       // AND isolates bit i; non-zero ⇒ that bit is 1
			count++ // bit i is set → record it
		}
	}
	return count
}

// ── Approach 2: Brian Kernighan's Trick ──────────────────────────────────────
//
// brianKernighan solves Number of 1 Bits by repeatedly erasing the lowest set
// bit with n & (n-1), looping once per SET bit instead of once per position.
//
// Intuition:
//
//	Subtracting 1 from n flips the lowest set bit to 0 and turns every bit
//	below it to 1 (e.g. 0b10100 - 1 = 0b10011). AND-ing n with n-1 therefore
//	wipes out exactly the lowest set bit and nothing else. Count how many
//	wipes it takes to reach 0 — that is the number of set bits.
//
// Algorithm:
//  1. Initialise count = 0.
//  2. While n != 0: replace n with n & (n-1) (drops lowest set bit), count++.
//  3. Return count.
//
// Time:  O(k) where k = number of set bits (≤ 32) — beats Approach 1 for
//
//	sparse numbers because zero bits are never visited.
//
// Space: O(1) — in-place bit arithmetic.
func brianKernighan(n uint32) int {
	count := 0
	for n != 0 { // one iteration per set bit, not per bit position
		n &= n - 1 // n-1 flips the lowest set bit and all zeros below it;
		//            AND-ing erases exactly that lowest set bit from n
		count++ // one set bit removed ⇒ one set bit counted
	}
	return count
}

// table8 caches the popcount of every possible byte value 0..255.
// Built once at start-up; afterwards any 32-bit popcount costs 4 lookups.
var table8 [256]int

// init fills table8 using the recurrence popcount(i) = popcount(i>>1) + (i&1):
// dropping the last bit of i gives a smaller, already-solved value.
func init() {
	for i := 1; i < 256; i++ {
		table8[i] = table8[i>>1] + (i & 1) // reuse the answer for i/2, add i's last bit
	}
}

// ── Approach 3: Lookup Table (8-Bit Chunks) ──────────────────────────────────
//
// lookupTable solves Number of 1 Bits by splitting n into four bytes and
// summing precomputed per-byte popcounts.
//
// Intuition:
//
//	This is the answer to the follow-up "what if the function is called many
//	times?": pay once to precompute popcount for all 256 byte values, then
//	every future query is just 4 table reads + 3 additions — no per-bit work.
//
// Algorithm:
//  1. (Once, at start-up) build table8[b] = popcount(b) for b in 0..255.
//  2. Slice n into 4 bytes with shifts and & 0xFF.
//  3. Return the sum of the 4 table entries.
//
// Time:  O(1) — exactly 4 lookups per call (O(256) one-time precomputation).
// Space: O(256) = O(1) — the fixed-size table, shared by all calls.
func lookupTable(n uint32) int {
	return table8[n&0xFF] + // byte 0: bits 0..7
		table8[(n>>8)&0xFF] + // byte 1: bits 8..15
		table8[(n>>16)&0xFF] + // byte 2: bits 16..23
		table8[(n>>24)&0xFF] // byte 3: bits 24..31
}

// ── Approach 4: Parallel Bit Count (SWAR / Divide and Conquer) ───────────────
//
// parallelCount solves Number of 1 Bits with the branch-free SWAR technique
// ("SIMD Within A Register"): it sums adjacent bit groups in parallel,
// doubling the group width each step — a divide-and-conquer on the bits.
//
// Intuition:
//
//	Treat the 32-bit word as 16 two-bit counters, then 8 four-bit counters,
//	then 4 byte counters. Each masked add merges neighbouring counters, so
//	after log2(32)=5 conceptual rounds the whole popcount is known. Constant
//	masks let one machine add process many counters simultaneously — this is
//	essentially how hardware POPCNT and math/bits work.
//
// Algorithm:
//  1. n - ((n>>1) & 0x55555555): every 2-bit field now holds the popcount
//     of the 2 bits it replaced (uses the identity popcount(ab) = ab - a).
//  2. (n & 0x33333333) + ((n>>2) & 0x33333333): merge into 4-bit field sums.
//  3. (n + (n>>4)) & 0x0F0F0F0F: merge into per-byte sums (each ≤ 8).
//  4. (n * 0x01010101) >> 24: the multiply adds all four bytes into the top
//     byte; shifting it down yields the total.
//
// Time:  O(1) — a fixed sequence of ~12 arithmetic ops, no loops or branches.
// Space: O(1) — pure register arithmetic.
func parallelCount(n uint32) int {
	n = n - ((n >> 1) & 0x55555555)                // each 2-bit field := popcount of its 2 bits
	n = (n & 0x33333333) + ((n >> 2) & 0x33333333) // each 4-bit field := sum of its two 2-bit fields
	n = (n + (n >> 4)) & 0x0F0F0F0F                // each byte := sum of its two nibbles (≤ 8, no overflow)
	return int((n * 0x01010101) >> 24)             // top byte accumulates b0+b1+b2+b3; shift it down
}

// ── Approach 5: Built-in Popcount (Optimal) ──────────────────────────────────
//
// builtinPopcount solves Number of 1 Bits with Go's math/bits.OnesCount32,
// which the compiler lowers to the hardware POPCNT instruction where available.
//
// Intuition:
//
//	Population count is so common that CPUs implement it in silicon and Go
//	exposes it as an intrinsic. In production code this is the correct answer;
//	in an interview, mention it after deriving Approaches 2–4 by hand.
//
// Algorithm:
//  1. Return bits.OnesCount32(n).
//
// Time:  O(1) — a single CPU instruction on amd64/arm64 (SWAR fallback otherwise).
// Space: O(1) — no allocations.
func builtinPopcount(n uint32) int {
	return bits.OnesCount32(n) // compiler intrinsic → hardware POPCNT
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Bit-by-Bit Mask Scan) ===")
	fmt.Println(bruteForce(11))         // 3  (1011 has three set bits)
	fmt.Println(bruteForce(128))        // 1  (10000000 has one set bit)
	fmt.Println(bruteForce(2147483645)) // 30 (1111111111111111111111111111101)

	fmt.Println("=== Approach 2: Brian Kernighan's Trick ===")
	fmt.Println(brianKernighan(11))         // 3
	fmt.Println(brianKernighan(128))        // 1
	fmt.Println(brianKernighan(2147483645)) // 30

	fmt.Println("=== Approach 3: Lookup Table (8-Bit Chunks) ===")
	fmt.Println(lookupTable(11))         // 3
	fmt.Println(lookupTable(128))        // 1
	fmt.Println(lookupTable(2147483645)) // 30

	fmt.Println("=== Approach 4: Parallel Bit Count (SWAR) ===")
	fmt.Println(parallelCount(11))         // 3
	fmt.Println(parallelCount(128))        // 1
	fmt.Println(parallelCount(2147483645)) // 30

	fmt.Println("=== Approach 5: Built-in Popcount (Optimal) ===")
	fmt.Println(builtinPopcount(11))         // 3
	fmt.Println(builtinPopcount(128))        // 1
	fmt.Println(builtinPopcount(2147483645)) // 30
}
