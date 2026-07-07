package main

import (
	"fmt"
	"sort"
)

// people[i] = [hi, ki] means person i has height hi and exactly ki people in
// FRONT of them (earlier in the queue) whose height is >= hi. Given the people
// in arbitrary order, rebuild the queue that satisfies every (h, k) pair.

// ── Approach 1: Brute Force (Slot Placement) ─────────────────────────────────
//
// bruteForce solves Queue Reconstruction by Height by scanning, for each person
// taken shortest-first, the still-empty slots and dropping them into the slot
// that has exactly k taller-or-equal people ahead.
//
// Intuition:
//
//	Process people from shortest to tallest. When we place a short person, the
//	only people already fixed are the ones we placed earlier — all of them are
//	SHORTER than the current person (or equal, handled by tie-breaking), so they
//	do NOT count toward k. The empty slots we leave behind will be filled by
//	TALLER people later, and every such taller person DOES count toward k. So the
//	current person must sit in the empty slot with exactly k empty slots before
//	it — those k slots are guaranteed to be filled by taller people.
//
// Algorithm:
//  1. Sort ascending by height; break height ties by DESCENDING k so that among
//     equal-height people the one needing more people ahead is placed first
//     (equal heights count toward each other).
//  2. result is a slice of n empty slots (marked with a sentinel).
//  3. For each person, walk result left→right counting empty slots; when the
//     count of empty slots seen equals k, drop the person into that slot.
//
// Time:  O(n^2) — for each of n people we may scan all n slots.
// Space: O(n) — the sentinel-filled result plus the sort.
func bruteForce(people [][]int) [][]int {
	n := len(people)
	// Sort ascending by height; ties broken by larger k first because two people
	// of the same height each count toward the other's k.
	sort.Slice(people, func(i, j int) bool {
		if people[i][0] == people[j][0] {
			return people[i][1] > people[j][1] // taller-or-equal tie → bigger k first
		}
		return people[i][0] < people[j][0] // shorter first
	})

	result := make([][]int, n) // fixed-size queue; nil entries = empty slots
	for _, p := range people {
		empties := 0 // how many empty slots we have passed so far
		for idx := 0; idx < n; idx++ {
			if result[idx] != nil {
				continue // occupied by a shorter person already placed — skip
			}
			// This slot is empty. If exactly k empties precede it, p belongs here:
			// its k taller people will land in those k earlier empty slots.
			if empties == p[1] {
				result[idx] = []int{p[0], p[1]} // claim the slot
				break
			}
			empties++ // count this empty slot and keep scanning
		}
	}
	return result
}

// ── Approach 2: Greedy Insertion (Tallest First) — Optimal-idea ───────────────
//
// greedyInsert solves Queue Reconstruction by Height by sorting tallest-first
// and inserting each person at index k of the result list.
//
// Intuition:
//
//	Flip the order: place the TALLEST people first. Once the tall people are
//	arranged correctly among themselves, inserting a SHORTER person cannot
//	disturb any already-placed person's k — a shorter person is invisible to a
//	taller one's count (k only counts height >= mine). And for the person being
//	inserted, EVERYONE already in the list is taller-or-equal, so its k value is
//	literally the index it must occupy: put it at position k and exactly k
//	taller-or-equal people sit in front of it.
//
// Algorithm:
//  1. Sort by DESCENDING height; break ties by ASCENDING k (so among equal
//     heights the smaller-k person is inserted first and ends up further front).
//  2. For each person in that order, insert it into result at index = its k.
//
// Time:  O(n^2) — each of n slice insertions shifts up to n elements.
// Space: O(n) — the result slice.
func greedyInsert(people [][]int) [][]int {
	// Tallest first; equal heights → smaller k first (smaller k must sit earlier).
	sort.Slice(people, func(i, j int) bool {
		if people[i][0] == people[j][0] {
			return people[i][1] < people[j][1] // equal height → smaller k first
		}
		return people[i][0] > people[j][0] // taller first
	})

	result := make([][]int, 0, len(people))
	for _, p := range people {
		k := p[1] // everyone already placed is >= p's height, so k IS the target index
		// Insert p at position k: grow by one, shift the tail right, drop p in.
		result = append(result, nil)   // extend length by one
		copy(result[k+1:], result[k:]) // shift elements from k rightward
		result[k] = []int{p[0], p[1]}  // place p exactly at index k
	}
	return result
}

// equal reports whether two queues are identical (helper for the demo output).
func equal(a, b [][]int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i][0] != b[i][0] || a[i][1] != b[i][1] {
			return false
		}
	}
	return true
}

func main() {
	ex1 := [][]int{{7, 0}, {4, 4}, {7, 1}, {5, 0}, {6, 1}, {5, 2}}
	want1 := [][]int{{5, 0}, {7, 0}, {5, 2}, {6, 1}, {4, 4}, {7, 1}}

	ex2 := [][]int{{6, 0}, {5, 0}, {4, 0}, {3, 2}, {2, 2}, {1, 4}}
	want2 := [][]int{{4, 0}, {5, 0}, {2, 2}, {3, 2}, {1, 4}, {6, 0}}

	fmt.Println("=== Approach 1: Brute Force (Slot Placement) ===")
	// bruteForce mutates the input order, so pass a fresh copy each time.
	in1a := [][]int{{7, 0}, {4, 4}, {7, 1}, {5, 0}, {6, 1}, {5, 2}}
	got1a := bruteForce(in1a)
	fmt.Printf("ex1 got=%v  match=%v  (expected %v)\n", got1a, equal(got1a, want1), want1)
	in2a := [][]int{{6, 0}, {5, 0}, {4, 0}, {3, 2}, {2, 2}, {1, 4}}
	got2a := bruteForce(in2a)
	fmt.Printf("ex2 got=%v  match=%v  (expected %v)\n", got2a, equal(got2a, want2), want2)

	fmt.Println("=== Approach 2: Greedy Insertion (Tallest First) ===")
	in1b := [][]int{{7, 0}, {4, 4}, {7, 1}, {5, 0}, {6, 1}, {5, 2}}
	got1b := greedyInsert(in1b)
	fmt.Printf("ex1 got=%v  match=%v  (expected %v)\n", got1b, equal(got1b, want1), want1)
	in2b := [][]int{{6, 0}, {5, 0}, {4, 0}, {3, 2}, {2, 2}, {1, 4}}
	got2b := greedyInsert(in2b)
	fmt.Printf("ex2 got=%v  match=%v  (expected %v)\n", got2b, equal(got2b, want2), want2)

	_ = ex1
	_ = ex2
}
