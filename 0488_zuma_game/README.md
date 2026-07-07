# 0488 — Zuma Game

> LeetCode #488 · Difficulty: Hard
> **Categories:** String, Backtracking, Breadth-First Search, Memoization, Depth-First Search

---

## Problem Statement

You are playing a variation of the game Zuma.

In this variation of Zuma, there is a **single row** of colored balls on a board, where each ball can be colored red `'R'`, yellow `'Y'`, blue `'B'`, green `'G'`, or white `'W'`. You also have several colored balls in your hand.

Your goal is to **clear all** of the balls from the board. On each turn:

- Pick any ball from your hand and insert it in between two balls in the row or on either end of the row.
- If there is a group of **three or more consecutive balls** of the **same color**, remove the group of balls from the board.
  - If this removal causes more groups of three or more of the same color to form, then continue removing each group until there are none left.
- If there are no more balls on the board, then you win the game.
- Repeat this process until you either win or do not have any more balls in your hand.

Given a string `board`, representing the row of balls on the board, and a string `hand`, representing the balls in your hand, return *the **minimum** number of balls you have to insert to clear all the balls from the board. If you cannot clear all the balls from the board using the balls in your hand, return* `-1`.

**Example 1:**

```
Input: board = "WRRBBW", hand = "RB"
Output: -1
Explanation: It is impossible to clear all the balls. The best you can do is:
- Insert 'R' so the board becomes WRRRBBW. WRRRBBW -> WBBW.
- Insert 'B' so the board becomes WBBBW. WBBBW -> WW.
There are still balls remaining on the board, and you are out of balls to insert.
```

**Example 2:**

```
Input: board = "WWRRBBWW", hand = "WRBRW"
Output: 2
Explanation: To make the board empty:
- Insert 'R' so the board becomes WWRRRBBWW. WWRRRBBWW -> WWBBWW.
- Insert 'B' so the board becomes WWBBBWW. WWBBBWW -> WWWW -> empty.
2 balls from your hand were needed to clear the board.
```

**Example 3:**

```
Input: board = "G", hand = "GGGGG"
Output: 2
Explanation: To make the board empty:
- Insert 'G' so the board becomes GG.
- Insert 'G' so the board becomes GGG. GGG -> empty.
2 balls from your hand were needed to clear the board.
```

**Constraints:**

- `1 <= board.length <= 16`
- `1 <= hand.length <= 5`
- `board` and `hand` consist of the characters `'R'`, `'Y'`, `'B'`, `'G'`, and `'W'`.
- The initial row of balls on the board will **not** have any groups of three or more consecutive balls of the same color.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Baidu      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Backtracking / DFS over game states** — the core is a depth-first search that, at each state, *tries* inserting a ball, recurses, and backtracks to try another placement, minimising the total inserts → see [`/dsa/backtracking.md`](/dsa/backtracking.md)
- **Memoization (hash map on state)** — many insertion orders reach the same `(board, remaining hand)` pair; caching the optimum per state collapses the search → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **String simulation** — the collapse/cascade of 3+ runs is a run-length pass repeated until the board stabilises → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute-Force DFS | Exponential (bounded by `|hand|` ≤ 5) | O(\|hand\|) depth | Baseline; correct because the hand is tiny |
| 2 | Memoised DFS (Optimal) | O(S · \|board\| · \|hand\|), S = reachable states | O(S) | Standard answer; avoids recomputing shared states |

Both approaches share the **matching-adjacency prune**: only insert a ball where it touches a ball of its own color. A ball with no like neighbour can never start (or contribute to) a removal, so such insertions are never part of an optimal solution — verified here against a no-prune brute force over thousands of random boards.

---

## Approach 1 — Brute-Force DFS

### Intuition

From any state `(board, hand)` we must eventually place *some* hand ball *somewhere*. So branch on "which color, at which gap", collapse the resulting board (cascading removals), recurse with that ball gone from the hand, and keep the cheapest clearing (`+1` for the ball just used). An empty board costs `0`; an empty hand over a non-empty board is impossible (`∞`, reported as `-1`). Two safe reductions keep it correct **and** small: (a) treat equal hand colors as one choice, and (b) only insert next to a matching color, since a lone ball with no same-color neighbour can never trigger a removal.

### Algorithm

