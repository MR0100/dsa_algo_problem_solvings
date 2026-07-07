package main

import "fmt"

// ── Approach 1: Two Pointers (Optimal) ───────────────────────────────────────
//
// twoPointers reverses the input character slice in place by swapping the ends
// inward.
//
// Intuition:
//
//	Reversing means the first char swaps with the last, the second with the
//	second-last, and so on. Walk one pointer from the left and one from the
//	right, swapping and stepping toward the middle. Stop when they cross.
//
// Algorithm:
//  1. left = 0, right = len-1.
//  2. While left < right: swap s[left], s[right]; left++; right--.
//
// Time:  O(n) — each element visited once. Space: O(1) — in place.
func twoPointers(s []byte) {
	left, right := 0, len(s)-1
	for left < right { // stop when the pointers meet/cross
		s[left], s[right] = s[right], s[left] // swap the mirror pair
		left++                                // move inward from the left
		right--                               // move inward from the right
	}
}

// ── Approach 2: Recursion (Swap Ends, Recurse Inward) ────────────────────────
//
// recursion reverses in place by swapping the outermost pair, then recursing on
// the inner sub-slice.
//
// Intuition:
//
//	Reverse(s[left..right]) = swap the ends, then Reverse(s[left+1..right-1]).
//	The base case is when the window has 0 or 1 element left.
//
// Algorithm:
//  1. helper(left, right): if left >= right return.
//  2. swap s[left], s[right]; recurse helper(left+1, right-1).
//
// Time:  O(n). Space: O(n) — recursion stack depth n/2.
func recursion(s []byte) {
	var helper func(left, right int)
	helper = func(left, right int) {
		if left >= right { // window shrunk to empty/single element
			return
		}
		s[left], s[right] = s[right], s[left] // swap the current ends
		helper(left+1, right-1)               // recurse on the inner window
	}
	helper(0, len(s)-1)
}

// reverseCopy is a small helper so main() can print results without mutating a
// shared input across the two approaches.
func reverseCopy(orig []byte, fn func([]byte)) string {
	cp := make([]byte, len(orig))
	copy(cp, orig)
	fn(cp)
	return string(cp)
}

func main() {
	ex1 := []byte("hello")
	ex2 := []byte("Hannah")

	fmt.Println("=== Approach 1: Two Pointers (Optimal) ===")
	fmt.Println(reverseCopy(ex1, twoPointers)) // expected olleh
	fmt.Println(reverseCopy(ex2, twoPointers)) // expected hannaH

	fmt.Println("=== Approach 2: Recursion ===")
	fmt.Println(reverseCopy(ex1, recursion)) // expected olleh
	fmt.Println(reverseCopy(ex2, recursion)) // expected hannaH
}
