package main

import (
	"fmt"
	"sort"
	"strings"
)

// The task: return every DISTINCT non-decreasing subsequence of nums that has
// at least two elements. Order of the outer list does not matter, so to make
// the output comparable across approaches we sort the collected subsequences
// into a canonical order before returning (LeetCode itself accepts any order).

// canonical sorts a list of subsequences so two approaches that discover them
// in different orders still print the same thing. It orders by length first,
// then lexicographically by element values — a stable, human-readable order.
func canonical(res [][]int) [][]int {
	sort.Slice(res, func(a, b int) bool {
		x, y := res[a], res[b]
		if len(x) != len(y) {
			return len(x) < len(y) // shorter subsequences first
		}
		for i := range x {
			if x[i] != y[i] {
				return x[i] < y[i] // then element-by-element ascending
			}
		}
		return false // identical
	})
	return res
}

// ── Approach 1: Backtracking + Hash-Set Dedup (Brute Force) ──────────────────
//
// backtrackSetDedup enumerates ALL non-decreasing subsequences by trying to
// keep or skip every element, then removes duplicates with a hash set keyed on
// the serialized subsequence.
//
// Intuition:
//
//	A subsequence is built by walking left→right and, at each index, deciding
//	"take this element" or "skip it". We may only take nums[i] if it is >= the
//	last element already taken, which guarantees the result is non-decreasing.
//	Because nums can contain repeats (e.g. [4,7,7]), the same subsequence can
//	be produced through different take/skip paths, so we deduplicate at the end
//	by serializing each subsequence to a string and storing it in a set.
//
// Algorithm:
//  1. DFS carrying the current start index and the path built so far.
//  2. If the path length >= 2, serialize it and add to a set of seen paths.
//  3. For i from start..n-1: if path is empty OR nums[i] >= path's last value,
//     append nums[i], recurse from i+1, then pop (backtrack).
//  4. Materialize the set back into a slice of subsequences.
//
// Time:  O(2^n * n) — up to 2^n subsequences, each costing O(n) to copy/hash.
// Space: O(2^n * n) — the set stores every distinct subsequence.
func backtrackSetDedup(nums []int) [][]int {
	seen := map[string][]int{} // serialized path -> the path itself (dedup store)
	path := []int{}            // current subsequence under construction

	var dfs func(start int)
	dfs = func(start int) {
		if len(path) >= 2 { // valid answer: at least two elements
			key := serialize(path) // stable string key for this exact sequence
			if _, ok := seen[key]; !ok {
				cp := make([]int, len(path)) // snapshot: path is mutated later
				copy(cp, path)
				seen[key] = cp // remember this distinct subsequence
			}
		}
		for i := start; i < len(nums); i++ {
			// Only extend when the non-decreasing property is preserved.
			if len(path) == 0 || nums[i] >= path[len(path)-1] {
				path = append(path, nums[i]) // take nums[i]
				dfs(i + 1)                   // recurse on the remaining suffix
				path = path[:len(path)-1]    // undo the take (backtrack)
			}
		}
	}
	dfs(0)

	res := make([][]int, 0, len(seen))
	for _, v := range seen {
		res = append(res, v) // collect all distinct subsequences
	}
	return canonical(res)
}

// serialize turns a subsequence into a comma-joined string usable as a map key.
// Values can be negative, so a plain join with a separator is required (we
// cannot pack into a fixed-width byte because the range is -100..100).
func serialize(path []int) string {
	parts := make([]string, len(path))
	for i, v := range path {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, ",")
}

// ── Approach 2: Backtracking + Per-Level Used Set (Optimal) ──────────────────
//
// backtrackLevelDedup enumerates the same subsequences but avoids ever
// GENERATING a duplicate, by refusing to start two branches with the same
// value at the same recursion depth.
//
// Intuition:
//
//	Duplicates arise only when, at a single recursion level, we choose the same
//	value from two different positions — e.g. picking the first 7 vs the second
//	7 as "the next element" leads to identical subsequences. So at each level we
//	keep a small "used this value already" set; the second time a value shows up
//	as a candidate at that level we skip it. This needs NO global set and no
//	post-processing: every path produced is unique by construction. The
//	non-decreasing constraint (nums[i] >= last taken) is still enforced.
//
// Algorithm:
//  1. DFS carrying start index and path.
//  2. If path length >= 2, record a copy immediately (it is guaranteed unique).
//  3. Maintain a per-call `used` set (map[int]bool) of values already chosen at
//     THIS level.
//  4. For i from start..n-1: skip if used[nums[i]]; skip if it would break the
//     non-decreasing order; otherwise mark used, take, recurse, backtrack.
//
// Time:  O(2^n * n) worst case (all-distinct increasing input has 2^n subseqs),
//
//	but no duplicates are ever built, so it beats Approach 1 on repeated data.
//
// Space: O(n) recursion depth + O(n) per-level used sets, excluding the output.
func backtrackLevelDedup(nums []int) [][]int {
	res := [][]int{}
	path := []int{}

	var dfs func(start int)
	dfs = func(start int) {
		if len(path) >= 2 {
			cp := make([]int, len(path)) // snapshot the current subsequence
			copy(cp, path)
			res = append(res, cp) // guaranteed distinct — no set needed
		}
		used := map[int]bool{} // values already used as a candidate at THIS level
		for i := start; i < len(nums); i++ {
			if used[nums[i]] {
				continue // same value already branched here → would duplicate
			}
			if len(path) > 0 && nums[i] < path[len(path)-1] {
				continue // taking it would break non-decreasing order
			}
			used[nums[i]] = true         // block this value for the rest of the level
			path = append(path, nums[i]) // take nums[i]
			dfs(i + 1)                   // explore subsequences continuing after i
			path = path[:len(path)-1]    // backtrack
		}
	}
	dfs(0)
	return canonical(res)
}

// format renders a list of subsequences in LeetCode's bracketed style so the
// printed output can be eyeballed against the expected answer.
func format(res [][]int) string {
	inner := make([]string, len(res))
	for i, seq := range res {
		parts := make([]string, len(seq))
		for j, v := range seq {
			parts[j] = fmt.Sprintf("%d", v)
		}
		inner[i] = "[" + strings.Join(parts, ",") + "]"
	}
	return "[" + strings.Join(inner, ",") + "]"
}

func main() {
	fmt.Println("=== Approach 1: Backtracking + Hash-Set Dedup (Brute Force) ===")
	fmt.Println(format(backtrackSetDedup([]int{4, 6, 7, 7})))    // expected [[4,6],[4,7],[6,7],[7,7],[4,6,7],[4,7,7],[6,7,7],[4,6,7,7]]
	fmt.Println(format(backtrackSetDedup([]int{4, 4, 3, 2, 1}))) // expected [[4,4]]

	fmt.Println("=== Approach 2: Backtracking + Per-Level Used Set (Optimal) ===")
	fmt.Println(format(backtrackLevelDedup([]int{4, 6, 7, 7})))    // expected [[4,6],[4,7],[6,7],[7,7],[4,6,7],[4,7,7],[6,7,7],[4,6,7,7]]
	fmt.Println(format(backtrackLevelDedup([]int{4, 4, 3, 2, 1}))) // expected [[4,4]]
}
