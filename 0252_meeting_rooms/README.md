# 0252 — Meeting Rooms

> LeetCode #252 · Difficulty: Easy
> **Categories:** Array, Sorting, Intervals

---

## Problem Statement

Given an array of meeting time `intervals` where `intervals[i] = [starti, endi]`, determine if a person could attend all meetings.

**Example 1:**

```
Input: intervals = [[0,30],[5,10],[15,20]]
Output: false
Explanation: The person cannot attend the meeting [0,30] and the meeting [5,10]
at the same time (they overlap).
```

**Example 2:**

```
Input: intervals = [[7,10],[2,4]]
Output: true
Explanation: The two meetings do not overlap, so the person can attend both.
```

**Constraints:**

- `0 <= intervals.length <= 10^4`
- `intervals[i].length == 2`
- `0 <= starti < endi <= 10^6`

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Facebook  | ★★★★☆ High       | 2023          |
| Amazon    | ★★★★☆ High       | 2023          |
| Google    | ★★★☆☆ Medium     | 2023          |
| Microsoft | ★★★☆☆ Medium     | 2022          |
| Bloomberg | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Intervals** — the core question is whether any two intervals overlap → see [`/dsa/intervals.md`](/dsa/intervals.md)
- **Sorting** — sorting by start time reduces the check to adjacent pairs → see [`/dsa/sorting.md`](/dsa/sorting.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force Pairwise | O(n²) | O(1) | Tiny inputs, or to reason about the overlap condition |
| 2 | Sort by Start (Optimal) | O(n log n) | O(1) | The standard/optimal answer |

---

## Approach 1 — Brute Force Pairwise Check

### Intuition
A person can attend all meetings iff no two meetings overlap. The most direct test is to check all O(n²) pairs. Two intervals `[s1,e1]`, `[s2,e2]` overlap when `s1 < e2` and `s2 < e1` — strict, because a meeting ending at 10 and another starting at 10 do not clash.

### Algorithm
1. For every pair `i < j`, test overlap: `intervals[i][0] < intervals[j][1] && intervals[j][0] < intervals[i][1]`.
2. If any pair overlaps, return `false`.
3. Otherwise return `true`.

### Complexity
- **Time:** O(n²) — every pair is compared.
- **Space:** O(1) — no extra storage.

### Code
```go
func bruteForce(intervals [][]int) bool {
	n := len(intervals)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if intervals[i][0] < intervals[j][1] && intervals[j][0] < intervals[i][1] {
				return false
			}
		}
	}
	return true
}
```

### Dry Run
Input `[[0,30],[5,10],[15,20]]`.

| i | j | intervals[i] | intervals[j] | overlap test                 | result |
|---|---|--------------|--------------|------------------------------|--------|
| 0 | 1 | [0,30]       | [5,10]       | 0<10 && 5<30 → true          | return false |

Overlap found immediately → answer `false`.

---

## Approach 2 — Sort by Start (Optimal)

### Intuition
Once meetings are sorted by start time, a conflict can only occur between consecutive meetings: if meeting `i` does not overlap `i+1`, and starts are non-decreasing, it cannot overlap any later meeting either. So one linear sweep after sorting suffices — a meeting conflicts exactly when it starts before the previous one ended.

### Algorithm
1. Sort `intervals` ascending by start time.
2. Walk `i = 1..n-1`; if `intervals[i][0] < intervals[i-1][1]`, return `false`.
3. Return `true`.

### Complexity
- **Time:** O(n log n) — dominated by the sort; the sweep is O(n).
- **Space:** O(1) extra (in-place sort).

### Code
```go
func sortByStart(intervals [][]int) bool {
	sort.Slice(intervals, func(a, b int) bool {
		return intervals[a][0] < intervals[b][0]
	})
	for i := 1; i < len(intervals); i++ {
		if intervals[i][0] < intervals[i-1][1] {
			return false
		}
	}
	return true
}
```

### Dry Run
Input `[[0,30],[5,10],[15,20]]`. After sort by start (already sorted): `[[0,30],[5,10],[15,20]]`.

| i | intervals[i-1] | intervals[i] | test: start[i] < end[i-1] | result |
|---|----------------|--------------|---------------------------|--------|
| 1 | [0,30]         | [5,10]       | 5 < 30 → true             | return false |

Conflict at the first adjacent pair → answer `false`.

---

## Key Takeaways
- The strict overlap condition `s1 < e2 && s2 < e1` correctly allows meetings that merely touch at an endpoint.
- Sorting by start collapses an all-pairs question into an adjacent-pairs sweep — a recurring interval trick.
- This is the "detector" version of the merge-intervals family; Meeting Rooms II asks *how many* rooms rather than just yes/no.

---

## Related Problems
- LeetCode #253 — Meeting Rooms II (count concurrent meetings)
- LeetCode #56 — Merge Intervals (combine overlapping intervals)
- LeetCode #57 — Insert Interval (merge one interval in)
- LeetCode #435 — Non-overlapping Intervals (remove minimum to de-conflict)
