package main

import "fmt"

// TreeNode is a binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// ── Approach 1: BFS + Reverse ─────────────────────────────────────────────────
//
// levelOrderBottom solves Binary Tree Level Order Traversal II using BFS then
// reversing the result.
//
// Intuition:
//   Do standard level-order BFS (#102), collect all levels, then reverse the
//   outer slice so leaves appear first.
//
// Time:  O(n)
// Space: O(w)
func levelOrderBottom(root *TreeNode) [][]int {
	if root == nil {
		return nil
	}
	var result [][]int
	queue := []*TreeNode{root}

	for len(queue) > 0 {
		levelSize := len(queue)
		level := make([]int, 0, levelSize)
		for i := 0; i < levelSize; i++ {
			node := queue[0]
			queue = queue[1:]
			level = append(level, node.Val)
			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
		result = append(result, level)
	}

	// reverse the levels
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return result
}

// ── Approach 2: DFS + Prepend ──────────────────────────────────────────────────
//
// levelOrderBottomDFS solves Binary Tree Level Order Traversal II using DFS.
//
// Intuition:
//   DFS with depth tracking. Instead of appending each new level, prepend it
//   so the root level ends up at the back (deepest leaves at front).
//   Since prepend is expensive, we append as in #102 and reverse at end —
//   or simply build in pre-order and reverse. Same complexity.
//
// Time:  O(n)
// Space: O(h)
func levelOrderBottomDFS(root *TreeNode) [][]int {
	var result [][]int
	var dfs func(node *TreeNode, depth int)
	dfs = func(node *TreeNode, depth int) {
		if node == nil {
			return
		}
		if depth == len(result) {
			result = append(result, []int{})
		}
		result[depth] = append(result[depth], node.Val)
		dfs(node.Left, depth+1)
		dfs(node.Right, depth+1)
	}
	dfs(root, 0)

	// reverse so leaves are first
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return result
}

func main() {
	// [3,9,20,null,null,15,7]
	t1 := &TreeNode{Val: 3,
		Left:  &TreeNode{Val: 9},
		Right: &TreeNode{Val: 20, Left: &TreeNode{Val: 15}, Right: &TreeNode{Val: 7}},
	}
	t2 := &TreeNode{Val: 1}
	var t3 *TreeNode

	fmt.Println("=== Approach 1: BFS + Reverse ===")
	fmt.Printf("tree=[3,9,20,null,null,15,7]  got=%v  expected [[15 7] [9 20] [3]]\n", levelOrderBottom(t1))
	fmt.Printf("tree=[1]  got=%v  expected [[1]]\n", levelOrderBottom(t2))
	fmt.Printf("tree=[]  got=%v  expected []\n", levelOrderBottom(t3))

	t4 := &TreeNode{Val: 3,
		Left:  &TreeNode{Val: 9},
		Right: &TreeNode{Val: 20, Left: &TreeNode{Val: 15}, Right: &TreeNode{Val: 7}},
	}

	fmt.Println("=== Approach 2: DFS + Reverse ===")
	fmt.Printf("tree=[3,9,20,null,null,15,7]  got=%v  expected [[15 7] [9 20] [3]]\n", levelOrderBottomDFS(t4))
	fmt.Printf("tree=[1]  got=%v  expected [[1]]\n", levelOrderBottomDFS(t2))
}
