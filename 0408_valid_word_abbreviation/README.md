# 0408 — Valid Word Abbreviation

> LeetCode #408 · Difficulty: Easy
> **Categories:** Two Pointers, String

---

## Problem Statement

A string can be **abbreviated** by replacing any number of **non-adjacent**, **non-empty** substrings with their lengths. The lengths **should not** have leading zeros.

For example, a string such as `"substitution"` could be abbreviated as (but not limited to):

- `"s10n"` (`"s ubstitutio n"`)
- `"sub4u4"` (`"sub stit u tion"`)
- `"12"` (`"substitution"`)
- `"su3i1u2on"` (`"su bst i t u ti on"`)
- `"substitution"` (no substrings replaced)

The following are **not valid** abbreviations:

- `"s55n"` (`"s ubsti tutio n"`, the replaced substrings are adjacent)
- `"s010n"` (has leading zeros)
- `"s0ubstitution"` (replaces an empty substring)

Given a string `word` and an abbreviation `abbr`, return *whether the string* `word` *matches the given abbreviation* `abbr`.

A substring is a contiguous **non-empty** sequence of characters within a string.

**Example 1:**

```
Input: word = "internationalization", abbr = "i12iz4n"
Output: true
Explanation: The word "internationalization" can be abbreviated as "i12iz4n" ("i nternational iz atio n").
```

**Example 2:**

```
Input: word = "apple", abbr = "a2e"
Output: false
Explanation: The word "apple" cannot be abbreviated as "a2e".
```

**Constraints:**

- `1 <= word.length <= 20`
- `word` consists of only lowercase English letters.
- `1 <= abbr.length <= 10`
- `abbr` consists of lowercase English letters and digits.
- All the integers in `abbr` will fit in a 32-bit integer.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Meta       | ★★★★★ Very High  | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Two Pointers** — advance one index over `word` and one over `abbr` in lockstep; letters step both by one, a number jumps the `word` pointer by its value → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **String Parsing** — the core subtlety is parsing multi-digit number tokens inline while enforcing the no-leading-zero rule, a classic scan-and-tokenise pattern → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Two Pointers (In-Place Parse) | O(n + m) | O(1) | The optimal, idiomatic answer; no allocation |
| 2 | Expand Then Compare (Reconstruct) | O(n + m) | O(1) extra | Same complexity; frames it as "rebuild and check", easier for some to reason about |

Here `n = len(word)`, `m = len(abbr)`.

---

## Approach 1 — Two Pointers (In-Place Parse)

### Intuition

Read `abbr` left to right while tracking where we are in `word`:

- A **letter** in `abbr` must match the current `word` character exactly — advance both pointers by one.
- A **digit** starts a number that means "skip this many characters of `word`" (they were compressed away). Parse the whole number, **reject a leading zero** (`"0…"` is never a valid length), then jump the `word` pointer forward by that amount.

The abbreviation is valid iff, once `abbr` is fully consumed, the `word` pointer sits **exactly** at the end — no leftover characters and no overshoot from a too-big number.

### Algorithm

1. `i = 0` over `word`, `j = 0` over `abbr`.
2. While `j < len(abbr)`:
   - If `abbr[j]` is a digit:
     - if it is `'0'`, return `false` (leading zero);
     - accumulate the consecutive digit run into `num`;
     - do `i += num`.
   - Else (a letter): if `i` is out of range or `word[i] != abbr[j]`, return `false`; otherwise advance both `i` and `j`.
3. Return `i == len(word)`.

### Complexity

- **Time:** O(n + m) — each character of `word` and `abbr` is visited once.
- **Space:** O(1) — two indices and an integer accumulator; nothing is allocated.

### Code

```go
func twoPointers(word string, abbr string) bool {
	i, j := 0, 0            // i indexes word, j indexes abbr
	n, m := len(word), len(abbr)
	for j < m {
		if abbr[j] >= '0' && abbr[j] <= '9' {
			// A number token: leading zero is illegal (e.g. "01", "0").
			if abbr[j] == '0' {
				return false
			}
			num := 0
			// Accumulate consecutive digits into the skip count.
			for j < m && abbr[j] >= '0' && abbr[j] <= '9' {
				num = num*10 + int(abbr[j]-'0')
				j++
			}
			i += num // jump over the compressed-away characters
		} else {
			// A literal letter: word must still have a char here and it must match.
			if i >= n || word[i] != abbr[j] {
				return false
			}
			i++ // consume the matched letter in word
			j++ // consume it in abbr
		}
	}
	// Both must land exactly at the end: leftover word chars ⇒ under-covered,
	// i overshooting ⇒ a number ran past the end.
	return i == n
}
```

### Dry Run

Example 1: `word = "internationalization"` (length 20), `abbr = "i12iz4n"`.

