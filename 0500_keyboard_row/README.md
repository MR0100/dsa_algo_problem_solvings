# 0500 — Keyboard Row

> LeetCode #500 · Difficulty: Easy
> **Categories:** Array, Hash Table, String

---

## Problem Statement

Given an array of strings `words`, return *the words that can be typed using letters of the alphabet on only one row of American keyboard like the image below*.

In the **American keyboard**:

- the first row consists of the characters `"qwertyuiop"`,
- the second row consists of the characters `"asdfghjkl"`, and
- the third row consists of the characters `"zxcvbnm"`.

**Example 1:**

```
Input: words = ["Hello","Alaska","Dad","Peace"]
Output: ["Alaska","Dad"]
```

**Explanation:** Both `"a"` and `"A"` are in the 2nd row of the American keyboard due to case insensitivity.

**Example 2:**

```
Input: words = ["omk"]
Output: []
```

**Example 3:**

```
Input: words = ["adsdf","sfd"]
Output: ["adsdf","sfd"]
```

**Constraints:**

- `1 <= words.length <= 20`
- `1 <= words[i].length <= 100`
- `words[i]` consists of English letters (both lowercase and uppercase).

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★☆☆☆ Low        | 2023          |
| Mathworks  | ★★☆☆☆ Low        | 2023          |
| Google     | ★☆☆☆☆ Rare       | 2022          |
| Bloomberg  | ★☆☆☆☆ Rare       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Table / Set Membership** — the crux is answering "which keyboard row does this character live on?" in O(1); either a per-row character set or a flat `letter → row` lookup table gives constant-time membership → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **String Processing** — case-folding each word and scanning its characters against a fixed alphabet partition → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Per-Word Row Scan with Sets | O(N·L) | O(1) | Most readable; mirrors the problem statement directly |
| 2 | Row-Index Lookup Table (Optimal) | O(N·L) | O(1) | Fewer constant-factor ops; one integer compare per char |

*(N = number of words, L = max word length; the alphabet is a fixed 26 letters, so the auxiliary structures are O(1).)*

---

## Approach 1 — Per-Word Row Scan with Sets

### Intuition

A word is typeable on one row **iff** all its letters share a row. So determine the row of the (lower-cased) first letter, then confirm every remaining letter belongs to that same row's character set. A per-row `set` gives O(1) membership tests. Keep the original word (with its original casing) when it passes.

### Algorithm

1. Build three character sets, one per keyboard row.
2. For each word:
   1. Lowercase it (case-insensitive keys).
   2. Find which row-set contains the first character → `target`.
   3. Scan every character; if one is **not** in `target`, reject.
3. Collect the words that pass (original spelling).

### Complexity

- **Time:** O(N·L) — each of the `N` words is scanned once, each character an O(1) set lookup.
- **Space:** O(1) — three fixed sets of ≤26 characters, independent of input size (excluding the output list).

### Code

```go
func bruteForce(words []string) []string {
	// map each character to a small set for O(1) membership
	rowSets := []map[rune]bool{{}, {}, {}}
	for i, row := range []string{row1, row2, row3} {
		for _, ch := range row {
			rowSets[i][ch] = true
		}
	}

	result := []string{}
	for _, word := range words {
		lower := strings.ToLower(word) // 'A' and 'a' share a key
		// pick the row that owns the first character
		var target map[rune]bool
		first := rune(lower[0])
		for _, set := range rowSets {
			if set[first] {
				target = set
				break
			}
		}
		// verify every character belongs to that same row
		ok := true
		for _, ch := range lower {
			if !target[ch] {
				ok = false // a letter from a different row → reject
				break
			}
		}
		if ok {
			result = append(result, word) // keep the ORIGINAL (preserve case)
		}
	}
	return result
}
```

### Dry Run

Example 1: `words = ["Hello","Alaska","Dad","Peace"]`. Rows: R1 `qwertyuiop`, R2 `asdfghjkl`, R3 `zxcvbnm`.

