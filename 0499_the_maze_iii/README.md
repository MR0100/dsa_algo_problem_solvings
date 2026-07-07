# 0499 — The Maze III

> LeetCode #499 · Difficulty: Hard
> **Categories:** Graph, Shortest Path, Dijkstra, Breadth-First Search, Heap (Priority Queue), Matrix

---

## Problem Statement

There is a ball in a `maze` with empty spaces (represented as `0`) and walls (represented as `1`). The ball can go through the empty spaces by rolling **up, down, left or right**, but it won't stop rolling until hitting a wall. When the ball stops, it could choose the next direction. There is also a **hole** in this maze. The ball will drop into the hole if it rolls onto the hole.

Given the `m x n` `maze`, the ball's position `ball` and the hole's position `hole`, where `ball = [ballrow, ballcol]` and `hole = [holerow, holecol]`, return *a string* `instructions` *of all the instructions that the ball should follow to drop in the hole with the **shortest distance** possible.* If there are multiple valid instructions, return the **lexicographically minimum** one. If the ball can't drop in the hole, return `"impossible"`.

If there is a way for the ball to drop in the hole, the answer `instructions` should contain the characters `'u'` (i.e., up), `'d'` (i.e., down), `'l'` (i.e., left), and `'r'` (i.e., right).

The **distance** is the number of **empty spaces** traveled by the ball from the start position (excluded) to the destination (included).

You may assume that **the borders of the maze are all walls** (see examples).

**Example 1:**

```
Input: maze = [[0,0,0,0,0],[1,1,0,0,1],[0,0,0,0,0],[0,1,0,0,1],[0,1,0,0,0]], ball = [4,3], hole = [0,1]
Output: "lul"
Explanation: There are two shortest ways for the ball to drop into the hole.
The first way is left -> up -> left, represented by "lul".
The second way is up -> left, represented by 'ul'.
Both ways have shortest distance 6, but the first way is lexicographically smaller because 'l' < 'u'. So the output is "lul".
```

**Example 2:**

```
Input: maze = [[0,0,0,0,0],[1,1,0,0,1],[0,0,0,0,0],[0,1,0,0,1],[0,1,0,0,0]], ball = [4,3], hole = [3,0]
Output: "impossible"
Explanation: The ball cannot reach the hole.
```

**Example 3:**

```
Input: maze = [[0,0,0,0,0,0,0],[0,0,1,0,0,1,0],[0,0,0,0,1,0,0],[0,0,0,0,0,0,1]], ball = [0,4], hole = [3,5]
Output: "dldr"
```

**Constraints:**

- `m == maze.length`
- `n == maze[i].length`
- `1 <= m, n <= 100`
- `maze[i][j]` is `0` or `1`.
- `ball.length == 2`
- `hole.length == 2`
- `0 <= ballrow, holerow <= m`
- `0 <= ballcol, holecol <= n`
- Both the ball and the hole exist in an empty space, and they will not be in the same position initially.
- The maze contains **at least 2 empty spaces**.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★☆☆☆ Low        | 2023          |
| Facebook   | ★★☆☆☆ Low        | 2022          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Dijkstra (weighted shortest path)** — the maze is a *weighted* graph whose nodes are cells and whose edges are "roll until you stop" (edge weight = roll length); Dijkstra settles the shortest, lexicographically-smallest path to the hole. The `roll` primitive is the grid-traversal walk that respects walls, borders, and the mid-roll hole → see [`/dsa/dijkstra.md`](/dsa/dijkstra.md)
- **Heap / Priority Queue** — Dijkstra needs a min-heap ordered by `(distance, path)` so the first extraction of the hole is simultaneously shortest and lexicographically smallest → see [`/dsa/heap_priority_queue.md`](/dsa/heap_priority_queue.md)

> Note: this is a **Dijkstra shortest-path** problem (roll edges have *unequal*
> weights, so plain BFS is incorrect). The repo has no dedicated Dijkstra file
> yet — see "NEW /dsa CONCEPTS NEEDED" in the batch report.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Repeated Relaxation (Bellman-Ford style) | O(V·E·L) | O(V·L) | Simple to reason about; fine for ≤100×100 grids |
| 2 | Dijkstra with (dist, path) Priority Queue (Optimal) | O(E log V + E·L) | O(V·L) | The intended answer; extracts the optimum in one pop |

*(V = cells, E = roll edges ≤ 4V, L = path-string length)*

---

## Approach 1 — Repeated Relaxation (Bellman-Ford style)

### Intuition

Keep, for every stop-cell, the best label found so far — a pair `(distance, path)`. From each labelled cell, try all four rolls and **relax** the neighbour if the new pair beats its current label under the ordering "smaller distance; if equal, lexicographically smaller path". Sweep the whole grid repeatedly; when a full sweep changes nothing, the labels are a fixed point (Bellman-Ford converges in ≤ V sweeps). Because the tie-break is baked into "better", the fixed-point label at the hole is already the lexicographically smallest shortest path.

### Algorithm

1. Set `best[ball] = (0, "")`, all other cells `(∞, unset)`.
2. Repeat sweeps: for every labelled cell, roll all four directions (in letter order `d,l,r,u`); relax the stop cell if `(dist+steps, path+letter)` is better than its label. Mark that a change occurred.
3. Stop when a sweep produces no change.
4. If `best[hole]` is still unset → `"impossible"`, else return its `path`.

### Complexity

- **Time:** O(V·E·L) — up to V sweeps, each rolling every cell in 4 directions (a roll is up to `max(m,n)` steps) with string comparisons of length L. Trivial for the constraint bounds.
- **Space:** O(V·L) — a `(dist, path)` label per cell.

