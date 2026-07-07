# 0295 — Find Median from Data Stream

> LeetCode #295 · Difficulty: Hard
> **Categories:** Heap (Priority Queue), Design, Two Heaps, Sorting, Data Stream

---

## Problem Statement

The **median** is the middle value in an ordered integer list. If the size of the list is even, there is no middle value, and the median is the mean of the two middle values.

- For example, for `arr = [2,3,4]`, the median is `3`.
- For example, for `arr = [2,3]`, the median is `(2 + 3) / 2 = 2.5`.

Implement the `MedianFinder` class:

- `MedianFinder()` initializes the `MedianFinder` object.
- `void addNum(int num)` adds the integer `num` from the data stream to the data structure.
- `double findMedian()` returns the median of all elements so far. Answers within `10⁻⁵` of the actual answer will be accepted.

**Example 1:**
```
Input
["MedianFinder", "addNum", "addNum", "findMedian", "addNum", "findMedian"]
[[], [1], [2], [], [3], []]
Output
[null, null, null, 1.5, null, 2.0]

Explanation
MedianFinder medianFinder = new MedianFinder();
medianFinder.addNum(1);    // arr = [1]
medianFinder.addNum(2);    // arr = [1, 2]
medianFinder.findMedian(); // return 1.5 (i.e., (1 + 2) / 2)
medianFinder.addNum(3);    // arr[1, 2, 3]
medianFinder.findMedian(); // return 2.0
```

**Constraints:**
- `-10⁵ <= num <= 10⁵`
- There will be at least one element in the data structure before calling `findMedian`.
- At most `5 * 10⁴` calls will be made to `addNum` and `findMedian`.

**Follow-up:**
- If all integer numbers from the stream are in the range `[0, 100]`, how would you optimize your solution? (Use a bucket count array — O(100) per query.)
- If 99% of all integer numbers from the stream are in the range `[0, 100]`, how would you optimize your solution? (Bucket the 99%, keep two overflow lists/heaps for the outliers.)

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★★ Very High  | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★★☆ High       | 2024          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used
- **Heap / Priority Queue** — two heaps straddling the median → see [`/dsa/heap_priority_queue.md`](/dsa/heap_priority_queue.md)
- **Design Data Structures** — build a stateful class with amortized operations → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)
- **Sorting / Binary Search insertion** — sorted-slice baseline → see [`/dsa/sorting.md`](/dsa/sorting.md)

---

## Approaches Overview
| # | Approach | addNum | findMedian | Space | When to use |
|---|----------|--------|-----------|-------|-------------|
| 1 | Sorted Slice | O(n) | O(1) | O(n) | Baseline; few inserts |
| 2 | Two Heaps (Optimal) | O(log n) | O(1) | O(n) | Many mixed operations |

---

## Approach 1 — Sorted Slice (Brute Force)

### Intuition
If the data is always sorted, the median is trivial: the middle element (odd count) or the mean of the two middles (even count). The cost is pushed to insertion — each `addNum` must place the value at its sorted position, shifting the tail.

### Algorithm
1. `addNum`: binary-search the insertion index with `sort.SearchInts`, then splice the value in (shifting the tail right).
2. `findMedian`: read the middle element, or average the two central elements for even counts.

### Complexity
- **Time:** `addNum` O(n) (the shift dominates the O(log n) search); `findMedian` O(1).
- **Space:** O(n) — all numbers are retained.

### Code
```go
type SortedSliceFinder struct {
	nums []int // kept in non-decreasing order after every AddNum
}

func NewSortedSliceFinder() *SortedSliceFinder { return &SortedSliceFinder{} }

func (f *SortedSliceFinder) AddNum(num int) {
	i := sort.SearchInts(f.nums, num)
	f.nums = append(f.nums, 0)     // grow by one (value overwritten below)
	copy(f.nums[i+1:], f.nums[i:]) // shift the tail right to open a gap at i
	f.nums[i] = num                // drop num into its sorted slot
}

func (f *SortedSliceFinder) FindMedian() float64 {
	n := len(f.nums)
	if n%2 == 1 {
		return float64(f.nums[n/2]) // single middle element
	}
	return float64(f.nums[n/2-1]+f.nums[n/2]) / 2.0
}
```

