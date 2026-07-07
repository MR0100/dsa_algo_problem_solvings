# 0435 — Non-overlapping Intervals

> LeetCode #435 · Difficulty: Medium
> **Categories:** Array, Dynamic Programming, Greedy, Sorting

---

## Problem Statement

Given an array of intervals `intervals` where `intervals[i] = [start_i, end_i]`, return *the minimum number of intervals you need to remove to make the rest of the intervals non-overlapping*.

Note that intervals which only touch at a point are **non-overlapping**. For example, `[1, 2]` and `[2, 3]` are non-overlapping.

**Example 1:**

```
Input: intervals = [[1,2],[2,3],[3,4],[1,3]]
Output: 1
Explanation: [1,3] can be removed and the rest of the intervals are non-overlapping.
```

**Example 2:**

```
Input: intervals = [[1,2],[1,2],[1,2]]
Output: 2
Explanation: You need to remove two [1,2] to make the rest of the intervals non-overlapping.
```

**Example 3:**

```
Input: intervals = [[1,2],[2,3]]
Output: 0
Explanation: You don't need to remove any of the intervals since they're already non-overlapping.
```

**Constraints:**

- `1 <= intervals.length <= 10^5`
- `intervals[i].length == 2`
- `-5 * 10^4 <= start_i < end_i <= 5 * 10^4`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Greedy** — the optimal solution is the classic activity-selection greedy: sort by end time and always keep the earliest-finishing interval, because it leaves the most room for the rest → see [`/dsa/greedy.md`](/dsa/greedy.md)
- **Intervals** — the whole problem is about interval overlap; the overlap test (`start ≥ prevEnd`, with touching endpoints allowed) and the sort-by-endpoint setup are core interval techniques → see [`/dsa/intervals.md`](/dsa/intervals.md)
- **Longest Increasing Subsequence** — the DP approach is LIS in disguise: "longest chain of non-overlapping intervals" replaces LIS's "increasing" test with "end ≤ next start" → see [`/dsa/longest_increasing_subsequence.md`](/dsa/longest_increasing_subsequence.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | DP (Longest Non-Overlapping Chain) | O(n²) | O(n) | Builds intuition (complement of max keep-set); too slow for n = 10⁵ |
| 2 | Greedy by Earliest End (Optimal) | O(n log n) | O(1) | The intended answer — activity selection, fast and simple |
| 3 | Greedy by Earliest Start (keep shorter) | O(n log n) | O(1) | Same optimum from a start-sorted view; keep the earlier-ending on clash |

`n` = number of intervals.

---

## Approach 1 — DP (Longest Non-Overlapping Chain)

### Intuition

"Minimum removals so the rest don't overlap" is the **complement** of "maximum number of intervals we can keep with no overlaps". If the largest non-overlapping subset has size `keep`, we must delete `n − keep`. Finding that max keep-set is a **longest-chain** DP: sort by start, and let `dp[i]` be the longest chain of non-overlapping intervals ending at interval `i`, formed by extending any earlier interval `j` that finishes at or before `i` starts. This is Longest Increasing Subsequence with "increasing" swapped for "non-overlapping".

### Algorithm

1. Sort intervals by start.
2. Set `dp[i] = 1` (interval `i` alone). For each `i`, for each `j < i` with `intervals[j].end ≤ intervals[i].start`: `dp[i] = max(dp[i], dp[j] + 1)`.
3. `keep = max(dp)`; return `n − keep`.

### Complexity

- **Time:** O(n²) — nested relaxation over `n` intervals × `n` predecessors.
- **Space:** O(n) — the `dp` table.

### Code

```go
func dpLongestChain(intervals [][]int) int {
	n := len(intervals)
	if n == 0 {
		return 0
	}
	sort.Slice(intervals, func(a, b int) bool {
		return intervals[a][0] < intervals[b][0]
	})
	dp := make([]int, n) // dp[i] = longest non-overlapping chain ending at i
	keep := 0
	for i := 0; i < n; i++ {
		dp[i] = 1
		for j := 0; j < i; j++ {
			// j precedes i iff it finishes at or before i starts (touching OK).
			if intervals[j][1] <= intervals[i][0] && dp[j]+1 > dp[i] {
				dp[i] = dp[j] + 1
			}
		}
		if dp[i] > keep {
			keep = dp[i]
		}
	}
	return n - keep
}
```

### Dry Run

Example 1: `[[1,2],[2,3],[3,4],[1,3]]`. After sort by start: `[1,2], [1,3], [2,3], [3,4]` (indices 0..3).

| i | interval | predecessors j with end ≤ start_i | dp[i] | keep |
|---|----------|-----------------------------------|-------|------|
| 0 | `[1,2]` | none | 1 | 1 |
| 1 | `[1,3]` | none (`[1,2]` end 2 > start 1) | 1 | 1 |
| 2 | `[2,3]` | `[1,2]` (end 2 ≤ start 2) → dp 1+1 | 2 | 2 |
| 3 | `[3,4]` | `[1,3]`(end 3≤3, dp1→2), `[2,3]`(end 3≤3, dp2→3) | 3 | 3 |

`keep = 3`, `n = 4` → removals `= 4 − 3 = 1` ✔

---

## Approach 2 — Greedy by Earliest End Time (Optimal)

### Intuition

The classic **activity-selection** argument. Sort intervals by **end** time and greedily keep an interval whenever it starts at or after the end of the last interval kept. Why keeping the earliest finisher is always safe: among any set of mutually clashing intervals, the one that ends soonest leaves the **most room** for everything to its right, so it can never be worse to keep it and drop a later-ending rival. Each time an interval clashes with the last kept one, that's a forced removal — count those directly and you skip the O(n²) DP.

