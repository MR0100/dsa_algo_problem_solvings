package main

import (
	"container/heap"
	"fmt"
)

// LeetCode 264 — Ugly Number II.
//
// An UGLY NUMBER is a positive integer whose prime factors are limited to 2, 3
// and 5. The sequence of ugly numbers begins 1, 2, 3, 4, 5, 6, 8, 9, 10, 12,
// ... Given n, return the nth ugly number (1-indexed).

// ── Approach 1: Min-Heap (Dijkstra-style expansion) ──────────────────────────
//
// minHeap solves Ugly Number II by always popping the smallest ugly number so
// far and pushing its 2×, 3×, 5× multiples, de-duplicating with a set.
//
// Intuition:
//
//	Every ugly number times 2, 3 or 5 is also ugly. Starting from 1, repeatedly
//	extract the current minimum and generate its three multiples; the kth number
//	extracted is the kth ugly number. A seen-set prevents duplicates like
//	6 = 2·3 = 3·2 being processed twice.
//
// Algorithm:
//  1. Push 1 into a min-heap and mark it seen.
//  2. Repeat n times: pop the smallest value (this is the next ugly number);
//     for each factor f in {2,3,5}, push f·value if not already seen.
//  3. The value popped on the nth iteration is the answer.
//
// Time:  O(n log n) — n pops, each with up to 3 heap pushes of O(log n).
// Space: O(n) — the heap and the seen-set both hold O(n) values.
func minHeap(n int) int {
	h := &intHeap{1}              // start the sequence at 1
	seen := map[int]bool{1: true} // avoid pushing the same value twice
	var val int
	for i := 0; i < n; i++ {
		val = heap.Pop(h).(int) // the (i+1)-th smallest ugly number
		for _, f := range []int{2, 3, 5} {
			next := val * f
			if !seen[next] { // only queue unseen multiples
				seen[next] = true
				heap.Push(h, next)
			}
		}
	}
	return val // last popped value is the nth ugly number
}

// intHeap is a min-heap of ints for the heap approach.
type intHeap []int

func (h intHeap) Len() int            { return len(h) }
func (h intHeap) Less(i, j int) bool  { return h[i] < h[j] } // min-heap ordering
func (h intHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *intHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *intHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// ── Approach 2: Dynamic Programming (Three Pointers) (Optimal) ────────────────
//
// dpThreePointers solves Ugly Number II by building the sorted sequence in an
// array, where each new term is the smallest of (next multiple of 2/3/5) using
// one pointer per factor.
//
// Intuition:
//
//	The next ugly number is always some earlier ugly number multiplied by 2, 3
//	or 5. Keep three indices i2, i3, i5 pointing at the earliest earlier term
//	that hasn't yet been multiplied by that factor. The next term is
//	min(ugly[i2]·2, ugly[i3]·3, ugly[i5]·5); advance every pointer whose product
//	equals that min (advancing all matching pointers dedupes ties like 6).
//
// Algorithm:
//  1. ugly[0] = 1; pointers i2 = i3 = i5 = 0.
//  2. For k = 1..n-1: compute candidates a=ugly[i2]·2, b=ugly[i3]·3, c=ugly[i5]·5;
//     ugly[k] = min(a,b,c); advance each pointer whose candidate equals ugly[k].
//  3. Return ugly[n-1].
//
// Time:  O(n) — one array fill; each of the n terms costs O(1).
// Space: O(n) — the ugly array of size n.
func dpThreePointers(n int) int {
	ugly := make([]int, n) // ugly[k] = (k+1)-th ugly number, in sorted order
	ugly[0] = 1            // the 1st ugly number is 1
	i2, i3, i5 := 0, 0, 0  // next index to multiply by 2, 3, 5 respectively
	for k := 1; k < n; k++ {
		a := ugly[i2] * 2 // smallest unused multiple of 2
		b := ugly[i3] * 3 // smallest unused multiple of 3
		c := ugly[i5] * 5 // smallest unused multiple of 5
		next := a
		if b < next {
			next = b
		}
		if c < next {
			next = c
		}
		ugly[k] = next // the next ugly number is the smallest candidate
		// advance ALL pointers that produced this value (handles duplicates)
		if next == a {
			i2++
		}
		if next == b {
			i3++
		}
		if next == c {
			i5++
		}
	}
	return ugly[n-1] // nth ugly number (0-indexed n-1)
}

func main() {
	// Example 1: n = 10 ⇒ 12  (sequence 1,2,3,4,5,6,8,9,10,12).
	// Example 2: n = 1  ⇒ 1.

	fmt.Println("=== Approach 1: Min-Heap ===")
	fmt.Println(minHeap(10)) // expected 12
	fmt.Println(minHeap(1))  // expected 1

	fmt.Println("=== Approach 2: DP Three Pointers (Optimal) ===")
	fmt.Println(dpThreePointers(10)) // expected 12
	fmt.Println(dpThreePointers(1))  // expected 1
}
