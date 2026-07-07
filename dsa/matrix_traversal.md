# Matrix Traversal

> **Category:** Arrays / 2-D Grids
> **Difficulty of mastery:** Easy fundamentals, Medium–Hard variations
> **Prerequisites:** 1-D array indexing, nested loops

---

## 1. What the concept is

**Matrix traversal** is the family of techniques for visiting the cells of a
2-D grid (`m` rows × `n` columns) in some prescribed order, or for moving
between cells along allowed directions. It is the 2-D generalisation of
iterating over an array, but the extra dimension unlocks many distinct
traversal orders, each of which is its own reusable pattern:

| Pattern | Visit order | Classic problem |
|---|---|---|
| **Row-major scan** | left→right, top→bottom | Set Matrix Zeroes, Valid Sudoku |
| **Column-major scan** | top→bottom, left→right | Search a 2D Matrix II |
| **Spiral order** | outer ring → inner ring, clockwise | Spiral Matrix I / II |
| **Diagonal order** | along anti-/main diagonals | Diagonal Traverse |
| **Boundary / layer** | one concentric ring at a time | Rotate Image |
| **Zigzag / snake** | alternate direction each row | Some level-order variants |
| **Directional walk** | move via a `(dr, dc)` direction array | Word Search, Surrounded Regions, island problems |
| **Staircase search** | start at a corner, prune a row/column per step | Search a 2D Matrix (sorted grids) |
| **DP sweep** | visit in dependency order (usually row-major) | Unique Paths, Minimum Path Sum, Maximal Rectangle |

The unifying skill: **map a 2-D position to loop variables + direction
vectors, and keep the visit order and boundaries airtight.**

---

## 2. How to recognise it — signals in the problem statement

Reach for matrix traversal when you see:

- The input is an `m x n` grid / matrix / board / image / maze:
  *"given an `m x n` matrix"*, *"an `n x n` 2D matrix representing an image"*,
  *"a 2D board"*, *"a grid of `'0'`s and `'1'`s"*.
- The required output is a **specific visit order**: *"return all elements in
  spiral order"*, *"diagonal order"*, *"level by level"*.
- An **in-place geometric transformation**: *"rotate the image by 90 degrees
  (in-place)"*, *"transpose"*, *"reflect"*.
- **Movement rules** are stated: *"you may move only down or right"*,
  *"adjacent cells are horizontally or vertically neighboring"*,
  *"8-directionally connected"*.
- **Connectivity / regions**: *"number of islands"*, *"surrounded regions"*,
  *"flood fill"* — these are matrix traversal + DFS/BFS.
- **Sorted-grid search**: *"each row is sorted; the first integer of each row
  is greater than the last integer of the previous row"* → binary search or
  staircase traversal instead of a full scan.
- **Sub-rectangle questions**: *"largest rectangle/square of 1s"* → traversal
  combined with DP or a stack per row.

Rule of thumb: the moment coordinates come in **pairs** `(row, col)`, decide
*which traversal order* the problem is really asking for before writing code.

---

## 3. General templates (Go)

### 3.1 Row-major scan (the default)

```go
// Visit every cell once, top-left to bottom-right.
// Time: O(m*n), Space: O(1).
for r := 0; r < len(matrix); r++ {          // each row, top to bottom
    for c := 0; c < len(matrix[0]); c++ {   // each column, left to right
        process(matrix[r][c])               // matrix[r][c] is the current cell
    }
}
```

### 3.2 Direction arrays — the workhorse for neighbours

```go
// dr/dc encode the 4 orthogonal moves: up, right, down, left.
// For 8-directional problems add the four diagonals.
var dr = []int{-1, 0, 1, 0}
var dc = []int{0, 1, 0, -1}

func neighbors(matrix [][]int, r, c int) {
    m, n := len(matrix), len(matrix[0])
    for d := 0; d < 4; d++ {          // try each direction
        nr, nc := r+dr[d], c+dc[d]    // candidate neighbour
        if nr < 0 || nr >= m || nc < 0 || nc >= n {
            continue                  // out of bounds — skip, never index
        }
        process(matrix[nr][nc])       // safe to visit now
    }
}
```

