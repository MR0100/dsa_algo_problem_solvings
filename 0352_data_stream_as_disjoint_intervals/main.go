package main

import (
	"fmt"
	"sort"
)

// The task: support a stream of integers. After each insertion we must be able
// to report the current set of covered integers as a minimal list of disjoint,
// sorted intervals [start, end]. E.g. after adding 1,3,7,2,6 the covered set is
// {1,2,3,6,7} → intervals [[1,3],[6,7]].
//
// We present two designs of the SummaryRanges type:
//   1. bruteForceSummaryRanges — keep a boolean/set of seen numbers, rebuild the
//      interval list on demand by scanning.
//   2. SummaryRanges (optimal) — keep the disjoint intervals themselves sorted,
//      and on each addNum binary-search the insert position and merge neighbours.

// ── Approach 1: Brute Force (Seen Set, Rebuild on Query) ─────────────────────
//
// bruteForceSummaryRanges stores every distinct number in a set and, when asked,
// scans the sorted distinct values to coalesce consecutive runs into intervals.
//
// Intuition:
//
//	Correctness first: remember which integers we have seen, then whenever a
//	summary is requested, sort the seen values and merge consecutive ones into
//	runs. AddNum is trivial; the work is deferred to getIntervals.
//
// Algorithm:
//
//	AddNum(v):       insert v into a map-set (dedup automatically).
//	GetIntervals():  collect keys, sort them, then sweep: start a new interval
//	                 whenever the current value is not exactly prev+1.
//
// Time:  AddNum O(1) amortized; GetIntervals O(k log k) for k distinct values.
// Space: O(k) for the set.
type bruteForceSummaryRanges struct {
	seen map[int]bool // set of every distinct number added so far
}

// newBruteForceSummaryRanges builds an empty stream summary.
func newBruteForceSummaryRanges() *bruteForceSummaryRanges {
	return &bruteForceSummaryRanges{seen: make(map[int]bool)}
}

// AddNum records value in the seen set.
func (s *bruteForceSummaryRanges) AddNum(value int) {
	s.seen[value] = true // duplicates are harmless: map keeps one copy
}

// GetIntervals rebuilds the disjoint interval list from the seen set.
func (s *bruteForceSummaryRanges) GetIntervals() [][]int {
	if len(s.seen) == 0 {
		return [][]int{}
	}
	nums := make([]int, 0, len(s.seen)) // gather distinct values
	for v := range s.seen {
		nums = append(nums, v)
	}
	sort.Ints(nums) // ascending so consecutive runs are adjacent

	res := [][]int{}
	start := nums[0] // current run's start
	prev := nums[0]  // last value placed in the current run
	for i := 1; i < len(nums); i++ {
		if nums[i] == prev+1 {
			prev = nums[i] // extend the current run
			continue
		}
		res = append(res, []int{start, prev}) // gap: close the run
		start, prev = nums[i], nums[i]        // begin a new run
	}
	res = append(res, []int{start, prev}) // flush the final run
	return res
}

// ── Approach 2: Sorted Disjoint Intervals + Binary Search (Optimal) ──────────
//
// SummaryRanges keeps the merged intervals themselves, sorted by start. AddNum
// binary-searches where the new value lands and merges with the left/right
// neighbours as needed, so getIntervals is a trivial O(k) copy.
//
// Intuition:
//
//	Maintain the invariant "intervals is always a sorted list of disjoint,
//	non-adjacent ranges". Adding a value can, at most, (a) sit inside an
//	existing interval (no-op), (b) extend one interval by one, or (c) bridge two
//	intervals into one, or (d) create a brand-new singleton interval. A binary
//	search locates the neighbourhood; a constant amount of stitching fixes it.
//
// Algorithm:
//
//	AddNum(v):
//	  1. Binary-search idx = first interval with start >= v.
//	  2. If v already covered by intervals[idx-1] (its end >= v), return.
//	  3. Determine mergeLeft  = intervals[idx-1].end == v-1.
//	     Determine mergeRight = idx < len && intervals[idx].start == v+1.
//	  4. Four cases: both → coalesce idx-1 and idx; left only → extend end;
//	     right only → lower start; neither → splice in [v,v].
//	GetIntervals(): return the stored slice.
//
// Time:  AddNum O(log k) search + O(k) worst-case slice splice; GetIntervals O(k).
// Space: O(k) intervals.
type SummaryRanges struct {
	intervals [][]int // sorted, disjoint, non-adjacent [start,end] ranges
}

