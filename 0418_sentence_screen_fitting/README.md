# 0418 — Sentence Screen Fitting

> LeetCode #418 · Difficulty: Medium · 🔒 Premium
> **Categories:** String, Dynamic Programming, Simulation

---

## Problem Statement

Given a `rows x cols` screen and a sentence represented by a list of non-empty words, find *how many times the given sentence can be fitted on the screen*.

**Note:**

- A word cannot be split into two lines.
- The order of words in the sentence must remain unchanged.
- Two consecutive words **in a line** must be separated by a single space.
- Total words in the sentence won't exceed `100`.
- Length of each word is greater than `0` and won't exceed `10`.
- `1 ≤ rows, cols ≤ 20,000`.

**Example 1:**

```
Input: rows = 2, cols = 8, sentence = ["hello", "world"]
Output: 1
Explanation:
hello---
world---

The character '-' signifies an empty space on the screen.
```

**Example 2:**

```
Input: rows = 3, cols = 6, sentence = ["a", "bcd", "e"]
Output: 2
Explanation:
a-bcd-
e-a---
bcd-e-

The character '-' signifies an empty space on the screen.
```

**Example 3:**

```
Input: rows = 4, cols = 5, sentence = ["I", "had", "apple", "pie"]
Output: 1
Explanation:
I-had
apple
pie-I
had--

The character '-' signifies an empty space on the screen.
```

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2023          |
| Amazon     | ★★★☆☆ Medium     | 2022          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Airbnb     | ★★☆☆☆ Low        | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Dynamic Programming (1D over starting word)** — a row's behaviour is a pure function of the word it starts on, giving only `n` distinct row "states"; precomputing `wordsPlaced[i]` and `nextStart[i]` and chaining them across rows is a classic 1D-DP / memoised-transition table → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)
- **String processing** — the whole task is packing space-separated words under a hard width, joining the sentence into a repeating character stream and reasoning about word boundaries → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (word-by-word) | O(rows·cols) | O(1) | Direct simulation; fine but touches up to 4·10⁸ cells at the limits |
| 2 | Next-Index stream simulation | O(rows·maxWordLen) | O(Σ word len) | Row-at-a-time via a repeated char stream; ~O(rows) |
| 3 | DP over starting word (Optimal) | O(n·cols) + O(rows) | O(n) | Best when `rows` is huge: precompute per-start-word jumps once |

---

## Approach 1 — Brute Force (Word-by-Word Simulation)

### Intuition

Do exactly what a typesetter does. Keep a pointer into the sentence and a count of free columns in the current row. Place the next word if it fits — remembering that any word that isn't first on its line must be preceded by a single space, so its true cost is `len(word)+1`. When a word doesn't fit, drop to the next row and reset the free columns. Every time the word pointer wraps past the last word, a whole sentence has landed on screen, so bump the counter. After `rows` rows, the counter is the answer.

### Algorithm

1. `wordIdx = 0`, `count = 0`.
2. For each row: `remaining = cols`. Repeatedly:
   - cost = `len(sentence[wordIdx])`, plus `1` if this isn't the first word on the line.
   - if cost `> remaining`, break to the next row.
   - else subtract cost, advance `wordIdx`; if it wraps to `0`, `count++`.
3. Return `count`.

### Complexity

- **Time:** O(rows·cols) — each placement consumes ≥ 1 column, so at most `cols` placements per row.
- **Space:** O(1) — a couple of counters.

### Code

```go
func bruteForce(rows, cols int, sentence []string) int {
	n := len(sentence)
	wordIdx := 0 // which word comes next
	count := 0   // completed sentences

	for r := 0; r < rows; r++ {
		remaining := cols // free columns left in this row
		for {
			wordLen := len(sentence[wordIdx])
			// A word that isn't first on the line needs a leading space, so its
			// "cost" is wordLen when remaining == cols, else wordLen+1.
			need := wordLen
			if remaining != cols {
				need = wordLen + 1 // account for the separating space
			}
			if need > remaining {
				break // this word can't start on the current row
			}
			remaining -= need // consume the word (and its space if any)
			wordIdx++         // move to the next word
			if wordIdx == n { // wrapped past the last word …
				wordIdx = 0 // … restart the sentence …
				count++     // … one full sentence placed
			}
		}
	}
	return count
}
```

