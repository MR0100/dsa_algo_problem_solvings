# 0432 — All O`one` Data Structure

> LeetCode #432 · Difficulty: Hard
> **Categories:** Hash Table, Linked List, Design, Doubly-Linked List

---

## Problem Statement

Design a data structure to store the strings' count with the ability to return the strings with minimum and maximum counts.

Implement the `AllOne` class:

- `AllOne()` Initializes the object of the data structure.
- `inc(String key)` Increments the count of the string `key` by `1`. If `key` does not exist in the data structure, insert it with count `1`.
- `dec(String key)` Decrements the count of the string `key` by `1`. If the count of `key` is `0` after the decrement, remove it from the data structure. It is guaranteed that `key` exists in the data structure before the decrement.
- `getMaxKey()` Returns one of the keys with the maximal count. If no element exists, return an empty string `""`.
- `getMinKey()` Returns one of the keys with the minimum count. If no element exists, return an empty string `""`.

**Note** that each function must run in `O(1)` average time complexity.

**Example 1:**

```
Input
["AllOne", "inc", "inc", "getMaxKey", "getMinKey", "inc", "getMaxKey", "getMinKey"]
[[], ["hello"], ["hello"], [], [], ["leet"], [], []]
Output
[null, null, null, "hello", "hello", null, "hello", "leet"]

Explanation
AllOne allOne = new AllOne();
allOne.inc("hello");
allOne.inc("hello");
allOne.getMaxKey(); // return "hello"
allOne.getMinKey(); // return "hello"
allOne.inc("leet");
allOne.getMaxKey(); // return "hello"
allOne.getMinKey(); // return "leet"
```

**Constraints:**

- `1 <= key.length <= 10`
- `key` consists of lowercase English letters.
- It is guaranteed that for each call to `dec`, `key` is existing in the data structure.
- At most `5 * 10^4` calls will be made to `inc`, `dec`, `getMaxKey`, and `getMinKey`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Facebook   | ★★★☆☆ Medium     | 2023          |
| LinkedIn   | ★★★☆☆ Medium     | 2023          |
| Uber       | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Design Data Structures** — the task is to compose primitives (a map + a linked list) into a new ADT that hits an O(1) budget on all four operations; the whole problem is a data-structure-design exercise → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)
- **Doubly-Linked List** — buckets are held on a doubly-linked list sorted by count so min = head, max = tail, and ±1 moves are adjacent-node splices (the same head/tail-sentinel technique as an LRU cache) → see [`/dsa/linked_list.md`](/dsa/linked_list.md)
- **Hash Map** — a `key → bucket` map turns "find this key's current count bucket" into an O(1) lookup, which is what lets Inc/Dec stay constant time → see [`/dsa/hash_map.md`](/dsa/hash_map.md)

---

## Approaches Overview

| # | Approach | Inc | Dec | GetMax | GetMin | Space | When to use |
|---|----------|-----|-----|--------|--------|-------|-------------|
| 1 | Hash Map + Linear Scan | O(1) | O(1) | O(n) | O(n) | O(n) | Quick to write; fine if min/max queries are rare |
| 2 | DLL of Count Buckets + Hash Map (Optimal) | O(1) | O(1) | O(1) | O(1) | O(n) | The intended answer — meets the "all O(1)" requirement |

`n` = number of distinct live keys.

---

## Approach 1 — Hash Map + Linear Scan (Brute Force)

### Intuition

The minimal structure that satisfies the API is a single map `key → count`. `Inc` and `Dec` are one-line map edits (O(1)). The only awkward part — returning *a* key with the min or max count — is solved by brute force: iterate the whole map and keep the best-so-far. This is obviously correct and a good baseline, but `getMaxKey`/`getMinKey` are O(n), which is exactly what the follow-up asks us to eliminate.

### Algorithm

1. `Inc(key)`: `counts[key]++` (an absent key defaults to `0`, so it becomes `1`).
2. `Dec(key)`: `counts[key]--`; if it hit `0`, delete the entry.
3. `GetMaxKey()`: scan every `(k, v)`, track the key with the largest `v`, return it (`""` if empty).
4. `GetMinKey()`: scan every `(k, v)`, track the key with the smallest `v`, return it (`""` if empty).

### Complexity

- **Time:** `Inc` O(1), `Dec` O(1); `GetMaxKey`/`GetMinKey` O(n) — each scans all live keys.
- **Space:** O(n) — one map entry per distinct live key.

### Code

