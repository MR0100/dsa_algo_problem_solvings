package main

import (
	"container/heap"
	"fmt"
)

// ── Approach 1: Brute Force ──────────────────────────────────────────────────
//
// bruteForce solves Sliding Window Maximum by scanning each window of size k and
// taking its maximum directly.
//
// Intuition:
//
//	There are n-k+1 windows; for each, look at all k elements and keep the
//	largest. Simple and obviously correct, but re-examines overlapping elements
//	repeatedly.
//
// Algorithm:
//  1. For each start index i from 0 to n-k:
//  2. scan nums[i..i+k-1], track the maximum.
//  3. append it to the result.
//
// Time:  O(n·k) — every window rescans k elements.
// Space: O(1) extra (output excluded).
func bruteForce(nums []int, k int) []int {
	n := len(nums)
	if n == 0 || k == 0 {
		return []int{}
	}
	result := make([]int, 0, n-k+1)
	for i := 0; i+k <= n; i++ {
		maxVal := nums[i] // start with the window's first element
		for j := i + 1; j < i+k; j++ {
			if nums[j] > maxVal {
				maxVal = nums[j] // track the running window maximum
			}
		}
		result = append(result, maxVal)
	}
	return result
}

// ── Approach 2: Max-Heap (indices) ───────────────────────────────────────────
//
// maxHeap solves Sliding Window Maximum with a max-heap of (value, index)
// pairs, lazily discarding heap-top entries that have slid out of the window.
//
// Intuition:
//
//	The window max is the largest value among current indices. A max-heap gives
//	the largest instantly; we only need to skip entries whose index has fallen
//	out of the window (lazy deletion), which we detect at pop time.
//
// Algorithm:
//  1. Push each element (value, index) into a max-heap.
//  2. Once i >= k-1, pop while the top's index <= i-k (out of window), then read
//     the top value as the answer for this window.
//
// Time:  O(n log n) — up to n heap operations.
// Space: O(n) — the heap.
func maxHeap(nums []int, k int) []int {
	n := len(nums)
	if n == 0 || k == 0 {
		return []int{}
	}
	h := &idxHeap{}
	result := make([]int, 0, n-k+1)
	for i := 0; i < n; i++ {
		heap.Push(h, [2]int{nums[i], i}) // store value and its index
		if i >= k-1 {
			// Discard maxima whose index has left the window [i-k+1, i].
			for (*h)[0][1] <= i-k {
				heap.Pop(h)
			}
			result = append(result, (*h)[0][0]) // top value is the window max
		}
	}
	return result
}

// idxHeap is a max-heap of [value, index] pairs ordered by value.
type idxHeap [][2]int

func (h idxHeap) Len() int            { return len(h) }
func (h idxHeap) Less(i, j int) bool  { return h[i][0] > h[j][0] } // max-heap on value
func (h idxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *idxHeap) Push(x interface{}) { *h = append(*h, x.([2]int)) }
func (h *idxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}

// ── Approach 3: Monotonic Deque (Optimal) ────────────────────────────────────
//
// monotonicDeque solves Sliding Window Maximum in O(n) using a deque of indices
// kept in decreasing order of their values; the front is always the window max.
//
// Intuition:
//
//	A smaller element that appears BEFORE a larger element can never be the
//	window max while the larger one is in range — so we discard it. Maintaining
//	the deque so its values are strictly decreasing means the front index always
//	holds the current window's maximum. Each index is pushed and popped once.
//
// Algorithm:
//  1. For each i: pop from the BACK while nums[back] <= nums[i] (they're dominated).
//  2. Push i at the back.
//  3. Pop from the FRONT if it has slid out of the window (front <= i-k).
//  4. Once i >= k-1, record nums[front] as this window's maximum.
//
// Time:  O(n) — each index enters and leaves the deque at most once.
// Space: O(k) — the deque holds at most one window's worth of indices.
func monotonicDeque(nums []int, k int) []int {
	n := len(nums)
	if n == 0 || k == 0 {
		return []int{}
	}
	deque := make([]int, 0, k) // holds indices, values strictly decreasing
	result := make([]int, 0, n-k+1)
	for i := 0; i < n; i++ {
		// Drop indices at the back whose values can't beat nums[i]; nums[i]
		// dominates them for every future window.
		for len(deque) > 0 && nums[deque[len(deque)-1]] <= nums[i] {
			deque = deque[:len(deque)-1]
		}
		deque = append(deque, i) // nums[i] is a candidate max
		// Evict the front if it has fallen out of the window [i-k+1, i].
		if deque[0] <= i-k {
			deque = deque[1:]
		}
		if i >= k-1 {
			result = append(result, nums[deque[0]]) // front = current window max
		}
	}
	return result
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(bruteForce([]int{1, 3, -1, -3, 5, 3, 6, 7}, 3)) // expected [3 3 5 5 6 7]
	fmt.Println(bruteForce([]int{1}, 1))                        // expected [1]

	fmt.Println("=== Approach 2: Max-Heap (indices) ===")
	fmt.Println(maxHeap([]int{1, 3, -1, -3, 5, 3, 6, 7}, 3)) // expected [3 3 5 5 6 7]
	fmt.Println(maxHeap([]int{1}, 1))                        // expected [1]

	fmt.Println("=== Approach 3: Monotonic Deque (Optimal) ===")
	fmt.Println(monotonicDeque([]int{1, 3, -1, -3, 5, 3, 6, 7}, 3)) // expected [3 3 5 5 6 7]
	fmt.Println(monotonicDeque([]int{1}, 1))                        // expected [1]
}
