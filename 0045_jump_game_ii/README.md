# 0045 — Jump Game II

> LeetCode #45 · Difficulty: Medium
> **Categories:** Array, Dynamic Programming, Greedy

---

## Problem Statement

You are given a **0-indexed** array of integers `nums` of length `n`. You are initially positioned at `nums[0]`.

Each element `nums[i]` represents the maximum length of a forward jump from index `i`. In other words, if you are at `nums[i]`, you can jump to any `nums[i + j]` where:
- `0 <= j <= nums[i]` and
- `i + j < n`

Return the **minimum number of jumps** to reach `nums[n - 1]`. The test cases are generated such that you can always reach `nums[n - 1]`.

**Example 1**
```
Input:  nums = [2,3,1,1,4]
Output: 2
Explanation: Jump 1 step from index 0 to 1, then 3 steps to the last index.
```

**Example 2**
```
Input:  nums = [2,3,0,1,4]
Output: 2
```

**Constraints**
- `1 <= nums.length <= 10⁴`
- `0 <= nums[i] <= 1000`
- The test cases are generated such that you can always reach `nums[n-1]`.

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

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Greedy / BFS Level Expansion** — the key insight: treat each "jump" as a BFS level. All positions reachable in exactly k jumps form one level. The answer is the level at which we first reach (or pass) the last index.
- **Dynamic Programming** — Approach 3 uses a dp array; O(n²) but good for understanding.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | BFS / Level Expansion ✅ | O(n) | O(1) | Intuitive; equivalent to greedy |
| 2 | Greedy ✅ | O(n) | O(1) | Cleaner code; same algorithm |
| 3 | DP | O(n²) | O(n) | Educational; clearly shows subproblem structure |

---

## Approach 1 — BFS / Level Expansion

### Intuition
Imagine the array as a graph: index 0 is the source, and from index `i` you can reach `i+1` to `i+nums[i]`. BFS finds the shortest path (in jumps). Each "level" in BFS is one jump.

Track:
- `curEnd` — the farthest index reachable with the current number of jumps.
- `farthest` — the farthest index reachable with one more jump from any index in the current level.

When we reach `curEnd` (end of current level), increment jumps and advance `curEnd = farthest`.

### Algorithm
```
jumps=0, curEnd=0, farthest=0
for i = 0 to n-2:
  farthest = max(farthest, i + nums[i])
  if i == curEnd:        // reached end of current BFS level
    jumps++
    curEnd = farthest
    if curEnd >= n-1: break
return jumps
```

### Complexity
- **Time:** O(n).
- **Space:** O(1).

### Code
```go
// bfs solves Jump Game II by treating each jump level as a BFS level.
//
// Time:  O(n)
// Space: O(1)
func bfs(nums []int) int {
    n := len(nums)
    if n <= 1 {
        return 0
    }
    jumps := 0
    curEnd := 0   // rightmost index reachable with current jump count
    farthest := 0 // rightmost index reachable with one more jump
    for i := 0; i < n-1; i++ {
        if i+nums[i] > farthest {
            farthest = i + nums[i] // update farthest reachable from i
        }
        if i == curEnd { // we've exhausted the current jump level
            jumps++
            curEnd = farthest
            if curEnd >= n-1 {
                break // already reached the end
            }
        }
    }
    return jumps
}
```

### Dry Run — `nums = [2,3,1,1,4]`
```
jumps=0, curEnd=0, farthest=0

i=0: farthest=max(0,0+2)=2. i==curEnd(0): jumps=1, curEnd=2
i=1: farthest=max(2,1+3)=4. i<curEnd(2)
i=2: farthest=max(4,2+1)=4. i==curEnd(2): jumps=2, curEnd=4. 4>=4→break

return 2 ✓
```

---

## Approach 2 — Greedy (Recommended ✅)

