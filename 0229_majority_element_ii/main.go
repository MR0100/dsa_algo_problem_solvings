package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Hash Map (Count and Filter) ──────────────────────────────────
//
// hashMapCount tallies every element's frequency, then returns those appearing
// more than n/3 times.
//
// Intuition:
//
//	"More than ⌊n/3⌋ times" is a direct frequency question. Count how many
//	times each value appears, then keep the values whose count clears the
//	threshold. There can be at most two such values (three distinct values each
//	exceeding n/3 would sum to more than n).
//
// Algorithm:
//  1. Build a map value → count over one pass.
//  2. threshold = n / 3 (integer division gives ⌊n/3⌋).
//  3. Collect every value whose count > threshold.
//  4. Sort the result for deterministic output.
//
// Time:  O(n) to count (+ O(1) sort of ≤2 elements).
// Space: O(n) — the frequency map holds up to n distinct keys.
func hashMapCount(nums []int) []int {
	counts := make(map[int]int) // value → number of occurrences
	for _, v := range nums {
		counts[v]++ // tally this value
	}
	threshold := len(nums) / 3 // must appear strictly MORE than ⌊n/3⌋ times
	res := []int{}
	for v, c := range counts {
		if c > threshold { // clears the majority-of-thirds bar
			res = append(res, v)
		}
	}
	sort.Ints(res) // map iteration order is random; sort for a stable answer
	return res
}

// ── Approach 2: Boyer–Moore Voting, Generalized (Optimal) ────────────────────
//
// boyerMooreVoting finds the up-to-two elements exceeding n/3 in O(1) extra
// space using two candidate/count pairs.
//
// Intuition:
//
//	At most two values can appear more than n/3 times. The classic Boyer–Moore
//	majority vote (for the >n/2 case) tracks one candidate; generalizing to
//	>n/k needs k-1 candidates. Here k=3, so we keep TWO candidates with two
//	independent counters. Each number either matches a candidate (increment its
//	count), fills an empty candidate slot, or "cancels" one vote from BOTH
//	candidates (decrementing both). Survivors are only POSSIBLE answers, so a
//	verification pass confirms each truly exceeds n/3.
//
// Algorithm:
//  1. Init cand1, cand2 with sentinel-distinct placeholders and count1=count2=0.
//  2. For each num:
//     - if it equals cand1 → count1++; else if equals cand2 → count2++;
//     - else if count1 == 0 → adopt num as cand1 (count1 = 1);
//     - else if count2 == 0 → adopt num as cand2 (count2 = 1);
//     - else → count1-- and count2-- (a three-way cancellation).
//  3. Recount cand1 and cand2 over the array; output those with count > n/3.
//
// Time:  O(n) — two linear passes.
// Space: O(1) — a fixed number of scalars.
func boyerMooreVoting(nums []int) []int {
	// Two candidate slots. Use distinct initial values so an all-same edge case
	// can't accidentally match an unset candidate; counts of 0 make them "empty".
	cand1, cand2 := 0, 1
	count1, count2 := 0, 0

	for _, v := range nums {
		switch {
		case v == cand1: // vote for candidate 1
			count1++
		case v == cand2: // vote for candidate 2
			count2++
		case count1 == 0: // slot 1 is empty → adopt v
			cand1, count1 = v, 1
		case count2 == 0: // slot 2 is empty → adopt v
			cand2, count2 = v, 1
		default: // v differs from both live candidates → cancel one vote each
			count1--
			count2--
		}
	}

	// The two survivors are only CANDIDATES — verify their true frequencies.
	count1, count2 = 0, 0
	for _, v := range nums {
		if v == cand1 {
			count1++
		} else if v == cand2 {
			count2++
		}
	}

	res := []int{}
	n := len(nums)
	if count1 > n/3 { // confirmed to exceed ⌊n/3⌋
		res = append(res, cand1)
	}
	if count2 > n/3 {
		res = append(res, cand2)
	}
	sort.Ints(res) // deterministic ordering
	return res
}

func main() {
	fmt.Println("=== Approach 1: Hash Map (Count and Filter) ===")
	fmt.Println(hashMapCount([]int{3, 2, 3})) // [3]
	fmt.Println(hashMapCount([]int{1}))       // [1]
	fmt.Println(hashMapCount([]int{1, 2}))    // [1 2]

	fmt.Println("=== Approach 2: Boyer–Moore Voting, Generalized (Optimal) ===")
	fmt.Println(boyerMooreVoting([]int{3, 2, 3})) // [3]
	fmt.Println(boyerMooreVoting([]int{1}))       // [1]
	fmt.Println(boyerMooreVoting([]int{1, 2}))    // [1 2]
}
