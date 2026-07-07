package main

import (
	"container/heap"
	"fmt"
)

// ── Approach 1: Min-Heap (Priority Queue of Candidates) ──────────────────────
//
// minHeap solves Super Ugly Number by repeatedly popping the smallest unseen
// multiple of a known super ugly number from a min-heap.
//
// Intuition:
//
//	Every super ugly number (other than 1) is prime * (an earlier super ugly
//	number). Seed the heap with each prime (== prime*1). Each time we pop the
//	global minimum, that value is the next super ugly number; push prime*value
//	for every prime to generate its successors. A `seen` set avoids emitting the
//	same value twice (e.g. 2*7 and 7*2 collide).
//
// Algorithm:
//
//  1. Start ugly=1 as the 1st super ugly number; push each prime into the heap.
//  2. Repeat n-1 times: pop the smallest value not seen before; that is the
//     next super ugly number. For each prime, push value*prime.
//  3. Return the n-th popped value.
//
// Time:  O(n*k*log(n*k)) — up to n*k pushes, each a log-time heap op (k = #primes).
// Space: O(n*k) — heap + seen set can hold that many candidates.
func minHeap(n int, primes []int) int {
	// pq is a min-heap of candidate super ugly numbers (int64 to avoid overflow
	// while generating; results fit in 32-bit int).
	pq := &int64Heap{}
	heap.Init(pq)
	seen := map[int64]bool{1: true}
	for _, p := range primes {
		heap.Push(pq, int64(p)) // prime*1 are the first successors of 1
		seen[int64(p)] = true
	}

	ugly := int64(1) // the 1st super ugly number is always 1
	for i := 1; i < n; i++ {
		ugly = heap.Pop(pq).(int64) // i-th smallest distinct super ugly number
		for _, p := range primes {
			cand := ugly * int64(p)
			if !seen[cand] { // dedupe: same product reachable multiple ways
				seen[cand] = true
				heap.Push(pq, cand)
			}
		}
	}
	return int(ugly)
}

// int64Heap is a minimal min-heap of int64 for the priority-queue approach.
type int64Heap []int64

func (h int64Heap) Len() int            { return len(h) }
func (h int64Heap) Less(i, j int) bool  { return h[i] < h[j] }
func (h int64Heap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *int64Heap) Push(x interface{}) { *h = append(*h, x.(int64)) }
func (h *int64Heap) Pop() interface{} {
	old := *h
	n := len(old)
	v := old[n-1]
	*h = old[:n-1]
	return v
}

// ── Approach 2: DP with k Pointers (Optimal Merge of k Sequences) ─────────────
//
// dpPointers solves Super Ugly Number by treating each prime as generating a
// sorted stream prime*ugly[0], prime*ugly[1], ... and merging the k streams,
// with one pointer per prime marking its current position in the ugly list.
//
// Intuition:
//
//	The full sorted list of super ugly numbers is the merge of k sorted
//	sequences: for prime p, the sequence p*ugly[0] < p*ugly[1] < ... Keep a
//	pointer idx[j] into `ugly` for each prime. The next super ugly number is
//	the minimum of primes[j]*ugly[idx[j]] across all j. Advance every pointer
//	whose product equals that minimum (this dedupes collisions naturally).
//
// Algorithm:
//
//  1. ugly[0] = 1; idx[j] = 0 for all primes.
//  2. For i from 1..n-1: next = min over j of primes[j]*ugly[idx[j]].
//  3. ugly[i] = next; for every j whose primes[j]*ugly[idx[j]] == next, idx[j]++.
//  4. Return ugly[n-1].
//
// Time:  O(n*k) — for each of n numbers, scan k primes twice (min then advance).
// Space: O(n + k) — the ugly array and the pointer array.
func dpPointers(n int, primes []int) int {
	k := len(primes)
	ugly := make([]int, n) // ugly[i] = (i+1)-th super ugly number
	ugly[0] = 1
	idx := make([]int, k) // idx[j] = position in ugly that prime j will multiply next

	for i := 1; i < n; i++ {
		next := int(^uint(0) >> 1) // start at max int
		// Find the minimum next candidate across all prime streams.
		for j := 0; j < k; j++ {
			if cand := primes[j] * ugly[idx[j]]; cand < next {
				next = cand
			}
		}
		ugly[i] = next
		// Advance every pointer that produced this minimum (handles duplicates).
		for j := 0; j < k; j++ {
			if primes[j]*ugly[idx[j]] == next {
				idx[j]++
			}
		}
	}
	return ugly[n-1]
}

func main() {
	// Official Example 1
	fmt.Println("=== Approach 1: Min-Heap ===")
	fmt.Println(minHeap(12, []int{2, 7, 13, 19})) // expected 32
	fmt.Println("=== Approach 2: DP with k Pointers (Optimal) ===")
	fmt.Println(dpPointers(12, []int{2, 7, 13, 19})) // expected 32

	// Official Example 2
	fmt.Println("=== Approach 1: Min-Heap (Example 2) ===")
	fmt.Println(minHeap(1, []int{2, 3, 5})) // expected 1
	fmt.Println("=== Approach 2: DP with k Pointers (Example 2) ===")
	fmt.Println(dpPointers(1, []int{2, 3, 5})) // expected 1
}
