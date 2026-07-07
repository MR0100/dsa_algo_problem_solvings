package main

import (
	"container/list"
	"fmt"
	"strings"
)

// lruCache is the common interface every approach implements so that main()
// can drive the same official example through all of them.
type lruCache interface {
	Get(key int) int
	Put(key int, value int)
}

// ── Approach 1: Brute Force (Ordered Slice) ──────────────────────────────────
//
// LRUCacheSlice solves LRU Cache using a plain slice kept in recency order.
//
// Intuition:
//
//	The simplest possible model of "least recently used" is a list ordered by
//	recency: index 0 is the least recently used entry, the last index is the
//	most recently used. Every Get/Put linearly searches the slice, and any
//	touched entry is moved to the back. Eviction is just dropping index 0.
//
// Algorithm:
//
//	Get:  linear-scan for key → if found, move entry to the back, return value.
//	Put:  linear-scan for key → if found, update value and move to back;
//	      otherwise append; if len > capacity, drop the front (LRU) entry.
//
// Time:  O(capacity) per operation — linear scan plus slice shifting.
// Space: O(capacity) — one pair stored per cached key.
type LRUCacheSlice struct {
	capacity int
	entries  []kvPair // entries[0] = least recently used, entries[len-1] = most recently used
}

// kvPair is one cached key/value entry for the slice-based approach.
type kvPair struct {
	key, value int
}

// NewLRUCacheSlice builds the brute-force cache with the given capacity.
func NewLRUCacheSlice(capacity int) *LRUCacheSlice {
	return &LRUCacheSlice{capacity: capacity, entries: make([]kvPair, 0, capacity)}
}

// Get returns the value for key or -1, promoting the entry to most-recent.
func (c *LRUCacheSlice) Get(key int) int {
	for i, e := range c.entries {
		if e.key == key {
			// remove the entry from its current slot...
			c.entries = append(c.entries[:i], c.entries[i+1:]...)
			// ...and re-append it at the back, marking it most recently used
			c.entries = append(c.entries, e)
			return e.value
		}
	}
	return -1 // key not present
}

// Put inserts/updates key and evicts the front (LRU) entry on overflow.
func (c *LRUCacheSlice) Put(key int, value int) {
	for i, e := range c.entries {
		if e.key == key {
			// key already cached: drop the old slot and re-append with new value
			c.entries = append(c.entries[:i], c.entries[i+1:]...)
			c.entries = append(c.entries, kvPair{key, value})
			return
		}
	}
	// brand-new key: append as most recently used
	c.entries = append(c.entries, kvPair{key, value})
	if len(c.entries) > c.capacity {
		// over capacity → evict entries[0], the least recently used
		c.entries = c.entries[1:]
	}
}

// ── Approach 2: Hash Map + Doubly Linked List (Optimal) ──────────────────────
//
// LRUCacheDLL solves LRU Cache with a hash map into a hand-rolled doubly
// linked list, achieving O(1) for both Get and Put.
//
// Intuition:
//
//	We need two things in O(1): (a) find a node by key → hash map, and
//	(b) move a node to the "most recent" end / evict the "least recent" end
//	→ doubly linked list, because unlinking a node needs prev AND next
//	pointers. Sentinel head/tail nodes remove all nil-edge special cases.
//
// Algorithm:
//
//	Layout: head <-> (LRU) ... (MRU) <-> tail, map[key]*node.
//	Get:  look up node in map; if present, unlink it and re-insert before
//	      tail (mark most recent); return its value.
//	Put:  if key exists, update value and move node before tail.
//	      Otherwise create a node, insert before tail, add to map;
//	      if size > capacity, unlink head.next (the LRU) and delete its key.
//
// Time:  O(1) per operation — map lookup plus constant pointer surgery.
// Space: O(capacity) — one map entry + one list node per cached key.
type LRUCacheDLL struct {
	capacity int
	index    map[int]*dllNode // key → node, for O(1) lookup
	head     *dllNode         // sentinel: head.next is the LRU node
	tail     *dllNode         // sentinel: tail.prev is the MRU node
}

// dllNode is one doubly-linked-list node holding a cached key/value.
type dllNode struct {
	key, value int
	prev, next *dllNode
}

// NewLRUCacheDLL builds the optimal cache with sentinel head/tail wired up.
func NewLRUCacheDLL(capacity int) *LRUCacheDLL {
	head, tail := &dllNode{}, &dllNode{}
	head.next = tail // empty list: head <-> tail
	tail.prev = head
	return &LRUCacheDLL{capacity: capacity, index: make(map[int]*dllNode, capacity), head: head, tail: tail}
}

// unlink removes n from the list in O(1) using its prev/next pointers.
func (c *LRUCacheDLL) unlink(n *dllNode) {
	n.prev.next = n.next // bypass n going forward
	n.next.prev = n.prev // bypass n going backward
}

