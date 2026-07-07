package main

import "fmt"

// LeetCode 303 — Range Sum Query - Immutable
//
// Given an integer array nums, handle multiple queries of the form: return the
// sum of nums[left..right] inclusive. The array is immutable (never updated),
// so we can preprocess it once and answer each query in O(1).

// ── Approach 1: Brute Force (Re-sum Each Query) ──────────────────────────────
//
// NumArrayBrute stores the raw array and re-adds the range on every query.
//
// Intuition:
//
//	The most literal implementation: keep the numbers, and when asked for a
//	range just loop and add. Correct, but every query is O(n); with q queries
//	the total is O(q·n), which is wasteful when the array never changes.
//
// Time:  constructor O(1); SumRange O(n) per query.
// Space: O(n) — stores the array.
type NumArrayBrute struct {
	nums []int // the original, immutable numbers
}

// NewNumArrayBrute stores the array by reference (it is never mutated).
func NewNumArrayBrute(nums []int) NumArrayBrute {
	return NumArrayBrute{nums: nums}
}

// SumRange adds nums[left..right] on demand.
//
// Time:  O(right − left + 1).
// Space: O(1).
func (a NumArrayBrute) SumRange(left, right int) int {
	sum := 0
	for i := left; i <= right; i++ { // walk the requested window
		sum += a.nums[i]
	}
	return sum
}

// ── Approach 2: Prefix Sum (Optimal) ─────────────────────────────────────────
//
// NumArray precomputes cumulative sums so any range answers in O(1).
//
// Intuition:
//
//	Let prefix[i] = nums[0] + ... + nums[i−1] (prefix[0] = 0). Then the sum of
//	nums[left..right] telescopes to prefix[right+1] − prefix[left]: the second
//	term subtracts off everything before `left`, leaving exactly the window.
//	One O(n) preprocessing pass buys O(1) queries forever after.
//
// Algorithm:
//  1. Build prefix of length n+1 with prefix[0] = 0 and
//     prefix[i+1] = prefix[i] + nums[i].
//  2. SumRange(left, right) = prefix[right+1] − prefix[left].
//
// Time:  constructor O(n); SumRange O(1) per query.
// Space: O(n) — the prefix array.
type NumArray struct {
	prefix []int // prefix[i] = sum of the first i elements
}

// NewNumArray builds the cumulative-sum table once.
//
// Time:  O(n).
// Space: O(n).
func NewNumArray(nums []int) NumArray {
	prefix := make([]int, len(nums)+1) // one extra slot; prefix[0] stays 0
	for i, v := range nums {
		prefix[i+1] = prefix[i] + v // running cumulative sum
	}
	return NumArray{prefix: prefix}
}

// SumRange returns nums[left..right] via a single subtraction.
//
// Time:  O(1).
// Space: O(1).
func (a NumArray) SumRange(left, right int) int {
	// Everything up to right, minus everything before left, = the window.
	return a.prefix[right+1] - a.prefix[left]
}

func main() {
	// Official example:
	// NumArray([-2,0,3,-5,2,-1])
	// sumRange(0,2) -> 1 ; sumRange(2,5) -> -1 ; sumRange(0,5) -> -3
	nums := []int{-2, 0, 3, -5, 2, -1}

	fmt.Println("=== Approach 1: Brute Force ===")
	nb := NewNumArrayBrute(nums)
	fmt.Println(nb.SumRange(0, 2)) // expected 1
	fmt.Println(nb.SumRange(2, 5)) // expected -1
	fmt.Println(nb.SumRange(0, 5)) // expected -3

	fmt.Println("=== Approach 2: Prefix Sum (Optimal) ===")
	na := NewNumArray(nums)
	fmt.Println(na.SumRange(0, 2)) // expected 1
	fmt.Println(na.SumRange(2, 5)) // expected -1
	fmt.Println(na.SumRange(0, 5)) // expected -3
}
