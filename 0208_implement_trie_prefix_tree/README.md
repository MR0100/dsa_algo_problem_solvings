# 0208 — Implement Trie (Prefix Tree)

> LeetCode #208 · Difficulty: Medium
> **Categories:** Trie, Design, Hash Table, String

---

## Problem Statement

A **trie** (pronounced as "try") or **prefix tree** is a tree data structure used to efficiently store and retrieve keys in a dataset of strings. There are various applications of this data structure, such as autocomplete and spellchecker.

Implement the Trie class:

- `Trie()` Initializes the trie object.
- `void insert(String word)` Inserts the string `word` into the trie.
- `boolean search(String word)` Returns `true` if the string `word` is in the trie (i.e., was inserted before), and `false` otherwise.
- `boolean startsWith(String prefix)` Returns `true` if there is a previously inserted string `word` that has the prefix `prefix`, and `false` otherwise.

**Example 1:**
```
Input
["Trie", "insert", "search", "search", "startsWith", "insert", "search"]
[[], ["apple"], ["apple"], ["app"], ["app"], ["app"], ["app"]]
Output
[null, null, true, false, true, null, true]

Explanation
Trie trie = new Trie();
trie.insert("apple");
trie.search("apple");   // return True
trie.search("app");     // return False
trie.startsWith("app"); // return True
trie.insert("app");
trie.search("app");     // return True
```

**Constraints:**
- `1 <= word.length, prefix.length <= 2000`
- `word` and `prefix` consist only of lowercase English letters.
- At most `3 * 10⁴` calls **in total** will be made to `insert`, `search`, and `startsWith`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2024          |
| Meta       | ★★★☆☆ Medium     | 2024          |
| Apple      | ★★★☆☆ Medium     | 2023          |
| Uber       | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Trie / Prefix Tree** — the target structure itself: store words as root-to-node paths so shared prefixes are stored (and traversed) exactly once → see [`/dsa/trie.md`](/dsa/trie.md)
- **Design / Data-Structure API** — a stateful class with an operation contract, judged per-operation rather than per-algorithm → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)
- **Hash Map / Hash Set** — the pre-trie baselines (word set, prefix set) and one variant of child storage inside trie nodes → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **String Handling** — prefix testing and byte-level indexing over the fixed `'a'..'z'` alphabet → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

Let n = number of inserted words, m = length of the word/prefix argument.

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Word List Scan) | Insert O(n·m), Search O(n·m), StartsWith O(n·m) | O(Σ word lengths) | Baseline; only when calls are few and words short |
| 2 | Hash Sets of Words + Prefixes | Insert O(m²), Search O(m), StartsWith O(m) | O(Σ m²) worst case | Quick hack when inserts are rare and words short; shows the time/space trade |
| 3 | Trie with Map Children | Insert O(m), Search O(m), StartsWith O(m) | O(unique prefix chars) | Sparse/huge alphabets (unicode); most flexible |
| 4 | Trie with Array Children (Optimal) | Insert O(m), Search O(m), StartsWith O(m) | O(unique prefix chars × 26) | Fixed small alphabet — the interview-standard answer |

---

## Approach 1 — Brute Force (Word List Scan)

### Intuition
The dumbest structure satisfying the API is literally "a list of the inserted words". `search` scans for an exact match; `startsWith` scans testing `strings.HasPrefix`. It is obviously correct, and it exposes precisely the waste a trie removes: every query re-compares the same shared prefixes against every stored word, over and over. With up to 3×10⁴ calls and words up to 2000 chars, that is ~n·m = tens of millions of byte comparisons per *call* in the worst case.

### Algorithm
1. `Insert(word)`: linear-scan the slice; if `word` is absent, append it.
2. `Search(word)`: linear-scan; return `true` on the first exact string equality.
3. `StartsWith(prefix)`: linear-scan; return `true` on the first word where `strings.HasPrefix(w, prefix)` holds.

### Complexity
- **Time:** Insert O(n·m), Search O(n·m), StartsWith O(n·m) — each scan compares against up to n words, each comparison costing up to m bytes.
- **Space:** O(Σ word lengths) — the words stored verbatim, nothing else.

### Code
```go
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
```

### Dry Run (Example 1)

| # | Operation | `words` before | Scan result | Return | Expected |
|---|-----------|----------------|-------------|--------|----------|
| 1 | `Trie()` | — | — | null | null |
| 2 | `insert("apple")` | [] | not found → append | null | null |
| 3 | `search("apple")` | [apple] | "apple" == "apple" ✓ | true | true |
| 4 | `search("app")` | [apple] | "apple" ≠ "app" (no exact match) | false | false |
| 5 | `startsWith("app")` | [apple] | HasPrefix("apple","app") ✓ | true | true |
| 6 | `insert("app")` | [apple] | not found → append | null | null |
| 7 | `search("app")` | [apple app] | "app" == "app" ✓ | true | true |

