# 0149 — Max Points on a Line

> LeetCode #149 · Difficulty: Hard
> **Categories:** Array, Hash Table, Math, Geometry

---

## Problem Statement

Given an array of `points` where `points[i] = [xi, yi]` represents a point on the **X-Y** plane, return the maximum number of points that lie on the same straight line.

**Example 1:**
```
Input: points = [[1,1],[2,2],[3,3]]
Output: 3
```

**Example 2:**
```
Input: points = [[1,1],[3,2],[5,3],[4,1],[2,3],[1,4]]
Output: 4
```

**Constraints:**
- `1 <= points.length <= 300`
- `points[i].length == 2`
- `-10^4 <= xi, yi <= 10^4`
- All the `points` are **unique**.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★★☆☆ Medium     | 2023          |
| LinkedIn   | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Map** — bucket points by slope relative to an anchor; biggest bucket = best line through that anchor → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Math / Number Theory (GCD)** — reduce (dy, dx) to lowest terms so equal slopes hash to the identical key, avoiding float precision entirely → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)
- **Cross product collinearity test** — `(x2−x1)(y3−y1) − (y2−y1)(x3−x1) == 0` checks three points on one line with pure integer arithmetic (no division, no vertical-line special case) → see [`/dsa/geometry.md`](/dsa/geometry.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Check Every Pair's Line) | O(n³) | O(1) | n ≤ 300 passes; zero precision risk, easiest to prove correct |
| 2 | Hash Map of Float Slopes per Anchor | O(n²) | O(n) | Fast to write; safe only when coordinates are small |
| 3 | Hash Map of GCD-Normalized Slopes (Optimal) | O(n² log C) | O(n) | The robust interview answer — exact for any integer input |

---

## Approach 1 — Brute Force (Check Every Pair's Line)

### Intuition
A line (containing at least 2 of our points) is determined by a pair of points. So try every pair `(i, j)` as "the line", and count how many points `k` are collinear with it. Collinearity of three points is a **cross product** test: vectors `i→j` and `i→k` are parallel iff

```
(xj − xi)·(yk − yi) − (yj − yi)·(xk − xi) == 0
```

Integer-only math — no slopes, no division, no vertical-line special case, no precision issues. With `|coord| ≤ 10^4`, each product is at most `4·10^8`, far inside int64 (and Go's `int` is 64-bit on all modern platforms).

### Algorithm
1. If `n ≤ 2`, return `n` (one or two points always lie on a single line).
2. For every pair `i < j`:
   1. `count = 0`; for every `k`, if `cross(i, j, k) == 0`, increment `count` (both `i` and `k = j` pass the test themselves).
   2. `best = max(best, count)`.
3. Return `best`.

### Complexity
- **Time:** O(n³) — n²/2 pairs, each scanned against all n points; ~4.5M checks at n = 300, fine.
- **Space:** O(1) — a few integer counters.

### Code
```go
func bruteForce(points [][]int) int {
	n := len(points)
	if n <= 2 {
		return n
	}
	best := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			count := 0
			for k := 0; k < n; k++ {
				cross := (points[j][0]-points[i][0])*(points[k][1]-points[i][1]) -
					(points[j][1]-points[i][1])*(points[k][0]-points[i][0])
				if cross == 0 {
					count++
				}
			}
			if count > best {
				best = count
			}
		}
	}
	return best
}
```

### Dry Run
Example 1: `points = [[1,1],[2,2],[3,3]]`, `n = 3`.

| Pair (i,j) | k | Cross product computation | Collinear? | count |
|------------|---|----------------------------|------------|-------|
| (0,1) = (1,1)-(2,2) | 0 = (1,1) | (2−1)(1−1) − (2−1)(1−1) = 0 | yes | 1 |
| | 1 = (2,2) | (2−1)(2−1) − (2−1)(2−1) = 1−1 = 0 | yes | 2 |
| | 2 = (3,3) | (2−1)(3−1) − (2−1)(3−1) = 2−2 = 0 | yes | 3 |
| (0,2) | all k | same line y = x → all zero | yes | 3 |
| (1,2) | all k | same line y = x → all zero | yes | 3 |

`best = 3` → Output: `3` ✓

---

## Approach 2 — Hash Map of Float Slopes per Anchor

### Intuition
Fix one point as the **anchor**. Every other point defines a line through the anchor, and two points define the *same* line through the anchor iff they have the same slope from it. So bucket the other `n−1` points by slope in a hash map; the biggest bucket plus the anchor itself is the best line through that anchor. Repeat for every anchor — every optimal line contains some point, so it is found when that point anchors.

Slope as `float64` is the quick version. Two traps and their fixes:
- **Vertical lines** (`dx == 0`): use `+Inf` as the shared key.
- **Horizontal lines** (`dy == 0`): `0/dx` gives `-0.0` when `dx < 0`; force the key to `+0` so the bucket doesn't split.

With `|coord| ≤ 10^4` a float64 quotient is precise enough to distinguish any two different reduced slopes here; for unbounded coordinates it is **not** — that's what Approach 3 fixes.

### Algorithm
1. If `n ≤ 2` → return `n`.
2. For each anchor `i`:
   1. Fresh map `slopes: map[float64]int`.
   2. For each `j ≠ i`: compute `dx`, `dy`; key = `+Inf` if `dx == 0`, else `0` if `dy == 0`, else `dy/dx`.
   3. `slopes[key]++`; update `best = max(best, slopes[key] + 1)` (the `+1` is the anchor).
3. Return `best`.

### Complexity
- **Time:** O(n²) — n anchors × (n−1) constant-time map updates.
- **Space:** O(n) — one slope map alive at a time.

### Code
```go
func hashMapFloatSlopes(points [][]int) int {
	n := len(points)
	if n <= 2 {
		return n
	}
	best := 0
	for i := 0; i < n; i++ {
		slopes := make(map[float64]int)
		for j := 0; j < n; j++ {
			if j == i {
				continue
			}
			dx := float64(points[j][0] - points[i][0])
			dy := float64(points[j][1] - points[i][1])
			var slope float64
			switch {
			case dx == 0:
				slope = math.Inf(1) // vertical
			case dy == 0:
				slope = 0 // normalize -0.0 → +0
			default:
				slope = dy / dx
			}
			slopes[slope]++
			if slopes[slope]+1 > best {
				best = slopes[slope] + 1
			}
		}
	}
	return best
}
```

### Dry Run
Example 1: `points = [[1,1],[2,2],[3,3]]`.

| Anchor | j | (dx, dy) | slope key | `slopes` after | best (bucket+1) |
|--------|---|----------|-----------|----------------|------------------|
| (1,1)  | (2,2) | (1, 1) | 1.0   | `{1.0: 1}`     | 2 |
| (1,1)  | (3,3) | (2, 2) | 1.0   | `{1.0: 2}`     | **3** |
| (2,2)  | (1,1) | (−1, −1) | 1.0 | `{1.0: 1}`     | 3 |
| (2,2)  | (3,3) | (1, 1) | 1.0   | `{1.0: 2}`     | 3 |
| (3,3)  | (1,1) | (−2, −2) | 1.0 | `{1.0: 1}`     | 3 |
| (3,3)  | (2,2) | (−1, −1) | 1.0 | `{1.0: 2}`     | 3 |

Output: `3` ✓

---

## Approach 3 — Hash Map of GCD-Normalized Slopes (Optimal)

### Intuition
Same anchor-and-bucket strategy, but make the slope key **exact**: represent it as the reduced fraction `(dy/g, dx/g)` where `g = gcd(|dy|, |dx|)`. Then normalize the sign — if `dx < 0` (or `dx == 0 && dy < 0`), negate both — so that, e.g., `(1,2)`, `(-1,-2)` and `(2,4)` all become `(1,2)`, and every vertical line becomes `(1,0)`. Equal slopes now produce byte-identical `[2]int` keys with pure integer arithmetic: correct for arbitrarily large coordinates, immune to float rounding.

### Algorithm
1. If `n ≤ 2` → return `n`.
2. For each anchor `i`:
   1. Fresh map `slopes: map[[2]int]int`.
   2. For each `j ≠ i`:
      - `dy = yj − yi`, `dx = xj − xi`.
      - `g = gcd(|dy|, |dx|)` (positive since points are distinct); `dy /= g`, `dx /= g`.
      - If `dx < 0 || (dx == 0 && dy < 0)`: negate both (canonical sign).
      - `slopes[[2]int{dy,dx}]++`; `best = max(best, bucket+1)`.
3. Return `best`.

### Complexity
- **Time:** O(n² log C) — n² pairs, each paying a gcd on values ≤ C = 2·10⁴ (log C ≈ 15 steps).
- **Space:** O(n) — one bucket map per anchor at a time.

### Code
```go
func hashMapGCDSlopes(points [][]int) int {
	n := len(points)
	if n <= 2 {
		return n
	}
	best := 0
	for i := 0; i < n; i++ {
		slopes := make(map[[2]int]int) // reduced (dy, dx) → count
		for j := 0; j < n; j++ {
			if j == i {
				continue
			}
			dy := points[j][1] - points[i][1]
			dx := points[j][0] - points[i][0]
			g := gcd(abs(dy), abs(dx))
			dy, dx = dy/g, dx/g
			if dx < 0 || (dx == 0 && dy < 0) {
				dy, dx = -dy, -dx // canonical sign
			}
			key := [2]int{dy, dx}
			slopes[key]++
			if slopes[key]+1 > best {
				best = slopes[key] + 1
			}
		}
	}
	return best
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
```

### Dry Run
Example 1: `points = [[1,1],[2,2],[3,3]]`.

| Anchor | j | raw (dy, dx) | g | reduced | sign-fixed key | `slopes` after | best |
|--------|---|--------------|---|---------|-----------------|----------------|------|
| (1,1)  | (2,2) | (1, 1)   | 1 | (1, 1)  | (1, 1)          | `{(1,1): 1}`   | 2 |
| (1,1)  | (3,3) | (2, 2)   | 2 | (1, 1)  | (1, 1)          | `{(1,1): 2}`   | **3** |
| (2,2)  | (1,1) | (−1, −1) | 1 | (−1, −1)| (1, 1) (negated)| `{(1,1): 1}`   | 3 |
| (2,2)  | (3,3) | (1, 1)   | 1 | (1, 1)  | (1, 1)          | `{(1,1): 2}`   | 3 |
| (3,3)  | (1,1) | (−2, −2) | 2 | (−1, −1)| (1, 1) (negated)| `{(1,1): 1}`   | 3 |
| (3,3)  | (2,2) | (−1, −1) | 1 | (−1, −1)| (1, 1) (negated)| `{(1,1): 2}`   | 3 |

Every pair collapses to the single key `(1,1)` — one line. Output: `3` ✓

---

## Key Takeaways

- **Cross product beats slope for collinearity**: `(x2−x1)(y3−y1) == (y2−y1)(x3−x1)` is integer-exact, division-free, and handles vertical lines with zero special-casing. Reach for it in any geometry problem.
- **Anchor + hash-bucket** turns "max points on any line" (a global question) into n local questions: "max points on a line *through this point*" — an O(n³) → O(n²) drop.
- **Never trust float slopes** beyond small bounded coordinates; the exact alternative is the reduced fraction `(dy/g, dx/g)` with a **canonical sign** (`dx > 0`, vertical = `(1,0)`).
- The `-0.0` map-key split and the `+Inf` vertical key are the two classic float-slope landmines — know them even if you jump straight to gcd.
- LeetCode guarantees unique points here; in interview variants with duplicates, count duplicates of the anchor separately and add them to every bucket.

---

## Related Problems

- LeetCode #2280 — Minimum Lines to Represent a Line Chart (gcd-normalized slope comparison)
- LeetCode #1232 — Check If It Is a Straight Line (single cross-product collinearity test)
- LeetCode #356 — Line Reflection (hashing point sets by geometric relation)
- LeetCode #939 — Minimum Area Rectangle (hash map over coordinate pairs)
- LeetCode #447 — Number of Boomerangs (anchor point + hash map of distances, same skeleton)
