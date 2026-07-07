package main

import "fmt"

// point is a 2-D lattice point / vertex of the polygon.
type point struct{ x, y int }

// cross returns the 2-D cross product of the vectors (b-a) and (c-b).
//
//	> 0  → the turn a→b→c is counter-clockwise (left turn)
//	< 0  → the turn is clockwise (right turn)
//	= 0  → a, b, c are collinear (straight)
//
// Using integer coordinates keeps this exact — no floating point, no division.
func cross(a, b, c point) int {
	// (b-a) = (b.x-a.x, b.y-a.y);  (c-b) = (c.x-b.x, c.y-b.y)
	// cross = (b-a).x*(c-b).y - (b-a).y*(c-b).x
	return (b.x-a.x)*(c.y-b.y) - (b.y-a.y)*(c.x-b.x)
}

// ── Approach 1: Orientation-First Cross-Product Check ─────────────────────────
//
// orientationFirst decides convexity by first finding the polygon's overall
// turning direction (the sign of the first non-zero cross product) and then
// verifying that EVERY vertex turns the same way (or is collinear).
//
// Intuition:
//
//	A simple polygon is convex iff, walking its vertices in order (and wrapping
//	around), you always turn the same direction — all left turns, or all right
//	turns; straight (collinear) steps are allowed. Each turn's direction is the
//	SIGN of the cross product of the two edge vectors meeting at that vertex. So:
//	find the reference sign once (skip leading collinear triples), then require
//	no vertex to have the opposite sign.
//
// Algorithm:
//  1. For each vertex i, take the triple (P[i], P[i+1], P[i+2]) with indices mod n
//     (edges wrap around the closed polygon).
//  2. Compute c = cross(...). Ignore c == 0 (collinear, allowed).
//  3. The first non-zero c fixes the reference sign. If any later non-zero c has
//     the opposite sign, the polygon bends both ways → not convex.
//  4. If no conflicting sign is found, it is convex.
//
// Time:  O(n) — one pass over the n vertices (each triple is O(1)).
// Space: O(1) — a couple of integer flags.
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

// ── Approach 2: Single-Pass Both-Signs Flags (Optimal) ───────────────────────
//
// bothSignsFlags decides convexity in one pass by recording whether ANY positive
// cross product and whether ANY negative cross product occurs; a convex polygon
// may produce only one of the two signs (plus zeros).
//
// Intuition:
//
//	Same orientation-consistency fact, phrased so no "reference sign" state is
//	needed: keep two booleans, hasPos and hasNeg. Every vertex's cross product
//	sets one of them. The instant BOTH are true the boundary has turned left
//	somewhere and right somewhere else, which a convex polygon cannot do — return
//	false immediately. If the scan finishes with at most one flag set, it is
//	convex.
//
// Algorithm:
//  1. hasPos = hasNeg = false.
//  2. For each vertex i: c = cross(P[i], P[i+1], P[i+2]) with wraparound indices.
//     If c > 0 set hasPos; if c < 0 set hasNeg. (c == 0 sets neither.)
//  3. If hasPos && hasNeg at any point, return false (bends both ways).
//  4. After the loop return true.
//
// Time:  O(n) — a single linear scan, with early exit on the first sign clash.
// Space: O(1) — two boolean flags.
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

// toPoints converts LeetCode's [][]int representation into []point.
func toPoints(points [][]int) []point {
	P := make([]point, len(points))
	for i, p := range points {
		P[i] = point{p[0], p[1]} // p[0]=x, p[1]=y
	}
	return P
}

func main() {
	square := [][]int{{0, 0}, {0, 1}, {1, 1}, {1, 0}}
	arrow := [][]int{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {5, 5}}
	triangle := [][]int{{0, 0}, {2, 0}, {1, 2}}
	withCollinear := [][]int{{0, 0}, {1, 0}, {2, 0}, {2, 2}, {0, 2}} // an edge has a mid-point

	fmt.Println("=== Approach 1: Orientation-First Cross-Product Check ===")
	fmt.Printf("square             got=%t  expected true\n", orientationFirst(square))
	fmt.Printf("arrow ([5,5] dent) got=%t  expected false\n", orientationFirst(arrow))
	fmt.Printf("triangle           got=%t  expected true\n", orientationFirst(triangle))
	fmt.Printf("rect+collinear pt  got=%t  expected true\n", orientationFirst(withCollinear))

	fmt.Println("=== Approach 2: Single-Pass Both-Signs Flags (Optimal) ===")
	fmt.Printf("square             got=%t  expected true\n", bothSignsFlags(square))
	fmt.Printf("arrow ([5,5] dent) got=%t  expected false\n", bothSignsFlags(arrow))
	fmt.Printf("triangle           got=%t  expected true\n", bothSignsFlags(triangle))
	fmt.Printf("rect+collinear pt  got=%t  expected true\n", bothSignsFlags(withCollinear))
}
