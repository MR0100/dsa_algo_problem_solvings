package main

import "fmt"

// ── Approach 1: Backtracking with Visited Array ───────────────────────────────
//
// backtracking solves Permutations by building each permutation position by
// position, using a boolean visited array to track which elements are used.
//
// Intuition: At each recursion level, choose one of the unused elements to fill
// the current position. When the path reaches length n, record it.
//
// Algorithm:
//  bt(path, visited):
//    if len(path) == n: record path; return
//    for i = 0 to n-1:
//      if not visited[i]:
//        visited[i]=true; bt(path+[nums[i]], visited); visited[i]=false
//
// Time:  O(n! * n) — n! permutations, each costs O(n) to copy
// Space: O(n) — recursion depth n; visited array n
func backtracking(nums []int) [][]int {
	n := len(nums)
	var result [][]int
	visited := make([]bool, n)

	var bt func(path []int)
	bt = func(path []int) {
		if len(path) == n {
			tmp := make([]int, n)
			copy(tmp, path)
			result = append(result, tmp)
			return
		}
		for i := 0; i < n; i++ {
			if !visited[i] {
				visited[i] = true
				bt(append(path, nums[i]))
				visited[i] = false
			}
		}
	}
	bt(nil)
	return result
}

// ── Approach 2: Swap-Based Backtracking (In-Place) ───────────────────────────
//
// swapBacktracking solves Permutations by swapping elements in-place.
//
// Intuition: Fix elements in the slice at each position by swapping nums[start]
// with each nums[i] (i >= start), recursing on the suffix, then swapping back.
// Avoids the visited array; generates all permutations of nums[:n].
//
// Time:  O(n! * n)
// Space: O(n) — recursion depth only (no visited array)
func swapBacktracking(nums []int) [][]int {
	var result [][]int
	n := len(nums)

	var bt func(start int)
	bt = func(start int) {
		if start == n {
			tmp := make([]int, n)
			copy(tmp, nums)
			result = append(result, tmp)
			return
		}
		for i := start; i < n; i++ {
			nums[start], nums[i] = nums[i], nums[start] // choose nums[i] for position start
			bt(start + 1)
			nums[start], nums[i] = nums[i], nums[start] // restore
		}
	}
	bt(0)
	return result
}

// ── Approach 3: Iterative (Insert Each Element into All Positions) ────────────
//
// iterative solves Permutations by iteratively inserting each number into every
// possible position of existing permutations.
//
// Intuition: Start with [[]], which represents the only permutation of an empty
// set. For each new number, take every existing permutation and insert the new
// number at every possible position (0 to len(perm)).
//
// Time:  O(n! * n)
// Space: O(n! * n)
func iterative(nums []int) [][]int {
	result := [][]int{{}} // start with one empty permutation
	for _, num := range nums {
		var next [][]int
		for _, perm := range result {
			// insert num at every position in perm
			for pos := 0; pos <= len(perm); pos++ {
				newPerm := make([]int, len(perm)+1)
				copy(newPerm[:pos], perm[:pos])
				newPerm[pos] = num
				copy(newPerm[pos+1:], perm[pos:])
				next = append(next, newPerm)
			}
		}
		result = next
	}
	return result
}

func main() {
	fmt.Println("=== Approach 1: Backtracking with Visited Array ===")
	r1 := backtracking([]int{1, 2, 3})
	fmt.Printf("nums=[1,2,3]  count=%d  expected 6\n%v\n", len(r1), r1)
	r2 := backtracking([]int{0, 1})
	fmt.Printf("nums=[0,1]    count=%d  expected 2  %v\n", len(r2), r2)

	fmt.Println("\n=== Approach 2: Swap Backtracking (In-Place) ===")
	r3 := swapBacktracking([]int{1, 2, 3})
	fmt.Printf("nums=[1,2,3]  count=%d  expected 6\n%v\n", len(r3), r3)

	fmt.Println("\n=== Approach 3: Iterative Insertion ===")
	r4 := iterative([]int{1, 2, 3})
	fmt.Printf("nums=[1,2,3]  count=%d  expected 6\n%v\n", len(r4), r4)
}
