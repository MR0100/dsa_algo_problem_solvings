package main

import (
	"fmt"
	"strings"
)

// TreeNode is the standard LeetCode binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// null marks missing children in level-order tree literals (like LeetCode's null).
const null = -1 << 60

// buildTree constructs a binary tree from a LeetCode-style level-order slice.
func buildTree(vals []int) *TreeNode {
	if len(vals) == 0 || vals[0] == null {
		return nil
	}
	root := &TreeNode{Val: vals[0]}
	queue := []*TreeNode{root} // nodes still waiting to receive children
	i := 1                     // next value to consume from the slice
	for len(queue) > 0 && i < len(vals) {
		node := queue[0]
		queue = queue[1:] // pop the front node
		if i < len(vals) {
			if vals[i] != null { // non-null → attach a left child
				node.Left = &TreeNode{Val: vals[i]}
				queue = append(queue, node.Left)
			}
			i++
		}
		if i < len(vals) {
			if vals[i] != null { // non-null → attach a right child
				node.Right = &TreeNode{Val: vals[i]}
				queue = append(queue, node.Right)
			}
			i++
		}
	}
	return root
}

// bstIterator lets main() drive every implementation through one test harness.
type bstIterator interface {
	Next() int
	HasNext() bool
}

// ── Approach 1: Brute Force (Flatten to Array) ───────────────────────────────
//
// FlattenIterator solves BST Iterator by precomputing the entire inorder
// traversal into a slice at construction time, then just moving an index.
//
// Intuition:
//
//	An inorder traversal of a BST yields its values in sorted order. If we
//	run the whole traversal up front and store it, Next() and HasNext()
//	degenerate to array indexing — trivially O(1). The cost: O(n) memory
//	even if the caller only ever asks for the first two elements, which is
//	exactly what the follow-up's O(h) requirement forbids.
//
// Algorithm:
//
//	Constructor: recursive inorder traversal appending every value to a
//	slice; set index = 0.
//	Next: return vals[idx], then idx++.
//	HasNext: idx < len(vals).
//
// Time:  Constructor O(n); Next O(1); HasNext O(1).
// Space: O(n) — the flattened value array (plus O(h) recursion during build).
type FlattenIterator struct {
	vals []int // full inorder sequence, i.e. sorted BST values
	idx  int   // position of the next unreturned value
}

func NewFlattenIterator(root *TreeNode) *FlattenIterator {
	it := &FlattenIterator{}
	var inorder func(*TreeNode)
	inorder = func(n *TreeNode) {
		if n == nil {
			return
		}
		inorder(n.Left)                  // everything smaller first
		it.vals = append(it.vals, n.Val) // then the node itself
		inorder(n.Right)                 // then everything larger
	}
	inorder(root)
	return it
}

func (it *FlattenIterator) Next() int {
	v := it.vals[it.idx] // problem guarantees Next() is only called when valid
	it.idx++             // advance the pointer past the returned value
	return v
}

func (it *FlattenIterator) HasNext() bool { return it.idx < len(it.vals) }

// ── Approach 2: Controlled Recursion with a Stack (Optimal) ──────────────────
//
// StackIterator solves BST Iterator by pausing an inorder traversal: an
// explicit stack holds the path of nodes whose values are still owed.
//
// Intuition:
//
//	Recursive inorder = "go left as deep as possible, visit, then recurse
//	right". Replace the call stack with our own stack and we can FREEZE the
//	traversal between Next() calls. Invariant: the stack holds exactly the
//	nodes whose left subtrees are fully consumed but which themselves are
//	unvisited — so the top of the stack is always the next-smallest value.
//	Popping a node may unlock its right subtree, whose left spine we push
//	immediately to restore the invariant.
//
// Algorithm:
//
//	Constructor: push the left spine of root (root, root.Left, ...).
//	Next: pop the top node (the answer); push the left spine of its right
//	      child; return the popped value.
//	HasNext: stack non-empty.
//
// Time:  Next amortized O(1) — each node is pushed and popped exactly once
//
//	over the whole iteration (single Next() worst case O(h)); HasNext O(1).
//
// Space: O(h) — the stack never holds more than one root-to-leaf path.
type StackIterator struct {
	stack []*TreeNode // partially-frozen inorder traversal; top = next value
}

func NewStackIterator(root *TreeNode) *StackIterator {
	it := &StackIterator{}
	it.pushLeftSpine(root) // prime the stack so the minimum is on top
	return it
}

// pushLeftSpine pushes node and all its left descendants onto the stack.
func (it *StackIterator) pushLeftSpine(node *TreeNode) {
	for node != nil {
		it.stack = append(it.stack, node) // owe this node's value later
		node = node.Left                  // its left subtree comes first
	}
}

func (it *StackIterator) Next() int {
	top := it.stack[len(it.stack)-1]      // smallest unvisited node
	it.stack = it.stack[:len(it.stack)-1] // pop it — we are visiting it now
	it.pushLeftSpine(top.Right)           // its right subtree is next in line
	return top.Val
}

