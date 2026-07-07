package main

import "fmt"

// TreeNode is a binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func inorderVals(root *TreeNode) []int {
	var res []int
	var dfs func(n *TreeNode)
	dfs = func(n *TreeNode) {
		if n == nil { return }
		dfs(n.Left); res = append(res, n.Val); dfs(n.Right)
	}
	dfs(root)
	return res
}

func preorderVals(root *TreeNode) []int {
	var res []int
	var dfs func(n *TreeNode)
	dfs = func(n *TreeNode) {
		if n == nil { return }
		res = append(res, n.Val); dfs(n.Left); dfs(n.Right)
	}
	dfs(root)
	return res
}

// ── Approach 1: Recursive with HashMap ───────────────────────────────────────
//
// buildTree solves Construct Binary Tree from Preorder and Inorder Traversal.
//
// Intuition:
//   preorder[0] is always the root.
//   Find root's position in inorder (using a hashmap for O(1) lookup).
//   Elements to the left of root in inorder = left subtree.
//   Elements to the right = right subtree.
//   The next `leftSize` elements in preorder = left subtree's preorder.
//   The remaining elements = right subtree's preorder.
//
// Algorithm:
//   build(preStart, preEnd, inStart, inEnd):
//     root = preorder[preStart]
//     rootIdx = inMap[root]
//     leftSize = rootIdx - inStart
//     root.Left = build(preStart+1, preStart+leftSize, inStart, rootIdx-1)
//     root.Right = build(preStart+leftSize+1, preEnd, rootIdx+1, inEnd)
//
// Time:  O(n) — n nodes, each looked up in O(1) via hashmap.
// Space: O(n) — hashmap + recursion stack.
func buildTree(preorder []int, inorder []int) *TreeNode {
	// build inorder index map for O(1) lookup
	inMap := make(map[int]int, len(inorder))
	for i, v := range inorder {
		inMap[v] = i
	}

	var build func(preStart, preEnd, inStart, inEnd int) *TreeNode
	build = func(preStart, preEnd, inStart, inEnd int) *TreeNode {
		if preStart > preEnd {
			return nil
		}
		rootVal := preorder[preStart]
		root := &TreeNode{Val: rootVal}
		rootIdx := inMap[rootVal]
		leftSize := rootIdx - inStart

		root.Left = build(preStart+1, preStart+leftSize, inStart, rootIdx-1)
		root.Right = build(preStart+leftSize+1, preEnd, rootIdx+1, inEnd)
		return root
	}

	return build(0, len(preorder)-1, 0, len(inorder)-1)
}

// ── Approach 2: Iterative Stack ───────────────────────────────────────────────
//
// buildTreeIterative solves Construct Binary Tree from Preorder and Inorder
// Traversal using an explicit stack.
//
// Intuition:
//   Walk preorder. The first element is the root. Maintain a stack.
//   Push each node. When preorder[i] matches inorder[inIdx], we've finished
//   the left subtree — keep popping until they don't match (right turn).
//   The next preorder element becomes the right child of the last popped node.
//
// Time:  O(n)
// Space: O(h)
func buildTreeIterative(preorder []int, inorder []int) *TreeNode {
	if len(preorder) == 0 {
		return nil
	}
	root := &TreeNode{Val: preorder[0]}
	stack := []*TreeNode{root}
	inIdx := 0

	for i := 1; i < len(preorder); i++ {
		node := &TreeNode{Val: preorder[i]}
		if stack[len(stack)-1].Val != inorder[inIdx] {
			// still going left
			stack[len(stack)-1].Left = node
		} else {
			// pop until mismatch — find where to attach as right child
			var parent *TreeNode
			for len(stack) > 0 && stack[len(stack)-1].Val == inorder[inIdx] {
				parent = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				inIdx++
			}
			parent.Right = node
		}
		stack = append(stack, node)
	}
	return root
}

func main() {
	fmt.Println("=== Approach 1: Recursive with HashMap ===")
	pre1 := []int{3, 9, 20, 15, 7}
	in1 := []int{9, 3, 15, 20, 7}
	t1 := buildTree(pre1, in1)
	fmt.Printf("preorder=%v inorder=%v\n  result inorder=%v  expected [9 3 15 20 7]\n", pre1, in1, inorderVals(t1))
	fmt.Printf("  result preorder=%v  expected [3 9 20 15 7]\n", preorderVals(t1))

	pre2 := []int{-1}
	in2 := []int{-1}
	t2 := buildTree(pre2, in2)
	fmt.Printf("preorder=%v inorder=%v  result=%v  expected [-1]\n", pre2, in2, inorderVals(t2))

	fmt.Println("=== Approach 2: Iterative Stack ===")
	pre3 := []int{3, 9, 20, 15, 7}
	in3 := []int{9, 3, 15, 20, 7}
	t3 := buildTreeIterative(pre3, in3)
	fmt.Printf("preorder=%v inorder=%v\n  result inorder=%v  expected [9 3 15 20 7]\n", pre3, in3, inorderVals(t3))
}
