# 0010 — Regular Expression Matching

> LeetCode #10 · Difficulty: Hard
> **Categories:** String, Dynamic Programming, Recursion

---

## Problem Statement

Given an input string `s` and a pattern `p`, implement regular expression matching with support for `'.'` and `'*'` where:
- `'.'` Matches any single character.
- `'*'` Matches zero or more of the preceding element.

The matching must cover the **entire** input string (not partial).

**Example 1**
```
Input:  s = "aa", p = "a"
Output: false
Explanation: "a" does not match the entire string "aa".
```

**Example 2**
```
Input:  s = "aa", p = "a*"
Output: true
Explanation: '*' means zero or more of the preceding element, 'a'. "aa" matches.
```

**Example 3**
```
Input:  s = "ab", p = ".*"
Output: true
Explanation: ".*" means "zero or more (*) of any character (.)".
```

**Constraints**
- `1 <= s.length <= 20`
- `1 <= p.length <= 20`
- `s` contains only lowercase English letters.
- `p` contains only lowercase English letters, `'.'`, and `'*'`.
- It is guaranteed for each appearance of `'*'`, there will be a previous valid character to match.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Google    | ★★★★★ Very High | 2024          |
| Amazon    | ★★★★★ Very High | 2024          |
| Meta      | ★★★★☆ High      | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |
| Apple     | ★★★☆☆ Medium    | 2023          |
| Uber      | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Recursion** — Approach 1 directly encodes the two-choice rule for `*` as recursive calls.
- **Memoisation (Top-Down DP)** — Approach 2 caches `(i, j)` results to avoid exponential repeated sub-problems.
- **Bottom-Up Dynamic Programming** — Approach 3 fills a 2-D table from base cases upward, the canonical interview answer for this problem. → see [`/dsa/dynamic_programming.md`](/dsa/dynamic_programming.md) (to be created).

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Recursion (brute) | O(2^(m+n)) | O(m+n) | Understanding the rules; never in practice |
| 2 | Top-Down DP (memo) | O(m×n) | O(m×n) | Natural extension of recursion; good when the pattern is clear |
| 3 | Bottom-Up DP ✅ | O(m×n) | O(m×n) | Standard interview answer; iterative, no stack overflow risk |

---

## Approach 1 — Recursion (Brute Force)

### Intuition
Encode the matching rules directly as a recursive function over positions `(i, j)` — `i` into `s`, `j` into `p`.

The key branching point is when `p[j+1] == '*'`:
- **Zero occurrences:** skip the `x*` pair — recurse with same `i`, advance `j` by 2.
- **One or more occurrences:** if `p[j]` matches `s[i]`, consume one char from `s` — recurse with `i+1`, same `j` (allowing more matches).

If no `*` follows, simply check the current character pair and advance both.

### Algorithm
```
match(i, j):
  if j == len(p): return i == len(s)   // pattern done
  firstMatch = i < len(s) && (p[j]==s[i] || p[j]=='.')
  if j+1 < len(p) && p[j+1] == '*':
    return match(i, j+2)               // zero occurrences
        || firstMatch && match(i+1, j) // one more occurrence
  else:
    return firstMatch && match(i+1, j+1)
```

### Complexity
- **Time:** O(2^(m+n)) — every `*` doubles the branching.
- **Space:** O(m+n) — maximum recursion depth.

---

## Approach 2 — Top-Down DP (Memoisation)

### Intuition
Observation: the same `(i, j)` pair is computed many times. A 2-D memo array of size `(m+1) × (n+1)` stores `1` (true), `-1` (false), or `0` (unvisited). On each recursive call, return the cached result immediately if known.

### Complexity
- **Time:** O(m×n) — at most `(m+1)(n+1)` unique states.
- **Space:** O(m×n) — the memo table, plus O(m+n) call stack.

---

## Approach 3 — Bottom-Up DP (Recommended ✅)

### Intuition
Define `dp[i][j]` = true iff `s[i:]` matches `p[j:]`. Fill from the end of both strings to the beginning.

**Base cases:**
- `dp[m][n] = true` — both exhausted.
- `dp[m][j]` for `j < n`: only true if the remaining pattern can match empty string, i.e., all remaining chars form `x*` pairs.

