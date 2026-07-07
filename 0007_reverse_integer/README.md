# 0007 — Reverse Integer

> LeetCode #7 · Difficulty: Medium
> **Categories:** Math, Integer Overflow

---

## Problem Statement

Given a signed 32-bit integer `x`, return `x` with its digits reversed. If reversing `x` causes the value to go outside the signed 32-bit integer range `[-2³¹, 2³¹ - 1]`, return `0`.

**Assume the environment does not allow you to store 64-bit integers** (signed or unsigned).

**Example 1**
```
Input:  x = 123
Output: 321
```

**Example 2**
```
Input:  x = -123
Output: -321
```

**Example 3**
```
Input:  x = 120
Output: 21
```

**Constraints**
- `-2³¹ <= x <= 2³¹ - 1`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |
| Adobe     | ★★★☆☆ Medium    | 2023          |
| Apple     | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Math / Digit Manipulation** — pop the last digit with `% 10` and push it onto the result with `* 10 + digit`. This is the core "digit reversal" pattern used throughout integer problems.
- **Overflow Detection** — checking before each push whether `result * 10 + digit` would exceed `INT32_MAX` or `INT32_MIN`.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | String Conversion | O(log x) | O(log x) | Simple; fine if strings are allowed |
| 2 | Math Pop/Push ✅ | O(log x) | O(1) | Correct no-64-bit approach; interview standard |

---

## Approach 1 — String Conversion

### Intuition
Convert `x` to a string (handling sign separately), reverse the characters, then parse back. Use `strconv.ParseInt` with 64-bit precision to detect overflow.

### Algorithm
1. Record sign, take `abs(x)`.
2. `s = strconv.Itoa(abs(x))`.
3. Reverse `s`.
4. `val, err = strconv.ParseInt(reversed, 10, 64)`.
5. If `err != nil` or `val > INT32_MAX` → return 0.
6. Return `sign * int(val)`.

### Complexity
- **Time:** O(log x) — number of digits.
- **Space:** O(log x) — the string.

### Note
This uses a 64-bit integer internally, which violates the problem constraint ("do not use 64-bit integers"). It is included for comparison only; the math approach is the correct interview answer.

### Code
```go
func stringConversion(x int) int {
    if x == 0 {
        return 0
    }

    sign := 1
    if x < 0 {
        sign = -1
        x = -x
    }

    s := strconv.Itoa(x)
    // Reverse the string.
    runes := []byte(s)
    for l, r := 0, len(runes)-1; l < r; l, r = l+1, r-1 {
        runes[l], runes[r] = runes[r], runes[l]
    }

    // Parse back; catch overflow.
    val, err := strconv.ParseInt(string(runes), 10, 64)
    if err != nil || val > math.MaxInt32 {
        return 0
    }
    return sign * int(val)
}
```

### Dry Run — `x = -123`

| Step | Action | State |
|------|--------|-------|
| 1 | `x != 0`, `x < 0` | `sign = -1`, `x = 123` |
| 2 | `strconv.Itoa(123)` | `s = "123"` |
| 3 | reverse bytes (`l=0,r=2` → `l=1,r=1` stop) | `runes = "321"` |
| 4 | `ParseInt("321")`, `321 ≤ MaxInt32` | `val = 321` |
| 5 | `return sign * val` | **-321** |

---

## Approach 2 — Math Pop and Push (Recommended ✅)

### Intuition
Digit reversal in pure arithmetic:
- **Pop:** `digit = x % 10`, `x /= 10`.
- **Push:** `result = result * 10 + digit`.

Overflow check happens **before** the push:
- If `result > INT32_MAX / 10`, then `result * 10` already overflows.
- If `result == INT32_MAX / 10` and `digit > 7`, then `result * 10 + digit > INT32_MAX`.
  (INT32_MAX = 2147483647, last digit = 7)
- Mirror for negative: `result < INT32_MIN / 10` or last digit `< -8`.
  (INT32_MIN = -2147483648, last digit = -8)

This works entirely in 32-bit range — no 64-bit type needed.

### Algorithm
1. While `x != 0`:
   - `digit = x % 10`, `x /= 10`.
   - Overflow check on `result` before push.
   - `result = result * 10 + digit`.
2. Return `result`.

### Complexity
- **Time:** O(log x) — one iteration per digit.
- **Space:** O(1) — no extra allocation.

### Code
```go
func mathPopPush(x int) int {
    result := 0
    for x != 0 {
        digit := x % 10
        x /= 10
        if result > math.MaxInt32/10 || (result == math.MaxInt32/10 && digit > 7) { return 0 }
        if result < math.MinInt32/10 || (result == math.MinInt32/10 && digit < -8) { return 0 }
        result = result*10 + digit
    }
    return result
}
```

### Dry Run — `x = -123`
```
Iteration 1: digit = -123 % 10 = -3, x = -12
  result=0: 0 > -214748364? No. 0 < -214748364? No.
  result = 0*10 + (-3) = -3

Iteration 2: digit = -12 % 10 = -2, x = -1
  result=-3: -3 > -214748364? No. -3 < -214748364? No.
  result = -3*10 + (-2) = -32

Iteration 3: digit = -1 % 10 = -1, x = 0
  result=-32: ok.
  result = -32*10 + (-1) = -321

x == 0 → stop. Return -321 ✓
```

### Dry Run — `x = 1534236469` (overflow)
```
...eventually result = 964323451
Iteration: digit = 1, result = 964323451
  964323451 > 214748364 (MaxInt32/10) → return 0 ✓
```

---

## Key Takeaways

- **`% 10` and `/ 10`** are the fundamental digit pop/push operations. Know them by heart.
- **Overflow check before push, not after** — after would already be undefined behavior in languages without big integers.
- **The last-digit threshold** — for positive overflow, the last digit threshold is 7 (MaxInt32 ends in 7). For negative, it is -8 (MinInt32 ends in -8 as absolute value). In Go, `%` on negative numbers returns a negative remainder, so this check works naturally.
- **Go's `%` for negatives** — in Go (unlike some languages), `x % 10` retains the sign of `x`. So `-123 % 10 = -3`, not `7`. This means the negative-overflow check compares `digit < -8`, which is correct.

---

## Implementation (Go)

See [main.go](main.go).

### Verification
```
x=123              expect=321   → 321
x=-123             expect=-321  → -321
x=120              expect=21    → 21
x=0                expect=0     → 0
x=1534236469       expect=0     → 0
x=-2147483648      expect=0     → 0
```

---

## Related Problems

- LeetCode #8 — String to Integer (atoi) (parsing direction)
- LeetCode #9 — Palindrome Number (check if reversing gives the same number)
- LeetCode #190 — Reverse Bits (similar pop/push but at bit level)
- LeetCode #191 — Number of 1 Bits (bit manipulation on integers)
