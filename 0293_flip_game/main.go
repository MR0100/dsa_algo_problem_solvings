package main

import (
	"fmt"
	"strings"
)

// ── Approach 1: Byte Slice Scan (Brute-ish) ──────────────────────────────────
//
// byteSliceScan solves Flip Game by copying the string into a mutable byte
// slice for every candidate "++" position and rebuilding the flipped string.
//
// Intuition:
//
//	A move flips a consecutive "++" into "--". Every place where s[i]=='+' and
//	s[i+1]=='+' is one legal move. For each such spot, produce the resulting
//	string. Using a fresh byte-slice copy per move keeps the logic obvious at
//	the cost of allocating a new slice for each result.
//
// Algorithm:
//
//  1. For i = 0..len-2:
//  2. if s[i]=='+' and s[i+1]=='+':
//  3. copy s into a byte slice, set positions i,i+1 to '-', append string.
//  4. Return the collected results.
//
// Time:  O(n²) — up to n moves, each rebuilds an O(n) string.
// Space: O(n)  per result string (O(n²) total output, inherent to the problem).
func byteSliceScan(s string) []string {
	res := []string{}
	for i := 0; i+1 < len(s); i++ {
		if s[i] == '+' && s[i+1] == '+' { // a flippable "++" pair
			b := []byte(s) // fresh mutable copy of the whole string
			b[i] = '-'     // flip left plus
			b[i+1] = '-'   // flip right plus
			res = append(res, string(b))
		}
	}
	return res
}

// ── Approach 2: String Concatenation by Slicing (Optimal) ────────────────────
//
// sliceConcat solves Flip Game by assembling each result directly from the
// unchanged prefix, a literal "--", and the unchanged suffix.
//
// Intuition:
//
//	The only characters that change are the two at positions i, i+1. So the
//	result is exactly s[:i] + "--" + s[i+2:] — no per-character copy needed,
//	just three slices glued together. Same asymptotics but tighter and clearer.
//
// Algorithm:
//
//  1. For i = 0..len-2 where s[i]==s[i+1]=='+':
//  2. append s[:i] + "--" + s[i+2:].
//  3. Return results.
//
// Time:  O(n²) — n candidate positions, each builds an O(n) string.
// Space: O(n)  per output string.
func sliceConcat(s string) []string {
	res := []string{}
	for i := 0; i+1 < len(s); i++ {
		if s[i] == '+' && s[i+1] == '+' { // flippable pair at i, i+1
			// Everything outside [i, i+1] is untouched; splice "--" in the middle.
			res = append(res, s[:i]+"--"+s[i+2:])
		}
	}
	return res
}

// fmtList renders a []string as LeetCode's ["a","b"] output for easy comparison.
func fmtList(xs []string) string {
	quoted := make([]string, len(xs))
	for i, x := range xs {
		quoted[i] = `"` + x + `"`
	}
	return "[" + strings.Join(quoted, ",") + "]"
}

func main() {
	fmt.Println("=== Approach 1: Byte Slice Scan ===")
	fmt.Println(fmtList(byteSliceScan("++++"))) // expected ["--++","+--+","++--"]
	fmt.Println(fmtList(byteSliceScan("+")))    // expected []

	fmt.Println("=== Approach 2: Slice Concatenation (Optimal) ===")
	fmt.Println(fmtList(sliceConcat("++++"))) // expected ["--++","+--+","++--"]
	fmt.Println(fmtList(sliceConcat("+")))    // expected []
}