Pseudocode for grid DFS (islands / regions / word search):

```text
dfs(r, c):
    if (r, c) out of bounds OR visited OR not matching: return
    mark (r, c) visited            # in-place sentinel or visited[][] array
    for each (dr, dc) in directions:
        dfs(r + dr, c + dc)
```

### 3.3 Spiral traversal — four shrinking boundaries

```go
// spiralOrder visits an m x n matrix in clockwise spiral order.
// Four boundary pointers shrink inward after each side is consumed.
// Time: O(m*n), Space: O(1) extra (excluding output).
func spiralOrder(matrix [][]int) []int {
    if len(matrix) == 0 {
        return nil
    }
    top, bottom := 0, len(matrix)-1      // first and last unvisited row
    left, right := 0, len(matrix[0])-1   // first and last unvisited column
    result := make([]int, 0, len(matrix)*len(matrix[0]))

    for top <= bottom && left <= right {
        for c := left; c <= right; c++ { // 1) top row, left → right
            result = append(result, matrix[top][c])
        }
        top++                            // top row consumed

        for r := top; r <= bottom; r++ { // 2) right column, top → bottom
            result = append(result, matrix[r][right])
        }
        right--                          // right column consumed

        if top <= bottom {               // guard: a row may remain
            for c := right; c >= left; c-- { // 3) bottom row, right → left
                result = append(result, matrix[bottom][c])
            }
            bottom--
        }
        if left <= right {               // guard: a column may remain
            for r := bottom; r >= top; r-- { // 4) left column, bottom → top
                result = append(result, matrix[r][left])
            }
            left++
        }
    }
    return result
}
```

### 3.4 In-place 90° rotation = transpose + reverse rows

```go
// rotate turns the matrix 90° clockwise in place.
// Clockwise 90° == transpose (mirror over main diagonal) then reverse each row.
// Time: O(n^2), Space: O(1).
func rotate(matrix [][]int) {
    n := len(matrix)
    for r := 0; r < n; r++ {
        for c := r + 1; c < n; c++ { // c starts at r+1: touch each pair once
            matrix[r][c], matrix[c][r] = matrix[c][r], matrix[r][c]
        }
    }
    for r := 0; r < n; r++ {
        for i, j := 0, n-1; i < j; i, j = i+1, j-1 { // two-pointer reverse
            matrix[r][i], matrix[r][j] = matrix[r][j], matrix[r][i]
        }
    }
}
```

(Counter-clockwise = transpose + reverse each **column**, or reverse rows first.)

### 3.5 Staircase search in a row- and column-sorted matrix

```go
// searchMatrix starts at the top-right corner; every step eliminates
// one full row or one full column. Time: O(m+n), Space: O(1).
func searchMatrix(matrix [][]int, target int) bool {
    r, c := 0, len(matrix[0])-1        // top-right corner
    for r < len(matrix) && c >= 0 {
        switch {
        case matrix[r][c] == target:
            return true
        case matrix[r][c] > target:
            c--                        // everything below in this column is bigger
        default:
            r++                        // everything left in this row is smaller
        }
    }
    return false
}
```

### 3.6 Flattened 1-D view (for strictly sorted matrices / binary search)

```go
// A matrix where rows are sorted and each row continues the previous one
// is just a sorted array of length m*n in disguise.
// index → cell:  r = idx / n,  c = idx % n
lo, hi := 0, m*n-1
for lo <= hi {
    mid := lo + (hi-lo)/2
    val := matrix[mid/n][mid%n]   // decode 1-D index back to (row, col)
    // ... standard binary search on val
}
```

### 3.7 Diagonal indexing

```go
// All cells on the same anti-diagonal share r + c.
// All cells on the same main diagonal share r - c.
for d := 0; d < m+n-1; d++ {          // there are m+n-1 anti-diagonals
    for r := max(0, d-n+1); r <= min(d, m-1); r++ {
        c := d - r                    // recover the column from the diagonal id
        process(matrix[r][c])
    }
}
```

