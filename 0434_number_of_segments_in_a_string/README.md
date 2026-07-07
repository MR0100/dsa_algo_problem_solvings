# 0434 — Number of Segments in a String

> LeetCode #434 · Difficulty: Easy
> **Categories:** String

---

## Problem Statement

Given a string `s`, return the number of segments in the string.

A **segment** is defined to be a contiguous sequence of **non-space characters**.

**Example 1:**

```
Input: s = "Hello, my name is John"
Output: 5
Explanation: The five segments are ["Hello,", "my", "name", "is", "John"]
```

**Example 2:**

```
Input: s = "Hello"
Output: 1
```

**Constraints:**

- `0 <= s.length <= 300`
- `s` consists of lowercase and uppercase English letters, digits, or one of the following characters `"!@#$%^&*()_+-=',.:"`.
- The only space character in `s` is `' '`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★☆☆☆ Low        | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Google     | ★☆☆☆☆ Rare       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **String Algorithms** — the entire problem is a single linear scan over the characters, counting token boundaries; recognising a "segment" as a non-space run preceded by a boundary is the core string-processing idea → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Built-in Split + Filter | O(n) | O(n) | Shortest to write; fine unless you must avoid allocation |
| 2 | Count Segment Starts (Optimal) | O(n) | O(1) | The clean interview answer — one pass, no allocation |
| 3 | Explicit State Machine | O(n) | O(1) | Same cost; generalises when "boundary" rules get complex |

`n` = length of `s`.

---

## Approach 1 — Built-in Split + Filter (Brute Force)

### Intuition

A segment is just a whitespace-delimited token. Go's `strings.Fields` splits a string on runs of whitespace *and* discards the empty pieces that leading/trailing/repeated spaces would otherwise create — which is precisely the segment definition. So the answer is simply the number of fields it returns. It is the least code, at the cost of allocating a slice of substrings only to count them.

### Algorithm

1. Call `strings.Fields(s)` → a slice of all non-empty whitespace-separated tokens.
2. Return `len(...)` of that slice.

### Complexity

- **Time:** O(n) — `Fields` scans the string once.
- **Space:** O(n) — it allocates the slice of token substrings.

### Code

```go
func builtinSplit(s string) int {
	return len(strings.Fields(s)) // Fields drops empties, so every field is a segment
}
```

### Dry Run

Example 1: `s = "Hello, my name is John"`.

| Step | Action | Result |
|------|--------|--------|
| 1 | `strings.Fields(s)` splits on spaces, drops empties | `["Hello,", "my", "name", "is", "John"]` |
| 2 | `len(...)` | `5` |

Output: `5` ✔

---

## Approach 2 — Count Segment Starts (One Pass, Optimal)

### Intuition

Instead of materialising the tokens, count each segment exactly once at the place it *begins*. A character starts a new segment iff it is **non-space** and its left neighbour is a **boundary** — either the start of the string or a space. Tallying those boundary positions counts the segments in a single pass with only a counter, and it handles empty strings, all-space strings, and repeated spaces automatically.

### Algorithm

1. Walk index `i` from `0` to `n-1`.
2. Increment `count` when `s[i] != ' '` **and** (`i == 0` **or** `s[i-1] == ' '`) — that marks the first character of a new segment.
3. Return `count`.

### Complexity

- **Time:** O(n) — one pass over the characters.
- **Space:** O(1) — just an integer counter.

### Code

```go
func countStarts(s string) int {
	count := 0
	for i := 0; i < len(s); i++ {
		// s[i] begins a segment iff it's a non-space preceded by a boundary
		// (start of string or a space).
		if s[i] != ' ' && (i == 0 || s[i-1] == ' ') {
			count++ // a new segment starts here
		}
	}
	return count
}
```

### Dry Run

Example 1: `s = "Hello, my name is John"` (spaces at indices 6, 9, 14, 17).

