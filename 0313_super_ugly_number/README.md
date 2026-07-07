# 0313 — Super Ugly Number

> LeetCode #313 · Difficulty: Medium
> **Categories:** Array, Math, Dynamic Programming, Heap (Priority Queue)

---

## Problem Statement

A **super ugly number** is a positive integer whose prime factors are in the array `primes`.

Given an integer `n` and an array of integers `primes`, return the `n`th super ugly number.

The `n`th super ugly number is **guaranteed** to fit in a **32-bit** signed integer.

**Example 1:**

```
Input: n = 12, primes = [2,7,13,19]
Output: 32
Explanation: [1,2,4,7,8,13,14,16,19,26,28,32] is the sequence of the first 12
super ugly numbers given primes = [2,7,13,19].
```

**Example 2:**

```
Input: n = 1, primes = [2,3,5]
Output: 1
Explanation: 1 has no prime factors, therefore all of its prime factors are in
the array primes = [2,3,5].
```

**Constraints:**

- `1 <= n <= 10^5`
- `1 <= primes.length <= 100`
- `2 <= primes[i] <= 1000`
- `primes[i]` is guaranteed to be a prime number.
- All the values of `primes` are **unique** and sorted in **ascending order**.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Heap / Priority Queue** — always extracting the globally smallest next candidate → see [`/dsa/heap_priority_queue.md`](/dsa/heap_priority_queue.md)
- **Dynamic Programming (1-D)** — building the sorted sequence incrementally, each new number from earlier ones → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)
- **Multi-pointer merge of sorted streams** — merging `k` sorted sequences `prime*ugly[...]` (k-way merge, same idea as merging sorted lists) → see [`/dsa/heap_priority_queue.md`](/dsa/heap_priority_queue.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Min-Heap of candidates | O(n·k·log(n·k)) | O(n·k) | Intuitive; generalizes to streaming |
| 2 | DP with k Pointers (Optimal) | O(n·k) | O(n + k) | Best time/space; standard answer |

---

## Approach 1 — Min-Heap

### Intuition
Every super ugly number besides `1` equals `prime * (an earlier super ugly number)`. Seed a min-heap with each prime (`prime*1`). Each pop of the global minimum yields the next super ugly number; push `prime*value` for every prime to generate its successors. A `seen` set avoids emitting the same value twice (e.g. `2*7` and `7*2` collide).

### Algorithm
1. `ugly = 1` is the 1st super ugly number; push each prime into the heap.
2. Repeat `n-1` times: pop the smallest unseen value → next super ugly number; for each prime push `value*prime` (if unseen).
3. Return the `n`th popped value.

### Complexity
- **Time:** O(n·k·log(n·k)) — up to `n·k` heap pushes, each `log`-time.
- **Space:** O(n·k) — heap plus `seen` set.

### Code
```go
func minHeap(n int, primes []int) int {
	pq := &int64Heap{}
	heap.Init(pq)
	seen := map[int64]bool{1: true}
	for _, p := range primes {
		heap.Push(pq, int64(p))
		seen[int64(p)] = true
	}

	ugly := int64(1)
	for i := 1; i < n; i++ {
		ugly = heap.Pop(pq).(int64)
		for _, p := range primes {
			cand := ugly * int64(p)
			if !seen[cand] {
				seen[cand] = true
				heap.Push(pq, cand)
			}
		}
	}
	return int(ugly)
}
```

### Dry Run
`n = 12`, `primes = [2,7,13,19]`. Heap seeded with `{2,7,13,19}`, `ugly=1` (1st).

| i | pop (i+1-th ugly) | push (unseen products) |
|---|-------------------|------------------------|
| 1 | 2 | 4, 14, 26, 38 |
| 2 | 4 | 8, 28, 52, 76 |
| 3 | 7 | (14 seen), 49, 91, 133 |
| 4 | 8 | 16, 56, 104, 152 |
| 5 | 13 | (26 seen), 91(seen), 169, 247 |
| 6 | 14 | 28(seen), 98, 182, 266 |
| 7 | 16 | 32, 112, 208, 304 |
| 8 | 19 | 38(seen), 133(seen), 247(seen), 361 |
| 9 | 26 | 52(seen), 182(seen), 338, 494 |
| 10 | 28 | 56(seen), 196, 364, 532 |
| 11 | 32 | ... |

The 12th popped value (`i=11`) is **32**.

---

## Approach 2 — DP with k Pointers (Optimal)

### Intuition
The full sorted list of super ugly numbers is the **merge of `k` sorted sequences**: for prime `p`, the sequence `p*ugly[0] < p*ugly[1] < ...`. Keep a pointer `idx[j]` into `ugly` for each prime. The next super ugly number is the minimum of `primes[j]*ugly[idx[j]]` over all `j`. Advance **every** pointer that achieves that minimum — this dedupes collisions naturally.

### Algorithm
1. `ugly[0] = 1`; `idx[j] = 0` for all primes.
2. For `i` from 1 to `n-1`: `next = min over j of primes[j]*ugly[idx[j]]`.
3. `ugly[i] = next`; for every `j` with `primes[j]*ugly[idx[j]] == next`, `idx[j]++`.
4. Return `ugly[n-1]`.

### Complexity
- **Time:** O(n·k) — for each of `n` numbers, scan `k` primes twice.
- **Space:** O(n + k) — the `ugly` array and the pointer array.

### Code
```go
func dpPointers(n int, primes []int) int {
	k := len(primes)
	ugly := make([]int, n)
	ugly[0] = 1
	idx := make([]int, k)

	for i := 1; i < n; i++ {
		next := int(^uint(0) >> 1)
		for j := 0; j < k; j++ {
			if cand := primes[j] * ugly[idx[j]]; cand < next {
				next = cand
			}
		}
		ugly[i] = next
		for j := 0; j < k; j++ {
			if primes[j]*ugly[idx[j]] == next {
				idx[j]++
			}
		}
	}
	return ugly[n-1]
}
```

### Dry Run
`n = 12`, `primes = [2,7,13,19]`. `ugly=[1]`, `idx=[0,0,0,0]`.

| i | candidates (2·, 7·, 13·, 19·) | next=ugly[i] | pointers advanced → idx |
|---|-------------------------------|--------------|--------------------------|
| 1 | 2, 7, 13, 19 | 2 | j0 → [1,0,0,0] |
| 2 | 4, 7, 13, 19 | 4 | j0 → [2,0,0,0] |
| 3 | 8, 7, 13, 19 | 7 | j1 → [2,1,0,0] |
| 4 | 8, 14, 13, 19 | 8 | j0 → [3,1,0,0] |
| 5 | 14, 14, 13, 19 | 13 | j2 → [3,1,1,0] |
| 6 | 14, 14, 26, 19 | 14 | j0,j1 → [4,2,1,0] |
| 7 | 16, 28, 26, 19 | 16 | j0 → [5,2,1,0] |
| 8 | 16→ next uses ugly[5]=13·2=... 26 etc; min=19 | 19 | j3 → [5,2,1,1] |
| 9 | 2·13=26, 7·4=28, 13·2=26, 19·2=38 | 26 | j0,j2 → [6,2,2,1] |
| 10 | 2·14=28, 28, 13·4=52, 38 | 28 | j0,j1 → [7,3,2,1] |
| 11 | 2·16=32, 7·7=49, 52, 38 | 32 | j0 → [8,3,2,1] |

`ugly[11] = 32`.

---

## Key Takeaways
- **Super ugly = k-way merge of `prime * earlier-ugly` streams.** Recognizing this turns the problem into a merge, giving the O(n·k) pointer DP.
- **Advance all tied pointers** to dedupe: if two primes produce the same next value, moving both prevents duplicates without a `seen` set.
- The **heap** approach is more intuitive and streams, but pays a `log` factor and stores `O(n·k)` candidates; the pointer DP is strictly better here.
- Generalization of LeetCode #264 (Ugly Number II, fixed primes `{2,3,5}`) to an arbitrary prime set.

---

## Related Problems
- LeetCode #264 — Ugly Number II (this problem with primes `{2,3,5}`)
- LeetCode #263 — Ugly Number (factor test)
- LeetCode #23 — Merge k Sorted Lists (k-way merge with a heap)
- LeetCode #378 — Kth Smallest Element in a Sorted Matrix (heap / merge)