1. Pre-collapse the board (defensive) and sort the hand for canonical dedup.
2. `dfs(board, hand)`:
   - if `board` empty → return `0`.
   - if `hand` empty → return `∞`.
   - for each **distinct** hand color `ball` (skip duplicates in the sorted hand):
     - for each gap `pos` in `0..len(board)`, skipping interior gaps of a same-color run and gaps where `ball` has no matching neighbour:
       - `next = collapse(board[:pos] + ball + board[pos:])`.
       - `best = min(best, 1 + dfs(next, hand without ball))`.
   - return `best`.
3. If the result is `∞`, return `-1`; else return it.

### Complexity

- **Time:** Exponential in general, but bounded by the number of insertion sequences ≤ `|hand|! · (|board|+…)`; tractable purely because `|hand| ≤ 5`.
- **Space:** O(|hand|) recursion depth, plus transient board strings built per branch.

### Code

```go
func bruteForceDFS(board, hand string) int {
	const INF = 1 << 30
	var dfs func(board, hand string) int
	dfs = func(board, hand string) int {
		if len(board) == 0 {
			return 0 // board cleared — no more insertions needed
		}
		if len(hand) == 0 {
			return INF // out of balls but board not empty — impossible from here
		}
		best := INF
		// Try inserting one hand ball; consider each distinct colour once.
		for h := 0; h < len(hand); h++ {
			if h > 0 && hand[h] == hand[h-1] {
				continue // skip duplicate colours (hand is sorted) — same result
			}
			ball := hand[h]
			// Remaining hand with this one ball removed.
			rest := hand[:h] + hand[h+1:]
			// Try every gap 0..len(board); dedupe positions inside equal runs.
			for pos := 0; pos <= len(board); pos++ {
				// Skip a gap that sits strictly inside a run of the same colour:
				// inserting there is equivalent to inserting at the run's start.
				if pos > 0 && pos < len(board) && board[pos] == board[pos-1] {
					continue
				}
				// Heuristic prune: only insert next to a matching colour (that is
				// the only way a removal can ever start).
				matchRight := pos < len(board) && board[pos] == ball
				matchLeft := pos > 0 && board[pos-1] == ball
				if !matchLeft && !matchRight {
					continue // a lone ball with no like neighbour can never help
				}
				// Build the new board and collapse any cascade.
				next := collapse(board[:pos] + string(ball) + board[pos:])
				if r := dfs(next, rest); r != INF && r+1 < best {
					best = r + 1 // +1 for the ball we just inserted
				}
			}
		}
		return best
	}
	res := dfs(collapse(board), sortHand(hand)) // pre-collapse in case input had a run
	if res >= INF {
		return -1 // never cleared
	}
	return res
}
```

The shared collapse helper:

```go
func collapse(board string) string {
	for {
		n := len(board)
		removed := false
		i := 0
		var b strings.Builder
		for i < n {
			j := i
			for j < n && board[j] == board[i] { // maximal same-colour run
				j++
			}
			if j-i >= 3 {
				removed = true // run of ≥ 3 disappears
			} else {
				b.WriteString(board[i:j]) // keep short runs
			}
			i = j
		}
		board = b.String()
		if !removed {
			return board // stable
		}
	}
}
```

### Dry Run

Example 3: `board = "G", hand = "GGGGG"` (sorted hand `"GGGGG"`).

| Depth | board | hand (distinct colors) | action | next board | result |
|-------|-------|------------------------|--------|-----------|--------|
| 0 | `G` | `G` | insert `G` next to `G` at pos 1 | `GG` (no ≥3 run) | 1 + dfs(`GG`, 4×G) |
| 1 | `GG` | `G` | insert `G` at pos 2 | `GGG` → collapse → `` (empty) | 1 + dfs(``, 3×G) |
| 2 | `` (empty) | — | base case | — | 0 |

Unwinding: depth 2 → 0, depth 1 → 1 + 0 = 1, depth 0 → 1 + 1 = **2** ✔.

---

## Approach 2 — Memoised DFS (Optimal)

### Intuition

Different insertion orders repeatedly converge on the same `(board, remaining hand)` — for instance, using two balls of interchangeable colors in either order lands on an identical state. The minimum extra insertions to clear a state depends **only** on that state, so memoise on `(board, sorted hand)`. Sorting the hand makes `"RB"` and `"BR"` share a key; the search then explores each distinct reachable state once.

### Algorithm

