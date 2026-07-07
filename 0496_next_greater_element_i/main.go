package main

import "fmt"

// ── Approach 1: Brute Force ──────────────────────────────────────────────────
//
// bruteForce solves Next Greater Element I by, for every value in nums1,
// locating it inside nums2 and then scanning forward for the first larger
// element.
//
// Intuition:
//
//	The definition is literal: the "next greater element" of x is the first
//	number to the right of x's position in nums2 that is bigger than x. So for
//	each query value we find where it lives in nums2, then walk rightwards
//	until we see something larger (answer) or fall off the end (-1).
//
// Algorithm:
//  1. For each value v in nums1:
//     a. Find its index j in nums2 (values are unique, so it is unique).
//     b. From j+1 to the end, return the first nums2[k] > v.
//     c. If none, record -1.
//
// Time:  O(m·n) — for each of the m queries we may scan all n of nums2.
// Space: O(1) — ignoring the output slice, only loop counters.
func bruteForce(nums1 []int, nums2 []int) []int {
	result := make([]int, len(nums1)) // one answer slot per query
	for i, v := range nums1 {         // handle each query value independently
		result[i] = -1 // default: assume no greater element exists
		j := 0
		for j < len(nums2) && nums2[j] != v { // locate v inside nums2
			j++
		}
		// walk to the right of v looking for the first strictly-greater number
		for k := j + 1; k < len(nums2); k++ {
			if nums2[k] > v {
				result[i] = nums2[k] // first greater to the right wins
				break                // stop at the very first one
			}
		}
	}
	return result
}

// ── Approach 2: Monotonic Stack + Hash Map (Optimal) ─────────────────────────
//
// monotonicStack solves Next Greater Element I by pre-computing, in one pass
// over nums2, the next-greater element for EVERY value, then answering each
// nums1 query with an O(1) map lookup.
//
// Intuition:
//
//	Keep a stack of values from nums2 that are still "waiting" for their next
//	greater element — the stack stays strictly decreasing from bottom to top.
//	When a new number cur arrives, it is the answer for every stacked value it
//	exceeds: pop each such value and record cur as its next-greater. Because a
//	value is pushed once and popped once, the whole scan is linear. Store the
//	pairings in a map so nums1 (a subset of nums2) can be answered instantly.
//
// Algorithm:
//  1. nextGreater = empty map (value -> its next greater element).
//  2. stack = empty (holds nums2 values with no answer yet, decreasing).
//  3. For each cur in nums2:
//     - While stack non-empty AND cur > stack top: pop t, set nextGreater[t]=cur.
//     - Push cur.
//  4. Any values still on the stack have no greater element (leave them absent).
//  5. For each v in nums1: answer = nextGreater[v] if present else -1.
//
// Time:  O(m + n) — each nums2 value is pushed/popped once; each query is O(1).
// Space: O(n) — the map and stack together hold at most n entries.
func monotonicStack(nums1 []int, nums2 []int) []int {
	nextGreater := make(map[int]int, len(nums2)) // value -> next greater element
	stack := make([]int, 0, len(nums2))          // decreasing stack of "unanswered" values

	for _, cur := range nums2 {
		// cur resolves every smaller value sitting on top of the stack
		for len(stack) > 0 && cur > stack[len(stack)-1] {
			top := stack[len(stack)-1]   // the value that was waiting
			stack = stack[:len(stack)-1] // pop it
			nextGreater[top] = cur       // cur is its first greater-to-the-right
		}
		stack = append(stack, cur) // cur now waits for its own next greater
	}
	// values remaining on the stack never found a greater element — omit them

	result := make([]int, len(nums1))
	for i, v := range nums1 {
		if g, ok := nextGreater[v]; ok { // O(1) lookup for this query
			result[i] = g
		} else {
			result[i] = -1 // no greater element existed
		}
	}
	return result
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(bruteForce([]int{4, 1, 2}, []int{1, 3, 4, 2})) // expected [-1 3 -1]
	fmt.Println(bruteForce([]int{2, 4}, []int{1, 2, 3, 4}))    // expected [3 -1]

	fmt.Println("=== Approach 2: Monotonic Stack + Hash Map (Optimal) ===")
	fmt.Println(monotonicStack([]int{4, 1, 2}, []int{1, 3, 4, 2})) // expected [-1 3 -1]
	fmt.Println(monotonicStack([]int{2, 4}, []int{1, 2, 3, 4}))    // expected [3 -1]
}
