package main

import "fmt"

// TreeNode is the standard LeetCode binary-tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// ── Approach 1: Brute Force (Full Traversal) ─────────────────────────────────
//
// bruteForce counts nodes by visiting every node once, ignoring the
// "complete tree" guarantee entirely.
//
// Intuition:
//
//	The count of any tree is 1 (this node) plus the counts of its two
//	subtrees. This works for arbitrary trees; the completeness property is
//	simply not exploited, so we pay for touching every node.
//
// Algorithm:
//
//	count(nil) = 0; count(node) = 1 + count(left) + count(right).
//
// Time:  O(n) — every node visited exactly once.
// Space: O(h) — recursion stack, h = tree height (≈ log n for a complete tree).
func bruteForce(root *TreeNode) int {
	if root == nil {
		return 0 // empty subtree contributes nothing
	}
	// this node + everything under it
	return 1 + bruteForce(root.Left) + bruteForce(root.Right)
}

// ── Approach 2: Perfect-Subtree Detection (Optimal) ──────────────────────────
//
// perfectSubtree counts nodes in O(log²n) by using the complete-tree
// structure: at every node it compares the leftmost and rightmost path
// depths; equal depths mean a *perfect* subtree whose size is 2^h − 1 with no
// recursion needed.
//
// Intuition:
//
//	In a complete tree, walk left-only from a node to get its left height and
//	right-only to get its right height. If they match, the subtree is
//	perfect and holds exactly 2^height − 1 nodes — computed instantly. If
//	they differ, only one side is "incomplete", so recurse into both but at
//	least one recursion terminates in the fast perfect case quickly. Each
//	level triggers one height walk (O(log n)) and we descend O(log n) levels.
//
// Algorithm:
//
//  1. If root is nil, return 0.
//  2. lh = length of the all-left path; rh = length of the all-right path.
//  3. If lh == rh, the subtree is perfect: return 2^lh − 1 (as (1<<lh) − 1).
//  4. Otherwise return 1 + count(left) + count(right).
//
// Time:  O(log²n) — O(log n) levels, each doing an O(log n) height walk.
// Space: O(log n) — recursion depth equals the tree height.
func perfectSubtree(root *TreeNode) int {
	if root == nil {
		return 0 // empty subtree
	}
	lh := leftHeight(root)  // depth following only left children
	rh := rightHeight(root) // depth following only right children
	if lh == rh {
		// perfect subtree of height lh → 2^lh − 1 nodes, no recursion
		return (1 << lh) - 1
	}
	// heights differ: split and recurse; completeness makes this cheap
	return 1 + perfectSubtree(root.Left) + perfectSubtree(root.Right)
}

// leftHeight returns how many edges the leftmost root-to-leaf path has.
func leftHeight(node *TreeNode) int {
	h := 0
	for node != nil {
		h++              // count this level
		node = node.Left // keep hugging the left edge
	}
	return h
}

// rightHeight returns how many edges the rightmost root-to-leaf path has.
func rightHeight(node *TreeNode) int {
	h := 0
	for node != nil {
		h++               // count this level
		node = node.Right // keep hugging the right edge
	}
	return h
}

// ── Approach 3: Binary Search on the Last Level ──────────────────────────────
//
// binarySearchLastLevel counts nodes by computing the number of full upper
// levels, then binary-searching how many leaves are present on the (possibly
// partial) last level.
//
// Intuition:
//
//	A complete tree of height h has (2^h − 1) nodes in its top h levels, all
//	guaranteed present. The last level holds between 1 and 2^h leaves, filled
//	left to right. Index those leaf slots 0 … 2^h − 1; "is slot i present?"
//	is monotone (present slots form a prefix), so binary-search the largest
//	present index. Checking a slot walks down the tree using the bits of i to
//	choose left/right — O(h) per check.
//
// Algorithm:
//
//  1. h = number of edges on the leftmost path (tree height).
//  2. If h == 0 the tree is a single node → return 1.
//  3. Binary-search lo=0, hi=2^h−1 over last-level leaf indices; `exists(i)`
//     walks h steps guided by the h bits of i.
//  4. Answer = (2^h − 1) upper nodes + (last present index + 1) last-level nodes.
//
// Time:  O(log²n) — O(log n) binary-search steps, each an O(log n) walk.
// Space: O(1) — iterative, no recursion.
func binarySearchLastLevel(root *TreeNode) int {
	if root == nil {
		return 0
	}
	// height h = edges from root down the leftmost path
	h := 0
	for n := root.Left; n != nil; n = n.Left {
		h++
	}
	if h == 0 {
		return 1 // only the root exists
	}
	// last-level leaf indices range over [0, 2^h − 1]
	lo, hi := 0, (1<<h)-1
	for lo <= hi {
		mid := (lo + hi) / 2
		if exists(root, mid, h) {
			lo = mid + 1 // slot present → answer is at least mid, search right
		} else {
			hi = mid - 1 // slot absent → search left
		}
	}
	// upper (2^h − 1) full nodes + `lo` present leaves (lo = count past hi)
	return (1<<h - 1) + lo
}

// exists reports whether the last-level leaf at index `idx` (0-based, among
// 2^h slots) is present, by walking h levels: bit (h-1-step) of idx picks
// right (1) or left (0).
func exists(root *TreeNode, idx, h int) bool {
	node := root
	for bit := h - 1; bit >= 0; bit-- {
		if idx&(1<<bit) != 0 { // this bit set → go right
			node = node.Right
		} else { // bit clear → go left
			node = node.Left
		}
		if node == nil {
			return false // path breaks → slot not filled
		}
	}
	return true // reached the leaf slot → present
}

// buildComplete builds a complete binary tree with `n` nodes valued 1..n in
// level order (node i's children are 2i and 2i+1, 1-based), matching how
// LeetCode serialises complete trees.
func buildComplete(n int) *TreeNode {
	if n == 0 {
		return nil
	}
	nodes := make([]*TreeNode, n+1) // 1-based
	for i := 1; i <= n; i++ {
		nodes[i] = &TreeNode{Val: i}
	}
	for i := 1; i <= n; i++ {
		if 2*i <= n {
			nodes[i].Left = nodes[2*i]
		}
		if 2*i+1 <= n {
			nodes[i].Right = nodes[2*i+1]
		}
	}
	return nodes[1]
}

func main() {
	// Example 1: root = [1,2,3,4,5,6] → 6 nodes.
	t1 := buildComplete(6)
	// Example 2: root = [] → 0 nodes.
	var t2 *TreeNode = buildComplete(0)
	// Example 3: root = [1] → 1 node.
	t3 := buildComplete(1)

	fmt.Println("=== Approach 1: Brute Force (Full Traversal) ===")
	fmt.Println(bruteForce(t1)) // expected 6
	fmt.Println(bruteForce(t2)) // expected 0
	fmt.Println(bruteForce(t3)) // expected 1

	fmt.Println("=== Approach 2: Perfect-Subtree Detection (Optimal) ===")
	fmt.Println(perfectSubtree(t1)) // expected 6
	fmt.Println(perfectSubtree(t2)) // expected 0
	fmt.Println(perfectSubtree(t3)) // expected 1

	fmt.Println("=== Approach 3: Binary Search on the Last Level ===")
	fmt.Println(binarySearchLastLevel(t1)) // expected 6
	fmt.Println(binarySearchLastLevel(t2)) // expected 0
	fmt.Println(binarySearchLastLevel(t3)) // expected 1
}
