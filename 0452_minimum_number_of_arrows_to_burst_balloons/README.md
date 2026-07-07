# 0452 — Minimum Number of Arrows to Burst Balloons

> LeetCode #452 · Difficulty: Medium
> **Categories:** Array, Greedy, Sorting, Intervals

---

## Problem Statement

There are some spherical balloons taped onto a flat wall that represents the XY-plane. The balloons are represented as a 2D integer array `points` where `points[i] = [xstart, xend]` denotes a balloon whose **horizontal diameter** stretches between `xstart` and `xend`. You do not know the exact y-coordinates of the balloons.

Arrows can be shot up **directly vertically** (in the positive y-direction) from different points along the x-axis. A balloon with `xstart` and `xend` is **burst** by an arrow shot at `x` if `xstart <= x <= xend`. There is **no limit** to the number of arrows that can be shot. A shot arrow keeps traveling up infinitely, bursting any balloons in its path.

Given the array `points`, return *the **minimum** number of arrows that must be shot to burst all balloons*.

**Example 1:**

```
Input: points = [[10,16],[2,8],[1,6],[7,12]]
Output: 2
Explanation: The balloons can be burst by 2 arrows:
- Shoot an arrow at x = 6, bursting the balloons [2,8] and [1,6].
- Shoot an arrow at x = 11, bursting the balloons [10,16] and [7,12].
```

**Example 2:**

```
Input: points = [[1,2],[3,4],[5,6],[7,8]]
Output: 4
Explanation: One arrow needs to be shot for each balloon for a total of 4 arrows.
```

**Example 3:**

```
Input: points = [[1,2],[2,3],[3,4],[4,5]]
Output: 2
Explanation: The balloons can be burst by 2 arrows:
- Shoot an arrow at x = 2, bursting the balloons [1,2] and [2,3].
- Shoot an arrow at x = 4, bursting the balloons [3,4] and [4,5].
```

**Constraints:**

