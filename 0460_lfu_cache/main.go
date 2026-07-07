package main

import (
	"container/list"
	"fmt"
	"strings"
)

// lfuADT is the common interface both implementations satisfy so main() can
// drive the same official operation sequence through each of them.
type lfuADT interface {
	Get(key int) int
	Put(key, value int)
}

// ── Approach 1: Brute Force (Map + Linear Eviction Scan) ─────────────────────
//
// BruteForceLFU solves LFU Cache with one map of entries; on a full insert it
// linearly scans to find the least-frequently-used key, breaking ties by the
// least-recently-used (smallest tick) key.
//
// Intuition:
//
//	Store per key: value, frequency (access count), and a monotonically
//	increasing "tick" recording the last time it was touched (for LRU
//	tie-breaking). Get/Put just update those fields. The only hard part —
//	eviction — is done by brute force: scan every entry, pick the one with the
//	smallest frequency, and among those the smallest tick (oldest use). Correct
//	and easy to reason about, but eviction is O(n).
//
// Time:  Get O(1), Put O(1) amortised except eviction O(n) (the scan).
// Space: O(capacity) — one entry per stored key.
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

// NewBruteForceLFU builds an empty brute-force LFU of the given capacity.
func NewBruteForceLFU(capacity int) *BruteForceLFU {
	return &BruteForceLFU{
		capacity: capacity,
		data:     map[int]*bfEntry{},
	}
}

// touch stamps an entry as freshly used: bump its frequency and record the time.
func (c *BruteForceLFU) touch(e *bfEntry) {
	e.freq++         // one more access
	c.clock++        // advance logical time
	e.tick = c.clock // remember when this access happened (newest wins LRU)
}

// Get returns the value for key (and counts as an access) or -1 if absent.
func (c *BruteForceLFU) Get(key int) int {
	e, ok := c.data[key]
	if !ok {
		return -1 // miss
	}
	c.touch(e) // a successful get raises frequency & recency
	return e.value
}

// Put inserts or updates key, evicting the LFU (then LRU) entry when full.
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
	c.clock++                                                    // stamp the insertion time
	c.data[key] = &bfEntry{value: value, freq: 1, tick: c.clock} // new keys start at freq 1
}

// evict removes the entry with minimum freq, breaking ties by minimum tick.
func (c *BruteForceLFU) evict() {
	victimKey := 0 // key to remove
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
	delete(c.data, victimKey) // drop the chosen victim
}

// ── Approach 2: HashMap + Frequency Buckets of DLLs (Optimal) ────────────────
//
// OptimalLFU solves LFU Cache in O(1) per operation using a key→node map plus a
// map from each frequency to a doubly-linked list of the keys at that
// frequency (ordered most-recent at the front), and a running minFreq.
//
// Intuition:
//
//	Two moving parts must be O(1): find an entry, and pick the eviction victim.
//	  • key→node map gives O(1) lookup.
//	  • freq→DLL groups all keys sharing a frequency; within a bucket the list
//	    is ordered by recency (front = most recent), so the LRU victim in a
//	    bucket is the list's back.
//	  • minFreq tracks the smallest occupied frequency, so eviction removes the
//	    back of freqList[minFreq] in O(1).
//	On each access, move the node from its bucket freq f to bucket f+1 (front),
//	and if bucket minFreq just emptied and it was the one we bumped, minFreq++.
//	On insert-when-full, evict the back of the minFreq bucket, then add the new
//	key at frequency 1 and reset minFreq = 1.
//
// Time:  Get O(1), Put O(1).
// Space: O(capacity) — a node per key plus the bucket lists.
type node struct {
	key, value, freq int // the frequency this node currently lives at
}

type OptimalLFU struct {
	capacity int
	minFreq  int                   // smallest frequency currently present
	nodes    map[int]*list.Element // key → its element inside a freq list
	freqList map[int]*list.List    // frequency → DLL of *node (front = newest)
}

// NewOptimalLFU builds an empty O(1) LFU of the given capacity.
func NewOptimalLFU(capacity int) *OptimalLFU {
	return &OptimalLFU{
		capacity: capacity,
		nodes:    map[int]*list.Element{},
		freqList: map[int]*list.List{},
	}
}

