package main

import (
	"fmt"
	"sort"
)

// LeetCode #287 — Find the Duplicate Number
//
// Given an array `nums` of n+1 integers where each integer is in the range
// [1, n] inclusive, there is exactly one repeated number. Return that number.
// You must NOT modify the array (Floyd approach) and use only O(1) extra space
// (Floyd approach) — but we show several approaches for learning.

// ── Approach 1: Brute Force (Nested Scan) ────────────────────────────────────
//
// bruteForce finds the duplicate by counting, for each element, how many later
// elements equal it.
//
// Intuition:
//
//	The most literal reading: a duplicate is a value that appears at two
//	different indices. Compare every pair; the value common to a matching
//	pair is the answer.
//
// Algorithm:
//  1. For each index i, scan j > i.
//  2. If nums[i] == nums[j], that value is the duplicate — return it.
//
// Time:  O(n^2) — all pairs.
// Space: O(1) — no extra structures; does not modify the array.
func bruteForce(nums []int) int {
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i] == nums[j] { // same value at two indices → duplicate
				return nums[i]
			}
		}
	}
	return -1 // unreachable given the problem guarantees a duplicate exists
}

// ── Approach 2: Hash Set ─────────────────────────────────────────────────────
//
// hashSet finds the duplicate by remembering values seen so far.
//
// Intuition:
//
//	Walk the array once; the first value you see for the second time is the
//	duplicate. A set gives O(1) membership tests.
//
// Algorithm:
//  1. Keep a set `seen`.
//  2. For each value: if already in `seen`, return it; else add it.
//
// Time:  O(n) — one pass, O(1) set operations.
// Space: O(n) — the set may hold up to n distinct values (violates the O(1)
//
//	follow-up, but is simple and fast).
func hashSet(nums []int) int {
	seen := make(map[int]struct{}, len(nums)) // set of already-observed values
	for _, v := range nums {
		if _, ok := seen[v]; ok { // second sighting → duplicate
			return v
		}
		seen[v] = struct{}{} // record first sighting
	}
	return -1
}

// ── Approach 3: Sort Then Adjacent Compare ───────────────────────────────────
//
// sortScan finds the duplicate by sorting a copy and looking for equal
// neighbours.
//
// Intuition:
//
//	After sorting, duplicate values sit next to each other. We sort a COPY so
//	the original array is not modified (respecting the constraint).
//
// Algorithm:
//  1. Copy nums, sort the copy.
//  2. Scan adjacent pairs; equal neighbours reveal the duplicate.
//
// Time:  O(n log n) — the sort dominates.
// Space: O(n) — the copy (sorting in place would modify the input, which is
//
//	disallowed).
func sortScan(nums []int) int {
	c := make([]int, len(nums))
	copy(c, nums) // copy so the original array is left untouched
	sort.Ints(c)  // duplicates become adjacent
	for i := 1; i < len(c); i++ {
		if c[i] == c[i-1] { // neighbours equal → duplicate
			return c[i]
		}
	}
	return -1
}

// ── Approach 4: Binary Search on Value Range (Pigeonhole) ────────────────────
//
// binarySearchCount finds the duplicate by binary-searching over the VALUE
// range [1, n] using a counting predicate.
//
// Intuition:
//
//	Pick a midpoint value `mid`. Count how many array elements are ≤ mid. If
//	there were no duplicate among [1..mid], that count would be exactly `mid`.
//	If count > mid, pigeonhole says the duplicate lies in [1..mid]; otherwise
//	it lies in [mid+1..n]. Binary-search the value that first makes the count
//	exceed the value.
//
// Algorithm:
//  1. lo = 1, hi = n (n = len-1).
//  2. mid = (lo+hi)/2; count = #{elements ≤ mid}.
//  3. If count > mid → answer is in [lo, mid], set hi = mid.
//     Else → answer is in [mid+1, hi], set lo = mid + 1.
//  4. Converge to lo == hi == the duplicate.
//
// Time:  O(n log n) — log n binary-search steps, each an O(n) count.
// Space: O(1) — no modification of the array, constant extra space.
func binarySearchCount(nums []int) int {
	lo, hi := 1, len(nums)-1 // value range [1, n]
	for lo < hi {
		mid := lo + (hi-lo)/2
		count := 0
		for _, v := range nums { // how many values are ≤ mid?
			if v <= mid {
				count++
			}
		}
		if count > mid { // too many small values → duplicate ≤ mid
			hi = mid
		} else { // duplicate is on the high side
			lo = mid + 1
		}
	}
	return lo // lo == hi is the duplicate
}

