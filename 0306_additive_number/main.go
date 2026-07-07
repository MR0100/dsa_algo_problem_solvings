package main

import (
	"fmt"
	"math/big"
)

// ── Approach 1: Brute Force Backtracking (string arithmetic) ──────────────────
//
// bruteForce solves Additive Number by trying every possible split of the first
// two numbers and verifying the rest of the string forms a valid additive
// sequence, using string-based addition to dodge integer overflow.
//
// Intuition:
//
//	An additive sequence is fully determined by its first two numbers. Once we
//	fix num1 and num2, the entire rest of the string is forced: num3 = num1 +
//	num2 must appear next, then num2 + num3, and so on. So we only need to
//	enumerate every (num1, num2) prefix pair and check whether the remainder
//	matches the forced continuation. Because the numbers can exceed 64 bits, we
//	add them as decimal strings rather than as ints.
//
// Algorithm:
//
//  1. For every end index i of the first number (num[0:i]) and every end index
//     j of the second number (num[i:j]):
//  2. Reject any candidate with a leading zero (length > 1 and first char '0').
//  3. Call a checker that, starting from num1 and num2 at position j, keeps
//     computing the next expected sum and matching it against the string.
//  4. Return true as soon as any split validates the whole string.
//
// Time:  O(n^2 · n) = O(n^3) — O(n^2) prefix pairs, each check scans O(n) with
//
//	O(n) string additions, but additions are bounded by remaining length.
//
// Space: O(n) — recursion/temporary strings for the addition buffers.
func bruteForce(num string) bool {
	n := len(num)
	// i is the end (exclusive) of the first number num[0:i].
	for i := 1; i <= n-2; i++ {
		if num[0] == '0' && i > 1 {
			break // first number has a leading zero → no longer split can help
		}
		// j is the end (exclusive) of the second number num[i:j].
		for j := i + 1; j <= n-1; j++ {
			if num[i] == '0' && j-i > 1 {
				break // second number has a leading zero → stop extending it
			}
			// Validate the forced continuation from position j onward.
			if isValid(num[0:i], num[i:j], num[j:]) {
				return true
			}
		}
	}
	return false
}

// isValid checks whether `rest` is the additive continuation of num1, num2.
// It computes the expected next sum (as a string) and peels it off `rest`,
// then recurses with the window slid forward by one.
func isValid(num1, num2, rest string) bool {
	if len(rest) == 0 {
		return true // consumed the whole string → valid sequence
	}
	sum := addStrings(num1, num2) // next expected number, overflow-safe
	// The remainder must start with exactly this sum.
	if len(rest) < len(sum) || rest[:len(sum)] != sum {
		return false
	}
	// Slide the window: (num2, sum) become the new pair, drop the matched sum.
	return isValid(num2, sum, rest[len(sum):])
}

// addStrings adds two non-negative decimal strings and returns the sum string.
// Manual column addition keeps arbitrarily large numbers exact.
func addStrings(a, b string) string {
	i, j := len(a)-1, len(b)-1 // start from least-significant digits
	carry := 0
	res := []byte{}
	for i >= 0 || j >= 0 || carry > 0 {
		sum := carry
		if i >= 0 {
			sum += int(a[i] - '0') // add digit of a
			i--
		}
		if j >= 0 {
			sum += int(b[j] - '0') // add digit of b
			j--
		}
		carry = sum / 10                    // carry into next column
		res = append(res, byte(sum%10)+'0') // store this column's digit
	}
	// res is built least-significant first → reverse it.
	for l, r := 0, len(res)-1; l < r; l, r = l+1, r-1 {
		res[l], res[r] = res[r], res[l]
	}
	return string(res)
}

// ── Approach 2: Iterative Split + big.Int Arithmetic (Optimal) ────────────────
//
// bigIntSplit solves Additive Number with the same two-pointer split over the
// first two numbers, but validates the continuation iteratively using
// math/big for clean arbitrary-precision addition.
//
// Intuition:
//
//	Identical search space as the brute force — every choice of the first two
//	numbers — but the verification is a tight loop instead of recursion, and we
//	lean on big.Int to add without hand-rolling column arithmetic. This is the
//	cleanest form to present in an interview once you have justified why 64-bit
//	ints are unsafe (numbers can have up to 35 digits given |num| ≤ 35).
//
// Algorithm:
//
//  1. Enumerate end index i of num1 and end index j of num2 (same leading-zero
//     pruning as approach 1).
//  2. Parse num1, num2 as big.Int; set pos = j.
//  3. While pos < n: compute next = num1 + num2, render as string; the string
//     from pos must start with it, else break; advance pos, shift the pair.
//  4. If pos reaches n exactly, the whole string was consumed → return true.
//
// Time:  O(n^2 · n) = O(n^3) — O(n^2) splits, each validated in O(n) big adds.
// Space: O(n) — big.Int operands proportional to number length.
func bigIntSplit(num string) bool {
	n := len(num)
	for i := 1; i <= n-2; i++ {
		if num[0] == '0' && i > 1 {
			break // leading zero in first number
		}
		first := new(big.Int)
		first.SetString(num[0:i], 10) // num1 as big integer
		for j := i + 1; j <= n-1; j++ {
			if num[i] == '0' && j-i > 1 {
				break // leading zero in second number
			}
			second := new(big.Int)
			second.SetString(num[i:j], 10) // num2 as big integer
			if validateFrom(num, j, new(big.Int).Set(first), second) {
				return true
			}
		}
	}
	return false
}

// validateFrom walks the suffix num[pos:] verifying each forced sum.
func validateFrom(num string, pos int, a, b *big.Int) bool {
	n := len(num)
	for pos < n {
		sum := new(big.Int).Add(a, b) // next expected number
		s := sum.String()             // decimal rendering
		// The suffix must begin with exactly this sum's digits.
		if pos+len(s) > n || num[pos:pos+len(s)] != s {
			return false
		}
		pos += len(s) // consume the matched number
		a, b = b, sum // slide the window forward
	}
	return pos == n // valid only if we consumed the entire string
}

func main() {
	fmt.Println("=== Approach 1: Brute Force Backtracking ===")
	fmt.Println(bruteForce("112358"))    // expected true
	fmt.Println(bruteForce("199100199")) // expected true
	fmt.Println(bruteForce("1203"))      // expected false
	fmt.Println(bruteForce("000"))       // expected true

	fmt.Println("=== Approach 2: Iterative Split + big.Int (Optimal) ===")
	fmt.Println(bigIntSplit("112358"))    // expected true
	fmt.Println(bigIntSplit("199100199")) // expected true
	fmt.Println(bigIntSplit("1203"))      // expected false
	fmt.Println(bigIntSplit("000"))       // expected true
}
