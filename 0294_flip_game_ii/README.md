# 0294 — Flip Game II

> LeetCode #294 · Difficulty: Medium
> **Categories:** Backtracking, Memoization, Game Theory, Dynamic Programming

---

## Problem Statement

You are playing a Flip Game with your friend.

You are given a string `currentState` that contains only `'+'` and `'-'`. You and your friend take turns to flip **two consecutive** `"++"` into `"--"`. The game ends when a person can no longer make a move, and therefore the other person will be the winner.

Return `true` if the **starting player** can guarantee a win, and `false` otherwise.

**Example 1:**
```
Input: currentState = "++++"
Output: true
Explanation: The starting player can guarantee a win by flipping the middle "++" to become "+--+".
```

**Example 2:**
```
Input: currentState = "+"
Output: false
```

**Constraints:**
- `1 <= currentState.length <= 60`
- `currentState[i]` is either `'+'` or `'-'`.

**Follow-up:** Derive your algorithm's runtime complexity.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★☆☆☆ Low        | 2022          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★☆☆☆☆ Rare       | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used
- **Backtracking** — try every move, recurse on the opponent's position, undo → see [`/dsa/backtracking.md`](/dsa/backtracking.md)
- **Memoization / DP** — cache each board's win/lose verdict → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)
- **Game Theory (Sprague–Grundy)** — decompose into independent runs, XOR their Grundy numbers → see [`/dsa/game_theory.md`](/dsa/game_theory.md)

---

## Approaches Overview
| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Backtracking | O(n!!) (exp.) | O(n·depth) | Baseline minimax |
| 2 | Backtracking + Memo | O(2ⁿ) reachable boards | O(2ⁿ) | Kill duplicate boards |
| 3 | Sprague–Grundy (Optimal) | O(m² + n) | O(m) | Polynomial; the intended follow-up |

(*n* = string length, *m* = longest run of `'+'`.)

---

## Approach 1 — Backtracking

### Intuition
This is a two-player, perfect-information game. The player to move wins iff there is at least one move after which the opponent (now to move on the new string) LOSES. That is a direct minimax recursion: `canWin(s)` is true iff some flippable `"++"` leads to a state where `canWin(next)` is false.

### Algorithm
1. Scan for `"++"`; for each, flip it to `"--"`.
2. Recurse on the new board; if the opponent cannot win from it, return `true`.
3. Undo the flip before trying the next candidate.
4. If no move wins, return `false`.

### Complexity
- **Time:** roughly O(n!!) (the game-tree size) — exponential; identical boards are re-explored.
- **Space:** O(n·depth) — a new string per move along a path of depth up to `n/2`.

### Code
```go
func backtracking(s string) bool {
	b := []byte(s) // mutable working copy of the board
	for i := 0; i+1 < len(b); i++ {
		if b[i] == '+' && b[i+1] == '+' { // a legal move exists here
			b[i], b[i+1] = '-', '-'                 // make the move
			opponentWins := backtracking(string(b)) // opponent plays on the new board
			b[i], b[i+1] = '+', '+'                 // undo the move
			if !opponentWins {                      // opponent can't win → we win here
				return true
			}
		}
	}
	return false // no move leaves the opponent losing → current player loses
}
```

### Dry Run
`s = "++++"`

| Depth | Board | Move tried | Opponent result | Verdict |
|-------|-------|-----------|-----------------|---------|
| 0 | `++++` | flip i=0 → `--++` | opponent on `--++` can flip → wins | not yet |
| 0 | `++++` | flip i=1 → `+--+` | opponent on `+--+` has **no move** → loses | **we win** |

The first player flips the middle pair to `+--+`, leaving the opponent stuck → `true`. ✓

---

## Approach 2 — Backtracking + Memoization

### Intuition
Many move sequences reach the same board (flips that are far apart commute). Memoizing on the exact string collapses that redundancy — a large improvement while keeping the same minimax logic.

### Algorithm
1. If `s` is in the memo, return the cached verdict.
2. Otherwise run the minimax scan (we win iff some move makes the opponent lose).
3. Store and return the result.

### Complexity
- **Time:** O(2ⁿ) worst case, bounded by the number of distinct reachable boards.
- **Space:** O(2ⁿ) memo in the worst case + O(n) recursion stack.

