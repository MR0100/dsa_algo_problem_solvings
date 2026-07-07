package main

import "fmt"

// ── Approach 1: Brute Force (Recursion) ──────────────────────────────────────
//
// bruteForce solves Word Break with plain recursion over split points.
//
// Intuition:
//
//	Try every dictionary word as a prefix of the remaining string. If some word
//	matches the front and the rest of the string can also be broken, the whole
//	string can. Without caching, the same suffix is re-solved exponentially
//	many times (e.g. "aaaaab" with dict ["a","aa","aaa"]).
//
// Algorithm:
//  1. canBreak(start): if start == len(s), the whole string is consumed → true.
//  2. For every end in (start, len(s)]: if s[start:end] is a dictionary word
//     and canBreak(end) is true → true.
//  3. Otherwise false.
//
// Time:  O(2^n · n) — 2^(n-1) split patterns, each with O(n) substring/hash work.
// Space: O(n) — recursion depth (plus the word set).
func bruteForce(s string, wordDict []string) bool {
	words := make(map[string]bool, len(wordDict)) // O(1) word lookup
	for _, w := range wordDict {
		words[w] = true
	}
	var canBreak func(start int) bool
	canBreak = func(start int) bool {
		if start == len(s) {
			return true // consumed the whole string successfully
		}
		for end := start + 1; end <= len(s); end++ {
			// try s[start:end] as the next word, then recurse on the suffix
			if words[s[start:end]] && canBreak(end) {
				return true
			}
		}
		return false // no word fits the front of s[start:]
	}
	return canBreak(0)
}

// ── Approach 2: DP Top-Down (Memoization) ────────────────────────────────────
//
// dpTopDown solves Word Break with recursion + memoized suffix results.
//
// Intuition:
//
//	The brute force recomputes canBreak(start) for the same start many times,
//	but its answer never changes — there are only n+1 distinct subproblems.
//	Cache them and the exponential tree collapses to O(n²) edges.
//
// Algorithm:
//  1. memo[start] stores the known answer for suffix s[start:].
//  2. Same recursion as brute force, but consult/populate memo.
//
// Time:  O(n² · L) — n subproblems × n split points, each doing an O(L)
//
//	substring hash (L = max word length; O(n³) worst case bound).
//
// Space: O(n) — memo array + recursion stack.
func dpTopDown(s string, wordDict []string) bool {
	words := make(map[string]bool, len(wordDict))
	for _, w := range wordDict {
		words[w] = true
	}
	memo := make(map[int]bool, len(s)) // start index → can s[start:] be broken?
	var canBreak func(start int) bool
	canBreak = func(start int) bool {
		if start == len(s) {
			return true // empty suffix always breaks
		}
		if res, ok := memo[start]; ok {
			return res // already solved this suffix
		}
		for end := start + 1; end <= len(s); end++ {
			if words[s[start:end]] && canBreak(end) {
				memo[start] = true // cache success
				return true
			}
		}
		memo[start] = false // cache failure so we never redo this suffix
		return false
	}
	return canBreak(0)
}

// ── Approach 3: DP Bottom-Up (Optimal) ───────────────────────────────────────
//
// dpBottomUp solves Word Break with a 1-D boolean table over prefixes.
//
// Intuition:
//
//	dp[i] = "can the prefix s[:i] be segmented?". A prefix of length i works if
//	some earlier breakable prefix s[:j] is extended by a dictionary word
//	s[j:i]. Build from the empty prefix up — no recursion, no stack.
//
// Algorithm:
//  1. dp[0] = true (empty prefix).
//  2. For i = 1..n: dp[i] = OR over j < i of (dp[j] AND s[j:i] ∈ words).
//  3. Answer is dp[n].
//
// Time:  O(n² · L) — n×n (i,j) pairs with O(L) substring hashing.
// Space: O(n) — the dp table.
func dpBottomUp(s string, wordDict []string) bool {
	words := make(map[string]bool, len(wordDict))
	for _, w := range wordDict {
		words[w] = true
	}
	n := len(s)
	dp := make([]bool, n+1) // dp[i] ⇔ s[:i] is segmentable
	dp[0] = true            // empty prefix is trivially segmentable
	for i := 1; i <= n; i++ {
		for j := 0; j < i; j++ {
			// s[:i] works if s[:j] works and the gap s[j:i] is a word
			if dp[j] && words[s[j:i]] {
				dp[i] = true
				break // one witness is enough
			}
		}
	}
	return dp[n]
}

