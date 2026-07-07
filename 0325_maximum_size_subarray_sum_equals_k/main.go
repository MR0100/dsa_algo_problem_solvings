package main

import "fmt"

// ── Approach 1: Brute Force (All Subarrays) ──────────────────────────────────
//
// bruteForce solves Maximum Size Subarray Sum Equals k by checking every
// contiguous subarray, tracking the length of the longest one whose sum is k.
//
// Intuition:
//
//	There are O(n^2) subarrays. Fix a start i, extend the end j while keeping a
//	running sum; whenever the sum equals k, the subarray i..j is a candidate and
//	its length j-i+1 competes for the maximum.
//
// Algorithm:
//  1. best = 0.
//  2. For each start i: sum = 0; for each end j >= i: sum += nums[j];
//     if sum == k, best = max(best, j-i+1).
//  3. Return best.
//
// Time:  O(n^2) — every (start, end) pair.
// Space: O(1) — a couple of scalars.
func bruteForce(nums []int, k int) int {
	best := 0
	for i := 0; i < len(nums); i++ {
		sum := 0 // running sum of subarray starting at i
		for j := i; j < len(nums); j++ {
			sum += nums[j] // extend the window to include nums[j]
			if sum == k && j-i+1 > best {
				best = j - i + 1 // longer valid subarray found
			}
		}
	}
	return best
}

// ── Approach 2: Prefix Sum + Hash Map (Optimal) ──────────────────────────────
//
// prefixSumHashMap solves Maximum Size Subarray Sum Equals k in one pass by
// storing the EARLIEST index at which each prefix sum first appeared; a subarray
// summing to k corresponds to two prefix sums differing by k.
//
// Intuition:
//
//	Let P(i) = nums[0]+...+nums[i-1] be the prefix sum before index i. A subarray
//	j..i has sum P(i+1) - P(j). We want that to equal k, i.e. P(j) = P(i+1) - k.
//	If we know the FIRST index where prefix sum (current - k) occurred, the
//	subarray from there to here is the longest ending here. Record each prefix
//	sum's earliest index only (never overwrite) to maximise length.
//
// Algorithm:
//  1. seen = {0: -1} — prefix sum 0 exists "before" index 0.
//  2. sum = 0, best = 0.
//  3. For i, x in nums: sum += x; if (sum-k) in seen, best =
//     max(best, i - seen[sum-k]); if sum not in seen, seen[sum] = i.
//  4. Return best.
//
// Time:  O(n) — single pass with O(1) hash operations.
// Space: O(n) — up to n distinct prefix sums stored.
func prefixSumHashMap(nums []int, k int) int {
	seen := map[int]int{0: -1} // prefix sum → earliest index it appeared at
	sum := 0                   // running prefix sum
	best := 0                  // longest valid subarray length
	for i, x := range nums {
		sum += x // prefix sum through index i (inclusive)
		// If some earlier prefix equals sum-k, the gap between them sums to k.
		if j, ok := seen[sum-k]; ok && i-j > best {
			best = i - j // length of that subarray
		}
		// Store only the FIRST time we see this prefix sum → longest reach later.
		if _, ok := seen[sum]; !ok {
			seen[sum] = i
		}
	}
	return best
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(bruteForce([]int{1, -1, 5, -2, 3}, 3)) // expected 4
	fmt.Println(bruteForce([]int{-2, -1, 2, 1}, 1))    // expected 2
	fmt.Println(bruteForce([]int{1, 2, 3}, 7))         // expected 0

	fmt.Println("=== Approach 2: Prefix Sum + Hash Map (Optimal) ===")
	fmt.Println(prefixSumHashMap([]int{1, -1, 5, -2, 3}, 3)) // expected 4
	fmt.Println(prefixSumHashMap([]int{-2, -1, 2, 1}, 1))    // expected 2
	fmt.Println(prefixSumHashMap([]int{1, 2, 3}, 7))         // expected 0
}
