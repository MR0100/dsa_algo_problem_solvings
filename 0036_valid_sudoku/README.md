# 0036 — Valid Sudoku

> LeetCode #36 · Difficulty: Medium
> **Categories:** Array, Hash Table, Matrix

---

## Problem Statement

Determine if a `9 x 9` Sudoku board is valid. Only the filled cells need to be validated according to the following rules:

1. Each row must contain the digits `1-9` without repetition.
2. Each column must contain the digits `1-9` without repetition.
3. Each of the nine `3 x 3` sub-boxes of the grid must contain the digits `1-9` without repetition.

**Note:** A Sudoku board (partially filled) could be valid but is not necessarily solvable. Only the filled cells need to be validated.

**Example 1**
```
Input:
[["5","3",".",".","7",".",".",".","."]
,["6",".",".","1","9","5",".",".","."]
,[".","9","8",".",".",".",".","6","."]
,["8",".",".",".","6",".",".",".","3"]
,["4",".",".","8",".","3",".",".","1"]
,["7",".",".",".","2",".",".",".","6"]
,[".","6",".",".",".",".","2","8","."]
,[".",".",".","4","1","9",".",".","5"]
,[".",".",".",".","8",".",".","7","9"]]
Output: true
```

**Example 2**
```
Input: Same board but first cell changed from "5" to "8"
       (8 appears twice in column 0)
Output: false
```

**Constraints**
- `board.length == 9`
- `board[i].length == 9`
- `board[i][j]` is a digit `1-9` or `'.'`.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★☆ High      | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Meta      | ★★★☆☆ Medium    | 2023          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |
| Apple     | ★★☆☆☆ Low       | 2022          |
| Adobe     | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Set / Boolean Array** — track seen digits per row, column, and box.
- **Index Mapping** — `boxID = (r/3)*3 + (c/3)` maps each cell to one of 9 box IDs (0–8).

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Three Separate Passes | O(1) | O(1) | Clearer to read; conceptually obvious |
| 2 | Single Pass (Three Arrays) ✅ | O(1) | O(1) | Optimal; one traversal of 81 cells |

Both are O(1) since the board is always 9×9 = 81 cells.

---

## Approach 1 — Three Separate Passes

### Intuition
Validate rows, columns, and boxes in three separate loops. Each loop uses a boolean array `seen[10]` reset between checks.

### Algorithm
- Pass 1: for each row r, check digits across 9 columns.
- Pass 2: for each column c, check digits across 9 rows.
- Pass 3: for each of 9 boxes (r in 0..2, c in 0..2), check digits across the 3×3 sub-grid.

### Complexity
- **Time:** O(81 × 3) = O(1).
- **Space:** O(9) per check = O(1).

### Code
```go
func threePasses(board [][]byte) bool {
    // check all rows
    for r := 0; r < 9; r++ {
        seen := [10]bool{}
        for c := 0; c < 9; c++ {
            if board[r][c] == '.' {
                continue
            }
            d := board[r][c] - '0'
            if seen[d] {
                return false
            }
            seen[d] = true
        }
    }
    // check all columns
    for c := 0; c < 9; c++ {
        seen := [10]bool{}
        for r := 0; r < 9; r++ {
            if board[r][c] == '.' {
                continue
            }
            d := board[r][c] - '0'
            if seen[d] {
                return false
            }
            seen[d] = true
        }
    }
    // check all 3×3 boxes
    for boxRow := 0; boxRow < 3; boxRow++ {
        for boxCol := 0; boxCol < 3; boxCol++ {
            seen := [10]bool{}
            for r := boxRow * 3; r < boxRow*3+3; r++ {
                for c := boxCol * 3; c < boxCol*3+3; c++ {
                    if board[r][c] == '.' {
                        continue
                    }
                    d := board[r][c] - '0'
                    if seen[d] {
                        return false
                    }
                    seen[d] = true
                }
            }
        }
    }
    return true
}
```

### Dry Run — Example 1 (valid board)
Each pass resets `seen[10]` per row/column/box and returns `false` on the first repeat. On the valid board every pass completes with no repeat, so the function returns `true`.

**Pass 1 — rows.** Row 0 = `5 3 . . 7 . . . .`:

| c | cell | d | seen[d] before | action |
|---|------|---|----------------|--------|
| 0 | 5 | 5 | false | mark seen[5] |
| 1 | 3 | 3 | false | mark seen[3] |
| 2 | . | — | —     | skip |
| 3 | . | — | —     | skip |
| 4 | 7 | 7 | false | mark seen[7] |
| 5–8 | . | — | —   | skip |

Row 0 has no repeat; rows 1–8 likewise pass.

**Pass 2 — columns.** Column 0 = `5 6 . 8 4 7 . . .` → digits {5,6,8,4,7}, all distinct → pass. Columns 1–8 likewise pass.

**Pass 3 — boxes.** Box 0 (rows 0–2, cols 0–2) = `5 3 . / 6 . . / . 9 8` → digits {5,3,6,9,8}, all distinct → pass. Boxes 1–8 likewise pass.

All three passes finish without a repeat → return `true` ✓.

---

## Approach 2 — Single Pass (Recommended ✅)

### Intuition
Traverse the board once. For each filled cell `(r, c)` with digit `d`:
- Check `rows[r][d]`, `cols[c][d]`, `boxes[boxID][d]`.
- If any is true: duplicate found → return false.
- Mark all three as true.

The box ID formula `(r/3)*3 + (c/3)` maps:
```
(0,0)–(2,2) → boxID 0    (0,3)–(2,5) → boxID 1    (0,6)–(2,8) → boxID 2
(3,0)–(5,2) → boxID 3    (3,3)–(5,5) → boxID 4    (3,6)–(5,8) → boxID 5
(6,0)–(8,2) → boxID 6    (6,3)–(8,5) → boxID 7    (6,6)–(8,8) → boxID 8
```

### Complexity
- **Time:** O(81) = O(1).
- **Space:** O(9 × 9 × 3) = O(1) — three fixed-size 9×10 arrays.

### Code
```go
func singlePass(board [][]byte) bool {
    var rows, cols, boxes [9][10]bool
    for r := 0; r < 9; r++ {
        for c := 0; c < 9; c++ {
            if board[r][c] == '.' { continue }
            d := board[r][c] - '0'
            boxID := (r/3)*3 + (c/3)
            if rows[r][d] || cols[c][d] || boxes[boxID][d] { return false }
            rows[r][d] = true; cols[c][d] = true; boxes[boxID][d] = true
        }
    }
    return true
}
```

### Dry Run — cell `(0,0)` with digit 5
```
d=5, boxID=(0/3)*3+(0/3)=0
rows[0][5]=false, cols[0][5]=false, boxes[0][5]=false → no conflict
Mark rows[0][5]=cols[0][5]=boxes[0][5]=true
```

---

## Key Takeaways

- **`boxID = (r/3)*3 + (c/3)`** — the most important formula. Integer division maps rows 0–2→0, 3–5→1, 6–8→2; multiply row-group by 3 and add col-group.
- **Validity ≠ solvability** — a valid board may have no solution. This problem only checks consistency of the filled-in digits.
- **Constant time and space** — the board size is fixed (9×9); there is no asymptotic growth. Technically O(1) despite looking like O(n²).

---

## Related Problems

- LeetCode #37 — Sudoku Solver (fill the board using backtracking)
