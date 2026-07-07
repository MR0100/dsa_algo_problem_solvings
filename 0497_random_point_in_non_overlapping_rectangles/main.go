package main

import (
	"fmt"
	"math/rand"
	"sort"
)

// Problem 497 — Random Point in Non-overlapping Rectangles.
//
// We must pick an integer lattice point uniformly at random from the union of
// several non-overlapping axis-aligned rectangles (borders included). "Uniform"
// means every integer point in the union is equally likely — so a rectangle
// that contains more integer points must be chosen proportionally more often.
//
// The number of integer points in rect [x1,y1,x2,y2] (inclusive corners) is
//
//	weight = (x2 - x1 + 1) * (y2 - y1 + 1).
//
// The two approaches below differ only in HOW they turn those weights into a
// uniform pick; both draw the point inside the chosen rectangle uniformly.
//
// Because the output is random, main() cannot assert exact coordinates. Instead
// it fixes the RNG seed and verifies the two *invariants* that define
// correctness:
//  1. every returned point lies inside some input rectangle, and
//  2. over many picks, each rectangle is chosen with frequency proportional to
//     its integer-point count (uniformity).

// ── Approach 1: Prefix-Sum of Weights + Binary Search (Optimal) ──────────────
//
// SolutionBinary picks a rectangle by treating the rectangles' point-counts as
// buckets laid end to end on a number line, drawing one index in [0, total),
// and binary-searching which bucket it fell into.
//
// Intuition:
//
//	Give every integer point a global id 0..total-1: rectangle 0 owns the first
//	w0 ids, rectangle 1 the next w1, and so on. Drawing a uniform id and finding
//	its owning rectangle picks each rectangle with probability wi/total —
//	exactly uniform over all points. Prefix sums of the weights let a binary
//	search find the owner in O(log k).
//
// Algorithm (construct):
//  1. prefix[i] = w0 + w1 + ... + w(i-1)  (running count; prefix[k] = total).
//
// Algorithm (pick):
//  1. Draw target uniformly in [0, total).
//  2. Binary-search the first prefix[i+1] > target → rectangle i owns target.
//  3. Return a uniform integer point inside rectangle i.
//
// Time:  O(k) construction, O(log k) per pick.
// Space: O(k) for the prefix array.
type SolutionBinary struct {
	rects  [][]int // original rectangles
	prefix []int   // prefix[i] = total integer points in rects[0..i-1]
	total  int     // total integer points across all rectangles
}

// ConstructorBinary builds the prefix-sum table over rectangle point counts.
func ConstructorBinary(rects [][]int) SolutionBinary {
	prefix := make([]int, len(rects)+1) // prefix[0]=0; prefix[k]=total
	for i, r := range rects {
		// integer points in an inclusive rectangle = (width+1)*(height+1)
		w := (r[2] - r[0] + 1) * (r[3] - r[1] + 1)
		prefix[i+1] = prefix[i] + w // running cumulative count
	}
	return SolutionBinary{rects: rects, prefix: prefix, total: prefix[len(rects)]}
}

// Pick returns a uniformly-random integer point from the union of rectangles.
func (s *SolutionBinary) Pick() []int {
	target := rand.Intn(s.total) // a uniform "point id" in [0, total)
	// Find the first rectangle whose cumulative count strictly exceeds target;
	// sort.Search returns the least i in [1,k] with prefix[i] > target.
	i := sort.Search(len(s.rects), func(i int) bool {
		return s.prefix[i+1] > target
	})
	r := s.rects[i] // the chosen rectangle
	// draw x uniformly in [x1, x2] and y uniformly in [y1, y2] (inclusive)
	x := r[0] + rand.Intn(r[2]-r[0]+1)
	y := r[1] + rand.Intn(r[3]-r[1]+1)
	return []int{x, y}
}

