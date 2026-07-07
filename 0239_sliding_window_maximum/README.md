# 0239 — Sliding Window Maximum

> LeetCode #239 · Difficulty: Hard
> **Categories:** Array, Queue, Sliding Window, Heap (Priority Queue), Monotonic Queue

---

## Problem Statement

You are given an array of integers `nums`, there is a sliding window of size `k` which is moving from the very left of the array to the very right. You can only see the `k` numbers in the window. Each time the sliding window moves right by one position.

Return *the max sliding window*.

**Example 1:**

```
Input: nums = [1,3,-1,-3,5,3,6,7], k = 3
Output: [3,3,5,5,6,7]
Explanation:
Window position                Max
---------------               -----
[1  3  -1] -3  5  3  6  7       3
 1 [3  -1  -3] 5  3  6  7       3
 1  3 [-1  -3  5] 3  6  7       5
 1  3  -1 [-3  5  3] 6  7       5
 1  3  -1  -3 [5  3  6] 7       6
 1  3  -1  -3  5 [3  6  7]      7
```

**Example 2:**

```
Input: nums = [1], k = 1
Output: [1]
```

**Constraints:**

- `1 <= nums.length <= 10^5`
- `-10^4 <= nums[i] <= 10^4`
- `1 <= k <= nums.length`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★★ Very High  | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Facebook   | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2024          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |
| Uber       | ★★★☆☆ Medium     | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Monotonic Deque** — a double-ended queue of indices kept in decreasing value order, so the front is always the window maximum; this is the O(n) optimal → see [`/dsa/monotonic_stack.md`](/dsa/monotonic_stack.md)
- **Sliding Window** — the core "fixed-size window moving right by one" pattern → see [`/dsa/sliding_window.md`](/dsa/sliding_window.md)
- **Queue / Deque** — the deque supports O(1) push/pop at both ends → see [`/dsa/queue_deque.md`](/dsa/queue_deque.md)
- **Heap (Priority Queue)** — the alternative O(n log n) solution uses a max-heap with lazy deletion → see [`/dsa/heap_priority_queue.md`](/dsa/heap_priority_queue.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(n·k) | O(1) extra | Baseline; TLEs when n·k is large |
| 2 | Max-Heap (indices) | O(n log n) | O(n) | Good when you also need order statistics; simpler to reason about |
| 3 | Monotonic Deque (Optimal) | O(n) | O(k) | The intended answer — linear time |

---

## Approach 1 — Brute Force

### Intuition

There are `n - k + 1` windows. For each, scan its `k` elements and keep the largest. Correct, but overlapping windows re-examine the same elements over and over.

### Algorithm

1. For each start `i` from `0` to `n-k`:
2. Scan `nums[i..i+k-1]`, tracking the max.
3. Append it to the result.

### Complexity

- **Time:** O(n·k) — every window rescans `k` elements.
- **Space:** O(1) extra (output excluded).

### Code

```go
func bruteForce(nums []int, k int) []int {
	n := len(nums)
	if n == 0 || k == 0 {
		return []int{}
	}
	result := make([]int, 0, n-k+1)
	for i := 0; i+k <= n; i++ {
		maxVal := nums[i] // start with the window's first element
		for j := i + 1; j < i+k; j++ {
			if nums[j] > maxVal {
				maxVal = nums[j] // track the running window maximum
			}
		}
		result = append(result, maxVal)
	}
	return result
}
```

### Dry Run

Example 1: `nums = [1,3,-1,-3,5,3,6,7]`, `k = 3`.

| Window start i | window | max |
|----------------|--------|-----|
| 0 | [1, 3, -1] | 3 |
| 1 | [3, -1, -3] | 3 |
| 2 | [-1, -3, 5] | 5 |
| 3 | [-3, 5, 3] | 5 |
| 4 | [5, 3, 6] | 6 |
| 5 | [3, 6, 7] | 7 |

Result: `[3, 3, 5, 5, 6, 7]` ✔

---

## Approach 2 — Max-Heap (indices)

### Intuition

The window maximum is the largest value among the current indices. A max-heap surfaces the largest instantly. The only complication is that the heap's top might refer to an index that has already slid out of the window — so we **lazily delete** those stale tops when we peek.

### Algorithm

1. Push each `(value, index)` pair into a max-heap ordered by value.
2. Once `i >= k-1` (first full window), pop while the top's index `<= i-k` (out of window), then read the top value as this window's answer.

### Complexity

- **Time:** O(n log n) — up to `n` pushes and pops.
- **Space:** O(n) — the heap may hold every element.

### Code

```go
func maxHeap(nums []int, k int) []int {
	n := len(nums)
	if n == 0 || k == 0 {
		return []int{}
	}
	h := &idxHeap{}
	result := make([]int, 0, n-k+1)
	for i := 0; i < n; i++ {
		heap.Push(h, [2]int{nums[i], i}) // store value and its index
		if i >= k-1 {
			// Discard maxima whose index has left the window [i-k+1, i].
			for (*h)[0][1] <= i-k {
				heap.Pop(h)
			}
			result = append(result, (*h)[0][0]) // top value is the window max
		}
	}
	return result
}
```

*(`idxHeap` implements `container/heap` as a max-heap of `[value, index]` ordered by value.)*

### Dry Run

Example 1: `nums = [1,3,-1,-3,5,3,6,7]`, `k = 3`. Heap shown as sorted values with indices.

| i | pushed | heap top after evicting stale (idx ≤ i-k) | recorded |
|---|--------|-------------------------------------------|----------|
| 2 | (-1,2) | (3,1) | 3 |
| 3 | (-3,3) | (3,1) | 3 |
| 4 | (5,4) | (5,4) | 5 |
| 5 | (3,5) | (5,4) | 5 |
| 6 | (6,6) | (6,6) | 6 |
| 7 | (7,7) | (7,7) | 7 |

Result: `[3, 3, 5, 5, 6, 7]` ✔

---

## Approach 3 — Monotonic Deque (Optimal)

### Intuition

Key observation: if a smaller element appears **before** a larger one, the smaller one can never be the window max while the larger is still in range — so throw it away. We keep a deque of **indices** whose corresponding values are strictly decreasing. Then the front index always holds the current window's maximum. Because every index is pushed once and popped once, the whole thing is O(n).

### Algorithm

1. For each `i`: pop from the **back** while `nums[back] <= nums[i]` (those are dominated by `nums[i]`).
2. Push `i` at the back.
3. Pop from the **front** if it has slid out of the window (`front <= i-k`).
4. Once `i >= k-1`, record `nums[front]` — the current window maximum.

### Complexity

- **Time:** O(n) — each index enters and leaves the deque at most once.
- **Space:** O(k) — the deque holds at most one window of indices.

### Code

```go
func monotonicDeque(nums []int, k int) []int {
	n := len(nums)
	if n == 0 || k == 0 {
		return []int{}
	}
	deque := make([]int, 0, k) // holds indices, values strictly decreasing
	result := make([]int, 0, n-k+1)
	for i := 0; i < n; i++ {
		for len(deque) > 0 && nums[deque[len(deque)-1]] <= nums[i] {
			deque = deque[:len(deque)-1]
		}
		deque = append(deque, i) // nums[i] is a candidate max
		if deque[0] <= i-k {
			deque = deque[1:]
		}
		if i >= k-1 {
			result = append(result, nums[deque[0]]) // front = current window max
		}
	}
	return result
}
```

### Dry Run

Example 1: `nums = [1,3,-1,-3,5,3,6,7]`, `k = 3`. Deque shown as indices (values in parens).

| i | nums[i] | pop back while ≤ nums[i] | push i | evict front if ≤ i-k | deque (idx:val) | record |
|---|---------|--------------------------|--------|----------------------|-----------------|--------|
| 0 | 1 | — | [0] | — | 0:1 | — |
| 1 | 3 | pop 0 (1≤3) | [1] | — | 1:3 | — |
| 2 | -1 | — | [1,2] | — | 1:3, 2:-1 | **3** |
| 3 | -3 | — | [1,2,3] | evict 1 (≤0) | 2:-1, 3:-3 | **3** |
| 4 | 5 | pop 3,2 (≤5) | [4] | — | 4:5 | **5** |
| 5 | 3 | — | [4,5] | — | 4:5, 5:3 | **5** |
| 6 | 6 | pop 5,4 (≤6) | [6] | — | 6:6 | **6** |
| 7 | 7 | pop 6 (≤7) | [7] | — | 7:7 | **7** |

Result: `[3, 3, 5, 5, 6, 7]` ✔

---

## Key Takeaways

- **Monotonic deque = O(1) amortized window extremum.** Keep indices with decreasing values; the front is the max. This is the go-to pattern for "max/min of every fixed window."
- **Store indices, not values**, so you can tell when the front has aged out of the window (`front <= i-k`).
- **Lazy deletion** makes the heap version simple: don't pay to remove stale entries eagerly — skip them when they surface at the top.
- Popping from the back "because a bigger, later element dominates you" is the same domination logic behind Next Greater Element and Daily Temperatures.

---

## Related Problems

- LeetCode #239 variants — Sliding Window Minimum (flip the comparison)
- LeetCode #480 — Sliding Window Median (two-heap / multiset)
- LeetCode #862 — Shortest Subarray with Sum at Least K (monotonic deque on prefix sums)
- LeetCode #739 — Daily Temperatures (monotonic stack domination)
- LeetCode #496 — Next Greater Element I
