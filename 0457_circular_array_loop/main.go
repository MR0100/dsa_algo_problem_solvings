package main

import "fmt"

// next computes the index reached by taking one jump from index `i` in a
// circular array of length n. Positive nums[i] jumps forward, negative jumps
// backward; Go's `%` can be negative, so we add n before the final mod to keep
// the result in [0, n).
func next(nums []int, i int) int {
	n := len(nums)
	// (i + nums[i]) may be negative or ≥ n; normalise into [0, n) circularly.
	return ((i+nums[i])%n + n) % n
}

// ── Approach 1: Brute Force (Per-Start Visited Set) ──────────────────────────
//
// bruteForce solves Circular Array Loop by starting a walk from every index and
// following jumps until it either closes a valid same-direction cycle of length
// > 1, or violates a rule.
//
// Intuition:
//
//	A valid cycle must (a) return to a starting index, (b) keep a single
//	direction throughout, and (c) have length > 1. Try each index as the start
//	of its own walk, remembering the indices visited on THIS walk in a set. If
//	we ever revisit an index from the current walk we have closed a loop — then
//	check length > 1. The direction rule is enforced greedily: the moment a
//	step's sign differs from the start's sign, or we land back on ourselves in
//	one hop (self-loop), abandon this start.
//
// Algorithm:
//  1. For each start index s:
//     a. Record the sign of nums[s] (the required direction).
//     b. Walk with a fresh `seen` set; at each node compute the next index.
//     c. If the next step's sign differs from the required direction → break.
//     d. If next == current (one-element self-loop) → break.
//     e. If next is already in `seen` → a cycle of length > 1 was found → true.
//  2. If no start yields a cycle, return false.
//
// Time:  O(n^2) — up to n starts, each walking up to n steps with a set.
// Space: O(n) — the per-walk visited set.
func bruteForce(nums []int) bool {
	n := len(nums)
	for s := 0; s < n; s++ {
		dir := nums[s] > 0     // required direction: true = forward, false = backward
		seen := map[int]bool{} // indices visited on this particular walk
		cur := s
		for {
			// Direction check: every element on the cycle must share the sign
			// of the start; a mixed-sign path is disqualified.
			if (nums[cur] > 0) != dir {
				break // direction flipped → this start cannot yield a valid cycle
			}
			nxt := next(nums, cur) // one circular jump from cur
			if nxt == cur {
				break // self-loop (length 1) is explicitly not allowed
			}
			if seen[cur] {
				return true // we have returned to a node already on this walk → cycle len > 1
			}
			seen[cur] = true // mark cur as visited on this walk
			cur = nxt        // advance
		}
	}
	return false // no start index produced a valid cycle
}

// ── Approach 2: Floyd's Tortoise & Hare, In-Place Marking (Optimal) ──────────
//
// floydFastSlow solves Circular Array Loop with fast/slow pointers per start,
// short-circuiting on direction changes and marking exhausted nodes with 0 so
// the total work stays linear.
//
// Intuition:
//
//	Detecting a cycle is exactly Floyd's algorithm: a slow pointer (one hop)
//	and a fast pointer (two hops); if they meet, a cycle exists. Two problem
//	specifics are layered on:
//	  • Direction: within one start attempt every visited element must keep the
//	    same sign; the moment `next` would change sign, this attempt is dead.
//	  • Length > 1: a self-loop (slow==next(slow)) is rejected.
//	After a start fails, every node touched on that failed path can never be
//	part of a valid cycle, so we overwrite them with 0 as a "dead" marker —
//	guaranteeing each element is processed O(1) amortised times → overall O(n).
//
// Algorithm:
//  1. For each start i (skip if nums[i] == 0, already dead):
//     a. dir = sign(nums[i]); slow = fast = i.
//     b. Repeat: advance slow one valid step, fast two valid steps, where a
//     "valid step" requires the destination to keep direction `dir` and not
//     be a self-loop; if any step is invalid, break out (bad path).
//     c. If slow == fast → a cycle of length > 1 was found → return true.
//  2. Walk the failed path again from i, marking each node 0 until a
//     direction change (so future starts skip them).
//  3. If no start succeeds, return false.
//
// Time:  O(n) — each index is fully traversed at most a constant number of
//
//	times thanks to the 0-marking of dead paths.
//
// Space: O(1) — pointers only; marking reuses the input array.
func floydFastSlow(nums []int) bool {
	n := len(nums)

	// sameDirection reports whether index i still moves in direction `dir`
	// (dir true = forward). A node whose sign disagrees ends the current path.
	sameDirection := func(i int, dir bool) bool {
		return (nums[i] > 0) == dir
	}

	for i := 0; i < n; i++ {
		if nums[i] == 0 {
			continue // already marked dead by a previous failed attempt
		}
		dir := nums[i] > 0 // required direction for this attempt
		slow, fast := i, i // both pointers start at i

		for {
			// Advance slow by one valid hop.
			slow = validNext(nums, slow, dir, sameDirection)
			if slow == -1 {
				break // slow hit a direction change or self-loop → path invalid
			}
			// Advance fast by two valid hops.
			fast = validNext(nums, fast, dir, sameDirection)
			if fast == -1 {
				break
			}
			fast = validNext(nums, fast, dir, sameDirection)
			if fast == -1 {
				break
			}
			if slow == fast {
				return true // pointers met inside a same-direction cycle of length > 1
			}
		}

		// This start failed: overwrite the whole path with 0 so it is never
		// re-explored, keeping the amortised cost linear.
		j := i
		for nums[j] != 0 && sameDirection(j, dir) {
			nxt := next(nums, j) // remember where this node pointed
			nums[j] = 0          // mark node j as dead
			if nxt == j {
				break // self-loop node; stop marking
			}
			j = nxt
		}
	}
	return false // no valid cycle from any start
}

// validNext returns the index reached by one hop from i IF that hop keeps
// direction `dir` and is not a self-loop; otherwise it returns -1 to signal an
// invalid step (used only by floydFastSlow).
func validNext(nums []int, i int, dir bool, sameDirection func(int, bool) bool) int {
	if !sameDirection(i, dir) {
		return -1 // sign of nums[i] disagrees with the cycle's direction
	}
	nxt := next(nums, i) // circular one-step jump
	if nxt == i {
		return -1 // length-1 self-loop is not a valid cycle
	}
	return nxt
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Per-Start Visited Set) ===")
	fmt.Println(bruteForce([]int{2, -1, 1, 2, 2}))        // expected true
	fmt.Println(bruteForce([]int{-1, -2, -3, -4, -5, 6})) // expected false
	fmt.Println(bruteForce([]int{1, -1, 5, 1, 4}))        // expected true

	// floydFastSlow mutates its input (0-marking), so pass fresh copies.
	fmt.Println("=== Approach 2: Floyd's Tortoise & Hare, In-Place Marking (Optimal) ===")
	fmt.Println(floydFastSlow([]int{2, -1, 1, 2, 2}))        // expected true
	fmt.Println(floydFastSlow([]int{-1, -2, -3, -4, -5, 6})) // expected false
	fmt.Println(floydFastSlow([]int{1, -1, 5, 1, 4}))        // expected true
}
