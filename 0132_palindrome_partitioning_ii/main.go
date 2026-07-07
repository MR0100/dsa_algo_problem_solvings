package main

import "fmt"

// ── Approach 1: Brute-Force Recursion ────────────────────────────────────────
//
// bruteForce solves Palindrome Partitioning II with plain recursion: try every
// palindromic prefix, recurse on the rest, take the minimum number of cuts.
//
// Intuition: if the whole string is a palindrome, zero cuts are needed.
// Otherwise the first cut splits off some palindromic prefix; the answer is
// 1 (that cut) plus the best answer for the remaining suffix. Trying every
// valid prefix and keeping the minimum explores all partitions.
//
// Algorithm:
//  1. If s is a palindrome, return 0.
//  2. For every split i in [1, len(s)-1]:
//     a. If s[:i] is a palindrome, candidate = 1 + bruteForce(s[i:]).
//     b. Keep the smallest candidate.
//  3. Return the minimum (worst case len(s)-1 cuts always works).
//
// Time:  O(n · 2^n) — explores every partition; each prefix check is O(n).
// Space: O(n) — recursion depth (substrings share the original backing array).
func bruteForce(s string) int {
	// isPalStr checks a whole string with two converging pointers.
	isPalStr := func(t string) bool {
		for i, j := 0, len(t)-1; i < j; i, j = i+1, j-1 {
			if t[i] != t[j] { // mismatch → not a palindrome
				return false
			}
		}
		return true
	}

	var solve func(t string) int
	solve = func(t string) int {
		if isPalStr(t) {
			return 0 // whole remainder is one palindrome: no cut needed
		}
		best := len(t) - 1 // cutting into single characters always works
		for i := 1; i < len(t); i++ {
			if isPalStr(t[:i]) { // prefix must be a palindrome to cut here
				if c := 1 + solve(t[i:]); c < best {
					best = c // found a partition with fewer cuts
				}
			}
		}
		return best
	}

	return solve(s)
}

// ── Approach 2: Top-Down DP (Memoization) ────────────────────────────────────
//
// dpTopDown solves Palindrome Partitioning II by memoizing the minimum cuts
// for every suffix and answering palindrome queries from a precomputed table.
//
// Intuition: the brute force recomputes the same suffixes exponentially many
// times, but minCuts(s[start:]) is independent of the prefix before start —
// a perfect overlapping subproblem. There are only n distinct suffixes, and
// palindromicity of every substring can be precomputed once in O(n²).
//
// Algorithm:
//  1. Build isPal[i][j] bottom-up (single chars true; s[i]==s[j] && inner).
//  2. solve(start): 0 if isPal[start][n-1]; else min over palindromic
//     prefixes s[start..end] of 1 + solve(end+1); memoize per start.
//  3. Answer is solve(0).
//
// Time:  O(n²) — n suffix states × n transitions, O(1) palindrome lookups.
// Space: O(n²) — the palindrome table dominates (memo is O(n)).
func dpTopDown(s string) int {
	n := len(s)

	// isPal[i][j] == true iff s[i..j] is a palindrome (bottom-up by length).
	isPal := make([][]bool, n)
	for i := range isPal {
		isPal[i] = make([]bool, n)
		isPal[i][i] = true // single characters are palindromes
	}
	for length := 2; length <= n; length++ {
		for i := 0; i+length-1 < n; i++ {
			j := i + length - 1
			if s[i] == s[j] {
				// ends match and the inside is a palindrome (or empty)
				isPal[i][j] = length == 2 || isPal[i+1][j-1]
			}
		}
	}

	memo := make([]int, n)
	for i := range memo {
		memo[i] = -1 // -1 marks "not computed yet"
	}

	var solve func(start int) int
	solve = func(start int) int {
		if isPal[start][n-1] {
			return 0 // whole suffix is a palindrome: zero cuts
		}
		if memo[start] != -1 {
			return memo[start] // reuse previously solved suffix
		}
		best := n - 1 - start // upper bound: cut every position
		for end := start; end < n-1; end++ {
			if isPal[start][end] { // first piece s[start..end] is valid
				if c := 1 + solve(end+1); c < best {
					best = c
				}
			}
		}
		memo[start] = best
		return best
	}

	return solve(0)
}

