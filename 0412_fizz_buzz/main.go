package main

import (
	"fmt"
	"strconv"
	"strings"
)

// ── Approach 1: Brute Force (Modulo Checks) ──────────────────────────────────
//
// bruteForce solves Fizz Buzz by testing divisibility with the modulo operator
// for each integer from 1 to n.
//
// Intuition:
//
//	The problem is a direct rule table. Walk every i in [1, n] and pick the label
//	by asking, in the right ORDER: is i divisible by both 3 and 5 (equivalently
//	by 15)? then by 3? then by 5? otherwise use the number itself. Order matters:
//	the "FizzBuzz" case must be tested first, otherwise a multiple of 15 would be
//	caught by the "Fizz" (or "Buzz") branch and never reach "FizzBuzz".
//
// Algorithm:
//  1. Allocate answer of length n.
//  2. For i = 1..n:
//     a. i%15 == 0 → "FizzBuzz"
//     b. else i%3 == 0 → "Fizz"
//     c. else i%5 == 0 → "Buzz"
//     d. else strconv.Itoa(i)
//  3. Return answer.
//
// Time:  O(n) — one constant-work pass over 1..n.
// Space: O(1) extra (ignoring the required output slice of n strings).
func bruteForce(n int) []string {
	answer := make([]string, n) // 0-indexed slice; answer[i-1] holds the label for value i
	for i := 1; i <= n; i++ {   // problem is 1-indexed, so iterate values 1..n
		switch {
		case i%15 == 0: // divisible by 3 AND 5 → must be checked FIRST
			answer[i-1] = "FizzBuzz"
		case i%3 == 0: // divisible by 3 only
			answer[i-1] = "Fizz"
		case i%5 == 0: // divisible by 5 only
			answer[i-1] = "Buzz"
		default: // divisible by neither → the number as text
			answer[i-1] = strconv.Itoa(i)
		}
	}
	return answer
}

// ── Approach 2: String Concatenation (No Modulo Combination) ──────────────────
//
// stringConcat solves Fizz Buzz by building each label from two independent
// divisibility tests and falling back to the number only when nothing matched.
//
// Intuition:
//
//	Instead of a special "divisible by 15" branch, append "Fizz" when 3 divides i
//	and "Buzz" when 5 divides i. A multiple of 15 satisfies both, so it naturally
//	becomes "Fizz"+"Buzz" = "FizzBuzz" with no separate case. If neither test
//	fired, the accumulated string is still empty, so we substitute the number.
//	This generalises cleanly to extra rules (e.g. add "Bazz" for 7) without an
//	exploding set of combined-modulo branches.
//
// Algorithm:
//  1. For i = 1..n: start with an empty string s.
//  2. If i%3 == 0, s += "Fizz". If i%5 == 0, s += "Buzz".
//  3. If s is still empty, s = strconv.Itoa(i).
//  4. Store s.
//
// Time:  O(n) — constant work per i (short, bounded string building).
// Space: O(1) extra beyond the output slice.
func stringConcat(n int) []string {
	answer := make([]string, n)
	for i := 1; i <= n; i++ {
		var sb strings.Builder // accumulate the label without a combined 15-check
		if i%3 == 0 {
			sb.WriteString("Fizz") // 3 contributes "Fizz"
		}
		if i%5 == 0 {
			sb.WriteString("Buzz") // 5 contributes "Buzz"; both → "FizzBuzz"
		}
		if sb.Len() == 0 { // no divisor matched → use the number itself
			sb.WriteString(strconv.Itoa(i))
		}
		answer[i-1] = sb.String()
	}
	return answer
}

// ── Approach 3: Counter Increments (No Modulo At All) (Optimal) ───────────────
//
// counterNoModulo solves Fizz Buzz while avoiding the modulo/division operation
// entirely — useful when % is expensive or disallowed — by keeping two running
// counters that reset when they reach 3 and 5.
//
// Intuition:
//
//	Divisibility by 3 recurs every 3 steps; by 5 every 5 steps. Track two small
//	counters fizz and buzz that increment each iteration. When fizz hits 3 we are
//	on a multiple of 3 (reset fizz to 0); when buzz hits 5 we are on a multiple of
//	5 (reset buzz to 0). Combine the two boolean "just reset?" flags exactly like
//	Approach 2. No modulo, no division — only additions and comparisons.
//
// Algorithm:
//  1. fizz = 0, buzz = 0.
//  2. For i = 1..n: fizz++, buzz++.
//  3. If fizz == 3 → multiple of 3 (reset fizz). If buzz == 5 → multiple of 5 (reset buzz).
//  4. Emit "Fizz"/"Buzz"/"FizzBuzz"/number from those two flags.
//
// Time:  O(n) — constant work per i, only additions and comparisons.
// Space: O(1) extra beyond the output slice.
func counterNoModulo(n int) []string {
	answer := make([]string, n)
	fizz, buzz := 0, 0 // steps since the last multiple of 3 / of 5
	for i := 1; i <= n; i++ {
		fizz++ // advance both cyclic counters
		buzz++
		isFizz := fizz == 3 // reached a multiple of 3 this step?
		isBuzz := buzz == 5 // reached a multiple of 5 this step?
		if isFizz {
			fizz = 0 // restart the 3-cycle
		}
		if isBuzz {
			buzz = 0 // restart the 5-cycle
		}
		switch {
		case isFizz && isBuzz: // both cycles landed → multiple of 15
			answer[i-1] = "FizzBuzz"
		case isFizz:
			answer[i-1] = "Fizz"
		case isBuzz:
			answer[i-1] = "Buzz"
		default:
			answer[i-1] = strconv.Itoa(i)
		}
	}
	return answer
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Modulo Checks) ===")
	fmt.Println(bruteForce(3))  // expected [1 2 Fizz]
	fmt.Println(bruteForce(5))  // expected [1 2 Fizz 4 Buzz]
	fmt.Println(bruteForce(15)) // expected [1 2 Fizz 4 Buzz Fizz 7 8 Fizz Buzz 11 Fizz 13 14 FizzBuzz]

	fmt.Println("=== Approach 2: String Concatenation ===")
	fmt.Println(stringConcat(3))  // expected [1 2 Fizz]
	fmt.Println(stringConcat(5))  // expected [1 2 Fizz 4 Buzz]
	fmt.Println(stringConcat(15)) // expected [1 2 Fizz 4 Buzz Fizz 7 8 Fizz Buzz 11 Fizz 13 14 FizzBuzz]

	fmt.Println("=== Approach 3: Counter Increments (No Modulo) (Optimal) ===")
	fmt.Println(counterNoModulo(3))  // expected [1 2 Fizz]
	fmt.Println(counterNoModulo(5))  // expected [1 2 Fizz 4 Buzz]
	fmt.Println(counterNoModulo(15)) // expected [1 2 Fizz 4 Buzz Fizz 7 8 Fizz Buzz 11 Fizz 13 14 FizzBuzz]
}
