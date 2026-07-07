package main

import (
	"fmt"
	"math/bits"
	"sort"
)

// A binary watch: 4 hour LEDs (values 1,2,4,8) show hours 0..11 and 6 minute
// LEDs (values 1,2,4,8,16,32) show minutes 0..59. `turnedOn` is the total
// number of LEDs lit across BOTH rows. We must return every valid "h:mm" the
// watch could display. Formatting: hour has no leading zero; minute is always
// two digits.

// ── Approach 1: Brute Force over all 12×60 times ─────────────────────────────
//
// bruteForce enumerates every possible clock time and keeps the ones whose
// total number of set bits (hour bits + minute bits) equals turnedOn.
//
// Intuition:
//
//	A binary watch simply shows a number in binary. The number of lit LEDs for
//	a given time is popcount(hour) + popcount(minute). There are only 12·60 =
//	720 times, so we can afford to test each one directly instead of choosing
//	which LEDs to light.
//
// Algorithm:
//  1. Loop hour h from 0..11 and minute m from 0..59.
//  2. If popcount(h) + popcount(m) == turnedOn, format "h:mm" and collect it.
//  3. Return all collected strings.
//
// Time:  O(12·60) = O(1) — a fixed 720 iterations regardless of input.
// Space: O(1) besides the output list (at most a few dozen strings).
func bruteForce(turnedOn int) []string {
	var result []string
	for h := 0; h < 12; h++ { // 4 hour LEDs cover 0..11
		for m := 0; m < 60; m++ { // 6 minute LEDs cover 0..59
			// bits.OnesCount counts the set bits (lit LEDs) of the value.
			if bits.OnesCount(uint(h))+bits.OnesCount(uint(m)) == turnedOn {
				// %d gives the hour with no leading zero; %02d pads the
				// minute to exactly two digits (e.g. 2 -> "02").
				result = append(result, fmt.Sprintf("%d:%02d", h, m))
			}
		}
	}
	return result
}

// ── Approach 2: Split the Budget (choose hour-bits, then minute-bits) ─────────
//
// splitBudget picks how many of the lit LEDs belong to the hour, then reuses
// the brute-force idea only within each half. It is essentially the same 720
// scan but framed as "distribute the budget", which generalises to watches
// with different LED counts.
//
// Intuition:
//
//	If the hour shows `hb` lit LEDs then the minute must show `turnedOn - hb`.
//	So iterate hb from 0..min(turnedOn,4); for that hb collect every hour with
//	exactly hb set bits and every minute with exactly (turnedOn-hb) set bits,
//	then take their cross product.
//
// Algorithm:
//  1. For hb = 0 .. turnedOn (but no more than 4 hour LEDs):
//     a. mb = turnedOn - hb; skip if mb < 0 or mb > 6 (too many minute LEDs).
//     b. Gather hours h in 0..11 with popcount(h) == hb.
//     c. Gather minutes m in 0..59 with popcount(m) == mb.
//     d. Emit every "h:mm" combination.
//  2. Return the collected list.
//
// Time:  O(12·60) = O(1) — still bounded by the fixed hour/minute ranges.
// Space: O(1) plus the output.
func splitBudget(turnedOn int) []string {
	var result []string
	// hb = number of lit LEDs assigned to the hour. At most 4 hour LEDs exist,
	// and it obviously cannot exceed the total budget.
	for hb := 0; hb <= turnedOn && hb <= 4; hb++ {
		mb := turnedOn - hb // the rest of the budget goes to the minute
		if mb < 0 || mb > 6 {
			continue // impossible: only 6 minute LEDs are available
		}
		// Collect all hours whose popcount matches the hour budget.
		var hours []int
		for h := 0; h < 12; h++ {
			if bits.OnesCount(uint(h)) == hb {
				hours = append(hours, h)
			}
		}
		// Collect all minutes whose popcount matches the minute budget.
		var mins []int
		for m := 0; m < 60; m++ {
			if bits.OnesCount(uint(m)) == mb {
				mins = append(mins, m)
			}
		}
		// Cross product: every valid hour with every valid minute.
		for _, h := range hours {
			for _, m := range mins {
				result = append(result, fmt.Sprintf("%d:%02d", h, m))
			}
		}
	}
	return result
}

