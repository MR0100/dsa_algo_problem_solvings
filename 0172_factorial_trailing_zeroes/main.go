package main

import (
	"fmt"
	"math/big"
)

// ── Approach 1: Brute Force (Compute n! with Big Integers) ───────────────────
//
// bruteForce solves Factorial Trailing Zeroes by literally computing n! with
// arbitrary-precision integers and counting the zeros at the end.
//
// Intuition:
//
//	The most direct reading of the problem: build the factorial, look at its
//	decimal representation, count trailing '0' characters. n! overflows int64
//	already at n = 21 (21! ≈ 5.1×10¹⁹), so math/big is mandatory — 10000!
//	has 35,660 digits. Correct, but wildly more work than needed.
//
// Algorithm:
//  1. f = 1 (big.Int); multiply f by every i in 2..n.
//  2. Render f as a decimal string.
//  3. Count '0' characters from the right until a non-zero digit appears.
//
// Time:  O(n² log n) bit-work — n big-integer multiplications, each linear in
//
//	the current digit count (which grows to Θ(n log n) digits).
//
// Space: O(n log n) — the digits of n! itself.
func bruteForce(n int) int {
	f := big.NewInt(1) // running factorial, arbitrary precision
	for i := 2; i <= n; i++ {
		f.Mul(f, big.NewInt(int64(i))) // f *= i, never overflows
	}
	s := f.String() // full decimal expansion of n!
	zeros := 0
	// Walk from the last digit backwards while we keep seeing '0'.
	for i := len(s) - 1; i >= 0 && s[i] == '0'; i-- {
		zeros++
	}
	return zeros
}

// ── Approach 2: Count Factors of 5 per Multiple ──────────────────────────────
//
// countFactorsOfFive solves Factorial Trailing Zeroes by counting how many
// times 5 divides each multiple of 5 up to n — without ever building n!.
//
// Intuition:
//
//	A trailing zero is a factor of 10 = 2×5. In n! = 1·2·…·n the factors of
//	2 vastly outnumber the factors of 5 (every 2nd number is even, only every
//	5th contributes a 5), so #zeros = #factors of 5 in n!. Only multiples of
//	5 contribute any, and numbers like 25 = 5² or 125 = 5³ contribute more
//	than one — hence the inner division loop per multiple.
//
// Algorithm:
//  1. zeros = 0.
//  2. For i = 5, 10, 15, ..., n:
//  3. While i is divisible by 5: zeros++, divide the temp copy by 5.
//  4. Return zeros.
//
// Time:  O(n) — n/5 multiples visited; the inner loop totals n/25 + n/125 + …
//
//	extra steps, still linear overall.
//
// Space: O(1) — two counters.
func countFactorsOfFive(n int) int {
	zeros := 0
	for i := 5; i <= n; i += 5 { // only multiples of 5 carry any factor of 5
		for x := i; x%5 == 0; x /= 5 {
			zeros++ // one zero per factor of 5 inside this multiple (25→2, 125→3, …)
		}
	}
	return zeros
}

// ── Approach 3: Logarithmic Division — Legendre's Formula (Optimal) ──────────
//
// logarithmicDivision solves Factorial Trailing Zeroes with Legendre's
// formula: the exponent of 5 in n! is ⌊n/5⌋ + ⌊n/25⌋ + ⌊n/125⌋ + …
//
// Intuition:
//
//	Instead of asking "how many 5s does each number contribute?", flip the
//	count: ⌊n/5⌋ numbers contribute at least one 5, ⌊n/25⌋ of those
//	contribute a second 5, ⌊n/125⌋ a third, and so on. Each term is one
//	integer division, and the terms shrink by 5× — so the whole sum takes
//	log₅(n) steps. This answers the follow-up exactly.
//
// Algorithm:
//  1. zeros = 0.
//  2. While n > 0: n /= 5; zeros += n.
//     (After the k-th division n holds ⌊original/5^k⌋, so we sum the series.)
//  3. Return zeros.
//
// Time:  O(log₅ n) — one division per power of 5 ≤ n (≤ 6 steps for n ≤ 10⁴).
// Space: O(1) — a single counter.
func logarithmicDivision(n int) int {
	zeros := 0
	for n > 0 {
		n /= 5     // n is now ⌊n/5⌋, ⌊n/25⌋, ⌊n/125⌋, ... on successive turns
		zeros += n // add how many numbers contribute yet another factor of 5
	}
	return zeros
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Compute n! with Big Integers) ===")
	fmt.Printf("n=3      got=%-5d expected 0\n", bruteForce(3))        // 3! = 6, no trailing zero
	fmt.Printf("n=5      got=%-5d expected 1\n", bruteForce(5))        // 5! = 120, one trailing zero
	fmt.Printf("n=0      got=%-5d expected 0\n", bruteForce(0))        // 0! = 1, no trailing zero
	fmt.Printf("n=10     got=%-5d expected 2\n", bruteForce(10))       // 10! = 3628800
	fmt.Printf("n=25     got=%-5d expected 6\n", bruteForce(25))       // 25 adds TWO fives (5·5)
	fmt.Printf("n=10000  got=%-5d expected 2499\n", bruteForce(10000)) // constraint upper bound

	fmt.Println("=== Approach 2: Count Factors of 5 per Multiple ===")
	fmt.Printf("n=3      got=%-5d expected 0\n", countFactorsOfFive(3))
	fmt.Printf("n=5      got=%-5d expected 1\n", countFactorsOfFive(5))
	fmt.Printf("n=0      got=%-5d expected 0\n", countFactorsOfFive(0))
	fmt.Printf("n=10     got=%-5d expected 2\n", countFactorsOfFive(10))
	fmt.Printf("n=25     got=%-5d expected 6\n", countFactorsOfFive(25))
	fmt.Printf("n=10000  got=%-5d expected 2499\n", countFactorsOfFive(10000))

	fmt.Println("=== Approach 3: Logarithmic Division — Legendre's Formula (Optimal) ===")
	fmt.Printf("n=3      got=%-5d expected 0\n", logarithmicDivision(3))
	fmt.Printf("n=5      got=%-5d expected 1\n", logarithmicDivision(5))
	fmt.Printf("n=0      got=%-5d expected 0\n", logarithmicDivision(0))
	fmt.Printf("n=10     got=%-5d expected 2\n", logarithmicDivision(10))
	fmt.Printf("n=25     got=%-5d expected 6\n", logarithmicDivision(25))
	fmt.Printf("n=10000  got=%-5d expected 2499\n", logarithmicDivision(10000))
}
