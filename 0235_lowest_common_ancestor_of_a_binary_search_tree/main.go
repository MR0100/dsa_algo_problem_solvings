package main

import "fmt"

// TreeNode is the standard binary tree node used by LeetCode.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// insert adds val into a BST rooted at root, returning the (possibly new) root.
// Used only to build the example trees.
func insert(root *TreeNode, val int) *TreeNode {
	if root == nil {
		return &TreeNode{Val: val}
	}
	if val < root.Val {
		root.Left = insert(root.Left, val) // smaller keys go left
	} else {
		root.Right = insert(root.Right, val) // larger keys go right
	}
	return root
}

// buildBST builds a BST by inserting vals in order.
func buildBST(vals []int) *TreeNode {
	var root *TreeNode
	for _, v := range vals {
		root = insert(root, v)
	}
	return root
}

// find returns the node whose Val == target (helper to pass p and q by node).
func find(root *TreeNode, target int) *TreeNode {
	cur := root
	for cur != nil {
		switch {
		case target < cur.Val:
			cur = cur.Left
		case target > cur.Val:
			cur = cur.Right
		default:
			return cur // found the node with this value
		}
	}
	return nil
}

// ── Approach 1: Recursive BST Walk ───────────────────────────────────────────
//
// recursiveBST solves LCA of a BST by exploiting the BST ordering to descend
// toward the split point of p and q.
//
// Intuition:
//
//	In a BST every node partitions keys: everything smaller is left, everything
//	larger is right. The lowest common ancestor of p and q is the first node
//	where p and q fall on different sides (or one equals the node). If both are
//	smaller than the current node they share an ancestor on the left; if both
//	larger, on the right; otherwise the current node is the split point — the
//	LCA.
//
// Algorithm:
//
//  1. At node root: if both p.Val and q.Val < root.Val, recurse left.
//  2. Else if both > root.Val, recurse right.
//  3. Else root is the LCA (the values straddle it, or one is root).
//
// Time:  O(h) — one node per level, h the tree height.
// Space: O(h) — recursion stack.
func recursiveBST(root, p, q *TreeNode) *TreeNode {
	if p.Val < root.Val && q.Val < root.Val {
		return recursiveBST(root.Left, p, q) // both keys live in the left subtree
	}
	if p.Val > root.Val && q.Val > root.Val {
		return recursiveBST(root.Right, p, q) // both keys live in the right subtree
	}
	return root // p and q split here (or one equals root) → lowest common ancestor
}

// ── Approach 2: Iterative BST Walk (Optimal) ─────────────────────────────────
//
// iterativeBST solves LCA of a BST with the same ordering logic but a loop,
// achieving O(1) extra space.
//
// Intuition:
//
//	The recursion is tail-shaped: each step just moves to one child. Replacing
//	it with a while loop removes the call stack entirely, giving the optimal
//	O(h) time and O(1) space solution.
//
// Algorithm:
//
//  1. Start cur = root.
//  2. While cur != nil:
//     - if both p,q < cur.Val: cur = cur.Left.
//     - else if both p,q > cur.Val: cur = cur.Right.
//     - else return cur (split point).
//
// Time:  O(h) — descends one level per iteration.
// Space: O(1) — a single pointer.
func iterativeBST(root, p, q *TreeNode) *TreeNode {
	cur := root
	for cur != nil {
		switch {
		case p.Val < cur.Val && q.Val < cur.Val:
			cur = cur.Left // both smaller → go left
		case p.Val > cur.Val && q.Val > cur.Val:
			cur = cur.Right // both larger → go right
		default:
			return cur // values straddle cur (or one equals it) → LCA
		}
	}
	return nil // unreachable given p and q exist in the tree
}

// ── Approach 3: General LCA via Root-to-Node Paths ───────────────────────────
//
// pathIntersection solves LCA by recording the root-to-p and root-to-q paths
// and returning the last shared node — a technique that works on any binary
// tree, not just BSTs.
//
// Intuition:
//
//	Ignore the BST property for a moment: the LCA of two nodes is the last
//	common node on their two root-to-node paths. Collect both paths as lists of
//	nodes, then walk them in parallel; the deepest position where they still
//	agree is the LCA. It is slower and heavier than the BST-specific walk but
//	generalises to arbitrary trees (see LC #236).
//
// Algorithm:
//
//  1. Using the BST ordering, build pathP = root..p and pathQ = root..q.
//  2. Walk both from the root while pathP[i] == pathQ[i].
//  3. The last matching node is the LCA.
//
// Time:  O(h) — building each path costs O(h); comparison is O(h).
// Space: O(h) — the two path slices.
func pathIntersection(root, p, q *TreeNode) *TreeNode {
	// pathTo returns the list of nodes from root down to the target, using the
	// BST ordering to choose a direction at each step.
	pathTo := func(target *TreeNode) []*TreeNode {
		path := []*TreeNode{}
		cur := root
		for cur != nil {
			path = append(path, cur) // record every node on the way down
			if target.Val < cur.Val {
				cur = cur.Left
			} else if target.Val > cur.Val {
				cur = cur.Right
			} else {
				break // reached the target node itself
			}
		}
		return path
	}

	pathP := pathTo(p)
	pathQ := pathTo(q)

	var lca *TreeNode
	// Walk both paths from the root; the last shared node is the LCA.
	for i := 0; i < len(pathP) && i < len(pathQ); i++ {
		if pathP[i] == pathQ[i] {
			lca = pathP[i] // still on the common prefix → update the candidate
		} else {
			break // paths diverge here; earlier node stays as the LCA
		}
	}
	return lca
}

func main() {
	// Example tree: [6,2,8,0,4,7,9,null,null,3,5]
	tree := buildBST([]int{6, 2, 8, 0, 4, 7, 9, 3, 5})

	fmt.Println("=== Approach 1: Recursive BST Walk ===")
	fmt.Println(recursiveBST(tree, find(tree, 2), find(tree, 8)).Val) // expected 6
	fmt.Println(recursiveBST(tree, find(tree, 2), find(tree, 4)).Val) // expected 2

	fmt.Println("=== Approach 2: Iterative BST Walk (Optimal) ===")
	fmt.Println(iterativeBST(tree, find(tree, 2), find(tree, 8)).Val) // expected 6
	fmt.Println(iterativeBST(tree, find(tree, 2), find(tree, 4)).Val) // expected 2

	fmt.Println("=== Approach 3: General LCA via Root-to-Node Paths ===")
	fmt.Println(pathIntersection(tree, find(tree, 2), find(tree, 8)).Val) // expected 6
	fmt.Println(pathIntersection(tree, find(tree, 2), find(tree, 4)).Val) // expected 2
}
