package main

import (
	"fmt"
	"sort"
)

// TreeNode is the standard LeetCode binary-tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// ── Approach 1: BFS with Column Index (Optimal / Canonical) ──────────────────
//
// bfs solves Binary Tree Vertical Order Traversal by a level-order traversal
// that tags each node with a column: root = 0, left child = col-1, right = col+1.
//
// Intuition:
//
//	Vertical order groups nodes by column. Within a column, nodes must appear
//	top-to-bottom, and left-before-right when they share a row. A BFS visits
//	nodes strictly top-to-bottom and, by enqueuing left child before right
//	child, left-before-right at each level — which is exactly the required tie
//	order. So simply appending each dequeued node to its column's bucket yields
//	the correct within-column order automatically, with NO sorting by row.
//
// Algorithm:
//
//  1. BFS from root, carrying each node's column with it in the queue.
//  2. Track minCol and maxCol seen; store cols[col] = list of values in visit order.
//  3. Emit columns from minCol to maxCol.
//
// Time:  O(n) — each node visited once; O(range) to assemble output.
// Space: O(n) — the queue and the column buckets.
func bfs(root *TreeNode) [][]int {
	if root == nil {
		return [][]int{} // empty tree -> empty result
	}

	cols := map[int][]int{} // column index -> values in top-to-bottom, L-to-R order
	minCol, maxCol := 0, 0  // track the horizontal range actually used

	// item pairs a node with its column so the queue carries both.
	type item struct {
		node *TreeNode
		col  int
	}
	queue := []item{{root, 0}} // root sits at column 0

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:] // dequeue front (FIFO -> level order, top to bottom)

		cols[cur.col] = append(cols[cur.col], cur.node.Val) // record in visit order
		if cur.col < minCol {
			minCol = cur.col
		}
		if cur.col > maxCol {
			maxCol = cur.col
		}

		// Enqueue left BEFORE right so same-row ties come out left-first.
		if cur.node.Left != nil {
			queue = append(queue, item{cur.node.Left, cur.col - 1})
		}
		if cur.node.Right != nil {
			queue = append(queue, item{cur.node.Right, cur.col + 1})
		}
	}

	// Assemble columns left-to-right by walking the known range.
	res := make([][]int, 0, maxCol-minCol+1)
	for c := minCol; c <= maxCol; c++ {
		res = append(res, cols[c])
	}
	return res
}

// ── Approach 2: DFS with (col, row) then Stable Sort ─────────────────────────
//
// dfsSort solves Binary Tree Vertical Order Traversal with a depth-first walk
// that records each node's (column, row), then sorts entries within a column by
// row (breaking ties by visit order) since DFS does NOT visit strictly top-down.
//
// Intuition:
//
//	DFS is a natural recursive traversal but it dives deep before going wide, so
//	nodes are not produced in top-to-bottom order. We therefore also record each
//	node's row (depth). To reconstruct vertical order we sort each column's
//	entries by row; a STABLE sort keeps left-before-right for equal (col,row)
//	because we recurse left before right, so left nodes are appended first.
//
// Algorithm:
//
//  1. DFS carrying (col, row); append (row, seq, val) into cols[col].
//  2. For each column, stable-sort by row (seq preserves L-before-R ties).
//  3. Emit columns from minCol to maxCol.
//
// Time:  O(n log n) — dominated by per-column sorting.
// Space: O(n) — stored entries and recursion stack.
func dfsSort(root *TreeNode) [][]int {
	if root == nil {
		return [][]int{}
	}

	// rowVal pairs a node's row (depth) with its value for later sorting.
	type rowVal struct {
		row int
		val int
	}
	cols := map[int][]rowVal{}
	minCol, maxCol := 0, 0

	var dfs func(node *TreeNode, col, row int)
	dfs = func(node *TreeNode, col, row int) {
		if node == nil {
			return
		}
		cols[col] = append(cols[col], rowVal{row, node.Val})
		if col < minCol {
			minCol = col
		}
		if col > maxCol {
			maxCol = col
		}
		dfs(node.Left, col-1, row+1)  // recurse left first: preserves L-before-R ties
		dfs(node.Right, col+1, row+1) // then right
	}
	dfs(root, 0, 0)

	res := make([][]int, 0, maxCol-minCol+1)
	for c := minCol; c <= maxCol; c++ {
		entries := cols[c]
		// Stable sort by row; equal rows keep DFS append order (left first).
		sort.SliceStable(entries, func(i, j int) bool {
			return entries[i].row < entries[j].row
		})
		vals := make([]int, len(entries))
		for i, e := range entries {
			vals[i] = e.val
		}
		res = append(res, vals)
	}
	return res
}

// buildExample1 builds [3,9,20,null,null,15,7].
func buildExample1() *TreeNode {
	return &TreeNode{
		Val:  3,
		Left: &TreeNode{Val: 9},
		Right: &TreeNode{
			Val:   20,
			Left:  &TreeNode{Val: 15},
			Right: &TreeNode{Val: 7},
		},
	}
}

// buildExample2 builds [3,9,8,4,0,1,7].
func buildExample2() *TreeNode {
	return &TreeNode{
		Val: 3,
		Left: &TreeNode{
			Val:   9,
			Left:  &TreeNode{Val: 4},
			Right: &TreeNode{Val: 0},
		},
		Right: &TreeNode{
			Val:   8,
			Left:  &TreeNode{Val: 1},
			Right: &TreeNode{Val: 7},
		},
	}
}

// buildExample3 builds [3,9,8,4,0,1,7,null,null,null,2,5].
// Node 0 (left child of 9) has right child 2; node 1 (left child of 8) has
// left child 5.
func buildExample3() *TreeNode {
	zero := &TreeNode{Val: 0, Right: &TreeNode{Val: 2}}
	one := &TreeNode{Val: 1, Left: &TreeNode{Val: 5}}
	return &TreeNode{
		Val: 3,
		Left: &TreeNode{
			Val:   9,
			Left:  &TreeNode{Val: 4},
			Right: zero,
		},
		Right: &TreeNode{
			Val:   8,
			Left:  one,
			Right: &TreeNode{Val: 7},
		},
	}
}

func main() {
	// Official Example 1: [[9],[3,15],[20],[7]]
	fmt.Println("=== Approach 1: BFS with Column Index ===")
	fmt.Println(bfs(buildExample1())) // expected [[9] [3 15] [20] [7]]
	fmt.Println("=== Approach 2: DFS with Row Sort ===")
	fmt.Println(dfsSort(buildExample1())) // expected [[9] [3 15] [20] [7]]

	// Official Example 2: [[4],[9],[3,0,1],[8],[7]]
	fmt.Println("=== Approach 1: BFS (Example 2) ===")
	fmt.Println(bfs(buildExample2())) // expected [[4] [9] [3 0 1] [8] [7]]
	fmt.Println("=== Approach 2: DFS (Example 2) ===")
	fmt.Println(dfsSort(buildExample2())) // expected [[4] [9] [3 0 1] [8] [7]]

	// Official Example 3: [[4],[9,5],[3,0,1],[8,2],[7]]
	fmt.Println("=== Approach 1: BFS (Example 3) ===")
	fmt.Println(bfs(buildExample3())) // expected [[4] [9 5] [3 0 1] [8 2] [7]]
	fmt.Println("=== Approach 2: DFS (Example 3) ===")
	fmt.Println(dfsSort(buildExample3())) // expected [[4] [9 5] [3 0 1] [8 2] [7]]
}
