package main

import "fmt"

// ── Approach 1: Brute Force (Extra Space) ────────────────────────────────────
//
// bruteForce solves Remove Duplicates from Sorted Array using a temporary slice.
//
// Intuition: Collect unique elements into a new slice, then copy back. Simple
// but uses O(n) extra space — violates the in-place spirit of the problem.
//
// Algorithm:
//  1. Walk nums; append element if it differs from the last unique element.
//  2. Copy the unique slice back into nums[:k].
//  3. Return k = len(unique).
//
// Time:  O(n)
// Space: O(n) — extra slice
func bruteForce(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	unique := []int{nums[0]}
	for i := 1; i < len(nums); i++ {
		if nums[i] != unique[len(unique)-1] { // new unique value
			unique = append(unique, nums[i])
		}
	}
	// copy unique values back into nums
	for i, v := range unique {
		nums[i] = v
	}
	return len(unique)
}

// ── Approach 2: Two Pointers / Write Pointer (Optimal) ───────────────────────
//
// twoPointers solves Remove Duplicates from Sorted Array in-place with O(1) space.
//
// Intuition: Maintain a write pointer `k` that tracks where the next unique
// element should be placed. A read pointer `i` scans the array. Whenever
// nums[i] differs from nums[k-1] (the last unique element written), copy
// nums[i] to nums[k] and advance k.
//
// Algorithm:
//  1. k = 1 (first element is always unique).
//  2. For i = 1 to n-1:
//     if nums[i] != nums[k-1]: nums[k] = nums[i]; k++.
//  3. Return k.
//
// Time:  O(n) — single pass
// Space: O(1) — in-place, only two index variables
func twoPointers(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	k := 1 // nums[0] is always the first unique element
	for i := 1; i < len(nums); i++ {
		if nums[i] != nums[k-1] { // found a new unique element
			nums[k] = nums[i] // place it at the write position
			k++
		}
	}
	return k
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	nums1a := []int{1, 1, 2}
	k := bruteForce(nums1a)
	fmt.Printf("nums=[1,1,2]   k=%d  nums[:k]=%v  expected k=2, [1,2]\n", k, nums1a[:k])

	nums1b := []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}
	k = bruteForce(nums1b)
	fmt.Printf("nums=[0,0,1,1,1,2,2,3,3,4]  k=%d  nums[:k]=%v  expected k=5, [0,1,2,3,4]\n", k, nums1b[:k])

	fmt.Println("\n=== Approach 2: Two Pointers (Optimal) ===")
	nums2a := []int{1, 1, 2}
	k = twoPointers(nums2a)
	fmt.Printf("nums=[1,1,2]   k=%d  nums[:k]=%v  expected k=2, [1,2]\n", k, nums2a[:k])

	nums2b := []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}
	k = twoPointers(nums2b)
	fmt.Printf("nums=[0,0,1,1,1,2,2,3,3,4]  k=%d  nums[:k]=%v  expected k=5, [0,1,2,3,4]\n", k, nums2b[:k])
}
