package main

import (
	"fmt"
	"sort"
)

// LeetCode 301 — Remove Invalid Parentheses
//
// Given a string s that contains parentheses and letters, remove the MINIMUM
// number of invalid parentheses to make the input string valid. Return ALL
// distinct valid strings reachable with the minimum number of removals, in any
// order.

// ── Approach 1: BFS Level-by-Level (Minimum Removals) ────────────────────────
//
// bfs solves Remove Invalid Parentheses by exploring the state space of strings
// one removal at a time, level by level, and stopping at the first level that
// contains any valid string.
//
// Intuition:
//
//	"Minimum number of removals" is a shortest-path notion: the number of
//	removals is the depth in a tree where each node's children are the string
//	with one more character deleted. BFS from the original string reaches the
//	shallowest (fewest-removal) valid strings first. The moment we find any
//	valid string on a level, that level's removal count is minimal — collect
//	every valid string on that level and stop; deeper levels remove more.
//
// Algorithm:
//  1. Start a BFS queue with the original string; use a visited set to dedupe.
//  2. For each level, scan every string; if any is valid, mark found = true
//     and collect the valid ones (do NOT expand further — deeper = more removals).
//  3. If nothing on the level was valid, generate all children by deleting each
//     single character once, enqueue unseen ones, and continue to the next level.
//  4. Return the collected valid strings (guaranteed minimal removals).
//
// Time:  O(2^n · n) worst case — up to 2^n subsequences, each O(n) to validate.
// Space: O(2^n · n) — the visited set / queue can hold exponentially many states.
func bfs(s string) []string {
	result := []string{}         // valid strings found on the minimal level
	visited := map[string]bool{} // strings already enqueued (dedupe)
	queue := []string{s}         // BFS frontier, starts with the whole string
	visited[s] = true            // mark the root as seen
	found := false               // becomes true once a valid string appears

	for len(queue) > 0 { // process the tree level by level
		next := []string{} // children to explore at the next depth
		for _, cur := range queue {
			if isValidParen(cur) { // this string needs no more removals
				result = append(result, cur)
				found = true // remember that this whole level is the answer level
			}
			if found {
				continue // once found on this level, never expand deeper
			}
			// Generate all strings with exactly one more character removed.
			for i := 0; i < len(cur); i++ {
				c := cur[i]
				if c != '(' && c != ')' {
					continue // only removing parentheses can help validity
				}
				child := cur[:i] + cur[i+1:] // delete character i
				if !visited[child] {
					visited[child] = true // avoid re-expanding duplicates
					next = append(next, child)
				}
			}
		}
		if found {
			break // minimal level fully collected — deeper levels remove more
		}
		queue = next // descend one level (one additional removal)
	}
	sort.Strings(result) // deterministic order for stable test output
	return result
}

// isValidParen reports whether s has balanced parentheses (letters ignored).
//
// Time:  O(n) — single pass.
// Space: O(1) — one counter.
func isValidParen(s string) bool {
	bal := 0 // running (open − close) count
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(':
			bal++
		case ')':
			bal--
			if bal < 0 { // a ')' with no matching '(' before it
				return false
			}
		}
	}
	return bal == 0 // every '(' must be closed
}

// ── Approach 2: DFS with Precomputed Removal Counts (Optimal) ─────────────────
//
// dfsBacktrack solves Remove Invalid Parentheses by first counting exactly how
// many '(' and ')' must be removed, then backtracking to try every way of
// deleting precisely that many, validating each candidate.
//
// Intuition:
//
//	The minimum removals are fully determined by one left-to-right pass:
//	each ')' with no available '(' is a mandatory close-removal, and any '('
//	still unmatched at the end is a mandatory open-removal. Knowing the exact
//	quotas (leftRem, rightRem) lets us prune aggressively: we only ever build
//	strings that delete exactly that many parens, and we skip consecutive
//	duplicate characters to avoid generating the same string twice.
//
// Algorithm:
//  1. Compute leftRem (unmatched '(') and rightRem (unmatched ')') in one pass.
//  2. DFS over indices carrying the current built string and remaining quotas:
//     - When both quotas hit 0, validate the remainder is balanced and record.
//     - At each index we may DELETE the current paren (if quota > 0), skipping
//     runs of identical chars so each distinct deletion is tried once.
//     - Or KEEP the current char and move on.
//  3. Collect all valid results (already distinct thanks to the skip rule).
//
// Time:  O(2^n · n) worst case, but pruned heavily by the exact quotas.
// Space: O(n) recursion depth (plus the output).
func dfsBacktrack(s string) []string {
	// One pass to compute the mandatory removal quotas.
	leftRem, rightRem := 0, 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(':
			leftRem++ // tentatively an unmatched '('
		case ')':
			if leftRem > 0 {
				leftRem-- // this ')' matches a pending '('
			} else {
				rightRem++ // ')' with nothing to match → must remove
			}
		}
	}

	resultSet := map[string]bool{} // dedupe defensively
	var dfs func(idx int, built string, open, lRem, rRem int)
	dfs = func(idx int, built string, open, lRem, rRem int) {
		if lRem < 0 || rRem < 0 || open < 0 {
			return // pruned: removed too many or unbalanced prefix
		}
		if idx == len(s) {
			if lRem == 0 && rRem == 0 && open == 0 {
				resultSet[built] = true // exact quotas used, balanced → valid
			}
			return
		}
		c := s[idx]
		// Option A: delete a parenthesis, spending the matching quota.
		if c == '(' && lRem > 0 {
			dfs(idx+1, built, open, lRem-1, rRem)
		} else if c == ')' && rRem > 0 {
			dfs(idx+1, built, open, lRem, rRem-1)
		}
		// Option B: keep the current character.
		switch c {
		case '(':
			dfs(idx+1, built+"(", open+1, lRem, rRem) // one more open
		case ')':
			dfs(idx+1, built+")", open-1, lRem, rRem) // closes an open
		default:
			dfs(idx+1, built+string(c), open, lRem, rRem) // letter: passthrough
		}
	}
	dfs(0, "", 0, leftRem, rightRem)

	result := make([]string, 0, len(resultSet))
	for str := range resultSet {
		result = append(result, str)
	}
	sort.Strings(result) // deterministic order for test output
	return result
}

func main() {
	fmt.Println("=== Approach 1: BFS Level-by-Level ===")
	fmt.Println(bfs("()())()"))  // expected [(())() ()()()]
	fmt.Println(bfs("(a)())()")) // expected [(a())() (a)()()]
	fmt.Println(bfs(")("))       // expected []

	fmt.Println("=== Approach 2: DFS with Precomputed Removal Counts (Optimal) ===")
	fmt.Println(dfsBacktrack("()())()"))  // expected [(())() ()()()]
	fmt.Println(dfsBacktrack("(a)())()")) // expected [(a())() (a)()()]
	fmt.Println(dfsBacktrack(")("))       // expected []
}
