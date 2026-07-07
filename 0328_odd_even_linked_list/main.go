package main

import (
	"fmt"
	"strings"
)

// ListNode is a singly-linked list node.
type ListNode struct {
	Val  int
	Next *ListNode
}

// build constructs a linked list from a slice of values and returns its head.
// An empty slice yields a nil head (the empty list).
func build(vals []int) *ListNode {
	dummy := &ListNode{} // sentinel so we never special-case the first node
	tail := dummy        // tail always points at the last appended node
	for _, v := range vals {
		tail.Next = &ListNode{Val: v} // link a fresh node after the current tail
		tail = tail.Next              // advance the tail onto it
	}
	return dummy.Next // node after the sentinel is the real head (nil if empty)
}

// toSlice walks a list and returns its values as a []int, for easy printing and
// output verification. A nil head yields an empty slice.
func toSlice(head *ListNode) []int {
	out := []int{}
	for cur := head; cur != nil; cur = cur.Next {
		out = append(out, cur.Val) // record each value in list order
	}
	return out
}

// toString renders a list as "1 -> 2 -> 3" (or "<empty>"), handy for dry-run
// style debugging. Not used by the graders below but kept as a helper.
func toString(head *ListNode) string {
	vals := toSlice(head)
	if len(vals) == 0 {
		return "<empty>"
	}
	parts := make([]string, len(vals))
	for i, v := range vals {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, " -> ")
}

// ── Approach 1: Two Sublists via Builder Tails ───────────────────────────────
//
// twoListsExtraNodes solves Odd Even Linked List by walking the original list
// once and threading each node onto one of two growing sublists — an "odd"
// sublist and an "even" sublist — by its 1-based position, then joining the
// odd tail to the even head.
//
// Despite the name, NO new nodes are allocated: we rewire the *existing* nodes
// by moving them under two builder tails, so this remains O(1) extra space
// (only a handful of pointers). It is the clarity-first cousin of the classic
// weave — the two sublists are built explicitly rather than interleaved in
// place, which makes the "odd positions first, then even positions" grouping
// obvious.
//
// Intuition:
//
//	The reordering keeps relative order inside each group, so a stable split
//	is enough: sweep the list, and for each node decide "is my position odd
//	or even?" and append it to the matching sublist. Finally the answer is
//	the odd sublist followed by the even sublist. Two dummy heads plus two
//	tail pointers let us append in O(1) without touching values.
//
// Algorithm:
//  1. If head is nil or head.Next is nil, return head unchanged (0 or 1 node —
//     nothing to reorder).
//  2. Create two sentinel dummies oddDummy and evenDummy with tails oddTail,
//     evenTail pointing at them.
//  3. Walk the original list with a running 1-based position counter:
//     - odd position  → append current node to oddTail, advance oddTail.
//     - even position → append current node to evenTail, advance evenTail.
//  4. Terminate the even sublist: evenTail.Next = nil (the last appended node
//     may still point into the original list).
//  5. Join: oddTail.Next = evenDummy.Next (odd group, then even group).
//  6. Return oddDummy.Next — the head of the odd group.
//
// Time:  O(n) — a single pass, O(1) work per node.
// Space: O(1) — four pointers (two dummies, two tails); no node allocation.
func twoListsExtraNodes(head *ListNode) *ListNode {
	// Edge cases: an empty list or a single node is already grouped correctly.
	if head == nil || head.Next == nil {
		return head
	}

	oddDummy := &ListNode{}  // sentinel before the odd-position sublist
	evenDummy := &ListNode{} // sentinel before the even-position sublist
	oddTail := oddDummy      // last node currently in the odd sublist
	evenTail := evenDummy    // last node currently in the even sublist

	pos := 1    // 1-based position of `cur` in the original list
	cur := head // node we are currently classifying
	for cur != nil {
		next := cur.Next // remember the successor before we rewire cur.Next
		if pos%2 == 1 {
			oddTail.Next = cur // odd position → extend the odd sublist
			oddTail = cur      // advance the odd tail onto this node
		} else {
			evenTail.Next = cur // even position → extend the even sublist
			evenTail = cur      // advance the even tail onto this node
		}
		cur = next // move to the successor we saved
		pos++      // its position is one greater
	}

	evenTail.Next = nil           // cap the even sublist so it can't loop back
	oddTail.Next = evenDummy.Next // stitch odd group in front of even group
	return oddDummy.Next          // head of the odd group is the new head
}

