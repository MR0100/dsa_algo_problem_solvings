package main

import (
	"fmt"
	"sort"
	"strings"
)

// ── Approach 1: Backtracking with Linear Prefix Scan (Brute Force) ────────────
//
// backtrackingBruteForce solves Word Squares by building the square row by row,
// and for each next row scanning ALL words to find those whose prefix matches
// the column constraint imposed by the rows already placed.
//
// Intuition:
//
//	In a valid square the k-th row equals the k-th column. So once rows 0..k-1
//	are fixed, the first k letters of row k are forced: they must equal
//	square[0][k], square[1][k], …, square[k-1][k] (reading down column k). Any
//	word we place at row k must therefore START WITH that prefix. Try every
//	word that qualifies, recurse to the next row, and when we have placed n rows
//	we have a complete square. The brute force finds qualifying words by simply
//	scanning the whole list and testing strings.HasPrefix.
//
// Algorithm:
//  1. n = word length. For the first row, try every word.
//  2. To fill row k (k ≥ 1): build prefix = column-k letters of rows 0..k-1.
//     Scan all words; each word with that prefix is a candidate.
//  3. Place a candidate, recurse to row k+1; on return, pop it (backtrack).
//  4. When the square has n rows, copy it into the results.
//
// Time:  O(n · N · L) per row of scanning × branching → roughly O(N^n · L) worst
//
//	case, where N = number of words, L = word length, n = L. Exponential; fine
//	for the tiny constraints (L ≤ 4) but slower than the indexed version.
//
// Space: O(n) recursion depth plus the output.
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

// ── Approach 2: Backtracking with Prefix Hash Map ────────────────────────────
//
// backtrackingPrefixMap solves Word Squares the same way but replaces the
// O(N) prefix scan with an O(1) hash-map lookup: precompute, for every prefix,
// the list of words that start with it.
//
// Intuition:
//
//	The only expensive step in the brute force is "find all words starting with
//	this prefix". Precompute it: for each word, register it under every one of
//	its prefixes ("", "w", "wa", "wal", "wall"). Then filling a row is a single
//	map lookup returning exactly the candidate words — no scanning.
//
// Algorithm:
//  1. Build prefixMap: prefix → []word for every prefix of every word.
//  2. Backtrack row by row as before, but get candidates via prefixMap[prefix].
//  3. Collect squares once n rows are placed.
//
// Time:  Preprocess O(N·L²) (each word has L prefixes of length ≤ L). Search
//
//	explores far fewer branches because each lookup is exact. Practically the
//	fastest for these constraints.
//
// Space: O(N·L²) for the prefix map + O(n) recursion.
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

// ── Approach 3: Backtracking with a Trie (Optimal) ───────────────────────────
//
// backtrackingTrie solves Word Squares using a trie so that "all words with
// this prefix" is answered by walking the prefix once and reading the word
// indices stored at that node — no per-prefix string keys, less memory churn.
//
// Intuition:
//
//	A trie stores every word as a path; at each node we keep the indices of all
//	words passing through it (i.e. sharing that prefix). To get the candidates
//	for a row, walk the trie along the required prefix and read that node's
//	index list. This is the textbook optimal structure: prefix queries in
//	O(prefix length) and memory proportional to shared prefixes.
//
// Algorithm:
//  1. Insert all words into a trie; at each visited node append the word index.
//  2. Backtrack row by row; candidates = wordsWithPrefix(trie, columnPrefix).
//  3. Collect squares when n rows are placed.
//
// Time:  Build O(N·L). Prefix query O(L + matches). Overall dominated by the
//
//	(pruned) search; asymptotically the best of the three.
//
// Space: O(N·L) trie nodes (plus index lists) + O(n) recursion.
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

// squareTrieNode is a lowercase-alphabet trie node that also stores the indices
// of every word sharing the prefix ending at this node.
type squareTrieNode struct {
	children [26]*squareTrieNode
	wordIdx  []int // indices into the original words slice
}

func newSquareTrieNode() *squareTrieNode { return &squareTrieNode{} }

// sortSquares canonicalises the output (sort each square is unnecessary since
// rows are ordered, but sort the LIST of squares) so the printed result is
// deterministic and easy to compare against the expected answer.
func sortSquares(sq [][]string) [][]string {
	sort.Slice(sq, func(i, j int) bool {
		return strings.Join(sq[i], ",") < strings.Join(sq[j], ",")
	})
	return sq
}

func main() {
	ex1 := []string{"area", "lead", "wall", "lady", "ball"}
	ex2 := []string{"abat", "baba", "atan", "atal"}

	fmt.Println("=== Approach 1: Backtracking with Linear Prefix Scan ===")
	fmt.Println(sortSquares(backtrackingBruteForce(ex1))) // expected [[ball area lead lady] [wall area lead lady]]
	fmt.Println(sortSquares(backtrackingBruteForce(ex2))) // expected [[baba abat baba atal] [baba abat baba atan]]

	fmt.Println("=== Approach 2: Backtracking with Prefix Hash Map ===")
	fmt.Println(sortSquares(backtrackingPrefixMap(ex1))) // expected [[ball area lead lady] [wall area lead lady]]
	fmt.Println(sortSquares(backtrackingPrefixMap(ex2))) // expected [[baba abat baba atal] [baba abat baba atan]]

	fmt.Println("=== Approach 3: Backtracking with a Trie (Optimal) ===")
	fmt.Println(sortSquares(backtrackingTrie(ex1))) // expected [[ball area lead lady] [wall area lead lady]]
	fmt.Println(sortSquares(backtrackingTrie(ex2))) // expected [[baba abat baba atal] [baba abat baba atan]]
}
