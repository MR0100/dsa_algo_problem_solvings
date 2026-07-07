# 0425 — Word Squares

> LeetCode #425 · Difficulty: Hard · 🔒 Premium
> **Categories:** Trie, Backtracking, Array, String, Hash Table

---

## Problem Statement

Given an array of **unique** strings `words`, return *all the* **[word squares](https://en.wikipedia.org/wiki/Word_square)** *you can build from* `words`. The same word from `words` can be used **multiple times**. You can return the answer in **any order**.

A sequence of strings forms a valid **word square** if the `k`th row and column read the same string, where `0 <= k < max(numRows, numColumns)`.

- For example, the word sequence `["ball","area","lead","lady"]` forms a word square because each word reads the same both horizontally and vertically.

**Example 1:**

```
Input: words = ["area","lead","wall","lady","ball"]
Output: [["ball","area","lead","lady"],["wall","area","lead","lady"]]
Explanation:
The output consists of two word squares. The order of output does not matter
(just the order of words in each word square matters).
```

**Example 2:**

```
Input: words = ["abat","baba","atan","atal"]
Output: [["baba","abat","baba","atal"],["baba","abat","baba","atan"]]
Explanation:
The output consists of two word squares. The order of output does not matter
(just the order of words in each word square matters).
```

**Constraints:**

- `1 <= words.length <= 1000`
- `1 <= words[i].length <= 4`
- All `words[i]` have the same length.
- `words[i]` consists of only lowercase English letters.
- All `words[i]` are **unique**.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Uber       | ★★☆☆☆ Low        | 2022          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Trie (prefix tree)** — the search repeatedly asks "which words start with this prefix?"; a trie storing word indices at each node answers that in O(prefix length) and is the canonical optimal structure → see [`/dsa/trie.md`](/dsa/trie.md)
- **Backtracking** — we build the square one row at a time, place a candidate word, recurse, and undo on return; the whole solution is a DFS over row choices → see [`/dsa/backtracking.md`](/dsa/backtracking.md)
- **Hash Table (prefix → words index)** — a lighter alternative to the trie: map every prefix to the list of words that start with it for O(1) candidate retrieval → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Matrix / grid symmetry** — the row-k-equals-column-k rule is what forces the first `k` letters of row `k`, which is the constraint that drives the prefix search → see [`/dsa/matrix_traversal.md`](/dsa/matrix_traversal.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Backtracking + Linear Prefix Scan | ~O(Nⁿ·L) worst | O(n) | Simplest to reason about; re-scans all words per row |
| 2 | Backtracking + Prefix Hash Map | Preprocess O(N·L²), then heavily pruned search | O(N·L²) | Easiest fast solution; one map lookup per row |
| 3 | Backtracking + Trie (Optimal) | Build O(N·L), prefix query O(L + matches) | O(N·L) | The textbook answer; least memory of the indexed variants |

`N` = number of words, `L` = word length (= side `n` of the square).

---

## The Core Idea (shared by all approaches)

Build the square **row by row**. In a valid square, row `k` equals column `k`. So once rows `0..k-1` are placed, the **first `k` characters of row `k` are already forced**: they must equal `square[0][k], square[1][k], …, square[k-1][k]` (reading down column `k`). Therefore any word we place at row `k` must **start with that prefix**. The three approaches differ only in *how they find the words matching that prefix*.

---

## Approach 1 — Backtracking + Linear Prefix Scan

### Intuition

Directly implement the core idea with the dumbest lookup: to fill row `k`, compute the column-imposed prefix, then scan the whole word list and take every word that `strings.HasPrefix` the prefix. Place it, recurse to the next row, and pop it when we return. When `n` rows are placed we have a complete square.

### Algorithm

1. `n = ` word length. For the first row (empty prefix) every word qualifies.
2. To fill row `k` (`k ≥ 1`): build `prefix` = column-`k` letters of rows `0..k-1`.
3. Scan all words; each word starting with `prefix` is a candidate — place it and recurse.
4. On return, remove it (backtrack). When the square has `n` rows, snapshot it into the results.

### Complexity

- **Time:** Roughly O(Nⁿ · L) in the worst case — up to `N` candidates per row, `n` rows, and each prefix scan is O(N·L). Fine for `L ≤ 4`.
- **Space:** O(n) recursion depth plus the output list.

### Code

```go
func backtrackingBruteForce(words []string) [][]string {
	var results [][]string
	if len(words) == 0 {
		return results
	}
	n := len(words[0]) // every word has the same length = side of the square
	var square []string

	// buildPrefix reads column k down the rows already placed to get the
	// mandatory prefix of the next word.
	buildPrefix := func(k int) string {
		var sb strings.Builder
		for r := 0; r < len(square); r++ {
			sb.WriteByte(square[r][k]) // the k-th char of row r is forced into column k
		}
		return sb.String()
	}

	var backtrack func()
	backtrack = func() {
		if len(square) == n { // n rows placed → a full, valid square
			row := make([]string, n)
			copy(row, square) // snapshot (square is mutated during search)
			results = append(results, row)
			return
		}
		k := len(square)         // index of the row we are about to fill
		prefix := buildPrefix(k) // letters column k demands at the start of this row
		for _, w := range words {
			// A candidate word must begin with the column-imposed prefix.
			if strings.HasPrefix(w, prefix) {
				square = append(square, w)      // tentatively place it as row k
				backtrack()                     // fill the remaining rows
				square = square[:len(square)-1] // undo and try the next word
			}
		}
	}

	backtrack()
	return results
}
```

### Dry Run

Input `words = ["area","lead","wall","lady","ball"]`, `n = 4`. Tracing the branch that yields `["ball","area","lead","lady"]`.

| row k | prefix (column-k so far) | candidate words with prefix | placed |
|-------|--------------------------|-----------------------------|--------|
| 0 | `""` (empty)            | all → try `ball`            | `ball` |
| 1 | col1 of {ball} = `a`    | `area`                      | `area` |
| 2 | col2 of {ball,area} = `le` | `lead`                   | `lead` |
| 3 | col3 of {ball,area,lead} = `lad` | `lady`             | `lady` |
| — | 4 rows placed → record `["ball","area","lead","lady"]` | | ✓ |

At row 0 the search also tries `wall`, producing the second square `["wall","area","lead","lady"]`; other first-row choices (`area`, `lead`, `lady`) dead-end when no word matches a later prefix.

---

## Approach 2 — Backtracking + Prefix Hash Map

### Intuition

The only slow part above is "find all words starting with `prefix`". Precompute it once: for every word register it under **each** of its prefixes (`""`, `"w"`, `"wa"`, `"wal"`, `"wall"`). Then filling a row is a single map lookup returning exactly the candidates — no scanning.

### Algorithm

1. Build `prefixMap`: `prefix → []word` for every prefix (including `""`) of every word.
2. Backtrack row by row exactly as before, but obtain candidates from `prefixMap[prefix]`.
3. Snapshot squares once `n` rows are placed.

### Complexity

- **Time:** Preprocess O(N·L²) — each word has `L+1` prefixes of length up to `L`. The search itself is strongly pruned (each lookup is exact), making this the fastest in practice for these constraints.
- **Space:** O(N·L²) for the prefix map + O(n) recursion.

### Code

```go
func backtrackingPrefixMap(words []string) [][]string {
	var results [][]string
	if len(words) == 0 {
		return results
	}
	n := len(words[0])

	// prefixMap[p] = all words that start with prefix p (including p == "").
	prefixMap := make(map[string][]string)
	for _, w := range words {
		for i := 0; i <= len(w); i++ {
			p := w[:i]                             // every leading slice, "" .. whole word
			prefixMap[p] = append(prefixMap[p], w) // index the word under this prefix
		}
	}

	var square []string
	buildPrefix := func(k int) string {
		var sb strings.Builder
		for r := 0; r < len(square); r++ {
			sb.WriteByte(square[r][k])
		}
		return sb.String()
	}

	var backtrack func()
	backtrack = func() {
		if len(square) == n {
			row := make([]string, n)
			copy(row, square)
			results = append(results, row)
			return
		}
		k := len(square)
		prefix := buildPrefix(k)
		// Exactly the words that fit column k's demand — no scanning.
		for _, w := range prefixMap[prefix] {
			square = append(square, w)
			backtrack()
			square = square[:len(square)-1]
		}
	}

	backtrack()
	return results
}
```

### Dry Run

Input `words = ["area","lead","wall","lady","ball"]`.

Prefix map (relevant entries):

| prefix | words indexed |
|--------|---------------|
| `""`   | area, lead, wall, lady, ball |
| `"a"`  | area |
| `"le"` | lead |
| `"lad"`| lady |

Backtracking the winning branch:

| row k | prefix | `prefixMap[prefix]` | placed |
|-------|--------|---------------------|--------|
| 0 | `""`   | {area, lead, wall, lady, ball} → try `ball` | `ball` |
| 1 | `"a"`  | {area}              | `area` |
| 2 | `"le"` | {lead}              | `lead` |
| 3 | `"lad"`| {lady}              | `lady` |
| — | record `["ball","area","lead","lady"]` | | ✓ |

The `wall` branch is found the same way; result matches the expected two squares.

---

## Approach 3 — Backtracking + Trie (Optimal)

### Intuition

A trie stores every word as a path; at each node we keep the indices of all words passing through it (i.e. sharing that prefix). Getting the candidates for a row = walk the trie along the required prefix and read that node's index list — O(prefix length) with memory proportional to *shared* prefixes rather than every prefix string. The empty prefix (row 0) is handled by seeding the root with every word index.

### Algorithm

1. Insert all words into a trie; at each visited node append the word's index. Seed the **root** with all indices (every word shares the empty prefix).
2. Backtrack row by row; candidates = `wordsWithPrefix(columnPrefix)` (walk + read the node's list).
3. Snapshot squares when `n` rows are placed.

### Complexity

- **Time:** Build O(N·L); each prefix query is O(L + matches). Overall dominated by the pruned search — asymptotically the best of the three.
- **Space:** O(N·L) trie nodes plus the per-node index lists + O(n) recursion.

### Code

```go
func backtrackingTrie(words []string) [][]string {
	var results [][]string
	if len(words) == 0 {
		return results
	}
	n := len(words[0])

	root := newSquareTrieNode()
	for idx, w := range words {
		node := root
		root.wordIdx = append(root.wordIdx, idx) // every word shares the empty prefix (row 0 candidates)
		for i := 0; i < len(w); i++ {
			c := w[i] - 'a'
			if node.children[c] == nil {
				node.children[c] = newSquareTrieNode()
			}
			node = node.children[c]
			node.wordIdx = append(node.wordIdx, idx) // this word passes through here
		}
	}

	// wordsWithPrefix walks the prefix and returns indices of all words under it.
	wordsWithPrefix := func(prefix string) []int {
		node := root
		for i := 0; i < len(prefix); i++ {
			c := prefix[i] - 'a'
			if node.children[c] == nil {
				return nil // no word has this prefix → dead end
			}
			node = node.children[c]
		}
		return node.wordIdx
	}

	var square []string
	buildPrefix := func(k int) string {
		var sb strings.Builder
		for r := 0; r < len(square); r++ {
			sb.WriteByte(square[r][k])
		}
		return sb.String()
	}

	var backtrack func()
	backtrack = func() {
		if len(square) == n {
			row := make([]string, n)
			copy(row, square)
			results = append(results, row)
			return
		}
		k := len(square)
		prefix := buildPrefix(k)
		for _, idx := range wordsWithPrefix(prefix) {
			square = append(square, words[idx]) // place candidate row
			backtrack()
			square = square[:len(square)-1] // backtrack
		}
	}

	backtrack()
	return results
}

type squareTrieNode struct {
	children [26]*squareTrieNode
	wordIdx  []int // indices into the original words slice
}
```

### Dry Run

Input `words = ["area","lead","wall","lady","ball"]` (indices 0..4). Trie (indices shown at each node along a path):

```
root [0,1,2,3,4]
 ├─ a → r → e → a        (word 0 "area")
 ├─ l → e → a → d        (word 1 "lead")
 │    └ a → d → y        (word 3 "lady")
 ├─ w → a → l → l        (word 2 "wall")
 └─ b → a → l → l        (word 4 "ball")
```

Backtracking the winning branch, reading each node's index list:

| row k | prefix | walk result → candidate indices | placed |
|-------|--------|--------------------------------|--------|
| 0 | `""`   | root → [0,1,2,3,4] → try idx 4 | `ball` |
| 1 | `"a"`  | a-node → [0] (only "area")     | `area` |
| 2 | `"le"` | l→e node → [1] ("lead")        | `lead` |
| 3 | `"lad"`| l→a→d node → [3] ("lady")      | `lady` |
| — | record `["ball","area","lead","lady"]` | | ✓ |

Trying idx 2 (`wall`) at row 0 yields the second square. Total: the expected two squares.

---

## Key Takeaways

- **Row k = column k forces a prefix.** The instant you place rows `0..k-1`, the start of row `k` is fixed. Every word-square constructor is "backtrack over rows, pruned by a prefix query".
- **Prefix queries want a trie (or a prefix hash map).** Registering word indices at each trie node turns "all words with this prefix" into a single walk; the hash-map version trades O(N·L²) memory for even simpler code.
- **Seed the root for the empty prefix.** A subtle bug: row 0 has an empty prefix, so the root node must carry *all* word indices — otherwise the search returns nothing.
- The three approaches are the same DFS; only the candidate-lookup engine changes (scan → map → trie), a classic "index your data to prune the search" progression.

---

## Related Problems

- LeetCode #422 — Valid Word Square (verify one square — the checking counterpart)
- LeetCode #212 — Word Search II (trie + backtracking on a grid)
- LeetCode #208 — Implement Trie (Prefix Tree) (the underlying structure)
- LeetCode #211 — Design Add and Search Words Data Structure (trie with wildcard search)
- LeetCode #79 — Word Search (backtracking template)
