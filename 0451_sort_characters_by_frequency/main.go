package main

import (
	"fmt"
	"sort"
	"strings"
)

// ── Approach 1: Count + Sort Unique Characters ───────────────────────────────
//
// countThenSort solves Sort Characters By Frequency by counting each
// character, sorting the distinct characters by descending frequency, and
// re-emitting each character its-count times.
//
// Intuition:
//
//	The output is just each character repeated `count` times, with the
//	characters ordered by how often they appear. So: tally counts, sort the
//	*distinct* characters (at most 62 of them) by count descending, then glue
//	together `char × count` for each.
//
// Algorithm:
//  1. Build freq[byte]int over the string.
//  2. Collect the distinct bytes into a slice.
//  3. Sort that slice by freq descending.
//  4. For each byte in sorted order, append it freq[b] times to a builder.
//
// Time:  O(n + k log k) where n = len(s), k = number of distinct chars (≤ 62).
//
//	Counting is O(n); sorting the ≤62 distinct chars is O(k log k); the
//	rebuild writes n bytes. So effectively O(n).
//
// Space: O(n) for the output (Go strings are immutable, we build a new one);
//
//	O(k) for the frequency map and the key slice.
func countThenSort(s string) string {
	freq := make(map[byte]int) // how many times each byte occurs
	for i := 0; i < len(s); i++ {
		freq[s[i]]++ // tally this character
	}

	// Collect the distinct characters so we can order them by frequency.
	chars := make([]byte, 0, len(freq))
	for b := range freq {
		chars = append(chars, b)
	}

	// Sort distinct characters by descending frequency (ties: any order is OK).
	sort.Slice(chars, func(i, j int) bool {
		return freq[chars[i]] > freq[chars[j]] // higher count comes first
	})

	var sb strings.Builder    // efficient string assembly, no repeated realloc
	sb.Grow(len(s))           // we know the final length up front
	for _, b := range chars { // most frequent first
		// Write this character exactly freq[b] times (block them together).
		sb.WriteString(strings.Repeat(string(b), freq[b]))
	}
	return sb.String()
}

// ── Approach 2: Bucket Sort by Frequency (Optimal) ───────────────────────────
//
// bucketSort solves Sort Characters By Frequency by placing characters into
// buckets indexed by their frequency, avoiding any comparison sort.
//
// Intuition:
//
//	A character's frequency is an integer in [1, n]. That is a perfect bucket
//	index. Put every character into bucket[freq]. Then walk the buckets from
//	high frequency down to low and emit whatever is there. This is counting
//	sort applied to the *frequencies*, so it is linear.
//
// Algorithm:
//  1. Count characters into freq[byte]int.
//  2. Create buckets: buckets[f] holds all characters whose count is f
//     (f ranges 0..n, so len(s)+1 buckets).
//  3. Walk buckets from f = n down to 1; for each character c in bucket f,
//     append c exactly f times.
//
// Time:  O(n) — counting is O(n); buckets total n+1 slots; the emit loop
//
//	writes n bytes overall. No log factor.
//
// Space: O(n) — buckets plus the output string.
func bucketSort(s string) string {
	n := len(s)
	freq := make(map[byte]int) // character → count
	for i := 0; i < n; i++ {
		freq[s[i]]++
	}

	// buckets[f] = list of characters that appear exactly f times.
	// A count can be at most n, so we need indices 0..n.
	buckets := make([][]byte, n+1)
	for b, f := range freq {
		buckets[f] = append(buckets[f], b) // drop char into its frequency bucket
	}

	var sb strings.Builder
	sb.Grow(n)
	// Walk from the highest possible frequency down to 1 (0 is always empty
	// of real characters). This yields descending-frequency output.
	for f := n; f >= 1; f-- {
		for _, b := range buckets[f] { // every char with this exact frequency
			// Emit the character f times so equal chars stay grouped.
			sb.WriteString(strings.Repeat(string(b), f))
		}
	}
	return sb.String()
}

func main() {
	fmt.Println("=== Approach 1: Count + Sort Unique Characters ===")
	fmt.Printf("s=%q -> %q  (valid: any arrangement grouping e,e first)\n", "tree", countThenSort("tree")) // expected "eert" or "eetr"
	fmt.Printf("s=%q -> %q  (valid: aaaccc or cccaaa)\n", "cccaaa", countThenSort("cccaaa"))               // expected "cccaaa" or "aaaccc"
	fmt.Printf("s=%q -> %q  (valid: bb then Aa in any order)\n", "Aabb", countThenSort("Aabb"))            // expected "bbAa" or "bbaA"

	fmt.Println("=== Approach 2: Bucket Sort by Frequency (Optimal) ===")
	fmt.Printf("s=%q -> %q  (valid: any arrangement grouping e,e first)\n", "tree", bucketSort("tree")) // expected "eert" or "eetr"
	fmt.Printf("s=%q -> %q  (valid: aaaccc or cccaaa)\n", "cccaaa", bucketSort("cccaaa"))               // expected "cccaaa" or "aaaccc"
	fmt.Printf("s=%q -> %q  (valid: bb then Aa in any order)\n", "Aabb", bucketSort("Aabb"))            // expected "bbAa" or "bbaA"

	// Deterministic self-check: for any valid answer the multiset of characters
	// and the "grouped + non-increasing frequency" property must hold.
	fmt.Println("=== Validation (frequency-sorted & same multiset) ===")
	for _, in := range []string{"tree", "cccaaa", "Aabb"} {
		out := bucketSort(in)
		fmt.Printf("s=%q  validSortedByFreq=%v\n", in, isFreqSorted(in, out)) // expected true
	}
}

// isFreqSorted verifies an answer without hard-coding one particular ordering:
// (1) output is a permutation of the input, (2) identical characters are
// contiguous, and (3) block lengths are non-increasing.
func isFreqSorted(in, out string) bool {
	if len(in) != len(out) {
		return false // must be a permutation
	}
	// Same multiset of characters?
	fin, fout := map[byte]int{}, map[byte]int{}
	for i := 0; i < len(in); i++ {
		fin[in[i]]++
		fout[out[i]]++
	}
	for b, c := range fin {
		if fout[b] != c {
			return false
		}
	}
	// Walk the output as consecutive equal-character blocks; block lengths must
	// be non-increasing, and no character may appear in two separate blocks.
	seen := map[byte]bool{}
	prevLen := 1 << 30 // +infinity: first block may be any length
	i := 0
	for i < len(out) {
		j := i
		for j < len(out) && out[j] == out[i] { // extend the current block
			j++
		}
		blockLen := j - i
		if seen[out[i]] { // same character split across blocks → not grouped
			return false
		}
		seen[out[i]] = true
		if blockLen > prevLen { // frequencies must not increase going forward
			return false
		}
		prevLen = blockLen
		i = j
	}
	return true
}