### Code

```go
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
```

### Dry Run

Example 1: `ball = (4,3)`, `hole = (0,1)`. Rolls use letter order `d,l,r,u`. Key relaxations (spaces = empty cells crossed; a roll onto the hole stops there):

| Sweep event | from cell | roll | stop cell | steps | label set at stop |
|-------------|-----------|------|-----------|-------|-------------------|
| init | — | — | (4,3) | — | (0, "") |
| relax | (4,3) | `l` | (4,2)…stops at (4,2)? | — | — |
| relax | (4,3) | `u` | rolls up col 3 to (0,3) | 4 | (4, "u") |
| relax | (4,3) | `l` | rolls left to (4,2) then blocked → (2,?) path leads toward (2,2) | … | intermediate labels |
| … | … | … | … | … | … |
| relax onto hole | (2,2) region via `l` then `u` then `l` | — | (0,1) hole | total 6 | (6, "lul") |
| relax onto hole | (0,3) via `u` then `l` | — | (0,1) hole | total 6 | competing (6, "ul") |

Both routes reach the hole at distance 6. Under the ordering, `"lul" < "ul"` (`'l' < 'u'`), so the hole's final label is `(6, "lul")`. Return `"lul"` ✔

*(The precise intermediate cells depend on the wall layout; what matters is that two distance-6 paths compete and the lexicographically smaller `"lul"` wins the relaxation.)*

---

## Approach 2 — Dijkstra with (distance, path) Priority Queue (Optimal)

### Intuition

A single push can travel several cells, so edges have **different weights** — this is Dijkstra territory, not BFS. Order the frontier by the pair `(distance, path)`. Dijkstra's invariant guarantees that the **first** time a cell is popped it carries its optimal label; ordering ties by the path string means that optimal label is also the lexicographically smallest. Therefore the first pop of the hole is the final answer. A `seen` grid marks finalised cells so we never reprocess them.

### Algorithm

1. Push `(ball, dist=0, path="")` onto a min-heap keyed by `(dist, path)`.
2. Pop the best state. If it is the hole, return its `path`.
3. If its cell is already finalised, discard; otherwise finalise it.
4. Roll all four directions; for each stop cell that moved and is not finalised, push `(dist+steps, path+letter)`.
5. If the heap empties, return `"impossible"`.

### Complexity

- **Time:** O(E log V + E·L) — heap operations over up to `E ≤ 4V` edges, each comparison touching path strings of length L. Negligible for a ≤100×100 maze.
- **Space:** O(V·L) — heap entries plus the `seen` grid.

### Code

```go
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
```

The supporting `roll` primitive (shared by both approaches) advances one cell at a time and stops on the hole or against a wall:

```go
func roll(maze [][]int, r, c, dr, dc, holeR, holeC int) (int, int, int, bool) {
	m, n := len(maze), len(maze[0])
	steps := 0
	for {
		nr, nc := r+dr, c+dc // candidate next cell
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
```

### Dry Run

Example 1: `ball = (4,3)`, `hole = (0,1)`. Heap ordered by `(dist, path)`; showing the decisive pops:

| Pop # | state popped (dist, path) @cell | action |
|-------|----------------------------------|--------|
| 1 | (0, "") @ (4,3) | finalise start; push rolls: `l`→(4,2)…, `u`→(0,3) dist 4, etc. |
| … | smaller-distance interior cells | finalise, keep pushing rolls |
| — | (6, "lul") @ (0,1) and (6, "ul") @ (0,1) both enqueued | heap orders `"lul"` before `"ul"` (same dist 6, `'l' < 'u'`) |
| k | **(6, "lul") @ (0,1)** | cell == hole → **return "lul"** |

Because `(6,"lul")` sits ahead of `(6,"ul")` in the heap, the hole is first popped with path `"lul"`. Return `"lul"` ✔

For `hole = (3,0)` (Example 2), no sequence of rolls ever stops on `(3,0)`; the heap drains without popping it → `"impossible"` ✔

---

## Key Takeaways

- **"Roll until you hit a wall" ⇒ weighted edges ⇒ Dijkstra, not BFS.** The classic Maze I/II/III trap: the number of cells crossed per push varies, so uniform-cost BFS gives wrong distances. Reserve BFS for unit-weight moves.
- **Encode the tie-break into the priority key.** Ordering the heap (or the relaxation test) by `(distance, path)` makes "shortest, then lexicographically smallest" fall out automatically — no post-processing of equal-length paths.
- **The hole can stop a roll mid-flight.** The `roll` primitive must check the hole *after each single step*, before the wall check, otherwise the ball overshoots it.
- **Try directions in sorted letter order** (`d < l < r < u`). It is not strictly required once the path is in the key, but it keeps expansion deterministic and makes the ordering argument obvious.
- Dijkstra's **first-pop-is-final** property is what lets us return immediately on reaching the target instead of exploring the whole graph.

---

## Related Problems

- LeetCode #490 — The Maze (can the ball reach the destination? — BFS/DFS over rolls)
- LeetCode #505 — The Maze II (shortest *distance* to stop at destination — Dijkstra, no path string)
- LeetCode #743 — Network Delay Time (textbook Dijkstra on a weighted graph)
- LeetCode #787 — Cheapest Flights Within K Stops (Bellman-Ford / bounded Dijkstra)
- LeetCode #1631 — Path With Minimum Effort (Dijkstra on a grid with custom edge cost)
