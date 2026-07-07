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
// maxDepth solves Maximum Depth of Binary Tree recursively.
//
// Intuition:
//   The depth of a tree is 1 + max(depth(left), depth(right)).
//   Base case: nil → depth 0.
//
// Time:  O(n)
// Space: O(h)
func maxDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}
	leftDepth := maxDepth(root.Left)
	rightDepth := maxDepth(root.Right)
	if leftDepth > rightDepth {
		return leftDepth + 1
	}
	return rightDepth + 1
}

// ── Approach 2: Iterative BFS ─────────────────────────────────────────────────
//
// maxDepthBFS solves Maximum Depth of Binary Tree using level-order BFS.
//
// Intuition:
//   Count the number of levels. Each completed level increments the depth.
//
// Time:  O(n)
// Space: O(w)
func maxDepthBFS(root *TreeNode) int {
	if root == nil {
		return 0
	}
	depth := 0
	queue := []*TreeNode{root}

	for len(queue) > 0 {
		levelSize := len(queue)
		depth++
		for i := 0; i < levelSize; i++ {
			node := queue[0]
			queue = queue[1:]
			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
	}
	return depth
}

// ── Approach 3: Iterative DFS (Stack) ────────────────────────────────────────
//
// maxDepthDFS solves Maximum Depth of Binary Tree using an explicit DFS stack
// storing (node, depth) pairs.
//
// Time:  O(n)
// Space: O(h)
func maxDepthDFS(root *TreeNode) int {
	if root == nil {
		return 0
	}
	type item struct {
		node  *TreeNode
		depth int
	}
	stack := []item{{root, 1}}
	maxD := 0

	for len(stack) > 0 {
		curr := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if curr.depth > maxD {
			maxD = curr.depth
		}
		if curr.node.Left != nil {
			stack = append(stack, item{curr.node.Left, curr.depth + 1})
		}
		if curr.node.Right != nil {
			stack = append(stack, item{curr.node.Right, curr.depth + 1})
		}
	}
	return maxD
}

func main() {
	// [3,9,20,null,null,15,7] — depth 3
	t1 := &TreeNode{Val: 3,
		Left:  &TreeNode{Val: 9},
		Right: &TreeNode{Val: 20, Left: &TreeNode{Val: 15}, Right: &TreeNode{Val: 7}},
	}
	// [1,null,2] — depth 2
	t2 := &TreeNode{Val: 1, Right: &TreeNode{Val: 2}}

	fmt.Println("=== Approach 1: Recursive DFS ===")
	fmt.Printf("tree=[3,9,20,null,null,15,7]  got=%d  expected 3\n", maxDepth(t1))
	fmt.Printf("tree=[1,null,2]  got=%d  expected 2\n", maxDepth(t2))

	fmt.Println("=== Approach 2: BFS ===")
	t3 := &TreeNode{Val: 3,
		Left:  &TreeNode{Val: 9},
		Right: &TreeNode{Val: 20, Left: &TreeNode{Val: 15}, Right: &TreeNode{Val: 7}},
	}
	fmt.Printf("tree=[3,9,20,null,null,15,7]  got=%d  expected 3\n", maxDepthBFS(t3))
	fmt.Printf("tree=[1,null,2]  got=%d  expected 2\n", maxDepthBFS(t2))

	fmt.Println("=== Approach 3: Iterative DFS ===")
	t4 := &TreeNode{Val: 3,
		Left:  &TreeNode{Val: 9},
		Right: &TreeNode{Val: 20, Left: &TreeNode{Val: 15}, Right: &TreeNode{Val: 7}},
	}
	fmt.Printf("tree=[3,9,20,null,null,15,7]  got=%d  expected 3\n", maxDepthDFS(t4))
}
