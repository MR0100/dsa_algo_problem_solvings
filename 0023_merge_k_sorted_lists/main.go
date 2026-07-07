package main

import (
	"container/heap"
	"fmt"
)

type ListNode struct {
	Val  int
	Next *ListNode
}

// ── Approach 1: Brute Force — Collect All Values and Sort ─────────────────────
//
// bruteForce collects every node's value, sorts the slice, and builds a new list.
//
// Intuition:
//   The simplest approach: ignore the "already sorted" property, pool all values,
//   sort, rebuild.
//
// Time:  O(N log N) where N = total nodes across all lists.
// Space: O(N) — the values slice.
func bruteForce(lists []*ListNode) *ListNode {
	var vals []int
	for _, head := range lists {
		for head != nil {
			vals = append(vals, head.Val)
			head = head.Next
		}
	}
	// Simple insertion sort to avoid importing sort (keeps imports clean).
	for i := 1; i < len(vals); i++ {
		for j := i; j > 0 && vals[j] < vals[j-1]; j-- {
			vals[j], vals[j-1] = vals[j-1], vals[j]
		}
	}
	dummy := &ListNode{}
	cur := dummy
	for _, v := range vals {
		cur.Next = &ListNode{Val: v}
		cur = cur.Next
	}
	return dummy.Next
}

// ── Approach 2: Compare One by One (Sequential Merge) ────────────────────────
//
// sequentialMerge merges lists one at a time into a running result.
//
// Intuition:
//   Merge lists[0] with lists[1], then merge that result with lists[2], etc.
//   Each merge is O(a+b) where a and b are the merged lengths so far.
//   This is like insertion sort on lists.
//
// Time:  O(k·N) — merging list i into the running result takes O(i·(N/k) + N/k)
//                  for k lists of N/k nodes each; total ≈ O(k·N).
// Space: O(1) — in-place relinking.
func sequentialMerge(lists []*ListNode) *ListNode {
	if len(lists) == 0 {
		return nil
	}
	result := lists[0]
	for i := 1; i < len(lists); i++ {
		result = mergeTwoLists(result, lists[i])
	}
	return result
}

func mergeTwoLists(l1, l2 *ListNode) *ListNode {
	dummy := &ListNode{}
	cur := dummy
	for l1 != nil && l2 != nil {
		if l1.Val <= l2.Val {
			cur.Next = l1
			l1 = l1.Next
		} else {
			cur.Next = l2
			l2 = l2.Next
		}
		cur = cur.Next
	}
	if l1 != nil {
		cur.Next = l1
	} else {
		cur.Next = l2
	}
	return dummy.Next
}

// ── Approach 3: Divide and Conquer ───────────────────────────────────────────
//
// divideAndConquer pairs up lists and merges them in rounds until one remains.
//
// Intuition:
//   Instead of merging lists one at a time (which accumulates O(k) overhead),
//   pair them and merge in parallel rounds — like merge sort's combine step.
//   Round 1: merge pairs → k/2 lists. Round 2: → k/4 lists. Total: log k rounds.
//   Each round touches every node once: O(N) per round × O(log k) rounds = O(N log k).
//
// Time:  O(N log k) — optimal for this problem.
// Space: O(log k) — recursion depth of the divide step.
func divideAndConquer(lists []*ListNode) *ListNode {
	if len(lists) == 0 {
		return nil
	}
	return dnc(lists, 0, len(lists)-1)
}

func dnc(lists []*ListNode, lo, hi int) *ListNode {
	if lo == hi {
		return lists[lo]
	}
	mid := (lo + hi) / 2
	left := dnc(lists, lo, mid)
	right := dnc(lists, mid+1, hi)
	return mergeTwoLists(left, right)
}

// ── Approach 4: Min-Heap (Priority Queue) ─────────────────────────────────────
//
// minHeap uses a priority queue holding one node per list, always extracting
// the globally minimum node.
//
// Intuition:
//   Maintain a min-heap of at most k nodes (one per list head).
//   Pop the minimum → add to result. Push its successor onto the heap.
//   Each of the N nodes is pushed and popped once: O(N log k).
//
// Time:  O(N log k) — same as divide and conquer.
// Space: O(k) — the heap holds at most k nodes.
func minHeap(lists []*ListNode) *ListNode {
	h := &nodeHeap{}
	heap.Init(h)
	for _, node := range lists {
		if node != nil {
			heap.Push(h, node)
		}
	}
	dummy := &ListNode{}
	cur := dummy
	for h.Len() > 0 {
		node := heap.Pop(h).(*ListNode)
		cur.Next = node
		cur = cur.Next
		if node.Next != nil {
			heap.Push(h, node.Next)
		}
	}
	return dummy.Next
}

// nodeHeap implements heap.Interface for *ListNode.
type nodeHeap []*ListNode

func (h nodeHeap) Len() int            { return len(h) }
func (h nodeHeap) Less(i, j int) bool  { return h[i].Val < h[j].Val }
func (h nodeHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *nodeHeap) Push(x interface{}) { *h = append(*h, x.(*ListNode)) }
func (h *nodeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
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
	type example struct {
		lists  [][]int
		expect []int
	}
	examples := []example{
		{[][]int{{1, 4, 5}, {1, 3, 4}, {2, 6}}, []int{1, 1, 2, 3, 4, 4, 5, 6}},
		{[][]int{}, nil},
		{[][]int{{}}, nil},
		{[][]int{{1}, {0}}, []int{0, 1}},
	}

	approaches := []struct {
		name string
		fn   func([]*ListNode) *ListNode
	}{
		{"Approach 1: Brute Force (sort all)    O(N log N) T | O(N)    S", bruteForce},
		{"Approach 2: Sequential Merge          O(kN)      T | O(1)    S", sequentialMerge},
		{"Approach 3: Divide and Conquer      ✅ O(N log k) T | O(log k) S", divideAndConquer},
		{"Approach 4: Min-Heap               ✅ O(N log k) T | O(k)    S", minHeap},
	}

	for _, ex := range examples {
		fmt.Printf("lists=%v  expect=%v\n", ex.lists, ex.expect)
		for _, ap := range approaches {
			lists := make([]*ListNode, len(ex.lists))
			for i, v := range ex.lists {
				lists[i] = makeList(v)
			}
			fmt.Printf("  %-65s → %v\n", ap.name, listToSlice(ap.fn(lists)))
		}
		fmt.Println()
	}
}
