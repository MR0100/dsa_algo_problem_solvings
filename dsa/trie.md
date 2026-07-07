# Trie (Prefix Tree)

> **Also known as:** prefix tree, digital tree, radix tree (compressed variant)
> **Core operations:** `Insert`, `Search`, `StartsWith` — all in O(L) where L = key length

---

## What it is

A **trie** is a tree that stores a *set of strings* by sharing common prefixes.
Each node represents a prefix; each edge is labelled with a single character.
Walking from the root along edges spells out a string one character at a time.

```
Insert: "car", "card", "care", "dog"

            (root)
            /    \
           c      d
           |      |
           a      o
           |      |
           r*     g*        * = isEnd marker (a complete word ends here)
          / \
         d*  e*
```

Key properties:

- All strings sharing a prefix share the **same path** from the root — the
  prefix is stored exactly once, no matter how many words contain it.
- Lookup cost depends on the **length of the key**, not on how many keys are
  stored. Searching among a million words costs the same O(L) as among ten.
- A node does *not* store its own character; the **edge** (i.e. the index/key
  in the children map) carries the character. A node stores:
  1. its children (array or map), and
  2. a flag `isEnd` — "some inserted word terminates exactly here".
  Optionally: a count, the full word, or any payload (value for a map-like trie).

### Trie vs. hash map — when each wins

| Question | Hash map | Trie |
|---|---|---|
| "Is exact word W present?" | O(L) average, simpler | O(L), no hashing |
| "Is any word with prefix P present?" | O(n·L) — must scan all keys | **O(P)** — walk the path |
| "List all words with prefix P" | O(n·L) | O(P + output) |
| "Match with wildcards (`.` = any char)?" | impossible without scan | DFS over branches |
| "Longest common prefix / shortest unique prefix?" | awkward | natural — follow single-child chain |
| Memory | compact | heavier (pointer-rich), mitigated by arrays / compression |

**Rule of thumb:** the moment a problem is about *prefixes* rather than *whole
keys*, a hash map stops helping and a trie becomes the right tool.

---

## How to recognise a trie problem — signals in the statement

Reach for a trie when you see any of these:

1. **The word "prefix" appears** — "starts with", "common prefix", "prefix of
   another word", "autocomplete", "search suggestions".
2. **A dictionary of words + repeated queries against it** — "given a list of
   words, answer Q queries of the form …". Building the trie once (O(total
   chars)) amortises across queries.
3. **Multi-pattern search** — find *many* words simultaneously inside a grid or
   text (Word Search II). A trie lets one DFS advance through *all* patterns at
   once instead of running one search per word.
4. **Wildcard / partial matching** — patterns containing `.` or "at most one
   character may differ". The trie's branching structure makes "try every
   child" a clean DFS.
5. **Design an autocomplete / spell-checker / search-suggestion system** —
   trie is the canonical interview answer.
6. **Maximum XOR of pairs of numbers** — the *binary trie* (bitwise trie):
   insert numbers as 32-bit paths, then greedily walk toward the opposite bit.
   Any "choose bits greedily from most-significant down" problem is a trie in
   disguise.
7. **"Replace words with their shortest root"** / stemming — walk each query
   word until the first `isEnd`.
8. **Constraints mention lowercase letters and total characters ≤ ~10⁵–10⁶** —
   a hint that O(total chars) preprocessing is intended.

Counter-signals (trie is overkill): only exact-match lookups (hash map), a
single pattern in a single text (KMP / Z-algorithm), or substring — not
prefix — queries (suffix automaton / suffix array territory).

---

## General templates (Go)

### 1. Classic lowercase-letter trie (array children — fastest)

