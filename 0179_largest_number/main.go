package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// toStrings converts the numbers to their decimal strings — every approach
// reasons about digit sequences, not numeric values.
func toStrings(nums []int) []string {
	strs := make([]string, len(nums))
	for i, v := range nums {
		strs[i] = strconv.Itoa(v)
	}
	return strs
}

// normalize collapses results like "00" to "0". If the largest arrangement
// starts with '0', every number must be 0 (a non-zero leader would have been
// placed first), so the whole answer is just "0".
func normalize(s string) string {
	if len(s) > 0 && s[0] == '0' {
		return "0" // all zeros — don't return "000...0"
	}
	return s
}

// ── Approach 1: Brute Force (All Permutations) ───────────────────────────────
//
// bruteForce solves Largest Number by trying every ordering of the numbers
// and keeping the largest concatenation.
//
// Intuition:
//
//	The answer IS some permutation of the input concatenated together, so
//	enumerate all n! permutations and take the maximum. Every candidate has
//	the same total number of digits, so plain lexicographic string comparison
//	equals numeric comparison — no big integers needed. Only viable for tiny
//	n, but it is the ground truth the clever approaches must match.
//
// Algorithm:
//  1. Convert all numbers to strings.
//  2. Recursively generate every permutation (swap-based backtracking).
//  3. Concatenate each permutation; keep the lexicographically largest.
//  4. Normalize the winner ("00" → "0") and return it.
//
// Time:  O(n! · n·k) — n! permutations, each joined in O(n·k) (k = max digits).
// Space: O(n·k) — recursion depth n plus the candidate strings.
func bruteForce(nums []int) string {
	strs := toStrings(nums)
	best := "" // lexicographic max so far (all candidates share one length)
	var permute func(k int)
	permute = func(k int) {
		// A full permutation is fixed — evaluate its concatenation.
		if k == len(strs) {
			cand := strings.Join(strs, "")
			// Same length ⇒ lexicographic comparison == numeric comparison.
			if cand > best {
				best = cand
			}
			return
		}
		for i := k; i < len(strs); i++ {
			strs[k], strs[i] = strs[i], strs[k] // choose strs[i] for slot k
			permute(k + 1)                      // permute the remaining slots
			strs[k], strs[i] = strs[i], strs[k] // undo the choice (backtrack)
		}
	}
	permute(0)
	return normalize(best)
}

// ── Approach 2: Greedy Selection ─────────────────────────────────────────────
//
// greedySelection solves Largest Number by repeatedly picking, for the next
// output slot, the number that concatenates best against all others.
//
// Intuition:
//
//	Build the answer left to right. For the next slot, the right choice is
//	the string s that "leads best": for every other remaining t, putting s
//	first is never worse (s+t >= t+s). Comparing the two concatenations
//	directly is the trick — comparing raw strings fails ("3" vs "30":
//	"330" > "303", so 3 must lead despite "3" < "30" lexicographically).
//	This is selection sort driven by the concatenation comparator.
//
// Algorithm:
//  1. Convert numbers to strings.
//  2. For each output slot pos = 0..n-1:
//     a. Scan the remaining strings; keep the candidate where
//     candidate+current > current+candidate (it leads better).
//     b. Swap the winner into position pos.
//  3. Join everything, normalize, and return.
//
// Time:  O(n² · k) — n² pairwise comparisons, each O(k) on ~2k-digit strings.
// Space: O(n·k) — the string copies of the numbers.
func greedySelection(nums []int) string {
	strs := toStrings(nums)
	for pos := 0; pos < len(strs); pos++ {
		bestIdx := pos // assume the current occupant leads best
		for i := pos + 1; i < len(strs); i++ {
			// Does strs[i] lead better than the current best? Compare the two
			// possible concatenations instead of the raw strings.
			if strs[i]+strs[bestIdx] > strs[bestIdx]+strs[i] {
				bestIdx = i
			}
		}
		// Place this slot's winner; the rest stay in the pool.
		strs[pos], strs[bestIdx] = strs[bestIdx], strs[pos]
	}
	return normalize(strings.Join(strs, ""))
}

// ── Approach 3: Custom Comparator Sort (Optimal) ─────────────────────────────
//
// customSort solves Largest Number by sorting all numbers with the
// "concatenation order": a before b iff a+b > b+a.
//
// Intuition:
//
//	If for two neighbours a+b < b+a, swapping them enlarges the result and
//	touches nothing else — so in the optimal arrangement every adjacent pair
//	satisfies a+b >= b+a. Sorting by exactly that relation produces such an
//	arrangement globally. The relation is a valid total order (transitive:
//	each string acts like a fixed "digit expansion" ratio; formally, if
//	a+b >= b+a and b+c >= c+b then a+c >= c+a), so sort.Slice is safe and
//	the greedy exchange argument proves the sorted order optimal.
//
// Algorithm:
//  1. Convert numbers to strings.
//  2. Sort descending under the comparator: strs[i]+strs[j] > strs[j]+strs[i].
//  3. Concatenate in sorted order.
//  4. Normalize the leading-zero case and return.
//
// Time:  O(n log n · k) — O(n log n) comparisons, each O(k) string work.
// Space: O(n·k) — the string slice (plus O(log n) sort recursion).
func customSort(nums []int) string {
	strs := toStrings(nums)
	sort.Slice(strs, func(i, j int) bool {
		// "i before j" exactly when i leading yields the bigger digit string.
		return strs[i]+strs[j] > strs[j]+strs[i]
	})
	return normalize(strings.Join(strs, ""))
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (All Permutations) ===")
	fmt.Printf("nums=[10,2]         got=%q  expected \"210\"\n", bruteForce([]int{10, 2}))
	fmt.Printf("nums=[3,30,34,5,9]  got=%q  expected \"9534330\"\n", bruteForce([]int{3, 30, 34, 5, 9}))
	fmt.Printf("nums=[0,0]          got=%q  expected \"0\"\n", bruteForce([]int{0, 0})) // leading-zero edge

	fmt.Println("=== Approach 2: Greedy Selection ===")
	fmt.Printf("nums=[10,2]         got=%q  expected \"210\"\n", greedySelection([]int{10, 2}))
	fmt.Printf("nums=[3,30,34,5,9]  got=%q  expected \"9534330\"\n", greedySelection([]int{3, 30, 34, 5, 9}))
	fmt.Printf("nums=[0,0]          got=%q  expected \"0\"\n", greedySelection([]int{0, 0}))

	fmt.Println("=== Approach 3: Custom Comparator Sort (Optimal) ===")
	fmt.Printf("nums=[10,2]         got=%q  expected \"210\"\n", customSort([]int{10, 2}))
	fmt.Printf("nums=[3,30,34,5,9]  got=%q  expected \"9534330\"\n", customSort([]int{3, 30, 34, 5, 9}))
	fmt.Printf("nums=[0,0]          got=%q  expected \"0\"\n", customSort([]int{0, 0}))
	fmt.Printf("nums=[432,43243]    got=%q  expected \"43243432\"\n", customSort([]int{432, 43243})) // prefix trap edge
}
