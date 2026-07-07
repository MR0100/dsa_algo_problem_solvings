package main

import (
	"fmt"
	"math/bits"
)

// hammingDistance returns the Hamming distance between two integers: the number
// of bit positions at which they differ. x ^ y has a 1 exactly where they
// disagree, so we just count set bits of the XOR.
func hammingDistance(x, y int) int {
	return bits.OnesCount(uint(x ^ y)) // popcount of the disagreement mask
}

// ── Approach 1: Brute Force (All Pairs) ──────────────────────────────────────
//
// bruteForce solves Total Hamming Distance by summing the Hamming distance of
// every unordered pair (i, j) with i < j.
//
// Intuition:
//
//	The definition is literally "sum over all pairs". Enumerate each pair once,
//	compute its Hamming distance via popcount(nums[i] ^ nums[j]), and add.
//
// Algorithm:
//  1. total = 0.
//  2. For i = 0..n-1, for j = i+1..n-1: total += hammingDistance(nums[i], nums[j]).
//  3. Return total.
//
// Time:  O(n² · w) where w ≤ 32 is the machine word — n²/2 pairs, each a
//
//	constant-width popcount. With n up to 10⁴ this is ~10⁸ and will TLE.
//
// Space: O(1).
func bruteForce(nums []int) int {
	total := 0
	for i := 0; i < len(nums); i++ { // pick the first element of the pair
		for j := i + 1; j < len(nums); j++ { // pair it with every later element
			total += hammingDistance(nums[i], nums[j]) // add this pair's distance
		}
	}
	return total
}

// ── Approach 2: Bit-Position Counting (Optimal) ──────────────────────────────
//
// bitColumnCount solves Total Hamming Distance by looking at each of the 32 bit
// positions independently and counting its contribution across all pairs.
//
// Intuition:
//
//	Total Hamming distance = Σ over all pairs of (number of differing bits).
//	Swap the order of summation: instead of "for each pair, count differing
//	bits", do "for each bit position, count differing pairs". At a fixed bit
//	position, a pair differs iff one number has a 0 there and the other a 1.
//	If `ones` numbers have a 1 at that position and the rest (`n - ones`) have a
//	0, the number of differing pairs is exactly ones · (n - ones). Sum that over
//	all 32 positions.
//
// Algorithm:
//  1. n = len(nums), total = 0.
//  2. For each bit position b in 0..31:
//     - ones = count of nums with bit b set (shift right b, AND 1).
//     - total += ones * (n - ones).
//  3. Return total.
//
// Time:  O(n · w), w = 32 → O(n). One pass per bit column.
// Space: O(1).
func bitColumnCount(nums []int) int {
	n := len(nums)
	total := 0
	for b := 0; b < 32; b++ { // examine each bit column independently
		ones := 0 // how many numbers have a 1 in this column
		for _, x := range nums {
			ones += (x >> b) & 1 // add 1 when bit b of x is set
		}
		zeros := n - ones     // the rest have a 0 in this column
		total += ones * zeros // each (one,zero) pair differs here → contributes 1
	}
	return total
}

// ── Approach 3: Bit-Position Counting, Early Exit ────────────────────────────
//
// bitColumnEarlyExit is the same column-counting idea but stops scanning bit
// positions once every remaining number is 0 there (no higher bits are set in
// any element), avoiding always looping the full 32 columns.
//
// Intuition:
//
//	The values are ≤ 10⁹ < 2³⁰, so high columns are usually all-zero and
//	contribute nothing (ones·zeros = 0·n = 0). We can OR all numbers together
//	to learn the highest bit that appears, then only iterate columns up to it.
//	This does not change the asymptotic bound but trims constant work.
//
// Algorithm:
//  1. OR all numbers into `orAll`; the highest set bit of orAll is the top
//     column worth scanning.
//  2. For each column b while (orAll >> b) != 0: accumulate ones·(n-ones) as
//     before, shifting every number down as we go.
//  3. Return total.
//
// Time:  O(n · maxBit) ≤ O(n · 30). Space: O(1).
func bitColumnEarlyExit(nums []int) int {
	n := len(nums)
	orAll := 0
	for _, x := range nums {
		orAll |= x // union of all bits present anywhere in the array
	}
	total := 0
	for b := 0; orAll>>b != 0; b++ { // stop once no number has a bit at/above b
		ones := 0
		for _, x := range nums {
			ones += (x >> b) & 1 // count 1s in column b
		}
		total += ones * (n - ones) // differing pairs contributed by column b
	}
	return total
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (All Pairs) ===")
	fmt.Printf("nums=[4,14,2]    got=%d  expected 6\n", bruteForce([]int{4, 14, 2}))
	fmt.Printf("nums=[4,14,4]    got=%d  expected 4\n", bruteForce([]int{4, 14, 4}))
	fmt.Printf("nums=[0]         got=%d  expected 0\n", bruteForce([]int{0}))

	fmt.Println("=== Approach 2: Bit-Position Counting (Optimal) ===")
	fmt.Printf("nums=[4,14,2]    got=%d  expected 6\n", bitColumnCount([]int{4, 14, 2}))
	fmt.Printf("nums=[4,14,4]    got=%d  expected 4\n", bitColumnCount([]int{4, 14, 4}))
	fmt.Printf("nums=[0]         got=%d  expected 0\n", bitColumnCount([]int{0}))

	fmt.Println("=== Approach 3: Bit-Position Counting, Early Exit ===")
	fmt.Printf("nums=[4,14,2]    got=%d  expected 6\n", bitColumnEarlyExit([]int{4, 14, 2}))
	fmt.Printf("nums=[4,14,4]    got=%d  expected 4\n", bitColumnEarlyExit([]int{4, 14, 4}))
	fmt.Printf("nums=[0]         got=%d  expected 0\n", bitColumnEarlyExit([]int{0}))
}
