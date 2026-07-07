package main

import (
	"container/heap"
	"fmt"
	"sort"
)

// ── Approach 1: Sweep Line + Max-Heap ────────────────────────────────────────
//
// sweepLineHeap solves The Skyline Problem by sweeping a vertical line across
// all distinct x-coordinates and, at each event, tracking the tallest building
// currently "alive" with a max-heap.
//
// Intuition:
//
//	The skyline only changes its height at building EDGES (a left edge may
//	raise the current tallest, a right edge may lower it). So the only x's we
//	care about are the buildings' left and right sides. Sweep left→right; at
//	every event x, ensure the heap holds exactly the heights of buildings that
//	are currently spanning x. The heap's max is the skyline height there. A key
//	point emerges whenever that max changes from the previous point.
//
// Algorithm:
//  1. Build events: for each building [L,R,H], a start event (L, -H) and an end
//     event (R, +H). Encode start heights as negative so that, when sorting,
//     at equal x, starts (taller first) come before ends, and among starts the
//     taller building is processed first — this avoids spurious key points.
//  2. Sort events by (x, height-code).
//  3. Maintain a max-heap of currently-active heights, plus a lazy-deletion
//     multiset count so we can "remove" an ended height.
//  4. For each event: add its height (start) or schedule removal (end). Then
//     read the current max alive height. If it differs from the previous max,
//     append [x, curMax] to the result.
//
// Time:  O(n log n) — sorting events and O(log n) heap ops per event.
// Space: O(n) — events, heap, and the removal multiset.
func sweepLineHeap(buildings [][]int) [][]int {
	// Each event is [x, heightCode]: negative code = building start (height H),
	// positive code = building end (height H).
	events := make([][2]int, 0, len(buildings)*2)
	for _, b := range buildings {
		L, R, H := b[0], b[1], b[2]
		events = append(events, [2]int{L, -H}) // start: negative height
		events = append(events, [2]int{R, H})  // end: positive height
	}
	// Sort by x; ties broken by heightCode so that at the same x:
	//   - starts (negative) precede ends (positive),
	//   - among starts, the taller (more negative) is processed first,
	//   - among ends, the shorter is processed first.
	sort.Slice(events, func(i, j int) bool {
		if events[i][0] != events[j][0] {
			return events[i][0] < events[j][0]
		}
		return events[i][1] < events[j][1]
	})

	result := make([][]int, 0)
	pq := &maxHeap{0}            // active heights; ground 0 always present
	removed := make(map[int]int) // lazy-deletion counts: height → pending removals
	prevMax := 0                 // skyline height before this event

	for _, e := range events {
		x, code := e[0], e[1]
		if code < 0 {
			h := -code       // a building starts here at height h
			heap.Push(pq, h) // it is now active
		} else {
			removed[code]++ // schedule removal of this ended height
		}
		// Pop off any heights whose removal is pending and that sit on top.
		for pq.Len() > 0 {
			top := (*pq)[0]
			if removed[top] > 0 { // this top height has actually ended
				removed[top]--
				heap.Pop(pq)
			} else {
				break // the real current max is on top
			}
		}
		curMax := (*pq)[0] // tallest active building (0 if none)
		if curMax != prevMax {
			// Height changed at x → this is a key point of the skyline.
			result = append(result, []int{x, curMax})
			prevMax = curMax
		}
	}
	return result
}

// maxHeap is a max-heap of ints (largest on top) for the active building
// heights. Implements heap.Interface.
type maxHeap []int

func (h maxHeap) Len() int            { return len(h) }
func (h maxHeap) Less(i, j int) bool  { return h[i] > h[j] } // ">" makes it a MAX-heap
func (h maxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *maxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	v := old[n-1]
	*h = old[:n-1]
	return v
}

