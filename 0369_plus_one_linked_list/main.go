package main

import "fmt"

// ListNode is the standard LeetCode singly-linked-list node.
type ListNode struct {
	Val  int
	Next *ListNode
}

// ── Approach 1: Reverse, Add, Reverse Back (Brute Force) ─────────────────────
//
// reverseAddReverse solves Plus One Linked List by reversing the list so the
// least-significant digit comes first, adding one with carry, then reversing
// back.
//
// Intuition:
//
//	Grade-school addition starts at the least-significant digit, but the list
//	stores the most-significant digit first. Reverse the list to put the ones
//	digit at the head, propagate the +1 carry left-to-right, and reverse again
//	to restore order. If a carry survives past the front, prepend a new "1".
//
// Algorithm:
//  1. Reverse the list.
//  2. Walk it adding carry (initially 1); digit = (val+carry); store digit%10;
//     carry = digit/10.
//  3. If carry remains after the last node, append a node with that carry.
//  4. Reverse back and return the head.
//
// Time:  O(n) — three linear passes.
// Space: O(1) — in-place pointer surgery.
func reverseAddReverse(head *ListNode) *ListNode {
	reverse := func(node *ListNode) *ListNode {
		var prev *ListNode
		for node != nil {
			node.Next, prev, node = prev, node, node.Next // classic 3-way swap
		}
		return prev
	}

	head = reverse(head) // least-significant digit now at the head
	carry := 1           // the "+1" we must add
	cur := head
	var tail *ListNode // remember the last node to append an overflow digit
	for cur != nil {
		sum := cur.Val + carry
		cur.Val = sum % 10 // keep the low digit
		carry = sum / 10   // propagate the carry (0 or 1)
		tail = cur
		cur = cur.Next
	}
	if carry > 0 { // overflow past the most-significant digit (e.g. 99 → 100)
		tail.Next = &ListNode{Val: carry}
	}
	return reverse(head) // restore most-significant-first order
}

// ── Approach 2: Rightmost Non-Nine (Optimal, No Reverse) ──────────────────────
//
// rightmostNonNine solves Plus One Linked List by finding the last digit that
// is not 9; incrementing it and zeroing every 9 after it is exactly "add one".
//
// Intuition:
//
//	Adding 1 to a number only ripples a carry through a trailing run of 9s. The
//	last non-9 digit is where the carry stops: bump it by one, and turn every 9
//	to its right into 0. If EVERY digit is 9 (no non-nine exists), the whole
//	number rolls over — prepend a leading 1 and zero out the rest. A sentinel
//	(dummy) node in front lets the all-nines case reuse the same code: treat the
//	dummy's 0 as the last non-nine.
//
// Algorithm:
//  1. Create dummy → head. Track lastNotNine = dummy.
//  2. Walk the real nodes; whenever a node's value != 9, update lastNotNine.
//  3. Increment lastNotNine.Val by 1.
//  4. Set every node after lastNotNine to 0.
//  5. Return dummy.Next if dummy stayed 0, else dummy itself (a leading 1 was
//     created).
//
// Time:  O(n) — a single pass to locate, a partial pass to zero out.
// Space: O(1).
func rightmostNonNine(head *ListNode) *ListNode {
	dummy := &ListNode{Val: 0, Next: head} // sentinel absorbs an all-nines carry
	lastNotNine := dummy                   // rightmost node whose value != 9
	for node := head; node != nil; node = node.Next {
		if node.Val != 9 {
			lastNotNine = node
		}
	}
	lastNotNine.Val++ // the carry lands here and stops
	// Everything to the right of the increment point becomes 0.
	for node := lastNotNine.Next; node != nil; node = node.Next {
		node.Val = 0
	}
	if dummy.Val == 0 { // no rollover ⇒ dummy is unused
		return dummy.Next
	}
	return dummy // dummy became the new leading 1 (e.g. 999 → 1000)
}

// ── Approach 3: Recursion Carry-Back (Optimal) ───────────────────────────────
//
// recursionCarry solves Plus One Linked List by recursing to the tail, adding
// the carry on the way back up the call stack.
//
// Intuition:
//
//	Recursion naturally reaches the least-significant digit last, so we can add
//	the carry as the stack unwinds — exactly the right-to-left order arithmetic
//	needs, without any reversal. The base case (past the tail) returns carry 1;
//	each node adds it, keeps the low digit, and hands the new carry upward. A
//	surviving carry at the head means prepend a 1.
//
// Algorithm:
//  1. helper(node): if node == nil return 1 (the "+1").
//  2. carry = helper(node.Next); sum = node.Val + carry.
//  3. node.Val = sum % 10; return sum / 10.
//  4. In the caller: if helper(head) == 1, prepend a node with value 1.
//
// Time:  O(n) — one node per stack frame.
// Space: O(n) — recursion depth.
func recursionCarry(head *ListNode) *ListNode {
	var add func(node *ListNode) int
	add = func(node *ListNode) int {
		if node == nil {
			return 1 // reached the end; this is the +1 to inject
		}
		carry := add(node.Next) // finish the lower digits first
		sum := node.Val + carry
		node.Val = sum % 10 // keep low digit
		return sum / 10     // pass carry up
	}
	if add(head) == 1 { // carry escaped the most-significant digit
		return &ListNode{Val: 1, Next: head}
	}
	return head
}

// build makes a linked list from a slice of digits.
func build(digits []int) *ListNode {
	dummy := &ListNode{}
	cur := dummy
	for _, d := range digits {
		cur.Next = &ListNode{Val: d}
		cur = cur.Next
	}
	return dummy.Next
}

// toSlice flattens a linked list to a slice for printing/comparison.
func toSlice(head *ListNode) []int {
	var out []int
	for node := head; node != nil; node = node.Next {
		out = append(out, node.Val)
	}
	return out
}

func main() {
	// Example 1: head = [1,2,3] → [1,2,4]
	// Example 2: head = [0]     → [1]
	// Edge:      head = [9,9]   → [1,0,0]

	fmt.Println("=== Approach 1: Reverse, Add, Reverse Back (Brute Force) ===")
	fmt.Println(toSlice(reverseAddReverse(build([]int{1, 2, 3})))) // expected [1 2 4]
	fmt.Println(toSlice(reverseAddReverse(build([]int{0}))))       // expected [1]
	fmt.Println(toSlice(reverseAddReverse(build([]int{9, 9}))))    // expected [1 0 0]

	fmt.Println("=== Approach 2: Rightmost Non-Nine (Optimal, No Reverse) ===")
	fmt.Println(toSlice(rightmostNonNine(build([]int{1, 2, 3})))) // expected [1 2 4]
	fmt.Println(toSlice(rightmostNonNine(build([]int{0}))))       // expected [1]
	fmt.Println(toSlice(rightmostNonNine(build([]int{9, 9}))))    // expected [1 0 0]

	fmt.Println("=== Approach 3: Recursion Carry-Back (Optimal) ===")
	fmt.Println(toSlice(recursionCarry(build([]int{1, 2, 3})))) // expected [1 2 4]
	fmt.Println(toSlice(recursionCarry(build([]int{0}))))       // expected [1]
	fmt.Println(toSlice(recursionCarry(build([]int{9, 9}))))    // expected [1 0 0]
}
