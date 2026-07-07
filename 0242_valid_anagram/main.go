package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Sorting ──────────────────────────────────────────────────────
//
// sortingApproach solves Valid Anagram by sorting both strings and comparing.
//
// Intuition:
//
//	Two strings are anagrams iff they are the same multiset of characters.
//	Sorting canonicalizes a multiset into a unique string, so anagrams sort to
//	identical strings.
//
// Algorithm:
//  1. If lengths differ, they cannot be anagrams — return false.
//  2. Sort the rune slice of each string.
//  3. Return whether the two sorted strings are equal.
//
// Time:  O(n log n) — dominated by the sort.
// Space: O(n) — rune slices for sorting.
func sortingApproach(s string, t string) bool {
	if len(s) != len(t) {
		return false // different sizes → not an anagram
	}
	a := []rune(s) // copy into a mutable slice we can sort
	b := []rune(t)
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] }) // sort s
	sort.Slice(b, func(i, j int) bool { return b[i] < b[j] }) // sort t
	return string(a) == string(b)                             // equal canonical forms?
}

// ── Approach 2: Hash Map Frequency Count ─────────────────────────────────────
//
// hashMap solves Valid Anagram by counting each character of s and cancelling
// them out with the characters of t.
//
// Intuition:
//
//	If we add one for every char in s and subtract one for every char in t, an
//	anagram leaves every count at exactly zero. A single map handles the full
//	(Unicode) alphabet without a fixed-size array.
//
// Algorithm:
//  1. If lengths differ, return false.
//  2. For each rune in s: count[r]++.
//  3. For each rune in t: count[r]--.
//  4. Every entry must be 0 → return true; any non-zero → false.
//
// Time:  O(n) — one pass to add, one to subtract, one to check.
// Space: O(k) — k distinct characters.
func hashMap(s string, t string) bool {
	if len(s) != len(t) {
		return false
	}
	count := make(map[rune]int) // char → net occurrence delta
	for _, r := range s {
		count[r]++ // s contributes +1
	}
	for _, r := range t {
		count[r]-- // t contributes -1
	}
	for _, c := range count {
		if c != 0 {
			return false // mismatch in some character's count
		}
	}
	return true
}

// ── Approach 3: Fixed Array Counter (Optimal, lowercase a–z) ─────────────────
//
// arrayCount solves Valid Anagram assuming lowercase English letters, using a
// 26-slot integer array as the counter.
//
// Intuition:
//
//	When the alphabet is a small known set (a–z), a fixed array indexed by
//	letter is faster and lighter than a map: no hashing, cache-friendly. We
//	increment on s, decrement on t, and can early-exit if any slot goes
//	negative (t has a letter s lacks).
//
// Algorithm:
//  1. If lengths differ, return false.
//  2. counts[26] all zero.
//  3. For i in range: counts[s[i]-'a']++ and counts[t[i]-'a']--.
//  4. Any non-zero slot → false; otherwise true.
//
// Time:  O(n) — single fused pass over both strings.
// Space: O(1) — a fixed 26-element array regardless of input size.
func arrayCount(s string, t string) bool {
	if len(s) != len(t) {
		return false
	}
	var counts [26]int // one bucket per lowercase letter
	for i := 0; i < len(s); i++ {
		counts[s[i]-'a']++ // s adds
		counts[t[i]-'a']-- // t removes, in the same loop
	}
	for _, c := range counts {
		if c != 0 {
			return false // some letter is unbalanced
		}
	}
	return true
}

func main() {
	fmt.Println("=== Approach 1: Sorting ===")
	fmt.Println(sortingApproach("anagram", "nagaram")) // expected true
	fmt.Println(sortingApproach("rat", "car"))         // expected false

	fmt.Println("=== Approach 2: Hash Map Frequency Count ===")
	fmt.Println(hashMap("anagram", "nagaram")) // expected true
	fmt.Println(hashMap("rat", "car"))         // expected false

	fmt.Println("=== Approach 3: Fixed Array Counter (Optimal, a–z) ===")
	fmt.Println(arrayCount("anagram", "nagaram")) // expected true
	fmt.Println(arrayCount("rat", "car"))         // expected false
}
