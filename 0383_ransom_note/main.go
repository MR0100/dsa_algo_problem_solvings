package main

import "fmt"

// ── Approach 1: Brute Force (Consume from a mutable copy) ─────────────────────
//
// bruteForce solves Ransom Note by, for each letter needed in the note,
// linearly scanning the (mutable) magazine for a matching, not-yet-used letter
// and crossing it off.
//
// Intuition:
//
//	The note can be built iff every one of its letters can be paired with a
//	distinct letter in the magazine. Model "distinct" by physically removing
//	the letter from a copy of the magazine once it is used, so it cannot be
//	reused. If any note letter finds no free match, construction fails.
//
// Algorithm:
//  1. Copy the magazine into a mutable byte slice.
//  2. For each character c of the note, scan the copy for c.
//  3. If found, blank that slot (mark it used) and move on.
//  4. If a scan finds nothing, return false. If all letters are matched, true.
//
// Time:  O(n·m) — for each of n note chars we may scan all m magazine chars.
// Space: O(m) — the mutable copy of the magazine.
func bruteForce(ransomNote string, magazine string) bool {
	mag := []byte(magazine) // mutable copy so we can cross letters off
	for i := 0; i < len(ransomNote); i++ {
		c := ransomNote[i] // the letter we currently need
		found := false
		for j := 0; j < len(mag); j++ {
			if mag[j] == c { // an unused magazine letter matches
				mag[j] = 0 // consume it so it can't be reused
				found = true
				break
			}
		}
		if !found { // no free copy of c left in the magazine
			return false
		}
	}
	return true // every note letter was matched to a distinct magazine letter
}

// ── Approach 2: Hash Map Counting ────────────────────────────────────────────
//
// hashMap solves Ransom Note by counting how many of each letter the magazine
// supplies, then decrementing as the note consumes them.
//
// Intuition:
//
//	Only the multiset of letters matters, not their order. Tally the
//	magazine's letters once; then each note letter just needs one unit of
//	that letter's remaining stock. This turns the nested scan into two linear
//	passes.
//
// Algorithm:
//  1. Build count[c] = number of times c appears in the magazine.
//  2. For each note char c: if count[c] == 0 → false; else count[c]--.
//  3. Survive the whole note → true.
//
// Time:  O(n + m) — one pass to count, one pass to spend.
// Space: O(k) — k = number of distinct letters (≤ 26 here, so effectively O(1)).
func hashMap(ransomNote string, magazine string) bool {
	count := map[byte]int{} // letter → available quantity
	for i := 0; i < len(magazine); i++ {
		count[magazine[i]]++ // tally everything the magazine offers
	}
	for i := 0; i < len(ransomNote); i++ {
		c := ransomNote[i]
		if count[c] == 0 { // note needs c but stock is exhausted
			return false
		}
		count[c]-- // spend one copy of c
	}
	return true // every needed letter was in stock
}

// ── Approach 3: Fixed Array Counting (Optimal) ───────────────────────────────
//
// arrayCount solves Ransom Note using a 26-slot integer array instead of a
// hash map, exploiting the "lowercase English letters" constraint.
//
// Intuition:
//
//	Same counting idea as the hash map, but the alphabet is exactly 'a'..'z',
//	so an index c-'a' into a length-26 array replaces every hash operation:
//	no hashing overhead, perfect cache locality, and it is the least code.
//	One early exit: if the note is longer than the magazine it is impossible.
//
// Algorithm:
//  1. If len(note) > len(magazine) → false immediately.
//  2. cnt[c-'a']++ for each magazine char.
//  3. For each note char, cnt[c-'a']--; if it goes below 0 → false.
//  4. Otherwise → true.
//
// Time:  O(n + m) — two linear passes, O(1) work each.
// Space: O(1) — a fixed 26-integer array regardless of input size.
func arrayCount(ransomNote string, magazine string) bool {
	if len(ransomNote) > len(magazine) { // can't cover more letters than exist
		return false
	}
	var cnt [26]int // cnt[i] = available count of letter 'a'+i
	for i := 0; i < len(magazine); i++ {
		cnt[magazine[i]-'a']++ // count magazine letters
	}
	for i := 0; i < len(ransomNote); i++ {
		cnt[ransomNote[i]-'a']-- // spend one for this note letter
		if cnt[ransomNote[i]-'a'] < 0 {
			return false // demand exceeded supply for this letter
		}
	}
	return true // all letters covered
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(bruteForce("a", "b"))    // expected false
	fmt.Println(bruteForce("aa", "ab"))  // expected false
	fmt.Println(bruteForce("aa", "aab")) // expected true

	fmt.Println("=== Approach 2: Hash Map Counting ===")
	fmt.Println(hashMap("a", "b"))    // expected false
	fmt.Println(hashMap("aa", "ab"))  // expected false
	fmt.Println(hashMap("aa", "aab")) // expected true

	fmt.Println("=== Approach 3: Fixed Array Counting (Optimal) ===")
	fmt.Println(arrayCount("a", "b"))    // expected false
	fmt.Println(arrayCount("aa", "ab"))  // expected false
	fmt.Println(arrayCount("aa", "aab")) // expected true
}
