# 0245 — Shortest Word Distance III

> LeetCode #245 · Difficulty: Medium · 🔒 Premium
> **Categories:** Array, String, Two Pointers

---

## Problem Statement

Given an array of strings `wordsDict` and two strings that already exist in the array `word1` and `word2`, return *the shortest distance between the occurrences of these two words in the list*.

**Note** that `word1` and `word2` **may be the same**. It is guaranteed that they represent **two individual words** in the list.

**Example 1:**

```
Input: wordsDict = ["practice", "makes", "perfect", "coding", "makes"], word1 = "makes", word2 = "coding"
Output: 1
```

**Example 2:**

```
Input: wordsDict = ["practice", "makes", "perfect", "coding", "makes"], word1 = "makes", word2 = "makes"
Output: 3
```

**Constraints:**

- `1 <= wordsDict.length <= 10^5`
- `1 <= wordsDict[i].length <= 10`
- `wordsDict[i]` consists of lowercase English letters.
- `word1` and `word2` are in `wordsDict`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| LinkedIn   | ★★★★☆ High       | 2023          |
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Meta       | ★★☆☆☆ Low        | 2023          |
| Google     | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Two Pointers (last-seen indices)** — track the most recent index of each word (or the previous occurrence when they coincide) in a single sweep → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **Array Traversal** — a linear scan with case-split logic → see [`/dsa/arrays.md`](/dsa/arrays.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (all valid pairs) | O(n²) | O(1) | Small inputs; easy to verify the `i != j` rule |
| 2 | One-Pass Two Pointers (Optimal) | O(n) | O(1) | Always preferred — handles both cases in one scan |

(n = length of `wordsDict`.)

---

## Approach 1 — Brute Force

### Intuition
The answer is `|i - j|` over pairs where `words[i] == word1` and `words[j] == word2`. The only new twist versus #243 is that when `word1 == word2` we must require `i != j`, otherwise we'd wrongly report distance 0 by pairing a position with itself.

### Algorithm
1. Initialize `min` to `len(words)`.
2. For each `i` with `words[i] == word1`:
   - For each `j` with `words[j] == word2` **and** `i != j`, update `min` with `|i - j|`.
3. Return `min`.

### Complexity
- **Time:** O(n²) — nested scan over matching positions.
- **Space:** O(1) — a running minimum.

### Code
```go
func bruteForce(words []string, word1 string, word2 string) int {
	min := len(words)
	for i := 0; i < len(words); i++ {
		if words[i] != word1 {
			continue
		}
		for j := 0; j < len(words); j++ {
			if i == j {
				continue
			}
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
Trace `words = ["practice","makes","perfect","coding","makes"]`, `word1="makes"`, `word2="coding"`:

| i (word1="makes") | j (word2="coding", j≠i) | \|i-j\| | min |
|-------------------|-------------------------|---------|-----|
| 1 | 3 | 2 | 2 |
| 4 | 3 | 1 | **1** |

Result: `1`. ✓

---

## Approach 2 — One-Pass Two Pointers (Optimal)

### Intuition
Split on whether the two words are equal:
- **Distinct words** → identical to #243: keep the last-seen index of each and measure the gap when either appears.
- **Same word** → every occurrence plays both roles. The closest pair of *different* positions of one word is always two **consecutive** occurrences, so keep a single `prev` index; each new occurrence pairs with `prev`, and we track the minimum consecutive gap.

### Algorithm
1. If `word1 == word2`: scan; on each occurrence, if `prev != -1` update `min` with `i - prev`, then set `prev = i`.
2. Else: `i1 = i2 = -1`; on `word1` set `i1` and try `i1-i2`; on `word2` set `i2` and try `i2-i1`.
3. Return `min`.

### Complexity
- **Time:** O(n) — a single linear pass.
- **Space:** O(1) — a couple of indices and a minimum.

### Code
```go
func twoPointers(words []string, word1 string, word2 string) int {
	min := len(words)

	if word1 == word2 {
		prev := -1
		for k, w := range words {
			if w == word1 {
				if prev != -1 && k-prev < min {
					min = k - prev
				}
				prev = k
			}
		}
		return min
	}

	i1, i2 := -1, -1
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
Trace the **same-word** case `word1 = word2 = "makes"` on `["practice","makes","perfect","coding","makes"]`:

| k | w | is "makes"? | prev before | candidate `k-prev` | min | prev after |
|---|---|-------------|-------------|--------------------|-----|-----------|
| 0 | practice | no | -1 | — | 5 | -1 |
| 1 | makes | yes | -1 | prev=-1 → skip | 5 | 1 |
| 2 | perfect | no | 1 | — | 5 | 1 |
| 3 | coding | no | 1 | — | 5 | 1 |
| 4 | makes | yes | 1 | 4-1 = 3 | **3** | 4 |

Result: `3`. ✓

For the **distinct-word** query `word1="makes"`, `word2="coding"` the scan behaves exactly like #243 and yields `1`.

---

## Key Takeaways
- **Case-split on `word1 == word2`.** The equal case is the interesting one: the answer is the minimum gap between *consecutive* occurrences, so a single `prev` pointer suffices.
- The `i != j` (or `prev != -1` before pairing) guard is what stops a word from being reported as distance 0 from itself.
- The distinct-word branch reuses the exact #243 last-seen-index technique — this problem generalizes it.
- One linear pass and O(1) space handle both cases; no need to store index lists here (that's #244's multi-query concern).

---

## Related Problems
- LeetCode #243 — Shortest Word Distance (distinct words only)
- LeetCode #244 — Shortest Word Distance II (design; repeated queries)
- LeetCode #821 — Shortest Distance to a Character (single-target variant)
