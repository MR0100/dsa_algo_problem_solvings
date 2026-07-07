package main

import "fmt"

// ── Approach 1: Brute Force Simulation (List) ────────────────────────────────
//
// bruteForce solves Elimination Game by literally building the list [1..n] and
// removing every other element, alternating direction each pass, until one
// number remains.
//
// Intuition:
//
//	Just do exactly what the problem says. Keep the surviving numbers in a
//	slice. On a left-to-right pass, keep indices 1,3,5,... (drop the first,
//	third-from-view, etc.). On a right-to-left pass, keep them symmetrically.
//	Flip direction each round. When one element is left, that is the answer.
//	Correct but O(n) memory and O(n) work rebuilding the list — TLE / MLE for
//	n up to 10⁹.
//
// Algorithm:
//  1. nums = [1..n]; leftToRight = true.
//  2. While len(nums) > 1: rebuild keeping every other element depending on
//     direction and parity of length; flip direction.
//  3. Return nums[0].
//
// Time:  O(n) total (list halves each pass: n + n/2 + ... ≈ 2n).
// Space: O(n) — stores the whole list.
func bruteForce(n int) int {
	nums := make([]int, n) // the surviving numbers, initially 1..n
	for i := 0; i < n; i++ {
		nums[i] = i + 1
	}
	leftToRight := true // direction of the current pass

	for len(nums) > 1 {
		next := make([]int, 0, len(nums)/2) // survivors of this pass
		if leftToRight {
			// Left→right: always remove the first, so keep indices 1,3,5,...
			for i := 1; i < len(nums); i += 2 {
				next = append(next, nums[i])
			}
		} else {
			// Right→left: mirror image. If length is even, we keep 0,2,4,...;
			// if odd, we keep 1,3,5,... (the leftmost is removed).
			start := 0
			if len(nums)%2 == 1 {
				start = 1
			}
			for i := start; i < len(nums); i += 2 {
				next = append(next, nums[i])
			}
		}
		nums = next                // advance to the survivors
		leftToRight = !leftToRight // alternate the sweep direction
	}
	return nums[0]
}

// ── Approach 2: Track the Head Pointer (Optimal, O(log n)) ───────────────────
//
// headPointer solves Elimination Game by tracking only the value of the FIRST
// remaining number ("head") and how many numbers remain, without storing them.
//
// Intuition:
//
//	The whole answer is determined by where the head lands after all passes.
//	Two facts drive the head's movement:
//	  • On a left→right pass, the head is ALWAYS removed, so head advances by
//	    the current step size.
//	  • On a right→left pass, the head moves ONLY when the count is odd (then
//	    the leftmost element is also removed); if the count is even the head
//	    survives that pass.
//	Each pass halves the remaining count and doubles the step. When only one
//	number remains, head is the answer.
//
// Algorithm:
//  1. head = 1, step = 1, remaining = n, leftToRight = true.
//  2. While remaining > 1:
//     - if leftToRight OR remaining is odd: head += step.
//     - remaining /= 2; step *= 2; flip direction.
//  3. Return head.
//
// Time:  O(log n) — remaining halves each iteration.
// Space: O(1) — a handful of scalars.
func headPointer(n int) int {
	head := 1      // value of the leftmost surviving number
	step := 1      // gap between consecutive survivors
	remaining := n // how many numbers are still in play
	leftToRight := true

	for remaining > 1 {
		// The head is eliminated (so it must move forward by one step) when:
		//  - we sweep left→right (head is always first to go), OR
		//  - we sweep right→left but the count is odd (leftmost also removed).
		if leftToRight || remaining%2 == 1 {
			head += step
		}
		remaining /= 2 // half the numbers are eliminated this pass
		step *= 2      // survivors are now twice as far apart
		leftToRight = !leftToRight
	}
	return head
}

func main() {
	fmt.Println("=== Approach 1: Brute Force Simulation ===")
	fmt.Println(bruteForce(9)) // expected 6
	fmt.Println(bruteForce(1)) // expected 1

	fmt.Println("=== Approach 2: Track Head Pointer (Optimal) ===")
	fmt.Println(headPointer(9)) // expected 6
	fmt.Println(headPointer(1)) // expected 1
}
