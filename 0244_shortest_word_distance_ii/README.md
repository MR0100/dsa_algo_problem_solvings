# 0244 — Shortest Word Distance II

> LeetCode #244 · Difficulty: Medium · 🔒 Premium
> **Categories:** Design, Array, Hash Table, Two Pointers, String

---

## Problem Statement

Design a data structure that will be initialized with a string array, and then it should answer queries of the shortest distance between two different strings from the array.

Implement the `WordDistance` class:

- `WordDistance(String[] wordsDict)` initializes the object with the strings array `wordsDict`.
- `int shortest(String word1, String word2)` returns the shortest distance between `word1` and `word2` in the array `wordsDict`.

**Example 1:**

```
Input
["WordDistance", "shortest", "shortest"]
[[["practice", "makes", "perfect", "coding", "makes"]], ["coding", "practice"], ["makes", "coding"]]
Output
[null, 3, 1]

Explanation
WordDistance wordDistance = new WordDistance(["practice", "makes", "perfect", "coding", "makes"]);
wordDistance.shortest("coding", "practice"); // return 3
wordDistance.shortest("makes", "coding");    // return 1
```

**Constraints:**

- `1 <= wordsDict.length <= 3 * 10^4`
- `1 <= wordsDict[i].length <= 10`
- `wordsDict[i]` consists of lowercase English letters.
- `word1` and `word2` are in `wordsDict`.
- `word1 != word2`
- At most `5000` calls will be made to `shortest`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| LinkedIn   | ★★★★★ Very High  | 2023          |
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Google     | ★★☆☆☆ Low        | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Design Data Structures** — precompute once in the constructor, answer many queries fast → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)
- **Hash Map** — map each word to its sorted list of occurrence indices → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Two Pointers (merge of two sorted lists)** — advance the smaller index to find the closest pair between two sorted index lists → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)

---

## Approaches Overview

| # | Approach | Init Time | Query Time | Space | When to use |
|---|----------|-----------|------------|-------|-------------|
| 1 | Index Map + Merge (Optimal) | O(n) | O(a + b) | O(n) | Many queries — the intended solution |
| 2 | Store Array + Rescan Per Query | O(1) | O(n) | O(n) | Few queries; simplest to code |

(n = length of `wordsDict`; a, b = occurrence counts of the two queried words.)

---

## Approach 1 — Index Map + Merge (Optimal)

### Intuition
Because we scan left to right, each word's occurrence indices are recorded in increasing order — already sorted. A `shortest` query then reduces to: given two sorted integer lists, find the smallest absolute difference between an element of one and an element of the other. Walk both with two pointers, always advancing the pointer at the smaller value (the only move that can shrink the gap).

### Algorithm
**Constructor:**
1. For each position `i`, append `i` to `wordIndex[words[i]]`.

**shortest(word1, word2):**
1. Let `l1 = wordIndex[word1]`, `l2 = wordIndex[word2]`, `p1 = p2 = 0`.
2. While both pointers in range: update `min` with `|l1[p1] - l2[p2]|`, then advance the pointer whose value is smaller.
3. Return `min`.

### Complexity
- **Time:** Constructor O(n); each query O(a + b), bounded by O(n).
- **Space:** O(n) — each index stored exactly once across all lists.

### Code
```go
type WordDistance struct {
	wordIndex map[string][]int
}

func Constructor(wordsDict []string) WordDistance {
	idx := make(map[string][]int, len(wordsDict))
	for i, w := range wordsDict {
		idx[w] = append(idx[w], i)
	}
	return WordDistance{wordIndex: idx}
}

func (wd *WordDistance) shortest(word1 string, word2 string) int {
	l1 := wordIndex(wd, word1)
	l2 := wordIndex(wd, word2)

	p1, p2 := 0, 0
	const maxInt = int(^uint(0) >> 1)
	min := maxInt
	for p1 < len(l1) && p2 < len(l2) {
		a, b := l1[p1], l2[p2]
		if d := abs(a - b); d < min {
			min = d
		}
		if a < b {
			p1++
		} else {
			p2++
		}
	}
	return min
}
```

