package main

import (
	"fmt"
	"math"
)

// ── Approach 1: Iterative Counting of States (Brute Force) ───────────────────
//
// iterativeStates solves Poor Pigs by incrementing the pig count until the
// number of distinguishable outcomes reaches the number of buckets.
//
// Intuition:
//
//	Each pig is an independent "sensor" that, over the whole experiment, ends
//	in one of (rounds + 1) distinguishable states: it dies in round 1, or
//	round 2, ..., or round `rounds`, or it survives every round. With
//	rounds = minutesToTest / minutesToDie, one pig encodes base = rounds + 1
//	states. p pigs jointly encode base^p distinct outcome vectors, and each
//	distinct vector can be assigned to a different bucket. So we need the
//	smallest p with base^p >= buckets. Brute force: keep multiplying by `base`
//	(reachable = base^p) and count how many pigs that took.
//
// Algorithm:
//  1. base = minutesToTest / minutesToDie + 1 (states per pig).
//  2. pigs = 0, reachable = 1 (base^0 = 1 bucket distinguishable with 0 pigs).
//  3. While reachable < buckets: reachable *= base; pigs++.
//  4. Return pigs.
//
// Time:  O(log_base(buckets)) — one multiply per pig, and pigs grows logarithmically.
// Space: O(1) — two integer counters.
func iterativeStates(buckets int, minutesToDie int, minutesToTest int) int {
	base := minutesToTest/minutesToDie + 1 // distinguishable states per pig
	pigs := 0                              // pigs used so far
	reachable := 1                         // buckets distinguishable = base^pigs
	// Grow the reach one pig at a time until it covers every bucket.
	for reachable < buckets {
		reachable *= base // adding one pig multiplies the outcome space by `base`
		pigs++            // account for that pig
	}
	return pigs
}

// ── Approach 2: Closed-Form Logarithm (Optimal) ─────────────────────────────
//
// logClosedForm solves Poor Pigs directly: the smallest p with base^p >=
// buckets is ceil(log(buckets) / log(base)).
//
// Intuition:
//
//	base^p >= buckets  ⇔  p >= log_base(buckets) = ln(buckets) / ln(base).
//	The smallest integer p satisfying this is the ceiling of that ratio. This
//	skips the loop and answers in constant time. Care is needed with floating
//	point: round the ratio to a nearby integer before ceiling so values like
//	log(1000)/log(10) = 2.9999999 don't wrongly ceil to 4 — we nudge with a
//	tiny epsilon and verify.
//
// Algorithm:
//  1. If buckets == 1, no pig is needed → return 0.
//  2. base = minutesToTest / minutesToDie + 1.
//  3. Compute ratio = ln(buckets) / ln(base).
//  4. Return ceil(ratio) with a small epsilon to counter FP error.
//
// Time:  O(1) — two logarithms and a ceiling.
// Space: O(1).
func logClosedForm(buckets int, minutesToDie int, minutesToTest int) int {
	if buckets == 1 {
		return 0 // only one bucket → it is trivially the poisonous one, no test needed
	}
	base := float64(minutesToTest/minutesToDie + 1) // states per pig, as float
	// p = ceil(log_base(buckets)); the epsilon absorbs floating-point noise so
	// exact powers (e.g. 1000 = 10^3) don't spuriously round up.
	ratio := math.Log(float64(buckets)) / math.Log(base)
	return int(math.Ceil(ratio - 1e-9))
}

func main() {
	fmt.Println("=== Approach 1: Iterative Counting of States (Brute Force) ===")
	fmt.Println(iterativeStates(4, 15, 15)) // expected 2
	fmt.Println(iterativeStates(4, 15, 30)) // expected 2
	fmt.Println(iterativeStates(1, 15, 15)) // expected 0

	fmt.Println("=== Approach 2: Closed-Form Logarithm (Optimal) ===")
	fmt.Println(logClosedForm(4, 15, 15)) // expected 2
	fmt.Println(logClosedForm(4, 15, 30)) // expected 2
	fmt.Println(logClosedForm(1, 15, 15)) // expected 0
}
