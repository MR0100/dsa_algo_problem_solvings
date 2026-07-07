package main

import (
	"fmt"
	"strconv"
	"strings"
)

// ── Approach 1: Brute Force (Long Division + Linear Remainder Scan) ──────────
//
// bruteForce solves Fraction to Recurring Decimal by simulating schoolbook
// long division and detecting a cycle by linearly scanning all remainders
// seen so far.
//
// Intuition:
//
//	A decimal expansion repeats exactly when a remainder repeats during long
//	division: once the same remainder shows up again, the same sequence of
//	digits must follow forever. The simplest cycle detector is to keep every
//	remainder in a slice and, before producing each new digit, scan the slice
//	to see whether the current remainder has appeared before.
//
// Algorithm:
//  1. Handle numerator == 0 → "0".
//  2. Determine the sign from the XOR of the operands' signs.
//  3. Work with int64 absolute values (|-2^31| overflows int32/int).
//  4. Emit the integer part (num / den) and stop if there is no remainder.
//  5. Long division: multiply remainder by 10, emit quotient digit, keep the
//     new remainder. Before each digit, linearly search past remainders; on
//     a match, wrap the digits from that position in parentheses.
//
// Time:  O(L^2) — L = answer length (≤ 10^4); each digit triggers a scan of
//
//	up to L stored remainders.
//
// Space: O(L) — the slice of remainders and the digit buffer.
func bruteForce(numerator int, denominator int) string {
	// Zero numerator short-circuits everything: 0 / anything = "0".
	if numerator == 0 {
		return "0"
	}

	var sb strings.Builder
	// The result is negative iff exactly one operand is negative.
	if (numerator < 0) != (denominator < 0) {
		sb.WriteByte('-')
	}

	// Promote to int64 before negating: -(-2^31) overflows a 32-bit int.
	num, den := int64(numerator), int64(denominator)
	if num < 0 {
		num = -num
	}
	if den < 0 {
		den = -den
	}

	// Integer part of the quotient.
	sb.WriteString(strconv.FormatInt(num/den, 10))

	rem := num % den
	// Division is exact — no fractional part needed.
	if rem == 0 {
		return sb.String()
	}
	sb.WriteByte('.')

	// remainders[i] is the remainder that produced fractional digit i.
	remainders := []int64{}
	// digits accumulates the fractional digits produced so far.
	digits := []byte{}

	for rem != 0 {
		// Linear scan: has this remainder produced a digit before?
		for i, r := range remainders {
			if r == rem {
				// Cycle found: digits[i:] repeat forever → wrap in parens.
				sb.Write(digits[:i])
				sb.WriteByte('(')
				sb.Write(digits[i:])
				sb.WriteByte(')')
				return sb.String()
			}
		}
		// Record the remainder responsible for the digit we emit next.
		remainders = append(remainders, rem)
		rem *= 10                                  // bring down a zero, as in long division
		digits = append(digits, byte('0'+rem/den)) // next fractional digit
		rem %= den                                 // remainder carried into the next step
	}

	// Remainder hit 0 → terminating decimal, no parentheses.
	sb.Write(digits)
	return sb.String()
}

// ── Approach 2: Hash Map of Remainders (Optimal) ─────────────────────────────
//
// hashMap solves Fraction to Recurring Decimal with long division plus a
// hash map from remainder → position of the digit it produced.
//
// Intuition:
//
//	Same cycle-detection insight as Approach 1, but the "have I seen this
//	remainder?" question is answered in O(1) with a map. The map stores the
//	index (inside the string being built) where each remainder first produced
//	a digit, so on a repeat we know exactly where to insert '('.
//
// Algorithm:
//  1. Handle numerator == 0, the sign, and int64 absolute values as before.
//  2. Emit the integer part; return early on exact division.
//  3. For each long-division step: if the current remainder is in the map,
//     splice "(" at the stored index and append ")" — done. Otherwise store
//     the current builder length, emit rem*10/den, keep rem = rem*10%den.
//
// Time:  O(L) — each of the ≤ 10^4 digits is produced with O(1) map work.
// Space: O(L) — the map holds at most one entry per emitted digit.
func hashMap(numerator int, denominator int) string {
	// Zero numerator short-circuits everything: 0 / anything = "0".
	if numerator == 0 {
		return "0"
	}

	var sb strings.Builder
	// Negative result iff signs differ.
	if (numerator < 0) != (denominator < 0) {
		sb.WriteByte('-')
	}

	// int64 avoids overflow when taking absolute values of -2^31.
	num, den := int64(numerator), int64(denominator)
	if num < 0 {
		num = -num
	}
	if den < 0 {
		den = -den
	}

	// Integer part.
	sb.WriteString(strconv.FormatInt(num/den, 10))

	rem := num % den
	if rem == 0 {
		return sb.String() // terminating with no fractional part
	}
	sb.WriteByte('.')

	// seen maps a remainder to the string index where its digit begins.
	seen := map[int64]int{}

	for rem != 0 {
		if pos, ok := seen[rem]; ok {
			// Repeat detected: everything from pos onward is the cycle.
			s := sb.String()
			return s[:pos] + "(" + s[pos:] + ")"
		}
		// The digit generated by this remainder starts at the current length.
		seen[rem] = sb.Len()
		rem *= 10                                      // long-division step
		sb.WriteString(strconv.FormatInt(rem/den, 10)) // emit next digit
		rem %= den                                     // carry the new remainder
	}

	// rem became 0 → the decimal terminates.
	return sb.String()
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Long Division + Linear Remainder Scan) ===")
	fmt.Printf("numerator=1, denominator=2      got=%q  expected \"0.5\"\n", bruteForce(1, 2))
	fmt.Printf("numerator=2, denominator=1      got=%q  expected \"2\"\n", bruteForce(2, 1))
	fmt.Printf("numerator=4, denominator=333    got=%q  expected \"0.(012)\"\n", bruteForce(4, 333))
	fmt.Printf("numerator=1, denominator=6      got=%q  expected \"0.1(6)\"\n", bruteForce(1, 6))                       // repeat starts after a non-repeating prefix
	fmt.Printf("numerator=-50, denominator=8    got=%q  expected \"-6.25\"\n", bruteForce(-50, 8))                      // negative terminating decimal
	fmt.Printf("numerator=-2147483648, denominator=-1  got=%q  expected \"2147483648\"\n", bruteForce(-2147483648, -1)) // int32 overflow edge

	fmt.Println("=== Approach 2: Hash Map of Remainders (Optimal) ===")
	fmt.Printf("numerator=1, denominator=2      got=%q  expected \"0.5\"\n", hashMap(1, 2))
	fmt.Printf("numerator=2, denominator=1      got=%q  expected \"2\"\n", hashMap(2, 1))
	fmt.Printf("numerator=4, denominator=333    got=%q  expected \"0.(012)\"\n", hashMap(4, 333))
	fmt.Printf("numerator=1, denominator=6      got=%q  expected \"0.1(6)\"\n", hashMap(1, 6))
	fmt.Printf("numerator=-50, denominator=8    got=%q  expected \"-6.25\"\n", hashMap(-50, 8))
	fmt.Printf("numerator=-2147483648, denominator=-1  got=%q  expected \"2147483648\"\n", hashMap(-2147483648, -1))
}
