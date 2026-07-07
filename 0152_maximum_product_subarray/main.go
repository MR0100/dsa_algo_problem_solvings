package main

import "fmt"

// ── Approach 1: Brute Force ──────────────────────────────────────────────────
//
// bruteForce solves Maximum Product Subarray by checking every subarray.
//
// Intuition:
//
//	Every subarray is defined by its start index i and end index j. Fix i,
//	then extend j one step at a time while maintaining a running product —
//	that turns the naive O(n^3) triple loop into O(n^2) because the product
//	of nums[i..j] is just product(nums[i..j-1]) * nums[j].
//
// Algorithm:
//  1. Initialize best to nums[0] (a subarray must be non-empty).
//  2. For every start index i:
//     a. Reset the running product to 1.
//     b. For every end index j >= i, multiply the running product by
//     nums[j] and update best.
//  3. Return best.
//
// Time:  O(n^2) — n choices of start, each extended up to n times.
// Space: O(1) — only two scalar accumulators.
func bruteForce(nums []int) int {
	best := nums[0] // best product seen so far; seeded with a valid subarray
	for i := 0; i < len(nums); i++ {
		prod := 1 // running product of nums[i..j]
		for j := i; j < len(nums); j++ {
			prod *= nums[j] // extend the subarray by one element
			if prod > best {
				best = prod // record a new maximum
			}
		}
	}
	return best
}

// ── Approach 2: DP with Max/Min Tracking (Kadane Variant, Optimal) ───────────
//
// dpMinMax solves Maximum Product Subarray by tracking both the maximum AND
// minimum product of a subarray ending at each index.
//
// Intuition:
//
//	Kadane's algorithm for maximum SUM fails here because a negative number
//	flips signs: a very NEGATIVE running product becomes very POSITIVE when
//	multiplied by another negative. So the "most promising" subarray ending
//	at i is not only the one with the largest product — the one with the
//	smallest (most negative) product is equally promising for the future.
//	Track both. When nums[i] is negative, the roles swap.
//
// Algorithm:
//  1. Initialize maxEnd, minEnd, best to nums[0].
//  2. For each i from 1 to n-1:
//     a. If nums[i] < 0, swap maxEnd and minEnd (multiplying by a negative
//     turns the largest into the smallest and vice versa).
//     b. maxEnd = max(nums[i], maxEnd*nums[i]) — either extend the best
//     subarray or start fresh at nums[i].
//     c. minEnd = min(nums[i], minEnd*nums[i]) — same choice for the worst.
//     d. best = max(best, maxEnd).
//  3. Return best.
//
// Time:  O(n) — single pass, O(1) work per element.
// Space: O(1) — three scalars; only the previous DP state is consulted.
func dpMinMax(nums []int) int {
	maxEnd := nums[0] // max product of a subarray ENDING exactly at i
	minEnd := nums[0] // min product of a subarray ENDING exactly at i
	best := nums[0]   // global answer over all end positions
	for i := 1; i < len(nums); i++ {
		n := nums[i]
		if n < 0 {
			// a negative factor turns the biggest product into the smallest
			// and the smallest into the biggest — swap before extending
			maxEnd, minEnd = minEnd, maxEnd
		}
		// either extend the previous subarray or start a new one at n
		maxEnd = max(n, maxEnd*n)
		minEnd = min(n, minEnd*n)
		if maxEnd > best {
			best = maxEnd // new global maximum
		}
	}
	return best
}

// ── Approach 3: Prefix/Suffix Sweep (Optimal, No DP State) ───────────────────
//
// prefixSuffix solves Maximum Product Subarray with two running-product
// sweeps: left→right and right→left, resetting at zeros.
//
// Intuition:
//
//	Between zeros, the answer is either the product of the WHOLE block (even
//	number of negatives) or the block with one negative-and-everything-before
//	-it (or after it) chopped off. Chopping "before the first negative" is
//	captured by the suffix sweep; chopping "after the last negative" by the
//	prefix sweep. So the best prefix or suffix product inside each zero-free
//	block always contains the answer. Zeros simply reset the running product
//	to 1 (a zero can still BE the answer, which the comparison before the
//	reset handles).
//
// Algorithm:
//  1. Sweep left→right keeping a running product; after each multiply,
//     update best; if the product hits 0, reset it to 1.
//  2. Sweep right→left the same way.
//  3. Return best.
//
// Time:  O(n) — two linear passes.
// Space: O(1) — one running product and the answer.
func prefixSuffix(nums []int) int {
	best := nums[0] // must hold at least one element
	prod := 1       // running prefix product
	for i := 0; i < len(nums); i++ {
		prod *= nums[i] // extend the prefix product
		if prod > best {
			best = prod
		}
		if prod == 0 {
			prod = 1 // a zero kills every product through it — restart after it
		}
	}
	prod = 1 // reset for the suffix sweep
	for i := len(nums) - 1; i >= 0; i-- {
		prod *= nums[i] // extend the suffix product
		if prod > best {
			best = prod
		}
		if prod == 0 {
			prod = 1 // restart past the zero
		}
	}
	return best
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(bruteForce([]int{2, 3, -2, 4})) // 6
	fmt.Println(bruteForce([]int{-2, 0, -1}))   // 0

	fmt.Println("=== Approach 2: DP with Max/Min Tracking ===")
	fmt.Println(dpMinMax([]int{2, 3, -2, 4})) // 6
	fmt.Println(dpMinMax([]int{-2, 0, -1}))   // 0

	fmt.Println("=== Approach 3: Prefix/Suffix Sweep (Optimal) ===")
	fmt.Println(prefixSuffix([]int{2, 3, -2, 4})) // 6
	fmt.Println(prefixSuffix([]int{-2, 0, -1}))   // 0
}
