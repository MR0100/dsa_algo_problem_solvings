package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Brute Force Simulation ───────────────────────────────────────
//
// bruteForce solves Minimum Moves to Equal Array Elements by literally
// performing the described operation until every element is equal.
//
// Intuition:
//
//	Each move increments n-1 elements by 1, i.e. every element EXCEPT the
//	current maximum. Equivalently, the maximum stays put while everyone else
//	catches up by 1. Repeat until the min equals the max. This mirrors the
//	problem statement exactly and is useful to build confidence in the O(1)
//	formula — but it runs in time proportional to the answer, which can be
//	up to ~10^9, so it is a teaching baseline only.
//
// Algorithm:
//  1. moves = 0.
//  2. Loop: find min and max of nums.
//     - if min == max, stop.
//     - else add 1 to every element that is NOT the (single) maximum, moves++.
//  3. Return moves.
//
// Time:  O(moves · n) — moves can be ~10^9; TLE on large gaps. Baseline only.
// Space: O(1) — in-place increments.
func bruteForce(nums []int) int {
	moves := 0
	for {
		// Find current min and the index of a maximum.
		minVal, maxVal, maxIdx := nums[0], nums[0], 0
		for i, v := range nums {
			if v < minVal {
				minVal = v
			}
			if v > maxVal {
				maxVal = v
				maxIdx = i
			}
		}
		if minVal == maxVal {
			return moves // all equal → done
		}
		// One move: increment every element except one chosen maximum.
		for i := range nums {
			if i != maxIdx {
				nums[i]++ // n-1 elements go up by 1
			}
		}
		moves++
	}
}

// ── Approach 2: Math — Increment ≡ Decrement (Optimal) ───────────────────────
//
// mathDecrement solves the problem in one pass using the key reframing:
// "add 1 to n-1 elements" changes nothing about the *relative* differences
// versus "subtract 1 from 1 element".
//
// Intuition:
//
//	Adding 1 to n-1 elements is the same, relatively, as subtracting 1 from
//	the one element left out — the gaps between elements are all that matter,
//	and both operations shrink exactly one gap by 1. So instead of raising
//	everyone up to some target, imagine lowering everyone DOWN to the current
//	minimum, one unit at a time. Element nums[i] needs (nums[i] - min)
//	decrements, and each decrement is one move. Total moves = Σ(nums[i] - min)
//	= sum(nums) - n * min(nums).
//
// Algorithm:
//  1. Track sum and min in a single pass.
//  2. Return sum - n * min.
//
// Time:  O(n) — one pass.
// Space: O(1).
func mathDecrement(nums []int) int {
	sum := 0
	minVal := nums[0]
	for _, v := range nums {
		sum += v // running total of all elements
		if v < minVal {
			minVal = v // smallest element = the level everyone descends to
		}
	}
	// Each element must drop to minVal; the drops summed = total moves.
	return sum - len(nums)*minVal
}

// ── Approach 3: Sort + Sum of Differences From Min ───────────────────────────
//
// sortAndSum solves the same problem after sorting; once sorted, nums[0] is
// the minimum and the answer is the sum of every element's distance to it.
//
// Intuition:
//
//	Identical math to Approach 2, just derived after sorting so the minimum is
//	simply nums[0]. Summing nums[i] - nums[0] over i is exactly
//	sum - n*min. Sorting adds a log factor and no benefit here, but it makes
//	the "everyone descends to the smallest" picture explicit and is a natural
//	first instinct many people have.
//
// Algorithm:
//  1. Sort nums ascending (nums[0] becomes the minimum).
//  2. moves = Σ_{i} (nums[i] - nums[0]).
//  3. Return moves.
//
// Time:  O(n log n) — the sort dominates.
// Space: O(n) — a copy is made so the caller's slice ordering is preserved.
func sortAndSum(nums []int) int {
	// Work on a copy so we don't disturb the caller's slice ordering.
	cp := make([]int, len(nums))
	copy(cp, nums)
	sort.Ints(cp) // ascending; cp[0] is now the minimum element

	moves := 0
	for _, v := range cp {
		moves += v - cp[0] // distance of each element down to the smallest
	}
	return moves
}

func main() {
	fmt.Println("=== Approach 1: Brute Force Simulation ===")
	// bruteForce mutates its slice, so pass fresh copies each call.
	fmt.Printf("nums=[1,2,3] -> %d  expected 3\n", bruteForce([]int{1, 2, 3}))
	fmt.Printf("nums=[1,1,1] -> %d  expected 0\n", bruteForce([]int{1, 1, 1}))

	fmt.Println("=== Approach 2: Math — Increment ≡ Decrement (Optimal) ===")
	fmt.Printf("nums=[1,2,3] -> %d  expected 3\n", mathDecrement([]int{1, 2, 3}))
	fmt.Printf("nums=[1,1,1] -> %d  expected 0\n", mathDecrement([]int{1, 1, 1}))

	fmt.Println("=== Approach 3: Sort + Sum of Differences From Min ===")
	fmt.Printf("nums=[1,2,3] -> %d  expected 3\n", sortAndSum([]int{1, 2, 3}))
	fmt.Printf("nums=[1,1,1] -> %d  expected 0\n", sortAndSum([]int{1, 1, 1}))

	fmt.Println("=== Extra checks ===")
	fmt.Printf("nums=[1,10]      -> %d  expected 9\n", mathDecrement([]int{1, 10}))         // one gap of 9
	fmt.Printf("nums=[5,6,8,8,5] -> %d  expected 7\n", mathDecrement([]int{5, 6, 8, 8, 5})) // sum 32 - 5*5 = 7
	fmt.Printf("nums=[5,6,8,8,5] -> %d  expected 7 (sort)\n", sortAndSum([]int{5, 6, 8, 8, 5}))
}
