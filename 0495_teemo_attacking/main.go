package main

import "fmt"

// Each attack at time t poisons the inclusive interval [t, t+duration-1]. A new
// attack resets the timer. Total poisoned seconds = size of the UNION of all
// these intervals. Because timeSeries is sorted, consecutive intervals either
// overlap (next attack lands during the current poison) or are disjoint.

// ── Approach 1: Timeline Simulation (Brute Force) ────────────────────────────
//
// simulate marks every poisoned second on an explicit timeline and counts the
// marked seconds. It is the most literal reading: paint each interval, then
// measure the painted length.
//
// Intuition:
//
//	Lay out a boolean timeline covering every second that could be poisoned
//	(from the first attack to the last attack + duration). For each attack,
//	mark seconds [t, t+duration-1] as poisoned. Overlaps naturally coalesce
//	because marking an already-marked second is idempotent. The answer is the
//	number of marked seconds.
//
// Algorithm:
//  1. If timeSeries is empty or duration == 0, return 0.
//  2. Offset all times by the first attack so the timeline array starts at 0.
//  3. For each attack, set poisoned[t..t+duration-1] = true.
//  4. Count and return the number of true entries.
//
// Time:  O(n * duration) — each of n attacks paints up to `duration` seconds.
// Space: O(span) — a boolean per second of the poisoned window (can be large).
func simulate(timeSeries []int, duration int) int {
	if len(timeSeries) == 0 || duration == 0 {
		return 0
	}
	base := timeSeries[0]                                   // shift so the array is 0-indexed
	span := timeSeries[len(timeSeries)-1] - base + duration // total seconds to model
	poisoned := make([]bool, span)                          // poisoned[k] == true if second base+k is poisoned
	for _, t := range timeSeries {
		start := t - base // this attack's start, relative to base
		for k := start; k < start+duration; k++ {
			poisoned[k] = true // paint each second of [t, t+duration-1]
		}
	}
	count := 0
	for _, p := range poisoned {
		if p {
			count++ // tally the painted (poisoned) seconds
		}
	}
	return count
}

// ── Approach 2: Single-Pass Gap Sum (Optimal) ────────────────────────────────
//
// gapSum walks consecutive attacks once and adds, for each adjacent pair, the
// smaller of (gap between attacks) and (duration): a full `duration` if the
// intervals are disjoint, or only the non-overlapping prefix if they overlap.
//
// Intuition:
//
//	Consider attacks i and i+1. Attack i alone would poison `duration` seconds,
//	but it gets cut short if the next attack arrives first. The gap
//	timeSeries[i+1] - timeSeries[i] is how many fresh seconds attack i actually
//	contributes before being reset: if gap >= duration the interval finished
//	uninterrupted (add duration); if gap < duration the next attack reset the
//	timer early (add only gap). The very last attack always contributes a full
//	duration (nothing resets it). Sum = min(gap, duration) over adjacent pairs,
//	plus duration for the final attack.
//
// Algorithm:
//  1. If empty or duration == 0, return 0.
//  2. total = 0. For i in 0..n-2: total += min(timeSeries[i+1]-timeSeries[i], duration).
//  3. total += duration (the last attack, never interrupted).
//  4. Return total.
//
// Time:  O(n) — one pass over adjacent pairs.
// Space: O(1) — a running total.
func gapSum(timeSeries []int, duration int) int {
	if len(timeSeries) == 0 || duration == 0 {
		return 0
	}
	total := 0
	for i := 0; i+1 < len(timeSeries); i++ {
		gap := timeSeries[i+1] - timeSeries[i] // seconds until the next reset
		if gap < duration {
			total += gap // interrupted early: only `gap` fresh seconds counted
		} else {
			total += duration // uninterrupted: full poison duration
		}
	}
	total += duration // the final attack always runs its full duration
	return total
}

// ── Approach 3: Interval Union Merge ─────────────────────────────────────────
//
// mergeIntervals treats each attack as the interval [t, t+duration-1] and sums
// the length of their union by merging overlapping intervals on the fly.
//
// Intuition:
//
//	The answer is literally the measure of the union of intervals
//	[t_i, t_i + duration - 1]. Since timeSeries is sorted, sweep left→right
//	maintaining the current merged interval [start, end]. If the next interval
//	starts at or before end+1 it overlaps/touches — extend end. Otherwise the
//	current block is finished: add its length and open a new block. This is the
//	general "merge intervals then sum lengths" pattern, which also works if the
//	input were unsorted (after a sort).
//
// Algorithm:
//  1. If empty or duration == 0, return 0.
//  2. start = timeSeries[0], end = start + duration - 1, total = 0.
//  3. For each later attack t with interval [t, t+duration-1]:
//     - if t <= end (overlap), extend end = max(end, t+duration-1);
//     - else close the block: total += end-start+1; reset start,end to t's interval.
//  4. Add the final block's length. Return total.
//
// Time:  O(n) given sorted input (O(n log n) if a sort were needed).
// Space: O(1).
func mergeIntervals(timeSeries []int, duration int) int {
	if len(timeSeries) == 0 || duration == 0 {
		return 0
	}
	total := 0
	start := timeSeries[0]      // current merged block's start second
	end := start + duration - 1 // current merged block's end second (inclusive)
	for i := 1; i < len(timeSeries); i++ {
		t := timeSeries[i]
		newEnd := t + duration - 1 // interval contributed by this attack
		if t <= end {              // overlaps or touches the current block
			if newEnd > end {
				end = newEnd // extend the block to cover this attack
			}
		} else {
			total += end - start + 1 // close the finished block, tally its length
			start, end = t, newEnd   // begin a fresh block at this attack
		}
	}
	total += end - start + 1 // add the last open block
	return total
}

func main() {
	fmt.Println("=== Approach 1: Timeline Simulation (Brute Force) ===")
	fmt.Println(simulate([]int{1, 4}, 2)) // expected 4
	fmt.Println(simulate([]int{1, 2}, 2)) // expected 3

	fmt.Println("=== Approach 2: Single-Pass Gap Sum (Optimal) ===")
	fmt.Println(gapSum([]int{1, 4}, 2)) // expected 4
	fmt.Println(gapSum([]int{1, 2}, 2)) // expected 3

	fmt.Println("=== Approach 3: Interval Union Merge ===")
	fmt.Println(mergeIntervals([]int{1, 4}, 2)) // expected 4
	fmt.Println(mergeIntervals([]int{1, 2}, 2)) // expected 3
}
