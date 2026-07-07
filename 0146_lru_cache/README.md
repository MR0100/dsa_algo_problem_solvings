# 0146 — LRU Cache

> LeetCode #146 · Difficulty: Medium
> **Categories:** Hash Table, Linked List, Design, Doubly-Linked List

---

## Problem Statement

Design a data structure that follows the constraints of a **Least Recently Used (LRU) cache**.

Implement the `LRUCache` class:

- `LRUCache(int capacity)` Initialize the LRU cache with **positive** size `capacity`.
- `int get(int key)` Return the value of the `key` if the key exists, otherwise return `-1`.
- `void put(int key, int value)` Update the value of the `key` if the `key` exists. Otherwise, add the `key-value` pair to the cache. If the number of keys exceeds the `capacity` from this operation, **evict** the least recently used key.

The functions `get` and `put` must each run in `O(1)` average time complexity.

**Example 1:**
```
Input
["LRUCache", "put", "put", "get", "put", "get", "put", "get", "get", "get"]
[[2], [1, 1], [2, 2], [1], [3, 3], [2], [4, 4], [1], [3], [4]]
Output
[null, null, null, 1, null, -1, null, -1, 3, 4]

Explanation
LRUCache lRUCache = new LRUCache(2);
lRUCache.put(1, 1); // cache is {1=1}
lRUCache.put(2, 2); // cache is {1=1, 2=2}
lRUCache.get(1);    // return 1
lRUCache.put(3, 3); // LRU key was 2, evicts key 2, cache is {1=1, 3=3}
lRUCache.get(2);    // returns -1 (not found)
lRUCache.put(4, 4); // LRU key was 1, evicts key 1, cache is {4=4, 3=3}
lRUCache.get(1);    // return -1 (not found)
lRUCache.get(3);    // return 3
lRUCache.get(4);    // return 4
```

**Constraints:**
- `1 <= capacity <= 3000`
- `0 <= key <= 10^4`
- `0 <= value <= 10^5`
- At most `2 * 10^5` calls will be made to `get` and `put`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★★ Very High  | 2024          |
| Microsoft  | ★★★★★ Very High  | 2024          |
| Facebook   | ★★★★★ Very High  | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Apple      | ★★★★☆ High       | 2024          |
| Bloomberg  | ★★★★☆ High       | 2024          |
| Oracle     | ★★★☆☆ Medium     | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Map** — O(1) key → node lookup; the "find it fast" half of the design → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Doubly Linked List** — O(1) unlink/re-insert of an arbitrary node; the "reorder it fast" half → see [`/dsa/linked_list.md`](/dsa/linked_list.md)
- **Sentinel (dummy) nodes** — dummy head/tail eliminate every nil-pointer edge case → see [`/dsa/linked_list.md`](/dsa/linked_list.md)

---

## Approaches Overview

| # | Approach | Time (per op) | Space | When to use |
|---|----------|---------------|-------|-------------|
| 1 | Brute Force (Ordered Slice) | O(capacity) | O(capacity) | Baseline; fine only for tiny capacities |
| 2 | Hash Map + Doubly Linked List (Optimal) | O(1) | O(capacity) | The interview answer — build it by hand |
| 3 | Hash Map + `container/list` (Optimal, Stdlib) | O(1) | O(capacity) | Idiomatic production Go |

---

## Approach 1 — Brute Force (Ordered Slice)

### Intuition
"Least recently used" is fundamentally an ordering by recency. The most literal model is a slice kept in that order: index `0` is the least recently used entry, the last index is the most recently used. Every operation linearly scans the slice; every touched entry is moved to the back; eviction drops index `0`.

### Algorithm
1. Store entries as a slice of `{key, value}` pairs, ordered LRU → MRU.
2. **Get(key):** linear-scan the slice. If found, remove the entry from its slot, re-append it at the back (now MRU), return its value. Otherwise return `-1`.
3. **Put(key, value):** linear-scan. If the key exists, remove the old slot and re-append with the new value.
4. If the key is new, append `{key, value}` at the back. If `len(entries) > capacity`, drop `entries[0]` — that is the LRU eviction.

### Complexity
- **Time:** O(capacity) per operation — the scan and the slice shift each touch up to `capacity` entries.
- **Space:** O(capacity) — one pair per cached key.

