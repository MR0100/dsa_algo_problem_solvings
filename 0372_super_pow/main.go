package main

import "fmt"

const mod = 1337 // the fixed modulus required by the problem

// ── Approach 1: Digit-by-Digit Horner Expansion ──────────────────────────────
//
// superPowHorner computes a^b mod 1337 where b is given as an array of decimal
// digits. It walks b left to right, using the identity
//
//	a^(10*x + d) = (a^x)^10 * a^d
//
// to fold one digit at a time — the classic way to raise to a giant exponent
// written out digit by digit (Horner's method on the exponent).
//
// Intuition:
//
//	If we have already computed r = a^(prefix), then appending digit d to the
//	exponent means the new exponent is prefix*10 + d, so the new result is
//	r^10 * a^d, all under mod 1337. Starting from r = 1 (empty prefix, a^0) and
//	consuming every digit reconstructs a^b without ever forming the huge b.
//
// Algorithm:
//  1. result = 1.
//  2. Reduce a modulo 1337 once up front.
//  3. For each digit d in b (most significant first):
//     result = powMod(result, 10) * powMod(a, d), taken mod 1337.
//  4. Return result.
//
// Time:  O(n) — n = len(b); each digit costs O(1) modular exponentiations with
//
//	tiny fixed exponents (10 and ≤9).
//
// Space: O(1) — a handful of scalars.
func superPowHorner(a int, b []int) int {
	a %= mod    // shrink the base first so products stay small
	result := 1 // a^0 for the empty prefix
	for _, d := range b {
		// Raise the running result to the 10th power (shift exponent one decimal
		// place) then multiply by a^d for the new least-significant digit.
		result = powMod(result, 10) * powMod(a, d) % mod
	}
	return result
}

// powMod returns base^exp mod 1337 using fast (binary) exponentiation.
//
// Time:  O(log exp). Space: O(1).
func powMod(base, exp int) int {
	base %= mod   // keep base in range
	res := 1      // multiplicative identity
	for exp > 0 { // square-and-multiply over the bits of exp
		if exp&1 == 1 { // this bit is set → include current base power
			res = res * base % mod
		}
		base = base * base % mod // square the base for the next bit
		exp >>= 1                // move to the next higher bit
	}
	return res
}

// ── Approach 2: Recursive Split on Last Digit (Optimal / idiomatic) ───────────
//
// superPowRecursive uses the same math but expressed recursively by peeling the
// LAST digit off b each call: a^b = a^(last) * (a^(b_without_last))^10.
//
// Intuition:
//
//	Let b = [b0, b1, ..., bk]. Then b = 10 * [b0..b_{k-1}] + bk, so
//	a^b = (a^[b0..b_{k-1}])^10 * a^bk. Recurse on the shorter prefix and combine.
//	Base case: empty exponent array → a^0 = 1.
//
// Algorithm:
//  1. If b is empty, return 1.
//  2. last = b[len-1]; prefix = b[:len-1].
//  3. part1 = powMod(a, last).
//  4. part2 = powMod(superPowRecursive(a, prefix), 10).
//  5. Return part1 * part2 mod 1337.
//
// Time:  O(n) recursive calls, each O(log) for the fixed small exponents.
// Space: O(n) recursion stack.
func superPowRecursive(a int, b []int) int {
	if len(b) == 0 { // empty exponent means a^0 = 1
		return 1
	}
	last := b[len(b)-1]                               // least-significant digit
	prefix := b[:len(b)-1]                            // everything above it
	part1 := powMod(a, last)                          // a^(last digit)
	part2 := powMod(superPowRecursive(a, prefix), 10) // (a^prefix)^10
	return part1 * part2 % mod                        // combine under the modulus
}

func main() {
	fmt.Println("=== Approach 1: Digit-by-Digit Horner ===")
	fmt.Println(superPowHorner(2, []int{3}))                // expected 8
	fmt.Println(superPowHorner(2, []int{1, 0}))             // expected 1024
	fmt.Println(superPowHorner(1, []int{4, 3, 3, 8, 5, 2})) // expected 1
	fmt.Println(superPowHorner(2147483647, []int{2, 0, 0})) // expected 1198

	fmt.Println("=== Approach 2: Recursive Split on Last Digit ===")
	fmt.Println(superPowRecursive(2, []int{3}))                // expected 8
	fmt.Println(superPowRecursive(2, []int{1, 0}))             // expected 1024
	fmt.Println(superPowRecursive(1, []int{4, 3, 3, 8, 5, 2})) // expected 1
	fmt.Println(superPowRecursive(2147483647, []int{2, 0, 0})) // expected 1198
}
