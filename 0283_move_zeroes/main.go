package main

import "fmt"

// ── Approach 1: Brute Force (Extra Array) ────────────────────────────────────
//
// bruteForce solves Move Zeroes by copying non-zero elements into a fresh
// buffer in order, then padding the rest with zeros, and copying back.
//
// Intuition:
//
//	The result is "all non-zeros first, in their original relative order, then
//	all the zeros". The most literal way to build that is: scan once collecting
//	non-zeros, then fill the tail with zeros. It ignores the in-place / minimal
//	-operations requirements but is the clearest correctness baseline.
//
// Algorithm:
//
//  1. Walk nums; append each non-zero to a buffer.
//  2. Extend the buffer with zeros until it has len(nums) entries.
//  3. Copy the buffer back over nums.
//
// Time:  O(n) — two linear passes.
// Space: O(n) — the auxiliary buffer (violates the in-place follow-up).
func bruteForce(nums []int) {
	buf := make([]int, 0, len(nums))
	for _, v := range nums {
		if v != 0 {
			buf = append(buf, v) // keep non-zeros in order
		}
	}
	for len(buf) < len(nums) {
		buf = append(buf, 0) // pad the remainder with zeros
	}
	copy(nums, buf) // write the arrangement back in place of the original
}

// ── Approach 2: Two Passes, Overwrite Then Fill ──────────────────────────────
//
// twoPass solves Move Zeroes in place: first compact non-zeros to the front by
// overwriting, then zero-fill the tail.
//
// Intuition:
//
//	Keep a write cursor `insert`. Scan left to right; every non-zero is written
//	at nums[insert] and insert advances. After the scan, positions
//	[insert..n-1] are leftovers that must all be zero — so overwrite them. This
//	is O(1) space but writes to the tail even where a zero already sat.
//
// Algorithm:
//
//  1. insert = 0. For each v in nums: if v != 0, nums[insert] = v; insert++.
//  2. For i from insert to n-1: nums[i] = 0.
//
// Time:  O(n) — one compaction pass plus one fill pass.
// Space: O(1) — in place.
func twoPass(nums []int) {
	insert := 0 // next slot to place a non-zero value
	for _, v := range nums {
		if v != 0 {
			nums[insert] = v // compact non-zeros toward the front
			insert++
		}
	}
	for i := insert; i < len(nums); i++ {
		nums[i] = 0 // remaining tail must be zeros
	}
}

// ── Approach 3: Two Pointers, Swap (Optimal) ─────────────────────────────────
//
// twoPointers solves Move Zeroes in place with a single pass that swaps each
// non-zero into the front region, minimising writes.
//
// Intuition:
//
//	Maintain `last` = boundary index where the next non-zero belongs (all
//	elements before it are non-zero). When we meet a non-zero at i, swap
//	nums[i] with nums[last] and advance last. If i == last nothing moves; a
//	real swap only happens when a zero sits at nums[last], so each element is
//	touched at most once — the fewest operations, order preserved.
//
// Algorithm:
//
//	last = 0. For i from 0 to n-1:
//	  if nums[i] != 0: swap(nums[i], nums[last]); last++.
//
// Time:  O(n) — single pass.
// Space: O(1) — in place, order-preserving.
func twoPointers(nums []int) {
	last := 0 // index where the next non-zero should land
	for i := 0; i < len(nums); i++ {
		if nums[i] != 0 {
			nums[i], nums[last] = nums[last], nums[i] // pull non-zero forward
			last++
		}
	}
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	a := []int{0, 1, 0, 3, 12}
	bruteForce(a)
	fmt.Println(a) // [1 3 12 0 0]
	b := []int{0}
	bruteForce(b)
	fmt.Println(b) // [0]

	fmt.Println("=== Approach 2: Two Passes ===")
	c := []int{0, 1, 0, 3, 12}
	twoPass(c)
	fmt.Println(c) // [1 3 12 0 0]
	d := []int{0}
	twoPass(d)
	fmt.Println(d) // [0]

	fmt.Println("=== Approach 3: Two Pointers (Optimal) ===")
	e := []int{0, 1, 0, 3, 12}
	twoPointers(e)
	fmt.Println(e) // [1 3 12 0 0]
	f := []int{0}
	twoPointers(f)
	fmt.Println(f) // [0]
}
