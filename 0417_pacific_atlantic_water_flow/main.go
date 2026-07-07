package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Brute Force (BFS/DFS From Every Cell) ─────────────────────────
//
// bruteForce solves Pacific Atlantic Water Flow by, for each cell, searching
// outward to check whether water starting there can reach the Pacific and the
// Atlantic separately.
//
// Intuition:
//
//	Water flows from a cell to a neighbour of equal-or-lower height. So a cell
//	drains to an ocean iff there exists a non-increasing path from it to that
//	ocean's border. Test every cell independently: run a flood (DFS) that only
//	steps "downhill or flat", and record whether it ever touches the top/left
//	edge (Pacific) and whether it touches the bottom/right edge (Atlantic). A
//	cell qualifies when both flags are set.
//
// Algorithm:
//  1. For each cell (r, c): DFS from it, moving to neighbours with height
//     <= current height; along the way note if we reach a Pacific border cell
//     and/or an Atlantic border cell.
//  2. If both oceans are reachable, add [r, c] to the result.
//
// Time:  O((m·n)²) — a full O(m·n) flood launched from each of m·n cells.
// Space: O(m·n) — visited grid + recursion stack per search.
func bruteForce(heights [][]int) [][]int {
	m, n := len(heights), len(heights[0])
	dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} // N, S, W, E
	var result [][]int

	// canReachBoth floods from (sr, sc) and returns whether the flood touches
	// each ocean's border.
	canReachBoth := func(sr, sc int) (pacific, atlantic bool) {
		visited := make([][]bool, m)
		for i := range visited {
			visited[i] = make([]bool, n)
		}
		var dfs func(r, c int)
		dfs = func(r, c int) {
			visited[r][c] = true
			if r == 0 || c == 0 { // top or left edge → Pacific
				pacific = true
			}
			if r == m-1 || c == n-1 { // bottom or right edge → Atlantic
				atlantic = true
			}
			for _, d := range dirs {
				nr, nc := r+d[0], c+d[1]
				if nr < 0 || nr >= m || nc < 0 || nc >= n {
					continue // off the board
				}
				if visited[nr][nc] {
					continue // already flooded
				}
				if heights[nr][nc] > heights[r][c] {
					continue // uphill — water can't flow there
				}
				dfs(nr, nc) // step downhill/flat
			}
		}
		dfs(sr, sc)
		return
	}

	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if p, a := canReachBoth(r, c); p && a {
				result = append(result, []int{r, c})
			}
		}
	}
	return result
}

// ── Approach 2: Reverse DFS From Ocean Borders (Optimal) ──────────────────────
//
// reverseDFS solves Pacific Atlantic Water Flow by flooding inward from each
// ocean's border, climbing to cells of equal-or-greater height, and returning
// the cells reachable from both floods.
//
// Intuition:
//
//	Instead of asking "can this cell reach the ocean?" for every cell, invert
//	the flow: start at the ocean and ask "which cells can send water HERE?"
//	Water flows downhill, so reversing it means climbing: from a border cell we
//	may move to a neighbour whose height is >= ours. One flood from the Pacific
//	border marks every cell that drains to the Pacific; one from the Atlantic
//	border marks every cell that drains to the Atlantic. The answer is the
//	intersection — two full-grid passes instead of m·n of them.
//
// Algorithm:
//  1. pacific[][] = cells reachable by climbing inward from top row + left col.
//  2. atlantic[][] = cells reachable by climbing inward from bottom row + right col.
//  3. Any cell marked in both drains to both oceans → add to result.
//
// Time:  O(m·n) — each cell is visited at most once per ocean flood.
// Space: O(m·n) — two boolean grids + recursion stack.
func reverseDFS(heights [][]int) [][]int {
	m, n := len(heights), len(heights[0])
	dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	pacific := make([][]bool, m)  // cells that can drain to the Pacific
	atlantic := make([][]bool, m) // cells that can drain to the Atlantic
	for i := 0; i < m; i++ {
		pacific[i] = make([]bool, n)
		atlantic[i] = make([]bool, n)
	}

	// dfs climbs inward: from (r,c) visit neighbours with height >= heights[r][c].
	var dfs func(r, c int, ocean [][]bool)
	dfs = func(r, c int, ocean [][]bool) {
		ocean[r][c] = true // this cell drains to the ocean we started from
		for _, d := range dirs {
			nr, nc := r+d[0], c+d[1]
			if nr < 0 || nr >= m || nc < 0 || nc >= n {
				continue // off the board
			}
			if ocean[nr][nc] {
				continue // already reached
			}
			if heights[nr][nc] < heights[r][c] {
				continue // reverse flow must go UPHILL-or-flat; lower neighbour can't feed us
			}
			dfs(nr, nc, ocean) // climb to the equal/higher neighbour
		}
	}

	// Seed the Pacific from the top row and left column.
	for c := 0; c < n; c++ {
		dfs(0, c, pacific) // top edge
	}
	for r := 0; r < m; r++ {
		dfs(r, 0, pacific) // left edge
	}
	// Seed the Atlantic from the bottom row and right column.
	for c := 0; c < n; c++ {
		dfs(m-1, c, atlantic) // bottom edge
	}
	for r := 0; r < m; r++ {
		dfs(r, n-1, atlantic) // right edge
	}

	var result [][]int
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if pacific[r][c] && atlantic[r][c] { // intersection of both floods
				result = append(result, []int{r, c})
			}
		}
	}
	return result
}

