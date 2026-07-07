# 0472 — Concatenated Words

> LeetCode #472 · Difficulty: Hard
> **Categories:** Array, String, Dynamic Programming, Depth-First Search, Trie

---

## Problem Statement

Given an array of strings `words` (**without duplicates**), return *all the **concatenated words** in the given list of `words`*.

A **concatenated word** is defined as a string that is comprised entirely of at least two shorter words (not necessarily distinct) in the given array.

**Example 1:**

```
Input: words = ["cat","cats","catsdogcats","dog","dogcatsdog","hippopotamuses","rat","ratcatdogcat"]
Output: ["catsdogcats","dogcatsdog","ratcatdogcat"]
Explanation: "catsdogcats" can be concatenated by "cats", "dog" and "cats";
"dogcatsdog" can be concatenated by "dog", "cats" and "dog";
"ratcatdogcat" can be concatenated by "rat", "cat", "dog" and "cat".
```

**Example 2:**

```
Input: words = ["cat","dog","catdog"]
Output: ["catdog"]
```

**Constraints:**

- `1 <= words.length <= 10^4`
- `1 <= words[i].length <= 30`
- `words[i]` consists of only lowercase English letters.
- All the strings of `words` are **unique**.
- `1 <= sum(words[i].length) <= 10^5`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Word Break (segmentation DP)** — each candidate word is a Word Break instance against the dictionary of all words; the twist is requiring at least two pieces → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)
- **Depth-First Search with memoisation** — segment the suffix recursively, caching per start index to kill the exponential blowup → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Trie (prefix tree)** — words share prefixes, so a trie enumerates all dictionary-word prefixes of a candidate in one character walk instead of hashing every substring → see [`/dsa/trie.md`](/dsa/trie.md)
- **Hash Set membership** — the O(1) "is this substring a word?" test underlying the simpler approaches → see [`/dsa/hash_map.md`](/dsa/hash_map.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force Word Break (no memo) | O(N·2^L) | O(N·L) | Explains the "≥ 2 pieces" rule; blows up on long words |
| 2 | Word Break DP (memoised DFS) | O(N·L²) | O(N·L) | The clean, accepted answer; easy to get right |
| 3 | Trie + DFS | O(N·L²) | O(Σ chars) | Avoids substring slicing; fastest constants with shared prefixes |

*(N = number of words, L = max word length ≤ 30.)*

---

## Approach 1 — Brute Force Word Break (Set + DFS, no memo)

### Intuition

A word counts if it splits into **≥ 2** dictionary words. Drop every word into a hash set, then test a word by trying each prefix that is itself a word and recursing on the remaining suffix. The only subtlety is the "at least two" rule: the *first* cut must be a **proper** prefix (strictly shorter than the whole word), otherwise a word would trivially "match itself" as a single piece. Once one real cut is made, the leftover suffix guarantees two or more pieces, so any dictionary word — even the whole original string — is allowed thereafter.

### Algorithm

1. Insert all words into a set `dict`.
2. For each word `w`, run `canBreak(w, 0)`.
3. `canBreak(s, i)`: if `i == len(s)` return `true`. For `end` in `i+1..len(s)`: skip the cut where `i == 0 && end == len(s)` (the whole word); if `s[i:end] ∈ dict` and `canBreak(s, end)` → `true`.
4. Collect every `w` for which it returns `true`.

### Complexity

- **Time:** O(N · 2^L) — without memoisation each word-break can branch on every prefix boundary, up to `2^L` combinations per word.
- **Space:** O(N·L) for the set + O(L) recursion depth.

### Code

```go
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
```

### Dry Run

Example 2: `words = ["cat","dog","catdog"]`, `dict = {cat, dog, catdog}`. Test `w = "catdog"` (`len = 6`), call `canBreak("catdog", 0)`:

| i | end | substring `s[i:end]` | in dict? | whole-word cut? | recurse | result |
|---|-----|----------------------|----------|-----------------|---------|--------|
| 0 | 1 | "c" | no | — | — | — |
| 0 | 2 | "ca" | no | — | — | — |
| 0 | 3 | "cat" | yes | no | `canBreak(3)` | see below |
| 3 | 4 | "d" | no | — | — | — |
| 3 | 5 | "do" | no | — | — | — |
| 3 | 6 | "dog" | yes | no | `canBreak(6)` → `i==len` → true | **true** |

`canBreak("catdog",0)` returns `true`. Testing `"cat"` / `"dog"`: their only whole-length cut is forbidden and no proper prefix is a word → excluded.

Result: `["catdog"]` ✔

---

## Approach 2 — Word Break DP (Memoised DFS)

### Intuition

Whether the suffix `s[i:]` can be segmented depends only on `i`, not on the path taken to reach it. So cache the verdict per start index and each suffix is solved once — the exponential search collapses to polynomial. The "≥ 2 words" rule is preserved exactly as before: forbid the single cut that equals the whole word at `i == 0`.

### Algorithm

1. Build set `dict`.
2. For each word `w`, keep `memo[i] ∈ {unknown, true, false}`.
3. `seg(i)`: if `i == len(w)` return `true`; if cached, return it. For `end` in `i+1..len(w)`: skip whole-word cut at `i==0`; if `w[i:end] ∈ dict` and `seg(end)` → cache `true`. Else cache `false`.
4. `w` is concatenated iff `seg(0)`.

### Complexity

- **Time:** O(N · L²) — per word there are `L` start indices, each scanning up to `L` end boundaries; the memo ensures each suffix is expanded once.
- **Space:** O(N·L) for the set + O(L) for the per-word memo and recursion.

### Code

```go
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
```

### Dry Run

Example 2: `w = "catdog"`, `dict = {cat, dog, catdog}`. `memo` starts empty.

| Call seg(i) | scan | cached memo[i] | returns |
|-------------|------|----------------|---------|
| seg(0) | "cat" (i=0..3) ∈ dict, whole-word "catdog" skipped | — | delegates to seg(3) |
| ↳ seg(3) | "dog" (3..6) ∈ dict → seg(6) | `memo[3]=1` | true |
| ↳↳ seg(6) | `6 == len` | — | true |
| seg(0) resumes | "cat" cut succeeded | `memo[0]=1` | **true** |

`seg(0) = true` → `"catdog"` included. Result: `["catdog"]` ✔

---

## Approach 3 — Trie + DFS (Optimal)

### Intuition

The building-block words overlap heavily on prefixes (`cat`, `cats`, `catdog`…). A trie lets us find *every* dictionary-word prefix of the current word by walking its characters once and following child links — no substring hashing. Each trie node flagged `isEnd` marks a legal cut; recurse on the remaining suffix. As always, forbid the whole word as a single first piece to enforce "≥ 2". Memoising start positions keeps it polynomial.

### Algorithm

1. Insert every word into a trie (set `isEnd` at terminals).
2. For each word `w`, keep `memo[start]`. `dfs(start, isFirst)`:
   - If `start == len(w)` return `true`; if cached, return it.
   - Walk from the root over `w[start..]`; when the current node `isEnd` at index `p`, let `end = p+1`. Skip if `isFirst && end == len(w)`. If `dfs(end, false)` → success.
   - Stop the walk when a child link is missing.
3. Collect `w` when `dfs(0, true)` succeeds.

### Complexity

- **Time:** O(N · L²) — per word, `L` start positions each drive an `O(L)` trie walk; memoised starts prevent recomputation. Building the trie is O(Σ chars).
- **Space:** O(Σ chars) for the trie + O(L) recursion.

### Code

```go
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

type trieNode struct {
	children [26]*trieNode
	isEnd    bool // true if some inserted word ends exactly here
}

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
```

### Dry Run

Example 2: trie holds `cat`, `dog`, `catdog`. Test `w = "catdog"`, `dfs(0, true)`:

| p | char | node after step | isEnd? | end | first-whole skip? | action |
|---|------|-----------------|--------|-----|-------------------|--------|
| 0 | c | c | no | — | — | keep walking |
| 1 | a | c→a | no | — | — | keep walking |
| 2 | t | c→a→t | yes ("cat") | 3 | no (end=3 ≠ 6) | `dfs(3,false)` |
| — | — | — | — | — | — | ↓ |
| 3 | d | d | no | — | — | keep walking |
| 4 | o | d→o | no | — | — | keep walking |
| 5 | g | d→o→g | yes ("dog") | 6 | not first | `dfs(6,false)` → `start==len` → true |

`dfs(3,false)` returns true → `dfs(0,true)` returns **true**. (Continuing the p-walk of `dfs(0)` would also reach `isEnd` at "catdog", but that whole-word cut is skipped because `isFirst` is set.)

Result: `["catdog"]` ✔

---

## Key Takeaways

- **"Made of ≥ 2 dictionary words" = Word Break + a floor of two pieces.** Enforce the floor by forbidding the single cut that swallows the whole string on the first segment; everything after that automatically has ≥ 2 pieces.
- **Memoise on the suffix start index.** Segmentability of `s[i:]` is path-independent — the single most important optimisation, turning `2^L` into `L²`.
- **A trie replaces substring hashing.** Following child links surfaces every dictionary prefix in one walk, which is why it wins on constants when the dictionary shares prefixes (this problem's whole flavour).
- Sorting words by length first (optional) lets you build the dictionary incrementally so a word is only ever tested against strictly shorter words — a common competitive tweak.

---

## Related Problems

- LeetCode #139 — Word Break (single word, "can it be segmented?")
- LeetCode #140 — Word Break II (enumerate all segmentations)
- LeetCode #208 — Implement Trie (the prefix-tree primitive used here)
- LeetCode #212 — Word Search II (trie + DFS over a grid)
