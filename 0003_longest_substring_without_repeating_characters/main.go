package main

import "fmt"

// ── Approach 1: Brute Force ───────────────────────────────────────────────────
//
// bruteForce checks every possible substring and tests each for uniqueness.
//
// Intuition:
//   Generate all substrings starting at every index i and ending at every
//   index j >= i. For each, check whether all characters are distinct using
//   a set. Track the maximum length seen.
//
// Algorithm:
//   1. For each start index i (0 to n-1):
//   2.   For each end index j (i to n-1):
//   3.     If s[i..j] has all unique characters, update maxLen.
//
// Time:  O(n³) — O(n²) substrings × O(n) uniqueness check each.
// Space: O(min(n,a)) — the set holds at most min(substring length, alphabet) chars.
func bruteForce(s string) int {
	n := len(s)
	maxLen := 0

	for i := 0; i < n; i++ {
		seen := make(map[byte]bool)
		for j := i; j < n; j++ {
			if seen[s[j]] {
				// Duplicate found — this substring and all longer ones starting
				// at i are invalid. Break and try next start.
				break
			}
			seen[s[j]] = true
			// All chars from i to j are unique; check if this is the longest.
			if j-i+1 > maxLen {
				maxLen = j - i + 1
			}
		}
	}
	return maxLen
}

// ── Approach 2: Sliding Window with HashSet ───────────────────────────────────
//
// slidingWindowSet maintains a window [left, right] of unique characters using
// a set. When a duplicate enters from the right, shrink from the left until
// the duplicate is removed, then expand right again.
//
// Intuition:
//   Use a set to represent the current window's character inventory. The right
//   pointer expands the window; when a duplicate is found the left pointer
//   shrinks it one step at a time until the window is valid again.
//
// Algorithm:
//   1. left = 0, set = {}.
//   2. For right = 0 to n-1:
//        While s[right] in set:
//          remove s[left] from set, left++.   (shrink until no dup)
//        Add s[right] to set.
//        Update maxLen = max(maxLen, right-left+1).
//   3. Return maxLen.
//
// Time:  O(n) — left and right each move at most n steps (2n total moves).
// Space: O(min(n,a)) — the set holds at most min(n, alphabet size) chars.
func slidingWindowSet(s string) int {
	charSet := make(map[byte]bool)
	left := 0
	maxLen := 0

	for right := 0; right < len(s); right++ {
		// Shrink window from the left until s[right] is no longer a duplicate.
		for charSet[s[right]] {
			delete(charSet, s[left])
			left++
		}
		// Now s[right] is safe to add; window [left,right] is duplicate-free.
		charSet[s[right]] = true
		if right-left+1 > maxLen {
			maxLen = right - left + 1
		}
	}
	return maxLen
}

// ── Approach 3: Sliding Window with HashMap — Jump Left (Optimal) ─────────────
//
// slidingWindowMap stores the last-seen index of each character. When a
// duplicate is found, left jumps directly past the previous occurrence in O(1)
// instead of creeping one step at a time.
//
// Intuition:
//   The inner while-loop in Approach 2 can be replaced by a direct jump.
//   The map tells us exactly where the duplicate was last seen; we can set
//   left = lastSeen[ch] + 1 instantly. The guard `lastSeen[ch] >= left`
//   ensures we never move left backwards (a character seen before the current
//   window start is irrelevant).
//
// Algorithm:
//   1. lastSeen = {}, left = 0, maxLen = 0.
//   2. For right = 0 to n-1:
//        ch = s[right]
//        if ch in lastSeen AND lastSeen[ch] >= left:
//          left = lastSeen[ch] + 1    (jump past the duplicate)
//        lastSeen[ch] = right
//        maxLen = max(maxLen, right-left+1)
//   3. Return maxLen.
//
// Time:  O(n) — single pass; each character processed once.
// Space: O(min(n,a)) — map holds at most alphabet-size entries.
func slidingWindowMap(s string) int {
	lastSeen := make(map[byte]int) // char → last index it was seen at
	left := 0
	maxLen := 0

	for right := 0; right < len(s); right++ {
		ch := s[right]

		// If ch was seen inside the current window, jump left past it.
		if idx, ok := lastSeen[ch]; ok && idx >= left {
			left = idx + 1
		}

		lastSeen[ch] = right // update to the most recent position

		if right-left+1 > maxLen {
			maxLen = right - left + 1
		}
	}
	return maxLen
}

// ── Approach 4: Sliding Window with Fixed Array (ASCII optimisation) ──────────
//
// slidingWindowArray replaces the hash map with a fixed 128-element array
// (one slot per ASCII character). Array access is faster than map lookup
// because there is no hashing overhead and everything fits in a single
// cache line.
//
// Intuition:
//   Identical logic to Approach 3 but uses an int array indexed by byte value
//   instead of a map. Works when the character set is bounded (ASCII = 128,
//   extended ASCII = 256). Initialise all entries to -1 (meaning "not seen").
//
// Time:  O(n) — same as Approach 3.
// Space: O(1) — the 128-slot array is a fixed-size constant.
func slidingWindowArray(s string) int {
	// lastSeen[c] = last index where character c appeared; -1 = not seen yet.
	var lastSeen [128]int
	for i := range lastSeen {
		lastSeen[i] = -1
	}

	left := 0
	maxLen := 0

	for right := 0; right < len(s); right++ {
		ch := s[right]
		// If ch was seen inside the current window, jump left.
		if lastSeen[ch] >= left {
			left = lastSeen[ch] + 1
		}
		lastSeen[ch] = right
		if right-left+1 > maxLen {
			maxLen = right - left + 1
		}
	}
	return maxLen
}

func main() {
	examples := []struct {
		s      string
		expect int
	}{
		{"abcabcbb", 3},
		{"bbbbb", 1},
		{"pwwkew", 3},
		{"", 0},
		{" ", 1},
	}

	approaches := []struct {
		name string
		fn   func(string) int
	}{
		{"Approach 1: Brute Force              O(n³) T | O(min(n,a)) S", bruteForce},
		{"Approach 2: Sliding Window + Set     O(n)  T | O(min(n,a)) S", slidingWindowSet},
		{"Approach 3: Sliding Window + Map ✅  O(n)  T | O(min(n,a)) S", slidingWindowMap},
		{"Approach 4: Sliding Window + Array   O(n)  T | O(1)        S", slidingWindowArray},
	}

	for _, ex := range examples {
		fmt.Printf("s=%q  expect=%d\n", ex.s, ex.expect)
		for _, ap := range approaches {
			fmt.Printf("  %-60s → %d\n", ap.name, ap.fn(ex.s))
		}
		fmt.Println()
	}
}
