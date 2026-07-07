package main

import (
	"fmt"
	"sort"
)

// Robot Room Cleaner (LeetCode 489) is an INTERACTIVE problem: the grader hands
// you a Robot exposing only four methods and you must clean every reachable
// cell without ever seeing the grid or your own coordinates.
//
//	Robot API:
//	  Move()      // step forward one cell in the current facing;
//	              // returns true and moves if the next cell is open,
//	              // returns false and stays put if it is a wall/obstacle.
//	  TurnLeft()  // rotate 90° counter-clockwise in place.
//	  TurnRight() // rotate 90° clockwise in place.
//	  Clean()     // clean the current cell.
//
// Because we cannot test against the real grader here, we implement a faithful
// SimRobot backed by a hidden grid (0 = wall, 1 = open) and a start pose, then
// run the cleaning algorithm against it and count how many open cells got
// cleaned — which must equal the number of reachable open cells.

// Robot is the exact interface the cleaning algorithm is allowed to use.
type Robot interface {
	Move() bool
	TurnLeft()
	TurnRight()
	Clean()
}

// SimRobot simulates the grader's robot over a known grid. The algorithm under
// test never inspects these fields — it only calls the interface methods.
type SimRobot struct {
	grid    [][]int         // 1 = open cell, 0 = wall/obstacle
	r, c    int             // current position (row, col)
	dir     int             // facing: 0=up, 1=right, 2=down, 3=left
	cleaned map[[2]int]bool // set of cells Clean() has been called on
}

// Direction vectors indexed by dir: up, right, down, left (clockwise order),
// so TurnRight is dir+1 and TurnLeft is dir+3 (mod 4).
var dRow = [4]int{-1, 0, 1, 0}
var dCol = [4]int{0, 1, 0, -1}

// NewSimRobot builds a simulated robot at (startR, startC) facing up.
func NewSimRobot(grid [][]int, startR, startC int) *SimRobot {
	return &SimRobot{grid: grid, r: startR, c: startC, dir: 0, cleaned: map[[2]int]bool{}}
}

// Move steps forward if the target cell is inside the grid and open.
func (s *SimRobot) Move() bool {
	nr, nc := s.r+dRow[s.dir], s.c+dCol[s.dir]
	// Out of bounds or a wall → blocked, stay put, report false.
	if nr < 0 || nr >= len(s.grid) || nc < 0 || nc >= len(s.grid[0]) || s.grid[nr][nc] == 0 {
		return false
	}
	s.r, s.c = nr, nc // advance one cell in the current facing
	return true
}

// TurnLeft rotates 90° counter-clockwise (dir − 1 mod 4).
func (s *SimRobot) TurnLeft() { s.dir = (s.dir + 3) % 4 }

// TurnRight rotates 90° clockwise (dir + 1 mod 4).
func (s *SimRobot) TurnRight() { s.dir = (s.dir + 1) % 4 }

// Clean records that the current cell has been cleaned.
func (s *SimRobot) Clean() { s.cleaned[[2]int{s.r, s.c}] = true }

// ── Approach 1: Backtracking DFS with Absolute Coordinates ────────────────────
//
// cleanRoomBacktrack cleans the whole room via DFS. Since the robot never
// reveals its position, the algorithm imposes its OWN coordinate frame: it
// starts the (virtual) origin at (0,0) facing "up" and tracks (row, col, dir)
// purely from the moves it issues. A `visited` set of virtual coordinates
// prevents re-cleaning, and after exploring each neighbour the robot physically
// backtracks to the cell (and facing) it came from.
//
// Intuition:
//
//	This is grid DFS where the twist is that "move to neighbour" is a physical
//	action with a cost you must undo. Clean the current cell, then for each of
//	the 4 directions (relative to the current facing): if the neighbour is
//	unvisited and Move() succeeds, recurse into it, then execute a precise
//	"go back" manoeuvre (turn 180°, move, turn 180° back) so the caller's pose
//	is exactly restored before trying the next direction.
//
// Algorithm:
//  1. Maintain visited set of virtual (row,col). Define backtrack(): TurnRight
//     twice, Move(), TurnRight twice — a U-turn, step, U-turn, restoring facing.
//  2. dfs(row, col, dir): mark visited, Clean().
//  3. For k = 0..3 (four turns): compute nextDir = (dir + k) % 4 and the
//     neighbour cell from (row,col) using that direction's delta.
//     - if the neighbour is unvisited and Move() returns true:
//     recurse dfs(neighbour, nextDir); then backtrack() to return.
//     - TurnRight() to advance to the next relative direction.
//  4. Start at dfs(0, 0, 0).
//
// Time:  O(cells) Clean/visited work; O(4·cells) robot operations overall.
// Space: O(cells) for the visited set and the recursion stack.
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
			nd := (dir + k) % 4  // absolute direction we are now facing
			nr := row + dRow[nd] // neighbour row in the virtual frame
			nc := col + dCol[nd] // neighbour col in the virtual frame
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

