package main

import (
	"fmt"
	"strconv"
	"strings"
)

// TreeNode is the standard LeetCode binary-tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// ── Approach 1: Recursive Delete (Inorder Successor) ─────────────────────────
//
// recursiveDelete solves Delete Node in a BST by recursing to the target using
// the BST ordering, then splicing it out with the three classic deletion cases.
//
// Intuition:
//
//	Deleting from a BST has three shapes once the node is found:
//	  (a) leaf / one child → just return the (possibly nil) child to reattach;
//	  (b) two children → we cannot simply drop it, so replace its value with its
//	      INORDER SUCCESSOR (smallest value in the right subtree, which has no
//	      left child), then recursively delete that successor from the right
//	      subtree. The successor is chosen because it is the next-larger value,
//	      so the BST ordering is preserved after the swap.
//	Recursion naturally reattaches subtrees: each call returns the new root of
//	the subtree it was handed, and the parent stores that return value.
//
// Algorithm:
//  1. If root is nil, return nil (key not present).
//  2. If key < root.Val, delete in left subtree; if key > root.Val, delete in right.
//  3. Else root is the target:
//     - if no left child, return root.Right; if no right child, return root.Left;
//     - otherwise find min of right subtree, copy its value into root, and
//     delete that value from root.Right.
//  4. Return root.
//
// Time:  O(h) — one root-to-target path plus a successor descent, h = tree height.
// Space: O(h) — recursion stack (O(log n) balanced, O(n) skewed).
func recursiveDelete(root *TreeNode, key int) *TreeNode {
	if root == nil {
		return nil // key not found; nothing to remove
	}
	if key < root.Val {
		root.Left = recursiveDelete(root.Left, key) // target is in the left subtree
	} else if key > root.Val {
		root.Right = recursiveDelete(root.Right, key) // target is in the right subtree
	} else {
		// Found the node to delete — handle by number of children.
		if root.Left == nil {
			return root.Right // 0 or 1 child: promote the right child
		}
		if root.Right == nil {
			return root.Left // exactly one (left) child: promote it
		}
		// Two children: find inorder successor = leftmost node of right subtree.
		succ := root.Right
		for succ.Left != nil {
			succ = succ.Left // walk to the smallest value greater than root
		}
		root.Val = succ.Val // overwrite value with successor's (ordering kept)
		// Remove the successor (a node with no left child) from the right subtree.
		root.Right = recursiveDelete(root.Right, succ.Val)
	}
	return root
}

// ── Approach 2: Iterative Delete (Parent Pointer) ────────────────────────────
//
// iterativeDelete solves Delete Node in a BST without recursion by first
// locating the node (and its parent) via a loop, then removing it in place.
//
// Intuition:
//
//	The same three-case deletion, but done iteratively to use O(1) auxiliary
//	space. Walk down tracking the parent so we know which child pointer to
//	rewire. When the target is found, compute its replacement subtree:
//	  - at most one child → the replacement is that child;
//	  - two children → detach the inorder successor from the right subtree
//	    (tracking ITS parent too) and splice it into the target's position,
//	    adopting the target's original left and right subtrees.
//	Finally hang the replacement off the parent (or make it the new root).
//
// Algorithm:
//  1. Find target and its parent with a search loop; if not found, return root.
//  2. If target has two children, detach its inorder successor and let that
//     successor become the replacement (adopting target's children).
//     Else the replacement is target's single (or nil) child.
//  3. Reconnect: if target was the root, return replacement; otherwise set the
//     correct child pointer of parent to replacement.
//
// Time:  O(h) — a search descent plus a successor descent.
// Space: O(1) — only a handful of pointers; no recursion/stack.
func iterativeDelete(root *TreeNode, key int) *TreeNode {
	// Step 1: locate the target node and remember its parent.
	var parent *TreeNode
	cur := root
	for cur != nil && cur.Val != key {
		parent = cur
		if key < cur.Val {
			cur = cur.Left // go left for smaller keys
		} else {
			cur = cur.Right // go right for larger keys
		}
	}
	if cur == nil {
		return root // key absent → tree unchanged
	}

	// Step 2: compute the subtree that will replace cur.
	var replacement *TreeNode
	if cur.Left == nil {
		replacement = cur.Right // 0/1 child: right child (or nil) takes over
	} else if cur.Right == nil {
		replacement = cur.Left // only a left child
	} else {
		// Two children: detach inorder successor (leftmost of right subtree).
		succParent := cur
		succ := cur.Right
		for succ.Left != nil {
			succParent = succ
			succ = succ.Left
		}
		if succParent != cur {
			// Successor is deeper: unlink it, then let it adopt cur.Right.
			succParent.Left = succ.Right // successor's right child fills its slot
			succ.Right = cur.Right       // successor adopts the whole right subtree
		}
		succ.Left = cur.Left // successor adopts cur's left subtree
		replacement = succ   // successor now stands in for cur
	}

	// Step 3: attach the replacement where cur used to hang.
	if parent == nil {
		return replacement // cur was the root
	}
	if parent.Left == cur {
		parent.Left = replacement // cur was a left child
	} else {
		parent.Right = replacement // cur was a right child
	}
	return root
}

