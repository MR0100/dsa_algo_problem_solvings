package main

import "fmt"

// ── Approach 1: Brute Force ──────────────────────────────────────────────────
//
// bruteForce solves Product of Array Except Self by, for each index, multiplying
// every OTHER element with a nested loop.
//
// Intuition:
//
//	The definition is literal: answer[i] is the product of all nums[j] with
//	j != i. Just compute that directly with two loops. It ignores the "no
//	division, O(n)" follow-up but is the obvious baseline.
//
// Algorithm:
//  1. For each i, initialise product = 1.
//  2. For each j != i, multiply product by nums[j].
//  3. Store product in answer[i].
//
// Time:  O(n²) — nested loops.
// Space: O(1) extra (output array excluded).
func bruteForce(nums []int) []int {
	n := len(nums)
	answer := make([]int, n)
	for i := 0; i < n; i++ {
		product := 1 // running product of every element except nums[i]
		for j := 0; j < n; j++ {
			if j != i {
				product *= nums[j] // multiply in every other element
			}
		}
		answer[i] = product
	}
	return answer
}

// ── Approach 2: Prefix × Suffix Arrays ───────────────────────────────────────
//
// prefixSuffixArrays solves Product of Array Except Self by precomputing, for
// each index, the product of everything to its left and everything to its right,
// then multiplying the two.
//
// Intuition:
//
//	The product of "all except i" splits cleanly into (product of elements left
//	of i) × (product of elements right of i). Precompute both directions in two
//	sweeps; answer[i] is just their product. No division needed.
//
// Algorithm:
//  1. prefix[i] = product of nums[0..i-1] (prefix[0] = 1).
//  2. suffix[i] = product of nums[i+1..n-1] (suffix[n-1] = 1).
//  3. answer[i] = prefix[i] * suffix[i].
//
// Time:  O(n) — three linear passes.
// Space: O(n) — two auxiliary arrays.
func prefixSuffixArrays(nums []int) []int {
	n := len(nums)
	prefix := make([]int, n) // prefix[i] = product of everything left of i
	suffix := make([]int, n) // suffix[i] = product of everything right of i
	answer := make([]int, n)

	prefix[0] = 1
	for i := 1; i < n; i++ {
		prefix[i] = prefix[i-1] * nums[i-1] // accumulate leftward products
	}
	suffix[n-1] = 1
	for i := n - 2; i >= 0; i-- {
		suffix[i] = suffix[i+1] * nums[i+1] // accumulate rightward products
	}
	for i := 0; i < n; i++ {
		answer[i] = prefix[i] * suffix[i] // left product × right product
	}
	return answer
}

// ── Approach 3: Two-Pass O(1) Extra Space (Optimal) ──────────────────────────
//
// prefixSuffixInPlace solves Product of Array Except Self in O(n) time and O(1)
// extra space by folding the prefix products into the output array, then
// multiplying the suffix products in on a second reverse pass.
//
// Intuition:
//
//	The output array itself can hold the prefix products (that slot doesn't
//	count as extra space per the problem). Then a single reverse pass carries a
//	running suffix product and multiplies it into each slot — combining left and
//	right without any second array.
//
// Algorithm:
//  1. First pass: answer[i] = product of all elements left of i.
//  2. Second pass (right→left): keep running suffix `R`; answer[i] *= R; then
//     R *= nums[i].
//
// Time:  O(n) — two passes.
// Space: O(1) extra — only the output array plus one scalar.
func prefixSuffixInPlace(nums []int) []int {
	n := len(nums)
	answer := make([]int, n)

	answer[0] = 1
	for i := 1; i < n; i++ {
		answer[i] = answer[i-1] * nums[i-1] // answer[i] = product of left part
	}
	right := 1 // running product of everything to the right of i
	for i := n - 1; i >= 0; i-- {
		answer[i] *= right // fold in the right-side product
		right *= nums[i]   // extend the suffix product to include nums[i]
	}
	return answer
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(bruteForce([]int{1, 2, 3, 4}))      // expected [24 12 8 6]
	fmt.Println(bruteForce([]int{-1, 1, 0, -3, 3})) // expected [0 0 9 0 0]

	fmt.Println("=== Approach 2: Prefix × Suffix Arrays ===")
	fmt.Println(prefixSuffixArrays([]int{1, 2, 3, 4}))      // expected [24 12 8 6]
	fmt.Println(prefixSuffixArrays([]int{-1, 1, 0, -3, 3})) // expected [0 0 9 0 0]

	fmt.Println("=== Approach 3: Two-Pass O(1) Extra Space (Optimal) ===")
	fmt.Println(prefixSuffixInPlace([]int{1, 2, 3, 4}))      // expected [24 12 8 6]
	fmt.Println(prefixSuffixInPlace([]int{-1, 1, 0, -3, 3})) // expected [0 0 9 0 0]
}
