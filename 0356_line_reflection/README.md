# 0356 — Line Reflection

> LeetCode #356 · Difficulty: Medium
> **Categories:** Hash Table, Math, Two Pointers, Sorting, Geometry

---

## Problem Statement

Given `n` points on a 2D plane, find if there is such a line parallel to the y-axis that reflects the given points symmetrically. In other words, answer whether or not there exists a line that after reflecting all points over the given line, the original set of points is the same as the reflected ones.

Note that there can be **repeated** points.

**Example 1:**

```
Input: points = [[1,1],[-1,1]]
Output: true
Explanation: We can choose the line x = 0.
```

**Example 2:**

```
Input: points = [[1,1],[-1,-1]]
Output: false
Explanation: We can't choose a line.
```

**Constraints:**

- `n == points.length`
- `1 <= n <= 10^4`
- `-10^8 <= points[i][0], points[i][1] <= 10^8`

**Follow-up:** Could you do better than `O(n^2)`?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★☆☆☆ Low        | 2022          |
| Facebook   | ★★☆☆☆ Low        | 2022          |
| Microsoft  | ★★☆☆☆ Low        | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Set** — store every point keyed by `(x, y)` so the mirror-membership test is O(1) → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Two Pointers** — after sorting by `(y, x)`, converging pointers pair each point with its reflection → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **Sorting** — groups equal-y points and lays them out symmetrically so pairs align → see [`/dsa/sorting.md`](/dsa/sorting.md)
- **Math / Geometry** — the mirror line is pinned to `x = (minX + maxX) / 2`; we compare `2·mirror` to avoid fractions → see [`/dsa/geometry.md`](/dsa/geometry.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (pairwise) | O(n²) | O(1) | Baseline; clear but quadratic |
| 2 | Sorting + Two Pointers | O(n log n) | O(1) | No hashing; deterministic order |
| 3 | Hash Set (Optimal) | O(n) | O(n) | Best; single-pass verification |

---

## Approach 1 — Brute Force

### Intuition

The mirror line can only be halfway between the extreme x-coordinates, so it is **not** something to search for — it is fixed by `sum = minX + maxX` (where the true mirror is `sum/2`). Working with `sum` keeps everything in integers: a point `(x, y)` reflects to `(sum - x, y)`. Once the line is pinned, symmetry means every point owns a partner at its mirror position. The brute-force way to check "partner exists" is to scan all points.

### Algorithm

1. Scan once to find `minX` and `maxX`; set `sum = minX + maxX`.
2. For each point `p`, compute `mirrorX = sum - p.x`.
3. Linear-scan all points for a `q` with `q.x == mirrorX` and `q.y == p.y`.
4. If any point lacks a partner, return `false`; otherwise `true`.

### Complexity

- **Time:** O(n²) — each of the n points triggers an O(n) scan for its mirror.
- **Space:** O(1) — only `min`, `max`, and loop indices.

### Code

```go
func bruteForce(points [][]int) bool {
	if len(points) == 0 {
		return true // vacuously symmetric
	}
	minX, maxX := points[0][0], points[0][0] // track horizontal extent
	for _, p := range points {
		if p[0] < minX {
			minX = p[0]
		}
		if p[0] > maxX {
			maxX = p[0]
		}
	}
	sum := minX + maxX // 2 * (mirror x), stays integer
	// For every point verify its mirror partner is present somewhere.
	for _, p := range points {
		mirrorX := sum - p[0] // x-coordinate the partner must have
		found := false
		for _, q := range points { // linear scan for the partner
			if q[0] == mirrorX && q[1] == p[1] {
				found = true
				break
			}
		}
		if !found {
			return false // this point had no reflection — not symmetric
		}
	}
	return true
}
```

### Dry Run

Example 1: `points = [[1,1],[-1,1]]`.

| Step | Detail | State |
|------|--------|-------|
| 1 | Find extremes | minX = -1, maxX = 1, sum = 0 |
| 2 | p = (1,1), mirrorX = 0 - 1 = -1 | scan finds (-1,1) → found |
| 3 | p = (-1,1), mirrorX = 0 - (-1) = 1 | scan finds (1,1) → found |
| 4 | all points partnered | return **true** ✔ |

---

## Approach 2 — Sorting + Two Pointers

### Intuition

Sort points by `y` then `x`. All points sharing a `y` cluster together and are laid out left→right, so their reflection pairs become the symmetric outer/inner elements. A left pointer and a right pointer walking inward must always meet at partners: same `y`, and x-values summing to the fixed `sum = minX + maxX`. A lone centre point (`x == mirror`) works too because `i == j` gives `2x == sum`.

### Algorithm

1. Compute `sum = minX + maxX`.
2. Sort points by `y` ascending, then `x` ascending.
3. Set `i = 0`, `j = n-1`; while `i <= j`, require `points[i].y == points[j].y` and `points[i].x + points[j].x == sum`, then move `i++`, `j--`.
4. Any mismatch ⇒ `false`.

### Complexity

- **Time:** O(n log n) — the sort dominates; the pointer walk is linear.
- **Space:** O(1) auxiliary — sort in place.

### Code

```go
func twoPointers(points [][]int) bool {
	if len(points) == 0 {
		return true
	}
	minX, maxX := points[0][0], points[0][0]
	for _, p := range points {
		if p[0] < minX {
			minX = p[0]
		}
		if p[0] > maxX {
			maxX = p[0]
		}
	}
	sum := minX + maxX
	// Sort by y first so equal-y points cluster; then by x so a group is
	// laid out left-to-right and its reflection pairs are symmetric.
	sort.Slice(points, func(a, b int) bool {
		if points[a][1] != points[b][1] {
			return points[a][1] < points[b][1]
		}
		return points[a][0] < points[b][0]
	})
	i, j := 0, len(points)-1 // converge from both ends
	for i <= j {
		// Partners must share the same y AND have x-values summing to `sum`.
		if points[i][1] != points[j][1] || points[i][0]+points[j][0] != sum {
			return false
		}
		i++
		j--
	}
	return true
}
```

### Dry Run

Example 1: `points = [[1,1],[-1,1]]`, `sum = -1 + 1 = 0`.

| Step | Sorted points | i | j | Check | Action |
|------|---------------|---|---|-------|--------|
| 0 | `[(-1,1),(1,1)]` | 0 | 1 | y: 1==1 ✔, x: -1+1=0==sum ✔ | i→1, j→0 |
| 1 | — | 1 | 0 | i > j | exit loop |

Return **true** ✔

---

## Approach 3 — Hash Set (Optimal)

### Intuition

The brute force re-scans for every mirror. Instead, dump all points into a hash set keyed by `(x, y)`. The candidate line is still `sum = minX + maxX`. Then "mirror exists?" collapses to a single O(1) lookup of `(sum - x, y)`. One pass builds the set, one pass verifies — beating `O(n²)` as the follow-up asks.

### Algorithm

1. Insert every point into a set; track `minX`, `maxX`.
2. `sum = minX + maxX`.
3. For each `(x, y)`, check membership of `(sum - x, y)`. Any miss ⇒ `false`.

### Complexity

- **Time:** O(n) — two linear passes, O(1) set ops.
- **Space:** O(n) — the set stores up to n distinct points.

### Code

```go
func hashSet(points [][]int) bool {
	if len(points) == 0 {
		return true
	}
	type pt struct{ x, y int } // composite key for the set
	set := make(map[pt]bool, len(points))
	minX, maxX := points[0][0], points[0][0]
	for _, p := range points {
		set[pt{p[0], p[1]}] = true // remember this point exists
		if p[0] < minX {
			minX = p[0]
		}
		if p[0] > maxX {
			maxX = p[0]
		}
	}
	sum := minX + maxX
	// Each point's mirror must also be in the set — O(1) lookup.
	for _, p := range points {
		if !set[pt{sum - p[0], p[1]}] {
			return false
		}
	}
	return true
}
```

### Dry Run

Example 1: `points = [[1,1],[-1,1]]`.

| Step | Detail | State |
|------|--------|-------|
| 1 | Build set | `{(1,1),(-1,1)}`, minX=-1, maxX=1, sum=0 |
| 2 | p=(1,1): lookup (0-1, 1)=(-1,1) | present ✔ |
| 3 | p=(-1,1): lookup (0-(-1),1)=(1,1) | present ✔ |
| 4 | all mirrors present | return **true** ✔ |

---

## Key Takeaways

- **The mirror line is not searched — it is pinned** to `x = (minX + maxX)/2`. Comparing `2·mirror = minX + maxX` against `x1 + x2` keeps everything in integers and dodges floating-point.
- **Reflection over a vertical line preserves y** and maps `x → sum - x`. That single fact turns a geometry question into a membership check.
- **Hash set of composite keys** is the go-to for "does the symmetric/partner element exist" in O(1), the same trick behind Two Sum and Valid Anagram.
- Duplicates are handled automatically: a set makes them idempotent, and the two-pointer pairing still matches identical y-groups symmetrically.

---

## Related Problems

- LeetCode #1 — Two Sum (hash-set complement lookup)
- LeetCode #149 — Max Points on a Line (2D-point geometry with hashing)
- LeetCode #447 — Number of Boomerangs (point pairing via hash map)
- LeetCode #391 — Perfect Rectangle (geometric symmetry / coverage checks)
