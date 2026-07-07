package main

import (
	"fmt"
	"sort"
	"strings"
)

// LeetCode #192 is one of the four Shell problems ("write a bash script...").
// Per this repo's Go-only rule, each approach re-implements the same pipeline
// in Go: the file words.txt is simulated as an in-memory string, and every
// function returns the report lines "word frequency" sorted by descending
// frequency. The problem guarantees every word's frequency is unique, so the
// order is deterministic; for robustness we still tie-break alphabetically.

// wordsTxt mirrors the official example content of words.txt.
const wordsTxt = `the day is sunny the the
the sunny is is`

// ── Approach 1: Brute Force (Parallel Slices + Selection Sort) ───────────────
//
// bruteForce solves Word Frequency without any hash map: it keeps two parallel
// slices (unique words, counts), scanning linearly for each word, then orders
// the result with a hand-rolled selection sort.
//
// Intuition:
//
//	The most primitive counting device is a list you scan front-to-back:
//	"have I filed this word already? If yes bump its tally, else open a new
//	entry." Sorting by count can likewise be done with the most primitive
//	sort — repeatedly select the remaining maximum. Everything O(1)-lookup-free.
//
// Algorithm:
//  1. Split the file into words on any whitespace run.
//  2. For each word, linearly scan `uniques`; bump counts[i] on a match,
//     otherwise append a new (word, 1) entry.
//  3. Selection sort both slices in lock-step: for each slot i pick the entry
//     with the highest count (alphabetically smallest on a tie) and swap it in.
//  4. Format each pair as "word count".
//
// Time:  O(W·U + U²) — W words each scanned against up to U uniques, plus the
//
//	U² selection sort (U = number of distinct words).
//
// Space: O(U) — the two parallel slices.
func bruteForce(fileContent string) []string {
	words := strings.Fields(fileContent) // split on any run of spaces/newlines/tabs

	uniques := []string{} // distinct words in first-seen order
	counts := []int{}     // counts[i] = occurrences of uniques[i]

	for _, w := range words {
		found := false
		for i := range uniques { // linear scan — the "no hash map" price
			if uniques[i] == w {
				counts[i]++ // seen before → bump its tally
				found = true
				break
			}
		}
		if !found { // first sighting → open a new entry with count 1
			uniques = append(uniques, w)
			counts = append(counts, 1)
		}
	}

	// Selection sort by count descending (alphabetical ascending on ties).
	for i := 0; i < len(uniques); i++ {
		best := i // index of the best remaining entry
		for j := i + 1; j < len(uniques); j++ {
			higher := counts[j] > counts[best]                                  // strictly more frequent wins
			tieAlpha := counts[j] == counts[best] && uniques[j] < uniques[best] // tie → smaller word wins
			if higher || tieAlpha {
				best = j
			}
		}
		// Swap the winner into slot i in BOTH slices to keep them aligned.
		counts[i], counts[best] = counts[best], counts[i]
		uniques[i], uniques[best] = uniques[best], uniques[i]
	}

	out := make([]string, len(uniques))
	for i := range uniques {
		out[i] = fmt.Sprintf("%s %d", uniques[i], counts[i]) // "word count" line
	}
	return out
}

// ── Approach 2: Hash Map + Comparison Sort ───────────────────────────────────
//
// hashMapSort solves Word Frequency the idiomatic way: count occurrences in a
// map for O(1) updates, then sort the distinct words by frequency.
//
// Intuition:
//
//	Counting is the textbook hash-map job — each word update is O(1) instead
//	of a linear scan. The only remaining work is ordering U distinct words,
//	which a comparison sort does in O(U log U). This mirrors the classic Unix
//	pipeline `sort | uniq -c | sort -nr` with the first sort replaced by a map.
//
// Algorithm:
//  1. Split the file into words; freq[word]++ for each.
//  2. Collect the map's keys into a slice.
//  3. sort.Slice by freq descending, word ascending on ties.
//  4. Format "word count" lines.
//
// Time:  O(W + U log U) — one counting pass plus the sort of distinct words.
// Space: O(U) — the frequency map and key slice.
func hashMapSort(fileContent string) []string {
	freq := make(map[string]int) // word → number of occurrences
	for _, w := range strings.Fields(fileContent) {
		freq[w]++ // O(1) amortised update per word
	}

	words := make([]string, 0, len(freq))
	for w := range freq { // gather distinct words (map order is random)
		words = append(words, w)
	}

	sort.Slice(words, func(i, j int) bool {
		if freq[words[i]] != freq[words[j]] {
			return freq[words[i]] > freq[words[j]] // primary: higher count first
		}
		return words[i] < words[j] // secondary: alphabetical (defensive; no ties guaranteed)
	})

	out := make([]string, len(words))
	for i, w := range words {
		out[i] = fmt.Sprintf("%s %d", w, freq[w])
	}
	return out
}

// ── Approach 3: Bucket Sort by Frequency (Optimal) ───────────────────────────
//
// bucketSort solves Word Frequency in linear time by exploiting that a
// frequency can never exceed W, the total word count: words are dropped into
// buckets indexed by their exact frequency, then buckets are read high→low.
//
// Intuition:
//
//	Comparison sorting is overkill when the sort key is a small bounded
//	integer. A word appearing c times satisfies 1 ≤ c ≤ W, so an array of W+1
//	buckets ("all words with count c") replaces the O(U log U) sort with two
//	linear sweeps — the same trick as LeetCode 347 Top K Frequent Elements.
//
// Algorithm:
//  1. Count words into a hash map (as in Approach 2).
//  2. Create buckets[0..W]; append each word w to buckets[freq[w]].
//  3. Walk c from W down to 1; for every word in buckets[c] (alphabetically
//     sorted for determinism) emit "w c".
//
// Time:  O(W + U) — counting pass + one sweep over W+1 buckets holding U words
//
//	(per-bucket sorts are no-ops given the unique-frequency guarantee).
//
// Space: O(W + U) — the bucket array dominates.
func bucketSort(fileContent string) []string {
	words := strings.Fields(fileContent)

	freq := make(map[string]int) // word → occurrences
	for _, w := range words {
		freq[w]++
	}

	// buckets[c] holds every distinct word occurring exactly c times.
	// A frequency can never exceed len(words), so W+1 buckets always suffice.
	buckets := make([][]string, len(words)+1)
	for w, c := range freq {
		buckets[c] = append(buckets[c], w)
	}

	out := []string{}
	for c := len(words); c >= 1; c-- { // highest frequency first
		if len(buckets[c]) == 0 {
			continue // no word has this exact count
		}
		sort.Strings(buckets[c]) // deterministic tie order (problem guarantees ≤1 word here)
		for _, w := range buckets[c] {
			out = append(out, fmt.Sprintf("%s %d", w, c))
		}
	}
	return out
}

// printLines prints each report line on its own row, matching the script output.
func printLines(lines []string) {
	for _, l := range lines {
		fmt.Println(l)
	}
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Parallel Slices + Selection Sort) ===")
	printLines(bruteForce(wordsTxt))
	// expected:
	// the 4
	// is 3
	// sunny 2
	// day 1

	fmt.Println("=== Approach 2: Hash Map + Comparison Sort ===")
	printLines(hashMapSort(wordsTxt))
	// expected:
	// the 4
	// is 3
	// sunny 2
	// day 1

	fmt.Println("=== Approach 3: Bucket Sort by Frequency (Optimal) ===")
	printLines(bucketSort(wordsTxt))
	// expected:
	// the 4
	// is 3
	// sunny 2
	// day 1
}
