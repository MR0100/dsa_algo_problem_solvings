package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Brute Force ──────────────────────────────────────────────────
//
// bruteForce solves Majority Element by counting the occurrences of every
// element with a nested loop.
//
// Intuition:
//
//	The majority element appears more than ⌊n/2⌋ times, so simply count each
//	candidate's occurrences and return the first one whose count crosses the
//	threshold.
//
// Algorithm:
//  1. For each element nums[i], scan the whole array counting equals.
//  2. If the count exceeds n/2, return nums[i].
//
// Time:  O(n^2) — a full counting pass for each of the n candidates.
// Space: O(1) — only two counters.
func bruteForce(nums []int) int {
	majorityCount := len(nums) / 2 // need strictly more than this many
	for _, candidate := range nums {
		count := 0
		for _, v := range nums {
			if v == candidate {
				count++ // tally occurrences of this candidate
			}
		}
		if count > majorityCount {
			return candidate // first value crossing the threshold wins
		}
	}
	return -1 // unreachable: a majority element always exists
}

// ── Approach 2: Hash Map Counting ────────────────────────────────────────────
//
// hashMap solves Majority Element with a single counting pass over a
// value → frequency map.
//
// Intuition:
//
//	One pass can build all counts at once. Since the majority element appears
//	more than n/2 times, we can even return the moment any count crosses the
//	threshold — no second pass needed.
//
// Algorithm:
//  1. Walk the array incrementing counts[v].
//  2. Return v as soon as counts[v] > n/2.
//
// Time:  O(n) — one pass with O(1) average map updates.
// Space: O(n) — up to n/2+1 distinct keys before the answer emerges.
func hashMap(nums []int) int {
	counts := map[int]int{} // value → occurrences seen so far
	majorityCount := len(nums) / 2
	for _, v := range nums {
		counts[v]++
		if counts[v] > majorityCount {
			return v // threshold crossed — must be the majority
		}
	}
	return -1 // unreachable: a majority element always exists
}

// ── Approach 3: Sorting ──────────────────────────────────────────────────────
//
// sorting solves Majority Element by sorting a copy and reading the middle.
//
// Intuition:
//
//	After sorting, equal values are contiguous. A block longer than n/2 must
//	cover the middle index ⌊n/2⌋ no matter where it starts, so the element
//	sitting there is guaranteed to be the majority.
//
// Algorithm:
//  1. Copy the array (keep the input intact) and sort the copy.
//  2. Return sorted[n/2].
//
// Time:  O(n log n) — dominated by the sort.
// Space: O(n) — the defensive copy (O(1) if sorting in place is allowed).
func sorting(nums []int) int {
	sorted := make([]int, len(nums))
	copy(sorted, nums) // don't mutate the caller's slice
	sort.Ints(sorted)
	// A run longer than n/2 must straddle the midpoint.
	return sorted[len(sorted)/2]
}

// ── Approach 4: Divide and Conquer ───────────────────────────────────────────
//
// divideAndConquer solves Majority Element by recursively finding the
// majority of each half and reconciling.
//
// Intuition:
//
//	If x is the majority of the whole range, x must be the majority of at
//	least one half (if it were a minority in both halves, its total count
//	couldn't exceed half the whole). So recurse on both halves; when the two
//	winners agree, done; when they disagree, count both across the full range
//	and keep the larger.
//
// Algorithm:
//  1. Base case: a 1-element range is its own majority.
//  2. Split at mid; recurse left and right.
//  3. If both halves agree, return that value.
//  4. Otherwise count each candidate in the full range and return the one
//     with more occurrences.
//
// Time:  O(n log n) — T(n) = 2T(n/2) + O(n) counting work per level.
// Space: O(log n) — recursion stack depth.
func divideAndConquer(nums []int) int {
	return majorityInRange(nums, 0, len(nums)-1)
}

// majorityInRange returns the majority element of nums[lo..hi] (inclusive).
func majorityInRange(nums []int, lo, hi int) int {
	// Base case: a single element is trivially the majority of its range.
	if lo == hi {
		return nums[lo]
	}
	mid := lo + (hi-lo)/2 // overflow-safe midpoint
	left := majorityInRange(nums, lo, mid)
	right := majorityInRange(nums, mid+1, hi)
	// Both halves elected the same value → it wins the merged range too.
	if left == right {
		return left
	}
	// Halves disagree → count both candidates over the whole range.
	if countInRange(nums, left, lo, hi) > countInRange(nums, right, lo, hi) {
		return left
	}
	return right
}

// countInRange counts occurrences of target in nums[lo..hi].
func countInRange(nums []int, target, lo, hi int) int {
	count := 0
	for i := lo; i <= hi; i++ {
		if nums[i] == target {
			count++
		}
	}
	return count
}

