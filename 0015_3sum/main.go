package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Brute Force ───────────────────────────────────────────────────
//
// bruteForce tries every triplet (i,j,k) and collects those that sum to zero.
// A set-of-strings deduplicates results.
//
// Intuition:
//   Try all O(n³) combinations. Use a sorted string key to avoid duplicate
//   triplets in the output.
//
// Time:  O(n³) — three nested loops.
// Space: O(n) — the output list (deduplicated via map).
func bruteForce(nums []int) [][]int {
	n := len(nums)
	seen := make(map[[3]int]bool)
	var result [][]int

	for i := 0; i < n-2; i++ {
		for j := i + 1; j < n-1; j++ {
			for k := j + 1; k < n; k++ {
				if nums[i]+nums[j]+nums[k] == 0 {
					// Sort the triple so duplicates have the same key.
					triple := [3]int{nums[i], nums[j], nums[k]}
					sort.Ints(triple[:])
					if !seen[triple] {
						seen[triple] = true
						result = append(result, []int{triple[0], triple[1], triple[2]})
					}
				}
			}
		}
	}
	return result
}

// ── Approach 2: Sort + Two Pointers (Optimal) ────────────────────────────────
//
// twoPointers sorts nums, then for each index i uses converging left/right
// pointers to find pairs summing to -nums[i].
//
// Intuition:
//   After sorting, fix the first element at index i. The problem reduces to
//   finding pairs in nums[i+1:] that sum to target = -nums[i]. Two pointers
//   (l = i+1, r = n-1) converge: if sum < target → l++, if sum > target → r--,
//   if equal → record and skip duplicates for both l and r.
//   Skip duplicate values of i to avoid duplicate triplets.
//
// Time:  O(n²) — O(n log n) sort + O(n) outer loop × O(n) two-pointer inner.
// Space: O(1) extra (sort is in-place; output space is O(k) where k = results).
func twoPointers(nums []int) [][]int {
	sort.Ints(nums)
	n := len(nums)
	var result [][]int

	for i := 0; i < n-2; i++ {
		// Optimisation: if smallest remaining element is positive, no solution.
		if nums[i] > 0 {
			break
		}
		// Skip duplicate values of i to avoid duplicate triplets.
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}

		l, r := i+1, n-1
		target := -nums[i]

		for l < r {
			sum := nums[l] + nums[r]
			if sum == target {
				result = append(result, []int{nums[i], nums[l], nums[r]})
				// Skip duplicates for the left pointer.
				for l < r && nums[l] == nums[l+1] {
					l++
				}
				// Skip duplicates for the right pointer.
				for l < r && nums[r] == nums[r-1] {
					r--
				}
				l++
				r--
			} else if sum < target {
				l++
			} else {
				r--
			}
		}
	}
	return result
}

// ── Approach 3: Hash Set per Pair ────────────────────────────────────────────
//
// hashSet fixes i, then for each j > i looks for -(nums[i]+nums[j]) in a set
// built from nums[i+1..j-1]. Avoids the sort but requires a deduplication step.
//
// Intuition:
//   For each pair (i, j), the required third element is complement = -nums[i]-nums[j].
//   We check if complement appeared before j (in the inner pass) using a set.
//   Deduplication is handled by sorting each found triplet and using a set of triplets.
//
// Time:  O(n²) — two nested loops, O(1) set ops.
// Space: O(n) — the per-pair set plus dedup set.
func hashSet(nums []int) [][]int {
	n := len(nums)
	seen := make(map[[3]int]bool)
	var result [][]int

	for i := 0; i < n-2; i++ {
		// Set of values seen in the current inner pass nums[i+1..j-1].
		inner := make(map[int]bool)
		for j := i + 1; j < n; j++ {
			complement := -nums[i] - nums[j]
			if inner[complement] {
				triple := [3]int{nums[i], nums[j], complement}
				sort.Ints(triple[:])
				if !seen[triple] {
					seen[triple] = true
					result = append(result, []int{triple[0], triple[1], triple[2]})
				}
			}
			inner[nums[j]] = true
		}
	}
	return result
}

func main() {
	examples := []struct {
		nums   []int
		expect string
	}{
		{[]int{-1, 0, 1, 2, -1, -4}, "[[-1 -1 2] [-1 0 1]]"},
		{[]int{0, 1, 1}, "[]"},
		{[]int{0, 0, 0}, "[[0 0 0]]"},
		{[]int{-2, 0, 0, 2, 2}, "[[-2 0 2]]"},
	}

	approaches := []struct {
		name string
		fn   func([]int) [][]int
	}{
		{"Approach 1: Brute Force      O(n³) T | O(n) S", bruteForce},
		{"Approach 2: Two Pointers   ✅ O(n²) T | O(1) S", twoPointers},
		{"Approach 3: Hash Set        O(n²) T | O(n) S", hashSet},
	}

	for _, ex := range examples {
		fmt.Printf("nums=%v\n", ex.nums)
		for _, ap := range approaches {
			cp := make([]int, len(ex.nums))
			copy(cp, ex.nums)
			fmt.Printf("  %-50s → %v\n", ap.name, ap.fn(cp))
		}
		fmt.Println()
	}
}
