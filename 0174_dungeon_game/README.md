# 0174 — Dungeon Game

> LeetCode #174 · Difficulty: Hard
> **Categories:** Array, Dynamic Programming, Matrix

---

## Problem Statement

The demons had captured the princess and imprisoned her in **the bottom-right corner** of a `dungeon`. The dungeon consists of `m x n` rooms laid out in a 2D grid. Our valiant knight was initially positioned in **the top-left room** and must fight his way through `dungeon` to rescue the princess.

The knight has an initial health point represented by a positive integer. If at any point his health point drops to `0` or below, he dies immediately.

Some of the rooms are guarded by demons (represented by negative integers), so the knight loses health upon entering these rooms; other rooms are either empty (represented as 0) or contain magic orbs that increase the knight's health (represented by positive integers).

To reach the princess as quickly as possible, the knight decides to move only **rightward** or **downward** in each step.

Return *the knight's minimum initial health so that he can rescue the princess*.

**Note** that any room can contain threats or power-ups, even the first room the knight enters and the bottom-right room where the princess is imprisoned.

**Example 1:**

```
Input: dungeon = [[-2,-3,3],[-5,-10,1],[10,30,-5]]
Output: 7
Explanation: The initial health of the knight must be at least 7 if he follows the optimal path:
RIGHT -> RIGHT -> DOWN -> DOWN.
```

**Example 2:**

```
Input: dungeon = [[0]]
Output: 1
```

**Constraints:**

