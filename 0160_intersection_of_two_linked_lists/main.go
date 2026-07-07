package main

import "fmt"

// ListNode is a singly linked list node.
type ListNode struct {
	Val  int
	Next *ListNode
}

// ── Approach 1: Brute Force (Nested Scan) ────────────────────────────────────
//
// bruteForce solves Intersection of Two Linked Lists by comparing every pair
// of nodes.
//
// Intuition:
//
//	Intersection is defined by POINTER identity, not value equality. So for
//	each node a in list A, walk all of list B and check whether some node b
//	is the very same node (a == b). The first match, scanning A front to
//	back, is the intersection start.
//
// Algorithm:
//  1. For each node a in A (in order):
//     a. For each node b in B: if a == b, return a.
//  2. No pair matched → return nil.
//
// Time:  O(m·n) — every A node may scan all of B.
// Space: O(1) — only two cursors.
func bruteForce(headA, headB *ListNode) *ListNode {
	for a := headA; a != nil; a = a.Next {
		for b := headB; b != nil; b = b.Next {
			if a == b { // same node in memory → lists merge here
				return a
			}
		}
	}
	return nil // exhausted all pairs: the lists never meet
}

// ── Approach 2: Hash Set ─────────────────────────────────────────────────────
//
// hashSet solves Intersection of Two Linked Lists by remembering A's nodes.
//
// Intuition:
//
//	Store every node pointer of list A in a set. Then walk list B: the first
//	node already present in the set is where B merges into A — after the
//	merge point both lists share ALL subsequent nodes, so the first hit is
//	the intersection start.
//
// Algorithm:
//  1. Insert each node of A (the pointer, not the value) into a hash set.
//  2. Walk B; return the first node found in the set.
//  3. If B ends with no hit, return nil.
//
// Time:  O(m+n) — one pass over each list.
// Space: O(m) — the set holds all nodes of A.
func hashSet(headA, headB *ListNode) *ListNode {
	seen := map[*ListNode]bool{} // identity set of A's nodes
	for a := headA; a != nil; a = a.Next {
		seen[a] = true // record the pointer itself
	}
	for b := headB; b != nil; b = b.Next {
		if seen[b] { // first shared node = intersection start
			return b
		}
	}
	return nil // B shares no node with A
}

// ── Approach 3: Length Difference ────────────────────────────────────────────
//
// lengthDifference solves Intersection of Two Linked Lists by aligning tails.
//
// Intuition:
//
//	If the lists intersect, they share a common TAIL. Measured from the end,
//	the intersection is at the same distance in both lists — the only
//	misalignment is the length difference of the prefixes. Advance the longer
//	list's cursor by |m-n| first; then both cursors are equally far from the
//	end and can walk in lockstep until they collide (or hit nil together).
//
// Algorithm:
//  1. Compute lenA and lenB with one pass each.
//  2. Advance the head of the longer list by |lenA-lenB| nodes.
//  3. Walk both cursors together; the first a == b is the answer
//     (nil == nil covers the disjoint case).
//
// Time:  O(m+n) — two measuring passes plus one aligned pass.
// Space: O(1) — a few cursors and two ints.
func lengthDifference(headA, headB *ListNode) *ListNode {
	// measure both lists
	lenA, lenB := 0, 0
	for n := headA; n != nil; n = n.Next {
		lenA++
	}
	for n := headB; n != nil; n = n.Next {
		lenB++
	}
	a, b := headA, headB
	// skip the surplus prefix of the longer list so both are aligned
	for ; lenA > lenB; lenA-- {
		a = a.Next
	}
	for ; lenB > lenA; lenB-- {
		b = b.Next
	}
	// lockstep walk: equal distance from the end guarantees they meet at the
	// intersection node, or reach nil simultaneously if disjoint
	for a != b {
		a = a.Next
		b = b.Next
	}
	return a
}

