package main

import "fmt"

// TreeNode is a binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// ── Approach 1: Recursive DFS ────────────────────────────────────────────────
//
// postorderRecursive solves Binary Tree Postorder Traversal with recursion.
//
// Intuition:
//
//	Postorder = left, right, root: a node is emitted only after BOTH of its
//	subtrees are fully done. The recursive definition maps 1:1 onto code —
//	recurse left, recurse right, then append the node.
//
// Algorithm:
//  1. If node == nil → return.
//  2. Recurse into node.Left.
//  3. Recurse into node.Right.
//  4. Append node.Val (root last — that's what "post" means).
//
// Time:  O(n) — each node visited exactly once.
// Space: O(h) — recursion stack of height h (O(n) skewed, O(log n) balanced).
func postorderRecursive(root *TreeNode) []int {
	result := []int{} // non-nil so an empty tree prints as [] not nil
	var dfs func(node *TreeNode)
	dfs = func(node *TreeNode) {
		if node == nil {
			return // empty subtree contributes nothing
		}
		dfs(node.Left)                    // finish the LEFT subtree first
		dfs(node.Right)                   // then the RIGHT subtree
		result = append(result, node.Val) // ROOT last
	}
	dfs(root)
	return result
}

// ── Approach 2: Reversed Modified Preorder (Stack + Reverse) ─────────────────
//
// reversedPreorder solves Binary Tree Postorder Traversal by producing a
// root-RIGHT-left preorder with a stack, then reversing it.
//
// Intuition:
//
//	Postorder is left, right, root. Read backwards that is root, right, left —
//	which is just preorder with the children swapped. Preorder is trivial with
//	a stack, so: run the mirrored preorder, reverse the output at the end.
//
// Algorithm:
//  1. If root == nil → return [].
//  2. stack = [root].
//  3. While stack non-empty: pop node, append node.Val,
//     push node.Left, then push node.Right (LEFT first so RIGHT pops first).
//  4. Reverse the result slice in place.
//
// Time:  O(n) — n pushes/pops plus an O(n) reversal.
// Space: O(n) — the stack can hold O(n) nodes (and the reversal is in-place).
func reversedPreorder(root *TreeNode) []int {
	result := []int{}
	if root == nil {
		return result // nothing to traverse
	}
	stack := []*TreeNode{root}
	for len(stack) > 0 {
		node := stack[len(stack)-1]       // peek top
		stack = stack[:len(stack)-1]      // pop it
		result = append(result, node.Val) // building root-right-left order
		if node.Left != nil {
			stack = append(stack, node.Left) // pushed FIRST → popped LAST
		}
		if node.Right != nil {
			stack = append(stack, node.Right) // pushed LAST → popped NEXT (right first)
		}
	}
	// Reverse root-right-left into left-right-root (= postorder).
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return result
}

// ── Approach 3: One Stack with Last-Visited Tracking ─────────────────────────
//
// oneStack solves Binary Tree Postorder Traversal with a single stack and a
// lastVisited pointer — a true streaming postorder (no final reversal).
//
// Intuition:
//
//	In postorder, a node may only be emitted once its right subtree is done.
//	Standing on the stack's top node after walking all the way left, there
//	are two cases: (a) it has an unvisited right child → dive right first;
//	(b) its right child is nil or was JUST emitted → emit the node now.
//	Remembering the last emitted node distinguishes the two cases.
//
// Algorithm:
//  1. curr = root; stack empty; lastVisited = nil.
//  2. While curr != nil or stack non-empty:
//     a. Slide left: push curr, curr = curr.Left, until curr == nil.
//     b. peek = stack top.
//     c. If peek.Right != nil and peek.Right != lastVisited →
//     curr = peek.Right (explore right subtree before emitting peek).
//     d. Else pop, emit peek, lastVisited = peek (leave curr nil so we
//     don't re-descend left).
//
// Time:  O(n) — every node is pushed once, popped once, peeked O(1) times.
// Space: O(h) — the stack never holds more than one root-to-leaf path.
func oneStack(root *TreeNode) []int {
	result := []int{}
	stack := []*TreeNode{}
	var lastVisited *TreeNode // most recently emitted node
	curr := root

	for curr != nil || len(stack) > 0 {
		// Phase a: push the whole left spine of the current subtree.
		for curr != nil {
			stack = append(stack, curr)
			curr = curr.Left
		}
		peek := stack[len(stack)-1] // candidate for emission
		if peek.Right != nil && peek.Right != lastVisited {
			// Right subtree exists and is NOT finished yet → traverse it first.
			curr = peek.Right
		} else {
			// Right subtree is absent or already emitted → safe to emit peek.
			stack = stack[:len(stack)-1] // pop
			result = append(result, peek.Val)
			lastVisited = peek // remember so the parent knows right is done
			// curr stays nil: next loop iteration re-examines the new top.
		}
	}
	return result
}

