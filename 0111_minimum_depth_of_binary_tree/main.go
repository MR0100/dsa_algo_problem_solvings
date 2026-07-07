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
// minDepth solves Minimum Depth of Binary Tree recursively.
//
// Intuition:
//   Minimum depth = shortest path from root to a leaf node.
//   A leaf has no children. If only one child exists, we MUST go down that
//   side (the nil side has no leaf). So: if left is nil, recurse only right
//   and vice versa. If both non-nil, return 1 + min(left, right).
//
// Time:  O(n)
// Space: O(h)
func minDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}
	if root.Left == nil && root.Right == nil {
		return 1 // leaf
	}
	if root.Left == nil {
		return 1 + minDepth(root.Right)
	}
	if root.Right == nil {
		return 1 + minDepth(root.Left)
	}
	// both children exist
	l := minDepth(root.Left)
	r := minDepth(root.Right)
	if l < r {
		return l + 1
	}
	return r + 1
}

// ── Approach 2: BFS (Optimal for Wide Trees) ──────────────────────────────────
//
// minDepthBFS solves Minimum Depth of Binary Tree using BFS.
//
// Intuition:
//   BFS level-by-level guarantees the first leaf we encounter is at the
//   minimum depth. Return immediately on finding the first leaf.
//   For wide trees this is faster than DFS in practice.
//
// Time:  O(n) worst case, but stops early at first leaf.
// Space: O(w)
func minDepthBFS(root *TreeNode) int {
	if root == nil {
		return 0
	}
	queue := []*TreeNode{root}
	depth := 0

	for len(queue) > 0 {
		levelSize := len(queue)
		depth++
		for i := 0; i < levelSize; i++ {
			node := queue[0]
			queue = queue[1:]
			if node.Left == nil && node.Right == nil {
				return depth // first leaf found
			}
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

func main() {
	// [3,9,20,null,null,15,7] — min depth 2 (path 3→9)
	t1 := &TreeNode{Val: 3,
		Left:  &TreeNode{Val: 9},
		Right: &TreeNode{Val: 20, Left: &TreeNode{Val: 15}, Right: &TreeNode{Val: 7}},
	}
	// [2,null,3,null,4,null,5,null,6] — min depth 5 (only right spine)
	t2 := &TreeNode{Val: 2,
		Right: &TreeNode{Val: 3,
			Right: &TreeNode{Val: 4,
				Right: &TreeNode{Val: 5,
					Right: &TreeNode{Val: 6},
				},
			},
		},
	}

	fmt.Println("=== Approach 1: Recursive DFS ===")
	fmt.Printf("tree=[3,9,20,null,null,15,7]  got=%d  expected 2\n", minDepth(t1))
	fmt.Printf("tree=[2,null,3,null,4,null,5,null,6]  got=%d  expected 5\n", minDepth(t2))

	t3 := &TreeNode{Val: 3,
		Left:  &TreeNode{Val: 9},
		Right: &TreeNode{Val: 20, Left: &TreeNode{Val: 15}, Right: &TreeNode{Val: 7}},
	}
	t4 := &TreeNode{Val: 2,
		Right: &TreeNode{Val: 3,
			Right: &TreeNode{Val: 4,
				Right: &TreeNode{Val: 5,
					Right: &TreeNode{Val: 6},
				},
			},
		},
	}

	fmt.Println("=== Approach 2: BFS ===")
	fmt.Printf("tree=[3,9,20,null,null,15,7]  got=%d  expected 2\n", minDepthBFS(t3))
	fmt.Printf("tree=[2,null,3,null,4,null,5,null,6]  got=%d  expected 5\n", minDepthBFS(t4))
}
