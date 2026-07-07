# 0392 — Is Subsequence

> LeetCode #392 · Difficulty: Easy
> **Categories:** Two Pointers, String, Dynamic Programming

---

## Problem Statement

Given two strings `s` and `t`, return `true` if `s` is a subsequence of `t`, or `false`
otherwise.

A subsequence of a string is a new string that is formed from the original string by
deleting some (can be none) of the characters without disturbing the relative positions of
the remaining characters. (i.e., `"ace"` is a subsequence of `"abcde"` while `"aec"` is
not).

**Example 1:**

```
Input: s = "abc", t = "ahbgdc"
Output: true
```

**Example 2:**

```
Input: s = "axc", t = "ahbgdc"
Output: false
```

**Constraints:**

- `0 <= s.length <= 100`
- `0 <= t.length <= 10^4`
- `s` and `t` consist only of lowercase English letters.

**Follow up:** Suppose there are lots of incoming `s`, say `s1, s2, ..., sk` where
`k >= 10^9`, and you want to check one by one to see if `t` has its subsequence. In this
scenario, how would you change your code?

---

## Company Frequency

| Company   | Frequency         | Last Reported |
|-----------|-------------------|---------------|
| Google    | ★★★★☆ High        | 2024          |
| Amazon    | ★★★☆☆ Medium      | 2024          |
| Meta      | ★★★☆☆ Medium      | 2023          |
| Microsoft | ★★☆☆☆ Low         | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Two Pointers** — greedily match `s` against `t` in one pass → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **Dynamic Programming (2D)** — subsequence-matching table generalizes to edit distance → see [`/dsa/dynamic_programming_2d.md`](/dsa/dynamic_programming_2d.md)
- **String algorithms / preprocessing** — next-occurrence jump table for the many-query follow-up → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Two Pointers | O(n) | O(1) | Single query — simplest and optimal |
| 2 | DP Table | O(m·n) | O(m·n) | Teaching the subsequence lattice; edit-distance kin |
| 3 | Preprocessed Jump Table | O(n·26) prep, O(m)/query | O(n·26) | Follow-up: billions of `s` against one `t` |

---

## Approach 1 — Two Pointers

### Intuition

A subsequence keeps order but may skip characters. Walk `t` once; whenever the current
`t`-character equals the next needed `s`-character, consume it (advance the `s` pointer).
Matching the earliest possible `t`-character is always safe — deferring can never help.
If all of `s` is consumed, it is a subsequence.

### Algorithm

1. `i = 0` (index into `s`), scan `j` over `t`.
2. If `s[i] == t[j]`, advance `i`; always advance `j`.
3. `s` is a subsequence iff `i` reached `len(s)`.

### Complexity

- **Time:** O(n), n = `len(t)` — one linear pass.
- **Space:** O(1) — two indices.

### Code

```go
func twoPointers(s string, t string) bool {
	i := 0
	for j := 0; i < len(s) && j < len(t); j++ {
		if s[i] == t[j] {
			i++
		}
	}
	return i == len(s)
}
```

### Dry Run

`s = "abc"`, `t = "ahbgdc"`:

| j | t[j] | s[i] needed | match? | i after |
|---|------|-------------|--------|---------|
| 0 | a | a (i=0) | yes | 1 |
| 1 | h | b (i=1) | no | 1 |
| 2 | b | b (i=1) | yes | 2 |
| 3 | g | c (i=2) | no | 2 |
| 4 | d | c (i=2) | no | 2 |
| 5 | c | c (i=2) | yes | 3 |

`i == 3 == len(s)` ⇒ **`true`**.

---

## Approach 2 — DP Table

### Intuition

Let `dp[i][j]` = "is `s[i:]` a subsequence of `t[j:]`?" An empty `s` suffix matches
anything (last row all true). Then `s[i:]` matches `t[j:]` if we skip `t[j]`
(`dp[i][j+1]`), or, when `s[i]==t[j]`, match it and recurse (`dp[i+1][j+1]`).

### Algorithm

1. Build `dp` of size `(m+1) × (n+1)`; set `dp[m][*] = true`.
2. Fill from `i = m-1` down, `j = n-1` down:
   `dp[i][j] = dp[i][j+1] || (s[i]==t[j] && dp[i+1][j+1])`.
3. Answer is `dp[0][0]`.

