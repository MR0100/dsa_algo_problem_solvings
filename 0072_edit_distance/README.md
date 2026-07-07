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

### Code
```go
func memoization(word1 string, word2 string) int {
	m, n := len(word1), len(word2)
	memo := make([][]int, m+1)
	for i := range memo {
		memo[i] = make([]int, n+1)
		for j := range memo[i] {
			memo[i][j] = -1
		}
	}

	var dp func(i, j int) int
	dp = func(i, j int) int {
		if i == m {
			return n - j // insert remaining chars of word2
		}
		if j == n {
			return m - i // delete remaining chars of word1
		}
		if memo[i][j] != -1 {
			return memo[i][j]
		}
		if word1[i] == word2[j] {
			memo[i][j] = dp(i+1, j+1)
		} else {
			del := dp(i+1, j)   // delete word1[i]
			ins := dp(i, j+1)   // insert word2[j] before word1[i]
			rep := dp(i+1, j+1) // replace word1[i] with word2[j]
			best := del
			if ins < best {
				best = ins
			}
			if rep < best {
				best = rep
			}
			memo[i][j] = 1 + best
		}
		return memo[i][j]
	}

	return dp(0, 0)
}
```

### Dry Run — `word1 = "horse"`, `word2 = "ros"`

`dp(i,j)` = min ops to convert `word1[i:]` to `word2[j:]`. Base rows/cols come from `i==5` (→ `3-j`) and `j==3` (→ `5-i`). Filling `memo[i][j]` (i = suffix of "horse", j = suffix of "ros"):

| dp(i,j) | word1[i]/word2[j] | rule | value |
|---------|-------------------|------|-------|
| dp(4,2) | e / s | mismatch: 1+min(dp(5,2)=1, dp(4,3)=1, dp(5,3)=0) | 1 |
| dp(3,2) | s / s | match: dp(4,3)=1 | 1 |
| dp(2,2) | r / s | mismatch: 1+min(dp(3,2)=1, dp(2,3)=2, dp(3,3)=1) | 2 |
| dp(2,1) | r / o | mismatch: 1+min(dp(3,1), dp(2,2)=2, dp(3,2)=1) | 2 |
| dp(1,1) | o / o | match: dp(2,2)=2 | 2 |
| dp(1,0) | o / r | mismatch: 1+min(dp(2,0), dp(1,1)=2, dp(2,1)=2) | 3 |
| dp(0,0) | h / r | mismatch: 1+min(dp(1,0)=3, dp(0,1), dp(1,1)=2) | 3 |

Return `dp(0,0) = 3` ✓ (replace h→r, delete r, delete e)

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

### Code
```go
func dpBottomUp(word1, word2 string) int {
	m, n := len(word1), len(word2)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
		dp[i][0] = i // cost to delete i chars
	}
	for j := 0; j <= n; j++ {
		dp[0][j] = j // cost to insert j chars
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if word1[i-1] == word2[j-1] {
				dp[i][j] = dp[i-1][j-1] // no op needed
			} else {
				best := dp[i-1][j] // delete
				if dp[i][j-1] < best {
					best = dp[i][j-1] // insert
				}
				if dp[i-1][j-1] < best {
					best = dp[i-1][j-1] // replace
				}
				dp[i][j] = 1 + best
			}
		}
	}
	return dp[m][n]
}
```

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

### Code
```go
func dpRolling(word1, word2 string) int {
	m, n := len(word1), len(word2)
	prev := make([]int, n+1)
	for j := 0; j <= n; j++ {
		prev[j] = j // dp[0][j] = j
	}
	for i := 1; i <= m; i++ {
		curr := make([]int, n+1)
		curr[0] = i // dp[i][0] = i
		for j := 1; j <= n; j++ {
			if word1[i-1] == word2[j-1] {
				curr[j] = prev[j-1]
			} else {
				best := prev[j] // delete
				if curr[j-1] < best {
					best = curr[j-1] // insert
				}
				if prev[j-1] < best {
					best = prev[j-1] // replace
				}
				curr[j] = 1 + best
			}
		}
		prev = curr
	}
	return prev[n]
}
```

### Dry Run — `word1 = "horse"`, `word2 = "ros"`

Only `prev` and `curr` rows are kept (columns index "" , r, o, s). `prev` starts as the base row `dp[0][*]`. Each iteration builds `curr` left-to-right, then `prev = curr`.

| i | char | curr[0] | curr (r) | curr (o) | curr (s) | full curr |
|---|------|---------|----------|----------|----------|-----------|
| 0 | — (base) | 0 | 1 | 2 | 3 | [0,1,2,3] |
| 1 | h | 1 | 1+min(prev1=1,curr0=1,prev0=0)=1 | 1+min(2,1,1)=2 | 1+min(3,2,2)=3 | [1,1,2,3] |
| 2 | o | 2 | 1+min(1,2,1)=2 | o==o: prev[r]=1 | 1+min(3,1,2)=2 | [2,2,1,2] |
| 3 | r | 3 | r==r: prev[""]=2 | 1+min(1,2,2)=2 | 1+min(2,2,1)=2 | [3,2,2,2] |
| 4 | s | 4 | 1+min(2,4,3)=3 | 1+min(2,3,2)=3 | s==s: prev[o]=2 | [4,3,3,2] |
| 5 | e | 5 | 1+min(3,5,4)=4 | 1+min(3,4,3)=4 | 1+min(2,4,3)=3 | [5,4,4,3] |

Return `prev[n] = prev[3] = 3` ✓ (rows match the 2D table in Approach 2)

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
