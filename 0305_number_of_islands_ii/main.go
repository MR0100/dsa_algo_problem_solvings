package main

import "fmt"

// LeetCode 305 — Number of Islands II
//
// An m x n grid starts as all water. We perform a sequence of "add land"
// operations at positions[i] = [r, c]. After each operation, report the current
// number of islands (4-directionally connected land groups). Return the list of
// counts, one per operation.

// ── Approach 1: BFS/Flood Fill After Each Operation (Brute Force) ─────────────
//
// floodFillCount rebuilds the island count from scratch after every add-land by
// flood-filling the whole grid.
//
// Intuition:
//
//	The most direct approach: maintain the actual grid, and after adding each
//	piece of land, recount islands by scanning every cell and BFS-flooding
//	unvisited land. Correct but expensive — each recount is O(m·n), so the
//	total is O(k·m·n) for k operations.
//
// Algorithm:
//  1. Keep a boolean grid; for each position, set that cell to land.
//  2. Recount: iterate all cells; each unvisited land cell starts a BFS that
//     marks its whole island; increment the count per island found.
//  3. Append the count after each operation.
//
// Time:  O(k·m·n) — k full grid scans.
// Space: O(m·n) — grid + visited.
func floodFillCount(m, n int, positions [][]int) []int {
	grid := make([][]bool, m) // true = land
	for i := range grid {
		grid[i] = make([]bool, n)
	}
	dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} // 4-neighbours
	result := make([]int, 0, len(positions))

	for _, p := range positions {
		grid[p[0]][p[1]] = true // add this land

		// Recount islands over the whole grid with BFS flood fill.
		visited := make([][]bool, m)
		for i := range visited {
			visited[i] = make([]bool, n)
		}
		count := 0
		for r := 0; r < m; r++ {
			for c := 0; c < n; c++ {
				if grid[r][c] && !visited[r][c] { // new island root
					count++
					queue := [][2]int{{r, c}} // BFS frontier
					visited[r][c] = true
					for len(queue) > 0 {
						cell := queue[0]
						queue = queue[1:]
						for _, d := range dirs {
							nr, nc := cell[0]+d[0], cell[1]+d[1]
							if nr >= 0 && nr < m && nc >= 0 && nc < n &&
								grid[nr][nc] && !visited[nr][nc] {
								visited[nr][nc] = true
								queue = append(queue, [2]int{nr, nc})
							}
						}
					}
				}
			}
		}
		result = append(result, count)
	}
	return result
}

// ── Approach 2: Union-Find / Disjoint Set (Optimal) ──────────────────────────
//
// unionFind maintains a running island count incrementally as land is added,
// merging with already-present neighbours via a disjoint-set structure.
//
// Intuition:
//
//	Adding one land cell provisionally creates a new island, so count++.
//	Then, for each of its up-to-four already-land neighbours, if the neighbour
//	belongs to a DIFFERENT island, merging them removes one island (count--).
//	Union-Find with path compression + union by rank makes each find/union
//	near O(1) (inverse-Ackermann), so total work is O(k·α(m·n)) — no rescans.
//	Each cell (r, c) maps to the flat index r*n + c.
//
// Algorithm:
//  1. parent[i] = i means "not yet land"; use a separate seen[] to mark land.
//  2. For each position: if already land, repeat the current count and skip.
//     Otherwise mark it land, count++.
//  3. For each land neighbour, union it in; every successful merge does count--.
//  4. Append count after each operation.
//
// Time:  O(k·α(m·n)) ≈ O(k) — near-constant find/union per operation.
// Space: O(m·n) — parent + rank + seen arrays.
func unionFind(m, n int, positions [][]int) []int {
	total := m * n
	parent := make([]int, total) // disjoint-set parent pointers
	rank := make([]int, total)   // tree height hint for union by rank
	seen := make([]bool, total)  // whether a cell is land yet
	for i := range parent {
		parent[i] = i // each node is initially its own root
	}

	var find func(x int) int
	find = func(x int) int {
		for parent[x] != x { // walk up to the root
			parent[x] = parent[parent[x]] // path halving (compression)
			x = parent[x]
		}
		return x
	}
	// union merges two sets; returns true iff they were previously separate.
	union := func(a, b int) bool {
		ra, rb := find(a), find(b)
		if ra == rb {
			return false // already the same island — no merge
		}
		if rank[ra] < rank[rb] { // attach shorter tree under taller
			ra, rb = rb, ra
		}
		parent[rb] = ra
		if rank[ra] == rank[rb] {
			rank[ra]++ // heights tied → resulting tree grew by one
		}
		return true
	}

	dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	result := make([]int, 0, len(positions))
	count := 0 // running number of islands

	for _, p := range positions {
		r, c := p[0], p[1]
		idx := r*n + c
		if seen[idx] { // duplicate add-land: island count unchanged
			result = append(result, count)
			continue
		}
		seen[idx] = true
		count++ // provisionally a brand-new island

		for _, d := range dirs {
			nr, nc := r+d[0], c+d[1]
			if nr < 0 || nr >= m || nc < 0 || nc >= n {
				continue // off-grid neighbour
			}
			nidx := nr*n + nc
			if !seen[nidx] {
				continue // neighbour is still water
			}
			if union(idx, nidx) { // merged two distinct islands
				count-- // two islands became one
			}
		}
		result = append(result, count)
	}
	return result
}

func main() {
	fmt.Println("=== Approach 1: BFS Flood Fill After Each Op ===")
	fmt.Println(floodFillCount(3, 3, [][]int{{0, 0}, {0, 1}, {1, 2}, {2, 1}})) // expected [1 1 2 3]
	fmt.Println(floodFillCount(1, 1, [][]int{{0, 0}}))                         // expected [1]
	// Duplicate add-land at (0,0) should not change the count.
	fmt.Println(floodFillCount(2, 2, [][]int{{0, 0}, {0, 0}, {1, 1}})) // expected [1 1 2]

	fmt.Println("=== Approach 2: Union-Find (Optimal) ===")
	fmt.Println(unionFind(3, 3, [][]int{{0, 0}, {0, 1}, {1, 2}, {2, 1}})) // expected [1 1 2 3]
	fmt.Println(unionFind(1, 1, [][]int{{0, 0}}))                         // expected [1]
	fmt.Println(unionFind(2, 2, [][]int{{0, 0}, {0, 0}, {1, 1}}))         // expected [1 1 2]
}
