package main

import (
	"fmt"
	"math"
)

// ── Approach 1: Brute Force (Check Every Pair's Line) ────────────────────────
//
// bruteForce solves Max Points on a Line by fixing every pair of points and
// counting how many points lie on the line through that pair.
//
// Intuition:
//
//	Every candidate line is determined by (at least) two of the given points.
//	So enumerate all O(n^2) pairs, and for each pair count the points that are
//	collinear with it. Collinearity is tested with the CROSS PRODUCT
//	(x2-x1)*(y3-y1) - (y2-y1)*(x3-x1) == 0 — pure integer math, no division,
//	no precision issues, and vertical lines need no special case.
//
// Algorithm:
//  1. If n <= 2, every point set is trivially on one line → return n.
//  2. For each pair (i, j), i < j:
//     a. Count every k with cross(points[i], points[j], points[k]) == 0.
//     b. Track the maximum count.
//  3. Return the maximum.
//
// Time:  O(n^3) — n^2 pairs × n collinearity checks.
// Space: O(1) — only counters.
func bruteForce(points [][]int) int {
	n := len(points)
	if n <= 2 {
		return n // 1 or 2 points always share a line
	}
	best := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			count := 0 // points on the line through points[i] and points[j]
			for k := 0; k < n; k++ {
				// cross product of vectors (i→j) and (i→k); zero ⇔ collinear
				cross := (points[j][0]-points[i][0])*(points[k][1]-points[i][1]) -
					(points[j][1]-points[i][1])*(points[k][0]-points[i][0])
				if cross == 0 {
					count++ // k lies on the line (i and j count themselves)
				}
			}
			if count > best {
				best = count
			}
		}
	}
	return best
}

// ── Approach 2: Hash Map of Float Slopes per Anchor ──────────────────────────
//
// hashMapFloatSlopes solves Max Points on a Line by anchoring each point and
// bucketing every other point by the float64 slope of the connecting line.
//
// Intuition:
//
//	All lines through a fixed anchor point differ only by slope. So for each
//	anchor, group the other n-1 points by slope in a hash map; the biggest
//	bucket + 1 (the anchor itself) is the best line through that anchor.
//	Checking every anchor covers every line. float64 slopes are exact enough
//	here because |coordinates| <= 10^4 keeps dy/dx well within double
//	precision — but beware: floats are NOT safe for unbounded inputs.
//
// Algorithm:
//  1. If n <= 2 → return n.
//  2. For each anchor i: build map[slope]count over all j != i.
//     - dx == 0 → slope = +Inf (vertical line).
//     - dy == 0 → slope = 0 (avoids the -0.0 vs +0.0 key split).
//  3. best = max over anchors of (largest bucket + 1).
//
// Time:  O(n^2) — n anchors × n slope insertions (O(1) each on average).
// Space: O(n) — one slope map at a time.
func hashMapFloatSlopes(points [][]int) int {
	n := len(points)
	if n <= 2 {
		return n // trivially one line
	}
	best := 0
	for i := 0; i < n; i++ {
		slopes := make(map[float64]int) // slope → how many points share it with anchor i
		for j := 0; j < n; j++ {
			if j == i {
				continue // don't pair the anchor with itself
			}
			dx := float64(points[j][0] - points[i][0])
			dy := float64(points[j][1] - points[i][1])
			var slope float64
			switch {
			case dx == 0:
				slope = math.Inf(1) // vertical line: one shared bucket
			case dy == 0:
				slope = 0 // force +0 so -0.0 and +0.0 don't split the bucket
			default:
				slope = dy / dx
			}
			slopes[slope]++
			if slopes[slope]+1 > best { // +1 counts the anchor itself
				best = slopes[slope] + 1
			}
		}
	}
	return best
}

// ── Approach 3: Hash Map of GCD-Normalized Slopes (Optimal) ──────────────────
//
// hashMapGCDSlopes solves Max Points on a Line like Approach 2 but keys each
// bucket by the slope as an exact reduced fraction (dy/g, dx/g).
//
// Intuition:
//
//	Same anchor-and-bucket idea, but the slope is stored as an integer pair
//	instead of a float: divide (dy, dx) by gcd(|dy|, |dx|) and normalize the
//	sign so that equal slopes always produce the identical key. This is fully
//	exact for any integer coordinates — the robust interview answer.
//
// Algorithm:
//  1. If n <= 2 → return n.
//  2. For each anchor i, for each j != i:
//     a. dy, dx = differences; g = gcd(|dy|, |dx|); reduce both by g.
//     b. Sign-normalize: if dx < 0, or dx == 0 and dy < 0, negate both.
//     (So slope 1/2 and -1/-2 collapse, and vertical is always (1, 0).)
//     c. Increment bucket [dy, dx]; track max bucket + 1.
//  3. Return best.
//
// Time:  O(n^2 log C) — n^2 pairs, each with a gcd on values up to C = 2*10^4.
// Space: O(n) — one bucket map per anchor at a time.
func hashMapGCDSlopes(points [][]int) int {
	n := len(points)
	if n <= 2 {
		return n // trivially one line
	}
	best := 0
	for i := 0; i < n; i++ {
		slopes := make(map[[2]int]int) // reduced (dy, dx) → count of points
		for j := 0; j < n; j++ {
			if j == i {
				continue // skip the anchor itself
			}
			dy := points[j][1] - points[i][1]
			dx := points[j][0] - points[i][0]
			g := gcd(abs(dy), abs(dx)) // g > 0 because points are distinct
			dy, dx = dy/g, dx/g        // reduce the fraction to lowest terms
			if dx < 0 || (dx == 0 && dy < 0) {
				dy, dx = -dy, -dx // canonical sign: dx > 0, or vertical (dy=1, dx=0)
			}
			key := [2]int{dy, dx}
			slopes[key]++
			if slopes[key]+1 > best { // +1 for the anchor
				best = slopes[key] + 1
			}
		}
	}
	return best
}

// gcd returns the greatest common divisor via Euclid's algorithm.
func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b // replace (a,b) with (b, a mod b) until b is 0
	}
	return a
}

// abs returns |x| for ints.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	// Official LeetCode examples.
	example1 := [][]int{{1, 1}, {2, 2}, {3, 3}}
	example2 := [][]int{{1, 1}, {3, 2}, {5, 3}, {4, 1}, {2, 3}, {1, 4}}

	fmt.Println("=== Approach 1: Brute Force (Check Every Pair's Line) ===")
	fmt.Println(bruteForce(example1)) // 3
	fmt.Println(bruteForce(example2)) // 4

	fmt.Println("=== Approach 2: Hash Map of Float Slopes per Anchor ===")
	fmt.Println(hashMapFloatSlopes(example1)) // 3
	fmt.Println(hashMapFloatSlopes(example2)) // 4

	fmt.Println("=== Approach 3: Hash Map of GCD-Normalized Slopes (Optimal) ===")
	fmt.Println(hashMapGCDSlopes(example1)) // 3
	fmt.Println(hashMapGCDSlopes(example2)) // 4
}
