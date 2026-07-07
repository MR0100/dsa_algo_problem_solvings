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
// hasPathSum solves Path Sum by subtracting node values as we descend.
//
// Intuition:
//   At each node, subtract its value from targetSum. When we reach a leaf
//   and the remaining target is 0, the path exists.
//
// Time:  O(n)
// Space: O(h)
func hasPathSum(root *TreeNode, targetSum int) bool {
	if root == nil {
		return false
	}
	remaining := targetSum - root.Val
	if root.Left == nil && root.Right == nil {
		return remaining == 0 // leaf: check if target exactly consumed
	}
	return hasPathSum(root.Left, remaining) || hasPathSum(root.Right, remaining)
}

// ── Approach 2: Iterative DFS (Stack) ────────────────────────────────────────
//
// hasPathSumIterative solves Path Sum using an explicit stack of (node, remaining).
//
// Intuition:
//   Same logic as recursive but with a (node, remaining) pair stack.
//   When we pop a leaf and remaining==0, return true.
//
// Time:  O(n)
// Space: O(h)
func hasPathSumIterative(root *TreeNode, targetSum int) bool {
	if root == nil {
		return false
	}
	type item struct {
		node      *TreeNode
		remaining int
	}
	stack := []item{{root, targetSum}}

	for len(stack) > 0 {
		curr := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		rem := curr.remaining - curr.node.Val
		if curr.node.Left == nil && curr.node.Right == nil && rem == 0 {
			return true
		}
		if curr.node.Left != nil {
			stack = append(stack, item{curr.node.Left, rem})
		}
		if curr.node.Right != nil {
			stack = append(stack, item{curr.node.Right, rem})
		}
	}
	return false
}

func main() {
	// [5,4,8,11,null,13,4,7,2,null,null,null,1], targetSum=22
	t1 := &TreeNode{Val: 5,
		Left: &TreeNode{Val: 4,
			Left: &TreeNode{Val: 11,
				Left: &TreeNode{Val: 7}, Right: &TreeNode{Val: 2}},
		},
		Right: &TreeNode{Val: 8,
			Left:  &TreeNode{Val: 13},
			Right: &TreeNode{Val: 4, Right: &TreeNode{Val: 1}},
		},
	}
	// [1,2,3], targetSum=5
	t2 := &TreeNode{Val: 1, Left: &TreeNode{Val: 2}, Right: &TreeNode{Val: 3}}

	fmt.Println("=== Approach 1: Recursive DFS ===")
	fmt.Printf("tree=[5,4,8,...] targetSum=22  got=%v  expected true\n", hasPathSum(t1, 22))
	fmt.Printf("tree=[1,2,3] targetSum=5  got=%v  expected false\n", hasPathSum(t2, 5))
	fmt.Printf("tree=nil targetSum=0  got=%v  expected false\n", hasPathSum(nil, 0))

	t3 := &TreeNode{Val: 5,
		Left: &TreeNode{Val: 4,
			Left: &TreeNode{Val: 11,
				Left: &TreeNode{Val: 7}, Right: &TreeNode{Val: 2}},
		},
		Right: &TreeNode{Val: 8,
			Left:  &TreeNode{Val: 13},
			Right: &TreeNode{Val: 4, Right: &TreeNode{Val: 1}},
		},
	}

	fmt.Println("=== Approach 2: Iterative DFS ===")
	fmt.Printf("tree=[5,4,8,...] targetSum=22  got=%v  expected true\n", hasPathSumIterative(t3, 22))
	fmt.Printf("tree=[1,2,3] targetSum=5  got=%v  expected false\n", hasPathSumIterative(t2, 5))
}
