package main

import "fmt"

type ListNode struct {
	Val  int
	Next *ListNode
}

// ── Approach 1: Iterative ─────────────────────────────────────────────────────
//
// iterative swaps consecutive pairs using a dummy head and a prev pointer.
//
// Intuition:
//   Keep a pointer `prev` that sits before the current pair.
//   For each pair (first, second):
//     1. prev.Next = second        (link prev to second)
//     2. first.Next = second.Next  (first skips to node after second)
//     3. second.Next = first       (second points back to first)
//     4. prev = first              (first is now the "tail" of the swapped pair)
//
// Time:  O(n) — one pass.
// Space: O(1) — in-place relinking.
func iterative(head *ListNode) *ListNode {
	dummy := &ListNode{Next: head}
	prev := dummy

	for prev.Next != nil && prev.Next.Next != nil {
		first := prev.Next
		second := prev.Next.Next

		// Swap: prev → second → first → (rest)
		prev.Next = second
		first.Next = second.Next
		second.Next = first

		// Advance prev to the tail of the swapped pair.
		prev = first
	}
	return dummy.Next
}

// ── Approach 2: Recursive ─────────────────────────────────────────────────────
//
// recursive swaps the first pair, then recurses for the tail.
//
// Intuition:
//   Base case: fewer than two nodes → return head unchanged.
//   Recursive case:
//     second = head.Next
//     head.Next = recursive(second.Next)  // head now points to swapped tail
//     second.Next = head                   // second now leads the pair
//     return second
//
// Time:  O(n) — one call per pair.
// Space: O(n) — recursion stack depth (n/2 calls).
func recursive(head *ListNode) *ListNode {
	// Base case: 0 or 1 node — nothing to swap.
	if head == nil || head.Next == nil {
		return head
	}

	second := head.Next
	head.Next = recursive(second.Next) // recurse on the rest
	second.Next = head                 // second now precedes head
	return second
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func makeList(vals []int) *ListNode {
	dummy := &ListNode{}
	cur := dummy
	for _, v := range vals {
		cur.Next = &ListNode{Val: v}
		cur = cur.Next
	}
	return dummy.Next
}

func listToSlice(head *ListNode) []int {
	var out []int
	for head != nil {
		out = append(out, head.Val)
		head = head.Next
	}
	return out
}

func main() {
	examples := []struct {
		vals   []int
		expect []int
	}{
		{[]int{1, 2, 3, 4}, []int{2, 1, 4, 3}},
		{[]int{}, nil},
		{[]int{1}, []int{1}},
		{[]int{1, 2, 3}, []int{2, 1, 3}},
	}

	approaches := []struct {
		name string
		fn   func(*ListNode) *ListNode
	}{
		{"Approach 1: Iterative  ✅ O(n) T | O(1) S", iterative},
		{"Approach 2: Recursive     O(n) T | O(n) S", recursive},
	}

	for _, ex := range examples {
		fmt.Printf("vals=%v  expect=%v\n", ex.vals, ex.expect)
		for _, ap := range approaches {
			fmt.Printf("  %-45s → %v\n", ap.name, listToSlice(ap.fn(makeList(ex.vals))))
		}
		fmt.Println()
	}
}
