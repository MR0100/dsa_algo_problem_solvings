package main

import "fmt"

// ── Approach 1: Brute Force Recursion ────────────────────────────────────────
//
// bruteForce solves Wiggle Subsequence by trying, for every element, whether to
// keep it as the next "up" or "down" wiggle and taking the best over all choices.
//
// Intuition:
//
//	A wiggle subsequence alternates strictly up, down, up, down, ... From a given
//	index, the length we can still build depends on whether the previous kept
//	difference was going up or going down. So define a recursion carrying that
//	direction as state: at each later element either extend (if it points the
//	required way) or skip it.
//
// Algorithm:
//  1. calc(i, isUp): best wiggle length starting after committing element i,
//     where isUp says the next needed difference must be positive.
//  2. Scan j > i: if the parity matches (nums[j] > nums[i] when isUp, or
//     nums[j] < nums[i] when down), recurse with the direction flipped.
//  3. Answer = 1 + max(start with an up move, start with a down move) from index 0.
//
// Time:  O(2^n) — every element branches into keep/skip across both directions.
// Space: O(n) — recursion depth.
func bruteForce(nums []int) int {
	n := len(nums)
	if n == 0 {
		return 0
	}
	// calc returns the longest wiggle subsequence achievable considering
	// indices after i, given the previous committed element index and whether
	// the next difference must be "up".
	var calc func(i int, isUp bool) int
	calc = func(i int, isUp bool) int {
		best := 0 // length contributed by elements after i
		for j := i + 1; j < len(nums); j++ {
			// Only follow j if it continues the required alternation direction.
			if isUp && nums[j] > nums[i] {
				best = max(best, 1+calc(j, false)) // took an up step; next must be down
			} else if !isUp && nums[j] < nums[i] {
				best = max(best, 1+calc(j, true)) // took a down step; next must be up
			}
		}
		return best
	}
	// The first move can be either up or down; add 1 for element 0 itself.
	return 1 + max(calc(0, true), calc(0, false))
}

// ── Approach 2: Dynamic Programming (up/down tables) ─────────────────────────
//
// dpTables solves Wiggle Subsequence by tracking, at each index, the longest
// wiggle ending there with a rising last step (up) and with a falling last step.
//
// Intuition:
//
//	up[i]   = longest wiggle subsequence ending at i whose last difference is +.
//	down[i] = same but whose last difference is −.
//	To extend to i we look at every earlier j: if nums[i] > nums[j] we can append
//	i after a "down"-ending sequence to make an up step; symmetrically for down.
//
// Algorithm:
//  1. up[i] = max over j<i with nums[i]>nums[j] of down[j]+1 (≥1).
//  2. down[i] = max over j<i with nums[i]<nums[j] of up[j]+1 (≥1).
//  3. Answer = max(up[n-1], down[n-1]).
//
// Time:  O(n²) — nested scan of earlier elements.
// Space: O(n) — the two tables.
func dpTables(nums []int) int {
	n := len(nums)
	if n == 0 {
		return 0
	}
	up := make([]int, n)   // best wiggle ending at i with a final up step
	down := make([]int, n) // best wiggle ending at i with a final down step
	for i := range nums {
		up[i], down[i] = 1, 1 // a single element is a valid length-1 wiggle
	}
	for i := 1; i < n; i++ {
		for j := 0; j < i; j++ {
			if nums[i] > nums[j] {
				up[i] = max(up[i], down[j]+1) // rise onto a sequence that fell into j
			} else if nums[i] < nums[j] {
				down[i] = max(down[i], up[j]+1) // fall onto a sequence that rose into j
			}
		}
	}
	return max(up[n-1], down[n-1])
}