// levelOrder renders a tree as a LeetCode-style level-order list (with "null"
// for missing children, trailing nulls trimmed) for verifying results.
func levelOrder(root *TreeNode) string {
	if root == nil {
		return "[]"
	}
	var out []string
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		if node == nil {
			out = append(out, "null")
			continue
		}
		out = append(out, strconv.Itoa(node.Val))
		queue = append(queue, node.Left, node.Right)
	}
	for len(out) > 0 && out[len(out)-1] == "null" {
		out = out[:len(out)-1]
	}
	return "[" + strings.Join(out, ",") + "]"
}

// buildBST inserts values (in the given order) into a BST to construct examples.
func buildBST(vals ...int) *TreeNode {
	var root *TreeNode
	var insert func(node *TreeNode, v int) *TreeNode
	insert = func(node *TreeNode, v int) *TreeNode {
		if node == nil {
			return &TreeNode{Val: v}
		}
		if v < node.Val {
			node.Left = insert(node.Left, v)
		} else {
			node.Right = insert(node.Right, v)
		}
		return node
	}
	for _, v := range vals {
		root = insert(root, v)
	}
	return root
}

func main() {
	// Example 1: [5,3,6,2,4,null,7], key = 3 → [5,4,6,2,null,null,7].
	// Insert order 5,3,6,2,4,7 reproduces that exact shape.
	fmt.Println("=== Approach 1: Recursive Delete (Inorder Successor) ===")
	fmt.Printf("delete key=3 from [5,3,6,2,4,null,7] → %s  expected [5,4,6,2,null,null,7]\n",
		levelOrder(recursiveDelete(buildBST(5, 3, 6, 2, 4, 7), 3)))
	fmt.Printf("delete key=0 from [5,3,6,2,4,null,7] → %s  expected [5,3,6,2,4,null,7]\n",
		levelOrder(recursiveDelete(buildBST(5, 3, 6, 2, 4, 7), 0))) // key absent
	fmt.Printf("delete key=0 from []                 → %s  expected []\n",
		levelOrder(recursiveDelete(nil, 0)))

	fmt.Println("=== Approach 2: Iterative Delete (Parent Pointer) ===")
	fmt.Printf("delete key=3 from [5,3,6,2,4,null,7] → %s  expected [5,4,6,2,null,null,7]\n",
		levelOrder(iterativeDelete(buildBST(5, 3, 6, 2, 4, 7), 3)))
	fmt.Printf("delete key=0 from [5,3,6,2,4,null,7] → %s  expected [5,3,6,2,4,null,7]\n",
		levelOrder(iterativeDelete(buildBST(5, 3, 6, 2, 4, 7), 0)))
	fmt.Printf("delete key=0 from []                 → %s  expected []\n",
		levelOrder(iterativeDelete(nil, 0)))
}
