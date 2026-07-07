package main

import (
	"fmt"
	"strconv"
	"strings"
)

// collapse returns the shortest run-length encoding of a string `s` IF `s` is
// made of a repeated block, else it returns `s` unchanged. It does NOT recurse
// into the block — callers supply the already-best encoding of the block.
//
// The classic period test: `s` is k copies of some block P (k ≥ 2) iff the first
// occurrence of `s` inside `(s+s)[1 : 2len-1]` starts at index p < len(s); then
// P = s[:p] and the repeat count is len(s)/p.
func collapse(s string, encodedBlock func(block string) string) string {
	n := len(s)
	// Find the smallest period by searching s within (s+s) excluding index 0.
	doubled := (s + s)[1 : 2*n] // drop first char so index 0 can't match
	idx := strings.Index(doubled, s)
	p := idx + 1 // period length (offset back for the dropped char)
	if p < n {
		// s repeats every p chars → count = n/p copies of block s[:p].
		count := n / p
		block := encodedBlock(s[:p])
		return strconv.Itoa(count) + "[" + block + "]"
	}
	return s // not periodic: no run-length form
}

// ── Approach 1: Brute Force Interval Recursion (No Memo) ──────────────────────
//
// bruteForce solves Encode String with Shortest Length by recursively computing,
// for a substring, the best of: (a) leaving it raw, (b) splitting it into two
// halves and concatenating their best encodings, (c) if it is periodic, encoding
// it as count[bestEncoding(block)].
//
// Intuition:
//
//	The shortest encoding of a substring is either the raw substring, or a
//	concatenation of the shortest encodings of a left part and a right part, or —
//	when the substring is a repeated block — the k[...] form wrapping the block's
//	own shortest encoding. Take the minimum length over all these choices.
//
// Algorithm:
//  1. encode(s): best = s (raw).
//  2. For each split point k in 1..len-1: cand = encode(s[:k]) + encode(s[k:]);
//     keep the shorter.
//  3. If s is periodic with block P: cand = count[encode(P)]; keep the shorter.
//  4. Return best.
//
// Time:  O(n^4) or worse without memo — O(n^2) substrings × O(n) splits ×
//
//	O(n) substring/period work, all recomputed repeatedly (exponential fan-out).
//
// Space: O(n) recursion depth (plus substrings).
func bruteForce(s string) string {
	var encode func(sub string) string
	encode = func(sub string) string {
		if len(sub) == 0 {
			return ""
		}
		best := sub // option A: keep it raw
		// Option B: try every split into two non-empty halves.
		for k := 1; k < len(sub); k++ {
			left := encode(sub[:k])
			right := encode(sub[k:])
			if len(left)+len(right) < len(best) {
				best = left + right
			}
		}
		// Option C: if periodic, wrap the block's best encoding as count[...].
		enc := collapse(sub, encode)
		if len(enc) < len(best) {
			best = enc
		}
		return best
	}
	return encode(s)
}

// ── Approach 2: Interval DP (Top-Down Memoised) ──────────────────────────────
//
// memoDP solves Encode String with Shortest Length with the same interval
// recurrence, memoised on the (i, j) substring bounds so each substring is
// solved exactly once.
//
// Intuition:
//
//	There are only O(n^2) distinct substrings s[i..j]. Cache the best encoding of
//	each. The recurrence is identical: raw, best split, or periodic wrap — but
//	now the block's encoding inside collapse() is fetched from the memo, so the
//	whole thing is polynomial.
//
// Algorithm:
//  1. memo[i][j] = best encoding of s[i..j], filled lazily.
//  2. solve(i,j): if cached return it. Start best = s[i..j+1] raw.
//     For split k in i..j-1: combine solve(i,k)+solve(k+1,j).
//     If s[i..j] periodic: count[solve over the block].
//  3. Return solve(0, n-1).
//
// Time:  O(n^4) — O(n^2) states, each doing O(n) splits and O(n) period work.
// Space: O(n^2) memo (+ recursion).
func memoDP(s string) string {
	n := len(s)
	if n == 0 {
		return ""
	}
	memo := make([][]string, n)
	for i := range memo {
		memo[i] = make([]string, n)
	}

	var solve func(i, j int) string
	solve = func(i, j int) string {
		if memo[i][j] != "" {
			return memo[i][j] // already computed this substring
		}
		sub := s[i : j+1]
		best := sub // raw substring
		// Split into s[i..k] + s[k+1..j] for every internal boundary k.
		for k := i; k < j; k++ {
			cand := solve(i, k) + solve(k+1, j)
			if len(cand) < len(best) {
				best = cand
			}
		}
		// Periodic wrap: encode block via the memoised solver (over block bounds).
		enc := collapseIndexed(s, i, j, solve)
		if len(enc) < len(best) {
			best = enc
		}
		memo[i][j] = best
		return best
	}
	return solve(0, n-1)
}

// collapseIndexed is collapse() operating on s[i..j] but resolving the block's
// best encoding through the (i,j)-indexed memoised solver, so no re-slicing loses
// the cache. It returns the count[...] form if s[i..j] is periodic, else s[i..j].
func collapseIndexed(s string, i, j int, solve func(a, b int) string) string {
	sub := s[i : j+1]
	n := len(sub)
	doubled := (sub + sub)[1 : 2*n]
	idx := strings.Index(doubled, sub)
	p := idx + 1 // smallest period length
	if p < n {
		count := n / p
		// The block occupies s[i .. i+p-1]; encode it via the memo.
		block := solve(i, i+p-1)
		return strconv.Itoa(count) + "[" + block + "]"
	}
	return sub
}

