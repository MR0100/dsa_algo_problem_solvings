package main

import (
	"fmt"
	"sort"
)

// The three approaches all revolve around the same reformulation.
//
// Build a prefix-sum array P of length n+1 where P[0] = 0 and
// P[k] = nums[0] + nums[1] + ... + nums[k-1]. Then the range sum of
// nums[i..j] (0-indexed, i <= j) equals P[j+1] - P[i].
//
// So counting range sums in [lower, upper] is exactly counting pairs
// (a, b) with a < b and lower <= P[b] - P[a] <= upper.
//
// int64 is used for prefix sums everywhere: with up to 10^5 elements each as
// large as ~2^31, a prefix sum can reach ~2*10^14, which overflows int32 and,
// on a 32-bit platform, plain int too.

// ── Approach 1: Brute Force ──────────────────────────────────────────────────
//
// bruteForce solves Count of Range Sum by testing every (i, j) pair directly.
//
// Intuition:
//
//	Compute the prefix sums once. Every range sum S(i, j) is a single
//	subtraction P[j+1] - P[i]. Enumerate all pairs a < b of prefix-sum
//	indices, and for each pair check whether P[b] - P[a] falls in the window.
//	No cleverness — just the definition made cheap by prefix sums.
//
// Algorithm:
//  1. Build P[0..n] with P[0] = 0 and P[k] = P[k-1] + nums[k-1].
//  2. For every pair a < b, compute diff = P[b] - P[a].
//  3. If lower <= diff <= upper, increment the counter.
//  4. Return the counter.
//
// Time:  O(n^2) — every pair of the n+1 prefix indices is examined.
// Space: O(n)  — the prefix-sum array (excluding the O(1) counter).
func bruteForce(nums []int, lower int, upper int) int {
	n := len(nums)
	prefix := make([]int64, n+1) // prefix[0] = 0 already (zero value)
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i] + int64(nums[i]) // running prefix sum in int64
	}

	lo, hi := int64(lower), int64(upper) // widen bounds once to compare cleanly
	count := 0
	for a := 0; a <= n; a++ { // a is the "start" prefix index
		for b := a + 1; b <= n; b++ { // b > a is the "end" prefix index
			diff := prefix[b] - prefix[a] // this equals range sum S(a, b-1)
			if diff >= lo && diff <= hi { // inside the inclusive window?
				count++ // valid range sum found
			}
		}
	}
	return count
}

