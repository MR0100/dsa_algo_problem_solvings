package main

import "fmt"

// ── Approach 1: Brute Force (Four Nested Loops) ──────────────────────────────
//
// bruteForce solves 4Sum II by trying every combination of one index from each
// of the four arrays and counting the quadruples that sum to zero.
//
// Intuition:
//
//	The definition is literally "count tuples (i,j,k,l) with
//	nums1[i]+nums2[j]+nums3[k]+nums4[l] == 0". So enumerate all of them. With
//	n up to 200 this is 200^4 = 1.6e9 iterations — too slow to submit, but it
//	is the ground truth the fast solution must match.
//
// Algorithm:
//  1. count = 0.
//  2. For every (a in nums1)(b in nums2)(c in nums3)(d in nums4):
//     if a+b+c+d == 0, count++.
//  3. Return count.
//
// Time:  O(n^4).
// Space: O(1).
func bruteForce(nums1, nums2, nums3, nums4 []int) int {
	count := 0
	for _, a := range nums1 {
		for _, b := range nums2 {
			for _, c := range nums3 {
				for _, d := range nums4 {
					if a+b+c+d == 0 { // this quadruple sums to zero
						count++ // every index tuple is counted separately
					}
				}
			}
		}
	}
	return count
}

// ── Approach 2: Two Hash Maps, Split 2+2 (Meet in the Middle, Optimal) ────────
//
// meetInTheMiddle solves 4Sum II by splitting the four arrays into two halves,
// hashing all pairwise sums of the first half, then for each pairwise sum of
// the second half looking up its negation.
//
// Intuition:
//
//	a + b + c + d == 0  ⇔  (a + b) == -(c + d). So precompute every possible
//	(a+b) with a from nums1, b from nums2 and store how many index pairs
//	produce each sum in a hash map. Then iterate every (c+d) with c from
//	nums3, d from nums4 and add the stored count of -(c+d): each such stored
//	pair combines with the current (c,d) into a valid zero-sum quadruple. This
//	turns O(n^4) into two O(n^2) passes — the classic meet-in-the-middle.
//
// Algorithm:
//  1. Build sumAB: for a in nums1, b in nums2, sumAB[a+b]++.
//  2. count = 0.
//  3. For c in nums3, d in nums4: count += sumAB[-(c+d)] (0 if absent).
//  4. Return count.
//
// Time:  O(n^2) — two double loops; map ops are O(1) average.
// Space: O(n^2) — the map holds up to n^2 distinct pair sums.
func meetInTheMiddle(nums1, nums2, nums3, nums4 []int) int {
	// sumAB[s] = number of (i, j) index pairs with nums1[i] + nums2[j] == s.
	sumAB := make(map[int]int, len(nums1)*len(nums2))
	for _, a := range nums1 {
		for _, b := range nums2 {
			sumAB[a+b]++ // record one more pair achieving this sum
		}
	}

	count := 0
	for _, c := range nums3 {
		for _, d := range nums4 {
			// We need a+b == -(c+d) to reach a total of 0. Every stored pair
			// with that sum pairs with the current (c,d) to form one quadruple.
			count += sumAB[-(c + d)] // missing key yields 0, which is correct
		}
	}
	return count
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Four Nested Loops) ===")
	fmt.Printf("[1,2],[-2,-1],[-1,2],[0,2] -> %d  expected 2\n",
		bruteForce([]int{1, 2}, []int{-2, -1}, []int{-1, 2}, []int{0, 2}))
	fmt.Printf("[0],[0],[0],[0]            -> %d  expected 1\n",
		bruteForce([]int{0}, []int{0}, []int{0}, []int{0}))

	fmt.Println("=== Approach 2: Two Hash Maps, Split 2+2 (Meet in the Middle, Optimal) ===")
	fmt.Printf("[1,2],[-2,-1],[-1,2],[0,2] -> %d  expected 2\n",
		meetInTheMiddle([]int{1, 2}, []int{-2, -1}, []int{-1, 2}, []int{0, 2}))
	fmt.Printf("[0],[0],[0],[0]            -> %d  expected 1\n",
		meetInTheMiddle([]int{0}, []int{0}, []int{0}, []int{0}))

	fmt.Println("=== Extra checks ===")
	// No quadruple sums to zero.
	fmt.Printf("[1],[1],[1],[1]            -> %d  expected 0\n",
		meetInTheMiddle([]int{1}, []int{1}, []int{1}, []int{1}))
	// All zeros, n=2 → every one of 2^4 = 16 tuples works.
	fmt.Printf("[0,0],[0,0],[0,0],[0,0]    -> %d  expected 16\n",
		meetInTheMiddle([]int{0, 0}, []int{0, 0}, []int{0, 0}, []int{0, 0}))
}
