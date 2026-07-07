package main

import "fmt"

// ── Approach 1: Brute Force ──────────────────────────────────────────────────
//
// bruteForce solves Longest Substring with At Most Two Distinct Characters by
// checking every starting index and extending as far as possible.
//
// Intuition:
//
//	Try every substring start i. Extend j to the right, tracking the distinct
//	characters seen in s[i..j] with a small set. The moment a third distinct
//	character appears, no longer substring starting at i can be valid, so
//	break and move to the next start.
//
// Algorithm:
//  1. For each start index i:
//     a. Reset an empty character set.
//     b. For each end index j ≥ i: add s[j] to the set.
//     c. If the set exceeds 2 distinct chars → break (extension is hopeless).
//     d. Otherwise update best with the window length j-i+1.
//  2. Return best.
//
// Time:  O(n²) — n starts, each extending up to n characters.
// Space: O(1) — the set never holds more than 3 characters.
func bruteForce(s string) int {
	best := 0
	for i := 0; i < len(s); i++ {
		distinct := map[byte]bool{} // chars present in s[i..j]
		for j := i; j < len(s); j++ {
			distinct[s[j]] = true // include the new right endpoint
			if len(distinct) > 2 {
				break // a 3rd distinct char: every longer j is also invalid
			}
			if j-i+1 > best {
				best = j - i + 1 // valid window: record its length
			}
		}
	}
	return best
}

// ── Approach 2: Sliding Window + Frequency Map ───────────────────────────────
//
// slidingWindow solves Longest Substring with At Most Two Distinct Characters
// with the classic shrinking window over character counts.
//
// Intuition:
//
//	Maintain a window [left..right] that always contains ≤ 2 distinct chars.
//	Grow right one char at a time; if the window now has 3 distinct chars,
//	shrink from the left (decrementing counts, deleting keys that hit 0)
//	until only 2 remain. Every position enters and leaves the window at most
//	once → linear time.
//
// Algorithm:
//  1. count = empty map (char → occurrences inside the window), left = 0.
//  2. For right from 0..n-1:
//     a. count[s[right]]++.
//     b. While len(count) > 2: count[s[left]]--; if it hits 0 delete the key;
//     left++.
//     c. best = max(best, right-left+1)  (window is valid again).
//  3. Return best.
//
// Time:  O(n) — right moves n times, left moves at most n times total.
// Space: O(1) — the map holds at most 3 keys at any moment.
func slidingWindow(s string) int {
	count := map[byte]int{} // frequency of each char inside the window
	left, best := 0, 0
	for right := 0; right < len(s); right++ {
		count[s[right]]++ // absorb the new right endpoint
		for len(count) > 2 {
			count[s[left]]-- // evict one occurrence from the left edge
			if count[s[left]] == 0 {
				delete(count, s[left]) // char fully gone → distinct count drops
			}
			left++ // shrink the window
		}
		if right-left+1 > best {
			best = right - left + 1 // window is valid: record its length
		}
	}
	return best
}

// ── Approach 3: Sliding Window + Last-Occurrence Map (Optimal) ───────────────
//
// slidingWindowLastIndex solves Longest Substring with At Most Two Distinct
// Characters by jumping the left edge instead of shrinking one step at a time.
//
// Intuition:
//
//	Track each window char's LAST occurrence index. When a 3rd distinct char
//	arrives, the char that must be evicted is the one whose last occurrence
//	is furthest left — and the window can jump directly to (that index + 1).
//	The map holds ≤ 3 entries, so finding the minimum is O(1). This variant
//	generalises cleanly to "at most k distinct" (LeetCode #340).
//
// Algorithm:
//  1. last = empty map (char → index of its most recent occurrence), left = 0.
//  2. For right from 0..n-1:
//     a. last[s[right]] = right.
//     b. If len(last) > 2: find the entry with the smallest index, delete it,
//     and set left = that smallest index + 1.
//     c. best = max(best, right-left+1).
//  3. Return best.
//
// Time:  O(n) — each step does O(3) map work; left only jumps forward.
// Space: O(1) — at most 3 map entries.
func slidingWindowLastIndex(s string) int {
	last := map[byte]int{} // char → last index where it was seen
	left, best := 0, 0
	for right := 0; right < len(s); right++ {
		last[s[right]] = right // refresh (or insert) the newest position
		if len(last) > 2 {
			// locate the char whose most recent occurrence is leftmost
			minIdx := right
			var evict byte
			for c, idx := range last {
				if idx < minIdx {
					minIdx = idx
					evict = c
				}
			}
			delete(last, evict) // that char leaves the window entirely
			left = minIdx + 1   // window jumps past its final occurrence
		}
		if right-left+1 > best {
			best = right - left + 1 // record the current valid window
		}
	}
	return best
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Printf("s=%q  got=%d  expected 3\n", "eceba", bruteForce("eceba"))
	fmt.Printf("s=%q  got=%d  expected 5\n", "ccaabbb", bruteForce("ccaabbb"))

	fmt.Println("=== Approach 2: Sliding Window + Frequency Map ===")
	fmt.Printf("s=%q  got=%d  expected 3\n", "eceba", slidingWindow("eceba"))
	fmt.Printf("s=%q  got=%d  expected 5\n", "ccaabbb", slidingWindow("ccaabbb"))

	fmt.Println("=== Approach 3: Sliding Window + Last-Occurrence Map (Optimal) ===")
	fmt.Printf("s=%q  got=%d  expected 3\n", "eceba", slidingWindowLastIndex("eceba"))
	fmt.Printf("s=%q  got=%d  expected 5\n", "ccaabbb", slidingWindowLastIndex("ccaabbb"))

	// extra edge cases (all approaches must agree)
	fmt.Println("=== Edge Cases ===")
	fmt.Printf("s=%q  brute=%d window=%d jump=%d  expected 1,1,1\n",
		"a", bruteForce("a"), slidingWindow("a"), slidingWindowLastIndex("a"))
	fmt.Printf("s=%q  brute=%d window=%d jump=%d  expected 4,4,4\n",
		"abaccc", bruteForce("abaccc"), slidingWindow("abaccc"), slidingWindowLastIndex("abaccc"))
}
