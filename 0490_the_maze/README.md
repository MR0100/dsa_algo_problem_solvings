# 0490 — The Maze

> LeetCode #490 · Difficulty: Medium
> **Categories:** Array, Depth-First Search, Breadth-First Search, Graph, Matrix

---

## Problem Statement

There is a ball in a `maze` with empty spaces (represented as `0`) and walls (represented as `1`). The ball can go through the empty spaces by rolling **up, down, left or right**, but it won't stop rolling until hitting a wall. When the ball stops, it could choose the next direction.

Given the `m x n` `maze`, the ball's `start` position and the `destination`, where `start = [startrow, startcol]` and `destination = [destinationrow, destinationcol]`, return `true` if the ball can stop at the destination, otherwise return `false`.

You may assume that **the borders of the maze are all walls** (see examples).

**Example 1:**

```
Input: maze = [[0,0,1,0,0],[0,0,0,0,0],[0,0,0,1,0],[1,1,0,1,1],[0,0,0,0,0]], start = [0,4], destination = [4,4]
Output: true
Explanation: One possible way is : left -> down -> left -> down -> right -> down -> right.
```

**Example 2:**

```
Input: maze = [[0,0,1,0,0],[0,0,0,0,0],[0,0,0,1,0],[1,1,0,1,1],[0,0,0,0,0]], start = [0,4], destination = [3,2]
Output: false
Explanation: There is no way for the ball to stop at the destination. Notice that you can pass through the destination but you cannot stop there.
```

**Example 3:**

```
Input: maze = [[0,0,0,0,0],[1,1,0,0,1],[0,0,0,0,0],[0,1,0,0,1],[0,1,0,0,0]], start = [4,3], destination = [0,1]
Output: false
```

**Constraints:**

