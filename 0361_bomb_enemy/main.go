package main

import "fmt"

// ── Approach 1: Brute Force (Scan Four Directions Per Cell) ───────────────────
//
// bruteForce solves Bomb Enemy by, for every empty cell, walking outward in all
// four directions and counting enemies until a wall or the border stops it.
//
// Intuition:
//
//	The problem literally describes a simulation: drop the bomb on an empty
//	cell, and it kills every enemy in the same row/column up to the nearest
//	wall. So try every empty cell and simulate the blast directly.
//
// Algorithm:
//  1. For each cell (r,c) that is empty ('0'):
//     a. Walk up, down, left, right from (r,c).
//     b. In each direction, count 'E' cells; stop at a 'W' wall or the border.
//  2. Track the maximum total kills over all empty cells.
//
// Time:  O(m*n*(m+n)) — each of the m*n cells may scan up to O(m+n) cells.
// Space: O(1) — only counters.
func bruteForce(grid [][]byte) int {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return 0 // no cells → nothing to bomb
	}
	m, n := len(grid), len(grid[0])
	best := 0
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if grid[r][c] != '0' {
				continue // bomb can only be placed on an empty cell
			}
			kills := 0
			// Walk up: decrement row until wall or top border.
			for i := r - 1; i >= 0 && grid[i][c] != 'W'; i-- {
				if grid[i][c] == 'E' {
					kills++ // an enemy in the blast path
				}
			}
			// Walk down.
			for i := r + 1; i < m && grid[i][c] != 'W'; i++ {
				if grid[i][c] == 'E' {
					kills++
				}
			}
			// Walk left.
			for j := c - 1; j >= 0 && grid[r][j] != 'W'; j-- {
				if grid[r][j] == 'E' {
					kills++
				}
			}
			// Walk right.
			for j := c + 1; j < n && grid[r][j] != 'W'; j++ {
				if grid[r][j] == 'E' {
					kills++
				}
			}
			if kills > best {
				best = kills // remember the strongest placement
			}
		}
	}
	return best
}

// ── Approach 2: Row/Column Running Count (Optimal) ───────────────────────────
//
// runningCount solves Bomb Enemy by reusing a row's enemy count until a wall
// resets it, plus a per-column running count, so each cell is O(1).
//
// Intuition:
//
//	Scanning left-to-right, the number of enemies a bomb kills to its LEFT
//	(up to the nearest wall) only changes when we cross a wall. So keep a
//	single `rowHits` that accumulates enemies and resets at every 'W'; it is
//	valid for every empty cell in that wall-bounded segment. The same idea
//	works per column with a `colHits[c]` array, recomputed whenever the top
//	of a column segment is reached. Adding rowHits + colHits[c] at an empty
//	cell gives the total kills for that cell in O(1).
//
// Algorithm:
//  1. Keep `rowHits` (enemies in current row segment) and `colHits[c]`
//     (enemies in current column segment for column c).
//  2. Traverse cells row by row, left to right. At (r,c):
//     - If it starts a new row segment (c==0 or left neighbour is 'W'),
//     recount enemies rightward until the next wall → rowHits.
//     - If it starts a new column segment (r==0 or cell above is 'W'),
//     recount enemies downward until the next wall → colHits[c].
//  3. On a wall, skip. On an empty cell, candidate = rowHits + colHits[c].
//  4. Return the max candidate.
//
// Time:  O(m*n) — each cell is visited O(1) times amortised (segment recounts
//
//	are charged once per segment).
//
// Space: O(n) — the per-column running-count array.
func runningCount(grid [][]byte) int {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return 0
	}
	m, n := len(grid), len(grid[0])
	best := 0
	rowHits := 0              // enemies in the current row segment
	colHits := make([]int, n) // enemies in each column's current segment
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			// Recompute rowHits at the start of a new row segment: either the
			// left border, or just after a wall on the left.
			if c == 0 || grid[r][c-1] == 'W' {
				rowHits = 0
				for k := c; k < n && grid[r][k] != 'W'; k++ {
					if grid[r][k] == 'E' {
						rowHits++ // enemies reachable rightward before a wall
					}
				}
			}
			// Recompute colHits[c] at the start of a new column segment: either
			// the top border, or just after a wall above.
			if r == 0 || grid[r-1][c] == 'W' {
				colHits[c] = 0
				for k := r; k < m && grid[k][c] != 'W'; k++ {
					if grid[k][c] == 'E' {
						colHits[c]++ // enemies reachable downward before a wall
					}
				}
			}
			// Only empty cells can hold the bomb; sum row + column kills.
			if grid[r][c] == '0' {
				if total := rowHits + colHits[c]; total > best {
					best = total
				}
			}
		}
	}
	return best
}

func main() {
	// Example 1: expected 3 (place bomb at grid[1][1]).
	grid1 := [][]byte{
		{'0', 'E', '0', '0'},
		{'E', '0', 'W', 'E'},
		{'0', 'E', '0', '0'},
	}
	// Example 2: expected 1 (walls block the enemies row from combining).
	grid2 := [][]byte{
		{'W', 'W', 'W'},
		{'0', '0', '0'},
		{'E', 'E', 'E'},
	}

	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(bruteForce(grid1)) // expected 3
	fmt.Println(bruteForce(grid2)) // expected 1

	fmt.Println("=== Approach 2: Row/Column Running Count (Optimal) ===")
	fmt.Println(runningCount(grid1)) // expected 3
	fmt.Println(runningCount(grid2)) // expected 1
}
