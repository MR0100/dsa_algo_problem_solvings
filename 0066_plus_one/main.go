package main

import "fmt"

// ── Approach 1: Scan from Right ───────────────────────────────────────────────
//
// plusOne solves Plus One by simulating addition from the least significant digit.
//
// Intuition:
//   Walk from right to left. If the current digit is < 9, increment and return.
//   If it's 9, set to 0 and carry over (continue left). If we exit the loop
//   (all digits were 9), prepend a 1.
//
// Algorithm:
//   for i = len-1 downto 0:
//     if digits[i] < 9: digits[i]++; return digits
//     digits[i] = 0  // carry
//   return prepend 1
//
// Time:  O(n) worst case (e.g., [9,9,9] → [1,0,0,0]).
//         O(1) amortised (most numbers have a non-9 last digit).
// Space: O(n) worst case (new slice when all 9s); O(1) otherwise.
func plusOne(digits []int) []int {
	for i := len(digits) - 1; i >= 0; i-- {
		if digits[i] < 9 {
			digits[i]++
			return digits
		}
		digits[i] = 0 // carry: 9+1=10, write 0
	}
	// all digits were 9 (e.g., [9,9] → [1,0,0])
	return append([]int{1}, digits...)
}

func main() {
	fmt.Println("=== Plus One ===")
	fmt.Printf("digits=[1,2,3]  got=%v  expected [1 2 4]\n", plusOne([]int{1, 2, 3}))
	fmt.Printf("digits=[4,3,2,1]  got=%v  expected [4 3 2 2]\n", plusOne([]int{4, 3, 2, 1}))
	fmt.Printf("digits=[9]  got=%v  expected [1 0]\n", plusOne([]int{9}))
	fmt.Printf("digits=[9,9,9]  got=%v  expected [1 0 0 0]\n", plusOne([]int{9, 9, 9}))
	fmt.Printf("digits=[1,9,9]  got=%v  expected [2 0 0]\n", plusOne([]int{1, 9, 9}))
}
