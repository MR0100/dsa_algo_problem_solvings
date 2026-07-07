package main

import "fmt"

// ListNode is a singly-linked list node.
type ListNode struct {
	Val  int
	Next *ListNode
}

// makeList builds a linked list from a slice.
func makeList(vals []int) *ListNode {
	dummy := &ListNode{}
	cur := dummy
	for _, v := range vals {
		cur.Next = &ListNode{Val: v}
		cur = cur.Next
	}
	return dummy.Next
}

// printList prints a linked list as a slice.
func printList(head *ListNode) []int {
	var result []int
	for head != nil {
		result = append(result, head.Val)
		head = head.Next
	}
	return result
}

// ── Approach 1: Two-Pointer with Dummy Head ───────────────────────────────────
//
// deleteDuplicates solves Remove Duplicates from Sorted List II by deleting
// all nodes that have duplicate numbers, leaving only distinct numbers.
//
// Intuition:
//   Use a dummy head so the result list's head can also be deleted if it's a
//   duplicate. Walk with pointer `prev` (last confirmed distinct node) and
//   peek ahead to skip entire runs of duplicate values.
//
// Algorithm:
//   dummy → head; prev = dummy
//   while prev.Next != nil:
//     cur = prev.Next
//     if cur.Next != nil && cur.Next.Val == cur.Val:
//       // cur is part of a duplicate run; skip the entire run
//       val = cur.Val
//       while prev.Next != nil && prev.Next.Val == val: prev.Next = prev.Next.Next
//     else:
//       prev = prev.Next  // cur is distinct; advance prev
//
// Time:  O(n)
// Space: O(1)
func deleteDuplicates(head *ListNode) *ListNode {
	dummy := &ListNode{Next: head}
	prev := dummy

	for prev.Next != nil {
		cur := prev.Next
		// check if cur starts a duplicate run
		if cur.Next != nil && cur.Next.Val == cur.Val {
			val := cur.Val
			// skip all nodes with this value
			for prev.Next != nil && prev.Next.Val == val {
				prev.Next = prev.Next.Next
			}
		} else {
			prev = prev.Next // cur is unique; advance prev
		}
	}
	return dummy.Next
}

// ── Approach 2: Recursive ─────────────────────────────────────────────────────
//
// deleteDuplicatesRecursive solves Remove Duplicates from Sorted List II using
// recursion: if the current node is a duplicate, skip the entire run and
// recurse on the remainder; otherwise keep it and recurse on the rest.
//
// Time:  O(n)
// Space: O(n) — recursion stack.
func deleteDuplicatesRecursive(head *ListNode) *ListNode {
	if head == nil {
		return nil
	}
	// current node is part of a duplicate run
	if head.Next != nil && head.Next.Val == head.Val {
		val := head.Val
		// skip all nodes in this run
		for head != nil && head.Val == val {
			head = head.Next
		}
		return deleteDuplicatesRecursive(head)
	}
	// current node is unique; keep it
	head.Next = deleteDuplicatesRecursive(head.Next)
	return head
}

func main() {
	fmt.Println("=== Approach 1: Two-Pointer with Dummy Head ===")
	l1 := makeList([]int{1, 2, 3, 3, 4, 4, 5})
	fmt.Printf("input=[1,2,3,3,4,4,5]  got=%v  expected [1 2 5]\n", printList(deleteDuplicates(l1)))

	l2 := makeList([]int{1, 1, 1, 2, 3})
	fmt.Printf("input=[1,1,1,2,3]  got=%v  expected [2 3]\n", printList(deleteDuplicates(l2)))

	l3 := makeList([]int{1, 1})
	fmt.Printf("input=[1,1]  got=%v  expected []\n", printList(deleteDuplicates(l3)))

	fmt.Println("=== Approach 2: Recursive ===")
	l4 := makeList([]int{1, 2, 3, 3, 4, 4, 5})
	fmt.Printf("input=[1,2,3,3,4,4,5]  got=%v  expected [1 2 5]\n", printList(deleteDuplicatesRecursive(l4)))

	l5 := makeList([]int{1, 1, 1, 2, 3})
	fmt.Printf("input=[1,1,1,2,3]  got=%v  expected [2 3]\n", printList(deleteDuplicatesRecursive(l5)))
}
