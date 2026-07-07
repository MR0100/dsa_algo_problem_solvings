# 0212 — Word Search II

> LeetCode #212 · Difficulty: Hard
> **Categories:** Trie, Depth-First Search, Backtracking, Matrix, String

---

## Problem Statement

Given an `m x n` `board` of characters and a list of strings `words`, return *all words on the board*.

Each word must be constructed from letters of sequentially adjacent cells, where **adjacent cells** are horizontally or vertically neighboring. The same letter cell may not be used more than once in a word.

**Example 1:**
```
Input: board = [["o","a","a","n"],["e","t","a","e"],["i","h","k","r"],["i","f","l","v"]],
       words = ["oath","pea","eat","rain"]
Output: ["eat","oath"]
```

**Example 2:**
```
Input: board = [["a","b"],["c","d"]], words = ["abcb"]
Output: []
```

**Constraints:**
- `m == board.length`
- `n == board[i].length`
- `1 <= m, n <= 12`
- `board[i][j]` is a lowercase English letter.
- `1 <= words.length <= 3 * 10⁴`
- `1 <= words[i].length <= 10`
- `words[i]` consists of lowercase English letters.
- All the strings of `words` are unique.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★★ Very High  | 2024          |
| Microsoft  | ★★★★☆ High       | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Meta       | ★★★☆☆ Medium     | 2024          |
| Uber       | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Trie / Prefix Tree** — all words share one trie so a single board DFS searches every word at once, pruning directions whose prefix no word extends → see [`/dsa/trie.md`](/dsa/trie.md)
- **Depth-First Search + Backtracking** — grid traversal marking cells visited on the path and restoring them on return → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Matrix Traversal** — 4-directional movement over a 2D grid → see [`/dsa/matrix_traversal.md`](/dsa/matrix_traversal.md)
- **Backtracking** — the classic in-place visited-mark / undo pattern → see [`/dsa/backtracking.md`](/dsa/backtracking.md)

---

## Approaches Overview

Let R×C = board size, W = number of words, L = max word length.

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (DFS per word) | O(W · R·C · 4ᴸ) | O(L) | Few words; conceptual baseline |
| 2 | Trie + Board DFS (Optimal) | O(R·C · 4·3ᴸ⁻¹) | O(Σ word chars) | Many words with shared prefixes — the intended answer |

---

## Approach 1 — Brute Force (DFS per Word)

### Intuition
Word Search I answers "is *this one* word on the board?" via a DFS from every cell that marks cells visited and backtracks. Word Search II is that, once per word: keep the words whose independent search succeeds. Correct and simple, but each word re-explores the board from scratch, so words sharing a prefix (`oath`/`oat`/`oa`) re-walk the same cells repeatedly — the redundancy Approach 2 eliminates.

### Algorithm
1. For each word `w`, call `existsOnBoard(board, w)`.
2. `existsOnBoard` tries to match `w` starting from every cell using a DFS: match `w[i]` at `(r,c)`; on success mark the cell `'#'`, recurse into the 4 neighbours for `w[i+1]`, then restore the cell.
3. Collect the words that match somewhere.

### Complexity
- **Time:** O(W · R·C · 4ᴸ) — W independent searches, each starting at R·C cells and branching up to 4 ways for L steps.
- **Space:** O(L) recursion depth; visited state handled in-place on the board.

### Code
```go
func bruteForce(board [][]byte, words []string) []string {
	var res []string
	for _, w := range words {
		if existsOnBoard(board, w) {
			res = append(res, w)
		}
	}
	return res
}

func existsOnBoard(board [][]byte, word string) bool {
	rows, cols := len(board), len(board[0])
	var dfs func(r, c, i int) bool
	dfs = func(r, c, i int) bool {
		if i == len(word) {
			return true
		}
		if r < 0 || r >= rows || c < 0 || c >= cols || board[r][c] != word[i] {
			return false
		}
		saved := board[r][c]
		board[r][c] = '#' // mark visited
		found := dfs(r+1, c, i+1) || dfs(r-1, c, i+1) ||
			dfs(r, c+1, i+1) || dfs(r, c-1, i+1)
		board[r][c] = saved // backtrack
		return found
	}
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if dfs(r, c, 0) {
				return true
			}
		}
	}
	return false
}
```

### Dry Run (Example 1)

Board:
```
o a a n
e t a e
i h k r
i f l v
```

Cells labelled `(row,col)`. Board recap: row0 `o a a n`, row1 `e t a e`, row2 `i h k r`, row3 `i f l v`.

| Word | Traced path (4-adjacent, no cell reused) | Kept? |
|------|-------------------------------------------|-------|
| `oath` | `o(0,0) → a(0,1) → t(1,1) → h(2,1)` — each step is up/down/left/right adjacent, spells `oath` | yes |
| `pea` | no `p` cell exists on the board → the DFS never starts | no |
| `eat` | `e(1,3) → a(1,2) → t(1,1)` — `(1,3),(1,2),(1,1)` are consecutive cells in row 1, spells `eat` | yes |
| `rain` | `r` occurs only at `(2,3)`; its neighbours are `e(1,3)`, `k(2,2)`, `v(3,3)` — none is `a`, so no `rain` path | no |

