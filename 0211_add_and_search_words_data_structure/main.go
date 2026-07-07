package main

import (
	"fmt"
	"strings"
)

// wordDictionaryADT is the common contract every approach implements so main()
// can drive the same official example through all of them.
type wordDictionaryADT interface {
	AddWord(word string)
	Search(word string) bool
}

// ── Approach 1: Brute Force (Grouped Word List + Pattern Scan) ────────────────
//
// BruteForceDict solves Add and Search Words by keeping every added word in a
// slice bucketed by length, and answering search with a linear pattern scan.
//
// Intuition:
//
//	The '.' wildcard matches any single letter, so a stored word matches a
//	query iff they have the same length and agree on every non-'.' position.
//	The dumbest correct structure is just "the list of words". Bucketing by
//	length prunes the obvious mismatches (only words of the query's length can
//	match) but each search still re-compares the query against many words.
//
// Algorithm:
//
//	AddWord: append word to buckets[len(word)].
//	Search:  for each stored word of the same length, compare character by
//	         character, treating '.' in the query as a match; return true on
//	         the first fully-matching word.
//
// Time:  AddWord O(1) amortised; Search O(n·m) — n words of the query length,
//
//	m the query length.
//
// Space: O(total characters stored).
type BruteForceDict struct {
	buckets map[int][]string // buckets[L] = every added word of length L
}

// NewBruteForceDict builds the empty brute-force structure.
func NewBruteForceDict() *BruteForceDict {
	return &BruteForceDict{buckets: map[int][]string{}}
}

// AddWord files word into the bucket for its length.
func (d *BruteForceDict) AddWord(word string) {
	d.buckets[len(word)] = append(d.buckets[len(word)], word) // O(1) amortised append
}

// matches reports whether stored word `w` fits the (possibly wildcarded) pattern.
func matches(pattern, w string) bool {
	for i := 0; i < len(pattern); i++ {
		if pattern[i] != '.' && pattern[i] != w[i] {
			return false // concrete letter mismatch → this word cannot match
		}
	}
	return true // every position matched (or was a '.')
}

// Search scans only the same-length bucket for a matching word.
func (d *BruteForceDict) Search(word string) bool {
	for _, w := range d.buckets[len(word)] { // only equal-length words can match
		if matches(word, w) {
			return true // found one word consistent with the pattern
		}
	}
	return false // scanned the whole bucket without a match
}

// ── Approach 2: Trie with Map Children + Wildcard DFS ─────────────────────────
//
// MapTrieDict solves Add and Search Words with a prefix tree whose children are
// a map[byte]*node, using DFS to branch over all children at each '.'.
//
// Intuition:
//
//	AddWord is a plain trie insert. Search is a trie walk, except a '.' must
//	try EVERY existing child at that node (any letter is allowed), which turns
//	the walk into a bounded DFS/backtracking. A concrete letter still follows a
//	single edge, so only '.' positions branch. The map keeps nodes sparse.
//
// Algorithm:
//
//	AddWord: standard trie insert, mark final node isEnd.
//	Search(word, node, i):
//	  - if i == len(word): return node.isEnd (a whole word ends here).
//	  - if word[i] == '.': recurse into every child; true if any subtree matches.
//	  - else: follow the single child edge for word[i] (fail if absent).
//
// Time:  AddWord O(m). Search O(m) with no wildcards, up to O(26^k · m) worst
//
//	case where k = number of '.' (each dot fans out over ≤ 26 children).
//
// Space: O(total characters over unique prefixes) + O(m) recursion depth.
type mapNode struct {
	children map[byte]*mapNode // outgoing edges keyed by character
	isEnd    bool              // a complete word terminates at this node
}

// MapTrieDict is the trie root wrapped so it satisfies wordDictionaryADT.
type MapTrieDict struct {
	root *mapNode
}

// NewMapTrieDict builds an empty map-backed trie dictionary.
func NewMapTrieDict() *MapTrieDict {
	return &MapTrieDict{root: &mapNode{children: map[byte]*mapNode{}}}
}

// AddWord inserts word into the trie, creating nodes on demand.
func (d *MapTrieDict) AddWord(word string) {
	cur := d.root
	for i := 0; i < len(word); i++ {
		c := word[i]
		if cur.children[c] == nil {
			cur.children[c] = &mapNode{children: map[byte]*mapNode{}} // create missing edge
		}
		cur = cur.children[c] // descend one level
	}
	cur.isEnd = true // path root→cur spells exactly `word`
}

// Search kicks off the wildcard DFS from the root.
func (d *MapTrieDict) Search(word string) bool {
	return dfsMap(d.root, word, 0)
}

