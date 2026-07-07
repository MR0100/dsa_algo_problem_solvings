package main

import "fmt"

// The permutation of [1..n] described by a length-(n-1) string of 'I'/'D':
//   s[i] == 'I'  ⇒  perm[i] < perm[i+1]   (increase)
//   s[i] == 'D'  ⇒  perm[i] > perm[i+1]   (decrease)
// We must return the LEXICOGRAPHICALLY SMALLEST such permutation.
//
// Core insight shared by all approaches: the identity [1,2,…,n] is already the
// lexicographically smallest sequence and it satisfies an all-'I' string. Each
// 'D' forces a local descent. To keep things as small as possible while
// obeying a maximal block of consecutive 'D's spanning positions i..j, we take
// the ascending numbers that would occupy indices i..j+1 and REVERSE just that
// window — the minimal disturbance that turns "increasing" into "decreasing".

// reverseInts reverses a[lo..hi] inclusive, in place (shared helper).
func reverseInts(a []int, lo, hi int) {
	for lo < hi {
		a[lo], a[hi] = a[hi], a[lo]
		lo++
		hi--
	}
}

// ── Approach 1: Brute Force via next_permutation Scan ────────────────────────
//
// bruteForce starts from the identity permutation and, using the standard
// "next lexicographic permutation" generator, walks permutations in increasing
// order, returning the FIRST one that matches the I/D pattern.
//
// Intuition:
//
//	We want the smallest permutation satisfying the pattern. Enumerate
//	permutations of [1..n] in lexicographic order (identity first) and test
//	each against s; the first that fits is by definition the smallest. This is
//	the naive baseline — correct, and a useful oracle to check the greedy
//	answer against for small n — but factorial in the worst case.
//
// Algorithm:
//  1. perm = [1,2,…,n].
//  2. Loop: if perm satisfies s, return it; else advance perm to its next
//     lexicographic permutation (until none remain).
//
// Time:  O(n! · n) worst case — up to n! permutations, each validated in O(n).
// Space: O(n) — the working permutation.
func bruteForce(s string) []int {
	n := len(s) + 1
	perm := make([]int, n)
	for i := range perm {
		perm[i] = i + 1 // identity = lexicographically first permutation
	}

	for {
		if satisfies(perm, s) { // first fit in lex order is the smallest
			return perm
		}
		if !nextPermutation(perm) { // exhausted all permutations
			return perm // (unreachable for valid input; identity handles all-'I')
		}
	}
}

// satisfies reports whether perm obeys every I/D constraint in s.
func satisfies(perm []int, s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == 'I' && !(perm[i] < perm[i+1]) {
			return false // 'I' demands an ascent here
		}
		if s[i] == 'D' && !(perm[i] > perm[i+1]) {
			return false // 'D' demands a descent here
		}
	}
	return true
}

// nextPermutation rearranges perm into the next lexicographically greater
// permutation in place; returns false if perm was the last (descending) one.
func nextPermutation(perm []int) bool {
	n := len(perm)
	i := n - 2
	for i >= 0 && perm[i] >= perm[i+1] { // find last ascent perm[i] < perm[i+1]
		i--
	}
	if i < 0 {
		return false // wholly descending → no greater permutation
	}
	j := n - 1
	for perm[j] <= perm[i] { // find rightmost element greater than perm[i]
		j--
	}
	perm[i], perm[j] = perm[j], perm[i] // swap pivot with its successor
	reverseInts(perm, i+1, n-1)         // reverse the suffix to its smallest order
	return true
}

