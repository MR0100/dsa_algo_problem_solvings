package main

import "fmt"

// ── Approach 1: Brute Force (all substrings) ─────────────────────────────────
//
// bruteForce solves Longest Substring with At Most K Distinct Characters by
// examining every substring and checking its distinct-character count.
//
// Intuition:
//
//	The direct definition: try every start i and every end j, count the
//	distinct characters in s[i..j], and keep the longest window whose count is
//	≤ k. No cleverness — a baseline to contrast with sliding window.
//
// Algorithm:
//  1. For each start i, build a frequency set while extending end j.
//  2. When the distinct count exceeds k, stop extending this start.
//  3. Track the maximum valid length.
//
// Time:  O(n^2) — n starts × up to n extensions, each an O(1) set update.
// Space: O(k) — the per-start distinct-character set (bounded by alphabet).
func bruteForce(s string, k int) int {
	if k == 0 {
		return 0 // no characters allowed → empty substring only
	}
	best := 0
	for i := 0; i < len(s); i++ {
		freq := map[byte]int{} // distinct chars in the current window s[i..j]
		for j := i; j < len(s); j++ {
			freq[s[j]]++ // extend the window by one character
			if len(freq) > k {
				break // too many distinct chars; no longer window from this i works
			}
			if j-i+1 > best {
				best = j - i + 1 // record a longer valid window
			}
		}
	}
	return best
}

// ── Approach 2: Sliding Window with Hash Map (Optimal) ───────────────────────
//
// slidingWindow solves it in one linear pass: grow the right edge, and whenever
// more than k distinct characters appear, shrink from the left until valid.
//
// Intuition:
//
//	Maintain a window [left, right] whose distinct-count never exceeds k. Push
//	right forward one char at a time. If the count exceeds k, advance left,
//	decrementing counts and dropping a character from the map when its count
//	hits 0, until the window is valid again. Every position enters and leaves
//	the window at most once → linear.
//
// Algorithm:
//  1. Expand right, incrementing freq[s[right]].
//  2. While len(freq) > k: decrement freq[s[left]], delete if 0, left++.
//  3. best = max(best, right-left+1).
//
// Time:  O(n) — each index is added and removed at most once.
// Space: O(k) — the map holds at most k+1 distinct characters.
func slidingWindow(s string, k int) int {
	if k == 0 {
		return 0 // vacuously no valid non-empty window
	}
	freq := map[byte]int{} // char → count within [left, right]
	best, left := 0, 0
	for right := 0; right < len(s); right++ {
		freq[s[right]]++ // include the new right character
		// Shrink until at most k distinct characters remain.
		for len(freq) > k {
			freq[s[left]]-- // one fewer occurrence of the leftmost char
			if freq[s[left]] == 0 {
				delete(freq, s[left]) // its last copy left the window
			}
			left++ // move the left edge inward
		}
		if right-left+1 > best {
			best = right - left + 1 // widest valid window so far
		}
	}
	return best
}

// ── Approach 3: Sliding Window with Fixed Array (Optimal, ASCII) ─────────────
//
// slidingWindowArray is the same linear window but replaces the hash map with a
// fixed [128]int counter and an explicit distinct counter — avoids hashing.
//
// Intuition:
//
//	For an ASCII alphabet we don't need a map: a 128-slot array of counts plus a
//	running `distinct` integer (incremented when a count goes 0→1, decremented
//	when it goes 1→0) tracks exactly what len(freq) did, with O(1) constant-time
//	updates and no allocation.
//
// Algorithm:
//  1. count[s[right]]++; if it became 1, distinct++.
//  2. While distinct > k: count[s[left]]--; if it became 0, distinct--; left++.
//  3. best = max(best, right-left+1).
//
// Time:  O(n) — single pass, O(1) per step.
// Space: O(1) — a fixed 128-entry counter.
func slidingWindowArray(s string, k int) int {
	if k == 0 {
		return 0
	}
	var count [128]int // ASCII code → occurrences within the window
	distinct := 0      // number of characters currently present (count > 0)
	best, left := 0, 0
	for right := 0; right < len(s); right++ {
		if count[s[right]] == 0 {
			distinct++ // a brand-new character entered the window
		}
		count[s[right]]++
		for distinct > k { // too many distinct → shrink from the left
			count[s[left]]--
			if count[s[left]] == 0 {
				distinct-- // the leftmost char's last copy left
			}
			left++
		}
		if right-left+1 > best {
			best = right - left + 1
		}
	}
	return best
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(bruteForce("eceba", 2)) // 3
	fmt.Println(bruteForce("aa", 1))    // 2

	fmt.Println("=== Approach 2: Sliding Window with Hash Map (Optimal) ===")
	fmt.Println(slidingWindow("eceba", 2)) // 3
	fmt.Println(slidingWindow("aa", 1))    // 2

	fmt.Println("=== Approach 3: Sliding Window with Fixed Array (Optimal) ===")
	fmt.Println(slidingWindowArray("eceba", 2)) // 3
	fmt.Println(slidingWindowArray("aa", 1))    // 2
}