```go
type BruteForceAllOne struct {
	counts map[string]int // key → current positive count
}

func NewBruteForceAllOne() *BruteForceAllOne {
	return &BruteForceAllOne{counts: map[string]int{}}
}

func (a *BruteForceAllOne) Inc(key string) {
	a.counts[key]++ // missing key defaults to 0, so this makes a new key 1
}

func (a *BruteForceAllOne) Dec(key string) {
	if _, ok := a.counts[key]; !ok {
		return
	}
	a.counts[key]--
	if a.counts[key] == 0 {
		delete(a.counts, key) // value 0 means the key no longer exists
	}
}

func (a *BruteForceAllOne) GetMaxKey() string {
	best, bestVal := "", -1
	for k, v := range a.counts {
		if v > bestVal {
			best, bestVal = k, v
		}
	}
	return best
}

func (a *BruteForceAllOne) GetMinKey() string {
	best, bestVal := "", int(^uint(0)>>1) // start at max int
	for k, v := range a.counts {
		if v < bestVal {
			best, bestVal = k, v
		}
	}
	return best
}
```

### Dry Run

Operation sequence from Example 1. `counts` is the map state after each op.

| Op | `counts` after | Return |
|----|----------------|--------|
| `inc("hello")` | `{hello:1}` | — |
| `inc("hello")` | `{hello:2}` | — |
| `getMaxKey()` | `{hello:2}` | scan → `"hello"` (val 2) |
| `getMinKey()` | `{hello:2}` | scan → `"hello"` (val 2) |
| `inc("leet")` | `{hello:2, leet:1}` | — |
| `getMaxKey()` | `{hello:2, leet:1}` | scan → `"hello"` (2 > 1) |
| `getMinKey()` | `{hello:2, leet:1}` | scan → `"leet"` (1 < 2) |

Output: `[null, null, null, "hello", "hello", null, "hello", "leet"]` ✔

---

## Approach 2 — DLL of Count Buckets + Hash Map (Optimal)

### Intuition

The brute force is slow only because it *searches* for the min/max count. Remove the search by keeping counts **sorted structurally**. Group all keys that currently share a count into one **bucket**, and keep the buckets on a **doubly-linked list sorted ascending by count**. Then:

- the **minimum**-count bucket is always right after the head sentinel — O(1),
- the **maximum**-count bucket is always right before the tail sentinel — O(1).

A hash map `key → bucket` finds any key's current bucket in O(1). Because `Inc`/`Dec` change a count by exactly `±1`, the destination bucket (`count±1`) is an **adjacent** node — so moving a key is a constant number of pointer splices: create the neighbour bucket if it doesn't exist yet, move the key, and unlink the old bucket if it became empty. Two sentinels (`head`, `tail`) mean the ends never need null-checks. Every operation is O(1).

### Algorithm

1. `Inc(key)`:
   - If `key` exists at bucket `cur` (count `c`): ensure `cur.next` is a bucket with count `c+1` (insert one after `cur` if not); add `key` there, remove it from `cur`.
   - Else (new key, target count `1`): ensure `head.next` is a count-`1` bucket (insert one after `head` if not); add `key` there.
   - If `cur` became empty, unlink it.
2. `Dec(key)`: let `cur` be `key`'s bucket (count `c`).
   - If `c == 1`: remove `key` from the map and from `cur` entirely.
   - Else: ensure `cur.prev` is a bucket with count `c-1` (insert one before `cur` if not); move `key` there.
   - If `cur` became empty, unlink it.
3. `GetMaxKey()`: any key in `tail.prev.keys`, or `""` if the list is empty (`tail.prev == head`).
4. `GetMinKey()`: any key in `head.next.keys`, or `""` if the list is empty (`head.next == tail`).

### Complexity

- **Time:** `Inc`, `Dec`, `GetMaxKey`, `GetMinKey` all O(1) — each is a fixed number of map operations and adjacent pointer splices; the `±1` step guarantees the target bucket is a neighbour, so there is never a scan.
- **Space:** O(n) — n live keys distributed over at most n buckets, plus the `key → bucket` map.

### Code

