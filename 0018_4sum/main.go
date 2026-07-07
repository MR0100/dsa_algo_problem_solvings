package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Brute Force ───────────────────────────────────────────────────
//
// bruteForce tries all O(n⁴) quadruplets and deduplicates via a sorted-quad map.
//
// Time:  O(n⁴).
// Space: O(n) — the dedup map.
func bruteForce(nums []int, target int) [][]int {
	n := len(nums)
	seen := make(map[[4]int]bool)
	var result [][]int

	for i := 0; i < n-3; i++ {
		for j := i + 1; j < n-2; j++ {
			for k := j + 1; k < n-1; k++ {
				for l := k + 1; l < n; l++ {
					if nums[i]+nums[j]+nums[k]+nums[l] == target {
						quad := [4]int{nums[i], nums[j], nums[k], nums[l]}
						sort.Ints(quad[:])
						if !seen[quad] {
							seen[quad] = true
							result = append(result, []int{quad[0], quad[1], quad[2], quad[3]})
						}
					}
				}
			}
		}
	}
	return result
}

// ── Approach 2: Sort + Two Outer Loops + Two Pointers (Optimal) ──────────────
//
// twoPointers extends the 3Sum two-pointer pattern by adding one more outer loop.
//
// Intuition:
//   Sort the array. Fix two elements (i and j), then find pairs in nums[j+1:]
//   that sum to (target - nums[i] - nums[j]) using converging two pointers.
//   Skip duplicate values of i and j to avoid duplicate quadruplets.
//   Overflow guard: cast to int64 for the sum comparison since values can be
//   large (|nums[i]| ≤ 10⁹).
//
// Time:  O(n³) — O(n log n) sort + O(n²) for the two outer loops × O(n) inner.
// Space: O(1) extra (sort in-place).
func twoPointers(nums []int, target int) [][]int {
	sort.Ints(nums)
	n := len(nums)
	var result [][]int

	for i := 0; i < n-3; i++ {
		// Skip duplicate first element.
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}
		for j := i + 1; j < n-2; j++ {
			// Skip duplicate second element.
			if j > i+1 && nums[j] == nums[j-1] {
				continue
			}

			l, r := j+1, n-1
			for l < r {
				// Use int64 to avoid overflow when nums can be ±10⁹.
				sum := int64(nums[i]) + int64(nums[j]) + int64(nums[l]) + int64(nums[r])
				if sum == int64(target) {
					result = append(result, []int{nums[i], nums[j], nums[l], nums[r]})
					for l < r && nums[l] == nums[l+1] {
						l++
					}
					for l < r && nums[r] == nums[r-1] {
						r--
					}
					l++
					r--
				} else if sum < int64(target) {
					l++
				} else {
					r--
				}
			}
		}
	}
	return result
}

func main() {
	examples := []struct {
		nums   []int
		target int
		expect string
	}{
		{[]int{1, 0, -1, 0, -2, 2}, 0, "[[−2,−1,1,2],[−2,0,0,2],[−1,0,0,1]]"},
		{[]int{2, 2, 2, 2, 2}, 8, "[[2,2,2,2]]"},
		{[]int{0, 0, 0, 0}, 0, "[[0,0,0,0]]"},
	}

	approaches := []struct {
		name string
		fn   func([]int, int) [][]int
	}{
		{"Approach 1: Brute Force    O(n⁴) T | O(n) S", bruteForce},
		{"Approach 2: Two Pointers ✅ O(n³) T | O(1) S", twoPointers},
	}

	for _, ex := range examples {
		fmt.Printf("nums=%v  target=%d\n", ex.nums, ex.target)
		for _, ap := range approaches {
			cp := make([]int, len(ex.nums))
			copy(cp, ex.nums)
			fmt.Printf("  %-46s → %v\n", ap.name, ap.fn(cp, ex.target))
		}
		fmt.Println()
	}
}
