package main

import "fmt"

// ── Approach 1: Naive Per-Element Updates (Brute Force) ───────────────────────
//
// bruteForce solves Range Addition by literally applying each update to every
// index in its [start, end] range.
//
// Intuition:
//
//	The statement says "increment arr[start..end] by inc". Do exactly that: for
//	each update, loop over its range and add inc to each cell. Correct but slow
//	— a single wide update can touch the whole array.
//
// Algorithm:
//  1. arr = zeros of size length.
//  2. For each update [start, end, inc]: for i in start..end, arr[i] += inc.
//  3. Return arr.
//
// Time:  O(length * k) worst case — k updates each spanning the full array.
// Space: O(length) for the output (O(1) auxiliary).
func bruteForce(length int, updates [][]int) []int {
	arr := make([]int, length) // starts all zeros
	for _, u := range updates {
		start, end, inc := u[0], u[1], u[2]
		for i := start; i <= end; i++ { // touch every index in the range
			arr[i] += inc
		}
	}
	return arr
}

// ── Approach 2: Difference Array (Optimal) ───────────────────────────────────
//
// differenceArray solves Range Addition with the difference-array trick:
// record each range update as two O(1) endpoint edits, then take a prefix sum.
//
// Intuition:
//
//	Adding inc to the whole range [start, end] shows up in the DIFFERENCE array
//	as exactly two changes: diff[start] += inc (the value steps up here) and
//	diff[end+1] -= inc (it steps back down just past the range). After all
//	updates are recorded this way, a running prefix sum of diff reconstructs the
//	real array — every element between start and end inherits the +inc, and
//	nothing outside does. Each update is O(1); one final O(length) sweep builds
//	the answer.
//
// Algorithm:
//  1. diff = zeros of size length (+1 conceptually; guard end+1 < length).
//  2. For each update [start, end, inc]: diff[start] += inc; if end+1 < length,
//     diff[end+1] -= inc.
//  3. Prefix-sum diff in place: arr[i] = arr[i-1] + diff[i].
//  4. Return diff (now holding the answer).
//
// Time:  O(length + k) — O(1) per update plus one prefix-sum pass.
// Space: O(length) — the difference/answer array.
func differenceArray(length int, updates [][]int) []int {
	diff := make([]int, length) // difference array, reused as the output
	for _, u := range updates {
		start, end, inc := u[0], u[1], u[2]
		diff[start] += inc // value steps up at the range start
		if end+1 < length {
			diff[end+1] -= inc // and steps back down just past the range end
		}
	}
	// Prefix sum turns the difference array back into actual values.
	for i := 1; i < length; i++ {
		diff[i] += diff[i-1]
	}
	return diff
}

func main() {
	// Example 1: length = 5,  updates = [[1,3,2],[2,4,3],[0,2,-2]]
	//            → [-2,0,3,5,3]
	// Example 2: length = 10, updates = [[2,4,6],[5,6,8],[1,9,-4]]
	//            → [0,-4,2,2,2,4,4,-4,-4,-4]

	ex1 := [][]int{{1, 3, 2}, {2, 4, 3}, {0, 2, -2}}
	ex2 := [][]int{{2, 4, 6}, {5, 6, 8}, {1, 9, -4}}

	fmt.Println("=== Approach 1: Naive Per-Element Updates (Brute Force) ===")
	fmt.Println(bruteForce(5, ex1))  // expected [-2 0 3 5 3]
	fmt.Println(bruteForce(10, ex2)) // expected [0 -4 2 2 2 4 4 -4 -4 -4]

	fmt.Println("=== Approach 2: Difference Array (Optimal) ===")
	fmt.Println(differenceArray(5, ex1))  // expected [-2 0 3 5 3]
	fmt.Println(differenceArray(10, ex2)) // expected [0 -4 2 2 2 4 4 -4 -4 -4]
}
