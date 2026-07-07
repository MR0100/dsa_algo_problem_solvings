package main

import "fmt"

// ListNode is a singly-linked list node.
type ListNode struct {
	Val  int
	Next *ListNode
}

// ── Approach 1: Brute Force (Nested Pointer Scan) ────────────────────────────
//
// bruteForce solves Linked List Cycle by checking, for every node reached,
// whether that exact node (pointer identity) already appeared earlier in the
// traversal.
//
// Intuition:
//
//	A cycle means the traversal revisits a node it has already passed through.
//	Without extra memory, we can detect a revisit by rescanning the prefix of
//	the list: the k-th node we reach is a revisit iff it equals one of the
//	first k-1 nodes (compared by pointer, not value — values may repeat).
//
// Algorithm:
//  1. Walk cur through the list, keeping a counter i = number of nodes
//     visited before cur.
//  2. For each cur, walk probe from head over the first i nodes.
//  3. If probe == cur, cur was already visited → cycle → return true.
//  4. If cur ever becomes nil, the list terminates → return false.
//
// Time:  O(n²) — the first revisit happens within n+1 steps, and each step
//
//	rescans up to n prefix nodes.
//
// Space: O(1) — only pointers and counters.
func bruteForce(head *ListNode) bool {
	i := 0 // number of nodes already visited before cur
	for cur := head; cur != nil; cur = cur.Next {
		// Rescan the first i nodes looking for cur by pointer identity.
		probe := head
		for j := 0; j < i; j++ {
			if probe == cur {
				return true // cur was seen before → we looped back → cycle
			}
			probe = probe.Next
		}
		i++ // one more distinct node confirmed visited
	}
	// Reached nil → the list has a tail → no cycle possible.
	return false
}

// ── Approach 2: Hash Set of Visited Nodes ────────────────────────────────────
//
// hashSet solves Linked List Cycle by remembering every node pointer seen so
// far in a set; a repeat means a cycle.
//
// Intuition:
//
//	Trade memory for time: store each visited *ListNode in a set. The moment
//	we're about to visit a node already in the set, we know the next pointers
//	have looped back.
//
// Algorithm:
//  1. Create an empty set of *ListNode.
//  2. Walk the list; for each node, if it is already in the set → return true.
//  3. Otherwise insert it and continue.
//  4. Reaching nil → return false.
//
// Time:  O(n) — each node is visited at most once; set ops are O(1) average.
// Space: O(n) — the set can hold every node.
func hashSet(head *ListNode) bool {
	seen := make(map[*ListNode]bool) // set keyed by node pointer, not value
	for cur := head; cur != nil; cur = cur.Next {
		if seen[cur] {
			return true // this exact node was visited before → cycle
		}
		seen[cur] = true // mark node as visited
	}
	return false // fell off the end → no cycle
}

// ── Approach 3: Pointer Rewiring (Destructive Marking) ───────────────────────
//
// markVisited solves Linked List Cycle by rewiring each visited node's Next
// pointer to point at itself, turning "visited" into a state readable in O(1).
//
// Intuition:
//
//	If we may mutate the list, we don't need a set: after leaving a node,
//	point its Next at itself. If we ever step onto a node whose Next is
//	itself, we must have arrived via a cycle back-edge (no fresh node
//	self-loops, since we create the self-loop only after visiting).
//
// Algorithm:
//
//  1. Walk cur through the list.
//
//  2. If cur.Next == cur → cur was visited before → return true.
//
//  3. Otherwise save next, set cur.Next = cur (mark), advance to next.
//
//  4. Reaching nil → return false.
//
//     Note: destroys the list — fine for detection puzzles, not for interviews
//     that forbid modification (LeetCode #142 explicitly forbids it).
//
// Time:  O(n) — each node visited at most twice.
// Space: O(1) — the list itself stores the "visited" flag.
func markVisited(head *ListNode) bool {
	cur := head
	for cur != nil {
		if cur.Next == cur {
			return true // self-loop marker found → we re-entered a visited node
		}
		next := cur.Next // save the onward pointer before overwriting it
		cur.Next = cur   // mark this node as visited via a self-loop
		cur = next       // continue the walk on the saved pointer
	}
	return false // hit nil → list ends → no cycle
}

// ── Approach 4: Floyd's Tortoise and Hare (Optimal) ──────────────────────────
//
// floydTwoPointers solves Linked List Cycle with two pointers moving at
// different speeds; they collide iff a cycle exists.
//
// Intuition:
//
//	Send a slow pointer (1 step/turn) and a fast pointer (2 steps/turn) down
//	the list. In an acyclic list, fast reaches nil. In a cyclic list, both
//	eventually enter the cycle, and inside the cycle the gap between them
//	shrinks by exactly 1 each turn (fast gains one step per turn), so fast
//	must land exactly on slow — they cannot "jump over" each other.
//
// Algorithm:
//  1. slow = head, fast = head.
//  2. While fast != nil and fast.Next != nil:
//     slow = slow.Next; fast = fast.Next.Next.
//     If slow == fast → return true.
//  3. Loop exit means fast hit the end → return false.
//
// Time:  O(n) — slow walks at most n + cycle-length steps before collision.
// Space: O(1) — two pointers only. Satisfies the follow-up.
func floydTwoPointers(head *ListNode) bool {
	slow, fast := head, head
	// fast needs two valid hops per turn, so guard fast and fast.Next.
	for fast != nil && fast.Next != nil {
		slow = slow.Next      // tortoise: 1 step
		fast = fast.Next.Next // hare: 2 steps
		if slow == fast {
			return true // collision can only happen inside a cycle
		}
	}
	return false // fast fell off the end → acyclic
}

// ── Helpers ───────────────────────────────────────────────────────────────────

// makeCycleList builds a list from vals and, when pos >= 0, connects the tail
// to the node at index pos (mirroring LeetCode's hidden `pos` input).
func makeCycleList(vals []int, pos int) *ListNode {
	dummy := &ListNode{}
	cur := dummy
	var nodes []*ListNode
	for _, v := range vals {
		cur.Next = &ListNode{Val: v}
		cur = cur.Next
		nodes = append(nodes, cur) // remember nodes by index for the back-edge
	}
	if pos >= 0 && len(nodes) > 0 {
		cur.Next = nodes[pos] // tail.Next → node at index pos ⇒ cycle
	}
	return dummy.Next
}

func main() {
	// Official LeetCode examples: (values, pos, expected).
	examples := []struct {
		vals   []int
		pos    int
		expect bool
	}{
		{[]int{3, 2, 0, -4}, 1, true}, // Example 1
		{[]int{1, 2}, 0, true},        // Example 2
		{[]int{1}, -1, false},         // Example 3
	}

	approaches := []struct {
		name string
		fn   func(*ListNode) bool
	}{
		{"Approach 1: Brute Force (Nested Scan)", bruteForce},
		{"Approach 2: Hash Set", hashSet},
		{"Approach 3: Pointer Rewiring", markVisited},
		{"Approach 4: Floyd's Tortoise & Hare (Optimal)", floydTwoPointers},
	}

	for _, ap := range approaches {
		fmt.Printf("=== %s ===\n", ap.name)
		for i, ex := range examples {
			// Build a fresh list per run — Approach 3 mutates its input.
			head := makeCycleList(ex.vals, ex.pos)
			got := ap.fn(head)
			fmt.Printf("Example %d: head=%v pos=%d → %v (expected %v)\n",
				i+1, ex.vals, ex.pos, got, ex.expect) // expected: true, true, false
		}
		fmt.Println()
	}
}
