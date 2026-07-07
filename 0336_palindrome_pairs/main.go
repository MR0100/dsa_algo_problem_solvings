package main

import (
	"fmt"
	"sort"
)

// isPalindrome reports whether s[lo..hi] (inclusive) reads the same both ways.
func isPalindrome(s string, lo, hi int) bool {
	for lo < hi {
		if s[lo] != s[hi] {
			return false // mismatched pair → not a palindrome
		}
		lo++
		hi--
	}
	return true
}

// isPal is the whole-string convenience wrapper.
func isPal(s string) bool { return isPalindrome(s, 0, len(s)-1) }

// reverse returns s reversed.
func reverse(s string) string {
	b := []byte(s)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i] // swap symmetric bytes
	}
	return string(b)
}

// ── Approach 1: Brute Force (all ordered pairs) ──────────────────────────────
//
// bruteForce solves Palindrome Pairs by testing every ordered pair (i, j),
// i != j, and checking whether words[i]+words[j] is a palindrome.
//
// Intuition:
//
//	The definition asks for all ordered index pairs whose concatenation is a
//	palindrome. Just try them all: n·(n−1) concatenations, each palindrome-
//	checked in linear time. Simple and obviously correct — the baseline.
//
// Algorithm:
//  1. For every ordered pair (i, j) with i != j:
//  2. if words[i]+words[j] is a palindrome, record [i, j].
//
// Time:  O(n^2 · m) — n^2 pairs, each concatenation checked in O(m), m avg len.
// Space: O(1) extra beyond the output.
func bruteForce(words []string) [][]int {
	res := [][]int{}
	for i := 0; i < len(words); i++ {
		for j := 0; j < len(words); j++ {
			if i == j {
				continue // a word cannot pair with itself
			}
			if isPal(words[i] + words[j]) { // concatenate and test
				res = append(res, []int{i, j})
			}
		}
	}
	return res
}

// ── Approach 2: Hash Map of Reversed Words (Optimal) ─────────────────────────
//
// hashMap solves Palindrome Pairs by, for each word, splitting it at every cut
// and using a map from reversed-word → index to find the complement.
//
// Intuition:
//
//	Split word w into (left, right) at each position. Two cases make w+other a
//	palindrome:
//	  (a) left is a palindrome and reverse(right) exists as another word →
//	      that word goes in FRONT:  other + w.
//	  (b) right is a palindrome and reverse(left) exists as another word →
//	      that word goes BEHIND:    w + other.
//	Look up reverse(substring) in a precomputed map to find the complement in
//	O(1). Careful de-duplication avoids counting the empty-cut cases twice.
//
// Algorithm:
//  1. Build lookup: reverse(word) → its index.
//  2. For each word i, for each cut position c in [0, len]:
//     left = w[:c], right = w[c:].
//     - If left is a palindrome and lookup[right] exists (≠ i): [lookup[right], i].
//     - If c != len (avoid empty right dup) and right is a palindrome and
//     lookup[left] exists (≠ i): [i, lookup[left]].
//
// Time:  O(n · m^2) — n words, m cut positions, each with an O(m) palindrome
//
//	check + O(m) reversed-substring hash.
//
// Space: O(n · m) — the reversed-word map.
func hashMap(words []string) [][]int {
	// lookup: reversed word → index, so we can find complements by exact match.
	lookup := make(map[string]int, len(words))
	for i, w := range words {
		lookup[reverse(w)] = i
	}
	res := [][]int{}
	for i, w := range words {
		for c := 0; c <= len(w); c++ {
			left, right := w[:c], w[c:] // split w at cut c
			// Case (a): left palindromic → complement (= reverse(right)) sits in front.
			if isPal(left) {
				if j, ok := lookup[right]; ok && j != i {
					res = append(res, []int{j, i}) // other + w
				}
			}
			// Case (b): right palindromic → complement (= reverse(left)) sits behind.
			// Guard c != len(w) so the empty-right split isn't double-counted with (a).
			if c != len(w) && isPal(right) {
				if j, ok := lookup[left]; ok && j != i {
					res = append(res, []int{i, j}) // w + other
				}
			}
		}
	}
	return res
}

// TrieNode is a node of the reversed-word trie used by the trie approach.
type TrieNode struct {
	children [26]*TrieNode // edges for 'a'..'z'
	wordIdx  int           // index of a word ENDING here, else -1
	palBelow []int         // indices of words whose remaining suffix (below this
	// node) is itself a palindrome — enables the "current word is longer" case
}

// newTrieNode allocates an empty node with no terminating word.
func newTrieNode() *TrieNode {
	return &TrieNode{wordIdx: -1}
}

