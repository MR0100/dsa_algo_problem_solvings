package main

import (
	"fmt"
	"strings"
)

// ── Approach 1: Brute Force (Word-by-Word Simulation) ────────────────────────
//
// bruteForce solves Sentence Screen Fitting by literally placing words onto
// the screen one at a time, row by row, wrapping when a word no longer fits.
//
// Intuition:
//
//	The screen is rows×cols characters. Fill it exactly as a person would:
//	keep a running column position; place the next word if it fits in the
//	remaining space of the current row (a following word needs a leading
//	space); otherwise move to the next row. Every time the word pointer wraps
//	past the last word, one full sentence has been placed — count it. After
//	`rows` rows, report how many complete sentences fit.
//
// Algorithm:
//  1. wordIdx = 0, count = 0.
//  2. For each of the `rows` rows: colsRemaining = cols.
//     While the current word (plus a space if not first on the line) fits:
//     place it, advance colsRemaining, advance wordIdx (wrapping, and
//     incrementing count on wrap).
//  3. Return count.
//
// Time:  O(rows · cols) worst case — filling nearly every cell (each word ≥ 1
//
//	char, so at most cols placements per row).
//
// Space: O(1) beyond the input word lengths.
func bruteForce(rows, cols int, sentence []string) int {
	n := len(sentence)
	wordIdx := 0 // which word comes next
	count := 0   // completed sentences

	for r := 0; r < rows; r++ {
		remaining := cols // free columns left in this row
		for {
			wordLen := len(sentence[wordIdx])
			// A word that isn't first on the line needs a leading space, so its
			// "cost" is wordLen when remaining == cols, else wordLen+1.
			need := wordLen
			if remaining != cols {
				need = wordLen + 1 // account for the separating space
			}
			if need > remaining {
				break // this word can't start on the current row
			}
			remaining -= need // consume the word (and its space if any)
			wordIdx++         // move to the next word
			if wordIdx == n { // wrapped past the last word …
				wordIdx = 0 // … restart the sentence …
				count++     // … one full sentence placed
			}
		}
	}
	return count
}

// ── Approach 2: Precomputed "Next Index" Simulation ──────────────────────────
//
// nextIndexSim solves Sentence Screen Fitting by simulating one row at a time
// (not one word at a time): for each possible starting word it precomputes,
// via a stream, how many words a single row consumes and where it lands.
//
// Intuition:
//
//	Brute force re-walks individual words; instead treat each row as a chunk of
//	`cols` characters cut out of the infinitely repeated "w0 w1 … wn-1 w0 …"
//	stream (words joined by single spaces). Maintain a global character pointer
//	`start` into that stream. Advancing one row means jumping `start` forward
//	by `cols`; then, if the character landing at the row's right boundary is a
//	space we're cleanly between words, otherwise we back up to the previous
//	space so we don't split a word. `start / lenSentence` at the end counts how
//	many full sentences were consumed.
//
// Algorithm:
//  1. Build s = words joined by spaces + a trailing space; L = len(s).
//  2. start = 0. For each row: start += cols.
//     - if s[start % L] == ' ': start++ (skip the boundary space).
//     - else: while start > 0 and s[(start-1) % L] != ' ': start-- (retreat to
//     a word boundary).
//  3. Return start / L (number of complete sentence copies consumed).
//
// Time:  O(rows · maxWordLen) — the retreat per row is bounded by the longest
//
//	word (≤ 10), so effectively O(rows).
//
// Space: O(total sentence length) for the joined string.
func nextIndexSim(rows, cols int, sentence []string) int {
	s := strings.Join(sentence, " ") + " " // one leading-normalised stream copy
	L := len(s)                            // length of a single sentence+space block
	start := 0                             // char index (into the infinite repeat) of the next row's first slot

	for r := 0; r < rows; r++ {
		start += cols // tentatively consume `cols` characters this row
		if s[start%L] == ' ' {
			start++ // landed exactly on a separating space — step over it
		} else {
			// Landed inside a word: retreat to the space before that word so we
			// don't cut it in half.
			for start > 0 && s[(start-1)%L] != ' ' {
				start--
			}
		}
	}
	return start / L // how many full sentence blocks were used up
}

// ── Approach 3: DP Over Starting Word (Optimal for many rows) ────────────────
//
// dpStartWord solves Sentence Screen Fitting by precomputing, for each word
// index i, how many words a single row can fit when it BEGINS with word i and
// which word index the next row will begin with. Rows then reduce to following
// these precomputed jumps.
//
// Intuition:
//
//	Because a row always starts on some word boundary, its behaviour depends
//	only on which word it starts with — there are just n distinct row "types".
//	Precompute two tables keyed by that starting word:
//	  wordsPlaced[i] = number of words fitted in a row that starts at word i,
//	  nextStart[i]   = the word index the following row starts at.
//	Then walk `rows` rows following nextStart, summing wordsPlaced; the grand
//	total of words placed, divided by n, is the number of full sentences.
//
// Algorithm:
//  1. For each starting word i: greedily pack words into a width-`cols` row,
//     recording how many words fit (wordsPlaced[i]) and the resulting next
//     starting index (nextStart[i]).
//  2. cur = 0, totalWords = 0. Repeat `rows` times:
//     totalWords += wordsPlaced[cur]; cur = nextStart[cur].
//  3. Return totalWords / n.
//
// Time:  O(n · cols) precompute + O(rows) simulation.
// Space: O(n) for the two tables.
func dpStartWord(rows, cols int, sentence []string) int {
	n := len(sentence)
	wordsPlaced := make([]int, n) // words fitted when a row starts at index i
	nextStart := make([]int, n)   // starting index of the next row

	for i := 0; i < n; i++ {
		length := 0 // characters used so far in this hypothetical row
		words := 0  // words placed so far
		idx := i    // walking word pointer
		// Greedily add words while they fit. After the first word, each new word
		// costs len(word)+1 (its leading space).
		for length+len(sentence[idx]) <= cols {
			length += len(sentence[idx]) + 1 // +1 reserves the trailing space
			words++
			idx = (idx + 1) % n // wrap around the sentence
		}
		wordsPlaced[i] = words
		nextStart[i] = idx // where the next row will begin
	}

	cur := 0        // starting word of the current row
	totalWords := 0 // words placed across all rows
	for r := 0; r < rows; r++ {
		totalWords += wordsPlaced[cur] // add this row's contribution
		cur = nextStart[cur]           // jump to next row's starting word
	}
	return totalWords / n // each full sentence is n words
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (word-by-word) ===")
	fmt.Println(bruteForce(2, 8, []string{"hello", "world"}))           // expected 1
	fmt.Println(bruteForce(3, 6, []string{"a", "bcd", "e"}))            // expected 2
	fmt.Println(bruteForce(4, 5, []string{"I", "had", "apple", "pie"})) // expected 1

	fmt.Println("=== Approach 2: Precomputed Next-Index Simulation ===")
	fmt.Println(nextIndexSim(2, 8, []string{"hello", "world"}))           // expected 1
	fmt.Println(nextIndexSim(3, 6, []string{"a", "bcd", "e"}))            // expected 2
	fmt.Println(nextIndexSim(4, 5, []string{"I", "had", "apple", "pie"})) // expected 1

	fmt.Println("=== Approach 3: DP Over Starting Word (Optimal) ===")
	fmt.Println(dpStartWord(2, 8, []string{"hello", "world"}))           // expected 1
	fmt.Println(dpStartWord(3, 6, []string{"a", "bcd", "e"}))            // expected 2
	fmt.Println(dpStartWord(4, 5, []string{"I", "had", "apple", "pie"})) // expected 1
}
