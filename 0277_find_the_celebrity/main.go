package main

import "fmt"

// The judge exposes a boolean primitive:
//
//	knows(a, b) == true  ⇔  person a knows person b.
//
// A celebrity is someone whom EVERYBODY knows but who knows NOBODY. There is at
// most one celebrity. We model the acquaintance graph as an n×n matrix and wrap
// it in a knows() helper so the solutions read exactly like the LeetCode API.

var graph [][]bool // graph[a][b] == true means a knows b

// knows is the provided API: does person a know person b?
func knows(a, b int) bool {
	return graph[a][b] // O(1) matrix lookup
}

// ── Approach 1: Brute Force (Count Degrees) ──────────────────────────────────
//
// bruteForce solves Find the Celebrity by directly checking the definition for
// every candidate.
//
// Intuition:
//
//	Person c is the celebrity iff for every other person i: knows(i, c) is true
//	(everyone knows c) and knows(c, i) is false (c knows no one). Test every c.
//
// Algorithm:
//  1. For each candidate c in 0..n-1:
//  2. For each other i: if knows(c, i) OR NOT knows(i, c), c fails.
//  3. If c passed every i, return c.
//  4. If nobody qualifies, return -1.
//
// Time:  O(n^2) — n candidates × n checks, each with O(1) knows() calls.
// Space: O(1) — no extra structures.
func bruteForce(n int) int {
	for c := 0; c < n; c++ { // try every person as the celebrity
		isCelebrity := true
		for i := 0; i < n; i++ {
			if i == c {
				continue // skip self
			}
			// c must NOT know i, and i MUST know c
			if knows(c, i) || !knows(i, c) {
				isCelebrity = false
				break // this candidate is disqualified
			}
		}
		if isCelebrity {
			return c
		}
	}
	return -1 // no celebrity present
}

// ── Approach 2: Two-Pass Elimination (Optimal) ───────────────────────────────
//
// twoPointers solves Find the Celebrity by first shrinking to a single
// candidate, then verifying it.
//
// Intuition:
//
//	Each knows(a, b) call eliminates exactly one person from contention:
//	  - if knows(a, b) is true, a knows someone, so a can't be the celebrity → a out.
//	  - if knows(a, b) is false, nobody-knows-b means b can't be the celebrity → b out.
//	Sweep a single "candidate" pointer across everyone: after n-1 calls only one
//	person can still possibly be the celebrity. Then verify that survivor fully,
//	because "possible" is not "confirmed".
//
// Algorithm:
//  1. candidate = 0.
//  2. For i = 1..n-1: if knows(candidate, i) then candidate = i (old one is out).
//  3. Verify candidate against everyone: it must know no one and be known by all.
//  4. Return candidate if verified, else -1.
//
// Time:  O(n) — n-1 calls to narrow + up to 2n calls to verify = O(n).
// Space: O(1) — a single pointer.
func twoPointers(n int) int {
	candidate := 0 // start by assuming person 0 might be the celebrity
	// Phase 1: narrow down to one candidate.
	for i := 1; i < n; i++ {
		if knows(candidate, i) {
			// candidate knows someone → candidate is disqualified; i survives
			candidate = i
		}
		// else i is disqualified (nobody-knows-i via candidate); keep candidate
	}
	// Phase 2: verify the survivor is really the celebrity.
	for i := 0; i < n; i++ {
		if i == candidate {
			continue
		}
		// celebrity knows no one AND is known by everyone
		if knows(candidate, i) || !knows(i, candidate) {
			return -1 // survivor failed verification → no celebrity
		}
	}
	return candidate
}

// buildGraph is a test helper: turns a matrix of ints into the bool graph.
func buildGraph(m [][]int) {
	graph = make([][]bool, len(m))
	for i := range m {
		graph[i] = make([]bool, len(m[i]))
		for j := range m[i] {
			graph[i][j] = m[i][j] == 1
		}
	}
}

func main() {
	// Example 1: graph = [[1,1,0],[0,1,0],[1,1,1]]. Person 1 knows nobody else
	// and is known by 0 and 2 → celebrity is 1.
	fmt.Println("=== Approach 1: Brute Force ===")
	buildGraph([][]int{{1, 1, 0}, {0, 1, 0}, {1, 1, 1}})
	fmt.Println(bruteForce(3)) // expected 1

	// Example 2: graph = [[1,0,1],[1,1,0],[0,1,1]]. No one is known by all &
	// knows none → -1.
	buildGraph([][]int{{1, 0, 1}, {1, 1, 0}, {0, 1, 1}})
	fmt.Println(bruteForce(3)) // expected -1

	fmt.Println("=== Approach 2: Two-Pass Elimination (Optimal) ===")
	buildGraph([][]int{{1, 1, 0}, {0, 1, 0}, {1, 1, 1}})
	fmt.Println(twoPointers(3)) // expected 1

	buildGraph([][]int{{1, 0, 1}, {1, 1, 0}, {0, 1, 1}})
	fmt.Println(twoPointers(3)) // expected -1
}
