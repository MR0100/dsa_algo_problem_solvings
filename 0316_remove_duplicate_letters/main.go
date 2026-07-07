package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Brute Force Recursion (Smallest via Rightmost Split) ──────────
//
// bruteForce solves Remove Duplicate Letters by recursively choosing, at each
// step, the smallest possible leading character that still allows the rest of
// the distinct letters to appear afterwards.
//
// Intuition:
//
//	The answer starts with the SMALLEST character c such that, after the last
//	position where all remaining distinct characters are still available, c
//	occurs. Concretely: scan a prefix; the prefix must extend at least up to
//	the FIRST occurrence of the rarest letter (the letter whose last index is
//	smallest). Within that prefix window pick the smallest character, drop
//	everything before its chosen occurrence and all further copies of it, then
//	recurse on the remainder.
//
// Algorithm:
//  1. If s is empty, return "".
//  2. Find pos = min over all chars of (last index of that char). The answer's
//     first letter must be chosen at or before pos.
//  3. Among s[0..pos] pick the smallest character; let i be its first index.
//  4. Emit that character, then recurse on s[i+1:] with every copy of that
//     character removed.
//
// Time:  O(k · n) where k = number of distinct letters (≤ 26). Each recursion
//
//	level fixes one output letter and scans O(n); at most 26 levels.
//
// Space: O(n) — recursion depth and the filtered substrings.
func bruteForce(s string) string {
	if len(s) == 0 {
		return ""
	}
	// last[c] = last index where byte c appears in s.
	last := map[byte]int{}
	for i := 0; i < len(s); i++ {
		last[s[i]] = i // overwrite → ends up as final occurrence
	}
	// pos = the smallest "last index" among all present letters. The answer's
	// first char must be picked no later than here, else some letter vanishes.
	pos := len(s)
	for _, idx := range last {
		if idx < pos {
			pos = idx
		}
	}
	// Within window s[0..pos], choose the smallest char and its first index.
	best := byte('z' + 1) // sentinel larger than any lowercase letter
	bestIdx := 0
	for i := 0; i <= pos; i++ {
		if s[i] < best {
			best = s[i]
			bestIdx = i
		}
	}
	// Build the remainder: everything after bestIdx, with all copies of `best`
	// removed (we've committed to this letter already).
	rest := make([]byte, 0, len(s))
	for i := bestIdx + 1; i < len(s); i++ {
		if s[i] != best {
			rest = append(rest, s[i])
		}
	}
	return string(best) + bruteForce(string(rest)) // greedy pick + recurse
}

// ── Approach 2: Greedy with Counts (Explicit) ────────────────────────────────
//
// greedyCounts solves Remove Duplicate Letters using a boolean "in result" set
// plus remaining-count bookkeeping, building the answer left to right.
//
// Intuition:
//
//	Sweep the string keeping a result buffer. Before appending a character c,
//	pop trailing buffer characters that are LARGER than c provided they still
//	appear later (their remaining count > 0) — a later copy can represent them,
//	and moving them after c yields a lexicographically smaller string. Skip any
//	character already present in the buffer.
//
// Algorithm:
//  1. count[c] = total occurrences of each letter.
//  2. inResult[c] = whether c currently sits in the buffer.
//  3. For each char c: decrement count[c] (one fewer left to come).
//     - If c already in buffer, continue.
//     - While buffer non-empty AND top > c AND count[top] > 0: pop top and
//     mark it not-in-buffer (a later copy will re-add it).
//     - Push c, mark in buffer.
//  4. Join buffer.
//
// Time:  O(n) — each char pushed and popped at most once.
// Space: O(1) — fixed 26-size arrays plus a buffer ≤ 26.
func greedyCounts(s string) string {
	var count [26]int     // how many of each letter remain to be seen
	var inResult [26]bool // is this letter already in the stack?
	for i := 0; i < len(s); i++ {
		count[s[i]-'a']++ // tally every letter first
	}
	stack := make([]byte, 0, 26) // result buffer (monotonic-ish)
	for i := 0; i < len(s); i++ {
		c := s[i]
		count[c-'a']-- // we are now consuming this occurrence
		if inResult[c-'a'] {
			continue // already placed; keep the earlier (better) position
		}
		// Pop larger trailing letters that still occur later — a later copy
		// can stand in for them, and demoting them shrinks the result.
		for len(stack) > 0 {
			top := stack[len(stack)-1]
			if top > c && count[top-'a'] > 0 {
				stack = stack[:len(stack)-1]
				inResult[top-'a'] = false // it can come back later
			} else {
				break
			}
		}
		stack = append(stack, c)
		inResult[c-'a'] = true
	}
	return string(stack)
}

