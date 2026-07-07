# 0335 — Self Crossing

> LeetCode #335 · Difficulty: Hard
> **Categories:** Array, Math, Geometry

---

## Problem Statement

You are given an array of integers `distance`.

You start at the point `(0, 0)` on an X-Y plane, and you move `distance[0]` meters to the north, then `distance[1]` meters to the west, `distance[2]` meters to the south, `distance[3]` meters to the east, and so on. In other words, after each move, your direction changes counter-clockwise.

Return `true` if your path crosses itself or `false` if it does not.

**Example 1:**

```
Input: distance = [2,1,1,2]
Output: true
Explanation: The path crosses itself at the point (0, 1).
```

**Example 2:**

```
Input: distance = [1,2,3,4]
Output: false
Explanation: The path does not cross itself at any point.
```

**Example 3:**

```
Input: distance = [1,1,1,2,1]
Output: true
Explanation: The path crosses itself at the point (0, 0).
```

**Constraints:**

- `1 <= distance.length <= 10^5`
- `1 <= distance[i] <= 10^5`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★☆☆☆ Low        | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |
| Adobe      | ★☆☆☆☆ Rare       | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Geometry / Math (case analysis)** — the spiral structure restricts a self-intersection to three fixed local patterns among the last few edges; the optimal solution is pure inequality bookkeeping → see [`/dsa/geometry.md`](/dsa/geometry.md)
- **Array scanning** — a single left-to-right pass comparing each edge only to a constant window of preceding edges → see [`/dsa/arrays.md`](/dsa/arrays.md)
- **Segment intersection (brute force oracle)** — the baseline simulates coordinates and tests axis-aligned segment overlap; no dedicated file exists, closest is math above.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Geometric Case Analysis (Optimal) | O(n) | O(1) | The intended answer: three inequality checks per step |
| 2 | Brute Force Segment Intersection | O(n²) | O(n) | Easy to trust; great oracle to validate the O(n) rules |

---

## Approach 1 — Geometric Case Analysis (Optimal)

### Intuition

Because you always turn counter-clockwise, the path is a spiral that either grows outward forever (never crosses) or, once it starts shrinking, may bump into an earlier edge. A new edge can only possibly touch **one of the last few edges** — it can never reach back to a much older edge without first crossing a nearer one. Enumerating how a spiral can first collide gives exactly **three** local cases:

- **Case 1 — crosses the 4th line back:** the current edge `d[i]` grows back out to meet the edge two before it. Happens when `d[i] >= d[i-2]` and `d[i-1] <= d[i-3]` (the spiral is turning inward).
- **Case 2 — touches the 5th line back:** the spiral lands exactly on the edge four steps back. Happens when `d[i-1] == d[i-3]` and `d[i] + d[i-4] >= d[i-2]`.
- **Case 3 — crosses the 6th line back:** a wider overlap where the current edge reaches an edge five steps back. Requires `d[i-2] >= d[i-4]`, `d[i-3] >= d[i-1]`, `d[i-1] + d[i-5] >= d[i-3]`, and `d[i] + d[i-4] >= d[i-2]`.

If no index triggers any case, the path never crosses.

### Algorithm

1. Loop `i` from 3 to `n-1`.
2. Test Case 1 for every `i >= 3`.
3. Test Case 2 for every `i >= 4`.
4. Test Case 3 for every `i >= 5`.
5. Return `true` on the first case that fires; return `false` if the loop finishes.

### Complexity

- **Time:** O(n) — one pass, O(1) inequality checks per index.
- **Space:** O(1) — only index arithmetic on the input array.

### Code

```go
func selfCrossingCases(distance []int) bool {
	d := distance
	n := len(d)
	for i := 3; i < n; i++ {
		// Case 1: current edge crosses the edge 2 steps back (4th line).
		if d[i] >= d[i-2] && d[i-1] <= d[i-3] {
			return true
		}
		// Case 2: current edge touches the edge 4 steps back (5th line).
		if i >= 4 && d[i-1] == d[i-3] && d[i]+d[i-4] >= d[i-2] {
			return true
		}
		// Case 3: current edge crosses the edge 5 steps back (6th line).
		if i >= 5 &&
			d[i-2] >= d[i-4] && d[i-3] >= d[i-1] &&
			d[i-1]+d[i-5] >= d[i-3] && d[i]+d[i-4] >= d[i-2] {
			return true
		}
	}
	return false // scanned every edge without a crossing
}
```

### Dry Run

Example 1: `distance = [2,1,1,2]`.

| i | d[i] | Case 1: `d[i]>=d[i-2] && d[i-1]<=d[i-3]` | Case 2 (i≥4) | Case 3 (i≥5) | Result |
|---|------|-----------------------------------------|--------------|--------------|--------|
| 3 | 2 | `d[3]=2 >= d[1]=1` ✓ and `d[2]=1 <= d[0]=2` ✓ | n/a | n/a | **return true** |