### Code
```go
type LRUCacheSlice struct {
	capacity int
	entries  []kvPair // entries[0] = LRU, entries[len-1] = MRU
}

type kvPair struct {
	key, value int
}

func NewLRUCacheSlice(capacity int) *LRUCacheSlice {
	return &LRUCacheSlice{capacity: capacity, entries: make([]kvPair, 0, capacity)}
}

func (c *LRUCacheSlice) Get(key int) int {
	for i, e := range c.entries {
		if e.key == key {
			c.entries = append(c.entries[:i], c.entries[i+1:]...) // remove
			c.entries = append(c.entries, e)                      // re-append as MRU
			return e.value
		}
	}
	return -1
}

func (c *LRUCacheSlice) Put(key int, value int) {
	for i, e := range c.entries {
		if e.key == key {
			c.entries = append(c.entries[:i], c.entries[i+1:]...)
			c.entries = append(c.entries, kvPair{key, value})
			return
		}
	}
	c.entries = append(c.entries, kvPair{key, value})
	if len(c.entries) > c.capacity {
		c.entries = c.entries[1:] // evict LRU at the front
	}
}
```

### Dry Run
Example 1, `capacity = 2`. Slice shown LRU → MRU.

| Step | Operation  | Slice before        | Action                              | Slice after         | Returns |
|------|------------|---------------------|-------------------------------------|---------------------|---------|
| 1    | `put(1,1)` | `[]`                | new key → append                    | `[(1,1)]`           | null    |
| 2    | `put(2,2)` | `[(1,1)]`           | new key → append                    | `[(1,1),(2,2)]`     | null    |
| 3    | `get(1)`   | `[(1,1),(2,2)]`     | found at 0 → move to back           | `[(2,2),(1,1)]`     | **1**   |
| 4    | `put(3,3)` | `[(2,2),(1,1)]`     | new key → append; len 3 > 2 → drop front `(2,2)` | `[(1,1),(3,3)]` | null |
| 5    | `get(2)`   | `[(1,1),(3,3)]`     | not found                           | `[(1,1),(3,3)]`     | **-1**  |
| 6    | `put(4,4)` | `[(1,1),(3,3)]`     | new key → append; drop front `(1,1)`| `[(3,3),(4,4)]`     | null    |
| 7    | `get(1)`   | `[(3,3),(4,4)]`     | not found                           | `[(3,3),(4,4)]`     | **-1**  |
| 8    | `get(3)`   | `[(3,3),(4,4)]`     | found → move to back                | `[(4,4),(3,3)]`     | **3**   |
| 9    | `get(4)`   | `[(4,4),(3,3)]`     | found → move to back                | `[(3,3),(4,4)]`     | **4**   |

Output: `[null, null, null, 1, null, -1, null, -1, 3, 4]` ✓

---

## Approach 2 — Hash Map + Doubly Linked List (Optimal)

### Intuition
The problem demands O(1) for both operations, which decomposes into two sub-requirements:

1. **Find a node by key in O(1)** → hash map `map[key]*node`.
2. **Reorder by recency in O(1)** → a *doubly* linked list. Removing a node from the middle needs its `prev` pointer, so a singly linked list will not do. Keep the list ordered `head → LRU … MRU ← tail`, using **sentinel** head/tail nodes so insertion and unlinking never branch on nil.

Each node stores its **key as well as its value** — when we evict the node at `head.next` we need the key to delete the map entry.

### Algorithm
1. Create sentinel nodes `head`, `tail` with `head.next = tail`, `tail.prev = head`; create empty map `index`.
2. **Get(key):** if `key` not in `index`, return `-1`. Otherwise `unlink(node)` (splice `prev`/`next` around it) and `pushBack(node)` (insert just before `tail`), then return `node.value`.
3. **Put(key, value)** when the key exists: overwrite `node.value`, unlink + pushBack (promote to MRU).
4. **Put(key, value)** when the key is new: create a node, `pushBack` it, add to `index`.
5. If `len(index) > capacity`: the LRU is `head.next`; unlink it and `delete(index, lru.key)`.

### Complexity
- **Time:** O(1) per operation — one hash lookup plus a constant number of pointer assignments.
- **Space:** O(capacity) — one map entry and one list node per cached key.

