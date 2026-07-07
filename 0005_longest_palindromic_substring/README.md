# 0005 — Longest Palindromic Substring

> LeetCode #5 · Difficulty: Medium
> **Categories:** String, Dynamic Programming, Two Pointers

---

## Problem Statement

Given a string `s`, return the **longest palindromic substring** in `s`.

**Example 1**
```
Input:  s = "babad"
Output: "bab"
Explanation: "aba" is also a valid answer.
```

**Example 2**
```
Input:  s = "cbbd"
Output: "bb"
```

**Constraints**
- `1 <= s.length <= 1000`
- `s` consists of only digits and English letters.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Microsoft | ★★★★★ Very High | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Meta      | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★★☆ High      | 2024          |
| Apple     | ★★★☆☆ Medium    | 2023          |
| Adobe     | ★★★☆☆ Medium    | 2023          |
| Uber      | ★★★☆☆ Medium    | 2023          |
| Netflix   | ★★☆☆☆ Low       | 2022          |
| LinkedIn  | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **String** — the problem operates directly on character indices; all approaches require positional access into the input.
- **Dynamic Programming** — Approach 2 builds a 2-D table from sub-problems (is `s[i..j]` a palindrome?) up to the full string.
- **Two Pointers / Expand Around Center** — Approach 3 uses two pointers (`l`, `r`) expanding outward from each center while `s[l]==s[r]`.
- **Manacher's Algorithm** — Approach 4 uses a "radius array" to skip already-computed palindrome lengths, achieving O(n). → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md) for related expand patterns.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(n³) | O(1) | Never in interviews; baseline only |
| 2 | Dynamic Programming | O(n²) | O(n²) | Clear and systematic; good for teaching |
| 3 | Expand Around Center ✅ | O(n²) | O(1) | Best practical choice; simple and space-efficient |
| 4 | Manacher's Algorithm | O(n) | O(n) | When O(n²) is too slow; harder to explain in an interview |

---

## Approach 1 — Brute Force

### Intuition
Generate every possible substring `s[i..j]` (O(n²) pairs) and check each one for palindromicity by comparing characters from both ends inward (O(n)). Keep the longest valid one.

### Algorithm
1. Start with `best = s[0:1]` (any single character is a palindrome).
2. For every pair `(i, j)` with `j > i`:
   - Only check if `j-i+1 > len(best)` (prune early).
   - If `isPalindrome(s, i, j)` → update `best`.
3. Return `best`.

### Complexity
- **Time:** O(n³) — O(n²) pairs × O(n) palindrome check.
- **Space:** O(1) — only index variables.

### Code
```go
func bruteForce(s string) string {
    n := len(s)
    best := s[0:1]
    for i := 0; i < n; i++ {
        for j := i + 1; j < n; j++ {
            if j-i+1 > len(best) && isPalindrome(s, i, j) {
                best = s[i : j+1]
            }
        }
    }
    return best
}
func isPalindrome(s string, l, r int) bool {
    for l < r {
        if s[l] != s[r] { return false }
        l++; r--
    }
    return true
}
```

### Dry Run — `s = "babad"`
```
i=0:
  j=1: "ba" → b≠a → false
  j=2: "bab" → b==b, middle a → true → best="bab"
  j=3: "baba" → b≠a → false
  j=4: "babad" → b≠d → false
i=1:
  j=3: "aba" → len=3 not > 3 → skip
  ...
Result: "bab"
```

---

## Approach 2 — Dynamic Programming

### Intuition
A substring `s[i..j]` is a palindrome iff `s[i]==s[j]` AND `s[i+1..j-1]` is also a palindrome. Fill a 2-D DP table from smaller substrings to larger.

### Algorithm
1. `dp[i][i] = true` for all `i` (single characters are palindromes).
2. For each length `L` from 2 to n:
   - For each start `i`, let `j = i + L - 1`.
   - `dp[i][j] = (s[i]==s[j]) && (L==2 || dp[i+1][j-1])`.
   - If true and `L > maxLen`, update `start, maxLen`.
3. Return `s[start : start+maxLen]`.

