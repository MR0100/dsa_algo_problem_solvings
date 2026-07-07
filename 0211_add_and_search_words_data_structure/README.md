# 0211 — Design Add and Search Words Data Structure

> LeetCode #211 · Difficulty: Medium
> **Categories:** Trie, Design, Depth-First Search, Backtracking, String

---

## Problem Statement

Design a data structure that supports adding new words and finding if a string matches any previously added string.

Implement the `WordDictionary` class:

- `WordDictionary()` Initializes the object.
- `void addWord(word)` Adds `word` to the data structure, it can be matched later.
- `bool search(word)` Returns `true` if there is any string in the data structure that matches `word` or `false` otherwise. `word` may contain dots `'.'` where dots can be matched with any letter.

**Example 1:**
```
Input
["WordDictionary","addWord","addWord","addWord","search","search","search","search"]
[[],["bad"],["dad"],["mad"],["pad"],["bad"],[".ad"],["b.."]]
Output
[null,null,null,null,false,true,true,true]

Explanation
WordDictionary wordDictionary = new WordDictionary();
wordDictionary.addWord("bad");
wordDictionary.addWord("dad");
wordDictionary.addWord("mad");
wordDictionary.search("pad"); // return False
wordDictionary.search("bad"); // return True
wordDictionary.search(".ad"); // return True
wordDictionary.search("b.."); // return True
```

**Constraints:**
- `1 <= word.length <= 25`
- `word` in `addWord` consists of lowercase English letters.
- `word` in `search` consists of `'.'` or lowercase English letters.
- There will be at most `2` dots in `word` for `search` queries.
- At most `10⁴` calls will be made to `addWord` and `search`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★★ Very High  | 2024          |
| Meta       | ★★★★☆ High       | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2024          |
| Apple      | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Trie / Prefix Tree** — words are stored as root-to-node paths so shared prefixes are traversed once; the natural structure for prefix/pattern queries → see [`/dsa/trie.md`](/dsa/trie.md)
- **Depth-First Search + Backtracking** — a `'.'` wildcard branches over every child at that node, turning the walk into a bounded DFS → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Design / Data-Structure API** — a stateful class judged per-operation against a scripted call sequence → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)
- **String Handling** — byte-level indexing over the fixed `'a'..'z'` alphabet plus the `'.'` sentinel → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

Let n = number of added words, m = query length, k = number of `'.'` in a query.

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Grouped Word List) | AddWord O(1), Search O(n·m) | O(Σ word lengths) | Baseline; few words, short queries |
| 2 | Trie (Map children) + Wildcard DFS | AddWord O(m), Search O(m) → O(26ᵏ·m) | O(Σ unique prefix chars) | Sparse / large alphabets |
| 3 | Trie (Array children) + Wildcard DFS (Optimal) | AddWord O(m), Search O(m) → O(26ᵏ·m) | O(unique prefix chars × 26) | Fixed `a`–`z` alphabet — interview standard |

---

## Approach 1 — Brute Force (Grouped Word List)

