package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Brute Force ───────────────────────────────────────────────────
//
// bruteForce checks every pair of indices (i, j) where i < j and returns
// the pair whose values sum to target.
//
// Intuition:
//   The simplest possible idea — try every combination. For each element,
//   scan every element that comes after it to see if the two add up to target.
//
// Algorithm:
//   1. Outer loop: pick index i from 0 to n-2.
//   2. Inner loop: pick index j from i+1 to n-1.
//   3. If nums[i] + nums[j] == target, return [i, j].
//
// Time:  O(n²) — two nested loops, each up to n iterations.
// Space: O(1)  — no extra data structures; only loop variables.
func bruteForce(nums []int, target int) []int {
	n := len(nums)
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			// Check if this pair sums to target.
			if nums[i]+nums[j] == target {
				return []int{i, j}
			}
		}
	}
	return nil // guaranteed one solution exists, so never reached
}

// ── Approach 2: Two-Pass Hash Map ────────────────────────────────────────────
//
// twoPassHashMap builds a complete value→index map first, then looks up each
// element's complement in a second pass.
//
// Intuition:
//   A hash map turns "does X exist in this array?" from O(n) scan to O(1)
//   lookup. Build the map once, then query it once per element.
//
// Algorithm:
//   Pass 1: Insert every nums[i] → i into the map.
//   Pass 2: For each i, compute complement = target - nums[i].
//           Look up complement in the map.
//           If found at index j and j != i, return [i, j].
//
// Time:  O(n) — two separate O(n) passes.
// Space: O(n) — the hash map holds up to n entries.
func twoPassHashMap(nums []int, target int) []int {
	// Pass 1: build value → index map.
	indexMap := make(map[int]int) // value → index
	for i, v := range nums {
		indexMap[v] = i
	}

	// Pass 2: for each element look up its complement.
	for i, v := range nums {
		complement := target - v
		// Check map; guard j != i to avoid using the same element twice.
		if j, ok := indexMap[complement]; ok && j != i {
			return []int{i, j}
		}
	}
	return nil
}

// ── Approach 3: One-Pass Hash Map (Optimal) ───────────────────────────────────
//
// onePassHashMap finds the answer in a single scan by checking the complement
// before inserting the current element into the map.
//
// Intuition:
//   While building the map we can simultaneously ask: "Have I already seen
//   the number I need to pair with the current element?" If yes, we're done.
//   Because we check before inserting, we never pair an element with itself,
//   and we never need a second pass.
//
// Algorithm:
//   For each index i:
//     1. Compute complement = target - nums[i].
//     2. If complement is already in the map, return [map[complement], i].
//     3. Otherwise, store nums[i] → i in the map and continue.
//
// Time:  O(n) — single pass; each element is visited exactly once.
// Space: O(n) — the map holds at most n entries.
func onePassHashMap(nums []int, target int) []int {
	seen := make(map[int]int) // value → index of elements visited so far

	for i, v := range nums {
		complement := target - v

		// If complement was already seen, we found our pair.
		if j, ok := seen[complement]; ok {
			return []int{j, i}
		}

		// Record this element for future lookups.
		seen[v] = i
	}
	return nil
}

// ── Approach 4: Sorting + Two Pointers ───────────────────────────────────────
//
// sortAndTwoPointers sorts the values (keeping original indices) then uses
// two pointers converging from both ends to find the pair.
//
// Intuition:
//   In a sorted array, if the sum of the two outermost elements is too small,
//   we need a larger left element (move left pointer right). If it's too large,
//   we need a smaller right element (move right pointer left). This eliminates
//   half the remaining candidates on each step — no need for a hash map.
//
//   Note: because the problem asks for original indices we must preserve the
//   index information through the sort by sorting pairs (value, originalIndex).
//
// Algorithm:
//   1. Create a slice of (value, originalIndex) pairs.
//   2. Sort by value.
//   3. Place l = 0, r = n-1.
//   4. While l < r:
//        sum = pairs[l].value + pairs[r].value
//        if sum == target → return [pairs[l].index, pairs[r].index]
//        if sum < target  → l++ (need larger value on left)
//        if sum > target  → r-- (need smaller value on right)
//
// Time:  O(n log n) — dominated by the sort; the two-pointer scan is O(n).
// Space: O(n)       — the auxiliary pairs slice.
func sortAndTwoPointers(nums []int, target int) []int {
	// Pair each value with its original index so we can recover it after sorting.
	type pair struct {
		val, idx int
	}
	pairs := make([]pair, len(nums))
	for i, v := range nums {
		pairs[i] = pair{v, i}
	}

	// Sort by value ascending.
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].val < pairs[j].val
	})

	l, r := 0, len(pairs)-1
	for l < r {
		sum := pairs[l].val + pairs[r].val
		switch {
		case sum == target:
			// Return the original indices in ascending order.
			a, b := pairs[l].idx, pairs[r].idx
			if a > b {
				a, b = b, a
			}
			return []int{a, b}
		case sum < target:
			l++ // left value too small, move right
		default:
			r-- // right value too large, move left
		}
	}
	return nil
}

func main() {
	examples := []struct {
		nums   []int
		target int
		expect []int
	}{
		{[]int{2, 7, 11, 15}, 9, []int{0, 1}},
		{[]int{3, 2, 4}, 6, []int{1, 2}},
		{[]int{3, 3}, 6, []int{0, 1}},
	}

	approaches := []struct {
		name string
		fn   func([]int, int) []int
	}{
		{"Approach 1: Brute Force         O(n²) T | O(1) S", bruteForce},
		{"Approach 2: Two-Pass Hash Map   O(n)  T | O(n) S", twoPassHashMap},
		{"Approach 3: One-Pass Hash Map   O(n)  T | O(n) S", onePassHashMap},
		{"Approach 4: Sort + Two Pointers O(n log n) T | O(n) S", sortAndTwoPointers},
	}

	for _, ex := range examples {
		fmt.Printf("nums=%v  target=%d  expect=%v\n", ex.nums, ex.target, ex.expect)
		for _, ap := range approaches {
			result := ap.fn(ex.nums, ex.target)
			fmt.Printf("  %-55s → %v\n", ap.name, result)
		}
		fmt.Println()
	}
}
