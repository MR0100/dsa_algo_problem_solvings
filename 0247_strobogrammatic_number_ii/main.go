package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Recursive Build from Outside In (Optimal) ─────────────────────
//
// recursiveBuild generates all strobogrammatic numbers of length n.
//
// Intuition:
//
//	A strobogrammatic number is built from the OUTSIDE toward the center: the
//	outermost pair must be one of the mirror pairs (0,0),(1,1),(8,8),(6,9),
//	(9,6); the next pair inside is again any mirror pair, and so on. So if we
//	know all valid strings of length n-2, we can wrap each one in a mirror pair
//	to get length n. Recurse down: length 0 has the single empty string "",
//	length 1 has the self-symmetric digits "0","1","8".
//
//	The one twist: the OUTERMOST pair may not be (0,0), because that would put
//	a leading zero on a multi-digit number (not allowed). But inner pairs MAY
//	use (0,0). We track whether we are at the outermost layer to enforce this.
//
// Algorithm:
//  1. helper(n, total): returns all strobogrammatic strings of length n where
//     `total` is the ORIGINAL requested length (used to detect the outer layer).
//  2. Base cases: n==0 → [""]; n==1 → ["0","1","8"].
//  3. Recurse on inner = helper(n-2, total).
//  4. For each inner string, wrap it with every mirror pair (a,b):
//     skip (0,0) when n == total (outermost, would be a leading zero) unless
//     n == total is a single-digit-handled case (n>=2 here).
//  5. Collect and return.
//
// Time:  O(5^(n/2)) — roughly five choices per outer pair, ~n/2 pairs.
// Space: O(5^(n/2) · n) — all results held plus recursion depth n/2.
func recursiveBuild(n int) []string {
	return strobHelper(n, n) // seed the recursion with total == n
}

// mirror pairs: first char + its 180° partner form the outer shell.
var strobPairs = [][2]byte{
	{'0', '0'},
	{'1', '1'},
	{'6', '9'},
	{'8', '8'},
	{'9', '6'},
}

func strobHelper(n, total int) []string {
	if n == 0 {
		return []string{""} // an even build reaching the center: empty middle
	}
	if n == 1 {
		// Odd build reaching the center: only self-symmetric digits.
		return []string{"0", "1", "8"}
	}

	inner := strobHelper(n-2, total) // all valid cores two shorter
	result := []string{}

	for _, core := range inner { // wrap each shorter core...
		for _, p := range strobPairs { // ...in every legal mirror pair
			// At the outermost layer, a leading '0' is forbidden.
			if n == total && p[0] == '0' {
				continue
			}
			// New string = left partner + core + right partner.
			result = append(result, string(p[0])+core+string(p[1]))
		}
	}
	return result
}

// ── Approach 2: Iterative Layer-by-Layer Build ───────────────────────────────
//
// iterativeBuild produces the same set without recursion, growing the strings
// two characters at a time from the center outward.
//
// Intuition:
//
//	Start from the innermost layer — "" for even n, or {"0","1","8"} for odd n —
//	then repeatedly wrap every current string in each mirror pair, adding two
//	characters per round, until the strings reach length n. On the FINAL wrap
//	(the outermost pair) we forbid the (0,0) pair to avoid leading zeros.
//
// Algorithm:
//  1. Seed `current` = [""] if n even, else ["0","1","8"].
//  2. Set len_ to 0 or 1 accordingly.
//  3. While len_ < n: len_ += 2; wrap each string in every pair; on the last
//     round (len_ == n) skip (0,0).
//  4. Return current.
//
// Time:  O(5^(n/2)).
// Space: O(5^(n/2) · n).
func iterativeBuild(n int) []string {
	var current []string
	length := 0
	if n%2 == 0 {
		current = []string{""} // even length grows from empty core
	} else {
		current = []string{"0", "1", "8"} // odd length grows from a center digit
		length = 1
	}

	for length < n { // each round adds one outer pair (two chars)
		length += 2
		next := []string{}
		for _, core := range current {
			for _, p := range strobPairs {
				// Only the final (outermost) round forbids the leading zero.
				if length == n && p[0] == '0' {
					continue
				}
				next = append(next, string(p[0])+core+string(p[1]))
			}
		}
		current = next // this layer becomes the core for the next round
	}
	return current
}

// sortedCopy returns a sorted copy so the two approaches can be compared
// deterministically (they may enumerate in different orders).
func sortedCopy(xs []string) []string {
	out := append([]string(nil), xs...)
	sort.Strings(out)
	return out
}

func main() {
	fmt.Println("=== Approach 1: Recursive Build (Optimal) ===")
	fmt.Printf("n=2  got=%v  expected [11 69 88 96]\n", sortedCopy(recursiveBuild(2))) // expected [11 69 88 96]
	fmt.Printf("n=1  got=%v  expected [0 1 8]\n", sortedCopy(recursiveBuild(1)))       // expected [0 1 8]
	fmt.Printf("n=3  count=%d  expected 12\n", len(recursiveBuild(3)))                 // expected 12
	fmt.Printf("n=4  count=%d  expected 20\n", len(recursiveBuild(4)))                 // expected 20

	fmt.Println("=== Approach 2: Iterative Layer-by-Layer Build ===")
	fmt.Printf("n=2  got=%v  expected [11 69 88 96]\n", sortedCopy(iterativeBuild(2))) // expected [11 69 88 96]
	fmt.Printf("n=1  got=%v  expected [0 1 8]\n", sortedCopy(iterativeBuild(1)))       // expected [0 1 8]
	fmt.Printf("n=3  count=%d  expected 12\n", len(iterativeBuild(3)))                 // expected 12
	fmt.Printf("n=4  count=%d  expected 20\n", len(iterativeBuild(4)))                 // expected 20
}
