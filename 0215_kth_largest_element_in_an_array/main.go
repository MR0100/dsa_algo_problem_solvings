package main

import (
	"container/heap"
	"fmt"
	"sort"
)

// ── Approach 1: Sort then Index (Brute Force) ────────────────────────────────
//
// sortIndex solves Kth Largest by sorting a copy ascending and indexing the
// k-th element from the end.
//
// Intuition:
//
//	The k-th largest is a positional statistic. Sorting puts every element in
//	order, after which the k-th largest sits at index n-k. Dead simple and
//	obviously correct; the cost is the full O(n log n) sort even though we only
//	need one position.
//
// Algorithm:
//
//  1. Copy nums (avoid mutating the caller) and sort ascending.
//  2. Return copy[n-k].
//
// Time:  O(n log n) — dominated by the sort.
// Space: O(n) — the sorted copy.
func sortIndex(nums []int, k int) int {
	cp := append([]int(nil), nums...) // copy so we don't mutate input
	sort.Ints(cp)                     // ascending order
	return cp[len(cp)-k]              // k-th largest = k-th from the end
}

// ── Approach 2: Min-Heap of Size k ───────────────────────────────────────────
//
// minHeapK solves Kth Largest by keeping a min-heap of the k largest elements
// seen so far; its root is the answer.
//
// Intuition:
//
//	We only care about the k biggest values. Maintain a min-heap capped at size
//	k: the smallest of the "current top k" sits at the root. For each new value,
//	push it; if the heap exceeds k, pop the smallest. After processing all
//	elements the heap holds exactly the k largest, and its root (the minimum of
//	those) is the k-th largest overall. Great when k ≪ n or for streaming data.
//
// Algorithm:
//
//  1. Push each number; whenever heap size > k, pop the minimum.
//  2. Return the heap root.
//
// Time:  O(n log k) — n pushes/pops on a heap of size ≤ k.
// Space: O(k) for the heap.
type intHeap []int

func (h intHeap) Len() int            { return len(h) }
func (h intHeap) Less(i, j int) bool  { return h[i] < h[j] } // min-heap
func (h intHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *intHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *intHeap) Pop() interface{} {
	old := *h
	n := len(old)
	v := old[n-1]
	*h = old[:n-1]
	return v
}

func minHeapK(nums []int, k int) int {
	h := &intHeap{}
	heap.Init(h)
	for _, v := range nums {
		heap.Push(h, v) // add the new value
		if h.Len() > k {
			heap.Pop(h) // drop the smallest → keep only the k largest
		}
	}
	return (*h)[0] // root = minimum of the k largest = k-th largest
}

// ── Approach 3: Quickselect (Optimal, average O(n)) ──────────────────────────
//
// quickselect solves Kth Largest by partitioning around a pivot (Hoare/Lomuto
// style) and recursing only into the side that contains the target index.
//
// Intuition:
//
//	Full sorting is wasteful: we only need the element at index target = n-k in
//	ascending order. Quickselect partitions the array so the pivot lands at its
//	final sorted position p; everything left is ≤ pivot, right is ≥ pivot. If
//	p == target we're done; otherwise recurse into just the half containing
//	target. On average each partition halves the work → O(n).
//
// Algorithm:
//
//  1. target = n - k (index of the k-th largest in ascending order).
//  2. partition(lo, hi): Lomuto scheme with pivot = a[hi]; returns pivot's
//     final index p.
//  3. If p == target return a[p]; if p < target search right; else left.
//
// Time:  O(n) average, O(n²) worst case (pathological pivots).
// Space: O(1) extra (in-place; iterative loop, no recursion stack growth here).
func quickselect(nums []int, k int) int {
	a := append([]int(nil), nums...) // copy so we don't mutate the caller
	target := len(a) - k             // ascending index of the k-th largest
	lo, hi := 0, len(a)-1
	for lo <= hi {
		p := partition(a, lo, hi) // place a[hi] at its sorted position p
		switch {
		case p == target:
			return a[p] // pivot is exactly the element we want
		case p < target:
			lo = p + 1 // target is to the right of the pivot
		default:
			hi = p - 1 // target is to the left of the pivot
		}
	}
	return -1 // unreachable for valid 1 <= k <= n
}

// partition uses the Lomuto scheme: pivot = a[hi]; after the loop every element
// left of the returned index is ≤ pivot and the pivot sits at that index.
func partition(a []int, lo, hi int) int {
	pivot := a[hi] // choose the last element as pivot
	i := lo        // boundary: a[lo..i) holds values ≤ pivot
	for j := lo; j < hi; j++ {
		if a[j] <= pivot {
			a[i], a[j] = a[j], a[i] // grow the ≤-pivot region
			i++
		}
	}
	a[i], a[hi] = a[hi], a[i] // drop the pivot into its final slot
	return i
}

func main() {
	fmt.Println("=== Approach 1: Sort then Index (Brute Force) ===")
	fmt.Println(sortIndex([]int{3, 2, 1, 5, 6, 4}, 2))          // 5
	fmt.Println(sortIndex([]int{3, 2, 3, 1, 2, 4, 5, 5, 6}, 4)) // 4

	fmt.Println("=== Approach 2: Min-Heap of Size k ===")
	fmt.Println(minHeapK([]int{3, 2, 1, 5, 6, 4}, 2))          // 5
	fmt.Println(minHeapK([]int{3, 2, 3, 1, 2, 4, 5, 5, 6}, 4)) // 4

	fmt.Println("=== Approach 3: Quickselect (Optimal) ===")
	fmt.Println(quickselect([]int{3, 2, 1, 5, 6, 4}, 2))          // 5
	fmt.Println(quickselect([]int{3, 2, 3, 1, 2, 4, 5, 5, 6}, 4)) // 4
}