// ── Approach 4: Morris Traversal (O(1) Space, Optimal) ───────────────────────
//
// morrisPostorder solves Binary Tree Postorder Traversal in O(1) extra space
// using threading plus reverse-emission of left-child right spines.
//
// Intuition:
//
//	Hang the tree under a dummy root (as its LEFT child). Run inorder Morris
//	threading; each time a thread brings us back to a node (its left subtree
//	is exhausted), emit the RIGHT SPINE of its left child in REVERSE order.
//	Summed over all nodes, those reversed spines produce exactly postorder.
//	Reversal is done by pointer-flipping the spine (like reversing a linked
//	list) and flipping it back — no extra memory.
//
// Algorithm:
//  1. dummy = TreeNode{Left: root}; curr = dummy.
//  2. While curr != nil:
//     a. If curr.Left == nil → curr = curr.Right.
//     b. Else find pred = rightmost of curr.Left (stop at thread):
//     - pred.Right == nil  → thread pred.Right = curr; curr = curr.Left.
//     - pred.Right == curr → unthread; emit the chain curr.Left … pred
//     in reverse (flip Right pointers, walk, flip back); curr = curr.Right.
//
// Time:  O(n) — each edge handled a constant number of times (thread,
// unthread, two spine flips).
// Space: O(1) — output slice aside, only pointers; tree fully restored.
func morrisPostorder(root *TreeNode) []int {
	result := []int{}
	dummy := &TreeNode{Left: root} // ensures the rightmost spine of the real
	curr := dummy                  // root is also emitted via the same rule

	for curr != nil {
		if curr.Left == nil {
			curr = curr.Right // nothing to emit here; keep moving
			continue
		}
		// Find inorder predecessor of curr within its left subtree.
		pred := curr.Left
		for pred.Right != nil && pred.Right != curr {
			pred = pred.Right
		}
		if pred.Right == nil {
			pred.Right = curr // first arrival: lay the return thread
			curr = curr.Left  // and descend left
		} else {
			pred.Right = nil                           // second arrival: remove thread
			emitReverseSpine(curr.Left, pred, &result) // emit spine curr.Left…pred backwards
			curr = curr.Right                          // move on past the finished left subtree
		}
	}
	return result
}

// emitReverseSpine appends the right-pointer chain from…to in reverse order,
// restoring the pointers afterwards. O(length) time, O(1) space.
func emitReverseSpine(from, to *TreeNode, result *[]int) {
	reverseRightChain(from) // flip: to → … → from
	// Walk from `to` back down to `from`, emitting values.
	for node := to; ; node = node.Right {
		*result = append(*result, node.Val)
		if node == from {
			break // reached the start of the original spine
		}
	}
	reverseRightChain(to) // flip back: from → … → to (tree restored)
}

// reverseRightChain reverses the linked chain formed by Right pointers
// starting at head (the chain must end in nil, guaranteed after unthreading).
func reverseRightChain(head *TreeNode) {
	var prev *TreeNode
	for node := head; node != nil; {
		next := node.Right // save onward pointer
		node.Right = prev  // flip the link
		prev = node        // advance prev
		node = next        // advance node
	}
}

// ── Helpers ───────────────────────────────────────────────────────────────────

// null marks a missing child in the level-order encodings below.
const null = -1 << 31

// buildTree constructs a binary tree from LeetCode's level-order array form.
func buildTree(levelOrder []int) *TreeNode {
	if len(levelOrder) == 0 || levelOrder[0] == null {
		return nil
	}
	root := &TreeNode{Val: levelOrder[0]}
	queue := []*TreeNode{root} // parents awaiting children
	i := 1
	for len(queue) > 0 && i < len(levelOrder) {
		parent := queue[0]
		queue = queue[1:]
		if i < len(levelOrder) && levelOrder[i] != null {
			parent.Left = &TreeNode{Val: levelOrder[i]}
			queue = append(queue, parent.Left)
		}
		i++ // consume the left-child slot even when it is null
		if i < len(levelOrder) && levelOrder[i] != null {
			parent.Right = &TreeNode{Val: levelOrder[i]}
			queue = append(queue, parent.Right)
		}
		i++ // consume the right-child slot even when it is null
	}
	return root
}

func main() {
	// Official LeetCode examples: (level-order input, expected postorder).
	examples := []struct {
		tree   []int
		expect []int
	}{
		{[]int{1, null, 2, 3}, []int{3, 2, 1}},                                                 // Example 1
		{[]int{1, 2, 3, 4, 5, null, 8, null, null, 6, 7, 9}, []int{4, 6, 7, 5, 2, 9, 8, 3, 1}}, // Example 2
		{[]int{}, []int{}},   // Example 3
		{[]int{1}, []int{1}}, // Example 4
	}

	approaches := []struct {
		name string
		fn   func(*TreeNode) []int
	}{
		{"Approach 1: Recursive DFS", postorderRecursive},
		{"Approach 2: Reversed Modified Preorder", reversedPreorder},
		{"Approach 3: One Stack + Last Visited", oneStack},
		{"Approach 4: Morris Traversal (Optimal Space)", morrisPostorder},
	}

	for _, ap := range approaches {
		fmt.Printf("=== %s ===\n", ap.name)
		for i, ex := range examples {
			root := buildTree(ex.tree) // fresh tree per run (Morris mutates temporarily)
			got := ap.fn(root)
			fmt.Printf("Example %d: → %v (expected %v)\n", i+1, got, ex.expect)
			// expected: [3 2 1], [4 6 7 5 2 9 8 3 1], [], [1]
		}
		fmt.Println()
	}
}
