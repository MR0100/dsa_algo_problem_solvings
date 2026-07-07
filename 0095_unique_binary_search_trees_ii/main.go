package main

import "fmt"

// TreeNode is a binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// printTree returns a level-order representation for display.
func inorder(root *TreeNode) []int {
	if root == nil {
		return nil
	}
	var result []int
	var dfs func(n *TreeNode)
	dfs = func(n *TreeNode) {
		if n == nil {
			return
		}
		dfs(n.Left)
		result = append(result, n.Val)
		dfs(n.Right)
	}
	dfs(root)
	return result
}

// ── Approach 1: Recursion (Divide and Conquer) ───────────────────────────────
//
// generateTrees solves Unique Binary Search Trees II by trying each value
// as the root and recursively building all left and right subtrees.
//
// Intuition:
//   For a range [start..end], each value `i` can be the root.
//   Left subtree: all BSTs with values [start..i-1].
//   Right subtree: all BSTs with values [i+1..end].
//   Combine: for each (leftTree, rightTree) pair, create a new root node.
//
// Time:  O(Catalan(n) × n) — Catalan(n) trees, each requiring O(n) to build.
// Space: O(Catalan(n) × n) — output.
func generateTrees(n int) []*TreeNode {
	if n == 0 {
		return nil
	}
	var generate func(start, end int) []*TreeNode
	generate = func(start, end int) []*TreeNode {
		if start > end {
			return []*TreeNode{nil} // nil represents an empty tree
		}
		var allTrees []*TreeNode
		for i := start; i <= end; i++ {
			leftTrees := generate(start, i-1)
			rightTrees := generate(i+1, end)
			for _, left := range leftTrees {
				for _, right := range rightTrees {
					root := &TreeNode{Val: i, Left: left, Right: right}
					allTrees = append(allTrees, root)
				}
			}
		}
		return allTrees
	}
	return generate(1, n)
}

// ── Approach 2: Memoized Recursion ───────────────────────────────────────────
//
// generateTreesMemo solves Unique Binary Search Trees II with memoization.
//
// Intuition:
//   The structure of the trees for range [1..k] is independent of the actual
//   values (they are relative). Cache by (start, end) pairs.
//
//   However, since node values differ for the same (start,end) in different
//   contexts, a true structural cache requires value-offset adjustment.
//   A simpler memo by (start, end) works directly since values are unique.
//
// Time:  O(Catalan(n) × n) — same asymptotic; memoization helps avoid
//         redundant subtree reconstruction in certain usage patterns.
// Space: O(Catalan(n) × n + n² memo entries)
func generateTreesMemo(n int) []*TreeNode {
	if n == 0 {
		return nil
	}
	memo := make(map[[2]int][]*TreeNode)
	var generate func(start, end int) []*TreeNode
	generate = func(start, end int) []*TreeNode {
		if start > end {
			return []*TreeNode{nil}
		}
		key := [2]int{start, end}
		if trees, ok := memo[key]; ok {
			return trees
		}
		var allTrees []*TreeNode
		for i := start; i <= end; i++ {
			leftTrees := generate(start, i-1)
			rightTrees := generate(i+1, end)
			for _, left := range leftTrees {
				for _, right := range rightTrees {
					allTrees = append(allTrees, &TreeNode{Val: i, Left: left, Right: right})
				}
			}
		}
		memo[key] = allTrees
		return allTrees
	}
	return generate(1, n)
}

func main() {
	fmt.Println("=== Approach 1: Recursion (Divide and Conquer) ===")
	trees1 := generateTrees(3)
	fmt.Printf("n=3  count=%d  expected 5\n", len(trees1))
	for _, t := range trees1 {
		fmt.Println("  inorder:", inorder(t))
	}

	trees2 := generateTrees(1)
	fmt.Printf("n=1  count=%d  expected 1\n", len(trees2))

	fmt.Println("=== Approach 2: Memoized Recursion ===")
	trees3 := generateTreesMemo(3)
	fmt.Printf("n=3  count=%d  expected 5\n", len(trees3))

	trees4 := generateTreesMemo(1)
	fmt.Printf("n=1  count=%d  expected 1\n", len(trees4))
}
