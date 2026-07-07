# 0008 — String to Integer (atoi)

> LeetCode #8 · Difficulty: Medium
> **Categories:** String, Two Pointers, Math

---

## Problem Statement

Implement the `myAtoi(string s)` function, which converts a string to a 32-bit signed integer.

The algorithm for `myAtoi(s)` is as follows:
1. **Whitespace**: Ignore any leading whitespace (`' '`).
2. **Signedness**: Determine the sign by checking if the next character is `'-'` or `'+'`, defaulting to positive.
3. **Conversion**: Read in next characters until the next non-digit character or end of input. The rest of the string is ignored.
4. **Rounding**: If the integer is out of the 32-bit range `[−2³¹, 2³¹ − 1]`, clamp to the range boundary.

**Example 1**
```
Input:  s = "42"
Output: 42
```

**Example 2**
```
Input:  s = "   -042"
Output: -42
```

**Example 3**
```
Input:  s = "1337c0d3"
Output: 1337
```

**Example 4**
```
Input:  s = "0-1"
Output: 0
```

**Example 5**
```
Input:  s = "words and 987"
Output: 0
```

**Constraints**
- `0 <= s.length <= 200`
- `s` consists of English letters (lower-case and upper-case), digits (`0-9`), `' '`, `'+'`, `'-'`, and `'.'`.

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
| Uber      | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **String / Linear Scan** — a single left-to-right pass through the string with a simple state machine: whitespace → sign → digits → stop.
- **Overflow / Clamping** — accumulate in `int64` and clamp to `[-2³¹, 2³¹-1]` to satisfy the 32-bit requirement.

---

## Approaches Overview

| # | Approach | Time | Space | Notes |
|---|----------|------|-------|-------|
| 1 | stdlib ParseInt (approx) | O(n) | O(1) | Approximation only; fails partial-parse cases like "1337c0d3" |
| 2 | Manual Linear Scan ✅ | O(n) | O(1) | Exact implementation; the correct interview answer |

---

## Approach 1 — stdlib ParseInt (Approximation)

### Intuition
Trim whitespace, call `strconv.ParseInt`. Map the range error to INT32 bounds.

### Limitation
`strconv.ParseInt` fails entirely on strings like `"1337c0d3"` (it does not do partial parsing — it stops and errors). The correct `myAtoi` should return `1337`. This approach is included to show the stdlib's behavior and why a custom parser is necessary.

### Complexity
- **Time:** O(n) — linear scan to trim whitespace plus `strconv.ParseInt`'s own linear scan of the remaining string.
- **Space:** O(1) — only the reslice offset and the parsed `int64`; no extra allocation.

### Code
```go
func useStdlib(s string) int {
    // Trim leading whitespace.
    i := 0
    for i < len(s) && s[i] == ' ' {
        i++
    }
    s = s[i:]

    val, err := strconv.ParseInt(s, 10, 64)
    if err != nil {
        // Range error: clamp.
        if numErr, ok := err.(*strconv.NumError); ok && numErr.Err == strconv.ErrRange {
            if len(s) > 0 && s[0] == '-' {
                return math.MinInt32
            }
            return math.MaxInt32
        }
        // Syntax error: partial parse — strconv doesn't do partial parse like atoi.
        // Fall back to 0 for this approximation.
        return 0
    }
    if val > math.MaxInt32 {
        return math.MaxInt32
    }
    if val < math.MinInt32 {
        return math.MinInt32
    }
    return int(val)
}
```

### Dry Run — `s = "   -042"`
```
Trim whitespace: i advances over 3 spaces → s = "-042"
strconv.ParseInt("-042", 10, 64):
  sees '-' → negative, then digits 0,4,2 → val = -42, err = nil
err == nil, so no clamp path
val=-42 is within [MinInt32, MaxInt32]
return int(-42) = -42 ✓
```

| step        | s / val      | note                              |
|-------------|--------------|-----------------------------------|
| trim spaces | `"-042"`     | i skipped 3 leading spaces        |
| ParseInt    | `val=-42`    | whole remaining string is numeric |
| err check   | `err==nil`   | no range/syntax error             |
| clamp check | within range | -42 ∈ [MinInt32, MaxInt32]        |
| return      | `-42`        | int(val)                          |

---

## Approach 2 — Manual Linear Scan (Correct ✅)

### Intuition
A textbook finite state machine with three states:
1. **Whitespace state** — advance `i` while `s[i] == ' '`.
2. **Sign state** — check for `+` / `-`, record sign, advance `i`.
3. **Digit state** — read digits, multiply running result by 10, add digit; stop at first non-digit.

Accumulate into `int64` so overflow can be detected cleanly; clamp before returning.

### Algorithm
```
i = 0
while s[i] == ' ': i++         // skip whitespace
if s[i] == '+' or '-': i++     // consume sign
while s[i] is digit:
    result = result*10 + digit
    clamp to [INT32_MIN, INT32_MAX] if needed
    i++
return sign * result
```

### Complexity
- **Time:** O(n) — at most one full pass.
- **Space:** O(1) — no extra allocation.

### Code
```go
func manualScan(s string) int {
    i, n := 0, len(s)
    for i < n && s[i] == ' ' { i++ }
    if i == n { return 0 }
    sign := 1
    if s[i] == '+' { i++ } else if s[i] == '-' { sign = -1; i++ }
    var result int64
    for i < n && s[i] >= '0' && s[i] <= '9' {
        result = result*10 + int64(s[i]-'0')
        if sign == 1 && result > math.MaxInt32 { return math.MaxInt32 }
        if sign == -1 && -result < math.MinInt32 { return math.MinInt32 }
        i++
    }
    return sign * int(result)
}
```

### Dry Run — `s = "   -042"`
```
i=0,1,2: spaces → i=3
i=3: '-' → sign=-1, i=4
i=4: '0' → result=0
i=5: '4' → result=4
i=6: '2' → result=42
i=7: end of string
return -1 * 42 = -42 ✓
```

### Dry Run — `s = "1337c0d3"`
```
i=0: no whitespace
i=0: no sign
digits: '1'→1, '3'→13, '3'→133, '7'→1337
i=4: 'c' is not a digit → stop
return 1 * 1337 = 1337 ✓
```

---

## Key Takeaways

- **Stop at the first non-digit** — do NOT skip non-digits. The parser terminates as soon as it sees a non-digit character after starting to read digits.
- **Sign comes before digits, whitespace before sign** — no whitespace is allowed between sign and first digit. `"  + 413"` returns 0 because after `+` there is a space.
- **Use int64 internally, clamp to int32 on output** — this is the cleanest pattern for avoiding mid-computation overflow.
- **Early clamp inside the loop** — clamp as soon as the running result exceeds the bound, not at the end. Otherwise you'd need to handle very long strings of digits safely.

---

## Implementation (Go)

See [main.go](main.go).

### Verification
```
s="42"                    expect=42      → Manual: 42
s="   -042"               expect=-42     → Manual: -42
s="1337c0d3"              expect=1337    → Manual: 1337
s="0-1"                   expect=0       → Manual: 0
s="words and 987"         expect=0       → Manual: 0
s="-91283472332"          expect=-2^31   → Manual: -2147483648
s="2147483648"            expect=2^31-1  → Manual: 2147483647
```

---

## Related Problems

- LeetCode #7 — Reverse Integer (overflow handling with same INT32 boundary)
- LeetCode #65 — Valid Number (full number validation, harder)
- LeetCode #273 — Integer to English Words (int → string conversion)
