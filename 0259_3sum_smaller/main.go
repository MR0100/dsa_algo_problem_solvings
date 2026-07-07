package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Brute Force (Triple Loop) ────────────────────────────────────
//
// bruteForce solves 3Sum Smaller by checking every triple i<j<k.
//
// Intuition:
//
//	The definition asks for the count of index triples with nums[i]+nums[j]+
//	nums[k] < target. Enumerate all of them with three nested loops.
//
// Algorithm:
//  1. For each i < j < k, if nums[i]+nums[j]+nums[k] < target, count++.
//  2. Return count.
//
// Time:  O(n^3) — three nested loops over n elements.
// Space: O(1) — just a counter.
func bruteForce(nums []int, target int) int {
	count := 0
	n := len(nums)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			for k := j + 1; k < n; k++ {
				if nums[i]+nums[j]+nums[k] < target { // valid triple
					count++
				}
			}
		}
	}
	return count
}

// ── Approach 2: Sort + Two Pointers (Optimal) ────────────────────────────────
//
// twoPointers solves 3Sum Smaller by sorting, then for each first element
// counting valid pairs in the remaining suffix with a two-pointer sweep.
//
// Intuition:
//
//	Sort the array. Fix the smallest element at index i. Now count pairs
//	(lo, hi) with i < lo < hi such that nums[i]+nums[lo]+nums[hi] < target.
//	With the suffix sorted, put lo=i+1 and hi=n-1. If the sum is < target,
//	then EVERY hi' in (lo, hi] also works (they are ≤ nums[hi]), so all
//	(hi - lo) pairs qualify at once — add them and advance lo. Otherwise the
//	sum is too big, so shrink hi.
//
// Algorithm:
//  1. Sort nums.
//  2. For i = 0..n-3: lo=i+1, hi=n-1.
//  3. While lo < hi:
//     - If nums[i]+nums[lo]+nums[hi] < target: count += hi-lo; lo++.
//     - Else hi--.
//  4. Return count.
//
// Time:  O(n^2) — sort O(n log n) + n outer iterations each with an O(n) sweep.
// Space: O(1) extra (in-place sort ignoring its stack).
func twoPointers(nums []int, target int) int {
	sort.Ints(nums) // sorting lets us count many pairs in one comparison
	count := 0
	n := len(nums)
	for i := 0; i < n-2; i++ { // fix the smallest of the triple
		lo, hi := i+1, n-1
		for lo < hi {
			if nums[i]+nums[lo]+nums[hi] < target {
				// nums[lo] with ANY hi' in (lo, hi] is also < target because
				// those values are ≤ nums[hi]; count all of them at once.
				count += hi - lo
				lo++ // move to the next (larger) middle element
			} else {
				hi-- // sum too large; the largest partner must shrink
			}
		}
	}
	return count
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Triple Loop) ===")
	fmt.Println(bruteForce([]int{-2, 0, 1, 3}, 2)) // expected 2
	fmt.Println(bruteForce([]int{}, 0))            // expected 0
	fmt.Println(bruteForce([]int{0}, 0))           // expected 0

	fmt.Println("=== Approach 2: Sort + Two Pointers (Optimal) ===")
	fmt.Println(twoPointers([]int{-2, 0, 1, 3}, 2)) // expected 2
	fmt.Println(twoPointers([]int{}, 0))            // expected 0
	fmt.Println(twoPointers([]int{0}, 0))           // expected 0
}
