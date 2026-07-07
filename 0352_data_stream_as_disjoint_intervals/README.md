# 0352 — Data Stream as Disjoint Intervals

> LeetCode #352 · Difficulty: Hard
> **Categories:** Design, Binary Search, Ordered Set, Intervals

---

## Problem Statement

Given a data stream input of non-negative integers `a1, a2, ..., an`, summarize
the numbers seen so far as a list of **disjoint intervals**.

Implement the `SummaryRanges` class:

- `SummaryRanges()` Initializes the object with an empty stream.
- `void addNum(int value)` Adds the integer `value` to the stream.
- `int[][] getIntervals()` Returns a summary of the integers in the stream
  currently as a list of disjoint intervals `[start_i, end_i]`. The answer should
  be sorted by `start_i`.

**Example 1:**

```
Input
["SummaryRanges", "addNum", "getIntervals", "addNum", "getIntervals", "addNum",
 "getIntervals", "addNum", "getIntervals", "addNum", "getIntervals"]
[[], [1], [], [3], [], [7], [], [2], [], [6], []]
Output
[null, null, [[1, 1]], null, [[1, 1], [3, 3]], null, [[1, 1], [3, 3], [7, 7]],
 null, [[1, 3], [7, 7]], null, [[1, 3], [6, 7]]]

Explanation
SummaryRanges summaryRanges = new SummaryRanges();
summaryRanges.addNum(1);      // arr = [1]
summaryRanges.getIntervals(); // return [[1, 1]]
summaryRanges.addNum(3);      // arr = [1, 3]
summaryRanges.getIntervals(); // return [[1, 1], [3, 3]]
summaryRanges.addNum(7);      // arr = [1, 3, 7]
summaryRanges.getIntervals(); // return [[1, 1], [3, 3], [7, 7]]
summaryRanges.addNum(2);      // arr = [1, 2, 3, 7]
summaryRanges.getIntervals(); // return [[1, 3], [7, 7]]
summaryRanges.addNum(6);      // arr = [1, 2, 3, 6, 7]
summaryRanges.getIntervals(); // return [[1, 3], [6, 7]]
```

**Constraints:**

- `0 <= value <= 10^4`
- At most `3 * 10^4` calls will be made to `addNum` and `getIntervals`.

**Follow up:** What if there are lots of merges and the number of disjoint
intervals is small compared to the size of the data stream?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Facebook   | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Interval merging** — insert a point into a set of disjoint ranges, merging
  with left/right neighbours → see [`/dsa/intervals.md`](/dsa/intervals.md)
- **Binary search on an ordered list** — locate the insertion position in the
  sorted interval list in O(log k) → see [`/dsa/binary_search.md`](/dsa/binary_search.md)
- **Design of a stateful data structure** — maintain an invariant across a
  sequence of operations → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)

---

## Approaches Overview

| # | Approach | addNum | getIntervals | Space | When to use |
|---|----------|--------|--------------|-------|-------------|
| 1 | Brute Force (seen set, rebuild) | O(1) | O(k log k) | O(k) | Many adds, rare queries |
| 2 | Sorted disjoint intervals + binary search (Optimal) | O(log k)+splice | O(k) | O(k) | Frequent queries; small interval count (follow-up) |

`k` = number of distinct values seen so far.

---

## Approach 1 — Brute Force (Seen Set, Rebuild on Query)

### Intuition
Keep every distinct number in a set. `addNum` just inserts (dedup is free).
`getIntervals` sorts the distinct values and coalesces consecutive runs into
`[start, end]` intervals. All the work is deferred to query time.

### Algorithm
1. `AddNum(v)`: `seen[v] = true`.
2. `GetIntervals()`: copy keys into a slice, sort ascending.
3. Sweep: open a run at the first value; whenever the next value equals
   `prev + 1`, extend the run; otherwise close the run and start a new one.
4. Flush the last run.

### Complexity
- **Time:** `AddNum` O(1) amortized. `GetIntervals` O(k log k) to sort k
  distinct values.
- **Space:** O(k) for the seen set.

### Code
```go
func (s *bruteForceSummaryRanges) AddNum(value int) {
	s.seen[value] = true // duplicates are harmless: map keeps one copy
}

func (s *bruteForceSummaryRanges) GetIntervals() [][]int {
	if len(s.seen) == 0 {
		return [][]int{}
	}
	nums := make([]int, 0, len(s.seen)) // gather distinct values
	for v := range s.seen {
		nums = append(nums, v)
	}
	sort.Ints(nums) // ascending so consecutive runs are adjacent

	res := [][]int{}
	start := nums[0] // current run's start
	prev := nums[0]  // last value placed in the current run
	for i := 1; i < len(nums); i++ {
		if nums[i] == prev+1 {
			prev = nums[i] // extend the current run
			continue
		}
		res = append(res, []int{start, prev}) // gap: close the run
		start, prev = nums[i], nums[i]        // begin a new run
	}
	res = append(res, []int{start, prev}) // flush the final run
	return res
}
```

### Dry Run
Stream adds `1, 3, 7, 2, 6`; query after each.

