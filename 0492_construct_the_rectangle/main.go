package main

import (
	"fmt"
	"math"
)

// Goal: given `area`, return [L, W] with L*W == area, L >= W, and L-W minimized.
// The minimal L-W means W should be the LARGEST divisor of area that is <= sqrt(area);
// then L = area / W is the smallest divisor >= sqrt(area), so the pair hugs the
// square root as tightly as possible.

// ── Approach 1: Brute Force Downward Scan ────────────────────────────────────
//
// bruteForce tries every candidate width from area down to 1 and returns the
// first one that divides area — but that alone is O(area). We instead start
// from area and walk down, which is the naive reading of the problem; it is
// here only to contrast with the optimal sqrt scan.
//
// Intuition:
//
//	We want the width closest to sqrt(area). A correct-but-slow way: iterate a
//	candidate width w from floor(area) ... actually the smallest search that is
//	obviously correct is "for w from area down to 1, first divisor with w <=
//	area/w wins". To keep it a genuine brute force yet finite, we scan w from
//	1..area and remember the best (largest) divisor not exceeding area/w.
//
// Algorithm:
//  1. best = 1 (width 1 always divides area).
//  2. For w = 1..area: if area % w == 0 and w <= area/w, update best = w.
//  3. Return [area/best, best].
//
// Time:  O(area) — scans every integer up to area. TLE for area up to 1e7 in
//
//	tight limits, but always correct.
//
// Space: O(1).
func bruteForce(area int) []int {
	best := 1 // width 1 is always a valid divisor
	for w := 1; w <= area; w++ {
		if area%w != 0 {
			continue // w must divide area to form an integer rectangle
		}
		l := area / w
		if w <= l { // keep W <= L; the largest such w gives the smallest L-W
			best = w // remember the widest width seen so far that stays <= its length
		}
	}
	return []int{area / best, best} // [L, W]
}

// ── Approach 2: Start at sqrt(area) and Walk Down (Optimal) ───────────────────
//
// sqrtScan starts the width at floor(sqrt(area)) and decreases until it divides
// area. That first divisor is the largest divisor <= sqrt(area), so its partner
// L = area/w is the closest length above the square root — minimal L-W.
//
// Intuition:
//
//	The pair minimizing L-W is the one nearest the square root: for a fixed
//	product, the difference L-W shrinks as W approaches sqrt(area) from below.
//	So begin W at floor(sqrt(area)) and step down; the first W that divides area
//	is the answer's width, and L = area/W its length.
//
// Algorithm:
//  1. w = floor(sqrt(area)).
//  2. While area % w != 0: w--.
//  3. Return [area/w, w].
//
// Time:  O(sqrt(area)) — at most sqrt(area) decrements before hitting a divisor.
// Space: O(1).
func sqrtScan(area int) []int {
	w := int(math.Sqrt(float64(area))) // largest candidate width to try first
	for area%w != 0 {                  // walk down until w divides area exactly
		w-- // the first divisor found is the largest divisor <= sqrt(area)
	}
	return []int{area / w, w} // L = area/w (>= w), W = w
}

func main() {
	fmt.Println("=== Approach 1: Brute Force Downward Scan ===")
	fmt.Println(bruteForce(4))      // expected [2 2]
	fmt.Println(bruteForce(37))     // expected [37 1]
	fmt.Println(bruteForce(122122)) // expected [427 286]

	fmt.Println("=== Approach 2: Start at sqrt(area) and Walk Down (Optimal) ===")
	fmt.Println(sqrtScan(4))      // expected [2 2]
	fmt.Println(sqrtScan(37))     // expected [37 1]
	fmt.Println(sqrtScan(122122)) // expected [427 286]
	fmt.Println(sqrtScan(1))      // expected [1 1]
}
