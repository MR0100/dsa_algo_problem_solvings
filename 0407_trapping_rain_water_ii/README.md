# 0407 — Trapping Rain Water II

> LeetCode #407 · Difficulty: Hard
> **Categories:** Array, Breadth-First Search, Heap (Priority Queue), Matrix

---

## Problem Statement

Given an `m x n` integer matrix `heightMap` representing the height of each unit cell in a 2D elevation map, return *the volume of water it can trap after raining*.

**Example 1:**

```
Input: heightMap = [[1,4,3,1,3,2],[3,2,1,3,2,4],[2,3,3,2,3,1]]
Output: 4
Explanation: After the rain, water is trapped between the blocks.
We have two small ponds 1 and 3 units trapped.
The total volume of water trapped is 4.
```

**Example 2:**

```
Input: heightMap = [[3,3,3,3,3],[3,2,2,2,3],[3,2,1,2,3],[3,2,2,2,3],[3,3,3,3,3]]
Output: 10
```

**Constraints:**

- `m == heightMap.length`
- `n == heightMap[i].length`
- `1 <= m, n <= 200`
- `0 <= heightMap[i][j] <= 2 * 10^4`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Uber       | ★★☆☆☆ Low        | 2022          |
| TikTok     | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Heap / Priority Queue** — the optimal method always pours over the *lowest wall on the current frontier*; a min-heap keyed on cell height delivers that minimum in O(log(mn)) → see [`/dsa/heap_priority_queue.md`](/dsa/heap_priority_queue.md)
- **BFS from the boundary (Dijkstra-like)** — we grow an inward frontier from the border, settling each cell the first time it is reached through its minimal enclosing rim — Dijkstra's shortest-bottleneck logic on a grid → see [`/dsa/dijkstra.md`](/dsa/dijkstra.md)
- **Matrix traversal** — 4-directional neighbour expansion with border seeding on an `m x n` grid → see [`/dsa/matrix_traversal.md`](/dsa/matrix_traversal.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Boundary Min-Heap (Priority-Queue BFS) | O(mn·log(mn)) | O(mn) | The optimal, expected answer; one heap pass over the grid |
| 2 | Iterative Relaxation (Fixed Point) | O((mn)²) | O(mn) | Easiest correctness proof; too slow for large grids but great for intuition |

Unlike the 1D version (#42), you **cannot** just take `min(leftMax, rightMax)` per column — water in 2D can leak diagonally around obstacles, so the bounding wall is the lowest wall on the *shortest escape path*, which needs the heap.

---

## Approach 1 — Boundary Min-Heap (Priority-Queue BFS)

### Intuition

Think of the grid as a landscape and imagine flooding it from the outside. Border cells leak instantly, so they trap nothing and form the initial **rim**. The water level any interior cell can hold equals the height of the **lowest wall on the cheapest path from that cell out to the border** — the classic "a basin fills only as high as its lowest rim point" fact.

To discover those minimal rims efficiently, expand from the border **through the lowest available wall each time** (a min-heap gives it). When we cross a rim of height `boundary` into a neighbour of floor height `h`:

- If `h < boundary`, that neighbour is a dip below the rim → it traps `boundary − h` units, and its water surface rises to `boundary`.
- If `h ≥ boundary`, it traps nothing and becomes a new, higher rim.

Because we always process the lowest frontier wall, the first time we touch a cell we have already found the minimal rim enclosing it — the same optimality argument as Dijkstra, applied to "minimise the maximum wall along a path".

### Algorithm

1. Push every **border** cell into a min-heap keyed on height; mark them visited.
2. Pop the lowest wall `boundary`. For each unvisited 4-neighbour:
   - add `max(0, boundary − heightMap[nbr])` to the answer;
   - mark it visited;
   - push it with height `max(boundary, heightMap[nbr])` (water raises its floor).
3. Repeat until the heap is empty; return the accumulated water.

### Complexity

- **Time:** O(mn·log(mn)) — every cell is pushed and popped exactly once, each heap operation is `O(log(mn))`.
- **Space:** O(mn) — the `visited` grid plus the heap (which holds up to all cells).

### Code

```go
func priorityQueueBFS(heightMap [][]int) int {
	m := len(heightMap)
	if m == 0 {
		return 0
	}
	n := len(heightMap[0])
	// A grid smaller than 3x3 has no interior cell that can be enclosed.
	if m < 3 || n < 3 {
		return 0
	}

	visited := make([][]bool, m) // whether a cell's final water level is settled
	for i := range visited {
		visited[i] = make([]bool, n)
	}

	h := &minHeap{}
	heap.Init(h)

	// Seed the heap with the entire border — border cells can never trap water.
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if i == 0 || i == m-1 || j == 0 || j == n-1 {
				heap.Push(h, cell{height: heightMap[i][j], row: i, col: j})
				visited[i][j] = true // border is our starting rim
			}
		}
	}

	dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} // up, down, left, right
	water := 0

	for h.Len() > 0 {
		// The lowest wall on the current frontier defines the rim we pour over.
		top := heap.Pop(h).(cell)
		boundary := top.height
		for _, d := range dirs {
			nr, nc := top.row+d[0], top.col+d[1]
			// Skip out-of-bounds and already-settled cells.
			if nr < 0 || nr >= m || nc < 0 || nc >= n || visited[nr][nc] {
				continue
			}
			visited[nr][nc] = true // this cell is now settled by the current rim
			// If the neighbour sits below the rim, it fills up to the rim height.
			if heightMap[nr][nc] < boundary {
				water += boundary - heightMap[nr][nc]
			}
			// The neighbour's effective wall = max(its own height, the rim). Water
			// resting on it raises the floor future cells must clear.
			newHeight := heightMap[nr][nc]
			if boundary > newHeight {
				newHeight = boundary
			}
			heap.Push(h, cell{height: newHeight, row: nr, col: nc})
		}
	}
	return water
}
```

### Dry Run

Example 1: `heightMap = [[1,4,3,1,3,2],[3,2,1,3,2,4],[2,3,3,2,3,1]]` (3 rows × 6 cols).

Only two cells are interior: `(1,1)=2`, `(1,2)=1`, `(1,3)=3`, `(1,4)=2` (row 1, cols 1..4).

The heap is seeded with all 14 border cells. We pop lowest-first. The relevant settling events:

| Event | Popped rim (h) | Neighbour settled | Floor | Water added | Notes |
|-------|----------------|-------------------|-------|-------------|-------|
| … | border pops in ascending height | — | — | 0 | border traps nothing |
| A | rim = 3 (from `(1,0)=3` / `(0,1)=4` region) | `(1,1)` floor 2 | 2 | `3 − 2 = 1` | dip below rim ⇒ +1 |
| B | rim = 2 (from `(1,4)`/`(2,4)=3`, `(0,4)=3`) | `(1,4)` floor 2 | 2 | `0` | equal to rim ⇒ +0 |
| C | rim = 3 (`(1,3)=3` acts as wall) | `(1,2)` floor 1 | 1 | `3 − 1 = 2`… capped by true rim | enclosed pond |

Summing the trapped contributions over the interior cells gives the two ponds (1 unit and 3 units) → total **4** ✔

(The precise pop order depends on heap tie-breaking, but the *first* time each interior cell is settled it is through its minimal enclosing rim, so the total is invariant.)

---

## Approach 2 — Iterative Relaxation (Fixed Point)

### Intuition

Let `level[i][j]` be the final water-surface height at each cell. Border cells are pinned to their own floor (they leak). For an interior cell, water can stand at most as high as its **lowest neighbour's surface**, but never below its own floor:

```
level[i][j] = max( heightMap[i][j], min over 4 neighbours of level[nbr] )
```

Start optimistically — every interior surface at the global maximum height — then sweep the grid repeatedly applying that rule, lowering surfaces until a full pass changes nothing. Trapped water at a cell is `level − floor`. It is a label-correcting / Bellman-Ford-style fixed point: slower than the heap but very easy to believe.

### Algorithm

1. Set `level = height` on the border; set interior `level = maxHeight`.
2. Repeat sweeps: for each interior cell, `newLevel = max(floor, min neighbour level)`; if it lowers the stored `level`, record a change.
3. Stop when a full sweep makes no change.
4. Sum `level[i][j] − heightMap[i][j]` over interior cells.

### Complexity

- **Time:** O((mn)²) worst case — each sweep is `O(mn)` and up to `O(mn)` sweeps may be needed for a change to propagate all the way across.
- **Space:** O(mn) — the `level` grid.

### Code

```go
func iterativeRelaxation(heightMap [][]int) int {
	m := len(heightMap)
	if m < 3 {
		return 0
	}
	n := len(heightMap[0])
	if n < 3 {
		return 0
	}

	// Find the tallest wall — no water surface can exceed it.
	maxHeight := 0
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if heightMap[i][j] > maxHeight {
				maxHeight = heightMap[i][j]
			}
		}
	}

	level := make([][]int, m) // final water surface height per cell
	for i := 0; i < m; i++ {
		level[i] = make([]int, n)
		for j := 0; j < n; j++ {
			if i == 0 || i == m-1 || j == 0 || j == n-1 {
				level[i][j] = heightMap[i][j] // border pinned to its own height
			} else {
				level[i][j] = maxHeight // interior starts optimistically high
			}
		}
	}

	dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	for {
		changed := false
		for i := 1; i < m-1; i++ {
			for j := 1; j < n-1; j++ {
				// The lowest neighbouring surface caps how high water can stand here.
				minNbr := maxHeight
				for _, d := range dirs {
					if v := level[i+d[0]][j+d[1]]; v < minNbr {
						minNbr = v
					}
				}
				// Surface can't be below the floor.
				newLevel := minNbr
				if heightMap[i][j] > newLevel {
					newLevel = heightMap[i][j]
				}
				if newLevel < level[i][j] {
					level[i][j] = newLevel // relax downward
					changed = true
				}
			}
		}
		if !changed {
			break // fixed point reached
		}
	}

	water := 0
	for i := 1; i < m-1; i++ {
		for j := 1; j < n-1; j++ {
			water += level[i][j] - heightMap[i][j] // surface minus floor = trapped
		}
	}
	return water
}
```

### Dry Run

Example 1: `heightMap = [[1,4,3,1,3,2],[3,2,1,3,2,4],[2,3,3,2,3,1]]`, `maxHeight = 4`.

Interior cells are row 1, cols 1..4: floors `2, 1, 3, 2`. Initial interior `level = 4` everywhere.

| Sweep | Cell | min neighbour level | max(floor, minNbr) | new level | changed? |
|-------|------|---------------------|--------------------|-----------|----------|
| 1 | (1,1) f=2 | neighbours: (0,1)=4,(2,1)=3,(1,0)=3,(1,2)=4 → min 3 | max(2,3)=3 | 4→3 | yes |
| 1 | (1,2) f=1 | (0,2)=3,(2,2)=3,(1,1)=3,(1,3)=4 → min 3 | max(1,3)=3 | 4→3 | yes |
| 1 | (1,3) f=3 | (0,3)=1,(2,3)=2,(1,2)=3,(1,4)=4 → min 1 | max(3,1)=3 | 4→3 | yes |
| 1 | (1,4) f=2 | (0,4)=3,(2,4)=3,(1,3)=3,(1,5)=4 → min 3 | max(2,3)=3 | 4→3 | yes |
| 2 | (1,2) f=1 | now (1,1)=3,(1,3)=3,(0,2)=3,(2,2)=3 → min 3 | max(1,3)=3 | 3 | no |
| 2 | others | … | … | stable | no |

Sweep 2 changes nothing → fixed point. Trapped = Σ(level − floor) over interior = `(3−2)+(3−1)+(3−3)+(3−2) = 1+2+0+1 = 4` ✔

---

## Key Takeaways

- **2D water ≠ 1D per-column.** In one dimension each bar's water is `min(leftMax, rightMax) − height`. In two dimensions water can escape *around* obstacles, so the binding wall is the lowest point on the cheapest escape path — you must search, not just take a per-axis min.
- **Process the frontier lowest-first.** A min-heap seeded with the border turns "flood inward" into a greedy that settles each cell through its minimal enclosing rim on first touch — Dijkstra's optimality applied to `minimise(max wall on path)`.
- **Water raises the floor.** After a dip fills to the rim, push it back with height `max(rim, floor)`; downstream cells must clear the *water surface*, not the original ground.
- **The relaxation view is a great sanity check.** Framing `level[i][j] = max(floor, min neighbour level)` as a fixed point makes correctness obvious and is a handy fallback when you can't recall the heap mechanics under pressure.

---

## Related Problems

- LeetCode #42 — Trapping Rain Water (the 1D original; two-pointer / prefix-max)
- LeetCode #778 — Swim in Rising Water (min-heap on max-height-along-path, same rim idea)
- LeetCode #1631 — Path With Minimum Effort (Dijkstra minimising the max edge)
- LeetCode #417 — Pacific Atlantic Water Flow (multi-source BFS from borders)
