package main

import "fmt"

// ── Approach 1: Brute Force (Re-scan From Each Start) ─────────────────────────
//
// bruteForce tries every index as the start of a run of 1's and counts how far
// the run extends, tracking the longest.
//
// Intuition:
//
//	The definition of the answer is "the longest stretch of consecutive 1's".
//	The most literal way to find it is: from each position i, if nums[i] is 1,
//	walk forward counting 1's until a 0 (or the end), and remember the biggest
//	such count. It re-examines characters many times but needs no cleverness.
//
// Algorithm:
//  1. best = 0.
//  2. For each start i: run = 0; extend j from i while nums[j] == 1, run++.
//     Update best = max(best, run).
//  3. Return best.
//
// Time:  O(n²) worst case — an all-1's array makes every start scan to the end.
// Space: O(1).
func bruteForce(nums []int) int {
	best := 0
	for i := 0; i < len(nums); i++ {
		run := 0
		for j := i; j < len(nums) && nums[j] == 1; j++ { // extend the run from i
			run++
		}
		if run > best {
			best = run // remember the longest run seen so far
		}
	}
	return best
}

// ── Approach 2: Single-Pass Running Count (Optimal) ──────────────────────────
//
// runningCount sweeps once, keeping a counter of the current consecutive-1's
// streak that resets on every 0.
//
// Intuition:
//
//	Consecutive 1's are contiguous, so a single left-to-right pass suffices:
//	maintain `cur`, the length of the streak ending at the current index. Each
//	1 extends the streak (cur++); each 0 breaks it (cur = 0). The answer is the
//	maximum value `cur` ever reaches. No re-scanning — every element is touched
//	once.
//
// Algorithm:
//  1. best = 0, cur = 0.
//  2. For each x in nums: if x == 1, cur++ and best = max(best, cur); else cur = 0.
//  3. Return best.
//
// Time:  O(n) — one pass.
// Space: O(1) — two counters.
func runningCount(nums []int) int {
	best, cur := 0, 0
	for _, x := range nums {
		if x == 1 {
			cur++ // extend the current streak of 1's
			if cur > best {
				best = cur // new longest streak
			}
		} else {
			cur = 0 // a 0 breaks the streak; start counting fresh
		}
	}
	return best
}

// ── Approach 3: Sliding Window (Framing for the Follow-Up) ────────────────────
//
// slidingWindow expresses the same O(n) scan as a window [left, right] that
// never contains a 0: whenever a 0 appears, the left edge jumps past it.
//
// Intuition:
//
//	Frame the streak as a window of all-1's. `right` advances through the
//	array; whenever nums[right] is 0, the window can hold no 1's ending here,
//	so snap `left` to right+1 (empty the window). Otherwise the window
//	[left..right] is a valid block of 1's whose length right-left+1 is a
//	candidate answer. This is deliberately overkill for the base problem, but
//	it is the exact template that generalises to "Max Consecutive Ones II/III"
//	(allow up to k zeros) — there you keep the window valid by shrinking left
//	only while it contains more than k zeros.
//
// Algorithm:
//  1. left = 0, best = 0.
//  2. For right = 0..n-1: if nums[right] == 0, set left = right+1 (reset window);
//     else best = max(best, right-left+1).
//  3. Return best.
//
// Time:  O(n) — left and right each move forward monotonically.
// Space: O(1).
func slidingWindow(nums []int) int {
	left, best := 0, 0
	for right := 0; right < len(nums); right++ {
		if nums[right] == 0 {
			left = right + 1 // a 0 can't be in an all-1's window → restart after it
		} else if right-left+1 > best {
			best = right - left + 1 // current all-1's window is the longest yet
		}
	}
	return best
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Re-scan From Each Start) ===")
	fmt.Printf("nums=[1 1 0 1 1 1]  got=%d  expected 3\n", bruteForce([]int{1, 1, 0, 1, 1, 1})) // 3
	fmt.Printf("nums=[1 0 1 1 0 1]  got=%d  expected 2\n", bruteForce([]int{1, 0, 1, 1, 0, 1})) // 2
	fmt.Printf("nums=[0 0 0]        got=%d  expected 0\n", bruteForce([]int{0, 0, 0}))          // 0
	fmt.Printf("nums=[1 1 1 1]      got=%d  expected 4\n", bruteForce([]int{1, 1, 1, 1}))       // 4

	fmt.Println("=== Approach 2: Single-Pass Running Count (Optimal) ===")
	fmt.Printf("nums=[1 1 0 1 1 1]  got=%d  expected 3\n", runningCount([]int{1, 1, 0, 1, 1, 1}))
	fmt.Printf("nums=[1 0 1 1 0 1]  got=%d  expected 2\n", runningCount([]int{1, 0, 1, 1, 0, 1}))
	fmt.Printf("nums=[0 0 0]        got=%d  expected 0\n", runningCount([]int{0, 0, 0}))
	fmt.Printf("nums=[1 1 1 1]      got=%d  expected 4\n", runningCount([]int{1, 1, 1, 1}))

	fmt.Println("=== Approach 3: Sliding Window (Framing for the Follow-Up) ===")
	fmt.Printf("nums=[1 1 0 1 1 1]  got=%d  expected 3\n", slidingWindow([]int{1, 1, 0, 1, 1, 1}))
	fmt.Printf("nums=[1 0 1 1 0 1]  got=%d  expected 2\n", slidingWindow([]int{1, 0, 1, 1, 0, 1}))
	fmt.Printf("nums=[0 0 0]        got=%d  expected 0\n", slidingWindow([]int{0, 0, 0}))
	fmt.Printf("nums=[1 1 1 1]      got=%d  expected 4\n", slidingWindow([]int{1, 1, 1, 1}))
}
