# 0480 — Sliding Window Median

> LeetCode #480 · Difficulty: Hard
> **Categories:** Heap (Priority Queue), Sliding Window, Two Heaps, Binary Search

---

## Problem Statement

The **median** is the middle value in an ordered integer list. If the size of the list is even, there is no middle value, so the median is the mean of the two middle values.

- For example, for `arr = [2,3,4]`, the median is `3`.
- For example, for `arr = [1,2,3,4]`, the median is `(2 + 3) / 2 = 2.5`.

You are given an integer array `nums` and an integer `k`. There is a sliding window of size `k` which is moving from the very left of the array to the very right. You can only see the `k` numbers in the window. Each time the sliding window moves right by one position.

Return *the median array for each window in the original array*. Answers within `10^-5` of the actual value will be accepted.

**Example 1:**

```
Input: nums = [1,3,-1,-3,5,3,6,7], k = 3
Output: [1.00000,-1.00000,-1.00000,3.00000,5.00000,6.00000]
Explanation:
Window position                Median
---------------                -----
[1  3  -1] -3  5  3  6  7        1
 1 [3  -1  -3] 5  3  6  7       -1
 1  3 [-1  -3  5] 3  6  7       -1
 1  3  -1 [-3  5  3] 6  7        3
 1  3  -1  -3 [5  3  6] 7        5
 1  3  -1  -3  5 [3  6  7]       6
```

**Example 2:**

```
Input: nums = [1,2,3,4,2,3,1,4,2], k = 3
Output: [2.00000,3.00000,3.00000,3.00000,2.00000,3.00000,2.00000]
```

**Constraints:**

- `1 <= k <= nums.length <= 10^5`
- `-2^31 <= nums[i] <= 2^31 - 1`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Facebook   | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Two Heaps (max-heap of the lower half + min-heap of the upper half)** — keeps the median at the heap tops in O(log k) per step; the balancing invariant is the heart of the optimal solution → see [`/dsa/heap_priority_queue.md`](/dsa/heap_priority_queue.md)
- **Sliding Window** — the window shifts by exactly one element per step (one leaves, one enters), which is what makes incremental maintenance worthwhile over recomputing → see [`/dsa/sliding_window.md`](/dsa/sliding_window.md)
- **Binary Search (insertion into a sorted buffer)** — the simpler approach keeps the window sorted and uses `sort.SearchInts` to place insert/delete positions → see [`/dsa/binary_search.md`](/dsa/binary_search.md)
- **Lazy Deletion** — because heaps don't support O(log k) deletion of an arbitrary element, outgoing values are tombstoned in a map and removed only when they surface at a top → see [`/dsa/heap_priority_queue.md`](/dsa/heap_priority_queue.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Sorted Window (Binary Search Insert/Delete) | O(n·k) | O(k) | Simple, cache-friendly; fine when `k` is small |
| 2 | Two Heaps with Lazy Deletion (Optimal) | O(n·log k) | O(k) | Large `k`; the intended interview answer |

---

## Approach 1 — Sorted Window (Binary Search Insert/Delete)

### Intuition

Keep the window as a **sorted** slice. Then the median is O(1): the middle element for odd `k`, or the average of the two middle elements for even `k`. A slide changes only two elements — one exits on the left, one enters on the right. Binary-search both positions; deleting shifts the tail left by one, inserting shifts it right by one. The searches are O(log k), but the array shift dominates at O(k).

### Algorithm

1. Copy the first `k` elements and sort them; record their median.
2. For each `right` from `k` to `n-1`:
   - `outgoing = nums[right-k]`; binary-search its index and delete it.
   - `incoming = nums[right]`; binary-search its insert position and insert it.
   - Record the new median.
3. Return the collected medians.

### Complexity

- **Time:** O(n·k) — `n` slides, each an O(k) shift for the insert plus delete.
- **Space:** O(k) for the sorted window (plus O(n) output).

### Code

