package main

import "fmt"

// ── Approach 1: Brute Force ───────────────────────────────────────────────────
//
// bruteForce checks every pair of lines and tracks the maximum area.
//
// Intuition:
//   The area formed by lines at indices i and j is:
//     area = min(height[i], height[j]) * (j - i)
//   Try all O(n²) pairs and return the maximum.
//
// Time:  O(n²) — all pairs.
// Space: O(1).
func bruteForce(height []int) int {
	n := len(height)
	maxArea := 0
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			// Area = shorter wall × width.
			h := height[i]
			if height[j] < h {
				h = height[j]
			}
			area := h * (j - i)
			if area > maxArea {
				maxArea = area
			}
		}
	}
	return maxArea
}

// ── Approach 2: Two Pointers (Optimal) ───────────────────────────────────────
//
// twoPointers uses converging left and right pointers, always moving the
// pointer with the shorter line inward.
//
// Intuition:
//   Start with the widest possible container (l=0, r=n-1).
//   The area is limited by the shorter wall. Moving the longer wall inward
//   cannot increase the area (width shrinks, height is still bounded by the
//   shorter wall). Moving the shorter wall inward might find a taller wall
//   and increase the area. So always move the shorter wall.
//
// Correctness argument:
//   When we move l inward (because height[l] <= height[r]), we are claiming
//   that no pair (l, j) for j < r can beat the current best. This is true
//   because height[l] is the bottleneck: any j < r gives area
//   min(height[l], height[j]) * (j-l) ≤ height[l] * (r-l) — either the
//   width shrinks or the height is still bounded by height[l]. So (l,r) is
//   the best pair involving l, and we can safely advance l.
//
// Time:  O(n) — each pointer moves at most n steps.
// Space: O(1).
func twoPointers(height []int) int {
	l, r := 0, len(height)-1
	maxArea := 0

	for l < r {
		// Compute area for the current pair.
		h := height[l]
		if height[r] < h {
			h = height[r]
		}
		area := h * (r - l)
		if area > maxArea {
			maxArea = area
		}

		// Move the shorter wall inward — it's the only move that might help.
		if height[l] <= height[r] {
			l++
		} else {
			r--
		}
	}
	return maxArea
}

func main() {
	examples := []struct {
		height []int
		expect int
	}{
		{[]int{1, 8, 6, 2, 5, 4, 8, 3, 7}, 49},
		{[]int{1, 1}, 1},
		{[]int{4, 3, 2, 1, 4}, 16},
		{[]int{1, 2, 1}, 2},
	}

	approaches := []struct {
		name string
		fn   func([]int) int
	}{
		{"Approach 1: Brute Force   O(n²) T | O(1) S", bruteForce},
		{"Approach 2: Two Pointers ✅ O(n) T | O(1) S", twoPointers},
	}

	for _, ex := range examples {
		fmt.Printf("height=%v  expect=%d\n", ex.height, ex.expect)
		for _, ap := range approaches {
			fmt.Printf("  %-47s → %d\n", ap.name, ap.fn(ex.height))
		}
		fmt.Println()
	}
}
