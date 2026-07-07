package main

import "fmt"

// ── Approach 1: Brute Force (Rotate One Step, k Times) ───────────────────────
//
// bruteForce solves Rotate Array by performing k single-step right rotations.
//
// Intuition:
//
//	A rotation by k is by definition k rotations by one. One right rotation
//	is easy: save the last element, shift everything else one slot right,
//	drop the saved element into slot 0. Repeat k times. Reducing k modulo n
//	first avoids useless full cycles (rotating by n is the identity).
//
// Algorithm:
//  1. k %= n (full cycles change nothing).
//  2. Repeat k times: save nums[n-1]; shift nums[i-1] → nums[i] for i from
//     n-1 down to 1; put the saved value in nums[0].
//
// Time:  O(n·k) — each of the k steps shifts all n elements.
// Space: O(1) — a single saved element.
func bruteForce(nums []int, k int) {
	n := len(nums)
	if n == 0 {
		return
	}
	k %= n // rotating by n is a no-op, so only the remainder matters
	for step := 0; step < k; step++ {
		last := nums[n-1] // the element that wraps around to the front
		// shift every element one slot to the right (right to left to avoid clobbering)
		for i := n - 1; i > 0; i-- {
			nums[i] = nums[i-1]
		}
		nums[0] = last // wrapped element lands at the front
	}
}

// ── Approach 2: Extra Array ──────────────────────────────────────────────────
//
// extraArray solves Rotate Array by writing every element directly to its
// final position in a second array.
//
// Intuition:
//
//	The final position of the element at index i is known in closed form:
//	(i + k) mod n. So build a fresh array, drop each element straight where
//	it belongs, and copy the result back — one pass, no shifting.
//
// Algorithm:
//  1. Allocate rotated[n].
//  2. For every i: rotated[(i+k) % n] = nums[i].
//  3. Copy rotated back into nums.
//
// Time:  O(n) — one placement pass plus one copy pass.
// Space: O(n) — the auxiliary rotated array (violates the O(1) follow-up).
func extraArray(nums []int, k int) {
	n := len(nums)
	if n == 0 {
		return
	}
	rotated := make([]int, n) // scratch array in final order
	for i, v := range nums {
		rotated[(i+k)%n] = v // closed-form destination of index i
	}
	copy(nums, rotated) // the problem wants nums itself mutated
}

// ── Approach 3: Cyclic Replacements ──────────────────────────────────────────
//
// cyclicReplacements solves Rotate Array in place by following each
// displacement cycle, dropping every element into its final slot directly.
//
// Intuition:
//
//	"Move i to (i+k) mod n" decomposes the indices into gcd(n,k) disjoint
//	cycles. Walk a cycle carrying the displaced value in hand: place it,
//	pick up the victim, hop k forward, repeat — the cycle closes exactly at
//	its start. Counting total placements tells us when all n elements are
//	home; if a cycle closes early (gcd > 1), start the next cycle one index
//	later.
//
// Algorithm:
//  1. k %= n; if k == 0 nothing moves.
//  2. For start = 0, 1, ... while fewer than n elements are placed:
//     walk current → (current+k) % n, swapping the carried value into each
//     visited slot, until the walk returns to start.
//
// Time:  O(n) — every element is picked up and placed exactly once.
// Space: O(1) — one carried value plus a few indices.
func cyclicReplacements(nums []int, k int) {
	n := len(nums)
	if n == 0 {
		return
	}
	k %= n
	if k == 0 {
		return // identity rotation: touching nothing is correct
	}
	count := 0 // how many elements have reached their final slot
	for start := 0; count < n; start++ {
		current := start
		carried := nums[start] // value in hand, waiting to be placed k ahead
		for {
			next := (current + k) % n // final slot of the carried value
			// place the carried value and pick up the displaced one in a single swap
			nums[next], carried = carried, nums[next]
			current = next
			count++
			if current == start { // cycle closed — everything on it is placed
				break
			}
		}
	}
}

// ── Approach 4: Reversal (Optimal) ───────────────────────────────────────────
//
// reversal solves Rotate Array with three in-place reversals.
//
// Intuition:
//
//	Rotating right by k moves the last k elements to the front (keeping
//	their order) and the first n-k behind them (keeping their order).
//	Reversing the whole array puts the last k elements at the front but
//	each block comes out internally backwards; reversing each block
//	separately repairs the internal order. Two reversals cancel within a
//	block, but the block-order flip from the full reversal survives.
//
// Algorithm:
//  1. k %= n.
//  2. Reverse the entire array.
//  3. Reverse the first k elements.
//  4. Reverse the remaining n-k elements.
//
// Time:  O(n) — each element is swapped at most twice across the three passes.
// Space: O(1) — reversals swap in place.
func reversal(nums []int, k int) {
	n := len(nums)
	if n == 0 {
		return
	}
	k %= n
	reverseRange(nums, 0, n-1) // whole array: tail block arrives at the front (backwards)
	reverseRange(nums, 0, k-1) // fix internal order of the first k (the old tail)
	reverseRange(nums, k, n-1) // fix internal order of the rest (the old head)
}

// reverseRange reverses nums[lo..hi] in place with two converging pointers.
func reverseRange(nums []int, lo, hi int) {
	for lo < hi {
		nums[lo], nums[hi] = nums[hi], nums[lo] // swap the outermost pair
		lo++
		hi--
	}
}

// clone returns an independent copy so each approach starts from the original input.
func clone(nums []int) []int {
	c := make([]int, len(nums))
	copy(c, nums)
	return c
}

func main() {
	// Example 1: nums = [1,2,3,4,5,6,7], k = 3 → [5,6,7,1,2,3,4]
	// Example 2: nums = [-1,-100,3,99],  k = 2 → [3,99,-1,-100]
	ex1, k1 := []int{1, 2, 3, 4, 5, 6, 7}, 3
	ex2, k2 := []int{-1, -100, 3, 99}, 2

	fmt.Println("=== Approach 1: Brute Force (Rotate One Step, k Times) ===")
	a1 := clone(ex1)
	bruteForce(a1, k1)
	fmt.Println(a1) // expected: [5 6 7 1 2 3 4]
	a2 := clone(ex2)
	bruteForce(a2, k2)
	fmt.Println(a2) // expected: [3 99 -1 -100]

	fmt.Println("=== Approach 2: Extra Array ===")
	b1 := clone(ex1)
	extraArray(b1, k1)
	fmt.Println(b1) // expected: [5 6 7 1 2 3 4]
	b2 := clone(ex2)
	extraArray(b2, k2)
	fmt.Println(b2) // expected: [3 99 -1 -100]

	fmt.Println("=== Approach 3: Cyclic Replacements ===")
	c1 := clone(ex1)
	cyclicReplacements(c1, k1)
	fmt.Println(c1) // expected: [5 6 7 1 2 3 4]
	c2 := clone(ex2)
	cyclicReplacements(c2, k2)
	fmt.Println(c2) // expected: [3 99 -1 -100]

	fmt.Println("=== Approach 4: Reversal (Optimal) ===")
	d1 := clone(ex1)
	reversal(d1, k1)
	fmt.Println(d1) // expected: [5 6 7 1 2 3 4]
	d2 := clone(ex2)
	reversal(d2, k2)
	fmt.Println(d2) // expected: [3 99 -1 -100]
}
