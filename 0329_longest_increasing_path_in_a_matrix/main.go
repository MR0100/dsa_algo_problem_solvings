package main

import "fmt"

// directions holds the four legal moves: up, down, left, right.
// No diagonals, so exactly these four (row-delta, col-delta) pairs.
var directions = [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

// ── Approach 1: Plain DFS (Brute Force) ──────────────────────────────────────
//
// dfsBrute solves Longest Increasing Path in a Matrix by exploring, from every
// cell, every strictly-increasing path with a naive DFS and no memoization.
//
// Intuition:
//
//	The longest increasing path starting at (i,j) is 1 plus the longest
//	increasing path starting at whichever strictly-larger neighbour gives the
//	best continuation. With no memo we simply recompute that recursion from
//	scratch every time we reach a cell, which re-walks overlapping suffixes
//	over and over — correct, but exponential.
//
// Algorithm:
//  1. For each cell (i,j), run dfs(i,j) = longest increasing path starting here.
//  2. dfs explores each of the four neighbours; if a neighbour is strictly
//     larger it recurses and keeps the best 1 + child length.
//  3. The answer is the maximum dfs value over all cells.
//
// Note: because edges only go small→large, no cell repeats within a single
// path, so no visited set is needed — but WITHOUT memo this TLEs on large
// grids (worst case exponential).
//
// Time:  O(2^(m*n)) worst case — overlapping subpaths recomputed repeatedly.
// Space: O(m*n) — recursion stack depth (length of the longest path).
func dfsBrute(matrix [][]int) int {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return 0
	}
	m, n := len(matrix), len(matrix[0])

	// dfs returns the length of the longest increasing path that STARTS at (i,j).
	var dfs func(i, j int) int
	dfs = func(i, j int) int {
		best := 1 // the cell itself is a path of length 1
		for _, d := range directions {
			ni, nj := i+d[0], j+d[1] // neighbour coordinates
			// Must stay in bounds AND be strictly larger to extend the path.
			if ni >= 0 && ni < m && nj >= 0 && nj < n && matrix[ni][nj] > matrix[i][j] {
				length := 1 + dfs(ni, nj) // this cell + best path from neighbour
				if length > best {
					best = length
				}
			}
		}
		return best
	}

	ans := 0
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if v := dfs(i, j); v > ans { // try starting from every cell
				ans = v
			}
		}
	}
	return ans
}

// ── Approach 2: DFS + Memoization (Optimal) ──────────────────────────────────
//
// dfsMemo solves Longest Increasing Path in a Matrix with the same DFS as
// Approach 1 but caches memo[i][j] = longest increasing path starting at (i,j).
//
// Intuition:
//
//	Point an edge from each cell to every strictly-larger neighbour. Because
//	values strictly increase along any path, this graph has NO cycles — it's a
//	DAG. In a DAG the answer for a node depends only on its successors, never on
//	the route taken to reach it, so memo[i][j] is well-defined and safe to
//	cache. Each cell is then solved exactly once.
//
// Algorithm:
//  1. memo[i][j] == 0 means "not computed yet".
//  2. dfs(i,j): if memo[i][j] != 0 return it. Otherwise best = 1, and for each
//     strictly-larger neighbour take 1 + dfs(neighbour), keeping the max.
//  3. Store best in memo[i][j] and return it.
//  4. Answer = max dfs over all cells.
//
// Time:  O(m*n) — each cell computed once; four O(1) neighbour checks each.
// Space: O(m*n) — memo table plus recursion stack.
func dfsMemo(matrix [][]int) int {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return 0
	}
	m, n := len(matrix), len(matrix[0])

	// memo[i][j] = longest increasing path starting at (i,j); 0 = uncomputed.
	memo := make([][]int, m)
	for i := range memo {
		memo[i] = make([]int, n)
	}

	var dfs func(i, j int) int
	dfs = func(i, j int) int {
		if memo[i][j] != 0 {
			return memo[i][j] // already solved this cell — reuse it
		}
		best := 1 // the cell alone
		for _, d := range directions {
			ni, nj := i+d[0], j+d[1]
			if ni >= 0 && ni < m && nj >= 0 && nj < n && matrix[ni][nj] > matrix[i][j] {
				length := 1 + dfs(ni, nj)
				if length > best {
					best = length
				}
			}
		}
		memo[i][j] = best // cache before returning so callers reuse it
		return best
	}

	ans := 0
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if v := dfs(i, j); v > ans {
				ans = v
			}
		}
	}
	return ans
}

