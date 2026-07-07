package main

import "fmt"

// ── Approach 1: Sorting ───────────────────────────────────────────────────────
//
// sortingApproach solves First Missing Positive by sorting the array and then
// scanning for the first missing positive integer.
//
// Intuition: After sorting, positive integers appear in order. Walk and track
// the expected next integer starting from 1.
//
// Time:  O(n log n) — dominated by the sort
// Space: O(1) — in-place sort (or O(log n) for the sort stack)
func sortingApproach(nums []int) int {
	sortInts(nums)
	expected := 1
	for _, v := range nums {
		if v == expected {
			expected++
		}
	}
	return expected
}

func sortInts(a []int) {
	for i := 1; i < len(a); i++ {
		for j := i; j > 0 && a[j] < a[j-1]; j-- {
			a[j], a[j-1] = a[j-1], a[j]
		}
	}
}

// ── Approach 2: Hash Set ──────────────────────────────────────────────────────
//
// hashSet solves First Missing Positive using a set for O(n) average lookup.
//
// Time:  O(n)
// Space: O(n) — the set
func hashSet(nums []int) int {
	s := make(map[int]bool, len(nums))
	for _, v := range nums {
		s[v] = true
	}
	for i := 1; ; i++ {
		if !s[i] {
			return i
		}
	}
}

// ── Approach 3: Index as Hash (Cyclic Sort Variant) — Optimal ────────────────
//
// indexAsHash solves First Missing Positive in O(n) time and O(1) space using
// the input array itself as a hash map.
//
// Intuition: The answer must be in the range [1, n+1] (where n = len(nums)).
//   If all integers 1..n are present, the answer is n+1.
//   Otherwise some integer in 1..n is missing, and that's the answer.
// So we only care about values in [1, n].
//
// We use the array indices as a hash: place value v at index v-1 (when v ∈ [1,n]).
// After placement, scan for the first index i where nums[i] != i+1; return i+1.
//
// Algorithm:
//  1. For i = 0 to n-1:
//     while nums[i] is in [1,n] and nums[nums[i]-1] != nums[i]:
//       swap nums[i] with nums[nums[i]-1]  // place nums[i] at its correct index
//  2. For i = 0 to n-1:
//     if nums[i] != i+1: return i+1
//  3. Return n+1
//
// Time:  O(n) — each element is swapped at most once (moves to its final slot)
// Space: O(1)
func indexAsHash(nums []int) int {
	n := len(nums)

	// place each value v in [1,n] at index v-1
	for i := 0; i < n; i++ {
		// keep swapping until nums[i] is out-of-range or already in place
		for nums[i] > 0 && nums[i] <= n && nums[nums[i]-1] != nums[i] {
			dest := nums[i] - 1
			nums[i], nums[dest] = nums[dest], nums[i] // swap to correct position
		}
	}

	// find first index where the value is wrong
	for i := 0; i < n; i++ {
		if nums[i] != i+1 {
			return i + 1
		}
	}
	return n + 1 // all 1..n present; answer is n+1
}

func main() {
	fmt.Println("=== Approach 1: Sorting ===")
	fmt.Printf("nums=[1,2,0]     => %d  expected 3\n", sortingApproach([]int{1, 2, 0}))
	fmt.Printf("nums=[3,4,-1,1]  => %d  expected 2\n", sortingApproach([]int{3, 4, -1, 1}))
	fmt.Printf("nums=[7,8,9,11,12] => %d  expected 1\n", sortingApproach([]int{7, 8, 9, 11, 12}))

	fmt.Println("\n=== Approach 2: Hash Set ===")
	fmt.Printf("nums=[1,2,0]     => %d  expected 3\n", hashSet([]int{1, 2, 0}))
	fmt.Printf("nums=[3,4,-1,1]  => %d  expected 2\n", hashSet([]int{3, 4, -1, 1}))
	fmt.Printf("nums=[7,8,9,11,12] => %d  expected 1\n", hashSet([]int{7, 8, 9, 11, 12}))

	fmt.Println("\n=== Approach 3: Index as Hash (Optimal) ===")
	fmt.Printf("nums=[1,2,0]     => %d  expected 3\n", indexAsHash([]int{1, 2, 0}))
	fmt.Printf("nums=[3,4,-1,1]  => %d  expected 2\n", indexAsHash([]int{3, 4, -1, 1}))
	fmt.Printf("nums=[7,8,9,11,12] => %d  expected 1\n", indexAsHash([]int{7, 8, 9, 11, 12}))
	fmt.Printf("nums=[1]         => %d  expected 2\n", indexAsHash([]int{1}))
	fmt.Printf("nums=[1,2,3]     => %d  expected 4\n", indexAsHash([]int{1, 2, 3}))
}
