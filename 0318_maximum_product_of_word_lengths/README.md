# 0318 — Maximum Product of Word Lengths

> LeetCode #318 · Difficulty: Medium
> **Categories:** Array, String, Bit Manipulation, Bitmask

---

## Problem Statement

Given a string array `words`, return the maximum value of
`length(word[i]) * length(word[j])` where the two words **do not share common
letters**. If no such two words exist, return `0`.

**Example 1:**

```
Input:  words = ["abcw","baz","foo","bar","xtfn","abcdef"]
Output: 16
```

Explanation: The two words can be `"abcw"`, `"xtfn"`.

**Example 2:**

```
Input:  words = ["a","ab","abc","d","cd","bcd","abcd"]
Output: 4
```

Explanation: The two words can be `"ab"`, `"cd"`.

**Example 3:**

```
Input:  words = ["a","aa","aaa","aaaa"]
Output: 0
```

Explanation: No such pair of words exists.

**Constraints:**

- `2 <= words.length <= 1000`
- `1 <= words[i].length <= 1000`
- `words[i]` consists only of lowercase English letters.

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Google    | ★★★★☆ High       | 2024          |
| Amazon    | ★★★☆☆ Medium     | 2023          |
| Microsoft | ★★★☆☆ Medium     | 2023          |
| Apple     | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Bit Manipulation (bitmask)** — encode each word's 26-letter set into an int
  so "share a letter?" is a single `AND` →
  see [`/dsa/bit_manipulation.md`](/dsa/bit_manipulation.md)
- **Hash Map / Set** — precompute per-word letter sets for O(26) intersection →
  see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **String Algorithms** — reducing a word to its distinct-letter signature →
  see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force pairwise char check | O(n²·L²) | O(1) | Tiny inputs; baseline |
| 2 | HashSet per word | O(ΣL + n²·26) | O(n·26) | When bit tricks feel obscure |
| 3 | Bitmask (Optimal) | O(ΣL + n²) | O(n) | Interview default — fastest, cleanest |

---

## Approach 1 — Brute Force

### Intuition
Two words qualify iff they share no letter. Test that directly for every pair by
scanning one word's characters against the other's, and track the largest length
product among qualifying pairs.

### Algorithm
1. For each pair `(i, j)` with `i < j`:
   - Scan `words[i] × words[j]` for a common character; if found, skip.
   - Otherwise compute `len(words[i]) * len(words[j])` and update the max.
2. Return the max (0 if no valid pair).

### Complexity
- **Time:** O(n²·L²) — `n²` pairs, each an O(L²) shared-letter scan.
- **Space:** O(1).

### Code
```go
func bruteForce(words []string) int {
	best := 0
	for i := 0; i < len(words); i++ {
		for j := i + 1; j < len(words); j++ {
			if !shareLetter(words[i], words[j]) {
				if p := len(words[i]) * len(words[j]); p > best {
					best = p
				}
			}
		}
	}
	return best
}

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
```

### Dry Run
Input `["abcw","baz","foo","bar","xtfn","abcdef"]`.

| Pair | words | share a letter? | product | best |
|------|-------|-----------------|---------|------|
| (0,1) | abcw / baz | yes ('a','b') | — | 0 |
| (0,4) | abcw / xtfn | no | 4·4 = 16 | 16 |
| (1,4) | baz / xtfn | no | 3·4 = 12 | 16 |
| ... | ... | ... | ≤ 16 | 16 |

Return `16`. ✓

---

## Approach 2 — HashSet Per Word

### Intuition
Building a distinct-letter set once per word turns the shared-letter test into a
bounded (≤ 26) membership scan instead of O(L²).

### Algorithm
1. For each word, build `sets[i]` = map of its distinct letters.
2. For each pair `(i, j)`: if no letter of `sets[i]` is in `sets[j]`, the pair is
   valid; update the max product.
3. Return the max.

### Complexity
- **Time:** O(ΣL + n²·26) — build sets in O(ΣL); each pair test ≤ 26 lookups.
- **Space:** O(n·26) — one set per word.

### Code
```go
func hashSet(words []string) int {
	sets := make([]map[byte]bool, len(words))
	for i, w := range words {
		s := make(map[byte]bool, 26)
		for k := 0; k < len(w); k++ {
			s[w[k]] = true
		}
		sets[i] = s
	}
	best := 0
	for i := 0; i < len(words); i++ {
		for j := i + 1; j < len(words); j++ {
			disjoint := true
			for c := range sets[i] {
				if sets[j][c] {
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
```

### Dry Run
Input `["abcw","baz","foo","bar","xtfn","abcdef"]`.

| Step | Detail |
|------|--------|
| sets[0] | {a,b,c,w} |
| sets[4] | {x,t,f,n} |
| pair (0,4) | scan {a,b,c,w} against sets[4] → none present → disjoint |
| product | 4·4 = 16 → best = 16 |
| other pairs | never exceed 16 |

Return `16`. ✓

---

## Approach 3 — Bitmask (Optimal)

### Intuition
Only the SET of distinct letters matters. Encode word `w` as a 26-bit integer:
bit `(c-'a')` is set when letter `c` appears. Then two words share no letter iff
`mask[i] & mask[j] == 0` — a single machine instruction. This makes the pair
loop O(1) per pair with a tiny constant.

### Algorithm
1. For each word, compute `masks[i]` = OR of `1 << (c-'a')` over its letters.
2. For each pair `(i, j)`: if `masks[i] & masks[j] == 0` they are disjoint;
   update the max product with `len[i]·len[j]`.
3. Return the max.

### Complexity
- **Time:** O(ΣL + n²) — build masks in O(ΣL); pair loop is O(1) per pair.
- **Space:** O(n) — one integer mask per word.

### Code
```go
func bitmask(words []string) int {
	masks := make([]int, len(words))
	for i, w := range words {
		m := 0
		for k := 0; k < len(w); k++ {
			m |= 1 << (w[k] - 'a')
		}
		masks[i] = m
	}
	best := 0
	for i := 0; i < len(words); i++ {
		for j := i + 1; j < len(words); j++ {
			if masks[i]&masks[j] == 0 {
				if p := len(words[i]) * len(words[j]); p > best {
					best = p
				}
			}
		}
	}
	return best
}
```

### Dry Run
Input `["abcw","baz","foo","bar","xtfn","abcdef"]`. Bits: a=0, b=1, c=2, ...

| word | letters | mask (binary, low bits) |
|------|---------|-------------------------|
| abcw | a,b,c,w | ...1 at bits 0,1,2,22 |
| xtfn | x,t,f,n | bits 5,13,19,23 |

`masks[0] & masks[4]`: no overlapping bits → `0` → disjoint → product 4·4 = 16.
No other pair beats 16, so return `16`. ✓

---

## Key Takeaways

- **Bitmask a small alphabet.** A set over 26 lowercase letters fits in an int;
  set membership → bit test, intersection → `AND`, union → `OR`.
- **Disjoint sets ⇔ `a & b == 0`.** This one-liner replaces an O(L²) or O(26)
  inner check and is the whole trick behind the optimal solution.
- **Reduce, then compare.** Precompute a compact signature (mask / set) per item
  once, so the O(n²) pair phase stays cheap.

---

## Related Problems

- LeetCode #187 — Repeated DNA Sequences (encode substrings into ints)
- LeetCode #421 — Maximum XOR of Two Numbers in an Array (bit tricks on pairs)
- LeetCode #526 — Beautiful Arrangement (bitmask over choices)
- LeetCode #1178 — Number of Valid Words for Each Puzzle (letter-set bitmasks)
