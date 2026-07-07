package main

import (
	"container/heap"
	"fmt"
)

// Problem 499 — The Maze III.
//
// A ball starts at `ball` in a 0/1 maze (1 = wall, 0 = empty). When pushed it
// rolls in one of four directions until it hits a wall, EXCEPT that if it rolls
// over the `hole` it drops in and stops there. We want the shortest path
// (fewest empty spaces travelled, hole counted, start excluded) that lands the
// ball in the hole; among equally short paths return the LEXICOGRAPHICALLY
// smallest direction string, or "impossible".
//
// This is a weighted shortest-path problem on the graph whose nodes are cells
// and whose edges are "roll until you stop", with edge weight = spaces rolled.
// Because a single push can travel many cells, edge weights differ, so plain
// BFS is wrong — we need Dijkstra (Approach 2) or repeated relaxation
// (Approach 1). The lexicographic tie-break is handled by comparing the path
// string alongside the distance.
//
// Directions are always considered in lexicographic order so that, on ties,
// the smaller letter wins: 'd' (down) < 'l' (left) < 'r' (right) < 'u' (up).

// dirs lists the four rolls in lexicographic letter order. Each entry is
// {row delta, col delta, letter}.
var dirs = []struct {
	dr, dc int
	ch     byte
}{
	{1, 0, 'd'},  // down
	{0, -1, 'l'}, // left
	{0, 1, 'r'},  // right
	{-1, 0, 'u'}, // up
}

// roll simulates pushing the ball from (r,c) in one direction. It advances one
// step at a time, stopping (a) the moment it lands on the hole, or (b) just
// before it would enter a wall / leave the grid. It returns the stop cell, the
// number of empty spaces travelled, and whether it dropped into the hole.
func roll(maze [][]int, r, c, dr, dc, holeR, holeC int) (int, int, int, bool) {
	m, n := len(maze), len(maze[0])
	steps := 0
	for {
		nr, nc := r+dr, c+dc // candidate next cell
		// stop if the next cell is off-grid or a wall
		if nr < 0 || nr >= m || nc < 0 || nc >= n || maze[nr][nc] == 1 {
			return r, c, steps, false // came to rest against a wall (not in hole)
		}
		r, c = nr, nc // move onto the empty cell
		steps++       // count this travelled space
		if r == holeR && c == holeC {
			return r, c, steps, true // fell into the hole mid-roll
		}
	}
}

// state is a Dijkstra priority-queue entry: the ball at (r,c), the total spaces
// travelled from the start (dist), and the direction string taken to get here.
type state struct {
	r, c int
	dist int
	path string
}

// pq is a min-heap of states ordered by (dist, path): fewer spaces first, and
// on ties the lexicographically smaller path first — exactly the problem's
// ranking, so the first time we pop the hole we have the optimal answer.
type pq []state

func (p pq) Len() int { return len(p) }
func (p pq) Less(i, j int) bool {
	if p[i].dist != p[j].dist {
		return p[i].dist < p[j].dist // primary: shortest distance
	}
	return p[i].path < p[j].path // tie-break: lexicographically smallest path
}
func (p pq) Swap(i, j int)       { p[i], p[j] = p[j], p[i] }
func (p *pq) Push(x interface{}) { *p = append(*p, x.(state)) }
func (p *pq) Pop() interface{} {
	old := *p
	n := len(old)
	item := old[n-1]
	*p = old[:n-1]
	return item
}

// ── Approach 1: Repeated Relaxation (Bellman-Ford style, Brute Force) ─────────
//
// bellmanFord solves The Maze III by relaxing every "roll" edge repeatedly
// until no (distance, path) pair improves, then reading the hole's best label.
//
// Intuition:
//
//	Treat each stop-cell as a node holding the best (dist, path) found so far.
//	From every settled cell, try all four rolls and relax the neighbour if the
//	new (dist, path) is better under the (shorter distance, else smaller path)
//	order. Repeat full sweeps until a sweep changes nothing — at most V passes,
//	Bellman-Ford style. Because we compare paths as part of "better", the fixed
//	point already encodes the lexicographic tie-break.
//
// Algorithm:
//  1. best[cell] = (∞, ""), best[ball] = (0, "").
//  2. Loop: for every cell with a finite label, roll in all 4 directions; relax
//     the stop cell if (newDist, newPath) beats its current label. Track if any
//     change happened this sweep.
//  3. Stop when a sweep makes no change.
//  4. Answer = best[hole].path, or "impossible" if still infinite.
//
// Time:  O(V · E · L) — up to V sweeps, each rolling every cell 4 directions
//
//	(each roll up to max(m,n) steps), with string compares of length L. For a
//	30×30 grid this is tiny.
//
// Space: O(V · L) — a (dist, path) label per cell.
func bellmanFord(maze [][]int, ball, hole []int) string {
	m, n := len(maze), len(maze[0])
	const inf = 1 << 30
	type label struct {
		dist int
		path string
		set  bool
	}
	best := make([][]label, m)
	for i := range best {
		best[i] = make([]label, n)
		for j := range best[i] {
			best[i][j] = label{dist: inf} // start everyone at infinity
		}
	}
	best[ball[0]][ball[1]] = label{dist: 0, path: "", set: true} // source label

	// better reports whether (d2,p2) is a strictly better label than (d1,p1)
	// under the (smaller dist, then lexicographically smaller path) ordering.
	better := func(d1 int, p1 string, d2 int, p2 string) bool {
		if d2 != d1 {
			return d2 < d1
		}
		return p2 < p1
	}

	for changed := true; changed; { // keep sweeping until a fixed point
		changed = false
		for r := 0; r < m; r++ {
			for c := 0; c < n; c++ {
				if !best[r][c].set { // skip unreached cells
					continue
				}
				cur := best[r][c]
				for _, d := range dirs { // try all four rolls, lexicographic order
					nr, nc, steps, _ := roll(maze, r, c, d.dr, d.dc, hole[0], hole[1])
					if steps == 0 {
						continue // immovable this direction (wall right there)
					}
					nd := cur.dist + steps        // spaces travelled to the stop cell
					np := cur.path + string(d.ch) // extend the direction string
					if !best[nr][nc].set || better(best[nr][nc].dist, best[nr][nc].path, nd, np) {
						best[nr][nc] = label{dist: nd, path: np, set: true}
						changed = true // relaxation happened → need another sweep
					}
				}
			}
		}
	}

	if !best[hole[0]][hole[1]].set {
		return "impossible" // hole never reached
	}
	return best[hole[0]][hole[1]].path
}

