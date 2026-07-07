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

// ── Approach 1: Collect and Rebuild ──────────────────────────────────────────
//
// reverseBetween solves Reverse Linked List II by collecting values into a
// slice, reversing the subarray, and rebuilding.
//
// Time:  O(n)
// Space: O(n)
func reverseBetweenSimple(head *ListNode, left int, right int) *ListNode {
	vals := []int{}
	for cur := head; cur != nil; cur = cur.Next {
		vals = append(vals, cur.Val)
	}
	// reverse the subarray [left-1 .. right-1]
	l, r := left-1, right-1
	for l < r {
		vals[l], vals[r] = vals[r], vals[l]
		l++
		r--
	}
	// rebuild
	cur := head
	for i := range vals {
		cur.Val = vals[i]
		cur = cur.Next
	}
	return head
}

// ── Approach 2: One-Pass In-Place (Optimal) ───────────────────────────────────
//
// reverseBetween solves Reverse Linked List II with a single pass,
// reversing the sublist in-place using the "insert at front" technique.
//
// Intuition:
//   Find the node just before position `left` (call it `prev`). Then,
//   repeatedly take the node after `curr` (the start of the to-be-reversed
//   segment) and move it to just after `prev`. This inserts each new node
//   at the front of the reversed segment. After `right - left` iterations,
//   the sublist is reversed.
//
// Algorithm:
//   dummy.Next = head; prev = dummy
//   advance prev to position left-1 (the node before the reversal segment)
//   curr = prev.Next (first node of reversal segment)
//   for _ in 0..right-left-1:
//     next = curr.Next
//     curr.Next = next.Next
//     next.Next = prev.Next
//     prev.Next = next
//
// Time:  O(n)
// Space: O(1)
func reverseBetween(head *ListNode, left int, right int) *ListNode {
	dummy := &ListNode{Next: head}
	prev := dummy

	// advance prev to node at position left-1
	for i := 1; i < left; i++ {
		prev = prev.Next
	}

	curr := prev.Next // first node of reversal segment
	for i := 0; i < right-left; i++ {
		next := curr.Next       // node to move to front
		curr.Next = next.Next   // unlink next
		next.Next = prev.Next   // next points to front of reversed segment
		prev.Next = next        // prev now points to next (new front)
	}
	return dummy.Next
}

func main() {
	fmt.Println("=== Approach 1: Collect and Rebuild ===")
	l1 := makeList([]int{1, 2, 3, 4, 5})
	fmt.Printf("head=[1,2,3,4,5] left=2 right=4  got=%v  expected [1 4 3 2 5]\n", printList(reverseBetweenSimple(l1, 2, 4)))

	l2 := makeList([]int{5})
	fmt.Printf("head=[5] left=1 right=1  got=%v  expected [5]\n", printList(reverseBetweenSimple(l2, 1, 1)))

	fmt.Println("=== Approach 2: One-Pass In-Place ===")
	l3 := makeList([]int{1, 2, 3, 4, 5})
	fmt.Printf("head=[1,2,3,4,5] left=2 right=4  got=%v  expected [1 4 3 2 5]\n", printList(reverseBetween(l3, 2, 4)))

	l4 := makeList([]int{5})
	fmt.Printf("head=[5] left=1 right=1  got=%v  expected [5]\n", printList(reverseBetween(l4, 1, 1)))

	l5 := makeList([]int{1, 2, 3, 4, 5})
	fmt.Printf("head=[1,2,3,4,5] left=1 right=5  got=%v  expected [5 4 3 2 1]\n", printList(reverseBetween(l5, 1, 5)))
}
