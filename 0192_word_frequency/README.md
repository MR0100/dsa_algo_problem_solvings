# 0192 — Word Frequency

> LeetCode #192 · Difficulty: Medium
> **Categories:** Shell, Hash Table, Sorting, String

---

## Problem Statement

Write a bash script to calculate the frequency of each word in a text file `words.txt`.

For simplicity sake, you may assume:

- `words.txt` contains only lowercase characters and space `' '` characters.
- Each word must consist of lowercase characters only.
- Words are separated by one or more whitespace characters.

**Example:**

Assume that `words.txt` has the following content:

```
the day is sunny the the
the sunny is is
```

Your script should output the following, sorted by descending frequency:

```
the 4
is 3
sunny 2
day 1
```

**Note:**

- Don't worry about handling ties, it is guaranteed that each word's frequency count is unique.
- Could you write it in one-line using [Unix pipes](http://tldp.org/LDP/abs/html/x17601.html)?

> **Repo note:** #192 is one of LeetCode's four Shell problems. Following this
> repo's Go-only convention, every approach below re-implements the identical
> word-count-and-rank pipeline in Go, treating the file's content as an
> in-memory string (`words.txt` → the `wordsTxt` constant). The canonical
> one-line bash answer appears in Key Takeaways for interview completeness.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |
| Google     | ★★☆☆☆ Low        | 2023          |
| LinkedIn   | ★☆☆☆☆ Rare       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Map** — counting occurrences is the textbook hash-map job: `freq[word]++` turns each update into O(1) instead of a linear scan → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Sorting / Bucket Sort** — ranking words by count is a sort on a bounded integer key, so both comparison sort (O(U log U)) and bucket sort (O(W)) apply → see [`/dsa/sorting.md`](/dsa/sorting.md)
- **String Algorithms** — whitespace tokenisation ("one or more whitespace characters") is the parsing half of the problem → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (parallel slices + selection sort) | O(W·U + U²) | O(U) | To show the baseline; no hash map, no library sort |
| 2 | Hash Map + Comparison Sort | O(W + U log U) | O(U) | The standard interview answer; simple and near-optimal |
| 3 | Bucket Sort by Frequency (Optimal) | O(W + U) | O(W + U) | When U is large and you want to beat the comparison-sort bound |

*W = total words in the file, U = distinct words (U ≤ W).*

---

## Approach 1 — Brute Force (Parallel Slices + Selection Sort)

### Intuition
The most primitive counting device is a list you scan front-to-back: "have I filed this word already? If yes bump its tally, else open a new entry." Sorting by count can likewise use the most primitive sort — repeatedly select the remaining maximum and swap it into place. Zero O(1)-lookup machinery, so every step's cost is visible.

### Algorithm
1. Split the file into words on any whitespace run (`strings.Fields`).
2. For each word, linearly scan the `uniques` slice; on a match bump the paired `counts[i]`, otherwise append a new `(word, 1)` entry.
3. Selection-sort both slices in lock-step: for each slot `i`, find the entry with the highest count among `i..end` (alphabetically smallest on a tie) and swap it into slot `i` in **both** slices.
4. Format each pair as `"word count"`.

### Complexity
- **Time:** O(W·U + U²) — each of the W words scans up to U uniques (counting phase), then selection sort costs U² comparisons.
- **Space:** O(U) — the two parallel slices of distinct words and counts.

### Code
```go
func bruteForce(fileContent string) []string {
	words := strings.Fields(fileContent) // split on any run of spaces/newlines/tabs

	uniques := []string{} // distinct words in first-seen order
	counts := []int{}     // counts[i] = occurrences of uniques[i]

	for _, w := range words {
		found := false
		for i := range uniques { // linear scan — the "no hash map" price
			if uniques[i] == w {
				counts[i]++ // seen before → bump its tally
				found = true
				break
			}
		}
		if !found { // first sighting → open a new entry with count 1
			uniques = append(uniques, w)
			counts = append(counts, 1)
		}
	}

	// Selection sort by count descending (alphabetical ascending on ties).
	for i := 0; i < len(uniques); i++ {
		best := i // index of the best remaining entry
		for j := i + 1; j < len(uniques); j++ {
			higher := counts[j] > counts[best]                                // strictly more frequent wins
			tieAlpha := counts[j] == counts[best] && uniques[j] < uniques[best] // tie → smaller word wins
			if higher || tieAlpha {
				best = j
			}
		}
		// Swap the winner into slot i in BOTH slices to keep them aligned.
		counts[i], counts[best] = counts[best], counts[i]
		uniques[i], uniques[best] = uniques[best], uniques[i]
	}

	out := make([]string, len(uniques))
	for i := range uniques {
		out[i] = fmt.Sprintf("%s %d", uniques[i], counts[i]) // "word count" line
	}
	return out
}
```

### Dry Run
Example: `words = [the day is sunny the the the sunny is is]` (10 words).

Counting phase (state of the parallel slices after each word):

| word processed | uniques | counts |
|----------------|---------|--------|
| the | `[the]` | `[1]` |
| day | `[the day]` | `[1 1]` |
| is | `[the day is]` | `[1 1 1]` |
| sunny | `[the day is sunny]` | `[1 1 1 1]` |
| the | `[the day is sunny]` | `[2 1 1 1]` |
| the | `[the day is sunny]` | `[3 1 1 1]` |
| the | `[the day is sunny]` | `[4 1 1 1]` |
| sunny | `[the day is sunny]` | `[4 1 1 2]` |
| is | `[the day is sunny]` | `[4 1 2 2]` |
| is | `[the day is sunny]` | `[4 1 3 2]` |

Selection-sort phase:

| i | best found (count) | swap | uniques after | counts after |
|---|--------------------|------|---------------|--------------|
| 0 | `the` (4) at 0 | none | `[the day is sunny]` | `[4 1 3 2]` |
| 1 | `is` (3) at 2 | day↔is | `[the is day sunny]` | `[4 3 1 2]` |
| 2 | `sunny` (2) at 3 | day↔sunny | `[the is sunny day]` | `[4 3 2 1]` |
| 3 | `day` (1) at 3 | none | `[the is sunny day]` | `[4 3 2 1]` |

Output: `the 4`, `is 3`, `sunny 2`, `day 1`. ✓

---

## Approach 2 — Hash Map + Comparison Sort

### Intuition
Counting is the textbook hash-map job — each word update is O(1) instead of a linear scan. The only remaining work is ordering U distinct words, which a comparison sort handles in O(U log U). This is the direct analogue of the Unix pipeline `sort | uniq -c | sort -nr`, with the first (expensive) sort replaced by a map.

### Algorithm
1. Split the file into words; `freq[word]++` for each.
2. Collect the map's keys into a slice (Go map iteration order is random, so an explicit sort is mandatory).
3. `sort.Slice` by frequency descending, word ascending on ties (defensive — the problem guarantees no ties).
4. Format `"word count"` lines in sorted order.

### Complexity
- **Time:** O(W + U log U) — one linear counting pass over W words plus the comparison sort of U distinct words.
- **Space:** O(U) — the frequency map and the key slice.

### Code
```go
func hashMapSort(fileContent string) []string {
	freq := make(map[string]int) // word → number of occurrences
	for _, w := range strings.Fields(fileContent) {
		freq[w]++ // O(1) amortised update per word
	}

	words := make([]string, 0, len(freq))
	for w := range freq { // gather distinct words (map order is random)
		words = append(words, w)
	}

	sort.Slice(words, func(i, j int) bool {
		if freq[words[i]] != freq[words[j]] {
			return freq[words[i]] > freq[words[j]] // primary: higher count first
		}
		return words[i] < words[j] // secondary: alphabetical (defensive; no ties guaranteed)
	})

	out := make([]string, len(words))
	for i, w := range words {
		out[i] = fmt.Sprintf("%s %d", w, freq[w])
	}
	return out
}
```

### Dry Run
Example 1, counting pass (map state after each word):

| word | freq map state |
|------|----------------|
| the | `{the:1}` |
| day | `{the:1 day:1}` |
| is | `{the:1 day:1 is:1}` |
| sunny | `{the:1 day:1 is:1 sunny:1}` |
| the, the, the | `{the:4 day:1 is:1 sunny:1}` |
| sunny | `{the:4 day:1 is:1 sunny:2}` |
| is, is | `{the:4 day:1 is:3 sunny:2}` |

Sort pass: keys `[day is sunny the]` (some random order) → compare by count: `the(4) > is(3) > sunny(2) > day(1)` → sorted `[the is sunny day]`.

Output: `the 4`, `is 3`, `sunny 2`, `day 1`. ✓

---

## Approach 3 — Bucket Sort by Frequency (Optimal)

### Intuition
Comparison sorting is overkill when the sort key is a small bounded integer: a word appearing `c` times always satisfies `1 ≤ c ≤ W`. So an array of W+1 buckets — bucket `c` holds "all words occurring exactly `c` times" — replaces the O(U log U) sort with two linear sweeps. This is the same trick as LeetCode 347 *Top K Frequent Elements*, and it makes the whole pipeline linear in the input size.

### Algorithm
1. Count words into a hash map (as in Approach 2).
2. Create `buckets[0..W]`; append each word `w` to `buckets[freq[w]]`.
3. Walk `c` from W down to 1; skip empty buckets; emit `"w c"` for every word in `buckets[c]` (alphabetically sorted for determinism — each bucket holds at most one word here thanks to the unique-frequency guarantee).

### Complexity
- **Time:** O(W + U) — one counting pass plus one sweep over W+1 buckets containing U words total; per-bucket sorts are no-ops given the guarantee.
- **Space:** O(W + U) — the bucket array (W+1 slots) dominates the map.

### Code
```go
func bucketSort(fileContent string) []string {
	words := strings.Fields(fileContent)

	freq := make(map[string]int) // word → occurrences
	for _, w := range words {
		freq[w]++
	}

	// buckets[c] holds every distinct word occurring exactly c times.
	// A frequency can never exceed len(words), so W+1 buckets always suffice.
	buckets := make([][]string, len(words)+1)
	for w, c := range freq {
		buckets[c] = append(buckets[c], w)
	}

	out := []string{}
	for c := len(words); c >= 1; c-- { // highest frequency first
		if len(buckets[c]) == 0 {
			continue // no word has this exact count
		}
		sort.Strings(buckets[c]) // deterministic tie order (problem guarantees ≤1 word here)
		for _, w := range buckets[c] {
			out = append(out, fmt.Sprintf("%s %d", w, c))
		}
	}
	return out
}
```

### Dry Run
Example 1: W = 10 words, map = `{the:4, day:1, is:3, sunny:2}`.

Bucket fill:

| word | count c | buckets state (non-empty only) |
|------|---------|--------------------------------|
| the | 4 | `b[4]=[the]` |
| day | 1 | `b[4]=[the] b[1]=[day]` |
| is | 3 | `b[4]=[the] b[3]=[is] b[1]=[day]` |
| sunny | 2 | `b[4]=[the] b[3]=[is] b[2]=[sunny] b[1]=[day]` |

Emit sweep (`c` from 10 down to 1):

| c | bucket content | emitted |
|---|----------------|---------|
| 10..5 | empty | — |
| 4 | `[the]` | `the 4` |
| 3 | `[is]` | `is 3` |
| 2 | `[sunny]` | `sunny 2` |
| 1 | `[day]` | `day 1` |

Output: `the 4`, `is 3`, `sunny 2`, `day 1`. ✓

---

## Key Takeaways

- **Count with a hash map, rank with a sort** — the universal two-phase shape of every frequency problem (words, characters, numbers).
- **Bucket sort beats comparison sort when the key is a bounded small integer** — frequency ≤ total elements, so `O(n)` buckets always suffice (same trick as LC 347 / LC 451).
- **Go maps iterate in random order** — any "sorted output" requirement forces an explicit key sort; never rely on map order.
- **Parallel slices must be swapped in lock-step** — when hand-rolling sorts over two aligned arrays, forgetting one swap silently corrupts pairings.
- The intended one-line bash answer, worth knowing verbatim:
  ```bash
  cat words.txt | tr -s ' ' '\n' | sort | uniq -c | sort -nr | awk '{print $2" "$1}'
  ```
  `tr -s` squeezes runs of spaces into single newlines (one word per line), `sort | uniq -c` counts duplicates, `sort -nr` ranks by count descending, `awk` swaps "count word" into "word count".

---

## Related Problems

- LeetCode #347 — Top K Frequent Elements (same hash-count + bucket-sort pattern)
- LeetCode #692 — Top K Frequent Words (adds explicit alphabetical tie-breaking + heap)
- LeetCode #451 — Sort Characters By Frequency (frequency bucket sort on characters)
- LeetCode #193 — Valid Phone Numbers (Shell problem set)
- LeetCode #194 — Transpose File (Shell problem set)
- LeetCode #195 — Tenth Line (Shell problem set)
