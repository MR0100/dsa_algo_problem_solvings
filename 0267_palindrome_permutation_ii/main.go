package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Backtracking over Half Permutations (Optimal) ────────────────
//
// backtrack solves Palindrome Permutation II by generating every distinct
// permutation of one half of the string and mirroring it around the (optional)
// middle character.
//
// Intuition:
//
//	A palindrome is fully determined by its left half plus at most one middle
//	character. First check feasibility: at most one character may have an odd
//	count. Take half of each character's count to form a multiset for the left
//	half; permute that multiset (handling duplicates so results are unique),
//	and for each permutation emit half + middle + reverse(half).
//
// Algorithm:
//  1. Count characters. If more than one has an odd count, return [] (no
//     palindrome possible).
//  2. Build `half` = each character repeated count/2 times; if a character has
//     an odd count, remember it as the single middle character.
//  3. Backtrack to permute `half` without producing duplicate arrangements
//     (sort, then skip a value equal to its previous sibling when the previous
//     was not used at this level).
//  4. For each complete half, form half + mid + reversed(half).
//
// Time:  O((n/2)! · n) — permutations of the half times the O(n) cost to build
//
//	each palindrome (far fewer when duplicates prune the tree).
//
// Space: O(n) — counts, the half buffer, and recursion depth n/2.
func backtrack(s string) []string {
	counts := make(map[byte]int)
	for i := 0; i < len(s); i++ {
		counts[s[i]]++ // frequency of each character
	}

	// Feasibility + build the half multiset and the middle character.
	var mid string
	oddCount := 0
	half := make([]byte, 0, len(s)/2)
	// Iterate characters in sorted order for deterministic, sorted output.
	keys := make([]byte, 0, len(counts))
	for c := range counts {
		keys = append(keys, c)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	for _, c := range keys {
		if counts[c]%2 == 1 { // this char has an odd count
			oddCount++
			mid = string(c) // candidate middle character
		}
		for k := 0; k < counts[c]/2; k++ {
			half = append(half, c) // each char contributes count/2 to the half
		}
	}
	if oddCount > 1 {
		return []string{} // more than one odd char -> impossible
	}

	var result []string
	used := make([]bool, len(half))
	cur := make([]byte, 0, len(half))

	var dfs func()
	dfs = func() {
		if len(cur) == len(half) { // a full left half is assembled
			// Build palindrome: half + mid + reverse(half).
			left := string(cur)
			rev := make([]byte, len(cur))
			for i := range cur {
				rev[len(cur)-1-i] = cur[i] // reversed copy
			}
			result = append(result, left+mid+string(rev))
			return
		}
		for i := 0; i < len(half); i++ {
			if used[i] {
				continue // already placed at an earlier position
			}
			// Skip duplicates: if this char equals the previous one and the
			// previous is unused at this level, choosing it now repeats work.
			if i > 0 && half[i] == half[i-1] && !used[i-1] {
				continue
			}
			used[i] = true
			cur = append(cur, half[i])
			dfs()
			cur = cur[:len(cur)-1] // undo choice
			used[i] = false
		}
	}
	dfs()
	return result
}

func main() {
	fmt.Println("=== Approach 1: Backtracking (Optimal) ===")
	fmt.Println(backtrack("aabb"))         // expected [abba baab]
	fmt.Println(backtrack("abc"))          // expected []
	fmt.Println(backtrack("aabbhijkkjih")) // one palindrome per distinct half
}