| i | s[i] | non-space? | boundary? (i==0 or s[i-1]==' ') | count after |
|---|------|-----------|----------------------------------|-------------|
| 0 | `H` | yes | yes (i==0) | 1 |
| 1–5 | `ello,` | yes | no (prev non-space) | 1 |
| 6 | `' '` | no | — | 1 |
| 7 | `m` | yes | yes (s[6]==' ') | 2 |
| 8 | `y` | yes | no | 2 |
| 9 | `' '` | no | — | 2 |
| 10 | `n` | yes | yes | 3 |
| … | `ame` | yes | no | 3 |
| 14 | `' '` | no | — | 3 |
| 15 | `i` | yes | yes | 4 |
| 16 | `s` | yes | no | 4 |
| 17 | `' '` | no | — | 4 |
| 18 | `J` | yes | yes | 5 |
| 19–21 | `ohn` | yes | no | 5 |

Output: `5` ✔

---

## Approach 3 — Explicit State Machine (In-Segment Flag)

### Intuition

The same boundary logic, framed as a two-state automaton: **OUTSIDE** (scanning spaces) and **INSIDE** (scanning a token). A rising edge OUTSIDE→INSIDE means a new segment just started, so count it; a space resets the machine to OUTSIDE. This form is handy when the boundary rule is richer than a one-character look-back (e.g. several whitespace kinds or escaped characters) — the flag captures "am I in a token?" without peeking backwards.

### Algorithm

1. Keep a boolean `inSegment`, initially `false`.
2. For each character:
   - if it's a space → set `inSegment = false`.
   - else if `inSegment` is `false` → it's a new segment: `count++` and set `inSegment = true`.
3. Return `count`.

### Complexity

- **Time:** O(n) — single pass.
- **Space:** O(1) — one flag and a counter.

### Code

```go
func stateMachine(s string) int {
	count := 0
	inSegment := false // are we currently scanning inside a token?
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' {
			inSegment = false // spaces end any current segment
		} else if !inSegment {
			count++          // rising edge: just entered a new segment
			inSegment = true // remember we're inside it now
		}
	}
	return count
}
```

### Dry Run

Example 1: `s = "Hello, my name is John"`.

| i | s[i] | branch | inSegment after | count after |
|---|------|--------|------------------|-------------|
| 0 | `H` | non-space, was OUTSIDE → count | true | 1 |
| 1–5 | `ello,` | non-space, already INSIDE | true | 1 |
| 6 | `' '` | space → OUTSIDE | false | 1 |
| 7 | `m` | non-space, was OUTSIDE → count | true | 2 |
| 9 | `' '` | space → OUTSIDE | false | 2 |
| 10 | `n` | non-space, was OUTSIDE → count | true | 3 |
| 14 | `' '` | space → OUTSIDE | false | 3 |
| 15 | `i` | non-space, was OUTSIDE → count | true | 4 |
| 17 | `' '` | space → OUTSIDE | false | 4 |
| 18 | `J` | non-space, was OUTSIDE → count | true | 5 |
| 19–21 | `ohn` | non-space, already INSIDE | true | 5 |

Output: `5` ✔

---

## Key Takeaways

- **Count boundaries, not tokens.** Anything of the form "how many maximal runs of X" is answered in O(1) space by counting *run starts* — a run of X begins where the current char is X and the previous char isn't.
- **`strings.Fields` ≠ `strings.Split(s, " ")`.** `Fields` splits on runs of whitespace and drops empties (matching "segment"); `Split` on a single space keeps empty tokens from adjacent/leading spaces and would over-count. Since this problem's only whitespace is `' '`, `Fields` is exact.
- The **state-machine framing** generalises the boundary rule cleanly — reach for it when "is this a separator?" is more than a single-character test.
- Watch the **edge cases**: empty string, all-spaces, and consecutive spaces should all yield the right count without special-casing — all three approaches do.

---

## Related Problems

- LeetCode #58 — Length of Last Word (scan for the final token from the right)
- LeetCode #557 — Reverse Words in a String III (per-segment processing)
- LeetCode #151 — Reverse Words in a String (tokenise and rejoin)
- LeetCode #1119 — Remove Vowels from a String (single linear string scan)
