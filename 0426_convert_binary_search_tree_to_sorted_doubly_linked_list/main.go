package main

import (
	"fmt"
	"strings"
)

// Node is the BST node reused as a doubly linked list node (LeetCode #426):
// after conversion, Left means "predecessor" and Right means "successor".
type Node struct {
	Val   int
	Left  *Node
	Right *Node
}

// ── Approach 1: In-order into a Slice, then Link (Brute Force) ────────────────
//
// inorderSlice solves Convert BST to Sorted DLL by first collecting nodes in
// sorted order into a slice, then wiring the doubly linked list in a second
// pass.
//
// Intuition:
//
//	An in-order traversal of a BST yields nodes in ascending value order — which
//	is exactly the order the circular list must have. So flatten to a slice
//	first (correctness is obvious), then stitch consecutive nodes with
//	Left/Right and finally close the circle between the last and first.
//
// Algorithm:
//  1. In-order traverse, appending every node pointer to `nodes`.
//  2. If empty, return nil.
//  3. For each adjacent pair, set a.Right = b and b.Left = a.
//  4. Close the ring: first.Left = last, last.Right = first.
//  5. Return nodes[0] (the smallest).
//
// Time:  O(n) — one traversal plus one linking pass.
// Space: O(n) — the slice of node pointers (plus O(h) recursion stack).
func inorderSlice(root *Node) *Node {
	if root == nil {
		return nil
	}
	nodes := []*Node{} // node pointers in ascending value order
	var inorder func(*Node)
	inorder = func(n *Node) {
		if n == nil {
			return
		}
		inorder(n.Left)          // left subtree: smaller values first
		nodes = append(nodes, n) // visit: record this node
		inorder(n.Right)         // right subtree: larger values
	}
	inorder(root)

	for i := 0; i < len(nodes); i++ {
		nodes[i].Right = nodes[(i+1)%len(nodes)]           // successor (wraps to head at the end)
		nodes[i].Left = nodes[(i-1+len(nodes))%len(nodes)] // predecessor (wraps to tail at start)
	}
	return nodes[0] // smallest element = head of the circular DLL
}

// ── Approach 2: In-order with a Running Previous Pointer (Optimal) ────────────
//
// inorderInPlace solves the problem in a single in-order pass, linking each
// visited node to the previously visited one on the fly — O(1) auxiliary space
// beyond the recursion stack.
//
// Intuition:
//
//	During an in-order walk, nodes are visited in the final sorted order, so at
//	the moment we visit node `cur`, the previously visited node `prev` is exactly
//	its predecessor. Wire prev.Right = cur and cur.Left = prev immediately. Track
//	the very first visited node as `head`; when the traversal ends, `prev` is the
//	tail, and we close the circle head↔tail.
//
// Algorithm:
//  1. Keep two pointers: head (first visited) and prev (last visited).
//  2. In-order: on visiting cur, if prev != nil link prev.Right = cur and
//     cur.Left = prev; else cur is the smallest → head = cur. Then prev = cur.
//  3. After traversal, if head != nil close the ring: head.Left = prev,
//     prev.Right = head.
//  4. Return head.
//
// Time:  O(n) — each node visited exactly once.
// Space: O(h) — recursion stack only; no auxiliary array.
func inorderInPlace(root *Node) *Node {
	if root == nil {
		return nil
	}
	var head, prev *Node // head = smallest node; prev = last node linked so far
	var inorder func(*Node)
	inorder = func(cur *Node) {
		if cur == nil {
			return
		}
		inorder(cur.Left) // recurse left first (smaller values)
		if prev != nil {
			prev.Right = cur // previous node's successor is this node
			cur.Left = prev  // this node's predecessor is the previous node
		} else {
			head = cur // no predecessor yet ⇒ this is the smallest node
		}
		prev = cur         // this node becomes the "previous" for the next visit
		inorder(cur.Right) // then recurse right (larger values)
	}
	inorder(root)

	// Close the circle: smallest's predecessor is the largest and vice versa.
	head.Left = prev
	prev.Right = head
	return head
}