// bump promotes elem from frequency f to f+1, moving it to the front of the
// higher bucket, and advances minFreq if the old bucket emptied at minFreq.
func (c *OptimalLFU) bump(elem *list.Element) {
	nd := elem.Value.(*node)
	f := nd.freq
	// Remove the node from its current frequency bucket.
	c.freqList[f].Remove(elem)
	// If that bucket is now empty and it was the minimum, the new minimum is f+1
	// (the node we are about to promote keeps frequencies contiguous upward).
	if c.freqList[f].Len() == 0 && c.minFreq == f {
		c.minFreq++
	}
	// Insert into the f+1 bucket at the front (most recently used).
	nd.freq = f + 1
	if c.freqList[nd.freq] == nil {
		c.freqList[nd.freq] = list.New()
	}
	c.nodes[nd.key] = c.freqList[nd.freq].PushFront(nd)
}

// Get returns the value for key (counting as an access) or -1 if absent.
func (c *OptimalLFU) Get(key int) int {
	elem, ok := c.nodes[key]
	if !ok {
		return -1 // miss
	}
	val := elem.Value.(*node).value
	c.bump(elem) // successful access raises this key's frequency
	return val
}

// Put inserts or updates key, evicting the LFU-then-LRU key when at capacity.
func (c *OptimalLFU) Put(key, value int) {
	if c.capacity == 0 {
		return // nothing can be stored
	}
	// Existing key: update value and bump frequency.
	if elem, ok := c.nodes[key]; ok {
		elem.Value.(*node).value = value
		c.bump(elem)
		return
	}
	// New key: evict first if full.
	if len(c.nodes) >= c.capacity {
		minList := c.freqList[c.minFreq]
		victim := minList.Back() // back of the min bucket = LFU + LRU
		delete(c.nodes, victim.Value.(*node).key)
		minList.Remove(victim)
	}
	// Insert the new key at frequency 1 and reset the running minimum to 1.
	nd := &node{key: key, value: value, freq: 1}
	if c.freqList[1] == nil {
		c.freqList[1] = list.New()
	}
	c.nodes[key] = c.freqList[1].PushFront(nd)
	c.minFreq = 1 // a brand-new freq-1 entry makes 1 the smallest frequency
}

// runExample drives the single official operation sequence through one
// implementation and returns the output list in LeetCode's format.
//
// Ops:  ["LFUCache","put","put","get","put","get","get","put","get","get","get"]
// Args: [[2],[1,1],[2,2],[1],[3,3],[2],[3],[4,4],[1],[3],[4]]
// Out:  [null, null, null, 1, null, -1, 3, null, -1, 3, 4]
func runExample(newCache func(int) lfuADT) string {
	c := newCache(2)        // LFUCache(2)          → null
	out := []string{"null"} // constructor yields no value
	c.Put(1, 1)             // put(1,1)             → null   cache={1:1}
	out = append(out, "null")
	c.Put(2, 2) // put(2,2)             → null   cache={1:1, 2:2}
	out = append(out, "null")
	out = append(out, fmt.Sprintf("%d", c.Get(1))) // get(1) → 1     (freq of 1 becomes 2)
	c.Put(3, 3)                                    // put(3,3): full → evict key 2 (LFU) cache={1:1, 3:3}
	out = append(out, "null")
	out = append(out, fmt.Sprintf("%d", c.Get(2))) // get(2) → -1    (2 was evicted)
	out = append(out, fmt.Sprintf("%d", c.Get(3))) // get(3) → 3     (freq of 3 becomes 2)
	c.Put(4, 4)                                    // put(4,4): full → tie freq(1)=freq(3)=2, evict LRU=1 cache={3:3, 4:4}
	out = append(out, "null")
	out = append(out, fmt.Sprintf("%d", c.Get(1))) // get(1) → -1    (1 was evicted)
	out = append(out, fmt.Sprintf("%d", c.Get(3))) // get(3) → 3
	out = append(out, fmt.Sprintf("%d", c.Get(4))) // get(4) → 4
	return "[" + strings.Join(out, ", ") + "]"
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Map + Linear Eviction Scan) ===")
	fmt.Println(runExample(func(cap int) lfuADT { return NewBruteForceLFU(cap) })) // [null, null, null, 1, null, -1, 3, null, -1, 3, 4]

	fmt.Println("=== Approach 2: HashMap + Frequency Buckets of DLLs (Optimal) ===")
	fmt.Println(runExample(func(cap int) lfuADT { return NewOptimalLFU(cap) })) // [null, null, null, 1, null, -1, 3, null, -1, 3, 4]
}