Output `[null, null, true, false, true, null, true]` ✓

---

## Approach 2 — Hash Sets of Words + Prefixes

### Intuition
If queries dominate, pay at insert time instead. A word of length m has exactly m non-empty prefixes — register *all of them* in a `prefixes` hash set, and keep a second `words` set for exact membership. Both queries collapse to a single hash lookup. The catch: inserting a length-m word hashes m strings of total length m(m+1)/2 = O(m²), and memory can balloon identically. With words up to 2000 chars, one insert may store ~2 million characters of keys — the exact duplication a trie's shared nodes eliminate. This approach is the conceptual midpoint: it *precomputes prefix membership* like a trie, but without sharing.

### Algorithm
1. `Insert(word)`: add `word` to `words`; for `i = 1..len(word)` add `word[:i]` to `prefixes`.
2. `Search(word)`: return `words[word]`.
3. `StartsWith(prefix)`: return `prefixes[prefix]`.

### Complexity
- **Time:** Insert O(m²) — m prefix keys of average length m/2 must be hashed; Search O(m), StartsWith O(m) — hashing the single query string.
- **Space:** O(Σ m²) worst case — every prefix of every word stored as an independent key (no sharing between words).

### Code
```go
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
```

### Dry Run (Example 1)

| # | Operation | `words` after | `prefixes` after | Return | Expected |
|---|-----------|---------------|------------------|--------|----------|
| 1 | `Trie()` | {} | {} | null | null |
| 2 | `insert("apple")` | {apple} | {a, ap, app, appl, apple} | null | null |
| 3 | `search("apple")` | {apple} | (unchanged) | words["apple"] → true | true |
| 4 | `search("app")` | {apple} | (unchanged) | words["app"] → false | false |
| 5 | `startsWith("app")` | {apple} | (unchanged) | prefixes["app"] → true | true |
| 6 | `insert("app")` | {apple, app} | {a, ap, app, appl, apple} (all dups) | null | null |
| 7 | `search("app")` | {apple, app} | (unchanged) | words["app"] → true | true |

Output `[null, null, true, false, true, null, true]` ✓

---

## Approach 3 — Trie with Map Children

### Intuition
Now share the prefixes. Store each word as a **path from the root**, one character per edge; words with a common prefix walk the same nodes for that prefix, so it is stored once no matter how many words share it. Each node needs (a) outgoing edges and (b) an `isEnd` flag marking "a complete word terminates here" — the flag is what distinguishes `search` from `startsWith` (path `a→p→p` exists after inserting "apple", but no word *ends* there until "app" is inserted). Holding children in a `map[byte]*node` keeps nodes sparse and works for arbitrary alphabets, at the price of hashing on every hop.

### Algorithm
1. `Insert(word)`: start at the root; for each character create the missing child, descend; mark the final node `isEnd = true`.
2. `walk(s)` (shared helper): follow the child edge for each character; return `nil` the moment an edge is missing, else the final node.
3. `Search(word)`: `walk(word)` must return a non-nil node **with** `isEnd == true`.
4. `StartsWith(prefix)`: `walk(prefix)` must return a non-nil node (flag irrelevant).

### Complexity
- **Time:** Insert O(m), Search O(m), StartsWith O(m) — exactly one map lookup/creation per character.
- **Space:** O(total characters over all *unique* prefixes) — shared prefixes stored once; each node's map holds only the edges that exist.

### Code
```go
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
```

### Dry Run (Example 1)

Nodes written as the prefix they spell; `*` marks `isEnd = true`.

| # | Operation | Trie state after | Walk trace | Return | Expected |
|---|-----------|------------------|------------|--------|----------|
| 1 | `Trie()` | (root) | — | null | null |
| 2 | `insert("apple")` | root→a→ap→app→appl→apple* | creates 5 nodes | null | null |
| 3 | `search("apple")` | (unchanged) | a✓ p✓ p✓ l✓ e✓, node "apple" isEnd=**true** | true | true |
| 4 | `search("app")` | (unchanged) | a✓ p✓ p✓, node "app" exists but isEnd=**false** | false | false |
| 5 | `startsWith("app")` | (unchanged) | a✓ p✓ p✓, node exists (flag ignored) | true | true |
| 6 | `insert("app")` | root→a→ap→app*→appl→apple* | creates 0 nodes, flips "app".isEnd | null | null |
| 7 | `search("app")` | (unchanged) | a✓ p✓ p✓, node "app" isEnd=**true** | true | true |

