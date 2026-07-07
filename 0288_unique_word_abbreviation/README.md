# 0288 — Unique Word Abbreviation

> LeetCode #288 · Difficulty: Medium
> **Categories:** Hash Table, String, Design

---

## Problem Statement

The **abbreviation** of a word is a concatenation of its first letter, the number of characters between the first and last letter, and its last letter. If a word has only two characters, then it is an abbreviation of itself.

For example:

- `dog` → `d1g` because there is one letter between the first letter `'d'` and the last letter `'g'`.
- `internationalization` → `i18n` because there are 18 letters between the first letter `'i'` and the last letter `'n'`.
- `it` → `it` because any word with only two characters is an abbreviation of itself.

Implement the `ValidWordAbbr` class:

- `ValidWordAbbr(String[] dictionary)` Initializes the object with a `dictionary` of words.
- `boolean isUnique(String word)` Returns `true` if **either** of the following conditions are met (otherwise returns `false`):
  - There is no word in `dictionary` whose abbreviation is equal to `word`'s abbreviation.
  - For any word in `dictionary` whose abbreviation is equal to `word`'s abbreviation, that word and `word` are **the same**.

**Example 1:**

```
Input:
["ValidWordAbbr", "isUnique", "isUnique", "isUnique", "isUnique", "isUnique"]
[[["deer", "door", "cake", "card"]], ["dear"], ["cart"], ["cane"], ["make"], ["cake"]]
Output:
[null, false, true, false, true, true]

Explanation:
ValidWordAbbr validWordAbbr = new ValidWordAbbr(["deer", "door", "cake", "card"]);
validWordAbbr.isUnique("dear"); // return false, dictionary word "deer" and word "dear" have the same abbreviation "d2r" but are not the same.
validWordAbbr.isUnique("cart"); // return true, no words in the dictionary have the abbreviation "c2t".
validWordAbbr.isUnique("cane"); // return false, dictionary word "cake" and word "cane" have the same abbreviation "c2e" but are not the same.
validWordAbbr.isUnique("make"); // return true, no words in the dictionary have the abbreviation "m2e".
validWordAbbr.isUnique("cake"); // return true, because "cake" is already in the dictionary and no other word in the dictionary has the abbreviation "c2e".
```

**Constraints:**

- `1 <= dictionary.length <= 3 * 10^4`
- `1 <= dictionary[i].length <= 20`
- `dictionary[i]` consists of lowercase English letters.
- `1 <= word.length <= 20`
- `word` consists of lowercase English letters.
- At most `5000` calls will be made to `isUnique`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★☆☆☆ Low        | 2022          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Map** — map each abbreviation to the set of distinct dictionary words producing it, giving O(1) `isUnique` lookups → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Design Data Structure** — build once, answer many queries; store just enough state (abbr → word set) to answer the bijection question → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)
- **String Encoding** — collapsing a word to `first + count + last` is a canonical fingerprint → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Abbreviation → set of words (Hash Map) | Build O(N·L), query O(L) | O(N·L) | The standard and optimal design |

> A single well-chosen structure solves this; the subtlety is storing the **set of distinct words** per abbreviation (not a single word or a count), which is what makes the "same word may appear multiple times" case correct.

---

## Approach 1 — Abbreviation → Set of Words (Hash Map)

### Intuition

`isUnique(word)` asks: among dictionary words, is `word` the *only* one whose abbreviation matches? So the state we need is, for each abbreviation, the set of distinct dictionary words that map to it. Then a query is unique when that set is empty (abbreviation unseen) or is exactly `{word}` (only `word`, even if it appeared several times in the dictionary, owns the abbreviation).

A tempting shortcut — map abbreviation → single word, and mark "conflicted" if a different word collides — also works, but storing the **distinct-word set** is the clearest correct model and handles duplicate dictionary entries naturally.

### Algorithm

