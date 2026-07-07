package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Sort then Swap Pairs (Brute Force) ───────────────────────────
//
// sortAndSwap solves Wiggle Sort by sorting ascending, then swapping adjacent
// pairs starting at index 1 to create the up/down pattern.
//
// Intuition:
//
//	After sorting, the array is fully ascending. Swapping (1,2), (3,4), … pushes
//	the larger of each pair into the odd "peak" positions, producing
//	nums[0] <= nums[1] >= nums[2] <= nums[3] ...
//
// Algorithm:
//  1. Sort nums ascending.
//  2. For i = 1; i < n-1; i += 2: swap nums[i] and nums[i+1].
//
// Time:  O(n log n) — dominated by the sort.
// Space: O(1) extra (in-place swaps; sort is in place).
func sortAndSwap(nums []int) {
	sort.Ints(nums) // ascending order first
	for i := 1; i < len(nums)-1; i += 2 {
		// swap each peak position with its successor to build the wiggle
		nums[i], nums[i+1] = nums[i+1], nums[i]
	}
}

// ── Approach 2: One-Pass Greedy Swap (Optimal) ───────────────────────────────
//
// greedy solves Wiggle Sort in a single linear scan by fixing each position as
// we go, without sorting.
//
// Intuition:
//
//	Walk left to right. At even index i we WANT nums[i] <= nums[i+1]; at odd
//	index i we WANT nums[i] >= nums[i+1]. Whenever the current pair violates the
//	desired relation, swap them. A local swap never breaks the already-satisfied
//	relation to the left, so one pass suffices.
//
// Algorithm:
//  1. For i = 0..n-2:
//     if i is even and nums[i] > nums[i+1]: swap.
//     if i is odd  and nums[i] < nums[i+1]: swap.
//
// Time:  O(n) — single pass.
// Space: O(1) — in place.
func greedy(nums []int) {
	for i := 0; i < len(nums)-1; i++ {
		if i%2 == 0 {
			// even index should be a "valley": nums[i] <= nums[i+1]
			if nums[i] > nums[i+1] {
				nums[i], nums[i+1] = nums[i+1], nums[i]
			}
		} else {
			// odd index should be a "peak": nums[i] >= nums[i+1]
			if nums[i] < nums[i+1] {
				nums[i], nums[i+1] = nums[i+1], nums[i]
			}
		}
	}
}

// clone copies a slice so each approach runs on fresh input.
func clone(nums []int) []int {
	c := make([]int, len(nums))
	copy(c, nums)
	return c
}

func main() {
	// Example 1: [3,5,2,1,6,4] → a valid wiggle such as [3,5,1,6,2,4].
	// (Multiple valid answers exist; both approaches produce a valid wiggle.)
	fmt.Println("=== Approach 1: Sort then Swap Pairs ===")
	a := clone([]int{3, 5, 2, 1, 6, 4})
	sortAndSwap(a)
	fmt.Println(a) // expected [1 3 2 5 4 6] (valid wiggle)

	// Example 2: [6,6,5,6,3,8] → a valid wiggle.
	b := clone([]int{6, 6, 5, 6, 3, 8})
	sortAndSwap(b)
	fmt.Println(b) // expected [3 6 5 6 6 8] (valid wiggle)

	fmt.Println("=== Approach 2: One-Pass Greedy Swap (Optimal) ===")
	c := clone([]int{3, 5, 2, 1, 6, 4})
	greedy(c)
	fmt.Println(c) // expected [3 5 1 6 2 4] (valid wiggle)

	d := clone([]int{6, 6, 5, 6, 3, 8})
	greedy(d)
	fmt.Println(d) // expected [6 6 5 6 3 8] (valid wiggle)
}
