package main

import "fmt"

// ── Approach 1: Count Sort ────────────────────────────────────────────────────
//
// countSort solves Sort Colors by counting 0s, 1s, and 2s, then overwriting.
//
// Intuition:
//   Since there are only 3 values, count each, then fill the array.
//
// Time:  O(n) — two passes.
// Space: O(1) — 3-element count array.
func countSort(nums []int) {
	count := [3]int{}
	for _, v := range nums {
		count[v]++
	}
	i := 0
	for color := 0; color <= 2; color++ {
		for j := 0; j < count[color]; j++ {
			nums[i] = color
			i++
		}
	}
}

// ── Approach 2: Dutch National Flag (One-Pass) ────────────────────────────────
//
// dutchFlag solves Sort Colors in one pass using the Dutch National Flag algorithm.
//
// Intuition:
//   Maintain three pointers:
//   - lo: boundary — everything left of lo is 0.
//   - hi: boundary — everything right of hi is 2.
//   - mid: current element being examined.
//
//   At each step:
//   - nums[mid]==0: swap(lo,mid); lo++; mid++. (0 goes to left section)
//   - nums[mid]==1: mid++. (1 stays in middle)
//   - nums[mid]==2: swap(mid,hi); hi--. (2 goes to right; DON'T advance mid)
//
//   We don't advance mid after swapping with hi because the newly swapped
//   element (which was at hi) hasn't been examined yet.
//
// Algorithm:
//   lo=0, mid=0, hi=n-1
//   while mid<=hi:
//     if 0: swap(lo,mid); lo++; mid++
//     if 1: mid++
//     if 2: swap(mid,hi); hi--
//
// Time:  O(n) — single pass; each element examined at most once.
// Space: O(1)
func dutchFlag(nums []int) {
	lo, mid, hi := 0, 0, len(nums)-1
	for mid <= hi {
		switch nums[mid] {
		case 0:
			nums[lo], nums[mid] = nums[mid], nums[lo]
			lo++
			mid++ // lo-1..mid-1 is now processed; move both forward
		case 1:
			mid++ // 1 is already in the correct middle section
		case 2:
			nums[mid], nums[hi] = nums[hi], nums[mid]
			hi-- // don't advance mid: the element swapped from hi is unexamined
		}
	}
}

func main() {
	fmt.Println("=== Approach 1: Count Sort ===")
	n1 := []int{2, 0, 2, 1, 1, 0}
	countSort(n1)
	fmt.Printf("got=%v  expected [0 0 1 1 2 2]\n", n1)
	n2 := []int{2, 0, 1}
	countSort(n2)
	fmt.Printf("got=%v  expected [0 1 2]\n", n2)

	fmt.Println("=== Approach 2: Dutch National Flag ===")
	n3 := []int{2, 0, 2, 1, 1, 0}
	dutchFlag(n3)
	fmt.Printf("got=%v  expected [0 0 1 1 2 2]\n", n3)
	n4 := []int{2, 0, 1}
	dutchFlag(n4)
	fmt.Printf("got=%v  expected [0 1 2]\n", n4)
	n5 := []int{0}
	dutchFlag(n5)
	fmt.Printf("got=%v  expected [0]\n", n5)
	n6 := []int{1, 2, 0, 1, 2}
	dutchFlag(n6)
	fmt.Printf("got=%v  expected [0 1 1 2 2]\n", n6)
}
