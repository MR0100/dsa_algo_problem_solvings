# 0151 — Reverse Words in a String

> LeetCode #151 · Difficulty: Medium
> **Categories:** String, Two Pointers, Stack

---

## Problem Statement

Given an input string `s`, reverse the order of the **words**.

A **word** is defined as a sequence of non-space characters. The **words** in `s` will be separated by at least one space.

Return a string of the words in reverse order concatenated by a single space.

**Note** that `s` may contain leading or trailing spaces or multiple spaces between two words. The returned string should only have a single space separating the words. Do not include any extra spaces.

**Example 1:**
```
Input: s = "the sky is blue"
Output: "blue is sky the"
```

**Example 2:**
```
Input: s = "  hello world  "
Output: "world hello"
Explanation: Your reversed string should not contain leading or trailing spaces.
```

**Example 3:**
```
Input: s = "a good   example"
Output: "example good a"
Explanation: You need to reduce multiple spaces between two words to a single space in the reversed string.
```

**Constraints:**
- `1 <= s.length <= 10^4`
- `s` contains English letters (upper-case and lower-case), digits, and spaces `' '`.
- There is **at least one** word in `s`.

**Follow-up:** If the string data type is mutable in your language, can you solve it **in-place** with `O(1)` extra space?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Microsoft  | ★★★★★ Very High  | 2024          |
| Amazon     | ★★★★☆ High       | 2024          |
| Apple      | ★★★★☆ High       | 2024          |
| Meta       | ★★★☆☆ Medium     | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Two Pointers** — a right-to-left scan extracts words in reversed order in one pass, and a read/write pointer pair compacts spaces in place → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **In-place reversal pattern** — "reverse the whole array, then reverse each segment" flips segment *order* while preserving segment *contents*; the same trick powers array rotation → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Split, Reverse, Join) | O(n) | O(n) | Production code / fastest to write |
| 2 | Two Pointers (Scan From the End) | O(n) | O(n) output only | Interview: shows you can parse without library helpers |
| 3 | In-Place Reversal (Optimal) | O(n) | O(1) extra | The follow-up: mutable buffer, no auxiliary structures |

---

## Approach 1 — Brute Force (Split, Reverse, Join)

### Intuition
The task is literally "reverse the order of the words". If we can obtain a clean list of words — with all the messy extra spaces already stripped — the answer is just that list reversed and re-joined with single spaces. `strings.Fields` does exactly the cleaning for us: it splits on *runs* of whitespace and never produces empty strings.

### Algorithm
1. Split `s` into a slice of words using `strings.Fields(s)`.
2. Reverse the slice in place with two converging pointers (`i` from the front, `j` from the back, swap and step until they cross).
3. Join the reversed slice with a single space and return it.

### Complexity
- **Time:** O(n) — splitting scans every character once, reversing swaps n/2 slice entries, joining copies every character once.
- **Space:** O(n) — the word slice plus the output string.

### Code
```go
func bruteForce(s string) string {
	// Fields splits around any run of spaces and never yields empty strings,
	// so "  hello   world  " becomes ["hello","world"].
	words := strings.Fields(s)
	// classic two-pointer in-place reversal of the slice
	for i, j := 0, len(words)-1; i < j; i, j = i+1, j-1 {
		words[i], words[j] = words[j], words[i] // swap symmetric positions
	}
	// join with exactly one space — output format requires single separators
	return strings.Join(words, " ")
}
```

### Dry Run
Example 1: `s = "the sky is blue"`

| Step | Action | State |
|------|--------|-------|
| 1 | `strings.Fields` | `words = ["the","sky","is","blue"]` |
| 2 | swap `i=0,j=3` | `["blue","sky","is","the"]` |
| 3 | swap `i=1,j=2` | `["blue","is","sky","the"]` |
| 4 | `i=2,j=1` → pointers crossed, stop | `["blue","is","sky","the"]` |
| 5 | `strings.Join(words, " ")` | `"blue is sky the"` ✓ |

---

## Approach 2 — Two Pointers (Scan From the End)

