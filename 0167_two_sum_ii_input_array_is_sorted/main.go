package main

import "fmt"

// ── Approach 1: Brute Force ──────────────────────────────────────────────────
//
// bruteForce solves Two Sum II by checking every pair of indices.
//
// Intuition:
//
//	Ignore the fact that the array is sorted and simply test every pair
//	(i, j) with i < j until one sums to target. Guaranteed to find the
//	unique answer, just slowly.
//
// Algorithm:
//  1. For every i from 0 to n-2:
//  2. For every j from i+1 to n-1:
//  3. If numbers[i] + numbers[j] == target, return the 1-indexed pair.
//
// Time:  O(n^2) — nested loops over all pairs.
// Space: O(1) — only loop counters.
func bruteForce(numbers []int, target int) []int {
	for i := 0; i < len(numbers)-1; i++ {
		for j := i + 1; j < len(numbers); j++ {
			// Exactly one solution exists, so return on the first hit.
			if numbers[i]+numbers[j] == target {
				return []int{i + 1, j + 1} // problem wants 1-indexed positions
			}
		}
	}
	return nil // unreachable: a solution is guaranteed
}

// ── Approach 2: Hash Map ─────────────────────────────────────────────────────
//
// hashMap solves Two Sum II exactly like classic Two Sum (#1): one pass with
// a value → index map.
//
// Intuition:
//
//	For each number x we need target-x. A map of already-seen values lets us
//	check for that complement in O(1). This ignores sortedness and violates
//	the O(1)-space requirement, but is the natural bridge from LeetCode #1
//	and worth knowing for comparison.
//
// Algorithm:
//  1. Walk the array left to right.
//  2. If target-numbers[i] is already in the map, return its stored index
//     and i (both +1 for 1-indexing).
//  3. Otherwise record numbers[i] → i and continue.
//
// Time:  O(n) — single pass, O(1) map operations.
// Space: O(n) — the map may hold nearly all elements.
func hashMap(numbers []int, target int) []int {
	seen := map[int]int{} // value → 0-based index where it was seen
	for i, x := range numbers {
		// Did an earlier element complete the pair?
		if j, ok := seen[target-x]; ok {
			return []int{j + 1, i + 1} // earlier index first, 1-indexed
		}
		seen[x] = i // remember this value for future complements
	}
	return nil // unreachable: a solution is guaranteed
}

// ── Approach 3: Binary Search ────────────────────────────────────────────────
//
// binarySearch solves Two Sum II by fixing each element and binary searching
// for its complement in the rest of the sorted array.
//
// Intuition:
//
//	Sortedness turns "find target-numbers[i]" into a binary search. For each
//	i, search the suffix (i+1..n-1) — searching only the suffix both keeps
//	O(1) space and prevents reusing the same element twice.
//
// Algorithm:
//  1. For each index i, compute need = target - numbers[i].
//  2. Binary search for need in numbers[i+1 .. n-1].
//  3. If found at index j, return [i+1, j+1].
//
// Time:  O(n log n) — n binary searches of O(log n) each.
// Space: O(1) — only search bounds.
func binarySearch(numbers []int, target int) []int {
	for i := 0; i < len(numbers)-1; i++ {
		need := target - numbers[i] // complement we must find
		lo, hi := i+1, len(numbers)-1
		for lo <= hi {
			mid := lo + (hi-lo)/2 // overflow-safe midpoint
			switch {
			case numbers[mid] == need:
				return []int{i + 1, mid + 1} // found the pair, 1-indexed
			case numbers[mid] < need:
				lo = mid + 1 // complement lies to the right
			default:
				hi = mid - 1 // complement lies to the left
			}
		}
	}
	return nil // unreachable: a solution is guaranteed
}

// ── Approach 4: Two Pointers (Optimal) ───────────────────────────────────────
//
// twoPointers solves Two Sum II with converging pointers from both ends.
//
// Intuition:
//
//	In a sorted array, the pair (left, right) has a tunable sum: moving left
//	rightward can only increase it, moving right leftward can only decrease
//	it. Compare the current sum with target and move the pointer that pushes
//	the sum toward target. Each step safely discards one element, because the
//	discarded element can't participate in any valid pair with the elements
//	still inside the window.
//
// Algorithm:
//  1. left = 0, right = n-1.
//  2. sum = numbers[left] + numbers[right].
//  3. If sum == target → return [left+1, right+1].
//  4. If sum < target  → left++ (need a bigger sum).
//  5. If sum > target  → right-- (need a smaller sum).
//  6. Repeat while left < right.
//
// Time:  O(n) — the pointers move at most n steps combined.
// Space: O(1) — two indices only (meets the problem's space requirement).
func twoPointers(numbers []int, target int) []int {
	left, right := 0, len(numbers)-1
	for left < right {
		sum := numbers[left] + numbers[right]
		switch {
		case sum == target:
			return []int{left + 1, right + 1} // 1-indexed answer
		case sum < target:
			left++ // sum too small → advance the small end
		default:
			right-- // sum too big → retreat the large end
		}
	}
	return nil // unreachable: a solution is guaranteed
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Printf("numbers=[2,7,11,15], target=9   got=%v  expected [1 2]\n", bruteForce([]int{2, 7, 11, 15}, 9))
	fmt.Printf("numbers=[2,3,4],     target=6   got=%v  expected [1 3]\n", bruteForce([]int{2, 3, 4}, 6))
	fmt.Printf("numbers=[-1,0],      target=-1  got=%v  expected [1 2]\n", bruteForce([]int{-1, 0}, -1))

	fmt.Println("=== Approach 2: Hash Map ===")
	fmt.Printf("numbers=[2,7,11,15], target=9   got=%v  expected [1 2]\n", hashMap([]int{2, 7, 11, 15}, 9))
	fmt.Printf("numbers=[2,3,4],     target=6   got=%v  expected [1 3]\n", hashMap([]int{2, 3, 4}, 6))
	fmt.Printf("numbers=[-1,0],      target=-1  got=%v  expected [1 2]\n", hashMap([]int{-1, 0}, -1))

	fmt.Println("=== Approach 3: Binary Search ===")
	fmt.Printf("numbers=[2,7,11,15], target=9   got=%v  expected [1 2]\n", binarySearch([]int{2, 7, 11, 15}, 9))
	fmt.Printf("numbers=[2,3,4],     target=6   got=%v  expected [1 3]\n", binarySearch([]int{2, 3, 4}, 6))
	fmt.Printf("numbers=[-1,0],      target=-1  got=%v  expected [1 2]\n", binarySearch([]int{-1, 0}, -1))

	fmt.Println("=== Approach 4: Two Pointers (Optimal) ===")
	fmt.Printf("numbers=[2,7,11,15], target=9   got=%v  expected [1 2]\n", twoPointers([]int{2, 7, 11, 15}, 9))
	fmt.Printf("numbers=[2,3,4],     target=6   got=%v  expected [1 3]\n", twoPointers([]int{2, 3, 4}, 6))
	fmt.Printf("numbers=[-1,0],      target=-1  got=%v  expected [1 2]\n", twoPointers([]int{-1, 0}, -1))
}
