package main

import (
	"fmt"
	"strings"
)

// Node is a multilevel doubly linked list node (LeetCode #430): besides the
// usual Prev/Next, it carries a Child pointer to a nested sublist.
type Node struct {
	Val   int
	Prev  *Node
	Next  *Node
	Child *Node
}

// ── Approach 1: Recursive Flatten (Return the Tail) ──────────────────────────
//
// recursiveFlatten flattens the list by recursively flattening each child
// sublist and splicing it in, tracking the tail of every flattened segment.
//
// Intuition:
//
//	A node with a child must have that entire child sublist inserted between it
//	and its original Next. If a helper can flatten a sublist AND hand back its
//	tail, splicing is three pointer updates: node → childHead, childTail →
//	oldNext. Recursion handles children-of-children automatically because
//	flattening a sublist recurses into its own children first.
//
// Algorithm:
//  1. Walk the current level with a cursor `cur`.
//  2. If cur.Child != nil:
//     a. Remember oldNext = cur.Next.
//     b. Recursively flatten the child, obtaining childTail.
//     c. Link cur ↔ cur.Child; clear cur.Child.
//     d. Link childTail ↔ oldNext (if oldNext exists).
//  3. Advance cur; the last non-nil cur is this level's tail.
//  4. Return that tail so a parent splice can chain onto it.
//
// Time:  O(n) — every node is visited a constant number of times.
// Space: O(d) — recursion depth equals the nesting depth d (up to n).
func recursiveFlatten(head *Node) *Node {
	flatten(head) // ignore the returned tail at the top level
	return head
}

// flatten flattens the sublist starting at `head` in place and returns its
// tail (last node). Returns nil for a nil head.
func flatten(head *Node) *Node {
	cur := head
	var tail *Node // last node seen so far on this level
	for cur != nil {
		next := cur.Next // the ORIGINAL successor, before any splice
		if cur.Child != nil {
			childTail := flatten(cur.Child) // recursively flatten the nested list

			cur.Next = cur.Child // node now points into the child list
			cur.Child.Prev = cur // back-link the child head to node
			cur.Child = nil      // child pointer must be cleared per the spec

			childTail.Next = next // stitch the child's tail to the old successor
			if next != nil {
				next.Prev = childTail // back-link if there was a successor
			}
		}
		tail = cur // cur is (so far) the furthest node on this level
		cur = next // continue with the original successor
	}
	return tail
}

// ── Approach 2: Iterative Flatten with an Explicit Stack ─────────────────────
//
// stackFlatten flattens the list without recursion by using a stack to remember
// the "detour return point" whenever we dive into a child.
//
// Intuition:
//
//	Descending into a child is a depth-first detour; the node we would have gone
//	to next (cur.Next) is where we must resume once the child branch is fully
//	flattened. A stack stores those pending Next nodes. Whenever the current
//	node has no Next, we pop the stack to reconnect the deferred successor.
//
// Algorithm:
//  1. cur = head. While cur != nil:
//     a. If cur.Child != nil: push cur.Next (if any) onto the stack, then splice
//     the child in as cur.Next and clear cur.Child.
//     b. If cur.Next == nil and the stack is non-empty: pop a node and splice it
//     as cur.Next (reconnecting a deferred branch).
//     c. Advance cur = cur.Next.
//
// Time:  O(n) — each node handled once.
// Space: O(d) — the stack holds at most one deferred Next per nesting level.
func stackFlatten(head *Node) *Node {
	if head == nil {
		return nil
	}
	stack := []*Node{} // deferred "resume here" successors, LIFO
	cur := head
	for cur != nil {
		if cur.Child != nil {
			if cur.Next != nil {
				stack = append(stack, cur.Next) // defer the current-level successor
			}
			cur.Next = cur.Child // dive into the child list
			cur.Child.Prev = cur
			cur.Child = nil // clear per the spec
		}
		if cur.Next == nil && len(stack) > 0 {
			top := stack[len(stack)-1]   // most recently deferred successor
			stack = stack[:len(stack)-1] // pop it
			cur.Next = top               // reconnect it after the finished branch
			top.Prev = cur
		}
		cur = cur.Next // move forward along the now-single-level list
	}
	return head
}