### Code
```go
type LRUCacheDLL struct {
	capacity int
	index    map[int]*dllNode // key → node
	head     *dllNode         // sentinel: head.next is the LRU
	tail     *dllNode         // sentinel: tail.prev is the MRU
}

type dllNode struct {
	key, value int
	prev, next *dllNode
}

func NewLRUCacheDLL(capacity int) *LRUCacheDLL {
	head, tail := &dllNode{}, &dllNode{}
	head.next = tail
	tail.prev = head
	return &LRUCacheDLL{capacity: capacity, index: make(map[int]*dllNode, capacity), head: head, tail: tail}
}

func (c *LRUCacheDLL) unlink(n *dllNode) {
	n.prev.next = n.next
	n.next.prev = n.prev
}

func (c *LRUCacheDLL) pushBack(n *dllNode) {
	n.prev = c.tail.prev
	n.next = c.tail
	c.tail.prev.next = n
	c.tail.prev = n
}

func (c *LRUCacheDLL) Get(key int) int {
	n, ok := c.index[key]
	if !ok {
		return -1
	}
	c.unlink(n)
	c.pushBack(n) // promote to most recently used
	return n.value
}

func (c *LRUCacheDLL) Put(key int, value int) {
	if n, ok := c.index[key]; ok {
		n.value = value
		c.unlink(n)
		c.pushBack(n)
		return
	}
	n := &dllNode{key: key, value: value}
	c.pushBack(n)
	c.index[key] = n
	if len(c.index) > c.capacity {
		lru := c.head.next // least recently used
		c.unlink(lru)
		delete(c.index, lru.key)
	}
}
```

### Dry Run
Example 1, `capacity = 2`. List shown `head → … → tail` (left = LRU).

| Step | Operation  | List before (LRU→MRU) | Map before  | Action                                    | List after (LRU→MRU) | Returns |
|------|------------|------------------------|-------------|-------------------------------------------|----------------------|---------|
| 1    | `put(1,1)` | `∅`                    | `{}`        | new node (1,1) pushBack, map[1]=n1        | `(1,1)`              | null    |
| 2    | `put(2,2)` | `(1,1)`                | `{1}`       | new node (2,2) pushBack, map[2]=n2        | `(1,1)·(2,2)`        | null    |
| 3    | `get(1)`   | `(1,1)·(2,2)`          | `{1,2}`     | hit n1 → unlink + pushBack                | `(2,2)·(1,1)`        | **1**   |
| 4    | `put(3,3)` | `(2,2)·(1,1)`          | `{1,2}`     | new node (3,3); size 3 > 2 → evict head.next = (2,2), delete key 2 | `(1,1)·(3,3)` | null |
| 5    | `get(2)`   | `(1,1)·(3,3)`          | `{1,3}`     | key 2 not in map                          | `(1,1)·(3,3)`        | **-1**  |
| 6    | `put(4,4)` | `(1,1)·(3,3)`          | `{1,3}`     | new node (4,4); evict (1,1), delete key 1 | `(3,3)·(4,4)`        | null    |
| 7    | `get(1)`   | `(3,3)·(4,4)`          | `{3,4}`     | key 1 not in map                          | `(3,3)·(4,4)`        | **-1**  |
| 8    | `get(3)`   | `(3,3)·(4,4)`          | `{3,4}`     | hit n3 → promote                          | `(4,4)·(3,3)`        | **3**   |
| 9    | `get(4)`   | `(4,4)·(3,3)`          | `{3,4}`     | hit n4 → promote                          | `(3,3)·(4,4)`        | **4**   |

Output: `[null, null, null, 1, null, -1, null, -1, 3, 4]` ✓

---

## Approach 3 — Hash Map + `container/list` (Optimal, Stdlib)

### Intuition
Approach 2's doubly linked list already exists in Go's standard library as `container/list`, complete with O(1) `PushBack`, `MoveToBack`, and `Remove`. The design is byte-for-byte the same — the map now stores `*list.Element` instead of a hand-rolled node. This is what you would ship in production Go; the hand-rolled version is what you write on a whiteboard to prove you understand the pointers.