// ── Approach 3: Iterative In-order with an Explicit Stack ────────
//
// iterativeInorder does the in-order linking iteratively with an explicit stack,
// avoiding recursion entirely while keeping the running-previous logic.
//
// Intuition:
//
//	The only thing recursion gave us was the in-order visit sequence. An explicit
//	stack reproduces that sequence iteratively: push all left descendants, pop to
//	visit, then move to the right child. The prev/head linking logic is identical
//	to Approach 2, just driven by a stack instead of the call stack.
//
// Algorithm:
//  1. Use a stack; cur = root.
//  2. Loop while cur != nil or stack non-empty:
//     a. Push every left descendant of cur.
//     b. Pop a node = the next in-order node; link it to prev (or set head).
//     c. Move cur to that node's Right.
//  3. Close the ring and return head.
//
// Time:  O(n) — each node pushed and popped once.
// Space: O(h) — the explicit stack (height of the tree).
func iterativeInorder(root *Node) *Node {
	if root == nil {
		return nil
	}
	var head, prev *Node
	stack := []*Node{}
	cur := root
	for cur != nil || len(stack) > 0 {
		for cur != nil {
			stack = append(stack, cur) // remember the path down the left spine
			cur = cur.Left
		}
		cur = stack[len(stack)-1]    // top = next node in ascending order
		stack = stack[:len(stack)-1] // pop it

		if prev != nil {
			prev.Right = cur // link predecessor → current
			cur.Left = prev
		} else {
			head = cur // first popped node is the minimum
		}
		prev = cur

		cur = cur.Right // explore the right subtree next
	}
	head.Left = prev
	prev.Right = head
	return head
}

// ── helpers for building the BST and printing the circular list ──────────────

// insert adds val into a BST rooted at root and returns the (possibly new) root.
func insert(root *Node, val int) *Node {
	if root == nil {
		return &Node{Val: val}
	}
	if val < root.Val {
		root.Left = insert(root.Left, val)
	} else {
		root.Right = insert(root.Right, val)
	}
	return root
}

// buildBST builds a BST by inserting vals in order.
func buildBST(vals ...int) *Node {
	var root *Node
	for _, v := range vals {
		root = insert(root, v)
	}
	return root
}

// forward walks Right pointers n steps from head, collecting values, to show
// the ascending order of the circular list.
func forward(head *Node) []int {
	out := []int{}
	if head == nil {
		return out
	}
	cur := head
	for {
		out = append(out, cur.Val)
		cur = cur.Right
		if cur == head { // came back around ⇒ full loop done
			break
		}
	}
	return out
}

// backward walks Left pointers from head (i.e. starts at the tail and goes
// down), collecting values, to prove the reverse links and circularity.
func backward(head *Node) []int {
	out := []int{}
	if head == nil {
		return out
	}
	cur := head.Left // tail (largest), since head.Left wraps to the end
	for {
		out = append(out, cur.Val)
		if cur == head {
			break
		}
		cur = cur.Left
	}
	return out
}

func fmtSlice(s []int) string {
	parts := make([]string, len(s))
	for i, v := range s {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return "[" + strings.Join(parts, ",") + "]"
}

func main() {
	// Example 1: root = [4,2,5,1,3] → sorted DLL 1<->2<->3<->4<->5 (circular).
	fmt.Println("=== Approach 1: In-order Slice + Link — Example 1 ===")
	h1 := inorderSlice(buildBST(4, 2, 5, 1, 3))
	fmt.Println("forward: ", fmtSlice(forward(h1)))  // [1,2,3,4,5]
	fmt.Println("backward:", fmtSlice(backward(h1))) // [5,4,3,2,1]

	fmt.Println("=== Approach 2: In-order Running-prev (Optimal) — Example 1 ===")
	h2 := inorderInPlace(buildBST(4, 2, 5, 1, 3))
	fmt.Println("forward: ", fmtSlice(forward(h2)))  // [1,2,3,4,5]
	fmt.Println("backward:", fmtSlice(backward(h2))) // [5,4,3,2,1]

	fmt.Println("=== Approach 3: Iterative In-order (stack) — Example 1 ===")
	h3 := iterativeInorder(buildBST(4, 2, 5, 1, 3))
	fmt.Println("forward: ", fmtSlice(forward(h3)))  // [1,2,3,4,5]
	fmt.Println("backward:", fmtSlice(backward(h3))) // [5,4,3,2,1]

	// Example 2: root = [2,1,3] → head should be node with value 1.
	fmt.Println("=== Example 2: root = [2,1,3] ===")
	h4 := inorderInPlace(buildBST(2, 1, 3))
	fmt.Println("head val:", h4.Val, "forward:", fmtSlice(forward(h4))) // head val: 1 forward: [1,2,3]

	// Example 3: empty tree → nil.
	fmt.Println("=== Example 3: root = [] ===")
	fmt.Println("head:", inorderInPlace(nil)) // head: <nil>
}
