package main

import (
	"fmt"
	"strings"
)

// Node is a Quad-Tree node (LeetCode #427). A node is either a leaf (isLeaf
// true, Val the uniform value of its sub-grid) or an internal node with four
// children covering the four quadrants.
type Node struct {
	Val         bool
	IsLeaf      bool
	TopLeft     *Node
	TopRight    *Node
	BottomLeft  *Node
	BottomRight *Node
}

// ── Approach 1: Recursive Divide and Conquer (Scan Each Sub-grid) ─────────────
//
// divideAndConquer builds the quad tree by recursively splitting the current
// square into four quadrants, checking uniformity by scanning the sub-grid.
//
// Intuition:
//
//	A square region is representable by one leaf iff every cell in it is equal.
//	Check that directly; if uniform, emit a leaf. Otherwise the region genuinely
//	needs four children — recurse on each of the four equal quadrants (each of
//	side length half the current one) and wrap them in an internal node. Because
//	n is a power of two, halving always yields four equal squares down to 1×1.
//
// Algorithm:
//  1. build(r, c, size): if the size×size block at (r,c) is uniform, return a
//     leaf with that value.
//  2. Otherwise split into four size/2 quadrants and recurse:
//     topLeft (r,c), topRight (r, c+half), bottomLeft (r+half, c),
//     bottomRight (r+half, c+half).
//  3. Return an internal node (IsLeaf=false) holding the four children.
//
// Time:  O(n² log n) worst case — a fully-mixed grid re-scans overlapping
//
//	regions across log n levels (Σ level-work = n² per level × log n levels).
//
// Space: O(log n) recursion depth plus O(#nodes) for the tree.
func divideAndConquer(grid [][]int) *Node {
	n := len(grid)
	var build func(r, c, size int) *Node
	build = func(r, c, size int) *Node {
		if isUniform(grid, r, c, size) {
			// Whole block equal ⇒ one leaf. grid[r][c]==1 → Val true.
			return &Node{Val: grid[r][c] == 1, IsLeaf: true}
		}
		half := size / 2 // each quadrant is half the side length
		return &Node{
			Val:         true, // arbitrary for internal nodes; true is conventional
			IsLeaf:      false,
			TopLeft:     build(r, c, half),
			TopRight:    build(r, c+half, half),
			BottomLeft:  build(r+half, c, half),
			BottomRight: build(r+half, c+half, half),
		}
	}
	return build(0, 0, n)
}

// isUniform reports whether every cell in the size×size block at (r,c) equals
// the top-left cell of that block.
func isUniform(grid [][]int, r, c, size int) bool {
	first := grid[r][c] // reference value for the block
	for i := r; i < r+size; i++ {
		for j := c; j < c+size; j++ {
			if grid[i][j] != first {
				return false // found a differing cell ⇒ not uniform
			}
		}
	}
	return true
}

// ── Approach 2: Divide and Conquer Merging Children (Optimal) ─────────────────
//
// mergeChildren builds the quad tree bottom-up: always split to the four
// children first, then MERGE them back into a single leaf when all four turn
// out to be leaves with the same value.
//
// Intuition:
//
//	Instead of scanning a whole block to test uniformity (which repeats work),
//	recurse to the four quadrants unconditionally down to 1×1 leaves, then ask a
//	local question: "are my four children all leaves carrying the same value?"
//	If yes, they collapse into one leaf (the block was uniform after all);
//	otherwise keep them as children. Each cell is read exactly once, so the total
//	work is linear in the grid.
//
// Algorithm:
//  1. build(r, c, size): if size == 1, return a leaf for that single cell.
//  2. Recurse into the four half-size quadrants.
//  3. If all four are leaves AND share the same Val, return a single merged
//     leaf with that Val (drop the children).
//  4. Otherwise return an internal node holding the four children.
//
// Time:  O(n²) — every cell contributes to exactly one base-case leaf; merges
//
//	are O(1) per internal node and there are O(n²) nodes total.
//
// Space: O(log n) recursion depth plus O(#nodes) for the tree.
func mergeChildren(grid [][]int) *Node {
	var build func(r, c, size int) *Node
	build = func(r, c, size int) *Node {
		if size == 1 {
			// Base case: a single cell is always a leaf.
			return &Node{Val: grid[r][c] == 1, IsLeaf: true}
		}
		half := size / 2
		tl := build(r, c, half)
		tr := build(r, c+half, half)
		bl := build(r+half, c, half)
		br := build(r+half, c+half, half)

		// Collapse iff all four children are leaves with an identical value.
		if tl.IsLeaf && tr.IsLeaf && bl.IsLeaf && br.IsLeaf &&
			tl.Val == tr.Val && tr.Val == bl.Val && bl.Val == br.Val {
			return &Node{Val: tl.Val, IsLeaf: true} // merged leaf; children discarded
		}
		return &Node{ // genuinely mixed region: keep the four children
			Val:         true,
			IsLeaf:      false,
			TopLeft:     tl,
			TopRight:    tr,
			BottomLeft:  bl,
			BottomRight: br,
		}
	}
	return build(0, 0, len(grid))
}

