package main

import "fmt"

// ── Approach 1: Brute-Force Simulation ───────────────────────────────────────
//
// bruteForce solves Count The Repetitions by walking the whole of str1 = [s1,n1]
// one character at a time and greedily matching s2 as a subsequence, counting how
// many complete copies of s2 (and therefore of str2 = [s2,n2]) we can carve out.
//
// Intuition:
//
//	"str2 = [s2,n2] can be obtained from str1" is exactly "s2 repeated (n2·m)
//	times is a subsequence of str1". Subsequence-matching is greedy: to match a
//	pattern p inside a text t, scan t left to right and advance a pointer into p
//	every time the current text character equals p[pointer]; each time the
//	pointer wraps past the end of p you have matched one full copy of p. Here
//	p = s2, t = s1 repeated n1 times. Count full s2 matches, then divide by n2.
//
// Algorithm:
//  1. j = index into s2 (0), cntS2 = number of completed s2 copies (0).
//  2. Repeat n1 times (one copy of s1 each):
//     for every character c of s1: if c == s2[j], advance j; when j reaches
//     len(s2), reset j to 0 and cntS2++.
//  3. The answer is cntS2 / n2 (integer division): how many [s2,n2] blocks fit.
//
// Time:  O(n1 · len(s1)) — every character of the fully expanded str1 is visited
//
//	once; with n1 up to 10^6 and |s1| up to 100 this is up to 10^8 steps, so it
//	is a correctness baseline that can TLE on the largest inputs.
//
// Space: O(1) — only the two counters j and cntS2.
func bruteForce(s1 string, n1 int, s2 string, n2 int) int {
	j := 0     // current position we are trying to match inside s2
	cntS2 := 0 // how many complete copies of s2 we have matched so far
	// Expand str1 = [s1, n1] copy by copy without materialising the big string.
	for i := 0; i < n1; i++ {
		// Consume one copy of s1, advancing the s2 pointer on every match.
		for k := 0; k < len(s1); k++ {
			if s1[k] == s2[j] { // this s1 char matches the char s2 needs next
				j++ // move on to the next character of s2
				if j == len(s2) {
					j = 0   // finished a whole s2 — wrap the pointer
					cntS2++ // and record the completed copy
				}
			}
		}
	}
	// cntS2 copies of s2 = cntS2/n2 copies of [s2, n2]; integer division floors it.
	return cntS2 / n2
}