// ── Approach 4: Two Pointers Switching Heads (Optimal) ───────────────────────
//
// twoPointers solves Intersection of Two Linked Lists with the elegant
// path-swap trick.
//
// Intuition:
//
//	Let A = a + c and B = b + c, where a/b are the exclusive prefixes and c
//	the shared tail (possibly empty). Walk pointer pa through A then restart
//	at B's head; walk pb through B then restart at A's head. Both traverse
//	exactly a + c + b steps, so they arrive TOGETHER at the intersection
//	(after a + b steps they are aligned at c's start) — or together at nil
//	when c is empty. Two passes, zero extra memory, no length arithmetic.
//
// Algorithm:
//  1. pa = headA, pb = headB.
//  2. While pa != pb:
//     a. pa = (pa == nil) ? headB : pa.Next
//     b. pb = (pb == nil) ? headA : pb.Next
//  3. Return pa (intersection node, or nil if none).
//
// Time:  O(m+n) — each pointer walks each list at most once.
// Space: O(1) — two pointers.
func twoPointers(headA, headB *ListNode) *ListNode {
	if headA == nil || headB == nil {
		return nil // an empty list cannot intersect anything
	}
	pa, pb := headA, headB
	for pa != pb {
		if pa == nil {
			pa = headB // finished path A → continue on path B
		} else {
			pa = pa.Next
		}
		if pb == nil {
			pb = headA // finished path B → continue on path A
		} else {
			pb = pb.Next
		}
	}
	return pa // meeting point: intersection node, or nil (both exhausted)
}

// buildList turns a value slice into a list, returning head and tail.
func buildList(vals []int) (head, tail *ListNode) {
	for _, v := range vals {
		node := &ListNode{Val: v}
		if head == nil {
			head = node // first node becomes the head
		} else {
			tail.Next = node // append at the end
		}
		tail = node
	}
	return head, tail
}

// buildCase constructs two lists that share the common suffix (by pointer).
func buildCase(prefixA, prefixB, common []int) (headA, headB *ListNode) {
	commonHead, _ := buildList(common)
	headA, tailA := buildList(prefixA)
	headB, tailB := buildList(prefixB)
	if tailA != nil {
		tailA.Next = commonHead // graft the SAME nodes onto A
	} else {
		headA = commonHead // A is entirely the common part
	}
	if tailB != nil {
		tailB.Next = commonHead // graft the SAME nodes onto B
	} else {
		headB = commonHead
	}
	return headA, headB
}

// describe formats a result the way LeetCode does.
func describe(n *ListNode) string {
	if n == nil {
		return "No intersection"
	}
	return fmt.Sprintf("Intersected at '%d'", n.Val)
}

func main() {
	solvers := []struct {
		name string
		fn   func(a, b *ListNode) *ListNode
	}{
		{"Approach 1: Brute Force (Nested Scan)", bruteForce},
		{"Approach 2: Hash Set", hashSet},
		{"Approach 3: Length Difference", lengthDifference},
		{"Approach 4: Two Pointers Switching Heads (Optimal)", twoPointers},
	}
	for _, s := range solvers {
		fmt.Printf("=== %s ===\n", s.name)
		// Example 1: A = [4,1] + [8,4,5], B = [5,6,1] + [8,4,5] → '8'
		a1, b1 := buildCase([]int{4, 1}, []int{5, 6, 1}, []int{8, 4, 5})
		fmt.Println(describe(s.fn(a1, b1))) // expected Intersected at '8'

		// Example 2: A = [1,9,1] + [2,4], B = [3] + [2,4] → '2'
		a2, b2 := buildCase([]int{1, 9, 1}, []int{3}, []int{2, 4})
		fmt.Println(describe(s.fn(a2, b2))) // expected Intersected at '2'

		// Example 3: A = [2,6,4], B = [1,5], no shared part → nil
		a3, _ := buildList([]int{2, 6, 4})
		b3, _ := buildList([]int{1, 5})
		fmt.Println(describe(s.fn(a3, b3))) // expected No intersection
	}
}