// ── Approach 3: Bottom-Up DP ─────────────────────────────────────────────────
//
// dpBottomUp solves Palindrome Partitioning II with an iterative 1-D cuts
// array over prefixes, backed by the O(n²) palindrome table.
//
// Intuition: let cuts[i] = minimum cuts for the prefix s[0..i]. If s[0..i]
// is itself a palindrome, cuts[i] = 0. Otherwise the last piece of an
// optimal partition is some palindrome s[j..i], and everything before it is
// an optimally cut prefix: cuts[i] = min over j of cuts[j-1] + 1.
//
// Algorithm:
//  1. Build isPal[i][j] exactly as in Approach 2.
//  2. For i from 0 to n-1:
//     a. If isPal[0][i], cuts[i] = 0.
//     b. Else cuts[i] = min{ cuts[j-1] + 1 : 1 <= j <= i, isPal[j][i] }.
//  3. Return cuts[n-1].
//
// Time:  O(n²) — table fill O(n²) + double loop over (i, j).
// Space: O(n²) — palindrome table (cuts array is O(n)).
func dpBottomUp(s string) int {
	n := len(s)

	// palindrome table, same construction as dpTopDown
	isPal := make([][]bool, n)
	for i := range isPal {
		isPal[i] = make([]bool, n)
		isPal[i][i] = true
	}
	for length := 2; length <= n; length++ {
		for i := 0; i+length-1 < n; i++ {
			j := i + length - 1
			if s[i] == s[j] {
				isPal[i][j] = length == 2 || isPal[i+1][j-1]
			}
		}
	}

	cuts := make([]int, n) // cuts[i] = min cuts for prefix s[0..i]
	for i := 0; i < n; i++ {
		if isPal[0][i] {
			cuts[i] = 0 // whole prefix is one palindrome
			continue
		}
		cuts[i] = i // worst case: i cuts → i+1 single characters
		for j := 1; j <= i; j++ {
			// last piece s[j..i] must be a palindrome; prefix s[0..j-1]
			// is already optimally solved in cuts[j-1]
			if isPal[j][i] && cuts[j-1]+1 < cuts[i] {
				cuts[i] = cuts[j-1] + 1
			}
		}
	}

	return cuts[n-1]
}

// ── Approach 4: Expand Around Center (Optimal) ───────────────────────────────
//
// expandAroundCenter solves Palindrome Partitioning II in O(n²) time but only
// O(n) extra space: instead of storing an n×n palindrome table, it generates
// every palindrome by expanding around its center and relaxes the cuts array
// on the fly.
//
// Intuition: every palindrome has a center (a character for odd length, a gap
// for even length). Expanding from all 2n-1 centers enumerates every
// palindromic substring exactly once. Whenever s[j..k] is discovered to be a
// palindrome, it can serve as the last piece of a partition of the prefix
// ending at k, giving the relaxation cuts[k+1] = min(cuts[k+1], cuts[j] + 1).
//
// Algorithm:
//  1. cuts[i] = min cuts for the prefix of length i; init cuts[i] = i-1
//     (cuts[0] = -1 so a full-prefix palindrome yields 0).
//  2. For each center c: expand odd (c,c) and even (c,c+1); for every
//     palindrome s[j..k] found, relax cuts[k+1] with cuts[j]+1.
//  3. Return cuts[n].
//
// Time:  O(n²) — 2n-1 centers, each expansion O(n).
// Space: O(n) — only the cuts array; no palindrome table.
func expandAroundCenter(s string) int {
	n := len(s)

	// cuts[i] = min cuts for prefix s[:i]; cuts[0] = -1 is a sentinel so
	// that a palindrome covering the whole prefix gives cuts = -1 + 1 = 0.
	cuts := make([]int, n+1)
	for i := 0; i <= n; i++ {
		cuts[i] = i - 1 // worst case: prefix of length i needs i-1 cuts
	}

	// expand grows a palindrome outward from (l, r) while ends match,
	// relaxing the cuts array for every palindrome it certifies.
	expand := func(l, r int) {
		for l >= 0 && r < n && s[l] == s[r] {
			// s[l..r] is a palindrome: use it as the final piece of the
			// prefix s[:r+1]; everything before it costs cuts[l], +1 cut.
			if cuts[l]+1 < cuts[r+1] {
				cuts[r+1] = cuts[l] + 1
			}
			l-- // widen the window symmetrically
			r++
		}
	}

	for c := 0; c < n; c++ {
		expand(c, c)   // odd-length palindromes centered at character c
		expand(c, c+1) // even-length palindromes centered at gap (c, c+1)
	}

	return cuts[n]
}

func main() {
	fmt.Println("=== Approach 1: Brute-Force Recursion ===")
	fmt.Println(bruteForce("aab")) // 1
	fmt.Println(bruteForce("a"))   // 0
	fmt.Println(bruteForce("ab"))  // 1

	fmt.Println("=== Approach 2: Top-Down DP (Memoization) ===")
	fmt.Println(dpTopDown("aab")) // 1
	fmt.Println(dpTopDown("a"))   // 0
	fmt.Println(dpTopDown("ab"))  // 1

	fmt.Println("=== Approach 3: Bottom-Up DP ===")
	fmt.Println(dpBottomUp("aab")) // 1
	fmt.Println(dpBottomUp("a"))   // 0
	fmt.Println(dpBottomUp("ab"))  // 1

	fmt.Println("=== Approach 4: Expand Around Center (Optimal) ===")
	fmt.Println(expandAroundCenter("aab")) // 1
	fmt.Println(expandAroundCenter("a"))   // 0
	fmt.Println(expandAroundCenter("ab"))  // 1
}
