package main

import "fmt"

// ── Approach 1: Brute Force (Enumerate + Set) ────────────────────────────────
//
// bruteForce solves Unique Substrings in Wraparound String by generating every
// substring of s, keeping only those that are contiguous runs of the wraparound
// alphabet, and counting the distinct ones with a set.
//
// Intuition:
//
//	base = "...xyzabcabc...". A substring appears in base iff every adjacent pair
//	(prev, cur) inside it steps forward by one letter cyclically, i.e.
//	(cur - prev + 26) % 26 == 1 (with 'z' -> 'a' wrapping). Enumerate all O(n^2)
//	substrings; as we extend a substring one character at a time we can stop the
//	moment the wraparound step breaks. Put each valid substring in a set so
//	duplicates like the two "c"s in "cac" are counted once.
//
// Algorithm:
//  1. For each start index i, walk end index j forward from i.
//  2. Keep a growing substring; at every extension check the wraparound step
//     between s[j-1] and s[j]. If it breaks, this start can produce no longer
//     valid substring, so break.
//  3. Insert each still-valid substring s[i..j] into a set.
//  4. Answer = size of the set.
//
// Time:  O(n^2) substrings, each hashed at O(length) → up to O(n^3) in the worst
//
//	case (a fully increasing string like "abc...") — only feasible for small s.
//
// Space: O(number of distinct valid substrings) — can be O(n^2).
func bruteForce(s string) int {
	seen := make(map[string]struct{}) // distinct valid substrings
	n := len(s)
	for i := 0; i < n; i++ {
		// Grow the window [i..j]; stop as soon as the wraparound chain breaks.
		for j := i; j < n; j++ {
			if j > i {
				// step from previous char to current char, cyclically
				step := (int(s[j]) - int(s[j-1]) + 26) % 26
				if step != 1 { // not a +1 wraparound move → chain broken
					break
				}
			}
			seen[s[i:j+1]] = struct{}{} // record this contiguous wraparound run
		}
	}
	return len(seen) // each distinct substring counted exactly once
}

// ── Approach 2: DP on Longest Run Ending at Each Letter (Optimal) ─────────────
//
// dpMaxEndingAt solves Unique Substrings in Wraparound String by the key counting
// insight: the number of DISTINCT valid substrings that END at a given letter c
// equals the length of the LONGEST valid wraparound run ending at c.
//
// Intuition:
//
//	Every valid substring ending at letter c is a suffix of the longest valid run
//	ending at c: if the longest run ending at 'd' has length L, the distinct
//	substrings ending at that 'd' are exactly the L suffixes ("...d", of lengths
//	1..L), and any two different runs ending at 'd' overlap on their shorter
//	suffixes — so distinctness is captured by keeping, per ending letter, only the
//	MAXIMUM run length ever seen. Summing max[c] over all 26 letters counts every
//	distinct substring exactly once (bucketed by its last character), with no set
//	needed.
//
// Algorithm:
//  1. maxEnd[26]: longest valid wraparound run ending at each letter (init 0).
//  2. Scan s; keep curLen = length of the current run ending at s[i].
//     If s[i] continues the wraparound from s[i-1] ((s[i]-s[i-1]+26)%26==1),
//     curLen++; otherwise reset curLen = 1.
//  3. maxEnd[s[i]] = max(maxEnd[s[i]], curLen).
//  4. Answer = sum of maxEnd[0..25].
//
// Time:  O(n) — one linear pass over s plus a fixed 26-element sum.
// Space: O(1) — a fixed array of 26 counters.
func dpMaxEndingAt(s string) int {
	var maxEnd [26]int // maxEnd[c] = longest wraparound run ending at letter c
	curLen := 0        // length of the run ending at the current position
	for i := 0; i < len(s); i++ {
		if i > 0 && (int(s[i])-int(s[i-1])+26)%26 == 1 {
			curLen++ // s[i] extends the wraparound run from s[i-1]
		} else {
			curLen = 1 // run restarts at s[i] (or first char)
		}
		c := s[i] - 'a'         // 0..25 index of the ending letter
		if curLen > maxEnd[c] { // keep the longest run ever ending at c
			maxEnd[c] = curLen
		}
	}
	total := 0
	for _, v := range maxEnd { // sum the per-letter maxima
		total += v // v distinct substrings end at that letter
	}
	return total
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Enumerate + Set) ===")
	fmt.Printf("s=a    got=%d  expected 1\n", bruteForce("a"))
	fmt.Printf("s=cac  got=%d  expected 2\n", bruteForce("cac"))
	fmt.Printf("s=zab  got=%d  expected 6\n", bruteForce("zab"))
	fmt.Printf("s=abcd got=%d  expected 10\n", bruteForce("abcd")) // 4+3+2+1 distinct runs

	fmt.Println("=== Approach 2: DP on Longest Run Ending at Each Letter (Optimal) ===")
	fmt.Printf("s=a    got=%d  expected 1\n", dpMaxEndingAt("a"))
	fmt.Printf("s=cac  got=%d  expected 2\n", dpMaxEndingAt("cac"))
	fmt.Printf("s=zab  got=%d  expected 6\n", dpMaxEndingAt("zab"))
	fmt.Printf("s=abcd got=%d  expected 10\n", dpMaxEndingAt("abcd"))
	fmt.Printf("s=zaba got=%d  expected 6\n", dpMaxEndingAt("zaba")) // trailing 'a' adds nothing new
}
