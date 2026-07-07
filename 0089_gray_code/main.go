package main

import "fmt"

// ── Approach 1: Bit Manipulation Formula ──────────────────────────────────────
//
// grayCode solves Gray Code by applying the formula i ^ (i >> 1).
//
// Intuition:
//   The n-bit Gray code for integer i is: i XOR (i >> 1).
//   Adjacent values i and i+1 differ in exactly 1 bit because:
//   - If i+1's lowest set bit is at position k, the bits above k flip between
//     i and i+1 by exactly one position in the XOR encoding.
//
//   This generates a valid Gray code sequence for all i = 0..2^n-1.
//
// Time:  O(2^n)
// Space: O(1) — output aside.
func grayCode(n int) []int {
	size := 1 << n // 2^n elements
	result := make([]int, size)
	for i := 0; i < size; i++ {
		result[i] = i ^ (i >> 1) // gray code formula
	}
	return result
}

// ── Approach 2: Recursive Mirror Construction ─────────────────────────────────
//
// grayCodeMirror solves Gray Code by reflecting the (n-1)-bit sequence.
//
// Intuition:
//   The n-bit Gray code is built from the (n-1)-bit Gray code by:
//   1. Prepend 0 to each code in the (n-1)-bit sequence.
//   2. Prepend 1 to each code in the REVERSED (n-1)-bit sequence.
//   Concatenate: this gives 2^n codes where consecutive elements (including
//   the join at the midpoint) differ by exactly 1 bit.
//
//   n=1: [0, 1]
//   n=2: [00,01] reflected → [00,01,11,10] = [0,1,3,2]
//   n=3: [000,001,011,010] reflected → [0,1,3,2,6,7,5,4]
//
// Time:  O(2^n)
// Space: O(2^n) — recursive call stack + output.
func grayCodeMirror(n int) []int {
	if n == 0 {
		return []int{0}
	}
	prev := grayCodeMirror(n - 1)
	result := make([]int, 0, 2*len(prev))
	// first half: prepend 0 (unchanged)
	result = append(result, prev...)
	// second half: reverse prev, prepend 1 (add 1 << (n-1))
	half := 1 << (n - 1)
	for i := len(prev) - 1; i >= 0; i-- {
		result = append(result, prev[i]+half)
	}
	return result
}

func main() {
	fmt.Println("=== Approach 1: Bit Manipulation Formula ===")
	fmt.Printf("n=2  got=%v  expected [0 1 3 2]\n", grayCode(2))
	fmt.Printf("n=1  got=%v  expected [0 1]\n", grayCode(1))
	fmt.Printf("n=3  got=%v  expected [0 1 3 2 6 7 5 4]\n", grayCode(3))

	fmt.Println("=== Approach 2: Recursive Mirror ===")
	fmt.Printf("n=2  got=%v  expected [0 1 3 2]\n", grayCodeMirror(2))
	fmt.Printf("n=1  got=%v  expected [0 1]\n", grayCodeMirror(1))
	fmt.Printf("n=3  got=%v  expected [0 1 3 2 6 7 5 4]\n", grayCodeMirror(3))
}
