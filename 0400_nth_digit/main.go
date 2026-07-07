package main

import (
	"fmt"
	"strconv"
)

// ── Approach 1: Brute Force ──────────────────────────────────────────────────
//
// bruteForce solves Nth Digit by generating the sequence digit by digit until
// it has emitted n digits.
//
// Intuition:
//
//	The sequence is just 1,2,3,... written out and concatenated. Walk the
//	integers in order, converting each to its decimal string, and decrement n
//	by the number's digit count; when n falls within the current number, the
//	target digit is inside it.
//
// Algorithm:
//  1. num = 1.
//  2. Loop: s = decimal string of num; if n <= len(s), the (n-1)-th char of s
//     is the answer; else n -= len(s), num++.
//
// Time:  O(n / averageDigitLen) ≈ O(n) — may iterate up to ~n numbers.
// Space: O(1) — reuses one small string.
func bruteForce(n int) int {
	num := 1
	for {
		s := strconv.Itoa(num) // decimal digits of the current number
		if n <= len(s) {
			// The n-th remaining digit is the (n-1)-th char of this number.
			return int(s[n-1] - '0')
		}
		n -= len(s) // skip all of this number's digits
		num++
	}
}

// ── Approach 2: Math / Digit Blocks (Optimal) ────────────────────────────────
//
// mathBlocks solves Nth Digit by jumping over whole blocks of equal-length
// numbers instead of counting one number at a time.
//
// Intuition:
//
//	Numbers group by length: there are 9 one-digit numbers (1..9), 90 two-digit
//	(10..99), 900 three-digit, ... In general 9·10^(len-1) numbers of a given
//	length, contributing len·9·10^(len-1) digits. Subtract these block sizes
//	from n until n lands inside a block; then arithmetic pinpoints exactly which
//	number and which digit within it.
//
// Algorithm:
//  1. length = 1, count = 9, start = 1.
//  2. While n > length*count: n -= length*count; length++; count *= 10;
//     start *= 10.
//  3. The target number is start + (n-1)/length.
//  4. The digit index within that number is (n-1)%length; return it.
//
// Time:  O(log n) — one iteration per digit-length block (≤ 10 for 32-bit n).
// Space: O(1) — a few integer accumulators.
func mathBlocks(n int) int {
	length := 1 // current block's number length (digits per number)
	count := 9  // how many numbers have this length: 9,90,900,...
	start := 1  // first number of this length: 1,10,100,...
	// Skip whole blocks while n overshoots the digits this block provides.
	for n > length*count {
		n -= length * count // consume this block's digits
		length++            // move to longer numbers
		count *= 10         // 10x as many of them
		start *= 10         // block now starts at the next power of ten
	}
	// n is now the 1-based offset within the current block.
	// Which number holds it: (n-1)/length steps past `start`.
	number := start + (n-1)/length
	// Which digit inside that number: (n-1)%length from the left.
	digitIndex := (n - 1) % length
	s := strconv.Itoa(number)
	return int(s[digitIndex] - '0')
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Printf("n=3           got=%d  expected 3\n", bruteForce(3))
	fmt.Printf("n=11          got=%d  expected 0\n", bruteForce(11))

	fmt.Println("=== Approach 2: Math / Digit Blocks (Optimal) ===")
	fmt.Printf("n=3           got=%d  expected 3\n", mathBlocks(3))
	fmt.Printf("n=11          got=%d  expected 0\n", mathBlocks(11))
	fmt.Printf("n=2147483647  got=%d  expected 2\n", mathBlocks(2147483647)) // large-n edge
}
