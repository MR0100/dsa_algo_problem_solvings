package main

import "fmt"

type ListNode struct {
	Val  int
	Next *ListNode
}

// ── Approach 1: Iterative with Dummy Head ─────────────────────────────────────
//
// iterative uses a dummy head node and a current pointer to build the merged
// list in-place by relinking existing nodes.
//
// Intuition:
//   Compare the front of each list; attach the smaller node to the result,
//   advance that list's pointer. When one list is exhausted, attach the rest
//   of the other list in O(1) (it is already sorted).
//
// Time:  O(m+n) — visit every node exactly once.
// Space: O(1)   — rearranges existing nodes; only the dummy head is extra.
func iterative(list1 *ListNode, list2 *ListNode) *ListNode {
	dummy := &ListNode{}
	cur := dummy

	for list1 != nil && list2 != nil {
		if list1.Val <= list2.Val {
			cur.Next = list1
			list1 = list1.Next
		} else {
			cur.Next = list2
			list2 = list2.Next
		}
		cur = cur.Next
	}

	// Attach the remaining non-empty list (already sorted).
	if list1 != nil {
		cur.Next = list1
	} else {
		cur.Next = list2
	}

	return dummy.Next
}

// ── Approach 2: Recursive ─────────────────────────────────────────────────────
//
// recursive picks the smaller head, recurses for the rest, and links the result.
//
// Intuition:
//   The merged list's head is whichever of list1.Val and list2.Val is smaller.
//   Its next pointer is the merge of the remaining nodes — a strictly smaller
//   sub-problem. Base case: either list is nil → return the other.
//
// Time:  O(m+n) — one call per node.
// Space: O(m+n) — recursion stack depth equals total nodes.
func recursive(list1 *ListNode, list2 *ListNode) *ListNode {
	if list1 == nil {
		return list2
	}
	if list2 == nil {
		return list1
	}

	if list1.Val <= list2.Val {
		list1.Next = recursive(list1.Next, list2)
		return list1
	}
	list2.Next = recursive(list1, list2.Next)
	return list2
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
	var out []int
	for head != nil {
		out = append(out, head.Val)
		head = head.Next
	}
	return out
}

func main() {
	examples := []struct {
		l1, l2 []int
		expect []int
	}{
		{[]int{1, 2, 4}, []int{1, 3, 4}, []int{1, 1, 2, 3, 4, 4}},
		{nil, nil, nil},
		{nil, []int{0}, []int{0}},
		{[]int{1, 3, 5}, []int{2, 4, 6}, []int{1, 2, 3, 4, 5, 6}},
	}

	approaches := []struct {
		name string
		fn   func(*ListNode, *ListNode) *ListNode
	}{
		{"Approach 1: Iterative  ✅ O(m+n) T | O(1)   S", iterative},
		{"Approach 2: Recursive     O(m+n) T | O(m+n) S", recursive},
	}

	for _, ex := range examples {
		fmt.Printf("l1=%v  l2=%v  expect=%v\n", ex.l1, ex.l2, ex.expect)
		for _, ap := range approaches {
			result := ap.fn(makeList(ex.l1), makeList(ex.l2))
			fmt.Printf("  %-50s → %v\n", ap.name, listToSlice(result))
		}
		fmt.Println()
	}
}
