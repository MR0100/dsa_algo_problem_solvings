package main

import (
	"fmt"
	"math"
	"sort"
)

// ── Approach 1: Brute Force (Sort + Deduplicate) ─────────────────────────────
//
// bruteForce solves Third Maximum Number by sorting the distinct values in
// descending order and reading the third one (or the first if fewer than three
// distinct values exist).
//
// Intuition:
//
//	"Third distinct maximum" is literally the 3rd element of the sorted-descending
//	list of *distinct* values. So collect the distinct values, sort them big→small,
//	and index. If there are fewer than 3 distinct values, the rule says return the
//	overall maximum, which is index 0.
//
// Algorithm:
//  1. Put nums into a set to drop duplicates.
//  2. Copy the set to a slice and sort descending.
//  3. If len >= 3 return slice[2], else return slice[0].
//
// Time:  O(n log n) — dominated by the sort of up to n distinct values.
// Space: O(n) — the set and the deduplicated slice.
func bruteForce(nums []int) int {
	seen := make(map[int]struct{}, len(nums)) // set of distinct values
	for _, v := range nums {
		seen[v] = struct{}{} // insert; duplicates collapse automatically
	}
	distinct := make([]int, 0, len(seen))
	for v := range seen {
		distinct = append(distinct, v) // materialise the distinct values
	}
	sort.Sort(sort.Reverse(sort.IntSlice(distinct))) // descending order
	if len(distinct) >= 3 {
		return distinct[2] // the 3rd distinct maximum
	}
	return distinct[0] // fewer than 3 distinct → return the maximum
}

// ── Approach 2: Three-Variable Scan (Optimal) ────────────────────────────────
//
// threeVariableScan solves Third Maximum Number in a single pass, tracking the
// top three DISTINCT values seen so far without any sorting.
//
// Intuition:
//
//	We only need the top three distinct values, so keep three "podium" slots
//	first > second > third. For each number: skip it if it already equals one of
//	the slots (must stay distinct); otherwise slot it into the podium, cascading
//	the smaller ones down. Using *int pointers (nil = "slot empty") avoids the
//	classic sentinel bug where nums[i] can legitimately be math.MinInt64 / -2^31.
//
// Algorithm:
//  1. first = second = third = nil (empty).
//  2. For each v: if v equals any filled slot, skip (distinctness).
//  3. Else compare against first/second/third and insert, shifting lower slots down.
//  4. If third is filled, return *third; else return *first (fewer than 3 distinct).
//
// Time:  O(n) — one pass, O(1) work per element.
// Space: O(1) — three pointer slots.
func threeVariableScan(nums []int) int {
	var first, second, third *int // podium slots; nil means "not yet filled"
	for i := range nums {
		v := nums[i]
		// Skip duplicates of any already-placed podium value: the three
		// maxima must be DISTINCT.
		if (first != nil && v == *first) ||
			(second != nil && v == *second) ||
			(third != nil && v == *third) {
			continue
		}
		switch {
		case first == nil || v > *first: // new overall maximum
			third = second // everyone slides down one place
			second = first
			nv := v // take an address of a fresh copy (loop var reuse safety)
			first = &nv
		case second == nil || v > *second: // fits between 1st and 2nd
			third = second
			nv := v
			second = &nv
		case third == nil || v > *third: // fits into 3rd place
			nv := v
			third = &nv
		}
	}
	if third != nil { // a genuine third distinct maximum exists
		return *third
	}
	return *first // fewer than 3 distinct values → the maximum
}

// ── Approach 3: Bounded Min-Set (Ordered Top-3) ──────────────────────────────
//
// boundedSet solves Third Maximum Number by maintaining a small sorted set of at
// most three distinct values, evicting the smallest whenever a fourth arrives.
//
// Intuition:
//
//	Keep a set capped at size 3. Insert each distinct value; if the set exceeds 3,
//	drop its minimum. After the scan the set holds the three largest distinct
//	values (or all of them if fewer). This mirrors how a size-limited min-heap or
//	TreeSet is used for "top-k" — here k = 3, so a tiny sorted slice is simplest
//	and keeps everything O(1) per step.
//
// Algorithm:
//  1. Maintain a set (map) and a sorted view of its (≤3) elements.
//  2. For each v: skip if already present; else insert; if size > 3 remove the min.
//  3. If final size >= 3 return the min of the set (the 3rd max); else return the max.
//
// Time:  O(n) — each step does O(1) work (the set never exceeds 3 elements).
// Space: O(1) — at most three tracked values.
func boundedSet(nums []int) int {
	top := make(map[int]struct{}, 4) // distinct candidates, capped at 3
	for _, v := range nums {
		if _, ok := top[v]; ok {
			continue // already a candidate → keep distinctness
		}
		top[v] = struct{}{}
		if len(top) > 3 { // one too many → evict the smallest
			minV := math.MaxInt64
			for k := range top {
				if k < minV {
					minV = k
				}
			}
			delete(top, minV)
		}
	}
	// Read out either the min (3rd max) or the max (fewer than 3 distinct).
	if len(top) >= 3 {
		minV := math.MaxInt64
		for k := range top {
			if k < minV {
				minV = k // smallest of the top 3 = the 3rd maximum
			}
		}
		return minV
	}
	maxV := math.MinInt64
	for k := range top {
		if k > maxV {
			maxV = k // fewer than 3 distinct → overall maximum
		}
	}
	return maxV
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Sort + Deduplicate) ===")
	fmt.Println(bruteForce([]int{3, 2, 1}))    // expected 1
	fmt.Println(bruteForce([]int{1, 2}))       // expected 2
	fmt.Println(bruteForce([]int{2, 2, 3, 1})) // expected 1

	fmt.Println("=== Approach 2: Three-Variable Scan (Optimal) ===")
	fmt.Println(threeVariableScan([]int{3, 2, 1}))    // expected 1
	fmt.Println(threeVariableScan([]int{1, 2}))       // expected 2
	fmt.Println(threeVariableScan([]int{2, 2, 3, 1})) // expected 1

	fmt.Println("=== Approach 3: Bounded Min-Set (Ordered Top-3) ===")
	fmt.Println(boundedSet([]int{3, 2, 1}))    // expected 1
	fmt.Println(boundedSet([]int{1, 2}))       // expected 2
	fmt.Println(boundedSet([]int{2, 2, 3, 1})) // expected 1
}
