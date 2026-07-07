package main

import (
	"fmt"
	"strings"
)

// trieADT is the common interface every approach implements so that main()
// can drive the same official example through all of them.
type trieADT interface {
	Insert(word string)
	Search(word string) bool
	StartsWith(prefix string) bool
}

// ── Approach 1: Brute Force (Word List Scan) ─────────────────────────────────
//
// BruteForceTrie solves Implement Trie by storing every inserted word in a
// plain slice and answering queries with a linear scan.
//
// Intuition:
//
//	The dumbest data structure that satisfies the API is "a list of the
//	words". Search is an equality scan; StartsWith is a prefix scan. It is
//	trivially correct and shows exactly what cost the real trie eliminates:
//	re-comparing the same shared prefixes against every stored word.
//
// Algorithm:
//
//	Insert:     append the word to the slice (skip if already present).
//	Search:     scan the slice for an exact string match.
//	StartsWith: scan the slice testing strings.HasPrefix on each word.
//
// Time:  Insert O(n·m) (duplicate check), Search O(n·m), StartsWith O(n·m) —
//
//	n stored words, m the query length.
//
// Space: O(total characters stored) — every word kept verbatim.
type BruteForceTrie struct {
	words []string // every distinct word inserted so far, in arrival order
}

// NewBruteForceTrie builds the empty brute-force structure.
func NewBruteForceTrie() *BruteForceTrie {
	return &BruteForceTrie{words: []string{}}
}

// Insert stores word (once) in the word list.
func (t *BruteForceTrie) Insert(word string) {
	for _, w := range t.words {
		if w == word {
			return // already stored — the structure is a set, not a bag
		}
	}
	t.words = append(t.words, word) // first time seen: remember it
}

// Search reports whether word was previously inserted (exact match).
func (t *BruteForceTrie) Search(word string) bool {
	for _, w := range t.words {
		if w == word {
			return true // exact match found
		}
	}
	return false // scanned everything without a hit
}

// StartsWith reports whether any inserted word begins with prefix.
func (t *BruteForceTrie) StartsWith(prefix string) bool {
	for _, w := range t.words {
		if strings.HasPrefix(w, prefix) {
			return true // some stored word extends this prefix
		}
	}
	return false // no stored word starts with prefix
}

// ── Approach 2: Hash Sets of Words + Prefixes ────────────────────────────────
//
// PrefixSetTrie solves Implement Trie by pre-registering every prefix of
// every inserted word in a hash set, making both queries O(m) lookups.
//
// Intuition:
//
//	Queries dominate? Pay at insert time instead. A word of length m has m
//	prefixes — dump them all into a "prefixes" set, and keep a second
//	"words" set for exact membership. Both queries become single hash
//	lookups. The price: inserting a word of length m hashes m strings whose
//	total length is O(m²), and memory can blow up the same way — the exact
//	waste a trie's shared nodes avoid.
//
// Algorithm:
//
//	Insert:     add word to the words set; add word[:1], word[:2], …,
//	            word[:m] to the prefixes set.
//	Search:     hash-lookup in the words set.
//	StartsWith: hash-lookup in the prefixes set.
//
// Time:  Insert O(m²) (m prefixes of average length m/2), Search O(m),
//
//	StartsWith O(m) — hashing a string of length m costs O(m).
//
// Space: O(Σ m²) worst case — every prefix of every word stored as its own key.
type PrefixSetTrie struct {
	words    map[string]bool // complete inserted words
	prefixes map[string]bool // every prefix (length ≥ 1) of every inserted word
}

// NewPrefixSetTrie builds the empty two-set structure.
func NewPrefixSetTrie() *PrefixSetTrie {
	return &PrefixSetTrie{
		words:    map[string]bool{},
		prefixes: map[string]bool{},
	}
}

// Insert registers the word and all of its prefixes.
func (t *PrefixSetTrie) Insert(word string) {
	t.words[word] = true // exact word membership
	for i := 1; i <= len(word); i++ {
		t.prefixes[word[:i]] = true // every leading slice is now a known prefix
	}
}

// Search reports exact membership via one hash lookup.
func (t *PrefixSetTrie) Search(word string) bool {
	return t.words[word]
}

// StartsWith reports prefix membership via one hash lookup.
func (t *PrefixSetTrie) StartsWith(prefix string) bool {
	return t.prefixes[prefix]
}

// ── Approach 3: Trie with Map Children ───────────────────────────────────────
//
// MapTrie solves Implement Trie with a real prefix tree whose children are
// held in a map[byte]*node.
//
// Intuition:
//
//	Store each word as a root-to-node path, one character per edge. Words
//	sharing a prefix share the path for that prefix, so common prefixes are
//	stored once — the core trie idea. A per-node hash map keeps children
//	sparse: memory is proportional to edges actually used, and the same code
//	handles any alphabet (unicode bytes, digits, …), at the cost of hashing
//	overhead per step versus a fixed array.
//
// Algorithm:
//
//	Insert:     walk from the root, creating missing child nodes per
//	            character; mark the final node isEnd = true.
//	Search:     walk the path; true iff every child exists AND the final
//	            node has isEnd set.
//	StartsWith: walk the path; true iff every child exists (isEnd ignored).
//
// Time:  Insert O(m), Search O(m), StartsWith O(m) — one map hop per character.
// Space: O(total characters across unique prefixes) — shared prefixes stored once.
type MapTrie struct {
	children map[byte]*MapTrie // outgoing edges keyed by character
	isEnd    bool              // true if a complete word terminates at this node
}

