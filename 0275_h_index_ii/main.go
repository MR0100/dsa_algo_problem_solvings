package main

import "fmt"

// H-Index II (LeetCode #275)
//
// Given citations sorted in ASCENDING order, where citations[i] is the number
// of citations for the i-th paper, return the researcher's h-index. You must
// design an algorithm running in logarithmic time.
//
// The h-index: the maximum h such that at least h papers each have >= h
// citations.

// ── Approach 1: Linear Scan (Brute Force baseline) ───────────────────────────
//
// linearScan walks the ascending array to find the first index where the
// number of papers from that index to the end is <= the citation at that index.
//
// Intuition:
//
//	The array is sorted ascending, so at index i there are (n - i) papers with
//	at least citations[i] citations. The h-index is the largest value of
//	(n - i) for which citations[i] >= (n - i). Scanning left to right, the
//	first i where citations[i] >= n - i gives h = n - i (that suffix all
//	qualifies, and it's the biggest qualifying suffix).
//
// Algorithm:
//  1. For i from 0..n-1: if citations[i] >= n - i, return n - i.
//  2. If none qualify, return 0.
//
// Time:  O(n) — a single pass (does NOT meet the log-time requirement; baseline).
// Space: O(1).
func linearScan(citations []int) int {
	n := len(citations)
	for i := 0; i < n; i++ {
		// n - i papers (from i to the end) each have >= citations[i] citations.
		if citations[i] >= n-i {
			return n - i // largest qualifying suffix length
		}
	}
	return 0
}

// ── Approach 2: Binary Search (Optimal, O(log n)) ────────────────────────────
//
// binarySearch finds, in logarithmic time, the leftmost index i such that
// citations[i] >= n - i; the answer is then n - i.
//
// Intuition:
//
//	Define f(i) = (citations[i] >= n - i). As i increases, citations[i] is
//	non-decreasing while (n - i) decreases, so once f(i) becomes true it stays
//	true — f is monotone. Binary-search the first i where f(i) holds. Every
//	paper from that i onward qualifies, giving h = n - i. The sorted input is
//	exactly what makes this monotone predicate binary-searchable.
//
// Algorithm:
//  1. lo = 0, hi = n. While lo < hi: mid = (lo+hi)/2.
//  2. If citations[mid] >= n - mid, hi = mid (first-true is at or left of mid);
//     else lo = mid + 1.
//  3. Return n - lo.
//
// Time:  O(log n) — halves the search range each step.
// Space: O(1).
func binarySearch(citations []int) int {
	n := len(citations)
	lo, hi := 0, n // search for the first index satisfying the predicate; hi=n means "none"
	for lo < hi {
		mid := (lo + hi) / 2
		// Papers mid..n-1 (that's n-mid of them) each have >= citations[mid] cites.
		if citations[mid] >= n-mid {
			hi = mid // predicate holds → answer index is here or to the left
		} else {
			lo = mid + 1 // predicate fails → move right
		}
	}
	// lo is the first index where citations[lo] >= n - lo; if none, lo == n → 0.
	return n - lo
}

func main() {
	ex1 := []int{0, 1, 3, 5, 6} // expected 3
	ex2 := []int{1, 2, 100}     // expected 2

	fmt.Println("=== Approach 1: Linear Scan ===")
	fmt.Println(linearScan(ex1)) // expected 3
	fmt.Println(linearScan(ex2)) // expected 2

	fmt.Println("=== Approach 2: Binary Search ===")
	fmt.Println(binarySearch(ex1)) // expected 3
	fmt.Println(binarySearch(ex2)) // expected 2
}
