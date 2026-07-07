# 0289 — Game of Life

> LeetCode #289 · Difficulty: Medium
> **Categories:** Array, Matrix, Simulation, Bit Manipulation

---

## Problem Statement

According to Wikipedia's article: "The **Game of Life**, also known simply as **Life**, is a cellular automaton devised by the British mathematician John Horton Conway in 1970."

The board is made up of an `m x n` grid of cells, where each cell has an initial state: **live** (represented by a `1`) or **dead** (represented by a `0`). Each cell interacts with its eight neighbors (horizontal, vertical, diagonal) using the following four rules (taken from the above Wikipedia article):

1. Any live cell with fewer than two live neighbors dies as if caused by under-population.
2. Any live cell with two or three live neighbors lives on to the next generation.
3. Any live cell with more than three live neighbors dies, as if by over-population.
4. Any dead cell with exactly three live neighbors becomes a live cell, as if by reproduction.

The next state of the board is determined by applying the above rules simultaneously to every cell in the current state of the `m x n` grid `board`. In this process, births and deaths occur simultaneously.

Given the current state of the `board`, update the `board` to reflect its next state.

**Note** that you do not need to return anything.

**Example 1:**

```
Input: board = [[0,1,0],[0,0,1],[1,1,1],[0,0,0]]
Output: [[0,0,0],[1,0,1],[0,1,1],[0,1,0]]
```

**Example 2:**

```
Input: board = [[1,1],[1,0]]
Output: [[1,1],[1,1]]
```

**Constraints:**

- `m == board.length`
- `n == board[i].length`
- `1 <= m, n <= 25`
- `board[i][j]` is `0` or `1`.

**Follow up:**