### Complexity
- **Time:** O(n²) — fill n² cells.
- **Space:** O(n²) — the DP table.

### Code
```go
func dpApproach(s string) string {
    n := len(s)
    dp := make([][]bool, n)
    for i := range dp { dp[i] = make([]bool, n); dp[i][i] = true }
    start, maxLen := 0, 1
    for length := 2; length <= n; length++ {
        for i := 0; i <= n-length; i++ {
            j := i + length - 1
            if length == 2 {
                dp[i][j] = s[i] == s[j]
            } else {
                dp[i][j] = (s[i] == s[j]) && dp[i+1][j-1]
            }
            if dp[i][j] && length > maxLen { start, maxLen = i, length }
        }
    }
    return s[start : start+maxLen]
}
```

### Dry Run — `s = "cbbd"`
```
Base: dp[0][0]=dp[1][1]=dp[2][2]=dp[3][3]=true

L=2:
  (0,1): c≠b → false
  (1,2): b==b → dp[1][2]=true → best="bb" (start=1,maxLen=2)
  (2,3): b≠d → false

L=3:
  (0,2): c≠b → false
  (1,3): b≠d → false

L=4:
  (0,3): c≠d → false

Result: s[1:3] = "bb" ✓
```

---

## Approach 3 — Expand Around Center (Recommended ✅)

### Intuition
Every palindrome has a center. Iterate over all 2n-1 possible centers (n single characters for odd-length + n-1 gaps for even-length) and expand outward while characters match.

### Algorithm
1. For each `i` from 0 to n-1:
   - Odd-length: `expand(s, i, i)`.
   - Even-length: `expand(s, i, i+1)`.
   - `expand` moves `l--`, `r++` while `s[l]==s[r]` and in-bounds, then steps back by 1.
   - Update `[start, end]` if the new span is wider.
2. Return `s[start : end+1]`.

### Complexity
- **Time:** O(n²) — 2n-1 centers × up to O(n) expansion.
- **Space:** O(1) — only index variables.

### Code
```go
func expandAroundCenter(s string) string {
    start, end := 0, 0
    for i := 0; i < len(s); i++ {
        l1, r1 := expand(s, i, i)
        l2, r2 := expand(s, i, i+1)
        if r1-l1 > end-start { start, end = l1, r1 }
        if r2-l2 > end-start { start, end = l2, r2 }
    }
    return s[start : end+1]
}
func expand(s string, l, r int) (int, int) {
    for l >= 0 && r < len(s) && s[l] == s[r] { l--; r++ }
    return l + 1, r - 1
}
```

### Dry Run — `s = "racecar"`
```
i=3 (center 'e'):
  odd: expand(3,3)
    l=3,r=3: 'e'=='e' → l=2,r=4
    l=2,r=4: 'c'=='c' → l=1,r=5
    l=1,r=5: 'a'=='a' → l=0,r=6
    l=0,r=6: 'r'=='r' → l=-1,r=7 → OOB
    return (0, 6) → span=6 → "racecar"
  even: expand(3,4): 'e'≠'c' → return (4,3) → r<l, ignored

Result: s[0:7] = "racecar" ✓
```

---

## Approach 4 — Manacher's Algorithm

### Intuition
Transform `s` into `T` by inserting `#` between every character (e.g. `"abc"` → `"#a#b#c#"`). All palindromes in `T` are odd-length. Maintain `p[i]` = palindrome radius at `T[i]`. When `T[i]` is inside the rightmost known palindrome `[center, right]`, its mirror's radius gives a free starting point, skipping redundant comparisons.

### Algorithm
1. Build `T = "#" + interleave(s, "#")`.
2. For each `i`: seed `p[i]` from mirror if `i < right`, then expand.
3. After expanding, update `[center, right]` if `i+p[i] > right`.
4. Find `maxRadius` → `start = (maxCenter - maxRadius) / 2`.

### Complexity
- **Time:** O(n) — `right` only increases; total expansions ≤ 2·|T|.
- **Space:** O(n) — `T` and `p`.