```go
// TrieNode holds 26 child pointers (one per lowercase letter) and an
// end-of-word marker. Array indexing beats a map when the alphabet is
// small and fixed — no hashing, cache-friendly.
type TrieNode struct {
	children [26]*TrieNode // children[c-'a'] is the child reached by letter c
	isEnd    bool          // true ⇢ an inserted word terminates at this node
}

// Trie wraps the root. The root represents the empty prefix "".
type Trie struct {
	root *TrieNode
}

func NewTrie() *Trie {
	return &Trie{root: &TrieNode{}}
}

// Insert adds word to the trie.
//
// Pseudocode:
//   node ← root
//   for each character c in word:
//       if node has no child for c: create it
//       node ← child for c
//   mark node as end-of-word
//
// Time:  O(L) — one step per character.
// Space: O(L) worst case — at most L new nodes.
func (t *Trie) Insert(word string) {
	node := t.root
	for i := 0; i < len(word); i++ {
		c := word[i] - 'a'              // map byte to index 0..25
		if node.children[c] == nil {    // path doesn't exist yet
			node.children[c] = &TrieNode{} // create the missing node
		}
		node = node.children[c] // descend
	}
	node.isEnd = true // the full word ends exactly here
}

// walk descends along prefix and returns the final node, or nil if the
// path breaks. Shared helper for Search and StartsWith.
func (t *Trie) walk(prefix string) *TrieNode {
	node := t.root
	for i := 0; i < len(prefix); i++ {
		c := prefix[i] - 'a'
		if node.children[c] == nil {
			return nil // path breaks ⇒ no stored word has this prefix
		}
		node = node.children[c]
	}
	return node
}

// Search reports whether the exact word was inserted.
// Time: O(L). Space: O(1).
func (t *Trie) Search(word string) bool {
	node := t.walk(word)
	return node != nil && node.isEnd // must land on a node AND be a word end
}

// StartsWith reports whether any inserted word has the given prefix.
// Time: O(L). Space: O(1).
func (t *Trie) StartsWith(prefix string) bool {
	return t.walk(prefix) != nil // landing anywhere on the path is enough
}
```

### 2. Map-children variant (arbitrary alphabet / Unicode)

```go
// Use when the alphabet is large, sparse, or unknown (Unicode, digits+letters).
// Trades a little speed for flexibility and memory on sparse alphabets.
type Node struct {
	children map[rune]*Node
	isEnd    bool
}

func newNode() *Node { return &Node{children: map[rune]*Node{}} }

func (n *Node) insert(word string) {
	cur := n
	for _, r := range word { // ranging over string yields runes (Unicode-safe)
		next, ok := cur.children[r]
		if !ok {
			next = newNode()
			cur.children[r] = next
		}
		cur = next
	}
	cur.isEnd = true
}
```

### 3. Wildcard search (`.` matches any single character) — DFS

```go
// searchPattern supports '.' as a one-character wildcard (LeetCode #211).
//
// Pseudocode:
//   dfs(node, i):
//       if i == len(pattern): return node.isEnd
//       if pattern[i] == '.': try EVERY non-nil child
//       else:                 follow the single matching child
//
// Time:  O(L) without wildcards; O(26^d · L) worst case with d dots.
func (t *Trie) SearchPattern(pattern string) bool {
	var dfs func(node *TrieNode, i int) bool
	dfs = func(node *TrieNode, i int) bool {
		if node == nil {
			return false // fell off the tree
		}
		if i == len(pattern) {
			return node.isEnd // consumed pattern: need a word ending here
		}
		if pattern[i] == '.' {
			for _, child := range node.children { // branch over all children
				if dfs(child, i+1) {
					return true // any branch succeeding is enough
				}
			}
			return false
		}
		return dfs(node.children[pattern[i]-'a'], i+1) // exact character step
	}
	return dfs(t.root, 0)
}
```

### 4. Binary trie for maximum XOR (LeetCode #421 pattern)

```go
// bitTrieNode has exactly two children: bit 0 and bit 1.
type bitTrieNode struct {
	children [2]*bitTrieNode
}

// insertNum stores the 32-bit big-endian path of num.
func insertNum(root *bitTrieNode, num int) {
	node := root
	for b := 31; b >= 0; b-- { // most-significant bit first — greedy order
		bit := (num >> b) & 1
		if node.children[bit] == nil {
			node.children[bit] = &bitTrieNode{}
		}
		node = node.children[bit]
	}
}

// maxXorWith greedily walks toward the OPPOSITE bit at each level,
// because a differing bit at position b contributes 2^b to the XOR —
// and taking it always beats anything the lower bits can add.
func maxXorWith(root *bitTrieNode, num int) int {
	node, res := root, 0
	for b := 31; b >= 0; b-- {
		bit := (num >> b) & 1
		want := 1 - bit                     // opposite bit maximises XOR
		if node.children[want] != nil {
			res |= 1 << b                   // we win this bit
			node = node.children[want]
		} else {
			node = node.children[bit]       // forced to match; bit stays 0
		}
	}
	return res
}
```

---

## Worked example — trace of Insert/Search/StartsWith

