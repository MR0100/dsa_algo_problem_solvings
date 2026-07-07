package main

import "fmt"

// ── Approach 1: Brute Force (Check Every Substring) ──────────────────────────
//
// bruteForce solves Longest Repeating Character Replacement by examining every
// substring, counting its most frequent letter, and checking whether the rest
// can be converted within k replacements.
//
// Intuition:
//
//	A substring can be turned into all-same letters using at most k operations
//	exactly when the number of characters that are NOT the most frequent one is
//	≤ k, i.e. (length − maxFrequency) ≤ k. So enumerate every substring, find
//	its max letter frequency, test the condition, and keep the longest that
//	passes. Correct and obvious; just slow.
//
// Algorithm:
//  1. For each start i, extend end j, maintaining letter counts of s[i..j].
//  2. Let maxFreq be the highest count in the window.
//  3. If (windowLength − maxFreq) ≤ k, the window is convertible — update best.
//  4. Return best.
//
// Time:  O(n²) — n² substrings; counts are updated incrementally, and maxFreq
//
//	over 26 letters is O(26) per step (a constant).
//
// Space: O(26) = O(1) — a fixed letter-count array per start.
func bruteForce(s string, k int) int {
	best := 0
	for i := 0; i < len(s); i++ {
		var count [26]int // letter frequencies for windows starting at i
		maxFreq := 0      // most frequent letter's count in s[i..j]
		for j := i; j < len(s); j++ {
			count[s[j]-'A']++ // extend window to include s[j]
			if count[s[j]-'A'] > maxFreq {
				maxFreq = count[s[j]-'A'] // s[j] may be the new most frequent letter
			}
			windowLen := j - i + 1 // current substring length
			// Chars other than the majority must be replaced; need at most k.
			if windowLen-maxFreq <= k && windowLen > best {
				best = windowLen // this substring is convertible and longer
			}
		}
	}
	return best
}

// ── Approach 2: Sliding Window with Recount ──────────────────────────────────
//
// slidingWindowRecount solves the problem with a variable-size window that
// shrinks whenever it becomes invalid, recomputing the window's max letter
// frequency from the count array on each check.
//
// Intuition:
//
//	Grow a window by moving right. A window is valid while (len − maxFreq) ≤ k.
//	When adding s[right] breaks that, shrink from the left until it holds again.
//	The largest width ever attained is the answer. Recomputing maxFreq by
//	scanning the 26 counts keeps the logic transparent; it costs an O(26)
//	factor but is still linear in n overall.
//
// Algorithm:
//  1. left = 0; grow right across the string, incrementing count[s[right]].
//  2. Compute maxFreq = max over the 26 counts.
//  3. While (windowLen − maxFreq) > k: decrement count[s[left]], left++,
//     recompute maxFreq.
//  4. Track the maximum windowLen.
//
// Time:  O(26·n) = O(n) — each index enters/leaves the window once; each check
//
//	scans 26 counters.
//
// Space: O(26) = O(1).
func slidingWindowRecount(s string, k int) int {
	var count [26]int
	best := 0
	left := 0
	for right := 0; right < len(s); right++ {
		count[s[right]-'A']++ // include the new right-hand character

		// Recompute the window's most frequent letter count.
		maxFreq := 0
		for _, c := range count {
			if c > maxFreq {
				maxFreq = c
			}
		}

		// Shrink while too many chars would need replacing.
		for (right-left+1)-maxFreq > k {
			count[s[left]-'A']-- // drop the leftmost character
			left++
			maxFreq = 0 // window changed → recompute the majority count
			for _, c := range count {
				if c > maxFreq {
					maxFreq = c
				}
			}
		}

		if right-left+1 > best {
			best = right - left + 1 // widest valid window so far
		}
	}
	return best
}

// ── Approach 3: Sliding Window, Non-Decreasing maxFreq (Optimal) ─────────────
//
// slidingWindowOptimal solves the problem in one pass by never recomputing
// maxFreq downward: it only ever grows. The window slides right by exactly one
// whenever the current width becomes invalid, and its width — since it never
// shrinks — is itself the running answer.
//
// Intuition:
//
//	Key insight: we only care about the LONGEST valid window, so we can let
//	maxFreq be a high-water mark that never decreases. If the window is valid
//	(width − maxFreq ≤ k), extend it. If not, we move left by one to keep the
//	width from growing — but we deliberately do NOT lower maxFreq. This is safe
//	because a smaller window with a smaller majority could never beat a length
//	we've already recorded; the width only advances when a genuinely better
//	(≥ maxFreq) majority appears. Thus `right − left + 1` is monotonically
//	non-decreasing and equals the answer.
//
// Algorithm:
//  1. left = 0, maxFreq = 0.
//  2. For each right: count[s[right]]++, maxFreq = max(maxFreq, that count).
//  3. If (windowLen − maxFreq) > k: count[s[left]]--, left++ (slide, don't shrink
//     maxFreq).
//  4. The answer is the final window width, len(s) − left.
//
// Time:  O(n) — single pass, O(1) work per character (no 26-scan).
// Space: O(26) = O(1).
func slidingWindowOptimal(s string, k int) int {
	var count [26]int
	maxFreq := 0 // high-water mark of any letter's count in the window; never decreases
	left := 0
	for right := 0; right < len(s); right++ {
		count[s[right]-'A']++ // add s[right] to the window
		if count[s[right]-'A'] > maxFreq {
			maxFreq = count[s[right]-'A'] // update the running majority count
		}

		// If more than k chars would need replacing, slide the window right by
		// one (keep its size fixed). We do not decrease maxFreq — a shorter
		// window can never improve on a length already achieved.
		if (right-left+1)-maxFreq > k {
			count[s[left]-'A']-- // remove the leftmost character
			left++
		}
	}
	// The window never shrank, so its final width is the longest valid length.
	return len(s) - left
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(bruteForce("ABAB", 2))    // expected 4
	fmt.Println(bruteForce("AABABBA", 1)) // expected 4

	fmt.Println("=== Approach 2: Sliding Window with Recount ===")
	fmt.Println(slidingWindowRecount("ABAB", 2))    // expected 4
	fmt.Println(slidingWindowRecount("AABABBA", 1)) // expected 4

	fmt.Println("=== Approach 3: Sliding Window, Non-Decreasing maxFreq (Optimal) ===")
	fmt.Println(slidingWindowOptimal("ABAB", 2))    // expected 4
	fmt.Println(slidingWindowOptimal("AABABBA", 1)) // expected 4
}
