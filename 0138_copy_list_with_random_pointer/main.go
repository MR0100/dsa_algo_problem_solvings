package main

import (
	"fmt"
	"strings"
)

// Node is a linked-list node with an extra random pointer, as defined by LeetCode.
type Node struct {
	Val    int
	Next   *Node
	Random *Node
}

// ── Approach 1: Hash Map (Two Pass) ──────────────────────────────────────────
//
// hashMapTwoPass solves Copy List with Random Pointer using an old→new node map.
//
// Intuition:
//
//	The random pointer may jump anywhere, so while cloning node-by-node we may
//	need a copy that doesn't exist yet. Fix: first create ALL copies and record
//	the original→copy mapping; then wire next/random by translating each
//	original pointer through the map.
//
// Algorithm:
//  1. Pass 1: for every original node, create a bare copy (value only) and
//     store oldToNew[orig] = copy.
//  2. Pass 2: for every original node, set
//     copy.Next   = oldToNew[orig.Next]
//     copy.Random = oldToNew[orig.Random]
//     (a nil original pointer maps to nil since the map returns the zero value).
//  3. Return oldToNew[head].
//
// Time:  O(n) — two linear passes with O(1) map operations.
// Space: O(n) — the map holds one entry per node (excluding the output itself).
func hashMapTwoPass(head *Node) *Node {
	if head == nil {
		return nil // empty list clones to empty list
	}
	oldToNew := make(map[*Node]*Node) // original node → its copy

	// Pass 1: create all copies so every target of next/random already exists.
	for curr := head; curr != nil; curr = curr.Next {
		oldToNew[curr] = &Node{Val: curr.Val} // copy carries only the value for now
	}

	// Pass 2: wire pointers by translating originals through the map.
	for curr := head; curr != nil; curr = curr.Next {
		copyNode := oldToNew[curr]
		copyNode.Next = oldToNew[curr.Next]     // nil key → nil value, handles tail
		copyNode.Random = oldToNew[curr.Random] // nil key → nil value, handles null random
	}

	return oldToNew[head] // the copy of the head is the new list
}

// ── Approach 2: Recursion + Memoization ──────────────────────────────────────
//
// recursiveMemo solves Copy List with Random Pointer by treating the list as a
// graph and deep-cloning it with memoized DFS.
//
// Intuition:
//
//	Each node has two outgoing edges (next, random), so the structure is a
//	graph. Cloning a graph = DFS where a memo map both prevents infinite
//	recursion on cycles and guarantees each original maps to exactly one copy.
//
// Algorithm:
//  1. If node is nil, return nil.
//  2. If node is already in the memo, return its existing copy (cycle/shared node).
//  3. Otherwise create the copy, memoize it BEFORE recursing (breaks cycles),
//     then recursively clone Next and Random.
//
// Time:  O(n) — every node is cloned exactly once; later visits hit the memo.
// Space: O(n) — memo map plus up to O(n) recursion stack.
func recursiveMemo(head *Node) *Node {
	memo := make(map[*Node]*Node) // original → copy, shared across the recursion
	var clone func(node *Node) *Node
	clone = func(node *Node) *Node {
		if node == nil {
			return nil // base case: null pointer clones to null
		}
		if copyNode, ok := memo[node]; ok {
			return copyNode // already cloned (revisited via random or a cycle)
		}
		copyNode := &Node{Val: node.Val}
		memo[node] = copyNode                // memoize BEFORE recursing to break cycles
		copyNode.Next = clone(node.Next)     // deep-clone the next edge
		copyNode.Random = clone(node.Random) // deep-clone the random edge
		return copyNode
	}
	return clone(head)
}

