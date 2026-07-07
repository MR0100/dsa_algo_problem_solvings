# 0050 — Pow(x, n)

> LeetCode #50 · Difficulty: Medium
> **Categories:** Math, Recursion

---

## Problem Statement

Implement `pow(x, n)`, which calculates `x` raised to the power `n` (i.e., `xⁿ`).

**Example 1**
```
Input:  x = 2.00000, n = 10
Output: 1024.00000
```

**Example 2**
```
Input:  x = 2.10000, n = 3
Output: 9.26100
```

**Example 3**
```
Input:  x = 2.00000, n = -2
Output: 0.25000
Explanation: 2⁻² = 1/4 = 0.25
```

**Constraints**
- `-100.0 < x < 100.0`
- `-2³¹ <= n <= 2³¹-1`
- `n` is an integer.
- Either `x != 0` or `n > 0`.
- `-10⁴ <= xⁿ <= 10⁴`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★★ Very High | 2024          |
| Meta      | ★★★★★ Very High | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★★☆ High      | 2023          |
| Apple     | ★★★☆☆ Medium    | 2023          |
| Adobe     | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Exponentiation by Squaring** — `x^n = (x^2)^(n/2)`. Halving the exponent at each step gives O(log n) multiplications.
- **Bit Manipulation** — `n % 2` detects odd exponents; `n >>= 1` halves n; both are O(1).

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (iterative multiply) | O(\|n\|) | O(1) | TLE for large n |
| 2 | Fast Power Iterative ✅ | O(log \|n\|) | O(1) | Optimal; no stack |
| 3 | Fast Power Recursive | O(log \|n\|) | O(log \|n\|) | Elegant; slight stack overhead |

---

## Approach 1 — Brute Force (Repeated Multiplication)

### Intuition
Multiply `x` by itself `|n|` times. If `n < 0`, compute `(1/x)^|n|`.

### Complexity
- **Time:** O(|n|) — up to 2³¹ multiplications. **TLE.**
- **Space:** O(1).

---

## Approach 2 — Fast Power Iterative (Recommended ✅)

### Intuition
Exponentiation by squaring uses the identity:

```
x^n = (x²)^(n/2)    if n is even
x^n = x × (x²)^(n/2) if n is odd
```

Iteratively: while `n > 0`, if the current bit of n is 1, multiply `result` by `x`; then square `x` and right-shift `n`. Each iteration processes one bit of n → O(log n) total.

### Algorithm
```
if n < 0: x = 1/x; n = -n
result = 1.0
while n > 0:
  if n % 2 == 1: result *= x
  x *= x
  n >>= 1
return result
```

### Complexity
- **Time:** O(log |n|) — log₂(|n|) iterations.
- **Space:** O(1).

### Code
```go
func fastPow(x float64, n int) float64 {
    if n < 0 { x = 1/x; n = -n }
    result := 1.0
    for n > 0 {
        if n%2 == 1 { result *= x }
        x *= x
        n >>= 1
    }
    return result
}
```

### Dry Run — `x = 2.0`, `n = 10`
```
n=10 (binary 1010): x=2, result=1
  n=10: n%2=0 → skip. x=4, n=5
  n=5:  n%2=1 → result=4. x=16, n=2
  n=2:  n%2=0 → skip. x=256, n=1
  n=1:  n%2=1 → result=4×256=1024. x=65536, n=0
Result: 1024.0 ✓ (read the set bits: 2^1 + 2^3 = 2+8 = 10; 4¹ × 16³ = 4 × 256 = 1024)
```

---

## Approach 3 — Fast Power Recursive

### Intuition
Base case: `n==0` → 1. Recursive case: compute `half = pow(x, n/2)` and return `half*half` (n even) or `x * half*half` (n odd).

### Code
```go
func fastPowRecursive(x float64, n int) float64 {
    if n == 0 { return 1 }
    if n < 0  { return fastPowRecursive(1/x, -n) }
    half := fastPowRecursive(x, n/2)
    if n%2 == 0 { return half*half }
    return x * half*half
}
```

### Complexity
- **Time:** O(log n).
- **Space:** O(log n) — call stack depth.

---

## Key Takeaways

- **Negative n** — compute `(1/x)^|n|`. Handle before entering the loop/recursion.
- **Absorb odd factor into `result`** — when `n` is odd, `result *= x` "absorbs" one extra factor; then we halve n (effectively treating it as even). The result accumulates the product of x at each set bit of n.
- **This technique is everywhere** — matrix exponentiation (compute Fibonacci in O(log n)), modular exponentiation (RSA), fast polynomial evaluation.
- **n = MinInt32** — `-n` of `MinInt32` overflows int32. In Go, int is 64-bit on 64-bit systems so `-MinInt32` is fine; in languages with fixed 32-bit int (Java/C), cast to long first.

---

## Related Problems

- LeetCode #29 — Divide Two Integers (bit shifting for division)
- LeetCode #372 — Super Pow (modular exponentiation)
- LeetCode #509 — Fibonacci Number (matrix exponentiation uses this)
