# 0020 — Valid Parentheses

> LeetCode #20 · Difficulty: Easy
> **Categories:** String, Stack

---

## Problem Statement

Given a string `s` containing just the characters `'('`, `')'`, `'{'`, `'}'`, `'['` and `']'`, determine if the input string is valid.

An input string is valid if:
1. Open brackets must be closed by the same type of brackets.
2. Open brackets must be closed in the correct order.
3. Every close bracket has a corresponding open bracket of the same type.

**Example 1**
```
Input:  s = "()"
Output: true
```

**Example 2**
```
Input:  s = "()[]{}"
Output: true
```

**Example 3**
```
Input:  s = "(]"
Output: false
```

**Constraints**
- `1 <= s.length <= 10⁴`
- `s` consists of parentheses only `'()[]{}'`.

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
| Uber      | ★★★☆☆ Medium    | 2023          |
| Netflix   | ★★★☆☆ Medium    | 2023          |
| LinkedIn  | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Stack (LIFO)** — the last opened bracket must be the first closed. A stack naturally enforces this: push openers, pop and verify on closers.

---

## Approaches Overview

| # | Approach | Time | Space | Notes |
|---|----------|------|-------|-------|
| 1 | Stack ✅ | O(n) | O(n) | Correct for all bracket types; the only valid general solution |
| 2 | Counter (educational) ⚠️ | O(n) | O(1) | Correct ONLY for a single bracket type; shown to demonstrate why a stack is necessary |

---

## Approach 1 — Stack (Recommended ✅)

### Intuition
Push every opening bracket. On a closing bracket, check if the stack top is the matching opener. If not (or if the stack is empty), the string is invalid. At the end, the stack must be empty.

The stack enforces the "correct nesting order" requirement. A counter cannot: `"([)]"` has balanced counts but is invalid because `]` tries to close `[` before `)` closes `(`.

### Algorithm
1. `match = {')':'(', ']':'[', '}':'{'}`.
2. For each char `ch`:
   - If `ch` is `(`, `[`, or `{`: push onto stack.
   - Else (closing): if stack empty or top ≠ `match[ch]` → return false; pop.
3. Return `len(stack) == 0`.

### Complexity
- **Time:** O(n) — one pass.
- **Space:** O(n) — at most n/2 openers on the stack.

### Code
```go
func stackApproach(s string) bool {
    stack := make([]byte, 0, len(s)/2)
    match := map[byte]byte{')': '(', ']': '[', '}': '{'}
    for i := 0; i < len(s); i++ {
        ch := s[i]
        if ch == '(' || ch == '[' || ch == '{' {
            stack = append(stack, ch)
        } else {
            if len(stack) == 0 || stack[len(stack)-1] != match[ch] { return false }
            stack = stack[:len(stack)-1]
        }
    }
    return len(stack) == 0
}
```

### Dry Run — `s = "{[]}"`
```
ch='{': push '{'. stack=[{]
ch='[': push '['. stack=[{,[]
ch=']': match[']']='[', top='[' ✓ → pop. stack=[{]
ch='}': match['}]='{', top='{' ✓ → pop. stack=[]
len(stack)==0 → true ✓
```

### Dry Run — `s = "([)]"`
```
ch='(': push. stack=[(]
ch='[': push. stack=[(,[]
ch=')': match[')']]='(', top='[' ✗ → return false ✓
```

---

## Approach 2 — Counter (Educational ⚠️)

### Intuition
For **single bracket type** only (e.g. only `()`): increment on `(`, decrement on `)`. Valid if count never goes negative and ends at 0. **This fails for mixed brackets** — `"([)]"` would return true (2 openers, 2 closers) but is actually invalid.

Included here to demonstrate concretely why a stack is required for multiple bracket types.

---

## Key Takeaways

- **Stack is mandatory for mixed brackets** — the counter approach shows why: a count can be balanced while nesting order is wrong. The stack preserves the exact nesting order because LIFO matches FILO bracket structure.
- **Three failure modes:**
  1. Closer found but stack is empty (e.g. `"]"`).
  2. Closer doesn't match the top (e.g. `"(]"`).
  3. Stack not empty at end — unclosed openers (e.g. `"(("`).
- **Dummy-match map trick** — using `map[byte]byte{')':'(', ']':'[', '}':'{'}` as a lookup is cleaner than a long `switch` statement. A single condition `stack[top] != match[closer]` covers all three bracket types.
- **Pre-allocate the stack** — `make([]byte, 0, len(s)/2)` avoids repeated allocations; at most half the characters can be openers.

---

## Related Problems

- LeetCode #22 — Generate Parentheses (backtracking to generate all valid strings)
- LeetCode #32 — Longest Valid Parentheses (longest valid substring — DP or stack)
- LeetCode #394 — Decode String (stack with nested bracket processing)
- LeetCode #856 — Score of Parentheses (evaluate nested bracket expressions)
