# 0292 — Nim Game

> LeetCode #292 · Difficulty: Easy
> **Categories:** Math, Brainteaser, Game Theory

---

## Problem Statement

You are playing the following Nim Game with your friend:

- Initially, there is a heap of stones on the table.
- You and your friend will alternate taking turns, and **you go first**.
- On each turn, the person whose turn it is will remove 1 to 3 stones from the heap.
- The one who removes the last stone is the winner.

Given `n`, the number of stones in the heap, return `true` if you can win the game assuming both you and your friend play optimally, otherwise return `false`.

**Example 1:**
```
Input: n = 4
Output: false
Explanation: Here are the possible outcomes:
1. You remove 1 stone. Your friend removes 3 stones, including the last stone. Your friend wins.
2. You remove 2 stones. Your friend removes 2 stones, including the last stone. Your friend wins.
3. You remove 3 stones. Your friend removes the last stone. Your friend wins.
In all outcomes, your friend wins.
```

**Example 2:**
```
Input: n = 1
Output: true
```

**Example 3:**
```
Input: n = 2
Output: true
```

**Constraints:**
- `1 <= n <= 2³¹ - 1`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Adobe      | ★★☆☆☆ Low        | 2022          |
| Google     | ★★☆☆☆ Low        | 2022          |
| Microsoft  | ★★☆☆☆ Low        | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used
- **Game Theory / Pattern Recognition** — winning vs. losing positions in an impartial game → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)
- **Dynamic Programming (1D)** — `win[k]` from `win[k-1..k-3]` (illustrative) → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)
- **Number Theory (modulo)** — closed-form `n % 4` → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)

---

## Approaches Overview
| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | DP Bottom-Up | O(n) | O(n) | Discover the pattern; small n |
| 2 | Math Modulo (Optimal) | O(1) | O(1) | The intended answer |

---

## Approach 1 — DP Bottom-Up

### Intuition
A player wins from a heap of size `k` if there EXISTS a move (take 1, 2, or 3) that leaves the opponent in a LOSING position. So `win[k]` is true iff at least one of `win[k-1]`, `win[k-2]`, `win[k-3]` is false — that "false" is the trap handed to the opponent. Base: `win[0] = false` (no stones to take → the player to move already lost).

### Algorithm
1. `win[0] = false`.
2. For `k = 1..n`: set `win[k] = true` if any of `win[k-1]`, `win[k-2]`, `win[k-3]` is `false` (treating out-of-range indices as "not a winning move").
3. Return `win[n]`.

### Complexity
- **Time:** O(n) — one boolean computed per heap size (up to 3 lookups each).
- **Space:** O(n) — the `win` table. (For `n` up to 2³¹ this is only illustrative — use Approach 2 in practice.)

### Code
```go
func dpBottomUp(n int) bool {
	if n <= 0 {
		return false // no stones: the player to move cannot take the last stone
	}
	win := make([]bool, n+1) // win[k] = current player wins with k stones
	win[0] = false           // 0 stones → current player has no move → loses
	for k := 1; k <= n; k++ {
		for take := 1; take <= 3 && take <= k; take++ {
			if !win[k-take] { // opponent would be stuck in a losing position
				win[k] = true
				break // one winning move is enough
			}
		}
	}
	return win[n]
}
```

### Dry Run
`n = 4`

| k | check win[k-1], win[k-2], win[k-3] | any false? | win[k] |
|---|-------------------------------------|------------|--------|
| 0 | — (base) | — | false |
| 1 | win[0]=false | yes | **true** |
| 2 | win[1]=true, win[0]=false | yes | **true** |
| 3 | win[2]=true, win[1]=true, win[0]=false | yes | **true** |
| 4 | win[3]=true, win[2]=true, win[1]=true | no | **false** |

`win[4] = false` → the first player loses. ✓

---

## Approach 2 — Math Modulo (Optimal)

### Intuition
Running the DP by hand reveals the pattern: sizes 1, 2, 3 are wins (take everything); size 4 is a loss (whatever you take leaves 3, 2, or 1 — all wins for the opponent). From 5, 6, 7 you can always drop the opponent to exactly 4 (a loss for them), so those are wins; 8 is a loss again. **The losing positions are exactly the multiples of 4.** Winning strategy: always leave a multiple of 4 for the opponent — possible only when `n` is not already a multiple of 4.

### Algorithm
1. Return `n % 4 != 0`.

### Complexity
- **Time:** O(1) — a single modulo operation.
- **Space:** O(1) — no extra memory.

### Code
```go
func mathModulo(n int) bool {
	return n%4 != 0 // first player loses iff n is divisible by 4
}
```

### Dry Run
`n = 4`

| Step | Computation | Result |
|------|-------------|--------|
| 1 | `4 % 4` | `0` |
| 2 | `0 != 0` | `false` |

Returns `false` → first player loses. ✓ (Check: `n=1 → 1%4=1 ≠ 0 → true`; `n=2 → true`.)

---

## Key Takeaways
- **Impartial games** classify positions into P-positions (previous player wins / player-to-move loses) and N-positions (next player wins). Here P-positions are the multiples of 4.
- A brute DP can *reveal a pattern*; recognizing the period (4) collapses it to O(1).
- The "leave the opponent in a losing state" recurrence is the general template for take-away games; the specific losing period depends on the move set (1..3 → period 4).

---

## Related Problems
- LeetCode #294 — Flip Game II (game theory, win/lose recursion)
- LeetCode #375 — Guess Number Higher or Lower II (minimax DP)
- LeetCode #486 — Predict the Winner (game DP)
- LeetCode #877 — Stone Game (game theory / DP)
