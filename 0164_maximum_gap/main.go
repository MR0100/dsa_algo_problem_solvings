package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Brute Force (Successor Search) ───────────────────────────────
//
// bruteForce solves Maximum Gap by finding, for every element, its successor
// in sorted order — the smallest value strictly greater than it — without
// ever sorting.
//
// Intuition:
//
//	In the sorted form, the element that follows v is the minimum of all
//	values strictly greater than v. The gap after v is (successor − v).
//	Computing that successor for every element by a full scan reproduces all
//	adjacent sorted gaps in O(n²) — no sorting, pure definition. (Duplicate
//	values create sorted gaps of 0, which can never beat a positive maximum,
//	so skipping them is safe.)
//
// Algorithm:
//  1. If n < 2, return 0.
//  2. For each element v:
//     a. Scan the array for the smallest value w with w > v.
//     b. If such w exists, candidate gap = w − v; keep the maximum.
//  3. Return the maximum candidate (0 if all values are equal).
//
// Time:  O(n²) — a full scan for each of the n elements.
// Space: O(1) — a few scalars.
func bruteForce(nums []int) int {
	if len(nums) < 2 {
		return 0 // fewer than two elements → no successive pair exists
	}
	maxGap := 0
	for _, v := range nums { // v plays the role of "left element of a sorted pair"
		successor := -1 // smallest value strictly greater than v; -1 = none found
		for _, w := range nums {
			if w > v && (successor == -1 || w < successor) {
				successor = w // tighter successor candidate
			}
		}
		if successor != -1 && successor-v > maxGap { // v is not the maximum value
			maxGap = successor - v
		}
	}
	return maxGap
}

// ── Approach 2: Sort and Scan ────────────────────────────────────────────────
//
// sortAndScan solves Maximum Gap by literally producing the sorted form and
// taking the largest adjacent difference.
//
// Intuition:
//
//	The problem statement is defined on the sorted array, so the obvious
//	solution is: sort, then one pass over adjacent pairs. Simple and correct,
//	but comparison sorting is O(n log n) — it violates the required linear
//	time bound, which is exactly what motivates Approaches 3 and 4.
//
// Algorithm:
//  1. If n < 2, return 0.
//  2. Sort a copy of the array.
//  3. Scan adjacent pairs and record the maximum difference.
//
// Time:  O(n log n) — dominated by the comparison sort.
// Space: O(n) — the defensive copy (O(log n)–O(n) also used by sort internals).
func sortAndScan(nums []int) int {
	if len(nums) < 2 {
		return 0 // no pair to compare
	}
	arr := make([]int, len(nums)) // copy so the caller's slice stays untouched
	copy(arr, nums)
	sort.Ints(arr) // produce the sorted form
	maxGap := 0
	for i := 1; i < len(arr); i++ {
		if gap := arr[i] - arr[i-1]; gap > maxGap { // adjacent sorted difference
			maxGap = gap
		}
	}
	return maxGap
}

// ── Approach 3: Radix Sort ───────────────────────────────────────────────────
//
// radixSort solves Maximum Gap in linear time by sorting with LSD radix sort
// (which beats the comparison-sort lower bound) and then scanning adjacents.
//
// Intuition:
//
//	The values are bounded (0 ≤ nums[i] ≤ 10⁹ < 2³²), so we can sort them
//	digit by digit in a fixed number of passes instead of comparing them.
//	Using base-256 digits, 4 stable counting-sort passes fully sort 32-bit
//	non-negative integers in O(4·n) = O(n) time — the linear-time
//	requirement is met, at the cost of O(n) auxiliary memory.
//
// Algorithm:
//  1. If n < 2, return 0.
//  2. For each 8-bit digit (shift = 0, 8, 16, 24):
//     a. Count occurrences of each of the 256 digit values.
//     b. Prefix-sum the counts into starting positions.
//     c. Stably place every element into the output buffer by its digit.
//     d. Swap buffers.
//  3. The array is now sorted — scan adjacent pairs for the maximum gap.
//
// Time:  O(d·(n + b)) with d = 4 passes, b = 256 buckets → O(n).
// Space: O(n + b) — one output buffer plus the 256-entry count table.
func radixSort(nums []int) int {
	if len(nums) < 2 {
		return 0 // no pair to compare
	}
	arr := make([]int, len(nums)) // working copy (also keeps input untouched)
	copy(arr, nums)
	buf := make([]int, len(arr))             // stable-scatter destination for each pass
	for shift := 0; shift < 32; shift += 8 { // 4 passes cover all 32 bits
		var counts [256]int
		for _, v := range arr {
			counts[(v>>shift)&0xFF]++ // histogram of the current 8-bit digit
		}
		pos := 0
		var starts [256]int
		for d := 0; d < 256; d++ { // exclusive prefix sums → first slot per digit
			starts[d] = pos
			pos += counts[d]
		}
		for _, v := range arr { // stable scatter: equal digits keep their order
			d := (v >> shift) & 0xFF
			buf[starts[d]] = v
			starts[d]++
		}
		arr, buf = buf, arr // sorted-by-this-digit buffer becomes the input
	}
	maxGap := 0
	for i := 1; i < len(arr); i++ {
		if gap := arr[i] - arr[i-1]; gap > maxGap { // adjacent sorted difference
			maxGap = gap
		}
	}
	return maxGap
}