### Dry Run

Example 2: `rows = 3, cols = 6, sentence = ["a","bcd","e"]` (n = 3).

| Row | remaining start | placements (word, cost, remaining after, wordIdx, count) |
|-----|-----------------|-----------------------------------------------------------|
| 0 | 6 | a(1)→5,idx1 · bcd(+1=4? cost4)→1,idx2 · e(+1=2)>1 stop | idx2, count0 |
| 1 | 6 | e(1)→5,idx0→**wrap count1** · a(+1=2)→3,idx1 · bcd(+1=4)>3 stop | idx1, count1 |
| 2 | 6 | bcd(3)→3,idx2 · e(+1=2)→1,idx0→**wrap count2** · a(+1=2)>1 stop | idx0, count2 |

Rows exhausted. Result: `2` ✔ — matching the layout `a-bcd- / e-a--- / bcd-e-`.

---

## Approach 2 — Precomputed "Next Index" Simulation

### Intuition

Instead of walking word by word, cut each row directly out of the infinitely repeated stream `"hello world hello world …"` (words joined by single spaces, with a trailing space so the last word also has a separator). Keep one global character index `start` into that stream. Advancing a row is just `start += cols`. Now inspect the character at the new right boundary: if it's a space, the row ended cleanly between words, so step over the space; if it's a letter, we've sliced a word in half, so retreat `start` back to the space before that word. At the end, `start / L` (where `L` is one sentence-block length) is exactly how many complete sentence copies were consumed.

### Algorithm

1. `s = join(sentence, " ") + " "`; `L = len(s)`.
2. `start = 0`. For each row:
   - `start += cols`.
   - if `s[start % L] == ' '`: `start++`.
   - else: while `start > 0 && s[(start-1) % L] != ' '`: `start--`.
3. Return `start / L`.

### Complexity

- **Time:** O(rows · maxWordLen) — each retreat backs up at most one word length (≤ 10), so effectively O(rows).
- **Space:** O(Σ word length) — the joined stream string.

### Code

```go
func nextIndexSim(rows, cols int, sentence []string) int {
	s := strings.Join(sentence, " ") + " " // one leading-normalised stream copy
	L := len(s)                            // length of a single sentence+space block
	start := 0                             // char index (into the infinite repeat) of the next row's first slot

	for r := 0; r < rows; r++ {
		start += cols // tentatively consume `cols` characters this row
		if s[start%L] == ' ' {
			start++ // landed exactly on a separating space — step over it
		} else {
			// Landed inside a word: retreat to the space before that word so we
			// don't cut it in half.
			for start > 0 && s[(start-1)%L] != ' ' {
				start--
			}
		}
	}
	return start / L // how many full sentence blocks were used up
}
```

### Dry Run

Example 2: `s = "a bcd e "`, `L = 8`, `cols = 6`.

| Row | start before | start += 6 | s[start%8] | Adjustment | start after |
|-----|--------------|------------|------------|------------|-------------|
| 0 | 0 | 6 | s[6]=`e` (letter) | retreat: s[5]=`d`,s[4]=`c`,s[3]=`b`,s[2]=`space` stop | 3 |
| 1 | 3 | 9 | s[9%8]=s[1]=`space` | step over: start++ | 10 |
| 2 | 10 | 16 | s[16%8]=s[0]=`a` (letter) | retreat: s[15%8]=s[7]=`space` stop | 16 |

`start / L = 16 / 8 = 2`. Result: `2` ✔

---

## Approach 3 — DP Over Starting Word (Optimal for many rows)

### Intuition

A row always begins on a word boundary, so its entire behaviour — how many words it fits and where the next row starts — is determined solely by **which word index it starts with**. There are only `n` such starting indices. Precompute two small tables once: `wordsPlaced[i]` (how many words a width-`cols` row packs when it starts at word `i`) and `nextStart[i]` (the word index the following row begins at). Then simulating all `rows` rows is just following `nextStart` and summing `wordsPlaced` — cheap even when `rows` is up to 20 000. Total words placed divided by `n` is the number of complete sentences.

