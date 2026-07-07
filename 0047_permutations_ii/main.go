package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Backtracking + Set Dedup ─────────────────────────────────────
//
// bruteForce solves Permutations II by generating all permutations and
// deduplicating with a map. Simple but uses extra space.
//
// Time:  O(n! * n) — generating all permutations + string key for dedup
// Space: O(n! * n)
func bruteForce(nums []int) [][]int {
	seen := map[string]bool{}
	var result [][]int
	visited := make([]bool, len(nums))
	var bt func(path []int)
	bt = func(path []int) {
		if len(path) == len(nums) {
			key := fmt.Sprint(path)
			if !seen[key] {
				seen[key] = true
				tmp := make([]int, len(nums))
				copy(tmp, path)
				result = append(result, tmp)
			}
			return
		}
		for i := 0; i < len(nums); i++ {
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

// ── Approach 2: Backtracking with Skip-Duplicate Pruning (Optimal) ────────────
//
// backtracking solves Permutations II by sorting first and skipping duplicate
// values at the same recursion level — no extra dedup structures needed.
//
// Intuition: Sort nums. At each level, if we pick nums[i] and nums[i] == nums[i-1]
// AND nums[i-1] was NOT visited in the current path, we'd generate a duplicate.
// The condition `!visited[i-1]` means: nums[i-1] has already been "put back"
// (was chosen and then un-chosen), so choosing nums[i] (same value) at this
// level would repeat the same subtree.
//
// Algorithm:
//  sort(nums)
//  bt(path, visited):
//    if len(path)==n: record; return
//    for i=0 to n-1:
//      if visited[i]: continue
//      if i>0 and nums[i]==nums[i-1] and NOT visited[i-1]: continue  // skip dup
//      visited[i]=true; bt(path+[nums[i]]); visited[i]=false
//
// Time:  O(n! * n) worst case; significantly better with duplicates
// Space: O(n) — visited array + recursion stack
func backtracking(nums []int) [][]int {
	sort.Ints(nums)
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
			if visited[i] {
				continue
			}
			// skip if same value as previous and previous was NOT in current path
			// (meaning it was already explored and backtracked at this level)
			if i > 0 && nums[i] == nums[i-1] && !visited[i-1] {
				continue
			}
			visited[i] = true
			bt(append(path, nums[i]))
			visited[i] = false
		}
	}
	bt(nil)
	return result
}

func main() {
	fmt.Println("=== Approach 1: Backtracking + Map Dedup ===")
	r1 := bruteForce([]int{1, 1, 2})
	fmt.Printf("nums=[1,1,2]  count=%d  expected 3  %v\n", len(r1), r1)
	r2 := bruteForce([]int{1, 2, 3})
	fmt.Printf("nums=[1,2,3]  count=%d  expected 6  %v\n", len(r2), r2)

	fmt.Println("\n=== Approach 2: Backtracking Skip-Dup (Optimal) ===")
	r3 := backtracking([]int{1, 1, 2})
	fmt.Printf("nums=[1,1,2]  count=%d  expected 3  %v\n", len(r3), r3)
	r4 := backtracking([]int{1, 2, 3})
	fmt.Printf("nums=[1,2,3]  count=%d  expected 6  %v\n", len(r4), r4)
	r5 := backtracking([]int{1, 1, 1, 2})
	fmt.Printf("nums=[1,1,1,2] count=%d  expected 4  %v\n", len(r5), r5)
}
