# 0215 — Kth Largest Element in an Array

> LeetCode #215 · Difficulty: Medium
> **Categories:** Array, Heap (Priority Queue), Divide and Conquer, Quickselect, Sorting

---

## Problem Statement

Given an integer array `nums` and an integer `k`, return *the* `kᵗʰ` *largest element in the array*.

Note that it is the `kᵗʰ` largest element in the sorted order, not the `kᵗʰ` distinct element.

Can you solve it without sorting?

**Example 1:**
```
Input: nums = [3,2,1,5,6,4], k = 2
Output: 5
```

**Example 2:**
```
Input: nums = [3,2,3,1,2,4,5,5,6], k = 4
Output: 4
```

**Constraints:**
- `1 <= k <= nums.length <= 10⁵`
- `-10⁴ <= nums[i] <= 10⁴`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★★ Very High  | 2024          |
| Meta       | ★★★★☆ High       | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2024          |
| Apple      | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Heap / Priority Queue** — a size-`k` min-heap keeps the k largest seen so far; its root is the k-th largest → see [`/dsa/heap_priority_queue.md`](/dsa/heap_priority_queue.md)
- **Quickselect** — partition around a pivot and recurse only into the side holding the target index; average O(n) selection → see [`/dsa/quickselect.md`](/dsa/quickselect.md)
- **Divide and Conquer** — quickselect is a one-sided quicksort → see [`/dsa/divide_and_conquer.md`](/dsa/divide_and_conquer.md)
- **Sorting** — the baseline: sort and index → see [`/dsa/sorting.md`](/dsa/sorting.md)

---

## Approaches Overview

Let n = `len(nums)`.

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Sort then index | O(n log n) | O(n) | Simplest; fine when n is small |
| 2 | Min-heap of size k | O(n log k) | O(k) | k ≪ n, or streaming/online data |
| 3 | Quickselect (Optimal) | O(n) avg, O(n²) worst | O(1) | "Without sorting" — best average time |

---

## Approach 1 — Sort then Index

### Intuition
The k-th largest is a positional statistic. Sort the array ascending and it sits at index `n-k`. Obviously correct; the only downside is paying for a full sort when we need a single position.