```go
type bucket struct {
	count int                 // the shared count of every key in this bucket
	keys  map[string]struct{} // set of keys currently holding `count`
	prev  *bucket             // neighbour with a smaller count
	next  *bucket             // neighbour with a larger count
}

type AllOne struct {
	keyBucket map[string]*bucket // key → the bucket that currently holds it
	head      *bucket            // sentinel BEFORE the smallest-count bucket
	tail      *bucket            // sentinel AFTER the largest-count bucket
}

func NewAllOne() *AllOne {
	head := &bucket{keys: map[string]struct{}{}}
	tail := &bucket{keys: map[string]struct{}{}}
	head.next = tail
	tail.prev = head
	return &AllOne{keyBucket: map[string]*bucket{}, head: head, tail: tail}
}

// insertAfter splices a new bucket with `count` right after `prev`.
func (a *AllOne) insertAfter(prev *bucket, count int) *bucket {
	b := &bucket{count: count, keys: map[string]struct{}{}}
	b.prev = prev
	b.next = prev.next
	prev.next.prev = b
	prev.next = b
	return b
}

// remove unlinks an (empty) bucket.
func (a *AllOne) remove(b *bucket) {
	b.prev.next = b.next
	b.next.prev = b.prev
}

func (a *AllOne) Inc(key string) {
	if cur, ok := a.keyBucket[key]; ok {
		next := cur.next
		if next == a.tail || next.count != cur.count+1 {
			next = a.insertAfter(cur, cur.count+1) // make the count+1 neighbour
		}
		next.keys[key] = struct{}{}
		a.keyBucket[key] = next
		delete(cur.keys, key)
		if len(cur.keys) == 0 {
			a.remove(cur)
		}
	} else {
		first := a.head.next
		if first == a.tail || first.count != 1 {
			first = a.insertAfter(a.head, 1) // count-1 bucket at the front
		}
		first.keys[key] = struct{}{}
		a.keyBucket[key] = first
	}
}

func (a *AllOne) Dec(key string) {
	cur, ok := a.keyBucket[key]
	if !ok {
		return
	}
	if cur.count == 1 {
		delete(cur.keys, key)
		delete(a.keyBucket, key) // count would hit 0 → key disappears
	} else {
		prev := cur.prev
		if prev == a.head || prev.count != cur.count-1 {
			prev = a.insertAfter(cur.prev, cur.count-1) // count-1 neighbour
		}
		prev.keys[key] = struct{}{}
		a.keyBucket[key] = prev
		delete(cur.keys, key)
	}
	if len(cur.keys) == 0 {
		a.remove(cur)
	}
}

func (a *AllOne) GetMaxKey() string {
	if a.tail.prev == a.head {
		return ""
	}
	for k := range a.tail.prev.keys {
		return k
	}
	return ""
}

func (a *AllOne) GetMinKey() string {
	if a.head.next == a.tail {
		return ""
	}
	for k := range a.head.next.keys {
		return k
	}
	return ""
}
```

### Dry Run

Example 1. List shown as `head <-> [count: {keys}] <-> ... <-> tail`.

| Op | List after | Return |
|----|-----------|--------|
| init | `head <-> tail` | — |
| `inc("hello")` | new key → count-1 bucket after head: `head <-> [1:{hello}] <-> tail` | — |
| `inc("hello")` | hello at 1 → move to count-2 neighbour, drop empty 1-bucket: `head <-> [2:{hello}] <-> tail` | — |
| `getMaxKey()` | unchanged | `tail.prev` = `[2:{hello}]` → `"hello"` |
| `getMinKey()` | unchanged | `head.next` = `[2:{hello}]` → `"hello"` |
| `inc("leet")` | new key → count-1 bucket inserted at front: `head <-> [1:{leet}] <-> [2:{hello}] <-> tail` | — |
| `getMaxKey()` | unchanged | `tail.prev` = `[2:{hello}]` → `"hello"` |
| `getMinKey()` | unchanged | `head.next` = `[1:{leet}]` → `"leet"` |

Output: `[null, null, null, "hello", "hello", null, "hello", "leet"]` ✔

---

## Key Takeaways

- **To make "find the extreme" O(1), keep the data sorted structurally instead of searching it.** A doubly-linked list of buckets ordered by count means the min is always the head and the max is always the tail.
- **Bucket by the changing quantity.** All keys with the same count share a node; a `±1` update just hops the key to the *adjacent* bucket, which is the crux of the O(1) bound. This "list of buckets" pattern also solves LFU Cache (#460).
- **Map → node pointer** is the recurring trick for O(1) structure edits (also LRU Cache #146): the map gives you the node instantly, the doubly-linked list lets you splice it in O(1).
- **Sentinel head/tail nodes** remove every boundary null-check — insertion, deletion, and the empty-list case all become uniform.
- Always **unlink emptied buckets immediately**, otherwise `head.next`/`tail.prev` stop pointing at real min/max buckets.

---

## Related Problems

- LeetCode #146 — LRU Cache (hash map + doubly-linked list, O(1) moves)
- LeetCode #460 — LFU Cache (buckets keyed by frequency — near-identical structure)
- LeetCode #895 — Maximum Frequency Stack (grouping by count in buckets)
- LeetCode #380 — Insert Delete GetRandom O(1) (map + array for O(1) ops)
- LeetCode #705 — Design HashSet (data-structure design fundamentals)
