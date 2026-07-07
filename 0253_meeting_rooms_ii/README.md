# 0253 — Meeting Rooms II

> LeetCode #253 · Difficulty: Medium
> **Categories:** Array, Sorting, Heap (Priority Queue), Intervals, Greedy, Prefix Sum

---

## Problem Statement

Given an array of meeting time intervals `intervals` where `intervals[i] = [starti, endi]`, return the minimum number of conference rooms required.

**Example 1:**

```
Input: intervals = [[0,30],[5,10],[15,20]]
Output: 2
```

**Example 2:**

```
Input: intervals = [[7,10],[2,4]]
Output: 1
```

**Constraints:**

- `1 <= intervals.length <= 10^4`
- `0 <= starti < endi <= 10^6`

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Google    | ★★★★★ Very High  | 2024          |
| Amazon    | ★★★★★ Very High  | 2024          |
| Facebook  | ★★★★☆ High       | 2023          |
| Microsoft | ★★★★☆ High       | 2023          |
| Bloomberg | ★★★☆☆ Medium     | 2023          |
| Uber      | ★★★☆☆ Medium     | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Intervals** — the answer is the maximum number of intervals overlapping at any instant → see [`/dsa/intervals.md`](/dsa/intervals.md)
- **Heap / priority queue** — a min-heap of end times tracks which room frees up soonest → see [`/dsa/heap_priority_queue.md`](/dsa/heap_priority_queue.md)
- **Sorting** — both approaches begin by sorting endpoints → see [`/dsa/sorting.md`](/dsa/sorting.md)
- **Greedy** — always reuse the earliest-freeing room before allocating a new one → see [`/dsa/greedy.md`](/dsa/greedy.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Min-Heap of End Times | O(n log n) | O(n) | Intuitive "reuse the soonest-free room" model |
| 2 | Sweep Line of Endpoints (Optimal) | O(n log n) | O(n) | Cleanest peak-concurrency count |

---

## Approach 1 — Min-Heap of End Times

### Intuition
Process meetings in start-time order. The room that frees up soonest is the one whose end time is smallest — the min-heap's top. When the next meeting starts: if it starts at or after that earliest end time, reuse that room (pop it, push the new end time); otherwise all rooms are busy, so allocate a new one (just push). The heap size is the number of rooms in use, and its peak is the answer.

### Algorithm
1. Sort `intervals` by start time.
2. Create an empty min-heap `h` of end times.
3. For each meeting `[s,e]`: if `h` is non-empty and `h.top() <= s`, pop (that room is now free). Then push `e`.
4. The final heap size equals the maximum concurrency required.

### Complexity
- **Time:** O(n log n) — sort plus n heap push/pop operations.
- **Space:** O(n) — the heap can hold all end times in the worst case.

### Code
```go
func minHeap(intervals [][]int) int {
	if len(intervals) == 0 {
		return 0
	}
	sort.Slice(intervals, func(a, b int) bool {
		return intervals[a][0] < intervals[b][0]
	})
	h := &intHeap{}
	heap.Init(h)
	for _, iv := range intervals {
		start, end := iv[0], iv[1]
		if h.Len() > 0 && (*h)[0] <= start {
			heap.Pop(h)
		}
		heap.Push(h, end)
	}
	return h.Len()
}
```

### Dry Run
Input `[[0,30],[5,10],[15,20]]`. Sorted by start: same order.

| meeting | heap top | reuse? (top ≤ start) | action        | heap after   | size |
|---------|----------|----------------------|---------------|--------------|------|
| [0,30]  | —        | no (empty)           | push 30       | [30]         | 1    |
| [5,10]  | 30       | 30 ≤ 5? no           | push 10       | [10,30]      | 2    |
| [15,20] | 10       | 10 ≤ 15? yes         | pop 10, push 20 | [20,30]    | 2    |

Peak / final heap size = **2**.

---

## Approach 2 — Sweep Line of Endpoints (Optimal)

### Intuition
The number of rooms needed at any instant equals the number of meetings in progress. Sort all start times and all end times independently. Sweep a pointer through the starts; before each start, release every meeting that has already ended (`end <= start`). Track the maximum live count — that is the minimum number of rooms.

### Algorithm
1. Build sorted arrays `starts` and `ends`.
2. Pointers `s = 0`, `e = 0`; `rooms = 0`, `maxRooms = 0`.
3. While `s < n`: if `starts[s] < ends[e]` a meeting begins before the next ends → `rooms++`, `s++`, update `maxRooms`. Else a meeting ends first → `rooms--`, `e++`.
4. Return `maxRooms`.

### Complexity
- **Time:** O(n log n) — two independent sorts; the merge sweep is O(n).
- **Space:** O(n) — the two endpoint arrays.

### Code
```go
func sweepLine(intervals [][]int) int {
	n := len(intervals)
	if n == 0 {
		return 0
	}
	starts := make([]int, n)
	ends := make([]int, n)
	for i, iv := range intervals {
		starts[i] = iv[0]
		ends[i] = iv[1]
	}
	sort.Ints(starts)
	sort.Ints(ends)

	rooms, maxRooms := 0, 0
	s, e := 0, 0
	for s < n {
		if starts[s] < ends[e] {
			rooms++
			s++
			if rooms > maxRooms {
				maxRooms = rooms
			}
		} else {
			rooms--
			e++
		}
	}
	return maxRooms
}
```

### Dry Run
Input `[[0,30],[5,10],[15,20]]`. `starts = [0,5,15]`, `ends = [10,20,30]`.

| step | s | e | starts[s] | ends[e] | starts[s] < ends[e]? | rooms | maxRooms |
|------|---|---|-----------|---------|----------------------|-------|----------|
| 1    | 0 | 0 | 0         | 10      | yes → start          | 1     | 1        |
| 2    | 1 | 0 | 5         | 10      | yes → start          | 2     | 2        |
| 3    | 2 | 0 | 15        | 10      | no → end             | 1     | 2        |
| 4    | 2 | 1 | 15        | 20      | yes → start          | 2     | 2        |
| —    | 3 |   | s == n, stop         |         |                      |       |          |

Answer = `maxRooms` = **2**.

---

## Key Takeaways
- "Minimum rooms" = "maximum number of intervals overlapping at any time" = peak concurrency.
- A min-heap keyed on end time is the canonical way to know which resource frees up first — reuse before allocating (greedy).
- Splitting each interval into a +1 start event and a −1 end event and sweeping in time order is a reusable pattern for concurrency/overlap counting.
- Using strict `<` when comparing a start against an end lets a meeting reuse a room the moment the prior one ends.

---

## Related Problems
- LeetCode #252 — Meeting Rooms (just detect any overlap)
- LeetCode #56 — Merge Intervals (combine overlapping intervals)
- LeetCode #1094 — Car Pooling (sweep line of capacity changes)
- LeetCode #218 — The Skyline Problem (heap sweep over building edges)
- LeetCode #732 — My Calendar III (max booking overlap over a stream)