---

## 4. Worked example — Spiral Matrix (LeetCode 54), traced step by step

Input:

```
matrix = [[1, 2, 3],
          [4, 5, 6],
          [7, 8, 9]]
```

Expected output: `[1, 2, 3, 6, 9, 8, 7, 4, 5]`

Initial state: `top=0, bottom=2, left=0, right=2, result=[]`.

| Step | Action | Cells appended | Boundaries after | result |
|---|---|---|---|---|
| 1 | Top row `top=0`, `c: left→right` (0→2) | `1, 2, 3` | `top=1` | `[1 2 3]` |
| 2 | Right col `right=2`, `r: top→bottom` (1→2) | `6, 9` | `right=1` | `[1 2 3 6 9]` |
| 3 | Guard `top(1) <= bottom(2)` ✓ → bottom row `bottom=2`, `c: right→left` (1→0) | `8, 7` | `bottom=1` | `[1 2 3 6 9 8 7]` |
| 4 | Guard `left(0) <= right(1)` ✓ → left col `left=0`, `r: bottom→top` (1→1) | `4` | `left=1` | `[1 2 3 6 9 8 7 4]` |
| 5 | Loop check: `top(1) <= bottom(1)` ✓ and `left(1) <= right(1)` ✓ → second lap. Top row `top=1`, `c: 1→1` | `5` | `top=2` | `[1 2 3 6 9 8 7 4 5]` |
| 6 | Right col: `r: top(2)→bottom(1)` — empty loop | — | `right=0` | unchanged |
| 7 | Guard `top(2) <= bottom(1)` ✗ → skip bottom row. Guard `left(1) <= right(0)` ✗ → skip left col | — | — | unchanged |
| 8 | Loop check: `top(2) <= bottom(1)` ✗ → **terminate** | | | |

Output: `[1, 2, 3, 6, 9, 8, 7, 4, 5]` ✓ — all 9 cells visited exactly once,
and step 7 shows exactly why the two mid-loop guards exist: without them, a
single leftover row or column would be traversed twice.

---

## 5. Common pitfalls and how to avoid them

1. **Row/column index swap.** `matrix[r][c]` — `r` selects the row (vertical
   position), `c` the column (horizontal). Mixing them up compiles fine on
   square matrices and explodes on rectangular ones. **Fix:** always name loop
   variables `r`/`c` (never `i`/`j`), and test on a non-square input like
   `2x3` before claiming done.

2. **Missing the two spiral guards.** After consuming the top row and right
   column, a single remaining row or column would be walked twice without the
   `top <= bottom` / `left <= right` checks. **Fix:** memorise the guards as
   part of the template; test with `1xN` and `Mx1` matrices.

3. **Bounds check after indexing instead of before.** `matrix[nr][nc]` with
   `nr == -1` panics in Go. **Fix:** the boundary condition is always the
   *first* line of the DFS / neighbour loop.

4. **Forgetting to mark visited (infinite DFS).** Two adjacent cells call DFS
   on each other forever. **Fix:** mark the cell *before* recursing —
   either a `visited [][]bool` or an in-place sentinel (e.g. flip `'1'`→`'#'`);
   restore the sentinel afterwards if the problem needs the grid intact
   (Word Search does).

5. **Empty-matrix / degenerate inputs.** `len(matrix) == 0` or
   `len(matrix[0]) == 0`, single row, single column, `1x1`. **Fix:** guard at
   the top; include these cases in `main()`.

6. **In-place rotation with a naive 4-way swap and wrong loop limits.** The
   layer-by-layer swap is easy to off-by-one. **Fix:** prefer
   transpose + reverse — two dead-simple loops, same O(1) space.

7. **Transpose touching each pair twice.** Looping `c` from `0` swaps every
   pair twice, restoring the original matrix. **Fix:** inner loop starts at
   `c = r + 1`.

