# 0436 — Find Right Interval

> LeetCode #436 · Difficulty: Medium
> **Categories:** Array, Binary Search, Sorting, Two Pointers

---

## Problem Statement

You are given an array of `intervals`, where `intervals[i] = [starti, endi]` and each `starti` is **unique**.

The **right interval** for an interval `i` is an interval `j` such that `startj >= endi` and `startj` is **minimized**. Note that `i` may equal `j`.

Return *an array of **right interval** indices for each interval `i`*. If no right interval exists for interval `i`, then put `-1` at index `i`.

**Example 1:**

```
Input: intervals = [[1,2]]
Output: [-1]
Explanation: There is only one interval in the collection, so it outputs -1.
```

**Example 2:**

```
Input: intervals = [[3,4],[2,3],[1,2]]
Output: [-1,0,1]
Explanation: There is no right interval for [3,4].
The right interval for [2,3] is [3,4] since start0 = 3 is the smallest start that is >= end1 = 3.
The right interval for [1,2] is [2,3] since start1 = 2 is the smallest start that is >= end2 = 2.
```

**Constraints:**

- `1 <= intervals.length <= 2 * 10^4`
- `intervals[i].length == 2`
- `-10^6 <= starti <= endi <= 10^6`
- The start point of each interval is **unique**.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |
| Facebook   | ★★☆☆☆ Low        | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Binary Search (lower bound)** — "smallest `start` that is `>= end_i`" is a textbook lower-bound query; sort the starts once and each lookup is O(log n) → see [`/dsa/binary_search.md`](/dsa/binary_search.md)
- **Sorting** — every efficient approach first sorts by start (and sometimes by end) so the search structure becomes monotone → see [`/dsa/sorting.md`](/dsa/sorting.md)
- **Two Pointers** — sweeping intervals in end-order lets a single forward pointer over start-order answer all queries without a per-query log factor → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **Intervals** — the problem is a relationship *between* intervals (which one begins right after another ends), the classic interval-processing shape → see [`/dsa/intervals.md`](/dsa/intervals.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(n²) | O(1) | Tiny inputs or to establish correctness; TLE at n = 2·10⁴ |
| 2 | Sort Starts + Binary Search | O(n log n) | O(n) | The standard answer; cleanest to reason about |
| 3 | Two Sorted Arrays + Two Pointers | O(n log n) | O(n) | Same asymptotics, no per-query log factor; shows the sweep idea |

---

## Approach 1 — Brute Force

### Intuition

The definition is constructive: the right interval of `i` is the qualifying `j` (with `start_j >= end_i`) whose start is smallest. So for each `i`, scan all `j`, keep the best qualifier. No cleverness — it simply materialises the definition, which is why it is O(n²).

### Algorithm

1. For each `i`: initialise `bestStart = +∞`, `answer = -1`.
2. For each `j`: if `start_j >= end_i` **and** `start_j < bestStart`, update `bestStart = start_j` and `answer = j`.
3. Set `res[i] = answer`.

### Complexity

- **Time:** O(n²) — the nested loop examines every ordered pair `(i, j)`.
- **Space:** O(1) beyond the O(n) output array.

### Code

```go
func bruteForce(intervals [][]int) []int {
	n := len(intervals)
	res := make([]int, n) // res[i] = index of the right interval of i, or -1
	for i := 0; i < n; i++ {
		endI := intervals[i][1] // the value every candidate start must reach
		best := -1              // index of the current best right interval
		bestStart := 1 << 62    // smallest qualifying start seen so far (huge sentinel)
		for j := 0; j < n; j++ {
			startJ := intervals[j][0]
			// j qualifies if its start is at least end_i; among qualifiers we
			// want the minimal start, so compare against bestStart.
			if startJ >= endI && startJ < bestStart {
				bestStart = startJ // tighten the minimum
				best = j           // remember which interval achieved it
			}
		}
		res[i] = best
	}
	return res
}
```

### Dry Run

Example 2: `intervals = [[3,4],[2,3],[1,2]]` (indices 0, 1, 2).

| i | end_i | scan j (start_j) | qualifiers (start_j ≥ end_i) | min start → answer |
|---|-------|------------------|------------------------------|--------------------|
| 0 | 4 | j0=3, j1=2, j2=1 | none (no start ≥ 4) | `-1` |
| 1 | 3 | j0=3, j1=2, j2=1 | j0 (start 3) | start 3 → `0` |
| 2 | 2 | j0=3, j1=2, j2=1 | j0 (3), j1 (2) | min start 2 → `1` |

Result: `[-1, 0, 1]` ✔

---

## Approach 2 — Sort Starts + Binary Search

### Intuition

The query "smallest `start_j` with `start_j >= end_i`" is a **lower bound** over the set of starts. Lower bounds are answered in O(log n) once the data is sorted. So pull out every `start` tagged with its **original index**, sort by start, and binary-search per interval. The tag matters because the required output is an index into the *input*, and sorting scrambles positions.

### Algorithm

1. Build `starts = [(start_i, i)]` for all `i`; sort ascending by `start`.
2. For each interval `i`, binary-search `starts` for the first pair with `start >= end_i` (lower bound via `sort.Search`).
3. If such a pair exists (`pos < n`), `res[i] = starts[pos].idx`; otherwise `res[i] = -1`.

### Complexity

- **Time:** O(n log n) — one sort (O(n log n)) plus n lower-bound searches (O(log n) each).
- **Space:** O(n) — the `(start, index)` array.

### Code

```go
func binarySearchSorted(intervals [][]int) []int {
	n := len(intervals)
	// starts[k] = {start value, original interval index}; sorting reorders these.
	type pair struct{ start, idx int }
	starts := make([]pair, n)
	for i := 0; i < n; i++ {
		starts[i] = pair{intervals[i][0], i}
	}
	// Sort ascending by start so binary search can find the lower bound.
	sort.Slice(starts, func(a, b int) bool { return starts[a].start < starts[b].start })

	res := make([]int, n)
	for i := 0; i < n; i++ {
		endI := intervals[i][1]
		// sort.Search returns the smallest index pos in [0, n] such that the
		// predicate is true; here: first start >= endI (the lower bound).
		pos := sort.Search(n, func(k int) bool { return starts[k].start >= endI })
		if pos < n {
			res[i] = starts[pos].idx // map the sorted hit back to its original index
		} else {
			res[i] = -1 // no start reaches endI → no right interval
		}
	}
	return res
}
```

### Dry Run

Example 2: `intervals = [[3,4],[2,3],[1,2]]`.

Sorted starts (value, original index): `[(1,2), (2,1), (3,0)]`.

| i | end_i | lower-bound search (first start ≥ end_i) | pos | starts[pos] | res[i] |
|---|-------|-------------------------------------------|-----|-------------|--------|
| 0 | 4 | scan 1,2,3 — none ≥ 4 | 3 (= n) | — | `-1` |
| 1 | 3 | first start ≥ 3 is `(3,0)` | 2 | idx 0 | `0` |
| 2 | 2 | first start ≥ 2 is `(2,1)` | 1 | idx 1 | `1` |

Result: `[-1, 0, 1]` ✔

---

## Approach 3 — Two Sorted Arrays + Two Pointers

### Intuition

Instead of a log-factor lookup per interval, exploit monotonicity. Process intervals from the **smallest end to the largest**. As `end` grows, the first start that reaches it can only move **rightward** in the start-sorted order — it never needs to back up. So one pointer `p` into the start-sorted list, advanced monotonically across the whole sweep, resolves every query. Two sorts set this up; the sweep itself is linear.

### Algorithm

1. Build `byStart` = indices sorted by start, and `byEnd` = indices sorted by end.
2. Walk `byEnd` from smallest end to largest, keeping a pointer `p` into `byStart` (starts at 0, only moves forward).
3. For the current interval (its end is the smallest not-yet-processed), advance `p` while `start(byStart[p]) < end`.
4. If `p < n`, its interval is the answer (`res[i] = byStart[p]`); else `res[i] = -1`.

### Complexity

- **Time:** O(n log n) — the two sorts dominate; the pointer `p` advances at most n times total (amortised O(n) sweep).
- **Space:** O(n) — the two index arrays.

### Code

```go
func twoPointers(intervals [][]int) []int {
	n := len(intervals)
	byStart := make([]int, n) // interval indices ordered by start ascending
	byEnd := make([]int, n)   // interval indices ordered by end ascending
	for i := range intervals {
		byStart[i] = i
		byEnd[i] = i
	}
	sort.Slice(byStart, func(a, b int) bool { return intervals[byStart[a]][0] < intervals[byStart[b]][0] })
	sort.Slice(byEnd, func(a, b int) bool { return intervals[byEnd[a]][1] < intervals[byEnd[b]][1] })

	res := make([]int, n)
	p := 0 // pointer into byStart; only ever moves forward across the whole sweep
	// Consider intervals from the smallest end to the largest.
	for _, i := range byEnd {
		endI := intervals[i][1]
		// Skip every start strictly smaller than endI; they can never be the
		// right interval for this end (or any larger end still to come).
		for p < n && intervals[byStart[p]][0] < endI {
			p++
		}
		if p < n {
			res[i] = byStart[p] // first start >= endI, in original-index terms
		} else {
			res[i] = -1 // all starts exhausted → no right interval
		}
	}
	return res
}
```

### Dry Run

Example 2: `intervals = [[3,4],[2,3],[1,2]]`.

- `byStart` sorted by start: indices `[2, 1, 0]` (starts 1, 2, 3).
- `byEnd` sorted by end: indices `[2, 1, 0]` (ends 2, 3, 4).

Sweep `byEnd`, pointer `p` starts at 0:

| processing i | end_i | advance p while start(byStart[p]) < end_i | p after | byStart[p] | res[i] |
|--------------|-------|-------------------------------------------|---------|------------|--------|
| 2 | 2 | start 1 < 2 → p→1; start 2 ≥ 2 stop | 1 | index 1 | `res[2]=1` |
| 1 | 3 | start 2 < 3 → p→2; start 3 ≥ 3 stop | 2 | index 0 | `res[1]=0` |
| 0 | 4 | start 3 < 4 → p→3; p = n stop | 3 (= n) | — | `res[0]=-1` |

Result assembled by index: `res = [-1, 0, 1]` ✔

---

## Key Takeaways

- **"Smallest x ≥ target" ⇒ lower bound.** The instant you see that phrasing, reach for sort + binary search (`sort.Search` in Go finds the first index where a predicate flips to true).
- **Carry the original index when you sort.** The answer is a position in the *input*; sorting destroys positions, so tag each element with where it came from.
- **A per-query log factor can collapse to a linear sweep** when the queries themselves are monotone — process ends in sorted order and a single forward pointer suffices (Approach 3). Same big-O as binary search, but a clean template worth recognising.
- **Unique starts** is the quiet gift here: it guarantees the lower bound is unambiguous, so no tie-breaking logic is needed.

---

## Related Problems

- LeetCode #435 — Non-overlapping Intervals (greedy interval scheduling)
- LeetCode #56 — Merge Intervals (sort-by-start interval processing)
- LeetCode #57 — Insert Interval (interval placement with binary search)
- LeetCode #253 — Meeting Rooms II (sweep starts and ends)
- LeetCode #1094 — Car Pooling (event sweep over intervals)
