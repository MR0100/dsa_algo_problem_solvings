package main

import "fmt"

// TreeNode is a standard binary-tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// ── Approach 1: In-Order Traversal (Brute Force) ─────────────────────────────
//
// inorderTraversal solves Inorder Successor in BST by producing the full
// in-order sequence, then returning the node that comes right after p.
//
// Intuition:
//
//	The in-order traversal of a BST lists nodes in ascending value. The
//	"successor" of p is simply the next node after p in that list. So flatten
//	the tree to a sorted node list and return the element following p (nil if p
//	is last). Ignores the BST shortcut but is obviously correct.
//
// Algorithm:
//
//  1. In-order traverse, appending node pointers to a slice.
//  2. Find p in the slice; return the next element, or nil if p is last.
//
// Time:  O(n) — visits every node.
// Space: O(n) — the flattened list plus recursion stack.
func inorderTraversal(root, p *TreeNode) *TreeNode {
	order := []*TreeNode{}
	var walk func(*TreeNode)
	walk = func(n *TreeNode) {
		if n == nil {
			return
		}
		walk(n.Left)             // smaller values first
		order = append(order, n) // visit
		walk(n.Right)            // larger values last
	}
	walk(root)
	for i, n := range order {
		if n == p && i+1 < len(order) {
			return order[i+1] // node immediately after p
		}
	}
	return nil // p was the maximum → no successor
}

// ── Approach 2: BST Property, Iterative (Optimal) ────────────────────────────
//
// bstSearch solves Inorder Successor in BST in O(h) by exploiting ordering:
// the successor is the smallest node whose value is strictly greater than p's.
//
// Intuition:
//
//	Two cases. If p has a right subtree, the successor is the leftmost node of
//	that right subtree (the smallest value still larger than p). Otherwise the
//	successor is the lowest ancestor for which p lies in its left subtree — we
//	find it by walking down from the root, and every time we go left we record
//	that node as a candidate successor (it's larger than p and we might find a
//	smaller such node deeper). This never visits both subtrees, so it's O(h).
//
// Algorithm:
//
//	successor = nil; cur = root.
//	While cur != nil:
//	  if p.Val < cur.Val: successor = cur; cur = cur.Left   (cur is a candidate)
//	  else:               cur = cur.Right                   (successor is larger)
//	Return successor.
//
// Time:  O(h) — one root-to-leaf descent (h = tree height).
// Space: O(1) — iterative, no recursion.
func bstSearch(root, p *TreeNode) *TreeNode {
	var successor *TreeNode // best candidate seen so far (> p.Val)
	cur := root
	for cur != nil {
		if p.Val < cur.Val {
			successor = cur // cur is larger than p — remember, then look for a smaller one
			cur = cur.Left
		} else {
			cur = cur.Right // cur <= p — the successor must be even larger
		}
	}
	return successor
}

// ── Approach 3: Right-Subtree + Ancestor (Explicit Cases) ────────────────────
//
// rightSubtreeCase solves Inorder Successor in BST by handling the two cases
// separately: leftmost of the right subtree, or the last "turn-left" ancestor.
//
// Intuition:
//
//	Same ordering facts as Approach 2, spelled out. If p.Right exists, dive
//	right once then go left as far as possible — that leaf-ward node is the
//	successor. If not, walk down from root comparing values; the deepest
//	ancestor we descended left from is the successor. Makes the "two cases"
//	structure explicit, which is how many people reason about it.
//
// Algorithm:
//
//	if p.Right != nil: cur = p.Right; while cur.Left != nil { cur = cur.Left };
//	                   return cur.
//	else: walk from root; each time p.Val < cur.Val, set successor = cur and go
//	      left; else go right; stop at p; return successor.
//
// Time:  O(h).
// Space: O(1).
func rightSubtreeCase(root, p *TreeNode) *TreeNode {
	// Case 1: p has a right subtree → successor is its leftmost node.
	if p.Right != nil {
		cur := p.Right
		for cur.Left != nil {
			cur = cur.Left // smallest value in the right subtree
		}
		return cur
	}
	// Case 2: no right subtree → nearest ancestor via a left turn.
	var successor *TreeNode
	cur := root
	for cur != nil && cur.Val != p.Val {
		if p.Val < cur.Val {
			successor = cur // turned left here → candidate successor
			cur = cur.Left
		} else {
			cur = cur.Right
		}
	}
	return successor
}

// val prints a node's value or "nil" for readable expected outputs.
func val(n *TreeNode) string {
	if n == nil {
		return "nil"
	}
	return fmt.Sprintf("%d", n.Val)
}

// buildTree1 constructs Example 1's BST:  [2,1,3], p = 1.
//
//	  2
//	 / \
//	1   3
func buildTree1() (*TreeNode, *TreeNode) {
	n1 := &TreeNode{Val: 1}
	n3 := &TreeNode{Val: 3}
	root := &TreeNode{Val: 2, Left: n1, Right: n3}
	return root, n1 // successor of 1 is 2
}

// buildTree2 constructs Example 2's BST:  [5,3,6,2,4,null,null,1], p = 6.
//
//	      5
//	     / \
//	    3   6
//	   / \
//	  2   4
//	 /
//	1
func buildTree2() (*TreeNode, *TreeNode) {
	n1 := &TreeNode{Val: 1}
	n2 := &TreeNode{Val: 2, Left: n1}
	n4 := &TreeNode{Val: 4}
	n3 := &TreeNode{Val: 3, Left: n2, Right: n4}
	n6 := &TreeNode{Val: 6}
	root := &TreeNode{Val: 5, Left: n3, Right: n6}
	return root, n6 // 6 is the maximum → no successor (nil)
}

func main() {
	fmt.Println("=== Approach 1: In-Order Traversal ===")
	r1, p1 := buildTree1()
	fmt.Println(val(inorderTraversal(r1, p1))) // 2
	r2, p2 := buildTree2()
	fmt.Println(val(inorderTraversal(r2, p2))) // nil

	fmt.Println("=== Approach 2: BST Property, Iterative (Optimal) ===")
	r1, p1 = buildTree1()
	fmt.Println(val(bstSearch(r1, p1))) // 2
	r2, p2 = buildTree2()
	fmt.Println(val(bstSearch(r2, p2))) // nil

	fmt.Println("=== Approach 3: Right-Subtree + Ancestor ===")
	r1, p1 = buildTree1()
	fmt.Println(val(rightSubtreeCase(r1, p1))) // 2
	r2, p2 = buildTree2()
	fmt.Println(val(rightSubtreeCase(r2, p2))) // nil
}
