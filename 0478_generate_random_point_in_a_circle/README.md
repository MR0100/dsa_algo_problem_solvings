# 0478 — Generate Random Point in a Circle

> LeetCode #478 · Difficulty: Medium
> **Categories:** Math, Geometry, Randomized, Rejection Sampling

---

## Problem Statement

Given the radius and the position of the center of a circle, implement the function `randPoint` which generates a uniform random point inside the circle.

Implement the `Solution` class:

- `Solution(double radius, double x_center, double y_center)` initializes the object with the radius of the circle `radius` and the position of the center `(x_center, y_center)`.
- `randPoint()` returns a random point inside the circle. A point on the circumference of the circle is considered to be in the circle. The answer is returned as an array `[x, y]`.

**Example 1:**

```
Input
["Solution", "randPoint", "randPoint", "randPoint"]
[[1.0, 0.0, 0.0], [], [], []]
Output
[null, [-0.02493, -0.38077], [0.82314, 0.38945], [0.36572, 0.17248]]

Explanation
Solution solution = new Solution(1.0, 0.0, 0.0);
solution.randPoint();  // return [-0.02493, -0.38077]
solution.randPoint();  // return [0.82314, 0.38945]
solution.randPoint();  // return [0.36572, 0.17248]
```

**Constraints:**

- `0 < radius <= 10^8`
- `-10^7 <= x_center, y_center <= 10^7`
- At most `3 * 10^4` calls will be made to `randPoint`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Facebook   | ★★★☆☆ Medium     | 2023          |
| Uber       | ★★★☆☆ Medium     | 2022          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Geometry / Uniform Sampling in 2-D** — the crux is that area in polar form is `r·dr·dθ`, so uniform *radius* is not uniform *area*; you must transform with `r = R·√u` → see [`/dsa/geometry.md`](/dsa/geometry.md)
- **Randomized Algorithms / Rejection Sampling** — the bounding-square method conditions a uniform square draw on the disk, a canonical rejection-sampling technique → see [`/dsa/reservoir_sampling.md`](/dsa/reservoir_sampling.md)
- **Math (inverse-transform sampling)** — invert the area CDF `(r/R)² = u` to draw the radius; a core probability tool → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)

---

## Approaches Overview

| # | Approach | Time / call | Space | When to use |
|---|----------|-------------|-------|-------------|
| 1 | Rejection Sampling (bounding square) | O(1) expected (~1.27 draws) | O(1) | Easiest to reason about correctness; no trig |
| 2 | Polar with `√` Radius (Optimal) | O(1) worst-case (no loop) | O(1) | Constant work always; the clean interview answer |

---

## Approach 1 — Rejection Sampling (Bounding Square)

### Intuition

Sampling a uniform point in the square `[-R, R] × [-R, R]` is trivial — draw `x` and `y` each uniformly. The disk is inscribed in that square and fills `π/4 ≈ 78.5%` of it. Draw a square point; if it lands inside the disk keep it, otherwise discard and redraw. Conditioning a uniform distribution on a sub-region gives a uniform distribution on that sub-region, so accepted points are uniform on the disk. Expected draws per point are `1 / (π/4) ≈ 1.27`, independent of `R`.

### Algorithm

1. Loop:
   - `x = (2·rand() − 1)·R`, `y = (2·rand() − 1)·R` (uniform in the square).
   - If `x² + y² ≤ R²`, accept: return `[x_center + x, y_center + y]`.
   - Otherwise repeat.

### Complexity

- **Time:** O(1) expected — geometric number of trials with success probability `π/4`.
- **Space:** O(1).

### Code

```go
func (s *SolutionReject) RandPoint() []float64 {
	for { // keep drawing until a point lands inside the disk
		// Map [0,1) → [-R, R): 2*u-1 ∈ [-1,1), times R.
		x := (2*s.rng.Float64() - 1) * s.radius // uniform x offset in the square
		y := (2*s.rng.Float64() - 1) * s.radius // uniform y offset in the square
		if x*x+y*y <= s.radius*s.radius {       // inside (or on) the circle?
			return []float64{s.xCenter + x, s.yCenter + y} // translate to true center
		}
		// else: point was in a corner of the square — reject and loop
	}
}
```

### Dry Run