### Dry Run
Ops: `addNum(1), addNum(2), findMedian, addNum(3), findMedian`

| Op | search idx | nums after | median |
|----|-----------|------------|--------|
| addNum(1) | 0 | `[1]` | — |
| addNum(2) | 1 | `[1,2]` | — |
| findMedian | — | `[1,2]` (even) | `(1+2)/2 = 1.5` |
| addNum(3) | 2 | `[1,2,3]` | — |
| findMedian | — | `[1,2,3]` (odd) | `nums[1] = 2.0` |

Output: `[null,null,null,1.5,null,2.0]` ✓

---

## Approach 2 — Two Heaps (Optimal)

### Intuition
Split the sorted stream at the middle into a **low** half (a max-heap, so its top is the largest small value) and a **high** half (a min-heap, so its top is the smallest large value). Keep the sizes balanced (low equals high, or has exactly one more). Then the median is either low's top (odd total) or the average of the two tops (even total) — both O(1) to read, with O(log n) inserts.

### Algorithm
1. `addNum`: push onto `low`; move `low`'s top into `high`; if `high` is now larger than `low`, move `high`'s top back to `low`. This routes the value to the correct side AND rebalances sizes in one pattern.
2. `findMedian`: if `low` is larger, return its top; otherwise average the two tops.

### Complexity
- **Time:** `addNum` O(log n) (a constant number of heap pushes/pops); `findMedian` O(1).
- **Space:** O(n) — every element lives in one of the two heaps.

### Code
```go
type TwoHeapFinder struct {
	low  *maxHeap // smaller half; top = largest of the small values
	high *minHeap // larger half;  top = smallest of the large values
}

func NewTwoHeapFinder() *TwoHeapFinder {
	return &TwoHeapFinder{low: &maxHeap{}, high: &minHeap{}}
}

func (f *TwoHeapFinder) AddNum(num int) {
	heap.Push(f.low, num)              // tentatively add to the low half
	heap.Push(f.high, heap.Pop(f.low)) // shift low's max into high (keeps order)
	if f.high.Len() > f.low.Len() {    // high grew too large...
		heap.Push(f.low, heap.Pop(f.high)) // ...move its min back to low
	}
}

func (f *TwoHeapFinder) FindMedian() float64 {
	if f.low.Len() > f.high.Len() {
		return float64((*f.low)[0]) // odd total: low holds the extra element
	}
	return float64((*f.low)[0]+(*f.high)[0]) / 2.0
}
```

### Dry Run
Ops: `addNum(1), addNum(2), findMedian, addNum(3), findMedian`

| Op | after push/pop dance | low (max-heap) | high (min-heap) | median |
|----|----------------------|----------------|-----------------|--------|
| addNum(1) | push 1→low, move to high, high>low → move back | `[1]` | `[]` | — |
| addNum(2) | push 2→low(`[2,1]`), move 2→high, sizes equal | `[1]` | `[2]` | — |
| findMedian | sizes equal | `[1]` | `[2]` | `(1+2)/2 = 1.5` |
| addNum(3) | push 3→low(`[3,1]`), move 3→high(`[2,3]`), high>low → move 2 back | `[2,1]` | `[3]` | — |
| findMedian | low bigger | `[2,1]` | `[3]` | `low[0] = 2.0` |

Output: `[null,null,null,1.5,null,2.0]` ✓

---

## Key Takeaways
- **Two heaps straddling the median** is the canonical "running median / running k-th" pattern: max-heap for the low half, min-heap for the high half, sizes kept within 1.
- The push→pop→rebalance trio (`push to low → pop into high → maybe pop back`) both places the value correctly and keeps sizes balanced without branching on comparisons.
- Go's `container/heap` needs a `sort.Interface` + `Push`/`Pop`; flip `Less` to turn a min-heap into a max-heap.
- Range-bounded follow-ups replace heaps with a **bucket-count array** for O(1) amortized updates and O(range) median queries.

---

## Related Problems
- LeetCode #480 — Sliding Window Median (two heaps with lazy deletion)
- LeetCode #4 — Median of Two Sorted Arrays (static median)
- LeetCode #703 — Kth Largest Element in a Stream (single heap)
- LeetCode #295 variants — running percentile / order statistics