8. **Using first-row/first-column as marker storage but clobbering the
   flags** (Set Matrix Zeroes). The cell `matrix[0][0]` serves two roles.
   **Fix:** keep a separate boolean for one of the two (e.g. `firstColZero`)
   and process the first row/column *last*.

9. **Copying Go slices of slices shallowly.** `dst := matrix` or
   `copy(dst, matrix)` copies row *headers*, not row contents — mutating the
   "copy" mutates the original. **Fix:** allocate and `copy` row by row.

10. **Wrong traversal order for DP dependencies.** `dp[r][c]` depending on
    `dp[r-1][c]` and `dp[r][c-1]` requires row-major order; a
    bottom-up dependency (Triangle) requires the reverse. **Fix:** write the
    recurrence first, then choose the sweep direction that guarantees
    dependencies are already computed.

11. **O(m·n) extra memory where O(1) was demanded.** Many matrix problems
    (Rotate Image, Set Matrix Zeroes) explicitly say "in place". Present the
    extra-space version as a first approach, then the in-place optimal.

---

## 6. Problems in this repo

Matrix traversal appears here in progressively richer forms:

- [`0036_valid_sudoku`](../0036_valid_sudoku/README.md) — row-major scan with simultaneous row / column / 3×3-box bookkeeping (`box = (r/3)*3 + c/3`).
- [`0037_sudoku_solver`](../0037_sudoku_solver/README.md) — grid scan + backtracking over empty cells.
- [`0048_rotate_image`](../0048_rotate_image/README.md) — in-place 90° rotation: transpose + reverse rows, and layer-by-layer swap.
- [`0054_spiral_matrix`](../0054_spiral_matrix/README.md) — the canonical four-boundary spiral traversal (traced above).
- [`0059_spiral_matrix_ii`](../0059_spiral_matrix_ii/README.md) — same spiral template, writing values instead of reading them.
- [`0062_unique_paths`](../0062_unique_paths/README.md) — row-major DP sweep over an implicit grid (moves: right / down).
- [`0063_unique_paths_ii`](../0063_unique_paths_ii/README.md) — same DP sweep with obstacle cells zeroing paths.
- [`0064_minimum_path_sum`](../0064_minimum_path_sum/README.md) — grid DP: cost accumulation in traversal (dependency) order.
- [`0073_set_matrix_zeroes`](../0073_set_matrix_zeroes/README.md) — two-pass scan using the first row/column as O(1) marker storage.
- [`0074_search_a_2d_matrix`](../0074_search_a_2d_matrix/README.md) — flattened 1-D binary search over a fully sorted grid.
- [`0079_word_search`](../0079_word_search/README.md) — directional DFS with in-place visited sentinel and restore.
- [`0085_maximal_rectangle`](../0085_maximal_rectangle/README.md) — row-by-row traversal building histogram heights, reducing 2-D to repeated 1-D.
- [`0097_interleaving_string`](../0097_interleaving_string/README.md) — 2-D DP table traversed row-major (grid-walk formulation of string interleaving).
- [`0120_triangle`](../0120_triangle/README.md) — triangular grid DP traversed bottom-up.
- [`0130_surrounded_regions`](../0130_surrounded_regions/README.md) — boundary-first DFS/BFS from border cells, then a full-grid rewrite pass.

> Problems 0131+ are being added continuously; a later pass will extend this list.

---

## 7. Key takeaways

- Decide the **visit order** before writing code — scan, spiral, layer,
  diagonal, directional walk, or DP-dependency order.
- The `(dr, dc)` **direction array + bounds-check-first** idiom covers every
  neighbour-based grid problem.
- **Spiral = four shrinking boundaries + two guards.**
- **Rotate = transpose + reverse** beats hand-rolled 4-way swaps.
- Sorted grids invite **O(m+n) staircase** or **O(log mn) flattened binary
  search** — never a full O(m·n) scan.
- Always test on **rectangular, single-row, single-column, and 1×1** inputs.
