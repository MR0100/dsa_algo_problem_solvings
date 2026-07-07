package main

import "fmt"

// ── Approach 1: Backtracking (DFS) ───────────────────────────────────────────
//
// exist solves Word Search using DFS backtracking on the board.
//
// Intuition:
//   For each cell, try to start a DFS that matches word[0..len-1].
//   At each step, mark the current cell as visited (e.g., by XOR-ing the
//   character with a sentinel) to prevent reuse, then try all 4 directions.
//   Unmark on backtrack.
//
// Algorithm:
//   for each cell (r,c):
//     if dfs(r, c, 0): return true
//
//   dfs(r, c, idx):
//     if idx == len(word): return true (all chars matched)
//     if out of bounds or board[r][c] != word[idx]: return false
//     mark board[r][c] visited
//     for each of 4 directions: if dfs(nr,nc,idx+1): return true
//     unmark (restore)
//     return false
//
// Time:  O(m × n × 4^L) — m×n starting cells, 4^L paths of length L.
//         Pruning makes it much faster in practice.
// Space: O(L) — recursion stack depth L (word length).
func exist(board [][]byte, word string) bool {
	m, n := len(board), len(board[0])
	dr := []int{0, 0, 1, -1}
	dc := []int{1, -1, 0, 0}

	var dfs func(r, c, idx int) bool
	dfs = func(r, c, idx int) bool {
		if idx == len(word) {
			return true // all characters matched
		}
		if r < 0 || r >= m || c < 0 || c >= n || board[r][c] != word[idx] {
			return false
		}
		// mark visited by XOR-ing with a non-letter byte
		board[r][c] ^= 255
		for d := 0; d < 4; d++ {
			if dfs(r+dr[d], c+dc[d], idx+1) {
				board[r][c] ^= 255 // restore before returning
				return true
			}
		}
		board[r][c] ^= 255 // restore (backtrack)
		return false
	}

	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if dfs(r, c, 0) {
				return true
			}
		}
	}
	return false
}

func main() {
	fmt.Println("=== Word Search (Backtracking DFS) ===")

	b1 := [][]byte{
		{'A', 'B', 'C', 'E'},
		{'S', 'F', 'C', 'S'},
		{'A', 'D', 'E', 'E'},
	}
	fmt.Printf("word=%q  got=%v  expected true\n", "ABCCED", exist(b1, "ABCCED"))

	b2 := [][]byte{
		{'A', 'B', 'C', 'E'},
		{'S', 'F', 'C', 'S'},
		{'A', 'D', 'E', 'E'},
	}
	fmt.Printf("word=%q  got=%v  expected true\n", "SEE", exist(b2, "SEE"))

	b3 := [][]byte{
		{'A', 'B', 'C', 'E'},
		{'S', 'F', 'C', 'S'},
		{'A', 'D', 'E', 'E'},
	}
	fmt.Printf("word=%q  got=%v  expected false\n", "ABCB", exist(b3, "ABCB"))

	b4 := [][]byte{{'a'}}
	fmt.Printf("word=%q  got=%v  expected true\n", "a", exist(b4, "a"))
}
