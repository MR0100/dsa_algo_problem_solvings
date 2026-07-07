package main

import "fmt"

// numArrayADT is the common interface every implementation satisfies so main()
// can drive the same official operation sequence through each approach.
type numArrayADT interface {
	Update(index, val int)        // set nums[index] = val
	SumRange(left, right int) int // return nums[left] + ... + nums[right]
}

// ── Approach 1: Brute Force (mutable array, linear range sum) ─────────────────
//
// BruteForce solves Range Sum Query - Mutable by keeping the raw array and
// summing the requested window on every query.
//
// Intuition:
//
//	The simplest structure satisfying the API: store nums directly. Update is a
//	single assignment; SumRange walks the window and adds. This is the baseline
//	the smarter structures must beat when updates and queries are interleaved.
//
// Algorithm:
//
//	Update:   nums[index] = val — O(1).
//	SumRange: loop i from left to right accumulating nums[i] — O(n).
//
// Time:  Update O(1), SumRange O(n).
// Space: O(n) — the stored array.
type BruteForce struct {
	nums []int // live copy of the array
}

// NewBruteForce copies the input array.
func NewBruteForce(nums []int) *BruteForce {
	cp := make([]int, len(nums))
	copy(cp, nums)
	return &BruteForce{nums: cp}
}

// Update assigns a new value in O(1).
func (b *BruteForce) Update(index, val int) {
	b.nums[index] = val // direct assignment
}

// SumRange adds the inclusive window in O(n).
func (b *BruteForce) SumRange(left, right int) int {
	sum := 0
	for i := left; i <= right; i++ {
		sum += b.nums[i] // accumulate each element
	}
	return sum
}

// ── Approach 2: Fenwick Tree / Binary Indexed Tree (Optimal) ─────────────────
//
// FenwickTree solves Range Sum Query - Mutable with a Binary Indexed Tree that
// supports point update and prefix sum in O(log n), so both operations are fast.
//
// Intuition:
//
//	A Fenwick tree stores partial sums indexed so that any prefix sum is the sum
//	of O(log n) tree cells, and a point update touches O(log n) cells. Cell i is
//	responsible for the range (i - lowbit(i), i], where lowbit(i) = i & -i is the
//	lowest set bit. Range sum [l, r] = prefix(r) - prefix(l-1). Because updates
//	come as "set to val" (not "add delta"), we track the current values and add
//	the difference val - old.
//
// Algorithm:
//
//	Build:    tree of size n+1 (1-indexed); add each initial value.
//	add(i,d): while i <= n: tree[i] += d; i += lowbit(i).
//	prefix(i):while i > 0: s += tree[i]; i -= lowbit(i).
//	Update:   d = val - nums[index]; nums[index] = val; add(index+1, d).
//	SumRange: prefix(right+1) - prefix(left).
//
// Time:  Update O(log n), SumRange O(log n), Build O(n log n).
// Space: O(n) — the tree plus the value cache.
type FenwickTree struct {
	n    int   // number of elements
	tree []int // 1-indexed BIT storage, length n+1
	nums []int // current values, to compute the delta on Update
}

// NewFenwickTree builds the BIT from the initial array.
func NewFenwickTree(nums []int) *FenwickTree {
	n := len(nums)
	f := &FenwickTree{
		n:    n,
		tree: make([]int, n+1), // index 0 unused
		nums: make([]int, n),
	}
	for i, v := range nums {
		f.nums[i] = v // remember the value for future deltas
		f.add(i+1, v) // seed the tree with this value (1-indexed)
	}
	return f
}

// add applies delta at 1-indexed position i, climbing to all responsible cells.
func (f *FenwickTree) add(i, delta int) {
	for ; i <= f.n; i += i & (-i) { // jump by the lowest set bit
		f.tree[i] += delta
	}
}

// prefix returns nums[0] + ... + nums[i-1] (sum of the first i elements).
func (f *FenwickTree) prefix(i int) int {
	s := 0
	for ; i > 0; i -= i & (-i) { // strip the lowest set bit each step
		s += f.tree[i]
	}
	return s
}

