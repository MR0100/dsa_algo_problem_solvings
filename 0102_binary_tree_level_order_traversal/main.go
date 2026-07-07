package main

import "fmt"

// TreeNode is a binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// ── Approach 1: BFS with Queue ────────────────────────────────────────────────
//
// levelOrder solves Binary Tree Level Order Traversal using BFS.
//
// Intuition:
//   Process nodes level by level. Use a queue. At the start of each level,
//   the queue size tells us how many nodes are in that level. Dequeue all
//   of them, record values, and enqueue their children.
//
// Time:  O(n)
// Space: O(w) — max width of tree (queue holds at most one level).
func levelOrder(root *TreeNode) [][]int {
	if root == nil {
		return nil
	}
	var result [][]int
	queue := []*TreeNode{root}

	for len(queue) > 0 {
		levelSize := len(queue) // number of nodes in this level
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
	return result
}

// ── Approach 2: DFS with Level Tracking ──────────────────────────────────────
//
// levelOrderDFS solves Binary Tree Level Order Traversal using DFS, passing
// the current depth to determine which level slice to append to.
//
// Intuition:
//   DFS pre-order traversal with depth parameter. When depth == len(result),
//   start a new level slice. Append node value to result[depth].
//
// Time:  O(n)
// Space: O(h) — recursion depth.
func levelOrderDFS(root *TreeNode) [][]int {
	var result [][]int
	var dfs func(node *TreeNode, depth int)
	dfs = func(node *TreeNode, depth int) {
		if node == nil {
			return
		}
		if depth == len(result) {
			result = append(result, []int{}) // start a new level
		}
		result[depth] = append(result[depth], node.Val)
		dfs(node.Left, depth+1)
		dfs(node.Right, depth+1)
	}
	dfs(root, 0)
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

	fmt.Println("=== Approach 1: BFS ===")
	fmt.Printf("tree=[3,9,20,null,null,15,7]  got=%v  expected [[3] [9 20] [15 7]]\n", levelOrder(t1))
	fmt.Printf("tree=[1]  got=%v  expected [[1]]\n", levelOrder(t2))
	fmt.Printf("tree=[]  got=%v  expected []\n", levelOrder(t3))

	// Rebuild
	t4 := &TreeNode{Val: 3,
		Left:  &TreeNode{Val: 9},
		Right: &TreeNode{Val: 20, Left: &TreeNode{Val: 15}, Right: &TreeNode{Val: 7}},
	}

	fmt.Println("=== Approach 2: DFS ===")
	fmt.Printf("tree=[3,9,20,null,null,15,7]  got=%v  expected [[3] [9 20] [15 7]]\n", levelOrderDFS(t4))
	fmt.Printf("tree=[1]  got=%v  expected [[1]]\n", levelOrderDFS(t2))
}