// pushBack inserts n right before the tail sentinel (most-recent position).
func (c *LRUCacheDLL) pushBack(n *dllNode) {
	n.prev = c.tail.prev // n sits after the current MRU
	n.next = c.tail
	c.tail.prev.next = n
	c.tail.prev = n
}

// Get returns the value for key or -1, promoting the node to most-recent.
func (c *LRUCacheDLL) Get(key int) int {
	n, ok := c.index[key]
	if !ok {
		return -1 // not cached
	}
	c.unlink(n)   // pull the node out of its current position
	c.pushBack(n) // re-insert as most recently used
	return n.value
}

// Put inserts/updates key in O(1) and evicts the LRU node on overflow.
func (c *LRUCacheDLL) Put(key int, value int) {
	if n, ok := c.index[key]; ok {
		// key already cached: overwrite value and promote to most-recent
		n.value = value
		c.unlink(n)
		c.pushBack(n)
		return
	}
	// new key: create node, mark most-recent, register in the map
	n := &dllNode{key: key, value: value}
	c.pushBack(n)
	c.index[key] = n
	if len(c.index) > c.capacity {
		lru := c.head.next // head.next is always the least recently used
		c.unlink(lru)
		delete(c.index, lru.key) // node stores its key exactly for this delete
	}
}

// ── Approach 3: Hash Map + container/list (Optimal, Stdlib) ─────────────────
//
// LRUCacheStdList solves LRU Cache with Go's container/list doing the
// doubly-linked-list bookkeeping.
//
// Intuition:
//
//	Identical design to Approach 2, but the standard library already ships a
//	doubly linked list with O(1) MoveToBack/Remove. In production Go code this
//	is the idiomatic version; in an interview, Approach 2 shows you can build
//	the pointer machinery yourself.
//
// Algorithm:
//
//	Get:  map lookup → MoveToBack(elem) → return stored value.
//	Put:  existing key → update payload, MoveToBack.
//	      New key → PushBack, store element in map;
//	      on overflow → Remove(list.Front()) and delete its key.
//
// Time:  O(1) per operation.
// Space: O(capacity).
type LRUCacheStdList struct {
	capacity int
	index    map[int]*list.Element // key → element, for O(1) lookup
	order    *list.List            // Front() = LRU, Back() = MRU
}

// NewLRUCacheStdList builds the stdlib-backed cache.
func NewLRUCacheStdList(capacity int) *LRUCacheStdList {
	return &LRUCacheStdList{capacity: capacity, index: make(map[int]*list.Element, capacity), order: list.New()}
}

// Get returns the value for key or -1, moving the element to the back.
func (c *LRUCacheStdList) Get(key int) int {
	elem, ok := c.index[key]
	if !ok {
		return -1 // not cached
	}
	c.order.MoveToBack(elem) // promote to most recently used
	return elem.Value.(kvPair).value
}

// Put inserts/updates key and evicts the front element on overflow.
func (c *LRUCacheStdList) Put(key int, value int) {
	if elem, ok := c.index[key]; ok {
		elem.Value = kvPair{key, value} // overwrite payload
		c.order.MoveToBack(elem)        // promote to most recently used
		return
	}
	c.index[key] = c.order.PushBack(kvPair{key, value}) // new MRU element
	if c.order.Len() > c.capacity {
		front := c.order.Front() // Front() is the least recently used
		c.order.Remove(front)
		delete(c.index, front.Value.(kvPair).key)
	}
}

// runExample drives the official LeetCode operation sequence through one
// cache implementation and returns the outputs formatted LeetCode-style.
func runExample(build func(capacity int) lruCache) string {
	ops := []string{"LRUCache", "put", "put", "get", "put", "get", "put", "get", "get", "get"}
	args := [][]int{{2}, {1, 1}, {2, 2}, {1}, {3, 3}, {2}, {4, 4}, {1}, {3}, {4}}

	var cache lruCache
	out := make([]string, 0, len(ops))
	for i, op := range ops {
		switch op {
		case "LRUCache":
			cache = build(args[i][0]) // constructor returns nothing → null
			out = append(out, "null")
		case "put":
			cache.Put(args[i][0], args[i][1]) // put returns nothing → null
			out = append(out, "null")
		case "get":
			out = append(out, fmt.Sprint(cache.Get(args[i][0])))
		}
	}
	return "[" + strings.Join(out, ", ") + "]"
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Ordered Slice) ===")
	fmt.Println(runExample(func(c int) lruCache { return NewLRUCacheSlice(c) })) // [null, null, null, 1, null, -1, null, -1, 3, 4]

	fmt.Println("=== Approach 2: Hash Map + Doubly Linked List (Optimal) ===")
	fmt.Println(runExample(func(c int) lruCache { return NewLRUCacheDLL(c) })) // [null, null, null, 1, null, -1, null, -1, 3, 4]

	fmt.Println("=== Approach 3: Hash Map + container/list (Optimal, Stdlib) ===")
	fmt.Println(runExample(func(c int) lruCache { return NewLRUCacheStdList(c) })) // [null, null, null, 1, null, -1, null, -1, 3, 4]
}
