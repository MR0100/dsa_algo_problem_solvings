package main

import (
	"fmt"
	"strconv"
)

// isPalindrome reports whether the decimal representation of x reads the same
// forwards and backwards. Used by the brute-force approach.
func isPalindrome(x int) bool {
	s := strconv.Itoa(x) // decimal digits as a string
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		if s[i] != s[j] { // mismatched mirror digits → not a palindrome
			return false
		}
	}
	return true
}

// ── Approach 1: Brute Force (Descending Product Search) ──────────────────────
//
// bruteForce solves Largest Palindrome Product by scanning products of two
// n-digit numbers from the largest downward and returning the first palindrome.
//
// Intuition:
//
//	n-digit numbers live in [10^(n-1), 10^n − 1]. To find the LARGEST palindrome
//	product, try big products first: iterate i from hi down, j from i down (so
//	each unordered pair once), and remember the biggest product i*j that is a
//	palindrome. Prune: once i*i ≤ best we can never beat `best`, so stop.
//
// Algorithm:
//  1. hi = 10^n − 1, lo = 10^(n-1).
//  2. best = 0.
//  3. For i = hi..lo: if i*i ≤ best break; for j = i..lo: if i*j ≤ best break;
//     if i*j is a palindrome, best = i*j (then break inner — j only shrinks).
//  4. Return best % 1337.
//
// Time:  O((10^n)²) worst case — fine for n ≤ 4, hopeless for n = 8 (10¹⁶ ops).
// Space: O(n) for the palindrome string check.
func bruteForce(n int) int {
	if n == 1 {
		return 9 // 3*3 = 9 is the largest single-digit-product palindrome
	}
	hi := pow10(n) - 1 // largest n-digit number, e.g. n=2 → 99
	lo := pow10(n - 1) // smallest n-digit number, e.g. n=2 → 10
	best := 0
	for i := hi; i >= lo; i-- {
		if i*i <= best { // even i*i can't exceed best → no larger product remains
			break
		}
		for j := i; j >= lo; j-- {
			prod := i * j
			if prod <= best { // products only shrink as j drops → give up on this i
				break
			}
			if isPalindrome(prod) {
				best = prod // record a new champion; inner products now all smaller
				break
			}
		}
	}
	return best % 1337
}

// ── Approach 2: Construct Palindrome, Test Factorization (Optimal) ────────────
//
// buildAndFactor solves Largest Palindrome Product by generating candidate
// palindromes from largest to smallest and checking whether each factors into
// two n-digit numbers.
//
// Intuition:
//
//	For n > 1 the answer has exactly 2n digits and is a palindrome, so it is
//	fully determined by its first half. Enumerate halves from the largest
//	(999… ) downward; mirror each half to build a full even-length palindrome
//	P. P is the biggest remaining candidate, so the FIRST one that has a divisor
//	d in [lo, hi] with P/d also in [lo, hi] is the answer. To test, trial-divide
//	P by candidates d from hi downward, stopping when d*d < P (past the square
//	root, the cofactor would exceed hi).
//
// Algorithm:
//  1. If n == 1 return 9.
//  2. hi = 10^n − 1, lo = 10^(n-1).
//  3. For half = hi down to lo:
//     - P = half concatenated with reverse(half)  (an even 2n-digit palindrome).
//     - For d = hi down while d*d >= P:
//     if P % d == 0 and P/d is n-digit (≤ hi, and automatically ≥ lo), return
//     P % 1337.
//  4. (Unreachable for valid n.)
//
// Why d*d >= P is the stop: if d < √P then P/d > √P ≥ d, and since we scan d
// from hi down, P/d would already have to be ≤ hi — but P/d > √P and the pair
// is symmetric, so all valid factor pairs are found before crossing √P.
//
// Time:  O(10^n · (hi − √P)) ≈ near-O(10^n) palindromes × short divisor scans;
//
//	the first candidate almost always succeeds, so it is effectively fast for
//	all n ≤ 8. Uses 64-bit products (P up to 10^16).
//
// Space: O(1) beyond a few 64-bit scalars.
func buildAndFactor(n int) int {
	if n == 1 {
		return 9
	}
	hi := pow10(n) - 1 // largest n-digit factor
	lo := pow10(n - 1) // smallest n-digit factor
	for half := hi; half >= lo; half-- {
		p := makePalindrome(half) // mirror half → full 2n-digit palindrome (int64-safe)
		// Trial divide P by n-digit numbers from the top; stop past sqrt(P).
		for d := hi; d*d >= p; d-- {
			if p%d == 0 { // d divides P
				cofactor := p / d
				if cofactor <= hi { // cofactor is n-digit (it is ≥ lo automatically here)
					return int(p % 1337)
				}
			}
		}
	}
	return -1 // not reachable for 1 <= n <= 8
}

// makePalindrome takes an n-digit half and returns the 2n-digit even-length
// palindrome formed by appending the reversed half. Returned as int (64-bit on
// this platform) because for n=8 the palindrome reaches ~10^16.
func makePalindrome(half int) int {
	pal := half
	rev := half
	for rev > 0 {
		pal = pal*10 + rev%10 // shift pal left and append reversed digit
		rev /= 10
	}
	return pal
}

// pow10 returns 10^e for small non-negative e (used for digit-range bounds).
func pow10(e int) int {
	result := 1
	for i := 0; i < e; i++ {
		result *= 10
	}
	return result
}

func main() {
	// The true largest palindromes (before mod 1337) for reference:
	//   n=1 → 9            (mod 1337 = 9)
	//   n=2 → 9009         (99*91)        → 987
	//   n=3 → 906609       (913*993)      → 123
	//   n=4 → 99000099     (9901*9999)    → 597
	//   n=5 → 9966006699                  → 677
	//   n=6 → 999000000999                → 1218
	//   n=7 → 99956644665999              → 877
	//   n=8 → 9999000000009999            → 475

	fmt.Println("=== Approach 1: Brute Force (Descending Product Search) ===")
	fmt.Printf("n=1  got=%d  expected 9\n", bruteForce(1))
	fmt.Printf("n=2  got=%d  expected 987\n", bruteForce(2))
	fmt.Printf("n=3  got=%d  expected 123\n", bruteForce(3))
	fmt.Printf("n=4  got=%d  expected 597\n", bruteForce(4)) // brute force still OK at n=4

	fmt.Println("=== Approach 2: Construct Palindrome, Test Factorization (Optimal) ===")
	fmt.Printf("n=1  got=%d  expected 9\n", buildAndFactor(1))
	fmt.Printf("n=2  got=%d  expected 987\n", buildAndFactor(2))
	fmt.Printf("n=3  got=%d  expected 123\n", buildAndFactor(3))
	fmt.Printf("n=4  got=%d  expected 597\n", buildAndFactor(4))
	fmt.Printf("n=5  got=%d  expected 677\n", buildAndFactor(5))
	fmt.Printf("n=6  got=%d  expected 1218\n", buildAndFactor(6))
	fmt.Printf("n=7  got=%d  expected 877\n", buildAndFactor(7))
	fmt.Printf("n=8  got=%d  expected 475\n", buildAndFactor(8))
}
