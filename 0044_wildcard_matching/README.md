# 0044 — Wildcard Matching

> LeetCode #44 · Difficulty: Hard
> **Categories:** String, Dynamic Programming, Greedy, Recursion

---

## Problem Statement

Given an input string `s` and a pattern `p`, implement wildcard pattern matching with support for `'?'` and `'*'` where:

- `'?'` Matches any single character.
- `'*'` Matches any sequence of characters (including the empty sequence).

The matching should cover the **entire** input string (not partial).

**Example 1**
```
Input:  s = "aa", p = "a"
Output: false
Explanation: "a" does not match the entire string "aa".
```

**Example 2**
```
Input:  s = "aa", p = "*"
Output: true
Explanation: '*' matches any sequence.
```

**Example 3**
```
Input:  s = "cb", p = "?a"
Output: false
Explanation: '?' matches 'c', but 'a' does not match 'b'.
```

**Example 4**
```
Input:  s = "adceb", p = "*a*b"
Output: true
Explanation: The first '*' matches "", 'a' matches 'a', second '*' matches "dce", 'b' matches 'b'.
```

**Constraints**
- `0 <= s.length, p.length <= 2000`
- `s` contains only lowercase English letters.
- `p` contains only lowercase English letters, `'?'` or `'*'`.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★★ Very High | 2024          |
| Meta      | ★★★★☆ High      | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |
| Apple     | ★★★☆☆ Medium    | 2023          |
| Adobe     | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Dynamic Programming** — `dp[i][j]` = can `s[0..i-1]` be matched by `p[0..j-1]`? Transitions cover the three pattern characters: literal, `?`, `*`.
- **Two Pointers with Bookmarking** — the greedy approach tracks the last `*` position and where in `s` it matched; extends when mismatches occur.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Recursion | O(2^(n+m)) | O(n+m) | TLE; shows the branching structure |
| 2 | DP Bottom-Up ✅ | O(n × m) | O(n × m) | Standard optimal interview answer |
| 3 | Two Pointers with Bookmark ✅ | O(n × m) worst | O(1) | O(1) space; optimal when '*' is rare |

n = len(s), m = len(p).

---

## Approach 1 — Recursion (Brute Force)

### Intuition
Match character by character. On `*`, branch: match zero chars (advance `p`) or one char (advance `s`).

### Complexity
- **Time:** O(2^(n+m)) — exponential branching on `*`.
- **Space:** O(n+m).

### Code
```go
// recursion solves Wildcard Matching with pure recursion.
//
// Time:  O(2^(len(s)+len(p))) — exponential due to '*' branching
// Space: O(len(s)+len(p)) — recursion stack
func recursion(s, p string) bool {
    if len(p) == 0 {
        return len(s) == 0
    }
    if p[0] == '*' {
        return recursion(s, p[1:]) || // '*' matches zero chars
            (len(s) > 0 && recursion(s[1:], p)) // '*' matches one or more
    }
    if len(s) > 0 && (p[0] == '?' || p[0] == s[0]) {
        return recursion(s[1:], p[1:])
    }
    return false
}
```

### Dry Run — `s = "aa"`, `p = "*"`

Each call peels one character. On `*`: try matching zero chars (`recursion(s, p[1:])`) OR one-or-more (`recursion(s[1:], p)`).

| Call | Branch taken | Result |
|------|--------------|--------|
| `recursion("aa","*")` | p[0]=='*': try `("aa","")` then `("a","*")` | see below → true |
| `recursion("aa","")` | len(p)==0, len(s)!=0 | false |
| `recursion("a","*")` | p[0]=='*': try `("a","")` then `("","*")` | see below → true |
| `recursion("a","")` | len(p)==0, len(s)!=0 | false |
| `recursion("","*")` | p[0]=='*': try `("","")` | true |
| `recursion("","")` | len(p)==0 and len(s)==0 | true |

The `("","")` base case returns true, which bubbles up: `("","*")` → true, `("a","*")` → true, `("aa","*")` → true.

Result: `true` ✓ — `*` consumed both characters via successive "match one more" branches.

---

## Approach 2 — DP Bottom-Up (Recommended ✅)

### Intuition
`dp[i][j]` = true if `s[0..i-1]` matches `p[0..j-1]`.

**Base cases:**
- `dp[0][0] = true` — empty matches empty.
- `dp[0][j] = dp[0][j-1]` if `p[j-1] == '*'` — leading `*`s match empty string.

