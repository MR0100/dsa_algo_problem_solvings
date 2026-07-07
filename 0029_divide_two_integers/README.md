# 0029 — Divide Two Integers

> LeetCode #29 · Difficulty: Medium
> **Categories:** Math, Bit Manipulation

---

## Problem Statement

Given two integers `dividend` and `divisor`, divide two integers **without** using multiplication, division, or mod operator.

The integer division should truncate toward zero, which means losing its fractional part.

Return the quotient after dividing `dividend` by `divisor`.

**Note:** Assume we are dealing with an environment that could only store integers within the **32-bit** signed integer range: `[−2³¹, 2³¹ − 1]`. For this problem, if the quotient is strictly greater than `2³¹ − 1`, then return `2³¹ − 1`, and if the quotient is strictly less than `−2³¹`, then return `−2³¹`.

**Example 1**
```
Input:  dividend = 10, divisor = 3
Output: 3
```

**Example 2**
```
Input:  dividend = 7, divisor = -2
Output: -3
```

**Constraints**
- `-2³¹ <= dividend, divisor <= 2³¹ - 1`
- `divisor != 0`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★☆ High      | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |
| Meta      | ★★★☆☆ Medium    | 2023          |
| Apple     | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Bit Manipulation** — left-shift (`<< 1`) doubles a value without multiplication, enabling an exponential speedup over repeated subtraction.
- **Overflow Handling** — `MinInt32 / -1` is the only case that overflows; handle it explicitly. Use `int64` internally to safely negate `MinInt32`.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Repeated Subtraction | O(dividend/divisor) | O(1) | TLE for large inputs; shows basic concept |
| 2 | Bit Shifting ✅ | O(log²n) | O(1) | Interview-optimal; no mul/div/mod |

---

## Approach 1 — Repeated Subtraction (Brute Force)

### Intuition
Division is "how many times does divisor fit into dividend?" Simply subtract divisor repeatedly and count.

### Algorithm
1. Handle overflow: `MinInt32 / -1` → return `MaxInt32`.
2. Determine sign; work with `int64` absolute values.
3. While `a >= b`: `a -= b; count++`.
4. Apply sign; return `int(±count)`.

### Complexity
- **Time:** O(|dividend/divisor|) — up to 2³¹ iterations in the worst case (`MinInt32 / 1`). **TLE in practice.**
- **Space:** O(1).

---

## Approach 2 — Bit Shifting (Recommended ✅)

### Intuition
Instead of subtracting `b` one copy at a time, subtract the largest `2^k * b` that still fits into `a`. This is like the binary representation of the quotient:
```
10 / 3:
  3*1=3  ≤ 10 → try 3*2=6 ≤ 10 → try 3*4=12 > 10: stop at k=1
  subtract 6 (= 3*2²¹); count += 2; a = 4
  3*1=3 ≤ 4 → try 3*2=6 > 4: stop at k=0
  subtract 3; count += 1; a = 1 < 3 → done
  quotient = 2+1 = 3 ✓
```

Each outer loop iteration at least halves `a` (because we subtract the largest fitting power). So the outer loop runs O(log n) times; the inner doubling also runs O(log n) times → O(log² n) total.

### Algorithm
```
if dividend == MinInt32 and divisor == -1: return MaxInt32
negative = sign(dividend) != sign(divisor)
a, b = |dividend|, |divisor|  (as int64)
quotient = 0
while a >= b:
  temp = b; multiple = 1
  while a >= temp*2:
    temp <<= 1; multiple <<= 1
  a -= temp
  quotient += multiple
return ±quotient
```

### Complexity
- **Time:** O(log² n) — outer loop O(log n); inner loop O(log n).
- **Space:** O(1).

### Code
```go
func bitShift(dividend, divisor int) int {
    if dividend == math.MinInt32 && divisor == -1 { return math.MaxInt32 }
    negative := (dividend < 0) != (divisor < 0)
    a, b := int64(dividend), int64(divisor)
    if a < 0 { a = -a }
    if b < 0 { b = -b }
    quotient := int64(0)
    for a >= b {
        temp, multiple := b, int64(1)
        for a >= (temp << 1) { temp <<= 1; multiple <<= 1 }
        a -= temp
        quotient += multiple
    }
    if negative { return int(-quotient) }
    return int(quotient)
}
```

### Dry Run — `dividend = 10, divisor = 3`
```
a=10, b=3, negative=false

Outer iteration 1:
  temp=3, multiple=1
  a=10 >= 3*2=6? yes → temp=6, multiple=2
  a=10 >= 6*2=12? no → stop
  a=10-6=4; quotient=2

Outer iteration 2:
  temp=3, multiple=1
  a=4 >= 3*2=6? no → stop
  a=4-3=1; quotient=3

a=1 < b=3 → exit loop
return 3 ✓
```

---

## Key Takeaways

- **The only overflow** — `MinInt32 / -1 = 2147483648 > MaxInt32`. Handle before anything else.
- **Use `int64` for negation** — `-MinInt32` overflows `int32`; negating in `int64` is safe.
- **Sign before absolute value** — compute the sign once from the originals, then work with absolute values. Apply the sign at the end.
- **Bit shift doubles without multiplying** — `temp <<= 1` is equivalent to `temp *= 2` but uses only bit operations, satisfying the "no multiplication" constraint.

---

## Related Problems

- LeetCode #50 — Pow(x, n) (same "exponential by squaring" idea)
- LeetCode #69 — Sqrt(x) (find quotient without division, using binary search)
- LeetCode #166 — Fraction to Recurring Decimal (long division simulation)
