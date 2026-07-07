package main

import (
	"fmt"
	"sort"
)

// A "reverse pair" is (i, j) with i < j and nums[i] > 2*nums[j]. We count them.
// Note: 2*nums[j] can overflow int32; Go's int is 64-bit on modern platforms,
// so plain int arithmetic is safe here (values fit in [-2^31, 2^31-1]).

// ── Approach 1: Brute Force (All Pairs) ──────────────────────────────────────
//
// bruteForce checks every ordered pair (i, j) with i < j directly.
//
// Intuition:
//
//	The definition is a double loop: for each earlier index i and later index
//	j, test nums[i] > 2*nums[j]. No cleverness — just count the ones that hold.
//	This is the ground truth we validate the fast solutions against.
//
// Algorithm:
//  1. count = 0.
//  2. For i in 0..n-1, for j in i+1..n-1: if nums[i] > 2*nums[j], count++.
//  3. Return count.
//
// Time:  O(n^2) — every pair examined. TLE for n up to 5e4 (2.5e9 pairs).
// Space: O(1).
func bruteForce(nums []int) int {
	count := 0
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			// int is 64-bit here, so 2*nums[j] cannot overflow the range.
			if nums[i] > 2*nums[j] {
				count++ // (i, j) is a reverse pair
			}
		}
	}
	return count
}

// ── Approach 2: Merge Sort (Divide and Conquer, Optimal) ─────────────────────
//
// mergeSortCount counts reverse pairs while sorting, by exploiting that within
// each merge step the left and right halves are already individually sorted.
//
// Intuition:
//
//	Split the array in half. Reverse pairs are either fully inside the left
//	half, fully inside the right half, or "cross" (i in left, j in right).
//	Recurse to count the within-half pairs; then, because both halves are
//	sorted, the crossing pairs can be counted in a single linear sweep: for a
//	rising left pointer i, advance a right pointer j while left[i] > 2*right[j];
//	every such j pairs with THIS i and all larger left values (they are even
//	bigger). Finally merge the two sorted halves so the parent can do the same.
//
// Algorithm:
//  1. If the segment has < 2 elements, return 0.
//  2. Recurse on left and right halves, summing their counts.
//  3. Count crossing pairs: two pointers i over left, j over right; for each i,
//     move j forward while left[i] > 2*left[right-relative]; add (j - midStart).
//  4. Standard merge of the two sorted halves back into place.
//  5. Return total count.
//
// Time:  O(n log n) — log n levels, each doing O(n) counting + merging.
// Space: O(n) — the temporary buffer used by the merge.
func mergeSortCount(nums []int) int {
	if len(nums) < 2 {
		return 0
	}
	arr := make([]int, len(nums)) // work on a copy so the input stays intact
	copy(arr, nums)
	tmp := make([]int, len(nums)) // reusable merge buffer
	return mergeCount(arr, tmp, 0, len(arr)-1)
}

// mergeCount sorts arr[lo..hi] in place and returns the reverse-pair count
// contained in that segment.
func mergeCount(arr, tmp []int, lo, hi int) int {
	if lo >= hi {
		return 0 // 0 or 1 element: no pairs
	}
	mid := lo + (hi-lo)/2
	// 1) Count pairs fully inside each half (and sort each half).
	count := mergeCount(arr, tmp, lo, mid) + mergeCount(arr, tmp, mid+1, hi)

	// 2) Count crossing pairs. Both arr[lo..mid] and arr[mid+1..hi] are sorted.
	//    For each i on the left, extend j on the right while arr[i] > 2*arr[j].
	j := mid + 1
	for i := lo; i <= mid; i++ {
		// arr grows with i, so j never needs to move backward (monotonic).
		for j <= hi && arr[i] > 2*arr[j] {
			j++
		}
		count += j - (mid + 1) // all right elements before j pair with arr[i]
	}

	// 3) Merge arr[lo..mid] and arr[mid+1..hi] into sorted order via tmp.
	i, k, r := lo, lo, mid+1
	for i <= mid && r <= hi {
		if arr[i] <= arr[r] {
			tmp[k] = arr[i]
			i++
		} else {
			tmp[k] = arr[r]
			r++
		}
		k++
	}
	for i <= mid { // drain leftovers from the left half
		tmp[k] = arr[i]
		i++
		k++
	}
	for r <= hi { // drain leftovers from the right half
		tmp[k] = arr[r]
		r++
		k++
	}
	copy(arr[lo:hi+1], tmp[lo:hi+1]) // write the merged run back in place
	return count
}