**Transitions:**
- `p[j-1] == '*'`: `dp[i][j] = dp[i][j-1]` (match 0 chars) `|| dp[i-1][j]` (match 1 more char of s).
- `p[j-1] == '?'` or `p[j-1] == s[i-1]`: `dp[i][j] = dp[i-1][j-1]`.
- Otherwise: `dp[i][j] = false`.

### Complexity
- **Time:** O(n × m).
- **Space:** O(n × m) — reducible to O(n) with a rolling array.

### Code
```go
// dpBottomUp solves Wildcard Matching using a 2D DP table.
//
// Time:  O(len(s) * len(p))
// Space: O(len(s) * len(p)) — can be reduced to O(len(s)) with rolling array
func dpBottomUp(s, p string) bool {
    m, n := len(s), len(p)
    dp := make([][]bool, m+1)
    for i := range dp {
        dp[i] = make([]bool, n+1)
    }
    dp[0][0] = true // empty matches empty

    // empty string can only match all-'*' pattern
    for j := 1; j <= n; j++ {
        if p[j-1] == '*' {
            dp[0][j] = dp[0][j-1]
        }
    }

    for i := 1; i <= m; i++ {
        for j := 1; j <= n; j++ {
            if p[j-1] == '*' {
                // '*' matches zero chars from s  → dp[i][j-1]
                // '*' matches one more char of s → dp[i-1][j]
                dp[i][j] = dp[i][j-1] || dp[i-1][j]
            } else if p[j-1] == '?' || p[j-1] == s[i-1] {
                dp[i][j] = dp[i-1][j-1]
            }
        }
    }
    return dp[m][n]
}
```

### Dry Run — `s="adceb"`, `p="*a*b"` (abbreviated)
```
dp[0][0]=T, dp[0][1]=T ('*'), dp[0][2]=F ('a'), dp[0][3..4]=F

Build row by row... dp[5][4] = true ✓
```

---

## Approach 3 — Two Pointers with Star Bookmark

### Intuition
Walk `s` and `p` with pointers `i` and `j`. When `*` is seen, record (`starIdx=j, match=i`) and assume it matches zero chars. On mismatch:
- If no `*` bookmarked: return false.
- If `*` bookmarked: extend it by one more char: `match++; i=match; j=starIdx+1`.

### Complexity
- **Time:** O(n × m) worst case (many backtracks); O(n+m) for simple patterns.
- **Space:** O(1).

### Code
```go
func twoPointers(s, p string) bool {
    i, j := 0, 0; starIdx, match := -1, 0
    for i < len(s) {
        if j < len(p) && (p[j] == '?' || p[j] == s[i]) { i++; j++ } else
        if j < len(p) && p[j] == '*' { starIdx = j; match = i; j++ } else
        if starIdx != -1 { match++; i = match; j = starIdx+1 } else { return false }
    }
    for j < len(p) && p[j] == '*' { j++ }
    return j == len(p)
}
```

### Dry Run — `s="aa"`, `p="*"`
```
i=0,j=0: p[0]='*' → starIdx=0, match=0, j=1
i=0,j=1: j==len(p)=1. mismatch. starIdx!=-1: match=1, i=1, j=1
i=1,j=1: j==len(p)=1. mismatch. starIdx!=-1: match=2, i=2, j=1
i=2: loop exits (i==len(s))
j=1<len(p): p[1] ∉ loop (j already at 1)
return j==len(p)=1? No→ Wait: trailing '*' consumed. j=1==len(p)=1 → return true ✓
```

---

## Key Takeaways

- **Wildcard `*` vs Regex `.*`** — in #10 (regular expression matching), `*` needs the preceding element to match; here `*` independently matches anything. The DP transitions differ: `dp[i][j] = dp[i][j-1] || dp[i-1][j]` for wildcard `*` vs `dp[i][j] = dp[i][j-2] || (firstMatch && dp[i-1][j])` for regex `.*`.
- **Star bookmark = greedy retry** — when `*` is seen, assume it matches zero. Only extend if forced by a mismatch. This greedy works because the last `*` is always the best fallback.
- **DP rolling array** — `dp[i]` only depends on `dp[i-1]`, so O(n) space suffices with careful update ordering.

---

## Related Problems

- LeetCode #10 — Regular Expression Matching (`.` and `*` with preceding-element semantics)
- LeetCode #28 — Find the Index of the First Occurrence in a String (pattern matching)
