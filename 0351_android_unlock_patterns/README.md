# 0351 — Android Unlock Patterns

> LeetCode #351 · Difficulty: Medium
> **Categories:** Backtracking, DFS, Dynamic Programming, Bit Manipulation

---

## Problem Statement

Android devices have a special lock screen with a `3 x 3` grid of dots. Users can
set an "unlock pattern" by connecting the dots in a specific sequence, forming a
series of joined line segments where each segment's endpoints are two distinct
dots in the sequence. A sequence of `k` dots is a **valid** unlock pattern if
both of the following are true:

- All the dots in the sequence are **distinct**.
- If the line segment connecting two consecutive dots in the sequence passes
  through the center of any other dot, that other dot **must have appeared
  previously** in the sequence. No jumps through unselected dots are allowed.

Here are some example valid and invalid unlock patterns:

- `4 - 1 - 3 - 6` is **invalid** because the line from `1` to `3` passes through
  `2`, which was not selected before.
- `4 - 1 - 9 - 2` is **invalid** because the line from `1` to `9` passes through
  `5`, which was not selected before.
- `2 - 4 - 1 - 3 - 6` is **valid**: connecting `1` to `3` is fine because `2`
  was selected previously.
- `6 - 5 - 4 - 1 - 9 - 2` is **valid**: connecting `1` to `9` is fine because
  `5` was selected previously.

Given two integers `m` and `n`, return *the number of unlock patterns of the
Android lock screen that consist of at least `m` keys and at most `n` keys.*

**Example 1:**

```
Input: m = 1, n = 1
Output: 9
```

**Example 2:**

```
Input: m = 1, n = 2
Output: 65
```

**Constraints:**

- `1 <= m <= n <= 9`

> Note: The 3×3 grid is numbered
> ```
> 1 2 3
> 4 5 6
> 7 8 9
> ```
> The total count over all lengths (m=1, n=9) is **389497** (per-length counts:
> 9, 56, 320, 1624, 7152, 26016, 72912, 140704, 140704).

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2023          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Uber       | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Backtracking / DFS** — grow the pattern one dot at a time, undoing the
  choice on the way back up → see [`/dsa/backtracking.md`](/dsa/backtracking.md)
