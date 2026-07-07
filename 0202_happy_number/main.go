package main

import "fmt"

// digitSquareSum returns the sum of the squares of the decimal digits of n —
// the "step" function every approach iterates.
// Time: O(log₁₀ n) (one loop per digit). Space: O(1).
func digitSquareSum(n int) int {
	sum := 0
	for n > 0 {
		digit := n % 10      // pop the lowest decimal digit
		sum += digit * digit // accumulate its square
		n /= 10              // shrink the number by one digit
	}
	return sum
}

// ── Approach 1: Hash Set Cycle Detection ─────────────────────────────────────
//
// hashSet solves Happy Number by remembering every value produced; revisiting
// one proves we are trapped in a cycle that never reaches 1.
//
// Intuition:
//
//	The process n → digitSquareSum(n) is a deterministic walk on integers.
//	Any number with ≥ 4 digits shrinks (a 13-digit number maps to at most
//	13·81 = 1053, and below 10⁴ the step maps into [1, 4·81=324]), so every
//	trajectory eventually lives inside a small finite set. A walk on a finite
//	set must either hit 1 or repeat a value — and a repeat means an endless
//	cycle excluding 1.
//
// Algorithm:
//  1. seen = empty set.
//  2. While n != 1 and n not in seen: add n to seen, n = digitSquareSum(n).
//  3. Return n == 1.
//
// Time:  O(log n) — the chain from a 2³¹-bound start reaches the sub-1000
//
//	zone in a handful of steps (each step costs O(log n) digit work,
//	and the number of distinct values visited is bounded by a constant
//	once below 1000).
//
// Space: O(log n) — the set stores the visited chain values.
func hashSet(n int) bool {
	seen := map[int]bool{} // every value the trajectory has produced
	for n != 1 && !seen[n] {
		seen[n] = true        // record before stepping, so a revisit is caught
		n = digitSquareSum(n) // advance the walk one step
	}
	return n == 1 // loop left either at 1 (happy) or on a repeat (cycle)
}

// ── Approach 2: Floyd's Cycle Detection (Two Pointers, Optimal) ──────────────
//
// floydCycleDetection solves Happy Number with the tortoise-and-hare trick —
// O(1) space instead of a hash set.
//
// Intuition:
//
//	Treat the sequence n, step(n), step(step(n)), ... as an implicit linked
//	list where step() is the Next pointer. "Unhappy" means this list has a
//	cycle not containing 1; "happy" means it terminates in the self-loop at
//	1. Exactly like LeetCode #141, advance a slow pointer one step and a
//	fast pointer two steps: inside any cycle the fast pointer must lap the
//	slow one, and if the number is happy the fast pointer reaches 1 first.
//
// Algorithm:
//  1. slow = n, fast = step(n).
//  2. While fast != 1 and slow != fast:
//     slow = step(slow); fast = step(step(fast)).
//  3. Return fast == 1 (meeting elsewhere proves a 1-free cycle).
//
// Time:  O(log n) — same trajectory bound as Approach 1; the pointers meet
//
//	within one lap of the constant-size cycle.
//
// Space: O(1) — two integers, no set.
func floydCycleDetection(n int) bool {
	slow := n                 // moves one step at a time
	fast := digitSquareSum(n) // starts one step ahead, moves two at a time
	for fast != 1 && slow != fast {
		slow = digitSquareSum(slow)                 // tortoise: one step
		fast = digitSquareSum(digitSquareSum(fast)) // hare: two steps
	}
	// Either fast reached the fixed point 1 (happy), or the pointers met
	// inside a cycle that never contains 1 (unhappy).
	return fast == 1
}

// ── Approach 3: Hardcoded Cycle (Math Fact) ──────────────────────────────────
//
// hardcodedCycle solves Happy Number using the number-theory fact that the
// ONLY cycle other than the fixed point 1 is
// 4 → 16 → 37 → 58 → 89 → 145 → 42 → 20 → 4.
//
// Intuition:
//
//	Because every trajectory falls below 1000 quickly, one can exhaustively
//	check all small numbers once (offline) and discover there is exactly one
//	non-trivial cycle, and it passes through 4. So at runtime we only need
//	to iterate until the value becomes 1 (happy) or 4 (doomed to loop).
//
// Algorithm:
//  1. While n != 1 and n != 4: n = digitSquareSum(n).
//  2. Return n == 1.
//
// Time:  O(log n) — the walk reaches {1, 4} within a bounded number of steps
//
//	(each step is O(log n) digit work).
//
// Space: O(1) — a single integer.
func hardcodedCycle(n int) bool {
	// 4 is the sentinel: every unhappy trajectory provably passes through it.
	for n != 1 && n != 4 {
		n = digitSquareSum(n)
	}
	return n == 1
}

func main() {
	fmt.Println("=== Approach 1: Hash Set Cycle Detection ===")
	fmt.Printf("n=19          got=%t  expected true\n", hashSet(19))
	fmt.Printf("n=2           got=%t  expected false\n", hashSet(2))
	fmt.Printf("n=1           got=%t  expected true\n", hashSet(1))           // fixed point edge
	fmt.Printf("n=7           got=%t  expected true\n", hashSet(7))           // 7→49→97→130→10→1
	fmt.Printf("n=2147483647  got=%t  expected false\n", hashSet(2147483647)) // constraint maximum

	fmt.Println("=== Approach 2: Floyd's Cycle Detection (Two Pointers, Optimal) ===")
	fmt.Printf("n=19          got=%t  expected true\n", floydCycleDetection(19))
	fmt.Printf("n=2           got=%t  expected false\n", floydCycleDetection(2))
	fmt.Printf("n=1           got=%t  expected true\n", floydCycleDetection(1))
	fmt.Printf("n=7           got=%t  expected true\n", floydCycleDetection(7))
	fmt.Printf("n=2147483647  got=%t  expected false\n", floydCycleDetection(2147483647))

	fmt.Println("=== Approach 3: Hardcoded Cycle (Math Fact) ===")
	fmt.Printf("n=19          got=%t  expected true\n", hardcodedCycle(19))
	fmt.Printf("n=2           got=%t  expected false\n", hardcodedCycle(2))
	fmt.Printf("n=1           got=%t  expected true\n", hardcodedCycle(1))
	fmt.Printf("n=7           got=%t  expected true\n", hardcodedCycle(7))
	fmt.Printf("n=2147483647  got=%t  expected false\n", hardcodedCycle(2147483647))
}