// cleanedCount returns how many distinct cells the simulated robot cleaned,
// used by main() to check the algorithm cleaned every reachable open cell.
func (s *SimRobot) cleanedCount() int { return len(s.cleaned) }

// reachableOpenCells counts the open cells reachable from (startR,startC) by a
// plain flood fill — the ground truth for how many cells SHOULD be cleaned.
func reachableOpenCells(grid [][]int, startR, startC int) int {
	seen := map[[2]int]bool{}
	var flood func(r, c int)
	flood = func(r, c int) {
		if r < 0 || r >= len(grid) || c < 0 || c >= len(grid[0]) || grid[r][c] == 0 || seen[[2]int{r, c}] {
			return
		}
		seen[[2]int{r, c}] = true
		flood(r-1, c)
		flood(r+1, c)
		flood(r, c-1)
		flood(r, c+1)
	}
	flood(startR, startC)
	return len(seen)
}

// cleanedCells returns the sorted list of cleaned cells, for a deterministic
// printout in main().
func (s *SimRobot) cleanedCells() [][2]int {
	out := make([][2]int, 0, len(s.cleaned))
	for cell := range s.cleaned {
		out = append(out, cell)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i][0] != out[j][0] {
			return out[i][0] < out[j][0]
		}
		return out[i][1] < out[j][1]
	})
	return out
}

func main() {
	// Official example grid (1 = open/accessible, 0 = wall), robot at row=1,col=3.
	//   [ [1,1,1,1,1,0,1,1],
	//     [1,1,1,1,1,0,1,1],
	//     [1,0,1,1,1,1,1,1],
	//     [0,0,0,1,0,0,0,0],
	//     [1,1,1,1,1,1,1,1] ]
	grid := [][]int{
		{1, 1, 1, 1, 1, 0, 1, 1},
		{1, 1, 1, 1, 1, 0, 1, 1},
		{1, 0, 1, 1, 1, 1, 1, 1},
		{0, 0, 0, 1, 0, 0, 0, 0},
		{1, 1, 1, 1, 1, 1, 1, 1},
	}
	startR, startC := 1, 3

	fmt.Println("=== Approach 1: Backtracking DFS (absolute coordinates) ===")
	robot := NewSimRobot(grid, startR, startC)
	cleanRoomBacktrack(robot)
	want := reachableOpenCells(grid, startR, startC)
	got := robot.cleanedCount()
	fmt.Printf("cleaned %d cells, reachable open cells = %d\n", got, want) // expected: cleaned 30 cells, reachable open cells = 30
	fmt.Println("all reachable cells cleaned:", got == want)               // expected: true
	fmt.Println("cleaned cells:", robot.cleanedCells())                    // expected: full sorted list of the 30 open cells

	// A second, tiny grid to show generality: an L-shaped open region, start top-left.
	//   [ [1,0],
	//     [1,1] ]
	grid2 := [][]int{
		{1, 0},
		{1, 1},
	}
	fmt.Println("=== Approach 1 on a 2x2 L-shaped room ===")
	robot2 := NewSimRobot(grid2, 0, 0)
	cleanRoomBacktrack(robot2)
	want2 := reachableOpenCells(grid2, 0, 0)
	got2 := robot2.cleanedCount()
	fmt.Printf("cleaned %d cells, reachable open cells = %d\n", got2, want2) // expected: cleaned 3 cells, reachable open cells = 3
	fmt.Println("all reachable cells cleaned:", got2 == want2)               // expected: true
}
