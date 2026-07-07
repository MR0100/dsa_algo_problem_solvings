package main

import "fmt"

// ── Approach 1: Brute Force (Flood Fill / DFS Connected Components) ───────────
//
// floodFill solves Battleships in a Board by counting connected components of
// 'X' cells: each maximal blob of horizontally/vertically adjacent 'X's is one
// battleship. It mutates a visited copy so nothing in the input changes.
//
// Intuition:
//
//	Forget the "no two ships touch" guarantee for a moment: a battleship is
//	simply a connected group of 'X' cells (a 1×k or k×1 line, but the code
//	need not care about the shape). Counting ships = counting connected
//	components. Scan the board; each time we hit an unvisited 'X', that's a new
//	ship — flood-fill its whole body so we don't recount its cells.
//
// Algorithm:
//  1. For every cell: if it's an unvisited 'X', increment the ship count and
//     DFS-flood all 4-directionally connected 'X' cells, marking them visited.
//  2. Return the count.
//
// Time:  O(m·n) — every cell visited a constant number of times.
// Space: O(m·n) — visited grid plus recursion stack (a ship can span a row/col).
func floodFill(board [][]byte) int {
	m, n := len(board), len(board[0])
	visited := make([][]bool, m)
	for i := range visited {
		visited[i] = make([]bool, n)
	}
	dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} // N, S, W, E

	var dfs func(r, c int)
	dfs = func(r, c int) {
		// Stop at walls, water, or already-counted cells.
		if r < 0 || r >= m || c < 0 || c >= n || board[r][c] != 'X' || visited[r][c] {
			return
		}
		visited[r][c] = true // absorb this cell into the current ship
		for _, d := range dirs {
			dfs(r+d[0], c+d[1]) // spread along the ship's body
		}
	}

	count := 0
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if board[r][c] == 'X' && !visited[r][c] {
				count++   // first time we touch this ship
				dfs(r, c) // consume the rest of it
			}
		}
	}
	return count
}

// ── Approach 2: Count Ship "Heads" — One Pass, O(1) Space (Optimal) ──────────
//
// countHeads solves Battleships in a Board in a single pass using O(1) extra
// memory and WITHOUT modifying the board, by counting only each ship's
// top-left "head" cell.
//
// Intuition:
//
//	Every battleship has exactly one unique cell: its top-left end — the cell
//	with NO 'X' immediately above it and NO 'X' immediately to its left.
//	(Because ships are straight 1×k / k×1 lines separated by gaps, this head
//	is well defined and unique per ship.) So just count cells that are 'X' but
//	whose up-neighbour and left-neighbour are not 'X'. One scan, two look-back
//	checks, no visited array, no board mutation — satisfying the follow-up.
//
// Algorithm:
//  1. For each 'X' cell (r, c): it is a ship head iff
//     (r == 0 or board[r-1][c] != 'X') AND (c == 0 or board[r][c-1] != 'X').
//  2. Count heads; that's the number of ships.
//
// Time:  O(m·n) — a single pass, O(1) work per cell.
// Space: O(1) — only a counter (board untouched).
func countHeads(board [][]byte) int {
	m, n := len(board), len(board[0])
	count := 0
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if board[r][c] != 'X' {
				continue // water — skip
			}
			// A ship's body extends right/down from its head, so a cell with an
			// 'X' above OR to the left is a continuation, not a new ship.
			if r > 0 && board[r-1][c] == 'X' {
				continue // part of a vertical ship already counted at its top
			}
			if c > 0 && board[r][c-1] == 'X' {
				continue // part of a horizontal ship already counted at its left
			}
			count++ // this 'X' is the unique top-left head of a new ship
		}
	}
	return count
}

// parse turns a slice of ASCII rows into a [][]byte board for convenience.
func parse(rows []string) [][]byte {
	b := make([][]byte, len(rows))
	for i, r := range rows {
		b[i] = []byte(r)
	}
	return b
}

func main() {
	ex1 := parse([]string{
		"X..X",
		"...X",
		"...X",
	})
	ex2 := parse([]string{"."})

	fmt.Println("=== Approach 1: Flood Fill (connected components) ===")
	fmt.Println(floodFill(ex1)) // expected 2
	fmt.Println(floodFill(ex2)) // expected 0

	fmt.Println("=== Approach 2: Count Ship Heads (one pass, O(1) space) ===")
	fmt.Println(countHeads(ex1)) // expected 2
	fmt.Println(countHeads(ex2)) // expected 0
}