| Step | j | abbr[j] | Action | i after | j after |
|------|---|---------|--------|---------|---------|
| 0 | 0 | `i` | letter, word[0]=`i` match | 1 | 1 |
| 1 | 1 | `1` | digit run "12" → num=12; i += 12 | 13 | 3 |
| 2 | 3 | `i` | letter, word[13]=`i` match | 14 | 4 |
| 3 | 4 | `z` | letter, word[14]=`z` match | 15 | 5 |
| 4 | 5 | `4` | digit run "4" → num=4; i += 4 | 19 | 6 |
| 5 | 6 | `n` | letter, word[19]=`n` match | 20 | 7 |

`abbr` consumed (`j=7=m`). Final `i = 20 = len(word)` → **true** ✔

(`"internationalization"[13] = 'i'`, `[14] = 'z'`, `[19] = 'n'` — the abbreviation carves it into `i · nternational · iz · atio · n`.)

---

## Approach 2 — Expand Then Compare (Reconstruct)

### Intuition

Instead of matching lazily, rebuild the string that `abbr` *claims to represent* and check it against `word`. Copy letters verbatim; when a number `k` appears (rejecting leading zeros), splice in the next `k` characters of `word` using a cursor into `word`. If the cursor would ever run past the end of `word`, the abbreviation over-claims and is invalid. It is valid iff the cursor consumes exactly all of `word`. Because we compare positionally against `word` as we go (rather than materialising a new string), this stays O(1) extra space.

### Algorithm

1. `cursor = 0` (position in `word`), scan `abbr` with index `j`.
2. On a digit run: reject leading `'0'`; parse `k`; do `cursor += k`; if `cursor > len(word)`, return `false`.
3. On a letter: if `cursor >= len(word)` or `word[cursor] != letter`, return `false`; else `cursor++`.
4. Return `cursor == len(word)`.

### Complexity

- **Time:** O(n + m) — one pass over `abbr`, and the cursor advances through `word` at most `n` times total.
- **Space:** O(1) extra — we validate against `word` in place via the cursor (the `[]rune(abbr)` copy is bounded by `m ≤ 10`).

### Code

```go
func expandThenCompare(word string, abbr string) bool {
	cursor := 0 // how many characters of word we have accounted for
	n := len(word)
	j := 0
	runes := []rune(abbr) // treat abbr as runes for clean digit checks
	for j < len(runes) {
		ch := runes[j]
		if unicode.IsDigit(ch) {
			if ch == '0' {
				return false // leading zero not allowed in a length token
			}
			// Read the whole number token.
			start := j
			for j < len(runes) && unicode.IsDigit(runes[j]) {
				j++
			}
			k, _ := strconv.Atoi(string(runes[start:j])) // token value
			cursor += k                                  // these k chars are compressed
			if cursor > n {
				return false // number claims more characters than word has
			}
		} else {
			// Literal letter must line up with word at the cursor.
			if cursor >= n || rune(word[cursor]) != ch {
				return false
			}
			cursor++
			j++
		}
	}
	return cursor == n // every character of word must be accounted for exactly
}
```

### Dry Run

Example 1: `word = "internationalization"` (length 20), `abbr = "i12iz4n"`, `cursor = 0`.

| Step | j | token | Action | cursor after |
|------|---|-------|--------|--------------|
| 0 | 0 | `i` | letter, word[0]=`i` match | 1 |
| 1 | 1 | `12` | number 12; cursor += 12 (≤20 ok) | 13 |
| 2 | 3 | `i` | letter, word[13]=`i` match | 14 |
| 3 | 4 | `z` | letter, word[14]=`z` match | 15 |
| 4 | 5 | `4` | number 4; cursor += 4 (≤20 ok) | 19 |
| 5 | 6 | `n` | letter, word[19]=`n` match | 20 |

`abbr` fully scanned; `cursor = 20 = len(word)` → **true** ✔

---

## Key Takeaways

- **Two-pointer lockstep on two strings** is the go-to for "does encoding X describe string Y" problems: a plain token advances both cursors, a compressed token advances only the source cursor.
- **Parse numbers inline, don't split.** Reading a digit run into an accumulator (`num = num*10 + d`) avoids allocating substrings and handles arbitrary lengths within the 32-bit bound.
- **Guard the edge cases explicitly:** leading zero (`'0'` starting a token) is invalid; a number must not push the pointer past the end; and after consuming `abbr` the source pointer must land **exactly** at the end — an off-by-one here is the most common bug.
- The "expand and compare" reframing is a useful mental fallback and generalises to the harder follow-ups (#320 Generalized Abbreviation, #411 Minimum Unique Word Abbreviation) where you *generate* abbreviations.

---

## Related Problems

- LeetCode #320 — Generalized Abbreviation (enumerate all abbreviations)
- LeetCode #411 — Minimum Unique Word Abbreviation (search over abbreviations)
- LeetCode #527 — Word Abbreviation (shortest unambiguous abbreviations)
- LeetCode #443 — String Compression (in-place run-length encoding, similar parsing)
