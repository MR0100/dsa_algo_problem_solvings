package main

import "fmt"

// TreeNode is the standard LeetCode binary-tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// ── Approach 1: Naive Recursion (exponential) ────────────────────────────────
//
// naiveRecursion solves House Robber III by, at each node, choosing the better
// of "rob this node (skip children, take grandchildren)" vs "skip this node
// (free to rob children)".
//
// Intuition:
//
//	The thief cannot rob two directly-linked houses. So at the root either:
//	  rob root  = root.Val + rob(root.Left.children) + rob(root.Right.children)
//	  skip root = rob(root.Left) + rob(root.Right)
//	Take the max. Correct, but it recomputes the same subtrees repeatedly:
//	robbing the root recurses into grandchildren, and skipping it recurses into
//	children which again recurse into those grandchildren → exponential work.
//
// Algorithm:
//  1. If node is nil, return 0.
//  2. robThis = node.Val + (recurse on all four grandchildren).
//  3. skipThis = recurse(node.Left) + recurse(node.Right).
//  4. Return max(robThis, skipThis).
//
// Time:  O(2^h) — overlapping subproblems recomputed; roughly exponential.
// Space: O(h) recursion stack, h = tree height.
func naiveRecursion(root *TreeNode) int {
	if root == nil {
		return 0 // empty subtree yields nothing
	}
	robThis := root.Val // we rob this house...
	if root.Left != nil {
		// ...so its children are off-limits; jump to grandchildren.
		robThis += naiveRecursion(root.Left.Left) + naiveRecursion(root.Left.Right)
	}
	if root.Right != nil {
		robThis += naiveRecursion(root.Right.Left) + naiveRecursion(root.Right.Right)
	}
	// Or skip this house and rob its children freely.
	skipThis := naiveRecursion(root.Left) + naiveRecursion(root.Right)
	return max(robThis, skipThis) // best of the two choices
}

// ── Approach 2: Memoized Recursion (Top-Down DP) ─────────────────────────────
//
// memoized solves House Robber III with the same rob/skip recurrence but caches
// each node's best result to kill the repeated work.
//
// Intuition:
//
//	The naive version recomputes rob(subtree) for the same node many times.
//	Memoize by node pointer: the first time we compute a node's best, store it;
//	later requests return the cached value. Each node is solved once.
//
// Algorithm:
//  1. memo maps *TreeNode → best amount rootable from that subtree.
//  2. dfs(node): if cached, return it; else compute robThis/skipThis exactly as
//     in Approach 1, cache max, return.
//
// Time:  O(n) — each node computed once, O(1) work per node.
// Space: O(n) memo + O(h) recursion stack.
func memoized(root *TreeNode) int {
	memo := map[*TreeNode]int{} // node → best amount for its subtree
	var dfs func(node *TreeNode) int
	dfs = func(node *TreeNode) int {
		if node == nil {
			return 0
		}
		if v, ok := memo[node]; ok {
			return v // already solved this subtree
		}
		robThis := node.Val
		if node.Left != nil {
			robThis += dfs(node.Left.Left) + dfs(node.Left.Right)
		}
		if node.Right != nil {
			robThis += dfs(node.Right.Left) + dfs(node.Right.Right)
		}
		skipThis := dfs(node.Left) + dfs(node.Right)
		best := max(robThis, skipThis)
		memo[node] = best // cache before returning
		return best
	}
	return dfs(root)
}

// ── Approach 3: DFS Returning a Pair (Optimal) ───────────────────────────────
//
// robTreeDP solves House Robber III with a single post-order traversal in which
// each node returns two numbers: the best if this node is robbed, and the best
// if it is not.
//
// Intuition:
//
//	The memo keyed by node still walks the tree twice conceptually (children and
//	grandchildren). Instead, let each node hand its parent a pair:
//	  withNode    = node.Val + left.withoutNode + right.withoutNode
//	  withoutNode = max(left pair) + max(right pair)
//	The parent decides using only its children's pairs — no grandchild peeking,
//	one clean post-order pass.
//
// Algorithm:
//  1. dfs(node) returns (rob, notRob).
//  2. Base: nil → (0, 0).
//  3. Combine children: rob = node.Val + l.notRob + r.notRob;
//     notRob = max(l.rob, l.notRob) + max(r.rob, r.notRob).
//  4. Answer = max(root.rob, root.notRob).
//
// Time:  O(n) — one visit per node.
// Space: O(h) recursion stack only.
func robTreeDP(root *TreeNode) int {
	// dfs returns {rob = best including node, notRob = best excluding node}.
	var dfs func(node *TreeNode) (int, int)
	dfs = func(node *TreeNode) (int, int) {
		if node == nil {
			return 0, 0 // nothing to rob, nothing to skip
		}
		lRob, lNot := dfs(node.Left)  // children solved first (post-order)
		rRob, rNot := dfs(node.Right) //
		// If we rob this node, both children must be skipped.
		rob := node.Val + lNot + rNot
		// If we skip this node, each child independently takes its own best.
		notRob := max(lRob, lNot) + max(rRob, rNot)
		return rob, notRob
	}
	rob, notRob := dfs(root)
	return max(rob, notRob) // whole-tree best
}

// ── helpers ──────────────────────────────────────────────────────────────────

// buildTree builds a binary tree from a level-order slice using nil for missing
// nodes (LeetCode's array format).
func buildTree(vals []interface{}) *TreeNode {
	if len(vals) == 0 || vals[0] == nil {
		return nil
	}
	root := &TreeNode{Val: vals[0].(int)}
	queue := []*TreeNode{root} // nodes still awaiting children
	i := 1
	for len(queue) > 0 && i < len(vals) {
		node := queue[0]
		queue = queue[1:]
		if i < len(vals) && vals[i] != nil { // left child
			node.Left = &TreeNode{Val: vals[i].(int)}
			queue = append(queue, node.Left)
		}
		i++
		if i < len(vals) && vals[i] != nil { // right child
			node.Right = &TreeNode{Val: vals[i].(int)}
			queue = append(queue, node.Right)
		}
		i++
	}
	return root
}

func main() {
	// Example 1: [3,2,3,null,3,null,1] → 7
	ex1 := buildTree([]interface{}{3, 2, 3, nil, 3, nil, 1})
	// Example 2: [3,4,5,1,3,null,1] → 9
	ex2 := buildTree([]interface{}{3, 4, 5, 1, 3, nil, 1})

	fmt.Println("=== Approach 1: Naive Recursion ===")
	fmt.Println(naiveRecursion(ex1)) // 7
	fmt.Println(naiveRecursion(ex2)) // 9

	fmt.Println("=== Approach 2: Memoized Recursion (Top-Down DP) ===")
	fmt.Println(memoized(ex1)) // 7
	fmt.Println(memoized(ex2)) // 9

	fmt.Println("=== Approach 3: DFS Returning a Pair (Optimal) ===")
	fmt.Println(robTreeDP(ex1)) // 7
	fmt.Println(robTreeDP(ex2)) // 9
}
