# Dynamic Programming — 1D

> A concept reference for the `/dsa/` library.
> Companion: [`/dsa/dynamic_programming.md`](/dsa/dynamic_programming.md) (general DP + 2D tables).

---

## What it is

**1D Dynamic Programming** solves a problem by defining a *single-index state*
`dp[i]` — "the answer for the prefix ending at (or of length) `i`" — and
computing each state from a **fixed, small set of earlier states** via a
recurrence. Because the state is one number (an index/position/amount), the
memo table is a one-dimensional array, giving `O(n)` time (or `O(n·k)` for a
constant/bounded `k` lookback) and `O(n)` space — usually reducible to `O(1)`.

The three ingredients you must be able to state out loud in an interview:

1. **State definition** — what exactly does `dp[i]` mean? (Most bugs are a
   fuzzy state definition, not a wrong loop.)
2. **Recurrence (transition)** — `dp[i] = f(dp[i-1], dp[i-2], …)` and *why*
   it's correct: the last decision partitions all solutions into cases.
3. **Base cases + answer location** — smallest states computed directly, and
   whether the final answer is `dp[n]`, `dp[n-1]`, or `max(dp[...])`.

Two equivalent implementation styles:

- **Top-down (memoized recursion):** write the natural recursion, cache
  results by index. Easier to derive; recursion depth `O(n)`.
- **Bottom-up (tabulation):** fill `dp[0..n]` in order with a loop. No stack,
  enables the rolling-variable space optimisation. This is the style you
  should be able to produce fastest.

### Overlapping subproblems + optimal substructure

DP applies when (a) the naive recursion recomputes the same subproblem
exponentially many times (e.g. `fib(n)` recursion tree), and (b) the optimal
answer for `i` is composed of optimal answers for smaller indices. If either
property is missing (e.g. subproblems don't overlap → plain divide & conquer;
future choices invalidate past ones → maybe greedy or search), 1D DP is the
wrong tool.

---

## How to recognise a 1D-DP problem (signals in the statement)

- **"Count the number of ways to …"** reach step `n`, decode a string, tile a
  board → counting DP (`dp[i] = sum of ways`). E.g. Climbing Stairs, Decode Ways.
- **"Maximum / minimum … over a sequence"** with a *contiguity or adjacency
  constraint* — max subarray sum, min cost to reach the end, "cannot pick two
  adjacent" (House Robber pattern).
- **"Can you reach / partition / segment …"** — feasibility DP where
  `dp[i] ∈ {true,false}` (Jump Game, Word Break).
- The input is a **linear sequence** (array, string, number `n`) and each
  position's answer plausibly depends only on **a few previous positions** or
  on **all previous positions via a max/min/sum** (`O(n²)` DP, sometimes
  optimisable).
- The obvious brute force is **exponential branching over choices at each
  index** ("take or skip", "jump 1 or 2", "decode 1 or 2 digits") — branching
  recursion over a prefix index is the classic tell.
- Constraints around `n ≤ 10⁴–10⁶` with an exponential brute force → the
  intended solution is almost certainly `O(n)`/`O(n²)` DP.
- **Greedy counter-signal:** if a locally best choice can be proven safe
  (exchange argument), greedy may beat DP (Jump Game II). When you can't prove
  it, DP is the safe interview answer; mention the greedy as the follow-up.

---

## General templates (Go)

### Template 1 — bottom-up tabulation (constant lookback)

```go
// dp[i] = <precise English definition of the state>
func solve1D(nums []int) int {
    n := len(nums)
    dp := make([]int, n+1) // often size n+1 so dp[0] is an "empty prefix" base

    // 1) Base cases: smallest subproblems answered directly.
    dp[0] = /* answer for empty prefix */
    dp[1] = /* answer for first element */

    // 2) Transition: each state from a FIXED set of earlier states.
    for i := 2; i <= n; i++ {
        // last-decision case analysis:
        //   option A: the solution for i ends with choice A -> uses dp[i-1]
        //   option B: the solution for i ends with choice B -> uses dp[i-2]
        dp[i] = best(dp[i-1] /* +cost of A */, dp[i-2] /* +cost of B */)
    }

    // 3) Answer: usually dp[n]; sometimes max over all dp[i] (e.g. Kadane).
    return dp[n]
}
```

### Template 2 — space-optimised rolling variables

When `dp[i]` depends only on the previous `k` states (constant `k`), keep just
those `k` values:

