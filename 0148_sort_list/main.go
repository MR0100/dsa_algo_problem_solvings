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
// bruteForceArraySort solves Sort List by dumping node values into a slice,
// sorting the slice, and writing the values back over the nodes.
//
// Intuition:
//
//	Arrays are the natural habitat of sorting algorithms. Move the values
//	into a slice, use the stdlib sort, and copy back in order. Meets the
//	O(n log n) time requirement but spends O(n) extra memory, so it fails
//	the follow-up's O(1)-space challenge.
//
// Algorithm:
//  1. Walk the list collecting values into a slice.
//  2. sort.Ints the slice.
//  3. Walk the list again, overwriting each node's value in order.
//
// Time:  O(n log n) — dominated by the array sort.
// Space: O(n) — the values slice.
func bruteForceArraySort(head *ListNode) *ListNode {
	var vals []int
	for n := head; n != nil; n = n.Next {
		vals = append(vals, n.Val) // collect every value
	}
	sort.Ints(vals) // stdlib sort on the array copy
	i := 0
	for n := head; n != nil; n = n.Next {
		n.Val = vals[i] // rewrite values in sorted order; links untouched
		i++
	}
	return head
}

// ── Approach 2: Merge Sort Top-Down (Recursive) ──────────────────────────────
//
// mergeSortTopDown solves Sort List with recursive merge sort: split the list
// in half via slow/fast pointers, sort each half, merge the sorted halves.
//
// Intuition:
//
//	Merge sort is THE sorting algorithm for linked lists: unlike quicksort it
//	needs no random access, and the merge step is pure pointer re-linking with
//	zero extra memory. Find the middle with slow/fast pointers, cut the list,
//	recursively sort both halves, then weave them together.
//
// Algorithm:
//  1. Base case: 0 or 1 node → already sorted.
//  2. slow/fast pointers find the node BEFORE the middle; cut the list there.
//  3. Recurse on both halves.
//  4. merge(): standard two-sorted-lists merge with a dummy head.
//
// Time:  O(n log n) — log n levels of splitting, O(n) merging per level.
// Space: O(log n) — recursion stack (the merges themselves are in-place).
func mergeSortTopDown(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head // 0 or 1 node: already sorted
	}

	// slow/fast: fast starts one ahead so slow stops at the END of the first half
	slow, fast := head, head.Next
	for fast != nil && fast.Next != nil {
		slow = slow.Next      // slow moves 1 step
		fast = fast.Next.Next // fast moves 2 steps
	}
	mid := slow.Next // head of the second half
	slow.Next = nil  // cut the list into two independent halves

	left := mergeSortTopDown(head) // sort first half
	right := mergeSortTopDown(mid) // sort second half
	return merge(left, right)      // weave the sorted halves together
}

// merge combines two sorted lists into one sorted list in O(len1+len2).
func merge(a, b *ListNode) *ListNode {
	dummy := &ListNode{} // sentinel so appending is uniform
	tail := dummy
	for a != nil && b != nil {
		if a.Val <= b.Val { // <= keeps the sort stable
			tail.Next = a
			a = a.Next
		} else {
			tail.Next = b
			b = b.Next
		}
		tail = tail.Next
	}
	// at most one of a/b is non-nil; append the leftover run
	if a != nil {
		tail.Next = a
	} else {
		tail.Next = b
	}
	return dummy.Next
}

