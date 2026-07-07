package main

import "fmt"

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
	var res []int
	for head != nil {
		res = append(res, head.Val)
		head = head.Next
	}
	return res
}

// ── Approach 1: Array Copy ────────────────────────────────────────────────────
//
// arrayCopy solves Rotate List by collecting all values, rotating the slice,
// and rebuilding the list.
//
// Intuition:
//   Collect all node values, compute effective rotation k%n, copy the rotated
//   values back into the original list nodes.
//
// Time:  O(n)
// Space: O(n) — extra slice.
func arrayCopy(head *ListNode, k int) *ListNode {
	if head == nil || head.Next == nil || k == 0 {
		return head
	}
	// collect values
	var vals []int
	for cur := head; cur != nil; cur = cur.Next {
		vals = append(vals, cur.Val)
	}
	n := len(vals)
	k = k % n
	if k == 0 {
		return head
	}
	// rotate: last k elements move to front
	rotated := append(vals[n-k:], vals[:n-k]...)
	// write back
	cur := head
	for _, v := range rotated {
		cur.Val = v
		cur = cur.Next
	}
	return head
}

// ── Approach 2: Find Tail and Reconnect ──────────────────────────────────────
//
// reconnect solves Rotate List in O(n) time O(1) space by making the list
// circular and then breaking it at the new tail position.
//
// Intuition:
//   Connect the tail to the head (make circular). The new tail is at position
//   n - k - 1 (0-indexed from head). The new head is the node after the new tail.
//   Break the circle there.
//
// Algorithm:
//   1. Find length n and the tail node.
//   2. k = k % n. If k == 0, return head unchanged.
//   3. Advance (n - k - 1) steps from head to reach new tail.
//   4. new head = new tail's next; break link.
//
// Time:  O(n)
// Space: O(1)
func reconnect(head *ListNode, k int) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}
	// find length and tail
	n := 1
	tail := head
	for tail.Next != nil {
		tail = tail.Next
		n++
	}
	k = k % n
	if k == 0 {
		return head
	}
	// make circular
	tail.Next = head
	// advance to new tail (n - k - 1 steps)
	newTail := head
	for i := 0; i < n-k-1; i++ {
		newTail = newTail.Next
	}
	newHead := newTail.Next
	newTail.Next = nil // break the circle
	return newHead
}

func main() {
	fmt.Println("=== Approach 1: Array Copy ===")
	fmt.Printf("head=[1,2,3,4,5] k=2  got=%v  expected [4 5 1 2 3]\n",
		printList(arrayCopy(makeList([]int{1, 2, 3, 4, 5}), 2)))
	fmt.Printf("head=[0,1,2] k=4  got=%v  expected [2 0 1]\n",
		printList(arrayCopy(makeList([]int{0, 1, 2}), 4)))
	fmt.Printf("head=[1,2] k=0  got=%v  expected [1 2]\n",
		printList(arrayCopy(makeList([]int{1, 2}), 0)))

	fmt.Println("=== Approach 2: Find Tail and Reconnect ===")
	fmt.Printf("head=[1,2,3,4,5] k=2  got=%v  expected [4 5 1 2 3]\n",
		printList(reconnect(makeList([]int{1, 2, 3, 4, 5}), 2)))
	fmt.Printf("head=[0,1,2] k=4  got=%v  expected [2 0 1]\n",
		printList(reconnect(makeList([]int{0, 1, 2}), 4)))
	fmt.Printf("head=[1,2] k=0  got=%v  expected [1 2]\n",
		printList(reconnect(makeList([]int{1, 2}), 0)))
}