### Algorithm
1. Copy `nums` (so the caller's array is untouched) and sort ascending.
2. Return `copy[n-k]`.

### Complexity
- **Time:** O(n log n) — dominated by the sort.
- **Space:** O(n) — the sorted copy.

### Code
```go
func sortIndex(nums []int, k int) int {
	cp := append([]int(nil), nums...)
	sort.Ints(cp)
	return cp[len(cp)-k]
}
```

### Dry Run (Example 1: `nums = [3,2,1,5,6,4]`, k = 2)

| Step | State |
|------|-------|
| copy | `[3,2,1,5,6,4]` |
| sort ascending | `[1,2,3,4,5,6]` |
| index `n-k = 6-2 = 4` | `cp[4] = 5` |

Return **5** ✓

---

## Approach 2 — Min-Heap of Size k

### Intuition
Only the k biggest values matter. Maintain a min-heap capped at size k: the smallest of the current top-k sits at the root. Push each value; whenever the heap exceeds k, pop the minimum (it cannot be among the k largest). After all elements, the heap holds exactly the k largest and its root — the minimum of those — is the k-th largest overall. Ideal when k is small or data arrives as a stream.

### Algorithm
1. Push each number onto the heap.
2. Whenever `heap.size > k`, pop the smallest.
3. Return the heap root.

### Complexity
- **Time:** O(n log k) — n pushes/pops on a heap of size ≤ k.
- **Space:** O(k) for the heap.

### Code
```go
type intHeap []int

func (h intHeap) Len() int            { return len(h) }
func (h intHeap) Less(i, j int) bool  { return h[i] < h[j] } // min-heap
func (h intHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *intHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *intHeap) Pop() interface{} {
	old := *h
	n := len(old)
	v := old[n-1]
	*h = old[:n-1]
	return v
}

func minHeapK(nums []int, k int) int {
	h := &intHeap{}
	heap.Init(h)
	for _, v := range nums {
		heap.Push(h, v)
		if h.Len() > k {
			heap.Pop(h) // drop the smallest → keep only the k largest
		}
	}
	return (*h)[0]
}
```

### Dry Run (Example 1: `nums = [3,2,1,5,6,4]`, k = 2)

| Push | Heap (min-root) | Size > k? Pop | Heap after |
|------|-----------------|---------------|------------|
| 3 | [3] | no | [3] |
| 2 | [2,3] | no | [2,3] |
| 1 | [1,3,2] | yes → pop 1 | [2,3] |
| 5 | [2,3,5] | yes → pop 2 | [3,5] |
| 6 | [3,5,6] | yes → pop 3 | [5,6] |
| 4 | [4,6,5] | yes → pop 4 | [5,6] |

Root of final heap `[5,6]` = **5** ✓

---

## Approach 3 — Quickselect (Optimal)

### Intuition
Full sorting is wasteful: we only need the element at ascending index `target = n-k`. Quickselect partitions the array so the pivot lands at its final sorted position `p` (everything left ≤ pivot, right ≥ pivot). If `p == target` we're done; otherwise recurse into just the side containing `target`. On average each partition roughly halves the remaining work → O(n).

### Algorithm
1. `target = n - k` (ascending index of the k-th largest).
2. `partition(lo, hi)`: Lomuto scheme with pivot `a[hi]`; returns the pivot's final index `p`.
3. Loop: if `p == target` return `a[p]`; if `p < target` set `lo = p+1`; else `hi = p-1`.

### Complexity
- **Time:** O(n) average; O(n²) worst case with adversarial pivots (mitigated in practice by random/median-of-three pivots).
- **Space:** O(1) extra — in-place partitioning, iterative loop.

### Code
```go
func quickselect(nums []int, k int) int {
	a := append([]int(nil), nums...)
	target := len(a) - k
	lo, hi := 0, len(a)-1
	for lo <= hi {
		p := partition(a, lo, hi)
		switch {
		case p == target:
			return a[p]
		case p < target:
			lo = p + 1
		default:
			hi = p - 1
		}
	}
	return -1
}

func partition(a []int, lo, hi int) int {
	pivot := a[hi]
	i := lo
	for j := lo; j < hi; j++ {
		if a[j] <= pivot {
			a[i], a[j] = a[j], a[i]
			i++
		}
	}
	a[i], a[hi] = a[hi], a[i]
	return i
}
```

### Dry Run (Example 1: `nums = [3,2,1,5,6,4]`, k = 2 → target = 4)

| lo | hi | pivot = a[hi] | array after partition | p | branch |
|----|----|---------------|-----------------------|---|--------|
| 0 | 5 | 4 | `[3,2,1,4,6,5]` | 3 | p(3) < target(4) → lo = 4 |
| 4 | 5 | 5 | `[3,2,1,4,5,6]` | 4 | p(4) == target → return a[4] |

`a[4] = 5` → return **5** ✓

---

## Key Takeaways

- **"k-th largest" = ascending index `n-k`.** Convert once and the problem becomes a selection at a known index.
- **Three tiers of effort:** sort O(n log n) → heap O(n log k) → quickselect O(n) average. Pick by constraints: small n → sort; streaming or tiny k → heap; "without sorting" / best average → quickselect.
- **Min-heap keeps the *largest* k** (root = the smallest survivor = the answer); a max-heap of all n and k pops is O(n + k log n) — worse when k is large but n huge is the other way, so know both.
- **Quickselect is one-sided quicksort:** recurse into only the partition containing the target index. Randomising the pivot avoids the O(n²) sorted-input trap.
- The heap `Push`/`Pop` interface signatures (`interface{}`, pointer receiver mutating the slice) are the standard Go `container/heap` boilerplate worth memorising.

---

## Related Problems

- LeetCode #973 — K Closest Points to Origin (heap / quickselect on distance)
- LeetCode #347 — Top K Frequent Elements (heap / bucket / quickselect)
- LeetCode #703 — Kth Largest Element in a Stream (size-k min-heap, online)
- LeetCode #692 — Top K Frequent Words (heap with tie-breaking)
- LeetCode #4 — Median of Two Sorted Arrays (selection by index)
- LeetCode #324 — Wiggle Sort II (quickselect + partition)
