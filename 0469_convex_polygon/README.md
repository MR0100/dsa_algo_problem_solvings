# 0469 — Convex Polygon

> LeetCode #469 · Difficulty: Medium · 🔒 Premium
> **Categories:** Math, Geometry

---

## Problem Statement

You are given an array of points on the **X-Y** plane `points` where `points[i] = [xi, yi]`. The points form a polygon when joined sequentially in the order they are given (the last point connects back to the first).

Return `true` if this polygon is [**convex**](https://en.wikipedia.org/wiki/Convex_polygon) and `false` otherwise.

You may assume the polygon formed by given points is always a [**simple polygon**](https://en.wikipedia.org/wiki/Simple_polygon). In other words, we ensure that exactly two edges intersect at each vertex and that edges otherwise don't intersect each other.

**Example 1:**

```
Input: points = [[0,0],[0,5],[5,5],[5,0]]
Output: true
```

**Example 2:**

```
Input: points = [[0,0],[0,10],[10,10],[10,0],[5,5]]
Output: false
```

**Constraints:**

- `3 <= points.length <= 10^4`
- `points[i].length == 2`
- `-10^4 <= xi, yi <= 10^4`
- All the given points are **unique**.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★☆☆☆ Low        | 2022          |
| Adobe      | ★★☆☆☆ Low        | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **2-D cross product / orientation test** — the whole problem reduces to "does every turn along the boundary go the same way?", answered by the sign of the integer cross product of consecutive edge vectors (`left`/`right`/`straight`), with zero floating point → see [`/dsa/geometry.md`](/dsa/geometry.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Orientation-First Cross-Product Check | O(n) | O(1) | Explicit "find the direction, then enforce it" phrasing |
| 2 | Single-Pass Both-Signs Flags (Optimal) | O(n) | O(1) | Cleanest one-pass version; two booleans, early exit |

`n = points.length`.

---

## Approach 1 — Orientation-First Cross-Product Check

### Intuition

A simple polygon is convex **iff** walking its vertices in order (wrapping past the last back to the first) you always turn the **same** direction — all left turns or all right turns; going straight (three collinear points) is allowed. The direction of the turn at a vertex is the **sign** of the cross product of the two edge vectors meeting there. So fix a reference sign from the first genuine (non-collinear) turn, then require that no vertex ever turns the opposite way.

The cross product of `(B−A)` and `(C−B)` is `(B.x−A.x)(C.y−B.y) − (B.y−A.y)(C.x−B.x)`: positive = counter-clockwise, negative = clockwise, zero = collinear. With integer coordinates it is **exact** — no floats, no division.

### Algorithm

1. For each vertex `i`, form the triple `(P[i], P[(i+1)%n], P[(i+2)%n])` — indices wrap because the polygon is closed.
2. Compute `c = cross(...)`. If `c == 0`, skip (collinear step, allowed).
3. The first non-zero `c` sets the reference `sign` (+1 or −1).
4. If any later non-zero `c` has the opposite sign, the boundary bends both ways → return `false`.
5. If no conflict is found, return `true`.

### Complexity

- **Time:** O(n) — one pass; each triple is O(1) integer arithmetic.
- **Space:** O(1) — a reference-sign integer (plus a converted point slice for readability, which can be inlined to true O(1)).

### Code

```go
func orientationFirst(points [][]int) bool {
	n := len(points)
	P := toPoints(points) // convert [][]int to []point for readability
	sign := 0             // reference orientation: +1, -1, or 0 (not yet set)
	for i := 0; i < n; i++ {
		// Edges wrap: after the last vertex we return to P[0], then P[1].
		c := cross(P[i], P[(i+1)%n], P[(i+2)%n])
		if c == 0 {
			continue // collinear triple contributes no turn — allowed
		}
		cur := 1
		if c < 0 {
			cur = -1 // this vertex is a right (clockwise) turn
		}
		if sign == 0 {
			sign = cur // first real turn fixes the polygon's direction
		} else if cur != sign {
			return false // a turn in the opposite direction → concave/reflex
		}
	}
	return true // every turn agreed (or all collinear) → convex
}
```

Helper:

```go
func cross(a, b, c point) int {
	return (b.x-a.x)*(c.y-b.y) - (b.y-a.y)*(c.x-b.x)
}
```

### Dry Run

Example 2 (as used in code): `points = [[0,0],[0,10],[10,10],[10,0],[5,5]]`, `n = 5`.

| i | triple (A,B,C) | cross value | sign of turn | reference sign | conflict? |
|---|----------------|-------------|--------------|----------------|-----------|
| 0 | (0,0),(0,10),(10,10) | (0)(0)−(10)(10) = −100 | −1 | set to −1 | no |
| 1 | (0,10),(10,10),(10,0) | (10)(−10)−(0)(0) = −100 | −1 | −1 | no |
| 2 | (10,10),(10,0),(5,5) | (0)(5)−(−10)(−5) = −50 | −1 | −1 | no |
| 3 | (10,0),(5,5),(0,0) | (−5)(−5)−(5)(−5) = 25+25 = **+50** | **+1** | −1 | **yes** → return false |

The dent at `[5,5]` produces a `+` turn while all others were `−` → `false` ✔

---

## Approach 2 — Single-Pass Both-Signs Flags (Optimal)

### Intuition

Exactly the same orientation-consistency fact, phrased without a "reference sign" variable. Keep two booleans: `hasPos` (saw a counter-clockwise turn) and `hasNeg` (saw a clockwise turn). Each vertex's cross product sets one of them (or neither, if zero). The instant **both** are true, the boundary has turned left somewhere and right somewhere else — impossible for a convex polygon — so return `false` immediately. Finish the loop with at most one flag set ⇒ convex.

### Algorithm

1. `hasPos = hasNeg = false`.
2. For each vertex `i`: `c = cross(P[i], P[(i+1)%n], P[(i+2)%n])`.
   - `c > 0` → `hasPos = true`; `c < 0` → `hasNeg = true`; `c == 0` → nothing.
3. If `hasPos && hasNeg`, return `false` (bends both ways).
4. After the loop, return `true`.

### Complexity

- **Time:** O(n) — a single linear scan, early-exiting on the first sign clash.
- **Space:** O(1) — two boolean flags.

### Code

```go
func bothSignsFlags(points [][]int) bool {
	n := len(points)
	P := toPoints(points)
	hasPos, hasNeg := false, false // did we see any CCW / any CW turn?
	for i := 0; i < n; i++ {
		c := cross(P[i], P[(i+1)%n], P[(i+2)%n]) // turn at vertex (i+1)
		if c > 0 {
			hasPos = true // a counter-clockwise turn appeared
		} else if c < 0 {
			hasNeg = true // a clockwise turn appeared
		}
		if hasPos && hasNeg { // both directions present → not convex
			return false
		}
	}
	return true // at most one turn direction seen → convex
}
```

### Dry Run

Example 2: `points = [[0,0],[0,10],[10,10],[10,0],[5,5]]`.

| i | cross value | hasPos | hasNeg | both? |
|---|-------------|--------|--------|-------|
| 0 | −100 | false | true  | no |
| 1 | −100 | false | true  | no |
| 2 | −50  | false | true  | no |
| 3 | +50  | **true** | true | **yes** → return false |

Both a `−` and a `+` turn occurred → `false` ✔ (Example 1, the square, produces only `−` turns → all four iterations keep `hasPos == false`, so it returns `true`.)

---

## Key Takeaways

- **Convex ⇔ monotone turning direction.** For a *simple* polygon, "always turn the same way (or go straight)" is exactly convexity; the sign of consecutive cross products is the turn direction.
- **The integer cross product is the one geometry primitive to memorise:** `(B−A)×(C−B)`; positive = left/CCW, negative = right/CW, zero = collinear. It answers orientation without trig or floating point, so lattice-point problems stay exact.
- **Wrap the indices** (`(i+1)%n`, `(i+2)%n`) because the polygon is closed — the turns at the last two vertices involve the first vertices.
- **Allow collinear triples** (`cross == 0`): three points on one edge are still convex, so treat zero as "no information", not as a failure. Forgetting this is the classic wrong answer.
- The "both signs seen ⇒ reject" flag pattern generalises to any monotonicity check (e.g. detecting whether a sequence is entirely non-increasing or non-decreasing).

---

## Related Problems

- LeetCode #587 — Erect the Fence (convex hull; same cross-product orientation core)
- LeetCode #1266 — Minimum Time Visiting All Points (Chebyshev distance geometry)
- LeetCode #149 — Max Points on a Line (collinearity via cross product)
- LeetCode #223 — Rectangle Area (axis-aligned geometry)
- LeetCode #836 — Rectangle Overlap (geometry case analysis)
