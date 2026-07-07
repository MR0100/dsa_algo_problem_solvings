package main

import (
	"fmt"
	"sort"
)

// abs returns the absolute value of an int (Go's stdlib abs is float-only).
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// ── Approach 1: Brute Force (Each House Scans Every Heater) ───────────────────
//
// bruteForce solves Heaters by, for every house, measuring the distance to
// every heater and keeping the closest; the answer is the largest such closest
// distance across all houses.
//
// Intuition:
//
//	A single global radius must cover the worst-off house. For each house the
//	needed radius is the distance to its NEAREST heater (min over heaters). The
//	whole array is covered only if the radius is at least the maximum of those
//	per-house minima. So answer = max_house( min_heater |house - heater| ).
//
// Algorithm:
//  1. ans = 0.
//  2. For each house h: best = min over heaters of |h - heater|.
//  3. ans = max(ans, best).
//  4. Return ans.
//
// Time:  O(H · K) — every house probes every heater (H houses, K heaters).
// Space: O(1) — only running extrema.
func bruteForce(houses []int, heaters []int) int {
	ans := 0 // the largest "distance to nearest heater" seen so far
	for _, h := range houses {
		best := -1 // distance from this house to its closest heater
		for _, ht := range heaters {
			d := abs(h - ht) // distance house→heater
			if best == -1 || d < best {
				best = d // this heater is closer
			}
		}
		if best > ans {
			ans = best // this house needs the radius bumped up
		}
	}
	return ans
}

// ── Approach 2: Sort + Binary Search ─────────────────────────────────────────
//
// binarySearch solves Heaters by sorting the heaters, then for each house
// binary-searching the insertion point to compare only the two neighbouring
// heaters instead of all of them.
//
// Intuition:
//
//	Once heaters are sorted, the nearest heater to a house is one of the two that
//	straddle it: the greatest heater ≤ house, or the least heater ≥ house. Binary
//	search finds that boundary in O(log K); take the smaller of the two gaps as
//	the house's required radius, and the global answer is the max over houses.
//
// Algorithm:
//  1. Sort heaters.
//  2. For each house h: pos = lower_bound(heaters, h).
//     - right neighbour: heaters[pos] if pos < K.
//     - left  neighbour: heaters[pos-1] if pos > 0.
//     dist = min gap to whichever neighbours exist.
//  3. ans = max over houses of dist.
//
// Time:  O((H + K) log K) — sort the heaters, then one binary search per house.
// Space: O(1) — sorting in place (ignoring sort's internal stack).
func binarySearch(houses []int, heaters []int) int {
	sort.Ints(heaters) // sorted heaters enable straddle lookup
	k := len(heaters)
	ans := 0
	for _, h := range houses {
		// pos = first index whose heater value is >= h (lower bound).
		pos := sort.SearchInts(heaters, h)
		dist := 1 << 62 // +inf placeholder
		if pos < k {
			// Heater at pos is >= h: gap to the right neighbour.
			if d := heaters[pos] - h; d < dist {
				dist = d
			}
		}
		if pos > 0 {
			// Heater at pos-1 is < h: gap to the left neighbour.
			if d := h - heaters[pos-1]; d < dist {
				dist = d
			}
		}
		if dist > ans {
			ans = dist // worst house so far dictates the radius
		}
	}
	return ans
}

// ── Approach 3: Sort Both + Two-Pointer Sweep (Optimal) ──────────────────────
//
// twoPointers solves Heaters by sorting BOTH arrays and sweeping a heater
// pointer forward as houses advance, so each house instantly knows its nearest
// heater without repeated searching.
//
// Intuition:
//
//	Sort houses and heaters. Walk houses left→right; keep a heater index j that
//	we advance while the NEXT heater is at least as close to the current house as
//	the current heater. When advancing no longer helps, heaters[j] is the nearest
//	heater to this house — record its distance. Because both arrays are sorted, j
//	never moves backward, giving a linear sweep after the sort.
//
// Algorithm:
//  1. Sort houses and heaters.
//  2. j = 0, ans = 0.
//  3. For each house h (in sorted order):
//     while j+1 < K and |heaters[j+1] - h| <= |heaters[j] - h|: j++.
//     ans = max(ans, |heaters[j] - h|).
//  4. Return ans.
//
// Time:  O(H log H + K log K) — dominated by the two sorts; the sweep is O(H+K).
// Space: O(1) — in-place sorts, two indices.
func twoPointers(houses []int, heaters []int) int {
	sort.Ints(houses)  // sweep houses in increasing order
	sort.Ints(heaters) // heaters aligned so the pointer only moves forward
	j := 0             // index of the current candidate heater
	ans := 0
	k := len(heaters)
	for _, h := range houses {
		// Advance while the next heater is no farther than the current one.
		// Once the next heater is strictly farther, heaters[j] is the closest.
		for j+1 < k && abs(heaters[j+1]-h) <= abs(heaters[j]-h) {
			j++
		}
		if d := abs(heaters[j] - h); d > ans {
			ans = d // this house needs at least distance d
		}
	}
	return ans
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Printf("houses=[1 2 3], heaters=[2]     => %d  (expected 1)\n", bruteForce([]int{1, 2, 3}, []int{2}))
	fmt.Printf("houses=[1 2 3 4], heaters=[1 4] => %d  (expected 1)\n", bruteForce([]int{1, 2, 3, 4}, []int{1, 4}))
	fmt.Printf("houses=[1 5], heaters=[2]       => %d  (expected 3)\n", bruteForce([]int{1, 5}, []int{2}))

	fmt.Println("=== Approach 2: Sort + Binary Search ===")
	fmt.Printf("houses=[1 2 3], heaters=[2]     => %d  (expected 1)\n", binarySearch([]int{1, 2, 3}, []int{2}))
	fmt.Printf("houses=[1 2 3 4], heaters=[1 4] => %d  (expected 1)\n", binarySearch([]int{1, 2, 3, 4}, []int{1, 4}))
	fmt.Printf("houses=[1 5], heaters=[2]       => %d  (expected 3)\n", binarySearch([]int{1, 5}, []int{2}))

	fmt.Println("=== Approach 3: Two-Pointer Sweep (Optimal) ===")
	fmt.Printf("houses=[1 2 3], heaters=[2]     => %d  (expected 1)\n", twoPointers([]int{1, 2, 3}, []int{2}))
	fmt.Printf("houses=[1 2 3 4], heaters=[1 4] => %d  (expected 1)\n", twoPointers([]int{1, 2, 3, 4}, []int{1, 4}))
	fmt.Printf("houses=[1 5], heaters=[2]       => %d  (expected 3)\n", twoPointers([]int{1, 5}, []int{2}))
}
