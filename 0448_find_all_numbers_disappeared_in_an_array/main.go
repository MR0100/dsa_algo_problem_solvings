package main

import "fmt"

// ── Approach 1: Brute Force (Seen Set) ───────────────────────────────────────
//
// bruteForce solves Find All Numbers Disappeared in an Array by recording which
// values appear in a boolean "seen" table and then collecting the unseen ones.
//
// Intuition:
//
//	Every value is in [1, n]. Mark each value present in a size-(n+1) boolean
//	array; afterwards, any index v in 1..n whose flag is still false never
//	appeared, so v is missing. Straightforward but uses O(n) auxiliary space.
//
// Algorithm:
//  1. seen := make([]bool, n+1); for each x in nums, seen[x] = true.
//  2. For v := 1..n: if !seen[v], append v to the result.
//  3. Return the result.
//
// Time:  O(n) — one pass to mark, one pass to collect.
// Space: O(n) — the boolean table (output not counted).
func bruteForce(nums []int) []int {
	n := len(nums)
	seen := make([]bool, n+1) // index 0 unused; values range over 1..n
	for _, x := range nums {
		seen[x] = true // this value is present
	}
	res := []int{}
	for v := 1; v <= n; v++ {
		if !seen[v] { // v in [1,n] but never marked → it is missing
			res = append(res, v)
		}
	}
	return res
}

// ── Approach 2: In-Place Negation Marking (Optimal) ──────────────────────────
//
// inPlaceMarking solves Find All Numbers Disappeared in an Array using the input
// array itself as the presence table, flipping signs to record which values
// occurred — O(1) extra space beyond the output.
//
// Intuition:
//
//	Values are in [1, n] and there are n slots, so value v maps to index v-1.
//	Walk the array; for each value v, negate nums[|v|-1] to stamp "value v was
//	seen" onto slot v-1. Use the absolute value because an earlier stamp may
//	have already flipped this element negative. After the pass, any slot i that
//	is still POSITIVE was never stamped, meaning value i+1 never appeared.
//	The signs are the only scratch space; the magnitudes stay intact for lookup.
//
// Algorithm:
//  1. For each element, let v = abs(nums[i]); if nums[v-1] > 0, negate it.
//  2. For each index i: if nums[i] > 0, value i+1 is missing → collect it.
//  3. Return the collected values. (Original signs could be restored if needed.)
//
// Time:  O(n) — two linear passes.
// Space: O(1) — mutates the input in place; only the output list is allocated.
func inPlaceMarking(nums []int) []int {
	n := len(nums)
	// Pass 1: stamp presence by making nums[value-1] negative.
	for i := 0; i < n; i++ {
		v := nums[i] // current value (may already be negated by a prior stamp)
		if v < 0 {
			v = -v // recover the true value 1..n
		}
		idx := v - 1 // the slot that represents value v
		if nums[idx] > 0 {
			nums[idx] = -nums[idx] // stamp: mark value v as seen
		}
	}
	// Pass 2: any still-positive slot i was never stamped → value i+1 missing.
	res := []int{}
	for i := 0; i < n; i++ {
		if nums[i] > 0 {
			res = append(res, i+1) // slot i encodes value i+1
		}
	}
	return res
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Seen Set) ===")
	fmt.Printf("nums=[4,3,2,7,8,2,3,1]  got=%v  expected [5 6]\n", bruteForce([]int{4, 3, 2, 7, 8, 2, 3, 1}))
	fmt.Printf("nums=[1,1]              got=%v  expected [2]\n", bruteForce([]int{1, 1}))
	fmt.Printf("nums=[1,2,3,4]          got=%v  expected []\n", bruteForce([]int{1, 2, 3, 4}))

	fmt.Println("=== Approach 2: In-Place Negation Marking (Optimal) ===")
	fmt.Printf("nums=[4,3,2,7,8,2,3,1]  got=%v  expected [5 6]\n", inPlaceMarking([]int{4, 3, 2, 7, 8, 2, 3, 1}))
	fmt.Printf("nums=[1,1]              got=%v  expected [2]\n", inPlaceMarking([]int{1, 1}))
	fmt.Printf("nums=[1,2,3,4]          got=%v  expected []\n", inPlaceMarking([]int{1, 2, 3, 4}))
}
