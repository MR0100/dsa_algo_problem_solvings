# 0447 — Number of Boomerangs

> LeetCode #447 · Difficulty: Medium
> **Categories:** Array, Hash Table, Math, Geometry

---

## Problem Statement

You are given `n` `points` in the plane that are all **distinct**, where `points[i] = [xi, yi]`. A **boomerang** is a tuple of points `(i, j, k)` such that the distance between `i` and `j` equals the distance between `i` and `k` **(the order of the tuple matters)**.

Return *the number of boomerangs*.

**Example 1:**

```
Input: points = [[0,0],[1,0],[2,0]]
Output: 2
Explanation: The two boomerangs are [[1,0],[0,0],[2,0]] and [[1,0],[2,0],[0,0]].
```

**Example 2:**

```
Input: points = [[1,1],[2,2],[3,3]]
Output: 2
```

**Example 3:**

```
Input: points = [[1,1]]
Output: 0
```

**Constraints:**

- `n == points.length`
- `1 <= n <= 500`
- `points[i].length == 2`
- `-10^4 <= xi, yi <= 10^4`
- All the points are **unique**.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Facebook   | ★★☆☆☆ Low        | 2022          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Apple      | ★★☆☆☆ Low        | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Map** — for a fixed apex, bucket the other points by their (squared) distance so equidistant points can be counted in one pass instead of a nested scan → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Geometry (distance)** — a boomerang is defined by equal Euclidean distances; comparing *squared* distances keeps the arithmetic integer and exact → see [`/dsa/geometry.md`](/dsa/geometry.md)
- **Counting / Permutations (Math)** — a distance class of `m` points yields `m·(m−1)` ordered pairs, the number of boomerangs through that apex at that radius → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Triple Loop) | O(n³) | O(1) | Baseline; n ≤ 500 → 1.25×10⁸ triples, borderline but conceptually clear |
| 2 | Hash Map by Center (Optimal) | O(n²) | O(n) | Intended solution; group distances per apex and use m·(m−1) |

---

## Approach 1 — Brute Force (Triple Loop)

### Intuition

The definition is directly testable. A boomerang is an ordered triple `(i, j, k)` where `i` is equidistant from `j` and `k`, and `j ≠ k`. So fix the apex `i`, then try every ordered pair `(j, k)` of *other* points and count those where the two distances from `i` are equal. Because the tuple order matters, `(i, j, k)` and `(i, k, j)` count as two. Compare **squared** distances — `dx² + dy²` — to stay in integer arithmetic and avoid `sqrt` rounding errors.

### Algorithm

1. For each center `i` (0..n−1):
   - For each `j ≠ i`:
     - For each `k ≠ i` and `k ≠ j`:
       - If `squaredDist(i, j) == squaredDist(i, k)` → `count++`.
2. Return `count`.

### Complexity

- **Time:** O(n³) — three nested loops over the `n` points.
- **Space:** O(1) — just the running counter.

### Code

```go
func bruteForce(points [][]int) int {
	n := len(points)
	count := 0
	// sq returns the squared Euclidean distance between points a and b; using
	// the square keeps everything integer and dodges sqrt rounding issues.
	sq := func(a, b []int) int {
		dx := a[0] - b[0]
		dy := a[1] - b[1]
		return dx*dx + dy*dy
	}
	for i := 0; i < n; i++ { // i is the boomerang's "apex" (equidistant point)
		for j := 0; j < n; j++ {
			if j == i {
				continue // j must differ from the center
			}
			for k := 0; k < n; k++ {
				if k == i || k == j {
					continue // k must differ from both center and j
				}
				// Ordered triple (i,j,k) is a boomerang iff |ij| == |ik|.
				if sq(points[i], points[j]) == sq(points[i], points[k]) {
					count++
				}
			}
		}
	}
	return count
}
```

### Dry Run

Example 1: `points = [[0,0],[1,0],[2,0]]` (index 0=A, 1=B, 2=C). Squared distances: AB=1, AC=4, BC=1.

