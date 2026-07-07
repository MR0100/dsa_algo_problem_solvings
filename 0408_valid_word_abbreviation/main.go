package main

import (
	"fmt"
	"strconv"
	"unicode"
)

// A word may be abbreviated by replacing some non-adjacent, non-empty substrings
// with their lengths (no leading zeros). Given `word` and `abbr`, decide whether
// `abbr` is a valid abbreviation of `word`. Example: "internationalization" →
// "i12iz4n" is valid; "apple" → "a2e" is not.

// ── Approach 1: Two Pointers (In-Place Parse) — Optimal ───────────────────────
//
// twoPointers solves Valid Word Abbreviation by walking one pointer over word
// and one over abbr, consuming digit runs as "skip counts" and letters as exact
// matches, without building any intermediate string.
//
// Intuition:
//
//	Scan abbr left to right. A letter must match the current word character
//	one-for-one. A digit begins a number that says "advance the word pointer that
//	many characters" (those characters were compressed away). Parse the whole
//	number, reject a leading zero (an abbreviation length can't start with 0),
//	jump the word pointer forward, and continue. The abbreviation is valid iff
//	both pointers finish exactly at the end together.
//
// Algorithm:
//  1. i over word, j over abbr, both from 0.
//  2. While j < len(abbr):
//     - if abbr[j] is a digit: it must not be '0' (leading-zero rule); read the
//     full run of digits into num; advance i by num.
//     - else (a letter): word[i] must exist and equal abbr[j]; advance both.
//  3. Valid iff i == len(word) after abbr is fully consumed.
//
// Time:  O(len(word) + len(abbr)) — single linear scan of each string.
// Space: O(1) — only integer indices and a running number.
func twoPointers(word string, abbr string) bool {
	i, j := 0, 0 // i indexes word, j indexes abbr
	n, m := len(word), len(abbr)
	for j < m {
		if abbr[j] >= '0' && abbr[j] <= '9' {
			// A number token: leading zero is illegal (e.g. "01", "0").
			if abbr[j] == '0' {
				return false
			}
			num := 0
			// Accumulate consecutive digits into the skip count.
			for j < m && abbr[j] >= '0' && abbr[j] <= '9' {
				num = num*10 + int(abbr[j]-'0')
				j++
			}
			i += num // jump over the compressed-away characters
		} else {
			// A literal letter: word must still have a char here and it must match.
			if i >= n || word[i] != abbr[j] {
				return false
			}
			i++ // consume the matched letter in word
			j++ // consume it in abbr
		}
	}
	// Both must land exactly at the end: leftover word chars ⇒ under-covered,
	// i overshooting ⇒ a number ran past the end.
	return i == n
}

// ── Approach 2: Expand Then Compare (Reconstruct) — brute force ───────────────
//
// expandThenCompare solves Valid Word Abbreviation by rebuilding the substring
// of `word` that `abbr` claims to describe and checking it equals `word`.
//
// Intuition:
//
//	Turn abbr back into the string it stands for. Copy letters verbatim; when we
//	meet a number k (rejecting leading zeros), splice in the next k characters of
//	word using a running cursor into word. If at any point the cursor would run
//	past word's end, the abbreviation is invalid. Finally the reconstruction is
//	valid iff we consumed exactly all of word.
//
// Algorithm:
//  1. cursor = 0 (position in word), scan abbr.
//  2. On a digit run: parse k (no leading zero); ensure cursor+k <= len(word);
//     advance cursor by k.
//  3. On a letter: ensure cursor < len(word) and word[cursor] == letter; append
//     and advance cursor.
//  4. Valid iff cursor == len(word) at the end.
//
// Time:  O(len(word) + len(abbr)).
// Space: O(1) extra — we compare against word directly via the cursor.
func expandThenCompare(word string, abbr string) bool {
	cursor := 0 // how many characters of word we have accounted for
	n := len(word)
	j := 0
	runes := []rune(abbr) // treat abbr as runes for clean digit checks
	for j < len(runes) {
		ch := runes[j]
		if unicode.IsDigit(ch) {
			if ch == '0' {
				return false // leading zero not allowed in a length token
			}
			// Read the whole number token.
			start := j
			for j < len(runes) && unicode.IsDigit(runes[j]) {
				j++
			}
			k, _ := strconv.Atoi(string(runes[start:j])) // token value
			cursor += k                                  // these k chars are compressed
			if cursor > n {
				return false // number claims more characters than word has
			}
		} else {
			// Literal letter must line up with word at the cursor.
			if cursor >= n || rune(word[cursor]) != ch {
				return false
			}
			cursor++
			j++
		}
	}
	return cursor == n // every character of word must be accounted for exactly
}

func main() {
	fmt.Println("=== Approach 1: Two Pointers (In-Place Parse) ===")
	fmt.Printf("word=\"internationalization\", abbr=\"i12iz4n\" got=%v  expected true\n", twoPointers("internationalization", "i12iz4n"))
	fmt.Printf("word=\"apple\", abbr=\"a2e\"                     got=%v  expected false\n", twoPointers("apple", "a2e"))
	fmt.Printf("word=\"substitution\", abbr=\"s10n\"             got=%v  expected true\n", twoPointers("substitution", "s10n"))
	fmt.Printf("word=\"substitution\", abbr=\"s55n\"             got=%v  expected false\n", twoPointers("substitution", "s55n"))
	fmt.Printf("word=\"substitution\", abbr=\"s010n\"            got=%v  expected false\n", twoPointers("substitution", "s010n"))
	fmt.Printf("word=\"a\", abbr=\"01\"                          got=%v  expected false\n", twoPointers("a", "01"))
	fmt.Printf("word=\"word\", abbr=\"1o1e\"                     got=%v  expected false\n", twoPointers("word", "1o1e"))

	fmt.Println("=== Approach 2: Expand Then Compare (Reconstruct) ===")
	fmt.Printf("word=\"internationalization\", abbr=\"i12iz4n\" got=%v  expected true\n", expandThenCompare("internationalization", "i12iz4n"))
	fmt.Printf("word=\"apple\", abbr=\"a2e\"                     got=%v  expected false\n", expandThenCompare("apple", "a2e"))
	fmt.Printf("word=\"substitution\", abbr=\"s10n\"             got=%v  expected true\n", expandThenCompare("substitution", "s10n"))
	fmt.Printf("word=\"substitution\", abbr=\"s55n\"             got=%v  expected false\n", expandThenCompare("substitution", "s55n"))
	fmt.Printf("word=\"substitution\", abbr=\"s010n\"            got=%v  expected false\n", expandThenCompare("substitution", "s010n"))
	fmt.Printf("word=\"a\", abbr=\"01\"                          got=%v  expected false\n", expandThenCompare("a", "01"))
	fmt.Printf("word=\"word\", abbr=\"1o1e\"                     got=%v  expected false\n", expandThenCompare("word", "1o1e"))
}
