# 0067 — Add Binary

> LeetCode #67 · Difficulty: Easy
> **Categories:** Math, String, Bit Manipulation, Simulation

---

## Problem Statement

Given two binary strings `a` and `b`, return their sum as a binary string.

**Example 1**
```
Input:  a = "11", b = "1"
Output: "100"
```

**Example 2**
```
Input:  a = "1010", b = "1011"
Output: "10101"
```

**Constraints**
- `1 <= a.length, b.length <= 10⁴`
- `a` and `b` consist only of `'0'` or `'1'` characters.
- Each string does not contain leading zeros except for the zero itself.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★☆ High      | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Meta      | ★★★☆☆ Medium    | 2023          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |
| Apple     | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Binary Arithmetic / Carry Propagation** — bit-by-bit addition from LSB to MSB with carry.
- **String Building** — append digits LSB-first, then reverse.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Bit-by-Bit Addition ✅ | O(max(m,n)) | O(max(m,n)) | The standard and only approach needed |

---

## Approach 1 — Bit-by-Bit Addition (Recommended ✅)

### Intuition
Walk both strings from right to left (LSB to MSB). At each position, add two bits plus the carry. The result bit is `sum % 2`; new carry is `sum / 2`. Build the result left-to-right (LSB-first), then reverse at the end.

### Algorithm
```
i = len(a)-1; j = len(b)-1; carry = 0
while i >= 0 or j >= 0 or carry > 0:
  sum = carry
  if i >= 0: sum += a[i]-'0'; i--
  if j >= 0: sum += b[j]-'0'; j--
  append sum%2 to result
  carry = sum/2
reverse result
```

### Complexity
- **Time:** O(max(m, n)) — iterate through the longer string.
- **Space:** O(max(m, n)) — output string.

### Code
```go
func addBinary(a string, b string) string {
    i, j, carry := len(a)-1, len(b)-1, 0
    var sb strings.Builder
    for i >= 0 || j >= 0 || carry > 0 {
        sum := carry
        if i >= 0 { sum += int(a[i]-'0'); i-- }
        if j >= 0 { sum += int(b[j]-'0'); j-- }
        sb.WriteByte(byte('0' + sum%2))
        carry = sum / 2
    }
    res := []byte(sb.String())
    for l, r := 0, len(res)-1; l < r; l, r = l+1, r-1 { res[l], res[r] = res[r], res[l] }
    return string(res)
}
```

### Dry Run — `a = "11"`, `b = "1"`
```
i=1, j=0, carry=0:
  sum=0 + a[1]='1'(1) + b[0]='1'(1) = 2. write 0. carry=1. i=0,j=-1.
i=0, j=-1, carry=1:
  sum=1 + a[0]='1'(1) = 2. write 0. carry=1. i=-1.
i=-1, j=-1, carry=1:
  sum=1. write 1. carry=0.

Built LSB-first: "001" → reversed: "100" ✓
```

---

## Key Takeaways

- **Build LSB-first then reverse** — simpler than prepending to the front (which would be O(n) per prepend).
- **The loop handles unequal lengths naturally** — once one string is exhausted, we stop adding its bits (just use 0) but continue for the carry.
- **Carries can cascade** — `"1" + "1" + carry=1` = 3 = 11₂, so carry can reach 1 again. The loop ends only when both pointers are exhausted AND carry is 0.
- **Generalises to any base** — replace base-2 literals with base B to add numbers in any base.

---

## Related Problems

- LeetCode #66 — Plus One (base-10 carry propagation)
- LeetCode #43 — Multiply Strings (grade-school multiplication with carry)
- LeetCode #371 — Sum of Two Integers (bit manipulation without `+`)
