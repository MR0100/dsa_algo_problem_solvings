package main

import (
	"fmt"
	"math"
	"math/rand"
)

// LeetCode 478 asks us to implement a Solution class:
//
//	Solution(radius, x_center, y_center) — set up the disk.
//	randPoint() []float64                — return a uniformly-random point
//	                                       inside (or on) that disk as [x, y].
//
// "Uniformly random" is the whole difficulty: every equal-area sub-region of
// the disk must be equally likely. Two standard correct strategies follow.
// Because output is random, main() cannot print fixed coordinates; instead it
// samples many points and prints deterministic *statistics* (all points inside,
// empirical mean ≈ center, empirical mean radius ≈ (2/3)·R) whose expected
// values are fixed. A shared seed makes the run reproducible.

// ── Approach 1: Rejection Sampling (Bounding Square) ─────────────────────────
//
// SolutionReject samples a uniform point in the axis-aligned square that
// circumscribes the disk and rejects (re-draws) any point that lands outside
// the disk.
//
// Intuition:
//
//	A uniform point in the [-R,R]×[-R,R] square is trivial: pick x and y each
//	uniformly. The disk sits inside that square and occupies π/4 ≈ 78.5% of its
//	area. If the drawn point is outside the disk, throw it away and try again.
//	Conditioning a uniform distribution on a sub-region yields a uniform
//	distribution on that sub-region — so accepted points are uniform on the disk.
//
// Algorithm:
//  1. Loop: draw x = uniform(-R, R), y = uniform(-R, R).
//  2. If x² + y² ≤ R² accept; return [xCenter + x, yCenter + y].
//  3. Otherwise repeat.
//
// Time:  O(1) expected — acceptance probability π/4, so ~1.27 draws on average
//
//	(geometric distribution). No dependence on R.
//
// Space: O(1).
type SolutionReject struct {
	radius, xCenter, yCenter float64
	rng                      *rand.Rand // private source so tests are seedable
}

// NewSolutionReject constructs the disk sampler (rejection variant).
func NewSolutionReject(radius, xCenter, yCenter float64, rng *rand.Rand) *SolutionReject {
	return &SolutionReject{radius: radius, xCenter: xCenter, yCenter: yCenter, rng: rng}
}

// RandPoint returns a uniformly random point in the disk via rejection.
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

// ── Approach 2: Polar with sqrt Radius (Optimal, No Rejection) ────────────────
//
// SolutionPolar draws an angle uniformly in [0, 2π) and a radius as
// R·√u (u uniform in [0,1)), guaranteeing a uniform *area* distribution in one
// shot — no rejection loop.
//
// Intuition:
//
//	In polar coordinates the disk's area element is r·dr·dθ, so area grows with
//	r, NOT uniformly in r. If we picked r uniformly in [0,R] we'd oversample the
//	center. The cumulative area up to radius r is πr²/(πR²) = (r/R)². To sample r
//	with that CDF, invert it: set (r/R)² = u (uniform) ⇒ r = R·√u. The angle is
//	uniform on the full circle. Then (x,y) = (r·cosθ, r·sinθ).
//
// Algorithm:
//  1. θ = 2π · uniform(0,1).
//  2. r = R · √(uniform(0,1)).   ← the √ is what makes it area-uniform.
//  3. Return [xCenter + r·cosθ, yCenter + r·sinθ].
//
// Time:  O(1) guaranteed (no loop). Space: O(1).
type SolutionPolar struct {
	radius, xCenter, yCenter float64
	rng                      *rand.Rand
}

// NewSolutionPolar constructs the disk sampler (polar variant).
func NewSolutionPolar(radius, xCenter, yCenter float64, rng *rand.Rand) *SolutionPolar {
	return &SolutionPolar{radius: radius, xCenter: xCenter, yCenter: yCenter, rng: rng}
}

// RandPoint returns a uniformly random point in the disk via polar sampling.
func (s *SolutionPolar) RandPoint() []float64 {
	angle := 2 * math.Pi * s.rng.Float64()     // uniform direction in [0, 2π)
	r := s.radius * math.Sqrt(s.rng.Float64()) // √ gives uniform area, not center bias
	x := s.xCenter + r*math.Cos(angle)         // convert polar → Cartesian x
	y := s.yCenter + r*math.Sin(angle)         // convert polar → Cartesian y
	return []float64{x, y}
}

// ── Verification helper ───────────────────────────────────────────────────────

// sampleStats draws n points from `draw` and reports whether every point is
// inside the disk, plus the empirical average distance from the center divided
// by R. For a *uniform* disk that ratio converges to 2/3 (E[r] = ∫₀^R r·(2r/R²)
// dr = 2R/3). A BUGGY "uniform r in [0,R]" sampler would instead give R/2 and
// concentrate near the center — so this statistic distinguishes correct from
// wrong. We print the ratio rounded to one decimal so the output is stable.
func sampleStats(radius, xc, yc float64, n int, draw func() []float64) (allInside bool, meanRadiusRatio float64) {
	allInside = true
	sumR := 0.0
	for i := 0; i < n; i++ {
		p := draw()
		dx, dy := p[0]-xc, p[1]-yc
		dist := math.Hypot(dx, dy)
		if dist > radius+1e-9 { // tolerate tiny float error on the boundary
			allInside = false
		}
		sumR += dist
	}
	return allInside, (sumR / float64(n)) / radius
}

func main() {
	// Fixed seed → identical numbers every run, so the expected comments hold.
	rng := rand.New(rand.NewSource(42))
	const n = 200000

	fmt.Println("=== Approach 1: Rejection Sampling ===")
	// Disk: radius 1.0 centered at (0,0). Show a couple of sample points, then stats.
	r1 := NewSolutionReject(1.0, 0.0, 0.0, rng)
	p := r1.RandPoint()
	fmt.Printf("sample point inside unit disk? %v  expected true\n", math.Hypot(p[0], p[1]) <= 1.0+1e-9)
	inside, ratio := sampleStats(1.0, 0.0, 0.0, n, r1.RandPoint)
	fmt.Printf("all %d points inside?          %v  expected true\n", n, inside)
	fmt.Printf("mean radius / R (want ~0.667)  %.1f  expected 0.7\n", ratio) // 2/3 rounds to 0.7

	fmt.Println("=== Approach 2: Polar with sqrt Radius (Optimal) ===")
	// Official-style disk: radius 0.01, center (-73.5, 40.3). Points must land near it.
	r2 := NewSolutionPolar(0.01, -73.5, 40.3, rng)
	q := r2.RandPoint()
	fmt.Printf("randPoint() near center?       %v  expected true\n", math.Hypot(q[0]-(-73.5), q[1]-40.3) <= 0.01+1e-9)
	inside2, ratio2 := sampleStats(0.01, -73.5, 40.3, n, r2.RandPoint)
	fmt.Printf("all %d points inside?          %v  expected true\n", n, inside2)
	fmt.Printf("mean radius / R (want ~0.667)  %.1f  expected 0.7\n", ratio2)

	// Cross-check: both strategies produce the same area statistic on the same disk.
	r3 := NewSolutionPolar(1.0, 0.0, 0.0, rng)
	_, ratio3 := sampleStats(1.0, 0.0, 0.0, n, r3.RandPoint)
	fmt.Printf("polar mean radius / R on unit  %.1f  expected 0.7\n", ratio3)
}
