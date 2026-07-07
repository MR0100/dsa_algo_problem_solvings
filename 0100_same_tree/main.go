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
// isSameTree solves Same Tree recursively.
//
// Intuition:
//   Two trees are the same iff their roots have the same value and their left
//   and right subtrees are also the same.
//
//   Base cases:
//   - Both nil: same (return true).
//   - One nil: not same (return false).
//   - Values differ: not same (return false).
//
// Time:  O(n) — n = min(|p|, |q|) nodes compared.
// Space: O(h) — recursion stack.
func isSameTree(p *TreeNode, q *TreeNode) bool {
	if p == nil && q == nil {
		return true
	}
	if p == nil || q == nil {
		return false
	}
	if p.Val != q.Val {
		return false
	}
	return isSameTree(p.Left, q.Left) && isSameTree(p.Right, q.Right)
}

// ── Approach 2: Iterative BFS ─────────────────────────────────────────────────
//
// isSameTreeBFS solves Same Tree using level-order BFS comparison.
//
// Intuition:
//   Use a queue of (p, q) node pairs. For each pair, check if both nil, one
//   nil, or different values. If equal, enqueue their children.
//
// Time:  O(n)
// Space: O(w) — w = max width of tree (O(n) worst case for complete tree).
func isSameTreeBFS(p *TreeNode, q *TreeNode) bool {
	type pair struct{ p, q *TreeNode }
	queue := []pair{{p, q}}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr.p == nil && curr.q == nil {
			continue
		}
		if curr.p == nil || curr.q == nil {
			return false
		}
		if curr.p.Val != curr.q.Val {
			return false
		}
		queue = append(queue, pair{curr.p.Left, curr.q.Left})
		queue = append(queue, pair{curr.p.Right, curr.q.Right})
	}
	return true
}

func main() {
	// [1,2,3] vs [1,2,3] → true
	p1 := &TreeNode{Val: 1, Left: &TreeNode{Val: 2}, Right: &TreeNode{Val: 3}}
	q1 := &TreeNode{Val: 1, Left: &TreeNode{Val: 2}, Right: &TreeNode{Val: 3}}

	// [1,2] vs [1,null,2] → false
	p2 := &TreeNode{Val: 1, Left: &TreeNode{Val: 2}}
	q2 := &TreeNode{Val: 1, Right: &TreeNode{Val: 2}}

	// [1,2,1] vs [1,1,2] → false
	p3 := &TreeNode{Val: 1, Left: &TreeNode{Val: 2}, Right: &TreeNode{Val: 1}}
	q3 := &TreeNode{Val: 1, Left: &TreeNode{Val: 1}, Right: &TreeNode{Val: 2}}

	fmt.Println("=== Approach 1: Recursive DFS ===")
	fmt.Printf("p=[1,2,3] q=[1,2,3]  got=%v  expected true\n", isSameTree(p1, q1))
	fmt.Printf("p=[1,2] q=[1,null,2]  got=%v  expected false\n", isSameTree(p2, q2))
	fmt.Printf("p=[1,2,1] q=[1,1,2]  got=%v  expected false\n", isSameTree(p3, q3))

	// Rebuild for BFS test
	p4 := &TreeNode{Val: 1, Left: &TreeNode{Val: 2}, Right: &TreeNode{Val: 3}}
	q4 := &TreeNode{Val: 1, Left: &TreeNode{Val: 2}, Right: &TreeNode{Val: 3}}
	p5 := &TreeNode{Val: 1, Left: &TreeNode{Val: 2}}
	q5 := &TreeNode{Val: 1, Right: &TreeNode{Val: 2}}

	fmt.Println("=== Approach 2: Iterative BFS ===")
	fmt.Printf("p=[1,2,3] q=[1,2,3]  got=%v  expected true\n", isSameTreeBFS(p4, q4))
	fmt.Printf("p=[1,2] q=[1,null,2]  got=%v  expected false\n", isSameTreeBFS(p5, q5))
}
