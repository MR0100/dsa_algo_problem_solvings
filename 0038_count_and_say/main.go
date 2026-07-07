package main

import (
	"fmt"
	"strings"
)

// ── Approach 1: Iterative Simulation ─────────────────────────────────────────
//
// iterative solves Count and Say by iteratively building each sequence from
// the previous one using run-length encoding.
//
// Intuition: Starting from "1", each subsequent term is the run-length encoding
// of the previous term. Walk the current string, count runs of identical digits,
// and emit "count + digit" for each run.
//
// Algorithm:
//  1. result = "1"
//  2. For i = 2 to n:
//     build next = rle(result).
//     result = next.
//  3. Return result.
//
// Run-Length Encoding helper:
//  - Walk the string; count consecutive identical chars; emit count+char.
//
// Time:  O(n * L) where L = average length of the sequence; L grows but is
//        bounded by O(n * 2^n/3) asymptotically (Conway's constant ≈ 1.30)
// Space: O(L) — the current sequence string
func iterative(n int) string {
	result := "1"
	for i := 2; i <= n; i++ {
		result = rle(result)
	}
	return result
}

// rle returns the run-length encoding of s.
func rle(s string) string {
	var sb strings.Builder
	j := 0
	for j < len(s) {
		ch := s[j]
		count := 1
		for j+count < len(s) && s[j+count] == ch { // count the run
			count++
		}
		sb.WriteByte(byte('0' + count)) // write count
		sb.WriteByte(ch)                // write digit
		j += count
	}
	return sb.String()
}

// ── Approach 2: Recursive ────────────────────────────────────────────────────
//
// recursive solves Count and Say using top-down recursion.
//
// Intuition: countAndSay(n) = rle(countAndSay(n-1)), with base case n==1 → "1".
//
// Time:  O(n * L) — same as iterative; n recursive calls each doing O(L) work.
// Space: O(n * L) — recursion stack holds n strings simultaneously.
func recursive(n int) string {
	if n == 1 {
		return "1"
	}
	return rle(recursive(n - 1))
}

func main() {
	fmt.Println("=== Approach 1: Iterative ===")
	for n := 1; n <= 6; n++ {
		fmt.Printf("n=%d  %s\n", n, iterative(n))
	}
	// LeetCode examples
	fmt.Printf("\nn=1  got=%s  expected=1\n", iterative(1))
	fmt.Printf("n=4  got=%s  expected=1211\n", iterative(4))

	fmt.Println("\n=== Approach 2: Recursive ===")
	fmt.Printf("n=1  got=%s  expected=1\n", recursive(1))
	fmt.Printf("n=4  got=%s  expected=1211\n", recursive(4))
	fmt.Printf("n=5  got=%s  expected=111221\n", recursive(5))
}
