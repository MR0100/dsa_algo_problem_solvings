# 0055 — Jump Game

> LeetCode #55 · Difficulty: Medium
> **Categories:** Array, Dynamic Programming, Greedy

---

## Problem Statement

You are given an integer array `nums`. You are initially positioned at the array's **first index**, and each element in the array represents your maximum jump length at that position.

Return `true` if you can reach the last index, or `false` otherwise.

**Example 1**
```
Input:  nums = [2,3,1,1,4]
Output: true
Explanation: Jump 1 step from index 0 to 1, then 3 steps to the last index.
```

**Example 2**
```
Input:  nums = [3,2,1,0,4]
Output: false
Explanation: You will always arrive at index 3 no matter what. Its maximum jump length is 0, which makes it impossible to reach the last index.
```

**Constraints**
- `1 <= nums.length <= 10⁴`
- `0 <= nums[i] <= 10⁵`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★★ Very High | 2024          |
| Meta      | ★★★★★ Very High | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★★☆ High      | 2023          |
| Apple     | ★★★☆☆ Medium    | 2023          |
| Adobe     | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Greedy (Farthest Reach)** — track the maximum index reachable; if current index exceeds it, we're stuck.
- **Dynamic Programming** — `dp[i]` = can we reach the last index from `i`? Computed right-to-left.
- **Memoization** — top-down DP caching which positions are "good" or "bad."

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Memoization (Top-Down DP) | O(n²) | O(n) | Good learning step; still quadratic |
| 2 | DP Bottom-Up | O(n²) | O(n) | Textbook DP; same complexity |
| 3 | Greedy ✅ | O(n) | O(1) | Optimal; the expected interview answer |

---

## Approach 1 — Memoization (Top-Down DP)

### Intuition
Label each index as "good" (can reach the end) or "bad" (cannot). The last index is trivially good. For any index `i`, if any reachable index `j` in `[i+1, i+nums[i]]` is good, then `i` is good.

Cache results to avoid repeated recursion on the same index.

### Complexity
- **Time:** O(n²) — n indices, each potentially scanning up to n successors.
- **Space:** O(n) — memo array + call stack.

### Code
```go
func memoization(nums []int) bool {
    n := len(nums)
    memo := make([]int, n) // 0=unknown, 1=good, 2=bad
    memo[n-1] = 1          // last index is always good

    var canJump func(i int) bool
    canJump = func(i int) bool {
        if memo[i] != 0 {
            return memo[i] == 1
        }
        maxReach := i + nums[i]
        if maxReach >= n-1 {
            memo[i] = 1
            return true
        }
        for j := i + 1; j <= maxReach; j++ {
            if canJump(j) {
                memo[i] = 1
                return true
            }
        }
        memo[i] = 2
        return false
    }

    return canJump(0)
}
```

### Dry Run — `nums = [2,3,1,1,4]` (n=5, `memo` init: index 4 = good)

`canJump(i)`: if `i+nums[i] >= n-1` mark good; else recurse on `j` in `[i+1, i+nums[i]]`.

| call | i | maxReach = i+nums[i] | check | result | memo update |
|------|---|----------------------|-------|--------|-------------|
| canJump(0) | 0 | 0+2=2 | 2 < 4 → recurse j=1,2 | true (via j=1) | memo[0]=good |
| ↳ canJump(1) | 1 | 1+3=4 | 4 >= 4 → reachable | true | memo[1]=good |

`canJump(1)` returns true immediately (maxReach 4 ≥ last index), so `canJump(0)` short-circuits and returns true. Result: true ✓

---

## Approach 2 — DP Bottom-Up

### Intuition
Same idea as memoization, computed iteratively from right to left. `dp[i]` is true if any `j` in `[i+1, i+nums[i]]` has `dp[j] = true`, or if `i+nums[i] >= n-1`.

### Complexity
- **Time:** O(n²).
- **Space:** O(n).

