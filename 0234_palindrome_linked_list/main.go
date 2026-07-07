package main

import "fmt"

// ListNode is the standard singly-linked list node used by LeetCode.
type ListNode struct {
	Val  int
	Next *ListNode
}

// build constructs a linked list from a slice and returns its head.
func build(vals []int) *ListNode {
	dummy := &ListNode{}
	cur := dummy
	for _, v := range vals { // append each value in order
		cur.Next = &ListNode{Val: v}
		cur = cur.Next
	}
	return dummy.Next
}

// ── Approach 1: Copy to Array + Two Pointers ─────────────────────────────────
//
// arrayTwoPointers solves Palindrome Linked List by copying values into a
// slice and checking the slice with two pointers.
//
// Intuition:
//
//	A singly-linked list can't be walked backwards, which is what a palindrome
//	check wants. The simplest fix is to materialise the values into an array
//	where random access is free, then compare the ends moving inward. Clear
//	and hard to get wrong, at the cost of O(n) extra memory.
//
// Algorithm:
//
//  1. Traverse the list once, appending each Val to a slice.
//  2. Set i = 0, j = len-1.
//  3. While i < j: if slice[i] != slice[j] return false; else i++, j--.
//  4. Return true.
//
// Time:  O(n) — one traversal plus one two-pointer pass.
// Space: O(n) — the values slice.
func arrayTwoPointers(head *ListNode) bool {
	vals := []int{}
	for node := head; node != nil; node = node.Next { // dump values into a slice
		vals = append(vals, node.Val)
	}
	i, j := 0, len(vals)-1
	for i < j { // compare symmetric ends moving toward the middle
		if vals[i] != vals[j] {
			return false // mismatch → not a palindrome
		}
		i++
		j--
	}
	return true
}

// ── Approach 2: Reverse Second Half In Place (Optimal) ───────────────────────
//
// reverseHalf solves Palindrome Linked List using O(1) extra space by finding
// the middle, reversing the second half, and comparing the two halves.
//
// Intuition:
//
//	We only need to compare the first half against the reversed second half.
//	A fast/slow pointer finds the midpoint in one pass; reversing the second
//	half in place then lets us walk both halves inward-out simultaneously with
//	no auxiliary array. This is the canonical O(1)-space solution.
//
// Algorithm:
//
//  1. Fast/slow walk: slow lands at the start of the second half.
//  2. Reverse the sublist starting at slow.
//  3. Walk `first` from head and `second` from the reversed head in lockstep,
//     comparing values until `second` runs out.
//  4. (Optional) restore the list; return whether all compared values matched.
//
// Time:  O(n) — midpoint scan + reversal + comparison, all linear.
// Space: O(1) — only a few pointers.
func reverseHalf(head *ListNode) bool {
	if head == nil || head.Next == nil {
		return true // empty or single node is trivially a palindrome
	}

	// Fast/slow: when fast reaches the end, slow is at the second-half start.
	slow, fast := head, head
	for fast != nil && fast.Next != nil {
		slow = slow.Next      // advances one node
		fast = fast.Next.Next // advances two nodes
	}

	// Reverse the second half beginning at slow.
	var prev *ListNode
	cur := slow
	for cur != nil {
		next := cur.Next // remember the rest of the list
		cur.Next = prev  // flip the pointer backward
		prev = cur       // advance prev
		cur = next       // advance cur
	}

	// Compare the first half (from head) with the reversed second half (prev).
	first, second := head, prev
	for second != nil { // second half is the shorter/equal one → drives the loop
		if first.Val != second.Val {
			return false // symmetric values differ
		}
		first = first.Next
		second = second.Next
	}
	return true
}

// ── Approach 3: Recursion (Front Pointer) ────────────────────────────────────
//
// recursive solves Palindrome Linked List by recursing to the tail and letting
// the call stack unwind back-to-front while a shared front pointer walks
// forward — comparing the ends inward.
//
// Intuition:
//
//	Recursion gives us a free reverse traversal: the deepest call sees the last
//	node, and as the stack unwinds we visit nodes from the back. Pairing each
//	unwinding node with a front pointer that advances from the head compares
//	position i against position n-1-i. Elegant, but O(n) stack space.
//
// Algorithm:
//
//  1. Keep a shared `front` pointer starting at head.
//  2. recurse(node): if node is nil return true. Recurse on node.Next; if that
//     sub-call already failed, propagate false. Compare node.Val with
//     front.Val, advance front, and return the equality.
//
// Time:  O(n) — each node visited once.
// Space: O(n) — recursion depth equals the list length.
func recursive(head *ListNode) bool {
	front := head // shared pointer sweeping from the front

	var recurse func(node *ListNode) bool
	recurse = func(node *ListNode) bool {
		if node == nil {
			return true // reached past the tail: base case
		}
		if !recurse(node.Next) { // dive to the end first
			return false // a deeper mismatch already failed
		}
		if node.Val != front.Val { // node unwinds from the back; front from the head
			return false
		}
		front = front.Next // advance the front pointer for the next comparison
		return true
	}
	return recurse(head)
}

func main() {
	fmt.Println("=== Approach 1: Copy to Array + Two Pointers ===")
	fmt.Println(arrayTwoPointers(build([]int{1, 2, 2, 1}))) // expected true
	fmt.Println(arrayTwoPointers(build([]int{1, 2})))       // expected false

	fmt.Println("=== Approach 2: Reverse Second Half In Place (Optimal) ===")
	fmt.Println(reverseHalf(build([]int{1, 2, 2, 1}))) // expected true
	fmt.Println(reverseHalf(build([]int{1, 2})))       // expected false

	fmt.Println("=== Approach 3: Recursion (Front Pointer) ===")
	fmt.Println(recursive(build([]int{1, 2, 2, 1}))) // expected true
	fmt.Println(recursive(build([]int{1, 2})))       // expected false
}
