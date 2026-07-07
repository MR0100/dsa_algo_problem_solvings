# Dynamic Programming — 2D (Table DP)

> **Category:** Dynamic Programming
> **Related:** [`dynamic_programming.md`](/dsa/dynamic_programming.md) (1D / general DP)

---

## What it is

2D dynamic programming solves a problem by defining a **two-parameter state**
`dp[i][j]` and filling a table so that each cell is computed from a constant
(or small) number of previously computed cells. The two indices almost always
come from one of three sources:

1. **Two sequences** — `i` indexes a prefix of string/array A, `j` a prefix of
   string/array B. `dp[i][j]` answers the question for `A[0..i)` vs `B[0..j)`.
   *Examples: Edit Distance, Regular Expression Matching, Distinct Subsequences.*
2. **A grid** — `(i, j)` is literally a cell of a 2D matrix; the DP walks the
   grid. *Examples: Unique Paths, Minimum Path Sum, Maximal Rectangle, Triangle.*
3. **An interval of one sequence** — `dp[i][j]` describes the substring/subarray
   `s[i..j]` ("interval DP"). *Examples: Longest Palindromic Substring,
   Burst Balloons, Matrix Chain Multiplication.*

The four questions you must answer for any 2D DP, in order:

| Step | Question | Example (Edit Distance) |
|------|----------|-------------------------|
| 1. State | What does `dp[i][j]` **mean**, precisely? | min ops to convert `word1[0..i)` → `word2[0..j)` |
| 2. Transition | How does `dp[i][j]` follow from smaller states? | `min(dp[i-1][j], dp[i][j-1], dp[i-1][j-1]) + 1`, or `dp[i-1][j-1]` on a match |
| 3. Base cases | What cells are known without recursion? | `dp[i][0] = i`, `dp[0][j] = j` |
| 4. Order + answer | What fill order makes dependencies ready? Where is the answer? | row by row; answer at `dp[m][n]` |

---

## How to recognise it — signals in the problem statement

- **Two strings / two arrays** and a question about how they relate:
  "convert one into the other", "is `s3` an interleaving of `s1` and `s2`",
  "does pattern `p` match string `s`", "number of subsequences of `s` equal to
  `t`", "longest common ...". Two independent positions ⇒ two indices ⇒ 2D.
- **A grid / matrix** with movement restricted to a DAG-like direction
  (usually "only right or down"), asking to **count paths**, **min/max path
  cost**, or **largest sub-rectangle/square** satisfying a property.
- **Optimal substructure over substrings**: "longest palindromic substring",
  "min cuts", "score of merging/bursting an interval" — the answer for
  `s[i..j]` depends only on answers for strictly shorter intervals inside it.
