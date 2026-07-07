package main

import "fmt"

// ── Approach 1: Backtracking ──────────────────────────────────────────────────
//
// backtracking solves Subsets by building all 2^n subsets via DFS.
//
// Intuition:
//   At each index, decide: include this element or skip it.
//   Record the current path at every level (not just at leaves).
//
// Algorithm:
//   bt(start, path):
//     record path
//     for i=start to n-1:
//       bt(i+1, path+[nums[i]])
//
// Time:  O(2^n × n) — 2^n subsets, each copied in O(n).
// Space: O(n) — recursion depth n.
func backtracking(nums []int) [][]int {
	var result [][]int
	var bt func(start int, path []int)
	bt = func(start int, path []int) {
		tmp := make([]int, len(path))
		copy(tmp, path)
		result = append(result, tmp)
		for i := start; i < len(nums); i++ {
			bt(i+1, append(path, nums[i]))
		}
	}
	bt(0, nil)
	return result
}

// ── Approach 2: Bit Manipulation ──────────────────────────────────────────────
//
// bitManipulation solves Subsets using bitmask enumeration.
//
// Intuition:
//   For n elements, there are 2^n subsets, each corresponding to a bitmask
//   from 0 to 2^n - 1. Bit j set in mask i means nums[j] is in subset i.
//
// Time:  O(2^n × n) — 2^n masks, each requiring O(n) to build the subset.
// Space: O(n) — current subset.
func bitManipulation(nums []int) [][]int {
	n := len(nums)
	total := 1 << n // 2^n
	result := make([][]int, 0, total)

	for mask := 0; mask < total; mask++ {
		subset := []int{}
		for j := 0; j < n; j++ {
			if mask&(1<<j) != 0 {
				subset = append(subset, nums[j])
			}
		}
		result = append(result, subset)
	}
	return result
}

// ── Approach 3: Cascading (Iterative) ────────────────────────────────────────
//
// cascading solves Subsets by starting with [[]] and adding each element to
// all existing subsets.
//
// Intuition:
//   Start: result = [[]]
//   For each num, create new subsets by adding num to every existing subset.
//   Append all new subsets to result.
//
// Time:  O(2^n × n) — after processing all n elements, 2^n subsets exist.
// Space: O(2^n × n) — output.
func cascading(nums []int) [][]int {
	result := [][]int{{}} // start with empty subset
	for _, num := range nums {
		n := len(result)
		for i := 0; i < n; i++ {
			// create a new subset by adding num to existing subset
			newSub := make([]int, len(result[i])+1)
			copy(newSub, result[i])
			newSub[len(result[i])] = num
			result = append(result, newSub)
		}
	}
	return result
}

func main() {
	nums1 := []int{1, 2, 3}
	nums2 := []int{0}

	fmt.Println("=== Approach 1: Backtracking ===")
	r1 := backtracking(nums1)
	fmt.Printf("nums=%v  count=%d  expected 8\n", nums1, len(r1))
	fmt.Println("subsets:", r1)

	r2 := backtracking(nums2)
	fmt.Printf("nums=%v  count=%d  expected 2\n", nums2, len(r2))
	fmt.Println("subsets:", r2)

	fmt.Println("=== Approach 2: Bit Manipulation ===")
	r3 := bitManipulation(nums1)
	fmt.Printf("nums=%v  count=%d  expected 8\n", nums1, len(r3))

	fmt.Println("=== Approach 3: Cascading ===")
	r4 := cascading(nums1)
	fmt.Printf("nums=%v  count=%d  expected 8\n", nums1, len(r4))
	fmt.Println("subsets:", r4)
}