Result (sorted): `[eat oath]` ✓

---

## Approach 2 — Trie + Board DFS (Optimal)

### Intuition
Searching words one at a time re-walks shared prefixes. Instead, insert **all** words into a single trie; now one DFS from each cell descends the board and the trie **together**. It only continues in a direction if the trie has a child for that letter, so every word sharing a prefix is searched simultaneously and dead directions are pruned immediately. Store the full word on its terminal trie node so collecting a hit needs no string rebuilding, and clear that marker after collecting so each word is reported once.

### Algorithm
1. Insert every word into a `[26]`-child trie; set `node.word = word` on each terminal node (this doubles as the "complete word" flag).
2. DFS from every cell `(r,c)` carrying the current trie node:
   - Let `ch = board[r][c]`; if `node.children[ch-'a']` is nil, stop (prune).
   - Descend to that child `next`; if `next.word != ""`, append it to results and clear it.
   - Mark the cell `'#'`, recurse into the 4 neighbours against `next`, then restore the cell.
3. Return the collected words.

### Complexity
- **Time:** O(R·C · 4·3ᴸ⁻¹) — one shared traversal (first step 4 directions, then ≤ 3 since we don't revisit the previous cell), pruned by the trie to only letters that exist.
- **Space:** O(total characters across all words) for the trie + O(L) recursion depth.

### Code
```go
type trieNode struct {
	children [26]*trieNode
	word     string // non-empty ⇒ a complete word ends here
}

func insert(root *trieNode, word string) {
	cur := root
	for i := 0; i < len(word); i++ {
		idx := word[i] - 'a'
		if cur.children[idx] == nil {
			cur.children[idx] = &trieNode{}
		}
		cur = cur.children[idx]
	}
	cur.word = word
}

func trieDFS(board [][]byte, words []string) []string {
	root := &trieNode{}
	for _, w := range words {
		insert(root, w)
	}
	rows, cols := len(board), len(board[0])
	var res []string

	var dfs func(r, c int, node *trieNode)
	dfs = func(r, c int, node *trieNode) {
		if r < 0 || r >= rows || c < 0 || c >= cols || board[r][c] == '#' {
			return
		}
		ch := board[r][c]
		next := node.children[ch-'a']
		if next == nil {
			return // prune: no word continues with this letter
		}
		if next.word != "" {
			res = append(res, next.word)
			next.word = "" // report each word once
		}
		saved := board[r][c]
		board[r][c] = '#'
		dfs(r+1, c, next)
		dfs(r-1, c, next)
		dfs(r, c+1, next)
		dfs(r, c-1, next)
		board[r][c] = saved
	}

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			dfs(r, c, root)
		}
	}
	return res
}
```

### Dry Run (Example 1)

Trie of `["oath","pea","eat","rain"]`: root children `o,p,e,r`.

| Cell start | Trie step | Outcome |
|------------|-----------|---------|
| (0,0)=`o` | root→`o` exists | descend; `o→a→t→h` follows board `o(0,0),a(0,1),t(1,1),h(2,1)` → node `oath`.word set → **collect "oath"**, clear it |
| (1,0)=`e` | root→`e` exists | `e→a→t` follows `e(1,3)`? explores `e` cells; a valid `e→a→t` path exists on the board → **collect "eat"**, clear it |
| any `p` cell | root→`p` | no `p` on board → branch never entered → `pea` never collected |
| (2,3)=`r` | root→`r` | `r→a→i→n`? board has no adjacent `a` after this `r` completing `rain` → not collected |

Collected: `{oath, eat}`; sorted → `[eat oath]` ✓

---

## Key Takeaways

- **One trie beats W separate searches.** Sharing prefixes across all words collapses W board explorations into one, and the trie prunes any direction whose letter no word needs — the whole point of the problem.
- **Store the word on the terminal node.** Tagging `node.word` avoids rebuilding the path string on every hit and gives a free "is this a word?" flag.
- **Clear the marker after collecting** (`next.word = ""`) to report each word exactly once and to prune already-found leaves.
- **Backtracking on the board = mark `'#'` before recursing, restore after.** In-place visited marking avoids an extra visited matrix.
- Optional extra pruning (not shown): delete exhausted leaf children after recursion so the trie shrinks as words are found — helps on adversarial inputs.
- This is the direct sequel to #208 (build the trie) and #211 (DFS over a trie); it also generalises #79 Word Search from one word to many.

---

## Related Problems

- LeetCode #79 — Word Search (single word, the brute-force core)
- LeetCode #208 — Implement Trie (Prefix Tree) (the trie used here)
- LeetCode #211 — Design Add and Search Words Data Structure (DFS over a trie)
- LeetCode #140 — Word Break II (trie/DFS word segmentation)
- LeetCode #472 — Concatenated Words (trie + DFS over a dictionary)