Operations (LeetCode #208's example):

```
Insert("apple")
Search("apple")   → true
Search("app")     → false
StartsWith("app") → true
Insert("app")
Search("app")     → true
```

**Step 1 — `Insert("apple")`.** Start at root; every child slot is nil, so
each character creates a node:

| i | char | child exists? | action | node after step |
|---|------|---------------|--------|-----------------|
| 0 | `a` | no | create node A | A |
| 1 | `p` | no | create node P1 | P1 |
| 2 | `p` | no | create node P2 | P2 |
| 3 | `l` | no | create node L | L |
| 4 | `e` | no | create node E | E |

Then set `E.isEnd = true`. Tree: `root → a → p → p → l → e*`.

**Step 2 — `Search("apple")`.** Walk `a,p,p,l,e` — every child exists, land on
E. `E.isEnd == true` → **true**.

**Step 3 — `Search("app")`.** Walk `a,p,p` — path exists, land on P2. But
`P2.isEnd == false` (no word *ends* there; "app" is only a prefix of "apple")
→ **false**. *This is the exact behavioural difference between Search and
StartsWith.*

**Step 4 — `StartsWith("app")`.** Walk `a,p,p` — path exists, land on P2.
StartsWith doesn't check `isEnd` → **true**.

**Step 5 — `Insert("app")`.** Walk `a,p,p` — all three children already exist,
so **zero new nodes** are allocated (prefix sharing at work). Set
`P2.isEnd = true`. Tree: `root → a → p* → p → l → e*` wait — careful: the
marker goes on the node reached after consuming all of "app", i.e. **P2**:
`root → a → p → p* → l → e*`.

**Step 6 — `Search("app")`.** Land on P2, now `isEnd == true` → **true**.

Complexities for a trie holding words of total length N, alphabet size Σ:

- Build: **O(N)** time, **O(N·Σ)** worst-case space (array children).
- Each query (Search / StartsWith / Insert): **O(L)** — independent of the
  number of stored words.

---

## Common pitfalls and how to avoid them

1. **Confusing `Search` and `StartsWith`.** `Search` must check `isEnd`;
   `StartsWith` must not. Forgetting `isEnd` makes `Search("app")` wrongly
   return true after only `Insert("apple")`. Write the shared `walk` helper
   and keep the `isEnd` check only in `Search`.
2. **Storing the character in the node instead of the edge.** The child's
   *index/key* already encodes the character; duplicating it wastes memory and
   invites inconsistency. The root corresponds to the empty string.
3. **Marking `isEnd` on the wrong node.** It goes on the node reached *after*
   consuming the last character — never on the parent, never at loop start.
   Off-by-one here makes every word register one character short.
4. **Forgetting to initialise map children.** With the map variant, a zero-value
   `Node{}` has a nil map; writing to it panics. Always construct via
   `newNode()` (or lazily allocate before first write). The array variant has
   no such problem — another reason to prefer it for `a–z`.
5. **Byte vs rune indexing.** `word[i] - 'a'` is byte-based and only valid for
   ASCII. For Unicode input, `for _, r := range word` with map children.
   Mixing the two silently corrupts multi-byte input.
6. **Memory blow-up on large alphabets.** `[26]*Node` per node is fine; a
   `[128]` or `[65536]` array is not. Switch to map children, or compress
   chains of single-child nodes (radix/Patricia trie) if memory is tight.
7. **Wildcard DFS without a nil guard.** In `SearchPattern`, recursing into
   `node.children[c]` may pass nil — handle `node == nil` at the top of the
   DFS or check before recursing, or you'll dereference nil.
8. **Not pruning in trie + grid backtracking (Word Search II).** Two classic
   optimisations: (a) after a word is found, clear its `isEnd` (or remove the
   word) to avoid duplicates; (b) delete leaf nodes once exhausted so later
   DFS paths die immediately. Skipping these turns an O(cells · 3^L) search
   into a TLE.
9. **Rebuilding the trie per query.** The whole point is: build once —
   O(total chars) — then answer many O(L) queries. Building inside the query
   loop throws away the asymptotic win.
10. **Recursion depth on very long keys.** Insert/Search are naturally
    iterative — keep them so. Reserve recursion for genuinely branching
    operations (wildcards, collect-all-words), where depth ≤ max word length.

---

## Problems in this repo

Problems currently in the repo that touch the trie concept:

- [0014 — Longest Common Prefix](/0014_longest_common_prefix/README.md) —
  tagged **Trie**: the LCP of a word set is exactly the chain of single-child,
  non-terminal nodes from a trie's root; solved there by simpler scans, with
  the trie as the conceptual model.
- [0079 — Word Search](/0079_word_search/README.md) — single-word grid
  backtracking; its multi-word extension (LeetCode #212, Word Search II) is
  the canonical **trie + backtracking** problem, noted in that README's
  Related Problems.

> Note: problems 0131–0400 are being written concurrently; a later pass will
> add the core trie problems (#208 Implement Trie, #211 Design Add and Search
> Words, #212 Word Search II, #336 Palindrome Pairs, #421 Maximum XOR) once
> their folders exist.
