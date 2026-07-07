package main

import "fmt"

// ── Approach 1: Brute Force with Shared-Letter Check ─────────────────────────
//
// bruteForce solves Maximum Product of Word Lengths by comparing every pair of
// words and, for each pair, checking character-by-character whether they share
// any letter.
//
// Intuition:
//
//	Two words qualify iff they share no common letter. Directly test that for
//	every pair using nested loops over their characters, and track the largest
//	length product among qualifying pairs.
//
// Algorithm:
//  1. For each pair (i, j):
//     a. Scan word[i] × word[j] for any equal character; if found, skip.
//     b. Otherwise product = len(word[i]) * len(word[j]); update the max.
//  2. Return the max (0 if no valid pair).
//
// Time:  O(n² · L²) where L is the max word length — nested pair loop with an
//
//	O(L²) shared-letter test.
//
// Space: O(1).
func bruteForce(words []string) int {
	best := 0
	for i := 0; i < len(words); i++ {
		for j := i + 1; j < len(words); j++ {
			if !shareLetter(words[i], words[j]) { // disjoint letter sets?
				if p := len(words[i]) * len(words[j]); p > best {
					best = p
				}
			}
		}
	}
	return best
}

// shareLetter reports whether a and b have at least one character in common.
func shareLetter(a, b string) bool {
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(b); j++ {
			if a[i] == b[j] {
				return true
			}
		}
	}
	return false
}

// ── Approach 2: HashSet Per Word ─────────────────────────────────────────────
//
// hashSet solves the problem by precomputing each word's set of letters, then
// intersecting sets pairwise.
//
// Intuition:
//
//	Build a letter set once per word so the shared-letter test becomes a set
//	membership scan (≤ 26 checks) rather than an O(L²) double loop.
//
// Algorithm:
//  1. For each word, build sets[i] = map of its distinct letters.
//  2. For each pair (i, j): if no letter of the smaller set is in the other,
//     the pair is valid; update the max product.
//  3. Return the max.
//
// Time:  O(Σ|word| + n² · 26) — building sets, then a bounded pair test.
// Space: O(n · 26) — a set per word.
func hashSet(words []string) int {
	sets := make([]map[byte]bool, len(words))
	for i, w := range words {
		s := make(map[byte]bool, 26)
		for k := 0; k < len(w); k++ {
			s[w[k]] = true // record each distinct letter
		}
		sets[i] = s
	}
	best := 0
	for i := 0; i < len(words); i++ {
		for j := i + 1; j < len(words); j++ {
			disjoint := true
			for c := range sets[i] { // check every letter of word i
				if sets[j][c] { // present in word j → not disjoint
					disjoint = false
					break
				}
			}
			if disjoint {
				if p := len(words[i]) * len(words[j]); p > best {
					best = p
				}
			}
		}
	}
	return best
}

// ── Approach 3: Bitmask (Optimal) ────────────────────────────────────────────
//
// bitmask solves Maximum Product of Word Lengths by encoding each word's letter
// set into a 26-bit integer, so the shared-letter test is a single AND.
//
// Intuition:
//
//	Only the SET of distinct letters matters. Encode word w as a bitmask where
//	bit (c-'a') is 1 if letter c appears. Two words share no letter iff
//	mask[i] & mask[j] == 0. That turns the inner test into one machine
//	operation, making the whole thing O(n²) with a tiny constant.
//
// Algorithm:
//  1. For each word, compute masks[i] = OR of (1 << (c-'a')).
//  2. For each pair (i, j): if masks[i] & masks[j] == 0, they are disjoint;
//     update max product with len[i]*len[j].
//  3. Return the max.
//
// Time:  O(Σ|word| + n²) — building masks, then O(1)-per-pair loop.
// Space: O(n) — one integer mask per word.
func bitmask(words []string) int {
	masks := make([]int, len(words)) // 26-bit letter set per word
	for i, w := range words {
		m := 0
		for k := 0; k < len(w); k++ {
			m |= 1 << (w[k] - 'a') // set the bit for this letter
		}
		masks[i] = m
	}
	best := 0
	for i := 0; i < len(words); i++ {
		for j := i + 1; j < len(words); j++ {
			if masks[i]&masks[j] == 0 { // no common bit → no common letter
				if p := len(words[i]) * len(words[j]); p > best {
					best = p
				}
			}
		}
	}
	return best
}

func main() {
	ex1 := []string{"abcw", "baz", "foo", "bar", "xtfn", "abcdef"}
	ex2 := []string{"a", "ab", "abc", "d", "cd", "bcd", "abcd"}
	ex3 := []string{"a", "aa", "aaa", "aaaa"}

	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Printf("ex1 -> %d  expected 16\n", bruteForce(ex1)) // "abcw" & "xtfn" -> 4*4
	fmt.Printf("ex2 -> %d  expected 4\n", bruteForce(ex2))  // "ab" & "cd" -> 2*2
	fmt.Printf("ex3 -> %d  expected 0\n", bruteForce(ex3))  // all share 'a'

	fmt.Println("=== Approach 2: HashSet Per Word ===")
	fmt.Printf("ex1 -> %d  expected 16\n", hashSet(ex1))
	fmt.Printf("ex2 -> %d  expected 4\n", hashSet(ex2))
	fmt.Printf("ex3 -> %d  expected 0\n", hashSet(ex3))

	fmt.Println("=== Approach 3: Bitmask (Optimal) ===")
	fmt.Printf("ex1 -> %d  expected 16\n", bitmask(ex1))
	fmt.Printf("ex2 -> %d  expected 4\n", bitmask(ex2))
	fmt.Printf("ex3 -> %d  expected 0\n", bitmask(ex3))
}
