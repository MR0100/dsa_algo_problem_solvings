package main

import "fmt"

// TreeNode is a binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// ── Approach 1: Post-Order DFS ────────────────────────────────────────────────
//
// maxPathSum solves Binary Tree Maximum Path Sum using DFS.
//
// Intuition:
//   A path can go: left-subtree → root → right-subtree (or any subset).
//   For each node, compute the max "gain" we can extend upward: max(0, gain(left)) + val + max(0, gain(right)).
//   The local max path through this node = val + max(0,left) + max(0,right).
//   Update a global max, but return only val + max(one side, 0) to the parent
//   (can only continue along one path upward).
//
// Algorithm:
//   gain(node):
//     if nil: return 0
//     leftGain  = max(0, gain(left))
//     rightGain = max(0, gain(right))
//     maxPath = max(maxPath, node.Val + leftGain + rightGain)
//     return node.Val + max(leftGain, rightGain)
//
// Time:  O(n)
// Space: O(h)
func maxPathSum(root *TreeNode) int {
	maxPath := -1 << 31

	var gain func(node *TreeNode) int
	gain = func(node *TreeNode) int {
		if node == nil {
			return 0
		}
		leftGain := gain(node.Left)
		if leftGain < 0 {
			leftGain = 0 // discard negative contribution
		}
		rightGain := gain(node.Right)
		if rightGain < 0 {
			rightGain = 0
		}
		// update global max with path through this node
		pathThroughNode := node.Val + leftGain + rightGain
		if pathThroughNode > maxPath {
			maxPath = pathThroughNode
		}
		// return max extension upward (only one side)
		if leftGain > rightGain {
			return node.Val + leftGain
		}
		return node.Val + rightGain
	}

	gain(root)
	return maxPath
}

func main() {
	// [1,2,3] → max path 2→1→3 = 6
	t1 := &TreeNode{Val: 1, Left: &TreeNode{Val: 2}, Right: &TreeNode{Val: 3}}
	// [-10,9,20,null,null,15,7] → max path 15→20→7 = 42
	t2 := &TreeNode{Val: -10,
		Left: &TreeNode{Val: 9},
		Right: &TreeNode{Val: 20, Left: &TreeNode{Val: 15}, Right: &TreeNode{Val: 7}},
	}

	fmt.Println("=== Approach 1: Post-Order DFS ===")
	fmt.Printf("tree=[1,2,3]  got=%d  expected 6\n", maxPathSum(t1))
	fmt.Printf("tree=[-10,9,20,null,null,15,7]  got=%d  expected 42\n", maxPathSum(t2))
}
