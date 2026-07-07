package main

import "fmt"

// TreeNode is the standard LeetCode binary-tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// Given the root of a binary tree, return the sum of all LEFT leaves. A leaf
// has no children; a *left* leaf is a leaf that is the LEFT child of its parent.
// The key subtlety: "left leaf" is decided by the parent-child link, not by a
// node's own value — so we must know whether we arrived via a left edge.

// ── Approach 1: Recursive DFS with an "isLeft" flag ──────────────────────────
//
// recursiveDFS walks the tree, carrying whether the current node is the left
// child of its parent; when a node is a leaf AND arrived via a left edge, its
// value is added.
//
// Intuition:
//
//	A node cannot tell on its own if it is a "left leaf" — that depends on the
//	edge the parent used to reach it. So pass a boolean `isLeft` down. At each
//	node: if it is a leaf and isLeft, contribute its value; otherwise recurse
//	into the left child with isLeft=true and the right child with isLeft=false.
//
// Algorithm:
//  1. helper(node, isLeft): if node is nil return 0.
//  2. If node is a leaf (no children) and isLeft, return node.Val.
//  3. Else return helper(node.Left, true) + helper(node.Right, false).
//  4. Answer = helper(root, false) — the root itself is never a "left leaf".
//
// Time:  O(n) — visits each node once.
// Space: O(h) — recursion stack, h = tree height.
func recursiveDFS(root *TreeNode) int {
	var helper func(node *TreeNode, isLeft bool) int
	helper = func(node *TreeNode, isLeft bool) int {
		if node == nil {
			return 0 // empty subtree adds nothing
		}
		// A leaf reached through a LEFT edge is exactly a "left leaf".
		if node.Left == nil && node.Right == nil {
			if isLeft {
				return node.Val
			}
			return 0 // a leaf, but it was a right child — ignore
		}
		// Recurse: left child is a left edge, right child is not.
		return helper(node.Left, true) + helper(node.Right, false)
	}
	// The root has no parent, so it is treated as "not a left child".
	return helper(root, false)
}

// ── Approach 2: Parent-Aware Recursion (check the left child directly) ────────
//
// parentCheckDFS inspects, at each node, whether its LEFT child is a leaf; if
// so it adds that child's value, then recurses into both children.
//
// Intuition:
//
//	Instead of passing a flag down, look one level down. Standing at a parent,
//	its left child is a "left leaf" precisely when that left child exists and
//	has no children of its own. Add it right here, then continue the walk.
//
// Algorithm:
//  1. dfs(node): if node is nil return 0. sum = 0.
//  2. If node.Left != nil and node.Left is a leaf: sum += node.Left.Val.
//     Else sum += dfs(node.Left) (left child is an internal subtree to explore).
//  3. sum += dfs(node.Right).
//  4. Return sum; answer = dfs(root).
//
// Time:  O(n) — each node visited once.
// Space: O(h) — recursion stack.
func parentCheckDFS(root *TreeNode) int {
	var dfs func(node *TreeNode) int
	dfs = func(node *TreeNode) int {
		if node == nil {
			return 0
		}
		sum := 0
		if node.Left != nil {
			// Is the left child itself a leaf? Then it's a left leaf — count it.
			if node.Left.Left == nil && node.Left.Right == nil {
				sum += node.Left.Val
			} else {
				sum += dfs(node.Left) // otherwise descend into it
			}
		}
		sum += dfs(node.Right) // right subtree may contain its own left leaves
		return sum
	}
	return dfs(root)
}

// ── Approach 3: Iterative DFS with an Explicit Stack (Optimal / no recursion) ─
//
// iterativeStack does the same traversal without recursion, pushing (node,
// isLeft) pairs onto a manual stack.
//
// Intuition:
//
//	Any recursive DFS can be made iterative with an explicit stack. Store each
//	node together with the "did I arrive via a left edge?" flag, exactly as the
//	recursive version carried it. Pop, and when a popped node is a leaf reached
//	via a left edge, add its value.
//
// Algorithm:
//  1. If root is nil return 0. Push (root, false).
//  2. While the stack is non-empty: pop (node, isLeft).
//     - If node is a leaf and isLeft: add node.Val.
//     - Else push (node.Left, true) and (node.Right, false) if non-nil.
//  3. Return the accumulated sum.
//
// Time:  O(n) — every node is pushed and popped once.
// Space: O(h) — the stack holds at most one root-to-leaf path's worth of nodes.
func iterativeStack(root *TreeNode) int {
	if root == nil {
		return 0
	}
	type frame struct {
		node   *TreeNode
		isLeft bool // whether this node is the left child of its parent
	}
	stack := []frame{{root, false}} // root has no parent -> not a left child
	sum := 0
	for len(stack) > 0 {
		// Pop the top frame.
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		node, isLeft := top.node, top.isLeft

		if node.Left == nil && node.Right == nil {
			if isLeft {
				sum += node.Val // leaf reached via a left edge = left leaf
			}
			continue // leaves have no children to push
		}
		// Push children with the correct edge label.
		if node.Left != nil {
			stack = append(stack, frame{node.Left, true})
		}
		if node.Right != nil {
			stack = append(stack, frame{node.Right, false})
		}
	}
	return sum
}

func main() {
	// Example 1: root = [3,9,20,null,null,15,7]
	//        3
	//       / \
	//      9  20
	//         / \
	//        15  7
	// Left leaves: 9 (left child, leaf) and 15 (left child, leaf) -> 24.
	ex1 := &TreeNode{
		Val:  3,
		Left: &TreeNode{Val: 9},
		Right: &TreeNode{
			Val:   20,
			Left:  &TreeNode{Val: 15},
			Right: &TreeNode{Val: 7},
		},
	}

	// Example 2: root = [1] -> the single node is the root (not a left leaf) -> 0.
	ex2 := &TreeNode{Val: 1}

	fmt.Println("=== Approach 1: Recursive DFS (isLeft flag) ===")
	fmt.Printf("[3,9,20,null,null,15,7] -> %d\n", recursiveDFS(ex1)) // expected 24
	fmt.Printf("[1]                     -> %d\n", recursiveDFS(ex2)) // expected 0

	fmt.Println("=== Approach 2: Parent-Check DFS ===")
	fmt.Printf("[3,9,20,null,null,15,7] -> %d\n", parentCheckDFS(ex1)) // expected 24
	fmt.Printf("[1]                     -> %d\n", parentCheckDFS(ex2)) // expected 0

	fmt.Println("=== Approach 3: Iterative DFS (explicit stack) ===")
	fmt.Printf("[3,9,20,null,null,15,7] -> %d\n", iterativeStack(ex1)) // expected 24
	fmt.Printf("[1]                     -> %d\n", iterativeStack(ex2)) // expected 0
}
