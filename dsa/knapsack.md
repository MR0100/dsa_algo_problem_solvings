# Knapsack DP Family (0/1, Unbounded, Subset-Sum)

> **Core idea:** make the **remaining capacity** a dimension of the DP state, and let each item either be taken or skipped.
> **The one rule that separates the variants:** iterate capacity **descending** for 0/1 (each item once), **ascending** for unbounded (item reusable).
> **Complexity:** O(n · C) time, O(C) space with the rolling-array trick — where n = #items, C = capacity.

---

## What it is

The knapsack problem asks: given items each with a **weight** and a **value**, and a bag of fixed **capacity**, which subset maximises total value without exceeding capacity? It is the archetype of a huge family of DP problems where the answer depends on a **budget that shrinks as you commit resources**. Recognising a problem as "knapsack-shaped" instantly gives you the recurrence, the loop structure, and the space optimisation.

The family members you must know:

| Variant | Each item | Question | This repo |
|---------|-----------|----------|-----------|
| **0/1 knapsack** | used **at most once** | max value within capacity | (foundation for all below) |
| **Unbounded knapsack** | used **any number of times** | max value / min count / #ways with reuse | #322 Coin Change, #518 Coin Change II |
| **Subset-sum** | in or out | *can* a subset hit exactly S? (bool) | #416 Partition Equal Subset Sum |
| **Subset-sum count** | in or out | *how many* subsets hit exactly S? | #494 Target Sum |
| **Partition / k-partition** | in or out | split multiset into groups of equal sum | #416 (2-way), #473 Matchsticks (4-way) |
| **2-D-capacity knapsack** | used once | max under **two independent** limits | #474 Ones and Zeroes |

### The unifying insight — capacity as a DP dimension

Define `dp[c]` = the best achievable answer using some prefix of the items with **capacity `c` remaining/available**. Each item transitions `dp[c]` from `dp[c - weight]` (take it) versus `dp[c]` (skip it). Because the only thing that matters going forward is *how much capacity is left* and *which items remain*, the whole exponential subset search collapses to an `n × C` table — this is the textbook example of overlapping subproblems + optimal substructure.

---

## When to recognise it — signals in the problem statement

