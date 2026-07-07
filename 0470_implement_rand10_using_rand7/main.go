package main

import "fmt"

// ─────────────────────────────────────────────────────────────────────────────
// Provided API (simulated).
//
// On LeetCode, rand7() is given: it returns a uniform integer in [1, 7]. To make
// this file self-contained and testable we implement rand7() with a tiny
// deterministic linear-congruential generator, so runs are reproducible and we
// can statistically check rand10()'s uniformity. The algorithms below only ever
// call rand7() — they never touch this internal state directly.
// ─────────────────────────────────────────────────────────────────────────────

var lcgState uint64 = 123456789 // fixed seed → deterministic, reproducible output

// rand7 returns a uniform random integer in [1, 7] (simulated).
func rand7() int {
	// Numerical Recipes LCG constants; high bits are the "more random" ones.
	lcgState = lcgState*6364136223846793005 + 1442695040888963407
	return int((lcgState>>33)%7) + 1 // map [0,6] → [1,7]
}

// ── Approach 1: Rejection Sampling on a 7×7 Grid (Optimal, standard) ──────────
//
// rejectionSampling implements rand10() by combining two rand7() calls into a
// uniform value in [1, 49], then rejecting the top 9 outcomes so the remaining
// 40 map perfectly (4 per output) onto [1, 10].
//
// Intuition:
//
//	One rand7() gives 7 equally likely values; two independent calls give a
//	uniform point in a 7×7 = 49-cell grid, i.e. a uniform integer in [1, 49].
//	40 is the largest multiple of 10 that is ≤ 49, so keep only values in
//	[1, 40] (reject 41..49 and retry) and reduce mod 10. Rejection is what keeps
//	the result EXACTLY uniform — you must not "fold" the leftover 9 into the
//	range, because 49 is not divisible by 10 and any reuse would bias it.
//
// Algorithm:
//  1. row = rand7() (1..7), col = rand7() (1..7).
//  2. idx = (row-1)*7 + col  → a uniform integer in [1, 49].
//  3. If idx > 40, discard and go back to step 1 (rejection).
//  4. Return (idx-1)%10 + 1 → a uniform integer in [1, 10].
//
// Time:  Expected O(1) — each trial succeeds with probability 40/49 ≈ 0.816, so
//
//	the expected number of trials is 49/40 ≈ 1.225, i.e. ~2.45 rand7() calls.
//
// Space: O(1) — no state beyond a few locals.
func rejectionSampling() int {
	for {
		row := rand7()         // 1..7 — chooses the grid row
		col := rand7()         // 1..7 — chooses the grid column
		idx := (row-1)*7 + col // 1..49, uniform over the 7×7 grid
		if idx <= 40 {         // keep only the first 40 (a multiple of 10)
			return (idx-1)%10 + 1 // fold 40 outcomes → [1,10], 4 each, unbiased
		}
		// idx in 41..49 → reject and retry so the distribution stays exact
	}
}

