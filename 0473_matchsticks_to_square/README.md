# 0473 — Matchsticks to Square

> LeetCode #473 · Difficulty: Medium
> **Categories:** Array, Backtracking, Bitmask, Dynamic Programming, Bit Manipulation

---

## Problem Statement

You are given an integer array `matchsticks` where `matchsticks[i]` is the length of the `ith` matchstick. You want to use **all the matchsticks** to make one square. You **should not break** any stick, but you can link them up, and each matchstick must be used **exactly one time**.

Return `true` if you can make this square and `false` otherwise.

**Example 1:**

```
Input: matchsticks = [1,1,2,2,2]
Output: true
Explanation: You can form a square with length 2, one side of the square came two sticks with length 1.
```

**Example 2:**

```
Input: matchsticks = [3,3,3,3,4]
Output: false
Explanation: You cannot find a way to form a square with all the matchsticks.
```

**Constraints:**

- `1 <= matchsticks.length <= 15`
- `1 <= matchsticks[i] <= 10^8`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |
| Apple      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Backtracking** — the core solution is a depth-first assignment of each stick to one of four sides, undoing a choice when it cannot complete a square → see [`/dsa/backtracking.md`](/dsa/backtracking.md)
- **Bitmask Dynamic Programming** — with `n ≤ 15`, a subset of used sticks fits in an `int`, so `dp[mask]` over `2^n` states gives a polynomial-in-`2^n` alternative to raw search → see [`/dsa/bit_manipulation.md`](/dsa/bit_manipulation.md)
- **Partition / Subset-Sum shape** — "split the multiset into 4 groups of equal sum" is a k-partition problem; the target side `total/4` is the subset-sum goal repeated four times → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Backtracking (fill 4 sides) | O(4^n) | O(n) | Clean baseline; may TLE on adversarial n=15 without pruning |
| 2 | Backtracking + pruning (sort desc, skip dup buckets) | O(4^n) worst, fast in practice | O(n) | The standard interview answer; passes easily |
| 3 | Bitmask DP | O(2^n · n) | O(2^n) | Deterministic bound; shines when you want no reliance on pruning |

---

## Approach 1 — Backtracking (Fill Four Sides)

### Intuition

A square is four sides of equal length `total/4`. Think of four buckets, each with capacity `side`. Drop each matchstick into some bucket without overflowing it; if every stick lands and the buckets fill exactly, a square exists. Because the buckets' capacities sum to the total, "all sticks placed" automatically means "all four buckets full" — no separate check needed.

### Algorithm

1. Sum all sticks. If the sum is `0` or not divisible by 4, return `false`.
2. Let `side = total/4` and keep `sides[4]` running lengths, all `0`.
3. `dfs(i)`: if `i == n`, return `true` (all placed ⇒ all sides exactly `side`).
4. For each bucket `b`: if `sides[b] + matchsticks[i] <= side`, add the stick, recurse on `i+1`, and undo on failure.
5. Return whether any branch reached the end.

### Complexity

- **Time:** O(4^n) — each of the `n` sticks may try all 4 buckets; without pruning the tree is exponential.
- **Space:** O(n) — recursion depth `n`, plus the fixed 4-length buckets array.

### Code

```go
func backtracking(matchsticks []int) bool {
	total := 0
	for _, m := range matchsticks {
		total += m // accumulate the perimeter
	}
	// A square needs a positive perimeter divisible by 4.
	if total == 0 || total%4 != 0 {
		return false
	}
	side := total / 4 // target length of every side
	sides := [4]int{} // running length of each of the four sides

	var dfs func(i int) bool
	dfs = func(i int) bool {
		if i == len(matchsticks) {
			// All sticks used; sides are valid iff each already reached `side`.
			// Because we never let a bucket exceed `side`, and the totals add up
			// to 4*side, reaching the end already guarantees all four equal side.
			return true
		}
		for b := 0; b < 4; b++ {
			// Only place stick i into bucket b if it still fits.
			if sides[b]+matchsticks[i] <= side {
				sides[b] += matchsticks[i] // place stick i on side b
				if dfs(i + 1) {
					return true // this placement led to a full square
				}
				sides[b] -= matchsticks[i] // undo — try the next bucket
			}
		}
		return false // stick i fit nowhere on a valid square
	}
	return dfs(0)
}
```

### Dry Run

Example 1: `matchsticks = [1,1,2,2,2]`, `total = 8`, `side = 2`.

