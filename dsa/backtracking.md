# Backtracking

> Systematic trial-and-error: build a candidate solution one **choice** at a
> time, and the moment a choice can no longer lead to a valid solution, **undo
> it** (backtrack) and try the next option. It is DFS over an implicit
> **decision tree** whose nodes are partial solutions.

---

## 1. What it is

Backtracking is an exhaustive-search technique with two upgrades over naive
brute force:

1. **Incremental construction** — instead of generating every complete
   candidate and then checking it, we extend a *partial* candidate one step at
   a time.
2. **Pruning** — the moment a partial candidate violates a constraint (or
   provably cannot be completed), we abandon the *entire subtree* rooted at it.
   One pruned node can eliminate exponentially many complete candidates.

The mental model is a tree walk:

```
                        []                 ← empty partial solution (root)
             ┌──────────┼──────────┐
            [1]        [2]        [3]      ← after 1st choice
          ┌──┴──┐    ┌──┴──┐    ┌──┴──┐
        [1,2] [1,3][2,1] [2,3][3,1] [3,2]  ← after 2nd choice
          │     │    │     │    │     │
        [1,2,3] …  [2,1,3] …  [3,1,2] …    ← leaves = complete solutions
```

Every backtracking algorithm answers three questions:

| Question | Meaning |
|----------|---------|
| **State** | What does a partial solution look like? (`path`, plus bookkeeping like `used[]`, `start` index, remaining target) |
| **Choices** | From the current state, what are the legal next moves? |
| **Goal / prune** | When is the state a complete answer? When is it hopeless and should be abandoned? |

The signature move is **choose → explore → unchoose**: mutate the shared
state, recurse, then restore the state *exactly* as it was so the next sibling
choice starts clean.

---

## 2. How to recognise a backtracking problem

Signals in the problem statement:

- **"Return ALL …"** — *all* combinations / permutations / subsets /
  partitions / valid arrangements. Enumeration of an exponential answer space
  is the #1 giveaway. (Contrast: "return the *count*" or "return the *best*"
  often means DP instead — see pitfalls.)
- **Small constraints** — `n ≤ 20`-ish, string length ≤ 16, board 9×9.
  Exponential (`2^n`, `n!`) search is only feasible when the input is tiny;
  problem setters signal backtracking by keeping `n` small.
- **"Generate" / "construct" wording** — generate parentheses, restore IP
  addresses, solve the Sudoku, place N queens.
- **Constraint-satisfaction flavour** — a set of slots must be filled subject
  to rules that interact (Sudoku rows/cols/boxes, queens attacking each
  other). No greedy or DP ordering works because a choice's validity depends
  on *all* previous choices.
- **Path finding with "used" state on a grid/graph** — Word Search style: a
  cell may not be revisited *within the current path*, so visited-state must
  be undone on return (plain DFS with a permanent `visited` set is wrong).
- **Decision-per-element structure** — every element is either
  included/excluded (subsets), or gets one of `k` labels (letter combinations:
  one letter per digit).

If the answer is a single number/optimum and constraints are large
(`n ≤ 10^5`), think DP/greedy, not backtracking.

---

## 3. General templates (Go)

### 3.1 Core skeleton

```go
// res collects complete solutions; path is the shared, mutable partial solution.
var res [][]int
var path []int

func backtrack(state /* whatever defines "where we are" */) {
    // 1. GOAL: is path a complete solution?
    if isComplete(state) {
        // Copy! path's backing array will be mutated after we return.
        res = append(res, append([]int(nil), path...))
        return
    }

    // 2. CHOICES: iterate over every legal next move from this state.
    for _, c := range choices(state) {
        if !valid(c, state) { // 3. PRUNE: skip doomed branches early
            continue
        }
        path = append(path, c)     // CHOOSE   – mutate shared state
        markUsed(c)                //            (any extra bookkeeping)
        backtrack(next(state, c))  // EXPLORE  – recurse one level deeper
        unmarkUsed(c)              // UNCHOOSE – restore state exactly
        path = path[:len(path)-1]  //            (reverse every mutation)
    }
}
```

### 3.2 Subsets / combinations — the `start` index

