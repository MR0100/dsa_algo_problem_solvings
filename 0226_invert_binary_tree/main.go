package main

import "fmt"

// TreeNode is the standard LeetCode binary-tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// ── Approach 1: Recursive DFS (Swap Children) ────────────────────────────────
//
// invertRecursive inverts a binary tree by swapping the left and right child of
// every node, recursing top-down.
//
// Intuition:
//
//	Inverting a tree is its own mirror image: at each node the left subtree and
//	right subtree simply trade places. Do that swap at every node and the whole
//	tree is mirrored. Recursion expresses this cleanly — invert this node's
//	children, then let recursion invert everything below.
//
// Algorithm:
//  1. If the node is nil, return nil (empty subtree is its own mirror).
//  2. Recursively invert the left subtree and the right subtree.
//  3. Swap the (now-inverted) left and right pointers.
//  4. Return the node.
//
// Time:  O(n) — every node is visited exactly once.
// Space: O(h) — recursion stack depth equals tree height (O(n) worst, O(log n) balanced).
func invertRecursive(root *TreeNode) *TreeNode {
	if root == nil { // empty subtree — nothing to mirror
		return nil
	}
	// Invert both sides first, then swap the returned (inverted) subtrees.
	left := invertRecursive(root.Left)   // fully-inverted left subtree
	right := invertRecursive(root.Right) // fully-inverted right subtree
	root.Left = right                    // right subtree now hangs on the left
	root.Right = left                    // left subtree now hangs on the right
	return root
}

// ── Approach 2: Iterative BFS (Queue) ────────────────────────────────────────
//
// invertBFS inverts a binary tree level by level using an explicit queue,
// swapping each dequeued node's children.
//
// Intuition:
//
//	The swap at each node is independent of every other node's swap, so the
//	order we visit nodes does not matter — only that we visit them all. A queue
//	gives a level-order sweep with no recursion, sidestepping deep-stack risk.
//
// Algorithm:
//  1. If root is nil, return nil.
//  2. Push root into a queue.
//  3. While the queue is non-empty: pop a node, swap its two children, and
//     enqueue any non-nil children to process later.
//  4. Return root.
//
// Time:  O(n) — each node enters and leaves the queue once.
// Space: O(w) — queue holds at most one tree level (O(n) worst case).
func invertBFS(root *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}
	queue := []*TreeNode{root} // FIFO of nodes still needing their swap
	for len(queue) > 0 {
		node := queue[0]  // dequeue the front node
		queue = queue[1:] // pop it off the queue
		// Swap this node's children — the core mirror operation.
		node.Left, node.Right = node.Right, node.Left
		if node.Left != nil { // enqueue children (post-swap) for later processing
			queue = append(queue, node.Left)
		}
		if node.Right != nil {
			queue = append(queue, node.Right)
		}
	}
	return root
}

// ── Approach 3: Iterative DFS (Stack) ────────────────────────────────────────
//
// invertDFSStack inverts a binary tree using an explicit stack instead of the
// call stack, swapping children as each node is popped.
//
// Intuition:
//
//	Same independent-swap observation as BFS, but a LIFO stack reproduces the
//	depth-first order of the recursive version without recursion — useful when
//	you want to avoid Go's growing goroutine stack or hit recursion limits.
//
// Algorithm:
//  1. If root is nil, return nil.
//  2. Push root onto a stack.
//  3. While the stack is non-empty: pop a node, swap its children, and push any
//     non-nil children.
//  4. Return root.
//
// Time:  O(n) — every node pushed and popped once.
// Space: O(h) — stack holds at most one root-to-leaf path plus siblings.
func invertDFSStack(root *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}
	stack := []*TreeNode{root} // LIFO of nodes awaiting their swap
	for len(stack) > 0 {
		node := stack[len(stack)-1]  // peek the top
		stack = stack[:len(stack)-1] // pop it
		// Swap children — the mirror operation.
		node.Left, node.Right = node.Right, node.Left
		if node.Left != nil { // push children to keep going deeper
			stack = append(stack, node.Left)
		}
		if node.Right != nil {
			stack = append(stack, node.Right)
		}
	}
	return root
}

// buildTree constructs a tree from a level-order slice using nil for missing
// nodes (LeetCode's array representation), so examples read like the site.
func buildTree(vals []interface{}) *TreeNode {
	if len(vals) == 0 || vals[0] == nil {
		return nil
	}
	root := &TreeNode{Val: vals[0].(int)}
	queue := []*TreeNode{root} // nodes still needing children attached
	i := 1                     // index into vals for the next child value
	for len(queue) > 0 && i < len(vals) {
		node := queue[0]
		queue = queue[1:]
		if i < len(vals) && vals[i] != nil { // attach left child if present
			node.Left = &TreeNode{Val: vals[i].(int)}
			queue = append(queue, node.Left)
		}
		i++
		if i < len(vals) && vals[i] != nil { // attach right child if present
			node.Right = &TreeNode{Val: vals[i].(int)}
			queue = append(queue, node.Right)
		}
		i++
	}
	return root
}

// serialize returns the level-order slice of a tree, trimming trailing nils so
// the output matches LeetCode's compact array form.
func serialize(root *TreeNode) []interface{} {
	if root == nil {
		return []interface{}{}
	}
	var out []interface{}
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		if node == nil {
			out = append(out, nil)
			continue
		}
		out = append(out, node.Val)
		queue = append(queue, node.Left, node.Right)
	}
	// Trim trailing nils to match LeetCode's compact representation.
	for len(out) > 0 && out[len(out)-1] == nil {
		out = out[:len(out)-1]
	}
	return out
}

func main() {
	fmt.Println("=== Approach 1: Recursive DFS (Swap Children) ===")
	fmt.Println(serialize(invertRecursive(buildTree([]interface{}{4, 2, 7, 1, 3, 6, 9})))) // [4 7 2 9 6 3 1]
	fmt.Println(serialize(invertRecursive(buildTree([]interface{}{2, 1, 3}))))             // [2 3 1]
	fmt.Println(serialize(invertRecursive(buildTree([]interface{}{}))))                    // []

	fmt.Println("=== Approach 2: Iterative BFS (Queue) ===")
	fmt.Println(serialize(invertBFS(buildTree([]interface{}{4, 2, 7, 1, 3, 6, 9})))) // [4 7 2 9 6 3 1]
	fmt.Println(serialize(invertBFS(buildTree([]interface{}{2, 1, 3}))))             // [2 3 1]
	fmt.Println(serialize(invertBFS(buildTree([]interface{}{}))))                    // []

	fmt.Println("=== Approach 3: Iterative DFS (Stack) ===")
	fmt.Println(serialize(invertDFSStack(buildTree([]interface{}{4, 2, 7, 1, 3, 6, 9})))) // [4 7 2 9 6 3 1]
	fmt.Println(serialize(invertDFSStack(buildTree([]interface{}{2, 1, 3}))))             // [2 3 1]
	fmt.Println(serialize(invertDFSStack(buildTree([]interface{}{}))))                    // []
}
