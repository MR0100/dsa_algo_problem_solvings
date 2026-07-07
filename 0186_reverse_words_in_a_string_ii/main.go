package main

import (
	"fmt"
	"strings"
)

// ── Approach 1: Brute Force (Split, Reverse, Copy Back) ──────────────────────
//
// bruteForce solves Reverse Words in a String II by splitting the character
// array into words, reversing the word list, and writing the result back.
//
// Intuition:
//
//	Ignore the in-place requirement first. The words are guaranteed to be
//	separated by exactly one space with no leading/trailing spaces, so a plain
//	split gives the word list, reversing that list gives the answer, and a
//	join with single spaces has exactly the original length — it can be copied
//	straight back over s.
//
// Algorithm:
//  1. Convert s to a string and split it on single spaces into words.
//  2. Reverse the words slice with two converging pointers.
//  3. Join with single spaces and copy the bytes back into s.
//
// Time:  O(n) — split, reverse, join, and copy each touch every byte a constant number of times.
// Space: O(n) — the word slice and the joined string are full-size auxiliary copies.
func bruteForce(s []byte) {
	// exactly one space between words, no leading/trailing spaces → Split is safe
	words := strings.Split(string(s), " ")
	// classic two-pointer reversal of the word order
	for i, j := 0, len(words)-1; i < j; i, j = i+1, j-1 {
		words[i], words[j] = words[j], words[i] // swap symmetric entries
	}
	// the joined string has the same length as s (same words, same separators)
	copy(s, strings.Join(words, " "))
}

// ── Approach 2: Two Pointers (Scan From the End) ─────────────────────────────
//
// twoPointers solves Reverse Words in a String II by scanning the array from
// right to left, emitting each word as it is found into an output buffer.
//
// Intuition:
//
//	The reversed sentence is just the words read back-to-front. Walking from
//	the last byte toward the first, every space marks the start of a word we
//	just walked over, so we can append words to a buffer in exactly the order
//	the answer needs — no split helper, no separate reverse step.
//
// Algorithm:
//  1. Keep end = index one past the current word (starts at len(s)).
//  2. Scan i from len(s)-1 down to 0; when s[i] is a space, the word is
//     s[i+1:end] — append it (preceded by a space if the buffer is non-empty)
//     and set end = i.
//  3. After the loop the leftmost word is s[0:end]; append it, then copy the
//     buffer back over s.
//
// Time:  O(n) — every byte is read once and written into the buffer once.
// Space: O(n) — the output buffer holds the full rebuilt sentence.
func twoPointers(s []byte) {
	out := make([]byte, 0, len(s)) // rebuilt sentence, filled right-to-left by word
	end := len(s)                  // one past the end of the word currently being scanned
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == ' ' {
			if len(out) > 0 {
				out = append(out, ' ') // separator before every word except the first emitted
			}
			out = append(out, s[i+1:end]...) // the word we just walked over
			end = i                          // next word ends right before this space
		}
	}
	if len(out) > 0 {
		out = append(out, ' ') // separator before the final (leftmost) word
	}
	out = append(out, s[:end]...) // leftmost word has no space before it
	copy(s, out)                  // write the answer back into the caller's array
}

// ── Approach 3: In-Place Double Reversal (Optimal) ───────────────────────────
//
// inPlaceReversal solves Reverse Words in a String II by reversing the whole
// array and then reversing each word individually — true O(1) extra space.
//
// Intuition:
//
//	Reversing the entire array puts the words into the correct (reversed)
//	order, but each word's letters come out backwards: "the sky is blue" →
//	"eulb si yks eht". A second pass that reverses every word in place fixes
//	the letters without disturbing the word order: "blue is sky the".
//
// Algorithm:
//  1. Reverse s[0:len(s)] entirely.
//  2. Scan for word boundaries (spaces or the end of the array); reverse each
//     word segment s[start:i] in place.
//
// Time:  O(n) — each byte is swapped at most twice (once per reversal pass).
// Space: O(1) — only index variables; all work happens inside s.
func inPlaceReversal(s []byte) {
	reverseRange(s, 0, len(s)-1) // step 1: whole-array reversal flips word order
	start := 0                   // start index of the word currently being scanned
	for i := 0; i <= len(s); i++ {
		// a boundary is either a space or one past the last byte
		if i == len(s) || s[i] == ' ' {
			reverseRange(s, start, i-1) // step 2: un-reverse this word's letters
			start = i + 1               // next word begins after the space
		}
	}
}

// reverseRange reverses s[lo..hi] in place with two converging pointers.
func reverseRange(s []byte, lo, hi int) {
	for lo < hi {
		s[lo], s[hi] = s[hi], s[lo] // swap the outermost pair
		lo++
		hi--
	}
}

func main() {
	// Example 1: s = ["t","h","e"," ","s","k","y"," ","i","s"," ","b","l","u","e"]
	// Example 2: s = ["a"]

	fmt.Println("=== Approach 1: Brute Force (Split, Reverse, Copy Back) ===")
	a1 := []byte("the sky is blue")
	bruteForce(a1)
	fmt.Printf("%q\n", string(a1)) // expected: "blue is sky the"
	a2 := []byte("a")
	bruteForce(a2)
	fmt.Printf("%q\n", string(a2)) // expected: "a"

	fmt.Println("=== Approach 2: Two Pointers (Scan From the End) ===")
	b1 := []byte("the sky is blue")
	twoPointers(b1)
	fmt.Printf("%q\n", string(b1)) // expected: "blue is sky the"
	b2 := []byte("a")
	twoPointers(b2)
	fmt.Printf("%q\n", string(b2)) // expected: "a"

	fmt.Println("=== Approach 3: In-Place Double Reversal (Optimal) ===")
	c1 := []byte("the sky is blue")
	inPlaceReversal(c1)
	fmt.Printf("%q\n", string(c1)) // expected: "blue is sky the"
	c2 := []byte("a")
	inPlaceReversal(c2)
	fmt.Printf("%q\n", string(c2)) // expected: "a"
}
