package main

import "fmt"

// ── Approach 1: Brute Force (try flipping each zero) ──────────────────────────
//
// bruteForce solves Max Consecutive Ones II by, for every position that could
// be the single flipped zero, measuring the run of 1s obtainable there.
//
// Intuition:
//
//	We may flip at most one 0. So the answer is either the longest existing run
//	of 1s (flip nothing), or — for some particular zero — the length obtained
//	by turning THAT zero into a 1 and joining the 1s on its left with the 1s on
//	its right. Try every zero as the flip target and take the maximum.
//
// Algorithm:
//  1. For each index i:
//     - count ones extending LEFT from i (stop at the first 0 to the left).
//     - count ones extending RIGHT from i (stop at the first 0 to the right).
//     - if nums[i] == 0, candidate = left + 1 + right (we spent our flip on i);
//     if nums[i] == 1, candidate = left + right + ... but a pure run is
//     already covered, so we only need the zero-flip case plus the all-ones case.
//  2. Track the running maximum; also handle the "no zero at all" array.
//
// Time:  O(n²) — each index scans left and right, O(n) per index.
// Space: O(1) — only counters.
func bruteForce(nums []int) int {
	n := len(nums)
	best := 0
	// First, the longest run with no flip (covers arrays of all 1s).
	run := 0
	for _, v := range nums {
		if v == 1 {
			run++ // extend the current run of ones
			if run > best {
				best = run
			}
		} else {
			run = 0 // a zero ends the flip-free run
		}
	}
	// Now try spending the single flip on each zero.
	for i := 0; i < n; i++ {
		if nums[i] != 0 {
			continue // only zeros are worth flipping
		}
		left := 0
		for j := i - 1; j >= 0 && nums[j] == 1; j-- {
			left++ // ones immediately to the left of i
		}
		right := 0
		for j := i + 1; j < n && nums[j] == 1; j++ {
			right++ // ones immediately to the right of i
		}
		if cand := left + 1 + right; cand > best { // +1 for the flipped zero itself
			best = cand
		}
	}
	return best
}

// ── Approach 2: Previous-Count DP (track length ending here) ──────────────────
//
// prevCountDP solves Max Consecutive Ones II in one pass by remembering, at
// each index, the length of the 1-run ending here WITH and WITHOUT a flip used.
//
// Intuition:
//
//	Scan left to right keeping two running lengths:
//	  cur  = length of consecutive 1s ending at i having used NO flip,
//	  prev = length of consecutive 1s ending at i having used ONE flip.
//	On a 1: both runs extend. On a 0: the "no flip" run resets to 0, and the
//	"one flip used" run becomes cur+1 (we flip THIS zero, gluing onto the run
//	that had no flip). The answer is the largest prev seen.
//
// Algorithm:
//  1. cur = prev = 0.
//  2. For each v: if v == 1, cur++ and prev++ (flip still available or already spent).
//     else prev = cur + 1 (flip this zero), cur = 0.
//  3. Answer = max prev over the scan.
//
// Time:  O(n) — single pass.
// Space: O(1) — two counters.
func prevCountDP(nums []int) int {
	cur := 0  // run of 1s ending here with NO flip used
	prev := 0 // run of 1s ending here with exactly ONE flip used
	best := 0
	for _, v := range nums {
		if v == 1 {
			cur++  // extend the flip-free run
			prev++ // the flipped run also grows over a real 1
		} else {
			prev = cur + 1 // flip THIS zero: glue onto the flip-free run + itself
			cur = 0        // the flip-free run is broken by a genuine 0
		}
		if prev > best {
			best = prev // prev always ≥ cur, so it alone bounds the answer
		}
	}
	return best
}

// ── Approach 3: Sliding Window (Optimal, streaming-friendly) ──────────────────
//
// slidingWindow solves Max Consecutive Ones II by keeping a window that
// contains at most one zero and reporting its greatest width.
//
// Intuition:
//
//	A valid answer is exactly a contiguous window we can turn into all-1s by
//	flipping ≤ 1 zero, i.e. a window holding at most ONE zero. Grow the window
//	to the right; whenever it would contain a second zero, shrink from the left
//	just past the older zero. The widest such window is the answer. This never
//	needs to look backwards, so it also answers the streaming follow-up: keep
//	only `left` and the index of the last zero.
//
// Algorithm:
//  1. left = 0, zeros = 0, best = 0.
//  2. For right = 0..n-1:
//     - if nums[right] == 0, zeros++.
//     - while zeros > 1: if nums[left] == 0 decrement zeros; left++.
//     - best = max(best, right − left + 1).
//  3. Return best.
//
// Time:  O(n) — each index enters and leaves the window at most once.
// Space: O(1) — two indices and a zero counter (stream-friendly).
func slidingWindow(nums []int) int {
	left := 0  // left edge of the current window
	zeros := 0 // number of zeros currently inside the window
	best := 0
	for right := 0; right < len(nums); right++ {
		if nums[right] == 0 {
			zeros++ // the new element on the right is a zero
		}
		// If we now hold two zeros, advance left until only one remains.
		for zeros > 1 {
			if nums[left] == 0 {
				zeros-- // the zero leaving on the left frees our single flip
			}
			left++ // shrink the window from the left
		}
		if w := right - left + 1; w > best {
			best = w // widest window with ≤ 1 zero seen so far
		}
	}
	return best
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(bruteForce([]int{1, 0, 1, 1, 0}))    // expected 4
	fmt.Println(bruteForce([]int{1, 0, 1, 1, 0, 1})) // expected 4
	fmt.Println(bruteForce([]int{0, 0, 0}))          // expected 1
	fmt.Println(bruteForce([]int{1, 1, 1}))          // expected 3

	fmt.Println("=== Approach 2: Previous-Count DP ===")
	fmt.Println(prevCountDP([]int{1, 0, 1, 1, 0}))    // expected 4
	fmt.Println(prevCountDP([]int{1, 0, 1, 1, 0, 1})) // expected 4
	fmt.Println(prevCountDP([]int{0, 0, 0}))          // expected 1
	fmt.Println(prevCountDP([]int{1, 1, 1}))          // expected 3

	fmt.Println("=== Approach 3: Sliding Window (Optimal) ===")
	fmt.Println(slidingWindow([]int{1, 0, 1, 1, 0}))    // expected 4
	fmt.Println(slidingWindow([]int{1, 0, 1, 1, 0, 1})) // expected 4
	fmt.Println(slidingWindow([]int{0, 0, 0}))          // expected 1
	fmt.Println(slidingWindow([]int{1, 1, 1}))          // expected 3
}
