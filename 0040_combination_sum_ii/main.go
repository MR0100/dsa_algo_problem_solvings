package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Brute Force (Backtracking + Set Dedup) ───────────────────────
//
// bruteForce solves Combination Sum II using backtracking but collecting results
// in a map to deduplicate. Simple but uses extra space.
//
// Intuition: Same backtracking as #39 but each element can only be used once
// (advance i+1 instead of i). Duplicates in candidates produce duplicate
// combinations; deduplicate by storing sorted combo keys in a map.
//
// Time:  O(2^N * N) — up to 2^N subsets, each takes O(N) to encode
// Space: O(2^N * N) — for the map
func bruteForce(candidates []int, target int) [][]int {
	sort.Ints(candidates)
	seen := map[string]bool{}
	var result [][]int
	var bt func(start, remaining int, path []int)
	bt = func(start, remaining int, path []int) {
		if remaining == 0 {
			// encode the combo as a string key for dedup
			key := fmt.Sprint(path)
			if !seen[key] {
				seen[key] = true
				tmp := make([]int, len(path))
				copy(tmp, path)
				result = append(result, tmp)
			}
			return
		}
		for i := start; i < len(candidates); i++ {
			if candidates[i] > remaining {
				break
			}
			bt(i+1, remaining-candidates[i], append(path, candidates[i]))
		}
	}
	bt(0, target, nil)
	return result
}

// ── Approach 2: Backtracking with Skip-Duplicate Pruning (Optimal) ────────────
//
// backtracking solves Combination Sum II without any extra data structures by
// skipping duplicate candidates at the same recursion level.
//
// Intuition: After sorting, identical values are adjacent. If we're at position
// i within a recursion level (start <= i), and candidates[i] == candidates[i-1],
// choosing candidates[i] again would produce the exact same combination that
// was already explored when we chose candidates[i-1] at this level. Skip it.
//
// Key distinction from #39:
//   - Recurse with i+1 (not i): each element used at most once per combination.
//   - Skip: if i > start && candidates[i] == candidates[i-1]: continue.
//
// Algorithm:
//  1. Sort candidates.
//  2. bt(start, remaining, path):
//     if remaining == 0: record path.
//     for i = start to n-1:
//       if candidates[i] > remaining: break.
//       if i > start && candidates[i] == candidates[i-1]: continue (skip dup).
//       bt(i+1, remaining-candidates[i], path+[candidates[i]]).
//
// Time:  O(2^N) — each element is either included or excluded; pruning reduces this
// Space: O(N) — recursion depth
func backtracking(candidates []int, target int) [][]int {
	sort.Ints(candidates)
	var result [][]int
	var bt func(start, remaining int, path []int)
	bt = func(start, remaining int, path []int) {
		if remaining == 0 {
			tmp := make([]int, len(path))
			copy(tmp, path)
			result = append(result, tmp)
			return
		}
		for i := start; i < len(candidates); i++ {
			if candidates[i] > remaining {
				break // sorted: no further candidates can help
			}
			// skip duplicates at the same recursion level
			if i > start && candidates[i] == candidates[i-1] {
				continue
			}
			bt(i+1, remaining-candidates[i], append(path, candidates[i]))
		}
	}
	bt(0, target, nil)
	return result
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Map Dedup) ===")
	fmt.Printf("candidates=[10,1,2,7,6,1,5] target=8  => %v\n  expected [[1 1 6] [1 2 5] [1 7] [2 6]]\n",
		bruteForce([]int{10, 1, 2, 7, 6, 1, 5}, 8))
	fmt.Printf("candidates=[2,5,2,1,2] target=5  => %v  expected [[1 2 2] [5]]\n",
		bruteForce([]int{2, 5, 2, 1, 2}, 5))

	fmt.Println("\n=== Approach 2: Backtracking Skip-Dup (Optimal) ===")
	fmt.Printf("candidates=[10,1,2,7,6,1,5] target=8  => %v\n  expected [[1 1 6] [1 2 5] [1 7] [2 6]]\n",
		backtracking([]int{10, 1, 2, 7, 6, 1, 5}, 8))
	fmt.Printf("candidates=[2,5,2,1,2] target=5  => %v  expected [[1 2 2] [5]]\n",
		backtracking([]int{2, 5, 2, 1, 2}, 5))
}
