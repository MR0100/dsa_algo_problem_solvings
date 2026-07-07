# 0347 — Top K Frequent Elements

> LeetCode #347 · Difficulty: Medium
> **Categories:** Array, Hash Table, Heap (Priority Queue), Bucket Sort, Counting, Divide and Conquer, Sorting, Quickselect

---

## Problem Statement

Given an integer array `nums` and an integer `k`, return the `k` most frequent elements. You may return the answer in **any order**.

**Example 1:**
```
Input: nums = [1,1,1,2,2,3], k = 2
Output: [1,2]
```

**Example 2:**
```
Input: nums = [1], k = 1
Output: [1]
```

**Constraints:**
- `1 <= nums.length <= 10⁵`
- `-10⁴ <= nums[i] <= 10⁴`
- `k` is in the range `[1, the number of unique elements in the array]`.
- It is **guaranteed** that the answer is **unique**.

**Follow-up:** Your algorithm's time complexity must be better than `O(n log n)`, where `n` is the array's size.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★★ Very High  | 2024          |
| Facebook   | ★★★★★ Very High  | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2024          |
| Apple      | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |
| Uber       | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used
- **Hash Map (frequency counting)** — a single pass tallies each value's count → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Heap / Priority Queue** — a size-`k` min-heap keeps the k largest frequencies for O(n log k) → see [`/dsa/heap_priority_queue.md`](/dsa/heap_priority_queue.md)
- **Bucket Sort / Counting Sort** — frequencies are integers in `[1, n]`, so they can index buckets for a linear scan → see [`/dsa/counting_sort.md`](/dsa/counting_sort.md)
- **Quickselect** — an alternative O(n)-average partition on frequencies (mentioned in Key Takeaways) → see [`/dsa/quickselect.md`](/dsa/quickselect.md)

---

## Approaches Overview
| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Count + Full Sort (Brute Force) | O(n log n) | O(n) | Simplest; fails the follow-up bound |
| 2 | Min-Heap of Size k | O(n + d log k) | O(n) | k ≪ number of distinct values |
| 3 | Bucket Sort by Frequency (Optimal) | O(n) | O(n) | Meets the sub-`n log n` follow-up |

*(d = number of distinct values ≤ n.)*

---

## Approach 1 — Count + Full Sort (Brute Force)

### Intuition
"Top k frequent" means: tally how often each value appears, rank the distinct values by that tally, and slice off the first k. A hash map gives the tallies; a comparator sort gives the ranking.

### Algorithm
1. Build `freq[value] = count` in one pass over `nums`.
2. Collect the distinct keys into a slice.
3. Sort the keys by **descending** frequency.
4. Return the first `k` keys.

### Complexity
- **Time:** O(n + d log d) — counting is O(n); sorting the d distinct keys is O(d log d), which is O(n log n) worst case (all distinct).
- **Space:** O(n) — the frequency map and the key slice.

### Code
```go
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
```

### Dry Run
`nums = [1,1,1,2,2,3]`, `k = 2`:

| step | state |
|------|-------|
| count | freq = {1:3, 2:2, 3:1} |
| keys | [1, 2, 3] (some order) |
| sort by freq desc | [1 (3×), 2 (2×), 3 (1×)] |
| take k=2 | **[1, 2]** |

---

## Approach 2 — Min-Heap of Size k

### Intuition
We do not have to sort **all** distinct values — only surface the top k. Stream each `(value, freq)` pair through a min-heap capped at size k. Whenever it overflows, pop the smallest frequency. What survives are exactly the k largest, and each push/pop costs O(log k) — cheaper than O(log n) when k ≪ n.

### Algorithm
1. Build the frequency map.
2. Push each `(value, freq)` onto a min-heap ordered by `freq`.
3. If the heap size exceeds k, pop the minimum (smallest frequency).
4. Drain the heap into the result slice.

### Complexity
- **Time:** O(n + d log k) — O(n) to count, then d distinct values each do O(log k) heap work.
- **Space:** O(n) — the frequency map; the heap holds at most k entries.

### Code
```go
type freqItem struct {
	val   int // the element
	count int // its frequency
}

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
```

### Dry Run
`nums = [1,1,1,2,2,3]`, `k = 2`, freq = {1:3, 2:2, 3:1} (iteration order varies; one possible order shown):

| push | heap (min on top) | size > k? | after eviction |
|------|-------------------|-----------|----------------|
| push (3,1) | [(3,1)]           | no        | [(3,1)]        |
| push (1,3) | [(3,1),(1,3)]     | no        | [(3,1),(1,3)]  |
| push (2,2) | [(3,1),(1,3),(2,2)] | yes (3>2) | pop min (3,1) → [(2,2),(1,3)] |

Drain heap → values `{2, 1}` → **[1, 2]** after sorting for display.

---

## Approach 3 — Bucket Sort by Frequency (Optimal, linear)

### Intuition
A value's frequency is an integer in `[1, n]`. So make `n+1` buckets and place each value into `bucket[frequency]`. Then walk the buckets from the highest frequency downward, collecting values until we have k. The frequency **is** the index — a counting/bucket sort in disguise — so no comparison sort is needed and the whole thing is linear, beating the follow-up bound.

### Algorithm
1. Build the frequency map.
2. Create `buckets[0..n]`; append each value to `buckets[freq[value]]`.
3. Iterate frequency `f` from `n` down to `1`, collecting values until k are gathered.

### Complexity
- **Time:** O(n) — counting is O(n); scanning n+1 buckets touches each distinct value exactly once.
- **Space:** O(n) — the frequency map and the buckets.

### Code
```go
func bucketSort(nums []int, k int) []int {
	freq := make(map[int]int) // value → count
	for _, v := range nums {
		freq[v]++
	}
	// buckets[f] holds every value whose frequency is exactly f.
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
```

### Dry Run
`nums = [1,1,1,2,2,3]`, `k = 2`, freq = {1:3, 2:2, 3:1}, `n = 6`:

| step | state |
|------|-------|
| build buckets | buckets[1]=[3], buckets[2]=[2], buckets[3]=[1], rest empty |
| f=6..4 | empty, skip |
| f=3 | value 1 → res=[1] (len 1 < 2) |
| f=2 | value 2 → res=[1,2] (len == k) → **return [1, 2]** |

---

## Key Takeaways
- **Frequency problems start with a hash-map tally** — everything else is just how you rank those tallies.
- **You rarely need a full sort for "top k".** A size-k heap turns O(n log n) into O(n log k); bucket sort exploits the bounded integer key (frequency ∈ `[1, n]`) for true O(n).
- **Bucket / counting sort works whenever the sort key is a small bounded integer** — use the key as the array index.
- **Quickselect** is another linear-average option: partition the distinct values around the k-th largest frequency (the same idea as LeetCode #215) — O(n) average, O(n²) worst.
- The problem explicitly allows **any output order**, which is exactly why the heap and bucket approaches are free to emit values in whatever order they finish.

---

## Related Problems
- LeetCode #215 — Kth Largest Element in an Array (heap / quickselect)
- LeetCode #692 — Top K Frequent Words (frequency + tie-break by lexical order)
- LeetCode #451 — Sort Characters By Frequency (bucket sort on char frequency)
- LeetCode #973 — K Closest Points to Origin (size-k heap / quickselect)
