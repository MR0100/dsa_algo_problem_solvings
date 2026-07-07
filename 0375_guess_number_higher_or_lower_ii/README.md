# 0375 — Guess Number Higher or Lower II

> LeetCode #375 · Difficulty: Medium
> **Categories:** Dynamic Programming, Math, Game Theory (Minimax)

---

## Problem Statement

We are playing the Guessing Game. The game will work as follows:

1. I pick a number between `1` and `n`.
2. You guess a number.
3. If you guess the right number, **you win the game**.
4. If you guess the wrong number, then I will tell you whether the number I picked is **higher or lower**, and you will continue guessing.
5. Every time you guess a wrong number `x`, you will pay `x` dollars. If you run out of money, **you lose the game**.

Given a particular `n`, return *the minimum amount of money you need to* **guarantee a win** *regardless of what number I pick.*

**Example 1:**

```
Input: n = 10
Output: 16
Explanation: The winning strategy is as follows:
- The range is [1,10]. Guess 7.
    - If this is my number, your total is $0. Otherwise, you pay $7.
    - If my number is higher, the range is [8,10]. Guess 9.
        - If this is my number, your total is $7. Otherwise, you pay $9.
        - If my number is higher, it must be 10. Guess 10. Your total is $7 + $9 = $16.
        - If my number is lower, it must be 8. Guess 8. Your total is $7 + $9 = $16.
    - If my number is lower, the range is [1,6]. Guess 3.
        - If this is my number, your total is $7. Otherwise, you pay $3.
        - If my number is higher, the range is [4,6]. Guess 5.
            - If this is my number, your total is $7 + $3 = $10. Otherwise, you pay $5.
            - If my number is higher, it must be 6. Guess 6. Your total is $7 + $3 + $5 = $15.
            - If my number is lower, it must be 4. Guess 4. Your total is $7 + $3 + $5 = $15.
        - If my number is lower, the range is [1,2]. Guess 1.
            - If this is my number, your total is $7 + $3 = $10. Otherwise, you pay $1.
            - If my number is higher, it must be 2. Guess 2. Your total is $7 + $3 + $1 = $11.
The worst case in all these scenarios is that you pay $16. Hence, you only need $16 to guarantee a win.
```

**Example 2:**

```
Input: n = 1
Output: 0
Explanation: There is only one possible number, so you can guess 1 and not have to pay anything.
```

**Example 3:**

```
Input: n = 2
Output: 1
Explanation: There are two possible numbers, 1 and 2.
- Guess 1.
    - If this is my number, your total is $0. Otherwise, you pay $1.
    - If my number is higher, it must be 2. Guess 2. Your total is $1.
The worst case is that you pay $1.
```

**Constraints:**

- `1 <= n <= 200`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Facebook   | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Interval Dynamic Programming** — `dp[lo][hi]` = min guaranteed cost for the range `[lo, hi]`, built from smaller sub-intervals → see [`/dsa/interval_dp.md`](/dsa/interval_dp.md)
- **Minimax / Game Theory** — the adversary picks the worse branch after each guess, so we minimise the maximum → see [`/dsa/game_theory.md`](/dsa/game_theory.md)
- **Memoization vs. tabulation** — the same recurrence solved top-down and bottom-up → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Top-down DP (memoized minimax) | O(n³) | O(n²) | Natural recursive framing |
| 2 | Bottom-up interval DP | O(n³) | O(n²) | Iterative, no recursion stack |

---

## Approach 1 — Top-Down DP (Memoized Minimax)

### Intuition
Define `cost(lo, hi)` = minimum money that **guarantees** finding any number in `[lo, hi]`. If we guess `x`, we pay `x` when wrong, and the adversary then forces us into whichever side is more expensive. So guessing `x` costs `x + max(cost(lo, x-1), cost(x+1, hi))`. We try every `x` and keep the cheapest — minimise the maximum. A range of size ≤ 1 costs `0` (we already know the answer).

### Algorithm
1. `memo[lo][hi]` caches solved ranges.
2. `solve(lo, hi)`: if `lo >= hi` return 0.
3. For `x = lo..hi`: `candidate = x + max(solve(lo, x-1), solve(x+1, hi))`. Keep the minimum. Store and return it.

### Complexity
- **Time:** O(n³) — O(n²) distinct `(lo, hi)` ranges, each scanning O(n) guesses.
- **Space:** O(n²) — memo table plus O(n) recursion depth.

