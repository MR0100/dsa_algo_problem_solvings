package main

import "fmt"

// TreeNode is the standard LeetCode binary-tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// ── Approach 1: Brute Force (Double DFS) ─────────────────────────────────────
//
// bruteForce solves Path Sum III by treating EVERY node as a possible start of
// a downward path, then counting downward paths from that start whose running
// sum hits targetSum.
//
// Intuition:
//
//	A valid path is any contiguous top-to-bottom chain summing to targetSum. If
//	we anchor the start at each node in turn, the problem shrinks to "count
//	downward paths from THIS node that sum to target" — a straight recursion
//	that keeps subtracting node values. Doing that from every anchor covers
//	all paths, but re-walks shared suffixes, hence the quadratic cost.
//
// Algorithm:
//  1. countFrom(node, remaining): if node is nil return 0; let rem =
//     remaining - node.Val; count = (rem == 0 ? 1 : 0); recurse into both
//     children with rem; sum them up.
//  2. bruteForce walks every node (outer DFS) and adds countFrom(node,
//     targetSum) for each.
//
// Time:  O(n^2) worst case (a path/skewed tree) — each of n anchors may walk an
//
//	O(n) suffix; O(n log n) for a balanced tree.
//
// Space: O(h) — recursion stack, h = tree height.
func bruteForce(root *TreeNode, targetSum int) int {
	if root == nil {
		return 0
	}
	// Paths starting exactly at root, plus all paths that start deeper (handled
	// by recursing the OUTER walk into each child with the full targetSum).
	return countFrom(root, targetSum) +
		bruteForce(root.Left, targetSum) +
		bruteForce(root.Right, targetSum)
}

// countFrom counts downward paths that START at node and sum to remaining.
func countFrom(node *TreeNode, remaining int) int {
	if node == nil {
		return 0
	}
	rem := remaining - node.Val // consume this node's value along the path
	count := 0
	if rem == 0 {
		count = 1 // a path ending here (root-of-this-call → node) hits the target
	}
	// A longer path may still reach the target further down; keep descending.
	// Note: we do NOT stop at rem == 0, because negative values later could also
	// form additional valid paths of greater length.
	count += countFrom(node.Left, rem)
	count += countFrom(node.Right, rem)
	return count
}

// ── Approach 2: Prefix Sum + Hash Map (Optimal) ──────────────────────────────
//
// prefixSumHashMap solves Path Sum III in one DFS by carrying the running sum
// from the root and asking, at each node, how many earlier ancestors closed a
// path of exactly targetSum — the classic subarray-sum-equals-k trick lifted
// onto a root-to-node path.
//
// Intuition:
//
//	On any root-to-node path the value is a prefix sum. A sub-path from ancestor
//	A (exclusive) down to the current node sums to target exactly when
//	curr - prefixAtA == target, i.e. prefixAtA == curr - target. So keep a
//	frequency map of prefix sums seen on the current path; at each node the
//	number of new valid paths ending here is freq[curr - target]. Insert curr
//	before recursing, and REMOVE it on the way back up so the map only ever
//	reflects the current root-to-node chain (not siblings).
//
// Algorithm:
//  1. Map freq{0: 1} (empty prefix, so a full path from the root counts).
//  2. DFS(node, curr): curr += node.Val; total += freq[curr - target];
//     freq[curr]++; recurse into children; then freq[curr]-- (backtrack).
//  3. Return the accumulated total.
//
// Time:  O(n) — each node visited once, O(1) map work per node.
// Space: O(n) — the prefix-sum map plus O(h) recursion.
func prefixSumHashMap(root *TreeNode, targetSum int) int {
	freq := map[int]int{0: 1} // prefix sum 0 seen once: the empty prefix at the root
	total := 0

	var dfs func(node *TreeNode, curr int)
	dfs = func(node *TreeNode, curr int) {
		if node == nil {
			return
		}
		curr += node.Val // running sum from root down to (and including) node
		// Any ancestor prefix equal to (curr - target) closes a path summing to
		// target ending at this node; count all such ancestors.
		total += freq[curr-targetSum]
		freq[curr]++ // register this node's prefix for its descendants
		dfs(node.Left, curr)
		dfs(node.Right, curr)
		freq[curr]-- // backtrack: leave the map holding only the current chain
	}

	dfs(root, 0)
	return total
}

// buildTree builds a binary tree from a level-order slice using -1<<62 as the
// null sentinel (real node values can be negative, so nil needs its own token).
func buildTree(vals []interface{}) *TreeNode {
	if len(vals) == 0 || vals[0] == nil {
		return nil
	}
	root := &TreeNode{Val: vals[0].(int)}
	queue := []*TreeNode{root} // BFS frontier of nodes still needing children
	i := 1                     // index into vals for the next child value
	for len(queue) > 0 && i < len(vals) {
		node := queue[0]
		queue = queue[1:]
		// Attach left child if present.
		if i < len(vals) && vals[i] != nil {
			node.Left = &TreeNode{Val: vals[i].(int)}
			queue = append(queue, node.Left)
		}
		i++
		// Attach right child if present.
		if i < len(vals) && vals[i] != nil {
			node.Right = &TreeNode{Val: vals[i].(int)}
			queue = append(queue, node.Right)
		}
		i++
	}
	return root
}

func main() {
	// Example 1: [10,5,-3,3,2,null,11,3,-2,null,1], targetSum = 8 → 3
	ex1 := buildTree([]interface{}{10, 5, -3, 3, 2, nil, 11, 3, -2, nil, 1})
	// Example 2: [5,4,8,11,null,13,4,7,2,null,null,5,1], targetSum = 22 → 3
	ex2 := buildTree([]interface{}{5, 4, 8, 11, nil, 13, 4, 7, 2, nil, nil, 5, 1})

	fmt.Println("=== Approach 1: Brute Force (Double DFS) ===")
	fmt.Println(bruteForce(ex1, 8))  // expected 3
	fmt.Println(bruteForce(ex2, 22)) // expected 3

	fmt.Println("=== Approach 2: Prefix Sum + Hash Map (Optimal) ===")
	fmt.Println(prefixSumHashMap(ex1, 8))  // expected 3
	fmt.Println(prefixSumHashMap(ex2, 22)) // expected 3
}
