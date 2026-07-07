# 0403 — Frog Jump

> LeetCode #403 · Difficulty: Hard
> **Categories:** Dynamic Programming, Hash Table, Array

---

## Problem Statement

A frog is crossing a river. The river is divided into some number of units, and at each unit, there may or may not exist a stone. The frog can jump on a stone, but it must not jump into the water.

Given a list of `stones`' positions (in units) in sorted **ascending order**, determine if the frog can cross the river by landing on the last stone. Initially, the frog is on the first stone and assumes the first jump must be `1` unit.

If the frog's last jump was `k` units, its next jump must be either `k - 1`, `k`, or `k + 1` units. The frog can only jump in the forward direction.

**Example 1:**

```
Input: stones = [0,1,3,5,6,8,12,17]
Output: true
Explanation: The frog can jump to the last stone by jumping 1 unit to the 2nd stone, then 2 units to the 3rd stone, then 2 units to the 4th stone, then 3 units to the 6th stone, 4 units to the 7th stone, and 5 units to the 8th stone.
```

**Example 2:**

```
Input: stones = [0,1,2,3,4,8,9,11]
Output: false
Explanation: There is no way to jump to the last stone as the gap between the 5th and 6th stone is too large.
```

**Constraints:**

- `2 <= stones.length <= 2000`
- `0 <= stones[i] <= 2^31 - 1`
- `stones[0] == 0`
- `stones` is sorted in a strictly increasing order.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★☆☆☆ Low        | 2022          |
| Apple      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Dynamic Programming (2D state)** — the state is `(stone, last jump k)`; the same stone can be reached with different last-jump sizes, giving an `O(n²)` state space over stones × jump-sizes → see [`/dsa/dynamic_programming_2d.md`](/dsa/dynamic_programming_2d.md)
- **Hash Map** — positions are sparse and up to `2³¹−1`, so we map position → stone index (and store reachable jump sets in maps) instead of indexing a giant array → see [`/dsa/hash_map.md`](/dsa/hash_map.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Plain Recursion | O(3ⁿ) | O(n) | Explains the state and transitions; exponential — only for understanding |
| 2 | Top-Down DP (memoized) | O(n²) | O(n²) | Same recursion made efficient; natural if you started from brute force |
| 3 | Bottom-Up DP (reachable jump sets) | O(n²) | O(n²) | Forward-propagation formulation; the standard optimal answer |

---

## Approach 1 — Plain Recursion

### Intuition

The frog's situation is fully described by `(current position, last jump size k)`. From there it may jump `k−1`, `k`, or `k+1` units, but only if the jump is at least `1` unit and lands **exactly** on a stone. Recurse on every legal jump; the frog succeeds the moment it stands on the last stone. Starting with a phantom last jump of `k = 0` makes the first real jump `k+1 = 1`, matching the rule that the opening jump is 1 unit. With no memoization, the same `(position, k)` is re-explored along many paths, so the work is exponential.

### Algorithm

1. Build a set of stone positions and note the target = `stones[last]`.
2. `dfs(position, k)`:
   1. If `position == target`, return `true`.
   2. For `step` in `{k−1, k, k+1}` with `step ≥ 1`: if `position + step` is a stone, recurse; return `true` if any succeeds.
3. Start `dfs(stones[0], 0)`.

### Complexity

- **Time:** O(3ⁿ) — up to three branches per stone with no reuse of subresults.
- **Space:** O(n) recursion depth plus O(n) for the position set.

### Code

```go
func bruteForceRecursion(stones []int) bool {
	target := stones[len(stones)-1]
	stoneSet := make(map[int]bool, len(stones)) // fast "is there a stone here?"
	for _, s := range stones {
		stoneSet[s] = true
	}

	var dfs func(position, k int) bool
	dfs = func(position, k int) bool {
		if position == target {
			return true // reached the far bank
		}
		// Try the three permitted next jump sizes.
		for _, step := range []int{k - 1, k, k + 1} {
			if step <= 0 {
				continue // must jump forward by at least 1 unit
			}
			next := position + step
			if stoneSet[next] && dfs(next, step) {
				return true // some continuation from `next` succeeds
			}
		}
		return false
	}
	// The problem fixes the first jump at 1: starting with k=0 makes k+1 == 1.
	return dfs(stones[0], 0)
}
```

### Dry Run

Example 1: `stones = [0,1,3,5,6,8,12,17]`, target = 17. Following the successful path only (the recursion explores others and backtracks):

| position | k (arrival) | steps tried {k−1,k,k+1} | chosen step lands on stone | next position |
|----------|-------------|-------------------------|----------------------------|---------------|
| 0 | 0 | {−1,0,1} → only 1 valid | 1 → stone 1 | 1 |
| 1 | 1 | {0,1,2} | 2 → stone 3 | 3 |
| 3 | 2 | {1,2,3} | 2 → stone 5 | 5 |
| 5 | 2 | {1,2,3} | 3 → stone 8 | 8 |
| 8 | 3 | {2,3,4} | 4 → stone 12 | 12 |
| 12 | 4 | {3,4,5} | 5 → stone 17 | 17 = target |

Reached 17 → `true` ✔

---

## Approach 2 — Top-Down DP

### Intuition

The exponential cost is pure repetition: many jump paths land on the same stone with the same last-jump `k`, and the answer from `(stone, k)` never changes. Cache it. There are at most `n` stones and any jump size is bounded by `n`, so the memo covers `O(n²)` distinct states, each solved once.

### Algorithm

1. Map each position to its stone index; create a per-index memo keyed by `k`.
2. `dfs(index, k)`:
   1. If `index == n−1`, return `true`.
   2. If `(index, k)` is memoized, return it.
   3. For `step` in `{k−1, k, k+1}` with `step ≥ 1`: if `stones[index] + step` is a stone at index `j`, recurse `dfs(j, step)`. Record and return whether any branch succeeds.
3. Start `dfs(0, 0)`.

### Complexity

- **Time:** O(n²) — `O(n²)` states, each doing constant work across 3 transitions.
- **Space:** O(n²) memo plus O(n) recursion depth.

### Code

```go
func dpTopDown(stones []int) bool {
	n := len(stones)
	indexOf := make(map[int]int, n) // stone position -> its index
	for i, s := range stones {
		indexOf[s] = i
	}

	// memo[index][k] caches whether (stone index, last jump k) can finish.
	// k ranges 0..n (a jump can't exceed n stones), so width n+1 is safe.
	memo := make([]map[int]bool, n)
	for i := range memo {
		memo[i] = make(map[int]bool)
	}

	var dfs func(index, k int) bool
	dfs = func(index, k int) bool {
		if index == n-1 {
			return true // standing on the last stone
		}
		if v, seen := memo[index][k]; seen {
			return v // already solved this exact state
		}
		res := false
		for _, step := range []int{k - 1, k, k + 1} {
			if step <= 0 {
				continue
			}
			if j, ok := indexOf[stones[index]+step]; ok {
				if dfs(j, step) {
					res = true
					break // one successful continuation is enough
				}
			}
		}
		memo[index][k] = res // remember for next time
		return res
	}
	return dfs(0, 0)
}
```

### Dry Run

Example 1: `stones = [0,1,3,5,6,8,12,17]`. Same successful chain as Approach 1, but each `(index, k)` is now cached the first time it is computed.

| call | index (position) | k | resolves via | memo written |
|------|------------------|---|--------------|--------------|
| dfs(0,0) | 0 (0) | 0 | step 1 → dfs(1,1) | memo[0][0]=true |
| dfs(1,1) | 1 (1) | 1 | step 2 → dfs(2,2) | memo[1][1]=true |
| dfs(2,2) | 2 (3) | 2 | step 2 → dfs(3,2) | memo[2][2]=true |
| dfs(3,2) | 3 (5) | 2 | step 3 → dfs(5,3) | memo[3][2]=true |
| dfs(5,3) | 5 (8) | 3 | step 4 → dfs(6,4) | memo[5][3]=true |
| dfs(6,4) | 6 (12) | 4 | step 5 → dfs(7,5) | memo[6][4]=true |
| dfs(7,5) | 7 (17) | 5 | index == n−1 | returns true |

Top call returns `true` ✔ — any revisit of a memoized state now returns instantly.

---

## Approach 3 — Bottom-Up DP (Reachable Jump Sets)

### Intuition

Flip the recursion into forward propagation. Let `reach[i]` be the set of last-jump sizes with which stone `i` can be reached. Seed `reach[0] = {0}` (we begin standing on stone 0 with a phantom jump of 0). For every stone `i` and every arrival jump `k`, the frog can leave with `k−1`, `k`, or `k+1` (each ≥ 1); if such a jump lands on a later stone `j`, add that jump size to `reach[j]`. The far bank is crossable exactly when `reach[last]` ends up non-empty.

### Algorithm

1. Build `indexOf`: position → stone index. Initialise `reach[i]` as empty sets.
2. Set `reach[0] = {0}`. For `i` from `0` to `n−1`, for each `k` in `reach[i]`, for `step` in `{k−1, k, k+1}` with `step ≥ 1`: if `stones[i] + step` is a stone `j` with `j > i`, add `step` to `reach[j]`.
3. Return `len(reach[n−1]) > 0`.

### Complexity

- **Time:** O(n²) — each stone accumulates at most `O(n)` distinct arrival jump sizes, and each yields 3 transitions.
- **Space:** O(n²) — the reachable-jump-size sets across all stones.

### Code

```go
func dpBottomUp(stones []int) bool {
	n := len(stones)
	indexOf := make(map[int]int, n)
	for i, s := range stones {
		indexOf[s] = i
	}

	// reach[i] holds every jump size by which stone i can be reached.
	reach := make([]map[int]bool, n)
	for i := range reach {
		reach[i] = make(map[int]bool)
	}
	reach[0][0] = true // start on stone 0, "arrived" with a phantom jump of 0

	for i := 0; i < n; i++ {
		for k := range reach[i] { // every way we could be standing on stone i
			for _, step := range []int{k - 1, k, k + 1} {
				if step <= 0 {
					continue // forward jumps only
				}
				if j, ok := indexOf[stones[i]+step]; ok && j > i {
					reach[j][step] = true // record how stone j was reached
				}
			}
		}
	}
	// The far bank is crossable iff we recorded any arrival on the last stone.
	return len(reach[n-1]) > 0
}
```

### Dry Run

Example 1: `stones = [0,1,3,5,6,8,12,17]`. Positions with their index in parentheses; each cell shows arrival jump sizes accumulated in `reach[j]` as we sweep `i` left to right.

| i (pos) | arrival k's in reach[i] | jumps emitted (step → landing stone) | updates |
|---------|-------------------------|--------------------------------------|---------|
| 0 (0) | {0} | 1 → pos 1 | reach[1] += 1 |
| 1 (1) | {1} | 2 → pos 3 | reach[2] += 2 |
| 2 (3) | {2} | 1 → (pos 4, no stone); 2 → pos 5; 3 → (pos 6=stone) | reach[3] += 2; reach[4(pos6)] += 3 |
| 3 (5) | {2} | 1 → pos 6; 3 → pos 8 | reach[4(pos6)] += 1; reach[5(pos8)] += 3 |
| 4 (6) | {3,1} | from 3: 2→pos8,4→(10 no); from 1: 2→pos8 | reach[5(pos8)] += {2} |
| 5 (8) | {3,2} | from 3: 4→pos12; from 2: … | reach[6(pos12)] += 4 |
| 6 (12) | {4} | 5 → pos 17 | reach[7(pos17)] += 5 |
| 7 (17) | {5} | last stone | — |

`reach[7]` (the last stone) = `{5}`, non-empty → `true` ✔

---

## Key Takeaways

- **When the same node is reachable in different "modes", the mode joins the DP state.** Here the mode is the last jump size `k`, turning a 1D "can I reach stone i?" into a 2D `(i, k)` state.
- **Sparse, huge coordinates ⇒ hash by position, not array index.** Positions up to `2³¹−1` forbid a positions-sized array; map position → stone index instead.
- **Top-down memo and bottom-up propagation are the same DP.** Choose top-down if you derived it from brute force; bottom-up "reachable sets" if you prefer explicit forward flow.
- **Seed the phantom start.** Modeling the initial state as "arrived at stone 0 with jump 0" cleanly forces the mandated first jump of 1 without a special case.

---

## Related Problems

- LeetCode #55 — Jump Game (reachability DP / greedy)
- LeetCode #45 — Jump Game II (min jumps, BFS-style DP)
- LeetCode #1306 — Jump Game III (reach any zero via ± jumps)
- LeetCode #1340 — Jump Game V (DP over jump constraints)
- LeetCode #70 — Climbing Stairs (1D step DP warm-up)