// ── Approach 2: Merge Sort Count (Divide and Conquer) ────────────────────────
//
// mergeSortCount solves Count of Range Sum with a modified merge sort over the
// prefix-sum array. This is the canonical O(n log n) divide-and-conquer answer.
//
// Intuition:
//
//	Sort the prefix sums recursively. During the merge of a sorted left half
//	and a sorted right half, every valid pair (a, b) we want to count has
//	a in the left half and b in the right half (a < b by construction because
//	left indices are all smaller than right indices). Because both halves are
//	sorted, for a fixed left value P[a] the set of right values P[b] with
//	lower <= P[b] - P[a] <= upper is a contiguous window; two monotone pointers
//	sweep that window in linear time as P[a] increases.
//
// Algorithm:
//  1. Build prefix sums P[0..n] in int64.
//  2. sort(lo, hi): if the block is trivial, return. Otherwise recurse on the
//     left [lo, mid] and right [mid+1, hi] halves.
//  3. After both halves are individually sorted, count cross pairs: for each
//     a in the left half, advance pointer `low` to the first right index whose
//     value is >= P[a] + lower, and pointer `high` to the first right index
//     whose value is > P[a] + upper. The gap (high - low) is the number of
//     right values inside [P[a]+lower, P[a]+upper], i.e. valid b's for this a.
//  4. Merge the two sorted halves back into a scratch buffer, then copy back.
//  5. Return the accumulated count.
//
// Time:  O(n log n) — standard merge-sort recursion; the counting sweep and the
//
//	merge are both linear per level, over log n levels.
//
// Space: O(n) — prefix array plus a scratch buffer for merging.
func mergeSortCount(nums []int, lower int, upper int) int {
	n := len(nums)
	prefix := make([]int64, n+1)
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i] + int64(nums[i]) // int64 prefix sums
	}

	lo, hi := int64(lower), int64(upper)
	scratch := make([]int64, len(prefix)) // reusable merge buffer

	// sortCount sorts prefix[left..right] in place and returns the number of
	// valid cross pairs discovered while merging its two halves.
	var sortCount func(left, right int) int
	sortCount = func(left, right int) int {
		if left >= right {
			return 0 // 0 or 1 element: nothing to pair, already sorted
		}
		mid := left + (right-left)/2 // split point (overflow-safe midpoint)
		// Count pairs fully inside each half first (both halves get sorted here).
		count := sortCount(left, mid) + sortCount(mid+1, right)

		// Count cross pairs: a in [left, mid], b in [mid+1, right].
		low, high := mid+1, mid+1 // moving window bounds into the right half
		for a := left; a <= mid; a++ {
			// low = first right index with prefix[low] - prefix[a] >= lower.
			for low <= right && prefix[low]-prefix[a] < lo {
				low++
			}
			// high = first right index with prefix[high] - prefix[a] > upper.
			for high <= right && prefix[high]-prefix[a] <= hi {
				high++
			}
			// [low, high) is the window of valid right partners for this a.
			count += high - low
		}

		// Standard merge of the two sorted halves into scratch, then copy back.
		i, j, k := left, mid+1, left
		for i <= mid && j <= right {
			if prefix[i] <= prefix[j] { // take the smaller front element
				scratch[k] = prefix[i]
				i++
			} else {
				scratch[k] = prefix[j]
				j++
			}
			k++
		}
		for i <= mid { // drain any remaining left elements
			scratch[k] = prefix[i]
			i++
			k++
		}
		for j <= right { // drain any remaining right elements
			scratch[k] = prefix[j]
			j++
			k++
		}
		copy(prefix[left:right+1], scratch[left:right+1]) // write sorted block back
		return count
	}

	return sortCount(0, n) // sort/count over all n+1 prefix indices
}

// ── Approach 3: Binary Indexed Tree (Fenwick, coordinate-compressed) ─────────
//
// binaryIndexedTree solves Count of Range Sum by streaming prefix sums into a
// Fenwick tree over compressed coordinates and querying a count-range per step.
//
// Intuition:
//
//	Process prefix indices b = 0, 1, ..., n in order. When we are about to
//	consider P[b] as the "end", every P[a] with a < b has already been inserted
//	into the tree. The number of previous prefix sums P[a] satisfying
//	lower <= P[b] - P[a] <= upper is the number of inserted values in the
//	closed interval [P[b] - upper, P[b] - lower]. A Fenwick tree gives that
//	count in O(log m) once all candidate values are coordinate-compressed to
//	dense ranks. Insert P[b] afterward so it becomes available to later ends.
//
// Algorithm:
//  1. Build int64 prefix sums P[0..n].
//  2. Collect every value that a query or an insert could reference:
//     each P[v], plus P[v]-lower and P[v]-upper. Sort and dedup them to get a
//     coordinate list; rank(x) = index of first value >= x (1-based).
//  3. For b = 0..n: let L = P[b] - upper, R = P[b] - lower. Add to the answer
//     the count of already-inserted values in [L, R] via two Fenwick prefix
//     queries: query(rank_upper(R)) - query(rank_upper(L)-1). Then insert P[b].
//  4. Return the accumulated answer.
//
// Time:  O(n log n) — sorting/compression is O(n log n); each of the n+1 steps
//
//	does O(log n) Fenwick work.
//
// Space: O(n) — the coordinate list and the Fenwick tree.
func binaryIndexedTree(nums []int, lower int, upper int) int {
	n := len(nums)
	prefix := make([]int64, n+1)
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i] + int64(nums[i]) // int64 prefix sums
	}
	lo, hi := int64(lower), int64(upper)

	// Gather every value that will ever be inserted or used as a query bound.
	all := make([]int64, 0, 3*(n+1))
	for _, p := range prefix {
		all = append(all, p)    // the value we insert
		all = append(all, p-lo) // R bound: P[b] - lower
		all = append(all, p-hi) // L bound: P[b] - upper
	}
	sort.Slice(all, func(i, j int) bool { return all[i] < all[j] }) // sort ascending
	// Dedup in place so equal values share one rank.
	uniq := all[:0:0] // fresh slice, cap 0, so appends don't alias `all`
	for i, v := range all {
		if i == 0 || v != all[i-1] {
			uniq = append(uniq, v)
		}
	}

	// rank returns the 1-based index of the first compressed value >= x.
	rank := func(x int64) int {
		return sort.Search(len(uniq), func(i int) bool { return uniq[i] >= x }) + 1
	}

	// Fenwick tree (1-indexed) storing counts of inserted compressed values.
	tree := make([]int, len(uniq)+1)
	update := func(i int) { // add 1 at position i
		for ; i < len(tree); i += i & (-i) { // climb via lowest-set-bit jumps
			tree[i]++
		}
	}
	query := func(i int) int { // prefix count over positions [1, i]
		s := 0
		for ; i > 0; i -= i & (-i) { // descend via lowest-set-bit jumps
			s += tree[i]
		}
		return s
	}

	count := 0
	for _, p := range prefix { // p plays the role of P[b], an "end" prefix sum
		L := p - hi // lower edge of the valid P[a] interval: P[b] - upper
		R := p - lo // upper edge of the valid P[a] interval: P[b] - lower
		// Count already-inserted values in [L, R].
		// rank(L) is the first value >= L; rank(R+1)-1 is the last value <= R.
		rightRank := rank(R+1) - 1 // last compressed index with value <= R
		leftRank := rank(L)        // first compressed index with value >= L
		count += query(rightRank) - query(leftRank-1)
		update(rank(p)) // now P[b] is available to future ends as a P[a]
	}
	return count
}

