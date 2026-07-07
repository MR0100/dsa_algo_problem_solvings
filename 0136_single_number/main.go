package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Brute Force ──────────────────────────────────────────────────
//
// bruteForce solves Single Number using nested-loop counting.
//
// Intuition:
//
//	For every element, count how many times it appears in the whole array.
//	The element whose count is exactly 1 is the answer.
//
// Algorithm:
//  1. For each index i, scan the entire array and count occurrences of nums[i].
//  2. If the count is 1, nums[i] is the single number — return it.
//
// Time:  O(n²) — for each of the n elements we scan all n elements again.
// Space: O(1) — only two loop counters.
func bruteForce(nums []int) int {
	for i := 0; i < len(nums); i++ { // candidate element
		count := 0
		for j := 0; j < len(nums); j++ { // count its occurrences everywhere
			if nums[j] == nums[i] {
				count++
			}
		}
		if count == 1 { // appears exactly once → it is the single number
			return nums[i]
		}
	}
	return -1 // unreachable per problem guarantee (a single number always exists)
}

// ── Approach 2: Hash Map ─────────────────────────────────────────────────────
//
// hashMap solves Single Number by counting occurrences in a map.
//
// Intuition:
//
//	One pass builds a frequency table; a second pass finds the key with
//	frequency 1. Trades O(n) memory for O(n) time.
//
// Algorithm:
//  1. Increment freq[num] for every num.
//  2. Return the key whose frequency is 1.
//
// Time:  O(n) — two linear passes with O(1) map operations.
// Space: O(n) — the map holds up to (n+1)/2 distinct keys.
func hashMap(nums []int) int {
	freq := make(map[int]int, len(nums)) // value → occurrence count
	for _, num := range nums {
		freq[num]++ // tally every element
	}
	for num, count := range freq {
		if count == 1 { // the unique element
			return num
		}
	}
	return -1 // unreachable per problem guarantee
}

// ── Approach 3: Sorting ──────────────────────────────────────────────────────
//
// sortAndScan solves Single Number by sorting and checking adjacent pairs.
//
// Intuition:
//
//	After sorting, every duplicate pair sits side by side. Walk the array two
//	at a time; the first pair that mismatches starts with the single number.
//
// Algorithm:
//  1. Sort a copy of the array (copy so the caller's slice is untouched).
//  2. Step i by 2: if arr[i] != arr[i+1], arr[i] broke the pairing — return it.
//  3. If every pair matched, the single number is the last element.
//
// Time:  O(n log n) — dominated by the sort.
// Space: O(n) — for the defensive copy (O(1) extra if sorting in place).
func sortAndScan(nums []int) int {
	arr := make([]int, len(nums)) // copy so we don't mutate the input
	copy(arr, nums)
	sort.Ints(arr) // duplicates become adjacent
	for i := 0; i+1 < len(arr); i += 2 {
		if arr[i] != arr[i+1] { // pairing broken → arr[i] is alone
			return arr[i]
		}
	}
	return arr[len(arr)-1] // all pairs matched → the last element is the single one
}

// ── Approach 4: Math (Set Sum) ───────────────────────────────────────────────
//
// mathSum solves Single Number using the identity 2·sum(set) − sum(nums).
//
// Intuition:
//
//	If every element appeared exactly twice, sum(nums) would equal
//	2·sum(distinct values). The single number appears once instead of twice,
//	so 2·sum(set) − sum(nums) = exactly the missing occurrence, i.e. the answer.
//
// Algorithm:
//  1. Collect distinct values in a set while accumulating sumAll and sumSet.
//  2. Return 2*sumSet − sumAll.
//
// Time:  O(n) — one pass.
// Space: O(n) — the set of distinct values.
func mathSum(nums []int) int {
	seen := make(map[int]bool, len(nums)) // distinct values
	sumAll, sumSet := 0, 0
	for _, num := range nums {
		sumAll += num // sum of everything, duplicates included
		if !seen[num] {
			seen[num] = true
			sumSet += num // sum of each distinct value once
		}
	}
	return 2*sumSet - sumAll // duplicates cancel, the single number remains
}

// ── Approach 5: Bitwise XOR (Optimal) ────────────────────────────────────────
//
// bitwiseXOR solves Single Number by XOR-folding the whole array.
//
// Intuition:
//
//	XOR is commutative/associative, a ⊕ a = 0, and a ⊕ 0 = a. XOR-ing all
//	elements cancels every pair to 0, leaving only the single number.
//
// Algorithm:
//  1. Start result at 0.
//  2. XOR every element into result.
//  3. result is the single number.
//
// Time:  O(n) — single pass.
// Space: O(1) — one accumulator, no extra memory.
func bitwiseXOR(nums []int) int {
	result := 0 // XOR identity
	for _, num := range nums {
		result ^= num // pairs cancel to 0; the loner survives
	}
	return result
}

func main() {
	examples := [][]int{
		{2, 2, 1},       // expected 1
		{4, 1, 2, 1, 2}, // expected 4
		{1},             // expected 1
	}
	expected := []int{1, 4, 1}

	fmt.Println("=== Approach 1: Brute Force ===")
	for i, ex := range examples {
		fmt.Printf("nums=%v  got=%d\n", ex, bruteForce(ex)) // expected 1, 4, 1
		_ = expected[i]
	}

	fmt.Println("=== Approach 2: Hash Map ===")
	for _, ex := range examples {
		fmt.Printf("nums=%v  got=%d\n", ex, hashMap(ex)) // expected 1, 4, 1
	}

	fmt.Println("=== Approach 3: Sorting ===")
	for _, ex := range examples {
		fmt.Printf("nums=%v  got=%d\n", ex, sortAndScan(ex)) // expected 1, 4, 1
	}

	fmt.Println("=== Approach 4: Math (Set Sum) ===")
	for _, ex := range examples {
		fmt.Printf("nums=%v  got=%d\n", ex, mathSum(ex)) // expected 1, 4, 1
	}

	fmt.Println("=== Approach 5: Bitwise XOR (Optimal) ===")
	for _, ex := range examples {
		fmt.Printf("nums=%v  got=%d\n", ex, bitwiseXOR(ex)) // expected 1, 4, 1
	}
}
