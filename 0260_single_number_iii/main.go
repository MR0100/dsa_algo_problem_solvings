package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Hash Map Counting (Brute Force) ──────────────────────────────
//
// hashMap solves Single Number III by counting occurrences and collecting the
// values that appear exactly once.
//
// Intuition:
//
//	Every number appears twice except two that appear once. Tally counts in a
//	map, then return the two keys with count 1.
//
// Algorithm:
//  1. Count each value's frequency in a map.
//  2. Collect keys whose count is 1 (there are exactly two).
//  3. Return them.
//
// Time:  O(n) — one pass to count, one over the map.
// Space: O(n) — the count map.
func hashMap(nums []int) []int {
	count := make(map[int]int)
	for _, x := range nums { // tally frequencies
		count[x]++
	}
	res := make([]int, 0, 2)
	for x, c := range count {
		if c == 1 { // appears exactly once
			res = append(res, x)
		}
	}
	sort.Ints(res) // stable output for testing (map order is random)
	return res
}

// ── Approach 2: XOR Partition by Differing Bit (Optimal) ─────────────────────
//
// xorPartition solves Single Number III in O(n) time and O(1) space using XOR.
//
// Intuition:
//
//	XOR of all numbers cancels every pair, leaving xorAll = a ^ b, where a and b
//	are the two uniques. Since a != b, xorAll has at least one set bit — pick
//	its lowest set bit (xorAll & -xorAll). That bit differs between a and b, so
//	it splits ALL numbers into two groups: one containing a (and pairs), the
//	other containing b (and pairs). XOR each group separately: pairs cancel,
//	leaving a in one group and b in the other.
//
// Algorithm:
//  1. xorAll = XOR of all nums (= a ^ b).
//  2. diff = xorAll & -xorAll (lowest set bit where a and b differ).
//  3. For each x: if x&diff != 0, XOR into groupA; else XOR into groupB.
//  4. Return {groupA, groupB}.
//
// Time:  O(n) — two linear passes.
// Space: O(1) — a few integer accumulators.
func xorPartition(nums []int) []int {
	xorAll := 0
	for _, x := range nums { // pairs cancel; left with a ^ b
		xorAll ^= x
	}
	// Isolate the lowest bit where a and b differ (two's-complement trick).
	diff := xorAll & (-xorAll)
	a, b := 0, 0
	for _, x := range nums {
		if x&diff != 0 { // this number has the differing bit set
			a ^= x // group A: pairs cancel, leaves one unique
		} else {
			b ^= x // group B: leaves the other unique
		}
	}
	res := []int{a, b}
	sort.Ints(res) // stable output for testing
	return res
}

func main() {
	fmt.Println("=== Approach 1: Hash Map Counting ===")
	fmt.Println(hashMap([]int{1, 2, 1, 3, 2, 5})) // expected [3 5]
	fmt.Println(hashMap([]int{-1, 0}))            // expected [-1 0]
	fmt.Println(hashMap([]int{0, 1}))             // expected [0 1]

	fmt.Println("=== Approach 2: XOR Partition by Differing Bit (Optimal) ===")
	fmt.Println(xorPartition([]int{1, 2, 1, 3, 2, 5})) // expected [3 5]
	fmt.Println(xorPartition([]int{-1, 0}))            // expected [-1 0]
	fmt.Println(xorPartition([]int{0, 1}))             // expected [0 1]
}