// ── Approach 3: Backtracking over the 10 LEDs (Optimal in spirit) ────────────
//
// backtracking treats the watch as 10 LEDs with fixed weights and chooses
// exactly `turnedOn` of them to switch on, pruning any partial choice that
// pushes the hour past 11 or the minute past 59.
//
// Intuition:
//
//	Model each LED as a (row, weight) pair: hours contribute {1,2,4,8}, minutes
//	contribute {1,2,4,8,16,32}. Pick LEDs left to right; each pick adds its
//	weight to either the hour or the minute total. When we have placed exactly
//	`turnedOn` LEDs, we have a concrete (hour, minute) — record it. Prune early
//	the instant hour>11 or minute>59, which is what makes this the "smart"
//	enumeration rather than a blind 720 scan.
//
// Algorithm:
//  1. Weights: hour LEDs indices 0..3 -> 1,2,4,8; minute LEDs 4..9 -> 1..32.
//  2. dfs(index, remaining, hour, minute):
//     - If hour>11 or minute>59: prune (return).
//     - If remaining==0: record "hour:minute" and return.
//     - For pick in index..9: light that LED (add to hour or minute) and recurse
//     with remaining-1, then move on (each combination chosen once).
//  3. Sort the results so output is deterministic.
//
// Time:  O(C(10, turnedOn)) subsets in the worst case, but pruning keeps it far
//
//	below that; effectively O(1) since 10 LEDs is a constant.
//
// Space: O(turnedOn) recursion depth plus the output list.
func backtracking(turnedOn int) []string {
	// Weight of each of the 10 LEDs. Indices 0..3 are hour bits, 4..9 minute.
	weights := []int{1, 2, 4, 8, 1, 2, 4, 8, 16, 32}

	var result []string
	var dfs func(index, remaining, hour, minute int)
	dfs = func(index, remaining, hour, minute int) {
		if hour > 11 || minute > 59 {
			return // invalid clock face — abandon this branch immediately
		}
		if remaining == 0 {
			// Exactly turnedOn LEDs are lit; emit the formatted time.
			result = append(result, fmt.Sprintf("%d:%02d", hour, minute))
			return
		}
		// Choose the next LED to light among the remaining LEDs (index..9),
		// which guarantees each subset of LEDs is generated exactly once.
		for i := index; i < 10; i++ {
			if i < 4 {
				// Hour LED: add its weight to the hour half.
				dfs(i+1, remaining-1, hour+weights[i], minute)
			} else {
				// Minute LED: add its weight to the minute half.
				dfs(i+1, remaining-1, hour, minute+weights[i])
			}
		}
	}
	dfs(0, turnedOn, 0, 0)

	// The recursion visits LEDs in a fixed order, so results are already grouped,
	// but sorting makes the output stable and easy to compare in tests.
	sort.Strings(result)
	return result
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Printf("turnedOn=1 -> %v\n", bruteForce(1)) // expected [0:01 0:02 0:04 0:08 0:16 0:32 1:00 2:00 4:00 8:00]
	fmt.Printf("turnedOn=9 -> %v\n", bruteForce(9)) // expected []

	fmt.Println("=== Approach 2: Split the Budget ===")
	fmt.Printf("turnedOn=1 -> %v\n", splitBudget(1)) // expected [0:01 0:02 0:04 0:08 0:16 0:32 1:00 2:00 4:00 8:00]
	fmt.Printf("turnedOn=9 -> %v\n", splitBudget(9)) // expected []

	fmt.Println("=== Approach 3: Backtracking ===")
	fmt.Printf("turnedOn=1 -> %v\n", backtracking(1)) // expected [0:01 0:02 0:04 0:08 0:16 0:32 1:00 2:00 4:00 8:00]
	fmt.Printf("turnedOn=9 -> %v\n", backtracking(9)) // expected []
}