### Intuition
The last word of `s` must appear first in the answer. So walk the string backwards: skip any spaces, mark where a word ends, keep walking left until the word starts, and emit that slice. Because we discover words back-to-front, they land in the output already in reversed order — no separate reversal pass, no intermediate word array.

### Algorithm
1. Initialize a `strings.Builder` and set `i = len(s) - 1`.
2. While `i >= 0`:
   1. Decrement `i` while `s[i]` is a space (skips trailing spaces and separators).
   2. If `i < 0`, stop — only spaces remained.
   3. Set `end = i` (last character of the current word).
   4. Decrement `i` while `i >= 0` and `s[i]` is not a space (find the word start).
   5. If the builder is non-empty, write one `' '` separator.
   6. Write `s[i+1 : end+1]` (the whole word) to the builder.
3. Return the built string.

### Complexity
- **Time:** O(n) — the pointer `i` moves strictly left; every character is inspected exactly once.
- **Space:** O(n) for the output builder only — no intermediate slice of words (O(1) auxiliary beyond the mandatory output).

### Code
```go
func twoPointers(s string) string {
	var sb strings.Builder // accumulates the answer without repeated copying
	sb.Grow(len(s))        // pre-allocate: the answer is never longer than s
	i := len(s) - 1        // right pointer, starts at the end of the string
	for i >= 0 {
		// skip the run of spaces between words (and trailing spaces)
		for i >= 0 && s[i] == ' ' {
			i--
		}
		if i < 0 { // nothing but spaces left → all words emitted
			break
		}
		end := i // inclusive index of the word's last character
		// walk left until we fall off the word (space or string start)
		for i >= 0 && s[i] != ' ' {
			i--
		}
		if sb.Len() > 0 { // separator only between words, never leading
			sb.WriteByte(' ')
		}
		// s[i+1:end+1] is the whole word: i stopped one left of its start
		sb.WriteString(s[i+1 : end+1])
	}
	return sb.String()
}
```

### Dry Run
Example 1: `s = "the sky is blue"` (indices 0–14)

| Iter | skip spaces → `i` | `end` | walk word → `i` | word emitted | builder |
|------|-------------------|-------|-----------------|--------------|---------|
| 1 | 14 (no spaces) | 14 | 10 (`s[10]=' '`) | `s[11:15]="blue"` | `"blue"` |
| 2 | 9 | 9 | 7 (`s[7]=' '`) | `s[8:10]="is"` | `"blue is"` |
| 3 | 6 | 6 | 3 (`s[3]=' '`) | `s[4:7]="sky"` | `"blue is sky"` |
| 4 | 2 | 2 | -1 (fell off front) | `s[0:3]="the"` | `"blue is sky the"` |
| 5 | `i=-1` → loop exits | — | — | — | return `"blue is sky the"` ✓ |

---

## Approach 3 — In-Place Reversal (Optimal)

