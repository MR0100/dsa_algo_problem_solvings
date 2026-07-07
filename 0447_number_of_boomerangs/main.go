package main

import "fmt"

// ── Approach 1: Brute Force (Triple Loop) ────────────────────────────────────
//
// bruteForce solves Number of Boomerangs by testing every ordered triple
// (i, j, k) and checking whether i is equidistant from j and k.
//
// Intuition:
//
//	A boomerang is an ordered triple (i, j, k) with dist(i,j) == dist(i,k) and
//	j != k. The definition is directly checkable: fix a center i, then try
//	every ordered pair (j, k) of the other points and count the ones whose two
//	distances from i match. Order matters, so (i,j,k) and (i,k,j) are both
//	counted. Comparing squared distances avoids floating-point sqrt.
//
// Algorithm:
//  1. For each center i, for each j != i, for each k != i with k != j:
//     if squaredDist(i,j) == squaredDist(i,k) → count++.
//  2. Return count.
//
// Time:  O(n^3) — three nested loops over n points.
// Space: O(1) — only a running counter.
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

// ── Approach 2: Hash Map by Center (Optimal) ─────────────────────────────────
//
// hashMap solves Number of Boomerangs by, for each center i, grouping the other
// points by their squared distance to i and using a permutation count.
//
// Intuition:
//
//	Fix the apex i. Among the remaining points, suppose m of them lie at the
//	SAME distance from i. Any ordered pair drawn from those m points forms a
//	valid boomerang (i, j, k): that is m * (m-1) ordered pairs (choose j, then a
//	different k). So instead of the cubic pair scan, bucket the distances in a
//	hash map (distance → how many points sit at it), then sum m*(m-1) over
//	buckets. This collapses the inner two loops into a single pass plus arithmetic.
//
// Algorithm:
//  1. For each center i:
//     a. Build map dist2 -> count over all j != i.
//     b. For each bucket with m points, add m*(m-1) to the answer.
//  2. Return the total.
//
// Time:  O(n^2) — for each of n centers, one O(n) pass to bucket distances.
// Space: O(n) — the per-center distance map holds up to n-1 entries.
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

func main() {
	fmt.Println("=== Approach 1: Brute Force (Triple Loop) ===")
	fmt.Printf("points=[[0,0],[1,0],[2,0]]  got=%d  expected 2\n", bruteForce([][]int{{0, 0}, {1, 0}, {2, 0}}))
	fmt.Printf("points=[[1,1],[2,2],[3,3]]  got=%d  expected 2\n", bruteForce([][]int{{1, 1}, {2, 2}, {3, 3}}))
	fmt.Printf("points=[[1,1]]              got=%d  expected 0\n", bruteForce([][]int{{1, 1}}))

	fmt.Println("=== Approach 2: Hash Map by Center (Optimal) ===")
	fmt.Printf("points=[[0,0],[1,0],[2,0]]  got=%d  expected 2\n", hashMap([][]int{{0, 0}, {1, 0}, {2, 0}}))
	fmt.Printf("points=[[1,1],[2,2],[3,3]]  got=%d  expected 2\n", hashMap([][]int{{1, 1}, {2, 2}, {3, 3}}))
	fmt.Printf("points=[[1,1]]              got=%d  expected 0\n", hashMap([][]int{{1, 1}}))
}
