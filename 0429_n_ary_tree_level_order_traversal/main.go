package main

import "fmt"

// Node is the N-ary tree node used by LeetCode #429: a value plus a slice of
// child pointers (arbitrary fan-out, no fixed N per node in practice).
type Node struct {
	Val      int
	Children []*Node
}

// ── Approach 1: BFS with an Explicit Queue (Optimal) ─────────────────────────
//
// bfs solves N-ary Tree Level Order Traversal by processing the tree one full
// level at a time using a FIFO queue.
//
// Intuition:
//
//	Level order == breadth-first search. A FIFO queue visits nodes in the exact
//	order we want, but a plain BFS loses level boundaries. The classic trick:
//	before draining, snapshot the queue's current length — that count is exactly
//	how many nodes live on the current level, because only this level's nodes
//	are in the queue at that instant. Pop that many, collect their values into
//	one bucket, and push their children (the next level) as you go.
//
// Algorithm:
//  1. If root is nil, return an empty result.
//  2. Seed a queue with root.
//  3. While the queue is non-empty:
//     a. width = len(queue) — the size of the current level.
//     b. Pop `width` nodes; append each Val to a level slice; enqueue each of
//     their children (already left-to-right in Children).
//     c. Append the level slice to the answer.
//
// Time:  O(n) — every node is enqueued and dequeued exactly once.
// Space: O(n) — the queue holds at most one level; the widest level can be O(n),
//
//	plus the O(n) output.
func bfs(root *Node) [][]int {
	result := [][]int{} // final list of levels; stays [] for an empty tree
	if root == nil {
		return result // nothing to traverse
	}

	queue := []*Node{root} // FIFO seeded with the single root node
	for len(queue) > 0 {
		width := len(queue) // # nodes on this level = everything currently queued
		level := []int{}    // values collected for this level, left to right
		for i := 0; i < width; i++ {
			node := queue[0]  // front of the queue
			queue = queue[1:] // dequeue (advance the head)
			level = append(level, node.Val)
			// Children are already stored left-to-right, so enqueueing them in
			// order preserves the left-to-right requirement on the next level.
			queue = append(queue, node.Children...)
		}
		result = append(result, level) // this whole level is done
	}
	return result
}

// ── Approach 2: DFS Carrying the Depth ───────────────────────────────────────
//
// dfs solves the same problem with recursion, using the recursion depth as the
// level index instead of a queue.
//
// Intuition:
//
//	A node at depth d belongs in result[d]. If we walk the tree depth-first but
//	always pass the current depth down, we can append each value to the right
//	bucket regardless of visit order — DFS and BFS disagree on *when* a node is
//	visited, but not on *which level* it sits on. We only need to create the
//	bucket for a depth the first time we reach it (result[d] doesn't exist yet).
//
// Algorithm:
//  1. Recurse from root at depth 0.
//  2. On entering a node at depth d: if result has no slice at index d yet,
//     append a fresh empty slice (this is the first node seen at that depth).
//  3. Append the node's value to result[d].
//  4. Recurse into every child at depth d+1 (left to right).
//
// Time:  O(n) — each node visited once.
// Space: O(n) — O(h) recursion stack (h = height, up to n for a degenerate
//
//	chain) plus the O(n) output.
func dfs(root *Node) [][]int {
	result := [][]int{}
	var walk func(node *Node, depth int)
	walk = func(node *Node, depth int) {
		if node == nil {
			return
		}
		if depth == len(result) {
			// First time we descend to this depth: open a new level bucket.
			result = append(result, []int{})
		}
		result[depth] = append(result[depth], node.Val) // place value on its level
		for _, c := range node.Children {
			walk(c, depth+1) // children live one level deeper
		}
	}
	walk(root, 0)
	return result
}

// build constructs a small N-ary tree from a parent-major description so main()
// can exercise both approaches on the official examples.
//
// It returns the root. Nodes are created by value; children are wired by the
// closure below — kept explicit for clarity over a generic deserializer.

func main() {
	// ── Example 1 ────────────────────────────────────────────────────────────
	// Serialized (level order, null = end-of-children): [1,null,3,2,4,null,5,6]
	//        1
	//      / | \
	//     3  2  4
	//    / \
	//   5   6
	n5 := &Node{Val: 5}
	n6 := &Node{Val: 6}
	n3 := &Node{Val: 3, Children: []*Node{n5, n6}}
	n2 := &Node{Val: 2}
	n4 := &Node{Val: 4}
	root1 := &Node{Val: 1, Children: []*Node{n3, n2, n4}}

	fmt.Println("=== Approach 1: BFS (queue) — Example 1 ===")
	fmt.Println(bfs(root1)) // [[1] [3 2 4] [5 6]]

	fmt.Println("=== Approach 2: DFS (depth) — Example 1 ===")
	fmt.Println(dfs(root1)) // [[1] [3 2 4] [5 6]]

	// ── Example 2 ────────────────────────────────────────────────────────────
	// Serialized: [1,null,2,3,4,5,null,null,6,7,null,8,null,9,10,null,null,11,null,12,null,13,null,null,14]
	//                      1
	//        /     /    \      \
	//       2     3      4      5
	//            / \     |     / \
	//           6   7    8    9  10
	//               |    |    |
	//              11   12   13
	//               |
	//              14
	n14 := &Node{Val: 14}
	n11 := &Node{Val: 11, Children: []*Node{n14}}
	n12 := &Node{Val: 12}
	n13 := &Node{Val: 13}
	n7 := &Node{Val: 7, Children: []*Node{n11}}
	n8 := &Node{Val: 8, Children: []*Node{n12}}
	n9 := &Node{Val: 9, Children: []*Node{n13}}
	n10 := &Node{Val: 10}
	n6b := &Node{Val: 6}
	n2b := &Node{Val: 2}
	n3b := &Node{Val: 3, Children: []*Node{n6b, n7}}
	n4b := &Node{Val: 4, Children: []*Node{n8}}
	n5b := &Node{Val: 5, Children: []*Node{n9, n10}}
	root2 := &Node{Val: 1, Children: []*Node{n2b, n3b, n4b, n5b}}

	fmt.Println("=== Approach 1: BFS (queue) — Example 2 ===")
	fmt.Println(bfs(root2)) // [[1] [2 3 4 5] [6 7 8 9 10] [11 12 13] [14]]

	fmt.Println("=== Approach 2: DFS (depth) — Example 2 ===")
	fmt.Println(dfs(root2)) // [[1] [2 3 4 5] [6 7 8 9 10] [11 12 13] [14]]

	// ── Example 3: empty tree ────────────────────────────────────────────────
	fmt.Println("=== Approach 1: BFS (queue) — Example 3 (empty) ===")
	fmt.Println(bfs(nil)) // []

	fmt.Println("=== Approach 2: DFS (depth) — Example 3 (empty) ===")
	fmt.Println(dfs(nil)) // []
}