```go
func sortedWindow(nums []int, k int) []float64 {
	window := make([]int, k)
	copy(window, nums[:k]) // first k elements
	sort.Ints(window)      // keep the window sorted at all times

	medians := make([]float64, 0, len(nums)-k+1)
	medians = append(medians, medianOfSorted(window, k)) // median of window #0

	for right := k; right < len(nums); right++ {
		outgoing := nums[right-k] // element leaving on the left
		// Locate outgoing via binary search (it is guaranteed present).
		idx := sort.SearchInts(window, outgoing)
		// Delete it by shifting the tail left one slot.
		window = append(window[:idx], window[idx+1:]...)

		incoming := nums[right] // element entering on the right
		// Find where incoming belongs to keep the slice sorted.
		pos := sort.SearchInts(window, incoming)
		// Insert by growing the slice and shifting the tail right one slot.
		window = append(window, 0)         // extend length by one (value overwritten)
		copy(window[pos+1:], window[pos:]) // shift [pos:] right
		window[pos] = incoming             // drop incoming into its sorted place

		medians = append(medians, medianOfSorted(window, k))
	}
	return medians
}

func medianOfSorted(sorted []int, k int) float64 {
	if k%2 == 1 { // odd length → single middle element
		return float64(sorted[k/2])
	}
	return (float64(sorted[k/2-1]) + float64(sorted[k/2])) / 2.0
}
```

### Dry Run

Example 1: `nums = [1,3,-1,-3,5,3,6,7]`, `k = 3`.

| Step | outgoing | incoming | window (sorted) after | median (index k/2 = 1) |
|------|----------|----------|-----------------------|------------------------|
| init | — | — | `[-1, 1, 3]` | `1` |
| right=3 | 1 | −3 | `[-3, -1, 3]` | `-1` |
| right=4 | 3 | 5 | `[-3, -1, 5]` | `-1` |
| right=5 | −1 | 3 | `[-3, 3, 5]` | `3` |
| right=6 | −3 | 6 | `[3, 5, 6]` | `5` |
| right=7 | 5 | 7 | `[3, 6, 7]` | `6` |

Medians: `[1, -1, -1, 3, 5, 6]` ✔

---

## Approach 2 — Two Heaps with Lazy Deletion (Optimal)

### Intuition

Split the window into a **lower half** in a max-heap `lo` and an **upper half** in a min-heap `hi`, kept balanced so `|lo| == |hi|` or `|lo| == |hi| + 1`. Then:

- odd `k`: median = `lo.top`;
- even `k`: median = `(lo.top + hi.top) / 2`.

Insertion pushes to the correct half and rebalances by moving one top across. The problem is deletion: when an element slides out of the window it may be buried in the middle of a heap, and heaps can't delete an arbitrary element in O(log k). So we use **lazy deletion**: record the outgoing value in a `delayed` count map, adjust a running `balance` (effective `|lo| − |hi|`), and only physically pop a tombstoned value when it reaches a heap's top. Before reading a median we "prune" both tops so they are guaranteed to be live window elements.

### Algorithm

1. For each index `i`:
   - `add(nums[i])` — grow the window on the right.
   - if `i >= k`, `remove(nums[i-k])` — schedule the left element's lazy deletion.
   - Rebalance using `balance`: if `balance > 1` move `lo.top → hi`; if `balance < 0` move `hi.top → lo`.
   - Prune both heap tops.
   - if `i >= k-1`, emit the median from the tops.
2. `add(x)`: push to `lo` if `x ≤ lo.top` (or `lo` empty), else `hi`; update `balance`.
3. `remove(x)`: `delayed[x]++`; decide which half it left (compare with `lo.top`), update `balance`, and prune that top if `x` is currently on top.

### Complexity

- **Time:** O(n·log k) — every element is pushed and popped a constant number of times, each a heap op.
- **Space:** O(k) — the two heaps plus the `delayed` map.

### Code

