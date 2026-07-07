package main

import "fmt"

// The 3x3 Android lock grid, keys numbered 1..9:
//
//	1 2 3
//	4 5 6
//	7 8 9
//
// A pattern is an ordered sequence of DISTINCT keys of length in [m, n].
// The one twist: the segment between two consecutive keys may not "jump over"
// an unvisited key. E.g. 1 -> 3 passes through 2, so 2 must already be used;
// 1 -> 9 passes through 5, so 5 must already be used. Adjacent / knight-style
// moves (like 1 -> 6, 1 -> 8, 2 -> 7) pass over nothing and are always allowed.
//
// `mid[a][b]` = the key that lies exactly between a and b when they are
// collinear with one key strictly between them; 0 means "nothing in between".

// buildMid returns the 10x10 "middle key" table (index 0 unused).
func buildMid() [10][10]int {
	var mid [10][10]int
	// Three horizontal lines.
	mid[1][3], mid[3][1] = 2, 2
	mid[4][6], mid[6][4] = 5, 5
	mid[7][9], mid[9][7] = 8, 8
	// Three vertical lines.
	mid[1][7], mid[7][1] = 4, 4
	mid[2][8], mid[8][2] = 5, 5
	mid[3][9], mid[9][3] = 6, 6
	// Two diagonals.
	mid[1][9], mid[9][1] = 5, 5
	mid[3][7], mid[7][3] = 5, 5
	return mid
}

// ── Approach 1: Plain Backtracking (DFS over all sequences) ──────────────────
//
// bruteForceBacktracking counts every valid pattern by exploring, from each
// starting key, all ways to extend the current path one distinct key at a time.
//
// Intuition:
//
//	A pattern is just a path through the 9 keys. We grow the path key by key.
//	At each step we may append any UNUSED key `to` provided the segment
//	current->to does not skip an unvisited middle key. Whenever the current
//	path length is within [m, n] we count it.
//
// Algorithm:
//  1. For each key 1..9 as the first key, mark it used and DFS.
//  2. In the DFS, if current depth is in [m, n], increment the count.
//  3. Try every unused `to`: legal iff mid[cur][to] == 0 (no key between) OR
//     that middle key is already used. Recurse, then unmark (backtrack).
//
// Time:  O(9!) worst case — it walks the tree of all key permutations up to
//
//	length n. Tiny in practice (n <= 9).
//
// Space: O(9) recursion depth + the used[] array.
func bruteForceBacktracking(m, n int) int {
	mid := buildMid()
	used := make([]bool, 10) // used[k] == true once key k is in the path
	count := 0

	var dfs func(cur, depth int)
	dfs = func(cur, depth int) {
		if depth >= m && depth <= n {
			count++ // current path is itself a valid pattern
		}
		if depth == n {
			return // cannot grow further; deeper paths exceed n
		}
		for to := 1; to <= 9; to++ {
			if used[to] {
				continue // keys must be distinct
			}
			jumped := mid[cur][to]           // key skipped over, or 0
			if jumped == 0 || used[jumped] { // legal move?
				used[to] = true  // choose
				dfs(to, depth+1) // explore
				used[to] = false // un-choose (backtrack)
			}
		}
	}

	for start := 1; start <= 9; start++ {
		used[start] = true
		dfs(start, 1)
		used[start] = false
	}
	return count
}

// ── Approach 2: Backtracking + Symmetry (Optimal) ────────────────────────────
//
// symmetryBacktracking counts patterns but only runs the DFS from three
// representative starting keys, then multiplies by their symmetry class size.
//
// Intuition:
//
//	The 3x3 grid has 8-fold symmetry (rotations + reflections). By that
//	symmetry every corner (1,3,7,9) starts the same number of patterns, every
//	edge (2,4,6,8) starts the same number, and the center (5) is unique. So we
//	only DFS from one corner, one edge, and the center, then combine:
//	total = 4*fromCorner + 4*fromEdge + 1*fromCenter.
//
// Algorithm:
//  1. Same legality rule (mid table) and DFS as Approach 1.
//  2. Run DFS once from key 1 (a corner), key 2 (an edge), key 5 (center).
//  3. Return 4*count(1) + 4*count(2) + count(5).
//
// Time:  O(9!) same asymptotic tree, but ~3/9 of the constant work.
// Space: O(9) recursion depth + used[].
func symmetryBacktracking(m, n int) int {
	mid := buildMid()
	used := make([]bool, 10)

	// dfs returns how many valid patterns start at `cur` given the current
	// visited set and path depth.
	var dfs func(cur, depth int) int
	dfs = func(cur, depth int) int {
		res := 0
		if depth >= m && depth <= n {
			res++ // this path is a valid pattern
		}
		if depth == n {
			return res
		}
		for to := 1; to <= 9; to++ {
			if used[to] {
				continue
			}
			jumped := mid[cur][to]
			if jumped == 0 || used[jumped] {
				used[to] = true
				res += dfs(to, depth+1)
				used[to] = false
			}
		}
		return res
	}

	countFrom := func(start int) int {
		used[start] = true
		c := dfs(start, 1)
		used[start] = false
		return c
	}

	corner := countFrom(1) // representative corner
	edge := countFrom(2)   // representative edge
	center := countFrom(5) // the center
	return 4*corner + 4*edge + center
}

func main() {
	fmt.Println("=== Approach 1: Plain Backtracking ===")
	fmt.Println(bruteForceBacktracking(1, 1)) // expected 9
	fmt.Println(bruteForceBacktracking(1, 2)) // expected 65
	fmt.Println(bruteForceBacktracking(1, 9)) // expected 389497

	fmt.Println("=== Approach 2: Backtracking + Symmetry (Optimal) ===")
	fmt.Println(symmetryBacktracking(1, 1)) // expected 9
	fmt.Println(symmetryBacktracking(1, 2)) // expected 65
	fmt.Println(symmetryBacktracking(1, 9)) // expected 389497
}
