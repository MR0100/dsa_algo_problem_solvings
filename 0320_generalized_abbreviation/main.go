package main

import (
	"fmt"
	"sort"
	"strconv"
)

// ── Approach 1: Backtracking — Keep or Abbreviate Each Char ───────────────────
//
// backtracking solves Generalized Abbreviation by, at each index, choosing to
// either keep the character as-is or fold it into a run of abbreviated
// characters (a count).
//
// Intuition:
//
//	Walk the word left to right carrying `count`, the number of consecutive
//	characters currently being abbreviated. At index i we branch:
//	  (a) Abbreviate word[i]: count+1, move on (don't emit yet).
//	  (b) Keep word[i]: first flush any pending count as a number, then append
//	      the literal character; reset count to 0.
//	When we reach the end, flush any trailing count. This naturally enforces
//	"non-adjacent" abbreviations: a kept char resets the run, so two number
//	segments can never touch.
//
// Algorithm:
//  1. dfs(i, count, cur):
//     - If i == len(word): append cur + (count>0 ? count : ""); return.
//     - Branch A (abbreviate): dfs(i+1, count+1, cur).
//     - Branch B (keep): let s = cur + (count>0 ? count : "") + word[i];
//     dfs(i+1, 0, s).
//  2. Start dfs(0, 0, "").
//
// Time:  O(2^n · n) — 2^n leaves, each building a string of length ≤ n.
// Space: O(n) recursion depth (excluding the output list).
func backtracking(word string) []string {
	res := []string{}
	var dfs func(i, count int, cur string)
	dfs = func(i, count int, cur string) {
		if i == len(word) { // reached the end of the word
			if count > 0 {
				cur += strconv.Itoa(count) // flush trailing abbreviated run
			}
			res = append(res, cur)
			return
		}
		// Branch A: abbreviate word[i] — extend the current numeric run.
		dfs(i+1, count+1, cur)
		// Branch B: keep word[i] — first flush any pending count, then append.
		next := cur
		if count > 0 {
			next += strconv.Itoa(count) // close the run before a literal char
		}
		next += string(word[i])
		dfs(i+1, 0, next) // reset the run after keeping a character
	}
	dfs(0, 0, "")
	return res
}

// ── Approach 2: Bitmask Enumeration ──────────────────────────────────────────
//
// bitmask solves Generalized Abbreviation by enumerating all 2^n subsets: bit i
// set means "abbreviate word[i]", bit clear means "keep word[i]".
//
// Intuition:
//
//	Each character is independently kept or abbreviated, so the 2^n binary masks
//	over n positions enumerate every abbreviation exactly once. For a given mask,
//	sweep left to right accumulating consecutive abbreviated chars into a count,
//	flushing the count as a number whenever we hit a kept char.
//
// Algorithm:
//  1. For mask = 0 .. 2^n - 1:
//     - count = 0, build sb.
//     - For i = 0..n-1: if bit i set → count++; else flush count (if >0) then
//     append word[i], reset count.
//     - Flush any trailing count. Append sb to results.
//
// Time:  O(2^n · n).
// Space: O(n) per built string.
func bitmask(word string) []string {
	n := len(word)
	res := make([]string, 0, 1<<n)
	for mask := 0; mask < (1 << n); mask++ {
		count := 0
		cur := ""
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 { // bit set → abbreviate this char
				count++
			} else { // bit clear → keep this char
				if count > 0 {
					cur += strconv.Itoa(count) // flush run before the literal
					count = 0
				}
				cur += string(word[i])
			}
		}
		if count > 0 {
			cur += strconv.Itoa(count) // trailing abbreviated run
		}
		res = append(res, cur)
	}
	return res
}

// sortedCopy returns a lexicographically sorted copy — used only to compare the
// two approaches deterministically in main (LeetCode allows any order).
func sortedCopy(xs []string) []string {
	c := append([]string(nil), xs...)
	sort.Strings(c)
	return c
}

func main() {
	fmt.Println("=== Approach 1: Backtracking ===")
	got1 := backtracking("word")
	fmt.Printf("word count -> %d  expected 16\n", len(got1)) // 2^4 = 16
	fmt.Printf("sorted     -> %v\n", sortedCopy(got1))
	// expected (sorted): [1o1d 1o2 1or1 2r1 2rd 3d 4 w1r1 w1rd w2d w3 wo1d wo2 word wor1]  (15 shown; 16 total incl "1ord")

	got1b := backtracking("a")
	fmt.Printf("a          -> %v  expected [a 1] (any order)\n", got1b)

	fmt.Println("=== Approach 2: Bitmask Enumeration ===")
	got2 := bitmask("word")
	fmt.Printf("word count -> %d  expected 16\n", len(got2))
	fmt.Printf("matches backtracking (as sets) -> %v  expected true\n",
		fmt.Sprint(sortedCopy(got1)) == fmt.Sprint(sortedCopy(got2)))

	got2b := bitmask("a")
	fmt.Printf("a          -> %v  expected [a 1] (any order)\n", got2b)

	fmt.Println("=== Full sorted list for \"word\" (both approaches agree) ===")
	fmt.Printf("%v\n", sortedCopy(got2))
	// expected: [1o1d 1o2 1or1 1ord 2r1 2rd 3d 4 w1r1 w1rd w2d w3 wo1d wo2 wor1 word]
}