| word | lower | first char → row | scan result | kept? |
|------|-------|------------------|-------------|-------|
| Hello | hello | `h` → R2 | `e` ∈ R1, not R2 → reject | no |
| Alaska | alaska | `a` → R2 | a,l,a,s,k,a all ∈ R2 → pass | **yes** |
| Dad | dad | `d` → R2 | d,a,d all ∈ R2 → pass | **yes** |
| Peace | peace | `p` → R1 | `e`∈R1, `a`∈R2 → reject | no |

Result: `["Alaska", "Dad"]` ✔

---

## Approach 2 — Row-Index Lookup Table (Optimal)

### Intuition

Collapse the three sets into one flat array `rowOf[letter] = rowIndex`. A word qualifies exactly when every letter's `rowOf` equals the first letter's — a single integer comparison per character, branch-free and cache-friendly, with no map hashing.

### Algorithm

1. Fill `rowOf[c-'a'] = r` for each character `c` in row `r` (`r ∈ {0,1,2}`).
2. For each word:
   1. `r0 = rowOf[firstLetter]`.
   2. If any later letter has `rowOf != r0`, reject; else keep.

### Complexity

- **Time:** O(N·L) — one pass per word, O(1) array index per character.
- **Space:** O(1) — a fixed 26-slot integer table plus the output.

### Code

```go
func rowIndexTable(words []string) []string {
	var rowOf [26]int // rowOf[letter-'a'] = which keyboard row (0,1,2)
	for r, row := range []string{row1, row2, row3} {
		for _, ch := range row {
			rowOf[ch-'a'] = r // record this letter's row
		}
	}

	result := []string{}
	for _, word := range words {
		lower := strings.ToLower(word)
		r0 := rowOf[lower[0]-'a'] // the row every letter must match
		ok := true
		for i := 1; i < len(lower); i++ {
			if rowOf[lower[i]-'a'] != r0 { // a letter from another row?
				ok = false
				break
			}
		}
		if ok {
			result = append(result, word) // preserve original casing
		}
	}
	return result
}
```

### Dry Run

`rowOf` after construction (partial): `q,w,e,r,t,y,u,i,o,p → 0`; `a,s,d,f,g,h,j,k,l → 1`; `z,x,c,v,b,n,m → 2`.

Example 1: `words = ["Hello","Alaska","Dad","Peace"]`.

| word | lower | r0 = rowOf[first] | mismatch found? | kept? |
|------|-------|-------------------|-----------------|-------|
| Hello | hello | rowOf['h']=1 | rowOf['e']=0 ≠ 1 → stop | no |
| Alaska | alaska | rowOf['a']=1 | all of l,a,s,k,a = 1 | **yes** |
| Dad | dad | rowOf['d']=1 | a=1, d=1 | **yes** |
| Peace | peace | rowOf['p']=0 | rowOf['e']=0 ok, rowOf['a']=1 ≠ 0 → stop | no |

Result: `["Alaska", "Dad"]` ✔

---

## Key Takeaways

- **Reduce "grouping" questions to O(1) membership.** Mapping each character to its group (a row here) up front turns the check into a constant-time lookup, so the scan is linear regardless of alphabet.
- **A flat index table beats several sets** when the universe is small and fixed (26 letters): `rowOf[c]` is one array read versus a hash lookup, and the comparison `rowOf[c] == r0` is branch-friendly.
- **Case-fold once per word** (`strings.ToLower`) so `'A'` and `'a'` collapse to the same key — but **keep the original string** in the output to preserve the input's casing.
- Fixed-alphabet auxiliary structures are **O(1) space**: their size does not grow with the input, only the output list does.

---

## Related Problems

- LeetCode #383 — Ransom Note (character-count membership with a fixed alphabet)
- LeetCode #242 — Valid Anagram (26-slot count table over letters)
- LeetCode #205 — Isomorphic Strings (character → group consistency)
- LeetCode #1160 — Find Words That Can Be Formed by Characters (alphabet counting filter)
- LeetCode #290 — Word Pattern (bijection membership between two domains)
