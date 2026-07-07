# 0079 — Word Search

> LeetCode #79 · Difficulty: Medium
> **Categories:** Array, Backtracking, Depth-First Search, Matrix

---

## Problem Statement

Given an `m x n` grid of characters `board` and a string `word`, return `true` if `word` exists in the grid.

The word can be constructed from letters of sequentially adjacent cells, where adjacent cells are horizontally or vertically neighboring. The same letter cell may not be used more than once.

**Example 1:**
```
Input: board = [["A","B","C","E"],
                ["S","F","C","S"],
                ["A","D","E","E"]], word = "ABCCED"
Output: true
```

**Example 2:**
```
Input: board = [["A","B","C","E"],
                ["S","F","C","S"],
                ["A","D","E","E"]], word = "SEE"
Output: true
```

**Example 3:**
```
Input: board = [["A","B","C","E"],
                ["S","F","C","S"],
                ["A","D","E","E"]], word = "ABCB"
Output: false
```

**Constraints:**
- `m == board.length`
- `n == board[i].length`
- `1 <= m, n <= 6`
- `1 <= word.length <= 15`
- `board` and `word` consists of only lowercase and uppercase English letters.

**Follow-up:** Could you use search pruning to make your solution faster with a larger board?

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Facebook  | ★★★★☆ High      | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Microsoft | ★★★☆☆ Medium    | 2024          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Backtracking / DFS on Grid** — try each direction; undo (restore) on failure. See [`/dsa/backtracking.md`](/dsa/backtracking.md)
- **In-Place Visited Marking** — XOR with 255 to mark/unmark without an extra visited matrix.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Backtracking DFS (XOR mark) | O(m×n×4^L) | O(L) | Only practical approach for this problem |

---

## Approach 1 — Backtracking DFS

### Intuition
Try to start the word from every cell. For each starting cell that matches `word[0]`, launch a DFS that matches subsequent characters by moving to adjacent cells. Mark each visited cell temporarily (to prevent reuse on the same path) and restore it when backtracking.

**In-place visited marking:** XOR the cell's byte value with 255. This corrupts the letter so it can no longer match any target character. XOR again to restore the original byte. This avoids allocating a separate `m×n` visited array.

### Algorithm
1. `dfs(r, c, idx)`:
   - If `idx == len(word)`: return `true` (all characters matched).
   - If out of bounds or `board[r][c] != word[idx]`: return `false`.
   - `board[r][c] ^= 255` (mark visited).
   - Try all 4 directions: `(r±1, c)`, `(r, c±1)`.
   - `board[r][c] ^= 255` (restore — backtrack).
   - Return `false` if no direction worked.
2. For each cell `(r, c)`: if `dfs(r, c, 0)` returns `true`, return `true`.
3. Return `false`.

### Complexity
- **Time:** O(m × n × 4^L) — m×n starting cells, each spawning a DFS of depth L with up to 4 branches (3 effective since we came from one direction).
- **Space:** O(L) — recursion stack depth equals word length L.

### Code
```go
func exist(board [][]byte, word string) bool {
    m, n := len(board), len(board[0])
    dr := []int{0, 0, 1, -1}
    dc := []int{1, -1, 0, 0}

    var dfs func(r, c, idx int) bool
    dfs = func(r, c, idx int) bool {
        if idx == len(word) {
            return true
        }
        if r < 0 || r >= m || c < 0 || c >= n || board[r][c] != word[idx] {
            return false
        }
        board[r][c] ^= 255
        for d := 0; d < 4; d++ {
            if dfs(r+dr[d], c+dc[d], idx+1) {
                board[r][c] ^= 255
                return true
            }
        }
        board[r][c] ^= 255
        return false
    }

    for r := 0; r < m; r++ {
        for c := 0; c < n; c++ {
            if dfs(r, c, 0) {
                return true
            }
        }
    }
    return false
}
```

### Dry Run (Example 1: word="ABCCED")

Board:
```
A B C E
S F C S
A D E E
```

Start at (0,0)='A' matches word[0]='A': DFS begins. Mark (0,0).
- (0,1)='B' matches word[1]='B': Mark (0,1).
  - (0,2)='C' matches word[2]='C': Mark (0,2).
    - (1,2)='C' matches word[3]='C': Mark (1,2).
      - (2,2)='E' matches word[4]='E': Mark (2,2).
        - (2,1)='D' matches word[5]='D': Mark (2,1). idx=6==len → return true!

Path traced: (0,0)→(0,1)→(0,2)→(1,2)→(2,2)→(2,1) = "ABCCED" ✓

---

## Key Takeaways
- XOR with 255 is an elegant O(1) space trick to mark visited cells in-place.
- Restore (`XOR` again) must happen both on success path (before returning true) and on failure path (after all directions fail).
- For boards where the starting character is rare, most DFS calls abort at depth 0 — very fast in practice.
- Pruning addition: before starting, if any required character doesn't appear in the board enough times, return false early.

---

## Related Problems
- LeetCode #212 — Word Search II (find multiple words using Trie + backtracking)
- LeetCode #200 — Number of Islands (DFS on grid)
- LeetCode #130 — Surrounded Regions (DFS/BFS on grid)
