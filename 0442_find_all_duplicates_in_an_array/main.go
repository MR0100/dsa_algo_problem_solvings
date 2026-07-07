package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Brute Force (Pairwise Compare) ───────────────────────────────
//
// bruteForce solves Find All Duplicates in an Array by comparing every pair of
// elements.
//
// Intuition:
//
//	A value is a duplicate iff some later index holds the same value. Check all
//	pairs (i, j) with i < j; whenever nums[i] == nums[j], nums[i] appears twice.
//	Since each value appears at most twice, one match per value is enough.
//
// Algorithm:
//  1. For each i from 0..n-1, for each j from i+1..n-1:
//  2. If nums[i] == nums[j], append nums[i] to the result and break inner loop.
//
// Time:  O(n²) — all pairs.
// Space: O(1) auxiliary (ignoring the output slice).
func bruteForce(nums []int) []int {
	result := []int{}
	for i := 0; i < len(nums); i++ {
		// Look for a later copy of nums[i].
		for j := i + 1; j < len(nums); j++ {
			if nums[i] == nums[j] {
				result = append(result, nums[i]) // found the second occurrence
				break                            // at most twice — stop scanning
			}
		}
	}
	return result
}

// ── Approach 2: Hash Set / Counting ──────────────────────────────────────────
//
// hashSet solves Find All Duplicates in an Array by remembering values already
// seen.
//
// Intuition:
//
//	Walk the array once; the first time a value shows up it is new, the second
//	time it is a duplicate. A hash set of "seen" values captures exactly that.
//
// Algorithm:
//  1. Create an empty set `seen`.
//  2. For each value v: if v is already in `seen`, it is a duplicate (record it);
//     otherwise add v to `seen`.
//
// Time:  O(n) — one pass, O(1) expected set ops.
// Space: O(n) — the set may hold up to n/2 distinct values (violates the O(1)
//
//	follow-up, but is the natural first "linear" answer).
func hashSet(nums []int) []int {
	result := []int{}
	seen := make(map[int]struct{}, len(nums)) // membership set; struct{} = 0 bytes
	for _, v := range nums {
		if _, ok := seen[v]; ok {
			result = append(result, v) // second time we meet v → duplicate
		} else {
			seen[v] = struct{}{} // first time we meet v → remember it
		}
	}
	return result
}

// ── Approach 3: Negative Marking / Index-as-Hash (Optimal) ───────────────────
//
// negativeMarking solves Find All Duplicates in an Array in O(n) time and O(1)
// extra space by using the sign of nums[|v|-1] as a "visited" flag.
//
// Intuition:
//
//	Values are in [1, n], so each value v maps to a home index v-1. We use the
//	SIGN bit at that home index as a one-bit "have I seen v before?" marker.
//	First sighting of v: nums[v-1] is positive → flip it negative. Second
//	sighting: nums[v-1] is already negative → v is a duplicate. The magnitudes
//	are preserved, so we always read the true value via abs().
//
// Algorithm:
//  1. For each element, take v = |nums[i]| (magnitude, since earlier steps may
//     have negated it), and idx = v-1.
//  2. If nums[idx] < 0, v was seen before → append v to result.
//     Else negate nums[idx] to mark v as seen.
//
// Time:  O(n) — a single pass.
// Space: O(1) — mutates the input in place; only the output slice is extra.
//
// Note: this destroys the sign of the input. Callers who need nums intact must
// restore signs (or copy first). We copy in main() to keep runs independent.
func negativeMarking(nums []int) []int {
	result := []int{}
	for i := 0; i < len(nums); i++ {
		v := nums[i] // may be negative if this position was marked earlier
		if v < 0 {
			v = -v // recover the true (positive) value
		}
		idx := v - 1 // home index for value v (values are 1..n)
		if nums[idx] < 0 {
			// Home already flagged → this is the second time we see v.
			result = append(result, v)
		} else {
			// First sighting → flag the home index by making it negative.
			nums[idx] = -nums[idx]
		}
	}
	return result
}

