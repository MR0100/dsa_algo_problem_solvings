package main

import "fmt"

// TreeNode is a binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// ── Approach 1: Backtracking DFS ─────────────────────────────────────────────
//
// pathSum solves Path Sum II by collecting all root-to-leaf paths with the
// given sum using backtracking.
//
// Intuition:
//   DFS down each path, maintaining the current path and remaining sum.
//   At a leaf, if remaining == 0, record a copy of the current path.
//   Backtrack by removing the last element after each recursive call.
//
// Time:  O(n^2) — at most n leaves, each path copy is O(n).
// Space: O(h) — recursion stack + path length.
func pathSum(root *TreeNode, targetSum int) [][]int {
	var result [][]int
	var path []int

	var dfs func(node *TreeNode, remaining int)
	dfs = func(node *TreeNode, remaining int) {
		if node == nil {
			return
		}
		path = append(path, node.Val)
		remaining -= node.Val

		if node.Left == nil && node.Right == nil && remaining == 0 {
			// copy path (don't append slice directly — it shares backing array)
			cp := make([]int, len(path))
			copy(cp, path)
			result = append(result, cp)
		}
		dfs(node.Left, remaining)
		dfs(node.Right, remaining)

		path = path[:len(path)-1] // backtrack
	}

	dfs(root, targetSum)
	return result
}

func main() {
	// [5,4,8,11,null,13,4,7,2,null,null,5,1], targetSum=22
	t1 := &TreeNode{Val: 5,
		Left: &TreeNode{Val: 4,
			Left: &TreeNode{Val: 11,
				Left: &TreeNode{Val: 7}, Right: &TreeNode{Val: 2}},
		},
		Right: &TreeNode{Val: 8,
			Left: &TreeNode{Val: 13},
			Right: &TreeNode{Val: 4,
				Left: &TreeNode{Val: 5}, Right: &TreeNode{Val: 1}},
		},
	}
	// [1,2,3], targetSum=5
	t2 := &TreeNode{Val: 1, Left: &TreeNode{Val: 2}, Right: &TreeNode{Val: 3}}

	fmt.Println("=== Approach 1: Backtracking DFS ===")
	fmt.Printf("tree=[5,4,8,...] targetSum=22  got=%v  expected [[5 4 11 2] [5 8 4 5]]\n", pathSum(t1, 22))
	fmt.Printf("tree=[1,2,3] targetSum=5  got=%v  expected []\n", pathSum(t2, 5))
}