### Code
```go
func dpTopDown(n int) int {
	// memo[lo][hi]; -1 marks "not computed yet". Size n+2 to allow x+1 up to n+1.
	memo := make([][]int, n+2)
	for i := range memo {
		memo[i] = make([]int, n+2)
		for j := range memo[i] {
			memo[i][j] = -1
		}
	}

	var solve func(lo, hi int) int
	solve = func(lo, hi int) int {
		if lo >= hi { // 0 or 1 candidate → no cost to be certain
			return 0
		}
		if memo[lo][hi] != -1 { // reuse a solved range
			return memo[lo][hi]
		}
		best := 1 << 30 // +infinity sentinel
		for x := lo; x <= hi; x++ {
			left := solve(lo, x-1)  // worst cost if the pick is below x
			right := solve(x+1, hi) // worst cost if the pick is above x
			worse := left           // adversary forces the more expensive side
			if right > worse {
				worse = right
			}
			cost := x + worse // pay x for the wrong guess, then the worse branch
			if cost < best {  // keep the cheapest guarantee
				best = cost
			}
		}
		memo[lo][hi] = best
		return best
	}
	return solve(1, n)
}
```

### Dry Run
`n = 2`, so we call `solve(1, 2)`:

| x | solve(lo,x-1) | solve(x+1,hi) | worse=max | cost = x + worse |
|---|---------------|---------------|-----------|------------------|
| 1 | solve(1,0)=0 | solve(2,2)=0 | 0 | 1 + 0 = **1** |
| 2 | solve(1,1)=0 | solve(3,2)=0 | 0 | 2 + 0 = 2 |

`best = min(1, 2) = 1`. Return `1`. ✓  (For `n = 10` the same recurrence yields `16`.)

---

## Approach 2 — Bottom-Up Interval DP (Optimal)

### Intuition
Same recurrence, filled iteratively so every sub-interval a guess depends on is already solved: `dp[lo][hi] = min over x of ( x + max(dp[lo][x-1], dp[x+1][hi]) )`. A range of length `L` depends only on ranges of length `< L`, so processing shorter ranges first (here, `lo` decreasing while `hi` increases) removes recursion.

### Algorithm
1. `dp` is `(n+2) × (n+2)`, all zero (empty/singleton ranges cost 0).
2. For `lo` from `n-1` down to `1`, for `hi` from `lo+1` to `n`:
   `dp[lo][hi] = min_{x in [lo,hi]} ( x + max(dp[lo][x-1], dp[x+1][hi]) )`.
3. Answer is `dp[1][n]`.

### Complexity
- **Time:** O(n³) — three nested loops over the range.
- **Space:** O(n²) — the DP table.

### Code
```go
func dpBottomUp(n int) int {
	// dp[lo][hi]; indices 0..n+1 so dp[x+1][hi] and dp[lo][x-1] stay in bounds.
	dp := make([][]int, n+2)
	for i := range dp {
		dp[i] = make([]int, n+2)
	}

	for lo := n - 1; lo >= 1; lo-- {
		for hi := lo + 1; hi <= n; hi++ {
			best := 1 << 30
			for x := lo; x <= hi; x++ {
				left := dp[lo][x-1]  // already computed (smaller hi)
				right := dp[x+1][hi] // already computed (larger lo)
				worse := left
				if right > worse {
					worse = right
				}
				cost := x + worse
				if cost < best {
					best = cost
				}
			}
			dp[lo][hi] = best
		}
	}
	return dp[1][n]
}
```

### Dry Run
`n = 2`. Table starts all zeros. Only interval processed: `lo=1, hi=2`:

| x | dp[lo][x-1] | dp[x+1][hi] | worse | cost = x + worse |
|---|-------------|-------------|-------|------------------|
| 1 | dp[1][0]=0 | dp[2][2]=0 | 0 | 1 + 0 = **1** |
| 2 | dp[1][1]=0 | dp[3][2]=0 | 0 | 2 + 0 = 2 |

`dp[1][2] = 1`. Answer `dp[1][2] = 1`. ✓

---

## Key Takeaways
- **"Guarantee a win / worst case" ⇒ minimax.** You minimise over your choices, the adversary maximises over the branch — hence `min_x ( x + max(left, right) )`.
- **This is interval DP**, not binary search: the optimal first guess is generally *not* the midpoint, because larger guesses cost more when wrong.
- **Cost is the guess value, not 1** — that asymmetry is what breaks naive binary search and forces the DP.
- Base case `lo >= hi` returns 0 (a known or empty range is free). Size the table `n+2` so `x-1` and `x+1` never index out of bounds.

---

## Related Problems
- LeetCode #374 — Guess Number Higher or Lower (the binary-search predecessor)
- LeetCode #312 — Burst Balloons (interval DP with a chosen split point)
- LeetCode #486 — Predict the Winner (minimax game DP)
- LeetCode #877 — Stone Game (game-theory DP)
