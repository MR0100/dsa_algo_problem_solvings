package main

import (
	"container/heap"
	"fmt"
	"sort"
)

// ── Approach 1: Brute Force (Generate All + Sort) ────────────────────────────
//
// bruteForce forms every pair (nums1[i], nums2[j]), sorts them all by sum, and
// returns the first k.
//
// Intuition:
//
//	The straightforward reading of the problem: enumerate the full m*n grid of
//	pairs, sort by sum, take the k smallest. Correct but wasteful when the grid
//	is large and k is small.
//
// Algorithm:
//  1. Build every pair [u, v] for u in nums1, v in nums2.
//  2. Stable-sort the list by u+v.
//  3. Return the first min(k, len) pairs.
//
// Time:  O(m*n log(m*n)) — building and sorting the whole grid.
// Space: O(m*n) — the list of all pairs.
func bruteForce(nums1, nums2 []int, k int) [][]int {
	pairs := make([][]int, 0, len(nums1)*len(nums2))
	for _, u := range nums1 { // every element of the first array
		for _, v := range nums2 { // paired with every element of the second
			pairs = append(pairs, []int{u, v})
		}
	}
	// Sort ascending by the pair sum; ties keep their relative order.
	sort.SliceStable(pairs, func(i, j int) bool {
		return pairs[i][0]+pairs[i][1] < pairs[j][0]+pairs[j][1]
	})
	if k > len(pairs) { // never ask for more pairs than exist
		k = len(pairs)
	}
	return pairs[:k] // the k smallest-sum pairs
}

// pairItem is one entry on the min-heap: indices into nums1/nums2 and the sum.
type pairItem struct {
	i, j int // index into nums1 (i) and nums2 (j)
	sum  int // nums1[i] + nums2[j], the heap key
}

// minHeap orders pairItems by ascending sum.
type minHeap []pairItem

func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(a, b int) bool  { return h[a].sum < h[b].sum } // smallest sum first
func (h minHeap) Swap(a, b int)       { h[a], h[b] = h[b], h[a] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(pairItem)) }
func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}

// ── Approach 2: Min-Heap over the Sorted Frontier (Optimal) ──────────────────
//
// heapFrontier exploits the fact that both arrays are sorted. It grows a
// "frontier" of candidate pairs on a min-heap and pops the k smallest sums,
// expanding neighbours as it goes — never materialising the full grid.
//
// Intuition:
//
//	Think of pairs as a grid where row i uses nums1[i] and column j uses
//	nums2[j]. Sums increase down each row and across each column. The globally
//	smallest unused pair is always on the "staircase" boundary. Seed the heap
//	with the first column (i, 0) for each i (bounded by k), then each time we pop
//	(i, j) we push its right neighbour (i, j+1) — the only new candidate that
//	could now be minimal.
//
// Algorithm:
//  1. Push (i, 0) for i in 0..min(len(nums1), k)-1 with sum nums1[i]+nums2[0].
//  2. Repeat k times (while heap non-empty):
//     a. Pop the smallest (i, j); append [nums1[i], nums2[j]] to the answer.
//     b. If j+1 < len(nums2), push (i, j+1).
//  3. Return the collected pairs.
//
// Time:  O(k log k) — the heap never holds more than ~k items; k pops/pushes.
// Space: O(k) — heap plus output.
func heapFrontier(nums1, nums2 []int, k int) [][]int {
	res := make([][]int, 0, k)
	if len(nums1) == 0 || len(nums2) == 0 || k == 0 {
		return res // nothing to pair
	}

	h := &minHeap{}
	heap.Init(h)
	// Seed with the first column: pairs (nums1[i], nums2[0]). Only the first k
	// rows can ever contribute to the k smallest sums.
	limit := len(nums1)
	if limit > k {
		limit = k
	}
	for i := 0; i < limit; i++ {
		heap.Push(h, pairItem{i: i, j: 0, sum: nums1[i] + nums2[0]})
	}

	// Pop k smallest sums, expanding the row's next column each time.
	for h.Len() > 0 && len(res) < k {
		it := heap.Pop(h).(pairItem)                       // current smallest-sum pair
		res = append(res, []int{nums1[it.i], nums2[it.j]}) // record it
		if it.j+1 < len(nums2) {                           // slide right within the same row
			heap.Push(h, pairItem{i: it.i, j: it.j + 1, sum: nums1[it.i] + nums2[it.j+1]})
		}
	}
	return res
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(bruteForce([]int{1, 7, 11}, []int{2, 4, 6}, 3)) // expected [[1 2] [1 4] [1 6]]
	fmt.Println(bruteForce([]int{1, 1, 2}, []int{1, 2, 3}, 2))  // expected [[1 1] [1 1]]
	fmt.Println(bruteForce([]int{1, 2}, []int{3}, 3))           // expected [[1 3] [2 3]]

	fmt.Println("=== Approach 2: Min-Heap Frontier ===")
	fmt.Println(heapFrontier([]int{1, 7, 11}, []int{2, 4, 6}, 3)) // expected [[1 2] [1 4] [1 6]]
	fmt.Println(heapFrontier([]int{1, 1, 2}, []int{1, 2, 3}, 2))  // expected [[1 1] [1 1]]
	fmt.Println(heapFrontier([]int{1, 2}, []int{3}, 3))           // expected [[1 3] [2 3]]
}
