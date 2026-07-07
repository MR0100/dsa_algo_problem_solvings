# 0290 — Word Pattern

> LeetCode #290 · Difficulty: Easy
> **Categories:** Hash Table, String

---

## Problem Statement

Given a `pattern` and a string `s`, find if `s` follows the same pattern.

Here **follow** means a full match, such that there is a **bijection** between a letter in `pattern` and a **non-empty** word in `s`. Specifically:

- Each letter in `pattern` maps to exactly one unique word in `s`.
- Each unique word in `s` maps to exactly one letter in `pattern`.
- No two letters map to the same word, and no two words map to the same letter.

**Example 1:**

```
Input: pattern = "abba", s = "dog cat cat dog"
Output: true

Explanation:
The bijection can be established as:
'a' maps to "dog".
'b' maps to "cat".
```

**Example 2:**

```
Input: pattern = "abba", s = "dog cat cat fish"
Output: false
```

**Example 3:**

```
Input: pattern = "aaaa", s = "dog cat cat dog"
Output: false
```

**Constraints:**

- `1 <= pattern.length <= 300`
- `pattern` contains only lower-case English letters.
- `1 <= s.length <= 3000`
- `s` contains only lowercase English letters and spaces `' '`.
- `s` does not contain any leading or trailing spaces.
- All the words in `s` are separated by a **single space**.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |
| Apple      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Map (Bijection)** — enforcing a one-to-one mapping needs BOTH directions (letter → word and word → letter), so two maps or a paired index map are used → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **String Tokenization** — split `s` on spaces into words to line up 1-to-1 with pattern letters → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Two hash maps (Optimal) | O(n) | O(k) | The standard answer; explicit bijection in both directions |
| 2 | Single first-seen index map | O(n) | O(k) | Elegant one-pass trick using first-occurrence indices |

(`n` = total characters in pattern + s; `k` = number of distinct letters/words.)

---

## Approach 1 — Two Hash Maps (Bijection Check)

### Intuition

A bijection must be consistent in **both** directions. Tracking only `letter → word` fails on `pattern="ab", s="dog dog"`: `'a'→"dog"` and `'b'→"dog"` each look fine individually, yet two letters share one word. Guard against this by also maintaining `word → letter` and rejecting any conflict either way.

### Algorithm

1. Split `s` into words; if the count differs from `len(pattern)`, return `false`.
2. For each aligned `(letter, word)`:
   - If `letter` is already mapped, it must map to this `word`.
   - If `word` is already mapped, it must map to this `letter`.
   - Otherwise, record both directions.
3. If no conflict occurs, return `true`.

### Complexity

- **Time:** O(n) — one pass; each map op is O(1) amortized (string keys hashed by length `≤ L`).
- **Space:** O(k) — the two maps over distinct letters/words.

### Code

```go
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
```

### Dry Run

Example 1: `pattern="abba"`, `s="dog cat cat dog"` → words = `[dog, cat, cat, dog]`.

| i | c | w | letterToWord check | wordToLetter check | maps after |
|---|---|---|--------------------|--------------------|------------|
| 0 | a | dog | unseen → add a→dog | unseen → add dog→a | {a:dog},{dog:a} |
| 1 | b | cat | unseen → add b→cat | unseen → add cat→b | +{b:cat},{cat:b} |
| 2 | b | cat | b→cat matches | cat→b matches | unchanged |
| 3 | a | dog | a→dog matches | dog→a matches | unchanged |

No conflict → `true` ✔

Counter-check `pattern="abba", s="dog dog dog dog"`: at `i=1`, `b` is new but `wordToLetter["dog"]` already = `a` ≠ `b` → returns `false` (correctly rejects the two-letters-one-word collision).

---

## Approach 2 — Single First-Seen Index Map

### Intuition

Two sequences share a pattern iff, at every position, "when did I last see this token?" matches for both. Concretely, if the **first occurrence** of `pattern[i]` sits at the same index as the first occurrence of `words[i]`, the two run in lockstep — a bijection. Record each token's first-seen index; the two must always agree.

### Algorithm

1. Split `s`; its length must equal `len(pattern)`.
2. Keep `letterIdx` (letter → first index) and `wordIdx` (word → first index).
3. At position `i`: the letter and the word must be "new together" or point to the same earlier index. Any mismatch breaks the bijection; then record `i` for any unseen token.

### Complexity

- **Time:** O(n) — single pass.
- **Space:** O(k) — first-index maps over distinct tokens.

### Code

```go
func indexMap(pattern string, s string) bool {
	words := strings.Fields(s)
	if len(words) != len(pattern) {
		return false
	}
	letterIdx := make(map[byte]int)   // letter → first index it appeared
	wordIdx := make(map[string]int)   // word → first index it appeared
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
```

### Dry Run

Example 1: `pattern="abba"`, words = `[dog, cat, cat, dog]`.

| i | c | w | letterIdx[c] | wordIdx[w] | agree? | record |
|---|---|---|--------------|------------|--------|--------|
| 0 | a | dog | (new) | (new) | both new ✓ | a→0, dog→0 |
| 1 | b | cat | (new) | (new) | both new ✓ | b→1, cat→1 |
| 2 | b | cat | 1 | 1 | 1 == 1 ✓ | — |
| 3 | a | dog | 0 | 0 | 0 == 0 ✓ | — |

All positions agree → `true` ✔

Counter-check `pattern="aaaa", s="dog cat cat dog"`: at `i=1`, `a` was seen (idx 0) but `cat` is new → `lok != wok` → `false` (correct).

---

## Key Takeaways

- **Bijection = check both directions.** A single forward map (`letter → word`) is not enough; add the reverse map (`word → letter`) or the collision `pattern="ab", s="dog dog"` slips through.
- **First-occurrence index is a compact bijection encoding.** Two token streams match iff each token's first-seen index aligns — a one-map-per-side trick that also solves Isomorphic Strings.
- Guard the **length mismatch** early: unequal counts of letters and words can never form a bijection.
- Use `strings.Fields` (not `strings.Split(s, " ")`) to robustly tokenize on whitespace.

---

## Related Problems

- LeetCode #205 — Isomorphic Strings (same bijection idea over characters)
- LeetCode #291 — Word Pattern II (backtracking to discover the mapping)
- LeetCode #890 — Find and Replace Pattern (bijection over multiple words)
- LeetCode #49 — Group Anagrams (canonical key grouping)
- LeetCode #288 — Unique Word Abbreviation (mapping / hashing design)