| After add | seen set          | sorted nums     | GetIntervals sweep result |
|-----------|-------------------|-----------------|---------------------------|
| 1 | {1}               | [1]             | [[1,1]] |
| 3 | {1,3}             | [1,3]           | 3≠1+1 → [[1,1],[3,3]] |
| 7 | {1,3,7}           | [1,3,7]         | [[1,1],[3,3],[7,7]] |
| 2 | {1,2,3,7}         | [1,2,3,7]       | 2=1+1,3=2+1 run 1..3; 7 new → [[1,3],[7,7]] |
| 6 | {1,2,3,6,7}       | [1,2,3,6,7]     | run 1..3; 6,7=6+1 run 6..7 → [[1,3],[6,7]] |

---

## Approach 2 — Sorted Disjoint Intervals + Binary Search (Optimal)

### Intuition
Maintain the answer directly: a sorted list of disjoint, non-adjacent intervals.
Each new value can only interact with its immediate neighbours in that list, so a
binary search plus O(1) stitching keeps the invariant. `getIntervals` is then a
free copy — ideal for the follow-up (few intervals, many adds/queries).

### Algorithm
1. `AddNum(v)`: binary-search `idx` = first interval with `start >= v`.
2. If the left neighbour `iv[idx-1]` already covers `v` (its `end >= v`), or the
   right neighbour starts exactly at `v`, return (already covered).
3. Compute `mergeLeft = iv[idx-1].end == v-1` and
   `mergeRight = iv[idx].start == v+1`.
4. Cases:
   - both → set `iv[idx-1].end = iv[idx].end`, delete `iv[idx]` (bridge).
   - left only → `iv[idx-1].end = v` (extend right edge).
   - right only → `iv[idx].start = v` (extend left edge).
   - neither → splice new `[v, v]` at position `idx`.
5. `GetIntervals()`: return the maintained slice.

### Complexity
- **Time:** `AddNum` O(log k) search + O(k) worst-case slice splice/delete;
  `GetIntervals` O(k) (or O(1) if returning a reference).
- **Space:** O(k) — one entry per disjoint interval.

### Code
```go
func (s *SummaryRanges) AddNum(value int) {
	iv := s.intervals
	// idx = first interval whose start is >= value.
	idx := sort.Search(len(iv), func(i int) bool { return iv[i][0] >= value })

	if idx > 0 && iv[idx-1][1] >= value {
		return // value lies inside an existing interval — nothing changes
	}
	if idx < len(iv) && iv[idx][0] == value {
		return // value is the start of an existing interval — already covered
	}

	mergeLeft := idx > 0 && iv[idx-1][1] == value-1       // touches left range
	mergeRight := idx < len(iv) && iv[idx][0] == value+1  // touches right range

	switch {
	case mergeLeft && mergeRight:
		iv[idx-1][1] = iv[idx][1]
		s.intervals = append(iv[:idx], iv[idx+1:]...)
	case mergeLeft:
		iv[idx-1][1] = value // extend the left interval's end by one
	case mergeRight:
		iv[idx][0] = value // lower the right interval's start by one
	default:
		s.intervals = append(iv, nil)
		copy(s.intervals[idx+1:], s.intervals[idx:])
		s.intervals[idx] = []int{value, value}
	}
}
```

### Dry Run
Adds `1, 3, 7, 2, 6`.

| add | binary-search idx | left | right | case | intervals after |
|-----|-------------------|------|-------|------|-----------------|
| 1 | 0 (empty) | — | — | neither → splice [1,1] | [[1,1]] |
| 3 | 1 (after [1,1]) | end=1 ≠ 2 | none | neither → splice [3,3] | [[1,1],[3,3]] |
| 7 | 2 | end=3 ≠ 6 | none | neither → splice [7,7] | [[1,1],[3,3],[7,7]] |
| 2 | 1 (first start≥2 is [3,3]) | [1,1].end=1=2-1 ✓ | [3,3].start=3=2+1 ✓ | both → bridge | [[1,3],[7,7]] |
| 6 | 1 (first start≥6 is [7,7]) | [1,3].end=3 ≠ 5 | [7,7].start=7=6+1 ✓ | right → lower start | [[1,3],[6,7]] |

Final: `[[1,3],[6,7]]`, matching the expected output.

---

## Key Takeaways

- **Choose where to pay.** If queries are rare, defer work (Approach 1). If
  queries are frequent or the follow-up demands it, keep the answer materialized
  and pay a little per insert (Approach 2).
- Inserting a point into disjoint intervals has exactly **four** local cases:
  inside (no-op), extend-left, extend-right, bridge-both, or create-new.
- `sort.Search` gives the clean "first interval with start ≥ v" boundary; then
  only `idx-1` and `idx` can possibly interact with `v`.
- A Go idiom for ordered-set behaviour without a balanced BST: maintain a sorted
  slice and splice with `append`/`copy`.

---

## Related Problems

- LeetCode #56 — Merge Intervals (batch version of the same merging logic)
- LeetCode #57 — Insert Interval (insert one interval, merge overlaps)
- LeetCode #715 — Range Module (add/remove/query ranges, harder follow-up)
- LeetCode #703 — Kth Largest Element in a Stream (streaming design)