// ── Approach 3: Trie of Reversed Words (Advanced Optimal) ────────────────────
//
// trieApproach solves Palindrome Pairs by inserting each word REVERSED into a
// trie (annotated with palindromic-suffix info) and then, for each word, walking
// it forward through the trie to discover all complements.
//
// Intuition:
//
//	For words[i] + words[j] to be a palindrome, words[j] reversed must "match"
//	the front of words[i]. Store all words reversed in a trie. Walk words[i]
//	character by character:
//	  • If we land on a node that terminates some word j (j != i) AND the REST
//	    of words[i] (unmatched suffix) is a palindrome → words[i]+words[j] is a
//	    palindrome (words[j] is shorter or equal).
//	  • If words[i] runs out but the trie has deeper words whose remaining part
//	    (palBelow) is palindromic → those j (longer words) also pair.
//	This is the classic trie solution; palBelow is precomputed at insert time.
//
// Algorithm:
//  1. Insert each word reversed. At every node along the path, if the REMAINING
//     suffix of that reversed word is a palindrome, append the word index to
//     the node's palBelow. Mark wordIdx at the terminal node.
//  2. For each word i, walk forward through the trie:
//     - At each step, if the current node terminates word j (j != i) and the
//     rest of word i is a palindrome, add [i, j].
//     - After consuming all of word i, add [i, j] for every j in node.palBelow
//     (j != i).
//
// Time:  O(n · m^2) — insertion and search each do O(m) palindrome checks per
//
//	character across n words.
//
// Space: O(n · m) — trie nodes plus palBelow lists.
func trieApproach(words []string) [][]int {
	root := newTrieNode()

	// insert places word (reversed) into the trie, tagging palindromic suffixes.
	insert := func(word string, idx int) {
		node := root
		n := len(word)
		// Walk the REVERSED word by iterating the original from the end.
		for pos := n - 1; pos >= 0; pos-- {
			// The characters not yet consumed of the reversed word correspond to
			// word[0..pos]; if that prefix is a palindrome, then words placed here
			// can pair with a longer counterpart (this word is the shorter one).
			if isPalindrome(word, 0, pos) {
				node.palBelow = append(node.palBelow, idx)
			}
			c := word[pos] - 'a'
			if node.children[c] == nil {
				node.children[c] = newTrieNode()
			}
			node = node.children[c]
		}
		node.palBelow = append(node.palBelow, idx) // empty remaining suffix is a palindrome
		node.wordIdx = idx                         // reversed word terminates here
	}
	for i, w := range words {
		insert(w, i)
	}

	res := [][]int{}
	// search walks word i forward through the reversed-word trie.
	search := func(word string, idx int) {
		node := root
		n := len(word)
		for k := 0; k < n; k++ {
			// If a (reversed) word j ends here and the REST of word i (word[k:])
			// is a palindrome, then word i + word j is a palindrome (j shorter).
			if node.wordIdx != -1 && node.wordIdx != idx && isPalindrome(word, k, n-1) {
				res = append(res, []int{idx, node.wordIdx})
			}
			c := word[k] - 'a'
			if node.children[c] == nil {
				return // no reversed word continues this way; dead end
			}
			node = node.children[c]
		}
		// Word i fully matched a path; any word j below whose remaining suffix is
		// a palindrome (palBelow) pairs as word i + word j (j longer or equal).
		for _, j := range node.palBelow {
			if j != idx {
				res = append(res, []int{idx, j})
			}
		}
	}
	for i, w := range words {
		search(w, i)
	}
	return res
}

// sortPairs canonicalises a result list (sort the pairs) so different valid
// orderings compare equal when we print them.
func sortPairs(p [][]int) [][]int {
	cp := make([][]int, len(p))
	copy(cp, p)
	sort.Slice(cp, func(a, b int) bool {
		if cp[a][0] != cp[b][0] {
			return cp[a][0] < cp[b][0]
		}
		return cp[a][1] < cp[b][1]
	})
	return cp
}

func main() {
	// Example 1: ["abcd","dcba","lls","s","sssll"]
	//   → [[0,1],[1,0],[3,2],[2,4]]  (sorted: [[0,1],[1,0],[2,4],[3,2]])
	ex1 := []string{"abcd", "dcba", "lls", "s", "sssll"}
	// Example 2: ["bat","tab","cat"] → [[0,1],[1,0]]
	ex2 := []string{"bat", "tab", "cat"}
	// Example 3: ["a",""] → [[0,1],[1,0]]
	ex3 := []string{"a", ""}

	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(sortPairs(bruteForce(ex1))) // [[0 1] [1 0] [2 4] [3 2]]
	fmt.Println(sortPairs(bruteForce(ex2))) // [[0 1] [1 0]]
	fmt.Println(sortPairs(bruteForce(ex3))) // [[0 1] [1 0]]

	fmt.Println("=== Approach 2: Hash Map of Reversed Words (Optimal) ===")
	fmt.Println(sortPairs(hashMap(ex1))) // [[0 1] [1 0] [2 4] [3 2]]
	fmt.Println(sortPairs(hashMap(ex2))) // [[0 1] [1 0]]
	fmt.Println(sortPairs(hashMap(ex3))) // [[0 1] [1 0]]

	fmt.Println("=== Approach 3: Trie of Reversed Words (Advanced Optimal) ===")
	fmt.Println(sortPairs(trieApproach(ex1))) // [[0 1] [1 0] [2 4] [3 2]]
	fmt.Println(sortPairs(trieApproach(ex2))) // [[0 1] [1 0]]
	fmt.Println(sortPairs(trieApproach(ex3))) // [[0 1] [1 0]]
}
