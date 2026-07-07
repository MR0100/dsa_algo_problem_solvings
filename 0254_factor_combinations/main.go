package main

import "fmt"

// ── Approach 1: Backtracking over Divisors ───────────────────────────────────
//
// backtracking solves Factor Combinations by recursively choosing factors in
// non-decreasing order, dividing n down as it goes.
//
// Intuition:
//
//	Every factorization of n (excluding the trivial [n] and factors 1 or n) is a
//	multiset of factors ≥ 2 whose product is n. To avoid duplicates like [2,6]
//	vs [6,2], we force factors to appear in non-decreasing order by only trying
//	factors ≥ the last one chosen. At each step we pick a factor f that divides
//	the remaining value, add it to the path, and recurse on remaining/f. The
//	remaining value itself (if ≥ start) also closes off a valid combination.
//
// Algorithm:
//  1. dfs(remaining, start, path):
//     for f = start; f*f <= remaining; f++:
//     if remaining % f == 0:
//     record path + [f, remaining/f]  (a complete factorization)
//     recurse dfs(remaining/f, f, path+[f])  (factor the quotient further)
//  2. Kick off with dfs(n, 2, []). Factors start at 2 and never exceed sqrt.
//
// Time:  O(number of factorizations · depth) — bounded by the divisor tree of n.
// Space: O(log n) recursion depth (each division at least halves… roughly).
func backtracking(n int) [][]int {
	result := [][]int{}
	var dfs func(remaining, start int, path []int)
	dfs = func(remaining, start int, path []int) {
		// Try each candidate factor f from `start` up to sqrt(remaining).
		for f := start; f*f <= remaining; f++ {
			if remaining%f == 0 { // f is a valid factor of what's left
				// Complete factorization: current path plus f and its cofactor.
				comb := make([]int, len(path))
				copy(comb, path)
				comb = append(comb, f, remaining/f)
				result = append(result, comb)

				// Recurse: keep factoring the quotient, factors ≥ f (non-decreasing).
				dfs(remaining/f, f, append(path, f))
			}
		}
	}
	dfs(n, 2, []int{}) // factors are ≥ 2; the whole number n is the target
	return result
}

// ── Approach 2: Backtracking with Explicit Result Aggregation (Alt Form) ──────
//
// dfsCollect solves Factor Combinations with a helper that returns the list of
// factorizations of a value using factors ≥ start. It is the same recursion
// framed as "combine current factor with sub-factorizations of the quotient".
//
// Intuition:
//
//	factorize(n, start) = for each divisor f in [start, n): pair f with n/f
//	directly ([f, n/f]), and also with every factorization of n/f using factors
//	≥ f. This makes the recursive structure explicit: the answer for n is built
//	from the answers for its quotients.
//
// Algorithm:
//  1. helper(value, start) returns [][]int:
//     for f = start; f*f <= value; f++ if value%f==0:
//     add [f, value/f]
//     for each sub in helper(value/f, f): add [f] ++ sub
//  2. Return helper(n, 2).
//
// Time:  O(number of factorizations · depth), same class as Approach 1.
// Space: O(log n) recursion plus the produced lists.
func dfsCollect(n int) [][]int {
	var helper func(value, start int) [][]int
	helper = func(value, start int) [][]int {
		res := [][]int{}
		for f := start; f*f <= value; f++ {
			if value%f == 0 { // f divides value
				co := value / f
				res = append(res, []int{f, co}) // the simple two-factor split

				// Prepend f to each deeper factorization of the cofactor.
				for _, sub := range helper(co, f) {
					combined := append([]int{f}, sub...)
					res = append(res, combined)
				}
			}
		}
		return res
	}
	return helper(n, 2)
}

// equalNested reports whether two [][]int are element-wise equal (helper for demo).
func equalNested(a, b [][]int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

func main() {
	// Official examples:
	// n = 1  -> []
	// n = 12 -> [[2,6],[2,2,3],[3,4]]
	// n = 37 -> []   (prime)
	fmt.Println("=== Approach 1: Backtracking ===")
	fmt.Println(backtracking(1))  // expected []
	fmt.Println(backtracking(12)) // expected [[2 6] [2 2 3] [3 4]]
	fmt.Println(backtracking(37)) // expected []
	fmt.Println(backtracking(32)) // expected [[2 16] [2 2 8] [2 2 2 4] [2 2 2 2 2] [2 4 4] [4 8]]

	fmt.Println("=== Approach 2: DFS Collect (Alt Form) ===")
	fmt.Println(dfsCollect(1))  // expected []
	fmt.Println(dfsCollect(12)) // expected [[2 6] [2 2 3] [3 4]]
	fmt.Println(dfsCollect(37)) // expected []

	// Sanity: both approaches agree on n = 12.
	fmt.Println("=== Consistency check (n=12) ===")
	fmt.Println(equalNested(backtracking(12), dfsCollect(12))) // expected true
}
