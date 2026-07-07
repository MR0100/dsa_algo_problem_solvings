package main

import (
	"fmt"
	"strings"
)

// ── Approach 1: Grade-School Multiplication ───────────────────────────────────
//
// multiply solves Multiply Strings by simulating grade-school long multiplication.
//
// Intuition: Multiplying two n-digit and m-digit numbers gives at most n+m digits.
// Allocate a result array of size n+m. For each digit pair (i from num1, j from num2),
// the product contributes to positions i+j and i+j+1 in the result array.
//
// Algorithm:
//  1. Allocate pos[n+m] zeroed.
//  2. For i = n-1 downto 0, j = m-1 downto 0:
//     mul = (num1[i]-'0') * (num2[j]-'0')
//     p1, p2 = i+j, i+j+1
//     sum = mul + pos[p2]
//     pos[p2] = sum % 10
//     pos[p1] += sum / 10
//  3. Convert pos to string, skip leading zeros.
//
// Time:  O(n * m)
// Space: O(n + m) — the result array
func multiply(num1, num2 string) string {
	if num1 == "0" || num2 == "0" {
		return "0"
	}
	n, m := len(num1), len(num2)
	pos := make([]int, n+m) // pos[i] holds the digit at that position

	for i := n - 1; i >= 0; i-- {
		for j := m - 1; j >= 0; j-- {
			mul := int(num1[i]-'0') * int(num2[j]-'0')
			p1, p2 := i+j, i+j+1   // p2 = units digit of this product; p1 = carry
			sum := mul + pos[p2]
			pos[p2] = sum % 10     // units digit
			pos[p1] += sum / 10    // carry propagates to p1
		}
	}

	var sb strings.Builder
	for _, d := range pos {
		if sb.Len() == 0 && d == 0 {
			continue // skip leading zeros
		}
		sb.WriteByte(byte('0' + d))
	}
	if sb.Len() == 0 {
		return "0"
	}
	return sb.String()
}

func main() {
	fmt.Println("=== Approach 1: Grade-School Multiplication ===")
	fmt.Printf("num1=2  num2=3    => %s  expected 6\n", multiply("2", "3"))
	fmt.Printf("num1=123 num2=456 => %s  expected 56088\n", multiply("123", "456"))
	fmt.Printf("num1=0  num2=0    => %s  expected 0\n", multiply("0", "0"))
	fmt.Printf("num1=99  num2=99  => %s  expected 9801\n", multiply("99", "99"))
	fmt.Printf("num1=9133 num2=0  => %s  expected 0\n", multiply("9133", "0"))
}
