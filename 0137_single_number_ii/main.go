package main

import "fmt"

// ── Approach 1: Brute Force ──────────────────────────────────────────────────
//
// bruteForce solves Single Number II using nested-loop counting.
//
// Intuition:
//
//	For every element, count its occurrences across the whole array. The one
//	whose count is 1 (not 3) is the answer.
//
// Algorithm:
//  1. For each index i, scan the array counting matches of nums[i].
//  2. Return nums[i] when the count is exactly 1.
//
// Time:  O(n²) — full counting scan per candidate.
// Space: O(1) — counters only.
func bruteForce(nums []int) int {
	for i := 0; i < len(nums); i++ { // candidate element
		count := 0
		for j := 0; j < len(nums); j++ { // count occurrences everywhere
			if nums[j] == nums[i] {
				count++
			}
		}
		if count == 1 { // appears once → the single number
			return nums[i]
		}
	}
	return -1 // unreachable per problem guarantee
}

// ── Approach 2: Hash Map ─────────────────────────────────────────────────────
//
// hashMap solves Single Number II with a frequency table.
//
// Intuition:
//
//	Tally every value; the key with count 1 is the single number. Simple and
//	general, but uses O(n) extra memory (disallowed by the follow-up).
//
// Algorithm:
//  1. Build freq[num]++ over one pass.
//  2. Return the key whose count is 1.
//
// Time:  O(n) — two linear passes.
// Space: O(n) — map of distinct values.
func hashMap(nums []int) int {
	freq := make(map[int]int, len(nums)) // value → count
	for _, num := range nums {
		freq[num]++ // tally
	}
	for num, count := range freq {
		if count == 1 { // unique element
			return num
		}
	}
	return -1 // unreachable per problem guarantee
}

// ── Approach 3: Math (Set Sum) ───────────────────────────────────────────────
//
// mathSum solves Single Number II via the identity (3·sum(set) − sum(all)) / 2.
//
// Intuition:
//
//	If every distinct value appeared exactly three times, the total would be
//	3·sum(distinct). The single number contributes once instead of thrice, so
//	the total is short by exactly two copies of it:
//	  3·sum(set) − sum(all) = 2 · single  →  divide by 2.
//
// Algorithm:
//  1. One pass: accumulate sumAll; add each value to sumSet the first time it
//     is seen (tracked with a set).
//  2. Return (3*sumSet − sumAll) / 2.
//
// Time:  O(n) — single pass.
// Space: O(n) — the distinct-value set.
func mathSum(nums []int) int {
	seen := make(map[int]bool, len(nums)) // distinct values
	sumAll, sumSet := 0, 0
	for _, num := range nums {
		sumAll += num // total including triplicates
		if !seen[num] {
			seen[num] = true
			sumSet += num // each distinct value once
		}
	}
	return (3*sumSet - sumAll) / 2 // missing two copies of the single number
}

// ── Approach 4: Bit Counting ─────────────────────────────────────────────────
//
// bitCount solves Single Number II by counting set bits per position modulo 3.
//
// Intuition:
//
//	Look at each of the 32 bit positions independently. Every value appearing
//	three times contributes 0 or 3 to a position's set-bit count, so
//	count % 3 isolates the single number's bit at that position. Working in
//	int32 makes the sign bit reconstruct correctly for negative answers.
//
// Algorithm:
//  1. For each bit position 0..31, count how many nums have that bit set.
//  2. If count % 3 != 0, set that bit in the result.
//  3. Convert the assembled int32 back to int (restores the sign).
//
// Time:  O(32·n) = O(n) — 32 passes of n bit tests.
// Space: O(1) — a single 32-bit accumulator.
func bitCount(nums []int) int {
	var result int32 // assemble the answer bit by bit (int32 keeps the sign bit honest)
	for bit := 0; bit < 32; bit++ {
		count := 0
		for _, num := range nums {
			if (int32(num)>>bit)&1 == 1 { // is this bit set in num?
				count++
			}
		}
		if count%3 != 0 { // triples contribute multiples of 3; remainder = single's bit
			result |= int32(1) << bit // plant the bit (bit 31 wraps to the sign bit correctly)
		}
	}
	return int(result) // int32 → int preserves the two's-complement value
}

// ── Approach 5: Ones/Twos Bitmask DFA (Optimal) ──────────────────────────────
//
// onesTwos solves Single Number II with two bitmasks acting as a mod-3 counter.
//
// Intuition:
//
//	Track, per bit position, how many times (mod 3) that bit has appeared:
//	`ones` holds bits seen 1 time, `twos` bits seen 2 times. A third sighting
//	clears the bit from both — a tiny finite-state machine 00 → 01 → 10 → 00
//	run in parallel across all 32 bit lanes. After the pass, `ones` holds
//	exactly the bits of the number seen once.
//
// Algorithm:
//  1. ones = (ones ^ num) &^ twos — toggle into ones unless already in twos.
//  2. twos = (twos ^ num) &^ ones — toggle into twos unless it just re-entered ones.
//  3. Return ones.
//
// Time:  O(n) — one pass, O(1) work per element.
// Space: O(1) — two integer registers.
func onesTwos(nums []int) int {
	ones, twos := 0, 0
	for _, num := range nums {
		ones = (ones ^ num) &^ twos // advance state 00→01, 10→(blocked, cleared below)
		twos = (twos ^ num) &^ ones // advance state 01→10, 10→00
	}
	return ones // bits seen exactly once (mod 3) = the single number
}

func main() {
	examples := [][]int{
		{2, 2, 3, 2},           // expected 3
		{0, 1, 0, 1, 0, 1, 99}, // expected 99
	}

	fmt.Println("=== Approach 1: Brute Force ===")
	for _, ex := range examples {
		fmt.Printf("nums=%v  got=%d\n", ex, bruteForce(ex)) // expected 3, then 99
	}

	fmt.Println("=== Approach 2: Hash Map ===")
	for _, ex := range examples {
		fmt.Printf("nums=%v  got=%d\n", ex, hashMap(ex)) // expected 3, then 99
	}

	fmt.Println("=== Approach 3: Math (Set Sum) ===")
	for _, ex := range examples {
		fmt.Printf("nums=%v  got=%d\n", ex, mathSum(ex)) // expected 3, then 99
	}

	fmt.Println("=== Approach 4: Bit Counting ===")
	for _, ex := range examples {
		fmt.Printf("nums=%v  got=%d\n", ex, bitCount(ex)) // expected 3, then 99
	}
	// extra sanity check: negative single number exercises the sign bit
	fmt.Printf("nums=[-2 -2 1 -2]  got=%d\n", bitCount([]int{-2, -2, 1, -2}))   // expected 1
	fmt.Printf("nums=[-4 -4 -4 -7]  got=%d\n", bitCount([]int{-4, -4, -4, -7})) // expected -7

	fmt.Println("=== Approach 5: Ones/Twos Bitmask DFA (Optimal) ===")
	for _, ex := range examples {
		fmt.Printf("nums=%v  got=%d\n", ex, onesTwos(ex)) // expected 3, then 99
	}
	fmt.Printf("nums=[-4 -4 -4 -7]  got=%d\n", onesTwos([]int{-4, -4, -4, -7})) // expected -7
}