### Code
```go
func manacher(s string) string {
    t := "#"
    for _, ch := range s { t += string(ch) + "#" }
    n := len(t)
    p := make([]int, n)
    center, right := 0, 0
    for i := 0; i < n; i++ {
        mirror := 2*center - i
        if i < right { p[i] = min(right-i, p[mirror]) }
        l, r := i-p[i]-1, i+p[i]+1
        for l >= 0 && r < n && t[l] == t[r] { p[i]++; l--; r++ }
        if i+p[i] > right { center = i; right = i + p[i] }
    }
    maxRadius, maxCenter := 0, 0
    for i, radius := range p {
        if radius > maxRadius { maxRadius, maxCenter = radius, i }
    }
    start := (maxCenter - maxRadius) / 2
    return s[start : start+maxRadius]
}
```

### Dry Run — `s = "babad"`, `T = "#b#a#b#a#d#"`
```
i=3 (T[3]='a'):
  Expand: T[2]='#'==T[4]='#' → p[3]=1
          T[1]='b'==T[5]='b' → p[3]=2
          T[0]='#'==T[6]='#' → p[3]=3
          T[-1] OOB → stop. right=6, center=3

i=5 (T[5]='b'):
  i<right → mirror=1, p[1]=1 → p[5]=min(1,1)=1
  Expand: T[3]='a'==T[7]='a' → p[5]=2
          T[2]='#'==T[8]='#' → p[5]=3
          T[1]='b'≠T[9]='d' → stop. right=8, center=5

max: p[3]=3 and p[5]=3 (tie at i=3)
start = (3-3)/2 = 0 → s[0:3] = "bab" ✓
```

---

## Key Takeaways

- **The center insight** — every palindrome has a center (a character or a gap). This transforms brute-force O(n³) into an elegant O(n²) by iterating over all 2n-1 centers.
- **Expand Around Center vs DP** — both are O(n²) time, but EAC wins on space: O(1) vs O(n²). In interviews, always prefer EAC.
- **Manacher's mirror trick** — the big O(n) win comes from re-using previously computed radii. When `i` is inside a known palindrome, its mirror's radius gives a free lower bound — no wasted expansions.
- **T transformation** — inserting `#` unifies odd and even palindromes so Manacher's needs only one code path.
- **Interview recommendation** — Implement Expand Around Center. Mention Manacher's exists; only code it if explicitly asked or if the interviewer pushes for O(n).

---

## Implementation (Go)

See [main.go](main.go).

```go
// (full code in main.go)
```

### Verification
```
s="babad"
  Approach 1: Brute Force             O(n³) T | O(1)   S       → "bab"
  Approach 2: Dynamic Programming     O(n²) T | O(n²)  S       → "bab"
  Approach 3: Expand Around Center ✅ O(n²) T | O(1)   S        → "bab"
  Approach 4: Manacher's Algorithm    O(n)  T | O(n)   S       → "bab"

s="cbbd"
  Approach 1: Brute Force             O(n³) T | O(1)   S       → "bb"
  Approach 2: Dynamic Programming     O(n²) T | O(n²)  S       → "bb"
  Approach 3: Expand Around Center ✅ O(n²) T | O(1)   S        → "bb"
  Approach 4: Manacher's Algorithm    O(n)  T | O(n)   S       → "bb"

s="racecar"
  Approach 1: Brute Force             O(n³) T | O(1)   S       → "racecar"
  Approach 2: Dynamic Programming     O(n²) T | O(n²)  S       → "racecar"
  Approach 3: Expand Around Center ✅ O(n²) T | O(1)   S        → "racecar"
  Approach 4: Manacher's Algorithm    O(n)  T | O(n)   S       → "racecar"
```

---

## Related Problems

- LeetCode #647 — Palindromic Substrings (count all palindromic substrings — same expand-around-center technique)
- LeetCode #9 — Palindrome Number (check if a number reads the same forward/backward)
- LeetCode #125 — Valid Palindrome (check if a string is palindrome ignoring case/spaces)
- LeetCode #516 — Longest Palindromic Subsequence (subsequence, not substring — DP)
- LeetCode #214 — Shortest Palindrome (palindrome construction using KMP or Manacher's)
