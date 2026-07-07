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

// ── Approach 1: Brute Force (Copy Values & Rebuild) ──────────────────────────
//
// bruteForce solves Reverse Linked List by copying every value into a slice
// and building a brand-new list from that slice back to front.
//
// Intuition:
//
//	If pointer surgery feels risky, sidestep it entirely: a linked list is
//	just a sequence of values, so record the sequence, then manufacture a
//	fresh list that emits those values in the opposite order. No existing
//	pointer is ever touched — correctness is trivial to see.
//
// Algorithm:
//  1. Walk the list once, appending each node's value to a slice vals.
//  2. Iterate vals from the last index down to 0, appending a new node per
//     value onto a growing result list (dummy-tail construction).
//  3. Return the node after the dummy — the head of the reversed copy.
//
// Time:  O(n) — one pass to collect values, one pass to rebuild.
// Space: O(n) — the slice of n values plus n brand-new nodes.
func bruteForce(head *ListNode) *ListNode {
	// Collect the values in original order.
	vals := []int{}
	for cur := head; cur != nil; cur = cur.Next {
		vals = append(vals, cur.Val) // remember each value as we pass it
	}
	// Rebuild a completely new list, reading the slice backwards.
	dummy := &ListNode{} // sentinel so the first append needs no special case
	tail := dummy        // tail always points at the last node built so far
	for i := len(vals) - 1; i >= 0; i-- {
		tail.Next = &ListNode{Val: vals[i]} // append a fresh node holding vals[i]
		tail = tail.Next                    // advance tail onto the new node
	}
	return dummy.Next // real head sits right after the sentinel
}

// ── Approach 2: Stack of Nodes ───────────────────────────────────────────────
//
// stackBased solves Reverse Linked List by pushing every node pointer onto a
// stack and re-linking them in pop (LIFO) order.
//
// Intuition:
//
//	A stack reverses order by nature: the last node pushed is the first node
//	popped. Push the whole list, then pop nodes one by one and chain each to
//	the next pop — the chain comes out reversed, reusing the original nodes.
//
// Algorithm:
//  1. Traverse the list, pushing each *ListNode onto a slice-backed stack.
//  2. If the stack is empty, return nil (empty input).
//  3. Pop the top node — it becomes the new head.
//  4. Keep popping, wiring the previous pop's Next to the current pop.
//  5. Set the final node's Next to nil to terminate the list.
//
// Time:  O(n) — each node is pushed once and popped once.
// Space: O(n) — the stack holds all n node pointers.
func stackBased(head *ListNode) *ListNode {
	// Push every node onto the stack in traversal order.
	stack := []*ListNode{}
	for cur := head; cur != nil; cur = cur.Next {
		stack = append(stack, cur) // node pointers, not values — we reuse nodes
	}
	if len(stack) == 0 {
		return nil // empty list reverses to the empty list
	}
	// The last node pushed (original tail) is the new head.
	newHead := stack[len(stack)-1]
	cur := newHead
	// Pop the remaining nodes, chaining each onto the reversed list.
	for i := len(stack) - 2; i >= 0; i-- {
		cur.Next = stack[i] // wire current node to the next-popped node
		cur = cur.Next      // advance along the growing reversed list
	}
	cur.Next = nil // the original head is now the tail — terminate it
	return newHead
}

// ── Approach 3: Recursion ────────────────────────────────────────────────────
//
// recursive solves Reverse Linked List by reversing the sublist after head
// first, then hooking head onto the end of that reversed sublist.
//
// Intuition:
//
//	Assume recursion magically reverses the rest of the list: head → (rest
//	reversed). Then head.Next is the *tail* of the reversed rest, so pointing
//	head.Next.Next back at head appends head to the end. The deepest call
//	returns the original tail, which stays the new head all the way up.
//
// Algorithm:
//  1. Base case: an empty list or single node is its own reversal — return head.
//  2. Recurse on head.Next; the returned newHead is the original list's tail.
//  3. head.Next currently points at the last node of the reversed sublist,
//     so set head.Next.Next = head to attach head after it.
//  4. Set head.Next = nil so head terminates the list (it may be the new tail).
//  5. Return newHead unchanged.
//
// Time:  O(n) — one recursive call per node, constant work in each.
// Space: O(n) — the recursion stack is n frames deep.
func recursive(head *ListNode) *ListNode {
	// Base case: nothing to reverse for an empty or single-node list.
	if head == nil || head.Next == nil {
		return head
	}
	// Reverse everything after head; newHead is the original tail.
	newHead := recursive(head.Next)
	// head.Next is the last node of the reversed sublist — hook head behind it.
	head.Next.Next = head
	// Break head's old forward pointer; it now ends the list (or gets
	// re-pointed by the caller one level up).
	head.Next = nil
	return newHead
}

