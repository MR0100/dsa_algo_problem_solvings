package main

import (
	"fmt"
	"sort"
)

// abs returns the absolute value of an int (Go's math.Abs is float64-only).
func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

// ── Approach 1: Brute Force (Try Every Target Value) ─────────────────────────
//
// bruteForce solves Minimum Moves to Equal Array Elements II by trying every
// candidate final value between min and max and summing the moves it costs.
//
// Intuition:
//
//	Each move changes one element by ±1, so making every element equal to some
//	target t costs Σ |nums[i] - t| moves. The optimal target must lie between
//	the array's min and max (going outside that range only adds distance). So
//	try each integer target in [min, max], compute its total cost, keep the
//	smallest. Correct, but the value range can be huge, so this is only a
//	reference baseline.
//
// Algorithm:
//  1. Find lo = min(nums), hi = max(nums).
//  2. For each target t in [lo, hi]: cost = Σ |nums[i] - t|.
//  3. Track and return the minimum cost seen.
//
// Time:  O(n · (hi - lo)) — a full array pass per candidate value.
// Space: O(1) — only counters.
func bruteForce(nums []int) int {
	lo, hi := nums[0], nums[0] // seed min and max with the first element
	for _, v := range nums {
		if v < lo {
			lo = v
		}
		if v > hi {
			hi = v
		}
	}
	best := -1 // sentinel: no cost computed yet
	// Every integer target in the value range is a candidate.
	for t := lo; t <= hi; t++ {
		cost := 0
		for _, v := range nums {
			cost += abs(v - t) // moves to drag v onto target t
		}
		if best == -1 || cost < best {
			best = cost // remember the cheapest target so far
		}
	}
	return best
}

// ── Approach 2: Sort + Median (Two Pointers) ─────────────────────────────────
//
// sortMedian solves Minimum Moves to Equal Array Elements II by sorting and
// pairing the smallest with the largest, accumulating their gaps.
//
// Intuition:
//
//	Σ |nums[i] - t| is minimised when t is the median. Equivalent (and it
//	avoids reasoning about odd/even length): after sorting, pair the outermost
//	elements. To bring nums[i] (smallest remaining) and nums[j] (largest
//	remaining) to any common point strictly between them costs exactly
//	nums[j] - nums[i], independent of where in that gap the meeting point is.
//	Summing (right - left) over all nested pairs gives the total distance to
//	the median without ever computing the median explicitly.
//
// Algorithm:
//  1. Sort nums ascending.
//  2. i = 0, j = n-1, moves = 0.
//  3. While i < j: moves += nums[j] - nums[i]; i++; j--.
//  4. Return moves.
//
// Time:  O(n log n) — dominated by the sort.
// Space: O(1) extra (in-place sort; O(log n) recursion aside).
func sortMedian(nums []int) int {
	sorted := make([]int, len(nums)) // copy so the input is not mutated
	copy(sorted, nums)
	sort.Ints(sorted) // ascending order

	moves := 0
	i, j := 0, len(sorted)-1 // outermost pair of pointers
	// Each pair contributes its span; the meeting point (the median) cancels.
	for i < j {
		moves += sorted[j] - sorted[i] // cost to converge this outer pair
		i++                            // shrink toward the center
		j--
	}
	return moves
}

// ── Approach 3: Quickselect Median (Optimal) ─────────────────────────────────
//
// quickselectMedian solves Minimum Moves to Equal Array Elements II by finding
// the median in expected linear time (no full sort), then summing distances.
//
// Intuition:
//
//	We do not need the array fully sorted — only the median element. Hoare's
//	quickselect partitions around a pivot and recurses into just the side
//	containing index n/2, giving the median in expected O(n). Once the median
//	m is known, the answer is Σ |nums[i] - m|. This is the asymptotically
//	optimal solution.
//
// Algorithm:
//  1. m = quickselect(nums, n/2)  — the lower median for even n (which is fine,
//     any value between the two central elements is optimal).
//  2. moves = Σ |nums[i] - m|.
//  3. Return moves.
//
// Time:  O(n) expected (quickselect) + O(n) for the sum = O(n) expected.
// Space: O(n) for the working copy (we avoid mutating the caller's slice).
func quickselectMedian(nums []int) int {
	work := make([]int, len(nums)) // partitioning reorders elements; copy first
	copy(work, nums)

	median := quickselect(work, 0, len(work)-1, len(work)/2) // k-th smallest, k=n/2
	moves := 0
	for _, v := range nums {
		moves += abs(v - median) // total distance to the median
	}
	return moves
}

// quickselect returns the element that would sit at index k if a[lo..hi] were
// sorted, using Lomuto partitioning and recursing into only one side.
func quickselect(a []int, lo, hi, k int) int {
	for lo < hi {
		p := partition(a, lo, hi) // pivot lands at its final sorted index p
		switch {
		case p == k:
			return a[k] // pivot is exactly the k-th element
		case p < k:
			lo = p + 1 // target is to the right of the pivot
		default:
			hi = p - 1 // target is to the left of the pivot
		}
	}
	return a[lo] // single-element window is the answer
}

// partition places a[hi] at its sorted position and returns that index,
// with everything smaller to its left and everything larger to its right.
func partition(a []int, lo, hi int) int {
	pivot := a[hi] // choose the last element as pivot
	i := lo        // boundary: a[lo..i-1] are < pivot
	for j := lo; j < hi; j++ {
		if a[j] < pivot {
			a[i], a[j] = a[j], a[i] // push a smaller element into the left region
			i++
		}
	}
	a[i], a[hi] = a[hi], a[i] // drop the pivot just past the smaller region
	return i
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Try Every Target Value) ===")
	fmt.Printf("[1,2,3]      got=%d  expected 2\n", bruteForce([]int{1, 2, 3}))      // 1→2, 3→2
	fmt.Printf("[1,10,2,9]   got=%d  expected 16\n", bruteForce([]int{1, 10, 2, 9})) // |1-9|+|10-9|+... etc

	fmt.Println("=== Approach 2: Sort + Median (Two Pointers) ===")
	fmt.Printf("[1,2,3]      got=%d  expected 2\n", sortMedian([]int{1, 2, 3}))
	fmt.Printf("[1,10,2,9]   got=%d  expected 16\n", sortMedian([]int{1, 10, 2, 9}))

	fmt.Println("=== Approach 3: Quickselect Median (Optimal) ===")
	fmt.Printf("[1,2,3]      got=%d  expected 2\n", quickselectMedian([]int{1, 2, 3}))
	fmt.Printf("[1,10,2,9]   got=%d  expected 16\n", quickselectMedian([]int{1, 10, 2, 9}))
}