- **State + visited pruning** — a `used[]` array both enforces distinctness and
  lets us test whether a "skipped" middle dot has already been visited → see
  [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Symmetry exploitation** — the grid's 8-fold symmetry collapses 9 start
  positions into 3 equivalence classes (corner / edge / center) → see
  [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Plain Backtracking (DFS) | O(9!) | O(9) | Baseline; clearest statement of the rules |
| 2 | Backtracking + Symmetry (Optimal) | O(9!) | O(9) | Same tree, ~1/3 the constant work |

Both are exponential in the tiny fixed grid; "optimal" here means fewer constant
passes via symmetry, not a lower asymptotic class.

---

## Approach 1 — Plain Backtracking (DFS)

### Intuition
A pattern is a path visiting distinct dots. Grow it dot by dot. From the current
dot, any unused dot is a legal next step **unless** the straight segment to it
jumps over a middle dot that has not yet been visited. Precompute, for every
ordered pair `(a, b)`, the single dot (if any) that lies exactly between them
when they are collinear — the `mid` table. A move `a -> b` is legal iff
`mid[a][b] == 0` (nothing between) or `mid[a][b]` is already used.

### Algorithm
1. Build the `mid[10][10]` table: horizontals (1-3→2, 4-6→5, 7-9→8), verticals
   (1-7→4, 2-8→5, 3-9→6), diagonals (1-9→5, 3-7→5).
2. For each starting dot `1..9`: mark it used, DFS at depth 1.
3. In the DFS, if `m <= depth <= n`, count the current path.
4. If `depth == n`, stop (deeper paths exceed `n`).
5. Try every unused `to`; if `mid[cur][to] == 0` or that middle dot is used,
   mark `to`, recurse at `depth+1`, then unmark (backtrack).

### Complexity
- **Time:** O(9!) — the search tree is bounded by the number of ordered
  distinct-dot sequences of length ≤ 9. Fixed and small.
- **Space:** O(9) — recursion depth plus the `used[]` array of size 10.

### Code
```go
func bruteForceBacktracking(m, n int) int {
	mid := buildMid()
	used := make([]bool, 10) // used[k] == true once key k is in the path
	count := 0

	var dfs func(cur, depth int)
	dfs = func(cur, depth int) {
		if depth >= m && depth <= n {
			count++ // current path is itself a valid pattern
		}
		if depth == n {
			return // cannot grow further; deeper paths exceed n
		}
		for to := 1; to <= 9; to++ {
			if used[to] {
				continue // keys must be distinct
			}
			jumped := mid[cur][to]                // key skipped over, or 0
			if jumped == 0 || used[jumped] {      // legal move?
				used[to] = true                   // choose
				dfs(to, depth+1)                  // explore
				used[to] = false                  // un-choose (backtrack)
			}
		}
	}

	for start := 1; start <= 9; start++ {
		used[start] = true
		dfs(start, 1)
		used[start] = false
	}
	return count
}
```

### Dry Run
Input `m = 1, n = 2`. We show the start dot `1` (all 9 starts behave alike by
symmetry for length-1, and each contributes its length-2 extensions).

| Step | cur | depth | Action | count |
|------|-----|-------|--------|-------|
| 1 | 1 | 1 | `1<=1<=2` → count the single-dot pattern `[1]` | 1 |
| 2 | 1 | 1 | try to=2: mid[1][2]=0 legal → recurse | 1 |
| 3 | 2 | 2 | `1<=2<=2` → count `[1,2]` | 2 |
| 4 | 2 | 2 | depth==n → return, backtrack (unmark 2) | 2 |
| 5 | 1 | 1 | try to=3: mid[1][3]=2, 2 unused → **illegal**, skip | 2 |
| 6 | 1 | 1 | try to=4,5,6,8: mid=0 → each gives a length-2 pattern | 6 |
| 7 | 1 | 1 | try to=7: mid[1][7]=4 unused → illegal, skip | 6 |
| 8 | 1 | 1 | try to=9: mid[1][9]=5 unused → illegal, skip | 6 |

Start `1` yields 1 (itself) + 5 length-2 patterns = 6 (neighbors 2,4,5,6,8).
Summing over all 9 starts with the same rule gives the total **65** for `n = 2`
(9 single-dot patterns + 56 two-dot patterns).

---

## Approach 2 — Backtracking + Symmetry (Optimal)

### Intuition
The 3×3 grid is symmetric under rotations and reflections (the dihedral group of
order 8). Under those symmetries the four **corners** (1,3,7,9) are
interchangeable, the four **edges** (2,4,6,8) are interchangeable, and the
**center** (5) is alone. So the number of patterns starting at any corner is the
same, likewise for edges. Compute the DFS once per class and combine:
`total = 4·f(corner) + 4·f(edge) + 1·f(center)`.

### Algorithm
1. Same `mid` table and legality rule as Approach 1, but make `dfs` **return**
   the number of valid patterns rooted at the current partial path.
2. Run `dfs` from dot `1` (a corner), dot `2` (an edge), and dot `5` (center).
3. Return `4*count(1) + 4*count(2) + count(5)`.

### Complexity
- **Time:** O(9!) asymptotically identical, but only 3 of the 9 start subtrees
  are explored → roughly one-third of the constant work.
- **Space:** O(9) — recursion depth plus `used[]`.

### Code
```go
func symmetryBacktracking(m, n int) int {
	mid := buildMid()
	used := make([]bool, 10)

	var dfs func(cur, depth int) int
	dfs = func(cur, depth int) int {
		res := 0
		if depth >= m && depth <= n {
			res++ // this path is a valid pattern
		}
		if depth == n {
			return res
		}
		for to := 1; to <= 9; to++ {
			if used[to] {
				continue
			}
			jumped := mid[cur][to]
			if jumped == 0 || used[jumped] {
				used[to] = true
				res += dfs(to, depth+1)
				used[to] = false
			}
		}
		return res
	}

	countFrom := func(start int) int {
		used[start] = true
		c := dfs(start, 1)
		used[start] = false
		return c
	}

	corner := countFrom(1) // representative corner
	edge := countFrom(2)   // representative edge
	center := countFrom(5) // the center
	return 4*corner + 4*edge + center
}
```

### Dry Run
Input `m = 1, n = 1` (only single-dot patterns).

| Step | Call | depth | `m<=depth<=n`? | returns |
|------|------|-------|----------------|---------|
| 1 | `countFrom(1)` → dfs(1,1) | 1 | yes → res=1; depth==n → return | 1 |
| 2 | `countFrom(2)` → dfs(2,1) | 1 | yes → res=1; return | 1 |
| 3 | `countFrom(5)` → dfs(5,1) | 1 | yes → res=1; return | 1 |
| 4 | combine | — | `4*1 + 4*1 + 1` | **9** |

Result `9`, matching Example 1. For `n = 2` the same three DFS calls return the
length-1 plus length-2 counts of a corner (6), edge (8), and center (9), giving
`4*6 + 4*8 + 9 = 24 + 32 + 9 = 65`.

---

## Key Takeaways

- **Precompute the "skipped middle" table.** The only non-trivial rule is the
  jump-over constraint; encoding it as `mid[a][b]` turns each move-legality
  check into O(1).
- The jump rule is equivalent to: for collinear dots, all integer lattice points
  strictly between must already be visited. On a 3×3 grid the only such pairs are
  the 8 listed lines — every "knight" move (e.g. 1→6, 2→7) skips nothing.
- **Exploit symmetry to cut constant factors.** When start states fall into a few
  equivalence classes, compute one representative each and multiply.
- The published total for lengths 1–9 is **389497**; a common web figure of
  "389112" is incorrect (verified by an independent GCD-based lattice check).

---

## Related Problems

- LeetCode #46 — Permutations (ordered distinct selections via backtracking)
- LeetCode #79 — Word Search (grid DFS with a visited set)
- LeetCode #52 — N-Queens II (count valid configurations via backtracking)
- LeetCode #980 — Unique Paths III (count constrained full-grid paths)