// ── Approach 2: Rejection Sampling that Reuses Rejects (Fewer rand7 Calls) ────
//
// reuseRejects answers the follow-up ("minimise rand7() calls") by not throwing
// the rejected randomness away: the 9 rejected values 41..49 still carry a
// uniform integer in [1, 9], which is combined with a fresh rand7() into [1, 63],
// and if that also lands in the reject zone its leftover [1, 3] is combined once
// more into [1, 21]. Each stage salvages entropy instead of restarting cold.
//
// Intuition:
//
//	After computing idx in [1, 49], if idx > 40 we still hold a uniform value in
//	[1, 9] (idx-40). Multiply that (minus 1) by 7 and add a new rand7() to get a
//	uniform value in [1, 63]; keep [1, 60] (→ [1,10]) and otherwise you hold a
//	uniform [1, 3]. Combine that with another rand7() into [1, 21]; keep [1, 20]
//	(→ [1,10]), else you hold a uniform [1, 1] and just loop. Reusing the reject
//	means fewer wasted rand7() calls than restarting from scratch each time.
//
// Algorithm:
//  1. a = rand7(), b = rand7(); idx = (a-1)*7 + b  (1..49). If idx ≤ 40 return
//     (idx-1)%10+1.
//  2. Else salvage rem = idx-40 (1..9). c = rand7(); idx = (rem-1)*7 + c (1..63).
//     If idx ≤ 60 return (idx-1)%10+1.
//  3. Else rem = idx-60 (1..3). d = rand7(); idx = (rem-1)*7 + d (1..21).
//     If idx ≤ 20 return (idx-1)%10+1.
//  4. Else loop back to step 1.
//
// Time:  Expected ~2.19 rand7() calls per rand10() (vs ~2.45 for Approach 1) —
//
//	still O(1) expected, but tighter because rejected entropy is recycled.
//
// Space: O(1).
func reuseRejects() int {
	for {
		a := rand7()
		b := rand7()
		idx := (a-1)*7 + b // 1..49
		if idx <= 40 {
			return (idx-1)%10 + 1
		}
		// Salvage the uniform [1,9] hiding in the rejected 41..49.
		rem := idx - 40 // 1..9
		c := rand7()
		idx = (rem-1)*7 + c // 1..63
		if idx <= 60 {
			return (idx-1)%10 + 1
		}
		// Salvage the uniform [1,3] hiding in the rejected 61..63.
		rem = idx - 60 // 1..3
		d := rand7()
		idx = (rem-1)*7 + d // 1..21
		if idx <= 20 {
			return (idx-1)%10 + 1
		}
		// Only 1 value left (21); nothing to salvage — loop and start over.
	}
}

// histogram calls fn n times and returns counts[1..10] of each outcome, plus a
// flag for whether every result stayed in range [1,10]. Used to sanity-check
// uniformity in main() since exact outputs are random.
func histogram(fn func() int, n int) ([11]int, bool) {
	var counts [11]int
	inRange := true
	for i := 0; i < n; i++ {
		v := fn()
		if v < 1 || v > 10 {
			inRange = false // must never happen for a correct rand10
		} else {
			counts[v]++
		}
	}
	return counts, inRange
}

func main() {
	// The official "examples" are just sample outputs of a random function
	// (e.g. n=3 → [3,8,10]); there is no fixed expected sequence. We instead
	// verify the two guarantees that DO hold: outputs land in [1,10], and the
	// distribution is (approximately) uniform over a large sample.
	const N = 700000
	const expectedPerBucket = N / 10 // 70000; uniform target

	fmt.Println("=== Approach 1: Rejection Sampling on a 7x7 Grid (Optimal) ===")
	fmt.Printf("first 3 rand10() outputs: [%d %d %d]  (any value in 1..10)\n", rejectionSampling(), rejectionSampling(), rejectionSampling())
	c1, ok1 := histogram(rejectionSampling, N)
	fmt.Printf("all outputs in [1,10]: %t  expected true\n", ok1)
	fmt.Printf("bucket counts over %d samples (each ~%d): %v\n", N, expectedPerBucket, c1[1:])
	fmt.Printf("uniform within 3%%: %t  expected true\n", withinTolerance(c1, expectedPerBucket, 0.03))

	fmt.Println("=== Approach 2: Rejection Sampling that Reuses Rejects (Fewer rand7 Calls) ===")
	fmt.Printf("first 3 rand10() outputs: [%d %d %d]  (any value in 1..10)\n", reuseRejects(), reuseRejects(), reuseRejects())
	c2, ok2 := histogram(reuseRejects, N)
	fmt.Printf("all outputs in [1,10]: %t  expected true\n", ok2)
	fmt.Printf("bucket counts over %d samples (each ~%d): %v\n", N, expectedPerBucket, c2[1:])
	fmt.Printf("uniform within 3%%: %t  expected true\n", withinTolerance(c2, expectedPerBucket, 0.03))
}

// withinTolerance reports whether every bucket count is within (1±tol)·expected.
func withinTolerance(counts [11]int, expected int, tol float64) bool {
	lo := float64(expected) * (1 - tol)
	hi := float64(expected) * (1 + tol)
	for v := 1; v <= 10; v++ {
		if float64(counts[v]) < lo || float64(counts[v]) > hi {
			return false // a bucket drifted too far from uniform
		}
	}
	return true
}
