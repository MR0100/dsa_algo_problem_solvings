package main

import "fmt"

// ListNode is a singly-linked list node.
type ListNode struct {
	Val  int
	Next *ListNode
}

// makeList builds a linked list from a slice (test helper).
func makeList(vals []int) *ListNode {
	dummy := &ListNode{}
	cur := dummy
	for _, v := range vals {
		cur.Next = &ListNode{Val: v}
		cur = cur.Next
	}
	return dummy.Next
}

// listToSlice flattens a linked list into a slice for printing (test helper).
func listToSlice(head *ListNode) []int {
	result := []int{}
	for head != nil {
		result = append(result, head.Val)
		head = head.Next
	}
	return result
}

// ── Approach 1: Iterative Without Dummy ──────────────────────────────────────
//
// iterativeWithoutDummy solves Remove Linked List Elements by treating the
// head as a special case, then splicing out matches in the body.
//
// Intuition:
//
//	Deleting an inner node is easy once you stand on the node BEFORE it
//	(prev.Next = prev.Next.Next). The only node with no "before" is the
//	head, so first discard leading matches until the head is safe, then
//	walk the rest with the standard splice.
//
// Algorithm:
//  1. While head != nil and head.Val == val: head = head.Next (drop leading
//     matches — possibly the entire list).
//  2. Walk cur from head; while cur != nil and cur.Next != nil:
//     if cur.Next.Val == val, bypass it (cur.Next = cur.Next.Next);
//     otherwise advance cur.
//  3. Return head.
//
// Time:  O(n) — every node is examined exactly once.
// Space: O(1) — pointer surgery in place.
func iterativeWithoutDummy(head *ListNode, val int) *ListNode {
	// Phase 1: the head has no predecessor — peel off matching heads first.
	for head != nil && head.Val == val {
		head = head.Next // old head becomes garbage; next node is new head
	}
	// Phase 2: cur always sits on a KEPT node, so cur.Next is deletable.
	cur := head
	for cur != nil && cur.Next != nil {
		if cur.Next.Val == val {
			cur.Next = cur.Next.Next // splice the matching node out
		} else {
			cur = cur.Next // next node survives; step onto it
		}
	}
	return head
}

// ── Approach 2: Recursion ────────────────────────────────────────────────────
//
// recursiveApproach solves Remove Linked List Elements by letting each call
// clean the suffix behind it, then deciding its own fate.
//
// Intuition:
//
//	"Remove val from a list" has a self-similar structure: a cleaned list is
//	(this node, if it survives) followed by the cleaned rest. Recurse to the
//	end; on the way back, each node links to the already-cleaned suffix and
//	either keeps itself or returns the suffix directly, deleting itself.
//
// Algorithm:
//  1. Base case: an empty list is already clean — return nil.
//  2. head.Next = recursiveApproach(head.Next, val) — clean the suffix.
//  3. If head.Val == val, return head.Next (skip self); else return head.
//
// Time:  O(n) — one call per node.
// Space: O(n) — recursion stack is as deep as the list (up to 10⁴ frames).
func recursiveApproach(head *ListNode, val int) *ListNode {
	// Base case: nothing to remove in an empty list.
	if head == nil {
		return nil
	}
	// Recursively clean everything after this node first.
	head.Next = recursiveApproach(head.Next, val)
	// Now decide this node's fate: matching nodes vanish by returning
	// the cleaned suffix instead of themselves.
	if head.Val == val {
		return head.Next
	}
	return head
}

// ── Approach 3: Dummy Node / Sentinel (Optimal) ──────────────────────────────
//
// dummyNode solves Remove Linked List Elements with a sentinel placed before
// the head, unifying head deletions with inner deletions.
//
// Intuition:
//
//	The head special-case exists only because the head lacks a predecessor.
//	Manufacture one: a dummy node whose Next is head. Now EVERY real node
//	has a node before it, one uniform loop handles all deletions (even
//	"delete the whole list"), and dummy.Next is always the true head.
//
// Algorithm:
//  1. dummy = &ListNode{Next: head}; prev = dummy.
//  2. While prev.Next != nil: if prev.Next.Val == val, splice it out
//     (prev.Next = prev.Next.Next); else advance prev.
//  3. Return dummy.Next.
//
// Time:  O(n) — single pass, O(1) work per node.
// Space: O(1) — one extra sentinel node regardless of list length.
func dummyNode(head *ListNode, val int) *ListNode {
	dummy := &ListNode{Next: head} // sentinel sits before the real head
	prev := dummy                  // prev is always the last KEPT node
	for prev.Next != nil {
		if prev.Next.Val == val {
			prev.Next = prev.Next.Next // unlink the matching node
		} else {
			prev = prev.Next // keep it; move the frontier forward
		}
	}
	return dummy.Next // real head, even if the original head was deleted
}

func main() {
	fmt.Println("=== Approach 1: Iterative Without Dummy ===")
	fmt.Printf("head=[1,2,6,3,4,5,6], val=6  got=%v  expected [1 2 3 4 5]\n", listToSlice(iterativeWithoutDummy(makeList([]int{1, 2, 6, 3, 4, 5, 6}), 6)))
	fmt.Printf("head=[], val=1               got=%v  expected []\n", listToSlice(iterativeWithoutDummy(makeList([]int{}), 1)))
	fmt.Printf("head=[7,7,7,7], val=7        got=%v  expected []\n", listToSlice(iterativeWithoutDummy(makeList([]int{7, 7, 7, 7}), 7)))

	fmt.Println("=== Approach 2: Recursion ===")
	fmt.Printf("head=[1,2,6,3,4,5,6], val=6  got=%v  expected [1 2 3 4 5]\n", listToSlice(recursiveApproach(makeList([]int{1, 2, 6, 3, 4, 5, 6}), 6)))
	fmt.Printf("head=[], val=1               got=%v  expected []\n", listToSlice(recursiveApproach(makeList([]int{}), 1)))
	fmt.Printf("head=[7,7,7,7], val=7        got=%v  expected []\n", listToSlice(recursiveApproach(makeList([]int{7, 7, 7, 7}), 7)))

	fmt.Println("=== Approach 3: Dummy Node / Sentinel (Optimal) ===")
	fmt.Printf("head=[1,2,6,3,4,5,6], val=6  got=%v  expected [1 2 3 4 5]\n", listToSlice(dummyNode(makeList([]int{1, 2, 6, 3, 4, 5, 6}), 6)))
	fmt.Printf("head=[], val=1               got=%v  expected []\n", listToSlice(dummyNode(makeList([]int{}), 1)))
	fmt.Printf("head=[7,7,7,7], val=7        got=%v  expected []\n", listToSlice(dummyNode(makeList([]int{7, 7, 7, 7}), 7)))
}
