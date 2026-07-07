package main

import "fmt"

// LeetCode #289 — Game of Life
//
// Given an m x n board of cells (1 = live, 0 = dead), compute the NEXT state
// using Conway's rules, applied simultaneously to every cell:
//   1. A live cell with < 2 live neighbours dies (underpopulation).
//   2. A live cell with 2 or 3 live neighbours lives on.
//   3. A live cell with > 3 live neighbours dies (overpopulation).
//   4. A dead cell with exactly 3 live neighbours becomes live (reproduction).
// Neighbours are the 8 surrounding cells. Update the board in place.

// ── Approach 1: Extra Copy Buffer ────────────────────────────────────────────
//
// extraCopy computes the next state by reading neighbour counts from an
// untouched snapshot copy, writing results into the real board.
//
// Intuition:
//
//	All cells must update "simultaneously". The simplest way to guarantee we
//	always read the OLD state while writing the NEW one is to keep a full
//	copy of the old board and read counts from it.
//
// Algorithm:
//  1. Copy the board into `snap`.
//  2. For each cell, count its 8 live neighbours in `snap`.
//  3. Apply Conway's rules; write 0/1 into the real board.
//
// Time:  O(m*n) — 8 neighbour reads per cell.
// Space: O(m*n) — the snapshot copy.
func extraCopy(board [][]int) {
	if len(board) == 0 || len(board[0]) == 0 {
		return
	}
	m, n := len(board), len(board[0])
	// Snapshot of the original board — the source of truth for counting.
	snap := make([][]int, m)
	for i := range board {
		snap[i] = make([]int, n)
		copy(snap[i], board[i])
	}
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			live := countLiveNeighbors(snap, i, j) // count from the snapshot
			// Live cell survives on 2 or 3 neighbours; dead cell born on 3.
			if snap[i][j] == 1 {
				if live == 2 || live == 3 {
					board[i][j] = 1
				} else {
					board[i][j] = 0
				}
			} else {
				if live == 3 {
					board[i][j] = 1
				} else {
					board[i][j] = 0
				}
			}
		}
	}
}

// countLiveNeighbors counts live cells among the 8 neighbours of (r, c) in g.
func countLiveNeighbors(g [][]int, r, c int) int {
	m, n := len(g), len(g[0])
	count := 0
	for dr := -1; dr <= 1; dr++ {
		for dc := -1; dc <= 1; dc++ {
			if dr == 0 && dc == 0 {
				continue // skip the cell itself
			}
			nr, nc := r+dr, c+dc
			if nr >= 0 && nr < m && nc >= 0 && nc < n && g[nr][nc] == 1 {
				count++
			}
		}
	}
	return count
}

// ── Approach 2: In-Place 2-Bit State Encoding (Optimal, O(1) Space) ──────────
//
// inPlaceBits computes the next state in place using O(1) extra space by
// encoding both the old and new bit of each cell in two bits.
//
// Intuition:
//
//	We need each cell to remember its old value (for neighbours still to be
//	processed) AND its new value. Pack both: keep the current state in bit 0
//	(value & 1) and stash the NEXT state in bit 1. Because we only ever read
//	bit 0 while counting, the snapshot is preserved without a copy. A final
//	pass shifts every cell right by one bit to reveal the next state.
//
//	Encoding: bit0 = old state, bit1 = new state.
//	  0 (00) old dead, new dead
//	  1 (01) old live, new dead
//	  2 (10) old dead, new live
//	  3 (11) old live, new live
//
// Algorithm:
//  1. For each cell, count live neighbours using (val & 1) so only old bits
//     are read.
//  2. If the cell will be live next, set bit 1 (val |= 2).
//  3. After the full sweep, shift every cell right by 1 (val >>= 1).
//
// Time:  O(m*n) — 8 neighbour reads per cell plus a final shift pass.
// Space: O(1) — no extra board; state packed into spare bits.
func inPlaceBits(board [][]int) {
	if len(board) == 0 || len(board[0]) == 0 {
		return
	}
	m, n := len(board), len(board[0])
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			live := 0
			// Count neighbours reading ONLY bit 0 (the original state).
			for dr := -1; dr <= 1; dr++ {
				for dc := -1; dc <= 1; dc++ {
					if dr == 0 && dc == 0 {
						continue
					}
					nr, nc := i+dr, j+dc
					if nr >= 0 && nr < m && nc >= 0 && nc < n {
						live += board[nr][nc] & 1 // old state lives in bit 0
					}
				}
			}
			old := board[i][j] & 1
			// Set bit 1 when the cell will be live in the next generation.
			if old == 1 && (live == 2 || live == 3) {
				board[i][j] |= 2 // live stays live
			} else if old == 0 && live == 3 {
				board[i][j] |= 2 // dead becomes live
			}
			// otherwise bit 1 stays 0 → next state dead
		}
	}
	// Reveal the next state: shift the stored bit 1 down into bit 0.
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			board[i][j] >>= 1
		}
	}
}

func deepCopy(g [][]int) [][]int {
	c := make([][]int, len(g))
	for i := range g {
		c[i] = make([]int, len(g[i]))
		copy(c[i], g[i])
	}
	return c
}

func main() {
	// Example 1:
	// Input:  [[0,1,0],[0,0,1],[1,1,1],[0,0,0]]
	// Output: [[0,0,0],[1,0,1],[0,1,1],[0,1,0]]
	example1 := [][]int{
		{0, 1, 0},
		{0, 0, 1},
		{1, 1, 1},
		{0, 0, 0},
	}
	// Example 2:
	// Input:  [[1,1],[1,0]]
	// Output: [[1,1],[1,1]]
	example2 := [][]int{
		{1, 1},
		{1, 0},
	}

	fmt.Println("=== Approach 1: Extra Copy Buffer ===")
	a1 := deepCopy(example1)
	extraCopy(a1)
	fmt.Println(a1) // expected [[0 0 0] [1 0 1] [0 1 1] [0 1 0]]
	a2 := deepCopy(example2)
	extraCopy(a2)
	fmt.Println(a2) // expected [[1 1] [1 1]]

	fmt.Println("=== Approach 2: In-Place 2-Bit Encoding (Optimal) ===")
	b1 := deepCopy(example1)
	inPlaceBits(b1)
	fmt.Println(b1) // expected [[0 0 0] [1 0 1] [0 1 1] [0 1 0]]
	b2 := deepCopy(example2)
	inPlaceBits(b2)
	fmt.Println(b2) // expected [[1 1] [1 1]]
}