// ── Approach 3: Topological Sort (BFS Peeling) ───────────────────────────────
//
// topoSortBFS solves Longest Increasing Path in a Matrix by treating it as a
// DAG and peeling off "peak" cells layer by layer via Kahn's algorithm.
//
// Intuition:
//
//	Orient each edge small→large. A cell's out-degree = how many of its four
//	neighbours are strictly larger. A cell with out-degree 0 is a local "peak":
//	no increasing path can continue past it, so it sits at the END of some
//	longest path. Remove all peaks simultaneously (that's one layer), which may
//	drop other cells' out-degree to 0, exposing the next layer of peaks. The
//	longest increasing path has exactly as many cells as the number of layers
//	we peel, because each layer contributes one step to the deepest chain.
//
// Algorithm:
//  1. Compute outDegree[i][j] = count of strictly-larger neighbours.
//  2. Queue every cell with outDegree 0 (the initial peaks).
//  3. BFS layer by layer: for each cell popped, look at its SMALLER neighbours
//     (predecessors); decrement their out-degree, and when it hits 0 enqueue
//     them for the next layer.
//  4. The number of layers processed is the answer.
//
// Time:  O(m*n) — every cell enqueued once, four neighbour checks each.
// Space: O(m*n) — the out-degree grid and the BFS queue.
func topoSortBFS(matrix [][]int) int {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return 0
	}
	m, n := len(matrix), len(matrix[0])

	// outDegree[i][j] = number of strictly-larger orthogonal neighbours.
	outDegree := make([][]int, m)
	for i := range outDegree {
		outDegree[i] = make([]int, n)
	}

	queue := make([][2]int, 0, m*n) // holds coordinates of current-layer peaks
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			for _, d := range directions {
				ni, nj := i+d[0], j+d[1]
				if ni >= 0 && ni < m && nj >= 0 && nj < n && matrix[ni][nj] > matrix[i][j] {
					outDegree[i][j]++ // (i,j) has an outgoing edge to a bigger neighbour
				}
			}
			if outDegree[i][j] == 0 {
				queue = append(queue, [2]int{i, j}) // a peak: nowhere larger to go
			}
		}
	}

	layers := 0 // how many BFS layers we peel = length of the longest path
	for len(queue) > 0 {
		layers++
		next := make([][2]int, 0) // cells that become peaks after this layer
		for _, cell := range queue {
			i, j := cell[0], cell[1]
			for _, d := range directions {
				ni, nj := i+d[0], j+d[1]
				// Look at SMALLER neighbours — they point INTO this cell.
				if ni >= 0 && ni < m && nj >= 0 && nj < n && matrix[ni][nj] < matrix[i][j] {
					outDegree[ni][nj]-- // remove the edge into the just-peeled cell
					if outDegree[ni][nj] == 0 {
						next = append(next, [2]int{ni, nj}) // newly exposed peak
					}
				}
			}
		}
		queue = next // advance to the next layer
	}
	return layers
}

func main() {
	ex1 := [][]int{{9, 9, 4}, {6, 6, 8}, {2, 1, 1}}
	ex2 := [][]int{{3, 4, 5}, {3, 2, 6}, {2, 2, 1}}
	ex3 := [][]int{{1}}

	fmt.Println("=== Approach 1: Plain DFS (Brute Force) ===")
	fmt.Println(dfsBrute(ex1)) // expected 4
	fmt.Println(dfsBrute(ex2)) // expected 4
	fmt.Println(dfsBrute(ex3)) // expected 1

	fmt.Println("=== Approach 2: DFS + Memoization (Optimal) ===")
	fmt.Println(dfsMemo(ex1)) // expected 4
	fmt.Println(dfsMemo(ex2)) // expected 4
	fmt.Println(dfsMemo(ex3)) // expected 1

	fmt.Println("=== Approach 3: Topological Sort (BFS Peeling) ===")
	fmt.Println(topoSortBFS(ex1)) // expected 4
	fmt.Println(topoSortBFS(ex2)) // expected 4
	fmt.Println(topoSortBFS(ex3)) // expected 1
}
