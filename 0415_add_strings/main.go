package main

import (
	"fmt"
	"strings"
)

// ── Approach 1: Two Pointers + Carry (Grade-School Addition) (Optimal) ────────
//
// twoPointers solves Add Strings by adding the two numbers digit-by-digit from
// the least-significant end, propagating a carry — exactly how addition is done
// by hand — without ever converting the whole string to an integer.
//
// Intuition:
//
//	Line the two numbers up at their right ends. Walk both from the last digit to
//	the first, summing the two current digits plus any carry. The output digit is
//	(sum % 10) and the new carry is (sum / 10). When one number runs out, treat
//	its missing digits as 0. After the loop, if a carry remains, it becomes a new
//	leading digit. We build the result right-to-left, then reverse it once.
//
// Algorithm:
//  1. i = len(num1)-1, j = len(num2)-1, carry = 0.
//  2. While i >= 0 or j >= 0 or carry > 0:
//     d1 = digit at i (or 0), d2 = digit at j (or 0).
//     sum = d1 + d2 + carry; append sum%10; carry = sum/10; i--, j--.
//  3. Reverse the collected digits into the final string.
//
// Time:  O(max(m, n)) — one pass over the longer number (m,n are the lengths).
// Space: O(max(m, n)) — the output buffer (unavoidable: the sum has that many digits).
func twoPointers(num1 string, num2 string) string {
	i, j := len(num1)-1, len(num2)-1 // start at the least-significant digit of each
	carry := 0                       // running carry into the current column
	var sb strings.Builder           // collects result digits in REVERSE order

	// Continue while either number has digits left OR a carry still needs placing.
	for i >= 0 || j >= 0 || carry > 0 {
		sum := carry // begin the column with the incoming carry
		if i >= 0 {
			sum += int(num1[i] - '0') // ASCII digit → integer value
			i--
		}
		if j >= 0 {
			sum += int(num2[j] - '0')
			j--
		}
		sb.WriteByte(byte(sum%10) + '0') // low digit of the column → output (as ASCII)
		carry = sum / 10                 // high part carries to the next column
	}

	return reverseString(sb.String()) // digits were emitted LSB-first; flip them
}

// reverseString returns s with its bytes reversed. All characters here are ASCII
// digits, so byte-level reversal is safe (no multi-byte runes to worry about).
func reverseString(s string) string {
	b := []byte(s)
	for lo, hi := 0, len(b)-1; lo < hi; lo, hi = lo+1, hi-1 {
		b[lo], b[hi] = b[hi], b[lo] // swap symmetric pairs toward the middle
	}
	return string(b)
}

// ── Approach 2: Pre-Sized Buffer, Fill Back-to-Front (No Reverse) ─────────────
//
// preSizedBuffer solves Add Strings with the same digit-by-digit addition but
// writes directly into a pre-allocated byte buffer from the rightmost cell
// backwards, so no final reverse pass is needed.
//
// Intuition:
//
//	The sum of an m-digit and an n-digit number has at most max(m,n)+1 digits.
//	Allocate exactly that many cells and fill them from the back — the natural
//	direction of addition — writing each column's digit into the next free cell
//	moving left. At the end, trim a single unused leading cell if the top carry
//	was 0. This trades the reverse for one up-front allocation and index math.
//
// Algorithm:
//  1. size = max(m, n) + 1; buf = make([]byte, size); pos = size-1.
//  2. Add columns as in Approach 1, writing buf[pos] = digit, pos--.
//  3. The written region is buf[pos+1:]; that is already most-significant-first.
//
// Time:  O(max(m, n)) — single pass.
// Space: O(max(m, n)) — the output buffer.
func preSizedBuffer(num1 string, num2 string) string {
	m, n := len(num1), len(num2)
	size := maxInt(m, n) + 1  // +1 leaves room for a possible final carry
	buf := make([]byte, size) // result digits, written back-to-front
	pos := size - 1           // next cell to write (rightmost first)

	i, j, carry := m-1, n-1, 0
	for i >= 0 || j >= 0 || carry > 0 {
		sum := carry
		if i >= 0 {
			sum += int(num1[i] - '0')
			i--
		}
		if j >= 0 {
			sum += int(num2[j] - '0')
			j--
		}
		buf[pos] = byte(sum%10) + '0' // place digit directly in its final position
		pos--                         // move one cell to the left
		carry = sum / 10
	}
	// buf[pos+1:] is the filled, correctly-ordered result (leading cell, if the
	// top carry was 0, is simply left out by starting the slice at pos+1).
	return string(buf[pos+1:])
}

// maxInt returns the larger of two ints. (Go 1.21+ ships a builtin max, but a
// named helper keeps this file explicit and avoids shadowing the builtin.)
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	fmt.Println("=== Approach 1: Two Pointers + Carry (Optimal) ===")
	fmt.Println(twoPointers("11", "123")) // expected 134
	fmt.Println(twoPointers("456", "77")) // expected 533
	fmt.Println(twoPointers("0", "0"))    // expected 0
	fmt.Println(twoPointers("99", "1"))   // expected 100 (carry ripples to a new digit)

	fmt.Println("=== Approach 2: Pre-Sized Buffer (No Reverse) ===")
	fmt.Println(preSizedBuffer("11", "123")) // expected 134
	fmt.Println(preSizedBuffer("456", "77")) // expected 533
	fmt.Println(preSizedBuffer("0", "0"))    // expected 0
	fmt.Println(preSizedBuffer("99", "1"))   // expected 100
}
