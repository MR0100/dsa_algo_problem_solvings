package main

import "fmt"

// LeetCode #288 — Unique Word Abbreviation
//
// The abbreviation of a word is: first letter + (number of characters strictly
// between first and last) + last letter. Words of length ≤ 2 abbreviate to
// themselves. Examples: "dog" → "d1g", "internationalization" → "i18n",
// "it" → "it", "a" → "a".
//
// Implement ValidWordAbbr(dictionary):
//   isUnique(word) returns true if EITHER
//     (a) no other word in the dictionary shares `word`'s abbreviation, OR
//     (b) every dictionary word sharing that abbreviation is exactly `word`.
// In other words: a word's abbreviation is "unique" iff, among dictionary
// words, the only word mapping to that abbreviation is `word` itself.

// ── Approach: Abbreviation → Set of Words (Hash Map, Optimal) ─────────────────
//
// abbreviate builds a word's abbreviation. Words of length ≤ 2 are their own
// abbreviation because "first + count + last" would not be shorter.
//
// Intuition:
//
//	The predicate "no OTHER dictionary word shares this abbreviation" needs,
//	for each abbreviation, the set of distinct dictionary words that produced
//	it. Store abbr → set-of-words. Then isUnique(word) is true when that set,
//	restricted to the abbreviation of `word`, is either empty or {word}.
//
// Algorithm (build):
//  1. For each dictionary word, compute its abbreviation.
//  2. Map[abbr] = set of distinct original words with that abbreviation.
//
// Time:  O(N * L) to build (N words, average length L).
// Space: O(N * L) for the map of abbreviations to word sets.
func abbreviate(word string) string {
	n := len(word)
	if n <= 2 {
		return word // too short to compress; itself is the "abbreviation"
	}
	// first letter + (middle length as a number) + last letter, e.g. i18n.
	return fmt.Sprintf("%c%d%c", word[0], n-2, word[n-1])
}

// ValidWordAbbr stores, per abbreviation, the set of distinct dictionary words
// that map to it.
type ValidWordAbbr struct {
	// abbr → set of distinct words producing that abbreviation.
	abbrToWords map[string]map[string]struct{}
}

// Constructor builds the ValidWordAbbr from the given dictionary.
//
// Time:  O(N * L). Space: O(N * L).
func Constructor(dictionary []string) ValidWordAbbr {
	m := make(map[string]map[string]struct{})
	for _, w := range dictionary {
		a := abbreviate(w)
		if m[a] == nil {
			m[a] = make(map[string]struct{})
		}
		m[a][w] = struct{}{} // dedupe: identical dictionary words collapse
	}
	return ValidWordAbbr{abbrToWords: m}
}

// isUnique reports whether `word`'s abbreviation is unique in the dictionary.
//
// Intuition:
//
//	Look up the set of dictionary words sharing `word`'s abbreviation. Unique
//	means that set contains nothing but `word` itself: either it is empty
//	(abbreviation unseen) or it is exactly {word} (only `word`, possibly
//	appearing multiple times in the dictionary, maps here).
//
// Algorithm:
//  1. a = abbreviate(word); words = map[a].
//  2. Return true if words is empty, OR (len(words) == 1 AND word ∈ words).
//
// Time:  O(L) — one abbreviation plus O(1) map lookups.
// Space: O(1) beyond the stored structure.
func (v *ValidWordAbbr) isUnique(word string) bool {
	a := abbreviate(word)
	words, ok := v.abbrToWords[a]
	if !ok || len(words) == 0 {
		return true // abbreviation never seen → unique
	}
	// Unique iff the ONLY word with this abbreviation is `word` itself.
	if _, contains := words[word]; contains && len(words) == 1 {
		return true
	}
	return false
}

func main() {
	// Official example:
	// ValidWordAbbr(["deer","door","cake","card"])
	//   isUnique("dear")  -> false  (abbr d2r == "deer"/"door", not "dear")
	//   isUnique("cart")  -> true   (abbr c2t matches no dictionary word)
	//   isUnique("cane")  -> false  (abbr c2e == "cake", which is not "cane")
	//   isUnique("make")  -> true   (abbr m2e matches no dictionary word)
	//   isUnique("cake")  -> true   (abbr c2e's only owner IS "cake")
	dict := []string{"deer", "door", "cake", "card"}
	vwa := Constructor(dict)

	fmt.Println("=== Approach: Abbreviation → Word Set (Hash Map) ===")
	fmt.Println(vwa.isUnique("dear")) // expected false
	fmt.Println(vwa.isUnique("cart")) // expected true
	fmt.Println(vwa.isUnique("cane")) // expected false
	fmt.Println(vwa.isUnique("make")) // expected true
	fmt.Println(vwa.isUnique("cake")) // expected true

	// Extra checks:
	// duplicate dictionary words collapse; short words abbreviate to themselves.
	vwa2 := Constructor([]string{"hello", "hello"})
	fmt.Println(vwa2.isUnique("hello")) // expected true  (only owner of h3o is "hello")
	vwa3 := Constructor([]string{"a", "a"})
	fmt.Println(vwa3.isUnique("a")) // expected true  ("a" abbreviates to "a", only owner)
	vwa4 := Constructor([]string{"hello"})
	fmt.Println(vwa4.isUnique("leetcode")) // expected true  (abbr l6e unseen)
}
