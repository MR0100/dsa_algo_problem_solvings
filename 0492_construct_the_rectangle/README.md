# 0492 — Construct the Rectangle

> LeetCode #492 · Difficulty: Easy
> **Categories:** Math

---

## Problem Statement

A web developer needs to know how to design a web page's size. So, given a specific rectangular web page's area, your job by now is to design a rectangular web page, whose length `L` and width `W` satisfy the following requirements:

1. The area of the rectangular web page you designed must equal to the given target area.
2. The width `W` should not be larger than the length `L`, which means `L >= W`.
3. The difference between length `L` and width `W` should be as small as possible.

Return *an array `[L, W]` where `L` and `W` are the length and width of the web page you designed in sequence.*

**Example 1:**

```
Input: area = 4
Output: [2,2]
Explanation: The target area is 4, and all the possible ways to construct it are [1,4], [2,2], [4,1].
But according to requirement 2, [1,4] is illegal; according to requirement 3, [4,1] is not optimal compared to [2,2]. So the length L is 2, and the width W is 2.
```

**Example 2:**

```
Input: area = 37
Output: [37,1]
```

**Example 3:**

```
Input: area = 122122
Output: [427,286]
```

**Constraints:**

- `1 <= area <= 10^7`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★☆☆☆ Low        | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Google     | ★☆☆☆☆ Rare       | 2022          |
| Apple      | ★☆☆☆☆ Rare       | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Factorization near the square root** — the pair `(L, W)` minimizing `L − W` is the divisor pair closest to `√area`; `W` is the largest divisor `≤ √area` and `L = area / W`. This is the same "iterate to √n" idea used for divisor/prime checks → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force Downward Scan | O(area) | O(1) | Conceptual baseline; TLE risk for area up to 1e7 |
| 2 | Start at √area, Walk Down (Optimal) | O(√area) | O(1) | The intended solution; only ~√area steps |

---

## Approach 1 — Brute Force Downward Scan

### Intuition

We want the width closest to `√area` (that minimizes `L − W`). The most literal correct method scans every integer `w` and keeps the largest one that (a) divides `area` and (b) does not exceed its partner `area/w` (so `W ≤ L`). The largest such `w` is exactly the width nearest `√area` from below. Correct, but it touches every integer up to `area`.

### Algorithm

1. Initialize `best = 1` (width `1` always divides `area`).
2. For `w` from `1` to `area`: if `area % w == 0` **and** `w <= area/w`, set `best = w`.
3. Return `[area/best, best]`.

### Complexity

- **Time:** O(area) — inspects every integer up to `area`; for `area = 10⁷` that is ten million iterations (borderline / TLE under tight limits).
- **Space:** O(1) — a couple of scalars.

### Code

```go
func bruteForce(area int) []int {
	best := 1 // width 1 is always a valid divisor
	for w := 1; w <= area; w++ {
		if area%w != 0 {
			continue // w must divide area to form an integer rectangle
		}
		l := area / w
		if w <= l { // keep W <= L; the largest such w gives the smallest L-W
			best = w // remember the widest width seen so far that stays <= its length
		}
	}
	return []int{area / best, best} // [L, W]
}
```

### Dry Run

Example 1: `area = 4`.

| w | area % w | l = area/w | w <= l? | best after |
|---|----------|-----------|---------|-----------|
| 1 | 0 | 4 | yes | 1 |
| 2 | 0 | 2 | yes | 2 |
| 3 | 1 | — | (not a divisor) | 2 |
| 4 | 0 | 1 | no (4 > 1) | 2 |

Loop ends. `best = 2`, so return `[4/2, 2] = [2, 2]` ✔

---

## Approach 2 — Start at √area and Walk Down (Optimal)

### Intuition

For a fixed product `area`, the difference `L − W` is smallest when `W` is as close as possible to `√area`. So don't scan from the bottom — **start** `W` at `floor(√area)` and decrease it one step at a time. The first `W` that divides `area` is the largest divisor `≤ √area`; its partner `L = area / W` is then the smallest divisor `≥ √area`, giving the tightest `L − W`. Because divisors are symmetric about `√area`, we never need to look above it.

### Algorithm

1. Set `w = floor(√area)`.
2. While `area % w != 0`, decrement `w`.
3. Return `[area/w, w]`.

### Complexity

- **Time:** O(√area) — at most `√area` decrements before landing on a divisor (`w = 1` always divides, so the loop is guaranteed to terminate).
- **Space:** O(1).

### Code

```go
func sqrtScan(area int) []int {
	w := int(math.Sqrt(float64(area))) // largest candidate width to try first
	for area%w != 0 {                  // walk down until w divides area exactly
		w-- // the first divisor found is the largest divisor <= sqrt(area)
	}
	return []int{area / w, w} // L = area/w (>= w), W = w
}
```

### Dry Run

Example 1: `area = 4`.

| Step | w | area % w | Action |
|------|---|----------|--------|
| init | `floor(√4) = 2` | 0 | `4 % 2 == 0` → divisor found, exit loop |

Return `[4/2, 2] = [2, 2]` ✔

Extra trace, Example 3: `area = 122122`. `floor(√122122) = 349`. Walk down: `349,348,…` none divide until `286` (`122122 / 286 = 427`, `286 · 427 = 122122`). Return `[427, 286]` ✔

---

## Key Takeaways

- **Closest-to-square divisor pair:** to split `n` into two factors with minimal difference, start at `floor(√n)` and walk down to the first divisor. This is a reusable micro-pattern (also: "is `n` prime?" only needs divisors up to `√n`).
- Casting through `float64` for `math.Sqrt` can be off by one for perfect squares near the floating-point boundary; walking *down* from `floor(√area)` is self-correcting because we only ever accept an exact divisor.
- Divisors pair up symmetrically around `√n`: if `w · l = n` and `w ≤ √n`, then `l ≥ √n`. That symmetry is why searching only up to `√n` suffices.
- The width `W` is returned second in `[L, W]` even though it is the value we actually search for.

---

## Related Problems

- LeetCode #593 — Valid Square (geometry from side/area reasoning)
- LeetCode #829 — Consecutive Numbers Sum (divisor enumeration up to √-bound)
- LeetCode #1492 — The kth Factor of n (iterate divisors)
- LeetCode #204 — Count Primes (trial division up to √n)
