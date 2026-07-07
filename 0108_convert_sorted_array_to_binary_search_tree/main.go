package main

import "fmt"

// TreeNode is a binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func inorderVals(root *TreeNode) []int {
	var res []int
	var dfs func(n *TreeNode)
	dfs = func(n *TreeNode) {
		if n == nil { return }
		dfs(n.Left); res = append(res, n.Val); dfs(n.Right)
	}
	dfs(root)
	return res
}

// ── Approach 1: Recursive Divide and Conquer ─────────────────────────────────
//
// sortedArrayToBST solves Convert Sorted Array to Binary Search Tree by always
// choosing the middle element as the root.
//
// Intuition:
//   A height-balanced BST has roughly equal numbers of nodes in left and right
//   subtrees. The middle element of the sorted array is the ideal root — it
//   splits the array into two equal halves. Recurse on each half.
//
// Algorithm:
//   build(lo, hi):
//     if lo > hi: return nil
//     mid = (lo + hi) / 2
//     root = nums[mid]
//     root.Left = build(lo, mid-1)
//     root.Right = build(mid+1, hi)
//
// Time:  O(n) — each element becomes a node.
// Space: O(log n) — recursion depth of balanced tree.
func sortedArrayToBST(nums []int) *TreeNode {
	var build func(lo, hi int) *TreeNode
	build = func(lo, hi int) *TreeNode {
		if lo > hi {
			return nil
		}
		mid := (lo + hi) / 2 // left-middle for even-length ranges
		root := &TreeNode{Val: nums[mid]}
		root.Left = build(lo, mid-1)
		root.Right = build(mid+1, hi)
		return root
	}
	return build(0, len(nums)-1)
}

// ── Approach 2: Iterative (Queue of Ranges) ──────────────────────────────────
//
// sortedArrayToBSTIterative solves Convert Sorted Array to Binary Search Tree
// iteratively using a queue of (node, lo, hi) triples.
//
// Intuition:
//   Same divide-and-conquer logic, but BFS-order using a queue.
//   Each queue entry holds the node that needs children and the range of nums
//   it should use.
//
// Time:  O(n)
// Space: O(n) — queue holds O(n) entries.
func sortedArrayToBSTIterative(nums []int) *TreeNode {
	if len(nums) == 0 {
		return nil
	}
	type entry struct {
		node   *TreeNode
		lo, hi int
		isLeft bool // true = attach as left child, false = right
		parent *TreeNode
	}

	mid := (0 + len(nums) - 1) / 2
	root := &TreeNode{Val: nums[mid]}

	type task struct {
		parent *TreeNode
		lo, hi int
		isLeft bool
	}
	queue := []task{
		{root, 0, mid - 1, true},
		{root, mid + 1, len(nums) - 1, false},
	}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		if curr.lo > curr.hi {
			continue
		}
		m := (curr.lo + curr.hi) / 2
		node := &TreeNode{Val: nums[m]}
		if curr.isLeft {
			curr.parent.Left = node
		} else {
			curr.parent.Right = node
		}
		queue = append(queue,
			task{node, curr.lo, m - 1, true},
			task{node, m + 1, curr.hi, false},
		)
	}
	return root
}

func main() {
	fmt.Println("=== Approach 1: Recursive Divide and Conquer ===")
	t1 := sortedArrayToBST([]int{-10, -3, 0, 5, 9})
	fmt.Printf("nums=[-10,-3,0,5,9]  inorder=%v  expected [-10 -3 0 5 9]\n", inorderVals(t1))

	t2 := sortedArrayToBST([]int{1, 3})
	fmt.Printf("nums=[1,3]  inorder=%v  expected [1 3]\n", inorderVals(t2))

	fmt.Println("=== Approach 2: Iterative Queue ===")
	t3 := sortedArrayToBSTIterative([]int{-10, -3, 0, 5, 9})
	fmt.Printf("nums=[-10,-3,0,5,9]  inorder=%v  expected [-10 -3 0 5 9]\n", inorderVals(t3))
}
