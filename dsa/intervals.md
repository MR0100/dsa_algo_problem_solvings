# Intervals

> **Category:** Array technique / Greedy / Sorting
> **Difficulty to master:** Medium
> **Prerequisites:** [Sorting](/dsa/sorting.md), basic greedy reasoning; some variants use a
> [Heap / Priority Queue](/dsa/heap_priority_queue.md) or [Line Sweep / Prefix Sum](/dsa/prefix_sum.md)

---

## What the concept is

An **interval** is a pair `[start, end]` describing a contiguous range on a number line —
a meeting from 9:00 to 10:30, a job from day 2 to day 5, an IP range, a segment of a
genome. Interval problems ask you to reason about how a *collection* of such ranges
relate to each other: do they overlap, can they be merged, how many coexist at one
moment, which subset can be kept without conflicts?

Almost every interval problem reduces to one master insight:

> **Sort the intervals (usually by start), then walk them left-to-right while
> maintaining a small amount of state (the "current" merged interval, the earliest
> ending time so far, a heap of active end times, ...).**

Sorting converts a messy 2-D relationship ("does any pair overlap?") into a 1-D
sequential one ("does *this* interval overlap the *previous* one?"). After sorting by
start, an interval can only overlap intervals that came before it if it overlaps the
one with the **largest end seen so far** — pairwise checks collapse into a single
running comparison. That is what turns O(n²) brute force into O(n log n).

### The overlap test — memorise it

Two intervals `a = [a.start, a.end]` and `b = [b.start, b.end]` **overlap** iff:

```
a.start <= b.end  &&  b.start <= a.end        // closed intervals, touching counts
a.start <  b.end  &&  b.start <  a.end        // open ends, touching does NOT count
```

