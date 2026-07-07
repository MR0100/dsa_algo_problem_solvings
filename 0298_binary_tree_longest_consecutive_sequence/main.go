package main

import "fmt"

// Binary Tree Longest Consecutive Sequence
//
// Given the root of a binary tree, return the length of the longest consecutive
// sequence path. The path must go top-down (parent -> child) but need not pass
// through the root; consecutive means each next value is exactly one greater
// than its parent (increasing by 1).

// TreeNode is the standard binary-tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// ── Approach 1: Top-Down DFS (pass the running length down) ───────────────────
//
// topDownDFS carries, into each recursive call, the length of the consecutive
// run ending at the current node.
//
// Intuition:
//
//	A consecutive path grows only when a child's value is exactly parent+1. So
//	pass down "how long the increasing run is so far, ending at me". If a child
//	continues the run, its length is parent's length + 1; otherwise it resets to
//	1. Track the global maximum as we go.
//
// Algorithm:
//  1. dfs(node, parentVal, lengthSoFar):
//     - if node is nil, return.
//     - length = (node.Val == parentVal+1) ? lengthSoFar+1 : 1.
//     - update global best with length.
//     - recurse into left and right with (node.Val, length).
//  2. Start with dfs(root, root.Val-1, 0) so the root itself begins a run of 1.
//
// Time:  O(n) — each node visited once.
// Space: O(h) — recursion stack, h = tree height.
func topDownDFS(root *TreeNode) int {
	best := 0
	var dfs func(node *TreeNode, parentVal, lengthSoFar int)
	dfs = func(node *TreeNode, parentVal, lengthSoFar int) {
		if node == nil {
			return
		}
		length := 1 // a lone node is always a run of length 1
		if node.Val == parentVal+1 {
			length = lengthSoFar + 1 // continues the parent's increasing run
		}
		if length > best {
			best = length // record the longest run seen so far
		}
		dfs(node.Left, node.Val, length)  // extend downward through left child
		dfs(node.Right, node.Val, length) // and through right child
	}
	if root != nil {
		// seed parentVal so root.Val == parentVal+1 is false initially
		dfs(root, root.Val-1, 0)
	}
	return best
}

// ── Approach 2: Bottom-Up DFS (return length upward) ─────────────────────────
//
// bottomUpDFS computes, for each node, the longest consecutive run STARTING at
// that node and returns it to the parent.
//
// Intuition:
//
//	Solve children first. The run starting at a node is 1 plus the child's run
//	IF that child is exactly node+1. Combine left and right, update the global
//	best, and hand the node's own run length back up.
//
// Algorithm:
//  1. dfs(node) returns the length of the longest increasing run starting here.
//  2. len = 1. If left exists and left.Val == node.Val+1, len = max(len, 1+dfs(left)).
//     If right exists and right.Val == node.Val+1, len = max(len, 1+dfs(right)).
//  3. Update global best with len; return len.
//
// Time:  O(n) — each node visited once.
// Space: O(h) — recursion stack.
func bottomUpDFS(root *TreeNode) int {
	best := 0
	var dfs func(node *TreeNode) int
	dfs = func(node *TreeNode) int {
		if node == nil {
			return 0
		}
		leftLen := dfs(node.Left)   // longest run starting at left child
		rightLen := dfs(node.Right) // longest run starting at right child

		length := 1 // this node alone
		if node.Left != nil && node.Left.Val == node.Val+1 {
			length = max(length, 1+leftLen) // extend through left if consecutive
		}
		if node.Right != nil && node.Right.Val == node.Val+1 {
			length = max(length, 1+rightLen) // extend through right if consecutive
		}
		if length > best {
			best = length // global maximum across all starting points
		}
		return length
	}
	dfs(root)
	return best
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	// Example 1: root = [1,null,3,2,4,null,null,null,5]
	//        1
	//         \
	//          3
	//         / \
	//        2   4
	//             \
	//              5
	// Longest consecutive path: 3-4-5, length 3.
	root1 := &TreeNode{
		Val: 1,
		Right: &TreeNode{
			Val:  3,
			Left: &TreeNode{Val: 2},
			Right: &TreeNode{
				Val:   4,
				Right: &TreeNode{Val: 5},
			},
		},
	}
	// Example 2: root = [2,null,3,2,null,1]
	//        2
	//         \
	//          3
	//         /
	//        2
	//       /
	//      1
	// Longest consecutive path: 2-3, length 2 (2->1 is decreasing).
	root2 := &TreeNode{
		Val: 2,
		Right: &TreeNode{
			Val: 3,
			Left: &TreeNode{
				Val:  2,
				Left: &TreeNode{Val: 1},
			},
		},
	}

	fmt.Println("=== Approach 1: Top-Down DFS ===")
	fmt.Println(topDownDFS(root1)) // expected 3
	fmt.Println(topDownDFS(root2)) // expected 2

	fmt.Println("=== Approach 2: Bottom-Up DFS ===")
	fmt.Println(bottomUpDFS(root1)) // expected 3
	fmt.Println(bottomUpDFS(root2)) // expected 2
}