// ── Approach 3: Interleaving (Optimal, O(1) extra space) ─────────────────────
//
// interleaving solves Copy List with Random Pointer by weaving copies into the
// original list, using the list itself as the old→new "map".
//
// Intuition:
//
//	The hash map exists only to answer "given an original node, where is its
//	copy?". If we insert each copy IMMEDIATELY AFTER its original
//	(A→A'→B→B'→C→C'), that question is answered by orig.Next — no map needed.
//
// Algorithm:
//  1. Weave: for each original node X, insert copy X' right after X.
//  2. Assign randoms: for each original X, if X.Random != nil then
//     X'.Random = X.Random.Next (the copy of X.Random sits right after it).
//  3. Unweave: split the mixed list back into original and copy lists,
//     restoring the original list intact.
//
// Time:  O(n) — three linear passes.
// Space: O(1) — only a few pointers beyond the output list itself.
func interleaving(head *Node) *Node {
	if head == nil {
		return nil // nothing to copy
	}

	// Step 1: weave copies into the original list: A→A'→B→B'→...
	for curr := head; curr != nil; curr = curr.Next.Next {
		copyNode := &Node{Val: curr.Val}
		copyNode.Next = curr.Next // copy points at the rest of the list
		curr.Next = copyNode      // original points at its copy
		// advance by two: skip over the copy we just inserted
	}

	// Step 2: set random pointers on the copies using the weave invariant.
	for curr := head; curr != nil; curr = curr.Next.Next {
		if curr.Random != nil {
			curr.Next.Random = curr.Random.Next // copy of X.Random is X.Random.Next
		}
	}

	// Step 3: unweave — detach the copies and restore the original list.
	newHead := head.Next
	for curr := head; curr != nil; {
		copyNode := curr.Next
		curr.Next = copyNode.Next // restore original's next
		if copyNode.Next != nil {
			copyNode.Next = copyNode.Next.Next // link copy to the next copy
		}
		curr = curr.Next // move to the next original node
	}

	return newHead
}

// ── Test helpers ─────────────────────────────────────────────────────────────

// buildList constructs a random-pointer list from LeetCode's [[val, randomIndex]]
// encoding. randoms[i] == -1 means the random pointer is null.
func buildList(vals []int, randoms []int) *Node {
	if len(vals) == 0 {
		return nil
	}
	nodes := make([]*Node, len(vals))
	for i, v := range vals {
		nodes[i] = &Node{Val: v} // create all nodes first
	}
	for i := range nodes {
		if i+1 < len(nodes) {
			nodes[i].Next = nodes[i+1] // chain next pointers
		}
		if randoms[i] >= 0 {
			nodes[i].Random = nodes[randoms[i]] // resolve random index to a node
		}
	}
	return nodes[0]
}

// serialize renders a list back into LeetCode's [[val,randomIndex],...] form so
// results can be compared against the expected output literally.
func serialize(head *Node) string {
	index := make(map[*Node]int) // node → position, to express randoms as indices
	for curr, i := head, 0; curr != nil; curr, i = curr.Next, i+1 {
		index[curr] = i
	}
	var parts []string
	for curr := head; curr != nil; curr = curr.Next {
		if curr.Random == nil {
			parts = append(parts, fmt.Sprintf("[%d,null]", curr.Val))
		} else {
			parts = append(parts, fmt.Sprintf("[%d,%d]", curr.Val, index[curr.Random]))
		}
	}
	return "[" + strings.Join(parts, ",") + "]"
}

// isDeepCopy verifies no node of the clone is shared with the original list.
func isDeepCopy(orig, clone *Node) bool {
	origSet := make(map[*Node]bool)
	for curr := orig; curr != nil; curr = curr.Next {
		origSet[curr] = true // collect all original node pointers
	}
	for curr := clone; curr != nil; curr = curr.Next {
		if origSet[curr] || (curr.Random != nil && origSet[curr.Random]) {
			return false // clone reuses an original node → not a deep copy
		}
	}
	return true
}

func main() {
	type example struct {
		vals    []int
		randoms []int
	}
	examples := []example{
		{[]int{7, 13, 11, 10, 1}, []int{-1, 0, 4, 2, 0}}, // expected [[7,null],[13,0],[11,4],[10,2],[1,0]]
		{[]int{1, 2}, []int{1, 1}},                       // expected [[1,1],[2,1]]
		{[]int{3, 3, 3}, []int{-1, 0, -1}},               // expected [[3,null],[3,0],[3,null]]
	}

	run := func(name string, solve func(*Node) *Node) {
		fmt.Printf("=== %s ===\n", name)
		for _, ex := range examples {
			orig := buildList(ex.vals, ex.randoms) // fresh list per approach (interleaving mutates+restores)
			clone := solve(orig)
			fmt.Printf("got=%s  deepCopy=%v  originalIntact=%s\n",
				serialize(clone), isDeepCopy(orig, clone), serialize(orig))
		}
	}

	run("Approach 1: Hash Map (Two Pass)", hashMapTwoPass)
	// expected [[7,null],[13,0],[11,4],[10,2],[1,0]] / [[1,1],[2,1]] / [[3,null],[3,0],[3,null]], all deepCopy=true
	run("Approach 2: Recursion + Memoization", recursiveMemo)
	// expected same outputs, all deepCopy=true
	run("Approach 3: Interleaving (Optimal)", interleaving)
	// expected same outputs, all deepCopy=true, originals restored intact
}
