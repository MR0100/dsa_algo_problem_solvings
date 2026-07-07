package main

import (
	"fmt"
	"math/rand"
)

// ListNode is the singly linked list node used by LeetCode.
type ListNode struct {
	Val  int
	Next *ListNode
}

// buildList builds a linked list from a slice and returns its head.
func buildList(vals []int) *ListNode {
	dummy := &ListNode{}
	cur := dummy
	for _, v := range vals {
		cur.Next = &ListNode{Val: v}
		cur = cur.Next
	}
	return dummy.Next
}

// ── Approach 1: Array Snapshot (Precompute) ──────────────────────────────────
//
// ArraySolution solves Linked List Random Node by copying every value into a
// slice once at construction, then indexing a uniformly random position on
// each getRandom call.
//
// Intuition:
//
//	A linked list has no O(1) random access, but a slice does. Pay one O(n)
//	pass up front to snapshot the values; afterwards each getRandom is a
//	single rand.Intn(n) index — trivially uniform because every index is
//	equally likely. The catch: it needs O(n) extra memory and assumes the
//	length is known/bounded, which the follow-up forbids.
//
// Algorithm:
//  1. Constructor: walk the list, appending each Val to vals.
//  2. getRandom: return vals[rand.Intn(len(vals))].
//
// Time:  Constructor O(n); getRandom O(1).
// Space: O(n) — the value snapshot.
type ArraySolution struct {
	vals []int      // snapshot of all node values
	rng  *rand.Rand // deterministic RNG for reproducible demos
}

// NewArraySolution snapshots the list into a slice.
func NewArraySolution(head *ListNode, rng *rand.Rand) *ArraySolution {
	vals := []int{}
	for n := head; n != nil; n = n.Next {
		vals = append(vals, n.Val) // record every value once
	}
	return &ArraySolution{vals: vals, rng: rng}
}

// GetRandom returns a uniformly random value via direct index.
func (s *ArraySolution) GetRandom() int {
	return s.vals[s.rng.Intn(len(s.vals))] // every index equally likely
}

// ── Approach 2: Reservoir Sampling (Optimal, streaming) ──────────────────────
//
// ReservoirSolution solves Linked List Random Node with reservoir sampling of
// size 1, so it works on a stream of unknown length and uses O(1) extra space.
//
// Intuition:
//
//	Walk the list once. Keep a single "chosen" value. When we meet the i-th
//	node (1-indexed), replace chosen with its value with probability 1/i.
//	By induction, after seeing k nodes every one of them is the current
//	choice with probability exactly 1/k — so at the end each of the n nodes
//	has probability 1/n. No length, no extra array: this is the answer to the
//	follow-up ("what if the list is huge and length unknown?").
//
// Algorithm:
//  1. chosen = head.Val, i = 1.
//  2. For each subsequent node, i++; with probability 1/i set chosen = node.Val
//     (implemented as rng.Intn(i) == 0).
//  3. Return chosen.
//
// Time:  Constructor O(1) (just store head); getRandom O(n) — one pass.
// Space: O(1) — a single running candidate.
type ReservoirSolution struct {
	head *ListNode
	rng  *rand.Rand
}

// NewReservoirSolution just remembers the head; no preprocessing.
func NewReservoirSolution(head *ListNode, rng *rand.Rand) *ReservoirSolution {
	return &ReservoirSolution{head: head, rng: rng}
}

// GetRandom returns a uniformly random value using reservoir sampling.
func (s *ReservoirSolution) GetRandom() int {
	chosen := s.head.Val // seed the reservoir with the first value
	i := 1               // number of nodes seen so far
	for n := s.head.Next; n != nil; n = n.Next {
		i++                     // now considering the i-th node
		if s.rng.Intn(i) == 0 { // pick it with probability 1/i
			chosen = n.Val
		}
	}
	return chosen
}

// sample runs getRandom `trials` times and returns a value→count histogram,
// used to demonstrate the (approximately) uniform distribution.
func sample(get func() int, trials int) map[int]int {
	hist := map[int]int{}
	for t := 0; t < trials; t++ {
		hist[get()]++
	}
	return hist
}

func main() {
	// Official example: list [1,2,3]; getRandom must return 1, 2, or 3
	// each with probability ~1/3.
	head := buildList([]int{1, 2, 3})

	fmt.Println("=== Approach 1: Array Snapshot ===")
	arr := NewArraySolution(head, rand.New(rand.NewSource(1)))
	fmt.Println(arr.GetRandom() >= 1 && arr.GetRandom() <= 3) // expected true (value in {1,2,3})
	fmt.Println(sample(arr.GetRandom, 30000))                 // expected ~10000 each: map[1:~ 2:~ 3:~]

	fmt.Println("=== Approach 2: Reservoir Sampling (Optimal) ===")
	res := NewReservoirSolution(head, rand.New(rand.NewSource(1)))
	fmt.Println(res.GetRandom() >= 1 && res.GetRandom() <= 3) // expected true (value in {1,2,3})
	fmt.Println(sample(res.GetRandom, 30000))                 // expected ~10000 each: map[1:~ 2:~ 3:~]
}
