package main

import "fmt"

// TreeNode is a binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// ── Approach 1: BFS with Direction Flag ──────────────────────────────────────
//
// zigzagLevelOrder solves Binary Tree Zigzag Level Order Traversal using BFS.
//
// Intuition:
//   Same as level-order BFS (#102), but alternate the order of values
//   collected per level. Track a boolean `leftToRight`: true = collect
//   normally, false = reverse the level's collected values before appending.
//
// Time:  O(n)
// Space: O(w)
func zigzagLevelOrder(root *TreeNode) [][]int {
	if root == nil {
		return nil
	}
	var result [][]int
	queue := []*TreeNode{root}
	leftToRight := true

	for len(queue) > 0 {
		levelSize := len(queue)
		level := make([]int, levelSize)
		for i := 0; i < levelSize; i++ {
			node := queue[0]
			queue = queue[1:]
			// fill position based on direction
			if leftToRight {
				level[i] = node.Val
			} else {
				level[levelSize-1-i] = node.Val // fill from right
			}
			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
		result = append(result, level)
		leftToRight = !leftToRight
	}
	return result
}

// ── Approach 2: DFS with Level Tracking ──────────────────────────────────────
//
// zigzagLevelOrderDFS solves Binary Tree Zigzag Level Order Traversal using DFS.
//
// Intuition:
//   DFS with depth parameter. For odd depths (right-to-left), prepend the
//   value to the level slice instead of appending.
//
// Time:  O(n)
// Space: O(h)
func zigzagLevelOrderDFS(root *TreeNode) [][]int {
	var result [][]int
	var dfs func(node *TreeNode, depth int)
	dfs = func(node *TreeNode, depth int) {
		if node == nil {
			return
		}
		if depth == len(result) {
			result = append(result, []int{})
		}
		if depth%2 == 0 {
			// left-to-right: append
			result[depth] = append(result[depth], node.Val)
		} else {
			// right-to-left: prepend
			result[depth] = append([]int{node.Val}, result[depth]...)
		}
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

	fmt.Println("=== Approach 1: BFS with Direction Flag ===")
	fmt.Printf("tree=[3,9,20,null,null,15,7]  got=%v  expected [[3] [20 9] [15 7]]\n", zigzagLevelOrder(t1))
	fmt.Printf("tree=[1]  got=%v  expected [[1]]\n", zigzagLevelOrder(t2))

	t3 := &TreeNode{Val: 3,
		Left:  &TreeNode{Val: 9},
		Right: &TreeNode{Val: 20, Left: &TreeNode{Val: 15}, Right: &TreeNode{Val: 7}},
	}

	fmt.Println("=== Approach 2: DFS with Level Tracking ===")
	fmt.Printf("tree=[3,9,20,null,null,15,7]  got=%v  expected [[3] [20 9] [15 7]]\n", zigzagLevelOrderDFS(t3))
	fmt.Printf("tree=[1]  got=%v  expected [[1]]\n", zigzagLevelOrderDFS(t2))
}