### Intuition
A `'.'` matches any single letter, so a stored word matches a query iff they are the **same length** and agree on every non-`'.'` position. The simplest correct structure is just the list of added words. Bucketing by length prunes the trivially-wrong candidates (only words of the query's length can possibly match), but each search still re-compares the pattern against every word in that bucket.

### Algorithm
1. `AddWord(word)`: append `word` to `buckets[len(word)]`.
2. `Search(word)`: for each stored word of the same length, compare position by position; a `'.'` in the query matches anything, a concrete letter must equal the stored letter. Return `true` on the first fully-matching word.
3. If no bucket word matches, return `false`.

### Complexity
- **Time:** AddWord O(1) amortised (slice append); Search O(n·m) — up to n same-length words, each compared over m characters.
- **Space:** O(Σ word lengths) — every word stored verbatim.

### Code
```go
type BruteForceDict struct {
	buckets map[int][]string // buckets[L] = every added word of length L
}

func NewBruteForceDict() *BruteForceDict {
	return &BruteForceDict{buckets: map[int][]string{}}
}

func (d *BruteForceDict) AddWord(word string) {
	d.buckets[len(word)] = append(d.buckets[len(word)], word) // O(1) amortised append
}

func matches(pattern, w string) bool {
	for i := 0; i < len(pattern); i++ {
		if pattern[i] != '.' && pattern[i] != w[i] {
			return false // concrete letter mismatch → this word cannot match
		}
	}
	return true // every position matched (or was a '.')
}

func (d *BruteForceDict) Search(word string) bool {
	for _, w := range d.buckets[len(word)] { // only equal-length words can match
		if matches(word, w) {
			return true
		}
	}
	return false
}
```

### Dry Run (Example 1)

| # | Operation | `buckets` after | Scan | Return | Expected |
|---|-----------|-----------------|------|--------|----------|
| 1 | `WordDictionary()` | {} | — | null | null |
| 2 | `addWord("bad")` | {3:[bad]} | — | null | null |
| 3 | `addWord("dad")` | {3:[bad,dad]} | — | null | null |
| 4 | `addWord("mad")` | {3:[bad,dad,mad]} | — | null | null |
| 5 | `search("pad")` | (unchanged) | bad✗ dad✗ mad✗ (p≠b,d,m) | false | false |
| 6 | `search("bad")` | (unchanged) | bad✓ | true | true |
| 7 | `search(".ad")` | (unchanged) | bad: `.`✓ a✓ d✓ | true | true |
| 8 | `search("b..")` | (unchanged) | bad: b✓ `.`✓ `.`✓ | true | true |

Output `[null, null, null, null, false, true, true, true]` ✓

---

## Approach 2 — Trie (Map children) + Wildcard DFS

### Intuition
`addWord` is an ordinary trie insert. `search` is a trie walk, except a `'.'` must try **every existing child** of the current node (any letter is allowed), which converts the linear walk into a bounded DFS with backtracking. A concrete letter still follows exactly one edge, so only `'.'` positions branch. Holding children in a `map[byte]*node` keeps nodes sparse and generalises to any alphabet.

### Algorithm
1. `AddWord(word)`: walk from the root creating missing children; mark the final node `isEnd`.
2. `Search(word)` runs `dfs(root, word, 0)`:
   - If `i == len(word)`: return `node.isEnd`.
   - If `word[i] == '.'`: recurse into every child; return `true` if any subtree matches.
   - Else: recurse into the single child for `word[i]` (fail if absent).

### Complexity
- **Time:** AddWord O(m); Search O(m) with no wildcards, up to O(26ᵏ·m) worst case where k = number of dots (each dot fans out over ≤ 26 children).
- **Space:** O(total characters over unique prefixes) for the trie + O(m) recursion depth.

### Code
```go
type mapNode struct {
	children map[byte]*mapNode
	isEnd    bool
}

type MapTrieDict struct{ root *mapNode }

func NewMapTrieDict() *MapTrieDict {
	return &MapTrieDict{root: &mapNode{children: map[byte]*mapNode{}}}
}

func (d *MapTrieDict) AddWord(word string) {
	cur := d.root
	for i := 0; i < len(word); i++ {
		c := word[i]
		if cur.children[c] == nil {
			cur.children[c] = &mapNode{children: map[byte]*mapNode{}}
		}
		cur = cur.children[c]
	}
	cur.isEnd = true
}

func (d *MapTrieDict) Search(word string) bool { return dfsMap(d.root, word, 0) }

func dfsMap(node *mapNode, word string, i int) bool {
	if node == nil {
		return false
	}
	if i == len(word) {
		return node.isEnd
	}
	c := word[i]
	if c == '.' {
		for _, child := range node.children { // wildcard: try every edge
			if dfsMap(child, word, i+1) {
				return true
			}
		}
		return false
	}
	return dfsMap(node.children[c], word, i+1) // concrete letter: single edge
}
```

### Dry Run (Example 1)

Trie after the three inserts (nodes spelled by their path; `*` = isEnd):
`root → b→a→d*`, `root → d→a→d*`, `root → m→a→d*`.

| # | Operation | DFS trace | Return | Expected |
|---|-----------|-----------|--------|----------|
| 5 | `search("pad")` | root has no child `p` → dead end | false | false |
| 6 | `search("bad")` | b✓ → a✓ → d✓, node `bad` isEnd | true | true |
| 7 | `search(".ad")` | `.`: try children b,d,m → into `b`: a✓ d✓ isEnd | true | true |
| 8 | `search("b..")` | b✓ → `.`: child a → `.`: child d, isEnd | true | true |

Output `[null, null, null, null, false, true, true, true]` ✓

---

## Approach 3 — Trie (Array children) + Wildcard DFS (Optimal)

### Intuition
Added words are lowercase `'a'..'z'`, so each node's children fit a fixed `[26]*node` array indexed by `c-'a'`. A concrete letter follows `children[c-'a']` in O(1); a `'.'` iterates the 26 slots and recurses into the non-nil ones. No hashing, cache-friendly, and the least code to get right — the canonical interview answer. Logic is identical to Approach 2; only child storage changes.

### Algorithm
1. `AddWord(word)`: array-trie insert (`idx = c-'a'`), mark `isEnd` on the last node.
2. `Search(word)` runs `dfs(root, word, 0)` with the same three cases; on `'.'` loop slots `0..25` and recurse into populated ones.

### Complexity
- **Time:** AddWord O(m); Search O(m) with no wildcards, up to O(26ᵏ·m) worst case.
- **Space:** O(unique-prefix chars × 26) pointers + O(m) recursion depth.

### Code
```go
type arrNode struct {
	children [26]*arrNode
	isEnd    bool
}

type ArrayTrieDict struct{ root *arrNode }

func NewArrayTrieDict() *ArrayTrieDict { return &ArrayTrieDict{root: &arrNode{}} }

func (d *ArrayTrieDict) AddWord(word string) {
	cur := d.root
	for i := 0; i < len(word); i++ {
		idx := word[i] - 'a'
		if cur.children[idx] == nil {
			cur.children[idx] = &arrNode{}
		}
		cur = cur.children[idx]
	}
	cur.isEnd = true
}

func (d *ArrayTrieDict) Search(word string) bool { return dfsArr(d.root, word, 0) }

func dfsArr(node *arrNode, word string, i int) bool {
	if node == nil {
		return false
	}
	if i == len(word) {
		return node.isEnd
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
```

### Dry Run (Example 1)

Indices: b→1, a→0, d→3, m→12. Trie after inserts: root slots 1(`b`),3(`d`),12(`m`) each → a(0) → d(3)`*`.

| # | Operation | DFS trace | Return | Expected |
|---|-----------|-----------|--------|----------|
| 5 | `search("pad")` | `p`=slot 15 nil at root → dead end | false | false |
| 6 | `search("bad")` | slot 1✓ → 0✓ → 3✓ isEnd | true | true |
| 7 | `search(".ad")` | `.`: scan 0..25, slot 1 (`b`) → a✓ d✓ isEnd | true | true |
| 8 | `search("b..")` | slot 1✓ → `.`: slot 0 (`a`) → `.`: slot 3 (`d`) isEnd | true | true |

Output `[null, null, null, null, false, true, true, true]` ✓

---

## Key Takeaways

- **A `'.'` turns a trie walk into a DFS.** Concrete letters follow a single edge (no branching); only wildcards fan out. With ≤ 2 dots per query (per constraints) the 26ᵏ blow-up is a tiny constant, so it is effectively O(m).
- **`isEnd` distinguishes `search` from prefix presence.** The recursion base case must check `node.isEnd`, not just "node exists" — otherwise `".ad"` would wrongly match a prefix that no word terminates.
- **Array vs map children is the same trade as in #208:** `[26]` array for fixed small alphabets (fast, interview default) vs `map[byte]` for sparse/unbounded alphabets. The DFS logic is identical.
- **Design pattern:** define the API as an interface and drive all implementations through the same scripted example → a free regression harness.
- This is #208 Implement Trie plus a wildcard; internalising it makes #212 Word Search II (trie guiding a board DFS) feel routine.

---

## Related Problems

- LeetCode #208 — Implement Trie (Prefix Tree) (the wildcard-free base)
- LeetCode #212 — Word Search II (trie guiding a board DFS)
- LeetCode #648 — Replace Words (shortest-root lookup in a trie)
- LeetCode #677 — Map Sum Pairs (trie with subtree value aggregation)
- LeetCode #1032 — Stream of Characters (suffix trie queried per character)
- LeetCode #10 — Regular Expression Matching (`.` and `*` wildcards, DP flavour)
