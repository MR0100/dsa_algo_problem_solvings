# 0056 — Merge Intervals

> LeetCode #56 · Difficulty: Medium
> **Categories:** Array, Sorting

---

## Problem Statement

Given an array of `intervals` where `intervals[i] = [starti, endi]`, merge all overlapping intervals, and return an array of the non-overlapping intervals that cover all the intervals in the input.

**Example 1**
```
Input:  intervals = [[1,3],[2,6],[8,10],[15,18]]
Output: [[1,6],[8,10],[15,18]]
Explanation: Since intervals [1,3] and [2,6] overlap, merge them into [1,6].
```

**Example 2**
```
Input:  intervals = [[1,4],[4,5]]
Output: [[1,5]]
Explanation: Intervals [1,4] and [4,5] are considered overlapping.
```

**Constraints**
- `1 <= intervals.length <= 10⁴`
- `intervals[i].length == 2`
- `0 <= starti <= endi <= 10⁴`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★★ Very High | 2024          |
| Meta      | ★★★★★ Very High | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★★☆ High      | 2023          |
| Apple     | ★★★★☆ High      | 2023          |
| Adobe     | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Sorting** — sort by start time so overlapping intervals become adjacent.
- **Greedy Merge** — at each step, either merge the current interval into the last result interval (extend its end) or append it as a new non-overlapping interval.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Sort + Linear Merge ✅ | O(n log n) | O(n) | The standard solution; works in every interview |

---

## Approach 1 — Sort + Linear Merge (Recommended ✅)

### Intuition
If intervals are sorted by start time, overlapping intervals are always adjacent. A greedy forward scan can then merge them:
- If the current interval's start is ≤ the last result interval's end → they overlap → extend the last interval's end.
- Otherwise → no overlap → append current as a new interval.

### Algorithm
```
sort intervals by start
result = [intervals[0]]
for each interval in intervals[1..]:
  last = result[-1]
  if interval[0] <= last[1]:   // overlap
    last[1] = max(last[1], interval[1])
  else:
    result.append(interval)
return result
```

### Complexity
- **Time:** O(n log n) — sorting dominates; merge pass is O(n).
- **Space:** O(n) — output array; sort is O(log n) auxiliary.

### Code
```go
func merge(intervals [][]int) [][]int {
    sort.Slice(intervals, func(i, j int) bool { return intervals[i][0] < intervals[j][0] })
    result := [][]int{intervals[0]}
    for _, curr := range intervals[1:] {
        last := result[len(result)-1]
        if curr[0] <= last[1] {
            if curr[1] > last[1] { last[1] = curr[1] }
        } else {
            result = append(result, curr)
        }
    }
    return result
}
```

### Dry Run — `intervals = [[1,3],[2,6],[8,10],[15,18]]`
```
After sort: [[1,3],[2,6],[8,10],[15,18]]  (already sorted)

result = [[1,3]]

curr=[2,6]: 2 <= 3 (overlap) → last[1] = max(3,6) = 6. result=[[1,6]]
curr=[8,10]: 8 > 6 (no overlap) → append. result=[[1,6],[8,10]]
curr=[15,18]: 15 > 10 → append. result=[[1,6],[8,10],[15,18]]

Output: [[1,6],[8,10],[15,18]] ✓
```

---

## Key Takeaways

- **Sort first — everything becomes adjacent** — this transforms a 2D problem (any pair could overlap) into a 1D problem (only adjacent intervals can overlap after sorting).
- **`curr[0] <= last[1]` is the overlap condition** — touching intervals (end == start) count as overlapping per the problem definition.
- **Modify in-place via the result slice** — `last := result[len(result)-1]` is a slice pointing into the result; modifying `last[1]` modifies the result directly.
- **This pattern recurs constantly** — interval scheduling, calendar merging, meeting room problems all use sort + greedy merge.

---

## Related Problems

- LeetCode #57 — Insert Interval (insert into already-merged list; no sort needed)
- LeetCode #252 — Meeting Rooms (can one person attend all? Sort + check overlap)
- LeetCode #253 — Meeting Rooms II (min rooms needed; sort + min-heap)
- LeetCode #435 — Non-overlapping Intervals (greedy interval scheduling)
