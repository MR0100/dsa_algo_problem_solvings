package main

import "fmt"

// ── Approach 1: Full Board Scan (Brute Force) ────────────────────────────────
//
// TicTacToeBrute stores the whole n×n grid and, after each move, checks whether
// the row, column, or (if on a diagonal) either diagonal is entirely the mover.
//
// Intuition:
//
//	The most obvious design keeps the board and, after placing a mark, re-scans
//	the affected lines to see if any is now filled with a single player. A move
//	at (row, col) can only complete: that row, that column, the main diagonal
//	(if row==col), and the anti-diagonal (if row+col==n-1). Check just those.
//
// Algorithm:
//  1. Write player into board[row][col].
//  2. Scan the row: if every cell equals player → win.
//  3. Scan the column likewise.
//  4. If row==col, scan the main diagonal.
//  5. If row+col==n-1, scan the anti-diagonal.
//  6. Otherwise return 0 (no winner yet).
//
// Time:  O(n) per move — each of the up-to-four line scans is length n.
// Space: O(n²) — the full board.
type TicTacToeBrute struct {
	n     int     // board dimension
	board [][]int // board[r][c] is 0 (empty), 1, or 2
}

// NewTicTacToeBrute builds an empty n×n board.
func NewTicTacToeBrute(n int) *TicTacToeBrute {
	board := make([][]int, n)
	for i := range board {
		board[i] = make([]int, n) // all zeros = empty
	}
	return &TicTacToeBrute{n: n, board: board}
}

// Move places player's mark at (row, col) and returns the winner (0 if none).
func (t *TicTacToeBrute) Move(row, col, player int) int {
	t.board[row][col] = player // record the move
	n := t.n

	rowWin := true // assume the whole row is player until disproved
	for c := 0; c < n; c++ {
		if t.board[row][c] != player {
			rowWin = false
			break
		}
	}
	if rowWin {
		return player
	}

	colWin := true // same for the column
	for r := 0; r < n; r++ {
		if t.board[r][col] != player {
			colWin = false
			break
		}
	}
	if colWin {
		return player
	}

	if row == col { // move sits on the main diagonal
		diagWin := true
		for i := 0; i < n; i++ {
			if t.board[i][i] != player {
				diagWin = false
				break
			}
		}
		if diagWin {
			return player
		}
	}

	if row+col == n-1 { // move sits on the anti-diagonal
		antiWin := true
		for i := 0; i < n; i++ {
			if t.board[i][n-1-i] != player {
				antiWin = false
				break
			}
		}
		if antiWin {
			return player
		}
	}

	return 0 // nobody has won yet
}

// ── Approach 2: Per-Line Counters (Optimal) ──────────────────────────────────
//
// TicTacToeOptimal keeps only running signed counts per row, per column, and
// for the two diagonals — no board at all — giving O(1) per move.
//
// Intuition:
//
//	A player wins a line when they own all n of its cells. Encode player 1 as
//	+1 and player 2 as −1 and keep a signed sum per line. Placing a mark adds
//	±1 to that line's counter; the line is fully owned by player 1 exactly when
//	its counter reaches +n, and by player 2 when it reaches −n. We only ever
//	touch four counters per move (its row, its column, and the two diagonals if
//	applicable), so no scanning is needed.
//
// Algorithm:
//  1. delta = +1 for player 1, −1 for player 2.
//  2. rows[row]+=delta; cols[col]+=delta.
//  3. If row==col: diag+=delta. If row+col==n-1: anti+=delta.
//  4. If any touched counter has absolute value n → player wins.
//
// Time:  O(1) per move — a fixed number of counter updates and comparisons.
// Space: O(n) — one counter per row and per column, plus two diagonal counters.
type TicTacToeOptimal struct {
	n    int   // board dimension and the win threshold
	rows []int // signed sum per row
	cols []int // signed sum per column
	diag int   // signed sum of the main diagonal
	anti int   // signed sum of the anti-diagonal
}

// NewTicTacToeOptimal allocates the per-line counters for an n×n board.
func NewTicTacToeOptimal(n int) *TicTacToeOptimal {
	return &TicTacToeOptimal{n: n, rows: make([]int, n), cols: make([]int, n)}
}

// Move records player's mark and returns the winner (0 if none) in O(1).
func (t *TicTacToeOptimal) Move(row, col, player int) int {
	delta := 1 // player 1 contributes +1
	if player == 2 {
		delta = -1 // player 2 contributes −1
	}

	t.rows[row] += delta // update this row's signed sum
	t.cols[col] += delta // update this column's signed sum
	if row == col {
		t.diag += delta // on the main diagonal
	}
	if row+col == t.n-1 {
		t.anti += delta // on the anti-diagonal
	}

	n := t.n
	// A line is fully owned when its |sum| == n. Only the four lines we just
	// touched can have changed, so those are the only ones worth checking.
	if abs(t.rows[row]) == n || abs(t.cols[col]) == n ||
		abs(t.diag) == n || abs(t.anti) == n {
		return player
	}
	return 0
}

// abs returns the absolute value of x.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	// Official Example (n=3):
	// Input:  ["TicTacToe","move","move","move","move","move","move","move"]
	//         [[3],[0,0,1],[0,2,2],[2,2,1],[1,1,2],[2,0,1],[0,1,2],[2,1,1]]
	// Output: [null,0,0,0,0,0,0,1]
	// Player 1 played (0,0),(2,2),(2,0),(2,1): the final move (2,1) fills all
	// three cells of row 2 for player 1, so it returns 1 (player 1 wins).

	moves := [][3]int{
		{0, 0, 1},
		{0, 2, 2},
		{2, 2, 1},
		{1, 1, 2},
		{2, 0, 1},
		{0, 1, 2},
		{2, 1, 1},
	}
	expected := []int{0, 0, 0, 0, 0, 0, 1}

	fmt.Println("=== Approach 1: Full Board Scan (Brute Force) ===")
	b := NewTicTacToeBrute(3)
	for i, m := range moves {
		got := b.Move(m[0], m[1], m[2])
		fmt.Printf("move(%d,%d,%d) -> %d  expected %d\n", m[0], m[1], m[2], got, expected[i])
	}

	fmt.Println("=== Approach 2: Per-Line Counters (Optimal) ===")
	o := NewTicTacToeOptimal(3)
	for i, m := range moves {
		got := o.Move(m[0], m[1], m[2])
		fmt.Printf("move(%d,%d,%d) -> %d  expected %d\n", m[0], m[1], m[2], got, expected[i])
	}
}