- `m == maze.length`
- `n == maze[i].length`
- `1 <= m, n <= 100`
- `maze[i][j]` is `0` or `1`.
- `start.length == 2`
- `destination.length == 2`
- `0 <= startrow, destinationrow < m`
- `0 <= startcol, destinationcol < n`
- Both the ball and the destination exist in an empty space, and they will not be in the same position initially.
- The maze contains **at least 2 empty spaces**.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Facebook   | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **BFS / DFS on an implicit graph** — nodes are *stopping cells*, and an edge links two stops if the ball can roll straight from one into a wall at the other; reachability is a plain graph search over that graph → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Matrix traversal with direction vectors** — rolling uses the four `(dRow, dCol)` deltas, sliding until the next cell is a wall or off-grid → see [`/dsa/matrix_traversal.md`](/dsa/matrix_traversal.md)
- **Queue (BFS frontier)** — the level-order search maintains a FIFO of stop cells to expand → see [`/dsa/queue_deque.md`](/dsa/queue_deque.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | BFS over stopping cells | O(m·n·max(m,n)) | O(m·n) | Natural framing; explores stops level by level |
| 2 | DFS over stopping cells | O(m·n·max(m,n)) | O(m·n) | Same graph, recursive; shortest-path not needed so DFS is fine |
| 3 | BFS with Explicit Stop-Set (Optimal) | O(m·n·max(m,n)) | O(m·n) | Cleanest form; early-returns the instant the destination is a stop |

> The three share one idea and one big-O; they differ only in traversal order
> and where the destination check sits. The **key** insight — "a move rolls all
> the way to a wall" — is identical in all three.

---

## Approach 1 — BFS over stopping cells

### Intuition

The ball can't stop mid-roll, so the only meaningful positions are cells where it **comes to rest** against a wall. Model each such stop cell as a graph node; rolling in a direction from a node lands on exactly one neighbour node (the wall it slams into). Then "can the ball stop at `destination`?" is just "is `destination` reachable in this graph?" — a textbook BFS from `start`, using a `visited` grid over stop cells to avoid cycles.

### Algorithm

1. Seed a queue with `start` and mark it visited.
2. Pop `(r, c)`. If it equals `destination`, return `true`.
3. For each of the four directions, `roll` to the stop cell `(nr, nc)`; if unvisited, mark it and enqueue.
4. If the queue empties without reaching `destination`, return `false`.

### Complexity

- **Time:** O(m·n·max(m,n)) — each of the `m·n` cells is enqueued at most once, and expanding a cell rolls up to O(max(m,n)) steps in each of four directions.
- **Space:** O(m·n) — the `visited` grid and the BFS queue.

### Code

```go
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
```

The shared `roll` helper (slide until the next cell is a wall/edge):

```go
func roll(maze [][]int, r, c, dr, dc int) (int, int) {
	m, n := len(maze), len(maze[0])
	for r+dr >= 0 && r+dr < m && c+dc >= 0 && c+dc < n && maze[r+dr][c+dc] == 0 {
		r += dr // step one cell further in the roll direction
		c += dc
	}
	return r, c // first cell against a wall/border in this direction
}
```

### Dry Run

Example 1: `start = [0,4]`, `destination = [4,4]`. Directions rolled: up, down, left, right.

| Pop (r,c) | == dest? | roll up | roll down | roll left | roll right | newly enqueued |
|-----------|----------|---------|-----------|-----------|------------|----------------|
| (0,4) | no | (0,4) self | (2,4) | (0,3) | (0,4) self | (2,4), (0,3) |
| (2,4) | no | (0,4)✓seen | (4,4) | (2,4) self* | (2,4) self | (4,4) |
| (0,3) | no | (0,3) self | (0,3) self† | (0,3) self | (0,4)✓seen | — |
| (4,4) | **yes** | — | — | — | — | **return true** |

\* rolling left from (2,4) is blocked by the wall at (2,3), so it stays. † (0,3) is boxed by walls at (0,2) and below/around within this trace. BFS dequeues (4,4) and it equals the destination → **true** ✔.

For Example 2 (`destination = [3,2]`), (3,2) is never produced as a *stop* cell — the ball only ever passes through it while rolling — so BFS drains the queue and returns **false**.

---

## Approach 2 — DFS over stopping cells

### Intuition

Reachability ignores path length, so depth-first search solves the same stop-cell graph just as correctly as BFS. From the current stop cell, roll in all four directions and recurse into each newly discovered stop cell, returning `true` immediately if any recursion finds the destination. The `visited` grid again prevents infinite loops on the cyclic roll-graph.

### Algorithm

1. `visit(r, c)`: if `(r, c)` is the destination, return `true`; otherwise mark it visited.
2. For each direction, `roll` to `(nr, nc)`; if unvisited, recurse `visit(nr, nc)` and propagate a `true`.
3. Return `false` if no direction leads to the destination.

### Complexity

- **Time:** O(m·n·max(m,n)) — identical node/edge bound to BFS.
- **Space:** O(m·n) — `visited` grid plus recursion stack (depth up to `m·n`).

### Code

```go
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
```

### Dry Run

Example 1: `start = [0,4]`, `destination = [4,4]`. DFS explores directions in order up, down, left, right.

| Depth | visit(r,c) | == dest? | first productive roll | descends to |
|-------|------------|----------|-----------------------|-------------|
| 0 | (0,4) | no | up → self (skip), down → (2,4) | (2,4) |
| 1 | (2,4) | no | up → (0,4) seen; down → (4,4) | (4,4) |
| 2 | (4,4) | **yes** | — | returns **true** up the stack |

Every frame propagates the `true`, so `dfs` returns **true** ✔. Example 2 recurses through all reachable stops, none equal `(3,2)`, and returns **false**.

---

## Approach 3 — BFS with Explicit Stop-Set (Optimal)

### Intuition

Functionally the same BFS as Approach 1, refactored so the intent is explicit: a cell is marked visited exactly when it is *discovered* as a stop, and the search early-returns the instant the destination first appears as a stop cell — the earliest possible moment, avoiding an extra dequeue-time comparison. A quick guard handles the degenerate `start == destination` case up front.

### Algorithm

1. If `start == destination`, return `true`. Seed the queue and `visited` with `start`.
2. Pop a cell; roll in all four directions. For each fresh stop cell `(nr, nc)`:
   - if it is the destination, return `true`;
   - else mark it visited and enqueue.
3. Empty queue → return `false`.

### Complexity

- **Time:** O(m·n·max(m,n)) — same as the other two.
- **Space:** O(m·n) — `visited` grid and queue.

### Code

```go
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
```

### Dry Run

Example 1: `start = [0,4]`, `destination = [4,4]`.

| Pop (r,c) | roll results (fresh stops) | destination hit? | queue after |
|-----------|----------------------------|------------------|-------------|
| (0,4) | down→(2,4), left→(0,3) | no | [(2,4), (0,3)] |
| (2,4) | down→(4,4) | **(4,4)?** yes → **return true** | — |

The destination is caught the moment `(4,4)` is generated from `(2,4)`, so the function returns **true** ✔ without enqueuing it. Example 2 never generates `(3,2)` as a stop and returns **false**.

---

## Key Takeaways

- **Redefine "a move" to match the physics.** The ball rolls until a wall, so an edge is a full slide, not a single step. Collapsing each roll into one neighbour turns a weird movement rule into an ordinary graph-search problem.
- **Reachability ⇒ BFS or DFS both work.** The Maze asks only *whether* the destination is a stop, so any traversal suffices; you'd reach for BFS/Dijkstra only in the follow-ups (#505 The Maze II asks for the shortest roll distance, #499 for the lexicographically smallest path).
- **The destination must be a STOP, not merely on the path.** Example 2's `(3,2)` is passed through but never a resting cell — the reason the answer is `false`. Always test the destination against stop cells, never against rolled-through cells.
- **Direction vectors + a shared `roll` helper** keep all three solutions to a few lines and eliminate copy-paste bugs in the boundary/wall checks.

---

## Related Problems

- LeetCode #505 — The Maze II (shortest rolling distance to the destination; BFS/Dijkstra with costs)
- LeetCode #499 — The Maze III (roll to a hole, lexicographically smallest path)
- LeetCode #200 — Number of Islands (grid BFS/DFS/flood fill)
- LeetCode #1197 — Minimum Knight Moves (implicit-graph BFS with custom moves)
- LeetCode #286 — Walls and Gates (multi-source BFS on a grid)
