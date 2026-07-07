package main

import (
	"fmt"
	"strings"
)

// ── Approach 1: Split and Count ───────────────────────────────────────────────
//
// splitCount solves Length of Last Word using strings.Fields to split on whitespace.
//
// Intuition:
//   strings.Fields splits on any whitespace and ignores leading/trailing spaces.
//   The last element of the resulting slice is the last word; return its length.
//
// Time:  O(n) — Fields scans the entire string.
// Space: O(n) — stores all words.
func splitCount(s string) int {
	words := strings.Fields(s)
	if len(words) == 0 {
		return 0
	}
	return len(words[len(words)-1])
}

// ── Approach 2: Reverse Scan ──────────────────────────────────────────────────
//
// reverseScan solves Length of Last Word by scanning from the end.
//
// Intuition:
//   Skip trailing spaces from the right. Then count characters until we hit
//   another space or reach the beginning.
//
// Algorithm:
//   i = len(s)-1
//   skip trailing spaces: while i>=0 and s[i]==' ': i--
//   count word: while i>=0 and s[i]!=' ': i--; count++
//   return count
//
// Time:  O(n) worst case (all spaces then a long word).
//         O(k) in practice where k = trailing_spaces + last_word_length.
// Space: O(1)
func reverseScan(s string) int {
	i := len(s) - 1
	// skip trailing spaces
	for i >= 0 && s[i] == ' ' {
		i--
	}
	count := 0
	// count the last word's characters
	for i >= 0 && s[i] != ' ' {
		count++
		i--
	}
	return count
}

func main() {
	fmt.Println("=== Approach 1: Split and Count ===")
	fmt.Printf("%q  got=%d  expected 5\n", "Hello World", splitCount("Hello World"))
	fmt.Printf("%q  got=%d  expected 4\n", "   fly me   to   the moon  ", splitCount("   fly me   to   the moon  "))
	fmt.Printf("%q  got=%d  expected 6\n", "luffy is still joyboy", splitCount("luffy is still joyboy"))

	fmt.Println("=== Approach 2: Reverse Scan ===")
	fmt.Printf("%q  got=%d  expected 5\n", "Hello World", reverseScan("Hello World"))
	fmt.Printf("%q  got=%d  expected 4\n", "   fly me   to   the moon  ", reverseScan("   fly me   to   the moon  "))
	fmt.Printf("%q  got=%d  expected 6\n", "luffy is still joyboy", reverseScan("luffy is still joyboy"))
}
