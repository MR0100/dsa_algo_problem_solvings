package main

import "fmt"

// ── Approach 1: Brute Force (Recompute Power per Character) ──────────────────
//
// bruteForce solves Excel Sheet Column Number by treating every letter as a
// digit whose place value is recomputed from scratch with an inner loop.
//
// Intuition:
//
//	"AB" means A·26¹ + B·26⁰ with digit values A=1..Z=26 — exactly like the
//	decimal number 28 means 2·10¹ + 8·10⁰. The most literal translation:
//	for each character, figure out how many positions sit to its right and
//	multiply its digit value by 26 raised to that count, recomputing the
//	power with a fresh inner loop every time.
//
// Algorithm:
//  1. For every index i in the title (left to right):
//  2. digit = title[i] - 'A' + 1  (map 'A'..'Z' → 1..26).
//  3. Compute power = 26^(len-1-i) with an inner multiplication loop.
//  4. Add digit * power to the running total.
//
// Time:  O(k²) — k characters, each recomputing a power in up to k-1 steps
//
//	(k = len(columnTitle) ≤ 7, so tiny in practice).
//
// Space: O(1) — only scalar accumulators.
func bruteForce(columnTitle string) int {
	total := 0
	k := len(columnTitle)
	for i := 0; i < k; i++ {
		digit := int(columnTitle[i]-'A') + 1 // 'A'→1 ... 'Z'→26 (no zero digit!)
		// Recompute 26^(k-1-i) naively — the "brute" part of this approach.
		power := 1
		for p := 0; p < k-1-i; p++ {
			power *= 26 // one factor of 26 per position right of index i
		}
		total += digit * power // place value contribution of this letter
	}
	return total
}

// ── Approach 2: Right-to-Left with Running Power ──────────────────────────────
//
// rightToLeftPower solves Excel Sheet Column Number by scanning from the last
// letter and carrying the place value along, multiplying it by 26 each step.
//
// Intuition:
//
//	Approach 1 wastes work: 26^(i+1) is just 26^i × 26. Walk the string from
//	the rightmost (least-significant) letter, keep the current place value in
//	a variable, and grow it by ×26 as we move one position left. Same math,
//	one pass, no inner loop.
//
// Algorithm:
//  1. power = 1, total = 0.
//  2. For i from len-1 down to 0:
//  3. total += (title[i]-'A'+1) * power.
//  4. power *= 26.
//
// Time:  O(k) — one constant-time step per character.
// Space: O(1) — two scalar accumulators.
func rightToLeftPower(columnTitle string) int {
	total := 0
	power := 1 // place value of the current position: 26^0, 26^1, ...
	for i := len(columnTitle) - 1; i >= 0; i-- {
		digit := int(columnTitle[i]-'A') + 1 // bijective digit 1..26
		total += digit * power               // add this letter's contribution
		power *= 26                          // next position left is 26× more significant
	}
	return total
}

// ── Approach 3: Left-to-Right Horner's Method (Optimal) ──────────────────────
//
// hornersMethod solves Excel Sheet Column Number with a single left-to-right
// accumulation: result = result*26 + digit, the same way you'd read "28"
// aloud digit by digit.
//
// Intuition:
//
//	Horner's rule: A·26¹ + B·26⁰ = (A)·26 + B. Every time one more letter
//	appears on the right, everything read so far becomes 26× more
//	significant. So keep one accumulator, multiply it by 26, add the new
//	digit — no powers tracked at all. This is the canonical string→number
//	pattern (identical to parsing a decimal int, just base 26 and digits
//	starting at 1).
//
// Algorithm:
//  1. result = 0.
//  2. For each character c left to right: result = result*26 + (c-'A'+1).
//  3. Return result.
//
// Time:  O(k) — one multiply-add per character.
// Space: O(1) — a single accumulator.
func hornersMethod(columnTitle string) int {
	result := 0
	for i := 0; i < len(columnTitle); i++ {
		// Shift everything seen so far one position left (×26),
		// then drop the new least-significant digit (1..26) in.
		result = result*26 + int(columnTitle[i]-'A') + 1
	}
	return result
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Recompute Power per Character) ===")
	fmt.Printf("columnTitle=%-11q got=%-10d expected 1\n", "A", bruteForce("A"))
	fmt.Printf("columnTitle=%-11q got=%-10d expected 28\n", "AB", bruteForce("AB"))
	fmt.Printf("columnTitle=%-11q got=%-10d expected 701\n", "ZY", bruteForce("ZY"))
	fmt.Printf("columnTitle=%-11q got=%-10d expected 26\n", "Z", bruteForce("Z"))                     // last 1-letter title
	fmt.Printf("columnTitle=%-11q got=%-10d expected 27\n", "AA", bruteForce("AA"))                   // first 2-letter title
	fmt.Printf("columnTitle=%-11q got=%-10d expected 2147483647\n", "FXSHRXW", bruteForce("FXSHRXW")) // max int32 (constraint upper bound)

	fmt.Println("=== Approach 2: Right-to-Left with Running Power ===")
	fmt.Printf("columnTitle=%-11q got=%-10d expected 1\n", "A", rightToLeftPower("A"))
	fmt.Printf("columnTitle=%-11q got=%-10d expected 28\n", "AB", rightToLeftPower("AB"))
	fmt.Printf("columnTitle=%-11q got=%-10d expected 701\n", "ZY", rightToLeftPower("ZY"))
	fmt.Printf("columnTitle=%-11q got=%-10d expected 26\n", "Z", rightToLeftPower("Z"))
	fmt.Printf("columnTitle=%-11q got=%-10d expected 27\n", "AA", rightToLeftPower("AA"))
	fmt.Printf("columnTitle=%-11q got=%-10d expected 2147483647\n", "FXSHRXW", rightToLeftPower("FXSHRXW"))

	fmt.Println("=== Approach 3: Left-to-Right Horner's Method (Optimal) ===")
	fmt.Printf("columnTitle=%-11q got=%-10d expected 1\n", "A", hornersMethod("A"))
	fmt.Printf("columnTitle=%-11q got=%-10d expected 28\n", "AB", hornersMethod("AB"))
	fmt.Printf("columnTitle=%-11q got=%-10d expected 701\n", "ZY", hornersMethod("ZY"))
	fmt.Printf("columnTitle=%-11q got=%-10d expected 26\n", "Z", hornersMethod("Z"))
	fmt.Printf("columnTitle=%-11q got=%-10d expected 27\n", "AA", hornersMethod("AA"))
	fmt.Printf("columnTitle=%-11q got=%-10d expected 2147483647\n", "FXSHRXW", hornersMethod("FXSHRXW"))
}
