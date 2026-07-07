package main

import "fmt"

// TreeNode is the standard LeetCode binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// ── Approach 1: Brute Force (Validate Every Subtree) ─────────────────────────
//
// bruteForce checks, for every node, whether its entire subtree is a BST and
// if so how many nodes it has, keeping the maximum.
//
// Intuition:
//
//	The direct reading of the problem: for each node, ask "is the subtree
//	rooted here a valid BST?" (bounds check) and "how big is it?" (count).
//	Take the largest size over all nodes whose subtree is a BST. Correct but
//	wasteful — it re-walks the same descendants many times.
//
// Algorithm:
//  1. For each node, run isBST(node, -inf, +inf) to test the whole subtree.
//  2. If it is a BST, its size (count of all nodes) is a candidate answer.
//  3. Recurse into children to consider their subtrees too; keep the max.
//
// Time:  O(n^2) — each of n nodes triggers an O(n) validate+count in the worst
//
//	case (a skewed tree).
//
// Space: O(h) recursion.
func bruteForce(root *TreeNode) int {
	if root == nil {
		return 0
	}
	// If this whole subtree is a BST, its node count is a candidate.
	if isBST(root, -1<<62, 1<<62) {
		return countNodes(root)
	}
	// Otherwise the best BST lies within one of the children's subtrees.
	l := bruteForce(root.Left)
	r := bruteForce(root.Right)
	if l > r {
		return l
	}
	return r
}

// isBST reports whether the subtree at node is a BST with all values strictly
// inside (lo, hi).
func isBST(node *TreeNode, lo, hi int) bool {
	if node == nil {
		return true // empty subtree is trivially a BST
	}
	if node.Val <= lo || node.Val >= hi {
		return false // value violates the inherited bound
	}
	// Left values must be < node.Val; right values must be > node.Val.
	return isBST(node.Left, lo, node.Val) && isBST(node.Right, node.Val, hi)
}

// countNodes returns the number of nodes in the subtree.
func countNodes(node *TreeNode) int {
	if node == nil {
		return 0
	}
	return 1 + countNodes(node.Left) + countNodes(node.Right)
}

// ── Approach 2: Post-Order Bottom-Up (Optimal, O(n)) ─────────────────────────
//
// postOrder computes, in one bottom-up pass, whether each subtree is a BST
// plus its size and value range, so a parent decides in O(1).
//
// Intuition:
//
//	A node forms a BST iff BOTH children are BSTs AND the node's value is
//	greater than the max of the left subtree and less than the min of the
//	right subtree. If every child returns (isBST, size, min, max), the parent
//	can validate and extend in constant time — no re-walking. Track the global
//	best size along the way.
//
// Algorithm:
//  1. DFS post-order returns info{isBST, size, min, max} for each node.
//  2. At a node: gather left/right info. If both BST and left.max < val < right.min,
//     this subtree is a BST of size left.size + right.size + 1; update the answer;
//     return the combined info. Otherwise return isBST=false (size still tracked
//     via the global best, so we just mark invalid).
//
// Time:  O(n) — each node visited once, O(1) work per node.
// Space: O(h) recursion.
func postOrder(root *TreeNode) int {
	best := 0
	type info struct {
		isBST    bool
		size     int
		min, max int
	}
	var dfs func(node *TreeNode) info
	dfs = func(node *TreeNode) info {
		if node == nil {
			// Empty subtree: a BST of size 0 with an "inverted" range so any
			// parent value satisfies max<val<min vacuously.
			return info{isBST: true, size: 0, min: 1 << 62, max: -1 << 62}
		}
		l := dfs(node.Left)  // info about left subtree
		r := dfs(node.Right) // info about right subtree
		// This node is a BST iff both sides are BSTs and val fits strictly between
		// the left subtree's max and the right subtree's min.
		if l.isBST && r.isBST && l.max < node.Val && node.Val < r.min {
			sz := l.size + r.size + 1 // combined size including this node
			if sz > best {
				best = sz // record a larger valid BST
			}
			return info{
				isBST: true,
				size:  sz,
				min:   min(node.Val, l.min), // smallest value in this BST
				max:   max(node.Val, r.max), // largest value in this BST
			}
		}
		// Not a BST: propagate invalid. Range fields are irrelevant now.
		return info{isBST: false}
	}
	dfs(root)
	return best
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	// Example 1: [10,5,15,1,8,null,7] → largest BST subtree has 3 nodes.
	//         10
	//        /  \
	//       5    15
	//      / \     \
	//     1   8     7
	ex1 := &TreeNode{Val: 10,
		Left:  &TreeNode{Val: 5, Left: &TreeNode{Val: 1}, Right: &TreeNode{Val: 8}},
		Right: &TreeNode{Val: 15, Right: &TreeNode{Val: 7}},
	}

	// Example 2: [4,2,7,2,3,5,null,2,null,null,null,null,null,1] → answer 2.
	// Decoded level-order (nulls are placeholders):
	//              4
	//            /   \
	//           2     7
	//          / \   /
	//         2   3 5
	//        /
	//       2
	//      /
	//     1
	// The largest BST subtree is {2,1} (node 2 with left child 1), size 2.
	ex2 := &TreeNode{Val: 4,
		Left: &TreeNode{Val: 2,
			Left: &TreeNode{Val: 2,
				Left: &TreeNode{Val: 2,
					Left: &TreeNode{Val: 1},
				},
			},
			Right: &TreeNode{Val: 3},
		},
		Right: &TreeNode{Val: 7,
			Left: &TreeNode{Val: 5},
		},
	}

	// A full BST for good measure.
	//        2
	//       / \
	//      1   3
	ex3 := &TreeNode{Val: 2, Left: &TreeNode{Val: 1}, Right: &TreeNode{Val: 3}}

	fmt.Println("=== Approach 1: Brute Force (Validate Every Subtree) ===")
	fmt.Println(bruteForce(ex1)) // expected 3
	fmt.Println(bruteForce(ex2)) // expected 2
	fmt.Println(bruteForce(ex3)) // expected 3
	fmt.Println(bruteForce(nil)) // expected 0

	fmt.Println("=== Approach 2: Post-Order Bottom-Up (Optimal) ===")
	fmt.Println(postOrder(ex1)) // expected 3
	fmt.Println(postOrder(ex2)) // expected 2
	fmt.Println(postOrder(ex3)) // expected 3
	fmt.Println(postOrder(nil)) // expected 0
}
