# 0223 — Rectangle Area

> LeetCode #223 · Difficulty: Medium
> **Categories:** Math, Geometry

---

## Problem Statement

Given the coordinates of two **rectilinear** rectangles in a 2D plane, return *the total area covered by the two rectangles*.

The first rectangle is defined by its **bottom-left** corner `(ax1, ay1)` and its **top-right** corner `(ax2, ay2)`.

The second rectangle is defined by its **bottom-left** corner `(bx1, by1)` and its **top-right** corner `(bx2, by2)`.

**Example 1:**

```
Input: ax1 = -3, ay1 = 0, ax2 = 3, ay2 = 4,
       bx1 = 0, by1 = -1, bx2 = 9, by2 = 2
Output: 45
```

Explanation: Rectangle A has area `(3−(−3))·(4−0) = 6·4 = 24`. Rectangle B has area `(9−0)·(2−(−1)) = 9·3 = 27`. Their overlap spans x `[0,3]` and y `[0,2]`, area `3·2 = 6`. Total covered `= 24 + 27 − 6 = 45`.

**Example 2:**

```
Input: ax1 = -2, ay1 = -2, ax2 = 2, ay2 = 2,
       bx1 = -2, by1 = -2, bx2 = 2, by2 = 2
Output: 16
```

Explanation: The two rectangles are identical (area 16 each) and fully overlap, so the covered area is just `16`.

**Constraints:**

