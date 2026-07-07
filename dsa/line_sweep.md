# Line Sweep (Sweep-Line Algorithm)

> **The idea:** turn every object (interval, rectangle, event) into a pair of
> **events** on one axis, sort them by coordinate, then move an imaginary line across
> the plane processing events left-to-right while maintaining an **active set** of
> whatever the line currently touches. **Signature move:** replace "compare all pairs"
> (O(n²)) with "sort events and sweep once" (O(n log n)).

---

## What it is

A **line sweep** imagines a vertical line starting at `x = -∞` and gliding to `+∞`.
The line only *changes state* at discrete **event points** — the left and right edges
of intervals/rectangles, or the coordinates where something starts or stops. Between
events nothing changes, so you never look at those gaps; you jump from event to event.

At each event you update an **active set** — the collection of objects the sweep line
currently intersects — and read off whatever the problem asks (the current height, the
number of overlapping intervals, the exposed length). The active set is usually a
**heap**, a **balanced BST / ordered multiset**, a **counter**, or a **+1/−1 delta**
depending on what you must query.

The three ingredients, every time:

1. **Events.** Decompose each object into `(coordinate, type)` events. An interval
   `[s, e)` becomes a **start** event at `s` and an **end** event at `e`.
2. **Sort.** Order events by coordinate. Ties need a deliberate rule (see pitfalls) —
   usually process starts before ends, or ends before starts, depending on whether
   touching counts as overlap.
3. **Sweep + active set.** Walk the sorted events once; on a start, insert into the
   active set; on an end, remove; after each step, record the answer.

Sweeping is the *ordering discipline*; the active-set data structure is the *engine*.
Pick the engine by the query:

| You need to know… | Active-set engine |
|-------------------|-------------------|
| how many objects overlap right now | an integer **counter** with +1/−1 deltas |
| the current **maximum** (tallest building) | a **max-heap** or ordered multiset |
| the set of live objects / their exact members | a **balanced BST / ordered set** |
| total exposed length on the other axis | a **counted segment structure** |

---

## When to recognise it