**Transition (for `i` from `m-1` down to 0, `j` from `n-1` down to 0):**
```
firstMatch = (p[j] == s[i]) || (p[j] == '.')

if j+1 < n && p[j+1] == '*':
    dp[i][j] = dp[i][j+2]                    // zero occurrences of x*
             || (firstMatch && dp[i+1][j])    // one more occurrence
else:
    dp[i][j] = firstMatch && dp[i+1][j+1]
```

### Complexity
- **Time:** O(m×n).
- **Space:** O(m×n) — reducible to O(n) with two rows.

### Code
```go
func bottomUpDP(s, p string) bool {
    m, n := len(s), len(p)
    dp := make([][]bool, m+1)
    for i := range dp { dp[i] = make([]bool, n+1) }
    dp[m][n] = true
    for j := n - 2; j >= 0; j -= 2 {
        if p[j+1] == '*' { dp[m][j] = dp[m][j+2] }
    }
    for i := m - 1; i >= 0; i-- {
        for j := n - 1; j >= 0; j-- {
            firstMatch := p[j] == s[i] || p[j] == '.'
            if j+1 < n && p[j+1] == '*' {
                dp[i][j] = dp[i][j+2] || (firstMatch && dp[i+1][j])
            } else {
                dp[i][j] = firstMatch && dp[i+1][j+1]
            }
        }
    }
    return dp[0][0]
}
```

### Dry Run — `s = "aab"`, `p = "c*a*b"`
```
m=3, n=5
dp[3][5]=true

Base row (i=3):
  j=3: p[3]='a', p[4]='b' — no * → dp[3][3]=false
  j=2: p[2]='*' ← skip (we step by 2 from j=n-2=3 down)
  j=1: p[1]='*' so dp[3][1]=dp[3][3]=false
  j=0: check p[0]='c', p[1]='*' → dp[3][0]=dp[3][2]
    dp[3][2]: p[2]='a', p[3]='b' — no → false
  → dp[3][0]=false... wait pattern "c*a*b" let's redo:
    p = c * a * b  → indices 0 1 2 3 4
  j=3 (down by 2 from j=4-2=3): p[4]='b' → dp[3][3]=false (j+1=4 no *)
                                 actually loop is j=n-2=3: p[3]='*' → dp[3][3]=dp[3][5]=true!
  j=1: p[1]='*' → dp[3][1]=dp[3][3]=true

Fill:
i=2 (s[2]='b'), j=4 (p[4]='b'): firstMatch=true, no * → dp[2][4]=dp[3][5]=true
i=2, j=3 (p[3]='*'): skipped (handled at j+1 from j=2)
i=2, j=2 (p[2]='a'): p[3]='*' → dp[2][2]=dp[2][4]||(firstMatch=false&&...)=true
i=2, j=0 (p[0]='c'): p[1]='*' → dp[2][0]=dp[2][2]||(false&&...)=true
...
dp[0][0]=true ✓
```

---

## Key Takeaways

- **The `*` choice** — when you see `x*`, you always have exactly two choices: skip the pair (zero occurrences) OR consume one `s[i]` if it matches. This two-choice rule drives the entire DP.
- **Fill direction matters** — `dp[i][j]` depends on `dp[i+1][j]` and `dp[i][j+2]`, so fill bottom-to-top and right-to-left.
- **Base row for empty string** — `dp[m][j]` requires careful handling: only `x*` pairs at the tail of the pattern can match an empty string. Step by 2 (skipping pairs) when initialising.
- **`'.'` is just a wildcard** — it participates in `firstMatch` exactly like any literal character; no special treatment beyond that.

---

## Implementation (Go)

See [main.go](main.go).

### Verification
```
s="aa"   p="a"     → false ✓
s="aa"   p="a*"    → true  ✓
s="ab"   p=".*"    → true  ✓
s="aab"  p="c*a*b" → true  ✓
s=""     p="a*"    → true  ✓
```

---

## Related Problems

- LeetCode #44 — Wildcard Matching (`?` and `*`; simpler because `*` matches any sequence, not tied to a preceding char)
- LeetCode #5 — Longest Palindromic Substring (2-D DP on substrings)
- LeetCode #72 — Edit Distance (2-D DP; classical string DP)
