package main

import (
	"fmt"
	"strings"
)

// ── Approach 1: Brute Force (Split, Reverse, Join) ───────────────────────────
//
// bruteForce solves Reverse Words in a String by splitting into a word list,
// reversing the list, and joining with single spaces.
//
// Intuition:
//
//	The problem is exactly "reverse the order of words". If we can get a clean
//	list of words (no empty strings from extra spaces), reversing that list
//	and joining with " " directly produces the answer. strings.Fields already
//	splits on runs of whitespace and drops leading/trailing spaces for us.
//
// Algorithm:
//  1. Split s into words with strings.Fields (handles multiple/leading/
//     trailing spaces).
//  2. Reverse the words slice in place with two pointers.
//  3. Join the reversed words with a single space.
//
// Time:  O(n) — split, reverse, and join each touch every character once.
// Space: O(n) — the word slice plus the output string.
func bruteForce(s string) string {
	// Fields splits around any run of spaces and never yields empty strings,
	// so "  hello   world  " becomes ["hello","world"].
	words := strings.Fields(s)
	// classic two-pointer in-place reversal of the slice
	for i, j := 0, len(words)-1; i < j; i, j = i+1, j-1 {
		words[i], words[j] = words[j], words[i] // swap symmetric positions
	}
	// join with exactly one space — output format requires single separators
	return strings.Join(words, " ")
}

// ── Approach 2: Two Pointers (Scan From the End) ─────────────────────────────
//
// twoPointers solves Reverse Words in a String by scanning the string from
// right to left and emitting each word as it is found.
//
// Intuition:
//
//	The last word of s must come first in the answer. So walk backwards:
//	skip trailing spaces, mark the end of a word, keep walking until the word
//	starts, and append s[start:end] to the output. Repeat until the front of
//	the string is reached — words come out already in reversed order, so no
//	separate reversal step is needed.
//
// Algorithm:
//  1. Set i to the last index of s.
//  2. While i >= 0:
//     a. Move i left past any spaces.
//     b. If i dropped below 0, stop — only spaces remained.
//     c. Record end = i, then move i left while s[i] is not a space.
//     d. The word is s[i+1 : end+1]; append it to the builder, preceded by
//     a single space unless it is the first word emitted.
//  3. Return the built string.
//
// Time:  O(n) — every character is visited exactly once.
// Space: O(n) — the output builder (no intermediate word slice).
func twoPointers(s string) string {
	var sb strings.Builder // accumulates the answer without repeated copying
	sb.Grow(len(s))        // pre-allocate: the answer is never longer than s
	i := len(s) - 1        // right pointer, starts at the end of the string
	for i >= 0 {
		// skip the run of spaces between words (and trailing spaces)
		for i >= 0 && s[i] == ' ' {
			i--
		}
		if i < 0 { // nothing but spaces left → all words emitted
			break
		}
		end := i // inclusive index of the word's last character
		// walk left until we fall off the word (space or string start)
		for i >= 0 && s[i] != ' ' {
			i--
		}
		if sb.Len() > 0 { // separator only between words, never leading
			sb.WriteByte(' ')
		}
		// s[i+1:end+1] is the whole word: i stopped one left of its start
		sb.WriteString(s[i+1 : end+1])
	}
	return sb.String()
}

// ── Approach 3: In-Place Reversal (Optimal for the Follow-Up) ────────────────
//
// inPlaceReversal solves Reverse Words in a String using the classic
// "reverse everything, then reverse each word" trick with O(1) extra space
// over a single mutable buffer.
//
// Intuition:
//
//	Reversing the entire string reverses the word ORDER (what we want) but
//	also reverses the LETTERS inside each word (what we don't want).
//	Reversing each individual word afterwards fixes the letters back while
//	keeping the new word order. Spaces are first compacted in place with a
//	read/write pointer pair so the buffer holds exactly single-space-separated
//	words before the reversals.
//
// Algorithm:
//  1. Copy s into a byte slice b (Go strings are immutable; in languages
//     with mutable strings this step disappears — hence O(1) extra space).
//  2. Compact spaces in place: read pointer r skips space runs, write
//     pointer w copies each word down, inserting one ' ' before every word
//     except the first. Truncate b to length w.
//  3. Reverse the whole buffer b[0..w-1].
//  4. Scan for word boundaries and reverse each word b[start..end] back.
//  5. Return string(b).
//
// Time:  O(n) — compaction, full reversal, and per-word reversal are each one linear pass.
// Space: O(1) extra beyond the single working buffer (the buffer itself is unavoidable in Go because strings are immutable).
func inPlaceReversal(s string) string {
	b := []byte(s) // single mutable working buffer

	// -- step 1: compact spaces in place ------------------------------------
	w := 0 // write pointer: next position to fill in the cleaned buffer
	r := 0 // read pointer: current position in the original content
	for r < len(b) {
		// skip a run of spaces (leading spaces or separators)
		for r < len(b) && b[r] == ' ' {
			r++
		}
		if r == len(b) { // trailing spaces only → done
			break
		}
		if w > 0 { // one space before every word except the first
			b[w] = ' '
			w++
		}
		// copy the word's bytes down to the write position
		for r < len(b) && b[r] != ' ' {
			b[w] = b[r]
			w++
			r++
		}
	}
	b = b[:w] // truncate: buffer now holds "word word word" exactly

	// -- step 2: reverse the entire buffer (reverses word order + letters) --
	reverseRange(b, 0, len(b)-1)

	// -- step 3: reverse each word back so letters read correctly -----------
	start := 0 // start index of the current word
	for i := 0; i <= len(b); i++ {
		// a word ends at a space or at the end of the buffer
		if i == len(b) || b[i] == ' ' {
			reverseRange(b, start, i-1) // fix this word's letters
			start = i + 1               // next word starts after the space
		}
	}
	return string(b)
}

// reverseRange reverses b[lo..hi] in place with two converging pointers.
func reverseRange(b []byte, lo, hi int) {
	for lo < hi {
		b[lo], b[hi] = b[hi], b[lo] // swap the outer pair
		lo++
		hi--
	}
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Split, Reverse, Join) ===")
	fmt.Printf("%q\n", bruteForce("the sky is blue"))  // "blue is sky the"
	fmt.Printf("%q\n", bruteForce("  hello world  "))  // "world hello"
	fmt.Printf("%q\n", bruteForce("a good   example")) // "example good a"

	fmt.Println("=== Approach 2: Two Pointers (Scan From the End) ===")
	fmt.Printf("%q\n", twoPointers("the sky is blue"))  // "blue is sky the"
	fmt.Printf("%q\n", twoPointers("  hello world  "))  // "world hello"
	fmt.Printf("%q\n", twoPointers("a good   example")) // "example good a"

	fmt.Println("=== Approach 3: In-Place Reversal (Optimal) ===")
	fmt.Printf("%q\n", inPlaceReversal("the sky is blue"))  // "blue is sky the"
	fmt.Printf("%q\n", inPlaceReversal("  hello world  "))  // "world hello"
	fmt.Printf("%q\n", inPlaceReversal("a good   example")) // "example good a"
}
