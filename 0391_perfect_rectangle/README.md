# 0391 — Perfect Rectangle

> LeetCode #391 · Difficulty: Hard
> **Categories:** Array, Line Sweep, Geometry, Hash Set

---

## Problem Statement

Given an array `rectangles` where `rectangles[i] = [xi, yi, ai, bi]` represents an
axis-aligned rectangle. The bottom-left point of the rectangle is `(xi, yi)` and the
top-right point of it is `(ai, bi)`.

Return `true` if all the rectangles together form an exact cover of a rectangular region.

**Example 1:**

```
Input: rectangles = [[1,1,3,3],[3,1,4,2],[3,2,4,4],[1,3,2,4],[2,3,3,4]]
Output: true
Explanation: All 5 rectangles together form an exact cover of a rectangular region.
```

**Example 2:**

```
Input: rectangles = [[1,1,2,3],[1,3,2,4],[3,1,4,2],[3,2,4,4]]
Output: false
Explanation: Because there is a gap between the two rectangular regions.
```

**Example 3:**

```
Input: rectangles = [[1,1,3,3],[3,1,4,2],[1,3,2,4],[2,2,4,4]]
Output: false
Explanation: Because two of the rectangles overlap with each other.
```

**Constraints:**

- `1 <= rectangles.length <= 2 * 10^4`
- `rectangles[i].length == 4`
- `-10^5 <= xi < ai <= 10^5`
- `-10^5 <= yi < bi <= 10^5`

---

## Company Frequency

