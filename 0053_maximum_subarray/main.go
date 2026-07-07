package main

import (
	"fmt"
	"math"
)

// ── Approach 1: Brute Force ──────────────────────────────────────────────────
//
// bruteForce solves Maximum Subarray by checking all subarrays.
//
// Intuition:
//   Try every pair (i, j) as subarray bounds, compute the sum, track the max.
//
// Algorithm:
//   for i=0 to n-1:
//     sum = 0
//     for j=i to n-1:
//       sum += nums[j]; max = max(max, sum)
//
// Time:  O(n²) — O(n) starting points × O(n) endpoint scan.
// Space: O(1)
func bruteForce(nums []int) int {
	best := math.MinInt64
	for i := 0; i < len(nums); i++ {
		sum := 0
		for j := i; j < len(nums); j++ {
			sum += nums[j]
			if sum > best {
				best = sum
			}
		}
	}
	return best
}

// ── Approach 2: Kadane's Algorithm ───────────────────────────────────────────
//
// kadane solves Maximum Subarray using Kadane's greedy approach.
//
// Intuition:
//   At each position, decide: extend the current subarray or start fresh from
//   this element? If the running sum drops below 0, it can only drag down
//   future sums, so reset it to 0 (start a new subarray from the next element).
//
// Algorithm:
//   curSum = 0; best = nums[0]
//   for each num:
//     curSum += num
//     best = max(best, curSum)
//     if curSum < 0: curSum = 0
//
// Time:  O(n) — single pass.
// Space: O(1)
func kadane(nums []int) int {
	best := nums[0]
	curSum := 0
	for _, num := range nums {
		curSum += num
		if curSum > best {
			best = curSum
		}
		if curSum < 0 {
			curSum = 0 // discard negative prefix; start fresh
		}
	}
	return best
}

// ── Approach 3: Divide and Conquer ───────────────────────────────────────────
//
// divideAndConquer solves Maximum Subarray using divide & conquer.
//
// Intuition:
//   Split the array at mid. The maximum subarray either lies entirely in the
//   left half, entirely in the right half, or crosses the midpoint.
//   The crossing subarray can be found by expanding from mid outward in both
//   directions and summing the best left and right arms.
//
// Algorithm:
//   maxSubArray(lo, hi):
//     if lo==hi: return nums[lo]
//     mid = (lo+hi)/2
//     left  = maxSubArray(lo, mid)
//     right = maxSubArray(mid+1, hi)
//     cross = maxCrossing(lo, mid, hi)
//     return max(left, right, cross)
//
// Time:  O(n log n) — T(n) = 2T(n/2) + O(n).
// Space: O(log n)   — recursion stack depth.
func divideAndConquer(nums []int) int {
	return dac(nums, 0, len(nums)-1)
}

func dac(nums []int, lo, hi int) int {
	if lo == hi {
		return nums[lo]
	}
	mid := (lo + hi) / 2
	left := dac(nums, lo, mid)
	right := dac(nums, mid+1, hi)
	cross := maxCross(nums, lo, mid, hi)
	return max3(left, right, cross)
}

// maxCross computes the max subarray sum that crosses mid.
func maxCross(nums []int, lo, mid, hi int) int {
	leftSum := math.MinInt64
	sum := 0
	for i := mid; i >= lo; i-- { // expand left from mid
		sum += nums[i]
		if sum > leftSum {
			leftSum = sum
		}
	}
	rightSum := math.MinInt64
	sum = 0
	for i := mid + 1; i <= hi; i++ { // expand right from mid+1
		sum += nums[i]
		if sum > rightSum {
			rightSum = sum
		}
	}
	return leftSum + rightSum
}

func max3(a, b, c int) int {
	if a >= b && a >= c {
		return a
	}
	if b >= c {
		return b
	}
	return c
}

// ── Approach 4: Dynamic Programming ──────────────────────────────────────────
//
// dpBottomUp solves Maximum Subarray with explicit DP array.
//
// Intuition:
//   dp[i] = max subarray sum ending at index i.
//   dp[i] = max(nums[i], dp[i-1] + nums[i])  — extend previous or start fresh.
//   Answer = max of all dp[i].
//
// Time:  O(n)
// Space: O(n)  — can be reduced to O(1) like Kadane's.
func dpBottomUp(nums []int) int {
	n := len(nums)
	dp := make([]int, n)
	dp[0] = nums[0]
	best := dp[0]
	for i := 1; i < n; i++ {
		if dp[i-1]+nums[i] > nums[i] {
			dp[i] = dp[i-1] + nums[i]
		} else {
			dp[i] = nums[i]
		}
		if dp[i] > best {
			best = dp[i]
		}
	}
	return best
}

func main() {
	nums1 := []int{-2, 1, -3, 4, -1, 2, 1, -5, 4}
	nums2 := []int{1}
	nums3 := []int{5, 4, -1, 7, 8}

	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Printf("nums=%v  got=%d  expected 6\n", nums1, bruteForce(nums1))
	fmt.Printf("nums=%v  got=%d  expected 1\n", nums2, bruteForce(nums2))
	fmt.Printf("nums=%v  got=%d  expected 23\n", nums3, bruteForce(nums3))

	fmt.Println("=== Approach 2: Kadane's Algorithm ===")
	fmt.Printf("nums=%v  got=%d  expected 6\n", nums1, kadane(nums1))
	fmt.Printf("nums=%v  got=%d  expected 1\n", nums2, kadane(nums2))
	fmt.Printf("nums=%v  got=%d  expected 23\n", nums3, kadane(nums3))

	fmt.Println("=== Approach 3: Divide and Conquer ===")
	fmt.Printf("nums=%v  got=%d  expected 6\n", nums1, divideAndConquer(nums1))
	fmt.Printf("nums=%v  got=%d  expected 1\n", nums2, divideAndConquer(nums2))
	fmt.Printf("nums=%v  got=%d  expected 23\n", nums3, divideAndConquer(nums3))

	fmt.Println("=== Approach 4: DP Bottom-Up ===")
	fmt.Printf("nums=%v  got=%d  expected 6\n", nums1, dpBottomUp(nums1))
	fmt.Printf("nums=%v  got=%d  expected 1\n", nums2, dpBottomUp(nums2))
	fmt.Printf("nums=%v  got=%d  expected 23\n", nums3, dpBottomUp(nums3))
}