// ── Approach 4: Cyclic Sort ──────────────────────────────────────────────────
//
// cyclicSort solves Find All Duplicates in an Array by placing each value at its
// home index, then reading off whatever is out of place.
//
// Intuition:
//
//	Because values are a near-permutation of 1..n, value v belongs at index v-1.
//	Repeatedly swap nums[i] toward its home. When the value that "should" be at a
//	home is already there (nums[i] == nums[nums[i]-1]) but sits at the wrong
//	index, it is the duplicate. After sorting, any index i whose value isn't i+1
//	holds a duplicate.
//
// Algorithm:
//  1. For each i: while nums[i] is not at its home (nums[i] != nums[nums[i]-1]),
//     swap it home.
//  2. Second pass: for each i, if nums[i] != i+1, nums[i] is a duplicate.
//
// Time:  O(n) — each swap puts one value in its final home, so total swaps ≤ n.
// Space: O(1) — in place (also mutates input; main() copies).
func cyclicSort(nums []int) []int {
	i := 0
	for i < len(nums) {
		home := nums[i] - 1 // where nums[i] wants to live
		// Swap toward home only if the home doesn't already hold this value.
		if nums[i] != nums[home] {
			nums[i], nums[home] = nums[home], nums[i]
		} else {
			i++ // either placed correctly or a duplicate blocks the home — advance
		}
	}
	result := []int{}
	// Anything not equal to index+1 is the extra copy squeezed out of place.
	for i := 0; i < len(nums); i++ {
		if nums[i] != i+1 {
			result = append(result, nums[i])
		}
	}
	return result
}

// copyInts returns a fresh copy so in-place approaches don't clobber shared data.
func copyInts(nums []int) []int {
	out := make([]int, len(nums))
	copy(out, nums)
	return out
}

// sortedCopy returns a sorted copy for order-independent comparison in prints.
func sortedCopy(nums []int) []int {
	out := copyInts(nums)
	sort.Ints(out)
	return out
}

func main() {
	ex1 := []int{4, 3, 2, 7, 8, 2, 3, 1} // expected {2,3}
	ex2 := []int{1, 1, 2}                // expected {1}
	ex3 := []int{}                       // expected {}

	fmt.Println("=== Approach 1: Brute Force (Pairwise Compare) ===")
	fmt.Printf("nums=[4 3 2 7 8 2 3 1]  got=%v  expected [2 3]\n", sortedCopy(bruteForce(copyInts(ex1))))
	fmt.Printf("nums=[1 1 2]            got=%v  expected [1]\n", sortedCopy(bruteForce(copyInts(ex2))))
	fmt.Printf("nums=[]                 got=%v  expected []\n", sortedCopy(bruteForce(copyInts(ex3))))

	fmt.Println("=== Approach 2: Hash Set / Counting ===")
	fmt.Printf("nums=[4 3 2 7 8 2 3 1]  got=%v  expected [2 3]\n", sortedCopy(hashSet(copyInts(ex1))))
	fmt.Printf("nums=[1 1 2]            got=%v  expected [1]\n", sortedCopy(hashSet(copyInts(ex2))))
	fmt.Printf("nums=[]                 got=%v  expected []\n", sortedCopy(hashSet(copyInts(ex3))))

	fmt.Println("=== Approach 3: Negative Marking (Optimal) ===")
	fmt.Printf("nums=[4 3 2 7 8 2 3 1]  got=%v  expected [2 3]\n", sortedCopy(negativeMarking(copyInts(ex1))))
	fmt.Printf("nums=[1 1 2]            got=%v  expected [1]\n", sortedCopy(negativeMarking(copyInts(ex2))))
	fmt.Printf("nums=[]                 got=%v  expected []\n", sortedCopy(negativeMarking(copyInts(ex3))))

	fmt.Println("=== Approach 4: Cyclic Sort ===")
	fmt.Printf("nums=[4 3 2 7 8 2 3 1]  got=%v  expected [2 3]\n", sortedCopy(cyclicSort(copyInts(ex1))))
	fmt.Printf("nums=[1 1 2]            got=%v  expected [1]\n", sortedCopy(cyclicSort(copyInts(ex2))))
	fmt.Printf("nums=[]                 got=%v  expected []\n", sortedCopy(cyclicSort(copyInts(ex3))))
}
