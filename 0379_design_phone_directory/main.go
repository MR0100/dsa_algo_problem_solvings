package main

import (
	"fmt"
	"strings"
)

// phoneDirectoryADT is the common interface every implementation satisfies so
// main() can drive the same official operation sequence through all of them.
type phoneDirectoryADT interface {
	Get() int // hand out any free number, or -1 if none
	Check(number int) bool
	Release(number int)
}

// ── Approach 1: Boolean Array + Linear Scan (Brute Force) ────────────────────
//
// LinearDirectory implements the directory with a used[] flag array and finds a
// free slot by scanning from the start on every Get.
//
// Intuition:
//
//	The simplest model: a boolean per number saying "in use". Get scans left to
//	right for the first false; Check reads the flag; Release clears it. Correct
//	and minimal, but Get is O(maxNumbers) because it re-scans every call.
//
// Algorithm:
//
//	Get:     scan used[0..max) for the first free index, mark it, return it.
//	Check:   return !used[number].
//	Release: used[number] = false.
//
// Time:  Get O(maxNumbers), Check O(1), Release O(1).
// Space: O(maxNumbers) — the flag array.
type LinearDirectory struct {
	used []bool // used[i] == true if number i is currently assigned
}

// NewLinearDirectory builds a directory over [0, maxNumbers).
func NewLinearDirectory(maxNumbers int) *LinearDirectory {
	return &LinearDirectory{used: make([]bool, maxNumbers)}
}

// Get returns the first free number, or -1 when the directory is full.
func (d *LinearDirectory) Get() int {
	for i := range d.used {
		if !d.used[i] {
			d.used[i] = true // claim this slot
			return i
		}
	}
	return -1 // nothing free
}

// Check reports whether `number` is currently available (not assigned).
func (d *LinearDirectory) Check(number int) bool {
	return !d.used[number]
}

// Release returns `number` to the pool of free slots.
func (d *LinearDirectory) Release(number int) {
	d.used[number] = false
}

// ── Approach 2: Queue of Free Slots + Used Set (Optimal) ─────────────────────
//
// QueueDirectory implements the directory with a FIFO queue of currently-free
// numbers plus a boolean used[] for O(1) Check, giving O(1) on all operations.
//
// Intuition:
//
//	The linear scan wastes time re-finding free slots. Instead keep an explicit
//	pool of free numbers in a queue: Get pops the front (O(1)); Release pushes a
//	number back (guarding against double-release with the used[] flags so the
//	same number is never queued twice).
//
// Algorithm:
//
//	Get:     if queue empty return -1; else pop front, mark used, return it.
//	Check:   return !used[number].
//	Release: if used[number], clear it and push number onto the queue.
//
// Time:  Get O(1), Check O(1), Release O(1) (amortized; slice pop of a growing
//
//	front index).
//
// Space: O(maxNumbers) — the queue plus the flag array.
type QueueDirectory struct {
	free []int  // FIFO pool of available numbers (front at index `head`)
	head int    // index of the current queue front (avoids reslicing churn)
	used []bool // used[i] == true when number i is assigned
}

// NewQueueDirectory builds a directory with every number [0, maxNumbers) free.
func NewQueueDirectory(maxNumbers int) *QueueDirectory {
	free := make([]int, maxNumbers)
	for i := range free {
		free[i] = i // initially every number is available, in order
	}
	return &QueueDirectory{free: free, head: 0, used: make([]bool, maxNumbers)}
}

// Get pops a free number in O(1), or returns -1 if none remain.
func (d *QueueDirectory) Get() int {
	if d.head >= len(d.free) {
		return -1 // queue exhausted — directory full
	}
	number := d.free[d.head] // take the front
	d.head++                 // advance the queue front
	d.used[number] = true    // mark assigned
	return number
}

// Check reports availability in O(1).
func (d *QueueDirectory) Check(number int) bool {
	return !d.used[number]
}

// Release frees `number`, guarding against releasing an already-free slot.
func (d *QueueDirectory) Release(number int) {
	if !d.used[number] {
		return // not assigned — do nothing, avoids duplicate queue entries
	}
	d.used[number] = false          // mark free
	d.free = append(d.free, number) // re-enqueue for future Gets
}

// runExample drives the single official LeetCode example through one
// implementation and returns the output list in LeetCode's format.
//
// Ops:  ["PhoneDirectory","get","get","check","get","check","release","check"]
// Args: [[3],[],[],[2],[],[2],[2],[2]]
func runExample(newDir func() phoneDirectoryADT) string {
	d := newDir()                                    // "PhoneDirectory"(3) → null
	out := []string{"null"}                          // constructor returns no value
	out = append(out, fmt.Sprintf("%d", d.Get()))    // get() → 0
	out = append(out, fmt.Sprintf("%d", d.Get()))    // get() → 1
	out = append(out, fmt.Sprintf("%t", d.Check(2))) // check(2) → true
	out = append(out, fmt.Sprintf("%d", d.Get()))    // get() → 2
	out = append(out, fmt.Sprintf("%t", d.Check(2))) // check(2) → false
	d.Release(2)                                     // release(2) → null
	out = append(out, "null")
	out = append(out, fmt.Sprintf("%t", d.Check(2))) // check(2) → true
	return "[" + strings.Join(out, ", ") + "]"
}

func main() {
	fmt.Println("=== Approach 1: Boolean Array + Linear Scan ===")
	fmt.Println(runExample(func() phoneDirectoryADT { return NewLinearDirectory(3) })) // [null, 0, 1, true, 2, false, null, true]

	fmt.Println("=== Approach 2: Queue of Free Slots + Used Set (Optimal) ===")
	fmt.Println(runExample(func() phoneDirectoryADT { return NewQueueDirectory(3) })) // [null, 0, 1, true, 2, false, null, true]
}
