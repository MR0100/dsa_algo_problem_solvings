package main

import "fmt"

// ── Approach 1: Brute Force (Generate All Combinations) ──────────────────────
//
// bruteForce solves Combination Sum by generating all combinations (with
// repetition) of all subsets and filtering by sum.
//
// Intuition: Try all subsets of candidates where each candidate may be chosen
// unlimited times. Filter for subsets with sum == target.
// Extremely inefficient; included to motivate the backtracking approach.
//
// Time:  O(N^(T/M)) where N=len(candidates), T=target, M=min(candidates)
// Space: O(T/M) — maximum recursion depth
func bruteForce(candidates []int, target int) [][]int {
	var result [][]int
	var dfs func(start, remaining int, path []int)
	dfs = func(start, remaining int, path []int) {
		if remaining == 0 {
			tmp := make([]int, len(path))
			copy(tmp, path)
			result = append(result, tmp)
			return
		}
		if remaining < 0 {
			return
		}
		// try all candidates (allow repeats by not advancing start)
		for i := start; i < len(candidates); i++ {
			dfs(i, remaining-candidates[i], append(path, candidates[i]))
		}
	}
	dfs(0, target, nil)
	return result
}

// ── Approach 2: Backtracking with Sorting Pruning (Optimal) ──────────────────
//
// backtracking solves Combination Sum using backtracking with an early-exit
// pruning: sort candidates first; if the current candidate exceeds the remaining
// target, all subsequent candidates also do (they are larger) → break.
//
// Intuition: Build combinations by deciding how many of each candidate to include.
// At each step, either:
//   a. Include candidates[i] again (remaining decreases; same i).
//   b. Move to the next candidate (i+1).
// Sorting + break avoids exploring branches that can never reach target.
//
// Algorithm:
//  1. Sort candidates.
//  2. backtrack(start, remaining, path):
//     if remaining == 0: record path.
//     for i = start to len(candidates)-1:
//       if candidates[i] > remaining: break (pruning).
//       add candidates[i] to path.
//       backtrack(i, remaining-candidates[i], path).  // i, not i+1, allows reuse
//       remove candidates[i] from path.
//
// Time:  O(N^(T/M)) worst case; pruning reduces this significantly in practice
// Space: O(T/M) — recursion depth (path length)
func backtracking(candidates []int, target int) [][]int {
	// sort to enable pruning
	sortInts(candidates)
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
				break // sorted: all further candidates are also too large
			}
			bt(i, remaining-candidates[i], append(path, candidates[i]))
		}
	}
	bt(0, target, nil)
	return result
}

func sortInts(a []int) {
	for i := 1; i < len(a); i++ {
		for j := i; j > 0 && a[j] < a[j-1]; j-- {
			a[j], a[j-1] = a[j-1], a[j]
		}
	}
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Printf("candidates=[2,3,6,7] target=7  => %v  expected [[2 2 3] [7]]\n",
		bruteForce([]int{2, 3, 6, 7}, 7))
	fmt.Printf("candidates=[2,3,5]   target=8  => %v  expected [[2 2 2 2] [2 3 3] [3 5]]\n",
		bruteForce([]int{2, 3, 5}, 8))
	fmt.Printf("candidates=[2]       target=1  => %v  expected []\n",
		bruteForce([]int{2}, 1))

	fmt.Println("\n=== Approach 2: Backtracking with Pruning (Optimal) ===")
	fmt.Printf("candidates=[2,3,6,7] target=7  => %v  expected [[2 2 3] [7]]\n",
		backtracking([]int{2, 3, 6, 7}, 7))
	fmt.Printf("candidates=[2,3,5]   target=8  => %v  expected [[2 2 2 2] [2 3 3] [3 5]]\n",
		backtracking([]int{2, 3, 5}, 8))
	fmt.Printf("candidates=[2]       target=1  => %v  expected []\n",
		backtracking([]int{2}, 1))
}
