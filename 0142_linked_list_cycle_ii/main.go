package main

import "fmt"

// ListNode is a singly-linked list node.
type ListNode struct {
	Val  int
	Next *ListNode
}

// ── Approach 1: Brute Force (Nested Pointer Scan) ────────────────────────────
//
// bruteForce solves Linked List Cycle II by finding the first node that is
// visited twice — that node is, by definition, the cycle entrance.
//
// Intuition:
//
//	Walking from head, every node before the cycle entrance is visited exactly
//	once; the entrance is the first node the walk returns to. So the first
//	revisited node IS the answer. Detect a revisit without extra memory by
//	rescanning the traversal prefix with pointer comparisons.
//
// Algorithm:
//  1. Walk cur through the list with counter i = nodes visited before cur.
//  2. For each cur, scan probe over the first i nodes from head.
//  3. If probe == cur → cur is the first revisited node → return cur.
//  4. If cur becomes nil → no cycle → return nil.
//
// Time:  O(n²) — first revisit occurs within n+1 steps; each step rescans ≤ n nodes.
// Space: O(1) — pointers and counters only.
func bruteForce(head *ListNode) *ListNode {
	i := 0 // number of nodes already visited before cur
	for cur := head; cur != nil; cur = cur.Next {
		// Rescan the first i nodes; a pointer match means cur is revisited.
		probe := head
		for j := 0; j < i; j++ {
			if probe == cur {
				return cur // first node seen twice = cycle entrance
			}
			probe = probe.Next
		}
		i++ // cur confirmed as a newly visited node
	}
	return nil // reached the tail → acyclic
}

// ── Approach 2: Hash Set of Visited Nodes ────────────────────────────────────
//
// hashSet solves Linked List Cycle II by storing visited node pointers in a
// set and returning the first node found already present.
//
// Intuition:
//
//	Same "first revisited node = entrance" fact as Approach 1, but the
//	membership test becomes O(1) with a hash set keyed by pointer.
//
// Algorithm:
//  1. seen = empty set of *ListNode.
//  2. Walk the list; if cur ∈ seen → return cur (entrance).
//  3. Else insert cur and advance.
//  4. Reaching nil → return nil.
//
// Time:  O(n) — one pass, O(1) average per set operation.
// Space: O(n) — set may store every node.
func hashSet(head *ListNode) *ListNode {
	seen := make(map[*ListNode]bool) // visited set keyed by node pointer
	for cur := head; cur != nil; cur = cur.Next {
		if seen[cur] {
			return cur // first repeat = where the tail's back-edge lands
		}
		seen[cur] = true // mark visited
	}
	return nil // no repeat before nil → no cycle
}

// ── Approach 3: Sentinel Rewiring (Destructive Marking) ──────────────────────
//
// sentinelRewiring solves Linked List Cycle II by redirecting each visited
// node's Next to a shared sentinel node; the first node already pointing at
// the sentinel is the entrance.
//
// Intuition:
//
//	Mutating the list lets it store its own "visited" flag: after leaving a
//	node, aim its Next at a sentinel. Arriving at a node whose Next is the
//	sentinel means we arrived via the cycle's back-edge, and the first such
//	arrival is exactly the cycle entrance.
//
//	NOTE: the problem statement says "Do not modify the linked list", so this
//	approach is educational — it shows the O(1)-space idea Floyd's achieves
//	without mutation.
//
// Algorithm:
//  1. Create sentinel node S.
//  2. Walk cur; if cur.Next == S → return cur.
//  3. Else save next, set cur.Next = S, advance.
//  4. nil reached → return nil.
//
// Time:  O(n) — each node entered at most twice.
// Space: O(1) — one sentinel node.
func sentinelRewiring(head *ListNode) *ListNode {
	sentinel := &ListNode{} // unique marker no real node ever points to
	cur := head
	for cur != nil {
		if cur.Next == sentinel {
			return cur // already marked → first revisited node → entrance
		}
		next := cur.Next    // save onward pointer before overwriting
		cur.Next = sentinel // mark cur as visited
		cur = next          // continue on the saved pointer
	}
	return nil // fell off the tail → no cycle
}

