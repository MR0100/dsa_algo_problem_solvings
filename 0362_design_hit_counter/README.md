# 0362 — Design Hit Counter

> LeetCode #362 · Difficulty: Medium
> **Categories:** Design, Queue, Data Stream, Binary Search

---

## Problem Statement

Design a hit counter which counts the number of hits received in the past `5` minutes (i.e., the past `300` seconds).

Your system should accept a `timestamp` parameter (**in seconds** granularity), and you may assume that calls are being made to the system in chronological order (i.e., `timestamp` is monotonically non-decreasing). Several hits may arrive roughly at the same time.

Implement the `HitCounter` class:

- `HitCounter()` Initializes the object of the hit counter system.
- `void hit(int timestamp)` Records a hit that happened at `timestamp` (**in seconds**). Several hits may happen at the same `timestamp`.
- `int getHits(int timestamp)` Returns the number of hits in the past 5 minutes from `timestamp` (i.e., the past `300` seconds).

**Example 1:**

```
Input
["HitCounter", "hit", "hit", "hit", "getHits", "hit", "getHits", "getHits"]
[[], [1], [2], [3], [4], [300], [300], [301]]
Output
[null, null, null, null, 3, null, 4, 3]

Explanation
HitCounter hitCounter = new HitCounter();
hitCounter.hit(1);       // hit at timestamp 1.
hitCounter.hit(2);       // hit at timestamp 2.
hitCounter.hit(3);       // hit at timestamp 3.
hitCounter.getHits(4);   // get hits at timestamp 4, return 3.
hitCounter.hit(300);     // hit at timestamp 300.
hitCounter.getHits(300); // get hits at timestamp 300, return 4.
hitCounter.getHits(301); // get hits at timestamp 301, return 3.
```

**Constraints:**

- `1 <= timestamp <= 2 * 10^9`
- All the calls are being made to the system in chronological order (i.e., `timestamp` is monotonically increasing).
- At most `300` calls will be made to `hit` and `getHits`.

**Follow up:** What if the number of hits per second could be huge? Does your design scale?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★★☆ High       | 2023          |
| Dropbox    | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Design a Data Structure** — the deliverable is a class with `hit`/`getHits` methods and internal state → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)
- **Queue / FIFO** — timestamps expire in arrival order, the classic use for a queue that pops from the front → see [`/dsa/queue_deque.md`](/dsa/queue_deque.md)
- **Fixed-size Circular Buffer (bucketing by time mod 300)** — the O(1)-space design maps second `t` to slot `t % 300`; closest existing reference is the design-structures file → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)

---

## Approaches Overview

| # | Approach | Hit | GetHits | Space | When to use |
|---|----------|-----|---------|-------|-------------|
| 1 | Queue of Timestamps | O(1) amortised | O(k) amortised | O(n) | Simple; fine when total hits are modest |
| 2 | Fixed 300-Slot Circular Buffer (Optimal) | O(1) | O(300)=O(1) | O(300)=O(1) | Scales to huge hits/sec — answers the follow-up |

---

## Approach 1 — Queue of Timestamps

### Intuition

A hit is "alive" for exactly 300 seconds. Store every hit timestamp in a FIFO queue in arrival order. Because timestamps are non-decreasing, the oldest hits sit at the front. On `getHits(t)`, pop every front timestamp that has expired (`≤ t-300`); whatever remains is precisely the hits inside the window, so the count is the queue length.

### Algorithm

1. `hit(t)`: append `t` to the back of the queue.
2. `getHits(t)`: while the front timestamp is `≤ t-300`, pop it (it expired).
3. Return the queue's length.

### Complexity

- **Time:** `hit` O(1); `getHits` O(k) where k is the number of newly-expired entries. Each timestamp is enqueued once and dequeued once, so it is O(1) amortised per hit overall.
- **Space:** O(n) — one entry per stored hit; a burst of hits in one second all occupy separate slots.

### Code

```go
type HitCounterQueue struct {
	q []int // timestamps in non-decreasing (arrival) order
}

func NewHitCounterQueue() *HitCounterQueue {
	return &HitCounterQueue{q: []int{}}
}

func (h *HitCounterQueue) Hit(timestamp int) {
	h.q = append(h.q, timestamp) // append keeps the queue sorted by time
}

func (h *HitCounterQueue) GetHits(timestamp int) int {
	// Evict from the front everything that fell out of the 5-minute window.
	for len(h.q) > 0 && h.q[0] <= timestamp-300 {
		h.q = h.q[1:] // pop the oldest expired timestamp
	}
	return len(h.q) // survivors are exactly the hits within the window
}
```

### Dry Run

Operations: `hit(1) hit(2) hit(3) getHits(4) hit(300) getHits(300) getHits(301)`.

