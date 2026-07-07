# 0264 — Ugly Number II

> LeetCode #264 · Difficulty: Medium
> **Categories:** Hash Table, Math, Dynamic Programming, Heap (Priority Queue)

---

## Problem Statement

An **ugly number** is a positive integer whose prime factors are limited to `2`, `3`, and `5`.

Given an integer `n`, return the `nth` ugly number.

**Example 1:**

```
Input: n = 10
Output: 12
Explanation: [1, 2, 3, 4, 5, 6, 8, 9, 10, 12] is the sequence of the first 10 ugly numbers.
```

**Example 2:**

```
Input: n = 1
Output: 1
Explanation: 1 has no prime factors, therefore all of its prime factors are limited to 2, 3, and 5.
```

**Constraints:**

- `1 <= n <= 1690`

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Amazon    | ★★★★☆ High       | 2024          |
| Microsoft | ★★★☆☆ Medium     | 2023          |
| Google    | ★★★☆☆ Medium     | 2023          |
| Bloomberg | ★★★☆☆ Medium     | 2023          |
| Meta      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Min-Heap (Priority Queue)** — repeatedly extract the smallest ugly value and expand it → see [`/dsa/heap_priority_queue.md`](/dsa/heap_priority_queue.md)
- **Dynamic Programming (multi-pointer merge)** — build the sorted sequence with three factor pointers → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)
- **Hash Set de-duplication** — the heap approach needs a seen-set to avoid duplicate multiples → see [`/dsa/hash_map.md`](/dsa/hash_map.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Min-Heap | O(n log n) | O(n) | Intuitive; generalises to arbitrary prime sets |
| 2 | DP Three Pointers (Optimal) | O(n) | O(n) | Fastest for a fixed prime set {2,3,5} |

---

## Approach 1 — Min-Heap

### Intuition
Every ugly number multiplied by 2, 3, or 5 is again ugly. Starting from 1, repeatedly pop the smallest value seen so far and push its three multiples. The value popped on the `n`th iteration is the `n`th ugly number. A `seen` set stops duplicates (e.g. `6 = 2·3 = 3·2`) from being enqueued twice.

### Algorithm
1. Push `1` into a min-heap and mark it seen.
2. Repeat `n` times: pop the smallest value (the next ugly number); for each `f` in `{2,3,5}`, push `f·value` if not seen.
3. The value popped on iteration `n` is the answer.

### Complexity
- **Time:** O(n log n) — `n` pops, each pushing up to 3 values at O(log n).
- **Space:** O(n) — heap and seen-set each hold O(n) entries.

### Code
```go
func minHeap(n int) int {
	h := &intHeap{1}              // start the sequence at 1
	seen := map[int]bool{1: true} // avoid pushing the same value twice
	var val int
	for i := 0; i < n; i++ {
		val = heap.Pop(h).(int) // the (i+1)-th smallest ugly number
		for _, f := range []int{2, 3, 5} {
			next := val * f
			if !seen[next] { // only queue unseen multiples
				seen[next] = true
				heap.Push(h, next)
			}
		}
	}
	return val // last popped value is the nth ugly number
}
```

### Dry Run
`n = 10`. Each row pops `val` then pushes unseen `2v,3v,5v`.

| iter | pop val | heap top values pushed | heap after (sorted) |
|------|---------|------------------------|---------------------|
| 1    | 1       | 2,3,5                  | 2,3,5               |
| 2    | 2       | 4,6,10                 | 3,4,5,6,10          |
| 3    | 3       | 6(seen),9,15           | 4,5,6,9,10,15       |
| 4    | 4       | 8,12,20                | 5,6,8,9,10,12,15,20 |
| 5    | 5       | 10(seen),15(seen),25   | 6,8,9,10,12,15,20,25|
| 6    | 6       | 12(seen),18,30         | 8,9,10,12,15,18,...  |
| 7    | 8       | 16,24,40               | 9,10,12,15,...       |
| 8    | 9       | 18(seen),27,45         | 10,12,15,...         |
| 9    | 10      | 20(seen),30(seen),50   | 12,15,...            |
| 10   | **12**  | —                      | —                    |

10th pop is `12` → answer `12`.

---

## Approach 2 — DP Three Pointers (Optimal)

### Intuition
The next ugly number is always an earlier ugly number times 2, 3, or 5. Keep three indices `i2, i3, i5` into the growing sequence — each pointing at the earliest term not yet multiplied by that factor. The next term is `min(ugly[i2]·2, ugly[i3]·3, ugly[i5]·5)`; advance **every** pointer whose product equals that min, which also dedupes ties.

### Algorithm
1. `ugly[0] = 1`; pointers `i2 = i3 = i5 = 0`.
2. For `k = 1..n-1`: `a = ugly[i2]·2`, `b = ugly[i3]·3`, `c = ugly[i5]·5`; `ugly[k] = min(a,b,c)`; advance each pointer whose candidate equals `ugly[k]`.
3. Return `ugly[n-1]`.

### Complexity
- **Time:** O(n) — one linear fill, each term O(1).
- **Space:** O(n) — the `ugly` array of size `n`.

### Code
```go
func dpThreePointers(n int) int {
	ugly := make([]int, n) // ugly[k] = (k+1)-th ugly number, in sorted order
	ugly[0] = 1            // the 1st ugly number is 1
	i2, i3, i5 := 0, 0, 0  // next index to multiply by 2, 3, 5 respectively
	for k := 1; k < n; k++ {
		a := ugly[i2] * 2 // smallest unused multiple of 2
		b := ugly[i3] * 3 // smallest unused multiple of 3
		c := ugly[i5] * 5 // smallest unused multiple of 5
		next := a
		if b < next {
			next = b
		}
		if c < next {
			next = c
		}
		ugly[k] = next // the next ugly number is the smallest candidate
		// advance ALL pointers that produced this value (handles duplicates)
		if next == a {
			i2++
		}
		if next == b {
			i3++
		}
		if next == c {
			i5++
		}
	}
	return ugly[n-1] // nth ugly number (0-indexed n-1)
}
```

### Dry Run
`n = 10`. Start `ugly=[1]`, `i2=i3=i5=0`.

| k | a=u[i2]·2 | b=u[i3]·3 | c=u[i5]·5 | min | ugly[k] | pointers advanced | i2,i3,i5 |
|---|-----------|-----------|-----------|-----|---------|-------------------|----------|
| 1 | 1·2=2     | 1·3=3     | 1·5=5     | 2   | 2       | i2                | 1,0,0    |
| 2 | 2·2=4     | 1·3=3     | 1·5=5     | 3   | 3       | i3                | 1,1,0    |
| 3 | 2·2=4     | 2·3=6     | 1·5=5     | 4   | 4       | i2                | 2,1,0    |
| 4 | 3·2=6     | 2·3=6     | 1·5=5     | 5   | 5       | i5                | 2,1,1    |
| 5 | 3·2=6     | 2·3=6     | 2·5=10    | 6   | 6       | i2,i3             | 3,2,1    |
| 6 | 4·2=8     | 3·3=9     | 2·5=10    | 8   | 8       | i2                | 4,2,1    |
| 7 | 5·2=10    | 3·3=9     | 2·5=10    | 9   | 9       | i3                | 4,3,1    |
| 8 | 5·2=10    | 4·3=12    | 2·5=10    | 10  | 10      | i2,i5             | 5,3,2    |
| 9 | 6·2=12    | 4·3=12    | 3·5=15    | 12  | **12**  | i2,i3             | 6,4,2    |

`ugly[9] = 12` → answer `12`.

---

## Key Takeaways

- **Generate in sorted order, don't test every integer.** Building the sequence from prior terms is far cheaper than checking `isUgly` on 1,2,3,… up to the answer.
- **Advance all tied pointers** in the DP to skip duplicates without a separate set — this is the crucial trick over the heap's explicit `seen` map.
- The three-pointer DP is the merge of three sorted streams `{2·ugly}`, `{3·ugly}`, `{5·ugly}` — the same k-way-merge pattern used in Super Ugly Number (#313) with more primes/pointers.

---

## Related Problems

- LeetCode #263 — Ugly Number (predicate check for a single number)
- LeetCode #313 — Super Ugly Number (arbitrary prime list → k pointers / heap)
- LeetCode #23 — Merge k Sorted Lists (same heap / k-way merge structure)
- LeetCode #1201 — Ugly Number III (binary search + inclusion–exclusion)
