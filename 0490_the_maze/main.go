package main

import "fmt"

// The Maze (LeetCode 490): a ball starts at `start` and can roll up/down/left/
// right, but once it starts rolling it CANNOT stop until it hits a wall. From a
// stopped position it may pick a new direction. Return whether the ball can stop
// at `destination`. 0 = empty, 1 = wall; the border is walled.
//
// The crucial modelling point: the graph's nodes are STOPPING cells, and an
// edge connects two stop cells if the ball can roll from one straight into a
// wall and halt at the other. So each "move" rolls all the way, not one step.

// dirs lists the four roll directions as (dRow, dCol): up, down, left, right.
var dirs = [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

// roll slides the ball from (r,c) along (dr,dc) until the next cell would be a
// wall or off-grid, and returns the cell where it STOPS. Shared by all approaches.
func roll(maze [][]int, r, c, dr, dc int) (int, int) {
	m, n := len(maze), len(maze[0])
	// Keep advancing while the NEXT cell is inside the grid and empty.
	for r+dr >= 0 && r+dr < m && c+dc >= 0 && c+dc < n && maze[r+dr][c+dc] == 0 {
		r += dr // step one cell further in the roll direction
		c += dc
	}
	return r, c // first cell against a wall/border in this direction
}

// ── Approach 1: BFS over stopping cells ───────────────────────────────────────
//
// bfs solves The Maze by exploring stopping positions level by level with a
// queue, rolling the ball fully in each of the four directions from every
// popped stop cell.
//
// Intuition:
//
//	Treat each cell where the ball can come to rest as a graph node. From a
//	node, rolling in a direction lands on exactly one neighbour node (the wall
//	it slams into). Standard BFS from `start` over these roll-neighbours reaches
//	every stoppable cell; if `destination` is among them, answer true. A
//	`visited` grid keyed by stop cell prevents revisiting.
//
// Algorithm:
//  1. Queue starts with `start`; mark it visited.
//  2. Pop (r,c). If it equals destination, return true.
//  3. For each direction, roll to the stop cell (nr,nc); if unvisited, mark and enqueue.
//  4. Empty queue → destination unreachable → false.
//
// Time:  O(m·n·max(m,n)) — each of the m·n cells is enqueued once and each pop
//
//	rolls up to O(max(m,n)) cells across four directions.
//
// Space: O(m·n) — visited grid and queue.
func bfs(maze [][]int, start, destination []int) bool {
	m, n := len(maze), len(maze[0])
	visited := make([][]bool, m) // visited[r][c] = ball has stopped here before
	for i := range visited {
		visited[i] = make([]bool, n)
	}
	queue := [][2]int{{start[0], start[1]}} // FIFO of stop cells to expand
	visited[start[0]][start[1]] = true
	for len(queue) > 0 {
		cur := queue[0] // dequeue the oldest stop cell
		queue = queue[1:]
		if cur[0] == destination[0] && cur[1] == destination[1] {
			return true // reached a stop exactly at the destination
		}
		for _, d := range dirs {
			nr, nc := roll(maze, cur[0], cur[1], d[0], d[1]) // roll until a wall
			if !visited[nr][nc] {
				visited[nr][nc] = true                // record this stop cell
				queue = append(queue, [2]int{nr, nc}) // explore it later
			}
		}
	}
	return false // exhausted all reachable stops without hitting destination
}

// ── Approach 2: DFS over stopping cells ───────────────────────────────────────
//
// dfs solves The Maze with the same stop-cell graph but explores depth-first
// via recursion instead of a queue.
//
// Intuition:
//
//	Reachability doesn't care about path length, so DFS works just as well as
//	BFS: from the current stop cell, roll in all four directions and recurse
//	into each newly discovered stop cell, short-circuiting the moment the
//	destination is found. The `visited` grid again guards against cycles.
//
// Algorithm:
//  1. visit(r,c): if it is the destination, return true; mark visited.
//  2. For each direction, roll to (nr,nc); if unvisited, recurse; propagate true.
//  3. Return false if no direction leads to the destination.
//
// Time:  O(m·n·max(m,n)) — same node/edge bound as BFS.
// Space: O(m·n) — visited grid plus recursion depth up to O(m·n).
func dfs(maze [][]int, start, destination []int) bool {
	m, n := len(maze), len(maze[0])
	visited := make([][]bool, m)
	for i := range visited {
		visited[i] = make([]bool, n)
	}
	var visit func(r, c int) bool
	visit = func(r, c int) bool {
		if r == destination[0] && c == destination[1] {
			return true // stopped exactly on the destination
		}
		visited[r][c] = true // mark before recursing to avoid cycles
		for _, d := range dirs {
			nr, nc := roll(maze, r, c, d[0], d[1]) // roll fully in this direction
			if !visited[nr][nc] && visit(nr, nc) {
				return true // some deeper roll reached the destination
			}
		}
		return false // no direction from here works
	}
	return visit(start[0], start[1])
}

// ── Approach 3: BFS with Explicit Stop-Set (Optimal, same big-O) ──────────────
//
// bfsStopSet solves The Maze identically to Approach 1 but checks the
// destination when GENERATING each neighbour (and seeds visited up front),
// making the intent — "we only ever mark true stop cells" — explicit and
// avoiding a redundant destination test after dequeue.
//
// Intuition:
//
//	Functionally the same BFS; the refactor marks a stop cell visited exactly
//	when it is discovered and can early-return as soon as the destination is
//	first reached as a stop cell, which is the earliest possible moment.
//
// Algorithm:
//  1. Seed queue and visited with `start`; if start == destination, true.
//  2. Pop, roll in four directions; for each fresh stop cell, if it is the
//     destination return true, else mark visited and enqueue.
//  3. Empty queue → false.
//
// Time:  O(m·n·max(m,n)). Space: O(m·n).
func bfsStopSet(maze [][]int, start, destination []int) bool {
	if start[0] == destination[0] && start[1] == destination[1] {
		return true // degenerate: already at the destination (a valid stop)
	}
	m, n := len(maze), len(maze[0])
	visited := make([][]bool, m)
	for i := range visited {
		visited[i] = make([]bool, n)
	}
	visited[start[0]][start[1]] = true
	queue := [][2]int{{start[0], start[1]}}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, d := range dirs {
			nr, nc := roll(maze, cur[0], cur[1], d[0], d[1])
			if visited[nr][nc] {
				continue // already known stop cell
			}
			if nr == destination[0] && nc == destination[1] {
				return true // first time we can stop on the destination
			}
			visited[nr][nc] = true
			queue = append(queue, [2]int{nr, nc})
		}
	}
	return false
}

func main() {
	// Official maze (0 = empty, 1 = wall).
	maze := [][]int{
		{0, 0, 1, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 0, 0, 1, 0},
		{1, 1, 0, 1, 1},
		{0, 0, 0, 0, 0},
	}

	fmt.Println("=== Approach 1: BFS over stopping cells ===")
	fmt.Println(bfs(maze, []int{0, 4}, []int{4, 4})) // expected true
	fmt.Println(bfs(maze, []int{0, 4}, []int{3, 2})) // expected false

	fmt.Println("=== Approach 2: DFS over stopping cells ===")
	fmt.Println(dfs(maze, []int{0, 4}, []int{4, 4})) // expected true
	fmt.Println(dfs(maze, []int{0, 4}, []int{3, 2})) // expected false

	fmt.Println("=== Approach 3: BFS with Explicit Stop-Set (Optimal) ===")
	fmt.Println(bfsStopSet(maze, []int{0, 4}, []int{4, 4})) // expected true
	fmt.Println(bfsStopSet(maze, []int{0, 4}, []int{3, 2})) // expected false
}