// ── Approach 2: In-Place Two-Pointer Weave (Optimal) ─────────────────────────
//
// inPlacePointers solves Odd Even Linked List using the canonical two-pointer
// weave: an `odd` pointer chasing odd-position nodes and an `even` pointer
// chasing even-position nodes, plus `evenHead` remembering where the even group
// starts so we can attach it after the odd group at the end.
//
// Intuition:
//
//	Odd-position and even-position nodes already alternate in the original
//	list (odd, even, odd, even, ...). So from any odd node, the *next* odd
//	node is exactly two hops away — i.e. odd.Next.Next — and likewise for
//	evens. We repeatedly splice each pointer past its neighbour, unzipping
//	the single list into two interleaved chains without moving values. When
//	the evens run out, we reconnect the tail of the odd chain to the head of
//	the even chain.
//
// Algorithm:
//  1. If head is nil or head.Next is nil, return head unchanged (0 or 1 node).
//  2. odd = head, even = head.Next, evenHead = even (start of even group).
//  3. While even != nil and even.Next != nil:
//     - odd.Next  = even.Next; odd  = odd.Next   (jump odd to next odd node)
//     - even.Next = odd.Next;  even = even.Next  (jump even to next even node)
//  4. odd.Next = evenHead (append the whole even group after the odd group).
//  5. Return head — still the first odd node, now head of the reordered list.
//
// Time:  O(n) — each node is visited a constant number of times in one sweep.
// Space: O(1) — three pointers (odd, even, evenHead); pure pointer surgery.
func inPlacePointers(head *ListNode) *ListNode {
	// Edge cases: 0 or 1 node — already grouped, nothing to weave.
	if head == nil || head.Next == nil {
		return head
	}

	odd := head       // walks the odd-position (1,3,5,...) chain
	even := head.Next // walks the even-position (2,4,6,...) chain
	evenHead := even  // remember the even group's head for the final join

	// Continue while there is a further even node to relink. Testing both
	// `even` and `even.Next` guards the two length parities safely.
	for even != nil && even.Next != nil {
		odd.Next = even.Next // odd skips the even node to reach the next odd
		odd = odd.Next       // advance odd onto that next-odd node
		even.Next = odd.Next // even skips the (new) odd node to the next even
		even = even.Next     // advance even onto that next-even node
	}

	odd.Next = evenHead // splice the even group onto the end of the odd group
	return head         // the original first node is still the overall head
}

func main() {
	// Example 1: [1,2,3,4,5] → [1,3,5,2,4]
	fmt.Println("=== Approach 1: Two Sublists via Builder Tails ===")
	fmt.Println(toSlice(twoListsExtraNodes(build([]int{1, 2, 3, 4, 5}))))       // expected [1 3 5 2 4]
	fmt.Println(toSlice(twoListsExtraNodes(build([]int{2, 1, 3, 5, 6, 4, 7})))) // expected [2 3 6 7 1 5 4]
	fmt.Println(toSlice(twoListsExtraNodes(build([]int{}))))                    // expected []
	fmt.Println(toSlice(twoListsExtraNodes(build([]int{1}))))                   // expected [1]

	fmt.Println("=== Approach 2: In-Place Two-Pointer Weave (Optimal) ===")
	fmt.Println(toSlice(inPlacePointers(build([]int{1, 2, 3, 4, 5}))))       // expected [1 3 5 2 4]
	fmt.Println(toSlice(inPlacePointers(build([]int{2, 1, 3, 5, 6, 4, 7})))) // expected [2 3 6 7 1 5 4]
	fmt.Println(toSlice(inPlacePointers(build([]int{}))))                    // expected []
	fmt.Println(toSlice(inPlacePointers(build([]int{1}))))                   // expected [1]
}
