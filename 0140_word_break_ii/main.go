package main

import (
	"fmt"
	"sort"
	"strings"
)

// ── Approach 1: Backtracking (Brute Force) ───────────────────────────────────
//
// backtracking solves Word Break II by exhaustively trying every word split.
//
// Intuition:
//
//	Build sentences left to right: at each position try every dictionary word
//	that matches the upcoming characters; on a match, append the word to the
//	current path and recurse on the rest. Reaching the end of the string means
//	the path is one complete sentence. Undo (backtrack) and try other words.
//
// Algorithm:
//  1. dfs(start, path): if start == len(s), join path with spaces and record it.
//  2. For each end in (start, len(s)]: if s[start:end] is a word, push it onto
//     path, recurse dfs(end), then pop it (backtrack).
//
// Time:  O(2^n · n) — worst case ("aaaa…" with dict {a, aa, aaa, …}) has
//
//	exponentially many sentences; each is built in O(n).
//
// Space: O(n) — recursion depth and path buffer (output not counted).
func backtracking(s string, wordDict []string) []string {
	words := make(map[string]bool, len(wordDict)) // O(1) membership test
	for _, w := range wordDict {
		words[w] = true
	}
	var results []string
	var path []string // words chosen so far on the current root-to-here branch
	var dfs func(start int)
	dfs = func(start int) {
		if start == len(s) {
			results = append(results, strings.Join(path, " ")) // full sentence found
			return
		}
		for end := start + 1; end <= len(s); end++ {
			word := s[start:end]
			if !words[word] {
				continue // this slice is not a dictionary word
			}
			path = append(path, word) // choose
			dfs(end)                  // explore the suffix
			path = path[:len(path)-1] // un-choose (backtrack)
		}
	}
	dfs(0)
	return results
}

// ── Approach 2: DP Top-Down (Memoized DFS) ───────────────────────────────────
//
// dpTopDown solves Word Break II by caching all sentences for each suffix.
//
// Intuition:
//
//	Plain backtracking rebuilds the sentence set of the same suffix repeatedly.
//	Since the set of sentences for s[start:] is fixed, compute it once and
//	memoize: sentences(start) = { word + " " + rest | word matches at start,
//	rest ∈ sentences(start+len(word)) }.
//
// Algorithm:
//  1. sentences(start): memo hit → return cached slice.
//  2. Base: start == len(s) → return [""] sentinel meaning "empty tail".
//  3. For each matching word at start, prepend it to every sentence of the
//     corresponding suffix (with a space unless the tail is empty).
//  4. Memoize and return.
//
// Time:  O(2^n · n) worst case — the output itself can be exponential; with few
//
//	breakable points it is near O(n²) thanks to memoization.
//
// Space: O(2^n · n) — memo stores every sentence of every suffix.
func dpTopDown(s string, wordDict []string) []string {
	words := make(map[string]bool, len(wordDict))
	for _, w := range wordDict {
		words[w] = true
	}
	memo := make(map[int][]string) // start index → all sentences for s[start:]
	var sentences func(start int) []string
	sentences = func(start int) []string {
		if res, ok := memo[start]; ok {
			return res // suffix already solved
		}
		if start == len(s) {
			return []string{""} // one way to break the empty tail: no words at all
		}
		var res []string
		for end := start + 1; end <= len(s); end++ {
			word := s[start:end]
			if !words[word] {
				continue // not a word, no branch here
			}
			for _, tail := range sentences(end) {
				if tail == "" {
					res = append(res, word) // word ends the sentence exactly
				} else {
					res = append(res, word+" "+tail) // glue word before the tail
				}
			}
		}
		memo[start] = res // cache for any other path that reaches this suffix
		return res
	}
	return sentences(0)
}

