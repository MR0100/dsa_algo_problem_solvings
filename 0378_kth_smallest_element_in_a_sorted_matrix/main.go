package main

import (
	"container/heap"
	"fmt"
	"sort"
)

// ── Approach 1: Brute Force (Flatten and Sort) ───────────────────────────────
//
// bruteForce solves Kth Smallest in a Sorted Matrix by copying all n² elements
// into a slice, sorting it, and indexing the k-th.
//
// Intuition:
//
//	Ignore the sortedness entirely: collect every value, sort ascending, and read
//	position k-1. Simple and obviously correct; the baseline the smarter methods
//	improve on.
//
// Algorithm:
//  1. Flatten the matrix into a slice of n² numbers.
//  2. Sort ascending.
//  3. Return element at index k-1.
//
// Time:  O(n² log n) — sorting n² values.
// Space: O(n²) — the flattened slice.
func bruteForce(matrix [][]int, k int) int {
	var flat []int
	for _, row := range matrix {
		flat = append(flat, row...) // gather every element
	}
	sort.Ints(flat)  // ascending order
	return flat[k-1] // k is 1-based
}

// minHeapItem tracks a value together with its matrix coordinates so we can push
// the "next in its row" after popping.
type minHeapItem struct {
	val, row, col int
}

// itemHeap is a min-heap of minHeapItem ordered by val.
type itemHeap []minHeapItem

func (h itemHeap) Len() int            { return len(h) }
func (h itemHeap) Less(i, j int) bool  { return h[i].val < h[j].val }
func (h itemHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *itemHeap) Push(x interface{}) { *h = append(*h, x.(minHeapItem)) }
func (h *itemHeap) Pop() interface{} {
	old := *h
	n := len(old)
	it := old[n-1]
	*h = old[:n-1]
	return it
}

// ── Approach 2: Min-Heap (K-Way Merge) ───────────────────────────────────────
//
// minHeapMerge solves Kth Smallest by merging the sorted rows: seed the heap
// with each row's first element and pop k times, pushing the next in the row.
//
// Intuition:
//
//	Each row is already sorted, so this is a k-way merge of n sorted lists. The
//	global smallest not-yet-emitted value is always the min across the current
//	front of each row — exactly what a min-heap of size ≤ n gives us. Pop the
//	smallest k times; the k-th pop is the answer.
//
// Algorithm:
//  1. Push (matrix[r][0], r, 0) for every row r.
//  2. Pop the min; if it has a right neighbour in its row, push that neighbour.
//  3. After k pops, return the last popped value.
//
// Time:  O(k log n) — k pops, each heap op O(log n) (heap holds ≤ n items).
// Space: O(n) — one entry per row.
func minHeapMerge(matrix [][]int, k int) int {
	n := len(matrix)
	h := &itemHeap{}
	for r := 0; r < n; r++ {
		heap.Push(h, minHeapItem{val: matrix[r][0], row: r, col: 0}) // row fronts
	}
	var popped minHeapItem
	for i := 0; i < k; i++ {
		popped = heap.Pop(h).(minHeapItem) // current global smallest
		if popped.col+1 < n {
			// advance within the popped item's row and offer that candidate
			heap.Push(h, minHeapItem{
				val: matrix[popped.row][popped.col+1],
				row: popped.row,
				col: popped.col + 1,
			})
		}
	}
	return popped.val // the k-th smallest overall
}

// ── Approach 3: Binary Search on Value (Optimal) ─────────────────────────────
//
// binarySearchValue solves Kth Smallest by binary-searching the answer VALUE in
// [matrix[0][0], matrix[n-1][n-1]] and counting how many elements are ≤ mid.
//
// Intuition:
//
//	The answer is some integer between the top-left (min) and bottom-right (max).
//	For a candidate value `mid`, countLessEqual(mid) — how many matrix entries are
//	≤ mid — is monotonically non-decreasing in mid. We want the smallest value
//	whose count is ≥ k; that value is guaranteed to be a matrix element. Counting
//	is O(n): start at the bottom-left corner and step right (count a column) when
//	the value ≤ mid, else step up.
//
// Algorithm:
//  1. lo = matrix[0][0], hi = matrix[n-1][n-1].
//  2. While lo < hi: mid = (lo+hi)/2; if countLessEqual(mid) < k, lo = mid+1
//     else hi = mid.
//  3. Return lo.
//
// Time:  O(n log(max-min)) — each of log(range) steps counts in O(n).
// Space: O(1).
func binarySearchValue(matrix [][]int, k int) int {
	n := len(matrix)
	lo, hi := matrix[0][0], matrix[n-1][n-1] // value range
	// countLessEqual returns the number of matrix entries ≤ target, walking from
	// the bottom-left corner in O(n) using both sortings.
	countLessEqual := func(target int) int {
		count := 0
		row, col := n-1, 0 // start bottom-left
		for row >= 0 && col < n {
			if matrix[row][col] <= target {
				count += row + 1 // this whole column up to `row` is ≤ target
				col++            // move right to a larger column
			} else {
				row-- // too big; move up to a smaller value
			}
		}
		return count
	}
	for lo < hi {
		mid := lo + (hi-lo)/2
		if countLessEqual(mid) < k {
			lo = mid + 1 // not enough values ≤ mid; the answer is larger
		} else {
			hi = mid // enough values; answer is mid or smaller
		}
	}
	return lo // smallest value with at least k entries ≤ it
}

func main() {
	m1 := [][]int{{1, 5, 9}, {10, 11, 13}, {12, 13, 15}}
	m2 := [][]int{{-5}}

	fmt.Println("=== Approach 1: Brute Force (Flatten and Sort) ===")
	fmt.Println(bruteForce(m1, 8)) // expected 13
	fmt.Println(bruteForce(m2, 1)) // expected -5

	fmt.Println("=== Approach 2: Min-Heap (K-Way Merge) ===")
	fmt.Println(minHeapMerge(m1, 8)) // expected 13
	fmt.Println(minHeapMerge(m2, 1)) // expected -5

	fmt.Println("=== Approach 3: Binary Search on Value (Optimal) ===")
	fmt.Println(binarySearchValue(m1, 8)) // expected 13
	fmt.Println(binarySearchValue(m2, 1)) // expected -5
}
