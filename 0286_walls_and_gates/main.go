package main

import "fmt"

// LeetCode #286 — Walls and Gates
//
// You are given an m x n grid `rooms` initialized with these three values:
//   -1  = a wall or an obstacle.
//    0  = a gate.
//   INF = 2147483647 (2^31 - 1) = an empty room.
//
// Fill each empty room with the distance to its nearest gate. If it is
// impossible to reach a gate, leave INF.

const INF = 2147483647 // 2^31 - 1, marks an unreachable / not-yet-visited empty room

// ── Approach 1: Brute Force BFS From Every Empty Room ────────────────────────
//
// bruteForceBFS solves Walls and Gates by, for EACH empty room, running a BFS
// outward until it hits the first gate, and recording that distance.
//
// Intuition:
//
//	The distance we want is "shortest path to nearest gate". BFS from a cell
//	finds the shortest path to everything. So do a BFS from every empty room;
//	the first gate it reaches is the nearest one. Correct but wasteful: no
//	sharing of work between rooms.
//
// Algorithm:
//  1. For every cell that is an empty room (== INF):
//     a. BFS outward, tracking distance, stepping only onto in-bounds
//     non-wall cells.
//     b. Stop at the first gate (== 0); write that distance into the room.
//  2. Rooms that never reach a gate keep INF.
//
// Time:  O((m*n)^2) — a BFS costing O(m*n) launched from up to m*n rooms.
// Space: O(m*n) — the visited grid and queue for one BFS.
func bruteForceBFS(rooms [][]int) {
	if len(rooms) == 0 || len(rooms[0]) == 0 {
		return
	}
	m, n := len(rooms), len(rooms[0])
	dirs := [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} // 4-neighbour moves

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if rooms[i][j] != INF { // only launch a search from empty rooms
				continue
			}
			// BFS from (i, j) to find the nearest gate.
			visited := make([][]bool, m)
			for k := range visited {
				visited[k] = make([]bool, n)
			}
			type cell struct{ r, c, d int } // row, col, distance-from-start
			queue := []cell{{i, j, 0}}
			visited[i][j] = true
			found := false
			for len(queue) > 0 && !found {
				cur := queue[0]
				queue = queue[1:]
				if rooms[cur.r][cur.c] == 0 { // reached a gate
					rooms[i][j] = cur.d // record distance in the origin room
					found = true
					break
				}
				for _, dir := range dirs {
					nr, nc := cur.r+dir[0], cur.c+dir[1]
					// stay in bounds, skip visited, skip walls (-1)
					if nr < 0 || nr >= m || nc < 0 || nc >= n {
						continue
					}
					if visited[nr][nc] || rooms[nr][nc] == -1 {
						continue
					}
					visited[nr][nc] = true
					queue = append(queue, cell{nr, nc, cur.d + 1})
				}
			}
			// if no gate found, the room stays INF (untouched).
		}
	}
}

// ── Approach 2: Multi-Source BFS From All Gates (Optimal) ────────────────────
//
// multiSourceBFS solves Walls and Gates by seeding a single BFS with ALL gates
// at once, so distances flow outward simultaneously.
//
// Intuition:
//
//	Instead of asking each room "where is my nearest gate?", flip it: push
//	water out from every gate at the same time. The first time the wave
//	touches a room, it arrives via the shortest path from the closest gate —
//	because all gates start at distance 0 and BFS expands in lockstep by
//	distance. Each room is filled exactly once, so the whole grid is solved
//	in one sweep.
//
// Algorithm:
//  1. Enqueue every gate (value 0) as a BFS source.
//  2. Pop a cell; for each empty (== INF) neighbour, set its value to
//     current + 1 and enqueue it. Setting the value marks it visited.
//  3. Continue until the queue drains.
//
// Time:  O(m*n) — every cell is enqueued/dequeued at most once.
// Space: O(m*n) — the BFS queue in the worst case.
func multiSourceBFS(rooms [][]int) {
	if len(rooms) == 0 || len(rooms[0]) == 0 {
		return
	}
	m, n := len(rooms), len(rooms[0])
	dirs := [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

	type point struct{ r, c int }
	queue := []point{}
	// Seed the queue with all gates. Their stored value (0) IS their distance.
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if rooms[i][j] == 0 {
				queue = append(queue, point{i, j})
			}
		}
	}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, dir := range dirs {
			nr, nc := cur.r+dir[0], cur.c+dir[1]
			if nr < 0 || nr >= m || nc < 0 || nc >= n {
				continue
			}
			// Only empty rooms (INF) are unvisited; walls and already-filled
			// rooms are skipped automatically because they are not INF.
			if rooms[nr][nc] != INF {
				continue
			}
			// Nearest-gate distance = neighbour's own distance + 1.
			rooms[nr][nc] = rooms[cur.r][cur.c] + 1
			queue = append(queue, point{nr, nc}) // expand the wave
		}
	}
}

// deepCopy clones a grid so each approach runs on a fresh copy.
func deepCopy(g [][]int) [][]int {
	c := make([][]int, len(g))
	for i := range g {
		c[i] = make([]int, len(g[i]))
		copy(c[i], g[i])
	}
	return c
}

func main() {
	// Example 1 grid. INF = 2147483647.
	// Input:
	//  [INF,  -1,   0, INF]
	//  [INF, INF, INF,  -1]
	//  [INF,  -1, INF,  -1]
	//  [  0,  -1, INF, INF]
	// Expected output:
	//  [ 3,  -1,  0,  1]
	//  [ 2,   2,  1, -1]
	//  [ 1,  -1,  2, -1]
	//  [ 0,  -1,  3,  4]
	example1 := [][]int{
		{INF, -1, 0, INF},
		{INF, INF, INF, -1},
		{INF, -1, INF, -1},
		{0, -1, INF, INF},
	}

	// Example 2: single gate.
	// Input:  [[0]]      Expected: [[0]]
	example2 := [][]int{{0}}

	fmt.Println("=== Approach 1: Brute Force BFS (per empty room) ===")
	g1 := deepCopy(example1)
	bruteForceBFS(g1)
	for _, row := range g1 {
		fmt.Println(row) // expected [3 -1 0 1] [2 2 1 -1] [1 -1 2 -1] [0 -1 3 4]
	}
	g2 := deepCopy(example2)
	bruteForceBFS(g2)
	fmt.Println(g2) // expected [[0]]

	fmt.Println("=== Approach 2: Multi-Source BFS (Optimal) ===")
	h1 := deepCopy(example1)
	multiSourceBFS(h1)
	for _, row := range h1 {
		fmt.Println(row) // expected [3 -1 0 1] [2 2 1 -1] [1 -1 2 -1] [0 -1 3 4]
	}
	h2 := deepCopy(example2)
	multiSourceBFS(h2)
	fmt.Println(h2) // expected [[0]]
}
