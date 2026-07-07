package main

import (
	"fmt"
	"math"
	"math/bits"
	"strconv"
)

// A "good base" k of n (k ≥ 2) writes n as all 1's, i.e.
//
//	n = 1 + k + k² + … + k^m   (that is m+1 ones in base k).
//
// We want the SMALLEST such k. Smaller base ⇒ more digits, so the smallest k
// corresponds to the LARGEST digit count m+1. Two facts bound the search:
//   - k = n-1 is always a good base: n = 11 in base n-1 (two 1's). So an answer
//     always exists, and it never exceeds n-1.
//   - With k ≥ 2 and m+1 digits, n ≥ 2^m, hence m ≤ log2(n) < 60 for n ≤ 10^18.
//
// So we try digit counts from many (large m) down to few, and for each m find
// the base k whose repunit equals n. The first hit (largest m) is the answer.

// repunitCmp compares 1 + k + … + k^m against target WITHOUT overflowing uint64.
// Returns -1 if the repunit < target, 0 if equal, +1 if greater. The instant a
// partial product/sum would exceed uint64 (or target) we bail out as +1, which
// is exactly what binary search needs ("this k is too big").
func repunitCmp(k uint64, m int, target uint64) int {
	var sum uint64 = 1  // the k^0 = 1 term
	var term uint64 = 1 // running k^i, starts at k^0
	for i := 1; i <= m; i++ {
		// term *= k, but detect overflow first via 128-bit multiply.
		hi, lo := bits.Mul64(term, k)
		if hi != 0 { // product needs more than 64 bits → definitely too big
			return 1
		}
		term = lo
		// sum += term with overflow guard.
		if sum > math.MaxUint64-term {
			return 1
		}
		sum += term
		if sum > target { // already past the goal — no need to keep going
			return 1
		}
	}
	if sum < target {
		return -1
	}
	return 0
}

// ── Approach 1: Brute Force Over Bases ───────────────────────────────────────
//
// bruteForce tries every candidate base k from 2 upward and, for each, checks
// whether repeatedly writing n in base k yields all 1's.
//
// Intuition:
//
//	"Smallest good base" literally asks for the least k ≥ 2 such that n in base
//	k is a string of 1's. So test k = 2, 3, 4, … and return the first that
//	works. Checking a base means: while n > 0, the remainder n % k must be 1,
//	then n /= k. If we ever see a remainder ≠ 1 the base fails. k = n-1 is a
//	guaranteed backstop, so the loop always terminates.
//
// Algorithm:
//  1. For k = 2 … n-1:
//     a. Set cur = n. While cur > 1: if cur % k != 1 break (fail); else cur /= k.
//     b. If the loop consumed cur down to exactly 1, k is good → return k.
//  2. Fallback return n-1 (always valid).
//
// Time:  O(n · log n) worst case — up to n candidate bases, each verified in
//
//	O(log_k n) divisions. Fine for tiny n; hopeless near 10^18 (shown only
//	for correctness on small inputs).
//
// Space: O(1).
func bruteForce(nStr string) string {
	n, _ := strconv.ParseUint(nStr, 10, 64)
	for k := uint64(2); k < n-1; k++ {
		cur := n
		good := true
		for cur > 1 { // peel digits of n in base k, low to high
			if cur%k != 1 { // every digit must be exactly 1
				good = false
				break
			}
			cur /= k // drop the digit we just verified
		}
		if good && cur == 1 { // consumed everything and ended on the leading 1
			return strconv.FormatUint(k, 10)
		}
	}
	return strconv.FormatUint(n-1, 10) // n = "11" in base n-1 always works
}

// ── Approach 2: Binary Search on the Base per Digit-Count ─────────────────────
//
// binarySearch iterates the number of digits m+1 from largest possible down to
// 2, and for each fixed m binary-searches the base k in [2, n^(1/m)] whose
// repunit 1+k+…+k^m equals n.
//
// Intuition:
//
//	Fix the digit count. Then f(k) = 1 + k + … + k^m is strictly increasing in
//	k, so there is at most one k making f(k) = n, and we can binary-search it.
//	Because the smallest base ⇔ the most digits, we scan m from high to low and
//	return the FIRST k we find — that necessarily uses the maximum digits and
//	hence is the minimum base. Upper bound for k at digit count m+1 is
//	⌊n^(1/m)⌋: since n > k^m, k < n^(1/m).
//
// Algorithm:
//  1. maxM = ⌊log2 n⌋ (largest possible top exponent, since n ≥ 2^m).
//  2. For m = maxM down to 1:
//     - lo = 2, hi = ⌊n^(1/m)⌋ + 1.
//     - Binary-search k: compare repunit(k,m) to n (overflow-safe). Return k on
//     an exact match.
//  3. If nothing matched, return n-1 (the m = 1 / two-digit case).
//
// Time:  O(log²n) — outer loop ~log n values of m, each an O(log n) binary
//
//	search, each comparison O(m) = O(log n)… bounded overall by ~log³n but in
//	practice tiny (≤ ~60 · 60 · 60 ops).
//
// Space: O(1).
func binarySearch(nStr string) string {
	n, _ := strconv.ParseUint(nStr, 10, 64)

	maxM := int(math.Log2(float64(n))) // most possible digits minus one
	for m := maxM; m >= 1; m-- {       // more digits first ⇒ smaller base first
		lo := uint64(2)
		hi := uint64(math.Pow(float64(n), 1.0/float64(m))) + 1 // k < n^(1/m), +1 for float slack
		for lo <= hi {
			mid := lo + (hi-lo)/2
			switch repunitCmp(mid, m, n) {
			case 0:
				return strconv.FormatUint(mid, 10) // exact repunit → smallest base found
			case -1:
				lo = mid + 1 // repunit too small → need a bigger base
			default:
				hi = mid - 1 // repunit too big → need a smaller base
			}
		}
	}
	return strconv.FormatUint(n-1, 10) // guaranteed two-digit fallback
}

