package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Brute Force ──────────────────────────────────────────────────
//
// bruteForce solves Contains Duplicate III by checking every pair of indices
// within indexDiff of each other against the valueDiff bound.
//
// Intuition:
//
//	We need indices i < j with (j − i) <= indexDiff and |nums[i] − nums[j]| <=
//	valueDiff. Directly test each i against the up-to-indexDiff elements after
//	it; the first pair that satisfies both bounds is our answer.
//
// Algorithm:
//  1. For each i, scan j from i+1 to min(i+indexDiff, n-1).
//  2. If abs(nums[i] − nums[j]) <= valueDiff, return true.
//  3. If nothing qualifies, return false.
//
// Time:  O(n·indexDiff) — each i inspects at most indexDiff neighbours.
// Space: O(1).
func bruteForce(nums []int, indexDiff int, valueDiff int) bool {
	n := len(nums)
	for i := 0; i < n; i++ {
		for j := i + 1; j <= i+indexDiff && j < n; j++ { // only indices within indexDiff
			if abs(nums[i]-nums[j]) <= valueDiff { // values close enough too
				return true
			}
		}
	}
	return false
}

// ── Approach 2: Sliding Window + Sorted Structure (TreeSet-style) ─────────────
//
// slidingWindowSorted solves Contains Duplicate III by keeping the last
// indexDiff values in sorted order and, for each new value, binary-searching
// for a neighbour within valueDiff.
//
// Intuition:
//
//	Restrict attention to a window of the last indexDiff values (this handles
//	the index bound). Within that window we need some value x with
//	|nums[i] − x| <= valueDiff, i.e. an x in [nums[i]−valueDiff,
//	nums[i]+valueDiff]. If the window is kept sorted, the smallest candidate
//	>= nums[i]−valueDiff is found by binary search; if that candidate is also
//	<= nums[i]+valueDiff, we have a valid pair. Slide the window by removing
//	the value that falls out of the index range.
//
// Algorithm:
//  1. Maintain a sorted slice `window` of the last indexDiff values.
//  2. For each i: binary-search the first element >= nums[i]−valueDiff.
//  3. If it exists and is <= nums[i]+valueDiff, return true.
//  4. Insert nums[i] into the sorted slice (keeping it ordered).
//  5. If the window now covers more than indexDiff indices, remove nums[i−indexDiff].
//  6. If loop ends, return false.
//
// Time:  O(n·log(indexDiff)) for search, but slice insert/delete is O(indexDiff),
//
//	so O(n·indexDiff) worst case with a plain slice. A balanced BST/skiplist
//	would make it O(n·log(indexDiff)) overall.
//
// Space: O(min(n, indexDiff)) — the window.
func slidingWindowSorted(nums []int, indexDiff int, valueDiff int) bool {
	window := make([]int, 0) // sorted values of the last indexDiff indices
	for i, v := range nums {
		// Find first element >= v - valueDiff (the lowest acceptable neighbour).
		pos := sort.SearchInts(window, v-valueDiff)
		// If such an element exists and is also <= v + valueDiff, it is within range.
		if pos < len(window) && window[pos] <= v+valueDiff {
			return true
		}
		// Insert v into the sorted window at its correct position.
		ins := sort.SearchInts(window, v)  // where v belongs
		window = append(window, 0)         // grow by one
		copy(window[ins+1:], window[ins:]) // shift right to open a gap
		window[ins] = v                    // place v
		// Evict the value that is now out of the index window.
		if i >= indexDiff {
			out := nums[i-indexDiff]            // value leaving the window
			del := sort.SearchInts(window, out) // find it (guaranteed present)
			window = append(window[:del], window[del+1:]...)
		}
	}
	return false
}

// ── Approach 3: Bucketing (Optimal) ──────────────────────────────────────────
//
// bucketing solves Contains Duplicate III in linear time by mapping values into
// buckets of width valueDiff+1 and checking each new value against its own and
// adjacent buckets.
//
// Intuition:
//
//	Give every value a bucket id = floor(v / (valueDiff+1)). Two values in the
//	SAME bucket differ by at most valueDiff automatically — instant hit. Values
//	within valueDiff can otherwise only sit in ADJACENT buckets (id−1 or id+1),
//	and there we still must verify the |difference| <= valueDiff explicitly.
//	Keep at most one value per bucket within the current index window: since a
//	same-bucket collision returns immediately, one slot per bucket suffices.
//	Slide by deleting the bucket of the value leaving the index window.
//
// Algorithm:
//  1. width = valueDiff + 1; bucket(v) = floor(v / width) (floor for negatives).
//  2. For each i with value v and b = bucket(v):
//  3. If bucket b already occupied → return true (same-bucket, diff <= valueDiff).
//  4. If bucket b−1 occupied and |v − its value| <= valueDiff → return true.
//  5. If bucket b+1 occupied and |v − its value| <= valueDiff → return true.
//  6. Put v in bucket b.
//  7. If i >= indexDiff, delete the bucket of nums[i−indexDiff].
//  8. If loop ends, return false.
//
// Time:  O(n) — each index does O(1) bucket work.
// Space: O(min(n, indexDiff)) — at most indexDiff+1 buckets alive.
func bucketing(nums []int, indexDiff int, valueDiff int) bool {
	buckets := make(map[int]int) // bucket id → the single value stored there
	width := valueDiff + 1       // bucket width so one bucket spans valueDiff+1 values

	// getBucket computes a floor-division bucket id that also works for negatives.
	getBucket := func(v int) int {
		if v >= 0 {
			return v / width
		}
		return (v+1)/width - 1 // arithmetic floor division for negative v
	}

	for i, v := range nums {
		b := getBucket(v)
		if _, ok := buckets[b]; ok { // same bucket → guaranteed within valueDiff
			return true
		}
		// Adjacent buckets can still hold a value within valueDiff — verify.
		if x, ok := buckets[b-1]; ok && abs(v-x) <= valueDiff {
			return true
		}
		if x, ok := buckets[b+1]; ok && abs(v-x) <= valueDiff {
			return true
		}
		buckets[b] = v // store v (at most one value per bucket in the window)
		if i >= indexDiff {
			// The value at index i-indexDiff is leaving the window; drop its bucket.
			delete(buckets, getBucket(nums[i-indexDiff]))
		}
	}
	return false
}

// abs returns the absolute value of x.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(bruteForce([]int{1, 2, 3, 1}, 3, 0))       // expected true
	fmt.Println(bruteForce([]int{1, 5, 9, 1, 5, 9}, 2, 3)) // expected false
	fmt.Println(bruteForce([]int{-3, 3}, 2, 4))            // expected false (|−3−3|=6>4)

	fmt.Println("=== Approach 2: Sliding Window + Sorted Structure ===")
	fmt.Println(slidingWindowSorted([]int{1, 2, 3, 1}, 3, 0))       // expected true
	fmt.Println(slidingWindowSorted([]int{1, 5, 9, 1, 5, 9}, 2, 3)) // expected false
	fmt.Println(slidingWindowSorted([]int{-3, 3}, 2, 4))            // expected false

	fmt.Println("=== Approach 3: Bucketing (Optimal) ===")
	fmt.Println(bucketing([]int{1, 2, 3, 1}, 3, 0))       // expected true
	fmt.Println(bucketing([]int{1, 5, 9, 1, 5, 9}, 2, 3)) // expected false
	fmt.Println(bucketing([]int{-3, 3}, 2, 4))            // expected false (checks negative-value floor bucketing)
}
