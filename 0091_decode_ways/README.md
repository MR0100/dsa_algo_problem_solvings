# 0091 — Decode Ways

> LeetCode #91 · Difficulty: Medium
> **Categories:** String, Dynamic Programming

---

## Problem Statement

A message containing letters from `A-Z` can be **encoded** into numbers using:
- `'A' -> "1"`, `'B' -> "2"`, ..., `'Z' -> "26"`

To **decode** an encoded message, all the digits must be grouped and then mapped back to letters using the reverse of the mapping above (there may be multiple ways).

Given a string `s` containing only digits, return the **number of ways** to decode it.

The test cases are generated such that the answer fits in a 32-bit integer.

**Example 1:**
```
Input: s = "12"
Output: 2
Explanation: "12" could be decoded as "AB" (1 2) or "L" (12).
```

**Example 2:**
```
Input: s = "226"
Output: 3
Explanation: "226" could be decoded as "BZ" (2 26), "VF" (22 6), or "BBF" (2 2 6).
```

**Example 3:**
```
Input: s = "06"
Output: 0
Explanation: "06" cannot be mapped to "F" because of the leading zero.
```

**Constraints:**
- `1 <= s.length <= 100`
- `s` contains only digits.
- `s` does not contain any leading zeros.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Facebook  | ★★★★☆ High      | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Google    | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Dynamic Programming (Linear)** — `dp[i]` = number of ways to decode `s[0:i]`. See [`/dsa/dynamic_programming.md`](/dsa/dynamic_programming.md)
- **Fibonacci-like recurrence** — each position depends on the previous 1 or 2 positions.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Memoized Recursion | O(n) | O(n) | Good starting point; intuitive |
| 2 | Bottom-Up DP | O(n) | O(n) | Standard iterative DP |
| 3 | O(1) Space DP | O(n) | O(1) | Optimal — only two previous values needed |

---

## Approach 1 — Memoized Recursion

### Intuition
`dp(i)` = number of ways to decode `s[i:]`.
- If `s[i] == '0'`: 0 ways (can't start with '0').
- Take 1 digit: always valid if `s[i] != '0'`. Add `dp(i+1)`.
- Take 2 digits `s[i:i+2]` if `10 <= val <= 26`. Add `dp(i+2)`.

### Algorithm
Recursive + memo. Base case: `dp(n) = 1` (empty suffix, 1 way).

### Complexity
- **Time:** O(n) — n unique states, O(1) work each.
- **Space:** O(n) — memo + call stack.

### Code
```go
func numDecodings(s string) int {
    memo := make(map[int]int)
    var dp func(i int) int
    dp = func(i int) int {
        if i == len(s) { return 1 }
        if s[i] == '0' { return 0 }
        if v, ok := memo[i]; ok { return v }
        result := dp(i + 1)
        if i+1 < len(s) {
            twoDigit := int(s[i]-'0')*10 + int(s[i+1]-'0')
            if twoDigit >= 10 && twoDigit <= 26 { result += dp(i + 2) }
        }
        memo[i] = result
        return result
    }
    return dp(0)
}
```

### Dry Run (s="226")

```
dp(0) = dp(1) + dp(2) [since 22<=26]
dp(1) = dp(2) + dp(3) [since 26<=26]
dp(2) = dp(3)         [6 only; no 2-digit 6_]
dp(3) = 1             [end of string]
dp(2) = 1
dp(1) = 1 + 1 = 2
dp(0) = 2 + 1 = 3
```

---

## Approach 2 — Bottom-Up DP

### Intuition
`dp[i]` = number of ways to decode `s[0:i]`.
- `dp[0] = 1` (empty string: one way).
- `dp[1] = 1` if `s[0] != '0'`, else 0.
- `dp[i]`:
  - If `s[i-1] != '0'`: add `dp[i-1]` (take 1 digit).
  - If `10 <= s[i-2:i] <= 26`: add `dp[i-2]` (take 2 digits).

### Complexity
- **Time:** O(n)
- **Space:** O(n)

### Code
```go
func numDecodingsDP(s string) int {
    n := len(s)
    dp := make([]int, n+1)
    dp[0] = 1
    if s[0] != '0' { dp[1] = 1 }
    for i := 2; i <= n; i++ {
        if s[i-1] != '0' { dp[i] += dp[i-1] }
        twoDigit := int(s[i-2]-'0')*10 + int(s[i-1]-'0')
        if twoDigit >= 10 && twoDigit <= 26 { dp[i] += dp[i-2] }
    }
    return dp[n]
}
```

### Dry Run (s="226", dp=[1,1,_,_])

| i | s[i-1] | 1-digit? | s[i-2:i] | 2-digit? | dp[i] |
|---|--------|----------|----------|----------|-------|
| 2 | '2'≠0 | dp[1]=1 | 22 in [10,26] | dp[0]=1 | 2 |
| 3 | '6'≠0 | dp[2]=2 | 26 in [10,26] | dp[1]=1 | 3 |

dp[3] = 3 ✓

---

## Approach 3 — O(1) Space DP

### Intuition
Only the previous two `dp` values are needed. Replace the array with two variables `prev2` (=`dp[i-2]`) and `prev1` (=`dp[i-1]`).

### Complexity
- **Time:** O(n)
- **Space:** O(1)

### Code
```go
func numDecodingsO1(s string) int {
    prev2, prev1 := 1, 0
    if s[0] != '0' { prev1 = 1 }
    for i := 2; i <= len(s); i++ {
        curr := 0
        if s[i-1] != '0' { curr += prev1 }
        two := int(s[i-2]-'0')*10 + int(s[i-1]-'0')
        if two >= 10 && two <= 26 { curr += prev2 }
        prev2, prev1 = prev1, curr
    }
    return prev1
}
```

---

## Key Takeaways
- `'0'` cannot be decoded as a single digit — it's only valid as part of `"10"` or `"20"`.
- Two-digit validity: `10 <= val <= 26`. Check the lower bound (≥10) to exclude leading zeros.
- This is essentially the Fibonacci recurrence with validity constraints at each step.
- Pattern: `dp[i]` from `dp[i-1]` (1 digit) and `dp[i-2]` (2 digits) — exactly like Climbing Stairs.

---

## Related Problems
- LeetCode #639 — Decode Ways II (wildcard `*` character)
- LeetCode #70 — Climbing Stairs (same Fibonacci-like DP)
- LeetCode #198 — House Robber (linear DP with 2-step skip)