### Complexity

- **Time:** O(m·n) — fill every cell once.
- **Space:** O(m·n) — the full table (compressible to O(n)).

### Code

```go
func dpBottomUp(s string, t string) bool {
	m, n := len(s), len(t)
	dp := make([][]bool, m+1)
	for i := range dp {
		dp[i] = make([]bool, n+1)
	}
	for j := 0; j <= n; j++ {
		dp[m][j] = true
	}
	for i := m - 1; i >= 0; i-- {
		for j := n - 1; j >= 0; j-- {
			skip := dp[i][j+1]
			take := s[i] == t[j] && dp[i+1][j+1]
			dp[i][j] = skip || take
		}
	}
	return dp[0][0]
}
```

### Dry Run

`s = "abc"` (m=3), `t = "ahbgdc"` (n=6). Base row `dp[3][*] = true`. Filling key cells
backwards (T = true, F = false):

| cell | meaning | value |
|------|---------|-------|
| `dp[2][5]` | "c" vs "c" | s[2]==t[5] && dp[3][6]=T → **T** |
| `dp[2][4]` | "c" vs "dc" | dp[2][5]=T → **T** |
| `dp[1][2]` | "bc" vs "bgdc" | s[1]==t[2] && dp[2][3] … resolves **T** |
| `dp[0][0]` | "abc" vs "ahbgdc" | s[0]==t[0] && dp[1][1] → **T** |

`dp[0][0] = true` ⇒ **`true`**.

---

## Approach 3 — Preprocessed Jump Table (Follow-up)

### Intuition

For billions of `s` against one fixed `t`, re-walking `t` each time is wasteful. Precompute
`nxt[j][c]` = the smallest index `>= j` in `t` where letter `c` occurs (or `len(t)` if
none). Each query then only jumps: at position `pos`, look up `nxt[pos][s[i]]`; if that is
`len(t)` the letter is missing, else move `pos` just past it.

### Algorithm

1. Build `nxt` with `len(t)+1` rows, 26 columns, bottom-up:
   `nxt[j] = nxt[j+1]` then override `nxt[j][t[j]] = j`; last row all `len(t)`.
2. Per query: `pos = 0`; for each char `c`, `j = nxt[pos][c]`; if `j == len(t)` return
   false, else `pos = j+1`.
3. Survive all chars ⇒ true.

### Complexity

- **Time:** O(n·26) preprocessing once, then O(m) per query.
- **Space:** O(n·26) for the jump table.

### Code

```go
func followUpManyQueries(s string, t string) bool {
	n := len(t)
	nxt := make([][26]int, n+1)
	for c := 0; c < 26; c++ {
		nxt[n][c] = n
	}
	for j := n - 1; j >= 0; j-- {
		nxt[j] = nxt[j+1]
		nxt[j][t[j]-'a'] = j
	}
	pos := 0
	for i := 0; i < len(s); i++ {
		c := s[i] - 'a'
		j := nxt[pos][c]
		if j == n {
			return false
		}
		pos = j + 1
	}
	return true
}
```

### Dry Run

`t = "ahbgdc"` (n=6). Relevant `nxt` values (next index of a letter at/after a position):
`nxt[0]['a']=0`, `nxt[1]['b']=2`, `nxt[3]['c']=5`.

Query `s = "abc"`:

| char | pos before | nxt[pos][char] | pos after |
|------|------------|----------------|-----------|
| a | 0 | 0 | 1 |
| b | 1 | 2 | 3 |
| c | 3 | 5 | 6 |

Never hit `n` before consuming all of `s` ⇒ **`true`**.

---

## Key Takeaways

- Greedy earliest-match two pointers is optimal for one subsequence check — O(n), O(1).
- The subsequence DP table is the skeleton of edit distance / LCS; know the transition.
- When one text answers *many* pattern queries, **preprocess a next-occurrence table** so
  each query costs only O(pattern length). This "jump table" pattern recurs widely.

---

## Related Problems

- LeetCode #524 — Longest Word in Dictionary through Deleting (many-query subsequence)
- LeetCode #792 — Number of Matching Subsequences (jump-table / bucketing follow-up)
- LeetCode #1143 — Longest Common Subsequence (2D DP generalization)
- LeetCode #115 — Distinct Subsequences (counting subsequences with DP)
