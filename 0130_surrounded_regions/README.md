# 0130 — Surrounded Regions

> LeetCode #130 · Difficulty: Medium
> **Categories:** Array, Depth-First Search, Breadth-First Search, Union Find, Matrix

---

## Problem Statement

Given an `m x n` matrix `board` containing `'X'` and `'O'`, capture all regions that are 4-directionally surrounded by `'X'`.

A region is **captured** by flipping all `'O'`s into `'X'`s in that surrounded region.

**Example 1:**
```
Input:  XXXX        Output: XXXX
        XOOX                XXXX
        XXOX                XXOX  ← wait
        XOXX                XOXX
```
Actually the standard example:
```
Input:  XXXX        Output: XXXX
        XOOX                XXXX
        XXOX                XXXX
        XOXX                XOXX
```

**Constraints:**
- `m == board.length`, `n == board[i].length`
- `1 <= m, n <= 200`
- `board[i][j]` is `'X'` or `'O'`.

---

## Company Frequency

| Company   | Frequency    | Last Reported |
|-----------|--------------|---------------|
| Amazon    | ★★★★☆ High  | 2024          |
| Google    | ★★★★☆ High  | 2024          |
| Microsoft | ★★★☆☆ Medium | 2023          |
| Bloomberg | ★★★☆☆ Medium | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **DFS/BFS from border** — inverse thinking: mark safe cells first → see [`/dsa/backtracking.md`](/dsa/backtracking.md)
- **Temporary marking** — use a sentinel ('S') to distinguish safe O's from capturable O's

---

## Approaches Overview

| # | Approach           | Time   | Space  | When to use         |
|---|--------------------|--------|--------|---------------------|
| 1 | DFS from Border    | O(m·n) | O(m·n) | Elegant recursive   |
| 2 | BFS from Border    | O(m·n) | O(m·n) | Avoids recursion    |

---

## Approach 1 — DFS from Border

### Intuition
**Inverse thinking**: instead of finding enclosed regions, find safe regions (border-connected 'O's). Mark them temporarily as 'S'. Then:
- Remaining 'O' → 'X' (captured).
- 'S' → 'O' (restore safe).

### Algorithm
1. DFS from every 'O' on the border, marking connected 'O's as 'S'.
2. Scan board:
   - 'O' → 'X'.
   - 'S' → 'O'.

### Complexity
- **Time:** O(m·n)
- **Space:** O(m·n) — recursion stack (worst case for DFS).

### Code
```go
func solve(board [][]byte) {
    m, n := len(board), len(board[0])
    var dfs func(r, c int)
    dfs = func(r, c int) {
        if r<0||r>=m||c<0||c>=n||board[r][c]!='O' { return }
        board[r][c]='S'
        dfs(r+1,c); dfs(r-1,c); dfs(r,c+1); dfs(r,c-1)
    }
    for c:=0;c<n;c++ { dfs(0,c); dfs(m-1,c) }
    for r:=0;r<m;r++ { dfs(r,0); dfs(r,n-1) }
    for r:=0;r<m;r++ {
        for c:=0;c<n;c++ {
            if board[r][c]=='O' { board[r][c]='X' } else if board[r][c]=='S' { board[r][c]='O' }
        }
    }
}
```

### Dry Run
```
XXXX        After DFS seeds:    After flip:
XOOX        XXXX                XXXX
XXOX    →   XSSX    →          XXXX
XOXX        XXOX                XXXX
            XOXX                XOXX
```

Wait — `XOXX`: the bottom-left 'O' at row=3,col=1 is on the border → DFS from col 0 of row 3 won't hit it; from row `m-1`, it scans all cols: dfs(3,1) marks it 'S'. So it's restored to 'O' in the output. ✓

---

## Approach 2 — BFS from Border

### Intuition
Same algorithm but uses BFS queue instead of recursion. Avoids potential stack overflow on large boards.

### Code
See `main.go` — `solveBFS`.

### Dry Run
Same result as DFS approach.

---

## Key Takeaways
- **Inverse approach**: don't find captured regions → find safe (border-connected) regions, mark, then flip the rest.
- Temporary sentinel character ('S') avoids extra visited array.
- DFS may stack overflow on large boards (200×200 = 40k deep). BFS is safer for production.

---

## Related Problems
- LeetCode #200 — Number of Islands (DFS/BFS flood fill)
- LeetCode #417 — Pacific Atlantic Water Flow
- LeetCode #695 — Max Area of Island