| Company   | Frequency         | Last Reported |
|-----------|-------------------|---------------|
| Google    | ★★★★☆ High        | 2024          |
| Amazon    | ★★★☆☆ Medium      | 2023          |
| Microsoft | ★★★☆☆ Medium      | 2023          |
| Apple     | ★★☆☆☆ Low         | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Set (corner parity)** — toggle each rectangle corner; interior corners cancel in pairs/quads, only the 4 outer corners survive → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Line Sweep over intervals** — process left/right edges by x, maintain active vertical intervals, reject overlaps → see [`/dsa/line_sweep.md`](/dsa/line_sweep.md)
- **Sorting** — order edge events by x for the sweep → see [`/dsa/sorting.md`](/dsa/sorting.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Corner Counting + Area | O(n) | O(n) | Best: linear, elegant necessary-and-sufficient test |
| 2 | Sweep Line | O(n log n) | O(n) | Geometric intuition; generalizes to overlap queries |

---

## Approach 1 — Corner Counting + Area (Optimal)

### Intuition

A tiling is perfect ⇔ two independent conditions both hold:

1. **Area** — the summed area of the small rectangles equals the area of the bounding
   box (min bottom-left to max top-right). This forbids a net shortfall (gap).
2. **Corner parity** — mark the 4 corners of every rectangle. Any point that is an
   *interior* meeting point is shared by an even number of rectangles and cancels; only
   the 4 corners of the whole bounding rectangle appear an odd number of times. So after
   XOR-toggling every corner, exactly those 4 outer corners must remain.

Area alone is fooled by "an overlap of size k balanced by a gap of size k"; the corner
parity catches precisely those cases. Together they are necessary and sufficient.

### Algorithm

1. Track bounding box `(minX, minY, maxX, maxY)` and running `area`.
2. For each rectangle: add its area; toggle each of its 4 corners in a set (present →
   delete, absent → insert).
3. After the loop, the set must have exactly the 4 bounding corners — no more, no fewer.
4. `area` must equal `(maxX-minX) * (maxY-minY)`.
5. Both true ⇒ return `true`.

### Complexity

- **Time:** O(n) — single pass, O(1) work per rectangle (4 corner toggles).
- **Space:** O(n) — the corner set can hold up to O(n) points before cancellation.

### Code

```go
func cornerCounting(rectangles [][]int) bool {
	type point struct{ x, y int }

	area := 0
	corners := map[point]bool{}
	minX, minY := 1<<62, 1<<62
	maxX, maxY := -(1 << 62), -(1 << 62)

	for _, r := range rectangles {
		x1, y1, x2, y2 := r[0], r[1], r[2], r[3]

		if x1 < minX {
			minX = x1
		}
		if y1 < minY {
			minY = y1
		}
		if x2 > maxX {
			maxX = x2
		}
		if y2 > maxY {
			maxY = y2
		}

		area += (x2 - x1) * (y2 - y1)

		for _, c := range []point{{x1, y1}, {x1, y2}, {x2, y1}, {x2, y2}} {
			if corners[c] {
				delete(corners, c)
			} else {
				corners[c] = true
			}
		}
	}

	if len(corners) != 4 {
		return false
	}
	for _, c := range []point{{minX, minY}, {minX, maxY}, {maxX, minY}, {maxX, maxY}} {
		if !corners[c] {
			return false
		}
	}

	return area == (maxX-minX)*(maxY-minY)
}
```

### Dry Run

Example 1: `[[1,1,3,3],[3,1,4,2],[3,2,4,4],[1,3,2,4],[2,3,3,4]]`

Small-rectangle areas: `[1,1,3,3]=4`, `[3,1,4,2]=1`, `[3,2,4,4]=2`, `[1,3,2,4]=1`,
`[2,3,3,4]=1` → summed area = **9**. Bounding box = `(1,1)..(4,4)` → `(4-1)*(4-1) = 9`.
Area test passes.

| Step | Rect processed | area so far | corner set after XOR toggle |
|------|----------------|-------------|-----------------------------|
| 1 | [1,1,3,3] | 4 | {(1,1),(1,3),(3,1),(3,3)} |
| 2 | [3,1,4,2] | 5 | (3,1) cancels; add (3,2),(4,1),(4,2) |
| 3 | [3,2,4,4] | 7 | (3,2),(4,2) cancel; add (3,4),(4,4) |
| 4 | [1,3,2,4] | 8 | (1,3) cancels; add (1,4),(2,3),(2,4) |
| 5 | [2,3,3,4] | 9 | (2,3),(2,4),(3,3),(3,4) cancel |

Surviving corners = `{(1,1),(4,1),(1,4),(4,4)}` — exactly the 4 bounding corners
(`len == 4`, all present). Area 9 == box area 9. Both conditions hold ⇒ **`true`**.

---

## Approach 2 — Sweep Line

### Intuition

Sweep a vertical line left → right. At each x, first close rectangles whose right edge is
here, then open rectangles whose left edge is here. For a perfect tiling, the currently
active vertical y-intervals must never overlap. If any newly opened interval overlaps an
existing active one, there is an overlap → not perfect. Combined with the area-equals-box
check (which forbids gaps), no-overlap ⇒ perfect cover.

### Algorithm

1. Build two events per rectangle: `(x1, open, y1, y2)` and `(x2, close, y1, y2)`.
2. Sort by x; at equal x process closes before opens.
3. Maintain active intervals sorted by `y1`. On each open, insert and reject if it
   overlaps the neighbour before or after it (touching is allowed).
4. On each close, remove the matching interval.
5. If no overlap ever occurs and total area equals the bounding-box area, return `true`.

### Complexity

- **Time:** O(n log n) — sorting the 2n edge events dominates.
- **Space:** O(n) — events plus the active interval list.

### Code

```go
func sweepLine(rectangles [][]int) bool {
	events := make([]event, 0, len(rectangles)*2)
	area := 0
	minY, maxY := 1<<62, -(1 << 62)
	for _, r := range rectangles {
		x1, y1, x2, y2 := r[0], r[1], r[2], r[3]
		events = append(events, event{x1, true, y1, y2})
		events = append(events, event{x2, false, y1, y2})
		area += (x2 - x1) * (y2 - y1)
		if y1 < minY {
			minY = y1
		}
		if y2 > maxY {
			maxY = y2
		}
	}

	sortEvents(events)

	active := []interval{}

	i := 0
	for i < len(events) {
		x := events[i].x
		for i < len(events) && events[i].x == x && !events[i].open {
			active = removeInterval(active, interval{events[i].y1, events[i].y2})
			i++
		}
		for i < len(events) && events[i].x == x && events[i].open {
			if !insertInterval(&active, interval{events[i].y1, events[i].y2}) {
				return false
			}
			i++
		}
	}

	minX, maxX := 1<<62, -(1 << 62)
	for _, r := range rectangles {
		if r[0] < minX {
			minX = r[0]
		}
		if r[2] > maxX {
			maxX = r[2]
		}
	}
	return area == (maxX-minX)*(maxY-minY)
}
```

### Dry Run

Example 1 events sorted by x (closes before opens at a tie):

| i | x | type | interval [y1,y2) | active after step | overlap? |
|---|---|------|------------------|-------------------|----------|
| 1 | 1 | open | [1,3) | {[1,3)} | no |
| 2 | 1 | open | [3,4) | {[1,3),[3,4)} | no (touches at 3) |
| 3 | 2 | open | [3,4) | {[1,3),[3,4),[3,4)} | no (touch) |
| 4 | 3 | close | [1,3) | {[3,4),[3,4)} | — |
| 5 | 3 | open | [1,2) | {[1,2),[3,4),[3,4)} | no |
| 6 | 3 | open | [2,4) | {[1,2),[2,4),[3,4),[3,4)} | no |
| … | … | close | … | shrinks toward empty | no |

No overlap raised; final `area = 9` equals bounding box `(1,1)-(4,4)` covered region →
`true`.

---

## Key Takeaways

- **Two orthogonal invariants**: area guards against gaps, corner-parity (or overlap
  detection) guards against overlaps. Neither alone suffices.
- **XOR / toggle a hash set** to find elements appearing an odd number of times — the same
  trick as "single number" but on 2D points.
- Perfect tiling ⇔ every interior corner shared an even count; only 4 outer corners odd.
- Sweep line turns a 2D covering question into 1D interval bookkeeping per x-slice.

---

## Related Problems

- LeetCode #223 — Rectangle Area (geometry of overlapping rectangles)
- LeetCode #218 — The Skyline Problem (line sweep over rectangle edges)
- LeetCode #850 — Rectangle Area II (sweep line + coordinate compression)
- LeetCode #136 — Single Number (XOR parity to find the odd one out)
