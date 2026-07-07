package main

import (
	"container/heap"
	"fmt"
	"sort"
)

// ── Approach 1: Min-Heap of End Times ────────────────────────────────────────
//
// minHeap solves Meeting Rooms II by keeping a min-heap of the end times of the
// rooms currently in use; the heap size at any moment is the number of rooms.
//
// Intuition:
//
//	Process meetings in start-time order. The room that frees up soonest is the
//	one whose end time is smallest — the heap's top. When the next meeting starts:
//	if it starts at or after that earliest end time, that room is free, so reuse
//	it (pop, then push the new end). Otherwise all rooms are busy, so allocate a
//	new one (just push). The peak heap size is the answer.
//
// Algorithm:
//  1. Sort intervals by start time.
//  2. Min-heap h of end times, initially empty.
//  3. For each meeting [s,e]: if h non-empty and h.top() <= s, pop (reuse room).
//     Then push e. Track max heap size seen (== current when we always push).
//  4. Answer is len(h) after processing all (never shrinks below peak with reuse),
//     which equals the maximum concurrency.
//
// Time:  O(n log n) — sort plus n heap operations.
// Space: O(n) — the heap.
type intHeap []int

func (h intHeap) Len() int            { return len(h) }
func (h intHeap) Less(i, j int) bool  { return h[i] < h[j] } // min-heap on end times
func (h intHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *intHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *intHeap) Pop() interface{} {
	old := *h
	n := len(old)
	v := old[n-1]
	*h = old[:n-1]
	return v
}

func minHeap(intervals [][]int) int {
	if len(intervals) == 0 {
		return 0
	}
	// Sort by start so we allocate rooms in chronological order.
	sort.Slice(intervals, func(a, b int) bool {
		return intervals[a][0] < intervals[b][0]
	})
	h := &intHeap{}
	heap.Init(h)
	for _, iv := range intervals {
		start, end := iv[0], iv[1]
		// If the soonest-freeing room is free by this meeting's start, reuse it.
		if h.Len() > 0 && (*h)[0] <= start {
			heap.Pop(h) // that room is done; remove its old end time
		}
		heap.Push(h, end) // occupy a room (reused or brand new) with this end time
	}
	// Heap size = number of rooms simultaneously in use at the peak.
	return h.Len()
}

// ── Approach 2: Chronological Sweep of Split Endpoints (Optimal) ──────────────
//
// sweepLine solves Meeting Rooms II by separately sorting all start times and
// all end times, then sweeping a two-pointer merge counting concurrent meetings.
//
// Intuition:
//
//	The number of rooms needed at any instant equals the number of meetings in
//	progress. Sort starts and ends independently. Walk a pointer through starts;
//	for each start, first release every meeting that has already ended (end <=
//	current start), decrementing the live count. Then this meeting starts, so
//	increment. Track the maximum live count — that's the minimum rooms.
//
// Algorithm:
//  1. Build sorted arrays starts and ends.
//  2. Pointers s=0, e=0; rooms=0, maxRooms=0.
//  3. While s < n: if starts[s] < ends[e], a new meeting begins before the next
//     end → rooms++, s++, update maxRooms. Else a meeting ends first → rooms--, e++.
//  4. Return maxRooms.
//
// Time:  O(n log n) — two sorts.
// Space: O(n) — the two endpoint arrays.
func sweepLine(intervals [][]int) int {
	n := len(intervals)
	if n == 0 {
		return 0
	}
	starts := make([]int, n)
	ends := make([]int, n)
	for i, iv := range intervals {
		starts[i] = iv[0] // collect all start times
		ends[i] = iv[1]   // collect all end times
	}
	sort.Ints(starts) // sort the two timelines independently
	sort.Ints(ends)

	rooms, maxRooms := 0, 0
	s, e := 0, 0
	for s < n {
		if starts[s] < ends[e] {
			// A meeting begins before the next one ends ⇒ need another room.
			rooms++
			s++
			if rooms > maxRooms {
				maxRooms = rooms // record new peak concurrency
			}
		} else {
			// The earliest end is at/or before this start ⇒ a room frees up.
			rooms--
			e++
		}
	}
	return maxRooms
}

func main() {
	fmt.Println("=== Approach 1: Min-Heap of End Times ===")
	fmt.Println(minHeap([][]int{{0, 30}, {5, 10}, {15, 20}})) // expected 2
	fmt.Println(minHeap([][]int{{7, 10}, {2, 4}}))            // expected 1

	fmt.Println("=== Approach 2: Sweep Line (Optimal) ===")
	fmt.Println(sweepLine([][]int{{0, 30}, {5, 10}, {15, 20}})) // expected 2
	fmt.Println(sweepLine([][]int{{7, 10}, {2, 4}}))            // expected 1
}
