package main

import "fmt"

// ── Approach 1: Sort After Copy ───────────────────────────────────────────────
//
// mergeSimple solves Merge Sorted Array by copying nums2 into nums1 then sorting.
//
// Intuition: Just append and sort. Doesn't use the sorted property.
//
// Time:  O((m+n) log(m+n))
// Space: O(1) — in-place sort (ignoring sort internals).
func mergeSimple(nums1 []int, m int, nums2 []int, n int) {
	// copy nums2 into the tail of nums1
	copy(nums1[m:], nums2[:n])
	// insertion sort on the combined array of size m+n
	for i := 1; i < m+n; i++ {
		key := nums1[i]
		j := i - 1
		for j >= 0 && nums1[j] > key {
			nums1[j+1] = nums1[j]
			j--
		}
		nums1[j+1] = key
	}
}

// ── Approach 2: Three Pointers from the End (Optimal) ────────────────────────
//
// merge solves Merge Sorted Array by filling nums1 from the back.
//
// Intuition:
//   Writing from the back avoids overwriting elements we haven't yet compared.
//   Use three pointers: p1 starts at m-1 (last valid element of nums1), p2 at
//   n-1 (last element of nums2), and write pointer p at m+n-1 (last position).
//   At each step, write the larger of nums1[p1] and nums2[p2] at nums1[p].
//   If p2 < 0, nums2 is exhausted and nums1's prefix is already in place.
//
// Algorithm:
//   p1=m-1, p2=n-1, p=m+n-1
//   while p2 >= 0:
//     if p1 >= 0 && nums1[p1] > nums2[p2]:
//       nums1[p] = nums1[p1]; p1--
//     else:
//       nums1[p] = nums2[p2]; p2--
//     p--
//
// Time:  O(m+n)
// Space: O(1)
func merge(nums1 []int, m int, nums2 []int, n int) {
	p1 := m - 1   // pointer to last valid element of nums1
	p2 := n - 1   // pointer to last element of nums2
	p := m + n - 1 // write position

	for p2 >= 0 {
		if p1 >= 0 && nums1[p1] > nums2[p2] {
			nums1[p] = nums1[p1]
			p1--
		} else {
			nums1[p] = nums2[p2]
			p2--
		}
		p--
	}
	// if p2 < 0, nums2 is exhausted; remaining nums1 elements are already in place
}

func main() {
	fmt.Println("=== Approach 1: Sort After Copy ===")
	n1 := []int{1, 2, 3, 0, 0, 0}
	mergeSimple(n1, 3, []int{2, 5, 6}, 3)
	fmt.Printf("nums1=[1,2,3,_,_,_] m=3 nums2=[2,5,6] n=3  got=%v  expected [1 2 2 3 5 6]\n", n1)

	n2 := []int{1}
	mergeSimple(n2, 1, []int{}, 0)
	fmt.Printf("nums1=[1] m=1 nums2=[] n=0  got=%v  expected [1]\n", n2)

	n3 := []int{0}
	mergeSimple(n3, 0, []int{1}, 1)
	fmt.Printf("nums1=[0] m=0 nums2=[1] n=1  got=%v  expected [1]\n", n3)

	fmt.Println("=== Approach 2: Three Pointers from End ===")
	n4 := []int{1, 2, 3, 0, 0, 0}
	merge(n4, 3, []int{2, 5, 6}, 3)
	fmt.Printf("nums1=[1,2,3,_,_,_] m=3 nums2=[2,5,6] n=3  got=%v  expected [1 2 2 3 5 6]\n", n4)

	n5 := []int{1}
	merge(n5, 1, []int{}, 0)
	fmt.Printf("nums1=[1] m=1 nums2=[] n=0  got=%v  expected [1]\n", n5)

	n6 := []int{0}
	merge(n6, 0, []int{1}, 1)
	fmt.Printf("nums1=[0] m=0 nums2=[1] n=1  got=%v  expected [1]\n", n6)
}
