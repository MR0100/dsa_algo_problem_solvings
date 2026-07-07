package main

import (
	"fmt"
	"strings"
)

// LeetCode #290 — Word Pattern
//
// Given a `pattern` and a string `s`, return true if `s` follows the same
// pattern. "Follows" means a BIJECTION between letters in `pattern` and
// non-empty words in `s`: each letter maps to exactly one word, and each word
// maps to exactly one letter.

// ── Approach 1: Two Hash Maps (Bijection Check, Optimal) ─────────────────────
//
// twoMaps checks the pattern by maintaining both directions of the mapping:
// letter → word and word → letter.
//
// Intuition:
//
//	A bijection requires consistency in BOTH directions. One map (letter →
//	word) alone fails cases like pattern="ab", s="dog dog": 'a'→"dog" and
//	'b'→"dog" would pass, yet two letters share one word. So track the
//	reverse map too and reject any conflict.
//
// Algorithm:
//  1. Split s into words; if count ≠ len(pattern), return false.
//  2. For each (letter, word):
//     - If letter already mapped, it must map to this word.
//     - If word already mapped, it must map to this letter.
//     - Otherwise record both directions.
//  3. If no conflict arises, return true.
//
// Time:  O(n) — n = total characters of pattern + s.
// Space: O(k) — k distinct letters/words in the maps.
func twoMaps(pattern string, s string) bool {
	words := strings.Fields(s) // split on whitespace into non-empty words
	if len(words) != len(pattern) {
		return false // lengths must line up 1-to-1
	}
	letterToWord := make(map[byte]string)
	wordToLetter := make(map[string]byte)
	for i := 0; i < len(pattern); i++ {
		c := pattern[i]
		w := words[i]
		// Forward direction: letter must map to the same word every time.
		if mapped, ok := letterToWord[c]; ok {
			if mapped != w {
				return false
			}
		} else {
			letterToWord[c] = w
		}
		// Reverse direction: word must map back to the same letter.
		if mapped, ok := wordToLetter[w]; ok {
			if mapped != c {
				return false
			}
		} else {
			wordToLetter[w] = c
		}
	}
	return true
}

// ── Approach 2: Single Map of First-Seen Index ───────────────────────────────
//
// indexMap checks the pattern by comparing, for each position, the index at
// which the current letter and the current word were FIRST seen.
//
// Intuition:
//
//	Two sequences are "the same pattern" iff, at every position, the last
//	time you saw this token matches for both sequences. Concretely, if the
//	first occurrence of pattern[i] happened at the same index as the first
//	occurrence of words[i], the two are in lockstep — a bijection. Store the
//	first-seen index of each letter and each word; they must always agree.
//
// Algorithm:
//  1. Split s; length must equal len(pattern).
//  2. Keep two maps: letter → first index, word → first index.
//  3. At position i, the recorded first index for pattern[i] must equal the
//     recorded first index for words[i]; then record i for any unseen token.
//
// Time:  O(n). Space: O(k).
func indexMap(pattern string, s string) bool {
	words := strings.Fields(s)
	if len(words) != len(pattern) {
		return false
	}
	letterIdx := make(map[byte]int) // letter → first index it appeared
	wordIdx := make(map[string]int) // word → first index it appeared
	for i := 0; i < len(pattern); i++ {
		c := pattern[i]
		w := words[i]
		li, lok := letterIdx[c]
		wi, wok := wordIdx[w]
		// Either both are new (first sighting) or both point to the same
		// earlier index. Any mismatch breaks the bijection.
		if lok != wok {
			return false
		}
		if lok && wok && li != wi {
			return false
		}
		if !lok {
			letterIdx[c] = i
		}
		if !wok {
			wordIdx[w] = i
		}
	}
	return true
}

func main() {
	// Example 1: pattern="abba", s="dog cat cat dog" → true
	// Example 2: pattern="abba", s="dog cat cat fish" → false
	// Example 3: pattern="aaaa", s="dog cat cat dog" → false
	fmt.Println("=== Approach 1: Two Hash Maps (Bijection) ===")
	fmt.Println(twoMaps("abba", "dog cat cat dog"))  // expected true
	fmt.Println(twoMaps("abba", "dog cat cat fish")) // expected false
	fmt.Println(twoMaps("aaaa", "dog cat cat dog"))  // expected false
	fmt.Println(twoMaps("abba", "dog dog dog dog"))  // expected false (b→dog & a→dog collide)

	fmt.Println("=== Approach 2: Single First-Seen Index Map ===")
	fmt.Println(indexMap("abba", "dog cat cat dog"))  // expected true
	fmt.Println(indexMap("abba", "dog cat cat fish")) // expected false
	fmt.Println(indexMap("aaaa", "dog cat cat dog"))  // expected false
	fmt.Println(indexMap("abba", "dog dog dog dog"))  // expected false
}
