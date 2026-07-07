package main

import (
	"fmt"
	"strings"
)

// ── Approach 1: Built-in Split + Filter (Brute Force) ────────────────────────
//
// builtinSplit counts segments by splitting the string on spaces and counting
// the non-empty pieces.
//
// Intuition:
//
//	A "segment" is a maximal run of non-space characters, i.e. a token. The
//	standard library's strings.Fields does exactly this — it splits on runs of
//	whitespace and drops empties — so the answer is simply how many fields it
//	returns. Clear and correct, but it allocates a slice of substrings just to
//	count them.
//
// Algorithm:
//
//  1. Call strings.Fields(s), which returns all whitespace-separated tokens
//     with empty strings already removed.
//  2. Return the length of that slice.
//
// Time:  O(n) — Fields scans the string once.
// Space: O(n) — it allocates the slice of token substrings.
func builtinSplit(s string) int {
	return len(strings.Fields(s)) // Fields drops empties, so every field is a segment
}

// ── Approach 2: Count Segment Starts (One Pass, Optimal) ──────────────────────
//
// countStarts counts a segment exactly once — at the position where it begins.
//
// Intuition:
//
//	Every segment has a unique first character: a non-space whose left neighbour
//	is either the string start or a space. Count those "boundary" positions and
//	you've counted the segments, with no allocation and a single pass. This is
//	the canonical O(1)-space answer.
//
// Algorithm:
//
//  1. Walk index i from 0 to n-1.
//  2. Increment the count when s[i] is NOT a space AND (i == 0 OR s[i-1] IS a
//     space) — that marks the first character of a new segment.
//  3. Return the count.
//
// Time:  O(n) — one pass over the characters.
// Space: O(1) — just a counter.
func countStarts(s string) int {
	count := 0
	for i := 0; i < len(s); i++ {
		// s[i] begins a segment iff it's a non-space preceded by a boundary
		// (start of string or a space).
		if s[i] != ' ' && (i == 0 || s[i-1] == ' ') {
			count++ // a new segment starts here
		}
	}
	return count
}

// ── Approach 3: Explicit State Machine (In-Segment Flag) ──────────────────────
//
// stateMachine tracks whether the scanner is currently inside a segment and
// counts each time it transitions from "outside" to "inside".
//
// Intuition:
//
//	Same idea as Approach 2, expressed as a two-state automaton: OUTSIDE (in
//	spaces) and INSIDE (in a token). A rising edge OUTSIDE→INSIDE means a new
//	segment just started, so count it. Useful when the "boundary" definition is
//	more complex than a single look-back (e.g. multiple whitespace kinds), since
//	the flag generalises cleanly.
//
// Algorithm:
//
//  1. Keep a boolean inSegment, initially false.
//  2. For each char: if it's non-space and inSegment is false, it's a new
//     segment → count++ and set inSegment = true. If it's a space, set
//     inSegment = false.
//  3. Return the count.
//
// Time:  O(n) — single pass.
// Space: O(1) — one flag and a counter.
func stateMachine(s string) int {
	count := 0
	inSegment := false // are we currently scanning inside a token?
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' {
			inSegment = false // spaces end any current segment
		} else if !inSegment {
			count++          // rising edge: just entered a new segment
			inSegment = true // remember we're inside it now
		}
	}
	return count
}

func main() {
	fmt.Println("=== Approach 1: Built-in Split + Filter (Brute Force) ===")
	fmt.Println(builtinSplit("Hello, my name is John")) // expected 5
	fmt.Println(builtinSplit("Hello"))                  // expected 1
	fmt.Println(builtinSplit(""))                       // expected 0
	fmt.Println(builtinSplit("   "))                    // expected 0 (all spaces)
	fmt.Println(builtinSplit("a b  c"))                 // expected 3 (double space)

	fmt.Println("=== Approach 2: Count Segment Starts (One Pass, Optimal) ===")
	fmt.Println(countStarts("Hello, my name is John")) // expected 5
	fmt.Println(countStarts("Hello"))                  // expected 1
	fmt.Println(countStarts(""))                       // expected 0
	fmt.Println(countStarts("   "))                    // expected 0
	fmt.Println(countStarts("a b  c"))                 // expected 3

	fmt.Println("=== Approach 3: Explicit State Machine (In-Segment Flag) ===")
	fmt.Println(stateMachine("Hello, my name is John")) // expected 5
	fmt.Println(stateMachine("Hello"))                  // expected 1
	fmt.Println(stateMachine(""))                       // expected 0
	fmt.Println(stateMachine("   "))                    // expected 0
	fmt.Println(stateMachine("a b  c"))                 // expected 3
}
