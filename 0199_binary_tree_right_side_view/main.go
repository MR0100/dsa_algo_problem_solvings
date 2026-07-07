package main

import "fmt"

// TreeNode is a binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// ── Approach 1: BFS Level Order ──────────────────────────────────────────────
//
// bfsLevelOrder solves Binary Tree Right Side View with a breadth-first
// level-order traversal, recording the last node of every level.
//
// Intuition:
//   Standing on the right you see exactly one node per level: the rightmost
//   one. BFS processes the tree level by level, so the final node dequeued
//   within each level is precisely the node visible from the right.
//
// Algorithm:
//   1. If the root is nil, return an empty view.
//   2. Push the root into a queue.
//   3. While the queue is non-empty: size = current level width; dequeue
//      exactly size nodes, enqueueing each node's children; when dequeuing
//      the size-th (last) node, append its value to the view.
//
// Time:  O(n) — every node is enqueued and dequeued exactly once.
// Space: O(w) — the queue holds at most one full level (w = max width,
//         up to ~n/2 for a complete tree).
func bfsLevelOrder(root *TreeNode) []int {
	view := []int{}
	if root == nil {
		return view // empty tree: nothing is visible
	}
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		size := len(queue) // number of nodes on the current level
		for i := 0; i < size; i++ {
			node := queue[0]
			queue = queue[1:]
			if i == size-1 {
				view = append(view, node.Val) // last node of the level = rightmost
			}
			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
	}
	return view
}

// ── Approach 2: DFS Left-First (Overwrite) ───────────────────────────────────
//
// dfsLeftFirst solves Binary Tree Right Side View with an ordinary
// left-to-right DFS that keeps overwriting each depth's slot; the last writer
// at every depth is the rightmost node.
//
// Intuition:
//   In a root → left → right traversal, the LAST node visited at each depth
//   is the rightmost node of that depth. So instead of recording first
//   arrivals, overwrite view[depth] on every visit and let the final write
//   win.
//
// Algorithm:
//   1. DFS with the current depth as a parameter, left child first.
//   2. On visiting a node: if depth == len(view), append its value (first
//      node ever seen this deep); otherwise overwrite view[depth].
//   3. Recurse into Left, then Right.
//
// Time:  O(n) — each node is visited once.
// Space: O(h) — recursion stack, h = tree height (O(n) for a skewed tree).
func dfsLeftFirst(root *TreeNode) []int {
	view := []int{}
	var dfs func(node *TreeNode, depth int)
	dfs = func(node *TreeNode, depth int) {
		if node == nil {
			return
		}
		if depth == len(view) {
			view = append(view, node.Val) // first node seen at this depth
		} else {
			view[depth] = node.Val // a node further right overwrites the slot
		}
		dfs(node.Left, depth+1) // left first: the rightmost node writes last
		dfs(node.Right, depth+1)
	}
	dfs(root, 0)
	return view
}

// ── Approach 3: DFS Right-First (Optimal) ────────────────────────────────────
//
// dfsRightFirst solves Binary Tree Right Side View with a depth-first
// traversal that always explores the right subtree before the left one.
//
// Intuition:
//   Walking root → right → left, the FIRST node reached at any depth is the
//   rightmost node of that depth — every node to its left is visited strictly
//   later. Record a value only on first arrival at a depth
//   (depth == len(view)); all later visits at that depth are ignored.
//
// Algorithm:
//   1. DFS with the current depth as a parameter, right child first.
//   2. On visiting a node: if depth == len(view), this is the first node seen
//      this deep → append its value to the view.
//   3. Recurse into Right, then Left.
//
// Time:  O(n) — each node is visited once.
// Space: O(h) — recursion stack only; no queue, and each answer is found at
//         the earliest possible moment.
func dfsRightFirst(root *TreeNode) []int {
	view := []int{}
	var dfs func(node *TreeNode, depth int)
	dfs = func(node *TreeNode, depth int) {
		if node == nil {
			return
		}
		if depth == len(view) {
			view = append(view, node.Val) // first arrival at this depth = rightmost
		}
		dfs(node.Right, depth+1) // right first, so the rightmost node wins each depth
		dfs(node.Left, depth+1)
	}
	dfs(root, 0)
	return view
}

func main() {
	// Example 1: root = [1,2,3,null,5,null,4]
	//        1
	//       / \
	//      2   3
	//       \   \
	//        5   4
	t1 := &TreeNode{
		Val:   1,
		Left:  &TreeNode{Val: 2, Right: &TreeNode{Val: 5}},
		Right: &TreeNode{Val: 3, Right: &TreeNode{Val: 4}},
	}

	// Example 2: root = [1,2,3,4,null,null,null,5]
	//          1
	//         / \
	//        2   3
	//       /
	//      4
	//     /
	//    5
	t2 := &TreeNode{
		Val:   1,
		Left:  &TreeNode{Val: 2, Left: &TreeNode{Val: 4, Left: &TreeNode{Val: 5}}},
		Right: &TreeNode{Val: 3},
	}

	// Example 3: root = [1,null,3]
	//    1
	//     \
	//      3
	t3 := &TreeNode{Val: 1, Right: &TreeNode{Val: 3}}

	// Example 4: root = []
	var t4 *TreeNode // nil — the empty tree

	fmt.Println("=== Approach 1: BFS Level Order ===")
	fmt.Println(bfsLevelOrder(t1)) // [1 3 4]
	fmt.Println(bfsLevelOrder(t2)) // [1 3 4 5]
	fmt.Println(bfsLevelOrder(t3)) // [1 3]
	fmt.Println(bfsLevelOrder(t4)) // []

	fmt.Println("=== Approach 2: DFS Left-First (Overwrite) ===")
	fmt.Println(dfsLeftFirst(t1)) // [1 3 4]
	fmt.Println(dfsLeftFirst(t2)) // [1 3 4 5]
	fmt.Println(dfsLeftFirst(t3)) // [1 3]
	fmt.Println(dfsLeftFirst(t4)) // []

	fmt.Println("=== Approach 3: DFS Right-First (Optimal) ===")
	fmt.Println(dfsRightFirst(t1)) // [1 3 4]
	fmt.Println(dfsRightFirst(t2)) // [1 3 4 5]
	fmt.Println(dfsRightFirst(t3)) // [1 3]
	fmt.Println(dfsRightFirst(t4)) // []
}
