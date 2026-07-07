package main

import "fmt"

// ── Approach 1: Brute Force (Count Per Number) ───────────────────────────────
//
// bruteForce solves Number of Digit One by counting the digit 1 in every
// integer from 1 to n.
//
// Intuition:
//
//	The definition is literal: sum, over all integers i in [1, n], the number
//	of times digit 1 appears in i. We can just do exactly that — for each i,
//	strip its decimal digits and tally the ones. Simple and obviously correct,
//	useful as a reference oracle, but O(n log n) so only for small n.
//
// Algorithm:
//
//  1. For each i from 1 to n:
//     a. While i has digits left, look at the lowest digit.
//     b. If it is 1, increment the total.
//     c. Drop the lowest digit (i /= 10).
//  2. Return the accumulated total.
//
// Time:  O(n log n) — n numbers, each with O(log n) digits.
// Space: O(1).
func bruteForce(n int) int {
	total := 0
	for i := 1; i <= n; i++ { // consider every integer in [1, n]
		x := i
		for x > 0 { // examine each decimal digit of x
			if x%10 == 1 { // lowest digit is a 1 → count it
				total++
			}
			x /= 10 // discard the lowest digit
		}
	}
	return total
}

// ── Approach 2: Digit-Position Counting (Optimal) ────────────────────────────
//
// digitPosition solves Number of Digit One by counting, for each decimal
// place, how many numbers in [1, n] carry a 1 in that place.
//
// Intuition:
//
//	Instead of counting per number, count per digit position. Fix a place
//	value `p` (1, 10, 100, ...). Split n around that place into high, cur, and
//	low parts. The count of 1s contributed at place p depends only on cur:
//	  - cur == 0 : high * p                        (only complete high cycles)
//	  - cur == 1 : high * p + (low + 1)            (partial cycle up to low)
//	  - cur >= 2 : (high + 1) * p                   (one extra full cycle)
//	Summing over all places gives the answer in O(log n).
//
// Algorithm:
//
//  1. For place p = 1, 10, 100, ... while p <= n:
//     high = n / (p*10); cur = (n / p) % 10; low = n % p.
//     if cur == 0: add high * p
//     if cur == 1: add high * p + low + 1
//     if cur >= 2: add (high + 1) * p
//  2. Return the running sum.
//
// Time:  O(log n) — one iteration per decimal place of n.
// Space: O(1).
func digitPosition(n int) int {
	count := 0
	for p := 1; p <= n; p *= 10 { // iterate over each place value 1,10,100,...
		high := n / (p * 10) // digits above the current place
		cur := (n / p) % 10  // the digit sitting at the current place
		low := n % p         // digits below the current place

		switch {
		case cur == 0:
			// current place never reaches 1 within the last partial cycle:
			// only the `high` completed cycles each contribute p ones.
			count += high * p
		case cur == 1:
			// completed cycles give high*p; the in-progress cycle contributes
			// ones for low+1 numbers (…10000 through …1<low>).
			count += high*p + low + 1
		default: // cur >= 2
			// the current place has fully passed 1 once more, so (high+1)
			// complete cycles each contribute p ones.
			count += (high + 1) * p
		}
	}
	return count
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Count Per Number) ===")
	fmt.Println(bruteForce(13)) // expected 6
	fmt.Println(bruteForce(0))  // expected 0

	fmt.Println("=== Approach 2: Digit-Position Counting (Optimal) ===")
	fmt.Println(digitPosition(13)) // expected 6
	fmt.Println(digitPosition(0))  // expected 0
}
