package main

import (
	"container/heap"
	"fmt"
	"math"
	"sort"
)

// Closest Binary Search Tree Value II (LeetCode #272)
//
// Given the root of a binary search tree, a target value, and an integer k,
// return the k values in the BST that are closest to target. The answer may be
// returned in any order; it is guaranteed to be unique.
//
// TreeNode is the standard binary-tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// ── Approach 1: Inorder + Sort by Distance (Brute Force) ─────────────────────
//
// bruteForce collects all node values, sorts them by |val - target|, and keeps
// the first k.
//
// Intuition:
//
//	"k closest to target" is a selection over ALL values. The most direct
//	route: gather every value, sort by distance to target, take k. Ignores
//	the BST structure entirely — a correctness-first baseline.
//
// Algorithm:
//  1. Inorder-traverse to collect all values into a slice.
//  2. Sort the slice by ascending |v - target|.
//  3. Return the first k entries.
//
// Time:  O(n log n) — the sort dominates.
// Space: O(n) — the collected values.
func bruteForce(root *TreeNode, target float64, k int) []int {
	var vals []int
	var inorder func(*TreeNode)
	inorder = func(n *TreeNode) {
		if n == nil {
			return
		}
		inorder(n.Left)            // left subtree first
		vals = append(vals, n.Val) // visit node
		inorder(n.Right)           // then right subtree
	}
	inorder(root)
	// Sort so that the closest-to-target values come first.
	sort.Slice(vals, func(i, j int) bool {
		return math.Abs(float64(vals[i])-target) < math.Abs(float64(vals[j])-target)
	})
	return vals[:k] // the k nearest
}

// ── Approach 2: Inorder + Sliding Window on Sorted Array ──────────────────────
//
// slidingWindow uses the fact that inorder traversal of a BST yields a SORTED
// array, then slides a size-k window to the position that minimises total
// distance to target.
//
// Intuition:
//
//	In a sorted array, the k closest values to a target always form a
//	CONTIGUOUS window of length k. Start with the window [0, k-1]; while the
//	element just past the right end is closer to target than the element at
//	the left end, slide the window right by one. When sliding no longer helps,
//	the window holds the k closest values.
//
// Algorithm:
//  1. Inorder-traverse to get a sorted slice `vals`.
//  2. lo = 0, hi = len(vals) - 1. Shrink [lo, hi] to width k:
//     while hi - lo >= k, drop whichever end is farther from target.
//  3. Return vals[lo : hi+1].
//
// Time:  O(n) — traversal is O(n); shrinking removes n-k elements, each O(1).
// Space: O(n) — the sorted array (O(k) extra beyond it).
func slidingWindow(root *TreeNode, target float64, k int) []int {
	var vals []int
	var inorder func(*TreeNode)
	inorder = func(n *TreeNode) {
		if n == nil {
			return
		}
		inorder(n.Left)
		vals = append(vals, n.Val)
		inorder(n.Right)
	}
	inorder(root) // vals is now sorted ascending

	lo, hi := 0, len(vals)-1
	// Repeatedly discard the endpoint farther from target until exactly k remain.
	for hi-lo >= k {
		if math.Abs(float64(vals[lo])-target) > math.Abs(float64(vals[hi])-target) {
			lo++ // left end is farther → drop it
		} else {
			hi-- // right end is farther (or tie) → drop it
		}
	}
	return vals[lo : hi+1] // the k contiguous nearest values
}

// ── Approach 3: Max-Heap of Size k (Optimal for k ≪ n) ───────────────────────
//
// maxHeapK does a single inorder pass, maintaining a max-heap of size k keyed
// by distance to target, so we never store more than k candidates.
//
// Intuition:
//
//	If k is much smaller than n, sorting all n values is wasteful. Keep only
//	the k best seen so far in a max-heap keyed by distance: the heap's root is
//	the WORST of the current best k. For each new value, if the heap has room,
//	push it; otherwise, if the new value beats the current worst, evict the
//	worst and push the new one.
//
// Algorithm:
//  1. Inorder-traverse. For each value v:
//     - push (distance, v) onto the heap;
//     - if heap size > k, pop the max (the farthest of the current best k).
//  2. Drain the heap into the answer.
//
// Time:  O(n log k) — n pushes/pops each O(log k).
// Space: O(k) — the heap.
func maxHeapK(root *TreeNode, target float64, k int) []int {
	h := &distHeap{}
	heap.Init(h)
	var inorder func(*TreeNode)
	inorder = func(n *TreeNode) {
		if n == nil {
			return
		}
		inorder(n.Left)
		// Push this value keyed by its distance to target.
		heap.Push(h, item{dist: math.Abs(float64(n.Val) - target), val: n.Val})
		// If we now hold more than k, evict the farthest (heap root).
		if h.Len() > k {
			heap.Pop(h)
		}
		inorder(n.Right)
	}
	inorder(root)
	// Whatever remains in the heap are the k closest values.
	res := make([]int, 0, k)
	for h.Len() > 0 {
		res = append(res, heap.Pop(h).(item).val)
	}
	return res
}

// item is one candidate: a BST value with its distance to target.
type item struct {
	dist float64
	val  int
}

// distHeap is a MAX-heap by distance (root = farthest), so we can evict the
// worst candidate in O(log k).
type distHeap []item

func (h distHeap) Len() int            { return len(h) }
func (h distHeap) Less(i, j int) bool  { return h[i].dist > h[j].dist } // '>' → max-heap
func (h distHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *distHeap) Push(x interface{}) { *h = append(*h, x.(item)) }
func (h *distHeap) Pop() interface{} {
	old := *h
	n := len(old)
	it := old[n-1]
	*h = old[:n-1]
	return it
}

// asMultiset compares two int slices ignoring order (the answer may be in any
// order), so main() can print stable PASS/FAIL results.
func asMultiset(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	ca, cb := append([]int(nil), a...), append([]int(nil), b...)
	sort.Ints(ca)
	sort.Ints(cb)
	for i := range ca {
		if ca[i] != cb[i] {
			return false
		}
	}
	return true
}

func main() {
	// Example 1 tree:      4
	//                    /   \
	//                   2     5
	//                  / \
	//                 1   3
	// root = [4,2,5,1,3], target = 3.714286, k = 2  → [4,3]
	root := &TreeNode{Val: 4,
		Left:  &TreeNode{Val: 2, Left: &TreeNode{Val: 1}, Right: &TreeNode{Val: 3}},
		Right: &TreeNode{Val: 5}}
	target1, k1 := 3.714286, 2
	want1 := []int{4, 3}

	// Example 2: root = [1], target = 0.000000, k = 1 → [1]
	root2 := &TreeNode{Val: 1}
	target2, k2 := 0.000000, 1
	want2 := []int{1}

	fmt.Println("=== Approach 1: Inorder + Sort by Distance ===")
	fmt.Println(asMultiset(bruteForce(root, target1, k1), want1))  // expected true
	fmt.Println(asMultiset(bruteForce(root2, target2, k2), want2)) // expected true

	fmt.Println("=== Approach 2: Sliding Window on Sorted Array ===")
	fmt.Println(asMultiset(slidingWindow(root, target1, k1), want1))  // expected true
	fmt.Println(asMultiset(slidingWindow(root2, target2, k2), want2)) // expected true

	fmt.Println("=== Approach 3: Max-Heap of Size k ===")
	fmt.Println(asMultiset(maxHeapK(root, target1, k1), want1))  // expected true
	fmt.Println(asMultiset(maxHeapK(root2, target2, k2), want2)) // expected true
}
