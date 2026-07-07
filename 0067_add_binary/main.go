package main

import (
	"fmt"
	"strings"
)

// ── Approach 1: Bit-by-Bit Addition ──────────────────────────────────────────
//
// addBinary solves Add Binary by simulating binary addition from LSB to MSB.
//
// Intuition:
//   Walk both strings from right to left. At each position, sum the two bits
//   (0 or 1) plus the carry. Write sum%2 as the result bit; carry = sum/2.
//   Continue until both strings are exhausted and carry is 0.
//
// Algorithm:
//   i = len(a)-1; j = len(b)-1; carry = 0
//   while i>=0 or j>=0 or carry>0:
//     sum = carry
//     if i>=0: sum += a[i]-'0'; i--
//     if j>=0: sum += b[j]-'0'; j--
//     prepend sum%2; carry = sum/2
//
// Time:  O(max(len(a), len(b)))
// Space: O(max(len(a), len(b))) — result string.
func addBinary(a string, b string) string {
	i, j := len(a)-1, len(b)-1
	carry := 0
	var sb strings.Builder

	for i >= 0 || j >= 0 || carry > 0 {
		sum := carry
		if i >= 0 {
			sum += int(a[i] - '0')
			i--
		}
		if j >= 0 {
			sum += int(b[j] - '0')
			j--
		}
		sb.WriteByte(byte('0' + sum%2)) // write LSB first
		carry = sum / 2
	}

	// reverse the result (we built it LSB-first)
	res := []byte(sb.String())
	for l, r := 0, len(res)-1; l < r; l, r = l+1, r-1 {
		res[l], res[r] = res[r], res[l]
	}
	return string(res)
}

func main() {
	fmt.Println("=== Add Binary ===")
	fmt.Printf("a=%q b=%q  got=%q  expected %q\n", "11", "1", addBinary("11", "1"), "100")
	fmt.Printf("a=%q b=%q  got=%q  expected %q\n", "1010", "1011", addBinary("1010", "1011"), "10101")
	fmt.Printf("a=%q b=%q  got=%q  expected %q\n", "0", "0", addBinary("0", "0"), "0")
	fmt.Printf("a=%q b=%q  got=%q  expected %q\n", "1111", "1111", addBinary("1111", "1111"), "11110")
}