### Intuition
Reversing the **entire** string reverses the word *order* (what we want) but also reverses the *letters* inside every word (what we don't). Reversing each individual word afterwards puts the letters back while keeping the new word order:

```
"the sky is blue"  --reverse all-->  "eulb si yks eht"  --reverse each word-->  "blue is sky the"
```

The extra spaces are handled first by an in-place compaction using a read pointer and a write pointer, so the buffer contains exactly single-space-separated words before any reversing happens. In a language with mutable strings this is a true O(1)-extra-space solution — the follow-up's answer.

### Algorithm
1. Copy `s` into a byte slice `b` (Go strings are immutable; in mutable-string languages this step disappears).
2. **Compact spaces in place:** read pointer `r` skips each run of spaces; write pointer `w` copies each word's bytes down, writing a single `' '` before every word except the first. Truncate `b` to length `w`.
3. **Reverse the whole buffer** `b[0..w-1]` — word order is now correct, letters are backwards.
4. **Reverse each word:** scan for word boundaries (space or end of buffer) and reverse each `b[start..end]` segment.
5. Return `string(b)`.

### Complexity
- **Time:** O(n) — compaction is one pass, the full reversal is one pass, and per-word reversal touches each character once more (three linear passes total).
- **Space:** O(1) extra beyond the single working buffer — no word list, no builder; the buffer itself is only required because Go strings are immutable.

### Code
```go
func inPlaceReversal(s string) string {
	b := []byte(s) // single mutable working buffer

	// -- step 1: compact spaces in place ------------------------------------
	w := 0 // write pointer: next position to fill in the cleaned buffer
	r := 0 // read pointer: current position in the original content
	for r < len(b) {
		// skip a run of spaces (leading spaces or separators)
		for r < len(b) && b[r] == ' ' {
			r++
		}
		if r == len(b) { // trailing spaces only → done
			break
		}
		if w > 0 { // one space before every word except the first
			b[w] = ' '
			w++
		}
		// copy the word's bytes down to the write position
		for r < len(b) && b[r] != ' ' {
			b[w] = b[r]
			w++
			r++
		}
	}
	b = b[:w] // truncate: buffer now holds "word word word" exactly

	// -- step 2: reverse the entire buffer (reverses word order + letters) --
	reverseRange(b, 0, len(b)-1)

	// -- step 3: reverse each word back so letters read correctly -----------
	start := 0 // start index of the current word
	for i := 0; i <= len(b); i++ {
		// a word ends at a space or at the end of the buffer
		if i == len(b) || b[i] == ' ' {
			reverseRange(b, start, i-1) // fix this word's letters
			start = i + 1               // next word starts after the space
		}
	}
	return string(b)
}

// reverseRange reverses b[lo..hi] in place with two converging pointers.
func reverseRange(b []byte, lo, hi int) {
	for lo < hi {
		b[lo], b[hi] = b[hi], b[lo] // swap the outer pair
		lo++
		hi--
	}
}
```

### Dry Run
Example 1: `s = "the sky is blue"` (already single-spaced, so compaction copies verbatim; see the table's first row for how `w`/`r` finish)

| Step | Operation | Buffer state |
|------|-----------|--------------|
| 1 | compact spaces (`r` ends at 15, `w` ends at 15) | `"the sky is blue"` (unchanged, len 15) |
| 2 | reverse whole buffer `[0..14]` | `"eulb si yks eht"` |
| 3 | word 1: boundary at `i=4` (space) → reverse `[0..3]` | `"blue si yks eht"`, `start=5` |
| 4 | word 2: boundary at `i=7` (space) → reverse `[5..6]` | `"blue is yks eht"`, `start=8` |
| 5 | word 3: boundary at `i=11` (space) → reverse `[8..10]` | `"blue is sky eht"`, `start=12` |
| 6 | word 4: boundary at `i=15` (end) → reverse `[12..14]` | `"blue is sky the"` |
| 7 | return | `"blue is sky the"` ✓ |

For Example 2 (`"  hello world  "`), step 1 compacts the buffer to `"hello world"` (w=11) before the reversals, which is why no leading/trailing spaces survive.

---

## Key Takeaways

- **Reverse-all-then-reverse-each** flips the *order* of segments while preserving their *contents* — the same trick solves array rotation (LeetCode #189) in O(1) space.
- **Read/write pointer compaction** is the standard in-place way to delete unwanted characters/elements: `r` scans everything, `w` only advances over kept bytes (same skeleton as Remove Element #27 and Remove Duplicates #26).
- Scanning **from the end** produces reversed order for free — whenever the output order is the reverse of the scan order, consider starting the scan from the other side instead of reversing afterwards.
- In Go, `strings.Fields` + `strings.Join` is the pragmatic answer; the byte-slice version is what demonstrates the O(1)-space follow-up in interviews.
- Emit separators *before* each element except the first (`if sb.Len() > 0`) — cleaner than trimming a trailing separator afterwards.

---

## Related Problems

- LeetCode #186 — Reverse Words in a String II (true in-place, char array input)
- LeetCode #557 — Reverse Words in a String III (reverse letters inside each word, keep word order)
- LeetCode #189 — Rotate Array (same reverse-all/reverse-parts trick)
- LeetCode #344 — Reverse String (the two-pointer reversal primitive)
- LeetCode #58 — Length of Last Word (same backwards word scan)
