# 0451 — Sort Characters By Frequency

> LeetCode #451 · Difficulty: Medium
> **Categories:** Hash Table, String, Bucket Sort, Sorting, Heap (Priority Queue), Counting

---

## Problem Statement

Given a string `s`, sort it in *decreasing order* based on the *frequency* of the characters. The frequency of a character is the number of times it appears in the string.

Return *the sorted string*. If there are multiple answers, return *any of them*.

**Example 1:**

```
Input: s = "tree"
Output: "eert"
Explanation: 'e' appears twice while 'r' and 't' both appear once.
So 'e' must appear before both 'r' and 't'. Therefore "eetr" is also a valid answer.
```

**Example 2:**

```
Input: s = "cccaaa"
Output: "aaaccc"
Explanation: Both 'c' and 'a' appear three times, so both "cccaaa" and "aaaccc" are valid answers.
Note that "cacaca" is incorrect, as the same characters must be together.
```

**Example 3:**

```
Input: s = "Aabb"
Output: "bbAa"
Explanation: "bbaA" is also a valid answer, but "Aabb" is incorrect.
Note that 'A' and 'a' are treated as two different characters.
```

**Constraints:**

- `1 <= s.length <= 5 * 10^5`
- `s` consists of uppercase and lowercase English letters and digits.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2024          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Apple      | ★★☆☆☆ Low        | 2022          |
| Adobe      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Map (frequency table)** — the first step of every approach is a character → count tally; a fixed-alphabet map (≤ 62 keys) makes counting O(n) → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Counting / Bucket Sort** — a character's frequency is an integer in `[1, n]`, a perfect bucket index; bucketing the characters by frequency avoids any comparison sort and gives a linear solution → see [`/dsa/counting_sort.md`](/dsa/counting_sort.md)
- **Sorting (comparison)** — the straightforward approach sorts the ≤ 62 distinct characters by descending count → see [`/dsa/sorting.md`](/dsa/sorting.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Count + Sort Unique Characters | O(n + k log k), k = distinct ≤ 62 | O(n) | Cleanest to write; the log factor is on a tiny alphabet so it is effectively O(n) |
| 2 | Bucket Sort by Frequency (Optimal) | O(n) | O(n) | Truly linear; the textbook "sort by count" trick, no comparison sort |

---

## Approach 1 — Count + Sort Unique Characters

### Intuition

The answer is simply each distinct character repeated as many times as it occurs, with the characters ordered by descending frequency. Two facts make this easy. First, "same characters must be together" means we never interleave — each character contributes one solid block. Second, the alphabet is tiny (upper + lower + digits = 62 symbols), so sorting the *distinct* characters is negligible work. Tally the counts, sort the distinct characters by count, and paste `char × count` blocks together.

### Algorithm

1. Build a frequency map `freq[b]` counting each byte `b` of `s`.
2. Collect the distinct bytes into a slice `chars`.
3. Sort `chars` by `freq` descending (ties may go in any order).
4. For each byte in sorted order, append it `freq[b]` times to a string builder; return the result.

### Complexity

- **Time:** O(n + k log k) — counting all `n` characters is O(n); sorting the `k ≤ 62` distinct characters is O(k log k); rebuilding writes `n` bytes. Since `k` is bounded by a constant, this is effectively O(n).
- **Space:** O(n) — the output string has length `n`; the map and key slice are O(k).

### Code

```go
func countThenSort(s string) string {
	freq := make(map[byte]int) // how many times each byte occurs
	for i := 0; i < len(s); i++ {
		freq[s[i]]++ // tally this character
	}

	// Collect the distinct characters so we can order them by frequency.
	chars := make([]byte, 0, len(freq))
	for b := range freq {
		chars = append(chars, b)
	}

	// Sort distinct characters by descending frequency (ties: any order is OK).
	sort.Slice(chars, func(i, j int) bool {
		return freq[chars[i]] > freq[chars[j]] // higher count comes first
	})

	var sb strings.Builder    // efficient string assembly, no repeated realloc
	sb.Grow(len(s))           // we know the final length up front
	for _, b := range chars { // most frequent first
		// Write this character exactly freq[b] times (block them together).
		sb.WriteString(strings.Repeat(string(b), freq[b]))
	}
	return sb.String()
}
```

### Dry Run

Example 1: `s = "tree"`.

| Step | Action | State |
|------|--------|-------|
| 1 | Count characters | `freq = {t:1, r:1, e:2}` |
| 2 | Collect distinct chars | `chars = [t, r, e]` (map order arbitrary) |
| 3 | Sort by freq desc | `chars = [e, t, r]` (e has count 2; t,r have count 1) |
| 4 | Emit `e`×2 | builder = `"ee"` |
| 4 | Emit `t`×1 | builder = `"eet"` |
| 4 | Emit `r`×1 | builder = `"eetr"` |

Result: `"eetr"` ✔ (a valid answer — 'e' block comes before the single 't' and 'r').

---

## Approach 2 — Bucket Sort by Frequency (Optimal)

### Intuition

Sorting by an integer key that lives in a small known range never needs comparisons — it needs buckets. Here the key is a character's frequency, an integer between `1` and `n`. Create one bucket per possible frequency, drop each character into `bucket[freq]`, then read the buckets from the highest frequency down to the lowest and emit each character `freq` times. This is counting sort applied to the frequencies, so the whole thing is linear regardless of alphabet size.

### Algorithm

1. Count characters into `freq[b]`.
2. Allocate `buckets` of length `n+1`; for each character `b`, append `b` to `buckets[freq[b]]`.
3. Iterate `f` from `n` down to `1`; for every character `b` in `buckets[f]`, append `b` exactly `f` times.
4. Return the assembled string.

### Complexity

- **Time:** O(n) — counting is O(n); there are `n+1` buckets; across the whole emit loop we write exactly `n` bytes. No log factor.
- **Space:** O(n) — the buckets array plus the output string are both bounded by `n`.

### Code

```go
func bucketSort(s string) string {
	n := len(s)
	freq := make(map[byte]int) // character → count
	for i := 0; i < n; i++ {
		freq[s[i]]++
	}

	// buckets[f] = list of characters that appear exactly f times.
	// A count can be at most n, so we need indices 0..n.
	buckets := make([][]byte, n+1)
	for b, f := range freq {
		buckets[f] = append(buckets[f], b) // drop char into its frequency bucket
	}

	var sb strings.Builder
	sb.Grow(n)
	// Walk from the highest possible frequency down to 1 (0 is always empty
	// of real characters). This yields descending-frequency output.
	for f := n; f >= 1; f-- {
		for _, b := range buckets[f] { // every char with this exact frequency
			// Emit the character f times so equal chars stay grouped.
			sb.WriteString(strings.Repeat(string(b), f))
		}
	}
	return sb.String()
}
```

### Dry Run

Example 1: `s = "tree"` (n = 4).

| Step | Action | State |
|------|--------|-------|
| 1 | Count | `freq = {t:1, r:1, e:2}` |
| 2 | Fill buckets (size 5) | `buckets[1] = [t, r]`, `buckets[2] = [e]`, others empty |
| 3 | `f = 4` | `buckets[4]` empty → nothing |
| 3 | `f = 3` | `buckets[3]` empty → nothing |
| 3 | `f = 2` | emit `e`×2 → `"ee"` |
| 3 | `f = 1` | emit `t`×1 → `"eet"`, then `r`×1 → `"eetr"` |

Result: `"eetr"` ✔ (bucket at frequency 2 is drained before frequency 1).

---

## Key Takeaways

- **"Sort by count" ⇒ bucket sort.** When the sort key is a frequency (an integer bounded by `n`), you never need a comparison sort — index directly into `buckets[frequency]` for O(n).
- **Small fixed alphabet ⇒ the sort is free.** With ≤ 62 possible characters, even `sort.Slice` on the distinct characters is effectively constant work, so the simple approach is competitive in practice.
- **"Same characters must be together"** is the tell that the answer is a concatenation of solid blocks, one per character — not an interleaving. Grouping is mandatory; the *order of equal-frequency blocks* is free.
- A **priority-max-heap of (count, char)** is a third route (pop the most frequent repeatedly) but adds a log factor for no benefit here; bucket sort dominates it.

---

## Related Problems

- LeetCode #347 — Top K Frequent Elements (bucket-by-frequency to pick the top k)
- LeetCode #692 — Top K Frequent Words (frequency + tie-break ordering)
- LeetCode #1636 — Sort Array by Increasing Frequency (same key, ascending)
- LeetCode #75 — Sort Colors (counting sort on a tiny value range)
- LeetCode #767 — Reorganize String (frequency-driven placement, but must *not* group)