### Code
```go
func memoBacktracking(s string) bool {
	memo := map[string]bool{}
	var solve func(cur string) bool
	solve = func(cur string) bool {
		if v, ok := memo[cur]; ok { // already decided this exact board
			return v
		}
		b := []byte(cur)
		result := false
		for i := 0; i+1 < len(b); i++ {
			if b[i] == '+' && b[i+1] == '+' {
				b[i], b[i+1] = '-', '-'
				win := !solve(string(b)) // we win if opponent loses on the new board
				b[i], b[i+1] = '+', '+'
				if win {
					result = true
					break // one winning move suffices
				}
			}
		}
		memo[cur] = result // remember the verdict for this board
		return result
	}
	return solve(s)
}
```

### Dry Run
`s = "++++"`

| Call | Board | Cached? | Explores | Stored |
|------|-------|---------|----------|--------|
| 1 | `++++` | no | flip i=0 → `--++` (recurse), flip i=1 → `+--+` (recurse) | true |
| 2 | `--++` | no | flip i=2 → `----` → no move → false; so `--++` is win | true |
| 3 | `+--+` | no | no `"++"` → no move → **false** | false |

At call 1, trying i=1 gives opponent board `+--+` (memo = false) → we win → `memo["++++"]=true`. ✓

---

## Approach 3 — Sprague–Grundy (Optimal)

### Intuition
A move flips two ADJACENT pluses inside some maximal run of consecutive `'+'`. Flipping inside a run of length `k` removes 2 pluses and splits it into two INDEPENDENT shorter runs. By the **Sprague–Grundy theorem**, the Grundy number of the whole game is the XOR of the Grundy numbers of its independent runs, and `G(k) = mex{ G(left) XOR G(right) }` over all legal splits (where `left + right = k - 2`). The first player wins iff the total nim-sum is nonzero.

### Algorithm
1. Split `currentState` into the lengths of maximal `'+'` runs.
2. Precompute `g[0..maxRun]`: for each `k`, take the mex over `g[left] XOR g[k-2-left]` for `left = 0..k-2`.
3. XOR `g[len]` for every run to get the nim-sum.
4. Return `nimSum != 0`.

### Complexity
- **Time:** O(m² + n) — building the Grundy table costs O(m²) for the longest run `m`; scanning the string is O(n).
- **Space:** O(m) — the Grundy table.

### Code
```go
func spragueGrundy(s string) bool {
	n := len(s)
	runs := []int{}
	count := 0
	for i := 0; i < n; i++ {
		if s[i] == '+' {
			count++
		} else if count > 0 {
			runs = append(runs, count)
			count = 0
		}
	}
	if count > 0 {
		runs = append(runs, count)
	}
	maxRun := 0
	for _, r := range runs {
		if r > maxRun {
			maxRun = r
		}
	}
	g := make([]int, maxRun+1)
	for k := 2; k <= maxRun; k++ {
		seen := map[int]bool{}
		for left := 0; left <= k-2; left++ {
			right := k - 2 - left
			seen[g[left]^g[right]] = true
		}
		mex := 0
		for seen[mex] {
			mex++
		}
		g[k] = mex
	}
	nim := 0
	for _, r := range runs {
		nim ^= g[r]
	}
	return nim != 0
}
```

### Dry Run
`s = "++++"`

| Step | Computation | Value |
|------|-------------|-------|
| runs | one maximal `+` run of length 4 | `[4]` |
| g[0] | base | 0 |
| g[1] | no split possible (needs ≥2) | 0 |
| g[2] | split left=0,right=0 → `g0^g0=0`; mex{0}=1 | 1 |
| g[3] | left∈{0,1}: `g0^g1=0`, `g1^g0=0`; mex{0}=1 | 1 |
| g[4] | left∈{0,1,2}: `g0^g2=1`, `g1^g1=0`, `g2^g0=1`; seen={0,1}; mex=2 | 2 |
| nim-sum | `g[4] = 2` | 2 |
| result | `2 != 0` | **true** |

Nonzero nim-sum → first player wins. ✓

---

## Key Takeaways
- Win/lose games follow the **"win iff some move leaves the opponent losing"** minimax recurrence.
- **Memoize on the board state** to kill the exponential re-exploration of order-independent moves.
- **Sprague–Grundy** turns an impartial game into arithmetic: decompose into independent sub-games, compute each Grundy number via `mex` of reachable XORs, and XOR them — nonzero ⇒ first player wins. This answers the follow-up with a polynomial algorithm.

---

## Related Problems
- LeetCode #293 — Flip Game (enumerate one-move states)
- LeetCode #292 — Nim Game (game theory, closed form)
- LeetCode #464 — Can I Win (game DP with bitmask memo)
- LeetCode #877 — Stone Game (game theory / DP)
