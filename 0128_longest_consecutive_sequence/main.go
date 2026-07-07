package main

import "fmt"

// ── Approach 1: Sort ──────────────────────────────────────────────────────────
//
// longestConsecutiveBrute solves Longest Consecutive Sequence by sorting.
//
// Intuition:
//   After sorting, consecutive integers will be adjacent. Scan linearly,
//   extending the current sequence or starting a new one.
//
// Time:  O(n log n)
// Space: O(1) if in-place sort.
func longestConsecutiveBrute(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	// simple sort
	sortInts(nums)
	maxLen, curr := 1, 1
	for i := 1; i < len(nums); i++ {
		if nums[i] == nums[i-1] {
			continue // skip duplicate
		}
		if nums[i] == nums[i-1]+1 {
			curr++
		} else {
			curr = 1
		}
		if curr > maxLen {
			maxLen = curr
		}
	}
	return maxLen
}

func sortInts(a []int) {
	// quicksort
	var qs func(lo, hi int)
	qs = func(lo, hi int) {
		if lo >= hi {
			return
		}
		p := a[hi]
		i := lo
		for j := lo; j < hi; j++ {
			if a[j] <= p {
				a[i], a[j] = a[j], a[i]
				i++
			}
		}
		a[i], a[hi] = a[hi], a[i]
		qs(lo, i-1)
		qs(i+1, hi)
	}
	qs(0, len(a)-1)
}

// ── Approach 2: HashSet (Optimal) ────────────────────────────────────────────
//
// longestConsecutive solves Longest Consecutive Sequence in O(n).
//
// Intuition:
//   Store all numbers in a HashSet. For each number that is the START of a
//   sequence (num-1 not in set), count how far the sequence extends.
//   This ensures each number is visited at most twice overall → O(n).
//
// Time:  O(n)
// Space: O(n)
func longestConsecutive(nums []int) int {
	numSet := make(map[int]bool, len(nums))
	for _, n := range nums {
		numSet[n] = true
	}

	maxLen := 0
	for n := range numSet {
		if !numSet[n-1] { // n is the start of a sequence
			curr := n
			length := 1
			for numSet[curr+1] {
				curr++
				length++
			}
			if length > maxLen {
				maxLen = length
			}
		}
	}
	return maxLen
}

func main() {
	fmt.Println("=== Approach 1: Sort ===")
	fmt.Printf("nums=[100,4,200,1,3,2]  got=%d  expected 4\n", longestConsecutiveBrute([]int{100, 4, 200, 1, 3, 2}))
	fmt.Printf("nums=[0,3,7,2,5,8,4,6,0,1]  got=%d  expected 9\n", longestConsecutiveBrute([]int{0, 3, 7, 2, 5, 8, 4, 6, 0, 1}))

	fmt.Println("=== Approach 2: HashSet ===")
	fmt.Printf("nums=[100,4,200,1,3,2]  got=%d  expected 4\n", longestConsecutive([]int{100, 4, 200, 1, 3, 2}))
	fmt.Printf("nums=[0,3,7,2,5,8,4,6,0,1]  got=%d  expected 9\n", longestConsecutive([]int{0, 3, 7, 2, 5, 8, 4, 6, 0, 1}))
	fmt.Printf("nums=[]  got=%d  expected 0\n", longestConsecutive([]int{}))
}
