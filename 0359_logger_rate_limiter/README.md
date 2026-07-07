# 0359 — Logger Rate Limiter

> LeetCode #359 · Difficulty: Easy
> **Categories:** Hash Table, Design, Data Stream, Queue

---

## Problem Statement

Design a logger system that receives a stream of messages along with their timestamps. Each **unique** message should only be printed **at most every 10 seconds** (i.e. a message printed at timestamp `t` will prevent other identical messages from being printed until timestamp `t + 10`).

All messages will come in chronological order. Several messages may arrive at the same timestamp.

Implement the `Logger` class:

- `Logger()` Initializes the `logger` object.
- `bool shouldPrintMessage(int timestamp, string message)` Returns `true` if the `message` should be printed in the given `timestamp`, otherwise returns `false`.

**Example 1:**

```
Input
["Logger", "shouldPrintMessage", "shouldPrintMessage", "shouldPrintMessage", "shouldPrintMessage", "shouldPrintMessage", "shouldPrintMessage"]
[[], [1, "foo"], [2, "bar"], [3, "foo"], [8, "bar"], [10, "foo"], [11, "foo"]]
Output
[null, true, true, false, false, false, true]

Explanation
Logger logger = new Logger();
logger.shouldPrintMessage(1, "foo");  // return true, next allowed timestamp for "foo" is 1 + 10 = 11
logger.shouldPrintMessage(2, "bar");  // return true, next allowed timestamp for "bar" is 2 + 10 = 12
logger.shouldPrintMessage(3, "foo");  // 3 < 11, return false
logger.shouldPrintMessage(8, "bar");  // 8 < 12, return false
logger.shouldPrintMessage(10, "foo"); // 10 < 11, return false
logger.shouldPrintMessage(11, "foo"); // 11 >= 11, return true, next allowed timestamp for "foo" is 11 + 10 = 21
```

**Constraints:**

- `0 <= timestamp <= 10^9`
- Every `timestamp` will be passed in non-decreasing order (chronological order).
- `1 <= message.length <= 30`
- At most `10^4` calls will be made to `shouldPrintMessage`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Facebook   | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Map** — map each message to the earliest timestamp it may print again, giving O(1) lookups → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Design (Data Structure API)** — implement a class with a stateful method over a data stream → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)
- **Queue (Sliding Window Eviction)** — a FIFO of recent prints models the literal 10-second window, evicting stale entries from the front → see [`/dsa/queue_deque.md`](/dsa/queue_deque.md)

---

## Approaches Overview

| # | Approach | Time / call | Space | When to use |
|---|----------|-------------|-------|-------------|
| 1 | Hash Map of next-allowed time (Optimal) | O(1) | O(m) messages | Simplest and fastest |
| 2 | Queue + Set (sliding window) | O(1) amortized | O(w) window | When memory must shrink to the active window |

---

## Approach 1 — Hash Map of Next-Allowed Time (Optimal)

### Intuition

The rule "not printed in the last 10 seconds" is equivalent to: once printed at time `t`, the message is blocked until `t + 10`. So store, per message, its **next-allowed time**. A request at `timestamp` prints iff `timestamp >= nextAllowed[message]`; on printing, set `nextAllowed[message] = timestamp + 10`. Messages never seen default to next-allowed `0`, so they always print the first time.

### Algorithm

