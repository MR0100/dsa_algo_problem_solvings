# 0087 — Scramble String

> LeetCode #87 · Difficulty: Hard
> **Categories:** String, Dynamic Programming, Recursion, Memoization

---

## Problem Statement

We can scramble a string `s` to get a string `t` using the following algorithm:

1. If the length of the string is 1, stop.
2. If the length of the string is > 1, do the following:
   - Split the string into two non-empty substrings at a random index, i.e., if the string is `s`, divide it to `x` and `y` where `s = x + y`.
   - **Randomly** decide to **swap** the two substrings or to **keep them in the same order**. i.e., after this step, `s` may become `s = x + y` or `s = y + x`.
   - Apply step 1 recursively on each of the two substrings `x` and `y`.

Given two strings `s1` and `s2` of **the same length**, return `true` if `s2` is a scrambled string of `s1`, otherwise, return `false`.

**Example 1:**
```
Input: s1 = "great", s2 = "rgeat"
Output: true
Explanation: "great" → split "gr"|"eat" → swap → "eat"+"gr" → split "e"|"at" → swap "at"+"e" → ... → "rgeat"
```

**Example 2:**
```
Input: s1 = "abcde", s2 = "caebd"
Output: false
```

**Example 3:**
```
Input: s1 = "a", s2 = "a"
Output: true
```

**Constraints:**
- `s1.length == s2.length`
- `1 <= s1.length <= 30`
- `s1` and `s2` consist of lowercase English letters.

---

## Company Frequency

| Company  | Frequency      | Last Reported |
|----------|----------------|---------------|
| Google   | ★★★☆☆ Medium   | 2024          |
| Amazon   | ★★☆☆☆ Low      | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Top-Down Memoized Recursion** — subproblems identified by the two substrings; memoize by `(a, b)` key. See [`/dsa/dynamic_programming_2d.md`](/dsa/dynamic_programming_2d.md)
- **3D Bottom-Up DP** — `dp[len][i][j]`: is `s1[i:i+len]` a scramble of `s2[j:j+len]`?
- **Interval DP** — overlapping subproblems on string intervals.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Memoized Recursion | O(n^4) | O(n^3) | Natural; easier to reason about |
| 2 | 3D Bottom-Up DP | O(n^4) | O(n^3) | Iterative; avoids recursion overhead |

---

## Approach 1 — Memoized Recursion

### Intuition
`s1` is a scramble of `s2` iff we can split `s1` at position `k` (for some `k` from 1 to n-1) and either:
- **(No-swap):** `s1[:k]` is a scramble of `s2[:k]` AND `s1[k:]` is a scramble of `s2[k:]`.
- **(Swap):** `s1[:k]` is a scramble of `s2[n-k:]` AND `s1[k:]` is a scramble of `s2[:n-k]`.

**Early exit:** if the character frequency multisets of `s1` and `s2` differ, they can't be scrambles — this prunes heavily.

**Memoization key:** `a + "#" + b` (the `#` separator prevents false matches like `"ab"+"c"` vs `"a"+"bc"`).

### Algorithm
1. Base: if `a == b` return `true`; if lengths differ return `false`.
2. Check memo; return if cached.
3. Check character frequency — if different, cache and return `false`.
4. For `k = 1` to `n-1`: try no-swap and swap splits. Cache and return `true` if any works.
5. Cache and return `false`.

### Complexity
- **Time:** O(n^4) — O(n^3) unique (a,b) pairs (bounded by O(n^3) substrings of a fixed-length string), O(n) splits per pair.
- **Space:** O(n^3) — memo table + call stack O(n).

### Code
```go
func isScramble(s1 string, s2 string) bool {
    memo := make(map[string]bool)
    var dp func(a, b string) bool
    dp = func(a, b string) bool {
        if a == b { return true }
        key := a + "#" + b
        if v, ok := memo[key]; ok { return v }
        freq := [26]int{}
        for i := 0; i < len(a); i++ { freq[a[i]-'a']++; freq[b[i]-'a']-- }
        for _, f := range freq { if f != 0 { memo[key] = false; return false } }
        n := len(a)
        for k := 1; k < n; k++ {
            if dp(a[:k], b[:k]) && dp(a[k:], b[k:]) { memo[key] = true; return true }
            if dp(a[:k], b[n-k:]) && dp(a[k:], b[:n-k]) { memo[key] = true; return true }
        }
        memo[key] = false
        return false
    }
    return dp(s1, s2)
}
```

### Dry Run (s1="great", s2="rgeat")

```
dp("great", "rgeat"):
  freq check: both have {g:1,r:1,e:1,a:1,t:1} ✓
  k=1: dp("g","r")=false; dp("g","t")=false; ...
  k=2: dp("gr","rg")? 
    dp("gr","rg"): k=1: dp("g","r")=false; dp("g","g")=true && dp("r","r")=true → true!
  dp("gr","rg")=true && dp("eat","eat")=true → return true!
```

---

## Approach 2 — 3D Bottom-Up DP

### Intuition
`dp[l][i][j]` = true if `s1[i:i+l]` is a scramble of `s2[j:j+l]`.

Fill by increasing length `l` from 1 to n. For length 1: `dp[1][i][j] = (s1[i] == s2[j])`.

For length `l`, split at `k`:
- No-swap: `dp[k][i][j] && dp[l-k][i+k][j+k]`.
- Swap: `dp[k][i][j+l-k] && dp[l-k][i+k][j]`.

### Complexity
- **Time:** O(n^4)
- **Space:** O(n^3)

### Code
```go
func isScrambleDP(s1 string, s2 string) bool {
    n := len(s1)
    dp := make([][][]bool, n+1)
    for l := 0; l <= n; l++ {
        dp[l] = make([][]bool, n)
        for i := range dp[l] { dp[l][i] = make([]bool, n) }
    }
    for i := 0; i < n; i++ {
        for j := 0; j < n; j++ { dp[1][i][j] = s1[i] == s2[j] }
    }
    for length := 2; length <= n; length++ {
        for i := 0; i <= n-length; i++ {
            for j := 0; j <= n-length; j++ {
                for k := 1; k < length; k++ {
                    if dp[k][i][j] && dp[length-k][i+k][j+k] { dp[length][i][j] = true; break }
                    if dp[k][i][j+length-k] && dp[length-k][i+k][j] { dp[length][i][j] = true; break }
                }
            }
        }
    }
    return dp[n][0][0]
}
```

### Dry Run (s1="a", s2="a")
`dp[1][0][0] = (s1[0]=='a') == (s2[0]=='a') = true`.
No length-2+ iterations needed. Return `dp[1][0][0] = true` ✓.

---

## Key Takeaways
- The "no-swap / swap" split structure is the key insight — think of the string as a binary tree where you can swap subtrees.
- Character frequency check as a pruning step dramatically reduces the actual runtime.
- The memoization key must separate the two strings (use `#`) to avoid hash collisions.
- 3D DP trades recursion stack for explicit O(n^3) memory; both are O(n^4) time.

---

## Related Problems
- LeetCode #140 — Word Break II (memoized recursion on string intervals)
- LeetCode #312 — Burst Balloons (interval DP)
- LeetCode #486 — Predict the Winner (interval DP, two-player game)
