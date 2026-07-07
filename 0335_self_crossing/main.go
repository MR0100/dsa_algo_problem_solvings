package main

import "fmt"

// ── Approach 1: Geometric Case Analysis (Optimal, O(n) time / O(1) space) ─────
//
// selfCrossingCases decides whether the spiral path crosses itself by checking
// three local geometric patterns against only the previous few edges.
//
// Intuition:
//
//	The moves go N, W, S, E, N, ... turning counter-clockwise each time. A new
//	edge can only ever intersect one of the recent edges — you cannot loop back
//	to touch an edge many steps ago without first crossing a nearer one. So a
//	crossing must fall into exactly one of three local cases:
//
//	  Case 1 (crosses 4th line back): the current edge d[i] reaches or passes
//	    the edge two before it. Condition: d[i] >= d[i-2] AND d[i-1] <= d[i-3].
//
//	  Case 2 (touches 5th line back): the spiral just closes onto the edge
//	    d[i-4]. Condition: i>=4, d[i-1] == d[i-3] AND d[i] + d[i-4] >= d[i-2].
//
//	  Case 3 (crosses 6th line back): a wider overlap.
//	    Condition: i>=5, d[i-2] >= d[i-4] AND d[i-3] >= d[i-1] AND
//	               d[i-1] + d[i-5] >= d[i-3] AND d[i] + d[i-4] >= d[i-2].
//
//	If none of these hold for any i, the path never crosses.
//
// Algorithm:
//  1. For each i from 3..n-1, test Case 1; from 4, add Case 2; from 5, add Case 3.
//  2. Return true on the first hit; false if the scan finishes clean.
//
// Time:  O(n) — a single pass with O(1) checks per index.
// Space: O(1) — only index arithmetic.
func selfCrossingCases(distance []int) bool {
	d := distance
	n := len(d)
	for i := 3; i < n; i++ {
		// Case 1: current edge crosses the edge 2 steps back (4th line).
		//   d[i] catches up to d[i-2] while shrinking inward (d[i-1] <= d[i-3]).
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

// ── Approach 2: Brute Force Segment Intersection ─────────────────────────────
//
// bruteForceSegments builds the actual line segments and tests each new
// segment against all earlier non-adjacent segments for intersection.
//
// Intuition:
//
//	Forget the geometry shortcuts: literally walk the path, record every
//	segment as (x1,y1)-(x2,y2), and for each new segment check it against all
//	previous segments (skipping the immediately adjacent one, which always
//	shares an endpoint). This is the definition of "path crosses itself".
//	O(n^2) but a rock-solid oracle to validate the clever O(n) rules.
//
// Algorithm:
//  1. Simulate positions in the fixed N,W,S,E cycle, producing n segments.
//  2. For each segment i, test intersection with segments 0..i-2.
//  3. Return true on the first intersection.
//
// Time:  O(n^2) — each of n segments compared with up to n earlier ones.
// Space: O(n) — stores all segments.
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

// seg is one axis-aligned path segment from (x1,y1) to (x2,y2).
type seg struct{ x1, y1, x2, y2 int }

// segmentsIntersect reports whether two axis-aligned segments intersect
// (touch or cross). Each segment is horizontal or vertical.
func segmentsIntersect(a, b seg) bool {
	// Normalize each segment to [lo,hi] ranges on x and y.
	aMinX, aMaxX := minI(a.x1, a.x2), maxI(a.x1, a.x2)
	aMinY, aMaxY := minI(a.y1, a.y2), maxI(a.y1, a.y2)
	bMinX, bMaxX := minI(b.x1, b.x2), maxI(b.x1, b.x2)
	bMinY, bMaxY := minI(b.y1, b.y2), maxI(b.y1, b.y2)
	// Two axis-aligned segments intersect iff their x-ranges and y-ranges overlap.
	return aMinX <= bMaxX && bMinX <= aMaxX &&
		aMinY <= bMaxY && bMinY <= aMaxY
}

func minI(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxI(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	ex1 := []int{2, 1, 1, 2}
	ex2 := []int{1, 2, 3, 4}
	ex3 := []int{1, 1, 1, 2, 1}

	fmt.Println("=== Approach 1: Geometric Case Analysis (Optimal) ===")
	fmt.Println(selfCrossingCases(ex1)) // expected true
	fmt.Println(selfCrossingCases(ex2)) // expected false
	fmt.Println(selfCrossingCases(ex3)) // expected true

	fmt.Println("=== Approach 2: Brute Force Segment Intersection ===")
	fmt.Println(bruteForceSegments(ex1)) // expected true
	fmt.Println(bruteForceSegments(ex2)) // expected false
	fmt.Println(bruteForceSegments(ex3)) // expected true
}