// ── Approach 3: DP Bottom-Up (Prefix Sentence Table) ─────────────────────────
//
// dpBottomUp solves Word Break II by building sentence lists for every prefix.
//
// Intuition:
//
//	Mirror of #139's boolean table, but dp[i] stores ALL sentences covering
//	the prefix s[:i] instead of a mere yes/no. dp[i] extends every sentence in
//	dp[j] with the word s[j:i]. A pre-computed feasibility table from #139
//	prunes dead prefixes so we never build sentences that cannot reach the end.
//
// Algorithm:
//  1. Run the #139 boolean DP → canBreak[i] for every prefix.
//  2. If canBreak[n] is false, return [] immediately.
//  3. dp[0] = [""]; for i = 1..n with canBreak[i]: for every j < i with
//     dp[j] non-empty and s[j:i] a word, append (sentence + " " + word).
//  4. Return dp[n].
//
// Time:  O(2^n · n) worst case (output-bound); feasibility pruning keeps
//
//	no-solution inputs at O(n²·L).
//
// Space: O(2^n · n) — sentence lists for all prefixes.
func dpBottomUp(s string, wordDict []string) []string {
	words := make(map[string]bool, len(wordDict))
	for _, w := range wordDict {
		words[w] = true
	}
	n := len(s)

	// Phase 1: cheap boolean feasibility (exactly LeetCode #139).
	canBreak := make([]bool, n+1)
	canBreak[0] = true
	for i := 1; i <= n; i++ {
		for j := 0; j < i; j++ {
			if canBreak[j] && words[s[j:i]] {
				canBreak[i] = true
				break // one witness suffices
			}
		}
	}
	if !canBreak[n] {
		return []string{} // no segmentation exists; skip the expensive phase
	}

	// Phase 2: build actual sentences only along feasible prefixes.
	dp := make([][]string, n+1) // dp[i] = all sentences spelling s[:i]
	dp[0] = []string{""}        // empty prefix = one empty sentence
	for i := 1; i <= n; i++ {
		if !canBreak[i] {
			continue // dead prefix: building sentences here is wasted work
		}
		for j := 0; j < i; j++ {
			if len(dp[j]) == 0 || !words[s[j:i]] {
				continue // either s[:j] unbuildable or gap isn't a word
			}
			word := s[j:i]
			for _, sentence := range dp[j] {
				if sentence == "" {
					dp[i] = append(dp[i], word) // first word needs no leading space
				} else {
					dp[i] = append(dp[i], sentence+" "+word)
				}
			}
		}
	}
	return dp[n]
}

// sortedCopy returns a sorted copy so results with "any order" semantics can be
// printed deterministically for comparison against the expected output.
func sortedCopy(list []string) []string {
	out := make([]string, len(list))
	copy(out, list)
	sort.Strings(out)
	return out
}

func main() {
	type example struct {
		s    string
		dict []string
	}
	examples := []example{
		{"catsanddog", []string{"cat", "cats", "and", "sand", "dog"}},
		// expected (any order): ["cat sand dog","cats and dog"]
		{"pineapplepenapple", []string{"apple", "pen", "applepen", "pine", "pineapple"}},
		// expected (any order): ["pine apple pen apple","pine applepen apple","pineapple pen apple"]
		{"catsandog", []string{"cats", "dog", "sand", "and", "cat"}},
		// expected: []
	}

	fmt.Println("=== Approach 1: Backtracking (Brute Force) ===")
	for _, ex := range examples {
		fmt.Printf("s=%q  got=%v\n", ex.s, sortedCopy(backtracking(ex.s, ex.dict)))
		// expected [cat sand dog, cats and dog] / [pine apple pen apple, pine applepen apple, pineapple pen apple] / []
	}

	fmt.Println("=== Approach 2: DP Top-Down (Memoized DFS) ===")
	for _, ex := range examples {
		fmt.Printf("s=%q  got=%v\n", ex.s, sortedCopy(dpTopDown(ex.s, ex.dict)))
		// expected same sentence sets as above
	}

	fmt.Println("=== Approach 3: DP Bottom-Up (Prefix Sentence Table) ===")
	for _, ex := range examples {
		fmt.Printf("s=%q  got=%v\n", ex.s, sortedCopy(dpBottomUp(ex.s, ex.dict)))
		// expected same sentence sets as above
	}
}
