package main

import (
	"fmt"
	"strings"
)

// TreeNode is a binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// ── Approach 1: Explicit Stack (Brute Force) ─────────────────────────────────
//
// bruteForceStack solves Binary Tree Upside Down using an explicit stack of
// the left spine.
//
// Intuition:
//
//	The guarantee "every right node has a sibling and no children" means the
//	whole tree hangs off the LEFT SPINE (root → root.Left → root.Left.Left…).
//	Right children are only decorative leaves attached to spine nodes.
//	If we collect the spine into a stack, the new root is the deepest (last
//	pushed) spine node, and rebuilding is just popping the stack and pointing
//	each spine node back at its parent.
//
// Algorithm:
//  1. Walk down the left spine, pushing every node onto a stack.
//  2. The top of the stack (leftmost leaf) is the new root.
//  3. Pop pairs (curr, parent): curr.Left = parent.Right, curr.Right = parent.
//  4. The original root ends up as a leaf: clear its children.
//
// Time:  O(n) — every node touched a constant number of times.
// Space: O(n) — the stack stores the whole left spine (up to n nodes).
func bruteForceStack(root *TreeNode) *TreeNode {
	if root == nil {
		return nil // empty tree stays empty
	}
	// step 1: push the entire left spine onto a stack
	stack := []*TreeNode{}
	for node := root; node != nil; node = node.Left {
		stack = append(stack, node)
	}
	// step 2: the deepest-left node becomes the new root
	newRoot := stack[len(stack)-1]
	// step 3: pop from the top; each spine node adopts its parent
	for i := len(stack) - 1; i > 0; i-- {
		curr := stack[i]         // deeper spine node (already the new "upper" node)
		parent := stack[i-1]     // its original parent
		curr.Left = parent.Right // parent's right sibling becomes new left child
		curr.Right = parent      // parent itself becomes new right child
	}
	// step 4: the original root is now a leaf — sever its old pointers
	root.Left, root.Right = nil, nil
	return newRoot
}

// ── Approach 2: Recursion (Bottom-Up) ────────────────────────────────────────
//
// recursion solves Binary Tree Upside Down with post-order recursion.
//
// Intuition:
//
//	Recurse all the way down the left spine first. The deepest-left node is
//	the answer's root. On the way back up, each frame rewires its own little
//	triangle: left child ← becomes parent of → (right sibling, root).
//
// Algorithm:
//  1. Base case: nil node or no left child → this node is the new root.
//  2. newRoot = recursion(root.Left)  (flip everything below first).
//  3. root.Left.Left  = root.Right  (rule 3: right child → new left child).
//     root.Left.Right = root        (rule 2: root → new right child).
//  4. Clear root's own Left/Right (it is now a leaf) and bubble newRoot up.
//
// Time:  O(n) — each node visited once.
// Space: O(n) — recursion depth equals the left-spine length (O(h)).
func recursion(root *TreeNode) *TreeNode {
	// base case: empty tree, or a node with no left child is the new root
	if root == nil || root.Left == nil {
		return root
	}
	left := root.Left   // will become the parent of root after flipping
	right := root.Right // right sibling, will become left's new left child
	// flip the subtree hanging below left first; its answer is the global answer
	newRoot := recursion(left)
	left.Left = right // original right child becomes the new left child
	left.Right = root // original root becomes the new right child
	// root is now a leaf of the flipped tree — remove stale pointers
	root.Left, root.Right = nil, nil
	return newRoot
}

// ── Approach 3: Iterative Pointer Rewiring (Optimal) ─────────────────────────
//
// iterative solves Binary Tree Upside Down like reversing a linked list.
//
// Intuition:
//
//	The left spine is effectively a singly linked list (Left pointers). Flip
//	it exactly like linked-list reversal, carrying one extra piece of state:
//	the PREVIOUS level's right sibling, which must become the current node's
//	new left child. No stack, no recursion → O(1) space.
//
// Algorithm:
//  1. prev = nil (parent already processed), prevRight = nil (its sibling).
//  2. For each spine node curr:
//     a. next = curr.Left            (save the walk direction)
//     b. curr.Left  = prevRight      (parent's sibling → new left child)
//     c. prevRight  = curr.Right     (stash sibling before overwriting)
//     d. curr.Right = prev           (parent → new right child)
//     e. prev = curr; curr = next    (advance down the spine)
//  3. When curr is nil, prev is the leftmost node = new root.
//
// Time:  O(n) — single pass down the left spine.
// Space: O(1) — only three extra pointers.
func iterative(root *TreeNode) *TreeNode {
	var prev, prevRight *TreeNode // processed parent and its right sibling
	curr := root
	for curr != nil {
		next := curr.Left      // remember where to go before rewiring
		curr.Left = prevRight  // rule 3: parent's right child → my left child
		prevRight = curr.Right // stash my sibling for the next iteration
		curr.Right = prev      // rule 2: my parent → my right child
		prev = curr            // I am now the processed parent
		curr = next            // continue down the original left spine
	}
	return prev // last processed node = original leftmost = new root
}

// serialize renders a tree in LeetCode level-order form, e.g. [4,5,2,null,null,3,1].
func serialize(root *TreeNode) string {
	if root == nil {
		return "[]"
	}
	vals := []string{}
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		if n == nil {
			vals = append(vals, "null") // placeholder to keep positions aligned
			continue
		}
		vals = append(vals, fmt.Sprintf("%d", n.Val))
		queue = append(queue, n.Left, n.Right) // children (possibly nil) keep shape
	}
	// LeetCode trims trailing nulls
	for len(vals) > 0 && vals[len(vals)-1] == "null" {
		vals = vals[:len(vals)-1]
	}
	return "[" + strings.Join(vals, ",") + "]"
}

// buildExample1 returns a fresh copy of [1,2,3,4,5] (trees are mutated in place).
func buildExample1() *TreeNode {
	return &TreeNode{1,
		&TreeNode{2, &TreeNode{4, nil, nil}, &TreeNode{5, nil, nil}},
		&TreeNode{3, nil, nil}}
}

func main() {
	fmt.Println("=== Approach 1: Explicit Stack (Brute Force) ===")
	fmt.Println(serialize(bruteForceStack(buildExample1())))   // expected [4,5,2,null,null,3,1]
	fmt.Println(serialize(bruteForceStack(nil)))               // expected []
	fmt.Println(serialize(bruteForceStack(&TreeNode{Val: 1}))) // expected [1]

	fmt.Println("=== Approach 2: Recursion (Bottom-Up) ===")
	fmt.Println(serialize(recursion(buildExample1())))   // expected [4,5,2,null,null,3,1]
	fmt.Println(serialize(recursion(nil)))               // expected []
	fmt.Println(serialize(recursion(&TreeNode{Val: 1}))) // expected [1]

	fmt.Println("=== Approach 3: Iterative Pointer Rewiring (Optimal) ===")
	fmt.Println(serialize(iterative(buildExample1())))   // expected [4,5,2,null,null,3,1]
	fmt.Println(serialize(iterative(nil)))               // expected []
	fmt.Println(serialize(iterative(&TreeNode{Val: 1}))) // expected [1]
}
