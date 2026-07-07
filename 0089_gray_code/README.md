# 0089 — Gray Code

> LeetCode #89 · Difficulty: Medium
> **Categories:** Math, Backtracking, Bit Manipulation

---

## Problem Statement

An **n-bit gray code sequence** is a sequence of `2^n` integers where:
- Every integer is in the **inclusive** range `[0, 2^n - 1]`.
- The first integer is `0`.
- An integer appears **no more than once** in the sequence.
- The binary representation of every pair of **adjacent** integers differs by **exactly one bit**.
- The binary representation of the **first** and **last** integers differs by exactly one bit.

Given an integer `n`, return any valid **n-bit gray code sequence**.

**Example 1:**
```
Input: n = 2
Output: [0,1,3,2]
Explanation:
The binary representation of [0,1,3,2] is [00,01,11,10].
- 00 and 01 differ by one bit
- 01 and 11 differ by one bit
- 11 and 10 differ by one bit
- 10 and 00 differ by one bit (wrap-around)
```

**Example 2:**
```
Input: n = 1
Output: [0,1]
```

**Constraints:**
- `1 <= n <= 16`

---

## Company Frequency

| Company  | Frequency      | Last Reported |
|----------|----------------|---------------|
| Amazon   | ★★★☆☆ Medium   | 2024          |
| Google   | ★★☆☆☆ Low      | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Bit Manipulation** — Gray code formula: `gray(i) = i XOR (i >> 1)`. See [`/dsa/bit_manipulation.md`](/dsa/bit_manipulation.md)
- **Recursion / Mirror Construction** — n-bit Gray code = (n-1)-bit code prepended with 0 + reversed (n-1)-bit code prepended with 1.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Bit Manipulation Formula | O(2^n) | O(1) | Simplest; O(1) per element |
| 2 | Recursive Mirror | O(2^n) | O(2^n) | Illustrates the structural insight |

---

## Approach 1 — Bit Manipulation Formula

### Intuition
The Gray code of integer `i` is `i XOR (i >> 1)`.

**Why this works:** for consecutive integers `i` and `i+1`, only one bit changes in their Gray codes. The XOR of adjacent Gray codes is `(i XOR (i>>1)) XOR ((i+1) XOR ((i+1)>>1))`, which always equals a power of 2 (single bit set). This is a known property of this particular encoding.

### Algorithm
1. For `i = 0` to `2^n - 1`: `result[i] = i ^ (i >> 1)`.

### Complexity
- **Time:** O(2^n)
- **Space:** O(1) extra (excluding output).

### Code
```go
func grayCode(n int) []int {
    size := 1 << n
    result := make([]int, size)
    for i := 0; i < size; i++ {
        result[i] = i ^ (i >> 1)
    }
    return result
}
```

### Dry Run (n=2, size=4)

| i | i>>1 | i XOR (i>>1) | binary |
|---|------|--------------|--------|
| 0 | 0 | 0 | 00 |
| 1 | 0 | 1 | 01 |
| 2 | 1 | 3 | 11 |
| 3 | 1 | 2 | 10 |

Adjacent pairs: 00↔01 (1 bit), 01↔11 (1 bit), 11↔10 (1 bit), 10↔00 (1 bit, wrap) ✓

---

## Approach 2 — Recursive Mirror Construction

### Intuition
Build the n-bit Gray code from the (n-1)-bit Gray code:
1. Take the (n-1)-bit sequence and **prepend 0** to each code (unchanged values).
2. Take the (n-1)-bit sequence in **reverse** and **prepend 1** (add `2^(n-1)` to each).

The boundary at the midpoint: the last element of the first half and the first element of the second half differ by only bit `n-1` (the prepended bit), making them adjacent in the output.

```
n=1: [0, 1]
n=2: [0,1] → [0,1, |mirror| 3,2] = [0,1,3,2]
n=3: [0,1,3,2] → [0,1,3,2, |mirror| 6,7,5,4]
```

### Algorithm
1. Base: `n=0` → `[0]`.
2. `prev = grayCodeMirror(n-1)`.
3. `result = prev + [v + 2^(n-1) for v in reversed(prev)]`.

### Complexity
- **Time:** O(2^n)
- **Space:** O(2^n) — recursion + output.

### Code
```go
func grayCodeMirror(n int) []int {
    if n == 0 { return []int{0} }
    prev := grayCodeMirror(n - 1)
    result := make([]int, 0, 2*len(prev))
    result = append(result, prev...)
    half := 1 << (n - 1)
    for i := len(prev) - 1; i >= 0; i-- {
        result = append(result, prev[i]+half)
    }
    return result
}
```

### Dry Run (n=3)

```
prev = [0,1,3,2] (from n=2)
first half: [0,1,3,2]
second half (reversed + 4): [2+4,3+4,1+4,0+4] = [6,7,5,4]
result: [0,1,3,2,6,7,5,4]
```

Check boundary: 2 (binary `010`) and 6 (binary `110`) differ by 1 bit ✓
Wrap-around: 4 (binary `100`) and 0 (binary `000`) differ by 1 bit ✓

---

## Key Takeaways
- `gray(i) = i ^ (i >> 1)` is the formula — O(1) per element, simplest implementation.
- Mirror construction reveals the recursive structure: each n-bit code is two mirrored (n-1)-bit codes with the MSB appended.
- Gray codes appear in angle encoders, error correction, and Hamiltonian cycles on hypercubes.
- Verify adjacency: consecutive pairs differ by 1 bit; wrap-around also differs by 1 bit.

---

## Related Problems
- LeetCode #1611 — Minimum One Bit Operations to Make Integers Zero (Gray code inverse)
- LeetCode #847 — Shortest Path Visiting All Nodes (Hamiltonian path; Gray code gives such a path on hypercubes)