// ── Approach 3: Binary Indexed Tree + Coordinate Compression ─────────────────
//
// bitCount sweeps j from left to right, and for each j counts how many already-
// seen nums[i] (i < j) satisfy nums[i] > 2*nums[j], using a Fenwick tree over
// compressed values for fast prefix counts.
//
// Intuition:
//
//	Process indices left→right. When we reach j, every i < j has already been
//	inserted into a frequency structure keyed by value. A reverse pair needs
//	nums[i] > 2*nums[j], i.e. we want "how many inserted values exceed
//	2*nums[j]". A Fenwick tree answers prefix-sum queries in O(log n); to bound
//	the index space we coordinate-compress the union of all nums[i] and all
//	2*nums[j] to ranks 1..m. Then "count values > 2*nums[j]" = total inserted −
//	prefix count up to rank(2*nums[j]).
//
// Algorithm:
//  1. Build a sorted, de-duplicated list of every nums[i] and every 2*nums[j];
//     map each to a 1-based rank (coordinate compression).
//  2. For j = 0..n-1 (left to right):
//     a. r = number of inserted values with rank <= rank(2*nums[j]); the
//     count of values strictly greater is (insertedSoFar - r), which is
//     added to the answer.
//     b. Insert nums[j] by rank into the tree.
//  3. Return the accumulated answer.
//
// Time:  O(n log n) — n updates and n queries, each O(log n), plus the sort.
// Space: O(n) — the compressed value list and the Fenwick array.
func bitCount(nums []int) int {
	n := len(nums)
	if n < 2 {
		return 0
	}
	// 1) Collect every value we will ever query or insert, then compress.
	vals := make([]int, 0, 2*n)
	for _, v := range nums {
		vals = append(vals, v)   // values we INSERT (the nums[i])
		vals = append(vals, 2*v) // values we QUERY against (the 2*nums[j])
	}
	sort.Ints(vals)
	uniq := vals[:0] // dedup in place
	for i, v := range vals {
		if i == 0 || v != vals[i-1] {
			uniq = append(uniq, v)
		}
	}
	// rank returns the 1-based position of x in the sorted unique list.
	rank := func(x int) int {
		return sort.SearchInts(uniq, x) + 1 // +1 → Fenwick indices start at 1
	}

	tree := make([]int, len(uniq)+1) // Fenwick tree of value frequencies
	// update adds 1 at position i.
	update := func(i int) {
		for ; i < len(tree); i += i & (-i) {
			tree[i]++
		}
	}
	// query returns the count of inserted values with rank in [1, i].
	query := func(i int) int {
		s := 0
		for ; i > 0; i -= i & (-i) {
			s += tree[i]
		}
		return s
	}

	answer := 0
	inserted := 0 // how many nums[i] are already in the tree (i < j)
	for j := 0; j < n; j++ {
		// Count already-inserted values strictly greater than 2*nums[j]:
		//   inserted - (# values with rank <= rank(2*nums[j])).
		r := query(rank(2 * nums[j]))
		answer += inserted - r
		update(rank(nums[j])) // now nums[j] becomes an eligible "i" for later j
		inserted++
	}
	return answer
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (All Pairs) ===")
	fmt.Println(bruteForce([]int{1, 3, 2, 3, 1})) // expected 2
	fmt.Println(bruteForce([]int{2, 4, 3, 5, 1})) // expected 3

	fmt.Println("=== Approach 2: Merge Sort (Divide and Conquer, Optimal) ===")
	fmt.Println(mergeSortCount([]int{1, 3, 2, 3, 1})) // expected 2
	fmt.Println(mergeSortCount([]int{2, 4, 3, 5, 1})) // expected 3

	fmt.Println("=== Approach 3: Binary Indexed Tree + Coordinate Compression ===")
	fmt.Println(bitCount([]int{1, 3, 2, 3, 1})) // expected 2
	fmt.Println(bitCount([]int{2, 4, 3, 5, 1})) // expected 3
}
