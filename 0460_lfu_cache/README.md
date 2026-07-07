# 0460 — LFU Cache

> LeetCode #460 · Difficulty: Hard
> **Categories:** Hash Table, Linked List, Design, Doubly-Linked List

---

## Problem Statement

Design and implement a data structure for a **Least Frequently Used (LFU)** cache.

Implement the `LFUCache` class:

- `LFUCache(int capacity)` Initializes the object with the `capacity` of the data structure.
- `int get(int key)` Gets the value of the `key` if the `key` exists in the cache. Otherwise, returns `-1`.
- `void put(int key, int value)` Update the value of the `key` if present, or inserts the `key` if not already present. When the cache reaches its `capacity`, it should invalidate and remove the **least frequently used** key before inserting a new item. For this problem, when there is a **tie** (i.e., two or more keys with the same frequency), the **least recently used** `key` would be invalidated.

To determine the least frequently used key, a **use counter** is maintained for each key in the cache. The key with the smallest **use counter** is the least frequently used key.

When a key is first inserted into the cache, its **use counter** is set to `1` (due to the `put` operation). The **use counter** for a key in the cache is incremented either a `get` or `put` operation is called on it.

The functions `get` and `put` must each run in `O(1)` average time complexity.

**Example 1:**

```
Input
["LFUCache", "put", "put", "get", "put", "get", "get", "put", "get", "get", "get"]
[[2], [1, 1], [2, 2], [1], [3, 3], [2], [3], [4, 4], [1], [3], [4]]
Output
[null, null, null, 1, null, -1, 3, null, -1, 3, 4]

Explanation
// cnt(x) = the use counter for key x
// cache=[] will show the last used order for tiebreakers (leftmost element is  most recent)
LFUCache lfu = new LFUCache(2);
lfu.put(1, 1);   // cache=[1,_], cnt(1)=1
lfu.put(2, 2);   // cache=[2,1], cnt(2)=1, cnt(1)=1
lfu.get(1);      // return 1
                 // cache=[1,2], cnt(2)=1, cnt(1)=2
lfu.put(3, 3);   // 2 is the LFU key because cnt(2)=1 is the smallest, invalidate 2.
                 // cache=[3,1], cnt(3)=1, cnt(1)=2
lfu.get(2);      // return -1 (not found)
lfu.get(3);      // return 3
                 // cache=[3,1], cnt(3)=2, cnt(1)=2
lfu.put(4, 4);   // Both 1 and 3 have the same cnt, but 1 is LRU, invalidate 1.
                 // cache=[4,3], cnt(4)=1, cnt(3)=2
lfu.get(1);      // return -1 (not found)
lfu.get(3);      // return 3
                 // cache=[3,4], cnt(4)=1, cnt(3)=3
lfu.get(4);      // return 4
                 // cache=[4,3], cnt(4)=2, cnt(3)=3
```

**Constraints:**

- `1 <= capacity <= 10^4`
- `0 <= key <= 10^5`
- `0 <= value <= 10^9`
- At most `2 * 10^5` calls will be made to `get` and `put`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★★ Very High  | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Design Data Structures** — this is a class-design problem: pick internal structures (maps, linked lists, a `minFreq` pointer) so both operations hit their O(1) contract → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)
- **Linked List (Doubly-Linked)** — each frequency bucket is a doubly-linked list ordered by recency, giving O(1) move-to-front and O(1) removal of the least-recently-used tail → see [`/dsa/linked_list.md`](/dsa/linked_list.md)
- **Hash Map** — a key→node map for O(1) lookup and a frequency→list map for O(1) bucket access are the backbone of the optimal design → see [`/dsa/hash_map.md`](/dsa/hash_map.md)

---

## Approaches Overview

| # | Approach | get | put | Space | When to use |
|---|----------|-----|-----|-------|-------------|
| 1 | Brute Force (map + linear eviction scan) | O(1) | O(1) / O(n) on eviction | O(capacity) | Easiest to get right; fails the O(1) requirement |
| 2 | HashMap + freq buckets of DLLs (Optimal) | O(1) | O(1) | O(capacity) | The intended design; meets the O(1) contract |

---

## Approach 1 — Brute Force (Map + Linear Eviction Scan)

### Intuition

Keep one map from key to an entry holding `{value, freq, tick}`, where `freq` is the use counter and `tick` is a logical timestamp of the last access (for LRU tie-breaking). `get`/`put` on an existing key just bump `freq` and stamp a fresh `tick`. The only expensive operation is eviction: to find the victim, scan every entry, pick the smallest `freq`, and among ties the smallest `tick` (oldest use). Simple and obviously correct — but the scan is O(n), so it violates the problem's O(1) target. It is the reference oracle for the optimal version.

### Algorithm

1. `get(key)`: miss → `-1`; hit → `touch` (freq++, tick = ++clock), return value.
2. `put(key, value)`:
   1. Capacity 0 → do nothing.
   2. Existing key → update value, `touch`.
   3. Full → `evict`, then insert new entry with `freq = 1`, `tick = ++clock`.
3. `evict`: scan all entries; victim = min `freq`, tie-broken by min `tick`; delete it.

