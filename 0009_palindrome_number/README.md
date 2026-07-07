# 0009 — Palindrome Number

> LeetCode #9 · Difficulty: Easy
> **Categories:** Math, Two Pointers

---

## Problem Statement

Given an integer `x`, return `true` if `x` is a **palindrome**, and `false` otherwise.

**Example 1**
```
Input:  x = 121
Output: true
Explanation: 121 reads as 121 from left to right and from right to left.
```

**Example 2**
```
Input:  x = -121
Output: false
Explanation: From left to right, it reads -121. From right to left, it reads 121-. Therefore it is not a palindrome.
```

**Example 3**
```
Input:  x = 10
Output: false
Explanation: Reads 01 from right to left. Therefore it is not a palindrome.
```

**Constraints**
- `-2³¹ <= x <= 2³¹ - 1`

**Follow-up:** Could you solve it without converting the integer to a string?

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★☆ High      | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Microsoft | ★★★☆☆ Medium    | 2024          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |
| Adobe     | ★★☆☆☆ Low       | 2022          |
| Apple     | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Math / Digit Manipulation** — all non-string approaches use `% 10` and `/ 10` to extract and rebuild digits.
- **Two Pointers (conceptual)** — Approach 3 compares the "front half" (remaining `x`) against the "back half" (accumulated `reversed`), conceptually a two-pointer convergence.

---

## Approaches Overview

| # | Approach | Time | Space | Notes |
|---|----------|------|-------|-------|
| 1 | String Conversion | O(log x) | O(log x) | Simplest; violates "no string" follow-up |
| 2 | Reverse Full Number | O(log x) | O(1) | Correct; risk of intermediate overflow (mitigated by int64) |
| 3 | Reverse Half ✅ | O(log x) | O(1) | Optimal; no overflow risk; cleanest math |

---

## Approach 1 — String Conversion

### Intuition
The simplest definition: a palindrome reads the same forwards and backwards. Convert to string, compare from both ends inward.

### Algorithm
1. If `x < 0` → return false.
2. `s = strconv.Itoa(x)`.
3. Two-pointer check: `l=0, r=len(s)-1`, compare `s[l]` and `s[r]`, advance inward.

### Complexity
- **Time:** O(log x) — O(d) where d = number of digits.
- **Space:** O(log x) — the string.

---

## Approach 2 — Reverse Full Number

### Intuition
Reverse all digits of `x` using integer arithmetic and compare against the original. Negative numbers are never palindromes. Numbers ending in 0 (except 0 itself) can't be palindromes (reversed, they'd have a leading 0).

### Algorithm
1. If `x < 0` or (`x % 10 == 0` and `x != 0`) → return false.
2. `original = x`.
3. `reversed = 0`; while `x > 0`: `reversed = reversed*10 + x%10; x /= 10`.
4. Return `original == reversed`.

### Complexity
- **Time:** O(log x).
- **Space:** O(1).

---

## Approach 3 — Reverse Only Half (Recommended ✅)

### Intuition
We only need to compare the first half of the digits with the reversed second half. Stop when `reversed >= x` — that means we've reversed exactly half (or half + 1 for odd-length numbers).

- **Even-length palindrome** (e.g. 1221): stop when `reversed == x`. `1221 → reverse 1,2 → reversed=12, x=12`. Return `x == reversed`.
- **Odd-length palindrome** (e.g. 12321): stop when `reversed > x`. `12321 → reverse 1,2,3 → reversed=123, x=12`. The middle digit is in `reversed`; ignore it with `reversed/10`. Return `x == reversed/10`.

This avoids reversing the whole number and has no overflow risk since `reversed` stays ≤ the original.

### Algorithm
1. If `x < 0` or (`x % 10 == 0` and `x != 0`) → false.
2. `reversed = 0`.
3. While `x > reversed`: `reversed = reversed*10 + x%10; x /= 10`.
4. Return `x == reversed || x == reversed/10`.

### Complexity
- **Time:** O(log x) — half the digits.
- **Space:** O(1).

### Code
```go
func reverseHalf(x int) bool {
    if x < 0 || (x%10 == 0 && x != 0) { return false }
    reversed := 0
    for x > reversed {
        reversed = reversed*10 + x%10
        x /= 10
    }
    return x == reversed || x == reversed/10
}
```

### Dry Run — `x = 12321`
```
Initial: x=12321, reversed=0

Iteration 1: x > reversed (12321 > 0)
  digit = 12321 % 10 = 1
  reversed = 0*10 + 1 = 1
  x = 12321 / 10 = 1232

Iteration 2: x > reversed (1232 > 1)
  digit = 1232 % 10 = 2
  reversed = 1*10 + 2 = 12
  x = 1232 / 10 = 123

Iteration 3: x > reversed (123 > 12)
  digit = 123 % 10 = 3
  reversed = 12*10 + 3 = 123
  x = 123 / 10 = 12

Check: x(12) > reversed(123)? No → stop.
Return: x(12) == reversed/10(12) → true ✓  (odd length: skip middle digit 3)
```

### Dry Run — `x = 1221`
```
Iteration 1: reversed=1, x=122
Iteration 2: reversed=12, x=12
Check: x(12) > reversed(12)? No → stop.
Return: x(12) == reversed(12) → true ✓  (even length: exact match)
```

### Dry Run — `x = 10`
```
x%10 == 0 and x != 0 → return false ✓
```

---

## Key Takeaways

- **Negative numbers and trailing zeros** — always false. Negative numbers have a leading `-` which is not a digit. A number ending in 0 (other than 0 itself) would need a leading 0 when reversed, which is impossible.
- **Reverse half, not all** — reversing the whole number risks overflow for large inputs. Reversing half is cleaner and naturally terminates when `reversed >= remaining`.
- **The `reversed/10` trick** — for odd-length palindromes, the middle digit ends up in `reversed` after the loop. Dividing by 10 removes it before comparison.
- **`x > reversed` as the loop condition** — this is the key: it stops exactly when we've consumed half the digits (or one more for odd-length).

---

## Implementation (Go)

See [main.go](main.go).

### Verification
```
x=121      expect=true   → true
x=-121     expect=false  → false
x=10       expect=false  → false
x=0        expect=true   → true
x=1221     expect=true   → true
x=12321    expect=true   → true
x=123      expect=false  → false
```

---

## Related Problems

- LeetCode #7 — Reverse Integer (same digit reversal technique)
- LeetCode #125 — Valid Palindrome (palindrome check on a string)
- LeetCode #234 — Palindrome Linked List (palindrome check without converting to array)
- LeetCode #5 — Longest Palindromic Substring (substring palindrome problems)