- `m == dungeon.length`
- `n == dungeon[i].length`
- `1 <= m, n <= 200`
- `-1000 <= dungeon[i][j] <= 1000`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2023          |
| Uber       | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **2D Dynamic Programming** — `need(i,j)` (minimum health on entering room (i,j)) depends only on the two rooms ahead; the twist is that the DP must run **backwards from the goal** → see [`/dsa/dynamic_programming_2d.md`](/dsa/dynamic_programming_2d.md)
- **Matrix Traversal (monotone right/down paths)** — the classic grid-path setting shared with Minimum Path Sum and Unique Paths → see [`/dsa/matrix_traversal.md`](/dsa/matrix_traversal.md)
- **Binary Search (on the answer)** — survivability is monotone in starting health, enabling a search-plus-simulation alternative → see [`/dsa/binary_search.md`](/dsa/binary_search.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (recursive exploration) | O(2^(m+n)) | O(m+n) | To discover the max(1, need−room) recurrence; only for tiny grids |
| 2 | DP Top-Down (memoization) | O(m·n) | O(m·n) + stack | When you have the recursion and want it fast with minimal changes |
| 3 | DP Bottom-Up (2D table) | O(m·n) | O(m·n) | The standard interview answer; easiest to trace and prove |
| 4 | DP Bottom-Up, 1D rolling row (Optimal) | O(m·n) | O(n) | Same speed, minimal memory — the polished final answer |
| 5 | Binary Search on Answer + simulation | O(m·n·log(1000(m+n))) | O(n) | Great fallback if you can't find the backward DP; shows monotonicity insight |

---

## Approach 1 — Brute Force (Recursive Path Exploration)

### Intuition

Reason **backwards from any room**: suppose I already know the minimum health needed to survive *starting from* the room to the right and the room below. The cheaper of those two is what I must still hold **after** this room's effect is applied. So *before* entering, I need `min(right, down) − dungeon[i][j]` — clamped up to 1, because health must stay positive at every instant (a huge healing orb cannot bank health below 1: `max(1, …)` is the heart of the problem). At the princess room, I must leave with at least 1 HP: `max(1, 1 − dungeon[m-1][n-1])`. Recursing this rule from `(0,0)` explores every monotone path.

### Algorithm

1. `need(m-1, n-1) = max(1, 1 − dungeon[m-1][n-1])` (base case: exit the princess room alive).
2. Otherwise `onward = min(need(i+1, j), need(i, j+1))`, considering only in-bounds moves.
3. `need(i, j) = max(1, onward − dungeon[i][j])`.
4. Answer = `need(0, 0)`.

### Complexity

- **Time:** O(2^(m+n)) — every interior cell branches into two calls and identical subproblems are recomputed exponentially often (cell (i,j) is reached via C(i+j, i) call paths).
- **Space:** O(m+n) — recursion depth is the path length.

### Code

```go
func bruteForce(dungeon [][]int) int {
	m, n := len(dungeon), len(dungeon[0])
	var need func(i, j int) int
	need = func(i, j int) int {
		// Base case: the princess room — leave it with at least 1 HP.
		if i == m-1 && j == n-1 {
			return max(1, 1-dungeon[i][j])
		}
		onward := math.MaxInt32 // cheapest requirement among legal next rooms
		if i+1 < m {
			onward = min(onward, need(i+1, j)) // option: move down
		}
		if j+1 < n {
			onward = min(onward, need(i, j+1)) // option: move right
		}
		// Must hold `onward` after this room's effect; clamp at 1 because
		// health may never be 0 or below, even inside a healing room.
		return max(1, onward-dungeon[i][j])
	}
	return need(0, 0)
}
```

### Dry Run

Example 1: `dungeon = [[-2,-3,3],[-5,-10,1],[10,30,-5]]`. Calls in actual evaluation order (down explored before right):

| Step | Call | Room value | Children (down, right) | onward = min(...) | Returns max(1, onward − room) |
|------|------|-----------|------------------------|-------------------|-------------------------------|
| 1 | need(0,0) | −2 | need(1,0), need(0,1) | *pending* | *pending* |
| 2 | need(1,0) | −5 | need(2,0), need(1,1) | *pending* | *pending* |
| 3 | need(2,0) | 10 | —, need(2,1) | *pending* | *pending* |
| 4 | need(2,1) | 30 | —, need(2,2) | *pending* | *pending* |
| 5 | need(2,2) | −5 | base case | — | max(1, 1−(−5)) = **6** |
| 6 | need(2,1) resolves | 30 | (—, 6) | 6 | max(1, 6−30) = **1** |
| 7 | need(2,0) resolves | 10 | (—, 1) | 1 | max(1, 1−10) = **1** |
| 8 | need(1,1) | −10 | need(2,1) **recomputed** → 1 (redoing steps 4–6), need(1,2) | *pending* | *pending* |
| 9 | need(1,2) | 1 | need(2,2) **recomputed** → 6, — | 6 | max(1, 6−1) = **5** |
| 10 | need(1,1) resolves | −10 | (1, 5) | 1 | max(1, 1+10) = **11** |
| 11 | need(1,0) resolves | −5 | (1, 11) | 1 | max(1, 1+5) = **6** |
| 12 | need(0,1) | −3 | need(1,1) **recomputed** → 11 (redoing 8–10), need(0,2) | *pending* | *pending* |
| 13 | need(0,2) | 3 | need(1,2) **recomputed** → 5, — | 5 | max(1, 5−3) = **2** |
| 14 | need(0,1) resolves | −3 | (11, 2) | 2 | max(1, 2+3) = **5** |
| 15 | need(0,0) resolves | −2 | (6, 5) | 5 | max(1, 5+2) = **7** |

Result: `7` ✔ — note steps 8, 9, 12, 13 recomputing whole subtrees: that duplication is what memoization removes.

---

## Approach 2 — DP Top-Down (Memoized Recursion)

### Intuition

`need(i, j)` depends only on the coordinates `(i, j)` — not on how the knight got there — so the exponential call tree contains just m×n *distinct* subproblems. Cache each cell's answer the first time it is computed; every later arrival becomes a table lookup. Since a real requirement is always ≥ 1, the zero value of an `int` table doubles as the "not computed yet" sentinel — no separate visited array needed.

### Algorithm

1. Allocate `memo[m][n]`, zero-initialized (0 = uncomputed).
2. Run the Approach-1 recursion; before computing a cell, return `memo[i][j]` if non-zero; after computing, store it.
3. Answer = `need(0, 0)`.

### Complexity

- **Time:** O(m·n) — each of the m·n subproblems computed exactly once with O(1) work; repeat visits are O(1) lookups.
- **Space:** O(m·n) for the memo table, plus O(m+n) recursion stack.

### Code

```go
func dpTopDown(dungeon [][]int) int {
	m, n := len(dungeon), len(dungeon[0])
	memo := make([][]int, m)
	for i := range memo {
		memo[i] = make([]int, n) // zero value 0 = "uncomputed" sentinel
	}
	var need func(i, j int) int
	need = func(i, j int) int {
		if i == m-1 && j == n-1 {
			return max(1, 1-dungeon[i][j]) // princess room base case
		}
		if memo[i][j] != 0 {
			return memo[i][j] // already solved — reuse
		}
		onward := math.MaxInt32
		if i+1 < m {
			onward = min(onward, need(i+1, j)) // requirement if we go down
		}
		if j+1 < n {
			onward = min(onward, need(i, j+1)) // requirement if we go right
		}
		memo[i][j] = max(1, onward-dungeon[i][j]) // clamp: health ≥ 1 always
		return memo[i][j]
	}
	return need(0, 0)
}
```

### Dry Run

Example 1 — same evaluation order as Approach 1, but now every repeat is a memo hit:

| Step | Call | Result | memo action |
|------|------|--------|-------------|
| 1 | need(2,2) | 6 | base case (not stored; recomputed O(1) on demand) |
| 2 | need(2,1) | 1 | store memo[2][1] = 1 |
| 3 | need(2,0) | 1 | store memo[2][0] = 1 |
| 4 | need(1,1) → asks need(2,1) | hit: 1 | — |
| 5 | need(1,2) | 5 | store memo[1][2] = 5 |
| 6 | need(1,1) | 11 | store memo[1][1] = 11 |
| 7 | need(1,0) | 6 | store memo[1][0] = 6 |
| 8 | need(0,1) → asks need(1,1) | hit: 11 | — |
| 9 | need(0,2) → asks need(1,2) | hit: 5 | store memo[0][2] = 2 |
| 10 | need(0,1) | 5 | store memo[0][1] = 5 |
| 11 | need(0,0) | max(1, min(6, 5) + 2) = **7** | store memo[0][0] = 7 |

Result: `7` ✔ — 9 stored cells + 3 memo hits replace the exponential tree.

---

## Approach 3 — DP Bottom-Up (2D Table from the Princess)

### Intuition

The recursion's dependency arrows all point down-right: `(i,j)` needs `(i+1,j)` and `(i,j+1)`. So fill the table in reverse — rows bottom→top, columns right→left — and both dependencies are always ready. A sentinel border of +∞ with **two 1-valued virtual cells flanking the princess** (`need[m][n-1]` and `need[m-1][n]`) lets a single formula handle the base case, the edges, and the interior identically.

Why not forward DP ("max health achievable at each cell")? Because greedy richness fails: a path that is healthiest *now* may be doomed *later* (Example 1's rich `10, 30` bottom row leads into the −10 trap being avoided, not sought). Correct forward state would need *two* numbers (current health, minimum margin so far); backward DP collapses everything ahead into a single number — that inversion is the whole trick.

### Algorithm

1. Allocate `(m+1)×(n+1)` table `need`, filled with +∞ (`math.MaxInt32`).
2. Set `need[m][n-1] = need[m-1][n] = 1` — virtual rooms "after" the princess meaning: exit alive with exactly 1 HP.
3. For `i = m-1 … 0`, `j = n-1 … 0`: `need[i][j] = max(1, min(need[i+1][j], need[i][j+1]) − dungeon[i][j])`.
4. Return `need[0][0]`.

### Complexity

- **Time:** O(m·n) — each cell filled once with constant work.
- **Space:** O(m·n) — the requirement table (one sentinel row + column extra).

### Code

```go
func dpBottomUp(dungeon [][]int) int {
	m, n := len(dungeon), len(dungeon[0])
	// One extra sentinel row and column so every cell uses the same formula.
	need := make([][]int, m+1)
	for i := range need {
		need[i] = make([]int, n+1)
		for j := range need[i] {
			need[i][j] = math.MaxInt32 // walls: never the min() winner
		}
	}
	// Virtual cells flanking the princess: surviving means exiting with 1 HP.
	need[m][n-1], need[m-1][n] = 1, 1
	for i := m - 1; i >= 0; i-- {
		for j := n - 1; j >= 0; j-- {
			// Cheapest requirement among the two onward rooms...
			req := min(need[i+1][j], need[i][j+1]) - dungeon[i][j]
			if req < 1 {
				req = 1 // ...clamped: health can never sit at 0 or below
			}
			need[i][j] = req
		}
	}
	return need[0][0]
}
```

### Dry Run

Example 1: `dungeon = [[-2,-3,3],[-5,-10,1],[10,30,-5]]` (∞ = MaxInt32 sentinel; virtual cells need[3][2] = need[2][3] = 1).

| Step | Cell (i,j) | room | below need[i+1][j] | right need[i][j+1] | min − room | need[i][j] (clamped ≥ 1) |
|------|-----------|------|--------------------|--------------------|------------|--------------------------|
| 1 | (2,2) | −5 | 1 (virtual) | 1 (virtual) | 1 − (−5) = 6 | **6** |
| 2 | (2,1) | 30 | ∞ | 6 | 6 − 30 = −24 | **1** |
| 3 | (2,0) | 10 | ∞ | 1 | 1 − 10 = −9 | **1** |
| 4 | (1,2) | 1 | 6 | ∞ | 6 − 1 = 5 | **5** |
| 5 | (1,1) | −10 | 1 | 5 | 1 + 10 = 11 | **11** |
| 6 | (1,0) | −5 | 1 | 11 | 1 + 5 = 6 | **6** |
| 7 | (0,2) | 3 | 5 | ∞ | 5 − 3 = 2 | **2** |
| 8 | (0,1) | −3 | 11 | 2 | 2 + 3 = 5 | **5** |
| 9 | (0,0) | −2 | 6 | 5 | 5 + 2 = 7 | **7** |

Final table:

```
need = [ 7  5  2 ]
       [ 6 11  5 ]
       [ 1  1  6 ]
```

Result: `need[0][0] = 7` ✔ (optimal path RIGHT→RIGHT→DOWN→DOWN reads 7→5→2→5→6: every step satisfies its requirement).

---

## Approach 4 — DP Bottom-Up, Space Optimized (1D Rolling Row)

### Intuition

Filling row `i` reads only row `i+1` (below) and already-updated cells of row `i` itself (right). One slice therefore suffices: just before writing index `j`, `dp[j]` still holds *row i+1*'s value ("below") while `dp[j+1]` already holds *row i*'s fresh value ("right"). Seed the slice with +∞ everywhere except `dp[n-1] = 1` (the virtual room below the princess), and the princess cell falls out of the exact same line of code on the first iteration.

### Algorithm

1. `dp` = slice of length `n+1`, all +∞; `dp[n-1] = 1`.
2. For `i = m-1 … 0`, `j = n-1 … 0`: `dp[j] = max(1, min(dp[j], dp[j+1]) − dungeon[i][j])` (overwrite in place).
3. Return `dp[0]`.

### Complexity

- **Time:** O(m·n) — identical arithmetic to Approach 3.
- **Space:** O(n) — one rolling row of n+1 ints instead of the full table.

### Code

```go
func dpSpaceOptimized(dungeon [][]int) int {
	m, n := len(dungeon), len(dungeon[0])
	dp := make([]int, n+1)
	for j := range dp {
		dp[j] = math.MaxInt32 // sentinel: off-grid moves are never chosen
	}
	dp[n-1] = 1 // virtual room below the princess: arrive there with 1 HP
	for i := m - 1; i >= 0; i-- {
		for j := n - 1; j >= 0; j-- {
			// dp[j] = requirement below (old row), dp[j+1] = right (new row).
			req := min(dp[j], dp[j+1]) - dungeon[i][j]
			if req < 1 {
				req = 1 // health floor: even big orbs can't push need below 1
			}
			dp[j] = req // overwrite in place: row i+1's slot becomes row i's
		}
	}
	return dp[0]
}
```

### Dry Run

Example 1 — `dp` shown after each cell write (indices 0..3; ∞ = MaxInt32):

| Step | (i,j) | room | min(dp[j], dp[j+1]) | dp after write |
|------|-------|------|---------------------|----------------|
| 0 | init | — | — | [∞, ∞, 1, ∞] |
| 1 | (2,2) | −5 | min(1, ∞) = 1 → 1+5 = 6 | [∞, ∞, **6**, ∞] |
| 2 | (2,1) | 30 | min(∞, 6) = 6 → −24 → 1 | [∞, **1**, 6, ∞] |
| 3 | (2,0) | 10 | min(∞, 1) = 1 → −9 → 1 | [**1**, 1, 6, ∞] |
| 4 | (1,2) | 1 | min(6, ∞) = 6 → 5 | [1, 1, **5**, ∞] |
| 5 | (1,1) | −10 | min(1, 5) = 1 → 11 | [1, **11**, 5, ∞] |
| 6 | (1,0) | −5 | min(1, 11) = 1 → 6 | [**6**, 11, 5, ∞] |
| 7 | (0,2) | 3 | min(5, ∞) = 5 → 2 | [6, 11, **2**, ∞] |
| 8 | (0,1) | −3 | min(11, 2) = 2 → 5 | [6, **5**, 2, ∞] |
| 9 | (0,0) | −2 | min(6, 5) = 5 → 7 | [**7**, 5, 2, ∞] |

Result: `dp[0] = 7` ✔ — the value sequence matches Approach 3 cell for cell.

---

## Approach 5 — Binary Search on Initial Health + Greedy Simulation

### Intuition

Survivability is **monotone** in starting health: if H hit points suffice, H+1 trivially does (health along any fixed path just shifts up by one, never crossing the death line sooner). Monotone yes/no ⇒ binary search the smallest yes. Checking a *fixed* start H is easy with a **forward** DP — the one that failed for the original problem now works, because with H pinned, health is purely additive, so arriving at each cell with maximum health is unambiguously best; any cell where even the maximum is ≤ 0 is a dead end.

### Algorithm

1. `feasible(H)`: sweep the grid top-left→bottom-right keeping `best[j]` = max health standing on `(i, j)` alive; transition `best = max(from above, from left) + dungeon[i][j]`, marking ≤ 0 as dead; return whether the princess cell is alive.
2. Binary search on `[1, 1000·(m+n)+1]` (any single path visits m+n−1 rooms, each draining ≤ 1000 HP, so the upper bound always survives).
3. Standard lower-bound loop: `feasible(mid)` → `hi = mid`, else `lo = mid+1`; answer `lo`.

### Complexity

- **Time:** O(m·n · log(1000·(m+n))) — one full-grid simulation per probe, ~log₂(6001) ≈ 13 probes for a 3×3, ~19 for 200×200. A log factor slower than the DP.
- **Space:** O(n) — rolling row inside the feasibility check.

### Code

```go
func binarySearchOnAnswer(dungeon [][]int) int {
	m, n := len(dungeon), len(dungeon[0])
	const dead = math.MinInt32 // marker: cannot reach this cell alive

	// feasible reports whether the knight can reach the princess alive
	// when starting with exactly `start` health points.
	feasible := func(start int) bool {
		best := make([]int, n+1) // best[j+1] = max health at (i, j); best[0] = left wall
		for j := range best {
			best[j] = dead
		}
		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				var h int
				if i == 0 && j == 0 {
					h = start + dungeon[0][0] // entering the first room applies its effect
				} else {
					// Best surviving predecessor: above (old best[j+1]) or left (new best[j]).
					prev := max(best[j+1], best[j])
					if prev == dead {
						best[j+1] = dead // no way to arrive alive
						continue
					}
					h = prev + dungeon[i][j] // health is purely additive
				}
				if h <= 0 {
					h = dead // died on entry — 0 or below is instant death
				}
				best[j+1] = h
			}
		}
		return best[n] != dead // alive at the princess cell?
	}

	// Any single path visits m+n-1 rooms, each draining at most 1000 HP,
	// so 1000·(m+n)+1 is always survivable → a valid search upper bound.
	lo, hi := 1, 1000*(m+n)+1
	for lo < hi {
		mid := lo + (hi-lo)/2 // avoid (lo+hi) overflow
		if feasible(mid) {
			hi = mid // mid works — try smaller
		} else {
			lo = mid + 1 // mid dies — need more health
		}
	}
	return lo // smallest health that survives
}
```

### Dry Run

Example 1: search range `[1, 6001]`. Feasibility simulation for the decisive probe `H = 7` (max health after entering each cell; ✝ = dead):

| Cell | from above / left | + room | best health |
|------|-------------------|--------|-------------|
| (0,0) | start 7 | −2 | 5 |
| (0,1) | left 5 | −3 | 2 |
| (0,2) | left 2 | +3 | 5 |
| (1,0) | above 5 | −5 | 0 → ✝ |
| (1,1) | max(above 2, left ✝) = 2 | −10 | −8 → ✝ |
| (1,2) | max(above 5, left ✝) = 5 | +1 | 6 |
| (2,0) | above ✝ | — | ✝ |
| (2,1) | max(✝, ✝) | — | ✝ |
| (2,2) | max(above 6, left ✝) = 6 | −5 | **1 → alive** ✔ |

(For contrast, `H = 6` reaches (2,2) with 6−5−... → the same path arrives at 0 → dead, so 6 is infeasible.)

Binary search probes:

| Step | lo | hi | mid | feasible(mid)? | next range |
|------|----|----|-----|----------------|------------|
| 1 | 1 | 6001 | 3001 | yes | hi = 3001 |
| 2 | 1 | 3001 | 1501 | yes | hi = 1501 |
| 3 | 1 | 1501 | 751 | yes | hi = 751 |
| 4 | 1 | 751 | 376 | yes | hi = 376 |
| 5 | 1 | 376 | 188 | yes | hi = 188 |
| 6 | 1 | 188 | 94 | yes | hi = 94 |
| 7 | 1 | 94 | 47 | yes | hi = 47 |
| 8 | 1 | 47 | 24 | yes | hi = 24 |
| 9 | 1 | 24 | 12 | yes | hi = 12 |
| 10 | 1 | 12 | 6 | **no** | lo = 7 |
| 11 | 7 | 12 | 9 | yes | hi = 9 |
| 12 | 7 | 9 | 8 | yes | hi = 8 |
| 13 | 7 | 8 | 7 | yes | hi = 7 |
| 14 | lo == hi == 7 → return | | | | |

Result: `7` ✔

---

## Key Takeaways

- **When "the future decides the present", run the DP backwards.** Forward "max health so far" is the natural—and wrong—first instinct: it would need two state numbers (health now, worst margin ahead). Anchoring at the goal collapses the future into one number, `need(i,j)`.
- **The clamp `max(1, …)` is the problem.** It encodes "health may never touch 0, and surplus health cannot be banked below the floor of 1". Forgetting it on healing rooms (positive cells) is the classic wrong-answer.
- Sentinel border trick: pad with +∞ and place two `1` cells beside the goal — base case, edge cells, and interior all share one formula.
- Rolling-row 2D→1D compression works whenever dependencies are "next row + already-written same row" — read `dp[j]` (old) vs `dp[j+1]` (new) carefully.
- **Binary search on the answer** is a general escape hatch: monotone feasibility + cheap simulation ≈ optimal answer at a log-factor cost. Recognizing the monotonicity ("more HP never hurts") is itself interview gold.

---

## Related Problems

- LeetCode #64 — Minimum Path Sum (same grid, no death floor — forward DP works there; contrast the two)
- LeetCode #62 — Unique Paths (the counting version of right/down grid DP)
- LeetCode #120 — Triangle (bottom-up DP anchored at the far end)
- LeetCode #741 — Cherry Pickup (harder grid DP where greedy/forward also fails)
- LeetCode #1631 — Path With Minimum Effort (binary search on answer + grid feasibility, same Approach-5 pattern)