1. **Build:** for each dictionary word, compute its abbreviation and insert the word into `abbrToWords[abbr]` (a set, so identical words dedupe).
2. **Query `isUnique(word)`:**
   1. `a = abbreviate(word)`; look up `words = abbrToWords[a]`.
   2. If `words` is empty → return `true`.
   3. If `words == {word}` (size 1 and contains `word`) → return `true`.
   4. Otherwise → return `false`.

The abbreviation itself is `first letter + (length-2) + last letter`, except words of length ≤ 2 abbreviate to themselves.

### Complexity

- **Time:** Build O(N·L) over `N` words of average length `L`; each `isUnique` is O(L) to abbreviate plus O(1) map/set operations.
- **Space:** O(N·L) — the map of abbreviations to word sets.

### Code

```go
func abbreviate(word string) string {
	n := len(word)
	if n <= 2 {
		return word // too short to compress; itself is the "abbreviation"
	}
	// first letter + (middle length as a number) + last letter, e.g. i18n.
	return fmt.Sprintf("%c%d%c", word[0], n-2, word[n-1])
}

type ValidWordAbbr struct {
	// abbr → set of distinct words producing that abbreviation.
	abbrToWords map[string]map[string]struct{}
}

func Constructor(dictionary []string) ValidWordAbbr {
	m := make(map[string]map[string]struct{})
	for _, w := range dictionary {
		a := abbreviate(w)
		if m[a] == nil {
			m[a] = make(map[string]struct{})
		}
		m[a][w] = struct{}{} // dedupe: identical dictionary words collapse
	}
	return ValidWordAbbr{abbrToWords: m}
}

func (v *ValidWordAbbr) isUnique(word string) bool {
	a := abbreviate(word)
	words, ok := v.abbrToWords[a]
	if !ok || len(words) == 0 {
		return true // abbreviation never seen → unique
	}
	// Unique iff the ONLY word with this abbreviation is `word` itself.
	if _, contains := words[word]; contains && len(words) == 1 {
		return true
	}
	return false
}
```

### Dry Run

Build from `["deer","door","cake","card"]`:

| word | abbreviation | map state |
|------|--------------|-----------|
| deer | d2r | {d2r:{deer}} |
| door | d2r | {d2r:{deer,door}} |
| cake | c2e | {d2r:{deer,door}, c2e:{cake}} |
| card | c2d | {…, c2d:{card}} |

Queries:

| word | abbr | set at abbr | unique? | reason |
|------|------|-------------|---------|--------|
| dear | d2r | {deer,door} | **false** | set not empty and not {dear} |
| cart | c2t | (none) | **true** | abbreviation unseen |
| cane | c2e | {cake} | **false** | owner is "cake", not "cane" |
| make | m2e | (none) | **true** | abbreviation unseen |
| cake | c2e | {cake} | **true** | size 1 and contains "cake" |

Output: `false, true, false, true, true` ✔

---

## Key Takeaways

- **Precompute a canonical key, then let a hash map do the work.** The abbreviation is a fingerprint; grouping dictionary words by it turns every query into an O(1) lookup.
- **Store the distinct-word SET per key, not a single word or a count.** That is exactly the state the bijection question needs and it cleanly handles duplicate dictionary entries (`["hello","hello"]` still leaves `"hello"` unique).
- Remember the length-≤-2 edge case: such words abbreviate to themselves, so `"it"` stays `"it"`.
- Classic "design" pattern: heavy build once, cheap repeated queries — invest the preprocessing to make `isUnique` trivial.

---

## Related Problems

- LeetCode #408 — Valid Word Abbreviation (verify an abbreviation matches a word)
- LeetCode #411 — Minimum Unique Word Abbreviation (generate abbreviations)
- LeetCode #49 — Group Anagrams (canonical key → hash-map grouping)
- LeetCode #205 — Isomorphic Strings (bijection via mapping)
- LeetCode #290 — Word Pattern (bijection check)
