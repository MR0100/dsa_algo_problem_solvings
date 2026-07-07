package main

import "fmt"

// ── Approach 1: Queue of Timestamps ──────────────────────────────────────────
//
// HitCounterQueue solves Design Hit Counter by storing every hit timestamp in a
// FIFO queue and evicting entries older than 300 seconds on each query.
//
// Intuition:
//
//	A hit "counts" only during the 300-second window ending at the current
//	timestamp. Keep every hit in arrival order (a queue). When asked for the
//	count at time t, drop all timestamps ≤ t-300 from the front — they have
//	expired — and the queue length is the answer.
//
// Time:  Hit  O(1) amortised — one append.
//
//	GetHits O(k) amortised where k = hits that just expired (each hit is
//	enqueued once and dequeued once over the object's lifetime).
//
// Space: O(n) — one entry per hit currently within (or about to leave) the
//
//	window; can grow large under a burst of hits in one second.
type HitCounterQueue struct {
	q []int // timestamps in non-decreasing (arrival) order
}

// NewHitCounterQueue constructs an empty queue-based counter.
func NewHitCounterQueue() *HitCounterQueue {
	return &HitCounterQueue{q: []int{}}
}

// Hit records a hit that happened at `timestamp` (seconds). Timestamps are
// non-decreasing across calls.
func (h *HitCounterQueue) Hit(timestamp int) {
	h.q = append(h.q, timestamp) // append keeps the queue sorted by time
}

// GetHits returns the number of hits in the past 300 seconds (i.e. with
// timestamp in (timestamp-300, timestamp]).
func (h *HitCounterQueue) GetHits(timestamp int) int {
	// Evict from the front everything that fell out of the 5-minute window.
	for len(h.q) > 0 && h.q[0] <= timestamp-300 {
		h.q = h.q[1:] // pop the oldest expired timestamp
	}
	return len(h.q) // survivors are exactly the hits within the window
}

// ── Approach 2: Fixed 300-Slot Circular Buffer (Optimal) ─────────────────────
//
// HitCounterBuckets solves Design Hit Counter with two fixed arrays of size 300
// indexed by timestamp mod 300: one holds a per-second hit count, the other the
// last timestamp that wrote that slot (to detect stale data).
//
// Intuition:
//
//	The window is exactly 300 seconds, so second t and second t+300 share the
//	same bucket t%300. Store, per bucket, both a count and the timestamp that
//	count belongs to. On a hit, if the bucket's stored timestamp differs from
//	now, the old value is stale (from a previous 300-cycle) — reset it. On a
//	query, sum only the buckets whose stored timestamp is still inside the
//	window (> t-300).
//
// Time:  Hit  O(1) — index one bucket.
//
//	GetHits O(300) = O(1) — scan the fixed 300 buckets.
//
// Space: O(300) = O(1) — two constant-size arrays regardless of hit volume.
type HitCounterBuckets struct {
	times  [300]int // times[i] = the timestamp that counts[i] refers to
	counts [300]int // counts[i] = number of hits during second times[i]
}

// NewHitCounterBuckets constructs an empty bucketed counter.
func NewHitCounterBuckets() *HitCounterBuckets {
	return &HitCounterBuckets{}
}

// Hit records a hit at `timestamp`.
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

// GetHits returns the number of hits in the past 300 seconds.
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

func main() {
	// Official example operation sequence:
	// HitCounter(); hit(1); hit(2); hit(3); getHits(4)=3;
	// hit(300); getHits(300)=4; getHits(301)=3.

	fmt.Println("=== Approach 1: Queue of Timestamps ===")
	q := NewHitCounterQueue()
	q.Hit(1)
	q.Hit(2)
	q.Hit(3)
	fmt.Println(q.GetHits(4)) // expected 3
	q.Hit(300)
	fmt.Println(q.GetHits(300)) // expected 4
	fmt.Println(q.GetHits(301)) // expected 3 (the hit at second 1 expired)

	fmt.Println("=== Approach 2: Fixed 300-Slot Circular Buffer (Optimal) ===")
	b := NewHitCounterBuckets()
	b.Hit(1)
	b.Hit(2)
	b.Hit(3)
	fmt.Println(b.GetHits(4)) // expected 3
	b.Hit(300)
	fmt.Println(b.GetHits(300)) // expected 4
	fmt.Println(b.GetHits(301)) // expected 3
}
