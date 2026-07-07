package main

import (
	"container/heap"
	"fmt"
	"sort"
)

// ── Approach 1: Count + Full Sort (Brute Force) ──────────────────────────────
//
// sortByFrequency counts every element, then sorts the distinct values by
// frequency and takes the top k.
//
// Intuition:
//
//	"Top k frequent" literally means: tally how often each value appears, then
//	rank the distinct values by that tally and slice off the first k. A hash map
//	gives the tallies; a comparator sort gives the ranking.
//
// Algorithm:
//  1. Build freq[value] = count with one pass.
//  2. Collect the distinct keys into a slice.
//  3. Sort keys by descending frequency.
//  4. Return the first k keys.
//
// Time:  O(n + d log d) where d = number of distinct values ≤ n; the sort
//
//	dominates in the worst case (all distinct) → O(n log n).
//
// Space: O(n) — the frequency map plus the key slice.
func sortByFrequency(nums []int, k int) []int {
	freq := make(map[int]int) // value → how many times it occurs
	for _, v := range nums {
		freq[v]++ // tally each occurrence
	}
	keys := make([]int, 0, len(freq)) // the distinct values
	for key := range freq {
		keys = append(keys, key)
	}
	// Rank distinct values by frequency, highest first.
	sort.Slice(keys, func(i, j int) bool {
		return freq[keys[i]] > freq[keys[j]]
	})
	return keys[:k] // the k most frequent
}

// ── Approach 2: Min-Heap of Size k (Optimal for small k) ─────────────────────
//
// minHeapTopK keeps a heap of the k highest-frequency values seen so far,
// evicting the smallest whenever the heap exceeds size k.
//
// Intuition:
//
//	We do not have to sort ALL distinct values — we only need the top k. Stream
//	the (value, freq) pairs through a min-heap capped at size k: whenever it
//	overflows, pop the smallest frequency. What survives are exactly the k
//	largest. Each push/pop costs O(log k), far cheaper than O(log n) when
//	k ≪ n.
//
// Algorithm:
//  1. Build the frequency map.
//  2. Push each (value, freq) onto a min-heap ordered by freq.
//  3. If the heap size exceeds k, pop the minimum.
//  4. Drain the heap into the result (order not required by the problem).
//
// Time:  O(n + d log k) — building the map is O(n); each of d distinct values
//
//	does O(log k) heap work.
//
// Space: O(n) — the frequency map; the heap holds at most k entries.
type freqItem struct {
	val   int // the element
	count int // its frequency
}

// freqHeap is a MIN-heap ordered by count: the smallest frequency sits on top,
// so it is the first to be evicted when the heap overflows size k.
type freqHeap []freqItem

func (h freqHeap) Len() int            { return len(h) }
func (h freqHeap) Less(i, j int) bool  { return h[i].count < h[j].count } // min by count
func (h freqHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *freqHeap) Push(x interface{}) { *h = append(*h, x.(freqItem)) }
func (h *freqHeap) Pop() interface{} {
	old := *h
	n := len(old)
	it := old[n-1] // the last element is the one heap.Pop moved to the end
	*h = old[:n-1] // shrink the slice
	return it
}

func minHeapTopK(nums []int, k int) []int {
	freq := make(map[int]int) // value → count
	for _, v := range nums {
		freq[v]++
	}
	h := &freqHeap{} // min-heap capped at size k
	heap.Init(h)
	for val, count := range freq {
		heap.Push(h, freqItem{val: val, count: count}) // add candidate
		if h.Len() > k {                               // heap too big?
			heap.Pop(h) // drop the smallest frequency currently held
		}
	}
	res := make([]int, h.Len())
	for i := range res {
		res[i] = heap.Pop(h).(freqItem).val // drain remaining k values
	}
	return res
}

// ── Approach 3: Bucket Sort by Frequency (Optimal, linear) ───────────────────
//
// bucketSort groups values by their exact frequency into buckets indexed by
// count, then scans from the highest possible frequency downward.
//
// Intuition:
//
//	A value's frequency is an integer in [1, n]. So make n+1 buckets, put each
//	value into bucket[frequency], and walk the buckets from the top (frequency
//	n) down, collecting values until we have k. No comparison sort needed — the
//	frequency IS the index, which is a counting/bucket sort in disguise.
//
// Algorithm:
//  1. Build the frequency map.
//  2. Create buckets[0..n]; append value to buckets[freq[value]].
//  3. Iterate frequency from n down to 1, collecting values until k are found.
//
// Time:  O(n) — counting is O(n), and scanning n+1 buckets visits each distinct
//
//	value once.
//
// Space: O(n) — the frequency map and the buckets.
func bucketSort(nums []int, k int) []int {
	freq := make(map[int]int) // value → count
	for _, v := range nums {
		freq[v]++
	}
	// buckets[f] holds every value whose frequency is exactly f.
	// The max possible frequency is len(nums), so we need len+1 slots.
	buckets := make([][]int, len(nums)+1)
	for val, count := range freq {
		buckets[count] = append(buckets[count], val) // drop value into its bucket
	}
	res := make([]int, 0, k)
	// Walk from the highest frequency downward — most frequent first.
	for f := len(buckets) - 1; f >= 1 && len(res) < k; f-- {
		for _, val := range buckets[f] { // every value with this exact frequency
			res = append(res, val)
			if len(res) == k { // collected enough
				return res
			}
		}
	}
	return res
}

// sortInts returns a sorted copy so output order is stable for verification.
func sortInts(s []int) []int {
	out := append([]int(nil), s...)
	sort.Ints(out)
	return out
}

func main() {
	// Example 1: nums=[1,1,1,2,2,3], k=2 → [1,2]
	// Example 2: nums=[1], k=1 → [1]
	// The problem allows any order; we sort each result before printing so the
	// output is deterministic and comparable to the expected set.

	fmt.Println("=== Approach 1: Count + Full Sort (Brute Force) ===")
	fmt.Println(sortInts(sortByFrequency([]int{1, 1, 1, 2, 2, 3}, 2))) // expected [1 2]
	fmt.Println(sortInts(sortByFrequency([]int{1}, 1)))                // expected [1]

	fmt.Println("=== Approach 2: Min-Heap of Size k ===")
	fmt.Println(sortInts(minHeapTopK([]int{1, 1, 1, 2, 2, 3}, 2))) // expected [1 2]
	fmt.Println(sortInts(minHeapTopK([]int{1}, 1)))                // expected [1]

	fmt.Println("=== Approach 3: Bucket Sort by Frequency (Optimal) ===")
	fmt.Println(sortInts(bucketSort([]int{1, 1, 1, 2, 2, 3}, 2))) // expected [1 2]
	fmt.Println(sortInts(bucketSort([]int{1}, 1)))                // expected [1]
}
