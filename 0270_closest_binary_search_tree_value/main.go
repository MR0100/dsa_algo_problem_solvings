package main

import (
	"fmt"
	"math"
)

// TreeNode is a standard binary-search-tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// ── Approach 1: In-order Traversal + Linear Scan ─────────────────────────────
//
// inorderScan solves Closest BST Value by flattening the tree into a sorted
// slice and scanning for the value nearest to target.
//
// Intuition:
//
//	Ignore the BST shape entirely: collect every node's value, then pick the
//	one with the smallest absolute distance to target. Simple, always correct.
//
// Algorithm:
//  1. In-order traverse to collect all values (sorted, though sorting is not
//     required for correctness).
//  2. Track the value with minimum |val - target|.
//  3. On ties, prefer the smaller value (per LeetCode's tie rule).
//
// Time:  O(n) — visits every node.
// Space: O(n) — the collected slice plus recursion stack.
func inorderScan(root *TreeNode, target float64) int {
	var vals []int
	var walk func(n *TreeNode)
	walk = func(n *TreeNode) {
		if n == nil {
			return
		}
		walk(n.Left)               // left subtree (smaller values)
		vals = append(vals, n.Val) // visit node
		walk(n.Right)              // right subtree (larger values)
	}
	walk(root)

	closest := vals[0]
	best := math.Abs(float64(closest) - target) // smallest distance so far
	for _, v := range vals[1:] {
		d := math.Abs(float64(v) - target)
		// strictly-smaller distance wins; equal distance keeps the smaller
		// value because we scan in ascending order and only replace on <.
		if d < best {
			best = d
			closest = v
		}
	}
	return closest
}

// ── Approach 2: Iterative BST Descent (Optimal) ──────────────────────────────
//
// bstDescent solves Closest BST Value by walking down the tree, using the BST
// ordering to head toward target and never revisiting nodes.
//
// Intuition:
//
//	At each node the BST property tells us which half of the values lie ahead:
//	go left if target < node.Val, else go right. Along that single root-to-leaf
//	path we pass the candidates that can possibly be closest, so tracking the
//	best on the way down suffices — no need to see the whole tree.
//
// Algorithm:
//  1. closest = root.Val.
//  2. While node != nil: update closest if this node is nearer (ties -> smaller
//     value). Then move left if target < node.Val, else right.
//  3. Return closest.
//
// Time:  O(h) — h = tree height (O(log n) balanced, O(n) skewed).
// Space: O(1) — no recursion, no extra storage.
func bstDescent(root *TreeNode, target float64) int {
	closest := root.Val
	node := root
	for node != nil {
		// Prefer this node if strictly closer, or equally close but smaller
		// (LeetCode breaks ties toward the smaller value).
		curD := math.Abs(float64(node.Val) - target)
		bestD := math.Abs(float64(closest) - target)
		if curD < bestD || (curD == bestD && node.Val < closest) {
			closest = node.Val
		}
		if target < float64(node.Val) {
			node = node.Left // target is smaller -> smaller values are left
		} else {
			node = node.Right // target is larger -> larger values are right
		}
	}
	return closest
}

func main() {
	// Example 1 tree: [4,2,5,1,3], target 3.714286
	root := &TreeNode{
		Val:   4,
		Left:  &TreeNode{Val: 2, Left: &TreeNode{Val: 1}, Right: &TreeNode{Val: 3}},
		Right: &TreeNode{Val: 5},
	}
	// Example 2 tree: [1], target 4.428571
	single := &TreeNode{Val: 1}

	fmt.Println("=== Approach 1: In-order Scan ===")
	fmt.Println(inorderScan(root, 3.714286))   // expected 4
	fmt.Println(inorderScan(single, 4.428571)) // expected 1

	fmt.Println("=== Approach 2: BST Descent (Optimal) ===")
	fmt.Println(bstDescent(root, 3.714286))   // expected 4
	fmt.Println(bstDescent(single, 4.428571)) // expected 1
}
