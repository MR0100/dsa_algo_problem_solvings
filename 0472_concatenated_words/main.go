package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Brute Force Word Break (Set + DFS, no memo) ───────────────────
//
// bruteForce solves Concatenated Words by, for each word, running a plain
// recursive word-break against the set of ALL other words, requiring at least
// one split so a word is never counted as a concatenation of just itself.
//
// Intuition:
//
//	A word is "concatenated" if it can be broken into ≥ 2 dictionary words. Put
//	every word in a hash set, then for the word under test try every prefix that
//	is in the set and recurse on the remaining suffix. To forbid the trivial
//	"the word itself" match, the very first cut must be a PROPER prefix (shorter
//	than the whole word); after that, any dictionary word — including the whole
//	original — is fine because the remaining suffix guarantees ≥ 2 pieces.
//
// Algorithm:
//  1. Insert all words into a set `dict`.
//  2. For each word w: canBreak(w, start=0, usedOne=false) via DFS.
//  3. canBreak(s, i, used): if i == len(s) return used (needed ≥ 1 real split).
//     For end in i+1..len(s): if s[i:end] in dict AND not (i==0 && end==len(s)),
//     and canBreak(s, end, true) → true.
//  4. Collect every w for which canBreak returns true.
//
// Time:  O(N · 2^L) worst case — N words, each word-break explores up to 2^L
//
//	prefix combinations (L = word length) without memoisation.
//
// Space: O(N·L) for the set + O(L) recursion.
func bruteForce(words []string) []string {
	dict := make(map[string]bool, len(words))
	for _, w := range words {
		dict[w] = true // every word is a candidate building block
	}

	var canBreak func(s string, i int) bool
	canBreak = func(s string, i int) bool {
		if i == len(s) {
			return true // consumed the whole word using only dictionary pieces
		}
		for end := i + 1; end <= len(s); end++ {
			// A cut of the FULL word (i==0 and end==len) would match the word
			// itself — that is a single word, not a concatenation. Forbid it.
			if i == 0 && end == len(s) {
				continue
			}
			if dict[s[i:end]] && canBreak(s, end) {
				return true // this prefix is a word and the rest also splits
			}
		}
		return false
	}

	res := []string{}
	for _, w := range words {
		if len(w) > 0 && canBreak(w, 0) {
			res = append(res, w) // w decomposes into ≥ 2 dictionary words
		}
	}
	return res
}

// ── Approach 2: Word Break DP (Memoised DFS) ─────────────────────────────────
//
// memoDFS solves Concatenated Words with the same word-break test but caches, per
// starting index, whether the suffix from that index can be fully segmented,
// turning the exponential search into a polynomial DP per word.
//
// Intuition:
//
//	The suffix s[i:] is segmentable or not, independent of how we reached i. So
//	memoise canBreak(i) per word. To respect "≥ 2 words", we still forbid the
//	single cut that equals the whole word. A clean way: define the memo over
//	"can the suffix starting at i be broken using dictionary words", and only
//	require that the SPLIT used at least one interior boundary — enforced by
//	starting with the whole-word cut disallowed at i==0.
//
// Algorithm:
//  1. Build set `dict`.
//  2. For each word w: memo = map[int]int8 (0 unknown, 1 true, -1 false).
//     seg(i): if i==len(w) true. For end in i+1..len: skip the whole-word cut
//     when i==0&&end==len; if w[i:end] in dict and seg(end) → memo true.
//  3. w is concatenated iff seg(0).
//
// Time:  O(N · L^2) amortised — per word, L start indices × O(L) substring cuts,
//
//	each suffix computed once thanks to the memo (substring hashing folded in).
//
// Space: O(N·L) set + O(L) memo/recursion per word.
func memoDFS(words []string) []string {
	dict := make(map[string]bool, len(words))
	for _, w := range words {
		dict[w] = true
	}

	res := []string{}
	for _, w := range words {
		if len(w) == 0 {
			continue
		}
		memo := make(map[int]int8) // per-word cache: suffix start → segmentable?

		var seg func(i int) bool
		seg = func(i int) bool {
			if i == len(w) {
				return true // fully segmented
			}
			if v, ok := memo[i]; ok {
				return v == 1 // reuse cached verdict for this suffix
			}
			ok := false
			for end := i + 1; end <= len(w); end++ {
				// Disallow matching the entire word as one piece.
				if i == 0 && end == len(w) {
					continue
				}
				if dict[w[i:end]] && seg(end) {
					ok = true
					break
				}
			}
			if ok {
				memo[i] = 1
			} else {
				memo[i] = -1
			}
			return ok
		}

		if seg(0) {
			res = append(res, w)
		}
	}
	return res
}

