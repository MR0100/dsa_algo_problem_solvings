# 0286 — Walls and Gates

> LeetCode #286 · Difficulty: Medium
> **Categories:** Breadth-First Search, Matrix, Graph

---

## Problem Statement

You are given an `m x n` grid `rooms` initialized with these three possible values:

- `-1` — A wall or an obstacle.
- `0` — A gate.
- `2147483647` — Infinity (`INF`), meaning an empty room. `2^31 - 1 = 2147483647` is used to represent `INF` since it is larger than any practical distance to a gate.

Fill each empty room with the distance to its **nearest** gate. If it is impossible to reach a gate, that room should remain filled with `INF`.

**Example 1:**

```
Input: rooms = [[2147483647,-1,0,2147483647],
                [2147483647,2147483647,2147483647,-1],
                [2147483647,-1,2147483647,-1],
                [0,-1,2147483647,2147483647]]
Output: [[3,-1,0,1],
         [2,2,1,-1],
         [1,-1,2,-1],
         [0,-1,3,4]]
```

**Example 2:**

```
Input: rooms = [[0]]
Output: [[0]]
```

**Constraints:**

- `m == rooms.length`
- `n == rooms[i].length`
- `1 <= m, n <= 250`
- `rooms[i][j]` is `-1`, `0`, or `2^31 - 1`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Meta       | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Uber       | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Multi-Source BFS** — seeding one BFS with every gate at once floods shortest-distance labels across the grid in a single sweep → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Matrix Traversal** — the grid is an implicit graph where each cell connects to its 4 orthogonal neighbours → see [`/dsa/matrix_traversal.md`](/dsa/matrix_traversal.md)
- **Queue / Deque** — BFS frontier expansion uses a FIFO queue → see [`/dsa/queue_deque.md`](/dsa/queue_deque.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force BFS per room | O((m·n)²) | O(m·n) | Baseline; independent BFS from each empty room, no shared work |
| 2 | Multi-Source BFS (Optimal) | O(m·n) | O(m·n) | The standard answer; one wave from all gates fills the whole grid |

---

## Approach 1 — Brute Force BFS From Every Empty Room

### Intuition

Distance to the *nearest* gate is a shortest-path question, and BFS answers shortest-path on an unweighted grid. So from each empty room, run a BFS outward; the first gate the wave touches is the closest one, and its BFS depth is the answer for that room. This is correct but repeats enormous amounts of work — nothing is shared between rooms.

### Algorithm

1. For every cell whose value is `INF` (an empty room):
   1. Run BFS from that cell, tracking distance, stepping onto in-bounds non-wall cells only.
   2. Stop at the first gate (`0`) reached and write its distance into the origin room.
2. Rooms that never reach a gate keep `INF`.

### Complexity

- **Time:** O((m·n)²) — up to `m·n` empty rooms, each launching a BFS that can visit O(m·n) cells.
- **Space:** O(m·n) — the visited grid and queue for one BFS at a time.

### Code

```go
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
		}
	}
}
```

### Dry Run

Example 1, focusing on the room at `(0,0)` (value `INF`). Nearest gate is `(0,2)`, but a wall at `(0,1)` blocks the straight path, so BFS routes down and around.

| BFS depth | Frontier cells reached from (0,0) | Gate found? |
|-----------|-----------------------------------|-------------|
| 0 | (0,0) | no |
| 1 | (1,0) | no (down; right is wall) |
| 2 | (2,0),(1,1) | no |
| 3 | (1,2) → then its neighbour (0,2) is a gate at depth… | detected at depth 3 |

BFS reaches gate `(0,2)` at distance `3`, so `rooms[0][0] = 3` ✔. The process repeats independently for every other `INF` room.

---

## Approach 2 — Multi-Source BFS From All Gates (Optimal)

### Intuition

Flip the question. Instead of each room searching for its gate, let **every gate push distance outward at the same time**. Seed a single BFS queue with all gates (each at distance `0`). Because BFS expands strictly in order of increasing distance and all sources start together, the first time the wave reaches a room it arrives from the *closest* gate along the shortest path. Each room is written exactly once, so the whole board is solved in one linear sweep.

### Algorithm

1. Enqueue every gate (cell value `0`) as a BFS source.
2. Pop a cell; for each neighbour that is still `INF`, set it to `current + 1` and enqueue it (writing the value also marks it visited).
3. Repeat until the queue empties. Walls (`-1`) and already-filled rooms are never `INF`, so they are skipped automatically.

### Complexity

- **Time:** O(m·n) — every cell is enqueued and dequeued at most once.
- **Space:** O(m·n) — the BFS queue in the worst case.

### Code

```go
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
			if rooms[nr][nc] != INF {
				continue
			}
			rooms[nr][nc] = rooms[cur.r][cur.c] + 1
			queue = append(queue, point{nr, nc}) // expand the wave
		}
	}
}
```

### Dry Run

Example 1 has two gates: `(0,2)` and `(3,0)`. Queue initialized with both.

| Round | Popped (dist) | Neighbours set to dist+1 |
|-------|---------------|--------------------------|
| init | queue = [(0,2)=0, (3,0)=0] | — |
| 1 | (0,2)=0 | (1,2)=1, (0,3)=1 |
| 2 | (3,0)=0 | (2,0)=1 |
| 3 | (1,2)=1 | (1,1)=2 |
| 4 | (0,3)=1 | — (neighbours are wall / gate) |
| 5 | (2,0)=1 | (1,0)=2 |
| 6 | (1,1)=2 | (1,0) already 2 → skip |
| 7 | (1,0)=2 | (0,0)=3 |
| … | … | eventually (2,2)=2, (3,2)=3, (3,3)=4 |

Final board matches:
```
[3,-1,0,1]
[2,2,1,-1]
[1,-1,2,-1]
[0,-1,3,4]
``` ✔

---

## Key Takeaways

- **Multi-source BFS turns "nearest of many targets" into one sweep.** Whenever you'd otherwise BFS from each cell to find the closest of several sources, seed the queue with *all* sources at distance 0 instead — you drop a factor of O(m·n).
- **Mutating the grid can double as the visited marker.** Here `INF` means "unvisited empty room"; writing any finite distance both records the answer and prevents re-visiting, saving a separate `visited` array.
- BFS on an unweighted grid gives shortest distances; the first time a cell is dequeued/labelled is optimal, so no relaxation (Dijkstra) is needed.
- Order of processing the two gates does not matter — the lockstep-by-distance property guarantees the minimum wins.

---

## Related Problems

- LeetCode #542 — 01 Matrix (multi-source BFS from all zeros)
- LeetCode #994 — Rotting Oranges (multi-source BFS wave, time steps)
- LeetCode #1162 — As Far from Land as Possible (multi-source BFS, maximize distance)
- LeetCode #200 — Number of Islands (grid flood fill)
- LeetCode #317 — Shortest Distance from All Buildings (BFS from each source, accumulate)