- **Counting or optimisation** verbs ("how many ways", "minimum number of
  operations", "maximum value") combined with any of the above. If the problem
  asked to *enumerate all* solutions instead, you'd want backtracking, not DP.
- A brute-force recursion whose call signature naturally has **two changing
  parameters** — that recursion tree has overlapping subproblems, and
  memoising it *is* top-down 2D DP.
- Constraints around `n, m ≤ 1000`–`5000`: an `O(n·m)` table is exactly what
  the setter intends. (`n ≤ 20` hints at bitmask/backtracking instead;
  `n ≥ 10^5` hints the 2D table is too big and a smarter idea is needed.)

Rule of thumb: **count the independent "positions" the subproblem needs to
remember. Two positions → 2D DP.**

---

## General templates (Go)

### Template A — two sequences (prefix vs prefix)

```go
// dp[i][j] = answer for the pair of prefixes a[0..i) and b[0..j).
// Table is (m+1) x (n+1): row/column 0 represent the EMPTY prefix.
func twoSequenceDP(a, b string) int {
    m, n := len(a), len(b)

    // allocate (m+1) x (n+1) table
    dp := make([][]int, m+1)
    for i := range dp {
        dp[i] = make([]int, n+1)
    }

    // ── base cases: one prefix empty ──
    for i := 0; i <= m; i++ {
        dp[i][0] = /* answer when b is empty, e.g. i for edit distance */ i
    }
    for j := 0; j <= n; j++ {
        dp[0][j] = /* answer when a is empty */ j
    }

    // ── fill row by row: every dependency (i-1,*) and (i,j-1) is ready ──
    for i := 1; i <= m; i++ {
        for j := 1; j <= n; j++ {
            if a[i-1] == b[j-1] { // NOTE: dp index i ↔ char index i-1
                dp[i][j] = dp[i-1][j-1] // characters match: extend the diagonal
            } else {
                // combine the three neighbours per the problem's transition
                dp[i][j] = 1 + min(dp[i-1][j-1], min(dp[i-1][j], dp[i][j-1]))
            }
        }
    }
    return dp[m][n] // answer = both full prefixes
}
```

### Template B — grid DP

```go
// dp[i][j] = best/count of ways to reach cell (i, j) moving only right/down.
func gridDP(grid [][]int) int {
    m, n := len(grid), len(grid[0])
    dp := make([][]int, m)
    for i := range dp {
        dp[i] = make([]int, n)
    }

    dp[0][0] = grid[0][0] // start cell

    for i := 0; i < m; i++ {
        for j := 0; j < n; j++ {
            if i == 0 && j == 0 {
                continue // already set
            }
            best := math.MaxInt // (or 0 for counting problems: sum instead of min)
            if i > 0 {
                best = min(best, dp[i-1][j]) // came from above
            }
            if j > 0 {
                best = min(best, dp[i][j-1]) // came from the left
            }
            dp[i][j] = best + grid[i][j]
        }
    }
    return dp[m-1][n-1]
}
```

### Template C — interval DP (fill by increasing length)

```go
// dp[i][j] = answer for the interval s[i..j] (inclusive).
// MUST be filled by increasing interval LENGTH, so that every inner
// interval is already computed when the outer one asks for it.
func intervalDP(s string) int {
    n := len(s)
    dp := make([][]int, n)
    for i := range dp {
        dp[i] = make([]int, n)
        dp[i][i] = /* base case: single character */ 1
    }

    for length := 2; length <= n; length++ { // interval length
        for i := 0; i+length-1 < n; i++ {    // left endpoint
            j := i + length - 1              // right endpoint
            // transition typically peels the endpoints...
            //   dp[i][j] from dp[i+1][j-1] when s[i]==s[j]
            // ...or tries every split point k in (i, j):
            //   dp[i][j] = best over k of combine(dp[i][k], dp[k+1][j])
            _ = j
        }
    }
    return dp[0][n-1] // answer = whole string
}
```

### Template D — top-down memoisation (same state, recursive shape)

```go
// Memoised recursion is often the fastest way to a correct 2D DP in an
// interview: write the brute-force recursion first, then cache (i, j).
func topDown(a, b string) int {
    memo := make([][]int, len(a)+1)
    for i := range memo {
        memo[i] = make([]int, len(b)+1)
        for j := range memo[i] {
            memo[i][j] = -1 // -1 = "not computed yet" (pick a sentinel outside the answer range)
        }
    }
    var solve func(i, j int) int
    solve = func(i, j int) int {
        if i == 0 { return j } // base cases mirror the bottom-up row/col 0
        if j == 0 { return i }
        if memo[i][j] != -1 {
            return memo[i][j] // reuse: this is what kills the exponential blow-up
        }
        var res int
        // ... same transition as bottom-up, but as recursive calls ...
        memo[i][j] = res
        return res
    }
    return solve(len(a), len(b))
}
```

### Space optimisation — rolling rows

When the transition only reads row `i-1` (and the current row), the full
table can be replaced with **two rows** — or **one row plus a saved diagonal**:

```go
// One-row edit distance: prev holds dp[i-1][*] values being overwritten.
prev := make([]int, n+1) // conceptually dp[i-1]
curr := make([]int, n+1) // conceptually dp[i]
for j := 0; j <= n; j++ { prev[j] = j }
for i := 1; i <= m; i++ {
    curr[0] = i
    for j := 1; j <= n; j++ {
        if a[i-1] == b[j-1] {
            curr[j] = prev[j-1]                          // diagonal
        } else {
            curr[j] = 1 + min(prev[j-1], min(prev[j], curr[j-1]))
        }
    }
    prev, curr = curr, prev // swap: O(1), no reallocation
}
// answer in prev[n] (after the final swap)
```

O(m·n) space → O(min(m, n)) space. Mention this in interviews even if you
code the full table first.

---

## Worked example — Edit Distance (LeetCode #72), traced step by step

**Problem:** minimum operations (insert / delete / replace) to convert
`word1 = "horse"` into `word2 = "ros"`.

**State:** `dp[i][j]` = min ops to convert `word1[0..i)` → `word2[0..j)`.

**Transition:**
- `word1[i-1] == word2[j-1]` → `dp[i][j] = dp[i-1][j-1]` (free match)
- else `dp[i][j] = 1 + min(dp[i-1][j-1] /*replace*/, dp[i-1][j] /*delete*/, dp[i][j-1] /*insert*/)`

**Base cases:** `dp[i][0] = i` (delete everything), `dp[0][j] = j` (insert everything).

Full table (rows = `"" h o r s e`, columns = `"" r o s`):

|       | "" | r | o | s |
|-------|----|---|---|---|
| **""**| 0  | 1 | 2 | 3 |
| **h** | 1  | 1 | 2 | 3 |
| **o** | 2  | 2 | 1 | 2 |
| **r** | 3  | 2 | 2 | 2 |
| **s** | 4  | 3 | 3 | 2 |
| **e** | 5  | 4 | 4 | 3 |

Cell-by-cell trace of the interesting cells:

1. `dp[1][1]` (`h` vs `r`): mismatch → `1 + min(dp[0][0]=0, dp[0][1]=1, dp[1][0]=1) = 1` — replace `h`→`r`.
2. `dp[2][2]` (`o` vs `o`): **match** → copy diagonal `dp[1][1] = 1`.
3. `dp[3][1]` (`r` vs `r`): **match** → copy diagonal `dp[2][0] = 2`.
4. `dp[3][2]` (`r` vs `o`): mismatch → `1 + min(dp[2][1]=2, dp[2][2]=1, dp[3][1]=2) = 2`.
5. `dp[4][3]` (`s` vs `s`): match → diagonal `dp[3][2] = 2`.
6. `dp[5][3]` (`e` vs `s`): mismatch → `1 + min(dp[4][2]=3, dp[4][3]=2, dp[5][2]=4) = 3`.

**Answer:** `dp[5][3] = 3` — replace `h`→`r`, delete `r` (of hor**r**se? no — delete `s`... concretely: `horse → rorse → rose → ros`).

Reading an optimal path back off the table (follow which neighbour produced
each cell) recovers the actual edit script — a common follow-up.

---

## Common pitfalls and how to avoid them

1. **Off-by-one between table index and character index.**
   With an `(m+1)×(n+1)` prefix table, `dp[i][j]` talks about *lengths* `i, j`,
   so the characters being compared are `a[i-1]` and `b[j-1]`. Write the state
   definition as a comment ("`dp[i][j]` = answer for first `i` chars of a...")
   and mechanically derive indices from it. Mixing "index" and "length"
   conventions mid-function is the #1 bug source.

2. **Forgetting (or wrongly initialising) row 0 / column 0.**
   The empty-prefix base cases carry real information — in Regular Expression
   Matching, `dp[0][j]` must be true for patterns like `a*b*` that can match
   the empty string; leaving it all-false silently breaks every `*` transition.
   Always ask: *"what is the answer when one input is empty?"* and code it
   explicitly.

3. **Wrong fill order for the dependencies.**
   Row-major order works for prefix/grid DP because dependencies point up/left.
   Interval DP **must** iterate by increasing length (Template C) — iterating
   `i` then `j` naively reads `dp[i+1][j-1]` before it exists. If unsure,
   write top-down memoisation instead: recursion computes dependencies on
   demand, so ordering bugs are impossible.

4. **Aliased rows when allocating in Go.**
   `dp := make([][]int, m+1)` allocates only the outer slice — each `dp[i]`
   still needs its own `make([]int, n+1)`. Copy-pasting a "2D make" that
   shares one backing row makes every row alias the same memory and produces
   baffling wrong answers.

5. **State that doesn't actually capture the subproblem.**
   If two different situations map to the same `(i, j)` but have different
   answers, the state is under-specified — add a dimension (e.g. stock
   problems need `dp[day][holding]`; Scramble String needs a length, giving
   `dp[len][i][j]`). Symptom: the recurrence "needs to know something extra"
   about how you got there.

6. **Premature space optimisation.**
   Rolling arrays are easy to get wrong (the diagonal `dp[i-1][j-1]` gets
   overwritten before it's read — save it in a temp, or iterate `j` in the
   right direction). Get the full table correct first, then optimise, and keep
   the full-table version around to diff against.

7. **Reconstructing the answer, not just its value.**
   If the follow-up asks for the actual path/edit script/subsequence, either
   keep parent pointers or walk the finished table backwards re-checking which
   transition was taken. Don't try to build the answer string during the fill.

8. **Assuming O(m·n) is affordable.**
   `m = n = 10^5` means `10^10` cells — the 2D formulation is the *model*, not
   the algorithm. Look for monotonicity, divide & conquer optimisation, or a
   1D reformulation before coding.

---

## Problems in this repo

Two-sequence DP:

- [0010 — Regular Expression Matching](/0010_regular_expression_matching/README.md) — `dp[i][j]`: does `p[0..j)` match `s[0..i)`; `*` transition reads both up and left.
- [0044 — Wildcard Matching](/0044_wildcard_matching/README.md) — same shape with simpler `*` semantics.
- [0072 — Edit Distance](/0072_edit_distance/README.md) — the canonical worked example above.
- [0097 — Interleaving String](/0097_interleaving_string/README.md) — `dp[i][j]`: can `s1[0..i)` + `s2[0..j)` interleave into `s3[0..i+j)`.
- [0115 — Distinct Subsequences](/0115_distinct_subsequences/README.md) — counting variant: sum transitions instead of min/max.

Grid DP:

- [0062 — Unique Paths](/0062_unique_paths/README.md) — count paths, `dp[i][j] = dp[i-1][j] + dp[i][j-1]`.
- [0063 — Unique Paths II](/0063_unique_paths_ii/README.md) — same with obstacle cells forced to 0.
- [0064 — Minimum Path Sum](/0064_minimum_path_sum/README.md) — min-cost path, the optimisation twin of #62.
- [0120 — Triangle](/0120_triangle/README.md) — grid DP on a triangular grid; clean bottom-up + 1D rolling row.
- [0085 — Maximal Rectangle](/0085_maximal_rectangle/README.md) — per-row DP arrays (left/right/height) over a 2D grid.

Interval / other 2D-table DP:

- [0005 — Longest Palindromic Substring](/0005_longest_palindromic_substring/README.md) — interval DP: `dp[i][j]` = is `s[i..j]` a palindrome, filled by length.
- [0087 — Scramble String](/0087_scramble_string/README.md) — state needs a third dimension (`dp[len][i][j]`): a worked example of pitfall #5.

> Problems 0131–0400 are being added concurrently; a later pass will extend
> this list (e.g. Longest Common Subsequence #1143, Maximal Square #221,
> Burst Balloons #312).
