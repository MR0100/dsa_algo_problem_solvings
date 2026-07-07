package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Brute Force (Sort Every Window) ──────────────────────────────
//
// bruteForce solves Find All Anagrams by taking every length-|p| window of s,
// sorting it, and comparing to sorted p.
//
// Intuition:
//
//	Two strings are anagrams iff their sorted forms are equal. So slide a window
//	of length m = len(p) across s, sort the window's characters, and test it
//	against the pre-sorted p. Correct and obvious, but sorting each of the ~n
//	windows costs O(m log m) apiece.
//
// Algorithm:
//  1. sortedP = sorted characters of p.
//  2. For each start i in [0, n-m]: take s[i:i+m], sort it, compare to sortedP.
//  3. Collect the starts that match.
//
// Time:  O((n - m) · m log m) — one sort per window.
// Space: O(m) — the per-window buffer.
func bruteForce(s string, p string) []int {
	n, m := len(s), len(p)
	res := []int{}
	if m > n {
		return res // p can't fit in s → no anagrams
	}
	sortedP := []byte(p)
	sort.Slice(sortedP, func(a, b int) bool { return sortedP[a] < sortedP[b] }) // canonical form of p

	for i := 0; i+m <= n; i++ {
		window := []byte(s[i : i+m])                                             // copy the current window
		sort.Slice(window, func(a, b int) bool { return window[a] < window[b] }) // canonicalise it
		if string(window) == string(sortedP) {                                   // same multiset of letters?
			res = append(res, i)
		}
	}
	return res
}

// ── Approach 2: Sliding Window with Count Comparison ─────────────────────────
//
// slidingWindowCompare solves Find All Anagrams by maintaining a 26-letter
// frequency array for the current window and comparing it, whole, to p's
// frequency array as the window slides one step at a time.
//
// Intuition:
//
//	Anagram == identical letter frequencies. Keep a rolling count of the current
//	window: add the incoming character, drop the outgoing one, and after each
//	slide check whether the window's 26-array equals p's 26-array. The window
//	itself is never re-sorted; only two counts change per step.
//
// Algorithm:
//  1. Build need[26] from p and win[26] from the first m characters of s.
//  2. If win == need, record start 0.
//  3. Slide i from m to n-1: win[s[i]]++ (enter), win[s[i-m]]-- (leave); if
//     win == need, record start i-m+1.
//
// Time:  O(n · 26) = O(n) — a constant 26-length array compare per position.
// Space: O(1) — two fixed 26-length arrays.
func slidingWindowCompare(s string, p string) []int {
	n, m := len(s), len(p)
	res := []int{}
	if m > n {
		return res
	}
	var need, win [26]int // fixed alphabet: 'a'..'z'
	for i := 0; i < m; i++ {
		need[p[i]-'a']++ // target letter frequencies
		win[s[i]-'a']++  // first window's letter frequencies
	}
	if win == need { // array value comparison in Go is a full element-wise check
		res = append(res, 0)
	}
	// Slide the window one character at a time.
	for i := m; i < n; i++ {
		win[s[i]-'a']++   // the new right-hand character enters the window
		win[s[i-m]-'a']-- // the leftmost character leaves the window
		if win == need {
			res = append(res, i-m+1) // window now covers s[i-m+1 .. i]
		}
	}
	return res
}

// ── Approach 3: Sliding Window with Match Counter (Optimal) ───────────────────
//
// slidingWindowMatchCounter solves Find All Anagrams without ever comparing the
// two 26-arrays: it tracks how many of the 26 letters currently have the exact
// required count, updating that tally in O(1) per slide.
//
// Intuition:
//
//	Re-scanning 26 slots each step is cheap but wasteful. Instead keep a single
//	`matches` counter = number of letters whose window-count already equals the
//	needed count. When a letter's count crosses INTO agreement, matches++; when
//	it crosses OUT, matches--. A window is an anagram exactly when matches == 26.
//	Each slide touches only two letters, so it adjusts matches in O(1).
//
// Algorithm:
//  1. need[26] from p. Initialise win[26]=0, matches = count of letters that
//     already agree (the 26 - (letters p uses) that are both zero).
//  2. Add first m chars via the same "adjust matches around equality" update.
//  3. Slide: add s[r], remove s[l]; whenever a letter's win count moves to/from
//     equal-to-need, bump matches. Record a hit when matches == 26 and the
//     window width is m.
//
// Time:  O(n) — O(1) work per character, no inner 26-loop after setup.
// Space: O(1) — two fixed arrays and a counter.
func slidingWindowMatchCounter(s string, p string) []int {
	n, m := len(s), len(p)
	res := []int{}
	if m > n {
		return res
	}
	var need, win [26]int
	for i := 0; i < m; i++ {
		need[p[i]-'a']++
	}
	matches := 0
	// Letters that p does NOT use start already satisfied (both counts 0).
	for c := 0; c < 26; c++ {
		if need[c] == 0 {
			matches++
		}
	}

	// add incorporates character c (index into window) and repairs `matches`.
	add := func(c int) {
		if win[c] == need[c] { // was equal → about to break equality
			matches--
		}
		win[c]++
		if win[c] == need[c] { // reached equality
			matches++
		}
	}
	// remove drops character c from the window and repairs `matches`.
	remove := func(c int) {
		if win[c] == need[c] { // was equal → about to break equality
			matches--
		}
		win[c]--
		if win[c] == need[c] { // reached equality
			matches++
		}
	}

	for r := 0; r < n; r++ {
		add(int(s[r] - 'a')) // extend window to the right
		if r >= m {
			remove(int(s[r-m] - 'a')) // shrink from the left to keep width m
		}
		if r >= m-1 && matches == 26 { // full-width window AND all 26 letters agree
			res = append(res, r-m+1)
		}
	}
	return res
}

// equal reports whether two int slices match — used to label expected output.
func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Sort Every Window) ===")
	fmt.Println(bruteForce("cbaebabacd", "abc"), equal(bruteForce("cbaebabacd", "abc"), []int{0, 6})) // [0 6] true
	fmt.Println(bruteForce("abab", "ab"), equal(bruteForce("abab", "ab"), []int{0, 1, 2}))            // [0 1 2] true

	fmt.Println("=== Approach 2: Sliding Window with Count Comparison ===")
	fmt.Println(slidingWindowCompare("cbaebabacd", "abc"), equal(slidingWindowCompare("cbaebabacd", "abc"), []int{0, 6})) // [0 6] true
	fmt.Println(slidingWindowCompare("abab", "ab"), equal(slidingWindowCompare("abab", "ab"), []int{0, 1, 2}))            // [0 1 2] true

	fmt.Println("=== Approach 3: Sliding Window with Match Counter (Optimal) ===")
	fmt.Println(slidingWindowMatchCounter("cbaebabacd", "abc"), equal(slidingWindowMatchCounter("cbaebabacd", "abc"), []int{0, 6})) // [0 6] true
	fmt.Println(slidingWindowMatchCounter("abab", "ab"), equal(slidingWindowMatchCounter("abab", "ab"), []int{0, 1, 2}))            // [0 1 2] true
}
