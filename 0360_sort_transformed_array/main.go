package main

import (
	"fmt"
	"sort"
)

// apply evaluates the quadratic f(x) = a*x^2 + b*x + c at x.
func apply(x, a, b, c int) int {
	return a*x*x + b*x + c
}

// ── Approach 1: Brute Force (Transform + Sort) ───────────────────────────────
//
// bruteForce solves Sort Transformed Array by applying f to every element and
// then sorting the results.
//
// Intuition:
//
//	Ignore the fact that nums is already sorted. Just compute f(x) for each x
//	and hand the resulting slice to a comparison sort. Correct for any a, b, c,
//	and a solid baseline to validate the O(n) solution against.
//
// Algorithm:
//  1. For each x in nums, compute y = a*x*x + b*x + c.
//  2. Sort the y-values ascending.
//  3. Return them.
//
// Time:  O(n log n) — the sort dominates.
// Space: O(n) — the output slice.
func bruteForce(nums []int, a, b, c int) []int {
	result := make([]int, len(nums))
	for i, x := range nums {
		result[i] = apply(x, a, b, c) // transform each element
	}
	sort.Ints(result) // then sort the transformed values
	return result
}

// ── Approach 2: Two Pointers (Optimal) ───────────────────────────────────────
//
// twoPointers solves Sort Transformed Array in linear time by exploiting the
// shape of a parabola over a sorted input.
//
// Intuition:
//
//	f(x) = a*x^2 + b*x + c is a parabola. On a sorted array, the transformed
//	values are largest at the two ENDS when a > 0 (upward parabola: extremes
//	are far from the vertex) and largest in the MIDDLE when a < 0 (downward
//	parabola). So compare the two ends and pour the more extreme value into the
//	output — filling from the back when a > 0 (largest first) and from the front
//	when a < 0. When a == 0 the function is linear/monotone, so either direction
//	works; treat it like a >= 0 (fill largest at the back). Two pointers from
//	both ends converge in a single pass.
//
// Algorithm:
//  1. Compute y-values lazily as fa = f(nums[left]), fb = f(nums[right]).
//  2. If a >= 0: fill the result from the back with the larger of the two ends,
//     advancing that pointer. If a < 0: fill from the front with the smaller of
//     the two ends.
//  3. Continue until the pointers cross.
//
// Time:  O(n) — each element is placed exactly once.
// Space: O(n) — the output slice (O(1) auxiliary beyond it).
func twoPointers(nums []int, a, b, c int) []int {
	n := len(nums)
	result := make([]int, n)
	left, right := 0, n-1 // scan sorted nums from both ends

	if a >= 0 {
		// Upward parabola (or linear): extremes are at the ends → largest first.
		idx := n - 1 // fill position, from the back
		for left <= right {
			fl := apply(nums[left], a, b, c)  // value at the left end
			fr := apply(nums[right], a, b, c) // value at the right end
			if fl >= fr {
				result[idx] = fl // left end is the larger extreme
				left++
			} else {
				result[idx] = fr // right end is the larger extreme
				right--
			}
			idx-- // next slot toward the front
		}
	} else {
		// Downward parabola: extremes (smallest values) are at the ends → fill
		// smallest first from the front.
		idx := 0 // fill position, from the front
		for left <= right {
			fl := apply(nums[left], a, b, c)
			fr := apply(nums[right], a, b, c)
			if fl <= fr {
				result[idx] = fl // left end is the smaller value
				left++
			} else {
				result[idx] = fr // right end is the smaller value
				right--
			}
			idx++ // next slot toward the back
		}
	}
	return result
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Transform + Sort) ===")
	fmt.Printf("nums=[-4,-2,2,4], a=1,b=3,c=5   got=%v  expected [3 9 15 33]\n", bruteForce([]int{-4, -2, 2, 4}, 1, 3, 5))
	fmt.Printf("nums=[-4,-2,2,4], a=-1,b=3,c=5  got=%v  expected [-23 -5 1 7]\n", bruteForce([]int{-4, -2, 2, 4}, -1, 3, 5))
	fmt.Printf("nums=[-4,-2,2,4], a=0,b=1,c=0   got=%v  expected [-4 -2 2 4]\n", bruteForce([]int{-4, -2, 2, 4}, 0, 1, 0))

	fmt.Println("=== Approach 2: Two Pointers (Optimal) ===")
	fmt.Printf("nums=[-4,-2,2,4], a=1,b=3,c=5   got=%v  expected [3 9 15 33]\n", twoPointers([]int{-4, -2, 2, 4}, 1, 3, 5))
	fmt.Printf("nums=[-4,-2,2,4], a=-1,b=3,c=5  got=%v  expected [-23 -5 1 7]\n", twoPointers([]int{-4, -2, 2, 4}, -1, 3, 5))
	fmt.Printf("nums=[-4,-2,2,4], a=0,b=1,c=0   got=%v  expected [-4 -2 2 4]\n", twoPointers([]int{-4, -2, 2, 4}, 0, 1, 0))
}
