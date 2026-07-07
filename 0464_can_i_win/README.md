# 0464 — Can I Win

> LeetCode #464 · Difficulty: Medium
> **Categories:** Math, Dynamic Programming, Bit Manipulation, Memoization, Game Theory

---

## Problem Statement

In the "100 game" two players take turns adding, to a running total, any integer from `1` to `10`. The player who first causes the running total to reach or exceed `100` wins.

What if we change the game so that players **cannot** re-use integers?

For example, two players might take turns drawing from a common pool of numbers from `1` to `15` without replacement until they reach a total `>= 100`.

Given two integers `maxChoosableInteger` and `desiredTotal`, return `true` if the first player to move can force a win, otherwise, return `false`. Assume both players play optimally.

**Example 1:**

```
Input: maxChoosableInteger = 10, desiredTotal = 11
Output: false
Explanation:
No matter which integer the first player chooses, the first player will lose.
The first player can choose an integer from 1 up to 10.
If the first player chooses 1, the second player can only choose integers from 2 up to 10.
The second player will win by choosing the integer 10 and get a total = 11, which is >= desiredTotal.
Same with other integers chosen by the first player, the second player will always win.
```

**Example 2:**

```
Input: maxChoosableInteger = 10, desiredTotal = 0
Output: true
```

**Example 3:**

```
Input: maxChoosableInteger = 10, desiredTotal = 1
Output: true
```

**Constraints:**

- `1 <= maxChoosableInteger <= 20`
- `0 <= desiredTotal <= 300`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Google     | ★★★☆☆ Medium     | 2023          |
| LinkedIn   | ★★☆☆☆ Low        | 2022          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Game Theory (Minimax)** — a two-player, perfect-information, zero-sum game: a position is a win iff *some* move leaves the opponent in a losing position. This win/lose recursion is the whole solution → see [`/dsa/game_theory.md`](/dsa/game_theory.md)
- **Bit Manipulation (bitmask state)** — the set of still-available numbers (max ≤ 20) is encoded as a 20-bit integer; bit `i` set ⇒ number `i+1` is used. The mask is both the memo key and the branching state → see [`/dsa/bitmask.md`](/dsa/bitmask.md)
- **Memoization / DP over subsets** — the outcome depends only on which numbers remain (order-independent), so each of the `2^max` masks is solved once and cached → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Minimax Backtracking (no memo) | O(n!) | O(n) | Establishes the game recursion; TLE for larger n |
| 2 | Minimax + Bitmask Memoization (Optimal) | O(2ⁿ · n) | O(2ⁿ) | The accepted solution; caches each used-set once |

---

## Approach 1 — Minimax Backtracking (No Memo)

### Intuition

This is a zero-sum game with perfect information, so it is decided by minimax. On your turn you **win** if either (a) some unused number `x` is large enough to reach the target immediately (`x ≥ remaining`), or (b) some `x` leaves the opponent in a position from which *they* cannot win. Represent "which numbers are still available" as a bitmask; the first player controls the root. Two guards handle edge cases: a target of `0` (or less) is already won, and if the sum `1..max` is below the target, nobody can ever reach it.

### Algorithm

1. If `desiredTotal ≤ 0` → `true` (target met before any move).
2. If `max·(max+1)/2 < desiredTotal` → `false` (unreachable even using every number).
3. `canWin(used, remaining)`: for each unused number `x`:
   - if `x ≥ remaining` → return `true` (immediate win), or
   - if `!canWin(used | bit(x), remaining − x)` → return `true` (opponent loses).
4. If no move wins, return `false`.

### Complexity

- **Time:** O(n!) worst case — every ordering of picks is explored because nothing is cached.
- **Space:** O(n) — recursion depth (at most `n` numbers taken).

### Code

```go
func bruteForceMinimax(maxChoosableInteger int, desiredTotal int) bool {
	if desiredTotal <= 0 {
		return true // target already met before anyone moves
	}
	// If even grabbing every number cannot reach the target, first player can't win.
	if maxChoosableInteger*(maxChoosableInteger+1)/2 < desiredTotal {
		return false
	}

	var canWin func(used, remaining int) bool
	canWin = func(used, remaining int) bool {
		for x := 1; x <= maxChoosableInteger; x++ {
			bit := 1 << uint(x-1) // bit representing "number x is taken"
			if used&bit != 0 {
				continue // x already used this game
			}
			// Win now (x meets/exceeds what's left) OR opponent loses after we take x.
			if x >= remaining || !canWin(used|bit, remaining-x) {
				return true
			}
		}
		return false // no move leads to a win → current player loses
	}
	return canWin(0, desiredTotal)
}
```

### Dry Run

Example 1: `max = 10, desiredTotal = 11`. Guards pass (`11 > 0`; sum `1..10 = 55 ≥ 11`). Root `canWin(used=0, remaining=11)` tries each first move `x`:

