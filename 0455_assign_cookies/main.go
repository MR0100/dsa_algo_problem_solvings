package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Greedy Two Pointers, Smallest-Fit (Optimal) ──────────────────
//
// greedyTwoPointers solves Assign Cookies by sorting both greed factors and
// cookie sizes ascending, then feeding each child the smallest cookie that
// still satisfies them.
//
// Intuition:
//
//	To content the MOST children, never "waste" a big cookie on a child a
//	small cookie would already satisfy. Sort children by greed and cookies by
//	size, both ascending. Sweep both with two pointers: for the current
//	(least greedy unserved) child, advance through cookies until one is big
//	enough — give it to this child (count++, move to next child), then move to
//	the next cookie. If a cookie is too small for the current child it is too
//	small for every remaining (greedier) child, so discard it. This exchange
//	argument guarantees optimality: satisfying the least-greedy child with the
//	smallest adequate cookie never blocks a better overall assignment.
//
// Algorithm:
//  1. Sort g (greed) and s (sizes) ascending.
//  2. child = 0, cookie = 0, count = 0.
//  3. While child < len(g) and cookie < len(s):
//     - if s[cookie] >= g[child]: this cookie contents the child →
//     count++, child++, cookie++.
//     - else: cookie too small → cookie++ (skip it).
//  4. Return count.
//
// Time:  O(n log n + m log m) — sorting both arrays dominates; the sweep is
//
//	O(n + m).
//
// Space: O(1) extra beyond sorting (in place); O(log) recursion stack.
func greedyTwoPointers(g []int, s []int) int {
	sort.Ints(g) // children by ascending greed (easiest to please first)
	sort.Ints(s) // cookies by ascending size (spend the smallest first)

	child := 0  // index into g: the current least-greedy unserved child
	cookie := 0 // index into s: the current smallest unused cookie
	count := 0  // contented children so far
	for child < len(g) && cookie < len(s) {
		if s[cookie] >= g[child] {
			// This cookie satisfies the current child: assign it and advance
			// both — this child is done, this cookie is spent.
			count++
			child++
			cookie++
		} else {
			// Cookie too small for the least greedy remaining child, hence too
			// small for everyone left; drop it and try the next larger cookie.
			cookie++
		}
	}
	return count
}

// ── Approach 2: Greedy From the Largest (Biggest-Fit) ────────────────────────
//
// greedyLargestFirst solves the same problem from the other end: sort both
// ascending but walk from the back, matching the greediest child with the
// largest cookie.
//
// Intuition:
//
//	Symmetric greedy. Consider the greediest child first and hand them the
//	largest cookie. If the largest cookie can content the greediest child,
//	that is a match (consume both). If it cannot, no cookie can content this
//	child (it is the greediest and that was the biggest cookie), so this child
//	goes unserved — move to the next-greediest child while keeping the cookie.
//	Same optimal count, just approached top-down; some find this direction's
//	"if the biggest can't please the greediest, give up on that child" easier
//	to reason about.
//
// Algorithm:
//  1. Sort g and s ascending.
//  2. child = len(g)-1, cookie = len(s)-1, count = 0.
//  3. While child >= 0 and cookie >= 0:
//     - if s[cookie] >= g[child]: match → count++, child--, cookie--.
//     - else: this child unservable → child-- (keep the cookie).
//  4. Return count.
//
// Time:  O(n log n + m log m) — sorting dominates.
// Space: O(1) extra + O(log) stack.
func greedyLargestFirst(g []int, s []int) int {
	sort.Ints(g)
	sort.Ints(s)

	child := len(g) - 1  // greediest child
	cookie := len(s) - 1 // largest cookie
	count := 0
	for child >= 0 && cookie >= 0 {
		if s[cookie] >= g[child] {
			// Largest remaining cookie contents the greediest remaining child.
			count++
			child--
			cookie--
		} else {
			// Even the biggest cookie can't satisfy this (greediest) child, so
			// nothing can; this child stays unhappy, keep the cookie for a less
			// greedy child.
			child--
		}
	}
	return count
}

func main() {
	fmt.Println("=== Approach 1: Greedy Two Pointers, Smallest-Fit (Optimal) ===")
	fmt.Printf("g=[1,2,3], s=[1,1] -> %d  expected 1\n", greedyTwoPointers([]int{1, 2, 3}, []int{1, 1}))
	fmt.Printf("g=[1,2],   s=[1,2,3] -> %d  expected 2\n", greedyTwoPointers([]int{1, 2}, []int{1, 2, 3}))

	fmt.Println("=== Approach 2: Greedy From the Largest (Biggest-Fit) ===")
	fmt.Printf("g=[1,2,3], s=[1,1] -> %d  expected 1\n", greedyLargestFirst([]int{1, 2, 3}, []int{1, 1}))
	fmt.Printf("g=[1,2],   s=[1,2,3] -> %d  expected 2\n", greedyLargestFirst([]int{1, 2}, []int{1, 2, 3}))

	fmt.Println("=== Edge cases ===")
	fmt.Printf("g=[1,2,3], s=[] (no cookies)  -> %d  expected 0\n", greedyTwoPointers([]int{1, 2, 3}, []int{}))
	fmt.Printf("g=[10,9,8,7], s=[5,6,7,8] -> %d  expected 2\n", greedyTwoPointers([]int{10, 9, 8, 7}, []int{5, 6, 7, 8}))
	fmt.Printf("g=[1,1,1], s=[1,1,1]      -> %d  expected 3\n", greedyLargestFirst([]int{1, 1, 1}, []int{1, 1, 1}))
}