- `-10^4 <= ax1 <= ax2 <= 10^4`
- `-10^4 <= ay1 <= ay2 <= 10^4`
- `-10^4 <= bx1 <= bx2 <= 10^4`
- `-10^4 <= by1 <= by2 <= 10^4`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Meta       | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |
| Apple      | ★☆☆☆☆ Rare       | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Inclusion–Exclusion Principle** — total covered area = sum of areas − their intersection, so the overlap is not double-counted → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)
- **Interval Intersection (per axis)** — the overlap rectangle is the intersection of the two `x` intervals and the two `y` intervals, `[max(lefts), min(rights)]` clamped at 0 → see [`/dsa/intervals.md`](/dsa/intervals.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Inclusion–Exclusion (Optimal) | O(1) | O(1) | The clean formula; `max(0, …)` handles disjoint cases automatically |
| 2 | Explicit Overlap-Detection Branch | O(1) | O(1) | When you also need an "do they overlap?" boolean, or prefer explicit branches |

---

## Approach 1 — Inclusion–Exclusion (Optimal)

### Intuition
Adding both rectangle areas counts the overlapping region twice, so subtract it once. The overlap is itself an axis-aligned rectangle: its x-extent runs from the *rightmost left edge* to the *leftmost right edge*, and similarly for y. If either extent is `≤ 0` the rectangles are disjoint and the overlap contributes nothing — a `max(0, …)` clamp encodes exactly that.

### Algorithm
1. `area1 = (ax2−ax1)·(ay2−ay1)`, `area2 = (bx2−bx1)·(by2−by1)`.
2. `overlapW = max(0, min(ax2,bx2) − max(ax1,bx1))`.
3. `overlapH = max(0, min(ay2,by2) − max(ay1,by1))`.
4. Return `area1 + area2 − overlapW·overlapH`.

### Complexity
- **Time:** O(1) — constant arithmetic.
- **Space:** O(1) — no allocation.

### Code
```go
func inclusionExclusion(ax1, ay1, ax2, ay2, bx1, by1, bx2, by2 int) int {
	area1 := (ax2 - ax1) * (ay2 - ay1) // area of the first rectangle
	area2 := (bx2 - bx1) * (by2 - by1) // area of the second rectangle

	// horizontal overlap: from the rightmost left-edge to the leftmost right-edge
	overlapW := maxInt(0, minInt(ax2, bx2)-maxInt(ax1, bx1))
	// vertical overlap: from the topmost bottom-edge to the bottommost top-edge
	overlapH := maxInt(0, minInt(ay2, by2)-maxInt(ay1, by1))

	overlap := overlapW * overlapH // 0 when the rectangles are disjoint

	return area1 + area2 - overlap
}
```

### Dry Run
Example 1: `A=(-3,0,3,4)`, `B=(0,-1,9,2)`.

| Quantity | Computation | Value |
|----------|-------------|-------|
| area1 | `(3−(−3))·(4−0)` | `6·4 = 24` |
| area2 | `(9−0)·(2−(−1))` | `9·3 = 27` |
| overlapW | `max(0, min(3,9) − max(−3,0))` = `max(0, 3−0)` | `3` |
| overlapH | `max(0, min(4,2) − max(0,−1))` = `max(0, 2−0)` | `2` |
| overlap | `3·2` | `6` |
| answer | `24 + 27 − 6` | **45** |

Answer `45`. ✔

---

## Approach 2 — Explicit Overlap-Detection Branch

### Intuition
Same math, but instead of the `max(0, …)` clamp we explicitly compute the four inner edges of a candidate intersection rectangle and test whether it is non-degenerate (`ix1 < ix2` and `iy1 < iy2`). Only then do we subtract its area. This makes the "do they overlap?" decision a visible boolean, handy if the caller needs it. Edge-touching rectangles (zero-area overlap) correctly fall through to no subtraction.

### Algorithm
1. Compute `area1`, `area2`.
2. Inner edges: `ix1=max(ax1,bx1)`, `iy1=max(ay1,by1)`, `ix2=min(ax2,bx2)`, `iy2=min(ay2,by2)`.
3. If `ix1 < ix2 && iy1 < iy2`, subtract `(ix2−ix1)·(iy2−iy1)`.
4. Otherwise return `area1 + area2`.

### Complexity
- **Time:** O(1).
- **Space:** O(1).

### Code
```go
func explicitOverlap(ax1, ay1, ax2, ay2, bx1, by1, bx2, by2 int) int {
	area1 := (ax2 - ax1) * (ay2 - ay1)
	area2 := (bx2 - bx1) * (by2 - by1)

	ix1 := maxInt(ax1, bx1) // left of overlap
	iy1 := maxInt(ay1, by1) // bottom of overlap
	ix2 := minInt(ax2, bx2) // right of overlap
	iy2 := minInt(ay2, by2) // top of overlap

	if ix1 < ix2 && iy1 < iy2 {
		overlap := (ix2 - ix1) * (iy2 - iy1)
		return area1 + area2 - overlap
	}
	return area1 + area2
}
```

### Dry Run
Example 1: `A=(-3,0,3,4)`, `B=(0,-1,9,2)`.

| Quantity | Computation | Value |
|----------|-------------|-------|
| area1 | `6·4` | `24` |
| area2 | `9·3` | `27` |
| ix1 | `max(−3,0)` | `0` |
| iy1 | `max(0,−1)` | `0` |
| ix2 | `min(3,9)` | `3` |
| iy2 | `min(4,2)` | `2` |
| overlap? | `0<3 && 0<2` → true | subtract |
| overlap area | `(3−0)·(2−0)` | `6` |
| answer | `24 + 27 − 6` | **45** |

Answer `45`. ✔

---

## Key Takeaways
- **Inclusion–exclusion** is the core idea for "area/count covered by two sets": add both, subtract the intersection once.
- The intersection of two axis-aligned rectangles decomposes into two independent 1D interval intersections — `[max(lefts), min(rights)]` per axis.
- A `max(0, …)` clamp elegantly folds the disjoint case into the same formula, avoiding a branch; watch that touching-edge rectangles have zero (not negative) overlap.
- All arithmetic fits comfortably in `int` given the `10^4` bounds (max area `~1.6·10^9`), but be mindful of overflow if constraints grow.

---

## Related Problems
- LeetCode #836 — Rectangle Overlap (just the boolean overlap test)
- LeetCode #850 — Rectangle Area II (union area of many rectangles; needs coordinate compression / sweep line)
- LeetCode #56 — Merge Intervals (1D analogue of interval union)
- LeetCode #391 — Perfect Rectangle (tiling check via corner/area accounting)