func main() {
	// Example 1: nums = [-2,5,-1], lower = -2, upper = 2  → 3
	ex1 := []int{-2, 5, -1}
	// Example 2: nums = [0], lower = 0, upper = 0          → 1
	ex2 := []int{0}
	// Edge:      nums = [2,-2,2,-2], lower = -1, upper = 1 → 5 (all pairs summing to 0)
	ex3 := []int{2, -2, 2, -2}

	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Printf("nums=[-2,5,-1], lower=-2, upper=2   got=%d  expected 3\n", bruteForce(ex1, -2, 2)) // expected 3
	fmt.Printf("nums=[0], lower=0, upper=0          got=%d  expected 1\n", bruteForce(ex2, 0, 0))  // expected 1
	fmt.Printf("nums=[2,-2,2,-2], lower=-1, upper=1 got=%d  expected 4\n", bruteForce(ex3, -1, 1)) // expected 4

	fmt.Println("=== Approach 2: Merge Sort Count (Divide and Conquer) ===")
	fmt.Printf("nums=[-2,5,-1], lower=-2, upper=2   got=%d  expected 3\n", mergeSortCount(ex1, -2, 2)) // expected 3
	fmt.Printf("nums=[0], lower=0, upper=0          got=%d  expected 1\n", mergeSortCount(ex2, 0, 0))  // expected 1
	fmt.Printf("nums=[2,-2,2,-2], lower=-1, upper=1 got=%d  expected 4\n", mergeSortCount(ex3, -1, 1)) // expected 4

	fmt.Println("=== Approach 3: Binary Indexed Tree (Fenwick) ===")
	fmt.Printf("nums=[-2,5,-1], lower=-2, upper=2   got=%d  expected 3\n", binaryIndexedTree(ex1, -2, 2)) // expected 3
	fmt.Printf("nums=[0], lower=0, upper=0          got=%d  expected 1\n", binaryIndexedTree(ex2, 0, 0))  // expected 1
	fmt.Printf("nums=[2,-2,2,-2], lower=-1, upper=1 got=%d  expected 4\n", binaryIndexedTree(ex3, -1, 1)) // expected 4
}
