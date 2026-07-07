package main

import "fmt"

// TreeNode is the standard LeetCode binary-tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// ── Approach 1: In-Order Traversal to Full List ──────────────────────────────
//
// inorderFullList collects every value in sorted order via a complete in-order
// traversal, then indexes the (k-1)-th element.
//
// Intuition:
//
//	An in-order traversal of a BST visits nodes in ascending value order. So
//	the k-th smallest is simply the k-th element produced by that traversal.
//	The simplest form materializes the whole sorted list and indexes into it.
//
// Algorithm:
//  1. Recursively in-order traverse (left, node, right), appending each value.
//  2. Return list[k-1] (1-indexed k → 0-indexed slice).
//
// Time:  O(n) — visits every node once.
// Space: O(n) — the full value list plus O(h) recursion stack.
func inorderFullList(root *TreeNode, k int) int {
	var vals []int
	var inorder func(node *TreeNode)
	inorder = func(node *TreeNode) {
		if node == nil {
			return
		}
		inorder(node.Left)            // all smaller values first
		vals = append(vals, node.Val) // then this node (ascending order)
		inorder(node.Right)           // then all larger values
	}
	inorder(root)
	return vals[k-1] // k is 1-indexed; slice is 0-indexed
}

// ── Approach 2: In-Order with Early Stop (Counter) ───────────────────────────
//
// inorderEarlyStop performs an in-order traversal but stops as soon as it has
// seen k nodes, avoiding traversal of the rest of the tree.
//
// Intuition:
//
//	We don't need the whole sorted list — only its k-th element. Count nodes as
//	the in-order walk emits them; when the count reaches k, that node's value is
//	the answer, and everything to its right is irrelevant. This trims work to
//	O(h + k) instead of O(n).
//
// Algorithm:
//  1. Keep a running count and an answer holder.
//  2. In-order recurse; on visiting a node, increment count. If count == k,
//     record the value and short-circuit further recursion.
//  3. Return the recorded value.
//
// Time:  O(h + k) — descend to the smallest (O(h)) then emit k nodes.
// Space: O(h) — recursion stack depth.
func inorderEarlyStop(root *TreeNode, k int) int {
	count := 0 // how many nodes emitted so far
	ans := 0   // the k-th smallest value once found
	found := false

	var inorder func(node *TreeNode)
	inorder = func(node *TreeNode) {
		if node == nil || found {
			return // stop descending once the answer is fixed
		}
		inorder(node.Left) // smaller values first
		if found {
			return // answer found in the left subtree — unwind
		}
		count++ // this node is the next value in ascending order
		if count == k {
			ans = node.Val // k-th smallest reached
			found = true
			return
		}
		inorder(node.Right) // otherwise continue into larger values
	}
	inorder(root)
	return ans
}

// ── Approach 3: Iterative In-Order (Stack, Optimal) ──────────────────────────
//
// iterativeInorder walks the BST in order using an explicit stack, decrementing
// k at each emitted node and returning the moment k hits zero.
//
// Intuition:
//
//	The iterative in-order pattern pushes left spines onto a stack, then pops to
//	emit nodes in ascending order. Emitting nodes one at a time lets us stop
//	exactly at the k-th without recursion and without building any list —
//	constant-per-node control flow that naturally supports the follow-up about
//	frequent modifications (you can pause/resume the iterator).
//
// Algorithm:
//  1. stack empty, curr = root.
//  2. Loop: push curr and all its left descendants (go as left as possible).
//  3. Pop a node — it is the next in ascending order. Decrement k; if k == 0,
//     return its value.
//  4. Move curr to the popped node's right child and repeat.
//
// Time:  O(h + k) — O(h) to reach the smallest, then k pops.
// Space: O(h) — the stack holds at most one root-to-leaf path.
func iterativeInorder(root *TreeNode, k int) int {
	stack := []*TreeNode{}
	curr := root
	for curr != nil || len(stack) > 0 {
		for curr != nil { // dive to the leftmost unvisited node
			stack = append(stack, curr)
			curr = curr.Left
		}
		// Pop: this node is the next smallest not yet emitted.
		curr = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		k--
		if k == 0 {
			return curr.Val // exactly the k-th smallest
		}
		curr = curr.Right // explore the right subtree (larger values)
	}
	return -1 // unreachable given valid input (1 ≤ k ≤ n)
}

// buildBST inserts values into a BST in the given order (helper for examples).
func buildBST(vals []int) *TreeNode {
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
	// Example 1: tree [3,1,4,null,2], k = 1 → 1.
	// Build the exact shape via insertion order 3,1,4,2.
	t1 := buildBST([]int{3, 1, 4, 2})
	// Example 2: tree [5,3,6,2,4,null,null,1], k = 3 → 3.
	t2 := buildBST([]int{5, 3, 6, 2, 4, 1})

	fmt.Println("=== Approach 1: In-Order Traversal to Full List ===")
	fmt.Println(inorderFullList(t1, 1)) // 1
	fmt.Println(inorderFullList(t2, 3)) // 3

	fmt.Println("=== Approach 2: In-Order with Early Stop (Counter) ===")
	fmt.Println(inorderEarlyStop(t1, 1)) // 1
	fmt.Println(inorderEarlyStop(t2, 3)) // 3

	fmt.Println("=== Approach 3: Iterative In-Order (Stack, Optimal) ===")
	fmt.Println(iterativeInorder(t1, 1)) // 1
	fmt.Println(iterativeInorder(t2, 3)) // 3
}
