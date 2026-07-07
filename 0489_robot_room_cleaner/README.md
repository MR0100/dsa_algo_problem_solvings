# 0489 — Robot Room Cleaner

> LeetCode #489 · Difficulty: Hard
> **Categories:** Backtracking, Depth-First Search, Interactive, Matrix

---

## Problem Statement

You are controlling a robot that is located somewhere in a room. The room is modeled as an `m x n` binary grid where `0` represents a wall and `1` represents an empty slot.

The robot starts at an unknown location in the room that is guaranteed to be empty, and you do not have access to the grid, but you can move the robot using the given API `Robot`.

You are tasked to use the robot to clean the entire room (i.e., clean every empty cell in the room). The robot with the four given APIs can move forward, turn left, or turn right. Each turn is `90` degrees.

When the robot tries to move into a wall cell, its bumper sensor detects the obstacle, and it stays on the current cell.

Design an algorithm to clean the entire room using the following APIs:

```
interface Robot {
  // returns true if the next cell is open and robot moves into the cell.
  // returns false if the next cell is an obstacle and robot stays on the current cell.
  boolean move();

  // Robot will stay on the same cell after calling turnLeft/turnRight.
  // Each turn will be 90 degrees.
  void turnLeft();
  void turnRight();

  // Clean the current cell.
  void clean();
}
```

**Note** that the initial direction of the robot will be facing up. You can assume all four edges of the grid are all surrounded by a wall.

**Example 1:**

```
Input: room = [[1,1,1,1,1,0,1,1],[1,1,1,1,1,0,1,1],[1,0,1,1,1,1,1,1],[0,0,0,1,0,0,0,0],[1,1,1,1,1,1,1,1]], row = 1, col = 3
Output: Robot cleaned all rooms.
Explanation: All grids in the room are marked by either 0 or 1.
0 means the cell is blocked, while 1 means the cell is accessible.
The robot initially starts at the position of row=1, col=3.
From the top left corner, its position is one row below and three columns right.
```

**Example 2:**

```
Input: room = [[1]], row = 0, col = 0
Output: Robot cleaned all rooms.
```

**Constraints:**

