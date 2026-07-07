package main

import "fmt"

// ── Approach 1: Count-Based Write Pointer ────────────────────────────────────
//
// removeDuplicates solves Remove Duplicates II by tracking how many times the
// current element has appeared and using a write pointer.
//
// Intuition:
//   Walk with a read pointer. For each element, track how many consecutive
//   times it has been seen. Allow at most 2 occurrences; write it only if
//   count <= 2.
//
// Algorithm:
//   k=0, count=0, prev='∅'
//   for each num:
//     if num == prev: count++
//     else: count=1; prev=num
//     if count<=2: nums[k]=num; k++
//
// Time:  O(n)
// Space: O(1)
func removeDuplicatesCount(nums []int) int {
	k := 0
	count := 0
	var prev *int // nil means no previous element

	for _, num := range nums {
		if prev != nil && num == *prev {
			count++
		} else {
			count = 1
			v := num
			prev = &v
		}
		if count <= 2 {
			nums[k] = num
			k++
		}
	}
	return k
}

// ── Approach 2: Compare with k-2 (Elegant) ───────────────────────────────────
//
// removeDuplicates solves Remove Duplicates II using the elegant k-2 comparison.
//
// Intuition:
//   Keep a write pointer k (starts at 0). For each element, it should be
//   written if k < 2 (first two elements always go in) OR if nums[i] != nums[k-2]
//   (not equal to the element two positions back in the output).
//
//   If nums[i] == nums[k-2], it means the last two written elements are the
//   same value — adding another copy would create 3 consecutive duplicates.
//
// Algorithm:
//   k = 0
//   for each num:
//     if k < 2 or num != nums[k-2]: nums[k]=num; k++
//
// Time:  O(n)
// Space: O(1)
func removeDuplicates(nums []int) int {
	k := 0
	for _, num := range nums {
		if k < 2 || num != nums[k-2] {
			nums[k] = num
			k++
		}
	}
	return k
}

func main() {
	fmt.Println("=== Approach 1: Count-Based ===")
	n1 := []int{1, 1, 1, 2, 2, 3}
	k1 := removeDuplicatesCount(n1)
	fmt.Printf("nums=[1,1,1,2,2,3]  k=%d  nums[:%d]=%v  expected k=5 [1 1 2 2 3]\n", k1, k1, n1[:k1])

	n2 := []int{0, 0, 1, 1, 1, 1, 2, 3, 3}
	k2 := removeDuplicatesCount(n2)
	fmt.Printf("nums=[0,0,1,1,1,1,2,3,3]  k=%d  nums[:%d]=%v  expected k=7 [0 0 1 1 2 3 3]\n", k2, k2, n2[:k2])

	fmt.Println("=== Approach 2: Compare with k-2 ===")
	n3 := []int{1, 1, 1, 2, 2, 3}
	k3 := removeDuplicates(n3)
	fmt.Printf("nums=[1,1,1,2,2,3]  k=%d  nums[:%d]=%v  expected k=5 [1 1 2 2 3]\n", k3, k3, n3[:k3])

	n4 := []int{0, 0, 1, 1, 1, 1, 2, 3, 3}
	k4 := removeDuplicates(n4)
	fmt.Printf("nums=[0,0,1,1,1,1,2,3,3]  k=%d  nums[:%d]=%v  expected k=7 [0 0 1 1 2 3 3]\n", k4, k4, n4[:k4])
}
