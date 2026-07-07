# 0069 — Sqrt(x)

> LeetCode #69 · Difficulty: Easy
> **Categories:** Math, Binary Search

---

## Problem Statement

Given a non-negative integer `x`, return the **square root of `x` rounded down to the nearest integer**. The returned integer should be non-negative as well.

You **must not** use any built-in exponent function or operator, such as `pow(x, 0.5)` in C++ or `x ** 0.5` in Python.

**Example 1**
```
Input:  x = 4
Output: 2
```

**Example 2**
```
Input:  x = 8
Output: 2
Explanation: The square root of 8 is 2.828..., and since we round it down, 2 is returned.
```

**Constraints**
- `0 <= x <= 2³¹ - 1`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★☆ High      | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★★☆ High      | 2023          |
| Microsoft | ★★★☆☆ Medium    | 2023          |
| Meta      | ★★★☆☆ Medium    | 2023          |
| Apple     | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Binary Search on Answer Space** — search for the largest k in [0, x] such that k² ≤ x.
- **Newton-Raphson Method** — iterative root-finding: converges quadratically to √x.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Linear Scan | O(√x) | O(1) | Reference; too slow for large x |
| 2 | Binary Search ✅ | O(log x) | O(1) | Standard interview answer |
| 3 | Newton's Method | O(log x) | O(1) | Elegant; quadratic convergence |

---

## Approach 1 — Linear Scan

### Intuition
Increment `k` until `k² > x`, then return `k-1`.

### Complexity
- **Time:** O(√x) — up to ~46,341 iterations for x = 2³¹-1.
- **Space:** O(1).

### Code
```go
// linearScan solves Sqrt(x) by trying each integer from 0 upward.
//
// Time:  O(√x)
// Space: O(1)
func linearScan(x int) int {
	if x == 0 {
		return 0
	}
	k := 1
	for k*k <= x {
		k++
	}
	return k - 1 // last k where k²<=x
}
```

### Dry Run — `x = 8`

| k | k*k | k*k <= 8? | action |
|---|-----|-----------|--------|
| 1 | 1 | yes | k++ → 2 |
| 2 | 4 | yes | k++ → 3 |
| 3 | 9 | no | stop |

Loop exits with `k = 3`; return `k - 1 = 2` ✓

---

## Approach 2 — Binary Search (Recommended ✅)

### Intuition
We want the largest `k` such that `k² ≤ x`. Binary search on `[0, x]`:
- If `mid² ≤ x`: `mid` is a candidate, try larger (`lo = mid+1`), save `ans = mid`.
- Else: too large, go smaller (`hi = mid-1`).

### Algorithm
```
lo=0, hi=x, ans=0
while lo <= hi:
  mid = (lo+hi)/2
  if mid*mid <= x: ans=mid; lo=mid+1
  else: hi=mid-1
return ans
```

### Complexity
- **Time:** O(log x) — log₂(2³¹) ≈ 31 iterations.
- **Space:** O(1).

### Code
```go
func binarySearch(x int) int {
    lo, hi, ans := 0, x, 0
    for lo <= hi {
        mid := lo + (hi-lo)/2
        if mid*mid <= x { ans = mid; lo = mid+1 } else { hi = mid-1 }
    }
    return ans
}
```

### Dry Run — `x = 8`
```
lo=0, hi=8: mid=4. 4*4=16>8 → hi=3.
lo=0, hi=3: mid=1. 1*1=1≤8 → ans=1, lo=2.
lo=2, hi=3: mid=2. 2*2=4≤8 → ans=2, lo=3.
lo=3, hi=3: mid=3. 3*3=9>8 → hi=2.
lo=3 > hi=2 → stop. return ans=2 ✓
```

---

## Approach 3 — Newton's Method

### Intuition
Newton-Raphson for `f(r) = r² - x = 0`:
```
r_{new} = r - f(r)/f'(r) = r - (r²-x)/(2r) = (r + x/r) / 2
```
Starting from `r = x`, this halves the error at each step (quadratic convergence).

In integer arithmetic: `r = (r + x/r) / 2`. Stop when `r² ≤ x`.

### Convergence
For `x = 8`: `r = 8 → 4 → 3 → 2 → 2` (converges in 4 steps, vs. 31 for binary search in the worst case).

### Code
```go
func newtonMethod(x int) int {
    r := x
    for r*r > x { r = (r + x/r) / 2 }
    return r
}
```

### Complexity
- **Time:** O(log x) — quadratic convergence in practice is very fast.
- **Space:** O(1).

### Dry Run — `x = 8`

| step | r | r*r | r*r > 8? | update r = (r + x/r)/2 |
|------|---|-----|----------|-------------------------|
| init | 8 | 64 | yes | (8 + 8/8)/2 = (8+1)/2 = 4 |
| 1 | 4 | 16 | yes | (4 + 8/4)/2 = (4+2)/2 = 3 |
| 2 | 3 | 9 | yes | (3 + 8/3)/2 = (3+2)/2 = 2 |
| 3 | 2 | 4 | no | stop |

Loop exits with `r = 2`; return `2` ✓

---

## Key Takeaways

- **`hi = x` is valid** — for x=1, binary search must be able to return 1; starting hi=x handles this. Could optimize to `hi = x/2 + 1` for x > 1.
- **Newton's method converges faster in practice** — but is less obvious to derive in an interview. Binary search is more universally recognized.
- **Integer overflow guard** — for large x, `mid*mid` can overflow int32. In Go, `int` is 64-bit on 64-bit systems so this is fine. In C++/Java, cast to `long`.
- **This is "binary search on the answer"** — the same pattern applies to #278 (First Bad Version), #374 (Guess Number), #1011 (Ship Packages in D Days), etc.

---

## Related Problems

- LeetCode #367 — Valid Perfect Square (check if √x is an integer)
- LeetCode #374 — Guess Number Higher or Lower (binary search on answer space)
- LeetCode #50 — Pow(x, n) (fast exponentiation; inverse operation)
