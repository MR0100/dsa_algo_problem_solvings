package main

import "fmt"

// ── Approach 1: Brute Force ──────────────────────────────────────────────────
//
// bruteForce solves Contains Duplicate II by checking every pair whose indices
// lie within k of each other.
//
// Intuition:
//
//	We need two EQUAL values at indices i and j with |i − j| <= k. For each i
//	we only have to look ahead at most k positions — beyond that the index
//	distance already exceeds k. If any such nearby pair is equal, return true.
//
// Algorithm:
//  1. For each i, scan j from i+1 up to min(i+k, n-1).
//  2. If nums[i] == nums[j], return true.
//  3. If nothing matches, return false.
//
// Time:  O(n·k) — each i inspects at most k neighbours.
// Space: O(1) — no auxiliary storage.
func bruteForce(nums []int, k int) bool {
	n := len(nums)
	for i := 0; i < n; i++ {
		// Only indices within k of i can satisfy |i-j| <= k.
		for j := i + 1; j <= i+k && j < n; j++ {
			if nums[i] == nums[j] { // equal values close enough in index
				return true
			}
		}
	}
	return false
}

// ── Approach 2: Hash Map of Last Index ───────────────────────────────────────
//
// hashMapLastIndex solves Contains Duplicate II by remembering, for each value,
// the most recent index at which it appeared.
//
// Intuition:
//
//	For a given value, the closest earlier occurrence is always its MOST
//	RECENT one. So we only need to store the last index per value. When we
//	meet a value again, compare the current index with that stored index; if
//	the gap is <= k we have a valid near-duplicate. Otherwise, update the
//	stored index to the current (closer) one and keep going.
//
// Algorithm:
//  1. Map value → last index seen.
//  2. For each i: if nums[i] is in the map and i − lastIndex <= k, return true.
//  3. Update map[nums[i]] = i (overwrite with the newest index).
//  4. If loop ends, return false.
//
// Time:  O(n) — one pass; O(1) average map operations.
// Space: O(n) — up to n distinct values stored.
func hashMapLastIndex(nums []int, k int) bool {
	lastIndex := make(map[int]int, len(nums)) // value → most recent index
	for i, v := range nums {
		if j, ok := lastIndex[v]; ok && i-j <= k { // seen before AND within k
			return true
		}
		lastIndex[v] = i // record/refresh the newest index of v
	}
	return false
}

// ── Approach 3: Sliding-Window Hash Set (Optimal) ────────────────────────────
//
// slidingWindowSet solves Contains Duplicate II by maintaining a set of the
// values in the last k indices and checking each new value against it.
//
// Intuition:
//
//	Keep a "window" set holding exactly the values at indices [i−k, i−1]. If
//	the current value nums[i] is already in that window, there is a duplicate
//	within distance k → true. As i advances, evict the value that falls out of
//	the window (the one at index i−k−1) so the set never exceeds size k.
//
// Algorithm:
//  1. Maintain a set `window` of the last k values.
//  2. For each i: if nums[i] ∈ window, return true.
//  3. Insert nums[i]; if the window now spans more than k indices, delete the
//     value at index i−k.
//  4. If loop ends, return false.
//
// Time:  O(n) — one pass; O(1) average set operations.
// Space: O(min(n, k)) — the window holds at most k values.
func slidingWindowSet(nums []int, k int) bool {
	window := make(map[int]struct{}) // values within the last k indices
	for i, v := range nums {
		if _, ok := window[v]; ok { // v already present in the k-window
			return true
		}
		window[v] = struct{}{} // add current value
		if len(window) > k {   // window grew beyond k values → evict oldest
			delete(window, nums[i-k]) // remove the value leaving the window
		}
	}
	return false
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(bruteForce([]int{1, 2, 3, 1}, 3))       // expected true
	fmt.Println(bruteForce([]int{1, 0, 1, 1}, 1))       // expected true
	fmt.Println(bruteForce([]int{1, 2, 3, 1, 2, 3}, 2)) // expected false

	fmt.Println("=== Approach 2: Hash Map of Last Index ===")
	fmt.Println(hashMapLastIndex([]int{1, 2, 3, 1}, 3))       // expected true
	fmt.Println(hashMapLastIndex([]int{1, 0, 1, 1}, 1))       // expected true
	fmt.Println(hashMapLastIndex([]int{1, 2, 3, 1, 2, 3}, 2)) // expected false

	fmt.Println("=== Approach 3: Sliding-Window Hash Set (Optimal) ===")
	fmt.Println(slidingWindowSet([]int{1, 2, 3, 1}, 3))       // expected true
	fmt.Println(slidingWindowSet([]int{1, 0, 1, 1}, 1))       // expected true
	fmt.Println(slidingWindowSet([]int{1, 2, 3, 1, 2, 3}, 2)) // expected false
}
