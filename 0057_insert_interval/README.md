# 0057 — Insert Interval

> LeetCode #57 · Difficulty: Medium
> **Categories:** Array

---

## Problem Statement

You are given an array of non-overlapping intervals `intervals` sorted in ascending order by `starti`, and an interval `newInterval = [start, end]`.

Insert `newInterval` into `intervals` such that `intervals` is still sorted in ascending order by `starti` and `intervals` still does not have any overlapping intervals (merge if necessary).

Return `intervals` after the insertion.

**Note:** You don't need to modify `intervals` in-place. You can make a new array and return it.

**Example 1**
```
Input:  intervals = [[1,3],[6,9]], newInterval = [2,5]
Output: [[1,5],[6,9]]
```

**Example 2**
```
Input:  intervals = [[1,2],[3,5],[6,7],[8,10],[12,16]], newInterval = [4,8]
Output: [[1,2],[3,10],[12,16]]
Explanation: The new interval [4,8] overlaps with [3,5],[6,7],[8,10].
```

**Constraints**
- `0 <= intervals.length <= 10⁴`
- `intervals[i].length == 2`
- `0 <= starti <= endi <= 10⁵`
- `intervals` is sorted by `starti` in ascending order.
- `0 <= start <= end <= 10⁵`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★★ Very High | 2024          |
| Meta      | ★★★★☆ High      | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★★☆ High      | 2023          |
| Apple     | ★★★☆☆ Medium    | 2023          |
| Adobe     | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Three-Phase Linear Scan** — exploit the sorted, non-overlapping property: intervals before, overlapping, and after the new interval can be processed in three distinct phases.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Three-Phase Linear Scan ✅ | O(n) | O(n) | The optimal and standard answer; input already sorted |

---

## Approach 1 — Three-Phase Linear Scan (Recommended ✅)

### Intuition
Unlike #56, the input is already sorted and non-overlapping. No need to sort. Use three phases:

1. **Before**: Collect all intervals that end strictly before `newInterval` starts. These cannot overlap.
2. **Merge**: Collect all intervals that overlap with `newInterval` (their start ≤ `newInterval[1]`). Expand `newInterval` to cover them.
3. **After**: Append the remaining intervals (they start after `newInterval` ends).

**Overlap condition**: interval `[a,b]` overlaps `newInterval = [x,y]` iff `a <= y AND b >= x`. Rearranging:
- No overlap to the LEFT: `b < x` (interval ends before new starts).
- No overlap to the RIGHT: `a > y` (interval starts after new ends).
- Everything else overlaps.

### Algorithm
```
i = 0
// Phase 1: before
while i<n and intervals[i][1] < newInterval[0]: append intervals[i]; i++
// Phase 2: merge
while i<n and intervals[i][0] <= newInterval[1]:
  newInterval[0] = min(newInterval[0], intervals[i][0])
  newInterval[1] = max(newInterval[1], intervals[i][1])
  i++
append newInterval
// Phase 3: after
append intervals[i..]
```

### Complexity
- **Time:** O(n) — each interval examined at most once.
- **Space:** O(n) — output array.

### Code
```go
func insert(intervals [][]int, newInterval []int) [][]int {
    result := [][]int{}; i, n := 0, len(intervals)
    for i < n && intervals[i][1] < newInterval[0] { result = append(result, intervals[i]); i++ }
    for i < n && intervals[i][0] <= newInterval[1] {
        if intervals[i][0] < newInterval[0] { newInterval[0] = intervals[i][0] }
        if intervals[i][1] > newInterval[1] { newInterval[1] = intervals[i][1] }
        i++
    }
    result = append(result, newInterval)
    return append(result, intervals[i:]...)
}
```

### Dry Run — `intervals = [[1,2],[3,5],[6,7],[8,10],[12,16]]`, `newInterval = [4,8]`
```
Phase 1: [1,2]: 2 < 4? yes → append. [3,5]: 5 < 4? no → stop. result=[[1,2]]

Phase 2: newInterval=[4,8]
  [3,5]: 3 <= 8? yes. min(4,3)=3, max(8,5)=8. newInterval=[3,8]. i++
  [6,7]: 6 <= 8? yes. min(3,6)=3, max(8,7)=8. newInterval=[3,8]. i++
  [8,10]: 8 <= 8? yes. min(3,8)=3, max(8,10)=10. newInterval=[3,10]. i++
  [12,16]: 12 <= 10? no → stop.
Append [3,10]. result=[[1,2],[3,10]]

Phase 3: append [[12,16]]. result=[[1,2],[3,10],[12,16]] ✓
```

---

## Key Takeaways

- **No sort needed** — the input invariant (sorted, non-overlapping) means a single O(n) pass is possible, beating the O(n log n) approach used in #56.
- **Three clear phases** — adding this structure to the solution makes edge cases (empty input, new interval before all, new interval after all, new interval swallowed by one, etc.) fall out naturally.
- **`append(result, intervals[i:]...)` to copy the tail** — in Go, `...` unpacks a slice as variadic args to append, making this a clean O(1) pointer copy.

---

## Related Problems

- LeetCode #56 — Merge Intervals (unsorted input; sort first)
- LeetCode #715 — Range Module (dynamic interval insertion/removal)
- LeetCode #352 — Data Stream as Disjoint Intervals (online interval merging)
