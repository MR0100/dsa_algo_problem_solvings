# 0348 — Design Tic-Tac-Toe

> LeetCode #348 · Difficulty: Medium
> **Categories:** Design, Array, Matrix, Hash Table, Simulation

---

## Problem Statement

Assume the following rules are for the tic-tac-toe game on an `n x n` board between two players:

1. A move is guaranteed to be valid and is placed on an empty block.
2. Once a winning condition is reached, no more moves are allowed.
3. A player who succeeds in placing `n` of their marks in a horizontal, vertical, or diagonal row wins the game.

Implement the `TicTacToe` class:

- `TicTacToe(int n)` Initializes the object the size of the board `n`.
- `int move(int row, int col, int player)` Indicates that the player with id `player` plays at the cell `(row, col)` of the board. The move is guaranteed to be a valid move, and the two players alternate in making moves. Return
  - `0` if there is **no winner** after the move,
  - `1` if **player 1** is the winner after the move, or
  - `2` if **player 2** is the winner after the move.

**Example 1:**
```
Input
["TicTacToe", "move", "move", "move", "move", "move", "move", "move"]
[[3], [0, 0, 1], [0, 2, 2], [2, 2, 1], [1, 1, 2], [2, 0, 1], [0, 1, 2], [2, 1, 1]]
Output
[null, 0, 0, 0, 0, 0, 0, 1]

Explanation
TicTacToe ticTacToe = new TicTacToe(3);
// Assume that player 1 is "X" and player 2 is "O" in the board.
ticTacToe.move(0, 0, 1); // return 0 (no one wins)
|X| | |
| | | |    // Player 1 makes a move at (0, 0).
| | | |

ticTacToe.move(0, 2, 2); // return 0 (no one wins)
|X| |O|
| | | |    // Player 2 makes a move at (0, 2).
| | | |

ticTacToe.move(2, 2, 1); // return 0 (no one wins)
|X| |O|
| | | |    // Player 1 makes a move at (2, 2).
| | |X|

ticTacToe.move(1, 1, 2); // return 0 (no one wins)
|X| |O|
| |O| |    // Player 2 makes a move at (1, 1).
| | |X|

ticTacToe.move(2, 0, 1); // return 0 (no one wins)
|X| |O|
| |O| |    // Player 1 makes a move at (2, 0).
|X| |X|

ticTacToe.move(0, 1, 2); // return 0 (no one wins)
|X|O|O|
| |O| |    // Player 2 makes a move at (0, 1).
|X| |X|

ticTacToe.move(2, 1, 1); // return 1 (player 1 wins)
|X|O|O|
| |O| |    // Player 1 makes a move at (2, 1).
|X|X|X|
```

**Constraints:**
- `2 <= n <= 100`
- `player` is `1` or `2`.
- `0 <= row, col < n`
- `(row, col)` are **unique** for each different call to `move`.
- At most `n²` calls will be made to `move`.

**Follow-up:** Could you do better than `O(n²)` per `move` operation?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Facebook   | ★★★☆☆ Medium     | 2024          |
| Apple      | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used
- **Design Data Structures** — a stateful class that answers a query incrementally after each update → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)
- **Matrix Traversal** — the brute-force check scans a row, column, and the two diagonals of the grid → see [`/dsa/matrix_traversal.md`](/dsa/matrix_traversal.md)
- **Counting / Signed Aggregates** — the optimal method keeps a signed sum per line so a win is detected by a counter reaching ±n → see [`/dsa/prefix_sum.md`](/dsa/prefix_sum.md)

---

## Approaches Overview
| # | Approach | Time per move | Space | When to use |
|---|----------|---------------|-------|-------------|
| 1 | Full Board Scan (Brute Force) | O(n) | O(n²) | Straightforward; keeps the visible board |
| 2 | Per-Line Counters (Optimal) | O(1) | O(n) | Meets the follow-up; no board needed |

---

## Approach 1 — Full Board Scan (Brute Force)

### Intuition
Keep the whole board. After placing a mark at `(row, col)`, only four lines could possibly have just been completed: that row, that column, the main diagonal (if `row == col`), and the anti-diagonal (if `row + col == n-1`). Scan just those and see if any is entirely the current player.

### Algorithm
1. Write `player` into `board[row][col]`.
2. Scan the row; if every cell equals `player`, return `player`.
3. Scan the column likewise.
4. If `row == col`, scan the main diagonal.
5. If `row + col == n-1`, scan the anti-diagonal.
6. Otherwise return `0`.

### Complexity
- **Time:** O(n) per move — up to four line scans, each of length n.
- **Space:** O(n²) — the full board is stored.

### Code
```go
type TicTacToeBrute struct {
	n     int     // board dimension
	board [][]int // board[r][c] is 0 (empty), 1, or 2
}

func NewTicTacToeBrute(n int) *TicTacToeBrute {
	board := make([][]int, n)
	for i := range board {
		board[i] = make([]int, n) // all zeros = empty
	}
	return &TicTacToeBrute{n: n, board: board}
}

func (t *TicTacToeBrute) Move(row, col, player int) int {
	t.board[row][col] = player // record the move
	n := t.n

	rowWin := true
	for c := 0; c < n; c++ {
		if t.board[row][c] != player {
			rowWin = false
			break
		}
	}
	if rowWin {
		return player
	}

	colWin := true
	for r := 0; r < n; r++ {
		if t.board[r][col] != player {
			colWin = false
			break
		}
	}
	if colWin {
		return player
	}

	if row == col { // main diagonal
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

	if row+col == n-1 { // anti-diagonal
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

	return 0
}
```

