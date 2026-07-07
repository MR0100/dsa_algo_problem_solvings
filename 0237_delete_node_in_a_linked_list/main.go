package main

import "fmt"

// ListNode is a singly-linked list node.
type ListNode struct {
	Val  int
	Next *ListNode
}

// ── Approach 1: Copy-Next-Value then Skip (Optimal) ──────────────────────────
//
// copyNextAndSkip solves Delete Node in a Linked List when you are given ONLY
// the node to delete (no head, guaranteed not the tail).
//
// Intuition:
//
//	We cannot reach the previous node to re-link around `node`, and we cannot
//	physically remove `node` from memory. But we CAN make `node` impersonate its
//	successor: copy the successor's value into `node`, then bypass the successor.
//	The list now reads as if `node` were deleted — the successor is the one
//	actually unlinked, but it carried the value we wanted gone.
//
// Algorithm:
//  1. Copy node.Next.Val into node.Val (node now looks like its successor).
//  2. Set node.Next = node.Next.Next (unlink the original successor).
//
// Time:  O(1) — two pointer/value assignments.
// Space: O(1) — no extra memory.
func copyNextAndSkip(node *ListNode) {
	node.Val = node.Next.Val   // steal the successor's value into this node
	node.Next = node.Next.Next // splice the (now duplicated) successor out
}

// ── Approach 2: Cascade-Shift Values (Illustrative Alternative) ───────────────
//
// cascadeShift achieves the same visible result by shifting every subsequent
// value one slot forward and dropping the final node, instead of the O(1)
// single-hop copy. It exists to contrast with the optimal trick.
//
// Intuition:
//
//	Deleting `node` is equivalent to shifting all following values left by one
//	position and then trimming the last node. This walks to the end, so it is
//	strictly worse than Approach 1 — but it shows the "delete = overwrite +
//	trim tail" mental model without relying on a single successor copy.
//
// Algorithm:
//  1. curr = node.
//  2. While curr.Next.Next != nil: curr.Val = curr.Next.Val; curr = curr.Next.
//  3. Now curr.Next is the last node: curr.Val = curr.Next.Val; curr.Next = nil.
//
// Time:  O(k) — k = number of nodes from `node` to the tail.
// Space: O(1).
func cascadeShift(node *ListNode) {
	curr := node // will slide down the tail copying values forward
	// Copy each successor's value back one slot until curr sits just before
	// the last node.
	for curr.Next.Next != nil {
		curr.Val = curr.Next.Val // overwrite with the next value
		curr = curr.Next         // advance
	}
	curr.Val = curr.Next.Val // absorb the last node's value
	curr.Next = nil          // trim the now-duplicate tail
}

// buildList constructs a linked list from a slice and returns head plus a
// value→node lookup (values are unique in the examples).
func buildList(vals []int) (*ListNode, map[int]*ListNode) {
	dummy := &ListNode{}
	curr := dummy
	lookup := map[int]*ListNode{}
	for _, v := range vals {
		curr.Next = &ListNode{Val: v}
		curr = curr.Next
		lookup[v] = curr
	}
	return dummy.Next, lookup
}

// toSlice serializes a list to a slice for easy comparison/printing.
func toSlice(head *ListNode) []int {
	out := []int{}
	for n := head; n != nil; n = n.Next {
		out = append(out, n.Val)
	}
	return out
}

func main() {
	fmt.Println("=== Approach 1: Copy-Next-Value then Skip (Optimal) ===")
	// Example 1: [4,5,1,9], delete node with value 5 → [4,1,9]
	head1, look1 := buildList([]int{4, 5, 1, 9})
	copyNextAndSkip(look1[5])
	fmt.Println(toSlice(head1)) // expected [4 1 9]

	// Example 2: [4,5,1,9], delete node with value 1 → [4,5,9]
	head2, look2 := buildList([]int{4, 5, 1, 9})
	copyNextAndSkip(look2[1])
	fmt.Println(toSlice(head2)) // expected [4 5 9]

	fmt.Println("=== Approach 2: Cascade-Shift Values (Illustrative Alternative) ===")
	head3, look3 := buildList([]int{4, 5, 1, 9})
	cascadeShift(look3[5])
	fmt.Println(toSlice(head3)) // expected [4 1 9]

	head4, look4 := buildList([]int{4, 5, 1, 9})
	cascadeShift(look4[1])
	fmt.Println(toSlice(head4)) // expected [4 5 9]
}
