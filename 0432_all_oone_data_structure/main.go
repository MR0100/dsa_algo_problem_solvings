package main

import "fmt"

// allOneADT is the common interface both implementations satisfy so main() can
// drive the same official operation sequence through each and compare output.
type allOneADT interface {
	Inc(key string)    // insert key with value 1, or increment its value
	Dec(key string)    // decrement key; remove it when its value hits 0
	GetMaxKey() string // any key with the maximum value ("" if empty)
	GetMinKey() string // any key with the minimum value ("" if empty)
}

// ── Approach 1: Hash Map + Linear Scan (Brute Force) ─────────────────────────
//
// BruteForceAllOne stores every key's count in a plain map and answers the
// min/max queries by scanning the whole map.
//
// Intuition:
//
//	The simplest structure that satisfies the API is "key → count". Inc/Dec are
//	O(1) map edits. The only hard part — finding a key with the min or max count
//	— is solved by the dumbest possible method: look at every entry and keep the
//	best. Correct and trivial, but GetMaxKey/GetMinKey are O(n) each, which is
//	exactly what the follow-up ("all operations O(1)") asks us to beat.
//
// Algorithm:
//
//	Inc:       counts[key]++.
//	Dec:       counts[key]--; delete the entry if it reaches 0.
//	GetMaxKey: scan counts, return a key achieving the largest value.
//	GetMinKey: scan counts, return a key achieving the smallest value.
//
// Time:  Inc O(1), Dec O(1), GetMaxKey/GetMinKey O(n).
// Space: O(n) — one map entry per distinct live key.
type BruteForceAllOne struct {
	counts map[string]int // key → current positive count
}

// NewBruteForceAllOne builds the empty brute-force structure.
func NewBruteForceAllOne() *BruteForceAllOne {
	return &BruteForceAllOne{counts: map[string]int{}}
}

// Inc inserts key at 1 or increments its existing count.
func (a *BruteForceAllOne) Inc(key string) {
	a.counts[key]++ // missing key defaults to 0, so this makes a new key 1
}

// Dec decrements key and removes it once its count drops to 0.
func (a *BruteForceAllOne) Dec(key string) {
	if _, ok := a.counts[key]; !ok {
		return // key absent → do nothing (per spec)
	}
	a.counts[key]-- // one occurrence removed
	if a.counts[key] == 0 {
		delete(a.counts, key) // value 0 means the key no longer exists
	}
}

// GetMaxKey scans for any key holding the maximum count.
func (a *BruteForceAllOne) GetMaxKey() string {
	best, bestVal := "", -1 // sentinel below any real count
	for k, v := range a.counts {
		if v > bestVal { // strictly larger → new champion
			best, bestVal = k, v
		}
	}
	return best // "" when the map is empty
}

// GetMinKey scans for any key holding the minimum count.
func (a *BruteForceAllOne) GetMinKey() string {
	best, bestVal := "", int(^uint(0)>>1) // start at max int
	for k, v := range a.counts {
		if v < bestVal { // strictly smaller → new champion
			best, bestVal = k, v
		}
	}
	return best // "" when the map is empty
}

// ── Approach 2: Doubly-Linked List of Count Buckets + Hash Map (Optimal) ──────
//
// AllOne achieves O(1) for every operation by grouping keys that share the same
// count into "buckets", and keeping the buckets on a doubly-linked list sorted
// by count in ascending order.
//
// Intuition:
//
//	The costly part of the brute force is finding the min/max count. Fix that by
//	keeping counts *sorted structurally*: a doubly-linked list of buckets, each
//	bucket = a set of keys with one specific count, ordered ascending. Then the
//	min bucket is always the head (after a sentinel) and the max bucket is always
//	the tail (before a sentinel) — O(1) to read. A hash map key→bucket lets Inc
//	and Dec find a key's bucket in O(1). Inc moves the key to the neighbouring
//	bucket for count+1 (creating it if absent); Dec moves it to count−1 (or drops
//	it). Because counts change by ±1, the destination bucket is always adjacent,
//	so the whole move is O(1). Empty buckets are unlinked immediately, keeping
//	head/tail meaningful.
//
// Algorithm (with sentinels head/tail so no null-checks at the ends):
//
//	Inc(key):
//	  if key unseen: its target count is 1; ensure a bucket right after head has
//	    count 1, add key there.
//	  else: let cur = keyBucket[key] with count c; ensure the next bucket has
//	    count c+1 (insert one if needed), move key there, remove key from cur.
//	  drop cur if it became empty.
//	Dec(key):
//	  let cur = keyBucket[key] with count c.
//	  if c == 1: remove key entirely (map + bucket).
//	  else: ensure the previous bucket has count c−1 (insert if needed), move key
//	    there.
//	  drop cur if it became empty.
//	GetMaxKey: any key in tail.prev's set (or "" if list empty).
//	GetMinKey: any key in head.next's set (or "" if list empty).
//
// Time:  Inc O(1), Dec O(1), GetMaxKey O(1), GetMinKey O(1) — amortised/worst
//
//	case, since every step is a constant number of pointer splices and one map op.
//
// Space: O(n) — n live keys spread across at most n buckets.
type bucket struct {
	count int                 // the shared count of every key in this bucket
	keys  map[string]struct{} // set of keys currently holding `count`
	prev  *bucket             // neighbour with a smaller count
	next  *bucket             // neighbour with a larger count
}

// AllOne is the optimal structure: a hash map to each key's bucket, plus a
// sentinel-bounded doubly-linked list of buckets sorted ascending by count.
type AllOne struct {
	keyBucket map[string]*bucket // key → the bucket that currently holds it
	head      *bucket            // sentinel BEFORE the smallest-count bucket
	tail      *bucket            // sentinel AFTER the largest-count bucket
}

