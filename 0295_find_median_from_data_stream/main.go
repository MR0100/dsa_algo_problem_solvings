package main

import (
	"container/heap"
	"fmt"
	"sort"
)

// medianFinder is the common interface every approach implements so main() can
// drive the same official operation sequence through all of them.
type medianFinder interface {
	AddNum(num int)
	FindMedian() float64
}

// ── Approach 1: Sorted Slice with Insertion (Brute Force) ────────────────────
//
// SortedSliceFinder keeps every number in a slice sorted at all times, so the
// median is just the middle element(s).
//
// Intuition:
//
//	If the data is always sorted, the median is trivial: the middle element
//	(odd count) or the mean of the two middle elements (even count). The cost is
//	pushed to insertion: each AddNum must place the new value at its sorted
//	position, shifting the tail.
//
// Algorithm:
//
//	AddNum:     binary-search the insert position, splice the value in.
//	FindMedian: read the middle (or average the two middles).
//
// Time:  AddNum O(n) (shift), FindMedian O(1).
// Space: O(n) — all numbers retained.
type SortedSliceFinder struct {
	nums []int // kept in non-decreasing order after every AddNum
}

// NewSortedSliceFinder builds an empty sorted-slice median finder.
func NewSortedSliceFinder() *SortedSliceFinder { return &SortedSliceFinder{} }

// AddNum inserts num while preserving sorted order.
func (f *SortedSliceFinder) AddNum(num int) {
	// sort.SearchInts finds the leftmost index where num can go to stay sorted.
	i := sort.SearchInts(f.nums, num)
	f.nums = append(f.nums, 0)     // grow by one (value overwritten below)
	copy(f.nums[i+1:], f.nums[i:]) // shift the tail right to open a gap at i
	f.nums[i] = num                // drop num into its sorted slot
}

// FindMedian returns the middle value, averaging the two middles when even.
func (f *SortedSliceFinder) FindMedian() float64 {
	n := len(f.nums)
	if n%2 == 1 {
		return float64(f.nums[n/2]) // single middle element
	}
	// Even count: average the two central elements.
	return float64(f.nums[n/2-1]+f.nums[n/2]) / 2.0
}

// ── Approach 2: Two Heaps (Optimal) ──────────────────────────────────────────
//
// TwoHeapFinder maintains a max-heap of the smaller half and a min-heap of the
// larger half so the median sits at the heap tops.
//
// Intuition:
//
//	Split the sorted stream at the middle: the "low" half (a max-heap so its top
//	is the largest small value) and the "high" half (a min-heap so its top is the
//	smallest large value). Keep sizes balanced (low has the same count as high,
//	or one more). Then the median is low's top (odd total) or the average of the
//	two tops (even total). Both are O(1) to read; inserts are O(log n) heap ops.
//
// Algorithm:
//
//	AddNum:     push onto low; move low's top to high; if high got bigger than
//	            low, move high's top back — this both routes the value to the
//	            correct side and rebalances sizes.
//	FindMedian: if low bigger → low top; else average of the two tops.
//
// Time:  AddNum O(log n), FindMedian O(1).
// Space: O(n).
type TwoHeapFinder struct {
	low  *maxHeap // smaller half; top = largest of the small values
	high *minHeap // larger half;  top = smallest of the large values
}

// NewTwoHeapFinder builds an empty two-heap median finder.
func NewTwoHeapFinder() *TwoHeapFinder {
	return &TwoHeapFinder{low: &maxHeap{}, high: &minHeap{}}
}

// AddNum routes num into the correct half and rebalances the heaps.
func (f *TwoHeapFinder) AddNum(num int) {
	heap.Push(f.low, num)              // tentatively add to the low half
	heap.Push(f.high, heap.Pop(f.low)) // shift low's max into high (keeps order)
	if f.high.Len() > f.low.Len() {    // high grew too large...
		heap.Push(f.low, heap.Pop(f.high)) // ...move its min back to low
	}
}

// FindMedian reads the median off the heap tops.
func (f *TwoHeapFinder) FindMedian() float64 {
	if f.low.Len() > f.high.Len() {
		return float64((*f.low)[0]) // odd total: low holds the extra element
	}
	// Even total: median is the average of the two heap tops.
	return float64((*f.low)[0]+(*f.high)[0]) / 2.0
}

// maxHeap is an int max-heap (largest at index 0) built on container/heap.
type maxHeap []int

func (h maxHeap) Len() int            { return len(h) }
func (h maxHeap) Less(i, j int) bool  { return h[i] > h[j] } // '>' → max on top
func (h maxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *maxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	v := old[n-1]
	*h = old[:n-1]
	return v
}

// minHeap is an int min-heap (smallest at index 0).
type minHeap []int

func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(i, j int) bool  { return h[i] < h[j] } // '<' → min on top
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	v := old[n-1]
	*h = old[:n-1]
	return v
}

// runExample drives the official operation sequence through one implementation
// and returns the outputs in LeetCode's list format.
//
// Ops:  ["MedianFinder","addNum","addNum","findMedian","addNum","findMedian"]
// Args: [[],[1],[2],[],[3],[]]
// Out:  [null,null,null,1.5,null,2.0]
func runExample(newFinder func() medianFinder) string {
	mf := newFinder()
	out := "[null" // constructor → null
	mf.AddNum(1)   // addNum(1) → null
	out += ",null"
	mf.AddNum(2) // addNum(2) → null
	out += ",null"
	out += fmt.Sprintf(",%.1f", mf.FindMedian()) // findMedian → 1.5
	mf.AddNum(3)                                 // addNum(3) → null
	out += ",null"
	out += fmt.Sprintf(",%.1f", mf.FindMedian()) // findMedian → 2.0
	out += "]"
	return out
}

func main() {
	fmt.Println("=== Approach 1: Sorted Slice (Brute Force) ===")
	fmt.Println(runExample(func() medianFinder { return NewSortedSliceFinder() })) // [null,null,null,1.5,null,2.0]

	fmt.Println("=== Approach 2: Two Heaps (Optimal) ===")
	fmt.Println(runExample(func() medianFinder { return NewTwoHeapFinder() })) // [null,null,null,1.5,null,2.0]
}
