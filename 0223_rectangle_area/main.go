package main

import "fmt"

// ── Approach 1: Inclusion–Exclusion (Optimal) ────────────────────────────────
//
// inclusionExclusion computes the total area covered by two axis-aligned
// rectangles as area1 + area2 − overlap, where the overlap is itself a
// rectangle found by intersecting the two coordinate intervals.
//
// Intuition:
//
//	If we simply add the two areas we double-count the region where they
//	overlap, so subtract that region once. The overlap in x spans
//	[max(left edges), min(right edges)]; likewise in y. If either span is
//	non-positive the rectangles don't overlap and the overlap area is 0.
//
// Algorithm:
//
//  1. area1 = (ax2−ax1)·(ay2−ay1); area2 = (bx2−bx1)·(by2−by1).
//  2. overlapW = max(0, min(ax2,bx2) − max(ax1,bx1)).
//  3. overlapH = max(0, min(ay2,by2) − max(ay1,by1)).
//  4. Return area1 + area2 − overlapW·overlapH.
//
// Time:  O(1) — a handful of arithmetic operations.
// Space: O(1) — no extra storage.
func inclusionExclusion(ax1, ay1, ax2, ay2, bx1, by1, bx2, by2 int) int {
	area1 := (ax2 - ax1) * (ay2 - ay1) // area of the first rectangle
	area2 := (bx2 - bx1) * (by2 - by1) // area of the second rectangle

	// horizontal overlap: from the rightmost left-edge to the leftmost right-edge
	overlapW := maxInt(0, minInt(ax2, bx2)-maxInt(ax1, bx1))
	// vertical overlap: from the topmost bottom-edge to the bottommost top-edge
	overlapH := maxInt(0, minInt(ay2, by2)-maxInt(ay1, by1))

	overlap := overlapW * overlapH // 0 when the rectangles are disjoint

	// add both areas, remove the double-counted intersection once
	return area1 + area2 - overlap
}

// ── Approach 2: Explicit Overlap-Detection Branch ────────────────────────────
//
// explicitOverlap computes the same result but first tests whether the two
// rectangles overlap at all, computing the intersection area only in that
// branch. It is functionally identical to Approach 1 with the max(0, …)
// clamp replaced by an explicit if — useful when you want the overlap flag.
//
// Intuition:
//
//	Two axis-aligned rectangles overlap iff one's left edge is strictly left
//	of the other's right edge in both dimensions. If that holds, the
//	intersection is a rectangle whose corners are the inner edges; otherwise
//	the covered area is just area1 + area2.
//
// Algorithm:
//
//  1. Compute area1, area2.
//  2. Compute intersection edges ix1,iy1,ix2,iy2 from the inner max/min edges.
//  3. If ix1 < ix2 AND iy1 < iy2, they overlap: subtract (ix2−ix1)·(iy2−iy1).
//  4. Otherwise return area1 + area2.
//
// Time:  O(1).
// Space: O(1).
func explicitOverlap(ax1, ay1, ax2, ay2, bx1, by1, bx2, by2 int) int {
	area1 := (ax2 - ax1) * (ay2 - ay1)
	area2 := (bx2 - bx1) * (by2 - by1)

	// inner edges of a potential intersection rectangle
	ix1 := maxInt(ax1, bx1) // left of overlap
	iy1 := maxInt(ay1, by1) // bottom of overlap
	ix2 := minInt(ax2, bx2) // right of overlap
	iy2 := minInt(ay2, by2) // top of overlap

	if ix1 < ix2 && iy1 < iy2 {
		// genuine overlap → subtract its area exactly once
		overlap := (ix2 - ix1) * (iy2 - iy1)
		return area1 + area2 - overlap
	}
	// disjoint or edge-touching (zero-area overlap): no subtraction
	return area1 + area2
}

// maxInt returns the larger of two ints.
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// minInt returns the smaller of two ints.
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	// Example 1: overlapping rectangles.
	fmt.Println("=== Approach 1: Inclusion–Exclusion (Optimal) ===")
	fmt.Println(inclusionExclusion(-3, 0, 3, 4, 0, -1, 9, 2))   // expected 45
	fmt.Println(inclusionExclusion(-2, -2, 2, 2, -2, -2, 2, 2)) // expected 16

	fmt.Println("=== Approach 2: Explicit Overlap-Detection Branch ===")
	fmt.Println(explicitOverlap(-3, 0, 3, 4, 0, -1, 9, 2))   // expected 45
	fmt.Println(explicitOverlap(-2, -2, 2, 2, -2, -2, 2, 2)) // expected 16
}