// Update sets nums[index] = val by pushing the difference into the tree.
func (f *FenwickTree) Update(index, val int) {
	delta := val - f.nums[index] // how much the value changes
	f.nums[index] = val          // record the new value
	f.add(index+1, delta)        // propagate the delta (1-indexed)
}

// SumRange returns the inclusive window sum via two prefix queries.
func (f *FenwickTree) SumRange(left, right int) int {
	return f.prefix(right+1) - f.prefix(left)
}

// ── Approach 3: Segment Tree (array-based, point update) ─────────────────────
//
// SegmentTree solves Range Sum Query - Mutable with a classic array-backed
// segment tree, the go-to structure when queries later generalize (min, max,
// lazy range updates, etc.).
//
// Intuition:
//
//	A segment tree stores range sums in a binary tree flattened into an array of
//	size 2n: leaves n..2n-1 hold the elements, and each internal node holds the
//	sum of its two children. A point update walks leaf → root fixing O(log n)
//	parents; a range query splits [l, r] into O(log n) canonical segments.
//
// Algorithm:
//
//	Build:    place values at [n, 2n); for i = n-1 down to 1, tree[i] =
//	          tree[2i] + tree[2i+1].
//	Update:   set leaf tree[index+n] = val; walk up halving the index, resumming.
//	SumRange: l = left+n, r = right+n; while l <= r, add boundary nodes when
//	          they are right/left children and move inward; halve both.
//
// Time:  Update O(log n), SumRange O(log n), Build O(n).
// Space: O(n) — 2n array.
type SegmentTree struct {
	n    int   // number of leaves
	tree []int // 2n-sized array; [n,2n) are leaves
}

// NewSegmentTree builds the tree bottom-up from the input array.
func NewSegmentTree(nums []int) *SegmentTree {
	n := len(nums)
	t := make([]int, 2*n)
	for i := 0; i < n; i++ {
		t[n+i] = nums[i] // fill the leaf layer
	}
	for i := n - 1; i > 0; i-- {
		t[i] = t[2*i] + t[2*i+1] // internal node = sum of children
	}
	return &SegmentTree{n: n, tree: t}
}

// Update sets nums[index] = val and repairs the path to the root.
func (s *SegmentTree) Update(index, val int) {
	i := index + s.n // leaf position in the flattened array
	s.tree[i] = val  // overwrite the leaf
	for i > 1 {      // climb to the root
		i /= 2                                  // move to parent
		s.tree[i] = s.tree[2*i] + s.tree[2*i+1] // recompute parent sum
	}
}

// SumRange returns the inclusive window sum by merging canonical segments.
func (s *SegmentTree) SumRange(left, right int) int {
	l, r := left+s.n, right+s.n // map to leaf indices
	sum := 0
	for l <= r {
		if l%2 == 1 { // l is a right child → it is fully inside; take it
			sum += s.tree[l]
			l++
		}
		if r%2 == 0 { // r is a left child → it is fully inside; take it
			sum += s.tree[r]
			r--
		}
		l /= 2 // ascend one level
		r /= 2
	}
	return sum
}

// runExample drives one official operation sequence through an implementation
// and returns the list of outputs in LeetCode's format (null for void ops).
//
// Ops:  ["NumArray","sumRange","update","sumRange"]
// Args: [[[1,3,5]],[0,2],[1,2],[0,2]]
func runExample(build func(nums []int) numArrayADT) []interface{} {
	na := build([]int{1, 3, 5})          // constructor → null
	out := []interface{}{nil}            //
	out = append(out, na.SumRange(0, 2)) // → 9
	na.Update(1, 2)                      // → null
	out = append(out, nil)               //
	out = append(out, na.SumRange(0, 2)) // → 8
	return out
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(runExample(func(n []int) numArrayADT { return NewBruteForce(n) })) // [<nil> 9 <nil> 8]

	fmt.Println("=== Approach 2: Fenwick Tree (Optimal) ===")
	fmt.Println(runExample(func(n []int) numArrayADT { return NewFenwickTree(n) })) // [<nil> 9 <nil> 8]

	fmt.Println("=== Approach 3: Segment Tree ===")
	fmt.Println(runExample(func(n []int) numArrayADT { return NewSegmentTree(n) })) // [<nil> 9 <nil> 8]
}