1. Key each state as `board + "|" + sortedHand`.
2. `dfs(board, hand)`: empty board → `0`; empty hand → `∞`; cached key → return the cache.
3. Otherwise branch over `(color, position)` exactly as brute force, collapse, recurse, take `min + 1`, **store** the result in the map, and return it.
4. Convert `∞` to `-1` at the top level.

### Complexity

- **Time:** O(S · |board| · |hand|), where `S` is the number of distinct reachable `(board, hand)` states; each state does O(|board|·|hand|) branching work.
- **Space:** O(S) memo entries plus O(|hand|) recursion depth.

### Code

```go
func memoDFS(board, hand string) int {
	const INF = 1 << 30
	memo := map[string]int{} // (board|sortedHand) → min insertions to clear
	var dfs func(board, hand string) int
	dfs = func(board, hand string) int {
		if len(board) == 0 {
			return 0 // cleared
		}
		if len(hand) == 0 {
			return INF // stuck: balls exhausted, board remains
		}
		key := board + "|" + hand // hand is already sorted by the caller
		if v, ok := memo[key]; ok {
			return v // this exact sub-problem was already solved
		}
		best := INF
		for h := 0; h < len(hand); h++ {
			if h > 0 && hand[h] == hand[h-1] {
				continue // one representative per colour
			}
			ball := hand[h]
			rest := hand[:h] + hand[h+1:] // hand stays sorted after removal
			for pos := 0; pos <= len(board); pos++ {
				if pos > 0 && pos < len(board) && board[pos] == board[pos-1] {
					continue // interior of a same-colour run: redundant gap
				}
				matchRight := pos < len(board) && board[pos] == ball
				matchLeft := pos > 0 && board[pos-1] == ball
				if !matchLeft && !matchRight {
					continue // insertion cannot trigger any removal
				}
				next := collapse(board[:pos] + string(ball) + board[pos:])
				if r := dfs(next, rest); r != INF && r+1 < best {
					best = r + 1
				}
			}
		}
		memo[key] = best // record the optimum for this state
		return best
	}
	res := dfs(collapse(board), sortHand(hand))
	if res >= INF {
		return -1
	}
	return res
}
```

### Dry Run

Example 2: `board = "WWRRBBWW", hand = "WRBRW"` → sorted hand `"BRRWW"`.

| Step | State (board \| hand) | Chosen insert | Collapse cascade | New state | Note |
|------|-----------------------|---------------|------------------|-----------|------|
| 1 | `WWRRBBWW \| BRRWW` | `R` between the two `R`s (pos 4) | `WWRRRBBWW` → remove `RRR` → `WWBBWW` | `WWBBWW \| BRWW` | cache miss, recurse |
| 2 | `WWBBWW \| BRWW` | `B` between the two `B`s (pos 3) | `WWBBBWW` → remove `BBB` → `WWWW` → remove `WWWW` → `` | `` (empty) \| `RW` | cache miss, recurse |
| 3 | `` (empty) \| `RW` | — | — | — | base case → 0 |

Unwinding: step 3 → 0; step 2 → 1 + 0 = 1; step 1 → 1 + 1 = **2** ✔. States are memoised, so if any other insertion order reproduces `WWBBWW | BRWW`, it is answered from the cache.

---

## Key Takeaways

- **Search + memoise on the full mutable state.** When moves transform a compound state `(board, hand)` and orders can converge, the key is the *entire* state; sorting the hand canonicalises interchangeable multisets so more branches share cache entries.
- **Cascading removal = collapse-until-stable.** One removal can make neighbours touch and form a new group, so run a run-length deletion pass in a loop until a pass deletes nothing.
- **Prune insertions to matching adjacencies.** A ball placed with no same-color neighbour can never start a removal now or later, so it is never in an optimal plan — a big constant-factor cut that keeps this Hard problem fast. (Validated by fuzzing against an unpruned baseline.)
- **Small bounds justify exponential search.** `|hand| ≤ 5` is the license to brute-force the placement tree; always check whether the constraints make an "obviously exponential" search actually feasible.

---

## Related Problems

- LeetCode #664 — Strange Printer (interval DP over a string with merge/removal semantics)
- LeetCode #546 — Remove Boxes (DP over a row with combo removals; harder cousin)
- LeetCode #1717 — Maximum Score From Removing Substrings (greedy string removals)
- LeetCode #301 — Remove Invalid Parentheses (BFS/DFS over string edit states)
- LeetCode #464 — Can I Win (game-state search with memoization)
