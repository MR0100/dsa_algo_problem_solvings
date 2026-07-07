package main

import "fmt"

// pick is the hidden number the judge chose. On real LeetCode it is hidden and
// only reachable via guess(); here we set it explicitly to exercise the code.
var pick int

// guess is the pre-defined API. It returns:
//
//	-1 if num is HIGHER than the picked number,
//	 1 if num is LOWER  than the picked number,
//	 0 if num IS the picked number.
func guess(num int) int {
	switch {
	case num > pick: // our guess overshot
		return -1
	case num < pick: // our guess undershot
		return 1
	default: // exact match
		return 0
	}
}

// ── Approach 1: Linear Scan (Brute Force) ────────────────────────────────────
//
// linearScan ignores the higher/lower hint and simply tries 1, 2, 3, ... until
// guess returns 0. Included as the naive baseline.
//
// Intuition:
//
//	guess(num)==0 pinpoints the answer, so walking upward from 1 must hit it.
//	Wastes the ordering information the API gives us.
//
// Algorithm:
//  1. For num = 1..n: if guess(num)==0 return num.
//  2. (Unreachable) return -1.
//
// Time:  O(n) — up to n guesses.
// Space: O(1).
func linearScan(n int) int {
	for num := 1; num <= n; num++ { // try every candidate in order
		if guess(num) == 0 { // found the picked number
			return num
		}
	}
	return -1 // per constraints the pick is always in [1, n]
}

// ── Approach 2: Binary Search (Optimal) ──────────────────────────────────────
//
// binarySearch uses the higher/lower feedback to halve the search space each
// guess.
//
// Intuition:
//
//	The response tells us which half of [lo, hi] contains the pick, so this is a
//	textbook binary search over a monotone predicate. We compute mid with
//	lo + (hi-lo)/2 to avoid overflow when n is near 2^31−1.
//
// Algorithm:
//  1. lo = 1, hi = n.
//  2. While lo <= hi:
//     a. mid = lo + (hi-lo)/2.
//     b. r = guess(mid). If r==0 return mid.
//     c. If r < 0 (mid too high) → hi = mid-1, else lo = mid+1.
//  3. (Unreachable) return -1.
//
// Time:  O(log n) — space halves each step.
// Space: O(1).
func binarySearch(n int) int {
	lo, hi := 1, n
	for lo <= hi {
		mid := lo + (hi-lo)/2 // overflow-safe midpoint
		switch guess(mid) {
		case 0: // exact hit
			return mid
		case -1: // mid is higher than pick → search lower half
			hi = mid - 1
		default: // guess returned 1: mid is lower → search upper half
			lo = mid + 1
		}
	}
	return -1
}

func main() {
	fmt.Println("=== Approach 1: Linear Scan ===")
	pick = 6
	fmt.Println(linearScan(10)) // expected 6
	pick = 1
	fmt.Println(linearScan(1)) // expected 1
	pick = 1
	fmt.Println(linearScan(2)) // expected 1

	fmt.Println("=== Approach 2: Binary Search ===")
	pick = 6
	fmt.Println(binarySearch(10)) // expected 6
	pick = 1
	fmt.Println(binarySearch(1)) // expected 1
	pick = 1702766719
	fmt.Println(binarySearch(2147483647)) // expected 1702766719
}
