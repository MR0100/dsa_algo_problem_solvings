package main

import (
	"fmt"
	"math"
	"sort"
)

// ── Approach 1: Brute Force ───────────────────────────────────────────────────
//
// bruteForce tries every triplet (i,j,k) and tracks the sum closest to target.
//
// Intuition:
//   Evaluate all O(n³) combinations, keep the one whose sum has the smallest
//   absolute difference from target.
//
// Time:  O(n³).
// Space: O(1).
func bruteForce(nums []int, target int) int {
	n := len(nums)
	closest := nums[0] + nums[1] + nums[2]

	for i := 0; i < n-2; i++ {
		for j := i + 1; j < n-1; j++ {
			for k := j + 1; k < n; k++ {
				sum := nums[i] + nums[j] + nums[k]
				if abs(sum-target) < abs(closest-target) {
					closest = sum
				}
			}
		}
	}
	return closest
}

// ── Approach 2: Sort + Two Pointers (Optimal) ────────────────────────────────
//
// twoPointers sorts nums, then for each i uses converging l/r pointers.
//
// Intuition:
//   Fix index i. Use two pointers l=i+1, r=n-1.
//   sum = nums[i]+nums[l]+nums[r].
//   If sum < target: sum is too small → l++ (increase sum).
//   If sum > target: sum is too big  → r-- (decrease sum).
//   If sum == target: exact match, return immediately.
//   Track the closest sum seen throughout.
//
// Time:  O(n²) — O(n log n) sort + O(n²) two-pointer.
// Space: O(1) extra.
func twoPointers(nums []int, target int) int {
	sort.Ints(nums)
	n := len(nums)
	closest := math.MaxInt32

	for i := 0; i < n-2; i++ {
		// Skip duplicate i values — same triplet candidates.
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}

		l, r := i+1, n-1
		for l < r {
			sum := nums[i] + nums[l] + nums[r]

			if abs(sum-target) < abs(closest-target) {
				closest = sum
			}

			if sum == target {
				return sum // exact match, can't do better
			} else if sum < target {
				l++
			} else {
				r--
			}
		}
	}
	return closest
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	examples := []struct {
		nums   []int
		target int
		expect int
	}{
		{[]int{-1, 2, 1, -4}, 1, 2},
		{[]int{0, 0, 0}, 1, 0},
		{[]int{1, 1, 1, 0}, -100, 2},
		{[]int{-1, 2, 1, -4}, 2, 2},
	}

	approaches := []struct {
		name string
		fn   func([]int, int) int
	}{
		{"Approach 1: Brute Force    O(n³) T | O(1) S", bruteForce},
		{"Approach 2: Two Pointers ✅ O(n²) T | O(1) S", twoPointers},
	}

	for _, ex := range examples {
		fmt.Printf("nums=%v  target=%d  expect=%d\n", ex.nums, ex.target, ex.expect)
		for _, ap := range approaches {
			cp := make([]int, len(ex.nums))
			copy(cp, ex.nums)
			fmt.Printf("  %-48s → %d\n", ap.name, ap.fn(cp, ex.target))
		}
		fmt.Println()
	}
}
