package main

import (
	"fmt"
	"math/bits"
)

// ── Approach 1: Brute Force (Kernighan per number) ───────────────────────────
//
// bruteForce solves Counting Bits by counting set bits of every value from 0
// to n independently with Brian Kernighan's lowest-set-bit loop.
//
// Intuition:
//
//	The definition asks for popcount(i) for each i in [0, n]. The simplest
//	honest answer is to count the 1-bits of each number on its own. Kernighan's
//	trick x &= x-1 clears the lowest set bit, so the number of iterations for a
//	given x equals popcount(x) — cheaper than testing all 32 bit positions.
//
// Algorithm:
//  1. Allocate ans of length n+1.
//  2. For each i in [0, n], count its set bits by repeatedly clearing the
//     lowest set bit until the value is 0.
//  3. Store the count in ans[i].
//
// Time:  O(n · k) where k = average popcount ≤ 32 → effectively O(n log n).
// Space: O(1) extra (besides the required output array).
func bruteForce(n int) []int {
	ans := make([]int, n+1) // ans[i] will hold popcount(i)
	for i := 0; i <= n; i++ {
		count := 0 // number of set bits found in i so far
		x := i     // work on a copy so i stays intact
		for x > 0 {
			x &= x - 1 // Kernighan: erase the lowest set bit
			count++    // each erase corresponds to exactly one set bit
		}
		ans[i] = count // record the popcount of i
	}
	return ans
}

// ── Approach 2: DP with Highest Power of Two (offset) ────────────────────────
//
// dpHighBit solves Counting Bits by reusing the answer of a smaller number
// obtained by stripping the highest power of two.
//
// Intuition:
//
//	If offset is the largest power of two ≤ i, then i has exactly one more set
//	bit than i - offset (we removed that single high bit). So
//	ans[i] = ans[i-offset] + 1, and offset doubles each time we reach a new
//	power of two.
//
// Algorithm:
//  1. offset = 1 (current highest power of two).
//  2. For i from 1 to n: if i == offset*2, offset doubles.
//  3. ans[i] = ans[i-offset] + 1.
//
// Time:  O(n) — one O(1) step per number.
// Space: O(1) extra beyond the output array.
func dpHighBit(n int) []int {
	ans := make([]int, n+1) // ans[0] = 0 by default
	offset := 1             // highest power of two seen so far (starts at 1 = 2^0)
	for i := 1; i <= n; i++ {
		if offset*2 == i {
			offset *= 2 // i just reached the next power of two → update the high bit
		}
		ans[i] = ans[i-offset] + 1 // i = (i-offset) plus one extra high bit
	}
	return ans
}

// ── Approach 3: DP with Right Shift (x >> 1) ─────────────────────────────────
//
// dpRightShift solves Counting Bits using the relation between i and i/2.
//
// Intuition:
//
//	i >> 1 drops i's lowest bit. That lowest bit is i & 1. So
//	popcount(i) = popcount(i >> 1) + (i & 1). Since i>>1 < i, its answer is
//	already computed.
//
// Algorithm:
//  1. For i from 1 to n: ans[i] = ans[i>>1] + (i & 1).
//
// Time:  O(n) — one O(1) step per number.
// Space: O(1) extra beyond the output array.
func dpRightShift(n int) []int {
	ans := make([]int, n+1)
	for i := 1; i <= n; i++ {
		ans[i] = ans[i>>1] + (i & 1) // half's popcount plus the bit we shifted off
	}
	return ans
}

// ── Approach 4: DP with Kernighan (x & (x-1)) (Optimal) ──────────────────────
//
// dpKernighan solves Counting Bits with the cleanest recurrence:
// ans[i] = ans[i & (i-1)] + 1.
//
// Intuition:
//
//	i & (i-1) equals i with its lowest set bit cleared — a strictly smaller
//	number already solved. It has exactly one fewer set bit than i, so
//	ans[i] = ans[i & (i-1)] + 1.
//
// Algorithm:
//  1. For i from 1 to n: ans[i] = ans[i & (i-1)] + 1.
//
// Time:  O(n) — one O(1) step per number.
// Space: O(1) extra beyond the output array.
func dpKernighan(n int) []int {
	ans := make([]int, n+1)
	for i := 1; i <= n; i++ {
		ans[i] = ans[i&(i-1)] + 1 // clear lowest set bit → one fewer 1-bit
	}
	return ans
}

// ── Approach 5: Library popcount (reference) ─────────────────────────────────
//
// libPopcount solves Counting Bits by calling the standard library's
// hardware-backed popcount — a sanity reference, not the intended solution.
//
// Intuition:
//
//	math/bits.OnesCount uses a hardware POPCNT-style routine. Handy to verify
//	the DP answers, though the DP solutions are what an interview wants.
//
// Algorithm:
//  1. For i in [0, n]: ans[i] = bits.OnesCount(uint(i)).
//
// Time:  O(n) — each OnesCount is O(1).
// Space: O(1) extra beyond the output array.
func libPopcount(n int) []int {
	ans := make([]int, n+1)
	for i := 0; i <= n; i++ {
		ans[i] = bits.OnesCount(uint(i)) // constant-time hardware popcount
	}
	return ans
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Kernighan per number) ===")
	fmt.Println(bruteForce(2)) // [0 1 1]
	fmt.Println(bruteForce(5)) // [0 1 1 2 1 2]

	fmt.Println("=== Approach 2: DP with Highest Power of Two ===")
	fmt.Println(dpHighBit(2)) // [0 1 1]
	fmt.Println(dpHighBit(5)) // [0 1 1 2 1 2]

	fmt.Println("=== Approach 3: DP with Right Shift ===")
	fmt.Println(dpRightShift(2)) // [0 1 1]
	fmt.Println(dpRightShift(5)) // [0 1 1 2 1 2]

	fmt.Println("=== Approach 4: DP with Kernighan (Optimal) ===")
	fmt.Println(dpKernighan(2)) // [0 1 1]
	fmt.Println(dpKernighan(5)) // [0 1 1 2 1 2]

	fmt.Println("=== Approach 5: Library popcount (reference) ===")
	fmt.Println(libPopcount(2)) // [0 1 1]
	fmt.Println(libPopcount(5)) // [0 1 1 2 1 2]
}
