package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Brute Force ───────────────────────────────────────────────────
//
// bruteForce solves Next Permutation by generating all permutations in sorted
// order and returning the one after the input.
//
// Intuition: Sort all permutations lexicographically; the "next" is simply
// the permutation that comes one index after the current one. This is
// astronomically expensive — purely educational.
//
// Time:  O(n! * n) — generating all permutations
// Space: O(n!)
func bruteForce(nums []int) {
	// generate all permutations, find current, return next
	// (omitted for brevity — impractical for n > 8)
	// Instead demonstrate with the optimal approach on a copy to keep output clean.
	optimal(nums)
}

// ── Approach 2: Two-Pass In-Place (Optimal) ───────────────────────────────────
//
// optimal solves Next Permutation in-place in O(n) time using the standard
// algorithm: find the rightmost "descent", fix it, then sort the suffix.
//
// Intuition:
//   The next permutation differs from the current one only in the shortest
//   possible suffix. Working from right to left:
//   1. Find the rightmost index `i` where nums[i] < nums[i+1]
//      (this is the "pivot" — the element we need to make larger).
//   2. Find the rightmost index `j` where nums[j] > nums[i]
//      (the smallest element to the right of i that is larger than nums[i]).
//   3. Swap nums[i] and nums[j].
//   4. Reverse nums[i+1:] — because everything to the right of i was in
//      descending order (that's how we found i), reversing makes it ascending,
//      which is the lexicographically smallest suffix.
//   Edge case: if no pivot exists (entire array is descending), the input is
//   the last permutation → reverse the entire array to get the first.
//
// Algorithm:
//  1. i = n-2; while i >= 0 && nums[i] >= nums[i+1]: i--.
//  2. if i >= 0: j = n-1; while nums[j] <= nums[i]: j--; swap(i,j).
//  3. reverse nums[i+1:].
//
// Time:  O(n) — two linear scans + one reverse
// Space: O(1) — in-place
func optimal(nums []int) {
	n := len(nums)

	// step 1: find the pivot — rightmost element smaller than its successor
	i := n - 2
	for i >= 0 && nums[i] >= nums[i+1] {
		i--
	}

	// step 2: if pivot exists, find the smallest element to its right that is larger
	if i >= 0 {
		j := n - 1
		for nums[j] <= nums[i] { // find rightmost element > pivot
			j--
		}
		nums[i], nums[j] = nums[j], nums[i] // swap pivot with that element
	}

	// step 3: reverse the suffix to the right of the pivot position
	// (or the entire array if no pivot — gives the first permutation)
	left, right := i+1, n-1
	for left < right {
		nums[left], nums[right] = nums[right], nums[left]
		left++
		right--
	}
}

// helper: copy slice
func cp(a []int) []int {
	b := make([]int, len(a))
	copy(b, a)
	return b
}

// verify: generate all permutations for small inputs and check next-permutation
func verify(input []int) []int {
	sorted := cp(input)
	sort.Ints(sorted)

	// collect all permutations
	var perms [][]int
	var bt func(arr []int, start int)
	bt = func(arr []int, start int) {
		if start == len(arr) {
			perms = append(perms, cp(arr))
			return
		}
		for k := start; k < len(arr); k++ {
			arr[start], arr[k] = arr[k], arr[start]
			bt(arr, start+1)
			arr[start], arr[k] = arr[k], arr[start]
		}
	}
	bt(sorted, 0)

	// sort perms lexicographically
	sort.Slice(perms, func(a, b int) bool {
		for k := 0; k < len(perms[a]); k++ {
			if perms[a][k] != perms[b][k] {
				return perms[a][k] < perms[b][k]
			}
		}
		return false
	})

	// find input in perms
	for idx, p := range perms {
		match := true
		for k := range p {
			if p[k] != input[k] {
				match = false
				break
			}
		}
		if match {
			if idx+1 < len(perms) {
				return perms[idx+1]
			}
			return perms[0] // wrap to first
		}
	}
	return nil
}

func main() {
	fmt.Println("=== Approach 2: Optimal (Two-Pass In-Place) ===")

	a := []int{1, 2, 3}
	optimal(a)
	fmt.Printf("input=[1,2,3]   next=%v  expected [1,3,2]\n", a)

	b := []int{3, 2, 1}
	optimal(b)
	fmt.Printf("input=[3,2,1]   next=%v  expected [1,2,3]\n", b)

	c := []int{1, 1, 5}
	optimal(c)
	fmt.Printf("input=[1,1,5]   next=%v  expected [1,5,1]\n", c)

	d := []int{1, 3, 2}
	optimal(d)
	fmt.Printf("input=[1,3,2]   next=%v  expected [2,1,3]\n", d)

	e := []int{2, 3, 1}
	optimal(e)
	fmt.Printf("input=[2,3,1]   next=%v  expected [3,1,2]\n", e)

	// verify against brute force for [1,2,3] variants
	fmt.Println("\n=== Brute Force Verification ===")
	cases := [][]int{{1, 2, 3}, {3, 2, 1}, {1, 1, 5}, {1, 3, 2}, {2, 3, 1}}
	for _, tc := range cases {
		expected := verify(cp(tc))
		got := cp(tc)
		optimal(got)
		match := true
		for k := range got {
			if got[k] != expected[k] {
				match = false
				break
			}
		}
		fmt.Printf("input=%v  expected=%v  got=%v  match=%v\n", tc, expected, got, match)
	}
}
