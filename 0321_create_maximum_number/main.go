package main

import "fmt"

// ── Approach 1: Merge of Best Sub-sequences (Optimal) ────────────────────────
//
// createMaxNumber solves Create Maximum Number by, for every split of k into
// (i digits taken from nums1, k-i digits taken from nums2), building the
// lexicographically-largest length-i sub-sequence of nums1 and the largest
// length-(k-i) sub-sequence of nums2, greedily merging them into one number,
// and keeping the largest number over all splits.
//
// Intuition:
//
//	Two independent decisions collapse into one loop. Within a single array the
//	best length-t sub-sequence (order preserved) is found with a monotonic
//	stack: pop smaller trailing digits while we still have enough digits left to
//	refill the stack. Across the two arrays we must interleave, and at each step
//	we take from whichever remaining suffix is lexicographically larger — ties
//	are broken by looking further ahead, which is exactly a suffix comparison.
//	Trying all i from 0..k and merging gives the global maximum.
//
// Algorithm:
//  1. For i in max(0, k-n) .. min(k, m):
//     a. best1 = maxSubsequence(nums1, i)
//     b. best2 = maxSubsequence(nums2, k-i)
//     c. candidate = merge(best1, best2)
//     d. keep candidate if greater than the current best.
//  2. Return the best candidate.
//
// Time:  O(k * (m + n)^2) — k splits, each merge is O((m+n)^2) worst case
//
//	because a merge tie triggers a suffix comparison.
//
// Space: O(m + n) — the sub-sequences and the merged candidate.
func createMaxNumber(nums1 []int, nums2 []int, k int) []int {
	m, n := len(nums1), len(nums2) // lengths of the two source arrays
	var best []int                 // best answer found so far
	// i = how many digits we draw from nums1; the rest (k-i) come from nums2.
	// i cannot exceed m, and k-i cannot exceed n → i >= k-n.
	start := 0
	if k-n > 0 {
		start = k - n // ensure k-i <= n
	}
	end := k
	if m < k {
		end = m // ensure i <= m
	}
	for i := start; i <= end; i++ {
		sub1 := maxSubsequence(nums1, i)   // largest length-i pick from nums1
		sub2 := maxSubsequence(nums2, k-i) // largest length-(k-i) pick from nums2
		cand := merge(sub1, sub2)          // interleave into one number
		if greater(cand, 0, best, 0) {     // compare full candidates
			best = cand // adopt the larger number
		}
	}
	return best
}

// maxSubsequence returns the lexicographically-largest length-t sub-sequence of
// nums with the original relative order preserved, using a monotonic stack.
//
// Time:  O(len(nums)) — each element is pushed and popped at most once.
// Space: O(t) — the resulting stack.
func maxSubsequence(nums []int, t int) []int {
	stack := make([]int, 0, t) // will hold the chosen digits (a decreasing-ish stack)
	drop := len(nums) - t      // how many digits we are still allowed to discard
	for _, x := range nums {
		// While the top is smaller than x and we can still afford to drop it,
		// pop it so a bigger digit takes a more significant position.
		for len(stack) > 0 && drop > 0 && stack[len(stack)-1] < x {
			stack = stack[:len(stack)-1] // discard the smaller trailing digit
			drop--                       // spent one of our allowed drops
		}
		stack = append(stack, x) // push current digit
	}
	return stack[:t] // keep only the first t (extra pushes stay at the tail)
}

// merge interleaves a and b into the lexicographically-largest sequence while
// preserving the internal order of each. At every step it takes from whichever
// remaining suffix is larger (compared element by element, deeper on ties).
//
// Time:  O((len(a)+len(b))^2) worst case — a tie forces a suffix comparison.
// Space: O(len(a)+len(b)) — the merged output.
func merge(a []int, b []int) []int {
	out := make([]int, 0, len(a)+len(b)) // merged result
	i, j := 0, 0                         // read cursors into a and b
	for i < len(a) || j < len(b) {
		// greater decides which suffix is larger starting at (i in a) vs (j in b).
		if greater(a, i, b, j) {
			out = append(out, a[i]) // take from a
			i++
		} else {
			out = append(out, b[j]) // take from b
			j++
		}
	}
	return out
}

// greater reports whether the suffix a[i:] is lexicographically greater than the
// suffix b[j:]. A longer suffix wins if it is an extension of the shorter one
// (i.e. the one that still has digits after the other is exhausted is larger).
//
// Time:  O(len(a)+len(b)) — walks both suffixes.
// Space: O(1).
func greater(a []int, i int, b []int, j int) bool {
	for i < len(a) && j < len(b) && a[i] == b[j] {
		i++ // skip equal prefixes
		j++
	}
	// If b is exhausted, a is >= b (and strictly > if a has anything left):
	// j == len(b) means a's suffix wins. Otherwise compare the differing digits.
	return j == len(b) || (i < len(a) && a[i] > b[j])
}

func main() {
	fmt.Println("=== Approach 1: Merge of Best Sub-sequences (Optimal) ===")
	fmt.Println(createMaxNumber([]int{3, 4, 6, 5}, []int{9, 1, 2, 5, 8, 3}, 5)) // expected [9 8 6 5 3]
	fmt.Println(createMaxNumber([]int{6, 7}, []int{6, 0, 4}, 5))                // expected [6 7 6 0 4]
	fmt.Println(createMaxNumber([]int{3, 9}, []int{8, 9}, 3))                   // expected [9 8 9]
}
