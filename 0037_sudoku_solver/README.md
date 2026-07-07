# 0037 — Sudoku Solver

> LeetCode #37 · Difficulty: Hard
> **Categories:** Array, Hash Table, Backtracking, Matrix

---

## Problem Statement

Write a program to solve a Sudoku puzzle by filling in the empty cells.

A sudoku solution must satisfy **all of the following rules**:
1. Each of the digits `1-9` must occur exactly once in each row.
2. Each of the digits `1-9` must occur exactly once in each column.
3. Each of the digits `1-9` must occur exactly once in each of the 9 `3x3` sub-boxes of the grid.

Empty cells are indicated by the character `'.'`.

**Note:** The input board is a `9x9` grid and it is guaranteed to have a unique solution.

**Example 1**
```
Input:
[["5","3",".",".","7",".",".",".","."]
,["6",".",".","1","9","5",".",".","."]
,...  (see main.go)]
Output: Board solved in-place.
```

**Constraints**
- `board.length == 9`
- `board[i].length == 9`
- `board[i][j]` is a digit or `'.'`.
- The input board has exactly one solution.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★☆ High      | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Meta      | ★★★☆☆ Medium    | 2023          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |
| Apple     | ★★★☆☆ Medium    | 2023          |
| Adobe     | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Backtracking** — try each digit 1–9 in each empty cell; undo if a conflict arises.
- **Constraint Tracking (Boolean Arrays)** — precompute `rows[r][d]`, `cols[c][d]`, `boxes[boxID][d]` for O(1) validity checking instead of scanning the board each time.
- **Index Mapping** — `boxID = (r/3)*3 + (c/3)` maps cells to their 3×3 box.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Backtracking with Constraint Arrays ✅ | O(9^m) | O(1) | The standard and only practical approach |

m = number of empty cells. In practice, constraint propagation reduces the branching factor dramatically.

---

## Approach 1 — Backtracking with Precomputed Constraints (Recommended ✅)

### Intuition
Try digits 1–9 for each empty cell in order. Before placing digit `d`, check the three constraint arrays in O(1). If valid, place `d`, update the constraint arrays, and recurse to the next empty cell. If recursion succeeds, return true (board is solved). If no digit works (all 1–9 conflict), unplace (backtrack) and return false to trigger backtracking in the caller.

The key optimisation is precomputing `rows/cols/boxes` from the initial clues — this avoids scanning 27 cells at each placement step.

### Algorithm
```
precompute rows[r][d], cols[c][d], boxes[boxID][d] from existing digits
backtrack(pos):
  advance pos to the next '.'
  if pos == 81: return true (board complete)
  r, c = pos/9, pos%9; boxID = (r/3)*3 + (c/3)
  for d = 1 to 9:
    if rows[r][d] or cols[c][d] or boxes[boxID][d]: skip
    place d; mark rows/cols/boxes; recurse(pos+1)
    if recursion returns true: return true
    unplace d; unmark rows/cols/boxes
  return false
```

### Complexity
- **Time:** O(9^m) worst case — m empty cells × 9 choices each. In practice exponentially faster due to constraints.
- **Space:** O(m) — recursion depth equals the number of empty cells (at most 81).

### Code
```go
func solveSudoku(board [][]byte) {
    var rows, cols, boxes [9][10]bool
    for r := 0; r < 9; r++ {
        for c := 0; c < 9; c++ {
            if board[r][c] != '.' {
                d := board[r][c] - '0'
                bid := (r/3)*3 + (c/3)
                rows[r][d] = true; cols[c][d] = true; boxes[bid][d] = true
            }
        }
    }
    var bt func(pos int) bool
    bt = func(pos int) bool {
        for pos < 81 && board[pos/9][pos%9] != '.' { pos++ }
        if pos == 81 { return true }
        r, c := pos/9, pos%9; bid := (r/3)*3 + (c/3)
        for d := 1; d <= 9; d++ {
            if rows[r][d] || cols[c][d] || boxes[bid][d] { continue }
            board[r][c] = byte('0' + d)
            rows[r][d] = true; cols[c][d] = true; boxes[bid][d] = true
            if bt(pos+1) { return true }
            board[r][c] = '.'
            rows[r][d] = false; cols[c][d] = false; boxes[bid][d] = false
        }
        return false
    }
    bt(0)
}
```

### Dry Run — first empty cell (0,2) in the standard example
```
r=0, c=2, boxID=0
rows[0] has: 5,3,7  → digits 5,3,7 used in row 0
cols[2] has: 8      → digit 8 used in column 2
boxes[0] has: 5,3,9,6,8 → digits 5,3,6,8,9 used in box 0

Try d=1: rows[0][1]=false, cols[2][1]=false, boxes[0][1]=false → place 1
  Recurse to next '.' ...eventually finds conflict → backtrack, remove 1
Try d=2: similar...
Try d=4: 4 not in row0, col2, box0 → place 4
  Recurse succeeds → board[0][2]='4' ✓
```

---

## Key Takeaways

- **Constraint arrays vs. board scan** — scanning the board for each placement is O(27) per attempt; precomputed boolean arrays make each check O(1), reducing constant factors significantly.
- **Backtracking terminates when pos == 81** — the linear scan `pos++` to find the next `'.'` means we never need to maintain a list of empty cells explicitly.
- **Guaranteed unique solution** — if the problem guarantees one solution, we can return true as soon as the board is complete. Without this guarantee, we'd need to collect all solutions.
- **Advanced: Dancing Links / Algorithm X** — for competitive programming, Donald Knuth's Algorithm X with Dancing Links solves any exact cover problem (including Sudoku) optimally. Not needed for interviews.

---

## Related Problems

- LeetCode #36 — Valid Sudoku (validation, not solving)
- LeetCode #51 — N-Queens (backtracking with row/col/diagonal constraints)
- LeetCode #52 — N-Queens II (count solutions)
