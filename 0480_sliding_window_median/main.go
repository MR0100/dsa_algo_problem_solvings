package main

import (
	"container/heap"
	"fmt"
	"sort"
)

// ── Approach 1: Sorted Window (Binary Search Insert / Delete) ────────────────
//
// sortedWindow solves Sliding Window Median by keeping the current window as a
// SORTED slice; each slide binary-searches to remove the outgoing value and to
// insert the incoming one, so the median is always the middle of the slice.
//
// Intuition:
//
//	If the k elements of the window are kept sorted, the median is O(1): the
//	middle element (odd k) or the average of the two middle elements (even k).
//	Sliding the window changes only two elements — one leaves on the left, one
//	enters on the right. Use binary search to locate both positions in the
//	sorted slice; deletion and insertion each shift O(k) elements.
//
// Algorithm:
//  1. Build the first window and sort it.
//  2. Record its median.
//  3. For each subsequent right end:
//     - binary-search the outgoing element and delete it (shift left).
//     - binary-search the insert position of the incoming element and insert.
//     - record the median.
//  4. Return the collected medians.
//
// Time:  O(n·k) — n slides, each an O(k) shift for insert+delete (search is
//
//	O(log k) but the memmove dominates).
//
// Space: O(k) for the window slice.
func sortedWindow(nums []int, k int) []float64 {
	window := make([]int, k)
	copy(window, nums[:k]) // first k elements
	sort.Ints(window)      // keep the window sorted at all times

	medians := make([]float64, 0, len(nums)-k+1)
	medians = append(medians, medianOfSorted(window, k)) // median of window #0

	for right := k; right < len(nums); right++ {
		outgoing := nums[right-k] // element leaving on the left
		// Locate outgoing via binary search (it is guaranteed present).
		idx := sort.SearchInts(window, outgoing)
		// Delete it by shifting the tail left one slot.
		window = append(window[:idx], window[idx+1:]...)

		incoming := nums[right] // element entering on the right
		// Find where incoming belongs to keep the slice sorted.
		pos := sort.SearchInts(window, incoming)
		// Insert by growing the slice and shifting the tail right one slot.
		window = append(window, 0)         // extend length by one (value overwritten)
		copy(window[pos+1:], window[pos:]) // shift [pos:] right
		window[pos] = incoming             // drop incoming into its sorted place

		medians = append(medians, medianOfSorted(window, k))
	}
	return medians
}

// medianOfSorted returns the median of an already-sorted slice of length k.
func medianOfSorted(sorted []int, k int) float64 {
	if k%2 == 1 { // odd length → single middle element
		return float64(sorted[k/2])
	}
	// even length → average the two middle elements (in float to avoid overflow/round)
	return (float64(sorted[k/2-1]) + float64(sorted[k/2])) / 2.0
}

// ── Approach 2: Two Heaps with Lazy Deletion (Optimal) ───────────────────────
//
// twoHeaps solves Sliding Window Median with a max-heap (`lo`, the smaller
// half) and a min-heap (`hi`, the larger half), balanced so the median sits at
// the heap tops. Out-of-window elements are removed lazily via a delete-count
// map and purged from a heap top only when they surface.
//
// Intuition:
//
//	Split the window into a lower half (max-heap) and an upper half (min-heap)
//	with |lo| == |hi| or |lo| == |hi|+1. Then the median is lo.top (odd k) or
//	the average of lo.top and hi.top (even k). Insertion pushes to the correct
//	half and rebalances. True removal from the middle of a heap is O(k), so we
//	defer it: mark the outgoing value in a `delayed` map, adjust a running
//	balance, and only pop it once it reaches a heap's top ("lazy deletion").
//
// Algorithm:
//  1. For each element add(incoming); once the window is full, also
//     remove(outgoing) lazily; then prune tops and read the median.
//  2. add: push to lo if ≤ lo.top else hi; rebalance sizes.
//  3. remove: delayed[x]++; if x ≤ lo.top it belonged to lo (balance--), else
//     hi (balance++); prune that heap's top if it now equals x.
//  4. prune tops before every median read and after rebalancing.
//
// Time:  O(n·log k) — each element pushed/popped a constant number of times.
// Space: O(k) for the heaps and the delayed map.
func twoHeaps(nums []int, k int) []float64 {
	lo := &maxHeap{}         // smaller half; top = largest of the small side
	hi := &minHeap{}         // larger half;  top = smallest of the large side
	delayed := map[int]int{} // value → how many pending (lazy) deletions

	medians := make([]float64, 0, len(nums)-k+1)

	// prune removes elements scheduled for deletion from the top of a heap.
	pruneMax := func() {
		for lo.Len() > 0 {
			top := (*lo)[0]
			if delayed[top] > 0 {
				delayed[top]--
				heap.Pop(lo)
			} else {
				break
			}
		}
	}
	pruneMin := func() {
		for hi.Len() > 0 {
			top := (*hi)[0]
			if delayed[top] > 0 {
				delayed[top]--
				heap.Pop(hi)
			} else {
				break
			}
		}
	}

	// balance is (effective size of lo) − (effective size of hi), where
	// "effective" ignores elements still sitting in `delayed`.
	balance := 0

	add := func(x int) {
		if lo.Len() == 0 || x <= (*lo)[0] {
			heap.Push(lo, x) // belongs to the lower half
			balance++
		} else {
			heap.Push(hi, x) // belongs to the upper half
			balance--
		}
	}
	remove := func(x int) {
		delayed[x]++ // schedule a lazy deletion
		if lo.Len() > 0 && x <= (*lo)[0] {
			balance-- // it logically leaves the lower half
			if x == (*lo)[0] {
				pruneMax() // if it's exactly the top, purge it now
			}
		} else {
			balance++ // it logically leaves the upper half
			if hi.Len() > 0 && x == (*hi)[0] {
				pruneMin()
			}
		}
	}

	for i := 0; i < len(nums); i++ {
		add(nums[i]) // grow the window on the right
		if i >= k {
			remove(nums[i-k]) // shrink it on the left once it's oversized
		}

		// Rebalance so lo holds either the same count as hi, or one more.
		// balance == 0 means equal; balance == 1 means lo has one extra.
		if balance > 1 { // lo too big → move its top to hi
			heap.Push(hi, heap.Pop(lo))
			balance -= 2
			pruneMax()
		} else if balance < 0 { // hi too big → move its top to lo
			heap.Push(lo, heap.Pop(hi))
			balance += 2
			pruneMin()
		}
		pruneMax() // ensure both tops are real (not tombstoned) before reading
		pruneMin()

		if i >= k-1 { // window is full → emit a median
			if k%2 == 1 {
				medians = append(medians, float64((*lo)[0]))
			} else {
				medians = append(medians, (float64((*lo)[0])+float64((*hi)[0]))/2.0)
			}
		}
	}
	return medians
}

