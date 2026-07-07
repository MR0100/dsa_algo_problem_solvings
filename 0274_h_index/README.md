# 0274 — H-Index

> LeetCode #274 · Difficulty: Medium
> **Categories:** Array, Sorting, Counting Sort

---

## Problem Statement

Given an array of integers `citations` where `citations[i]` is the number of citations a researcher received for their `i`th paper, return *the researcher's h-index*.

According to the [definition of h-index on Wikipedia](https://en.wikipedia.org/wiki/H-index): The h-index is defined as the maximum value of `h` such that the given researcher has published at least `h` papers that have each been cited at least `h` times.

**Example 1:**

```
Input: citations = [3,0,6,1,5]
Output: 3
Explanation: [3,0,6,1,5] means the researcher has 5 papers in total and each of them had received 3, 0, 6, 1, 5 citations respectively.
Since the researcher has 3 papers with at least 3 citations each and the remaining two with no more than 3 citations each, their h-index is 3.
```

**Example 2:**

```
Input: citations = [1,3,1]
Output: 1
```

**Constraints:**

- `n == citations.length`
- `1 <= n <= 5000`
- `0 <= citations[i] <= 1000`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2024          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Facebook   | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Sorting** — sorting citations descending makes the h-index a simple rank scan → see [`/dsa/sorting.md`](/dsa/sorting.md)
- **Counting Sort / Buckets** — because h ≤ n, bucketing by citation count gives an O(n) solution → see [`/dsa/counting_sort.md`](/dsa/counting_sort.md)
- **Binary Search on the Answer** — the "≥ h papers with ≥ h citations" predicate is monotone, so h can be binary-searched → see [`/dsa/binary_search.md`](/dsa/binary_search.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Sort Descending | O(n log n) | O(1) | Cleanest; the natural first solution |
| 2 | Counting Buckets | O(n) | O(n) | Optimal time; exploits h ≤ n |
| 3 | Binary Search on the Answer | O(n log n) | O(1) | Demonstrates monotonicity; good when input is huge/streamed |

---

## Approach 1 — Sort Descending

### Intuition

Sort papers by citations, most-cited first. Walk down the list: after examining `i` papers, all `i` of them have at least `citations[i-1]` citations (because the list is sorted descending). So `h` can be at least `i` as long as the `i`-th paper still has `>= i` citations. The h-index is the largest such `i`; once a paper's citation count drops below its 1-based rank, no larger `h` is possible.

### Algorithm

1. Sort `citations` in descending order.
2. For `i` from `0` to `n-1`: if `citations[i] >= i+1`, set `h = i+1`; else break.
3. Return `h`.

### Complexity

- **Time:** O(n log n) — dominated by the sort.
- **Space:** O(1) — in-place sort (aside from the sort library's internals).

### Code

```go
func sortDescending(citations []int) int {
	// Sort a copy would be cleaner, but sorting in place is standard here.
	sort.Sort(sort.Reverse(sort.IntSlice(citations)))
	h := 0
	for i := 0; i < len(citations); i++ {
		// citations[i] is the (i+1)-th largest; if it still has >= i+1 cites,
		// we have i+1 papers each cited >= i+1 times.
		if citations[i] >= i+1 {
			h = i + 1
		} else {
			break // counts only decrease from here — no larger h possible
		}
	}
	return h
}
```

### Dry Run

Example 1: `citations = [3,0,6,1,5]`.

| Step | action | state |
|------|--------|-------|
| 1 | sort desc | `[6,5,3,1,0]` |
| 2 | i=0: 6 >= 1? yes | h = 1 |
| 3 | i=1: 5 >= 2? yes | h = 2 |
| 4 | i=2: 3 >= 3? yes | h = 3 |
| 5 | i=3: 1 >= 4? no | break |

Result: `h = 3` ✔

---

## Approach 2 — Counting Buckets

### Intuition

The h-index can never exceed `n` (you can't have more qualifying papers than you published). So bucket papers by citation count, lumping everything `>= n` into bucket `n`. Then scan buckets from high to low, accumulating how many papers have *at least* that many citations. The first citation level `h` where the running count `>= h` is the h-index.

### Algorithm

1. Build `buckets[0..n]`; for each citation `c`, increment `buckets[min(c, n)]`.
2. `total = 0`; for `h` from `n` down to `0`: `total += buckets[h]`; if `total >= h`, return `h`.

### Complexity

- **Time:** O(n) — one pass to bucket, one pass over `n+1` buckets.
- **Space:** O(n) — the bucket array.

### Code

```go
func countingBuckets(citations []int) int {
	n := len(citations)
	buckets := make([]int, n+1) // index = citation count, capped at n
	for _, c := range citations {
		if c >= n {
			buckets[n]++ // everything >= n lands in the top bucket
		} else {
			buckets[c]++
		}
	}
	total := 0 // papers with AT LEAST h citations, accumulated high→low
	for h := n; h >= 0; h-- {
		total += buckets[h]
		if total >= h {
			// h papers each have >= h citations → h is achievable, and since we
			// scan from the top, this is the maximum such h.
			return h
		}
	}
	return 0
}
```

### Dry Run

Example 1: `citations = [3,0,6,1,5]`, `n = 5`.

Bucketing (cap at 5): `3→b[3]`, `0→b[0]`, `6→b[5]`, `1→b[1]`, `5→b[5]`.

`buckets = [1, 1, 0, 1, 0, 2]` (indices 0..5).

| h | buckets[h] | total (running) | total >= h? |
|---|------------|-----------------|-------------|
| 5 | 2 | 2 | 2 >= 5? no |
| 4 | 0 | 2 | 2 >= 4? no |
| 3 | 1 | 3 | 3 >= 3? **yes → return 3** |

Result: `h = 3` ✔

---

## Approach 3 — Binary Search on the Answer

### Intuition

The predicate "at least `h` papers have `>= h` citations" is monotone in `h`: if it holds for some `h`, it holds for every smaller `h`. That monotonicity lets us binary-search the largest feasible `h` rather than scanning linearly.

### Algorithm

1. `lo = 0`, `hi = n`. While `lo < hi`: `mid = (lo+hi+1)/2` (bias high).
2. Count papers with `>= mid` citations. If `count >= mid`, `lo = mid`; else `hi = mid-1`.
3. Return `lo`.

### Complexity

- **Time:** O(n log n) — `log n` iterations, each an O(n) count.
- **Space:** O(1).

### Code

```go
func binarySearchAnswer(citations []int) int {
	n := len(citations)
	lo, hi := 0, n
	for lo < hi {
		mid := (lo + hi + 1) / 2 // bias upward to make progress toward larger h
		count := 0
		for _, c := range citations {
			if c >= mid {
				count++ // this paper qualifies for a candidate h of `mid`
			}
		}
		if count >= mid {
			lo = mid // feasible: at least mid papers with >= mid cites → try larger
		} else {
			hi = mid - 1 // infeasible → h must be smaller
		}
	}
	return lo
}
```

### Dry Run

Example 1: `citations = [3,0,6,1,5]`, `n = 5`.

| lo | hi | mid = (lo+hi+1)/2 | count(≥ mid) | count ≥ mid? | update |
|----|----|-------------------|--------------|--------------|--------|
| 0 | 5 | 3 | {3,6,5} → 3 | 3 ≥ 3 yes | lo = 3 |
| 3 | 5 | 4 | {6,5} → 2 | 2 ≥ 4 no | hi = 3 |
| 3 | 3 | — | loop ends | — | — |

Result: `lo = 3` ✔

---

## Key Takeaways

- **h-index is bounded by n.** That single observation (`h ≤ n`) is what enables the O(n) counting-bucket solution — you never need buckets beyond `n`.
- **Sort-descending + rank scan** is the intuitive O(n log n) baseline: compare each paper's citation count against its 1-based rank; stop at the first failure.
- **"At least h with ≥ h" is monotone**, so it's binary-searchable — a reusable pattern for "find the largest threshold that still satisfies a count condition".
- Counting sort beats comparison sort whenever the key range is bounded (here, capped at `n`).

---

## Related Problems

- LeetCode #275 — H-Index II (same problem on an already-sorted array → O(log n))
- LeetCode #1122 — Relative Sort Array (counting-sort style bucketing)
- LeetCode #912 — Sort an Array (sorting fundamentals)
- LeetCode #1237 — Find Positive Integer Solution (monotone-predicate search)
