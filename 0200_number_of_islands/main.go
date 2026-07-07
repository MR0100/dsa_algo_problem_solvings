package main

import "fmt"

// cloneGrid deep-copies a grid so that mutating approaches (DFS/BFS sink the
// land cells they visit) each start from untouched input.
func cloneGrid(grid [][]byte) [][]byte {
	clone := make([][]byte, len(grid))
	for i, row := range grid {
		clone[i] = make([]byte, len(row))
		copy(clone[i], row)
	}
	return clone
}

// ── Approach 1: DFS Flood Fill ───────────────────────────────────────────────
//
// dfsFloodFill solves Number of Islands by scanning every cell and, on each
// still-unvisited land cell, recursively "sinking" its entire island.
//
// Intuition:
//   Whenever the scan touches a cell that is still '1', that cell belongs to
//   an island nobody has counted yet. Count it once, then flood-fill (DFS)
//   every land cell reachable from it, flipping '1' → '0' so the island can
//   never be counted again. Islands found = flood fills started.
//
// Algorithm:
//   1. For every cell (r, c): if grid[r][c] == '1', increment the counter and
//      call sinkDFS(r, c).
//   2. sinkDFS marks the cell '0' and recurses into its four orthogonal
//      neighbours that are still '1'.
//   3. Return the counter.
//
// Time:  O(m·n) — every cell is visited a constant number of times.
// Space: O(m·n) — recursion stack in the worst case (grid entirely land).
func dfsFloodFill(grid [][]byte) int {
	if len(grid) == 0 {
		return 0
	}
	count := 0
	for r := range grid {
		for c := range grid[r] {
			if grid[r][c] == '1' {
				count++            // a new, never-seen island starts here
				sinkDFS(grid, r, c) // erase it so it is counted exactly once
			}
		}
	}
	return count
}

// sinkDFS turns the whole island containing (r, c) into water.
func sinkDFS(grid [][]byte, r, c int) {
	// Stop at the grid border and at water / already-sunk cells.
	if r < 0 || r >= len(grid) || c < 0 || c >= len(grid[0]) || grid[r][c] != '1' {
		return
	}
	grid[r][c] = '0'      // sink: acts as the "visited" mark
	sinkDFS(grid, r+1, c) // down
	sinkDFS(grid, r-1, c) // up
	sinkDFS(grid, r, c+1) // right
	sinkDFS(grid, r, c-1) // left
}

// ── Approach 2: BFS Flood Fill ───────────────────────────────────────────────
//
// bfsFloodFill counts islands with the same sinking idea but explores each
// island with an iterative breadth-first queue instead of recursion.
//
// Intuition:
//   The counting argument is identical; BFS just expands the island ring by
//   ring. Prefer it when recursion depth is a concern — a 300×300 all-land
//   grid drives the DFS stack to 90,000 frames, while BFS keeps only a thin
//   frontier in memory.
//
// Algorithm:
//   1. Scan all cells; when grid[r][c] == '1': count++, mark it '0', and
//      enqueue (r, c).
//   2. While the queue is non-empty, dequeue a cell and, for each of its four
//      in-bounds neighbours still equal to '1', mark '0' immediately and
//      enqueue it (mark-on-enqueue prevents duplicate queue entries).
//   3. Return the counter.
//
// Time:  O(m·n) — every cell enters the queue at most once.
// Space: O(min(m,n)) — the BFS frontier grows at worst like a diagonal band
//         across the island (O(m·n) as a loose upper bound).
func bfsFloodFill(grid [][]byte) int {
	if len(grid) == 0 {
		return 0
	}
	rows, cols := len(grid), len(grid[0])
	dirs := [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} // the four orthogonal moves
	count := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if grid[r][c] != '1' {
				continue // water or an already-sunk island cell
			}
			count++          // new island discovered
			grid[r][c] = '0' // mark before enqueueing to avoid duplicates
			queue := [][2]int{{r, c}}
			for len(queue) > 0 {
				cell := queue[0]
				queue = queue[1:]
				for _, d := range dirs {
					nr, nc := cell[0]+d[0], cell[1]+d[1]
					// Sink in-bounds land neighbours the moment we see them.
					if nr >= 0 && nr < rows && nc >= 0 && nc < cols && grid[nr][nc] == '1' {
						grid[nr][nc] = '0'
						queue = append(queue, [2]int{nr, nc})
					}
				}
			}
		}
	}
	return count
}