| Signal | Which variant |
|--------|---------------|
| "pick a subset to **maximise value** within a weight/budget limit" | 0/1 knapsack |
| "**can** the array be split into two equal-sum halves?" / "reach exactly target sum" | subset-sum (bool) |
| "assign `+`/`−` signs so the expression equals a target" | subset-sum count (algebraic reduction, see #494) |
| "split into **k groups** of equal sum" / "form a square from matchsticks" | k-partition (subset-sum repeated k times; often backtracking) |
| "coins/items may be used **unlimited** times" — fewest coins, #ways to make an amount | unbounded knapsack |
| "at most `m` of resource A **and** `n` of resource B" | 2-D-capacity 0/1 knapsack |
| a **target number** small enough that `n · target` fits in time/memory | capacity-as-dimension is viable |

**When *not* to use it:** the "capacity" (target sum) is astronomically large (e.g. 10¹⁸) → pseudo-polynomial DP won't fit; look for math/greedy instead. Items have real-valued weights → classic knapsack DP needs integer capacities (or scaling).

---

## General templates (Go)

### 1. Subset-sum — "can we hit exactly `target`?" (0/1, boolean)

```go
// canPartitionToSum reports whether some subset of nums sums to exactly target.
func canPartitionToSum(nums []int, target int) bool {
    dp := make([]bool, target+1) // dp[s] = "some subset seen so far sums to s"
    dp[0] = true                 // the empty subset sums to 0

    for _, num := range nums {
        // DESCENDING: guarantees each item is used at most once this round —
        // dp[s-num] on the right still refers to the PREVIOUS item's table,
        // not a value we just updated with num. (Ascending would let one item
        // be reused, which is a different, unbounded problem.)
        for s := target; s >= num; s-- {
            if dp[s-num] {
                dp[s] = true // reachable: take `num` on top of a subset summing to s-num
            }
        }
    }
    return dp[target]
}
```

### 2. Subset-sum **count** — "how many subsets sum to `target`?" (0/1)

Swap `bool` for an integer count; `+=` instead of `||`. This is the engine behind #494 Target Sum.

```go
// countSubsetsWithSum returns the number of subsets of nums summing to target.
func countSubsetsWithSum(nums []int, target int) int {
    dp := make([]int, target+1) // dp[s] = number of subsets summing to s
    dp[0] = 1                   // one way to make 0: pick nothing

    for _, num := range nums {
        for s := target; s >= num; s-- { // DESCENDING → each item counted once
            dp[s] += dp[s-num]
        }
    }
    return dp[target]
}
```

> **#494 reduction:** partition into a `+` group (sum `P`) and `−` group (sum `N`). From `P − N = target` and `P + N = total`, we get `P = (total + target) / 2`. So the count of valid sign assignments equals the number of subsets summing to `P`. If `total + target` is odd or `|target| > total`, the answer is `0`.

### 3. Classic 0/1 knapsack — maximise value within capacity `W`

```go
// knapsack01 returns the max total value with total weight <= W, each item used once.
func knapsack01(weights, values []int, W int) int {
    dp := make([]int, W+1) // dp[c] = best value achievable with capacity c
    for i := range weights {
        // DESCENDING so item i is considered at most once.
        for c := W; c >= weights[i]; c-- {
            take := dp[c-weights[i]] + values[i]
            if take > dp[c] {
                dp[c] = take // better to take item i here
            }
        }
    }
    return dp[W]
}
```

### 4. Unbounded knapsack — items reusable (min count / max value / #ways)

The **only** change from 0/1 is the inner loop direction: **ascending**, so `dp[c-w]` may already include the current item, allowing it to be picked again.

```go
// coinChangeMin returns the fewest coins summing to amount (unbounded), or -1.
func coinChangeMin(coins []int, amount int) int {
    const INF = 1 << 30
    dp := make([]int, amount+1) // dp[a] = fewest coins to make amount a
    for a := 1; a <= amount; a++ {
        dp[a] = INF
    }
    dp[0] = 0

    for _, coin := range coins {
        // ASCENDING → coin may be reused any number of times.
        for a := coin; a <= amount; a++ {
            if dp[a-coin]+1 < dp[a] {
                dp[a] = dp[a-coin] + 1
            }
        }
    }
    if dp[amount] >= INF {
        return -1
    }
    return dp[amount]
}

// coinChangeWays returns the number of combinations that make amount (order-independent).
func coinChangeWays(coins []int, amount int) int {
    dp := make([]int, amount+1) // dp[a] = #combinations summing to a
    dp[0] = 1

    // Coins in the OUTER loop → combinations (each coin type once), not permutations.
    // (Swapping the loop order counts ordered sequences instead — a different problem.)
    for _, coin := range coins {
        for a := coin; a <= amount; a++ { // ASCENDING for reuse
            dp[a] += dp[a-coin]
        }
    }
    return dp[amount]
}
```

> **Loop-order subtlety (unbounded counting):** *coin outer, amount inner* counts **combinations** (`{1,2}` = `{2,1}`). *amount outer, coin inner* counts **permutations** (ordered) — that is the "Combination Sum IV" / staircase problem, not Coin Change II. Memorise which one you want.

### 5. 2-D-capacity 0/1 knapsack — two independent limits (used by #474)

Each item consumes from **two** budgets at once (here: zeros and ones). Add a dimension; iterate **both** capacities descending.

```go
// findMaxForm returns the largest subset of strs using <= m zeros and <= n ones.
func findMaxForm(strs []string, m, n int) int {
    // dp[i][j] = max #strings using at most i zeros and j ones.
    dp := make([][]int, m+1)
    for i := range dp {
        dp[i] = make([]int, n+1)
    }

    for _, s := range strs {
        zeros, ones := count01(s) // this item's two "weights"; value = 1
        // BOTH capacities descending → each string used at most once.
        for i := m; i >= zeros; i-- {
            for j := n; j >= ones; j-- {
                if dp[i-zeros][j-ones]+1 > dp[i][j] {
                    dp[i][j] = dp[i-zeros][j-ones] + 1
                }
            }
        }
    }
    return dp[m][n]
}

func count01(s string) (zeros, ones int) {
    for _, ch := range s {
        if ch == '0' {
            zeros++
        } else {
            ones++
        }
    }
    return
}
```

### Rolling-array space optimisation — why it works

All templates above already use a **1-D (or minimal-D) rolling array** instead of the full `n × C` table. The 2-D recurrence `dp[i][c]` only reads row `i-1`, so one row suffices *if* you traverse `c` in the direction that preserves the "previous row" semantics:

- **Descending `c`** → `dp[c-w]` still holds the *old* (row `i-1`) value → each item used **once** (0/1).
- **Ascending `c`** → `dp[c-w]` may hold the *new* (row `i`) value → item reusable (**unbounded**).

That single directional choice is the entire distinction between the two most important variants — remember it and you can reconstruct every template above.

---

## Worked example — step-by-step trace

**Subset-sum (#416 shape):** can `nums = [1, 5, 11, 5]` be split into two equal-sum halves? Total = 22, so target = 11. We fill `dp[0..11]` (`dp[s]` = "some subset sums to s"), starting `dp[0]=true`.

Initial: `dp = [T,F,F,F,F,F,F,F,F,F,F,F]` (index 0..11).

**Process `num = 1`** (s from 11 down to 1): only `dp[1] |= dp[0]` fires.

```
dp: [T,T,F,F,F,F,F,F,F,F,F,F]     reachable sums: {0,1}
```

**Process `num = 5`** (s from 11 down to 5): `dp[6]|=dp[1]`, `dp[5]|=dp[0]`.

```
dp: [T,T,F,F,F,T,T,F,F,F,F,F]     reachable: {0,1,5,6}
```

**Process `num = 11`** (s from 11 down to 11): `dp[11] |= dp[0]` → **true**.

```
dp: [T,T,F,F,F,T,T,F,F,F,F,T]     reachable: {0,1,5,6,11}
```

`dp[11]` is already true — we can stop early. (For the trace's sake, the last item `5` would additionally set `dp[10],dp[11]` from `dp[5],dp[6]`.)

**Result:** `dp[11] = true` → the array **can** be partitioned (e.g. `{11}` and `{1,5,5}`). ✔

Notice the **descending** scan: when processing `num=5`, `dp[5] |= dp[0]` reads `dp[0]` which belongs to the *previous* item's table, so the single `5` isn't accidentally reused to also make `dp[10]` in the same pass.

---

## Complexity

For n items and capacity (target sum) C:

| Variant | Time | Space (rolling) |
|---------|------|-----------------|
| 0/1 knapsack / subset-sum / count | O(n · C) | O(C) |
| Unbounded (coin change) | O(n · C) | O(C) |
| 2-D-capacity (#474) | O(L · m · n) | O(m · n) |
| k-partition (#473, backtracking) | O(kⁿ) worst, heavy pruning in practice | O(n) |

- These bounds are **pseudo-polynomial**: linear in the *numeric value* C, not in its bit-length. That is why knapsack is NP-hard in general yet tractable when C is small — the catch is that C = 10⁹ blows the table up.
- Space drops from the naive O(n · C) table to O(C) via the rolling array; the full table is only needed if you must **reconstruct** the chosen items.

---

## Common pitfalls

1. **Wrong loop direction.** The single most common bug: iterating capacity **ascending** for a 0/1 problem (silently reuses items) or **descending** for an unbounded one (forbids reuse). Descending = once, ascending = many. Everything else about the two is identical.

2. **Coin Change II loop order.** Coins **outer**, amount **inner** counts *combinations*. Reversing the loops counts *permutations* — a different answer. If your "#ways to make change" is too large, you probably counted ordered sequences.

3. **Skipping the `dp[0]` base case.** `dp[0] = true` (subset-sum) / `dp[0] = 1` (counting) / `dp[0] = 0` (min-count) seeds "the empty selection". Forget it and every count comes out 0.

4. **Not pre-checking feasibility of the target.** For partition/#494: if `total` is odd (can't halve) or `total + target` is odd or `|target| > total`, the target isn't an integer subset sum — bail out with `false`/`0` *before* building `dp`, else you index negatively or return garbage.

5. **Off-by-one / negative index in the inner loop.** Guard `c >= weight` (equivalently loop `c := C; c >= weight; c--`). Reading `dp[c-weight]` with `c < weight` is an out-of-range panic.

6. **Integer overflow in counting variants.** Subset counts (#494, #518) can exceed 32-bit range for large inputs; use `int` (64-bit on most platforms) or the problem's specified modulus.

7. **Confusing k-partition with plain subset-sum.** "Split into k equal groups" (#473) is **not** one subset-sum — a value fitting one bucket's total doesn't guarantee all k buckets pack simultaneously. The robust solution is backtracking (assign each item to a bucket, prune) — sort descending and skip duplicate/empty buckets to make it fast.

8. **Using knapsack DP when C is huge.** Pseudo-polynomial means "polynomial in C's magnitude". If the target sum can be 10¹²+, the array won't fit — the problem wants number theory or greedy, not a DP table.

9. **Trying to reconstruct items from a rolling array.** The 1-D array loses per-item history. If you must output *which* items were chosen, keep the full 2-D `dp[i][c]` table (or a parent/choice trace) and walk it backwards.

---

## Problems in this repo that use Knapsack DP

- [0416 — Partition Equal Subset Sum](/0416_partition_equal_subset_sum/README.md) — **subset-sum (boolean)**: total must be even; ask whether a subset sums to `total/2`. Textbook 1-D DP with the descending scan.
- [0494 — Target Sum](/0494_target_sum/README.md) — **subset-sum count** via the algebraic reduction `P = (total + target)/2`; `dp[s] += dp[s-num]` counts sign assignments.
- [0322 — Coin Change](/0322_coin_change/README.md) — **unbounded knapsack, minimisation**: fewest coins to make an amount; ascending inner loop for reuse; `-1` if unreachable.
- [0474 — Ones and Zeroes](/0474_ones_and_zeroes/README.md) — **2-D-capacity 0/1 knapsack**: each string is an item weighing `(zeros, ones)` with value 1; maximise count under two independent budgets `m` and `n`, both scanned descending.
- [0473 — Matchsticks to Square](/0473_matchsticks_to_square/README.md) — **4-way partition** (k-partition of the subset-sum family): target side = `total/4`; solved by pruned backtracking (assign each stick to one of four sides), the k-group generalisation of #416.

### Related classics to know (not yet in repo)

- LeetCode #518 — Coin Change II (**unbounded #ways** / combinations count — coins outer, amount inner; the counting cousin of #322)
- LeetCode #377 — Combination Sum IV (unbounded **permutations** count — amount outer, coin inner: the loop-order contrast to #518)
- LeetCode #698 — Partition to K Equal Sum Subsets (the general k-partition #473 specialises)
- LeetCode #1049 — Last Stone Weight II (subset-sum reframed as "minimise |sum(A) − sum(B)|")
- LeetCode #279 — Perfect Squares (unbounded knapsack: fewest squares summing to n)
- LeetCode #139 — Word Break (unbounded-knapsack-shaped reachability over a dictionary)