// NewAllOne wires up the two sentinels into an empty list.
func NewAllOne() *AllOne {
	head := &bucket{keys: map[string]struct{}{}} // left sentinel (count unused)
	tail := &bucket{keys: map[string]struct{}{}} // right sentinel (count unused)
	head.next = tail                             // empty list: head <-> tail
	tail.prev = head
	return &AllOne{
		keyBucket: map[string]*bucket{},
		head:      head,
		tail:      tail,
	}
}

// insertAfter splices a brand-new bucket with the given count immediately after
// node `prev`, and returns it. O(1) pointer surgery.
func (a *AllOne) insertAfter(prev *bucket, count int) *bucket {
	b := &bucket{count: count, keys: map[string]struct{}{}}
	b.prev = prev      // link back
	b.next = prev.next // link forward to whatever followed prev
	prev.next.prev = b // old successor now points back to b
	prev.next = b      // prev now points to b
	return b
}

// remove unlinks an (empty) bucket from the list. O(1).
func (a *AllOne) remove(b *bucket) {
	b.prev.next = b.next // bypass b going forward
	b.next.prev = b.prev // bypass b going backward
}

// Inc inserts key at count 1 or bumps it to count+1, moving it one bucket right.
func (a *AllOne) Inc(key string) {
	if cur, ok := a.keyBucket[key]; ok {
		// Existing key at count c → move to count c+1 (the next bucket).
		next := cur.next
		if next == a.tail || next.count != cur.count+1 {
			// No adjacent bucket for c+1 yet → create one right after cur.
			next = a.insertAfter(cur, cur.count+1)
		}
		next.keys[key] = struct{}{} // key now lives at count+1
		a.keyBucket[key] = next     // update its home pointer
		delete(cur.keys, key)       // leave the old bucket
		if len(cur.keys) == 0 {
			a.remove(cur) // no keys left at count c → unlink the empty bucket
		}
	} else {
		// Brand-new key → target count is 1, which belongs right after head.
		first := a.head.next
		if first == a.tail || first.count != 1 {
			// No count-1 bucket at the front yet → create one.
			first = a.insertAfter(a.head, 1)
		}
		first.keys[key] = struct{}{} // register the key at count 1
		a.keyBucket[key] = first
	}
}

// Dec decrements key: removes it entirely at count 1, else moves it one bucket
// left to count−1.
func (a *AllOne) Dec(key string) {
	cur, ok := a.keyBucket[key]
	if !ok {
		return // key absent → no-op
	}
	if cur.count == 1 {
		// Dropping from 1 means the key disappears completely.
		delete(cur.keys, key)
		delete(a.keyBucket, key)
	} else {
		// Move to count-1 (the previous bucket), creating it if missing.
		prev := cur.prev
		if prev == a.head || prev.count != cur.count-1 {
			prev = a.insertAfter(cur.prev, cur.count-1) // insert BEFORE cur
		}
		prev.keys[key] = struct{}{} // key now at count-1
		a.keyBucket[key] = prev
		delete(cur.keys, key) // leave the old bucket
	}
	if len(cur.keys) == 0 {
		a.remove(cur) // clean up any bucket we emptied
	}
}

// GetMaxKey returns any key in the highest-count bucket (tail.prev).
func (a *AllOne) GetMaxKey() string {
	if a.tail.prev == a.head {
		return "" // list empty
	}
	for k := range a.tail.prev.keys { // any element of the max bucket
		return k
	}
	return ""
}

// GetMinKey returns any key in the lowest-count bucket (head.next).
func (a *AllOne) GetMinKey() string {
	if a.head.next == a.tail {
		return "" // list empty
	}
	for k := range a.head.next.keys { // any element of the min bucket
		return k
	}
	return ""
}

// runExample drives the official LeetCode operation sequence through one
// implementation and returns the output list in LeetCode's null-padded format.
//
// Ops:  ["AllOne","inc","inc","getMaxKey","getMinKey","inc","getMaxKey","getMinKey"]
// Args: [[],["hello"],["hello"],[],[],["leet"],[],[]]
// Out:  [null,null,null,"hello","hello",null,"hello","leet"]
func runExample(newAllOne func() allOneADT) []string {
	a := newAllOne()
	out := []string{"null"} // constructor returns nothing
	a.Inc("hello")
	out = append(out, "null")
	a.Inc("hello")
	out = append(out, "null")
	out = append(out, quote(a.GetMaxKey())) // "hello" (count 2)
	out = append(out, quote(a.GetMinKey())) // "hello" (only key)
	a.Inc("leet")
	out = append(out, "null")
	out = append(out, quote(a.GetMaxKey())) // "hello" (count 2 > 1)
	out = append(out, quote(a.GetMinKey())) // "leet"  (count 1 < 2)
	return out
}

// quote wraps a key in double quotes to match LeetCode's printed output format.
func quote(s string) string { return "\"" + s + "\"" }

func main() {
	fmt.Println("=== Approach 1: Hash Map + Linear Scan (Brute Force) ===")
	fmt.Println(runExample(func() allOneADT { return NewBruteForceAllOne() }))
	// [null null null "hello" "hello" null "hello" "leet"]

	fmt.Println("=== Approach 2: DLL of Buckets + Hash Map (Optimal) ===")
	fmt.Println(runExample(func() allOneADT { return NewAllOne() }))
	// [null null null "hello" "hello" null "hello" "leet"]
}
