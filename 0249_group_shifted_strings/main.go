package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Hash Map by Shift Signature (Optimal) ─────────────────────────
//
// hashMapBySignature groups strings that belong to the same shifting sequence.
//
// Intuition:
//
//	"Shifting" a string moves every letter forward by the same amount, wrapping
//	z→a. Two strings are in the same group iff the GAPS between consecutive
//	letters are identical (mod 26). For "abc" the gaps are (1,1); for "bcd" they
//	are also (1,1); shifting preserves every gap. So a canonical KEY for a group
//	is the tuple of consecutive differences taken modulo 26. Strings with the
//	same key go in the same bucket. Single-character strings all share the empty
//	gap sequence, so they group together.
//
// Algorithm:
//  1. For each string s, compute its signature: for i in 1..len-1,
//     diff = (s[i] - s[i-1] + 26) % 26; join diffs into a key string.
//  2. Append s to the map bucket for that key.
//  3. Collect all buckets as the answer.
//
// Time:  O(N · L) — N strings, each length L, one pass to build each key.
// Space: O(N · L) — the map stores every string plus keys.
func hashMapBySignature(strings []string) [][]string {
	groups := map[string][]string{} // signature -> list of strings

	for _, s := range strings {
		key := shiftSignature(s)             // canonical key for the group
		groups[key] = append(groups[key], s) // bucket the string
	}

	result := make([][]string, 0, len(groups))
	for _, g := range groups {
		result = append(result, g) // one entry per distinct signature
	}
	return result
}

// shiftSignature builds the canonical key: the consecutive letter differences
// modulo 26, separated by commas so multi-digit diffs don't collide (e.g. the
// diff 1 followed by 2 must not read the same as the diff 12).
func shiftSignature(s string) string {
	if len(s) <= 1 {
		return "" // all single chars (and empties) share one group
	}
	key := ""
	for i := 1; i < len(s); i++ {
		// Forward distance from previous letter to this one, wrapping z→a.
		diff := (int(s[i]) - int(s[i-1]) + 26) % 26
		key += fmt.Sprintf("%d,", diff) // comma delimiter avoids ambiguity
	}
	return key
}

// ── Approach 2: Normalize to Canonical Form Starting at 'a' ───────────────────
//
// normalizeToBase shifts every string so that it starts at 'a', using the
// resulting normalized string itself as the group key.
//
// Intuition:
//
//	If we shift a string by exactly the amount that turns its first letter into
//	'a', every string in the same shift-group collapses to the SAME normalized
//	string. E.g. "abc" → "abc" (first is already 'a'), and "bcd" shifted back
//	by 1 → "abc" too. That normalized form is a natural, human-readable key —
//	an alternative to the difference signature that produces identical grouping.
//
// Algorithm:
//  1. For each string s, compute shift = s[0]-'a'.
//  2. Build normalized: each char c → 'a' + (c - 'a' - shift + 26) % 26.
//  3. Bucket s under its normalized key; collect buckets.
//
// Time:  O(N · L).
// Space: O(N · L).
func normalizeToBase(strings []string) [][]string {
	groups := map[string][]string{}

	for _, s := range strings {
		key := normalize(s)                  // shift so first char is 'a'
		groups[key] = append(groups[key], s) // bucket by normalized form
	}

	result := make([][]string, 0, len(groups))
	for _, g := range groups {
		result = append(result, g)
	}
	return result
}

// normalize shifts s backward so its first character becomes 'a'.
func normalize(s string) string {
	if len(s) == 0 {
		return ""
	}
	shift := int(s[0] - 'a') // how far to shift back
	buf := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		// Shift each letter back by `shift`, wrapping around the alphabet.
		buf[i] = byte('a' + (int(s[i]-'a')-shift+26)%26)
	}
	return string(buf)
}

// canonicalize sorts the groups and each group's contents so the two approaches
// can be printed deterministically for comparison.
func canonicalize(groups [][]string) [][]string {
	out := make([][]string, 0, len(groups))
	for _, g := range groups {
		gc := append([]string(nil), g...)
		sort.Strings(gc)
		out = append(out, gc)
	}
	sort.Slice(out, func(i, j int) bool {
		return fmt.Sprint(out[i]) < fmt.Sprint(out[j])
	})
	return out
}

func main() {
	ex1 := []string{"abc", "bcd", "acef", "xyz", "az", "ba", "a", "z"}
	ex2 := []string{"a"}

	fmt.Println("=== Approach 1: Hash Map by Shift Signature (Optimal) ===")
	// expected groups: [[a z] [abc bcd xyz] [acef] [az ba]]
	fmt.Printf("ex1 got=%v\n", canonicalize(hashMapBySignature(ex1)))
	fmt.Printf("ex1 numGroups=%d  expected 4\n", len(hashMapBySignature(ex1))) // expected 4
	fmt.Printf("ex2 got=%v  expected [[a]]\n", canonicalize(hashMapBySignature(ex2)))

	fmt.Println("=== Approach 2: Normalize to Canonical Form ===")
	fmt.Printf("ex1 got=%v\n", canonicalize(normalizeToBase(ex1)))
	fmt.Printf("ex1 numGroups=%d  expected 4\n", len(normalizeToBase(ex1))) // expected 4
	fmt.Printf("ex2 got=%v  expected [[a]]\n", canonicalize(normalizeToBase(ex2)))
}