### Code
```go
func dpBottomUp(nums []int) bool {
    n := len(nums)
    dp := make([]bool, n)
    dp[n-1] = true // last index is reachable from itself

    for i := n - 2; i >= 0; i-- {
        maxReach := i + nums[i]
        if maxReach >= n-1 {
            dp[i] = true // can jump directly to or past the end
            continue
        }
        for j := i + 1; j <= maxReach; j++ {
            if dp[j] {
                dp[i] = true
                break
            }
        }
    }

    return dp[0]
}
```

### Dry Run — `nums = [2,3,1,1,4]` (n=5)

`dp[n-1]=true`. Iterate `i` from n-2 down to 0. `dp[i]=true` if `i+nums[i] >= n-1`, or if any `dp[j]` (`j` in `[i+1, i+nums[i]]`) is true.

| i | nums[i] | maxReach = i+nums[i] | condition | dp[i] |
|---|---------|----------------------|-----------|-------|
| 4 | 4       | —                    | init      | true  |
| 3 | 1       | 4                    | 4 >= 4 → reaches end | true |
| 2 | 1       | 3                    | dp[3]=true | true |
| 1 | 3       | 4                    | 4 >= 4 → reaches end | true |
| 0 | 2       | 2                    | dp[1] or dp[2] = true | true |

`dp[0] = true`. Result: true ✓

---

## Approach 3 — Greedy (Recommended ✅)

### Intuition
The key insight: we don't need to know exactly which path reaches the end — we only need to know **how far we can reach**.

Track `farthest` = the farthest index reachable from any position visited so far. As we scan left to right:
- If `i > farthest`: index `i` is unreachable → return false.
- Update `farthest = max(farthest, i + nums[i])`.
- If `farthest >= n-1`: we can reach the last index → return true.

### Algorithm
```
farthest = 0
for i = 0 to n-1:
  if i > farthest: return false
  farthest = max(farthest, i + nums[i])
return true
```

### Complexity
- **Time:** O(n) — single left-to-right pass.
- **Space:** O(1).

### Code
```go
func greedy(nums []int) bool {
    farthest := 0
    for i, num := range nums {
        if i > farthest { return false }
        if i+num > farthest { farthest = i+num }
        if farthest >= len(nums)-1 { return true }
    }
    return true
}
```

### Dry Run — `nums = [2,3,1,1,4]` (n=5, need to reach index 4)
```
i=0, num=2: farthest=max(0,0+2)=2. 2<4.
i=1, num=3: 1<=2 ✓. farthest=max(2,1+3)=4. 4>=4 → return true ✓
```

### Dry Run — `nums = [3,2,1,0,4]` (n=5)
```
i=0, num=3: farthest=max(0,3)=3. 3<4.
i=1, num=2: 1<=3 ✓. farthest=max(3,1+2)=3.
i=2, num=1: 2<=3 ✓. farthest=max(3,2+1)=3.
i=3, num=0: 3<=3 ✓. farthest=max(3,3+0)=3.
i=4, num=4: 4>3 → return false ✓
```

---

## Key Takeaways

- **Greedy works because reach is monotone** — extending the farthest reach is always at least as good as any other strategy. No backtracking needed.
- **`nums[i] = 0` is the trap** — if we reach a "0-jump" index and it's the only way forward, we're stuck. The greedy correctly detects this via `i > farthest`.
- **#45 Jump Game II (minimum jumps) uses the same `farthest` idea** — it adds a `curEnd` pointer to count "levels" (jumps), which is the BFS / greedy jump-count extension.
- **DP O(n²) is needed only for the jump count variant with constraints** — for this yes/no problem, O(n) greedy is sufficient.

---

## Related Problems

- LeetCode #45 — Jump Game II (minimum number of jumps; greedy with level tracking)
- LeetCode #1345 — Jump Game IV (BFS for minimum jumps with value-indexed graph)
- LeetCode #1696 — Jump Game VI (DP with sliding window max deque)