func (it *StackIterator) HasNext() bool { return len(it.stack) > 0 }

// ── Approach 3: Morris Threading (O(1) Extra Space) ──────────────────────────
//
// MorrisIterator solves BST Iterator with Morris inorder traversal: instead
// of a stack, it temporarily rewires right pointers of inorder predecessors
// as "threads" back to their successors.
//
// Intuition:
//
//	The stack exists only so we can climb back up after finishing a left
//	subtree. Morris traversal stores that return path IN the tree: before
//	descending into a left subtree, point the rightmost node of that subtree
//	(the current node's inorder predecessor) back at the current node. When
//	iteration later walks off that predecessor's right pointer, it lands
//	exactly where it must resume — and the second visit removes the thread,
//	restoring the tree. State between Next() calls is a single pointer.
//
// Algorithm (one Next() call):
//  1. Loop on cur:
//  2. If cur has no left child → its value is next; move cur = cur.Right
//     (which may be a thread) and return the value.
//  3. Otherwise find pred = rightmost node of cur.Left, stopping early if
//     pred.Right already threads to cur.
//  4. If pred.Right == nil: create thread pred.Right = cur; descend
//     cur = cur.Left.
//  5. Else (thread found — left subtree done): remove the thread, return
//     cur's value, and move cur = cur.Right.
//
// Time:  Next amortized O(1) — every tree edge is walked at most twice
//
//	(once creating a thread, once removing it) across the full iteration.
//
// Space: O(1) — one pointer of state; the tree itself stores the traversal
//
//	(mutated during iteration, fully restored once iteration completes).
type MorrisIterator struct {
	cur *TreeNode // where the paused traversal will resume
}

func NewMorrisIterator(root *TreeNode) *MorrisIterator {
	return &MorrisIterator{cur: root}
}

func (it *MorrisIterator) Next() int {
	for it.cur != nil {
		if it.cur.Left == nil {
			// No left subtree → cur itself is the next inorder value.
			v := it.cur.Val
			it.cur = it.cur.Right // real child or thread to the successor
			return v
		}
		// Find cur's inorder predecessor: rightmost node in the left subtree.
		pred := it.cur.Left
		for pred.Right != nil && pred.Right != it.cur {
			pred = pred.Right
		}
		if pred.Right == nil {
			pred.Right = it.cur  // lay a thread so we can return to cur later
			it.cur = it.cur.Left // now safe to descend into the left subtree
		} else {
			pred.Right = nil      // second arrival: left subtree finished — unthread
			v := it.cur.Val       // cur is the next inorder value
			it.cur = it.cur.Right // continue with the right subtree
			return v
		}
	}
	return null // unreachable: problem guarantees Next() calls are valid
}

func (it *MorrisIterator) HasNext() bool { return it.cur != nil }

// runExample drives one implementation through the official LeetCode example:
//
//	ops:  ["BSTIterator","next","next","hasNext","next","hasNext","next","hasNext","next","hasNext"]
//	args: [[[7,3,15,null,null,9,20]],[],[],[],[],[],[],[],[],[]]
//
// expected output: [null,3,7,true,9,true,15,true,20,false]
func runExample(makeIter func(*TreeNode) bstIterator) {
	// Fresh tree per implementation (Morris mutates it mid-iteration).
	root := buildTree([]int{7, 3, 15, null, null, 9, 20})
	it := makeIter(root)
	out := []string{"null"}                     // constructor produces null
	out = append(out, fmt.Sprint(it.Next()))    // 3
	out = append(out, fmt.Sprint(it.Next()))    // 7
	out = append(out, fmt.Sprint(it.HasNext())) // true
	out = append(out, fmt.Sprint(it.Next()))    // 9
	out = append(out, fmt.Sprint(it.HasNext())) // true
	out = append(out, fmt.Sprint(it.Next()))    // 15
	out = append(out, fmt.Sprint(it.HasNext())) // true
	out = append(out, fmt.Sprint(it.Next()))    // 20
	out = append(out, fmt.Sprint(it.HasNext())) // false
	fmt.Printf("[%s]\n", strings.Join(out, ","))
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Flatten to Array) ===")
	runExample(func(r *TreeNode) bstIterator { return NewFlattenIterator(r) }) // [null,3,7,true,9,true,15,true,20,false]

	fmt.Println("=== Approach 2: Controlled Recursion with a Stack (Optimal) ===")
	runExample(func(r *TreeNode) bstIterator { return NewStackIterator(r) }) // [null,3,7,true,9,true,15,true,20,false]

	fmt.Println("=== Approach 3: Morris Threading (O(1) Extra Space) ===")
	runExample(func(r *TreeNode) bstIterator { return NewMorrisIterator(r) }) // [null,3,7,true,9,true,15,true,20,false]
}