// dfsMap matches word[i:] starting at node, branching on '.'.
func dfsMap(node *mapNode, word string, i int) bool {
	if node == nil {
		return false // fell off the trie
	}
	if i == len(word) {
		return node.isEnd // consumed the pattern; match iff a word ends here
	}
	c := word[i]
	if c == '.' {
		for _, child := range node.children { // wildcard: try every existing edge
			if dfsMap(child, word, i+1) {
				return true // some branch completed the match
			}
		}
		return false // no child led to a match
	}
	return dfsMap(node.children[c], word, i+1) // concrete letter: single edge
}

// ── Approach 3: Trie with Array Children + Wildcard DFS (Optimal) ─────────────
//
// ArrayTrieDict solves Add and Search Words with a fixed [26]*node trie — the
// canonical interview answer. Identical DFS logic; child access is an O(1)
// array index instead of a hash lookup.
//
// Intuition:
//
//	The added words are lowercase 'a'..'z', so children fit in a [26] array
//	indexed by c-'a'. A concrete letter follows children[c-'a']; a '.' iterates
//	all 26 slots, recursing into the non-nil ones. No hashing, cache-friendly,
//	shortest to write correctly under interview pressure.
//
// Algorithm:
//
//	AddWord: array-trie insert, mark isEnd on the last node.
//	Search(word, node, i): same three cases as Approach 2, but '.' loops slots
//	0..25 and a letter indexes children[word[i]-'a'].
//
// Time:  AddWord O(m). Search O(m) with no wildcards; worst case O(26^k · m).
// Space: O(unique-prefix chars × 26) pointers + O(m) recursion depth.
type arrNode struct {
	children [26]*arrNode // children[i] = subtree for letter 'a'+i
	isEnd    bool         // a complete word terminates here
}

// ArrayTrieDict is the array-trie root satisfying wordDictionaryADT.
type ArrayTrieDict struct {
	root *arrNode
}

// NewArrayTrieDict builds an empty array-backed trie dictionary.
func NewArrayTrieDict() *ArrayTrieDict {
	return &ArrayTrieDict{root: &arrNode{}}
}

// AddWord inserts word into the array trie.
func (d *ArrayTrieDict) AddWord(word string) {
	cur := d.root
	for i := 0; i < len(word); i++ {
		idx := word[i] - 'a' // letter → slot 0..25
		if cur.children[idx] == nil {
			cur.children[idx] = &arrNode{} // create missing edge
		}
		cur = cur.children[idx] // descend
	}
	cur.isEnd = true // mark end of a stored word
}

// Search launches the wildcard DFS from the root.
func (d *ArrayTrieDict) Search(word string) bool {
	return dfsArr(d.root, word, 0)
}

// dfsArr matches word[i:] from node, branching over all 26 slots on '.'.
func dfsArr(node *arrNode, word string, i int) bool {
	if node == nil {
		return false // no such path
	}
	if i == len(word) {
		return node.isEnd // pattern consumed; must end on a word
	}
	c := word[i]
	if c == '.' {
		for j := 0; j < 26; j++ { // wildcard: try every populated slot
			if node.children[j] != nil && dfsArr(node.children[j], word, i+1) {
				return true
			}
		}
		return false
	}
	return dfsArr(node.children[c-'a'], word, i+1) // concrete letter: one slot
}

// runExample drives the single official example through one implementation and
// returns the output list in LeetCode's format.
//
// Ops:  ["WordDictionary","addWord","addWord","addWord","search","search","search","search"]
// Args: [[],["bad"],["dad"],["mad"],["pad"],["bad"],[".ad"],["b.."]]
func runExample(newDict func() wordDictionaryADT) string {
	d := newDict()
	out := []string{"null"} // constructor → null
	d.AddWord("bad")
	d.AddWord("dad")
	d.AddWord("mad")
	out = append(out, "null", "null", "null")             // three addWord → null each
	out = append(out, fmt.Sprintf("%t", d.Search("pad"))) // false
	out = append(out, fmt.Sprintf("%t", d.Search("bad"))) // true
	out = append(out, fmt.Sprintf("%t", d.Search(".ad"))) // true
	out = append(out, fmt.Sprintf("%t", d.Search("b.."))) // true
	return "[" + strings.Join(out, ", ") + "]"
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Grouped Word List) ===")
	fmt.Println(runExample(func() wordDictionaryADT { return NewBruteForceDict() })) // [null, null, null, null, false, true, true, true]

	fmt.Println("=== Approach 2: Trie with Map Children + Wildcard DFS ===")
	fmt.Println(runExample(func() wordDictionaryADT { return NewMapTrieDict() })) // [null, null, null, null, false, true, true, true]

	fmt.Println("=== Approach 3: Trie with Array Children + Wildcard DFS (Optimal) ===")
	fmt.Println(runExample(func() wordDictionaryADT { return NewArrayTrieDict() })) // [null, null, null, null, false, true, true, true]
}
