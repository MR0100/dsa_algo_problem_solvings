package main

import "fmt"

// ── Approach 1: Brute Force (All Substrings) ─────────────────────────────────
//
// bruteForce checks every substring and keeps the longest whose every distinct
// character occurs at least k times.
//
// Intuition:
//
//	Directly encode the definition. For each start i, extend end j, maintaining
//	a frequency count of the window s[i..j]. A window qualifies when all present
//	characters have count >= k. Track the max qualifying length.
//
// Algorithm:
//  1. For each start i: reset a 26-frequency array.
//  2. For each end j >= i: count[s[j]]++, then test if the whole window is valid
//     (every non-zero count >= k). If valid, update the answer.
//  3. Return the longest valid length.
//
// Time:  O(n² · 26) — O(n²) substrings, O(26) validity check each.
// Space: O(1) — a fixed 26-length count array.
func bruteForce(s string, k int) int {
	n := len(s)
	best := 0
	for i := 0; i < n; i++ {
		var count [26]int // fresh frequency table for windows starting at i
		for j := i; j < n; j++ {
			count[s[j]-'a']++ // extend the window to include s[j]
			if allAtLeastK(count, k) && j-i+1 > best {
				best = j - i + 1 // this window qualifies and is longer
			}
		}
	}
	return best
}

// allAtLeastK reports whether every character present (count > 0) occurs >= k times.
func allAtLeastK(count [26]int, k int) bool {
	for _, c := range count {
		if c > 0 && c < k {
			return false // some present char is too rare
		}
	}
	return true
}

// ── Approach 2: Divide and Conquer (on rare "splitter" characters) ───────────
//
// divideAndConquer uses the key insight that any character appearing fewer than
// k times can NEVER be part of a valid substring, so it acts as a splitter.
//
// Intuition:
//
//	Count characters in the segment. If every character already occurs >= k
//	times, the whole segment is valid — return its length. Otherwise pick a
//	character that appears but fewer than k times: no valid substring may
//	contain it, so the answer must lie entirely within one of the pieces
//	obtained by splitting the segment on every occurrence of that character.
//	Recurse into each piece and take the best.
//
// Algorithm:
//  1. Count frequencies in s[lo:hi].
//  2. If no character is present-but-rare (all >= k), return hi-lo.
//  3. Else pick such a splitter; walk the segment, and for each maximal piece
//     free of the splitter, recurse; return the max over pieces.
//
// Time:  O(n · 26) typical, O(n²) worst case (each level removes one char class).
// Space: O(26 · depth) recursion; depth <= 26.
func divideAndConquer(s string, k int) int {
	var solve func(lo, hi int) int
	solve = func(lo, hi int) int {
		if hi-lo < k {
			return 0 // too short to possibly satisfy the k requirement
		}
		// Count frequencies within s[lo:hi].
		var count [26]int
		for i := lo; i < hi; i++ {
			count[s[i]-'a']++
		}
		// Find a character that appears but fewer than k times: a splitter.
		splitter := byte(0)
		found := false
		for c := 0; c < 26; c++ {
			if count[c] > 0 && count[c] < k {
				splitter = byte('a' + c)
				found = true
				break
			}
		}
		// No splitter ⇒ every present char meets the threshold ⇒ whole segment valid.
		if !found {
			return hi - lo
		}
		// Split on every occurrence of the splitter and recurse into the pieces.
		best := 0
		start := lo
		for i := lo; i < hi; i++ {
			if s[i] == splitter {
				if r := solve(start, i); r > best {
					best = r
				}
				start = i + 1 // next piece begins just past the splitter
			}
		}
		// Trailing piece after the last splitter.
		if r := solve(start, hi); r > best {
			best = r
		}
		return best
	}
	return solve(0, len(s))
}

// ── Approach 3: Sliding Window over Fixed Unique Count (Optimal) ─────────────
//
// slidingWindow makes the problem monotone by fixing the number of distinct
// characters allowed in the window, turning it into a standard sliding window.
//
// Intuition:
//
//	A plain sliding window fails here because the "valid" condition is not
//	monotone as the window grows. Fix the target number of DISTINCT characters
//	the window may contain, from 1 to 26. For each fixed target `unique`, run a
//	sliding window that maintains exactly `unique` distinct chars, and track how
//	many of them already meet the >= k count. When distinct == unique AND every
//	one of them has count >= k, the window is valid — record its length. Because
//	the distinct count is now bounded/monotone, the window slides cleanly.
//
// Algorithm:
//  1. For unique = 1..26:
//     • Expand right pointer; on adding a char, update distinct count and the
//     "countAtLeastK" tally.
//     • While distinct > unique, shrink from the left, updating tallies.
//     • When distinct == unique and countAtLeastK == unique, update the answer.
//  2. Return the best length found.
//
// Time:  O(26 · n) — 26 passes, each a linear two-pointer sweep.
// Space: O(1) — a fixed 26-length frequency array per pass.
func slidingWindow(s string, k int) int {
	n := len(s)
	best := 0
	// Try every possible number of distinct characters in the answer window.
	for unique := 1; unique <= 26; unique++ {
		var count [26]int // frequency within the current window
		distinct := 0     // number of distinct chars currently in the window
		atLeastK := 0     // how many of those chars have count >= k
		left := 0
		for right := 0; right < n; right++ {
			// Add s[right] to the window.
			ri := s[right] - 'a'
			if count[ri] == 0 {
				distinct++ // a brand-new distinct character
			}
			count[ri]++
			if count[ri] == k {
				atLeastK++ // this char just reached the threshold
			}
			// Too many distinct chars: shrink from the left until distinct == unique.
			for distinct > unique {
				li := s[left] - 'a'
				if count[li] == k {
					atLeastK-- // this char is about to drop below k
				}
				count[li]--
				if count[li] == 0 {
					distinct-- // char fully left the window
				}
				left++
			}
			// Valid window: exactly `unique` distinct chars, all with count >= k.
			if distinct == unique && atLeastK == unique {
				if right-left+1 > best {
					best = right - left + 1
				}
			}
		}
	}
	return best
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Printf("s=\"aaabb\", k=3:   got=%d  expected 3\n", bruteForce("aaabb", 3))  // expected 3
	fmt.Printf("s=\"ababbc\", k=2:  got=%d  expected 5\n", bruteForce("ababbc", 2)) // expected 5

	fmt.Println("=== Approach 2: Divide and Conquer ===")
	fmt.Printf("s=\"aaabb\", k=3:   got=%d  expected 3\n", divideAndConquer("aaabb", 3))  // expected 3
	fmt.Printf("s=\"ababbc\", k=2:  got=%d  expected 5\n", divideAndConquer("ababbc", 2)) // expected 5

	fmt.Println("=== Approach 3: Sliding Window (Optimal) ===")
	fmt.Printf("s=\"aaabb\", k=3:   got=%d  expected 3\n", slidingWindow("aaabb", 3))  // expected 3
	fmt.Printf("s=\"ababbc\", k=2:  got=%d  expected 5\n", slidingWindow("ababbc", 2)) // expected 5
}
