package main

import "fmt"

// ── Approach 1: Brute Force (Digit-by-Digit Search) ──────────────────────────
//
// bruteForce solves Excel Sheet Column Title by determining the answer's
// length first, then choosing each letter greedily from the left.
//
// Intuition:
//
//	Titles of length k cover a contiguous range of numbers: length 1 covers
//	1..26, length 2 covers 27..702, and in general length k covers the next
//	26^k values. So we can first find the answer's length by subtracting
//	block sizes, then pick each letter left to right: the first letter
//	partitions the remaining range into 26 equal chunks of size 26^(k-1),
//	and simple division tells us which chunk (letter) the number falls in.
//
// Algorithm:
//  1. Find length k: subtract 26, 26^2, 26^3 ... while columnNumber remains
//     positive; convert to a 0-based offset within the length-k block.
//  2. For each position from left to right, chunk = 26^(remaining letters);
//     letter index = offset / chunk; offset %= chunk.
//  3. Append 'A'+index for every position.
//
// Time:  O(log26(n)^2 ) — k ≤ 7 positions, each computing a 26^i power
//
//	(tiny in practice; effectively O(k^2) with k ≤ 7).
//
// Space: O(log26(n)) — the k output letters.
func bruteForce(columnNumber int) string {
	// Step 1: find the title length k and the 0-based offset inside that block.
	offset := columnNumber
	blockSize := 26 // number of titles having the current length (26^k)
	k := 1
	for offset > blockSize {
		offset -= blockSize // skip all titles of length k
		blockSize *= 26     // next block: titles of length k+1
		k++
	}
	offset-- // make the offset 0-based within the length-k block

	// Step 2: pick letters left to right by dividing into 26 equal chunks.
	out := make([]byte, k)
	// chunk = 26^(k-1): how many titles share the same first letter.
	chunk := 1
	for i := 0; i < k-1; i++ {
		chunk *= 26
	}
	for i := 0; i < k; i++ {
		out[i] = byte('A' + offset/chunk) // which of the 26 chunks we are in
		offset %= chunk                   // descend into that chunk
		if chunk > 1 {
			chunk /= 26 // next position partitions 26x finer
		}
	}
	return string(out)
}

// ── Approach 2: Iterative Base-26 with Offset (Optimal) ──────────────────────
//
// iterativeBase26 solves Excel Sheet Column Title as a base-26 conversion
// where digits run 1..26 instead of 0..25.
//
// Intuition:
//
//	This is base-26, except there is no zero digit: A=1 ... Z=26. The classic
//	fix for such "bijective numeration" is to subtract 1 before every mod/div
//	step, which shifts the digit range from 1..26 down to 0..25 so ordinary
//	% 26 and / 26 work. Digits come out least-significant first, so build the
//	string backwards.
//
// Algorithm:
//  1. While columnNumber > 0:
//  2. columnNumber-- (shift 1..26 → 0..25).
//  3. Emit letter 'A' + columnNumber%26 (prepend, or append and reverse).
//  4. columnNumber /= 26.
//
// Time:  O(log26(n)) — one letter produced per loop iteration (≤ 7 for int32).
// Space: O(log26(n)) — the output buffer.
func iterativeBase26(columnNumber int) string {
	out := []byte{}
	for columnNumber > 0 {
		columnNumber-- // shift digits from 1..26 to 0..25 (no zero in this system)
		// Least-significant "digit" is the last letter of the title.
		out = append(out, byte('A'+columnNumber%26))
		columnNumber /= 26 // drop the digit we just emitted
	}
	// Digits were produced right-to-left → reverse in place.
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return string(out)
}

// ── Approach 3: Recursive Base-26 ────────────────────────────────────────────
//
// recursiveBase26 solves Excel Sheet Column Title with the same subtract-one
// trick expressed recursively.
//
// Intuition:
//
//	The title of n is the title of (n-1)/26 followed by the letter for
//	(n-1)%26. Recursion emits the most-significant letters first, so no
//	reversal is needed — the call stack does the reordering.
//
// Algorithm:
//  1. Base case: n == 0 → empty string.
//  2. Recurse on (n-1)/26, then append 'A' + (n-1)%26.
//
// Time:  O(log26(n)) — one call per output letter.
// Space: O(log26(n)) — recursion depth (≤ 7 frames for int32 inputs).
func recursiveBase26(columnNumber int) string {
	// Base case: nothing left to convert.
	if columnNumber == 0 {
		return ""
	}
	columnNumber-- // bijective base-26: shift 1..26 → 0..25 before splitting
	// Prefix letters first (recursion), then this position's letter.
	return recursiveBase26(columnNumber/26) + string(byte('A'+columnNumber%26))
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Digit-by-Digit Search) ===")
	fmt.Printf("columnNumber=1           got=%q  expected \"A\"\n", bruteForce(1))
	fmt.Printf("columnNumber=28          got=%q  expected \"AB\"\n", bruteForce(28))
	fmt.Printf("columnNumber=701         got=%q  expected \"ZY\"\n", bruteForce(701))
	fmt.Printf("columnNumber=26          got=%q  expected \"Z\"\n", bruteForce(26))               // last 1-letter title
	fmt.Printf("columnNumber=27          got=%q  expected \"AA\"\n", bruteForce(27))              // first 2-letter title
	fmt.Printf("columnNumber=2147483647  got=%q  expected \"FXSHRXW\"\n", bruteForce(2147483647)) // max int32

	fmt.Println("=== Approach 2: Iterative Base-26 with Offset (Optimal) ===")
	fmt.Printf("columnNumber=1           got=%q  expected \"A\"\n", iterativeBase26(1))
	fmt.Printf("columnNumber=28          got=%q  expected \"AB\"\n", iterativeBase26(28))
	fmt.Printf("columnNumber=701         got=%q  expected \"ZY\"\n", iterativeBase26(701))
	fmt.Printf("columnNumber=26          got=%q  expected \"Z\"\n", iterativeBase26(26))
	fmt.Printf("columnNumber=27          got=%q  expected \"AA\"\n", iterativeBase26(27))
	fmt.Printf("columnNumber=2147483647  got=%q  expected \"FXSHRXW\"\n", iterativeBase26(2147483647))

	fmt.Println("=== Approach 3: Recursive Base-26 ===")
	fmt.Printf("columnNumber=1           got=%q  expected \"A\"\n", recursiveBase26(1))
	fmt.Printf("columnNumber=28          got=%q  expected \"AB\"\n", recursiveBase26(28))
	fmt.Printf("columnNumber=701         got=%q  expected \"ZY\"\n", recursiveBase26(701))
	fmt.Printf("columnNumber=26          got=%q  expected \"Z\"\n", recursiveBase26(26))
	fmt.Printf("columnNumber=27          got=%q  expected \"AA\"\n", recursiveBase26(27))
	fmt.Printf("columnNumber=2147483647  got=%q  expected \"FXSHRXW\"\n", recursiveBase26(2147483647))
}
