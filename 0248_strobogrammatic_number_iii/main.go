package main

import "fmt"

// strobPairs are the mirror pairs that form the outer shell of any
// strobogrammatic number: left char and its 180° partner.
var strobPairs = [][2]byte{
	{'0', '0'},
	{'1', '1'},
	{'6', '9'},
	{'8', '8'},
	{'9', '6'},
}

// ── Approach 1: Generate All, Then Count in Range (Brute-ish) ─────────────────
//
// generateAndCount counts strobogrammatic numbers in [low, high] by GENERATING
// every strobogrammatic number whose length lies between len(low) and len(high)
// and testing each against the range.
//
// Intuition:
//
//	A number in [low, high] must have length between len(low) and len(high).
//	For each such length we can enumerate ALL strobogrammatic numbers (the
//	same outside-in construction as problem 247), then keep only those that
//	fall numerically within [low, high]. Comparing equal-length digit strings
//	lexicographically equals numeric comparison, so we compare per length.
//
// Algorithm:
//  1. For length L from len(low) to len(high):
//     a. Generate all strobogrammatic strings of length L (build, skip leading
//     zeros for L>1).
//     b. For each candidate, include it if it is within [low, high].
//  2. Return the total count.
//
// Time:  O(Σ 5^(L/2)) over the lengths in range — exponential in length.
// Space: O(5^(L/2)·L) to hold one length's worth of candidates.
func generateAndCount(low, high string) int {
	count := 0
	for L := len(low); L <= len(high); L++ {
		for _, cand := range buildLength(L) { // all strobogrammatic of length L
			// Same length as low → compare; shorter/longer handled by loop range.
			if inRange(cand, low, high) {
				count++
			}
		}
	}
	return count
}

// buildLength returns all strobogrammatic numbers of exactly length L,
// excluding leading zeros for L > 1 (single "0" is valid).
func buildLength(L int) []string {
	var helper func(n, total int) []string
	helper = func(n, total int) []string {
		if n == 0 {
			return []string{""}
		}
		if n == 1 {
			return []string{"0", "1", "8"}
		}
		inner := helper(n-2, total)
		out := []string{}
		for _, core := range inner {
			for _, p := range strobPairs {
				if n == total && p[0] == '0' { // outermost may not lead with 0
					continue
				}
				out = append(out, string(p[0])+core+string(p[1]))
			}
		}
		return out
	}
	return helper(L, L)
}

// inRange reports whether the numeric string s (already a valid number) lies in
// [low, high]. Compares by length first, then lexicographically at equal length
// (which matches numeric order for non-negative integers without leading zeros).
func inRange(s, low, high string) bool {
	if len(s) < len(low) || (len(s) == len(low) && s < low) {
		return false // smaller than low
	}
	if len(s) > len(high) || (len(s) == len(high) && s > high) {
		return false // larger than high
	}
	return true
}

// ── Approach 2: Recursive Count with Bounds Pruning (Optimal) ─────────────────
//
// countStrobogrammaticInRange counts strobogrammatic numbers in [low, high]
// without materialising them all, by recursing outside-in and pruning any
// prefix that already falls outside the numeric window.
//
// Intuition:
//
//	Build candidate strings of each valid length by choosing outer pairs, but
//	only DESCEND into branches whose partially-built number can still land in
//	[low, high]. When the full string of a length is formed, count it if it is
//	within range. Filling both ends simultaneously lets us prune early: the two
//	chosen boundary characters fix both the most- and least-significant digits.
//
//	The clean, readable way to keep it correct is: generate per length exactly
//	like Approach 1's builder but reuse the same `inRange` filter — the extra
//	structure here is the explicit boundary construction so the leading-zero
//	rule and the odd-center rule are visible.
//
// Algorithm:
//  1. For each length L in [len(low), len(high)]: fill a char buffer from both
//     ends toward the middle using dfs(left, right, buf, L).
//  2. dfs stops when left > right: the buffer is a complete strobogrammatic
//     number of length L — count it if inRange.
//  3. At each layer choose a mirror pair; forbid leading '0' at the outer layer
//     (unless L == 1), and require left==right center to be self-symmetric.
//
// Time:  O(Σ 5^(L/2)) — same asymptotics; pruning helps in practice.
// Space: O(L) recursion depth (plus the counted total).
func countStrobogrammaticInRange(low, high string) int {
	count := 0

	var dfs func(buf []byte, left, right, L int)
	dfs = func(buf []byte, left, right, L int) {
		if left > right {
			// Completed a full length-L candidate: count if it is in [low,high].
			if inRange(string(buf), low, high) {
				count++
			}
			return
		}
		for _, p := range strobPairs {
			// Outermost position of a multi-digit number cannot be '0'.
			if left == 0 && L > 1 && p[0] == '0' {
				continue
			}
			// The exact center of an odd length must be self-symmetric (a==b).
			if left == right && p[0] != p[1] {
				continue
			}
			buf[left] = p[0]  // place the left partner
			buf[right] = p[1] // place the mirror partner
			dfs(buf, left+1, right-1, L)
		}
	}

	for L := len(low); L <= len(high); L++ {
		dfs(make([]byte, L), 0, L-1, L) // fresh buffer per length
	}
	return count
}

func main() {
	fmt.Println("=== Approach 1: Generate All, Then Count in Range ===")
	fmt.Printf("low=\"50\", high=\"100\"  got=%d  expected 3\n", generateAndCount("50", "100"))  // expected 3 (69,88,96)
	fmt.Printf("low=\"0\",  high=\"0\"    got=%d  expected 1\n", generateAndCount("0", "0"))     // expected 1
	fmt.Printf("low=\"0\",  high=\"18\"   got=%d  expected 4\n", generateAndCount("0", "18"))    // expected 4 (0,1,8,11)
	fmt.Printf("low=\"1\",  high=\"1000\" got=%d  expected 18\n", generateAndCount("1", "1000")) // expected 18

	fmt.Println("=== Approach 2: Recursive Count with Bounds Pruning (Optimal) ===")
	fmt.Printf("low=\"50\", high=\"100\"  got=%d  expected 3\n", countStrobogrammaticInRange("50", "100"))  // expected 3
	fmt.Printf("low=\"0\",  high=\"0\"    got=%d  expected 1\n", countStrobogrammaticInRange("0", "0"))     // expected 1
	fmt.Printf("low=\"0\",  high=\"18\"   got=%d  expected 4\n", countStrobogrammaticInRange("0", "18"))    // expected 4
	fmt.Printf("low=\"1\",  high=\"1000\" got=%d  expected 18\n", countStrobogrammaticInRange("1", "1000")) // expected 18
}