// maxHeap is a max-heap of ints (largest at index 0) for the lower half.
type maxHeap []int

func (h maxHeap) Len() int           { return len(h) }
func (h maxHeap) Less(i, j int) bool { return h[i] > h[j] } // '>' → max-heap
func (h maxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *maxHeap) Pop() any {
	old := *h
	n := len(old)
	v := old[n-1]
	*h = old[:n-1]
	return v
}

// minHeap is a min-heap of ints (smallest at index 0) for the upper half.
type minHeap []int

func (h minHeap) Len() int           { return len(h) }
func (h minHeap) Less(i, j int) bool { return h[i] < h[j] } // '<' → min-heap
func (h minHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *minHeap) Pop() any {
	old := *h
	n := len(old)
	v := old[n-1]
	*h = old[:n-1]
	return v
}

func main() {
	fmt.Println("=== Approach 1: Sorted Window (Binary Search) ===")
	fmt.Printf("nums=[1,3,-1,-3,5,3,6,7], k=3  got=%v\n", sortedWindow([]int{1, 3, -1, -3, 5, 3, 6, 7}, 3))
	fmt.Println("expected                        [1 -1 -1 3 5 6]")
	fmt.Printf("nums=[1,2,3,4,2,3,1,4,2], k=3   got=%v\n", sortedWindow([]int{1, 2, 3, 4, 2, 3, 1, 4, 2}, 3))
	fmt.Println("expected                        [2 3 3 3 2 3 2]")
	fmt.Printf("nums=[1,4,2,3], k=4             got=%v\n", sortedWindow([]int{1, 4, 2, 3}, 4))
	fmt.Println("expected                        [2.5]")
	fmt.Printf("nums=[2147483647,2147483647], k=2  got=%v\n", sortedWindow([]int{2147483647, 2147483647}, 2))
	fmt.Println("expected                            [2.147483647e+09]")

	fmt.Println("=== Approach 2: Two Heaps with Lazy Deletion (Optimal) ===")
	fmt.Printf("nums=[1,3,-1,-3,5,3,6,7], k=3  got=%v\n", twoHeaps([]int{1, 3, -1, -3, 5, 3, 6, 7}, 3))
	fmt.Println("expected                        [1 -1 -1 3 5 6]")
	fmt.Printf("nums=[1,2,3,4,2,3,1,4,2], k=3   got=%v\n", twoHeaps([]int{1, 2, 3, 4, 2, 3, 1, 4, 2}, 3))
	fmt.Println("expected                        [2 3 3 3 2 3 2]")
	fmt.Printf("nums=[1,4,2,3], k=4             got=%v\n", twoHeaps([]int{1, 4, 2, 3}, 4))
	fmt.Println("expected                        [2.5]")
	fmt.Printf("nums=[2147483647,2147483647], k=2  got=%v\n", twoHeaps([]int{2147483647, 2147483647}, 2))
	fmt.Println("expected                            [2.147483647e+09]")
}
