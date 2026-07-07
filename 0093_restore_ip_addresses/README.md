# 0093 — Restore IP Addresses

> LeetCode #93 · Difficulty: Medium
> **Categories:** String, Backtracking

---

## Problem Statement

A **valid IP address** consists of exactly four integers separated by single dots. Each integer is between `0` and `255` (inclusive) and cannot have leading zeros.

Given a string `s` containing only digits, return all possible valid IP addresses that can be formed by inserting dots into `s`. You are **not** allowed to reorder or remove any digits in `s`. You may return the valid IP addresses in **any** order.

**Example 1:**
```
Input: s = "25525511135"
Output: ["255.255.11.135","255.255.111.35"]
```

**Example 2:**
```
Input: s = "0000"
Output: ["0.0.0.0"]
```

**Example 3:**
```
Input: s = "101023"
Output: ["1.0.10.23","1.0.102.3","10.1.0.23","10.10.2.3","101.0.2.3"]
```

**Constraints:**
- `1 <= s.length <= 20`

---

## Company Frequency

| Company   | Frequency      | Last Reported |
|-----------|----------------|---------------|
| Amazon    | ★★★☆☆ Medium   | 2024          |
| Google    | ★★★☆☆ Medium   | 2024          |
| Microsoft | ★★☆☆☆ Low      | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Backtracking** — try lengths 1, 2, 3 for each octet; prune early on invalid values. See [`/dsa/backtracking.md`](/dsa/backtracking.md)
- **Constraint-based pruning** — leading zeros, value > 255, insufficient remaining characters.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Backtracking | O(1) | O(1) | Clean and generalizable |
| 2 | Three Nested Loops | O(1) | O(1) | Simple; bounded loops since n≤20 |

Both are O(1) because the string length is bounded (≤20) and there are at most 3³ = 27 possible dot placements (3 choices × 3 positions), so the total work is constant.

---

## Approach 1 — Backtracking

### Intuition
At each step, try taking 1, 2, or 3 characters as the next octet. Prune if:
- The segment has a leading zero (length > 1 and `segment[0] == '0'`).
- The segment's value > 255.
- We've already placed 4 octets but haven't consumed all characters.

When we have exactly 4 octets and have consumed all characters, record the result.

### Algorithm
1. `bt(start, parts, current)`:
   - If `parts == 4 && start == len(s)`: record `current` (remove trailing dot).
   - If `parts == 4` or `start == len(s)`: return.
   - For `length = 1, 2, 3`:
     - If `start+length > len(s)`: break.
     - `segment = s[start:start+length]`.
     - If leading zero (length>1 and segment[0]=='0'): break.
     - If `val > 255`: break.
     - Recurse with `bt(start+length, parts+1, current+segment+".")`.

### Complexity
- **Time:** O(1) — at most 3^4 = 81 leaf paths.
- **Space:** O(1) — recursion depth 4.

### Code
```go
func restoreIpAddresses(s string) []string {
    var result []string
    var bt func(start, parts int, current string)
    bt = func(start, parts int, current string) {
        if parts == 4 && start == len(s) {
            result = append(result, current[:len(current)-1])
            return
        }
        if parts == 4 || start == len(s) { return }
        for length := 1; length <= 3; length++ {
            if start+length > len(s) { break }
            segment := s[start : start+length]
            if length > 1 && segment[0] == '0' { break }
            val, _ := strconv.Atoi(segment)
            if val > 255 { break }
            bt(start+length, parts+1, current+segment+".")
        }
    }
    bt(0, 0, "")
    return result
}
```

### Dry Run (s="0000")

```
bt(0,0,""):
  length=1: seg="0" → bt(1,1,"0.")
    length=1: seg="0" → bt(2,2,"0.0.")
      length=1: seg="0" → bt(3,3,"0.0.0.")
        length=1: seg="0" → bt(4,4,"0.0.0.0.")
          parts==4 && start==4 → record "0.0.0.0"
```

Output: `["0.0.0.0"]` ✓

---

## Approach 2 — Three Nested Loops

### Intuition
Place three dots at positions `i`, `j`, `k` (0-indexed). The four segments are `s[:i]`, `s[i:j]`, `s[j:k]`, `s[k:]`. Check each segment with an `isValid` helper.

The loops are bounded: each dot adds at most 3 characters, so `i ∈ [1,3]`, `j ∈ [i+1,i+3]`, `k ∈ [j+1,j+3]`.

### Complexity
- **Time:** O(1) — at most 3×3×3 = 27 iterations.
- **Space:** O(1)

### Code
```go
func restoreIpAddressesIter(s string) []string {
    var result []string
    n := len(s)
    isValid := func(seg string) bool {
        if len(seg) == 0 || len(seg) > 3 { return false }
        if len(seg) > 1 && seg[0] == '0' { return false }
        v, _ := strconv.Atoi(seg); return v <= 255
    }
    for i := 1; i <= 3 && i < n; i++ {
        for j := i+1; j <= i+3 && j < n; j++ {
            for k := j+1; k <= j+3 && k < n; k++ {
                a, b, c, d := s[:i], s[i:j], s[j:k], s[k:]
                if isValid(a) && isValid(b) && isValid(c) && isValid(d) {
                    result = append(result, a+"."+b+"."+c+"."+d)
                }
            }
        }
    }
    return result
}
```

### Dry Run (s="25525511135")

Some iterations: i=3,j=6,k=8 → "255"."255"."11"."135" → valid ✓; i=3,j=6,k=9 → "255"."255"."111"."35" → valid ✓; all others invalid.

---

## Key Takeaways
- Pruning on leading zeros: `length > 1 && segment[0] == '0'` — break (not continue!) since longer segments from this start will also have a leading zero.
- Pruning on value > 255: similarly break since extending the segment only increases its value.
- The problem has O(1) complexity (bounded input length, bounded branching factor).

---

## Related Problems
- LeetCode #78 — Subsets (similar "try length 1,2,3" partitioning)
- LeetCode #131 — Palindrome Partitioning (partition string into valid parts)
- LeetCode #751 — IP to CIDR (IP address manipulation)