// ── Approach 2: Dijkstra with (distance, path) Priority Queue (Optimal) ───────
//
// dijkstra solves The Maze III by expanding stop-cells in nondecreasing
// (distance, path) order using a min-heap; the first time the hole is popped,
// its path is the optimal, lexicographically-smallest answer.
//
// Intuition:
//
//	Roll edges have varying weights, so we need Dijkstra rather than BFS. Order
//	the frontier by (dist, path): the very first extraction of the hole is
//	guaranteed to be both shortest and, among shortest, lexicographically
//	smallest — because the heap never pops a worse label before a better one. A
//	`seen` set of finalised cells prevents reprocessing.
//
// Algorithm:
//  1. Push (ball, dist=0, path="").
//  2. Pop the best state. If it is the hole, return its path.
//  3. If the cell is already finalised, skip; else finalise it.
//  4. Roll all four directions; push each stop cell with (dist+steps, path+ch).
//  5. If the heap empties without popping the hole, return "impossible".
//
// Time:  O(E log V + E · L) — heap ops with string comparisons of length L;
//
//	trivial for a 30×30 maze.
//
// Space: O(V · L) — heap entries and the seen grid.
func dijkstra(maze [][]int, ball, hole []int) string {
	m, n := len(maze), len(maze[0])
	seen := make([][]bool, m) // finalised cells (shortest label already settled)
	for i := range seen {
		seen[i] = make([]bool, n)
	}

	h := &pq{{r: ball[0], c: ball[1], dist: 0, path: ""}} // frontier starts at the ball
	heap.Init(h)

	for h.Len() > 0 {
		cur := heap.Pop(h).(state) // best (dist, path) not yet finalised
		if cur.r == hole[0] && cur.c == hole[1] {
			return cur.path // first hole pop = optimal lexicographic answer
		}
		if seen[cur.r][cur.c] {
			continue // an equal-or-better label for this cell was already settled
		}
		seen[cur.r][cur.c] = true // finalise this cell

		for _, d := range dirs { // expand all four rolls (lexicographic order)
			nr, nc, steps, _ := roll(maze, cur.r, cur.c, d.dr, d.dc, hole[0], hole[1])
			if steps == 0 || seen[nr][nc] {
				continue // no movement, or neighbour already finalised
			}
			heap.Push(h, state{
				r:    nr,
				c:    nc,
				dist: cur.dist + steps,        // add spaces rolled this push
				path: cur.path + string(d.ch), // append the direction letter
			})
		}
	}
	return "impossible" // heap drained without reaching the hole
}

func main() {
	maze := [][]int{
		{0, 0, 0, 0, 0},
		{1, 1, 0, 0, 1},
		{0, 0, 0, 0, 0},
		{0, 1, 0, 0, 1},
		{0, 1, 0, 0, 0},
	}

	fmt.Println("=== Approach 1: Repeated Relaxation (Bellman-Ford style) ===")
	fmt.Println(bellmanFord(maze, []int{4, 3}, []int{0, 1})) // expected lul
	fmt.Println(bellmanFord(maze, []int{4, 3}, []int{3, 0})) // expected impossible

	fmt.Println("=== Approach 2: Dijkstra with (dist, path) PQ (Optimal) ===")
	fmt.Println(dijkstra(maze, []int{4, 3}, []int{0, 1})) // expected lul
	fmt.Println(dijkstra(maze, []int{4, 3}, []int{3, 0})) // expected impossible
}
