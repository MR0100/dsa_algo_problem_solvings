package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Sort Each Word as Key ────────────────────────────────────────
//
// groupBySorted solves Group Anagrams by using the sorted version of each
// word as a canonical key in a hash map.
//
// Intuition: Two strings are anagrams iff they have the same character
// frequencies, which means their sorted versions are identical. Group by key.
//
// Time:  O(n * k log k) where n = len(strs), k = max word length
// Space: O(n * k) — the hash map
func groupBySorted(strs []string) [][]string {
	groups := make(map[string][]string)
	for _, s := range strs {
		key := sortedKey(s)
		groups[key] = append(groups[key], s)
	}
	result := make([][]string, 0, len(groups))
	for _, g := range groups {
		result = append(result, g)
	}
	return result
}

func sortedKey(s string) string {
	b := []byte(s)
	sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
	return string(b)
}

// ── Approach 2: Character Frequency Array as Key (Optimal) ───────────────────
//
// groupByFreq solves Group Anagrams using a character frequency array as the
// key instead of sorting. Avoids the O(k log k) sort cost per word.
//
// Intuition: Encode the character frequency of each word as a [26]int array,
// then use it as a map key (converted to string for Go map compatibility).
//
// Time:  O(n * k) — O(k) to build frequency array per word; 26 fixed overhead
// Space: O(n * k)
func groupByFreq(strs []string) [][]string {
	groups := make(map[[26]int][]string)
	for _, s := range strs {
		var freq [26]int
		for _, ch := range s {
			freq[ch-'a']++
		}
		groups[freq] = append(groups[freq], s)
	}
	result := make([][]string, 0, len(groups))
	for _, g := range groups {
		result = append(result, g)
	}
	return result
}

func main() {
	strs1 := []string{"eat", "tea", "tan", "ate", "nat", "bat"}
	strs2 := []string{""}
	strs3 := []string{"a"}

	fmt.Println("=== Approach 1: Sort Key ===")
	r1 := groupBySorted(strs1)
	fmt.Printf("strs=%v\ngroups: count=%d  %v\nexpected 3 groups\n\n", strs1, len(r1), r1)
	fmt.Printf("strs=%v  groups=%v  expected [[\"\"]] \n", strs2, groupBySorted(strs2))
	fmt.Printf("strs=%v  groups=%v  expected [[\"a\"]]\n", strs3, groupBySorted(strs3))

	fmt.Println("\n=== Approach 2: Frequency Array Key (Optimal) ===")
	r2 := groupByFreq(strs1)
	fmt.Printf("strs=%v\ngroups: count=%d  %v\nexpected 3 groups\n\n", strs1, len(r2), r2)
	fmt.Printf("strs=%v  groups=%v  expected [[\"\"]] \n", strs2, groupByFreq(strs2))
	fmt.Printf("strs=%v  groups=%v  expected [[\"a\"]]\n", strs3, groupByFreq(strs3))
}
