package main

import "fmt"

// ListNode is a singly-linked list node.
type ListNode struct {
	Val  int
	Next *ListNode
}

// ── Approach 1: Brute Force (Repeated Tail Splice) ───────────────────────────
//
// bruteForce solves Reorder List by repeatedly walking to the current tail and
// splicing it right after the current front node.
//
// Intuition:
//
//	The target order L0, Ln, L1, Ln-1, ... interleaves "next from the front"
//	with "next from the back". A singly-linked list can't step backwards, but
//	we can always *find* the last node by walking to the end. So: for each
//	front node, walk to the tail, detach it, and insert it right after the
//	front; then the next front is the node after the freshly placed tail.
//
// Algorithm:
//  1. cur = head.
//  2. While cur has at least two nodes after it:
//     a. Walk (prev,last) to the tail so prev is the tail's predecessor.
//     b. If last is already cur.Next, the remainder is fully reordered → stop.
//     c. Detach last (prev.Next = nil), splice: last.Next = cur.Next,
//     cur.Next = last.
//     d. Advance cur past the spliced pair: cur = last.Next.
//
// Time:  O(n²) — each of the n/2 splices rescans up to n nodes to find the tail.
// Space: O(1) — pointer variables only.
func bruteForce(head *ListNode) {
	cur := head
	for cur != nil && cur.Next != nil {
		// Walk to the tail, tracking its predecessor for detachment.
		prev, last := cur, cur.Next
		for last.Next != nil {
			prev = last
			last = last.Next
		}
		if last == cur.Next {
			break // tail is directly after cur → nothing left to interleave
		}
		prev.Next = nil      // detach the tail from the list
		last.Next = cur.Next // tail now points at the old front-successor
		cur.Next = last      // front now points at the tail
		cur = last.Next      // skip over the pair (front, tail) just fixed
	}
}

// ── Approach 2: Array of Nodes (Two-Pointer Rebuild) ─────────────────────────
//
// arrayTwoPointers solves Reorder List by collecting node pointers into a
// slice, then rewiring Next pointers with two indices closing from both ends.
//
// Intuition:
//
//	The pain point is O(1)-time access to the k-th node from the back. A slice
//	of node pointers gives random access: nodes[i] from the front, nodes[j]
//	from the back. Alternate i and j while rewiring Next pointers in place.
//
// Algorithm:
//  1. One pass: append every node pointer to nodes[].
//  2. i = 0, j = len-1.
//  3. While i < j: nodes[i].Next = nodes[j]; i++; if i == j break;
//     nodes[j].Next = nodes[i]; j--.
//  4. Terminate the list: nodes[i].Next = nil (i == j is the new tail).
//
// Time:  O(n) — one collection pass + one rewiring pass.
// Space: O(n) — the slice of n node pointers.
func arrayTwoPointers(head *ListNode) {
	// Collect pointers for O(1) indexed access.
	var nodes []*ListNode
	for cur := head; cur != nil; cur = cur.Next {
		nodes = append(nodes, cur)
	}

	i, j := 0, len(nodes)-1
	for i < j {
		nodes[i].Next = nodes[j] // front node links to current back node
		i++                      // next front
		if i == j {
			break // pointers met → all nodes placed
		}
		nodes[j].Next = nodes[i] // back node links to next front node
		j--                      // next back
	}
	nodes[i].Next = nil // whoever the pointers met on becomes the new tail
}

// ── Approach 3: Middle + Reverse + Merge (Optimal) ───────────────────────────
//
// middleReverseMerge solves Reorder List in three O(n)/O(1) sub-steps:
// find the middle (slow/fast), reverse the second half, merge alternately.
//
// Intuition:
//
//	The result interleaves the first half (in order) with the second half
//	(in reverse): [L0..Lmid] and [Ln..Lmid+1]. Each sub-step is a classic
//	linked-list primitive we already know how to do in O(1) space:
//	#876 middle of list, #206 reverse list, then a zip-merge.
//
// Algorithm:
//  1. Fast/slow pointers: slow lands on the end of the first half
//     (for even n) or the exact middle (odd n).
//  2. Split: second = slow.Next; slow.Next = nil.
//  3. Reverse the second half iteratively (prev/cur/next rotation).
//  4. Zip: alternately link first-half node → second-half node → next
//     first-half node, until the (shorter) second half is exhausted.
//
// Time:  O(n) — three linear passes (find middle, reverse, merge).
// Space: O(1) — in-place pointer surgery only. Optimal.
func middleReverseMerge(head *ListNode) {
	if head == nil || head.Next == nil {
		return // 0 or 1 node → already reordered
	}

	// Step 1: find the end of the first half.
	// Using fast.Next/fast.Next.Next stops slow at ⌈n/2⌉-th node, so the
	// first half is never shorter than the second — required for the zip.
	slow, fast := head, head
	for fast.Next != nil && fast.Next.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}

	// Step 2: split the list after slow.
	second := slow.Next
	slow.Next = nil // terminate the first half

	// Step 3: reverse the second half in place.
	var prev *ListNode
	for cur := second; cur != nil; {
		next := cur.Next // save onward pointer
		cur.Next = prev  // flip the link backwards
		prev = cur       // advance prev
		cur = next       // advance cur
	}
	second = prev // prev is the new head of the reversed half

	// Step 4: zip-merge the two halves, alternating first → second.
	first := head
	for second != nil {
		n1, n2 := first.Next, second.Next // save both onward pointers
		first.Next = second               // front node → back node
		second.Next = n1                  // back node → next front node
		first, second = n1, n2            // advance both halves
	}
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
	var result []int
	for ; head != nil; head = head.Next {
		result = append(result, head.Val)
	}
	return result
}

func main() {
	// Official LeetCode examples: (input, expected reordering).
	examples := []struct {
		vals   []int
		expect []int
	}{
		{[]int{1, 2, 3, 4}, []int{1, 4, 2, 3}},       // Example 1
		{[]int{1, 2, 3, 4, 5}, []int{1, 5, 2, 4, 3}}, // Example 2
	}

	approaches := []struct {
		name string
		fn   func(*ListNode)
	}{
		{"Approach 1: Brute Force (Tail Splice)", bruteForce},
		{"Approach 2: Array + Two Pointers", arrayTwoPointers},
		{"Approach 3: Middle + Reverse + Merge (Optimal)", middleReverseMerge},
	}

	for _, ap := range approaches {
		fmt.Printf("=== %s ===\n", ap.name)
		for i, ex := range examples {
			head := makeList(ex.vals) // fresh list — every approach mutates
			ap.fn(head)
			fmt.Printf("Example %d: %v → %v (expected %v)\n",
				i+1, ex.vals, listToSlice(head), ex.expect) // expected: [1 4 2 3], [1 5 2 4 3]
		}
		fmt.Println()
	}
}
