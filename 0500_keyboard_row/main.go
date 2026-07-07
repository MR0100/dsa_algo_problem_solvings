package main

import (
	"fmt"
	"strings"
)

// The three rows of a US QWERTY keyboard (lowercase). Each character in a word
// must come from exactly one of these rows for the word to qualify.
var (
	row1 = "qwertyuiop"
	row2 = "asdfghjkl"
	row3 = "zxcvbnm"
)

// ── Approach 1: Per-Word Row Scan with Sets (Brute Force) ─────────────────────
//
// bruteForce solves Keyboard Row by, for each word, discovering which row its
// first letter belongs to, then verifying every remaining letter lives in that
// same row using a set membership test.
//
// Intuition:
//
//	A word is typeable on one row iff all its letters share a row. Find the row
//	of the first (lower-cased) letter, then check that each subsequent letter is
//	in that row's character set. If any letter escapes the row, reject the word.
//
// Algorithm:
//  1. Build a set of characters for each of the three rows.
//  2. For each word: lowercase it, look up the row-set of its first char.
//  3. Scan the rest; if any char is not in that set, discard the word.
//  4. Otherwise keep the original word.
//
// Time:  O(N·L) — N words, each of length L, each char an O(1) set lookup.
// Space: O(1) — the three fixed row-sets (≤26 entries) plus the output.
func bruteForce(words []string) []string {
	// map each character to a small set for O(1) membership
	rowSets := []map[rune]bool{{}, {}, {}}
	for i, row := range []string{row1, row2, row3} {
		for _, ch := range row {
			rowSets[i][ch] = true
		}
	}

	result := []string{}
	for _, word := range words {
		lower := strings.ToLower(word) // 'A' and 'a' share a key
		// pick the row that owns the first character
		var target map[rune]bool
		first := rune(lower[0])
		for _, set := range rowSets {
			if set[first] {
				target = set
				break
			}
		}
		// verify every character belongs to that same row
		ok := true
		for _, ch := range lower {
			if !target[ch] {
				ok = false // a letter from a different row → reject
				break
			}
		}
		if ok {
			result = append(result, word) // keep the ORIGINAL (preserve case)
		}
	}
	return result
}

// ── Approach 2: Row-Index Lookup Table (Optimal) ─────────────────────────────
//
// rowIndexTable solves Keyboard Row by precomputing, for every letter, an
// integer row id (0/1/2), then accepting a word iff all its letters map to the
// same row id.
//
// Intuition:
//
//	Replace three set lookups with a single 26-entry array rowOf[letter] = row.
//	A word qualifies exactly when every letter has the same rowOf value as the
//	first letter — one integer comparison per character, branch-free and cache
//	friendly.
//
// Algorithm:
//  1. Fill rowOf[c] = r for each character c in row r (a..z indexed 0..25).
//  2. For each word: r0 = rowOf[first letter].
//  3. If every other letter has rowOf == r0, keep the word.
//
// Time:  O(N·L) — same asymptotics, fewer constant-factor operations.
// Space: O(1) — a fixed 26-slot table plus the output.
func rowIndexTable(words []string) []string {
	var rowOf [26]int // rowOf[letter-'a'] = which keyboard row (0,1,2)
	for r, row := range []string{row1, row2, row3} {
		for _, ch := range row {
			rowOf[ch-'a'] = r // record this letter's row
		}
	}

	result := []string{}
	for _, word := range words {
		lower := strings.ToLower(word)
		r0 := rowOf[lower[0]-'a'] // the row every letter must match
		ok := true
		for i := 1; i < len(lower); i++ {
			if rowOf[lower[i]-'a'] != r0 { // a letter from another row?
				ok = false
				break
			}
		}
		if ok {
			result = append(result, word) // preserve original casing
		}
	}
	return result
}

func main() {
	fmt.Println("=== Approach 1: Per-Word Row Scan with Sets (Brute Force) ===")
	fmt.Println(bruteForce([]string{"Hello", "Alaska", "Dad", "Peace"})) // expected [Alaska Dad]
	fmt.Println(bruteForce([]string{"omk"}))                             // expected []
	fmt.Println(bruteForce([]string{"adsdf", "sfd"}))                    // expected [adsdf sfd]

	fmt.Println("=== Approach 2: Row-Index Lookup Table (Optimal) ===")
	fmt.Println(rowIndexTable([]string{"Hello", "Alaska", "Dad", "Peace"})) // expected [Alaska Dad]
	fmt.Println(rowIndexTable([]string{"omk"}))                             // expected []
	fmt.Println(rowIndexTable([]string{"adsdf", "sfd"}))                    // expected [adsdf sfd]
}
