package main

import "fmt"

// ── Approach 1: Brute Force (Three Nested Loops) ─────────────────────────────
//
// bruteForce checks every ordered triple (i<j<k) for nums[i]<nums[j]<nums[k].
//
// Intuition:
//
//	The definition is "does some increasing triple of indices exist?" — so
//	enumerate all triples i<j<k and test the inequality. Trivially correct,
//	hopelessly slow for large n.
//
// Algorithm:
//  1. For every i, for every j>i with nums[j]>nums[i], for every k>j with
//     nums[k]>nums[j]: return true.
//  2. If no triple qualifies, return false.
//
// Time:  O(n^3) — three nested loops.
// Space: O(1).
func bruteForce(nums []int) bool {
	n := len(nums)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if nums[j] <= nums[i] {
				continue // need nums[i] < nums[j]
			}
			for k := j + 1; k < n; k++ {
				if nums[k] > nums[j] {
					return true // found nums[i] < nums[j] < nums[k]
				}
			}
		}
	}
	return false
}

// ── Approach 2: Precomputed Min-Left / Max-Right ─────────────────────────────
//
// minLeftMaxRight finds a middle index j that has a smaller element to its
// left and a larger element to its right.
//
// Intuition:
//
//	A triple exists iff some index j is the "middle": there is something
//	smaller before it and something larger after it. Precompute prefixMin[j]
//	= min of nums[0..j] and suffixMax[j] = max of nums[j..n-1]. Then j works
//	iff prefixMin[j-1] < nums[j] < suffixMax[j+1].
//
// Algorithm:
//  1. Build prefixMin (running min from the left).
//  2. Build suffixMax (running max from the right).
//  3. For each middle j in [1, n-2]: if prefixMin[j-1] < nums[j] < suffixMax[j+1] → true.
//
// Time:  O(n) — three linear passes.
// Space: O(n) — the two auxiliary arrays.
func minLeftMaxRight(nums []int) bool {
	n := len(nums)
	if n < 3 {
		return false // need at least three elements
	}
	prefixMin := make([]int, n) // prefixMin[i] = min(nums[0..i])
	prefixMin[0] = nums[0]
	for i := 1; i < n; i++ {
		prefixMin[i] = min(prefixMin[i-1], nums[i]) // extend the running min
	}
	suffixMax := make([]int, n) // suffixMax[i] = max(nums[i..n-1])
	suffixMax[n-1] = nums[n-1]
	for i := n - 2; i >= 0; i-- {
		suffixMax[i] = max(suffixMax[i+1], nums[i]) // extend the running max
	}
	for j := 1; j < n-1; j++ {
		// j is a valid middle if something smaller is on its left and
		// something larger is on its right.
		if prefixMin[j-1] < nums[j] && nums[j] < suffixMax[j+1] {
			return true
		}
	}
	return false
}

// ── Approach 3: Two Smallest (Greedy, Optimal) ───────────────────────────────
//
// twoSmallest tracks the smallest and second-smallest tails seen so far; a
// third element beating both proves a triple exists.
//
// Intuition:
//
//	Keep two values: `first` = smallest number so far, `second` = smallest
//	number that has some strictly smaller number before it (a valid "middle").
//	Scan left to right. If a number ≤ first, it's a new potential start. Else
//	if it's ≤ second, it's a better middle (and it implicitly had a smaller
//	`first` before it). Else it's larger than both → it can be the third,
//	completing an increasing triple. Reassigning `first` later is safe: it
//	only records a smaller start that occurred BEFORE the current `second`
//	was set, so the ordering invariant still holds.
//
// Algorithm:
//  1. first = second = +∞.
//  2. For each x: if x <= first, first = x; else if x <= second, second = x;
//     else return true (x > second > (some earlier first)).
//  3. Return false.
//
// Time:  O(n) — one pass.
// Space: O(1) — two scalars.
func twoSmallest(nums []int) bool {
	first, second := 1<<62, 1<<62 // smallest and second-smallest valid tails
	for _, x := range nums {
		switch {
		case x <= first:
			first = x // new smallest candidate for the triple's start
		case x <= second:
			second = x // x can serve as a middle (some smaller `first` precedes it)
		default:
			return true // x beats both → increasing triple exists
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	ex1 := []int{1, 2, 3, 4, 5}
	ex2 := []int{5, 4, 3, 2, 1}
	ex3 := []int{2, 1, 5, 0, 4, 6}

	fmt.Println("=== Approach 1: Brute Force (Three Nested Loops) ===")
	fmt.Println(bruteForce(ex1)) // expected true
	fmt.Println(bruteForce(ex2)) // expected false
	fmt.Println(bruteForce(ex3)) // expected true

	fmt.Println("=== Approach 2: Precomputed Min-Left / Max-Right ===")
	fmt.Println(minLeftMaxRight(ex1)) // expected true
	fmt.Println(minLeftMaxRight(ex2)) // expected false
	fmt.Println(minLeftMaxRight(ex3)) // expected true

	fmt.Println("=== Approach 3: Two Smallest (Greedy, Optimal) ===")
	fmt.Println(twoSmallest(ex1)) // expected true
	fmt.Println(twoSmallest(ex2)) // expected false
	fmt.Println(twoSmallest(ex3)) // expected true
}