Output `[null, null, true, false, true, null, true]` ✓

---

## Approach 4 — Trie with Array Children (Optimal)

### Intuition
The constraints promise lowercase `'a'..'z'` only — an alphabet of exactly 26. So replace each node's hash map with a fixed `[26]*node` array indexed by `c - 'a'`: child access becomes a raw array index (no hashing, no map overhead, cache-friendly), and the code gets *shorter*. This is the canonical trie every interviewer expects. The logic is identical to Approach 3 — only the child-storage representation changes.

### Algorithm
1. `Insert(word)`: for each character compute `idx = c - 'a'`; allocate `children[idx]` if nil; descend; mark the last node `isEnd = true`.
2. `walk(s)`: follow `children[s[i]-'a']` per character; `nil` on a missing edge.
3. `Search(word)`: `walk(word) != nil && walk-result.isEnd`.
4. `StartsWith(prefix)`: `walk(prefix) != nil`.

### Complexity
- **Time:** Insert O(m), Search O(m), StartsWith O(m) — one O(1) array hop per character; with ≤ 3×10⁴ calls and m ≤ 2000 this is at most ~6×10⁷ trivial steps overall.
- **Space:** O(unique-prefix characters × 26) pointers — each node reserves 26 slots even if one is used; the constant-factor price for guaranteed O(1) hops.

### Code
```go
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
```

### Dry Run (Example 1)

Indices: 'a'→0, 'p'→15, 'l'→11, 'e'→4. `*` marks `isEnd = true`.

| # | Operation | Nodes traversed (slot taken) | Allocation / flag | Return | Expected |
|---|-----------|------------------------------|-------------------|--------|----------|
| 1 | `Trie()` | root created, all 26 slots nil | — | null | null |
| 2 | `insert("apple")` | root─[0]→a─[15]→ap─[15]→app─[11]→appl─[4]→apple | 5 nodes allocated, "apple".isEnd = true | null | null |
| 3 | `search("apple")` | root→a→ap→app→appl→apple | node non-nil, isEnd true | true | true |
| 4 | `search("app")` | root→a→ap→app | node non-nil, isEnd **false** | false | false |
| 5 | `startsWith("app")` | root→a→ap→app | node non-nil (flag ignored) | true | true |
| 6 | `insert("app")` | root→a→ap→app (all slots already allocated) | 0 allocations, "app".isEnd = true | null | null |
| 7 | `search("app")` | root→a→ap→app | node non-nil, isEnd true | true | true |

Output `[null, null, true, false, true, null, true]` ✓

---

## Key Takeaways

- **`search` vs `startsWith` is one boolean.** A trie node existing means "some word passes through here"; only `isEnd` says "some word *ends* here". Forgetting the flag (or checking it in `startsWith`) is the classic bug — step 4 of the example exists precisely to catch it.
- **Trie = precomputed prefix membership with sharing.** The prefix-set approach (Approach 2) proves prefix lookups can be O(m); the trie gets the same query cost while storing each shared prefix once instead of once per word.
- **Child storage is a knob:** `[26]*node` array for small fixed alphabets (fastest, interview default) vs `map[byte]*node` for sparse or unbounded alphabets (unicode, digits). Logic is unchanged — factor the descent into a shared `walk` helper.
- All trie operations cost **O(m) in the key length, independent of how many words are stored** — the property that makes tries beat hash sets for prefix queries (a hash set simply cannot answer `startsWith` without enumerating).
- The trie built here is the skeleton for a whole family: add DFS over nodes → #212 Word Search II; add `.` wildcard matching → #211; store values at `isEnd` nodes → #677 Map Sum.
- Design-problem pattern: define the API as an interface and drive all implementations through the same scripted example — instant regression harness.

---

## Related Problems

- LeetCode #211 — Design Add and Search Words Data Structure (trie + `.` wildcard DFS)
- LeetCode #212 — Word Search II (trie guiding a board DFS)
- LeetCode #648 — Replace Words (shortest-root lookup in a trie)
- LeetCode #677 — Map Sum Pairs (trie with values aggregated over subtrees)
- LeetCode #1032 — Stream of Characters (suffix trie queried per incoming char)
- LeetCode #14 — Longest Common Prefix (degenerate prefix problem a trie also solves)
- LeetCode #421 — Maximum XOR of Two Numbers in an Array (bitwise trie)