Use a `start` index so each element is considered only *after* the previous
pick → no duplicate sets like `{1,2}` and `{2,1}`.

```go
// subsets enumerates all 2^n subsets of nums.
func subsets(nums []int) [][]int {
    res, path := [][]int{}, []int{}
    var bt func(start int)
    bt = func(start int) {
        // Every node is an answer for subsets (record at every node, not just leaves).
        res = append(res, append([]int(nil), path...))
        for i := start; i < len(nums); i++ {
            path = append(path, nums[i]) // choose nums[i]
            bt(i + 1)                    // only elements after i may follow
            path = path[:len(path)-1]    // unchoose
        }
    }
    bt(0)
    return res
}
```

Variants:
- **Combination Sum (#39)** — element reusable: recurse with `bt(i)` instead
  of `bt(i+1)`; prune when `remaining < 0` (sort first to `break` early).
- **Combinations (#77)** — fixed size `k`: goal test `len(path)==k`; prune
  when not enough elements remain: `if len(nums)-i < k-len(path) { break }`.

### 3.3 Permutations — the `used[]` array

Order matters, every element appears exactly once → track which are taken.

```go
// permute enumerates all n! orderings of nums.
func permute(nums []int) [][]int {
    res, path := [][]int{}, []int{}
    used := make([]bool, len(nums))
    var bt func()
    bt = func() {
        if len(path) == len(nums) { // goal: path holds every element
            res = append(res, append([]int(nil), path...))
            return
        }
        for i, v := range nums {
            if used[i] { // choice already consumed on this path
                continue
            }
            used[i] = true
            path = append(path, v)
            bt()
            path = path[:len(path)-1]
            used[i] = false
        }
    }
    bt()
    return res
}
```

### 3.4 Duplicate elements — sort + same-level skip

For Subsets II (#90), Combination Sum II (#40), Permutations II (#47): sort
first, then **skip a value if it equals the previous value and the previous
value was not chosen on this path** — i.e. the duplicate is being tried as a
*sibling* (same tree level), which would rebuild an identical subtree.

```go
sort.Ints(nums)
for i := start; i < len(nums); i++ {
    // nums[i] == nums[i-1] and nums[i-1] was skipped at THIS level →
    // choosing nums[i] now would duplicate the branch that started with nums[i-1].
    if i > start && nums[i] == nums[i-1] {
        continue
    }
    // ... choose / explore / unchoose
}
// Permutations II uses the used[] form of the same guard:
// if i > 0 && nums[i] == nums[i-1] && !used[i-1] { continue }
```

### 3.5 Grid backtracking — mark/unmark cells (Word Search #79)

```go
func exist(board [][]byte, word string) bool {
    m, n := len(board), len(board[0])
    var bt func(r, c, k int) bool
    bt = func(r, c, k int) bool {
        if r < 0 || r >= m || c < 0 || c >= n || board[r][c] != word[k] {
            return false // out of bounds, mismatch, or cell already in path
        }
        if k == len(word)-1 {
            return true // matched the whole word
        }
        saved := board[r][c]
        board[r][c] = '#' // CHOOSE: mark cell as used in the current path
        found := bt(r+1, c, k+1) || bt(r-1, c, k+1) ||
            bt(r, c+1, k+1) || bt(r, c-1, k+1)
        board[r][c] = saved // UNCHOOSE: cell reusable by other paths
        return found
    }
    for r := 0; r < m; r++ {
        for c := 0; c < n; c++ {
            if bt(r, c, 0) {
                return true
            }
        }
    }
    return false
}
```

Note the **found-one-answer** shape: return `bool` and short-circuit instead
of collecting into `res` (also used by Sudoku Solver #37).

### 3.6 Constraint satisfaction with O(1) checks (N-Queens #51/#52)

Keep hash sets / boolean arrays / bitmasks describing the constraints, so
`valid()` is O(1) instead of rescanning the partial board:

```go
cols  := make([]bool, n)     // column c occupied?
diag1 := make([]bool, 2*n-1) // "/" diagonal: r + c constant
diag2 := make([]bool, 2*n-1) // "\" diagonal: r - c + n - 1 constant

var bt func(r int)
bt = func(r int) {
    if r == n { count++; return } // all rows filled → solution
    for c := 0; c < n; c++ {
        if cols[c] || diag1[r+c] || diag2[r-c+n-1] {
            continue // square attacked → prune
        }
        cols[c], diag1[r+c], diag2[r-c+n-1] = true, true, true
        bt(r + 1)                                   // one queen per row by construction
        cols[c], diag1[r+c], diag2[r-c+n-1] = false, false, false
    }
}
```

---

## 4. Worked example — Permutations of `[1, 2, 3]`

Using the §3.3 template. `path=[]`, `used=[F F F]`.

| Step | Action | path | used | res |
|------|--------|------|------|-----|
| 1 | choose 1 | `[1]` | `[T F F]` | — |
| 2 | choose 2 | `[1 2]` | `[T T F]` | — |
| 3 | choose 3 → goal | `[1 2 3]` | `[T T T]` | `{[1 2 3]}` |
| 4 | unchoose 3, unchoose 2 | `[1]` | `[T F F]` | |
| 5 | choose 3 | `[1 3]` | `[T F T]` | — |
| 6 | choose 2 → goal | `[1 3 2]` | `[T T T]` | `{…, [1 3 2]}` |
| 7 | unwind back to root (unchoose 2, 3, 1) | `[]` | `[F F F]` | |
| 8 | choose 2 | `[2]` | `[F T F]` | — |
| 9 | choose 1, choose 3 → goal | `[2 1 3]` | `[T T T]` | `{…, [2 1 3]}` |
| 10 | backtrack; choose 3, choose 1 → goal | `[2 3 1]` | `[T T T]` | `{…, [2 3 1]}` |
| 11 | unwind to root; choose 3 | `[3]` | `[F F T]` | — |
| 12 | choose 1, choose 2 → goal | `[3 1 2]` | `[T T T]` | `{…, [3 1 2]}` |
| 13 | backtrack; choose 2, choose 1 → goal | `[3 2 1]` | `[T T T]` | all 6 found |

Key observations:

- Step 4 is the essence: after emitting `[1 2 3]` we *pop 3 and 2* so that the
  loop at the `[1]` level can try `3` as the second element (step 5). The
  shared `path`/`used` state is identical on re-entry to a level as on first
  entry — that invariant is what "unchoose must mirror choose" guarantees.
- The trace is a pre-order DFS of the 3-level decision tree; leaves (6 of
  them = 3!) are exactly the answers.
- Each answer is **copied** into `res` (step 3's append copies `path`); if we
  appended `path` itself, steps 4–5 would corrupt the stored answer.

Complexity: O(n · n!) time (n! leaves, O(n) to copy each answer), O(n)
auxiliary space for recursion + path (output excluded).

---

## 5. Common pitfalls

1. **Appending `path` without copying.** `res = append(res, path)` stores a
   slice header sharing the backing array that later `append`/truncate calls
   mutate → all stored "answers" end up garbled or identical.
   Fix: `res = append(res, append([]int(nil), path...))`.
2. **Asymmetric choose/unchoose.** Every mutation made before the recursive
   call (path push, `used[i]=true`, grid mark, remaining -= v, open++) must be
   reversed after it, in reverse order. One forgotten unmark poisons every
   sibling branch. Grid version: restore the *original* character, don't
   assume it.
3. **Wrong duplicate handling.** Using a set of final answers to dedupe
   "works" but still explores the duplicate subtrees (exponential waste).
   Correct: sort + skip duplicates **at the same tree level**
   (`i > start && nums[i]==nums[i-1]`, or `!used[i-1]` for permutations).
   Also don't over-skip: the guard must allow the duplicate when it *extends*
   the previous copy (different level).
4. **Forgetting `start` (or misusing it).** Combinations/subsets without a
   `start` index generate every ordering of every set. Passing `bt(start+1)`
   instead of `bt(i+1)` reuses skipped elements; `bt(i)` vs `bt(i+1)` is
   exactly the reuse-allowed vs use-once distinction (#39 vs #40).
5. **Pruning too late (or not at all).** Validate *before* recursing, not at
   the leaf. Sort candidates so a failed bound check can `break` the whole
   loop rather than `continue`. In Generate Parentheses, checking validity
   only at length `2n` is 2^(2n) work; pruning with `open < n` / `close < open`
   collapses it to the Catalan number.
6. **Permanent `visited` where path-local is needed.** In Word Search, marking
   a cell visited and never unmarking blocks *other* paths that legally pass
   through it. Conversely, in pure reachability DFS (e.g. flood fill),
   unmarking turns O(V+E) into exponential — know which problem you have.
7. **Backtracking when DP is asked for.** "How many ways…" or "minimum cost…"
   with large `n` shouldn't enumerate; overlapping subproblems + optimal
   substructure → memoise/DP. Backtracking is for when the *listing itself*
   is the output (which is inherently exponential-size).
8. **Recursion-depth / state-passing costs.** Passing big value-type state
   (copying a board each call) multiplies cost; mutate-and-restore a single
   shared structure instead. In Go, declare `var bt func(...)` before
   assigning the closure so it can recurse.
9. **Off-by-one goal tests.** Decide whether the answer is recorded at
   *every* node (Subsets) or only at *leaves* (Permutations, N-Queens) — the
   two templates record in different places.

---

## 6. Complexity cheat sheet

| Shape | Tree size | Typical time |
|-------|-----------|--------------|
| Subsets (include/exclude) | 2^n nodes | O(n · 2^n) |
| Permutations | n! leaves | O(n · n!) |
| Combinations C(n,k) | C(n,k) leaves | O(k · C(n,k)) |
| k-ary labelling (phone digits) | k^n leaves | O(n · k^n) |
| CSP (Sudoku, N-Queens) | exponential worst case, small in practice due to pruning | bounded by constraint propagation |

Space is O(depth) for the recursion stack + path, plus the output.

---

## 7. Problems in this repo

Core backtracking (the concept is the main solution technique):

- [0017 — Letter Combinations of a Phone Number](../0017_letter_combinations_of_a_phone_number/README.md) — one letter per digit; k-ary labelling template
- [0022 — Generate Parentheses](../0022_generate_parentheses/README.md) — prune with open/close counts (Catalan-sized tree)
- [0037 — Sudoku Solver](../0037_sudoku_solver/README.md) — CSP: try digits 1–9 per empty cell, undo on conflict; bool-return short-circuit
- [0039 — Combination Sum](../0039_combination_sum/README.md) — `start` index, element reusable (`bt(i)`), sort + bound prune
- [0040 — Combination Sum II](../0040_combination_sum_ii/README.md) — use-once (`bt(i+1)`) + same-level duplicate skip
- [0046 — Permutations](../0046_permutations/README.md) — `used[]` template; the §4 worked example
- [0047 — Permutations II](../0047_permutations_ii/README.md) — duplicates: sort + `!used[i-1]` skip guard
- [0051 — N-Queens](../0051_n_queens/README.md) — row-by-row placement, O(1) column/diagonal sets
- [0052 — N-Queens II](../0052_n_queens_ii/README.md) — same search, count only; bitmask state
- [0077 — Combinations](../0077_combinations/README.md) — fixed-size k; not-enough-left pruning
- [0078 — Subsets](../0078_subsets/README.md) — record at every node, not just leaves
- [0079 — Word Search](../0079_word_search/README.md) — grid backtracking, mark/unmark cells
- [0090 — Subsets II](../0090_subsets_ii/README.md) — subsets + sort + same-level duplicate skip
- [0093 — Restore IP Addresses](../0093_restore_ip_addresses/README.md) — segment lengths 1–3 per octet, prune invalid values early
- [0113 — Path Sum II](../0113_path_sum_ii/README.md) — tree DFS with path push/pop (choose/unchoose on a tree)

Backtracking appears as one approach or a related tool:

- [0089 — Gray Code](../0089_gray_code/README.md) — can be generated by backtracking over bit flips (formula/mirror approach is optimal)
- [0095 — Unique Binary Search Trees II](../0095_unique_binary_search_trees_ii/README.md) — recursive construction of all trees (enumeration flavour)
- [0126 — Word Ladder II](../0126_word_ladder_ii/README.md) — BFS layering + backtracking DFS through the parent map to emit all shortest paths

> Problems 0131+ are being added concurrently; a later pass will extend this list.
