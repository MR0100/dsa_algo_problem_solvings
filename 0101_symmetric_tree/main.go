package main

import "fmt"

// TreeNode is a binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// ── Approach 1: Recursive ─────────────────────────────────────────────────────
//
// isSymmetric solves Symmetric Tree recursively by checking if the left and
// right subtrees are mirrors of each other.
//
// Intuition:
//   A tree is symmetric iff its left subtree is a mirror of its right subtree.
//   Two trees are mirrors iff:
//   - Their roots have the same value.
//   - The left's right subtree is a mirror of the right's left subtree.
//   - The left's left subtree is a mirror of the right's right subtree.
//
// Time:  O(n) — each node visited once.
// Space: O(h) — recursion stack.
func isSymmetric(root *TreeNode) bool {
	var mirror func(left, right *TreeNode) bool
	mirror = func(left, right *TreeNode) bool {
		if left == nil && right == nil {
			return true
		}
		if left == nil || right == nil {
			return false
		}
		return left.Val == right.Val &&
			mirror(left.Left, right.Right) &&
			mirror(left.Right, right.Left)
	}
	if root == nil {
		return true
	}
	return mirror(root.Left, root.Right)
}

// ── Approach 2: Iterative BFS ─────────────────────────────────────────────────
//
// isSymmetricIterative solves Symmetric Tree using a queue to compare
// mirror pairs level by level.
//
// Intuition:
//   At each step, compare pairs of nodes that should be mirrors.
//   Start with (root.Left, root.Right). For each pair (l, r):
//   - Both nil: ok.
//   - One nil or different values: not symmetric.
//   - Otherwise: enqueue (l.Left, r.Right) and (l.Right, r.Left).
//
// Time:  O(n)
// Space: O(w) — width of tree.
func isSymmetricIterative(root *TreeNode) bool {
	if root == nil {
		return true
	}
	type pair struct{ l, r *TreeNode }
	queue := []pair{{root.Left, root.Right}}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr.l == nil && curr.r == nil {
			continue
		}
		if curr.l == nil || curr.r == nil || curr.l.Val != curr.r.Val {
			return false
		}
		queue = append(queue, pair{curr.l.Left, curr.r.Right})
		queue = append(queue, pair{curr.l.Right, curr.r.Left})
	}
	return true
}

func main() {
	// [1,2,2,3,4,4,3] — symmetric
	t1 := &TreeNode{Val: 1,
		Left:  &TreeNode{Val: 2, Left: &TreeNode{Val: 3}, Right: &TreeNode{Val: 4}},
		Right: &TreeNode{Val: 2, Left: &TreeNode{Val: 4}, Right: &TreeNode{Val: 3}},
	}
	// [1,2,2,null,3,null,3] — not symmetric
	t2 := &TreeNode{Val: 1,
		Left:  &TreeNode{Val: 2, Right: &TreeNode{Val: 3}},
		Right: &TreeNode{Val: 2, Right: &TreeNode{Val: 3}},
	}

	fmt.Println("=== Approach 1: Recursive ===")
	fmt.Printf("tree=[1,2,2,3,4,4,3]  got=%v  expected true\n", isSymmetric(t1))
	fmt.Printf("tree=[1,2,2,null,3,null,3]  got=%v  expected false\n", isSymmetric(t2))

	// Rebuild for iterative
	t3 := &TreeNode{Val: 1,
		Left:  &TreeNode{Val: 2, Left: &TreeNode{Val: 3}, Right: &TreeNode{Val: 4}},
		Right: &TreeNode{Val: 2, Left: &TreeNode{Val: 4}, Right: &TreeNode{Val: 3}},
	}
	t4 := &TreeNode{Val: 1,
		Left:  &TreeNode{Val: 2, Right: &TreeNode{Val: 3}},
		Right: &TreeNode{Val: 2, Right: &TreeNode{Val: 3}},
	}

	fmt.Println("=== Approach 2: Iterative BFS ===")
	fmt.Printf("tree=[1,2,2,3,4,4,3]  got=%v  expected true\n", isSymmetricIterative(t3))
	fmt.Printf("tree=[1,2,2,null,3,null,3]  got=%v  expected false\n", isSymmetricIterative(t4))
}