```go
func solve1DConstantSpace(n int) int {
    prev2, prev1 := base0, base1 // dp[i-2], dp[i-1]
    for i := 2; i <= n; i++ {
        cur := combine(prev1, prev2) // the recurrence
        prev2, prev1 = prev1, cur    // slide the window forward
    }
    return prev1
}
```

### Template 3 — O(n²) lookback ("best over all earlier j")

Used when position `i` can extend *any* earlier position `j` (LIS, Jump Game
DP, Word Break):

```go
func solveQuadratic(nums []int) int {
    n := len(nums)
    dp := make([]int, n)
    for i := 0; i < n; i++ {
        dp[i] = baseValue // e.g. 1 for LIS: the element alone
        for j := 0; j < i; j++ {
            if canExtend(j, i) {           // problem-specific predicate
                dp[i] = best(dp[i], dp[j]+gain(i))
            }
        }
    }
    return maxOf(dp) // answer often max over all states, not dp[n-1]
}
```

### Template 4 — top-down memoization

```go
func solveTopDown(n int) int {
    memo := make([]int, n+1)
    for i := range memo {
        memo[i] = -1 // sentinel meaning "not computed"; careful if -1 is a
    }                // legal answer — use a bool array or map instead
    var rec func(i int) int
    rec = func(i int) int {
        if i <= 1 {            // base cases first
            return base(i)
        }
        if memo[i] != -1 {     // cache hit: each state computed once
            return memo[i]
        }
        memo[i] = combine(rec(i-1), rec(i-2))
        return memo[i]
    }
    return rec(n)
}
```

### State-machine variant (still 1D)

Some problems need 2–4 *named* running values per index instead of one —
e.g. "max profit while holding a stock" vs "while not holding". The index is
still the only dimension; the states are constants:

```go
hold, free := -prices[0], 0
for _, p := range prices[1:] {
    hold = max(hold, free-p) // buy today or keep holding
    free = max(free, hold+p) // sell today or stay out
}
return free
```

---

## Worked example — Climbing Stairs (LeetCode #70)

**Problem:** you climb `n` stairs taking 1 or 2 steps at a time; count the
distinct ways to reach the top.

**1. State:** `dp[i]` = number of distinct ways to stand on step `i`.

**2. Recurrence:** the *last* move to reach step `i` was either a 1-step from
`i-1` or a 2-step from `i-2`. These cases are disjoint and exhaustive, so:

```
dp[i] = dp[i-1] + dp[i-2]
```

**3. Base cases:** `dp[0] = 1` (one way: do nothing), `dp[1] = 1`.

```go
// climbStairs counts ways to climb n steps taking 1 or 2 at a time.
// Time: O(n) — one pass. Space: O(1) — rolling variables.
func climbStairs(n int) int {
    if n <= 1 {
        return 1
    }
    prev2, prev1 := 1, 1 // dp[0], dp[1]
    for i := 2; i <= n; i++ {
        prev2, prev1 = prev1, prev1+prev2 // dp[i] = dp[i-1] + dp[i-2]
    }
    return prev1 // dp[n]
}
```

**Step-by-step trace for `n = 5`:**

| i | dp[i-2] (`prev2`) | dp[i-1] (`prev1`) | dp[i] = dp[i-1]+dp[i-2] | meaning |
|---|------|------|------|---------|
| 2 | 1 | 1 | **2** | `1+1`, `2` |
| 3 | 1 | 2 | **3** | `1+1+1`, `1+2`, `2+1` |
| 4 | 2 | 3 | **5** | 5 sequences |
| 5 | 3 | 5 | **8** | 8 sequences ← answer |

The recursion tree for a naive `f(5)` would call `f(3)` twice and `f(2)` three
times — the *overlapping subproblems* the memo/table eliminates. Note the
answer is the Fibonacci sequence shifted by one: many 1D counting DPs reduce
to Fibonacci-like recurrences.

---

## Common pitfalls (and how to avoid them)

1. **Vague state definition.** "dp[i] is the answer at i" — at, ending at, or
   up to? For Max Subarray, `dp[i]` must be "best sum of a subarray **ending
   exactly at** `i`", and the final answer is `max(dp)`, not `dp[n-1]`.
   *Fix:* write the definition as a full sentence before coding.
2. **Off-by-one between "index `i`" and "prefix of length `i`".** Sizing the
   array `n` vs `n+1` and mapping `dp[i]` ↔ `s[i-1]` is the #1 bug source in
   string DPs like Decode Ways. *Fix:* prefer `dp[0] =` empty-prefix base and
   `dp[i]` = "first `i` characters"; be consistent everywhere.
