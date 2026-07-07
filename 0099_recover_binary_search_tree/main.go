package main

import "fmt"

// TreeNode is a binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// buildTree builds a BST from sorted values (for testing).
func buildBST(vals []int, start, end int) *TreeNode {
	if start > end {
		return nil
	}
	mid := (start + end) / 2
	return &TreeNode{Val: vals[mid], Left: buildBST(vals, start, mid-1), Right: buildBST(vals, mid+1, end)}
}

func inorderVals(root *TreeNode) []int {
	var result []int
	var dfs func(n *TreeNode)
	dfs = func(n *TreeNode) {
		if n == nil { return }
		dfs(n.Left); result = append(result, n.Val); dfs(n.Right)
	}
	dfs(root)
	return result
}

// ── Approach 1: Inorder Traversal + Find Two Swapped Nodes ───────────────────
//
// recoverTree solves Recover Binary Search Tree by finding the two nodes that
// are out of order in the inorder traversal.
//
// Intuition:
//   In a valid BST's inorder traversal, all values are strictly increasing.
//   When exactly two nodes are swapped, the inorder sequence has exactly two
//   "inversion" points where current value < previous value.
//
//   Case 1 (adjacent swap): only 1 inversion point. Swap prev and curr.
//   Case 2 (non-adjacent swap): 2 inversion points. Swap first prev and last curr.
//
//   Track `first` (prev at first inversion) and `second` (curr at second inversion
//   — or curr at first inversion if only one).
//
// Time:  O(n)
// Space: O(h) — recursion stack.
func recoverTree(root *TreeNode) {
	var first, second, prev *TreeNode

	var inorder func(node *TreeNode)
	inorder = func(node *TreeNode) {
		if node == nil {
			return
		}
		inorder(node.Left)
		if prev != nil && prev.Val > node.Val {
			if first == nil {
				first = prev // first inversion: record prev
			}
			second = node // second (or only) inversion: always update second
		}
		prev = node
		inorder(node.Right)
	}
	inorder(root)

	// swap the values of the two misplaced nodes
	first.Val, second.Val = second.Val, first.Val
}

// ── Approach 2: Morris Inorder Traversal (O(1) Space) ────────────────────────
//
// recoverTreeMorris solves Recover Binary Search Tree in O(1) extra space
// using Morris inorder traversal.
//
// Time:  O(n)
// Space: O(1)
func recoverTreeMorris(root *TreeNode) {
	var first, second, prev *TreeNode
	curr := root

	for curr != nil {
		if curr.Left == nil {
			// visit curr
			if prev != nil && prev.Val > curr.Val {
				if first == nil { first = prev }
				second = curr
			}
			prev = curr
			curr = curr.Right
		} else {
			// find inorder predecessor
			pred := curr.Left
			for pred.Right != nil && pred.Right != curr {
				pred = pred.Right
			}
			if pred.Right == nil {
				pred.Right = curr // thread
				curr = curr.Left
			} else {
				pred.Right = nil // unthread
				// visit curr
				if prev != nil && prev.Val > curr.Val {
					if first == nil { first = prev }
					second = curr
				}
				prev = curr
				curr = curr.Right
			}
		}
	}
	first.Val, second.Val = second.Val, first.Val
}

func main() {
	// Example 1: [1,3,null,null,2] — swap 1 and 3 → should recover to [3,1,null,null,2]
	// Actually: tree with root=1, left=3, 3.right=2
	// Inorder: 3,2,1 → wrong (should be 1,2,3)
	// After fix: swap 1 and 3 → root=3, left=1, 1.right=2
	t1 := &TreeNode{Val: 1, Left: &TreeNode{Val: 3, Right: &TreeNode{Val: 2}}}
	fmt.Println("=== Approach 1: Inorder Traversal ===")
	fmt.Printf("before: inorder=%v\n", inorderVals(t1))
	recoverTree(t1)
	fmt.Printf("after:  inorder=%v  expected [1 2 3]\n", inorderVals(t1))

	// Example 2: [3,1,4,null,null,2] — swap 2 and 3
	t2 := &TreeNode{Val: 3,
		Left:  &TreeNode{Val: 1},
		Right: &TreeNode{Val: 4, Left: &TreeNode{Val: 2}},
	}
	fmt.Printf("before: inorder=%v\n", inorderVals(t2))
	recoverTree(t2)
	fmt.Printf("after:  inorder=%v  expected [1 2 3 4]\n", inorderVals(t2))

	fmt.Println("=== Approach 2: Morris Traversal ===")
	t3 := &TreeNode{Val: 1, Left: &TreeNode{Val: 3, Right: &TreeNode{Val: 2}}}
	fmt.Printf("before: inorder=%v\n", inorderVals(t3))
	recoverTreeMorris(t3)
	fmt.Printf("after:  inorder=%v  expected [1 2 3]\n", inorderVals(t3))

	t4 := &TreeNode{Val: 3,
		Left:  &TreeNode{Val: 1},
		Right: &TreeNode{Val: 4, Left: &TreeNode{Val: 2}},
	}
	fmt.Printf("before: inorder=%v\n", inorderVals(t4))
	recoverTreeMorris(t4)
	fmt.Printf("after:  inorder=%v  expected [1 2 3 4]\n", inorderVals(t4))
}
