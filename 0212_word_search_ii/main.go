package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Brute Force (DFS per Word) ───────────────────────────────────
//
// bruteForce solves Word Search II by running an independent board DFS (the
// Word Search I algorithm) for every word in the list.
//
// Intuition:
//
//	Word Search I answers "is this one word on the board?" via DFS from every
//	cell. Word Search II is just that, once per word: keep the words whose
//	search succeeds. Correct and simple, but each word re-explores the board
//	from scratch, so shared prefixes ("oath"/"oat"/"oa") are re-walked over
//	and over — exactly the redundancy Approach 2 removes with a trie.
//
// Algorithm:
//
//	For each word: try to match it starting from each cell with a DFS that
//	marks visited cells (temporarily), backtracking on return. Collect words
//	that match anywhere.
//
// Time:  O(W · R·C · 4^L) — W words, R·C start cells, L = word length, up to 4
//
//	directions per step.
//
// Space: O(L) recursion depth (visited handled in-place).
func bruteForce(board [][]byte, words []string) []string {
	var res []string
	for _, w := range words { // independently search each word
		if existsOnBoard(board, w) {
			res = append(res, w)
		}
	}
	return res
}

// existsOnBoard reports whether word can be traced on the board (Word Search I).
func existsOnBoard(board [][]byte, word string) bool {
	rows, cols := len(board), len(board[0])
	var dfs func(r, c, i int) bool
	dfs = func(r, c, i int) bool {
		if i == len(word) {
			return true // matched every character
		}
		if r < 0 || r >= rows || c < 0 || c >= cols || board[r][c] != word[i] {
			return false // off-board or wrong letter
		}
		saved := board[r][c] // remember the letter
		board[r][c] = '#'    // mark visited so we don't reuse this cell
		found := dfs(r+1, c, i+1) || dfs(r-1, c, i+1) ||
			dfs(r, c+1, i+1) || dfs(r, c-1, i+1) // explore 4 neighbours
		board[r][c] = saved // backtrack: restore the cell
		return found
	}
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if dfs(r, c, 0) { // try starting the word at every cell
				return true
			}
		}
	}
	return false
}

// ── Approach 2: Trie + Board DFS (Optimal) ───────────────────────────────────
//
// trieDFS solves Word Search II by inserting all words into a trie, then doing
// ONE board DFS that walks the trie in lockstep, collecting a word whenever the
// current trie node marks a complete word.
//
// Intuition:
//
//	Searching words one at a time re-walks shared prefixes. Instead, build a
//	trie of all words: now a single DFS from each cell descends the board and
//	the trie together — it only continues down a direction if the trie has a
//	child for that letter, so all words sharing a prefix are searched at once.
//	Storing the finished word ON the trie node avoids rebuilding strings, and
//	de-duplicating found words plus pruning dead trie branches keep it fast.
//
// Algorithm:
//
//  1. Insert every word into a [26]-child trie; store the whole word on its
//     terminal node (node.word) as the "this is a complete word" marker.
//  2. DFS from every board cell, carrying the current trie node:
//     - letter = board[r][c]; if the node has no child for it, stop.
//     - descend to that child; if child.word != "" collect it and clear the
//     marker (so each word is reported once).
//     - mark the cell visited, recurse into 4 neighbours, then backtrack.
//  3. (Pruning) after exploring, drop leaf trie children that are exhausted.
//
// Time:  O(R·C·4·3^(L-1)) where L = max word length — one shared traversal
//
//	instead of one per word; the trie caps branching to letters that exist.
//
// Space: O(total characters in words) for the trie + O(L) recursion depth.
type trieNode struct {
	children [26]*trieNode // children[i] = subtree for letter 'a'+i
	word     string        // non-empty ⇒ a complete word ends here (stores the word)
}

// insert adds word to the trie rooted at root, tagging its terminal node.
func insert(root *trieNode, word string) {
	cur := root
	for i := 0; i < len(word); i++ {
		idx := word[i] - 'a'
		if cur.children[idx] == nil {
			cur.children[idx] = &trieNode{} // create missing edge
		}
		cur = cur.children[idx]
	}
	cur.word = word // remember the full word at its end node
}

// trieDFS returns every word that can be traced on the board.
func trieDFS(board [][]byte, words []string) []string {
	root := &trieNode{}
	for _, w := range words {
		insert(root, w) // build the combined trie once
	}
	rows, cols := len(board), len(board[0])
	var res []string

	var dfs func(r, c int, node *trieNode)
	dfs = func(r, c int, node *trieNode) {
		if r < 0 || r >= rows || c < 0 || c >= cols || board[r][c] == '#' {
			return // off-board or already on the current path
		}
		ch := board[r][c]
		next := node.children[ch-'a']
		if next == nil {
			return // no word continues with this letter → prune the branch
		}
		if next.word != "" {
			res = append(res, next.word) // a full word ends here: collect it
			next.word = ""               // clear marker → report each word once
		}
		saved := board[r][c]
		board[r][c] = '#' // mark visited
		dfs(r+1, c, next) // explore 4 neighbours, all against `next`
		dfs(r-1, c, next)
		dfs(r, c+1, next)
		dfs(r, c-1, next)
		board[r][c] = saved // backtrack
	}

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			dfs(r, c, root) // start every DFS at the trie root
		}
	}
	return res
}

// toGrid converts a slice of strings into a byte board for convenience.
func toGrid(rows []string) [][]byte {
	g := make([][]byte, len(rows))
	for i, row := range rows {
		g[i] = []byte(row)
	}
	return g
}

// sortedCopy returns res sorted so output order is deterministic for comparison.
func sortedCopy(res []string) []string {
	out := append([]string(nil), res...)
	sort.Strings(out)
	return out
}

func main() {
	// Example 1
	board1 := []string{"oaan", "etae", "ihkr", "iflv"}
	words1 := []string{"oath", "pea", "eat", "rain"}

	fmt.Println("=== Approach 1: Brute Force (DFS per Word) ===")
	fmt.Println(sortedCopy(bruteForce(toGrid(board1), words1))) // [eat oath]
	fmt.Println("=== Approach 2: Trie + Board DFS (Optimal) ===")
	fmt.Println(sortedCopy(trieDFS(toGrid(board1), words1))) // [eat oath]

	// Example 2
	board2 := []string{"ab", "cd"}
	words2 := []string{"abcb"}

	fmt.Println("=== Approach 1: Brute Force (DFS per Word) ===")
	fmt.Println(sortedCopy(bruteForce(toGrid(board2), words2))) // []
	fmt.Println("=== Approach 2: Trie + Board DFS (Optimal) ===")
	fmt.Println(sortedCopy(trieDFS(toGrid(board2), words2))) // []
}
