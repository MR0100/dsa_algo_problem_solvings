# 0275 — H-Index II

> LeetCode #275 · Difficulty: Medium
> **Categories:** Array, Binary Search

---

## Problem Statement

Given an array of integers `citations` where `citations[i]` is the number of citations a researcher received for their `i`th paper and `citations` is sorted in **ascending order**, return *the researcher's h-index*.

According to the [definition of h-index on Wikipedia](https://en.wikipedia.org/wiki/H-index): The h-index is defined as the maximum value of `h` such that the given researcher has published at least `h` papers that have each been cited at least `h` times.

You must write an algorithm that runs in logarithmic time.

**Example 1:**

```
Input: citations = [0,1,3,5,6]
Output: 3
Explanation: [0,1,3,5,6] means the researcher has 5 papers in total and each of them had received 0, 1, 3, 5, 6 citations respectively.
Since the researcher has 3 papers with at least 3 citations each and the remaining two with no more than 3 citations each, their h-index is 3.
```

**Example 2:**

```
Input: citations = [1,2,100]
Output: 2
```

**Constraints:**

- `n == citations.length`
- `1 <= n <= 10^5`
- `0 <= citations[i] <= 1000`
- `citations` is sorted in ascending order.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Facebook   | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Binary Search** — the input is pre-sorted and the qualifying predicate is monotone, so the h-index is found in O(log n) → see [`/dsa/binary_search.md`](/dsa/binary_search.md)
- **Sorting (given)** — the array arrives sorted ascending, which is precisely what makes binary search applicable → see [`/dsa/sorting.md`](/dsa/sorting.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Linear Scan | O(n) | O(1) | Baseline; simplest, but violates the log-time requirement |
| 2 | Binary Search (Optimal) | O(log n) | O(1) | The required solution — exploits the sorted input |

---

## Approach 1 — Linear Scan

### Intuition

The array is sorted ascending, so at index `i` there are `n - i` papers (from `i` to the end) that each have at least `citations[i]` citations. The h-index is the largest value of `n - i` for which `citations[i] >= n - i`. Scanning left to right, the *first* `i` where `citations[i] >= n - i` gives `h = n - i` — that suffix all qualifies, and being the earliest qualifying index it yields the biggest qualifying suffix.

### Algorithm

1. For `i` from `0` to `n-1`: if `citations[i] >= n - i`, return `n - i`.
2. If none qualify, return `0`.

### Complexity

- **Time:** O(n) — a single pass; this does **not** meet the log-time requirement, hence it is only a baseline.
- **Space:** O(1).

### Code

```go
func linearScan(citations []int) int {
	n := len(citations)
	for i := 0; i < n; i++ {
		// n - i papers (from i to the end) each have >= citations[i] citations.
		if citations[i] >= n-i {
			return n - i // largest qualifying suffix length
		}
	}
	return 0
}
```

### Dry Run

Example 1: `citations = [0,1,3,5,6]`, `n = 5`.

| i | citations[i] | n - i | citations[i] >= n - i? | action |
|---|--------------|-------|------------------------|--------|
| 0 | 0 | 5 | 0 >= 5? no | continue |
| 1 | 1 | 4 | 1 >= 4? no | continue |
| 2 | 3 | 3 | 3 >= 3? **yes** | return n - i = 3 |

Result: `3` ✔

---

## Approach 2 — Binary Search (Optimal)

### Intuition

Define the predicate `f(i) = (citations[i] >= n - i)`. As `i` increases, `citations[i]` is non-decreasing (sorted) while `n - i` decreases — so once `f(i)` becomes true it stays true. `f` is monotone, which means we can **binary-search** the first index where it holds. Every paper from that index onward qualifies, giving `h = n - i`.

### Algorithm

1. `lo = 0`, `hi = n`. While `lo < hi`: `mid = (lo + hi) / 2`.
2. If `citations[mid] >= n - mid`, set `hi = mid` (first-true is at or left of `mid`); else `lo = mid + 1`.
3. Return `n - lo`.

### Complexity

- **Time:** O(log n) — the search range halves each iteration.
- **Space:** O(1).

### Code

```go
func binarySearch(citations []int) int {
	n := len(citations)
	lo, hi := 0, n // search for the first index satisfying the predicate; hi=n means "none"
	for lo < hi {
		mid := (lo + hi) / 2
		// Papers mid..n-1 (that's n-mid of them) each have >= citations[mid] cites.
		if citations[mid] >= n-mid {
			hi = mid // predicate holds → answer index is here or to the left
		} else {
			lo = mid + 1 // predicate fails → move right
		}
	}
	// lo is the first index where citations[lo] >= n - lo; if none, lo == n → 0.
	return n - lo
}
```

### Dry Run

Example 1: `citations = [0,1,3,5,6]`, `n = 5`.

| lo | hi | mid | citations[mid] | n - mid | citations[mid] >= n-mid? | update |
|----|----|-----|----------------|---------|--------------------------|--------|
| 0 | 5 | 2 | 3 | 3 | 3 >= 3 yes | hi = 2 |
| 0 | 2 | 1 | 1 | 4 | 1 >= 4 no | lo = 2 |
| 2 | 2 | — | — | — | loop ends | — |

Return `n - lo = 5 - 2 = 3` ✔

---

## Key Takeaways

- **Sorted input → binary search.** The explicit log-time requirement is a strong hint: exploit the sort with a monotone predicate.
- The reframing `citations[i] >= n - i` turns "count papers with ≥ h citations" into a per-index condition, because in an ascending array the suffix length `n - i` *is* the count of papers with at least `citations[i]` citations.
- **Find-first-true binary search:** when the predicate flips from false to true exactly once, the standard `lo/hi` with `hi = mid` on true and `lo = mid + 1` on false converges to that boundary.
- This is the sorted-array specialization of #274 — same h-index definition, but the pre-sort drops the cost from O(n log n)/O(n) to O(log n).

---

## Related Problems

- LeetCode #274 — H-Index (unsorted; sort or counting buckets)
- LeetCode #35 — Search Insert Position (find-first-index binary search)
- LeetCode #278 — First Bad Version (monotone predicate, binary search the boundary)
- LeetCode #33 — Search in Rotated Sorted Array (binary search on structured input)