// ── Approach 4: BFS over Indices ─────────────────────────────────────────────
//
// bfsApproach solves Word Break as reachability in an implicit graph.
//
// Intuition:
//
//	Treat each index 0..n as a graph node; an edge start→end exists when
//	s[start:end] is a dictionary word. The string breaks iff node n is
//	reachable from node 0. BFS with a visited set explores each index once.
//
// Algorithm:
//  1. Queue = {0}, visited = {0}.
//  2. Pop start; for each end with s[start:end] ∈ words: if end == n → true;
//     else push end if unvisited.
//  3. Queue empty → false.
//
// Time:  O(n² · L) — each index expands at most once over ≤ n edges.
// Space: O(n) — queue + visited set.
func bfsApproach(s string, wordDict []string) bool {
	words := make(map[string]bool, len(wordDict))
	for _, w := range wordDict {
		words[w] = true
	}
	n := len(s)
	if n == 0 {
		return true // empty string trivially breaks
	}
	visited := make([]bool, n+1) // indices already expanded (avoid rework)
	queue := []int{0}            // frontier of reachable split points
	visited[0] = true
	for len(queue) > 0 {
		start := queue[0]
		queue = queue[1:] // pop front
		for end := start + 1; end <= n; end++ {
			if !words[s[start:end]] {
				continue // no edge start→end
			}
			if end == n {
				return true // reached the end of the string
			}
			if !visited[end] {
				visited[end] = true // mark before enqueue to prevent duplicates
				queue = append(queue, end)
			}
		}
	}
	return false // end index unreachable
}

// ── Approach 5: Trie + DP ────────────────────────────────────────────────────
//
// trieDP solves Word Break with a trie replacing per-substring hashing.
//
// Intuition:
//
//	From each breakable index j, instead of hashing every substring s[j:i],
//	walk the trie character by character; every trie node marked as a word end
//	reveals a valid i in one pass. Matching all words starting at j costs
//	O(maxWordLen) instead of O(n·L) hashing.
//
// Algorithm:
//  1. Insert all dictionary words into a trie.
//  2. dp[0] = true. For each j with dp[j] true, walk the trie along
//     s[j], s[j+1], ...; whenever a word ends at position i, set dp[i] = true.
//  3. Answer is dp[n].
//
// Time:  O(n · M + totalWordChars) — M = longest word length; each start walks
//
//	at most M characters in the trie.
//
// Space: O(totalWordChars + n) — trie nodes plus the dp table.
func trieDP(s string, wordDict []string) bool {
	type trieNode struct {
		children map[byte]*trieNode
		isWord   bool
	}
	newNode := func() *trieNode { return &trieNode{children: map[byte]*trieNode{}} }

	// Build the trie of dictionary words.
	root := newNode()
	for _, w := range wordDict {
		curr := root
		for i := 0; i < len(w); i++ {
			c := w[i]
			if curr.children[c] == nil {
				curr.children[c] = newNode() // create the path lazily
			}
			curr = curr.children[c]
		}
		curr.isWord = true // mark the end of a full word
	}

	n := len(s)
	dp := make([]bool, n+1) // dp[i] ⇔ s[:i] is segmentable
	dp[0] = true
	for j := 0; j <= n; j++ {
		if !dp[j] {
			continue // can't start a word from an unreachable index
		}
		curr := root
		for i := j; i < n; i++ {
			curr = curr.children[s[i]] // extend the match by one character
			if curr == nil {
				break // no dictionary word continues this way
			}
			if curr.isWord {
				dp[i+1] = true // s[j:i+1] is a word extending a breakable prefix
			}
		}
	}
	return dp[n]
}

func main() {
	type example struct {
		s    string
		dict []string
	}
	examples := []example{
		{"leetcode", []string{"leet", "code"}},                       // expected true
		{"applepenapple", []string{"apple", "pen"}},                  // expected true
		{"catsandog", []string{"cats", "dog", "sand", "and", "cat"}}, // expected false
	}

	fmt.Println("=== Approach 1: Brute Force (Recursion) ===")
	for _, ex := range examples {
		fmt.Printf("s=%q dict=%v  got=%v\n", ex.s, ex.dict, bruteForce(ex.s, ex.dict)) // expected true, true, false
	}

	fmt.Println("=== Approach 2: DP Top-Down (Memoization) ===")
	for _, ex := range examples {
		fmt.Printf("s=%q dict=%v  got=%v\n", ex.s, ex.dict, dpTopDown(ex.s, ex.dict)) // expected true, true, false
	}

	fmt.Println("=== Approach 3: DP Bottom-Up (Optimal) ===")
	for _, ex := range examples {
		fmt.Printf("s=%q dict=%v  got=%v\n", ex.s, ex.dict, dpBottomUp(ex.s, ex.dict)) // expected true, true, false
	}

	fmt.Println("=== Approach 4: BFS over Indices ===")
	for _, ex := range examples {
		fmt.Printf("s=%q dict=%v  got=%v\n", ex.s, ex.dict, bfsApproach(ex.s, ex.dict)) // expected true, true, false
	}

	fmt.Println("=== Approach 5: Trie + DP ===")
	for _, ex := range examples {
		fmt.Printf("s=%q dict=%v  got=%v\n", ex.s, ex.dict, trieDP(ex.s, ex.dict)) // expected true, true, false
	}
}
