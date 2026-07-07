# 0379 — Design Phone Directory

> LeetCode #379 · Difficulty: Medium
> **Categories:** Design, Hash Set, Queue, Linked List, Array

---

## Problem Statement

Design a phone directory that initially has `maxNumbers` empty slots that can store numbers. The directory should store numbers, check if a certain slot is empty or not, and empty a given slot.

Implement the `PhoneDirectory` class:

- `PhoneDirectory(int maxNumbers)` Initializes the phone directory with the number of available slots `maxNumbers`.
- `int get()` Provides a number that is not assigned to anyone. Returns `-1` if no number is available.
- `bool check(int number)` Returns `true` if the slot `number` is available and `false` otherwise.
- `void release(int number)` Recycles or releases the slot `number`.

**Example 1:**

```
Input
["PhoneDirectory", "get", "get", "check", "get", "check", "release", "check"]
[[3], [], [], [2], [], [2], [2], [2]]
Output
[null, 0, 1, true, 2, false, null, true]

Explanation
PhoneDirectory phoneDirectory = new PhoneDirectory(3);
phoneDirectory.get();      // It can return any available phone number. Here we assume it returns 0.
phoneDirectory.get();      // Assume it returns 1.
phoneDirectory.check(2);   // The number 2 is available, so return true.
phoneDirectory.get();      // It returns 2, the only number that is left.
phoneDirectory.check(2);   // The number 2 is no longer available, so return false.
phoneDirectory.release(2); // Release number 2 back to the pool.
phoneDirectory.check(2);   // Number 2 is available again, return true.
```

**Constraints:**

- `1 <= maxNumbers <= 10^4`
- `0 <= number < maxNumbers`
- At most `2 * 10^4` calls will be made to `get`, `check`, and `release`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Design Data Structures** — build a class with a small operation API and pick backing structures to hit the required time bounds → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)
- **Queue / Deque** — a FIFO pool of free numbers makes `get`/`release` O(1) → see [`/dsa/queue_deque.md`](/dsa/queue_deque.md)
- **Hash Map** — a used-flag lookup answers `check` in O(1) (an array here since keys are dense) → see [`/dsa/hash_map.md`](/dsa/hash_map.md)

---

## Approaches Overview

| # | Approach | get | check | release | Space | When to use |
|---|----------|-----|-------|---------|-------|-------------|
| 1 | Boolean Array + Linear Scan | O(maxNumbers) | O(1) | O(1) | O(maxNumbers) | Simplest; fine if get is rare |
| 2 | Queue of Free Slots + Used Set (Optimal) | O(1) | O(1) | O(1) | O(maxNumbers) | The intended answer; all ops O(1) |

---

## Approach 1 — Boolean Array + Linear Scan

### Intuition

Model each number with a boolean "in use" flag. `get` scans for the first free flag, `check` reads a flag, `release` clears a flag. Dead simple and obviously correct — but `get` re-scans from the start every time, costing O(maxNumbers) per call.

### Algorithm

1. `get`: scan `used[0..max)` for the first `false`; mark it `true`, return the index; `-1` if none.
2. `check(number)`: return `!used[number]`.
3. `release(number)`: set `used[number] = false`.

### Complexity

- **Time:** `get` O(maxNumbers), `check` O(1), `release` O(1).
- **Space:** O(maxNumbers) — the flag array.

### Code

```go
type LinearDirectory struct {
	used []bool
}

func NewLinearDirectory(maxNumbers int) *LinearDirectory {
	return &LinearDirectory{used: make([]bool, maxNumbers)}
}

func (d *LinearDirectory) Get() int {
	for i := range d.used {
		if !d.used[i] {
			d.used[i] = true
			return i
		}
	}
	return -1
}

func (d *LinearDirectory) Check(number int) bool {
	return !d.used[number]
}

func (d *LinearDirectory) Release(number int) {
	d.used[number] = false
}
```

### Dry Run

