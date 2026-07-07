package main

import "fmt"

// ── Approach 1: Iterative Digit Sum (Brute Force) ────────────────────────────
//
// iterativeDigitSum solves Add Digits by literally repeating the process:
// sum the digits, and keep repeating while the result has more than one digit.
//
// Intuition:
//
//	The problem is defined as a loop: replace num by the sum of its digits
//	until a single digit remains. Just implement that loop directly.
//
// Algorithm:
//  1. While num >= 10 (more than one digit):
//     a. sum = 0; while num > 0: sum += num%10; num /= 10.
//     b. num = sum.
//  2. Return num.
//
// Time:  O(log num) per pass, O(log* num) passes ≈ O(log num) overall — the
//
//	value shrinks extremely fast (each pass roughly to its digit sum).
//
// Space: O(1) — a couple of scalars.
func iterativeDigitSum(num int) int {
	for num >= 10 { // keep going while more than one digit
		sum := 0
		for num > 0 { // add up the decimal digits
			sum += num % 10 // take the lowest digit
			num /= 10       // drop it
		}
		num = sum // replace num with its digit sum, repeat
	}
	return num
}

// ── Approach 2: Single-Pass Recursion ────────────────────────────────────────
//
// recursiveDigitSum solves Add Digits recursively: one pass computes the digit
// sum, then recurse on that sum until it is a single digit.
//
// Intuition:
//
//	Same repeated digit-sum process, expressed as recursion. If num < 10 it is
//	already the digital root; otherwise recurse on its digit sum.
//
// Algorithm:
//  1. If num < 10, return num.
//  2. Compute digit sum s of num.
//  3. Return recursiveDigitSum(s).
//
// Time:  O(log num) amortized — same as the iterative version.
// Space: O(log* num) recursion depth ≈ O(1) in practice.
func recursiveDigitSum(num int) int {
	if num < 10 { // base case: already a single digit
		return num
	}
	sum := 0
	for n := num; n > 0; n /= 10 { // sum the digits once
		sum += n % 10
	}
	return recursiveDigitSum(sum) // recurse on the reduced value
}

// ── Approach 3: Digital Root O(1) Formula (Optimal) ──────────────────────────
//
// digitalRoot solves Add Digits in constant time using the closed-form
// digital-root formula.
//
// Intuition:
//
//	Repeatedly summing digits preserves a number's value mod 9 (because
//	10 ≡ 1 mod 9, so every digit contributes its face value mod 9). The result
//	is the "digital root": 0 for num == 0, and 1 + (num-1) mod 9 otherwise.
//	The (num-1) shift maps multiples of 9 to 9 instead of 0.
//
// Algorithm:
//  1. If num == 0, return 0.
//  2. Return 1 + (num-1) % 9.
//
// Time:  O(1) — pure arithmetic.
// Space: O(1).
func digitalRoot(num int) int {
	if num == 0 {
		return 0 // 0 stays 0 (the formula below would also give 0, guarded for clarity)
	}
	// 1 + (num-1)%9 collapses to 9 for multiples of 9, else num%9.
	return 1 + (num-1)%9
}

func main() {
	fmt.Println("=== Approach 1: Iterative Digit Sum ===")
	fmt.Println(iterativeDigitSum(38)) // expected 2
	fmt.Println(iterativeDigitSum(0))  // expected 0

	fmt.Println("=== Approach 2: Single-Pass Recursion ===")
	fmt.Println(recursiveDigitSum(38)) // expected 2
	fmt.Println(recursiveDigitSum(0))  // expected 0

	fmt.Println("=== Approach 3: Digital Root O(1) Formula (Optimal) ===")
	fmt.Println(digitalRoot(38)) // expected 2
	fmt.Println(digitalRoot(0))  // expected 0
	fmt.Println(digitalRoot(9))  // expected 9 (multiple of 9 maps to 9, not 0)
}