// ── Approach 3: Monotonic Stack with Last-Index (Optimal) ────────────────────
//
// monotonicStack solves Remove Duplicate Letters with a monotonic stack keyed
// on each letter's LAST occurrence index, avoiding a separate running count.
//
// Intuition:
//
//	Same greedy as Approach 2, but instead of remaining-counts we precompute
//	last[c] = final index of c. We may safely pop a larger top letter iff it
//	appears again later, i.e. its last index > current index i.
//
// Algorithm:
//  1. last[c] = final index of each letter.
//  2. seen[c] tracks buffer membership.
//  3. For i, c := range s:
//     - If seen[c], skip.
//     - While stack non-empty AND top > c AND last[top] > i: pop, unset seen.
//     - Push c, set seen.
//  4. Join stack.
//
// Time:  O(n) — one pass, amortized O(1) stack ops.
// Space: O(1) — 26-size arrays plus a bounded buffer.
func monotonicStack(s string) string {
	var last [26]int  // last index of each letter
	var seen [26]bool // membership in the stack
	for i := 0; i < len(s); i++ {
		last[s[i]-'a'] = i // final occurrence wins
	}
	stack := make([]byte, 0, 26)
	for i := 0; i < len(s); i++ {
		c := s[i]
		if seen[c-'a'] {
			continue // keep its first, already-optimal placement
		}
		// Pop a larger top only if it recurs after i (so we don't lose it).
		for len(stack) > 0 {
			top := stack[len(stack)-1]
			if top > c && last[top-'a'] > i {
				stack = stack[:len(stack)-1]
				seen[top-'a'] = false
			} else {
				break
			}
		}
		stack = append(stack, c)
		seen[c-'a'] = true
	}
	return string(stack)
}

// verifyDistinctSorted is a tiny helper for the demo output: it confirms a
// result contains exactly the distinct letters of the input.
func distinctLettersSorted(s string) string {
	set := map[rune]bool{}
	for _, r := range s {
		set[r] = true
	}
	out := make([]string, 0, len(set))
	for r := range set {
		out = append(out, string(r))
	}
	sort.Strings(out)
	res := ""
	for _, x := range out {
		res += x
	}
	return res
}

func main() {
	fmt.Println("=== Approach 1: Brute Force Recursion ===")
	fmt.Printf("bcabc      -> %q  expected \"abc\"\n", bruteForce("bcabc"))
	fmt.Printf("cbacdcbc   -> %q  expected \"acdb\"\n", bruteForce("cbacdcbc"))

	fmt.Println("=== Approach 2: Greedy with Counts ===")
	fmt.Printf("bcabc      -> %q  expected \"abc\"\n", greedyCounts("bcabc"))
	fmt.Printf("cbacdcbc   -> %q  expected \"acdb\"\n", greedyCounts("cbacdcbc"))

	fmt.Println("=== Approach 3: Monotonic Stack (Optimal) ===")
	fmt.Printf("bcabc      -> %q  expected \"abc\"\n", monotonicStack("bcabc"))
	fmt.Printf("cbacdcbc   -> %q  expected \"acdb\"\n", monotonicStack("cbacdcbc"))

	fmt.Println("=== Sanity: distinct letters preserved ===")
	fmt.Printf("cbacdcbc distinct sorted = %q, result sorted uses same set = %v\n",
		distinctLettersSorted("cbacdcbc"),
		distinctLettersSorted("cbacdcbc") == distinctLettersSorted(monotonicStack("cbacdcbc"))) // expected true
}