| center i | (j, k) tried | sq(i,j) | sq(i,k) | equal? | count |
|----------|--------------|---------|---------|--------|-------|
| 0 (A) | (B, C) | 1 | 4 | no | 0 |
| 0 (A) | (C, B) | 4 | 1 | no | 0 |
| 1 (B) | (A, C) | 1 | 1 | **yes** | 1 |
| 1 (B) | (C, A) | 1 | 1 | **yes** | 2 |
| 2 (C) | (A, B) | 4 | 1 | no | 2 |
| 2 (C) | (B, A) | 1 | 4 | no | 2 |

Only apex B has two points at equal distance. Result: **2**. ✔

---

## Approach 2 — Hash Map by Center (Optimal)

### Intuition

Fix the apex `i`. Suppose `m` of the other points lie at the *same* distance from `i`. Every **ordered** pair drawn from those `m` points is a valid boomerang `(i, j, k)`: pick `j` (`m` ways), then a different `k` (`m−1` ways) → `m·(m−1)` boomerangs. So we do not need the inner pair loop at all: make one pass over the other points, bucket them by squared distance in a hash map, and sum `m·(m−1)` across the buckets. Repeat for every apex.

### Algorithm

1. For each center `i`:
   - Build a map `buckets: squaredDistance → count`, iterating all `j ≠ i`.
   - For each bucket holding `m` points, add `m·(m−1)` to the answer.
2. Return the total.

### Complexity

- **Time:** O(n²) — for each of the `n` apexes, a single O(n) pass buckets the distances; the per-bucket summation is bounded by the same pass.
- **Space:** O(n) — one distance map per apex, at most `n−1` entries, reused each iteration.

### Code

```go
func hashMap(points [][]int) int {
	n := len(points)
	count := 0
	for i := 0; i < n; i++ { // choose the apex
		// buckets: squared distance from i -> number of points at that distance
		buckets := make(map[int]int)
		for j := 0; j < n; j++ {
			if j == i {
				continue // skip the apex itself
			}
			dx := points[i][0] - points[j][0]
			dy := points[i][1] - points[j][1]
			buckets[dx*dx+dy*dy]++ // tally this distance class
		}
		// For each distance class of size m, there are m*(m-1) ordered (j,k)
		// pairs, each a distinct boomerang with apex i.
		for _, m := range buckets {
			count += m * (m - 1)
		}
	}
	return count
}
```

### Dry Run

Example 1: `points = [[0,0],[1,0],[2,0]]` (A, B, C).

| apex i | buckets built (dist² → count) | contribution Σ m·(m−1) | running count |
|--------|-------------------------------|------------------------|----------------|
| 0 (A) | {1: 1 (B), 4: 1 (C)} | 1·0 + 1·0 = 0 | 0 |
| 1 (B) | {1: 2 (A, C)} | 2·1 = **2** | 2 |
| 2 (C) | {4: 1 (A), 1: 1 (B)} | 1·0 + 1·0 = 0 | 2 |

Apex B has a bucket of size 2 at distance 1 → `2·1 = 2` boomerangs. Result: **2**. ✔

---

## Key Takeaways

- **"Equal distance / equal something" ⇒ bucket by that key.** Grouping equidistant points in a hash map turns an O(n³) pair search into an O(n²) tally.
- **Ordered pairs from a group of `m`: `m·(m−1)`.** (Unordered would be `m·(m−1)/2`.) The problem's "order matters" is exactly why there is no `/2` here.
- **Compare squared distances, never `sqrt`.** With integer coordinates the squared distance is exact; `sqrt` introduces floating-point ties/mismatches.
- The apex is the *first* element of the tuple — the point that must be equidistant. Always identify which role is fixed before counting.

---

## Related Problems

- LeetCode #149 — Max Points on a Line (group points by slope through a fixed point)
- LeetCode #1 — Two Sum (canonical hash-map counting)
- LeetCode #454 — 4Sum II (count pairs via hash map of sums)
- LeetCode #356 — Line Reflection (geometry + hashing coordinates)
