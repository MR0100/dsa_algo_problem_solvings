package main

import "fmt"

// ── Approach 1: Recursion (Top-Down) ─────────────────────────────────────────
//
// recursion solves Integer Replacement by branching on both choices when n is
// odd and taking the minimum.
//
// Intuition:
//
//	The rules: if n is even, halve it (unambiguous, one step); if n is odd, we
//	may either +1 or -1 (both give an even number). Since we don't know which
//	odd move leads to fewer total steps, try both and take the smaller. Even
//	moves are forced, so no branching there.
//
// Algorithm:
//  1. Base case: n == 1 → 0 steps.
//  2. If n is even → 1 + recursion(n/2).
//  3. If n is odd → 1 + min(recursion(n+1), recursion(n-1)).
//
// Time:  O(log n) typical but up to O(2^log n) branching in the worst
//
//	odd-heavy case (no memo); each step at least halves within a branch.
//
// Space: O(log n) — recursion stack depth.
//
// Note: use uint to avoid int overflow when n = 2^31 - 1 and we compute n+1.
func recursion(n int) int {
	var solve func(x uint) int
	solve = func(x uint) int {
		if x == 1 {
			return 0 // reached the target, no more operations
		}
		if x%2 == 0 {
			return 1 + solve(x/2) // even: the only legal move is to halve
		}
		// odd: try both neighbours (each becomes even) and take the cheaper.
		return 1 + min(solve(x+1), solve(x-1))
	}
	return solve(uint(n))
}

// ── Approach 2: Greedy Bit Manipulation (Optimal) ────────────────────────────
//
// greedyBits solves Integer Replacement by deciding each odd step from the two
// lowest bits, clearing trailing 1s efficiently.
//
// Intuition:
//
//	Even → shift right (halve). For odd n, the choice +1 vs -1 should remove as
//	many low set bits as possible. Look at the low two bits:
//	  - n == 3 is the special case: 3 -> 2 -> 1 (subtract) beats 3 -> 4 -> 2 ->
//	    1 (add), so subtract.
//	  - Otherwise if the second-lowest bit is 1 (n ends in ...11), adding 1
//	    triggers a carry that clears a run of trailing 1s → n++.
//	  - If it is 0 (n ends in ...01), subtracting 1 clears the single low 1 → n--.
//	Each step then halves, so we march the number down to 1.
//
// Algorithm:
//  1. count = 0.
//  2. While n != 1:
//     a. If n even → n >>= 1.
//     b. Else if n == 3 or bit1 is 0 (n&2 == 0) → n-- (strip the lone low 1).
//     c. Else → n++ (carry away a run of trailing 1s).
//     d. count++.
//  3. Return count.
//
// Time:  O(log n) — each iteration removes at least one bit of magnitude.
// Space: O(1) — a counter, no recursion.
//
// Note: uint again guards the n = 2^31 - 1 (+1) overflow.
func greedyBits(n int) int {
	x := uint(n)
	count := 0
	for x != 1 {
		if x%2 == 0 {
			x >>= 1 // even: halve (drop the trailing 0)
		} else if x == 3 || x&2 == 0 {
			// x == 3: subtract to avoid the longer 3->4->2->1 path.
			// x&2 == 0 (ends in 01): subtracting clears the single low 1 bit.
			x--
		} else {
			// ends in 11: adding 1 carries and clears a run of trailing 1s.
			x++
		}
		count++ // every branch above performed one operation
	}
	return count
}

func main() {
	fmt.Println("=== Approach 1: Recursion (Top-Down) ===")
	fmt.Printf("n=8           got=%d  expected 3\n", recursion(8))
	fmt.Printf("n=7           got=%d  expected 4\n", recursion(7))
	fmt.Printf("n=4           got=%d  expected 2\n", recursion(4))

	fmt.Println("=== Approach 2: Greedy Bit Manipulation (Optimal) ===")
	fmt.Printf("n=8           got=%d  expected 3\n", greedyBits(8))
	fmt.Printf("n=7           got=%d  expected 4\n", greedyBits(7))
	fmt.Printf("n=4           got=%d  expected 2\n", greedyBits(4))
	fmt.Printf("n=2147483647  got=%d  expected 32\n", greedyBits(2147483647)) // overflow-safe edge
}
