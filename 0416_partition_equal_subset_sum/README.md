# 0416 — Partition Equal Subset Sum

> LeetCode #416 · Difficulty: Medium
> **Categories:** Array, Dynamic Programming, 0/1 Knapsack

---

## Problem Statement

Given an integer array `nums`, return `true` if you can partition the array into two subsets such that the sum of the elements in both subsets is equal or `false` otherwise.

**Example 1:**

```
Input: nums = [1,5,11,5]
Output: true
Explanation: The array can be partitioned as [1, 5, 5] and [11].
```

**Example 2:**

```
Input: nums = [1,2,3,5]
Output: false
Explanation: The array cannot be partitioned into equal sum subsets.
```

**Constraints:**

- `1 <= nums.length <= 200`
- `1 <= nums[i] <= 100`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Adobe      | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Dynamic Programming (2D → 1D)** — this is a 0/1-knapsack in disguise: "can a subset reach sum `total/2`?" is a two-dimensional reachability table over (items, sum) that compresses to one row → see [`/dsa/dynamic_programming_2d.md`](/dsa/dynamic_programming_2d.md)
- **Array** — the input is a flat array of positive integers whose total and half-total drive the entire solution → see [`/dsa/arrays.md`](/dsa/arrays.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (recursive subset sum) | O(2ⁿ) | O(n) | Understanding only; n ≤ 200 with values ≤ 100 → TLE |
| 2 | Top-Down DP (memoised) | O(n·target) | O(n·target) | Natural memoisation of the recursion; easiest correctness argument |
| 3 | Bottom-Up 2D DP | O(n·target) | O(n·target) | Textbook knapsack table; clear to explain in interview |
| 4 | Bottom-Up 1D DP (Optimal) | O(n·target) | O(target) | Best space; the answer to write once you know the trick |

`target = sum/2`, so `target ≤ 100·200/2 = 10⁴`.

---

## Approach 1 — Brute Force (Recursive Subset Sum)

### Intuition

Two subsets with equal sums means each subset sums to exactly half the total. So the problem is really the **subset-sum decision problem**: is there a subset of `nums` summing to `target = total/2`? If the total is odd, the split is impossible outright. Otherwise, walk the array making a binary choice at every element — take it toward the target or skip it — and succeed if any path drives the remaining target to exactly 0.

### Algorithm

1. Compute `sum`; if `sum` is odd, return `false`.
2. Let `target = sum / 2`.
3. Recurse `canReach(index, remaining)`:
   - if `remaining == 0` → `true` (found a subset).
   - if `index` past the end or `remaining < 0` → `false` (dead end).
   - else return `canReach(index+1, remaining - nums[index])` **OR** `canReach(index+1, remaining)`.
4. Answer is `canReach(0, target)`.

### Complexity

- **Time:** O(2ⁿ) — each of the n elements doubles the number of paths.
- **Space:** O(n) — recursion depth is at most n frames.

### Code

```go
func bruteForce(nums []int) bool {
	sum := 0
	for _, v := range nums { // total of all elements
		sum += v
	}
	if sum%2 != 0 { // an odd total can never split into two equal integer halves
		return false
	}
	target := sum / 2 // each side must reach exactly this

	var canReach func(index, remaining int) bool
	canReach = func(index, remaining int) bool {
		if remaining == 0 {
			return true // exact subset found
		}
		if index >= len(nums) || remaining < 0 {
			return false // ran out of items or overshot the target
		}
		// Branch: use nums[index] toward the target, OR leave it out.
		return canReach(index+1, remaining-nums[index]) ||
			canReach(index+1, remaining)
	}
	return canReach(0, target)
}
```

### Dry Run

Example 1: `nums = [1,5,11,5]`, `sum = 22`, `target = 11`. Trace the first branch that succeeds (take-first bias explores subtracting each element).

| Call | index | remaining | Action taken | Result |
|------|-------|-----------|--------------|--------|
| 1 | 0 | 11 | take nums[0]=1 → remaining 10 | (recurse) |
| 2 | 1 | 10 | take nums[1]=5 → remaining 5 | (recurse) |
| 3 | 2 | 5 | take nums[2]=11 → remaining −6 < 0 | false, backtrack |
| 4 | 2 | 5 | skip nums[2]=11 → remaining 5 | (recurse) |
| 5 | 3 | 5 | take nums[3]=5 → remaining 0 | **true** |

Subset `{1, 5, 5}` sums to 11. Result: `true` ✔

---

## Approach 2 — Top-Down DP (Memoised Subset Sum)

### Intuition

The recursion tree revisits identical `(index, remaining)` states through different take/skip orderings — e.g. skipping then taking vs. taking then skipping can land on the same `(index, remaining)`. The outcome for a state is fixed, so cache it. Only `n × (target+1)` distinct states exist; solving each once collapses the exponential blow-up to a polynomial.

### Algorithm

1. Same reduction to `target = sum/2` (odd sum → `false`).
2. Keep `memo[index][remaining]` with three states: unvisited / known-false / known-true.
3. In `canReach`, before recursing check the cache; after computing, store the result.

### Complexity

- **Time:** O(n·target) — each `(index, remaining)` pair evaluated at most once.
- **Space:** O(n·target) memo table + O(n) recursion stack.

### Code

```go
func dpTopDown(nums []int) bool {
	sum := 0
	for _, v := range nums {
		sum += v
	}
	if sum%2 != 0 {
		return false
	}
	target := sum / 2

	// memo[i][r]: 0 = unvisited, 1 = false, 2 = true (avoids a separate seen map).
	memo := make([][]int8, len(nums))
	for i := range memo {
		memo[i] = make([]int8, target+1)
	}

	var canReach func(index, remaining int) bool
	canReach = func(index, remaining int) bool {
		if remaining == 0 {
			return true
		}
		if index >= len(nums) || remaining < 0 {
			return false
		}
		if memo[index][remaining] != 0 { // seen this exact state before
			return memo[index][remaining] == 2
		}
		res := canReach(index+1, remaining-nums[index]) || // take
			canReach(index+1, remaining) // skip
		if res {
			memo[index][remaining] = 2 // cache true
		} else {
			memo[index][remaining] = 1 // cache false
		}
		return res
	}
	return canReach(0, target)
}
```

### Dry Run

Example 1: `nums = [1,5,11,5]`, `target = 11`. States are cached as they resolve.

| Step | State (index, remaining) | Cache before | Computed | Cache after |
|------|--------------------------|--------------|----------|-------------|
| 1 | (0, 11) | unvisited | needs children | pending |
| 2 | (1, 10) | unvisited | needs children | pending |
| 3 | (2, 5) via take-11 | unvisited | −6 branch false; skip branch → (3,5) | pending |
| 4 | (3, 5) | unvisited | take-5 → (4,0)=true | store **true** |
| 5 | (2, 5) resolves | pending | skip branch true | store **true** |
| 6 | (1, 10), (0, 11) resolve | pending | propagate true up | store **true** |

Result: `true` ✔ — later probes of `(2,5)` would hit the cached `true` instead of recursing.

---

## Approach 3 — Bottom-Up 2D DP

### Intuition

Turn the recursion inside out into an explicit table. `dp[i][s]` answers "using the first `i` numbers, can I make sum `s`?" You can make `s` with the first `i` items iff either the first `i-1` already make `s` (skip item `i`) or they make `s - nums[i-1]` (take item `i`). Sum `0` is always reachable via the empty subset. Read the final answer off `dp[n][target]`.

### Algorithm

1. Reduce to `target = sum/2` (odd → `false`).
2. Initialise `dp[i][0] = true` for every `i`.
3. For each item `i` (1..n) and each sum `s` (1..target):
   `dp[i][s] = dp[i-1][s] || (s >= nums[i-1] && dp[i-1][s-nums[i-1]])`.
4. Return `dp[n][target]`.

### Complexity

- **Time:** O(n·target) — one pass over an (n+1) × (target+1) grid.
- **Space:** O(n·target) — the full table.

### Code

```go
func dp2D(nums []int) bool {
	sum := 0
	for _, v := range nums {
		sum += v
	}
	if sum%2 != 0 {
		return false
	}
	target := sum / 2
	n := len(nums)

	// dp[i][s] = can the first i numbers form subset-sum s.
	dp := make([][]bool, n+1)
	for i := range dp {
		dp[i] = make([]bool, target+1)
		dp[i][0] = true // sum 0 always achievable via the empty subset
	}

	for i := 1; i <= n; i++ {
		v := nums[i-1] // the i-th number (0-indexed)
		for s := 1; s <= target; s++ {
			dp[i][s] = dp[i-1][s] // case 1: don't use v
			if s >= v {
				// case 2: use v, provided the remainder s-v was reachable before
				dp[i][s] = dp[i][s] || dp[i-1][s-v]
			}
		}
	}
	return dp[n][target]
}
```

### Dry Run

Example 1: `nums = [1,5,11,5]`, `target = 11`. Showing which sums each prefix can reach (✓ = reachable). Column 0 is always ✓.

| After item | reachable sums (0..11) |
|------------|------------------------|
| none | 0 |
| +1 | 0, 1 |
| +5 | 0, 1, 5, 6 |
| +11 | 0, 1, 5, 6, 11 |
| +5 | 0, 1, 5, 6, 10, 11 |

`dp[4][11] = ✓`. Result: `true` ✔ (11 first appears after adding the `11`; the final `5` adds 10 = 5+5.)

---

## Approach 4 — Bottom-Up 1D DP (Optimal)

### Intuition

Row `i` of the 2D table depends only on row `i-1`, so a single boolean array `dp[s]` suffices — overwrite it in place per item. The subtlety is the **iteration direction**: each item may be used at most once (0/1 knapsack), so process sums from `target` **down** to `v`. Descending guarantees that when we read `dp[s-v]` it still reflects the state *before* this item was folded in; iterating upward would allow the same item to contribute to `s-v` and then again to `s`, i.e. unbounded knapsack.

### Algorithm

1. Reduce to `target = sum/2` (odd → `false`).
2. `dp[0] = true`.
3. For each value `v`: for `s` from `target` down to `v`: `dp[s] = dp[s] || dp[s-v]`.
   Short-circuit to `true` as soon as `dp[target]` is set.
4. Return `dp[target]`.

### Complexity

- **Time:** O(n·target) — same state count, single row.
- **Space:** O(target) — one boolean array.

### Code

```go
func dp1D(nums []int) bool {
	sum := 0
	for _, v := range nums {
		sum += v
	}
	if sum%2 != 0 {
		return false
	}
	target := sum / 2

	dp := make([]bool, target+1) // dp[s] = some processed subset sums to s
	dp[0] = true                 // empty subset

	for _, v := range nums {
		// Descend so each v folds in at most once (0/1-knapsack requirement).
		for s := target; s >= v; s-- {
			if dp[s-v] { // s-v was reachable without v → s reachable with v
				dp[s] = true
			}
		}
		if dp[target] {
			return true // target already reachable — no need to process more items
		}
	}
	return dp[target]
}
```

### Dry Run

Example 1: `nums = [1,5,11,5]`, `target = 11`. Each row shows the reachable sums after folding one value in (scanning `s` high→low).

| Fold value | dp reachable sums (0..11) after this value | dp[11]? |
|------------|--------------------------------------------|---------|
| start | 0 | no |
| 1 | 0, 1 | no |
| 5 | 0, 1, 5, 6 | no |
| 11 | 0, 1, 5, 6, 11 | **yes** → return true |

Result: `true` ✔ — reached `dp[11]` after the third value, so the final `5` is never processed.

---

## Key Takeaways

- **"Partition into two equal-sum subsets" ≡ "subset-sum to total/2".** Halving the total is the reduction that unlocks knapsack DP. Always parity-check the total first: odd total → instant `false`.
- **0/1 knapsack in 1D scans the capacity dimension backwards.** High→low reuse of `dp[s-v]` enforces "each item at most once"; low→high would model unbounded quantities. This direction rule is the single most error-prone detail in knapsack problems.
- **Three DP forms, one recurrence.** Memoised recursion, 2D table, and 1D rolling array all encode `take v` OR `skip v`; know how to derive each from the last and which to reach for (1D for space, 2D/memo for clarity).
- **Bounded values bound the table.** Because `nums[i] ≤ 100` and `n ≤ 200`, `target ≤ 10⁴`, keeping the pseudo-polynomial `O(n·target)` firmly in range.

---

## Related Problems

- LeetCode #494 — Target Sum (assign ± signs → subset-sum reduction)
- LeetCode #698 — Partition to K Equal Sum Subsets (generalises to k parts)
- LeetCode #1049 — Last Stone Weight II (minimise |sum difference| = subset-sum)
- LeetCode #322 — Coin Change (unbounded knapsack: forward capacity scan)
- LeetCode #474 — Ones and Zeroes (2D-capacity knapsack)
