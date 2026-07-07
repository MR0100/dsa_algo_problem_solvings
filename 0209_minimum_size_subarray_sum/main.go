package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Brute Force (Try Every Start, Extend) ────────────────────────
//
// bruteForce solves Minimum Size Subarray Sum by trying every possible start
// index and extending the subarray rightwards until its sum reaches target.
//
// Intuition:
//
//	A subarray is fixed by its two endpoints, so enumerate them. For each
//	start i, grow the end j while accumulating a running sum; because all
//	numbers are positive, the sum only grows as j advances, so the FIRST j
//	where sum ≥ target gives the shortest valid subarray starting at i —
//	stop there and try the next start. (Recomputing each window's sum from
//	scratch would be O(n³); the running sum trims one factor of n.)
//
// Algorithm:
//  1. For every start index i:
//     a. Reset sum to 0.
//     b. Extend j from i to n-1, adding nums[j] to sum.
//     c. The first time sum ≥ target, record length j-i+1 if it beats the
//     best so far, then break (longer j can only do worse for this i).
//  2. Return 0 if no window ever reached the target, else the best length.
//
// Time:  O(n²) — n starts, each extending up to n elements.
// Space: O(1) — a few scalar trackers.
func bruteForce(target int, nums []int) int {
	n := len(nums)
	best := n + 1 // sentinel: longer than any real subarray means "not found yet"

	for i := 0; i < n; i++ { // try every start position
		sum := 0
		for j := i; j < n; j++ { // extend the end one element at a time
			sum += nums[j] // running sum of nums[i..j]
			if sum >= target {
				if j-i+1 < best {
					best = j - i + 1 // shorter qualifying window found
				}
				break // any longer window from this i is worse — next start
			}
		}
	}

	if best == n+1 {
		return 0 // no subarray ever reached target
	}
	return best
}

// ── Approach 2: Prefix Sums + Binary Search ──────────────────────────────────
//
// prefixSumBinarySearch solves Minimum Size Subarray Sum by precomputing
// prefix sums and binary-searching, for each start, the nearest end that
// makes the window sum reach target.
//
// Intuition:
//
//	sum(i..j) = prefix[j+1] - prefix[i]. Since every nums[k] ≥ 1, the prefix
//	array is strictly increasing — and a sorted array invites binary search.
//	For each start i we need the smallest index e > i with
//	prefix[e] ≥ prefix[i] + target; then the window i..e-1 has length e-i.
//	This is the O(n log n) alternative the follow-up asks for.
//
// Algorithm:
//  1. Build prefix[0..n] with prefix[k] = nums[0] + … + nums[k-1].
//  2. For each start i in 0..n-1:
//     a. need := prefix[i] + target.
//     b. Binary-search the smallest e in (i, n] with prefix[e] ≥ need
//     (sort.SearchInts does exactly this lower-bound search).
//     c. If such e exists (e ≤ n), candidate length is e - i.
//  3. Return the smallest candidate, or 0 if none exists.
//
// Time:  O(n log n) — n starts × one O(log n) binary search each.
// Space: O(n) — the prefix-sum array.
func prefixSumBinarySearch(target int, nums []int) int {
	n := len(nums)
	// prefix[k] = sum of the first k elements; strictly increasing since nums[i] ≥ 1.
	prefix := make([]int, n+1)
	for k := 0; k < n; k++ {
		prefix[k+1] = prefix[k] + nums[k]
	}

	best := n + 1 // sentinel meaning "no valid window found yet"
	for i := 0; i < n; i++ {
		need := prefix[i] + target // window i..e-1 works iff prefix[e] ≥ need
		// Lower bound: smallest e with prefix[e] ≥ need (searches the whole
		// array; results e ≤ i are impossible since that window would be empty
		// or negative-length, and prefix[e] < need there anyway).
		e := sort.SearchInts(prefix, need)
		if e <= n { // found a real prefix index → window nums[i..e-1] qualifies
			if e-i < best {
				best = e - i // record the shorter window length
			}
		}
	}

	if best == n+1 {
		return 0 // target unreachable from any start
	}
	return best
}

// ── Approach 3: Sliding Window (Optimal) ─────────────────────────────────────
//
// slidingWindow solves Minimum Size Subarray Sum with a variable-size window
// that expands right to gain sum and shrinks left while the sum suffices.
//
// Intuition:
//
//	All elements are positive, so a window's sum grows when the right edge
//	moves right and shrinks when the left edge moves right — both moves are
//	monotone. That means each pointer only ever needs to move forward:
//	expand right until the window is valid (sum ≥ target), then shrink from
//	the left as far as validity survives, recording the length at every
//	valid moment. No candidate is missed because for each right edge we find
//	the tightest possible left edge.
//
// Algorithm:
//  1. left = 0, sum = 0, best = ∞.
//  2. For right = 0..n-1: add nums[right] to sum.
//  3. While sum ≥ target: record right-left+1 if smaller than best, then
//     remove nums[left] from sum and advance left.
//  4. Return best, or 0 if it never got set.
//
// Time:  O(n) — left and right each advance at most n times total (amortised).
// Space: O(1) — two indices and two accumulators.
func slidingWindow(target int, nums []int) int {
	n := len(nums)
	best := n + 1 // sentinel: "no valid window seen"
	sum := 0      // sum of the current window nums[left..right]
	left := 0     // left edge of the window

	for right := 0; right < n; right++ {
		sum += nums[right] // expand: pull nums[right] into the window
		// Shrink while still valid — finds the tightest window ending at right.
		for sum >= target {
			if right-left+1 < best {
				best = right - left + 1 // new shortest qualifying window
			}
			sum -= nums[left] // expel the leftmost element
			left++            // window's left edge moves right
		}
	}

	if best == n+1 {
		return 0 // total array sum < target — impossible
	}
	return best
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Try Every Start, Extend) ===")
	fmt.Println(bruteForce(7, []int{2, 3, 1, 2, 4, 3}))        // 2
	fmt.Println(bruteForce(4, []int{1, 4, 4}))                 // 1
	fmt.Println(bruteForce(11, []int{1, 1, 1, 1, 1, 1, 1, 1})) // 0

	fmt.Println("=== Approach 2: Prefix Sums + Binary Search ===")
	fmt.Println(prefixSumBinarySearch(7, []int{2, 3, 1, 2, 4, 3}))        // 2
	fmt.Println(prefixSumBinarySearch(4, []int{1, 4, 4}))                 // 1
	fmt.Println(prefixSumBinarySearch(11, []int{1, 1, 1, 1, 1, 1, 1, 1})) // 0

	fmt.Println("=== Approach 3: Sliding Window (Optimal) ===")
	fmt.Println(slidingWindow(7, []int{2, 3, 1, 2, 4, 3}))        // 2
	fmt.Println(slidingWindow(4, []int{1, 4, 4}))                 // 1
	fmt.Println(slidingWindow(11, []int{1, 1, 1, 1, 1, 1, 1, 1})) // 0
}