### Dry Run
Constructor on `["practice","makes","perfect","coding","makes"]` builds:

```
practice → [0]
makes    → [1, 4]
perfect  → [2]
coding   → [3]
```

Query `shortest("coding","practice")`: `l1 = [3]`, `l2 = [0]`.

| p1 | p2 | a=l1[p1] | b=l2[p2] | \|a-b\| | min | move |
|----|----|----------|----------|---------|-----|------|
| 0 | 0 | 3 | 0 | 3 | 3 | a>b → p2++ |
| 0 | 1 | — | out of range | — | 3 | loop ends |

Result: `3`. ✓

Query `shortest("makes","coding")`: `l1 = [1,4]`, `l2 = [3]`.

| p1 | p2 | a | b | \|a-b\| | min | move |
|----|----|---|---|---------|-----|------|
| 0 | 0 | 1 | 3 | 2 | 2 | a<b → p1++ |
| 1 | 0 | 4 | 3 | 1 | 1 | a>b → p2++ |
| 1 | 1 | 4 | out | — | 1 | loop ends |

Result: `1`. ✓

---

## Approach 2 — Store Array + Rescan Per Query (Baseline)

### Intuition
Skip the precomputation and instead re-run problem 243's single-pass, last-seen-index scan on every query. Simple, but every query costs a full sweep of the array regardless of how rare the two words are.

### Algorithm
1. Store the original `wordsDict`.
2. On each query, scan tracking `i1`/`i2` (last-seen indices of `word1`/`word2`); whenever both are known, update the minimum gap.
3. Return the minimum.

### Complexity
- **Time:** O(n) per query — no matter what, the whole array is scanned.
- **Space:** O(n) — keep the original array.

### Code
```go
type WordDistanceRescan struct {
	words []string
}

func NewRescan(wordsDict []string) WordDistanceRescan {
	return WordDistanceRescan{words: wordsDict}
}

func (wd *WordDistanceRescan) shortest(word1 string, word2 string) int {
	i1, i2 := -1, -1
	min := len(wd.words)
	for k, w := range wd.words {
		switch w {
		case word1:
			i1 = k
			if i2 != -1 && i1-i2 < min {
				min = i1 - i2
			}
		case word2:
			i2 = k
			if i1 != -1 && i2-i1 < min {
				min = i2 - i1
			}
		}
	}
	return min
}
```

### Dry Run
Query `shortest("makes","coding")` on `["practice","makes","perfect","coding","makes"]`:

| k | w | i1 (makes) | i2 (coding) | candidate | min |
|---|---|------------|-------------|-----------|-----|
| 0 | practice | -1 | -1 | — | 5 |
| 1 | makes | 1 | -1 | i2=-1 → skip | 5 |
| 2 | perfect | 1 | -1 | — | 5 |
| 3 | coding | 1 | 3 | i2-i1 = 2 | **2** |
| 4 | makes | 4 | 3 | i1-i2 = 1 | **1** |

Result: `1`. ✓

---

## Key Takeaways
- **Design pattern: precompute in the constructor to speed up repeated queries.** With up to 5000 `shortest` calls, per-query O(n) rescans are wasteful — an index map pays off.
- Occurrence indices come out **already sorted** from a left-to-right pass; exploit that instead of sorting again.
- Finding the closest pair across two sorted lists is a **two-pointer merge**: always advance the smaller value.
- This is the multi-query evolution of #243 and the base for #245 (allowing `word1 == word2`).

---

## Related Problems
- LeetCode #243 — Shortest Word Distance (single query)
- LeetCode #245 — Shortest Word Distance III (word1 may equal word2)
- LeetCode #170 — Two Sum III – Data structure design (precompute vs per-query trade-off)
- LeetCode #348 — Design Tic-Tac-Toe (design; precompute state)