// ── Approach 3: Interval DP (Bottom-Up, Optimal) ─────────────────────────────
//
// bottomUpDP solves Encode String with Shortest Length by filling dp[i][j] for
// increasing substring lengths, so every smaller substring is ready before the
// substrings that depend on it.
//
// Intuition:
//
//	Order the states by length. dp[i][j] uses only shorter substrings (splits and
//	the periodic block are all shorter than s[i..j]), so a length-ascending sweep
//	needs no recursion. Same three candidates: raw, best split, periodic wrap.
//
// Algorithm:
//  1. dp[i][i] = s[i] for all i (single chars).
//  2. For length L = 2..n, for each start i (j = i+L-1):
//     - dp[i][j] = s[i..j] raw.
//     - For split k in i..j-1: dp[i][j] = shorter(dp[i][j], dp[i][k]+dp[k+1][j]).
//     - If s[i..j] periodic with block length p: candidate =
//     count[dp[i][i+p-1]]; keep if shorter.
//  3. Return dp[0][n-1].
//
// Time:  O(n^4) — O(n^2) cells × O(n) splits (+ O(n) period detection per cell).
// Space: O(n^2) — the dp table.
func bottomUpDP(s string) string {
	n := len(s)
	if n == 0 {
		return ""
	}
	dp := make([][]string, n)
	for i := range dp {
		dp[i] = make([]string, n)
	}
	// Base case: length-1 substrings encode to themselves.
	for i := 0; i < n; i++ {
		dp[i][i] = string(s[i])
	}
	// Grow by substring length so dependencies (shorter substrings) are ready.
	for length := 2; length <= n; length++ {
		for i := 0; i+length-1 < n; i++ {
			j := i + length - 1
			sub := s[i : j+1]
			best := sub // raw
			// Best two-way split.
			for k := i; k < j; k++ {
				if len(dp[i][k])+len(dp[k+1][j]) < len(best) {
					best = dp[i][k] + dp[k+1][j]
				}
			}
			// Periodic wrap using the already-computed block cell.
			doubled := (sub + sub)[1 : 2*length]
			idx := strings.Index(doubled, sub)
			p := idx + 1 // smallest period
			if p < length {
				count := length / p
				cand := strconv.Itoa(count) + "[" + dp[i][i+p-1] + "]"
				if len(cand) < len(best) {
					best = cand
				}
			}
			dp[i][j] = best
		}
	}
	return dp[0][n-1]
}

func main() {
	// Note: LeetCode only encodes when it makes the string STRICTLY shorter.
	//   "aaa"          -> "aaa"      (3 raw; "3[a]" would be 4 — not shorter)
	//   "aaaaa"        -> "5[a]"     (4 < 5)
	//   "abcabc"       -> "abcabc"   (6 raw; "2[abc]" is also 6 — not shorter)
	//   "abcabcabcabc" -> "4[abc]"   (6 < 12)
	//   "abababab"     -> "4[ab]"    (5 < 8)
	//   "abbbabbbcabbbabbbc" -> "2[2[abbb]c]"  (nested; 11 < 18)

	fmt.Println("=== Approach 1: Brute Force Interval Recursion ===")
	fmt.Printf("%q => %q  (expected \"aaa\")\n", "aaa", bruteForce("aaa"))
	fmt.Printf("%q => %q  (expected \"5[a]\")\n", "aaaaa", bruteForce("aaaaa"))
	fmt.Printf("%q => %q  (expected \"4[abc]\")\n", "abcabcabcabc", bruteForce("abcabcabcabc"))
	fmt.Printf("%q => %q  (expected \"abcabc\")\n", "abcabc", bruteForce("abcabc"))

	fmt.Println("=== Approach 2: Interval DP (Top-Down Memoised) ===")
	fmt.Printf("%q => %q  (expected \"aaa\")\n", "aaa", memoDP("aaa"))
	fmt.Printf("%q => %q  (expected \"5[a]\")\n", "aaaaa", memoDP("aaaaa"))
	fmt.Printf("%q => %q  (expected \"4[abc]\")\n", "abcabcabcabc", memoDP("abcabcabcabc"))
	fmt.Printf("%q => %q  (expected \"4[ab]\")\n", "abababab", memoDP("abababab"))

	fmt.Println("=== Approach 3: Interval DP (Bottom-Up, Optimal) ===")
	fmt.Printf("%q => %q  (expected \"aaa\")\n", "aaa", bottomUpDP("aaa"))
	fmt.Printf("%q => %q  (expected \"5[a]\")\n", "aaaaa", bottomUpDP("aaaaa"))
	fmt.Printf("%q => %q  (expected \"4[abc]\")\n", "abcabcabcabc", bottomUpDP("abcabcabcabc"))
	fmt.Printf("%q => %q  (expected \"4[ab]\")\n", "abababab", bottomUpDP("abababab"))
	// Nested repeat: block "2[abbb]c" is itself repeated twice.
	fmt.Printf("%q => %q  (expected \"2[2[abbb]c]\")\n", "abbbabbbcabbbabbbc", bottomUpDP("abbbabbbcabbbabbbc"))
}
