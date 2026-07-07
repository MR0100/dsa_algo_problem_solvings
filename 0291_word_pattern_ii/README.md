# 0291 — Word Pattern II

> LeetCode #291 · Difficulty: Medium
> **Categories:** Backtracking, Hash Table, String

---

## Problem Statement

Given a `pattern` and a string `s`, return `true` if `s` **matches** the `pattern`.

A string `s` **matches** a `pattern` if there is some **bijective mapping** of single characters to **non-empty** strings such that if each character in `pattern` is replaced by the string it maps to, then the resulting string is `s`. A **bijective mapping** means that no two characters map to the same string, and no character maps to two different strings.

**Example 1:**
```
Input: pattern = "abab", s = "redblueredblue"
Output: true
Explanation: One possible mapping is as follows:
'a' -> "red"
'b' -> "blue"
```

**Example 2:**
```
Input: pattern = "aaaa", s = "asdasdasdasd"
Output: true
Explanation: One possible mapping is as follows:
'a' -> "asd"
```

**Example 3:**
```
Input: pattern = "aabb", s = "xyzabcxzyabc"
Output: false
```

**Constraints:**
- `1 <= pattern.length <= 20`
- `1 <= s.length <= 100`
- `pattern` and `s` consist of only lowercase English letters.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |
| Uber       | ★★☆☆☆ Low        | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used
- **Backtracking** — try each split of `s`, recurse, and undo on failure → see [`/dsa/backtracking.md`](/dsa/backtracking.md)
- **Hash Map (bijection)** — enforce char→word and word→char consistency in O(1) → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **String Algorithms** — substring slicing and prefix matching → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview
| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Backtracking (One Map) | O(nᵐ)·O(m) check | O(m) | Simplest correct search |
| 2 | Backtracking (Two Maps) (Optimal) | O(nᵐ) | O(m) | O(1) bijection check |

(*m* = pattern length, *n* = length of `s`.)

---

## Approach 1 — Backtracking (One Map)

### Intuition
We need a bijection between pattern characters and non-empty substrings so that concatenating each character's substring in order rebuilds `s`. Substring lengths are unknown, so we search: for the current pattern character, if it's already mapped, the next slice of `s` must equal that word; otherwise we try every possible prefix of the remaining string as its word. A single `char→word` map enforces "same letter → same word"; to also enforce "different letters → different words", we linearly scan the map's values.

### Algorithm
1. `dfs(pi, si)` walks pattern index `pi` and string index `si` together.
2. If `pi == len(pattern)`: success iff `si == len(s)` (both fully consumed).
3. Let `c = pattern[pi]`. If `c` is already mapped to `w`, the slice `s[si:si+len(w)]` must equal `w`; recurse to `(pi+1, si+len(w))`.
4. Otherwise, for every `end` from `si+1..len(s)`, take `w = s[si:end]`; skip it if another character already uses `w`; tentatively map `c→w`, recurse; unmap on failure.

### Complexity
- **Time:** O(nᵐ) — each of the `m` pattern characters tries up to `n` split lengths, with an O(m) value scan per new assignment.
- **Space:** O(m) — the map holds at most one entry per distinct pattern character; recursion depth is O(m).

### Code
```go
func backtrackOneMap(pattern string, s string) bool {
	charToWord := map[byte]string{} // pattern letter → assigned substring
	var dfs func(pi, si int) bool
	dfs = func(pi, si int) bool {
		if pi == len(pattern) { // consumed the whole pattern...
			return si == len(s) // ...succeed only if we also consumed all of s
		}
		c := pattern[pi]
		if w, ok := charToWord[c]; ok {
			end := si + len(w)
			if end > len(s) || s[si:end] != w {
				return false // mismatch → this branch is dead
			}
			return dfs(pi+1, end) // consume w and move on
		}
		for end := si + 1; end <= len(s); end++ {
			w := s[si:end]
			used := false
			for _, existing := range charToWord {
				if existing == w {
					used = true
					break
				}
			}
			if used {
				continue // would break the bijection
			}
			charToWord[c] = w // tentatively assign
			if dfs(pi+1, end) {
				return true
			}
			delete(charToWord, c) // backtrack: undo the assignment
		}
		return false
	}
	return dfs(0, 0)
}
```

### Dry Run
`pattern = "abab"`, `s = "redblueredblue"`

