package main

import "fmt"

// TreeNode is a binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// ── Approach 1: Recursive ─────────────────────────────────────────────────────
//
// inorderTraversal solves Binary Tree Inorder Traversal recursively.
//
// Intuition:
//   Inorder = left, root, right. Recurse left, append root.Val, recurse right.
//
// Time:  O(n) — each node visited once.
// Space: O(h) — call stack depth h (height of tree).
func inorderRecursive(root *TreeNode) []int {
	var result []int
	var inorder func(node *TreeNode)
	inorder = func(node *TreeNode) {
		if node == nil {
			return
		}
		inorder(node.Left)
		result = append(result, node.Val)
		inorder(node.Right)
	}
	inorder(root)
	return result
}

// ── Approach 2: Iterative (Stack) ────────────────────────────────────────────
//
// inorderIterative solves Binary Tree Inorder Traversal using an explicit stack.
//
// Intuition:
//   Simulate the recursive call stack. Push all left nodes first. When we
//   reach nil, pop the stack (this is the node to visit), then go right.
//
// Algorithm:
//   stack = []; curr = root
//   while curr != nil or stack not empty:
//     while curr != nil: push curr; curr = curr.Left
//     curr = pop; result.append(curr.Val); curr = curr.Right
//
// Time:  O(n)
// Space: O(h)
func inorderIterative(root *TreeNode) []int {
	var result []int
	stack := []*TreeNode{}
	curr := root

	for curr != nil || len(stack) > 0 {
		// go all the way left
		for curr != nil {
			stack = append(stack, curr)
			curr = curr.Left
		}
		// pop and visit
		curr = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		result = append(result, curr.Val)
		// then go right
		curr = curr.Right
	}
	return result
}

// ── Approach 3: Morris Traversal (O(1) Space) ────────────────────────────────
//
// inorderMorris solves Binary Tree Inorder Traversal with O(1) extra space
// by temporarily threading right pointers (Morris threading).
//
// Intuition:
//   For each node, if it has no left child, visit it and go right.
//   If it has a left child, find its inorder predecessor (rightmost node in
//   the left subtree). If the predecessor's right is nil, thread it to current
//   node and go left. If the predecessor's right is current (already threaded),
//   unthread, visit current, go right.
//
// Time:  O(n) — each edge traversed at most twice.
// Space: O(1) — no stack or recursion.
func inorderMorris(root *TreeNode) []int {
	var result []int
	curr := root

	for curr != nil {
		if curr.Left == nil {
			// no left subtree: visit and go right
			result = append(result, curr.Val)
			curr = curr.Right
		} else {
			// find inorder predecessor (rightmost in left subtree)
			pred := curr.Left
			for pred.Right != nil && pred.Right != curr {
				pred = pred.Right
			}
			if pred.Right == nil {
				// thread: predecessor.right → curr; go left
				pred.Right = curr
				curr = curr.Left
			} else {
				// unthread; visit curr; go right
				pred.Right = nil
				result = append(result, curr.Val)
				curr = curr.Right
			}
		}
	}
	return result
}

func main() {
	// Build tree: [1,null,2,3] → 1 → right=2, 2.left=3
	root1 := &TreeNode{Val: 1}
	root1.Right = &TreeNode{Val: 2}
	root1.Right.Left = &TreeNode{Val: 3}

	root2 := (*TreeNode)(nil) // empty tree

	root3 := &TreeNode{Val: 1}

	fmt.Println("=== Approach 1: Recursive ===")
	fmt.Printf("tree=[1,null,2,3]  got=%v  expected [1 3 2]\n", inorderRecursive(root1))
	fmt.Printf("tree=[]  got=%v  expected []\n", inorderRecursive(root2))
	fmt.Printf("tree=[1]  got=%v  expected [1]\n", inorderRecursive(root3))

	// Rebuild trees (modified by Morris)
	root4 := &TreeNode{Val: 1}
	root4.Right = &TreeNode{Val: 2}
	root4.Right.Left = &TreeNode{Val: 3}

	root5 := &TreeNode{Val: 1}
	root5.Right = &TreeNode{Val: 2}
	root5.Right.Left = &TreeNode{Val: 3}

	fmt.Println("=== Approach 2: Iterative (Stack) ===")
	fmt.Printf("tree=[1,null,2,3]  got=%v  expected [1 3 2]\n", inorderIterative(root4))

	fmt.Println("=== Approach 3: Morris Traversal ===")
	fmt.Printf("tree=[1,null,2,3]  got=%v  expected [1 3 2]\n", inorderMorris(root5))
}