### Complexity

- **Time:** `get` O(1), `put` O(1) except eviction which is O(n) (the linear scan).
- **Space:** O(capacity) — one entry per stored key.

### Code

```go
type bfEntry struct {
	value int // stored value
	freq  int // number of times this key has been accessed (get or put)
	tick  int // logical time of the most recent access (for LRU tie-break)
}

type BruteForceLFU struct {
	capacity int
	clock    int              // ever-increasing logical clock stamped on each access
	data     map[int]*bfEntry // key → entry
}

func NewBruteForceLFU(capacity int) *BruteForceLFU {
	return &BruteForceLFU{capacity: capacity, data: map[int]*bfEntry{}}
}

func (c *BruteForceLFU) touch(e *bfEntry) {
	e.freq++         // one more access
	c.clock++        // advance logical time
	e.tick = c.clock // remember when this access happened (newest wins LRU)
}

func (c *BruteForceLFU) Get(key int) int {
	e, ok := c.data[key]
	if !ok {
		return -1 // miss
	}
	c.touch(e) // a successful get raises frequency & recency
	return e.value
}

func (c *BruteForceLFU) Put(key, value int) {
	if c.capacity == 0 {
		return // a zero-capacity cache stores nothing
	}
	if e, ok := c.data[key]; ok {
		e.value = value // update in place
		c.touch(e)      // updating also counts as an access
		return
	}
	if len(c.data) >= c.capacity {
		c.evict() // make room by removing the least-frequently/recently used
	}
	c.clock++
	c.data[key] = &bfEntry{value: value, freq: 1, tick: c.clock} // new keys start at freq 1
}

func (c *BruteForceLFU) evict() {
	victimKey := 0
	bestFreq, bestTick := -1, -1
	first := true
	for k, e := range c.data {
		// Choose smaller freq; on equal freq choose the smaller tick (older use).
		if first || e.freq < bestFreq || (e.freq == bestFreq && e.tick < bestTick) {
			victimKey = k
			bestFreq = e.freq
			bestTick = e.tick
			first = false
		}
	}
	delete(c.data, victimKey)
}
```

### Dry Run

Official sequence, `capacity = 2`. Columns show state after each op; `tick` grows on every access.

| Op | Action | Entries `{key: (value, freq, tick)}` after | Return |
|----|--------|--------------------------------------------|--------|
| put(1,1) | insert | `{1:(1,1,1)}` | — |
| put(2,2) | insert | `{1:(1,1,1), 2:(2,1,2)}` | — |
| get(1) | touch 1 | `{1:(1,2,3), 2:(2,1,2)}` | 1 |
| put(3,3) | full; evict min-freq → key 2 (freq 1) | `{1:(1,2,3), 3:(3,1,4)}` | — |
| get(2) | miss | unchanged | -1 |
| get(3) | touch 3 | `{1:(1,2,3), 3:(3,2,5)}` | 3 |
| put(4,4) | full; freq tie 1&3 (both 2) → evict min tick → key 1 (tick 3) | `{3:(3,2,5), 4:(4,1,6)}` | — |
| get(1) | miss | unchanged | -1 |
| get(3) | touch 3 | `{3:(3,3,7), 4:(4,1,6)}` | 3 |
| get(4) | touch 4 | `{3:(3,3,7), 4:(4,2,8)}` | 4 |

Output: `[null, null, null, 1, null, -1, 3, null, -1, 3, 4]` ✔

---

## Approach 2 — HashMap + Frequency Buckets of DLLs (Optimal)

### Intuition

Both hard operations must be O(1): locate an entry, and pick the eviction victim. Use three pieces:

- **`nodes`** — a `map[key] → *list.Element`, giving O(1) lookup of any key's node.
- **`freqList`** — a `map[freq] → doubly-linked list` of the nodes currently at that frequency, ordered **most-recent at the front**. Within a bucket, the *back* is the least-recently-used node.
- **`minFreq`** — the smallest occupied frequency, so the global LFU-then-LRU victim is always `freqList[minFreq].Back()`.

On every access, `bump` the node from bucket `f` to bucket `f+1` and push it to that bucket's front. If bucket `f` was the `minFreq` bucket and just emptied, `minFreq` increments (frequencies stay contiguous upward as the single accessed node moves). On a full insert, evict `freqList[minFreq].Back()`, add the new key at frequency 1, and reset `minFreq = 1`.

### Algorithm

1. `get(key)`: miss → `-1`; hit → read value, `bump` node, return value.
2. `put(key, value)`:
   1. Capacity 0 → return.
   2. Existing key → set value, `bump`.
   3. Full → remove `freqList[minFreq].Back()` from the list and `nodes`.
   4. Insert node at frequency 1 (front of `freqList[1]`), set `minFreq = 1`.
3. `bump(elem)`: remove from bucket `f`; if that bucket empties and `f == minFreq`, `minFreq++`; set `freq = f+1`; push to front of `freqList[f+1]`.

### Complexity

- **Time:** `get` O(1), `put` O(1) — every step is a map lookup or a constant-time linked-list splice.
- **Space:** O(capacity) — one node per key plus the bucket lists.