| Step | pi | si | c | Action | Map state |
|------|----|----|---|--------|-----------|
| 1 | 0 | 0 | a | unbound; try `w="r"` | `{a:r}` |
| 2 | 1 | 1 | b | unbound; try `w="e"` | `{a:r,b:e}` |
| 3 | 2 | 2 | a | bound to "r"; `s[2:3]="d"≠"r"` → fail, backtrack | `{a:r}` |
| … | | | | eventually `a="red"` tried | `{a:red}` |
| 4 | 1 | 3 | b | unbound; try `w="b"` … up to `w="blue"` | `{a:red,b:blue}` |
| 5 | 2 | 7 | a | bound "red"; `s[7:10]="red"` ✓ | `{a:red,b:blue}` |
| 6 | 3 | 10 | b | bound "blue"; `s[10:14]="blue"` ✓ | — |
| 7 | 4 | 14 | — | `pi==len(pattern)` and `si==len(s)` → **true** | — |

---

## Approach 2 — Backtracking (Two Maps) (Optimal)

### Intuition
A bijection needs BOTH directions enforced: `char→word` and `word→char`. Keeping an inverse `word→char` map lets us reject in O(1) both "this letter already maps elsewhere" and "this word is already claimed by another letter", instead of scanning all values as Approach 1 does.

### Algorithm
1. `dfs(pi, si)` as before.
2. If `pi == len(pattern)`: return `si == len(s)`.
3. If `c = pattern[pi]` is bound to `w`, verify `s[si:si+len(w)] == w`; recurse.
4. Otherwise, for each candidate `w = s[si:end]`: if `w` is already in `wordToChar`, skip; else set both maps, recurse, and unset both on backtrack.

### Complexity
- **Time:** O(nᵐ) — same split search, but each bijection check is O(1).
- **Space:** O(m) — two maps with at most `m` entries plus O(m) recursion depth.

### Code
```go
func backtrackTwoMaps(pattern string, s string) bool {
	charToWord := map[byte]string{} // forward map: letter → word
	wordToChar := map[string]byte{} // inverse map: word → letter
	var dfs func(pi, si int) bool
	dfs = func(pi, si int) bool {
		if pi == len(pattern) {
			return si == len(s) // both fully consumed ⇒ valid matching
		}
		c := pattern[pi]
		if w, ok := charToWord[c]; ok {
			end := si + len(w)
			if end > len(s) || s[si:end] != w {
				return false
			}
			return dfs(pi+1, end)
		}
		for end := si + 1; end <= len(s); end++ {
			w := s[si:end]
			if _, taken := wordToChar[w]; taken {
				continue // another letter already owns this word → injectivity broken
			}
			charToWord[c] = w // bind both directions
			wordToChar[w] = c
			if dfs(pi+1, end) {
				return true
			}
			delete(charToWord, c) // backtrack both maps together
			delete(wordToChar, w)
		}
		return false
	}
	return dfs(0, 0)
}
```

### Dry Run
`pattern = "abab"`, `s = "redblueredblue"`

| Step | pi | si | c | Action | charToWord / wordToChar |
|------|----|----|---|--------|-------------------------|
| 1 | 0 | 0 | a | try `w="red"` (after shorter tries fail) | `{a:red}` / `{red:a}` |
| 2 | 1 | 3 | b | try `w="blue"` | `{a:red,b:blue}` / `{red:a,blue:b}` |
| 3 | 2 | 7 | a | bound; `s[7:10]="red"` ✓ (O(1) map hit) | unchanged |
| 4 | 3 | 10 | b | bound; `s[10:14]="blue"` ✓ | unchanged |
| 5 | 4 | 14 | — | `pi==len` and `si==len` → **true** | — |

---

## Key Takeaways
- Unknown-length splits + a consistency constraint = **backtracking with a map**. Try prefixes, recurse, undo.
- A **bijection** requires enforcing *both* directions; a second inverse map turns the injectivity check from O(m) to O(1).
- Always pair the success base case with "consumed all of `s`" — reaching the end of the pattern is not enough.
- Backtracking cleanliness: whatever you set before recursing, you must `delete` after it fails.

---

## Related Problems
- LeetCode #290 — Word Pattern (same idea, but words are pre-split by spaces)
- LeetCode #139 — Word Break (split a string against a dictionary)
- LeetCode #140 — Word Break II (enumerate all splits — backtracking)
- LeetCode #472 — Concatenated Words