// dsu is a disjoint-set union (union-find) structure with path compression
// and union by rank.
type dsu struct {
	parent []int // parent[i] = parent of i; i is a root when parent[i] == i
	rank   []int // rank[i] = upper bound on the height of the tree rooted at i
	count  int   // number of live components (islands not yet merged)
}

// newDSU builds a DSU over the grid with one singleton set per land cell.
func newDSU(grid [][]byte) *dsu {
	rows, cols := len(grid), len(grid[0])
	d := &dsu{parent: make([]int, rows*cols), rank: make([]int, rows*cols)}
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			id := r*cols + c // flatten (r, c) into a single integer id
			d.parent[id] = id
			if grid[r][c] == '1' {
				d.count++ // only land cells start as components
			}
		}
	}
	return d
}

// find returns the root of x, compressing the path as it walks.
func (d *dsu) find(x int) int {
	for d.parent[x] != x {
		d.parent[x] = d.parent[d.parent[x]] // path halving: skip a generation
		x = d.parent[x]
	}
	return x
}

// union merges the sets containing a and b; a successful merge removes one
// component from the count.
func (d *dsu) union(a, b int) {
	ra, rb := d.find(a), d.find(b)
	if ra == rb {
		return // already part of the same island
	}
	// Union by rank: attach the shorter tree under the taller one.
	if d.rank[ra] < d.rank[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	if d.rank[ra] == d.rank[rb] {
		d.rank[ra]++ // both trees equally tall → the merged tree grows by one
	}
	d.count-- // two islands fused into one
}

// ── Approach 3: Union-Find (Disjoint Set Union) ──────────────────────────────
//
// unionFind counts islands by starting every land cell as its own component
// and merging orthogonally adjacent land cells; the surviving component count
// is the number of islands.
//
// Intuition:
//   "Number of islands" is literally "number of connected components" in the
//   graph whose vertices are land cells and whose edges connect orthogonal
//   land neighbours. A DSU maintains component counts under edge insertions,
//   so one scan that unions every adjacent land pair yields the answer. This
//   is the approach that generalises to the dynamic follow-up (LeetCode #305)
//   where land is added one cell at a time.
//
// Algorithm:
//   1. Initialise the DSU: one singleton per land cell; count = land cells.
//   2. For every land cell (r, c), union it with its right (r, c+1) and down
//      (r+1, c) neighbours when they are land — those two directions cover
//      every adjacent pair exactly once.
//   3. Return the surviving component count.
//
// Time:  O(m·n · α(m·n)) — α is the inverse Ackermann function (≤ 5 for any
//         realistic input), so effectively linear.
// Space: O(m·n) — the parent and rank arrays.
func unionFind(grid [][]byte) int {
	if len(grid) == 0 {
		return 0
	}
	rows, cols := len(grid), len(grid[0])
	d := newDSU(grid)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if grid[r][c] != '1' {
				continue // only land cells create edges
			}
			if r+1 < rows && grid[r+1][c] == '1' {
				d.union(r*cols+c, (r+1)*cols+c) // edge to the cell below
			}
			if c+1 < cols && grid[r][c+1] == '1' {
				d.union(r*cols+c, r*cols+c+1) // edge to the cell on the right
			}
		}
	}
	return d.count
}

func main() {
	// Example 1 — one big island in the top-left corner.
	grid1 := [][]byte{
		[]byte("11110"),
		[]byte("11010"),
		[]byte("11000"),
		[]byte("00000"),
	}
	// Example 2 — three separate islands.
	grid2 := [][]byte{
		[]byte("11000"),
		[]byte("11000"),
		[]byte("00100"),
		[]byte("00011"),
	}

	fmt.Println("=== Approach 1: DFS Flood Fill ===")
	fmt.Println(dfsFloodFill(cloneGrid(grid1))) // 1
	fmt.Println(dfsFloodFill(cloneGrid(grid2))) // 3

	fmt.Println("=== Approach 2: BFS Flood Fill ===")
	fmt.Println(bfsFloodFill(cloneGrid(grid1))) // 1
	fmt.Println(bfsFloodFill(cloneGrid(grid2))) // 3

	fmt.Println("=== Approach 3: Union-Find ===")
	fmt.Println(unionFind(cloneGrid(grid1))) // 1
	fmt.Println(unionFind(cloneGrid(grid2))) // 3
}
