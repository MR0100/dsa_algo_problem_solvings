# 0076 — Minimum Window Substring

> LeetCode #76 · Difficulty: Hard
> **Categories:** Hash Table, String, Sliding Window, Two Pointers

---

## Problem Statement

Given two strings `s` and `t` of lengths `m` and `n`, return the **minimum window substring** of `s` such that every character in `t` (including duplicates) is included in the window. If there is no such substring, return the empty string `""`.

The testcases will be generated such that the answer is **unique**.

**Example 1:**
```
Input: s = "ADOBECODEBANC", t = "ABC"
Output: "BANC"
Explanation: The minimum window substring "BANC" includes 'A', 'B', and 'C' from string t.
```

**Example 2:**
```
Input: s = "a", t = "a"
Output: "a"
Explanation: The entire string s is the minimum window.
```

**Example 3:**
```
Input: s = "a", t = "aa"
Output: ""
Explanation: Both 'a's from t must be included in the window. Since the largest window of s only has one 'a', return empty string.
```

**Constraints:**
- `m == s.length`
- `n == t.length`
- `1 <= m, n <= 10^5`
- `s` and `t` consist of uppercase and lowercase English letters.

**Follow-up:** Could you find an algorithm that runs in `O(m + n)` time?

---

## Company Frequency

| Company    | Frequency       | Last Reported |
|------------|-----------------|---------------|
| Facebook   | ★★★★★ Very High | 2024          |
| Amazon     | ★★★★☆ High      | 2024          |
| Google     | ★★★★☆ High      | 2024          |
| Microsoft  | ★★★☆☆ Medium    | 2024          |
| Bloomberg  | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Sliding Window** — expand right to find valid windows, shrink left to minimize. See [`/dsa/sliding_window.md`](/dsa/sliding_window.md)
- **Hash Map / Frequency Count** — track character requirements and current window character counts.
- **Two Pointers** — left and right define the current window boundary.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(n² × \|t\|) | O(\|Σ\|) | Never in interview; understanding only |
| 2 | Sliding Window | O(\|s\| + \|t\|) | O(\|Σ\|) | Always — optimal and expected |

---

## Approach 1 — Brute Force

### Intuition
Try every possible substring of `s`. For each starting index `i`, extend the right boundary `j` until the window contains all characters of `t`. Record the shortest valid window found.

