package main

import "fmt"

// Shortest Word Distance III — like #243, but word1 and word2 MAY be the same
// string. When they are the same, we must find the minimum distance between two
// DIFFERENT occurrences of that word.

// ── Approach 1: Brute Force (All Valid Pairs) ────────────────────────────────
//
// bruteForce solves it by trying every index of word1 against every index of
// word2, skipping the case i == j (a word is never distance 0 from itself).
//
// Intuition:
//
//	The answer is |i - j| over valid (i, j) pairs. "Valid" adds one rule: when
//	word1 == word2 we need i != j so we compare two distinct positions of the
//	same word.
//
// Algorithm:
//  1. For each i with words[i] == word1:
//     For each j with words[j] == word2 and i != j: min = min(min, |i - j|).
//  2. Return min.
//
// Time:  O(n²) — nested scan over matching positions.
// Space: O(1) — a running minimum.
func bruteForce(words []string, word1 string, word2 string) int {
	min := len(words)
	for i := 0; i < len(words); i++ {
		if words[i] != word1 {
			continue
		}
		for j := 0; j < len(words); j++ {
			if i == j {
				continue // never pair a position with itself
			}
			if words[j] == word2 {
				if d := abs(i - j); d < min {
					min = d
				}
			}
		}
	}
	return min
}

// ── Approach 2: One-Pass Two Pointers (Optimal) ──────────────────────────────
//
// twoPointers solves it in one scan, branching on whether word1 == word2.
//
// Intuition:
//
//	Case A (word1 != word2): identical to #243 — keep the last-seen index of
//	each word and measure the gap whenever we see either.
//
//	Case B (word1 == word2): every occurrence is both a "word1" and a "word2".
//	Keep a single prev index of the last occurrence; each new occurrence pairs
//	with the previous one, giving the gap between consecutive occurrences (the
//	closest possible for the same word). Track the minimum of those.
//
// Algorithm:
//  1. If word1 == word2: scan; on each occurrence, if prev set update min with
//     i - prev, then set prev = i.
//  2. Else: i1 = i2 = -1; on word1 set i1 and try i1-i2; on word2 set i2 and
//     try i2-i1.
//  3. Return min.
//
// Time:  O(n) — one pass.
// Space: O(1) — a couple of indices and a minimum.
func twoPointers(words []string, word1 string, word2 string) int {
	min := len(words)

	if word1 == word2 {
		prev := -1 // index of the previous occurrence of the shared word
		for k, w := range words {
			if w == word1 {
				if prev != -1 && k-prev < min {
					min = k - prev // gap to the previous occurrence
				}
				prev = k // this occurrence becomes the new "previous"
			}
		}
		return min
	}

	// Distinct words: classic last-seen-index tracking.
	i1, i2 := -1, -1
	for k, w := range words {
		switch w {
		case word1:
			i1 = k
			if i2 != -1 && i1-i2 < min {
				min = i1 - i2
			}
		case word2:
			i2 = k
			if i1 != -1 && i2-i1 < min {
				min = i2 - i1
			}
		}
	}
	return min
}

// abs returns the absolute value of x.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	words := []string{"practice", "makes", "perfect", "coding", "makes"}

	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(bruteForce(words, "makes", "coding")) // expected 1
	fmt.Println(bruteForce(words, "makes", "makes"))  // expected 3

	fmt.Println("=== Approach 2: One-Pass Two Pointers (Optimal) ===")
	fmt.Println(twoPointers(words, "makes", "coding")) // expected 1
	fmt.Println(twoPointers(words, "makes", "makes"))  // expected 3
}