// ── Approach 3: Trie + DFS (Optimal) ─────────────────────────────────────────
//
// trieDFS solves Concatenated Words by storing all words in a trie and, for each
// word, walking the trie to enumerate dictionary-word prefixes without slicing
// strings, then recursing on the rest — segmentation guided by shared prefixes.
//
// Intuition:
//
//	Building block words share prefixes; a trie lets us discover every dictionary
//	prefix of the current word in a single character walk (following child links)
//	instead of hashing every substring. At each trie node marking end-of-word we
//	have a valid cut; recurse on the remaining suffix. Require ≥ 2 pieces by not
//	accepting the cut that consumes the whole word on the first segment.
//
// Algorithm:
//  1. Insert all words into a trie (isEnd flags).
//  2. For each word w: dfs(start): from the root, walk characters w[start..];
//     whenever a node isEnd at position p (and it is not the whole word on the
//     first segment), if p==len(w) succeed, else dfs(p). Cache visited starts.
//  3. Collect w when dfs(0) succeeds.
//
// Time:  O(N · L^2) — per word, L start positions × O(L) trie walk; memoised
//
//	starts avoid recomputation.
//
// Space: O(total characters) for the trie + O(L) recursion.
func trieDFS(words []string) []string {
	root := &trieNode{children: [26]*trieNode{}}
	for _, w := range words {
		root.insert(w) // build the prefix tree of building blocks
	}

	res := []string{}
	for _, w := range words {
		if len(w) == 0 {
			continue
		}
		memo := make([]int8, len(w)+1) // 0 unknown, 1 true, -1 false, indexed by start

		var dfs func(start int, isFirst bool) bool
		dfs = func(start int, isFirst bool) bool {
			if start == len(w) {
				return true // reached the end via valid cuts
			}
			if memo[start] != 0 {
				return memo[start] == 1
			}
			node := root
			ok := false
			// Walk characters from `start`, following trie child links.
			for p := start; p < len(w); p++ {
				c := w[p] - 'a'
				if node.children[c] == nil {
					break // no dictionary word shares this prefix — stop
				}
				node = node.children[c]
				if node.isEnd {
					// A dictionary word ends at index p (inclusive) → cut after p.
					end := p + 1
					// Forbid the whole word as a single first piece.
					if isFirst && end == len(w) {
						continue
					}
					if dfs(end, false) {
						ok = true
						break
					}
				}
			}
			if ok {
				memo[start] = 1
			} else {
				memo[start] = -1
			}
			return ok
		}

		if dfs(0, true) {
			res = append(res, w)
		}
	}
	return res
}

// trieNode is a 26-ary lowercase trie node.
type trieNode struct {
	children [26]*trieNode
	isEnd    bool // true if some inserted word ends exactly here
}

// insert adds a word into the trie, creating nodes as needed.
func (t *trieNode) insert(word string) {
	node := t
	for i := 0; i < len(word); i++ {
		c := word[i] - 'a'
		if node.children[c] == nil {
			node.children[c] = &trieNode{}
		}
		node = node.children[c]
	}
	node.isEnd = true // mark the terminal node
}

// sortedCopy returns a sorted copy so output order is deterministic for tests.
func sortedCopy(xs []string) []string {
	cp := append([]string(nil), xs...)
	sort.Strings(cp)
	return cp
}

func main() {
	ex1 := []string{"cat", "cats", "catsdogcats", "dog", "dogcatsdog", "hippopotamuses", "rat", "ratcatdogcat"}
	ex2 := []string{"cat", "dog", "catdog"}

	fmt.Println("=== Approach 1: Brute Force Word Break ===")
	fmt.Printf("ex1 => %v\n", sortedCopy(bruteForce(ex1))) // expected [catsdogcats dogcatsdog ratcatdogcat]
	fmt.Printf("ex2 => %v\n", sortedCopy(bruteForce(ex2))) // expected [catdog]

	fmt.Println("=== Approach 2: Word Break DP (Memoised) ===")
	fmt.Printf("ex1 => %v\n", sortedCopy(memoDFS(ex1))) // expected [catsdogcats dogcatsdog ratcatdogcat]
	fmt.Printf("ex2 => %v\n", sortedCopy(memoDFS(ex2))) // expected [catdog]

	fmt.Println("=== Approach 3: Trie + DFS (Optimal) ===")
	fmt.Printf("ex1 => %v\n", sortedCopy(trieDFS(ex1))) // expected [catsdogcats dogcatsdog ratcatdogcat]
	fmt.Printf("ex2 => %v\n", sortedCopy(trieDFS(ex2))) // expected [catdog]
}