// ── Approach 3: Linear DP (rolling up/down) ──────────────────────────────────
//
// dpLinear solves Wiggle Subsequence keeping only the best up- and down-ending
// lengths seen so far, updated as we sweep left to right.
//
// Intuition:
//
//	The O(n²) tables are wasteful: when nums[i] > nums[i-1], the only useful
//	previous state is the best down-ending length, and it can only get better as
//	i advances. So collapse the arrays into two running scalars.
//
// Algorithm:
//  1. up, down = 1, 1.
//  2. If nums[i] > nums[i-1]: up = down + 1.
//  3. Else if nums[i] < nums[i-1]: down = up + 1.
//  4. Equal elements change nothing. Answer = max(up, down).
//
// Time:  O(n) — one pass.
// Space: O(1) — two scalars.
func dpLinear(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	up, down := 1, 1 // best wiggle so far ending with an up / down step
	for i := 1; i < len(nums); i++ {
		if nums[i] > nums[i-1] {
			up = down + 1 // an ascent extends the best descending wiggle
		} else if nums[i] < nums[i-1] {
			down = up + 1 // a descent extends the best ascending wiggle
		}
		// equal → neither direction advances
	}
	return max(up, down)
}

// ── Approach 4: Greedy Slope Counting (Optimal) ──────────────────────────────
//
// greedy solves Wiggle Subsequence by counting how many times the sign of the
// consecutive difference flips — each flip is a genuine wiggle turning point.
//
// Intuition:
//
//	The longest wiggle subsequence uses exactly the "peaks and valleys" of the
//	array. Walk through consecutive differences; every time the difference
//	switches from positive to negative (or vice versa) we have found a new
//	turning point, so bump the count. Flat stretches are ignored.
//
// Algorithm:
//  1. count = 1, prevDiff = 0.
//  2. For each i≥1: diff = nums[i]-nums[i-1].
//  3. If (diff>0 && prevDiff<=0) or (diff<0 && prevDiff>=0): count++, prevDiff=diff.
//  4. Return count.
//
// Time:  O(n) — single pass.
// Space: O(1).
func greedy(nums []int) int {
	if len(nums) < 2 {
		return len(nums) // 0 or 1 element is already the answer
	}
	count := 1    // the first element always starts a wiggle
	prevDiff := 0 // sign of the last accepted difference (0 = none yet)
	for i := 1; i < len(nums); i++ {
		diff := nums[i] - nums[i-1]
		// A rising step counts only if the previous accepted step was not rising,
		// and symmetrically for a falling step — that is exactly a direction flip.
		if (diff > 0 && prevDiff <= 0) || (diff < 0 && prevDiff >= 0) {
			count++
			prevDiff = diff // remember the new direction
		}
	}
	return count
}

func main() {
	fmt.Println("=== Approach 1: Brute Force Recursion ===")
	fmt.Println(bruteForce([]int{1, 7, 4, 9, 2, 5}))                   // expected 6
	fmt.Println(bruteForce([]int{1, 17, 5, 10, 13, 15, 10, 5, 16, 8})) // expected 7
	fmt.Println(bruteForce([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}))          // expected 2

	fmt.Println("=== Approach 2: Dynamic Programming (up/down tables) ===")
	fmt.Println(dpTables([]int{1, 7, 4, 9, 2, 5}))                   // expected 6
	fmt.Println(dpTables([]int{1, 17, 5, 10, 13, 15, 10, 5, 16, 8})) // expected 7
	fmt.Println(dpTables([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}))          // expected 2

	fmt.Println("=== Approach 3: Linear DP (rolling up/down) ===")
	fmt.Println(dpLinear([]int{1, 7, 4, 9, 2, 5}))                   // expected 6
	fmt.Println(dpLinear([]int{1, 17, 5, 10, 13, 15, 10, 5, 16, 8})) // expected 7
	fmt.Println(dpLinear([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}))          // expected 2

	fmt.Println("=== Approach 4: Greedy Slope Counting (Optimal) ===")
	fmt.Println(greedy([]int{1, 7, 4, 9, 2, 5}))                   // expected 6
	fmt.Println(greedy([]int{1, 17, 5, 10, 13, 15, 10, 5, 16, 8})) // expected 7
	fmt.Println(greedy([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}))          // expected 2
}
