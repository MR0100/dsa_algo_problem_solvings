package main

import "fmt"

// ── Approach 1: Brute Force (Hash Set + Full Range Sweep) ────────────────────
//
// bruteForce solves Missing Ranges by testing every number in [lower, upper]
// for membership and grouping consecutive missing numbers into ranges.
//
// Intuition:
//
//	The definition of "missing" is direct: x ∈ [lower, upper] and x ∉ nums.
//	So put nums into a hash set, sweep x from lower to upper, and whenever we
//	walk into a stretch of missing numbers, remember where it started; when
//	the stretch ends (a present number or the end of the range), emit it.
//	This is only viable when upper − lower is small — the constraints allow
//	a span of ~2·10⁹, which makes this sweep impractical in general.
//
// Algorithm:
//  1. Insert every element of nums into a set.
//  2. Sweep x = lower … upper:
//     a. If x is missing and no stretch is open, open one at x.
//     b. If x is present and a stretch is open, close it as [start, x−1].
//  3. After the sweep, close any still-open stretch as [start, upper].
//
// Time:  O(n + U) where U = upper − lower + 1 — the sweep visits every value.
// Space: O(n) — the hash set of present numbers.
func bruteForce(nums []int, lower, upper int) [][]int {
	present := make(map[int]bool, len(nums)) // O(1) membership tests
	for _, v := range nums {
		present[v] = true
	}
	result := [][]int{}
	start := 0       // first number of the currently open missing stretch
	inRange := false // whether a missing stretch is currently open
	for x := lower; x <= upper; x++ {
		if !present[x] { // x is missing
			if !inRange { // a new missing stretch begins here
				start = x
				inRange = true
			}
		} else if inRange { // x is present → the open stretch ended at x−1
			result = append(result, []int{start, x - 1})
			inRange = false
		}
	}
	if inRange { // the range ended while a stretch was still open
		result = append(result, []int{start, upper})
	}
	return result
}

// ── Approach 2: Linear Scan Over nums (Optimal) ──────────────────────────────
//
// linearScan solves Missing Ranges by walking only the (sorted, unique)
// array and emitting the gap before each element.
//
// Intuition:
//
//	nums is sorted and every element lies inside [lower, upper], so the
//	missing numbers are exactly the gaps: before the first element, between
//	consecutive elements, and after the last element. Tracking "next" — the
//	smallest value not yet covered — turns each gap into one range in O(1),
//	independent of how wide the gap is (crucial when the span is 2·10⁹).
//
// Algorithm:
//  1. next = lower (smallest value still unaccounted for).
//  2. For each v in nums:
//     a. If v > next, the values [next, v−1] are missing → emit that range.
//     b. Set next = v + 1 (everything up to v is now covered).
//  3. If next <= upper after the loop, emit the final range [next, upper].
//
// Time:  O(n) — one pass over nums; each element does O(1) work.
// Space: O(1) — excluding the output list.
func linearScan(nums []int, lower, upper int) [][]int {
	result := [][]int{}
	next := lower // smallest number in [lower, upper] not yet covered
	for _, v := range nums {
		if v > next { // gap [next, v-1] is entirely missing
			result = append(result, []int{next, v - 1})
		}
		next = v + 1 // v itself is present → coverage advances past it
	}
	if next <= upper { // tail gap after the last element of nums
		result = append(result, []int{next, upper})
	}
	return result
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Hash Set + Full Range Sweep) ===")
	fmt.Println(bruteForce([]int{0, 1, 3, 50, 75}, 0, 99)) // expected [[2 2] [4 49] [51 74] [76 99]]
	fmt.Println(bruteForce([]int{-1}, -1, -1))             // expected []

	fmt.Println("=== Approach 2: Linear Scan Over nums (Optimal) ===")
	fmt.Println(linearScan([]int{0, 1, 3, 50, 75}, 0, 99)) // expected [[2 2] [4 49] [51 74] [76 99]]
	fmt.Println(linearScan([]int{-1}, -1, -1))             // expected []
}
