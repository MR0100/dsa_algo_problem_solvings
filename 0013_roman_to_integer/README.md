# 0013 — Roman to Integer

> LeetCode #13 · Difficulty: Easy
> **Categories:** Hash Table, Math, String

---

## Problem Statement

Given a Roman numeral, convert it to an integer.

Roman numerals are usually written largest to smallest from left to right. However, when a smaller value precedes a larger value, it is subtracted. The six subtraction cases are:
- `I` before `V` (5) or `X` (10) → 4 or 9
- `X` before `L` (50) or `C` (100) → 40 or 90
- `C` before `D` (500) or `M` (1000) → 400 or 900

**Example 1**
```
Input:  s = "III"
Output: 3
```

**Example 2**
```
Input:  s = "LVIII"
Output: 58
Explanation: L=50, V=5, III=3.
```

**Example 3**
```
Input:  s = "MCMXCIV"
Output: 1994
Explanation: M=1000, CM=900, XC=90, IV=4.
```

**Constraints**
- `1 <= s.length <= 15`
- `s` contains only the characters `('I', 'V', 'X', 'L', 'C', 'D', 'M')`.
- It is guaranteed that `s` is a valid Roman numeral in the range `[1, 3999]`.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★★☆ High      | 2024          |
| Meta      | ★★★☆☆ Medium    | 2023          |
| Apple     | ★★★☆☆ Medium    | 2023          |
| Adobe     | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Table / Value Map** — a `map[byte]int` maps each Roman character to its integer value in O(1).
- **String / Linear Scan** — both approaches make a single left-to-right or right-to-left pass.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Left-to-Right with Lookahead | O(n) | O(1) | Natural reading order; explicit subtractive check |
| 2 | Right-to-Left Running Total ✅ | O(n) | O(1) | Cleaner code; no explicit lookahead needed |

---

## Approach 1 — Left-to-Right with Lookahead

### Intuition
Read left to right. If the current symbol has a smaller value than the next symbol, it must be in a subtractive position — subtract it. Otherwise, add it.

### Algorithm
1. Build `val` map.
2. For each index `i`:
   - If `i+1 < len(s)` and `val[s[i+1]] > val[s[i]]` → subtract `val[s[i]]`.
   - Else → add `val[s[i]]`.

### Complexity
- **Time:** O(n).
- **Space:** O(1) — 7-entry constant map.

### Code
```go
func leftToRightLookahead(s string) int {
    val := map[byte]int{'I':1,'V':5,'X':10,'L':50,'C':100,'D':500,'M':1000}
    result := 0
    for i := 0; i < len(s); i++ {
        cur := val[s[i]]
        if i+1 < len(s) && val[s[i+1]] > cur {
            result -= cur
        } else {
            result += cur
        }
    }
    return result
}
```

### Dry Run — `s = "MCMXCIV"`
```
i=0 'M'=1000: next 'C'=100, 100<1000 → add 1000. total=1000
i=1 'C'=100:  next 'M'=1000, 1000>100 → sub 100.  total=900
i=2 'M'=1000: next 'X'=10,  10<1000  → add 1000. total=1900
i=3 'X'=10:   next 'C'=100, 100>10   → sub 10.   total=1890
i=4 'C'=100:  next 'I'=1,   1<100    → add 100.  total=1990
i=5 'I'=1:    next 'V'=5,   5>1      → sub 1.    total=1989
i=6 'V'=5:    no next                → add 5.    total=1994 ✓
```

---

## Approach 2 — Right-to-Left Running Total (Recommended ✅)

### Intuition
Scan from right to left, maintaining a `prev` value. If the current symbol's value is less than `prev`, it's in a subtractive position → subtract it. Otherwise add it. This eliminates the lookahead entirely.

### Algorithm
1. `result = 0`, `prev = 0`.
2. For `i` from `len(s)-1` down to 0:
   - `cur = val[s[i]]`.
   - If `cur < prev` → `result -= cur`.
   - Else → `result += cur`.
   - `prev = cur`.

### Complexity
- **Time:** O(n).
- **Space:** O(1).

### Code
```go
func rightToLeft(s string) int {
    val := map[byte]int{'I':1,'V':5,'X':10,'L':50,'C':100,'D':500,'M':1000}
    result, prev := 0, 0
    for i := len(s) - 1; i >= 0; i-- {
        cur := val[s[i]]
        if cur < prev { result -= cur } else { result += cur }
        prev = cur
    }
    return result
}
```

### Dry Run — `s = "MCMXCIV"` (right to left)
```
i=6 'V'=5:    prev=0,   5>=0  → add 5.    total=5,    prev=5
i=5 'I'=1:    prev=5,   1<5   → sub 1.    total=4,    prev=1
i=4 'C'=100:  prev=1,   100>=1 → add 100. total=104,  prev=100
i=3 'X'=10:   prev=100, 10<100 → sub 10.  total=94,   prev=10
i=2 'M'=1000: prev=10,  1000>=10 → add 1000. total=1094, prev=1000
i=1 'C'=100:  prev=1000, 100<1000 → sub 100. total=994, prev=100
i=0 'M'=1000: prev=100, 1000>=100 → add 1000. total=1994 ✓
```

---

## Key Takeaways

- **Subtractive rule** — a symbol is subtractive when it appears immediately before a larger-valued symbol. Both approaches detect this; the right-to-left scan does it implicitly through the `cur < prev` check.
- **Right-to-left is cleaner** — no lookahead, no index arithmetic for `i+1`. Just compare current with the previous (already-processed) value.
- **Seven symbols, six subtractive cases** — memorise: `IV`, `IX`, `XL`, `XC`, `CD`, `CM`.
- **This is the exact inverse of LeetCode #12** — the same value map is used in both directions.

---

## Implementation (Go)

See [main.go](main.go).

### Verification
```
"III"       → 3    ✓
"LVIII"     → 58   ✓
"MCMXCIV"   → 1994 ✓
"MMMCMXCIX" → 3999 ✓
```

---

## Related Problems

- LeetCode #12 — Integer to Roman (the reverse direction)
- LeetCode #273 — Integer to English Words (number to string, similar lookup approach)
