package main

import "fmt"

// ListNode is a singly-linked list node.
type ListNode struct {
	Val  int
	Next *ListNode
}

// ── Approach 1: Two-Pass (Compute Length) ────────────────────────────────────
//
// twoPass first counts the list length, then removes the (length-n)-th node
// from the beginning (0-indexed).
//
// Intuition:
//   The n-th node from the end is the (L-n+1)-th node from the start (1-indexed).
//   First pass: count L. Second pass: walk to node L-n, unlink the next node.
//   A dummy head node simplifies edge cases (removing the actual head).
//
// Time:  O(L) — two traversals of length L.
// Space: O(1) — only pointer variables.
func twoPass(head *ListNode, n int) *ListNode {
	dummy := &ListNode{Next: head}

	// First pass: count list length.
	length := 0
	for cur := head; cur != nil; cur = cur.Next {
		length++
	}

	// Walk to the node just before the one to remove.
	// We need to reach index (length - n) from the dummy head.
	cur := dummy
	for i := 0; i < length-n; i++ {
		cur = cur.Next
	}
	cur.Next = cur.Next.Next // unlink the target node

	return dummy.Next
}

// ── Approach 2: One-Pass with Two Pointers (Optimal) ─────────────────────────
//
// onePass uses a fast pointer n steps ahead of the slow pointer. When fast
// reaches the end, slow is exactly at the node before the target.
//
// Intuition:
//   Advance `fast` n+1 steps ahead of `slow` (both start at a dummy head).
//   Then advance both together until `fast` is nil. At that point, `slow`
//   is at the node just before the n-th from the end. Unlink `slow.Next`.
//
//   Why n+1? We want slow to stop at the predecessor of the target, not the
//   target itself. So the gap between fast and slow is n+1 nodes.
//
// Time:  O(L) — single traversal.
// Space: O(1).
func onePass(head *ListNode, n int) *ListNode {
	dummy := &ListNode{Next: head}
	fast, slow := dummy, dummy

	// Advance fast n+1 steps so the gap between fast and slow is n+1.
	for i := 0; i <= n; i++ {
		fast = fast.Next
	}

	// Move both until fast reaches nil.
	for fast != nil {
		fast = fast.Next
		slow = slow.Next
	}

	// slow is now the node before the n-th from the end.
	slow.Next = slow.Next.Next

	return dummy.Next
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func makeList(vals []int) *ListNode {
	dummy := &ListNode{}
	cur := dummy
	for _, v := range vals {
		cur.Next = &ListNode{Val: v}
		cur = cur.Next
	}
	return dummy.Next
}

func listToSlice(head *ListNode) []int {
	var result []int
	for head != nil {
		result = append(result, head.Val)
		head = head.Next
	}
	return result
}

func main() {
	examples := []struct {
		vals   []int
		n      int
		expect []int
	}{
		{[]int{1, 2, 3, 4, 5}, 2, []int{1, 2, 3, 5}},
		{[]int{1}, 1, []int{}},
		{[]int{1, 2}, 1, []int{1}},
		{[]int{1, 2}, 2, []int{2}},
	}

	approaches := []struct {
		name string
		fn   func(*ListNode, int) *ListNode
	}{
		{"Approach 1: Two-Pass        O(L) T | O(1) S", twoPass},
		{"Approach 2: One-Pass (2ptr) ✅ O(L) T | O(1) S", onePass},
	}

	for _, ex := range examples {
		fmt.Printf("list=%v  n=%d  expect=%v\n", ex.vals, ex.n, ex.expect)
		for _, ap := range approaches {
			head := makeList(ex.vals)
			out := ap.fn(head, ex.n)
			fmt.Printf("  %-50s → %v\n", ap.name, listToSlice(out))
		}
		fmt.Println()
	}
}
