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

// ── Approach 1: Recursive with HashMap ───────────────────────────────────────
//
// buildTree solves Construct Binary Tree from Inorder and Postorder Traversal.
//
// Intuition:
//   postorder[last] is always the root of the current subtree.
//   Find root's index in inorder (O(1) via hashmap).
//   Left of root in inorder = left subtree; right = right subtree.
//   rightSize = inEnd - rootIdx determines how many postorder elements go right.
//   Right subtree postorder ends at postEnd-1; left ends before that.
//
// Algorithm:
//   build(postStart, postEnd, inStart, inEnd):
//     root = postorder[postEnd]
//     rootIdx = inMap[root]
//     rightSize = inEnd - rootIdx
//     root.Right = build(postEnd-rightSize, postEnd-1, rootIdx+1, inEnd)
//     root.Left  = build(postStart, postEnd-rightSize-1, inStart, rootIdx-1)
//
// Time:  O(n)
// Space: O(n) — hashmap + recursion stack.
func buildTree(inorder []int, postorder []int) *TreeNode {
	inMap := make(map[int]int, len(inorder))
	for i, v := range inorder {
		inMap[v] = i
	}

	var build func(postStart, postEnd, inStart, inEnd int) *TreeNode
	build = func(postStart, postEnd, inStart, inEnd int) *TreeNode {
		if postStart > postEnd {
			return nil
		}
		rootVal := postorder[postEnd]
		root := &TreeNode{Val: rootVal}
		rootIdx := inMap[rootVal]
		rightSize := inEnd - rootIdx

		// right subtree uses the last `rightSize` elements of postorder (before current root)
		root.Right = build(postEnd-rightSize, postEnd-1, rootIdx+1, inEnd)
		root.Left = build(postStart, postEnd-rightSize-1, inStart, rootIdx-1)
		return root
	}

	return build(0, len(postorder)-1, 0, len(inorder)-1)
}

// ── Approach 2: Iterative (Reverse Postorder + Stack) ────────────────────────
//
// buildTreeIterative solves Construct Binary Tree from Inorder and Postorder
// Traversal iteratively using a stack.
//
// Intuition:
//   Reversed postorder is: root → right → left (mirror of preorder).
//   We can apply the same iterative trick as #105 but mirrored:
//   walk reversed postorder, use a reversed inorder pointer.
//   When a match is found, the next element becomes the left child.
//
// Time:  O(n)
// Space: O(h)
func buildTreeIterative(inorder []int, postorder []int) *TreeNode {
	n := len(postorder)
	if n == 0 {
		return nil
	}
	root := &TreeNode{Val: postorder[n-1]}
	stack := []*TreeNode{root}
	inIdx := n - 1 // walk inorder from right

	for i := n - 2; i >= 0; i-- {
		node := &TreeNode{Val: postorder[i]}
		if stack[len(stack)-1].Val != inorder[inIdx] {
			// still going right (mirror of going left in #105)
			stack[len(stack)-1].Right = node
		} else {
			var parent *TreeNode
			for len(stack) > 0 && stack[len(stack)-1].Val == inorder[inIdx] {
				parent = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				inIdx--
			}
			parent.Left = node
		}
		stack = append(stack, node)
	}
	return root
}

func main() {
	fmt.Println("=== Approach 1: Recursive with HashMap ===")
	in1 := []int{9, 3, 15, 20, 7}
	post1 := []int{9, 15, 7, 20, 3}
	t1 := buildTree(in1, post1)
	fmt.Printf("inorder=%v postorder=%v\n  result inorder=%v  expected [9 3 15 20 7]\n", in1, post1, inorderVals(t1))

	in2 := []int{-1}
	post2 := []int{-1}
	t2 := buildTree(in2, post2)
	fmt.Printf("inorder=%v postorder=%v  result=%v  expected [-1]\n", in2, post2, inorderVals(t2))

	fmt.Println("=== Approach 2: Iterative Stack ===")
	in3 := []int{9, 3, 15, 20, 7}
	post3 := []int{9, 15, 7, 20, 3}
	t3 := buildTreeIterative(in3, post3)
	fmt.Printf("inorder=%v postorder=%v\n  result inorder=%v  expected [9 3 15 20 7]\n", in3, post3, inorderVals(t3))
}
