package main

import "fmt"

// ── Approach 1: Brute Force ──────────────────────────────────────────────────
//
// bruteForce solves Bitwise AND of Numbers Range by literally ANDing every
// number from left to right.
//
// Intuition:
//
//	The definition of the problem is an accumulation: result = left & (left+1)
//	& ... & right. Just run the loop. One saving grace makes it survivable on
//	many inputs: AND can only ever turn bits OFF, so the moment the running
//	result reaches 0 it can never recover — bail out immediately.
//
// Algorithm:
//  1. result = left.
//  2. For num = left+1 .. right: result &= num.
//  3. If result becomes 0, return 0 early (AND is monotonically decreasing).
//  4. Return result.
//
// Time:  O(right − left) worst case — e.g. [2³⁰, 2³¹−1] never hits 0 and
//
//	iterates ~10⁹ times; the early exit only saves ranges that cross a
//	power of two quickly.
//
// Space: O(1) — a single accumulator.
func bruteForce(left int, right int) int {
	result := left // accumulate the AND starting from the low end
	for num := left + 1; num <= right; num++ {
		result &= num // fold the next number into the running AND
		if result == 0 {
			return 0 // bits only turn off — once zero, always zero
		}
	}
	return result
}

// ── Approach 2: Common Prefix (Bit Shifting) ─────────────────────────────────
//
// commonPrefixShift solves Bitwise AND of Numbers Range by finding the common
// binary prefix of left and right.
//
// Intuition:
//
//	Below the highest bit where left and right differ, every bit pattern
//	occurs somewhere in [left, right]: counting up from left to right must
//	flip that differing bit, and on the way each lower bit passes through 0.
//	So every bit below the common prefix is ANDed to 0, and the answer is
//	exactly the shared prefix of left and right, zero-padded on the right.
//
// Algorithm:
//  1. shift = 0.
//  2. While left != right: halve both (drop the lowest bit) and shift++.
//  3. When they meet, the surviving value is the common prefix; restore its
//     position with left << shift.
//
// Time:  O(log right) — at most 31 shifts for 32-bit inputs.
// Space: O(1) — two scalars.
func commonPrefixShift(left int, right int) int {
	shift := 0 // how many low bits we discarded
	// Keep chopping the lowest bit until the remaining prefixes agree.
	for left != right {
		left >>= 1  // drop lowest bit of left
		right >>= 1 // drop lowest bit of right
		shift++     // remember how far we shifted
	}
	// left == right == common prefix; move it back to its true position,
	// filling the discarded (differing) bits with zeros.
	return left << shift
}

// ── Approach 3: Brian Kernighan's Trick (Optimal) ────────────────────────────
//
// brianKernighan solves Bitwise AND of Numbers Range by clearing the lowest
// set bit of right until right no longer exceeds left.
//
// Intuition:
//
//	right &= right−1 erases the lowest set bit of right in one operation.
//	Any bit set in right below the common prefix cannot survive the AND
//	(some number in the range has a 0 there), so we may keep erasing
//	right's trailing bits while right > left. The first value of right that
//	is ≤ left consists purely of prefix bits shared with left — and a prefix
//	of right that is ≤ left must be a prefix of left too, i.e. exactly the
//	common prefix. That is the answer.
//
// Algorithm:
//  1. While left < right: right &= right − 1 (drop right's lowest set bit).
//  2. Return right.
//
// Time:  O(popcount(right)) ≤ O(log right) — each step removes one set bit,
//
//	and it usually stops well before all of them are gone.
//
// Space: O(1) — in-place on the two parameters.
func brianKernighan(left int, right int) int {
	// Erase right's lowest set bit until right sinks down to (or below) left.
	for left < right {
		right &= right - 1 // Kernighan: clears exactly the lowest set bit
	}
	return right // now the shared prefix of the original [left, right]
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Printf("left=5, right=7                    got=%d  expected 4\n", bruteForce(5, 7))
	fmt.Printf("left=0, right=0                    got=%d  expected 0\n", bruteForce(0, 0))
	fmt.Printf("left=1, right=2147483647           got=%d  expected 0\n", bruteForce(1, 2147483647)) // early exit saves this one

	fmt.Println("=== Approach 2: Common Prefix (Bit Shifting) ===")
	fmt.Printf("left=5, right=7                    got=%d  expected 4\n", commonPrefixShift(5, 7))
	fmt.Printf("left=0, right=0                    got=%d  expected 0\n", commonPrefixShift(0, 0))
	fmt.Printf("left=1, right=2147483647           got=%d  expected 0\n", commonPrefixShift(1, 2147483647))
	fmt.Printf("left=2147483646, right=2147483647  got=%d  expected 2147483646\n", commonPrefixShift(2147483646, 2147483647)) // adjacent-pair edge

	fmt.Println("=== Approach 3: Brian Kernighan's Trick (Optimal) ===")
	fmt.Printf("left=5, right=7                    got=%d  expected 4\n", brianKernighan(5, 7))
	fmt.Printf("left=0, right=0                    got=%d  expected 0\n", brianKernighan(0, 0))
	fmt.Printf("left=1, right=2147483647           got=%d  expected 0\n", brianKernighan(1, 2147483647))
	fmt.Printf("left=2147483646, right=2147483647  got=%d  expected 2147483646\n", brianKernighan(2147483646, 2147483647))
}
