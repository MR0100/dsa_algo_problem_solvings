# 0387 — First Unique Character in a String

> LeetCode #387 · Difficulty: Easy
> **Categories:** Hash Table, String, Queue, Counting

---

## Problem Statement

Given a string `s`, find the first non-repeating character in it and return its index. If it does not exist, return `-1`.

**Example 1:**

```
Input: s = "leetcode"
Output: 0
Explanation: The character 'l' at index 0 is the first character that does not occur at any other index.
```

**Example 2:**

```
Input: s = "loveleetcode"
Output: 2
```

**Example 3:**

```
Input: s = "aabb"
Output: -1
```

**Constraints:**

- `1 <= s.length <= 10^5`
- `s` consists of only lowercase English letters.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★★ Very High  | 2024          |
| Microsoft  | ★★★★☆ High       | 2024          |
| Bloomberg  | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Map / Counting** — count each character's frequency, then the first character with count 1 is the answer → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **String Processing** — a straightforward left-to-right scan over the characters → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Scan-and-Compare) | O(n²) | O(1) | Tiny inputs; no extra memory allowed |
| 2 | Hash Map Frequency (Two Pass) | O(n) | O(k) | General strings, any alphabet |
| 3 | Fixed 26-Bucket Array (Optimal) | O(n) | O(1) | Lowercase-only constraint; fastest |

---

## Approach 1 — Brute Force (Scan-and-Compare)

### Intuition

"First non-repeating character" means: for each position, does this character appear anywhere else? The first index whose character has no twin is the answer. Checking each index needs a full pass over the string, giving `O(n²)` but zero extra storage.

### Algorithm

1. For `i` from `0` to `len-1`:
   1. Scan `j` over the whole string.
   2. If `i != j` and `s[j] == s[i]`, mark `i` as duplicated and stop scanning.
   3. If no duplicate was found, return `i`.
2. Return `-1` if every character repeats.

### Complexity

- **Time:** O(n²) — for each of n indices we may scan up to n characters.
- **Space:** O(1) — only a boolean flag.

### Code

```go
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
```

### Dry Run

Input `s = "leetcode"`:

| i | s[i] | inner scan finds duplicate? | action |
|---|------|-----------------------------|--------|
| 0 | 'l' | no other 'l' | `unique` stays true → **return 0** |

Answer: `0`.

---

## Approach 2 — Hash Map Frequency (Two Pass)

### Intuition

A character is "unique" iff it occurs exactly once. Count every character first (pass 1), then walk the string left-to-right and return the first index whose character has count `1` (pass 2). The second pass preserves original order.

### Algorithm

1. Build `count[c]` = number of occurrences of `c` in one pass.
2. Scan `i` left-to-right; return the first `i` with `count[s[i]] == 1`.
3. Return `-1` if none.

### Complexity

- **Time:** O(n) — two linear passes.
- **Space:** O(k) — `k` distinct characters (≤ 26 here → effectively O(1)).

### Code

```go
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
```

### Dry Run

Input `s = "leetcode"`. Pass 1 counts: `l:1, e:3, t:1, c:1, o:1, d:1`.

| i | s[i] | count[s[i]] | action |
|---|------|-------------|--------|
| 0 | 'l' | 1 | **return 0** |

Answer: `0`. (For `"loveleetcode"`: `l:2,o:2,v:1,e:4,...` → first count-1 is `'v'` at index 2.)

---

## Approach 3 — Fixed 26-Bucket Array (Optimal)

### Intuition

Same two-pass idea, but since the alphabet is a fixed constant set (`a`–`z`), index a plain `[26]int` array by `c - 'a'` instead of hashing. No hash overhead, and space is strictly `O(1)` — 26 counters regardless of input length.

### Algorithm

1. `count[s[i]-'a']++` for every character (pass 1).
2. Return the first `i` with `count[s[i]-'a'] == 1` (pass 2).
3. Return `-1` otherwise.

### Complexity

- **Time:** O(n) — two linear passes.
- **Space:** O(1) — exactly 26 integer counters.

### Code

```go
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
```

### Dry Run

Input `s = "leetcode"`. Pass 1 fills `count`: index `'l'-'a'=11`→1, `'e'-'a'=4`→3, `'t'`→1, `'c'`→1, `'o'`→1, `'d'`→1.

| i | s[i] | count[s[i]-'a'] | action |
|---|------|-----------------|--------|
| 0 | 'l' | 1 | **return 0** |

Answer: `0`.

---

## Key Takeaways

- **Count then re-scan for order.** Two-pass counting is the canonical pattern for "first/only element with property P": pass 1 aggregates, pass 2 preserves original order.
- **Fixed alphabet ⇒ array beats map.** When the character set is a known constant (`a`–`z`, ASCII), a fixed-size array indexed by `c-'a'` is both faster and truly O(1) space.
- A one-pass "queue of candidates" variant also exists (LeetCode #451/#387 follow-ups) but two passes is simplest and optimal here.

---

## Related Problems

- LeetCode #451 — Sort Characters By Frequency (same counting primitive)
- LeetCode #383 — Ransom Note (fixed 26-bucket counting)
- LeetCode #242 — Valid Anagram (character frequency comparison)
- LeetCode #205 — Isomorphic Strings (character mapping)