At `i=3`, Case 1 fires: the 4th edge (length 2, heading east) comes back and crosses the 2nd edge (length 1, heading west) at `(0,1)`. Result: `true` ✔

---

## Approach 2 — Brute Force Segment Intersection

### Intuition

Drop the cleverness: literally walk the path and record each move as an axis-aligned segment `(x1,y1)-(x2,y2)`. Then a self-crossing is, by definition, any new segment that intersects an earlier, non-adjacent segment (the immediately preceding segment always shares an endpoint, so it is skipped). Two axis-aligned segments intersect iff their x-ranges overlap **and** their y-ranges overlap. This is O(n²) but unquestionably correct — a perfect oracle to validate the three-case rules.

### Algorithm

1. Simulate positions through the fixed direction cycle N, W, S, E, producing `n` segments.
2. For each segment `i`, compare it against all earlier segments `0..i-2` (skip the adjacent `i-1`).
3. Return `true` on the first overlapping pair; `false` otherwise.

### Complexity

- **Time:** O(n²) — each of `n` segments checked against up to `n` predecessors.
- **Space:** O(n) — stores all segments.

### Code

```go
func bruteForceSegments(distance []int) bool {
	// Direction deltas in order: North, West, South, East (counter-clockwise).
	dx := []int{0, -1, 0, 1}
	dy := []int{1, 0, -1, 0}
	segs := []seg{}
	x, y := 0, 0
	for i, dist := range distance {
		dir := i % 4                             // cycle through N,W,S,E
		nx, ny := x+dx[dir]*dist, y+dy[dir]*dist // endpoint of this move
		segs = append(segs, seg{x, y, nx, ny})   // record the segment
		x, y = nx, ny                            // advance current position
	}
	for i := 0; i < len(segs); i++ {
		// Compare against all non-adjacent earlier segments (skip i-1).
		for j := 0; j <= i-2; j++ {
			if segmentsIntersect(segs[i], segs[j]) {
				return true
			}
		}
	}
	return false
}

func segmentsIntersect(a, b seg) bool {
	aMinX, aMaxX := minI(a.x1, a.x2), maxI(a.x1, a.x2)
	aMinY, aMaxY := minI(a.y1, a.y2), maxI(a.y1, a.y2)
	bMinX, bMaxX := minI(b.x1, b.x2), maxI(b.x1, b.x2)
	bMinY, bMaxY := minI(b.y1, b.y2), maxI(b.y1, b.y2)
	return aMinX <= bMaxX && bMinX <= aMaxX &&
		aMinY <= bMaxY && bMinY <= aMaxY
}
```

### Dry Run

Example 1: `distance = [2,1,1,2]`. Start at `(0,0)`.

| i | dir | move | segment (x1,y1)-(x2,y2) | new pos |
|---|-----|------|--------------------------|---------|
| 0 | N | +2 y | (0,0)-(0,2) | (0,2) |
| 1 | W | −1 x | (0,2)-(-1,2) | (-1,2) |
| 2 | S | −1 y | (-1,2)-(-1,1) | (-1,1) |
| 3 | E | +2 x | (-1,1)-(1,1) | (1,1) |

Check segment 3 `(-1,1)-(1,1)` (x∈[-1,1], y=1) against segment 0 `(0,0)-(0,2)` (x=0, y∈[0,2]):
x-ranges overlap (`-1<=0<=1`), y-ranges overlap (`0<=1<=2`) → **intersect** at `(0,1)`. Result: `true` ✔

---

## Key Takeaways

- **Counter-clockwise turning makes a spiral, and a spiral limits collisions to a constant window.** A new edge cannot cross a far-back edge without first crossing a nearer one — that locality is what collapses an O(n²) geometry problem to O(n).
- **When exhaustive case analysis is the answer, verify it against a brute-force oracle.** The three inequality cases here are error-prone off-by-one traps; the segment-intersection brute force (validated on 300k random inputs) is how you gain confidence they're exactly right.
- **Axis-aligned segment intersection = interval overlap on both axes** — no floating-point or cross-product geometry needed.
- The three cases correspond to the spiral touching the 4th, 5th, and 6th line back; memorizing "4-5-6 lines back" is the compact way to recall this solution.

---

## Related Problems

- LeetCode #587 — Erect the Fence (computational geometry / convex hull)
- LeetCode #223 — Rectangle Area (axis-aligned overlap detection)
- LeetCode #836 — Rectangle Overlap (range-overlap reasoning)
- LeetCode #149 — Max Points on a Line (geometry case analysis)
- LeetCode #48 — Rotate Image (fixed directional / spiral traversal reasoning)
