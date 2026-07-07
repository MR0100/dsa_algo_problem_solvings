package main

import "fmt"

// ── Approach 1: BFS From Every Empty Cell (Brute Force) ───────────────────────
//
// bfsFromEmpty solves Shortest Distance from All Buildings by running a BFS
// outward from each empty land cell to measure the total distance to reach all
// buildings, then taking the minimum over all empty cells.
//
// Intuition:
//
//	The house must sit on an empty cell. For a candidate empty cell, its cost
//	is the sum of shortest path lengths to every building. A single BFS from
//	that cell finds all those shortest paths at once. If the BFS fails to reach
//	some building, the cell is invalid. Try every empty cell and keep the best.
//
// Algorithm:
//  1. Count total buildings.
//  2. For each empty cell (grid==0): BFS over passable empties, summing the
//     distance to each building reached; count how many buildings were reached.
//  3. If all buildings reached, update the answer with the summed distance.
//  4. Return the minimum, or -1 if none reaches all.
//
// Time:  O((m·n)^2) — one BFS (O(m·n)) started from each of up to m·n empties.
// Space: O(m·n) — the visited matrix and BFS queue.
func bfsFromEmpty(grid [][]int) int {
	m := len(grid)
	if m == 0 {
		return -1
	}
	n := len(grid[0])
	// Count all buildings so we know when a BFS has reached every one.
	totalBuildings := 0
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == 1 {
				totalBuildings++
			}
		}
	}
	dirs := [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	best := -1
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] != 0 {
				continue // house can only sit on empty land
			}
			// BFS from this empty cell.
			visited := make([][]bool, m)
			for r := range visited {
				visited[r] = make([]bool, n)
			}
			queue := [][2]int{{i, j}}
			visited[i][j] = true
			dist := 0      // current BFS layer distance
			totalDist := 0 // sum of distances to reached buildings
			reached := 0   // how many buildings reached
			for len(queue) > 0 {
				dist++ // step out one ring
				next := [][2]int{}
				for _, cell := range queue {
					for _, d := range dirs {
						nr, nc := cell[0]+d[0], cell[1]+d[1]
						if nr < 0 || nr >= m || nc < 0 || nc >= n || visited[nr][nc] {
							continue
						}
						visited[nr][nc] = true
						if grid[nr][nc] == 1 { // hit a building
							totalDist += dist
							reached++
						} else if grid[nr][nc] == 0 { // keep walking through empties
							next = append(next, [2]int{nr, nc})
						}
						// grid==2 (obstacle): marked visited, never expanded
					}
				}
				queue = next
			}
			if reached == totalBuildings { // valid meeting point
				if best == -1 || totalDist < best {
					best = totalDist
				}
			}
		}
	}
	return best
}

// ── Approach 2: BFS From Every Building (Optimal) ─────────────────────────────
//
// bfsFromBuildings solves the problem by running one BFS from EACH building,
// accumulating distances into every reachable empty cell, and requiring an
// empty cell to be reached by ALL buildings.
//
// Intuition:
//
//	Buildings are usually far fewer than empty cells, so starting BFS from
//	buildings is cheaper. For each building we add its shortest distance into a
//	`total[r][c]` accumulator for every empty cell, and bump a `reach[r][c]`
//	counter. After all buildings, an empty cell is a valid house iff its reach
//	count equals the number of buildings; its cost is `total[r][c]`. The trick
//	that keeps this a plain BFS (no per-building visited reset cost issue): use
//	an `emptyMarker` that decrements each building round, and only step onto a
//	cell whose current grid value equals the marker — this both enforces
//	"reachable by all previous buildings" and doubles as the visited flag.
//
// Algorithm:
//  1. total[r][c] = 0, and set emptyMarker = 0.
//  2. For each building: BFS; a neighbour is walkable iff grid[nr][nc] ==
//     emptyMarker. On visiting, add dist to total[nr][nc] and set grid[nr][nc]
//     = emptyMarker-1 (so only cells reached by ALL buildings so far survive).
//  3. emptyMarker--.
//  4. Answer = min total over cells whose grid == emptyMarker (reached by all).
//
// Time:  O(B · m · n) where B = number of buildings — one BFS per building.
// Space: O(m·n) — the total accumulator and BFS queue.
func bfsFromBuildings(grid [][]int) int {
	m := len(grid)
	if m == 0 {
		return -1
	}
	n := len(grid[0])
	total := make([][]int, m) // summed distance from all buildings so far
	for i := range total {
		total[i] = make([]int, n)
	}
	// Work on a copy so we don't mutate the caller's grid (and so the two
	// approaches see identical inputs in main).
	g := make([][]int, m)
	for i := range g {
		g[i] = make([]int, n)
		copy(g[i], grid[i])
	}
	dirs := [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	emptyMarker := 0 // cells walkable in the current round have this value
	best := -1
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if g[i][j] != 1 {
				continue // start BFS only from buildings
			}
			queue := [][2]int{{i, j}}
			dist := 0
			for len(queue) > 0 {
				dist++
				next := [][2]int{}
				for _, cell := range queue {
					for _, d := range dirs {
						nr, nc := cell[0]+d[0], cell[1]+d[1]
						// Walkable iff it is empty AND was reached by every
						// previous building (its marker matches emptyMarker).
						if nr < 0 || nr >= m || nc < 0 || nc >= n || g[nr][nc] != emptyMarker {
							continue
						}
						g[nr][nc]-- // consume: now only next building can step here
						total[nr][nc] += dist
						next = append(next, [2]int{nr, nc})
					}
				}
				queue = next
			}
			emptyMarker-- // tighten the "reachable by all" requirement
		}
	}
	// After all rounds, valid cells are exactly those whose marker equals
	// emptyMarker (reached by every building). Pick the minimum total.
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if g[i][j] == emptyMarker && total[i][j] > 0 {
				if best == -1 || total[i][j] < best {
					best = total[i][j]
				}
			}
		}
	}
	return best
}

func main() {
	// Example 1 (LeetCode canonical):
	//   1 - 0 - 2 - 0 - 1
	//   0 - 0 - 0 - 0 - 0
	//   0 - 0 - 1 - 0 - 0
	// Optimal house at (1,2) reaches all buildings in 3+3+1 = 7 steps.
	grid1 := [][]int{
		{1, 0, 2, 0, 1},
		{0, 0, 0, 0, 0},
		{0, 0, 1, 0, 0},
	}
	// Example 2: grid = [[1,0]] -> house at (0,1), distance 1.
	grid2 := [][]int{{1, 0}}
	// Example 3: grid = [[1]] -> no empty land, return -1.
	grid3 := [][]int{{1}}

	fmt.Println("=== Approach 1: BFS From Every Empty Cell ===")
	fmt.Printf("grid1 -> %d  expected 7\n", bfsFromEmpty(grid1))
	fmt.Printf("grid2 -> %d  expected 1\n", bfsFromEmpty(grid2))
	fmt.Printf("grid3 -> %d  expected -1\n", bfsFromEmpty(grid3))

	fmt.Println("=== Approach 2: BFS From Every Building (Optimal) ===")
	fmt.Printf("grid1 -> %d  expected 7\n", bfsFromBuildings(grid1))
	fmt.Printf("grid2 -> %d  expected 1\n", bfsFromBuildings(grid2))
	fmt.Printf("grid3 -> %d  expected -1\n", bfsFromBuildings(grid3))
}
