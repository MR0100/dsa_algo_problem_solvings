package main

import "fmt"

// ── Approach 1: Hash Map Count ───────────────────────────────────────────────
//
// hashMap solves Palindrome Permutation by counting each character's frequency.
//
// Intuition:
//
//	A string can be rearranged into a palindrome iff at most one distinct
//	character has an odd count. (Even counts pair up around the centre; a
//	single odd-count character can sit in the middle.)
//
// Algorithm:
//  1. Count how many times each character appears.
//  2. Count how many characters have an odd frequency.
//  3. Return true if that number of odd-count characters is 0 or 1.
//
// Time:  O(n) — one pass to count, one pass over the map.
// Space: O(k) — k distinct characters.
func hashMap(s string) bool {
	counts := make(map[rune]int) // char -> frequency
	for _, c := range s {
		counts[c]++ // tally every character
	}
	odd := 0 // number of characters seen an odd number of times
	for _, v := range counts {
		if v%2 == 1 { // odd frequency
			odd++
		}
	}
	return odd <= 1 // palindrome possible with at most one odd-count char
}

// ── Approach 2: Single-Pass Set Toggle (Optimal) ─────────────────────────────
//
// setToggle solves Palindrome Permutation by tracking the parity of each char
// in a set, toggling membership as characters are seen.
//
// Intuition:
//
//	We do not need the exact counts — only their parity. Keep a set of chars
//	currently seen an odd number of times: add on first sight, remove on the
//	second, add again on the third, and so on. After the pass the set holds
//	exactly the odd-count characters; its size must be ≤ 1.
//
// Algorithm:
//  1. For each char: if present in the set, delete it; else insert it.
//  2. After the pass, return len(set) <= 1.
//
// Time:  O(n) — single pass over the string.
// Space: O(k) — set of distinct odd-parity characters.
func setToggle(s string) bool {
	seen := make(map[rune]struct{}) // chars currently at odd parity
	for _, c := range s {
		if _, ok := seen[c]; ok { // second (even) occurrence
			delete(seen, c) // parity flips back to even
		} else {
			seen[c] = struct{}{} // odd occurrence
		}
	}
	return len(seen) <= 1 // at most one leftover odd char
}

func main() {
	fmt.Println("=== Approach 1: Hash Map ===")
	fmt.Println(hashMap("code")) // expected false
	fmt.Println(hashMap("aab"))  // expected true
	fmt.Println(hashMap("carerac"))

	fmt.Println("=== Approach 2: Set Toggle (Optimal) ===")
	fmt.Println(setToggle("code")) // expected false
	fmt.Println(setToggle("aab"))  // expected true
	fmt.Println(setToggle("carerac"))
}
