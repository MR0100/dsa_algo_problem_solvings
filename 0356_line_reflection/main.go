package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Brute Force (Pairwise Check) ─────────────────────────────────
//
// bruteForce solves Line Reflection by first guessing the mirror line, then
// verifying every point has a partner reflected across it.
//
// Intuition:
//
//	If a vertical line x = k reflects the whole set onto itself, then it must
//	sit exactly halfway between the leftmost and rightmost x. So the candidate
//	is fixed: sum = minX + maxX (we compare 2*k against x1+x2 to avoid halves).
//	With the line pinned, the only thing left is to confirm that for EVERY
//	point (x, y) its mirror image (sum - x, y) also exists. The brute-force
//	way to answer "does the mirror exist" is a linear scan over all points.
//
// Algorithm:
//  1. Find minX and maxX; let sum = minX + maxX (2 * mirror line).
//  2. For each point p, scan the whole list for a point q with
//     q.x == sum - p.x and q.y == p.y.
//  3. If any point has no such partner, return false; else true.
//
// Time:  O(n^2) — for each of n points, a full O(n) scan for its mirror.
// Space: O(1) — only min/max and loop indices.
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

// ── Approach 2: Sorting + Two Pointers ───────────────────────────────────────
//
// twoPointers solves Line Reflection by sorting points and matching them from
// both ends inward around the candidate mirror line.
//
// Intuition:
//
//	Sort points by (y, x). Within each y-group the reflected pairs are the
//	symmetric outer/inner elements, so walking one pointer from the left and
//	one from the right must always meet at partners whose x-values sum to the
//	fixed 2*mirror. Duplicates and single centre points (x == mirror) fall out
//	naturally because the left/right x still sum correctly.
//
// Algorithm:
//  1. Compute sum = minX + maxX.
//  2. Sort points by y ascending, then x ascending.
//  3. Use i from the start and j from the end; require points[i].y ==
//     points[j].y and points[i].x + points[j].x == sum, advancing inward.
//  4. Any mismatch ⇒ false.
//
// Time:  O(n log n) — dominated by the sort; the scan is linear.
// Space: O(1) auxiliary (sort in place; ignoring sort's own recursion).
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

// ── Approach 3: Hash Set (Optimal) ───────────────────────────────────────────
//
// hashSet solves Line Reflection by storing every point in a set and checking
// each point's mirror membership in O(1).
//
// Intuition:
//
//	The brute force wastes time re-scanning for each mirror. Put every point
//	into a hash set keyed by (x, y). The candidate mirror is still fixed by
//	sum = minX + maxX. Now "does the mirror exist" is a single O(1) set
//	lookup for (sum - x, y). One pass builds the set, one pass verifies.
//
// Algorithm:
//  1. Insert every point into a set; track minX, maxX.
//  2. sum = minX + maxX.
//  3. For each point (x, y), test membership of (sum - x, y). Missing ⇒ false.
//
// Time:  O(n) — two linear passes with O(1) set operations.
// Space: O(n) — the set holds up to n distinct points.
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

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Printf("[[1,1],[-1,1]]           got=%t  expected true\n", bruteForce([][]int{{1, 1}, {-1, 1}}))
	fmt.Printf("[[1,1],[-1,-1]]          got=%t  expected false\n", bruteForce([][]int{{1, 1}, {-1, -1}}))
	fmt.Printf("[[0,0],[1,0]]            got=%t  expected true\n", bruteForce([][]int{{0, 0}, {1, 0}}))
	fmt.Printf("[[1,1],[-1,1],[0,1]]     got=%t  expected true\n", bruteForce([][]int{{1, 1}, {-1, 1}, {0, 1}}))

	fmt.Println("=== Approach 2: Sorting + Two Pointers ===")
	fmt.Printf("[[1,1],[-1,1]]           got=%t  expected true\n", twoPointers([][]int{{1, 1}, {-1, 1}}))
	fmt.Printf("[[1,1],[-1,-1]]          got=%t  expected false\n", twoPointers([][]int{{1, 1}, {-1, -1}}))
	fmt.Printf("[[0,0],[1,0]]            got=%t  expected true\n", twoPointers([][]int{{0, 0}, {1, 0}}))
	fmt.Printf("[[1,1],[-1,1],[0,1]]     got=%t  expected true\n", twoPointers([][]int{{1, 1}, {-1, 1}, {0, 1}}))

	fmt.Println("=== Approach 3: Hash Set (Optimal) ===")
	fmt.Printf("[[1,1],[-1,1]]           got=%t  expected true\n", hashSet([][]int{{1, 1}, {-1, 1}}))
	fmt.Printf("[[1,1],[-1,-1]]          got=%t  expected false\n", hashSet([][]int{{1, 1}, {-1, -1}}))
	fmt.Printf("[[0,0],[1,0]]            got=%t  expected true\n", hashSet([][]int{{0, 0}, {1, 0}}))
	fmt.Printf("[[1,1],[-1,1],[0,1]]     got=%t  expected true\n", hashSet([][]int{{1, 1}, {-1, 1}, {0, 1}}))
}
