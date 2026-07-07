package main

import "fmt"

// TreeNode is a binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// ── Approach 1: Recursive DFS ────────────────────────────────────────────────
//
// preorderRecursive solves Binary Tree Preorder Traversal with plain recursion.
//
// Intuition:
//
//	Preorder = root, left, right. The definition is itself recursive, so the
//	code writes itself: visit the node, then recurse into each subtree.
//
// Algorithm:
//  1. If node == nil → return.
//  2. Append node.Val (visit root first — that's what "pre" means).
//  3. Recurse left, then recurse right.
//
// Time:  O(n) — every node is visited exactly once.
// Space: O(h) — recursion stack of the tree height h (O(n) worst, O(log n) balanced).
func preorderRecursive(root *TreeNode) []int {
	result := []int{} // non-nil so an empty tree prints as [] not nil
	var dfs func(node *TreeNode)
	dfs = func(node *TreeNode) {
		if node == nil {
			return // empty subtree contributes nothing
		}
		result = append(result, node.Val) // ROOT first
		dfs(node.Left)                    // then the whole LEFT subtree
		dfs(node.Right)                   // then the whole RIGHT subtree
	}
	dfs(root)
	return result
}

// ── Approach 2: Iterative (Explicit Stack) ───────────────────────────────────
//
// preorderIterative solves Binary Tree Preorder Traversal by replacing the
// call stack with an explicit slice-based stack.
//
// Intuition:
//
//	A stack pops the most recently pushed node. If we visit a node the moment
//	we pop it, we just need its children visited left-before-right — so push
//	RIGHT first, LEFT second (LIFO order flips them back).
//
// Algorithm:
//  1. If root == nil → return [].
//  2. stack = [root].
//  3. While stack non-empty: pop node, append node.Val,
//     push node.Right (if any), then push node.Left (if any).
//
// Time:  O(n) — each node pushed and popped exactly once.
// Space: O(h) — stack holds at most one full root-to-leaf frontier
// (O(n) worst-case skewed/bushy, O(log n) balanced).
func preorderIterative(root *TreeNode) []int {
	result := []int{}
	if root == nil {
		return result // nothing to traverse
	}
	stack := []*TreeNode{root}
	for len(stack) > 0 {
		node := stack[len(stack)-1]       // peek top
		stack = stack[:len(stack)-1]      // pop it
		result = append(result, node.Val) // visit on pop → root before children
		if node.Right != nil {
			stack = append(stack, node.Right) // pushed FIRST → popped LAST
		}
		if node.Left != nil {
			stack = append(stack, node.Left) // pushed LAST → popped NEXT (left first)
		}
	}
	return result
}

// ── Approach 3: Morris Traversal (O(1) Space, Optimal) ───────────────────────
//
// preorderMorris solves Binary Tree Preorder Traversal in O(1) extra space by
// temporarily threading each left subtree's rightmost node back to the root.
//
// Intuition:
//
//	The stack exists only so we can return to a node after finishing its left
//	subtree. Morris threading stores that return address inside the tree: the
//	inorder predecessor's nil Right pointer temporarily points back at us.
//	For PREorder we emit the node on FIRST arrival (before descending left),
//	unlike inorder Morris which emits on the second arrival.
//
// Algorithm:
//  1. curr = root.
//  2. If curr.Left == nil → visit curr, go right.
//  3. Else find pred = rightmost node of curr.Left (stopping at threads):
//     a. pred.Right == nil  → visit curr (first arrival), thread
//     pred.Right = curr, descend left.
//     b. pred.Right == curr → second arrival: remove thread (pred.Right=nil),
//     go right (left subtree already emitted).
//
// Time:  O(n) — each edge is walked at most twice (once to thread, once to unthread).
// Space: O(1) — no stack, no recursion; the tree is restored before returning.
func preorderMorris(root *TreeNode) []int {
	result := []int{}
	curr := root
	for curr != nil {
		if curr.Left == nil {
			result = append(result, curr.Val) // leaf-ward: visit and move on
			curr = curr.Right
			continue
		}
		// Find the inorder predecessor: rightmost node in the left subtree,
		// stopping early if we hit a thread pointing back to curr.
		pred := curr.Left
		for pred.Right != nil && pred.Right != curr {
			pred = pred.Right
		}
		if pred.Right == nil {
			result = append(result, curr.Val) // FIRST arrival → preorder visit
			pred.Right = curr                 // lay the return thread
			curr = curr.Left                  // descend into the left subtree
		} else {
			pred.Right = nil  // second arrival → remove thread (restore tree)
			curr = curr.Right // left subtree done; continue rightward
		}
	}
	return result
}

// ── Helpers ───────────────────────────────────────────────────────────────────

// null marks a missing child in the level-order encodings below.
const null = -1 << 31

// buildTree constructs a binary tree from LeetCode's level-order array form.
func buildTree(levelOrder []int) *TreeNode {
	if len(levelOrder) == 0 || levelOrder[0] == null {
		return nil
	}
	root := &TreeNode{Val: levelOrder[0]}
	queue := []*TreeNode{root} // parents awaiting children
	i := 1
	for len(queue) > 0 && i < len(levelOrder) {
		parent := queue[0]
		queue = queue[1:]
		if i < len(levelOrder) && levelOrder[i] != null {
			parent.Left = &TreeNode{Val: levelOrder[i]}
			queue = append(queue, parent.Left)
		}
		i++ // consume the left-child slot even when it is null
		if i < len(levelOrder) && levelOrder[i] != null {
			parent.Right = &TreeNode{Val: levelOrder[i]}
			queue = append(queue, parent.Right)
		}
		i++ // consume the right-child slot even when it is null
	}
	return root
}

func main() {
	// Official LeetCode examples: (level-order input, expected preorder).
	examples := []struct {
		tree   []int
		expect []int
	}{
		{[]int{1, null, 2, 3}, []int{1, 2, 3}},                                                 // Example 1
		{[]int{1, 2, 3, 4, 5, null, 8, null, null, 6, 7, 9}, []int{1, 2, 4, 5, 6, 7, 3, 8, 9}}, // Example 2
		{[]int{}, []int{}},   // Example 3
		{[]int{1}, []int{1}}, // Example 4
	}

	approaches := []struct {
		name string
		fn   func(*TreeNode) []int
	}{
		{"Approach 1: Recursive DFS", preorderRecursive},
		{"Approach 2: Iterative (Stack)", preorderIterative},
		{"Approach 3: Morris Traversal (Optimal Space)", preorderMorris},
	}

	for _, ap := range approaches {
		fmt.Printf("=== %s ===\n", ap.name)
		for i, ex := range examples {
			root := buildTree(ex.tree) // fresh tree per run (Morris mutates temporarily)
			got := ap.fn(root)
			fmt.Printf("Example %d: → %v (expected %v)\n", i+1, got, ex.expect)
			// expected: [1 2 3], [1 2 4 5 6 7 3 8 9], [], [1]
		}
		fmt.Println()
	}
}
