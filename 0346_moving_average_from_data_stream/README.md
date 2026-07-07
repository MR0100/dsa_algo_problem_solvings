# 0346 — Moving Average from Data Stream

> LeetCode #346 · Difficulty: Easy
> **Categories:** Design, Queue, Array, Data Stream, Sliding Window

---

## Problem Statement

Given a stream of integers and a window size, calculate the moving average of all integers in the sliding window.

Implement the `MovingAverage` class:

- `MovingAverage(int size)` Initializes the object with the size of the window `size`.
- `double next(int val)` Returns the moving average of the last `size` values of the stream.

**Example 1:**
```
Input
["MovingAverage", "next", "next", "next", "next"]
[[3], [1], [10], [3], [5]]
Output
[null, 1.0, 5.5, 4.666666666666667, 6.0]

Explanation
MovingAverage movingAverage = new MovingAverage(3);
movingAverage.next(1); // return 1.0 = 1 / 1
movingAverage.next(10); // return 5.5 = (1 + 10) / 2
movingAverage.next(3); // return 4.66667 = (1 + 10 + 3) / 3
movingAverage.next(5); // return 6.0 = (10 + 3 + 5) / 3
```

**Constraints:**
- `1 <= size <= 1000`
- `-10⁵ <= val <= 10⁵`
- At most `10⁴` calls will be made to `next`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2024          |
| Facebook   | ★★☆☆☆ Low        | 2023          |
| Uber       | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used
- **Queue / Circular Buffer** — the window is a FIFO structure: newest in at the back, oldest out at the front → see [`/dsa/queue_deque.md`](/dsa/queue_deque.md)
- **Sliding Window** — a fixed-length window slides across the data stream one element at a time → see [`/dsa/sliding_window.md`](/dsa/sliding_window.md)
- **Design Data Structures** — implementing a stateful class with incremental updates → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)

---

## Approaches Overview
| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Queue with Full Recompute (Brute Force) | O(size) per next | O(size) | Simplest to reason about; fine for tiny windows |
| 2 | Circular Buffer + Running Sum (Optimal) | O(1) per next | O(size) | The standard answer; constant-time updates |

---

## Approach 1 — Queue with Full Recompute (Brute Force)

### Intuition
The moving average of the last `size` values is just the average of a sliding window. Store the window explicitly as a slice: push the newcomer, drop the oldest if we overflow capacity, then sum whatever remains and divide by the count. It is correct by construction — no invariants to maintain beyond "keep at most `size` values".

### Algorithm
1. Append `val` to the window slice (newest at the back).
2. If `len(window) > size`, slice off the front element (the oldest value).
3. Sum every element currently in the window.
4. Return `sum / len(window)` as a `float64`.

### Complexity
- **Time:** O(size) per `next` — step 3 re-walks the entire window each call.
- **Space:** O(size) — the window slice holds at most `size` integers.

### Code
```go
type MovingAverageBrute struct {
	size   int   // maximum number of values the window may hold
	window []int // the values currently inside the window, oldest at front
}

func NewMovingAverageBrute(size int) *MovingAverageBrute {
	return &MovingAverageBrute{size: size}
}

func (m *MovingAverageBrute) Next(val int) float64 {
	m.window = append(m.window, val) // push newest value to the back
	if len(m.window) > m.size {      // window overflowed its capacity?
		m.window = m.window[1:] // evict the oldest value at the front
	}
	sum := 0
	for _, v := range m.window { // re-add everything currently in the window
		sum += v
	}
	return float64(sum) / float64(len(m.window))
}
```

### Dry Run
Trace Example 1 with `size = 3`:

| next(val) | window after push/evict | sum | count | average |
|-----------|-------------------------|-----|-------|---------|
| next(1)   | [1]                     | 1   | 1     | 1.0     |
| next(10)  | [1, 10]                 | 11  | 2     | 5.5     |
| next(3)   | [1, 10, 3]              | 14  | 3     | 4.666666666666667 |
| next(5)   | [10, 3, 5] (1 evicted)  | 18  | 3     | 6.0     |

---

## Approach 2 — Circular Buffer + Running Sum (Optimal)

### Intuition
We do not need to re-sum the window every call. When a new value slides in, exactly one old value slides out — the one being overwritten. Keep a **running sum**: add the incoming value and subtract the outgoing one. A fixed ring buffer of length `size` gives O(1) access to the slot being overwritten and never grows. While the window is still filling, the overwritten slot holds `0`, so the subtraction is harmless.

### Algorithm
1. Compute `head = count % size` — the slot the new value occupies (wraps around).
2. Update `sum += val − buf[head]` — add the newcomer, remove whatever it evicts.
3. Store `buf[head] = val` and increment `count`.
4. Let `live = min(count, size)` be the number of real values in the window.
5. Return `sum / live`.

### Complexity
- **Time:** O(1) per `next` — constant index arithmetic and one add/subtract.
- **Space:** O(size) — a single ring buffer allocated once at construction.

### Code
```go
type MovingAverageCircular struct {
	size  int   // window capacity and ring length
	buf   []int // ring buffer of the last `size` values
	count int   // total number of Next() calls so far
	sum   int   // running sum of the values currently in the window
}

func NewMovingAverageCircular(size int) *MovingAverageCircular {
	return &MovingAverageCircular{size: size, buf: make([]int, size)}
}

func (m *MovingAverageCircular) Next(val int) float64 {
	head := m.count % m.size   // slot to (over)write; wraps around the ring
	m.sum += val - m.buf[head] // add newcomer, remove whatever it evicts
	m.buf[head] = val          // store the new value in its slot
	m.count++                  // one more value has entered the stream
	live := m.count            // how many real values are in the window
	if live > m.size {         // once past capacity, only `size` are live
		live = m.size
	}
	return float64(m.sum) / float64(live)
}
```

### Dry Run
Trace Example 1 with `size = 3`, `buf = [0,0,0]`:

| next(val) | head = count%3 | buf[head] before | sum += val−old | buf after      | count | live | average |
|-----------|----------------|------------------|----------------|----------------|-------|------|---------|
| next(1)   | 0              | 0                | 0 + (1−0)=1    | [1,0,0]        | 1     | 1    | 1.0     |
| next(10)  | 1              | 0                | 1 + (10−0)=11  | [1,10,0]       | 2     | 2    | 5.5     |
| next(3)   | 2              | 0                | 11 + (3−0)=14  | [1,10,3]       | 3     | 3    | 4.666666666666667 |
| next(5)   | 0              | 1                | 14 + (5−1)=18  | [5,10,3]       | 4     | 3    | 6.0     |

Note how `next(5)` overwrites slot 0 (which held the oldest value `1`), and the running sum subtracts that `1` automatically.

---

## Key Takeaways
- **Fixed-window streaming → running sum.** Whenever a window of constant size slides one step, maintain an aggregate incrementally instead of recomputing: add the entering element, subtract the leaving one.
- **Ring buffer index trick:** `count % size` gives the slot to overwrite, and that slot always holds exactly the value about to leave the window — no separate bookkeeping of "the oldest" is needed.
- **Guard the divisor with `min(count, size)`** so the average is correct while the window is still filling up.
- Both approaches share the same O(size) space; the win is turning O(size) per query into O(1).

---

## Related Problems
- LeetCode #239 — Sliding Window Maximum (fixed window, monotonic deque)
- LeetCode #480 — Sliding Window Median (fixed window, two heaps)
- LeetCode #362 — Design Hit Counter (time-windowed stream counting)
- LeetCode #1352 — Product of the Last K Numbers (streaming aggregate over a window)
