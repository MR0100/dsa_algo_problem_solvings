# 0032 — Longest Valid Parentheses

> LeetCode #32 · Difficulty: Hard
> **Categories:** String, Dynamic Programming, Stack

---

## Problem Statement

Given a string containing just the characters `'('` and `')'`, return the length of the longest valid (well-formed) parentheses substring.

**Example 1**
```
Input:  s = "(()"
Output: 2
Explanation: The longest valid parentheses substring is "()".
```

**Example 2**
```
Input:  s = ")()())"
Output: 4
Explanation: The longest valid parentheses substring is "()()".
```

**Example 3**
```
Input:  s = ""
Output: 0
```

**Constraints**
- `0 <= s.length <= 3 * 10⁴`
- `s[i]` is `'('` or `')'`.

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

- **Stack** — push indices to track unmatched characters; the top of the stack is always the base before the current valid run.
- **Dynamic Programming** — `dp[i]` = length of the longest valid substring ending at index i. Transition requires looking back by `dp[i-1]+1` to find the potential match.
- **Two Pointers** — the two-pass counter approach uses left/right counts and two scans to cover all cases in O(1) space.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(n³) | O(1) | Never; TLE |
| 2 | Stack | O(n) | O(n) | Clearest to explain in interviews |
| 3 | Dynamic Programming | O(n) | O(n) | Good follow-up when asked for DP solution |
| 4 | Two-Pass Counters ✅ | O(n) | O(1) | Optimal; satisfies O(1)-space follow-up |

---

## Approach 1 — Brute Force

### Intuition
Try every even-length substring and validate it with a balance counter. O(n³) — too slow.

### Complexity
- **Time:** O(n³).
- **Space:** O(1).

---

## Approach 2 — Stack

### Intuition
Push indices onto a stack. The stack always holds the index of the **last unmatched character** at its top — this acts as the "base" from which we measure the length of the current valid run.

Initialize with `-1` as a sentinel.

- `'('` → push its index.
- `')'` → pop the top:
  - If stack is empty after pop: push current index (new unmatched base).
  - Else: `current_length = i - stack.top`.

### Dry Run — `s = ")()())"`, target = 4
```
stack=[-1]
i=0 ')': pop → stack=[]. Empty → push 0. stack=[0]
i=1 '(': push 1. stack=[0,1]
i=2 ')': pop → stack=[0]. length=2-0=2. best=2
i=3 '(': push 3. stack=[0,3]
i=4 ')': pop → stack=[0]. length=4-0=4. best=4
i=5 ')': pop → stack=[]. Empty → push 5. stack=[5]
Result: 4 ✓
```

### Complexity
- **Time:** O(n).
- **Space:** O(n) — the stack.

---

## Approach 3 — Dynamic Programming

### Intuition
`dp[i]` = length of the longest valid substring ending exactly at index `i`.

If `s[i] == '('`: `dp[i] = 0` (valid substring can't end with `(`).

If `s[i] == ')'`:
- Let `j = i - dp[i-1] - 1` — the index just before the current valid run ending at `i-1`.
- If `s[j] == '('`: `dp[i] = dp[i-1] + 2 + (j>0 ? dp[j-1] : 0)`.
  - `dp[i-1] + 2`: current pair + the run to its left.
  - `dp[j-1]`: valid substring that ended just before the matching `(`.

### Dry Run — `s = "(()"`
```
dp = [0, 0, 0]
i=1: s[1]=')'. j=1-dp[0]-1=0. s[0]='('→ dp[1]=dp[0]+2=2. best=2
i=2: s[2]='('. dp[2]=0.
Result: 2 ✓
```

### Complexity
- **Time:** O(n).
- **Space:** O(n).

---

## Approach 4 — Two-Pass Counters (Recommended ✅)

### Intuition
Scan left→right:
- Increment `open` on `(`, `close` on `)`.
- `open == close` → found a valid substring of length `2*close`. Update best.
- `close > open` → this `)` can never be matched; reset both to 0.

This catches all valid substrings where `close` catches up to `open`. But it misses cases like `"(()"` where the excess `(` is never balanced (forward pass never resets, so `close` never equals `open`).

Fix: scan **right→left** with roles swapped — reset when `open > close`.

### Code
```go
func twoPass(s string) int {
    best := 0
    open, close := 0, 0
    for _, ch := range s {
        if ch == '(' { open++ } else { close++ }
        if open == close { best = max(best, 2*close) } else if close > open { open, close = 0, 0 }
    }
    open, close = 0, 0
    for i := len(s) - 1; i >= 0; i-- {
        if s[i] == '(' { open++ } else { close++ }
        if open == close { best = max(best, 2*open) } else if open > close { open, close = 0, 0 }
    }
    return best
}
```

### Complexity
- **Time:** O(n) — two passes.
- **Space:** O(1) — four integer variables.

### Dry Run — `s = "(()"`, forward pass
```
open=0,close=0
'(': open=1. 1≠0, 0<1 → ok
'(': open=2. 2≠0 → ok
')': close=1. 2≠1 → ok
End: no match found (close never equals open since extra '(' unmatched)

Backward pass (right→left):
')': close=1. 0<1 → ok
'(': open=1. 1==1 → best=2
'(': open=2. 2>1 → reset.
Result: 2 ✓
```

---

## Key Takeaways

- **Stack sentinel `-1`** — the base before any valid string prevents an empty-stack check inside the loop; every pop that reveals the base directly gives length `i - (-1) = i+1` (or `i - stack.top`).
- **DP look-back index** — `j = i - dp[i-1] - 1` is the key formula: `dp[i-1]` tells us how far the current run extends to the left of `i`; subtract one more to reach the potential matching `(`.
- **Two passes cover each other's blind spots** — left→right catches excess `)` resets; right→left catches excess `(` resets.
- **O(1) space is the hardest follow-up** — if asked, the two-pass counter approach is the only O(1)-space O(n)-time solution.

---

## Related Problems

- LeetCode #20 — Valid Parentheses (validation, not length)
- LeetCode #22 — Generate Parentheses (generation)
- LeetCode #301 — Remove Invalid Parentheses (minimum removals)
