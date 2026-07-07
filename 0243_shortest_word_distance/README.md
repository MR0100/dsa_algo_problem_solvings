# 0243 — Shortest Word Distance

> LeetCode #243 · Difficulty: Easy · 🔒 Premium
> **Categories:** Array, String, Two Pointers

---

## Problem Statement

Given an array of strings `wordsDict` and two different strings that already exist in the array `word1` and `word2`, return *the shortest distance between these two words in the list*.

**Example 1:**

```
Input: wordsDict = ["practice", "makes", "perfect", "coding", "makes"], word1 = "coding", word2 = "practice"
Output: 3
```

**Example 2:**

```
Input: wordsDict = ["practice", "makes", "perfect", "coding", "makes"], word1 = "makes", word2 = "coding"
Output: 1
```

**Constraints:**

- `2 <= wordsDict.length <= 3 * 10^4`
- `1 <= wordsDict[i].length <= 10`
- `wordsDict[i]` consists of lowercase English letters.
- `word1` and `word2` are in `wordsDict`.
- `word1 != word2`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| LinkedIn   | ★★★★☆ High       | 2023          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Two Pointers (last-seen indices)** — track the most recent position of each word and measure the gap in a single sweep → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **Array Traversal** — a linear scan over the word list → see [`/dsa/arrays.md`](/dsa/arrays.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (all pairs) | O(n²) | O(1) | Tiny inputs; easiest to reason about |
| 2 | One-Pass Two Pointers (Optimal) | O(n) | O(1) | Always preferred — single scan |

(n = length of `wordsDict`.)

---

## Approach 1 — Brute Force

### Intuition
The answer is `|i - j|` for some index `i` of `word1` and some index `j` of `word2`. Enumerate every such pair and keep the minimum gap.

### Algorithm
1. Initialize `min` to a large sentinel (`len(words)` is a safe upper bound).
2. For each `i` where `words[i] == word1`:
   - For each `j` where `words[j] == word2`, update `min` with `|i - j|`.
3. Return `min`.

### Complexity
- **Time:** O(n²) — nested scan over matching positions in the worst case.
- **Space:** O(1) — only a running minimum.

### Code
```go
func bruteForce(words []string, word1 string, word2 string) int {
	min := len(words)
	for i := 0; i < len(words); i++ {
		if words[i] != word1 {
			continue
		}
		for j := 0; j < len(words); j++ {
			if words[j] == word2 {
				if d := abs(i - j); d < min {
					min = d
				}
			}
		}
	}
	return min
}
```

### Dry Run
Trace `words = ["practice","makes","perfect","coding","makes"]`, `word1="coding"`, `word2="practice"`:

| i (word1="coding") | j (word2="practice") | \|i-j\| | min |
|--------------------|----------------------|---------|-----|
| 3 | 0 | 3 | 3 |

Only one occurrence of each. Result: `3`. ✓

---

## Approach 2 — One-Pass Two Pointers (Optimal)

### Intuition
The nearest `word2` to any `word1` is the most recently seen one. Scan once; each time you see either word, record its latest index, and if the other word already has an index, the current gap is a candidate answer. Always pairing with the *latest* index of the other word guarantees no closer pair is missed.

### Algorithm
1. `i1 = i2 = -1` (last-seen indices of `word1`, `word2`).
2. For each `(k, w)` in `words`:
   - if `w == word1`: set `i1 = k`; if `i2 != -1`, update `min` with `i1 - i2`.
   - if `w == word2`: set `i2 = k`; if `i1 != -1`, update `min` with `i2 - i1`.
3. Return `min`.

### Complexity
- **Time:** O(n) — a single linear pass.
- **Space:** O(1) — two indices and a minimum.

### Code
```go
func twoPointers(words []string, word1 string, word2 string) int {
	i1, i2 := -1, -1
	min := len(words)
	for k, w := range words {
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
Trace `words = ["practice","makes","perfect","coding","makes"]`, `word1="coding"`, `word2="practice"`:

| k | w | i1 | i2 | candidate | min |
|---|---|----|----|-----------|-----|
| 0 | practice | -1 | 0 | i1=-1 → skip | 5 |
| 1 | makes | -1 | 0 | — | 5 |
| 2 | perfect | -1 | 0 | — | 5 |
| 3 | coding | 3 | 0 | i1-i2 = 3 | **3** |
| 4 | makes | 3 | 0 | — | 3 |

Result: `3`. ✓

---

## Key Takeaways
- **"Nearest occurrence" ⇒ track last-seen index** of each target while scanning once.
- Updating the gap only when *both* indices are known (`!= -1`) avoids spurious pairs before both words have appeared.
- Using `len(words)` as the initial minimum is a clean sentinel since no valid distance can reach it.
- This last-seen-index pattern generalizes: #244 (many queries → precompute index lists) and #245 (same word allowed → track two positions of the same word).

---

## Related Problems
- LeetCode #244 — Shortest Word Distance II (design; repeated queries)
- LeetCode #245 — Shortest Word Distance III (word1 may equal word2)
- LeetCode #821 — Shortest Distance to a Character (single-target variant)
