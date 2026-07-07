package main

import "fmt"

// ── Approach 1: Cell Scan, Count Exposed Sides ───────────────────────────────
//
// cellScan solves Island Perimeter by visiting every land cell and, for each
// of its four sides, adding 1 to the perimeter whenever that side faces water
// or the grid boundary.
//
// Intuition:
//
//	Each land cell is a 1×1 square with 4 unit sides. A side belongs to the
//	perimeter exactly when the neighbour on that side is NOT land — i.e. it is
//	water or off the edge of the grid. So for every land cell, look at its 4
//	neighbours and count how many are non-land; that count is the cell's
//	contribution to the perimeter.
//
// Algorithm:
//  1. perimeter = 0.
//  2. For each cell (r, c) with grid[r][c] == 1:
//     for each of the 4 directions, if the neighbour is out of bounds or 0,
//     perimeter++.
//  3. Return perimeter.
//
// Time:  O(rows · cols) — a constant 4 checks per cell.
// Space: O(1) — only a counter.
func cellScan(grid [][]int) int {
	rows, cols := len(grid), len(grid[0])
	// The four orthogonal neighbour offsets: up, down, left, right.
	dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	perimeter := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if grid[r][c] != 1 {
				continue // water contributes nothing
			}
			for _, d := range dirs {
				nr, nc := r+d[0], c+d[1] // neighbour coordinates
				// A side is on the perimeter if the neighbour is off-grid or water.
				if nr < 0 || nr >= rows || nc < 0 || nc >= cols || grid[nr][nc] == 0 {
					perimeter++
				}
			}
		}
	}
	return perimeter
}

// ── Approach 2: Contribution Counting (Optimal) ──────────────────────────────
//
// contributionCount solves Island Perimeter with the identity
// perimeter = 4·(land cells) − 2·(adjacent land pairs).
//
// Intuition:
//
//	Start by giving every land cell its full 4 sides. Whenever two land cells
//	are adjacent, they share one internal edge — that edge is counted twice
//	(once from each cell) but belongs to neither's perimeter, so subtract 2.
//	Count each shared edge exactly once by only looking LEFT and UP from each
//	cell (the right/down neighbours will look back). One pass, no direction
//	array to re-scan.
//
// Algorithm:
//  1. land = 0, shared = 0.
//  2. For each cell (r, c) with grid[r][c] == 1:
//     land++; if the cell above is land, shared++; if the cell to the left is
//     land, shared++.
//  3. Return 4·land − 2·shared.
//
// Time:  O(rows · cols) — one pass, constant work per cell.
// Space: O(1).
func contributionCount(grid [][]int) int {
	rows, cols := len(grid), len(grid[0])
	land, shared := 0, 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if grid[r][c] != 1 {
				continue
			}
			land++ // this cell brings 4 sides to the table
			if r > 0 && grid[r-1][c] == 1 {
				shared++ // shares an edge with the land cell above
			}
			if c > 0 && grid[r][c-1] == 1 {
				shared++ // shares an edge with the land cell to the left
			}
		}
	}
	// Each shared edge was double-counted, and it is not part of the perimeter.
	return 4*land - 2*shared
}

// ── Approach 3: DFS Flood Fill (Count Boundary Edges) ────────────────────────
//
// dfsPerimeter solves Island Perimeter by flood-filling the single island and,
// at each land cell, adding 1 for every side that borders water or the edge.
//
// Intuition:
//
//	Because the grid has exactly one island, a DFS from any land cell visits
//	all of it. As DFS crosses each cell, count its perimeter-contributing
//	sides: stepping off the grid or into water is a boundary edge (+1);
//	stepping into unvisited land recurses; stepping into visited land does
//	nothing (that shared edge was already handled). A "visited" grid prevents
//	infinite recursion. Same answer as Approach 1, but framed as graph search
//	— the natural template when the island itself must be explored.
//
// Algorithm:
//  1. Find any land cell; DFS from it.
//  2. In dfs(r, c): if (r, c) is off-grid or water, return 1 (a boundary edge).
//     If already visited, return 0. Otherwise mark visited and return the sum
//     of dfs over the 4 neighbours.
//  3. The DFS return value is the perimeter.
//
// Time:  O(rows · cols) — each cell visited once; boundary lookups are O(1).
// Space: O(rows · cols) — visited grid plus recursion stack.
func dfsPerimeter(grid [][]int) int {
	rows, cols := len(grid), len(grid[0])
	visited := make([][]bool, rows) // guards against revisiting land cells
	for i := range visited {
		visited[i] = make([]bool, cols)
	}

	var dfs func(r, c int) int
	dfs = func(r, c int) int {
		// Off-grid or water: the side we crossed to get here is a boundary edge.
		if r < 0 || r >= rows || c < 0 || c >= cols || grid[r][c] == 0 {
			return 1
		}
		if visited[r][c] {
			return 0 // internal shared edge, already accounted for
		}
		visited[r][c] = true // mark before recursing to avoid cycles
		// Sum boundary edges contributed through each of the 4 sides.
		return dfs(r-1, c) + dfs(r+1, c) + dfs(r, c-1) + dfs(r, c+1)
	}

	// Launch DFS from the first land cell (exactly one island guaranteed).
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if grid[r][c] == 1 {
				return dfs(r, c)
			}
		}
	}
	return 0 // no land at all (not possible under constraints)
}

func main() {
	grid1 := [][]int{{0, 1, 0, 0}, {1, 1, 1, 0}, {0, 1, 0, 0}, {1, 1, 0, 0}}
	grid2 := [][]int{{1}}
	grid3 := [][]int{{1, 0}}

	fmt.Println("=== Approach 1: Cell Scan, Count Exposed Sides ===")
	fmt.Printf("grid1  got=%d  expected 16\n", cellScan(grid1))
	fmt.Printf("[[1]]  got=%d  expected 4\n", cellScan(grid2))
	fmt.Printf("[[1,0]] got=%d  expected 4\n", cellScan(grid3))

	fmt.Println("=== Approach 2: Contribution Counting (Optimal) ===")
	fmt.Printf("grid1  got=%d  expected 16\n", contributionCount(grid1))
	fmt.Printf("[[1]]  got=%d  expected 4\n", contributionCount(grid2))
	fmt.Printf("[[1,0]] got=%d  expected 4\n", contributionCount(grid3))

	fmt.Println("=== Approach 3: DFS Flood Fill (Count Boundary Edges) ===")
	fmt.Printf("grid1  got=%d  expected 16\n", dfsPerimeter(grid1))
	fmt.Printf("[[1]]  got=%d  expected 4\n", dfsPerimeter(grid2))
	fmt.Printf("[[1,0]] got=%d  expected 4\n", dfsPerimeter(grid3))
}