// ── Approach 3: Reverse BFS From Ocean Borders ───────────────────────────────
//
// reverseBFS solves Pacific Atlantic Water Flow with the same reverse-flow idea
// as Approach 2 but using an explicit queue (breadth-first) instead of
// recursion — avoids deep call stacks on large grids.
//
// Intuition:
//
//	Identical marking logic to the reverse DFS: seed a queue with all border
//	cells of an ocean, then repeatedly pop a cell and enqueue neighbours that
//	are equal-or-higher (water could flow from them down to us). BFS just
//	replaces the implicit recursion stack with a manual FIFO queue, which is
//	safer when m·n is large enough to overflow the call stack.
//
// Algorithm:
//  1. For each ocean, enqueue every border cell and mark it.
//  2. Pop cells; for each higher-or-equal unvisited neighbour, mark + enqueue.
//  3. Intersect the two reachable sets.
//
// Time:  O(m·n) — every cell enqueued at most once per ocean.
// Space: O(m·n) — the reachable grids plus the BFS queue.
func reverseBFS(heights [][]int) [][]int {
	m, n := len(heights), len(heights[0])
	dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	// bfs marks every cell that can reach the ocean whose border cells seed it.
	bfs := func(starts [][2]int) [][]bool {
		reach := make([][]bool, m)
		for i := range reach {
			reach[i] = make([]bool, n)
		}
		queue := make([][2]int, 0, len(starts))
		for _, s := range starts {
			reach[s[0]][s[1]] = true // border cells trivially drain to their ocean
			queue = append(queue, s)
		}
		for len(queue) > 0 {
			cell := queue[0] // dequeue front (FIFO)
			queue = queue[1:]
			r, c := cell[0], cell[1]
			for _, d := range dirs {
				nr, nc := r+d[0], c+d[1]
				if nr < 0 || nr >= m || nc < 0 || nc >= n {
					continue
				}
				if reach[nr][nc] {
					continue
				}
				if heights[nr][nc] < heights[r][c] {
					continue // neighbour lower → cannot flow uphill into us
				}
				reach[nr][nc] = true // neighbour drains to this ocean too
				queue = append(queue, [2]int{nr, nc})
			}
		}
		return reach
	}

	// Border seeds for each ocean.
	var pacStarts, atlStarts [][2]int
	for c := 0; c < n; c++ {
		pacStarts = append(pacStarts, [2]int{0, c})     // top row
		atlStarts = append(atlStarts, [2]int{m - 1, c}) // bottom row
	}
	for r := 0; r < m; r++ {
		pacStarts = append(pacStarts, [2]int{r, 0})     // left col
		atlStarts = append(atlStarts, [2]int{r, n - 1}) // right col
	}

	pacific := bfs(pacStarts)
	atlantic := bfs(atlStarts)

	var result [][]int
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if pacific[r][c] && atlantic[r][c] {
				result = append(result, []int{r, c})
			}
		}
	}
	return result
}

// sortCoords sorts coordinate pairs (row, then col) so outputs are comparable
// regardless of the order cells were discovered. LeetCode accepts any order;
// we sort only to make the printed lines deterministic against the expected
// comments.
func sortCoords(coords [][]int) [][]int {
	sort.Slice(coords, func(i, j int) bool {
		if coords[i][0] != coords[j][0] {
			return coords[i][0] < coords[j][0]
		}
		return coords[i][1] < coords[j][1]
	})
	return coords
}

func main() {
	ex1 := [][]int{
		{1, 2, 2, 3, 5},
		{3, 2, 3, 4, 4},
		{2, 4, 5, 3, 1},
		{6, 7, 1, 4, 5},
		{5, 1, 1, 2, 4},
	}
	ex2 := [][]int{{1}}

	fmt.Println("=== Approach 1: Brute Force (flood from every cell) ===")
	fmt.Println(sortCoords(bruteForce(ex1))) // expected [[0 4] [1 3] [1 4] [2 2] [3 0] [3 1] [4 0]]
	fmt.Println(sortCoords(bruteForce(ex2))) // expected [[0 0]]

	fmt.Println("=== Approach 2: Reverse DFS From Ocean Borders (Optimal) ===")
	fmt.Println(sortCoords(reverseDFS(ex1))) // expected [[0 4] [1 3] [1 4] [2 2] [3 0] [3 1] [4 0]]
	fmt.Println(sortCoords(reverseDFS(ex2))) // expected [[0 0]]

	fmt.Println("=== Approach 3: Reverse BFS From Ocean Borders ===")
	fmt.Println(sortCoords(reverseBFS(ex1))) // expected [[0 4] [1 3] [1 4] [2 2] [3 0] [3 1] [4 0]]
	fmt.Println(sortCoords(reverseBFS(ex2))) // expected [[0 0]]
}