### Algorithm

1. For each `i` in `0..n-1`: greedily pack words starting at `i` into a `cols`-wide row (each word costs `len+1` including its trailing space, and the row fits while `length + len(word) <= cols`). Record `wordsPlaced[i]` and the wrapped `nextStart[i]`.
2. `cur = 0`, `totalWords = 0`. Repeat `rows` times: `totalWords += wordsPlaced[cur]`; `cur = nextStart[cur]`.
3. Return `totalWords / n`.

### Complexity

- **Time:** O(n·cols) to build the tables + O(rows) to walk them.
- **Space:** O(n) — the two per-word tables.

### Code

```go
func dpStartWord(rows, cols int, sentence []string) int {
	n := len(sentence)
	wordsPlaced := make([]int, n) // words fitted when a row starts at index i
	nextStart := make([]int, n)   // starting index of the next row

	for i := 0; i < n; i++ {
		length := 0 // characters used so far in this hypothetical row
		words := 0  // words placed so far
		idx := i    // walking word pointer
		// Greedily add words while they fit. After the first word, each new word
		// costs len(word)+1 (its leading space).
		for length+len(sentence[idx]) <= cols {
			length += len(sentence[idx]) + 1 // +1 reserves the trailing space
			words++
			idx = (idx + 1) % n // wrap around the sentence
		}
		wordsPlaced[i] = words
		nextStart[i] = idx // where the next row will begin
	}

	cur := 0        // starting word of the current row
	totalWords := 0 // words placed across all rows
	for r := 0; r < rows; r++ {
		totalWords += wordsPlaced[cur] // add this row's contribution
		cur = nextStart[cur]           // jump to next row's starting word
	}
	return totalWords / n // each full sentence is n words
}
```

### Dry Run

Example 2: `sentence = ["a","bcd","e"]` (lengths 1,3,1), `cols = 6`, `n = 3`.

Precompute (each word's cost includes its trailing space; row fits while `length + len(word) ≤ 6`):

| start i | packing trace | wordsPlaced[i] | nextStart[i] |
|---------|---------------|----------------|--------------|
| 0 (`a`) | a→len2, bcd→len6, e? 6+1>6 stop | 2 | 2 |
| 1 (`bcd`) | bcd→len4, e→len6, a? 6+1>6 stop | 2 | 0 |
| 2 (`e`) | e→len2, a→len4, bcd? 4+3>6 stop | 2 | 1 |

Simulate 3 rows from `cur = 0`:

| Row | cur | wordsPlaced | totalWords | next cur |
|-----|-----|-------------|------------|----------|
| 0 | 0 | 2 | 2 | 2 |
| 1 | 2 | 2 | 4 | 1 |
| 2 | 1 | 2 | 6 | 0 |

`totalWords / n = 6 / 3 = 2`. Result: `2` ✔

---

## Key Takeaways

- **State = the word a row starts on.** Recognising that only `n` row-types exist collapses a potentially `rows`-long simulation into an O(n·cols) precompute plus an O(rows) chase — the key insight when `rows` dwarfs the sentence.
- **Turn wrapping text into a repeated character stream.** Joining words with spaces (plus a trailing space) lets a row be a modular slice `[start, start+cols)`; boundary handling becomes "is `s[boundary]` a space, else retreat" — no per-word bookkeeping.
- **Count by division, not by flag.** Both the stream (`start / L`) and DP (`totalWords / n`) approaches recover the sentence count from a running total, avoiding fragile wrap-detection.
- **Reserve the trailing space in the fit test.** Using cost `len(word)+1` uniformly and testing `length + len(word) ≤ cols` cleanly encodes "single space between consecutive words" without special-casing the first word.

---

## Related Problems

- LeetCode #68 — Text Justification (word packing under a width, then padding)
- LeetCode #1055 — Shortest Way to Form String (greedy consumption of a repeated string)
- LeetCode #686 — Repeated String Match (how many copies of a string are needed)
- LeetCode #1044 — Longest Duplicate Substring (reasoning over a repeated char stream)