// ── Approach 5: Bit Manipulation ─────────────────────────────────────────────
//
// bitManipulation solves Majority Element by reconstructing it bit by bit.
//
// Intuition:
//
//	Look at any bit position: the majority element contributes its bit value
//	to more than n/2 of the numbers. So for each of the 32 bit positions,
//	whichever bit (0 or 1) appears in more than n/2 elements is the majority
//	element's bit at that position. Assemble all 32 winners into the answer.
//	Working on int32 makes the sign bit position 31, so negatives reconstruct
//	correctly.
//
// Algorithm:
//  1. For bit = 0..31: count elements whose int32 representation has that
//     bit set.
//  2. If the count exceeds n/2, set the bit in the int32 answer.
//  3. Return the assembled value widened back to int.
//
// Time:  O(32·n) = O(n) — one pass per bit position.
// Space: O(1) — a couple of counters and the accumulator.
func bitManipulation(nums []int) int {
	n := len(nums)
	var answer int32 // assemble in int32 so bit 31 doubles as the sign bit
	for bit := 0; bit < 32; bit++ {
		ones := 0
		for _, v := range nums {
			// Extract this bit from the two's-complement int32 form.
			if (int32(v)>>bit)&1 == 1 {
				ones++
			}
		}
		// The majority element dictates the majority bit at every position.
		if ones > n/2 {
			answer |= 1 << bit
		}
	}
	return int(answer) // widen back; sign extends automatically for negatives
}

// ── Approach 6: Boyer–Moore Voting (Optimal) ─────────────────────────────────
//
// boyerMooreVoting solves Majority Element in one pass with O(1) space —
// answering the follow-up.
//
// Intuition:
//
//	Maintain a candidate and a vote counter. Each equal element votes +1,
//	each different element votes -1; when the counter hits 0 the next
//	element becomes the new candidate. Pairing off one majority occurrence
//	against one non-majority occurrence can never exhaust the majority
//	element (it has > n/2 copies), so whoever survives as candidate at the
//	end must be the majority.
//
// Algorithm:
//  1. count = 0, candidate = undefined.
//  2. For each v: if count == 0, candidate = v. Then count += 1 if
//     v == candidate else count -= 1.
//  3. Return candidate (a verification pass would be required if a majority
//     were not guaranteed).
//
// Time:  O(n) — a single pass.
// Space: O(1) — one candidate and one counter.
func boyerMooreVoting(nums []int) int {
	count := 0
	candidate := 0
	for _, v := range nums {
		// Counter exhausted → previous candidate fully paired off; adopt v.
		if count == 0 {
			candidate = v
		}
		if v == candidate {
			count++ // a vote for the current candidate
		} else {
			count-- // a vote against — cancels one supporter
		}
	}
	// Majority is guaranteed to exist, so no verification pass is needed.
	return candidate
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Printf("nums=[3,2,3]          got=%d  expected 3\n", bruteForce([]int{3, 2, 3}))
	fmt.Printf("nums=[2,2,1,1,1,2,2]  got=%d  expected 2\n", bruteForce([]int{2, 2, 1, 1, 1, 2, 2}))

	fmt.Println("=== Approach 2: Hash Map Counting ===")
	fmt.Printf("nums=[3,2,3]          got=%d  expected 3\n", hashMap([]int{3, 2, 3}))
	fmt.Printf("nums=[2,2,1,1,1,2,2]  got=%d  expected 2\n", hashMap([]int{2, 2, 1, 1, 1, 2, 2}))

	fmt.Println("=== Approach 3: Sorting ===")
	fmt.Printf("nums=[3,2,3]          got=%d  expected 3\n", sorting([]int{3, 2, 3}))
	fmt.Printf("nums=[2,2,1,1,1,2,2]  got=%d  expected 2\n", sorting([]int{2, 2, 1, 1, 1, 2, 2}))

	fmt.Println("=== Approach 4: Divide and Conquer ===")
	fmt.Printf("nums=[3,2,3]          got=%d  expected 3\n", divideAndConquer([]int{3, 2, 3}))
	fmt.Printf("nums=[2,2,1,1,1,2,2]  got=%d  expected 2\n", divideAndConquer([]int{2, 2, 1, 1, 1, 2, 2}))

	fmt.Println("=== Approach 5: Bit Manipulation ===")
	fmt.Printf("nums=[3,2,3]          got=%d  expected 3\n", bitManipulation([]int{3, 2, 3}))
	fmt.Printf("nums=[2,2,1,1,1,2,2]  got=%d  expected 2\n", bitManipulation([]int{2, 2, 1, 1, 1, 2, 2}))
	fmt.Printf("nums=[-1,-1,2]        got=%d  expected -1\n", bitManipulation([]int{-1, -1, 2})) // negative majority edge

	fmt.Println("=== Approach 6: Boyer–Moore Voting (Optimal) ===")
	fmt.Printf("nums=[3,2,3]          got=%d  expected 3\n", boyerMooreVoting([]int{3, 2, 3}))
	fmt.Printf("nums=[2,2,1,1,1,2,2]  got=%d  expected 2\n", boyerMooreVoting([]int{2, 2, 1, 1, 1, 2, 2}))
	fmt.Printf("nums=[-1,-1,2]        got=%d  expected -1\n", boyerMooreVoting([]int{-1, -1, 2}))
}