| First move x | remaining after | Opponent's position | Opponent can win? | Does x win for P1? |
|--------------|-----------------|----------------------|-------------------|--------------------|
| 1 | 10 | can pick 10 (≥10) | yes | no |
| 2 | 9 | can pick 10 (≥9) | yes | no |
| … | … | can always pick some `y ≥ remaining` (10 is still there) | yes | no |
| 10 | 1 | any pick ≥ 1 wins instantly | yes | no |

For every first move, `remaining ≤ 10`, and the number `10` (or another large one) is still available to the opponent, so the opponent immediately reaches the total. No first move gives P1 a win → root returns `false`.

Result: `false` ✔

---

## Approach 2 — Minimax + Bitmask Memoization (Optimal)

### Intuition

The outcome of a position depends **only on which numbers remain**, not on the order they were drawn — because `remaining = desiredTotal − Σ(used numbers)`, and the used-set determines that sum. So the bitmask `used` alone is a complete state key. There are only `2^max` distinct masks (`max ≤ 20`), so memoize each: solve a mask once, store win/lose, reuse everywhere. This collapses the factorial tree into `O(2^max · max)`.

### Algorithm

1. Same two guards as Approach 1.
2. `memo[used] ∈ {0 unknown, 1 win, −1 lose}`.
3. `canWin(used, remaining)`: if `memo[used]` is known, return it; otherwise try each unused `x` (win if `x ≥ remaining` or the opponent loses on the child); store and return the verdict.

### Complexity

- **Time:** O(2^max · max) — each of `2^max` masks does at most `max` work once.
- **Space:** O(2^max) — the memo table, plus O(max) recursion depth.

### Code

```go
func bitmaskMemo(maxChoosableInteger int, desiredTotal int) bool {
	if desiredTotal <= 0 {
		return true
	}
	if maxChoosableInteger*(maxChoosableInteger+1)/2 < desiredTotal {
		return false
	}

	// memo maps a used-mask to the current mover's outcome: 0 unknown, 1 win, -1 lose.
	memo := make([]int8, 1<<uint(maxChoosableInteger))

	var canWin func(used, remaining int) bool
	canWin = func(used, remaining int) bool {
		if memo[used] != 0 {
			return memo[used] == 1 // already solved this position
		}
		result := false
		for x := 1; x <= maxChoosableInteger; x++ {
			bit := 1 << uint(x-1)
			if used&bit != 0 {
				continue // x taken
			}
			// Immediate win, or push opponent into a proven loss.
			if x >= remaining || !canWin(used|bit, remaining-x) {
				result = true
				break // one winning move is enough
			}
		}
		if result { // cache the verdict for this mask
			memo[used] = 1
		} else {
			memo[used] = -1
		}
		return result
	}
	return canWin(0, desiredTotal)
}
```

### Dry Run

Example 1: `max = 10, desiredTotal = 11`. The recursion is identical to Approach 1; memoization only changes *how often* each mask is evaluated (once). Trace the root and one representative child:

| Call | used (bits) | remaining | evaluation | memo write | returns |
|------|-------------|-----------|------------|------------|---------|
| canWin(0, 11) | `0000000000` | 11 | try x=1: child below is a win for P2, so not a P1 win; try x=2..10 likewise all let P2 win | memo[0] = −1 | false |
| ↳ canWin(bit1, 10) | `0000000001` | 10 | x=10 unused and `10 ≥ 10` → immediate win for the mover (P2) | memo[1] = 1 | true |

Because every first move hands P2 an immediate-win position (some `y ≥ remaining` is still free), the root caches `memo[0] = −1` and returns `false`.

Result: `false` ✔

---

## Key Takeaways

- **Win/lose game recursion:** a position is a **win** iff at least one move leads to a **losing** position for the opponent; it is a **loss** iff every move leads to a winning position for the opponent. This single rule solves most impartial turn-based games.
- **State = the set, encoded as a bitmask.** When picks are without replacement and order does not affect the value, the used-set is the state — a bitmask makes it a cheap integer memo key. `max ≤ 20` is the tell that `2^max` states are intended.
- **Prune impossible games up front:** target `≤ 0` is an instant win; target above the full sum `n(n+1)/2` is an instant loss. These guards also keep `remaining` positive inside the recursion.
- **`remaining` is derivable from the mask**, so it need not be part of the memo key — a common space-saving observation for subset-DP games.

---

## Related Problems

- LeetCode #294 — Flip Game II (win/lose game recursion with memoization)
- LeetCode #486 — Predict the Winner (minimax over an array, take from ends)
- LeetCode #877 — Stone Game (minimax / parity)
- LeetCode #1140 — Stone Game II (game DP with a growing move limit)
- LeetCode #698 — Partition to K Equal Sum Subsets (bitmask subset search)
