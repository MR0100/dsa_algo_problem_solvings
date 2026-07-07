# 0186 — Reverse Words in a String II

> LeetCode #186 · Difficulty: Medium (Premium)
> **Categories:** Two Pointers, String

---

## Problem Statement

Given a character array `s`, reverse the order of the **words**.

A **word** is defined as a sequence of non-space characters. The **words** in `s` will be separated by a single space.

Your code must solve the problem **in place**, i.e. without allocating extra space.

**Example 1:**
```
Input: s = ["t","h","e"," ","s","k","y"," ","i","s"," ","b","l","u","e"]
Output: ["b","l","u","e"," ","i","s"," ","s","k","y"," ","t","h","e"]
```

**Example 2:**
```
Input: s = ["a"]
Output: ["a"]
```

**Constraints:**
- `1 <= s.length <= 10^5`
- `s[i]` is an English letter (uppercase or lowercase), digit, or space `' '`.
- There is **at least one** word in `s`.
- `s` does not contain leading or trailing spaces.
- All the words in `s` are guaranteed to be separated by a single space.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Microsoft  | ★★★★☆ High       | 2024          |
| Amazon     | ★★★★☆ High       | 2024          |
| Meta       | ★★★☆☆ Medium     | 2024          |
| Apple      | ★★★☆☆ Medium     | 2023          |
| Uber       | ★★★☆☆ Medium     | 2023          |
| Google     | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Two Pointers** — converging read/write pointers perform every reversal, and a right-to-left scan extracts words without any split helper → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **String Algorithms (in-place reversal pattern)** — "reverse the whole thing, then reverse each piece" flips segment *order* while restoring segment *contents*; the exact same trick solves array rotation → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Split, Reverse, Copy Back) | O(n) | O(n) | Quick correctness baseline; ignores the in-place requirement |
| 2 | Two Pointers (Scan From the End) | O(n) | O(n) | Shows manual word parsing, still uses an output buffer |
| 3 | In-Place Double Reversal (Optimal) | O(n) | O(1) | The required answer: mutable buffer, zero auxiliary structures |

---

## Approach 1 — Brute Force (Split, Reverse, Copy Back)

### Intuition
Set the in-place constraint aside to get a correct baseline. Because the input guarantees *exactly one* space between words and no leading/trailing spaces, splitting on `" "` yields the clean word list, reversing that list is literally the task, and re-joining with single spaces produces a string of **exactly the same length** as `s` — so it can be copied straight back over the original array.

### Algorithm
1. Convert `s` to a `string` and split it on single spaces into a `words` slice.
2. Reverse `words` in place with two converging pointers (`i` from the front, `j` from the back; swap and step until they cross).
3. Join the reversed slice with single spaces and `copy` the bytes back into `s`.

### Complexity
- **Time:** O(n) — split, reverse, join, and copy each touch every byte a constant number of times.
- **Space:** O(n) — the words slice plus the joined string are full-size auxiliary copies (this is exactly what the problem asks us to avoid — see Approach 3).

### Code
```go
func bruteForce(s []byte) {
	// exactly one space between words, no leading/trailing spaces → Split is safe
	words := strings.Split(string(s), " ")
	// classic two-pointer reversal of the word order
	for i, j := 0, len(words)-1; i < j; i, j = i+1, j-1 {
		words[i], words[j] = words[j], words[i] // swap symmetric entries
	}
	// the joined string has the same length as s (same words, same separators)
	copy(s, strings.Join(words, " "))
}
```

### Dry Run
Example 1: `s = "the sky is blue"` (as a char array)

| Step | Action | State |
|------|--------|-------|
| 1 | `strings.Split(string(s), " ")` | `words = ["the","sky","is","blue"]` |
| 2 | swap `i=0, j=3` | `["blue","sky","is","the"]` |
| 3 | swap `i=1, j=2` | `["blue","is","sky","the"]` |
| 4 | `i=2, j=1` → pointers crossed, stop | `["blue","is","sky","the"]` |
| 5 | join + `copy` into `s` | `s = "blue is sky the"` ✓ |

---

## Approach 2 — Two Pointers (Scan From the End)

### Intuition
The answer is just the words read back-to-front. Walk the array from the last byte toward the first; every space we hit marks the *start* of the word we just walked across. Appending each word to an output buffer the moment it is discovered produces the words in exactly the reversed order — no library split, no separate reverse step.

### Algorithm
1. Keep `end` = index one past the current word (initialised to `len(s)`).
2. Scan `i` from `len(s)-1` down to `0`. When `s[i] == ' '`, the word is `s[i+1:end]`: append a separating space (if the buffer is non-empty) then the word, and set `end = i`.
3. After the loop the leftmost word is `s[0:end]`; append it last.
4. `copy` the buffer back over `s`.

### Complexity
- **Time:** O(n) — every byte is read once during the scan and written once into the buffer.
- **Space:** O(n) — the output buffer holds the entire rebuilt sentence, so the in-place requirement is still not met.