### Algorithm
1. Build a frequency map `need` from `t`.
2. For each `i` from 0 to n-1:
   - Maintain a running window frequency map `have`.
   - For each `j` from `i` to n-1:
     - Add `s[j]` to `have`.
     - If `have` satisfies `need` for all characters: update best, break (extending further can't help from the same `i`).

### Complexity
- **Time:** O(n² × |t|) — for each of O(n²) substrings we check O(|t|) required chars.
- **Space:** O(|Σ|) — frequency maps.

### Code
```go
func bruteForce(s string, t string) string {
    if len(s) == 0 || len(t) == 0 {
        return ""
    }
    need := make(map[byte]int)
    for i := 0; i < len(t); i++ {
        need[t[i]]++
    }
    best := ""
    for i := 0; i < len(s); i++ {
        have := make(map[byte]int)
        for j := i; j < len(s); j++ {
            have[s[j]]++
            valid := true
            for ch, cnt := range need {
                if have[ch] < cnt {
                    valid = false
                    break
                }
            }
            if valid {
                if best == "" || j-i+1 < len(best) {
                    best = s[i : j+1]
                }
                break
            }
        }
    }
    return best
}
```

### Dry Run (Example 1: s="ADOBECODEBANC", t="ABC")

| i | j | window | have | valid? | best |
|---|---|--------|------|--------|------|
| 0 | 0 | A | {A:1} | No | — |
| 0 | 1 | AB | {A:1,B:1} | No | — |
| 0 | 2 | ABO | ... | No | — |
| 0 | 5 | ADOBEC | {A:1,B:1,C:1,...} | Yes | "ADOBEC" (len 6) |
| 1 | 9 | DOBECODEBA | ... | Yes | "DOBECODEBA"? No, shorter not found yet |
| ... | | | | | |
| 9 | 12 | BANC | {B:1,A:1,N:1,C:1} | Yes | "BANC" (len 4) ← best |

Final: `"BANC"`

---

## Approach 2 — Sliding Window (Optimal)

### Intuition
Instead of restarting from scratch for each `i`, maintain a live window `[l, r]`. Expand `r` to include characters. Track how many distinct required characters are currently satisfied (`formed`). Once all are satisfied (`formed == required`), try shrinking `l` to find the minimum window — shrinking continues until the window becomes invalid.

Key insight: a character `c`'s requirement is *satisfied* the moment `windowCnt[c] == need[c]`. Incrementing `formed` only at that exact threshold (not for every addition) avoids double-counting duplicates.

### Algorithm
1. Build `need` from `t`; `required = len(need)`.
2. `l = 0, formed = 0, windowCnt = {}`.
3. For each `r`:
   - Add `s[r]` to `windowCnt`.
   - If `windowCnt[s[r]] == need[s[r]]`: `formed++`.
   - While `formed == required`:
     - Update best window if `r-l+1 < bestLen`.
     - Remove `s[l]` from `windowCnt`; if it drops below its requirement: `formed--`.
     - `l++`.

### Complexity
- **Time:** O(|s| + |t|) — each character is added and removed at most once.
- **Space:** O(|Σ|) — frequency maps; |Σ| ≤ 128 for ASCII.

### Code
```go
func slidingWindow(s string, t string) string {
    need := make(map[byte]int)
    for i := 0; i < len(t); i++ {
        need[t[i]]++
    }
    required := len(need)

    windowCnt := make(map[byte]int)
    formed := 0
    l := 0
    bestLen := -1
    bestL, bestR := 0, 0

    for r := 0; r < len(s); r++ {
        ch := s[r]
        windowCnt[ch]++
        if need[ch] > 0 && windowCnt[ch] == need[ch] {
            formed++
        }
        for formed == required {
            if bestLen == -1 || r-l+1 < bestLen {
                bestLen = r - l + 1
                bestL, bestR = l, r
            }
            lch := s[l]
            windowCnt[lch]--
            if need[lch] > 0 && windowCnt[lch] < need[lch] {
                formed--
            }
            l++
        }
    }
    if bestLen == -1 {
        return ""
    }
    return s[bestL : bestR+1]
}
```

### Dry Run (Example 1: s="ADOBECODEBANC", t="ABC", need={A:1,B:1,C:1})

| r | s[r] | formed | window | Action |
|---|------|--------|--------|--------|
| 0 | A | 1 | A | A satisfies need[A]=1 |
| 1 | D | 1 | AD | — |
| 2 | O | 1 | ADO | — |
| 3 | B | 2 | ADOB | B satisfies need[B]=1 |
| 4 | E | 2 | ADOBE | — |
| 5 | C | 3 | ADOBEC | C satisfies need[C]=1; formed==3 → shrink |
| shrink l=0 | A | → remove A; formed=2 | bestLen=6 "ADOBEC" |
| 6 | O | 2 | DOBECO | — |
| 7 | D | 2 | DOBECOD | — |
| 8 | E | 2 | DOBECODE | — |
| 9 | B | 2 | DOBECODEБ | B: windowCnt=2 ≠ need[B]=1, no formed++ |
| 10 | A | 3 | ...CODEBA | A satisfies → formed=3; shrink |
| shrink l=1..9: remove D,O,B,E,C,O,D,E until B removed (formed=2) | bestLen=10? no, window "ODEBANC"? |
| ... continue | |
| 11 | N | 2 | — | — |
| 12 | C | 3 | BANC | formed=3 → shrink; bestLen=4 "BANC" |

Final: `"BANC"`

---

## Key Takeaways
- The `formed` counter elegantly tracks how many requirements are currently met without iterating the full `need` map every step.
- Only increment `formed` at the exact threshold (`windowCnt[c] == need[c]`), not for every occurrence.
- The inner `while` loop runs at most O(n) total across the entire outer loop — so overall is O(n), not O(n²).
- This pattern (expand right → contract left) applies to any "minimum window satisfying condition" problem.

---

## Related Problems
- LeetCode #3 — Longest Substring Without Repeating Characters (sliding window)
- LeetCode #159 — Fruit Into Baskets (sliding window with at most K distinct)
- LeetCode #567 — Permutation in String (fixed-size sliding window)
- LeetCode #438 — Find All Anagrams in a String (fixed-size sliding window)
