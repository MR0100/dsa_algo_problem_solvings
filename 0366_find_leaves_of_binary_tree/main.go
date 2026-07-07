package main

import "fmt"

// TreeNode is the standard LeetCode binary-tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// ── Approach 1: Repeated Leaf Stripping (Brute Force) ────────────────────────
//
// repeatedStripping solves Find Leaves of Binary Tree by literally simulating
// the problem statement: find every current leaf, record its value, physically
// detach it from the tree, and repeat until the tree is empty.
//
// Intuition:
//
//	The statement is procedural — "collect leaves, remove them, repeat". We can
//	obey it directly. Each pass does two things: (1) collect the values of all
//	nodes that are currently leaves, and (2) sever the parent→leaf pointers so
//	those leaves disappear before the next pass. The root itself needs special
//	handling because it has no parent; once the root becomes a leaf, the pass
//	that collects it also empties the tree.
//
// Algorithm:
//  1. While the tree is non-empty:
//     a. If the root is a leaf, append [root.Val] and set the tree to nil.
//     b. Otherwise walk the tree; for each node whose child is a leaf, collect
//     that child's value and null out the pointer to it.
//     c. Append this pass's collected values as one group.
//
// Time:  O(n * h) worst case — up to h ≈ n passes, each pass touches O(n) nodes
//
//	(a degenerate "stick" tree strips one leaf per pass).
//
// Space: O(h) recursion stack per pass plus O(n) for the output.
func repeatedStripping(root *TreeNode) [][]int {
	var result [][]int // groups of values, one per stripping pass

	// removeLeaves detaches every current leaf below node, appending their
	// values to *collected. It returns the (possibly nil) node to keep — nil if
	// node itself was a leaf and should be removed by its caller.
	var removeLeaves func(node *TreeNode, collected *[]int) *TreeNode
	removeLeaves = func(node *TreeNode, collected *[]int) *TreeNode {
		if node == nil {
			return nil
		}
		if node.Left == nil && node.Right == nil {
			// node is a leaf: record it and tell the parent to drop it.
			*collected = append(*collected, node.Val)
			return nil
		}
		// Recurse first, then re-link the surviving children.
		node.Left = removeLeaves(node.Left, collected)
		node.Right = removeLeaves(node.Right, collected)
		return node
	}

	for root != nil {
		var pass []int // values collected in this single stripping pass
		root = removeLeaves(root, &pass)
		result = append(result, pass) // this pass forms one output group
	}
	return result
}

// ── Approach 2: Postorder Height Grouping (Optimal) ──────────────────────────
//
// heightGrouping solves Find Leaves of Binary Tree by recognising that the pass
// in which a node is removed equals its *height* (distance to its deepest leaf).
//
// Intuition:
//
//	A node is stripped on pass k exactly when the longest downward path from it
//	has length k (0-indexed): true leaves have height 0 and go first, a node
//	whose deepest child leaf is one level down has height 1, and so on. Height
//	is computable bottom-up in a single postorder DFS: height(node) =
//	1 + max(height(left), height(right)). We use that height directly as the
//	index of the output group to append the node's value to — one traversal,
//	no tree mutation, no repeated passes.
//
// Algorithm:
//  1. DFS postorder. For a nil child return height −1 so a leaf gets height 0.
//  2. h = 1 + max(dfs(left), dfs(right)).
//  3. If result has no group at index h yet, create it; append node.Val to
//     result[h].
//  4. Return h to the parent.
//
// Time:  O(n) — each node visited once.
// Space: O(h) recursion stack (O(n) worst case) plus O(n) output.
func heightGrouping(root *TreeNode) [][]int {
	var result [][]int // result[h] holds every node of height h

	// dfs returns the height of node (leaf = 0, nil = -1) and files node into
	// the group indexed by its height along the way.
	var dfs func(node *TreeNode) int
	dfs = func(node *TreeNode) int {
		if node == nil {
			return -1 // so a leaf computes 1 + max(-1,-1) = 0
		}
		left := dfs(node.Left)    // height of left subtree
		right := dfs(node.Right)  // height of right subtree
		h := 1 + max(left, right) // this node's height
		if h == len(result) {     // first node discovered at this height
			result = append(result, []int{}) // open a new group
		}
		result[h] = append(result[h], node.Val) // file node under its height
		return h
	}

	dfs(root)
	return result
}

// max returns the larger of two ints (Go 1.21+ has a builtin, but we spell it
// out for clarity and older toolchains).
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// buildTree builds a binary tree from a level-order slice using nil for missing
// nodes (LeetCode's array format). Used only to construct the examples.
func buildTree(vals []interface{}) *TreeNode {
	if len(vals) == 0 || vals[0] == nil {
		return nil
	}
	nodes := make([]*TreeNode, len(vals))
	for i, v := range vals {
		if v != nil {
			nodes[i] = &TreeNode{Val: v.(int)}
		}
	}
	// Level-order wiring via a queue: each dequeued node claims the next one or
	// two positions of the compact slice as its children.
	queue := []*TreeNode{nodes[0]}
	i := 1
	for len(queue) > 0 && i < len(vals) {
		node := queue[0]
		queue = queue[1:]
		if node == nil {
			continue
		}
		if i < len(vals) { // left child
			node.Left = nodes[i]
			queue = append(queue, nodes[i])
			i++
		}
		if i < len(vals) { // right child
			node.Right = nodes[i]
			queue = append(queue, nodes[i])
			i++
		}
	}
	return nodes[0]
}

func main() {
	// Example 1: root = [1,2,3,4,5] → [[4,5,3],[2],[1]]
	// Example 2: root = [1]         → [[1]]
	// (Approach 1 mutates the tree, so build a fresh copy for each call.)

	fmt.Println("=== Approach 1: Repeated Leaf Stripping (Brute Force) ===")
	fmt.Println(repeatedStripping(buildTree([]interface{}{1, 2, 3, 4, 5}))) // expected [[4 5 3] [2] [1]]
	fmt.Println(repeatedStripping(buildTree([]interface{}{1})))             // expected [[1]]

	fmt.Println("=== Approach 2: Postorder Height Grouping (Optimal) ===")
	fmt.Println(heightGrouping(buildTree([]interface{}{1, 2, 3, 4, 5}))) // expected [[4 5 3] [2] [1]]
	fmt.Println(heightGrouping(buildTree([]interface{}{1})))             // expected [[1]]
}
