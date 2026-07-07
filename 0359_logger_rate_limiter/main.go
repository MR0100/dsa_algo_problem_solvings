package main

import (
	"container/list"
	"fmt"
)

// loggerADT is the common interface every approach implements so main() can
// drive the same official call sequence through all of them.
type loggerADT interface {
	ShouldPrintMessage(timestamp int, message string) bool
}

// ── Approach 1: Hash Map of Next-Allowed Time (Optimal) ──────────────────────
//
// HashMapLogger solves Logger Rate Limiter by remembering, for each message,
// the earliest timestamp at which it may be printed again.
//
// Intuition:
//
//	The rule: a message prints only if it hasn't been printed in the last 10
//	seconds. Equivalently, once printed at time t, it is blocked until t+10.
//	Store per message the value t+10 = "next allowed time". A new request at
//	`timestamp` prints iff timestamp >= nextAllowed[message]; if it prints,
//	update nextAllowed[message] = timestamp + 10. Unknown messages default to
//	next-allowed 0, so they always print the first time.
//
// Algorithm:
//  1. Look up nextAllowed[message] (0 if absent).
//  2. If timestamp < nextAllowed ⇒ return false (still cooling down).
//  3. Else set nextAllowed[message] = timestamp + 10 and return true.
//
// Time:  O(1) average per call — one map read + one map write.
// Space: O(m) — one entry per distinct message ever printed (never shrinks).
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

// ── Approach 2: Queue + Set (Sliding Window Eviction) ────────────────────────
//
// QueueSetLogger solves Logger Rate Limiter by maintaining a sliding 10-second
// window of recently printed messages, evicting expired ones from the front.
//
// Intuition:
//
//	Keep a FIFO queue of (timestamp, message) for everything printed in the
//	last 10 seconds, plus a set of the messages currently in that window. On a
//	new request at `timestamp`, first evict from the front every entry with
//	time <= timestamp - 10 (outside the window), removing it from the set. Then
//	the message may print iff it is NOT in the set. If it prints, enqueue it and
//	add it to the set. This models the "last 10 seconds" window literally.
//
// Algorithm:
//  1. Pop front entries whose timestamp <= timestamp - 10, deleting from set.
//  2. If message is in set ⇒ return false.
//  3. Else enqueue (timestamp, message), add to set, return true.
//
// Time:  O(1) amortized — each entry is enqueued and dequeued at most once.
// Space: O(w) — entries within the 10-second window (bounded by distinct
//
//	messages printed in that span).
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

func main() {
	// Official example call sequence:
	//   Logger();                       // constructor
	//   shouldPrintMessage(1, "foo");   // true
	//   shouldPrintMessage(2, "bar");   // true
	//   shouldPrintMessage(3, "foo");   // false (foo at 1, 3-1 < 10)
	//   shouldPrintMessage(8, "bar");   // false (bar at 2, 8-2 < 10)
	//   shouldPrintMessage(10, "foo");  // false (foo at 1, 10-1 < 10)
	//   shouldPrintMessage(11, "foo");  // true  (foo at 1, 11-1 == 10)
	calls := []struct {
		ts  int
		msg string
		exp bool
	}{
		{1, "foo", true},
		{2, "bar", true},
		{3, "foo", false},
		{8, "bar", false},
		{10, "foo", false},
		{11, "foo", true},
	}

	run := func(name string, l loggerADT) {
		fmt.Printf("=== %s ===\n", name)
		for _, c := range calls {
			got := l.ShouldPrintMessage(c.ts, c.msg)
			fmt.Printf("shouldPrintMessage(%2d, %q) got=%-5t expected %t\n", c.ts, c.msg, got, c.exp)
		}
	}

	run("Approach 1: Hash Map (Optimal)", NewHashMapLogger())
	run("Approach 2: Queue + Set (Sliding Window)", NewQueueSetLogger())
}
