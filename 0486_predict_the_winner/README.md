# 0486 — Predict the Winner

> LeetCode #486 · Difficulty: Medium
> **Categories:** Array, Math, Dynamic Programming, Recursion, Game Theory

---

## Problem Statement

You are given an integer array `nums`. Two players are playing a game with this array: player 1 and player 2.

Player 1 and player 2 take turns, with player 1 starting first. Both players start the game with a score of `0`. At each turn, the player takes one of the numbers from either end of the array (i.e., `nums[0]` or `nums[nums.length - 1]`) which reduces the size of the array by `1`. The player adds the chosen number to their score. The game ends when there are no more elements in the array.

Return `true` if Player 1 can win the game. If the scores of both players are equal, then player 1 is still the winner, and you should also return `true`. You may assume that both players are playing optimally.

**Example 1:**

```
Input: nums = [1,5,2]
Output: false
Explanation: Initially, player 1 can choose between 1 and 2.
If he chooses 2 (or 1), then player 2 can choose from 1 (or 2) and 5.
If player 2 chooses 5, then player 1 will be left with 1 (or 2).
So, final score of player 1 is 1 + 2 = 3, and player 2 is 5.
Hence, player 1 will never be the winner and you need to return false.
```

**Example 2:**

```
Input: nums = [1,5,233,7]
Output: true
Explanation: Player 1 first chooses 1. Then player 2 has to choose between 5 and 7.
No matter which number player 2 chooses, player 1 can choose 233.
Finally, player 1 has more score (234) than player 2 (12), so you need to return true representing player1 can win.
```

**Constraints:**

