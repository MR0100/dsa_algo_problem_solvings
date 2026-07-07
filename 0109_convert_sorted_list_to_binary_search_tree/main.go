package main

import "fmt"

// ListNode is a singly linked list node.
type ListNode struct {
	Val  int
	Next *ListNode
}

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

func makeList(vals []int) *ListNode {
	dummy := &ListNode{}
	cur := dummy
	for _, v := range vals {
		cur.Next = &ListNode{Val: v}
		cur = cur.Next
	}
	return dummy.Next
}

// ── Approach 1: Convert to Array, then #108 ──────────────────────────────────
//
// sortedListToBST solves Convert Sorted List to BST by first collecting all
// values into a slice, then applying the divide-and-conquer mid-element approach.
//
// Intuition:
//   Same as #108 once we have random access. Extra O(n) space for the array.
//
// Time:  O(n)
// Space: O(n) — array copy + recursion stack.
func sortedListToBST(head *ListNode) *TreeNode {
	// collect all values
	var nums []int
	for cur := head; cur != nil; cur = cur.Next {
		nums = append(nums, cur.Val)
	}

	var build func(lo, hi int) *TreeNode
	build = func(lo, hi int) *TreeNode {
		if lo > hi {
			return nil
		}
		mid := (lo + hi) / 2
		root := &TreeNode{Val: nums[mid]}
		root.Left = build(lo, mid-1)
		root.Right = build(mid+1, hi)
		return root
	}
	return build(0, len(nums)-1)
}

// ── Approach 2: Slow/Fast Pointer (In-Order Simulation) ─────────────────────
//
// sortedListToBSTInOrder solves Convert Sorted List to BST using inorder
// simulation: advance the list pointer as we recurse.
//
// Intuition:
//   In inorder traversal, we visit left → root → right.
//   If we recurse in the same order and advance the list pointer at root-visit
//   time, each node gets exactly its correct value without random access.
//   We use the size of the range [lo, hi] to know where to split.
//
// Algorithm:
//   build(lo, hi):
//     if lo > hi: return nil
//     mid = (lo+hi)/2
//     left = build(lo, mid-1)      ← advances cur to mid-th element
//     root = &TreeNode{Val: cur.Val}
//     cur = cur.Next               ← consume mid-th element
//     root.Right = build(mid+1, hi)
//     return root
//
// Time:  O(n) — each node touched once.
// Space: O(log n) — recursion stack (balanced BST depth).
func sortedListToBSTInOrder(head *ListNode) *TreeNode {
	// count list length
	length := 0
	for cur := head; cur != nil; cur = cur.Next {
		length++
	}

	cur := head // shared pointer advanced during inorder traversal
	var build func(lo, hi int) *TreeNode
	build = func(lo, hi int) *TreeNode {
		if lo > hi {
			return nil
		}
		mid := (lo + hi) / 2
		left := build(lo, mid-1) // builds left subtree, advances cur
		root := &TreeNode{Val: cur.Val}
		cur = cur.Next // consume cur value for this root
		root.Left = left
		root.Right = build(mid+1, hi)
		return root
	}
	return build(0, length-1)
}

func main() {
	head1 := makeList([]int{-10, -3, 0, 5, 9})
	head2 := makeList([]int{})
	head3 := makeList([]int{-10, -3, 0, 5, 9})

	fmt.Println("=== Approach 1: Array Copy ===")
	t1 := sortedListToBST(head1)
	fmt.Printf("list=[-10,-3,0,5,9]  inorder=%v  expected [-10 -3 0 5 9]\n", inorderVals(t1))
	fmt.Printf("list=[]  result=%v  expected []\n", inorderVals(sortedListToBST(head2)))

	fmt.Println("=== Approach 2: Inorder Simulation (O(log n) space) ===")
	t3 := sortedListToBSTInOrder(head3)
	fmt.Printf("list=[-10,-3,0,5,9]  inorder=%v  expected [-10 -3 0 5 9]\n", inorderVals(t3))
}