// NewMapTrie builds an empty root node.
func NewMapTrie() *MapTrie {
	return &MapTrie{children: map[byte]*MapTrie{}}
}

// Insert adds word to the trie, creating nodes on demand.
func (t *MapTrie) Insert(word string) {
	cur := t // start walking at the root
	for i := 0; i < len(word); i++ {
		c := word[i]
		if cur.children[c] == nil {
			cur.children[c] = NewMapTrie() // first word through here: create the edge
		}
		cur = cur.children[c] // descend along the character edge
	}
	cur.isEnd = true // the path root→cur spells exactly `word`
}

// walk descends the trie along s, returning the final node or nil if the
// path breaks (shared by Search and StartsWith).
func (t *MapTrie) walk(s string) *MapTrie {
	cur := t
	for i := 0; i < len(s); i++ {
		cur = cur.children[s[i]] // follow the edge for this character
		if cur == nil {
			return nil // path missing → nothing stored starts with s[:i+1]
		}
	}
	return cur
}

// Search is true iff the full path exists and ends on a word marker.
func (t *MapTrie) Search(word string) bool {
	n := t.walk(word)
	return n != nil && n.isEnd // node must exist AND terminate a stored word
}

// StartsWith is true iff the full path exists (word marker irrelevant).
func (t *MapTrie) StartsWith(prefix string) bool {
	return t.walk(prefix) != nil
}

// ── Approach 4: Trie with Array Children (Optimal) ───────────────────────────
//
// ArrayTrie solves Implement Trie with a prefix tree whose children live in a
// fixed [26]*node array indexed by letter — the canonical interview trie.
//
// Intuition:
//
//	The alphabet is exactly 'a'..'z' (26 letters, per the constraints), so a
//	child pointer array indexed by c-'a' replaces the hash map: O(1) child
//	access with zero hashing, better cache behaviour, and the least code.
//	This is the textbook trie and the version to write in interviews.
//
// Algorithm:
//
//	Insert:     for each character, index children[c-'a'], allocating the
//	            node if nil; mark the last node isEnd.
//	Search:     follow the indices; true iff no nil hit and final isEnd.
//	StartsWith: follow the indices; true iff no nil hit.
//
// Time:  Insert O(m), Search O(m), StartsWith O(m) — one array index per character.
// Space: O(unique prefix characters × 26) pointers — constant-factor tradeoff
//
//	versus the map for guaranteed O(1) steps.
type ArrayTrie struct {
	children [26]*ArrayTrie // children[i] = subtree for letter 'a'+i
	isEnd    bool           // true if a complete word terminates here
}

// NewArrayTrie builds an empty root node.
func NewArrayTrie() *ArrayTrie {
	return &ArrayTrie{}
}

// Insert adds word to the trie, allocating child nodes on demand.
func (t *ArrayTrie) Insert(word string) {
	cur := t // start at the root
	for i := 0; i < len(word); i++ {
		idx := word[i] - 'a' // map letter to slot 0..25
		if cur.children[idx] == nil {
			cur.children[idx] = NewArrayTrie() // create the missing edge
		}
		cur = cur.children[idx] // descend one level
	}
	cur.isEnd = true // mark that a whole word ends on this node
}

// walk descends along s and returns the final node, or nil if the path breaks.
func (t *ArrayTrie) walk(s string) *ArrayTrie {
	cur := t
	for i := 0; i < len(s); i++ {
		cur = cur.children[s[i]-'a'] // O(1) array hop for this character
		if cur == nil {
			return nil // no stored word continues this way
		}
	}
	return cur
}

// Search is true iff the path exists and the final node ends a word.
func (t *ArrayTrie) Search(word string) bool {
	n := t.walk(word)
	return n != nil && n.isEnd // "app" fails on apple-only tries: node exists, isEnd false
}

// StartsWith is true iff the path exists at all.
func (t *ArrayTrie) StartsWith(prefix string) bool {
	return t.walk(prefix) != nil
}

// runExample drives the single official LeetCode example through one
// implementation and returns the output list in LeetCode's format.
//
// Ops:  ["Trie","insert","search","search","startsWith","insert","search"]
// Args: [[],["apple"],["apple"],["app"],["app"],["app"],["app"]]
func runExample(newTrie func() trieADT) string {
	t := newTrie()          // "Trie"        → null
	out := []string{"null"} // constructor produces no value
	t.Insert("apple")       // insert("apple") → null
	out = append(out, "null")
	out = append(out, fmt.Sprintf("%t", t.Search("apple")))   // → true
	out = append(out, fmt.Sprintf("%t", t.Search("app")))     // → false
	out = append(out, fmt.Sprintf("%t", t.StartsWith("app"))) // → true
	t.Insert("app")                                           // insert("app")  → null
	out = append(out, "null")
	out = append(out, fmt.Sprintf("%t", t.Search("app"))) // → true
	return "[" + strings.Join(out, ", ") + "]"
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Word List Scan) ===")
	fmt.Println(runExample(func() trieADT { return NewBruteForceTrie() })) // [null, null, true, false, true, null, true]

	fmt.Println("=== Approach 2: Hash Sets of Words + Prefixes ===")
	fmt.Println(runExample(func() trieADT { return NewPrefixSetTrie() })) // [null, null, true, false, true, null, true]

	fmt.Println("=== Approach 3: Trie with Map Children ===")
	fmt.Println(runExample(func() trieADT { return NewMapTrie() })) // [null, null, true, false, true, null, true]

	fmt.Println("=== Approach 4: Trie with Array Children (Optimal) ===")
	fmt.Println(runExample(func() trieADT { return NewArrayTrie() })) // [null, null, true, false, true, null, true]
}