// ── Approach 4: Iterative Pointer Reversal (Optimal) ─────────────────────────
//
// iterative solves Reverse Linked List by walking the list once and flipping
// each node's Next pointer to face the previous node.
//
// Intuition:
//
//	Reversing a list is exactly "make every arrow point the other way". Walk
//	with two pointers — prev (already-reversed portion) and cur (first
//	unprocessed node) — and at each step redirect cur.Next from its successor
//	to prev. One temporary pointer saves the successor before the flip so the
//	walk can continue. When cur falls off the end, prev is the new head.
//
// Algorithm:
//  1. Initialise prev = nil (reversed part is empty) and cur = head.
//  2. While cur != nil:
//     a. next := cur.Next   — save the rest of the list before overwriting.
//     b. cur.Next = prev    — flip the arrow to point backwards.
//     c. prev = cur         — the reversed part grows by one node.
//     d. cur = next         — step forward into the unprocessed part.
//  3. Return prev — the last real node processed, now the head.
//
// Time:  O(n) — exactly one visit per node with O(1) work.
// Space: O(1) — three pointers regardless of list length.
func iterative(head *ListNode) *ListNode {
	var prev *ListNode // head of the already-reversed prefix (starts empty)
	cur := head        // first node not yet reversed
	for cur != nil {
		next := cur.Next // save successor — cur.Next is about to be overwritten
		cur.Next = prev  // flip: cur now points backwards at the reversed prefix
		prev = cur       // reversed prefix grows to include cur
		cur = next       // advance into the untouched suffix
	}
	// cur is nil, so prev holds the original tail — the new head.
	return prev
}

// buildList constructs a linked list from a slice of values and returns its head.
func buildList(vals []int) *ListNode {
	dummy := &ListNode{} // sentinel node removes the empty-list special case
	tail := dummy        // tail tracks the last node appended so far
	for _, v := range vals {
		tail.Next = &ListNode{Val: v} // append a fresh node holding v
		tail = tail.Next              // advance to the newly appended node
	}
	return dummy.Next // actual head lives after the sentinel
}

// listString renders a list as "[1 2 3]" (or "[]") for printing in main().
func listString(head *ListNode) string {
	var sb strings.Builder
	sb.WriteByte('[')
	for cur := head; cur != nil; cur = cur.Next {
		if cur != head {
			sb.WriteByte(' ') // space-separate every value after the first
		}
		fmt.Fprintf(&sb, "%d", cur.Val)
	}
	sb.WriteByte(']')
	return sb.String()
}

func main() {
	// Each approach gets a freshly built list because reversal mutates nodes
	// (except brute force, which copies — rebuilt anyway for uniformity).

	fmt.Println("=== Approach 1: Brute Force (Copy Values & Rebuild) ===")
	fmt.Println(listString(bruteForce(buildList([]int{1, 2, 3, 4, 5})))) // [5 4 3 2 1]
	fmt.Println(listString(bruteForce(buildList([]int{1, 2}))))          // [2 1]
	fmt.Println(listString(bruteForce(buildList([]int{}))))              // []

	fmt.Println("=== Approach 2: Stack of Nodes ===")
	fmt.Println(listString(stackBased(buildList([]int{1, 2, 3, 4, 5})))) // [5 4 3 2 1]
	fmt.Println(listString(stackBased(buildList([]int{1, 2}))))          // [2 1]
	fmt.Println(listString(stackBased(buildList([]int{}))))              // []

	fmt.Println("=== Approach 3: Recursion ===")
	fmt.Println(listString(recursive(buildList([]int{1, 2, 3, 4, 5})))) // [5 4 3 2 1]
	fmt.Println(listString(recursive(buildList([]int{1, 2}))))          // [2 1]
	fmt.Println(listString(recursive(buildList([]int{}))))              // []

	fmt.Println("=== Approach 4: Iterative Pointer Reversal (Optimal) ===")
	fmt.Println(listString(iterative(buildList([]int{1, 2, 3, 4, 5})))) // [5 4 3 2 1]
	fmt.Println(listString(iterative(buildList([]int{1, 2}))))          // [2 1]
	fmt.Println(listString(iterative(buildList([]int{}))))              // []
}