- `1 <= points.length <= 10^5`
- `points[i].length == 2`
- `-2^31 <= xstart < xend <= 2^31 - 1`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Meta       | ★★★☆☆ Medium     | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |
| Salesforce | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Intervals (overlap / activity selection)** — each balloon is a closed interval; one arrow can burst a set of balloons iff they share a common point, i.e. their intervals mutually overlap. The answer is the minimum number of "stabbing points" to hit every interval, the classic interval-partitioning result → see [`/dsa/intervals.md`](/dsa/intervals.md)
- **Greedy** — sort by an endpoint and commit an arrow at the earliest end; local optimality (place the arrow as far left as it can still reach the current balloon) yields the global minimum → see [`/dsa/greedy.md`](/dsa/greedy.md)
- **Sorting** — both approaches begin by sorting the intervals (by end, or by start) to expose the overlap structure in a single left-to-right sweep → see [`/dsa/sorting.md`](/dsa/sorting.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Greedy, Sort by End (Optimal) | O(n log n) | O(1) extra | The canonical answer; equivalent to counting non-overlapping intervals |
| 2 | Greedy, Sort by Start (Shrink Overlap) | O(n log n) | O(1) extra | Same result via the "running intersection" viewpoint; handy intuition |

---

## Approach 1 — Greedy, Sort by End (Optimal)

### Intuition

One arrow shot at position `x` bursts every balloon whose interval contains `x`. So the task is: cover all intervals with as few stabbing points as possible. Greedy strategy — **sort balloons by their right edge** and fire the first arrow at the smallest right edge. That is the furthest-left position guaranteed to still hit the first balloon, which leaves the arrow maximal reach over the balloons that come after. Every balloon whose `start ≤ arrowX` is already popped; the first balloon whose `start` is strictly greater needs a fresh arrow, again placed at *its* right edge. Because we sorted by end, whenever `start ≤ arrowX` we automatically have `arrowX ≤ end`, so the arrow really is inside `[start, end]`.

> Comparing endpoints directly (never computing midpoints or `start+end`) sidesteps the 32-bit overflow the constraints hint at.

### Algorithm

1. If `points` is empty, return `0`.
2. Sort `points` by `xend` ascending.
3. Initialise `arrows = 1` and `arrowX = points[0][1]`.
4. For each subsequent balloon `[start, end]`:
   - If `start > arrowX`, this balloon is not covered → `arrows++`, `arrowX = end`.
   - Otherwise it is already burst by the current arrow → skip.
5. Return `arrows`.

### Complexity

- **Time:** O(n log n) — the sort dominates; the single sweep is O(n).
- **Space:** O(1) extra — sorting is in place; only O(log n) recursion stack.

### Code

```go
func greedySortByEnd(points [][]int) int {
	if len(points) == 0 {
		return 0 // nothing to burst
	}

	// Sort by the right edge so we always know the earliest place an arrow
	// "must" go to still catch the current balloon.
	sort.Slice(points, func(i, j int) bool {
		return points[i][1] < points[j][1] // ascending by xend
	})

	arrows := 1                 // the first balloon always needs an arrow
	arrowX := points[0][1]      // fire it at the first balloon's right edge
	for i := 1; i < len(points); i++ {
		start := points[i][0]
		// If this balloon starts strictly after our current arrow's x, the
		// arrow (which sits at arrowX) cannot reach it — need a fresh arrow.
		if start > arrowX {
			arrows++             // one more arrow required
			arrowX = points[i][1] // place it at this balloon's right edge
		}
		// Otherwise start <= arrowX <= end (because sorted by end, end >= arrowX),
		// so the current arrow already bursts this balloon — do nothing.
	}
	return arrows
}
```

### Dry Run

Example 1: `points = [[10,16],[2,8],[1,6],[7,12]]`.

After sorting by end: `[[1,6],[2,8],[7,12],[10,16]]`.

| Step | balloon | start | arrowX before | start > arrowX? | Action | arrows | arrowX after |
|------|---------|-------|---------------|-----------------|--------|--------|--------------|
| init | `[1,6]` | — | — | — | first arrow at end 6 | 1 | 6 |
| 1 | `[2,8]` | 2 | 6 | no (2 ≤ 6) | burst by current arrow | 1 | 6 |
| 2 | `[7,12]` | 7 | 6 | yes (7 > 6) | new arrow at end 12 | 2 | 12 |
| 3 | `[10,16]` | 10 | 12 | no (10 ≤ 12) | burst by current arrow | 2 | 12 |

Result: `2` ✔ — arrows at x = 6 and x = 12.

---

## Approach 2 — Greedy, Sort by Start (Shrink Overlap)

### Intuition

Same problem seen as maintaining the **shared overlap** of the current group. Sort by start. Keep a value `curEnd` = the right edge of the region that one arrow could still cover for every balloon in the current group. For each next balloon: if its `start` is still `≤ curEnd`, the group still shares a common point, so tighten `curEnd = min(curEnd, end)` (the shared point can be no further right than the smallest end). If its `start` exceeds `curEnd`, the overlap is broken — that group's arrow is committed and a new group (new arrow) begins. This is the "intersection of intervals" mirror image of Approach 1 and gives the identical count.

### Algorithm

1. If `points` is empty, return `0`.
2. Sort by `xstart` ascending (break ties by `xend`).
3. Initialise `arrows = 1` and `curEnd = points[0][1]`.
4. For each next `[s, e]`:
   - If `s > curEnd` → new group: `arrows++`, `curEnd = e`.
   - Else same group: `curEnd = min(curEnd, e)` (shrink the shared overlap).
5. Return `arrows`.

### Complexity

- **Time:** O(n log n) — sorting dominates; the sweep is O(n).
- **Space:** O(1) extra + O(log n) sort stack.

### Code

```go
func greedySortByStart(points [][]int) int {
	if len(points) == 0 {
		return 0
	}

	// Sort by left edge; ties broken by right edge for determinism.
	sort.Slice(points, func(i, j int) bool {
		if points[i][0] != points[j][0] {
			return points[i][0] < points[j][0]
		}
		return points[i][1] < points[j][1]
	})

	arrows := 1              // first balloon opens the first group
	curEnd := points[0][1]   // the group's shared overlap currently ends here
	for i := 1; i < len(points); i++ {
		s, e := points[i][0], points[i][1]
		if s > curEnd {
			// This balloon starts past the current group's overlap → the
			// group's single arrow can't reach it; commit a new arrow.
			arrows++
			curEnd = e // the new group's overlap starts as this balloon's span
		} else {
			// Still overlapping the group; the shared point can only be as far
			// right as the smallest end seen so far.
			if e < curEnd {
				curEnd = e // shrink overlap to keep it valid for all in group
			}
		}
	}
	return arrows
}
```

### Dry Run

Example 1: `points = [[10,16],[2,8],[1,6],[7,12]]`.

After sorting by start: `[[1,6],[2,8],[7,12],[10,16]]`.

| Step | balloon | s | curEnd before | s > curEnd? | Action | arrows | curEnd after |
|------|---------|---|---------------|-------------|--------|--------|--------------|
| init | `[1,6]` | — | — | — | open group | 1 | 6 |
| 1 | `[2,8]` | 2 | 6 | no | same group, shrink to min(6,8) | 1 | 6 |
| 2 | `[7,12]` | 7 | 6 | yes | new group | 2 | 12 |
| 3 | `[10,16]` | 10 | 12 | no | same group, shrink to min(12,16) | 2 | 12 |

Result: `2` ✔ — two groups, hence two arrows.

---

## Key Takeaways

- **"Minimum arrows/points to stab all intervals" = maximum non-overlapping intervals.** Both equal "number of groups when you sweep sorted-by-end and start a new group each time an interval falls outside the last committed point."
- **Sort by end, commit greedily at the end.** Placing the arrow at the earliest end keeps it as far left as possible, preserving reach for later balloons — the exchange argument proves no strategy does better.
- **Compare endpoints, never midpoints.** With `xstart`/`xend` spanning the full `int32` range, `mid = (a+b)/2` or `a+b` would overflow; endpoint comparisons avoid it entirely (and Go's `int` is 64-bit here anyway, but the habit matters).
- Two symmetric framings — *sort by end + track the arrow position* vs *sort by start + track the shrinking overlap* — give the same answer; recognising both cements the interval-overlap pattern.

---

## Related Problems

- LeetCode #435 — Non-overlapping Intervals (remove fewest intervals = n − arrows)
- LeetCode #56 — Merge Intervals (combine overlapping intervals)
- LeetCode #57 — Insert Interval (merge one interval into a sorted set)
- LeetCode #253 — Meeting Rooms II (minimum resources for overlapping intervals)
- LeetCode #1288 — Remove Covered Intervals (interval containment via sorting)