// ── Approach 4: Bucket Sort + Pigeonhole (Optimal) ───────────────────────────
//
// bucketPigeonhole solves Maximum Gap in O(n) time without fully sorting, by
// exploiting the pigeonhole principle.
//
// Intuition:
//
//	With n numbers spanning [min, max], the n−1 sorted gaps must average
//	(max−min)/(n−1), so the maximum gap is at least ceil((max−min)/(n−1)).
//	Choose buckets of exactly that width: two numbers inside the same bucket
//	can then differ by at most bucketSize − 1 < maxGap, so the answer can
//	NEVER be an intra-bucket gap. Only gaps between the maximum of one
//	non-empty bucket and the minimum of the next non-empty bucket matter —
//	and those need just the per-bucket min and max, not a full sort.
//
// Algorithm:
//  1. If n < 2, return 0. Find min and max; if equal, return 0.
//  2. bucketSize = ceil((max − min) / (n − 1)); bucketCount = (max−min)/bucketSize + 1.
//  3. For every value, compute its bucket (v − min) / bucketSize and update
//     that bucket's min and max.
//  4. Sweep the buckets in order, keeping prevMax = max of the last
//     non-empty bucket; candidate gap = currentBucketMin − prevMax.
//  5. Return the largest candidate.
//
// Time:  O(n) — one pass to find min/max, one to fill buckets, one bucket sweep
//
//	(at most n+1 buckets).
//
// Space: O(n) — the two per-bucket arrays.
func bucketPigeonhole(nums []int) int {
	n := len(nums)
	if n < 2 {
		return 0 // no pair to compare
	}
	minV, maxV := nums[0], nums[0]
	for _, v := range nums { // one pass for the global extremes
		if v < minV {
			minV = v
		}
		if v > maxV {
			maxV = v
		}
	}
	if minV == maxV {
		return 0 // all elements equal → every sorted gap is 0
	}
	// Ceil division keeps bucketSize ≥ 1 and guarantees the answer is inter-bucket.
	bucketSize := (maxV - minV + n - 2) / (n - 1)
	bucketCount := (maxV-minV)/bucketSize + 1 // enough buckets to cover [minV, maxV]
	bucketMin := make([]int, bucketCount)     // per-bucket minimum
	bucketMax := make([]int, bucketCount)     // per-bucket maximum
	for i := range bucketMin {
		bucketMin[i] = -1 // -1 marks an empty bucket (values are ≥ 0 by constraint)
		bucketMax[i] = -1
	}
	for _, v := range nums {
		b := (v - minV) / bucketSize // bucket index of this value
		if bucketMin[b] == -1 || v < bucketMin[b] {
			bucketMin[b] = v
		}
		if bucketMax[b] == -1 || v > bucketMax[b] {
			bucketMax[b] = v
		}
	}
	maxGap := 0
	prevMax := minV // max of the last non-empty bucket seen (bucket 0 holds minV)
	for b := 0; b < bucketCount; b++ {
		if bucketMin[b] == -1 {
			continue // empty bucket — the gap simply spans across it
		}
		if gap := bucketMin[b] - prevMax; gap > maxGap { // inter-bucket gap
			maxGap = gap
		}
		prevMax = bucketMax[b] // this bucket's max feeds the next gap
	}
	return maxGap
}

func main() {
	examples := [][]int{
		{3, 6, 9, 1}, // expected 3 (sorted [1 3 6 9]: gaps 2, 3, 3)
		{10},         // expected 0 (fewer than two elements)
	}

	fmt.Println("=== Approach 1: Brute Force (Successor Search) ===")
	for _, ex := range examples {
		fmt.Printf("nums=%v  got=%d\n", ex, bruteForce(ex)) // expected 3, 0
	}

	fmt.Println("=== Approach 2: Sort and Scan ===")
	for _, ex := range examples {
		fmt.Printf("nums=%v  got=%d\n", ex, sortAndScan(ex)) // expected 3, 0
	}

	fmt.Println("=== Approach 3: Radix Sort ===")
	for _, ex := range examples {
		fmt.Printf("nums=%v  got=%d\n", ex, radixSort(ex)) // expected 3, 0
	}

	fmt.Println("=== Approach 4: Bucket Sort + Pigeonhole (Optimal) ===")
	for _, ex := range examples {
		fmt.Printf("nums=%v  got=%d\n", ex, bucketPigeonhole(ex)) // expected 3, 0
	}
}
