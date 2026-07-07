package main

import "fmt"

// LeetCode 263 — Ugly Number.
//
// An UGLY NUMBER is a positive integer whose prime factors are limited to 2, 3
// and 5. Given an integer n, return true if n is ugly, else false. Note: 1 is
// ugly (it has no prime factors), and every n <= 0 is NOT ugly.

// ── Approach 1: Trial Division (Brute Force / Optimal) ───────────────────────
//
// trialDivision solves Ugly Number by repeatedly dividing out the factors 2, 3
// and 5; n is ugly iff nothing but 1 remains.
//
// Intuition:
//
//	Strip every factor of 2, then every factor of 3, then every factor of 5.
//	If the number was built only from those primes, the leftover is 1. If any
//	other prime factor exists (e.g. 7), it survives the divisions and the
//	leftover is > 1.
//
// Algorithm:
//  1. If n <= 0, return false (ugly numbers are positive).
//  2. For each p in {2, 3, 5}: while n % p == 0, divide n by p.
//  3. Return true iff the remaining n == 1.
//
// Time:  O(log n) — each division at least halves n (factor 2 dominates).
// Space: O(1) — only the running value n.
func trialDivision(n int) bool {
	if n <= 0 { // 0 and negatives are never ugly
		return false
	}
	for _, p := range []int{2, 3, 5} { // divide out each allowed prime fully
		for n%p == 0 { // while p divides n evenly
			n /= p // remove one factor of p
		}
	}
	return n == 1 // only 2/3/5 factors ⇒ nothing but 1 is left
}

// ── Approach 2: Recursion ────────────────────────────────────────────────────
//
// recursive solves Ugly Number by peeling one allowed prime factor per call
// and recursing on the quotient.
//
// Intuition:
//
//	n is ugly iff it is divisible by one of 2/3/5 and the quotient is also
//	ugly, with 1 as the base case. This is the same divide-out logic expressed
//	as recursion instead of loops.
//
// Algorithm:
//  1. Base cases: n <= 0 ⇒ false; n == 1 ⇒ true.
//  2. If n divisible by 2/3/5, recurse on n / that prime.
//  3. Otherwise n has a forbidden prime factor ⇒ false.
//
// Time:  O(log n) — one factor removed per recursive call.
// Space: O(log n) — recursion stack depth equals the number of prime factors.
func recursive(n int) bool {
	if n <= 0 { // non-positive ⇒ not ugly
		return false
	}
	if n == 1 { // fully reduced ⇒ ugly
		return true
	}
	for _, p := range []int{2, 3, 5} {
		if n%p == 0 { // divisible by an allowed prime
			return recursive(n / p) // recurse on the quotient
		}
	}
	return false // divisible by none of 2/3/5 ⇒ has a forbidden factor
}

func main() {
	// Example 1: n = 6 = 2·3 ⇒ true.
	// Example 2: n = 1 ⇒ true (no prime factors).
	// Example 3: n = 14 = 2·7 ⇒ false (factor 7 is forbidden).

	fmt.Println("=== Approach 1: Trial Division (Optimal) ===")
	fmt.Println(trialDivision(6))  // expected true
	fmt.Println(trialDivision(1))  // expected true
	fmt.Println(trialDivision(14)) // expected false

	fmt.Println("=== Approach 2: Recursion ===")
	fmt.Println(recursive(6))  // expected true
	fmt.Println(recursive(1))  // expected true
	fmt.Println(recursive(14)) // expected false
}
