package main

import "fmt"

// ── Approach 1: Recursive Range Partition (Divide & Conquer) ─────────────────
//
// recursivePartition verifies a preorder sequence of a BST by splitting each
// subarray into the left and right subtree ranges around the root.
//
// Intuition:
//
//	In preorder, the first element of any (sub)sequence is the subtree's root.
//	Everything after it that is smaller belongs to the left subtree, and it must
//	form a contiguous prefix; everything that follows must all be greater than
//	the root (the right subtree). If any element after the left-run is not > root,
//	the ordering is impossible. Recurse into both parts with the tightened bounds.
//
// Algorithm:
//  1. verify(lo, hi, min, max): the first element preorder[lo] is the root; it
//     must lie strictly within (min, max).
//  2. Scan forward from lo+1 while values < root — that's the left subtree; the
//     rest is the right subtree, and every right-subtree value must be > root.
//  3. Recurse on left with bounds (min, root) and right with (root, max).
//
// Time:  O(n²) worst case (degenerate/skewed splits rescan).
// Space: O(n) recursion depth.
func recursivePartition(preorder []int) bool {
	var verify func(lo, hi, min, max int) bool
	verify = func(lo, hi, min, max int) bool {
		if lo > hi {
			return true // empty range is a valid (empty) BST
		}
		root := preorder[lo]
		// Root must respect the inherited bounds from ancestors.
		if root <= min || root >= max {
			return false
		}
		// Find where the left subtree (values < root) ends.
		i := lo + 1
		for i <= hi && preorder[i] < root {
			i++
		}
		// Everything from i..hi is the right subtree; all must exceed root.
		for j := i; j <= hi; j++ {
			if preorder[j] <= root {
				return false // a "right" value not greater than root ⇒ invalid
			}
		}
		// Recurse into left (bounded above by root) and right (bounded below).
		return verify(lo+1, i-1, min, root) && verify(i, hi, root, max)
	}
	return verify(0, len(preorder)-1, -1<<63, 1<<63-1)
}

// ── Approach 2: Monotonic Stack with Lower Bound (Optimal) ────────────────────
//
// monotonicStack verifies the preorder sequence in one linear pass using a
// decreasing stack and a running lower bound.
//
// Intuition:
//
//	Walk the preorder left to right. While we keep descending left children the
//	values decrease, so we push them onto a stack (kept decreasing). When we hit
//	a value bigger than the stack top, we've turned into a right subtree: we pop
//	all smaller values, and the last popped value becomes a new lower bound —
//	nothing that follows may be ≤ it (it lives in that popped node's right
//	subtree). If any later value violates the lower bound, the sequence is invalid.
//
// Algorithm:
//  1. lowerBound = -inf; empty stack.
//  2. For each value v: if v <= lowerBound, return false.
//  3. While stack non-empty and v > stack.top(): lowerBound = pop().
//  4. Push v.
//  5. Return true.
//
// Time:  O(n) — each element pushed and popped at most once.
// Space: O(n) — stack, O(1) if you overwrite the input in place.
func monotonicStack(preorder []int) bool {
	lowerBound := -1 << 63 // nothing may be ≤ this yet (negative infinity)
	stack := []int{}
	for _, v := range preorder {
		// Once we've moved into a right subtree, everything must exceed the bound.
		if v <= lowerBound {
			return false
		}
		// Turning right: pop all ancestors we are now to the right of; the last
		// popped is the closest ancestor whose right subtree we entered.
		for len(stack) > 0 && v > stack[len(stack)-1] {
			lowerBound = stack[len(stack)-1] // raise the floor
			stack = stack[:len(stack)-1]     // pop
		}
		stack = append(stack, v) // push current node (descending left chain)
	}
	return true
}

func main() {
	fmt.Println("=== Approach 1: Recursive Range Partition ===")
	fmt.Println(recursivePartition([]int{5, 2, 1, 3, 6})) // expected true
	fmt.Println(recursivePartition([]int{5, 2, 6, 1, 3})) // expected false
	fmt.Println(recursivePartition([]int{1, 3, 2}))       // expected true
	fmt.Println(recursivePartition([]int{2, 1, 3}))       // expected true

	fmt.Println("=== Approach 2: Monotonic Stack (Optimal) ===")
	fmt.Println(monotonicStack([]int{5, 2, 1, 3, 6})) // expected true
	fmt.Println(monotonicStack([]int{5, 2, 6, 1, 3})) // expected false
	fmt.Println(monotonicStack([]int{1, 3, 2}))       // expected true
	fmt.Println(monotonicStack([]int{2, 1, 3}))       // expected true
}
