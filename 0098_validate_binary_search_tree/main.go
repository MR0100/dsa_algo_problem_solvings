package main

import (
	"fmt"
	"math"
)

// TreeNode is a binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// ── Approach 1: Inorder Traversal ────────────────────────────────────────────
//
// isValidBST solves Validate Binary Search Tree by checking that the inorder
// traversal is strictly increasing.
//
// Intuition:
//   A valid BST's inorder traversal yields values in strictly ascending order.
//   Track the previously seen value; if current <= prev, invalid.
//
// Time:  O(n)
// Space: O(h) — call stack.
func isValidBSTInorder(root *TreeNode) bool {
	prev := math.MinInt64
	valid := true
	var inorder func(node *TreeNode)
	inorder = func(node *TreeNode) {
		if node == nil || !valid {
			return
		}
		inorder(node.Left)
		if node.Val <= prev {
			valid = false
			return
		}
		prev = node.Val
		inorder(node.Right)
	}
	inorder(root)
	return valid
}

// ── Approach 2: Min/Max Range Check (Recursive) ──────────────────────────────
//
// isValidBST solves Validate Binary Search Tree by passing valid (min, max)
// ranges for each node.
//
// Intuition:
//   For each node, its value must be in the range (min, max).
//   Root: (-∞, +∞).
//   Left child: (min, parent.Val).
//   Right child: (parent.Val, max).
//
// Time:  O(n)
// Space: O(h)
func isValidBST(root *TreeNode) bool {
	var check func(node *TreeNode, min, max int) bool
	check = func(node *TreeNode, min, max int) bool {
		if node == nil {
			return true
		}
		if node.Val <= min || node.Val >= max {
			return false
		}
		return check(node.Left, min, node.Val) && check(node.Right, node.Val, max)
	}
	return check(root, math.MinInt64, math.MaxInt64)
}

// ── Approach 3: Iterative Inorder ────────────────────────────────────────────
//
// isValidBSTIterative solves Validate Binary Search Tree using an iterative
// inorder traversal with a stack.
//
// Time:  O(n)
// Space: O(h)
func isValidBSTIterative(root *TreeNode) bool {
	stack := []*TreeNode{}
	prev := math.MinInt64
	curr := root

	for curr != nil || len(stack) > 0 {
		for curr != nil {
			stack = append(stack, curr)
			curr = curr.Left
		}
		curr = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if curr.Val <= prev {
			return false
		}
		prev = curr.Val
		curr = curr.Right
	}
	return true
}

func main() {
	// tree [2,1,3]: valid BST
	t1 := &TreeNode{Val: 2, Left: &TreeNode{Val: 1}, Right: &TreeNode{Val: 3}}
	// tree [5,1,4,null,null,3,6]: invalid (4 < 5)
	t2 := &TreeNode{Val: 5,
		Left:  &TreeNode{Val: 1},
		Right: &TreeNode{Val: 4, Left: &TreeNode{Val: 3}, Right: &TreeNode{Val: 6}},
	}
	// tree [2,2,2]: invalid (duplicates)
	t3 := &TreeNode{Val: 2, Left: &TreeNode{Val: 2}, Right: &TreeNode{Val: 2}}

	fmt.Println("=== Approach 1: Inorder Traversal ===")
	fmt.Printf("tree=[2,1,3]  got=%v  expected true\n", isValidBSTInorder(t1))
	fmt.Printf("tree=[5,1,4,null,null,3,6]  got=%v  expected false\n", isValidBSTInorder(t2))

	fmt.Println("=== Approach 2: Min/Max Range Check ===")
	fmt.Printf("tree=[2,1,3]  got=%v  expected true\n", isValidBST(t1))
	fmt.Printf("tree=[5,1,4,null,null,3,6]  got=%v  expected false\n", isValidBST(t2))
	fmt.Printf("tree=[2,2,2]  got=%v  expected false\n", isValidBST(t3))

	fmt.Println("=== Approach 3: Iterative Inorder ===")
	fmt.Printf("tree=[2,1,3]  got=%v  expected true\n", isValidBSTIterative(t1))
	fmt.Printf("tree=[5,1,4,null,null,3,6]  got=%v  expected false\n", isValidBSTIterative(t2))
}
