package main

import "fmt"

// ── Approach 1: Queue with Full Recompute (Brute Force) ──────────────────────
//
// MovingAverageBrute keeps every value that is currently inside the sliding
// window in a slice and, on each next(), sums the whole window from scratch.
//
// Intuition:
//
//	A "moving average of the last size values" is just the average of a sliding
//	window. The most direct implementation stores the window explicitly: push
//	the new value, drop the oldest if we exceeded `size`, then add up whatever
//	remains and divide by the count. No cleverness — correct by construction.
//
// Algorithm:
//  1. Append val to the window slice.
//  2. If len(window) > size, remove the front element (the oldest).
//  3. Sum all remaining elements and divide by their count.
//
// Time:  O(size) per next() — the re-summation walks the whole window.
// Space: O(size) — the window slice holds at most `size` elements.
type MovingAverageBrute struct {
	size   int   // maximum number of values the window may hold
	window []int // the values currently inside the window, oldest at front
}

// NewMovingAverageBrute initialises the object with the window capacity.
func NewMovingAverageBrute(size int) *MovingAverageBrute {
	return &MovingAverageBrute{size: size}
}

// Next adds val to the stream and returns the average of the last `size` values.
func (m *MovingAverageBrute) Next(val int) float64 {
	m.window = append(m.window, val) // push newest value to the back
	if len(m.window) > m.size {      // window overflowed its capacity?
		m.window = m.window[1:] // evict the oldest value at the front
	}
	sum := 0
	for _, v := range m.window { // re-add everything currently in the window
		sum += v
	}
	// average = total / count; count is len(window), never zero after a push
	return float64(sum) / float64(len(m.window))
}

// ── Approach 2: Circular Buffer + Running Sum (Optimal) ──────────────────────
//
// MovingAverageCircular stores the window in a fixed-size ring buffer and keeps
// a running sum, so each next() is O(1).
//
// Intuition:
//
//	We do not need to re-sum the window every time. When a new value slides in,
//	exactly one old value (the one being overwritten) slides out. Maintain a
//	running sum: add the incoming value, subtract the outgoing one. A ring
//	buffer of length `size` gives O(1) access to the slot being overwritten and
//	never grows.
//
// Algorithm:
//  1. head = count % size selects the slot the new value will occupy.
//  2. sum += val − buf[head]  (the slot currently holds the value leaving the
//     window, which is 0 while the window is still filling).
//  3. buf[head] = val; count++.
//  4. Divide sum by min(count, size) — the true number of live values.
//
// Time:  O(1) per next() — constant index math and one add/subtract.
// Space: O(size) — one fixed ring buffer, allocated once.
type MovingAverageCircular struct {
	size  int   // window capacity and ring length
	buf   []int // ring buffer of the last `size` values
	count int   // total number of Next() calls so far
	sum   int   // running sum of the values currently in the window
}

// NewMovingAverageCircular initialises the ring buffer to the window capacity.
func NewMovingAverageCircular(size int) *MovingAverageCircular {
	return &MovingAverageCircular{size: size, buf: make([]int, size)}
}

// Next adds val to the stream and returns the average of the last `size` values.
func (m *MovingAverageCircular) Next(val int) float64 {
	head := m.count % m.size   // slot to (over)write; wraps around the ring
	m.sum += val - m.buf[head] // add newcomer, remove whatever it evicts
	m.buf[head] = val          // store the new value in its slot
	m.count++                  // one more value has entered the stream
	live := m.count            // how many real values are in the window
	if live > m.size {         // once past capacity, only `size` are live
		live = m.size
	}
	return float64(m.sum) / float64(live)
}

func main() {
	// Official Example 1:
	// Input:  ["MovingAverage","next","next","next","next"]
	//         [[3],[1],[10],[3],[5]]
	// Output: [null, 1.0, 5.5, 4.666666666666667, 6.0]

	fmt.Println("=== Approach 1: Queue with Full Recompute (Brute Force) ===")
	b := NewMovingAverageBrute(3)
	fmt.Println(b.Next(1))  // expected 1
	fmt.Println(b.Next(10)) // expected 5.5
	fmt.Println(b.Next(3))  // expected 4.666666666666667
	fmt.Println(b.Next(5))  // expected 6

	fmt.Println("=== Approach 2: Circular Buffer + Running Sum (Optimal) ===")
	c := NewMovingAverageCircular(3)
	fmt.Println(c.Next(1))  // expected 1
	fmt.Println(c.Next(10)) // expected 5.5
	fmt.Println(c.Next(3))  // expected 4.666666666666667
	fmt.Println(c.Next(5))  // expected 6
}