1. Read `nextAllowed[message]` (Go's zero value `0` when absent).
2. If `timestamp < nextAllowed[message]`, return `false` (still cooling down).
3. Otherwise set `nextAllowed[message] = timestamp + 10` and return `true`.

### Complexity

- **Time:** O(1) average — one map read, one map write.
- **Space:** O(m) — one entry per distinct message ever printed (the map never shrinks).

### Code

```go
type HashMapLogger struct {
	nextAllowed map[string]int // message → earliest timestamp it may print again
}

func NewHashMapLogger() *HashMapLogger {
	return &HashMapLogger{nextAllowed: make(map[string]int)}
}

func (l *HashMapLogger) ShouldPrintMessage(timestamp int, message string) bool {
	// Absent messages return the zero value 0, so first-ever calls pass.
	if timestamp < l.nextAllowed[message] {
		return false // printed within the last 10 seconds — suppress
	}
	l.nextAllowed[message] = timestamp + 10 // block this message until +10s
	return true
}
```

### Dry Run

Example 1 call sequence:

| Call | timestamp | message | nextAllowed[msg] before | timestamp < before? | print? | nextAllowed after |
|------|-----------|---------|-------------------------|---------------------|--------|-------------------|
| 1 | 1 | foo | 0 | no | **true** | foo → 11 |
| 2 | 2 | bar | 0 | no | **true** | bar → 12 |
| 3 | 3 | foo | 11 | 3 < 11 yes | **false** | unchanged |
| 4 | 8 | bar | 12 | 8 < 12 yes | **false** | unchanged |
| 5 | 10 | foo | 11 | 10 < 11 yes | **false** | unchanged |
| 6 | 11 | foo | 11 | 11 < 11 no | **true** | foo → 21 |

Output: `[true, true, false, false, false, true]` ✔

---

## Approach 2 — Queue + Set (Sliding Window Eviction)

### Intuition

Model the "last 10 seconds" literally. Keep a FIFO queue of `(timestamp, message)` for everything printed in the current window plus a set of the messages inside it. On a new request at `timestamp`, first evict every front entry with `time <= timestamp - 10` (aged out), removing it from the set. Then the message prints iff it is **not** in the set. Because entries leave the window as the clock advances, memory tracks the active window rather than all messages ever seen.

### Algorithm

1. Pop front entries while `front.timestamp <= timestamp - 10`, deleting each from the set.
2. If `message` is in the set, return `false`.
3. Otherwise enqueue `(timestamp, message)`, add to the set, return `true`.

### Complexity

- **Time:** O(1) amortized — each entry is enqueued and dequeued at most once.
- **Space:** O(w) — entries within the 10-second window.

### Code

```go
type QueueSetLogger struct {
	window *list.List      // FIFO of entries printed in the last 10 seconds
	inSet  map[string]bool // messages currently inside the window
}

type entry struct {
	timestamp int
	message   string
}

func NewQueueSetLogger() *QueueSetLogger {
	return &QueueSetLogger{window: list.New(), inSet: make(map[string]bool)}
}

func (l *QueueSetLogger) ShouldPrintMessage(timestamp int, message string) bool {
	// Evict everything that has aged out of the 10-second window.
	for l.window.Len() > 0 {
		front := l.window.Front().Value.(entry)
		if front.timestamp <= timestamp-10 {
			l.window.Remove(l.window.Front()) // drop the stale entry
			delete(l.inSet, front.message)    // and forget its message
		} else {
			break // front is still fresh ⇒ the rest are too (FIFO by time)
		}
	}
	if l.inSet[message] {
		return false // still inside the window — suppress
	}
	// Record this print and admit it into the window.
	l.window.PushBack(entry{timestamp, message})
	l.inSet[message] = true
	return true
}
```

### Dry Run

Example 1 call sequence (window keeps entries with `t > timestamp - 10`):

| Call | ts | msg | evict (t ≤ ts-10) | in set? | print? | window after |
|------|----|-----|-------------------|---------|--------|--------------|
| 1 | 1 | foo | none | no | **true** | [(1,foo)] |
| 2 | 2 | bar | none | no | **true** | [(1,foo),(2,bar)] |
| 3 | 3 | foo | none (1 > -7) | yes | **false** | unchanged |
| 4 | 8 | bar | none (1 > -2) | yes | **false** | unchanged |
| 5 | 10 | foo | (1,foo): 1 ≤ 0? no | yes | **false** | unchanged |
| 6 | 11 | foo | (1,foo): 1 ≤ 1 → evict foo | no | **true** | [(2,bar),(11,foo)] |

Output: `[true, true, false, false, false, true]` ✔

---

## Key Takeaways

- **"Not seen in the last T seconds" → store next-allowed = t + T.** Turning a window constraint into a single scalar per key is the cleanest design and gives O(1) with a plain hash map.
- **Guaranteed chronological input** is what lets both approaches be trivial — no need to handle out-of-order timestamps.
- **Hash-map vs. queue trade-off**: the map is O(1) but never frees old messages; the queue+set frees anything outside the window at the cost of eviction bookkeeping. Pick based on whether messages are unbounded and rarely repeat.
- Go's map zero-value (`0` for absent `int`) removes the need for an explicit "seen before?" branch.

---

## Related Problems

- LeetCode #362 — Design Hit Counter (sliding-window counting)
- LeetCode #379 — Design Phone Directory (design with a data stream)
- LeetCode #146 — LRU Cache (hash map + eviction order)
- LeetCode #1352 — Product of the Last K Numbers (streaming design)