### Algorithm

1. Sort intervals by end.
2. `prevEnd =` end of the first (earliest-ending) interval; that one is kept.
3. For each subsequent interval: if its `start ≥ prevEnd`, keep it and set `prevEnd =` its end; otherwise it overlaps → `removals++`.
4. Return `removals`.

### Complexity

- **Time:** O(n log n) — the sort dominates; the scan is O(n).
- **Space:** O(1) — a couple of scalars (besides the in-place sort).

### Code

```go
func greedyEarliestEnd(intervals [][]int) int {
	n := len(intervals)
	if n == 0 {
		return 0
	}
	// Sort by END time: the earliest finisher is the safest to keep.
	sort.Slice(intervals, func(a, b int) bool {
		return intervals[a][1] < intervals[b][1]
	})
	removals := 0
	prevEnd := intervals[0][1] // keep the very first (earliest-ending) interval
	for i := 1; i < n; i++ {
		if intervals[i][0] >= prevEnd {
			prevEnd = intervals[i][1] // no overlap → keep, advance boundary
		} else {
			removals++ // overlaps and ends no earlier → drop this one
		}
	}
	return removals
}
```

### Dry Run

Example 1: `[[1,2],[2,3],[3,4],[1,3]]`. After sort by end: `[1,2](2), [1,3](3), [2,3](3), [3,4](4)`.

| i | interval | prevEnd before | start ≥ prevEnd? | action | prevEnd after | removals |
|---|----------|----------------|-------------------|--------|---------------|----------|
| 0 | `[1,2]` | — | (kept as first) | keep | 2 | 0 |
| 1 | `[1,3]` | 2 | 1 ≥ 2? no | remove | 2 | 1 |
| 2 | `[2,3]` | 2 | 2 ≥ 2? yes | keep | 3 | 1 |
| 3 | `[3,4]` | 3 | 3 ≥ 3? yes | keep | 4 | 1 |

Output: `1` ✔

---

## Approach 3 — Greedy by Earliest Start (Keep Shorter on Clash)

### Intuition

The same greedy optimum, reached from a **start-sorted** view. Walking in start order, whenever the current interval overlaps the last kept one, one of the two must be removed — keep the one that **ends earlier** (it blocks the least future space) and drop the other. Because the list is sorted by start, "retain the earlier finisher" is just `prevEnd = min(prevEnd, currentEnd)`. This makes explicit that the greedy decision is fundamentally about **end** times regardless of the sort key.

### Algorithm

1. Sort intervals by start.
2. `prevEnd =` end of the first interval.
3. For each subsequent interval: if its `start ≥ prevEnd` → disjoint, keep it (`prevEnd =` its end). Else overlap → `removals++` and `prevEnd = min(prevEnd, its end)`.
4. Return `removals`.

### Complexity

- **Time:** O(n log n) — the sort dominates.
- **Space:** O(1) — scalar state.

### Code

```go
func greedyEarliestStart(intervals [][]int) int {
	n := len(intervals)
	if n == 0 {
		return 0
	}
	sort.Slice(intervals, func(a, b int) bool {
		return intervals[a][0] < intervals[b][0]
	})
	removals := 0
	prevEnd := intervals[0][1]
	for i := 1; i < n; i++ {
		if intervals[i][0] >= prevEnd {
			prevEnd = intervals[i][1] // disjoint → keep
		} else {
			removals++ // overlap → drop one
			if intervals[i][1] < prevEnd {
				prevEnd = intervals[i][1] // keep whichever ends earlier
			}
		}
	}
	return removals
}
```

### Dry Run

Example 1: `[[1,2],[2,3],[3,4],[1,3]]`. After sort by start: `[1,2], [1,3], [2,3], [3,4]`.

| i | interval | prevEnd before | start ≥ prevEnd? | action | prevEnd after | removals |
|---|----------|----------------|-------------------|--------|---------------|----------|
| 0 | `[1,2]` | — | (kept as first) | keep | 2 | 0 |
| 1 | `[1,3]` | 2 | 1 ≥ 2? no | remove; min(2,3)=2 | 2 | 1 |
| 2 | `[2,3]` | 2 | 2 ≥ 2? yes | keep | 3 | 1 |
| 3 | `[3,4]` | 3 | 3 ≥ 3? yes | keep | 4 | 1 |

Output: `1` ✔

---

## Key Takeaways

- **"Minimum removals to make disjoint" = total − maximum non-overlapping subset.** Flip a deletion problem into a keep-the-most problem; the keep-set is a classic activity-selection result.
- **Sort by end time, keep the earliest finisher.** This single greedy rule (the exchange argument: an earlier end never hurts) solves the whole family of interval-scheduling problems optimally in O(n log n).
- **Touching endpoints don't overlap here** — the test is `start ≥ prevEnd`, not `>`. Read each problem's overlap definition carefully; a `>` vs `≥` flip changes the answer.
- **The greedy choice is about ends, not starts.** Whether you sort by start or end, on a clash you retain the interval with the smaller end — the two greedy variants are the same decision viewed from opposite ends.
- The **DP (O(n²))** is great for intuition and small inputs but TLEs at `n = 10⁵`; know both, reach for greedy.

---

## Related Problems

- LeetCode #452 — Minimum Number of Arrows to Burst Balloons (same sort-by-end greedy)
- LeetCode #56 — Merge Intervals (sort-by-start interval sweep)
- LeetCode #57 — Insert Interval (interval merging)
- LeetCode #253 — Meeting Rooms II (interval overlap counting)
- LeetCode #646 — Maximum Length of Pair Chain (identical longest-chain / greedy)
- LeetCode #1288 — Remove Covered Intervals (interval domination via sorting)