// ── Approach 3: Merge Sort Bottom-Up (Optimal, O(1) Space) ───────────────────
//
// mergeSortBottomUp solves Sort List iteratively: merge runs of size 1, then
// 2, then 4, ... until one run covers the whole list. Answers the follow-up.
//
// Intuition:
//
//	Top-down recursion costs O(log n) stack. Flip the direction: instead of
//	splitting from the top, start from the bottom where every single node is
//	already a sorted run of length 1. Pass 1 merges neighbouring runs into
//	sorted runs of 2, pass 2 into runs of 4, ... After ceil(log2 n) passes the
//	whole list is one sorted run — no recursion, O(1) extra space.
//
// Algorithm:
//  1. Count the list length n.
//  2. For width = 1, 2, 4, ... < n:
//     a. Walk the list, repeatedly: split off run A (width nodes) and run B
//     (next width nodes), merge them, and append the merged run to the
//     growing result via a tail pointer.
//  3. Return dummy.Next.
//
// Time:  O(n log n) — ceil(log2 n) passes, each touching all n nodes.
// Space: O(1) — a fixed set of pointers, no recursion.
func mergeSortBottomUp(head *ListNode) *ListNode {
	// count the length so we know when runs cover the whole list
	n := 0
	for node := head; node != nil; node = node.Next {
		n++
	}

	dummy := &ListNode{Next: head} // sentinel; dummy.Next is always the current list
	for width := 1; width < n; width *= 2 {
		tail := dummy      // end of the already-merged prefix for this pass
		curr := dummy.Next // start of the not-yet-merged remainder
		for curr != nil {
			left := curr                        // run A: `width` nodes starting at curr
			right := split(left, width)         // run B starts where A was cut
			curr = split(right, width)          // remainder starts where B was cut
			tail = mergeTail(left, right, tail) // merge A+B, append after tail
		}
	}
	return dummy.Next
}

// split cuts off the first n nodes of head and returns the head of the rest.
func split(head *ListNode, n int) *ListNode {
	for i := 1; head != nil && i < n; i++ {
		head = head.Next // advance to the n-th node (run may be shorter)
	}
	if head == nil {
		return nil // run shorter than n: nothing left after it
	}
	rest := head.Next // remainder begins after the n-th node
	head.Next = nil   // terminate the first run
	return rest
}

// mergeTail merges sorted runs a and b, appends the result after tail,
// and returns the new tail (last node of the merged run).
func mergeTail(a, b *ListNode, tail *ListNode) *ListNode {
	curr := tail // build directly onto the result list
	for a != nil && b != nil {
		if a.Val <= b.Val { // <= keeps the sort stable
			curr.Next = a
			a = a.Next
		} else {
			curr.Next = b
			b = b.Next
		}
		curr = curr.Next
	}
	// attach whichever run still has nodes
	if a != nil {
		curr.Next = a
	} else {
		curr.Next = b
	}
	for curr.Next != nil {
		curr = curr.Next // walk to the end so the caller gets the new tail
	}
	return curr
}

func main() {
	// Official LeetCode examples.
	example1 := []int{4, 2, 1, 3}
	example2 := []int{-1, 5, 3, 4, 0}
	example3 := []int{} // empty list

	fmt.Println("=== Approach 1: Brute Force (Copy Values + Array Sort) ===")
	fmt.Println(listString(bruteForceArraySort(buildList(example1)))) // [1,2,3,4]
	fmt.Println(listString(bruteForceArraySort(buildList(example2)))) // [-1,0,3,4,5]
	fmt.Println(listString(bruteForceArraySort(buildList(example3)))) // []

	fmt.Println("=== Approach 2: Merge Sort Top-Down (Recursive) ===")
	fmt.Println(listString(mergeSortTopDown(buildList(example1)))) // [1,2,3,4]
	fmt.Println(listString(mergeSortTopDown(buildList(example2)))) // [-1,0,3,4,5]
	fmt.Println(listString(mergeSortTopDown(buildList(example3)))) // []

	fmt.Println("=== Approach 3: Merge Sort Bottom-Up (Optimal, O(1) Space) ===")
	fmt.Println(listString(mergeSortBottomUp(buildList(example1)))) // [1,2,3,4]
	fmt.Println(listString(mergeSortBottomUp(buildList(example2)))) // [-1,0,3,4,5]
	fmt.Println(listString(mergeSortBottomUp(buildList(example3)))) // []
}
