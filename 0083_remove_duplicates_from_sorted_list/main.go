package main

import "fmt"

// ListNode is a singly-linked list node.
type ListNode struct {
	Val  int
	Next *ListNode
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

func printList(head *ListNode) []int {
	var result []int
	for head != nil {
		result = append(result, head.Val)
		head = head.Next
	}
	return result
}

// ── Approach 1: Iterative ─────────────────────────────────────────────────────
//
// deleteDuplicates solves Remove Duplicates from Sorted List by keeping
// exactly one occurrence of each value.
//
// Intuition:
//   Walk through the list. If the current node's value equals its next node's
//   value, skip the next node (by setting cur.Next = cur.Next.Next). Only
//   advance cur when the next node has a different value.
//
// Algorithm:
//   cur = head
//   while cur != nil && cur.Next != nil:
//     if cur.Val == cur.Next.Val: cur.Next = cur.Next.Next
//     else: cur = cur.Next
//   return head
//
// Time:  O(n)
// Space: O(1)
func deleteDuplicates(head *ListNode) *ListNode {
	cur := head
	for cur != nil && cur.Next != nil {
		if cur.Val == cur.Next.Val {
			cur.Next = cur.Next.Next // skip duplicate
		} else {
			cur = cur.Next // distinct value; advance
		}
	}
	return head
}

// ── Approach 2: Recursive ─────────────────────────────────────────────────────
//
// deleteDuplicatesRecursive solves Remove Duplicates from Sorted List
// recursively.
//
// Intuition:
//   Recursively deduplicate the rest of the list. If the current node has the
//   same value as its (now-deduped) next node, skip the current node.
//
// Time:  O(n)
// Space: O(n) — recursion stack.
func deleteDuplicatesRecursive(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}
	head.Next = deleteDuplicatesRecursive(head.Next)
	if head.Val == head.Next.Val {
		return head.Next // skip current; its duplicate is next
	}
	return head
}

func main() {
	fmt.Println("=== Approach 1: Iterative ===")
	l1 := makeList([]int{1, 1, 2})
	fmt.Printf("input=[1,1,2]  got=%v  expected [1 2]\n", printList(deleteDuplicates(l1)))

	l2 := makeList([]int{1, 1, 2, 3, 3})
	fmt.Printf("input=[1,1,2,3,3]  got=%v  expected [1 2 3]\n", printList(deleteDuplicates(l2)))

	l3 := makeList([]int{1})
	fmt.Printf("input=[1]  got=%v  expected [1]\n", printList(deleteDuplicates(l3)))

	fmt.Println("=== Approach 2: Recursive ===")
	l4 := makeList([]int{1, 1, 2})
	fmt.Printf("input=[1,1,2]  got=%v  expected [1 2]\n", printList(deleteDuplicatesRecursive(l4)))

	l5 := makeList([]int{1, 1, 2, 3, 3})
	fmt.Printf("input=[1,1,2,3,3]  got=%v  expected [1 2 3]\n", printList(deleteDuplicatesRecursive(l5)))
}