// ── Approach 2: Divide and Conquer ───────────────────────────────────────────
//
// divideAndConquer solves The Skyline Problem by splitting the buildings into
// two halves, computing each half's skyline recursively, and merging them like
// the merge step of merge sort.
//
// Intuition:
//
//	A skyline of one building is trivial: [[L,H],[R,0]]. Two skylines can be
//	merged by sweeping both key-point lists together, at each x taking the max
//	of the two current heights, and emitting a point whenever that max changes.
//	Recursing on halves and merging mirrors merge sort exactly.
//
// Algorithm:
//  1. Base case: 0 buildings → empty; 1 building → [[L,H],[R,0]].
//  2. Split buildings in half; recurse on each to get left/right skylines.
//  3. Merge: walk both lists by x, maintaining leftH and rightH (current height
//     contributed by each side). At each processed x, the merged height is
//     max(leftH, rightH); emit [x, merged] whenever it differs from the last.
//
// Time:  O(n log n) — T(n) = 2T(n/2) + O(n) merge.
// Space: O(n) — recursion and output lists.
func divideAndConquer(buildings [][]int) [][]int {
	if len(buildings) == 0 {
		return [][]int{}
	}
	if len(buildings) == 1 {
		b := buildings[0]
		// One building becomes two key points: rise at L, fall to 0 at R.
		return [][]int{{b[0], b[2]}, {b[1], 0}}
	}
	mid := len(buildings) / 2
	left := divideAndConquer(buildings[:mid])  // skyline of the left half
	right := divideAndConquer(buildings[mid:]) // skyline of the right half
	return mergeSkylines(left, right)
}

// mergeSkylines merges two skylines (each a list of [x, height] key points) by
// a two-pointer sweep, emitting a point wherever the combined height changes.
func mergeSkylines(left, right [][]int) [][]int {
	merged := make([][]int, 0, len(left)+len(right))
	i, j := 0, 0
	leftH, rightH := 0, 0 // current height contributed by each skyline
	for i < len(left) && j < len(right) {
		var x, h int
		if left[i][0] < right[j][0] { // left key point comes first
			x = left[i][0]
			leftH = left[i][1] // update left's current height
			i++
		} else if left[i][0] > right[j][0] { // right key point comes first
			x = right[j][0]
			rightH = right[j][1]
			j++
		} else { // same x: consume both
			x = left[i][0]
			leftH = left[i][1]
			rightH = right[j][1]
			i++
			j++
		}
		h = max(leftH, rightH) // combined skyline height at x
		if len(merged) == 0 || merged[len(merged)-1][1] != h {
			merged = append(merged, []int{x, h}) // only emit real changes
		}
	}
	// Append the leftovers (only one list can have remaining points).
	for i < len(left) {
		x, h := left[i][0], left[i][1]
		if len(merged) == 0 || merged[len(merged)-1][1] != h {
			merged = append(merged, []int{x, h})
		}
		i++
	}
	for j < len(right) {
		x, h := right[j][0], right[j][1]
		if len(merged) == 0 || merged[len(merged)-1][1] != h {
			merged = append(merged, []int{x, h})
		}
		j++
	}
	return merged
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	fmt.Println("=== Approach 1: Sweep Line + Max-Heap ===")
	fmt.Println(sweepLineHeap([][]int{{2, 9, 10}, {3, 7, 15}, {5, 12, 12}, {15, 20, 10}, {19, 24, 8}}))
	// expected [[2 10] [3 15] [7 12] [12 0] [15 10] [20 8] [24 0]]
	fmt.Println(sweepLineHeap([][]int{{0, 2, 3}, {2, 5, 3}}))
	// expected [[0 3] [5 0]]

	fmt.Println("=== Approach 2: Divide and Conquer ===")
	fmt.Println(divideAndConquer([][]int{{2, 9, 10}, {3, 7, 15}, {5, 12, 12}, {15, 20, 10}, {19, 24, 8}}))
	// expected [[2 10] [3 15] [7 12] [12 0] [15 10] [20 8] [24 0]]
	fmt.Println(divideAndConquer([][]int{{0, 2, 3}, {2, 5, 3}}))
	// expected [[0 3] [5 0]]
}
