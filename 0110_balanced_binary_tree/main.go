package main

import "fmt"

// TreeNode is a binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// ── Approach 1: Top-Down Recursive (Naive) ───────────────────────────────────
//
// isBalanced solves Balanced Binary Tree top-down.
//
// Intuition:
//   A tree is balanced iff at every node: |height(left) - height(right)| <= 1,
//   and both subtrees are themselves balanced.
//   Compute height separately for each node — O(n log n) due to repeated work.
//
// Time:  O(n log n) — height called at each node, each taking O(n).
// Space: O(h)
func isBalanced(root *TreeNode) bool {
	var height func(node *TreeNode) int
	height = func(node *TreeNode) int {
		if node == nil {
			return 0
		}
		l := height(node.Left)
		r := height(node.Right)
		if l > r {
			return l + 1
		}
		return r + 1
	}

	if root == nil {
		return true
	}
	l := height(root.Left)
	r := height(root.Right)
	diff := l - r
	if diff < 0 {
		diff = -diff
	}
	return diff <= 1 && isBalanced(root.Left) && isBalanced(root.Right)
}

// ── Approach 2: Bottom-Up Recursive (Optimal) ────────────────────────────────
//
// isBalancedOptimal solves Balanced Binary Tree bottom-up in a single pass.
//
// Intuition:
//   Compute height in a post-order traversal. Return -1 (sentinel) as height
//   if any subtree is already unbalanced — this propagates failure upward
//   without redundant recomputation.
//
// Algorithm:
//   check(node):
//     if nil: return 0
//     l = check(left);  if l==-1 return -1
//     r = check(right); if r==-1 return -1
//     if |l-r| > 1: return -1
//     return 1 + max(l, r)
//
// Time:  O(n) — each node visited once.
// Space: O(h)
func isBalancedOptimal(root *TreeNode) bool {
	var check func(node *TreeNode) int
	check = func(node *TreeNode) int {
		if node == nil {
			return 0
		}
		l := check(node.Left)
		if l == -1 {
			return -1 // left subtree already unbalanced
		}
		r := check(node.Right)
		if r == -1 {
			return -1 // right subtree already unbalanced
		}
		diff := l - r
		if diff < 0 {
			diff = -diff
		}
		if diff > 1 {
			return -1 // current node violates balance
		}
		if l > r {
			return l + 1
		}
		return r + 1
	}
	return check(root) != -1
}

func main() {
	// [3,9,20,null,null,15,7] — balanced
	t1 := &TreeNode{Val: 3,
		Left:  &TreeNode{Val: 9},
		Right: &TreeNode{Val: 20, Left: &TreeNode{Val: 15}, Right: &TreeNode{Val: 7}},
	}
	// [1,2,2,3,3,null,null,4,4] — not balanced
	t2 := &TreeNode{Val: 1,
		Left: &TreeNode{Val: 2,
			Left: &TreeNode{Val: 3,
				Left: &TreeNode{Val: 4}, Right: &TreeNode{Val: 4}},
			Right: &TreeNode{Val: 3}},
		Right: &TreeNode{Val: 2},
	}

	fmt.Println("=== Approach 1: Top-Down (O(n log n)) ===")
	fmt.Printf("tree=[3,9,20,null,null,15,7]  got=%v  expected true\n", isBalanced(t1))
	fmt.Printf("tree=[1,2,2,3,3,null,null,4,4]  got=%v  expected false\n", isBalanced(t2))
	fmt.Printf("tree=nil  got=%v  expected true\n", isBalanced(nil))

	t3 := &TreeNode{Val: 3,
		Left:  &TreeNode{Val: 9},
		Right: &TreeNode{Val: 20, Left: &TreeNode{Val: 15}, Right: &TreeNode{Val: 7}},
	}
	t4 := &TreeNode{Val: 1,
		Left: &TreeNode{Val: 2,
			Left: &TreeNode{Val: 3,
				Left: &TreeNode{Val: 4}, Right: &TreeNode{Val: 4}},
			Right: &TreeNode{Val: 3}},
		Right: &TreeNode{Val: 2},
	}

	fmt.Println("=== Approach 2: Bottom-Up Optimal (O(n)) ===")
	fmt.Printf("tree=[3,9,20,null,null,15,7]  got=%v  expected true\n", isBalancedOptimal(t3))
	fmt.Printf("tree=[1,2,2,3,3,null,null,4,4]  got=%v  expected false\n", isBalancedOptimal(t4))
	fmt.Printf("tree=nil  got=%v  expected true\n", isBalancedOptimal(nil))
}