| Signal in the problem | Why sweep helps |
|-----------------------|-----------------|
| "maximum number of overlapping intervals" | count active intervals; the peak is the answer (#253) |
| "minimum rooms / CPUs / platforms needed" | same as max overlap — each concurrent interval needs a resource |
| "outline / skyline of overlapping rectangles" | sweep x, keep the max active height; emit a point when it changes (#218) |
| "area / perimeter covered by rectangles" | sweep one axis, track covered length on the other |
| "do these rectangles tile a region perfectly" | sweep + corner/area accounting (#391) |
| "at each x, what is the tallest / count / sum" | maintain an active aggregate indexed by the sweep coordinate |
| "events with start and end times, find a conflict/peak" | the canonical start/end event decomposition |

Reach for a sweep when a brute force would compare **every pair** of intervals or
rectangles (O(n²)) and you can instead **order the boundaries** and process each once.
If there are no natural start/end boundaries to sort on, a sweep probably isn't it.

---

## General templates / pseudocode

### 1. Max overlap via signed deltas (no heap) — Meeting Rooms II

The lightest sweep: forget the objects, keep only a **±1 timeline**. `+1` at every
start, `−1` at every end; sweep coordinates in order and track the running sum's peak.

```go
// minMeetingRooms = the maximum number of simultaneously-active intervals.
func minMeetingRooms(intervals [][]int) int {
    n := len(intervals)
    starts := make([]int, n)
    ends := make([]int, n)
    for i, iv := range intervals {
        starts[i], ends[i] = iv[0], iv[1]
    }
    sort.Ints(starts)
    sort.Ints(ends)

    rooms, maxRooms := 0, 0
    s, e := 0, 0
    for s < n {
        if starts[s] < ends[e] { // a meeting begins before the earliest one ends
            rooms++               // need one more room
            s++
            if rooms > maxRooms {
                maxRooms = rooms
            }
        } else { // a meeting ended: free its room, then reconsider this start
            rooms--
            e++
        }
    }
    return maxRooms
}
```

Note the tie rule: `starts[s] < ends[e]` (strict) means a meeting **ending** exactly
when another **starts** frees the room in time — `[0,10]` and `[10,20]` need one room.

### 2. Explicit event list (the general shape)

When you need more than a count — e.g. to know *which* objects are active, or to emit
output as state changes — build an explicit, sorted event list.

```go
type Event struct {
    x    int // sweep coordinate
    typ  int // +1 = start (enter), -1 = end (leave)
    // ... payload: height, id, the other-axis interval, etc.
}

func sweep(events []Event) {
    sort.Slice(events, func(i, j int) bool {
        if events[i].x != events[j].x {
            return events[i].x < events[j].x
        }
        return events[i].typ > events[j].typ // tie rule: process starts before ends
    })

    active := /* heap / ordered set / counter */ nil
    for _, ev := range events {
        if ev.typ == +1 {
            // insert ev's payload into active
        } else {
            // remove ev's payload from active
        }
        // read the current state of `active` and record the answer
    }
}
```

### 3. Skyline — sweep x, keep max active height (heap)

Decompose each building `(left, right, height)` into a **start** event `(left, -h)`
and an **end** event `(right, +h)` (encode start height as negative so that, at the
same x, taller starts sort first and ends sort after starts). Keep a multiset of
active heights; whenever the current max changes, emit a key point.

```go
// getSkyline returns the key points [x, height] outlining the buildings.
func getSkyline(buildings [][]int) [][]int {
    type ev struct{ x, h int }
    var events []ev
    for _, b := range buildings {
        events = append(events, ev{b[0], -b[2]}) // start: negative height
        events = append(events, ev{b[1], b[2]})  // end:   positive height
    }
    sort.Slice(events, func(i, j int) bool {
        if events[i].x != events[j].x {
            return events[i].x < events[j].x
        }
        return events[i].h < events[j].h // starts (neg) before ends (pos) at same x
    })

    var res [][]int
    active := map[int]int{0: 1} // multiset of live heights; ground (0) always present
    prevMax := 0
    // (a real impl uses a max-heap or an ordered multiset for O(log n) max queries)
    curMaxOf := func() int {
        m := 0
        for h, c := range active {
            if c > 0 && h > m {
                m = h
            }
        }
        return m
    }
    for _, e := range events {
        if e.h < 0 { // start: add height |e.h|
            active[-e.h]++
        } else { // end: remove height e.h
            active[e.h]--
        }
        if cur := curMaxOf(); cur != prevMax { // the skyline height changed here
            res = append(res, []int{e.x, cur})
            prevMax = cur
        }
    }
    return res
}
```

> The `curMaxOf` linear scan above is written for clarity. In production replace the
> `map` + scan with a **max-heap** (lazy deletion) or an **ordered multiset**, making
> each event O(log n) and the whole sweep **O(n log n)** — see the pitfalls.

### 4. Rectangle-cover / perfect-rectangle accounting

For "do these rectangles form a perfect (gapless, overlap-free) rectangle", the sweep
mindset pairs with **corner counting**: sum every rectangle's area, track the overall
bounding box, and count how many times each corner appears. A perfect tiling requires
`Σ area == bounding-box area` **and** exactly the four outer corners appearing an odd
number of times (every interior corner is shared by an even number of rectangles).

```go
func isRectangleCover(rectangles [][]int) bool {
    area := 0
    minX, minY := math.MaxInt, math.MaxInt
    maxX, maxY := math.MinInt, math.MinInt
    corners := map[[2]int]int{} // corner -> times seen (parity is what matters)

    for _, r := range rectangles {
        x1, y1, x2, y2 := r[0], r[1], r[2], r[3]
        area += (x2 - x1) * (y2 - y1)
        minX, minY = min(minX, x1), min(minY, y1)
        maxX, maxY = max(maxX, x2), max(maxY, y2)
        for _, c := range [][2]int{{x1, y1}, {x1, y2}, {x2, y1}, {x2, y2}} {
            corners[c]++ // interior corners cancel out to even counts
        }
    }
    if area != (maxX-minX)*(maxY-minY) {
        return false // wrong total area ⇒ gap or overlap
    }
    // exactly the 4 bounding corners may have odd parity; every other corner even
    for c, cnt := range corners {
        odd := cnt%2 == 1
        isBoundCorner := (c[0] == minX || c[0] == maxX) && (c[1] == minY || c[1] == maxY) &&
            (c == [2]int{minX, minY} || c == [2]int{minX, maxY} ||
                c == [2]int{maxX, minY} || c == [2]int{maxX, maxY})
        if odd != isBoundCorner {
            return false
        }
    }
    return true
}
```

---

## Worked example — Meeting Rooms II sweep

Input: `intervals = [[0,30], [5,10], [15,20]]`. Expected: **2** rooms.

`starts = [0, 5, 15]`, `ends = [10, 20, 30]` (each sorted independently). Sweep with
two cursors `s` (starts) and `e` (ends), running `rooms`:

| step | compare `starts[s] < ends[e]` | action | s | e | rooms | maxRooms |
|------|-------------------------------|--------|---|---|-------|----------|
| init | —                             | —      | 0 | 0 | 0 | 0 |
| 1 | `0 < 10` ✓ | start → room++ | 1 | 0 | 1 | 1 |
| 2 | `5 < 10` ✓ | start → room++ | 2 | 0 | 2 | **2** |
| 3 | `15 < 10` ✗ | end → room−−   | 2 | 1 | 1 | 2 |
| 4 | `15 < 20` ✓ | start → room++ | 3 | 1 | 2 | 2 |
| —  | `s == n`, stop | | | | | |

Peak concurrency was **2** (meetings `[0,30]` and `[5,10]` overlapped, then `[0,30]`
and `[15,20]`), so 2 rooms suffice. We never compared all `3×3` interval pairs — just
walked 6 sorted boundaries.

---

## Complexity

| Phase | Cost | Why |
|-------|------|-----|
| Build events | O(n) | two events per interval/rectangle |
| Sort events | **O(n log n)** | the dominant term |
| Sweep with heap / ordered-set active set | O(n log n) | each of the O(n) events does an O(log n) insert/remove/max |
| Sweep with plain counter (deltas only) | O(n) after the sort | no per-event log factor |
| Space | O(n) | the event list plus the active set |

**Overall O(n log n)** — the whole reason to sweep instead of the O(n²) all-pairs
comparison. If you only need a *count* of overlaps you can drop the heap and stay at
sort-plus-linear.

---

## Common pitfalls

1. **Tie-breaking at equal coordinates — the #1 sweep bug.** When a start and an end share an x, the order decides whether touching intervals "overlap". If `[1,2]` and `[2,3]` should **not** count as overlapping, process **ends before starts** at equal x (or, in the two-array version, use strict `starts[s] < ends[e]`). If they *should* touch-overlap, do the reverse. Pick deliberately and write a comment.

2. **Sorting starts and ends together but losing which is which.** The signed-delta trick (`+1`/`−1`) works, but if you merge starts and ends into one sorted list you must keep the sign, or sort starts and ends into **two separate arrays** and merge with two cursors (template 1).

3. **Linear max-scan on the active set → back to O(n²).** The skyline/"current tallest" query must be O(log n). A `map` scanned every event (as in template 3's `curMaxOf`) is O(n) per event and O(n²) overall. Use a **max-heap with lazy deletion** or an **ordered multiset**.

4. **Heap without lazy deletion.** Standard binary heaps can't remove an arbitrary interior element. The idiom is **lazy deletion**: don't remove on the end event; instead pop from the top only while the top is stale (its interval has already ended). Track "ended" via a counter or an expiry key.

5. **Forgetting the ground level in skyline.** Seed the active multiset with height `0` so that when the last building ends, the skyline correctly drops to ground and emits `[x, 0]`.

6. **Emitting duplicate or redundant key points.** Only record a skyline point when the current max **actually changes** from the previous one; equal consecutive heights must be suppressed.

7. **Half-open vs. closed intervals.** Model intervals as `[start, end)` consistently. Mixing closed `[s, e]` with the tie rules above produces off-by-one overlaps at the boundaries.

8. **Coordinate scale / overflow.** Event coordinates and covered-length products (rectangle cover) can be large; keep them in 64-bit `int` and, for very large or sparse coordinates, **coordinate-compress** (map the distinct coordinates to `0..k`) before building a segment structure.

---

## Problems in this repo that use it

- [0218 — The Skyline Problem](/0218_the_skyline_problem/README.md) — the archetypal sweep: events at building edges, a max-heap/ordered-set active set, emit a key point when the tallest active height changes (template 3)
- [0253 — Meeting Rooms II](/0253_meeting_rooms_ii/README.md) — max-overlap sweep; peak concurrency = rooms needed (template 1, traced above)
- [0391 — Perfect Rectangle](/0391_perfect_rectangle/README.md) — sweep-style area + corner-parity accounting for a gapless, overlap-free tiling (template 4)

### Closely related in this repo (interval mechanics that feed a sweep)

- [0252 — Meeting Rooms](/0252_meeting_rooms/README.md) — the yes/no precursor to #253: sort by start, check for any overlap
- [0056 — Merge Intervals](/0056_merge_intervals/README.md) — sort by start, then coalesce; a one-sided sweep
- [0057 — Insert Interval](/0057_insert_interval/README.md) — merge a new interval into a sorted set

See also [`/dsa/intervals.md`](/dsa/intervals.md) for the sort-and-merge interval
toolkit, [`/dsa/heap_priority_queue.md`](/dsa/heap_priority_queue.md) for the active-set
engine, and [`/dsa/geometry.md`](/dsa/geometry.md) for the rectangle/area problems a
sweep is often paired with.

### Classics worth knowing (not necessarily in repo yet)

- LeetCode #1288 / #56 family — interval merging and removal
- LeetCode #850 — Rectangle Area II (sweep + coordinate compression + covered length)
- LeetCode #759 — Employee Free Time (merge many interval lists via sweep)
- LeetCode #1094 — Car Pooling (the pure ±1 delta timeline, no heap)
