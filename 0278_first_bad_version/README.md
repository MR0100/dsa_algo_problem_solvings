# 0278 — First Bad Version

> LeetCode #278 · Difficulty: Easy
> **Categories:** Binary Search, Interactive

---

## Problem Statement

You are a product manager and currently leading a team to develop a new product.
Unfortunately, the latest version of your product fails the quality check. Since
each version is developed based on the previous version, all the versions after
a bad version are also bad.

Suppose you have `n` versions `[1, 2, ..., n]` and you want to find out the first
bad one, which causes all the following ones to be bad.

You are given an API `bool isBadVersion(version)` which returns whether `version`
is bad. Implement a function to find the first bad version. You should minimize
the number of calls to the API.

**Example 1:**

```
Input: n = 5, bad = 4
Output: 4
Explanation:
call isBadVersion(3) -> false
call isBadVersion(5) -> true
call isBadVersion(4) -> true
Then 4 is the first bad version.
```

**Example 2:**

```
Input: n = 1, bad = 1
Output: 1
```

**Constraints:**

- `1 <= bad <= n <= 2^31 - 1`

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Facebook  | ★★★★☆ High       | 2024          |
| Google    | ★★★☆☆ Medium     | 2023          |
| Amazon    | ★★★☆☆ Medium     | 2023          |
| Microsoft | ★★★☆☆ Medium     | 2023          |
| Apple     | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Binary Search on a boolean boundary** — the versions form a sorted sequence
  `[good…good, bad…bad]`; we search for the leftmost `true` → see [`/dsa/binary_search.md`](/dsa/binary_search.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Linear Scan | O(n) | O(1) | Simple, but too many API calls for large `n` |
| 2 | Binary Search (Optimal) | O(log n) | O(1) | Minimizes API calls; the intended solution |

---

## Approach 1 — Linear Scan

### Intuition
Since badness is monotone, the first version that returns `true` is the answer.
Walk from version 1 upward and stop at the first bad one.

### Algorithm
1. For `v = 1..n`: if `isBadVersion(v)`, return `v`.
2. (A bad version is guaranteed to exist.)

### Complexity
- **Time:** O(n) — up to `n` API calls.
- **Space:** O(1).

### Code
```go
func linearScan(n int) int {
	for v := 1; v <= n; v++ { // scan versions in order
		if isBadVersion(v) { // first true = first bad version
			return v
		}
	}
	return -1 // unreachable per problem guarantees
}
```

### Dry Run
`n = 5`, first bad = 4:

| v | isBadVersion(v) | Action |
|---|-----------------|--------|
| 1 | false | continue |
| 2 | false | continue |
| 3 | false | continue |
| 4 | true  | return 4 |

Return **4**.

---

## Approach 2 — Binary Search (Optimal)

### Intuition
The sequence is `[good, good, …, good, bad, bad, …, bad]`. Finding the first bad
version is finding the boundary in a sorted boolean array — the leftmost `true`.
Halve the interval each step: if `mid` is bad, the boundary is at `mid` or to its
left; if `mid` is good, it is strictly to the right.

### Algorithm
1. `lo = 1`, `hi = n`.
2. While `lo < hi`:
   - `mid = lo + (hi-lo)/2` (overflow-safe).
   - if `isBadVersion(mid)`: `hi = mid` (keep `mid` as a candidate).
   - else: `lo = mid + 1`.
3. When `lo == hi`, that is the first bad version.

### Complexity
- **Time:** O(log n) — the interval halves each iteration.
- **Space:** O(1).

### Code
```go
func binarySearch(n int) int {
	lo, hi := 1, n // search space is [1, n]
	for lo < hi {  // stop when the interval collapses to one version
		mid := lo + (hi-lo)/2 // midpoint without integer overflow
		if isBadVersion(mid) {
			hi = mid // mid is bad → first bad is mid or earlier; keep mid
		} else {
			lo = mid + 1 // mid is good → first bad is strictly after mid
		}
	}
	return lo // lo == hi points at the boundary: the first bad version
}
```

### Dry Run
`n = 5`, first bad = 4:

| lo | hi | mid | isBadVersion(mid) | Update |
|----|----|-----|-------------------|--------|
| 1 | 5 | 3 | false | lo = 4 |
| 4 | 5 | 4 | true  | hi = 4 |
| 4 | 4 | — | — | loop ends |

`lo == hi == 4` → return **4**.

---

## Key Takeaways

- "First `true` in a monotone boolean sequence" is a canonical binary-search
  template: move `hi` inward on true, `lo` past `mid` on false.
- Use `mid = lo + (hi-lo)/2` instead of `(lo+hi)/2` — `n` can be up to `2^31 - 1`,
  so the naive sum overflows a 32-bit int.
- Using `hi = mid` (not `mid - 1`) is what makes the search converge on the
  boundary rather than skipping past it.

---

## Related Problems

- LeetCode #35 — Search Insert Position (leftmost boundary search)
- LeetCode #34 — Find First and Last Position of Element in Sorted Array
- LeetCode #704 — Binary Search (baseline template)
- LeetCode #374 — Guess Number Higher or Lower (interactive binary search)