// ── Approach 2: Greedy Reverse of D-Runs (Optimal, array walk) ────────────────
//
// greedyReverseRuns fills the identity permutation, then reverses each maximal
// consecutive block of 'D's in one sweep.
//
// Intuition:
//
//	Start from [1,2,…,n] (already lex-smallest, satisfies all 'I'). Wherever s
//	has a maximal run of 'D' at positions i..j, indices i..j+1 must strictly
//	decrease. The values sitting there are the consecutive ascending numbers
//	i+1..j+2; reversing exactly that window makes them descend while leaving
//	every value outside untouched — the smallest possible change. Isolated 'I'
//	positions need no work. Sweeping left to right and reversing each D-block
//	yields the global lex-smallest permutation.
//
// Algorithm:
//  1. perm = [1,2,…,n].
//  2. i = 0. While i < len(s): if s[i] == 'D', find the run end j (last
//     consecutive 'D'), reverse perm[i .. j+1], then jump i to j+1; else i++.
//  3. Return perm.
//
// Time:  O(n) — each element is moved by at most one reversal; runs are disjoint.
// Space: O(n) — the output permutation (O(1) beyond it).
func greedyReverseRuns(s string) []int {
	n := len(s) + 1
	perm := make([]int, n)
	for i := range perm {
		perm[i] = i + 1 // baseline ascending permutation
	}

	i := 0
	for i < len(s) {
		if s[i] == 'D' {
			j := i
			for j < len(s) && s[j] == 'D' { // extend over the whole D-run
				j++
			}
			// s[i..j-1] are 'D' ⇒ indices i..j must descend; reverse that window.
			reverseInts(perm, i, j)
			i = j // continue past the run (position j was the boundary)
		} else {
			i++ // 'I' already satisfied by the ascending baseline
		}
	}
	return perm
}

// ── Approach 3: Stack-Based Emit-on-I (Optimal, one pass) ─────────────────────
//
// stackEmit pushes the numbers 1..n in order onto a stack and flushes the stack
// whenever it meets an 'I' (or the end), which naturally reverses each pending
// descending block.
//
// Intuition:
//
//	Push 1,2,3,… onto a stack. A stack reverses order on pop, so numbers held
//	between flushes come out descending — exactly what a run of 'D' needs. At
//	every 'I' boundary (and after the final number) we pop everything queued so
//	far into the output. The delayed, LIFO flush produces the same lex-smallest
//	permutation as reversing D-runs, in a single left-to-right pass without an
//	explicit reverse.
//
// Algorithm:
//  1. For pos = 0 … n-1: push (pos+1) on the stack.
//  2. If pos == n-1 OR s[pos] == 'I': pop the whole stack into the output.
//  3. The output, filled in pop order, is the answer.
//
// Time:  O(n) — each number is pushed once and popped once.
// Space: O(n) — stack plus output.
func stackEmit(s string) []int {
	n := len(s) + 1
	out := make([]int, 0, n)
	stack := make([]int, 0, n)

	for pos := 0; pos < n; pos++ {
		stack = append(stack, pos+1) // queue the next natural number
		// Flush at an 'I' boundary or at the very end: everything queued since
		// the last flush belongs to one descending block, reversed by the pops.
		if pos == n-1 || s[pos] == 'I' {
			for len(stack) > 0 {
				out = append(out, stack[len(stack)-1]) // pop top (LIFO ⇒ reversed)
				stack = stack[:len(stack)-1]
			}
		}
	}
	return out
}

func main() {
	fmt.Println("=== Approach 1: Brute Force via next_permutation Scan ===")
	fmt.Printf("s=%q    got=%v  expected [1 2]\n", "I", bruteForce("I"))          // [1 2]
	fmt.Printf("s=%q   got=%v  expected [2 1 3]\n", "DI", bruteForce("DI"))       // [2 1 3]
	fmt.Printf("s=%q got=%v  expected [3 2 1 5 4]\n", "DDID", bruteForce("DDID")) // [3 2 1 5 4]

	fmt.Println("=== Approach 2: Greedy Reverse of D-Runs (Optimal, array walk) ===")
	fmt.Printf("s=%q    got=%v  expected [1 2]\n", "I", greedyReverseRuns("I"))
	fmt.Printf("s=%q   got=%v  expected [2 1 3]\n", "DI", greedyReverseRuns("DI"))
	fmt.Printf("s=%q got=%v  expected [3 2 1 5 4]\n", "DDID", greedyReverseRuns("DDID"))

	fmt.Println("=== Approach 3: Stack-Based Emit-on-I (Optimal, one pass) ===")
	fmt.Printf("s=%q    got=%v  expected [1 2]\n", "I", stackEmit("I"))
	fmt.Printf("s=%q   got=%v  expected [2 1 3]\n", "DI", stackEmit("DI"))
	fmt.Printf("s=%q got=%v  expected [3 2 1 5 4]\n", "DDID", stackEmit("DDID"))
}