- `1 <= nums.length <= 20`
- `0 <= nums[i] <= 10^7`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Adobe      | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Game Theory (Minimax / zero-sum)** — two players alternate optimally on a shared resource; collapsing the two scores into one signed "difference" is the canonical minimax reduction → see [`/dsa/game_theory.md`](/dsa/game_theory.md)
- **Interval DP** — the state is a contiguous sub-array `nums[lo..hi]`, and each interval's answer is built from strictly shorter intervals — the defining shape of interval dynamic programming → see [`/dsa/interval_dp.md`](/dsa/interval_dp.md)
- **2D Dynamic Programming** — the memo/table is indexed by the two endpoints `(lo, hi)`, giving an n×n triangular state space → see [`/dsa/dynamic_programming_2d.md`](/dsa/dynamic_programming_2d.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force Recursion | O(2ⁿ) | O(n) | Explains the minimax reduction; fine for n ≤ 20 but wasteful |
| 2 | Top-Down DP (Memoised) | O(n²) | O(n²) | Easiest optimal version — add a cache to the recursion |
| 3 | Bottom-Up Interval DP | O(n²) | O(n²) | Iterative, no stack; canonical interval-DP fill order |
| 4 | Space-Optimised 1D DP (Optimal) | O(n²) | O(n) | When memory matters; same time, one rolling row |

---

## Approach 1 — Brute Force Recursion

### Intuition

Both players play optimally on a zero-sum game, so instead of juggling two running totals we track a single signed number: the **net advantage** of whoever is about to move on the current sub-array `nums[lo..hi]`. If the mover takes the left end `nums[lo]`, they immediately earn `nums[lo]`, then hand a smaller array to the opponent, who will now play optimally and rack up their *own* best difference on `[lo+1..hi]`. Because roles flip, the opponent's advantage counts against the mover — so we **subtract** it. The mover simply picks whichever end maximises `pick − (opponent's best difference)`. Player 1 wins iff the difference over the whole array is `≥ 0` (ties favour Player 1).

### Algorithm

1. Define `score(lo, hi)` = best `(mover − opponent)` difference on `nums[lo..hi]`.
2. Base case: if `lo == hi`, one number remains → `return nums[lo]`.
3. Otherwise return `max(nums[lo] − score(lo+1, hi), nums[hi] − score(lo, hi-1))`.
4. The answer is `score(0, n-1) >= 0`.

### Complexity

- **Time:** O(2ⁿ) — every call spawns two recursive calls and nothing is cached, so the call tree doubles with each level.
- **Space:** O(n) — recursion depth equals the number of picks (one per array element).

### Code

```go
func bruteForce(nums []int) bool {
	// score returns the best achievable (mover − opponent) difference on nums[lo..hi].
	var score func(lo, hi int) int
	score = func(lo, hi int) int {
		if lo == hi {
			return nums[lo] // last remaining number goes straight to the mover
		}
		// Take the left end: gain nums[lo], then subtract the opponent's best
		// difference on the smaller range (roles flip → subtract).
		takeLeft := nums[lo] - score(lo+1, hi)
		// Take the right end symmetrically.
		takeRight := nums[hi] - score(lo, hi-1)
		if takeLeft > takeRight {
			return takeLeft // prefer the end that yields the larger net advantage
		}
		return takeRight
	}
	return score(0, len(nums)-1) >= 0 // ≥ 0 ⇒ Player 1 at least ties ⇒ wins
}
```

### Dry Run

Example 1: `nums = [1, 5, 2]`. We compute `score(0, 2)`.

| Call | lo | hi | takeLeft | takeRight | returns |
|------|----|----|----------|-----------|---------|
| score(2,2) | 2 | 2 | — | — | 2 |
| score(1,1) | 1 | 1 | — | — | 5 |
| score(0,0) | 0 | 0 | — | — | 1 |
| score(1,2) | 1 | 2 | 5 − score(2,2) = 5−2 = 3 | 2 − score(1,1) = 2−5 = −3 | max(3,−3) = 3 |
| score(0,1) | 0 | 1 | 1 − score(1,1) = 1−5 = −4 | 5 − score(0,0) = 5−1 = 4 | max(−4,4) = 4 |
| score(0,2) | 0 | 2 | 1 − score(1,2) = 1−3 = −2 | 2 − score(0,1) = 2−4 = −2 | max(−2,−2) = −2 |

`score(0,2) = −2 < 0` → return **false** ✔ (Player 1 trails by 2 no matter what.)

---

## Approach 2 — Top-Down DP (Memoised)

### Intuition

`score(lo, hi)` depends only on the pair `(lo, hi)`, and there are just `n²` such pairs. The exponential blow-up in brute force came purely from recomputing the same sub-array difference again and again. Cache each result the first time it's computed and the `2ⁿ` call tree collapses into `n²` distinct states, each solved once.

### Algorithm

1. Keep `memo[lo][hi]` for the value and `seen[lo][hi]` for whether it is filled.
2. `score(lo, hi)`: if `lo == hi` return `nums[lo]`; if `seen[lo][hi]` return the cache.
3. Otherwise compute `max(nums[lo] − score(lo+1,hi), nums[hi] − score(lo,hi-1))`, store it, mark `seen`, and return.
4. Answer is `score(0, n-1) >= 0`.

### Complexity

- **Time:** O(n²) — `n²` states, each doing O(1) work once memoised.
- **Space:** O(n²) for the memo table plus O(n) recursion stack.

### Code

```go
func dpTopDown(nums []int) bool {
	n := len(nums)
	memo := make([][]int, n)  // memo[lo][hi] = best difference on nums[lo..hi]
	seen := make([][]bool, n) // seen[lo][hi] = has this state been computed?
	for i := range memo {
		memo[i] = make([]int, n)
		seen[i] = make([]bool, n)
	}
	var score func(lo, hi int) int
	score = func(lo, hi int) int {
		if lo == hi {
			return nums[lo] // base case: single number
		}
		if seen[lo][hi] {
			return memo[lo][hi] // already solved this sub-array
		}
		takeLeft := nums[lo] - score(lo+1, hi)  // take left end
		takeRight := nums[hi] - score(lo, hi-1) // take right end
		best := takeLeft
		if takeRight > best {
			best = takeRight
		}
		seen[lo][hi] = true // memoise before returning
		memo[lo][hi] = best
		return best
	}
	return score(0, n-1) >= 0
}
```

### Dry Run

Example 1: `nums = [1, 5, 2]`. Recursion visits the same six `(lo,hi)` states as Approach 1, but now each is written to the cache exactly once.

| Order computed | State (lo,hi) | Value stored | Cache hits later |
|----------------|---------------|--------------|------------------|
| 1 | (1,2) via (2,2),(1,1) | 3 | — |
| 2 | (2,2) | 2 | reused by (1,2) |
| 3 | (1,1) | 5 | reused by (1,2) and (0,1) |
| 4 | (0,1) via (1,1),(0,0) | 4 | (1,1) served from cache |
| 5 | (0,0) | 1 | reused by (0,1) |
| 6 | (0,2) | −2 | reads cached (1,2)=3, (0,1)=4 |

`memo[0][2] = −2 < 0` → **false** ✔. No state is computed twice.

---

## Approach 3 — Bottom-Up Interval DP

### Intuition

`dp[lo][hi]` is the best `(mover − opponent)` difference on `nums[lo..hi]`. A length-1 interval is trivial: `dp[i][i] = nums[i]`. Every longer interval reads only from strictly shorter ones (`[lo+1..hi]` and `[lo..hi-1]`). So if we fill intervals from shortest to longest, every dependency is already computed — the textbook interval-DP order — and we never touch a recursion stack.

### Algorithm

1. Initialise `dp[i][i] = nums[i]` for all `i`.
2. For `length` = 2..n, and each `lo` with `hi = lo + length − 1`:
   `dp[lo][hi] = max(nums[lo] − dp[lo+1][hi], nums[hi] − dp[lo][hi-1])`.
3. Return `dp[0][n-1] >= 0`.

### Complexity

- **Time:** O(n²) — the outer length loop and inner start loop together visit each of the `~n²/2` intervals once.
- **Space:** O(n²) — the triangular table of interval answers.

### Code

```go
func dpBottomUp(nums []int) bool {
	n := len(nums)
	dp := make([][]int, n) // dp[lo][hi] = best difference on nums[lo..hi]
	for i := range dp {
		dp[i] = make([]int, n)
		dp[i][i] = nums[i] // base case: a single number is a pure gain for the mover
	}
	// Grow the interval length; every state below depends only on shorter ones.
	for length := 2; length <= n; length++ {
		for lo := 0; lo+length-1 < n; lo++ {
			hi := lo + length - 1
			takeLeft := nums[lo] - dp[lo+1][hi]  // take the left end
			takeRight := nums[hi] - dp[lo][hi-1] // take the right end
			if takeLeft > takeRight {
				dp[lo][hi] = takeLeft
			} else {
				dp[lo][hi] = takeRight
			}
		}
	}
	return dp[0][n-1] >= 0
}
```

### Dry Run

Example 1: `nums = [1, 5, 2]`, `n = 3`.

| Phase | Cells set | Computation | Result |
|-------|-----------|-------------|--------|
| length 1 | dp[0][0], dp[1][1], dp[2][2] | = nums[i] | 1, 5, 2 |
| length 2, lo=0 | dp[0][1] | max(1 − dp[1][1], 5 − dp[0][0]) = max(1−5, 5−1) | max(−4, 4) = 4 |
| length 2, lo=1 | dp[1][2] | max(5 − dp[2][2], 2 − dp[1][1]) = max(5−2, 2−5) | max(3, −3) = 3 |
| length 3, lo=0 | dp[0][2] | max(1 − dp[1][2], 2 − dp[0][1]) = max(1−3, 2−4) | max(−2, −2) = −2 |

`dp[0][2] = −2 < 0` → **false** ✔.

---

## Approach 4 — Space-Optimised 1D DP (Optimal)

### Intuition

In the 2D fill, `dp[lo][hi]` reads exactly two neighbours: `dp[lo+1][hi]` (same column `hi`, from the previous shorter length) and `dp[lo][hi-1]` (one column left, already updated this pass). Reuse a single array `dp` indexed by `hi` and iterate `lo` **downward**. When we reach `dp[hi]`, its *current* contents are still the `lo+1` value (not yet overwritten), and `dp[hi-1]` already holds the freshly computed `lo` value. Those are precisely the two inputs — so one row of memory replaces the whole table.

### Algorithm

1. Start with `dp[i] = nums[i]` (all length-1 intervals `dp[lo][lo]`).
2. For `lo` from `n−2` down to `0`, and `hi` from `lo+1` to `n−1`:
   `dp[hi] = max(nums[lo] − dp[hi], nums[hi] − dp[hi-1])`
   (`dp[hi]` before the write is old `dp[lo+1][hi]`; `dp[hi-1]` is new `dp[lo][hi-1]`).
3. Return `dp[n-1] >= 0`.

### Complexity

- **Time:** O(n²) — identical double loop over `(lo, hi)`.
- **Space:** O(n) — a single rolling row instead of the `n×n` table.

### Code

```go
func dpOneDim(nums []int) bool {
	n := len(nums)
	dp := make([]int, n) // dp[hi] doubles as dp[lo][hi] for the current lo
	copy(dp, nums)       // length-1 intervals: dp[i] = nums[i]
	// Sweep lo from the second-to-last start downward so shorter intervals are ready.
	for lo := n - 2; lo >= 0; lo-- {
		for hi := lo + 1; hi < n; hi++ {
			// dp[hi]   still holds the value for start lo+1 (previous outer pass).
			// dp[hi-1] already holds the value for start lo   (this pass).
			takeLeft := nums[lo] - dp[hi]    // take left end: opponent solves [lo+1..hi]
			takeRight := nums[hi] - dp[hi-1] // take right end: opponent solves [lo..hi-1]
			if takeLeft > takeRight {
				dp[hi] = takeLeft
			} else {
				dp[hi] = takeRight
			}
		}
	}
	return dp[n-1] >= 0
}
```

### Dry Run

Example 1: `nums = [1, 5, 2]`, `n = 3`. Start `dp = [1, 5, 2]`.

| lo | hi | dp before | takeLeft = nums[lo]−dp[hi] | takeRight = nums[hi]−dp[hi-1] | dp[hi] ← | dp after |
|----|----|-----------|----------------------------|-------------------------------|----------|----------|
| 1 | 2 | [1, 5, 2] | 5 − dp[2]=5−2 = 3 | 2 − dp[1]=2−5 = −3 | 3 | [1, 5, 3] |
| 0 | 1 | [1, 5, 3] | 1 − dp[1]=1−5 = −4 | 5 − dp[0]=5−1 = 4 | 4 | [1, 4, 3] |
| 0 | 2 | [1, 4, 3] | 1 − dp[2]=1−3 = −2 | 2 − dp[1]=2−4 = −2 | −2 | [1, 4, −2] |

`dp[n-1] = dp[2] = −2 < 0` → **false** ✔. Note `dp[2]` was `3` (the `lo=1` answer) exactly when the `lo=0, hi=2` step needed it as `dp[lo+1][hi]`.

---

## Key Takeaways

- **Two-player zero-sum → track one signed difference, not two scores.** `mover − opponent` turns "maximise mine, minimise theirs" into a single `max(pick − recurse(rest))`. This minimax-on-difference trick recurs across game problems (Stone Game, Nim-like takeaways).
- **"Pick from either end of an array" is an interval-DP fingerprint.** The state is `(lo, hi)`; the answer for an interval is built from the two intervals one element shorter.
- **Ties go to Player 1 ⇒ compare against `>= 0`**, not `> 0`. Read the tiebreak rule carefully.
- **Optimisation ladder:** brute recursion → memoise (top-down) → invert to bottom-up table → drop a dimension when each cell needs only its neighbours. The 2D→1D collapse hinges on iterating the dropped dimension in the direction that keeps both needed values alive.

---

## Related Problems

- LeetCode #877 — Stone Game (same pick-from-ends minimax; even length guarantees P1 win)
- LeetCode #1140 — Stone Game II (interval/suffix DP with a variable take size)
- LeetCode #1406 — Stone Game III (three-at-a-time takeaway, difference DP)
- LeetCode #375 — Guess Number Higher or Lower II (interval DP, minimax cost)
- LeetCode #464 — Can I Win (game theory with bitmask memoisation)