// ── helpers to build and print lists for the examples ────────────────────────

// buildLevel links a slice of values into a doubly linked list and returns its
// head; used to construct each level of the example.
func buildLevel(vals ...int) *Node {
	var head, prev *Node
	for _, v := range vals {
		n := &Node{Val: v}
		if prev == nil {
			head = n
		} else {
			prev.Next = n
			n.Prev = prev
		}
		prev = n
	}
	return head
}

// nodeAt returns the k-th node (0-indexed) of a single-level list, for wiring
// Child pointers when constructing the example.
func nodeAt(head *Node, k int) *Node {
	for k > 0 && head != nil {
		head = head.Next
		k--
	}
	return head
}

// toSlice walks Next pointers and returns the values, so we can print/verify
// the flattened order.
func toSlice(head *Node) []int {
	out := []int{}
	for head != nil {
		out = append(out, head.Val)
		head = head.Next
	}
	return out
}

// prevOK verifies the Prev pointers are consistent with Next (a proper doubly
// linked list) and that no Child pointer survived flattening.
func prevOK(head *Node) bool {
	var prev *Node
	for head != nil {
		if head.Prev != prev || head.Child != nil {
			return false
		}
		prev = head
		head = head.Next
	}
	return true
}

// buildExample1 constructs the canonical multilevel example:
//
//	1---2---3---4---5---6--NULL
//	        |
//	        7---8---9---10--NULL
//	            |
//	            11--12--NULL
//
// Flattened order: 1 2 3 7 8 11 12 9 10 4 5 6.
func buildExample1() *Node {
	top := buildLevel(1, 2, 3, 4, 5, 6)
	mid := buildLevel(7, 8, 9, 10)
	low := buildLevel(11, 12)

	nodeAt(top, 2).Child = mid // node 3's child is the 7..10 list
	nodeAt(mid, 1).Child = low // node 8's child is the 11..12 list
	return top
}

// buildExample2 constructs: 1---2---NULL with 1's child = 3---NULL.
// Flattened order: 1 3 2.
func buildExample2() *Node {
	top := buildLevel(1, 2)
	top.Child = buildLevel(3) // node 1's child is the single-node list [3]
	return top
}

func fmtSlice(s []int) string {
	parts := make([]string, len(s))
	for i, v := range s {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return "[" + strings.Join(parts, ",") + "]"
}

func main() {
	fmt.Println("=== Approach 1: Recursive Flatten ===")
	r1 := recursiveFlatten(buildExample1())
	fmt.Println(fmtSlice(toSlice(r1)), "prev/child ok:", prevOK(r1)) // [1,2,3,7,8,11,12,9,10,4,5,6] prev/child ok: true
	r1b := recursiveFlatten(buildExample2())
	fmt.Println(fmtSlice(toSlice(r1b)), "prev/child ok:", prevOK(r1b)) // [1,3,2] prev/child ok: true
	fmt.Println(fmtSlice(toSlice(recursiveFlatten(nil))))              // []

	fmt.Println("=== Approach 2: Iterative Flatten (stack) ===")
	r2 := stackFlatten(buildExample1())
	fmt.Println(fmtSlice(toSlice(r2)), "prev/child ok:", prevOK(r2)) // [1,2,3,7,8,11,12,9,10,4,5,6] prev/child ok: true
	r2b := stackFlatten(buildExample2())
	fmt.Println(fmtSlice(toSlice(r2b)), "prev/child ok:", prevOK(r2b)) // [1,3,2] prev/child ok: true
	fmt.Println(fmtSlice(toSlice(stackFlatten(nil))))                  // []
}