### Dry Run
Board size 3. Only the final move is a win; earlier moves fail their line scans.

| move(row,col,player) | board[row] after | row all player? | col all player? | on diag? | result |
|----------------------|------------------|-----------------|-----------------|----------|--------|
| (0,0,1) | [1,0,0] | no (0s) | no | yes but col not full | 0 |
| (0,2,2) | [1,0,2] | no | no | anti-diag not full | 0 |
| (2,2,1) | [0,0,1] | no | no | diag not full | 0 |
| (1,1,2) | [0,2,0] | no | no | diag/anti not full | 0 |
| (2,0,1) | [1,0,1] | no | no | anti not full | 0 |
| (0,1,2) | [1,2,2] | no | no | — | 0 |
| (2,1,1) | [1,1,1] | **yes** | — | — | **1** |

On `(2,1,1)`, row 2 = `[1,1,1]`, so player 1 wins.

---

## Approach 2 — Per-Line Counters (Optimal)

### Intuition
A player wins a line when they own all n cells. Encode player 1 as `+1` and player 2 as `−1`, and keep a **signed sum per line** — one per row, one per column, and one for each diagonal. Placing a mark adds `±1` to that line's counter. Player 1 owns a line exactly when its counter hits `+n`; player 2 when it hits `−n`. Each move touches at most four counters, so the win check is O(1) and no board is stored.

### Algorithm
1. `delta = +1` for player 1, `−1` for player 2.
2. `rows[row] += delta`; `cols[col] += delta`.
3. If `row == col`: `diag += delta`. If `row + col == n-1`: `anti += delta`.
4. If any touched counter has absolute value `n`, `player` wins; else `0`.

### Complexity
- **Time:** O(1) per move — a constant number of counter updates and comparisons.
- **Space:** O(n) — one counter per row and per column, plus two diagonal scalars.

### Code
```go
type TicTacToeOptimal struct {
	n    int   // board dimension and the win threshold
	rows []int // signed sum per row
	cols []int // signed sum per column
	diag int   // signed sum of the main diagonal
	anti int   // signed sum of the anti-diagonal
}

func NewTicTacToeOptimal(n int) *TicTacToeOptimal {
	return &TicTacToeOptimal{n: n, rows: make([]int, n), cols: make([]int, n)}
}

func (t *TicTacToeOptimal) Move(row, col, player int) int {
	delta := 1 // player 1 contributes +1
	if player == 2 {
		delta = -1 // player 2 contributes −1
	}

	t.rows[row] += delta
	t.cols[col] += delta
	if row == col {
		t.diag += delta
	}
	if row+col == t.n-1 {
		t.anti += delta
	}

	n := t.n
	if abs(t.rows[row]) == n || abs(t.cols[col]) == n ||
		abs(t.diag) == n || abs(t.anti) == n {
		return player
	}
	return 0
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
```

### Dry Run
Board size 3, threshold `n = 3`. rows/cols start at 0, diag=0, anti=0.

| move(row,col,player) | delta | rows[row] | cols[col] | diag | anti | any \|·\|==3? | result |
|----------------------|-------|-----------|-----------|------|------|--------------|--------|
| (0,0,1) | +1 | rows[0]=1 | cols[0]=1 | 1 | — | no | 0 |
| (0,2,2) | −1 | rows[0]=0 | cols[2]=−1 | — | anti=−1 | no | 0 |
| (2,2,1) | +1 | rows[2]=1 | cols[2]=0 | diag=2 | — | no | 0 |
| (1,1,2) | −1 | rows[1]=−1 | cols[1]=−1 | diag=1 | anti=−2 | no | 0 |
| (2,0,1) | +1 | rows[2]=2 | cols[0]=2 | — | anti=−1 | no | 0 |
| (0,1,2) | −1 | rows[0]=−1 | cols[1]=−2 | — | — | no | 0 |
| (2,1,1) | +1 | **rows[2]=3** | cols[1]=−1 | — | — | **yes (rows[2]=3)** | **1** |

The final move drives `rows[2]` to `+3 == n`, so player 1 wins in O(1).

---

## Key Takeaways
- **Turn "line is fully owned" into a counter.** Encoding the two players as `+1`/`−1` collapses each row/column/diagonal into a single signed sum; a win is `|sum| == n`, checkable in O(1).
- **Only update what the move touches.** A move at `(r,c)` affects at most its row, its column, and the two diagonals — never rescan the board.
- The diagonal membership tests are the classic `row == col` (main) and `row + col == n-1` (anti).
- This is the canonical "replace an O(n) recompute with an O(1) incremental aggregate" design pattern, the same idea behind running-sum stream problems.

---

## Related Problems
- LeetCode #794 — Valid Tic-Tac-Toe State (board validity reasoning)
- LeetCode #1275 — Find Winner on a Tic Tac Toe Game (same line-counting idea)
- LeetCode #361 — Bomb Enemy (row/column aggregate scanning on a grid)
- LeetCode #2013 — Detect Squares (incremental geometric counting design)