// ── Approach 2: Cycle Detection / Pattern Fast-Forward (Optimal) ──────────────
//
// cycleDetection solves Count The Repetitions by noticing that the subsequence
// matcher is a finite-state machine: its only state between s1 copies is "which
// index of s2 are we in the middle of matching". With at most len(s2) states, the
// per-s1-copy behaviour MUST fall into a cycle within the first len(s2)+1 copies;
// once the cycle is found we jump over the bulk of the n1 copies arithmetically.
//
// Intuition:
//
//	Process s1 copies one by one, but after each copy record two running totals
//	keyed by the current s2 index j:
//	  countRec[i]  = total s2 copies completed after i copies of s1,
//	  indexRec[i]  = the value of j after i copies of s1.
//	The moment we see a j we have seen before (at an earlier copy `start`), the
//	block of copies (start .. i) is a repeating cycle that always produces the
//	same number of s2 copies. We can therefore fast-forward through as many whole
//	cycles as fit in the remaining n1 copies in O(1), then simulate the small
//	leftover tail. This turns O(n1·|s1|) into O(|s2|·|s1|).
//
// Algorithm:
//  1. Simulate s1 copies one at a time (like brute force), but after finishing
//     copy i store countRec[i] and indexRec[i], indexed by copy count.
//  2. Before storing, check whether the current j already appeared at some earlier
//     copy `start` (scan the recorded indices). If so we found a cycle.
//  3. cycleLen  = i - start (copies per cycle);
//     cycleCnt  = countRec[i] - countRec[start] (s2 copies per cycle).
//     cyclesLeft = (n1 - 1 - start) / cycleLen whole cycles fit after `start`.
//     total = countRec[start] (prefix) + cyclesLeft*cycleCnt (the jumped cycles).
//  4. Advance the copy counter past the jumped cycles and simulate the remaining
//     tail copies normally, adding their completed s2 copies.
//  5. Return total / n2.
//
// Time:  O(len(s2) · len(s1)) — at most len(s2)+1 distinct copies are simulated
//
//	before a state repeats; the rest are skipped by arithmetic. Independent of
//	n1's magnitude.
//
// Space: O(len(s2)) — the countRec / indexRec bookkeeping arrays, one slot per
//
//	simulated s1 copy (bounded by the number of distinct s2 indices + 1).
func cycleDetection(s1 string, n1 int, s2 string, n2 int) int {
	if n1 == 0 { // no copies of s1 at all → nothing can be matched
		return 0
	}
	// countRec[i] / indexRec[i]: state AFTER i copies of s1 have been consumed.
	// Index 0 means "before any copy": 0 completed, pointer at s2 index 0.
	countRec := make([]int, n1+1)
	indexRec := make([]int, n1+1)

	j := 0   // current index inside s2 we are matching
	cnt := 0 // total complete s2 copies matched so far
	countRec[0] = 0
	indexRec[0] = 0

	for i := 1; i <= n1; i++ { // i = number of s1 copies consumed so far
		// Consume the i-th copy of s1.
		for k := 0; k < len(s1); k++ {
			if s1[k] == s2[j] { // matched the char s2 currently needs
				j++
				if j == len(s2) {
					j = 0 // completed one full s2
					cnt++ // count it
				}
			}
		}
		countRec[i] = cnt // snapshot totals after this copy
		indexRec[i] = j

		// Look for an earlier copy that ended in the SAME s2 index j.
		for start := 0; start < i; start++ {
			if indexRec[start] == j { // state repeats → cycle from start..i
				cycleLen := i - start                     // s1 copies per cycle
				cycleCnt := countRec[i] - countRec[start] // s2 copies per cycle
				remaining := n1 - start                   // copies left after prefix
				cyclesLeft := remaining / cycleLen        // whole cycles that fit
				tail := remaining % cycleLen              // leftover copies to simulate

				// Prefix copies (0..start) already contribute countRec[start].
				total := countRec[start] + cyclesLeft*cycleCnt
				// The leftover `tail` copies produce the same delta the pattern
				// produced over its first `tail` copies past `start`.
				total += countRec[start+tail] - countRec[start]
				return total / n2 // how many [s2, n2] blocks fit
			}
		}
	}
	// No cycle detected within n1 copies (n1 small): plain division of the total.
	return countRec[n1] / n2
}

func main() {
	fmt.Println("=== Approach 1: Brute-Force Simulation ===")
	fmt.Printf("s1=acb n1=4 s2=ab n2=2   got=%d  expected 2\n", bruteForce("acb", 4, "ab", 2))
	fmt.Printf("s1=acb n1=1 s2=acb n2=1  got=%d  expected 1\n", bruteForce("acb", 1, "acb", 1))
	fmt.Printf("s1=aaa n1=3 s2=aa n2=1   got=%d  expected 4\n", bruteForce("aaa", 3, "aa", 1)) // 9 a's / 2 = 4
	fmt.Printf("s1=abc n1=0 s2=a n2=1    got=%d  expected 0\n", bruteForce("abc", 0, "a", 1))  // no s1 copies

	fmt.Println("=== Approach 2: Cycle Detection / Pattern Fast-Forward (Optimal) ===")
	fmt.Printf("s1=acb n1=4 s2=ab n2=2   got=%d  expected 2\n", cycleDetection("acb", 4, "ab", 2))
	fmt.Printf("s1=acb n1=1 s2=acb n2=1  got=%d  expected 1\n", cycleDetection("acb", 1, "acb", 1))
	fmt.Printf("s1=aaa n1=3 s2=aa n2=1   got=%d  expected 4\n", cycleDetection("aaa", 3, "aa", 1))
	fmt.Printf("s1=abc n1=0 s2=a n2=1    got=%d  expected 0\n", cycleDetection("abc", 0, "a", 1))
	fmt.Printf("s1=aaa n1=1000000 s2=a n2=1 got=%d  expected 3000000\n", cycleDetection("aaa", 1000000, "a", 1)) // large n1: cycle jump
}