### Intuition
At each position, greedily track the farthest index reachable so far. When forced to use a jump (when we've exhausted the current range), increment the jump count and extend the range.

This is the same algorithm as Approach 1, framed as greedy.

### Complexity
- **Time:** O(n).
- **Space:** O(1).

### Code
```go
// greedy solves Jump Game II identically to the BFS but framed as a greedy:
// at each position, track the farthest reachable index. Use a jump whenever we
// must advance to the next "frontier".
//
// Time:  O(n)
// Space: O(1)
func greedy(nums []int) int {
    jumps, curEnd, farthest := 0, 0, 0
    for i := 0; i < len(nums)-1; i++ {
        if i+nums[i] > farthest {
            farthest = i + nums[i]
        }
        if i == curEnd {
            jumps++
            curEnd = farthest
        }
    }
    return jumps
}
```

### Dry Run — `nums = [2,3,1,1,4]`

Same mechanics as Approach 1 (no early break in this variant): update `farthest`, and whenever `i == curEnd` spend a jump and extend the frontier. Loop runs `i = 0 … n-2`.

| `i` | `nums[i]` | `farthest = max(farthest, i+nums[i])` | `i == curEnd`? | `jumps` | `curEnd` |
|-----|-----------|----------------------------------------|----------------|---------|----------|
| start | — | 0 | — | 0 | 0 |
| 0 | 2 | max(0, 2) = 2 | yes (0==0) | 1 | 2 |
| 1 | 3 | max(2, 4) = 4 | no (1<2) | 1 | 2 |
| 2 | 1 | max(4, 3) = 4 | yes (2==2) | 2 | 4 |
| 3 | 1 | max(4, 4) = 4 | no (3<4) | 2 | 4 |

Loop ends at `i = 3` (`n-1 = 4`). Return `jumps = 2` ✓.

---

## Approach 3 — Dynamic Programming

### Intuition
`dp[i]` = minimum jumps to reach index `i`. For each `i`, look back at all `j < i` where `j + nums[j] >= i`, and take `dp[j] + 1`.

### Complexity
- **Time:** O(n²).
- **Space:** O(n).

### Code
```go
// dpApproach solves Jump Game II with backward DP.
//
// dp[i] = minimum jumps needed to reach index i.
// For each i, look back at all j < i where j + nums[j] >= i:
//   dp[i] = min(dp[j] + 1).
//
// Time:  O(n²)
// Space: O(n)
func dpApproach(nums []int) int {
    n := len(nums)
    dp := make([]int, n)
    for i := range dp {
        dp[i] = 1<<31 - 1 // infinity
    }
    dp[0] = 0
    for i := 1; i < n; i++ {
        for j := 0; j < i; j++ {
            if j+nums[j] >= i && dp[j]+1 < dp[i] {
                dp[i] = dp[j] + 1
            }
        }
    }
    return dp[n-1]
}
```

### Dry Run — `nums = [2,3,1,1,4]`

`dp[i]` = min jumps to reach `i`. Init `dp = [0, ∞, ∞, ∞, ∞]`. For each `i`, scan `j < i`; if `j + nums[j] >= i`, relax `dp[i] = min(dp[i], dp[j]+1)`.

| `i` | reachable `j` (where `j+nums[j] >= i`) | `dp[j]+1` candidates | `dp[i]` |
|-----|-----------------------------------------|----------------------|---------|
| 1 | j=0 (0+2≥1) | dp[0]+1 = 1 | 1 |
| 2 | j=0 (0+2≥2), j=1 (1+3≥2) | 0+1=1, 1+1=2 | 1 |
| 3 | j=1 (1+3≥3), j=2 (2+1≥3) | dp[1]+1=2, dp[2]+1=2 | 2 |
| 4 | j=1 (1+3≥4), j=3 (3+1≥4) | dp[1]+1=2, dp[3]+1=3 | 2 |

`dp = [0, 1, 1, 2, 2]`. Return `dp[n-1] = dp[4] = 2` ✓.

---

## Key Takeaways

- **Loop to `n-2`, not `n-1`** — we process positions as sources; the last index `n-1` is the destination, not a source. Processing it would incorrectly increment `jumps`.
- **`farthest` is updated before the `if i==curEnd` check** — so when we increment `jumps` at `i==curEnd`, `farthest` already reflects the best jump from the entire current level.
- **Greedy choice is optimal** — at each "forced jump" point, taking the maximum reach is always best (no future choice can be hurt by reaching farther now, since reaching farther gives more options).
- **Difference from #55 (Jump Game I)** — #55 asks CAN you reach the end; this asks what is the MINIMUM number of jumps.

---

## Related Problems

- LeetCode #55 — Jump Game (can you reach the end; no count needed)
- LeetCode #1306 — Jump Game III (bidirectional jumps; BFS)
- LeetCode #1345 — Jump Game IV (different connections; BFS)