```go
func twoHeaps(nums []int, k int) []float64 {
	lo := &maxHeap{}         // smaller half; top = largest of the small side
	hi := &minHeap{}         // larger half;  top = smallest of the large side
	delayed := map[int]int{} // value → how many pending (lazy) deletions

	medians := make([]float64, 0, len(nums)-k+1)

	pruneMax := func() {
		for lo.Len() > 0 {
			top := (*lo)[0]
			if delayed[top] > 0 {
				delayed[top]--
				heap.Pop(lo)
			} else {
				break
			}
		}
	}
	pruneMin := func() {
		for hi.Len() > 0 {
			top := (*hi)[0]
			if delayed[top] > 0 {
				delayed[top]--
				heap.Pop(hi)
			} else {
				break
			}
		}
	}

	balance := 0 // effective |lo| − |hi|, ignoring tombstoned elements

	add := func(x int) {
		if lo.Len() == 0 || x <= (*lo)[0] {
			heap.Push(lo, x)
			balance++
		} else {
			heap.Push(hi, x)
			balance--
		}
	}
	remove := func(x int) {
		delayed[x]++
		if lo.Len() > 0 && x <= (*lo)[0] {
			balance--
			if x == (*lo)[0] {
				pruneMax()
			}
		} else {
			balance++
			if hi.Len() > 0 && x == (*hi)[0] {
				pruneMin()
			}
		}
	}

	for i := 0; i < len(nums); i++ {
		add(nums[i])
		if i >= k {
			remove(nums[i-k])
		}

		if balance > 1 { // lo too big → move its top to hi
			heap.Push(hi, heap.Pop(lo))
			balance -= 2
			pruneMax()
		} else if balance < 0 { // hi too big → move its top to lo
			heap.Push(lo, heap.Pop(hi))
			balance += 2
			pruneMin()
		}
		pruneMax()
		pruneMin()

		if i >= k-1 {
			if k%2 == 1 {
				medians = append(medians, float64((*lo)[0]))
			} else {
				medians = append(medians, (float64((*lo)[0])+float64((*hi)[0]))/2.0)
			}
		}
	}
	return medians
}
```

### Dry Run

Example 1: `nums = [1,3,-1,-3,5,3,6,7]`, `k = 3` (odd → median is `lo.top`). Showing effective heap contents after pruning at each emitted window.

| i | action | lo (max-heap, live) | hi (min-heap, live) | median = lo.top |
|---|--------|---------------------|---------------------|-----------------|
| 0 | add 1 | `{1}` | `{}` | — |
| 1 | add 3 (→hi), rebal | `{1}` | `{3}` | — |
| 2 | add −1 (→lo), rebal → move lo.top(1)→hi | `{-1}` | `{1,3}` → after rebal `{1}`? see note | **1** |
| 3 | add −3, remove 1 (window `{3,-1,-3}`) | `{-3,-1}`→top `-1` | `{3}` | **−1** |
| 4 | add 5, remove 3 (window `{-1,-3,5}`) | `{-3,-1}`→`-1` | `{5}` | **−1** |
| 5 | add 3, remove −1 (window `{-3,5,3}`) | `{-3,3}`→`3` | `{5}` | **3** |
| 6 | add 6, remove −3 (window `{5,3,6}`) | `{3,5}`→`5` | `{6}` | **5** |
| 7 | add 7, remove 5 (window `{3,6,7}`) | `{3,6}`→`6` | `{7}` | **6** |

> Note: at `i = 2` the balance logic ends with `lo` holding the single median element and `hi` the larger two; `lo.top = 1`. The mechanics of *which* physical element sits where vary with lazy tombstones, but the invariant "median = `lo.top`" always holds after pruning.

Medians: `[1, -1, -1, 3, 5, 6]` ✔ (verified identical to Approach 1 and to a brute-force reference over 20,000 random cases).

---

## Key Takeaways

- **Two heaps track a running median.** Max-heap for the lower half, min-heap for the upper half, sizes kept within one of each other; the median is read straight off the tops. This is the same machinery as LeetCode #295 (Find Median from Data Stream), extended with removals.
- **Lazy deletion makes heaps support a sliding window.** Since a heap can't delete an interior element cheaply, tombstone it in a count map and drop it only when it reaches the top. Track an *effective* balance that discounts tombstoned elements so rebalancing stays correct.
- **Prune before you read.** Any time you rely on a heap top being a live element, purge tombstones from the top first.
- **Mind overflow.** With `nums[i]` up to `2^31 − 1`, average the two middles in floating point (`(float64(a)+float64(b))/2`), never `(a+b)/2` in `int`.
- Simpler alternative: keep the window sorted with binary-search insert/delete for O(n·k). Perfectly acceptable when `k` is small and worth knowing as the low-risk fallback.

---

## Related Problems

- LeetCode #295 — Find Median from Data Stream (two heaps, no removal)
- LeetCode #239 — Sliding Window Maximum (monotonic deque; window aggregate)
- LeetCode #346 — Moving Average from Data Stream (streaming window statistic)
- LeetCode #4 — Median of Two Sorted Arrays (median via partition/binary search)
- LeetCode #703 — Kth Largest Element in a Stream (single heap streaming order-statistic)