- `m == room.length`
- `n == room[i].length`
- `1 <= m <= 100`
- `1 <= n <= 200`
- `room[i][j]` is either `0` or `1`.
- `0 <= row < m`
- `0 <= col < n`
- `room[row][col] == 1`
- All the empty cells can be visited from the starting position.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★★ Very High  | 2024          |
| Facebook   | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Uber       | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Backtracking** — after exploring a neighbour the robot must physically return to the previous cell and restore its facing (a precise "undo the move"), which is backtracking made literal → see [`/dsa/backtracking.md`](/dsa/backtracking.md)
- **DFS on an implicit grid graph** — cells are nodes, open neighbours are edges; the room is explored depth-first from the start cell → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Matrix traversal with direction vectors** — moves use `(dRow, dCol)` arrays in clockwise order so `TurnRight` is `+1 mod 4`, letting relative turns map cleanly to absolute directions → see [`/dsa/matrix_traversal.md`](/dsa/matrix_traversal.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Backtracking DFS (absolute coordinates) | O(cells) cleans, O(4·cells) robot ops | O(cells) | The canonical (and essentially only) solution to this interactive problem |

> This problem has one accepted technique. The "approaches" here are really the
> one correct algorithm; what makes it Hard is reasoning about a robot that hides
> its own position, forcing you to build a private coordinate frame and to undo
> every physical move.

---

## Approach 1 — Backtracking DFS (absolute coordinates)

### Intuition

The robot never reveals where it is, so the algorithm imposes its **own** coordinate frame: pretend the start cell is `(0, 0)` facing "up", and update a virtual `(row, col, dir)` from the moves you issue. This is ordinary grid DFS, except "go to a neighbour" is a *physical* action that costs a move and must be undone afterward. At each cell: `Clean()`, then for each of the four directions relative to the current facing, if the neighbour is unvisited and `Move()` succeeds, recurse into it and then perform a fixed **backtrack manoeuvre** — turn 180°, step forward, turn 180° back — so the robot's pose is *exactly* what it was before the excursion. A `visited` set of virtual coordinates stops re-cleaning and infinite loops.

Key device: keep the direction arrays in **clockwise** order (up, right, down, left). Then issuing one `TurnRight()` per loop iteration walks the robot's facing through `dir, dir+1, dir+2, dir+3`; after four turns the facing returns to where it started, so the cell is left in a clean, predictable state for its parent.

### Algorithm

1. Maintain `visited` (virtual `(row,col)`). Define `backtrack()` = `TurnRight, TurnRight, Move, TurnRight, TurnRight` (U-turn → step back → U-turn), which restores facing.
2. `dfs(row, col, dir)`: mark `visited[(row,col)]`, call `Clean()`.
3. For `k = 0..3`:
   - `nd = (dir + k) % 4`; neighbour `(nr, nc) = (row + dRow[nd], col + dCol[nd])`.
   - if `(nr,nc)` is unvisited and `Move()` returns true: `dfs(nr, nc, nd)` then `backtrack()`.
   - `TurnRight()` to line up the next relative direction.
4. Start with `dfs(0, 0, 0)`.

### Complexity

- **Time:** O(cells) `Clean`/`visited` operations; the robot performs O(4·cells) moves and turns overall (a constant amount of physical work per cell).
- **Space:** O(cells) for the `visited` set plus the DFS recursion stack (depth ≤ number of cells).

### Code

```go
func cleanRoomBacktrack(robot Robot) {
	visited := map[[2]int]bool{} // virtual coordinates already cleaned

	// backtrack physically returns the robot to the previous cell, restoring
	// its original facing: U-turn, step forward, U-turn again.
	backtrack := func() {
		robot.TurnRight()
		robot.TurnRight() // now facing 180° from before
		robot.Move()      // step back into the cell we came from
		robot.TurnRight()
		robot.TurnRight() // restore the original facing
	}

	var dfs func(row, col, dir int)
	dfs = func(row, col, dir int) {
		visited[[2]int{row, col}] = true // record this virtual cell
		robot.Clean()                    // clean where we currently stand

		// Explore all four directions relative to the current facing.
		for k := 0; k < 4; k++ {
			nd := (dir + k) % 4          // absolute direction we are now facing
			nr := row + dRow[nd]         // neighbour row in the virtual frame
			nc := col + dCol[nd]         // neighbour col in the virtual frame
			if !visited[[2]int{nr, nc}] && robot.Move() {
				dfs(nr, nc, nd) // Move() succeeded → we are in the neighbour now
				backtrack()     // return to (row,col) with facing == nd
			}
			robot.TurnRight() // rotate to the next relative direction (dir+k+1)
		}
		// Four TurnRight() calls net to a full 360°, so facing is unchanged here.
	}

	dfs(0, 0, 0) // origin at (0,0), facing "up"
}
```

Direction vectors (shared, clockwise so `TurnRight` = `dir+1`):

```go
var dRow = [4]int{-1, 0, 1, 0} // up, right, down, left
var dCol = [4]int{0, 1, 0, -1}
```

### Dry Run

Start the robot at virtual `(0,0)` facing up (`dir = 0`). Trace the first excursion; `U` = up, `R` = right, `D` = down, `L` = left.

| Step | At (row,col,dir) | k | nd (absolute) | neighbour | visited? | Move()? | Action |
|------|------------------|---|---------------|-----------|----------|---------|--------|
| 1 | (0,0,U) | — | — | — | — | — | mark (0,0), Clean() |
| 2 | (0,0), facing U | 0 | U | (−1,0) | no | **false** (wall/edge) | can't go; TurnRight → facing R |
| 3 | (0,0), facing R | 1 | R | (0,1) | no | **true** | recurse dfs(0,1,R) … |
| 3a | (0,1,R) | — | — | — | — | — | mark (0,1), Clean(), explore its 4 dirs |
| 3b | … eventually returns | — | — | — | — | — | backtrack(): U-turn, Move back to (0,0), U-turn → facing R again |
| 4 | (0,0), facing R | — | — | — | — | — | TurnRight → facing D |
| 5 | (0,0), facing D | 2 | D | (1,0) | no | **true** | recurse dfs(1,0,D) … then backtrack, TurnRight → facing L |
| 6 | (0,0), facing L | 3 | L | (0,−1) | no | **false** (edge) | TurnRight → facing U (back to start facing) |

After the loop the robot faces up again at `(0,0)`, exactly as it began — so its parent (there is none for the root) would find the pose untouched. Recursion into `(0,1)` and `(1,0)` repeats this pattern, and because `visited` blocks re-entry, every reachable open cell is `Clean()`ed exactly once. In `main()` the simulator confirms **30 of 30** reachable open cells were cleaned for the official grid. ✔

---

## Key Takeaways

- **No self-localization? Invent a coordinate frame.** Fix the origin at the start `(0,0)` and track `(row, col, dir)` from the moves you command; the absolute grid position is irrelevant — only relative geometry matters.
- **Backtracking is physical here.** Every recursive descent that issues a `Move()` must be paired with an exact inverse manoeuvre (`TurnRight×2, Move, TurnRight×2`) so the caller's pose is restored. Forgetting the pose restore is the #1 bug.
- **Clockwise direction arrays make turns arithmetic.** With `[up, right, down, left]`, `TurnRight` is `+1 mod 4` and issuing one `TurnRight` per loop cycles cleanly through all four headings, netting a full rotation over four iterations.
- **`visited` uses your virtual coordinates, not the robot's.** The set keyed by `(row, col)` in your own frame is what prevents infinite loops and double-cleaning.

---

## Related Problems

- LeetCode #200 — Number of Islands (grid DFS/flood fill)
- LeetCode #79 — Word Search (grid DFS with backtracking/undo)
- LeetCode #490 — The Maze (rolling-ball BFS/DFS on a grid)
- LeetCode #874 — Walking Robot Simulation (direction-vector robot movement)
- LeetCode #588 — Design In-Memory File System (interactive/design-style modeling)
