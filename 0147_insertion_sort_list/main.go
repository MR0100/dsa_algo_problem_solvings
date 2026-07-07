package main

import (
	"fmt"
	"sort"
	"strings"
)

// ListNode is the singly-linked list node used by LeetCode.
type ListNode struct {
	Val  int
	Next *ListNode
}

// buildList turns a slice of values into a linked list (test helper).
func buildList(vals []int) *ListNode {
	dummy := &ListNode{} // sentinel so we can append uniformly
	curr := dummy
	for _, v := range vals {
		curr.Next = &ListNode{Val: v}
		curr = curr.Next
	}
	return dummy.Next
}

// listString renders a list as "[a,b,c]" for printing (test helper).
func listString(head *ListNode) string {
	var parts []string
	for n := head; n != nil; n = n.Next {
		parts = append(parts, fmt.Sprint(n.Val))
	}
	return "[" + strings.Join(parts, ",") + "]"
}

// ── Approach 1: Brute Force (Copy Values + Array Sort) ───────────────────────
//
// bruteForceArraySort solves Insertion Sort List by dumping the node values
// into a slice, sorting the slice, and writing the values back.
//
// Intuition:
//
//	If we ignore the "insertion sort" instruction entirely, the easiest way
//	to sort a linked list is to move the problem into an array where sorting
//	is a one-liner, then copy the sorted values back over the same nodes.
//
// Algorithm:
//  1. Walk the list, appending every value to a slice.
//  2. sort.Ints the slice (O(n log n) introsort).
//  3. Walk the list again, overwriting node values in order.
//
// Time:  O(n log n) — dominated by the array sort.
// Space: O(n) — the value slice (violates the spirit of the problem).
func bruteForceArraySort(head *ListNode) *ListNode {
	var vals []int
	for n := head; n != nil; n = n.Next {
		vals = append(vals, n.Val) // collect every value
	}
	sort.Ints(vals) // let the stdlib sort the array
	i := 0
	for n := head; n != nil; n = n.Next {
		n.Val = vals[i] // overwrite values in sorted order; links untouched
		i++
	}
	return head
}

// ── Approach 2: Insertion Sort (Classic) ─────────────────────────────────────
//
// insertionSortList solves Insertion Sort List with the textbook algorithm:
// grow a sorted prefix list one node at a time.
//
// Intuition:
//
//	Exactly like sorting playing cards in your hand: take the next card
//	(node) from the unsorted remainder, scan the already-sorted hand from the
//	left, and splice the card in front of the first bigger card. A dummy head
//	makes "insert at the very front" the same code path as any other insert.
//
// Algorithm:
//  1. Create dummy → nil; dummy.Next will always be the sorted list.
//  2. For each node curr taken off the input list:
//     a. Detach curr (remember curr.Next first).
//     b. Scan from dummy while scan.Next.Val < curr.Val.
//     c. Splice curr between scan and scan.Next.
//  3. Return dummy.Next.
//
// Time:  O(n^2) worst/average — each insert may scan the whole sorted prefix.
//
//	O(n) best case is NOT achieved here (we always scan from the head).
//
// Space: O(1) — pointers only, nodes re-linked in place.
func insertionSortList(head *ListNode) *ListNode {
	dummy := &ListNode{} // sentinel head of the sorted result list
	curr := head         // next node to insert
	for curr != nil {
		next := curr.Next // save the rest of the unsorted list before relinking

		// find the insertion point: last sorted node with value < curr.Val
		scan := dummy
		for scan.Next != nil && scan.Next.Val < curr.Val {
			scan = scan.Next
		}

		// splice curr between scan and scan.Next
		curr.Next = scan.Next
		scan.Next = curr

		curr = next // move on to the next unsorted node
	}
	return dummy.Next
}

// ── Approach 3: Insertion Sort + Tail Shortcut (Optimal) ─────────────────────
//
// insertionSortTailOptimized solves Insertion Sort List with the classic
// algorithm plus a tail pointer that skips the scan for already-in-order nodes.
//
// Intuition:
//
//	In the classic version we re-scan from the head even when the incoming
//	node is >= everything sorted so far — which is every single step when the
//	input is already (nearly) sorted. Track the sorted list's tail: if
//	curr.Val >= tail.Val, just extend the tail in O(1). Only genuinely
//	out-of-order nodes pay for a scan. This restores insertion sort's famous
//	O(n) best case on sorted input while remaining O(1) space.
//
// Algorithm:
//  1. dummy sentinel + tail pointer to the last sorted node.
//  2. For each detached node curr:
//     a. If tail exists and curr.Val >= tail.Val → append after tail (O(1)).
//     b. Else scan from dummy for the insertion point and splice curr in;
//     if curr landed at the very end, it becomes the new tail.
//  3. Return dummy.Next.
//
// Time:  O(n^2) worst case, O(n) best case (sorted / nearly-sorted input).
// Space: O(1) — in-place pointer surgery.
func insertionSortTailOptimized(head *ListNode) *ListNode {
	dummy := &ListNode{} // sentinel head of the sorted result list
	var tail *ListNode   // last node of the sorted list (nil while empty)
	for curr := head; curr != nil; {
		next := curr.Next // save the unsorted remainder

		if tail != nil && curr.Val >= tail.Val {
			// fast path: curr belongs at the end — no scan needed
			tail.Next = curr
			curr.Next = nil
			tail = curr
		} else {
			// slow path: scan from the head for the first bigger node
			scan := dummy
			for scan.Next != nil && scan.Next.Val < curr.Val {
				scan = scan.Next
			}
			curr.Next = scan.Next
			scan.Next = curr
			if curr.Next == nil {
				tail = curr // inserted at the end → curr is the new tail
			}
		}

		curr = next
	}
	return dummy.Next
}

func main() {
	// Official LeetCode examples.
	example1 := []int{4, 2, 1, 3}
	example2 := []int{-1, 5, 3, 4, 0}

	fmt.Println("=== Approach 1: Brute Force (Copy Values + Array Sort) ===")
	fmt.Println(listString(bruteForceArraySort(buildList(example1)))) // [1,2,3,4]
	fmt.Println(listString(bruteForceArraySort(buildList(example2)))) // [-1,0,3,4,5]

	fmt.Println("=== Approach 2: Insertion Sort (Classic) ===")
	fmt.Println(listString(insertionSortList(buildList(example1)))) // [1,2,3,4]
	fmt.Println(listString(insertionSortList(buildList(example2)))) // [-1,0,3,4,5]

	fmt.Println("=== Approach 3: Insertion Sort + Tail Shortcut (Optimal) ===")
	fmt.Println(listString(insertionSortTailOptimized(buildList(example1)))) // [1,2,3,4]
	fmt.Println(listString(insertionSortTailOptimized(buildList(example2)))) // [-1,0,3,4,5]
}
