# 0051 — N-Queens

> LeetCode #51 · Difficulty: Hard
> **Categories:** Array, Backtracking

---

## Problem Statement

The **n-queens** puzzle is the problem of placing `n` queens on an `n × n` chessboard such that no two queens attack each other.

Given an integer `n`, return all distinct solutions to the n-queens puzzle. You may return the answer in **any order**.

Each solution contains a distinct board configuration of the n-queens' placement, where `'Q'` and `'.'` both indicate a queen and an empty space, respectively.

**Example 1**
```
Input:  n = 4
Output: [[".Q..","...Q","Q...","..Q."],["..Q.","Q...","...Q",".Q.."]]
Explanation: There exist two distinct solutions to the 4-queens puzzle.
```

**Example 2**
```
Input:  n = 1
Output: [["Q"]]
```

**Constraints**
- `1 <= n <= 9`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★★ Very High | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Meta      | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |
| Apple     | ★★★☆☆ Medium    | 2023          |
| Adobe     | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Backtracking** — place one queen per row; at each step try all valid columns and undo on failure.
- **Diagonal Conflict Detection** — queen at (r, c) attacks diagonal r-c (constant on `\` diagonal) and anti-diagonal r+c (constant on `/` diagonal).

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Backtracking + Map Sets | O(n!) | O(n) | Clear; uses maps for attacked sets |
| 2 | Backtracking + Bool Arrays ✅ | O(n!) | O(n) | Same idea; faster constant with arrays |

---

## Approach 1 — Backtracking with Map Sets

### Intuition
Place queens row by row. At row `r`, try each column `c`. A queen at `(r, c)` conflicts with existing queens if:
- Column `c` is occupied (`cols` set).
- Diagonal `r-c` is occupied (`diag` set) — same value means same `\` diagonal.
- Anti-diagonal `r+c` is occupied (`anti` set) — same value means same `/` diagonal.

Choose → recurse → unchoose (backtrack).

### Algorithm
```
bt(row):
  if row == n: record board; return
  for c = 0 to n-1:
    if cols[c] or diag[row-c] or anti[row+c]: continue
    place queen at (row, c); mark sets
    bt(row+1)
    remove queen; unmark sets
```

### Complexity
- **Time:** O(n!) — n! paths in the worst case (n choices in row 0, n-1 in row 1, …).
- **Space:** O(n) — sets of size ≤ n; recursion stack depth n.

### Code
```go
func backtracking(n int) [][]string {
    var result [][]string
    cols := make(map[int]bool); diag := make(map[int]bool); anti := make(map[int]bool)
    board := make([][]byte, n)
    for i := range board { board[i] = make([]byte, n); for j := range board[i] { board[i][j] = '.' } }

    var bt func(row int)
    bt = func(row int) {
        if row == n {
            snap := make([]string, n)
            for i, r := range board { snap[i] = string(r) }
            result = append(result, snap); return
        }
        for c := 0; c < n; c++ {
            if cols[c] || diag[row-c] || anti[row+c] { continue }
            board[row][c] = 'Q'; cols[c] = true; diag[row-c] = true; anti[row+c] = true
            bt(row+1)
            board[row][c] = '.'; cols[c] = false; diag[row-c] = false; anti[row+c] = false
        }
    }
    bt(0); return result
}
```

### Dry Run — `n = 4`
```
bt(row=0):
  c=0: anti[0]=occ? no. place Q at (0,0). diag[-0]=0, anti[0]=0.
  bt(row=1):
    c=0: cols[0]=true → skip
    c=1: diag[1-1=0]=true → skip (same \ diagonal as (0,0))
    c=2: anti[1+2=3]=? no. place Q at (1,2).
    bt(row=2):
      c=0: diag[2-0=2]? no, but anti[2+0=2]? no. cols[0]=true → skip
      c=1: anti[2+1=3]=true → skip
      c=2: cols[2]=true → skip
      c=3: anti[2+3=5]? no. But diag[2-3=-1]? no. cols[3]? no. place Q at (2,3).
      bt(row=3): all columns blocked → backtrack
    ... continue
  Eventually: (0,1),(1,3),(2,0),(3,2) → solution 1 ✓
              (0,2),(1,0),(2,3),(3,1) → solution 2 ✓
```

---

## Approach 2 — Backtracking with Boolean Arrays (Recommended ✅)

### Intuition
Same algorithm; replace maps with pre-allocated boolean slices for O(1) array access instead of hash map lookup. Diagonal index `r-c` is offset by `n-1` to keep it non-negative (range: `[0, 2n-2]`). Anti-diagonal `r+c` is already non-negative (range: `[0, 2n-2]`).

### Complexity
- **Time:** O(n!).
- **Space:** O(n).

---

## Key Takeaways

- **`r-c` is constant on `\` diagonals** — two queens at `(r1,c1)` and `(r2,c2)` share a `\` diagonal iff `r1-c1 == r2-c2`.
- **`r+c` is constant on `/` anti-diagonals** — same for `r1+c1 == r2+c2`.
- **One queen per row** — the crucial observation that reduces the problem from placing n queens on n² squares to choosing one column per row.
- **No column/row/diagonal re-check needed** — `cols`, `diag`, and `anti` sets encode all attacks; checking these three is sufficient.
- **#52 is #51 + counting** — the same backtracking skeleton; just increment a counter instead of recording the board.

---

## Related Problems

- LeetCode #52 — N-Queens II (count solutions only; add bitmask optimization)
- LeetCode #37 — Sudoku Solver (backtracking with constraint arrays — same pattern)
- LeetCode #46 — Permutations (backtracking template)