// ── Approach 2: Weighted Reservoir Sampling (one pass, O(k) per pick) ─────────
//
// SolutionReservoir picks a rectangle in a single streaming pass without any
// precomputed prefix array, using weighted reservoir sampling (A-Res, k=1).
//
// Intuition:
//
//	Scan rectangles while tracking the running total of points seen so far. When
//	rectangle i (weight wi) arrives, replace the current choice with i with
//	probability wi / runningTotal. A classic induction shows that after the last
//	rectangle, rectangle i is held with probability wi / total — the same
//	uniform-over-points distribution, but with no O(k) prefix storage.
//
// Algorithm (pick):
//  1. running = 0, chosen = -1.
//  2. For each rectangle i with weight wi:
//     running += wi; with probability wi/running set chosen = i.
//  3. Return a uniform integer point inside rectangle chosen.
//
// Time:  O(k) per pick (streams all rectangles each time).
// Space: O(1) beyond the stored rectangles.
type SolutionReservoir struct {
	rects [][]int // original rectangles
}

// ConstructorReservoir just stores the rectangles; no precomputation needed.
func ConstructorReservoir(rects [][]int) SolutionReservoir {
	return SolutionReservoir{rects: rects}
}

// Pick returns a uniformly-random integer point via weighted reservoir sampling.
func (s *SolutionReservoir) Pick() []int {
	running := 0 // total integer points seen so far in the stream
	chosen := -1 // index of the currently-held rectangle
	for i, r := range s.rects {
		w := (r[2] - r[0] + 1) * (r[3] - r[1] + 1) // this rectangle's point count
		running += w                               // extend the seen-so-far total
		// keep rectangle i with probability w/running (rand.Intn(running) < w)
		if rand.Intn(running) < w {
			chosen = i
		}
	}
	r := s.rects[chosen]               // the survivor of the reservoir
	x := r[0] + rand.Intn(r[2]-r[0]+1) // uniform x in [x1, x2]
	y := r[1] + rand.Intn(r[3]-r[1]+1) // uniform y in [y1, y2]
	return []int{x, y}
}

// pointInsideAnyRect reports whether (x,y) lies within any rectangle (inclusive
// borders). Used only to validate outputs deterministically.
func pointInsideAnyRect(rects [][]int, x, y int) bool {
	for _, r := range rects {
		if x >= r[0] && x <= r[2] && y >= r[1] && y <= r[3] {
			return true
		}
	}
	return false
}

func main() {
	rand.Seed(1) // fix the seed so this program prints the same lines every run

	// Example rectangles: rect0 is 5x5 (25 points), rect1 is 2x2 (4 points).
	rects := [][]int{{1, 1, 5, 5}, {-2, -2, -1, -1}}
	const trials = 200000

	fmt.Println("=== Approach 1: Prefix-Sum + Binary Search (Optimal) ===")
	sb := ConstructorBinary(rects)
	fmt.Println("total integer points:", sb.total) // expected 29 (25 + 4)
	// Validity: sample a few points and confirm each lands inside a rectangle.
	allInsideB := true
	countRect0B := 0
	for i := 0; i < trials; i++ {
		p := sb.Pick()
		if !pointInsideAnyRect(rects, p[0], p[1]) {
			allInsideB = false
		}
		if p[0] >= 1 { // x>=1 uniquely identifies rect0 here
			countRect0B++
		}
	}
	fmt.Println("all picks inside a rectangle:", allInsideB)                                  // expected true
	fmt.Printf("rect0 share ~25/29=%.3f, got %.3f\n", 25.0/29.0, float64(countRect0B)/trials) // expected close

	fmt.Println("=== Approach 2: Weighted Reservoir Sampling ===")
	sr := ConstructorReservoir(rects)
	allInsideR := true
	countRect0R := 0
	for i := 0; i < trials; i++ {
		p := sr.Pick()
		if !pointInsideAnyRect(rects, p[0], p[1]) {
			allInsideR = false
		}
		if p[0] >= 1 {
			countRect0R++
		}
	}
	fmt.Println("all picks inside a rectangle:", allInsideR)                                  // expected true
	fmt.Printf("rect0 share ~25/29=%.3f, got %.3f\n", 25.0/29.0, float64(countRect0R)/trials) // expected close
}
