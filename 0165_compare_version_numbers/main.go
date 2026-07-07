package main

import (
	"fmt"
	"strconv"
	"strings"
)

// ── Approach 1: Brute Force (Split + Parse + Compare) ────────────────────────
//
// splitAndCompare solves Compare Version Numbers by splitting both strings
// on '.', converting every revision to an integer, and comparing position by
// position with missing revisions treated as 0.
//
// Intuition:
//
//	The problem statement describes the algorithm directly: revisions are
//	the dot-separated chunks, their value is the integer conversion (which
//	automatically ignores leading zeros), and a shorter version is padded
//	with zero revisions. Materialise both revision lists and walk them
//	in parallel up to the longer length.
//
// Algorithm:
//  1. Split version1 and version2 on ".".
//  2. For i from 0 to max(len1, len2) − 1:
//     a. r1 = Atoi(parts1[i]) if i < len1 else 0.
//     b. r2 = Atoi(parts2[i]) if i < len2 else 0.
//     c. If r1 < r2 return −1; if r1 > r2 return 1.
//  3. All revisions equal → return 0.
//
// Time:  O(m + n) — splitting, parsing and comparing each character once.
// Space: O(m + n) — the two slices of revision substrings.
func splitAndCompare(version1, version2 string) int {
	parts1 := strings.Split(version1, ".") // revision chunks of version1
	parts2 := strings.Split(version2, ".") // revision chunks of version2
	longer := len(parts1)
	if len(parts2) > longer {
		longer = len(parts2) // compare up to the longer revision list
	}
	for i := 0; i < longer; i++ {
		r1, r2 := 0, 0 // missing revisions count as 0
		if i < len(parts1) {
			r1, _ = strconv.Atoi(parts1[i]) // Atoi drops leading zeros ("001" → 1)
		}
		if i < len(parts2) {
			r2, _ = strconv.Atoi(parts2[i])
		}
		if r1 < r2 { // first differing revision decides the order
			return -1
		}
		if r1 > r2 {
			return 1
		}
	}
	return 0 // every revision matched → the versions are equal
}

// ── Approach 2: Two Pointers In-Place Parse (Optimal) ────────────────────────
//
// twoPointers solves Compare Version Numbers with one pointer per string,
// parsing each revision on the fly without allocating substrings.
//
// Intuition:
//
//	We never need the whole revision lists at once — only the current pair
//	of revision values. Two cursors accumulate digits into an int until they
//	hit a '.' or the end of their string. A cursor that has run off the end
//	simply yields 0, which implements the "pad with zero revisions" rule for
//	free. This keeps the memory footprint constant.
//
// Algorithm:
//  1. i = 0 over version1, j = 0 over version2.
//  2. While i or j has characters left:
//     a. r1 = digits of version1 from i until '.'/end (0 if i is exhausted).
//     b. r2 = digits of version2 from j until '.'/end (0 if j is exhausted).
//     c. If r1 < r2 return −1; if r1 > r2 return 1.
//     d. Step both cursors past their '.' separator.
//  3. Return 0.
//
// Time:  O(m + n) — every character is visited exactly once.
// Space: O(1) — two cursors and two revision accumulators, no substrings.
func twoPointers(version1, version2 string) int {
	i, j := 0, 0 // cursors into version1 and version2
	for i < len(version1) || j < len(version2) {
		r1 := 0
		for i < len(version1) && version1[i] != '.' {
			r1 = r1*10 + int(version1[i]-'0') // accumulate digit (leading zeros vanish)
			i++
		}
		r2 := 0
		for j < len(version2) && version2[j] != '.' {
			r2 = r2*10 + int(version2[j]-'0')
			j++
		}
		if r1 < r2 { // first differing revision decides the order
			return -1
		}
		if r1 > r2 {
			return 1
		}
		i++ // skip the '.' (harmless past-the-end increment when exhausted)
		j++
	}
	return 0 // both strings fully consumed with equal revisions
}

func main() {
	type example struct {
		v1, v2 string
	}
	examples := []example{
		{"1.2", "1.10"},    // expected -1 (revision 2 < revision 10)
		{"1.01", "1.001"},  // expected 0  (both second revisions equal 1)
		{"1.0", "1.0.0.0"}, // expected 0  (missing revisions count as 0)
	}

	fmt.Println("=== Approach 1: Brute Force (Split + Parse + Compare) ===")
	for _, ex := range examples {
		fmt.Printf("v1=%q v2=%q  got=%d\n", ex.v1, ex.v2, splitAndCompare(ex.v1, ex.v2)) // expected -1, 0, 0
	}

	fmt.Println("=== Approach 2: Two Pointers In-Place Parse (Optimal) ===")
	for _, ex := range examples {
		fmt.Printf("v1=%q v2=%q  got=%d\n", ex.v1, ex.v2, twoPointers(ex.v1, ex.v2)) // expected -1, 0, 0
	}
}