// ── Approach 3: Direct m-th Root Estimate (Optimal) ──────────────────────────
//
// mthRootEstimate skips the inner binary search: for each digit count m+1 it
// estimates the base directly as k ≈ ⌊n^(1/m)⌋ (from the dominant term k^m ≈ n)
// and verifies that single candidate.
//
// Intuition:
//
//	For m+1 digits, n = k^m + k^(m-1) + … + 1 is dominated by k^m, so
//	k = ⌊n^(1/m)⌋ is essentially forced — the only base that could possibly
//	work for that digit count. Compute it with a real m-th root, then verify
//	the exact repunit with the overflow-safe comparison (guarding against the
//	floating-point estimate being off by one). Scanning m from large to small
//	and returning the first candidate that verifies gives the smallest base
//	with just ONE check per m instead of a whole binary search.
//
// Algorithm:
//  1. For m = ⌊log2 n⌋ down to 2:
//     - k = ⌊n^(1/m)⌋ via math.Pow (this is the mandatory base for m digits).
//     - If k < 2, skip. Else verify repunit(k,m) == n; if so return k.
//  2. Fallback: return n-1 (the m = 1 case, base n-1, digits "11").
//
// Time:  O(log n · log n) — ~log n values of m, each an O(m)=O(log n) verify.
// Space: O(1).
func mthRootEstimate(nStr string) string {
	n, _ := strconv.ParseUint(nStr, 10, 64)

	// m is the top exponent; digit count is m+1. Largest m first ⇒ smallest base.
	for m := int(math.Log2(float64(n))); m >= 2; m-- {
		// Dominant-term estimate: n ≈ k^m ⇒ k ≈ n^(1/m).
		k := uint64(math.Pow(float64(n), 1.0/float64(m)))
		if k < 2 {
			continue // a valid base must be ≥ 2
		}
		// The float root can be off by one; verify the true repunit exactly.
		if repunitCmp(k, m, n) == 0 {
			return strconv.FormatUint(k, 10) // first verified ⇒ maximum digits ⇒ min base
		}
	}
	return strconv.FormatUint(n-1, 10) // m = 1: n written as "11" in base n-1
}

func main() {
	fmt.Println("=== Approach 1: Brute Force Over Bases ===")
	fmt.Printf("n=%q       got=%q  expected \"3\"\n", "13", bruteForce("13"))   // 13 = 111 base 3
	fmt.Printf("n=%q     got=%q  expected \"8\"\n", "4681", bruteForce("4681")) // 4681 = 11111 base 8
	fmt.Printf("n=%q        got=%q  expected \"2\"\n", "7", bruteForce("7"))    // 7 = 111 base 2
	// n = 10^18 skipped for brute force (would scan ~10^18 bases — TLE by design)

	fmt.Println("=== Approach 2: Binary Search on the Base per Digit-Count ===")
	fmt.Printf("n=%q       got=%q  expected \"3\"\n", "13", binarySearch("13"))
	fmt.Printf("n=%q     got=%q  expected \"8\"\n", "4681", binarySearch("4681"))
	fmt.Printf("n=%q        got=%q  expected \"2\"\n", "7", binarySearch("7"))
	fmt.Printf("n=%q  got=%q  expected \"999999999999999999\"\n", "1000000000000000000", binarySearch("1000000000000000000"))

	fmt.Println("=== Approach 3: Direct m-th Root Estimate (Optimal) ===")
	fmt.Printf("n=%q       got=%q  expected \"3\"\n", "13", mthRootEstimate("13"))
	fmt.Printf("n=%q     got=%q  expected \"8\"\n", "4681", mthRootEstimate("4681"))
	fmt.Printf("n=%q        got=%q  expected \"2\"\n", "7", mthRootEstimate("7"))
	fmt.Printf("n=%q  got=%q  expected \"999999999999999999\"\n", "1000000000000000000", mthRootEstimate("1000000000000000000"))
}