Example 1: `maxNumbers = 3`, `used = [F,F,F]`.

| Op | Scan / read | Result | used after |
|----|-------------|--------|-----------|
| get() | first free = 0 | 0 | [T,F,F] |
| get() | first free = 1 | 1 | [T,T,F] |
| check(2) | !used[2] | true | [T,T,F] |
| get() | first free = 2 | 2 | [T,T,T] |
| check(2) | !used[2] | false | [T,T,T] |
| release(2) | clear | null | [T,T,F] |
| check(2) | !used[2] | true | [T,T,F] |

Output `[null, 0, 1, true, 2, false, null, true]` ✔

---

## Approach 2 — Queue of Free Slots + Used Set (Optimal)

### Intuition

Instead of re-searching for a free slot, keep an explicit **pool** of free numbers in a FIFO queue. `get` pops the front in O(1); `release` pushes a number back. A `used[]` flag array answers `check` in O(1) and — critically — guards `release` so the same number is never queued twice (a double-release would otherwise corrupt the pool).

### Algorithm

1. Initialise the queue with every number `0..maxNumbers-1`.
2. `get`: if the queue is empty return `-1`; else pop the front, set `used[number] = true`, return it.
3. `check(number)`: return `!used[number]`.
4. `release(number)`: if `used[number]`, set it `false` and enqueue `number`; otherwise do nothing.

### Complexity

- **Time:** `get` O(1), `check` O(1), `release` O(1).
- **Space:** O(maxNumbers) — queue plus flag array.

### Code

```go
type QueueDirectory struct {
	free []int
	head int
	used []bool
}

func NewQueueDirectory(maxNumbers int) *QueueDirectory {
	free := make([]int, maxNumbers)
	for i := range free {
		free[i] = i
	}
	return &QueueDirectory{free: free, head: 0, used: make([]bool, maxNumbers)}
}

func (d *QueueDirectory) Get() int {
	if d.head >= len(d.free) {
		return -1
	}
	number := d.free[d.head]
	d.head++
	d.used[number] = true
	return number
}

func (d *QueueDirectory) Check(number int) bool {
	return !d.used[number]
}

func (d *QueueDirectory) Release(number int) {
	if !d.used[number] {
		return
	}
	d.used[number] = false
	d.free = append(d.free, number)
}
```

### Dry Run

Example 1: `maxNumbers = 3`, `free = [0,1,2] head=0`, `used = [F,F,F]`.

| Op | Action | Result | free / head | used |
|----|--------|--------|-------------|------|
| get() | pop front 0 | 0 | head=1 | [T,F,F] |
| get() | pop front 1 | 1 | head=2 | [T,T,F] |
| check(2) | !used[2] | true | — | [T,T,F] |
| get() | pop front 2 | 2 | head=3 | [T,T,T] |
| check(2) | !used[2] | false | — | [T,T,T] |
| release(2) | enqueue 2 | null | free=[0,1,2,2] head=3 | [T,T,F] |
| check(2) | !used[2] | true | — | [T,T,F] |

Output `[null, 0, 1, true, 2, false, null, true]` ✔

---

## Key Takeaways

- **Separate concerns with two structures:** a queue for "what's free" (fast allocation) and a flag array for "is X free" (fast membership). Neither alone gives O(1) on all three ops.
- **Guard `release`** against releasing an already-free number — otherwise the same number gets queued twice and `get` hands out a duplicate.
- Because keys are dense (`0..maxNumbers-1`), a boolean **array** beats a hash set — same O(1) with lower constant factors.
- Tracking the queue front with a moving `head` index avoids repeated reslice/copy churn on pops.

---

## Related Problems

- LeetCode #146 — LRU Cache (design with combined structures)
- LeetCode #380 — Insert Delete GetRandom O(1) (array + hash map for O(1) ops)
- LeetCode #705 — Design HashSet (membership design)
- LeetCode #1845 — Seat Reservation Manager (heap/queue pool of free ids)