Disk `R = 1`, center `(0, 0)`. Suppose the RNG yields these `Float64()` pairs:

| Trial | u₁, u₂ | x = 2u₁−1 | y = 2u₂−1 | x² + y² | ≤ 1? | Action |
|-------|--------|-----------|-----------|---------|------|--------|
| 1 | 0.95, 0.95 | 0.90 | 0.90 | 1.62 | no | reject (corner) |
| 2 | 0.60, 0.30 | 0.20 | −0.40 | 0.20 | yes | accept → `[0.20, −0.40]` |

Returned point `[0.20, −0.40]` lies inside the unit disk. ✔

---

## Approach 2 — Polar with `√` Radius (Optimal)

### Intuition

Go polar: pick an angle `θ` uniformly on `[0, 2π)` and a radius `r`. The subtlety is the radius. A disk's area element is `r·dr·dθ` — area grows *linearly with r*, so picking `r` uniformly on `[0, R]` would oversample the center (too many points near `r = 0`). The fraction of area within radius `r` is `πr² / πR² = (r/R)²`. This is the CDF of `r`; to sample from it, apply **inverse-transform sampling**: set `(r/R)² = u` for a uniform `u`, giving `r = R·√u`. The `√` is exactly what compensates for the center bias. Then convert to Cartesian.

### Algorithm

1. `θ = 2π · rand()`.
2. `r = R · √(rand())`.
3. Return `[x_center + r·cos θ, y_center + r·sin θ]`.

### Complexity

- **Time:** O(1) worst case — one angle draw, one radius draw, two trig calls; no rejection loop.
- **Space:** O(1).

### Code

```go
func (s *SolutionPolar) RandPoint() []float64 {
	angle := 2 * math.Pi * s.rng.Float64()     // uniform direction in [0, 2π)
	r := s.radius * math.Sqrt(s.rng.Float64()) // √ gives uniform area, not center bias
	x := s.xCenter + r*math.Cos(angle)         // convert polar → Cartesian x
	y := s.yCenter + r*math.Sin(angle)         // convert polar → Cartesian y
	return []float64{x, y}
}
```

### Dry Run

Disk `R = 1`, center `(0, 0)`. Suppose the RNG yields `Float64()` values `u_angle = 0.25`, then `u_r = 0.25`:

| Step | Expression | Value |
|------|------------|-------|
| 1 | `θ = 2π · 0.25` | `π/2` (90°) |
| 2 | `r = 1 · √0.25` | `0.5` |
| 3 | `x = 0 + 0.5·cos(π/2)` | `0.5 · 0 = 0.0` |
| 4 | `y = 0 + 0.5·sin(π/2)` | `0.5 · 1 = 0.5` |

Returned `[0.0, 0.5]`: distance `0.5 ≤ 1`, inside the disk. Note that because `r = √u`, `u = 0.25` maps to `r = 0.5` (not `0.25`), pushing mass outward to keep the *area* uniform. ✔

---

## Key Takeaways

- **Uniform radius ≠ uniform area.** In 2-D, sample `r = R·√u`; in 3-D (ball) it's `r = R·u^(1/3)`. The exponent comes from inverting the volume CDF.
- **Inverse-transform sampling:** to draw from any distribution, set its CDF equal to a uniform `u` and solve for the variable. Here CDF `= (r/R)²`.
- **Rejection sampling** is the fallback when the target region sits inside an easy-to-sample region: draw from the easy region, keep only what lands in the target. Efficiency = area ratio (`π/4` for disk-in-square).
- A neat sanity statistic: for a uniform disk, `E[r] = 2R/3` (≈ 0.667·R). A wrong "uniform-r" sampler gives `R/2` and visibly clumps near the center — this is exactly how `main()` here validates correctness deterministically despite randomness.

---

## Related Problems

- LeetCode #470 — Implement Rand10() Using Rand7() (rejection sampling)
- LeetCode #497 — Random Point in Non-overlapping Rectangles (weighted region + uniform sampling)
- LeetCode #528 — Random Pick with Weight (inverse-transform via prefix sums + binary search)
- LeetCode #519 — Random Flip Matrix (uniform sampling without replacement)
- LeetCode #398 — Random Pick Index (reservoir sampling)