// ── Approach 5: Floyd's Cycle Detection (Optimal) ────────────────────────────
//
// floydCycle finds the duplicate by treating the array as a linked list and
// detecting the entrance of its cycle.
//
// Intuition:
//
//	Read nums as a function f(i) = nums[i], i.e. follow index → value → next
//	index. Starting from index 0, this walk must eventually revisit a value
//	because a value appears twice, and that repeated value is the node where
//	two "pointers" merge — the entrance of a cycle. Floyd's tortoise & hare
//	finds a meeting point inside the cycle; a second walk from the start
//	locks onto the cycle's entrance, which is exactly the duplicate value.
//
// Algorithm:
//  1. slow = nums[0], fast = nums[nums[0]] (advance 1 and 2 steps).
//  2. Advance until they meet inside the cycle.
//  3. Reset slow to 0; advance slow and fast one step each until equal.
//  4. That equal value is the duplicate.
//
// Time:  O(n) — linear number of pointer moves.
// Space: O(1) — two integer pointers; the array is never modified.
func floydCycle(nums []int) int {
	// Phase 1: find an intersection point inside the cycle.
	slow, fast := nums[0], nums[nums[0]]
	for slow != fast {
		slow = nums[slow]       // one step
		fast = nums[nums[fast]] // two steps
	}
	// Phase 2: find the cycle entrance = duplicate value.
	slow = 0
	for slow != fast {
		slow = nums[slow] // both now move one step at a time
		fast = nums[fast]
	}
	return slow // entrance node == duplicated value
}

func main() {
	// Example 1: Input [1,3,4,2,2] → Output 2
	// Example 2: Input [3,1,3,4,2] → Output 3
	// Example 3: Input [3,3,3,3,3] → Output 3
	ex1 := []int{1, 3, 4, 2, 2}
	ex2 := []int{3, 1, 3, 4, 2}
	ex3 := []int{3, 3, 3, 3, 3}

	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(bruteForce(ex1)) // expected 2
	fmt.Println(bruteForce(ex2)) // expected 3
	fmt.Println(bruteForce(ex3)) // expected 3

	fmt.Println("=== Approach 2: Hash Set ===")
	fmt.Println(hashSet(ex1)) // expected 2
	fmt.Println(hashSet(ex2)) // expected 3
	fmt.Println(hashSet(ex3)) // expected 3

	fmt.Println("=== Approach 3: Sort Then Adjacent Compare ===")
	fmt.Println(sortScan(ex1)) // expected 2
	fmt.Println(sortScan(ex2)) // expected 3
	fmt.Println(sortScan(ex3)) // expected 3

	fmt.Println("=== Approach 4: Binary Search on Value Range ===")
	fmt.Println(binarySearchCount(ex1)) // expected 2
	fmt.Println(binarySearchCount(ex2)) // expected 3
	fmt.Println(binarySearchCount(ex3)) // expected 3

	fmt.Println("=== Approach 5: Floyd's Cycle Detection (Optimal) ===")
	fmt.Println(floydCycle(ex1)) // expected 2
	fmt.Println(floydCycle(ex2)) // expected 3
	fmt.Println(floydCycle(ex3)) // expected 3
}