| Step | i | stick | Action | sides after |
|------|---|-------|--------|-------------|
| 1 | 0 | 1 | place in bucket 0 | `[1,0,0,0]` |
| 2 | 1 | 1 | bucket 0 → `1+1=2 ≤ 2`, place | `[2,0,0,0]` |
| 3 | 2 | 2 | bucket 0 full (`2+2>2`); bucket 1 fits | `[2,2,0,0]` |
| 4 | 3 | 2 | bucket 1 full; bucket 2 fits | `[2,2,2,0]` |
| 5 | 4 | 2 | bucket 2 full; bucket 3 fits | `[2,2,2,2]` |
| 6 | 5 | — | `i == n` → return `true` | — |

Result: `true` ✔

---

## Approach 2 — Backtracking + Pruning (Sort Desc, Skip Duplicate Buckets)

### Intuition

The naive search wastes time in two ways, and two prunes fix both:

1. **Sort descending.** Large sticks are the most constraining, so place them first — a stick longer than `side` is rejected at depth 0, and big pieces slot into buckets before the small ones fill the gaps. This turns many exponential branches into instant failures.
2. **Skip symmetric buckets.** If two buckets currently hold the same length, dropping the stick into either is *the same decision*; trying both just re-explores an identical subtree. So within one `dfs` call, skip a bucket whose running length equals one already tried (in particular, all empty buckets are interchangeable).

### Algorithm

1. Divisibility check; `side = total/4`.
2. Sort `matchsticks` descending. If `matchsticks[0] > side`, return `false`.
3. `dfs(i)`: for each bucket `b`, skip if some earlier bucket `k < b` has `sides[k] == sides[b]` (duplicate state). Otherwise, if the stick fits, place, recurse, undo.
4. Succeed when `i == n`.

### Complexity

- **Time:** O(4^n) in the worst case, but the two prunes make it effectively linear-ish for `n ≤ 15` — it is the accepted solution.
- **Space:** O(n) — recursion stack.

### Code

```go
func backtrackingPruned(matchsticks []int) bool {
	total := 0
	for _, m := range matchsticks {
		total += m
	}
	if total == 0 || total%4 != 0 {
		return false
	}
	side := total / 4
	// Sort descending so large sticks are placed first (fail fast).
	sort.Sort(sort.Reverse(sort.IntSlice(matchsticks)))
	// Largest stick cannot exceed one side.
	if matchsticks[0] > side {
		return false
	}
	sides := [4]int{}

	var dfs func(i int) bool
	dfs = func(i int) bool {
		if i == len(matchsticks) {
			return true // all placed, all buckets full
		}
		for b := 0; b < 4; b++ {
			// Symmetry prune: buckets with the same running length are
			// interchangeable — only try the first such bucket.
			dup := false
			for k := 0; k < b; k++ {
				if sides[k] == sides[b] {
					dup = true
					break
				}
			}
			if dup {
				continue
			}
			if sides[b]+matchsticks[i] <= side {
				sides[b] += matchsticks[i]
				if dfs(i + 1) {
					return true
				}
				sides[b] -= matchsticks[i]
			}
		}
		return false
	}
	return dfs(0)
}
```

### Dry Run

Example 1: `matchsticks = [1,1,2,2,2]` → sorted desc `[2,2,2,1,1]`, `side = 2`.

| Step | i | stick | Bucket chosen (dup-skips) | sides after |
|------|---|-------|---------------------------|-------------|
| 1 | 0 | 2 | buckets 1,2,3 all equal bucket 0 (all 0) → only bucket 0 tried; fits (`2`) | `[2,0,0,0]` |
| 2 | 1 | 2 | bucket 0 full; buckets 2,3 dup of bucket 1 → only bucket 1; fits | `[2,2,0,0]` |
| 3 | 2 | 2 | buckets 0,1 full; bucket 3 dup of bucket 2 → only bucket 2; fits | `[2,2,2,0]` |
| 4 | 3 | 1 | only bucket 3 has room; fits (`1`) | `[2,2,2,1]` |
| 5 | 4 | 1 | bucket 3 → `1+1=2 ≤ 2`; fits | `[2,2,2,2]` |
| 6 | 5 | — | `i == n` → `true` | — |

Result: `true` ✔ — the descending sort places the three 2s first, one per side, and the two 1s finish the fourth side.

---

## Approach 3 — Bitmask DP (Optimal)

### Intuition

With `n ≤ 15`, "which sticks are used" fits in a 15-bit mask (≤ 32768 values). Build sides one stick at a time and track only how full the **current** side is. Let `dp[mask]` = the length occupied on the in-progress side after using exactly the sticks in `mask`, defined only for masks reachable as (some completed sides) + (this partial side). From a reachable `mask`, add any unused stick `j` that fits in the room `side - dp[mask]`; the new occupancy is `(dp[mask] + len) mod side`, where the `mod` snaps a just-closed side back to `0` and starts the next one. When `mask` is full and `dp == 0`, all four sides closed exactly.

