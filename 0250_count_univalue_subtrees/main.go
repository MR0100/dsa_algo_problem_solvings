package main

import "fmt"

// TreeNode is a standard binary-tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// ── Approach 1: Post-Order DFS with Return Flag (Optimal) ─────────────────────
//
// countUnivalPostOrder counts the univalue subtrees of a binary tree in a
// single post-order traversal.
//
// Intuition:
//
//	A subtree is "univalue" if every node in it shares the same value. Whether a
//	node's subtree is univalue depends ENTIRELY on its children's subtrees:
//	- a leaf is always univalue;
//	- an internal node is univalue iff BOTH child subtrees are univalue AND each
//	  existing child's value equals the node's own value.
//	That is a bottom-up (post-order) fact, so we recurse to the leaves first,
//	then combine. We return a boolean "is this subtree univalue" up the call
//	stack and bump a shared counter each time the answer is yes.
//
// Algorithm:
//  1. dfs(node) returns whether node's subtree is univalue.
//  2. nil is treated as univalue (identity — a missing child never breaks it).
//  3. Recurse left and right FIRST (post-order) so both flags are known.
//  4. node is univalue iff left and right are both univalue AND every present
//     child has the same value as node.
//  5. If univalue, count++.
//  6. Return the flag.
//
// Time:  O(n) — each node visited once.
// Space: O(h) — recursion stack, h = tree height (O(n) worst, O(log n) balanced).
func countUnivalPostOrder(root *TreeNode) int {
	count := 0

	var dfs func(node *TreeNode) bool
	dfs = func(node *TreeNode) bool {
		if node == nil {
			return true // a missing child imposes no constraint
		}

		// Post-order: resolve BOTH children before deciding this node.
		leftUni := dfs(node.Left)
		rightUni := dfs(node.Right)

		// If either subtree isn't univalue, this one can't be either.
		if !leftUni || !rightUni {
			return false
		}
		// Every present child's value must match this node's value.
		if node.Left != nil && node.Left.Val != node.Val {
			return false
		}
		if node.Right != nil && node.Right.Val != node.Val {
			return false
		}

		count++     // this subtree is univalue — tally it
		return true // and report that fact to the parent
	}

	dfs(root)
	return count
}

// ── Approach 2: Post-Order Returning (isUnival, count) Pair ───────────────────
//
// countUnivalPair is the same post-order idea, but instead of mutating a shared
// counter via closure it threads the running count UP through the return value
// as a (isUnival, count) pair. Purely functional style — no captured variable.
//
// Intuition:
//
//	Each recursive call reports two things about the subtree rooted at node:
//	whether it is univalue, and how many univalue subtrees it contains in total
//	(including itself if applicable). The parent adds the two children's counts,
//	then decides its own univalue status and adds 1 if so.
//
// Algorithm:
//  1. dfs(node) returns (isUnival bool, count int).
//  2. nil → (true, 0).
//  3. Recurse both children; total = leftCount + rightCount.
//  4. Determine isUnival with the same rule as Approach 1.
//  5. If isUnival, total++.
//  6. Return (isUnival, total).
//
// Time:  O(n).
// Space: O(h) recursion.
func countUnivalPair(root *TreeNode) int {
	var dfs func(node *TreeNode) (bool, int)
	dfs = func(node *TreeNode) (bool, int) {
		if node == nil {
			return true, 0 // no node, no univalue subtrees below
		}

		leftUni, leftCount := dfs(node.Left)    // resolve left subtree
		rightUni, rightCount := dfs(node.Right) // resolve right subtree
		total := leftCount + rightCount         // univalue subtrees seen so far

		// This node is univalue only if both children are AND values agree.
		isUni := leftUni && rightUni
		if isUni && node.Left != nil && node.Left.Val != node.Val {
			isUni = false
		}
		if isUni && node.Right != nil && node.Right.Val != node.Val {
			isUni = false
		}

		if isUni {
			total++ // count this subtree itself
		}
		return isUni, total
	}

	_, count := dfs(root)
	return count
}

func main() {
	// Example 1:  root = [5,1,5,5,5,null,5]
	//        5
	//       / \
	//      1   5
	//     / \    \
	//    5   5    5
	// Univalue subtrees: the three leaf 5's, the right child 5 (with leaf 5),
	// and the two leaf 5's under node 1 → 4 total.
	ex1 := &TreeNode{
		Val:  5,
		Left: &TreeNode{Val: 1, Left: &TreeNode{Val: 5}, Right: &TreeNode{Val: 5}},
		Right: &TreeNode{Val: 5,
			Right: &TreeNode{Val: 5}},
	}

	// Example 2: root = [] → 0
	var ex2 *TreeNode = nil

	// Example 3: root = [5,5,5,5,5,null,5] → 6
	//        5
	//       / \
	//      5   5
	//     / \    \
	//    5   5    5
	ex3 := &TreeNode{
		Val:  5,
		Left: &TreeNode{Val: 5, Left: &TreeNode{Val: 5}, Right: &TreeNode{Val: 5}},
		Right: &TreeNode{Val: 5,
			Right: &TreeNode{Val: 5}},
	}

	fmt.Println("=== Approach 1: Post-Order DFS with Return Flag (Optimal) ===")
	fmt.Printf("ex1 [5,1,5,5,5,null,5]  got=%d  expected 4\n", countUnivalPostOrder(ex1)) // expected 4
	fmt.Printf("ex2 []                  got=%d  expected 0\n", countUnivalPostOrder(ex2)) // expected 0
	fmt.Printf("ex3 [5,5,5,5,5,null,5]  got=%d  expected 6\n", countUnivalPostOrder(ex3)) // expected 6

	fmt.Println("=== Approach 2: Post-Order Returning (isUnival, count) Pair ===")
	fmt.Printf("ex1 [5,1,5,5,5,null,5]  got=%d  expected 4\n", countUnivalPair(ex1)) // expected 4
	fmt.Printf("ex2 []                  got=%d  expected 0\n", countUnivalPair(ex2)) // expected 0
	fmt.Printf("ex3 [5,5,5,5,5,null,5]  got=%d  expected 6\n", countUnivalPair(ex3)) // expected 6
}
