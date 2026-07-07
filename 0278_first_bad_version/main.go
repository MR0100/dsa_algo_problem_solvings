package main

import "fmt"

// The judge provides isBadVersion(version) which returns true once a version is
// bad. Versions are monotone: once one version is bad, every later version is
// bad too. We model that with a global "firstBad" threshold.

var firstBad int // versions >= firstBad are bad

// isBadVersion is the provided API.
func isBadVersion(version int) bool {
	return version >= firstBad
}

// ── Approach 1: Linear Scan (Brute Force) ────────────────────────────────────
//
// linearScan solves First Bad Version by testing versions one by one.
//
// Intuition:
//
//	Because badness is monotone, the very first version that returns true is the
//	answer. Walk from 1 upward and stop at the first bad one.
//
// Algorithm:
//  1. For v = 1..n: if isBadVersion(v) return v.
//  2. (Guaranteed to find one within n.)
//
// Time:  O(n) — up to n API calls.
// Space: O(1).
func linearScan(n int) int {
	for v := 1; v <= n; v++ { // scan versions in order
		if isBadVersion(v) { // first true = first bad version
			return v
		}
	}
	return -1 // unreachable per problem guarantees
}

// ── Approach 2: Binary Search (Optimal) ──────────────────────────────────────
//
// binarySearch solves First Bad Version by halving the search interval.
//
// Intuition:
//
//	The sequence looks like [good, good, ..., good, bad, bad, ..., bad]. Finding
//	the first bad is finding the boundary in a sorted boolean array — a textbook
//	binary search for the leftmost true.
//
// Algorithm:
//  1. lo = 1, hi = n.
//  2. While lo < hi:
//     mid = lo + (hi-lo)/2   // avoid overflow
//     if isBadVersion(mid): hi = mid   // answer is mid or to its left
//     else:                 lo = mid+1 // answer is strictly right
//  3. lo == hi is the first bad version.
//
// Time:  O(log n) — the interval halves each step.
// Space: O(1).
func binarySearch(n int) int {
	lo, hi := 1, n // search space is [1, n]
	for lo < hi {  // stop when the interval collapses to one version
		mid := lo + (hi-lo)/2 // midpoint without integer overflow
		if isBadVersion(mid) {
			hi = mid // mid is bad → first bad is mid or earlier; keep mid
		} else {
			lo = mid + 1 // mid is good → first bad is strictly after mid
		}
	}
	return lo // lo == hi points at the boundary: the first bad version
}

func main() {
	// Example 1: n = 5, first bad = 4 → answer 4.
	fmt.Println("=== Approach 1: Linear Scan ===")
	firstBad = 4
	fmt.Println(linearScan(5)) // expected 4

	// Example 2: n = 1, first bad = 1 → answer 1.
	firstBad = 1
	fmt.Println(linearScan(1)) // expected 1

	fmt.Println("=== Approach 2: Binary Search (Optimal) ===")
	firstBad = 4
	fmt.Println(binarySearch(5)) // expected 4

	firstBad = 1
	fmt.Println(binarySearch(1)) // expected 1
}