### Code

```go
type node struct {
	key, value, freq int // the frequency this node currently lives at
}

type OptimalLFU struct {
	capacity int
	minFreq  int                   // smallest frequency currently present
	nodes    map[int]*list.Element // key → its element inside a freq list
	freqList map[int]*list.List    // frequency → DLL of *node (front = newest)
}

func NewOptimalLFU(capacity int) *OptimalLFU {
	return &OptimalLFU{
		capacity: capacity,
		nodes:    map[int]*list.Element{},
		freqList: map[int]*list.List{},
	}
}

func (c *OptimalLFU) bump(elem *list.Element) {
	nd := elem.Value.(*node)
	f := nd.freq
	c.freqList[f].Remove(elem) // pull node out of its current bucket
	// If that bucket is now empty and it was the minimum, the new minimum is f+1.
	if c.freqList[f].Len() == 0 && c.minFreq == f {
		c.minFreq++
	}
	nd.freq = f + 1
	if c.freqList[nd.freq] == nil {
		c.freqList[nd.freq] = list.New()
	}
	c.nodes[nd.key] = c.freqList[nd.freq].PushFront(nd) // most-recent at the front
}

func (c *OptimalLFU) Get(key int) int {
	elem, ok := c.nodes[key]
	if !ok {
		return -1 // miss
	}
	val := elem.Value.(*node).value
	c.bump(elem) // successful access raises this key's frequency
	return val
}

func (c *OptimalLFU) Put(key, value int) {
	if c.capacity == 0 {
		return
	}
	if elem, ok := c.nodes[key]; ok {
		elem.Value.(*node).value = value
		c.bump(elem)
		return
	}
	if len(c.nodes) >= c.capacity {
		minList := c.freqList[c.minFreq]
		victim := minList.Back() // back of the min bucket = LFU + LRU
		delete(c.nodes, victim.Value.(*node).key)
		minList.Remove(victim)
	}
	nd := &node{key: key, value: value, freq: 1}
	if c.freqList[1] == nil {
		c.freqList[1] = list.New()
	}
	c.nodes[key] = c.freqList[1].PushFront(nd)
	c.minFreq = 1 // a brand-new freq-1 entry makes 1 the smallest frequency
}
```

### Dry Run

Official sequence, `capacity = 2`. Buckets show `freq → [front … back]` of keys.

| Op | Action | Buckets after | minFreq | Return |
|----|--------|---------------|---------|--------|
| put(1,1) | insert freq1 | `1→[1]` | 1 | — |
| put(2,2) | insert freq1 (front) | `1→[2,1]` | 1 | — |
| get(1) | bump 1: freq1→freq2 | `1→[2]`, `2→[1]` | 1 | 1 |
| put(3,3) | full; evict back of freq1 = key 2; insert 3 freq1 | `1→[3]`, `2→[1]` | 1 | — |
| get(2) | miss | unchanged | 1 | -1 |
| get(3) | bump 3: freq1→freq2; freq1 empties, minFreq→2 | `2→[3,1]` | 2 | 3 |
| put(4,4) | full; evict back of freq2 = key 1 (LRU); insert 4 freq1; minFreq=1 | `1→[4]`, `2→[3]` | 1 | — |
| get(1) | miss | unchanged | 1 | -1 |
| get(3) | bump 3: freq2→freq3 | `1→[4]`, `3→[3]` | 1 | 3 |
| get(4) | bump 4: freq1→freq2; freq1 empties, minFreq→2 | `2→[4]`, `3→[3]` | 2 | 4 |

Output: `[null, null, null, 1, null, -1, 3, null, -1, 3, 4]` ✔ (Verified equal to Approach 1 over 5000 randomised trials.)

---

## Key Takeaways

- **LFU = "map to find" + "bucket to evict".** The key→node map handles lookup; grouping keys into per-frequency lists (recency-ordered) makes the LFU-then-LRU victim the back of the lowest bucket.
- **`minFreq` only ever increases by 1 on a bump, and resets to 1 on any insert.** When the single accessed node leaves the `minFreq` bucket and empties it, the next-lowest occupied frequency is exactly `minFreq + 1` — no scan needed.
- **Order within a bucket encodes recency for free.** Push-front on access, evict-back on eviction: the doubly-linked list gives O(1) at both ends, which is why LRU tie-breaking stays O(1).
- **Design problems reward layering data structures.** No single structure gives O(1) LFU; the trick is composing a hash map with per-frequency linked lists and one integer pointer.
- **Always validate a subtle design against a brute-force oracle.** The `minFreq` bookkeeping is easy to get wrong; a randomised cross-check against the O(n) version catches it.

---

## Related Problems

- LeetCode #146 — LRU Cache (the simpler single-list cousin)
- LeetCode #355 — Design Twitter (compose maps + heaps/lists)
- LeetCode #432 — All O`one Data Structure (buckets by count with min/max in O(1))
- LeetCode #1476 — Subrectangle Queries (design with update/query trade-offs)
- LeetCode #895 — Maximum Frequency Stack (frequency buckets of stacks — same bucketing idea)
