# 0072 — Edit Distance

> LeetCode #72 · Difficulty: Hard
> **Categories:** String, Dynamic Programming

---

## Problem Statement

Given two strings `word1` and `word2`, return the **minimum number of operations** required to convert `word1` to `word2`.

You have the following three operations permitted on a word:
- Insert a character
- Delete a character
- Replace a character

**Example 1**
```
Input:  word1 = "horse", word2 = "ros"
Output: 3
Explanation:
horse → rorse (replace 'h' with 'r')
rorse → rose  (delete 'r')
rose  → ros   (delete 'e')
```

**Example 2**
```
Input:  word1 = "intention", word2 = "execution"
Output: 5
```

**Constraints**
- `0 <= word1.length, word2.length <= 500`
- `word1` and `word2` consist of lowercase English letters.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★★ Very High | 2024          |
| Meta      | ★★★★★ Very High | 2024          |
| Microsoft | ★★★★★ Very High | 2024          |
| Bloomberg | ★★★★☆ High      | 2024          |
| Apple     | ★★★★☆ High      | 2023          |
| Adobe     | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Dynamic Programming (2D)** — `dp[i][j]` = min operations to convert `word1[0..i-1]` to `word2[0..j-1]`.
- **Three Operations** — insert, delete, replace each have a direct interpretation in the DP recurrence.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Memoization (Top-Down) | O(m × n) | O(m × n) | Explicit recursion with caching |
| 2 | DP 2D Table ✅ | O(m × n) | O(m × n) | Standard; clearest to explain |
| 3 | DP Rolling Row | O(m × n) | O(n) | Space-optimised; same algorithm |

---

## Approach 1 — Memoization (Top-Down DP)

### Intuition
`dp(i, j)` = min ops to convert `word1[i:]` to `word2[j:]`.
- Base: `dp(i, n) = m - i` (delete remaining chars of word1); `dp(m, j) = n - j` (insert remaining chars of word2).
- Match: `word1[i] == word2[j]` → `dp(i,j) = dp(i+1, j+1)`.
- No match: `dp(i,j) = 1 + min(dp(i+1,j), dp(i,j+1), dp(i+1,j+1))`.

### Complexity
- **Time:** O(m × n).
- **Space:** O(m × n) — memo + O(m+n) stack.

---

## Approach 2 — DP 2D Table (Recommended ✅)

### Intuition
`dp[i][j]` = min ops to convert `word1[0..i-1]` to `word2[0..j-1]`.

**Base cases:**
- `dp[i][0] = i` — delete i chars from word1.
- `dp[0][j] = j` — insert j chars to get word2[0..j-1].

**Recurrence:**
```
if word1[i-1] == word2[j-1]:
  dp[i][j] = dp[i-1][j-1]  // characters match; no op needed
else:
  dp[i][j] = 1 + min(
    dp[i-1][j],    // delete word1[i-1]: now must convert word1[0..i-2] to word2[0..j-1]
    dp[i][j-1],    // insert word2[j-1]: now word2[0..j-2] is the remaining target
    dp[i-1][j-1]   // replace word1[i-1] with word2[j-1]
  )
```

### Complexity
- **Time:** O(m × n).
- **Space:** O(m × n).

### Dry Run — `word1 = "horse"`, `word2 = "ros"`

|   |   | r | o | s |
|---|---|---|---|---|
|   | 0 | 1 | 2 | 3 |
| h | 1 | 1 | 2 | 3 |
| o | 2 | 2 | 1 | 2 |
| r | 3 | 2 | 2 | 2 |
| s | 4 | 3 | 3 | 2 |
| e | 5 | 4 | 4 | 3 |

`dp[5][3] = 3` ✓

---

## Approach 3 — DP Rolling Row

### Intuition
`dp[i][j]` only depends on `dp[i-1][j]`, `dp[i][j-1]`, and `dp[i-1][j-1]`. Use `prev` and `curr` rows, updating `curr` left-to-right. Before updating `curr[j]`, it still holds `prev[j]` (top neighbor); `curr[j-1]` (just updated) is the left neighbor; `prev[j-1]` would be the diagonal — save it before overwriting.

### Complexity
- **Time:** O(m × n).
- **Space:** O(n).

---

## Key Takeaways

- **The three operations map directly to the DP transitions:**
  - `dp[i-1][j]`: delete word1[i-1] (shift word1 pointer, stay on word2).
  - `dp[i][j-1]`: insert word2[j-1] (stay on word1, shift word2 pointer).
  - `dp[i-1][j-1]`: replace word1[i-1] with word2[j-1] (shift both).
- **This is one of the most important DP patterns** — used in diff tools, spell checkers, and bioinformatics (Levenshtein distance).
- **When characters match, no cost is incurred** — `dp[i][j] = dp[i-1][j-1]`, not `dp[i-1][j-1] + 0`. This is why matching is free but mismatching costs 1.

---

## Related Problems

- LeetCode #161 — One Edit Distance (check if edit distance exactly equals 1)
- LeetCode #583 — Delete Operation for Two Strings (use LCS-based DP)
- LeetCode #712 — Minimum ASCII Delete Sum (weighted edit distance)
