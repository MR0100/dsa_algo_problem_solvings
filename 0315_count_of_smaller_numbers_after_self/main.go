package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Brute Force (Count Pairs Directly) ────────────────────────────
//
// bruteForce solves Count of Smaller Numbers After Self by, for each index i,
// scanning everything to its right and counting strictly-smaller values.
//
// Intuition:
//
//	The definition is literal: counts[i] = number of j > i with nums[j] < nums[i].
//	Just check every such pair. It is obviously correct and a good baseline.
//
// Algorithm:
//
//  1. For each i, initialise c = 0.
//  2. For each j > i, if nums[j] < nums[i], increment c.
//  3. counts[i] = c.
//
// Time:  O(n^2) — every pair (i, j) with j > i is examined.
// Space: O(1) extra (excluding output).
func bruteForce(nums []int) []int {
	n := len(nums)
	counts := make([]int, n)
	for i := 0; i < n; i++ {
		c := 0
		for j := i + 1; j < n; j++ {
			if nums[j] < nums[i] { // strictly smaller element to the right
				c++
			}
		}
		counts[i] = c
	}
	return counts
}

// ── Approach 2: Binary Indexed Tree (Fenwick) over Ranks ──────────────────────
//
// fenwick solves Count of Smaller Numbers After Self by scanning right-to-left
// and, for each value, asking a Fenwick tree "how many already-seen values are
// strictly smaller than me?", then inserting the current value.
//
// Intuition:
//
//	Process from the right so "already inserted" == "to my right". A Fenwick tree
//	over compressed value-ranks supports two O(log n) operations: prefix-sum
//	(count of inserted values with rank < r) and point-update (insert a rank).
//	Coordinate-compress values to ranks 1..m to bound the tree size.
//
// Algorithm:
//
//  1. Compress values to ranks (sorted distinct order), rank starting at 1.
//  2. Iterate i from n-1 down to 0: counts[i] = query(rank[i]-1) = #inserted
//     values with smaller rank; then update(rank[i]) to insert nums[i].
//
// Time:  O(n log n) — n updates and queries, each O(log n).
// Space: O(n) — the Fenwick tree and the rank map.
func fenwick(nums []int) []int {
	n := len(nums)
	counts := make([]int, n)

	// Coordinate-compress: sorted distinct values -> rank (1-based).
	sorted := make([]int, n)
	copy(sorted, nums)
	sort.Ints(sorted)
	rank := map[int]int{}
	r := 0
	for _, v := range sorted {
		if _, ok := rank[v]; !ok {
			r++
			rank[v] = r // smallest distinct value gets rank 1
		}
	}

	// Fenwick tree (1-indexed) counting inserted values by rank.
	tree := make([]int, r+1)
	update := func(i int) { // add 1 at position i
		for ; i <= r; i += i & (-i) {
			tree[i]++
		}
	}
	query := func(i int) int { // prefix sum over positions 1..i
		s := 0
		for ; i > 0; i -= i & (-i) {
			s += tree[i]
		}
		return s
	}

	for i := n - 1; i >= 0; i-- {
		ri := rank[nums[i]]
		counts[i] = query(ri - 1) // #already-inserted (to the right) with smaller rank
		update(ri)                // now insert nums[i]
	}
	return counts
}

// ── Approach 3: Merge Sort (Count Inversions, Optimal) ────────────────────────
//
// mergeSortCount solves Count of Smaller Numbers After Self by an index-based
// merge sort: while merging two sorted halves, when a right-half element is
// placed before a left-half element, every remaining left element gains one
// "smaller to the right".
//
// Intuition:
//
//	counts[i] is the number of inversions (i, j), j > i, nums[j] < nums[i].
//	Merge sort naturally counts inversions. We sort INDICES (not values) so each
//	original index keeps its own tally. When merging, if the right candidate is
//	strictly smaller than the left candidate, it is smaller-and-to-the-right of
//	the left candidate and of every left element still waiting — but we credit it
//	lazily: each time we PLACE a left index, we add how many right elements have
//	already been merged ahead of it.
//
// Algorithm:
//
//  1. idx = [0..n-1]; counts = zeros.
//  2. Recursively sort idx by nums value. In merge, keep a counter of right-half
//     elements already merged; when we take a LEFT index, add that counter to
//     its count (those right elements are smaller and were to its right).
//  3. Return counts.
//
// Time:  O(n log n) — merge sort.
// Space: O(n) — index arrays and recursion.
func mergeSortCount(nums []int) []int {
	n := len(nums)
	counts := make([]int, n)
	idx := make([]int, n) // indices, reordered by value during the sort
	for i := range idx {
		idx[i] = i
	}

	var sortRange func(lo, hi int)
	sortRange = func(lo, hi int) {
		if hi-lo <= 1 {
			return // single element (or empty) is already sorted
		}
		mid := (lo + hi) / 2
		sortRange(lo, mid)
		sortRange(mid, hi)

		merged := make([]int, 0, hi-lo)
		i, j := lo, mid
		rightMerged := 0 // count of right-half indices already placed
		for i < mid && j < hi {
			if nums[idx[j]] < nums[idx[i]] {
				// Right value strictly smaller -> place it; it sits to the right
				// of every not-yet-placed left index. Credit given when left is placed.
				rightMerged++
				merged = append(merged, idx[j])
				j++
			} else {
				// Left value <= right value: idx[i] is placed now. Exactly
				// rightMerged right-half elements were smaller AND to its right.
				counts[idx[i]] += rightMerged
				merged = append(merged, idx[i])
				i++
			}
		}
		for i < mid { // remaining left indices: all rightMerged right elems were smaller
			counts[idx[i]] += rightMerged
			merged = append(merged, idx[i])
			i++
		}
		for j < hi { // remaining right indices: nothing to credit
			merged = append(merged, idx[j])
			j++
		}
		copy(idx[lo:hi], merged) // write the sorted order back
	}
	sortRange(0, n)
	return counts
}

func main() {
	// Official Example 1
	e1 := []int{5, 2, 6, 1}
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(bruteForce(e1)) // expected [2 1 1 0]
	fmt.Println("=== Approach 2: Fenwick (BIT) ===")
	fmt.Println(fenwick(e1)) // expected [2 1 1 0]
	fmt.Println("=== Approach 3: Merge Sort (Optimal) ===")
	fmt.Println(mergeSortCount(e1)) // expected [2 1 1 0]

	// Official Example 2
	e2 := []int{-1}
	fmt.Println("=== Approach 1: Brute Force (Example 2) ===")
	fmt.Println(bruteForce(e2)) // expected [0]
	fmt.Println("=== Approach 2: Fenwick (Example 2) ===")
	fmt.Println(fenwick(e2)) // expected [0]
	fmt.Println("=== Approach 3: Merge Sort (Example 2) ===")
	fmt.Println(mergeSortCount(e2)) // expected [0]

	// Official Example 3
	e3 := []int{-1, -1}
	fmt.Println("=== Approach 1: Brute Force (Example 3) ===")
	fmt.Println(bruteForce(e3)) // expected [0 0]
	fmt.Println("=== Approach 2: Fenwick (Example 3) ===")
	fmt.Println(fenwick(e3)) // expected [0 0]
	fmt.Println("=== Approach 3: Merge Sort (Example 3) ===")
	fmt.Println(mergeSortCount(e3)) // expected [0 0]
}