### Code
```go
func twoPointers(s []byte) {
	out := make([]byte, 0, len(s)) // rebuilt sentence, filled right-to-left by word
	end := len(s)                  // one past the end of the word currently being scanned
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == ' ' {
			if len(out) > 0 {
				out = append(out, ' ') // separator before every word except the first emitted
			}
			out = append(out, s[i+1:end]...) // the word we just walked over
			end = i                          // next word ends right before this space
		}
	}
	if len(out) > 0 {
		out = append(out, ' ') // separator before the final (leftmost) word
	}
	out = append(out, s[:end]...) // leftmost word has no space before it
	copy(s, out)                  // write the answer back into the caller's array
}
```

### Dry Run
Example 1: `s = "the sky is blue"` (indices 0–14; spaces at 3, 7, 10)

| Step | i | s[i] | Action | `out` | `end` |
|------|---|------|--------|-------|-------|
| 0 | — | — | init | `""` | 15 |
| 1 | 14→11 | letters | keep scanning | `""` | 15 |
| 2 | 10 | `' '` | append `s[11:15]` = `blue` | `"blue"` | 10 |
| 3 | 9→8 | letters | keep scanning | `"blue"` | 10 |
| 4 | 7 | `' '` | append `' '` + `s[8:10]` = `is` | `"blue is"` | 7 |
| 5 | 6→4 | letters | keep scanning | `"blue is"` | 7 |
| 6 | 3 | `' '` | append `' '` + `s[4:7]` = `sky` | `"blue is sky"` | 3 |
| 7 | 2→0 | letters | loop ends | `"blue is sky"` | 3 |
| 8 | — | — | append `' '` + `s[0:3]` = `the` | `"blue is sky the"` | 3 |
| 9 | — | — | `copy(s, out)` | `s = "blue is sky the"` ✓ | — |

---

## Approach 3 — In-Place Double Reversal (Optimal)

### Intuition
Reversing the **entire** array gets the words into the right order but leaves each word's letters backwards: `"the sky is blue"` → `"eulb si yks eht"`. Notice `eulb` sits exactly where `blue` must go — it is just internally reversed. So a second pass that reverses **each word individually** repairs the letters without moving the words: `"blue is sky the"`. Two reversals cancel inside each word but their *order* flip survives. All swaps happen inside `s`, so extra space is O(1).

### Algorithm
1. Reverse the whole array `s[0..n-1]`.
2. Scan left to right tracking `start`, the beginning of the current word.
3. At every boundary (a space, or one past the last byte), reverse the segment `s[start..i-1]` and set `start = i + 1`.

### Complexity
- **Time:** O(n) — each byte participates in at most two swaps (one per reversal pass) plus one boundary scan.
- **Space:** O(1) — only the `lo`/`hi`/`start` index variables; every write lands inside `s` itself.

### Code
```go
func inPlaceReversal(s []byte) {
	reverseRange(s, 0, len(s)-1) // step 1: whole-array reversal flips word order
	start := 0                   // start index of the word currently being scanned
	for i := 0; i <= len(s); i++ {
		// a boundary is either a space or one past the last byte
		if i == len(s) || s[i] == ' ' {
			reverseRange(s, start, i-1) // step 2: un-reverse this word's letters
			start = i + 1               // next word begins after the space
		}
	}
}

// reverseRange reverses s[lo..hi] in place with two converging pointers.
func reverseRange(s []byte, lo, hi int) {
	for lo < hi {
		s[lo], s[hi] = s[hi], s[lo] // swap the outermost pair
		lo++
		hi--
	}
}
```

### Dry Run
Example 1: `s = "the sky is blue"`

| Step | Action | s (as string) | start |
|------|--------|----------------|-------|
| 1 | reverse whole array `[0..14]` | `"eulb si yks eht"` | 0 |
| 2 | `i=4` is space → reverse `[0..3]` (`eulb`) | `"blue si yks eht"` | 5 |
| 3 | `i=7` is space → reverse `[5..6]` (`si`) | `"blue is yks eht"` | 8 |
| 4 | `i=11` is space → reverse `[8..10]` (`yks`) | `"blue is sky eht"` | 12 |
| 5 | `i=15 == len(s)` → reverse `[12..14]` (`eht`) | `"blue is sky the"` ✓ | 16 |

Example 2: `s = "a"` — step 1 reverses `[0..0]` (no-op); `i=1 == len(s)` reverses `[0..0]` (no-op) → `"a"` ✓

---

## Key Takeaways

- **Reverse-then-reverse** is the canonical O(1)-space trick for reordering segments of a mutable buffer: reverse the whole array to fix segment *order*, then reverse each segment to fix its *contents*. It reappears verbatim in LeetCode #189 (Rotate Array) and #151/#557 (word reversals).
- A right-to-left scan naturally yields words in reversed order — useful whenever you must emit tokens back-to-front in a single pass.
- When a problem guarantees clean formatting (single separators, no leading/trailing spaces), in-place algorithms get much simpler because the output occupies **exactly** the same bytes as the input.
- The premium variant differs from #151 precisely in mutability: `[]byte` (mutable) makes O(1) space achievable, while an immutable `string` forces O(n) output space.

---

## Related Problems

- LeetCode #151 — Reverse Words in a String (same task on an immutable string with messy spaces)
- LeetCode #557 — Reverse Words in a String III (reverse letters inside each word, keep word order)
- LeetCode #189 — Rotate Array (same double-reversal trick)
- LeetCode #344 — Reverse String (the two-pointer reversal primitive)
- LeetCode #61 — Rotate List (rotation/reordering on a linked list)
