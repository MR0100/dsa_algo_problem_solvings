package main

import "fmt"

// ── Approach 1: Brute Force (Scan-and-Compare) ───────────────────────────────
//
// bruteForce solves First Unique Character in a String by, for each index,
// scanning the whole string to check whether that character appears anywhere
// else.
//
// Intuition:
//
//	"First non-repeating character" — so for each position ask: does this
//	character occur at any other position? The first index for which the
//	answer is "no" is the answer. Checking each position costs a full pass,
//	so this is O(n²) but needs no extra storage.
//
// Algorithm:
//  1. For i in 0..len-1:
//  2. Scan j over the whole string; if s[j]==s[i] and j!=i, mark duplicate.
//  3. If no duplicate found, return i.
//  4. Return -1 if all characters repeat.
//
// Time:  O(n²) — for each of n indices we scan up to n characters.
// Space: O(1) — no auxiliary structures.
func bruteForce(s string) int {
	for i := 0; i < len(s); i++ {
		unique := true // assume s[i] is unique until proven otherwise
		for j := 0; j < len(s); j++ {
			if i != j && s[j] == s[i] { // found the same char elsewhere
				unique = false
				break // one duplicate is enough to disqualify i
			}
		}
		if unique {
			return i // first index with no duplicate → answer
		}
	}
	return -1 // every character repeats
}

// ── Approach 2: Hash Map Frequency (Two Pass) ────────────────────────────────
//
// hashMap solves First Unique Character in a String by counting occurrences of
// every character in one pass, then scanning again for the first count-1 char.
//
// Intuition:
//
//	A character is "unique" iff it occurs exactly once. So count everything
//	first (pass 1), then walk left-to-right and return the first index whose
//	character has count 1 (pass 2). Order matters only on the second pass,
//	which is why we re-scan the original string rather than the map.
//
// Algorithm:
//  1. Build count[c] = number of occurrences of c (one pass).
//  2. Scan i left-to-right; return first i with count[s[i]] == 1.
//  3. Return -1 if none.
//
// Time:  O(n) — two linear passes.
// Space: O(k) — k = distinct characters (≤ 26 for lowercase letters, O(1)).
func hashMap(s string) int {
	count := make(map[byte]int) // char -> occurrences
	for i := 0; i < len(s); i++ {
		count[s[i]]++ // pass 1: tally every character
	}
	for i := 0; i < len(s); i++ {
		if count[s[i]] == 1 { // pass 2: first char seen exactly once
			return i
		}
	}
	return -1 // no unique character exists
}

// ── Approach 3: Fixed 26-Bucket Array (Optimal) ──────────────────────────────
//
// arrayCount solves First Unique Character in a String using a fixed-size
// integer array of 26 counters (constraint: only lowercase English letters).
//
// Intuition:
//
//	Same two-pass idea as the hash map, but since the alphabet is a known,
//	small, constant set (a–z) we can index a plain array by (c - 'a') instead
//	of hashing. This removes hashing overhead and makes the space strictly
//	O(1) — 26 ints regardless of input length.
//
// Algorithm:
//  1. count[c-'a']++ for every character (pass 1).
//  2. Return the first i with count[s[i]-'a'] == 1 (pass 2).
//  3. Return -1 otherwise.
//
// Time:  O(n) — two linear passes.
// Space: O(1) — exactly 26 counters.
func arrayCount(s string) int {
	var count [26]int // fixed table for 'a'..'z'
	for i := 0; i < len(s); i++ {
		count[s[i]-'a']++ // pass 1: bucket by letter index
	}
	for i := 0; i < len(s); i++ {
		if count[s[i]-'a'] == 1 { // pass 2: first letter with a single count
			return i
		}
	}
	return -1 // all letters repeat
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(bruteForce("leetcode"))     // expected 0
	fmt.Println(bruteForce("loveleetcode")) // expected 2
	fmt.Println(bruteForce("aabb"))         // expected -1

	fmt.Println("=== Approach 2: Hash Map Frequency ===")
	fmt.Println(hashMap("leetcode"))     // expected 0
	fmt.Println(hashMap("loveleetcode")) // expected 2
	fmt.Println(hashMap("aabb"))         // expected -1

	fmt.Println("=== Approach 3: Fixed 26-Bucket Array (Optimal) ===")
	fmt.Println(arrayCount("leetcode"))     // expected 0
	fmt.Println(arrayCount("loveleetcode")) // expected 2
	fmt.Println(arrayCount("aabb"))         // expected -1
}
