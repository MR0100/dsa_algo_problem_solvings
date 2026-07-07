package main

import "fmt"

// file simulates the hidden file object that sits behind the read4 API.
type file struct {
	data string // full file contents
	pos  int    // read4's own file pointer (persists across read4 calls)
}

// read4 is the provided API: it reads up to 4 consecutive characters from the
// file into buf4 and returns the number of characters actually read.
func (f *file) read4(buf4 []byte) int {
	count := 0
	// copy until we have 4 chars or the file is exhausted
	for count < 4 && f.pos < len(f.data) {
		buf4[count] = f.data[f.pos]
		count++
		f.pos++
	}
	return count
}

// ── Approach 1: Leftover Queue (Brute Force) ─────────────────────────────────
//
// queueReader solves Read N Characters Given Read4 II by staging every
// character through a growable queue.
//
// Intuition:
//
//	The whole difficulty vs. #157 is that a read4 chunk can straddle two read()
//	calls: read(1) on "abc" consumes chunk "abc" but delivers only "a" — the
//	"bc" must survive until the next call. Easiest fix: dump every chunk into
//	a persistent queue, and let each read() pop from the queue's front. The
//	queue IS the memory between calls.
//
// Algorithm (per read call):
//  1. While queue holds fewer than n chars: cnt = read4(buf4);
//     if cnt == 0 break (EOF); else append buf4[:cnt] to the queue.
//  2. total = min(n, len(queue)).
//  3. Copy queue[:total] into buf and pop them off the queue.
//  4. Return total.
//
// Time:  O(n) per call — each character enters and leaves the queue once.
// Space: O(n) — the queue can hold up to n+3 characters after over-reading.
type queueReader struct {
	read4 func([]byte) int
	queue []byte // characters fetched from the file but not yet delivered
}

func (r *queueReader) read(buf []byte, n int) int {
	buf4 := make([]byte, 4) // scratch buffer for the API
	// fill the queue until it can satisfy the request or the file ends
	for len(r.queue) < n {
		cnt := r.read4(buf4)
		if cnt == 0 {
			break // EOF: nothing more to enqueue
		}
		r.queue = append(r.queue, buf4[:cnt]...) // stage every char, even extras
	}
	total := n
	if len(r.queue) < n {
		total = len(r.queue) // file ran out before n chars
	}
	copy(buf, r.queue[:total]) // deliver from the front of the queue
	r.queue = r.queue[total:]  // pop delivered chars; leftovers persist
	return total
}

// ── Approach 2: Persistent Buffer + Pointers (Optimal) ───────────────────────
//
// pointerReader solves Read N Characters Given Read4 II with a fixed 4-byte
// buffer and two indices that persist across calls.
//
// Intuition:
//
//	At most 3 characters can ever be "left over" (a chunk is ≤ 4 and at least
//	1 was consumed, or none were — either way the surplus fits in buf4). So
//	instead of a growable queue, keep the LAST read4 chunk plus two cursors:
//	i4 = next unconsumed index inside buf4, n4 = number of valid chars in
//	buf4. Drain buf4 first; refill only when i4 == n4.
//
// Algorithm (per read call):
//  1. total = 0.
//  2. While total < n:
//     a. If i4 == n4 (internal buffer drained): n4 = read4(buf4); i4 = 0;
//     if n4 == 0 → EOF, break.
//     b. While total < n and i4 < n4: buf[total++] = buf4[i4++].
//  3. Return total.
//
// Time:  O(n) per call — each character is copied exactly once overall.
// Space: O(1) — a fixed 4-byte buffer and two ints, regardless of n or file.
type pointerReader struct {
	read4 func([]byte) int
	buf4  [4]byte // last chunk fetched from the file (persists across calls)
	i4    int     // index of the next unconsumed char inside buf4
	n4    int     // number of valid chars currently in buf4
}

func (r *pointerReader) read(buf []byte, n int) int {
	total := 0 // chars delivered to buf in THIS call
	for total < n {
		if r.i4 == r.n4 { // internal buffer fully consumed → refill from file
			r.n4 = r.read4(r.buf4[:])
			r.i4 = 0
			if r.n4 == 0 {
				break // EOF: file exhausted and no leftovers remain
			}
		}
		// drain the internal buffer into buf without exceeding the request
		for total < n && r.i4 < r.n4 {
			buf[total] = r.buf4[r.i4]
			total++
			r.i4++
		}
	}
	return total
}

// reader is the common interface both approaches satisfy, for uniform testing.
type reader interface {
	read(buf []byte, n int) int
}

// runQueries executes a sequence of read(n) calls against one reader instance
// (the whole point: state must persist BETWEEN calls).
func runQueries(r reader, data string, queries []int) {
	fmt.Printf("file=%q queries=%v\n", data, queries)
	for _, n := range queries {
		buf := make([]byte, n) // destination sized for this request
		got := r.read(buf, n)
		fmt.Printf("  read(%d) → %d %q\n", n, got, string(buf[:got]))
	}
}

func main() {
	fmt.Println("=== Approach 1: Leftover Queue (Brute Force) ===")
	f1 := &file{data: "abc"}
	runQueries(&queueReader{read4: f1.read4}, "abc", []int{1, 2, 1}) // expected 1 "a", 2 "bc", 0 ""
	f2 := &file{data: "abc"}
	runQueries(&queueReader{read4: f2.read4}, "abc", []int{4, 1}) // expected 3 "abc", 0 ""

	fmt.Println("=== Approach 2: Persistent Buffer + Pointers (Optimal) ===")
	f3 := &file{data: "abc"}
	runQueries(&pointerReader{read4: f3.read4}, "abc", []int{1, 2, 1}) // expected 1 "a", 2 "bc", 0 ""
	f4 := &file{data: "abc"}
	runQueries(&pointerReader{read4: f4.read4}, "abc", []int{4, 1}) // expected 3 "abc", 0 ""
}