3. **Wrong or missing base cases.** `dp[0] = 1` (one *empty* way) vs
   `dp[0] = 0` changes every counting answer. *Fix:* hand-compute `dp[0..2]`
   from the problem statement and check the loop reproduces `dp[2]`.
4. **Transition reads a state not yet computed** (or, after space
   optimisation, one already overwritten). In-place row updates (Pascal's
   Triangle II, 0/1-knapsack-style) often must iterate **right-to-left** to
   read old values. *Fix:* for each `dp[i] = f(dp[j])`, verify `j`'s value at
   that moment is the intended generation.
5. **Memo sentinel collides with a real answer.** `-1` as "uncomputed" breaks
   when `-1` is a legal result (min-cost DPs with negatives). *Fix:* a
   parallel `seen []bool` or `map[int]int`.
6. **Forgetting where the answer lives.** `dp[n]`? `dp[n-1]`? `max(dp[...])`?
   Kadane and LIS need the max over all states. *Fix:* decide this when you
   define the state, not after the loop.
7. **Using DP when greedy is required (or vice-versa).** `O(n²)` DP for Jump
   Game II TLEs on large inputs where the `O(n)` greedy passes. *Fix:* present
   DP as the safe baseline, then ask whether a local choice is provably safe.
8. **Recursion depth / TLE in top-down Go.** `n = 10⁶` recursive calls risk
   stack growth and constant-factor overhead. *Fix:* convert to bottom-up for
   large `n`.
9. **Integer overflow in counting DPs.** Way-counts explode; many problems ask
   for the answer "mod 1e9+7". *Fix:* apply the modulus inside the transition,
   every time you add or multiply.
10. **Optimising space before the table version works.** Rolling variables are
    a refactor, not a starting point. *Fix:* get the `O(n)` array correct and
    tested, then collapse it.

---

## Problems in this repo

Problems currently in the repo whose solutions use a 1D DP state
(problems 0131+ will be appended in a later pass):

- [0022 — Generate Parentheses](../0022_generate_parentheses/README.md) — `dp[k]` = all valid strings with `k` pairs (Catalan-style build-up).
- [0032 — Longest Valid Parentheses](../0032_longest_valid_parentheses/README.md) — `dp[i]` = length of the longest valid substring **ending at** `i`.
- [0045 — Jump Game II](../0045_jump_game_ii/README.md) — `dp[i]` = min jumps to reach `i` (O(n²) baseline; greedy is optimal).
- [0053 — Maximum Subarray](../0053_maximum_subarray/README.md) — Kadane: `dp[i]` = best sum of a subarray ending at `i`; answer is `max(dp)`.
- [0055 — Jump Game](../0055_jump_game/README.md) — feasibility DP: `dp[i]` = "is index `i` reachable?" (greedy is optimal).
- [0070 — Climbing Stairs](../0070_climbing_stairs/README.md) — the canonical Fibonacci recurrence (worked example above).
- [0091 — Decode Ways](../0091_decode_ways/README.md) — `dp[i]` = ways to decode the first `i` chars; take 1 or 2 digits as the last move.
- [0096 — Unique Binary Search Trees](../0096_unique_binary_search_trees/README.md) — `dp[n]` = Catalan number via `Σ dp[left]·dp[right]`.
- [0119 — Pascal's Triangle II](../0119_pascals_triangle_ii/README.md) — in-place 1D row update; must sweep right-to-left (pitfall 4).
- [0120 — Triangle](../0120_triangle/README.md) — min path sum with a rolling 1D row (`dp[col]`), bottom-up over rows.
- [0121 — Best Time to Buy and Sell Stock](../0121_best_time_to_buy_and_sell_stock/README.md) — running-min / best-profit-so-far state per index.
- [0122 — Best Time to Buy and Sell Stock II](../0122_best_time_to_buy_and_sell_stock_ii/README.md) — hold/free two-state machine over one index.
- [0123 — Best Time to Buy and Sell Stock III](../0123_best_time_to_buy_and_sell_stock_iii/README.md) — four-state machine (two transactions), still one pass.

Related but **not** 1D: string/grid problems with two indices
(`dp[i][j]`) — Longest Palindromic Substring, Edit Distance, Unique Paths —
live in [`/dsa/dynamic_programming.md`](/dsa/dynamic_programming.md).
