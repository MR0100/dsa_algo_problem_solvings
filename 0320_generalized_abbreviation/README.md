# 0320 — Generalized Abbreviation

> LeetCode #320 · Difficulty: Medium
> **Categories:** String, Backtracking, Bit Manipulation, Recursion

---

## Problem Statement

A word's **generalized abbreviation** can be constructed by taking any number of
**non-overlapping** and **non-adjacent** substrings and replacing them with their
respective lengths.

- For example, `"abcde"` can be abbreviated into:
  - `"a3e"` (`"bcd"` turned into `"3"`)
  - `"1bcd1"` (`"a"` and `"e"` both turned into `"1"`)
  - `"5"` (`"abcde"` turned into `"5"`)
  - `"abcde"` (no substrings replaced)
- However, these abbreviations are **invalid**:
  - `"23"` (`"ab"` turned into `"2"` and `"cde"` turned into `"3"`) has adjacent
    substrings chosen.
  - `"22de"` (`"ab"` turned into `"2"` and `"bc"` turned into `"2"`) has
    overlapping substrings chosen.

Given a string `word`, return *a list of all the possible generalized
abbreviations of* `word`. Return the answer in **any order**.

**Example 1:**

```
Input:  word = "word"
Output: ["4","3d","2r1","2rd","1o2","1o1d","1or1","1ord","w3","w2d","w1r1","w1rd","wo2","wo1d","wor1","word"]
```

**Example 2:**

```
Input:  word = "a"
Output: ["1","a"]
```

**Constraints:**

- `1 <= word.length <= 15`
- `word` consists of only lowercase English letters.

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Google    | ★★★★☆ High       | 2023          |
| Facebook  | ★★★☆☆ Medium     | 2022          |
| Amazon    | ★★★☆☆ Medium     | 2023          |
| Microsoft | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Backtracking** — at each character, branch into keep vs. abbreviate and
  recurse → see [`/dsa/backtracking.md`](/dsa/backtracking.md)
- **Bit Manipulation (subset enumeration)** — a 2ⁿ bitmask assigns each char to
  keep/abbreviate; enumerate all masks →
  see [`/dsa/bit_manipulation.md`](/dsa/bit_manipulation.md)
- **String Algorithms** — flushing consecutive abbreviated runs into count
  numbers → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Backtracking (keep/abbreviate) | O(2ⁿ·n) | O(n) | Interview default; clean recursion |
| 2 | Bitmask enumeration | O(2ⁿ·n) | O(n) | Iterative, no recursion; easy to reason about |

Both are optimal: there are exactly 2ⁿ abbreviations, so any solution is Ω(2ⁿ).

---

## Approach 1 — Backtracking (Keep or Abbreviate Each Char)

### Intuition
Walk the word carrying `count`, the number of consecutive characters currently
being abbreviated. At index `i` branch two ways: **abbreviate** `word[i]`
(`count+1`, emit nothing yet), or **keep** `word[i]` (first flush any pending
`count` as a number, then append the literal, and reset `count`). Flushing on a
kept character is what guarantees abbreviations are **non-adjacent** — a number
segment is always separated from the next by at least one literal.

### Algorithm
1. `dfs(i, count, cur)`:
   - If `i == len(word)`: append `cur + (count>0 ? count : "")`; return.
   - Branch A (abbreviate): `dfs(i+1, count+1, cur)`.
   - Branch B (keep): `next = cur + (count>0 ? count : "") + word[i]`;
     `dfs(i+1, 0, next)`.
2. Start with `dfs(0, 0, "")`.

### Complexity
- **Time:** O(2ⁿ·n) — 2ⁿ leaves, each assembling a string of length ≤ n.
- **Space:** O(n) recursion depth (output list not counted).

### Code
```go
func backtracking(word string) []string {
	res := []string{}
	var dfs func(i, count int, cur string)
	dfs = func(i, count int, cur string) {
		if i == len(word) {
			if count > 0 {
				cur += strconv.Itoa(count)
			}
			res = append(res, cur)
			return
		}
		dfs(i+1, count+1, cur)
		next := cur
		if count > 0 {
			next += strconv.Itoa(count)
		}
		next += string(word[i])
		dfs(i+1, 0, next)
	}
	dfs(0, 0, "")
	return res
}
```

### Dry Run
Input `word = "a"` (n = 1). Start `dfs(0, 0, "")`.

| Call | i | count | cur | action |
|------|---|-------|-----|--------|
| dfs(0,0,"") | 0 | 0 | "" | Branch A → dfs(1,1,"") |
| dfs(1,1,"") | 1 | 1 | "" | end: count>0 → append "1" |
| back to dfs(0) | 0 | 0 | "" | Branch B: keep 'a' → dfs(1,0,"a") |
| dfs(1,0,"a") | 1 | 0 | "a" | end: count=0 → append "a" |

Result: `["1", "a"]`. ✓

---

## Approach 2 — Bitmask Enumeration

### Intuition
Each of the `n` characters is independently kept or abbreviated, so the `2ⁿ`
binary masks enumerate every abbreviation exactly once. Bit `i` set means
"abbreviate `word[i]`". For each mask, sweep left to right accumulating a run of
abbreviated chars into `count`, and flush that count as a number whenever a kept
character appears.

### Algorithm
1. For `mask = 0 .. 2ⁿ - 1`:
   - `count = 0`, `cur = ""`.
   - For `i = 0..n-1`: if bit `i` set → `count++`; else flush `count` (if > 0),
     append `word[i]`, reset `count`.
   - Flush any trailing `count`. Append `cur` to results.

### Complexity
- **Time:** O(2ⁿ·n) — `2ⁿ` masks, each an O(n) sweep.
- **Space:** O(n) per built string.

### Code
```go
func bitmask(word string) []string {
	n := len(word)
	res := make([]string, 0, 1<<n)
	for mask := 0; mask < (1 << n); mask++ {
		count := 0
		cur := ""
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				count++
			} else {
				if count > 0 {
					cur += strconv.Itoa(count)
					count = 0
				}
				cur += string(word[i])
			}
		}
		if count > 0 {
			cur += strconv.Itoa(count)
		}
		res = append(res, cur)
	}
	return res
}
```

### Dry Run
Input `word = "a"` (n = 1). Masks 0 and 1.

| mask | i=0 bit | count/cur trace | result string |
|------|---------|-----------------|---------------|
| 0    | clear   | keep 'a' → cur="a" | `"a"`      |
| 1    | set     | count=1; end flush → cur="1" | `"1"` |

Result: `["a", "1"]`. ✓

---

## Key Takeaways

- **Two-way per-element branch (keep / transform)** is a textbook backtracking
  shape; carrying a running `count` avoids emitting partial numbers early.
- **`2ⁿ` independent binary choices ⇒ bitmask enumeration.** Any "each element is
  in or out" generation problem can be written iteratively over `0..2ⁿ-1`.
- **Flush-on-literal enforces the non-adjacency rule for free** — you never need
  an explicit "were the last two both numbers?" check.
- The answer size is exactly `2ⁿ`, so both approaches are asymptotically optimal;
  choice is a matter of recursion vs. iteration preference.

---

## Related Problems

- LeetCode #78 — Subsets (2ⁿ enumeration, keep/drop per element)
- LeetCode #93 — Restore IP Addresses (backtracking string partition)
- LeetCode #17 — Letter Combinations of a Phone Number (backtracking build)
- LeetCode #408 — Valid Word Abbreviation (the verification counterpart)
