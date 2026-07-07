package main

import (
	"fmt"
	"math"
)

// ── Approach 1: Brute Force (Three Nested Loops) ─────────────────────────────
//
// bruteForce solves 132 Pattern by trying every triple of indices i < j < k and
// checking the 132 value relation directly.
//
// Intuition:
//
//	The definition is a triple (i, j, k) with i < j < k and
//	nums[i] < nums[k] < nums[j]. The most literal thing we can do is enumerate
//	all such ordered triples and test the inequality. It is obviously correct
//	and serves as the baseline the smarter approaches beat.
//
// Algorithm:
//  1. For every i, for every j > i, for every k > j:
//  2. If nums[i] < nums[k] AND nums[k] < nums[j], a 132 pattern exists → true.
//  3. If no triple qualifies, return false.
//
// Time:  O(n^3) — three nested loops over the array.
// Space: O(1) — only loop counters.
func bruteForce(nums []int) bool {
	n := len(nums)
	for i := 0; i < n; i++ { // choose the "1" (smallest, leftmost)
		for j := i + 1; j < n; j++ { // choose the "3" (largest, middle index)
			for k := j + 1; k < n; k++ { // choose the "2" (middle value, rightmost)
				// 132 means nums[i] < nums[k] < nums[j].
				if nums[i] < nums[k] && nums[k] < nums[j] {
					return true // found one qualifying triple
				}
			}
		}
	}
	return false // exhausted all triples, none matched
}

// ── Approach 2: Fix j, Track Min-i So Far (Two Loops) ────────────────────────
//
// minPrefix solves 132 Pattern by fixing the middle element as "3" and reusing
// the running prefix minimum as the best possible "1".
//
// Intuition:
//
//	For a fixed j (the "3"), the best "1" is the smallest value strictly to its
//	left — call it minLeft. Then we only need some k > j with
//	minLeft < nums[k] < nums[j]. So sweep j left→right maintaining minLeft
//	(prefix minimum before j), and for each j scan k to the right looking for a
//	value in the open interval (minLeft, nums[j]). This removes the innermost i
//	loop by precomputing "the only i worth trying".
//
// Algorithm:
//  1. minLeft = nums[0]; treat it as the running minimum of nums[0..j-1].
//  2. For j from 1..n-1:
//     a. For k from j+1..n-1: if minLeft < nums[k] < nums[j], return true.
//     b. Update minLeft = min(minLeft, nums[j]) for the next j.
//  3. Return false.
//
// Time:  O(n^2) — for each j an inner scan over k.
// Space: O(1) — a single running minimum.
func minPrefix(nums []int) bool {
	n := len(nums)
	if n < 3 {
		return false // need at least three elements for a triple
	}
	minLeft := nums[0] // smallest value seen strictly left of j (best "1")
	for j := 1; j < n; j++ {
		// Look for a "2": some k > j whose value sits strictly between the best
		// "1" (minLeft) and the current "3" (nums[j]).
		for k := j + 1; k < n; k++ {
			if minLeft < nums[k] && nums[k] < nums[j] {
				return true // minLeft (i) < nums[k] (k) < nums[j] (j): a 132 pattern
			}
		}
		if nums[j] < minLeft {
			minLeft = nums[j] // extend the prefix minimum for future j's
		}
	}
	return false
}

// ── Approach 3: Monotonic Stack, Scan Right→Left (Optimal) ───────────────────
//
// monotonicStack solves 132 Pattern in one right-to-left pass using a
// decreasing stack of candidate "3" values and a running "2" ceiling.
//
// Intuition:
//
//	Scan from the right and treat each element as the "1". We want to know: has
//	some pair (j, k) already appeared to the right with nums[k] < nums[j]
//	(k > j) — i.e. a valid "2" whose matching "3" was even larger? Maintain the
//	largest such "2" seen so far, call it `third`. Any value we pop off a
//	decreasing stack becomes a "3" that had a strictly larger element to its
//	right serving as... no — precisely: the stack holds potential "3"s in
//	decreasing order; when the current value exceeds the stack top, that top is
//	a "2" (there was a bigger "3" further right), so we raise `third` to it.
//	Then if any later (leftward) element is strictly less than `third`, that
//	element is a valid "1" and we have i < j < k with nums[i] < nums[k] <
//	nums[j]. Because the current value is always a larger "3" sitting to the
//	LEFT of the popped "2", the k index (the popped "2") is to the right of j.
//
// Algorithm:
//  1. Initialise third = -inf (the best "2" value with a larger "3" to its right).
//  2. Use an empty stack (slice) of candidate "3" values.
//  3. Scan i from n-1 down to 0:
//     a. If nums[i] < third → a 132 pattern exists (nums[i] is the "1"). Return true.
//     b. While the stack is non-empty and nums[i] > stack top: pop it into
//     `third` (that popped value is a valid "2" paired with nums[i] as "3").
//     c. Push nums[i] as a new candidate "3".
//  4. Return false.
//
// Time:  O(n) — each element is pushed and popped at most once.
// Space: O(n) — the stack may hold up to n candidates.
func monotonicStack(nums []int) bool {
	n := len(nums)
	if n < 3 {
		return false // impossible to form a triple
	}
	third := math.MinInt64 // best "2" that already has a larger "3" to its right
	stack := []int{}       // decreasing stack of candidate "3" values (right side)
	for i := n - 1; i >= 0; i-- {
		// nums[i] plays the "1": if it is below the best known "2", we win,
		// because that "2" had a strictly larger "3" between i and it.
		if nums[i] < third {
			return true
		}
		// Everything smaller than nums[i] on the stack is a "2" that pairs with
		// nums[i] as its larger "3"; raise `third` to the largest such value.
		for len(stack) > 0 && nums[i] > stack[len(stack)-1] {
			third = stack[len(stack)-1] // this popped value is a valid "2"
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, nums[i]) // nums[i] is a fresh candidate "3"
	}
	return false
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Three Nested Loops) ===")
	fmt.Println(bruteForce([]int{1, 2, 3, 4}))  // expected false
	fmt.Println(bruteForce([]int{3, 1, 4, 2}))  // expected true
	fmt.Println(bruteForce([]int{-1, 3, 2, 0})) // expected true

	fmt.Println("=== Approach 2: Fix j, Track Min-i So Far (Two Loops) ===")
	fmt.Println(minPrefix([]int{1, 2, 3, 4}))  // expected false
	fmt.Println(minPrefix([]int{3, 1, 4, 2}))  // expected true
	fmt.Println(minPrefix([]int{-1, 3, 2, 0})) // expected true

	fmt.Println("=== Approach 3: Monotonic Stack, Scan Right→Left (Optimal) ===")
	fmt.Println(monotonicStack([]int{1, 2, 3, 4}))  // expected false
	fmt.Println(monotonicStack([]int{3, 1, 4, 2}))  // expected true
	fmt.Println(monotonicStack([]int{-1, 3, 2, 0})) // expected true
}