// ── serialization for verification (LeetCode's level-order format) ───────────

// serialize renders a quad tree in LeetCode's exact output form: a level-order
// (BFS) list where each present node prints as [isLeaf, val] (1/0). Every node —
// leaf OR internal — contributes four child slots; a leaf's four children are
// null, and a null node contributes no further slots. Trailing nulls are then
// trimmed, matching LeetCode's serializer.
func serialize(root *Node) string {
	if root == nil {
		return "[]"
	}
	parts := []string{}
	queue := []*Node{root}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		if node == nil {
			parts = append(parts, "null") // placeholder for a missing child
			continue
		}
		parts = append(parts, fmt.Sprintf("[%d,%d]", b2i(node.IsLeaf), b2i(node.Val)))
		if node.IsLeaf {
			// A leaf still occupies four child slots, all null.
			queue = append(queue, nil, nil, nil, nil)
		} else {
			queue = append(queue, node.TopLeft, node.TopRight, node.BottomLeft, node.BottomRight)
		}
	}
	// Trim trailing "null" placeholders (LeetCode omits them).
	end := len(parts)
	for end > 0 && parts[end-1] == "null" {
		end--
	}
	return "[" + strings.Join(parts[:end], ",") + "]"
}

// b2i maps a bool to 1/0 for the [isLeaf,val] pair.
func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

func main() {
	// Example 1: grid = [[0,1],[1,0]] — every cell differs from its neighbours,
	// so the root splits into four 1×1 leaves.
	grid1 := [][]int{
		{0, 1},
		{1, 0},
	}
	fmt.Println("=== Approach 1: Divide & Conquer (scan) — Example 1 ===")
	fmt.Println(serialize(divideAndConquer(grid1))) // [[0,1],[1,0],[1,1],[1,1],[1,0]]
	fmt.Println("=== Approach 2: Merge Children (Optimal) — Example 1 ===")
	fmt.Println(serialize(mergeChildren(grid1))) // [[0,1],[1,0],[1,1],[1,1],[1,0]]

	// Example 2: 8×8 grid. TopLeft quadrant is all 1s (→ single leaf); the other
	// three quadrants are mixed.
	grid2 := [][]int{
		{1, 1, 1, 1, 0, 0, 0, 0},
		{1, 1, 1, 1, 0, 0, 0, 0},
		{1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 0, 0, 0, 0},
		{1, 1, 1, 1, 0, 0, 0, 0},
		{1, 1, 1, 1, 0, 0, 0, 0},
		{1, 1, 1, 1, 0, 0, 0, 0},
	}
	fmt.Println("=== Approach 1: Divide & Conquer (scan) — Example 2 ===")
	fmt.Println(serialize(divideAndConquer(grid2))) // [[0,1],[1,1],[0,1],[1,1],[1,0],null,null,null,null,[1,0],[1,0],[1,1],[1,1]]
	fmt.Println("=== Approach 2: Merge Children (Optimal) — Example 2 ===")
	fmt.Println(serialize(mergeChildren(grid2))) // [[0,1],[1,1],[0,1],[1,1],[1,0],null,null,null,null,[1,0],[1,0],[1,1],[1,1]]
	// Note: [1,1]=leaf 1, [1,0]=leaf 0, [0,1]=internal node. The four null,null,
	// null,null are the TopLeft leaf's (absent) children; trailing nulls trimmed.
}
