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

// ── Approach 1: Two Dummy Lists ───────────────────────────────────────────────
//
// partition solves Partition List by splitting into two sublists: one for
// nodes < x and one for nodes >= x, then concatenating.
//
// Intuition:
//   Create two separate lists using dummy heads: "less" collects nodes with
//   Val < x, "greater" collects nodes with Val >= x. Walk through the original
//   list and append each node to the appropriate sublist. Then concatenate.
//
// Algorithm:
//   lessHead, greaterHead = dummy nodes
//   less, greater = their tails
//   for each node:
//     if node.Val < x: less.Next = node; less = less.Next
//     else: greater.Next = node; greater = greater.Next
//   greater.Next = nil  // terminate greater list
//   less.Next = greaterHead.Next  // join the two lists
//   return lessHead.Next
//
// Time:  O(n)
// Space: O(1) — rearranges existing nodes, no new allocations.
func partition(head *ListNode, x int) *ListNode {
	lessHead := &ListNode{}  // dummy head for the "less than x" list
	greaterHead := &ListNode{} // dummy head for the "greater or equal to x" list

	less := lessHead
	greater := greaterHead

	for head != nil {
		if head.Val < x {
			less.Next = head
			less = less.Next
		} else {
			greater.Next = head
			greater = greater.Next
		}
		head = head.Next
	}
	greater.Next = nil         // prevent cycle (last node may point somewhere old)
	less.Next = greaterHead.Next // connect the two sublists

	return lessHead.Next
}

func main() {
	fmt.Println("=== Approach 1: Two Dummy Lists ===")
	l1 := makeList([]int{1, 4, 3, 2, 5, 2})
	fmt.Printf("head=[1,4,3,2,5,2] x=3  got=%v  expected [1 2 2 4 3 5]\n", printList(partition(l1, 3)))

	l2 := makeList([]int{2, 1})
	fmt.Printf("head=[2,1] x=2  got=%v  expected [1 2]\n", printList(partition(l2, 2)))

	l3 := makeList([]int{1})
	fmt.Printf("head=[1] x=2  got=%v  expected [1]\n", printList(partition(l3, 2)))

	l4 := makeList([]int{3, 1, 2})
	fmt.Printf("head=[3,1,2] x=3  got=%v  expected [1 2 3]\n", printList(partition(l4, 3)))
}
