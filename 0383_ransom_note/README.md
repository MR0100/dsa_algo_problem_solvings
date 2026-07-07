# 0383 — Ransom Note

> LeetCode #383 · Difficulty: Easy
> **Categories:** Hash Table, String, Counting

---

## Problem Statement

Given two strings `ransomNote` and `magazine`, return `true` if `ransomNote` can be constructed by using the letters from `magazine` and `false` otherwise.

Each letter in `magazine` can only be used once in `ransomNote`.

**Example 1:**
```
Input: ransomNote = "a", magazine = "b"
Output: false
```

**Example 2:**
```
Input: ransomNote = "aa", magazine = "ab"
Output: false
```

**Example 3:**
```
Input: ransomNote = "aa", magazine = "aab"
Output: true
```

**Constraints:**
- `1 <= ransomNote.length, magazine.length <= 10⁵`
- `ransomNote` and `magazine` consist of lowercase English letters.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Frequency Counting / Hash Map** — reduce the strings to letter multisets and compare supply vs. demand → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **String Handling** — byte-level iteration over the fixed `'a'..'z'` alphabet, enabling an array in place of a map → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

Let n = len(ransomNote), m = len(magazine).

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (consume from copy) | O(n·m) | O(m) | Baseline; tiny inputs only |
| 2 | Hash Map Counting | O(n + m) | O(k) distinct letters | Arbitrary/large alphabets (unicode) |
| 3 | Fixed Array Counting (Optimal) | O(n + m) | O(1) | Known small alphabet — the interview answer |

---

## Approach 1 — Brute Force

### Intuition
Each note letter must pair with a *distinct* magazine letter. Model "distinct" by physically crossing a letter off a mutable copy of the magazine as soon as it is used, so it cannot be reused. If any note letter finds no free match, construction is impossible.

### Algorithm
1. Copy `magazine` into a mutable byte slice.
2. For each character `c` of `ransomNote`, scan the copy for an unused `c`.
3. If found, blank that slot (mark used) and continue.
4. If a scan finds nothing, return `false`.
5. If every letter is matched, return `true`.

### Complexity
- **Time:** O(n·m) — for each of n note chars we may scan all m magazine chars.
- **Space:** O(m) — the mutable copy of the magazine.

### Code
```go
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
```

### Dry Run
`ransomNote = "aa"`, `magazine = "aab"`, `mag = ['a','a','b']`.

| Step | note char | scan of mag | action | mag after |
|------|-----------|-------------|--------|-----------|
| 1 | `a` (i=0) | j=0 matches | blank slot 0 | `[0,'a','b']` |
| 2 | `a` (i=1) | j=0 is 0, j=1 matches | blank slot 1 | `[0,0,'b']` |
| end | — | note exhausted | return `true` | — |

Result: `true`.

---

## Approach 2 — Hash Map Counting

### Intuition
Only the multiset of letters matters, not order. Tally the magazine's letters once, then each note letter needs one unit of that letter's remaining stock. Two linear passes replace the nested scan.

### Algorithm
1. Build `count[c]` = number of times `c` appears in `magazine`.
2. For each note char `c`: if `count[c] == 0` → `false`; else `count[c]--`.
3. Survive the whole note → `true`.

### Complexity
- **Time:** O(n + m) — one pass to count, one to spend.
- **Space:** O(k) — k distinct letters (≤ 26, effectively O(1)).

### Code
```go
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
```

### Dry Run
`ransomNote = "aa"`, `magazine = "aab"`.

| Step | action | count map |
|------|--------|-----------|
| count `a` | ++ | `{a:1}` |
| count `a` | ++ | `{a:2}` |
| count `b` | ++ | `{a:2, b:1}` |
| note `a` | count[a]=2>0, -- | `{a:1, b:1}` |
| note `a` | count[a]=1>0, -- | `{a:0, b:1}` |
| end | return `true` | — |

Result: `true`.

---

## Approach 3 — Fixed Array Counting (Optimal)

### Intuition
Same counting idea, but the alphabet is exactly `'a'..'z'`, so an index `c-'a'` into a length-26 array replaces every hash operation: no hashing, perfect cache locality, least code. Early exit: if the note is longer than the magazine it is impossible.

### Algorithm
1. If `len(note) > len(magazine)` → `false` immediately.
2. `cnt[c-'a']++` for each magazine char.
3. For each note char, `cnt[c-'a']--`; if it drops below 0 → `false`.
4. Otherwise → `true`.

### Complexity
- **Time:** O(n + m) — two linear passes, O(1) work each.
- **Space:** O(1) — a fixed 26-integer array.

### Code
```go
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
```

### Dry Run
`ransomNote = "aa"`, `magazine = "aab"`. Track only slot `a` (index 0).

| Step | action | cnt[a] |
|------|--------|--------|
| length check | 2 > 3? no | — |
| count `a` | ++ | 1 |
| count `a` | ++ | 2 |
| count `b` | (slot b) | 2 |
| note `a` | -- → 1 ≥ 0 | 1 |
| note `a` | -- → 0 ≥ 0 | 0 |
| end | return `true` | — |

Result: `true`.

---

## Key Takeaways
- "Can string A be built from the letters of string B" is a **multiset containment** check — count B, decrement per A char.
- A fixed-size array beats a hash map whenever the alphabet is small and known; index arithmetic replaces hashing.
- The length pre-check (`|note| > |magazine|`) is a free early exit.
- Decrement-then-check (going below 0) fuses the "has stock?" and "spend" steps into one.

---

## Related Problems
- LeetCode #242 — Valid Anagram (equality of two letter-count arrays)
- LeetCode #387 — First Unique Character in a String (frequency array)
- LeetCode #691 — Stickers to Spell Word (multiset coverage, harder)
- LeetCode #49 — Group Anagrams (letter-count signature as a key)