// ── Approach 4: Floyd's Cycle Detection (Optimal) ────────────────────────────
//
// floydCycle solves Linked List Cycle II in two phases: detect a collision
// with fast/slow pointers, then walk two pointers at equal speed from head
// and the collision point — they meet at the entrance.
//
// Intuition:
//
//	Let F = distance head→entrance, C = cycle length, and let the pointers
//	collide a distance k past the entrance (inside the cycle). At collision:
//	  slow travelled  F + k
//	  fast travelled  F + k + m·C   (m ≥ 1 full extra laps)
//	fast moves twice as fast, so 2(F + k) = F + k + m·C  ⇒  F + k = m·C
//	⇒  F = m·C − k. Meaning: from the collision point, walking F more steps
//	lands exactly on the entrance (k + F = m·C ≡ 0 (mod C)). So start one
//	pointer at head, keep one at the collision, step both by 1 — they meet
//	at the entrance after F steps.
//
// Algorithm:
//
//	Phase 1 (detect):
//	  1. slow = fast = head.
//	  2. While fast and fast.Next non-nil: slow += 1 step, fast += 2 steps.
//	  3. If they collide → go to phase 2. If fast hits nil → return nil.
//	Phase 2 (locate entrance):
//	  4. ptr1 = head, ptr2 = collision node.
//	  5. While ptr1 != ptr2: advance both by 1.
//	  6. Return ptr1 (== ptr2) — the entrance.
//
// Time:  O(n) — each phase is at most a constant number of passes.
// Space: O(1) — four pointers. Satisfies the follow-up without mutation.
func floydCycle(head *ListNode) *ListNode {
	slow, fast := head, head
	// Phase 1: advance until collision or end of list.
	for fast != nil && fast.Next != nil {
		slow = slow.Next      // tortoise: 1 step
		fast = fast.Next.Next // hare: 2 steps
		if slow == fast {
			// Phase 2: F = m·C − k ⇒ head and collision point are the same
			// distance (mod C) from the entrance.
			ptr1, ptr2 := head, slow
			for ptr1 != ptr2 {
				ptr1 = ptr1.Next // both move at speed 1
				ptr2 = ptr2.Next
			}
			return ptr1 // meeting point = cycle entrance
		}
	}
	return nil // fast fell off the end → acyclic
}

// ── Helpers ───────────────────────────────────────────────────────────────────

// makeCycleList builds a list from vals; when pos >= 0 the tail is linked back
// to the node at index pos. It also returns the node slice so results can be
// reported as an index even after destructive approaches ran.
func makeCycleList(vals []int, pos int) (*ListNode, []*ListNode) {
	dummy := &ListNode{}
	cur := dummy
	var nodes []*ListNode
	for _, v := range vals {
		cur.Next = &ListNode{Val: v}
		cur = cur.Next
		nodes = append(nodes, cur) // keep index → node mapping
	}
	if pos >= 0 && len(nodes) > 0 {
		cur.Next = nodes[pos] // create the back-edge
	}
	return dummy.Next, nodes
}

// describe renders a returned node as "index i (value v)" or "nil".
func describe(node *ListNode, nodes []*ListNode) string {
	if node == nil {
		return "nil (no cycle)"
	}
	for i, n := range nodes {
		if n == node { // pointer identity locates the original index
			return fmt.Sprintf("index %d (value %d)", i, n.Val)
		}
	}
	return "unknown node" // unreachable for correct solutions
}

func main() {
	// Official LeetCode examples: (values, pos, expected description).
	examples := []struct {
		vals   []int
		pos    int
		expect string
	}{
		{[]int{3, 2, 0, -4}, 1, "index 1 (value 2)"}, // Example 1
		{[]int{1, 2}, 0, "index 0 (value 1)"},        // Example 2
		{[]int{1}, -1, "nil (no cycle)"},             // Example 3
	}

	approaches := []struct {
		name string
		fn   func(*ListNode) *ListNode
	}{
		{"Approach 1: Brute Force (Nested Scan)", bruteForce},
		{"Approach 2: Hash Set", hashSet},
		{"Approach 3: Sentinel Rewiring", sentinelRewiring},
		{"Approach 4: Floyd's Cycle Detection (Optimal)", floydCycle},
	}

	for _, ap := range approaches {
		fmt.Printf("=== %s ===\n", ap.name)
		for i, ex := range examples {
			// Fresh list per run — Approach 3 mutates its input.
			head, nodes := makeCycleList(ex.vals, ex.pos)
			got := ap.fn(head)
			fmt.Printf("Example %d: head=%v pos=%d → %s (expected %s)\n",
				i+1, ex.vals, ex.pos, describe(got, nodes), ex.expect)
		}
		fmt.Println()
	}
}
