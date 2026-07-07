package main

import (
	"fmt"
	"sort"
	"strconv"
)

// ── Approach 1: Brute Force (Convert + String Sort) ──────────────────────────
//
// bruteForce solves Lexicographical Numbers by generating [1..n], converting
// every number to its decimal string, sorting the strings lexicographically,
// then converting back to ints.
//
// Intuition:
//
//	"Lexicographic order" is literally dictionary order of the decimal strings.
//	So the most direct thing that could possibly work is: make the strings,
//	sort them as strings, parse them back. It is O(n log n) and O(n) extra —
//	the problem asks for O(n) time / O(1) space, so this is a baseline only.
//
// Algorithm:
//  1. Build nums = [1, 2, ..., n].
//  2. Sort by comparing the decimal STRINGS (not the numeric values).
//  3. Return nums.
//
// Time:  O(n · L · log n) — sort does O(n log n) comparisons, each comparing
//
//	two strings of length up to L = O(log n) digits.
//
// Space: O(n · L) — we hold n strings (or n ints plus scratch) at once.
func bruteForce(n int) []int {
	nums := make([]int, n) // will hold 1..n
	for i := 0; i < n; i++ {
		nums[i] = i + 1 // fill with the values 1..n
	}
	// Sort by the DECIMAL STRING of each value, i.e. dictionary order:
	// "10" < "2" because '1' < '2' at the first character.
	sort.Slice(nums, func(a, b int) bool {
		return strconv.Itoa(nums[a]) < strconv.Itoa(nums[b])
	})
	return nums
}

// ── Approach 2: Preorder DFS on the 10-ary Digit Tree ────────────────────────
//
// dfs solves Lexicographical Numbers by walking the trie of numbers: the root
// has children 1..9, and every node x has children 10x .. 10x+9. A preorder
// traversal of this tree visits numbers in exactly lexicographic order.
//
// Intuition:
//
//	Picture numbers as paths in a tree keyed by digits. From "1" you can go
//	deeper to "10", "11", ..., "19" (append a digit) before you ever move to
//	the sibling "2". That "go deep, append 0..9, then move to next sibling"
//	order IS dictionary order. A depth-first preorder visit prints them
//	correctly. We just prune any branch whose value exceeds n.
//
// Algorithm:
//  1. For each starting digit d in 1..9: DFS(d).
//  2. DFS(cur): if cur > n return; append cur; then for next in 0..9 recurse
//     DFS(cur*10 + next).
//
// Time:  O(n) — every value in [1..n] is appended exactly once; the pruned
//
//	recursive calls are bounded by a constant factor per emitted node.
//
// Space: O(log n) — recursion depth equals the number of digits of n (the
//
//	output slice itself is required output, not counted as auxiliary).
func dfs(n int) []int {
	result := make([]int, 0, n) // preallocate for the n emitted values

	// visit performs the preorder walk rooted at cur.
	var visit func(cur int)
	visit = func(cur int) {
		if cur > n { // cur (and all its 10x children) exceed n → prune
			return
		}
		result = append(result, cur) // preorder: record before descending
		// Children of cur are cur*10 + 0 .. cur*10 + 9, in ascending digit order.
		for next := 0; next <= 9; next++ {
			visit(cur*10 + next) // append the next digit and go deeper
		}
	}

	for d := 1; d <= 9; d++ { // the tree has 9 roots: 1..9 (no leading zero)
		visit(d)
	}
	return result
}

// ── Approach 3: Iterative Successor (Optimal, O(n) time / O(1) space) ─────────
//
// iterativeNext solves Lexicographical Numbers by starting at 1 and repeatedly
// computing the lexicographic SUCCESSOR in place, using only the current value.
//
// Intuition:
//
//	Given the current number, the next number in dictionary order is found by
//	three rules that mirror the DFS moves without a stack:
//	  1. Try to go DEEPER: cur*10. If cur*10 <= n, that is next (append a 0).
//	  2. Otherwise go to a SIBLING/ancestor: while the current path ends in 9
//	     (cur % 10 == 9) or cur+1 would exceed n, climb up one level (cur /= 10),
//	     then increment: cur++.
//	This reproduces preorder traversal with O(1) memory.
//
// Algorithm:
//  1. cur = 1; repeat n times: append cur, then advance cur to its successor.
//  2. Advance: if cur*10 <= n then cur *= 10; else { while cur%10==9 ||
//     cur+1 > n: cur /= 10; cur++ }.
//
// Time:  O(n) — n iterations; the inner "climb" is amortised O(1) per step.
// Space: O(1) — only the counter cur (output slice is required output).
func iterativeNext(n int) []int {
	result := make([]int, 0, n) // n values will be produced
	cur := 1                    // lexicographically smallest positive number
	for i := 0; i < n; i++ {
		result = append(result, cur) // emit current number
		if cur*10 <= n {             // rule 1: can we append a '0' and stay ≤ n?
			cur *= 10 // go deeper: e.g. 1 -> 10
		} else {
			// rule 2: cannot go deeper; move right, climbing when blocked.
			// Climb while the last digit is 9 (no right sibling) OR
			// incrementing would exceed n (sibling out of range).
			for cur%10 == 9 || cur+1 > n {
				cur /= 10 // step up to the parent
			}
			cur++ // move to the next sibling at this (possibly higher) level
		}
	}
	return result
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (String Sort) ===")
	fmt.Println(bruteForce(13)) // expected [1 10 11 12 13 2 3 4 5 6 7 8 9]
	fmt.Println(bruteForce(2))  // expected [1 2]

	fmt.Println("=== Approach 2: Preorder DFS ===")
	fmt.Println(dfs(13)) // expected [1 10 11 12 13 2 3 4 5 6 7 8 9]
	fmt.Println(dfs(2))  // expected [1 2]

	fmt.Println("=== Approach 3: Iterative Successor (Optimal) ===")
	fmt.Println(iterativeNext(13)) // expected [1 10 11 12 13 2 3 4 5 6 7 8 9]
	fmt.Println(iterativeNext(2))  // expected [1 2]
}
