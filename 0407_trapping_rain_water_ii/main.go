package main

import (
	"container/heap"
	"fmt"
)

// heightMap is an m x n grid of wall heights. After infinite rain, water pools
// wherever it is enclosed by taller walls on all escape routes. Return the total
// trapped volume. This is the 2D generalisation of LeetCode #42.

// ── cell is one grid position tagged with its height, ordered by height in the
//
//	priority queue so the LOWEST wall is always processed first.
type cell struct {
	height int // wall height at this position
	row    int // grid row
	col    int // grid col
}

// minHeap is a container/heap of cells keyed on ascending height. The smallest
// wall on the current water frontier sits at the top.
type minHeap []cell

func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(i, j int) bool  { return h[i].height < h[j].height } // min-heap on height
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(cell)) }
func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	top := old[n-1]
	*h = old[:n-1]
	return top
}

// ── Approach: Boundary Min-Heap (Priority-Queue BFS) — Optimal ───────────────
//
// priorityQueueBFS solves Trapping Rain Water II by growing an inward frontier
// from the border, always crossing the lowest wall on that frontier first.
//
// Intuition:
//
//	Water on the border cell escapes for free, so the border traps nothing. For
//	any interior cell, the water level it can hold is bounded by the LOWEST wall
//	on the shortest "leak path" out to the border — exactly like a topographic
//	basin whose rim height sets the pond level. Process cells from the OUTSIDE
//	in, always through the current lowest rim wall (a min-heap gives that in
//	O(log). When we step from a rim of height `boundary` into a neighbour of
//	height `h`:
//	  - if h < boundary, the neighbour is a dip below the rim → it holds
//	    (boundary - h) units of water, and its effective wall height becomes
//	    `boundary` (water fills it up to the rim);
//	  - if h >= boundary, no water there, its own height is the new local rim.
//	Because we always expand through the lowest available wall, the first time we
//	reach any cell we have found the minimal rim that encloses it — Dijkstra's
//	"shortest bottleneck" guarantee applied to max-height along a path.
//
// Algorithm:
//  1. Push every border cell into a min-heap (keyed on height) and mark visited.
//  2. Repeatedly pop the lowest wall `boundary`. For each unvisited 4-neighbour:
//     add max(0, boundary - neighbourHeight) to the answer, mark it visited, and
//     push it with height max(boundary, neighbourHeight) (water raises the floor).
//  3. When the heap empties, every interior cell has been settled.
//
// Time:  O(m*n log(m*n)) — each cell is pushed/popped once; heap ops are log.
// Space: O(m*n) — the visited grid and the heap.
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

// ── Approach 2: Iterative Relaxation (Bellman-Ford style) — brute-ish ─────────
//
// iterativeRelaxation solves Trapping Rain Water II by initialising every
// interior cell's water level to the global max and repeatedly lowering it until
// no cell can be lowered further (a fixed-point / label-correcting method).
//
// Intuition:
//
//	Define level[i][j] = final water surface height at that cell. Border cells
//	are pinned to their own height (they can't hold water). For an interior cell,
//	its surface can be at most the LOWEST surface among its 4 neighbours, but
//	never below its own floor:
//	    level[i][j] = max(heightMap[i][j], min over neighbours of level[nbr]).
//	Start optimistically with every interior surface at the maximum height, then
//	sweep the grid relaxing this rule until a full pass changes nothing. Trapped
//	water at a cell is level - floor. This is easier to reason about than the
//	heap but slower.
//
// Algorithm:
//  1. Set level[i][j] = heightMap[i][j] on the border, = maxHeight inside.
//  2. Repeat: for every interior cell, newLevel = max(floor, min neighbour
//     level); if it drops the stored level, record the change. Stop when a pass
//     makes no change.
//  3. Sum level - floor over interior cells.
//
// Time:  O((m*n)^2) worst case — up to O(m*n) passes, each O(m*n).
// Space: O(m*n) — the level grid.
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

func main() {
	ex1 := [][]int{
		{1, 4, 3, 1, 3, 2},
		{3, 2, 1, 3, 2, 4},
		{2, 3, 3, 2, 3, 1},
	}
	ex2 := [][]int{
		{3, 3, 3, 3, 3},
		{3, 2, 2, 2, 3},
		{3, 2, 1, 2, 3},
		{3, 2, 2, 2, 3},
		{3, 3, 3, 3, 3},
	}

	fmt.Println("=== Approach 1: Boundary Min-Heap (Priority-Queue BFS) ===")
	fmt.Printf("ex1 got=%d  expected 4\n", priorityQueueBFS(ex1))
	fmt.Printf("ex2 got=%d  expected 10\n", priorityQueueBFS(ex2))

	fmt.Println("=== Approach 2: Iterative Relaxation (Fixed Point) ===")
	fmt.Printf("ex1 got=%d  expected 4\n", iterativeRelaxation(ex1))
	fmt.Printf("ex2 got=%d  expected 10\n", iterativeRelaxation(ex2))
}
