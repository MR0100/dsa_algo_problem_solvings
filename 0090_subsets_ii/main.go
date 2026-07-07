package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Backtracking with Skip-Duplicate Guard ────────────────────────
//
// subsetsWithDup solves Subsets II (with duplicates) using backtracking.
//
// Intuition:
//   Sort the array first so duplicates are adjacent. During backtracking,
//   if nums[i] == nums[i-1] and i > start (i.e., nums[i-1] was not used as
//   the starting element of this recursion level), skip nums[i] to avoid
//   generating duplicate subsets.
//
//   The guard condition is: i > start (not i > 0). This ensures we skip
//   duplicates at the same recursion level but still allow the first element
//   at each level to be used.
//
// Time:  O(2^n × n) — 2^n subsets in the worst case (all distinct), each O(n).
// Space: O(n) — recursion depth.
func subsetsWithDup(nums []int) [][]int {
	sort.Ints(nums)
	var result [][]int
	var bt func(start int, path []int)
	bt = func(start int, path []int) {
		tmp := make([]int, len(path))
		copy(tmp, path)
		result = append(result, tmp)
		for i := start; i < len(nums); i++ {
			// skip duplicate at the same recursion level
			if i > start && nums[i] == nums[i-1] {
				continue
			}
			bt(i+1, append(path, nums[i]))
		}
	}
	bt(0, nil)
	return result
}

// ── Approach 2: Cascading with Duplicate Detection ────────────────────────────
//
// subsetsWithDupCascading solves Subsets II using iterative cascading,
// skipping duplicates by only appending new element to subsets from the
// current element's group.
//
// Intuition:
//   For distinct elements, cascading doubles the result. For duplicates, only
//   the subsets added in the PREVIOUS step (not all existing subsets) should
//   be extended — otherwise we'd create duplicates.
//
//   Track `startIdx` = index of first subset added in the last iteration.
//   If current element == previous: only extend subsets from startIdx onward.
//   If current element != previous: extend all subsets.
//
// Time:  O(2^n × n)
// Space: O(2^n × n)
func subsetsWithDupCascading(nums []int) [][]int {
	sort.Ints(nums)
	result := [][]int{{}}
	startIdx := 0 // index of first subset added in the previous round

	for i, num := range nums {
		start := 0
		if i > 0 && nums[i] == nums[i-1] {
			start = startIdx // only extend subsets from the last batch
		}
		startIdx = len(result) // record where new subsets will start
		for j := start; j < startIdx; j++ {
			newSub := make([]int, len(result[j])+1)
			copy(newSub, result[j])
			newSub[len(result[j])] = num
			result = append(result, newSub)
		}
	}
	return result
}

func main() {
	fmt.Println("=== Approach 1: Backtracking with Skip-Dup Guard ===")
	r1 := subsetsWithDup([]int{1, 2, 2})
	fmt.Printf("nums=[1,2,2]  count=%d  subsets=%v  expected 6\n", len(r1), r1)

	r2 := subsetsWithDup([]int{0})
	fmt.Printf("nums=[0]  count=%d  subsets=%v  expected 2\n", len(r2), r2)

	r3 := subsetsWithDup([]int{1, 1, 2, 2})
	fmt.Printf("nums=[1,1,2,2]  count=%d  expected 9\n", len(r3))

	fmt.Println("=== Approach 2: Cascading with Dup Detection ===")
	r4 := subsetsWithDupCascading([]int{1, 2, 2})
	fmt.Printf("nums=[1,2,2]  count=%d  expected 6\n", len(r4))

	r5 := subsetsWithDupCascading([]int{1, 1, 2, 2})
	fmt.Printf("nums=[1,1,2,2]  count=%d  expected 9\n", len(r5))
}
