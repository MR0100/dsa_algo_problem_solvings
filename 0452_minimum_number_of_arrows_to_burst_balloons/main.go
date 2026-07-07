package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Greedy, Sort by End (Optimal) ────────────────────────────────
//
// greedySortByEnd solves Minimum Number of Arrows to Burst Balloons by sorting
// the balloons by their right edge and shooting an arrow at the right edge of
// the first balloon of each new group.
//
// Intuition:
//
//	This is the classic "maximum number of non-overlapping intervals" /
//	activity-selection problem in disguise. Think of each arrow as covering a
//	point on the x-axis; a single arrow bursts every balloon whose interval
//	contains that point. To minimise arrows, be greedy: sort balloons by their
//	END coordinate, shoot the very first arrow at the smallest end. That arrow
//	pops every balloon that starts at or before it. The first balloon whose
//	start is strictly past the arrow needs a NEW arrow — again placed at that
//	balloon's end. Placing the arrow at the earliest end keeps it as far left
//	as possible, maximising how many later balloons it can still reach.
//
// Algorithm:
//  1. If there are no balloons, return 0.
//  2. Sort points by xend ascending.
//  3. arrows = 1; arrowX = points[0][1] (end of the first balloon).
//  4. For each subsequent balloon [s, e]:
//     - if s > arrowX, it is not hit → fire a new arrow: arrows++, arrowX = e.
//     - else it is already burst by the current arrow → skip.
//  5. Return arrows.
//
// Time:  O(n log n) — dominated by the sort; the sweep is O(n).
// Space: O(1) extra beyond the sort (Go's sort is in place); O(log n) stack.
func greedySortByEnd(points [][]int) int {
	if len(points) == 0 {
		return 0 // nothing to burst
	}

	// Sort by the right edge so we always know the earliest place an arrow
	// "must" go to still catch the current balloon.
	sort.Slice(points, func(i, j int) bool {
		return points[i][1] < points[j][1] // ascending by xend
	})

	arrows := 1            // the first balloon always needs an arrow
	arrowX := points[0][1] // fire it at the first balloon's right edge
	for i := 1; i < len(points); i++ {
		start := points[i][0]
		// If this balloon starts strictly after our current arrow's x, the
		// arrow (which sits at arrowX) cannot reach it — need a fresh arrow.
		if start > arrowX {
			arrows++              // one more arrow required
			arrowX = points[i][1] // place it at this balloon's right edge
		}
		// Otherwise start <= arrowX <= end (because sorted by end, end >= arrowX),
		// so the current arrow already bursts this balloon — do nothing.
	}
	return arrows
}

// ── Approach 2: Greedy, Sort by Start (Shrink Overlap) ───────────────────────
//
// greedySortByStart solves the same problem sorting by the LEFT edge instead,
// maintaining the shared overlap of the current group and shrinking it.
//
// Intuition:
//
//	Sort by start. Walk the balloons keeping a running "common overlap"
//	[curStart, curEnd] that one arrow could cover for the current group. For
//	each next balloon, if its start is still within curEnd, the group still
//	shares a point — tighten curEnd to min(curEnd, thisEnd) so the overlap
//	stays valid. If its start is beyond curEnd, the group is broken: that
//	arrow is committed, start a brand-new group. Same answer as Approach 1,
//	just tracked from the other side; useful to see the "intersection of
//	intervals" viewpoint.
//
// Algorithm:
//  1. If empty, return 0.
//  2. Sort by xstart ascending.
//  3. arrows = 1; curEnd = points[0][1].
//  4. For each next [s, e]:
//     - if s > curEnd → new group: arrows++, curEnd = e.
//     - else same group: curEnd = min(curEnd, e) (shrink the shared overlap).
//  5. Return arrows.
//
// Time:  O(n log n) — the sort dominates.
// Space: O(1) extra + O(log n) sort stack.
func greedySortByStart(points [][]int) int {
	if len(points) == 0 {
		return 0
	}

	// Sort by left edge; ties broken by right edge for determinism.
	sort.Slice(points, func(i, j int) bool {
		if points[i][0] != points[j][0] {
			return points[i][0] < points[j][0]
		}
		return points[i][1] < points[j][1]
	})

	arrows := 1            // first balloon opens the first group
	curEnd := points[0][1] // the group's shared overlap currently ends here
	for i := 1; i < len(points); i++ {
		s, e := points[i][0], points[i][1]
		if s > curEnd {
			// This balloon starts past the current group's overlap → the
			// group's single arrow can't reach it; commit a new arrow.
			arrows++
			curEnd = e // the new group's overlap starts as this balloon's span
		} else {
			// Still overlapping the group; the shared point can only be as far
			// right as the smallest end seen so far.
			if e < curEnd {
				curEnd = e // shrink overlap to keep it valid for all in group
			}
		}
	}
	return arrows
}

func main() {
	fmt.Println("=== Approach 1: Greedy, Sort by End (Optimal) ===")
	fmt.Printf("points=[[10,16],[2,8],[1,6],[7,12]] -> %d  expected 2\n", greedySortByEnd([][]int{{10, 16}, {2, 8}, {1, 6}, {7, 12}}))
	fmt.Printf("points=[[1,2],[3,4],[5,6],[7,8]]     -> %d  expected 4\n", greedySortByEnd([][]int{{1, 2}, {3, 4}, {5, 6}, {7, 8}}))
	fmt.Printf("points=[[1,2],[2,3],[3,4],[4,5]]     -> %d  expected 2\n", greedySortByEnd([][]int{{1, 2}, {2, 3}, {3, 4}, {4, 5}}))

	fmt.Println("=== Approach 2: Greedy, Sort by Start (Shrink Overlap) ===")
	fmt.Printf("points=[[10,16],[2,8],[1,6],[7,12]] -> %d  expected 2\n", greedySortByStart([][]int{{10, 16}, {2, 8}, {1, 6}, {7, 12}}))
	fmt.Printf("points=[[1,2],[3,4],[5,6],[7,8]]     -> %d  expected 4\n", greedySortByStart([][]int{{1, 2}, {3, 4}, {5, 6}, {7, 8}}))
	fmt.Printf("points=[[1,2],[2,3],[3,4],[4,5]]     -> %d  expected 2\n", greedySortByStart([][]int{{1, 2}, {2, 3}, {3, 4}, {4, 5}}))

	// Edge cases.
	fmt.Println("=== Edge cases ===")
	fmt.Printf("points=[[1,2]]         -> %d  expected 1\n", greedySortByEnd([][]int{{1, 2}}))                                     // single balloon
	fmt.Printf("points=[[1,10],[2,3],[4,5],[6,7]] -> %d  expected 3\n", greedySortByEnd([][]int{{1, 10}, {2, 3}, {4, 5}, {6, 7}})) // wide one straddles
}