### Algorithm
1. Keep `order *list.List` (Front = LRU, Back = MRU) and `index map[int]*list.Element`.
2. **Get(key):** map lookup → if hit, `order.MoveToBack(elem)` and return the payload value; else `-1`.
3. **Put(key, value)** existing key: overwrite `elem.Value`, `MoveToBack`.
4. **Put(key, value)** new key: `index[key] = order.PushBack(kvPair{key, value})`.
5. If `order.Len() > capacity`: `front := order.Front()`, `order.Remove(front)`, delete `front`'s key from the map.

### Complexity
- **Time:** O(1) per operation — same pointer operations, performed by the stdlib.
- **Space:** O(capacity).

### Code
```go
type LRUCacheStdList struct {
	capacity int
	index    map[int]*list.Element // key → element
	order    *list.List            // Front() = LRU, Back() = MRU
}

func NewLRUCacheStdList(capacity int) *LRUCacheStdList {
	return &LRUCacheStdList{capacity: capacity, index: make(map[int]*list.Element, capacity), order: list.New()}
}

func (c *LRUCacheStdList) Get(key int) int {
	elem, ok := c.index[key]
	if !ok {
		return -1
	}
	c.order.MoveToBack(elem)
	return elem.Value.(kvPair).value
}

func (c *LRUCacheStdList) Put(key int, value int) {
	if elem, ok := c.index[key]; ok {
		elem.Value = kvPair{key, value}
		c.order.MoveToBack(elem)
		return
	}
	c.index[key] = c.order.PushBack(kvPair{key, value})
	if c.order.Len() > c.capacity {
		front := c.order.Front()
		c.order.Remove(front)
		delete(c.index, front.Value.(kvPair).key)
	}
}
```

### Dry Run
Identical state evolution to Approach 2 — only the mechanism differs (`MoveToBack` / `Remove` instead of manual pointer surgery).

| Step | Operation  | `order` (Front→Back) after | `index` keys after | Returns |
|------|------------|-----------------------------|--------------------|---------|
| 1    | `put(1,1)` | `(1,1)`                     | `{1}`              | null    |
| 2    | `put(2,2)` | `(1,1)·(2,2)`               | `{1,2}`            | null    |
| 3    | `get(1)`   | `(2,2)·(1,1)`               | `{1,2}`            | **1**   |
| 4    | `put(3,3)` | `(1,1)·(3,3)` (evict 2)     | `{1,3}`            | null    |
| 5    | `get(2)`   | `(1,1)·(3,3)`               | `{1,3}`            | **-1**  |
| 6    | `put(4,4)` | `(3,3)·(4,4)` (evict 1)     | `{3,4}`            | null    |
| 7    | `get(1)`   | `(3,3)·(4,4)`               | `{3,4}`            | **-1**  |
| 8    | `get(3)`   | `(4,4)·(3,3)`               | `{3,4}`            | **3**   |
| 9    | `get(4)`   | `(3,3)·(4,4)`               | `{3,4}`            | **4**   |

Output: `[null, null, null, 1, null, -1, null, -1, 3, 4]` ✓

---

## Key Takeaways

- **O(1) design problems decompose into capabilities.** "Find by key fast" → hash map; "reorder/evict by recency fast" → doubly linked list. Composite data structures are the standard answer to composite O(1) requirements.
- **Doubly, not singly:** removing an arbitrary node in O(1) requires a `prev` pointer. Whenever a problem needs "delete this exact node fast", think doubly linked list.
- **Sentinel head/tail nodes** turn 4 nil-check branches into 0. Always use dummies in linked-list design problems.
- **Store the key inside the node.** Eviction walks list → map (node tells you which map key to delete); lookups walk map → list. The two structures must reference each other both ways.
- **`get` mutates state** in an LRU — reading an entry promotes it. Easy to forget in interviews.
- Same pattern powers LeetCode #460 (LFU Cache — add a frequency dimension) and real systems (Redis' approximate LRU, OS page replacement, CDN caches).

---

## Related Problems

- LeetCode #460 — LFU Cache (same map + linked-list design, plus frequency buckets)
- LeetCode #432 — All O`one Data Structure (doubly linked list of buckets + hash map)
- LeetCode #380 — Insert Delete GetRandom O(1) (composite structure for composite O(1) ops)
- LeetCode #705/#706 — Design HashSet / HashMap (design-a-structure fundamentals)
- LeetCode #1472 — Design Browser History (recency-ordered navigation)