Equivalently (often easier to reason about): they do **not** overlap iff one ends
strictly before the other starts. Whether `[1,4]` and `[4,5]` "overlap" is
**problem-specific** — Merge Intervals (#56) says yes (merge them), Non-overlapping
Intervals (#435) says no (they can coexist). Getting this boundary wrong is the #1
interval bug; always check the problem's examples for a touching pair.

### The five canonical sub-patterns

| Pattern | Question asked | Core state while sweeping | Example problems |
|---|---|---|---|
| **Merge** | Combine all overlapping intervals | Current merged interval `[curStart, curEnd]` | #56, #57 |
| **Scheduling / erase** | Max non-overlapping subset (or min removals) | End of last kept interval; greedy: sort by **end** | #435, #452 |
| **Point-in-time counting** | Max simultaneous intervals ("min rooms") | Min-heap of end times, or event-sweep counter | #253, #1094 |
| **Intersection** | Common parts of two interval lists | Two pointers, one per list | #986 |
| **Coverage / gaps** | Is a range fully covered? Free time? | Furthest end reached so far | #1288, #759, #45 (disguised) |

---

## How to recognise it — signals in the problem statement

Reach for the intervals toolbox when you see:

- Input is literally pairs: `intervals[i] = [starti, endi]`, "meetings", "events",
  "bookings", "flights", "tasks with start/end times", "ranges".
- Verbs: **merge**, **insert**, **overlap**, **intersect**, **conflict**,
  **cover**, **schedule**, **book**, **attend**.
- Questions of the form:
  - "Merge all overlapping..." → merge pattern.
  - "Minimum number of intervals to remove so the rest don't overlap" /
    "maximum number of events you can attend" → greedy scheduling (sort by **end**).
  - "Minimum number of rooms / platforms / resources needed" → max concurrent
    intervals → heap or event sweep.
  - "Can a person attend all meetings?" → sort + adjacent-pair overlap check.
  - "Return the intersection of these two sorted lists" → two pointers.
- Hidden interval problems: a car's fuel range, a jump's reach (`[i, i+nums[i]]`),
  a sprinkler's watering span, a burst balloon's x-range — anything that maps an
  item to a range it *covers* is an interval problem in disguise.

**Not** this concept: "interval DP" (Burst Balloons #312, Scramble String #87) splits
an *index range of an array/string* into subproblems — that's dynamic programming on
substructure, not geometric interval overlap. Different file, different toolbox.

---

## General templates (Go)

### Template 1 — Merge overlapping intervals (the workhorse)

```go
// mergeIntervals merges all overlapping intervals in-place of the pattern:
// sort by start, then extend-or-emit.
//
// Time:  O(n log n) — sorting dominates; the sweep is O(n).
// Space: O(n) for the output (O(log n) sort stack if output not counted).
func mergeIntervals(intervals [][]int) [][]int {
    if len(intervals) == 0 {
        return intervals
    }

    // 1. Sort by start so any overlap is with the interval we are building.
    sort.Slice(intervals, func(i, j int) bool {
        return intervals[i][0] < intervals[j][0]
    })

    merged := [][]int{intervals[0]} // seed with the first interval

    for _, cur := range intervals[1:] {
        last := merged[len(merged)-1] // interval currently being built

        if cur[0] <= last[1] { // overlap (touching counts here)
            // 2a. Extend: new end is the max — cur may be nested inside last!
            if cur[1] > last[1] {
                last[1] = cur[1]
            }
        } else {
            // 2b. Gap: last is final, start building a new interval.
            merged = append(merged, cur)
        }
    }
    return merged
}
```

### Template 2 — Greedy scheduling (max non-overlapping / min removals)

```go
// maxNonOverlapping returns the size of the largest subset of mutually
// non-overlapping intervals. (Min removals = n - this.)
//
// KEY: sort by END, not start. The interval that ends earliest leaves the
// most room for everything after it — a classic exchange-argument greedy.
//
// Time:  O(n log n)   Space: O(1) extra
func maxNonOverlapping(intervals [][]int) int {
    sort.Slice(intervals, func(i, j int) bool {
        return intervals[i][1] < intervals[j][1] // by end ascending
    })

    count := 0
    lastEnd := math.MinInt // end of the last interval we kept

    for _, in := range intervals {
        if in[0] >= lastEnd { // starts after (or exactly when) last one ends
            count++          // keep it
            lastEnd = in[1]  // it becomes the new boundary
        }
        // else: skip it — it collides with a kept interval that ends earlier
    }
    return count
}
```

### Template 3 — Max concurrent intervals (event sweep / "meeting rooms")

```go
// minRooms returns the max number of intervals alive at any single moment.
//
// Idea: turn each interval into two events — +1 at start, -1 at end —
// sort all events, sweep, track the running count's peak.
//
// Time:  O(n log n)   Space: O(n) for the events
func minRooms(intervals [][]int) int {
    events := make([][2]int, 0, 2*len(intervals)) // [time, delta]
    for _, in := range intervals {
        events = append(events, [2]int{in[0], +1}, [2]int{in[1], -1})
    }
    sort.Slice(events, func(i, j int) bool {
        if events[i][0] != events[j][0] {
            return events[i][0] < events[j][0]
        }
        // Ties: process -1 before +1 so back-to-back meetings share a room.
        return events[i][1] < events[j][1]
    })

    active, peak := 0, 0
    for _, e := range events {
        active += e[1] // apply the delta
        if active > peak {
            peak = active
        }
    }
    return peak
}
```

(The min-heap-of-end-times version is equivalent: sort by start; for each interval,
pop the heap while `heap top <= start`; push its end; peak = max heap size.)

### Template 4 — Intersection of two sorted interval lists (two pointers)

```go
// intervalIntersection returns all intersections of two lists that are each
// internally disjoint and sorted by start.
//
// Time:  O(m + n)   Space: O(1) extra beyond the output
func intervalIntersection(a, b [][]int) [][]int {
    res := [][]int{}
    i, j := 0, 0
    for i < len(a) && j < len(b) {
        // Candidate intersection: latest start, earliest end.
        lo := max(a[i][0], b[j][0])
        hi := min(a[i][1], b[j][1])
        if lo <= hi { // non-empty overlap
            res = append(res, []int{lo, hi})
        }
        // Advance whichever interval ends first — it can't intersect anything else.
        if a[i][1] < b[j][1] {
            i++
        } else {
            j++
        }
    }
    return res
}
```

---

## Worked example — Merge Intervals (LeetCode #56), step by step

Input: `intervals = [[1,3],[2,6],[8,10],[15,18]]` (already sorted by start here,
but the sort step still runs). Using Template 1:

| Step | `cur` | `last` (top of `merged`) | Test `cur[0] <= last[1]` | Action | `merged` after |
|------|-------|--------------------------|--------------------------|--------|----------------|
| init | —     | —                        | —                        | seed with `[1,3]` | `[[1,3]]` |
| 1    | `[2,6]`   | `[1,3]`   | `2 <= 3` → **overlap** | extend end: `max(3,6)=6` | `[[1,6]]` |
| 2    | `[8,10]`  | `[1,6]`   | `8 <= 6` → **gap**     | append new interval      | `[[1,6],[8,10]]` |
| 3    | `[15,18]` | `[8,10]`  | `15 <= 10` → **gap**   | append new interval      | `[[1,6],[8,10],[15,18]]` |

Output: `[[1,6],[8,10],[15,18]]` ✓ — `[1,3]` and `[2,6]` merged into `[1,6]`.

Now the boundary case `[[1,4],[4,5]]`: at step 1, `cur=[4,5]`, `last=[1,4]`,
test `4 <= 4` → overlap (touching counts for #56) → extend to `[1,5]`.
Output: `[[1,5]]` ✓.

And the nesting case `[[1,10],[2,3]]`: `cur=[2,3]` overlaps `last=[1,10]`, and the
`max` on the end keeps `10` — the result stays `[[1,10]]`, not the corrupted `[[1,3]]`
you'd get by blindly assigning `last[1] = cur[1]`.

---

## Common pitfalls and how to avoid them

1. **Assigning instead of taking the max when extending.**
   `last[1] = cur[1]` breaks on nested intervals (`[1,10]` then `[2,3]` → wrongly
   `[1,3]`). Always `last[1] = max(last[1], cur[1])`. Sorting by start does *not*
   sort ends.

2. **Getting the touching boundary wrong (`<=` vs `<`).**
   Decide from the problem's own examples whether `[1,4]` and `[4,5]` overlap.
   Merge problems: usually yes (`cur[0] <= last[1]`). Scheduling problems: usually
   no (`in[0] >= lastEnd` is a valid keep). Room-counting: a meeting ending at 10
   usually frees the room for one starting at 10 — hence "-1 before +1" on ties in
   the event sweep.

3. **Sorting by the wrong key.**
   - Merge / insert / counting → sort by **start**.
   - "Max non-overlapping" / "min removals" / "min arrows" greedy → sort by **end**.
   Sorting a scheduling problem by start and greedily keeping the first-starting
   interval is wrong (a long early interval blocks many short ones); the
   earliest-*ending* interval is provably safe by an exchange argument.

4. **Forgetting to sort at all.**
   The extend-or-emit sweep is only correct on start-sorted input. If the problem
   guarantees sorted, disjoint input (like #57 Insert Interval), you may exploit it
   for O(n) — but say so explicitly.

5. **Comparing only adjacent-after-sort pairs when ends can leapfrog.**
   Checking `intervals[i]` only against `intervals[i-1]` (rather than the running
   merged end) misses overlaps like `[1,10],[2,3],[4,5]` where `[4,5]` doesn't
   overlap `[2,3]` but does overlap the merged `[1,10]`. The running `last`/`lastEnd`
   state exists precisely to fix this.

6. **Mutating the input when the caller might not expect it.**
   `sort.Slice` reorders the caller's slice, and Template 1 mutates inner slices
   (`last[1] = ...`). Fine on LeetCode; in production code, copy first.

7. **Empty input / single interval.**
   Guard `len(intervals) == 0` before seeding `merged` with `intervals[0]`.
   A single interval must round-trip unchanged.

8. **Overflow-adjacent tricks.** Some problems encode intervals with sentinel
   values (`math.MinInt` as initial `lastEnd`). If ends can themselves be
   `math.MinInt`/`math.MaxInt`, use a boolean "first" flag instead of a sentinel.

9. **Insert-interval three-phase structure.** For #57, don't re-sort and re-merge
   from scratch when the input is already sorted and disjoint — sweep in three
   phases: (a) copy intervals ending before `newInterval` starts, (b) absorb every
   interval overlapping `newInterval` by min/max, (c) copy the rest. O(n), no sort.

---

## Problems in this repo

- [0056 — Merge Intervals](/0056_merge_intervals/README.md) — the canonical merge
  pattern (Template 1): sort by start, extend-or-emit sweep.
- [0057 — Insert Interval](/0057_insert_interval/README.md) — merging one new
  interval into an already-sorted disjoint list; three-phase O(n) sweep, no sort
  needed.

*(Problems 0131–0400 are being added concurrently; a later pass will link
additional interval problems — e.g. #252/#253 Meeting Rooms, #435 Non-overlapping
Intervals, #452 Minimum Arrows, #986 Interval List Intersections — once their
folders exist.)*

Related-but-different in this repo: [0087 — Scramble String](/0087_scramble_string/README.md)
and [0095 — Unique BSTs II](/0095_unique_binary_search_trees_ii/README.md) use
**interval DP** (splitting index ranges into subproblems), which is a dynamic
programming pattern, not the geometric interval-overlap toolbox covered here.
