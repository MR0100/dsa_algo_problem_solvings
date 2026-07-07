package main

import "fmt"

// TreeNode is a binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func printList(root *TreeNode) []int {
	var res []int
	for root != nil {
		res = append(res, root.Val)
		root = root.Right
	}
	return res
}

// ── Approach 1: Collect Preorder then Rebuild ─────────────────────────────────
//
// flatten solves Flatten Binary Tree to Linked List by collecting all nodes
// in preorder then rewiring them.
//
// Intuition:
//   Preorder traversal (root → left → right) gives the order of the flattened
//   list. Collect all nodes, then wire each node's Right to the next and set
//   Left to nil.
//
// Time:  O(n)
// Space: O(n) — nodes slice.
func flatten(root *TreeNode) {
	var nodes []*TreeNode
	var preorder func(node *TreeNode)
	preorder = func(node *TreeNode) {
		if node == nil {
			return
		}
		nodes = append(nodes, node)
		preorder(node.Left)
		preorder(node.Right)
	}
	preorder(root)

	for i := 0; i < len(nodes)-1; i++ {
		nodes[i].Left = nil
		nodes[i].Right = nodes[i+1]
	}
}

// ── Approach 2: Morris-Like In-Place (O(1) Space) ────────────────────────────
//
// flattenInPlace solves Flatten Binary Tree to Linked List in-place.
//
// Intuition:
//   For each node with a left child:
//   1. Find the rightmost node of the left subtree (preorder predecessor).
//   2. Attach the current right subtree to that node's right.
//   3. Move the left subtree to the right, set left to nil.
//   Repeat until no node has a left child.
//
// Time:  O(n) — each node visited at most twice.
// Space: O(1)
func flattenInPlace(root *TreeNode) {
	curr := root
	for curr != nil {
		if curr.Left != nil {
			// find rightmost node of left subtree
			rightmost := curr.Left
			for rightmost.Right != nil {
				rightmost = rightmost.Right
			}
			// attach current right to rightmost
			rightmost.Right = curr.Right
			// move left subtree to right
			curr.Right = curr.Left
			curr.Left = nil
		}
		curr = curr.Right
	}
}

// ── Approach 3: Recursive Post-Order (Reverse Preorder) ──────────────────────
//
// flattenReverse solves Flatten Binary Tree to Linked List using a reverse
// preorder traversal (right → left → root) with a `prev` pointer.
//
// Intuition:
//   Process nodes in reverse preorder: right, then left, then root.
//   Maintain a `prev` pointer to the previously processed node.
//   Set current node's right = prev, left = nil, then update prev = current.
//   This wires up the list in reverse order — correct because we go right→left→root.
//
// Time:  O(n)
// Space: O(h)
func flattenReverse(root *TreeNode) {
	var prev *TreeNode
	var dfs func(node *TreeNode)
	dfs = func(node *TreeNode) {
		if node == nil {
			return
		}
		dfs(node.Right)  // process right first
		dfs(node.Left)
		node.Right = prev // wire up
		node.Left = nil
		prev = node
	}
	dfs(root)
}

func main() {
	// [1,2,5,3,4,null,6]
	build1 := func() *TreeNode {
		return &TreeNode{Val: 1,
			Left: &TreeNode{Val: 2,
				Left: &TreeNode{Val: 3}, Right: &TreeNode{Val: 4}},
			Right: &TreeNode{Val: 5, Right: &TreeNode{Val: 6}},
		}
	}

	fmt.Println("=== Approach 1: Preorder Collect ===")
	t1 := build1()
	flatten(t1)
	fmt.Printf("got=%v  expected [1 2 3 4 5 6]\n", printList(t1))

	fmt.Println("=== Approach 2: Morris-Like In-Place ===")
	t2 := build1()
	flattenInPlace(t2)
	fmt.Printf("got=%v  expected [1 2 3 4 5 6]\n", printList(t2))

	fmt.Println("=== Approach 3: Reverse Preorder ===")
	t3 := build1()
	flattenReverse(t3)
	fmt.Printf("got=%v  expected [1 2 3 4 5 6]\n", printList(t3))

	// edge: single node
	t4 := &TreeNode{Val: 0}
	flatten(t4)
	fmt.Printf("single node got=%v  expected [0]\n", printList(t4))
}
