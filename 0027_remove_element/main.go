package main

import "fmt"

// ── Approach 1: Brute Force (Shift Left) ─────────────────────────────────────
//
// bruteForce solves Remove Element by shifting elements left when a match is found.
//
// Intuition: When we encounter nums[i] == val, shift all subsequent elements
// one position left to overwrite it. Decrement both length and index to
// re-examine the position that now holds a new element.
//
// Algorithm:
//  1. n = len(nums); i = 0.
//  2. While i < n:
//     if nums[i] == val: shift nums[i+1..n-1] left by one; n--.
//     else: i++.
//  3. Return n.
//
// Time:  O(n²) — each match causes an O(n) shift
// Space: O(1)
func bruteForce(nums []int, val int) int {
	n := len(nums)
	i := 0
	for i < n {
		if nums[i] == val {
			// shift everything after i left by one
			for j := i; j < n-1; j++ {
				nums[j] = nums[j+1]
			}
			n-- // one fewer element in the valid portion
		} else {
			i++
		}
	}
	return n
}

// ── Approach 2: Two Pointers — Write Pointer (Optimal) ───────────────────────
//
// twoPointers solves Remove Element in-place with a write pointer.
//
// Intuition: Identical to LeetCode #26 but the filter condition is "not equal
// to val" rather than "not a duplicate". A write pointer k starts at 0.
// Every element that should be kept is copied to nums[k], then k advances.
//
// Algorithm:
//  1. k = 0.
//  2. For each element e in nums:
//     if e != val: nums[k] = e; k++.
//  3. Return k.
//
// Time:  O(n)
// Space: O(1)
func twoPointers(nums []int, val int) int {
	k := 0 // write pointer: next position for a "kept" element
	for _, e := range nums {
		if e != val { // keep this element
			nums[k] = e
			k++
		}
	}
	return k
}

// ── Approach 3: Two Pointers — Swap from End (Optimal, fewer writes) ─────────
//
// swapFromEnd solves Remove Element with minimal writes when val is rare.
//
// Intuition: Use two pointers: left starts at 0, right starts at n-1.
// Whenever nums[left] == val, swap it with nums[right] and shrink right.
// The right pointer marks the end of the "valid" region. Elements at or
// beyond right may or may not equal val — we simply stop including them.
//
// This approach does at most k swaps where k = number of occurrences of val.
// Preferred when val is rare (few writes) but order may change.
//
// Algorithm:
//  1. left = 0, right = len(nums).
//  2. While left < right:
//     if nums[left] == val: swap nums[left] with nums[right-1]; right--.
//     else: left++.
//  3. Return right (== left at termination).
//
// Time:  O(n)
// Space: O(1)
func swapFromEnd(nums []int, val int) int {
	left, right := 0, len(nums)
	for left < right {
		if nums[left] == val {
			// overwrite the match with the last element in the valid region
			nums[left] = nums[right-1]
			right-- // shrink valid region from the right
		} else {
			left++
		}
	}
	return right
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	a1 := []int{3, 2, 2, 3}
	fmt.Printf("nums=[3,2,2,3] val=3  k=%d  nums[:k]=%v  expected k=2, [2,2]\n", bruteForce(a1, 3), a1[:bruteForce([]int{3, 2, 2, 3}, 3)])
	k1 := bruteForce([]int{3, 2, 2, 3}, 3)
	a2 := []int{3, 2, 2, 3}
	fmt.Printf("nums=[3,2,2,3] val=3  k=%d  nums[:k]=%v\n", bruteForce(a2, 3), a2[:k1])

	b1 := []int{0, 1, 2, 2, 3, 0, 4, 2}
	k2 := bruteForce(b1, 2)
	fmt.Printf("nums=[0,1,2,2,3,0,4,2] val=2  k=%d  nums[:k]=%v  expected k=5\n", k2, b1[:k2])

	fmt.Println("\n=== Approach 2: Two Pointers Write Pointer ===")
	c1 := []int{3, 2, 2, 3}
	k3 := twoPointers(c1, 3)
	fmt.Printf("nums=[3,2,2,3] val=3  k=%d  nums[:k]=%v  expected k=2\n", k3, c1[:k3])

	d1 := []int{0, 1, 2, 2, 3, 0, 4, 2}
	k4 := twoPointers(d1, 2)
	fmt.Printf("nums=[0,1,2,2,3,0,4,2] val=2  k=%d  nums[:k]=%v  expected k=5\n", k4, d1[:k4])

	fmt.Println("\n=== Approach 3: Swap from End (Optimal) ===")
	e1 := []int{3, 2, 2, 3}
	k5 := swapFromEnd(e1, 3)
	fmt.Printf("nums=[3,2,2,3] val=3  k=%d  nums[:k]=%v  expected k=2\n", k5, e1[:k5])

	f1 := []int{0, 1, 2, 2, 3, 0, 4, 2}
	k6 := swapFromEnd(f1, 2)
	fmt.Printf("nums=[0,1,2,2,3,0,4,2] val=2  k=%d  nums[:k]=%v  expected k=5\n", k6, f1[:k6])
}
