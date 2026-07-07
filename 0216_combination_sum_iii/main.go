package main

import "fmt"

// ── Approach 1: Backtracking (Optimal) ───────────────────────────────────────
//
// backtracking solves Combination Sum III by exploring the digits 1..9 in
// increasing order, choosing each digit at most once, and pruning branches
// that can no longer produce a valid combination.
//
// Intuition:
//
//	We must pick exactly k DISTINCT digits from 1..9 that sum to n. Because
//	every digit may be used at most once and order does not matter, we walk
//	the digits left-to-right and, at each digit, decide "include it or skip
//	it". Forcing the next digit to be strictly greater than the current one
//	(via a `start` index) guarantees each combination is generated in sorted
//	order exactly once, so there are no duplicates to filter out.
//
// Algorithm:
//  1. DFS carries (start, remaining count still to pick, remaining sum) and the
//     current partial combination `path`.
//  2. Base success: count == 0 AND sum == 0 → record a copy of path.
//  3. Prune: if count == 0 or sum <= 0 (and not the success case) → dead end.
//  4. For digit d from start..9: if d > remaining sum, break (all later digits
//     are even larger — sorted, so no point continuing).
//  5. Choose d, recurse with (d+1, count-1, sum-d), then un-choose (backtrack).
//
// Time:  O(C(9, k) · k) — at most C(9,k) valid combinations, each copied in
//
//	O(k). The search tree over 9 digits is tiny and bounded.
//
// Space: O(k) — recursion depth and the path buffer; output not counted.
func backtracking(k int, n int) [][]int {
	var result [][]int        // collected valid combinations
	path := make([]int, 0, k) // current partial combination being built

	// dfs tries digits from `start` upward, needing `count` more digits that
	// together sum to `sum`.
	var dfs func(start, count, sum int)
	dfs = func(start, count, sum int) {
		if count == 0 { // we have chosen exactly k digits
			if sum == 0 { // and they hit the target sum
				combo := make([]int, len(path)) // copy: path is mutated after return
				copy(combo, path)
				result = append(result, combo)
			}
			return // either recorded or over-shot the sum with k digits — stop
		}
		// Try each candidate digit in strictly increasing order.
		for d := start; d <= 9; d++ {
			if d > sum { // digits are sorted; d and everything after overshoots
				break
			}
			path = append(path, d)    // choose d
			dfs(d+1, count-1, sum-d)  // recurse: next digit must exceed d
			path = path[:len(path)-1] // un-choose d (backtrack)
		}
	}

	dfs(1, k, n) // start from digit 1, needing k digits summing to n
	return result
}

// ── Approach 2: Bitmask Enumeration ──────────────────────────────────────────
//
// bitmaskEnumeration solves Combination Sum III by enumerating every subset of
// the 9 digits via a 9-bit mask and keeping subsets of exactly size k whose
// elements sum to n.
//
// Intuition:
//
//	There are only 9 digits, so there are just 2^9 = 512 subsets. That is
//	small enough to brute-force: represent a subset as a bitmask where bit i
//	(0-indexed) means "digit i+1 is included". For each mask, count its set
//	bits (that is the subset size) and sum the corresponding digits, then
//	keep the mask iff size == k and sum == n.
//
// Algorithm:
//  1. For mask from 0 to 511:
//  2. Walk bits 0..8; if bit i is set, add (i+1) to sum and to a temp list.
//  3. If the temp list has exactly k elements and sum == n, record it.
//  4. Return all recorded lists. Each list is internally sorted (bits scanned
//     low→high), though the lists themselves come out in mask order, not sorted
//     lexicographically. LeetCode accepts combinations in any order.
//
// Time:  O(2^9 · 9) = O(4608) — constant work; independent of k and n.
// Space: O(k) scratch per mask; output not counted.
func bitmaskEnumeration(k int, n int) [][]int {
	var result [][]int
	// 1<<9 == 512 subsets of the digit set {1..9}.
	for mask := 0; mask < (1 << 9); mask++ {
		sum := 0                   // running sum of chosen digits
		combo := make([]int, 0, k) // chosen digits for this mask
		for i := 0; i < 9; i++ {   // inspect each of the 9 candidate digits
			if mask&(1<<i) != 0 { // bit i set → digit (i+1) is included
				digit := i + 1
				sum += digit
				combo = append(combo, digit)
			}
		}
		// Keep only subsets of the right size that hit the target sum.
		if len(combo) == k && sum == n {
			result = append(result, combo)
		}
	}
	return result
}

func main() {
	fmt.Println("=== Approach 1: Backtracking (Optimal) ===")
	fmt.Println(backtracking(3, 7)) // expected [[1 2 4]]
	fmt.Println(backtracking(3, 9)) // expected [[1 2 6] [1 3 5] [2 3 4]]
	fmt.Println(backtracking(4, 1)) // expected [] (no 4 distinct digits sum to 1)

	fmt.Println("=== Approach 2: Bitmask Enumeration ===")
	fmt.Println(bitmaskEnumeration(3, 7)) // expected [[1 2 4]]
	fmt.Println(bitmaskEnumeration(3, 9)) // expected [[2 3 4] [1 3 5] [1 2 6]] (any order accepted)
	fmt.Println(bitmaskEnumeration(4, 1)) // expected []
}