### Algorithm

1. Divisibility check; `side = total/4`. If any stick `> side`, return `false`.
2. `dp` of size `2^n`, all `-1` (unreachable) except `dp[0] = 0`.
3. For each `mask` with `dp[mask] >= 0`, for each stick `j` not in `mask`:
   if `dp[mask] + matchsticks[j] <= side`, set `dp[mask | (1<<j)] = (dp[mask] + matchsticks[j]) % side` (if not already set).
4. Return `dp[fullMask] == 0`.

### Complexity

- **Time:** O(2^n · n) — every one of the `2^n` masks scans `n` sticks. For `n = 15` that is ≈ 490K operations, trivially fast and independent of stick values.
- **Space:** O(2^n) — the `dp` table.

### Code

```go
func bitmaskDP(matchsticks []int) bool {
	n := len(matchsticks)
	total := 0
	for _, m := range matchsticks {
		total += m
	}
	if total == 0 || total%4 != 0 {
		return false
	}
	side := total / 4
	// Any single stick longer than a side is an immediate no.
	for _, m := range matchsticks {
		if m > side {
			return false
		}
	}
	full := (1 << n) - 1 // mask with all sticks used
	// dp[mask] = length occupied on the current (in-progress) side, or -1 if the
	// subset `mask` cannot be arranged into complete-sides + this partial side.
	dp := make([]int, 1<<n)
	for i := range dp {
		dp[i] = -1
	}
	dp[0] = 0 // empty set: current side is empty, reachable
	for mask := 0; mask <= full; mask++ {
		if dp[mask] < 0 {
			continue // this configuration is not reachable
		}
		for j := 0; j < n; j++ {
			if mask&(1<<j) != 0 {
				continue // stick j already used
			}
			// Stick j must fit in the room left on the current side.
			if dp[mask]+matchsticks[j] <= side {
				// Adding it; wrap to 0 when the side closes exactly.
				used := (dp[mask] + matchsticks[j]) % side
				next := mask | (1 << j)
				// Prefer any reachable state; all lead to the same closure logic.
				if dp[next] == -1 {
					dp[next] = used
				}
			}
		}
	}
	// All sticks used AND the current side ended exactly closed ⇒ square.
	return dp[full] == 0
}
```

### Dry Run

Example 1: `matchsticks = [1,1,2,2,2]` (indices 0..4), `side = 2`. Bits are `b4 b3 b2 b1 b0`. Showing one path that reaches `dp[full] == 0`:

| mask (used) | dp[mask] | add stick j (len) | fits? room = 2 − dp | new mask | new dp = (dp+len) % 2 |
|-------------|----------|-------------------|---------------------|----------|-----------------------|
| `00000` | 0 | j0 (1) | 1 ≤ 2 ✔ | `00001` | `1` |
| `00001` | 1 | j1 (1) | 1 ≤ 1 ✔ | `00011` | `(1+1)%2 = 0` |
| `00011` | 0 | j2 (2) | 2 ≤ 2 ✔ | `00111` | `(0+2)%2 = 0` |
| `00111` | 0 | j3 (2) | 2 ≤ 2 ✔ | `01111` | `0` |
| `01111` | 0 | j4 (2) | 2 ≤ 2 ✔ | `11111` | `0` |

`dp[11111] = dp[full] = 0` → all sticks used, current side closed exactly.

Result: `true` ✔

---

## Key Takeaways

- **"Form k equal groups" = k-way partition.** The target group sum is `total/k`; bail immediately if `total % k != 0` or any element exceeds the target.
- **Backtracking prunes that matter:** sort *descending* (place the most constraining items first) and *skip interchangeable buckets/positions* (equal running sums are the same decision). These two alone rescue an `O(4^n)` search on `n = 15`.
- **`n ≤ 20`-ish ⇒ consider a bitmask.** A subset over `2^n` states with an `O(n)` transition is often cleaner and gives a hard worst-case bound, unlike search that leans on pruning.
- **`(x) % side` is a neat "close-and-reset" trick** when you fill fixed-capacity bins one after another.

---

## Related Problems

- LeetCode #698 — Partition to K Equal Sum Subsets (identical technique, general `k`)
- LeetCode #416 — Partition Equal Subset Sum (`k = 2`, pure subset-sum DP)
- LeetCode #1723 — Find Minimum Time to Finish All Jobs (bucket-filling backtracking / bitmask)
- LeetCode #93 — Restore IP Addresses (fixed-count partition via backtracking)