| Call | queue before | expiry cutoff (t-300) | pops | queue after | return |
|------|--------------|-----------------------|------|-------------|--------|
| hit(1)       | []          | —    | —          | [1]         | — |
| hit(2)       | [1]         | —    | —          | [1,2]       | — |
| hit(3)       | [1,2]       | —    | —          | [1,2,3]     | — |
| getHits(4)   | [1,2,3]     | -296 | none (all > -296) | [1,2,3] | **3** |
| hit(300)     | [1,2,3]     | —    | —          | [1,2,3,300] | — |
| getHits(300) | [1,2,3,300] | 0    | none (all > 0)    | [1,2,3,300] | **4** |
| getHits(301) | [1,2,3,300] | 1    | pop 1 (1 ≤ 1)     | [2,3,300]   | **3** |

Answers `3, 4, 3` ✔

---

## Approach 2 — Fixed 300-Slot Circular Buffer (Optimal)

### Intuition

The window is always exactly 300 seconds, so seconds `t` and `t+300` collide on the same bucket `t % 300`. Give each bucket two fields: the **timestamp** it currently represents and the **hit count** for that second. On a hit, if the bucket's stored timestamp is stale (a different second from a previous 300-cycle), overwrite it and reset the count; otherwise increment. On a query, sum only buckets whose stored timestamp is still within the window. Memory is a fixed 300 slots no matter how many hits arrive — this is the answer to the follow-up.

### Algorithm

1. Keep `times[300]` and `counts[300]`.
2. `hit(t)`: let `i = t % 300`. If `times[i] != t`, set `times[i]=t, counts[i]=1` (stale slot, reset); else `counts[i]++`.
3. `getHits(t)`: sum `counts[i]` over all `i` where `t - times[i] < 300`.

### Complexity

- **Time:** `hit` O(1) (one bucket); `getHits` O(300) = O(1) (fixed scan).
- **Space:** O(300) = O(1) — two constant-size arrays independent of hit volume. Bursts of thousands of hits in one second collapse into a single incremented counter.

### Code

```go
type HitCounterBuckets struct {
	times  [300]int // times[i] = the timestamp that counts[i] refers to
	counts [300]int // counts[i] = number of hits during second times[i]
}

func NewHitCounterBuckets() *HitCounterBuckets {
	return &HitCounterBuckets{}
}

func (h *HitCounterBuckets) Hit(timestamp int) {
	i := timestamp % 300 // which of the 300 slots this second maps to
	if h.times[i] != timestamp {
		// This slot last held a DIFFERENT second (some multiple of 300 ago).
		// Overwrite it: rebind the slot to `timestamp` and start counting fresh.
		h.times[i] = timestamp
		h.counts[i] = 1
	} else {
		// Same second as the slot already tracks — just increment.
		h.counts[i]++
	}
}

func (h *HitCounterBuckets) GetHits(timestamp int) int {
	total := 0
	for i := 0; i < 300; i++ {
		// A slot contributes only if the second it holds is still within the
		// window (strictly greater than timestamp-300).
		if timestamp-h.times[i] < 300 {
			total += h.counts[i]
		}
	}
	return total
}
```

### Dry Run

Same operations. Slots shown only for the ones touched (`1→slot 1`, `2→slot 2`, `3→slot 3`, `300→slot 0`).

| Call | slot i | times[i] before | action | times[i]/counts[i] after | return |
|------|--------|-----------------|--------|--------------------------|--------|
| hit(1)   | 1 | 0 | ≠1 → reset | times[1]=1, counts[1]=1 | — |
| hit(2)   | 2 | 0 | ≠2 → reset | times[2]=2, counts[2]=1 | — |
| hit(3)   | 3 | 0 | ≠3 → reset | times[3]=3, counts[3]=1 | — |
| getHits(4)   | — | — | sum slots with `4-times<300` → slots 1,2,3 | — | **3** |
| hit(300) | 0 | 0 | ≠300 → reset | times[0]=300, counts[0]=1 | — |
| getHits(300) | — | — | `300-times<300`: slot0(300→0<300 ✔), slots1,2,3 (299/298/297<300 ✔) | — | **4** |
| getHits(301) | — | — | slot1: `301-1=300` **not <300** ✗; slots2,3 ✔; slot0 `301-300=1` ✔ | — | **3** |

Answers `3, 4, 3` ✔ — at `getHits(301)` the hit at second 1 is exactly 300 seconds old and correctly excluded.

---

## Key Takeaways

- **Sliding time window → drop expired entries.** A FIFO queue models "expire in arrival order" perfectly; pop the front lazily on read.
- **When the window is a fixed constant, bucket by `time % window`.** This bounds memory to O(window) and both operations to O(1), independent of throughput — the standard scaling answer.
- **Store timestamp + count per bucket** so you can distinguish a fresh hit from a stale value left over from a previous cycle; comparing the stored timestamp is what makes overwrite-vs-increment correct.
- The follow-up (huge hits/sec) is exactly why the bucket design wins: a million hits in one second is one counter increment, not a million queue entries.

---

## Related Problems

- LeetCode #359 — Logger Rate Limiter (time-window design)
- LeetCode #346 — Moving Average from Data Stream (fixed-window queue)
- LeetCode #933 — Number of Recent Calls (queue of timestamps, 3000ms window)
- LeetCode #1352 — Product of the Last K Numbers (streaming window design)