// Constructor builds an empty SummaryRanges.
func Constructor() SummaryRanges {
	return SummaryRanges{intervals: [][]int{}}
}

// AddNum inserts value, maintaining the sorted disjoint-interval invariant.
func (s *SummaryRanges) AddNum(value int) {
	iv := s.intervals
	// idx = first interval whose start is >= value.
	idx := sort.Search(len(iv), func(i int) bool { return iv[i][0] >= value })

	// Already covered by the interval just to the left?
	if idx > 0 && iv[idx-1][1] >= value {
		return // value lies inside an existing interval — nothing changes
	}
	// Exactly equal to iv[idx][0] means the right neighbour already starts here.
	if idx < len(iv) && iv[idx][0] == value {
		return // value is the start of an existing interval — already covered
	}

	mergeLeft := idx > 0 && iv[idx-1][1] == value-1      // touches left range
	mergeRight := idx < len(iv) && iv[idx][0] == value+1 // touches right range

	switch {
	case mergeLeft && mergeRight:
		// Bridge left and right into a single interval, drop the right one.
		iv[idx-1][1] = iv[idx][1]
		s.intervals = append(iv[:idx], iv[idx+1:]...)
	case mergeLeft:
		iv[idx-1][1] = value // extend the left interval's end by one
	case mergeRight:
		iv[idx][0] = value // lower the right interval's start by one
	default:
		// Isolated value: splice a new singleton interval at position idx.
		s.intervals = append(iv, nil)                // grow by one (reuse capacity)
		copy(s.intervals[idx+1:], s.intervals[idx:]) // shift tail right
		s.intervals[idx] = []int{value, value}       // place [v,v]
	}
}

// GetIntervals returns the current disjoint interval list (already maintained).
func (s *SummaryRanges) GetIntervals() [][]int {
	return s.intervals
}

func main() {
	// Official example operation sequence:
	//   addNum 1 -> [[1,1]]
	//   addNum 3 -> [[1,1],[3,3]]
	//   addNum 7 -> [[1,1],[3,3],[7,7]]
	//   addNum 2 -> [[1,3],[7,7]]
	//   addNum 6 -> [[1,3],[6,7]]

	fmt.Println("=== Approach 1: Brute Force (Seen Set) ===")
	b := newBruteForceSummaryRanges()
	b.AddNum(1)
	fmt.Println(b.GetIntervals()) // expected [[1 1]]
	b.AddNum(3)
	fmt.Println(b.GetIntervals()) // expected [[1 1] [3 3]]
	b.AddNum(7)
	fmt.Println(b.GetIntervals()) // expected [[1 1] [3 3] [7 7]]
	b.AddNum(2)
	fmt.Println(b.GetIntervals()) // expected [[1 3] [7 7]]
	b.AddNum(6)
	fmt.Println(b.GetIntervals()) // expected [[1 3] [6 7]]

	fmt.Println("=== Approach 2: Sorted Intervals + Binary Search (Optimal) ===")
	s := Constructor()
	s.AddNum(1)
	fmt.Println(s.GetIntervals()) // expected [[1 1]]
	s.AddNum(3)
	fmt.Println(s.GetIntervals()) // expected [[1 1] [3 3]]
	s.AddNum(7)
	fmt.Println(s.GetIntervals()) // expected [[1 1] [3 3] [7 7]]
	s.AddNum(2)
	fmt.Println(s.GetIntervals()) // expected [[1 3] [7 7]]
	s.AddNum(6)
	fmt.Println(s.GetIntervals()) // expected [[1 3] [6 7]]
}