- Could you solve it in-place? Remember that the board needs to be updated simultaneously: You cannot update some cells first and then use their updated values to update other cells.
- In this question, we represent the board using a 2D array. In principle, the board is infinite, which would cause problems when the active area encroaches upon the border of the array (i.e., live cells reach the border). How would you address these problems?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Dropbox    | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Matrix Traversal** — every cell scans its 8 neighbours; the core loop is a bounded 3×3 sweep → see [`/dsa/matrix_traversal.md`](/dsa/matrix_traversal.md)
- **Bit Manipulation** — packing the next state into a spare bit (bit 1) while the current state stays in bit 0 gives O(1) extra space → see [`/dsa/bit_manipulation.md`](/dsa/bit_manipulation.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Extra copy buffer | O(m·n) | O(m·n) | Simplest correct simulation; read old state from a snapshot |
| 2 | In-place 2-bit encoding (Optimal) | O(m·n) | O(1) | The follow-up: update simultaneously with constant extra space |

---

## Approach 1 — Extra Copy Buffer

### Intuition

The rules must apply *simultaneously*: every cell's next value depends only on the *current* generation. The bullet-proof way to always read the old state while writing the new one is to keep a full snapshot copy and count neighbours from it.

### Algorithm

1. Copy `board` into `snap`.
2. For each cell, count its 8 live neighbours in `snap`.
3. Apply Conway's rules and write `0`/`1` into the real `board`.

### Complexity

- **Time:** O(m·n) — a constant 8 neighbour reads per cell.
- **Space:** O(m·n) — the snapshot copy.

### Code

```go
func extraCopy(board [][]int) {
	if len(board) == 0 || len(board[0]) == 0 {
		return
	}
	m, n := len(board), len(board[0])
	// Snapshot of the original board — the source of truth for counting.
	snap := make([][]int, m)
	for i := range board {
		snap[i] = make([]int, n)
		copy(snap[i], board[i])
	}
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			live := countLiveNeighbors(snap, i, j) // count from the snapshot
			if snap[i][j] == 1 {
				if live == 2 || live == 3 {
					board[i][j] = 1
				} else {
					board[i][j] = 0
				}
			} else {
				if live == 3 {
					board[i][j] = 1
				} else {
					board[i][j] = 0
				}
			}
		}
	}
}
```

### Dry Run

Example 1 snapshot:
```
0 1 0
0 0 1
1 1 1
0 0 0
```

A few representative cells (neighbour count from `snap`):

| cell | old | live neighbours | rule | new |
|------|-----|-----------------|------|-----|
| (0,0) | 0 | 1 (just (0,1)) | dead, ≠3 | 0 |
| (1,0) | 0 | 3 ((0,1),(2,0),(2,1)) | dead, ==3 born | 1 |
| (1,2) | 1 | 2 ((2,1),(2,2)) wait count = 3 | live, 2–3 survive | 1 |
| (2,1) | 1 | 3 ((1,2),(2,0),(2,2)) | live, survive | 1 |
| (2,0) | 1 | 1 ((2,1)) | live, <2 dies | 0 |
| (3,1) | 0 | 3 ((2,0),(2,1),(2,2)) | dead, born | 1 |

Full next board:
```
0 0 0
1 0 1
0 1 1
0 1 0
``` ✔

---

## Approach 2 — In-Place 2-Bit State Encoding (Optimal)

### Intuition

We need each cell to keep its **old** value (for neighbours still to be visited) while also recording its **new** value — without a copy. Pack both into two bits: bit 0 holds the current state, bit 1 stashes the next state. While counting, read only bit 0 (`value & 1`), so the original snapshot is preserved implicitly. After the whole sweep, shift every cell right by one bit to promote the next state into bit 0.

Encoding:

| stored | bit1 (new) | bit0 (old) | meaning |
|--------|------------|------------|---------|
| 0 | 0 | 0 | old dead, new dead |
| 1 | 0 | 1 | old live, new dead |
| 2 | 1 | 0 | old dead, new live |
| 3 | 1 | 1 | old live, new live |

### Algorithm

1. For each cell, count live neighbours using `board[nr][nc] & 1` (old bits only).
2. If the cell will be live next generation, set bit 1: `board[i][j] |= 2`.
3. After the full sweep, `board[i][j] >>= 1` everywhere to reveal the next state.

### Complexity

- **Time:** O(m·n) — 8 neighbour reads per cell plus a final shift pass.
- **Space:** O(1) — no extra board; the second bit is free storage.

### Code

```go
func inPlaceBits(board [][]int) {
	if len(board) == 0 || len(board[0]) == 0 {
		return
	}
	m, n := len(board), len(board[0])
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			live := 0
			// Count neighbours reading ONLY bit 0 (the original state).
			for dr := -1; dr <= 1; dr++ {
				for dc := -1; dc <= 1; dc++ {
					if dr == 0 && dc == 0 {
						continue
					}
					nr, nc := i+dr, j+dc
					if nr >= 0 && nr < m && nc >= 0 && nc < n {
						live += board[nr][nc] & 1 // old state lives in bit 0
					}
				}
			}
			old := board[i][j] & 1
			if old == 1 && (live == 2 || live == 3) {
				board[i][j] |= 2 // live stays live
			} else if old == 0 && live == 3 {
				board[i][j] |= 2 // dead becomes live
			}
		}
	}
	// Reveal the next state: shift the stored bit 1 down into bit 0.
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			board[i][j] >>= 1
		}
	}
}
```

### Dry Run

Example 1, first pass sets bit 1 where the cell becomes live. Tracking a few cells (reading neighbours via `&1`, which is unaffected by earlier `|= 2` writes):

| cell | old bit0 | live nbrs (via &1) | next live? | stored after |
|------|----------|--------------------|------------|--------------|
| (1,0) | 0 | 3 | yes | 0 \| 2 = 2 |
| (2,0) | 1 | 1 | no | 1 (bit1 stays 0) |
| (2,1) | 1 | 3 | yes | 1 \| 2 = 3 |
| (3,1) | 0 | 3 | yes | 0 \| 2 = 2 |
| (0,0) | 0 | 1 | no | 0 |

After the sweep, board holds the packed values; the final `>>= 1` pass turns e.g. `2 → 1`, `3 → 1`, `1 → 0`, `0 → 0`, giving:
```
0 0 0
1 0 1
0 1 1
0 1 0
``` ✔

---

## Key Takeaways

- **Simultaneous update = never read a value you've already overwritten.** Either snapshot the old state (Approach 1) or hide the new state in unused bits so the old state survives (Approach 2).
- **Spare bits are free scratch space.** Since cells are only `0`/`1`, bit 1 is idle; use it to carry the next state and finish with a single shift pass — O(1) extra space.
- The 3×3 neighbour sweep with a `dr,dc ∈ {-1,0,1}` double loop (skipping `(0,0)`) is a reusable matrix idiom.
- **Infinite-board follow-up:** track only live cells and their neighbours in a hash set of coordinates, so memory scales with the active region rather than the full (unbounded) grid.

---

## Related Problems

- LeetCode #73 — Set Matrix Zeroes (in-place marking with sentinel encoding)
- LeetCode #200 — Number of Islands (grid neighbour traversal)
- LeetCode #529 — Minesweeper (8-neighbour grid simulation)
- LeetCode #419 — Battleships in a Board (in-place grid scan)
- LeetCode #48 — Rotate Image (in-place matrix transform)
